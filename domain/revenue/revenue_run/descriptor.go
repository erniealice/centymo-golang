package revenuerun

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "revenue.revenue_run",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "revenue_run"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "revenue.json", Key: "revenue_run"},
		LabelName: "RevenueRunLabels",
		Templates: TemplatesFS,
	}
}
