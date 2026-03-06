package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.DisbursementRoutes
	DB           centymo.DataSource
	Labels       centymo.DisbursementLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PageData holds the data for the disbursement detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Disbursement    map[string]any
	Labels          centymo.DisbursementLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	// Convenience fields for template rendering
	Reference     string
	StatusLabel   string
	StatusVariant string
	Amount        string
	Currency      string

	AuditTable *types.TableConfig
}

// NewView creates the disbursement detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		disbursement, err := deps.DB.Read(ctx, "disbursement", id)
		if err != nil {
			log.Printf("Failed to read disbursement %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load disbursement: %w", err))
		}

		refNumber, _ := disbursement["reference_number"].(string)
		status, _ := disbursement["status"].(string)
		currency, _ := disbursement["currency"].(string)
		amount, _ := disbursement["amount"].(string)
		l := deps.Labels
		headerTitle := l.Detail.TitlePrefix + refNumber

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}
		tabItems := buildTabItems(l, id, deps.Routes)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "cash",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.PageTitle,
				HeaderIcon:     "icon-arrow-up-right",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "disbursement-detail-content",
			Disbursement:    disbursement,
			Labels:          l,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
			Reference:       refNumber,
			StatusLabel:     status,
			StatusVariant:   statusVariant(status),
			Amount:          amount,
			Currency:        currency,
		}

		// Load tab-specific data
		switch activeTab {
		case "info":
			// Disbursement map has everything
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		}

		return view.OK("disbursement-detail", pageData)
	})
}

func buildTabItems(l centymo.DisbursementLabels, id string, routes centymo.DisbursementRoutes) []pyeza.TabItem {
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

		disbursement, err := deps.DB.Read(ctx, "disbursement", id)
		if err != nil {
			log.Printf("Failed to read disbursement %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load disbursement: %w", err))
		}

		status, _ := disbursement["status"].(string)
		currency, _ := disbursement["currency"].(string)
		amount, _ := disbursement["amount"].(string)
		refNumber, _ := disbursement["reference_number"].(string)

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Disbursement:  disbursement,
			Labels:        l,
			ActiveTab:     tab,
			TabItems:      buildTabItems(l, id, deps.Routes),
			Reference:     refNumber,
			StatusLabel:   status,
			StatusVariant: statusVariant(status),
			Amount:        amount,
			Currency:      currency,
		}

		switch tab {
		case "info":
			// disbursement map has everything
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		}

		templateName := "disbursement-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

func buildAuditTable(l centymo.DisbursementLabels, tableLabels types.TableLabels) *types.TableConfig {
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

func statusVariant(status string) string {
	switch status {
	case "draft":
		return "default"
	case "pending":
		return "warning"
	case "approved":
		return "info"
	case "paid":
		return "success"
	case "cancelled":
		return "danger"
	case "overdue":
		return "danger"
	default:
		return "default"
	}
}
