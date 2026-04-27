// Package recurrence_drafts renders the Recurrence Drafts queue —
// Expenditure rows with status='draft' AND supplier_contract_id IS NOT NULL.
// This is the AP team's review/approve queue fed by the recurrence engine.
//
// P5 deferred: the recurrence engine does not exist yet. This view will
// always return empty state until P5 ships. The view is wired now so the
// sidebar slot exists and lights up automatically when P5 lands.
package recurrence_drafts

import (
	"context"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// DraftRow holds display data for one recurrence draft expenditure.
type DraftRow struct {
	ExpenditureID     string
	ContractID        string
	SupplierName      string
	Amount            int64 // centavos (GetTotalAmount)
	DateCreatedString string
	Status            string
}

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.ProcurementRoutes
	Labels       centymo.ProcurementLabels
	CommonLabels pyeza.CommonLabels

	// nil-safe: view renders empty state when not provided (pre-P5)
	ListExpenditures func(ctx context.Context, req *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error)
}

// PageData holds the data for the recurrence drafts queue page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          centymo.ProcurementLabels
	Routes          centymo.ProcurementRoutes
	Rows            []DraftRow
	Empty           bool
	// P5Pending signals to the template that the recurrence engine is not yet
	// active, so the empty state can show an informative "coming soon" message
	// rather than a generic "no data" message.
	P5Pending bool
}

// NewView creates the recurrence drafts queue view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.RecurrenceTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "procurement",
				ActiveSubNav: "recurrence-drafts",
				HeaderTitle:  deps.Labels.RecurrenceTitle,
				HeaderIcon:   "icon-refresh-cw",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "procurement-recurrence-drafts-content",
			Labels:          deps.Labels,
			Routes:          deps.Routes,
			Empty:           true,
			P5Pending:       deps.ListExpenditures == nil,
		}

		if deps.ListExpenditures == nil {
			// Recurrence engine not wired — P5 deferred. Show empty state.
			return view.OK("procurement-recurrence-drafts", pageData)
		}

		resp, err := deps.ListExpenditures(ctx, &expenditurepb.ListExpendituresRequest{})
		if err != nil {
			log.Printf("recurrence_drafts: ListExpenditures: %v", err)
			return view.OK("procurement-recurrence-drafts", pageData)
		}

		// Filter: status=draft AND supplier_contract_id is set
		// (The ListExpenditures request may not yet support supplier_contract_id filter —
		// post-P5 this can be pushed down to the use case for efficiency.)
		var rows []DraftRow
		for _, e := range resp.GetData() {
			if e.GetStatus() != "draft" {
				continue
			}
			scID := e.GetSupplierContractId()
			if scID == "" {
				continue
			}
			rows = append(rows, DraftRow{
				ExpenditureID:     e.GetId(),
				ContractID:        scID,
				SupplierName:      e.GetSupplierId(), // P3c enriches
				Amount:            e.GetTotalAmount(),
				DateCreatedString: e.GetDateCreatedString(),
				Status:            e.GetStatus(),
			})
		}

		pageData.Rows = rows
		pageData.Empty = len(rows) == 0

		return view.OK("procurement-recurrence-drafts", pageData)
	})
}
