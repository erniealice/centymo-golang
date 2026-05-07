package cost_schedule

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	costscheduleaction "github.com/erniealice/centymo-golang/views/cost_schedule/action"
	costscheduledetail "github.com/erniealice/centymo-golang/views/cost_schedule/detail"
	costschedulelist "github.com/erniealice/centymo-golang/views/cost_schedule/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	costschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_schedule"
)

// ModuleDeps holds all dependencies for the cost_schedule module.
type ModuleDeps struct {
	Routes       centymo.CostScheduleRoutes
	Labels       centymo.CostScheduleLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	CreateCostSchedule          func(ctx context.Context, req *costschedulepb.CreateCostScheduleRequest) (*costschedulepb.CreateCostScheduleResponse, error)
	ReadCostSchedule            func(ctx context.Context, req *costschedulepb.ReadCostScheduleRequest) (*costschedulepb.ReadCostScheduleResponse, error)
	UpdateCostSchedule          func(ctx context.Context, req *costschedulepb.UpdateCostScheduleRequest) (*costschedulepb.UpdateCostScheduleResponse, error)
	DeleteCostSchedule          func(ctx context.Context, req *costschedulepb.DeleteCostScheduleRequest) (*costschedulepb.DeleteCostScheduleResponse, error)
	GetCostScheduleListPageData func(ctx context.Context, req *costschedulepb.GetCostScheduleListPageDataRequest) (*costschedulepb.GetCostScheduleListPageDataResponse, error)
	GetCostScheduleItemPageData func(ctx context.Context, req *costschedulepb.GetCostScheduleItemPageDataRequest) (*costschedulepb.GetCostScheduleItemPageDataResponse, error)

	// SetCostScheduleActive performs a raw DB update to toggle active.
	SetCostScheduleActive func(ctx context.Context, id string, active bool) error
}

// Module holds all constructed cost_schedule views.
type Module struct {
	routes        centymo.CostScheduleRoutes
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

// NewModule creates the cost_schedule module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &costscheduleaction.Deps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		CreateCostSchedule:          deps.CreateCostSchedule,
		ReadCostSchedule:            deps.ReadCostSchedule,
		UpdateCostSchedule:          deps.UpdateCostSchedule,
		DeleteCostSchedule:          deps.DeleteCostSchedule,
		GetCostScheduleItemPageData: deps.GetCostScheduleItemPageData,
		SetCostScheduleActive:       deps.SetCostScheduleActive,
	}

	listDeps := &costschedulelist.ListViewDeps{
		Routes:                      deps.Routes,
		GetCostScheduleListPageData: deps.GetCostScheduleListPageData,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
	}
	listView := costschedulelist.NewView(listDeps)
	tableView := costschedulelist.NewTableView(listDeps)

	detailDeps := &costscheduledetail.DetailViewDeps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
		ReadCostSchedule:            deps.ReadCostSchedule,
		GetCostScheduleItemPageData: deps.GetCostScheduleItemPageData,
	}

	return &Module{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           costscheduleaction.NewAddAction(actionDeps),
		Edit:          costscheduleaction.NewEditAction(actionDeps),
		Delete:        costscheduleaction.NewDeleteAction(actionDeps),
		BulkDelete:    costscheduleaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     costscheduleaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: costscheduleaction.NewBulkSetStatusAction(actionDeps),
		Detail:        costscheduledetail.NewView(detailDeps),
		TabAction:     costscheduledetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all cost_schedule routes.
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
