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
	costplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_plan"
)

// ListViewDeps holds view dependencies for the cost_plan list.
type ListViewDeps struct {
	Routes                  centymo.CostPlanRoutes
	GetCostPlanListPageData func(ctx context.Context, req *costplanpb.GetCostPlanListPageDataRequest) (*costplanpb.GetCostPlanListPageDataResponse, error)
	Labels                  centymo.CostPlanLabels
	CommonLabels            pyeza.CommonLabels
	TableLabels             types.TableLabels
}

// PageData holds the data for the cost_plan list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// CostPlanSortSpec is the canonical sort spec.
var CostPlanSortSpec = espynahttp.SortSpec{
	AllowedCols: []string{"name", "billing_kind", "billing_amount", "billing_currency", "date_created"},
	DefaultCol:  "name",
	DefaultDir:  "asc",
}

var costPlanSearchFields = []string{"name"}

// NewView creates the cost_plan list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		p, err := espynahttp.ParseTableParamsFromSpec(viewCtx.Request, CostPlanSortSpec)
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
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "cost-plan-list-content",
			Table:           tableConfig,
		}
		return view.OK("cost-plan-list", pageData)
	})
}

// NewTableView returns only the table-card partial.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		p, err := espynahttp.ParseTableParamsFromSpec(viewCtx.Request, CostPlanSortSpec)
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
	listParams := espynahttp.ToListParams(p, costPlanSearchFields)

	activeValue := status != "inactive"
	if listParams.Filters == nil {
		listParams.Filters = &commonpb.FilterRequest{}
	}
	listParams.Filters.Filters = append(listParams.Filters.Filters, &commonpb.TypedFilter{
		Field: "cp.active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: activeValue},
		},
	})

	resp, err := deps.GetCostPlanListPageData(ctx, &costplanpb.GetCostPlanListPageDataRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list cost plans: %v", err)
		return nil, fmt.Errorf("failed to load cost plans: %w", err)
	}

	l := deps.Labels
	columns := costPlanColumns(l)
	rows := buildTableRows(resp.GetCostPlanList(), status, l, deps.Routes, perms)
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
		ID:                   "cost-plans-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    CostPlanSortSpec.DefaultCol,
		DefaultSortDirection: CostPlanSortSpec.DefaultDir,
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
			Label:           l.Buttons.AddCostPlan,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("cost_plan", "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
		}
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig, nil
}

func costPlanColumns(l centymo.CostPlanLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "billing_kind", Label: l.Columns.BillingKind},
		{Key: "billing_amount", Label: l.Columns.Amount},
		{Key: "billing_currency", Label: l.Columns.Currency, NoFilter: true, NoSort: true},
		{Key: "supplier_plan", Label: l.Columns.SupplierPlan, NoFilter: true, NoSort: true},
	}
}

func buildTableRows(plans []*costplanpb.CostPlan, status string, l centymo.CostPlanLabels, routes centymo.CostPlanRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, cp := range plans {
		active := cp.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}
		id := cp.GetId()
		name := cp.GetName()
		billingKind := cp.GetBillingKind().String()
		amount := formatCentavos(cp.GetBillingAmount())
		currency := cp.GetBillingCurrency()

		spLabel := cp.GetSupplierPlanId()
		if sp := cp.GetSupplierPlan(); sp != nil && sp.GetName() != "" {
			spLabel = sp.GetName()
		}

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
				Disabled:        !perms.Can("cost_plan", "update"),
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
				Disabled:        !perms.Can("cost_plan", "update"),
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
				Disabled:        !perms.Can("cost_plan", "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		}
		actions = append(actions, types.TableAction{
			Type:            "delete",
			Label:           l.Actions.Delete,
			Action:          "delete",
			URL:             routes.DeleteURL,
			ItemName:        name,
			Disabled:        !perms.Can("cost_plan", "delete"),
			DisabledTooltip: l.Errors.PermissionDenied,
		})

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: billingKind},
				{Type: "number", Value: amount},
				{Type: "text", Value: currency},
				{Type: "text", Value: spLabel},
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

func statusTitle(l centymo.CostPlanLabels, status string) string {
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

func formatCentavos(centavos int64) string {
	whole := centavos / 100
	frac := centavos % 100
	if frac < 0 {
		frac = -frac
	}
	return fmt.Sprintf("%d.%02d", whole, frac)
}
