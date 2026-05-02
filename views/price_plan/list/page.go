package list

import (
	"context"
	"fmt"
	"log"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
)

type ListViewDeps struct {
	Routes               centymo.PricePlanRoutes
	ListPricePlans       func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ListPlans            func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules   func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	Labels               centymo.PricePlanLabels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels
	GetPricePlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var pricePlanSearchFields = []string{"name", "description"}

func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := pricePlanColumns(deps.Labels)
		p, err := espynahttp.ParseTableParams(viewCtx.Request, types.SortableKeys(columns), "name", "asc")
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
				HeaderSubtitle: statusPageCaption(deps.Labels, status),
				HeaderIcon:     "icon-tag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "price-plan-list-content",
			Table:           tableConfig,
		}

		return view.OK("price-plan-list", pageData)
	})
}

func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		columns := pricePlanColumns(deps.Labels)
		p, err := espynahttp.ParseTableParams(viewCtx.Request, types.SortableKeys(columns), "name", "asc")
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
	listParams := espynahttp.ToListParams(p, pricePlanSearchFields)

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

	resp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list price plans: %v", err)
		return nil, fmt.Errorf("failed to load price plans: %w", err)
	}

	var inUseIDs map[string]bool
	if deps.GetPricePlanInUseIDs != nil {
		var itemIDs []string
		for _, item := range resp.GetData() {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetPricePlanInUseIDs(ctx, itemIDs)
	}

	// Build plan name lookup map
	planNames := map[string]string{}
	if deps.ListPlans != nil {
		planResp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
		if err != nil {
			log.Printf("Failed to list plans for price plan table: %v", err)
		} else {
			for _, p := range planResp.GetData() {
				planNames[p.GetId()] = p.GetName()
			}
		}
	}

	// Build schedule name lookup map
	scheduleNames := map[string]string{}
	if deps.ListPriceSchedules != nil {
		schedResp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
		if err != nil {
			log.Printf("Failed to list price schedules for price plan table: %v", err)
		} else {
			for _, s := range schedResp.GetData() {
				scheduleNames[s.GetId()] = s.GetName()
			}
		}
	}

	l := deps.Labels
	rows := buildTableRows(resp.GetData(), status, l, deps.Routes, inUseIDs, perms, planNames, scheduleNames, deps.CommonLabels.DurationUnit)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:              "delete",
			Label:            l.Bulk.DeleteTitle,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         deps.Routes.BulkDeleteURL,
			ConfirmTitle:     l.Bulk.DeleteTitle,
			ConfirmMessage:   l.Bulk.DeleteMessage,
			RequiresDataAttr: "deletable",
		},
	}

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)
	tableConfig := &types.TableConfig{
		ID:                   "price-plans-table",
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
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.Add,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("price_plan", "create"),
			DisabledTooltip: l.Errors.Unauthorized,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig, nil
}

func pricePlanColumns(l centymo.PricePlanLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "amount", Label: l.Columns.Amount, WidthClass: "col-2xl"},
		{Key: "duration", Label: l.Columns.Duration, NoFilter: true, WidthClass: "col-2xl"},
		{Key: "plan", Label: l.Columns.Plan, NoSort: true, NoFilter: true},
		{Key: "schedule", Label: l.Columns.Schedule, NoSort: true, NoFilter: true},
		{Key: "status", Label: l.Columns.Status, NoFilter: true, WidthClass: "col-2xl"},
	}
}

func buildTableRows(pricePlans []*priceplanpb.PricePlan, status string, l centymo.PricePlanLabels, routes centymo.PricePlanRoutes, inUseIDs map[string]bool, perms *types.UserPermissions, planNames, scheduleNames map[string]string, durationLabels pyeza.DurationUnitLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, pp := range pricePlans {
		recordStatus := "active"
		if !pp.GetActive() {
			recordStatus = "inactive"
		}

		id := pp.GetId()
		name := pp.GetName()

		currency := pp.GetBillingCurrency()

		durationDisplay := ""
		if pp.GetDurationValue() > 0 {
			durationDisplay = pyeza.FormatDuration(pp.GetDurationValue(), pp.GetDurationUnit(), durationLabels)
		}

		planName := planNames[pp.GetPlanId()]

		scheduleName := "—"
		if schedID := pp.GetPriceScheduleId(); schedID != "" {
			if n, ok := scheduleNames[schedID]; ok && n != "" {
				scheduleName = n
			} else {
				scheduleName = schedID
			}
		}

		isInUse := inUseIDs[id]

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
		}
		if !perms.Can("price_plan", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = l.Errors.Unauthorized
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				types.MoneyCell(float64(pp.GetBillingAmount()), currency, true),
				{Type: "text", Value: durationDisplay},
				{Type: "text", Value: planName},
				{Type: "text", Value: scheduleName},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":      name,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(!isInUse),
			},
			Actions: []types.TableAction{
				{Type: "view", Label: l.Buttons.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
				{Type: "edit", Label: l.Buttons.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Buttons.Edit, Disabled: !perms.Can("price_plan", "update"), DisabledTooltip: l.Errors.Unauthorized},
				deleteAction,
			},
		})
	}
	return rows
}

func statusPageTitle(l centymo.PricePlanLabels, status string) string {
	switch status {
	case "active":
		return l.Page.ActiveTitle
	case "inactive":
		return l.Page.InactiveTitle
	default:
		return l.Page.Title
	}
}

func statusPageCaption(l centymo.PricePlanLabels, status string) string {
	switch status {
	case "active":
		return l.Page.Subtitle
	case "inactive":
		return l.Page.Subtitle
	default:
		return l.Page.Subtitle
	}
}

func statusEmptyTitle(l centymo.PricePlanLabels, status string) string {
	return l.Empty.Title
}

func statusEmptyMessage(l centymo.PricePlanLabels, status string) string {
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
