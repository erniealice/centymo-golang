package detail

import (
	"context"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	supplierplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_plan"
)

// DetailViewDeps holds view dependencies for the supplier_plan detail page.
type DetailViewDeps struct {
	Routes       centymo.SupplierPlanRoutes
	Labels       centymo.SupplierPlanLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadSupplierPlan            func(ctx context.Context, req *supplierplanpb.ReadSupplierPlanRequest) (*supplierplanpb.ReadSupplierPlanResponse, error)
	GetSupplierPlanItemPageData func(ctx context.Context, req *supplierplanpb.GetSupplierPlanItemPageDataRequest) (*supplierplanpb.GetSupplierPlanItemPageDataResponse, error)
}

// TabItem represents one tab in the detail page.
type TabItem struct {
	Key    string
	Label  string
	Active bool
}

// PageData holds the data for the supplier_plan detail page.
type PageData struct {
	types.PageData
	ActiveTab string
	TabItems  []TabItem
	Record    *supplierplanpb.SupplierPlan
	EditURL   string
	DeleteURL string
}

func loadRecord(ctx context.Context, deps *DetailViewDeps, id string) (*supplierplanpb.SupplierPlan, error) {
	if deps.GetSupplierPlanItemPageData != nil {
		resp, err := deps.GetSupplierPlanItemPageData(ctx, &supplierplanpb.GetSupplierPlanItemPageDataRequest{
			SupplierPlanId: id,
		})
		if err != nil || resp == nil || resp.GetSupplierPlan() == nil {
			return nil, err
		}
		return resp.GetSupplierPlan(), nil
	}
	resp, err := deps.ReadSupplierPlan(ctx, &supplierplanpb.ReadSupplierPlanRequest{
		Data: &supplierplanpb.SupplierPlan{Id: id},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return nil, err
	}
	return resp.GetData()[0], nil
}

// NewView creates the supplier_plan detail page view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}
		record, err := loadRecord(ctx, deps, id)
		if err != nil || record == nil {
			log.Printf("Failed to load supplier plan detail %s: %v", id, err)
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
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ActiveTab: activeTab,
			TabItems:  buildTabs(l, activeTab),
			Record:    record,
			EditURL:   route.ResolveURL(deps.Routes.EditURL, "id", id),
			DeleteURL: deps.Routes.DeleteURL,
		}
		return view.OK("supplier-plan-detail", pageData)
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
		return view.OK("supplier-plan-detail-tab-"+tab, pageData)
	})
}

func buildTabs(l centymo.SupplierPlanLabels, activeTab string) []TabItem {
	tabs := []struct {
		Key   string
		Label string
	}{
		{"info", l.Tabs.Info},
		{"cost_plans", l.Tabs.CostPlans},
		{"product_plans", l.Tabs.ProductPlans},
		{"activity", l.Tabs.Activity},
	}
	items := make([]TabItem, 0, len(tabs))
	for _, t := range tabs {
		items = append(items, TabItem{Key: t.Key, Label: t.Label, Active: t.Key == activeTab})
	}
	return items
}
