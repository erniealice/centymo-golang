package recovery_document

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the compose descriptor for the recovery document pages. The
// unit is mounted by centymo's RecoveryDocumentUnit only when the app opts in
// (block.WithRecoveryCharges); it is not part of the default centymo unit set.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "revenue.recovery_document",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "recovery_document"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "recovery_document.json", Key: "recovery_document"},
		LabelName: "RecoveryDocumentLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "recovery_document:list",
			Items: []compose.NavItem{
				{Key: "recovery-documents", Route: "recovery_document.list", Params: map[string]string{"status": "issued"},
					Label: "Recovery Documents", Icon: "icon-file-text", Permission: "recovery_document:list",
					LabelKey: "recovery_documents_label", IconKey: "recovery_documents_icon"},
			},
		},
	}
}
