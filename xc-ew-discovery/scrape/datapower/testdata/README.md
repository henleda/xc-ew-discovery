# DataPower fixtures

DataPower Virtual Edition is amd64 only and will not run on the dev box. These fixtures let M5 proceed without an appliance.

Capture real responses from an appliance over the network and commit them here, sanitized. Needed:

- `get-config-wsgateway.xml`, SOMA get-config response for WSGateway objects
- `get-config-mpgw.xml`, SOMA get-config response for MultiProtocolGateway objects
- `wsdl-sample.xml`, a WSDL referenced by a WSP so operation resolution has something to parse
- `rest-mgmt-domains.json`, REST management interface domain listing

Sanitize hostnames, certificates, and any customer identifiers before committing.
