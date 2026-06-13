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
	}
}
