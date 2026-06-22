package subscription

import (
	"context"

	subscriptiongroupaction "github.com/erniealice/centymo-golang/domain/subscription/subscription_group/action"
	subscriptiongroupdetail "github.com/erniealice/centymo-golang/domain/subscription/subscription_group/detail"
	subscriptiongrouplist "github.com/erniealice/centymo-golang/domain/subscription/subscription_group/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
)

// SubscriptionGroupModuleDeps holds all dependencies for the
// subscription_group module (the education "section / cohort" cohort).
type SubscriptionGroupModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListSubscriptionGroups  func(ctx context.Context, req *subscriptiongrouppb.ListSubscriptionGroupsRequest) (*subscriptiongrouppb.ListSubscriptionGroupsResponse, error)
	ReadSubscriptionGroup   func(ctx context.Context, req *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error)
	CreateSubscriptionGroup func(ctx context.Context, req *subscriptiongrouppb.CreateSubscriptionGroupRequest) (*subscriptiongrouppb.CreateSubscriptionGroupResponse, error)
	UpdateSubscriptionGroup func(ctx context.Context, req *subscriptiongrouppb.UpdateSubscriptionGroupRequest) (*subscriptiongrouppb.UpdateSubscriptionGroupResponse, error)
	DeleteSubscriptionGroup func(ctx context.Context, req *subscriptiongrouppb.DeleteSubscriptionGroupRequest) (*subscriptiongrouppb.DeleteSubscriptionGroupResponse, error)

	// Program (plan) + period (price_schedule) pickers / display lookups.
	ListPlans          func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)

	// Optional reference checker; nil disables delete gating (no ref-checker
	// method exists for subscription_group yet).
	GetSubscriptionGroupInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// SubscriptionGroupModule holds all constructed subscription_group views.
type SubscriptionGroupModule struct {
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

// NewSubscriptionGroupModule creates the subscription_group module with all
// views wired.
func NewSubscriptionGroupModule(deps *SubscriptionGroupModuleDeps) *SubscriptionGroupModule {
	actionDeps := &subscriptiongroupaction.Deps{
		Routes:                       deps.Routes,
		Labels:                       deps.Labels,
		CreateSubscriptionGroup:      deps.CreateSubscriptionGroup,
		ReadSubscriptionGroup:        deps.ReadSubscriptionGroup,
		UpdateSubscriptionGroup:      deps.UpdateSubscriptionGroup,
		DeleteSubscriptionGroup:      deps.DeleteSubscriptionGroup,
		ListPlans:                    deps.ListPlans,
		ListPriceSchedules:           deps.ListPriceSchedules,
		GetSubscriptionGroupInUseIDs: deps.GetSubscriptionGroupInUseIDs,
	}

	listDeps := &subscriptiongrouplist.ListViewDeps{
		Routes:                       deps.Routes,
		ListSubscriptionGroups:       deps.ListSubscriptionGroups,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
		TableLabels:                  deps.TableLabels,
		GetSubscriptionGroupInUseIDs: deps.GetSubscriptionGroupInUseIDs,
	}
	listView := subscriptiongrouplist.NewView(listDeps)
	tableView := subscriptiongrouplist.NewTableView(listDeps)

	detailDeps := &subscriptiongroupdetail.DetailViewDeps{
		Routes:                deps.Routes,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		TableLabels:           deps.TableLabels,
		ReadSubscriptionGroup: deps.ReadSubscriptionGroup,
		ListPlans:             deps.ListPlans,
		ListPriceSchedules:    deps.ListPriceSchedules,
	}

	return &SubscriptionGroupModule{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           subscriptiongroupaction.NewAddAction(actionDeps),
		Edit:          subscriptiongroupaction.NewEditAction(actionDeps),
		Delete:        subscriptiongroupaction.NewDeleteAction(actionDeps),
		BulkDelete:    subscriptiongroupaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     subscriptiongroupaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: subscriptiongroupaction.NewBulkSetStatusAction(actionDeps),
		Detail:        subscriptiongroupdetail.NewView(detailDeps),
		TabAction:     subscriptiongroupdetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all subscription_group routes.
func (m *SubscriptionGroupModule) RegisterRoutes(r view.RouteRegistrar) {
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
