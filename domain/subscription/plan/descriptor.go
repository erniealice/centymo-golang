package plan

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.plan",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "plan"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "plan.json", Key: "plan"},
		LabelName: "PlanLabels",
		Templates: TemplatesFS,
	}
}
