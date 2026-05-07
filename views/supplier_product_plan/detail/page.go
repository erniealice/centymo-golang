package detail

import (
	"context"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	supplierproductplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_plan"
)

// DetailViewDeps holds view dependencies for the supplier_product_plan detail page.
type DetailViewDeps struct {
	Routes       centymo.SupplierProductPlanRoutes
	Labels       centymo.SupplierProductPlanLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadSupplierProductPlan            func(ctx context.Context, req *supplierproductplanpb.ReadSupplierProductPlanRequest) (*supplierproductplanpb.ReadSupplierProductPlanResponse, error)
	GetSupplierProductPlanItemPageData func(ctx context.Context, req *supplierproductplanpb.GetSupplierProductPlanItemPageDataRequest) (*supplierproductplanpb.GetSupplierProductPlanItemPageDataResponse, error)
}

// TabItem represents one tab in the detail page.
type TabItem struct {
	Key    string
	Label  string
	Active bool
}

// PageData holds the data for the supplier_product_plan detail page.
type PageData struct {
	types.PageData
	ActiveTab string
	TabItems  []TabItem
	Record    *supplierproductplanpb.SupplierProductPlan
	EditURL   string
	DeleteURL string
}

func loadRecord(ctx context.Context, deps *DetailViewDeps, id string) (*supplierproductplanpb.SupplierProductPlan, error) {
	if deps.GetSupplierProductPlanItemPageData != nil {
		resp, err := deps.GetSupplierProductPlanItemPageData(ctx, &supplierproductplanpb.GetSupplierProductPlanItemPageDataRequest{
			SupplierProductPlanId: id,
		})
		if err != nil || resp == nil || resp.GetSupplierProductPlan() == nil {
			return nil, err
		}
		return resp.GetSupplierProductPlan(), nil
	}
	resp, err := deps.ReadSupplierProductPlan(ctx, &supplierproductplanpb.ReadSupplierProductPlanRequest{
		Data: &supplierproductplanpb.SupplierProductPlan{Id: id},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return nil, err
	}
	return resp.GetData()[0], nil
}

// NewView creates the supplier_product_plan detail page view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}
		record, err := loadRecord(ctx, deps, id)
		if err != nil || record == nil {
			log.Printf("Failed to load supplier product plan detail %s: %v", id, err)
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          record.GetName(),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    record.GetName(),
				HeaderSubtitle: l.Detail.InfoSection,
				HeaderIcon:     "icon-box",
				CommonLabels:   deps.CommonLabels,
			},
			ActiveTab: activeTab,
			TabItems:  buildTabs(l, activeTab),
			Record:    record,
			EditURL:   route.ResolveURL(deps.Routes.EditURL, "id", id),
			DeleteURL: deps.Routes.DeleteURL,
		}
		return view.OK("supplier-product-plan-detail", pageData)
	})
}

// NewTabAction handles HTMX tab-swap requests.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}
		record, err := loadRecord(ctx, deps, id)
		if err != nil || record == nil {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		pageData := &PageData{
			ActiveTab: tab,
			TabItems:  buildTabs(deps.Labels, tab),
			Record:    record,
			EditURL:   route.ResolveURL(deps.Routes.EditURL, "id", id),
			DeleteURL: deps.Routes.DeleteURL,
		}
		return view.OK("supplier-product-plan-detail-tab-"+tab, pageData)
	})
}

func buildTabs(l centymo.SupplierProductPlanLabels, activeTab string) []TabItem {
	tabs := []struct {
		Key   string
		Label string
	}{
		{"info", l.Tabs.Info},
		{"cost_plan_lines", l.Tabs.CostPlanLines},
		{"activity", l.Tabs.Activity},
	}
	items := make([]TabItem, 0, len(tabs))
	for _, t := range tabs {
		items = append(items, TabItem{Key: t.Key, Label: t.Label, Active: t.Key == activeTab})
	}
	return items
}
