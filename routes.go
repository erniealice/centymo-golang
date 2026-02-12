package centymo

// Default route constants for centymo views.
// Consumer apps can use these or define their own.
const (
	PlanListURL              = "/app/plans/list/{status}"
	SubscriptionListURL      = "/app/subscriptions/list/{status}"
	ProductListURL           = "/app/products/list/{status}"
	ProductDetailURL         = "/app/products/detail/{id}"
	PaymentCollectionListURL = "/app/payment-collections/list/{status}"

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
	ProductTabActionURL = "/action/products/detail/{id}/tab/{tab}"

	// Product variant routes (within product detail)
	ProductVariantTableURL  = "/action/products/detail/{id}/variants/table"
	ProductVariantAssignURL = "/action/products/detail/{id}/variants/assign"
	ProductVariantEditURL   = "/action/products/detail/{id}/variants/edit/{vid}"
	ProductVariantRemoveURL = "/action/products/detail/{id}/variants/remove"

	// Product attribute routes (within product detail)
	ProductAttributeTableURL  = "/action/products/detail/{id}/attributes/table"
	ProductAttributeAssignURL = "/action/products/detail/{id}/attributes/assign"
	ProductAttributeRemoveURL = "/action/products/detail/{id}/attributes/remove"

	// Inventory status routes
	InventorySetStatusURL     = "/action/inventory/set-status"
	InventoryBulkSetStatusURL = "/action/inventory/bulk-set-status"

	// Inventory tab action route
	InventoryTabActionURL = "/action/inventory/detail/{id}/tab/{tab}"

	// Inventory movements (global)
	InventoryMovementsURL = "/app/inventory/movements"

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

	// Inventory table refresh (per-location)
	InventoryTableURL = "/action/inventory/table/{location}"

	// Sales (revenue) routes
	SalesDashboardURL     = "/app/sales/dashboard"
	SalesListURL          = "/app/sales/list/{status}"
	SalesDetailURL        = "/app/sales/detail/{id}"
	SalesAddURL           = "/action/sales/add"
	SalesEditURL          = "/action/sales/edit/{id}"
	SalesDeleteURL        = "/action/sales/delete"
	SalesBulkDeleteURL    = "/action/sales/delete/bulk"
	SalesSetStatusURL     = "/action/sales/status/set"
	SalesBulkSetStatusURL = "/action/sales/status/set/bulk"

	// Sales line item routes (within sale detail)
	SalesLineItemTableURL    = "/action/sales/detail/{id}/items/table"
	SalesLineItemAddURL      = "/action/sales/detail/{id}/items/add"
	SalesLineItemEditURL     = "/action/sales/detail/{id}/items/edit/{itemId}"
	SalesLineItemRemoveURL   = "/action/sales/detail/{id}/items/remove"
	SalesLineItemDiscountURL = "/action/sales/detail/{id}/items/add-discount"

	// Sales payment routes (within sale detail)
	SalesPaymentTableURL  = "/action/sales/detail/{id}/payment/table"
	SalesPaymentAddURL    = "/action/sales/detail/{id}/payment/add"
	SalesPaymentEditURL   = "/action/sales/detail/{id}/payment/edit/{pid}"
	SalesPaymentRemoveURL = "/action/sales/detail/{id}/payment/remove"

	// Price List routes
	PriceListListURL       = "/app/price-lists/list/{status}"
	PriceListDetailURL     = "/app/price-lists/{id}"
	PriceListAddURL        = "/action/price-lists/add"
	PriceListEditURL       = "/action/price-lists/edit/{id}"
	PriceListDeleteURL     = "/action/price-lists/delete"
	PriceListBulkDeleteURL = "/action/price-lists/bulk-delete"

	// Price Product routes (within price list detail)
	PriceProductAddURL    = "/action/price-lists/{id}/products/add"
	PriceProductDeleteURL = "/action/price-lists/{id}/products/delete"
)
