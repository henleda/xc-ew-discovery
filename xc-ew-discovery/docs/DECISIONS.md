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

## D9. Authenticated flag ratified on L7 and Endpoint

Decided. `L7.Authenticated` and `Endpoint.Authenticated` stay in the frozen
schema. They were present before this entry, which D3 forbids; this entry
ratifies them rather than removing a field the product needs, and `authenticated`
is added to the SPEC 2.4 field list in the same change so the two agree.

The "unauthenticated" sort key in CLAUDE.md is a concrete, non-composite fact the
inventory sorts on. Surfacing it requires the observation to carry authentication
state, so the field is load-bearing, not speculative like the v2 policy fields.

Forecloses: nothing.

Reason: the field is needed and already shipped. The defect was process, not the
field. Recorded here so the schema, SPEC 2.4, and `/freeze-check` agree.

## D10. The freeze covers all of model/, not only event.go

Decided. `model/inventory.go` (State, Endpoint, Coverage, Inventory) is frozen on
the same terms as `model/event.go`. Changing a field in either requires an entry
here first. CLAUDE.md constraint 4 and `/freeze-check` are updated to name both.

SPEC 2.7 calls `/model` the "shared event and inventory schema, single source of
truth" and says to treat it as frozen. M4 and M5 merge into Endpoint and
Coverage, so a retype there after M3 forces the same cross-component rework D3
exists to prevent.

Forecloses: silent edits to the inventory contract.

Reason: D3 named only event.go, leaving the inventory half of the same contract
ungated. This closes that gap.

## D11. Attributed is per-peer; Confidence.Attributed is the observation roll-up

Decided. `Peer.Attributed` records whether that one peer resolved to a workload.
`Confidence.Attributed` is the observation-level roll-up and is true only when
both peers resolved: `source.Attributed && destination.Attributed`.

An observation is a caller-graph edge, and an edge needs both ends named to key on
identity (D4). A half-resolved observation counts as unattributed for coverage, so
`Confidence.Attributed` is the value `Coverage.UnattributedRate` is computed from;
the per-peer flags say which end failed.

Forecloses: reading either flag as a synonym for the other.

Reason: both flags existed with no rule for which governs, so enrich, the sensor,
and Coverage could diverge. This fixes the meaning before enrich is written.

## M0 result

Pending. Record the Java TLS finding here before starting M1.
