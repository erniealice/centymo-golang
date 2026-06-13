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
		Nav: compose.NavContrib{
			Permission: "supplier_plan:list",
			Items: []compose.NavItem{
				// supplier app — Supplier Plans section
				{Key: "supplier-plans-active", Route: "supplier_plan.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "supplier_plan:list"},
				{Key: "supplier-plans-inactive", Route: "supplier_plan.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "supplier_plan:list"},
			},
		},
	}
}
