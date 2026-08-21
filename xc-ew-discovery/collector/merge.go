package collector

import "github.com/henleda/xc-ew-discovery/model"

// Merge folds observations from every source — sensor, k8s scrape, gateway
// scrape, DataPower — into one inventory keyed on workload identity, and
// classifies each endpoint declared-and-observed, shadow, zombie, or drift.
//
// This is the declared-versus-observed diff that carries the M4 demo. It is a
// pure function so both the M4 stdout path (`ewd inventory`) and the M6 CE run it
// unchanged. Observed observations (Confidence.Source == SourceSensor) set
// LastSeen; declared ones from a scrape leave it zero. Declared but never
// observed is a zombie; observed but never declared is shadow. See P0#11, D11,
// and model/inventory.go.
func Merge(obs []model.Observation) model.Inventory {
	// TODO(M6): deduplicate observations before folding, per SPEC 2.3.
	// TODO(M4): group observations by Owner plus L7 identity into Endpoints.
	// TODO(M4): classify shadow (observed, never declared) and zombie (declared,
	// never observed) from which Sources contributed each endpoint.
	// TODO(M4): when a declared and an observed endpoint disagree on shape, mark
	// StateDrift.
	// TODO(M4): compute Coverage from Confidence across the batch; the
	// observation-level roll-up for UnattributedRate is Confidence.Attributed (D11).
	panic("not implemented: M4")
}
