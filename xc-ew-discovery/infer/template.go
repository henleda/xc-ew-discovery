package infer

// Template collapses variable path segments so /users/1234/orders and
// /users/5678/orders produce one endpoint.
//
// M3. Segments collapse when the set of observed values at that position
// exceeds a cardinality threshold and the values share a shape (numeric, uuid,
// hex, base64). A segment observed with one value stays literal.
type Template struct {
	// MinObservations is the count required before a segment is eligible to
	// collapse. Below it, keep the literal and revisit.
	MinObservations int
	// CardinalityRatio is the distinct-to-total ratio above which a segment
	// collapses.
	CardinalityRatio float64
}

// Observe records a raw path for later templating.
func (t *Template) Observe(host, rawPath string) {
	// TODO(M3): tokenize and accumulate per-position value sets.
	panic("not implemented: M3")
}

// Resolve returns the templated form for a raw path, or the raw path when no
// template has been established yet.
func (t *Template) Resolve(host, rawPath string) string {
	// TODO(M3): match against established templates, longest prefix wins.
	panic("not implemented: M3")
}
