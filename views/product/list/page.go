package list

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	linepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"

	centymo "github.com/erniealice/centymo-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes centymo.ProductRoutes
	// Mode selects the product_kind filter set for the list query.
	// "service" (default/empty) or "inventory" — selects product_kind filter set.
	// Zero-value ("") maps to the service behaviour (product_kind = 'service').
	Mode         string
	ListProducts func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ListLines    func(ctx context.Context, req *linepb.ListLinesRequest) (*linepb.ListLinesResponse, error)
	GetInUseIDs  func(ctx context.Context, ids []string) (map[string]bool, error)
	Labels       centymo.ProductLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	// PermissionEntity is the first argument to perms.Can(entity, action) on
	// the list view's primary/row/bulk action buttons. Defaults to "product".
	// See centymo-golang/views/product/module.go ModuleDeps.PermissionEntity.
	PermissionEntity string
}

// permEntity returns the configured PermissionEntity with a safe default so
// the disabled-state logic never nil-guards.
func (d *ListViewDeps) permEntity() string {
	if d == nil || d.PermissionEntity == "" {
		return "product"
	}
	return d.PermissionEntity
}

// modeProductKinds maps a ListViewDeps.Mode value to the set of
// product_kind strings that should be included in the list filter.
// Zero-value ("") maps to the service behaviour so that callers which
// do not set Mode continue to see service products.
var modeProductKinds = map[string][]string{
	"":          {"service"},                         // zero-value = service behaviour
	"service":   {"service"},                         // explicit service mount
	"inventory": {"stocked_good", "non_stocked_good"}, // resold goods only
	"supplies":  {"consumable"},                      // consumables used in service delivery
}

// PageData holds the data for the product list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var productSearchFields = []string{"name", "description"}

// NewView creates the product list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := productColumns(deps.Labels)
		p, err := espynahttp.ParseTableParams(viewCtx.Request, types.SortableKeys(columns), "name", "asc")
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, status, p)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    statusPageTitle(deps.Labels, status),
				HeaderSubtitle: statusPageCaption(deps.Labels, status),
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "product-list-content",
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "product"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("product-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations so that only the table
// is swapped (not the entire page content).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		columns := productColumns(deps.Labels)
		p, err := espynahttp.ParseTableParams(viewCtx.Request, types.SortableKeys(columns), "name", "asc")
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, status, p)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches product data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, columns []types.TableColumn, status string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, productSearchFields)

	// Inject active BooleanFilter so the repository returns records matching
	// the requested status. dbOps.List defaults to active=true; by supplying an
	// explicit filter here we override that default for the inactive list.
	activeValue := status != "inactive"
	activeFilter := &commonpb.TypedFilter{
		Field: "active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: activeValue},
		},
	}

	// Inject product_kind IN (...) filter so the list shows only products
	// that belong to this mount's surface. The filter values are selected by
	// deps.Mode — "service" → service, "inventory" → stocked_good/non_stocked_good/consumable.
	// Zero-value / unknown Mode falls back to the service behaviour.
	values, ok := modeProductKinds[deps.Mode]
	if !ok || len(values) == 0 {
		values = []string{"service"}
	}
	productKindFilter := &commonpb.TypedFilter{
		Field: "product_kind",
		FilterType: &commonpb.TypedFilter_ListFilter{
			ListFilter: &commonpb.ListFilter{
				Values:   values,
				Operator: commonpb.ListOperator_LIST_IN,
			},
		},
	}
	filters := listParams.Filters
	if filters == nil {
		filters = &commonpb.FilterRequest{}
	}
	filters.Filters = append(filters.Filters, activeFilter, productKindFilter)

	resp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{
		Search:     listParams.Search,
		Filters:    filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list products: %v", err)
		return nil, fmt.Errorf("failed to load products: %w", err)
	}

	var inUseIDs map[string]bool
	if deps.GetInUseIDs != nil {
		var itemIDs []string
		for _, item := range resp.GetData() {
			itemIDs = append(itemIDs, item.GetId())
		}
		inUseIDs, _ = deps.GetInUseIDs(ctx, itemIDs)
	}

	// Build line name lookup map for the line column.
	lineNameByID := map[string]string{}
	if deps.ListLines != nil {
		lineResp, lerr := deps.ListLines(ctx, &linepb.ListLinesRequest{})
		if lerr != nil {
			log.Printf("Failed to list lines for product table: %v", lerr)
		} else {
			for _, line := range lineResp.GetData() {
				if line != nil {
					lineNameByID[line.GetId()] = line.GetName()
				}
			}
		}
	}

	l := deps.Labels
	rows := buildTableRows(resp.GetData(), status, l, deps.CommonLabels, deps.Routes, inUseIDs, perms, deps.permEntity(), lineNameByID)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:              "activate",
			Label:            l.Status.Activate,
			Icon:             "icon-check-circle",
			Variant:          "success",
			Endpoint:         deps.Routes.BulkSetStatusURL,
			ExtraParamsJSON:  `{"target_status":"active"}`,
			ConfirmTitle:     l.Status.Activate,
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
			ConfirmTitle:     l.Status.Deactivate,
			ConfirmMessage:   l.Confirm.BulkDeactivateMessage,
			RequiresDataAttr: "deactivatable",
		},
		{
			Key:              "delete",
			Label:            l.Bulk.Delete,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         deps.Routes.BulkDeleteURL,
			ConfirmTitle:     l.Bulk.Delete,
			ConfirmMessage:   l.Confirm.BulkDeleteMessage,
			RequiresDataAttr: "deletable",
		},
	}

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	// Build ServerPagination
	totalRows := len(rows)
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

	var primaryAction *types.PrimaryAction
	if status == "active" {
		primaryAction = &types.PrimaryAction{
			Label:           l.Buttons.AddProduct,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can(deps.permEntity(), "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
		}
	}

	tableConfig := &types.TableConfig{
		ID:                   "products-table",
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
		PrimaryAction:    primaryAction,
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func productColumns(l centymo.ProductLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Filterable: true, FilterType: types.FilterTypeString},
		{Key: "description", Label: l.Columns.Description, NoSort: true},
		{Key: "line", Label: l.Columns.Line, NoSort: true},
		{Key: "price", Label: l.Columns.Price, WidthClass: "col-4xl"},
		{Key: "date_created", Label: "Date Created", Filterable: true, FilterType: types.FilterTypeDate},
	}
}

func buildTableRows(products []*productpb.Product, status string, l centymo.ProductLabels, cl pyeza.CommonLabels, routes centymo.ProductRoutes, inUseIDs map[string]bool, perms *types.UserPermissions, permEntity string, lineNameByID map[string]string) []types.TableRow {
	rows := []types.TableRow{}
	for _, p := range products {
		active := p.GetActive()
		recordStatus := "active"
		if !active {
			recordStatus = "inactive"
		}

		id := p.GetId()
		name := p.GetName()
		description := p.GetDescription()
		currency := p.GetCurrency()
		isInUse := inUseIDs[id]
		lineName := lineNameByID[p.GetLineId()]

		deleteAction := types.TableAction{
			Type:     "delete",
			Label:    l.Actions.Delete,
			Action:   "delete",
			URL:      routes.DeleteURL,
			ItemName: name,
		}
		if isInUse {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = l.Errors.CannotDelete
		}
		if !perms.Can(permEntity, "delete") {
			deleteAction.Disabled = true
			deleteAction.DisabledTooltip = l.Errors.PermissionDenied
		}

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can(permEntity, "update"), DisabledTooltip: l.Errors.PermissionDenied},
		}
		if recordStatus == "active" {
			actions = append(actions, types.TableAction{
				Type:            "clone",
				Label:           cl.Actions.Clone,
				Action:          "clone",
				URL:             route.ResolveURL(routes.EditURL, "id", id),
				DrawerTitle:     cl.Actions.Clone,
				Disabled:        !perms.Can(permEntity, "create"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
			actions = append(actions, types.TableAction{
				Type:            "deactivate",
				Label:           l.Status.Deactivate,
				Action:          "deactivate",
				URL:             routes.SetStatusURL + "?status=inactive",
				ItemName:        name,
				ConfirmTitle:    l.Status.Deactivate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.DeactivateMessage, name),
				Disabled:        !perms.Can(permEntity, "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type:            "activate",
				Label:           l.Status.Activate,
				Action:          "activate",
				URL:             routes.SetStatusURL + "?status=active",
				ItemName:        name,
				ConfirmTitle:    l.Status.Activate,
				ConfirmMessage:  fmt.Sprintf(l.Confirm.ActivateMessage, name),
				Disabled:        !perms.Can(permEntity, "update"),
				DisabledTooltip: l.Errors.PermissionDenied,
			})
		}
		actions = append(actions, deleteAction)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: description},
				{Type: "text", Value: lineName},
				types.MoneyCell(float64(p.GetPrice()), currency, true),
				types.DateTimeCell(p.GetDateCreatedString(), types.DateReadable),
			},
			DataAttrs: map[string]string{
				"name":          name,
				"price":         fmt.Sprintf("%d", p.GetPrice()),
				"status":        recordStatus,
				"deletable":     strconv.FormatBool(!isInUse),
				"activatable":   strconv.FormatBool(recordStatus == "inactive"),
				"deactivatable": strconv.FormatBool(recordStatus == "active"),
			},
			Actions: actions,
		})
	}
	return rows
}

func statusPageTitle(l centymo.ProductLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusPageCaption(l centymo.ProductLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l centymo.ProductLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l centymo.ProductLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveMessage
	case "inactive":
		return l.Empty.InactiveMessage
	default:
		return l.Empty.ActiveMessage
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
