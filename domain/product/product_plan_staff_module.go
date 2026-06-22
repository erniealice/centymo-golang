package product

import (
	"context"

	productplanstaffaction "github.com/erniealice/centymo-golang/domain/product/product_plan_staff/action"
	productplanstaffdetail "github.com/erniealice/centymo-golang/domain/product/product_plan_staff/detail"
	productplanstafflist "github.com/erniealice/centymo-golang/domain/product/product_plan_staff/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/product/product_plan_staff"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
)

// ProductPlanStaffModuleDeps holds all dependencies for the
// product_plan_staff module (the staff-eligibility pool for a product_plan).
type ProductPlanStaffModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListProductPlanStaffs  func(ctx context.Context, req *productplanstaffpb.ListProductPlanStaffsRequest) (*productplanstaffpb.ListProductPlanStaffsResponse, error)
	ReadProductPlanStaff   func(ctx context.Context, req *productplanstaffpb.ReadProductPlanStaffRequest) (*productplanstaffpb.ReadProductPlanStaffResponse, error)
	CreateProductPlanStaff func(ctx context.Context, req *productplanstaffpb.CreateProductPlanStaffRequest) (*productplanstaffpb.CreateProductPlanStaffResponse, error)
	UpdateProductPlanStaff func(ctx context.Context, req *productplanstaffpb.UpdateProductPlanStaffRequest) (*productplanstaffpb.UpdateProductPlanStaffResponse, error)
	DeleteProductPlanStaff func(ctx context.Context, req *productplanstaffpb.DeleteProductPlanStaffRequest) (*productplanstaffpb.DeleteProductPlanStaffResponse, error)

	// Optional reference checker; nil disables delete gating.
	GetProductPlanStaffInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// ProductPlanStaffModule holds all constructed product_plan_staff views.
type ProductPlanStaffModule struct {
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

// NewProductPlanStaffModule creates the product_plan_staff module with all
// views wired.
func NewProductPlanStaffModule(deps *ProductPlanStaffModuleDeps) *ProductPlanStaffModule {
	actionDeps := &productplanstaffaction.Deps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CreateProductPlanStaff:      deps.CreateProductPlanStaff,
		ReadProductPlanStaff:        deps.ReadProductPlanStaff,
		UpdateProductPlanStaff:      deps.UpdateProductPlanStaff,
		DeleteProductPlanStaff:      deps.DeleteProductPlanStaff,
		GetProductPlanStaffInUseIDs: deps.GetProductPlanStaffInUseIDs,
	}

	listDeps := &productplanstafflist.ListViewDeps{
		Routes:                      deps.Routes,
		ListProductPlanStaffs:       deps.ListProductPlanStaffs,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
		GetProductPlanStaffInUseIDs: deps.GetProductPlanStaffInUseIDs,
	}
	listView := productplanstafflist.NewView(listDeps)
	tableView := productplanstafflist.NewTableView(listDeps)

	detailDeps := &productplanstaffdetail.DetailViewDeps{
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
		ReadProductPlanStaff: deps.ReadProductPlanStaff,
	}

	return &ProductPlanStaffModule{
		routes:        deps.Routes,
		Dashboard:     listView,
		List:          listView,
		Table:         tableView,
		Add:           productplanstaffaction.NewAddAction(actionDeps),
		Edit:          productplanstaffaction.NewEditAction(actionDeps),
		Delete:        productplanstaffaction.NewDeleteAction(actionDeps),
		BulkDelete:    productplanstaffaction.NewBulkDeleteAction(actionDeps),
		SetStatus:     productplanstaffaction.NewSetStatusAction(actionDeps),
		BulkSetStatus: productplanstaffaction.NewBulkSetStatusAction(actionDeps),
		Detail:        productplanstaffdetail.NewView(detailDeps),
		TabAction:     productplanstaffdetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all product_plan_staff routes.
func (m *ProductPlanStaffModule) RegisterRoutes(r view.RouteRegistrar) {
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
