package subscription

import (
	"context"

	pswuaction "github.com/erniealice/centymo-golang/domain/subscription/price_schedule_workspace_user/action"
	pswudetail "github.com/erniealice/centymo-golang/domain/subscription/price_schedule_workspace_user/detail"
	pswulist "github.com/erniealice/centymo-golang/domain/subscription/price_schedule_workspace_user/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/subscription/price_schedule_workspace_user"
	pswupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule_workspace_user"
)

// PriceScheduleWorkspaceUserModuleDeps holds all dependencies for the
// price_schedule_workspace_user module (period-level coordinator access records).
type PriceScheduleWorkspaceUserModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListPriceScheduleWorkspaceUsers  func(ctx context.Context, req *pswupb.ListPriceScheduleWorkspaceUsersRequest) (*pswupb.ListPriceScheduleWorkspaceUsersResponse, error)
	ReadPriceScheduleWorkspaceUser   func(ctx context.Context, req *pswupb.ReadPriceScheduleWorkspaceUserRequest) (*pswupb.ReadPriceScheduleWorkspaceUserResponse, error)
	CreatePriceScheduleWorkspaceUser func(ctx context.Context, req *pswupb.CreatePriceScheduleWorkspaceUserRequest) (*pswupb.CreatePriceScheduleWorkspaceUserResponse, error)
	UpdatePriceScheduleWorkspaceUser func(ctx context.Context, req *pswupb.UpdatePriceScheduleWorkspaceUserRequest) (*pswupb.UpdatePriceScheduleWorkspaceUserResponse, error)
	DeletePriceScheduleWorkspaceUser func(ctx context.Context, req *pswupb.DeletePriceScheduleWorkspaceUserRequest) (*pswupb.DeletePriceScheduleWorkspaceUserResponse, error)

	// Optional reference checker; nil disables delete gating.
	GetPriceScheduleWorkspaceUserInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// PriceScheduleWorkspaceUserModule holds all constructed price_schedule_workspace_user views.
type PriceScheduleWorkspaceUserModule struct {
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

// NewPriceScheduleWorkspaceUserModule creates the price_schedule_workspace_user
// module with all views wired.
func NewPriceScheduleWorkspaceUserModule(deps *PriceScheduleWorkspaceUserModuleDeps) *PriceScheduleWorkspaceUserModule {
	actionDeps := &pswuaction.Deps{
		Routes:                                deps.Routes,
		Labels:                                deps.Labels,
		CreatePriceScheduleWorkspaceUser:      deps.CreatePriceScheduleWorkspaceUser,
		ReadPriceScheduleWorkspaceUser:        deps.ReadPriceScheduleWorkspaceUser,
		UpdatePriceScheduleWorkspaceUser:      deps.UpdatePriceScheduleWorkspaceUser,
		DeletePriceScheduleWorkspaceUser:      deps.DeletePriceScheduleWorkspaceUser,
		GetPriceScheduleWorkspaceUserInUseIDs: deps.GetPriceScheduleWorkspaceUserInUseIDs,
	}

	listDeps := &pswulist.ListViewDeps{
		Routes:                                deps.Routes,
		ListPriceScheduleWorkspaceUsers:       deps.ListPriceScheduleWorkspaceUsers,
		Labels:                                deps.Labels,
		CommonLabels:                          deps.CommonLabels,
		TableLabels:                           deps.TableLabels,
		GetPriceScheduleWorkspaceUserInUseIDs: deps.GetPriceScheduleWorkspaceUserInUseIDs,
	}
	listView := pswulist.NewView(listDeps)
	tableView := pswulist.NewTableView(listDeps)

	detailDeps := &pswudetail.DetailViewDeps{
		Routes:                         deps.Routes,
		Labels:                         deps.Labels,
		CommonLabels:                   deps.CommonLabels,
		TableLabels:                    deps.TableLabels,
		ReadPriceScheduleWorkspaceUser: deps.ReadPriceScheduleWorkspaceUser,
	}

	return &PriceScheduleWorkspaceUserModule{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           pswuaction.NewAddAction(actionDeps),
		Edit:          pswuaction.NewEditAction(actionDeps),
		Delete:        pswuaction.NewDeleteAction(actionDeps),
		BulkDelete:    pswuaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     pswuaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: pswuaction.NewBulkSetStatusAction(actionDeps),
		Detail:        pswudetail.NewView(detailDeps),
		TabAction:     pswudetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all price_schedule_workspace_user routes.
func (m *PriceScheduleWorkspaceUserModule) RegisterRoutes(r view.RouteRegistrar) {
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
