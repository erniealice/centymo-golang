package action

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// RecognizeFormLabels mirrors the lyngua-driven label set that the recognize
// drawer surfaces. Built from centymo.SubscriptionRecognizeLabels +
// centymo.SubscriptionInvoicesLabels at handler time so the template only
// touches a single ".Labels.X" path per string.
type RecognizeFormLabels struct {
	Title                 string
	Subtitle              string
	ContextSection        string
	ClientLabel           string
	PlanLabel             string
	QuantityLabel         string
	PeriodSection         string
	PeriodStart           string
	PeriodEnd             string
	RevenueDate           string
	LineItemsSection      string
	ColumnDescription     string
	ColumnUnitPrice       string
	ColumnQuantity        string
	ColumnLineTotal       string
	ColumnTreatment       string
	TotalLabel            string
	RemoveLine            string
	TreatmentRecurring    string
	TreatmentFirstCycle   string
	TreatmentUsageBased   string
	TreatmentOneTime      string
	NotesLabel            string
	NotesPlaceholder      string
	Generate              string
	Cancel                string
	Timezone              string
	StartDateInfo         string
	EndDateInfo           string
	StartTimeInfo         string
	EndTimeInfo           string

	// Blocking error banners.
	CurrencyMismatchError   string
	IdempotencyError        string
	IdempotencyExistingLink string
	NoLinesError            string
}

// RecognizePreviewLine is the row shape consumed by the drawer template.
// Mirrors revenuepb.PreviewLineItem but exposes only fields the template
// actually renders.
type RecognizePreviewLine struct {
	ProductPricePlanID string
	Description        string
	UnitPrice          int64
	Quantity           float64
	TotalPrice         int64
	Currency           string
	Treatment          string
	TreatmentLabel     string
}

// RecognizeFormData is the template data for the recognize-revenue drawer.
type RecognizeFormData struct {
	FormAction       string
	SubscriptionID   string
	SubscriptionName string
	ClientLabel      string
	PlanLabel        string
	Quantity         int32
	Currency         string

	// Period (date + time grid, IANA tz aware — same pattern as the standard
	// subscription drawer).
	PeriodStartDate string
	PeriodStartTime string
	PeriodStartISO  string
	PeriodEndDate   string
	PeriodEndTime   string
	PeriodEndISO    string
	DefaultTZ       string

	// Revenue date (single date input).
	RevenueDate string

	// Notes — auto-prefixed with the period marker.
	Notes string

	// Line items preview.
	PreviewLines []RecognizePreviewLine
	TotalAmount  int64

	// Blocking-error state.
	CurrencyMismatch        bool
	ClientCurrency          string
	PlanCurrency            string
	IdempotencyConflict     bool
	ConflictingRevenueID    string
	ConflictingRevenueURL   string
	NoLinesToInvoice        bool

	// Non-blocking warnings (e.g. usage-based skipped notice).
	Warnings []string

	Labels       RecognizeFormLabels
	CommonLabels any
}

// recognizeFormLabels builds the label bag from the typed centymo labels.
func recognizeFormLabels(l centymo.SubscriptionLabels) RecognizeFormLabels {
	return RecognizeFormLabels{
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
	}
}

// NewRecognizeAction creates the subscription "Recognize Revenue" view.
//
//	GET  → renders the drawer with a dry_run preview from the use case.
//	POST → calls the use case for real, returning HTMXSuccess on success so
//	       the invoices table refreshes inline (no redirect).
//
// The use case enforces the currency mismatch + idempotency hard blocks; on
// those errors the drawer re-renders with the appropriate error banner and
// the Generate button disabled (template handles that based on flags).
func NewRecognizeAction(deps *Deps) view.View {
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
		// determination. Done once per call (GET preview + POST submit both
		// need it).
		sub, client, pricePlan := loadRecognizeContext(ctx, deps, subscriptionID)
		if sub == nil {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}

		tz := pyezatypes.LocationFromContext(ctx)

		switch viewCtx.Request.Method {
		case http.MethodGet:
			return renderRecognizeDrawer(ctx, deps, viewCtx, subscriptionID, sub, client, pricePlan, tz)
		case http.MethodPost:
			return submitRecognizeDrawer(ctx, deps, viewCtx, subscriptionID, sub, client, pricePlan, tz)
		default:
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
	})
}

// loadRecognizeContext resolves the (subscription, client, pricePlan) triple
// that the drawer header uses. Best-effort — nil values fall through to empty
// labels in the template.
func loadRecognizeContext(
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

// renderRecognizeDrawer handles GET — runs a dry_run through the use case and
// renders the drawer with the preview.
func renderRecognizeDrawer(
	ctx context.Context, deps *Deps, viewCtx *view.ViewContext,
	subscriptionID string,
	sub *subscriptionpb.Subscription,
	client *clientpb.Client,
	pricePlan *priceplanpb.PricePlan,
	tz *time.Location,
) view.ViewResult {
	defaultPeriodStart, defaultPeriodEnd := defaultPeriodBoundsForSubscription(sub, pricePlan, tz)
	revenueDate := time.Now().In(tz).Format(pyezatypes.DateInputLayout)

	data := buildBaseFormData(deps, subscriptionID, sub, client, pricePlan, tz,
		defaultPeriodStart, defaultPeriodEnd, revenueDate)

	// Run the use case in dry_run mode to compute the preview. If unavailable,
	// the drawer still renders without lines (operator sees an empty preview).
	if deps.RecognizeRevenueFromSubscription != nil {
		req := buildPreviewRequest(subscriptionID, defaultPeriodStart, defaultPeriodEnd, revenueDate)
		req.DryRun = boolPtr(true)
		resp, err := deps.RecognizeRevenueFromSubscription(ctx, req)
		applyResponseToFormData(&data, resp, err, deps, client, pricePlan)
	}

	return view.OK("subscription-recognize-drawer-form", &data)
}

// submitRecognizeDrawer handles POST — runs the use case for real and (on
// success) returns the HTMXSuccess header bundle that closes the drawer and
// triggers the invoices table refresh.
func submitRecognizeDrawer(
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

	if deps.RecognizeRevenueFromSubscription == nil {
		return centymo.HTMXError("recognize-revenue use case not configured")
	}

	req := buildPreviewRequest(subscriptionID, periodStart, periodEnd, revenueDate)
	req.DryRun = boolPtr(false)
	if notes != "" {
		req.Data = &revenuepb.Revenue{Notes: stringPtr(notes)}
	}

	resp, err := deps.RecognizeRevenueFromSubscription(ctx, req)
	if err != nil {
		// The use case returned a hard block (currency mismatch or idempotency
		// conflict) or an internal failure. Re-render the drawer with the
		// banner state so the operator sees the cause inline.
		data := buildBaseFormData(deps, subscriptionID, sub, client, pricePlan, tz,
			periodStart, periodEnd, revenueDate)
		applyResponseToFormData(&data, resp, err, deps, client, pricePlan)
		// Still need a preview for the table — re-run dry_run so the operator
		// can adjust and retry.
		previewReq := buildPreviewRequest(subscriptionID, periodStart, periodEnd, revenueDate)
		previewReq.DryRun = boolPtr(true)
		if previewResp, _ := deps.RecognizeRevenueFromSubscription(ctx, previewReq); previewResp != nil {
			data.PreviewLines = convertPreviewLines(previewResp.GetPreviewLines(), pricePlan, deps.Labels)
			data.TotalAmount = sumPreview(data.PreviewLines)
		}
		log.Printf("Recognize revenue from subscription %s failed: %v", subscriptionID, err)
		return view.OK("subscription-recognize-drawer-form", &data)
	}

	// Success — close the drawer and refresh the invoices table.
	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": `{"formSuccess":true,"refreshTable":"subscription-invoices-table","refresh-invoices":true}`,
		},
	}
}

// buildBaseFormData populates the form-data with the read-only context (no
// preview lines yet — those are layered on after the use case dry-run).
func buildBaseFormData(
	deps *Deps,
	subscriptionID string,
	sub *subscriptionpb.Subscription,
	client *clientpb.Client,
	pricePlan *priceplanpb.PricePlan,
	tz *time.Location,
	periodStart, periodEnd, revenueDate string,
) RecognizeFormData {
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

	formAction := resolveRecognizeURL(deps.Routes.RecognizeURL, subscriptionID)

	return RecognizeFormData{
		FormAction:       formAction,
		SubscriptionID:   subscriptionID,
		SubscriptionName: sub.GetName(),
		ClientLabel:      clientLabel,
		PlanLabel:        planLabel,
		Quantity:         sub.GetQuantity(),
		Currency:         planCurrency,
		PeriodStartDate:  startDate,
		PeriodStartTime:  startTime,
		PeriodStartISO:   startISO,
		PeriodEndDate:    endDate,
		PeriodEndTime:    endTime,
		PeriodEndISO:     endISO,
		DefaultTZ:        tz.String(),
		RevenueDate:      revenueDate,
		ClientCurrency:   clientCurrency,
		PlanCurrency:     planCurrency,
		Labels:           recognizeFormLabels(deps.Labels),
		CommonLabels:     nil, // injected by ViewAdapter
	}
}

// applyResponseToFormData overlays preview lines / warnings / error banners
// from a use-case response onto the form data.
func applyResponseToFormData(
	data *RecognizeFormData,
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
			// Build a link to the existing revenue's detail page (same template
			// the invoices table uses).
			data.ConflictingRevenueURL = strings.ReplaceAll(centymo.RevenueDetailURL, "{id}", cid)
		}
	}
	// Currency mismatch is a hard block — detect from raw client/plan compare
	// since the use case returns an error before populating the response.
	if client != nil && pricePlan != nil &&
		client.GetBillingCurrency() != "" &&
		pricePlan.GetBillingCurrency() != "" &&
		client.GetBillingCurrency() != pricePlan.GetBillingCurrency() {
		data.CurrencyMismatch = true
	}
	// no_lines_to_invoice is a hard block returned BEFORE the response is
	// populated — detect it via the error string. The translation key the use
	// case emits is revenue.errors.no_lines_to_invoice; deps.Labels carries
	// the resolved English/professional copy.
	if err != nil && strings.Contains(err.Error(), "no line items") {
		data.NoLinesToInvoice = true
	}
}

// buildPreviewRequest assembles a CreateRevenueWithLineItemsRequest with the
// minimum fields the use case needs to compute defaults.
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

// convertPreviewLines maps proto preview lines into the template's row shape,
// translating the treatment token into a localized badge label.
func convertPreviewLines(
	in []*revenuepb.PreviewLineItem,
	pricePlan *priceplanpb.PricePlan,
	labels centymo.SubscriptionLabels,
) []RecognizePreviewLine {
	currency := ""
	if pricePlan != nil {
		currency = pricePlan.GetBillingCurrency()
	}
	out := make([]RecognizePreviewLine, 0, len(in))
	for _, p := range in {
		lineCurrency := currency
		if p.GetCurrency() != "" {
			lineCurrency = p.GetCurrency()
		}
		out = append(out, RecognizePreviewLine{
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

func sumPreview(lines []RecognizePreviewLine) int64 {
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

// defaultPeriodBoundsForSubscription returns sane RFC3339 defaults for the
// drawer's period inputs. Mirrors plan §3.2 — uses the subscription window
// when present and falls back to today + 1 month.
func defaultPeriodBoundsForSubscription(
	sub *subscriptionpb.Subscription, pricePlan *priceplanpb.PricePlan, tz *time.Location,
) (string, string) {
	now := time.Now().In(tz)
	// Period start defaults to either subscription start (if known) or now.
	var start time.Time
	if dts := sub.GetDateTimeStart(); dts != nil && dts.IsValid() {
		start = dts.AsTime().In(tz)
	} else {
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tz)
	}
	// Period end defaults to start + cycle (or +1 month when cycle missing).
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
// Prefers the hidden ISO field set by JS (preserves the operator's tz offset);
// falls back to date+time-in-tz when JS is disabled.
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

// resolveRecognizeURL inlines the small {id} substitution so the package does
// not need to import pyeza-golang/route here.
func resolveRecognizeURL(template, id string) string {
	return strings.ReplaceAll(template, "{id}", id)
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
