package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

type LocationOption struct {
	Id   string
	Name string
}

type PlanOption struct {
	Id   string
	Name string
}

type FormData struct {
	FormAction         string
	IsEdit             bool
	ID                 string
	Name               string
	Description        string
	Amount             string
	Currency           string
	DurationValue      string
	DurationUnit       string
	Active             bool
	PlanID             string
	Locations          []*LocationOption
	SelectedLocationID string
	Plans              []*PlanOption
	SelectedPlanID     string
	Labels             centymo.PricePlanFormLabels
	CommonLabels       any
}

type Deps struct {
	Routes          centymo.PricePlanRoutes
	Labels          centymo.PricePlanLabels
	CreatePricePlan func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)
	ReadPricePlan   func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	UpdatePricePlan func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)
	ListLocations   func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	ListPlans       func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
}

func loadLocations(ctx context.Context, deps *Deps) []*LocationOption {
	if deps.ListLocations == nil {
		return nil
	}
	resp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
	if err != nil {
		return nil
	}
	opts := make([]*LocationOption, 0, len(resp.GetData()))
	for _, loc := range resp.GetData() {
		opts = append(opts, &LocationOption{Id: loc.GetId(), Name: loc.GetName()})
	}
	return opts
}

func loadPlans(ctx context.Context, deps *Deps) []*PlanOption {
	if deps.ListPlans == nil {
		return nil
	}
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return nil
	}
	opts := make([]*PlanOption, 0, len(resp.GetData()))
	for _, p := range resp.GetData() {
		opts = append(opts, &PlanOption{Id: p.GetId(), Name: p.GetName()})
	}
	return opts
}

func parseAmount(s string) int64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f * 100))
}

func formatAmount(centavos int64) string {
	return strconv.FormatFloat(float64(centavos)/100.0, 'f', 2, 64)
}

func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("price-plan-drawer-form", &FormData{
				FormAction: deps.Routes.AddURL,
				Active:     true,
				Currency:   "PHP",
				Locations:  loadLocations(ctx, deps),
				Plans:      loadPlans(ctx, deps),
				Labels:     deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		dv, _ := strconv.ParseInt(r.FormValue("duration_value"), 10, 32)
		locationID := r.FormValue("location_id")
		req := &priceplanpb.CreatePricePlanRequest{
			Data: &priceplanpb.PricePlan{
				PlanId:        r.FormValue("plan_id"),
				Name:          r.FormValue("name"),
				Description:   r.FormValue("description"),
				Amount:        parseAmount(r.FormValue("amount")),
				Currency:      r.FormValue("currency"),
				DurationValue: int32(dv),
				DurationUnit:  r.FormValue("duration_unit"),
				Active:        active,
			},
		}
		if locationID != "" {
			req.Data.LocationId = &locationID
		}
		if _, err := deps.CreatePricePlan(ctx, req); err != nil {
			log.Printf("Failed to create price plan: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}

func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}})
			if err != nil || len(resp.GetData()) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := resp.GetData()[0]
			return view.OK("price-plan-drawer-form", &FormData{
				FormAction:         route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:             true,
				ID:                 id,
				Name:               record.GetName(),
				Description:        record.GetDescription(),
				Amount:             formatAmount(record.GetAmount()),
				Currency:           record.GetCurrency(),
				DurationValue:      fmt.Sprintf("%d", record.GetDurationValue()),
				DurationUnit:       record.GetDurationUnit(),
				Active:             record.GetActive(),
				PlanID:             record.GetPlanId(),
				SelectedPlanID:     record.GetPlanId(),
				SelectedLocationID: record.GetLocationId(),
				Locations:          loadLocations(ctx, deps),
				Plans:              loadPlans(ctx, deps),
				Labels:             deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		dv, _ := strconv.ParseInt(r.FormValue("duration_value"), 10, 32)
		locationID := r.FormValue("location_id")
		req := &priceplanpb.UpdatePricePlanRequest{
			Data: &priceplanpb.PricePlan{
				Id:            id,
				PlanId:        r.FormValue("plan_id"),
				Name:          r.FormValue("name"),
				Description:   r.FormValue("description"),
				Amount:        parseAmount(r.FormValue("amount")),
				Currency:      r.FormValue("currency"),
				DurationValue: int32(dv),
				DurationUnit:  r.FormValue("duration_unit"),
				Active:        active,
			},
		}
		if locationID != "" {
			req.Data.LocationId = &locationID
		}
		if _, err := deps.UpdatePricePlan(ctx, req); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}

func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		if _, err := deps.DeletePricePlan(ctx, &priceplanpb.DeletePricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}}); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}

func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.DeleteFailed)
		}
		for _, id := range viewCtx.Request.Form["id"] {
			if id != "" {
				_, _ = deps.DeletePricePlan(ctx, &priceplanpb.DeletePricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}})
			}
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}

func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}})
		if err != nil || len(readResp.GetData()) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdatePricePlan(ctx, &priceplanpb.UpdatePricePlanRequest{
			Data: &priceplanpb.PricePlan{
				Id: id, PlanId: record.GetPlanId(), Name: record.GetName(),
				Description: record.GetDescription(), Amount: record.GetAmount(),
				Currency: record.GetCurrency(), DurationValue: record.GetDurationValue(),
				DurationUnit: record.GetDurationUnit(), Active: status == "active",
				LocationId: record.LocationId,
			},
		})
		if err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}

func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdatePricePlan(ctx, &priceplanpb.UpdatePricePlanRequest{
				Data: &priceplanpb.PricePlan{
					Id: id, PlanId: record.GetPlanId(), Name: record.GetName(),
					Description: record.GetDescription(), Amount: record.GetAmount(),
					Currency: record.GetCurrency(), DurationValue: record.GetDurationValue(),
					DurationUnit: record.GetDurationUnit(), Active: status == "active",
					LocationId: record.LocationId,
				},
			})
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}
