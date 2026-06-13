package treasuryadvancesdashboard

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "treasury.treasuryadvancesdashboard",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "treasury_advances"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "advances_dashboard.json", Key: "advancesDashboard"},
		LabelName: "AdvancesDashboardLabels",
		Templates: TemplatesFS,
	}
}
