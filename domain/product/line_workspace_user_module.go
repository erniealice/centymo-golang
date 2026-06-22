package product

import (
	"context"

	lineworkspaceuseraction "github.com/erniealice/centymo-golang/domain/product/line_workspace_user/action"
	lineworkspaceuserdetail "github.com/erniealice/centymo-golang/domain/product/line_workspace_user/detail"
	lineworkspaceuserlist "github.com/erniealice/centymo-golang/domain/product/line_workspace_user/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/product/line_workspace_user"
	lineworkspaceuserpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line_workspace_user"
)

// LineWorkspaceUserModuleDeps holds all dependencies for the
// line_workspace_user module.
type LineWorkspaceUserModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListLineWorkspaceUsers  func(ctx context.Context, req *lineworkspaceuserpb.ListLineWorkspaceUsersRequest) (*lineworkspaceuserpb.ListLineWorkspaceUsersResponse, error)
	ReadLineWorkspaceUser   func(ctx context.Context, req *lineworkspaceuserpb.ReadLineWorkspaceUserRequest) (*lineworkspaceuserpb.ReadLineWorkspaceUserResponse, error)
	CreateLineWorkspaceUser func(ctx context.Context, req *lineworkspaceuserpb.CreateLineWorkspaceUserRequest) (*lineworkspaceuserpb.CreateLineWorkspaceUserResponse, error)
	UpdateLineWorkspaceUser func(ctx context.Context, req *lineworkspaceuserpb.UpdateLineWorkspaceUserRequest) (*lineworkspaceuserpb.UpdateLineWorkspaceUserResponse, error)
	DeleteLineWorkspaceUser func(ctx context.Context, req *lineworkspaceuserpb.DeleteLineWorkspaceUserRequest) (*lineworkspaceuserpb.DeleteLineWorkspaceUserResponse, error)

	// Optional reference checker; nil disables delete gating.
	GetLineWorkspaceUserInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// LineWorkspaceUserModule holds all constructed line_workspace_user views.
type LineWorkspaceUserModule struct {
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

// NewLineWorkspaceUserModule creates the line_workspace_user module with all
// views wired.
func NewLineWorkspaceUserModule(deps *LineWorkspaceUserModuleDeps) *LineWorkspaceUserModule {
	actionDeps := &lineworkspaceuseraction.Deps{
		Routes:                       deps.Routes,
		Labels:                       deps.Labels,
		CreateLineWorkspaceUser:      deps.CreateLineWorkspaceUser,
		ReadLineWorkspaceUser:        deps.ReadLineWorkspaceUser,
		UpdateLineWorkspaceUser:      deps.UpdateLineWorkspaceUser,
		DeleteLineWorkspaceUser:      deps.DeleteLineWorkspaceUser,
		GetLineWorkspaceUserInUseIDs: deps.GetLineWorkspaceUserInUseIDs,
	}

	listDeps := &lineworkspaceuserlist.ListViewDeps{
		Routes:                       deps.Routes,
		ListLineWorkspaceUsers:       deps.ListLineWorkspaceUsers,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
		TableLabels:                  deps.TableLabels,
		GetLineWorkspaceUserInUseIDs: deps.GetLineWorkspaceUserInUseIDs,
	}
	listView := lineworkspaceuserlist.NewView(listDeps)
	tableView := lineworkspaceuserlist.NewTableView(listDeps)

	detailDeps := &lineworkspaceuserdetail.DetailViewDeps{
		Routes:                deps.Routes,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		TableLabels:           deps.TableLabels,
		ReadLineWorkspaceUser: deps.ReadLineWorkspaceUser,
	}

	return &LineWorkspaceUserModule{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           lineworkspaceuseraction.NewAddAction(actionDeps),
		Edit:          lineworkspaceuseraction.NewEditAction(actionDeps),
		Delete:        lineworkspaceuseraction.NewDeleteAction(actionDeps),
		BulkDelete:    lineworkspaceuseraction.NewBulkDeleteAction(actionDeps),
		SetStatus:     lineworkspaceuseraction.NewSetStatusAction(actionDeps),
		BulkSetStatus: lineworkspaceuseraction.NewBulkSetStatusAction(actionDeps),
		Detail:        lineworkspaceuserdetail.NewView(detailDeps),
		TabAction:     lineworkspaceuserdetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all line_workspace_user routes.
func (m *LineWorkspaceUserModule) RegisterRoutes(r view.RouteRegistrar) {
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
