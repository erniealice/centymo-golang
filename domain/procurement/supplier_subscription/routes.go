package supplier_subscription

// routes.go — supplier_subscription route constants + Routes type (centymo W7).
//
// Extracted from the procurement-domain routes.go per the domain-first
// restructure. Pure structural move — no behaviour change; route strings are
// byte-identical.

// Default route constants for the supplier_subscription module
// (20260506-supplier-subscriptions).
const (
	ListURL             = "/supplier-subscriptions/list/{status}"
	TableURL            = "/action/supplier-subscription/table/{status}"
	DetailURL           = "/supplier-subscriptions/detail/{id}"
	AddURL              = "/action/supplier-subscription/add"
	EditURL             = "/action/supplier-subscription/edit/{id}"
	DeleteURL           = "/action/supplier-subscription/delete"
	BulkDeleteURL       = "/action/supplier-subscription/bulk-delete"
	SetStatusURL        = "/action/supplier-subscription/set-status"
	BulkSetStatusURL    = "/action/supplier-subscription/bulk-set-status"
	TabActionURL        = "/action/supplier-subscription/detail/{id}/tab/{tab}"
	SearchCostPlanURL   = "/action/supplier-subscription/search/cost-plans"
	SearchSupplierURL   = "/action/supplier-subscription/search/suppliers"
	RecognizeExpenseURL = "/action/supplier-subscription/recognize-expense/{id}"
)

// Routes holds all route paths for supplier_subscription views.
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

	// Search autocomplete endpoints for the add/edit drawer
	SearchCostPlanURL string `json:"search_cost_plan_url"`
	SearchSupplierURL string `json:"search_supplier_url"`

	// Recognition CTA — POST; opens the recognize-expense drawer on the detail page.
	RecognizeExpenseURL string `json:"recognize_expense_url"`

	// ExpenseRecognitionRunURL — GET; opens the per-SupplierSubscription Expense
	// Recognition Run drawer (Surface C). Resolved by resolveRecognitionsPrimaryAction
	// for CostPlan.billing_kind RECURRING / CONTRACT-with-cycle.
	// Plan A 20260517-expense-run Phase 4 / Surface C.
	ExpenseRecognitionRunURL string `json:"expense_recognition_run_url"`
}

// DefaultRoutes returns a Routes using route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-subscriptions",

		ListURL:             ListURL,
		TableURL:            TableURL,
		DetailURL:           DetailURL,
		AddURL:              AddURL,
		EditURL:             EditURL,
		DeleteURL:           DeleteURL,
		BulkDeleteURL:       BulkDeleteURL,
		SetStatusURL:        SetStatusURL,
		BulkSetStatusURL:    BulkSetStatusURL,
		TabActionURL:        TabActionURL,
		SearchCostPlanURL:   SearchCostPlanURL,
		SearchSupplierURL:   SearchSupplierURL,
		RecognizeExpenseURL: RecognizeExpenseURL,
		// centymo W7: ExpenseRecognitionRunPerSubscriptionDrawerURL is an
		// expenditure-domain URL const. Inlining the literal here avoids
		// a procurement→expenditure domain import while preserving the exact
		// same value — same approach W5/W6 used for analogous route strings.
		ExpenseRecognitionRunURL: "/action/supplier-subscription/expense-recognition-run/{id}",
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_subscription.list":                    r.ListURL,
		"supplier_subscription.table":                   r.TableURL,
		"supplier_subscription.detail":                  r.DetailURL,
		"supplier_subscription.add":                     r.AddURL,
		"supplier_subscription.edit":                    r.EditURL,
		"supplier_subscription.delete":                  r.DeleteURL,
		"supplier_subscription.bulk_delete":             r.BulkDeleteURL,
		"supplier_subscription.set_status":              r.SetStatusURL,
		"supplier_subscription.bulk_set_status":         r.BulkSetStatusURL,
		"supplier_subscription.tab_action":              r.TabActionURL,
		"supplier_subscription.search_cost_plan":        r.SearchCostPlanURL,
		"supplier_subscription.search_supplier":         r.SearchSupplierURL,
		"supplier_subscription.recognize_expense":       r.RecognizeExpenseURL,
		"supplier_subscription.expense_recognition_run": r.ExpenseRecognitionRunURL,
	}
}
