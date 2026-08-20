package collector

import (
	"errors"

	"github.com/henleda/xc-ew-discovery/model"
)

// ErrBodyLeak is returned when an outbound payload carries anything beyond
// schemas, fingerprints, and metadata.
var ErrBodyLeak = errors.New("egress payload contains body content")

// AssertNoBodies is the egress boundary. Every outbound batch passes through
// it, and CI runs it against fixtures containing known payload markers.
//
// Constraint 2 in CLAUDE.md is enforced here. Do not add a bypass, a debug
// flag, or a verbose mode that skips this.
func AssertNoBodies(batch []model.Observation) error {
	// TODO(M6): reflect over the batch and reject any string field carrying
	// content beyond the declared metadata set.
	// TODO(M6): assert ClassifierHit has no value-bearing field.
	panic("not implemented: M6")
}
