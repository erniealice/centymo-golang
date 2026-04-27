package procurement_request_line

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	procurementrequestlineaction "github.com/erniealice/centymo-golang/views/procurement_request_line/action"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	procurementrequestlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request_line"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the procurement_request_line module.
type ModuleDeps struct {
	Routes       centymo.ProcurementRequestRoutes
	Labels       centymo.ProcurementRequestLabels
	CommonLabels pyeza.CommonLabels

	CreateProcurementRequestLine func(ctx context.Context, req *procurementrequestlinepb.CreateProcurementRequestLineRequest) (*procurementrequestlinepb.CreateProcurementRequestLineResponse, error)
	ReadProcurementRequestLine   func(ctx context.Context, req *procurementrequestlinepb.ReadProcurementRequestLineRequest) (*procurementrequestlinepb.ReadProcurementRequestLineResponse, error)
	UpdateProcurementRequestLine func(ctx context.Context, req *procurementrequestlinepb.UpdateProcurementRequestLineRequest) (*procurementrequestlinepb.UpdateProcurementRequestLineResponse, error)
	DeleteProcurementRequestLine func(ctx context.Context, req *procurementrequestlinepb.DeleteProcurementRequestLineRequest) (*procurementrequestlinepb.DeleteProcurementRequestLineResponse, error)

	// Optional — for the product picker
	ListProducts func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
}

// Module holds all constructed procurement_request_line views.
type Module struct {
	routes centymo.ProcurementRequestRoutes
	Add    view.View
	Edit   view.View
	Delete view.View
}

// NewModule creates the procurement_request_line module.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &procurementrequestlineaction.Deps{
		Routes:                       deps.Routes,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
		CreateProcurementRequestLine: deps.CreateProcurementRequestLine,
		ReadProcurementRequestLine:   deps.ReadProcurementRequestLine,
		UpdateProcurementRequestLine: deps.UpdateProcurementRequestLine,
		DeleteProcurementRequestLine: deps.DeleteProcurementRequestLine,
		ListProducts:                 deps.ListProducts,
	}

	return &Module{
		routes: deps.Routes,
		Add:    procurementrequestlineaction.NewAddAction(actionDeps),
		Edit:   procurementrequestlineaction.NewEditAction(actionDeps),
		Delete: procurementrequestlineaction.NewDeleteAction(actionDeps),
	}
}

// RegisterRoutes registers all procurement_request_line action routes.
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
