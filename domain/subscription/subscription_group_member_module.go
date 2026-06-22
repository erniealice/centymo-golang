package subscription

import (
	"context"

	subscriptiongroupmemberaction "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_member/action"
	subscriptiongroupmemberdetail "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_member/detail"
	subscriptiongroupmemberlist "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_member/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_member"
	subscriptiongroupmemberpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_member"
)

// SubscriptionGroupMemberModuleDeps holds all dependencies for the
// subscription_group_member module.
type SubscriptionGroupMemberModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListSubscriptionGroupMembers     func(ctx context.Context, req *subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest) (*subscriptiongroupmemberpb.ListSubscriptionGroupMembersResponse, error)
	ReadSubscriptionGroupMember      func(ctx context.Context, req *subscriptiongroupmemberpb.ReadSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.ReadSubscriptionGroupMemberResponse, error)
	CreateSubscriptionGroupMember    func(ctx context.Context, req *subscriptiongroupmemberpb.CreateSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.CreateSubscriptionGroupMemberResponse, error)
	UpdateSubscriptionGroupMember    func(ctx context.Context, req *subscriptiongroupmemberpb.UpdateSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.UpdateSubscriptionGroupMemberResponse, error)
	DeleteSubscriptionGroupMember    func(ctx context.Context, req *subscriptiongroupmemberpb.DeleteSubscriptionGroupMemberRequest) (*subscriptiongroupmemberpb.DeleteSubscriptionGroupMemberResponse, error)

	// Optional reference checker; nil disables delete gating.
	GetSubscriptionGroupMemberInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// SubscriptionGroupMemberModule holds all constructed subscription_group_member views.
type SubscriptionGroupMemberModule struct {
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

// NewSubscriptionGroupMemberModule creates the subscription_group_member module
// with all views wired.
func NewSubscriptionGroupMemberModule(deps *SubscriptionGroupMemberModuleDeps) *SubscriptionGroupMemberModule {
	actionDeps := &subscriptiongroupmemberaction.Deps{
		Routes:                             deps.Routes,
		Labels:                             deps.Labels,
		CreateSubscriptionGroupMember:      deps.CreateSubscriptionGroupMember,
		ReadSubscriptionGroupMember:        deps.ReadSubscriptionGroupMember,
		UpdateSubscriptionGroupMember:      deps.UpdateSubscriptionGroupMember,
		DeleteSubscriptionGroupMember:      deps.DeleteSubscriptionGroupMember,
		GetSubscriptionGroupMemberInUseIDs: deps.GetSubscriptionGroupMemberInUseIDs,
	}

	listDeps := &subscriptiongroupmemberlist.ListViewDeps{
		Routes:                             deps.Routes,
		ListSubscriptionGroupMembers:       deps.ListSubscriptionGroupMembers,
		Labels:                             deps.Labels,
		CommonLabels:                       deps.CommonLabels,
		TableLabels:                        deps.TableLabels,
		GetSubscriptionGroupMemberInUseIDs: deps.GetSubscriptionGroupMemberInUseIDs,
	}
	listView := subscriptiongroupmemberlist.NewView(listDeps)
	tableView := subscriptiongroupmemberlist.NewTableView(listDeps)

	detailDeps := &subscriptiongroupmemberdetail.DetailViewDeps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
		ReadSubscriptionGroupMember: deps.ReadSubscriptionGroupMember,
	}

	return &SubscriptionGroupMemberModule{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           subscriptiongroupmemberaction.NewAddAction(actionDeps),
		Edit:          subscriptiongroupmemberaction.NewEditAction(actionDeps),
		Delete:        subscriptiongroupmemberaction.NewDeleteAction(actionDeps),
		BulkDelete:    subscriptiongroupmemberaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     subscriptiongroupmemberaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: subscriptiongroupmemberaction.NewBulkSetStatusAction(actionDeps),
		Detail:        subscriptiongroupmemberdetail.NewView(detailDeps),
		TabAction:     subscriptiongroupmemberdetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all subscription_group_member routes.
func (m *SubscriptionGroupMemberModule) RegisterRoutes(r view.RouteRegistrar) {
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
