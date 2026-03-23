package settings

import (
	templateview "github.com/erniealice/hybra-golang/views/template"
	"github.com/erniealice/pyeza-golang/view"
)

// NewView creates the purchases settings templates list view.
// Route: GET /app/purchases/settings/templates
func NewView(deps *SettingsViewDeps) view.View {
	return templateview.NewListView(templateConfig(deps))
}

// NewUploadAction creates the upload handler for purchase order templates.
// Route: GET/POST /action/purchases/settings/templates/upload
func NewUploadAction(deps *SettingsViewDeps) view.View {
	return templateview.NewUploadAction(templateConfig(deps))
}

// NewDeleteAction creates the delete handler for purchase order templates.
// Route: POST /action/purchases/settings/templates/delete
func NewDeleteAction(deps *SettingsViewDeps) view.View {
	return templateview.NewDeleteAction(templateConfig(deps))
}

// NewSetDefaultAction creates the set-default handler for purchase order templates.
// Route: POST /action/purchases/settings/templates/set-default/{id}
func NewSetDefaultAction(deps *SettingsViewDeps) view.View {
	return templateview.NewSetDefaultAction(templateConfig(deps))
}
