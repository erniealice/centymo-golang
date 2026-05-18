package supplier_billing_event

import "embed"

// TemplatesFS exposes the templates dir for service-admin's container to
// register at app startup, mirroring the pattern used by every other
// centymo-golang view module.
//
//go:embed templates/*.html
var TemplatesFS embed.FS
