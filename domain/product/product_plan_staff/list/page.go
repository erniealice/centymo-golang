package list

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	product_plan_staff "github.com/erniealice/centymo-golang/domain/product/product_plan_staff"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/shared/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
)

type ListViewDeps struct {
	Routes                      product_plan_staff.Routes
	ListProductPlanStaffs       func(ctx context.Context, req *productplanstaffpb.ListProductPlanStaffsRequest) (*productplanstaffpb.ListProductPlanStaffsResponse, error)
	Labels                      product_plan_staff.Labels
	CommonLabels                pyeza.CommonLabels
	TableLabels                 types.TableLabels
	GetProductPlanStaffInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var productPlanStaffSearchFields = []string{"staff_id", "product_plan_id", "role"}

func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan_staff", "list") {
			return view.Forbidden("product_plan_staff:list")
		}
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := productPlanStaffColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "staff_id", "asc")
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
			ContentTemplate: "product-plan-staff-list-content",
			Table:           tableConfig,
		}

		return view.OK("product-plan-staff-list", pageData)
	})
}

func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := productPlanStaffColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "staff_id", "asc")
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
	listParams := espynahttp.ToListParams(p, productPlanStaffSearchFields)

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

	resp, err := deps.ListProductPlanStaffs(ctx, &productplanstaffpb.ListProductPlanStaffsRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list product plan staffs: %v", err)
		return nil, err
	}

	items := resp.GetData()

	var inUseIDs map[string]bool
	if deps.GetProductPlanStaffInUseIDs != nil {
		var itemIDs []string
		for _, item := range items {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetProductPlanStaffInUseIDs(ctx, itemIDs)
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
			Disabled:        !perms.Can("product_plan_staff", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "product_plan_staff:create"),
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "product-plan-staffs-table",
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
		DefaultSortColumn:    "staff_id",
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

func productPlanStaffColumns(l product_plan_staff.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "staff_id", Label: l.Columns.StaffID},
		{Key: "product_plan_id", Label: l.Columns.ProductPlan},
		{Key: "role", Label: l.Columns.Role, WidthClass: "col-2xl"},
	}
}

func buildTableRows(
	items []*productplanstaffpb.ProductPlanStaff,
	status string,
	l product_plan_staff.Labels,
	cl pyeza.CommonLabels,
	routes product_plan_staff.Routes,
	inUseIDs map[string]bool,
	perms *types.UserPermissions,
) []types.TableRow {
	rows := []types.TableRow{}
	for _, ps := range items {
		recordStatus := "active"
		if !ps.GetActive() {
			recordStatus = "inactive"
		}

		id := ps.GetId()
		staffID := ps.GetStaffId()
		productPlanID := ps.GetProductPlanId()
		role := ps.GetRole()
		if role == "" {
			role = l.Detail.NoRole
		}

		isInUse := inUseIDs[id]

		// Use staff_id as the display name for confirmation dialogs.
		displayName := staffID
		if displayName == "" {
			displayName = id
		}

		cells := []types.TableCell{
			{Type: "text", Value: staffID},
			{Type: "text", Value: productPlanID},
			{Type: "text", Value: role},
		}

		rows = append(rows, types.TableRow{
			ID:    id,
			Cells: cells,
			DataAttrs: map[string]string{
				"name":      displayName,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: buildRowActions(id, displayName, ps.GetActive(), isInUse, l, cl, routes, perms),
		})
	}
	return rows
}

func buildRowActions(id, name string, active, isInUse bool, l product_plan_staff.Labels, cl pyeza.CommonLabels, routes product_plan_staff.Routes, perms *types.UserPermissions) []types.TableAction {
	actions := []types.TableAction{
		{Type: "view", Label: l.Buttons.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		{Type: "edit", Label: l.Buttons.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Buttons.Edit,
			Disabled: !perms.Can("product_plan_staff", "update"), DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "product_plan_staff:update")},
	}

	if active {
		actions = append(actions, types.TableAction{
			Type:            "clone",
			Label:           cl.Actions.Clone,
			Action:          "clone",
			URL:             route.ResolveURL(routes.EditURL, "id", id),
			DrawerTitle:     cl.Actions.Clone,
			Disabled:        !perms.Can("product_plan_staff", "create"),
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "product_plan_staff:create"),
		})
		actions = append(actions, types.TableAction{
			Type: "deactivate", Label: l.Buttons.Deactivate, Action: "deactivate",
			URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
			ConfirmTitle:    l.Confirm.DeactivateTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.DeactivateMessage, "{{name}}", name),
			Disabled:        !perms.Can("product_plan_staff", "update"),
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "product_plan_staff:update"),
		})
	} else {
		actions = append(actions, types.TableAction{
			Type: "activate", Label: l.Buttons.Activate, Action: "activate",
			URL: routes.SetStatusURL + "?status=active", ItemName: name,
			ConfirmTitle:    l.Confirm.ActivateTitle,
			ConfirmMessage:  strings.ReplaceAll(l.Confirm.ActivateMessage, "{{name}}", name),
			Disabled:        !perms.Can("product_plan_staff", "update"),
			DisabledTooltip: fmt.Sprintf(cl.Errors.MissingPermission, "product_plan_staff:update"),
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
	} else if !perms.Can("product_plan_staff", "delete") {
		deleteAction.Disabled = true
		deleteAction.DisabledTooltip = fmt.Sprintf(cl.Errors.MissingPermission, "product_plan_staff:delete")
	}
	actions = append(actions, deleteAction)
	return actions
}

func buildBulkActions(l product_plan_staff.Labels, cl pyeza.CommonLabels, status string, routes product_plan_staff.Routes) []types.BulkAction {
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

func statusPageTitle(l product_plan_staff.Labels, status string) string {
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
