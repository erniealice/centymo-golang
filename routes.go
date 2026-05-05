package centymo

// Default route constants for centymo views.
// Consumer apps can use these or define their own.
const (
	PlanListURL             = "/app/plans/list/{status}"
	PlanDetailURL           = "/app/plans/detail/{id}"
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
	PlanPricePlanDetailURL             = "/app/plans/detail/{id}/price/{ppid}"
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
	PricePlanDashboardURL        = "/app/price-plans/dashboard"
	PricePlanListURL             = "/app/price-plans/list/{status}"
	PricePlanTableURL            = "/action/price-plan/table/{status}"
	PricePlanDetailURL           = "/app/price-plans/detail/{id}"
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
	PriceScheduleDashboardURL        = "/app/price-schedules/dashboard"
	PriceScheduleListURL             = "/app/price-schedules/list/{status}"
	PriceScheduleTableURL            = "/action/price-schedule/table/{status}"
	PriceScheduleDetailURL           = "/app/price-schedules/detail/{id}"
	PriceScheduleAddURL              = "/action/price-schedule/add"
	PriceScheduleEditURL             = "/action/price-schedule/edit/{id}"
	PriceScheduleDeleteURL           = "/action/price-schedule/delete"
	PriceScheduleBulkDeleteURL       = "/action/price-schedule/bulk-delete"
	PriceScheduleSetStatusURL        = "/action/price-schedule/set-status"
	PriceScheduleBulkSetStatusURL    = "/action/price-schedule/bulk-set-status"
	PriceScheduleTabActionURL        = "/action/price-schedule/{id}/tab/{tab}"
	PriceSchedulePlanAddURL          = "/action/price-schedule/{id}/plan/add"
	// Schedule-scoped price_plan detail. Mirrors /app/price-plans/detail/{id} but nests
	// under the schedule so sidebar context stays on price-schedules (price_plan is no
	// longer a top-level sidebar entry).
	PriceSchedulePlanDetailURL             = "/app/price-schedules/detail/{id}/plan/{ppid}"
	PriceSchedulePlanTabActionURL          = "/action/price-schedule/{id}/plan/{ppid}/tab/{tab}"
	PriceSchedulePlanEditURL               = "/action/price-schedule/{id}/plan/{ppid}/edit"
	PriceSchedulePlanDeleteURL             = "/action/price-schedule/{id}/plan/{ppid}/delete"
	PriceSchedulePlanProductPriceAddURL    = "/action/price-schedule/{id}/plan/{ppid}/product-prices/add"
	PriceSchedulePlanProductPriceEditURL   = "/action/price-schedule/{id}/plan/{ppid}/product-prices/edit/{pppid}"
	PriceSchedulePlanProductPriceDeleteURL = "/action/price-schedule/{id}/plan/{ppid}/product-prices/delete"

	// 2026-05-04 — Engagements (subscriptions) tab on the schedule-scoped
	// price_plan detail page. Same handler as SubscriptionDetailURL; the
	// nested URL alone activates the rate-card → plan → engagement breadcrumb.
	// See docs/plan/20260504-price-plan-engagements-tab/.
	PriceSchedulePlanEngagementDetailURL = "/app/price-schedules/detail/{id}/plan/{ppid}/engagement/{eid}"

	SubscriptionListURL             = "/app/subscriptions/list/{status}"
	// SubscriptionTableURL returns ONLY the table-card partial — used as the
	// data-refresh-url so HTMX swaps the table without re-rendering the whole page.
	SubscriptionTableURL            = "/action/subscription/table/{status}"
	SubscriptionDetailURL           = "/app/subscriptions/detail/{id}"
	// SubscriptionUnderClientDetailURL is the nested subscription-detail path
	// rendered with a client breadcrumb. Same view as SubscriptionDetailURL.
	SubscriptionUnderClientDetailURL = "/app/clients/detail/{client_id}/subscriptions/{id}"
	SubscriptionAddURL              = "/action/subscription/add"
	SubscriptionEditURL             = "/action/subscription/edit/{id}"
	SubscriptionDeleteURL           = "/action/subscription/delete"
	SubscriptionBulkDeleteURL       = "/action/subscription/bulk-delete"
	SubscriptionSetStatusURL        = "/action/subscription/set-status"
	SubscriptionBulkSetStatusURL    = "/action/subscription/bulk-set-status"
	SubscriptionTabActionURL        = "/action/subscription/detail/{id}/tab/{tab}"
	SubscriptionAttachmentUploadURL = "/action/subscription/detail/{id}/attachments/upload"
	SubscriptionAttachmentDeleteURL = "/action/subscription/detail/{id}/attachments/delete"
	SubscriptionSearchPlanURL       = "/action/subscription/search/plans"
	SubscriptionSearchClientURL     = "/action/subscription/search/clients"
	// SubscriptionRecognizeURL opens the "Recognize Revenue" drawer for a
	// subscription. GET = preview drawer (dry_run); POST = generate the Revenue.
	// Verb-first to avoid Go ServeMux ambiguity with /action/subscription/edit/{id}
	// — id-first and static-prefix patterns at the same depth can't disambiguate
	// (e.g. "/action/subscription/edit/recognize-revenue" matches both).
	SubscriptionRecognizeURL        = "/action/subscription/recognize-revenue/{id}"

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
	CollectionListURL             = "/app/collections/list/{status}"
	CollectionDetailURL           = "/app/collections/detail/{id}"
	CollectionDashboardURL        = "/app/collections/dashboard"
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
	DisbursementListURL             = "/app/disbursements/list/{status}"
	DisbursementDetailURL           = "/app/disbursements/detail/{id}"
	DisbursementDashboardURL        = "/app/disbursements/dashboard"
	DisbursementAddURL              = "/action/disbursement/add"
	DisbursementEditURL             = "/action/disbursement/edit/{id}"
	DisbursementDeleteURL           = "/action/disbursement/delete"
	DisbursementBulkDeleteURL       = "/action/disbursement/bulk-delete"
	DisbursementSetStatusURL        = "/action/disbursement/set-status"
	DisbursementBulkSetStatusURL    = "/action/disbursement/bulk-set-status"
	DisbursementTabActionURL        = "/action/disbursement/detail/{id}/tab/{tab}"
	DisbursementAttachmentUploadURL = "/action/disbursement/detail/{id}/attachments/upload"
	DisbursementAttachmentDeleteURL = "/action/disbursement/detail/{id}/attachments/delete"

	ProductListURL   = "/app/products/list/{status}"
	ProductTableURL  = "/action/product/table/{status}"
	ProductDetailURL = "/app/products/detail/{id}"

	// Service dashboard — services are products filtered to product_kind="service".
	// The dashboard sits ABOVE the service-mount product list at this URL.
	ServiceDashboardURL = "/app/services/dashboard"

	// Inventory routes — list is per-location
	InventoryDashboardURL  = "/app/inventory/dashboard"
	InventoryListURL       = "/app/inventory/list/{location}"
	InventoryAddURL        = "/action/inventory/add"
	InventoryEditURL       = "/action/inventory/edit/{id}"
	InventoryDeleteURL     = "/action/inventory/delete"
	InventoryBulkDeleteURL = "/action/inventory/bulk-delete"
	InventoryDetailURL     = "/app/inventory/detail/{id}"

	// Product action routes
	ProductAddURL        = "/action/product/add"
	ProductEditURL       = "/action/product/edit/{id}"
	ProductDeleteURL     = "/action/product/delete"
	ProductBulkDeleteURL = "/action/product/bulk-delete"

	// Product status routes
	ProductSetStatusURL     = "/action/product/set-status"
	ProductBulkSetStatusURL = "/action/product/bulk-set-status"

	// Product detail tab action route
	ProductTabActionURL        = "/action/product/detail/{id}/tab/{tab}"
	ProductAttachmentUploadURL = "/action/product/detail/{id}/attachments/upload"
	ProductAttachmentDeleteURL = "/action/product/detail/{id}/attachments/delete"

	// Product variant routes (within product detail)
	ProductVariantTableURL  = "/action/product/detail/{id}/variants/table"
	ProductVariantAssignURL = "/action/product/detail/{id}/variants/assign"
	ProductVariantEditURL   = "/action/product/detail/{id}/variants/edit/{vid}"
	ProductVariantRemoveURL = "/action/product/detail/{id}/variants/remove"

	// Product attribute routes (within product detail)
	ProductAttributeTableURL  = "/action/product/detail/{id}/attributes/table"
	ProductAttributeAssignURL = "/action/product/detail/{id}/attributes/assign"
	ProductAttributeRemoveURL = "/action/product/detail/{id}/attributes/remove"

	// Product option routes (within product detail)
	ProductOptionTableURL  = "/action/product/detail/{id}/options/table"
	ProductOptionAddURL    = "/action/product/detail/{id}/options/add"
	ProductOptionEditURL   = "/action/product/detail/{id}/options/edit/{oid}"
	ProductOptionDeleteURL = "/action/product/detail/{id}/options/delete"

	// Product option detail page (option values management)
	ProductOptionDetailURL = "/app/products/detail/{id}/option/{oid}"

	// Product variant detail page (variant info, pricing, stock, audit, images)
	ProductVariantDetailURL    = "/app/products/detail/{id}/variant/{vid}"
	ProductVariantTabActionURL = "/action/product/detail/{id}/variant/{vid}/tab/{tab}"

	// Product variant image routes (upload/delete within variant detail)
	ProductVariantImageUploadURL = "/action/product/detail/{id}/variant/{vid}/images/upload"
	ProductVariantImageDeleteURL = "/action/product/detail/{id}/variant/{vid}/images/delete"

	// Product variant attachment routes
	ProductVariantAttachmentUploadURL = "/action/product/detail/{id}/variant/{vid}/attachments/upload"
	ProductVariantAttachmentDeleteURL = "/action/product/detail/{id}/variant/{vid}/attachments/delete"

	// Product variant stock detail (inventory item within variant context)
	ProductVariantStockDetailURL    = "/app/products/detail/{id}/variant/{vid}/stock/{iid}"
	ProductVariantStockTabActionURL = "/action/product/detail/{id}/variant/{vid}/stock/{iid}/tab/{tab}"

	// Product variant stock attachment routes
	ProductVariantStockAttachmentUploadURL = "/action/product/detail/{id}/variant/{vid}/stock/{iid}/attachments/upload"
	ProductVariantStockAttachmentDeleteURL = "/action/product/detail/{id}/variant/{vid}/stock/{iid}/attachments/delete"

	// Inventory serial detail (individual serial within inventory item)
	ProductVariantSerialDetailURL = "/app/products/detail/{id}/variant/{vid}/stock/{iid}/serial/{sid}"

	// Product option value routes (within product option)
	ProductOptionValueTableURL  = "/action/product/detail/{id}/options/{oid}/values/table"
	ProductOptionValueAddURL    = "/action/product/detail/{id}/options/{oid}/values/add"
	ProductOptionValueEditURL   = "/action/product/detail/{id}/options/{oid}/values/edit/{vid}"
	ProductOptionValueDeleteURL = "/action/product/detail/{id}/options/{oid}/values/delete"

	// Product line routes
	ProductLineDashboardURL        = "/app/product-lines/dashboard"
	ProductLineListURL             = "/app/product-lines/list/{status}"
	ProductLineTableURL            = "/action/product-line/table/{status}"
	ProductLineDetailURL           = "/app/product-lines/detail/{id}"
	ProductLineAddURL              = "/action/product-line/add"
	ProductLineEditURL             = "/action/product-line/edit/{id}"
	ProductLineDeleteURL           = "/action/product-line/delete"
	ProductLineBulkDeleteURL       = "/action/product-line/bulk-delete"
	ProductLineSetStatusURL        = "/action/product-line/set-status"
	ProductLineBulkSetStatusURL    = "/action/product-line/bulk-set-status"
	ProductLineTabActionURL        = "/action/product-line/{id}/tab/{tab}"
	ProductLineAttachmentUploadURL = "/action/product-line/{id}/attachments/upload"
	ProductLineAttachmentDeleteURL = "/action/product-line/{id}/attachments/delete"

	// Inventory status routes
	InventorySetStatusURL     = "/action/inventory/set-status"
	InventoryBulkSetStatusURL = "/action/inventory/bulk-set-status"

	// Inventory tab action route
	InventoryTabActionURL        = "/action/inventory/detail/{id}/tab/{tab}"
	InventoryAttachmentUploadURL = "/action/inventory/detail/{id}/attachments/upload"
	InventoryAttachmentDeleteURL = "/action/inventory/detail/{id}/attachments/delete"

	// Inventory movements (global transaction history)
	InventoryMovementsURL       = "/app/inventory/movements"
	InventoryMovementsTableURL  = "/action/inventory/movements/table"
	InventoryMovementsExportURL = "/action/inventory/movements/export"

	// Inventory serial routes
	InventorySerialTableURL  = "/action/inventory/detail/{id}/serials/table"
	InventorySerialAssignURL = "/action/inventory/detail/{id}/serials/assign"
	InventorySerialEditURL   = "/action/inventory/detail/{id}/serials/edit/{sid}"
	InventorySerialRemoveURL = "/action/inventory/detail/{id}/serials/remove"

	// Inventory transaction routes
	InventoryTransactionTableURL  = "/action/inventory/detail/{id}/transactions/table"
	InventoryTransactionAssignURL = "/action/inventory/detail/{id}/transactions/assign"

	// Inventory depreciation routes
	InventoryDepreciationAssignURL = "/action/inventory/detail/{id}/depreciation/assign"
	InventoryDepreciationEditURL   = "/action/inventory/detail/{id}/depreciation/edit/{did}"

	// Inventory attribute routes (within detail)
	InventoryAttributeTableURL = "/action/inventory/detail/{id}/attributes/table"

	// Inventory dashboard partial routes
	InventoryDashboardStatsURL     = "/action/inventory/dashboard/stats"
	InventoryDashboardChartURL     = "/action/inventory/dashboard/chart"
	InventoryDashboardMovementsURL = "/action/inventory/dashboard/movements"
	InventoryDashboardAlertsURL    = "/action/inventory/dashboard/alerts"

	// Inventory product-context detail routes
	InventoryProductDetailURL    = "/app/products/detail/{pid}/inventory/detail/{iid}"
	InventoryProductTabActionURL = "/action/product/{pid}/inventory/{iid}/tab/{tab}"

	// Inventory table refresh (per-location)
	InventoryTableURL = "/action/inventory/table/{location}"

	// Revenue routes
	RevenueDashboardURL     = "/app/sales/dashboard"
	RevenueListURL          = "/app/sales/list/{status}"
	RevenueTableURL         = "/action/revenue/table/{status}"
	RevenueDetailURL        = "/app/sales/detail/{id}"
	RevenueAddURL           = "/action/revenue/add"
	RevenueEditURL          = "/action/revenue/edit/{id}"
	RevenueDeleteURL        = "/action/revenue/delete"
	RevenueBulkDeleteURL    = "/action/revenue/bulk-delete"
	RevenueSetStatusURL     = "/action/revenue/set-status"
	RevenueBulkSetStatusURL = "/action/revenue/bulk-set-status"

	// Revenue tab action route
	RevenueTabActionURL        = "/action/revenue/detail/{id}/tab/{tab}"
	RevenueAttachmentUploadURL = "/action/revenue/detail/{id}/attachments/upload"
	RevenueAttachmentDeleteURL = "/action/revenue/detail/{id}/attachments/delete"

	// Revenue line item routes (within revenue detail)
	RevenueLineItemTableURL    = "/action/revenue/detail/{id}/items/table"
	RevenueLineItemAddURL      = "/action/revenue/detail/{id}/items/add"
	RevenueLineItemEditURL     = "/action/revenue/detail/{id}/items/edit/{itemId}"
	RevenueLineItemRemoveURL   = "/action/revenue/detail/{id}/items/remove"
	RevenueLineItemDiscountURL = "/action/revenue/detail/{id}/items/add-discount"

	// Revenue payment routes (within revenue detail)
	RevenuePaymentTableURL  = "/action/revenue/detail/{id}/payment/table"
	RevenuePaymentAddURL    = "/action/revenue/detail/{id}/payment/add"
	RevenuePaymentEditURL   = "/action/revenue/detail/{id}/payment/edit/{pid}"
	RevenuePaymentRemoveURL = "/action/revenue/detail/{id}/payment/remove"

	// Revenue report routes
	RevenueSummaryURL = "/app/sales/reports/sales-summary"

	// Revenue invoice document routes
	RevenueInvoiceDownloadURL = "/action/revenue/detail/{id}/invoice/download"
	RevenueEmailURL           = "/action/revenue/detail/{id}/invoice/send-email"

	// Revenue settings routes (template management)
	RevenueSettingsTemplatesURL       = "/app/sales/settings/templates"
	RevenueSettingsTemplateUploadURL  = "/action/revenue/settings/templates/upload"
	RevenueSettingsTemplateDeleteURL  = "/action/revenue/settings/templates/delete"
	RevenueSettingsTemplateDefaultURL = "/action/revenue/settings/templates/set-default/{id}"
	RevenueSearchClientURL            = "/action/revenue/search/clients"
	RevenueSearchSubscriptionURL      = "/action/revenue/search/subscriptions"
	RevenueSearchLocationURL          = "/action/revenue/search/locations"
	RevenueSearchProductURL           = "/action/revenue/search/products"
	RevenuePriceLookupURL             = "/action/revenue/price-lookup"

	// Expenditure (purchase + expense) routes
	ExpenditurePurchaseListURL      = "/app/purchases/list/{status}"
	ExpenditurePurchaseDashboardURL = "/app/purchases/dashboard"
	ExpenditureExpenseListURL       = "/app/expenses/list/{status}"
	ExpenditureExpenseDashboardURL  = "/app/expenses/dashboard"

	// Expenditure expense CRUD action routes
	ExpenditureExpenseAddURL       = "/action/expense/add"
	ExpenditureExpenseEditURL      = "/action/expense/edit/{id}"
	ExpenditureExpenseDeleteURL    = "/action/expense/delete"
	ExpenditureExpenseSetStatusURL = "/action/expense/set-status"
	ExpenditureExpenseDetailURL    = "/app/expenses/detail/{id}"
	ExpenditureExpenseTableURL     = "/action/expense/table/{status}"
	ExpenditureExpenseTabActionURL = "/action/expense/detail/{id}/tab/{tab}"

	// Expenditure expense line item action routes
	ExpenditureExpenseLineItemAddURL    = "/action/expense/detail/{id}/items/add"
	ExpenditureExpenseLineItemEditURL   = "/action/expense/detail/{id}/items/edit/{itemId}"
	ExpenditureExpenseLineItemRemoveURL = "/action/expense/detail/{id}/items/remove"
	ExpenditureExpenseLineItemTableURL  = "/action/expense/detail/{id}/items/table"

	// Expenditure pay action route (creates a pre-linked disbursement)
	ExpenditureExpensePayURL = "/action/expense/detail/{id}/pay"

	// Expenditure report routes
	PurchasesSummaryURL = "/app/purchases/reports/purchases-summary"
	ExpensesSummaryURL  = "/app/expenses/reports/expenses-summary"

	// Expenditure settings (template management) routes
	ExpenditureSettingsTemplatesURL       = "/app/purchases/settings/templates"
	ExpenditureSettingsTemplateUploadURL  = "/action/purchase/settings/templates/upload"
	ExpenditureSettingsTemplateDeleteURL  = "/action/purchase/settings/templates/delete"
	ExpenditureSettingsTemplateDefaultURL = "/action/purchase/settings/templates/set-default/{id}"

	// Purchase Order routes
	PurchaseOrderListURL      = "/app/purchase-orders/list/{status}"
	PurchaseOrderDetailURL    = "/app/purchase-orders/detail/{id}"
	PurchaseOrderAddURL       = "/action/purchase-order/add"
	PurchaseOrderEditURL      = "/action/purchase-order/edit/{id}"
	PurchaseOrderDeleteURL    = "/action/purchase-order/delete"
	PurchaseOrderSetStatusURL = "/action/purchase-order/set-status"
	PurchaseOrderTableURL     = "/action/purchase-order/table/{status}"
	PurchaseOrderTabActionURL = "/action/purchase-order/detail/{id}/tab/{tab}"

	// Purchase Order line item routes (within PO detail)
	PurchaseOrderLineItemTableURL  = "/action/purchase-order/detail/{id}/items/table"
	PurchaseOrderLineItemAddURL    = "/action/purchase-order/detail/{id}/items/add"
	PurchaseOrderLineItemEditURL   = "/action/purchase-order/detail/{id}/items/edit/{itemId}"
	PurchaseOrderLineItemRemoveURL = "/action/purchase-order/detail/{id}/items/remove"

	// Purchase Order receipt action
	PurchaseOrderConfirmReceiptURL = "/action/purchase-order/{id}/confirm-receipt"

	// Expense category settings routes
	ExpenditureExpenseCategoryListURL   = "/app/expenses/categories/list"
	ExpenditureExpenseCategoryAddURL    = "/action/expense/categories/add"
	ExpenditureExpenseCategoryEditURL   = "/action/expense/categories/edit/{id}"
	ExpenditureExpenseCategoryDeleteURL = "/action/expense/categories/delete"
	ExpenditureExpenseCategoryTableURL  = "/action/expense/categories/table"

	// Resource routes (person or equipment linked to a Product for billing)
	ResourceListURL          = "/app/resources/list/{status}"
	ResourceTableURL         = "/action/resource/table/{status}"
	ResourceDetailURL        = "/app/resources/detail/{id}"
	ResourceAddURL           = "/action/resource/add"
	ResourceEditURL          = "/action/resource/edit/{id}"
	ResourceDeleteURL        = "/action/resource/delete"
	ResourceBulkDeleteURL    = "/action/resource/bulk-delete"
	ResourceSetStatusURL     = "/action/resource/set-status"
	ResourceBulkSetStatusURL = "/action/resource/bulk-set-status"

	// Price List routes — canonical home is the inventory accordion (/app/inventory/price-lists/*)
	PriceListListURL       = "/app/inventory/price-lists/list/{status}"
	PriceListTableURL      = "/action/inventory-price-list/table/{status}"
	PriceListDetailURL     = "/app/inventory/price-lists/detail/{id}"
	PriceListAddURL        = "/action/inventory-price-list/add"
	PriceListEditURL       = "/action/inventory-price-list/edit/{id}"
	PriceListDeleteURL     = "/action/inventory-price-list/delete"
	PriceListBulkDeleteURL = "/action/inventory-price-list/bulk-delete"

	PriceListTabActionURL        = "/action/inventory-price-list/{id}/tab/{tab}"
	PriceListAttachmentUploadURL = "/action/inventory-price-list/{id}/attachments/upload"
	PriceListAttachmentDeleteURL = "/action/inventory-price-list/{id}/attachments/delete"

	// Price Product routes (within price list detail)
	PriceProductAddURL    = "/action/inventory-price-list/{id}/products/add"
	PriceProductDeleteURL = "/action/inventory-price-list/{id}/products/delete"

	// ---------------------------------------------------------------------------
	// P3a — SupplierContract + SupplierContractLine route constants
	// ---------------------------------------------------------------------------

	// SupplierContract master routes
	SupplierContractListURL          = "/app/supplier-contracts/list/{status}"
	SupplierContractDetailURL        = "/app/supplier-contracts/detail/{id}"
	SupplierContractAddURL           = "/action/supplier-contract/add"
	SupplierContractEditURL          = "/action/supplier-contract/edit/{id}"
	SupplierContractDeleteURL        = "/action/supplier-contract/delete"
	SupplierContractSetStatusURL     = "/action/supplier-contract/set-status"
	SupplierContractBulkSetStatusURL = "/action/supplier-contract/bulk-set-status"
	SupplierContractTabActionURL     = "/action/supplier-contract/detail/{id}/tab/{tab}"
	SupplierContractApproveURL       = "/action/supplier-contract/approve/{id}"
	SupplierContractTerminateURL     = "/action/supplier-contract/terminate/{id}"

	// SupplierContractLine routes (child of contract detail)
	SupplierContractLineAddURL    = "/action/supplier-contract/{id}/lines/add"
	SupplierContractLineEditURL   = "/action/supplier-contract/{id}/lines/edit/{lid}"
	SupplierContractLineDeleteURL = "/action/supplier-contract/{id}/lines/delete"

	// ---------------------------------------------------------------------------
	// P3a — ProcurementRequest + ProcurementRequestLine route constants
	// ---------------------------------------------------------------------------

	// ProcurementRequest routes
	ProcurementRequestListURL          = "/app/procurement-requests/list/{status}"
	ProcurementRequestDetailURL        = "/app/procurement-requests/detail/{id}"
	ProcurementRequestAddURL           = "/action/procurement-request/add"
	ProcurementRequestEditURL          = "/action/procurement-request/edit/{id}"
	ProcurementRequestDeleteURL        = "/action/procurement-request/delete"
	ProcurementRequestSetStatusURL     = "/action/procurement-request/set-status"
	ProcurementRequestBulkSetStatusURL = "/action/procurement-request/bulk-set-status"
	ProcurementRequestTabActionURL     = "/action/procurement-request/detail/{id}/tab/{tab}"
	ProcurementRequestSubmitURL        = "/action/procurement-request/submit/{id}"
	ProcurementRequestApproveURL       = "/action/procurement-request/approve/{id}"
	ProcurementRequestRejectURL        = "/action/procurement-request/reject/{id}"
	ProcurementRequestSpawnPOURL       = "/action/procurement-request/spawn-po/{id}"

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
	ProcurementDashboardURL        = "/app/procurement/dashboard"
	ProcurementRenewalCalendarURL  = "/app/procurement/renewals"
	ProcurementVarianceURL         = "/app/procurement/variance"
	ProcurementUtilizationURL      = "/app/procurement/utilization"
	ProcurementRecurrenceDraftsURL = "/app/procurement/recurrence-drafts/list/{status}"

	// ---------------------------------------------------------------------------
	// SPS P7 — SupplierContractPriceSchedule + SupplierContractPriceScheduleLine
	// ---------------------------------------------------------------------------

	// SupplierContractPriceSchedule master routes
	SupplierContractPriceScheduleListURL          = "/app/supplier-contract-price-schedules/list/{status}"
	SupplierContractPriceScheduleDetailURL        = "/app/supplier-contract-price-schedules/detail/{id}"
	SupplierContractPriceScheduleAddURL           = "/action/supplier-contract-price-schedule/add"
	SupplierContractPriceScheduleEditURL          = "/action/supplier-contract-price-schedule/edit/{id}"
	SupplierContractPriceScheduleDeleteURL        = "/action/supplier-contract-price-schedule/delete"
	SupplierContractPriceScheduleSetStatusURL     = "/action/supplier-contract-price-schedule/set-status"
	SupplierContractPriceScheduleBulkSetStatusURL = "/action/supplier-contract-price-schedule/bulk-set-status"
	SupplierContractPriceScheduleTabActionURL     = "/action/supplier-contract-price-schedule/detail/{id}/tab/{tab}"
	SupplierContractPriceScheduleActivateURL      = "/action/supplier-contract-price-schedule/activate/{id}"
	SupplierContractPriceScheduleSupersedeURL     = "/action/supplier-contract-price-schedule/supersede/{id}"

	// SupplierContractPriceScheduleLine routes (child of schedule detail)
	SupplierContractPriceScheduleLineAddURL    = "/action/supplier-contract-price-schedule/{id}/lines/add"
	SupplierContractPriceScheduleLineEditURL   = "/action/supplier-contract-price-schedule/{id}/lines/edit/{lid}"
	SupplierContractPriceScheduleLineDeleteURL = "/action/supplier-contract-price-schedule/{id}/lines/delete"

	// ---------------------------------------------------------------------------
	// SPS P10 — ExpenseRecognition + ExpenseRecognitionLine route constants
	// ---------------------------------------------------------------------------

	// ExpenseRecognition master routes (no add/edit drawer — created BY use case)
	ExpenseRecognitionListURL                    = "/app/expense-recognitions/list/{status}"
	ExpenseRecognitionDetailURL                  = "/app/expense-recognitions/detail/{id}"
	ExpenseRecognitionDeleteURL                  = "/action/expense-recognition/delete"
	ExpenseRecognitionTabActionURL               = "/action/expense-recognition/detail/{id}/tab/{tab}"
	ExpenseRecognitionReverseURL                 = "/action/expense-recognition/reverse/{id}"
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
	AccruedExpenseListURL          = "/app/accrued-expenses/list/{status}"
	AccruedExpenseDetailURL        = "/app/accrued-expenses/detail/{id}"
	AccruedExpenseAddURL           = "/action/accrued-expense/add"
	AccruedExpenseEditURL          = "/action/accrued-expense/edit/{id}"
	AccruedExpenseDeleteURL        = "/action/accrued-expense/delete"
	AccruedExpenseSetStatusURL     = "/action/accrued-expense/set-status"
	AccruedExpenseBulkSetStatusURL = "/action/accrued-expense/bulk-set-status"
	AccruedExpenseTabActionURL     = "/action/accrued-expense/detail/{id}/tab/{tab}"
	AccruedExpenseSettleURL            = "/action/accrued-expense/settle/{id}"
	AccruedExpenseReverseURL           = "/action/accrued-expense/reverse/{id}"
	AccruedExpenseAccrueFromContractURL = "/action/accrued-expense/accrue-from-contract"

	// AccruedExpenseSettlement routes (child of accrual detail — inline CRUD)
	AccruedExpenseSettlementAddURL    = "/action/accrued-expense/{id}/settlements/add"
	AccruedExpenseSettlementEditURL   = "/action/accrued-expense/{id}/settlements/edit/{sid}"
	AccruedExpenseSettlementDeleteURL = "/action/accrued-expense/{id}/settlements/delete"
)
