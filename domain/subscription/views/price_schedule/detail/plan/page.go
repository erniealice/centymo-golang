// Package plan renders the price_plan detail page nested under its parent PriceSchedule
// at /app/price-schedules/detail/{id}/plan/{ppid}. The sidebar stays on price-schedules
// because price_plan is no longer a top-level sidebar entry.
package plan

import (
	"context"
	"fmt"
	"log"
	"strings"

	subscription "github.com/erniealice/centymo-golang/domain/subscription"
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
	Routes                 subscription.PriceScheduleRoutes
	ScheduleLabels         subscription.PriceScheduleLabels
	PlanLabels             subscription.PricePlanLabels
	ProductPricePlanLabels subscription.ProductPricePlanLabels
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
	// subscription URL when PlanSubscriptionDetailURL is non-empty.
	SubscriptionEditURL   string
	SubscriptionDeleteURL string
	// SubscriptionAddURL is the route for the "Add Subscription" primary
	// action on the subscriptions tab. The drawer is opened in
	// price-plan-locked mode via ?price_plan_id=&plan_label=&client_id=...
	// query params, mirroring how the client-detail subscriptions tab opens
	// the same drawer in client-locked mode.
	SubscriptionAddURL string
	// PlanSubscriptionDetailURL is the schedule-scoped subscription URL template
	// /app/price-schedules/detail/{id}/plan/{ppid}/subscription/{eid}. When
	// set, the row's "View" action targets this nested URL so the
	// subscription detail page renders with a rate-card → plan → subscription
	// breadcrumb. Empty falls back to SubscriptionDetailURL.
	PlanSubscriptionDetailURL string
	SubscriptionDetailURL     string

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
	Labels          subscription.PricePlanLabels
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
		pageData.SubscriptionsTable = buildSubscriptionsTable(ctx, deps, sid, ppid, pp, effectiveName, rows)
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
