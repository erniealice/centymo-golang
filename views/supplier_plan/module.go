package supplier_plan

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	supplierplanaction "github.com/erniealice/centymo-golang/views/supplier_plan/action"
	supplierplandetail "github.com/erniealice/centymo-golang/views/supplier_plan/detail"
	supplierplanlist "github.com/erniealice/centymo-golang/views/supplier_plan/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	supplierplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_plan"
)

// ModuleDeps holds all dependencies for the supplier_plan module.
type ModuleDeps struct {
	Routes       centymo.SupplierPlanRoutes
	Labels       centymo.SupplierPlanLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	CreateSupplierPlan          func(ctx context.Context, req *supplierplanpb.CreateSupplierPlanRequest) (*supplierplanpb.CreateSupplierPlanResponse, error)
	ReadSupplierPlan            func(ctx context.Context, req *supplierplanpb.ReadSupplierPlanRequest) (*supplierplanpb.ReadSupplierPlanResponse, error)
	UpdateSupplierPlan          func(ctx context.Context, req *supplierplanpb.UpdateSupplierPlanRequest) (*supplierplanpb.UpdateSupplierPlanResponse, error)
	DeleteSupplierPlan          func(ctx context.Context, req *supplierplanpb.DeleteSupplierPlanRequest) (*supplierplanpb.DeleteSupplierPlanResponse, error)
	GetSupplierPlanListPageData func(ctx context.Context, req *supplierplanpb.GetSupplierPlanListPageDataRequest) (*supplierplanpb.GetSupplierPlanListPageDataResponse, error)
	GetSupplierPlanItemPageData func(ctx context.Context, req *supplierplanpb.GetSupplierPlanItemPageDataRequest) (*supplierplanpb.GetSupplierPlanItemPageDataResponse, error)

	// SetSupplierPlanActive performs a raw DB update to toggle active.
	SetSupplierPlanActive func(ctx context.Context, id string, active bool) error

	// SearchSupplierURL is threaded into the form for the supplier autocomplete.
	SearchSupplierURL string
}

// Module holds all constructed supplier_plan views.
type Module struct {
	routes        centymo.SupplierPlanRoutes
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

// NewModule creates the supplier_plan module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &supplierplanaction.Deps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		CreateSupplierPlan:          deps.CreateSupplierPlan,
		ReadSupplierPlan:            deps.ReadSupplierPlan,
		UpdateSupplierPlan:          deps.UpdateSupplierPlan,
		DeleteSupplierPlan:          deps.DeleteSupplierPlan,
		GetSupplierPlanItemPageData: deps.GetSupplierPlanItemPageData,
		SetSupplierPlanActive:       deps.SetSupplierPlanActive,
		SearchSupplierURL:           deps.SearchSupplierURL,
	}

	listDeps := &supplierplanlist.ListViewDeps{
		Routes:                     deps.Routes,
		GetSupplierPlanListPageData: deps.GetSupplierPlanListPageData,
		Labels:                     deps.Labels,
		CommonLabels:               deps.CommonLabels,
		TableLabels:                deps.TableLabels,
	}
	listView := supplierplanlist.NewView(listDeps)
	tableView := supplierplanlist.NewTableView(listDeps)

	detailDeps := &supplierplandetail.DetailViewDeps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
		ReadSupplierPlan:            deps.ReadSupplierPlan,
		GetSupplierPlanItemPageData: deps.GetSupplierPlanItemPageData,
	}

	return &Module{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           supplierplanaction.NewAddAction(actionDeps),
		Edit:          supplierplanaction.NewEditAction(actionDeps),
		Delete:        supplierplanaction.NewDeleteAction(actionDeps),
		BulkDelete:    supplierplanaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     supplierplanaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: supplierplanaction.NewBulkSetStatusAction(actionDeps),
		Detail:        supplierplandetail.NewView(detailDeps),
		TabAction:     supplierplandetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all supplier_plan routes.
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
