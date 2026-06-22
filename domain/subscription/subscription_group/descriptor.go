package subscription_group

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.subscription_group",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "subscription_group"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "subscription_group.json", Key: "subscriptionGroup"},
		LabelName: "SubscriptionGroupLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "subscription_group:list",
			Items: []compose.NavItem{
				// service app — "Sections / Cohorts" section
				{Key: "subscription-groups-active", Route: "subscription_group.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "subscription_group:list"},
				{Key: "subscription-groups-inactive", Route: "subscription_group.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "subscription_group:list"},
			},
		},
	}
}
