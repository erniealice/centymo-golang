// Package dashboard renders the Service dashboard (product_kind=service
// surface). Phase 5 — sits ABOVE the service-mount product list.
//
// The dashboard's data comes from the espyna product dashboard adapters,
// orchestrated through the GetPageData callback.  Workspace-scoped via the
// callback (orchestrator pulls workspace_id from request context).
//
// Note: this package lives under views/product/ — services are not their
// own domain, they're products filtered to product_kind="service". The
// `service-dashboard` template name and "service" sidebar key are used
// because that's how the surface is presented to users.
package dashboard

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"sort"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Stats holds tile values for the service dashboard.
type Stats struct {
	TotalActive      int64
	TopRevenueName   string
	TopRevenueValue  int64 // centavos — placeholder until line-item join lands
	LineCount        int64
	RecentlyAddedCnt int64
}

// LineRow is a row in the by-line custom table.
type LineRow struct {
	LineID string
	Count  int64
}

// TopRevenueRow is one row in the top-revenue services widget.
// `Total` is centavos; in this initial slice it's a placeholder rank value
// (line-item revenue join is deferred — see plan).
type TopRevenueRow struct {
	ProductID   string
	ProductName string
	Total       int64
}

// Request is the input to the GetPageData callback.
type Request struct {
	Now time.Time
}

// Response is the projection the view consumes.
type Response struct {
	Stats Stats

	// Services-by-line for chart-bar (parallel slices).
	LineLabels []string
	LineValues []float64

	// Top revenue services. Deferred slice — populated as best-effort by the
	// orchestrator; falls back to top-by-recency if revenue join unavailable.
	TopRevenue []TopRevenueRow

	// Recent service additions.
	Recent []*productpb.Product
}

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.ProductRoutes
	Labels       centymo.ProductLabels
	CommonLabels pyeza.CommonLabels

	// GetPageData is nil-safe; orchestrator wraps the espyna product dashboard
	// adapter (kind="service", workspace_id from ctx).
	GetPageData func(ctx context.Context, req *Request) (*Response, error)
}

// PageData is what the dashboard template receives.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the service dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.ServiceDashboard
		now := time.Now()

		var resp *Response
		if deps.GetPageData != nil {
			r, err := deps.GetPageData(ctx, &Request{Now: now})
			if err == nil {
				resp = r
			}
		}
		if resp == nil {
			resp = &Response{}
		}

		// By-line bar chart — values are counts, not centavos.
		lineChart := &types.ChartData{
			Labels: resp.LineLabels,
			Series: []types.ChartSeries{{
				Name:   l.WidgetByLine,
				Values: resp.LineValues,
				Color:  "amber",
			}},
		}
		if len(lineChart.Labels) == 0 {
			lineChart.Labels = []string{"-"}
			lineChart.Series[0].Values = []float64{0}
		}
		lineChart.AutoScale()

		// Recent additions list.
		recentItems := make([]types.ActivityItem, 0, len(resp.Recent))
		for i, p := range resp.Recent {
			title := p.GetName()
			if title == "" {
				title = l.NewService
			}
			desc := p.GetDescription()
			recentItems = append(recentItems, types.ActivityItem{
				IconName:    "icon-briefcase",
				IconVariant: "client",
				Title:       title,
				Description: desc,
				Time:        p.GetDateCreatedString(),
				TestID:      fmt.Sprintf("service-list-item-%d", i),
			})
		}

		// Top-revenue table — sort defensively in case the orchestrator returns unsorted.
		topRev := append([]TopRevenueRow{}, resp.TopRevenue...)
		sort.Slice(topRev, func(i, j int) bool { return topRev[i].Total > topRev[j].Total })
		if len(topRev) > 5 {
			topRev = topRev[:5]
		}

		topRevenueLabel := l.EmptyTopRevenue
		if resp.Stats.TopRevenueName != "" {
			topRevenueLabel = resp.Stats.TopRevenueName
		}

		dash := types.DashboardData{
			Title:    l.Title,
			Icon:     "icon-briefcase",
			Subtitle: l.Subtitle,
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickNew, Href: deps.Routes.AddURL, Variant: "primary", TestID: "service-action-new"},
				{Icon: "icon-package", Label: l.QuickBundleBuilder, Href: deps.Routes.ListURL, TestID: "service-action-bundle"},
				{Icon: "icon-tag", Label: l.QuickTagService, Href: deps.Routes.ListURL, TestID: "service-action-tag"},
				{Icon: "icon-dollar-sign", Label: l.QuickPriceSchedule, Href: deps.Routes.ListURL, TestID: "service-action-pricing"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-briefcase", Value: fmt.Sprintf("%d", resp.Stats.TotalActive), Label: l.StatTotalActive, Color: "terracotta", TestID: "service-stat-total"},
				{Icon: "icon-trending-up", Value: topRevenueLabel, Label: l.StatTopRevenue, Color: "sage", TestID: "service-stat-top-revenue"},
				{Icon: "icon-grid", Value: fmt.Sprintf("%d", resp.Stats.LineCount), Label: l.StatByLineCount, Color: "amber", TestID: "service-stat-lines"},
				{Icon: "icon-plus-circle", Value: fmt.Sprintf("%d", resp.Stats.RecentlyAddedCnt), Label: l.StatRecentlyAdded, Color: "navy", TestID: "service-stat-recent"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID:        "by-line",
					Title:     l.WidgetByLine,
					Type:      "chart",
					ChartKind: "bar",
					ChartData: lineChart,
					Span:      2,
				},
				{
					ID:    "top-revenue",
					Title: l.WidgetTopRevenue,
					Type:  "custom",
					Span:  2,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.ListURL},
					},
					Custom: renderTopRevenueTable(topRev, l),
				},
				{
					ID:    "recent",
					Title: l.WidgetRecent,
					Type:  "list",
					Span:  1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.ListURL},
					},
					ListItems: recentItems,
					EmptyState: &types.EmptyStateData{
						Icon:  "icon-briefcase",
						Title: l.EmptyRecentTitle,
						Desc:  l.EmptyRecentDesc,
					},
				},
			},
		}

		dashboardURL := deps.Routes.DashboardURL
		if dashboardURL == "" {
			dashboardURL = centymo.ServiceDashboardURL
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "service",
				ActiveSubNav: "dashboard",
				HeaderTitle:  l.Title,
				HeaderIcon:   "icon-briefcase",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "service-dashboard-content",
			Dashboard:       dash,
		}
		_ = dashboardURL // reserved for future HTMX refresh wiring

		return view.OK("service-dashboard", pageData)
	})
}

// renderTopRevenueTable renders the top-revenue services table.
// Until the line-item revenue join is wired, this displays rank + name.
func renderTopRevenueTable(rows []TopRevenueRow, l centymo.ServiceDashboardLabels) template.HTML {
	if len(rows) == 0 {
		return template.HTML(fmt.Sprintf(
			`<div class="empty-state" data-testid="service-dashboard-top-empty"><p>%s</p></div>`,
			template.HTMLEscapeString(l.EmptyTopRevenue),
		))
	}
	var buf bytes.Buffer
	buf.WriteString(`<table class="data-table" id="service-dashboard-top-table"><thead><tr><th>`)
	buf.WriteString(template.HTMLEscapeString(l.ColRank))
	buf.WriteString(`</th><th>`)
	buf.WriteString(template.HTMLEscapeString(l.ColService))
	buf.WriteString(`</th></tr></thead><tbody>`)
	for i, r := range rows {
		buf.WriteString(fmt.Sprintf(`<tr data-testid="service-table-row-%d"><td>%d</td><td>`, i, i+1))
		buf.WriteString(template.HTMLEscapeString(r.ProductName))
		buf.WriteString(`</td></tr>`)
	}
	buf.WriteString(`</tbody></table>`)
	return template.HTML(buf.String())
}

// LineRowsToChart projects []LineRow onto parallel label/value slices used
// by the chart-bar widget. Exposed as a helper for orchestrator wrappers.
func LineRowsToChart(rows []LineRow) (labels []string, values []float64) {
	for _, r := range rows {
		labels = append(labels, r.LineID)
		values = append(values, float64(r.Count))
	}
	return labels, values
}
