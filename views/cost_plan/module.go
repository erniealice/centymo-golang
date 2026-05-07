package cost_plan

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	costplanaction "github.com/erniealice/centymo-golang/views/cost_plan/action"
	costplandetail "github.com/erniealice/centymo-golang/views/cost_plan/detail"
	costplanlist "github.com/erniealice/centymo-golang/views/cost_plan/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	costplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_plan"
	supplierproductcostplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_cost_plan"
)

// ModuleDeps holds all dependencies for the cost_plan module.
type ModuleDeps struct {
	Routes       centymo.CostPlanRoutes
	Labels       centymo.CostPlanLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// SupplierProductCostPlan labels for the inline detail tab editor.
	ProductCostLabels centymo.SupplierProductCostPlanLabels

	CreateCostPlan          func(ctx context.Context, req *costplanpb.CreateCostPlanRequest) (*costplanpb.CreateCostPlanResponse, error)
	ReadCostPlan            func(ctx context.Context, req *costplanpb.ReadCostPlanRequest) (*costplanpb.ReadCostPlanResponse, error)
	UpdateCostPlan          func(ctx context.Context, req *costplanpb.UpdateCostPlanRequest) (*costplanpb.UpdateCostPlanResponse, error)
	DeleteCostPlan          func(ctx context.Context, req *costplanpb.DeleteCostPlanRequest) (*costplanpb.DeleteCostPlanResponse, error)
	GetCostPlanListPageData func(ctx context.Context, req *costplanpb.GetCostPlanListPageDataRequest) (*costplanpb.GetCostPlanListPageDataResponse, error)
	GetCostPlanItemPageData func(ctx context.Context, req *costplanpb.GetCostPlanItemPageDataRequest) (*costplanpb.GetCostPlanItemPageDataResponse, error)

	// SetCostPlanActive performs a raw DB update to toggle active.
	SetCostPlanActive func(ctx context.Context, id string, active bool) error

	// SupplierProductCostPlan CRUD (inline editor in the Lines tab).
	CreateSupplierProductCostPlan func(ctx context.Context, req *supplierproductcostplanpb.CreateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.CreateSupplierProductCostPlanResponse, error)
	ReadSupplierProductCostPlan   func(ctx context.Context, req *supplierproductcostplanpb.ReadSupplierProductCostPlanRequest) (*supplierproductcostplanpb.ReadSupplierProductCostPlanResponse, error)
	UpdateSupplierProductCostPlan func(ctx context.Context, req *supplierproductcostplanpb.UpdateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.UpdateSupplierProductCostPlanResponse, error)
	DeleteSupplierProductCostPlan func(ctx context.Context, req *supplierproductcostplanpb.DeleteSupplierProductCostPlanRequest) (*supplierproductcostplanpb.DeleteSupplierProductCostPlanResponse, error)

	// Autocomplete search URLs threaded into form/drawers.
	SearchSupplierPlanURL        string
	SearchCostScheduleURL        string
	SearchSupplierProductPlanURL string
}

// Module holds all constructed cost_plan views.
type Module struct {
	routes             centymo.CostPlanRoutes
	Dashboard          view.View
	List               view.View
	Table              view.View
	Add                view.View
	Edit               view.View
	Delete             view.View
	BulkDelete         view.View
	SetStatus          view.View
	BulkSetStatus      view.View
	Detail             view.View
	TabAction          view.View
	ProductCostAdd     view.View
	ProductCostEdit    view.View
	ProductCostDelete  view.View
}

// NewModule creates the cost_plan module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &costplanaction.Deps{
		Routes:                  deps.Routes,
		Labels:                  deps.Labels,
		CommonLabels:            deps.CommonLabels,
		CreateCostPlan:          deps.CreateCostPlan,
		ReadCostPlan:            deps.ReadCostPlan,
		UpdateCostPlan:          deps.UpdateCostPlan,
		DeleteCostPlan:          deps.DeleteCostPlan,
		GetCostPlanItemPageData: deps.GetCostPlanItemPageData,
		SetCostPlanActive:       deps.SetCostPlanActive,
		SearchSupplierPlanURL:   deps.SearchSupplierPlanURL,
		SearchCostScheduleURL:   deps.SearchCostScheduleURL,
	}

	listDeps := &costplanlist.ListViewDeps{
		Routes:                  deps.Routes,
		GetCostPlanListPageData: deps.GetCostPlanListPageData,
		Labels:                  deps.Labels,
		CommonLabels:            deps.CommonLabels,
		TableLabels:             deps.TableLabels,
	}
	listView := costplanlist.NewView(listDeps)
	tableView := costplanlist.NewTableView(listDeps)

	detailDeps := &costplandetail.DetailViewDeps{
		Routes:                        deps.Routes,
		Labels:                        deps.Labels,
		ProductCostLabels:             deps.ProductCostLabels,
		CommonLabels:                  deps.CommonLabels,
		TableLabels:                   deps.TableLabels,
		ReadCostPlan:                  deps.ReadCostPlan,
		GetCostPlanItemPageData:       deps.GetCostPlanItemPageData,
		CreateSupplierProductCostPlan: deps.CreateSupplierProductCostPlan,
		ReadSupplierProductCostPlan:   deps.ReadSupplierProductCostPlan,
		UpdateSupplierProductCostPlan: deps.UpdateSupplierProductCostPlan,
		DeleteSupplierProductCostPlan: deps.DeleteSupplierProductCostPlan,
		SearchSupplierProductPlanURL:  deps.SearchSupplierProductPlanURL,
	}

	return &Module{
		routes:            deps.Routes,
		Dashboard:         listView,
		List:              listView,
		Table:             tableView,
		Add:               costplanaction.NewAddAction(actionDeps),
		Edit:              costplanaction.NewEditAction(actionDeps),
		Delete:            costplanaction.NewDeleteAction(actionDeps),
		BulkDelete:        costplanaction.NewBulkDeleteAction(actionDeps),
		SetStatus:         costplanaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:     costplanaction.NewBulkSetStatusAction(actionDeps),
		Detail:            costplandetail.NewView(detailDeps),
		TabAction:         costplandetail.NewTabAction(detailDeps),
		ProductCostAdd:    nil, // stub — inline SPCP add action wired via block.go
		ProductCostEdit:   nil, // stub
		ProductCostDelete: nil, // stub
	}
}

// RegisterRoutes registers all cost_plan routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	if m.Detail != nil && m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil && m.routes.TabActionURL != "" {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
	if m.ProductCostAdd != nil && m.routes.ProductCostAddURL != "" {
		r.GET(m.routes.ProductCostAddURL, m.ProductCostAdd)
		r.POST(m.routes.ProductCostAddURL, m.ProductCostAdd)
	}
	if m.ProductCostEdit != nil && m.routes.ProductCostEditURL != "" {
		r.GET(m.routes.ProductCostEditURL, m.ProductCostEdit)
		r.POST(m.routes.ProductCostEditURL, m.ProductCostEdit)
	}
	if m.ProductCostDelete != nil && m.routes.ProductCostDeleteURL != "" {
		r.POST(m.routes.ProductCostDeleteURL, m.ProductCostDelete)
	}
}
