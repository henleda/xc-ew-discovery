package infer

import "github.com/henleda/xc-ew-discovery/model"

// Classifier detects sensitive data and records the type and location only.
//
// There is no code path in this package that stores, logs, or returns a matched
// value. See constraint 3 in CLAUDE.md. The test suite asserts this.
type Classifier struct{}

// Scan returns hits for the given body. Field paths are structural, so
// "$.customer.ssn" is fine and the value at that path is not.
func (c *Classifier) Scan(location string, body []byte) []model.ClassifierHit {
	// TODO(M3): pattern and structural detection for common types.
	// TODO(P1): east/west specific types. Internal identifiers and bearer
	// tokens matter more here than card numbers do.
	panic("not implemented: M3")
}
