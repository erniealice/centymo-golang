package list

import (
	"context"
	"fmt"
	"log"
	"math"

	centymo "github.com/erniealice/centymo-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	supplierplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_plan"
)

// ListViewDeps holds view dependencies for the supplier_plan list.
type ListViewDeps struct {
	Routes                     centymo.SupplierPlanRoutes
	GetSupplierPlanListPageData func(ctx context.Context, req *supplierplanpb.GetSupplierPlanListPageDataRequest) (*supplierplanpb.GetSupplierPlanListPageDataResponse, error)
	Labels                     centymo.SupplierPlanLabels
	CommonLabels               pyeza.CommonLabels
	TableLabels                types.TableLabels
}

// PageData holds the data for the supplier_plan list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// SupplierPlanSortSpec is the canonical sort spec.
var SupplierPlanSortSpec = espynahttp.SortSpec{
	AllowedCols: []string{"name", "date_created"},
	DefaultCol:  "name",
	DefaultDir:  "asc",
}

var supplierPlanSearchFields = []string{"name"}

// NewView creates the supplier_plan list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		p, err := espynahttp.ParseTableParamsFromSpec(viewCtx.Request, SupplierPlanSortSpec)
		if err != nil {
			return view.Error(err)
		}
		tableConfig, err := buildTableConfig(ctx, deps, status, p)
		if err != nil {
			return view.Error(err)
		}
		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(l, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav + "-" + status,
				HeaderTitle:    statusTitle(l, status),
				HeaderSubtitle: l.Page.Caption,
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "supplier-plan-list-content",
			Table:           tableConfig,
		}
		return view.OK("supplier-plan-list", pageData)
	})
}

// NewTableView returns only the table-card partial.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		p, err := espynahttp.ParseTableParamsFromSpec(viewCtx.Request, SupplierPlanSortSpec)
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

func buildTableConfig(ctx context.Context, deps *ListViewDeps, status string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)
	listParams := espynahttp.ToListParams(p, supplierPlanSearchFields)

	activeValue := status != "inactive"
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "sp.active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: activeValue},
		},
	})

	resp, err := deps.GetSupplierPlanListPageData(ctx, &supplierplanpb.GetSupplierPlanListPageDataRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list supplier plans: %v", err)
		return nil, fmt.Errorf("failed to load supplier plans: %w", err)
	}

	l := deps.Labels
	columns := supplierPlanColumns(l)
	rows := buildTableRows(ctx, resp.GetSupplierPlanList(), status, l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	totalRows := int(resp.GetPagination().GetTotalItems())
	sp := &types.ServerPagination{
		Enabled:       true,
		Mode:          "offset",
		CurrentPage:   p.Page,
		PageSize:      p.PageSize,
		TotalRows:     totalRows,
		TotalPages:    int(math.Ceil(float64(totalRows) / float64(p.PageSize))),
		SearchQuery:   p.Search,
		SortColumn:    p.SortColumn,
		SortDirection: p.SortDir,
		FiltersJSON:   p.FiltersRaw,
		PaginationURL: refreshURL,
	}
	sp.BuildDisplay()

	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:              "activate",
			Label:            l.Status.Activate,
			Icon:             "icon-check-circle",
			Variant:          "success",
			Endpoint:         deps.Routes.BulkSetStatusURL,
			ExtraParamsJSON:  `{"target_status":"active"}`,
			ConfirmTitle:     l.Confirm.BulkActivate,
			ConfirmMessage:   l.Confirm.BulkActivateMessage,
			RequiresDataAttr: "activatable",
		},
		{
			Key:              "deactivate",
			Label:            l.Status.Deactivate,
			Icon:             "icon-x-circle",
			Variant:          "warning",
			Endpoint:         deps.Routes.BulkSetStatusURL,
			ExtraParamsJSON:  `{"target_status":"inactive"}`,
			ConfirmTitle:     l.Confirm.BulkDeactivate,
			ConfirmMessage:   l.Confirm.BulkDeactivateMessage,
			RequiresDataAttr: "deactivatable",
		},
		{
			Key:              "delete",
			Label:            l.Bulk.Delete,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         deps.Routes.BulkDeleteURL,
			ConfirmTitle:     l.Confirm.BulkDelete,
			ConfirmMessage:   l.Confirm.BulkDeleteMessage,
			RequiresDataAttr: "deletable",
		},
	}

	tableConfig := &types.TableConfig{
		ID:                   "supplier-plans-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    SupplierPlanSortSpec.DefaultCol,
		DefaultSortDirection: SupplierPlanSortSpec.DefaultDir,
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		ServerPagination: sp,
		BulkActions:      &bulkCfg,
	}
	if status == "active" {
		tableConfig.PrimaryAction = &types.PrimaryAction{
			Label:           l.Buttons.AddSupplierPlan,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("supplier_plan", "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
		}
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig, nil
}

func supplierPlanColumns(l centymo.SupplierPlanLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "supplier", Label: l.Columns.Supplier, NoSort: true, NoFilter: true},
	}
}

func buildTableRows(ctx context.Context, plans []*supplierplanpb.SupplierPlan, status string, l centymo.SupplierPlanLabels, routes centymo.SupplierPlanRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, p := range plans {
		active := p.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		id := p.GetId()
		name := p.GetName()
		supplierName := p.GetSupplierId()

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
		}
		if recordStatus == "active" {
			actions = append(actions, types.TableAction{
				Type:            "edit",
				Label:           l.Actions.Edit,
				Action:          "edit",
				URL:             route.ResolveURL(routes.EditURL, "id", id),
				DrawerTitle:     l.Actions.Edit,
				Disabled:        !perms.Can("supplier_plan", "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
			actions = append(actions, types.TableAction{
				Type:            "deactivate",
				Label:           l.Actions.Deactivate,
				Action:          "deactivate",
				URL:             routes.SetStatusURL + "?status=inactive",
				ItemName:        name,
				ConfirmTitle:    l.Confirm.Deactivate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.DeactivateMessage, name),
				Disabled:        !perms.Can("supplier_plan", "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
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
				Disabled:        !perms.Can("supplier_plan", "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		}
		actions = append(actions, types.TableAction{
			Type:            "delete",
			Label:           l.Actions.Delete,
			Action:          "delete",
			URL:             routes.DeleteURL,
			ItemName:        name,
			Disabled:        !perms.Can("supplier_plan", "delete"),
			DisabledTooltip: l.Errors.PermissionDenied,
		})

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: supplierName},
			},
			DataAttrs: map[string]string{
				"name":          name,
				"status":        recordStatus,
				"deletable":     "true",
				"activatable":   boolAttr(recordStatus == "inactive"),
				"deactivatable": boolAttr(recordStatus == "active"),
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l centymo.SupplierPlanLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func boolAttr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
