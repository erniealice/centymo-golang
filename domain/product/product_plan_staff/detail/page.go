package detail

import (
	"context"
	"fmt"
	"log"
	"time"

	product_plan_staff "github.com/erniealice/centymo-golang/domain/product/product_plan_staff"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
)

// DetailViewDeps holds view dependencies for the product_plan_staff detail page.
type DetailViewDeps struct {
	Routes               product_plan_staff.Routes
	Labels               product_plan_staff.Labels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels
	ReadProductPlanStaff func(ctx context.Context, req *productplanstaffpb.ReadProductPlanStaffRequest) (*productplanstaffpb.ReadProductPlanStaffResponse, error)
}

// PageData holds the data for the product_plan_staff detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Record          *productplanstaffpb.ProductPlanStaff
	Labels          product_plan_staff.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID            string
	StaffID       string
	ProductPlanID string
	Role          string
	Status        string
	StatusVariant string
	CreatedDate   string
	ModifiedDate  string
}

// NewView creates the product_plan_staff detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan_staff", "read") {
			return view.Forbidden("product_plan_staff:read")
		}
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("product-plan-staff-detail", pageData)
	})
}

// NewTabAction handles GET /action/product-plan-staff/{id}/tab/{tab}.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}
		pageData, err := buildPageData(ctx, deps, id, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("product-plan-staff-tab-"+tab, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadProductPlanStaff(ctx, &productplanstaffpb.ReadProductPlanStaffRequest{
		Data: &productplanstaffpb.ProductPlanStaff{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read product plan staff %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	ps := data[0]

	l := deps.Labels

	role := ps.GetRole()
	if role == "" {
		role = l.Detail.NoRole
	}

	status := "active"
	statusVariant := "success"
	if !ps.GetActive() {
		status = "inactive"
		statusVariant = "warning"
	}

	tabItems := []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info,
			Href:  route.ResolveURL(deps.Routes.DetailURL, "id", id) + "?tab=info",
			HxGet: route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "info"),
			Icon:  "icon-info"},
	}

	tz := types.LocationFromContext(ctx)
	createdDate := ""
	if ms := ps.GetDateCreated(); ms > 0 {
		createdDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}
	modifiedDate := ""
	if ms := ps.GetDateModified(); ms > 0 {
		modifiedDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}

	// Use role as the subtitle if set; otherwise fall back to productPlanId.
	headerSubtitle := ps.GetProductPlanId()
	if headerSubtitle == "" {
		headerSubtitle = l.Detail.NoSubtitle
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          ps.GetStaffId(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    ps.GetStaffId(),
			HeaderSubtitle: headerSubtitle,
			HeaderIcon:     "icon-user",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "product-plan-staff-detail-content",
		Record:          ps,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		StaffID:         ps.GetStaffId(),
		ProductPlanID:   ps.GetProductPlanId(),
		Role:            role,
		Status:          status,
		StatusVariant:   statusVariant,
		CreatedDate:     createdDate,
		ModifiedDate:    modifiedDate,
	}
	return pageData, nil
}
