package product_plan_staff

// ProductPlanStaff domain route constants. A product_plan_staff row is the
// eligibility edge — the qualified-staff pool for a subject-in-program
// (product_plan). It constrains the deliverer picker; it is not an access
// grant.
const (
	DashboardURL     = "/product-plan-staffs/dashboard"
	ListURL          = "/product-plan-staffs/list/{status}"
	TableURL         = "/action/product-plan-staff/table/{status}"
	DetailURL        = "/product-plan-staffs/detail/{id}"
	AddURL           = "/action/product-plan-staff/add"
	EditURL          = "/action/product-plan-staff/edit/{id}"
	DeleteURL        = "/action/product-plan-staff/delete"
	BulkDeleteURL    = "/action/product-plan-staff/bulk-delete"
	SetStatusURL     = "/action/product-plan-staff/set-status"
	BulkSetStatusURL = "/action/product-plan-staff/bulk-set-status"
	TabActionURL     = "/action/product-plan-staff/{id}/tab/{tab}"
)

// Routes holds all route paths for product_plan_staff views and actions.
type Routes struct {
	ActiveNav        string `json:"active_nav"`
	ActiveSubNav     string `json:"active_sub_nav"`
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route
// constants defined above.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:        "product",
		ActiveSubNav:     "product-plan-staffs",
		DashboardURL:     DashboardURL,
		ListURL:          ListURL,
		TableURL:         TableURL,
		DetailURL:        DetailURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,
		TabActionURL:     TabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// product_plan_staff routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"product_plan_staff.dashboard":       r.DashboardURL,
		"product_plan_staff.list":            r.ListURL,
		"product_plan_staff.table":           r.TableURL,
		"product_plan_staff.detail":          r.DetailURL,
		"product_plan_staff.add":             r.AddURL,
		"product_plan_staff.edit":            r.EditURL,
		"product_plan_staff.delete":          r.DeleteURL,
		"product_plan_staff.bulk_delete":     r.BulkDeleteURL,
		"product_plan_staff.set_status":      r.SetStatusURL,
		"product_plan_staff.bulk_set_status": r.BulkSetStatusURL,
		"product_plan_staff.tab_action":      r.TabActionURL,
	}
}
