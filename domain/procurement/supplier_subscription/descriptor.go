package supplier_subscription

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "procurement.supplier_subscription",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "supplier_subscription"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "supplier_subscription.json", Key: "supplierSubscription"},
		LabelName: "SupplierSubscriptionLabels",
		Templates: templateFS,
		Nav: compose.NavContrib{
			Permission: "supplier_subscription:list",
			Items: []compose.NavItem{
				// supplier app — Supplier Subscriptions section
				{Key: "supplier-subscriptions-active", Route: "supplier_subscription.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "supplier_subscription:list"},
				{Key: "supplier-subscriptions-inactive", Route: "supplier_subscription.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "supplier_subscription:list"},
			},
		},
	}
}
