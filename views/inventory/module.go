package inventory

import (
	"context"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang"
	inventoryaction "github.com/erniealice/centymo-golang/views/inventory/action"
	inventorydashboard "github.com/erniealice/centymo-golang/views/inventory/dashboard"
	inventorydepreciation "github.com/erniealice/centymo-golang/views/inventory/depreciation"
	inventorydetail "github.com/erniealice/centymo-golang/views/inventory/detail"
	inventorylist "github.com/erniealice/centymo-golang/views/inventory/list"
	inventorymovements "github.com/erniealice/centymo-golang/views/inventory/movements"
	inventoryserial "github.com/erniealice/centymo-golang/views/inventory/serial"
	inventorytransaction "github.com/erniealice/centymo-golang/views/inventory/transaction"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	inventorydepreciationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_depreciation"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	inventorytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_transaction"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
)

// ModuleDeps holds all dependencies for the inventory module.
type ModuleDeps struct {
	Routes       centymo.InventoryRoutes
	Labels       centymo.InventoryLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Inventory Item CRUD
	ListInventoryItems  func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	CreateInventoryItem func(ctx context.Context, req *inventoryitempb.CreateInventoryItemRequest) (*inventoryitempb.CreateInventoryItemResponse, error)
	ReadInventoryItem   func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	UpdateInventoryItem func(ctx context.Context, req *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error)
	DeleteInventoryItem func(ctx context.Context, req *inventoryitempb.DeleteInventoryItemRequest) (*inventoryitempb.DeleteInventoryItemResponse, error)
	SetItemActive       func(ctx context.Context, id string, active bool) error

	// Inventory Serial CRUD
	ListInventorySerials  func(ctx context.Context, req *inventoryserialpb.ListInventorySerialsRequest) (*inventoryserialpb.ListInventorySerialsResponse, error)
	CreateInventorySerial func(ctx context.Context, req *inventoryserialpb.CreateInventorySerialRequest) (*inventoryserialpb.CreateInventorySerialResponse, error)
	ReadInventorySerial   func(ctx context.Context, req *inventoryserialpb.ReadInventorySerialRequest) (*inventoryserialpb.ReadInventorySerialResponse, error)
	UpdateInventorySerial func(ctx context.Context, req *inventoryserialpb.UpdateInventorySerialRequest) (*inventoryserialpb.UpdateInventorySerialResponse, error)
	DeleteInventorySerial func(ctx context.Context, req *inventoryserialpb.DeleteInventorySerialRequest) (*inventoryserialpb.DeleteInventorySerialResponse, error)

	// Inventory Transaction
	ListInventoryTransactions         func(ctx context.Context, req *inventorytransactionpb.ListInventoryTransactionsRequest) (*inventorytransactionpb.ListInventoryTransactionsResponse, error)
	CreateInventoryTransaction        func(ctx context.Context, req *inventorytransactionpb.CreateInventoryTransactionRequest) (*inventorytransactionpb.CreateInventoryTransactionResponse, error)
	GetInventoryMovementsListPageData func(ctx context.Context, req *inventorytransactionpb.GetInventoryMovementsListPageDataRequest) (*inventorytransactionpb.GetInventoryMovementsListPageDataResponse, error)

	// Inventory Depreciation
	ListInventoryDepreciations  func(ctx context.Context, req *inventorydepreciationpb.ListInventoryDepreciationsRequest) (*inventorydepreciationpb.ListInventoryDepreciationsResponse, error)
	CreateInventoryDepreciation func(ctx context.Context, req *inventorydepreciationpb.CreateInventoryDepreciationRequest) (*inventorydepreciationpb.CreateInventoryDepreciationResponse, error)
	ReadInventoryDepreciation   func(ctx context.Context, req *inventorydepreciationpb.ReadInventoryDepreciationRequest) (*inventorydepreciationpb.ReadInventoryDepreciationResponse, error)
	UpdateInventoryDepreciation func(ctx context.Context, req *inventorydepreciationpb.UpdateInventoryDepreciationRequest) (*inventorydepreciationpb.UpdateInventoryDepreciationResponse, error)

	// Cross-domain deps (product, location)
	ReadProduct               func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	ListProductVariantOptions func(ctx context.Context, req *productvariantoptionpb.ListProductVariantOptionsRequest) (*productvariantoptionpb.ListProductVariantOptionsResponse, error)
	ListProductOptionValues   func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	ListProductOptions        func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ListLocations             func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string

	// Audit history
	ListAuditHistory func(ctx context.Context, req *auditlog.ListAuditRequest) (*auditlog.ListAuditResponse, error)
}

// Module holds all constructed inventory views.
type Module struct {
	routes             centymo.InventoryRoutes
	Dashboard          view.View
	List               view.View
	Detail             view.View
	Add                view.View
	Edit               view.View
	Delete             view.View
	BulkDelete         view.View
	SetStatus          view.View
	BulkSetStatus      view.View
	TabAction          view.View
	SerialTable        view.View
	SerialAssign       view.View
	SerialEdit         view.View
	SerialRemove       view.View
	TransactionTable   view.View
	TransactionAssign  view.View
	DepreciationAssign view.View
	DepreciationEdit   view.View
	DashboardStats     view.View
	DashboardChart     view.View
	DashboardMovements view.View
	DashboardAlerts    view.View
	Movements          view.View
	MovementsTable     view.View
	MovementsExport    http.HandlerFunc
	ProductDetail      view.View
	ProductTabAction   view.View
	AttachmentUpload   view.View
	AttachmentDelete   view.View
}

func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &inventoryaction.Deps{
		Routes:              deps.Routes,
		Labels:              deps.Labels,
		CreateInventoryItem: deps.CreateInventoryItem,
		ReadInventoryItem:   deps.ReadInventoryItem,
		UpdateInventoryItem: deps.UpdateInventoryItem,
		DeleteInventoryItem: deps.DeleteInventoryItem,
	}

	dashboardDeps := &inventorydashboard.Deps{
		Routes:                     deps.Routes,
		ListInventoryItems:         deps.ListInventoryItems,
		ListInventorySerials:       deps.ListInventorySerials,
		ListInventoryTransactions:  deps.ListInventoryTransactions,
		ListInventoryDepreciations: deps.ListInventoryDepreciations,
		Labels:                     deps.Labels,
		CommonLabels:               deps.CommonLabels,
	}

	detailDeps := &inventorydetail.DetailViewDeps{
		ReadInventoryItem:          deps.ReadInventoryItem,
		ListInventorySerials:       deps.ListInventorySerials,
		ListInventoryTransactions:  deps.ListInventoryTransactions,
		ListInventoryDepreciations: deps.ListInventoryDepreciations,
		ListProductVariantOptions:  deps.ListProductVariantOptions,
		ListProductOptionValues:    deps.ListProductOptionValues,
		ListProductOptions:         deps.ListProductOptions,
		Labels:                     deps.Labels,
		CommonLabels:               deps.CommonLabels,
		TableLabels:                deps.TableLabels,
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
		AuditOps: auditlog.AuditOps{
			ListAuditHistory: deps.ListAuditHistory,
		},
	}

	productDetailDeps := &inventorydetail.ProductDetailDeps{
		ReadInventoryItem: deps.ReadInventoryItem,
		ReadProduct:       deps.ReadProduct,
		DetailDeps:        detailDeps,
		Labels:            deps.Labels,
		CommonLabels:      deps.CommonLabels,
		TableLabels:       deps.TableLabels,
	}

	movementsDeps := &inventorymovements.Deps{
		GetInventoryMovementsListPageData: deps.GetInventoryMovementsListPageData,
		ListInventoryItems:                deps.ListInventoryItems,
		ListInventoryTransactions:         deps.ListInventoryTransactions,
		ListLocations:                     deps.ListLocations,
		Labels:                            deps.Labels,
		CommonLabels:                      deps.CommonLabels,
		TableLabels:                       deps.TableLabels,
	}

	depreciationDeps := &inventorydepreciation.Deps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CreateInventoryDepreciation: deps.CreateInventoryDepreciation,
		ReadInventoryDepreciation:   deps.ReadInventoryDepreciation,
		UpdateInventoryDepreciation: deps.UpdateInventoryDepreciation,
	}

	serialDeps := &inventoryserial.Deps{
		Routes:                deps.Routes,
		Labels:                deps.Labels,
		CreateInventorySerial: deps.CreateInventorySerial,
		ReadInventorySerial:   deps.ReadInventorySerial,
		UpdateInventorySerial: deps.UpdateInventorySerial,
		DeleteInventorySerial: deps.DeleteInventorySerial,
	}

	transactionDeps := &inventorytransaction.Deps{
		Routes:                     deps.Routes,
		Labels:                     deps.Labels,
		CreateInventoryTransaction: deps.CreateInventoryTransaction,
		ReadInventoryItem:          deps.ReadInventoryItem,
		UpdateInventoryItem:        deps.UpdateInventoryItem,
	}

	return &Module{
		routes:    deps.Routes,
		Dashboard: inventorydashboard.NewView(dashboardDeps),
		List: inventorylist.NewView(&inventorylist.ListViewDeps{
			Routes:             deps.Routes,
			ListInventoryItems: deps.ListInventoryItems,
			Labels:             deps.Labels,
			CommonLabels:       deps.CommonLabels,
			TableLabels:        deps.TableLabels,
		}),
		Detail:             inventorydetail.NewView(detailDeps),
		Add:                inventoryaction.NewAddAction(actionDeps),
		Edit:               inventoryaction.NewEditAction(actionDeps),
		Delete:             inventoryaction.NewDeleteAction(actionDeps),
		BulkDelete:         inventoryaction.NewBulkDeleteAction(actionDeps),
		SetStatus:          inventoryaction.NewSetStatusAction(deps.SetItemActive, deps.Labels),
		BulkSetStatus:      inventoryaction.NewBulkSetStatusAction(deps.SetItemActive, deps.Labels),
		TabAction:          inventorydetail.NewTabAction(detailDeps),
		SerialTable:        inventoryserial.NewTableAction(serialDeps),
		SerialAssign:       inventoryserial.NewAssignAction(serialDeps),
		SerialEdit:         inventoryserial.NewEditAction(serialDeps),
		SerialRemove:       inventoryserial.NewRemoveAction(serialDeps),
		TransactionTable:   inventorytransaction.NewTableAction(transactionDeps),
		TransactionAssign:  inventorytransaction.NewAssignAction(transactionDeps),
		DepreciationAssign: inventorydepreciation.NewAssignAction(depreciationDeps),
		DepreciationEdit:   inventorydepreciation.NewEditAction(depreciationDeps),
		DashboardStats:     inventorydashboard.NewDashboardStatsAction(dashboardDeps),
		DashboardChart:     inventorydashboard.NewDashboardChartAction(dashboardDeps),
		DashboardMovements: inventorydashboard.NewDashboardMovementsAction(dashboardDeps),
		DashboardAlerts:    inventorydashboard.NewDashboardAlertsAction(dashboardDeps),
		Movements:          inventorymovements.NewView(movementsDeps),
		MovementsTable:     inventorymovements.NewFilterView(movementsDeps),
		MovementsExport:    inventorymovements.NewExportHandler(movementsDeps),
		ProductDetail:      inventorydetail.NewProductDetailView(productDetailDeps),
		ProductTabAction:   inventorydetail.NewProductDetailTabAction(productDetailDeps),
		AttachmentUpload:   inventorydetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete:   inventorydetail.NewAttachmentDeleteAction(detailDeps),
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	r.GET(m.routes.TabActionURL, m.TabAction)
	r.GET(m.routes.SerialTableURL, m.SerialTable)
	r.GET(m.routes.SerialAssignURL, m.SerialAssign)
	r.POST(m.routes.SerialAssignURL, m.SerialAssign)
	r.GET(m.routes.SerialEditURL, m.SerialEdit)
	r.POST(m.routes.SerialEditURL, m.SerialEdit)
	r.POST(m.routes.SerialRemoveURL, m.SerialRemove)
	r.GET(m.routes.TransactionTableURL, m.TransactionTable)
	r.GET(m.routes.TransactionAssignURL, m.TransactionAssign)
	r.POST(m.routes.TransactionAssignURL, m.TransactionAssign)
	r.GET(m.routes.DepreciationAssignURL, m.DepreciationAssign)
	r.POST(m.routes.DepreciationAssignURL, m.DepreciationAssign)
	r.GET(m.routes.DepreciationEditURL, m.DepreciationEdit)
	r.POST(m.routes.DepreciationEditURL, m.DepreciationEdit)
	r.GET(m.routes.DashboardStatsURL, m.DashboardStats)
	r.GET(m.routes.DashboardChartURL, m.DashboardChart)
	r.GET(m.routes.DashboardMovementsURL, m.DashboardMovements)
	r.GET(m.routes.DashboardAlertsURL, m.DashboardAlerts)
	r.GET(m.routes.MovementsURL, m.Movements)
	r.GET(m.routes.MovementsTableURL, m.MovementsTable)

	// Product-context inventory detail
	r.GET(m.routes.ProductDetailURL, m.ProductDetail)
	r.GET(m.routes.ProductTabActionURL, m.ProductTabAction)
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}
