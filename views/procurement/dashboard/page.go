// Package dashboard renders the Procurement Operations dashboard.
// Widgets: pending-approval request count, expiring contracts (30 days),
// top-5 variance alerts, recurrence drafts pending review.
package dashboard

import (
	"context"
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

// PageData holds dashboard widget data.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          centymo.ProcurementLabels
	Routes          centymo.ProcurementRoutes

	// Widget: pending approval count
	PendingApprovalCount int

	// Widget: contracts expiring in 30 days
	ExpiringContracts []ExpiringRow

	// Widget: top-5 contracts by billed/committed utilization
	TopVarianceAlerts []VarianceRow

	// Widget: recurrence drafts count (0 until P5 ships)
	RecurrenceDraftsCount int
}

// NewView creates the procurement dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.DashboardTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "procurement",
				ActiveSubNav: "dashboard",
				HeaderTitle:  deps.Labels.DashboardTitle,
				HeaderIcon:   "icon-file-text",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "procurement-dashboard-content",
			Labels:          deps.Labels,
			Routes:          deps.Routes,
		}

		now := time.Now().UTC()
		horizon30 := now.AddDate(0, 0, 30)

		// Widget: pending-approval procurement request count
		if deps.ListProcurementRequests != nil {
			pendingStatus := procurementrequestpb.ProcurementRequestStatus_PROCUREMENT_REQUEST_STATUS_PENDING_APPROVAL
			resp, err := deps.ListProcurementRequests(ctx, &procurementrequestpb.ListProcurementRequestsRequest{
				Status: &pendingStatus,
			})
			if err != nil {
				log.Printf("procurement dashboard: ListProcurementRequests: %v", err)
			} else {
				pageData.PendingApprovalCount = len(resp.GetData())
			}
		}

		// Widget: expiring contracts (next 30 days) + top-5 variance
		if deps.ListSupplierContracts != nil {
			resp, err := deps.ListSupplierContracts(ctx, &suppliercontractpb.ListSupplierContractsRequest{})
			if err != nil {
				log.Printf("procurement dashboard: ListSupplierContracts: %v", err)
			} else {
				var varianceRows []VarianceRow
				for _, c := range resp.GetData() {
					// Expiring widget: contracts with a finite end date within 30 days
					if dateEnd := c.GetDateTimeEnd(); dateEnd != "" {
						endDate, parseErr := time.Parse("2006-01-02", dateEnd)
						if parseErr == nil && endDate.After(now) && !endDate.After(horizon30) {
							pageData.ExpiringContracts = append(pageData.ExpiringContracts, ExpiringRow{
								ContractID:   c.GetId(),
								SupplierName: c.GetSupplierId(), // P3c will enrich with supplier name
								DateTimeEnd:  dateEnd,
								DaysUntil:    int(endDate.Sub(now).Hours() / 24),
							})
						}
					}

					// Variance: collect all contracts with committed > 0
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

				// Sort variance rows descending by utilization pct, take top 5
				sortVarianceDesc(varianceRows)
				if len(varianceRows) > 5 {
					varianceRows = varianceRows[:5]
				}
				pageData.TopVarianceAlerts = varianceRows
			}
		}

		// Widget: recurrence drafts (always 0 until P5 ships the recurrence engine)
		// When P5 lands, wire ListExpenditures here with status=draft + supplier_contract_id IS NOT NULL.
		pageData.RecurrenceDraftsCount = 0

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
