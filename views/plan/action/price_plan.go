package action

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/price_plan/form"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
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
	Routes             centymo.PlanRoutes
	Labels             centymo.PlanLabels
	CommonLabels       pyeza.CommonLabels
	CreatePricePlan    func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)
	ReadPricePlan      func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	UpdatePricePlan    func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan    func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)
	ListPriceSchedules func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)

	// ReadPlan resolves the parent plan's name for display in the locked
	// "Package" field on the drawer.
	ReadPlan func(ctx context.Context, req *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)

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

// lookupPlanName reads the parent plan and returns its name for the locked
// Package display. Falls back to the planID on any error.
func lookupPlanName(ctx context.Context, deps *PricePlanDeps, planID string) string {
	if deps.ReadPlan == nil {
		return planID
	}
	resp, err := deps.ReadPlan(ctx, &planpb.ReadPlanRequest{Data: &planpb.Plan{Id: &planID}})
	if err != nil || len(resp.GetData()) == 0 {
		return planID
	}
	if name := resp.GetData()[0].GetName(); name != "" {
		return name
	}
	return planID
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
		if productID == "" {
			continue
		}
		var price int64
		rowCurrency := currency
		if rowCurrency == "" {
			rowCurrency = "PHP"
		}
		if prod := products[productID]; prod != nil {
			price = prod.GetPrice()
			if c := prod.GetCurrency(); c != "" {
				rowCurrency = c
			}
		}
		if _, err := deps.CreateProductPricePlan(ctx, &productpriceplanpb.CreateProductPricePlanRequest{
			Data: &productpriceplanpb.ProductPricePlan{
				PricePlanId: pricePlanID,
				ProductId:   productID,
				Price:       price,
				Currency:    rowCurrency,
				Active:      true,
			},
		}); err != nil {
			log.Printf("auto-seed product_price_plan failed for %s/%s: %v", pricePlanID, productID, err)
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
			return view.OK("price-plan-drawer-form", &form.Data{
				FormAction:      route.ResolveURL(deps.Routes.PricePlanAddURL, "id", planID),
				Context:         form.ContextPlan,
				PlanID:          planID,
				PlanName:        lookupPlanName(ctx, deps, planID),
				Active:          true,
				Currency:        "PHP",
				DurationUnit:    "months",
				ScheduleOptions: form.BuildOptions(schedules, ""),
				Labels:          form.LabelsFromPricePlan(deps.Labels.PricePlanForm),
				CommonLabels:    deps.CommonLabels,
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

		durationValue := int32(0)
		if v, err := strconv.ParseInt(r.FormValue("duration_value"), 10, 32); err == nil {
			durationValue = int32(v)
		}

		currency := r.FormValue("currency")

		ppName := r.FormValue("name")
		ppDescription := r.FormValue("description")
		pp := &priceplanpb.PricePlan{
			PlanId:        planID,
			Name:          &ppName,
			Description:   &ppDescription,
			Amount:        amount,
			Currency:      currency,
			DurationValue: durationValue,
			DurationUnit:  r.FormValue("duration_unit"),
			Active:        active,
		}
		if schedID := r.FormValue("price_schedule_id"); schedID != "" {
			pp.PriceScheduleId = &schedID
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

			amountStr := strconv.FormatFloat(float64(pp.GetAmount())/100.0, 'f', 2, 64)
			durationStr := strconv.FormatInt(int64(pp.GetDurationValue()), 10)
			selectedScheduleID := pp.GetPriceScheduleId()
			schedules := loadScheduleOptions(ctx, deps, scheduleLocationHintPrefix)

			inUse := false
			lockMsg := ""
			if deps.GetPricePlanInUseIDs != nil {
				if m, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppID}); m[ppID] {
					inUse = true
					lockMsg = "This price plan is in use by active subscriptions. Pricing changes are disabled."
				}
			}

			return view.OK("price-plan-drawer-form", &form.Data{
				FormAction:            route.ResolveURL(deps.Routes.PricePlanEditURL, "id", planID, "ppid", ppID),
				IsEdit:                true,
				Context:               form.ContextPlan,
				ID:                    ppID,
				PlanID:                planID,
				PlanName:              lookupPlanName(ctx, deps, planID),
				ScheduleID:            selectedScheduleID,
				Name:                  pp.GetName(),
				Description:           pp.GetDescription(),
				Amount:                amountStr,
				Currency:              pp.GetCurrency(),
				DurationValue:         durationStr,
				DurationUnit:          pp.GetDurationUnit(),
				Active:                pp.GetActive(),
				ScheduleOptions:       form.BuildOptions(schedules, selectedScheduleID),
				SelectedScheduleID:    selectedScheduleID,
				SelectedScheduleLabel: form.FindLabel(schedules, selectedScheduleID),
				InUse:                 inUse,
				LockMessage:           lockMsg,
				Labels:                form.LabelsFromPricePlan(deps.Labels.PricePlanForm),
				CommonLabels:          deps.CommonLabels,
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

		durationValue := int32(0)
		if v, err := strconv.ParseInt(r.FormValue("duration_value"), 10, 32); err == nil {
			durationValue = int32(v)
		}

		currency := r.FormValue("currency")

		editPPName := r.FormValue("name")
		editPPDescription := r.FormValue("description")
		pp := &priceplanpb.PricePlan{
			Id:            ppID,
			PlanId:        planID,
			Name:          &editPPName,
			Description:   &editPPDescription,
			Amount:        amount,
			Currency:      currency,
			DurationValue: durationValue,
			DurationUnit:  r.FormValue("duration_unit"),
			Active:        active,
		}
		if schedID := r.FormValue("price_schedule_id"); schedID != "" {
			pp.PriceScheduleId = &schedID
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
