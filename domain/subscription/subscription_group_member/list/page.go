package list

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	subscription_group_member "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_member"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/shared/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	subscriptiongroupmemberpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_member"
)

// ListViewDeps holds dependencies for the subscription_group_member list views.
type ListViewDeps struct {
	Routes                             subscription_group_member.Routes
	ListSubscriptionGroupMembers       func(ctx context.Context, req *subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest) (*subscriptiongroupmemberpb.ListSubscriptionGroupMembersResponse, error)
	Labels                             subscription_group_member.Labels
	CommonLabels                       pyeza.CommonLabels
	TableLabels                        types.TableLabels
	GetSubscriptionGroupMemberInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// PageData holds the data for the subscription_group_member list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var memberSearchFields = []string{"subscription_group_id", "subscription_id", "client_id"}

// NewView creates the full-page subscription_group_member list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group_member", "list") {
			return view.Forbidden("subscription_group_member:list")
		}
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := memberColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "date_created", "desc")
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
			ContentTemplate: "subscription-group-member-list-content",
			Table:           tableConfig,
		}

		return view.OK("subscription-group-member-list", pageData)
	})
}

// NewTableView creates the HTMX table-only subscription_group_member list view.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := memberColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "date_created", "desc")
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
	listParams := espynahttp.ToListParams(p, memberSearchFields)

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

	resp, err := deps.ListSubscriptionGroupMembers(ctx, &subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		return nil, err
	}

	items := resp.GetData()

	var inUseIDs map[string]bool
	if deps.GetSubscriptionGroupMemberInUseIDs != nil {
		var itemIDs []string
		for _, item := range items {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetSubscriptionGroupMemberInUseIDs(ctx, itemIDs)
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
			Disabled:        !perms.Can("subscription_group_member", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "subscription_group_member:create"),
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "subscription-group-members-table",
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
		DefaultSortColumn:    "date_created",
		DefaultSortDirection: "desc",
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

func memberColumns(l subscription_group_member.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "subscription_group_id", Label: l.Columns.SubscriptionGroupId},
		{Key: "subscription_id", Label: l.Columns.SubscriptionId},
		{Key: "client_id", Label: l.Columns.ClientId},
	}
}

func buildTableRows(members []*subscriptiongroupmemberpb.SubscriptionGroupMember, status string, l subscription_group_member.Labels, cl pyeza.CommonLabels, routes subscription_group_member.Routes, inUseIDs map[string]bool, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, m := range members {
		recordStatus := "active"
		if !m.GetActive() {
			recordStatus = "inactive"
		}

		id := m.GetId()
		groupID := m.GetSubscriptionGroupId()
		subID := m.GetSubscriptionId()
		clientID := m.GetClientId()

		isInUse := inUseIDs[id]

		cells := []types.TableCell{
			{Type: "text", Value: groupID},
			{Type: "text", Value: subID},
			{Type: "text", Value: clientID},
		}

		rows = append(rows, types.TableRow{
			ID:    id,
			Cells: cells,
			DataAttrs: map[string]string{
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: buildRowActions(id, id, m.GetActive(), isInUse, l, cl, routes, perms),
		})
	}
	return rows
}

func buildRowActions(id, name string, active, isInUse bool, l subscription_group_member.Labels, cl pyeza.CommonLabels, routes subscription_group_member.Routes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Buttons.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		{Type: "edit", Label: l.Buttons.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Buttons.Edit,
			Disabled: !perms.Can("subscription_group_member", "update"), DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "subscription_group_member:update")},
	}

	if active {
		actions = append(actions, types.TableAction{
			Type: "deactivate", Label: l.Buttons.Deactivate, Action: "deactivate",
			URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
			ConfirmTitle:    l.Confirm.DeactivateTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.DeactivateMessage, "{{name}}", name),
			Disabled:        !perms.Can("subscription_group_member", "update"),
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "subscription_group_member:update"),
		})
	} else {
		actions = append(actions, types.TableAction{
			Type: "activate", Label: l.Buttons.Activate, Action: "activate",
			URL: routes.SetStatusURL + "?status=active", ItemName: name,
			ConfirmTitle:    l.Confirm.ActivateTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.ActivateMessage, "{{name}}", name),
			Disabled:        !perms.Can("subscription_group_member", "update"),
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "subscription_group_member:update"),
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
	} else if !perms.Can("subscription_group_member", "delete") {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = fmt.Sprintf(cl.Errors.MissingPermission, "subscription_group_member:delete")
	}
	actions = append(actions, deleteAction)
	return actions
}

func buildBulkActions(l subscription_group_member.Labels, cl pyeza.CommonLabels, status string, routes subscription_group_member.Routes) []types.BulkAction {
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

func statusPageTitle(l subscription_group_member.Labels, status string) string {
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
