package subscription

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.subscription",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "subscription"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "subscription.json", Key: "subscription"},
		LabelName: "SubscriptionLabels",
		Templates: TemplatesFS,
	}
}
