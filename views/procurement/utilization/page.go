// Package utilization renders the Utilization Report — released_amount vs
// committed_amount per active contract, presented as a sortable table.
// Read-only; no drawer forms.
package utilization

import (
	"context"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// UtilizationRow holds per-contract utilization data.
// All amounts are in centavos (CLAUDE.md centavo convention).
type UtilizationRow struct {
	ContractID      string
	SupplierName    string
	ContractRef     string
	Status          string
	CommittedAmount int64 // centavos
	ReleasedAmount  int64 // centavos — sum of POs against this contract
	BilledAmount    int64 // centavos — sum of posted expenditures
	RemainingAmount int64 // centavos — committed - billed
	// ReleasedPct is released/committed * 100 (0 when committed == 0)
	ReleasedPct int
	// BilledPct is billed/committed * 100 (0 when committed == 0)
	BilledPct int
}

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.ProcurementRoutes
	Labels       centymo.ProcurementLabels
	CommonLabels pyeza.CommonLabels

	// nil-safe: view renders empty state when not provided
	ListSupplierContracts func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
}

// PageData holds the data for the utilization report page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          centymo.ProcurementLabels
	Routes          centymo.ProcurementRoutes
	Rows            []UtilizationRow
	Empty           bool

	// Totals row
	TotalCommitted int64
	TotalReleased  int64
	TotalBilled    int64
	TotalRemaining int64
}

// NewView creates the utilization report view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.UtilizationTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "procurement",
				ActiveSubNav: "utilization",
				HeaderTitle:  deps.Labels.UtilizationTitle,
				HeaderIcon:   "icon-bar-chart-2",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "procurement-utilization-report-content",
			Labels:          deps.Labels,
			Routes:          deps.Routes,
			Empty:           true,
		}

		if deps.ListSupplierContracts == nil {
			return view.OK("procurement-utilization-report", pageData)
		}

		activeStatus := suppliercontractpb.SupplierContractStatus_SUPPLIER_CONTRACT_STATUS_ACTIVE
		resp, err := deps.ListSupplierContracts(ctx, &suppliercontractpb.ListSupplierContractsRequest{
			Status: &activeStatus,
		})
		if err != nil {
			log.Printf("utilization: ListSupplierContracts: %v", err)
			return view.OK("procurement-utilization-report", pageData)
		}

		var rows []UtilizationRow
		var totalCommitted, totalReleased, totalBilled, totalRemaining int64

		for _, c := range resp.GetData() {
			committed := c.GetCommittedAmount()
			released := c.GetReleasedAmount()
			billed := c.GetBilledAmount()
			remaining := c.GetRemainingAmount()

			releasedPct := 0
			billedPct := 0
			if committed > 0 {
				releasedPct = int(released * 100 / committed)
				billedPct = int(billed * 100 / committed)
			}

			rows = append(rows, UtilizationRow{
				ContractID:      c.GetId(),
				SupplierName:    c.GetSupplierId(), // P3c enriches
				ContractRef:     c.GetId(),
				Status:          c.GetStatus().String(),
				CommittedAmount: committed,
				ReleasedAmount:  released,
				BilledAmount:    billed,
				RemainingAmount: remaining,
				ReleasedPct:     releasedPct,
				BilledPct:       billedPct,
			})

			totalCommitted += committed
			totalReleased += released
			totalBilled += billed
			totalRemaining += remaining
		}

		pageData.Rows = rows
		pageData.Empty = len(rows) == 0
		pageData.TotalCommitted = totalCommitted
		pageData.TotalReleased = totalReleased
		pageData.TotalBilled = totalBilled
		pageData.TotalRemaining = totalRemaining

		return view.OK("procurement-utilization-report", pageData)
	})
}
