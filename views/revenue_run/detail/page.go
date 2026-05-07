// Package detail implements the revenue-run detail page (Surface D).
// Pattern mirrors packages/centymo-golang/views/revenue/detail/page.go.
package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	detailform "github.com/erniealice/centymo-golang/views/revenue_run/detail/form"
	rrshared "github.com/erniealice/centymo-golang/views/revenue_run/shared"
	"github.com/erniealice/hybra-golang/views/attachment"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
)

// DetailViewDeps holds view dependencies for the detail page.
type DetailViewDeps struct {
	Routes       centymo.RevenueRunRoutes
	Labels       centymo.RevenueRunLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ReadRevenueRun fetches a run + all attempts by ID.
	ReadRevenueRun func(ctx context.Context, id string) (*rrshared.RevenueRunWithAttempts, error)

	// ListRevenueByRunID fetches invoice records for the Invoices tab.
	ListRevenueByRunID func(ctx context.Context, runID string) ([]rrshared.RevenueRow, error)

	attachment.AttachmentOps
}

// NewView creates the full-page revenue-run detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		runWithAttempts, err := deps.ReadRevenueRun(ctx, id)
		if err != nil {
			log.Printf("Failed to read revenue run %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load run: %w", err))
		}
		if runWithAttempts == nil {
			log.Printf("Revenue run %s not found", id)
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
			ContentTemplate:       "revenue-run-detail-content",
			Run:                   run,
			Attempts:              runWithAttempts.Attempts,
			IsPossiblyInterrupted: run.IsStalePending,
			ActiveTab:             activeTab,
			TabItems:              tabItems,
			Labels:                l,
		}

		loadTabData(ctx, pageData, deps, runWithAttempts, activeTab)

		return view.OK("revenue-run-detail", pageData)
	})
}

// NewTabAction creates a partial view that returns only the active tab content.
// Called via HTMX when the user clicks a tab button.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "summary"
		}

		runWithAttempts, err := deps.ReadRevenueRun(ctx, id)
		if err != nil {
			log.Printf("Failed to read revenue run %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load run: %w", err))
		}
		if runWithAttempts == nil {
			log.Printf("Revenue run %s not found", id)
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

		templateName := "revenue-run-" + tab + "-tab"
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// loadTabData populates tab-specific fields on pageData.
func loadTabData(
	ctx context.Context,
	pageData *detailform.PageData,
	deps *DetailViewDeps,
	runWithAttempts *rrshared.RevenueRunWithAttempts,
	tab string,
) {
	l := deps.Labels
	attempts := runWithAttempts.Attempts

	switch tab {
	case "summary":
		// All data is already in Run; nothing extra to load.

	case "selections":
		pageData.SelectionsTable = buildSelectionsTable(attempts, l, deps.TableLabels)

	case "results":
		pageData.ResultsTable = buildResultsTable(attempts, l, deps.TableLabels)

	case "invoices":
		if deps.ListRevenueByRunID != nil {
			revenues, err := deps.ListRevenueByRunID(ctx, runWithAttempts.Run.ID)
			if err != nil {
				log.Printf("Failed to load invoices for run %s: %v", runWithAttempts.Run.ID, err)
				revenues = []rrshared.RevenueRow{}
			}
			pageData.InvoicesTable = buildInvoicesTable(revenues, l, deps.TableLabels)
		} else {
			pageData.InvoicesTable = buildInvoicesTable(nil, l, deps.TableLabels)
		}

	case "audit-history":
		// Deferred — rendered as an info alert in the template.

	case "attachments":
		if deps.ListAttachments != nil {
			cfg := attachmentConfig(deps)
			var attachItems []*attachmentpb.Attachment
			if resp, err := deps.ListAttachments(ctx, cfg.EntityType, runWithAttempts.Run.ID); err == nil && resp != nil {
				attachItems = resp.GetData()
			}
			pageData.AttachmentTable = attachment.BuildTable(attachItems, cfg, runWithAttempts.Run.ID)
		}
	}
}

// buildTabItems constructs the tab bar items.
func buildTabItems(l centymo.RevenueRunLabels, id string, routes centymo.RevenueRunRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.DetailTabActionURL, "id", id, "tab", "")
	lt := l.Detail.Tabs
	attachmentsLabel := lt.Attachments
	if attachmentsLabel == "" {
		attachmentsLabel = "Attachments"
	}
	return []pyeza.TabItem{
		{Key: "summary", Label: lt.Summary, Href: base + "?tab=summary", HxGet: action + "summary", Icon: "icon-info"},
		{Key: "selections", Label: lt.Selections, Href: base + "?tab=selections", HxGet: action + "selections", Icon: "icon-list"},
		{Key: "results", Label: lt.Results, Href: base + "?tab=results", HxGet: action + "results", Icon: "icon-check-circle"},
		{Key: "invoices", Label: lt.Invoices, Href: base + "?tab=invoices", HxGet: action + "invoices", Icon: "icon-file-text"},
		{Key: "audit-history", Label: lt.AuditHistory, Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
		{Key: "attachments", Label: attachmentsLabel, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
	}
}

// buildSelectionsTable builds the TableConfig for the Selections tab.
// Shows all attempts (all are selections — created, skipped, errored).
func buildSelectionsTable(attempts []rrshared.RevenueRunAttemptRow, l centymo.RevenueRunLabels, tableLabels types.TableLabels) *types.TableConfig {
	ls := l.Detail.Selections
	columns := []types.TableColumn{
		{Key: "subscription", Label: ls.ColSubscription, NoSort: true},
		{Key: "period_start", Label: ls.ColPeriodStart, NoSort: true, WidthClass: "col-3xl"},
		{Key: "period_end", Label: ls.ColPeriodEnd, NoSort: true, WidthClass: "col-3xl"},
		{Key: "period_marker", Label: ls.ColPeriodMarker, NoSort: true, WidthClass: "col-3xl"},
	}

	rows := make([]types.TableRow, 0, len(attempts))
	for _, a := range attempts {
		subscriptionDisplay := a.SubscriptionID
		if a.SubscriptionName != "" {
			subscriptionDisplay = a.SubscriptionName
		}
		rows = append(rows, types.TableRow{
			ID: a.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: subscriptionDisplay},
				{Type: "text", Value: a.PeriodStart},
				{Type: "text", Value: a.PeriodEnd},
				{Type: "text", Value: a.PeriodMarker},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "revenue-run-selections-table",
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
// Same data as selections but adds Outcome + Error columns.
func buildResultsTable(attempts []rrshared.RevenueRunAttemptRow, l centymo.RevenueRunLabels, tableLabels types.TableLabels) *types.TableConfig {
	lr := l.Detail.Results
	columns := []types.TableColumn{
		{Key: "subscription", Label: lr.ColSubscription, NoSort: true},
		{Key: "period_start", Label: lr.ColPeriodStart, NoSort: true, WidthClass: "col-3xl"},
		{Key: "period_end", Label: lr.ColPeriodEnd, NoSort: true, WidthClass: "col-3xl"},
		{Key: "outcome", Label: lr.ColOutcome, NoSort: true, WidthClass: "col-3xl"},
		{Key: "error_code", Label: lr.ColErrorCode, NoSort: true, WidthClass: "col-4xl"},
	}

	rows := make([]types.TableRow, 0, len(attempts))
	for _, a := range attempts {
		subscriptionDisplay := a.SubscriptionID
		if a.SubscriptionName != "" {
			subscriptionDisplay = a.SubscriptionName
		}
		outcomeLabel, outcomeVariant := outcomeCell(l, a.Outcome)
		rows = append(rows, types.TableRow{
			ID: a.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: subscriptionDisplay},
				{Type: "text", Value: a.PeriodStart},
				{Type: "text", Value: a.PeriodEnd},
				{Type: "badge", Value: outcomeLabel, Variant: outcomeVariant},
				{Type: "text", Value: a.ErrorCode},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "revenue-run-results-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   lr.EmptyTitle,
			Message: lr.EmptyMessage,
		},
	}
}

// buildInvoicesTable builds the TableConfig for the Invoices tab.
func buildInvoicesTable(revenues []rrshared.RevenueRow, l centymo.RevenueRunLabels, tableLabels types.TableLabels) *types.TableConfig {
	li := l.Detail.Invoices
	columns := []types.TableColumn{
		{Key: "reference", Label: li.ColReference, WidthClass: "col-5xl"},
		{Key: "date", Label: li.ColDate, WidthClass: "col-3xl"},
		{Key: "amount", Label: li.ColAmount, WidthClass: "col-3xl", Align: "right"},
		{Key: "status", Label: li.ColStatus, WidthClass: "col-3xl", NoSort: true},
	}

	rows := make([]types.TableRow, 0, len(revenues))
	for _, rv := range revenues {
		statusLabel, statusVariant := revenueStatusCell(rv.Status)
		var actions []types.TableAction
		if rv.DetailURL != "" {
			actions = append(actions, types.TableAction{
				Type:  "view",
				Label: l.Actions.ViewRun,
				Href:  rv.DetailURL,
			})
		}
		rows = append(rows, types.TableRow{
			ID:   rv.ID,
			Href: rv.DetailURL,
			Cells: []types.TableCell{
				{Type: "text", Value: rv.ReferenceNumber},
				types.DateTimeCell(rv.RevenueDate, types.DateReadable),
				types.MoneyCell(float64(rv.TotalAmount), rv.Currency, true),
				{Type: "badge", Value: statusLabel, Variant: statusVariant},
			},
			Actions: actions,
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "revenue-run-invoices-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   li.EmptyTitle,
			Message: li.EmptyMessage,
		},
	}
}

func outcomeCell(l centymo.RevenueRunLabels, outcome string) (label, variant string) {
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

func revenueStatusCell(status string) (label, variant string) {
	switch status {
	case "complete":
		return status, "success"
	case "draft":
		return status, "info"
	case "cancelled":
		return status, "warning"
	default:
		return status, "info"
	}
}

