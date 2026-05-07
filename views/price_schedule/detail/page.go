package detail

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/price_plan/form"
	"github.com/erniealice/hybra-golang/views/attachment"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// DetailViewDeps holds view dependencies for the price schedule detail page.
type DetailViewDeps struct {
	Routes centymo.PriceScheduleRoutes
	Labels centymo.PriceScheduleLabels
	// PricePlanLabels is the authoritative source for the price-plan drawer
	// form (sourced from lyngua price_plan.json → price_plan.form). Used when
	// opening the schedule-scoped Add drawer so all Wave 2 fields render.
	PricePlanLabels centymo.PricePlanLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels

	ReadPriceSchedule func(ctx context.Context, req *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error)
	ListLocations     func(ctx context.Context, req *locationpb.ListLocationsRequest) (*locationpb.ListLocationsResponse, error)
	ListPricePlans    func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ListPlans         func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	CreatePricePlan   func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)

	// Auto-seed product_price_plan rows on PricePlan create. When a package is added
	// under a rate-card, one ProductPricePlan row is created per linked product_plan,
	// copying price/currency from the Product record so the newly-created PricePlan's
	// "product-prices" tab is pre-populated.
	ListProductPlans       func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProducts           func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	CreateProductPricePlan func(ctx context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)

	// Reference checker: returns a map of price_plan_id → true for plans in use by active subscriptions.
	// Delete is disabled for in-use plans; Edit remains enabled (Pricing fields lock inside the drawer).
	GetPricePlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	attachment.AttachmentOps
}

// PageData holds the data for the price schedule detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Schedule        *priceschedulepb.PriceSchedule
	Labels          centymo.PriceScheduleLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID             string
	Name           string
	Description    string
	DateStart      string
	DateEnd        string
	LocationName   string
	Status         string
	StatusVariant  string
	CreatedDate    string
	ModifiedDate   string

	PlansTable      *types.TableConfig
	AttachmentTable *types.TableConfig
}

// NewView creates the price schedule detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := deps.Labels.Tabs.CanonicalizeTab(viewCtx.Request.URL.Query().Get("tab"))
		if activeTab == "" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("price-schedule-detail", pageData)
	})
}

// NewTabAction handles GET /action/price-schedule/{id}/tab/{tab}.
// Returns only the tab partial template, for HTMX tab switching.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := deps.Labels.Tabs.CanonicalizeTab(viewCtx.Request.PathValue("tab"))
		if tab == "" {
			tab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		templateName := "price-schedule-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// NewPlanAddAction handles GET/POST /action/price-schedule/{id}/plan/add.
// GET renders a drawer with price_schedule_id pre-locked; POST creates the PricePlan.
func NewPlanAddAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if deps.CreatePricePlan == nil {
			return centymo.HTMXError(deps.Labels.Errors.PricePlanCreateUnavailable)
		}
		scheduleID := viewCtx.Request.PathValue("id")
		if scheduleID == "" {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}

		if viewCtx.Request.Method == http.MethodGet {
			scheduleName, scheduleClientID := lookupScheduleNameAndClient(ctx, deps, scheduleID)
			// 2026-05-03 — when the parent schedule is client-scoped, the
			// Package picker is filtered to plans that are either general-scope
			// (no client_id) or assigned to the same client.
			planOptions := buildPlanOptions(ctx, deps, scheduleClientID)
			formLabels := deps.PricePlanLabels.Form
			return view.OK("price-plan-drawer-form", &form.Data{
				FormAction:             route.ResolveURL(deps.Routes.PlanAddURL, "id", scheduleID),
				Context:                form.ContextSchedule,
				ScheduleID:             scheduleID,
				ScheduleName:           scheduleName,
				ParentScheduleClientID: scheduleClientID,
				Active:                 true,
				Currency:               "PHP",
				DurationUnit:           "months",
				BillingKind:            "BILLING_KIND_RECURRING",
				AmountBasis:            "AMOUNT_BASIS_PER_CYCLE",
				BillingCycleUnit:       "month",
				TermUnit:               "month",
				PlanOptions:            planOptions,
				DurationUnitOptions:    form.BuildDurationUnitOptions(deps.CommonLabels),
				Labels:                 form.LabelsFromPricePlan(formLabels),
				CommonLabels:           deps.CommonLabels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		r := viewCtx.Request
		planID := r.FormValue("plan_id")
		if planID == "" {
			return centymo.HTMXError(deps.Labels.Detail.PlanRequired)
		}
		amount := int64(0)
		if v, err := strconv.ParseFloat(r.FormValue("amount"), 64); err == nil {
			amount = int64(math.Round(v * 100))
		}
		currency := r.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}

		ppName := r.FormValue("name")
		ppDesc := r.FormValue("description")
		pp := &priceplanpb.PricePlan{
			PlanId:          planID,
			Name:            &ppName,
			Description:     &ppDesc,
			BillingAmount:   amount,
			BillingCurrency: currency,
			Active:          true,
		}
		pp.PriceScheduleId = &scheduleID
		// Phase 1 legacy dual-write — proto fields now optional; only set when present.
		if dvStr := r.FormValue("duration_value"); dvStr != "" {
			if parsed, err := strconv.ParseInt(dvStr, 10, 32); err == nil {
				dv32 := int32(parsed)
				pp.DurationValue = &dv32
			}
		}
		if du := r.FormValue("duration_unit"); du != "" {
			pp.DurationUnit = &du
		}
		applyBillingFieldsFromRequest(pp, r)

		createResp, err := deps.CreatePricePlan(ctx, &priceplanpb.CreatePricePlanRequest{Data: pp})
		if err != nil {
			log.Printf("Failed to create price plan from schedule %s: %v", scheduleID, err)
			return centymo.HTMXError(err.Error())
		}

		// Auto-seed product_price_plan rows: one per product_plan linked to this plan_id,
		// copying price/currency from the Product record. Failures here are non-fatal —
		// the main PricePlan create already succeeded.
		autoSeedProductPricePlans(ctx, deps, createResp, planID)

		return centymo.HTMXSuccess("price-schedule-plans-table")
	})
}

// autoSeedProductPricePlans creates one ProductPricePlan row per product_plan linked
// to planID, copying price/currency from the underlying Product record. This runs after
// a successful CreatePricePlan so the newly-created PricePlan's "product-prices" tab is
// pre-populated. All failures are logged and non-fatal.
func autoSeedProductPricePlans(ctx context.Context, deps *DetailViewDeps, createResp *priceplanpb.CreatePricePlanResponse, planID string) {
	if deps.ListProductPlans == nil || deps.ListProducts == nil || deps.CreateProductPricePlan == nil {
		return
	}
	if createResp == nil || len(createResp.GetData()) == 0 {
		log.Printf("auto-seed product_price_plan skipped: CreatePricePlan response had no data")
		return
	}
	createdID := createResp.GetData()[0].GetId()
	if createdID == "" {
		log.Printf("auto-seed product_price_plan skipped: created PricePlan has no ID")
		return
	}

	// Load product_plans for this plan_id (same filter pattern as
	// views/plan/action/product_plan.go loadExistingProductIDs).
	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
		Filters: &commonpb.FilterRequest{
			Logic: commonpb.FilterLogic_AND,
			Filters: []*commonpb.TypedFilter{
				{
					Field: "plan_id",
					FilterType: &commonpb.TypedFilter_StringFilter{
						StringFilter: &commonpb.StringFilter{
							Value:    planID,
							Operator: commonpb.StringOperator_STRING_EQUALS,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("auto-seed product_price_plan: failed to list product_plans for plan %s: %v", planID, err)
		return
	}
	if ppResp == nil || len(ppResp.GetData()) == 0 {
		return
	}

	// Build product_id → Product map for price/currency lookup.
	prodResp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		log.Printf("auto-seed product_price_plan: failed to list products: %v", err)
		return
	}
	products := map[string]*productpb.Product{}
	for _, p := range prodResp.GetData() {
		if p != nil {
			products[p.GetId()] = p
		}
	}

	for _, pp := range ppResp.GetData() {
		if pp == nil {
			continue
		}
		productID := pp.GetProductId()
		productPlanID := pp.GetId()
		if productPlanID == "" {
			continue
		}
		// Seed a row regardless of whether the Product can be resolved — zero
		// price + default currency when missing, so the nested "package-item-prices"
		// tab is always pre-populated and the user can edit values from there.
		var price int64
		currency := "PHP"
		if prod := products[productID]; prod != nil {
			// Product.price is optional under Model D (configurable products
			// have per-variant overrides instead). Only dereference when set.
			if prod.Price != nil {
				price = prod.GetPrice()
			}
			if c := prod.GetCurrency(); c != "" {
				currency = c
			}
		}
		if _, err := deps.CreateProductPricePlan(ctx, &productpriceplanpb.CreateProductPricePlanRequest{
			Data: &productpriceplanpb.ProductPricePlan{
				PricePlanId:     createdID,
				ProductPlanId:   productPlanID, // Model D — FK to catalog line
				BillingAmount:   price,
				BillingCurrency: currency,
				Active:          true,
			},
		}); err != nil {
			log.Printf("auto-seed product_price_plan failed for %s/%s: %v", createdID, productPlanID, err)
		}
	}
}

func lookupScheduleNameAndClient(ctx context.Context, deps *DetailViewDeps, scheduleID string) (name, clientID string) {
	if deps.ReadPriceSchedule == nil {
		return scheduleID, ""
	}
	resp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{
		Data: &priceschedulepb.PriceSchedule{Id: scheduleID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return scheduleID, ""
	}
	ps := resp.GetData()[0]
	name = ps.GetName()
	if name == "" {
		name = scheduleID
	}
	return name, ps.GetClientId()
}

// buildPlanOptions returns the plan picker options with a strict, mutually-
// exclusive client-scope filter:
//
//   - Client-scoped schedule (scheduleClientID != ""): only plans whose
//     client_id equals scheduleClientID. General-scope plans are excluded.
//   - General-scope schedule (scheduleClientID == ""): only plans with no
//     client_id. Client-scoped plans cannot be attached to a general schedule.
func buildPlanOptions(ctx context.Context, deps *DetailViewDeps, scheduleClientID string) []map[string]any {
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
		planClientID := p.GetClientId()
		if scheduleClientID != "" {
			if planClientID != scheduleClientID {
				continue
			}
		} else {
			if planClientID != "" {
				continue
			}
		}
		opts = append(opts, map[string]any{
			"Value":       p.GetId(),
			"Label":       p.GetName(),
			"Description": p.GetDescription(),
		})
	}
	return opts
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{
		Data: &priceschedulepb.PriceSchedule{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read price schedule %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}

	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	ps := data[0]

	locationName := deps.Labels.Detail.NoLocation
	if locID := ps.GetLocationId(); locID != "" {
		if n := lookupLocationName(ctx, deps, locID); n != "" {
			locationName = n
		} else {
			locationName = locID
		}
	}

	status := "active"
	statusVariant := "success"
	if !ps.GetActive() {
		status = "inactive"
		statusVariant = "warning"
	}

	l := deps.Labels
	planCount := countPlansForSchedule(ctx, deps, ps)
	tabItems := buildTabItems(id, l, planCount, deps.Routes)

	headerSubtitle := strings.TrimSpace(ps.GetDescription())
	if headerSubtitle == "" {
		headerSubtitle = l.Detail.NoDescriptionSubtitle
	}

	tz := types.LocationFromContext(ctx)

	createdDate := ""
	if ms := ps.GetDateCreated(); ms > 0 {
		createdDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}
	modifiedDate := ""
	if ms := ps.GetDateModified(); ms > 0 {
		modifiedDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          ps.GetName(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    ps.GetName(),
			HeaderSubtitle: headerSubtitle,
			HeaderIcon:     "icon-calendar",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "price-schedule-detail-content",
		Schedule:        ps,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		Name:            ps.GetName(),
		Description:     ps.GetDescription(),
		DateStart:       types.FormatTimestampInTZ(ps.GetDateTimeStart(), tz, types.DateTimeReadable),
		DateEnd:         types.FormatTimestampInTZ(ps.GetDateTimeEnd(), tz, types.DateTimeReadable),
		LocationName:    locationName,
		Status:          status,
		StatusVariant:   statusVariant,
		CreatedDate:     createdDate,
		ModifiedDate:    modifiedDate,
	}

	if activeTab == "pricePlan" {
		pageData.PlansTable = buildPlansTable(ctx, deps, ps)
	}
	if activeTab == "attachments" && deps.ListAttachments != nil {
		cfg := attachmentConfig(deps)
		var attachItems []*attachmentpb.Attachment
		if resp, err := deps.ListAttachments(ctx, cfg.EntityType, id); err == nil && resp != nil {
			attachItems = resp.GetData()
		}
		pageData.AttachmentTable = attachment.BuildTable(attachItems, cfg, id)
	}

	return pageData, nil
}

func buildTabItems(id string, l centymo.PriceScheduleLabels, planCount int, routes centymo.PriceScheduleRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	pricePlanSlug := l.Tabs.ResolveTabSlug("pricePlan")
	attachmentsLabel := l.Detail.TabAttachments
	if attachmentsLabel == "" {
		attachmentsLabel = "Attachments"
	}
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "pricePlan", Label: l.Tabs.PricePlan, Href: base + "?tab=" + pricePlanSlug, HxGet: action + pricePlanSlug, Icon: "icon-layers", Count: planCount},
		{Key: "attachments", Label: attachmentsLabel, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
	}
}

// countPlansForSchedule counts price_plans linked to this schedule via price_schedule_id FK.
func countPlansForSchedule(ctx context.Context, deps *DetailViewDeps, ps *priceschedulepb.PriceSchedule) int {
	if deps.ListPricePlans == nil {
		return 0
	}
	resp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
	if err != nil {
		return 0
	}
	schedID := ps.GetId()
	count := 0
	for _, pp := range resp.GetData() {
		if pp != nil && pp.GetPriceScheduleId() == schedID {
			count++
		}
	}
	return count
}

func buildPlansTable(ctx context.Context, deps *DetailViewDeps, ps *priceschedulepb.PriceSchedule) *types.TableConfig {
	l := deps.Labels

	columns := []types.TableColumn{
		{Key: "name", Label: l.Detail.PlanColumnPlan},
		{Key: "amount", Label: l.Detail.PlanColumnAmount, WidthClass: "col-3xl", Align: "right"},
		{Key: "billing_kind", Label: l.Detail.PlanColumnBillingKind, NoSort: true, WidthClass: "col-2xl"},
		{Key: "amount_basis", Label: l.Detail.PlanColumnAmountBasis, NoSort: true, WidthClass: "col-2xl"},
		{Key: "cadence", Label: l.Detail.PlanColumnCadence, NoSort: true, WidthClass: "col-3xl"},
		{Key: "status", Label: l.Detail.PlanColumnStatus, NoSort: true, WidthClass: "col-2xl"},
	}

	perms := view.GetUserPermissions(ctx)
	rows := []types.TableRow{}
	if deps.ListPricePlans != nil {
		resp, err := deps.ListPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
		if err != nil {
			log.Printf("Failed to list price plans for schedule %s: %v", ps.GetId(), err)
		} else {
			schedID := ps.GetId()
			// 2026-04-27 plan-client-scope plan §6.4. When the schedule itself is
			// master, the default Plans tab hides client-scoped PricePlans; the
			// `?show_client_specific=1` toggle on the toolbar opts back in. When
			// the schedule is client-scoped, the §3.2 cascade guarantees every
			// PricePlan inside shares its client_id, so the filter is a no-op.
			scheduleIsClientScoped := ps.GetClientId() != ""

			// Collect IDs for the reference checker (one batch call for the whole table).
			var ppIDs []string
			for _, pp := range resp.GetData() {
				if pp != nil && pp.GetPriceScheduleId() == schedID {
					if !scheduleIsClientScoped && pp.GetClientId() != "" {
						// Hidden by default on master schedules; the toolbar
						// toggle in the template opts in via a query param,
						// which the View handler can pass through if needed.
						continue
					}
					ppIDs = append(ppIDs, pp.GetId())
				}
			}
			inUseIDs := map[string]bool{}
			if deps.GetPricePlanInUseIDs != nil && len(ppIDs) > 0 {
				inUseIDs, _ = deps.GetPricePlanInUseIDs(ctx, ppIDs)
			}

			// Build plan ID → name map for fallback display when price_plan.Name is blank.
			planNames := map[string]string{}
			if deps.ListPlans != nil {
				planResp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
				if err == nil {
					for _, p := range planResp.GetData() {
						if p != nil {
							planNames[p.GetId()] = p.GetName()
						}
					}
				}
			}

			for _, pp := range resp.GetData() {
				if pp == nil || pp.GetPriceScheduleId() != schedID {
					continue
				}
				// Same client-scope filter as the ID collection loop above.
				if !scheduleIsClientScoped && pp.GetClientId() != "" {
					continue
				}
				ppID := pp.GetId()
				name := strings.TrimSpace(pp.GetName())
				if name == "" {
					name = planNames[pp.GetPlanId()]
				}
				currency := pp.GetBillingCurrency()
				if currency == "" {
					currency = "PHP"
				}
				cycle := ""
				if cv := pp.GetBillingCycleValue(); cv > 0 {
					cycle = pyeza.FormatDuration(cv, pp.GetBillingCycleUnit(), deps.CommonLabels.DurationUnit)
				} else if dv := pp.GetDurationValue(); dv > 0 {
					// Fallback to legacy duration_* dual-write fields.
					cycle = pyeza.FormatDuration(dv, pp.GetDurationUnit(), deps.CommonLabels.DurationUnit)
				}
				cadence := formatCadence(pp.GetBillingKind(), cycle, l.Detail)
				billingKindLabel := formatBillingKindLabel(pp.GetBillingKind(), l.Detail)
				amountBasisLabel := formatAmountBasisLabel(pp.GetAmountBasis(), l.Detail)
				planStatus := "active"
				planVariant := "success"
				if !pp.GetActive() {
					planStatus = "inactive"
					planVariant = "warning"
				}

				inUse := inUseIDs[ppID]
				deleteDisabled := !perms.Can("price_plan", "delete") || inUse
				deleteTooltip := l.Errors.Unauthorized
				if inUse {
					deleteTooltip = l.Detail.PlanInUseTooltip
				}

				rows = append(rows, types.TableRow{
					ID: ppID,
					Cells: []types.TableCell{
						{Type: "text", Value: name},
						types.MoneyCell(float64(pp.GetBillingAmount()), currency, true),
						{Type: "badge", Value: billingKindLabel, Variant: "info"},
						{Type: "text", Value: amountBasisLabel},
						{Type: "text", Value: cadence},
						{Type: "badge", Value: planStatus, Variant: planVariant},
					},
					Actions: []types.TableAction{
						{
							Type:   "view",
							Label:  l.Detail.PlanView,
							Action: "view",
							Href:   route.ResolveURL(deps.Routes.PlanDetailURL, "id", schedID, "ppid", ppID),
						},
						{
							Type:            "edit",
							Label:           l.Detail.PlanEdit,
							Action:          "edit",
							URL:             route.ResolveURL(deps.Routes.PlanEditURL, "id", schedID, "ppid", ppID),
							DrawerTitle:     l.Detail.PlanEditDrawerTitle,
							Disabled:        !perms.Can("price_plan", "update"),
							DisabledTooltip: l.Errors.Unauthorized,
						},
						{
							Type:            "delete",
							Label:           l.Detail.PlanDelete,
							Action:          "delete",
							URL:             route.ResolveURL(deps.Routes.PlanDeleteURL, "id", schedID, "ppid", ppID),
							ItemName:        name,
							ConfirmTitle:    l.Detail.PlanDeleteTitle,
							ConfirmMessage:  fmt.Sprintf(l.Detail.PlanDeleteMsg, name),
							Disabled:        deleteDisabled,
							DisabledTooltip: deleteTooltip,
						},
					},
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)


	refreshURL := route.ResolveURL(deps.Routes.TabActionURL, "id", ps.GetId(), "tab", l.Tabs.ResolveTabSlug("pricePlan"))
	tableConfig := &types.TableConfig{
		ID:                   "price-schedule-plans-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.PlansEmptyTitle,
			Message: l.Detail.PlansEmptyMsg,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Detail.PlanAdd,
			ActionURL:       route.ResolveURL(deps.Routes.PlanAddURL, "id", ps.GetId()),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("price_plan", "create"),
			DisabledTooltip: l.Errors.Unauthorized,
		},
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig
}

func lookupLocationName(ctx context.Context, deps *DetailViewDeps, locationID string) string {
	if deps.ListLocations == nil || locationID == "" {
		return ""
	}
	resp, err := deps.ListLocations(ctx, &locationpb.ListLocationsRequest{})
	if err != nil {
		return ""
	}
	for _, loc := range resp.GetData() {
		if loc.GetId() == locationID {
			return loc.GetName()
		}
	}
	return ""
}

// applyBillingFieldsFromRequest writes Wave 2 billing-semantics fields
// (billing_kind, amount_basis, billing_cycle_*, default_term_*) from the POST
// body onto pp. Mirrors the equivalent helper in views/plan/action.
func applyBillingFieldsFromRequest(pp *priceplanpb.PricePlan, r *http.Request) {
	if v := r.FormValue("billing_kind"); v != "" {
		if bk, ok := priceplanpb.BillingKind_value[v]; ok {
			pp.BillingKind = priceplanpb.BillingKind(bk)
		}
	}
	if v := r.FormValue("amount_basis"); v != "" {
		if ab, ok := priceplanpb.AmountBasis_value[v]; ok {
			pp.AmountBasis = priceplanpb.AmountBasis(ab)
		}
	}
	if s := r.FormValue("billing_cycle_value"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil {
			v32 := int32(n)
			pp.BillingCycleValue = &v32
		}
	}
	if u := r.FormValue("billing_cycle_unit"); u != "" {
		pp.BillingCycleUnit = &u
	}
	if s := r.FormValue("default_term_value"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil {
			v32 := int32(n)
			pp.DefaultTermValue = &v32
		}
	}
	if u := r.FormValue("default_term_unit"); u != "" {
		pp.DefaultTermUnit = &u
	}
	if s := r.FormValue("entitled_occurrences"); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil {
			v32 := int32(n)
			pp.EntitledOccurrences = &v32
		}
	}
}

// formatCadence renders a billing-kind-aware cadence string for the plans
// table. Recurring and contract kinds embed the formatted cycle string (e.g.
// "monthly") into a tier-translatable template via fmt.Sprintf("%s"); the
// other kinds return a kind-specific phrase that doesn't need a cycle.
func formatCadence(kind priceplanpb.BillingKind, cycle string, l centymo.PriceScheduleDetailLabels) string {
	switch kind {
	case priceplanpb.BillingKind_BILLING_KIND_ONE_TIME:
		return l.CadenceOneTime
	case priceplanpb.BillingKind_BILLING_KIND_RECURRING:
		if cycle == "" {
			return l.CadenceUnspecified
		}
		return fmt.Sprintf(l.CadenceRecurring, cycle)
	case priceplanpb.BillingKind_BILLING_KIND_CONTRACT:
		if cycle == "" {
			return l.BillingKindContract
		}
		return fmt.Sprintf(l.CadenceContract, cycle)
	case priceplanpb.BillingKind_BILLING_KIND_MILESTONE:
		return l.CadenceMilestone
	case priceplanpb.BillingKind_BILLING_KIND_AD_HOC:
		return l.CadenceAdHoc
	default:
		if cycle != "" {
			return cycle
		}
		return l.CadenceUnspecified
	}
}

func formatBillingKindLabel(kind priceplanpb.BillingKind, l centymo.PriceScheduleDetailLabels) string {
	switch kind {
	case priceplanpb.BillingKind_BILLING_KIND_ONE_TIME:
		return l.BillingKindOneTime
	case priceplanpb.BillingKind_BILLING_KIND_RECURRING:
		return l.BillingKindRecurring
	case priceplanpb.BillingKind_BILLING_KIND_CONTRACT:
		return l.BillingKindContract
	case priceplanpb.BillingKind_BILLING_KIND_MILESTONE:
		return l.BillingKindMilestone
	case priceplanpb.BillingKind_BILLING_KIND_AD_HOC:
		return l.BillingKindAdHoc
	default:
		return ""
	}
}

func formatAmountBasisLabel(basis priceplanpb.AmountBasis, l centymo.PriceScheduleDetailLabels) string {
	switch basis {
	case priceplanpb.AmountBasis_AMOUNT_BASIS_PER_CYCLE:
		return l.AmountBasisPerCycle
	case priceplanpb.AmountBasis_AMOUNT_BASIS_TOTAL_PACKAGE:
		return l.AmountBasisTotalPackage
	case priceplanpb.AmountBasis_AMOUNT_BASIS_DERIVED_FROM_LINES:
		return l.AmountBasisDerivedFromLines
	case priceplanpb.AmountBasis_AMOUNT_BASIS_PER_OCCURRENCE:
		return l.AmountBasisPerOccurrence
	default:
		return ""
	}
}
