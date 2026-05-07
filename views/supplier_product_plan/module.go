package supplier_product_plan

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	supplierproductplanaction "github.com/erniealice/centymo-golang/views/supplier_product_plan/action"
	supplierproductplandetail "github.com/erniealice/centymo-golang/views/supplier_product_plan/detail"
	supplierproductplanlist "github.com/erniealice/centymo-golang/views/supplier_product_plan/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	supplierproductplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_plan"
)

// ModuleDeps holds all dependencies for the supplier_product_plan module.
type ModuleDeps struct {
	Routes       centymo.SupplierProductPlanRoutes
	Labels       centymo.SupplierProductPlanLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	CreateSupplierProductPlan          func(ctx context.Context, req *supplierproductplanpb.CreateSupplierProductPlanRequest) (*supplierproductplanpb.CreateSupplierProductPlanResponse, error)
	ReadSupplierProductPlan            func(ctx context.Context, req *supplierproductplanpb.ReadSupplierProductPlanRequest) (*supplierproductplanpb.ReadSupplierProductPlanResponse, error)
	UpdateSupplierProductPlan          func(ctx context.Context, req *supplierproductplanpb.UpdateSupplierProductPlanRequest) (*supplierproductplanpb.UpdateSupplierProductPlanResponse, error)
	DeleteSupplierProductPlan          func(ctx context.Context, req *supplierproductplanpb.DeleteSupplierProductPlanRequest) (*supplierproductplanpb.DeleteSupplierProductPlanResponse, error)
	GetSupplierProductPlanListPageData func(ctx context.Context, req *supplierproductplanpb.GetSupplierProductPlanListPageDataRequest) (*supplierproductplanpb.GetSupplierProductPlanListPageDataResponse, error)
	GetSupplierProductPlanItemPageData func(ctx context.Context, req *supplierproductplanpb.GetSupplierProductPlanItemPageDataRequest) (*supplierproductplanpb.GetSupplierProductPlanItemPageDataResponse, error)

	// SetSupplierProductPlanActive performs a raw DB update to toggle active.
	SetSupplierProductPlanActive func(ctx context.Context, id string, active bool) error

	// Autocomplete search URLs threaded into form.
	SearchSupplierPlanURL string
	SearchProductURL      string
}

// Module holds all constructed supplier_product_plan views.
type Module struct {
	routes        centymo.SupplierProductPlanRoutes
	Dashboard     view.View
	List          view.View
	Table         view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	BulkDelete    view.View
	SetStatus     view.View
	BulkSetStatus view.View
	Detail        view.View
	TabAction     view.View
}

// NewModule creates the supplier_product_plan module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &supplierproductplanaction.Deps{
		Routes:                             deps.Routes,
		Labels:                             deps.Labels,
		CommonLabels:                       deps.CommonLabels,
		CreateSupplierProductPlan:          deps.CreateSupplierProductPlan,
		ReadSupplierProductPlan:            deps.ReadSupplierProductPlan,
		UpdateSupplierProductPlan:          deps.UpdateSupplierProductPlan,
		DeleteSupplierProductPlan:          deps.DeleteSupplierProductPlan,
		GetSupplierProductPlanItemPageData: deps.GetSupplierProductPlanItemPageData,
		SetSupplierProductPlanActive:       deps.SetSupplierProductPlanActive,
		SearchSupplierPlanURL:              deps.SearchSupplierPlanURL,
		SearchProductURL:                   deps.SearchProductURL,
	}

	listDeps := &supplierproductplanlist.ListViewDeps{
		Routes:                             deps.Routes,
		GetSupplierProductPlanListPageData: deps.GetSupplierProductPlanListPageData,
		Labels:                             deps.Labels,
		CommonLabels:                       deps.CommonLabels,
		TableLabels:                        deps.TableLabels,
	}
	listView := supplierproductplanlist.NewView(listDeps)
	tableView := supplierproductplanlist.NewTableView(listDeps)

	detailDeps := &supplierproductplandetail.DetailViewDeps{
		Routes:                             deps.Routes,
		Labels:                             deps.Labels,
		CommonLabels:                       deps.CommonLabels,
		TableLabels:                        deps.TableLabels,
		ReadSupplierProductPlan:            deps.ReadSupplierProductPlan,
		GetSupplierProductPlanItemPageData: deps.GetSupplierProductPlanItemPageData,
	}

	return &Module{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           supplierproductplanaction.NewAddAction(actionDeps),
		Edit:          supplierproductplanaction.NewEditAction(actionDeps),
		Delete:        supplierproductplanaction.NewDeleteAction(actionDeps),
		BulkDelete:    supplierproductplanaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     supplierproductplanaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: supplierproductplanaction.NewBulkSetStatusAction(actionDeps),
		Detail:        supplierproductplandetail.NewView(detailDeps),
		TabAction:     supplierproductplandetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all supplier_product_plan routes.
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
}
