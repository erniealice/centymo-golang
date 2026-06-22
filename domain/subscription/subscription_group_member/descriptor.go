package subscription_group_member

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.subscription_group_member",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "subscription_group_member"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "subscription_group_member.json", Key: "subscriptionGroupMember"},
		LabelName: "SubscriptionGroupMemberLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "subscription_group_member:list",
			Items: []compose.NavItem{
				{Key: "subscription-group-members-active", Route: "subscription_group_member.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "subscription_group_member:list"},
				{Key: "subscription-group-members-inactive", Route: "subscription_group_member.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "subscription_group_member:list"},
			},
		},
	}
}
