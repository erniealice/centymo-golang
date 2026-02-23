package list

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"

	"github.com/erniealice/centymo-golang"
)

// Deps holds view dependencies.
type Deps struct {
	ListInventoryItems func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	RefreshURL         string
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

// NewView creates the inventory list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		location := viewCtx.Request.PathValue("location")
		if location == "" {
			location = "ayala-central-bloc"
		}

		resp, err := deps.ListInventoryItems(ctx, &inventoryitempb.ListInventoryItemsRequest{
			LocationId: &location,
		})
		if err != nil {
			log.Printf("Failed to list inventory: %v", err)
			return view.Error(fmt.Errorf("failed to load inventory: %w", err))
		}

		l := deps.Labels
		columns := inventoryColumns(l)
		rows := buildTableRows(resp.GetData(), l)
		types.ApplyColumnStyles(columns, rows)

		bulkCfg := centymo.MapBulkConfig(deps.CommonLabels)
		bulkCfg.Actions = []types.BulkAction{
			{
				Key:             "activate",
				Label:           l.Status.Activate,
				Icon:            "icon-check-circle",
				Variant:         "success",
				Endpoint:        "/action/inventory/bulk-set-status",
				ConfirmTitle:    l.Status.Activate,
				ConfirmMessage:  "Are you sure you want to activate {{count}} item(s)?",
				ExtraParamsJSON: `{"target_status":"active"}`,
			},
			{
				Key:             "deactivate",
				Label:           l.Status.Deactivate,
				Icon:            "icon-x-circle",
				Variant:         "warning",
				Endpoint:        "/action/inventory/bulk-set-status",
				ConfirmTitle:    l.Status.Deactivate,
				ConfirmMessage:  "Are you sure you want to deactivate {{count}} item(s)?",
				ExtraParamsJSON: `{"target_status":"inactive"}`,
			},
			{
				Key:            "delete",
				Label:          deps.CommonLabels.Bulk.Delete,
				Icon:           "icon-trash-2",
				Variant:        "danger",
				Endpoint:       "/action/inventory/bulk-delete",
				ConfirmTitle:   deps.CommonLabels.Bulk.Delete,
				ConfirmMessage: "Are you sure you want to delete {{count}} item(s)? This action cannot be undone.",
			},
		}

		tableConfig := &types.TableConfig{
			ID:                   "inventory-table",
			RefreshURL:           deps.RefreshURL,
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
				Title:   l.Empty.Title,
				Message: l.Empty.Message,
			},
			PrimaryAction: &types.PrimaryAction{
				Label:     l.Buttons.AddItem,
				ActionURL: "/action/inventory/add",
				Icon:      "icon-plus",
			},
			BulkActions: &bulkCfg,
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          "Inventory \u2014 " + centymo.LocationDisplayName(location),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "inventory",
				ActiveSubNav:   location,
				HeaderTitle:    "Inventory \u2014 " + centymo.LocationDisplayName(location),
				HeaderSubtitle: l.Page.Caption,
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "inventory-list-content",
			Table:           tableConfig,
		}

		return view.OK("inventory-list", pageData)
	})
}

func inventoryColumns(l centymo.InventoryLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.ProductName, Sortable: true},
		{Key: "sku", Label: l.Columns.SKU, Sortable: true, Width: "150px"},
		{Key: "item_type", Label: "Type", Sortable: true, Width: "130px"},
		{Key: "on_hand", Label: l.Columns.OnHand, Sortable: true, Width: "120px"},
		{Key: "available", Label: l.Columns.Available, Sortable: true, Width: "120px"},
		{Key: "reorder_level", Label: l.Columns.ReorderLvl, Sortable: true, Width: "140px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(items []*inventoryitempb.InventoryItem, l centymo.InventoryLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, item := range items {
		id := item.GetId()
		name := item.GetName()
		sku := item.GetSku()
		onHand := item.GetQuantityOnHand()
		reserved := item.GetQuantityReserved()
		reorderLvl := item.GetReorderLevel()
		itemType := item.GetProduct().GetItemType()
		if itemType == "" {
			itemType = "non_serialized"
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

		detailURL := "/app/inventory/detail/" + id

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
				{Type: "badge", Value: status, Variant: statusVariant(status)},
			},
			DataAttrs: map[string]string{
				"name":        name,
				"sku":         sku,
				"item_type":   itemType,
				"on_hand":     onHandStr,
				"reserved":    reservedStr,
				"available":   available,
				"reorder_lvl": reorderStr,
				"status":      status,
			},
			Actions: []types.TableAction{
				{Type: "view", Label: l.Actions.View, Action: "view", Href: detailURL},
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: "/action/inventory/edit/" + id, DrawerTitle: l.Actions.Edit},
				{Type: "delete", Label: l.Actions.Delete, Action: "delete", URL: "/action/inventory/delete", ItemName: name},
			},
		})
	}
	return rows
}

func itemTypeLabel(itemType string, l centymo.InventoryLabels) string {
	switch itemType {
	case "serialized":
		return l.ItemType.Serialized
	case "non_serialized":
		return l.ItemType.NonSerialized
	case "consumable":
		return l.ItemType.Consumable
	default:
		return itemType
	}
}

func itemTypeVariant(itemType string) string {
	switch itemType {
	case "serialized":
		return "info"
	case "non_serialized":
		return "default"
	case "consumable":
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
