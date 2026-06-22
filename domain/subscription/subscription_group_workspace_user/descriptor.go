package subscription_group_workspace_user

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.subscription_group_workspace_user",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "subscription_group_workspace_user"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "subscription_group_workspace_user.json", Key: "subscriptionGroupWorkspaceUser"},
		LabelName: "SubscriptionGroupWorkspaceUserLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "subscription_group_workspace_user:list",
			Items: []compose.NavItem{
				{Key: "subscription-group-workspace-users-active", Route: "subscription_group_workspace_user.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "subscription_group_workspace_user:list"},
				{Key: "subscription-group-workspace-users-inactive", Route: "subscription_group_workspace_user.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "subscription_group_workspace_user:list"},
			},
		},
	}
}
