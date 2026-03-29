package centymo

// Default route constants for centymo views.
// Consumer apps can use these or define their own.
const (
	PlanListURL             = "/app/plans/list/{status}"
	PlanDetailURL           = "/app/plans/{id}"
	PlanAddURL              = "/action/plans/add"
	PlanEditURL             = "/action/plans/edit/{id}"
	PlanDeleteURL           = "/action/plans/delete"
	PlanTabActionURL        = "/action/plans/detail/{id}/tab/{tab}"
	PlanAttachmentUploadURL = "/action/plans/detail/{id}/attachments/upload"
	PlanAttachmentDeleteURL = "/action/plans/detail/{id}/attachments/delete"

	SubscriptionListURL             = "/app/subscriptions/list/{status}"
	SubscriptionDetailURL           = "/app/subscriptions/{id}"
	SubscriptionAddURL              = "/action/subscriptions/add"
	SubscriptionEditURL             = "/action/subscriptions/edit/{id}"
	SubscriptionDeleteURL           = "/action/subscriptions/delete"
	SubscriptionTabActionURL        = "/action/subscriptions/detail/{id}/tab/{tab}"
	SubscriptionAttachmentUploadURL = "/action/subscriptions/detail/{id}/attachments/upload"
	SubscriptionAttachmentDeleteURL = "/action/subscriptions/detail/{id}/attachments/delete"
	SubscriptionSearchPlanURL       = "/action/subscriptions/search/plans"
	SubscriptionSearchClientURL     = "/action/subscriptions/search/clients"

	// Collection (money IN) routes
	CollectionListURL             = "/app/collections/list/{status}"
	CollectionDetailURL           = "/app/collections/detail/{id}"
	CollectionDashboardURL        = "/app/collections/dashboard"
	CollectionAddURL              = "/action/collections/add"
	CollectionEditURL             = "/action/collections/edit/{id}"
	CollectionDeleteURL           = "/action/collections/delete"
	CollectionBulkDeleteURL       = "/action/collections/bulk-delete"
	CollectionSetStatusURL        = "/action/collections/set-status"
	CollectionBulkSetStatusURL    = "/action/collections/bulk-set-status"
	CollectionTabActionURL        = "/action/collections/detail/{id}/tab/{tab}"
	CollectionAttachmentUploadURL = "/action/collections/detail/{id}/attachments/upload"
	CollectionAttachmentDeleteURL = "/action/collections/detail/{id}/attachments/delete"

	// Disbursement (money OUT) routes
	DisbursementListURL             = "/app/disbursements/list/{status}"
	DisbursementDetailURL           = "/app/disbursements/detail/{id}"
	DisbursementDashboardURL        = "/app/disbursements/dashboard"
	DisbursementAddURL              = "/action/disbursements/add"
	DisbursementEditURL             = "/action/disbursements/edit/{id}"
	DisbursementDeleteURL           = "/action/disbursements/delete"
	DisbursementBulkDeleteURL       = "/action/disbursements/bulk-delete"
	DisbursementSetStatusURL        = "/action/disbursements/set-status"
	DisbursementBulkSetStatusURL    = "/action/disbursements/bulk-set-status"
	DisbursementTabActionURL        = "/action/disbursements/detail/{id}/tab/{tab}"
	DisbursementAttachmentUploadURL = "/action/disbursements/detail/{id}/attachments/upload"
	DisbursementAttachmentDeleteURL = "/action/disbursements/detail/{id}/attachments/delete"

	ProductListURL   = "/app/products/list/{status}"
	ProductTableURL  = "/action/products/table/{status}"
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
	ProductAddURL        = "/action/products/add"
	ProductEditURL       = "/action/products/edit/{id}"
	ProductDeleteURL     = "/action/products/delete"
	ProductBulkDeleteURL = "/action/products/bulk-delete"

	// Product status routes
	ProductSetStatusURL     = "/action/products/set-status"
	ProductBulkSetStatusURL = "/action/products/bulk-set-status"

	// Product detail tab action route
	ProductTabActionURL        = "/action/products/detail/{id}/tab/{tab}"
	ProductAttachmentUploadURL = "/action/products/detail/{id}/attachments/upload"
	ProductAttachmentDeleteURL = "/action/products/detail/{id}/attachments/delete"

	// Product variant routes (within product detail)
	ProductVariantTableURL  = "/action/products/detail/{id}/variants/table"
	ProductVariantAssignURL = "/action/products/detail/{id}/variants/assign"
	ProductVariantEditURL   = "/action/products/detail/{id}/variants/edit/{vid}"
	ProductVariantRemoveURL = "/action/products/detail/{id}/variants/remove"

	// Product attribute routes (within product detail)
	ProductAttributeTableURL  = "/action/products/detail/{id}/attributes/table"
	ProductAttributeAssignURL = "/action/products/detail/{id}/attributes/assign"
	ProductAttributeRemoveURL = "/action/products/detail/{id}/attributes/remove"

	// Product option routes (within product detail)
	ProductOptionTableURL  = "/action/products/detail/{id}/options/table"
	ProductOptionAddURL    = "/action/products/detail/{id}/options/add"
	ProductOptionEditURL   = "/action/products/detail/{id}/options/edit/{oid}"
	ProductOptionDeleteURL = "/action/products/detail/{id}/options/delete"

	// Product option detail page (option values management)
	ProductOptionDetailURL = "/app/products/detail/{id}/option/{oid}"

	// Product variant detail page (variant info, pricing, stock, audit, images)
	ProductVariantDetailURL    = "/app/products/detail/{id}/variant/{vid}"
	ProductVariantTabActionURL = "/action/products/detail/{id}/variant/{vid}/tab/{tab}"

	// Product variant image routes (upload/delete within variant detail)
	ProductVariantImageUploadURL = "/action/products/detail/{id}/variant/{vid}/images/upload"
	ProductVariantImageDeleteURL = "/action/products/detail/{id}/variant/{vid}/images/delete"

	// Product variant attachment routes
	ProductVariantAttachmentUploadURL = "/action/products/detail/{id}/variant/{vid}/attachments/upload"
	ProductVariantAttachmentDeleteURL = "/action/products/detail/{id}/variant/{vid}/attachments/delete"

	// Product variant stock detail (inventory item within variant context)
	ProductVariantStockDetailURL    = "/app/products/detail/{id}/variant/{vid}/stock/{iid}"
	ProductVariantStockTabActionURL = "/action/products/detail/{id}/variant/{vid}/stock/{iid}/tab/{tab}"

	// Product variant stock attachment routes
	ProductVariantStockAttachmentUploadURL = "/action/products/detail/{id}/variant/{vid}/stock/{iid}/attachments/upload"
	ProductVariantStockAttachmentDeleteURL = "/action/products/detail/{id}/variant/{vid}/stock/{iid}/attachments/delete"

	// Inventory serial detail (individual serial within inventory item)
	ProductVariantSerialDetailURL = "/app/products/detail/{id}/variant/{vid}/stock/{iid}/serial/{sid}"

	// Product option value routes (within product option)
	ProductOptionValueTableURL  = "/action/products/detail/{id}/options/{oid}/values/table"
	ProductOptionValueAddURL    = "/action/products/detail/{id}/options/{oid}/values/add"
	ProductOptionValueEditURL   = "/action/products/detail/{id}/options/{oid}/values/edit/{vid}"
	ProductOptionValueDeleteURL = "/action/products/detail/{id}/options/{oid}/values/delete"

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
	RevenueTableURL         = "/action/sales/table/{status}"
	RevenueDetailURL        = "/app/sales/detail/{id}"
	RevenueAddURL           = "/action/sales/add"
	RevenueEditURL          = "/action/sales/edit/{id}"
	RevenueDeleteURL        = "/action/sales/delete"
	RevenueBulkDeleteURL    = "/action/sales/delete/bulk"
	RevenueSetStatusURL     = "/action/sales/status/set"
	RevenueBulkSetStatusURL = "/action/sales/status/set/bulk"

	// Revenue tab action route
	RevenueTabActionURL        = "/action/sales/detail/{id}/tab/{tab}"
	RevenueAttachmentUploadURL = "/action/sales/detail/{id}/attachments/upload"
	RevenueAttachmentDeleteURL = "/action/sales/detail/{id}/attachments/delete"

	// Revenue line item routes (within sale detail)
	RevenueLineItemTableURL    = "/action/sales/detail/{id}/items/table"
	RevenueLineItemAddURL      = "/action/sales/detail/{id}/items/add"
	RevenueLineItemEditURL     = "/action/sales/detail/{id}/items/edit/{itemId}"
	RevenueLineItemRemoveURL   = "/action/sales/detail/{id}/items/remove"
	RevenueLineItemDiscountURL = "/action/sales/detail/{id}/items/add-discount"

	// Revenue payment routes (within sale detail)
	RevenuePaymentTableURL  = "/action/sales/detail/{id}/payment/table"
	RevenuePaymentAddURL    = "/action/sales/detail/{id}/payment/add"
	RevenuePaymentEditURL   = "/action/sales/detail/{id}/payment/edit/{pid}"
	RevenuePaymentRemoveURL = "/action/sales/detail/{id}/payment/remove"

	// Revenue report routes
	RevenueSummaryURL = "/app/sales/reports/sales-summary"

	// Revenue invoice document routes
	RevenueInvoiceDownloadURL = "/action/sales/detail/{id}/invoice/download"
	RevenueEmailURL           = "/action/sales/detail/{id}/invoice/send-email"

	// Revenue settings routes (template management)
	RevenueSettingsTemplatesURL       = "/app/sales/settings/templates"
	RevenueSettingsTemplateUploadURL  = "/action/sales/settings/templates/upload"
	RevenueSettingsTemplateDeleteURL  = "/action/sales/settings/templates/delete"
	RevenueSettingsTemplateDefaultURL = "/action/sales/settings/templates/set-default/{id}"

	// Expenditure (purchase + expense) routes
	ExpenditurePurchaseListURL      = "/app/purchases/list/{status}"
	ExpenditurePurchaseDashboardURL = "/app/purchases/dashboard"
	ExpenditureExpenseListURL       = "/app/expenses/list/{status}"
	ExpenditureExpenseDashboardURL  = "/app/expenses/dashboard"

	// Expenditure expense CRUD action routes
	ExpenditureExpenseAddURL       = "/action/expenses/add"
	ExpenditureExpenseEditURL      = "/action/expenses/edit/{id}"
	ExpenditureExpenseDeleteURL    = "/action/expenses/delete"
	ExpenditureExpenseSetStatusURL = "/action/expenses/set-status"
	ExpenditureExpenseDetailURL    = "/app/expenses/detail/{id}"
	ExpenditureExpenseTableURL     = "/action/expenses/table/{status}"
	ExpenditureExpenseTabActionURL = "/action/expenses/detail/{id}/tab/{tab}"

	// Expenditure expense line item action routes
	ExpenditureExpenseLineItemAddURL    = "/action/expenses/detail/{id}/items/add"
	ExpenditureExpenseLineItemEditURL   = "/action/expenses/detail/{id}/items/edit/{itemId}"
	ExpenditureExpenseLineItemRemoveURL = "/action/expenses/detail/{id}/items/remove"
	ExpenditureExpenseLineItemTableURL  = "/action/expenses/detail/{id}/items/table"

	// Expenditure report routes
	PurchasesSummaryURL = "/app/purchases/reports/purchases-summary"
	ExpensesSummaryURL  = "/app/expenses/reports/expenses-summary"

	// Expenditure settings (template management) routes
	ExpenditureSettingsTemplatesURL       = "/app/purchases/settings/templates"
	ExpenditureSettingsTemplateUploadURL  = "/action/purchases/settings/templates/upload"
	ExpenditureSettingsTemplateDeleteURL  = "/action/purchases/settings/templates/delete"
	ExpenditureSettingsTemplateDefaultURL = "/action/purchases/settings/templates/set-default/{id}"

	// Purchase Order routes
	PurchaseOrderListURL      = "/app/purchase-orders/list/{status}"
	PurchaseOrderDetailURL    = "/app/purchase-orders/{id}"
	PurchaseOrderAddURL       = "/action/purchase-orders/add"
	PurchaseOrderEditURL      = "/action/purchase-orders/edit/{id}"
	PurchaseOrderDeleteURL     = "/action/purchase-orders/delete"
	PurchaseOrderSetStatusURL  = "/action/purchase-orders/set-status"
	PurchaseOrderTableURL      = "/action/purchase-orders/table/{status}"
	PurchaseOrderTabActionURL  = "/action/purchase-orders/detail/{id}/tab/{tab}"

	// Expense category settings routes
	ExpenditureExpenseCategoryListURL   = "/app/expenses/categories/list"
	ExpenditureExpenseCategoryAddURL    = "/action/expenses/categories/add"
	ExpenditureExpenseCategoryEditURL   = "/action/expenses/categories/edit/{id}"
	ExpenditureExpenseCategoryDeleteURL = "/action/expenses/categories/delete"
	ExpenditureExpenseCategoryTableURL  = "/action/expenses/categories/table"

	// Price List routes
	PriceListListURL       = "/app/price-lists/list/{status}"
	PriceListTableURL      = "/action/price-lists/table/{status}"
	PriceListDetailURL     = "/app/price-lists/{id}"
	PriceListAddURL        = "/action/price-lists/add"
	PriceListEditURL       = "/action/price-lists/edit/{id}"
	PriceListDeleteURL     = "/action/price-lists/delete"
	PriceListBulkDeleteURL = "/action/price-lists/bulk-delete"

	PriceListTabActionURL        = "/action/price-lists/{id}/tab/{tab}"
	PriceListAttachmentUploadURL = "/action/price-lists/{id}/attachments/upload"
	PriceListAttachmentDeleteURL = "/action/price-lists/{id}/attachments/delete"

	// Price Product routes (within price list detail)
	PriceProductAddURL    = "/action/price-lists/{id}/products/add"
	PriceProductDeleteURL = "/action/price-lists/{id}/products/delete"

)
