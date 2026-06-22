package product

import (
	"context"

	plangroupaction "github.com/erniealice/centymo-golang/domain/product/plan_group/action"
	plangroupdetail "github.com/erniealice/centymo-golang/domain/product/plan_group/detail"
	plangrouplist "github.com/erniealice/centymo-golang/domain/product/plan_group/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/product/plan_group"
	plangroupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/plan_group"
)

// PlanGroupModuleDeps holds all dependencies for the plan_group module.
type PlanGroupModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListPlanGroups  func(ctx context.Context, req *plangroupb.ListPlanGroupsRequest) (*plangroupb.ListPlanGroupsResponse, error)
	ReadPlanGroup   func(ctx context.Context, req *plangroupb.ReadPlanGroupRequest) (*plangroupb.ReadPlanGroupResponse, error)
	CreatePlanGroup func(ctx context.Context, req *plangroupb.CreatePlanGroupRequest) (*plangroupb.CreatePlanGroupResponse, error)
	UpdatePlanGroup func(ctx context.Context, req *plangroupb.UpdatePlanGroupRequest) (*plangroupb.UpdatePlanGroupResponse, error)
	DeletePlanGroup func(ctx context.Context, req *plangroupb.DeletePlanGroupRequest) (*plangroupb.DeletePlanGroupResponse, error)

	// Optional reference checker; nil disables delete gating.
	GetPlanGroupInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// PlanGroupModule holds all constructed plan_group views.
type PlanGroupModule struct {
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

// NewPlanGroupModule creates the plan_group module with all views wired.
func NewPlanGroupModule(deps *PlanGroupModuleDeps) *PlanGroupModule {
	actionDeps := &plangroupaction.Deps{
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		CreatePlanGroup:      deps.CreatePlanGroup,
		ReadPlanGroup:        deps.ReadPlanGroup,
		UpdatePlanGroup:      deps.UpdatePlanGroup,
		DeletePlanGroup:      deps.DeletePlanGroup,
		ListPlanGroups:       deps.ListPlanGroups,
		GetPlanGroupInUseIDs: deps.GetPlanGroupInUseIDs,
	}

	listDeps := &plangrouplist.ListViewDeps{
		Routes:               deps.Routes,
		ListPlanGroups:       deps.ListPlanGroups,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
		GetPlanGroupInUseIDs: deps.GetPlanGroupInUseIDs,
	}
	listView := plangrouplist.NewView(listDeps)
	tableView := plangrouplist.NewTableView(listDeps)

	detailDeps := &plangroupdetail.DetailViewDeps{
		Routes:         deps.Routes,
		Labels:         deps.Labels,
		CommonLabels:   deps.CommonLabels,
		TableLabels:    deps.TableLabels,
		ReadPlanGroup:  deps.ReadPlanGroup,
		ListPlanGroups: deps.ListPlanGroups,
	}

	return &PlanGroupModule{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           plangroupaction.NewAddAction(actionDeps),
		Edit:          plangroupaction.NewEditAction(actionDeps),
		Delete:        plangroupaction.NewDeleteAction(actionDeps),
		BulkDelete:    plangroupaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     plangroupaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: plangroupaction.NewBulkSetStatusAction(actionDeps),
		Detail:        plangroupdetail.NewView(detailDeps),
		TabAction:     plangroupdetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all plan_group routes.
func (m *PlanGroupModule) RegisterRoutes(r view.RouteRegistrar) {
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
