package action

import (
	"context"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
)

type LocationOption struct {
	Id   string
	Name string
}

type FormData struct {
	FormAction            string
	IsEdit                bool
	ID                    string
	Name                  string
	Description           string
	DateStart             string
	DateEnd               string
	Active                bool
	Locations             []*LocationOption
	SelectedLocationID    string
	SelectedLocationLabel string
	LocationOptions       []map[string]any
	Labels                centymo.PriceScheduleFormLabels
	CommonLabels          any
}

type Deps struct {
	Routes              centymo.PriceScheduleRoutes
	Labels              centymo.PriceScheduleLabels
	CreatePriceSchedule func(ctx context.Context, req *priceschedulepb.CreatePriceScheduleRequest) (*priceschedulepb.CreatePriceScheduleResponse, error)
	ReadPriceSchedule   func(ctx context.Context, req *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error)
	UpdatePriceSchedule func(ctx context.Context, req *priceschedulepb.UpdatePriceScheduleRequest) (*priceschedulepb.UpdatePriceScheduleResponse, error)
	DeletePriceSchedule func(ctx context.Context, req *priceschedulepb.DeletePriceScheduleRequest) (*priceschedulepb.DeletePriceScheduleResponse, error)
	ListLocations       func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
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

func buildLocationAutoCompleteOptions(locations []*LocationOption, selectedID string) []map[string]any {
	opts := make([]map[string]any, 0, len(locations))
	for _, loc := range locations {
		opts = append(opts, map[string]any{
			"Value":    loc.Id,
			"Label":    loc.Name,
			"Selected": loc.Id == selectedID,
		})
	}
	return opts
}

func findLocationLabel(locations []*LocationOption, id string) string {
	for _, loc := range locations {
		if loc.Id == id {
			return loc.Name
		}
	}
	return ""
}

func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule", "create") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			locations := loadLocations(ctx, deps)
			return view.OK("price-schedule-drawer-form", &FormData{
				FormAction:      deps.Routes.AddURL,
				Active:          true,
				Locations:       locations,
				LocationOptions: buildLocationAutoCompleteOptions(locations, ""),
				Labels:          deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		locationID := r.FormValue("location_id")
		dateEnd := r.FormValue("date_end")
		description := r.FormValue("description")
		req := &priceschedulepb.CreatePriceScheduleRequest{
			Data: &priceschedulepb.PriceSchedule{
				Name:      r.FormValue("name"),
				DateStart: r.FormValue("date_start"),
				Active:    active,
			},
		}
		if description != "" {
			req.Data.Description = &description
		}
		if dateEnd != "" {
			req.Data.DateEnd = &dateEnd
		}
		if locationID != "" {
			req.Data.LocationId = &locationID
		}
		if _, err := deps.CreatePriceSchedule(ctx, req); err != nil {
			log.Printf("Failed to create price schedule: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedules-table")
	})
}

func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule", "update") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{Data: &priceschedulepb.PriceSchedule{Id: id}})
			if err != nil || len(resp.GetData()) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := resp.GetData()[0]
			locations := loadLocations(ctx, deps)
			selectedLocationID := record.GetLocationId()
			return view.OK("price-schedule-drawer-form", &FormData{
				FormAction:            route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                true,
				ID:                    id,
				Name:                  record.GetName(),
				Description:           record.GetDescription(),
				DateStart:             record.GetDateStart(),
				DateEnd:               record.GetDateEnd(),
				Active:                record.GetActive(),
				SelectedLocationID:    selectedLocationID,
				SelectedLocationLabel: findLocationLabel(locations, selectedLocationID),
				Locations:             locations,
				LocationOptions:       buildLocationAutoCompleteOptions(locations, selectedLocationID),
				Labels:                deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		locationID := r.FormValue("location_id")
		dateEnd := r.FormValue("date_end")
		description := r.FormValue("description")
		req := &priceschedulepb.UpdatePriceScheduleRequest{
			Data: &priceschedulepb.PriceSchedule{
				Id:        id,
				Name:      r.FormValue("name"),
				DateStart: r.FormValue("date_start"),
				Active:    active,
			},
		}
		if description != "" {
			req.Data.Description = &description
		}
		if dateEnd != "" {
			req.Data.DateEnd = &dateEnd
		}
		if locationID != "" {
			req.Data.LocationId = &locationID
		}
		if _, err := deps.UpdatePriceSchedule(ctx, req); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedules-table")
	})
}

func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule", "delete") {
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
		if _, err := deps.DeletePriceSchedule(ctx, &priceschedulepb.DeletePriceScheduleRequest{Data: &priceschedulepb.PriceSchedule{Id: id}}); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedules-table")
	})
}

func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.DeleteFailed)
		}
		for _, id := range viewCtx.Request.Form["id"] {
			if id != "" {
				_, _ = deps.DeletePriceSchedule(ctx, &priceschedulepb.DeletePriceScheduleRequest{Data: &priceschedulepb.PriceSchedule{Id: id}})
			}
		}
		return centymo.HTMXSuccess("price-schedules-table")
	})
}

func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule", "update") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{Data: &priceschedulepb.PriceSchedule{Id: id}})
		if err != nil || len(readResp.GetData()) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdatePriceSchedule(ctx, &priceschedulepb.UpdatePriceScheduleRequest{
			Data: &priceschedulepb.PriceSchedule{
				Id:          id,
				Name:        record.GetName(),
				Description: record.Description,
				DateStart:   record.GetDateStart(),
				DateEnd:     record.DateEnd,
				Active:      status == "active",
				LocationId:  record.LocationId,
			},
		})
		if err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedules-table")
	})
}

func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule", "update") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{Data: &priceschedulepb.PriceSchedule{Id: id}})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdatePriceSchedule(ctx, &priceschedulepb.UpdatePriceScheduleRequest{
				Data: &priceschedulepb.PriceSchedule{
					Id:          id,
					Name:        record.GetName(),
					Description: record.Description,
					DateStart:   record.GetDateStart(),
					DateEnd:     record.DateEnd,
					Active:      status == "active",
					LocationId:  record.LocationId,
				},
			})
		}
		return centymo.HTMXSuccess("price-schedules-table")
	})
}
