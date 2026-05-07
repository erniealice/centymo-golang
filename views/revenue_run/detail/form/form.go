// Package form holds the data types shared between the revenue-run detail page
// view and its templates.
package form

import (
	centymo "github.com/erniealice/centymo-golang"
	rrshared "github.com/erniealice/centymo-golang/views/revenue_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// PageData is the full data context passed to the revenue-run-detail template.
type PageData struct {
	types.PageData
	ContentTemplate string

	// Run holds the view-typed run row.
	Run rrshared.RevenueRunRow

	// Attempts holds all attempt rows for the run.
	Attempts []rrshared.RevenueRunAttemptRow

	// IsPossiblyInterrupted is true when Status=pending AND initiated_at is stale.
	IsPossiblyInterrupted bool

	// ActiveTab is the currently active tab key.
	ActiveTab string

	// TabItems is the slice of tab buttons rendered by {{template "tabs" ...}}.
	TabItems []pyeza.TabItem

	// Labels is the revenue-run label bundle.
	Labels centymo.RevenueRunLabels

	// SelectionsTable is the TableConfig for the Selections tab.
	SelectionsTable *types.TableConfig

	// ResultsTable is the TableConfig for the Results tab.
	ResultsTable *types.TableConfig

	// InvoicesTable is the TableConfig for the Invoices tab.
	InvoicesTable *types.TableConfig

	// AttachmentTable is the TableConfig for the Attachments tab.
	AttachmentTable *types.TableConfig
}
