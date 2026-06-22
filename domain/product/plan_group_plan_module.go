package product

import (
	"context"

	plangroupplanaction "github.com/erniealice/centymo-golang/domain/product/plan_group_plan/action"
	plangroupplandetail "github.com/erniealice/centymo-golang/domain/product/plan_group_plan/detail"
	plangroupplanlist "github.com/erniealice/centymo-golang/domain/product/plan_group_plan/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/product/plan_group_plan"
	plangroupplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/plan_group_plan"
)

// PlanGroupPlanModuleDeps holds all dependencies for the plan_group_plan module.
type PlanGroupPlanModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListPlanGroupPlans  func(ctx context.Context, req *plangroupplanpb.ListPlanGroupPlansRequest) (*plangroupplanpb.ListPlanGroupPlansResponse, error)
	ReadPlanGroupPlan   func(ctx context.Context, req *plangroupplanpb.ReadPlanGroupPlanRequest) (*plangroupplanpb.ReadPlanGroupPlanResponse, error)
	CreatePlanGroupPlan func(ctx context.Context, req *plangroupplanpb.CreatePlanGroupPlanRequest) (*plangroupplanpb.CreatePlanGroupPlanResponse, error)
	UpdatePlanGroupPlan func(ctx context.Context, req *plangroupplanpb.UpdatePlanGroupPlanRequest) (*plangroupplanpb.UpdatePlanGroupPlanResponse, error)
	DeletePlanGroupPlan func(ctx context.Context, req *plangroupplanpb.DeletePlanGroupPlanRequest) (*plangroupplanpb.DeletePlanGroupPlanResponse, error)

	// Optional reference checker; nil disables delete gating.
	GetPlanGroupPlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// PlanGroupPlanModule holds all constructed plan_group_plan views.
type PlanGroupPlanModule struct {
	routes        epkg.Routes
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

// NewPlanGroupPlanModule creates the plan_group_plan module with all views wired.
func NewPlanGroupPlanModule(deps *PlanGroupPlanModuleDeps) *PlanGroupPlanModule {
	actionDeps := &plangroupplanaction.Deps{
		Routes:                   deps.Routes,
		Labels:                   deps.Labels,
		CreatePlanGroupPlan:      deps.CreatePlanGroupPlan,
		ReadPlanGroupPlan:        deps.ReadPlanGroupPlan,
		UpdatePlanGroupPlan:      deps.UpdatePlanGroupPlan,
		DeletePlanGroupPlan:      deps.DeletePlanGroupPlan,
		GetPlanGroupPlanInUseIDs: deps.GetPlanGroupPlanInUseIDs,
	}

	listDeps := &plangroupplanlist.ListViewDeps{
		Routes:                   deps.Routes,
		ListPlanGroupPlans:       deps.ListPlanGroupPlans,
		Labels:                   deps.Labels,
		CommonLabels:             deps.CommonLabels,
		TableLabels:              deps.TableLabels,
		GetPlanGroupPlanInUseIDs: deps.GetPlanGroupPlanInUseIDs,
	}
	listView := plangroupplanlist.NewView(listDeps)
	tableView := plangroupplanlist.NewTableView(listDeps)

	detailDeps := &plangroupplandetail.DetailViewDeps{
		Routes:            deps.Routes,
		Labels:            deps.Labels,
		CommonLabels:      deps.CommonLabels,
		TableLabels:       deps.TableLabels,
		ReadPlanGroupPlan: deps.ReadPlanGroupPlan,
	}

	return &PlanGroupPlanModule{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           plangroupplanaction.NewAddAction(actionDeps),
		Edit:          plangroupplanaction.NewEditAction(actionDeps),
		Delete:        plangroupplanaction.NewDeleteAction(actionDeps),
		BulkDelete:    plangroupplanaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     plangroupplanaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: plangroupplanaction.NewBulkSetStatusAction(actionDeps),
		Detail:        plangroupplandetail.NewView(detailDeps),
		TabAction:     plangroupplandetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all plan_group_plan routes.
func (m *PlanGroupPlanModule) RegisterRoutes(r view.RouteRegistrar) {
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

	if m.Detail != nil && m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil && m.routes.TabActionURL != "" {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
}
