package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// DetailViewDeps holds view dependencies for the price plan detail page.
type DetailViewDeps struct {
	Routes                  centymo.PricePlanRoutes
	Labels                  centymo.PricePlanLabels
	ProductPricePlanLabels  centymo.ProductPricePlanLabels
	CommonLabels            pyeza.CommonLabels
	TableLabels             types.TableLabels

	ReadPricePlan            func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	ListProductPlans         func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProducts             func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ListProductPricePlans    func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	CreateProductPricePlan   func(ctx context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)
	UpdateProductPricePlan   func(ctx context.Context, req *productpriceplanpb.UpdateProductPricePlanRequest) (*productpriceplanpb.UpdateProductPricePlanResponse, error)
	DeleteProductPricePlan   func(ctx context.Context, req *productpriceplanpb.DeleteProductPricePlanRequest) (*productpriceplanpb.DeleteProductPricePlanResponse, error)
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
type ProductPricePlanFormData struct {
	FormAction     string
	IsEdit         bool
	ID             string
	PricePlanID    string
	ProductID      string
	Price          string
	Currency       string
	ProductOptions []types.SelectOption
	CommonLabels   pyeza.CommonLabels

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
			// Load the linked plan's products for the selector
			planID := loadPricePlanPlanID(ctx, deps, id)
			productOptions := loadProductOptions(ctx, deps, planID, "")
			currency := loadPricePlanCurrency(ctx, deps, id)
			pplLabels := deps.ProductPricePlanLabels.Form
			return view.OK("product-price-plan-drawer-form", &ProductPricePlanFormData{
				FormAction:              route.ResolveURL(deps.Routes.ProductPriceAddURL, "id", id),
				PricePlanID:             id,
				Currency:                currency,
				ProductOptions:          productOptions,
				CommonLabels:            deps.CommonLabels,
				BillingTreatmentOptions: buildBillingTreatmentOptions(pplLabels),
				Labels:                  pplLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		productID := viewCtx.Request.FormValue("product_id")
		if productID == "" {
			return centymo.HTMXError("Product is required")
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
			PricePlanId: id,
			ProductId:   productID,
			Price:       priceCentavos,
			Currency:    currency,
			Active:      true,
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
			productOptions := loadProductOptions(ctx, deps, planID, existing.GetProductId())
			currency := existing.GetCurrency()
			if currency == "" {
				currency = "PHP"
			}
			return view.OK("product-price-plan-drawer-form", &ProductPricePlanFormData{
				FormAction:              route.ResolveURL(deps.Routes.ProductPriceEditURL, "id", id, "ppid", ppid),
				IsEdit:                  true,
				ID:                      ppid,
				PricePlanID:             id,
				ProductID:               existing.GetProductId(),
				Price:                   fmt.Sprintf("%.2f", float64(existing.GetPrice())/100.0),
				Currency:                currency,
				ProductOptions:          productOptions,
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

		productID := viewCtx.Request.FormValue("product_id")
		if productID == "" {
			return centymo.HTMXError("Product is required")
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
			Id:          ppid,
			PricePlanId: id,
			ProductId:   productID,
			Price:       priceCentavos,
			Currency:    currency,
			Active:      existing.GetActive(),
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
	currency := pp.GetCurrency()
	if currency == "" {
		currency = "PHP"
	}
	amountFormatted := types.MoneyCell(float64(pp.GetAmount()), currency, true)

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
				productID := item.GetProductId()
				productName := productNames[productID]
				if productName == "" {
					productName = productID
				}
				itemCurrency := item.GetCurrency()
				if itemCurrency == "" {
					itemCurrency = "PHP"
				}
				priceCell := types.MoneyCell(float64(item.GetPrice()), itemCurrency, true)

				rows = append(rows, types.TableRow{
					ID: itemID,
					Cells: []types.TableCell{
						{Type: "text", Value: productName},
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
	c := resp.GetData()[0].GetCurrency()
	if c == "" {
		return "PHP"
	}
	return c
}

// loadProductOptions builds a select option list from ProductPlans of a given plan.
// It only shows products assigned to this plan (via product_plan) so the selector
// is scoped to the correct product set.
func loadProductOptions(ctx context.Context, deps *DetailViewDeps, planID, selectedProductID string) []types.SelectOption {
	// Build product ID → name map
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

	// List ProductPlans for this plan to scope the options
	if deps.ListProductPlans == nil || planID == "" {
		// Fallback: show all products if we can't scope
		options := make([]types.SelectOption, 0, len(productNames))
		for pid, name := range productNames {
			options = append(options, types.SelectOption{
				Value:    pid,
				Label:    name,
				Selected: pid == selectedProductID,
			})
		}
		return options
	}

	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
	if err != nil {
		log.Printf("Failed to list product plans for plan %s: %v", planID, err)
		return nil
	}

	options := []types.SelectOption{}
	for _, pp := range ppResp.GetData() {
		if pp == nil || pp.GetPlanId() != planID {
			continue
		}
		pid := pp.GetProductId()
		name := productNames[pid]
		if name == "" {
			name = pid
		}
		options = append(options, types.SelectOption{
			Value:    pid,
			Label:    name,
			Selected: pid == selectedProductID,
		})
	}
	return options
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

