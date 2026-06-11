// Package form holds the page data types for the
// expense-recognition-run queue page (Surface B).
//
// Mirror of packages/centymo-golang/views/revenue_run/queue/form/form.go.
// Plan A 20260517-expense-run Phase 4 / Surface B.
package form

import (
	"github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// PageData is the full data context passed to the
// expense-recognition-run-queue template.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	// AsOfDate is the YYYY-MM-DD date currently selected for the queue filter.
	AsOfDate string
	// AsOfDateMax is the maximum selectable date (today, YYYY-MM-DD).
	AsOfDateMax string
	// RefreshURL is the HTMX target for the date picker to hit.
	RefreshURL   string
	Labels       expense_recognition_run.QueueLabels
	CommonLabels pyeza.CommonLabels
}

// QueueRow is the view-layer representation of one supplier row in the queue.
type QueueRow struct {
	SupplierID               string
	SupplierName             string
	SupplierDetailURL        string
	DrawerURL                string // resolved per-supplier Surface-A drawer URL
	SubscriptionCount        int
	AdvanceDisbursementCount int
	PendingPeriods           int
	TotalAmount              int64
	Currency                 string
	// MultiCurrency is true when the supplier has candidates in more than one currency.
	MultiCurrency bool
	// ErrorMessage is set when the fan-out goroutine for this supplier returned an error
	// or panicked. The row still renders but shows an error chip instead of amounts.
	ErrorMessage string
}
