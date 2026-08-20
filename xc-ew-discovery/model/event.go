// Package model holds the contract every component in this repo depends on.
//
// FROZEN. Do not add, remove, rename, or retype a field without writing an
// entry in docs/DECISIONS.md first. See D3.
package model

import "time"

// Direction classifies an observation relative to the cluster boundary.
type Direction string

const (
	DirectionEastWest   Direction = "east_west"
	DirectionNorthSouth Direction = "north_south"
	DirectionEgress     Direction = "egress"
)

// Source identifies which component produced an observation. Coverage
// reporting in the UI groups by this field.
type Source string

const (
	SourceSensor        Source = "sensor"
	SourceK8sScrape     Source = "k8s_scrape"
	SourceGatewayScrape Source = "gateway_scrape"
	SourceDataPowerSOMA Source = "datapower_soma"
	SourceDataPowerLog  Source = "datapower_log"
)

// TLSState records what the sensor saw on the wire.
type TLSState string

const (
	TLSPlaintext TLSState = "plaintext"
	TLSOneWay    TLSState = "tls"
	TLSMutual    TLSState = "mtls"
)

// Observation is the single unit of output from every source: the eBPF
// sensor, the control-plane scrapers, and the DataPower connector.
//
// Fields under Source, Destination, L7, and Direction are what a policy
// synthesizer in v2 needs to emit a Cilium NetworkPolicy or an Istio
// AuthorizationPolicy. Nothing in v1 reads them for that purpose. Populate
// them anyway. Retrofitting identity into an observation-shaped schema is a
// rewrite, not a migration.
type Observation struct {
	ID        string    `json:"observation_id"`
	Timestamp time.Time `json:"timestamp"`
	SensorID  string    `json:"sensor_id"`
	ClusterID string    `json:"cluster_id"`

	Source      Peer `json:"source"`
	Destination Peer `json:"destination"`

	Transport   Transport       `json:"transport"`
	L7          L7              `json:"l7"`
	Schema      SchemaRef       `json:"schema"`
	Sensitivity []ClassifierHit `json:"sensitivity"`

	Direction  Direction  `json:"direction"`
	Confidence Confidence `json:"confidence"`
}

// Peer is one end of an observed call, keyed on workload identity.
//
// Attributed is false when enrichment failed to resolve the workload. Such
// observations are counted in the unattributed bucket and never emitted with
// ClusterIP standing in for identity. See D4.
type Peer struct {
	Workload   Workload `json:"workload"`
	Service    string   `json:"service,omitempty"`
	Port       int32    `json:"port,omitempty"`
	ClusterIP  string   `json:"cluster_ip,omitempty"`
	Attributed bool     `json:"attributed"`
}

// Workload is the identity a policy would name.
type Workload struct {
	Namespace      string `json:"namespace,omitempty"`
	ServiceAccount string `json:"service_account,omitempty"`
	Kind           string `json:"workload_kind,omitempty"`
	Name           string `json:"workload_name,omitempty"`
	Pod            string `json:"pod,omitempty"`
	Node           string `json:"node,omitempty"`
	ContainerImage string `json:"container_image,omitempty"`
	PID            int32  `json:"pid,omitempty"`
}

// Transport records protocol and encryption state.
type Transport struct {
	Protocol         string   `json:"protocol"`
	TLSState         TLSState `json:"tls_state"`
	MTLSPeerIdentity string   `json:"mtls_peer_identity,omitempty"`
}

// L7 holds application-layer facts. PathTemplate is the collapsed form and is
// the key the inventory groups on. PathRaw stays for debugging and never
// reaches the inventory.
type L7 struct {
	Method       string `json:"method,omitempty"`
	Host         string `json:"host,omitempty"`
	PathRaw      string `json:"path_raw,omitempty"`
	PathTemplate string `json:"path_template,omitempty"`
	Status       int    `json:"status,omitempty"`
	ContentType  string `json:"content_type,omitempty"`
	GRPCService  string `json:"grpc_service,omitempty"`
	GRPCMethod   string `json:"grpc_method,omitempty"`
	SOAPAction   string `json:"soap_action,omitempty"`
	Authenticated bool  `json:"authenticated"`
}

// SchemaRef carries shape hashes, never payloads. See constraint 2 in CLAUDE.md.
type SchemaRef struct {
	RequestFingerprint  string `json:"request_fingerprint,omitempty"`
	ResponseFingerprint string `json:"response_fingerprint,omitempty"`
}

// ClassifierHit records that sensitive data was detected and where.
//
// There is no Value field and there will never be one. See constraint 3 in
// CLAUDE.md.
type ClassifierHit struct {
	Type      string `json:"type"`
	Location  string `json:"location"`
	FieldPath string `json:"field_path,omitempty"`
}

// Confidence is what the UI reports next to every count. See D7.
//
// Every source populates this on every observation. An unpopulated Confidence
// is a bug, not a default.
type Confidence struct {
	Source     Source  `json:"source"`
	Sampled    bool    `json:"sampled"`
	SampleRate float64 `json:"sample_rate"`
	Attributed bool    `json:"attributed"`
}
