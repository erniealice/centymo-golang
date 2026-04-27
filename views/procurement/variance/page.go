// Package variance renders the Variance Alerts view — contracts where
// billed_amount / committed_amount > 0.85 (budget pressure threshold).
// Read-only; no drawer forms.
package variance

import (
	"context"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// budgetPressureThreshold is the billed/committed ratio above which a contract
// is flagged as a variance alert (85%).
const budgetPressureThreshold = 0.85

// VarianceRow holds display data for one contract in budget pressure.
type VarianceRow struct {
	ContractID      string
	SupplierName    string
	ContractRef     string
	CommittedAmount int64 // centavos
	BilledAmount    int64 // centavos
	RemainingAmount int64 // centavos
	UtilizationPct  int   // billed/committed * 100
}

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.ProcurementRoutes
	Labels       centymo.ProcurementLabels
	CommonLabels pyeza.CommonLabels

	// nil-safe: view renders empty state when not provided
	ListSupplierContracts func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
}

// PageData holds the data for the variance alerts page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          centymo.ProcurementLabels
	Routes          centymo.ProcurementRoutes
	Rows            []VarianceRow
	Empty           bool
	ThresholdPct    int // display: 85
}

// NewView creates the variance alerts view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.VarianceTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "procurement",
				ActiveSubNav: "variance",
				HeaderTitle:  deps.Labels.VarianceTitle,
				HeaderIcon:   "icon-alert-triangle",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "procurement-variance-alerts-content",
			Labels:          deps.Labels,
			Routes:          deps.Routes,
			Empty:           true,
			ThresholdPct:    int(budgetPressureThreshold * 100),
		}

		if deps.ListSupplierContracts == nil {
			return view.OK("procurement-variance-alerts", pageData)
		}

		resp, err := deps.ListSupplierContracts(ctx, &suppliercontractpb.ListSupplierContractsRequest{})
		if err != nil {
			log.Printf("variance: ListSupplierContracts: %v", err)
			return view.OK("procurement-variance-alerts", pageData)
		}

		var rows []VarianceRow
		for _, c := range resp.GetData() {
			committed := c.GetCommittedAmount()
			if committed == 0 {
				continue // skip contracts with no committed amount
			}

			billed := c.GetBilledAmount()
			ratio := float64(billed) / float64(committed)
			if ratio <= budgetPressureThreshold {
				continue // under threshold — not a variance alert
			}

			rows = append(rows, VarianceRow{
				ContractID:      c.GetId(),
				SupplierName:    c.GetSupplierId(), // P3c enriches
				ContractRef:     c.GetId(),
				CommittedAmount: committed,
				BilledAmount:    billed,
				RemainingAmount: c.GetRemainingAmount(),
				UtilizationPct:  int(ratio * 100),
			})
		}

		// Sort descending by utilization pct (highest pressure first)
		sortVarianceDesc(rows)

		pageData.Rows = rows
		pageData.Empty = len(rows) == 0

		return view.OK("procurement-variance-alerts", pageData)
	})
}

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
