package plan

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/price_plan/form"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

// EditFormData is the drawer form for editing a price_plan under a schedule.
type EditFormData struct {
	FormAction    string
	ScheduleID    string
	ScheduleName  string
	ID            string
	PlanID        string
	PlanLabel     string // display label for the currently-selected plan (for SelectedLabel on auto-complete)
	PlanOptions   []map[string]any
	Name          string
	Description   string
	Amount        string
	Currency      string
	DurationValue string // Phase 1 legacy dual-write
	DurationUnit  string // Phase 1 legacy dual-write
	CommonLabels  pyeza.CommonLabels
	Labels        centymo.PriceScheduleLabels

	// Wave 2: new billing semantics fields.
	//
	// 2026-04-30 enum-select-canonicalize — BillingKindOptions /
	// AmountBasisOptions removed; the drawer template hardcodes the option
	// list. Only the selected value is passed in.
	BillingKind       string
	AmountBasis       string
	BillingCycleValue string
	BillingCycleUnit  string
	// kept as DefaultTermValue/Unit on the wire (form input names) but renamed
	// in the form.Data struct.
	DurationUnitOptions []types.SelectOption

	// PricingLocked is true when the price_plan is referenced by active subscriptions.
	// The Pricing section fields (Amount, Currency, Duration, DurationUnit) are rendered
	// as read-only in the drawer, but all other fields remain editable.
	PricingLocked       bool
	PricingLockedReason string
}

// NewEditAction handles GET/POST /action/price-schedule/{id}/plan/{ppid}/edit.
func NewEditAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "update") {
			return centymo.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: ppid}})
			if err != nil || len(resp.GetData()) == 0 {
				return centymo.HTMXError(deps.PlanLabels.Errors.NotFound)
			}
			pp := resp.GetData()[0]

			// Check whether this plan is referenced by active subscriptions.
			// When true, the Pricing section is rendered read-only in the drawer.
			pricingLocked := false
			pricingLockedReason := ""
			if deps.GetPricePlanInUseIDs != nil {
				inUseMap, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppid})
				if inUseMap[ppid] {
					pricingLocked = true
					pricingLockedReason = deps.PlanLabels.Messages.PricingLockedReason
				}
			}

			// Populate new billing fields from existing record.
			billingCycleValue := ""
			if v := pp.GetBillingCycleValue(); v > 0 {
				billingCycleValue = fmt.Sprintf("%d", v)
			}
			defaultTermValue := ""
			if v := pp.GetDefaultTermValue(); v > 0 {
				defaultTermValue = fmt.Sprintf("%d", v)
			}
			formLabels := deps.PlanLabels.Form
			scheduleName, scheduleClientID := lookupScheduleNameAndClient(ctx, deps, sid)
			// 2026-05-03 — same mutually-exclusive plan-scope filter as the
			// schedule-add drawer: client-scoped schedule shows only matching
			// client's plans, master schedule shows only master plans.
			planOpts := buildPlanOptions(ctx, deps, pp.GetPlanId(), scheduleClientID)
			return view.OK("price-plan-drawer-form", &form.Data{
				FormAction:             route.ResolveURL(deps.Routes.PlanEditURL, "id", sid, "ppid", ppid),
				IsEdit:                 true,
				Context:                form.ContextSchedule,
				ID:                     ppid,
				PlanID:                 pp.GetPlanId(),
				ScheduleID:             sid,
				ScheduleName:           scheduleName,
				ParentScheduleClientID: scheduleClientID,
				Name:                   pp.GetName(),
				Description:            pp.GetDescription(),
				Amount:                 strconv.FormatFloat(float64(pp.GetBillingAmount())/100.0, 'f', 2, 64),
				Currency:               pp.GetBillingCurrency(),
				DurationValue:          fmt.Sprintf("%d", pp.GetDurationValue()),
				DurationUnit:           pp.GetDurationUnit(),
				Active:                 pp.GetActive(),
				// Wave 2: populate new billing fields.
				BillingKind:         pp.GetBillingKind().String(),
				AmountBasis:         pp.GetAmountBasis().String(),
				BillingCycleValue:   billingCycleValue,
				BillingCycleUnit:    pp.GetBillingCycleUnit(),
				TermValue:           defaultTermValue,
				TermUnit:            pp.GetDefaultTermUnit(),
				DurationUnitOptions: buildDurationUnitOptions(deps.CommonLabels),
				PlanOptions:         planOpts,
				SelectedPlanID:      pp.GetPlanId(),
				SelectedPlanLabel:   labelFromOptions(planOpts, pp.GetPlanId()),
				InUse:               pricingLocked,
				LockMessage:         pricingLockedReason,
				Labels:              form.LabelsFromPricePlan(formLabels),
				CommonLabels:        deps.CommonLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.PlanLabels.Errors.UpdateFailed)
		}
		r := viewCtx.Request
		amount := int64(0)
		if v, err := strconv.ParseFloat(r.FormValue("amount"), 64); err == nil {
			amount = int64(math.Round(v * 100))
		}
		dvStr := r.FormValue("duration_value")
		currency := r.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}
		// Read existing to preserve active state (not in form) and to enforce
		// pricing-field immutability when the plan is in use by active subscriptions.
		existing, _ := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: ppid}})
		active := true
		if existing != nil && len(existing.GetData()) > 0 {
			active = existing.GetData()[0].GetActive()
		}

		// Server-side guard: if this plan is referenced by active subscriptions,
		// overwrite the four pricing fields with the existing DB values so a client
		// cannot bypass the read-only drawer by editing the HTML.
		durationUnit := r.FormValue("duration_unit")
		var existingDurationValue *int32
		var existingDurationUnit *string
		if deps.GetPricePlanInUseIDs != nil && existing != nil && len(existing.GetData()) > 0 {
			inUseMap, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppid})
			if inUseMap[ppid] {
				ex := existing.GetData()[0]
				amount = ex.GetBillingAmount()
				currency = ex.GetBillingCurrency()
				existingDurationValue = ex.DurationValue
				existingDurationUnit = ex.DurationUnit
				dvStr = ""        // sentinel: prefer existingDurationValue below
				durationUnit = "" // sentinel: prefer existingDurationUnit below
			}
		}

		// Wave 2: new billing semantics fields.
		bcvStr := r.FormValue("billing_cycle_value")
		bcv, _ := strconv.ParseInt(bcvStr, 10, 32)
		bcu := r.FormValue("billing_cycle_unit")
		dtvStr := r.FormValue("default_term_value")
		dtv, _ := strconv.ParseInt(dtvStr, 10, 32)
		dtu := r.FormValue("default_term_unit")
		billingKindStr := r.FormValue("billing_kind")
		amountBasisStr := r.FormValue("amount_basis")

		planPageName := r.FormValue("name")
		planPageDesc := r.FormValue("description")
		req := &priceplanpb.UpdatePricePlanRequest{
			Data: &priceplanpb.PricePlan{
				Id:              ppid,
				PlanId:          r.FormValue("plan_id"),
				Name:            &planPageName,
				Description:     &planPageDesc,
				BillingAmount:   amount,
				BillingCurrency: currency,
				Active:          active,
			},
		}
		req.Data.PriceScheduleId = &sid
		// Phase 1 legacy dual-write — proto fields now optional. Prefer the
		// in-use snapshot when locked; otherwise read from form input.
		if existingDurationValue != nil {
			req.Data.DurationValue = existingDurationValue
		} else if dvStr != "" {
			if parsed, err := strconv.ParseInt(dvStr, 10, 32); err == nil {
				dv32 := int32(parsed)
				req.Data.DurationValue = &dv32
			}
		}
		if existingDurationUnit != nil {
			req.Data.DurationUnit = existingDurationUnit
		} else if durationUnit != "" {
			req.Data.DurationUnit = &durationUnit
		}
		// Set new enum fields.
		if billingKindStr != "" {
			if bk, ok := priceplanpb.BillingKind_value[billingKindStr]; ok {
				req.Data.BillingKind = priceplanpb.BillingKind(bk)
			}
		}
		if amountBasisStr != "" {
			if ab, ok := priceplanpb.AmountBasis_value[amountBasisStr]; ok {
				req.Data.AmountBasis = priceplanpb.AmountBasis(ab)
			}
		}
		// Set new optional duration fields.
		if bcvStr != "" {
			bcv32 := int32(bcv)
			req.Data.BillingCycleValue = &bcv32
		}
		if bcu != "" {
			req.Data.BillingCycleUnit = &bcu
		}
		if dtvStr != "" {
			dtv32 := int32(dtv)
			req.Data.DefaultTermValue = &dtv32
		}
		if dtu != "" {
			req.Data.DefaultTermUnit = &dtu
		}
		if _, err := deps.UpdatePricePlan(ctx, req); err != nil {
			log.Printf("Failed to update price plan %s under schedule %s: %v", ppid, sid, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedule-plans-table")
	})
}

// buildPlanOptions returns the plan picker options with a strict, mutually-
// exclusive client-scope filter mirroring the schedule-add drawer:
//
//   - Client-scoped schedule (scheduleClientID != ""): only plans whose
//     client_id matches scheduleClientID. Master plans are excluded.
//   - Master schedule (scheduleClientID == ""): only master plans
//     (client_id empty). Client-scoped plans cannot attach to a master schedule.
//
// `selectedID` is preserved for edit-mode drawer pre-selection regardless of
// the filter (the row is already wired; we never silently drop the operator's
// existing selection).
func buildPlanOptions(ctx context.Context, deps *DetailViewDeps, selectedID, scheduleClientID string) []map[string]any {
	if deps.ListPlans == nil {
		return nil
	}
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return nil
	}
	opts := make([]map[string]any, 0, len(resp.GetData()))
	for _, p := range resp.GetData() {
		if p == nil || !p.GetActive() {
			continue
		}
		planClientID := p.GetClientId()
		if scheduleClientID != "" {
			if planClientID != scheduleClientID && p.GetId() != selectedID {
				continue
			}
		} else {
			if planClientID != "" && p.GetId() != selectedID {
				continue
			}
		}
		opts = append(opts, map[string]any{
			"Value":       p.GetId(),
			"Label":       p.GetName(),
			"Description": p.GetDescription(),
			"Selected":    p.GetId() == selectedID,
		})
	}
	return opts
}

func loadProductOptions(ctx context.Context, deps *DetailViewDeps, planID, selectedProductID string) []types.SelectOption {
	productNames := map[string]string{}
	if deps.ListProducts != nil {
		prodResp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
		if err == nil {
			for _, p := range prodResp.GetData() {
				if p != nil {
					productNames[p.GetId()] = p.GetName()
				}
			}
		}
	}

	if deps.ListProductPlans == nil || planID == "" {
		// Fallback: all products
		options := make([]types.SelectOption, 0, len(productNames))
		for pid, name := range productNames {
			options = append(options, types.SelectOption{
				Value:    pid,
				Label:    name,
				Selected: pid == selectedProductID,
			})
		}
		return options
	}

	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
	if err != nil {
		return nil
	}
	options := []types.SelectOption{}
	for _, pp := range ppResp.GetData() {
		if pp == nil || pp.GetPlanId() != planID {
			continue
		}
		pid := pp.GetProductId()
		name := productNames[pid]
		if name == "" {
			name = pid
		}
		options = append(options, types.SelectOption{
			Value:    pid,
			Label:    name,
			Selected: pid == selectedProductID,
		})
	}
	return options
}

// labelFromOptions returns the Label string for the option whose Value matches id.
// Used to populate SelectedLabel on the edit-drawer auto-complete.
func labelFromOptions(opts []map[string]any, id string) string {
	for _, opt := range opts {
		if v, ok := opt["Value"].(string); ok && v == id {
			if label, ok := opt["Label"].(string); ok {
				return label
			}
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// Option builder helpers — non-proto-enum only
// ---------------------------------------------------------------------------
//
// 2026-04-30 enum-select-canonicalize plan §6 — the proto-enum option
// builders (BillingKind, AmountBasis, BillingTreatment) are gone. Their
// option lists now live as hardcoded <option> tags in the drawer
// templates, and a checked-in drift test (price_plan/templates/templates_test.go)
// keeps them aligned with the proto enum's _name map.

// buildDurationUnitOptions builds select options for billing_cycle_unit / default_term_unit
// reusing the existing DurationUnit labels from CommonLabels. duration_unit is
// a plain string column (not a proto enum), so its option builder stays.
func buildDurationUnitOptions(cl pyeza.CommonLabels) []types.SelectOption {
	du := cl.DurationUnit
	return []types.SelectOption{
		{Value: "day", Label: du.DaySelect},
		{Value: "week", Label: du.WeekSelect},
		{Value: "month", Label: du.MonthSelect},
		{Value: "year", Label: du.YearSelect},
	}
}
