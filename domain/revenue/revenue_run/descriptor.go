package revenuerun

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "revenue.revenue_run",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "revenue_run"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "revenue.json", Key: "revenueRun"},
		LabelName: "RevenueRunLabels",
		Templates: TemplatesFS,
	}
}
