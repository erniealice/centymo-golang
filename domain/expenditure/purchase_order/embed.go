package purchaseorder

import "embed"

// TemplatesFS embeds HTML templates for the purchase order views.
// Templates should be placed in the templates/ subdirectory.
//
//go:embed templates/*.html
var TemplatesFS embed.FS
