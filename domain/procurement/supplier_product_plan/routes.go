package supplier_product_plan

// routes.go — supplier_product_plan route constants + Routes type (centymo W7).
//
// Extracted from the procurement-domain routes.go per the domain-first
// restructure. Pure structural move — no behaviour change; route strings are
// byte-identical.

// Default route constants for the supplier_product_plan module.
const (
	ListURL          = "/supplier-product-plans/list/{status}"
	TableURL         = "/action/supplier-product-plan/table/{status}"
	DetailURL        = "/supplier-product-plans/detail/{id}"
	AddURL           = "/action/supplier-product-plan/add"
	EditURL          = "/action/supplier-product-plan/edit/{id}"
	DeleteURL        = "/action/supplier-product-plan/delete"
	BulkDeleteURL    = "/action/supplier-product-plan/bulk-delete"
	SetStatusURL     = "/action/supplier-product-plan/set-status"
	BulkSetStatusURL = "/action/supplier-product-plan/bulk-set-status"
	TabActionURL     = "/action/supplier-product-plan/detail/{id}/tab/{tab}"
)

// Routes holds all route paths for supplier_product_plan views.
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

	// Autocomplete search URLs for add/edit form selects.
	SearchSupplierPlanURL string `json:"search_supplier_plan_url"`
	SearchProductURL      string `json:"search_product_url"`
}

// DefaultRoutes returns a Routes using route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-product-plans",

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

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_product_plan.list":            r.ListURL,
		"supplier_product_plan.table":           r.TableURL,
		"supplier_product_plan.detail":          r.DetailURL,
		"supplier_product_plan.add":             r.AddURL,
		"supplier_product_plan.edit":            r.EditURL,
		"supplier_product_plan.delete":          r.DeleteURL,
		"supplier_product_plan.bulk_delete":     r.BulkDeleteURL,
		"supplier_product_plan.set_status":      r.SetStatusURL,
		"supplier_product_plan.bulk_set_status": r.BulkSetStatusURL,
		"supplier_product_plan.tab_action":      r.TabActionURL,
	}
}
