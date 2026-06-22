package subscription_group_product_plan_staff

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.subscription_group_product_plan_staff",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "subscription_group_product_plan_staff"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "subscription_group_product_plan_staff.json", Key: "subscriptionGroupProductPlanStaff"},
		LabelName: "SubscriptionGroupProductPlanStaffLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "subscription_group_product_plan_staff:list",
			Items: []compose.NavItem{
				{Key: "sgpps-active", Route: "subscription_group_product_plan_staff.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "subscription_group_product_plan_staff:list"},
				{Key: "sgpps-inactive", Route: "subscription_group_product_plan_staff.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "subscription_group_product_plan_staff:list"},
			},
		},
	}
}
