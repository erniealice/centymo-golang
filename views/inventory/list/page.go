package list

import (
	"context"
	"fmt"
	"log"
	"math"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/tableparams"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"

	"github.com/erniealice/centymo-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes             centymo.InventoryRoutes
	ListInventoryItems func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	Labels             centymo.InventoryLabels
	CommonLabels       pyeza.CommonLabels
	TableLabels        types.TableLabels
}

// PageData holds the data for the inventory list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

var inventorySearchFields = []string{"product_name", "sku"}

// NewView creates the inventory list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		location := viewCtx.Request.PathValue("location")
		if location == "" {
			location = "ayala-central-bloc"
		}

		columns := inventoryColumns(deps.Labels)
		p, err := espynahttp.ParseTableParams(viewCtx.Request, types.SortableKeys(columns), "date_created", "desc")
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, location, p)
		if err != nil {
			return view.Error(err)
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading + " \u2014 " + centymo.LocationDisplayName(location),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "inventory",
				ActiveSubNav:   location,
				HeaderTitle:    deps.Labels.Page.Heading + " \u2014 " + centymo.LocationDisplayName(location),
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "inventory-list-content",
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "inventory"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("inventory-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations so that only the table
// is swapped (not the entire page content).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		location := viewCtx.Request.PathValue("location")
		if location == "" {
			location = "ayala-central-bloc"
		}

		columns := inventoryColumns(deps.Labels)
		p, err := espynahttp.ParseTableParams(viewCtx.Request, types.SortableKeys(columns), "date_created", "desc")
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, location, p)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches inventory data and builds the table configuration.
func buildTableConfig(ctx context.Context, deps *ListViewDeps, columns []types.TableColumn, location string, p tableparams.TableQueryParams) (*types.TableConfig, error) {
	perms := view.GetUserPermissions(ctx)

	listParams := espynahttp.ToListParams(p, inventorySearchFields)
	resp, err := deps.ListInventoryItems(ctx, &inventoryitempb.ListInventoryItemsRequest{
		LocationId: &location,
		Search:     listParams.Search,
		Filters:    listParams.Filters,
		Sort:       listParams.Sort,
		Pagination: listParams.Pagination,
	})
	if err != nil {
		log.Printf("Failed to list inventory: %v", err)
		return nil, fmt.Errorf("failed to load inventory: %w", err)
	}

	l := deps.Labels
	rows := buildTableRows(resp.GetData(), l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:             "activate",
			Label:           l.Status.Activate,
			Icon:            "icon-check-circle",
			Variant:         "success",
			Endpoint:        deps.Routes.BulkSetStatusURL,
			ConfirmTitle:    l.Status.Activate,
			ConfirmMessage:  l.Confirm.BulkActivateMessage,
			ExtraParamsJSON: `{"target_status":"active"}`,
		},
		{
			Key:             "deactivate",
			Label:           l.Status.Deactivate,
			Icon:            "icon-x-circle",
			Variant:         "warning",
			Endpoint:        deps.Routes.BulkSetStatusURL,
			ConfirmTitle:    l.Status.Deactivate,
			ConfirmMessage:  l.Confirm.BulkDeactivateMessage,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		},
		{
			Key:            "delete",
			Label:          deps.CommonLabels.Bulk.Delete,
			Icon:           "icon-trash-2",
			Variant:        "danger",
			Endpoint:       deps.Routes.BulkDeleteURL,
			ConfirmTitle:   deps.CommonLabels.Bulk.Delete,
			ConfirmMessage: l.Confirm.BulkDeleteMessage,
		},
	}

	refreshURL := deps.Routes.TableURL

	// Build ServerPagination
	totalRows := len(rows) // TODO: migrate to GetInventoryItemListPageData (CTE variant) to get resp.GetPagination().GetTotalItems(); ListInventoryItemsResponse has no pagination field
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

	tableConfig := &types.TableConfig{
		ID:                   "inventory-table",
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
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddItem,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("inventory_item", "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
		},
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func inventoryColumns(l centymo.InventoryLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "product_name", Label: l.Columns.ProductName, Filterable: true, FilterType: types.FilterTypeString},
		{Key: "sku", Label: l.Columns.SKU, Filterable: false, WidthClass: "col-4xl"},
		{Key: "tracking_mode", Label: l.Columns.Type, Filterable: false, WidthClass: "col-3xl"},
		{Key: "quantity", Label: l.Columns.OnHand, Filterable: true, FilterType: types.FilterTypeNumeric, WidthClass: "col-2xl"},
		{Key: "available", Label: l.Columns.Available, NoSort: true, Filterable: false, WidthClass: "col-2xl"},
		{Key: "reorder_level", Label: l.Columns.ReorderLvl, NoSort: true, Filterable: false, WidthClass: "col-3xl"},
		{Key: "date_created", Label: "Date Created", Filterable: true, FilterType: types.FilterTypeDate},
		{Key: "status", Label: l.Columns.Status, Filterable: false, WidthClass: "col-2xl"},
	}
}

func buildTableRows(items []*inventoryitempb.InventoryItem, l centymo.InventoryLabels, routes centymo.InventoryRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, item := range items {
		id := item.GetId()
		name := item.GetName()
		sku := item.GetSku()
		onHand := item.GetQuantityOnHand()
		reserved := item.GetQuantityReserved()
		reorderLvl := item.GetReorderLevel()
		itemType := item.GetProduct().GetTrackingMode()
		if itemType == "" {
			itemType = "bulk"
		}

		avail := onHand - reserved
		if avail < 0 {
			avail = 0
		}
		available := formatFloat(avail)
		onHandStr := formatFloat(onHand)
		reservedStr := formatFloat(reserved)
		reorderStr := formatFloat(reorderLvl)

		status := "active"
		if !item.GetActive() {
			status = "inactive"
		}

		// Low stock alert: if available quantity is at or below reorder level
		reorderDisplay := reorderStr
		if reorderLvl > 0 && avail <= reorderLvl {
			reorderDisplay = reorderStr + " (!)"
		}

		dateCreated := item.GetDateCreatedString()
		detailURL := route.ResolveURL(routes.DetailURL, "id", id)

		rows = append(rows, types.TableRow{
			ID:   id,
			Href: detailURL,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: sku},
				{Type: "badge", Value: itemTypeLabel(itemType, l), Variant: itemTypeVariant(itemType)},
				{Type: "text", Value: onHandStr},
				{Type: "text", Value: available},
				{Type: "text", Value: reorderDisplay},
				types.DateTimeCell(dateCreated, types.DateReadable),
				{Type: "badge", Value: status, Variant: statusVariant(status)},
			},
			DataAttrs: map[string]string{
				"name":        name,
				"sku":         sku,
				"tracking_mode": itemType,
				"on_hand":     onHandStr,
				"reserved":    reservedStr,
				"available":   available,
				"reorder_lvl": reorderStr,
				"status":      status,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: l.Actions.View, Action: "view", Href: detailURL},
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("inventory_item", "update"), DisabledTooltip: l.Errors.PermissionDenied},
				{Type: "delete", Label: l.Actions.Delete, Action: "delete", URL: routes.DeleteURL, ItemName: name, Disabled: !perms.Can("inventory_item", "delete"), DisabledTooltip: l.Errors.PermissionDenied},
			},
		})
	}
	return rows
}

func itemTypeLabel(itemType string, l centymo.InventoryLabels) string {
	switch itemType {
	case "none":
		return l.TrackingMode.None
	case "bulk":
		return l.TrackingMode.Bulk
	case "serialized":
		return l.TrackingMode.Serialized
	default:
		return itemType
	}
}

func itemTypeVariant(itemType string) string {
	switch itemType {
	case "none":
		return "neutral"
	case "bulk":
		return "info"
	case "serialized":
		return "success"
	default:
		return "default"
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

func formatFloat(f float64) string {
	if f == float64(int64(f)) {
		return fmt.Sprintf("%d", int64(f))
	}
	return fmt.Sprintf("%.2f", f)
}
