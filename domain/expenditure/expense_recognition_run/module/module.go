// Package expense_recognition_run wires the buying-side Expense Recognition
// Run view module — Plan A mirror of the shipped Revenue Run module.
//
// Surfaces wired by this package:
//   - Surface B: workspace queue page (queue/ sub-package)
//   - Surface D: run history list + detail (list/ + detail/ sub-packages)
//
// Surface A (per-supplier drawer) lives in entydad-golang/views/supplier;
// Surface C (per-supplier-subscription drawer) lives in
// centymo-golang/views/supplier_subscription. Both call into espyna consumer
// wrappers via the same callback shapes wired here.
//
// Reference: docs/plan/20260517-expense-run/plan.md §"Phase 4".
package expense_recognition_runmodule

import (
	"context"

	epkg "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run"
	errdetail "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run/detail"
	errdrawer "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run/drawer"
	errlist "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run/list"
	errqueue "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run/queue"
	errqueueaction "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run/queue/action"
	errshared "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// DrawerCandidateRow is re-exported so consumers can shape the
// ListExpenseRunCandidates callback's return type without importing the drawer
// sub-package directly.
type DrawerCandidateRow = errdrawer.CandidateRow

// DrawerListScope is the scope-input shape for the drawer's candidate listing.
type DrawerListScope = errdrawer.ListScope

// DrawerGenerateInput is the scope-input shape for the drawer's Generate call.
type DrawerGenerateInput = errdrawer.GenerateInput

// ---------------------------------------------------------------------------
// Re-export shared view-typed data shapes so block.go callers can reference
// them via the top-level expense_recognition_run package.
// ---------------------------------------------------------------------------

// ExpenseRecognitionRunRow is the view-layer representation of a single run.
type ExpenseRecognitionRunRow = errshared.ExpenseRecognitionRunRow

// ExpenseRecognitionRunWithAttempts bundles a run + its attempts.
type ExpenseRecognitionRunWithAttempts = errshared.ExpenseRecognitionRunWithAttempts

// ExpenseRecognitionRunAttemptRow is a view-layer representation of an attempt.
type ExpenseRecognitionRunAttemptRow = errshared.ExpenseRecognitionRunAttemptRow

// ExpenseRecognitionRow is a minimal recognition row for the Recognitions tab.
type ExpenseRecognitionRow = errshared.ExpenseRecognitionRow

// ExpenditureRow is a minimal draft-bill row for the Bills tab.
type ExpenditureRow = errshared.ExpenditureRow

// ListExpenseRecognitionRunsScope carries list-page filter params.
type ListExpenseRecognitionRunsScope = errshared.ListExpenseRecognitionRunsScope

// QueueSupplierRecord is a minimal supplier row for queue fan-out.
type QueueSupplierRecord = errshared.QueueSupplierRecord

// QueueCandidateInput is the per-candidate input for the queue fan-out.
type QueueCandidateInput = errshared.QueueCandidateInput

// BatchRunInput is the per-supplier input for the batch handler.
type BatchRunInput = errshared.BatchRunInput

// BatchRunOutput is the per-supplier output from the batch handler.
type BatchRunOutput = errshared.BatchRunOutput

// ---------------------------------------------------------------------------
// ModuleDeps — typed callbacks; no espyna/proto types cross this boundary.
// ---------------------------------------------------------------------------

// ModuleDeps holds all dependencies for the expense-recognition-run view module.
type ModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Surface D callbacks — run history list + detail.

	// ListExpenseRecognitionRuns returns a page of run rows matching the scope.
	ListExpenseRecognitionRuns func(ctx context.Context, scope ListExpenseRecognitionRunsScope) ([]ExpenseRecognitionRunRow, string, error)

	// ReadExpenseRecognitionRun fetches a single run plus all attempts.
	ReadExpenseRecognitionRun func(ctx context.Context, id string) (*ExpenseRecognitionRunWithAttempts, error)

	// ListExpenseRecognitionsByRunID returns expense-recognition rows for the Recognitions tab.
	ListExpenseRecognitionsByRunID func(ctx context.Context, runID string) ([]ExpenseRecognitionRow, error)

	// ListExpendituresByRunID returns expenditure rows for the Bills tab.
	ListExpendituresByRunID func(ctx context.Context, runID string) ([]ExpenditureRow, error)

	// Surface B callbacks — workspace queue page.

	// SupplierDetailURLTemplate is the path template for the supplier-detail
	// page (e.g. "/app/suppliers/detail/{id}"). Optional.
	SupplierDetailURLTemplate string

	// SupplierDrawerURLTemplate is the path template for the Surface-A per-supplier
	// drawer (e.g. "/action/supplier/expense-recognition-run/{id}").
	SupplierDrawerURLTemplate string

	// ListSuppliers returns all suppliers visible to the current workspace user.
	ListSuppliers func(ctx context.Context, cursor string) ([]QueueSupplierRecord, string, error)

	// ListExpenseRunCandidates returns pending candidates for one supplier.
	ListExpenseRunCandidates func(ctx context.Context, supplierID, asOfDate string) ([]QueueCandidateInput, error)

	// GenerateExpenseRun executes the expense run for a single supplier.
	GenerateExpenseRun func(ctx context.Context, in BatchRunInput) (*BatchRunOutput, error)

	// Surface A + C drawer callbacks.

	// ListExpenseRunCandidatesForDrawer returns the candidate list for the
	// per-supplier (Surface A) or per-supplier-subscription (Surface C) drawer.
	// The scope discriminator is on the ListScope arg.
	ListExpenseRunCandidatesForDrawer func(ctx context.Context, scope DrawerListScope) ([]DrawerCandidateRow, error)

	// GenerateExpenseRunForDrawer executes the run for a per-supplier or
	// per-supplier-subscription drawer submit.
	GenerateExpenseRunForDrawer func(ctx context.Context, in DrawerGenerateInput) (*BatchRunOutput, error)
}

// Module holds all constructed expense-recognition-run views.
type Module struct {
	routes epkg.Routes
	// Surface D.
	List      view.View
	Table     view.View
	Detail    view.View
	TabAction view.View
	// Surface B.
	Queue      view.View
	QueueTable view.View
	BatchRun   view.View
	// Surface A + C drawer.
	SupplierDrawer     view.View
	SubscriptionDrawer view.View
	GenerateAction     view.View
}

// NewModule constructs the expense-recognition-run module from the deps.
func NewModule(deps *ModuleDeps) *Module {
	listDeps := &errlist.ListViewDeps{
		Routes:                     deps.Routes,
		Labels:                     deps.Labels,
		CommonLabels:               deps.CommonLabels,
		TableLabels:                deps.TableLabels,
		ListExpenseRecognitionRuns: deps.ListExpenseRecognitionRuns,
	}
	detailDeps := &errdetail.DetailViewDeps{
		Routes:                         deps.Routes,
		Labels:                         deps.Labels,
		CommonLabels:                   deps.CommonLabels,
		TableLabels:                    deps.TableLabels,
		ReadExpenseRecognitionRun:      deps.ReadExpenseRecognitionRun,
		ListExpenseRecognitionsByRunID: deps.ListExpenseRecognitionsByRunID,
		ListExpendituresByRunID:        deps.ListExpendituresByRunID,
	}
	queueDeps := &errqueue.QueueViewDeps{
		Routes:                    deps.Routes,
		Labels:                    deps.Labels,
		CommonLabels:              deps.CommonLabels,
		TableLabels:               deps.TableLabels,
		SupplierDetailURLTemplate: deps.SupplierDetailURLTemplate,
		SupplierDrawerURLTemplate: deps.SupplierDrawerURLTemplate,
		ListSuppliers:             deps.ListSuppliers,
		ListExpenseRunCandidates:  deps.ListExpenseRunCandidates,
	}
	batchRunDeps := &errqueueaction.BatchRunDeps{
		Routes:             deps.Routes,
		Labels:             deps.Labels,
		GenerateExpenseRun: deps.GenerateExpenseRun,
	}
	drawerDeps := &errdrawer.Deps{
		Routes:                   deps.Routes,
		Labels:                   deps.Labels,
		CommonLabels:             deps.CommonLabels,
		TableLabels:              deps.TableLabels,
		ListExpenseRunCandidates: deps.ListExpenseRunCandidatesForDrawer,
		GenerateExpenseRun:       deps.GenerateExpenseRunForDrawer,
	}
	return &Module{
		routes:             deps.Routes,
		List:               errlist.NewView(listDeps),
		Table:              errlist.NewTableView(listDeps),
		Detail:             errdetail.NewView(detailDeps),
		TabAction:          errdetail.NewTabAction(detailDeps),
		Queue:              errqueue.NewView(queueDeps),
		QueueTable:         errqueue.NewTableView(queueDeps),
		BatchRun:           errqueueaction.NewBatchRunAction(batchRunDeps),
		SupplierDrawer:     errdrawer.NewSupplierDrawer(drawerDeps),
		SubscriptionDrawer: errdrawer.NewSubscriptionDrawer(drawerDeps),
		GenerateAction:     errdrawer.NewGenerateAction(drawerDeps),
	}
}

// RegisterRoutes registers all expense-recognition-run routes on the registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	// Surface D — run history list + detail.
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.ListTableURL, m.Table)
	r.POST(m.routes.ListTableURL, m.Table)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.DetailTabActionURL, m.TabAction)
	// Surface B — workspace queue page.
	r.GET(m.routes.QueueURL, m.Queue)
	r.GET(m.routes.QueueTableURL, m.QueueTable)
	r.POST(m.routes.QueueTableURL, m.QueueTable)
	r.POST(m.routes.SubmitBatchURL, m.BatchRun)
	// Surface A + C drawer.
	if m.routes.PerSupplierDrawerURL != "" && m.SupplierDrawer != nil {
		r.GET(m.routes.PerSupplierDrawerURL, m.SupplierDrawer)
	}
	if m.routes.PerSubscriptionDrawerURL != "" && m.SubscriptionDrawer != nil {
		r.GET(m.routes.PerSubscriptionDrawerURL, m.SubscriptionDrawer)
	}
	if m.routes.GenerateURL != "" && m.GenerateAction != nil {
		r.POST(m.routes.GenerateURL, m.GenerateAction)
	}
}
