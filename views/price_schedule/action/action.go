package action

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/route"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// parseScheduleDate parses a YYYY-MM-DD form input as start-of-day in tz and
// returns the UTC timestamp. Empty input → nil.
func parseScheduleDate(input string, tz *time.Location) *timestamppb.Timestamp {
	if input == "" {
		return nil
	}
	t, err := time.ParseInLocation("2006-01-02", input, tz)
	if err != nil {
		return nil
	}
	return timestamppb.New(t.UTC())
}

// formatScheduleDate formats a Timestamp as YYYY-MM-DD in tz for the form input.
func formatScheduleDate(ts *timestamppb.Timestamp, tz *time.Location) string {
	if ts == nil || !ts.IsValid() {
		return ""
	}
	return ts.AsTime().In(tz).Format(pyezatypes.DateInputLayout)
}

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
	Routes                   centymo.PriceScheduleRoutes
	Labels                   centymo.PriceScheduleLabels
	CreatePriceSchedule      func(ctx context.Context, req *priceschedulepb.CreatePriceScheduleRequest) (*priceschedulepb.CreatePriceScheduleResponse, error)
	ReadPriceSchedule        func(ctx context.Context, req *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error)
	UpdatePriceSchedule      func(ctx context.Context, req *priceschedulepb.UpdatePriceScheduleRequest) (*priceschedulepb.UpdatePriceScheduleResponse, error)
	DeletePriceSchedule      func(ctx context.Context, req *priceschedulepb.DeletePriceScheduleRequest) (*priceschedulepb.DeletePriceScheduleResponse, error)
	ListLocations            func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	GetPriceScheduleInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
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
		description := r.FormValue("description")
		tz := pyezatypes.LocationFromContext(ctx)
		req := &priceschedulepb.CreatePriceScheduleRequest{
			Data: &priceschedulepb.PriceSchedule{
				Name:          r.FormValue("name"),
				DateTimeStart: parseScheduleDate(r.FormValue("date_start"), tz),
				DateTimeEnd:   parseScheduleDate(r.FormValue("date_end"), tz),
				Active:        active,
			},
		}
		if description != "" {
			req.Data.Description = &description
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

// NewEditAction creates the price-schedule edit action (GET = form, POST = update).
// When the GET request includes ?clone=1, the handler returns the drawer form
// pre-populated from the source record but wired to AddURL (submission creates
// a new price schedule) with " (Copy)" appended to the name.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("price_schedule", requiredAction) {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{Data: &priceschedulepb.PriceSchedule{Id: id}})
			if err != nil || len(resp.GetData()) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := resp.GetData()[0]
			locations := loadLocations(ctx, deps)
			selectedLocationID := record.GetLocationId()

			name := record.GetName()
			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id)
			formID := id
			if isClone {
				name = strings.TrimSpace(name) + viewCtx.T("actions.copySuffix")
				formAction = deps.Routes.AddURL
				formID = ""
			}
			tz := pyezatypes.LocationFromContext(ctx)
			return view.OK("price-schedule-drawer-form", &FormData{
				FormAction:            formAction,
				IsEdit:                !isClone,
				ID:                    formID,
				Name:                  name,
				Description:           record.GetDescription(),
				DateStart:             formatScheduleDate(record.GetDateTimeStart(), tz),
				DateEnd:               formatScheduleDate(record.GetDateTimeEnd(), tz),
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
		description := r.FormValue("description")
		tz := pyezatypes.LocationFromContext(ctx)
		req := &priceschedulepb.UpdatePriceScheduleRequest{
			Data: &priceschedulepb.PriceSchedule{
				Id:            id,
				Name:          r.FormValue("name"),
				DateTimeStart: parseScheduleDate(r.FormValue("date_start"), tz),
				DateTimeEnd:   parseScheduleDate(r.FormValue("date_end"), tz),
				Active:        active,
			},
		}
		if description != "" {
			req.Data.Description = &description
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
		if deps.GetPriceScheduleInUseIDs != nil {
			if inUse, _ := deps.GetPriceScheduleInUseIDs(ctx, []string{id}); inUse[id] {
				return centymo.HTMXError(deps.Labels.Errors.InUse)
			}
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
		ids := viewCtx.Request.Form["id"]
		var inUse map[string]bool
		if deps.GetPriceScheduleInUseIDs != nil && len(ids) > 0 {
			inUse, _ = deps.GetPriceScheduleInUseIDs(ctx, ids)
		}
		for _, id := range ids {
			if id == "" || inUse[id] {
				continue
			}
			_, _ = deps.DeletePriceSchedule(ctx, &priceschedulepb.DeletePriceScheduleRequest{Data: &priceschedulepb.PriceSchedule{Id: id}})
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
				Id:            id,
				Name:          record.GetName(),
				Description:   record.Description,
				DateTimeStart: record.GetDateTimeStart(),
				DateTimeEnd:   record.GetDateTimeEnd(),
				Active:        status == "active",
				LocationId:    record.LocationId,
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
					Id:            id,
					Name:          record.GetName(),
					Description:   record.Description,
					DateTimeStart: record.GetDateTimeStart(),
					DateTimeEnd:   record.GetDateTimeEnd(),
					Active:        status == "active",
					LocationId:    record.LocationId,
				},
			})
		}
		return centymo.HTMXSuccess("price-schedules-table")
	})
}
