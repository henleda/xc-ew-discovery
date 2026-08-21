# DataPower fixtures

DataPower Gateway for Developers is amd64 and runs on the amd64 dev VM alongside k3s (D8). These fixtures let M5's tests run without the container.

Capture real responses from an appliance over the network and commit them here, sanitized. Needed:

- `get-config-wsgateway.xml`, SOMA get-config response for WSGateway objects
- `get-config-mpgw.xml`, SOMA get-config response for MultiProtocolGateway objects
- `wsdl-sample.xml`, a WSDL referenced by a WSP so operation resolution has something to parse
- `rest-mgmt-domains.json`, REST management interface domain listing

Sanitize hostnames, certificates, and any customer identifiers before committing.
