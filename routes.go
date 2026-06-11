package centymo

// Default route constants for centymo views.
// Consumer apps can use these or define their own.
const (
	// SupplierBillingEvent (buying-side MILESTONE anchor).
	SupplierBillingEventListURL      = "/supplier-billing-events/list/{status}"
	SupplierBillingEventDetailURL    = "/supplier-billing-events/detail/{id}"
	SupplierBillingEventRecognizeURL = "/action/supplier-billing-event/recognize/{id}"

	// Expense Recognition Run (buying-side) routes — Plan A 20260517-expense-run.
	ExpenseRecognitionRunQueueURL                   = "/expense-recognition-run/queue"
	ExpenseRecognitionRunQueueTableURL              = "/action/expense-recognition-run/queue/table"
	ExpenseRecognitionRunListURL                    = "/expense-recognition-run/list/{status}"
	ExpenseRecognitionRunListTableURL               = "/action/expense-recognition-run/table/{status}"
	ExpenseRecognitionRunDetailURL                  = "/expense-recognition-run/detail/{id}"
	ExpenseRecognitionRunDetailTabActionURL         = "/action/expense-recognition-run/detail/{id}/tab/{tab}"
	ExpenseRecognitionRunNewURL                     = "/expense-recognition-run/new"
	ExpenseRecognitionRunGenerateURL                = "/action/expense-recognition-run/generate"
	ExpenseRecognitionRunSubmitBatchURL             = "/action/expense-recognition-run/submit-batch"
	ExpenseRecognitionRunPerSupplierDrawerURL       = "/action/supplier/expense-recognition-run/{id}"
	ExpenseRecognitionRunPerSubscriptionDrawerURL   = "/action/supplier-subscription/expense-recognition-run/{id}"

	// Expenditure (purchase + expense) routes
	ExpenditurePurchaseListURL      = "/purchases/list/{status}"
	ExpenditurePurchaseDashboardURL = "/purchases/dashboard"
	ExpenditureExpenseListURL       = "/expenses/list/{status}"
	ExpenditureExpenseDashboardURL  = "/expenses/dashboard"

	// Expenditure expense CRUD action routes
	ExpenditureExpenseAddURL       = "/action/expense/add"
	ExpenditureExpenseEditURL      = "/action/expense/edit/{id}"
	ExpenditureExpenseDeleteURL    = "/action/expense/delete"
	ExpenditureExpenseSetStatusURL = "/action/expense/set-status"
	ExpenditureExpenseDetailURL    = "/expenses/detail/{id}"
	ExpenditureExpenseTableURL     = "/action/expense/table/{status}"
	ExpenditureExpenseTabActionURL = "/action/expense/detail/{id}/tab/{tab}"
	ExpenditureAttachmentUploadURL = "/action/expense/detail/{id}/attachments/upload"
	ExpenditureAttachmentDeleteURL = "/action/expense/detail/{id}/attachments/delete"

	// Expenditure expense line item action routes
	ExpenditureExpenseLineItemAddURL    = "/action/expense/detail/{id}/items/add"
	ExpenditureExpenseLineItemEditURL   = "/action/expense/detail/{id}/items/edit/{itemId}"
	ExpenditureExpenseLineItemRemoveURL = "/action/expense/detail/{id}/items/remove"
	ExpenditureExpenseLineItemTableURL  = "/action/expense/detail/{id}/items/table"

	// Expenditure pay action route (creates a pre-linked disbursement)
	ExpenditureExpensePayURL = "/action/expense/detail/{id}/pay"

	// Expenditure report routes
	PurchasesSummaryURL = "/purchases/reports/purchases-summary"
	ExpensesSummaryURL  = "/expenses/reports/expenses-summary"

	// Expenditure settings (template management) routes
	ExpenditureSettingsTemplatesURL       = "/purchases/settings/templates"
	ExpenditureSettingsTemplateUploadURL  = "/action/purchase/settings/templates/upload"
	ExpenditureSettingsTemplateDeleteURL  = "/action/purchase/settings/templates/delete"
	ExpenditureSettingsTemplateDefaultURL = "/action/purchase/settings/templates/set-default/{id}"

	// Purchase Order routes
	PurchaseOrderListURL             = "/purchase-orders/list/{status}"
	PurchaseOrderDetailURL           = "/purchase-orders/detail/{id}"
	PurchaseOrderAddURL              = "/action/purchase-order/add"
	PurchaseOrderEditURL             = "/action/purchase-order/edit/{id}"
	PurchaseOrderDeleteURL           = "/action/purchase-order/delete"
	PurchaseOrderSetStatusURL        = "/action/purchase-order/set-status"
	PurchaseOrderTableURL            = "/action/purchase-order/table/{status}"
	PurchaseOrderTabActionURL        = "/action/purchase-order/detail/{id}/tab/{tab}"
	PurchaseOrderAttachmentUploadURL = "/action/purchase-order/detail/{id}/attachments/upload"
	PurchaseOrderAttachmentDeleteURL = "/action/purchase-order/detail/{id}/attachments/delete"

	// Purchase Order line item routes (within PO detail)
	PurchaseOrderLineItemTableURL  = "/action/purchase-order/detail/{id}/items/table"
	PurchaseOrderLineItemAddURL    = "/action/purchase-order/detail/{id}/items/add"
	PurchaseOrderLineItemEditURL   = "/action/purchase-order/detail/{id}/items/edit/{itemId}"
	PurchaseOrderLineItemRemoveURL = "/action/purchase-order/detail/{id}/items/remove"

	// Purchase Order receipt action
	PurchaseOrderConfirmReceiptURL = "/action/purchase-order/{id}/confirm-receipt"

	// Expense category settings routes
	ExpenditureExpenseCategoryListURL   = "/expenses/categories/list"
	ExpenditureExpenseCategoryAddURL    = "/action/expense/categories/add"
	ExpenditureExpenseCategoryEditURL   = "/action/expense/categories/edit/{id}"
	ExpenditureExpenseCategoryDeleteURL = "/action/expense/categories/delete"
	ExpenditureExpenseCategoryTableURL  = "/action/expense/categories/table"

	// ---------------------------------------------------------------------------
	// P3a — SupplierContract + SupplierContractLine route constants
	// ---------------------------------------------------------------------------

	// SupplierContract master routes
	SupplierContractListURL             = "/supplier-contracts/list/{status}"
	SupplierContractDetailURL           = "/supplier-contracts/detail/{id}"
	SupplierContractAddURL              = "/action/supplier-contract/add"
	SupplierContractEditURL             = "/action/supplier-contract/edit/{id}"
	SupplierContractDeleteURL           = "/action/supplier-contract/delete"
	SupplierContractSetStatusURL        = "/action/supplier-contract/set-status"
	SupplierContractBulkSetStatusURL    = "/action/supplier-contract/bulk-set-status"
	SupplierContractTabActionURL        = "/action/supplier-contract/detail/{id}/tab/{tab}"
	SupplierContractAttachmentUploadURL = "/action/supplier-contract/detail/{id}/attachments/upload"
	SupplierContractAttachmentDeleteURL = "/action/supplier-contract/detail/{id}/attachments/delete"
	SupplierContractApproveURL          = "/action/supplier-contract/approve/{id}"
	SupplierContractTerminateURL        = "/action/supplier-contract/terminate/{id}"

	// SupplierContractLine routes (child of contract detail)
	SupplierContractLineAddURL    = "/action/supplier-contract/{id}/lines/add"
	SupplierContractLineEditURL   = "/action/supplier-contract/{id}/lines/edit/{lid}"
	SupplierContractLineDeleteURL = "/action/supplier-contract/{id}/lines/delete"

	// ---------------------------------------------------------------------------
	// P3a — ProcurementRequest + ProcurementRequestLine route constants
	// ---------------------------------------------------------------------------

	// ProcurementRequest routes
	ProcurementRequestListURL             = "/procurement-requests/list/{status}"
	ProcurementRequestDetailURL           = "/procurement-requests/detail/{id}"
	ProcurementRequestAddURL              = "/action/procurement-request/add"
	ProcurementRequestEditURL             = "/action/procurement-request/edit/{id}"
	ProcurementRequestDeleteURL           = "/action/procurement-request/delete"
	ProcurementRequestSetStatusURL        = "/action/procurement-request/set-status"
	ProcurementRequestBulkSetStatusURL    = "/action/procurement-request/bulk-set-status"
	ProcurementRequestTabActionURL        = "/action/procurement-request/detail/{id}/tab/{tab}"
	ProcurementRequestAttachmentUploadURL = "/action/procurement-request/detail/{id}/attachments/upload"
	ProcurementRequestAttachmentDeleteURL = "/action/procurement-request/detail/{id}/attachments/delete"
	ProcurementRequestSubmitURL           = "/action/procurement-request/submit/{id}"
	ProcurementRequestApproveURL          = "/action/procurement-request/approve/{id}"
	ProcurementRequestRejectURL           = "/action/procurement-request/reject/{id}"
	ProcurementRequestSpawnPOURL          = "/action/procurement-request/spawn-po/{id}"

	// ProcurementRequestLine routes (child of request detail)
	ProcurementRequestLineAddURL    = "/action/procurement-request/{id}/lines/add"
	ProcurementRequestLineEditURL   = "/action/procurement-request/{id}/lines/edit/{lid}"
	ProcurementRequestLineDeleteURL = "/action/procurement-request/{id}/lines/delete"

	// SPS Wave 3 — CRIT-3 spawn-retry placeholder route. Wired into the line-row
	// "Retry" button; the actual retry use case lands in a later wave so the
	// handler is currently a no-op redirect (see action/action.go::NewRetrySpawnAction).
	// NOTE: pattern uses `/retry-spawn/{lid}` (not `/{lid}/retry-spawn`) to avoid
	// stdlib ServeMux conflict with the existing `/lines/edit/{lid}` pattern.
	ProcurementRequestLineRetrySpawnURL = "/action/procurement-request/{id}/lines/retry-spawn/{lid}"

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
	// SPS P7 — SupplierContractPriceSchedule + SupplierContractPriceScheduleLine
	// ---------------------------------------------------------------------------

	// SupplierContractPriceSchedule master routes
	SupplierContractPriceScheduleListURL             = "/supplier-contract-price-schedules/list/{status}"
	SupplierContractPriceScheduleDetailURL           = "/supplier-contract-price-schedules/detail/{id}"
	SupplierContractPriceScheduleAddURL              = "/action/supplier-contract-price-schedule/add"
	SupplierContractPriceScheduleEditURL             = "/action/supplier-contract-price-schedule/edit/{id}"
	SupplierContractPriceScheduleDeleteURL           = "/action/supplier-contract-price-schedule/delete"
	SupplierContractPriceScheduleSetStatusURL        = "/action/supplier-contract-price-schedule/set-status"
	SupplierContractPriceScheduleBulkSetStatusURL    = "/action/supplier-contract-price-schedule/bulk-set-status"
	SupplierContractPriceScheduleTabActionURL        = "/action/supplier-contract-price-schedule/detail/{id}/tab/{tab}"
	SupplierContractPriceScheduleAttachmentUploadURL = "/action/supplier-contract-price-schedule/detail/{id}/attachments/upload"
	SupplierContractPriceScheduleAttachmentDeleteURL = "/action/supplier-contract-price-schedule/detail/{id}/attachments/delete"
	SupplierContractPriceScheduleActivateURL         = "/action/supplier-contract-price-schedule/activate/{id}"
	SupplierContractPriceScheduleSupersedeURL        = "/action/supplier-contract-price-schedule/supersede/{id}"

	// SupplierContractPriceScheduleLine routes (child of schedule detail)
	SupplierContractPriceScheduleLineAddURL    = "/action/supplier-contract-price-schedule/{id}/lines/add"
	SupplierContractPriceScheduleLineEditURL   = "/action/supplier-contract-price-schedule/{id}/lines/edit/{lid}"
	SupplierContractPriceScheduleLineDeleteURL = "/action/supplier-contract-price-schedule/{id}/lines/delete"

	// ---------------------------------------------------------------------------
	// SPS P10 — ExpenseRecognition + ExpenseRecognitionLine route constants
	// ---------------------------------------------------------------------------

	// ExpenseRecognition master routes (no add/edit drawer — created BY use case)
	ExpenseRecognitionListURL                     = "/expense-recognitions/list/{status}"
	ExpenseRecognitionDetailURL                   = "/expense-recognitions/detail/{id}"
	ExpenseRecognitionDeleteURL                   = "/action/expense-recognition/delete"
	ExpenseRecognitionTabActionURL                = "/action/expense-recognition/detail/{id}/tab/{tab}"
	ExpenseRecognitionAttachmentUploadURL         = "/action/expense-recognition/detail/{id}/attachments/upload"
	ExpenseRecognitionAttachmentDeleteURL         = "/action/expense-recognition/detail/{id}/attachments/delete"
	ExpenseRecognitionReverseURL                  = "/action/expense-recognition/reverse/{id}"
	ExpenseRecognitionRecognizeFromExpenditureURL = "/action/expense-recognition/recognize-from-expenditure"
	ExpenseRecognitionRecognizeFromContractURL    = "/action/expense-recognition/recognize-from-contract"

	// ExpenseRecognitionLine routes (child of recognition detail — inline CRUD)
	ExpenseRecognitionLineAddURL    = "/action/expense-recognition/{id}/lines/add"
	ExpenseRecognitionLineEditURL   = "/action/expense-recognition/{id}/lines/edit/{lid}"
	ExpenseRecognitionLineDeleteURL = "/action/expense-recognition/{id}/lines/delete"

	// ---------------------------------------------------------------------------
	// SPS P10 — AccruedExpense + AccruedExpenseSettlement route constants
	// ---------------------------------------------------------------------------

	// AccruedExpense master routes (manual create drawer is secondary — primary path is AccrueFromContract use case)
	AccruedExpenseListURL               = "/accrued-expenses/list/{status}"
	AccruedExpenseDetailURL             = "/accrued-expenses/detail/{id}"
	AccruedExpenseAddURL                = "/action/accrued-expense/add"
	AccruedExpenseEditURL               = "/action/accrued-expense/edit/{id}"
	AccruedExpenseDeleteURL             = "/action/accrued-expense/delete"
	AccruedExpenseSetStatusURL          = "/action/accrued-expense/set-status"
	AccruedExpenseBulkSetStatusURL      = "/action/accrued-expense/bulk-set-status"
	AccruedExpenseTabActionURL          = "/action/accrued-expense/detail/{id}/tab/{tab}"
	AccruedExpenseAttachmentUploadURL   = "/action/accrued-expense/detail/{id}/attachments/upload"
	AccruedExpenseAttachmentDeleteURL   = "/action/accrued-expense/detail/{id}/attachments/delete"
	AccruedExpenseSettleURL             = "/action/accrued-expense/settle/{id}"
	AccruedExpenseReverseURL            = "/action/accrued-expense/reverse/{id}"
	AccruedExpenseAccrueFromContractURL = "/action/accrued-expense/accrue-from-contract"

	// AccruedExpenseSettlement routes (child of accrual detail — inline CRUD)
	AccruedExpenseSettlementAddURL    = "/action/accrued-expense/{id}/settlements/add"
	AccruedExpenseSettlementEditURL   = "/action/accrued-expense/{id}/settlements/edit/{sid}"
	AccruedExpenseSettlementDeleteURL = "/action/accrued-expense/{id}/settlements/delete"

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
