package sensor

import (
	"context"

	"github.com/henleda/xc-ew-discovery/model"
)

// Reader converts Beyla's export stream into model.Observation.
//
// This package wraps Beyla. It does not implement eBPF. See D2.
type Reader struct {
	Config Config
}

// Config points at the Beyla export endpoint and carries the sampling policy.
type Config struct {
	// OTLPEndpoint is where the Beyla DaemonSet ships spans.
	OTLPEndpoint string
	// SampleRate is recorded on every observation so Coverage is computable.
	SampleRate float64
	SensorID   string
	ClusterID  string
}

// Run streams observations until ctx is cancelled.
//
// M1. Every emitted observation sets Confidence.Source to SourceSensor,
// Confidence.SampleRate from Config, and Attributed to false. The enrich
// package flips Attributed once identity resolves.
func (r *Reader) Run(ctx context.Context, out chan<- model.Observation) error {
	// TODO(M1): subscribe to the Beyla OTLP stream.
	// TODO(M1): map Beyla span attributes to model.L7 for http, http2, grpc.
	// TODO(M1): carry source and destination PIDs through so enrich can resolve them.
	// TODO(M1): verify TLS capture on arm64 specifically. See docs/ENVIRONMENT.md.
	panic("not implemented: M1")
}
