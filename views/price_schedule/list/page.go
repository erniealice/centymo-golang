package list

import (
	"context"
	"log"
	"strconv"
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
)

type ListViewDeps struct {
	Routes                   centymo.PriceScheduleRoutes
	ListPriceSchedules       func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	ListLocations            func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	Labels                   centymo.PriceScheduleLabels
	CommonLabels             pyeza.CommonLabels
	TableLabels              types.TableLabels
	GetPriceScheduleInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// 2026-04-27 plan-client-scope plan §6.1 — optional client name lookup
	// for the optional Client column when scope ≠ master.
	ListClientNames func(ctx context.Context) map[string]string
}

// ScopeChip represents one chip in the scope filter bar on the list page.
// 2026-04-27 plan-client-scope plan §6.1.
type ScopeChip struct {
	Label  string
	URL    string
	Active bool
}

// ScopeFilterData holds the rendered chip group data for the price schedule list toolbar.
type ScopeFilterData struct {
	Label  string
	Active string
	Master ScopeChip
	Client ScopeChip
	All    ScopeChip
}

type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	ScopeFilter     ScopeFilterData
}

var priceScheduleAllowedSortCols = []string{"date_created", "date_modified", "name", "status"}
var priceScheduleSearchFields = []string{"name", "description"}

// buildScopeFilterData constructs the three scope chips for the price schedule list toolbar.
// baseListURL is the fully-resolved list URL without query params.
// clientFilter is preserved across chip clicks when set.
func buildScopeFilterData(baseListURL, scope, clientFilter string, l centymo.PriceScheduleLabels) ScopeFilterData {
	chipURL := func(s string) string {
		u := baseListURL + "?scope=" + s
		if clientFilter != "" {
			u += "&client_id=" + clientFilter
		}
		return u
	}
	return ScopeFilterData{
		Label:  l.Filters.ScopeChipLabel,
		Active: scope,
		Master: ScopeChip{Label: l.Filters.ScopeMaster, URL: chipURL("master"), Active: scope == "master"},
		Client: ScopeChip{Label: l.Filters.ScopeClient, URL: chipURL("client"), Active: scope == "client"},
		All:    ScopeChip{Label: l.Filters.ScopeAll, URL: chipURL("all"), Active: scope == "all"},
	}
}

func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		p, err := espynahttp.ParseTableParams(viewCtx.Request, priceScheduleAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}
		scope := normalizeScope(viewCtx.Request.URL.Query().Get("scope"))
		clientFilter := viewCtx.Request.URL.Query().Get("client_id")
		tableConfig, err := buildTableConfig(ctx, deps, status, scope, clientFilter, p)
		if err != nil {
			return view.Error(err)
		}

		baseListURL := route.ResolveURL(deps.Routes.ListURL, "status", status)
		scopeFilter := buildScopeFilterData(baseListURL, scope, clientFilter, deps.Labels)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   statusSubNav(deps.Routes.ActiveSubNav, status),
				HeaderTitle:    statusPageTitle(deps.Labels, status),
				HeaderSubtitle: statusPageCaption(deps.Labels, status),
				HeaderIcon:     "icon-calendar",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "price-schedule-list-content",
			Table:           tableConfig,
			ScopeFilter:     scopeFilter,
		}

		return view.OK("price-schedule-list", pageData)
	})
}

func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		p, err := espynahttp.ParseTableParams(viewCtx.Request, priceScheduleAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}
		scope := normalizeScope(viewCtx.Request.URL.Query().Get("scope"))
		clientFilter := viewCtx.Request.URL.Query().Get("client_id")
		tableConfig, err := buildTableConfig(ctx, deps, status, scope, clientFilter, p)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("table-card", tableConfig)
	})
}

// normalizeScope coerces an arbitrary client-scope filter value to one of
// the three supported chips: "master" (default), "client", "all".
//
// 2026-04-27 plan-client-scope plan §6.1.
func normalizeScope(s string) string {
	switch s {
	case "client", "all":
		return s
	default:
		return "master"
	}
}

// applyClientScopeFilter narrows the schedule list per §6.1.
func applyClientScopeFilter(items []*priceschedulepb.PriceSchedule, scope, clientFilter string) []*priceschedulepb.PriceSchedule {
	out := make([]*priceschedulepb.PriceSchedule, 0, len(items))
	for _, ps := range items {
		cid := ps.GetClientId()
		switch scope {
		case "master":
			if cid != "" {
				continue
			}
		case "client":
			if cid == "" {
				continue
			}
		}
		if clientFilter != "" && cid != clientFilter {
			continue
		}
		out = append(out, ps)
	}
	return out
}

func buildTableConfig(ctx context.Context, deps *ListViewDeps, status, scope, clientFilter string, p espynahttp.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)
	listParams := espynahttp.ToListParams(p, priceScheduleSearchFields)

	// Inject status filter for server-side pagination
	activeValue := status != "inactive"
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: activeValue},
		},
	})

	resp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list price schedules: %v", err)
		return nil, err
	}

	// 2026-04-27 plan-client-scope plan §6.1 — client-scope filter chip.
	filtered := applyClientScopeFilter(resp.GetData(), scope, clientFilter)

	var inUseIDs map[string]bool
	if deps.GetPriceScheduleInUseIDs != nil {
		var itemIDs []string
		for _, item := range filtered {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetPriceScheduleInUseIDs(ctx, itemIDs)
	}

	// Build location name lookup map
	locationNames := map[string]string{}
	if deps.ListLocations != nil {
		locResp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
		if err != nil {
			log.Printf("Failed to list locations for price schedule table: %v", err)
		} else {
			for _, loc := range locResp.GetData() {
				locationNames[loc.GetId()] = loc.GetName()
			}
		}
	}

	clientNames := map[string]string{}
	if deps.ListClientNames != nil && scope != "master" {
		clientNames = deps.ListClientNames(ctx)
	}

	l := deps.Labels
	includeClientColumn := scope != "master"
	columns := priceScheduleColumns(l, includeClientColumn)
	rows := buildTableRows(ctx, filtered, status, l, deps.CommonLabels, deps.Routes, inUseIDs, perms, locationNames, clientNames, includeClientColumn)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.CommonLabels, status, deps.Routes)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	var primaryAction *types.PrimaryAction
	if status == "active" {
		primaryAction = &types.PrimaryAction{
			Label:           l.Buttons.Add,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("price_schedule", "create"),
			DisabledTooltip: l.Errors.Unauthorized,
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "price-schedules-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   statusEmptyTitle(l, status),
			Message: statusEmptyMessage(l, status),
		},
		PrimaryAction: primaryAction,
		BulkActions:   &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig, nil
}

func priceScheduleColumns(l centymo.PriceScheduleLabels, includeClient bool) []types.TableColumn {
	cols := []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true, Filterable: true, FilterType: types.FilterTypeString},
		{Key: "description", Label: l.Columns.Description, Sortable: false},
		{Key: "date_start", Label: l.Columns.DateStart, Sortable: true, WidthClass: "col-2xl"},
		{Key: "date_end", Label: l.Columns.DateEnd, Sortable: true, WidthClass: "col-2xl"},
		{Key: "location", Label: l.Columns.Location, Sortable: false},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Filterable: false, WidthClass: "col-2xl"},
	}
	if includeClient {
		cols = append(cols, types.TableColumn{Key: "client", Label: l.Form.ClientLabel, Sortable: false, WidthClass: "col-3xl"})
	}
	return cols
}

func buildTableRows(ctx context.Context, priceSchedules []*priceschedulepb.PriceSchedule, status string, l centymo.PriceScheduleLabels, cl pyeza.CommonLabels, routes centymo.PriceScheduleRoutes, inUseIDs map[string]bool, perms *types.UserPermissions, locationNames map[string]string, clientNames map[string]string, includeClient bool) []types.TableRow {
	tz := types.LocationFromContext(ctx)
	rows := []types.TableRow{}
	for _, ps := range priceSchedules {
		recordStatus := "active"
		if !ps.GetActive() {
			recordStatus = "inactive"
		}

		id := ps.GetId()
		name := ps.GetName()
		description := ps.GetDescription()
		dateStart := types.FormatTimestampInTZ(ps.GetDateTimeStart(), tz, types.DateTimeReadable)
		dateEnd := types.FormatTimestampInTZ(ps.GetDateTimeEnd(), tz, types.DateTimeReadable)

		locationName := "—"
		if locID := ps.GetLocationId(); locID != "" {
			if n, ok := locationNames[locID]; ok && n != "" {
				locationName = n
			} else {
				locationName = locID
			}
		}

		isInUse := inUseIDs[id]

		clientID := ps.GetClientId()
		clientLabel := ""
		if clientID != "" {
			if n, ok := clientNames[clientID]; ok {
				clientLabel = n
			} else {
				clientLabel = clientID
			}
		}

		cells := []types.TableCell{
			{Type: "text", Value: name},
			{Type: "text", Value: description},
			{Type: "datetime", Value: dateStart},
			{Type: "datetime", Value: dateEnd},
			{Type: "text", Value: locationName},
			{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
		}
		if includeClient {
			if clientLabel != "" {
				cells = append(cells, types.TableCell{Type: "badge", Value: clientLabel, Variant: "info"})
			} else {
				cells = append(cells, types.TableCell{Type: "text", Value: ""})
			}
		}

		rows = append(rows, types.TableRow{
			ID:    id,
			Cells: cells,
			DataAttrs: map[string]string{
				"name":      name,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
				"client_id": clientID,
			},
			Actions: buildRowActions(id, name, ps.GetActive(), isInUse, l, cl, routes, perms),
		})
	}
	return rows
}

func buildRowActions(id, name string, active, isInUse bool, l centymo.PriceScheduleLabels, cl pyeza.CommonLabels, routes centymo.PriceScheduleRoutes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Buttons.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		{Type: "edit", Label: l.Buttons.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Buttons.Edit,
			Disabled: !perms.Can("price_schedule", "update"), DisabledTooltip: l.Errors.Unauthorized},
	}

	if active {
		actions = append(actions, types.TableAction{
			Type:            "clone",
			Label:           cl.Actions.Clone,
			Action:          "clone",
			URL:             route.ResolveURL(routes.EditURL, "id", id),
			DrawerTitle:     cl.Actions.Clone,
			Disabled:        !perms.Can("price_schedule", "create"),
			DisabledTooltip: l.Errors.Unauthorized,
		})
		actions = append(actions, types.TableAction{
			Type: "deactivate", Label: l.Buttons.Deactivate, Action: "deactivate",
			URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
			ConfirmTitle:    l.Confirm.DeactivateTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.DeactivateMessage, "{{name}}", name),
			Disabled:        !perms.Can("price_schedule", "update"),
			DisabledTooltip: l.Errors.Unauthorized,
		})
	} else {
		actions = append(actions, types.TableAction{
			Type: "activate", Label: l.Buttons.Activate, Action: "activate",
			URL: routes.SetStatusURL + "?status=active", ItemName: name,
			ConfirmTitle:    l.Confirm.ActivateTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.ActivateMessage, "{{name}}", name),
			Disabled:        !perms.Can("price_schedule", "update"),
			DisabledTooltip: l.Errors.Unauthorized,
		})
	}

	deleteAction := types.TableAction{
		Type:     "delete",
		Label:    l.Buttons.Delete,
		Action:   "delete",
		URL:      routes.DeleteURL,
		ItemName: name,
	}
	if isInUse {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = l.Errors.InUse
	} else if !perms.Can("price_schedule", "delete") {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = l.Errors.Unauthorized
	}
	actions = append(actions, deleteAction)
	return actions
}

func buildBulkActions(l centymo.PriceScheduleLabels, cl pyeza.CommonLabels, status string, routes centymo.PriceScheduleRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           cl.Bulk.Deactivate,
			Icon:            "icon-pause",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Bulk.DeactivateTitle,
			ConfirmMessage:  l.Bulk.DeactivateMessage,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           cl.Bulk.Activate,
			Icon:            "icon-play",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Bulk.ActivateTitle,
			ConfirmMessage:  l.Bulk.ActivateMessage,
			ExtraParamsJSON: `{"target_status":"active"}`,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:              "delete",
		Label:            cl.Bulk.Delete,
		Icon:             "icon-trash-2",
		Variant:          "danger",
		Endpoint:         routes.BulkDeleteURL,
		ConfirmTitle:     l.Bulk.DeleteTitle,
		ConfirmMessage:   l.Bulk.DeleteMessage,
		RequiresDataAttr: "deletable",
	})

	return actions
}

func statusPageTitle(l centymo.PriceScheduleLabels, status string) string {
	switch status {
	case "active":
		return l.Page.ActiveTitle
	case "inactive":
		return l.Page.InactiveTitle
	default:
		return l.Page.Title
	}
}

func statusPageCaption(l centymo.PriceScheduleLabels, status string) string {
	return l.Page.Subtitle
}

func statusEmptyTitle(l centymo.PriceScheduleLabels, status string) string {
	return l.Empty.Title
}

func statusEmptyMessage(l centymo.PriceScheduleLabels, status string) string {
	return l.Empty.Message
}

func statusVariant(status string) string {
	switch status {
	case "active":
		return "success"
	case "inactive":
		return "warning"
	default:
		return "default"
	}
}

func statusSubNav(base, status string) string {
	if base == "" {
		return status
	}
	return base + "-" + status
}
