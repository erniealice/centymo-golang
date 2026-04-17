package list

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	resourcepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/resource"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes        centymo.ResourceRoutes
	ListResources func(ctx context.Context, req *resourcepb.ListResourcesRequest) (*resourcepb.ListResourcesResponse, error)
	Labels        centymo.ResourceLabels
	CommonLabels  pyeza.CommonLabels
	TableLabels   types.TableLabels
}

// PageData holds the data for the resource list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var resourceAllowedSortCols = []string{"date_created", "name", "status"}
var resourceSearchFields = []string{"name", "description"}

// NewView creates the resource list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		p, err := espynahttp.ParseTableParams(viewCtx.Request, resourceAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, status, p)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav + "-" + status,
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-box",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "resource-list-content",
			Table:           tableConfig,
		}

		return view.OK("resource-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		p, err := espynahttp.ParseTableParams(viewCtx.Request, resourceAllowedSortCols)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, status, p)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches resource data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, status string, p espynahttp.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, resourceSearchFields)

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

	resp, err := deps.ListResources(ctx, &resourcepb.ListResourcesRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list resources: %v", err)
		return nil, fmt.Errorf("failed to load resources: %w", err)
	}

	l := deps.Labels
	columns := resourceColumns(l)
	rows := buildTableRows(resp.GetData(), status, l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:             "activate",
			Label:           l.Status.Activate,
			Icon:            "icon-check-circle",
			Variant:         "success",
			Endpoint:        deps.Routes.BulkSetStatusURL,
			ExtraParamsJSON: `{"target_status":"active"}`,
			ConfirmTitle:    l.Confirm.BulkActivate,
			ConfirmMessage:  l.Confirm.BulkActivateMessage,
		},
		{
			Key:            "delete",
			Label:          l.Bulk.Delete,
			Icon:           "icon-trash-2",
			Variant:        "danger",
			Endpoint:       deps.Routes.BulkDeleteURL,
			ConfirmTitle:   l.Confirm.BulkDelete,
			ConfirmMessage: l.Confirm.BulkDeleteMessage,
		},
	}

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	tableConfig := &types.TableConfig{
		ID:                   "resources-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.Add,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("resource", "create"),
			DisabledTooltip: l.Errors.NoPermission,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func resourceColumns(l centymo.ResourceLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "description", Label: l.Columns.Description, Sortable: false},
		{Key: "product", Label: l.Columns.Product, Sortable: false},
		{Key: "status", Label: l.Columns.Status, Sortable: true, WidthClass: "col-2xl"},
	}
}

func buildTableRows(resources []*resourcepb.Resource, status string, l centymo.ResourceLabels, routes centymo.ResourceRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, res := range resources {
		active := res.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := res.GetId()
		name := res.GetName()
		description := res.GetDescription()
		productId := res.GetProductId()

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("resource", "update"), DisabledTooltip: l.Errors.NoPermission},
		}

		if recordStatus == "active" {
			actions = append(actions, types.TableAction{
				Type:            "deactivate",
				Label:           l.Actions.Deactivate,
				Action:          "deactivate",
				URL:             routes.SetStatusURL + "?status=inactive",
				ItemName:        name,
				ConfirmTitle:    l.Confirm.Deactivate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.DeactivateMessage, name),
				Disabled:        !perms.Can("resource", "update"),
				DisabledTooltip: l.Errors.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type:            "activate",
				Label:           l.Actions.Activate,
				Action:          "activate",
				URL:             routes.SetStatusURL + "?status=active",
				ItemName:        name,
				ConfirmTitle:    l.Confirm.Activate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.ActivateMessage, name),
				Disabled:        !perms.Can("resource", "update"),
				DisabledTooltip: l.Errors.NoPermission,
			})
		}

		deleteAction := types.TableAction{
			Type:     "delete",
			Label:    l.Actions.Delete,
			Action:   "delete",
			URL:      routes.DeleteURL,
			ItemName: name,
		}
		if !perms.Can("resource", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = l.Errors.NoPermission
		}
		actions = append(actions, deleteAction)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: description},
				{Type: "text", Value: productId},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":   name,
				"status": recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l centymo.ResourceLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l centymo.ResourceLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
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
