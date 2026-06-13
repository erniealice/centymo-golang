package supplier_plan

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "procurement.supplier_plan",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "supplier_plan"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "supplier_plan.json", Key: "supplierPlan"},
		LabelName: "SupplierPlanLabels",
		Templates: templateFS,
	}
}
