package detail

import (
	"context"
	"fmt"
	"log"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	priceproductpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_product"
)

// Deps holds view dependencies.
type Deps struct {
	ReadProduct       func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	ListPriceProducts func(ctx context.Context, req *priceproductpb.ListPriceProductsRequest) (*priceproductpb.ListPriceProductsResponse, error)
	Labels            centymo.ProductLabels
	CommonLabels      pyeza.CommonLabels
	TableLabels       types.TableLabels
	// DataSource for variant/attribute queries (raw DB)
	DB centymo.DataSource
}

// PageData holds the data for the product detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Product         *productpb.Product
	Labels          centymo.ProductLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	ID              string
	ProductName     string
	ProductDesc     string
	ProductPrice    string
	ProductCurrency string
	ProductStatus   string
	StatusVariant   string
	Collections     []string
	VariantsTable   *types.TableConfig
	AttributesTable *types.TableConfig
	PricingTable    *types.TableConfig
}

// NewView creates the product detail view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("product-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
// Handles GET /action/products/detail/{id}/tab/{tab}
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		// Return only the tab partial template
		templateName := "product-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

// buildPageData loads product data and builds the PageData for the given active tab.
func buildPageData(ctx context.Context, deps *Deps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
		Data: &productpb.Product{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read product %s: %v", id, err)
		return nil, fmt.Errorf("failed to load product: %w", err)
	}

	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("product not found")
	}
	product := data[0]

	name := product.GetName()
	description := product.GetDescription()
	currency := product.GetCurrency()
	if currency == "" {
		currency = "PHP"
	}
	priceFormatted := formatPrice(currency, product.GetPrice())

	productStatus := "active"
	if !product.GetActive() {
		productStatus = "inactive"
	}
	statusVariant := "success"
	if productStatus == "inactive" {
		statusVariant = "warning"
	}

	// Extract collection names
	var collections []string
	for _, pc := range product.GetProductCollections() {
		if pc != nil {
			collections = append(collections, pc.GetCollectionId())
		}
	}

	// Get counts for tab badges
	variantCount := 0
	attributeCount := 0

	if deps.DB != nil {
		variants, err := deps.DB.ListSimple(ctx, "product_variant")
		if err == nil {
			for _, v := range variants {
				if pid, _ := v["product_id"].(string); pid == id {
					variantCount++
				}
			}
		}
		attributes, err := deps.DB.ListSimple(ctx, "product_attribute")
		if err == nil {
			for _, a := range attributes {
				if pid, _ := a["product_id"].(string); pid == id {
					attributeCount++
				}
			}
		}
	}

	l := deps.Labels
	tabItems := buildTabItems(id, l, variantCount, attributeCount)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          name,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "products",
			HeaderTitle:    name,
			HeaderSubtitle: description,
			HeaderIcon:     "icon-package",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "product-detail-content",
		Product:         product,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		ProductName:     name,
		ProductDesc:     description,
		ProductPrice:    priceFormatted,
		ProductCurrency: currency,
		ProductStatus:   productStatus,
		StatusVariant:   statusVariant,
		Collections:     collections,
	}

	// Load tab-specific data
	switch activeTab {
	case "variants":
		tableConfig := buildVariantsTable(ctx, deps, id)
		pageData.VariantsTable = tableConfig
	case "attributes":
		tableConfig := buildAttributesTable(ctx, deps, id)
		pageData.AttributesTable = tableConfig
	case "pricing":
		tableConfig := buildPricingTable(ctx, deps, id)
		pageData.PricingTable = tableConfig
	}

	return pageData, nil
}

func buildTabItems(id string, l centymo.ProductLabels, variantCount, attributeCount int) []pyeza.TabItem {
	base := "/app/products/detail/" + id
	action := "/action/products/detail/" + id + "/tab/"
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "variants", Label: l.Tabs.Variants, Href: base + "?tab=variants", HxGet: action + "variants", Icon: "icon-layers", Count: variantCount, Disabled: false},
		{Key: "attributes", Label: l.Tabs.Attributes, Href: base + "?tab=attributes", HxGet: action + "attributes", Icon: "icon-sliders", Count: attributeCount, Disabled: false},
		{Key: "pricing", Label: l.Tabs.Pricing, Href: base + "?tab=pricing", HxGet: action + "pricing", Icon: "icon-tag", Count: 0, Disabled: false},
	}
}

// ---------------------------------------------------------------------------
// Variants tab table
// ---------------------------------------------------------------------------

func buildVariantsTable(ctx context.Context, deps *Deps, productID string) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "sku", Label: l.Variant.SKU, Sortable: true},
		{Key: "priceOverride", Label: l.Variant.PriceOverride, Sortable: true, Width: "150px"},
		{Key: "attributes", Label: l.Variant.Attributes, Sortable: false},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}

	rows := []types.TableRow{}

	if deps.DB != nil {
		variants, err := deps.DB.ListSimple(ctx, "product_variant")
		if err != nil {
			log.Printf("Failed to list product variants: %v", err)
		} else {
			for _, v := range variants {
				pid, _ := v["product_id"].(string)
				if pid != productID {
					continue
				}

				vid, _ := v["id"].(string)
				sku, _ := v["sku"].(string)
				priceStr := ""
				if po, ok := v["price_override"]; ok && po != nil {
					priceStr = fmt.Sprintf("%v", po)
				}
				active, _ := v["active"].(bool)
				status := "active"
				if !active {
					status = "inactive"
				}

				// Collect attribute values as a display string
				attrDisplay := ""
				if attrs, ok := v["attribute_values"].(string); ok {
					attrDisplay = attrs
				}

				actions := []types.TableAction{
					{
						Type: "edit", Label: l.Variant.Edit, Action: "edit",
						URL:         fmt.Sprintf("/action/products/detail/%s/variants/edit/%s", productID, vid),
						DrawerTitle: l.Variant.Edit,
					},
					{
						Type: "delete", Label: l.Variant.Remove, Action: "delete",
						URL:            fmt.Sprintf("/action/products/detail/%s/variants/remove", productID),
						ItemName:       sku,
						ConfirmTitle:   l.Variant.Remove,
						ConfirmMessage: fmt.Sprintf("Are you sure you want to remove variant %s?", sku),
					},
				}

				rows = append(rows, types.TableRow{
					ID: vid,
					Cells: []types.TableCell{
						{Type: "text", Value: sku},
						{Type: "text", Value: priceStr},
						{Type: "text", Value: attrDisplay},
						{Type: "badge", Value: status, Variant: statusVariant(status)},
					},
					DataAttrs: map[string]string{
						"sku":    sku,
						"status": status,
					},
					Actions: actions,
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "product-variants-table",
		RefreshURL:           fmt.Sprintf("/action/products/detail/%s/variants/table", productID),
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
			Title:   l.Variant.Empty,
			Message: "No variants have been added to this product yet.",
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     l.Variant.Assign,
			ActionURL: fmt.Sprintf("/action/products/detail/%s/variants/assign", productID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Attributes tab table
// ---------------------------------------------------------------------------

func buildAttributesTable(ctx context.Context, deps *Deps, productID string) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "code", Label: "Code", Sortable: true},
		{Key: "dataType", Label: "Data Type", Sortable: true, Width: "120px"},
		{Key: "defaultValue", Label: l.Attribute.DefaultValue, Sortable: false, Width: "150px"},
	}

	rows := []types.TableRow{}

	if deps.DB != nil {
		attributes, err := deps.DB.ListSimple(ctx, "product_attribute")
		if err != nil {
			log.Printf("Failed to list product attributes: %v", err)
		} else {
			for _, a := range attributes {
				pid, _ := a["product_id"].(string)
				if pid != productID {
					continue
				}

				aid, _ := a["id"].(string)
				name, _ := a["attribute_name"].(string)
				code, _ := a["attribute_code"].(string)
				dataType, _ := a["data_type"].(string)
				defaultVal, _ := a["default_value"].(string)

				actions := []types.TableAction{
					{
						Type: "delete", Label: l.Attribute.Remove, Action: "delete",
						URL:            fmt.Sprintf("/action/products/detail/%s/attributes/remove", productID),
						ItemName:       name,
						ConfirmTitle:   l.Attribute.Remove,
						ConfirmMessage: fmt.Sprintf("Are you sure you want to remove attribute %s?", name),
					},
				}

				rows = append(rows, types.TableRow{
					ID: aid,
					Cells: []types.TableCell{
						{Type: "text", Value: name},
						{Type: "text", Value: code},
						{Type: "badge", Value: dataType, Variant: "info"},
						{Type: "text", Value: defaultVal},
					},
					DataAttrs: map[string]string{
						"name": name,
						"code": code,
					},
					Actions: actions,
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "product-attributes-table",
		RefreshURL:           fmt.Sprintf("/action/products/detail/%s/attributes/table", productID),
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
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Attribute.Empty,
			Message: "No attributes have been assigned to this product yet.",
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     l.Attribute.Assign,
			ActionURL: fmt.Sprintf("/action/products/detail/%s/attributes/assign", productID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Pricing tab table
// ---------------------------------------------------------------------------

func buildPricingTable(ctx context.Context, deps *Deps, productID string) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "priceListName", Label: "Price List", Sortable: true},
		{Key: "currency", Label: l.Detail.Currency, Sortable: true, Width: "100px"},
		{Key: "customPrice", Label: "Custom Price", Sortable: true, Width: "150px"},
		{Key: "validFrom", Label: "Valid From", Sortable: true, Width: "140px"},
		{Key: "validTo", Label: "Valid To", Sortable: true, Width: "140px"},
	}

	rows := []types.TableRow{}

	if deps.ListPriceProducts != nil {
		resp, err := deps.ListPriceProducts(ctx, &priceproductpb.ListPriceProductsRequest{})
		if err != nil {
			log.Printf("Failed to list price products: %v", err)
		} else {
			for _, pp := range resp.GetData() {
				if pp.GetProductId() != productID {
					continue
				}

				ppID := pp.GetId()
				priceListName := pp.GetPriceListId()
				currency := pp.GetCurrency()
				customPrice := fmt.Sprintf("%d", pp.GetAmount())
				validFrom := pp.GetDateStartString()
				validTo := pp.GetDateEndString()

				rows = append(rows, types.TableRow{
					ID: ppID,
					Cells: []types.TableCell{
						{Type: "text", Value: priceListName},
						{Type: "text", Value: currency},
						{Type: "text", Value: customPrice},
						{Type: "text", Value: validFrom},
						{Type: "text", Value: validTo},
					},
					DataAttrs: map[string]string{
						"priceListName": priceListName,
					},
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "product-pricing-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          false,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "priceListName",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   "No Price Lists",
			Message: "This product has not been added to any price lists yet.",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func formatPrice(currency string, price float64) string {
	if currency == "" {
		currency = "PHP"
	}
	raw := fmt.Sprintf("%.2f", price)
	parts := strings.SplitN(raw, ".", 2)
	intPart := parts[0]
	decPart := "00"
	if len(parts) > 1 {
		decPart = parts[1]
	}

	n := len(intPart)
	if n <= 3 {
		return currency + " " + intPart + "." + decPart
	}
	var result []byte
	for i, c := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, byte(c))
	}
	return currency + " " + string(result) + "." + decPart
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
