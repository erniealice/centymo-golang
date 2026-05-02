package purchaseorder

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	purchaseorderlineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order_line_item"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventorymovementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/inventory_movement"

	purchaseorderaction "github.com/erniealice/centymo-golang/views/purchase_order/action"
	purchaseorderdetail "github.com/erniealice/centymo-golang/views/purchase_order/detail"
	purchaseorderlineitem "github.com/erniealice/centymo-golang/views/purchase_order/line_item"
	purchaseorderlist "github.com/erniealice/centymo-golang/views/purchase_order/list"
	purchaseorderreceipt "github.com/erniealice/centymo-golang/views/purchase_order/receipt"
)

// ModuleDeps holds all dependencies for the purchase order module.
type ModuleDeps struct {
	Routes       centymo.ExpenditureRoutes
	DB           centymo.DataSource
	Labels       centymo.ExpenditureLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Purchase order CRUD operations (nil-guarded — only built when provided)
	ListPurchaseOrders  func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)
	CreatePurchaseOrder func(ctx context.Context, req *purchaseorderpb.CreatePurchaseOrderRequest) (*purchaseorderpb.CreatePurchaseOrderResponse, error)
	ReadPurchaseOrder   func(ctx context.Context, req *purchaseorderpb.ReadPurchaseOrderRequest) (*purchaseorderpb.ReadPurchaseOrderResponse, error)
	UpdatePurchaseOrder func(ctx context.Context, req *purchaseorderpb.UpdatePurchaseOrderRequest) (*purchaseorderpb.UpdatePurchaseOrderResponse, error)
	DeletePurchaseOrder func(ctx context.Context, req *purchaseorderpb.DeletePurchaseOrderRequest) (*purchaseorderpb.DeletePurchaseOrderResponse, error)

	// Purchase order line item operations (optional — used by detail and line item action views)
	ListPurchaseOrderLineItems  func(ctx context.Context, req *purchaseorderlineitempb.ListPurchaseOrderLineItemsRequest) (*purchaseorderlineitempb.ListPurchaseOrderLineItemsResponse, error)
	CreatePurchaseOrderLineItem func(ctx context.Context, req *purchaseorderlineitempb.CreatePurchaseOrderLineItemRequest) (*purchaseorderlineitempb.CreatePurchaseOrderLineItemResponse, error)
	ReadPurchaseOrderLineItem   func(ctx context.Context, req *purchaseorderlineitempb.ReadPurchaseOrderLineItemRequest) (*purchaseorderlineitempb.ReadPurchaseOrderLineItemResponse, error)
	UpdatePurchaseOrderLineItem func(ctx context.Context, req *purchaseorderlineitempb.UpdatePurchaseOrderLineItemRequest) (*purchaseorderlineitempb.UpdatePurchaseOrderLineItemResponse, error)
	DeletePurchaseOrderLineItem func(ctx context.Context, req *purchaseorderlineitempb.DeletePurchaseOrderLineItemRequest) (*purchaseorderlineitempb.DeletePurchaseOrderLineItemResponse, error)

	// Inventory operations (optional — used by confirm-receipt for goods lines)
	CreateInventoryMovement func(ctx context.Context, req *inventorymovementpb.CreateInventoryMovementRequest) (*inventorymovementpb.CreateInventoryMovementResponse, error)
	ReadInventoryItem       func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	UpdateInventoryItem     func(ctx context.Context, req *inventoryitempb.UpdateInventoryItemRequest) (*inventoryitempb.UpdateInventoryItemResponse, error)
}

// Module holds all constructed purchase order views.
type Module struct {
	routes                       centymo.ExpenditureRoutes
	PurchaseOrderList            view.View
	PurchaseOrderAdd             view.View
	PurchaseOrderEdit            view.View
	PurchaseOrderDelete          view.View
	PurchaseOrderSetStatus       view.View
	PurchaseOrderDetail          view.View
	PurchaseOrderTabAction       view.View
	PurchaseOrderLineItemTable   view.View
	PurchaseOrderLineItemAdd     view.View
	PurchaseOrderLineItemEdit    view.View
	PurchaseOrderLineItemRemove  view.View
	PurchaseOrderConfirmReceipt  view.View
}

// NewModule creates the purchase order module views.
func NewModule(deps *ModuleDeps) *Module {
	m := &Module{
		routes: deps.Routes,
	}

	// List view (nil-guarded — only built when ListPurchaseOrders is provided)
	if deps.ListPurchaseOrders != nil {
		m.PurchaseOrderList = purchaseorderlist.NewView(&purchaseorderlist.ListViewDeps{
			ListPurchaseOrders: deps.ListPurchaseOrders,
			RefreshURL:         deps.Routes.PurchaseOrderTableURL,
			AddURL:             deps.Routes.PurchaseOrderAddURL,
			Labels:             deps.Labels,
			CommonLabels:       deps.CommonLabels,
			TableLabels:        deps.TableLabels,
		})
	}

	// CRUD action views (nil-guarded — only built when CreatePurchaseOrder is provided)
	if deps.CreatePurchaseOrder != nil {
		actionDeps := &purchaseorderaction.Deps{
			Routes:              deps.Routes,
			Labels:              deps.Labels,
			CreatePurchaseOrder: deps.CreatePurchaseOrder,
			ReadPurchaseOrder:   deps.ReadPurchaseOrder,
			UpdatePurchaseOrder: deps.UpdatePurchaseOrder,
			DeletePurchaseOrder: deps.DeletePurchaseOrder,
		}
		m.PurchaseOrderAdd = purchaseorderaction.NewAddAction(actionDeps)
		m.PurchaseOrderEdit = purchaseorderaction.NewEditAction(actionDeps)
		m.PurchaseOrderDelete = purchaseorderaction.NewDeleteAction(actionDeps)
		m.PurchaseOrderSetStatus = purchaseorderaction.NewSetStatusAction(actionDeps)
	}

	// Detail view (nil-guarded — only built when ReadPurchaseOrder is provided)
	if deps.ReadPurchaseOrder != nil {
		detailDeps := &purchaseorderdetail.DetailViewDeps{
			Routes:                     deps.Routes,
			Labels:                     deps.Labels,
			CommonLabels:               deps.CommonLabels,
			TableLabels:                deps.TableLabels,
			ReadPurchaseOrder:          deps.ReadPurchaseOrder,
			ListPurchaseOrderLineItems: deps.ListPurchaseOrderLineItems,
		}
		m.PurchaseOrderDetail = purchaseorderdetail.NewView(detailDeps)
		m.PurchaseOrderTabAction = purchaseorderdetail.NewTabAction(detailDeps)
	}

	// Line item table view (nil-guarded — only built when ListPurchaseOrderLineItems is provided)
	if deps.ListPurchaseOrderLineItems != nil {
		lineItemTableDeps := &purchaseorderdetail.LineItemDeps{
			Routes:                     deps.Routes,
			Labels:                     deps.Labels,
			TableLabels:                deps.TableLabels,
			ReadPurchaseOrder:          deps.ReadPurchaseOrder,
			ListPurchaseOrderLineItems: deps.ListPurchaseOrderLineItems,
		}
		m.PurchaseOrderLineItemTable = purchaseorderdetail.NewLineItemTableView(lineItemTableDeps)
	}

	// Line item action views (nil-guarded — only built when CreatePurchaseOrderLineItem is provided)
	if deps.CreatePurchaseOrderLineItem != nil {
		lineItemActionDeps := &purchaseorderlineitem.Deps{
			Routes:                      deps.Routes,
			Labels:                      deps.Labels,
			CreatePurchaseOrderLineItem: deps.CreatePurchaseOrderLineItem,
			ReadPurchaseOrderLineItem:   deps.ReadPurchaseOrderLineItem,
			UpdatePurchaseOrderLineItem: deps.UpdatePurchaseOrderLineItem,
			DeletePurchaseOrderLineItem: deps.DeletePurchaseOrderLineItem,
		}
		m.PurchaseOrderLineItemAdd = purchaseorderlineitem.NewAddAction(lineItemActionDeps)
		m.PurchaseOrderLineItemEdit = purchaseorderlineitem.NewEditAction(lineItemActionDeps)
		m.PurchaseOrderLineItemRemove = purchaseorderlineitem.NewRemoveAction(lineItemActionDeps)
	}

	// Confirm-receipt action (nil-guarded — only built when ReadPurchaseOrder + ListPurchaseOrderLineItems + UpdatePurchaseOrderLineItem are provided)
	if deps.ReadPurchaseOrder != nil && deps.ListPurchaseOrderLineItems != nil && deps.UpdatePurchaseOrderLineItem != nil {
		receiptDeps := &purchaseorderreceipt.Deps{
			Routes:                      deps.Routes,
			Labels:                      deps.Labels,
			ReadPurchaseOrder:           deps.ReadPurchaseOrder,
			UpdatePurchaseOrder:         deps.UpdatePurchaseOrder,
			ListPurchaseOrderLineItems:  deps.ListPurchaseOrderLineItems,
			ReadPurchaseOrderLineItem:   deps.ReadPurchaseOrderLineItem,
			UpdatePurchaseOrderLineItem: deps.UpdatePurchaseOrderLineItem,
			CreateInventoryMovement:     deps.CreateInventoryMovement,
			ReadInventoryItem:           deps.ReadInventoryItem,
			UpdateInventoryItem:         deps.UpdateInventoryItem,
		}
		m.PurchaseOrderConfirmReceipt = purchaseorderreceipt.NewConfirmReceiptAction(receiptDeps)
	}

	return m
}

// RegisterRoutes registers all purchase order routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	// List routes (nil-guarded)
	if m.PurchaseOrderList != nil {
		r.GET(m.routes.PurchaseOrderListURL, m.PurchaseOrderList)
		r.GET(m.routes.PurchaseOrderTableURL, m.PurchaseOrderList)
	}

	// CRUD action routes (nil-guarded)
	if m.PurchaseOrderAdd != nil {
		r.GET(m.routes.PurchaseOrderAddURL, m.PurchaseOrderAdd)
		r.POST(m.routes.PurchaseOrderAddURL, m.PurchaseOrderAdd)
		r.GET(m.routes.PurchaseOrderEditURL, m.PurchaseOrderEdit)
		r.POST(m.routes.PurchaseOrderEditURL, m.PurchaseOrderEdit)
		r.POST(m.routes.PurchaseOrderDeleteURL, m.PurchaseOrderDelete)
		r.POST(m.routes.PurchaseOrderSetStatusURL, m.PurchaseOrderSetStatus)
	}

	// Detail routes (nil-guarded)
	if m.PurchaseOrderDetail != nil {
		r.GET(m.routes.PurchaseOrderDetailURL, m.PurchaseOrderDetail)
		r.GET(m.routes.PurchaseOrderTabActionURL, m.PurchaseOrderTabAction)
	}

	// Line item table route (nil-guarded)
	if m.PurchaseOrderLineItemTable != nil {
		r.GET(m.routes.PurchaseOrderLineItemTableURL, m.PurchaseOrderLineItemTable)
	}

	// Line item action routes (nil-guarded)
	if m.PurchaseOrderLineItemAdd != nil {
		r.GET(m.routes.PurchaseOrderLineItemAddURL, m.PurchaseOrderLineItemAdd)
		r.POST(m.routes.PurchaseOrderLineItemAddURL, m.PurchaseOrderLineItemAdd)
		r.GET(m.routes.PurchaseOrderLineItemEditURL, m.PurchaseOrderLineItemEdit)
		r.POST(m.routes.PurchaseOrderLineItemEditURL, m.PurchaseOrderLineItemEdit)
		r.POST(m.routes.PurchaseOrderLineItemRemoveURL, m.PurchaseOrderLineItemRemove)
	}

	// Confirm-receipt action routes (nil-guarded)
	if m.PurchaseOrderConfirmReceipt != nil {
		r.GET(m.routes.PurchaseOrderConfirmReceiptURL, m.PurchaseOrderConfirmReceipt)
		r.POST(m.routes.PurchaseOrderConfirmReceiptURL, m.PurchaseOrderConfirmReceipt)
	}
}
