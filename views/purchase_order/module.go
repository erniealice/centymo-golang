package purchaseorder

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	purchaseorderlineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order_line_item"

	purchaseorderaction "github.com/erniealice/centymo-golang/views/purchase_order/action"
	purchaseorderdetail "github.com/erniealice/centymo-golang/views/purchase_order/detail"
	purchaseorderlist "github.com/erniealice/centymo-golang/views/purchase_order/list"
)

// ModuleDeps holds all dependencies for the purchase order module.
type ModuleDeps struct {
	Routes       centymo.ExpenditureRoutes
	DB           centymo.DataSource
	Labels       centymo.ExpenditureLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Purchase order CRUD operations (nil-guarded — only built when provided)
	ListPurchaseOrders   func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)
	CreatePurchaseOrder  func(ctx context.Context, req *purchaseorderpb.CreatePurchaseOrderRequest) (*purchaseorderpb.CreatePurchaseOrderResponse, error)
	ReadPurchaseOrder    func(ctx context.Context, req *purchaseorderpb.ReadPurchaseOrderRequest) (*purchaseorderpb.ReadPurchaseOrderResponse, error)
	UpdatePurchaseOrder  func(ctx context.Context, req *purchaseorderpb.UpdatePurchaseOrderRequest) (*purchaseorderpb.UpdatePurchaseOrderResponse, error)
	DeletePurchaseOrder  func(ctx context.Context, req *purchaseorderpb.DeletePurchaseOrderRequest) (*purchaseorderpb.DeletePurchaseOrderResponse, error)

	// Purchase order line item operations (optional — used by detail view)
	ListPurchaseOrderLineItems func(ctx context.Context, req *purchaseorderlineitempb.ListPurchaseOrderLineItemsRequest) (*purchaseorderlineitempb.ListPurchaseOrderLineItemsResponse, error)
}

// Module holds all constructed purchase order views.
type Module struct {
	routes                  centymo.ExpenditureRoutes
	PurchaseOrderList       view.View
	PurchaseOrderAdd        view.View
	PurchaseOrderEdit       view.View
	PurchaseOrderDelete     view.View
	PurchaseOrderSetStatus  view.View
	PurchaseOrderDetail     view.View
	PurchaseOrderTabAction  view.View
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
}
