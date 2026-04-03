package action

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

// LocationOption is a minimal struct for rendering location options in the price plan form.
type LocationOption struct {
	Id   string
	Name string
}

// PricePlanFormLabels holds i18n labels for the price plan drawer form template.
type PricePlanFormLabels struct {
	Name                string
	NamePlaceholder     string
	Description         string
	DescPlaceholder     string
	Amount              string
	AmountPlaceholder   string
	Currency            string
	CurrencyPlaceholder string
	DurationValue       string
	DurationUnit        string
	Location            string
	LocationPlaceholder string
	SelectLocation      string
	Active              string
}

// PricePlanFormData is the template data for the price plan drawer form.
type PricePlanFormData struct {
	FormAction         string
	IsEdit             bool
	ID                 string
	PlanID             string
	Name               string
	Description        string
	Amount             string
	Currency           string
	DurationValue      string
	DurationUnit       string
	Active             bool
	Locations          []*LocationOption
	SelectedLocationID string
	Labels             PricePlanFormLabels
	CommonLabels       any
}

// PricePlanDeps holds dependencies for price plan action handlers.
type PricePlanDeps struct {
	Routes            centymo.PlanRoutes
	Labels            centymo.PlanLabels
	CreatePricePlan   func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)
	ReadPricePlan     func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	UpdatePricePlan   func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan   func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)
	ListLocations     func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
}

// pricePlanFormLabels converts centymo.PricePlanFormLabels into the local type.
func pricePlanFormLabels(l centymo.PricePlanFormLabels) PricePlanFormLabels {
	return PricePlanFormLabels{
		Name:                l.Name,
		NamePlaceholder:     l.NamePlaceholder,
		Description:         l.Description,
		DescPlaceholder:     l.DescPlaceholder,
		Amount:              l.Amount,
		AmountPlaceholder:   l.AmountPlaceholder,
		Currency:            l.Currency,
		CurrencyPlaceholder: l.CurrencyPlaceholder,
		DurationValue:       l.DurationValue,
		DurationUnit:        l.DurationUnit,
		Location:            l.Location,
		LocationPlaceholder: l.LocationPlaceholder,
		SelectLocation:      l.SelectLocation,
		Active:              l.Active,
	}
}

// loadLocationOptions fetches the location list and converts to options.
// Returns nil slice on error (graceful degradation).
func loadLocationOptions(ctx context.Context, deps *PricePlanDeps) []*LocationOption {
	if deps.ListLocations == nil {
		return nil
	}
	resp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
	if err != nil {
		log.Printf("Failed to load locations for price plan form: %v", err)
		return nil
	}
	var options []*LocationOption
	for _, loc := range resp.GetData() {
		options = append(options, &LocationOption{
			Id:   loc.GetId(),
			Name: loc.GetName(),
		})
	}
	return options
}

// NewPricePlanAddAction creates the price plan add action (GET = form, POST = create).
// URL: /action/plans/{id}/pricelists/add
func NewPricePlanAddAction(deps *PricePlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		planID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("price-plan-drawer-form", &PricePlanFormData{
				FormAction:   route.ResolveURL(deps.Routes.PricePlanAddURL, "id", planID),
				PlanID:       planID,
				Active:       true,
				Currency:     "PHP",
				DurationUnit: "months",
				Locations:    loadLocationOptions(ctx, deps),
				Labels:       pricePlanFormLabels(deps.Labels.PricePlanForm),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create price plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		amount := float64(0)
		if v, err := strconv.ParseFloat(r.FormValue("amount"), 64); err == nil {
			amount = v
		}

		durationValue := int32(0)
		if v, err := strconv.ParseInt(r.FormValue("duration_value"), 10, 32); err == nil {
			durationValue = int32(v)
		}

		pp := &priceplanpb.PricePlan{
			PlanId:        planID,
			Name:          r.FormValue("name"),
			Description:   r.FormValue("description"),
			Amount:        amount,
			Currency:      r.FormValue("currency"),
			DurationValue: durationValue,
			DurationUnit:  r.FormValue("duration_unit"),
			Active:        active,
		}
		if locID := r.FormValue("location_id"); locID != "" {
			pp.LocationId = &locID
		}

		_, err := deps.CreatePricePlan(ctx, &priceplanpb.CreatePricePlanRequest{
			Data: pp,
		})
		if err != nil {
			log.Printf("Failed to create price plan for plan %s: %v", planID, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plan-pricelists-table")
	})
}

// NewPricePlanEditAction creates the price plan edit action (GET = form, POST = update).
// URL: /action/plans/{id}/pricelists/edit/{ppid}
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

			amountStr := strconv.FormatFloat(pp.GetAmount(), 'f', 2, 64)
			durationStr := strconv.FormatInt(int64(pp.GetDurationValue()), 10)

			return view.OK("price-plan-drawer-form", &PricePlanFormData{
				FormAction:         route.ResolveURL(deps.Routes.PricePlanEditURL, "id", planID, "ppid", ppID),
				IsEdit:             true,
				ID:                 ppID,
				PlanID:             planID,
				Name:               pp.GetName(),
				Description:        pp.GetDescription(),
				Amount:             amountStr,
				Currency:           pp.GetCurrency(),
				DurationValue:      durationStr,
				DurationUnit:       pp.GetDurationUnit(),
				Active:             pp.GetActive(),
				Locations:          loadLocationOptions(ctx, deps),
				SelectedLocationID: pp.GetLocationId(),
				Labels:             pricePlanFormLabels(deps.Labels.PricePlanForm),
				CommonLabels:       nil, // injected by ViewAdapter
			})
		}

		// POST — update price plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		amount := float64(0)
		if v, err := strconv.ParseFloat(r.FormValue("amount"), 64); err == nil {
			amount = v
		}

		durationValue := int32(0)
		if v, err := strconv.ParseInt(r.FormValue("duration_value"), 10, 32); err == nil {
			durationValue = int32(v)
		}

		pp := &priceplanpb.PricePlan{
			Id:            ppID,
			PlanId:        planID,
			Name:          r.FormValue("name"),
			Description:   r.FormValue("description"),
			Amount:        amount,
			Currency:      r.FormValue("currency"),
			DurationValue: durationValue,
			DurationUnit:  r.FormValue("duration_unit"),
			Active:        active,
		}
		if locID := r.FormValue("location_id"); locID != "" {
			pp.LocationId = &locID
		}

		_, err := deps.UpdatePricePlan(ctx, &priceplanpb.UpdatePricePlanRequest{
			Data: pp,
		})
		if err != nil {
			log.Printf("Failed to update price plan %s: %v", ppID, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plan-pricelists-table")
	})
}

// NewPricePlanDeleteAction creates the price plan delete action (POST only).
// URL: /action/plans/{id}/pricelists/delete  (id=price_plan_id via query param)
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

		return centymo.HTMXSuccess("plan-pricelists-table")
	})
}
