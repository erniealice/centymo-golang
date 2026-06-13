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
	}
}
