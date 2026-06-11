package inventory

// Default route constants for inventory views.
// Consumer apps can use these or define their own via lyngua route.json overrides.
const (
	// Inventory routes — list is per-location
	DashboardURL  = "/inventory/dashboard"
	ListURL       = "/inventory/list/{location}"
	AddURL        = "/action/inventory/add"
	EditURL       = "/action/inventory/edit/{id}"
	DeleteURL     = "/action/inventory/delete"
	BulkDeleteURL = "/action/inventory/bulk-delete"
	DetailURL     = "/inventory/detail/{id}"

	// Inventory status routes
	SetStatusURL     = "/action/inventory/set-status"
	BulkSetStatusURL = "/action/inventory/bulk-set-status"

	// Inventory tab action route
	TabActionURL        = "/action/inventory/detail/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/inventory/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/inventory/detail/{id}/attachments/delete"

	// Inventory movements (global transaction history)
	MovementsURL       = "/inventory/movements"
	MovementsTableURL  = "/action/inventory/movements/table"
	MovementsExportURL = "/action/inventory/movements/export"

	// Inventory serial routes
	SerialTableURL  = "/action/inventory/detail/{id}/serials/table"
	SerialAssignURL = "/action/inventory/detail/{id}/serials/assign"
	SerialEditURL   = "/action/inventory/detail/{id}/serials/edit/{sid}"
	SerialRemoveURL = "/action/inventory/detail/{id}/serials/remove"

	// Inventory transaction routes
	TransactionTableURL  = "/action/inventory/detail/{id}/transactions/table"
	TransactionAssignURL = "/action/inventory/detail/{id}/transactions/assign"

	// Inventory depreciation routes
	DepreciationAssignURL = "/action/inventory/detail/{id}/depreciation/assign"
	DepreciationEditURL   = "/action/inventory/detail/{id}/depreciation/edit/{did}"

	// Inventory attribute routes (within detail)
	AttributeTableURL = "/action/inventory/detail/{id}/attributes/table"

	// Inventory dashboard partial routes
	DashboardStatsURL     = "/action/inventory/dashboard/stats"
	DashboardChartURL     = "/action/inventory/dashboard/chart"
	DashboardMovementsURL = "/action/inventory/dashboard/movements"
	DashboardAlertsURL    = "/action/inventory/dashboard/alerts"

	// Inventory product-context detail routes
	ProductDetailURL    = "/products/detail/{pid}/inventory/detail/{iid}"
	ProductTabActionURL = "/action/product/{pid}/inventory/{iid}/tab/{tab}"

	// Inventory table refresh (per-location)
	TableURL = "/action/inventory/table/{location}"
)

// Routes holds all route paths for inventory views and actions,
// including serial, transaction, depreciation, dashboard, and movement sub-routes.
type Routes struct {
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

// DefaultRoutes returns an Routes populated from the
// package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		DashboardURL:  DashboardURL,
		ListURL:       ListURL,
		AddURL:        AddURL,
		EditURL:       EditURL,
		DeleteURL:     DeleteURL,
		BulkDeleteURL: BulkDeleteURL,
		DetailURL:     DetailURL,
		TableURL:      TableURL,

		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,

		TabActionURL: TabActionURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		MovementsURL:       MovementsURL,
		MovementsTableURL:  MovementsTableURL,
		MovementsExportURL: MovementsExportURL,

		SerialTableURL:  SerialTableURL,
		SerialAssignURL: SerialAssignURL,
		SerialEditURL:   SerialEditURL,
		SerialRemoveURL: SerialRemoveURL,

		TransactionTableURL:  TransactionTableURL,
		TransactionAssignURL: TransactionAssignURL,

		DepreciationAssignURL: DepreciationAssignURL,
		DepreciationEditURL:   DepreciationEditURL,

		AttributeTableURL: AttributeTableURL,

		DashboardStatsURL:     DashboardStatsURL,
		DashboardChartURL:     DashboardChartURL,
		DashboardMovementsURL: DashboardMovementsURL,
		DashboardAlertsURL:    DashboardAlertsURL,

		ProductDetailURL:    ProductDetailURL,
		ProductTabActionURL: ProductTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// inventory routes.
func (r Routes) RouteMap() map[string]string {
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
