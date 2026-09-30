package cost_source_component

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the compose descriptor for the recoverable cost lines. The
// unit is opt-in (block.WithKnownCostRecovery): it carries no sidebar entry —
// the lines are managed from the expenditure detail tab.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.cost_source_component",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "cost_source_component"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "cost_source_component.json", Key: "cost_source_component"},
		LabelName: "CostSourceComponentLabels",
		Templates: TemplatesFS,
	}
}
