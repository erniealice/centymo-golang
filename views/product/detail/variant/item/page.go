package item

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/hybra-golang/views/attachment"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	detail "github.com/erniealice/centymo-golang/views/product/detail"
	"github.com/erniealice/centymo-golang/views/product/detail/variant"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
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

// StockDetailPageData holds data for the inventory item detail page.
type StockDetailPageData struct {
	types.PageData
	ContentTemplate string
	Breadcrumbs     []detail.Breadcrumb
	ProductID       string
	VariantID       string
	InventoryItemID string
	ActiveTab       string
	TabItems        []pyeza.TabItem
	// Item info fields
	ItemName            string
	ItemSKU             string
	TrackingMode        string
	TrackingModeLabel   string
	TrackingModeVariant string
	LocationName        string
	QuantityOnHand      string
	QuantityReserved    string
	AvailableQty        string
	ItemStatus          string
	ItemStatusVariant   string
	// Serial data
	SerialTable   *types.TableConfig
	SerialSummary *SerialSummary
	// Attachments tab
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
	Labels              centymo.ProductLabels
}

// NewPageView creates the inventory item detail view (full page).
// Route: /app/products/detail/{id}/variant/{vid}/stock/{iid}
func NewPageView(deps *variant.DetailViewDeps) view.View {
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
		item, err := ReadInventoryItem(ctx, deps, itemID)
		if err != nil {
			return view.Error(err)
		}

		name := item.GetName()
		locationID := item.GetLocationId()
		locationName := centymo.LocationDisplayName(locationID)
		headerTitle := name + " \u2014 " + locationName

		trackingMode := product.GetTrackingMode()
		if trackingMode == "" {
			trackingMode = "bulk"
		}

		active := item.GetActive()
		itemStatus := "active"
		if !active {
			itemStatus = "inactive"
		}

		available := ComputeAvailable(item.GetQuantityOnHand(), item.GetQuantityReserved())

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		l := deps.Labels

		breadcrumbs := []detail.Breadcrumb{
			{Label: l.Breadcrumb.Products, Href: route.ResolveURL(deps.Routes.ListURL, "status", "active")},
			{Label: productName, Href: route.ResolveURL(deps.Routes.DetailURL, "id", productID) + "?tab=variants"},
			{Label: variantSKU, Href: route.ResolveURL(deps.Routes.VariantDetailURL, "id", productID, "vid", variantID) + "?tab=stock"},
			{Label: name + " @ " + locationName, Href: ""},
		}

		tabItems := BuildStockTabItems(productID, variantID, itemID, l, deps.Routes)

		pageData := &StockDetailPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    headerTitle,
				HeaderSubtitle: item.GetSku(),
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:     "variant-stock-detail-content",
			Breadcrumbs:         breadcrumbs,
			ProductID:           productID,
			VariantID:           variantID,
			InventoryItemID:     itemID,
			ActiveTab:           activeTab,
			TabItems:            tabItems,
			ItemName:            name,
			ItemSKU:             item.GetSku(),
			TrackingMode:        trackingMode,
			TrackingModeLabel:   TrackingModeDisplayLabel(trackingMode),
			TrackingModeVariant: TrackingModeDisplayVariant(trackingMode),
			LocationName:        locationName,
			QuantityOnHand:      fmt.Sprintf("%v", item.GetQuantityOnHand()),
			QuantityReserved:    fmt.Sprintf("%v", item.GetQuantityReserved()),
			AvailableQty:        available,
			ItemStatus:          itemStatus,
			ItemStatusVariant:   detail.StatusVariant(itemStatus),
			Labels:              l,
		}

		// Load tab-specific data
		switch activeTab {
		case "serials":
			serials := LoadSerials(ctx, deps, itemID)
			pageData.SerialTable = BuildSerialTable(serials, deps.TableLabels, productID, variantID, itemID, l, deps.Routes)
			pageData.SerialSummary = ComputeSerialSummary(serials)
		case "attachments":
			if deps.ListAttachments != nil {
				cfg := stockAttachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, itemID)
				if err != nil {
					log.Printf("Failed to list attachments: %v", err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, itemID)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.VariantStockAttachmentUploadURL, "id", productID, "vid", variantID, "iid", itemID)
		}

		return view.OK("variant-stock-detail", pageData)
	})
}

// NewTabAction creates the HTMX tab action view for inventory item detail (partial).
// Route: /action/products/detail/{id}/variant/{vid}/stock/{iid}/tab/{tab}
func NewTabAction(deps *variant.DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")
		variantID := viewCtx.Request.PathValue("vid")
		itemID := viewCtx.Request.PathValue("iid")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		// Load product for tracking mode
		prodResp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
			Data: &productpb.Product{Id: productID},
		})
		if err != nil || len(prodResp.GetData()) == 0 {
			log.Printf("Failed to read product %s: %v", productID, err)
			return view.Error(fmt.Errorf("failed to load product: %w", err))
		}
		product := prodResp.GetData()[0]

		trackingMode := product.GetTrackingMode()
		if trackingMode == "" {
			trackingMode = "bulk"
		}

		// Load inventory item
		item, err := ReadInventoryItem(ctx, deps, itemID)
		if err != nil {
			return view.Error(err)
		}

		name := item.GetName()
		locationID := item.GetLocationId()
		locationName := centymo.LocationDisplayName(locationID)

		active := item.GetActive()
		itemStatus := "active"
		if !active {
			itemStatus = "inactive"
		}

		available := ComputeAvailable(item.GetQuantityOnHand(), item.GetQuantityReserved())

		l := deps.Labels

		pageData := &StockDetailPageData{
			ProductID:           productID,
			VariantID:           variantID,
			InventoryItemID:     itemID,
			ActiveTab:           tab,
			ItemName:            name,
			ItemSKU:             item.GetSku(),
			TrackingMode:        trackingMode,
			TrackingModeLabel:   TrackingModeDisplayLabel(trackingMode),
			TrackingModeVariant: TrackingModeDisplayVariant(trackingMode),
			LocationName:        locationName,
			QuantityOnHand:      fmt.Sprintf("%v", item.GetQuantityOnHand()),
			QuantityReserved:    fmt.Sprintf("%v", item.GetQuantityReserved()),
			AvailableQty:        available,
			ItemStatus:          itemStatus,
			ItemStatusVariant:   detail.StatusVariant(itemStatus),
			Labels:              l,
		}

		// Load tab-specific data
		switch tab {
		case "serials":
			serials := LoadSerials(ctx, deps, itemID)
			pageData.SerialTable = BuildSerialTable(serials, deps.TableLabels, productID, variantID, itemID, l, deps.Routes)
			pageData.SerialSummary = ComputeSerialSummary(serials)
		case "attachments":
			if deps.ListAttachments != nil {
				cfg := stockAttachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, itemID)
				if err != nil {
					log.Printf("Failed to list attachments: %v", err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, itemID)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.VariantStockAttachmentUploadURL, "id", productID, "vid", variantID, "iid", itemID)
		}

		templateName := "stock-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// BuildStockTabItems creates the tab items for the inventory item detail page.
func BuildStockTabItems(productID, variantID, itemID string, l centymo.ProductLabels, routes centymo.ProductRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.VariantStockDetailURL, "id", productID, "vid", variantID, "iid", itemID)
	action := route.ResolveURL(routes.VariantStockTabActionURL, "id", productID, "vid", variantID, "iid", itemID, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "serials", Label: l.Tabs.Serials, Href: base + "?tab=serials", HxGet: action + "serials", Icon: "icon-hash", Count: 0, Disabled: false},
		{Key: "pricing-history", Label: l.Tabs.PricingHistory, Href: base + "?tab=pricing-history", HxGet: action + "pricing-history", Icon: "icon-tag", Count: 0, Disabled: false},
		{Key: "attachments", Label: l.Tabs.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip", Count: 0, Disabled: false},
		{Key: "audit-trail", Label: l.Tabs.AuditTrail, Href: base + "?tab=audit-trail", HxGet: action + "audit-trail", Icon: "icon-clock", Count: 0, Disabled: false},
	}
}

// ReadInventoryItem reads a single inventory item by ID.
func ReadInventoryItem(ctx context.Context, deps *variant.DetailViewDeps, itemID string) (*inventoryitempb.InventoryItem, error) {
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
	for _, i := range resp.GetData() {
		if i.GetId() == itemID {
			return i, nil
		}
	}
	return nil, fmt.Errorf("inventory item not found")
}

// LoadSerials loads serial numbers for an inventory item.
func LoadSerials(ctx context.Context, deps *variant.DetailViewDeps, inventoryItemID string) []*inventoryserialpb.InventorySerial {
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

// BuildSerialTable builds the serial numbers table with view actions linking to serial detail.
func BuildSerialTable(serials []*inventoryserialpb.InventorySerial, tableLabels types.TableLabels, productID, variantID, itemID string, l centymo.ProductLabels, routes centymo.ProductRoutes) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "serial_number", Label: l.Detail.SerialNumber},
		{Key: "imei", Label: l.Detail.IMEI, NoSort: true},
		{Key: "status", Label: l.Detail.Status},
		{Key: "warranty_end", Label: l.Detail.WarrantyEnd},
		{Key: "purchase_order", Label: l.Detail.PurchaseOrder, NoSort: true},
	}

	rows := []types.TableRow{}
	for _, s := range serials {
		id := s.GetId()
		serial := s.GetSerialNumber()
		imei := s.GetImei()
		status := s.GetStatus()
		warrantyEnd := s.GetWarrantyEnd()
		po := s.GetPurchaseOrder()

		actions := []types.TableAction{
			{
				Type:  "view",
				Label: l.Actions.View,
				Href:  route.ResolveURL(routes.VariantSerialDetailURL, "id", productID, "vid", variantID, "iid", itemID, "sid", id),
			},
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: serial},
				{Type: "text", Value: imei},
				{Type: "badge", Value: status, Variant: SerialStatusVariant(status)},
				{Type: "text", Value: warrantyEnd},
				{Type: "text", Value: po},
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   "serial-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowEntries:          true,
		ShowSort:             true,
		ShowDensity:          true,
		ShowActions:          true,
		DefaultSortColumn:    "serial_number",
		DefaultSortDirection: "asc",
		Labels:               tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.NoSerialNumbers,
			Message: l.Detail.NoSerialNumbersMsg,
		},
	}
	types.ApplyTableSettings(cfg)

	return cfg
}

// SerialStatusVariant returns the badge variant for a serial status.
func SerialStatusVariant(status string) string {
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

// ComputeSerialSummary computes serial count totals.
func ComputeSerialSummary(serials []*inventoryserialpb.InventorySerial) *SerialSummary {
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

// ComputeAvailable computes available quantity from on-hand minus reserved.
func ComputeAvailable(onHand, reserved float64) string {
	avail := onHand - reserved
	if avail < 0 {
		avail = 0
	}
	if avail == float64(int64(avail)) {
		return fmt.Sprintf("%d", int64(avail))
	}
	return fmt.Sprintf("%.2f", avail)
}

// TrackingModeDisplayLabel returns a human-readable label for the tracking mode.
func TrackingModeDisplayLabel(trackingMode string) string {
	switch trackingMode {
	case "none":
		return "None"
	case "bulk":
		return "Bulk"
	case "serialized":
		return "Serialized"
	default:
		return trackingMode
	}
}

// TrackingModeDisplayVariant returns the badge variant for the tracking mode.
func TrackingModeDisplayVariant(trackingMode string) string {
	switch trackingMode {
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
