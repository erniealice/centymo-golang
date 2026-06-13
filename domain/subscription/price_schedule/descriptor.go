package price_schedule

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.price_schedule",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "price_schedule"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "price_schedule.json", Key: "priceSchedule"},
		LabelName: "PriceScheduleLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "price_schedule:list",
			Items: []compose.NavItem{
				// service app — "Price Lists" / "Rate Cards" section
				{Key: "price-schedules-active", Route: "price_schedule.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "price_schedule:list"},
				{Key: "price-schedules-inactive", Route: "price_schedule.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "price_schedule:list"},
			},
		},
	}
}
