// Package dashboard renders the Cash (collection) dashboard.
//
// Phase 5 (2026-05-02): real aggregates wired via the GetPageData callback,
// which the orchestrator backs with the espyna
// internal/application/usecases/treasury/collection/dashboard use case. To
// keep this package free of espyna internal imports, the contract is
// expressed by locally-defined Stats/Response types — orchestrator adapts.
//
// Workspace-scoped via the GetPageData call (orchestrator pulls workspace_id
// from context inside the wrapper).
package dashboard

import (
	"context"
	"fmt"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Stats is the cash dashboard's KPI tile values (centavos).
type Stats struct {
	Pending           int64
	Overdue           int64
	CollectedToday    int64
	CollectedThisWeek int64
}

// Request is the input to the GetPageData callback.
type Request struct {
	Now time.Time
}

// Response is the projection the view consumes.
type Response struct {
	Stats Stats

	// Daily series (last 30 days) — Labels parallel to Values, centavos.
	DailyLabels []string
	DailyValues []float64

	// Payment-mode mix — Labels parallel to Values, centavos.
	ModeLabels []string
	ModeValues []float64

	Recent []*collectionpb.Collection
}

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.CollectionRoutes
	Labels       centymo.CollectionLabels
	CommonLabels pyeza.CommonLabels

	// GetPageData is nil-safe; when nil the dashboard renders zero values.
	// Orchestrator wraps the espyna treasury/collection/dashboard use case
	// (workspace_id pulled from request context inside the wrapper).
	GetPageData func(ctx context.Context, req *Request) (*Response, error)
}

// PageData is what the cash dashboard template receives.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the cash dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.Dashboard
		now := time.Now()

		// Load aggregates (nil-safe).
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

		// Daily trend chart-line.
		dailyChart := &types.ChartData{
			Labels: resp.DailyLabels,
			Series: []types.ChartSeries{{
				Name:   l.WidgetDailyTrend,
				Values: resp.DailyValues,
				Color:  "sage",
			}},
			Currency: "PHP",
		}
		if len(dailyChart.Labels) == 0 {
			dailyChart.Labels = []string{"-"}
			dailyChart.Series[0].Values = []float64{0}
		}
		dailyChart.AutoScale()

		// Payment-mode chart-pie.
		modeChart := &types.ChartData{
			Labels: resp.ModeLabels,
			Series: []types.ChartSeries{{
				Name:   l.WidgetByMode,
				Values: resp.ModeValues,
				Color:  "terracotta",
			}},
			Currency: "PHP",
		}
		if len(modeChart.Labels) == 0 {
			modeChart.Labels = []string{"-"}
			modeChart.Series[0].Values = []float64{0}
		}
		modeChart.AutoScale()

		// Recent collections list.
		recentItems := make([]types.ActivityItem, 0, len(resp.Recent))
		for i, c := range resp.Recent {
			title := c.GetReferenceNumber()
			if title == "" {
				title = c.GetName()
			}
			if title == "" {
				title = l.NewCollection
			}
			amount := fmt.Sprintf("₱%s", formatCentavos(c.GetAmount()))
			desc := amount
			if c.GetCollectionMethodId() != "" {
				desc = fmt.Sprintf("%s — %s", amount, c.GetCollectionMethodId())
			}
			recentItems = append(recentItems, types.ActivityItem{
				IconName:    "icon-dollar-sign",
				IconVariant: "client",
				Title:       title,
				Description: desc,
				Time:        c.GetPaymentDate(),
				TestID:      fmt.Sprintf("cash-list-item-%d", i),
			})
		}

		dash := types.DashboardData{
			Title:    l.Title,
			Icon:     "icon-dollar-sign",
			Subtitle: l.Subtitle,
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickRecord, Href: deps.Routes.AddURL, Variant: "primary", TestID: "cash-action-record"},
				{Icon: "icon-check-circle", Label: l.QuickReconcile, Href: deps.Routes.ListURL, TestID: "cash-action-reconcile"},
				{Icon: "icon-clock", Label: l.QuickAging, Href: deps.Routes.ListURL, TestID: "cash-action-aging"},
				{Icon: "icon-check", Label: l.QuickMarkCleared, Href: deps.Routes.ListURL, TestID: "cash-action-cleared"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-clock", Value: formatPesoSummary(resp.Stats.Pending), Label: l.StatPending, Color: "amber", TestID: "cash-stat-pending"},
				{Icon: "icon-alert-triangle", Value: formatPesoSummary(resp.Stats.Overdue), Label: l.StatOverdue, Color: "navy", TestID: "cash-stat-overdue"},
				{Icon: "icon-dollar-sign", Value: formatPesoSummary(resp.Stats.CollectedToday), Label: l.StatCollectedToday, Color: "terracotta", TestID: "cash-stat-today"},
				{Icon: "icon-trending-up", Value: formatPesoSummary(resp.Stats.CollectedThisWeek), Label: l.StatCollectedWeek, Color: "sage", TestID: "cash-stat-week"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID:        "daily-trend",
					Title:     l.WidgetDailyTrend,
					Type:      "chart",
					ChartKind: "line",
					ChartData: dailyChart,
					Span:      2,
				},
				{
					ID:        "by-mode",
					Title:     l.WidgetByMode,
					Type:      "chart",
					ChartKind: "pie",
					ChartData: modeChart,
					Span:      1,
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
						Icon:  "icon-dollar-sign",
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
				ActiveNav:    "cash",
				ActiveSubNav: "dashboard",
				HeaderTitle:  l.Title,
				HeaderIcon:   "icon-dollar-sign",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "cash-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("cash-dashboard", pageData)
	})
}

// formatCentavos renders centavos as "1,234.50" (no currency symbol).
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

// formatPesoSummary returns a compact peso amount like "₱4.8M" for stat cards.
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

// withThousandsSeparators inserts comma thousands separators.
func withThousandsSeparators(n int64) string {
	s := fmt.Sprintf("%d", n)
	if n < 0 {
		return s
	}
	if len(s) <= 3 {
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
