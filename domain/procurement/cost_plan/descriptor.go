package cost_plan

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "procurement.cost_plan",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "cost_plan"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "cost_plan.json", Key: "costPlan"},
		LabelName: "CostPlanLabels",
		Templates: templateFS,
	}
}
