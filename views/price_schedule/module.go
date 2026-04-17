package price_schedule

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	pricescheduleaction "github.com/erniealice/centymo-golang/views/price_schedule/action"
	pricescheduledetail "github.com/erniealice/centymo-golang/views/price_schedule/detail"
	priceschedulePlan "github.com/erniealice/centymo-golang/views/price_schedule/detail/plan"
	priceschedulelist "github.com/erniealice/centymo-golang/views/price_schedule/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// ModuleDeps holds all dependencies for the price_schedule module.
type ModuleDeps struct {
	Routes       centymo.PriceScheduleRoutes
	Labels       centymo.PriceScheduleLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListPriceSchedules  func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	ReadPriceSchedule   func(ctx context.Context, req *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error)
	CreatePriceSchedule func(ctx context.Context, req *priceschedulepb.CreatePriceScheduleRequest) (*priceschedulepb.CreatePriceScheduleResponse, error)
	UpdatePriceSchedule func(ctx context.Context, req *priceschedulepb.UpdatePriceScheduleRequest) (*priceschedulepb.UpdatePriceScheduleResponse, error)
	DeletePriceSchedule func(ctx context.Context, req *priceschedulepb.DeletePriceScheduleRequest) (*priceschedulepb.DeletePriceScheduleResponse, error)

	ListLocations   func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	ListPricePlans  func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ListPlans       func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	CreatePricePlan func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)
	ReadPricePlan   func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	UpdatePricePlan func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)

	// Plan detail page (schedule-scoped) — ProductPricePlan CRUD + supporting lists
	ListProducts           func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ListProductPlans       func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProductPricePlans  func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	CreateProductPricePlan func(ctx context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)
	UpdateProductPricePlan func(ctx context.Context, req *productpriceplanpb.UpdateProductPricePlanRequest) (*productpriceplanpb.UpdateProductPricePlanResponse, error)
	DeleteProductPricePlan func(ctx context.Context, req *productpriceplanpb.DeleteProductPricePlanRequest) (*productpriceplanpb.DeleteProductPricePlanResponse, error)

	GetPriceScheduleInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// Module holds all constructed price_schedule views.
type Module struct {
	routes        centymo.PriceScheduleRoutes
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
	PlanAdd       view.View

	PlanDetail              view.View
	PlanTabAction           view.View
	PlanEdit                view.View
	PlanDelete              view.View
	PlanProductPriceAdd     view.View
	PlanProductPriceEdit    view.View
	PlanProductPriceDelete  view.View
}

// NewModule creates the price_schedule module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &pricescheduleaction.Deps{
		Routes:              deps.Routes,
		Labels:              deps.Labels,
		CreatePriceSchedule: deps.CreatePriceSchedule,
		ReadPriceSchedule:   deps.ReadPriceSchedule,
		UpdatePriceSchedule: deps.UpdatePriceSchedule,
		DeletePriceSchedule: deps.DeletePriceSchedule,
		ListLocations:       deps.ListLocations,
	}

	listDeps := &priceschedulelist.ListViewDeps{
		Routes:                   deps.Routes,
		ListPriceSchedules:       deps.ListPriceSchedules,
		ListLocations:            deps.ListLocations,
		Labels:                   deps.Labels,
		CommonLabels:             deps.CommonLabels,
		TableLabels:              deps.TableLabels,
		GetPriceScheduleInUseIDs: deps.GetPriceScheduleInUseIDs,
	}
	listView := priceschedulelist.NewView(listDeps)
	tableView := priceschedulelist.NewTableView(listDeps)

	detailDeps := &pricescheduledetail.DetailViewDeps{
		Routes:            deps.Routes,
		Labels:            deps.Labels,
		CommonLabels:      deps.CommonLabels,
		TableLabels:       deps.TableLabels,
		ReadPriceSchedule: deps.ReadPriceSchedule,
		ListLocations:     deps.ListLocations,
		ListPricePlans:    deps.ListPricePlans,
		ListPlans:         deps.ListPlans,
		CreatePricePlan:   deps.CreatePricePlan,
	}

	planDetailDeps := &priceschedulePlan.DetailViewDeps{
		Routes:                 deps.Routes,
		ScheduleLabels:         deps.Labels,
		PlanLabels:             pricePlanLabelsFromDeps(deps),
		CommonLabels:           deps.CommonLabels,
		TableLabels:            deps.TableLabels,
		ReadPriceSchedule:      deps.ReadPriceSchedule,
		ReadPricePlan:          deps.ReadPricePlan,
		UpdatePricePlan:        deps.UpdatePricePlan,
		DeletePricePlan:        deps.DeletePricePlan,
		ListPlans:              deps.ListPlans,
		ListProducts:           deps.ListProducts,
		ListProductPlans:       deps.ListProductPlans,
		ListProductPricePlans:  deps.ListProductPricePlans,
		CreateProductPricePlan: deps.CreateProductPricePlan,
		UpdateProductPricePlan: deps.UpdateProductPricePlan,
		DeleteProductPricePlan: deps.DeleteProductPricePlan,
	}

	return &Module{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           pricescheduleaction.NewAddAction(actionDeps),
		Edit:          pricescheduleaction.NewEditAction(actionDeps),
		Delete:        pricescheduleaction.NewDeleteAction(actionDeps),
		BulkDelete:    pricescheduleaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     pricescheduleaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: pricescheduleaction.NewBulkSetStatusAction(actionDeps),
		Detail:        pricescheduledetail.NewView(detailDeps),
		TabAction:     pricescheduledetail.NewTabAction(detailDeps),
		PlanAdd:       pricescheduledetail.NewPlanAddAction(detailDeps),

		PlanDetail:             priceschedulePlan.NewView(planDetailDeps),
		PlanTabAction:          priceschedulePlan.NewTabAction(planDetailDeps),
		PlanEdit:               priceschedulePlan.NewEditAction(planDetailDeps),
		PlanDelete:             priceschedulePlan.NewDeleteAction(planDetailDeps),
		PlanProductPriceAdd:    priceschedulePlan.NewProductPriceAddAction(planDetailDeps),
		PlanProductPriceEdit:   priceschedulePlan.NewProductPriceEditAction(planDetailDeps),
		PlanProductPriceDelete: priceschedulePlan.NewProductPriceDeleteAction(planDetailDeps),
	}
}

// pricePlanLabelsFromDeps provides a fallback PricePlanLabels for the schedule-scoped
// plan detail. The module intentionally does not force callers to pass a full PricePlanLabels —
// only error strings are referenced.
func pricePlanLabelsFromDeps(deps *ModuleDeps) centymo.PricePlanLabels {
	return centymo.PricePlanLabels{
		Errors: centymo.PricePlanErrorLabels{
			NotFound:     deps.Labels.Errors.NotFound,
			LoadFailed:   deps.Labels.Errors.LoadFailed,
			Unauthorized: deps.Labels.Errors.Unauthorized,
			CreateFailed: deps.Labels.Errors.CreateFailed,
			UpdateFailed: deps.Labels.Errors.UpdateFailed,
			DeleteFailed: deps.Labels.Errors.DeleteFailed,
		},
	}
}

// RegisterRoutes registers all price_schedule routes.
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

	if m.Detail != nil && m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil && m.routes.TabActionURL != "" {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
	if m.PlanAdd != nil && m.routes.PlanAddURL != "" {
		r.GET(m.routes.PlanAddURL, m.PlanAdd)
		r.POST(m.routes.PlanAddURL, m.PlanAdd)
	}

	// Schedule-scoped price_plan detail page + CRUD
	if m.PlanDetail != nil && m.routes.PlanDetailURL != "" {
		r.GET(m.routes.PlanDetailURL, m.PlanDetail)
	}
	if m.PlanTabAction != nil && m.routes.PlanTabActionURL != "" {
		r.GET(m.routes.PlanTabActionURL, m.PlanTabAction)
	}
	if m.PlanEdit != nil && m.routes.PlanEditURL != "" {
		r.GET(m.routes.PlanEditURL, m.PlanEdit)
		r.POST(m.routes.PlanEditURL, m.PlanEdit)
	}
	if m.PlanDelete != nil && m.routes.PlanDeleteURL != "" {
		r.POST(m.routes.PlanDeleteURL, m.PlanDelete)
	}
	if m.PlanProductPriceAdd != nil && m.routes.PlanProductPriceAddURL != "" {
		r.GET(m.routes.PlanProductPriceAddURL, m.PlanProductPriceAdd)
		r.POST(m.routes.PlanProductPriceAddURL, m.PlanProductPriceAdd)
	}
	if m.PlanProductPriceEdit != nil && m.routes.PlanProductPriceEditURL != "" {
		r.GET(m.routes.PlanProductPriceEditURL, m.PlanProductPriceEdit)
		r.POST(m.routes.PlanProductPriceEditURL, m.PlanProductPriceEdit)
	}
	if m.PlanProductPriceDelete != nil && m.routes.PlanProductPriceDeleteURL != "" {
		r.POST(m.routes.PlanProductPriceDeleteURL, m.PlanProductPriceDelete)
	}
}
