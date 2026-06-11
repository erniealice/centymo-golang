// Package inventory is the inventory-domain consumer facade (centymo restructure).
//
// PURE RE-EXPORT — zero behaviour. The inventory domain's data/route types,
// Default* constructors, and URL consts moved into per-entity packages under
// domain/inventory/<entity>/ with entity-local names (the <Entity> prefix stripped).
// This facade re-adds the original prefixed names so existing consumers
// (block/, service-admin) keep resolving inventory.<Entity>Labels /
// inventory.Default<Entity>Routes() / inventory.<Entity>ListURL unchanged.
//
// An entity package MUST NEVER import this facade (that would be an import
// cycle inventory -> <entity> -> inventory); cross-entity references go DIRECT to the
// sibling package.
package inventory

import (
	inventorypkg "github.com/erniealice/centymo-golang/domain/inventory/inventory"
)

// Re-exported data/route types (type aliases — identity-preserving).
type (
	InventoryActionLabels       = inventorypkg.ActionLabels
	InventoryBreadcrumbLabels   = inventorypkg.BreadcrumbLabels
	InventoryBulkLabels         = inventorypkg.BulkLabels
	InventoryButtonLabels       = inventorypkg.ButtonLabels
	InventoryColumnLabels       = inventorypkg.ColumnLabels
	InventoryConfirmLabels      = inventorypkg.ConfirmLabels
	InventoryDashboardLabels    = inventorypkg.DashboardLabels
	InventoryDepreciationLabels = inventorypkg.DepreciationLabels
	InventoryDetailLabels       = inventorypkg.DetailLabels
	InventoryEmptyLabels        = inventorypkg.EmptyLabels
	InventoryErrorLabels        = inventorypkg.ErrorLabels
	InventoryFormLabels         = inventorypkg.FormLabels
	InventoryLabels             = inventorypkg.Labels
	InventoryMovementsLabels    = inventorypkg.MovementsLabels
	InventoryPageLabels         = inventorypkg.PageLabels
	InventoryRoutes             = inventorypkg.Routes
	InventorySerialLabels       = inventorypkg.SerialLabels
	InventoryStatusLabels       = inventorypkg.StatusLabels
	InventoryTabLabels          = inventorypkg.TabLabels
	InventoryTransactionLabels  = inventorypkg.TransactionLabels
)

// Re-exported URL route consts (const-identity preserved).
const (
	InventoryAddURL                = inventorypkg.AddURL
	InventoryAttachmentDeleteURL   = inventorypkg.AttachmentDeleteURL
	InventoryAttachmentUploadURL   = inventorypkg.AttachmentUploadURL
	InventoryAttributeTableURL     = inventorypkg.AttributeTableURL
	InventoryBulkDeleteURL         = inventorypkg.BulkDeleteURL
	InventoryBulkSetStatusURL      = inventorypkg.BulkSetStatusURL
	InventoryDashboardAlertsURL    = inventorypkg.DashboardAlertsURL
	InventoryDashboardChartURL     = inventorypkg.DashboardChartURL
	InventoryDashboardMovementsURL = inventorypkg.DashboardMovementsURL
	InventoryDashboardStatsURL     = inventorypkg.DashboardStatsURL
	InventoryDashboardURL          = inventorypkg.DashboardURL
	InventoryDeleteURL             = inventorypkg.DeleteURL
	InventoryDepreciationAssignURL = inventorypkg.DepreciationAssignURL
	InventoryDepreciationEditURL   = inventorypkg.DepreciationEditURL
	InventoryDetailURL             = inventorypkg.DetailURL
	InventoryEditURL               = inventorypkg.EditURL
	InventoryListURL               = inventorypkg.ListURL
	InventoryMovementsExportURL    = inventorypkg.MovementsExportURL
	InventoryMovementsTableURL     = inventorypkg.MovementsTableURL
	InventoryMovementsURL          = inventorypkg.MovementsURL
	InventoryProductDetailURL      = inventorypkg.ProductDetailURL
	InventoryProductTabActionURL   = inventorypkg.ProductTabActionURL
	InventorySerialAssignURL       = inventorypkg.SerialAssignURL
	InventorySerialEditURL         = inventorypkg.SerialEditURL
	InventorySerialRemoveURL       = inventorypkg.SerialRemoveURL
	InventorySerialTableURL        = inventorypkg.SerialTableURL
	InventorySetStatusURL          = inventorypkg.SetStatusURL
	InventoryTabActionURL          = inventorypkg.TabActionURL
	InventoryTableURL              = inventorypkg.TableURL
	InventoryTransactionAssignURL  = inventorypkg.TransactionAssignURL
	InventoryTransactionTableURL   = inventorypkg.TransactionTableURL
)

// Re-exported Default* constructors (function values).
var (
	DefaultInventoryRoutes = inventorypkg.DefaultRoutes
)
