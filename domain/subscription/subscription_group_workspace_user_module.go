package subscription

import (
	"context"

	sgwuaction "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_workspace_user/action"
	sgwudetail "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_workspace_user/detail"
	sgwulist "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_workspace_user/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_workspace_user"
	sgwupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_workspace_user"
)

// SubscriptionGroupWorkspaceUserModuleDeps holds all dependencies for the
// subscription_group_workspace_user module (operator assignment to a cohort).
type SubscriptionGroupWorkspaceUserModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListSubscriptionGroupWorkspaceUsers  func(ctx context.Context, req *sgwupb.ListSubscriptionGroupWorkspaceUsersRequest) (*sgwupb.ListSubscriptionGroupWorkspaceUsersResponse, error)
	ReadSubscriptionGroupWorkspaceUser   func(ctx context.Context, req *sgwupb.ReadSubscriptionGroupWorkspaceUserRequest) (*sgwupb.ReadSubscriptionGroupWorkspaceUserResponse, error)
	CreateSubscriptionGroupWorkspaceUser func(ctx context.Context, req *sgwupb.CreateSubscriptionGroupWorkspaceUserRequest) (*sgwupb.CreateSubscriptionGroupWorkspaceUserResponse, error)
	UpdateSubscriptionGroupWorkspaceUser func(ctx context.Context, req *sgwupb.UpdateSubscriptionGroupWorkspaceUserRequest) (*sgwupb.UpdateSubscriptionGroupWorkspaceUserResponse, error)
	DeleteSubscriptionGroupWorkspaceUser func(ctx context.Context, req *sgwupb.DeleteSubscriptionGroupWorkspaceUserRequest) (*sgwupb.DeleteSubscriptionGroupWorkspaceUserResponse, error)

	// Optional reference checker; nil disables delete gating.
	GetSubscriptionGroupWorkspaceUserInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// SubscriptionGroupWorkspaceUserModule holds all constructed views.
type SubscriptionGroupWorkspaceUserModule struct {
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

// NewSubscriptionGroupWorkspaceUserModule creates the module with all views wired.
func NewSubscriptionGroupWorkspaceUserModule(deps *SubscriptionGroupWorkspaceUserModuleDeps) *SubscriptionGroupWorkspaceUserModule {
	actionDeps := &sgwuaction.Deps{
		Routes:                                    deps.Routes,
		Labels:                                    deps.Labels,
		CreateSubscriptionGroupWorkspaceUser:      deps.CreateSubscriptionGroupWorkspaceUser,
		ReadSubscriptionGroupWorkspaceUser:        deps.ReadSubscriptionGroupWorkspaceUser,
		UpdateSubscriptionGroupWorkspaceUser:      deps.UpdateSubscriptionGroupWorkspaceUser,
		DeleteSubscriptionGroupWorkspaceUser:      deps.DeleteSubscriptionGroupWorkspaceUser,
		GetSubscriptionGroupWorkspaceUserInUseIDs: deps.GetSubscriptionGroupWorkspaceUserInUseIDs,
	}

	listDeps := &sgwulist.ListViewDeps{
		Routes:                              deps.Routes,
		ListSubscriptionGroupWorkspaceUsers: deps.ListSubscriptionGroupWorkspaceUsers,
		Labels:                              deps.Labels,
		CommonLabels:                        deps.CommonLabels,
		TableLabels:                         deps.TableLabels,
		GetSubscriptionGroupWorkspaceUserInUseIDs: deps.GetSubscriptionGroupWorkspaceUserInUseIDs,
	}
	listView := sgwulist.NewView(listDeps)
	tableView := sgwulist.NewTableView(listDeps)

	detailDeps := &sgwudetail.DetailViewDeps{
		Routes:                             deps.Routes,
		Labels:                             deps.Labels,
		CommonLabels:                       deps.CommonLabels,
		TableLabels:                        deps.TableLabels,
		ReadSubscriptionGroupWorkspaceUser: deps.ReadSubscriptionGroupWorkspaceUser,
	}

	return &SubscriptionGroupWorkspaceUserModule{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           sgwuaction.NewAddAction(actionDeps),
		Edit:          sgwuaction.NewEditAction(actionDeps),
		Delete:        sgwuaction.NewDeleteAction(actionDeps),
		BulkDelete:    sgwuaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     sgwuaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: sgwuaction.NewBulkSetStatusAction(actionDeps),
		Detail:        sgwudetail.NewView(detailDeps),
		TabAction:     sgwudetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all subscription_group_workspace_user routes.
func (m *SubscriptionGroupWorkspaceUserModule) RegisterRoutes(r view.RouteRegistrar) {
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
