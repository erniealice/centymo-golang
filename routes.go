package centymo

// routes.go — RESIDUAL after centymo W6.
//
// Expenditure-domain route constants moved to domain/expenditure/routes.go (W6).
// What remains is the procurement-domain (W7) route constants: the Procurement
// Operations composition app + SupplierSubscription / CostSchedule / SupplierPlan /
// CostPlan / SupplierProductPlan(+cost) consts.
// Consumer apps can use these or define their own.
const (
	// ---------------------------------------------------------------------------
	// P3b — Procurement Operations app route constants
	// (composition surface; no proto entity)
	// ---------------------------------------------------------------------------

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
