// Package purchase_dashboard renders the Purchases dashboard
// (expenditure_type=purchase surface). Phase 5 — replaces the previous
// list-as-dashboard wiring.
//
// Workspace-scoped via the GetPageData callback (orchestrator pulls
// workspace_id from request context inside the wrapper).
package purchase_dashboard

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// TopSupplierRow mirrors the shape returned by the espyna adapter, copied
// here so the view package stays free of espyna internal imports.
type TopSupplierRow struct {
	SupplierID   string
	SupplierName string
	Total        int64 // centavos
}

// Stats holds tile values (centavos for monetary, count for non-).
type Stats struct {
	OpenCount        int64
	AwaitingCount    int64
	SpentMTD         int64
	TopSupplierName  string
	TopSupplierTotal int64
}

// Request is the input to the GetPageData callback.
type Request struct {
	Now time.Time
}

// Response is the projection the view consumes.
type Response struct {
	Stats        Stats
	MonthLabels  []string
	MonthValues  []float64 // centavos
	TopSuppliers []TopSupplierRow
	Recent       []*expenditurepb.Expenditure
}

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.ExpenditureRoutes
	Labels       centymo.ExpenditureLabels
	CommonLabels pyeza.CommonLabels

	// GetPageData is nil-safe; orchestrator wraps the espyna expenditure
	// dashboard use case (kind="purchase", workspace_id from ctx).
	GetPageData func(ctx context.Context, req *Request) (*Response, error)
}

// PageData is what the dashboard template receives.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the purchase dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.PurchaseDashboard
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

		// Spend-per-month bar chart.
		monthChart := &types.ChartData{
			Labels: resp.MonthLabels,
			Series: []types.ChartSeries{{
				Name:   l.WidgetMonthly,
				Values: resp.MonthValues,
				Color:  "terracotta",
			}},
			Currency: "PHP",
		}
		if len(monthChart.Labels) == 0 {
			monthChart.Labels = []string{"-"}
			monthChart.Series[0].Values = []float64{0}
		}
		monthChart.AutoScale()

		// Recent purchases activity list.
		recentItems := make([]types.ActivityItem, 0, len(resp.Recent))
		for i, e := range resp.Recent {
			title := e.GetReferenceNumber()
			if title == "" {
				title = e.GetName()
			}
			if title == "" {
				title = l.NewPurchase
			}
			amount := fmt.Sprintf("₱%s", formatCentavos(e.GetTotalAmount()))
			desc := amount
			if e.GetStatus() != "" {
				desc = fmt.Sprintf("%s — %s", amount, e.GetStatus())
			}
			recentItems = append(recentItems, types.ActivityItem{
				IconName:    "icon-shopping-bag",
				IconVariant: "client",
				Title:       title,
				Description: desc,
				Time:        e.GetExpenditureDateString(),
				TestID:      fmt.Sprintf("purchase-list-item-%d", i),
			})
		}

		topSupplierStat := l.EmptySuppliers
		if resp.Stats.TopSupplierName != "" {
			topSupplierStat = resp.Stats.TopSupplierName
		}

		dash := types.DashboardData{
			Title:    l.Title,
			Icon:     "icon-shopping-bag",
			Subtitle: l.Subtitle,
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickNew, Href: deps.Routes.AddURL, Variant: "primary", TestID: "purchase-action-new"},
				{Icon: "icon-package", Label: l.QuickReceive, Href: deps.Routes.PurchaseListURL, TestID: "purchase-action-receive"},
				{Icon: "icon-file-text", Label: l.QuickMatch, Href: deps.Routes.PurchaseListURL, TestID: "purchase-action-match"},
				{Icon: "icon-users", Label: l.QuickSuppliers, Href: deps.Routes.PurchaseListURL, TestID: "purchase-action-suppliers"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-shopping-bag", Value: fmt.Sprintf("%d", resp.Stats.OpenCount), Label: l.StatOpenPOs, Color: "terracotta", TestID: "purchase-stat-open"},
				{Icon: "icon-package", Value: fmt.Sprintf("%d", resp.Stats.AwaitingCount), Label: l.StatAwaiting, Color: "amber", TestID: "purchase-stat-awaiting"},
				{Icon: "icon-dollar-sign", Value: formatPesoSummary(resp.Stats.SpentMTD), Label: l.StatSpentMTD, Color: "navy", TestID: "purchase-stat-mtd"},
				{Icon: "icon-trending-up", Value: topSupplierStat, Label: l.StatTopSupplier, Color: "sage", TestID: "purchase-stat-supplier"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID:        "monthly",
					Title:     l.WidgetMonthly,
					Type:      "chart",
					ChartKind: "bar",
					ChartData: monthChart,
					Span:      2,
				},
				{
					ID:    "top-suppliers",
					Title: l.WidgetTopSupplier,
					Type:  "custom",
					Span:  2,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.PurchaseListURL},
					},
					Custom: renderTopSuppliersTable(resp.TopSuppliers, l),
				},
				{
					ID:    "recent",
					Title: l.WidgetRecent,
					Type:  "list",
					Span:  1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.PurchaseListURL},
					},
					ListItems: recentItems,
					EmptyState: &types.EmptyStateData{
						Icon:  "icon-shopping-bag",
						Title: l.EmptyRecentTitle,
						Desc:  l.EmptyRecentDesc,
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "purchase",
				ActiveSubNav: "dashboard",
				HeaderTitle:  l.Title,
				HeaderIcon:   "icon-shopping-bag",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "purchase-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("purchase-dashboard", pageData)
	})
}

// renderTopSuppliersTable renders the top-suppliers table inside a custom widget.
func renderTopSuppliersTable(rows []TopSupplierRow, l centymo.PurchaseDashboardLabels) template.HTML {
	if len(rows) == 0 {
		return template.HTML(fmt.Sprintf(
			`<div class="empty-state" data-testid="purchase-dashboard-suppliers-empty"><p>%s</p></div>`,
			template.HTMLEscapeString(l.EmptySuppliers),
		))
	}
	var buf bytes.Buffer
	buf.WriteString(`<table class="data-table" id="purchase-dashboard-suppliers-table"><thead><tr><th>`)
	buf.WriteString(template.HTMLEscapeString(l.ColSupplier))
	buf.WriteString(`</th><th>`)
	buf.WriteString(template.HTMLEscapeString(l.ColTotal))
	buf.WriteString(`</th></tr></thead><tbody>`)
	for i, r := range rows {
		buf.WriteString(fmt.Sprintf(`<tr data-testid="purchase-table-row-%d"><td>`, i))
		buf.WriteString(template.HTMLEscapeString(r.SupplierName))
		buf.WriteString(`</td><td>₱`)
		buf.WriteString(template.HTMLEscapeString(formatCentavos(r.Total)))
		buf.WriteString(`</td></tr>`)
	}
	buf.WriteString(`</tbody></table>`)
	return template.HTML(buf.String())
}

func formatCentavos(centavos int64) string {
	negative := centavos < 0
	if negative {
		centavos = -centavos
	}
	whole := centavos / 100
	cents := centavos % 100
	wholeStr := withThousandsSeparators(whole)
	out := fmt.Sprintf("%s.%02d", wholeStr, cents)
	if negative {
		out = "-" + out
	}
	return out
}

func formatPesoSummary(centavos int64) string {
	if centavos == 0 {
		return "₱0"
	}
	pesos := float64(centavos) / 100.0
	abs := pesos
	if abs < 0 {
		abs = -abs
	}
	switch {
	case abs >= 1_000_000:
		return fmt.Sprintf("₱%.1fM", pesos/1_000_000)
	case abs >= 10_000:
		return fmt.Sprintf("₱%.0fK", pesos/1_000)
	case abs >= 1_000:
		return fmt.Sprintf("₱%.1fK", pesos/1_000)
	default:
		return fmt.Sprintf("₱%.0f", pesos)
	}
}

func withThousandsSeparators(n int64) string {
	s := fmt.Sprintf("%d", n)
	if n < 0 || len(s) <= 3 {
		return s
	}
	out := make([]byte, 0, len(s)+len(s)/3)
	pre := len(s) % 3
	if pre > 0 {
		out = append(out, s[:pre]...)
		if len(s) > pre {
			out = append(out, ',')
		}
	}
	for i := pre; i < len(s); i += 3 {
		out = append(out, s[i:i+3]...)
		if i+3 < len(s) {
			out = append(out, ',')
		}
	}
	return string(out)
}
