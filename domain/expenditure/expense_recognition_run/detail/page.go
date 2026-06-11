// Package detail implements the expense-recognition-run detail page (Surface D).
// Pattern mirrors packages/centymo-golang/views/revenue_run/detail/page.go.
// Plan A 20260517-expense-run Phase 4.
package detail

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run"

	detailform "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run/detail/form"
	errshared "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// DetailViewDeps holds view dependencies for the detail page.
type DetailViewDeps struct {
	Routes       expense_recognition_run.Routes
	Labels       expense_recognition_run.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ReadExpenseRecognitionRun fetches a run + all attempts by ID.
	ReadExpenseRecognitionRun func(ctx context.Context, id string) (*errshared.ExpenseRecognitionRunWithAttempts, error)

	// ListExpenseRecognitionsByRunID fetches expense_recognition rows whose run_id matches.
	// Used to populate the Recognitions tab on the detail page.
	ListExpenseRecognitionsByRunID func(ctx context.Context, runID string) ([]errshared.ExpenseRecognitionRow, error)

	// ListExpendituresByRunID fetches expenditure rows whose run_id matches.
	// Used to populate the Draft Bills tab on the detail page.
	ListExpendituresByRunID func(ctx context.Context, runID string) ([]errshared.ExpenditureRow, error)
}

// NewView creates the full-page expense-recognition-run detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense_recognition_run", "read") {
			return view.Forbidden("expense_recognition_run:read")
		}
		_ = perms
		id := viewCtx.Request.PathValue("id")

		if deps.ReadExpenseRecognitionRun == nil {
			return view.Error(fmt.Errorf("expense-recognition-run detail: ReadExpenseRecognitionRun callback not wired"))
		}
		runWithAttempts, err := deps.ReadExpenseRecognitionRun(ctx, id)
		if err != nil {
			log.Printf("Failed to read expense recognition run %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load run: %w", err))
		}
		if runWithAttempts == nil {
			log.Printf("Expense recognition run %s not found", id)
			return view.Error(fmt.Errorf("run not found"))
		}

		l := deps.Labels
		run := runWithAttempts.Run
		headerTitle := l.Detail.Title + " — " + run.ID

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "summary"
		}
		tabItems := buildTabItems(l, id, deps.Routes)

		pageData := &detailform.PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.Title,
				HeaderIcon:     "icon-zap",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:       "expense-recognition-run-detail-content",
			Run:                   run,
			Attempts:              runWithAttempts.Attempts,
			IsPossiblyInterrupted: run.IsStalePending,
			ActiveTab:             activeTab,
			TabItems:              tabItems,
			Labels:                l,
		}

		loadTabData(ctx, pageData, deps, runWithAttempts, activeTab)

		return view.OK("expense-recognition-run-detail", pageData)
	})
}

// NewTabAction creates a partial view that returns only the active tab content.
// Called via HTMX when the user clicks a tab button.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense_recognition_run", "read") {
			return view.Forbidden("expense_recognition_run:read")
		}
		_ = perms
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "summary"
		}

		if deps.ReadExpenseRecognitionRun == nil {
			return view.Error(fmt.Errorf("expense-recognition-run detail: ReadExpenseRecognitionRun callback not wired"))
		}
		runWithAttempts, err := deps.ReadExpenseRecognitionRun(ctx, id)
		if err != nil {
			log.Printf("Failed to read expense recognition run %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load run: %w", err))
		}
		if runWithAttempts == nil {
			return view.Error(fmt.Errorf("run not found"))
		}

		l := deps.Labels
		pageData := &detailform.PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Run:                   runWithAttempts.Run,
			Attempts:              runWithAttempts.Attempts,
			IsPossiblyInterrupted: runWithAttempts.Run.IsStalePending,
			ActiveTab:             tab,
			TabItems:              buildTabItems(l, id, deps.Routes),
			Labels:                l,
		}

		loadTabData(ctx, pageData, deps, runWithAttempts, tab)

		templateName := "expense-recognition-run-" + tab + "-tab"
		return view.OK(templateName, pageData)
	})
}

// loadTabData populates tab-specific fields on pageData.
func loadTabData(
	ctx context.Context,
	pageData *detailform.PageData,
	deps *DetailViewDeps,
	runWithAttempts *errshared.ExpenseRecognitionRunWithAttempts,
	tab string,
) {
	l := deps.Labels
	attempts := runWithAttempts.Attempts

	switch tab {
	case "summary":
		// All summary data is already in Run.

	case "selections":
		pageData.SelectionsTable = buildSelectionsTable(attempts, l, deps.TableLabels)

	case "results":
		pageData.ResultsTable = buildResultsTable(attempts, l, deps.TableLabels)

	case "bills":
		if deps.ListExpendituresByRunID != nil {
			rows, err := deps.ListExpendituresByRunID(ctx, runWithAttempts.Run.ID)
			if err != nil {
				log.Printf("Failed to load expenditures for run %s: %v", runWithAttempts.Run.ID, err)
				rows = []errshared.ExpenditureRow{}
			}
			pageData.BillsTable = buildBillsTable(rows, l, deps.TableLabels)
		} else {
			pageData.BillsTable = buildBillsTable(nil, l, deps.TableLabels)
		}

	case "recognitions":
		if deps.ListExpenseRecognitionsByRunID != nil {
			rows, err := deps.ListExpenseRecognitionsByRunID(ctx, runWithAttempts.Run.ID)
			if err != nil {
				log.Printf("Failed to load expense recognitions for run %s: %v", runWithAttempts.Run.ID, err)
				rows = []errshared.ExpenseRecognitionRow{}
			}
			pageData.RecognitionsTable = buildRecognitionsTable(rows, l, deps.TableLabels)
		} else {
			pageData.RecognitionsTable = buildRecognitionsTable(nil, l, deps.TableLabels)
		}

	case "audit-history":
		// Deferred — rendered as an info alert in the template.
	}
}

// buildTabItems constructs the tab bar items.
func buildTabItems(l expense_recognition_run.Labels, id string, routes expense_recognition_run.Routes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.DetailTabActionURL, "id", id, "tab", "")
	lt := l.Detail.Tabs
	return []pyeza.TabItem{
		{Key: "summary", Label: lt.Summary, Href: base + "?tab=summary", HxGet: action + "summary", Icon: "icon-info"},
		{Key: "selections", Label: lt.Selections, Href: base + "?tab=selections", HxGet: action + "selections", Icon: "icon-list"},
		{Key: "results", Label: lt.Results, Href: base + "?tab=results", HxGet: action + "results", Icon: "icon-check-circle"},
		{Key: "bills", Label: lt.Bills, Href: base + "?tab=bills", HxGet: action + "bills", Icon: "icon-receipt"},
		{Key: "recognitions", Label: lt.Recognitions, Href: base + "?tab=recognitions", HxGet: action + "recognitions", Icon: "icon-trending-down"},
		{Key: "audit-history", Label: lt.AuditHistory, Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

// buildSelectionsTable builds the TableConfig for the Selections tab.
// Shows all attempts (selections include created, skipped, errored).
func buildSelectionsTable(attempts []errshared.ExpenseRecognitionRunAttemptRow, l expense_recognition_run.Labels, tableLabels types.TableLabels) *types.TableConfig {
	ls := l.Detail.Selections
	columns := []types.TableColumn{
		{Key: "source", Label: ls.ColSource, NoSort: true, WidthClass: "col-3xl"},
		{Key: "supplier_subscription", Label: ls.ColSupplierSubscription, NoSort: true},
		{Key: "advance_disbursement", Label: ls.ColAdvanceDisbursement, NoSort: true},
		{Key: "period_start", Label: ls.ColPeriodStart, NoSort: true, WidthClass: "col-3xl"},
		{Key: "period_end", Label: ls.ColPeriodEnd, NoSort: true, WidthClass: "col-3xl"},
		{Key: "period_marker", Label: ls.ColPeriodMarker, NoSort: true, WidthClass: "col-3xl"},
	}

	rows := make([]types.TableRow, 0, len(attempts))
	for _, a := range attempts {
		rows = append(rows, types.TableRow{
			ID: a.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: sourceKindLabel(l, a.SourceKind)},
				{Type: "text", Value: displayOr(a.SupplierSubscriptionName, a.SupplierSubscriptionID)},
				{Type: "text", Value: displayOr(a.AdvanceDisbursementName, a.AdvanceDisbursementID)},
				{Type: "text", Value: a.PeriodStart},
				{Type: "text", Value: a.PeriodEnd},
				{Type: "text", Value: a.PeriodMarker},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "expense-run-selections-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   ls.EmptyTitle,
			Message: ls.EmptyMessage,
		},
	}
}

// buildResultsTable builds the TableConfig for the Results tab.
func buildResultsTable(attempts []errshared.ExpenseRecognitionRunAttemptRow, l expense_recognition_run.Labels, tableLabels types.TableLabels) *types.TableConfig {
	lr := l.Detail.Results
	columns := []types.TableColumn{
		{Key: "source", Label: lr.ColSource, NoSort: true, WidthClass: "col-3xl"},
		{Key: "supplier_subscription", Label: lr.ColSupplierSubscription, NoSort: true},
		{Key: "advance_disbursement", Label: lr.ColAdvanceDisbursement, NoSort: true},
		{Key: "period_start", Label: lr.ColPeriodStart, NoSort: true, WidthClass: "col-3xl"},
		{Key: "period_end", Label: lr.ColPeriodEnd, NoSort: true, WidthClass: "col-3xl"},
		{Key: "outcome", Label: lr.ColOutcome, NoSort: true, WidthClass: "col-3xl"},
		{Key: "error_code", Label: lr.ColErrorCode, NoSort: true, WidthClass: "col-4xl"},
	}

	rows := make([]types.TableRow, 0, len(attempts))
	for _, a := range attempts {
		outcomeLabel, outcomeVariant := outcomeCell(l, a.Outcome)
		rows = append(rows, types.TableRow{
			ID: a.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: sourceKindLabel(l, a.SourceKind)},
				{Type: "text", Value: displayOr(a.SupplierSubscriptionName, a.SupplierSubscriptionID)},
				{Type: "text", Value: displayOr(a.AdvanceDisbursementName, a.AdvanceDisbursementID)},
				{Type: "text", Value: a.PeriodStart},
				{Type: "text", Value: a.PeriodEnd},
				{Type: "badge", Value: outcomeLabel, Variant: outcomeVariant},
				{Type: "text", Value: a.ErrorCode},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "expense-run-results-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   lr.EmptyTitle,
			Message: lr.EmptyMessage,
		},
	}
}

// buildBillsTable builds the TableConfig for the Draft Bills tab.
func buildBillsTable(rows []errshared.ExpenditureRow, l expense_recognition_run.Labels, tableLabels types.TableLabels) *types.TableConfig {
	lb := l.Detail.Bills
	columns := []types.TableColumn{
		{Key: "reference", Label: lb.ColReference, WidthClass: "col-5xl"},
		{Key: "date", Label: lb.ColDate, WidthClass: "col-3xl"},
		{Key: "amount", Label: lb.ColAmount, WidthClass: "col-3xl", Align: "right"},
		{Key: "status", Label: lb.ColStatus, WidthClass: "col-3xl", NoSort: true},
	}

	tableRows := make([]types.TableRow, 0, len(rows))
	for _, r := range rows {
		statusLabel, statusVariant := genericStatusCell(r.Status)
		var actions []types.TableAction
		if r.DetailURL != "" {
			actions = append(actions, types.TableAction{
				Type:  "view",
				Label: l.Actions.ViewRun,
				Href:  r.DetailURL,
			})
		}
		tableRows = append(tableRows, types.TableRow{
			ID:   r.ID,
			Href: r.DetailURL,
			Cells: []types.TableCell{
				{Type: "text", Value: r.ReferenceNumber},
				types.DateTimeCell(r.ExpenditureDate, types.DateReadable),
				types.MoneyCell(float64(r.TotalAmount), r.Currency, true),
				{Type: "badge", Value: statusLabel, Variant: statusVariant},
			},
			Actions: actions,
		})
	}
	types.ApplyColumnStyles(columns, tableRows)

	return &types.TableConfig{
		ID:      "expense-run-bills-table",
		Columns: columns,
		Rows:    tableRows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   lb.EmptyTitle,
			Message: lb.EmptyMessage,
		},
	}
}

// buildRecognitionsTable builds the TableConfig for the Recognitions tab.
func buildRecognitionsTable(rows []errshared.ExpenseRecognitionRow, l expense_recognition_run.Labels, tableLabels types.TableLabels) *types.TableConfig {
	lr := l.Detail.Recognitions
	columns := []types.TableColumn{
		{Key: "reference", Label: lr.ColReference, WidthClass: "col-5xl"},
		{Key: "date", Label: lr.ColDate, WidthClass: "col-3xl"},
		{Key: "amount", Label: lr.ColAmount, WidthClass: "col-3xl", Align: "right"},
		{Key: "source_kind", Label: lr.ColSourceKind, WidthClass: "col-4xl", NoSort: true},
		{Key: "status", Label: lr.ColStatus, WidthClass: "col-3xl", NoSort: true},
	}

	tableRows := make([]types.TableRow, 0, len(rows))
	for _, r := range rows {
		statusLabel, statusVariant := genericStatusCell(r.Status)
		var actions []types.TableAction
		if r.DetailURL != "" {
			actions = append(actions, types.TableAction{
				Type:  "view",
				Label: l.Actions.ViewRun,
				Href:  r.DetailURL,
			})
		}
		tableRows = append(tableRows, types.TableRow{
			ID:   r.ID,
			Href: r.DetailURL,
			Cells: []types.TableCell{
				{Type: "text", Value: r.ReferenceNumber},
				types.DateTimeCell(r.RecognitionDate, types.DateReadable),
				types.MoneyCell(float64(r.TotalAmount), r.Currency, true),
				{Type: "text", Value: sourceKindLabel(l, r.SourceKind)},
				{Type: "badge", Value: statusLabel, Variant: statusVariant},
			},
			Actions: actions,
		})
	}
	types.ApplyColumnStyles(columns, tableRows)

	return &types.TableConfig{
		ID:      "expense-run-recognitions-table",
		Columns: columns,
		Rows:    tableRows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   lr.EmptyTitle,
			Message: lr.EmptyMessage,
		},
	}
}

func outcomeCell(l expense_recognition_run.Labels, outcome string) (label, variant string) {
	switch outcome {
	case "created":
		return l.AttemptOutcome.Created, "success"
	case "skipped":
		return l.AttemptOutcome.Skipped, "info"
	case "errored":
		return l.AttemptOutcome.Errored, "error"
	default:
		return outcome, "info"
	}
}

func sourceKindLabel(l expense_recognition_run.Labels, kind string) string {
	switch kind {
	case "subscription", "subscription_cycle":
		return l.SourceKind.SubscriptionCycle
	case "advance_disbursement":
		return l.SourceKind.AdvanceDisbursement
	default:
		return kind
	}
}

func genericStatusCell(status string) (label, variant string) {
	switch status {
	case "complete", "posted":
		return status, "success"
	case "draft":
		return status, "info"
	case "cancelled", "reversed":
		return status, "warning"
	default:
		return status, "info"
	}
}

func displayOr(primary, fallback string) string {
	if primary != "" {
		return primary
	}
	return fallback
}
