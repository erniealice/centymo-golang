// Package recognize handles the "Recognize Revenue" feature for subscriptions.
// Drawer template: subscription-recognize-drawer-form.html (stays flat at view root).
package recognize

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	recognizeform "github.com/erniealice/centymo-golang/views/subscription/recognize/form"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	billingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billing_event"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// Deps is the dependency subset needed by the recognize feature.
// A subset of action.Deps is threaded through from block.go.
type Deps struct {
	Routes centymo.SubscriptionRoutes
	Labels centymo.SubscriptionLabels

	ReadSubscription                 func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
	ListClients                      func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ReadPricePlan                    func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	RecognizeRevenueFromSubscription func(ctx context.Context, req *revenuepb.CreateRevenueWithLineItemsRequest) (*revenuepb.CreateRevenueWithLineItemsResponse, error)
	ListBillingEventsBySubscription  func(ctx context.Context, req *billingeventpb.ListBillingEventsBySubscriptionRequest) (*billingeventpb.ListBillingEventsBySubscriptionResponse, error)
}

// NewAction creates the subscription "Recognize Revenue" view.
//
//	GET  → renders the drawer with a dry_run preview from the use case.
//	POST → calls the use case for real, returning HTMXSuccess on success so
//	       the invoices table refreshes inline (no redirect).
func NewAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("revenue", "create") || !perms.Can("subscription", "read") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		subscriptionID := viewCtx.Request.PathValue("id")
		if subscriptionID == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		// Resolve subscription for the read-only context block + currency
		// determination. Done once per call (GET preview + POST submit both need it).
		sub, client, pricePlan := loadContext(ctx, deps, subscriptionID)
		if sub == nil {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}

		tz := pyezatypes.LocationFromContext(ctx)

		switch viewCtx.Request.Method {
		case http.MethodGet:
			return renderDrawer(ctx, deps, subscriptionID, sub, client, pricePlan, tz)
		case http.MethodPost:
			return submitDrawer(ctx, deps, viewCtx, subscriptionID, sub, client, pricePlan, tz)
		default:
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
	})
}

// loadContext resolves the (subscription, client, pricePlan) triple.
// Best-effort — nil values fall through to empty labels in the template.
func loadContext(
	ctx context.Context, deps *Deps, subscriptionID string,
) (*subscriptionpb.Subscription, *clientpb.Client, *priceplanpb.PricePlan) {
	if deps.ReadSubscription == nil {
		return nil, nil, nil
	}
	subResp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
		Data: &subscriptionpb.Subscription{Id: subscriptionID},
	})
	if err != nil || subResp == nil || len(subResp.GetData()) == 0 {
		return nil, nil, nil
	}
	sub := subResp.GetData()[0]

	var client *clientpb.Client
	if deps.ListClients != nil && sub.GetClientId() != "" {
		if cResp, cErr := deps.ListClients(ctx, &clientpb.ListClientsRequest{}); cErr == nil {
			for _, c := range cResp.GetData() {
				if c.GetId() == sub.GetClientId() {
					client = c
					break
				}
			}
		}
	}

	var pricePlan *priceplanpb.PricePlan
	if deps.ReadPricePlan != nil && sub.GetPricePlanId() != "" {
		if ppResp, ppErr := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
			Data: &priceplanpb.PricePlan{Id: sub.GetPricePlanId()},
		}); ppErr == nil && len(ppResp.GetData()) > 0 {
			pricePlan = ppResp.GetData()[0]
		}
	}

	return sub, client, pricePlan
}

// renderDrawer handles GET — runs a dry_run through the use case and
// renders the drawer with the preview.
func renderDrawer(
	ctx context.Context, deps *Deps,
	subscriptionID string,
	sub *subscriptionpb.Subscription,
	client *clientpb.Client,
	pricePlan *priceplanpb.PricePlan,
	tz *time.Location,
) view.ViewResult {
	defaultPeriodStart, defaultPeriodEnd := defaultPeriodBounds(sub, pricePlan, tz)
	revenueDate := time.Now().In(tz).Format(pyezatypes.DateInputLayout)

	data := buildBaseData(deps, subscriptionID, sub, client, pricePlan, tz,
		defaultPeriodStart, defaultPeriodEnd, revenueDate)

	// 2026-04-29 milestone-billing plan §5 — detect MILESTONE plans.
	isMilestone := pricePlan != nil &&
		pricePlan.GetBillingKind() == priceplanpb.BillingKind_BILLING_KIND_MILESTONE
	data.IsMilestone = isMilestone

	var selectedEvent *billingeventpb.BillingEvent
	if isMilestone {
		options, picked := loadMilestoneOptions(ctx, deps, subscriptionID, "")
		data.MilestoneOptions = options
		if picked != nil {
			selectedEvent = picked
			data.SelectedBillingEventID = picked.GetId()
			data.BillAmountDisplay = formatCentavos(picked.GetBillableAmount())
		}
	}

	// Run the use case in dry_run mode to compute the preview.
	if deps.RecognizeRevenueFromSubscription != nil {
		var req *revenuepb.CreateRevenueWithLineItemsRequest
		if isMilestone && selectedEvent != nil {
			req = buildMilestoneRequest(subscriptionID, selectedEvent.GetId(), revenueDate)
		} else {
			req = buildPreviewRequest(subscriptionID, defaultPeriodStart, defaultPeriodEnd, revenueDate)
		}
		req.DryRun = boolPtr(true)
		resp, err := deps.RecognizeRevenueFromSubscription(ctx, req)
		applyResponse(&data, resp, err, deps, client, pricePlan)
	}

	return view.OK("subscription-recognize-drawer-form", &data)
}

// loadMilestoneOptions fetches BillingEvent rows for the given subscription,
// converts them to the drawer's option shape, and returns the pre-selected
// event (first READY → first DEFERRED → nil). UNSPECIFIED (pending) events
// are marked Hidden=true so the template can skip rendering them.
// selectedID="" means "auto-pick"; non-empty restores the operator's prior choice.
func loadMilestoneOptions(
	ctx context.Context, deps *Deps, subscriptionID, selectedID string,
) ([]recognizeform.MilestoneOption, *billingeventpb.BillingEvent) {
	if deps.ListBillingEventsBySubscription == nil {
		return nil, nil
	}
	resp, err := deps.ListBillingEventsBySubscription(ctx, &billingeventpb.ListBillingEventsBySubscriptionRequest{
		SubscriptionId: subscriptionID,
	})
	if err != nil || resp == nil {
		return nil, nil
	}
	events := resp.GetBillingEvents()
	mLabels := deps.Labels.Milestone
	options := make([]recognizeform.MilestoneOption, 0, len(events))
	var firstReady, firstDeferred *billingeventpb.BillingEvent
	var requested *billingeventpb.BillingEvent
	for _, ev := range events {
		status := ev.GetStatus()
		statusKey := statusKey(status)
		hidden := status == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_UNSPECIFIED
		selectable := status == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_READY ||
			status == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_DEFERRED
		opt := recognizeform.MilestoneOption{
			EventID:         ev.GetId(),
			SequenceLabel:   sequenceLabel(ev),
			Status:          statusKey,
			StatusLabel:     statusLabel(status, mLabels),
			BillableAmount:  ev.GetBillableAmount(),
			BillableDisplay: formatCentavos(ev.GetBillableAmount()),
			Currency:        ev.GetBillingCurrency(),
			Selectable:      selectable,
			Hidden:          hidden,
		}
		if firstReady == nil && status == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_READY {
			firstReady = ev
		}
		if firstDeferred == nil && status == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_DEFERRED {
			firstDeferred = ev
		}
		if selectedID != "" && ev.GetId() == selectedID {
			requested = ev
		}
		options = append(options, opt)
	}
	picked := requested
	if picked == nil {
		picked = firstReady
	}
	if picked == nil {
		picked = firstDeferred
	}
	if picked != nil {
		for i := range options {
			if options[i].EventID == picked.GetId() {
				options[i].Selected = true
			}
		}
	}
	return options, picked
}

// submitDrawer handles POST — runs the use case for real and (on success)
// returns the HTMXSuccess header bundle that closes the drawer and triggers
// the invoices table refresh.
func submitDrawer(
	ctx context.Context, deps *Deps, viewCtx *view.ViewContext,
	subscriptionID string,
	sub *subscriptionpb.Subscription,
	client *clientpb.Client,
	pricePlan *priceplanpb.PricePlan,
	tz *time.Location,
) view.ViewResult {
	if err := viewCtx.Request.ParseForm(); err != nil {
		return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
	}
	r := viewCtx.Request
	periodStart := readISODateTime(r.FormValue("period_start_iso"),
		r.FormValue("period_start_date"), r.FormValue("period_start_time"), tz)
	periodEnd := readISODateTime(r.FormValue("period_end_iso"),
		r.FormValue("period_end_date"), r.FormValue("period_end_time"), tz)
	revenueDate := strings.TrimSpace(r.FormValue("revenue_date"))
	notes := strings.TrimSpace(r.FormValue("notes"))

	// 2026-04-29 milestone-billing — extract milestone-specific fields.
	billingEventID := strings.TrimSpace(r.FormValue("billing_event_id"))
	billAmountStr := strings.TrimSpace(r.FormValue("bill_amount"))
	leaveRemainder := r.FormValue("leave_remainder_open") == "true" || r.FormValue("leave_remainder_open") == "on"
	closeShort := r.FormValue("close_short") == "true" || r.FormValue("close_short") == "on"
	partialReason := strings.TrimSpace(r.FormValue("partial_reason"))
	isMilestone := pricePlan != nil &&
		pricePlan.GetBillingKind() == priceplanpb.BillingKind_BILLING_KIND_MILESTONE

	if deps.RecognizeRevenueFromSubscription == nil {
		return centymo.HTMXError("recognize-revenue use case not configured")
	}

	var req *revenuepb.CreateRevenueWithLineItemsRequest
	if isMilestone && billingEventID != "" {
		req = buildMilestoneRequest(subscriptionID, billingEventID, revenueDate)
		if billAmountStr != "" {
			if cents, ok := parseCentavoDecimal(billAmountStr); ok {
				if leaveRemainder || closeShort {
					override := cents
					req.OverrideTotalAmount = &override
				}
			}
		}
		if partialReason != "" {
			pr := partialReason
			req.PartialReason = &pr
		}
		if leaveRemainder {
			lr := true
			req.LeaveRemainderOpen = &lr
		}
	} else {
		req = buildPreviewRequest(subscriptionID, periodStart, periodEnd, revenueDate)
	}
	req.DryRun = boolPtr(false)
	if notes != "" {
		req.Data = &revenuepb.Revenue{Notes: stringPtr(notes)}
	}

	resp, err := deps.RecognizeRevenueFromSubscription(ctx, req)
	if err != nil {
		data := buildBaseData(deps, subscriptionID, sub, client, pricePlan, tz,
			periodStart, periodEnd, revenueDate)
		applyResponse(&data, resp, err, deps, client, pricePlan)
		if isMilestone {
			data.IsMilestone = true
			options, picked := loadMilestoneOptions(ctx, deps, subscriptionID, billingEventID)
			data.MilestoneOptions = options
			if picked != nil {
				data.SelectedBillingEventID = picked.GetId()
			} else {
				data.SelectedBillingEventID = billingEventID
			}
			data.BillAmountDisplay = billAmountStr
			data.LeaveRemainderOpen = leaveRemainder
			data.CloseShort = closeShort
			data.PartialReasonValue = partialReason
			if err != nil && strings.Contains(strings.ToLower(err.Error()), "over") &&
				strings.Contains(strings.ToLower(err.Error()), "bill") {
				data.OverBillingError = true
			}
		}
		// Re-run dry_run so the operator can adjust and retry.
		var previewReq *revenuepb.CreateRevenueWithLineItemsRequest
		if isMilestone && data.SelectedBillingEventID != "" {
			previewReq = buildMilestoneRequest(subscriptionID, data.SelectedBillingEventID, revenueDate)
		} else {
			previewReq = buildPreviewRequest(subscriptionID, periodStart, periodEnd, revenueDate)
		}
		previewReq.DryRun = boolPtr(true)
		if previewResp, _ := deps.RecognizeRevenueFromSubscription(ctx, previewReq); previewResp != nil {
			data.PreviewLines = convertPreviewLines(previewResp.GetPreviewLines(), pricePlan, deps.Labels)
			data.TotalAmount = sumPreview(data.PreviewLines)
		}
		log.Printf("Recognize revenue from subscription %s failed: %v", subscriptionID, err)
		return view.ViewResult{
			Template:   "subscription-recognize-drawer-form",
			Data:       &data,
			StatusCode: http.StatusUnprocessableEntity,
			Headers: map[string]string{
				"HX-Reswap":   "outerHTML",
				"HX-Retarget": "#sheet form",
			},
		}
	}

	// Success — close the drawer and refresh the invoices table + milestones.
	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": `{"formSuccess":true,"refreshTable":"subscription-invoices-table","refresh-invoices":true,"refresh-milestones":true}`,
		},
	}
}

// --- dep-bearing helpers ---

// buildBaseData populates the form data with the read-only context (no
// preview lines yet — those are layered on after the use case dry-run).
func buildBaseData(
	deps *Deps,
	subscriptionID string,
	sub *subscriptionpb.Subscription,
	client *clientpb.Client,
	pricePlan *priceplanpb.PricePlan,
	tz *time.Location,
	periodStart, periodEnd, revenueDate string,
) recognizeform.Data {
	startDate, startTime, startISO := splitISOForInputs(periodStart, tz)
	endDate, endTime, endISO := splitISOForInputs(periodEnd, tz)

	clientLabel := ""
	clientCurrency := ""
	if client != nil {
		clientLabel = client.GetName()
		if u := client.GetUser(); clientLabel == "" && u != nil {
			clientLabel = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
		}
		clientCurrency = client.GetBillingCurrency()
	}

	planLabel := ""
	planCurrency := ""
	if pricePlan != nil {
		planLabel = pricePlan.GetName()
		if planLabel == "" {
			if p := pricePlan.GetPlan(); p != nil {
				planLabel = p.GetName()
			}
		}
		planCurrency = pricePlan.GetBillingCurrency()
	}

	formAction := strings.ReplaceAll(deps.Routes.RecognizeURL, "{id}", subscriptionID)

	return recognizeform.Data{
		FormAction:      formAction,
		SubscriptionID:  subscriptionID,
		SubscriptionName: sub.GetName(),
		ClientLabel:     clientLabel,
		PlanLabel:       planLabel,
		Quantity:        sub.GetQuantity(),
		Currency:        planCurrency,
		PeriodStartDate: startDate,
		PeriodStartTime: startTime,
		PeriodStartISO:  startISO,
		PeriodEndDate:   endDate,
		PeriodEndTime:   endTime,
		PeriodEndISO:    endISO,
		DefaultTZ:       tz.String(),
		RevenueDate:     revenueDate,
		ClientCurrency:  clientCurrency,
		PlanCurrency:    planCurrency,
		Labels:          buildLabels(deps.Labels),
		CommonLabels:    nil, // injected by ViewAdapter
	}
}

// buildLabels assembles the recognize drawer's flat Labels from centymo.SubscriptionLabels.
func buildLabels(l centymo.SubscriptionLabels) recognizeform.Labels {
	return recognizeform.Labels{
		Title:                   l.Invoices.RecognizeTitle,
		Subtitle:                l.Invoices.RecognizeSubtitle,
		ContextSection:          l.Recognize.ContextSection,
		ClientLabel:             l.Recognize.ClientLabel,
		PlanLabel:               l.Recognize.PlanLabel,
		QuantityLabel:           l.Recognize.QuantityLabel,
		PeriodSection:           l.Recognize.PeriodSection,
		PeriodStart:             l.Recognize.PeriodStart,
		PeriodEnd:               l.Recognize.PeriodEnd,
		RevenueDate:             l.Recognize.RevenueDate,
		LineItemsSection:        l.Recognize.LineItemsSection,
		ColumnDescription:       l.Recognize.ColumnDescription,
		ColumnUnitPrice:         l.Recognize.ColumnUnitPrice,
		ColumnQuantity:          l.Recognize.ColumnQuantity,
		ColumnLineTotal:         l.Recognize.ColumnLineTotal,
		ColumnTreatment:         l.Recognize.ColumnTreatment,
		TotalLabel:              l.Recognize.TotalLabel,
		RemoveLine:              l.Recognize.RemoveLine,
		TreatmentRecurring:      l.Recognize.TreatmentRecurring,
		TreatmentFirstCycle:     l.Recognize.TreatmentFirstCycle,
		TreatmentUsageBased:     l.Recognize.TreatmentUsageBased,
		TreatmentOneTime:        l.Recognize.TreatmentOneTime,
		NotesLabel:              l.Recognize.NotesLabel,
		NotesPlaceholder:        l.Recognize.NotesPlaceholder,
		Generate:                l.Recognize.Generate,
		Cancel:                  l.Recognize.Cancel,
		Timezone:                l.Form.Timezone,
		StartDateInfo:           l.Form.StartDateInfo,
		EndDateInfo:             l.Form.EndDateInfo,
		StartTimeInfo:           l.Form.StartTimeInfo,
		EndTimeInfo:             l.Form.EndTimeInfo,
		CurrencyMismatchError:   l.Recognize.CurrencyMismatchError,
		IdempotencyError:        l.Recognize.IdempotencyError,
		IdempotencyExistingLink: l.Recognize.IdempotencyExistingLink,
		NoLinesError:            l.Recognize.NoLinesError,
		MilestoneSelect:            l.Recognize.MilestoneSelect,
		MilestoneSelectPlaceholder: l.Recognize.MilestoneSelectPlaceholder,
		NoReadyMilestone:           l.Recognize.NoReadyMilestone,
		MilestoneNotApplicable:     l.Recognize.MilestoneNotApplicable,
		BillAmount:                 l.Recognize.BillAmount,
		LeaveRemainderOpen:         l.Recognize.LeaveRemainderOpen,
		CloseShort:                 l.Recognize.CloseShort,
		PartialReason:              l.Recognize.PartialReason,
		PartialReasonRequired:      l.Recognize.PartialReasonRequired,
		OverBillingRejected:        l.Recognize.OverBillingRejected,
	}
}

// applyResponse overlays preview lines / warnings / error banners
// from a use-case response onto the form data.
func applyResponse(
	data *recognizeform.Data,
	resp *revenuepb.CreateRevenueWithLineItemsResponse,
	err error,
	deps *Deps,
	client *clientpb.Client,
	pricePlan *priceplanpb.PricePlan,
) {
	if resp != nil {
		data.PreviewLines = convertPreviewLines(resp.GetPreviewLines(), pricePlan, deps.Labels)
		data.TotalAmount = sumPreview(data.PreviewLines)
		data.Warnings = resp.GetWarnings()
		if cid := resp.GetConflictingRevenueId(); cid != "" {
			data.IdempotencyConflict = true
			data.ConflictingRevenueID = cid
			data.ConflictingRevenueURL = strings.ReplaceAll(centymo.RevenueDetailURL, "{id}", cid)
		}
	}
	if client != nil && pricePlan != nil &&
		client.GetBillingCurrency() != "" &&
		pricePlan.GetBillingCurrency() != "" &&
		client.GetBillingCurrency() != pricePlan.GetBillingCurrency() {
		data.CurrencyMismatch = true
	}
	if err != nil && strings.Contains(err.Error(), "no line items") {
		data.NoLinesToInvoice = true
	}
}

// convertPreviewLines maps proto preview lines into the template's row shape.
func convertPreviewLines(
	in []*revenuepb.PreviewLineItem,
	pricePlan *priceplanpb.PricePlan,
	labels centymo.SubscriptionLabels,
) []recognizeform.PreviewLine {
	currency := ""
	if pricePlan != nil {
		currency = pricePlan.GetBillingCurrency()
	}
	out := make([]recognizeform.PreviewLine, 0, len(in))
	for _, p := range in {
		lineCurrency := currency
		if p.GetCurrency() != "" {
			lineCurrency = p.GetCurrency()
		}
		out = append(out, recognizeform.PreviewLine{
			ProductPricePlanID: p.GetProductPricePlanId(),
			Description:        p.GetDescription(),
			UnitPrice:          p.GetUnitPrice(),
			Quantity:           p.GetQuantity(),
			TotalPrice:         p.GetTotalPrice(),
			Currency:           lineCurrency,
			Treatment:          p.GetTreatment(),
			TreatmentLabel:     treatmentLabel(p.GetTreatment(), labels),
		})
	}
	return out
}

func sumPreview(lines []recognizeform.PreviewLine) int64 {
	var total int64
	for _, l := range lines {
		total += l.TotalPrice
	}
	return total
}

// treatmentLabel maps the use-case treatment token to a localized badge label.
func treatmentLabel(t string, labels centymo.SubscriptionLabels) string {
	switch t {
	case "recurring":
		return labels.Recognize.TreatmentRecurring
	case "first_cycle":
		return labels.Recognize.TreatmentFirstCycle
	case "usage_based":
		return labels.Recognize.TreatmentUsageBased
	case "one_time":
		return labels.Recognize.TreatmentOneTime
	default:
		return ""
	}
}

// defaultPeriodBounds returns sane RFC3339 defaults for the drawer's period inputs.
func defaultPeriodBounds(
	sub *subscriptionpb.Subscription, pricePlan *priceplanpb.PricePlan, tz *time.Location,
) (string, string) {
	now := time.Now().In(tz)
	var start time.Time
	if dts := sub.GetDateTimeStart(); dts != nil && dts.IsValid() {
		start = dts.AsTime().In(tz)
	} else {
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz)
	}
	var end time.Time
	if pricePlan != nil && pricePlan.GetBillingCycleValue() > 0 {
		end = addBillingCycle(start, int(pricePlan.GetBillingCycleValue()), pricePlan.GetBillingCycleUnit())
	} else if dte := sub.GetDateTimeEnd(); dte != nil && dte.IsValid() {
		end = dte.AsTime().In(tz)
	} else {
		end = start.AddDate(0, 1, 0)
	}
	return start.Format(time.RFC3339), end.Format(time.RFC3339)
}

func addBillingCycle(t time.Time, value int, unit string) time.Time {
	switch strings.ToLower(unit) {
	case "day", "days":
		return t.AddDate(0, 0, value)
	case "week", "weeks":
		return t.AddDate(0, 0, value*7)
	case "month", "months":
		return t.AddDate(0, value, 0)
	case "year", "years":
		return t.AddDate(value, 0, 0)
	default:
		return t.AddDate(0, value, 0)
	}
}

// splitISOForInputs splits an RFC3339 timestamp into (date, time, RFC3339)
// suitable for the drawer's date+time grid.
func splitISOForInputs(iso string, tz *time.Location) (string, string, string) {
	if iso == "" {
		return "", "", ""
	}
	parsed, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return "", "", ""
	}
	moment := parsed.In(tz)
	return moment.Format(pyezatypes.DateInputLayout),
		moment.Format(pyezatypes.TimeInputLayout),
		moment.Format(time.RFC3339)
}

// readISODateTime returns a canonical RFC3339 string from the drawer form.
func readISODateTime(iso, dateStr, timeStr string, tz *time.Location) string {
	iso = strings.TrimSpace(iso)
	if iso != "" {
		if _, err := time.Parse(time.RFC3339, iso); err == nil {
			return iso
		}
	}
	if dateStr == "" {
		return ""
	}
	if timeStr == "" {
		timeStr = "00:00"
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04", dateStr+" "+timeStr, tz)
	if err != nil {
		return ""
	}
	return parsed.Format(time.RFC3339)
}

// buildPreviewRequest assembles the minimum request for a period-based preview.
func buildPreviewRequest(subscriptionID, periodStart, periodEnd, revenueDate string) *revenuepb.CreateRevenueWithLineItemsRequest {
	subID := subscriptionID
	req := &revenuepb.CreateRevenueWithLineItemsRequest{
		Data:           &revenuepb.Revenue{},
		SubscriptionId: &subID,
	}
	if periodStart != "" {
		ps := periodStart
		req.PeriodStart = &ps
	}
	if periodEnd != "" {
		pe := periodEnd
		req.PeriodEnd = &pe
	}
	if revenueDate != "" {
		rd := revenueDate
		req.RevenueDate = &rd
	}
	return req
}

// buildMilestoneRequest constructs the dry-run request for a milestone.
func buildMilestoneRequest(subscriptionID, eventID, revenueDate string) *revenuepb.CreateRevenueWithLineItemsRequest {
	subID := subscriptionID
	evID := eventID
	req := &revenuepb.CreateRevenueWithLineItemsRequest{
		Data:           &revenuepb.Revenue{},
		SubscriptionId: &subID,
		BillingEventId: &evID,
	}
	if revenueDate != "" {
		rd := revenueDate
		req.RevenueDate = &rd
	}
	return req
}

// parseCentavoDecimal parses an operator-typed decimal ("80000.00") into centavos.
func parseCentavoDecimal(s string) (int64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || f < 0 {
		return 0, false
	}
	return int64(f * 100), true
}

// formatCentavos formats a centavo amount as a 2-decimal string ("150000.00").
func formatCentavos(c int64) string {
	whole := c / 100
	frac := c % 100
	if frac < 0 {
		frac = -frac
	}
	return fmt.Sprintf("%d.%02d", whole, frac)
}

// statusLabel returns the localized status label for a BillingEvent.status value.
func statusLabel(s billingeventpb.BillingEventStatus, l centymo.SubscriptionMilestoneLabels) string {
	switch s {
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_READY:
		return l.StatusReady
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_BILLED:
		return l.StatusBilled
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_WAIVED:
		return l.StatusWaived
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_DEFERRED:
		return l.StatusDeferred
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_CANCELLED:
		return l.StatusCancelled
	default:
		return l.StatusPending
	}
}

// statusKey returns the lowercase status token used in data-testid attributes.
func statusKey(s billingeventpb.BillingEventStatus) string {
	switch s {
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_READY:
		return "ready"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_BILLED:
		return "billed"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_WAIVED:
		return "waived"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_DEFERRED:
		return "deferred"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_CANCELLED:
		return "cancelled"
	default:
		return "pending"
	}
}

// sequenceLabel returns a friendly label for the milestone select.
func sequenceLabel(ev *billingeventpb.BillingEvent) string {
	if ev == nil {
		return ""
	}
	if s := strings.TrimSpace(ev.GetSequenceLabel()); s != "" {
		return s
	}
	id := ev.GetId()
	if len(id) > 8 {
		return "Event " + id[len(id)-6:]
	}
	return "Event " + id
}

func boolPtr(b bool) *bool   { return &b }
func stringPtr(s string) *string { return &s }
