package settings

import (
	templateview "github.com/erniealice/hybra-golang/views/template"
	"github.com/erniealice/pyeza-golang/view"
)

// NewView creates the sales settings templates list view.
// Route: GET /app/sales/settings/templates
func NewView(deps *Deps) view.View {
	return templateview.NewListView(templateConfig(deps))
}
