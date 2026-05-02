// Package dashboard renders the Procurement Operations dashboard.
//
// Phase 1 refactor (2026-05-02): wired onto the pyeza "dashboard" block.
// Same aggregate data as before — pending-approval count, expiring contracts,
// top variance alerts, recurrence drafts — projected into typed Stats /
// Widgets / QuickActions on DashboardData. No new aggregate methods.
package dashboard

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	procurementrequestpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// VarianceRow is a summarised contract row for the top-5 variance widget.
type VarianceRow struct {
	ContractID      string
	SupplierName    string
	CommittedAmount int64
	BilledAmount    int64
	// UtilizationPct is billed/committed * 100, truncated to integer.
	UtilizationPct int
}

// ExpiringRow is a summarised contract row for the expiring-soon widget.
type ExpiringRow struct {
	ContractID   string
	SupplierName string
	DateTimeEnd  string
	DaysUntil    int
}

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.ProcurementRoutes
	Labels       centymo.ProcurementLabels
	CommonLabels pyeza.CommonLabels

	// nil-safe: widget shows 0 when not provided
	ListSupplierContracts   func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
	ListProcurementRequests func(ctx context.Context, req *procurementrequestpb.ListProcurementRequestsRequest) (*procurementrequestpb.ListProcurementRequestsResponse, error)
}

// PageData is what the procurement dashboard template receives.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the procurement dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels

		now := time.Now().UTC()
		horizon30 := now.AddDate(0, 0, 30)

		// Widget aggregates — preserved from the pre-refactor view.
		pendingApprovalCount := 0
		var expiring []ExpiringRow
		var varianceRows []VarianceRow

		if deps.ListProcurementRequests != nil {
			pendingStatus := procurementrequestpb.ProcurementRequestStatus_PROCUREMENT_REQUEST_STATUS_PENDING_APPROVAL
			resp, err := deps.ListProcurementRequests(ctx, &procurementrequestpb.ListProcurementRequestsRequest{
				Status: &pendingStatus,
			})
			if err != nil {
				log.Printf("procurement dashboard: ListProcurementRequests: %v", err)
			} else {
				pendingApprovalCount = len(resp.GetData())
			}
		}

		if deps.ListSupplierContracts != nil {
			resp, err := deps.ListSupplierContracts(ctx, &suppliercontractpb.ListSupplierContractsRequest{})
			if err != nil {
				log.Printf("procurement dashboard: ListSupplierContracts: %v", err)
			} else {
				for _, c := range resp.GetData() {
					if dateEnd := c.GetDateTimeEnd(); dateEnd != "" {
						endDate, parseErr := time.Parse("2006-01-02", dateEnd)
						if parseErr == nil && endDate.After(now) && !endDate.After(horizon30) {
							expiring = append(expiring, ExpiringRow{
								ContractID:   c.GetId(),
								SupplierName: c.GetSupplierId(),
								DateTimeEnd:  dateEnd,
								DaysUntil:    int(endDate.Sub(now).Hours() / 24),
							})
						}
					}

					committed := c.GetCommittedAmount()
					if committed > 0 {
						billed := c.GetBilledAmount()
						pct := int(billed * 100 / committed)
						varianceRows = append(varianceRows, VarianceRow{
							ContractID:      c.GetId(),
							SupplierName:    c.GetSupplierId(),
							CommittedAmount: committed,
							BilledAmount:    billed,
							UtilizationPct:  pct,
						})
					}
				}

				sortVarianceDesc(varianceRows)
				if len(varianceRows) > 5 {
					varianceRows = varianceRows[:5]
				}
			}
		}

		// Recurrence drafts always 0 until P5 ships the recurrence engine.
		recurrenceDraftsCount := 0

		// Build the dashboard input dict.
		dash := types.DashboardData{
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.PendingApprovalsTitle, Href: deps.Routes.DashboardURL, Variant: "primary", TestID: "procurement-action-new-request"},
				{Icon: "icon-calendar", Label: l.RenewalsTitle, Href: deps.Routes.RenewalCalendarURL, TestID: "procurement-action-renewals"},
				{Icon: "icon-alert-triangle", Label: l.VarianceTitle, Href: deps.Routes.VarianceURL, TestID: "procurement-action-variance"},
				{Icon: "icon-bar-chart-2", Label: l.UtilizationTitle, Href: deps.Routes.UtilizationURL, TestID: "procurement-action-utilization"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-inbox", Value: fmt.Sprintf("%d", pendingApprovalCount), Label: l.PendingApprovalsTitle, Color: "amber", TestID: "procurement-stat-pending"},
				{Icon: "icon-calendar", Value: fmt.Sprintf("%d", len(expiring)), Label: l.ExpiringTitle, Color: "terracotta", TestID: "procurement-stat-expiring"},
				{Icon: "icon-alert-triangle", Value: fmt.Sprintf("%d", len(varianceRows)), Label: l.VarianceTitle, Color: "navy", TestID: "procurement-stat-variance"},
				{Icon: "icon-refresh-cw", Value: fmt.Sprintf("%d", recurrenceDraftsCount), Label: l.RecurrenceTitle, Color: "sage", TestID: "procurement-stat-recurrence"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID:    "expiring",
					Title: l.ExpiringTitle,
					Type:  "custom",
					Span:  2,
					HeaderActions: []types.QuickAction{
						{Label: l.RenewalsTitle, Href: deps.Routes.RenewalCalendarURL},
					},
					Custom: renderExpiringTable(expiring, l),
				},
				{
					ID:    "variance",
					Title: l.VarianceTitle,
					Type:  "custom",
					Span:  2,
					HeaderActions: []types.QuickAction{
						{Label: l.BudgetPressureLabel, Href: deps.Routes.VarianceURL},
					},
					Custom: renderVarianceTable(varianceRows, l),
				},
				{
					ID:    "recurrence",
					Title: l.RecurrenceTitle,
					Type:  "custom",
					Span:  2,
					HeaderActions: []types.QuickAction{
						{Label: l.RecurrenceTitle, Href: deps.Routes.RecurrenceDraftsURL},
					},
					Custom: renderRecurrencePanel(recurrenceDraftsCount, l),
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.DashboardTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "procurement",
				ActiveSubNav: "dashboard",
				HeaderTitle:  l.DashboardTitle,
				HeaderIcon:   "icon-file-text",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "procurement-dashboard-content",
			Dashboard:       dash,
		}

		return view.OK("procurement-dashboard", pageData)
	})
}

// sortVarianceDesc is a simple insertion sort (N≤200 contracts typical).
func sortVarianceDesc(rows []VarianceRow) {
	for i := 1; i < len(rows); i++ {
		key := rows[i]
		j := i - 1
		for j >= 0 && rows[j].UtilizationPct < key.UtilizationPct {
			rows[j+1] = rows[j]
			j--
		}
		rows[j+1] = key
	}
}

// renderExpiringTable renders the expiring contracts table as raw HTML so it
// can ride inside a "custom" widget. Same markup as the pre-refactor template.
func renderExpiringTable(rows []ExpiringRow, l centymo.ProcurementLabels) template.HTML {
	if len(rows) == 0 {
		return template.HTML(fmt.Sprintf(
			`<div class="empty-state" data-testid="procurement-dashboard-expiring-empty"><p>%s</p></div>`,
			template.HTMLEscapeString(l.EmptyRenewals),
		))
	}
	var buf bytes.Buffer
	buf.WriteString(`<table class="data-table" id="procurement-dashboard-expiring-table"><thead><tr>`)
	buf.WriteString(`<th>Supplier</th><th>`)
	buf.WriteString(template.HTMLEscapeString(l.DaysUntilExpiry))
	buf.WriteString(`</th><th>End Date</th></tr></thead><tbody>`)
	for _, r := range rows {
		buf.WriteString(`<tr data-testid="expiring-contract-row"><td>`)
		buf.WriteString(template.HTMLEscapeString(r.SupplierName))
		buf.WriteString(`</td><td>`)
		buf.WriteString(fmt.Sprintf("%d", r.DaysUntil))
		buf.WriteString(`</td><td>`)
		buf.WriteString(template.HTMLEscapeString(r.DateTimeEnd))
		buf.WriteString(`</td></tr>`)
	}
	buf.WriteString(`</tbody></table>`)
	return template.HTML(buf.String())
}

// renderVarianceTable renders the top-variance contracts table.
func renderVarianceTable(rows []VarianceRow, l centymo.ProcurementLabels) template.HTML {
	if len(rows) == 0 {
		return template.HTML(fmt.Sprintf(
			`<div class="empty-state" data-testid="procurement-dashboard-variance-empty"><p>%s</p></div>`,
			template.HTMLEscapeString(l.EmptyVariance),
		))
	}
	var buf bytes.Buffer
	buf.WriteString(`<table class="data-table" id="procurement-dashboard-variance-table"><thead><tr>`)
	buf.WriteString(`<th>Supplier</th><th>`)
	buf.WriteString(template.HTMLEscapeString(l.UtilizationPercent))
	buf.WriteString(`</th></tr></thead><tbody>`)
	for _, r := range rows {
		buf.WriteString(`<tr data-testid="variance-alert-row"><td>`)
		buf.WriteString(template.HTMLEscapeString(r.SupplierName))
		buf.WriteString(`</td><td>`)
		buf.WriteString(fmt.Sprintf("%d%%", r.UtilizationPct))
		buf.WriteString(`</td></tr>`)
	}
	buf.WriteString(`</tbody></table>`)
	return template.HTML(buf.String())
}

// renderRecurrencePanel renders the recurrence-drafts summary panel.
func renderRecurrencePanel(count int, l centymo.ProcurementLabels) template.HTML {
	if count <= 0 {
		return template.HTML(fmt.Sprintf(
			`<div class="empty-state" data-testid="procurement-dashboard-recurrence-empty"><p>%s</p></div>`,
			template.HTMLEscapeString(l.EmptyRecurrence),
		))
	}
	return template.HTML(fmt.Sprintf(
		`<div class="centymo-activity-item" data-testid="procurement-dashboard-recurrence-count"><span>%d %s</span></div>`,
		count,
		template.HTMLEscapeString(l.EmptyRecurrence),
	))
}
