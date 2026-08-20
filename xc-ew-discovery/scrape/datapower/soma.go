package datapower

import (
	"context"

	"github.com/henleda/xc-ew-discovery/model"
)

// Client reads declared services from a DataPower appliance over SOMA and the
// REST management interface.
//
// Declared inventory only in v1. Log-target ingestion for observed calls is
// P1, because log format comes from each customer's own service policy and
// never becomes a connector shipped once. See D6.
//
// Wire capture is not available on this platform and is not attempted.
type Client struct {
	// Endpoint is the appliance management interface, typically port 5550 for
	// SOMA or 5554 for REST.
	Endpoint string
	Domain   string
}

// Scrape returns declared WSP and MPGW services as observations with
// Confidence.Source set to SourceDataPowerSOMA.
//
// M5. The appliance is amd64 only and will not run on the dev box. Develop
// against testdata fixtures and validate over the network against a real
// appliance before closing the milestone. See docs/ENVIRONMENT.md.
func (c *Client) Scrape(ctx context.Context) ([]model.Observation, error) {
	// TODO(M5): get-config for WSGateway and MultiProtocolGateway objects.
	// TODO(M5): resolve WSDL references to operations so a SOAP action becomes
	// an endpoint row rather than a service row.
	// TODO(M5): map each service to a synthetic Workload so DataPower rows key
	// the same way cluster rows do. Namespace becomes the appliance domain.
	panic("not implemented: M5")
}
