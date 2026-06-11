package cost_plan

// routes.go — cost_plan route constants + Routes type (centymo W7).
//
// Extracted from the procurement-domain routes.go per the domain-first
// restructure. Pure structural move — no behaviour change; route strings are
// byte-identical.

// Default route constants for the cost_plan module.
const (
	ListURL          = "/cost-plans/list/{status}"
	TableURL         = "/action/cost-plan/table/{status}"
	DetailURL        = "/cost-plans/detail/{id}"
	AddURL           = "/action/cost-plan/add"
	EditURL          = "/action/cost-plan/edit/{id}"
	DeleteURL        = "/action/cost-plan/delete"
	BulkDeleteURL    = "/action/cost-plan/bulk-delete"
	SetStatusURL     = "/action/cost-plan/set-status"
	BulkSetStatusURL = "/action/cost-plan/bulk-set-status"
	TabActionURL     = "/action/cost-plan/detail/{id}/tab/{tab}"

	// SupplierProductCostPlan CRUD routes (inline within CostPlan detail)
	ProductCostAddURL    = "/action/cost-plan/{id}/product-costs/add"
	ProductCostEditURL   = "/action/cost-plan/{id}/product-costs/edit/{pcid}"
	ProductCostDeleteURL = "/action/cost-plan/{id}/product-costs/delete"
)

// Routes holds all route paths for cost_plan views.
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

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

	// SupplierProductCostPlan inline CRUD within cost_plan detail
	ProductCostAddURL    string `json:"product_cost_add_url"`
	ProductCostEditURL   string `json:"product_cost_edit_url"`
	ProductCostDeleteURL string `json:"product_cost_delete_url"`

	// Autocomplete search URLs for add/edit form selects.
	SearchSupplierPlanURL        string `json:"search_supplier_plan_url"`
	SearchCostScheduleURL        string `json:"search_cost_schedule_url"`
	SearchSupplierProductPlanURL string `json:"search_supplier_product_plan_url"`
}

// DefaultRoutes returns a Routes using route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "supplier",
		ActiveSubNav: "cost-plans",

		ListURL:              ListURL,
		TableURL:             TableURL,
		DetailURL:            DetailURL,
		AddURL:               AddURL,
		EditURL:              EditURL,
		DeleteURL:            DeleteURL,
		BulkDeleteURL:        BulkDeleteURL,
		SetStatusURL:         SetStatusURL,
		BulkSetStatusURL:     BulkSetStatusURL,
		TabActionURL:         TabActionURL,
		ProductCostAddURL:    ProductCostAddURL,
		ProductCostEditURL:   ProductCostEditURL,
		ProductCostDeleteURL: ProductCostDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"cost_plan.list":                r.ListURL,
		"cost_plan.table":               r.TableURL,
		"cost_plan.detail":              r.DetailURL,
		"cost_plan.add":                 r.AddURL,
		"cost_plan.edit":                r.EditURL,
		"cost_plan.delete":              r.DeleteURL,
		"cost_plan.bulk_delete":         r.BulkDeleteURL,
		"cost_plan.set_status":          r.SetStatusURL,
		"cost_plan.bulk_set_status":     r.BulkSetStatusURL,
		"cost_plan.tab_action":          r.TabActionURL,
		"cost_plan.product_cost.add":    r.ProductCostAddURL,
		"cost_plan.product_cost.edit":   r.ProductCostEditURL,
		"cost_plan.product_cost.delete": r.ProductCostDeleteURL,
	}
}
