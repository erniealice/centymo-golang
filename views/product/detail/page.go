package detail

import (
	"context"
	"fmt"
	"log"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
)

// Deps holds view dependencies.
type Deps struct {
	Routes       centymo.ProductRoutes
	ReadProduct  func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	Labels       centymo.ProductLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	// DataSource for backward compat
	DB centymo.DataSource

	// Typed proto funcs
	ListProductVariants       func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	ListProductOptions        func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ListProductOptionValues   func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	ListProductVariantOptions func(ctx context.Context, req *productvariantoptionpb.ListProductVariantOptionsRequest) (*productvariantoptionpb.ListProductVariantOptionsResponse, error)
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
	VariantsTable *types.TableConfig
	OptionsTable  *types.TableConfig
}

// NewView creates the product detail view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" || activeTab == "pricing" {
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
	priceFormatted := FormatPrice(currency, product.GetPrice())

	productStatus := "active"
	if !product.GetActive() {
		productStatus = "inactive"
	}
	StatusVariant := "success"
	if productStatus == "inactive" {
		StatusVariant = "warning"
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
	if deps.ListProductVariants != nil {
		varResp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{})
		if err == nil {
			for _, v := range varResp.GetData() {
				if v.GetProductId() == id {
					variantCount++
				}
			}
		}
	}

	optionCount := getOptionCountTyped(ctx, deps, id)

	l := deps.Labels
	tabItems := buildTabItems(id, l, variantCount, optionCount, deps.Routes)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          name,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
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
		StatusVariant:   StatusVariant,
		Collections:     collections,
	}

	// Load tab-specific data
	perms := view.GetUserPermissions(ctx)
	switch activeTab {
	case "variants":
		tableConfig := BuildVariantsTable(ctx, deps, id, perms)
		pageData.VariantsTable = tableConfig
	case "options":
		tableConfig := buildOptionsTable(ctx, deps, id)
		pageData.OptionsTable = tableConfig
	}

	return pageData, nil
}

func buildTabItems(id string, l centymo.ProductLabels, variantCount, optionCount int, routes centymo.ProductRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "options", Label: l.Tabs.Options, Href: base + "?tab=options", HxGet: action + "options", Icon: "icon-settings", Count: optionCount, Disabled: false},
		{Key: "variants", Label: l.Tabs.Variants, Href: base + "?tab=variants", HxGet: action + "variants", Icon: "icon-layers", Count: variantCount, Disabled: false},
	}
}

// ---------------------------------------------------------------------------
// Variants tab table
// ---------------------------------------------------------------------------

func BuildVariantsTable(ctx context.Context, deps *Deps, productID string, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "sku", Label: l.Variant.SKU, Sortable: true},
		{Key: "priceOverride", Label: l.Variant.PriceOverride, Sortable: true, Width: "150px"},
		{Key: "options", Label: "Options", Sortable: false},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}

	rows := []types.TableRow{}

	// Build variant → option value labels map
	variantOptionLabels := make(map[string][]string)
	if deps.ListProductVariantOptions != nil && deps.ListProductOptionValues != nil {
		voResp, _ := deps.ListProductVariantOptions(ctx, &productvariantoptionpb.ListProductVariantOptionsRequest{})
		ovResp, _ := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
		valLabelMap := make(map[string]string)
		if ovResp != nil {
			for _, ov := range ovResp.GetData() {
				vid := ov.GetId()
				label := ov.GetLabel()
				if vid != "" {
					valLabelMap[vid] = label
				}
			}
		}
		if voResp != nil {
			for _, vo := range voResp.GetData() {
				varID := vo.GetProductVariantId()
				valID := vo.GetProductOptionValueId()
				if label, ok := valLabelMap[valID]; ok && varID != "" {
					variantOptionLabels[varID] = append(variantOptionLabels[varID], label)
				}
			}
		}
	}

	if deps.ListProductVariants != nil {
		varResp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{})
		if err != nil {
			log.Printf("Failed to list product variants: %v", err)
		} else {
			for _, v := range varResp.GetData() {
				pid := v.GetProductId()
				if pid != productID {
					continue
				}

				vid := v.GetId()
				sku := v.GetSku()
				priceStr := ""
				if po := v.GetPriceOverride(); po != 0 {
					priceStr = fmt.Sprintf("%v", po)
				}
				active := v.GetActive()
				status := "active"
				if !active {
					status = "inactive"
				}

				optionsDisplay := strings.Join(variantOptionLabels[vid], ", ")

				actions := []types.TableAction{
					{
						Type: "view", Label: l.Actions.View,
						Href: route.ResolveURL(deps.Routes.VariantDetailURL, "id", productID, "vid", vid),
					},
					{
						Type: "edit", Label: l.Variant.Edit, Action: "edit",
						URL:             route.ResolveURL(deps.Routes.VariantEditURL, "id", productID, "vid", vid),
						DrawerTitle:     l.Variant.Edit,
						Disabled:        !perms.Can("product", "update"),
						DisabledTooltip: "No permission",
					},
					{
						Type: "delete", Label: l.Variant.Remove, Action: "delete",
						URL:             route.ResolveURL(deps.Routes.VariantRemoveURL, "id", productID),
						ItemName:        sku,
						ConfirmTitle:    l.Variant.Remove,
						ConfirmMessage:  fmt.Sprintf("Are you sure you want to remove variant %s?", sku),
						Disabled:        !perms.Can("product", "delete"),
						DisabledTooltip: "No permission",
					},
				}

				rows = append(rows, types.TableRow{
					ID: vid,
					Cells: []types.TableCell{
						{Type: "text", Value: sku},
						{Type: "text", Value: priceStr},
						{Type: "text", Value: optionsDisplay},
						{Type: "badge", Value: status, Variant: StatusVariant(status)},
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
		RefreshURL:           route.ResolveURL(deps.Routes.VariantTableURL, "id", productID),
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
			Label:           l.Variant.Assign,
			ActionURL:       route.ResolveURL(deps.Routes.VariantAssignURL, "id", productID),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("product", "create"),
			DisabledTooltip: "No permission",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func FormatPrice(currency string, price float64) string {
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

func StatusVariant(status string) string {
	switch status {
	case "active":
		return "success"
	case "inactive":
		return "warning"
	default:
		return "default"
	}
}
