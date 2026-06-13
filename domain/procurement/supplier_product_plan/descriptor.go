package supplier_product_plan

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "procurement.supplier_product_plan",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "supplier_product_plan"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "supplier_product_plan.json", Key: "supplierProductPlan"},
		LabelName: "SupplierProductPlanLabels",
		Templates: templateFS,
	}
}
