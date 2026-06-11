package inventory

// Default route constants for inventory views.
// Consumer apps can use these or define their own via lyngua route.json overrides.
const (
	// Inventory routes — list is per-location
	InventoryDashboardURL  = "/inventory/dashboard"
	InventoryListURL       = "/inventory/list/{location}"
	InventoryAddURL        = "/action/inventory/add"
	InventoryEditURL       = "/action/inventory/edit/{id}"
	InventoryDeleteURL     = "/action/inventory/delete"
	InventoryBulkDeleteURL = "/action/inventory/bulk-delete"
	InventoryDetailURL     = "/inventory/detail/{id}"

	// Inventory status routes
	InventorySetStatusURL     = "/action/inventory/set-status"
	InventoryBulkSetStatusURL = "/action/inventory/bulk-set-status"

	// Inventory tab action route
	InventoryTabActionURL        = "/action/inventory/detail/{id}/tab/{tab}"
	InventoryAttachmentUploadURL = "/action/inventory/detail/{id}/attachments/upload"
	InventoryAttachmentDeleteURL = "/action/inventory/detail/{id}/attachments/delete"

	// Inventory movements (global transaction history)
	InventoryMovementsURL       = "/inventory/movements"
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
	InventoryProductDetailURL    = "/products/detail/{pid}/inventory/detail/{iid}"
	InventoryProductTabActionURL = "/action/product/{pid}/inventory/{iid}/tab/{tab}"

	// Inventory table refresh (per-location)
	InventoryTableURL = "/action/inventory/table/{location}"
)

// InventoryRoutes holds all route paths for inventory views and actions,
// including serial, transaction, depreciation, dashboard, and movement sub-routes.
type InventoryRoutes struct {
	DashboardURL  string `json:"dashboard_url"`
	ListURL       string `json:"list_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`
	DetailURL     string `json:"detail_url"`
	TableURL      string `json:"table_url"`

	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`

	TabActionURL string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Movement routes
	MovementsURL       string `json:"movements_url"`
	MovementsTableURL  string `json:"movements_table_url"`
	MovementsExportURL string `json:"movements_export_url"`

	// Serial routes
	SerialTableURL  string `json:"serial_table_url"`
	SerialAssignURL string `json:"serial_assign_url"`
	SerialEditURL   string `json:"serial_edit_url"`
	SerialRemoveURL string `json:"serial_remove_url"`

	// Transaction routes
	TransactionTableURL  string `json:"transaction_table_url"`
	TransactionAssignURL string `json:"transaction_assign_url"`

	// Depreciation routes
	DepreciationAssignURL string `json:"depreciation_assign_url"`
	DepreciationEditURL   string `json:"depreciation_edit_url"`

	// Attribute routes
	AttributeTableURL string `json:"attribute_table_url"`

	// Dashboard partial routes
	DashboardStatsURL     string `json:"dashboard_stats_url"`
	DashboardChartURL     string `json:"dashboard_chart_url"`
	DashboardMovementsURL string `json:"dashboard_movements_url"`
	DashboardAlertsURL    string `json:"dashboard_alerts_url"`

	// Product-context detail routes
	ProductDetailURL    string `json:"product_detail_url"`
	ProductTabActionURL string `json:"product_tab_action_url"`
}

// DefaultInventoryRoutes returns an InventoryRoutes populated from the
// package-level route constants.
func DefaultInventoryRoutes() InventoryRoutes {
	return InventoryRoutes{
		DashboardURL:  InventoryDashboardURL,
		ListURL:       InventoryListURL,
		AddURL:        InventoryAddURL,
		EditURL:       InventoryEditURL,
		DeleteURL:     InventoryDeleteURL,
		BulkDeleteURL: InventoryBulkDeleteURL,
		DetailURL:     InventoryDetailURL,
		TableURL:      InventoryTableURL,

		SetStatusURL:     InventorySetStatusURL,
		BulkSetStatusURL: InventoryBulkSetStatusURL,

		TabActionURL: InventoryTabActionURL,

		AttachmentUploadURL: InventoryAttachmentUploadURL,
		AttachmentDeleteURL: InventoryAttachmentDeleteURL,

		MovementsURL:       InventoryMovementsURL,
		MovementsTableURL:  InventoryMovementsTableURL,
		MovementsExportURL: InventoryMovementsExportURL,

		SerialTableURL:  InventorySerialTableURL,
		SerialAssignURL: InventorySerialAssignURL,
		SerialEditURL:   InventorySerialEditURL,
		SerialRemoveURL: InventorySerialRemoveURL,

		TransactionTableURL:  InventoryTransactionTableURL,
		TransactionAssignURL: InventoryTransactionAssignURL,

		DepreciationAssignURL: InventoryDepreciationAssignURL,
		DepreciationEditURL:   InventoryDepreciationEditURL,

		AttributeTableURL: InventoryAttributeTableURL,

		DashboardStatsURL:     InventoryDashboardStatsURL,
		DashboardChartURL:     InventoryDashboardChartURL,
		DashboardMovementsURL: InventoryDashboardMovementsURL,
		DashboardAlertsURL:    InventoryDashboardAlertsURL,

		ProductDetailURL:    InventoryProductDetailURL,
		ProductTabActionURL: InventoryProductTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// inventory routes.
func (r InventoryRoutes) RouteMap() map[string]string {
	return map[string]string{
		"inventory.dashboard":   r.DashboardURL,
		"inventory.list":        r.ListURL,
		"inventory.add":         r.AddURL,
		"inventory.edit":        r.EditURL,
		"inventory.delete":      r.DeleteURL,
		"inventory.bulk_delete": r.BulkDeleteURL,
		"inventory.detail":      r.DetailURL,
		"inventory.table":       r.TableURL,

		"inventory.set_status":      r.SetStatusURL,
		"inventory.bulk_set_status": r.BulkSetStatusURL,

		"inventory.tab_action": r.TabActionURL,

		"inventory.attachment.upload": r.AttachmentUploadURL,
		"inventory.attachment.delete": r.AttachmentDeleteURL,

		"inventory.movements":        r.MovementsURL,
		"inventory.movements.table":  r.MovementsTableURL,
		"inventory.movements.export": r.MovementsExportURL,

		"inventory.serial.table":  r.SerialTableURL,
		"inventory.serial.assign": r.SerialAssignURL,
		"inventory.serial.edit":   r.SerialEditURL,
		"inventory.serial.remove": r.SerialRemoveURL,

		"inventory.transaction.table":  r.TransactionTableURL,
		"inventory.transaction.assign": r.TransactionAssignURL,

		"inventory.depreciation.assign": r.DepreciationAssignURL,
		"inventory.depreciation.edit":   r.DepreciationEditURL,

		"inventory.attribute.table": r.AttributeTableURL,

		"inventory.dashboard.stats":     r.DashboardStatsURL,
		"inventory.dashboard.chart":     r.DashboardChartURL,
		"inventory.dashboard.movements": r.DashboardMovementsURL,
		"inventory.dashboard.alerts":    r.DashboardAlertsURL,

		"inventory.product.detail":     r.ProductDetailURL,
		"inventory.product.tab_action": r.ProductTabActionURL,
	}
}
