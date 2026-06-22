package subscription_group_product_plan_staff

// SubscriptionGroupProductPlanStaff-domain route constants. The class-edge
// records which staff member delivers which subject (product_plan) in which
// section (subscription_group) and in what role.
const (
	DashboardURL     = "/subscription-group-product-plan-staffs/dashboard"
	ListURL          = "/subscription-group-product-plan-staffs/list/{status}"
	TableURL         = "/action/subscription-group-product-plan-staff/table/{status}"
	DetailURL        = "/subscription-group-product-plan-staffs/detail/{id}"
	AddURL           = "/action/subscription-group-product-plan-staff/add"
	EditURL          = "/action/subscription-group-product-plan-staff/edit/{id}"
	DeleteURL        = "/action/subscription-group-product-plan-staff/delete"
	BulkDeleteURL    = "/action/subscription-group-product-plan-staff/bulk-delete"
	SetStatusURL     = "/action/subscription-group-product-plan-staff/set-status"
	BulkSetStatusURL = "/action/subscription-group-product-plan-staff/bulk-set-status"
	TabActionURL     = "/action/subscription-group-product-plan-staff/{id}/tab/{tab}"
)

// Routes holds all route paths for subscription_group_product_plan_staff views
// and actions.
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
		ActiveNav:        "service",
		ActiveSubNav:     "subscription-group-product-plan-staffs",
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
// subscription_group_product_plan_staff routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"subscription_group_product_plan_staff.dashboard":       r.DashboardURL,
		"subscription_group_product_plan_staff.list":            r.ListURL,
		"subscription_group_product_plan_staff.table":           r.TableURL,
		"subscription_group_product_plan_staff.detail":          r.DetailURL,
		"subscription_group_product_plan_staff.add":             r.AddURL,
		"subscription_group_product_plan_staff.edit":            r.EditURL,
		"subscription_group_product_plan_staff.delete":          r.DeleteURL,
		"subscription_group_product_plan_staff.bulk_delete":     r.BulkDeleteURL,
		"subscription_group_product_plan_staff.set_status":      r.SetStatusURL,
		"subscription_group_product_plan_staff.bulk_set_status": r.BulkSetStatusURL,
		"subscription_group_product_plan_staff.tab_action":      r.TabActionURL,
	}
}
