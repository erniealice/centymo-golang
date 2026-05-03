package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/price_plan/form"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// PricePlanDeps holds dependencies for price plan action handlers.
type PricePlanDeps struct {
	Routes centymo.PlanRoutes
	// Labels carries plan-level strings (errors, actions). Form-field labels
	// live on PricePlanLabels.Form (sourced from lyngua price_plan.json →
	// price_plan.form, which is the single source for the drawer).
	Labels              centymo.PlanLabels
	PricePlanLabels     centymo.PricePlanLabels
	// PriceScheduleLabels surfaces the customClientPriceScheduleLabelSuffix
	// used to derive the readonly Schedule label when the parent Plan is
	// client-scoped (plan §6.7). Optional — when zero-value, the helper
	// falls back to the proto-generic "Price Schedule".
	PriceScheduleLabels centymo.PriceScheduleLabels
	CommonLabels        pyeza.CommonLabels
	CreatePricePlan    func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)
	ReadPricePlan      func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	UpdatePricePlan    func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan    func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)
	ListPriceSchedules func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)

	// ReadPlan resolves the parent plan's name for display in the locked
	// "Package" field on the drawer, plus its client_id for the
	// client-scope schedule lock (plan §6.7).
	ReadPlan func(ctx context.Context, req *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)

	// ListClients resolves a client_id → display name for the readonly
	// schedule label and the lock tooltip (plan §6.7). Optional — when
	// nil, the lock falls back to a label derived from the client_id alone.
	ListClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)

	// ListLocations resolves price_schedule.location_id → location.name for
	// the form-hint below the rate-card auto-complete.
	ListLocations func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)

	// Reference checker: PricePlans in use by active subscriptions render
	// the drawer's Pricing section as read-only via InUse + LockMessage.
	GetPricePlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// Auto-seed ProductPricePlan rows on create (mirror of the
	// PriceSchedule-side behavior). All optional — when nil, auto-seed skips.
	ListProducts           func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ListProductPlans       func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	CreateProductPricePlan func(ctx context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)
	ListProductPricePlans  func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
}

// loadScheduleOptions fetches active price schedules as form.Option entries.
// Each option's Description carries the resolved location name so the
// drawer's rate-card auto-complete can render a location hint right below
// the field — dynamically updating as the user switches selection.
func loadScheduleOptions(ctx context.Context, deps *PricePlanDeps, hintPrefix string) []form.Option {
	if deps.ListPriceSchedules == nil {
		return nil
	}
	resp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
	if err != nil {
		log.Printf("Failed to load price schedules for price plan form: %v", err)
		return nil
	}

	locationNames := map[string]string{}
	if deps.ListLocations != nil {
		locResp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
		if err == nil {
			for _, l := range locResp.GetData() {
				locationNames[l.GetId()] = l.GetName()
			}
		}
	}

	var options []form.Option
	for _, s := range resp.GetData() {
		if !s.GetActive() {
			continue
		}
		description := ""
		if locID := s.GetLocationId(); locID != "" {
			if n := locationNames[locID]; n != "" {
				description = hintPrefix + n
			}
		}
		options = append(options, form.Option{
			ID:          s.GetId(),
			Name:        s.GetName(),
			Description: description,
		})
	}
	return options
}

// scheduleLocationHintPrefix is the literal prefix used when rendering the
// location hint under the rate-card auto-complete. Kept in sync with
// form.Labels.LocationHintPrefix so the hint reads e.g. "Location: Manila".
const scheduleLocationHintPrefix = "Location: "

// readPlan returns the parent plan record (for both name and client_id), or
// nil when unwired or missing. Single read shared by add + edit handlers.
func readPlan(ctx context.Context, deps *PricePlanDeps, planID string) *planpb.Plan {
	if deps.ReadPlan == nil || planID == "" {
		return nil
	}
	resp, err := deps.ReadPlan(ctx, &planpb.ReadPlanRequest{Data: &planpb.Plan{Id: &planID}})
	if err != nil || len(resp.GetData()) == 0 {
		return nil
	}
	return resp.GetData()[0]
}

// resolvePricePlanClientNameAndTz returns the display name AND the
// *time.Location anchoring derived-name suffixes. Tz comes from
// Client.User.Timezone (the only client-side tz the proto carries),
// falling back to fallbackTz when missing or invalid. Bad IANA names log
// and fall through. fallbackTz must not be nil.
func resolvePricePlanClientNameAndTz(ctx context.Context, clientID string, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error), fallbackTz *time.Location) (string, *time.Location) {
	if clientID == "" || listClients == nil {
		return "", fallbackTz
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		return clientID, fallbackTz
	}
	for _, c := range resp.GetData() {
		if c.GetId() != clientID {
			continue
		}
		tz := fallbackTz
		if u := c.GetUser(); u != nil {
			if tzName := strings.TrimSpace(u.GetTimezone()); tzName != "" {
				if loc, err := time.LoadLocation(tzName); err == nil {
					tz = loc
				} else {
					log.Printf("invalid client timezone %q for client %s: %v; falling back to request tz", tzName, clientID, err)
				}
			}
		}
		if name := c.GetName(); name != "" {
			return name, tz
		}
		if u := c.GetUser(); u != nil {
			full := strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			if full != "" {
				return full, tz
			}
		}
		return clientID, tz
	}
	return clientID, fallbackTz
}

// findClientPriceSchedule searches the workspace's active price schedules
// for the one bound to the supplied client_id. Returns the first match's
// full row, or nil when none exists (the use case will auto-create on save).
// Returning the row (not just the ID) lets the caller render the schedule's
// actual saved Name in the drawer's readonly label, instead of re-deriving
// "{client} - {suffix}" — the two can disagree if an operator has manually
// renamed the schedule.
func findClientPriceSchedule(ctx context.Context, deps *PricePlanDeps, clientID string) *priceschedulepb.PriceSchedule {
	if clientID == "" || deps.ListPriceSchedules == nil {
		return nil
	}
	resp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
	if err != nil {
		log.Printf("Failed to list price schedules for client schedule lock: %v", err)
		return nil
	}
	for _, s := range resp.GetData() {
		if !s.GetActive() {
			continue
		}
		if s.GetClientId() == clientID {
			return s
		}
	}
	return nil
}

// resolveScheduleLock computes the (mode, scheduleID, scheduleLabel,
// clientName) tuple for the price-schedule field on the PricePlan drawer.
// When the parent plan carries a client_id, the field collapses to readonly:
//   - if a client schedule already exists, the label is the schedule's
//     saved Name (so the drawer never disagrees with the persisted row);
//   - if none exists yet, the label is a timestamped preview built via
//     pyezatypes.AppendTimestamp, in the client's tz when set, otherwise
//     the request tz. The base mirrors espyna's
//     applyClientScopedScheduleRule (parentPlan.Name → "{client} - {suffix}"
//     → suffix-only) so the persisted name will share the same shape.
//     Render-time vs save-time timestamps will differ by a few seconds —
//     accepted.
//
// When the plan is master, returns ("picker", "", "", "") so the caller
// falls through to the standard auto-complete branch.
func resolveScheduleLock(ctx context.Context, deps *PricePlanDeps, plan *planpb.Plan) (mode, scheduleID, scheduleLabel, clientName string) {
	if plan == nil {
		return "picker", "", "", ""
	}
	clientID := plan.GetClientId()
	if clientID == "" {
		return "picker", "", "", ""
	}
	requestTz := pyezatypes.LocationFromContext(ctx)
	resolvedName, tz := resolvePricePlanClientNameAndTz(ctx, clientID, deps.ListClients, requestTz)
	clientName = resolvedName

	// Existing schedule → use the row's saved Name verbatim. If the
	// operator renamed it, the drawer reflects that — the alternative
	// (re-derive "{client} - {suffix}") silently lies.
	if existing := findClientPriceSchedule(ctx, deps, clientID); existing != nil {
		return "readonly", existing.GetId(), existing.GetName(), clientName
	}

	// No existing schedule → preview the name the use case will mint on
	// save. Base is ALWAYS "{client} - {suffix}" so the rate-card list
	// scans by client, regardless of how the parent Plan is named. Suffix
	// is lyngua-driven (general: "Price Schedule"; professional: "Rate
	// Cards"). Espyna's applyClientScopedScheduleRule mirrors this shape
	// using the suffix passed via context (see ctx plumbing in the POST
	// handler) so the persisted name matches.
	suffix := deps.PriceScheduleLabels.Form.CustomClientPriceScheduleLabelSuffix
	if suffix == "" {
		suffix = "Price Schedule"
	}
	var base string
	if cn := strings.TrimSpace(clientName); cn != "" {
		base = cn + " - " + suffix
	} else {
		base = suffix
	}
	preview := pyezatypes.AppendTimestamp(base, time.Now(), tz)
	return "readonly", "", preview, clientName
}

// buildScheduleAutoHint picks the right info-line for the readonly schedule
// field: "will create" when no client schedule exists yet, "will reuse" when
// one was found. Empty string when the field is in picker mode (parent Plan
// is master). Substitutes {{.ClientName}} server-side because html/template
// doesn't recursively render label values.
func buildScheduleAutoHint(formLabels centymo.PricePlanFormLabels, mode, scheduleID, clientName string) string {
	if mode != "readonly" {
		return ""
	}
	tmpl := formLabels.ScheduleAutoCreateHint
	if scheduleID != "" {
		tmpl = formLabels.ScheduleAutoReuseHint
	}
	if tmpl == "" || clientName == "" {
		return tmpl
	}
	return strings.ReplaceAll(tmpl, "{{.ClientName}}", clientName)
}

// applyScheduleLockTooltip substitutes the {{.ClientName}} placeholder in
// the lyngua-driven scheduleLockedTooltip with the resolved client name.
// html/template doesn't recursively render strings stored in label fields,
// so we do the swap server-side. When clientName is empty, returns the
// raw template (still readable in the unlikely case ListClients is unwired).
func applyScheduleLockTooltip(template, clientName string) string {
	if template == "" || clientName == "" {
		return template
	}
	return strings.ReplaceAll(template, "{{.ClientName}}", clientName)
}

// autoSeedProductPricePlans creates one ProductPricePlan row per product_plan
// linked to planID, copying price/currency from the underlying Product record.
// Mirrors price_schedule/detail/page.go's behavior so both contexts auto-seed
// the newly-created PricePlan's product-prices tab. Non-fatal on failure.
func autoSeedProductPricePlans(ctx context.Context, deps *PricePlanDeps, pricePlanID, planID, currency string) {
	if pricePlanID == "" || planID == "" {
		return
	}
	if deps.ListProductPlans == nil || deps.ListProducts == nil || deps.CreateProductPricePlan == nil {
		return
	}
	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
		Filters: &commonpb.FilterRequest{
			Logic: commonpb.FilterLogic_AND,
			Filters: []*commonpb.TypedFilter{{
				Field: "plan_id",
				FilterType: &commonpb.TypedFilter_StringFilter{
					StringFilter: &commonpb.StringFilter{
						Value: planID, Operator: commonpb.StringOperator_STRING_EQUALS,
					},
				},
			}},
		},
	})
	if err != nil || ppResp == nil {
		return
	}
	prodResp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		return
	}
	products := map[string]*productpb.Product{}
	for _, p := range prodResp.GetData() {
		if p != nil {
			products[p.GetId()] = p
		}
	}
	for _, pp := range ppResp.GetData() {
		if pp == nil {
			continue
		}
		productID := pp.GetProductId()
		productPlanID := pp.GetId()
		if productPlanID == "" {
			continue
		}
		var price int64
		rowCurrency := currency
		if rowCurrency == "" {
			rowCurrency = "PHP"
		}
		// Product.price is now optional (Model D). Only dereference when set.
		if prod := products[productID]; prod != nil {
			if prod.Price != nil {
				price = prod.GetPrice()
			}
			if c := prod.GetCurrency(); c != "" {
				rowCurrency = c
			}
		}
		if _, err := deps.CreateProductPricePlan(ctx, &productpriceplanpb.CreateProductPricePlanRequest{
			Data: &productpriceplanpb.ProductPricePlan{
				PricePlanId:     pricePlanID,
				ProductPlanId:   productPlanID, // Model D — FK to catalog line, not product
				BillingAmount:   price,
				BillingCurrency: rowCurrency,
				Active:          true,
			},
		}); err != nil {
			log.Printf("auto-seed product_price_plan failed for %s/%s: %v", pricePlanID, productPlanID, err)
		}
	}
}

// removed: buildScheduleAutoCompleteOptions / findScheduleLabel (replaced by
// form.BuildOptions + form.FindLabel), loadProductPlansForPlan (no longer
// rendered in the drawer — per-product prices are seeded automatically on
// create and edited from the PricePlan detail page).

// NewPricePlanAddAction creates the price plan add action (GET = form, POST = create).
// URL: /action/plan/{id}/price-plans/add
func NewPricePlanAddAction(deps *PricePlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		planID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			schedules := loadScheduleOptions(ctx, deps, scheduleLocationHintPrefix)
			formLabels := deps.PricePlanLabels.Form
			parentPlan := readPlan(ctx, deps, planID)
			planName := planID
			if parentPlan != nil {
				if n := parentPlan.GetName(); n != "" {
					planName = n
				}
			}
			scheduleMode, scheduleLockID, scheduleLockLabel, scheduleLockClient := resolveScheduleLock(ctx, deps, parentPlan)
			labels := form.LabelsFromPricePlan(formLabels)
			// Context-specific overrides; intentionally inlined (Decision 2 — only 2 fields).
			labels.ScheduleLockedTooltip = applyScheduleLockTooltip(labels.ScheduleLockedTooltip, scheduleLockClient)
			scheduleAutoHint := buildScheduleAutoHint(formLabels, scheduleMode, scheduleLockID, scheduleLockClient)
			return view.OK("price-plan-drawer-form", &form.Data{
				FormAction:          route.ResolveURL(deps.Routes.PricePlanAddURL, "id", planID),
				Context:             form.ContextPlan,
				PlanID:              planID,
				PlanName:            planName,
				Active:              true,
				Currency:            "PHP",
				DurationUnit:        "months",
				BillingKind:         "BILLING_KIND_RECURRING",
				AmountBasis:         "AMOUNT_BASIS_PER_CYCLE",
				BillingCycleUnit:    "month",
				TermUnit:            "month",
				ScheduleOptions:     form.BuildOptions(schedules, ""),
				DurationUnitOptions: form.BuildDurationUnitOptions(deps.CommonLabels),
				// Plan §6.7 — readonly schedule field when parent plan is
				// client-scoped. ScheduleID may be empty (auto-create on save).
				ScheduleFieldMode:        scheduleMode,
				ScheduleID:               scheduleLockID,
				ScheduleLabel:            scheduleLockLabel,
				ScheduleLockedClientName: scheduleLockClient,
				ScheduleAutoHint:         scheduleAutoHint,
				// 2026-04-30 cyclic-subscription-jobs plan §7.4 — surface the
				// parent Plan's cyclic flag so the drawer's
				// applyBasisOptionGuards() can disable MILESTONE.
				ParentPlanIsCyclic: parentPlan != nil && parentPlan.GetVisitsPerCycle() > 1,
				Labels:             labels,
				CommonLabels:       deps.CommonLabels,
			})
		}

		// POST — create price plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		amount := int64(0)
		if v, err := strconv.ParseFloat(r.FormValue("amount"), 64); err == nil {
			amount = int64(math.Round(v * 100))
		}

		currency := r.FormValue("currency")

		ppName := r.FormValue("name")
		ppDescription := r.FormValue("description")
		pp := &priceplanpb.PricePlan{
			PlanId:          planID,
			Name:            &ppName,
			Description:     &ppDescription,
			BillingAmount:   amount,
			BillingCurrency: currency,
			Active:          active,
		}
		// Phase 1 legacy dual-write — proto fields now optional; only set when present.
		if dvStr := r.FormValue("duration_value"); dvStr != "" {
			if v, err := strconv.ParseInt(dvStr, 10, 32); err == nil {
				dv32 := int32(v)
				pp.DurationValue = &dv32
			}
		}
		if du := r.FormValue("duration_unit"); du != "" {
			pp.DurationUnit = &du
		}
		if schedID := r.FormValue("price_schedule_id"); schedID != "" {
			pp.PriceScheduleId = &schedID
		}
		applyBillingFields(pp, r)

		// Plumb the lyngua-resolved suffix through context so espyna's
		// applyClientScopedScheduleRule can name an auto-created schedule
		// "{client} - {suffix} - {timestamp} {tz}" with the right
		// tier-specific noun (general: "Price Schedule"; professional:
		// "Rate Cards"). String key matches espyna's
		// ExtractClientScheduleSuffixFromContext reader.
		if suffix := deps.PriceScheduleLabels.Form.CustomClientPriceScheduleLabelSuffix; suffix != "" {
			ctx = context.WithValue(ctx, "clientScheduleSuffix", suffix)
		}

		createResp, err := deps.CreatePricePlan(ctx, &priceplanpb.CreatePricePlanRequest{
			Data: pp,
		})
		if err != nil {
			log.Printf("Failed to create price plan for plan %s: %v", planID, err)
			return centymo.HTMXError(err.Error())
		}

		// Auto-seed ProductPricePlan rows for the new PricePlan — mirrors the
		// schedule-side behavior so the drawer doesn't need to collect prices.
		if createResp != nil && len(createResp.GetData()) > 0 {
			autoSeedProductPricePlans(ctx, deps, createResp.GetData()[0].GetId(), planID, currency)
		}

		return centymo.HTMXSuccess("plan-price-plans-table")
	})
}

// NewPricePlanEditAction creates the price plan edit action (GET = form, POST = update).
// URL: /action/plan/{id}/price-plans/edit/{ppid}
func NewPricePlanEditAction(deps *PricePlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		planID := viewCtx.Request.PathValue("id")
		ppID := viewCtx.Request.PathValue("ppid")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
				Data: &priceplanpb.PricePlan{Id: ppID},
			})
			if err != nil {
				log.Printf("Failed to read price plan %s: %v", ppID, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			data := resp.GetData()
			if len(data) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			pp := data[0]

			amountStr := strconv.FormatFloat(float64(pp.GetBillingAmount())/100.0, 'f', 2, 64)
			durationStr := strconv.FormatInt(int64(pp.GetDurationValue()), 10)
			selectedScheduleID := pp.GetPriceScheduleId()
			schedules := loadScheduleOptions(ctx, deps, scheduleLocationHintPrefix)

			inUse := false
			lockMsg := ""
			if deps.GetPricePlanInUseIDs != nil {
				if m, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppID}); m[ppID] {
					inUse = true
					lockMsg = deps.PricePlanLabels.Messages.PricingLockedReason
				}
			}

			billingCycleStr := ""
			if v := pp.GetBillingCycleValue(); v > 0 {
				billingCycleStr = fmt.Sprintf("%d", v)
			}
			defaultTermStr := ""
			if v := pp.GetDefaultTermValue(); v > 0 {
				defaultTermStr = fmt.Sprintf("%d", v)
			}
			formLabels := deps.PricePlanLabels.Form
			parentPlan := readPlan(ctx, deps, planID)
			planName := planID
			if parentPlan != nil {
				if n := parentPlan.GetName(); n != "" {
					planName = n
				}
			}
			// Plan §6.7 — apply the same readonly-schedule rule on edit.
			// When the parent Plan is client-scoped, prefer the existing
			// PricePlan's price_schedule_id over the resolver lookup so the
			// drawer never silently rebinds an in-flight record.
			editScheduleMode, editScheduleLockID, editScheduleLockLabel, editScheduleLockClient := resolveScheduleLock(ctx, deps, parentPlan)
			editScheduleID := selectedScheduleID
			if editScheduleMode == "readonly" {
				if selectedScheduleID != "" {
					editScheduleLockID = selectedScheduleID
				}
				editScheduleID = editScheduleLockID
			}
			editLabels := form.LabelsFromPricePlan(formLabels)
			// Context-specific overrides; intentionally inlined (Decision 2 — only 2 fields).
			editLabels.ScheduleLockedTooltip = applyScheduleLockTooltip(editLabels.ScheduleLockedTooltip, editScheduleLockClient)
			editScheduleAutoHint := buildScheduleAutoHint(formLabels, editScheduleMode, editScheduleLockID, editScheduleLockClient)

			return view.OK("price-plan-drawer-form", &form.Data{
				FormAction:            route.ResolveURL(deps.Routes.PricePlanEditURL, "id", planID, "ppid", ppID),
				IsEdit:                true,
				Context:               form.ContextPlan,
				ID:                    ppID,
				PlanID:                planID,
				PlanName:              planName,
				ScheduleID:            editScheduleID,
				Name:                  pp.GetName(),
				Description:           pp.GetDescription(),
				Amount:                amountStr,
				Currency:              pp.GetBillingCurrency(),
				DurationValue:         durationStr,
				DurationUnit:          pp.GetDurationUnit(),
				Active:                pp.GetActive(),
				BillingKind:           pp.GetBillingKind().String(),
				AmountBasis:           pp.GetAmountBasis().String(),
				BillingCycleValue:     billingCycleStr,
				BillingCycleUnit:      pp.GetBillingCycleUnit(),
				TermValue:             defaultTermStr,
				TermUnit:              pp.GetDefaultTermUnit(),
				EntitledOccurrences: func() string {
					if pp.GetEntitledOccurrences() > 0 {
						return strconv.FormatInt(int64(pp.GetEntitledOccurrences()), 10)
					}
					return ""
				}(),
				DurationUnitOptions:   form.BuildDurationUnitOptions(deps.CommonLabels),
				ScheduleOptions:       form.BuildOptions(schedules, selectedScheduleID),
				SelectedScheduleID:    selectedScheduleID,
				SelectedScheduleLabel: form.FindLabel(schedules, selectedScheduleID),
				// Plan §6.7 — when the parent Plan is client-scoped, the
				// schedule field collapses to readonly carrying the
				// existing/derived schedule label. ScheduleFieldMode stays
				// "" (= picker) when the plan is master.
				ScheduleFieldMode:        editScheduleMode,
				ScheduleLabel:            editScheduleLockLabel,
				ScheduleLockedClientName: editScheduleLockClient,
				ScheduleAutoHint:         editScheduleAutoHint,
				InUse:                    inUse,
				LockMessage:              lockMsg,
				// 2026-04-30 cyclic-subscription-jobs plan §7.4 — surface the
				// parent Plan's cyclic flag so the drawer's
				// applyBasisOptionGuards() can disable MILESTONE.
				ParentPlanIsCyclic: parentPlan != nil && parentPlan.GetVisitsPerCycle() > 1,
				Labels:             editLabels,
				CommonLabels:       deps.CommonLabels,
			})
		}

		// POST — update price plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		amount := int64(0)
		if v, err := strconv.ParseFloat(r.FormValue("amount"), 64); err == nil {
			amount = int64(math.Round(v * 100))
		}

		currency := r.FormValue("currency")

		editPPName := r.FormValue("name")
		editPPDescription := r.FormValue("description")
		pp := &priceplanpb.PricePlan{
			Id:              ppID,
			PlanId:          planID,
			Name:            &editPPName,
			Description:     &editPPDescription,
			BillingAmount:   amount,
			BillingCurrency: currency,
			Active:          active,
		}
		// Phase 1 legacy dual-write — proto fields now optional; only set when present.
		if dvStr := r.FormValue("duration_value"); dvStr != "" {
			if v, err := strconv.ParseInt(dvStr, 10, 32); err == nil {
				dv32 := int32(v)
				pp.DurationValue = &dv32
			}
		}
		if du := r.FormValue("duration_unit"); du != "" {
			pp.DurationUnit = &du
		}
		if schedID := r.FormValue("price_schedule_id"); schedID != "" {
			pp.PriceScheduleId = &schedID
		}
		applyBillingFields(pp, r)

		// Plumb the lyngua-resolved suffix through context — see comment
		// in NewPricePlanAddAction. The update path also routes through
		// applyClientScopedScheduleRule when the operator clears
		// price_schedule_id on edit.
		if suffix := deps.PriceScheduleLabels.Form.CustomClientPriceScheduleLabelSuffix; suffix != "" {
			ctx = context.WithValue(ctx, "clientScheduleSuffix", suffix)
		}

		if _, err := deps.UpdatePricePlan(ctx, &priceplanpb.UpdatePricePlanRequest{Data: pp}); err != nil {
			log.Printf("Failed to update price plan %s: %v", ppID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("plan-price-plans-table")
	})
}

// NewPricePlanDeleteAction creates the price plan delete action (POST only).
// URL: /action/plan/{id}/price-plans/delete  (id=price_plan_id via query param)
func NewPricePlanDeleteAction(deps *PricePlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		ppID := viewCtx.Request.URL.Query().Get("id")
		if ppID == "" {
			_ = viewCtx.Request.ParseForm()
			ppID = viewCtx.Request.FormValue("id")
		}
		if ppID == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeletePricePlan(ctx, &priceplanpb.DeletePricePlanRequest{
			Data: &priceplanpb.PricePlan{Id: ppID},
		})
		if err != nil {
			log.Printf("Failed to delete price plan %s: %v", ppID, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plan-price-plans-table")
	})
}

// applyBillingFields reads the Wave 2 billing-semantics fields from the form
// and writes them onto pp. Mirrors the standalone price_plan action so the
// plan-nested drawer persists billing_kind, amount_basis, billing_cycle_*
// and default_term_* alongside the deprecated duration_* dual-write.
func applyBillingFields(pp *priceplanpb.PricePlan, r *http.Request) {
	if v := r.FormValue("billing_kind"); v != "" {
		if bk, ok := priceplanpb.BillingKind_value[v]; ok {
			pp.BillingKind = priceplanpb.BillingKind(bk)
		}
	}
	if v := r.FormValue("amount_basis"); v != "" {
		if ab, ok := priceplanpb.AmountBasis_value[v]; ok {
			pp.AmountBasis = priceplanpb.AmountBasis(ab)
		}
	}
	if s := r.FormValue("billing_cycle_value"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil {
			v32 := int32(n)
			pp.BillingCycleValue = &v32
		}
	}
	if u := r.FormValue("billing_cycle_unit"); u != "" {
		pp.BillingCycleUnit = &u
	}
	if s := r.FormValue("default_term_value"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil {
			v32 := int32(n)
			pp.DefaultTermValue = &v32
		}
	}
	if u := r.FormValue("default_term_unit"); u != "" {
		pp.DefaultTermUnit = &u
	}
	// 2026-05-01 ad-hoc-subscription-billing plan §2.3 — entitled_occurrences
	// only meaningful on AD_HOC × TOTAL_PACKAGE; the drawer JS clears the
	// input on every other (kind × basis) so r.FormValue is empty for those
	// combos and the field stays nil. Server-side validate_ad_hoc.go
	// enforces the same rule (codex MAJ-1 + MAJ-4).
	if s := r.FormValue("entitled_occurrences"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil {
			v32 := int32(n)
			pp.EntitledOccurrences = &v32
		}
	}
}

