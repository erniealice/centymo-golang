package price_plan

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	priceplanaction "github.com/erniealice/centymo-golang/views/price_plan/action"
	priceplandetail "github.com/erniealice/centymo-golang/views/price_plan/detail"
	priceplanlist "github.com/erniealice/centymo-golang/views/price_plan/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	jobtemplatephasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_phase"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// ModuleDeps holds all dependencies for the price_plan module.
type ModuleDeps struct {
	Routes                    centymo.PricePlanRoutes
	Labels                    centymo.PricePlanLabels
	ProductPricePlanLabels    centymo.ProductPricePlanLabels
	PriceScheduleDetailLabels centymo.PriceScheduleDetailLabels // for the PPP drawer's basis banner
	CommonLabels              pyeza.CommonLabels
	TableLabels               types.TableLabels

	ListPricePlans  func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ReadPricePlan   func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	CreatePricePlan func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)
	UpdatePricePlan func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)

	ListPlans          func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)

	GetPricePlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// Detail page — ProductPricePlan CRUD
	ListProductPlans         func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProducts             func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	// ListProductVariants / ListProductOptions / ListProductOptionValues /
	// ListProductVariantOptions feed the enriched catalog-line picker labels
	// ("Product — SKU — Red / Large / Cotton") in the price drawer.
	ListProductVariants       func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	ListProductOptions        func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ListProductOptionValues   func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	ListProductVariantOptions func(ctx context.Context, req *productvariantoptionpb.ListProductVariantOptionsRequest) (*productvariantoptionpb.ListProductVariantOptionsResponse, error)
	ListProductPricePlans    func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	CreateProductPricePlan   func(ctx context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)
	UpdateProductPricePlan   func(ctx context.Context, req *productpriceplanpb.UpdateProductPricePlanRequest) (*productpriceplanpb.UpdateProductPricePlanResponse, error)
	DeleteProductPricePlan   func(ctx context.Context, req *productpriceplanpb.DeleteProductPricePlanRequest) (*productpriceplanpb.DeleteProductPricePlanResponse, error)

	// 2026-04-27 plan-client-scope plan §6.7 — used by the price-plan drawer
	// to resolve the parent PriceSchedule's client name for the info banner.
	ListClientNames func(ctx context.Context) map[string]string

	// 2026-04-29 milestone-billing plan §5 / Phase D — used by the PPP drawer
	// to populate the optional milestone (job_template_phase) select when the
	// parent PricePlan has billing_kind = MILESTONE.
	ReadPlan                           func(ctx context.Context, req *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)
	ListJobTemplatePhasesByJobTemplate func(ctx context.Context, req *jobtemplatephasepb.ListByJobTemplateRequest) (*jobtemplatephasepb.ListByJobTemplateResponse, error)
}

// Module holds all constructed price_plan views.
type Module struct {
	routes                 centymo.PricePlanRoutes
	Dashboard              view.View
	List                   view.View
	Table                  view.View
	Add                    view.View
	Edit                   view.View
	Delete                 view.View
	BulkDelete             view.View
	SetStatus              view.View
	BulkSetStatus          view.View
	Detail                 view.View
	TabAction              view.View
	ProductPriceAdd        view.View
	ProductPriceEdit       view.View
	ProductPriceDelete     view.View
}

// NewModule creates the price_plan module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &priceplanaction.Deps{
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		CreatePricePlan:      deps.CreatePricePlan,
		ReadPricePlan:        deps.ReadPricePlan,
		UpdatePricePlan:      deps.UpdatePricePlan,
		DeletePricePlan:      deps.DeletePricePlan,
		ListPlans:            deps.ListPlans,
		ListPriceSchedules:   deps.ListPriceSchedules,
		GetPricePlanInUseIDs: deps.GetPricePlanInUseIDs,
		// 2026-04-27 plan-client-scope plan §6.7.
		ListClientNames: deps.ListClientNames,
	}

	listDeps := &priceplanlist.ListViewDeps{
		Routes:               deps.Routes,
		ListPricePlans:       deps.ListPricePlans,
		ListPlans:            deps.ListPlans,
		ListPriceSchedules:   deps.ListPriceSchedules,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
		GetPricePlanInUseIDs: deps.GetPricePlanInUseIDs,
	}
	listView := priceplanlist.NewView(listDeps)
	tableView := priceplanlist.NewTableView(listDeps)

	detailDeps := &priceplandetail.DetailViewDeps{
		Routes:                             deps.Routes,
		Labels:                             deps.Labels,
		ProductPricePlanLabels:             deps.ProductPricePlanLabels,
		PriceScheduleDetailLabels:          deps.PriceScheduleDetailLabels,
		CommonLabels:                       deps.CommonLabels,
		TableLabels:                        deps.TableLabels,
		ReadPricePlan:                      deps.ReadPricePlan,
		ListProductPlans:                   deps.ListProductPlans,
		ListProducts:                       deps.ListProducts,
		ListProductVariants:                deps.ListProductVariants,
		ListProductOptions:                 deps.ListProductOptions,
		ListProductOptionValues:            deps.ListProductOptionValues,
		ListProductVariantOptions:          deps.ListProductVariantOptions,
		ListProductPricePlans:              deps.ListProductPricePlans,
		CreateProductPricePlan:             deps.CreateProductPricePlan,
		UpdateProductPricePlan:             deps.UpdateProductPricePlan,
		DeleteProductPricePlan:             deps.DeleteProductPricePlan,
		ListPlans:                          deps.ListPlans,
		ListPriceSchedules:                 deps.ListPriceSchedules,
		ReadPlan:                           deps.ReadPlan,
		ListJobTemplatePhasesByJobTemplate: deps.ListJobTemplatePhasesByJobTemplate,
	}

	return &Module{
		routes:             deps.Routes,
		Dashboard:          listView,
		List:               listView,
		Table:              tableView,
		Add:                priceplanaction.NewAddAction(actionDeps),
		Edit:               priceplanaction.NewEditAction(actionDeps),
		Delete:             priceplanaction.NewDeleteAction(actionDeps),
		BulkDelete:         priceplanaction.NewBulkDeleteAction(actionDeps),
		SetStatus:          priceplanaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:      priceplanaction.NewBulkSetStatusAction(actionDeps),
		Detail:             priceplandetail.NewView(detailDeps),
		TabAction:          priceplandetail.NewTabAction(detailDeps),
		ProductPriceAdd:    priceplandetail.NewProductPriceAddAction(detailDeps),
		ProductPriceEdit:   priceplandetail.NewProductPriceEditAction(detailDeps),
		ProductPriceDelete: priceplandetail.NewProductPriceDeleteAction(detailDeps),
	}
}

// RegisterRoutes registers all price_plan routes.
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

	// Detail page
	if m.Detail != nil {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}

	// ProductPricePlan CRUD within detail
	if m.ProductPriceAdd != nil && m.routes.ProductPriceAddURL != "" {
		r.GET(m.routes.ProductPriceAddURL, m.ProductPriceAdd)
		r.POST(m.routes.ProductPriceAddURL, m.ProductPriceAdd)
	}
	if m.ProductPriceEdit != nil && m.routes.ProductPriceEditURL != "" {
		r.GET(m.routes.ProductPriceEditURL, m.ProductPriceEdit)
		r.POST(m.routes.ProductPriceEditURL, m.ProductPriceEdit)
	}
	if m.ProductPriceDelete != nil && m.routes.ProductPriceDeleteURL != "" {
		r.POST(m.routes.ProductPriceDeleteURL, m.ProductPriceDelete)
	}
}
