// Package plan renders the price_plan detail page nested under its parent PriceSchedule
// at /app/price-schedules/detail/{id}/plan/{ppid}. The sidebar stays on price-schedules
// because price_plan is no longer a top-level sidebar entry.
package plan

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// DetailViewDeps holds all dependencies for the schedule-scoped price_plan detail page.
type DetailViewDeps struct {
	Routes         centymo.PriceScheduleRoutes
	ScheduleLabels centymo.PriceScheduleLabels
	PlanLabels     centymo.PricePlanLabels
	CommonLabels   pyeza.CommonLabels
	TableLabels    types.TableLabels

	ReadPriceSchedule func(ctx context.Context, req *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error)
	ReadPricePlan     func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	UpdatePricePlan   func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan   func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)

	ListPlans        func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListProducts     func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ListProductPlans func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)

	ListProductPricePlans  func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	CreateProductPricePlan func(ctx context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)
	UpdateProductPricePlan func(ctx context.Context, req *productpriceplanpb.UpdateProductPricePlanRequest) (*productpriceplanpb.UpdateProductPricePlanResponse, error)
	DeleteProductPricePlan func(ctx context.Context, req *productpriceplanpb.DeleteProductPricePlanRequest) (*productpriceplanpb.DeleteProductPricePlanResponse, error)

	// Reference checker: returns a map of price_plan_id → true for plans in use by active subscriptions.
	// When a plan is in use, Pricing fields in the Edit drawer are read-only.
	GetPricePlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// PageData is the template data for the schedule-scoped plan detail page.
type PageData struct {
	types.PageData
	ContentTemplate string

	ScheduleID      string
	ScheduleName    string
	ScheduleBackURL string
	PricePlan       *priceplanpb.PricePlan
	Labels          centymo.PricePlanLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID            string
	Name          string
	Description   string
	Amount        types.TableCell
	Currency      string
	Duration      string
	Status        string
	StatusVariant string
	CreatedDate   string
	ModifiedDate  string

	EditURL                string
	ProductPricesTable     *types.TableConfig
	ProductPriceEmptyTitle string
	ProductPriceEmptyMsg   string
}

// EditFormData is the drawer form for editing a price_plan under a schedule.
type EditFormData struct {
	FormAction    string
	ScheduleID    string
	ScheduleName  string
	ID            string
	PlanID        string
	PlanLabel     string // display label for the currently-selected plan (for SelectedLabel on auto-complete)
	PlanOptions   []map[string]any
	Name          string
	Description   string
	Amount        string
	Currency      string
	DurationValue string
	DurationUnit  string
	CommonLabels  pyeza.CommonLabels

	// PricingLocked is true when the price_plan is referenced by active subscriptions.
	// The Pricing section fields (Amount, Currency, Duration, DurationUnit) are rendered
	// as read-only in the drawer, but all other fields remain editable.
	PricingLocked       bool
	PricingLockedReason string
}

// ProductPriceFormData is the drawer form for adding/editing a ProductPricePlan.
// SectionTitle + ProductFieldLabel pull from PriceScheduleDetailLabels so the
// drawer reads "Service Price" / "Service" in the professional tier.
type ProductPriceFormData struct {
	FormAction        string
	IsEdit            bool
	ID                string
	ScheduleID        string
	PricePlanID       string
	ProductID         string
	Price             string
	Currency          string
	ProductOptions    []types.SelectOption
	SectionTitle      string
	ProductFieldLabel string
	CommonLabels      pyeza.CommonLabels
}

// NewView renders the full detail page at /app/price-schedules/detail/{id}/plan/{ppid}.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, sid, ppid, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("price-schedule-plan-detail", pageData)
	})
}

// NewTabAction handles GET /action/price-schedule/{id}/plan/{ppid}/tab/{tab}.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		pageData, err := buildPageData(ctx, deps, sid, ppid, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("price-schedule-plan-tab-"+tab, pageData)
	})
}

// NewEditAction handles GET/POST /action/price-schedule/{id}/plan/{ppid}/edit.
func NewEditAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "update") {
			return centymo.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: ppid}})
			if err != nil || len(resp.GetData()) == 0 {
				return centymo.HTMXError(deps.PlanLabels.Errors.NotFound)
			}
			pp := resp.GetData()[0]

			// Check whether this plan is referenced by active subscriptions.
			// When true, the Pricing section is rendered read-only in the drawer.
			pricingLocked := false
			pricingLockedReason := ""
			if deps.GetPricePlanInUseIDs != nil {
				inUseMap, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppid})
				if inUseMap[ppid] {
					pricingLocked = true
					// TODO: lyngua — pull from deps.ScheduleLabels.Detail.PricingLockedReason if added
					pricingLockedReason = "This plan is in use by active subscriptions. Pricing changes are disabled. You can still rename or reassign the package."
				}
			}

			planOpts := buildPlanOptions(ctx, deps, pp.GetPlanId())
			return view.OK("price-schedule-plan-edit-drawer", &EditFormData{
				FormAction:          route.ResolveURL(deps.Routes.PlanEditURL, "id", sid, "ppid", ppid),
				ScheduleID:          sid,
				ScheduleName:        lookupScheduleName(ctx, deps, sid),
				ID:                  ppid,
				PlanID:              pp.GetPlanId(),
				PlanLabel:           labelFromOptions(planOpts, pp.GetPlanId()),
				PlanOptions:         planOpts,
				Name:                pp.GetName(),
				Description:         pp.GetDescription(),
				Amount:              strconv.FormatFloat(float64(pp.GetAmount())/100.0, 'f', 2, 64),
				Currency:            pp.GetCurrency(),
				DurationValue:       fmt.Sprintf("%d", pp.GetDurationValue()),
				DurationUnit:        pp.GetDurationUnit(),
				CommonLabels:        deps.CommonLabels,
				PricingLocked:       pricingLocked,
				PricingLockedReason: pricingLockedReason,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.PlanLabels.Errors.UpdateFailed)
		}
		r := viewCtx.Request
		amount := int64(0)
		if v, err := strconv.ParseFloat(r.FormValue("amount"), 64); err == nil {
			amount = int64(math.Round(v * 100))
		}
		dv, _ := strconv.ParseInt(r.FormValue("duration_value"), 10, 32)
		currency := r.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}
		// Read existing to preserve active state (not in form) and to enforce
		// pricing-field immutability when the plan is in use by active subscriptions.
		existing, _ := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: ppid}})
		active := true
		if existing != nil && len(existing.GetData()) > 0 {
			active = existing.GetData()[0].GetActive()
		}

		// Server-side guard: if this plan is referenced by active subscriptions,
		// overwrite the four pricing fields with the existing DB values so a client
		// cannot bypass the read-only drawer by editing the HTML.
		durationUnit := r.FormValue("duration_unit")
		if deps.GetPricePlanInUseIDs != nil && existing != nil && len(existing.GetData()) > 0 {
			inUseMap, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppid})
			if inUseMap[ppid] {
				ex := existing.GetData()[0]
				amount = ex.GetAmount()
				currency = ex.GetCurrency()
				dv = int64(ex.GetDurationValue())
				durationUnit = ex.GetDurationUnit()
			}
		}

		planPageName := r.FormValue("name")
		planPageDesc := r.FormValue("description")
		req := &priceplanpb.UpdatePricePlanRequest{
			Data: &priceplanpb.PricePlan{
				Id:            ppid,
				PlanId:        r.FormValue("plan_id"),
				Name:          &planPageName,
				Description:   &planPageDesc,
				Amount:        amount,
				Currency:      currency,
				DurationValue: int32(dv),
				DurationUnit:  durationUnit,
				Active:        active,
			},
		}
		req.Data.PriceScheduleId = &sid
		if _, err := deps.UpdatePricePlan(ctx, req); err != nil {
			log.Printf("Failed to update price plan %s under schedule %s: %v", ppid, sid, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedule-plans-table")
	})
}

// NewDeleteAction handles POST /action/price-schedule/{id}/plan/{ppid}/delete.
// Hard delete — PricePlan rows are removed permanently (matches price_schedule's delete semantics).
func NewDeleteAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "delete") {
			return centymo.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		ppid := viewCtx.Request.PathValue("ppid")
		if ppid == "" {
			_ = viewCtx.Request.ParseForm()
			ppid = viewCtx.Request.FormValue("id")
		}
		if ppid == "" {
			return centymo.HTMXError(deps.PlanLabels.Errors.NotFound)
		}
		if _, err := deps.DeletePricePlan(ctx, &priceplanpb.DeletePricePlanRequest{Data: &priceplanpb.PricePlan{Id: ppid}}); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedule-plans-table")
	})
}

// NewProductPriceAddAction handles add under the schedule namespace.
func NewProductPriceAddAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_price_plan", "create") {
			return centymo.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		if deps.CreateProductPricePlan == nil {
			return centymo.HTMXError("Product price plan create is not available")
		}
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")

		if viewCtx.Request.Method == http.MethodGet {
			planID := loadPricePlanPlanID(ctx, deps, ppid)
			return view.OK("price-schedule-plan-product-price-drawer", &ProductPriceFormData{
				FormAction:        route.ResolveURL(deps.Routes.PlanProductPriceAddURL, "id", sid, "ppid", ppid),
				ScheduleID:        sid,
				PricePlanID:       ppid,
				Currency:          loadPricePlanCurrency(ctx, deps, ppid),
				ProductOptions:    loadProductOptions(ctx, deps, planID, ""),
				SectionTitle:      deps.ScheduleLabels.Detail.ProductPriceSection,
				ProductFieldLabel: deps.ScheduleLabels.Detail.ProductField,
				CommonLabels:      deps.CommonLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		productID := viewCtx.Request.FormValue("product_id")
		if productID == "" {
			return centymo.HTMXError("Product is required")
		}
		priceCentavos, ok := parsePriceCentavos(viewCtx.Request.FormValue("price"))
		if !ok {
			return centymo.HTMXError("Invalid price value")
		}
		currency := viewCtx.Request.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}
		record := &productpriceplanpb.ProductPricePlan{
			PricePlanId: ppid,
			ProductId:   productID,
			Price:       priceCentavos,
			Currency:    currency,
			Active:      true,
		}
		if _, err := deps.CreateProductPricePlan(ctx, &productpriceplanpb.CreateProductPricePlanRequest{Data: record}); err != nil {
			log.Printf("Failed to create product price plan for plan %s (schedule %s): %v", ppid, sid, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedule-plan-product-prices-table")
	})
}

// NewProductPriceEditAction handles edit under the schedule namespace.
func NewProductPriceEditAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_price_plan", "update") {
			return centymo.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		if deps.UpdateProductPricePlan == nil {
			return centymo.HTMXError("Product price plan update is not available")
		}
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")
		pppid := viewCtx.Request.PathValue("pppid")

		existing, err := findProductPricePlan(ctx, deps, pppid)
		if err != nil {
			return centymo.HTMXError(err.Error())
		}

		if viewCtx.Request.Method == http.MethodGet {
			planID := loadPricePlanPlanID(ctx, deps, ppid)
			currency := existing.GetCurrency()
			if currency == "" {
				currency = "PHP"
			}
			return view.OK("price-schedule-plan-product-price-drawer", &ProductPriceFormData{
				FormAction:        route.ResolveURL(deps.Routes.PlanProductPriceEditURL, "id", sid, "ppid", ppid, "pppid", pppid),
				IsEdit:            true,
				ID:                pppid,
				ScheduleID:        sid,
				PricePlanID:       ppid,
				ProductID:         existing.GetProductId(),
				Price:             fmt.Sprintf("%.2f", float64(existing.GetPrice())/100.0),
				Currency:          currency,
				ProductOptions:    loadProductOptions(ctx, deps, planID, existing.GetProductId()),
				SectionTitle:      deps.ScheduleLabels.Detail.ProductPriceSection,
				ProductFieldLabel: deps.ScheduleLabels.Detail.ProductField,
				CommonLabels:      deps.CommonLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		productID := viewCtx.Request.FormValue("product_id")
		if productID == "" {
			return centymo.HTMXError("Product is required")
		}
		priceCentavos, ok := parsePriceCentavos(viewCtx.Request.FormValue("price"))
		if !ok {
			return centymo.HTMXError("Invalid price value")
		}
		currency := viewCtx.Request.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}
		updated := &productpriceplanpb.ProductPricePlan{
			Id:          pppid,
			PricePlanId: ppid,
			ProductId:   productID,
			Price:       priceCentavos,
			Currency:    currency,
			Active:      existing.GetActive(),
		}
		if _, err := deps.UpdateProductPricePlan(ctx, &productpriceplanpb.UpdateProductPricePlanRequest{Data: updated}); err != nil {
			log.Printf("Failed to update product price plan %s: %v", pppid, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedule-plan-product-prices-table")
	})
}

// NewProductPriceDeleteAction handles delete under the schedule namespace.
func NewProductPriceDeleteAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_price_plan", "delete") {
			return centymo.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		if deps.DeleteProductPricePlan == nil {
			return centymo.HTMXError("Product price plan delete is not available")
		}
		_ = viewCtx.Request.ParseForm()
		pppid := viewCtx.Request.FormValue("id")
		if pppid == "" {
			pppid = viewCtx.Request.URL.Query().Get("id")
		}
		if pppid == "" {
			return centymo.HTMXError("ID is required")
		}
		if _, err := deps.DeleteProductPricePlan(ctx, &productpriceplanpb.DeleteProductPricePlanRequest{
			Data: &productpriceplanpb.ProductPricePlan{Id: pppid},
		}); err != nil {
			log.Printf("Failed to delete product price plan %s: %v", pppid, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-schedule-plan-product-prices-table")
	})
}

// ---------------------------------------------------------------------------
// buildPageData + helpers
// ---------------------------------------------------------------------------

func buildPageData(ctx context.Context, deps *DetailViewDeps, sid, ppid, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: ppid},
	})
	if err != nil {
		log.Printf("Failed to read price plan %s under schedule %s: %v", ppid, sid, err)
		return nil, fmt.Errorf("failed to load price plan: %w", err)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("price plan not found")
	}
	pp := data[0]

	currency := pp.GetCurrency()
	if currency == "" {
		currency = "PHP"
	}
	amountCell := types.MoneyCell(float64(pp.GetAmount()), currency, true)

	duration := ""
	if dv := pp.GetDurationValue(); dv > 0 {
		duration = fmt.Sprintf("%d %s", dv, pp.GetDurationUnit())
	}

	status := "active"
	statusVariant := "success"
	if !pp.GetActive() {
		status = "inactive"
		statusVariant = "warning"
	}

	scheduleName := lookupScheduleName(ctx, deps, sid)
	scheduleBack := route.ResolveURL(deps.Routes.DetailURL, "id", sid) + "?tab=plans"

	// Product price count for tab badge
	count := 0
	if deps.ListProductPricePlans != nil {
		pppResp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{})
		if err == nil {
			for _, item := range pppResp.GetData() {
				if item.GetPricePlanId() == ppid {
					count++
				}
			}
		}
	}

	base := route.ResolveURL(deps.Routes.PlanDetailURL, "id", sid, "ppid", ppid)
	action := route.ResolveURL(deps.Routes.PlanTabActionURL, "id", sid, "ppid", ppid, "tab", "")
	tabItems := []pyeza.TabItem{
		{Key: "info", Label: deps.ScheduleLabels.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "product-prices", Label: deps.ScheduleLabels.Tabs.ProductPrices, Href: base + "?tab=product-prices", HxGet: action + "product-prices", Icon: "icon-package", Count: count},
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          pp.GetName(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    pp.GetName(),
			HeaderSubtitle: fmt.Sprintf("under %s", scheduleName),
			HeaderIcon:     "icon-tag",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate:        "price-schedule-plan-detail-content",
		ScheduleID:             sid,
		ScheduleName:           scheduleName,
		ScheduleBackURL:        scheduleBack,
		PricePlan:              pp,
		Labels:                 deps.PlanLabels,
		ActiveTab:              activeTab,
		TabItems:               tabItems,
		ID:                     ppid,
		Name:                   pp.GetName(),
		Description:            pp.GetDescription(),
		Amount:                 amountCell,
		Currency:               currency,
		Duration:               duration,
		Status:                 status,
		StatusVariant:          statusVariant,
		CreatedDate:            pp.GetDateCreatedString(),
		ModifiedDate:           pp.GetDateModifiedString(),
		EditURL:                route.ResolveURL(deps.Routes.PlanEditURL, "id", sid, "ppid", ppid),
		ProductPriceEmptyTitle: deps.ScheduleLabels.Detail.ProductPriceEmptyTitle,
		ProductPriceEmptyMsg:   deps.ScheduleLabels.Detail.ProductPriceEmptyMsg,
	}

	if activeTab == "product-prices" {
		pageData.ProductPricesTable = buildProductPricesTable(ctx, deps, sid, ppid)
	}
	return pageData, nil
}

func buildProductPricesTable(ctx context.Context, deps *DetailViewDeps, sid, ppid string) *types.TableConfig {
	perms := view.GetUserPermissions(ctx)
	l := deps.PlanLabels

	columns := []types.TableColumn{
		{Key: "product", Label: deps.ScheduleLabels.Detail.ProductPriceColumnProduct, Sortable: true},
		{Key: "price", Label: deps.ScheduleLabels.Detail.ProductPriceColumnPrice, Sortable: true, WidthClass: "col-4xl"},
	}

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

	refreshURL := route.ResolveURL(deps.Routes.PlanTabActionURL, "id", sid, "ppid", ppid, "tab", "product-prices")
	rows := []types.TableRow{}
	if deps.ListProductPricePlans != nil {
		pppResp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{})
		if err != nil {
			log.Printf("Failed to list product price plans: %v", err)
		} else {
			for _, item := range pppResp.GetData() {
				if item == nil || item.GetPricePlanId() != ppid {
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
							Label:           deps.ScheduleLabels.Detail.ProductPriceEdit,
							Action:          "edit",
							URL:             route.ResolveURL(deps.Routes.PlanProductPriceEditURL, "id", sid, "ppid", ppid, "pppid", itemID),
							DrawerTitle:     deps.ScheduleLabels.Detail.ProductPriceEdit,
							Disabled:        !perms.Can("product_price_plan", "update"),
							DisabledTooltip: l.Errors.Unauthorized,
						},
						{
							Type:            "delete",
							Label:           deps.ScheduleLabels.Detail.ProductPriceDelete,
							Action:          "delete",
							URL:             deps.Routes.PlanProductPriceDeleteURL,
							ItemName:        productName,
							ConfirmTitle:    deps.ScheduleLabels.Detail.ProductPriceDelete,
							ConfirmMessage:  fmt.Sprintf(deps.ScheduleLabels.Detail.ProductPriceDeleteConfirm, productName),
							Disabled:        !perms.Can("product_price_plan", "delete"),
							DisabledTooltip: l.Errors.Unauthorized,
						},
					},
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   "price-schedule-plan-product-prices-table",
		RefreshURL:           refreshURL,
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
			Title:   deps.ScheduleLabels.Detail.ProductPriceEmptyTitle,
			Message: deps.ScheduleLabels.Detail.ProductPriceEmptyMsg,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           deps.ScheduleLabels.Detail.ProductPriceAdd,
			ActionURL:       route.ResolveURL(deps.Routes.PlanProductPriceAddURL, "id", sid, "ppid", ppid),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("product_price_plan", "create"),
			DisabledTooltip: l.Errors.Unauthorized,
		},
	}
	types.ApplyTableSettings(cfg)
	return cfg
}

func findProductPricePlan(ctx context.Context, deps *DetailViewDeps, pppid string) (*productpriceplanpb.ProductPricePlan, error) {
	if deps.ListProductPricePlans == nil {
		return nil, fmt.Errorf("product price plans not available")
	}
	resp, err := deps.ListProductPricePlans(ctx, &productpriceplanpb.ListProductPricePlansRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to load product price plans")
	}
	for _, item := range resp.GetData() {
		if item != nil && item.GetId() == pppid {
			return item, nil
		}
	}
	return nil, fmt.Errorf("product price plan not found")
}

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

func loadProductOptions(ctx context.Context, deps *DetailViewDeps, planID, selectedProductID string) []types.SelectOption {
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

	if deps.ListProductPlans == nil || planID == "" {
		// Fallback: all products
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

func buildPlanOptions(ctx context.Context, deps *DetailViewDeps, selectedID string) []map[string]any {
	if deps.ListPlans == nil {
		return nil
	}
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return nil
	}
	opts := make([]map[string]any, 0, len(resp.GetData()))
	for _, p := range resp.GetData() {
		if p == nil || !p.GetActive() {
			continue
		}
		opts = append(opts, map[string]any{
			"Value":       p.GetId(),
			"Label":       p.GetName(),
			"Description": p.GetDescription(),
			"Selected":    p.GetId() == selectedID,
		})
	}
	return opts
}

func lookupScheduleName(ctx context.Context, deps *DetailViewDeps, scheduleID string) string {
	if deps.ReadPriceSchedule == nil {
		return scheduleID
	}
	resp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{
		Data: &priceschedulepb.PriceSchedule{Id: scheduleID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return scheduleID
	}
	if n := resp.GetData()[0].GetName(); n != "" {
		return n
	}
	return scheduleID
}

func parsePriceCentavos(s string) (int64, bool) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || f < 0 {
		return 0, false
	}
	return int64(math.Round(f * 100)), true
}

// labelFromOptions returns the Label string for the option whose Value matches id.
// Used to populate SelectedLabel on the edit-drawer auto-complete.
func labelFromOptions(opts []map[string]any, id string) string {
	for _, opt := range opts {
		if v, ok := opt["Value"].(string); ok && v == id {
			if label, ok := opt["Label"].(string); ok {
				return label
			}
		}
	}
	return ""
}
