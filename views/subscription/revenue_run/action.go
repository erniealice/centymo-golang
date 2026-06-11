// Package revenue_run handles the per-subscription "Run Invoices" drawer
// (Surface C — CYCLE billing_kind). Template:
// subscription-revenue-run-drawer-form.html.
package revenue_run

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	appcontext "github.com/erniealice/espyna-golang/appcontext"

	centymo "github.com/erniealice/centymo-golang"
	revenuedomain "github.com/erniealice/centymo-golang/domain/revenue"
	revenuerunform "github.com/erniealice/centymo-golang/views/subscription/revenue_run/form"
)

// ---------------------------------------------------------------------------
// View-local types — NOT imported from espyna. Views never import espyna
// internals; block.go provides typed callbacks that translate between
// consumer.* shapes and these view-local shapes.
// ---------------------------------------------------------------------------

// RevenueRunScope is the view-layer scope passed to the list / generate callbacks.
type RevenueRunScope struct {
	WorkspaceID    string
	ClientID       string
	SubscriptionID string
	AsOfDate       string // YYYY-MM-DD; empty → use today
	Cursor         string
	Limit          int32
}

// RevenueRunCandidate is the view-layer representation of one pending period.
type RevenueRunCandidate struct {
	SubscriptionID    string
	SubscriptionName  string
	ClientID          string
	ClientName        string
	PlanName          string
	BillingCycleLabel string
	Currency          string
	PeriodStart       string // YYYY-MM-DD
	PeriodEnd         string // YYYY-MM-DD
	PeriodLabel       string
	PeriodMarker      string
	Amount            int64
	AmountDisplay     string
	LineItemCount     int
	Eligible          bool
	BlockerReason     string
	// SourceKind discriminates SUBSCRIPTION_CYCLE vs ADVANCE_COLLECTION
	// (Plan B Phase 5b). Empty defaults to subscription-cycle.
	SourceKind string
	// AdvanceCollectionID is set when SourceKind == ADVANCE_COLLECTION.
	AdvanceCollectionID string
	// SuppressingAdvanceCollectionID is set on cycle rows overlapped by an
	// active TIME_BASED advance Collection (Decision A).
	SuppressingAdvanceCollectionID string
}

// SelectedRevenueRunCandidate is one confirmed selection.
type SelectedRevenueRunCandidate struct {
	SubscriptionID string
	PeriodStart    string
	PeriodEnd      string
	PeriodMarker   string
	// SourceKind discriminates dispatcher branch. Empty defaults to
	// SUBSCRIPTION_CYCLE. "ADVANCE_COLLECTION" routes to
	// AmortizeAdvanceCollection.
	SourceKind string
	// AdvanceCollectionID is required when SourceKind == ADVANCE_COLLECTION.
	AdvanceCollectionID string
}

// RevenueRunSelections carries either an explicit list or a filter token.
type RevenueRunSelections struct {
	ExplicitList []SelectedRevenueRunCandidate
	FilterToken  string
}

// RevenueRunResult is the output of a successful GenerateRevenueRun call.
type RevenueRunResult struct {
	RunID   string
	Status  string
	Created int32
	Skipped int32
	Errored int32
}

// ---------------------------------------------------------------------------
// Deps + NewAction
// ---------------------------------------------------------------------------

// Deps is the dependency subset needed by the per-subscription revenue-run drawer.
// A subset of subscriptionaction.Deps is threaded through from block.go.
type Deps struct {
	Routes       centymo.SubscriptionRoutes
	Labels       centymo.SubscriptionLabels
	CommonLabels pyeza.CommonLabels

	// ListRevenueRunCandidates enumerates un-invoiced billing periods for this
	// subscription. nil-safe — drawer renders empty state when unset.
	ListRevenueRunCandidates func(ctx context.Context, scope RevenueRunScope) ([]RevenueRunCandidate, string, error)

	// GenerateRevenueRun executes the batch run for the selected periods.
	// nil-safe — POST returns an error when unset.
	GenerateRevenueRun func(ctx context.Context, scope RevenueRunScope, sels RevenueRunSelections) (*RevenueRunResult, error)
}

// NewAction returns a view.View that serves the per-subscription Invoice Run drawer.
//
// GET  → renders the drawer form populated with ListRevenueRunCandidates.
// POST → submits the selected periods via GenerateRevenueRun; on success
//
//	returns HX-Trigger headers to close the drawer, fire the toast, and
//	refresh the invoices table.
func NewAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("revenue", "create") || !perms.Can("subscription", "read") {
			return view.HTMXError(deps.Labels.RevenueRun.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.HTMXError(deps.Labels.RevenueRun.Errors.IDRequired)
		}

		if deps.ListRevenueRunCandidates == nil || deps.GenerateRevenueRun == nil {
			return view.HTMXError(deps.Labels.RevenueRun.Errors.UseCaseUnavailable)
		}

		switch viewCtx.Request.Method {
		case http.MethodGet:
			return renderDrawer(ctx, viewCtx, deps, id)
		case http.MethodPost:
			return submitDrawer(ctx, viewCtx, deps, id)
		default:
			return view.HTMXError(deps.Labels.RevenueRun.Errors.InvalidFormData)
		}
	})
}

// ---------------------------------------------------------------------------
// GET handler
// ---------------------------------------------------------------------------

func renderDrawer(
	ctx context.Context,
	viewCtx *view.ViewContext,
	deps *Deps,
	subscriptionID string,
) view.ViewResult {
	l := deps.Labels.RevenueRun

	// Resolve as-of date: prefer query param, fall back to today in workspace TZ.
	tz := pyezatypes.LocationFromContext(ctx)
	today := time.Now().In(tz).Format(pyezatypes.DateInputLayout)

	asOfDate := viewCtx.Request.URL.Query().Get("as_of_date")
	if asOfDate == "" {
		asOfDate = today
	}

	scope := RevenueRunScope{
		WorkspaceID:    appcontext.GetWorkspaceIDFromContext(ctx),
		SubscriptionID: subscriptionID,
		AsOfDate:       asOfDate,
	}

	candidates, _, err := deps.ListRevenueRunCandidates(ctx, scope)
	if err != nil {
		log.Printf("revenue_run.renderDrawer: ListRevenueRunCandidates for sub %s failed: %v", subscriptionID, err)
		return view.HTMXError(l.Errors.UseCaseUnavailable)
	}

	formAction := route.ResolveURL(deps.Routes.RevenueRunURL, "id", subscriptionID)
	fragmentURL := formAction + "?partial=candidates&as_of_date=" + asOfDate

	data := buildDrawerData(candidates, subscriptionID, asOfDate, today, formAction, fragmentURL, l, deps.CommonLabels)

	// Determine which template to render: the outer form or the inner partial.
	// The HTMX inner-swap on date change targets the candidates partial.
	templateName := "subscription-revenue-run-drawer-form"
	if viewCtx.Request.URL.Query().Get("partial") == "candidates" {
		templateName = "subscription-revenue-run-candidates"
	}

	return view.OK(templateName, data)
}

// ---------------------------------------------------------------------------
// POST handler
// ---------------------------------------------------------------------------

func submitDrawer(
	ctx context.Context,
	viewCtx *view.ViewContext,
	deps *Deps,
	subscriptionID string,
) view.ViewResult {
	l := deps.Labels.RevenueRun

	if err := viewCtx.Request.ParseForm(); err != nil {
		return view.HTMXError(l.Errors.InvalidFormData)
	}

	asOfDate := viewCtx.Request.FormValue("as_of_date")
	if asOfDate == "" {
		tz := pyezatypes.LocationFromContext(ctx)
		asOfDate = time.Now().In(tz).Format(pyezatypes.DateInputLayout)
	}

	// Parse "selection" form values: each is "{sub_id}|{start}|{end}|{marker}".
	rawSelections := viewCtx.Request.Form["selection"]
	if len(rawSelections) == 0 {
		return view.HTMXError(l.Errors.SelectOne)
	}

	var sels RevenueRunSelections
	for _, raw := range rawSelections {
		parts := strings.Split(raw, "|")
		// Selection value formats (Plan B Phase 5b):
		//   - SUBSCRIPTION_CYCLE: "{sub_id}|{start}|{end}|{marker}"  (4 parts)
		//   - ADVANCE_COLLECTION: "{advance_id}|{start}|{end}|{marker}|ADVANCE_COLLECTION"
		//                         (5 parts)
		if len(parts) < 4 {
			continue
		}
		sel := SelectedRevenueRunCandidate{
			PeriodStart:  parts[1],
			PeriodEnd:    parts[2],
			PeriodMarker: parts[3],
		}
		if len(parts) >= 5 && parts[4] == "ADVANCE_COLLECTION" {
			sel.SourceKind = "ADVANCE_COLLECTION"
			sel.AdvanceCollectionID = parts[0]
		} else {
			sel.SubscriptionID = parts[0]
		}
		sels.ExplicitList = append(sels.ExplicitList, sel)
	}
	if len(sels.ExplicitList) == 0 {
		return view.HTMXError(l.Errors.SelectOne)
	}

	scope := RevenueRunScope{
		WorkspaceID:    appcontext.GetWorkspaceIDFromContext(ctx),
		SubscriptionID: subscriptionID,
		AsOfDate:       asOfDate,
	}

	result, err := deps.GenerateRevenueRun(ctx, scope, sels)
	if err != nil {
		log.Printf("revenue_run.submitDrawer: GenerateRevenueRun for sub %s failed: %v", subscriptionID, err)
		return view.HTMXError(l.Errors.UseCaseUnavailable)
	}
	if result == nil {
		return view.HTMXError(l.Errors.UseCaseUnavailable)
	}

	// Resolve the lyngua-translated toast text. Substitute Go-template
	// placeholders here so the JS-side lf.Toast receives a plain string.
	toastMessage := strings.NewReplacer(
		"{{.Created}}", fmt.Sprintf("%d", result.Created),
		"{{.Skipped}}", fmt.Sprintf("%d", result.Skipped),
		"{{.Errored}}", fmt.Sprintf("%d", result.Errored),
	).Replace(l.ToastSuccess)

	toastPayload := map[string]any{
		"message": toastMessage,
		"state":   toastStateFromCounts(result.Created, result.Skipped, result.Errored),
	}
	if result.RunID != "" && l.ViewRunLink != "" {
		toastPayload["link"] = map[string]any{
			"url":   route.ResolveURL(revenuedomain.RevenueRunDetailURL, "id", result.RunID),
			"label": l.ViewRunLink,
		}
	}

	triggerPayload, _ := json.Marshal(map[string]any{
		"pyeza:toast":  toastPayload,
		"refreshTable": "subscription-invoices-table",
	})

	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": string(triggerPayload),
		},
	}
}

// toastStateFromCounts maps create/skip/error counts to a toast state.
// All-errored = error, any-errored = warning, otherwise success.
func toastStateFromCounts(created, _ /*skipped*/, errored int32) string {
	if errored > 0 && created == 0 {
		return "error"
	}
	if errored > 0 {
		return "warning"
	}
	return "success"
}

// ---------------------------------------------------------------------------
// Builder helpers
// ---------------------------------------------------------------------------

// buildClientHint substitutes the {client} token in the label template with
// the actual client name. Returns an empty string when either argument is empty.
func buildClientHint(template, name string) string {
	if template == "" || name == "" {
		return ""
	}
	return strings.Replace(template, "{client}", name, 1)
}

// buildDrawerData constructs the template-facing Data from the raw candidate slice.
//
// Plan B Phase 5b: candidates may be either subscription-cycle or
// advance-Collection rows (discriminated by c.SourceKind). The builder splits
// them into two separate template sections.
func buildDrawerData(
	candidates []RevenueRunCandidate,
	subscriptionID string,
	asOfDate, maxAsOfDate string,
	formAction, fragmentURL string,
	l centymo.SubscriptionRevenueRunLabels,
	commonLabels pyeza.CommonLabels,
) *revenuerunform.Data {
	periods := make([]revenuerunform.Period, 0, len(candidates))
	advanceRows := make([]revenuerunform.AdvanceRow, 0)
	eligibleCount := 0

	var subName, planName, currency, clientName string
	for _, c := range candidates {
		if subName == "" && c.SubscriptionName != "" {
			subName = c.SubscriptionName
		}
		if planName == "" && c.PlanName != "" {
			planName = c.PlanName
		}
		if currency == "" && c.Currency != "" {
			currency = c.Currency
		}
		if clientName == "" && c.ClientName != "" {
			clientName = c.ClientName
		}

		if c.SourceKind == "REVENUE_RUN_SOURCE_KIND_ADVANCE_COLLECTION" || c.SourceKind == "ADVANCE_COLLECTION" {
			advID := c.AdvanceCollectionID
			advanceRows = append(advanceRows, revenuerunform.AdvanceRow{
				AdvanceCollectionID: advID,
				Currency:            c.Currency,
				PeriodStart:         c.PeriodStart,
				PeriodEnd:           c.PeriodEnd,
				PeriodMarker:        c.PeriodMarker,
				PeriodLabel:         c.PeriodLabel,
				Amount:              c.Amount,
				AmountDisplay:       c.AmountDisplay,
				Eligible:            c.Eligible,
				BlockerReason:       c.BlockerReason,
				SelectionValue: fmt.Sprintf("%s|%s|%s|%s|ADVANCE_COLLECTION",
					advID, c.PeriodStart, c.PeriodEnd, c.PeriodMarker),
			})
			if c.Eligible {
				eligibleCount++
			}
			continue
		}

		period := revenuerunform.Period{
			SubscriptionID: c.SubscriptionID,
			PeriodStart:    c.PeriodStart,
			PeriodEnd:      c.PeriodEnd,
			PeriodMarker:   c.PeriodMarker,
			PeriodLabel:    c.PeriodLabel,
			Amount:         c.Amount,
			AmountDisplay:  c.AmountDisplay,
			LineItemCount:  c.LineItemCount,
			Eligible:       c.Eligible,
			BlockerReason:  c.BlockerReason,
			SelectionValue: fmt.Sprintf("%s|%s|%s|%s",
				c.SubscriptionID, c.PeriodStart, c.PeriodEnd, c.PeriodMarker),
			SuppressingAdvanceCollectionID: c.SuppressingAdvanceCollectionID,
		}
		if c.Eligible {
			eligibleCount++
		}
		periods = append(periods, period)
	}

	return &revenuerunform.Data{
		FormAction:            formAction,
		FragmentURL:           fragmentURL,
		SubscriptionID:        subscriptionID,
		SubscriptionName:      subName,
		ClientHint:            buildClientHint(l.ClientHintTemplate, clientName),
		PlanName:              planName,
		AsOfDate:              asOfDate,
		MaxAsOfDate:           maxAsOfDate,
		EligibleCount:         eligibleCount,
		Periods:               periods,
		AdvanceCollectionRows: advanceRows,
		Currency:              currency,
		Labels:                l,
		CommonLabels:          commonLabels,
	}
}
