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
	PricePlanAddURL    = "/action/plan/{id}/pricelists/add"
	PricePlanEditURL   = "/action/plan/{id}/pricelists/edit/{ppid}"
	PricePlanDeleteURL = "/action/plan/{id}/pricelists/delete"

	// ProductPlan CRUD routes (within plan context)
	PlanProductPlanAddURL    = "/action/plan/{id}/products/add"
	PlanProductPlanEditURL   = "/action/plan/{id}/products/edit/{ppid}"
	PlanProductPlanDeleteURL = "/action/plan/{id}/products/delete"

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

	SubscriptionListURL             = "/app/subscriptions/list/{status}"
	SubscriptionDetailURL           = "/app/subscriptions/detail/{id}"
	SubscriptionAddURL              = "/action/subscription/add"
	SubscriptionEditURL             = "/action/subscription/edit/{id}"
	SubscriptionDeleteURL           = "/action/subscription/delete"
	SubscriptionTabActionURL        = "/action/subscription/detail/{id}/tab/{tab}"
	SubscriptionAttachmentUploadURL = "/action/subscription/detail/{id}/attachments/upload"
	SubscriptionAttachmentDeleteURL = "/action/subscription/detail/{id}/attachments/delete"
	SubscriptionSearchPlanURL       = "/action/subscription/search/plans"
	SubscriptionSearchClientURL     = "/action/subscription/search/clients"

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

	// Price List routes
	PriceListListURL       = "/app/price-lists/list/{status}"
	PriceListTableURL      = "/action/price-list/table/{status}"
	PriceListDetailURL     = "/app/price-lists/detail/{id}"
	PriceListAddURL        = "/action/price-list/add"
	PriceListEditURL       = "/action/price-list/edit/{id}"
	PriceListDeleteURL     = "/action/price-list/delete"
	PriceListBulkDeleteURL = "/action/price-list/bulk-delete"

	PriceListTabActionURL        = "/action/price-list/{id}/tab/{tab}"
	PriceListAttachmentUploadURL = "/action/price-list/{id}/attachments/upload"
	PriceListAttachmentDeleteURL = "/action/price-list/{id}/attachments/delete"

	// Price Product routes (within price list detail)
	PriceProductAddURL    = "/action/price-list/{id}/products/add"
	PriceProductDeleteURL = "/action/price-list/{id}/products/delete"
)
