package procurement

// routes.go — procurement-domain route constants (centymo W7).
//
// Extracted from the root routes.go (Procurement Operations URL consts) and
// routes_config.go (ProcurementRoutes type, DefaultProcurementRoutes() constructor,
// RouteMap() method) into the procurement domain package per the domain-first
// restructure. Pure structural move — no behaviour change; route strings are
// byte-identical.

// Default route constants for the Procurement Operations composition app.
// Consumer apps can use these or define their own.
const (
	// Procurement Operations app — all GET, read-only views
	ProcurementDashboardURL        = "/procurement/dashboard"
	ProcurementRenewalCalendarURL  = "/procurement/renewals"
	ProcurementVarianceURL         = "/procurement/variance"
	ProcurementUtilizationURL      = "/procurement/utilization"
	ProcurementRecurrenceDraftsURL = "/procurement/recurrence-drafts/list/{status}"

	// ---------------------------------------------------------------------------
	// P3 — SupplierSubscription route constants (20260506-supplier-subscriptions)
	// ---------------------------------------------------------------------------

	SupplierSubscriptionListURL             = "/supplier-subscriptions/list/{status}"
	SupplierSubscriptionTableURL            = "/action/supplier-subscription/table/{status}"
	SupplierSubscriptionDetailURL           = "/supplier-subscriptions/detail/{id}"
	SupplierSubscriptionAddURL              = "/action/supplier-subscription/add"
	SupplierSubscriptionEditURL             = "/action/supplier-subscription/edit/{id}"
	SupplierSubscriptionDeleteURL           = "/action/supplier-subscription/delete"
	SupplierSubscriptionBulkDeleteURL       = "/action/supplier-subscription/bulk-delete"
	SupplierSubscriptionSetStatusURL        = "/action/supplier-subscription/set-status"
	SupplierSubscriptionBulkSetStatusURL    = "/action/supplier-subscription/bulk-set-status"
	SupplierSubscriptionTabActionURL        = "/action/supplier-subscription/detail/{id}/tab/{tab}"
	SupplierSubscriptionSearchCostPlanURL   = "/action/supplier-subscription/search/cost-plans"
	SupplierSubscriptionSearchSupplierURL   = "/action/supplier-subscription/search/suppliers"
	SupplierSubscriptionRecognizeExpenseURL = "/action/supplier-subscription/recognize-expense/{id}"

	// ---------------------------------------------------------------------------
	// P3 — CostSchedule route constants
	// ---------------------------------------------------------------------------

	CostScheduleListURL          = "/cost-schedules/list/{status}"
	CostScheduleTableURL         = "/action/cost-schedule/table/{status}"
	CostScheduleDetailURL        = "/cost-schedules/detail/{id}"
	CostScheduleAddURL           = "/action/cost-schedule/add"
	CostScheduleEditURL          = "/action/cost-schedule/edit/{id}"
	CostScheduleDeleteURL        = "/action/cost-schedule/delete"
	CostScheduleBulkDeleteURL    = "/action/cost-schedule/bulk-delete"
	CostScheduleSetStatusURL     = "/action/cost-schedule/set-status"
	CostScheduleBulkSetStatusURL = "/action/cost-schedule/bulk-set-status"
	CostScheduleTabActionURL     = "/action/cost-schedule/detail/{id}/tab/{tab}"

	// ---------------------------------------------------------------------------
	// P3 — SupplierPlan route constants
	// ---------------------------------------------------------------------------

	SupplierPlanListURL          = "/supplier-plans/list/{status}"
	SupplierPlanTableURL         = "/action/supplier-plan/table/{status}"
	SupplierPlanDetailURL        = "/supplier-plans/detail/{id}"
	SupplierPlanAddURL           = "/action/supplier-plan/add"
	SupplierPlanEditURL          = "/action/supplier-plan/edit/{id}"
	SupplierPlanDeleteURL        = "/action/supplier-plan/delete"
	SupplierPlanBulkDeleteURL    = "/action/supplier-plan/bulk-delete"
	SupplierPlanSetStatusURL     = "/action/supplier-plan/set-status"
	SupplierPlanBulkSetStatusURL = "/action/supplier-plan/bulk-set-status"
	SupplierPlanTabActionURL     = "/action/supplier-plan/detail/{id}/tab/{tab}"

	// ---------------------------------------------------------------------------
	// P3 — CostPlan route constants
	// ---------------------------------------------------------------------------

	CostPlanListURL          = "/cost-plans/list/{status}"
	CostPlanTableURL         = "/action/cost-plan/table/{status}"
	CostPlanDetailURL        = "/cost-plans/detail/{id}"
	CostPlanAddURL           = "/action/cost-plan/add"
	CostPlanEditURL          = "/action/cost-plan/edit/{id}"
	CostPlanDeleteURL        = "/action/cost-plan/delete"
	CostPlanBulkDeleteURL    = "/action/cost-plan/bulk-delete"
	CostPlanSetStatusURL     = "/action/cost-plan/set-status"
	CostPlanBulkSetStatusURL = "/action/cost-plan/bulk-set-status"
	CostPlanTabActionURL     = "/action/cost-plan/detail/{id}/tab/{tab}"

	// SupplierProductCostPlan CRUD routes (inline within CostPlan detail)
	CostPlanProductCostAddURL    = "/action/cost-plan/{id}/product-costs/add"
	CostPlanProductCostEditURL   = "/action/cost-plan/{id}/product-costs/edit/{pcid}"
	CostPlanProductCostDeleteURL = "/action/cost-plan/{id}/product-costs/delete"

	// ---------------------------------------------------------------------------
	// P3 — SupplierProductPlan route constants
	// ---------------------------------------------------------------------------

	SupplierProductPlanListURL          = "/supplier-product-plans/list/{status}"
	SupplierProductPlanTableURL         = "/action/supplier-product-plan/table/{status}"
	SupplierProductPlanDetailURL        = "/supplier-product-plans/detail/{id}"
	SupplierProductPlanAddURL           = "/action/supplier-product-plan/add"
	SupplierProductPlanEditURL          = "/action/supplier-product-plan/edit/{id}"
	SupplierProductPlanDeleteURL        = "/action/supplier-product-plan/delete"
	SupplierProductPlanBulkDeleteURL    = "/action/supplier-product-plan/bulk-delete"
	SupplierProductPlanSetStatusURL     = "/action/supplier-product-plan/set-status"
	SupplierProductPlanBulkSetStatusURL = "/action/supplier-product-plan/bulk-set-status"
	SupplierProductPlanTabActionURL     = "/action/supplier-product-plan/detail/{id}/tab/{tab}"
)

// ---------------------------------------------------------------------------
// P3 — SupplierSubscription routes (20260506-supplier-subscriptions)
// ---------------------------------------------------------------------------

// SupplierSubscriptionRoutes holds all route paths for supplier_subscription views.
type SupplierSubscriptionRoutes struct {
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

// DefaultSupplierSubscriptionRoutes returns a SupplierSubscriptionRoutes using route constants.
func DefaultSupplierSubscriptionRoutes() SupplierSubscriptionRoutes {
	return SupplierSubscriptionRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-subscriptions",

		ListURL:             SupplierSubscriptionListURL,
		TableURL:            SupplierSubscriptionTableURL,
		DetailURL:           SupplierSubscriptionDetailURL,
		AddURL:              SupplierSubscriptionAddURL,
		EditURL:             SupplierSubscriptionEditURL,
		DeleteURL:           SupplierSubscriptionDeleteURL,
		BulkDeleteURL:       SupplierSubscriptionBulkDeleteURL,
		SetStatusURL:        SupplierSubscriptionSetStatusURL,
		BulkSetStatusURL:    SupplierSubscriptionBulkSetStatusURL,
		TabActionURL:        SupplierSubscriptionTabActionURL,
		SearchCostPlanURL:   SupplierSubscriptionSearchCostPlanURL,
		SearchSupplierURL:   SupplierSubscriptionSearchSupplierURL,
		RecognizeExpenseURL: SupplierSubscriptionRecognizeExpenseURL,
		// centymo W7: ExpenseRecognitionRunPerSubscriptionDrawerURL is an
		// expenditure-domain URL const. Inlining the literal here avoids
		// a procurement→expenditure domain import while preserving the exact
		// same value — same approach W5/W6 used for analogous route strings.
		ExpenseRecognitionRunURL: "/action/supplier-subscription/expense-recognition-run/{id}",
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierSubscriptionRoutes) RouteMap() map[string]string {
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

// ---------------------------------------------------------------------------
// P3 — CostSchedule routes
// ---------------------------------------------------------------------------

// CostScheduleRoutes holds all route paths for cost_schedule views.
type CostScheduleRoutes struct {
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
}

// DefaultCostScheduleRoutes returns a CostScheduleRoutes using route constants.
func DefaultCostScheduleRoutes() CostScheduleRoutes {
	return CostScheduleRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "cost-schedules",

		ListURL:          CostScheduleListURL,
		TableURL:         CostScheduleTableURL,
		DetailURL:        CostScheduleDetailURL,
		AddURL:           CostScheduleAddURL,
		EditURL:          CostScheduleEditURL,
		DeleteURL:        CostScheduleDeleteURL,
		BulkDeleteURL:    CostScheduleBulkDeleteURL,
		SetStatusURL:     CostScheduleSetStatusURL,
		BulkSetStatusURL: CostScheduleBulkSetStatusURL,
		TabActionURL:     CostScheduleTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r CostScheduleRoutes) RouteMap() map[string]string {
	return map[string]string{
		"cost_schedule.list":            r.ListURL,
		"cost_schedule.table":           r.TableURL,
		"cost_schedule.detail":          r.DetailURL,
		"cost_schedule.add":             r.AddURL,
		"cost_schedule.edit":            r.EditURL,
		"cost_schedule.delete":          r.DeleteURL,
		"cost_schedule.bulk_delete":     r.BulkDeleteURL,
		"cost_schedule.set_status":      r.SetStatusURL,
		"cost_schedule.bulk_set_status": r.BulkSetStatusURL,
		"cost_schedule.tab_action":      r.TabActionURL,
	}
}

// ---------------------------------------------------------------------------
// P3 — SupplierPlan routes
// ---------------------------------------------------------------------------

// SupplierPlanRoutes holds all route paths for supplier_plan views.
type SupplierPlanRoutes struct {
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

	// Autocomplete search URL for the supplier select in add/edit forms.
	SearchSupplierURL string `json:"search_supplier_url"`
}

// DefaultSupplierPlanRoutes returns a SupplierPlanRoutes using route constants.
func DefaultSupplierPlanRoutes() SupplierPlanRoutes {
	return SupplierPlanRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-plans",

		ListURL:          SupplierPlanListURL,
		TableURL:         SupplierPlanTableURL,
		DetailURL:        SupplierPlanDetailURL,
		AddURL:           SupplierPlanAddURL,
		EditURL:          SupplierPlanEditURL,
		DeleteURL:        SupplierPlanDeleteURL,
		BulkDeleteURL:    SupplierPlanBulkDeleteURL,
		SetStatusURL:     SupplierPlanSetStatusURL,
		BulkSetStatusURL: SupplierPlanBulkSetStatusURL,
		TabActionURL:     SupplierPlanTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierPlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"supplier_plan.list":            r.ListURL,
		"supplier_plan.table":           r.TableURL,
		"supplier_plan.detail":          r.DetailURL,
		"supplier_plan.add":             r.AddURL,
		"supplier_plan.edit":            r.EditURL,
		"supplier_plan.delete":          r.DeleteURL,
		"supplier_plan.bulk_delete":     r.BulkDeleteURL,
		"supplier_plan.set_status":      r.SetStatusURL,
		"supplier_plan.bulk_set_status": r.BulkSetStatusURL,
		"supplier_plan.tab_action":      r.TabActionURL,
	}
}

// ---------------------------------------------------------------------------
// P3 — CostPlan routes
// ---------------------------------------------------------------------------

// CostPlanRoutes holds all route paths for cost_plan views.
type CostPlanRoutes struct {
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

// DefaultCostPlanRoutes returns a CostPlanRoutes using route constants.
func DefaultCostPlanRoutes() CostPlanRoutes {
	return CostPlanRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "cost-plans",

		ListURL:              CostPlanListURL,
		TableURL:             CostPlanTableURL,
		DetailURL:            CostPlanDetailURL,
		AddURL:               CostPlanAddURL,
		EditURL:              CostPlanEditURL,
		DeleteURL:            CostPlanDeleteURL,
		BulkDeleteURL:        CostPlanBulkDeleteURL,
		SetStatusURL:         CostPlanSetStatusURL,
		BulkSetStatusURL:     CostPlanBulkSetStatusURL,
		TabActionURL:         CostPlanTabActionURL,
		ProductCostAddURL:    CostPlanProductCostAddURL,
		ProductCostEditURL:   CostPlanProductCostEditURL,
		ProductCostDeleteURL: CostPlanProductCostDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r CostPlanRoutes) RouteMap() map[string]string {
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

// ---------------------------------------------------------------------------
// P3 — SupplierProductPlan routes
// ---------------------------------------------------------------------------

// SupplierProductPlanRoutes holds all route paths for supplier_product_plan views.
type SupplierProductPlanRoutes struct {
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

// DefaultSupplierProductPlanRoutes returns a SupplierProductPlanRoutes using route constants.
func DefaultSupplierProductPlanRoutes() SupplierProductPlanRoutes {
	return SupplierProductPlanRoutes{
		ActiveNav:    "supplier",
		ActiveSubNav: "supplier-product-plans",

		ListURL:          SupplierProductPlanListURL,
		TableURL:         SupplierProductPlanTableURL,
		DetailURL:        SupplierProductPlanDetailURL,
		AddURL:           SupplierProductPlanAddURL,
		EditURL:          SupplierProductPlanEditURL,
		DeleteURL:        SupplierProductPlanDeleteURL,
		BulkDeleteURL:    SupplierProductPlanBulkDeleteURL,
		SetStatusURL:     SupplierProductPlanSetStatusURL,
		BulkSetStatusURL: SupplierProductPlanBulkSetStatusURL,
		TabActionURL:     SupplierProductPlanTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r SupplierProductPlanRoutes) RouteMap() map[string]string {
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
