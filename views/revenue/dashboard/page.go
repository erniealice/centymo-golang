package dashboard

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	Labels       centymo.RevenueLabels
	Routes       centymo.RevenueRoutes
	CommonLabels pyeza.CommonLabels
}

// PageData is what the revenue dashboard template receives.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the revenue dashboard view.
//
// Phase 1 refactor (2026-05-02): wired onto the pyeza "dashboard" block.
// Aggregates remain placeholder until Phase 5 wires real Revenue/Invoice
// repository aggregate methods.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.Dashboard

		trend := &types.ChartData{
			Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
			Series: []types.ChartSeries{{
				Name:   l.Revenue,
				Values: []float64{42, 58, 71, 64, 88, 95, 102, 89, 110, 124, 118, 134},
				Color:  "terracotta",
			}},
		}
		trend.AutoScale()

		dash := types.DashboardData{
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickNewRevenue, Href: deps.Routes.AddURL, Variant: "primary", TestID: "revenue-action-new"},
				{Icon: "icon-list", Label: l.QuickViewAll, Href: deps.Routes.ListURL, TestID: "revenue-action-list"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-shopping-bag", Value: "156", Label: l.TotalRevenue, Trend: "+18%", TrendUp: true, Color: "terracotta", TestID: "revenue-stat-total"},
				{Icon: "icon-dollar-sign", Value: "₱284K", Label: l.Revenue, Trend: "+12%", TrendUp: true, Color: "sage", TestID: "revenue-stat-amount"},
				{Icon: "icon-check-circle", Value: "89", Label: l.Completed, Trend: "+7", TrendUp: true, Color: "navy", TestID: "revenue-stat-completed"},
				{Icon: "icon-clock", Value: "42", Label: l.Active, Trend: "+3", TrendUp: true, Color: "amber", TestID: "revenue-stat-active"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID: "trend", Title: l.RevenueTrend, Type: "chart", ChartKind: "line",
					ChartData: trend, Span: 2,
					HeaderActions: []types.QuickAction{
						{Label: l.Week, Href: "#"},
						{Label: l.Month, Href: "#"},
						{Label: l.Year, Href: "#", Variant: "primary"},
					},
				},
				{
					ID: "recent", Title: l.RecentRevenue, Type: "list", Span: 1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.ListURL},
					},
					ListItems: []types.ActivityItem{
						{IconName: "icon-shopping-bag", IconVariant: "client", Title: l.NewRevenueCreated, Description: "REF-2024-001 — ₱12,500", Time: "30m ago", TestID: "revenue-activity-new"},
						{IconName: "icon-check-circle", IconVariant: "quote", Title: l.RevenueCompleted, Description: "REF-2024-098 — ₱45,000", Time: "2h ago", TestID: "revenue-activity-completed"},
						{IconName: "icon-edit", IconVariant: "award", Title: l.RevenueUpdated, Description: "REF-2024-095 amount adjusted", Time: "4h ago", TestID: "revenue-activity-updated"},
						{IconName: "icon-x-circle", IconVariant: "integration", Title: l.RevenueCancelled, Description: "REF-2024-087 — ₱8,200", Time: "1d ago", TestID: "revenue-activity-cancelled"},
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "revenue",
				ActiveSubNav: "dashboard",
				HeaderTitle:  l.Title,
				HeaderIcon:   "icon-shopping-bag",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "revenue-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("revenue-dashboard", pageData)
	})
}
