package variant

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/erniealice/hybra-golang/views/attachment"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	detail "github.com/erniealice/centymo-golang/views/product/detail"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
)

// OptionEntry represents a name-value pair for a variant's option selection.
type OptionEntry struct {
	Name  string
	Value string
}

// VariantPageData holds data for the variant detail page.
type VariantPageData struct {
	types.PageData
	ContentTemplate string
	Breadcrumbs     []detail.Breadcrumb
	ProductID       string
	VariantID       string
	ActiveTab       string
	TabItems        []pyeza.TabItem
	// Info tab
	VariantName   string
	VariantSKU    string
	VariantPrice  string
	VariantStatus string
	StatusVariant string
	OptionEntries []OptionEntry
	// Stock tab
	StockTable *types.TableConfig
	// Images tab
	Images []ImageData
	// Attachments tab
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
	Labels              centymo.ProductLabels
}

// NewPageView creates the variant detail view (full page).
// Route: /app/products/detail/{id}/variant/{vid}
func NewPageView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		vid := viewCtx.Request.PathValue("vid")

		// Load product for breadcrumb context
		prodResp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
			Data: &productpb.Product{Id: id},
		})
		if err != nil || len(prodResp.GetData()) == 0 {
			log.Printf("Failed to read product %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load product: %w", err))
		}
		product := prodResp.GetData()[0]
		productName := product.GetName()
		currency := product.GetCurrency()
		if currency == "" {
			currency = "PHP"
		}

		// Load variant
		varResp, err := deps.ReadProductVariant(ctx, &productvariantpb.ReadProductVariantRequest{
			Data: &productvariantpb.ProductVariant{Id: vid},
		})
		if err != nil || len(varResp.GetData()) == 0 {
			log.Printf("Failed to read product_variant %s: %v", vid, err)
			return view.Error(fmt.Errorf("failed to load variant: %w", err))
		}
		variant := varResp.GetData()[0]
		sku := variant.GetSku()
		active := variant.GetActive()

		variantStatus := "active"
		if !active {
			variantStatus = "inactive"
		}

		// Format price override
		priceOverride := variant.GetPriceOverride()
		variantPrice := ""
		if priceOverride != 0 {
			variantPrice = detail.FormatPrice(currency, float64(priceOverride))
		}

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		l := deps.Labels

		// Load option entries early — needed for page title and info tab
		optionEntries := loadVariantOptionEntries(ctx, deps, id, vid)

		// Build title from product name + option values (e.g., "iPhone 16 Pro Max — Black / 256GB")
		headerTitle := productName
		if len(optionEntries) > 0 {
			vals := make([]string, len(optionEntries))
			for i, e := range optionEntries {
				vals[i] = e.Value
			}
			headerTitle = productName + " - " + strings.Join(vals, " | ")
		}

		breadcrumbs := []detail.Breadcrumb{
			{Label: l.Breadcrumb.Products, Href: route.ResolveURL(deps.Routes.ListURL, "status", "active")},
			{Label: productName, Href: route.ResolveURL(deps.Routes.DetailURL, "id", id) + "?tab=variants"},
			{Label: sku, Href: ""},
		}

		tabItems := buildVariantTabItems(id, vid, l, deps.Routes)

		pageData := &VariantPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    headerTitle,
				HeaderSubtitle: sku,
				HeaderIcon:     "icon-layers",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "variant-detail-content",
			Breadcrumbs:     breadcrumbs,
			ProductID:       id,
			VariantID:       vid,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
			VariantName:     productName,
			VariantSKU:      sku,
			VariantPrice:    variantPrice,
			VariantStatus:   variantStatus,
			StatusVariant:   detail.StatusVariant(variantStatus),
			OptionEntries:   optionEntries,
			Labels:          l,
		}

		// Load tab-specific data
		switch activeTab {
		case "stock":
			pageData.StockTable = buildStockTable(ctx, deps, id, vid)
		case "images":
			pageData.Images = loadVariantImages(ctx, deps, vid)
		case "attachments":
			if deps.ListAttachments != nil {
				cfg := variantAttachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, vid)
				if err != nil {
					log.Printf("Failed to list attachments: %v", err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, vid)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.VariantAttachmentUploadURL, "id", id, "vid", vid)
		}

		return view.OK("variant-detail", pageData)
	})
}

// NewTabAction creates the HTMX tab action view for variant detail (partial).
// Route: /action/products/detail/{id}/variant/{vid}/tab/{tab}
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		vid := viewCtx.Request.PathValue("vid")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		// Load product for context
		prodResp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
			Data: &productpb.Product{Id: id},
		})
		if err != nil || len(prodResp.GetData()) == 0 {
			log.Printf("Failed to read product %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load product: %w", err))
		}
		product := prodResp.GetData()[0]
		productName := product.GetName()
		currency := product.GetCurrency()
		if currency == "" {
			currency = "PHP"
		}

		// Load variant
		varResp, err := deps.ReadProductVariant(ctx, &productvariantpb.ReadProductVariantRequest{
			Data: &productvariantpb.ProductVariant{Id: vid},
		})
		if err != nil || len(varResp.GetData()) == 0 {
			log.Printf("Failed to read product_variant %s: %v", vid, err)
			return view.Error(fmt.Errorf("failed to load variant: %w", err))
		}
		variant := varResp.GetData()[0]
		sku := variant.GetSku()
		active := variant.GetActive()

		variantStatus := "active"
		if !active {
			variantStatus = "inactive"
		}

		priceOverride := variant.GetPriceOverride()
		variantPrice := ""
		if priceOverride != 0 {
			variantPrice = detail.FormatPrice(currency, float64(priceOverride))
		}

		l := deps.Labels

		pageData := &VariantPageData{
			ProductID:     id,
			VariantID:     vid,
			ActiveTab:     tab,
			VariantName:   productName,
			VariantSKU:    sku,
			VariantPrice:  variantPrice,
			VariantStatus: variantStatus,
			StatusVariant: detail.StatusVariant(variantStatus),
			Labels:        l,
		}

		// Load tab-specific data
		switch tab {
		case "info":
			pageData.OptionEntries = loadVariantOptionEntries(ctx, deps, id, vid)
		case "stock":
			pageData.StockTable = buildStockTable(ctx, deps, id, vid)
		case "images":
			pageData.Images = loadVariantImages(ctx, deps, vid)
		case "attachments":
			if deps.ListAttachments != nil {
				cfg := variantAttachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, vid)
				if err != nil {
					log.Printf("Failed to list attachments: %v", err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, vid)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.VariantAttachmentUploadURL, "id", id, "vid", vid)
		}

		templateName := "variant-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// buildVariantTabItems creates the tab items for the variant detail page.
func buildVariantTabItems(productID, variantID string, l centymo.ProductLabels, routes centymo.ProductRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.VariantDetailURL, "id", productID, "vid", variantID)
	action := route.ResolveURL(routes.VariantTabActionURL, "id", productID, "vid", variantID, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "images", Label: l.Tabs.Images, Href: base + "?tab=images", HxGet: action + "images", Icon: "icon-image", Count: 0, Disabled: false},
		{Key: "pricing", Label: l.Tabs.Pricing, Href: base + "?tab=pricing", HxGet: action + "pricing", Icon: "icon-tag", Count: 0, Disabled: false},
		{Key: "stock", Label: l.Tabs.Stock, Href: base + "?tab=stock", HxGet: action + "stock", Icon: "icon-package", Count: 0, Disabled: false},
		{Key: "attachments", Label: l.Tabs.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip", Count: 0, Disabled: false},
		{Key: "audit-trail", Label: l.Tabs.AuditTrail, Href: base + "?tab=audit-trail", HxGet: action + "audit-trail", Icon: "icon-clock", Count: 0, Disabled: false},
	}
}

// loadVariantOptionEntries loads ALL product options with their assigned values for a variant.
// Options without an assigned value show "\u2014" as the value.
func loadVariantOptionEntries(ctx context.Context, deps *DetailViewDeps, productID, variantID string) []OptionEntry {
	if deps.ListProductOptions == nil || deps.ListProductOptionValues == nil {
		return nil
	}

	// 1. Load ALL active product options for this product
	optResp, err := deps.ListProductOptions(ctx, &productoptionpb.ListProductOptionsRequest{})
	if err != nil {
		log.Printf("Failed to list product_option: %v", err)
		return nil
	}

	var productOptions []*productoptionpb.ProductOption
	for _, o := range optResp.GetData() {
		if o.GetProductId() == productID && o.GetActive() {
			productOptions = append(productOptions, o)
		}
	}
	if len(productOptions) == 0 {
		return nil
	}

	// 2. Load product option values for label lookup
	valResp, err := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
	if err != nil {
		log.Printf("Failed to list product_option_value: %v", err)
		return nil
	}
	valueMap := map[string]*productoptionvaluepb.ProductOptionValue{}
	for _, v := range valResp.GetData() {
		valueMap[v.GetId()] = v
	}

	// 3. Load variant option assignments (if any)
	assignedValues := map[string]string{} // optionID -> valueLabel
	if deps.ListProductVariantOptions != nil {
		voResp, err := deps.ListProductVariantOptions(ctx, &productvariantoptionpb.ListProductVariantOptionsRequest{})
		if err == nil {
			for _, pvo := range voResp.GetData() {
				if pvo.GetProductVariantId() != variantID {
					continue
				}
				vid := pvo.GetProductOptionValueId()
				if pov, ok := valueMap[vid]; ok {
					oid := pov.GetProductOptionId()
					assignedValues[oid] = pov.GetLabel()
				}
			}
		}
	}

	// 4. Build entries for ALL options
	var entries []OptionEntry
	for _, o := range productOptions {
		oid := o.GetId()
		name := o.GetName()
		value := "\u2014"
		if label, ok := assignedValues[oid]; ok && label != "" {
			value = label
		}
		entries = append(entries, OptionEntry{Name: name, Value: value})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	return entries
}

// buildStockTable builds the stock table showing inventory items for this variant.
func buildStockTable(ctx context.Context, deps *DetailViewDeps, productID, variantID string) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "sku", Label: l.Variant.SKU, Sortable: true},
		{Key: "location", Label: l.Variant.Location, Sortable: true},
		{Key: "quantity", Label: l.Variant.QtyOnHand, Sortable: true, Width: "120px"},
		{Key: "serials", Label: l.Variant.SerialCount, Sortable: true, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}

	rows := []types.TableRow{}

	if deps.ListInventoryItems == nil {
		return nil
	}

	// Load all inventory items, filter by product_variant_id
	itemResp, err := deps.ListInventoryItems(ctx, &inventoryitempb.ListInventoryItemsRequest{})
	if err != nil {
		log.Printf("Failed to list inventory_item for stock tab: %v", err)
		return nil
	}

	// Filter items for this variant first, then count serials per item
	var variantItems []*inventoryitempb.InventoryItem
	for _, item := range itemResp.GetData() {
		if item.GetProductVariantId() == variantID {
			variantItems = append(variantItems, item)
		}
	}

	// Count serials per inventory item using filtered per-item queries
	serialCounts := make(map[string]int)
	if deps.ListInventorySerials != nil {
		for _, item := range variantItems {
			iid := item.GetId()
			serialResp, err := deps.ListInventorySerials(ctx, &inventoryserialpb.ListInventorySerialsRequest{
				InventoryItemId: &iid,
			})
			if err == nil {
				serialCounts[iid] = len(serialResp.GetData())
			}
		}
	}

	for _, item := range variantItems {
		iid := item.GetId()
		sku := item.GetSku()
		locationID := item.GetLocationId()
		locationName := centymo.LocationDisplayName(locationID)

		qtyStr := fmt.Sprintf("%v", item.GetQuantityOnHand())

		serialCount := fmt.Sprintf("%d", serialCounts[iid])

		active := item.GetActive()
		status := "active"
		if !active {
			status = "inactive"
		}

		actions := []types.TableAction{
			{
				Type: "view", Label: l.Actions.View,
				Href: route.ResolveURL(deps.Routes.VariantStockDetailURL, "id", productID, "vid", variantID, "iid", iid),
			},
		}

		rows = append(rows, types.TableRow{
			ID: iid,
			Cells: []types.TableCell{
				{Type: "text", Value: sku},
				{Type: "text", Value: locationName},
				{Type: "text", Value: qtyStr},
				{Type: "text", Value: serialCount},
				{Type: "badge", Value: status, Variant: detail.StatusVariant(status)},
			},
			Actions: actions,
		})
	}

	if len(rows) == 0 {
		return nil
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "variant-stock-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "sku",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Variant.NoStock,
			Message: l.Variant.NoStockMsg,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}
