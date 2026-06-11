package centymo

import (
	"github.com/erniealice/centymo-golang/domain/subscription"
	"github.com/erniealice/centymo-golang/domain/treasury"
)

// ── centymo W4 subscription-domain compatibility shim ────────────────────────
// The Subscription/PriceSchedule route types + their Default* constructors moved
// to domain/subscription/ (centymo W4). entydad-golang/block/route_loading.go is
// an EXTERNAL consumer (outside this wave's edit scope) that still references
// centymo.SubscriptionRoutes / centymo.PriceScheduleRoutes and their Default*
// constructors. These thin aliases + forwarders keep that consumer compiling
// with ZERO behaviour change (pure type-identity aliases). Remove once entydad
// is re-pointed to domain/subscription directly (W9 / entydad-coordinated).
type SubscriptionRoutes = subscription.SubscriptionRoutes
type PriceScheduleRoutes = subscription.PriceScheduleRoutes

func DefaultSubscriptionRoutes() SubscriptionRoutes { return subscription.DefaultSubscriptionRoutes() }
func DefaultPriceScheduleRoutes() PriceScheduleRoutes {
	return subscription.DefaultPriceScheduleRoutes()
}

// ── centymo W5 treasury-domain compatibility shim ────────────────────────────
// Treasury types (Collection/Disbursement labels+routes, TreasuryAdvancesRoutes,
// the AdvanceRecognizeMilestone view I/O) moved to domain/treasury/ (centymo W5).
// The not-yet-migrated W6 view packages still reference a subset of them via the
// centymo root:
//   - views/expenditure/*            -> DisbursementRoutes / DisbursementLabels /
//     DisbursementFormLabels (the expense "pay"
//     flow creates a pre-linked disbursement)
//   - views/supplier_billing_event/* -> TreasuryAdvancesRoutes (+ its Default*)
//     and AdvanceRecognizeMilestoneInput/Output
//   - domain/subscription/views/...  -> AdvanceRecognizeMilestoneInput/Output
//     (already-migrated W4 billing-event action)
//
// These thin aliases + forwarders keep those consumers compiling with ZERO
// behaviour change. Removed as each consuming domain migrates (W6 / W9).
type DisbursementRoutes = treasury.DisbursementRoutes
type DisbursementLabels = treasury.DisbursementLabels
type DisbursementFormLabels = treasury.DisbursementFormLabels
type TreasuryAdvancesRoutes = treasury.TreasuryAdvancesRoutes
type AdvanceRecognizeMilestoneInput = treasury.AdvanceRecognizeMilestoneInput
type AdvanceRecognizeMilestoneOutput = treasury.AdvanceRecognizeMilestoneOutput

func DefaultTreasuryAdvancesRoutes() TreasuryAdvancesRoutes {
	return treasury.DefaultTreasuryAdvancesRoutes()
}

// Three-level routing system for centymo views:
//
// Level 1: Generic defaults from Go consts (this file).
//   DefaultXxxRoutes() constructors return structs populated from the route
//   constants defined in routes.go. These are sensible defaults that work
//   out of the box for any app.
//
// Level 2: Industry-specific overrides via JSON (loaded by consumer apps).
//   Consumer apps can load a JSON config that partially overrides the
//   default routes. Struct fields carry json tags for unmarshalling.
//
// Level 3: App-specific overrides via Go field assignment (optional).
//   After loading defaults and/or JSON, consumer apps can programmatically
//   set individual fields to further customize routing.
//
// Each route struct also exposes a RouteMap() method that returns a
// map[string]string keyed by dot-notation identifiers (e.g. "product.list"),
// useful for template rendering, URL resolution, and debugging.

// ---------------------------------------------------------------------------
// P3b — Procurement Operations app routes
// (composition surface; no proto entity — mirrors the schedule/cyta pattern)
// ---------------------------------------------------------------------------

// ProcurementRoutes holds the URL constants for the Procurement Operations app.
// These are defined in the centymo package so service-admin composition (P3c)
// can wire them into SidebarRoutes.Operations.Procurement.
type ProcurementRoutes struct {
	// Dashboard
	DashboardURL string `json:"dashboard_url"`

	// Contract operations (views over SupplierContract)
	RenewalCalendarURL string `json:"renewal_calendar_url"`
	VarianceURL        string `json:"variance_url"`
	UtilizationURL     string `json:"utilization_url"`

	// Recurrence drafts queue (lights up when P5 ships the recurrence engine)
	RecurrenceDraftsURL string `json:"recurrence_drafts_url"`
}

// DefaultProcurementRoutes returns a ProcurementRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultProcurementRoutes() ProcurementRoutes {
	return ProcurementRoutes{
		DashboardURL:        ProcurementDashboardURL,
		RenewalCalendarURL:  ProcurementRenewalCalendarURL,
		VarianceURL:         ProcurementVarianceURL,
		UtilizationURL:      ProcurementUtilizationURL,
		RecurrenceDraftsURL: ProcurementRecurrenceDraftsURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// procurement operations app routes.
func (r ProcurementRoutes) RouteMap() map[string]string {
	return map[string]string{
		"procurement.dashboard":         r.DashboardURL,
		"procurement.renewals":          r.RenewalCalendarURL,
		"procurement.variance":          r.VarianceURL,
		"procurement.utilization":       r.UtilizationURL,
		"procurement.recurrence_drafts": r.RecurrenceDraftsURL,
	}
}

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
		// centymo W6: ExpenseRecognitionRunPerSubscriptionDrawerURL is an
		// expenditure-domain URL const that moved to domain/expenditure/routes.go.
		// SupplierSubscription is procurement-domain (W7, still at root). Inlining
		// the literal here avoids a root->domain/expenditure import (root must not
		// import domain) while preserving the exact same value — same approach W5
		// used for the SupplierBillingEvent route strings.
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

// MapTableLabels is a shared helper used across all centymo view modules to
// produce a types.TableLabels from pyeza CommonLabels. Defined here to avoid
// duplication; all block module wirings call this.
func mapTableLabelsFromStrings(search, searchPlaceholder, sortAsc, sortDesc, noResults, loading string) struct{} {
	// Placeholder — actual implementation lives in the block package; this
	// comment documents the cross-module convention.
	return struct{}{}
}
