// Package list implements the revenue-run history list page (Surface D).
// Mirror of packages/centymo-golang/views/revenue/list/page.go.
package list

import (
	"context"
	"fmt"
	"log"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	rrshared "github.com/erniealice/centymo-golang/views/revenue_run/shared"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ListViewDeps holds view dependencies for the list page.
type ListViewDeps struct {
	Routes          centymo.RevenueRunRoutes
	Labels          centymo.RevenueRunLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
	ListRevenueRuns func(ctx context.Context, scope rrshared.ListRevenueRunsScope) ([]rrshared.RevenueRunRow, string, error)
}

// PageData is the full data context passed to the revenue-run-list template.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the full-page revenue-run list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "pending"
		}

		columns := revenueRunColumns(deps.Labels)
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
			ContentTemplate: "revenue-run-list-content",
			Table:           tableConfig,
		}

		return view.OK("revenue-run-list", pageData)
	})
}

// NewTableView returns only the table-card HTML (used as HTMX refresh target).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "pending"
		}

		columns := revenueRunColumns(deps.Labels)
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

// buildTableConfig fetches revenue-run data and builds the table configuration.
func buildTableConfig(
	ctx context.Context,
	deps *ListViewDeps,
	columns []types.TableColumn,
	status string,
	p tableparams.TableQueryParams,
) (*types.TableConfig, error) {
	if deps.ListRevenueRuns == nil {
		log.Printf("revenue-run list: ListRevenueRuns callback is nil — returning empty table")
	}

	var rows []rrshared.RevenueRunRow
	var nextCursor string
	if deps.ListRevenueRuns != nil {
		var err error
		rows, nextCursor, err = deps.ListRevenueRuns(ctx, rrshared.ListRevenueRunsScope{
			Status: status,
			Limit:  int32(p.PageSize),
		})
		if err != nil {
			log.Printf("Failed to list revenue runs: %v", err)
			return nil, fmt.Errorf("failed to load revenue runs: %w", err)
		}
	}
	if rows == nil {
		rows = []rrshared.RevenueRunRow{}
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
		ID:                   "revenue-run-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 tableRows,
		ShowSearch:           false, // cursor pagination doesn't combine with search
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

func revenueRunColumns(l centymo.RevenueRunLabels) []types.TableColumn {
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

func buildTableRows(rows []rrshared.RevenueRunRow, l centymo.RevenueRunLabels, routes centymo.RevenueRunRoutes) []types.TableRow {
	tableRows := make([]types.TableRow, 0, len(rows))
	for _, r := range rows {
		detailURL := route.ResolveURL(routes.DetailURL, "id", r.ID)

		// Status cell: if IsStalePending, override the badge label and variant
		// to show "Possibly interrupted" — this is the most accurate signal.
		// A separate second badge inside one cell isn't natively supported by
		// pyeza's table-cell types without Type:"html". We prefer the single
		// consolidated badge over injecting raw HTML, and surface the
		// "possibly interrupted" text as the primary status signal when stale.
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

				// Scope cell: show kind label + scope name when available
		scopeDisplay := scopeKindLabel(l, r.ScopeKind)
		if r.ScopeLabel != "" {
			scopeDisplay = scopeKindLabel(l, r.ScopeKind) + ": " + r.ScopeLabel
		}

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.ViewRun, Action: "view", Href: detailURL},
		}

		tableRows = append(tableRows, types.TableRow{
			ID:   r.ID,
			Href: detailURL,
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

func statusBadge(l centymo.RevenueRunLabels, status string) (label, variant string) {
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

func scopeKindLabel(l centymo.RevenueRunLabels, kind string) string {
	switch kind {
	case "subscription":
		return l.ScopeKind.Subscription
	case "client":
		return l.ScopeKind.Client
	case "workspace":
		return l.ScopeKind.Workspace
	default:
		return kind
	}
}

func statusPageTitle(l centymo.RevenueRunLabels, status string) string {
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

func statusEmptyTitle(l centymo.RevenueRunLabels, status string) string {
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

func statusEmptyMessage(l centymo.RevenueRunLabels, status string) string {
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
