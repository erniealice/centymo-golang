package price_schedule_workspace_user

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.price_schedule_workspace_user",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "price_schedule_workspace_user"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "price_schedule_workspace_user.json", Key: "priceScheduleWorkspaceUser"},
		LabelName: "PriceScheduleWorkspaceUserLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "price_schedule_workspace_user:list",
			Items: []compose.NavItem{
				{Key: "price-schedule-workspace-users-active", Route: "price_schedule_workspace_user.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "price_schedule_workspace_user:list"},
				{Key: "price-schedule-workspace-users-inactive", Route: "price_schedule_workspace_user.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "price_schedule_workspace_user:list"},
			},
		},
	}
}
