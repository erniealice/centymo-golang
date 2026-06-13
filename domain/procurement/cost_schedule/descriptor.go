package cost_schedule

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "procurement.cost_schedule",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "cost_schedule"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "cost_schedule.json", Key: "costSchedule"},
		LabelName: "CostScheduleLabels",
		Templates: templateFS,
	}
}
