package supplier_contract_line

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	suppliercontractlineaction "github.com/erniealice/centymo-golang/views/supplier_contract_line/action"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	suppliercontractlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_line"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the supplier_contract_line module.
type ModuleDeps struct {
	Routes       centymo.SupplierContractRoutes
	Labels       centymo.SupplierContractLabels
	CommonLabels pyeza.CommonLabels

	CreateSupplierContractLine func(ctx context.Context, req *suppliercontractlinepb.CreateSupplierContractLineRequest) (*suppliercontractlinepb.CreateSupplierContractLineResponse, error)
	ReadSupplierContractLine   func(ctx context.Context, req *suppliercontractlinepb.ReadSupplierContractLineRequest) (*suppliercontractlinepb.ReadSupplierContractLineResponse, error)
	UpdateSupplierContractLine func(ctx context.Context, req *suppliercontractlinepb.UpdateSupplierContractLineRequest) (*suppliercontractlinepb.UpdateSupplierContractLineResponse, error)
	DeleteSupplierContractLine func(ctx context.Context, req *suppliercontractlinepb.DeleteSupplierContractLineRequest) (*suppliercontractlinepb.DeleteSupplierContractLineResponse, error)

	// Optional — for the product picker in the line drawer form
	ListProducts func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
}

// Module holds all constructed supplier_contract_line views.
type Module struct {
	routes centymo.SupplierContractRoutes
	Add    view.View
	Edit   view.View
	Delete view.View
}

// NewModule creates the supplier_contract_line module.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &suppliercontractlineaction.Deps{
		Routes:                     deps.Routes,
		Labels:                     deps.Labels,
		CommonLabels:               deps.CommonLabels,
		CreateSupplierContractLine: deps.CreateSupplierContractLine,
		ReadSupplierContractLine:   deps.ReadSupplierContractLine,
		UpdateSupplierContractLine: deps.UpdateSupplierContractLine,
		DeleteSupplierContractLine: deps.DeleteSupplierContractLine,
		ListProducts:               deps.ListProducts,
	}

	return &Module{
		routes: deps.Routes,
		Add:    suppliercontractlineaction.NewAddAction(actionDeps),
		Edit:   suppliercontractlineaction.NewEditAction(actionDeps),
		Delete: suppliercontractlineaction.NewDeleteAction(actionDeps),
	}
}

// RegisterRoutes registers all supplier_contract_line action routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	if m.Add != nil {
		r.GET(m.routes.LineAddURL, m.Add)
		r.POST(m.routes.LineAddURL, m.Add)
	}
	if m.Edit != nil {
		r.GET(m.routes.LineEditURL, m.Edit)
		r.POST(m.routes.LineEditURL, m.Edit)
	}
	if m.Delete != nil {
		r.POST(m.routes.LineDeleteURL, m.Delete)
	}
}
