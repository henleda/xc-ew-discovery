package k8s

import (
	"context"

	"github.com/henleda/xc-ew-discovery/model"
)

// Scraper reads declared inventory from Kubernetes objects. This is the
// declared side of the diff that carries the demo. See M4.
type Scraper struct {
	ClusterID string
}

// Scrape reads HTTPRoute, GRPCRoute, and Ingress and returns them as
// observations with Confidence.Source set to SourceK8sScrape.
//
// Declared observations carry no timestamps from traffic. LastSeen stays zero
// and the merge step uses that to classify zombies.
func (s *Scraper) Scrape(ctx context.Context) ([]model.Observation, error) {
	// TODO(M4): list HTTPRoute and GRPCRoute from gateway.networking.k8s.io.
	// TODO(M4): list Ingress from networking.k8s.io.
	// TODO(M4): resolve backendRefs to the owning workload so declared and
	// observed rows key the same way.
	panic("not implemented: M4")
}
