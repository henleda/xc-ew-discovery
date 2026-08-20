package infer

import "github.com/henleda/xc-ew-discovery/model"

// Inferrer builds OpenAPI 3.1 documents from observed traffic.
//
// M3. Inference runs inside the cluster. Bodies are read here and discarded
// here. Only the derived shape and its fingerprint leave this package. See
// constraint 2 in CLAUDE.md.
type Inferrer struct {
	Templater *Template
}

// Ingest folds one observation's payload shape into the model for its endpoint.
//
// The caller passes the body. This function must not retain it, log it, or
// attach it to anything returned.
func (i *Inferrer) Ingest(o model.Observation, requestBody, responseBody []byte) error {
	// TODO(M3): derive JSON shape, merge into the endpoint's accumulated schema.
	// TODO(M3): compute a stable fingerprint over the shape, not the values.
	// TODO(M3): for grpc without proto descriptors, record the binary shape and
	// mark the schema incomplete rather than guessing.
	panic("not implemented: M3")
}

// WriteSpecs emits one OpenAPI 3.1 document per discovered service into dir.
func (i *Inferrer) WriteSpecs(dir string) error {
	// TODO(M3): emit and validate against the OpenAPI 3.1 schema before writing.
	panic("not implemented: M3")
}
