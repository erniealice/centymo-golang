package price_plan

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.price_plan",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "price_plan"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "price_plan.json", Key: "price_plan"},
		LabelName: "PricePlanLabels",
		Templates: TemplateFS,
	}
}
