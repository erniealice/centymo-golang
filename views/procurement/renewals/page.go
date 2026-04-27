// Package renewals renders the Renewal Calendar — contracts approaching
// date_time_end within the next 90 days, sorted by urgency (days until expiry).
package renewals

import (
	"context"
	"log"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

const renewalHorizonDays = 90

// RenewalRow holds display data for one expiring contract.
type RenewalRow struct {
	ContractID   string
	SupplierName string
	ContractRef  string
	DateTimeEnd  string
	DaysUntil    int
	Status       string
}

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.ProcurementRoutes
	Labels       centymo.ProcurementLabels
	CommonLabels pyeza.CommonLabels

	// nil-safe: view renders empty state when not provided
	ListSupplierContracts func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
}

// PageData holds the data for the renewal calendar page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          centymo.ProcurementLabels
	Routes          centymo.ProcurementRoutes
	Rows            []RenewalRow
	Empty           bool
}

// NewView creates the renewal calendar view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.RenewalsTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "procurement",
				ActiveSubNav: "renewals",
				HeaderTitle:  deps.Labels.RenewalsTitle,
				HeaderIcon:   "icon-calendar",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "procurement-renewal-calendar-content",
			Labels:          deps.Labels,
			Routes:          deps.Routes,
			Empty:           true,
		}

		if deps.ListSupplierContracts == nil {
			return view.OK("procurement-renewal-calendar", pageData)
		}

		resp, err := deps.ListSupplierContracts(ctx, &suppliercontractpb.ListSupplierContractsRequest{})
		if err != nil {
			log.Printf("renewals: ListSupplierContracts: %v", err)
			return view.OK("procurement-renewal-calendar", pageData)
		}

		now := time.Now().UTC()
		horizon := now.AddDate(0, 0, renewalHorizonDays)

		var rows []RenewalRow
		for _, c := range resp.GetData() {
			dateEnd := c.GetDateTimeEnd()
			if dateEnd == "" {
				continue // open-ended contract — not in renewal calendar
			}
			endDate, parseErr := time.Parse("2006-01-02", dateEnd)
			if parseErr != nil {
				continue
			}
			if endDate.Before(now) || endDate.After(horizon) {
				continue // already expired or beyond horizon
			}
			rows = append(rows, RenewalRow{
				ContractID:   c.GetId(),
				SupplierName: c.GetSupplierId(), // P3c enriches with supplier name via join
				ContractRef:  c.GetId(),
				DateTimeEnd:  dateEnd,
				DaysUntil:    int(endDate.Sub(now).Hours() / 24),
				Status:       c.GetStatus().String(),
			})
		}

		// Sort ascending by days until expiry (most urgent first)
		sortRenewalAsc(rows)

		pageData.Rows = rows
		pageData.Empty = len(rows) == 0

		return view.OK("procurement-renewal-calendar", pageData)
	})
}

func sortRenewalAsc(rows []RenewalRow) {
	for i := 1; i < len(rows); i++ {
		key := rows[i]
		j := i - 1
		for j >= 0 && rows[j].DaysUntil > key.DaysUntil {
			rows[j+1] = rows[j]
			j--
		}
		rows[j+1] = key
	}
}
