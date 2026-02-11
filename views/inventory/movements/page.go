package movements

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	DB           centymo.DataSource
	Labels       centymo.InventoryLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PageData holds the data for the movements page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the global inventory movements page.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels

		transactions, err := deps.DB.ListSimple(ctx, "inventory_transaction")
		if err != nil {
			log.Printf("Failed to list inventory_transaction: %v", err)
			transactions = []map[string]any{}
		}

		// Build item name lookup
		items, err := deps.DB.ListSimple(ctx, "inventory_item")
		if err != nil {
			log.Printf("Failed to list inventory_item for movements: %v", err)
			items = []map[string]any{}
		}
		itemMap := map[string]map[string]any{}
		for _, item := range items {
			id, _ := item["id"].(string)
			itemMap[id] = item
		}

		columns := []types.TableColumn{
			{Key: "transaction_date", Label: l.Detail.Date, Sortable: true, Width: "130px"},
			{Key: "item_name", Label: l.Columns.ProductName, Sortable: true},
			{Key: "location", Label: l.Detail.Location, Sortable: true, Width: "160px"},
			{Key: "transaction_type", Label: l.Detail.Type, Sortable: true, Width: "120px"},
			{Key: "quantity", Label: l.Detail.Quantity, Sortable: true, Width: "100px"},
			{Key: "serial_number", Label: l.Detail.Serial, Sortable: false, Width: "150px"},
			{Key: "reference", Label: l.Detail.Reference, Sortable: false},
			{Key: "performed_by", Label: l.Detail.PerformedBy, Sortable: false, Width: "150px"},
		}

		rows := []types.TableRow{}
		for _, t := range transactions {
			id, _ := t["id"].(string)
			txDate, _ := t["transaction_date"].(string)
			txType, _ := t["transaction_type"].(string)
			qty := formatQuantity(t["quantity"], txType)
			ref, _ := t["reference"].(string)
			serial, _ := t["serial_number"].(string)
			performer, _ := t["performed_by"].(string)

			inventoryItemID, _ := t["inventory_item_id"].(string)
			itemName := inventoryItemID
			locationName := ""
			if item, ok := itemMap[inventoryItemID]; ok {
				name, _ := item["name"].(string)
				if name != "" {
					itemName = name
				}
				locID, _ := item["location_id"].(string)
				locationName = centymo.LocationDisplayName(locID)
			}

			rows = append(rows, types.TableRow{
				ID: id,
				Cells: []types.TableCell{
					{Type: "text", Value: txDate},
					{Type: "text", Value: itemName},
					{Type: "text", Value: locationName},
					{Type: "badge", Value: txType, Variant: txTypeVariant(txType)},
					{Type: "text", Value: qty},
					{Type: "text", Value: serial},
					{Type: "text", Value: ref},
					{Type: "text", Value: performer},
				},
			})
		}

		types.ApplyColumnStyles(columns, rows)

		tableConfig := &types.TableConfig{
			ID:                   "movements-table",
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
			DefaultSortColumn:    "transaction_date",
			DefaultSortDirection: "desc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   l.Detail.TransactionEmptyTitle,
				Message: l.Detail.TransactionEmptyMessage,
			},
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          "Inventory Movements",
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "inventory",
				ActiveSubNav:   "movements",
				HeaderTitle:    "Inventory Movements",
				HeaderSubtitle: l.Movements.Subtitle,
				HeaderIcon:     "icon-repeat",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "inventory-movements-content",
			Table:           tableConfig,
		}

		return view.OK("inventory-movements", pageData)
	})
}

func txTypeVariant(txType string) string {
	switch txType {
	case "received":
		return "success"
	case "sold":
		return "default"
	case "adjusted":
		return "info"
	case "transferred":
		return "warning"
	case "returned":
		return "danger"
	case "write_off":
		return "danger"
	default:
		return "default"
	}
}

func formatQuantity(qty any, txType string) string {
	val := fmt.Sprintf("%v", qty)
	switch txType {
	case "received", "returned":
		return "+" + val
	case "sold", "transferred", "write_off":
		return "-" + val
	default:
		return val
	}
}
