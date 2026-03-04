package detail

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.CollectionRoutes
	DB           centymo.DataSource
	Labels       centymo.CollectionLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PageData holds the data for the collection detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Collection      map[string]any
	Labels          centymo.CollectionLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	AuditTable      *types.TableConfig
}

// NewView creates the collection detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		collection, err := deps.DB.Read(ctx, "collection", id)
		if err != nil {
			log.Printf("Failed to read collection %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load collection: %w", err))
		}

		refNumber, _ := collection["reference_number"].(string)
		headerTitle := "Collection #" + refNumber

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}

		l := deps.Labels
		tabItems := buildTabItems(l, id, deps.Routes)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "cash",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.PageTitle,
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "collection-detail-content",
			Collection:      collection,
			Labels:          l,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
		}

		switch activeTab {
		case "info":
			// collection map has everything
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		}

		return view.OK("collection-detail", pageData)
	})
}

func buildTabItems(l centymo.CollectionLabels, id string, routes centymo.CollectionRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Detail.TabBasicInfo, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "audit", Label: l.Detail.TabAuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
	}
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		collection, err := deps.DB.Read(ctx, "collection", id)
		if err != nil {
			log.Printf("Failed to read collection %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load collection: %w", err))
		}

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Collection: collection,
			Labels:     l,
			ActiveTab:  tab,
			TabItems:   buildTabItems(l, id, deps.Routes),
		}

		switch tab {
		case "info":
			// collection map has everything
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		}

		templateName := "collection-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

// buildAuditTable creates the audit trail table.
func buildAuditTable(l centymo.CollectionLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: l.Detail.Date, Sortable: true, Width: "160px"},
		{Key: "action", Label: l.Detail.AuditAction, Sortable: true},
		{Key: "user", Label: l.Detail.AuditUser, Sortable: true, Width: "180px"},
	}

	rows := []types.TableRow{}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   "audit-trail-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowEntries:          true,
		DefaultSortColumn:    "date",
		DefaultSortDirection: "desc",
		Labels:               tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.AuditEmptyTitle,
			Message: l.Detail.AuditEmptyMessage,
		},
	}
	types.ApplyTableSettings(cfg)

	return cfg
}
