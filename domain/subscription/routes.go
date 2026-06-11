package subscription

import "strings"

// Subscription-domain route constants. Moved verbatim from the centymo root
// routes.go god-file (centymo W4). Pure structural relocation.
const (
	PlanListURL             = "/plans/list/{status}"
	PlanDetailURL           = "/plans/detail/{id}"
	PlanAddURL              = "/action/plan/add"
	PlanEditURL             = "/action/plan/edit/{id}"
	PlanDeleteURL           = "/action/plan/delete"
	PlanBulkDeleteURL       = "/action/plan/bulk-delete"
	PlanSetStatusURL        = "/action/plan/set-status"
	PlanBulkSetStatusURL    = "/action/plan/bulk-set-status"
	PlanTableURL            = "/action/plan/table/{status}"
	PlanTabActionURL        = "/action/plan/detail/{id}/tab/{tab}"
	PlanAttachmentUploadURL = "/action/plan/detail/{id}/attachments/upload"
	PlanAttachmentDeleteURL = "/action/plan/detail/{id}/attachments/delete"

	// PricePlan CRUD routes (within plan context)
	PricePlanAddURL    = "/action/plan/{id}/price-plans/add"
	PricePlanEditURL   = "/action/plan/{id}/price-plans/edit/{ppid}"
	PricePlanDeleteURL = "/action/plan/{id}/price-plans/delete"

	// Plan-scoped PricePlan detail — mirrors PriceSchedulePlanDetailURL but
	// keeps users in the Package (Plan) URL namespace so ActiveNav stays
	// anchored to the Services accordion's Packages section.
	// {id}=plan id, {ppid}=price_plan id.
	PlanPricePlanDetailURL             = "/plans/detail/{id}/price/{ppid}"
	PlanPricePlanTabActionURL          = "/action/plan/{id}/price/{ppid}/tab/{tab}"
	PlanPricePlanEditURL               = "/action/plan/{id}/price/{ppid}/edit"
	PlanPricePlanDeleteURL             = "/action/plan/{id}/price/{ppid}/delete"
	PlanPricePlanProductPriceAddURL    = "/action/plan/{id}/price/{ppid}/product-prices/add"
	PlanPricePlanProductPriceEditURL   = "/action/plan/{id}/price/{ppid}/product-prices/edit/{pppid}"
	PlanPricePlanProductPriceDeleteURL = "/action/plan/{id}/price/{ppid}/product-prices/delete"

	// ProductPlan CRUD routes (within plan context)
	PlanProductPlanAddURL    = "/action/plan/{id}/products/add"
	PlanProductPlanEditURL   = "/action/plan/{id}/products/edit/{ppid}"
	PlanProductPlanDeleteURL = "/action/plan/{id}/products/delete"
	PlanProductPlanPickerURL = "/action/plan/{id}/products/picker"

	// PricePlan standalone routes (rate cards as independent entity)
	PricePlanDashboardURL        = "/price-plans/dashboard"
	PricePlanListURL             = "/price-plans/list/{status}"
	PricePlanTableURL            = "/action/price-plan/table/{status}"
	PricePlanDetailURL           = "/price-plans/detail/{id}"
	PricePlanStandaloneAddURL    = "/action/price-plan/add"
	PricePlanStandaloneEditURL   = "/action/price-plan/edit/{id}"
	PricePlanStandaloneDeleteURL = "/action/price-plan/delete"
	PricePlanBulkDeleteURL       = "/action/price-plan/bulk-delete"
	PricePlanSetStatusURL        = "/action/price-plan/set-status"
	PricePlanBulkSetStatusURL    = "/action/price-plan/bulk-set-status"
	PricePlanTabActionURL        = "/action/price-plan/{id}/tab/{tab}"
	PricePlanAttachmentUploadURL = "/action/price-plan/{id}/attachments/upload"
	PricePlanAttachmentDeleteURL = "/action/price-plan/{id}/attachments/delete"

	// ProductPricePlan CRUD routes (within price plan / rate card detail)
	PricePlanProductPriceAddURL    = "/action/price-plan/{id}/product-prices/add"
	PricePlanProductPriceEditURL   = "/action/price-plan/{id}/product-prices/edit/{ppid}"
	PricePlanProductPriceDeleteURL = "/action/price-plan/{id}/product-prices/delete"

	// PriceSchedule routes (date-bounded pricing container for plans)
	PriceScheduleDashboardURL        = "/price-schedules/dashboard"
	PriceScheduleListURL             = "/price-schedules/list/{status}"
	PriceScheduleTableURL            = "/action/price-schedule/table/{status}"
	PriceScheduleDetailURL           = "/price-schedules/detail/{id}"
	PriceScheduleAddURL              = "/action/price-schedule/add"
	PriceScheduleEditURL             = "/action/price-schedule/edit/{id}"
	PriceScheduleDeleteURL           = "/action/price-schedule/delete"
	PriceScheduleBulkDeleteURL       = "/action/price-schedule/bulk-delete"
	PriceScheduleSetStatusURL        = "/action/price-schedule/set-status"
	PriceScheduleBulkSetStatusURL    = "/action/price-schedule/bulk-set-status"
	PriceScheduleTabActionURL        = "/action/price-schedule/{id}/tab/{tab}"
	PriceScheduleAttachmentUploadURL = "/action/price-schedule/{id}/attachments/upload"
	PriceScheduleAttachmentDeleteURL = "/action/price-schedule/{id}/attachments/delete"
	PriceSchedulePlanAddURL          = "/action/price-schedule/{id}/plan/add"
	// Schedule-scoped price_plan detail. Mirrors /app/price-plans/detail/{id} but nests
	// under the schedule so sidebar context stays on price-schedules (price_plan is no
	// longer a top-level sidebar entry).
	PriceSchedulePlanDetailURL             = "/price-schedules/detail/{id}/plan/{ppid}"
	PriceSchedulePlanTabActionURL          = "/action/price-schedule/{id}/plan/{ppid}/tab/{tab}"
	PriceSchedulePlanEditURL               = "/action/price-schedule/{id}/plan/{ppid}/edit"
	PriceSchedulePlanDeleteURL             = "/action/price-schedule/{id}/plan/{ppid}/delete"
	PriceSchedulePlanProductPriceAddURL    = "/action/price-schedule/{id}/plan/{ppid}/product-prices/add"
	PriceSchedulePlanProductPriceEditURL   = "/action/price-schedule/{id}/plan/{ppid}/product-prices/edit/{pppid}"
	PriceSchedulePlanProductPriceDeleteURL = "/action/price-schedule/{id}/plan/{ppid}/product-prices/delete"

	// Attachments on the nested price_schedule/plan detail page.
	PriceSchedulePlanAttachmentUploadURL = "/action/price-schedule/{id}/plan/{ppid}/attachments/upload"
	PriceSchedulePlanAttachmentDeleteURL = "/action/price-schedule/{id}/plan/{ppid}/attachments/delete"

	// 2026-05-04 — Subscriptions tab on the schedule-scoped
	// price_plan detail page. Same handler as SubscriptionDetailURL; the
	// nested URL alone activates the rate-card → plan → subscription breadcrumb.
	// See docs/plan/20260504-price-plan-engagements-tab/.
	PriceSchedulePlanSubscriptionDetailURL = "/price-schedules/detail/{id}/plan/{ppid}/subscription/{eid}"

	SubscriptionListURL = "/subscriptions/list/{status}"
	// SubscriptionTableURL returns ONLY the table-card partial — used as the
	// data-refresh-url so HTMX swaps the table without re-rendering the whole page.
	SubscriptionTableURL  = "/action/subscription/table/{status}"
	SubscriptionDetailURL = "/subscriptions/detail/{id}"
	// SubscriptionUnderClientDetailURL is the nested subscription-detail path
	// rendered with a client breadcrumb. Same view as SubscriptionDetailURL.
	SubscriptionUnderClientDetailURL  = "/clients/detail/{client_id}/subscriptions/{id}"
	SubscriptionAddURL                = "/action/subscription/add"
	SubscriptionEditURL               = "/action/subscription/edit/{id}"
	SubscriptionDeleteURL             = "/action/subscription/delete"
	SubscriptionBulkDeleteURL         = "/action/subscription/bulk-delete"
	SubscriptionSetStatusURL          = "/action/subscription/set-status"
	SubscriptionBulkSetStatusURL      = "/action/subscription/bulk-set-status"
	SubscriptionTabActionURL          = "/action/subscription/detail/{id}/tab/{tab}"
	SubscriptionAttachmentUploadURL   = "/action/subscription/detail/{id}/attachments/upload"
	SubscriptionAttachmentDeleteURL   = "/action/subscription/detail/{id}/attachments/delete"
	SubscriptionAttachmentDownloadURL = "/action/subscription/detail/{id}/attachments/download"
	SubscriptionSearchPlanURL         = "/action/subscription/search/plans"
	SubscriptionSearchClientURL       = "/action/subscription/search/clients"
	// SubscriptionRecognizeURL opens the "Recognize Revenue" drawer for a
	// subscription. GET = preview drawer (dry_run); POST = generate the Revenue.
	// Verb-first to avoid Go ServeMux ambiguity with /action/subscription/edit/{id}
	// — id-first and static-prefix patterns at the same depth can't disambiguate
	// (e.g. "/action/subscription/edit/recognize-revenue" matches both).
	SubscriptionRecognizeURL = "/action/subscription/recognize-revenue/{id}"

	// SubscriptionRevenueRunURL opens the "Invoice Run" drawer for a single
	// subscription (Surface C — per-subscription drawer, CYCLE billing_kind only).
	// Verb-first ("revenue-run") to avoid ServeMux ambiguity consistent with
	// SubscriptionRecognizeURL above.
	SubscriptionRevenueRunURL = "/action/subscription/revenue-run/{id}"

	// SubscriptionCustomizePackageURL is the POST endpoint that drives the
	// "Customize this package for {ClientName}" CTA on the subscription
	// detail's Package tab. Calls espyna's CustomizePlanForClient use case
	// and HX-redirects to the new (cloned) PricePlan's package page.
	// Verb-first ("customize-package") to avoid the same ServeMux ambiguity
	// SubscriptionRecognizeURL above guards against.
	SubscriptionCustomizePackageURL = "/action/subscription/customize-package/{id}"

	// 2026-04-29 milestone-billing plan §5 / Phase D — mark-ready + waive
	// handlers for BillingEvent rows on the subscription Package tab.
	// Both POST through the espyna BillingEvent.SetStatus domain service.
	MilestoneMarkReadyURL = "/action/subscription/{id}/billing-event/{eventId}/mark-ready"
	MilestoneWaiveURL     = "/action/subscription/{id}/billing-event/{eventId}/waive"

	// 20260517-advance-cash-events Plan B Phase 7 — Recognize handler for a
	// BillingEvent row when it is linked to a MILESTONE advance Collection via
	// the collection_billing_event junction. POSTs through the
	// espyna RecognizeMilestoneAdvanceCollection use case.
	MilestoneRecognizeURL = "/action/subscription/{id}/billing-event/{eventId}/recognize"

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — retroactive
	// spawn drawer endpoint (GET = drawer, POST = spawn) and HTMX-driven
	// partial that re-renders the Spawn Jobs section in the create drawer
	// when the operator changes the selected Plan / PricePlan.
	//
	// Verb-first ("spawn-jobs") to avoid the Go ServeMux ambiguity that would
	// otherwise pit `{subscriptionId}/spawn-jobs` against
	// `table/{status}` (and similar id-first/static-prefix patterns at the
	// same depth — same root cause as SubscriptionRecognizeURL above).
	SubscriptionSpawnJobsURL        = "/action/subscription/spawn-jobs/{subscriptionId}"
	SubscriptionSpawnJobsPartialURL = "/action/subscription/_partial/spawn-jobs-section"

	// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D — manual cycle
	// spawn + backfill triggers. Both routes call into espyna's
	// MaterializeInstanceJobsForSubscription consumer (single-cycle vs.
	// multi-cycle modes). Verb-first ("spawn-cycle-jobs", "backfill-cycle-jobs")
	// to keep ServeMux disambiguation consistent with existing
	// SubscriptionRecognizeURL / SubscriptionSpawnJobsURL.
	SubscriptionSpawnCycleJobsURL    = "/action/subscription/spawn-cycle-jobs/{subscriptionId}"
	SubscriptionBackfillCycleJobsURL = "/action/subscription/backfill-cycle-jobs/{subscriptionId}"

	// 2026-05-01 ad-hoc-subscription-billing plan §5.2 — operator-driven CTA
	// for AD_HOC subscriptions. Pool-Generate-Invoice reuses the existing
	// SubscriptionRecognizeURL; Extend-Pool deferred to v1.5.5 (needs new
	// espyna use case for Subscription.entitled_occurrences_override write).
	SubscriptionRequestUsageURL = "/action/subscription/request-usage/{subscriptionId}"
)

type PlanRoutes struct {
	// Sidebar navigation context — set via defaults or routes.json override
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// PricePlan CRUD routes (within plan context)
	PricePlanAddURL    string `json:"price_plan_add_url"`
	PricePlanEditURL   string `json:"price_plan_edit_url"`
	PricePlanDeleteURL string `json:"price_plan_delete_url"`

	// Plan-scoped PricePlan detail (mirrors PriceSchedulePlanRoutes.PlanDetailURL
	// but anchored under /app/plans/detail/{id}/price/{ppid}). Lets the package
	// detail's package-prices tab keep ActiveNav on Services > Packages instead
	// of jumping to the rate-cards namespace.
	PricePlanDetailURL             string `json:"price_plan_detail_url"`
	PricePlanTabActionURL          string `json:"price_plan_tab_action_url"`
	PricePlanProductPriceAddURL    string `json:"price_plan_product_price_add_url"`
	PricePlanProductPriceEditURL   string `json:"price_plan_product_price_edit_url"`
	PricePlanProductPriceDeleteURL string `json:"price_plan_product_price_delete_url"`

	// ProductPlan CRUD routes (within plan context)
	ProductPlanAddURL    string `json:"product_plan_add_url"`
	ProductPlanEditURL   string `json:"product_plan_edit_url"`
	ProductPlanDeleteURL string `json:"product_plan_delete_url"`
	ProductPlanPickerURL string `json:"product_plan_picker_url"`
}

// DefaultPlanRoutes returns a PlanRoutes populated from the package-level
// route constants defined in routes.go.
func DefaultPlanRoutes() PlanRoutes {
	return PlanRoutes{
		ActiveNav:    "service",
		ActiveSubNav: "plans",

		ListURL:          PlanListURL,
		TableURL:         PlanTableURL,
		DetailURL:        PlanDetailURL,
		AddURL:           PlanAddURL,
		EditURL:          PlanEditURL,
		DeleteURL:        PlanDeleteURL,
		BulkDeleteURL:    PlanBulkDeleteURL,
		SetStatusURL:     PlanSetStatusURL,
		BulkSetStatusURL: PlanBulkSetStatusURL,
		TabActionURL:     PlanTabActionURL,

		AttachmentUploadURL: PlanAttachmentUploadURL,
		AttachmentDeleteURL: PlanAttachmentDeleteURL,

		PricePlanAddURL:    PricePlanAddURL,
		PricePlanEditURL:   PricePlanEditURL,
		PricePlanDeleteURL: PricePlanDeleteURL,

		PricePlanDetailURL:             PlanPricePlanDetailURL,
		PricePlanTabActionURL:          PlanPricePlanTabActionURL,
		PricePlanProductPriceAddURL:    PlanPricePlanProductPriceAddURL,
		PricePlanProductPriceEditURL:   PlanPricePlanProductPriceEditURL,
		PricePlanProductPriceDeleteURL: PlanPricePlanProductPriceDeleteURL,

		ProductPlanAddURL:    PlanProductPlanAddURL,
		ProductPlanEditURL:   PlanProductPlanEditURL,
		ProductPlanDeleteURL: PlanProductPlanDeleteURL,
		ProductPlanPickerURL: PlanProductPlanPickerURL,
	}
}

// DefaultPlanBundleRoutes returns a PlanRoutes with every URL namespace-shifted
// from the services namespace onto the inventory accordion namespace. Used as
// the route base for the Plan inventory-mount registration in block.go; a lyngua
// `plan_bundle` override can layer additional tweaks on top.
//
// Bundle-mount variant of Plan routes — shifts every page + action URL from the
// services namespace (/app/plans/*) onto the inventory accordion namespace
// (/app/inventory/bundles/*). Used as the route base for the Plan
// inventory-mount registration in block.go; a lyngua `plan_bundle` override can
// layer additional tweaks on top.
//
// Shift rules:
//   - "/app/plans/"   → "/app/inventory/bundles/"
//   - "/action/plan/" → "/action/inventory-bundle/"
func DefaultPlanBundleRoutes() PlanRoutes {
	r := DefaultPlanRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "bundles-active"
	// shift matches both pre-P4 (`/app/plans/*`) and post-P4 (`/plans/*`)
	// constant shapes — see DefaultProductInventoryRoutes shift comment
	// for the P4 regression context.
	shift := func(s string) string {
		s = strings.Replace(s, "/app/plans/", "/app/inventory/bundles/", 1)
		s = strings.Replace(s, "/plans/", "/inventory/bundles/", 1)
		s = strings.Replace(s, "/action/plan/", "/action/inventory-bundle/", 1)
		return s
	}
	r.ListURL = shift(r.ListURL)
	r.TableURL = shift(r.TableURL)
	r.DetailURL = shift(r.DetailURL)
	r.AddURL = shift(r.AddURL)
	r.EditURL = shift(r.EditURL)
	r.DeleteURL = shift(r.DeleteURL)
	r.BulkDeleteURL = shift(r.BulkDeleteURL)
	r.SetStatusURL = shift(r.SetStatusURL)
	r.BulkSetStatusURL = shift(r.BulkSetStatusURL)
	r.TabActionURL = shift(r.TabActionURL)
	r.AttachmentUploadURL = shift(r.AttachmentUploadURL)
	r.AttachmentDeleteURL = shift(r.AttachmentDeleteURL)
	r.PricePlanAddURL = shift(r.PricePlanAddURL)
	r.PricePlanEditURL = shift(r.PricePlanEditURL)
	r.PricePlanDeleteURL = shift(r.PricePlanDeleteURL)
	r.PricePlanDetailURL = shift(r.PricePlanDetailURL)
	r.PricePlanTabActionURL = shift(r.PricePlanTabActionURL)
	r.PricePlanProductPriceAddURL = shift(r.PricePlanProductPriceAddURL)
	r.PricePlanProductPriceEditURL = shift(r.PricePlanProductPriceEditURL)
	r.PricePlanProductPriceDeleteURL = shift(r.PricePlanProductPriceDeleteURL)
	r.ProductPlanAddURL = shift(r.ProductPlanAddURL)
	r.ProductPlanEditURL = shift(r.ProductPlanEditURL)
	r.ProductPlanDeleteURL = shift(r.ProductPlanDeleteURL)
	r.ProductPlanPickerURL = shift(r.ProductPlanPickerURL)
	return r
}

// RouteMap returns a map of dot-notation keys to route paths for all
// plan routes.
func (r PlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"plan.list":            r.ListURL,
		"plan.table":           r.TableURL,
		"plan.detail":          r.DetailURL,
		"plan.add":             r.AddURL,
		"plan.edit":            r.EditURL,
		"plan.delete":          r.DeleteURL,
		"plan.bulk_delete":     r.BulkDeleteURL,
		"plan.set_status":      r.SetStatusURL,
		"plan.bulk_set_status": r.BulkSetStatusURL,
		"plan.tab_action":      r.TabActionURL,

		"plan.attachment.upload": r.AttachmentUploadURL,
		"plan.attachment.delete": r.AttachmentDeleteURL,

		"plan.pricelist.add":    r.PricePlanAddURL,
		"plan.pricelist.edit":   r.PricePlanEditURL,
		"plan.pricelist.delete": r.PricePlanDeleteURL,

		"plan.price.detail":               r.PricePlanDetailURL,
		"plan.price.tab_action":           r.PricePlanTabActionURL,
		"plan.price.product_price.add":    r.PricePlanProductPriceAddURL,
		"plan.price.product_price.edit":   r.PricePlanProductPriceEditURL,
		"plan.price.product_price.delete": r.PricePlanProductPriceDeleteURL,

		"plan.product_plan.add":    r.ProductPlanAddURL,
		"plan.product_plan.edit":   r.ProductPlanEditURL,
		"plan.product_plan.delete": r.ProductPlanDeleteURL,
		"plan.product_plan.picker": r.ProductPlanPickerURL,
	}
}

// PricePlanRoutes holds all route paths for price plan (rate card) views and actions.
type PricePlanRoutes struct {
	ActiveNav           string `json:"active_nav"`
	ActiveSubNav        string `json:"active_sub_nav"`
	DashboardURL        string `json:"dashboard_url"`
	ListURL             string `json:"list_url"`
	TableURL            string `json:"table_url"`
	DetailURL           string `json:"detail_url"`
	AddURL              string `json:"add_url"`
	EditURL             string `json:"edit_url"`
	DeleteURL           string `json:"delete_url"`
	BulkDeleteURL       string `json:"bulk_delete_url"`
	SetStatusURL        string `json:"set_status_url"`
	BulkSetStatusURL    string `json:"bulk_set_status_url"`
	TabActionURL        string `json:"tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// ProductPricePlan CRUD routes (within rate card detail)
	ProductPriceAddURL    string `json:"product_price_add_url"`
	ProductPriceEditURL   string `json:"product_price_edit_url"`
	ProductPriceDeleteURL string `json:"product_price_delete_url"`
}

// DefaultPricePlanRoutes returns a PricePlanRoutes populated from the package-level
// route constants defined in routes.go.
func DefaultPricePlanRoutes() PricePlanRoutes {
	return PricePlanRoutes{
		ActiveNav:             "service",
		ActiveSubNav:          "rate-cards",
		DashboardURL:          PricePlanDashboardURL,
		ListURL:               PricePlanListURL,
		TableURL:              PricePlanTableURL,
		DetailURL:             PricePlanDetailURL,
		AddURL:                PricePlanStandaloneAddURL,
		EditURL:               PricePlanStandaloneEditURL,
		DeleteURL:             PricePlanStandaloneDeleteURL,
		BulkDeleteURL:         PricePlanBulkDeleteURL,
		SetStatusURL:          PricePlanSetStatusURL,
		BulkSetStatusURL:      PricePlanBulkSetStatusURL,
		TabActionURL:          PricePlanTabActionURL,
		AttachmentUploadURL:   PricePlanAttachmentUploadURL,
		AttachmentDeleteURL:   PricePlanAttachmentDeleteURL,
		ProductPriceAddURL:    PricePlanProductPriceAddURL,
		ProductPriceEditURL:   PricePlanProductPriceEditURL,
		ProductPriceDeleteURL: PricePlanProductPriceDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// price plan routes.
func (r PricePlanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"price_plan.dashboard":            r.DashboardURL,
		"price_plan.list":                 r.ListURL,
		"price_plan.table":                r.TableURL,
		"price_plan.detail":               r.DetailURL,
		"price_plan.add":                  r.AddURL,
		"price_plan.edit":                 r.EditURL,
		"price_plan.delete":               r.DeleteURL,
		"price_plan.bulk_delete":          r.BulkDeleteURL,
		"price_plan.set_status":           r.SetStatusURL,
		"price_plan.bulk_set_status":      r.BulkSetStatusURL,
		"price_plan.tab_action":           r.TabActionURL,
		"price_plan.attachment.upload":    r.AttachmentUploadURL,
		"price_plan.attachment.delete":    r.AttachmentDeleteURL,
		"price_plan.product_price.add":    r.ProductPriceAddURL,
		"price_plan.product_price.edit":   r.ProductPriceEditURL,
		"price_plan.product_price.delete": r.ProductPriceDeleteURL,
	}
}

// PriceScheduleRoutes holds all route paths for price schedule views and actions.
type PriceScheduleRoutes struct {
	ActiveNav                 string `json:"active_nav"`
	ActiveSubNav              string `json:"active_sub_nav"`
	DashboardURL              string `json:"dashboard_url"`
	ListURL                   string `json:"list_url"`
	TableURL                  string `json:"table_url"`
	DetailURL                 string `json:"detail_url"`
	AddURL                    string `json:"add_url"`
	EditURL                   string `json:"edit_url"`
	DeleteURL                 string `json:"delete_url"`
	BulkDeleteURL             string `json:"bulk_delete_url"`
	SetStatusURL              string `json:"set_status_url"`
	BulkSetStatusURL          string `json:"bulk_set_status_url"`
	TabActionURL              string `json:"tab_action_url"`
	AttachmentUploadURL       string `json:"attachment_upload_url"`
	AttachmentDeleteURL       string `json:"attachment_delete_url"`
	PlanAddURL                string `json:"plan_add_url"`
	PlanDetailURL             string `json:"plan_detail_url"`
	PlanTabActionURL          string `json:"plan_tab_action_url"`
	PlanEditURL               string `json:"plan_edit_url"`
	PlanDeleteURL             string `json:"plan_delete_url"`
	PlanProductPriceAddURL    string `json:"plan_product_price_add_url"`
	PlanProductPriceEditURL   string `json:"plan_product_price_edit_url"`
	PlanProductPriceDeleteURL string `json:"plan_product_price_delete_url"`
	PlanAttachmentUploadURL   string `json:"plan_attachment_upload_url"`
	PlanAttachmentDeleteURL   string `json:"plan_attachment_delete_url"`
	// 2026-05-04 — Subscription detail nested under the
	// schedule-scoped price_plan path. Activates the rate-card → plan →
	// subscription breadcrumb in the subscription detail view. Empty string
	// disables the nested route. See
	// docs/plan/20260504-price-plan-engagements-tab/.
	PlanSubscriptionDetailURL string `json:"plan_subscription_detail_url"`
}

// DefaultPriceScheduleRoutes returns a PriceScheduleRoutes populated from the package-level
// route constants defined in routes.go.
func DefaultPriceScheduleRoutes() PriceScheduleRoutes {
	return PriceScheduleRoutes{
		ActiveNav:                 "service",
		ActiveSubNav:              "price-schedules",
		DashboardURL:              PriceScheduleDashboardURL,
		ListURL:                   PriceScheduleListURL,
		TableURL:                  PriceScheduleTableURL,
		DetailURL:                 PriceScheduleDetailURL,
		AddURL:                    PriceScheduleAddURL,
		EditURL:                   PriceScheduleEditURL,
		DeleteURL:                 PriceScheduleDeleteURL,
		BulkDeleteURL:             PriceScheduleBulkDeleteURL,
		SetStatusURL:              PriceScheduleSetStatusURL,
		BulkSetStatusURL:          PriceScheduleBulkSetStatusURL,
		TabActionURL:              PriceScheduleTabActionURL,
		AttachmentUploadURL:       PriceScheduleAttachmentUploadURL,
		AttachmentDeleteURL:       PriceScheduleAttachmentDeleteURL,
		PlanAddURL:                PriceSchedulePlanAddURL,
		PlanDetailURL:             PriceSchedulePlanDetailURL,
		PlanTabActionURL:          PriceSchedulePlanTabActionURL,
		PlanEditURL:               PriceSchedulePlanEditURL,
		PlanDeleteURL:             PriceSchedulePlanDeleteURL,
		PlanProductPriceAddURL:    PriceSchedulePlanProductPriceAddURL,
		PlanProductPriceEditURL:   PriceSchedulePlanProductPriceEditURL,
		PlanProductPriceDeleteURL: PriceSchedulePlanProductPriceDeleteURL,
		PlanAttachmentUploadURL:   PriceSchedulePlanAttachmentUploadURL,
		PlanAttachmentDeleteURL:   PriceSchedulePlanAttachmentDeleteURL,
		PlanSubscriptionDetailURL: PriceSchedulePlanSubscriptionDetailURL,
	}
}

// DefaultPriceScheduleInventoryRoutes returns a PriceScheduleRoutes with every
// URL namespace-shifted from the services namespace onto the inventory accordion
// namespace. Used as the route base for the PriceSchedule inventory-mount
// registration in block.go; a lyngua `price_schedule_inventory` override can
// layer additional tweaks on top.
//
// Shift rules:
//   - "/app/price-schedules/"   → "/app/inventory/price-schedules/"
//   - "/action/price-schedule/" → "/action/inventory-price-schedule/"
func DefaultPriceScheduleInventoryRoutes() PriceScheduleRoutes {
	r := DefaultPriceScheduleRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "inventory-price-schedules-active"
	// shift matches both pre-P4 (`/app/price-schedules/*`) and post-P4
	// (`/price-schedules/*`) constant shapes — see DefaultProductInventoryRoutes
	// shift comment for the P4 regression context.
	shift := func(s string) string {
		s = strings.Replace(s, "/app/price-schedules/", "/app/inventory/price-schedules/", 1)
		s = strings.Replace(s, "/price-schedules/", "/inventory/price-schedules/", 1)
		s = strings.Replace(s, "/action/price-schedule/", "/action/inventory-price-schedule/", 1)
		return s
	}
	r.DashboardURL = shift(r.DashboardURL)
	r.ListURL = shift(r.ListURL)
	r.TableURL = shift(r.TableURL)
	r.DetailURL = shift(r.DetailURL)
	r.AddURL = shift(r.AddURL)
	r.EditURL = shift(r.EditURL)
	r.DeleteURL = shift(r.DeleteURL)
	r.BulkDeleteURL = shift(r.BulkDeleteURL)
	r.SetStatusURL = shift(r.SetStatusURL)
	r.BulkSetStatusURL = shift(r.BulkSetStatusURL)
	r.TabActionURL = shift(r.TabActionURL)
	r.AttachmentUploadURL = shift(r.AttachmentUploadURL)
	r.AttachmentDeleteURL = shift(r.AttachmentDeleteURL)
	r.PlanAddURL = shift(r.PlanAddURL)
	r.PlanDetailURL = shift(r.PlanDetailURL)
	r.PlanTabActionURL = shift(r.PlanTabActionURL)
	r.PlanEditURL = shift(r.PlanEditURL)
	r.PlanDeleteURL = shift(r.PlanDeleteURL)
	r.PlanProductPriceAddURL = shift(r.PlanProductPriceAddURL)
	r.PlanProductPriceEditURL = shift(r.PlanProductPriceEditURL)
	r.PlanProductPriceDeleteURL = shift(r.PlanProductPriceDeleteURL)
	r.PlanAttachmentUploadURL = shift(r.PlanAttachmentUploadURL)
	r.PlanAttachmentDeleteURL = shift(r.PlanAttachmentDeleteURL)
	r.PlanSubscriptionDetailURL = shift(r.PlanSubscriptionDetailURL)
	return r
}

// RouteMap returns a map of dot-notation keys to route paths for all
// price schedule routes.
func (r PriceScheduleRoutes) RouteMap() map[string]string {
	return map[string]string{
		"price_schedule.dashboard":                 r.DashboardURL,
		"price_schedule.list":                      r.ListURL,
		"price_schedule.table":                     r.TableURL,
		"price_schedule.detail":                    r.DetailURL,
		"price_schedule.add":                       r.AddURL,
		"price_schedule.edit":                      r.EditURL,
		"price_schedule.delete":                    r.DeleteURL,
		"price_schedule.bulk_delete":               r.BulkDeleteURL,
		"price_schedule.set_status":                r.SetStatusURL,
		"price_schedule.bulk_set_status":           r.BulkSetStatusURL,
		"price_schedule.tab_action":                r.TabActionURL,
		"price_schedule.attachment.upload":         r.AttachmentUploadURL,
		"price_schedule.attachment.delete":         r.AttachmentDeleteURL,
		"price_schedule.plan.add":                  r.PlanAddURL,
		"price_schedule.plan.detail":               r.PlanDetailURL,
		"price_schedule.plan.tab_action":           r.PlanTabActionURL,
		"price_schedule.plan.edit":                 r.PlanEditURL,
		"price_schedule.plan.delete":               r.PlanDeleteURL,
		"price_schedule.plan.product_price.add":    r.PlanProductPriceAddURL,
		"price_schedule.plan.product_price.edit":   r.PlanProductPriceEditURL,
		"price_schedule.plan.product_price.delete": r.PlanProductPriceDeleteURL,
		"price_schedule.plan.attachment.upload":    r.PlanAttachmentUploadURL,
		"price_schedule.plan.attachment.delete":    r.PlanAttachmentDeleteURL,
		"price_schedule.plan.subscription.detail":  r.PlanSubscriptionDetailURL,
	}
}

// SubscriptionRoutes holds all route paths for subscription views and actions.
type SubscriptionRoutes struct {
	// Sidebar navigation context — set via defaults or routes.json override
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL              string `json:"list_url"`
	TableURL             string `json:"table_url"`
	DetailURL            string `json:"detail_url"`
	UnderClientDetailURL string `json:"under_client_detail_url"`
	AddURL               string `json:"add_url"`
	EditURL              string `json:"edit_url"`
	DeleteURL            string `json:"delete_url"`
	BulkDeleteURL        string `json:"bulk_delete_url"`
	SetStatusURL         string `json:"set_status_url"`
	BulkSetStatusURL     string `json:"bulk_set_status_url"`
	TabActionURL         string `json:"tab_action_url"`
	SearchPlanURL        string `json:"search_plan_url"`
	SearchClientURL      string `json:"search_client_url"`

	// RecognizeURL opens the "Recognize Revenue" drawer for a subscription.
	// GET = preview drawer (dry_run); POST = generate the Revenue.
	RecognizeURL string `json:"recognize_url"`

	// CustomizePackageURL is the POST endpoint that drives the "Customize
	// this package for {ClientName}" CTA on the subscription detail's
	// Package tab. Server clones the source Plan + PricePlan into a
	// client-scoped copy and HX-redirects to the new package URL.
	CustomizePackageURL string `json:"customize_package_url"`

	// 2026-04-29 milestone-billing plan §5 / Phase D — mark-ready + waive
	// endpoints for BillingEvent rows on the subscription Package tab.
	MilestoneMarkReadyURL string `json:"milestone_mark_ready_url"`
	MilestoneWaiveURL     string `json:"milestone_waive_url"`

	// 20260517-advance-cash-events Plan B Phase 7 — Recognize button on a
	// BillingEvent row when it is linked to a MILESTONE advance Collection.
	// POSTs through espyna RecognizeMilestoneAdvanceCollection.
	MilestoneRecognizeURL string `json:"milestone_recognize_url"`

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / Phase D —
	// retroactive spawn drawer URL + HTMX-driven partial URL for the
	// Spawn Jobs section on the create form.
	SpawnJobsURL        string `json:"spawn_jobs_url"`
	SpawnJobsPartialURL string `json:"spawn_jobs_partial_url"`

	// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D — manual cycle
	// spawn + backfill triggers. Both POST through espyna's
	// MaterializeInstanceJobsForSubscription consumer. Backfill GET renders
	// a preview drawer; POST commits the spawn.
	SpawnCycleJobsURL    string `json:"spawn_cycle_jobs_url"`
	BackfillCycleJobsURL string `json:"backfill_cycle_jobs_url"`

	// 2026-05-01 ad-hoc-subscription-billing — operator-driven Request Usage CTA.
	RequestUsageURL string `json:"request_usage_url"`

	// 2026-05-06 revenue-run — per-subscription Invoice Run drawer (Surface C,
	// CYCLE billing_kind only). Empty string when revenue-run module is not wired.
	RevenueRunURL string `json:"revenue_run_url"`

	// Attachment routes
	AttachmentUploadURL   string `json:"attachment_upload_url"`
	AttachmentDeleteURL   string `json:"attachment_delete_url"`
	AttachmentDownloadURL string `json:"attachment_download_url"`
}

// DefaultSubscriptionRoutes returns a SubscriptionRoutes populated from the
// package-level route constants defined in routes.go.
func DefaultSubscriptionRoutes() SubscriptionRoutes {
	return SubscriptionRoutes{
		ActiveNav:    "client",
		ActiveSubNav: "subscriptions",

		ListURL:              SubscriptionListURL,
		TableURL:             SubscriptionTableURL,
		DetailURL:            SubscriptionDetailURL,
		UnderClientDetailURL: SubscriptionUnderClientDetailURL,
		AddURL:               SubscriptionAddURL,
		EditURL:              SubscriptionEditURL,
		DeleteURL:            SubscriptionDeleteURL,
		BulkDeleteURL:        SubscriptionBulkDeleteURL,
		SetStatusURL:         SubscriptionSetStatusURL,
		BulkSetStatusURL:     SubscriptionBulkSetStatusURL,
		TabActionURL:         SubscriptionTabActionURL,
		SearchPlanURL:        SubscriptionSearchPlanURL,
		SearchClientURL:      SubscriptionSearchClientURL,
		RecognizeURL:         SubscriptionRecognizeURL,
		CustomizePackageURL:  SubscriptionCustomizePackageURL,

		// 2026-04-29 milestone-billing.
		MilestoneMarkReadyURL: MilestoneMarkReadyURL,
		MilestoneWaiveURL:     MilestoneWaiveURL,

		// 20260517-advance-cash-events Plan B Phase 7.
		MilestoneRecognizeURL: MilestoneRecognizeURL,

		// 2026-04-29 auto-spawn-jobs-from-subscription.
		SpawnJobsURL:        SubscriptionSpawnJobsURL,
		SpawnJobsPartialURL: SubscriptionSpawnJobsPartialURL,

		// 2026-04-30 cyclic-subscription-jobs.
		SpawnCycleJobsURL:    SubscriptionSpawnCycleJobsURL,
		BackfillCycleJobsURL: SubscriptionBackfillCycleJobsURL,

		// 2026-05-01 ad-hoc-subscription-billing.
		RequestUsageURL: SubscriptionRequestUsageURL,

		// 2026-05-06 revenue-run — per-subscription drawer.
		RevenueRunURL: SubscriptionRevenueRunURL,

		AttachmentUploadURL:   SubscriptionAttachmentUploadURL,
		AttachmentDeleteURL:   SubscriptionAttachmentDeleteURL,
		AttachmentDownloadURL: SubscriptionAttachmentDownloadURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// subscription routes.
func (r SubscriptionRoutes) RouteMap() map[string]string {
	return map[string]string{
		"subscription.list":                r.ListURL,
		"subscription.table":               r.TableURL,
		"subscription.detail":              r.DetailURL,
		"subscription.under_client_detail": r.UnderClientDetailURL,
		"subscription.add":                 r.AddURL,
		"subscription.edit":                r.EditURL,
		"subscription.delete":              r.DeleteURL,
		"subscription.bulk_delete":         r.BulkDeleteURL,
		"subscription.set_status":          r.SetStatusURL,
		"subscription.bulk_set_status":     r.BulkSetStatusURL,
		"subscription.tab_action":          r.TabActionURL,
		"subscription.search_plan":         r.SearchPlanURL,
		"subscription.search_client":       r.SearchClientURL,
		"subscription.recognize":           r.RecognizeURL,
		"subscription.customize_package":   r.CustomizePackageURL,

		// 2026-04-29 milestone-billing routes.
		"milestone.mark_ready": r.MilestoneMarkReadyURL,
		"milestone.waive":      r.MilestoneWaiveURL,

		// 20260517-advance-cash-events Plan B Phase 7.
		"milestone.recognize": r.MilestoneRecognizeURL,

		// 2026-04-29 auto-spawn-jobs-from-subscription routes.
		"subscription.spawn_jobs":         r.SpawnJobsURL,
		"subscription.spawn_jobs_partial": r.SpawnJobsPartialURL,

		// 2026-04-30 cyclic-subscription-jobs routes.
		"subscription.spawn_cycle_jobs":    r.SpawnCycleJobsURL,
		"subscription.backfill_cycle_jobs": r.BackfillCycleJobsURL,

		// 2026-05-01 ad-hoc-subscription-billing routes.
		"subscription.request_usage": r.RequestUsageURL,

		// 2026-05-06 revenue-run per-subscription drawer.
		"subscription.revenue_run": r.RevenueRunURL,

		"subscription.attachment.upload":   r.AttachmentUploadURL,
		"subscription.attachment.delete":   r.AttachmentDeleteURL,
		"subscription.attachment.download": r.AttachmentDownloadURL,
	}
}

// CollectionRoutes holds all route paths for collection (money IN) views
// and actions.
