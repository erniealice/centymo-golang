// Package expense_dashboard renders the Expenses dashboard
// (expenditure_type=expense surface). Phase 5 — replaces the previous
// list-as-dashboard wiring.
//
// Workspace-scoped via the GetPageData callback (orchestrator pulls
// workspace_id from request context inside the wrapper).
package expense_dashboard

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"sort"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Stats holds tile values for the expense surface (centavos for monetary).
type Stats struct {
	PendingApprovalCount int64
	ApprovedMTD          int64
	ReimbursableMTD      int64
	CategoriesUsed       int64
}

// catRow is one row in the top-categories table widget.
type catRow struct {
	Category string
	Total    int64 // centavos
}

// Request is the input to the GetPageData callback.
type Request struct {
	Now time.Time
}

// Response is the projection the view consumes.
type Response struct {
	Stats          Stats
	CategoryLabels []string
	CategoryValues []float64 // centavos
	Recent         []*expenditurepb.Expenditure
}

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.ExpenditureRoutes
	Labels       centymo.ExpenditureLabels
	CommonLabels pyeza.CommonLabels

	// GetPageData is nil-safe; orchestrator wraps the espyna expenditure
	// dashboard use case (kind="expense", workspace_id from ctx).
	GetPageData func(ctx context.Context, req *Request) (*Response, error)
}

// PageData is what the dashboard template receives.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the expense dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.ExpenseDashboard
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

		// By-category bar chart.
		catChart := &types.ChartData{
			Labels: resp.CategoryLabels,
			Series: []types.ChartSeries{{
				Name:   l.WidgetByCategory,
				Values: resp.CategoryValues,
				Color:  "plum",
			}},
			Currency: "PHP",
		}
		if len(catChart.Labels) == 0 {
			catChart.Labels = []string{"-"}
			catChart.Series[0].Values = []float64{0}
		}
		catChart.AutoScale()

		// Top categories rows for table widget — sorted desc by amount.
		topCats := make([]catRow, 0, len(resp.CategoryLabels))
		for i := range resp.CategoryLabels {
			topCats = append(topCats, catRow{
				Category: resp.CategoryLabels[i],
				Total:    int64(resp.CategoryValues[i]),
			})
		}
		sort.Slice(topCats, func(i, j int) bool { return topCats[i].Total > topCats[j].Total })
		if len(topCats) > 5 {
			topCats = topCats[:5]
		}

		// Recent expenses.
		recentItems := make([]types.ActivityItem, 0, len(resp.Recent))
		for i, e := range resp.Recent {
			title := e.GetReferenceNumber()
			if title == "" {
				title = e.GetName()
			}
			if title == "" {
				title = l.NewExpense
			}
			amount := fmt.Sprintf("₱%s", formatCentavos(e.GetTotalAmount()))
			desc := amount
			if e.GetExpenditureCategoryId() != "" {
				desc = fmt.Sprintf("%s — %s", amount, e.GetExpenditureCategoryId())
			}
			recentItems = append(recentItems, types.ActivityItem{
				IconName:    "icon-file-text",
				IconVariant: "integration",
				Title:       title,
				Description: desc,
				Time:        e.GetExpenditureDateString(),
				TestID:      fmt.Sprintf("expense-list-item-%d", i),
			})
		}

		dash := types.DashboardData{
			Title:    l.Title,
			Icon:     "icon-file-text",
			Subtitle: l.Subtitle,
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickNew, Href: deps.Routes.AddURL, Variant: "primary", TestID: "expense-action-new"},
				{Icon: "icon-check-circle", Label: l.QuickApprove, Href: deps.Routes.ExpenseListURL, TestID: "expense-action-approve"},
				{Icon: "icon-credit-card", Label: l.QuickReimburse, Href: deps.Routes.ExpenseListURL, TestID: "expense-action-reimburse"},
				{Icon: "icon-settings", Label: l.QuickCategorySettings, Href: deps.Routes.ExpenseCategoryListURL, TestID: "expense-action-categories"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-clock", Value: fmt.Sprintf("%d", resp.Stats.PendingApprovalCount), Label: l.StatPendingApproval, Color: "amber", TestID: "expense-stat-pending"},
				{Icon: "icon-check-circle", Value: formatPesoSummary(resp.Stats.ApprovedMTD), Label: l.StatApprovedMTD, Color: "sage", TestID: "expense-stat-approved"},
				{Icon: "icon-credit-card", Value: formatPesoSummary(resp.Stats.ReimbursableMTD), Label: l.StatReimbursable, Color: "terracotta", TestID: "expense-stat-reimbursable"},
				{Icon: "icon-pie-chart", Value: fmt.Sprintf("%d", resp.Stats.CategoriesUsed), Label: l.StatCategoriesUsed, Color: "navy", TestID: "expense-stat-categories"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID:        "by-category",
					Title:     l.WidgetByCategory,
					Type:      "chart",
					ChartKind: "bar",
					ChartData: catChart,
					Span:      2,
				},
				{
					ID:    "top-categories",
					Title: l.WidgetTopCategory,
					Type:  "custom",
					Span:  2,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.ExpenseCategoryListURL},
					},
					Custom: renderTopCategoriesTable(topCats, l),
				},
				{
					ID:    "recent",
					Title: l.WidgetRecent,
					Type:  "list",
					Span:  1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.ExpenseListURL},
					},
					ListItems: recentItems,
					EmptyState: &types.EmptyStateData{
						Icon:  "icon-file-text",
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
				ActiveNav:    "expense",
				ActiveSubNav: "dashboard",
				HeaderTitle:  l.Title,
				HeaderIcon:   "icon-file-text",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "expense-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("expense-dashboard", pageData)
	})
}

// renderTopCategoriesTable renders top expense categories inside a custom widget.
func renderTopCategoriesTable(rows []catRow, l centymo.ExpenseDashboardLabels) template.HTML {
	if len(rows) == 0 {
		return template.HTML(fmt.Sprintf(
			`<div class="empty-state" data-testid="expense-dashboard-categories-empty"><p>%s</p></div>`,
			template.HTMLEscapeString(l.EmptyCategories),
		))
	}
	var buf bytes.Buffer
	buf.WriteString(`<table class="data-table" id="expense-dashboard-categories-table"><thead><tr><th>`)
	buf.WriteString(template.HTMLEscapeString(l.ColCategory))
	buf.WriteString(`</th><th>`)
	buf.WriteString(template.HTMLEscapeString(l.ColTotal))
	buf.WriteString(`</th></tr></thead><tbody>`)
	for i, r := range rows {
		buf.WriteString(fmt.Sprintf(`<tr data-testid="expense-table-row-%d"><td>`, i))
		buf.WriteString(template.HTMLEscapeString(r.Category))
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
