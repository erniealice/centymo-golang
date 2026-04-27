package supplier_contract

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	suppliercontractaction "github.com/erniealice/centymo-golang/views/supplier_contract/action"
	suppliercontractdetail "github.com/erniealice/centymo-golang/views/supplier_contract/detail"
	suppliercontractlist "github.com/erniealice/centymo-golang/views/supplier_contract/list"

	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	suppliercontractlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_line"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the supplier_contract module.
type ModuleDeps struct {
	Routes       centymo.SupplierContractRoutes
	Labels       centymo.SupplierContractLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Core CRUD
	ListSupplierContracts  func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
	ReadSupplierContract   func(ctx context.Context, req *suppliercontractpb.ReadSupplierContractRequest) (*suppliercontractpb.ReadSupplierContractResponse, error)
	CreateSupplierContract func(ctx context.Context, req *suppliercontractpb.CreateSupplierContractRequest) (*suppliercontractpb.CreateSupplierContractResponse, error)
	UpdateSupplierContract func(ctx context.Context, req *suppliercontractpb.UpdateSupplierContractRequest) (*suppliercontractpb.UpdateSupplierContractResponse, error)
	DeleteSupplierContract func(ctx context.Context, req *suppliercontractpb.DeleteSupplierContractRequest) (*suppliercontractpb.DeleteSupplierContractResponse, error)

	// Workflow actions
	SetSupplierContractStatus func(ctx context.Context, id, status string) error

	// Child entity — lines
	ListSupplierContractLines func(ctx context.Context, req *suppliercontractlinepb.ListSupplierContractLinesRequest) (*suppliercontractlinepb.ListSupplierContractLinesResponse, error)

	// Related entities for dropdowns + linked tabs
	ListSuppliers      func(ctx context.Context, req *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error)
	ListPurchaseOrders func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)
	ListExpenditures   func(ctx context.Context, req *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error)
}

// Module holds all constructed supplier_contract views.
type Module struct {
	routes        centymo.SupplierContractRoutes
	List          view.View
	Detail        view.View
	TabAction     view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	SetStatus     view.View
	BulkSetStatus view.View
	Approve       view.View
	Terminate     view.View
}

// NewModule creates the supplier_contract module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &suppliercontractaction.Deps{
		Routes:                    deps.Routes,
		Labels:                    deps.Labels,
		CommonLabels:              deps.CommonLabels,
		CreateSupplierContract:    deps.CreateSupplierContract,
		ReadSupplierContract:      deps.ReadSupplierContract,
		UpdateSupplierContract:    deps.UpdateSupplierContract,
		DeleteSupplierContract:    deps.DeleteSupplierContract,
		SetSupplierContractStatus: deps.SetSupplierContractStatus,
		ListSuppliers:             deps.ListSuppliers,
	}

	detailDeps := &suppliercontractdetail.DetailViewDeps{
		Routes:                    deps.Routes,
		Labels:                    deps.Labels,
		CommonLabels:              deps.CommonLabels,
		TableLabels:               deps.TableLabels,
		ReadSupplierContract:      deps.ReadSupplierContract,
		ListSupplierContractLines: deps.ListSupplierContractLines,
		ListPurchaseOrders:        deps.ListPurchaseOrders,
		ListExpenditures:          deps.ListExpenditures,
	}

	listDeps := &suppliercontractlist.ListViewDeps{
		Routes:                deps.Routes,
		ListSupplierContracts: deps.ListSupplierContracts,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		TableLabels:           deps.TableLabels,
	}

	m := &Module{
		routes:        deps.Routes,
		List:          suppliercontractlist.NewView(listDeps),
		Add:           suppliercontractaction.NewAddAction(actionDeps),
		Edit:          suppliercontractaction.NewEditAction(actionDeps),
		Delete:        suppliercontractaction.NewDeleteAction(actionDeps),
		SetStatus:     suppliercontractaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: suppliercontractaction.NewBulkSetStatusAction(actionDeps),
	}

	if deps.ReadSupplierContract != nil {
		m.Detail = suppliercontractdetail.NewView(detailDeps)
		m.TabAction = suppliercontractdetail.NewTabAction(detailDeps)
		m.Approve = suppliercontractdetail.NewApproveAction(detailDeps)
		m.Terminate = suppliercontractdetail.NewTerminateAction(detailDeps)
	}

	return m
}

// RegisterRoutes registers all supplier_contract routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)

	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)

	if m.Detail != nil {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
	if m.Approve != nil {
		r.POST(m.routes.ApproveURL, m.Approve)
	}
	if m.Terminate != nil {
		r.POST(m.routes.TerminateURL, m.Terminate)
	}
}
