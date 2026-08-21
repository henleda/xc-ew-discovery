# XC east/west API discovery

Out-of-band API discovery for Kubernetes, third-party gateways, and IBM DataPower. Produces a merged service-to-service API inventory with inferred schemas, sensitive-data flags, and a caller graph. Ships as an add-on to F5 XC API Security.

Read `docs/SPEC.md` for requirements and `docs/MILESTONES.md` for what to build next. Work one milestone at a time. Stop at its definition of done and report, rather than continuing into the next one.

## Hard constraints

Breaking any of these breaks the product thesis, not a test.

1. **Nothing goes inline.** No proxying, no sidecar injection, no CNI change, no request modification, no blocking. The sales argument rests entirely on out-of-band posture.
2. **No request or response body leaves the cluster.** Schemas, fingerprints, and metadata only. `collector/egress.go` holds the assertion and CI enforces it.
3. **Classifier hits record type and location, never value.** No code path stores, logs, or transmits a matched sensitive value.
4. **`model/event.go` and `model/inventory.go` are frozen.** Adding, removing, renaming, or retyping a field in either requires an entry in `docs/DECISIONS.md` written first. See D3 and D10.
5. **Identity, not IP.** Every output keys on namespace and service account. An observation resolving only to an IP gets `Attributed: false` and lands in the unattributed bucket, never emitted with an IP standing in for identity.
6. **The sensor is adopted, not written.** Use Grafana Beyla. Do not write eBPF programs in this repo.

## Architecture

Seven top-level packages, six producing components. Every source writes `model.Observation`; `collector` also builds the merged `model.Inventory` (see `model/inventory.go`).

| Package | Role |
|---|---|
| `sensor/` | Wraps Beyla. Converts its export format to `model.Observation`. |
| `enrich/` | Watches kube-apiserver. Resolves socket tuple, PID, and cgroup to workload identity. |
| `infer/` | Path templating, schema inference to OpenAPI 3.1, sensitive-data classification. |
| `scrape/k8s/` | Gateway API HTTPRoute and GRPCRoute, Ingress, Istio VirtualService. Declared inventory. |
| `scrape/gateways/` | Kong, Apigee, Azure APIM, AWS API Gateway. Declared inventory. |
| `scrape/datapower/` | SOMA and REST management interfaces. Declared services. |
| `collector/` | Receives observations, merges, holds the egress boundary. |

Flow: sources produce observations, `enrich` attributes them, `infer` derives schemas, `collector` merges and emits upstream.

## Environment

Development target is one amd64 Ubuntu 24.04 cloud VM running k3s with DataPower for Developers alongside it. A managed EKS cluster gets created on demand for M6 and for node-exclusion evidence, then torn down. Read `docs/ENVIRONMENT.md` before the first build.

Do not develop this on macOS. Docker Desktop and kind run Linux in a VM whose uprobe behavior does not represent a customer node.

## Conventions

- Go 1.23. Standard library first. Add a dependency only when it removes real work.
- Every source populates `Confidence` on every observation. UI coverage numbers are computed from it, so an unpopulated `Confidence` is a bug.
- Table-driven tests. `testdata/` fixtures are committed so development works without a live appliance or cluster.
- Wrap errors with context using `fmt.Errorf("...: %w", err)`.
- No `interface{}` or `any` in the model package.

## Do not

- Do not build a global service graph as a landing view. The graph is a two-hop drill-down entered from a workload. Global force-directed maps become unreadable above roughly fifty services.
- Do not invent a composite risk score. Sort on concrete facts: undeclared, unauthenticated, sensitive data, new this week.
- Do not start P1 or P2 work while a P0 milestone is open. `docs/SPEC.md` carries the priority of every requirement.
- Do not add top-level packages. The top-level set is `model`, `sensor`, `enrich`, `infer`, `scrape`, `collector`, `cmd`. Add subpackages instead.
- Do not put Java TLS capture into the main path before M0 resolves. It is P1 and gated on a spike.
