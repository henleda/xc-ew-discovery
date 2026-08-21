# XC East/West API Discovery

**Add-on proposal and product spec**

Owner: Dan Henley, VP Product Marketing, XC
Status: Draft for internal socialization and prototype build
Date: August 2026

---

# Part 1: The Case

## 1.1 Problem

XC API Security sees what crosses the load balancer. Once traffic enters the cluster or the ESB, visibility stops. Customers pay for an API inventory covering the minority of their real API surface, with no clean number for what fraction is missing.

Three gaps show up in every large account:

1. Service-to-service calls inside Kubernetes, invisible unless a mesh with L7 telemetry is already deployed and instrumented.
2. Legacy SOAP and MQ services behind DataPower, declared in config nobody has read in years and called by systems nobody owns.
3. Internal APIs exposed through third-party gateways bought by a different team on a different budget.

Regulated accounts feel this first. Examiners ask for a complete API inventory. AppSec produces the north/south list and a shrug.

The cost of leaving it alone is category position. Akamai bought Noname. Harness bought Traceable. Cisco holds Isovalent plus Panoptica. Wiz, Upwind, and Sweet attack runtime visibility from the CNAPP side. F5 sells the device in front of these environments and sees the least of what happens inside them.

## 1.2 What it is

An out-of-band eBPF node sensor, plus control-plane scrapers, plus a DataPower connector, feeding a CE running in-cluster as collector and inference point. Output is a service-to-service API inventory with inferred schemas, sensitive-data flags, and a caller graph, merged into the existing XC API Security inventory alongside north/south data.

Nothing sits in the data path. Nothing blocks in v1. Payload bodies never leave the cluster.

## 1.3 Right to win

The CE fleet is already in these environments and already carries tenant identity, mTLS transport to the RE, RBAC, and fleet lifecycle. Competitors ship a Helm chart and then spend a year rebuilding those four things.

The install base overlap is exact. Accounts running DataPower, a service mesh, and XC load balancers are the same accounts.

North/south and east/west in one inventory is a claim no pure-play makes. Salt and Noname come from the edge inward. Wiz and Upwind come from the workload outward. Nobody owns both sides of the load balancer.

The honest counter: F5 has no kernel engineering bench, and an eBPF sensor is not proxy work. Section 2 resolves this by adopting rather than building.

## 1.4 Buyer and approval path

AppSec buys. Platform engineering approves. Two approvers, and the second has never heard of the product.

This converts directly into requirements. The sensor must be read-only, out of the data path, resource-capped, observable, and removable without residue. A platform engineer must reach that conclusion from the docs in ten minutes without a call. Treat the install story as a P0 product requirement rather than a field enablement problem.

## 1.5 Packaging and price

**Recommendation: metered add-on attached to the API Security SKU, priced per instrumented node, with a per-cluster floor.**

Per node matches how platform teams already budget for Datadog, Wiz, and Cilium Enterprise, and it ties cost to sensor footprint, which is the variable the customer controls.

Per discovered endpoint is how the API security category prices, and it is wrong here. East/west endpoint counts run five to fifty times north/south counts. The first invoice becomes a churn event.

Attached rather than standalone keeps one paper, one renewal, and one buyer of record, and it makes east/west an expansion motion instead of a new-logo motion. Requiring an active API Security entitlement also guards against the add-on cannibalizing the SKU it extends.

## 1.6 Competitive frame

| Vendor | Comes from | East/west method | Gap to exploit |
|---|---|---|---|
| Salt, Noname (Akamai) | North/south edge | Mirror or gateway integration | Weak inside the cluster, no appliance story |
| Traceable (Harness) | APM lineage | Agent plus mesh | Requires app instrumentation buy-in |
| Wiz, Upwind, Sweet | CNAPP runtime | eBPF sensor | Strong sensor, no API schema depth, no north/south |
| Cisco (Isovalent, Panoptica) | CNI | eBPF native | Best sensor in market, weakest API security product |
| Kong, Apigee | Gateway | Own traffic only | Blind to anything not through them |

Nobody covers DataPower. It is unglamorous and it sits in the accounts with the largest budgets.

## 1.7 What would kill it

**Java TLS.** The JVM terminates TLS in SunJSSE, not OpenSSL, so uprobes on libssl miss every Java-to-Java mTLS call. The target accounts are the most Java-heavy accounts in existence. This is the highest-risk unknown in the entire program and it gets spiked before anything else gets built.

**Managed Kubernetes exclusions.** EKS Fargate refuses privileged DaemonSets outright. GKE Autopilot constrains them. Quantify the excluded share of the target fleet before committing a number to a roadmap.

**Kernel floor.** BTF and CO-RE put the floor around RHEL 8.6 and newer. Older enterprise fleets drop out, and those fleets overlap with the DataPower install base.

**Internal SKU boundary.** Settling whether this extends API Security or competes with it will take longer than the engineering work.

## 1.8 Success measures

Leading, measured 30 to 60 days after first install:
- Time from Helm install to first populated inventory, target under 30 minutes
- Ratio of east/west endpoints found to north/south endpoints already known, target 5x or better
- Share of pilot accounts where platform engineering approved the sensor without a security review escalation, target 70 percent

Lagging, measured over two renewal cycles:
- Attach rate to API Security renewals
- Win rate in deals with a CNAPP vendor also bidding
- Expansion ARR per attached account

---

# Part 2: The Spec

## 2.1 Goals

1. Produce a merged API inventory spanning north/south and east/west in one view, keyed by workload identity rather than IP.
2. Reach a populated inventory within 30 minutes of install with no application changes, no sidecar injection, and no data-path change.
3. Diff declared inventory against observed traffic to surface shadow endpoints (observed, never declared) and zombie endpoints (declared, never observed).
4. Keep payload bodies inside the customer boundary. Ship schemas, fingerprints, and metadata only.
5. Emit observations in a shape sufficient to synthesize enforcement policy later, without shipping enforcement in v1.

## 2.2 Non-goals for v1

- No inline proxying, blocking, or request modification. The posture is out-of-band and the sales argument depends on holding it.
- No policy generation or enforcement actions. The event schema supports it. The product does not do it.
- No Java TLS capture until the spike in 2.7 resolves. Scoping it in before the answer is known puts the roadmap on an unowned dependency.
- No Kafka, AMQP, GraphQL, or database protocol parsing. HTTP, HTTP/2, gRPC, and SOAP over HTTP only.
- No agentless mode. Anything without kernel access is a control-plane scrape, covered separately.
- No support for EKS Fargate or GKE Autopilot. Document the exclusion rather than working around it.

## 2.3 Architecture

Six components across four planes.

**Sensor plane.** Privileged DaemonSet on each node. Socket-layer kprobes and tracepoints for plaintext. Uprobes on libssl, libcrypto, and Go crypto/tls for encrypted traffic. Emits raw L7 events to the in-cluster collector over Beyla's OTLP export. **Adopt Grafana Beyla rather than build.** Apache 2, already handles Go symbol offset resolution across ABI changes, already parses HTTP, HTTP/2, gRPC, and SQL. The fork point is the export format, not the capture logic. Building this in-house adds two quarters and a kernel engineering hire.

**Enrichment plane.** Watches kube-apiserver. Maps socket tuple, PID, and cgroup to pod, namespace, service account, and workload owner. This component turns an endpoint list into a caller graph, which is the difference between a report and a product.

**Collector plane.** CE running in-cluster. Receives sensor events, performs path templating, schema inference, sensitive-data classification, and deduplication locally. Emits schemas and metadata upstream. Never emits bodies.

**Ingest plane.** Existing XC API Security pipeline in the RE, extended with a direction field and a workload identity field so east/west and north/south records merge into one inventory rather than two.

**Control-plane scrapers.** Read declared inventory from Kubernetes Gateway API HTTPRoute and GRPCRoute, Ingress objects, Istio VirtualService, Kong Admin API, Apigee, Azure APIM, and AWS API Gateway. Produces the declared side of the diff.

**DataPower connector.** SOMA and REST management interfaces for declared services: WSPs, multi-protocol gateways, and API Connect definitions. Log-target ingestion for observed calls, second and separately. Wire capture is not available on this platform and is not attempted.

## 2.4 Event schema

Every observation carries these fields. The identity, verb, path template, and direction fields are the minimum needed to synthesize a Cilium NetworkPolicy or an Istio AuthorizationPolicy later. Emit them from day one even though no v1 feature reads them. Retrofitting identity into an observation-shaped schema is a rewrite.

```
observation_id, timestamp, sensor_id, cluster_id

source:      namespace, service_account, workload_kind, workload_name,
             pod, node, container_image, pid, attributed
destination: namespace, service_account, workload_kind, workload_name,
             pod, node, container_image, pid,
             service, port, cluster_ip, attributed
transport:   protocol, tls_state, mtls_peer_identity
l7:          method, host, path_raw, path_template, status,
             content_type, grpc_service, grpc_method, soap_action,
             authenticated
schema:      request_fingerprint, response_fingerprint
sensitivity: classifier_hits[] (type and location, never value)
direction:   east_west | north_south | egress
confidence:  source, sampled, sample_rate, attributed   (see D7)
```

## 2.5 Requirements

### P0

1. **Sensor deploys as a single Helm chart** with no application restart, no sidecar injection, and no CNI change. Given a running cluster, when the chart installs, then existing workloads continue serving with no restart recorded.
2. **Sensor captures plaintext HTTP, HTTP/2, and gRPC** at the socket layer with source and destination PIDs attached.
3. **Sensor captures TLS payloads** via uprobes on OpenSSL, BoringSSL, and Go crypto/tls, including stripped Go binaries.
4. **Sensor enforces a hard resource cap** with a configurable CPU and memory ceiling, and sheds load by sampling rather than by growing. Given sustained traffic above capacity, when the cap is reached, then the sensor drops events and reports the drop rate rather than consuming node resources.
5. **Enrichment resolves every observation to workload identity** within one reconciliation interval. Observations resolving only to an IP are marked unattributed and counted.
6. **Schema inference runs locally** and emits OpenAPI 3.1 per discovered service.
7. **Path templating collapses variable segments** so `/users/1234/orders` and `/users/5678/orders` produce one endpoint.
8. **No request or response body leaves the cluster.** Verified by an egress capture test in CI.
9. **Control-plane scrapers ingest declared inventory** from Gateway API, Ingress, and at least two third-party gateways.
10. **DataPower connector reads declared services** through the management interface and merges them into the same inventory model.
11. **Declared versus observed diff** surfaces shadow and zombie endpoints as first-class objects in the inventory.
12. **Uninstall removes all sensor artifacts** including pinned BPF maps and programs. Verified by a clean-node assertion.
13. **Dry-run mode** reports what the sensor would attach to, resource estimate included, without loading programs. This is the artifact platform engineering approves.

### P1

14. Java TLS capture, gated on the spike in 2.7.
15. DataPower log-target ingestion for observed calls, with a configurable format parser.
16. Sensitive-data classification tuned for east/west patterns, including internal identifiers and tokens.
17. Sensor supports non-Kubernetes Linux hosts for VM-based east/west traffic.

### P2

18. Policy synthesis from observations, targeting Cilium NetworkPolicy and Istio AuthorizationPolicy.
19. Kafka and database protocol parsing.
20. Mesh-native mode reading Envoy tap output where a mesh already exists, removing the need for node access.

## 2.6 Open questions

| Question | Owner | Blocking |
|---|---|---|
| Does Beyla capture Java-to-Java mTLS, and at what cost | Engineering | Yes, for FSI accounts |
| Beyla versus Pixie versus build | Engineering | Yes, gates the build sequence |
| Does the RE ingest schema accept direction and identity, or does east/west need a separate index | Product and engineering | Yes |
| Is the in-cluster collector a full CE or a lighter component, and what is its footprint | Engineering | Yes, drives the approval story |
| SKU boundary against API Security | Product and pricing | No, not for prototype |
| Does schema-only egress satisfy EU and Japan residency requirements | Legal | No, not for prototype |

## 2.7 Prototype build sequence

Ordered to resolve the largest unknown first and reach a demo before the hard integration work.

**M0. Java spike. Two days. Do this before writing anything else.**
Two Spring Boot services calling each other over mTLS on the k3s cluster. Run Beyla against them. Document exactly what is captured and what is missing. This is the go/no-go for the FSI story and it determines whether the roadmap carries a JVMTI dependency.

**M1. Sensor to stdout.**
Beyla DaemonSet on k3s with a polyglot demo app. Confirm HTTP, HTTP/2, and gRPC capture in both plaintext and TLS. Definition of done: raw L7 events printing with PIDs attached.

**M2. Enrichment.**
Add a kube-apiserver watch. Resolve events to namespace, service account, and workload. Definition of done: a caller graph emitted as JSON with no IPs in it.

**M3. Inference.**
Path templating and schema inference locally. Definition of done: valid OpenAPI 3.1 per service, generated from observed traffic alone.

**M4. Declared versus observed. This is the demo.**
Scrape Gateway API and Ingress objects. Diff against M3 output. Definition of done: a screen showing shadow endpoints and zombie endpoints side by side.

**M5. DataPower.**
SOMA scrape against DataPower Gateway for Developers for declared WSP and MPGW services. Merge into the M4 inventory model. Definition of done: a SOAP service appearing in the same inventory as a gRPC service.

**M6. CE integration.**
Replace the stdout collector with a CE in-cluster. Confirm schema-only egress with a packet capture on the CE uplink.

**Install-approval and coverage work.** Three P0 requirements sit outside this demo-first sequence: the resource cap and drop reporting (P0#4), dry-run install preview (P0#13), and uninstall with a clean-node assertion (P0#12). The two third-party gateways P0#9 requires (Kong, Apigee, APIM, AWS) and Istio VirtualService also land after the demo. These are tracked as M7 through M10 in `docs/MILESTONES.md`. They gate GA, not the demo, and none depends on M0.

### Suggested repo layout

```
/sensor        Beyla fork or wrapper, export format only
/enrich        kube-apiserver watch, identity resolution
/infer         path templating, schema inference, classification
/scrape
  /k8s         Gateway API, Ingress, Istio VirtualService
  /gateways    Kong, Apigee, APIM, AWS
  /datapower   SOMA and REST management client
/model         shared event and inventory schema, single source of truth
/collector     CE integration and egress
/demo          k3s manifests, polyglot app, DataPower Gateway for Developers
```

Build `/model` first and treat it as frozen — both `event.go` and `inventory.go`, per D3 and D10. Every other component depends on the schema in `/model`, and changing it after M3 forces rework across all six.
