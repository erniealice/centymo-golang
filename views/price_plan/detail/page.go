package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// DetailViewDeps holds view dependencies for the price plan detail page.
type DetailViewDeps struct {
	Routes                 centymo.PricePlanRoutes
	Labels                 centymo.PricePlanLabels
	ProductPricePlanLabels centymo.ProductPricePlanLabels
	CommonLabels           pyeza.CommonLabels
	TableLabels            types.TableLabels

	ReadPricePlan          func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	ListProductPlans       func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProducts           func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ListProductVariants    func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	// ListProductOptions / ListProductOptionValues / ListProductVariantOptions
	// power the enriched variant label in the catalog-line picker
	// ("Product — SKU — Red / Large / Cotton"). Optional — when nil the
	// label falls back to the plain "Product — SKU" form.
	ListProductOptions        func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ListProductOptionValues   func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	ListProductVariantOptions func(ctx context.Context, req *productvariantoptionpb.ListProductVariantOptionsRequest) (*productvariantoptionpb.ListProductVariantOptionsResponse, error)
	ListProductPricePlans  func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	CreateProductPricePlan func(ctx context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)
	UpdateProductPricePlan func(ctx context.Context, req *productpriceplanpb.UpdateProductPricePlanRequest) (*productpriceplanpb.UpdateProductPricePlanResponse, error)
	DeleteProductPricePlan func(ctx context.Context, req *productpriceplanpb.DeleteProductPricePlanRequest) (*productpriceplanpb.DeleteProductPricePlanResponse, error)
}

// ProductPlanGroup groups catalog-line options by parent product so the
// drawer's <select> can render <optgroup label="Product Name">…</optgroup>
// blocks. Within each group, Options is sorted alphabetically by Label.
type ProductPlanGroup struct {
	ProductName string
	Options     []types.SelectOption
}

// PageData holds the data for the price plan detail page.
type PageData struct {
	types.PageData
	ContentTemplate      string
	PricePlan            *priceplanpb.PricePlan
	Labels               centymo.PricePlanLabels
	ActiveTab            string
	TabItems             []pyeza.TabItem
	ID                   string
	PricePlanName        string
	PricePlanDesc        string
	PricePlanAmount      types.TableCell
	PricePlanCurrency    string
	PricePlanLocation    string
	PricePlanDuration    string
	PricePlanStatus      string
	StatusVariant        string
	CreatedDate          string
	ModifiedDate         string
	ProductPricesTable   *types.TableConfig
}

// ProductPricePlanFormData holds data for the add/edit drawer form.
//
// Model D — the drawer now references a catalog line (ProductPlan) rather
// than a bare product. ProductPlanID is the FK written to ProductPricePlan;
// ProductPlanOptions drives the picker and is scoped to the PricePlan's
// parent Plan. The read-only SelectedProductName + SelectedVariantName
// surface the resolved product + variant context above the price input.
type ProductPricePlanFormData struct {
	FormAction         string
	IsEdit             bool
	ID                 string
	PricePlanID        string
	ProductPlanID      string
	// Grouped by parent product — each ProductPlanGroup renders as an
	// <optgroup>. Within each group, Options is alphabetically sorted; the
	// outer slice is sorted by ProductName.
	ProductPlanOptions []ProductPlanGroup
	// Read-only display of the selected catalog line's product + variant
	// (surfaced above the price input).
	SelectedProductName string
	SelectedVariantName string
	Price               string
	Currency            string
	CommonLabels        pyeza.CommonLabels

	// Wave 2: billing treatment + effective date fields.
	BillingTreatment        string
	BillingTreatmentOptions []types.SelectOption
	DateStart               string // ISO 8601 (YYYY-MM-DD) or empty
	DateEnd                 string // ISO 8601 (YYYY-MM-DD) or empty

	// Wave 2: labels for the new fields (populated from ProductPricePlanLabels).
	Labels centymo.ProductPricePlanFormLabels
}

// NewView creates the price plan detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
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

		return view.OK("price-plan-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
// Handles GET /action/price-plans/{id}/tab/{tab}
func NewTabAction(deps *DetailViewDeps) view.View {
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

		templateName := "price-plan-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

// NewProductPriceAddAction handles GET (render form) and POST (submit) for adding a ProductPricePlan.
func NewProductPriceAddAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_price_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if deps.CreateProductPricePlan == nil {
			return centymo.HTMXError("Product price plan create is not available")
		}

		if viewCtx.Request.Method == http.MethodGet {
			// Model D — load catalog-line options scoped to the PricePlan's parent Plan.
			planID := loadPricePlanPlanID(ctx, deps, id)
			productPlanOptions := loadProductPlanOptions(ctx, deps, planID, id, "")
			currency := loadPricePlanCurrency(ctx, deps, id)
			pplLabels := deps.ProductPricePlanLabels.Form
			return view.OK("product-price-plan-drawer-form", &ProductPricePlanFormData{
				FormAction:              route.ResolveURL(deps.Routes.ProductPriceAddURL, "id", id),
				PricePlanID:             id,
				Currency:                currency,
				ProductPlanOptions:      productPlanOptions,
				CommonLabels:            deps.CommonLabels,
				BillingTreatmentOptions: buildBillingTreatmentOptions(pplLabels),
				Labels:                  pplLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		productPlanID := viewCtx.Request.FormValue("product_plan_id")
		if productPlanID == "" {
			return centymo.HTMXError("Catalog line is required")
		}
		priceStr := viewCtx.Request.FormValue("price")
		currency := viewCtx.Request.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}

		priceFloat, err := strconv.ParseFloat(priceStr, 64)
		if err != nil || priceFloat < 0 {
			return centymo.HTMXError("Invalid price value")
		}
		priceCentavos := int64(priceFloat * 100)

		dateStart := viewCtx.Request.FormValue("date_start")
		dateEnd := viewCtx.Request.FormValue("date_end")
		billingTreatment := viewCtx.Request.FormValue("billing_treatment")

		record := &productpriceplanpb.ProductPricePlan{
			PricePlanId:   id,
			ProductPlanId: productPlanID,
			BillingAmount: priceCentavos,
			BillingCurrency: currency,
			Active:        true,
		}
		if billingTreatment != "" {
			if bt, ok := productpriceplanpb.BillingTreatment_value[billingTreatment]; ok {
				record.BillingTreatment = productpriceplanpb.BillingTreatment(bt)
			}
		}
		if dateStart != "" {
			record.DateStart = &dateStart
		}
		if dateEnd != "" {
			record.DateEnd = &dateEnd
		}

		if _, err := deps.CreateProductPricePlan(ctx, &productpriceplanpb.CreateProductPricePlanRequest{Data: record}); err != nil {
			log.Printf("Failed to create product price plan for price plan %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("price-plan-product-prices-table")
	})
}

// NewProductPriceEditAction handles GET (render pre-filled form) and POST (submit) for editing a ProductPricePlan.
func NewProductPriceEditAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")

		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_price_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if deps.UpdateProductPricePlan == nil {
			return centymo.HTMXError("Product price plan update is not available")
		}

		existing, err := findProductPricePlan(ctx, deps, ppid)
		if err != nil {
			return centymo.HTMXError(err.Error())
		}

		pplLabels := deps.ProductPricePlanLabels.Form

		if viewCtx.Request.Method == http.MethodGet {
			planID := loadPricePlanPlanID(ctx, deps, id)
			existingProductPlanID := existing.GetProductPlanId()
			productPlanOptions := loadProductPlanOptions(ctx, deps, planID, id, existingProductPlanID)
			productName, variantName := resolveProductPlanDisplay(ctx, deps, existingProductPlanID)
			currency := existing.GetBillingCurrency()
			if currency == "" {
				currency = "PHP"
			}
			return view.OK("product-price-plan-drawer-form", &ProductPricePlanFormData{
				FormAction:              route.ResolveURL(deps.Routes.ProductPriceEditURL, "id", id, "ppid", ppid),
				IsEdit:                  true,
				ID:                      ppid,
				PricePlanID:             id,
				ProductPlanID:           existingProductPlanID,
				ProductPlanOptions:      productPlanOptions,
				SelectedProductName:     productName,
				SelectedVariantName:     variantName,
				Price:                   fmt.Sprintf("%.2f", float64(existing.GetBillingAmount())/100.0),
				Currency:                currency,
				CommonLabels:            deps.CommonLabels,
				BillingTreatment:        existing.GetBillingTreatment().String(),
				BillingTreatmentOptions: buildBillingTreatmentOptions(pplLabels),
				DateStart:               existing.GetDateStart(),
				DateEnd:                 existing.GetDateEnd(),
				Labels:                  pplLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		// On edit, the catalog line is pinned (the drawer disables the picker);
		// still accept a posted product_plan_id when present, else preserve.
		productPlanID := viewCtx.Request.FormValue("product_plan_id")
		if productPlanID == "" {
			productPlanID = existing.GetProductPlanId()
		}
		if productPlanID == "" {
			return centymo.HTMXError("Catalog line is required")
		}
		priceStr := viewCtx.Request.FormValue("price")
		currency := viewCtx.Request.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}

		priceFloat, err := strconv.ParseFloat(priceStr, 64)
		if err != nil || priceFloat < 0 {
			return centymo.HTMXError("Invalid price value")
		}
		priceCentavos := int64(priceFloat * 100)

		dateStart := viewCtx.Request.FormValue("date_start")
		dateEnd := viewCtx.Request.FormValue("date_end")
		billingTreatment := viewCtx.Request.FormValue("billing_treatment")

		updated := &productpriceplanpb.ProductPricePlan{
			Id:              ppid,
			PricePlanId:     id,
			ProductPlanId:   productPlanID,
			BillingAmount:   priceCentavos,
			BillingCurrency: currency,
			Active:          existing.GetActive(),
		}
		if billingTreatment != "" {
			if bt, ok := productpriceplanpb.BillingTreatment_value[billingTreatment]; ok {
				updated.BillingTreatment = productpriceplanpb.BillingTreatment(bt)
			}
		}
		if dateStart != "" {
			updated.DateStart = &dateStart
		}
		if dateEnd != "" {
			updated.DateEnd = &dateEnd
		}

		if _, err := deps.UpdateProductPricePlan(ctx, &productpriceplanpb.UpdateProductPricePlanRequest{Data: updated}); err != nil {
			log.Printf("Failed to update product price plan %s: %v", ppid, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("price-plan-product-prices-table")
	})
}

// NewProductPriceDeleteAction handles POST for deleting a ProductPricePlan.
func NewProductPriceDeleteAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_price_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if deps.DeleteProductPricePlan == nil {
			return centymo.HTMXError("Product price plan delete is not available")
		}

		_ = viewCtx.Request.ParseForm()
		ppid := viewCtx.Request.FormValue("id")
		if ppid == "" {
			ppid = viewCtx.Request.URL.Query().Get("id")
		}
		if ppid == "" {
			return centymo.HTMXError("ID is required")
		}

		if _, err := deps.DeleteProductPricePlan(ctx, &productpriceplanpb.DeleteProductPricePlanRequest{
			Data: &productpriceplanpb.ProductPricePlan{Id: ppid},
		}); err != nil {
			log.Printf("Failed to delete product price plan %s for price plan %s: %v", ppid, id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("price-plan-product-prices-table")
	})
}

// ---------------------------------------------------------------------------
// buildPageData
// ---------------------------------------------------------------------------

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read price plan %s: %v", id, err)
		return nil, fmt.Errorf("failed to load rate card: %w", err)
	}

	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("rate card not found")
	}
	pp := data[0]

	name := pp.GetName()
	description := pp.GetDescription()
	currency := pp.GetBillingCurrency()
	if currency == "" {
		currency = "PHP"
	}
	amountFormatted := types.MoneyCell(float64(pp.GetBillingAmount()), currency, true)

	duration := ""
	if dv := pp.GetDurationValue(); dv > 0 {
		duration = pyeza.FormatDuration(dv, pp.GetDurationUnit(), deps.CommonLabels.DurationUnit)
	}

	status := "active"
	if !pp.GetActive() {
		status = "inactive"
	}
	statusVariant := "success"
	if status == "inactive" {
		statusVariant = "warning"
	}

	l := deps.Labels

	// Count for product-prices tab badge
	productPriceCount := 0
	if deps.ListProductPricePlans != nil {
		pppResp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{})
		if err == nil {
			for _, item := range pppResp.GetData() {
				if item.GetPricePlanId() == id {
					productPriceCount++
				}
			}
		}
	}

	tabItems := buildTabItems(id, l, productPriceCount, deps.Routes)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          name,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    name,
			HeaderSubtitle: description,
			HeaderIcon:     "icon-tag",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate:   "price-plan-detail-content",
		PricePlan:         pp,
		Labels:            l,
		ActiveTab:         activeTab,
		TabItems:          tabItems,
		ID:                id,
		PricePlanName:     name,
		PricePlanDesc:     description,
		PricePlanAmount:   amountFormatted,
		PricePlanCurrency: currency,
		PricePlanLocation: pp.GetPriceScheduleId(),
		PricePlanDuration: duration,
		PricePlanStatus:   status,
		StatusVariant:     statusVariant,
		CreatedDate:       pp.GetDateCreatedString(),
		ModifiedDate:      pp.GetDateModifiedString(),
	}

	// Load tab-specific data
	switch activeTab {
	case "product-prices":
		tableConfig := buildProductPricesTable(ctx, deps, id, pp.GetPlanId())
		pageData.ProductPricesTable = tableConfig
	}

	return pageData, nil
}

func buildTabItems(id string, l centymo.PricePlanLabels, productPriceCount int, routes centymo.PricePlanRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "product-prices", Label: l.Tabs.Products, Href: base + "?tab=product-prices", HxGet: action + "product-prices", Icon: "icon-package", Count: productPriceCount},
	}
}

// ---------------------------------------------------------------------------
// Product prices tab table
// ---------------------------------------------------------------------------

func buildProductPricesTable(ctx context.Context, deps *DetailViewDeps, pricePlanID, planID string) *types.TableConfig {
	l := deps.Labels
	perms := view.GetUserPermissions(ctx)

	columns := []types.TableColumn{
		{Key: "product", Label: "Product", Sortable: true},
		{Key: "billing_treatment", Label: "Billing", Sortable: true, WidthClass: "col-3xl"},
		{Key: "price", Label: "Price", Sortable: true, WidthClass: "col-4xl"},
	}

	// Build product ID → name map for display
	productNames := map[string]string{}
	if deps.ListProducts != nil {
		prodResp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
		if err == nil {
			for _, p := range prodResp.GetData() {
				if p != nil {
					productNames[p.GetId()] = p.GetName()
				}
			}
		}
	}

	// Model D — build product_plan_id → (product_id, variant_id) map so we
	// can display "Product (SKU)" for rows that reference a catalog line
	// carrying a variant.
	type productPlanRef struct {
		productID string
		variantID string
	}
	productPlans := map[string]productPlanRef{}
	if deps.ListProductPlans != nil {
		ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
		if err == nil {
			for _, pp := range ppResp.GetData() {
				if pp == nil {
					continue
				}
				productPlans[pp.GetId()] = productPlanRef{
					productID: pp.GetProductId(),
					variantID: pp.GetProductVariantId(),
				}
			}
		}
	}

	// Build variant_id → SKU map (fall back to id).
	variantLabels := map[string]string{}
	if deps.ListProductVariants != nil {
		vResp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{})
		if err == nil {
			for _, v := range vResp.GetData() {
				if v == nil {
					continue
				}
				sku := v.GetSku()
				if sku == "" {
					sku = v.GetId()
				}
				variantLabels[v.GetId()] = sku
			}
		}
	}

	billingLabels := map[productpriceplanpb.BillingTreatment]string{
		productpriceplanpb.BillingTreatment_BILLING_TREATMENT_RECURRING:          deps.ProductPricePlanLabels.Form.BillingTreatmentRecurring,
		productpriceplanpb.BillingTreatment_BILLING_TREATMENT_ONE_TIME_INITIAL:   deps.ProductPricePlanLabels.Form.BillingTreatmentOneTimeInitial,
		productpriceplanpb.BillingTreatment_BILLING_TREATMENT_USAGE_BASED:        deps.ProductPricePlanLabels.Form.BillingTreatmentUsageBased,
	}

	rows := []types.TableRow{}
	if deps.ListProductPricePlans != nil {
		pppResp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{})
		if err != nil {
			log.Printf("Failed to list product price plans: %v", err)
		} else {
			for _, item := range pppResp.GetData() {
				if item == nil || item.GetPricePlanId() != pricePlanID {
					continue
				}

				itemID := item.GetId()
				// Model D — resolve product + variant from the referenced
				// ProductPlan row. When the adapter eventually populates the
				// embedded ProductPlan we can read straight from
				// item.GetProductPlan(); for now we fall back to the local map.
				ppID := item.GetProductPlanId()
				ref := productPlans[ppID]
				if embed := item.GetProductPlan(); embed != nil {
					if pid := embed.GetProductId(); pid != "" {
						ref.productID = pid
					}
					if vid := embed.GetProductVariantId(); vid != "" {
						ref.variantID = vid
					}
				}
				productName := productNames[ref.productID]
				if productName == "" {
					productName = ref.productID
				}
				if ref.variantID != "" {
					if label := variantLabels[ref.variantID]; label != "" {
						productName = fmt.Sprintf("%s (%s)", productName, label)
					}
				}
				itemCurrency := item.GetBillingCurrency()
				if itemCurrency == "" {
					itemCurrency = "PHP"
				}
				priceCell := types.MoneyCell(float64(item.GetBillingAmount()), itemCurrency, true)
				btLabel := billingLabels[item.GetBillingTreatment()]
				if btLabel == "" {
					btLabel = item.GetBillingTreatment().String()
				}

				rows = append(rows, types.TableRow{
					ID: itemID,
					Cells: []types.TableCell{
						{Type: "text", Value: productName},
						{Type: "text", Value: btLabel},
						priceCell,
					},
					Actions: []types.TableAction{
						{
							Type:            "edit",
							Label:           "Edit",
							Action:          "edit",
							URL:             route.ResolveURL(deps.Routes.ProductPriceEditURL, "id", pricePlanID, "ppid", itemID),
							DrawerTitle:     l.ProductPrice.EditTitle,
							Disabled:        !perms.Can("product_price_plan", "update"),
							DisabledTooltip: l.Errors.Unauthorized,
						},
						{
							Type:            "delete",
							Label:           "Delete",
							Action:          "delete",
							URL:             deps.Routes.ProductPriceDeleteURL,
							ItemName:        productName,
							ConfirmTitle:    l.ProductPrice.DeleteTitle,
							ConfirmMessage:  fmt.Sprintf("Remove %s from this rate card?", productName),
							Disabled:        !perms.Can("product_price_plan", "delete"),
							DisabledTooltip: l.Errors.Unauthorized,
						},
					},
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "price-plan-product-prices-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "product",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.ProductPrice.EmptyTitle,
			Message: l.ProductPrice.EmptyMsg,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           "Add Product Price",
			ActionURL:       route.ResolveURL(deps.Routes.ProductPriceAddURL, "id", pricePlanID),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("product_price_plan", "create"),
			DisabledTooltip: l.Errors.Unauthorized,
		},
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func findProductPricePlan(ctx context.Context, deps *DetailViewDeps, ppid string) (*productpriceplanpb.ProductPricePlan, error) {
	if deps.ListProductPricePlans == nil {
		return nil, fmt.Errorf("product price plans not available")
	}
	resp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to load product price plans")
	}
	for _, item := range resp.GetData() {
		if item != nil && item.GetId() == ppid {
			return item, nil
		}
	}
	return nil, fmt.Errorf("product price plan not found")
}

// loadPricePlanPlanID reads the price plan to get its linked plan_id.
func loadPricePlanPlanID(ctx context.Context, deps *DetailViewDeps, pricePlanID string) string {
	if deps.ReadPricePlan == nil {
		return ""
	}
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: pricePlanID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return ""
	}
	return resp.GetData()[0].GetPlanId()
}

// loadPricePlanCurrency reads the price plan to get its currency.
func loadPricePlanCurrency(ctx context.Context, deps *DetailViewDeps, pricePlanID string) string {
	if deps.ReadPricePlan == nil {
		return "PHP"
	}
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: pricePlanID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return "PHP"
	}
	c := resp.GetData()[0].GetBillingCurrency()
	if c == "" {
		return "PHP"
	}
	return c
}

// loadProductPlanOptions builds a grouped select list of catalog lines
// (ProductPlan rows) belonging to the given Plan. This is the Model D
// replacement for the bare product picker — users now pick *which catalog
// line to price*, not a bare product, so the selector is intrinsically
// scoped to the plan.
//
// The returned []ProductPlanGroup is grouped by parent product (one
// <optgroup> per product). Within each group, Options is sorted
// alphabetically by Label (case-insensitive); the outer slice is sorted by
// ProductName.
//
// Label format:
//   - Product-only line: "Product Name"
//   - Variant line with option values: "Product Name — SKU — Red / Large / Cotton"
//     (option values ordered by product_option.sort_order ASC,
//     joined by " / ")
//   - Variant line without option-value rows: "Product Name — SKU"
//
// Lines that already have a ProductPricePlan row in the current PricePlan
// (pricePlanID) are excluded so the user cannot double-price the same
// catalog line. The currently-selected line (selectedProductPlanID, set in
// edit mode) is always included even if it is "already priced" — otherwise
// the edit drawer would render with an empty selection.
func loadProductPlanOptions(ctx context.Context, deps *DetailViewDeps, planID, pricePlanID, selectedProductPlanID string) []ProductPlanGroup {
	if deps.ListProductPlans == nil || planID == "" {
		return nil
	}
	// Build product ID → name map.
	productNames := map[string]string{}
	if deps.ListProducts != nil {
		prodResp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
		if err == nil {
			for _, p := range prodResp.GetData() {
				if p != nil {
					productNames[p.GetId()] = p.GetName()
				}
			}
		}
	}
	// Build variant ID → SKU + active map. We track active state alongside SKU
	// so the picker can filter out catalog lines pinned to deactivated variants
	// (with an edit-mode exception so the previously-selected line round-trips).
	variantSKUs := map[string]string{}
	variantActive := map[string]bool{}
	if deps.ListProductVariants != nil {
		vResp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{})
		if err == nil {
			for _, v := range vResp.GetData() {
				if v == nil {
					continue
				}
				sku := v.GetSku()
				if sku == "" {
					sku = v.GetId()
				}
				variantSKUs[v.GetId()] = sku
				variantActive[v.GetId()] = v.GetActive()
			}
		}
	}

	// Build variant ID → []labelInOptionSortOrder. Mirrors the variantOptionLabels
	// map in product/detail/page.go BuildVariantsTable. Two lookups are needed:
	//   1. product_option_value.id → (label, product_option_id)
	//   2. product_option.id → sort_order
	// Then for each product_variant_option row, place the value's label at the
	// position dictated by its parent option's sort_order.
	variantOptionLabels := map[string][]string{}
	if deps.ListProductVariantOptions != nil && deps.ListProductOptionValues != nil {
		// product_option_id → sort_order
		optionSortOrder := map[string]int32{}
		if deps.ListProductOptions != nil {
			if optResp, err := deps.ListProductOptions(ctx, &productoptionpb.ListProductOptionsRequest{}); err == nil {
				for _, opt := range optResp.GetData() {
					if opt != nil {
						optionSortOrder[opt.GetId()] = opt.GetSortOrder()
					}
				}
			}
		}
		// product_option_value.id → (label, product_option_id)
		type valueRef struct {
			label    string
			optionID string
		}
		valueByID := map[string]valueRef{}
		if ovResp, err := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{}); err == nil {
			for _, ov := range ovResp.GetData() {
				if ov != nil {
					valueByID[ov.GetId()] = valueRef{
						label:    ov.GetLabel(),
						optionID: ov.GetProductOptionId(),
					}
				}
			}
		}
		// Per-variant tuples (sortOrder, label) so we can sort, then project to []string.
		type ordered struct {
			sortOrder int32
			label     string
		}
		acc := map[string][]ordered{}
		if voResp, err := deps.ListProductVariantOptions(ctx, &productvariantoptionpb.ListProductVariantOptionsRequest{}); err == nil {
			for _, vo := range voResp.GetData() {
				if vo == nil {
					continue
				}
				vid := vo.GetProductVariantId()
				ref, ok := valueByID[vo.GetProductOptionValueId()]
				if !ok || vid == "" {
					continue
				}
				acc[vid] = append(acc[vid], ordered{
					sortOrder: optionSortOrder[ref.optionID],
					label:     ref.label,
				})
			}
		}
		for vid, list := range acc {
			sort.SliceStable(list, func(i, j int) bool {
				return list[i].sortOrder < list[j].sortOrder
			})
			labels := make([]string, 0, len(list))
			for _, o := range list {
				labels = append(labels, o.label)
			}
			variantOptionLabels[vid] = labels
		}
	}

	// Build the set of ProductPlan IDs already priced under this PricePlan,
	// so we can exclude them from the picker. Edit mode preserves the
	// currently-selected line via the selectedProductPlanID exception.
	alreadyPriced := map[string]bool{}
	if deps.ListProductPricePlans != nil && pricePlanID != "" {
		if pppResp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{}); err == nil {
			for _, item := range pppResp.GetData() {
				if item == nil || item.GetPricePlanId() != pricePlanID {
					continue
				}
				if ppid := item.GetProductPlanId(); ppid != "" {
					alreadyPriced[ppid] = true
				}
			}
		}
	}

	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
	if err != nil {
		log.Printf("Failed to list product plans for plan %s: %v", planID, err)
		return nil
	}

	// Group options by parent product name. We keep the underlying productID
	// in a parallel map so blank product names from missing lookups don't
	// collide across different products.
	byProduct := map[string]*ProductPlanGroup{}
	for _, pp := range ppResp.GetData() {
		if pp == nil || pp.GetPlanId() != planID {
			continue
		}
		ppID := pp.GetId()
		isSelected := ppID == selectedProductPlanID
		// Filter already-priced lines unless this is the line the user is
		// currently editing.
		if alreadyPriced[ppID] && !isSelected {
			continue
		}
		// Filter inactive catalog lines (and lines pinned to a deactivated
		// variant) so the picker only surfaces things you can actually price
		// today. Edit mode keeps the selected line visible regardless so the
		// form round-trips a known-good FK.
		if !pp.GetActive() && !isSelected {
			continue
		}
		if vid := pp.GetProductVariantId(); vid != "" {
			if active, known := variantActive[vid]; known && !active && !isSelected {
				continue
			}
		}

		pid := pp.GetProductId()
		productName := productNames[pid]
		if productName == "" {
			productName = pp.GetName()
		}
		if productName == "" {
			productName = pid
		}

		// Build the option label per the rules above.
		label := productName
		if vid := pp.GetProductVariantId(); vid != "" {
			sku := variantSKUs[vid]
			if sku == "" {
				sku = vid
			}
			label = fmt.Sprintf("%s — %s", productName, sku)
			if values := variantOptionLabels[vid]; len(values) > 0 {
				label = fmt.Sprintf("%s — %s", label, strings.Join(values, centymo.OptionValueSeparator))
			}
		}

		groupKey := pid
		if groupKey == "" {
			groupKey = productName
		}
		group, ok := byProduct[groupKey]
		if !ok {
			group = &ProductPlanGroup{ProductName: productName}
			byProduct[groupKey] = group
		}
		group.Options = append(group.Options, types.SelectOption{
			Value:    ppID,
			Label:    label,
			Selected: ppID == selectedProductPlanID,
		})
	}

	// Materialise + sort.
	groups := make([]ProductPlanGroup, 0, len(byProduct))
	for _, g := range byProduct {
		sort.SliceStable(g.Options, func(i, j int) bool {
			return strings.ToLower(g.Options[i].Label) < strings.ToLower(g.Options[j].Label)
		})
		groups = append(groups, *g)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		return strings.ToLower(groups[i].ProductName) < strings.ToLower(groups[j].ProductName)
	})
	return groups
}

// resolveProductPlanDisplay returns the product name + variant SKU (if any)
// for the given ProductPlan.id — surfaced read-only in the price drawer's
// "selected catalog line" context row.
func resolveProductPlanDisplay(ctx context.Context, deps *DetailViewDeps, productPlanID string) (productName, variantName string) {
	if productPlanID == "" || deps.ListProductPlans == nil {
		return "", ""
	}
	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
	if err != nil {
		return "", ""
	}
	var (
		pid string
		vid string
	)
	for _, pp := range ppResp.GetData() {
		if pp != nil && pp.GetId() == productPlanID {
			pid = pp.GetProductId()
			vid = pp.GetProductVariantId()
			break
		}
	}
	if pid != "" && deps.ListProducts != nil {
		if prodResp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{}); err == nil {
			for _, p := range prodResp.GetData() {
				if p != nil && p.GetId() == pid {
					productName = p.GetName()
					break
				}
			}
		}
	}
	if vid != "" && deps.ListProductVariants != nil {
		if vResp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{}); err == nil {
			for _, v := range vResp.GetData() {
				if v != nil && v.GetId() == vid {
					variantName = v.GetSku()
					if variantName == "" {
						variantName = v.GetId()
					}
					break
				}
			}
		}
	}
	return productName, variantName
}

// buildBillingTreatmentOptions builds the select options for the BillingTreatment
// enum using lyngua-provided labels. Values are proto enum string names.
func buildBillingTreatmentOptions(labels centymo.ProductPricePlanFormLabels) []types.SelectOption {
	return []types.SelectOption{
		{Value: "BILLING_TREATMENT_RECURRING", Label: labels.BillingTreatmentRecurring},
		{Value: "BILLING_TREATMENT_ONE_TIME_INITIAL", Label: labels.BillingTreatmentOneTimeInitial},
		{Value: "BILLING_TREATMENT_USAGE_BASED", Label: labels.BillingTreatmentUsageBased},
	}
}

