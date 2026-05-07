// Package form holds the page data types for the revenue-run queue page (Surface B).
package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// PageData is the full data context passed to the revenue-run-queue template.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	// AsOfDate is the YYYY-MM-DD date currently selected for the queue filter.
	AsOfDate string
	// AsOfDateMax is the maximum selectable date (today, YYYY-MM-DD).
	AsOfDateMax string
	// RefreshURL is the HTMX target for the date picker to hit.
	RefreshURL string
	Filters     QueueFilters
	Labels      centymo.RevenueRunQueueLabels
	CommonLabels pyeza.CommonLabels
}

// QueueFilters holds the active filter values for the queue page.
type QueueFilters struct {
	// ClientNamePrefix filters the client list to names starting with this prefix.
	ClientNamePrefix string
	// Currency filters to clients that have pending periods in this currency.
	Currency string
	// MinCandidates filters to clients with at least this many pending periods.
	MinCandidates int
}

// QueueRow is the view-layer representation of one client row in the queue.
type QueueRow struct {
	ClientID          string
	ClientName        string
	ClientDetailURL   string
	DrawerURL         string // resolved per-client Surface-A drawer URL
	SubscriptionCount int
	PeriodCount       int
	TotalAmount       int64
	Currency          string
	// MultiCurrency is true when the client has candidates in more than one currency.
	// The row renders a warning badge when true.
	MultiCurrency bool
	// ErrorMessage is set when the fan-out goroutine for this client returned an error
	// or panicked. The row still renders but shows an error chip instead of amounts.
	ErrorMessage string
}
