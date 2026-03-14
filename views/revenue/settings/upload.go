package settings

import (
	templateview "github.com/erniealice/hybra-golang/views/template"
	"github.com/erniealice/pyeza-golang/view"
)

// NewUploadAction creates the upload handler for invoice templates.
// Route: GET/POST /action/sales/settings/templates/upload
func NewUploadAction(deps *Deps) view.View {
	return templateview.NewUploadAction(templateConfig(deps))
}
