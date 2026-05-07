// Package queue implements the revenue-run workspace queue page (Surface B).
// Phase 7 of the 20260506-subscription-invoice-run plan.
//
// The queue lists every client that has at least one pending billing period as
// of the selected AsOfDate. For each visible client the view fans out a goroutine
// (bounded by REVENUE_RUN_QUEUE_FANOUT, default 16) to call
// ListRevenueRunCandidates. Per-client panics are recovered; a failed client
// renders an error chip on its row rather than failing the whole page.
package queue

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	queueform "github.com/erniealice/centymo-golang/views/revenue_run/queue/form"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	"golang.org/x/sync/errgroup"
)

// ClientRecord is a minimal client row used for queue population.
type ClientRecord struct {
	ID   string
	Name string
}

// CandidateSummary holds the aggregated candidate data for one client.
type CandidateSummary struct {
	// CandidatesBySubscription groups candidates by subscription ID.
	// Key = subscription ID, value = count of pending periods.
	SubscriptionIDs   map[string]struct{}
	TotalByCurrency   map[string]int64
	PeriodCount       int
}

// RevenueRunCandidateInput is the minimal shape the queue needs from a candidate.
// Populated by the ListRevenueRunCandidates callback shim in block.go.
type RevenueRunCandidateInput struct {
	SubscriptionID string
	Currency       string
	Amount         int64
	Eligible       bool
}

// QueueViewDeps holds all dependencies for the queue page views.
type QueueViewDeps struct {
	Routes       centymo.RevenueRunRoutes
	Labels       centymo.RevenueRunLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ClientDetailURLTemplate is the path template for the client detail page
	// (e.g. "/app/clients/detail/{id}"). Optional — rows are not linked when empty.
	ClientDetailURLTemplate string

	// ClientDrawerURLTemplate is the path template for the Surface-A per-client
	// revenue-run drawer (e.g. "/action/client/revenue-run/{id}").
	// Provided via WithClientRevenueRunDrawerURL BlockOption. Optional.
	ClientDrawerURLTemplate string

	// ListClients returns all clients visible to the current workspace user.
	// The returned slice is paginated; cursor-based pagination is used.
	ListClients func(ctx context.Context, cursor string) ([]ClientRecord, string, error)

	// ListRevenueRunCandidates returns pending billing periods for the given client
	// and as-of-date. Called once per client row in a bounded fan-out goroutine.
	ListRevenueRunCandidates func(ctx context.Context, clientID, asOfDate string) ([]RevenueRunCandidateInput, error)
}

// NewView creates the full-page revenue-run queue view.
func NewView(deps *QueueViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		asOfDate, asOfDateMax := resolveAsOfDate(viewCtx.Request.URL.Query().Get("as_of_date"))
		cursor := viewCtx.Request.URL.Query().Get("cursor")

		tableConfig, rows, err := buildTableConfig(ctx, deps, asOfDate, cursor)
		if err != nil {
			return view.Error(err)
		}

		l := deps.Labels.Queue
		pageData := &queueform.PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "queue",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-list",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "revenue-run-queue-content",
			Table:           tableConfig,
			AsOfDate:        asOfDate,
			AsOfDateMax:     asOfDateMax,
			RefreshURL:      deps.Routes.QueueTableURL,
			Labels:          l,
			CommonLabels:    deps.CommonLabels,
		}
		_ = rows

		return view.OK("revenue-run-queue", pageData)
	})
}

// NewTableView returns only the table-card HTML (used as HTMX refresh target).
func NewTableView(deps *QueueViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		asOfDate, _ := resolveAsOfDate(viewCtx.Request.URL.Query().Get("as_of_date"))
		cursor := viewCtx.Request.URL.Query().Get("cursor")

		tableConfig, _, err := buildTableConfig(ctx, deps, asOfDate, cursor)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches client + candidate data and assembles a TableConfig.
func buildTableConfig(
	ctx context.Context,
	deps *QueueViewDeps,
	asOfDate, cursor string,
) (*types.TableConfig, []queueform.QueueRow, error) {
	l := deps.Labels
	lq := l.Queue
	perms := view.GetUserPermissions(ctx)

	// 1. Fetch the visible client page.
	var clients []ClientRecord
	var nextCursor string
	if deps.ListClients != nil {
		var err error
		clients, nextCursor, err = deps.ListClients(ctx, cursor)
		if err != nil {
			log.Printf("revenue-run queue: failed to list clients: %v", err)
			return nil, nil, fmt.Errorf("failed to load client list: %w", err)
		}
	}

	// 2. Fan-out: one goroutine per client, bounded by REVENUE_RUN_QUEUE_FANOUT.
	fanout := fanoutLimit()
	type clientResult struct {
		idx     int
		summary CandidateSummary
		errMsg  string
	}
	results := make([]clientResult, len(clients))

	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(fanout)

	var mu sync.Mutex
	for i, c := range clients {
		idx := i
		clientID := c.ID
		eg.Go(func() (retErr error) {
			defer func() {
				if r := recover(); r != nil {
					msg := fmt.Sprintf("panic: %v", r)
					log.Printf("revenue-run queue: panic for client %s: %s", clientID, msg)
					mu.Lock()
					results[idx] = clientResult{idx: idx, errMsg: msg}
					mu.Unlock()
				}
			}()

			var summary CandidateSummary
			var rowErr string
			if deps.ListRevenueRunCandidates != nil {
				candidates, err := deps.ListRevenueRunCandidates(egCtx, clientID, asOfDate)
				if err != nil {
					log.Printf("revenue-run queue: ListRevenueRunCandidates error for client %s: %v", clientID, err)
					rowErr = err.Error()
				} else {
					summary = aggregateCandidates(candidates)
				}
			}

			mu.Lock()
			results[idx] = clientResult{idx: idx, summary: summary, errMsg: rowErr}
			mu.Unlock()
			return nil
		})
	}
	// Discard errgroup error — individual client errors are captured per-row.
	_ = eg.Wait()

	// 3. Build QueueRow list, filtering out clients with zero pending periods.
	queueRows := make([]queueform.QueueRow, 0, len(clients))
	for i, c := range clients {
		res := results[i]
		if res.errMsg == "" && res.summary.PeriodCount == 0 {
			// No pending periods — skip row.
			continue
		}

		drawerURL := ""
		if deps.ClientDrawerURLTemplate != "" {
			drawerURL = route.ResolveURL(deps.ClientDrawerURLTemplate, "id", c.ID)
		}
		clientDetailURL := ""
		if deps.ClientDetailURLTemplate != "" {
			clientDetailURL = route.ResolveURL(deps.ClientDetailURLTemplate, "id", c.ID)
		}

		var totalAmount int64
		currency := ""
		multiCurrency := false
		if res.errMsg == "" {
			totalAmount, currency, multiCurrency = flattenCurrencyTotals(res.summary.TotalByCurrency)
		}

		queueRows = append(queueRows, queueform.QueueRow{
			ClientID:          c.ID,
			ClientName:        c.Name,
			ClientDetailURL:   clientDetailURL,
			DrawerURL:         drawerURL,
			SubscriptionCount: len(res.summary.SubscriptionIDs),
			PeriodCount:       res.summary.PeriodCount,
			TotalAmount:       totalAmount,
			Currency:          currency,
			MultiCurrency:     multiCurrency,
			ErrorMessage:      res.errMsg,
		})
	}

	// 4. Convert to table rows.
	columns := queueColumns(l)
	tableRows := buildTableRows(queueRows, l, perms)
	types.ApplyColumnStyles(columns, tableRows)

	// 5. Build server pagination.
	sp := &types.ServerPagination{
		Enabled:       true,
		Mode:          "cursor",
		PaginationURL: deps.Routes.QueueTableURL,
	}
	if nextCursor != "" {
		sp.NextCursor = nextCursor
	}
	sp.BuildDisplay()

	// 6. Bulk actions.
	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:      "run-selected",
			Label:    lq.Bulk.RunSelected,
			Icon:     "icon-zap",
			Variant:  "primary",
			Endpoint: deps.Routes.SubmitBatchURL,
			ExtraParamsJSON: `{"selection_mode":"selected"}`,
		},
		{
			Key:      "run-all-matching",
			Label:    lq.Bulk.RunAllMatching,
			Icon:     "icon-zap",
			Variant:  "warning",
			Endpoint: deps.Routes.SubmitBatchURL,
			ExtraParamsJSON: `{"selection_mode":"all_matching"}`,
		},
	}

	tableConfig := &types.TableConfig{
		ID:          "revenue-run-queue-table",
		RefreshURL:  deps.Routes.QueueTableURL,
		Columns:     columns,
		Rows:        tableRows,
		ShowSearch:  false,
		ShowActions: true,
		ShowFilters: false,
		ShowSort:    false,
		ShowColumns: false,
		ShowExport:  false,
		ShowDensity: true,
		ShowEntries: false,
		Labels:      deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   lq.Empty.Title,
			Message: lq.Empty.Message,
		},
		BulkActions: &bulkCfg,
		ServerPagination: sp,
	}

	// The primary action (run for all matching) is covered by bulk actions.
	// PrimaryAction intentionally nil — the AsOfDate picker is rendered above
	// the table in the page template, not as a button.

	if !perms.Can("revenue", "create") {
		// Disable all bulk actions when the user lacks revenue:create.
		tableConfig.BulkActions = nil
	}

	types.ApplyTableSettings(tableConfig)
	return tableConfig, queueRows, nil
}

func queueColumns(l centymo.RevenueRunLabels) []types.TableColumn {
	lc := l.Queue.Columns
	return []types.TableColumn{
		{Key: "client", Label: lc.Client, WidthClass: "col-9xl"},
		{Key: "subscriptions", Label: lc.Subscriptions, WidthClass: "col-3xl", Align: "right", NoSort: true, NoFilter: true},
		{Key: "pending_periods", Label: lc.PendingPeriods, WidthClass: "col-3xl", Align: "right", NoSort: true, NoFilter: true},
		{Key: "total", Label: lc.Total, WidthClass: "col-4xl", Align: "right", NoSort: true, NoFilter: true},
		{Key: "currency", Label: lc.Currency, WidthClass: "col-2xl", NoSort: true, NoFilter: true},
	}
}

func buildTableRows(
	rows []queueform.QueueRow,
	l centymo.RevenueRunLabels,
	perms *types.UserPermissions,
) []types.TableRow {
	tableRows := make([]types.TableRow, 0, len(rows))
	lq := l.Queue

	for _, r := range rows {
		// Client name cell — linked to client detail when URL is set.
		var clientCell types.TableCell
		if r.ClientDetailURL != "" {
			clientCell = types.TableCell{Type: "link", Value: r.ClientName, Href: r.ClientDetailURL}
		} else {
			clientCell = types.TableCell{Type: "text", Value: r.ClientName}
		}

		// Currency cell — badge variant=info; warning badge when multi-currency.
		var currencyCell types.TableCell
		if r.ErrorMessage != "" {
			currencyCell = types.TableCell{Type: "badge", Value: l.Errors.InvalidSelection, Variant: "error"}
		} else if r.MultiCurrency {
			currencyCell = types.TableCell{Type: "badge", Value: r.Currency, Variant: "warning"}
		} else {
			currencyCell = types.TableCell{Type: "badge", Value: r.Currency, Variant: "info"}
		}

		// Total cell — centavos via MoneyCell (centMode=true).
		var totalCell types.TableCell
		if r.ErrorMessage != "" {
			totalCell = types.TableCell{Type: "text", Value: "—"}
		} else {
			totalCell = types.MoneyCell(float64(r.TotalAmount), r.Currency, true)
		}

		// Per-row [Run] action — loads Surface A drawer via HTMX GET into
		// #sheetContent, then opens the sheet. Using HxGet (not Href) avoids
		// the hx-boost navigation path which would replace #main-content instead
		// of the sheet portal.
		actions := []types.TableAction{}
		if r.DrawerURL != "" {
			actions = append(actions, types.TableAction{
				Type:            "run",
				Label:           lq.Columns.Run,
				Action:          "run",
				HxGet:           r.DrawerURL,
				HxTarget:        "#sheetContent",
				HxSwap:          "innerHTML",
				OnClick:         "lf.Sheet.open()",
				Disabled:        !perms.Can("revenue", "create"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		}

		tableRows = append(tableRows, types.TableRow{
			ID:   r.ClientID,
			Href: r.ClientDetailURL,
			Cells: []types.TableCell{
				clientCell,
				{Type: "text", Value: strconv.Itoa(r.SubscriptionCount), Align: "right"},
				{Type: "text", Value: strconv.Itoa(r.PeriodCount), Align: "right"},
				totalCell,
				currencyCell,
			},
			Actions: actions,
		})
	}
	return tableRows
}

// aggregateCandidates accumulates subscription IDs, period counts, and totals by
// currency from a flat candidate list.
func aggregateCandidates(candidates []RevenueRunCandidateInput) CandidateSummary {
	s := CandidateSummary{
		SubscriptionIDs: make(map[string]struct{}),
		TotalByCurrency: make(map[string]int64),
	}
	for _, c := range candidates {
		if !c.Eligible {
			continue
		}
		s.SubscriptionIDs[c.SubscriptionID] = struct{}{}
		s.TotalByCurrency[c.Currency] += c.Amount
		s.PeriodCount++
	}
	return s
}

// flattenCurrencyTotals returns the primary currency + total, and whether there
// are multiple currencies. When multi-currency, returns the first currency key
// and the sum of all amounts as the "total" (which will be displayed alongside
// a warning badge on the row).
func flattenCurrencyTotals(totals map[string]int64) (total int64, currency string, multi bool) {
	if len(totals) == 0 {
		return 0, "", false
	}
	if len(totals) == 1 {
		for cur, amt := range totals {
			return amt, cur, false
		}
	}
	// Multi-currency: sum everything, pick first key (alphabetically stable via map).
	// In practice, this is rare and the UI will show a warning badge.
	for cur, amt := range totals {
		total += amt
		if currency == "" || cur < currency {
			currency = cur
		}
	}
	return total, currency, true
}

// resolveAsOfDate returns the asOfDate to use and the max allowed date (today).
// Falls back to today when input is empty or invalid.
func resolveAsOfDate(input string) (asOfDate, maxDate string) {
	today := time.Now().UTC().Format("2006-01-02")
	if input == "" {
		return today, today
	}
	if _, err := time.Parse("2006-01-02", input); err != nil {
		return today, today
	}
	return input, today
}

// fanoutLimit reads REVENUE_RUN_QUEUE_FANOUT env and returns the concurrency
// cap. Falls back to 16 when the variable is absent or unparseable.
func fanoutLimit() int {
	if v := os.Getenv("REVENUE_RUN_QUEUE_FANOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 16
}
