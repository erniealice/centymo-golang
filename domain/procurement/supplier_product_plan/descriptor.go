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
		Nav: compose.NavContrib{
			Permission: "supplier_product_plan:list",
			Items: []compose.NavItem{
				// supplier app — Supplier Product Plans section
				{Key: "supplier-product-plans-active", Route: "supplier_product_plan.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "supplier_product_plan:list"},
				{Key: "supplier-product-plans-inactive", Route: "supplier_product_plan.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "supplier_product_plan:list"},
			},
		},
	}
}
