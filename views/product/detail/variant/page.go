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
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
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
	// Pricing tab
	PricingTable *types.TableConfig
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

		// Breadcrumb data is retained for any legacy consumers; the visible
		// product → variant crumb is rendered via HeaderBreadcrumb on the
		// shared page-header partial (see option_page.go for the same pattern).
		breadcrumbs := []detail.Breadcrumb{
			{Label: l.Breadcrumb.Products, Href: route.ResolveURL(deps.Routes.ListURL, "status", "active")},
			{Label: productName, Href: route.ResolveURL(deps.Routes.DetailURL, "id", id) + "?tab=variants"},
			{Label: sku, Href: ""},
		}

		tabItems := buildVariantTabItems(id, vid, l, deps.Routes)

		pageData := &VariantPageData{
			PageData: types.PageData{
				CacheVersion:        viewCtx.CacheVersion,
				Title:               headerTitle,
				CurrentPath:         viewCtx.CurrentPath,
				ActiveNav:           deps.Routes.ActiveNav,
				ActiveSubNav:        deps.Routes.ActiveSubNav,
				HeaderTitle:         headerTitle,
				HeaderSubtitle:      sku,
				HeaderBreadcrumb:    productName,
				HeaderBreadcrumbURL: route.ResolveURL(deps.Routes.DetailURL, "id", id) + "?tab=variants",
				HeaderIcon:          "icon-layers",
				CommonLabels:        deps.CommonLabels,
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
		case "pricing":
			pageData.PricingTable = buildPricingTable(ctx, deps, id, vid)
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
		case "pricing":
			pageData.PricingTable = buildPricingTable(ctx, deps, id, vid)
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
		{Key: "quantity", Label: l.Variant.QtyOnHand, Sortable: true, WidthClass: "col-2xl"},
		{Key: "serials", Label: l.Variant.SerialCount, Sortable: true, WidthClass: "col-2xl"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, WidthClass: "col-2xl"},
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

// buildPricingTable builds the pricing table for the variant detail page.
// It joins: product_plan → product_price_plan → price_plan → price_schedule → plan.
// Only product_plan rows where both product_id and product_variant_id match are included.
func buildPricingTable(ctx context.Context, deps *DetailViewDeps, productID, variantID string) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "start", Label: "Start Date", Sortable: true},
		{Key: "end", Label: "End Date", Sortable: true},
		{Key: "package", Label: "Package", Sortable: true},
		{Key: "rate_card", Label: "Rate Card", Sortable: true},
		{Key: "amount", Label: "Amount", Sortable: true, WidthClass: "col-2xl"},
	}

	emptyTitle := l.Detail.VariantPricing
	if emptyTitle == "" {
		emptyTitle = "No Pricing"
	}
	emptyMsg := "No pricing plans have been set for this variant yet."

	noData := func() *types.TableConfig {
		return &types.TableConfig{
			ID:      "variant-pricing-table",
			Columns: columns,
			Rows:    []types.TableRow{},
			Labels:  deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   emptyTitle,
				Message: emptyMsg,
			},
		}
	}

	// Guard: all five list deps required.
	if deps.ListProductPlans == nil || deps.ListProductPricePlans == nil ||
		deps.ListPricePlans == nil || deps.ListPriceSchedules == nil || deps.ListPlans == nil {
		return noData()
	}

	// Step 1: find product_plan rows for this (product, variant) pair.
	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
	if err != nil {
		log.Printf("buildPricingTable: ListProductPlans error: %v", err)
		return noData()
	}
	productPlanIDs := map[string]bool{}
	for _, pp := range ppResp.GetData() {
		if pp.GetProductId() == productID && pp.GetProductVariantId() == variantID {
			productPlanIDs[pp.GetId()] = true
		}
	}
	if len(productPlanIDs) == 0 {
		return noData()
	}

	// Step 2: find product_price_plan rows pointing at those product_plan IDs.
	pppResp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{})
	if err != nil {
		log.Printf("buildPricingTable: ListProductPricePlans error: %v", err)
		return noData()
	}
	pricePlanIDs := map[string]bool{}
	for _, ppp := range pppResp.GetData() {
		if productPlanIDs[ppp.GetProductPlanId()] {
			pricePlanIDs[ppp.GetPricePlanId()] = true
		}
	}
	if len(pricePlanIDs) == 0 {
		return noData()
	}

	// Step 3: load price_plan rows.
	plResp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
	if err != nil {
		log.Printf("buildPricingTable: ListPricePlans error: %v", err)
		return noData()
	}
	type pricePlanMeta struct {
		name            string
		billingAmount   int64
		billingCurrency string
		planID          string
		scheduleID      string
	}
	relevantPlans := map[string]pricePlanMeta{}
	scheduleIDs := map[string]bool{}
	planIDs := map[string]bool{}
	for _, pp := range plResp.GetData() {
		if !pricePlanIDs[pp.GetId()] {
			continue
		}
		n := ""
		if pp.Name != nil {
			n = pp.GetName()
		}
		relevantPlans[pp.GetId()] = pricePlanMeta{
			name:            n,
			billingAmount:   pp.GetBillingAmount(),
			billingCurrency: pp.GetBillingCurrency(),
			planID:          pp.GetPlanId(),
			scheduleID:      pp.GetPriceScheduleId(),
		}
		if pp.GetPriceScheduleId() != "" {
			scheduleIDs[pp.GetPriceScheduleId()] = true
		}
		if pp.GetPlanId() != "" {
			planIDs[pp.GetPlanId()] = true
		}
	}

	// Step 4: load price_schedule rows for date range + rate-card name.
	schedResp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
	if err != nil {
		log.Printf("buildPricingTable: ListPriceSchedules error: %v", err)
		return noData()
	}
	type scheduleMeta struct {
		name      string
		dateStart string
		dateEnd   string
	}
	scheduleByID := map[string]scheduleMeta{}
	for _, s := range schedResp.GetData() {
		if !scheduleIDs[s.GetId()] {
			continue
		}
		startStr := ""
		if ts := s.GetDateTimeStart(); ts != nil {
			startStr = ts.AsTime().Format("2006-01-02")
		}
		endStr := "—" // em-dash for null
		if ts := s.GetDateTimeEnd(); ts != nil {
			endStr = ts.AsTime().Format("2006-01-02")
		}
		scheduleByID[s.GetId()] = scheduleMeta{
			name:      s.GetName(),
			dateStart: startStr,
			dateEnd:   endStr,
		}
	}

	// Step 5: load plan rows for package name fallback.
	planResp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		log.Printf("buildPricingTable: ListPlans error: %v", err)
		return noData()
	}
	planNameByID := map[string]string{}
	for _, p := range planResp.GetData() {
		if planIDs[p.GetId()] {
			planNameByID[p.GetId()] = p.GetName()
		}
	}

	// Step 6: build table rows, one per relevant price_plan.
	type pricingRow struct {
		dateStart string
		dateEnd   string
		pkg       string
		rateCard  string
		amount    string
	}
	var rows []pricingRow
	for ppID, meta := range relevantPlans {
		_ = ppID
		sched := scheduleByID[meta.scheduleID]

		pkgName := meta.name
		if pkgName == "" {
			pkgName = planNameByID[meta.planID]
		}

		rateCardName := sched.name
		if rateCardName == "" {
			rateCardName = meta.scheduleID
		}

		currency := meta.billingCurrency
		if currency == "" {
			currency = "PHP"
		}
		amountStr := detail.FormatPrice(currency, float64(meta.billingAmount)/100.0)

		rows = append(rows, pricingRow{
			dateStart: sched.dateStart,
			dateEnd:   sched.dateEnd,
			pkg:       pkgName,
			rateCard:  rateCardName,
			amount:    amountStr,
		})
	}

	// Sort by start date ASC for predictability.
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].dateStart < rows[j].dateStart
	})

	tableRows := make([]types.TableRow, 0, len(rows))
	for _, r := range rows {
		tableRows = append(tableRows, types.TableRow{
			Cells: []types.TableCell{
				{Type: "text", Value: r.dateStart},
				{Type: "text", Value: r.dateEnd},
				{Type: "text", Value: r.pkg},
				{Type: "text", Value: r.rateCard},
				{Type: "text", Value: r.amount},
			},
		})
	}

	types.ApplyColumnStyles(columns, tableRows)

	tableConfig := &types.TableConfig{
		ID:                   "variant-pricing-table",
		Columns:              columns,
		Rows:                 tableRows,
		ShowSearch:           false,
		ShowActions:          false,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "start",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   emptyTitle,
			Message: emptyMsg,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}
