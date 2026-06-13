package expenditure

import "github.com/erniealice/pyeza-golang/compose"

// Describe returns the composition-v2 descriptor for the expenditure entity.
// Labels are not yet exposed via DefaultLabels() — the LabelJSON binding is
// left empty until a DefaultLabels factory is added.
func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "expenditure.expenditure",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "expenditure"},
		Templates: TemplatesFS,
	}
}
