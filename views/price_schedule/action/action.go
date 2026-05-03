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

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/erniealice/centymo-golang/views/price_schedule/form"
)

var clientNameSort = &commonpb.SortRequest{
	Fields: []*commonpb.SortField{{Field: "name", Direction: commonpb.SortDirection_ASC}},
}

// splitScheduleDateTimeForInputs renders ts in tz as a (date, time) pair
// suitable for the drawer's date+time input grid. Nil ts → ("", "").
func splitScheduleDateTimeForInputs(ts *timestamppb.Timestamp, tz *time.Location) (date, t string) {
	if ts == nil || !ts.IsValid() {
		return "", ""
	}
	moment := ts.AsTime().In(tz)
	return moment.Format(pyezatypes.DateInputLayout), moment.Format(pyezatypes.TimeInputLayout)
}

// parseScheduleDateTime combines a date input (YYYY-MM-DD) and an OPTIONAL
// time input (HH:MM) into a UTC timestamp, anchored to tz.
//
// When time is empty:
//   - For start-of-range (isEnd=false): defaults to 00:00:00 (start-of-day).
//   - For end-of-range (isEnd=true): defaults to 23:59:59 (end-of-day) so an
//     "end" date without a time still includes the full day.
//
// Empty date → nil. The 2026-04-28 date+time field plan §4 anchors to
// types.LocationFromContext(ctx) for the operator's display timezone.
func parseScheduleDateTime(date, t string, tz *time.Location, isEnd bool) *timestamppb.Timestamp {
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

type Deps struct {
	Routes                   centymo.PriceScheduleRoutes
	Labels                   centymo.PriceScheduleLabels
	CreatePriceSchedule      func(ctx context.Context, req *priceschedulepb.CreatePriceScheduleRequest) (*priceschedulepb.CreatePriceScheduleResponse, error)
	ReadPriceSchedule        func(ctx context.Context, req *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error)
	UpdatePriceSchedule      func(ctx context.Context, req *priceschedulepb.UpdatePriceScheduleRequest) (*priceschedulepb.UpdatePriceScheduleResponse, error)
	DeletePriceSchedule      func(ctx context.Context, req *priceschedulepb.DeletePriceScheduleRequest) (*priceschedulepb.DeletePriceScheduleResponse, error)
	ListLocations            func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	GetPriceScheduleInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// 2026-04-27 plan-client-scope plan §6.7 / §4.4.1.
	ListClients      func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	SearchClientsURL string
}

func loadLocations(ctx context.Context, deps *Deps) []*form.LocationOption {
	if deps.ListLocations == nil {
		return nil
	}
	resp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
	if err != nil {
		return nil
	}
	opts := make([]*form.LocationOption, 0, len(resp.GetData()))
	for _, loc := range resp.GetData() {
		opts = append(opts, &form.LocationOption{Id: loc.GetId(), Name: loc.GetName()})
	}
	return opts
}

// loadClientOptions fetches the workspace's clients and converts them into
// auto-complete options. Returns nil when the dep is unwired.
//
// 2026-04-27 plan-client-scope plan §6.7 / §4.4.1.
func loadClientOptions(ctx context.Context, deps *Deps, selectedID string) []map[string]any {
	if deps.ListClients == nil {
		return nil
	}
	resp, err := deps.ListClients(ctx, &clientpb.ListClientsRequest{Sort: clientNameSort})
	if err != nil {
		return nil
	}
	opts := make([]map[string]any, 0, len(resp.GetData()))
	for _, c := range resp.GetData() {
		label := c.GetName()
		if label == "" {
			if u := c.GetUser(); u != nil {
				label = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			}
		}
		if label == "" {
			label = c.GetId()
		}
		opts = append(opts, map[string]any{
			"Value":    c.GetId(),
			"Label":    label,
			"Selected": c.GetId() == selectedID,
		})
	}
	return opts
}

// resolveClientName looks up a client_id in the workspace and returns its
// display name. Falls back to the rep full name and finally the bare ID,
// mirroring resolveClientBreadcrumb in the subscription detail view.
func resolveClientName(ctx context.Context, deps *Deps, clientID string) string {
	if clientID == "" || deps.ListClients == nil {
		return ""
	}
	resp, err := deps.ListClients(ctx, &clientpb.ListClientsRequest{Sort: clientNameSort})
	if err != nil {
		return clientID
	}
	for _, c := range resp.GetData() {
		if c.GetId() != clientID {
			continue
		}
		if name := c.GetName(); name != "" {
			return name
		}
		if u := c.GetUser(); u != nil {
			full := strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			if full != "" {
				return full
			}
		}
		return clientID
	}
	return clientID
}

func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule", "create") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			// 2026-04-27 plan-client-scope plan §4.4.1 — when the request
			// asks for a name suggestion, render only the Name input partial
			// so the schedule-drawer's HTMX picker swap can update the name
			// without reloading the entire form.
			// 2026-05-03 — Trailing timestamp ("YYYY-MM-DD HH:MM" in operator
			// TZ) appended to derived names so duplicate-name collisions are
			// avoided. Same suffix used by both the suggest-name partial and
			// the initial render.
			tz := pyezatypes.LocationFromContext(ctx)
			derivedTimestamp := time.Now().In(tz).Format("2006-01-02 15:04:05")

			if viewCtx.Request.URL.Query().Get("suggest_name") == "1" {
				clientID := viewCtx.Request.URL.Query().Get("client_id")
				clientName := resolveClientName(ctx, deps, clientID)
				derived := form.BuildDerivedScheduleName(clientName, deps.Labels.Form.CustomClientPriceScheduleLabelSuffix, derivedTimestamp)
				return view.OK("price-schedule-name-suggest", map[string]any{
					"Value":           derived,
					"NamePlaceholder": deps.Labels.Form.NamePlaceholder,
				})
			}

			locations := loadLocations(ctx, deps)
			// 2026-04-27 plan-client-scope plan §4.4.1 / §6.7 — Client picker
			// + name pre-fill via HTMX swap when a client is selected.
			pinnedClientID := viewCtx.Request.URL.Query().Get("client_id")
			clientLabel := resolveClientName(ctx, deps, pinnedClientID)
			// 2026-05-03 — Pre-fill the name in BOTH contexts so duplicates
			// are unlikely. Client context cascades the client name; standalone
			// uses just the lyngua suffix + timestamp.
			defaultName := form.BuildDerivedScheduleName(clientLabel, deps.Labels.Form.CustomClientPriceScheduleLabelSuffix, derivedTimestamp)
			// 2026-04-28 — Scope radio default. `location` unless the URL
			// pins a client (?client_id=...) which implies client scope.
			scope := "location"
			clientLocked := false
			if pinnedClientID != "" {
				scope = "client"
				clientLocked = true
			}
			return view.OK("price-schedule-drawer-form", &form.Data{
				FormAction:      deps.Routes.AddURL,
				Active:          true,
				Name:            defaultName,
				Locations:       locations,
				LocationOptions: form.BuildLocationAutoCompleteOptions(locations, ""),
				ClientID:        pinnedClientID,
				ClientLabel:     clientLabel,
				ClientOptions:   loadClientOptions(ctx, deps, pinnedClientID),
				SearchClientURL: deps.SearchClientsURL,
				SuggestNameURL:  deps.Routes.AddURL + "?suggest_name=1",
				Scope:           scope,
				ClientLocked:    clientLocked,
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
		clientID := strings.TrimSpace(r.FormValue("client_id"))
		// 2026-04-28 — Scope radio enforces mutual exclusion server-side:
		// when scope=location, drop any submitted client_id; when scope=client,
		// drop any submitted location_id. Defaults to "location" so legacy
		// callers without a scope field keep prior behaviour.
		scope := r.FormValue("scope")
		if scope == "" {
			scope = "location"
		}
		if scope == "location" {
			clientID = ""
		} else if scope == "client" {
			locationID = ""
		}
		tz := pyezatypes.LocationFromContext(ctx)
		req := &priceschedulepb.CreatePriceScheduleRequest{
			Data: &priceschedulepb.PriceSchedule{
				Name:          r.FormValue("name"),
				DateTimeStart: parseScheduleDateTime(r.FormValue("date_start_date"), r.FormValue("date_start_time"), tz, false),
				DateTimeEnd:   parseScheduleDateTime(r.FormValue("date_end_date"), r.FormValue("date_end_time"), tz, true),
				Active:        active,
			},
		}
		if description != "" {
			req.Data.Description = &description
		}
		if locationID != "" {
			req.Data.LocationId = &locationID
		}
		if clientID != "" {
			req.Data.ClientId = &clientID
		}
		if _, err := deps.CreatePriceSchedule(ctx, req); err != nil {
			log.Printf("Failed to create price schedule: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedules-table")
	})
}

// NewSuggestNameAction renders the per-tier "{ClientName} - {suffix}" name
// as a partial HTML <input> for the schedule add drawer's Client picker
// HTMX swap (plan §4.4.1 fallback path). Idempotent GET.
//
// Wired at GET deps.Routes.AddURL with `?suggest_name=1`. Centymo block
// registers it alongside the regular Add handler.
func NewSuggestNameAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		clientID := viewCtx.Request.URL.Query().Get("client_id")
		clientName := resolveClientName(ctx, deps, clientID)
		tz := pyezatypes.LocationFromContext(ctx)
		derivedTimestamp := time.Now().In(tz).Format("2006-01-02 15:04:05")
		derived := form.BuildDerivedScheduleName(clientName, deps.Labels.Form.CustomClientPriceScheduleLabelSuffix, derivedTimestamp)
		return view.OK("price-schedule-name-suggest", map[string]any{
			"Value":           derived,
			"NamePlaceholder": deps.Labels.Form.NamePlaceholder,
		})
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
			selectedClientID := record.GetClientId()
			clientLabel := resolveClientName(ctx, deps, selectedClientID)
			// 2026-04-28 — Scope detection on edit: a populated client_id
			// means the schedule is client-scoped; otherwise location-scoped.
			scope := "location"
			if selectedClientID != "" {
				scope = "client"
			}
			startDate, startTime := splitScheduleDateTimeForInputs(record.GetDateTimeStart(), tz)
			endDate, endTime := splitScheduleDateTimeForInputs(record.GetDateTimeEnd(), tz)
			return view.OK("price-schedule-drawer-form", &form.Data{
				FormAction:            formAction,
				IsEdit:                !isClone,
				ID:                    formID,
				Name:                  name,
				Description:           record.GetDescription(),
				DateStartDate:         startDate,
				DateStartTime:         startTime,
				DateEndDate:           endDate,
				DateEndTime:           endTime,
				Active:                record.GetActive(),
				SelectedLocationID:    selectedLocationID,
				SelectedLocationLabel: form.FindLocationLabel(locations, selectedLocationID),
				Locations:             locations,
				LocationOptions:       form.BuildLocationAutoCompleteOptions(locations, selectedLocationID),
				// 2026-04-27 plan-client-scope plan §6.7 — Client picker.
				ClientID:        selectedClientID,
				ClientLabel:     clientLabel,
				ClientOptions:   loadClientOptions(ctx, deps, selectedClientID),
				SearchClientURL: deps.SearchClientsURL,
				SuggestNameURL:  deps.Routes.AddURL + "?suggest_name=1",
				Scope:           scope,
				Labels:          deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		locationID := r.FormValue("location_id")
		description := r.FormValue("description")
		clientID := strings.TrimSpace(r.FormValue("client_id"))
		// 2026-04-28 — Mirror the Add POST scope-driven mutual exclusion.
		scope := r.FormValue("scope")
		if scope == "" {
			scope = "location"
		}
		if scope == "location" {
			clientID = ""
		} else if scope == "client" {
			locationID = ""
		}
		tz := pyezatypes.LocationFromContext(ctx)
		req := &priceschedulepb.UpdatePriceScheduleRequest{
			Data: &priceschedulepb.PriceSchedule{
				Id:            id,
				Name:          r.FormValue("name"),
				DateTimeStart: parseScheduleDateTime(r.FormValue("date_start_date"), r.FormValue("date_start_time"), tz, false),
				DateTimeEnd:   parseScheduleDateTime(r.FormValue("date_end_date"), r.FormValue("date_end_time"), tz, true),
				Active:        active,
			},
		}
		if description != "" {
			req.Data.Description = &description
		}
		// 2026-04-28 — Always set both pointers explicitly so that flipping
		// scope on edit actively clears the inverse FK rather than leaving
		// the previous value behind.
		req.Data.LocationId = &locationID
		req.Data.ClientId = &clientID
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
