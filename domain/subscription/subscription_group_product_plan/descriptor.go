package subscription_group_product_plan

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.subscription_group_product_plan",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "subscription_group_product_plan"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "subscription_group_product_plan.json", Key: "subscription_group_product_plan"},
		LabelName: "SubscriptionGroupProductPlanLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "subscription_group_product_plan:list",
			Items: []compose.NavItem{
				{Key: "subscription-group-product-plans-active", Route: "subscription_group_product_plan.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "subscription_group_product_plan:list"},
				{Key: "subscription-group-product-plans-inactive", Route: "subscription_group_product_plan.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "subscription_group_product_plan:list"},
			},
		},
	}
}
