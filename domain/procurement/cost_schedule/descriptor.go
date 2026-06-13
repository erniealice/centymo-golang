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
		Nav: compose.NavContrib{
			Permission: "cost_schedule:list",
			Items: []compose.NavItem{
				// supplier app — Cost Schedules section
				{Key: "cost-schedules-active", Route: "cost_schedule.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "cost_schedule:list"},
				{Key: "cost-schedules-inactive", Route: "cost_schedule.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "cost_schedule:list"},
			},
		},
	}
}
