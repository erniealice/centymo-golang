package product_plan_staff

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "product.product_plan_staff",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "product_plan_staff"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "product_plan_staff.json", Key: "productPlanStaff"},
		LabelName: "ProductPlanStaffLabels",
		Templates: TemplateFS,
		Nav: compose.NavContrib{
			Permission: "product_plan_staff:list",
			Items: []compose.NavItem{
				{Key: "product-plan-staffs-active", Route: "product_plan_staff.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "product_plan_staff:list"},
				{Key: "product-plan-staffs-inactive", Route: "product_plan_staff.list", Params: map[string]string{"status": "inactive"},
					Label: "Inactive", Icon: "icon-circle", Permission: "product_plan_staff:list"},
			},
		},
	}
}
