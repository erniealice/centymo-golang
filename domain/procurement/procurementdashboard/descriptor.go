package procurementdashboard

import "github.com/erniealice/pyeza-golang/compose"

// Describe returns the composition-v2 descriptor for the procurement dashboard.
// This is a composition-surface (no proto entity) so it has no Labels DefaultLabels
// factory yet — LabelJSON is left empty until one is added.
func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "procurement.procurementdashboard",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "procurement"},
		Templates: TemplatesFS,
	}
}
