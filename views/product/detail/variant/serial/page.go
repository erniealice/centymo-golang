package serial

import (
	"context"
	"fmt"
	"log"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	detail "github.com/erniealice/centymo-golang/views/product/detail"
	"github.com/erniealice/centymo-golang/views/product/detail/variant"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
)

// SerialSummary holds serial count totals.
type SerialSummary struct {
	Total     int
	Available int
	Sold      int
	Reserved  int
}

// StockDetailPageData holds data for the variant-scoped stock detail page.
type StockDetailPageData struct {
	types.PageData
	ContentTemplate  string
	Breadcrumbs      []detail.Breadcrumb
	ProductID        string
	VariantID        string
	InventoryItemID  string
	ActiveTab        string
	TabItems         []pyeza.TabItem
	// Item info fields
	ItemName         string
	ItemSKU          string
	ItemType         string
	ItemTypeLabel    string
	ItemTypeVariant  string
	LocationName     string
	QuantityOnHand   string
	QuantityReserved string
	AvailableQty     string
	ItemStatus       string
	ItemStatusVariant string
	// Serial data
	SerialTable   *types.TableConfig
	SerialSummary *SerialSummary
	Labels        centymo.ProductLabels
}

// NewPageView creates the variant-scoped stock detail view (full page).
// Route: /app/products/detail/{id}/variant/{vid}/stock/{iid}
func NewPageView(deps *variant.Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")
		variantID := viewCtx.Request.PathValue("vid")
		itemID := viewCtx.Request.PathValue("iid")

		// Load product for breadcrumb context
		prodResp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
			Data: &productpb.Product{Id: productID},
		})
		if err != nil || len(prodResp.GetData()) == 0 {
			log.Printf("Failed to read product %s: %v", productID, err)
			return view.Error(fmt.Errorf("failed to load product: %w", err))
		}
		product := prodResp.GetData()[0]
		productName := product.GetName()

		// Load variant for breadcrumb context
		varResp, err := deps.ReadProductVariant(ctx, &productvariantpb.ReadProductVariantRequest{
			Data: &productvariantpb.ProductVariant{Id: variantID},
		})
		if err != nil || len(varResp.GetData()) == 0 {
			log.Printf("Failed to read product_variant %s: %v", variantID, err)
			return view.Error(fmt.Errorf("failed to load variant: %w", err))
		}
		variantData := varResp.GetData()[0]
		variantSKU := variantData.GetSku()

		// Load inventory item
		item, err := readInventoryItem(ctx, deps, itemID)
		if err != nil {
			return view.Error(err)
		}

		name := item.GetName()
		locationID := item.GetLocationId()
		locationName := centymo.LocationDisplayName(locationID)
		headerTitle := name + " \u2014 " + locationName

		itemType := item.GetItemType()
		if itemType == "" {
			itemType = "non_serialized"
		}

		active := item.GetActive()
		itemStatus := "active"
		if !active {
			itemStatus = "inactive"
		}

		available := computeAvailable(item.GetQuantityOnHand(), item.GetQuantityReserved())

		l := deps.Labels

		breadcrumbs := []detail.Breadcrumb{
			{Label: "Products", Href: "/app/products/list/active"},
			{Label: productName, Href: fmt.Sprintf("/app/products/detail/%s?tab=variants", productID)},
			{Label: variantSKU, Href: fmt.Sprintf("/app/products/detail/%s/variant/%s?tab=stock", productID, variantID)},
			{Label: name + " @ " + locationName, Href: ""},
		}

		// Build variant-context tabs (same as variant page, stock tab active)
		tabItems := buildVariantTabItems(productID, variantID, l)

		// Load serials for this inventory item
		serials := loadSerials(ctx, deps, itemID)
		serialTable := buildSerialTable(serials, deps.TableLabels, itemID)
		serialSummary := computeSerialSummary(serials)

		pageData := &StockDetailPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "products",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: item.GetSku(),
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:   "variant-stock-detail-content",
			Breadcrumbs:       breadcrumbs,
			ProductID:         productID,
			VariantID:         variantID,
			InventoryItemID:   itemID,
			ActiveTab:         "stock",
			TabItems:          tabItems,
			ItemName:          name,
			ItemSKU:           item.GetSku(),
			ItemType:          itemType,
			ItemTypeLabel:     itemTypeDisplayLabel(itemType),
			ItemTypeVariant:   itemTypeDisplayVariant(itemType),
			LocationName:      locationName,
			QuantityOnHand:    fmt.Sprintf("%v", item.GetQuantityOnHand()),
			QuantityReserved:  fmt.Sprintf("%v", item.GetQuantityReserved()),
			AvailableQty:      available,
			ItemStatus:        itemStatus,
			ItemStatusVariant: detail.StatusVariant(itemStatus),
			SerialTable:       serialTable,
			SerialSummary:     serialSummary,
			Labels:            l,
		}

		return view.OK("variant-stock-detail", pageData)
	})
}

// readInventoryItem reads a single inventory item by ID.
// Uses ReadInventoryItem if available, otherwise falls back to ListInventoryItems with filtering.
func readInventoryItem(ctx context.Context, deps *variant.Deps, itemID string) (*inventoryitempb.InventoryItem, error) {
	if deps.ReadInventoryItem != nil {
		resp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
			Data: &inventoryitempb.InventoryItem{Id: itemID},
		})
		if err != nil {
			log.Printf("Failed to read inventory_item %s: %v", itemID, err)
			return nil, fmt.Errorf("failed to load inventory item: %w", err)
		}
		items := resp.GetData()
		if len(items) == 0 {
			return nil, fmt.Errorf("inventory item not found")
		}
		return items[0], nil
	}

	// Fallback: list all and filter
	if deps.ListInventoryItems == nil {
		return nil, fmt.Errorf("no inventory item reader available")
	}
	resp, err := deps.ListInventoryItems(ctx, &inventoryitempb.ListInventoryItemsRequest{})
	if err != nil {
		log.Printf("Failed to list inventory_items: %v", err)
		return nil, fmt.Errorf("failed to load inventory items: %w", err)
	}
	for _, item := range resp.GetData() {
		if item.GetId() == itemID {
			return item, nil
		}
	}
	return nil, fmt.Errorf("inventory item not found")
}

// loadSerials loads serial numbers for an inventory item.
func loadSerials(ctx context.Context, deps *variant.Deps, inventoryItemID string) []*inventoryserialpb.InventorySerial {
	if deps.ListInventorySerials == nil {
		return nil
	}
	resp, err := deps.ListInventorySerials(ctx, &inventoryserialpb.ListInventorySerialsRequest{
		InventoryItemId: &inventoryItemID,
	})
	if err != nil {
		log.Printf("Failed to list inventory_serial: %v", err)
		return nil
	}
	return resp.GetData()
}

// buildSerialTable builds the serial numbers table.
func buildSerialTable(serials []*inventoryserialpb.InventorySerial, tableLabels types.TableLabels, inventoryItemID string) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "serial_number", Label: "Serial Number", Sortable: true},
		{Key: "imei", Label: "IMEI", Sortable: false, Width: "180px"},
		{Key: "status", Label: "Status", Sortable: true, Width: "120px"},
		{Key: "warranty_end", Label: "Warranty End", Sortable: true, Width: "140px"},
		{Key: "purchase_order", Label: "Purchase Order", Sortable: false, Width: "140px"},
	}

	rows := []types.TableRow{}
	for _, s := range serials {
		id := s.GetId()
		serial := s.GetSerialNumber()
		imei := s.GetImei()
		status := s.GetStatus()
		warrantyEnd := s.GetWarrantyEnd()
		po := s.GetPurchaseOrder()

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: serial},
				{Type: "text", Value: imei},
				{Type: "badge", Value: status, Variant: serialStatusVariant(status)},
				{Type: "text", Value: warrantyEnd},
				{Type: "text", Value: po},
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   "variant-serial-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowEntries:          true,
		ShowSort:             true,
		ShowDensity:          true,
		DefaultSortColumn:    "serial_number",
		DefaultSortDirection: "asc",
		Labels:               tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   "No Serial Numbers",
			Message: "No serial numbers have been recorded for this inventory item.",
		},
	}
	types.ApplyTableSettings(cfg)

	return cfg
}

// serialStatusVariant returns the badge variant for a serial status.
func serialStatusVariant(status string) string {
	switch status {
	case "available":
		return "success"
	case "sold":
		return "default"
	case "reserved":
		return "warning"
	case "defective":
		return "danger"
	case "returned":
		return "warning"
	default:
		return "default"
	}
}

// computeSerialSummary computes serial count totals.
func computeSerialSummary(serials []*inventoryserialpb.InventorySerial) *SerialSummary {
	summary := &SerialSummary{Total: len(serials)}
	for _, s := range serials {
		switch s.GetStatus() {
		case "available":
			summary.Available++
		case "sold":
			summary.Sold++
		case "reserved":
			summary.Reserved++
		}
	}
	return summary
}

// computeAvailable computes available quantity from on-hand minus reserved.
func computeAvailable(onHand, reserved float64) string {
	avail := onHand - reserved
	if avail < 0 {
		avail = 0
	}
	if avail == float64(int64(avail)) {
		return fmt.Sprintf("%d", int64(avail))
	}
	return fmt.Sprintf("%.2f", avail)
}

// itemTypeDisplayLabel returns a human-readable label for the item type.
func itemTypeDisplayLabel(itemType string) string {
	switch itemType {
	case "serialized":
		return "Serialized"
	case "non_serialized":
		return "Non-Serialized"
	case "consumable":
		return "Consumable"
	default:
		return itemType
	}
}

// itemTypeDisplayVariant returns the badge variant for the item type.
func itemTypeDisplayVariant(itemType string) string {
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

// buildVariantTabItems creates the tab items for the variant detail page.
// Reuses the same tabs as the variant page so the user stays in context.
func buildVariantTabItems(productID, variantID string, l centymo.ProductLabels) []pyeza.TabItem {
	base := fmt.Sprintf("/app/products/detail/%s/variant/%s", productID, variantID)
	action := fmt.Sprintf("/action/products/detail/%s/variant/%s/tab/", productID, variantID)
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "pricing", Label: l.Tabs.Pricing, Href: base + "?tab=pricing", HxGet: action + "pricing", Icon: "icon-tag", Count: 0, Disabled: false},
		{Key: "stock", Label: "Stock", Href: base + "?tab=stock", HxGet: action + "stock", Icon: "icon-package", Count: 0, Disabled: false},
		{Key: "audit-trail", Label: "Audit Trail", Href: base + "?tab=audit-trail", HxGet: action + "audit-trail", Icon: "icon-clock", Count: 0, Disabled: false},
	}
}
