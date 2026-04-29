package action

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/erniealice/pyeza-golang/route"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	billingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billing_event"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FormLabels holds i18n labels for the subscription drawer form template.
type FormLabels struct {
	Customer                  string
	CustomerPlaceholder       string
	Plan                      string
	PlanPlaceholder           string
	StartDate                 string
	EndDate                   string
	StartTime                 string
	EndTime                   string
	TimePlaceholder           string
	Timezone                  string
	Notes                     string
	NotesPlaceholder          string
	CustomerSearchPlaceholder string
	PlanSearchPlaceholder     string
	CustomerNoResults         string
	PlanNoResults             string
	Code                      string
	CodePlaceholder           string
	CustomerInfo              string
	PlanInfo                  string
	CodeInfo                  string
	StartDateInfo             string
	EndDateInfo               string
	StartTimeInfo             string
	EndTimeInfo               string
	NotesInfo                 string

	// 2026-04-27 plan-client-scope plan §5.1 / §7 — group headers in the
	// grouped Plan / PricePlan auto-complete picker.
	PlanGroupForClient string // "For {ClientName}" — pre-resolved with ClientName injected.
	PlanGroupGeneral   string
}

// PlanOptionGroup is one optgroup in the grouped Plan/PricePlan auto-complete
// on the subscription drawer (plan §5.1). Field name `GroupLabel` matches the
// pyeza auto-complete component's expected SelectOptionGroup shape — see
// templates/components/auto-complete.html.
type PlanOptionGroup struct {
	GroupLabel string              // group header
	Options    []map[string]string // {Value, Label} entries
}

// FormData is the template data for the subscription drawer form.
type FormData struct {
	FormAction      string
	IsEdit          bool
	ID              string
	Code            string
	ClientID        string
	PricePlanID     string
	// Date/Time form values, split for the two-row date+time grid.
	// Stored in the operator's display TZ (DefaultTZ) for the date/time inputs;
	// JS recombines + converts to UTC RFC 3339 for the hidden field.
	DateStartDate string
	DateStartTime string
	DateEndDate   string
	DateEndTime   string
	// Pre-computed RFC 3339 hidden values; JS overwrites on every change.
	DateStartISO string
	DateEndISO   string
	// DefaultTZ is the IANA name of the operator's display timezone, surfaced as
	// data-default-tz on the form for client-side recombination.
	DefaultTZ string
	Notes     string

	Clients         []map[string]string
	PricePlans      []map[string]string
	SearchClientURL string
	SearchPlanURL   string
	ClientLabel     string
	PlanLabel       string
	ClientLocked    bool
	// ClientBillingCurrency is the selected client's billing currency, passed to
	// the plan search URL so the grouped auto-complete only shows plans in that
	// currency. Empty = no currency filter.
	ClientBillingCurrency string

	// 2026-04-27 plan-client-scope plan §5 — grouped picker options. When
	// non-empty, the template renders the grouped variant instead of the
	// flat search auto-complete.
	PlanOptionGroups []PlanOptionGroup

	Labels                FormLabels
	CommonLabels          any
}

// Deps holds dependencies for subscription action handlers.
type Deps struct {
	Routes centymo.SubscriptionRoutes
	Labels centymo.SubscriptionLabels

	CreateSubscription  func(ctx context.Context, req *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error)
	ReadSubscription    func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
	// GetSubscriptionItemPageData returns the subscription with its joined
	// Client (+ User) and PricePlan (+ Plan) populated. Edit drawer uses it
	// to render the customer name (not the bare client_id) without depending
	// on a separate ListClients-and-iterate fallback.
	GetSubscriptionItemPageData func(ctx context.Context, req *subscriptionpb.GetSubscriptionItemPageDataRequest) (*subscriptionpb.GetSubscriptionItemPageDataResponse, error)
	UpdateSubscription  func(ctx context.Context, req *subscriptionpb.UpdateSubscriptionRequest) (*subscriptionpb.UpdateSubscriptionResponse, error)
	DeleteSubscription  func(ctx context.Context, req *subscriptionpb.DeleteSubscriptionRequest) (*subscriptionpb.DeleteSubscriptionResponse, error)
	ListClients         func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ListPlans           func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ReadPlan            func(ctx context.Context, req *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)
	SearchClientsByName func(ctx context.Context, req *clientpb.SearchClientsByNameRequest) (*clientpb.SearchClientsByNameResponse, error)
	SearchPlansByName   func(ctx context.Context, req *planpb.SearchPlansByNameRequest) (*planpb.SearchPlansByNameResponse, error)
	ListPricePlans      func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ReadPricePlan       func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	ListPriceSchedules  func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)

	// SetSubscriptionActive performs a raw DB update of the active field.
	// Required for set-status and bulk-set-status handlers.
	// Uses raw update (not proto) because proto3 omits bool=false on serialization.
	SetSubscriptionActive func(ctx context.Context, id string, active bool) error

	// GetInUseIDs checks whether subscription IDs are referenced by dependent records.
	// Used by the bulk-delete handler to skip in-use rows.
	GetInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// RecognizeRevenueFromSubscription invokes the espyna use case that
	// materializes a Revenue + N RevenueLineItems for a billing period.
	// Set when the centymo block wiring threads the use case through. Used
	// by NewRecognizeAction (drawer GET dry-run + POST commit) and by the
	// existing manual revenue-add flow's auto-populate path (skip_header=true).
	RecognizeRevenueFromSubscription func(ctx context.Context, req *revenuepb.CreateRevenueWithLineItemsRequest) (*revenuepb.CreateRevenueWithLineItemsResponse, error)

	// CustomClientPriceScheduleLabelSuffix carries the lyngua-resolved suffix
	// appended to a client's name when constructing the default custom
	// PriceSchedule name (e.g. "Price Schedule" / "Rate Cards"). Read by the
	// customize handler; sourced from PriceScheduleLabels.Form by block.go.
	CustomClientPriceScheduleLabelSuffix string

	// CustomizePlanForClient invokes the espyna use case that clones the
	// source Plan + PricePlan into a client-scoped copy and (optionally)
	// repoints the subscription onto the new PricePlan. See plan §4.
	// Wired by the centymo block when the use case is available.
	CustomizePlanForClient func(ctx context.Context, req *CustomizePlanForClientRequest) (*CustomizePlanForClientResponse, error)

	// 2026-04-29 milestone-billing plan §5 / Phase D — BillingEvent operations
	// for the subscription Package tab Milestones section + recognize drawer
	// milestone select. nil-safe: when unset (no adapter registered), the
	// drawer falls back to the legacy non-milestone branches and the Package
	// tab milestone section is skipped.
	ListBillingEventsBySubscription func(ctx context.Context, req *billingeventpb.ListBillingEventsBySubscriptionRequest) (*billingeventpb.ListBillingEventsBySubscriptionResponse, error)
	SetBillingEventStatus           func(ctx context.Context, req *billingeventpb.SetBillingEventStatusRequest) (*billingeventpb.SetBillingEventStatusResponse, error)
}

// CustomizePlanForClientRequest mirrors the espyna use-case request shape
// (plan §4.1). Centymo handlers build this and pass it through Deps.
// The `derivedName` carries the per-tier "{Client.name} - {suffix}" label
// resolved on the centymo side from typed labels (plan §4.4.1 step 2-3).
type CustomizePlanForClientRequest struct {
	SourcePlanID      string
	SourcePricePlanID string
	ClientID          string
	SubscriptionID    string
	NewScheduleName   string
}

// CustomizePlanForClientResponse mirrors the espyna use-case response shape
// (plan §4.1). Only the fields centymo's POST handler needs are surfaced;
// extend if a future caller needs the cloned proto records.
type CustomizePlanForClientResponse struct {
	NewPlanID      string
	NewPricePlanID string
	NewScheduleID  string
	Reused         bool
}

func formLabels(l centymo.SubscriptionLabels) FormLabels {
	return FormLabels{
		Customer:                  l.Form.Customer,
		CustomerPlaceholder:       l.Form.CustomerPlaceholder,
		Plan:                      l.Form.Plan,
		PlanPlaceholder:           l.Form.PlanPlaceholder,
		StartDate:                 l.Form.StartDate,
		EndDate:                   l.Form.EndDate,
		StartTime:                 l.Form.StartTime,
		EndTime:                   l.Form.EndTime,
		TimePlaceholder:           l.Form.TimePlaceholder,
		Timezone:                  l.Form.Timezone,
		Notes:                     l.Form.Notes,
		NotesPlaceholder:          l.Form.NotesPlaceholder,
		CustomerSearchPlaceholder: l.Form.CustomerSearchPlaceholder,
		PlanSearchPlaceholder:     l.Form.PlanSearchPlaceholder,
		CustomerNoResults:         l.Form.CustomerNoResults,
		PlanNoResults:             l.Form.PlanNoResults,
		Code:                      l.Form.Code,
		CodePlaceholder:           l.Form.CodePlaceholder,
		CustomerInfo:              l.Form.CustomerInfo,
		PlanInfo:                  l.Form.PlanInfo,
		CodeInfo:                  l.Form.CodeInfo,
		StartDateInfo:             l.Form.StartDateInfo,
		EndDateInfo:               l.Form.EndDateInfo,
		StartTimeInfo:             l.Form.StartTimeInfo,
		EndTimeInfo:               l.Form.EndTimeInfo,
		NotesInfo:                 l.Form.NotesInfo,
		PlanGroupForClient:        l.Form.PlanGroupForClient,
		PlanGroupGeneral:          l.Form.PlanGroupGeneral,
	}
}

// resolvePlanGroupForClientLabel renders the {{.ClientName}}-templated
// "For {ClientName}" group header. Falls back gracefully when the label
// has no template directive or the client name is empty.
func resolvePlanGroupForClientLabel(template, clientName string) string {
	if clientName == "" {
		return template
	}
	return strings.ReplaceAll(template, "{{.ClientName}}", clientName)
}

// generateCode returns a random 7-character uppercase alphanumeric code,
// using chars that are visually unambiguous (no O, I, 0, 1).
func generateCode() string {
	const chars = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	b := make([]byte, 7)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// loadClientOptions fetches the client list and converts to select options.
func loadClientOptions(ctx context.Context, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)) []map[string]string {
	if listClients == nil {
		return nil
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		log.Printf("Failed to load clients for dropdown: %v", err)
		return nil
	}
	var options []map[string]string
	for _, c := range resp.GetData() {
		label := c.GetId()
		if u := c.GetUser(); u != nil {
			first := u.GetFirstName()
			last := u.GetLastName()
			if first != "" || last != "" {
				label = first + " " + last
			}
		}
		options = append(options, map[string]string{
			"Value": c.GetId(),
			"Label": label,
		})
	}
	return options
}

// loadPlanOptions fetches the plan list and converts to select options.
func loadPlanOptions(ctx context.Context, listPlans func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)) []map[string]string {
	if listPlans == nil {
		return nil
	}
	resp, err := listPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		log.Printf("Failed to load plans for dropdown: %v", err)
		return nil
	}
	var options []map[string]string
	for _, p := range resp.GetData() {
		if !p.GetActive() {
			continue
		}
		options = append(options, map[string]string{
			"Value": p.GetId(),
			"Label": p.GetName(),
		})
	}
	return options
}

// loadPlanOptionGroups builds the grouped Plan picker for the subscription
// drawer per plan §5.1. Group order: client-scoped first ("For {ClientName}"),
// general ("General packages") second. Empty groups are omitted.
//
// The list is filtered post-fetch in Go because TypedFilter doesn't yet
// expose a NULL/NOT-NULL primitive on string fields. Volume is small.
func loadPlanOptionGroups(ctx context.Context, listPlans func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error), clientID, clientName string, l FormLabels) []PlanOptionGroup {
	if listPlans == nil {
		return nil
	}
	resp, err := listPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		log.Printf("Failed to load plans for grouped picker: %v", err)
		return nil
	}

	var clientPlans, masterPlans []map[string]string
	for _, p := range resp.GetData() {
		if !p.GetActive() {
			continue
		}
		entry := map[string]string{"Value": p.GetId(), "Label": p.GetName()}
		switch cid := p.GetClientId(); {
		case cid == "":
			masterPlans = append(masterPlans, entry)
		case clientID != "" && cid == clientID:
			clientPlans = append(clientPlans, entry)
		}
	}

	var groups []PlanOptionGroup
	if len(clientPlans) > 0 {
		groups = append(groups, PlanOptionGroup{
			GroupLabel: resolvePlanGroupForClientLabel(l.PlanGroupForClient, clientName),
			Options:    clientPlans,
		})
	}
	if len(masterPlans) > 0 {
		groups = append(groups, PlanOptionGroup{
			GroupLabel: l.PlanGroupGeneral,
			Options:    masterPlans,
		})
	}
	return groups
}

// loadPricePlanOptionGroups is the same shape as loadPlanOptionGroups but
// keyed off PricePlan.client_id. Used by the subscription edit drawer's
// PricePlan picker.
func loadPricePlanOptionGroups(ctx context.Context, listPricePlans func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error), clientID, clientName string, l FormLabels) []PlanOptionGroup {
	if listPricePlans == nil {
		return nil
	}
	resp, err := listPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
	if err != nil {
		log.Printf("Failed to load price plans for grouped picker: %v", err)
		return nil
	}

	var clientPP, masterPP []map[string]string
	for _, pp := range resp.GetData() {
		if !pp.GetActive() {
			continue
		}
		label := pp.GetName()
		if label == "" {
			if pl := pp.GetPlan(); pl != nil {
				label = pl.GetName()
			}
			if label == "" {
				label = pp.GetId()
			}
		}
		entry := map[string]string{"Value": pp.GetId(), "Label": label}
		switch cid := pp.GetClientId(); {
		case cid == "":
			masterPP = append(masterPP, entry)
		case clientID != "" && cid == clientID:
			clientPP = append(clientPP, entry)
		}
	}

	var groups []PlanOptionGroup
	if len(clientPP) > 0 {
		groups = append(groups, PlanOptionGroup{
			GroupLabel: resolvePlanGroupForClientLabel(l.PlanGroupForClient, clientName),
			Options:    clientPP,
		})
	}
	if len(masterPP) > 0 {
		groups = append(groups, PlanOptionGroup{
			GroupLabel: l.PlanGroupGeneral,
			Options:    masterPP,
		})
	}
	return groups
}

// resolveClientBillingCurrency finds the billing_currency for a client by ID.
// Returns empty string when the client has no billing_currency set (caller should
// treat empty as "no currency filter" in the drawer's plan search).
func resolveClientBillingCurrency(ctx context.Context, clientID string, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)) string {
	if clientID == "" || listClients == nil {
		return ""
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		return ""
	}
	for _, c := range resp.GetData() {
		if c.GetId() == clientID {
			return c.GetBillingCurrency()
		}
	}
	return ""
}

// resolveClientLabel finds the display name for a client by ID.
func resolveClientLabel(ctx context.Context, clientID string, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)) string {
	if clientID == "" || listClients == nil {
		return ""
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		return clientID
	}
	for _, c := range resp.GetData() {
		if c.GetId() == clientID {
			if cn := c.GetName(); cn != "" {
				return cn
			}
			if u := c.GetUser(); u != nil {
				first := u.GetFirstName()
				last := u.GetLastName()
				if first != "" || last != "" {
					return strings.TrimSpace(first + " " + last)
				}
			}
			return clientID
		}
	}
	return clientID
}

// resolvePlanLabel finds the display name for a plan by ID.
func resolvePlanLabel(ctx context.Context, planID string, listPlans func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)) string {
	if planID == "" || listPlans == nil {
		return ""
	}
	resp, err := listPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return planID
	}
	for _, p := range resp.GetData() {
		if p.GetId() == planID {
			return p.GetName()
		}
	}
	return planID
}

// resolvePricePlanName looks up a PricePlan by ID and returns its display name.
// Prefers a single-row ReadPricePlan lookup over a full list scan.
// Falls back to the legacy ListPlans scan only when ReadPricePlan is nil or errors.
func resolvePricePlanName(ctx context.Context, pricePlanID string, deps *Deps) string {
	if pricePlanID == "" {
		return ""
	}
	// Single-row lookup is preferred over ListPricePlans for one-shot resolution.
	if deps.ReadPricePlan != nil {
		if resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
			Data: &priceplanpb.PricePlan{Id: pricePlanID},
		}); err == nil && resp != nil && len(resp.GetData()) > 0 {
			pp := resp.GetData()[0]
			if name := pp.GetName(); name != "" {
				return name
			}
			if pl := pp.GetPlan(); pl != nil && pl.GetName() != "" {
				return pl.GetName()
			}
			if deps.ReadPlan != nil && pp.GetPlanId() != "" {
				planID := pp.GetPlanId()
				if rr, err := deps.ReadPlan(ctx, &planpb.ReadPlanRequest{Data: &planpb.Plan{Id: &planID}}); err == nil && len(rr.GetData()) > 0 {
					if n := rr.GetData()[0].GetName(); n != "" {
						return n
					}
				}
			}
			return pricePlanID
		}
	}
	// Last-resort fallback: legacy Plan list (handles cases where the
	// submitted ID is actually a plan_id rather than a price_plan_id).
	if deps.ListPlans != nil {
		if name := resolvePlanLabel(ctx, pricePlanID, deps.ListPlans); name != "" && name != pricePlanID {
			return name
		}
	}
	return pricePlanID
}

// splitTimestampForInputs renders ts in tz as a (date, time, RFC3339) triple
// suitable for the drawer's two-input grid + hidden ISO field. Nil ts → empties.
func splitTimestampForInputs(ts *timestamppb.Timestamp, tz *time.Location) (date, t, iso string) {
	if ts == nil || !ts.IsValid() {
		return "", "", ""
	}
	moment := ts.AsTime().In(tz)
	return moment.Format(pyezatypes.DateInputLayout), moment.Format(pyezatypes.TimeInputLayout), moment.Format(time.RFC3339)
}

// parseFormDateTime combines a date input ("2026-04-17"), a time input ("09:00"),
// and an explicit RFC3339 ISO string (set by JS with the chosen TZ offset) into
// a *timestamppb.Timestamp. The hidden ISO wins when present so the operator's
// chosen offset is preserved exactly. Falls back to date+time-in-tz when JS is
// disabled or the hidden field is empty. Empty all → nil.
//
// 2026-04-28 date+time field plan §4 — when no time is provided, isEnd
// switches the default between 00:00:00 (start) and 23:59:59 (end) so that an
// "end" date without a time still includes the full day.
func parseFormDateTime(date, t, iso string, tz *time.Location, isEnd bool) *timestamppb.Timestamp {
	if iso != "" {
		if parsed, err := time.Parse(time.RFC3339, iso); err == nil {
			return timestamppb.New(parsed.UTC())
		}
	}
	if date == "" {
		return nil
	}
	if t == "" {
		if isEnd {
			t = "23:59:59"
		} else {
			t = "00:00:00"
		}
	} else if len(t) == 5 {
		// Browser time inputs default to HH:MM precision; pad seconds.
		t = t + ":00"
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", date+" "+t, tz)
	if err != nil {
		return nil
	}
	return timestamppb.New(parsed.UTC())
}

// NewAddAction creates the subscription add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			clientID := viewCtx.Request.URL.Query().Get("client_id")
			clientName := viewCtx.Request.URL.Query().Get("client_name")
			clientBillingCurrency := viewCtx.Request.URL.Query().Get("billing_currency")
			clientLocked := clientID != ""

			tz := pyezatypes.LocationFromContext(ctx)
			// Default new engagement to "today, 00:00" in the operator's TZ.
			today := time.Now().In(tz)
			defaultDate := today.Format(pyezatypes.DateInputLayout)
			defaultISO := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, tz).Format(time.RFC3339)
			labels := formLabels(deps.Labels)
			return view.OK("subscription-drawer-form", &FormData{
				FormAction:            deps.Routes.AddURL,
				SearchClientURL:       deps.Routes.SearchClientURL,
				SearchPlanURL:         deps.Routes.SearchPlanURL,
				ClientID:              clientID,
				ClientLabel:           clientName,
				ClientLocked:          clientLocked,
				ClientBillingCurrency: clientBillingCurrency,
				Code:                  generateCode(),
				DateStartDate:         defaultDate,
				DateStartISO:          defaultISO,
				DefaultTZ:             tz.String(),
				PlanOptionGroups:      loadPricePlanOptionGroups(ctx, deps.ListPricePlans, clientID, clientName, labels),
				Labels:                labels,
				CommonLabels:          nil, // injected by ViewAdapter
			})
		}

		// POST — create subscription
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		tz := pyezatypes.LocationFromContext(ctx)
		dateTimeStart := parseFormDateTime(
			r.FormValue("date_start_date"),
			r.FormValue("date_start_time"),
			r.FormValue("date_time_start_iso"),
			tz,
			false,
		)
		dateTimeEnd := parseFormDateTime(
			r.FormValue("date_end_date"),
			r.FormValue("date_end_time"),
			r.FormValue("date_time_end_iso"),
			tz,
			true,
		)

		pricePlanID := r.FormValue("price_plan_id")

		code := r.FormValue("code")
		if code == "" {
			code = generateCode()
		}

		// Resolve plan name for auto-generated subscription name. The drawer
		// submits a price_plan_id, so look up the PricePlan (not the Plan).
		planName := resolvePricePlanName(ctx, pricePlanID, deps)
		name := planName
		if code != "" {
			name = planName + " [" + code + "]"
		}

		resp, err := deps.CreateSubscription(ctx, &subscriptionpb.CreateSubscriptionRequest{
			Data: &subscriptionpb.Subscription{
				Name:          name,
				ClientId:      r.FormValue("client_id"),
				PricePlanId:   pricePlanID,
				Code:          strPtr(code),
				DateTimeStart: dateTimeStart,
				DateTimeEnd:   dateTimeEnd,
				Active:        true,
			},
		})
		if err != nil {
			log.Printf("Failed to create subscription: %v", err)
			return centymo.HTMXError(err.Error())
		}

		_ = resp
		return centymo.HTMXSuccess("subscriptions-table")
	})
}

// NewEditAction creates the subscription edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			// Prefer the joined item-page-data path so Client (+ User) is populated
			// without a second ListClients-and-iterate roundtrip. Falls back to
			// ReadSubscription only if the dep is unwired.
			var record *subscriptionpb.Subscription
			if deps.GetSubscriptionItemPageData != nil {
				resp, err := deps.GetSubscriptionItemPageData(ctx, &subscriptionpb.GetSubscriptionItemPageDataRequest{
					SubscriptionId: id,
				})
				if err != nil || resp == nil || resp.GetSubscription() == nil {
					log.Printf("Failed to read subscription %s: %v", id, err)
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetSubscription()
			} else {
				readResp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
					Data: &subscriptionpb.Subscription{Id: id},
				})
				if err != nil {
					log.Printf("Failed to read subscription %s: %v", id, err)
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				readData := readResp.GetData()
				if len(readData) == 0 {
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = readData[0]
			}

			// Prefer the joined client (populated by GetSubscriptionItemPageData);
			// fall back to the ListClients lookup for the legacy ReadSubscription path.
			clientLabel := ""
			if c := record.GetClient(); c != nil {
				if name := c.GetName(); name != "" {
					clientLabel = name
				} else if u := c.GetUser(); u != nil {
					clientLabel = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
				}
			}
			if clientLabel == "" {
				clientLabel = resolveClientLabel(ctx, record.GetClientId(), deps.ListClients)
			}
			clientBillingCurrency := ""
			if c := record.GetClient(); c != nil {
				clientBillingCurrency = c.GetBillingCurrency()
			}
			if clientBillingCurrency == "" {
				clientBillingCurrency = resolveClientBillingCurrency(ctx, record.GetClientId(), deps.ListClients)
			}
			// PricePlanID, not a plan_id — resolve via PricePlan so the selected
			// label matches the autocomplete dropdown's display.
			planLabel := resolvePricePlanName(ctx, record.GetPricePlanId(), deps)

			// Lock client field when opened from client detail page
			clientLocked := viewCtx.Request.URL.Query().Get("client_id") != ""

			tz := pyezatypes.LocationFromContext(ctx)
			startDate, startTime, startISO := splitTimestampForInputs(record.GetDateTimeStart(), tz)
			endDate, endTime, endISO := splitTimestampForInputs(record.GetDateTimeEnd(), tz)

			labels := formLabels(deps.Labels)
			return view.OK("subscription-drawer-form", &FormData{
				FormAction:            route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                true,
				ID:                    id,
				Code:                  record.GetCode(),
				ClientID:              record.GetClientId(),
				PricePlanID:           record.GetPricePlanId(),
				DateStartDate:         startDate,
				DateStartTime:         startTime,
				DateStartISO:          startISO,
				DateEndDate:           endDate,
				DateEndTime:           endTime,
				DateEndISO:            endISO,
				DefaultTZ:             tz.String(),
				SearchClientURL:       deps.Routes.SearchClientURL,
				SearchPlanURL:         deps.Routes.SearchPlanURL,
				ClientLabel:           clientLabel,
				ClientLocked:          clientLocked,
				ClientBillingCurrency: clientBillingCurrency,
				PlanLabel:             planLabel,
				PlanOptionGroups:      loadPricePlanOptionGroups(ctx, deps.ListPricePlans, record.GetClientId(), clientLabel, labels),
				Labels:                labels,
				CommonLabels:          nil, // injected by ViewAdapter
			})
		}

		// POST — update subscription
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		tz := pyezatypes.LocationFromContext(ctx)
		dateTimeStart := parseFormDateTime(
			r.FormValue("date_start_date"),
			r.FormValue("date_start_time"),
			r.FormValue("date_time_start_iso"),
			tz,
			false,
		)
		dateTimeEnd := parseFormDateTime(
			r.FormValue("date_end_date"),
			r.FormValue("date_end_time"),
			r.FormValue("date_time_end_iso"),
			tz,
			true,
		)

		pricePlanID := r.FormValue("price_plan_id")
		if pricePlanID == "" {
			pricePlanID = r.FormValue("plan_id")
		}

		code := r.FormValue("code")
		if code == "" {
			code = generateCode()
		}

		// Resolve plan name for auto-generated subscription name. The drawer
		// submits a price_plan_id, so look up the PricePlan (not the Plan).
		planName := resolvePricePlanName(ctx, pricePlanID, deps)
		name := planName
		if code != "" {
			name = planName + " [" + code + "]"
		}

		_, err := deps.UpdateSubscription(ctx, &subscriptionpb.UpdateSubscriptionRequest{
			Data: &subscriptionpb.Subscription{
				Id:            id,
				Name:          name,
				ClientId:      r.FormValue("client_id"),
				PricePlanId:   pricePlanID,
				Code:          strPtr(code),
				DateTimeStart: dateTimeStart,
				DateTimeEnd:   dateTimeEnd,
			},
		})
		if err != nil {
			log.Printf("Failed to update subscription %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("subscriptions-table")
	})
}

// NewDeleteAction creates the subscription delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeleteSubscription(ctx, &subscriptionpb.DeleteSubscriptionRequest{
			Data: &subscriptionpb.Subscription{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete subscription %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("subscriptions-table")
	})
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}

