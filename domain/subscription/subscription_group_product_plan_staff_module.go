package subscription

import (
	"context"

	sgppsaction "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan_staff/action"
	sgppsdetail "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan_staff/detail"
	sgppslist "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan_staff/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan_staff"
	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan_staff/form"
	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

// SubscriptionGroupProductPlanStaffModuleDeps holds all dependencies for the
// subscription_group_product_plan_staff module (the class-edge).
type SubscriptionGroupProductPlanStaffModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListSubscriptionGroupProductPlanStaffs  func(ctx context.Context, req *sgppspb.ListSubscriptionGroupProductPlanStaffsRequest) (*sgppspb.ListSubscriptionGroupProductPlanStaffsResponse, error)
	ReadSubscriptionGroupProductPlanStaff   func(ctx context.Context, req *sgppspb.ReadSubscriptionGroupProductPlanStaffRequest) (*sgppspb.ReadSubscriptionGroupProductPlanStaffResponse, error)
	CreateSubscriptionGroupProductPlanStaff func(ctx context.Context, req *sgppspb.CreateSubscriptionGroupProductPlanStaffRequest) (*sgppspb.CreateSubscriptionGroupProductPlanStaffResponse, error)
	UpdateSubscriptionGroupProductPlanStaff func(ctx context.Context, req *sgppspb.UpdateSubscriptionGroupProductPlanStaffRequest) (*sgppspb.UpdateSubscriptionGroupProductPlanStaffResponse, error)
	DeleteSubscriptionGroupProductPlanStaff func(ctx context.Context, req *sgppspb.DeleteSubscriptionGroupProductPlanStaffRequest) (*sgppspb.DeleteSubscriptionGroupProductPlanStaffResponse, error)

	// Optional FK picker loaders — nil disables the picker.
	ListSubscriptionGroupOptions func(ctx context.Context) []form.Pair
	ListProductPlanOptions       func(ctx context.Context) []form.Pair
	ListStaffOptions             func(ctx context.Context) []form.Pair

	// Optional reference checker; nil disables delete gating.
	GetSubscriptionGroupProductPlanStaffInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// SubscriptionGroupProductPlanStaffModule holds all constructed views.
type SubscriptionGroupProductPlanStaffModule struct {
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

// NewSubscriptionGroupProductPlanStaffModule creates the module with all views
// wired.
func NewSubscriptionGroupProductPlanStaffModule(deps *SubscriptionGroupProductPlanStaffModuleDeps) *SubscriptionGroupProductPlanStaffModule {
	actionDeps := &sgppsaction.Deps{
		Routes:                                       deps.Routes,
		Labels:                                       deps.Labels,
		CreateSubscriptionGroupProductPlanStaff:      deps.CreateSubscriptionGroupProductPlanStaff,
		ReadSubscriptionGroupProductPlanStaff:        deps.ReadSubscriptionGroupProductPlanStaff,
		UpdateSubscriptionGroupProductPlanStaff:      deps.UpdateSubscriptionGroupProductPlanStaff,
		DeleteSubscriptionGroupProductPlanStaff:      deps.DeleteSubscriptionGroupProductPlanStaff,
		GetSubscriptionGroupProductPlanStaffInUseIDs: deps.GetSubscriptionGroupProductPlanStaffInUseIDs,
		ListSubscriptionGroupOptions:                 deps.ListSubscriptionGroupOptions,
		ListProductPlanOptions:                       deps.ListProductPlanOptions,
		ListStaffOptions:                             deps.ListStaffOptions,
	}

	listDeps := &sgppslist.ListViewDeps{
		Routes:                                 deps.Routes,
		ListSubscriptionGroupProductPlanStaffs: deps.ListSubscriptionGroupProductPlanStaffs,
		Labels:                                 deps.Labels,
		CommonLabels:                           deps.CommonLabels,
		TableLabels:                            deps.TableLabels,
		GetSubscriptionGroupProductPlanStaffInUseIDs: deps.GetSubscriptionGroupProductPlanStaffInUseIDs,
	}
	listView := sgppslist.NewView(listDeps)
	tableView := sgppslist.NewTableView(listDeps)

	detailDeps := &sgppsdetail.DetailViewDeps{
		Routes:                                deps.Routes,
		Labels:                                deps.Labels,
		CommonLabels:                          deps.CommonLabels,
		TableLabels:                           deps.TableLabels,
		ReadSubscriptionGroupProductPlanStaff: deps.ReadSubscriptionGroupProductPlanStaff,
	}

	return &SubscriptionGroupProductPlanStaffModule{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           sgppsaction.NewAddAction(actionDeps),
		Edit:          sgppsaction.NewEditAction(actionDeps),
		Delete:        sgppsaction.NewDeleteAction(actionDeps),
		BulkDelete:    sgppsaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     sgppsaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: sgppsaction.NewBulkSetStatusAction(actionDeps),
		Detail:        sgppsdetail.NewView(detailDeps),
		TabAction:     sgppsdetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all subscription_group_product_plan_staff routes.
func (m *SubscriptionGroupProductPlanStaffModule) RegisterRoutes(r view.RouteRegistrar) {
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
