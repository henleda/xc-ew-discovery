package model

// FROZEN. inventory.go is part of the shared contract and is frozen on the same
// terms as event.go: no field added, removed, renamed, or retyped without a
// docs/DECISIONS.md entry first. See D3 and D10.

import "time"

// State is the declared-versus-observed classification. This drives the
// landing view and it is the number the purchase gets made on.
type State string

const (
	StateDeclaredAndObserved State = "declared_and_observed"
	StateShadow              State = "shadow" // observed, never declared
	StateZombie              State = "zombie" // declared, never observed
	StateDrift               State = "drift"  // declared and observed, shapes disagree
)

// Endpoint is one row in the inventory. Grouped on Owner plus L7 identity, not
// on IP or pod.
type Endpoint struct {
	Key   string   `json:"key"`
	Owner Workload `json:"owner"`

	Method       string `json:"method,omitempty"`
	PathTemplate string `json:"path_template,omitempty"`
	GRPCService  string `json:"grpc_service,omitempty"`
	GRPCMethod   string `json:"grpc_method,omitempty"`
	SOAPAction   string `json:"soap_action,omitempty"`
	Protocol     string `json:"protocol"`

	Direction Direction `json:"direction"`
	State     State     `json:"state"`

	DeclaredBy []Source   `json:"declared_by,omitempty"`
	FirstSeen  time.Time  `json:"first_seen"`
	LastSeen   time.Time  `json:"last_seen"`
	Callers    []Workload `json:"callers"`

	Sensitivity   []ClassifierHit `json:"sensitivity,omitempty"`
	Authenticated bool            `json:"authenticated"`

	// SpecPath points at the generated OpenAPI document. The document lives on
	// disk, never inline in the inventory.
	SpecPath string `json:"spec_path,omitempty"`
}

// Coverage is the trust mechanic. Every count in the UI renders next to these
// numbers, because the product asks an analyst to trust a negative. See D7.
type Coverage struct {
	NodesTotal        int       `json:"nodes_total"`
	NodesInstrumented int       `json:"nodes_instrumented"`
	NodesExcluded     int       `json:"nodes_excluded"`
	ExclusionReasons  []string  `json:"exclusion_reasons,omitempty"`
	UnattributedRate  float64   `json:"unattributed_rate"`
	SampleLossRate    float64   `json:"sample_loss_rate"`
	WindowStart       time.Time `json:"window_start"`
	WindowEnd         time.Time `json:"window_end"`
}

// Inventory is the merged result. North/south and east/west live in one list
// separated by Direction, never in two lists. See CLAUDE.md "Do not".
type Inventory struct {
	Endpoints []Endpoint `json:"endpoints"`
	Coverage  Coverage   `json:"coverage"`
}
