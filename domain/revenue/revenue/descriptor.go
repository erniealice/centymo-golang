package revenue

import "github.com/erniealice/pyeza-golang/compose"

// Describe returns the composition-v2 descriptor for the revenue entity.
// Labels are not yet exposed via DefaultLabels() — the LabelJSON binding is
// left empty until a DefaultLabels factory is added.
func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "revenue.revenue",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "revenue"},
		Templates: TemplatesFS,
	}
}
