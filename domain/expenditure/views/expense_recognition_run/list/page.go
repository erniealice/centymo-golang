// Package list implements the expense-recognition-run history list page
// (Surface D of the Plan A Expense Run epic).
//
// Mirror of packages/centymo-golang/views/revenue_run/list/page.go.
// Plan A 20260517-expense-run Phase 4.
package list

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/erniealice/centymo-golang/domain/expenditure"
	errshared "github.com/erniealice/centymo-golang/domain/expenditure/views/expense_recognition_run/shared"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ListViewDeps holds view dependencies for the list page.
type ListViewDeps struct {
	Routes                     expenditure.ExpenseRecognitionRunRoutes
	Labels                     expenditure.ExpenseRecognitionRunLabels
	CommonLabels               pyeza.CommonLabels
	TableLabels                types.TableLabels
	ListExpenseRecognitionRuns func(ctx context.Context, scope errshared.ListExpenseRecognitionRunsScope) ([]errshared.ExpenseRecognitionRunRow, string, error)
}

// PageData is the full data context passed to the expense-recognition-run-list
// template.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the full-page expense-recognition-run list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense_recognition_run", "list") {
			return view.Forbidden("expense_recognition_run:list")
		}
		_ = perms
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "pending"
		}

		columns := expenseRunColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(
			viewCtx.Request,
			types.SortableKeys(columns),
			types.FilterableKeys(columns),
			"initiated_at",
			"desc",
		)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, status, p)
		if err != nil {
			return view.Error(err)
		}

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(l, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(l, status),
				HeaderSubtitle: l.List.Subtitle,
				HeaderIcon:     "icon-zap",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "expense-recognition-run-list-content",
			Table:           tableConfig,
		}

		return view.OK("expense-recognition-run-list", pageData)
	})
}

// NewTableView returns only the table-card HTML (used as HTMX refresh target).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "pending"
		}

		columns := expenseRunColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(
			viewCtx.Request,
			types.SortableKeys(columns),
			types.FilterableKeys(columns),
			"initiated_at",
			"desc",
		)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, status, p)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches expense-run data and builds the table configuration.
func buildTableConfig(
	ctx context.Context,
	deps *ListViewDeps,
	columns []types.TableColumn,
	status string,
	p tableparams.TableQueryParams,
) (*types.TableConfig, error) {
	if deps.ListExpenseRecognitionRuns == nil {
		log.Printf("expense-recognition-run list: ListExpenseRecognitionRuns callback is nil — returning empty table")
	}

	var rows []errshared.ExpenseRecognitionRunRow
	var nextCursor string
	if deps.ListExpenseRecognitionRuns != nil {
		var err error
		rows, nextCursor, err = deps.ListExpenseRecognitionRuns(ctx, errshared.ListExpenseRecognitionRunsScope{
			Status: status,
			Limit:  int32(p.PageSize),
		})
		if err != nil {
			log.Printf("Failed to list expense recognition runs: %v", err)
			return nil, fmt.Errorf("failed to load expense recognition runs: %w", err)
		}
	}
	if rows == nil {
		rows = []errshared.ExpenseRecognitionRunRow{}
	}

	l := deps.Labels
	tableRows := buildTableRows(rows, l, deps.Routes)
	types.ApplyColumnStyles(columns, tableRows)

	refreshURL := route.ResolveURL(deps.Routes.ListTableURL, "status", status)

	sp := &types.ServerPagination{
		Enabled:       true,
		Mode:          "cursor",
		SortColumn:    p.SortColumn,
		SortDirection: p.SortDir,
		FiltersJSON:   p.FiltersRaw,
		PaginationURL: refreshURL,
	}
	if nextCursor != "" {
		sp.NextCursor = nextCursor
	}
	sp.BuildDisplay()

	tableConfig := &types.TableConfig{
		ID:                   "expense-run-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 tableRows,
		ShowSearch:           false,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "initiated_at",
		DefaultSortDirection: "desc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   statusEmptyTitle(l, status),
			Message: statusEmptyMessage(l, status),
		},
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func expenseRunColumns(l expenditure.ExpenseRecognitionRunLabels) []types.TableColumn {
	lc := l.List.Columns
	return []types.TableColumn{
		{Key: "id", Label: lc.ID, WidthClass: "col-5xl"},
		{Key: "scope", Label: lc.Scope, NoSort: true, NoFilter: true},
		{Key: "as_of_date", Label: lc.AsOfDate, WidthClass: "col-3xl"},
		{Key: "initiator", Label: lc.Initiator, WidthClass: "col-6xl", NoSort: true},
		{Key: "initiated_at", Label: lc.InitiatedAt, WidthClass: "col-4xl"},
		{Key: "status", Label: lc.Status, WidthClass: "col-3xl", NoSort: true},
		{Key: "created", Label: lc.Created, WidthClass: "col-md", Align: "right"},
		{Key: "skipped", Label: lc.Skipped, WidthClass: "col-md", Align: "right"},
		{Key: "errored", Label: lc.Errored, WidthClass: "col-md", Align: "right"},
	}
}

func buildTableRows(rows []errshared.ExpenseRecognitionRunRow, l expenditure.ExpenseRecognitionRunLabels, routes expenditure.ExpenseRecognitionRunRoutes) []types.TableRow {
	tableRows := make([]types.TableRow, 0, len(rows))
	for _, r := range rows {
		detailURL := route.ResolveURL(routes.DetailURL, "id", r.ID)

		var statusCell types.TableCell
		if r.IsStalePending {
			statusCell = types.TableCell{
				Type:    "badge",
				Value:   l.StatusBadges.PossiblyInterrupted,
				Variant: "warning",
			}
		} else {
			label, variant := statusBadge(l, r.Status)
			statusCell = types.TableCell{
				Type:    "badge",
				Value:   label,
				Variant: variant,
			}
		}

		scopeDisplay := scopeKindLabel(l, r.ScopeKind)
		if r.ScopeLabel != "" {
			scopeDisplay = scopeDisplay + ": " + r.ScopeLabel
		}

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.ViewRun, Action: "view", Href: detailURL},
		}

		tableRows = append(tableRows, types.TableRow{
			ID:   r.ID,
			Href: detailURL,
			DataAttrs: map[string]string{
				"testid": "expense-run-row-" + r.ID,
			},
			Cells: []types.TableCell{
				{Type: "text", Value: r.ID},
				{Type: "text", Value: scopeDisplay},
				{Type: "text", Value: r.AsOfDate},
				{Type: "text", Value: r.InitiatorName},
				types.DateTimeCell(r.InitiatedAt, types.DateTimeFull),
				statusCell,
				{Type: "text", Value: strconv.Itoa(int(r.CreatedCount)), Align: "right"},
				{Type: "text", Value: strconv.Itoa(int(r.SkippedCount)), Align: "right"},
				{Type: "text", Value: strconv.Itoa(int(r.ErroredCount)), Align: "right"},
			},
			Actions: actions,
		})
	}
	return tableRows
}

func statusBadge(l expenditure.ExpenseRecognitionRunLabels, status string) (label, variant string) {
	switch status {
	case "pending":
		return l.StatusBadges.Pending, "warning"
	case "complete":
		return l.StatusBadges.Complete, "success"
	case "failed":
		return l.StatusBadges.Failed, "error"
	default:
		return status, "info"
	}
}

func scopeKindLabel(l expenditure.ExpenseRecognitionRunLabels, kind string) string {
	switch kind {
	case "supplier":
		return l.ScopeKind.Supplier
	case "subscription":
		return l.ScopeKind.Subscription
	case "workspace":
		return l.ScopeKind.Workspace
	default:
		return kind
	}
}

func statusPageTitle(l expenditure.ExpenseRecognitionRunLabels, status string) string {
	switch status {
	case "pending":
		return l.List.Title + " — " + l.List.Filters.Pending
	case "complete":
		return l.List.Title + " — " + l.List.Filters.Complete
	case "failed":
		return l.List.Title + " — " + l.List.Filters.Failed
	default:
		return l.List.Title
	}
}

func statusEmptyTitle(l expenditure.ExpenseRecognitionRunLabels, status string) string {
	switch status {
	case "pending":
		return l.List.Empty.Pending.Title
	case "complete":
		return l.List.Empty.Complete.Title
	case "failed":
		return l.List.Empty.Failed.Title
	default:
		return l.List.Empty.Pending.Title
	}
}

func statusEmptyMessage(l expenditure.ExpenseRecognitionRunLabels, status string) string {
	switch status {
	case "pending":
		return l.List.Empty.Pending.Message
	case "complete":
		return l.List.Empty.Complete.Message
	case "failed":
		return l.List.Empty.Failed.Message
	default:
		return l.List.Empty.Pending.Message
	}
}
