package detail

import (
	"context"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	costplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_plan"
	supplierproductcostplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_cost_plan"
)

// DetailViewDeps holds view dependencies for the cost_plan detail page.
type DetailViewDeps struct {
	Routes       centymo.CostPlanRoutes
	Labels       centymo.CostPlanLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// SupplierProductCostPlan labels for the inline editor.
	ProductCostLabels centymo.SupplierProductCostPlanLabels

	ReadCostPlan            func(ctx context.Context, req *costplanpb.ReadCostPlanRequest) (*costplanpb.ReadCostPlanResponse, error)
	GetCostPlanItemPageData func(ctx context.Context, req *costplanpb.GetCostPlanItemPageDataRequest) (*costplanpb.GetCostPlanItemPageDataResponse, error)

	// SupplierProductCostPlan CRUD for the inline lines editor.
	CreateSupplierProductCostPlan func(ctx context.Context, req *supplierproductcostplanpb.CreateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.CreateSupplierProductCostPlanResponse, error)
	ReadSupplierProductCostPlan   func(ctx context.Context, req *supplierproductcostplanpb.ReadSupplierProductCostPlanRequest) (*supplierproductcostplanpb.ReadSupplierProductCostPlanResponse, error)
	UpdateSupplierProductCostPlan func(ctx context.Context, req *supplierproductcostplanpb.UpdateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.UpdateSupplierProductCostPlanResponse, error)
	DeleteSupplierProductCostPlan func(ctx context.Context, req *supplierproductcostplanpb.DeleteSupplierProductCostPlanRequest) (*supplierproductcostplanpb.DeleteSupplierProductCostPlanResponse, error)

	// SearchSupplierProductPlanURL for the inline cost-line drawer autocomplete.
	SearchSupplierProductPlanURL string
}

// TabItem represents one tab in the detail page.
type TabItem struct {
	Key    string
	Label  string
	Active bool
}

// PageData holds the data for the cost_plan detail page.
type PageData struct {
	types.PageData
	ActiveTab string
	TabItems  []TabItem
	Record    *costplanpb.CostPlan
	EditURL   string
	DeleteURL string
}

func loadRecord(ctx context.Context, deps *DetailViewDeps, id string) (*costplanpb.CostPlan, error) {
	if deps.GetCostPlanItemPageData != nil {
		resp, err := deps.GetCostPlanItemPageData(ctx, &costplanpb.GetCostPlanItemPageDataRequest{
			CostPlanId: id,
		})
		if err != nil || resp == nil || resp.GetCostPlan() == nil {
			return nil, err
		}
		return resp.GetCostPlan(), nil
	}
	resp, err := deps.ReadCostPlan(ctx, &costplanpb.ReadCostPlanRequest{
		Data: &costplanpb.CostPlan{Id: id},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return nil, err
	}
	return resp.GetData()[0], nil
}

// NewView creates the cost_plan detail page view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}
		record, err := loadRecord(ctx, deps, id)
		if err != nil || record == nil {
			log.Printf("Failed to load cost plan detail %s: %v", id, err)
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
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ActiveTab: activeTab,
			TabItems:  buildTabs(l, activeTab),
			Record:    record,
			EditURL:   route.ResolveURL(deps.Routes.EditURL, "id", id),
			DeleteURL: deps.Routes.DeleteURL,
		}
		return view.OK("cost-plan-detail", pageData)
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
		return view.OK("cost-plan-detail-tab-"+tab, pageData)
	})
}

func buildTabs(l centymo.CostPlanLabels, activeTab string) []TabItem {
	tabs := []struct {
		Key   string
		Label string
	}{
		{"info", l.Tabs.Info},
		{"lines", l.Tabs.Lines},
		{"linked_subscriptions", l.Tabs.LinkedSubscriptions},
		{"activity", l.Tabs.Activity},
	}
	items := make([]TabItem, 0, len(tabs))
	for _, t := range tabs {
		items = append(items, TabItem{Key: t.Key, Label: t.Label, Active: t.Key == activeTab})
	}
	return items
}
