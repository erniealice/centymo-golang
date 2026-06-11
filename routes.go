package centymo

// Default route constants for centymo views.
// Consumer apps can use these or define their own.
const (
	PlanListURL             = "/plans/list/{status}"
	PlanDetailURL           = "/plans/detail/{id}"
	PlanAddURL              = "/action/plan/add"
	PlanEditURL             = "/action/plan/edit/{id}"
	PlanDeleteURL           = "/action/plan/delete"
	PlanBulkDeleteURL       = "/action/plan/bulk-delete"
	PlanSetStatusURL        = "/action/plan/set-status"
	PlanBulkSetStatusURL    = "/action/plan/bulk-set-status"
	PlanTableURL            = "/action/plan/table/{status}"
	PlanTabActionURL        = "/action/plan/detail/{id}/tab/{tab}"
	PlanAttachmentUploadURL = "/action/plan/detail/{id}/attachments/upload"
	PlanAttachmentDeleteURL = "/action/plan/detail/{id}/attachments/delete"

	// PricePlan CRUD routes (within plan context)
	PricePlanAddURL    = "/action/plan/{id}/price-plans/add"
	PricePlanEditURL   = "/action/plan/{id}/price-plans/edit/{ppid}"
	PricePlanDeleteURL = "/action/plan/{id}/price-plans/delete"

	// Plan-scoped PricePlan detail — mirrors PriceSchedulePlanDetailURL but
	// keeps users in the Package (Plan) URL namespace so ActiveNav stays
	// anchored to the Services accordion's Packages section.
	// {id}=plan id, {ppid}=price_plan id.
	PlanPricePlanDetailURL             = "/plans/detail/{id}/price/{ppid}"
	PlanPricePlanTabActionURL          = "/action/plan/{id}/price/{ppid}/tab/{tab}"
	PlanPricePlanEditURL               = "/action/plan/{id}/price/{ppid}/edit"
	PlanPricePlanDeleteURL             = "/action/plan/{id}/price/{ppid}/delete"
	PlanPricePlanProductPriceAddURL    = "/action/plan/{id}/price/{ppid}/product-prices/add"
	PlanPricePlanProductPriceEditURL   = "/action/plan/{id}/price/{ppid}/product-prices/edit/{pppid}"
	PlanPricePlanProductPriceDeleteURL = "/action/plan/{id}/price/{ppid}/product-prices/delete"

	// ProductPlan CRUD routes (within plan context)
	PlanProductPlanAddURL    = "/action/plan/{id}/products/add"
	PlanProductPlanEditURL   = "/action/plan/{id}/products/edit/{ppid}"
	PlanProductPlanDeleteURL = "/action/plan/{id}/products/delete"
	PlanProductPlanPickerURL = "/action/plan/{id}/products/picker"

	// PricePlan standalone routes (rate cards as independent entity)
	PricePlanDashboardURL        = "/price-plans/dashboard"
	PricePlanListURL             = "/price-plans/list/{status}"
	PricePlanTableURL            = "/action/price-plan/table/{status}"
	PricePlanDetailURL           = "/price-plans/detail/{id}"
	PricePlanStandaloneAddURL    = "/action/price-plan/add"
	PricePlanStandaloneEditURL   = "/action/price-plan/edit/{id}"
	PricePlanStandaloneDeleteURL = "/action/price-plan/delete"
	PricePlanBulkDeleteURL       = "/action/price-plan/bulk-delete"
	PricePlanSetStatusURL        = "/action/price-plan/set-status"
	PricePlanBulkSetStatusURL    = "/action/price-plan/bulk-set-status"
	PricePlanTabActionURL        = "/action/price-plan/{id}/tab/{tab}"
	PricePlanAttachmentUploadURL = "/action/price-plan/{id}/attachments/upload"
	PricePlanAttachmentDeleteURL = "/action/price-plan/{id}/attachments/delete"

	// ProductPricePlan CRUD routes (within price plan / rate card detail)
	PricePlanProductPriceAddURL    = "/action/price-plan/{id}/product-prices/add"
	PricePlanProductPriceEditURL   = "/action/price-plan/{id}/product-prices/edit/{ppid}"
	PricePlanProductPriceDeleteURL = "/action/price-plan/{id}/product-prices/delete"

	// PriceSchedule routes (date-bounded pricing container for plans)
	PriceScheduleDashboardURL        = "/price-schedules/dashboard"
	PriceScheduleListURL             = "/price-schedules/list/{status}"
	PriceScheduleTableURL            = "/action/price-schedule/table/{status}"
	PriceScheduleDetailURL           = "/price-schedules/detail/{id}"
	PriceScheduleAddURL              = "/action/price-schedule/add"
	PriceScheduleEditURL             = "/action/price-schedule/edit/{id}"
	PriceScheduleDeleteURL           = "/action/price-schedule/delete"
	PriceScheduleBulkDeleteURL       = "/action/price-schedule/bulk-delete"
	PriceScheduleSetStatusURL        = "/action/price-schedule/set-status"
	PriceScheduleBulkSetStatusURL    = "/action/price-schedule/bulk-set-status"
	PriceScheduleTabActionURL        = "/action/price-schedule/{id}/tab/{tab}"
	PriceScheduleAttachmentUploadURL = "/action/price-schedule/{id}/attachments/upload"
	PriceScheduleAttachmentDeleteURL = "/action/price-schedule/{id}/attachments/delete"
	PriceSchedulePlanAddURL          = "/action/price-schedule/{id}/plan/add"
	// Schedule-scoped price_plan detail. Mirrors /app/price-plans/detail/{id} but nests
	// under the schedule so sidebar context stays on price-schedules (price_plan is no
	// longer a top-level sidebar entry).
	PriceSchedulePlanDetailURL             = "/price-schedules/detail/{id}/plan/{ppid}"
	PriceSchedulePlanTabActionURL          = "/action/price-schedule/{id}/plan/{ppid}/tab/{tab}"
	PriceSchedulePlanEditURL               = "/action/price-schedule/{id}/plan/{ppid}/edit"
	PriceSchedulePlanDeleteURL             = "/action/price-schedule/{id}/plan/{ppid}/delete"
	PriceSchedulePlanProductPriceAddURL    = "/action/price-schedule/{id}/plan/{ppid}/product-prices/add"
	PriceSchedulePlanProductPriceEditURL   = "/action/price-schedule/{id}/plan/{ppid}/product-prices/edit/{pppid}"
	PriceSchedulePlanProductPriceDeleteURL = "/action/price-schedule/{id}/plan/{ppid}/product-prices/delete"

	// Attachments on the nested price_schedule/plan detail page.
	PriceSchedulePlanAttachmentUploadURL = "/action/price-schedule/{id}/plan/{ppid}/attachments/upload"
	PriceSchedulePlanAttachmentDeleteURL = "/action/price-schedule/{id}/plan/{ppid}/attachments/delete"

	// 2026-05-04 — Subscriptions tab on the schedule-scoped
	// price_plan detail page. Same handler as SubscriptionDetailURL; the
	// nested URL alone activates the rate-card → plan → subscription breadcrumb.
	// See docs/plan/20260504-price-plan-engagements-tab/.
	PriceSchedulePlanSubscriptionDetailURL = "/price-schedules/detail/{id}/plan/{ppid}/subscription/{eid}"

	SubscriptionListURL = "/subscriptions/list/{status}"
	// SubscriptionTableURL returns ONLY the table-card partial — used as the
	// data-refresh-url so HTMX swaps the table without re-rendering the whole page.
	SubscriptionTableURL  = "/action/subscription/table/{status}"
	SubscriptionDetailURL = "/subscriptions/detail/{id}"
	// SubscriptionUnderClientDetailURL is the nested subscription-detail path
	// rendered with a client breadcrumb. Same view as SubscriptionDetailURL.
	SubscriptionUnderClientDetailURL  = "/clients/detail/{client_id}/subscriptions/{id}"
	SubscriptionAddURL                = "/action/subscription/add"
	SubscriptionEditURL               = "/action/subscription/edit/{id}"
	SubscriptionDeleteURL             = "/action/subscription/delete"
	SubscriptionBulkDeleteURL         = "/action/subscription/bulk-delete"
	SubscriptionSetStatusURL          = "/action/subscription/set-status"
	SubscriptionBulkSetStatusURL      = "/action/subscription/bulk-set-status"
	SubscriptionTabActionURL          = "/action/subscription/detail/{id}/tab/{tab}"
	SubscriptionAttachmentUploadURL   = "/action/subscription/detail/{id}/attachments/upload"
	SubscriptionAttachmentDeleteURL   = "/action/subscription/detail/{id}/attachments/delete"
	SubscriptionAttachmentDownloadURL = "/action/subscription/detail/{id}/attachments/download"
	SubscriptionSearchPlanURL         = "/action/subscription/search/plans"
	SubscriptionSearchClientURL       = "/action/subscription/search/clients"
	// SubscriptionRecognizeURL opens the "Recognize Revenue" drawer for a
	// subscription. GET = preview drawer (dry_run); POST = generate the Revenue.
	// Verb-first to avoid Go ServeMux ambiguity with /action/subscription/edit/{id}
	// — id-first and static-prefix patterns at the same depth can't disambiguate
	// (e.g. "/action/subscription/edit/recognize-revenue" matches both).
	SubscriptionRecognizeURL = "/action/subscription/recognize-revenue/{id}"

	// SubscriptionRevenueRunURL opens the "Invoice Run" drawer for a single
	// subscription (Surface C — per-subscription drawer, CYCLE billing_kind only).
	// Verb-first ("revenue-run") to avoid ServeMux ambiguity consistent with
	// SubscriptionRecognizeURL above.
	SubscriptionRevenueRunURL = "/action/subscription/revenue-run/{id}"

	// SubscriptionCustomizePackageURL is the POST endpoint that drives the
	// "Customize this package for {ClientName}" CTA on the subscription
	// detail's Package tab. Calls espyna's CustomizePlanForClient use case
	// and HX-redirects to the new (cloned) PricePlan's package page.
	// Verb-first ("customize-package") to avoid the same ServeMux ambiguity
	// SubscriptionRecognizeURL above guards against.
	SubscriptionCustomizePackageURL = "/action/subscription/customize-package/{id}"

	// 2026-04-29 milestone-billing plan §5 / Phase D — mark-ready + waive
	// handlers for BillingEvent rows on the subscription Package tab.
	// Both POST through the espyna BillingEvent.SetStatus domain service.
	MilestoneMarkReadyURL = "/action/subscription/{id}/billing-event/{eventId}/mark-ready"
	MilestoneWaiveURL     = "/action/subscription/{id}/billing-event/{eventId}/waive"

	// 20260517-advance-cash-events Plan B Phase 7 — Recognize handler for a
	// BillingEvent row when it is linked to a MILESTONE advance Collection via
	// the collection_billing_event junction. POSTs through the
	// espyna RecognizeMilestoneAdvanceCollection use case.
	MilestoneRecognizeURL = "/action/subscription/{id}/billing-event/{eventId}/recognize"

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — retroactive
	// spawn drawer endpoint (GET = drawer, POST = spawn) and HTMX-driven
	// partial that re-renders the Spawn Jobs section in the create drawer
	// when the operator changes the selected Plan / PricePlan.
	//
	// Verb-first ("spawn-jobs") to avoid the Go ServeMux ambiguity that would
	// otherwise pit `{subscriptionId}/spawn-jobs` against
	// `table/{status}` (and similar id-first/static-prefix patterns at the
	// same depth — same root cause as SubscriptionRecognizeURL above).
	SubscriptionSpawnJobsURL        = "/action/subscription/spawn-jobs/{subscriptionId}"
	SubscriptionSpawnJobsPartialURL = "/action/subscription/_partial/spawn-jobs-section"

	// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D — manual cycle
	// spawn + backfill triggers. Both routes call into espyna's
	// MaterializeInstanceJobsForSubscription consumer (single-cycle vs.
	// multi-cycle modes). Verb-first ("spawn-cycle-jobs", "backfill-cycle-jobs")
	// to keep ServeMux disambiguation consistent with existing
	// SubscriptionRecognizeURL / SubscriptionSpawnJobsURL.
	SubscriptionSpawnCycleJobsURL    = "/action/subscription/spawn-cycle-jobs/{subscriptionId}"
	SubscriptionBackfillCycleJobsURL = "/action/subscription/backfill-cycle-jobs/{subscriptionId}"

	// 2026-05-01 ad-hoc-subscription-billing plan §5.2 — operator-driven CTA
	// for AD_HOC subscriptions. Pool-Generate-Invoice reuses the existing
	// SubscriptionRecognizeURL; Extend-Pool deferred to v1.5.5 (needs new
	// espyna use case for Subscription.entitled_occurrences_override write).
	SubscriptionRequestUsageURL = "/action/subscription/request-usage/{subscriptionId}"

	// Collection (money IN) routes
	CollectionListURL             = "/collections/list/{status}"
	CollectionDetailURL           = "/collections/detail/{id}"
	CollectionDashboardURL        = "/collections/dashboard"
	CollectionAddURL              = "/action/collection/add"
	CollectionEditURL             = "/action/collection/edit/{id}"
	CollectionDeleteURL           = "/action/collection/delete"
	CollectionBulkDeleteURL       = "/action/collection/bulk-delete"
	CollectionSetStatusURL        = "/action/collection/set-status"
	CollectionBulkSetStatusURL    = "/action/collection/bulk-set-status"
	CollectionTabActionURL        = "/action/collection/detail/{id}/tab/{tab}"
	CollectionAttachmentUploadURL = "/action/collection/detail/{id}/attachments/upload"
	CollectionAttachmentDeleteURL = "/action/collection/detail/{id}/attachments/delete"

	// Disbursement (money OUT) routes
	DisbursementListURL             = "/disbursements/list/{status}"
	DisbursementDetailURL           = "/disbursements/detail/{id}"
	DisbursementDashboardURL        = "/disbursements/dashboard"
	DisbursementAddURL              = "/action/disbursement/add"
	DisbursementEditURL             = "/action/disbursement/edit/{id}"
	DisbursementDeleteURL           = "/action/disbursement/delete"
	DisbursementBulkDeleteURL       = "/action/disbursement/bulk-delete"
	DisbursementSetStatusURL        = "/action/disbursement/set-status"
	DisbursementBulkSetStatusURL    = "/action/disbursement/bulk-set-status"
	DisbursementTabActionURL        = "/action/disbursement/detail/{id}/tab/{tab}"
	DisbursementAttachmentUploadURL = "/action/disbursement/detail/{id}/attachments/upload"
	DisbursementAttachmentDeleteURL = "/action/disbursement/detail/{id}/attachments/delete"

	// ---------------------------------------------------------------------------
	// 20260517-advance-cash-events Plan B Phase 3 — Advance Cash Events routes.
	// "Advances" is a Cash-app section that surfaces TreasuryCollection /
	// TreasuryDisbursement rows whose advance_kind != NONE plus a workspace
	// dashboard. These are first-class operator actions (Settle / Refund /
	// Cancel) anchored on the existing TreasuryCollection / TreasuryDisbursement
	// detail pages — there is no separate "advance" entity.
	// ---------------------------------------------------------------------------

	// Advances Dashboard — workspace-level summary (both sides).
	AdvancesDashboardURL = "/cash/advances/dashboard"

	// Filtered list URLs (advance_kind != NONE) — point at the existing
	// Collection / Disbursement list pages with the chip pre-applied via a
	// query string the list page interprets. These are sidebar Href targets,
	// NOT ServeMux patterns — the list pages are registered at the underlying
	// pattern (CollectionListURL / DisbursementListURL) and read advance_kind
	// from the request query string.
	AdvanceCollectionListURL   = "/collections/list/pending?advance_kind=any"
	AdvanceDisbursementListURL = "/disbursements/list/pending?advance_kind=any"

	// TreasuryCollection / TreasuryDisbursement Advance Schedule tab partials
	// (loaded via HTMX, sit beside info / attachments / audit / advance-schedule
	// values in the detail page's tab switch).
	TreasuryCollectionAdvanceScheduleTabURL   = "/action/collection/detail/{id}/tab/advance-schedule"
	TreasuryDisbursementAdvanceScheduleTabURL = "/action/disbursement/detail/{id}/tab/advance-schedule"

	// UNSCHEDULED workflow drawers — Settle / Refund / Cancel on both sides.
	// Verb-first to avoid Go ServeMux ambiguity with the existing edit/{id}
	// patterns at the same depth (same rationale as SubscriptionRecognizeURL).
	TreasuryCollectionSettleURL   = "/action/collection/settle/{id}"
	TreasuryCollectionRefundURL   = "/action/collection/refund/{id}"
	TreasuryCollectionCancelURL   = "/action/collection/cancel/{id}"
	TreasuryDisbursementSettleURL = "/action/disbursement/settle/{id}"
	TreasuryDisbursementRefundURL = "/action/disbursement/refund/{id}"
	TreasuryDisbursementCancelURL = "/action/disbursement/cancel/{id}"

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
