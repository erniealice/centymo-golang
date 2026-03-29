package expenditure

import "embed"

//go:embed templates/*.html detail/templates/*.html
var TemplatesFS embed.FS
