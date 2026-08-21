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

	// CPUMillicores and MemoryMiB are the hard ceiling. The sensor sheds load by
	// sampling as it approaches them, never by growing, and reports the loss as
	// Coverage.SampleLossRate. Enforced with the Beyla DaemonSet resources.limits.
	// P0#4.
	CPUMillicores int
	MemoryMiB     int
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
	// TODO(M1): verify TLS capture on amd64, Beyla's best-tested uprobe path. See docs/ENVIRONMENT.md.
	// TODO(M7): enforce the CPU/memory ceiling by sampling, not growth, and
	// report the drop rate as Coverage.SampleLossRate. P0#4.
	panic("not implemented: M1")
}

// DryRunReport is the attach plan: what the sensor would instrument and its
// estimated footprint. It is neither an Observation nor an Inventory, so it lives
// outside the frozen model package. P0#13.
type DryRunReport struct {
	// Targets are the executables and libraries Beyla would attach uprobes to.
	Targets []string
	// EstCPUMillicores and EstMemoryMiB are the projected footprint.
	EstCPUMillicores int
	EstMemoryMiB     int
}

// DryRun reports what the sensor would attach to and an estimated resource
// footprint without loading any BPF program. This is the artifact platform
// engineering approves before install. P0#13.
func (r *Reader) DryRun(ctx context.Context) (DryRunReport, error) {
	// TODO(M8): enumerate attachable executables and libraries on the node and
	// the uprobe/kprobe targets Beyla would attach, without loading programs.
	// TODO(M8): estimate CPU and memory from process count and a traffic sample.
	panic("not implemented: M8")
}
