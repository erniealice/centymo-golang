package list

import (
	"context"
	"fmt"
	"log"
	"strings"

	sgwu "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_workspace_user"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/shared/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	sgwupb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_workspace_user"
)

// ListViewDeps holds view dependencies for the list and table views.
type ListViewDeps struct {
	Routes                                    sgwu.Routes
	ListSubscriptionGroupWorkspaceUsers       func(ctx context.Context, req *sgwupb.ListSubscriptionGroupWorkspaceUsersRequest) (*sgwupb.ListSubscriptionGroupWorkspaceUsersResponse, error)
	Labels                                    sgwu.Labels
	CommonLabels                              pyeza.CommonLabels
	TableLabels                               types.TableLabels
	GetSubscriptionGroupWorkspaceUserInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// PageData holds data for the subscription_group_workspace_user list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var sgwuSearchFields = []string{"scope", "role"}

// NewView creates the full-page list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_workspace_user", "list") {
			return view.Forbidden("subscription_group_workspace_user:list")
		}
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := sgwuColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "scope", "asc")
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
			ContentTemplate: "subscription-group-workspace-user-list-content",
			Table:           tableConfig,
		}

		return view.OK("subscription-group-workspace-user-list", pageData)
	})
}

// NewTableView creates the HTMX table-only view.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := sgwuColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "scope", "asc")
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
	listParams := espynahttp.ToListParams(p, sgwuSearchFields)

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

	resp, err := deps.ListSubscriptionGroupWorkspaceUsers(ctx, &sgwupb.ListSubscriptionGroupWorkspaceUsersRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list subscription_group_workspace_users: %v", err)
		return nil, err
	}

	items := resp.GetData()

	var inUseIDs map[string]bool
	if deps.GetSubscriptionGroupWorkspaceUserInUseIDs != nil {
		var itemIDs []string
		for _, item := range items {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetSubscriptionGroupWorkspaceUserInUseIDs(ctx, itemIDs)
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
			Disabled:        !perms.Can("subscription_group_workspace_user", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "subscription_group_workspace_user:create"),
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "subscription-group-workspace-users-table",
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
		DefaultSortColumn:    "scope",
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

func sgwuColumns(l sgwu.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "workspace_user_id", Label: l.Columns.WorkspaceUser},
		{Key: "subscription_group_id", Label: l.Columns.SubscriptionGroup, WidthClass: "col-2xl"},
		{Key: "scope", Label: l.Columns.Scope, WidthClass: "col-2xl"},
		{Key: "role", Label: l.Columns.Role, WidthClass: "col-xl"},
		{Key: "is_owner", Label: l.Columns.IsOwner, NoSort: true, NoFilter: true, WidthClass: "col-xl"},
	}
}

func buildTableRows(items []*sgwupb.SubscriptionGroupWorkspaceUser, status string, l sgwu.Labels, cl pyeza.CommonLabels, routes sgwu.Routes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, item := range items {
		recordStatus := "active"
		if !item.GetActive() {
			recordStatus = "inactive"
		}

		id := item.GetId()
		workspaceUserID := item.GetWorkspaceUserId()
		subscriptionGroupID := item.GetSubscriptionGroupId()
		scopeVal := item.GetScope()
		if scopeVal == "" {
			scopeVal = l.Detail.NoScope
		}
		roleVal := item.GetRole()
		if roleVal == "" {
			roleVal = l.Detail.NoRole
		}
		ownerLabel := l.Detail.OwnerNo
		if item.GetIsOwner() {
			ownerLabel = l.Detail.OwnerYes
		}

		isInUse := inUseIDs[id]

		// Use workspace_user_id as the display name for confirm messages.
		displayName := workspaceUserID
		if displayName == "" {
			displayName = id
		}

		cells := []types.TableCell{
			{Type: "text", Value: workspaceUserID},
			{Type: "text", Value: subscriptionGroupID},
			{Type: "text", Value: scopeVal},
			{Type: "text", Value: roleVal},
			{Type: "text", Value: ownerLabel},
		}

		rows = append(rows, types.TableRow{
			ID:    id,
			Cells: cells,
			DataAttrs: map[string]string{
				"name":      displayName,
				"status":    recordStatus,
				"deletable": fmt.Sprintf("%v", !isInUse),
			},
			Actions: buildRowActions(id, displayName, item.GetActive(), isInUse, l, cl, routes, perms),
		})
	}
	return rows
}

func buildRowActions(id, displayName string, active, isInUse bool, l sgwu.Labels, cl pyeza.CommonLabels, routes sgwu.Routes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Buttons.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		{Type: "edit", Label: l.Buttons.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Buttons.Edit,
			Disabled: !perms.Can("subscription_group_workspace_user", "update"), DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "subscription_group_workspace_user:update")},
	}

	if active {
		actions = append(actions, types.TableAction{
			Type:            "clone",
			Label:           cl.Actions.Clone,
			Action:          "clone",
			URL:             route.ResolveURL(routes.EditURL, "id", id),
			DrawerTitle:     cl.Actions.Clone,
			Disabled:        !perms.Can("subscription_group_workspace_user", "create"),
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "subscription_group_workspace_user:create"),
		})
		actions = append(actions, types.TableAction{
			Type: "deactivate", Label: l.Buttons.Deactivate, Action: "deactivate",
			URL: routes.SetStatusURL + "?status=inactive", ItemName: displayName,
			ConfirmTitle:    l.Confirm.DeactivateTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.DeactivateMessage, "{{name}}", displayName),
			Disabled:        !perms.Can("subscription_group_workspace_user", "update"),
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "subscription_group_workspace_user:update"),
		})
	} else {
		actions = append(actions, types.TableAction{
			Type: "activate", Label: l.Buttons.Activate, Action: "activate",
			URL: routes.SetStatusURL + "?status=active", ItemName: displayName,
			ConfirmTitle:    l.Confirm.ActivateTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.ActivateMessage, "{{name}}", displayName),
			Disabled:        !perms.Can("subscription_group_workspace_user", "update"),
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "subscription_group_workspace_user:update"),
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
	} else if !perms.Can("subscription_group_workspace_user", "delete") {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = fmt.Sprintf(cl.Errors.MissingPermission, "subscription_group_workspace_user:delete")
	}
	actions = append(actions, deleteAction)
	return actions
}

func buildBulkActions(l sgwu.Labels, cl pyeza.CommonLabels, status string, routes sgwu.Routes) []types.BulkAction {
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

func statusPageTitle(l sgwu.Labels, status string) string {
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
