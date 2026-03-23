package settings

import (
	templateview "github.com/erniealice/hybra-golang/views/template"
	"github.com/erniealice/pyeza-golang/view"
)

// NewDeleteAction creates the delete handler for invoice templates.
// Route: POST /action/sales/settings/templates/delete
func NewDeleteAction(deps *SettingsViewDeps) view.View {
	return templateview.NewDeleteAction(templateConfig(deps))
}
