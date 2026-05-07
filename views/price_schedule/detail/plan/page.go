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
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/price_plan/form"
	"github.com/erniealice/hybra-golang/views/attachment"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// DetailViewDeps holds all dependencies for the schedule-scoped price_plan detail page.
type DetailViewDeps struct {
	Routes                 centymo.PriceScheduleRoutes
	ScheduleLabels         centymo.PriceScheduleLabels
	PlanLabels             centymo.PricePlanLabels
	ProductPricePlanLabels centymo.ProductPricePlanLabels
	CommonLabels           pyeza.CommonLabels
	TableLabels            types.TableLabels

	ReadPriceSchedule func(ctx context.Context, req *priceschedulepb.ReadPriceScheduleRequest) (*priceschedulepb.ReadPriceScheduleResponse, error)
	ReadPricePlan     func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	UpdatePricePlan   func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan   func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)

	ListPlans           func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListProducts        func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ListProductPlans    func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProductVariants func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)

	ListProductPricePlans  func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	CreateProductPricePlan func(ctx context.Context, req *productpriceplanpb.CreateProductPricePlanRequest) (*productpriceplanpb.CreateProductPricePlanResponse, error)
	UpdateProductPricePlan func(ctx context.Context, req *productpriceplanpb.UpdateProductPricePlanRequest) (*productpriceplanpb.UpdateProductPricePlanResponse, error)
	DeleteProductPricePlan func(ctx context.Context, req *productpriceplanpb.DeleteProductPricePlanRequest) (*productpriceplanpb.DeleteProductPricePlanResponse, error)

	// Reference checker: returns a map of price_plan_id → true for plans in use by active subscriptions.
	// When a plan is in use, Pricing fields in the Edit drawer are read-only.
	GetPricePlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// Mount overrides — populated only by the plan-scoped entry points
	// (NewPlanScopedView / NewPlanScopedTabAction) so the same render path
	// can be served under /app/plans/detail/{id}/price/{ppid} with the
	// sidebar anchored to Services > Packages instead of Rate Cards. When
	// any of these is empty, buildPageData falls back to the rate-card
	// defaults derived from Routes (PriceScheduleRoutes).
	ActiveNavOverride    string
	ActiveSubNavOverride string
	// PlanDetailBackURL is the package-detail URL pattern — e.g.
	// "/app/services/packages/detail/{id}". When set, buildPageData resolves
	// it with the loaded price_plan's plan_id and appends `?tab=` +
	// PlanDetailBackTab to point the breadcrumb back at the package's
	// package-prices tab. The breadcrumb label becomes the plan name.
	PlanDetailBackURL string
	PlanDetailBackTab string
	// PlanScopedDetailURL / TabActionURL are the new package-scoped routes
	// (PlanRoutes.PricePlanDetailURL / PricePlanTabActionURL). When set, tab
	// items render with these URLs so the address bar stays under the plan
	// namespace as the operator switches tabs.
	PlanScopedDetailURL    string
	PlanScopedTabActionURL string

	// 2026-05-04 — Engagements (subscriptions) tab dependencies. See
	// docs/plan/20260504-price-plan-engagements-tab/.
	ListSubscriptionsByPricePlan func(ctx context.Context, req *subscriptionpb.ListSubscriptionsByPricePlanRequest) (*subscriptionpb.ListSubscriptionsByPricePlanResponse, error)
	// SubscriptionEditURL / SubscriptionDeleteURL drive the row actions on
	// the subscriptions tab table. When empty, the row's edit/delete actions
	// render disabled (display-only). DetailURL is overridden by the nested
	// engagement URL when PlanEngagementDetailURL is non-empty.
	SubscriptionEditURL   string
	SubscriptionDeleteURL string
	// PlanEngagementDetailURL is the schedule-scoped engagement URL template
	// /app/price-schedules/detail/{id}/plan/{ppid}/engagement/{eid}. When
	// set, the row's "View" action targets this nested URL so the
	// subscription detail page renders with a rate-card → plan → engagement
	// breadcrumb. Empty falls back to SubscriptionDetailURL.
	PlanEngagementDetailURL string
	SubscriptionDetailURL   string

	attachment.AttachmentOps
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

	// 2026-05-04 — Engagements (subscriptions) tab payload.
	SubscriptionsTable *types.TableConfig

	AttachmentTable *types.TableConfig
}

// SubscriptionRow is one row in the schedule-scoped plan detail
// "Engagements" / "Subscriptions" tab table. Names + dates are pre-formatted
// for the tier's display TZ; the view-layer table builder consumes this struct
// directly into pyeza TableRow cells.
type SubscriptionRow struct {
	ID         string
	Name       string
	ClientID   string
	ClientName string
	Plan       string
	DateStart  string
	DateEnd    string
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
	DurationValue string // Phase 1 legacy dual-write
	DurationUnit  string // Phase 1 legacy dual-write
	CommonLabels  pyeza.CommonLabels
	Labels        centymo.PriceScheduleLabels

	// Wave 2: new billing semantics fields.
	//
	// 2026-04-30 enum-select-canonicalize — BillingKindOptions /
	// AmountBasisOptions removed; the drawer template hardcodes the option
	// list. Only the selected value is passed in.
	BillingKind         string
	AmountBasis         string
	BillingCycleValue   string
	BillingCycleUnit    string
	// kept as DefaultTermValue/Unit on the wire (form input names) but renamed
	// in the form.Data struct.
	DurationUnitOptions []types.SelectOption

	// PricingLocked is true when the price_plan is referenced by active subscriptions.
	// The Pricing section fields (Amount, Currency, Duration, DurationUnit) are rendered
	// as read-only in the drawer, but all other fields remain editable.
	PricingLocked       bool
	PricingLockedReason string
}

// ProductPriceFormData is the drawer form for editing a ProductPricePlan.
// Plan + Product sections are display-only context; only Price + Currency are editable.
// Rows are auto-seeded from product_plan assignments on PricePlan create, so the
// Model D catalog-line selection is fixed per row.
type ProductPriceFormData struct {
	FormAction    string
	IsEdit        bool
	ID            string
	ScheduleID    string
	PricePlanID   string
	ProductPlanID string
	Price         string
	Currency      string
	CommonLabels  pyeza.CommonLabels

	// Display-only context (read-only).
	PlanName           string
	PlanDescription    string
	ProductName        string
	ProductDescription string
	VariantName        string // SKU of the catalog line's variant, when any

	// Wave 2: billing treatment + effective date fields.
	//
	// 2026-04-30 enum-select-canonicalize — BillingTreatmentOptions removed;
	// the drawer template (_ppp-fields.html) hardcodes the option list.
	BillingTreatment        string
	DateStart               string // ISO 8601 (YYYY-MM-DD) or empty
	DateEnd                 string // ISO 8601 (YYYY-MM-DD) or empty

	// Parent PricePlan context — drives field visibility and currency lock.
	// Why: billing_treatment is meaningless when parent.billing_kind=ONE_TIME
	// (no cycles); per-line currency must match parent.billing_currency.
	ParentBillingKind  string // proto enum string, e.g. "BILLING_KIND_RECURRING"
	ParentAmountBasis  string // proto enum string, e.g. "AMOUNT_BASIS_PER_CYCLE"
	ShowTreatment      bool   // false when parent.billing_kind=ONE_TIME
	BasisBannerMessage string // contextual hint about the parent's amount_basis

	// Read-only "package context" block (Plan, Rate card, Billing model,
	// Amount basis, Billing cycle, Term, Currency) rendered above the
	// editable fields via the shared `ppp-parent-context` partial.
	BillingKindDisplay    string
	AmountBasisDisplay    string
	BillingCycleDisplay   string
	TermDisplay           string
	ParentCurrencyDisplay string
	RateCardName          string

	// Wave 2: labels for the new fields.
	ProductPricePlanLabels centymo.ProductPricePlanFormLabels
	PriceScheduleLabels    centymo.PriceScheduleDetailLabels // for section labels

	// PricingLocked is true when the parent PricePlan is referenced by an active
	// subscription — editing the per-item price would shift revenue allocation
	// on live engagements. Mirrors the PricePlan edit drawer's lock rule.
	PricingLocked       bool
	PricingLockedReason string
}

// NewView renders the full detail page at /app/price-schedules/detail/{id}/plan/{ppid}.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")

		activeTab := deps.ScheduleLabels.Tabs.CanonicalizeTab(viewCtx.Request.URL.Query().Get("tab"))
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

// NewPlanScopedView renders the same detail page but mounted under
// /app/plans/detail/{id}/price/{ppid}, where {id} is the plan_id (not the
// schedule_id). The handler resolves the underlying schedule_id by reading the
// price_plan, then delegates to buildPageData. The mount-override fields on
// deps anchor ActiveNav to Services > Packages and point the breadcrumb back
// at the plan detail's package-prices tab.
func NewPlanScopedView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		ppid := viewCtx.Request.PathValue("ppid")
		sid, err := resolveScheduleIDFromPricePlan(ctx, deps, ppid)
		if err != nil {
			return view.Error(err)
		}
		activeTab := deps.ScheduleLabels.Tabs.CanonicalizeTab(viewCtx.Request.URL.Query().Get("tab"))
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

// NewPlanScopedTabAction handles GET /action/plan/{id}/price/{ppid}/tab/{tab}.
// {id} is plan_id; resolves schedule_id from the price_plan record.
func NewPlanScopedTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		ppid := viewCtx.Request.PathValue("ppid")
		sid, err := resolveScheduleIDFromPricePlan(ctx, deps, ppid)
		if err != nil {
			return view.Error(err)
		}
		tab := deps.ScheduleLabels.Tabs.CanonicalizeTab(viewCtx.Request.PathValue("tab"))
		if tab == "" {
			tab = "info"
		}
		pageData, err := buildPageData(ctx, deps, sid, ppid, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}
		templateName := "price-schedule-plan-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// resolveScheduleIDFromPricePlan looks up a price_plan and returns its
// price_schedule_id. Used by the plan-scoped entry points to translate
// {plan_id, ppid} URL params into the {schedule_id, ppid} pair the existing
// buildPageData helper expects.
func resolveScheduleIDFromPricePlan(ctx context.Context, deps *DetailViewDeps, ppid string) (string, error) {
	if deps.ReadPricePlan == nil || ppid == "" {
		return "", fmt.Errorf("price plan not found")
	}
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: ppid},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return "", fmt.Errorf("price plan not found")
	}
	return resp.GetData()[0].GetPriceScheduleId(), nil
}

// NewTabAction handles GET /action/price-schedule/{id}/plan/{ppid}/tab/{tab}.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")
		tab := deps.ScheduleLabels.Tabs.CanonicalizeTab(viewCtx.Request.PathValue("tab"))
		if tab == "" {
			tab = "info"
		}

		pageData, err := buildPageData(ctx, deps, sid, ppid, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}

		templateName := "price-schedule-plan-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		return view.OK(templateName, pageData)
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
					pricingLockedReason = deps.PlanLabels.Messages.PricingLockedReason
				}
			}

			// Populate new billing fields from existing record.
			billingCycleValue := ""
			if v := pp.GetBillingCycleValue(); v > 0 {
				billingCycleValue = fmt.Sprintf("%d", v)
			}
			defaultTermValue := ""
			if v := pp.GetDefaultTermValue(); v > 0 {
				defaultTermValue = fmt.Sprintf("%d", v)
			}
			formLabels := deps.PlanLabels.Form
			scheduleName, scheduleClientID := lookupScheduleNameAndClient(ctx, deps, sid)
			// 2026-05-03 — same mutually-exclusive plan-scope filter as the
			// schedule-add drawer: client-scoped schedule shows only matching
			// client's plans, master schedule shows only master plans.
			planOpts := buildPlanOptions(ctx, deps, pp.GetPlanId(), scheduleClientID)
			return view.OK("price-plan-drawer-form", &form.Data{
				FormAction:             route.ResolveURL(deps.Routes.PlanEditURL, "id", sid, "ppid", ppid),
				IsEdit:                 true,
				Context:                form.ContextSchedule,
				ID:                     ppid,
				PlanID:                 pp.GetPlanId(),
				ScheduleID:             sid,
				ScheduleName:           scheduleName,
				ParentScheduleClientID: scheduleClientID,
				Name:                  pp.GetName(),
				Description:           pp.GetDescription(),
				Amount:                strconv.FormatFloat(float64(pp.GetBillingAmount())/100.0, 'f', 2, 64),
				Currency:              pp.GetBillingCurrency(),
				DurationValue:         fmt.Sprintf("%d", pp.GetDurationValue()),
				DurationUnit:          pp.GetDurationUnit(),
				Active:                pp.GetActive(),
				// Wave 2: populate new billing fields.
				BillingKind:         pp.GetBillingKind().String(),
				AmountBasis:         pp.GetAmountBasis().String(),
				BillingCycleValue:   billingCycleValue,
				BillingCycleUnit:    pp.GetBillingCycleUnit(),
				TermValue:           defaultTermValue,
				TermUnit:            pp.GetDefaultTermUnit(),
				DurationUnitOptions: buildDurationUnitOptions(deps.CommonLabels),
				PlanOptions:           planOpts,
				SelectedPlanID:        pp.GetPlanId(),
				SelectedPlanLabel:     labelFromOptions(planOpts, pp.GetPlanId()),
				InUse:                 pricingLocked,
				LockMessage:           pricingLockedReason,
				Labels:                form.LabelsFromPricePlan(formLabels),
				CommonLabels:          deps.CommonLabels,
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
		dvStr := r.FormValue("duration_value")
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
		var existingDurationValue *int32
		var existingDurationUnit *string
		if deps.GetPricePlanInUseIDs != nil && existing != nil && len(existing.GetData()) > 0 {
			inUseMap, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppid})
			if inUseMap[ppid] {
				ex := existing.GetData()[0]
				amount = ex.GetBillingAmount()
				currency = ex.GetBillingCurrency()
				existingDurationValue = ex.DurationValue
				existingDurationUnit = ex.DurationUnit
				dvStr = ""        // sentinel: prefer existingDurationValue below
				durationUnit = "" // sentinel: prefer existingDurationUnit below
			}
		}

		// Wave 2: new billing semantics fields.
		bcvStr := r.FormValue("billing_cycle_value")
		bcv, _ := strconv.ParseInt(bcvStr, 10, 32)
		bcu := r.FormValue("billing_cycle_unit")
		dtvStr := r.FormValue("default_term_value")
		dtv, _ := strconv.ParseInt(dtvStr, 10, 32)
		dtu := r.FormValue("default_term_unit")
		billingKindStr := r.FormValue("billing_kind")
		amountBasisStr := r.FormValue("amount_basis")

		planPageName := r.FormValue("name")
		planPageDesc := r.FormValue("description")
		req := &priceplanpb.UpdatePricePlanRequest{
			Data: &priceplanpb.PricePlan{
				Id:              ppid,
				PlanId:          r.FormValue("plan_id"),
				Name:            &planPageName,
				Description:     &planPageDesc,
				BillingAmount:   amount,
				BillingCurrency: currency,
				Active:          active,
			},
		}
		req.Data.PriceScheduleId = &sid
		// Phase 1 legacy dual-write — proto fields now optional. Prefer the
		// in-use snapshot when locked; otherwise read from form input.
		if existingDurationValue != nil {
			req.Data.DurationValue = existingDurationValue
		} else if dvStr != "" {
			if parsed, err := strconv.ParseInt(dvStr, 10, 32); err == nil {
				dv32 := int32(parsed)
				req.Data.DurationValue = &dv32
			}
		}
		if existingDurationUnit != nil {
			req.Data.DurationUnit = existingDurationUnit
		} else if durationUnit != "" {
			req.Data.DurationUnit = &durationUnit
		}
		// Set new enum fields.
		if billingKindStr != "" {
			if bk, ok := priceplanpb.BillingKind_value[billingKindStr]; ok {
				req.Data.BillingKind = priceplanpb.BillingKind(bk)
			}
		}
		if amountBasisStr != "" {
			if ab, ok := priceplanpb.AmountBasis_value[amountBasisStr]; ok {
				req.Data.AmountBasis = priceplanpb.AmountBasis(ab)
			}
		}
		// Set new optional duration fields.
		if bcvStr != "" {
			bcv32 := int32(bcv)
			req.Data.BillingCycleValue = &bcv32
		}
		if bcu != "" {
			req.Data.BillingCycleUnit = &bcu
		}
		if dtvStr != "" {
			dtv32 := int32(dtv)
			req.Data.DefaultTermValue = &dtv32
		}
		if dtu != "" {
			req.Data.DefaultTermUnit = &dtu
		}
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
			return centymo.HTMXError(deps.PlanLabels.Messages.CreateNotAvailable)
		}
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")

		pplLabels := deps.ProductPricePlanLabels.Form
		if viewCtx.Request.Method == http.MethodGet {
			planName, planDesc := lookupPackageNameDesc(ctx, deps, ppid)
			parent, _ := loadParentContext(ctx, deps, ppid)
			currency := parent.Currency
			if currency == "" {
				currency = "PHP"
			}
			showTreatment := parent.BillingKind != "BILLING_KIND_ONE_TIME"
			return view.OK("price-schedule-plan-product-price-drawer", &ProductPriceFormData{
				FormAction:              route.ResolveURL(deps.Routes.PlanProductPriceAddURL, "id", sid, "ppid", ppid),
				ScheduleID:              sid,
				PricePlanID:             ppid,
				Currency:                currency,
				CommonLabels:            deps.CommonLabels,
				PlanName:                planName,
				PlanDescription:         planDesc,
				ParentBillingKind:       parent.BillingKind,
				ParentAmountBasis:       parent.AmountBasis,
				ShowTreatment:           showTreatment,
				BasisBannerMessage:      basisBannerMessage(parent.AmountBasis, deps.ScheduleLabels.Detail),
				BillingKindDisplay:      parent.BillingKindDisplay,
				AmountBasisDisplay:      parent.AmountBasisDisplay,
				BillingCycleDisplay:     parent.BillingCycleDisplay,
				TermDisplay:             parent.TermDisplay,
				ParentCurrencyDisplay:   parent.ParentCurrencyDisplay,
				RateCardName:            parent.RateCardName,
				ProductPricePlanLabels:  pplLabels,
				PriceScheduleLabels:     deps.ScheduleLabels.Detail,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		productPlanID := viewCtx.Request.FormValue("product_plan_id")
		if productPlanID == "" {
			// Backward-compatible: old form posts may still send product_id;
			// resolve to a product_plan_id on this plan when possible.
			if legacyProductID := viewCtx.Request.FormValue("product_id"); legacyProductID != "" {
				productPlanID = resolveProductPlanIDForProduct(ctx, deps, ppid, legacyProductID)
			}
		}
		if productPlanID == "" {
			return centymo.HTMXError(deps.PlanLabels.Messages.ProductRequired)
		}
		priceCentavos, ok := parsePriceCentavos(viewCtx.Request.FormValue("price"))
		if !ok {
			return centymo.HTMXError(deps.PlanLabels.Messages.InvalidPrice)
		}
		currency := viewCtx.Request.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}
		dateStart := viewCtx.Request.FormValue("date_start")
		dateEnd := viewCtx.Request.FormValue("date_end")
		billingTreatment := viewCtx.Request.FormValue("billing_treatment")
		parent, _ := loadParentContext(ctx, deps, ppid)
		// Currency must match parent PricePlan.billing_currency (proto invariant).
		if parent.Currency != "" && currency != parent.Currency {
			return centymo.HTMXError(deps.PlanLabels.Messages.CurrencyMismatch)
		}
		// billing_treatment is meaningless when parent has no cycles. Drop the
		// posted value so we never persist a stale treatment on a ONE_TIME plan.
		if parent.BillingKind == "BILLING_KIND_ONE_TIME" {
			billingTreatment = ""
		}
		record := &productpriceplanpb.ProductPricePlan{
			PricePlanId:     ppid,
			ProductPlanId:   productPlanID,
			BillingAmount:   priceCentavos,
			BillingCurrency: currency,
			Active:          true,
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
			return centymo.HTMXError(deps.PlanLabels.Messages.UpdateNotAvailable)
		}
		sid := viewCtx.Request.PathValue("id")
		ppid := viewCtx.Request.PathValue("ppid")
		pppid := viewCtx.Request.PathValue("pppid")

		existing, err := findProductPricePlan(ctx, deps, pppid)
		if err != nil {
			return centymo.HTMXError(err.Error())
		}

		pplLabels := deps.ProductPricePlanLabels.Form
		if viewCtx.Request.Method == http.MethodGet {
			parent, _ := loadParentContext(ctx, deps, ppid)
			currency := existing.GetBillingCurrency()
			if currency == "" {
				currency = parent.Currency
			}
			if currency == "" {
				currency = "PHP"
			}
			showTreatment := parent.BillingKind != "BILLING_KIND_ONE_TIME"
			planName, planDesc := lookupPackageNameDesc(ctx, deps, ppid)
			// Model D — resolve product + variant via the referenced ProductPlan row.
			existingProductPlanID := existing.GetProductPlanId()
			prodName, prodDesc, variantName := lookupProductPlanDisplay(ctx, deps, existingProductPlanID)

			pricingLocked := false
			pricingLockedReason := ""
			if deps.GetPricePlanInUseIDs != nil {
				if inUse, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppid}); inUse[ppid] {
					pricingLocked = true
					pricingLockedReason = deps.PlanLabels.Messages.ItemPricingLockedReason
				}
			}

			return view.OK("price-schedule-plan-product-price-drawer", &ProductPriceFormData{
				FormAction:              route.ResolveURL(deps.Routes.PlanProductPriceEditURL, "id", sid, "ppid", ppid, "pppid", pppid),
				IsEdit:                  true,
				ID:                      pppid,
				ScheduleID:              sid,
				PricePlanID:             ppid,
				ProductPlanID:           existingProductPlanID,
				Price:                   fmt.Sprintf("%.2f", float64(existing.GetBillingAmount())/100.0),
				Currency:                currency,
				CommonLabels:            deps.CommonLabels,
				PlanName:                planName,
				PlanDescription:         planDesc,
				ProductName:             prodName,
				ProductDescription:      prodDesc,
				VariantName:             variantName,
				PricingLocked:           pricingLocked,
				PricingLockedReason:     pricingLockedReason,
				// Wave 2: populate billing treatment and dates from existing record.
				BillingTreatment:        existing.GetBillingTreatment().String(),
				DateStart:               existing.GetDateStart(),
				DateEnd:                 existing.GetDateEnd(),
				ParentBillingKind:       parent.BillingKind,
				ParentAmountBasis:       parent.AmountBasis,
				ShowTreatment:           showTreatment,
				BasisBannerMessage:      basisBannerMessage(parent.AmountBasis, deps.ScheduleLabels.Detail),
				BillingKindDisplay:      parent.BillingKindDisplay,
				AmountBasisDisplay:      parent.AmountBasisDisplay,
				BillingCycleDisplay:     parent.BillingCycleDisplay,
				TermDisplay:             parent.TermDisplay,
				ParentCurrencyDisplay:   parent.ParentCurrencyDisplay,
				RateCardName:            parent.RateCardName,
				ProductPricePlanLabels:  pplLabels,
				PriceScheduleLabels:     deps.ScheduleLabels.Detail,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.PlanLabels.Errors.Unauthorized)
		}
		// Server-side lock enforcement: if the parent PricePlan is in use by an
		// active subscription, reject price/currency changes (client may have
		// bypassed the disabled inputs).
		if deps.GetPricePlanInUseIDs != nil {
			if inUse, _ := deps.GetPricePlanInUseIDs(ctx, []string{ppid}); inUse[ppid] {
				return centymo.HTMXError(deps.PlanLabels.Messages.InUseCannotModify)
			}
		}
		// The catalog-line assignment is display-only in the drawer — preserve
		// the existing product_plan_id.
		priceCentavos, ok := parsePriceCentavos(viewCtx.Request.FormValue("price"))
		if !ok {
			return centymo.HTMXError(deps.PlanLabels.Messages.InvalidPrice)
		}
		currency := viewCtx.Request.FormValue("currency")
		if currency == "" {
			currency = "PHP"
		}
		dateStart := viewCtx.Request.FormValue("date_start")
		dateEnd := viewCtx.Request.FormValue("date_end")
		billingTreatment := viewCtx.Request.FormValue("billing_treatment")
		parent, _ := loadParentContext(ctx, deps, ppid)
		if parent.Currency != "" && currency != parent.Currency {
			return centymo.HTMXError(deps.PlanLabels.Messages.CurrencyMismatch)
		}
		if parent.BillingKind == "BILLING_KIND_ONE_TIME" {
			billingTreatment = ""
		}
		updated := &productpriceplanpb.ProductPricePlan{
			Id:              pppid,
			PricePlanId:     ppid,
			ProductPlanId:   existing.GetProductPlanId(),
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
			return centymo.HTMXError(deps.PlanLabels.Messages.DeleteNotAvailable)
		}
		_ = viewCtx.Request.ParseForm()
		pppid := viewCtx.Request.FormValue("id")
		if pppid == "" {
			pppid = viewCtx.Request.URL.Query().Get("id")
		}
		if pppid == "" {
			return centymo.HTMXError(deps.PlanLabels.Messages.IDRequired)
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

	currency := pp.GetBillingCurrency()
	if currency == "" {
		currency = "PHP"
	}
	amountCell := types.MoneyCell(float64(pp.GetBillingAmount()), currency, true)

	duration := ""
	if dv := pp.GetDurationValue(); dv > 0 {
		duration = pyeza.FormatDuration(dv, pp.GetDurationUnit(), deps.CommonLabels.DurationUnit)
	}

	status := "active"
	statusVariant := "success"
	if !pp.GetActive() {
		status = "inactive"
		statusVariant = "warning"
	}

	scheduleName := lookupScheduleName(ctx, deps, sid)
	scheduleBack := route.ResolveURL(deps.Routes.DetailURL, "id", sid) + "?tab=" + deps.ScheduleLabels.Tabs.ResolveTabSlug("pricePlan")

	// Plan-scoped mount overrides — keep the page anchored to Services > Packages
	// when invoked under /app/plans/detail/{id}/price/{ppid}. The breadcrumb
	// resolves to the loaded price_plan's plan_id so the back link survives
	// even when the schedule is missing.
	breadcrumbLabel := scheduleName
	breadcrumbURL := scheduleBack
	if deps.PlanDetailBackURL != "" {
		back := route.ResolveURL(deps.PlanDetailBackURL, "id", pp.GetPlanId())
		if deps.PlanDetailBackTab != "" {
			back += "?tab=" + deps.PlanDetailBackTab
		}
		breadcrumbURL = back
	}
	activeNav := deps.Routes.ActiveNav
	if deps.ActiveNavOverride != "" {
		activeNav = deps.ActiveNavOverride
	}
	activeSubNav := deps.Routes.ActiveSubNav
	if deps.ActiveSubNavOverride != "" {
		activeSubNav = deps.ActiveSubNavOverride
	}

	// Fallback to linked Plan's name/description when price_plan values are blank —
	// mirrors the rate-card packages-tab table convention.
	planName, planDesc := lookupPlanNameDesc(ctx, deps, pp.GetPlanId())
	if deps.PlanDetailBackURL != "" && planName != "" {
		breadcrumbLabel = planName
	}
	effectiveName := strings.TrimSpace(pp.GetName())
	if effectiveName == "" {
		effectiveName = planName
	}
	effectiveDesc := strings.TrimSpace(pp.GetDescription())
	if effectiveDesc == "" {
		effectiveDesc = planDesc
	}

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

	// Subscriptions count for tab badge — sourced from the new
	// ListSubscriptionsByPricePlan use case so the value matches what the tab
	// table will render. When the dep is unwired the badge stays at 0.
	subscriptionCount := countSubscriptionsForPricePlan(ctx, deps, ppid)

	// Tab item URLs reflect the active mount: when plan-scoped, they live under
	// /app/plans/detail/{plan_id}/price/{ppid}; otherwise the rate-card-scoped
	// defaults derived from PriceScheduleRoutes.
	detailURLTemplate := deps.Routes.PlanDetailURL
	tabActionURLTemplate := deps.Routes.PlanTabActionURL
	pathID := sid
	if deps.PlanScopedDetailURL != "" {
		detailURLTemplate = deps.PlanScopedDetailURL
		// Plan-scoped routes use plan_id as {id}; resolve from the loaded price_plan.
		pathID = pp.GetPlanId()
	}
	if deps.PlanScopedTabActionURL != "" {
		tabActionURLTemplate = deps.PlanScopedTabActionURL
	}
	base := route.ResolveURL(detailURLTemplate, "id", pathID, "ppid", ppid)
	action := route.ResolveURL(tabActionURLTemplate, "id", pathID, "ppid", ppid, "tab", "")
	productPricesSlug := deps.ScheduleLabels.Tabs.ResolveTabSlug("product-prices")
	subscriptionsSlug := deps.ScheduleLabels.Tabs.ResolveTabSlug("subscriptions")
	subscriptionsLabel := deps.ScheduleLabels.Tabs.Subscriptions
	if subscriptionsLabel == "" {
		subscriptionsLabel = deps.PlanLabels.Tabs.Subscriptions
	}
	if subscriptionsLabel == "" {
		subscriptionsLabel = "Subscriptions"
	}
	attachmentsLabel := deps.ScheduleLabels.Detail.TabAttachments
	if attachmentsLabel == "" {
		attachmentsLabel = deps.PlanLabels.Detail.AttachmentsTab
	}
	if attachmentsLabel == "" {
		attachmentsLabel = "Attachments"
	}
	tabItems := []pyeza.TabItem{
		{Key: "info", Label: deps.ScheduleLabels.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "product-prices", Label: deps.ScheduleLabels.Tabs.ProductPrices, Href: base + "?tab=" + productPricesSlug, HxGet: action + productPricesSlug, Icon: "icon-package", Count: count},
		{Key: "subscriptions", Label: subscriptionsLabel, Href: base + "?tab=" + subscriptionsSlug, HxGet: action + subscriptionsSlug, Icon: "icon-briefcase", Count: subscriptionCount},
		{Key: "attachments", Label: attachmentsLabel, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
	}

	headerSubtitle := effectiveDesc
	if headerSubtitle == "" {
		headerSubtitle = deps.ScheduleLabels.Detail.NoDescriptionSubtitle
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:        viewCtx.CacheVersion,
			Title:               effectiveName,
			CurrentPath:         viewCtx.CurrentPath,
			ActiveNav:           activeNav,
			ActiveSubNav:        activeSubNav,
			HeaderTitle:         effectiveName,
			HeaderSubtitle:      headerSubtitle,
			HeaderBreadcrumb:    breadcrumbLabel,
			HeaderBreadcrumbURL: breadcrumbURL,
			HeaderIcon:          "icon-tag",
			CommonLabels:        deps.CommonLabels,
		},
		ContentTemplate:        "price-schedule-plan-detail-content",
		ScheduleID:             sid,
		ScheduleName:           scheduleName,
		ScheduleBackURL:        breadcrumbURL,
		PricePlan:              pp,
		Labels:                 deps.PlanLabels,
		ActiveTab:              activeTab,
		TabItems:               tabItems,
		ID:                     ppid,
		Name:                   effectiveName,
		Description:            effectiveDesc,
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
	if activeTab == "subscriptions" {
		rows := loadSubscriptionsForPricePlan(ctx, deps, ppid)
		pageData.SubscriptionsTable = buildSubscriptionsTable(ctx, deps, sid, ppid, rows)
	}
	if activeTab == "attachments" && deps.ListAttachments != nil {
		cfg := attachmentConfig(deps)
		var attachItems []*attachmentpb.Attachment
		if resp, err := deps.ListAttachments(ctx, cfg.EntityType, ppid); err == nil && resp != nil {
			attachItems = resp.GetData()
		}
		pageData.AttachmentTable = attachment.BuildTable(attachItems, cfg, sid, "ppid", ppid)
	}
	return pageData, nil
}

// countSubscriptionsForPricePlan returns the count of active subscriptions
// referencing the given PricePlan. The use case is the same one the tab body
// renders against — so the badge count and row count cannot drift. Returns 0
// (no badge) when the dep is unwired.
func countSubscriptionsForPricePlan(ctx context.Context, deps *DetailViewDeps, pricePlanID string) int {
	if deps.ListSubscriptionsByPricePlan == nil {
		return 0
	}
	activeOnly := true
	resp, err := deps.ListSubscriptionsByPricePlan(ctx, &subscriptionpb.ListSubscriptionsByPricePlanRequest{
		PricePlanId: pricePlanID,
		ActiveOnly:  &activeOnly,
	})
	if err != nil {
		log.Printf("Failed to count subscriptions for price plan %s: %v", pricePlanID, err)
		return 0
	}
	return len(resp.GetSubscriptionList())
}

// loadSubscriptionsForPricePlan fetches active subscriptions for the given
// PricePlan and shapes them into SubscriptionRow values for the tab table.
// Hydration of Client + PricePlan + Plan is provided by the espyna use case
// (single CTE-based JOIN) — the view layer does not chain N+1 lookups.
func loadSubscriptionsForPricePlan(ctx context.Context, deps *DetailViewDeps, pricePlanID string) []SubscriptionRow {
	if deps.ListSubscriptionsByPricePlan == nil {
		return nil
	}
	activeOnly := true
	resp, err := deps.ListSubscriptionsByPricePlan(ctx, &subscriptionpb.ListSubscriptionsByPricePlanRequest{
		PricePlanId: pricePlanID,
		ActiveOnly:  &activeOnly,
	})
	if err != nil {
		log.Printf("Failed to load subscriptions for price plan %s: %v", pricePlanID, err)
		return nil
	}

	tz := types.LocationFromContext(ctx)
	rows := make([]SubscriptionRow, 0, len(resp.GetSubscriptionList()))
	for _, s := range resp.GetSubscriptionList() {
		if s == nil {
			continue
		}
		clientName := ""
		clientID := s.GetClientId()
		if c := s.GetClient(); c != nil {
			clientName = c.GetName()
			if clientName == "" {
				if u := c.GetUser(); u != nil {
					first := u.GetFirstName()
					last := u.GetLastName()
					if first != "" || last != "" {
						clientName = strings.TrimSpace(first + " " + last)
					}
					if clientName == "" {
						clientName = u.GetEmailAddress()
					}
				}
			}
			if clientID == "" {
				clientID = c.GetId()
			}
		}

		planName := ""
		if pp := s.GetPricePlan(); pp != nil {
			if p := pp.GetPlan(); p != nil {
				planName = p.GetName()
			}
			if planName == "" {
				planName = pp.GetName()
			}
		}

		rows = append(rows, SubscriptionRow{
			ID:         s.GetId(),
			Name:       s.GetName(),
			ClientID:   clientID,
			ClientName: clientName,
			Plan:       planName,
			DateStart:  types.FormatTimestampInTZ(s.GetDateTimeStart(), tz, types.DateTimeReadable),
			DateEnd:    types.FormatTimestampInTZ(s.GetDateTimeEnd(), tz, types.DateTimeReadable),
		})
	}
	return rows
}

// buildSubscriptionsTable assembles the TableConfig for the schedule-scoped
// plan detail's "Engagements"/"Subscriptions" tab. Columns:
// Name → Client → Plan → Start Date → End Date. The View action targets the
// nested engagement URL when configured so the breadcrumb chains
// rate-card → plan → engagement.
func buildSubscriptionsTable(ctx context.Context, deps *DetailViewDeps, sid, ppid string, rows []SubscriptionRow) *types.TableConfig {
	perms := view.GetUserPermissions(ctx)
	subLabels := deps.PlanLabels.Detail.Subscriptions

	columns := []types.TableColumn{
		{Key: "name", Label: subLabels.ColumnName},
		{Key: "client", Label: subLabels.ColumnClient},
		{Key: "plan", Label: subLabels.ColumnPlan},
		{Key: "start_date", Label: subLabels.ColumnStartDate, WidthClass: "col-3xl"},
		{Key: "end_date", Label: subLabels.ColumnEndDate, WidthClass: "col-3xl"},
	}

	tableRows := make([]types.TableRow, 0, len(rows))
	for _, r := range rows {
		viewURL := ""
		if deps.PlanEngagementDetailURL != "" {
			viewURL = route.ResolveURL(deps.PlanEngagementDetailURL, "id", sid, "ppid", ppid, "eid", r.ID)
		} else if deps.SubscriptionDetailURL != "" {
			viewURL = route.ResolveURL(deps.SubscriptionDetailURL, "id", r.ID)
		}

		actions := []types.TableAction{}
		if viewURL != "" {
			actions = append(actions, types.TableAction{Type: "view", Label: deps.CommonLabels.Actions.View, Action: "view", Href: viewURL})
		}
		if perms.Can("subscription", "update") && deps.SubscriptionEditURL != "" {
			editURL := route.ResolveURL(deps.SubscriptionEditURL, "id", r.ID)
			actions = append(actions, types.TableAction{Type: "edit", Label: deps.CommonLabels.Actions.Edit, Action: "edit", URL: editURL, DrawerTitle: r.Name})
		}
		if perms.Can("subscription", "delete") && deps.SubscriptionDeleteURL != "" {
			actions = append(actions, types.TableAction{
				Type:           "delete",
				Label:          deps.CommonLabels.Actions.Delete,
				Action:         "delete",
				URL:            deps.SubscriptionDeleteURL,
				ItemName:       r.Name,
				ConfirmTitle:   subLabels.ConfirmDeleteTitle,
				ConfirmMessage: fmt.Sprintf(subLabels.ConfirmDeleteMessage, r.Name),
			})
		}

		tableRows = append(tableRows, types.TableRow{
			ID: r.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: r.Name},
				{Type: "text", Value: r.ClientName},
				{Type: "text", Value: r.Plan},
				{Type: "text", Value: r.DateStart},
				{Type: "text", Value: r.DateEnd},
			},
			DataAttrs: map[string]string{
				"name":   r.Name,
				"client": r.ClientName,
				"plan":   r.Plan,
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, tableRows)

	tc := &types.TableConfig{
		ID:                   "subscriptions-table",
		Columns:              columns,
		Rows:                 tableRows,
		Labels:               deps.TableLabels,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		EmptyState: types.TableEmptyState{
			Title:   subLabels.EmptyTitle,
			Message: subLabels.EmptyMessage,
		},
	}
	types.ApplyTableSettings(tc)
	return tc
}

func buildProductPricesTable(ctx context.Context, deps *DetailViewDeps, sid, ppid string) *types.TableConfig {
	perms := view.GetUserPermissions(ctx)
	l := deps.PlanLabels

	pplLabels := deps.ProductPricePlanLabels.Form
	parent, _ := loadParentContext(ctx, deps, ppid)
	showTreatment := parent.BillingKind != "BILLING_KIND_ONE_TIME"

	columns := []types.TableColumn{
		{Key: "product", Label: deps.ScheduleLabels.Detail.ProductPriceColumnProduct},
		{Key: "price", Label: deps.ScheduleLabels.Detail.ProductPriceColumnPrice, WidthClass: "col-4xl", Align: "right"},
		{Key: "currency", Label: deps.ScheduleLabels.Detail.ProductPriceColumnCurrency, NoSort: true, WidthClass: "col-2xl"},
	}
	if showTreatment {
		columns = append(columns, types.TableColumn{Key: "treatment", Label: deps.ScheduleLabels.Detail.ProductPriceColumnTreatment, NoSort: true, WidthClass: "col-3xl"})
	}
	columns = append(columns, types.TableColumn{Key: "effective", Label: deps.ScheduleLabels.Detail.ProductPriceColumnEffective, NoSort: true, WidthClass: "col-4xl"})

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
	// resolve row display via the catalog line's FK.
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

	refreshURL := route.ResolveURL(deps.Routes.PlanTabActionURL, "id", sid, "ppid", ppid, "tab", deps.ScheduleLabels.Tabs.ResolveTabSlug("product-prices"))
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
				ref := productPlans[item.GetProductPlanId()]
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
					productName = fmt.Sprintf("%s (%s)", productName, ref.variantID)
				}
				itemCurrency := item.GetBillingCurrency()
				if itemCurrency == "" {
					itemCurrency = "PHP"
				}
				priceCell := types.MoneyCell(float64(item.GetBillingAmount()), itemCurrency, true)
				cells := []types.TableCell{
					{Type: "text", Value: productName},
					priceCell,
					{Type: "text", Value: itemCurrency},
				}
				if showTreatment {
					cells = append(cells, types.TableCell{Type: "text", Value: billingTreatmentDisplay(item.GetBillingTreatment().String(), pplLabels)})
				}
				cells = append(cells, types.TableCell{Type: "text", Value: effectiveRangeDisplay(item.GetDateStart(), item.GetDateEnd())})
				rows = append(rows, types.TableRow{
					ID:    itemID,
					Cells: cells,
					// No delete action: rows are auto-seeded from product_plan assignments,
					// so deletion here would desync the two tables. Use the plan's
					// Products tab to remove the product_plan link, which in turn
					// should remove its product_price_plan rows.
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
		// No PrimaryAction: product_price_plan rows are auto-seeded from product_plan
		// assignments when the parent PricePlan is created, so manual Add is disabled here —
		// users Edit existing rows instead.
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
	parent, ok := loadParentContext(ctx, deps, pricePlanID)
	if !ok || parent.Currency == "" {
		return "PHP"
	}
	return parent.Currency
}

// parentPricePlanContext captures the parent PricePlan fields the PPP drawer
// needs to know about: currency (locks the per-line currency), billing_kind
// (decides whether billing_treatment renders), amount_basis (drives the
// banner explaining what the line prices mean), and pre-formatted display
// strings for the read-only context block above the editable fields.
type parentPricePlanContext struct {
	Currency    string
	BillingKind string
	AmountBasis string

	// Display strings — empty when the corresponding source data is missing.
	BillingKindDisplay    string
	AmountBasisDisplay    string
	BillingCycleDisplay   string
	TermDisplay           string
	ParentCurrencyDisplay string
	RateCardName          string
}

func loadParentContext(ctx context.Context, deps *DetailViewDeps, pricePlanID string) (parentPricePlanContext, bool) {
	if deps.ReadPricePlan == nil || pricePlanID == "" {
		return parentPricePlanContext{}, false
	}
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: pricePlanID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return parentPricePlanContext{}, false
	}
	pp := resp.GetData()[0]
	pc := parentPricePlanContext{
		Currency:              pp.GetBillingCurrency(),
		BillingKind:           pp.GetBillingKind().String(),
		AmountBasis:           pp.GetAmountBasis().String(),
		ParentCurrencyDisplay: pp.GetBillingCurrency(),
		BillingKindDisplay:    formatBillingKindLabel(pp.GetBillingKind().String(), deps.PlanLabels.Form),
		AmountBasisDisplay:    formatAmountBasisLabel(pp.GetAmountBasis().String(), deps.PlanLabels.Form),
	}
	if v := pp.GetBillingCycleValue(); v > 0 {
		pc.BillingCycleDisplay = pyeza.FormatDuration(v, pp.GetBillingCycleUnit(), deps.CommonLabels.DurationUnit)
	}
	if v := pp.GetDefaultTermValue(); v > 0 {
		pc.TermDisplay = pyeza.FormatDuration(v, pp.GetDefaultTermUnit(), deps.CommonLabels.DurationUnit)
	}
	if scheduleID := pp.GetPriceScheduleId(); scheduleID != "" && deps.ReadPriceSchedule != nil {
		if schedResp, err := deps.ReadPriceSchedule(ctx, &priceschedulepb.ReadPriceScheduleRequest{
			Data: &priceschedulepb.PriceSchedule{Id: scheduleID},
		}); err == nil && len(schedResp.GetData()) > 0 {
			pc.RateCardName = schedResp.GetData()[0].GetName()
		}
	}
	return pc, true
}

func formatBillingKindLabel(kind string, l centymo.PricePlanFormLabels) string {
	switch kind {
	case "BILLING_KIND_ONE_TIME":
		return l.BillingKindOneTime
	case "BILLING_KIND_RECURRING":
		return l.BillingKindRecurring
	case "BILLING_KIND_CONTRACT":
		return l.BillingKindContract
	}
	return kind
}

func formatAmountBasisLabel(basis string, l centymo.PricePlanFormLabels) string {
	switch basis {
	case "AMOUNT_BASIS_PER_CYCLE":
		return l.AmountBasisPerCycle
	case "AMOUNT_BASIS_TOTAL_PACKAGE":
		return l.AmountBasisTotalPackage
	case "AMOUNT_BASIS_DERIVED_FROM_LINES":
		return l.AmountBasisDerivedFromLines
	}
	return basis
}

// billingTreatmentDisplay maps the proto enum string to its human label.
// Returns "—" for unspecified so the table cell stays visually quiet.
func billingTreatmentDisplay(value string, l centymo.ProductPricePlanFormLabels) string {
	switch value {
	case "BILLING_TREATMENT_RECURRING":
		return l.BillingTreatmentRecurring
	case "BILLING_TREATMENT_ONE_TIME_INITIAL":
		return l.BillingTreatmentOneTimeInitial
	case "BILLING_TREATMENT_USAGE_BASED":
		return l.BillingTreatmentUsageBased
	}
	return "—"
}

// effectiveRangeDisplay renders the per-line effective dates as "start → end",
// "from start", "until end", or "Always" when both ends are empty.
func effectiveRangeDisplay(start, end string) string {
	switch {
	case start == "" && end == "":
		return "Always"
	case start != "" && end == "":
		return "from " + start
	case start == "" && end != "":
		return "until " + end
	default:
		return start + " → " + end
	}
}

// basisBannerMessage returns a one-line explanation for the user about the
// relationship between the parent PricePlan's amount_basis and the per-line
// price they're editing. Sourced from PriceScheduleDetailLabels so tier
// overrides flow through lyngua.
func basisBannerMessage(amountBasis string, l centymo.PriceScheduleDetailLabels) string {
	switch amountBasis {
	case "AMOUNT_BASIS_DERIVED_FROM_LINES":
		return l.BasisBannerDerived
	case "AMOUNT_BASIS_TOTAL_PACKAGE":
		return l.BasisBannerTotalPackage
	case "AMOUNT_BASIS_PER_CYCLE":
		return l.BasisBannerPerCycle
	}
	return ""
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

// buildPlanOptions returns the plan picker options with a strict, mutually-
// exclusive client-scope filter mirroring the schedule-add drawer:
//
//   - Client-scoped schedule (scheduleClientID != ""): only plans whose
//     client_id matches scheduleClientID. Master plans are excluded.
//   - Master schedule (scheduleClientID == ""): only master plans
//     (client_id empty). Client-scoped plans cannot attach to a master schedule.
//
// `selectedID` is preserved for edit-mode drawer pre-selection regardless of
// the filter (the row is already wired; we never silently drop the operator's
// existing selection).
func buildPlanOptions(ctx context.Context, deps *DetailViewDeps, selectedID, scheduleClientID string) []map[string]any {
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
			if planClientID != scheduleClientID && p.GetId() != selectedID {
				continue
			}
		} else {
			if planClientID != "" && p.GetId() != selectedID {
				continue
			}
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
	name, _ := lookupScheduleNameAndClient(ctx, deps, scheduleID)
	return name
}

// lookupScheduleNameAndClient returns the schedule's display name + client_id.
// Used by edit handlers that need the client_id to scope-filter the plan picker.
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

// lookupPlanNameDesc returns the linked Plan's name and description (trimmed).
// Used as fallback when price_plan.Name / price_plan.Description are blank.
func lookupPlanNameDesc(ctx context.Context, deps *DetailViewDeps, planID string) (string, string) {
	if planID == "" || deps.ListPlans == nil {
		return "", ""
	}
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return "", ""
	}
	for _, p := range resp.GetData() {
		if p == nil || p.GetId() != planID {
			continue
		}
		return strings.TrimSpace(p.GetName()), strings.TrimSpace(p.GetDescription())
	}
	return "", ""
}

// lookupProductNameDesc reads the Product and returns its trimmed name + description.
func lookupProductNameDesc(ctx context.Context, deps *DetailViewDeps, productID string) (string, string) {
	if productID == "" || deps.ListProducts == nil {
		return "", ""
	}
	resp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		return "", ""
	}
	for _, p := range resp.GetData() {
		if p == nil || p.GetId() != productID {
			continue
		}
		return strings.TrimSpace(p.GetName()), strings.TrimSpace(p.GetDescription())
	}
	return "", ""
}

// lookupProductPlanDisplay resolves product name, description, and (optional)
// variant SKU for a ProductPlan.id — used to render the read-only context
// rows on the schedule-nested product-price drawer under Model D.
func lookupProductPlanDisplay(ctx context.Context, deps *DetailViewDeps, productPlanID string) (name, desc, variant string) {
	if productPlanID == "" || deps.ListProductPlans == nil {
		return "", "", ""
	}
	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
	if err != nil {
		return "", "", ""
	}
	var (
		productID string
		variantID string
	)
	for _, pp := range ppResp.GetData() {
		if pp != nil && pp.GetId() == productPlanID {
			productID = pp.GetProductId()
			variantID = pp.GetProductVariantId()
			break
		}
	}
	name, desc = lookupProductNameDesc(ctx, deps, productID)
	variant = variantID // caller may translate via SKU lookup if needed
	return name, desc, variant
}

// resolveProductPlanIDForProduct finds the ProductPlan row in the parent
// Plan of the given PricePlan that references the supplied product_id. Used
// as a transitional fallback so old form posts (which send product_id) still
// work while Model D rolls out.
func resolveProductPlanIDForProduct(ctx context.Context, deps *DetailViewDeps, pricePlanID, productID string) string {
	if productID == "" || deps.ListProductPlans == nil {
		return ""
	}
	planID := loadPricePlanPlanID(ctx, deps, pricePlanID)
	if planID == "" {
		return ""
	}
	ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{})
	if err != nil {
		return ""
	}
	for _, pp := range ppResp.GetData() {
		if pp != nil && pp.GetPlanId() == planID && pp.GetProductId() == productID {
			return pp.GetId()
		}
	}
	return ""
}

// lookupPackageNameDesc resolves the display name + description for a price_plan,
// falling back to the linked Plan's values when the price_plan fields are blank.
// Used to populate the read-only Package section on the product-price drawer.
func lookupPackageNameDesc(ctx context.Context, deps *DetailViewDeps, pricePlanID string) (string, string) {
	if pricePlanID == "" || deps.ReadPricePlan == nil {
		return "", ""
	}
	resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{
		Data: &priceplanpb.PricePlan{Id: pricePlanID},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return "", ""
	}
	pp := resp.GetData()[0]
	name := strings.TrimSpace(pp.GetName())
	desc := strings.TrimSpace(pp.GetDescription())
	if name != "" && desc != "" {
		return name, desc
	}
	planName, planDesc := lookupPlanNameDesc(ctx, deps, pp.GetPlanId())
	if name == "" {
		name = planName
	}
	if desc == "" {
		desc = planDesc
	}
	return name, desc
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

// ---------------------------------------------------------------------------
// Option builder helpers — non-proto-enum only
// ---------------------------------------------------------------------------
//
// 2026-04-30 enum-select-canonicalize plan §6 — the proto-enum option
// builders (BillingKind, AmountBasis, BillingTreatment) are gone. Their
// option lists now live as hardcoded <option> tags in the drawer
// templates, and a checked-in drift test (price_plan/templates/templates_test.go)
// keeps them aligned with the proto enum's _name map.

// buildDurationUnitOptions builds select options for billing_cycle_unit / default_term_unit
// reusing the existing DurationUnit labels from CommonLabels. duration_unit is
// a plain string column (not a proto enum), so its option builder stays.
func buildDurationUnitOptions(cl pyeza.CommonLabels) []types.SelectOption {
	du := cl.DurationUnit
	return []types.SelectOption{
		{Value: "day", Label: du.DaySelect},
		{Value: "week", Label: du.WeekSelect},
		{Value: "month", Label: du.MonthSelect},
		{Value: "year", Label: du.YearSelect},
	}
}
