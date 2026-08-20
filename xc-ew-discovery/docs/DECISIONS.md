# Decisions

Append-only. Every entry states the decision, the reasoning, and what it forecloses. Changing `model/event.go` requires an entry here written first.

## D1. Out-of-band only, no inline path

Decided. The sensor never blocks and never sits in the data path. Enforcement in v2 works by generating policy for enforcement points already in the path: Cilium NetworkPolicy, Istio AuthorizationPolicy, the customer gateway, the XC load balancer for north/south.

Forecloses: any feature requiring request modification or a synchronous decision.

Reason: the install approval story depends on it. Platform engineering approves a read-only sensor. It does not approve an inline component on a two-week timeline.

## D2. Adopt Grafana Beyla rather than write eBPF

Decided. Beyla is Apache 2, already handles Go symbol offset resolution across ABI changes, and already parses HTTP, HTTP/2, gRPC, and SQL. The fork point is the export format, not the capture logic.

Forecloses: capture behavior Beyla does not support without upstream work.

Reason: writing the sensor adds roughly two quarters and a kernel engineering hire.

## D3. Event schema frozen at M0

Decided. `model/event.go` is the contract between all six components. Changing it after M3 forces rework across every one.

The identity, verb, path template, and direction fields exist to make v2 policy synthesis possible without a schema migration. Nothing in v1 reads them. Populate them anyway.

## D4. Identity-keyed, never IP-keyed

Decided. Every output keys on namespace and service account. An unresolved observation is marked and counted rather than emitted with an IP as identity.

Reason: pod IPs churn. An inventory keyed on IPs is stale the moment it is written, and the caller graph is the thing customers pay for.

## D5. Go throughout

Decided. Beyla is Go, client-go is Go, and one language keeps `model/` genuinely shared rather than duplicated across runtimes.

## D6. DataPower is declared-inventory first

Decided. SOMA and REST management interfaces for declared services in v1. Log-target ingestion for observed calls is P1.

Reason: log format comes from each customer's own service policy, so log parsing never becomes a connector shipped once. Declared inventory needs no behavior change on the appliance and carries the demo.

## D7. Confidence on every observation

Decided. Sampling state, attribution state, and source are on the observation itself.

Reason: the UI reports coverage next to every count, because customers are being asked to trust a negative. A clean result means nothing without the share of nodes covered, the unattributed rate, and the sampling loss. Those numbers are computable only if every observation carries them.

## D8. Dev environment is one amd64 cloud VM, not the DGX Spark

Decided, superseding the earlier arm64 target.

amd64 puts Beyla on its best-tested path and removes the DataPower problem entirely, since IBM ships a free developer-edition container that runs alongside k3s on the same host and exposes SOMA and REST management.

Forecloses: nothing in v1. arm64 validation moves to a Graviton EKS node group before GA, because arm64 nodes are common in enterprise clusters.

Reason: the arm64 path bought a hardware constraint and paid for it with two workarounds and an untested uprobe surface.

## M0 result

Pending. Record the Java TLS finding here before starting M1.
