package gateways

import (
	"context"

	"github.com/henleda/xc-ew-discovery/model"
)

// Scraper reads declared inventory from third-party API gateways: Kong, Apigee,
// Azure APIM, and AWS API Gateway. These are the gateways a different team buys
// on a different budget, invisible to anything that only watches the cluster.
//
// P0#9 requires at least two of these. Each returns observations with
// Confidence.Source set to SourceGatewayScrape. This is the declared side of the
// diff, exactly like scrape/k8s but for gateways outside the platform team.
type Scraper struct {
	// Kind selects the gateway flavor: "kong", "apigee", "apim", or "aws".
	Kind     string
	Endpoint string
}

// Scrape reads declared routes or API products from the selected gateway's admin
// or management API and returns them as observations.
//
// Declared observations carry no timestamps from traffic. LastSeen stays zero and
// the merge step uses that to classify zombies, the same contract scrape/k8s and
// scrape/datapower follow.
func (s *Scraper) Scrape(ctx context.Context) ([]model.Observation, error) {
	// TODO(M10): Kong Admin API — list routes and services.
	// TODO(M10): Apigee — list API proxies and their basepaths.
	// TODO(M10): Azure APIM and AWS API Gateway — list APIs and operations.
	// TODO(M10): resolve each to a synthetic Workload so gateway rows key the
	// same way cluster rows do, and record which gateway declared each endpoint.
	panic("not implemented: M10")
}
