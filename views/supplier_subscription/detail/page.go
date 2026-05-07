package detail

import (
	"context"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	suppliersubscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_subscription"
)

// DetailViewDeps holds view dependencies for the supplier_subscription detail page.
type DetailViewDeps struct {
	Routes       centymo.SupplierSubscriptionRoutes
	Labels       centymo.SupplierSubscriptionLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadSupplierSubscription            func(ctx context.Context, req *suppliersubscriptionpb.ReadSupplierSubscriptionRequest) (*suppliersubscriptionpb.ReadSupplierSubscriptionResponse, error)
	GetSupplierSubscriptionItemPageData func(ctx context.Context, req *suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataRequest) (*suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataResponse, error)
}

// TabItem represents one tab in the detail page tab strip.
type TabItem struct {
	Key    string
	Label  string
	Active bool
}

// PageData holds the data for the supplier_subscription detail page.
type PageData struct {
	types.PageData
	ActiveTab string
	TabItems  []TabItem
	Record    *suppliersubscriptionpb.SupplierSubscription
	EditURL   string
	DeleteURL string
}

func loadRecord(ctx context.Context, deps *DetailViewDeps, id string) (*suppliersubscriptionpb.SupplierSubscription, error) {
	if deps.GetSupplierSubscriptionItemPageData != nil {
		resp, err := deps.GetSupplierSubscriptionItemPageData(ctx, &suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataRequest{
			SupplierSubscriptionId: id,
		})
		if err != nil || resp == nil || resp.GetSupplierSubscription() == nil {
			return nil, err
		}
		return resp.GetSupplierSubscription(), nil
	}
	resp, err := deps.ReadSupplierSubscription(ctx, &suppliersubscriptionpb.ReadSupplierSubscriptionRequest{
		Data: &suppliersubscriptionpb.SupplierSubscription{Id: id},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return nil, err
	}
	return resp.GetData()[0], nil
}

// NewView creates the supplier_subscription detail page view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		record, err := loadRecord(ctx, deps, id)
		if err != nil || record == nil {
			log.Printf("Failed to load supplier subscription detail %s: %v", id, err)
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}

		l := deps.Labels
		tabItems := buildTabs(l, activeTab)

		pageTitle := record.GetName()
		if pageTitle == "" {
			pageTitle = record.GetCode()
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          pageTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    pageTitle,
				HeaderSubtitle: l.Detail.InfoSection,
				HeaderIcon:     "icon-refresh-cw",
				CommonLabels:   deps.CommonLabels,
			},
			ActiveTab: activeTab,
			TabItems:  tabItems,
			Record:    record,
			EditURL:   route.ResolveURL(deps.Routes.EditURL, "id", id),
			DeleteURL: deps.Routes.DeleteURL,
		}
		return view.OK("supplier-subscription-detail", pageData)
	})
}

// NewTabAction handles HTMX tab-swap requests for the detail page.
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

		l := deps.Labels
		pageData := &PageData{
			ActiveTab: tab,
			TabItems:  buildTabs(l, tab),
			Record:    record,
			EditURL:   route.ResolveURL(deps.Routes.EditURL, "id", id),
			DeleteURL: deps.Routes.DeleteURL,
		}
		return view.OK("supplier-subscription-detail-tab-"+tab, pageData)
	})
}

func buildTabs(l centymo.SupplierSubscriptionLabels, activeTab string) []TabItem {
	tabs := []struct {
		Key   string
		Label string
	}{
		{"info", l.Tabs.Info},
		{"cost_plan", l.Tabs.CostPlan},
		{"activity", l.Tabs.Activity},
	}
	items := make([]TabItem, 0, len(tabs))
	for _, t := range tabs {
		items = append(items, TabItem{
			Key:    t.Key,
			Label:  t.Label,
			Active: t.Key == activeTab,
		})
	}
	return items
}
