package price_plan

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	priceplanaction "github.com/erniealice/centymo-golang/views/price_plan/action"
	priceplanlist "github.com/erniealice/centymo-golang/views/price_plan/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

// ModuleDeps holds all dependencies for the price_plan module.
type ModuleDeps struct {
	Routes       centymo.PricePlanRoutes
	Labels       centymo.PricePlanLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListPricePlans  func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ReadPricePlan   func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	CreatePricePlan func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)
	UpdatePricePlan func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)

	ListLocations func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	ListPlans     func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
}

// Module holds all constructed price_plan views.
type Module struct {
	routes        centymo.PricePlanRoutes
	Dashboard     view.View
	List          view.View
	Table         view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	BulkDelete    view.View
	SetStatus     view.View
	BulkSetStatus view.View
}

// NewModule creates the price_plan module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &priceplanaction.Deps{
		Routes:          deps.Routes,
		Labels:          deps.Labels,
		CreatePricePlan: deps.CreatePricePlan,
		ReadPricePlan:   deps.ReadPricePlan,
		UpdatePricePlan: deps.UpdatePricePlan,
		DeletePricePlan: deps.DeletePricePlan,
		ListLocations:   deps.ListLocations,
		ListPlans:       deps.ListPlans,
	}

	listView := priceplanlist.NewView(&priceplanlist.ListViewDeps{
		Routes:         deps.Routes,
		ListPricePlans: deps.ListPricePlans,
		Labels:         deps.Labels,
		CommonLabels:   deps.CommonLabels,
		TableLabels:    deps.TableLabels,
	})

	tableView := priceplanlist.NewTableView(&priceplanlist.ListViewDeps{
		Routes:         deps.Routes,
		ListPricePlans: deps.ListPricePlans,
		Labels:         deps.Labels,
		CommonLabels:   deps.CommonLabels,
		TableLabels:    deps.TableLabels,
	})

	return &Module{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           priceplanaction.NewAddAction(actionDeps),
		Edit:          priceplanaction.NewEditAction(actionDeps),
		Delete:        priceplanaction.NewDeleteAction(actionDeps),
		BulkDelete:    priceplanaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     priceplanaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: priceplanaction.NewBulkSetStatusAction(actionDeps),
	}
}

// RegisterRoutes registers all price_plan routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
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
}
