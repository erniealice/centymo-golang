// Package form holds the data types shared between the
// expense-recognition-run detail page view and its templates.
//
// Mirror of packages/centymo-golang/views/revenue_run/detail/form/form.go.
// Plan A 20260517-expense-run Phase 4 / Surface D.
package form

import (
	"github.com/erniealice/centymo-golang/domain/expenditure"
	errshared "github.com/erniealice/centymo-golang/domain/expenditure/views/expense_recognition_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// PageData is the full data context passed to the
// expense-recognition-run-detail template.
type PageData struct {
	types.PageData
	ContentTemplate string

	// Run holds the view-typed run row.
	Run errshared.ExpenseRecognitionRunRow

	// Attempts holds all attempt rows for the run.
	Attempts []errshared.ExpenseRecognitionRunAttemptRow

	// IsPossiblyInterrupted is true when Status=pending AND initiated_at is stale.
	IsPossiblyInterrupted bool

	// ActiveTab is the currently active tab key.
	ActiveTab string

	// TabItems is the slice of tab buttons rendered by {{template "tabs" ...}}.
	TabItems []pyeza.TabItem

	// Labels is the expense-recognition-run label bundle.
	Labels expenditure.ExpenseRecognitionRunLabels

	// SelectionsTable is the TableConfig for the Selections tab.
	SelectionsTable *types.TableConfig

	// ResultsTable is the TableConfig for the Results tab.
	ResultsTable *types.TableConfig

	// BillsTable is the TableConfig for the Draft Bills tab.
	BillsTable *types.TableConfig

	// RecognitionsTable is the TableConfig for the Recognitions tab.
	RecognitionsTable *types.TableConfig
}
