package centymo

// Default route constants for centymo views.
// Consumer apps can use these or define their own.
const (
	PlanListURL              = "/app/plans/list/{status}"
	SubscriptionListURL      = "/app/subscriptions/list/{status}"
	ProductListURL           = "/app/products/list/{status}"
	ProductDetailURL         = "/app/products/{id}"
	PaymentCollectionListURL = "/app/payment-collections/list/{status}"

	// Inventory routes — list is per-location
	InventoryListURL       = "/app/inventory/list/{location}"
	InventoryAddURL        = "/action/inventory/add"
	InventoryEditURL       = "/action/inventory/edit/{id}"
	InventoryDeleteURL     = "/action/inventory/delete"
	InventoryBulkDeleteURL = "/action/inventory/bulk-delete"

	// Product action routes
	ProductAddURL        = "/action/products/add"
	ProductEditURL       = "/action/products/edit/{id}"
	ProductDeleteURL     = "/action/products/delete"
	ProductBulkDeleteURL = "/action/products/bulk-delete"

	// Sales (revenue) routes
	SalesListURL       = "/app/sales/list/{status}"
	SalesDetailURL     = "/app/sales/{id}"
	SalesAddURL        = "/action/sales/add"
	SalesEditURL       = "/action/sales/edit/{id}"
	SalesDeleteURL     = "/action/sales/delete"
	SalesBulkDeleteURL = "/action/sales/bulk-delete"
)
