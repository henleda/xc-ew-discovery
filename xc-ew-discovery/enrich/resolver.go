package enrich

import (
	"context"

	"github.com/henleda/xc-ew-discovery/model"
)

// Resolver turns a socket tuple, PID, and cgroup path into workload identity.
//
// This is the component that makes an endpoint list into a caller graph, which
// is the difference between a report and a product.
type Resolver struct {
	// Unattributed counts observations that failed to resolve. This number
	// reaches the UI. Never drop it and never fall back to an IP. See D4.
	Unattributed uint64
	Total        uint64
}

// Start begins the kube-apiserver watch and keeps the local index warm.
//
// M2. Index pods by IP and by cgroup path. Both change on restart, so the
// watch is authoritative and any cache is best-effort.
func (r *Resolver) Start(ctx context.Context) error {
	// TODO(M2): informer over Pods, Services, and EndpointSlices.
	// TODO(M2): build cgroup path index. k3s uses containerd. See docs/ENVIRONMENT.md.
	// TODO(M2): resolve workload owner by walking ownerReferences to the top controller.
	panic("not implemented: M2")
}

// Attribute fills Source and Destination workload identity on o.
//
// On failure it sets Attributed false on the unresolved peer, increments the
// counter, and returns the observation unchanged otherwise. It never writes an
// IP into a Workload field.
func (r *Resolver) Attribute(o *model.Observation) {
	// TODO(M2): resolve both peers.
	panic("not implemented: M2")
}

// Rate returns the unattributed share for Coverage reporting.
func (r *Resolver) Rate() float64 {
	if r.Total == 0 {
		return 0
	}
	return float64(r.Unattributed) / float64(r.Total)
}
