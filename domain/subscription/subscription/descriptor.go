package subscription

import "github.com/erniealice/espyna-golang/consumer/compose"

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
		Nav: compose.NavContrib{
			Permission: "subscription:list",
			Items: []compose.NavItem{
				// job app (Engagements / Subscriptions section)
				{Key: "subscriptions-active", Route: "subscription.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "subscription:list"},
				{Key: "subscriptions-inactive", Route: "subscription.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "subscription:list"},
			},
		},
	}
}
