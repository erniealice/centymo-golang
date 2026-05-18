// Package queue implements the expense-recognition-run workspace queue page
// (Surface B of the Plan A Expense Run epic).
//
// The queue lists every supplier that has at least one pending recognition
// candidate (subscription cycle or advance Disbursement tranche) as of the
// selected AsOfDate. For each visible supplier the view fans out a goroutine
// (bounded by EXPENSE_RUN_QUEUE_FANOUT, default 16) to call
// ListExpenseRunCandidates. Per-supplier panics are recovered; a failed supplier
// renders an error chip on its row rather than failing the whole page.
//
// Mirror of packages/centymo-golang/views/revenue_run/queue/page.go.
// Plan A 20260517-expense-run Phase 4.
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
	queueform "github.com/erniealice/centymo-golang/views/expense_recognition_run/queue/form"
	errshared "github.com/erniealice/centymo-golang/views/expense_recognition_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	"golang.org/x/sync/errgroup"
)

// CandidateSummary holds the aggregated candidate data for one supplier.
type CandidateSummary struct {
	SubscriptionIDs        map[string]struct{}
	AdvanceDisbursementIDs map[string]struct{}
	TotalByCurrency        map[string]int64
	PeriodCount            int
}

// QueueViewDeps holds all dependencies for the queue page views.
type QueueViewDeps struct {
	Routes       centymo.ExpenseRecognitionRunRoutes
	Labels       centymo.ExpenseRecognitionRunLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// SupplierDetailURLTemplate is the path template for the supplier detail page
	// (e.g. "/app/suppliers/detail/{id}"). Optional — rows are not linked when empty.
	SupplierDetailURLTemplate string

	// SupplierDrawerURLTemplate is the path template for the Surface-A per-supplier
	// expense-recognition-run drawer (e.g. "/action/supplier/expense-recognition-run/{id}").
	SupplierDrawerURLTemplate string

	// ListSuppliers returns all suppliers visible to the current workspace user.
	ListSuppliers func(ctx context.Context, cursor string) ([]errshared.QueueSupplierRecord, string, error)

	// ListExpenseRunCandidates returns pending recognition candidates for one supplier.
	// Called per-supplier in a bounded fan-out goroutine on the queue page.
	ListExpenseRunCandidates func(ctx context.Context, supplierID, asOfDate string) ([]errshared.QueueCandidateInput, error)
}

// NewView creates the full-page expense-recognition-run queue view.
func NewView(deps *QueueViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense_recognition_run", "list") {
			return view.Forbidden("expense_recognition_run:list")
		}
		asOfDate, asOfDateMax := resolveAsOfDate(ctx, viewCtx.Request.URL.Query().Get("as_of_date"))
		cursor := viewCtx.Request.URL.Query().Get("cursor")

		tableConfig, _, err := buildTableConfig(ctx, deps, asOfDate, cursor)
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
			ContentTemplate: "expense-recognition-run-queue-content",
			Table:           tableConfig,
			AsOfDate:        asOfDate,
			AsOfDateMax:     asOfDateMax,
			RefreshURL:      deps.Routes.QueueTableURL,
			Labels:          l,
			CommonLabels:    deps.CommonLabels,
		}

		return view.OK("expense-recognition-run-queue", pageData)
	})
}

// NewTableView returns only the table-card HTML (used as HTMX refresh target).
func NewTableView(deps *QueueViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		asOfDate, _ := resolveAsOfDate(ctx, viewCtx.Request.URL.Query().Get("as_of_date"))
		cursor := viewCtx.Request.URL.Query().Get("cursor")

		tableConfig, _, err := buildTableConfig(ctx, deps, asOfDate, cursor)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches supplier + candidate data and assembles a TableConfig.
func buildTableConfig(
	ctx context.Context,
	deps *QueueViewDeps,
	asOfDate, cursor string,
) (*types.TableConfig, []queueform.QueueRow, error) {
	l := deps.Labels
	lq := l.Queue
	perms := view.GetUserPermissions(ctx)

	// 1. Fetch the visible supplier page.
	var suppliers []errshared.QueueSupplierRecord
	var nextCursor string
	if deps.ListSuppliers != nil {
		var err error
		suppliers, nextCursor, err = deps.ListSuppliers(ctx, cursor)
		if err != nil {
			log.Printf("expense-recognition-run queue: failed to list suppliers: %v", err)
			return nil, nil, fmt.Errorf("failed to load supplier list: %w", err)
		}
	}

	// 2. Fan-out: one goroutine per supplier, bounded by EXPENSE_RUN_QUEUE_FANOUT.
	fanout := fanoutLimit()
	// Cap at 50 — surface deferred filtering per plan.
	if len(suppliers) > 50 {
		suppliers = suppliers[:50]
	}

	type supplierResult struct {
		idx     int
		summary CandidateSummary
		errMsg  string
	}
	results := make([]supplierResult, len(suppliers))

	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(fanout)

	var mu sync.Mutex
	for i, s := range suppliers {
		idx := i
		supplierID := s.ID
		eg.Go(func() (retErr error) {
			defer func() {
				if r := recover(); r != nil {
					msg := fmt.Sprintf("panic: %v", r)
					log.Printf("expense-recognition-run queue: panic for supplier %s: %s", supplierID, msg)
					mu.Lock()
					results[idx] = supplierResult{idx: idx, errMsg: msg}
					mu.Unlock()
				}
			}()

			var summary CandidateSummary
			var rowErr string
			if deps.ListExpenseRunCandidates != nil {
				candidates, err := deps.ListExpenseRunCandidates(egCtx, supplierID, asOfDate)
				if err != nil {
					log.Printf("expense-recognition-run queue: ListExpenseRunCandidates error for supplier %s: %v", supplierID, err)
					rowErr = err.Error()
				} else {
					summary = aggregateCandidates(candidates)
				}
			}

			mu.Lock()
			results[idx] = supplierResult{idx: idx, summary: summary, errMsg: rowErr}
			mu.Unlock()
			return nil
		})
	}
	_ = eg.Wait()

	// 3. Build QueueRow list, filtering out suppliers with zero pending periods.
	queueRows := make([]queueform.QueueRow, 0, len(suppliers))
	for i, s := range suppliers {
		res := results[i]
		if res.errMsg == "" && res.summary.PeriodCount == 0 {
			continue
		}

		drawerURL := ""
		if deps.SupplierDrawerURLTemplate != "" {
			drawerURL = route.ResolveURL(deps.SupplierDrawerURLTemplate, "id", s.ID)
		}
		supplierDetailURL := ""
		if deps.SupplierDetailURLTemplate != "" {
			supplierDetailURL = route.ResolveURL(deps.SupplierDetailURLTemplate, "id", s.ID)
		}

		var totalAmount int64
		currency := ""
		multiCurrency := false
		if res.errMsg == "" {
			totalAmount, currency, multiCurrency = flattenCurrencyTotals(res.summary.TotalByCurrency)
		}

		queueRows = append(queueRows, queueform.QueueRow{
			SupplierID:               s.ID,
			SupplierName:             s.Name,
			SupplierDetailURL:        supplierDetailURL,
			DrawerURL:                drawerURL,
			SubscriptionCount:        len(res.summary.SubscriptionIDs),
			AdvanceDisbursementCount: len(res.summary.AdvanceDisbursementIDs),
			PendingPeriods:           res.summary.PeriodCount,
			TotalAmount:              totalAmount,
			Currency:                 currency,
			MultiCurrency:            multiCurrency,
			ErrorMessage:             res.errMsg,
		})
	}

	// 4. Convert to table rows.
	columns := queueColumns(l)
	tableRows := buildTableRows(queueRows, l, perms)
	types.ApplyColumnStyles(columns, tableRows)

	sp := &types.ServerPagination{
		Enabled:       true,
		Mode:          "cursor",
		PaginationURL: deps.Routes.QueueTableURL,
	}
	if nextCursor != "" {
		sp.NextCursor = nextCursor
	}
	sp.BuildDisplay()

	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:             "run-selected",
			Label:           lq.Bulk.RunSelected,
			Icon:            "icon-zap",
			Variant:         "primary",
			Endpoint:        deps.Routes.SubmitBatchURL,
			ExtraParamsJSON: `{"selection_mode":"selected"}`,
		},
		{
			Key:             "run-all-matching",
			Label:           lq.Bulk.RunAllMatching,
			Icon:            "icon-zap",
			Variant:         "warning",
			Endpoint:        deps.Routes.SubmitBatchURL,
			ExtraParamsJSON: `{"selection_mode":"all_matching"}`,
		},
	}

	tableConfig := &types.TableConfig{
		ID:          "expense-run-queue-table",
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
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}

	if !perms.Can("expense_recognition_run", "create") {
		tableConfig.BulkActions = nil
	}

	types.ApplyTableSettings(tableConfig)
	return tableConfig, queueRows, nil
}

func queueColumns(l centymo.ExpenseRecognitionRunLabels) []types.TableColumn {
	lc := l.Queue.Columns
	return []types.TableColumn{
		{Key: "supplier", Label: lc.Supplier, WidthClass: "col-8xl"},
		{Key: "subscriptions", Label: lc.Subscriptions, WidthClass: "col-3xl", Align: "right", NoSort: true, NoFilter: true},
		{Key: "advance_disbursements", Label: lc.AdvanceDisbursements, WidthClass: "col-3xl", Align: "right", NoSort: true, NoFilter: true},
		{Key: "pending_periods", Label: lc.PendingPeriods, WidthClass: "col-3xl", Align: "right", NoSort: true, NoFilter: true},
		{Key: "total", Label: lc.Total, WidthClass: "col-4xl", Align: "right", NoSort: true, NoFilter: true},
		{Key: "currency", Label: lc.Currency, WidthClass: "col-2xl", NoSort: true, NoFilter: true},
	}
}

func buildTableRows(
	rows []queueform.QueueRow,
	l centymo.ExpenseRecognitionRunLabels,
	perms *types.UserPermissions,
) []types.TableRow {
	tableRows := make([]types.TableRow, 0, len(rows))
	lq := l.Queue

	for _, r := range rows {
		var supplierCell types.TableCell
		if r.SupplierDetailURL != "" {
			supplierCell = types.TableCell{Type: "link", Value: r.SupplierName, Href: r.SupplierDetailURL}
		} else {
			supplierCell = types.TableCell{Type: "text", Value: r.SupplierName}
		}

		var currencyCell types.TableCell
		if r.ErrorMessage != "" {
			currencyCell = types.TableCell{Type: "badge", Value: l.Errors.InvalidSelection, Variant: "error"}
		} else if r.MultiCurrency {
			currencyCell = types.TableCell{Type: "badge", Value: r.Currency, Variant: "warning"}
		} else {
			currencyCell = types.TableCell{Type: "badge", Value: r.Currency, Variant: "info"}
		}

		var totalCell types.TableCell
		if r.ErrorMessage != "" {
			totalCell = types.TableCell{Type: "text", Value: "—"}
		} else {
			totalCell = types.MoneyCell(float64(r.TotalAmount), r.Currency, true)
		}

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
				Disabled:        !perms.Can("expense_recognition_run", "create"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		}

		tableRows = append(tableRows, types.TableRow{
			ID:   r.SupplierID,
			Href: r.SupplierDetailURL,
			DataAttrs: map[string]string{
				"testid": "queue-row-" + r.SupplierID,
			},
			Cells: []types.TableCell{
				supplierCell,
				{Type: "text", Value: strconv.Itoa(r.SubscriptionCount), Align: "right"},
				{Type: "text", Value: strconv.Itoa(r.AdvanceDisbursementCount), Align: "right"},
				{Type: "text", Value: strconv.Itoa(r.PendingPeriods), Align: "right"},
				totalCell,
				currencyCell,
			},
			Actions: actions,
		})
	}
	return tableRows
}

func aggregateCandidates(candidates []errshared.QueueCandidateInput) CandidateSummary {
	s := CandidateSummary{
		SubscriptionIDs:        make(map[string]struct{}),
		AdvanceDisbursementIDs: make(map[string]struct{}),
		TotalByCurrency:        make(map[string]int64),
	}
	for _, c := range candidates {
		if !c.Eligible {
			continue
		}
		if c.SupplierSubscriptionID != "" {
			s.SubscriptionIDs[c.SupplierSubscriptionID] = struct{}{}
		}
		if c.AdvanceDisbursementID != "" {
			s.AdvanceDisbursementIDs[c.AdvanceDisbursementID] = struct{}{}
		}
		s.TotalByCurrency[c.Currency] += c.Amount
		s.PeriodCount++
	}
	return s
}

func flattenCurrencyTotals(totals map[string]int64) (total int64, currency string, multi bool) {
	if len(totals) == 0 {
		return 0, "", false
	}
	if len(totals) == 1 {
		for cur, amt := range totals {
			return amt, cur, false
		}
	}
	for cur, amt := range totals {
		total += amt
		if currency == "" || cur < currency {
			currency = cur
		}
	}
	return total, currency, true
}

func resolveAsOfDate(ctx context.Context, input string) (asOfDate, maxDate string) {
	tz := types.LocationFromContext(ctx)
	today := time.Now().In(tz).Format(types.DateInputLayout)
	if input == "" {
		return today, today
	}
	if _, err := time.Parse(types.DateInputLayout, input); err != nil {
		return today, today
	}
	return input, today
}

func fanoutLimit() int {
	if v := os.Getenv("EXPENSE_RUN_QUEUE_FANOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 16
}
