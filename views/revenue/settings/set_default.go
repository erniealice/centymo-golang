package settings

import (
	templateview "github.com/erniealice/hybra-golang/views/template"
	"github.com/erniealice/pyeza-golang/view"
)

// NewSetDefaultAction creates the set-default handler for invoice templates.
// Route: POST /action/sales/settings/templates/set-default/{id}
func NewSetDefaultAction(deps *Deps) view.View {
	return templateview.NewSetDefaultAction(templateConfig(deps))
}
