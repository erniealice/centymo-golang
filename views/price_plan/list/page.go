package list

import (
	"context"
	"fmt"
	"log"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

type ListViewDeps struct {
	Routes         centymo.PricePlanRoutes
	ListPricePlans func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ListLocations  func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	ListPlans      func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	Labels         centymo.PricePlanLabels
	CommonLabels   pyeza.CommonLabels
	TableLabels    types.TableLabels
}

type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var pricePlanAllowedSortCols = []string{"date_created", "date_modified", "name", "status"}
var pricePlanSearchFields = []string{"name", "description"}

func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}
		p, err := espynahttp.ParseTableParams(viewCtx.Request, pricePlanAllowedSortCols)
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
		p, err := espynahttp.ParseTableParams(viewCtx.Request, pricePlanAllowedSortCols)
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

func buildTableConfig(ctx context.Context, deps *ListViewDeps, status string, p espynahttp.TableQueryParams) (*types.TableConfig, error) {
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

	// Build location name lookup map
	locationNames := map[string]string{}
	if deps.ListLocations != nil {
		locResp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
		if err != nil {
			log.Printf("Failed to list locations for price plan table: %v", err)
		} else {
			for _, loc := range locResp.GetData() {
				locationNames[loc.GetId()] = loc.GetName()
			}
		}
	}

	l := deps.Labels
	columns := pricePlanColumns(l)
	rows := buildTableRows(resp.GetData(), status, l, deps.Routes, perms, planNames, locationNames)
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
		{Key: "name", Label: l.Columns.Name, Sortable: true, Filterable: true, FilterType: types.FilterTypeString},
		{Key: "amount", Label: l.Columns.Amount, Sortable: true, WidthClass: "col-2xl"},
		{Key: "duration", Label: l.Columns.Duration, Sortable: false, WidthClass: "col-2xl"},
		{Key: "plan", Label: l.Columns.Plan, Sortable: false},
		{Key: "location", Label: l.Columns.Location, Sortable: false},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Filterable: false, WidthClass: "col-2xl"},
	}
}

func buildTableRows(pricePlans []*priceplanpb.PricePlan, status string, l centymo.PricePlanLabels, routes centymo.PricePlanRoutes, perms *types.UserPermissions, planNames, locationNames map[string]string) []types.TableRow {
	rows := []types.TableRow{}
	for _, pp := range pricePlans {
		recordStatus := "active"
		if !pp.GetActive() {
			recordStatus = "inactive"
		}

		id := pp.GetId()
		name := pp.GetName()

		amountDisplay := strconv.FormatFloat(float64(pp.GetAmount())/100.0, 'f', 2, 64)
		currency := pp.GetCurrency()

		durationDisplay := ""
		if pp.GetDurationValue() > 0 {
			durationDisplay = strconv.FormatInt(int64(pp.GetDurationValue()), 10) + " " + pp.GetDurationUnit()
		}

		planName := planNames[pp.GetPlanId()]

		locationName := "—"
		if locID := pp.GetLocationId(); locID != "" {
			if n, ok := locationNames[locID]; ok && n != "" {
				locationName = n
			} else {
				locationName = locID
			}
		}

		deleteAction := types.TableAction{
			Type:     "delete",
			Label:    l.Buttons.Delete,
			Action:   "delete",
			URL:      routes.DeleteURL,
			ItemName: name,
		}
		if !perms.Can("price_plan", "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = l.Errors.Unauthorized
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "money", Value: amountDisplay, Currency: currency},
				{Type: "text", Value: durationDisplay},
				{Type: "text", Value: planName},
				{Type: "text", Value: locationName},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":      name,
				"status":    recordStatus,
				"deletable": strconv.FormatBool(true),
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
