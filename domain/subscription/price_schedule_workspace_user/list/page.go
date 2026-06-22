package list

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	pswu "github.com/erniealice/centymo-golang/domain/subscription/price_schedule_workspace_user"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	pswupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule_workspace_user"
)

type ListViewDeps struct {
	Routes                                pswu.Routes
	ListPriceScheduleWorkspaceUsers       func(ctx context.Context, req *pswupb.ListPriceScheduleWorkspaceUsersRequest) (*pswupb.ListPriceScheduleWorkspaceUsersResponse, error)
	Labels                                pswu.Labels
	CommonLabels                          pyeza.CommonLabels
	TableLabels                           types.TableLabels
	GetPriceScheduleWorkspaceUserInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var pswuSearchFields = []string{"price_schedule_id", "workspace_user_id", "scope", "role"}

func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_schedule_workspace_user", "list") {
			return view.Forbidden("price_schedule_workspace_user:list")
		}
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := pswuColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "price_schedule_id", "asc")
		if err != nil {
			return view.Error(err)
		}
		tableConfig, err := buildTableConfig(ctx, deps, status, columns, p)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   statusSubNav(deps.Routes.ActiveSubNav, status),
				HeaderTitle:    statusPageTitle(deps.Labels, status),
				HeaderSubtitle: deps.Labels.Page.Subtitle,
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "price-schedule-workspace-user-list-content",
			Table:           tableConfig,
		}

		return view.OK("price-schedule-workspace-user-list", pageData)
	})
}

func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := pswuColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "price_schedule_id", "asc")
		if err != nil {
			return view.Error(err)
		}
		tableConfig, err := buildTableConfig(ctx, deps, status, columns, p)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("table-card", tableConfig)
	})
}

func buildTableConfig(ctx context.Context, deps *ListViewDeps, status string, columns []types.TableColumn, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)
	listParams := espynahttp.ToListParams(p, pswuSearchFields)

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

	resp, err := deps.ListPriceScheduleWorkspaceUsers(ctx, &pswupb.ListPriceScheduleWorkspaceUsersRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list price_schedule_workspace_users: %v", err)
		return nil, err
	}

	items := resp.GetData()

	var inUseIDs map[string]bool
	if deps.GetPriceScheduleWorkspaceUserInUseIDs != nil {
		var itemIDs []string
		for _, item := range items {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetPriceScheduleWorkspaceUserInUseIDs(ctx, itemIDs)
	}

	l := deps.Labels
	rows := buildTableRows(items, status, l, deps.CommonLabels, deps.Routes, inUseIDs, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := pyeza.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.CommonLabels, status, deps.Routes)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	var primaryAction *types.PrimaryAction
	if status == "active" {
		primaryAction = &types.PrimaryAction{
			Label:           l.Buttons.Add,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("price_schedule_workspace_user", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "price_schedule_workspace_user:create"),
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "price-schedule-workspace-users-table",
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
		DefaultSortColumn:    "price_schedule_id",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: primaryAction,
		BulkActions:   &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig, nil
}

func pswuColumns(l pswu.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "price_schedule_id", Label: l.Columns.PriceScheduleId},
		{Key: "workspace_user_id", Label: l.Columns.WorkspaceUserId},
		{Key: "scope", Label: l.Columns.Scope, WidthClass: "col-2xl"},
		{Key: "role", Label: l.Columns.Role, WidthClass: "col-2xl"},
	}
}

func buildTableRows(items []*pswupb.PriceScheduleWorkspaceUser, status string, l pswu.Labels, cl pyeza.CommonLabels, routes pswu.Routes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, item := range items {
		recordStatus := "active"
		if !item.GetActive() {
			recordStatus = "inactive"
		}

		id := item.GetId()
		scheduleID := item.GetPriceScheduleId()
		userID := item.GetWorkspaceUserId()
		scope := item.GetScope()
		role := item.GetRole()

		isInUse := inUseIDs[id]

		cells := []types.TableCell{
			{Type: "text", Value: scheduleID},
			{Type: "text", Value: userID},
			{Type: "text", Value: scope},
			{Type: "text", Value: role},
		}

		rows = append(rows, types.TableRow{
			ID:    id,
			Cells: cells,
			DataAttrs: map[string]string{
				"id":        id,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: buildRowActions(id, scheduleID, item.GetActive(), isInUse, l, cl, routes, perms),
		})
	}
	return rows
}

func buildRowActions(id, displayName string, active, isInUse bool, l pswu.Labels, cl pyeza.CommonLabels, routes pswu.Routes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Buttons.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		{Type: "edit", Label: l.Buttons.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Buttons.Edit,
			Disabled: !perms.Can("price_schedule_workspace_user", "update"), DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "price_schedule_workspace_user:update")},
	}

	if active {
		actions = append(actions, types.TableAction{
			Type: "deactivate", Label: l.Buttons.Deactivate, Action: "deactivate",
			URL: routes.SetStatusURL + "?status=inactive", ItemName: displayName,
			ConfirmTitle:    l.Confirm.DeactivateTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.DeactivateMessage, "{{id}}", id),
			Disabled:        !perms.Can("price_schedule_workspace_user", "update"),
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "price_schedule_workspace_user:update"),
		})
	} else {
		actions = append(actions, types.TableAction{
			Type: "activate", Label: l.Buttons.Activate, Action: "activate",
			URL: routes.SetStatusURL + "?status=active", ItemName: displayName,
			ConfirmTitle:    l.Confirm.ActivateTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.ActivateMessage, "{{id}}", id),
			Disabled:        !perms.Can("price_schedule_workspace_user", "update"),
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "price_schedule_workspace_user:update"),
		})
	}

	deleteAction := types.TableAction{
		Type:     "delete",
		Label:    l.Buttons.Delete,
		Action:   "delete",
		URL:      routes.DeleteURL,
		ItemName: displayName,
	}
	if isInUse {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = l.Errors.InUse
	} else if !perms.Can("price_schedule_workspace_user", "delete") {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = fmt.Sprintf(cl.Errors.MissingPermission, "price_schedule_workspace_user:delete")
	}
	actions = append(actions, deleteAction)
	return actions
}

func buildBulkActions(l pswu.Labels, cl pyeza.CommonLabels, status string, routes pswu.Routes) []types.BulkAction {
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

func statusPageTitle(l pswu.Labels, status string) string {
	switch status {
	case "active":
		return l.Page.ActiveTitle
	case "inactive":
		return l.Page.InactiveTitle
	default:
		return l.Page.Title
	}
}

func statusSubNav(base, status string) string {
	if base == "" {
		return status
	}
	return base + "-" + status
}
