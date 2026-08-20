# Milestones

Ordered to resolve the largest unknown first and reach a demo before the hard integration work. Work one milestone per session. Stop at the definition of done and report.

## M0. Java TLS spike. Two days. Do this before writing anything else.

Two Spring Boot services calling each other over mTLS on the k3s cluster. Run Beyla against them. Document exactly what is captured and what is missing.

Work in `hack/m0-java-spike/`. This is throwaway code and it does not follow repo conventions.

**Definition of done:** `docs/DECISIONS.md` carries an entry stating whether Beyla reaches SunJSSE payloads on arm64, what the workaround costs if it does not, and whether Java capture stays P1 or becomes a blocker for financial services accounts.

## M1. Sensor to stdout

Beyla DaemonSet on k3s with the polyglot demo app. Confirm HTTP, HTTP/2, and gRPC capture in both plaintext and TLS. Convert Beyla's export format into `model.Observation` in `sensor/`.

**Definition of done:** `ewd sensor --stdout` prints valid observations with source and destination PIDs populated, for all three protocols, plaintext and TLS. `Confidence.Source` is `SourceSensor` on every one.

## M2. Enrichment

kube-apiserver watch in `enrich/`. Resolve socket tuple, PID, and cgroup to namespace, service account, workload kind, and workload name.

**Definition of done:** `ewd sensor --graph` emits a caller graph as JSON containing no IP addresses. Observations failing to resolve appear with `Attributed: false` and are counted in a summary line. Unattributed rate under 10 percent on the demo app.

## M3. Inference

Path templating and OpenAPI 3.1 schema inference in `infer/`, running locally against the observation stream.

**Definition of done:** `ewd infer --out ./specs` writes one valid OpenAPI 3.1 document per discovered service. `/users/1234/orders` and `/users/5678/orders` collapse to one templated endpoint. Documents validate against the OpenAPI schema.

## M4. Declared versus observed. This is the demo.

Scrape Gateway API HTTPRoute, GRPCRoute, and Ingress objects in `scrape/k8s/`. Diff against M3 output.

**Definition of done:** `ewd inventory` prints shadow endpoints (observed, never declared) and zombie endpoints (declared, never observed) as separate sections with counts. Deleting an HTTPRoute moves its endpoint from declared to shadow on the next run.

## M5. DataPower

SOMA scrape in `scrape/datapower/` against a DataPower Gateway for Developers container running on the dev VM. Read WSP and MPGW declared services. Merge into the M4 inventory model. Commit sanitized fixtures as you go so tests run without the container.

**Definition of done:** a SOAP service from DataPower appears in `ewd inventory` output alongside a gRPC service from the cluster, in one list, keyed the same way.

## M6. CE integration

Replace the stdout collector with a CE running in-cluster. Confirm schema-only egress.

**Definition of done:** a packet capture on the CE uplink during a run containing known payload markers shows none of those markers in the captured traffic. The assertion in `collector/egress.go` runs in CI.
