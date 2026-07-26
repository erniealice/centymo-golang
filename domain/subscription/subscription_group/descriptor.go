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
		LabelJSON: compose.JSONBinding{File: "subscription_group.json", Key: "subscription_group"},
		LabelName: "SubscriptionGroupLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "subscription_group:list",
			Items: []compose.NavItem{
				// service app — "Sections / Cohorts" section. Items split by the
				// lifecycle status category (current/completed/draft), not the
				// active visibility bool; labels reuse the status enum labels.
				// The legacy active/inactive list URLs stay routable directly.
				{Key: "subscription-groups-current", Route: "subscription_group.list", Params: map[string]string{"status": "current"},
					Label: l.Form.StatusCurrent, Icon: "icon-check-circle", Permission: "subscription_group:list"},
				{Key: "subscription-groups-completed", Route: "subscription_group.list", Params: map[string]string{"status": "completed"},
					Label: l.Form.StatusCompleted, Icon: "icon-archive", Permission: "subscription_group:list"},
				{Key: "subscription-groups-draft", Route: "subscription_group.list", Params: map[string]string{"status": "draft"},
					Label: l.Form.StatusDraft, Icon: "icon-edit", Permission: "subscription_group:list"},
			},
		},
	}
}
