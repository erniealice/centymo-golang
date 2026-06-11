package subscription

// Subscription-domain route constants. Relocated from the subscription
// routes.go god-file (entity-local extraction). Pure structural move.
const (
	ListURL = "/subscriptions/list/{status}"
	// TableURL returns ONLY the table-card partial — used as the
	// data-refresh-url so HTMX swaps the table without re-rendering the whole page.
	TableURL  = "/action/subscription/table/{status}"
	DetailURL = "/subscriptions/detail/{id}"
	// UnderClientDetailURL is the nested subscription-detail path
	// rendered with a client breadcrumb. Same view as DetailURL.
	UnderClientDetailURL  = "/clients/detail/{client_id}/subscriptions/{id}"
	AddURL                = "/action/subscription/add"
	EditURL               = "/action/subscription/edit/{id}"
	DeleteURL             = "/action/subscription/delete"
	BulkDeleteURL         = "/action/subscription/bulk-delete"
	SetStatusURL          = "/action/subscription/set-status"
	BulkSetStatusURL      = "/action/subscription/bulk-set-status"
	TabActionURL          = "/action/subscription/detail/{id}/tab/{tab}"
	AttachmentUploadURL   = "/action/subscription/detail/{id}/attachments/upload"
	AttachmentDeleteURL   = "/action/subscription/detail/{id}/attachments/delete"
	AttachmentDownloadURL = "/action/subscription/detail/{id}/attachments/download"
	SearchPlanURL         = "/action/subscription/search/plans"
	SearchClientURL       = "/action/subscription/search/clients"
	// RecognizeURL opens the "Recognize Revenue" drawer for a
	// subscription. GET = preview drawer (dry_run); POST = generate the Revenue.
	// Verb-first to avoid Go ServeMux ambiguity with /action/subscription/edit/{id}
	// — id-first and static-prefix patterns at the same depth can't disambiguate
	// (e.g. "/action/subscription/edit/recognize-revenue" matches both).
	RecognizeURL = "/action/subscription/recognize-revenue/{id}"

	// RevenueRunURL opens the "Invoice Run" drawer for a single
	// subscription (Surface C — per-subscription drawer, CYCLE billing_kind only).
	// Verb-first ("revenue-run") to avoid ServeMux ambiguity consistent with
	// RecognizeURL above.
	RevenueRunURL = "/action/subscription/revenue-run/{id}"

	// CustomizePackageURL is the POST endpoint that drives the
	// "Customize this package for {ClientName}" CTA on the subscription
	// detail's Package tab. Calls espyna's CustomizePlanForClient use case
	// and HX-redirects to the new (cloned) PricePlan's package page.
	// Verb-first ("customize-package") to avoid the same ServeMux ambiguity
	// RecognizeURL above guards against.
	CustomizePackageURL = "/action/subscription/customize-package/{id}"

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
	// same depth — same root cause as RecognizeURL above).
	SpawnJobsURL        = "/action/subscription/spawn-jobs/{subscriptionId}"
	SpawnJobsPartialURL = "/action/subscription/_partial/spawn-jobs-section"

	// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D — manual cycle
	// spawn + backfill triggers. Both routes call into espyna's
	// MaterializeInstanceJobsForSubscription consumer (single-cycle vs.
	// multi-cycle modes). Verb-first ("spawn-cycle-jobs", "backfill-cycle-jobs")
	// to keep ServeMux disambiguation consistent with existing
	// RecognizeURL / SpawnJobsURL.
	SpawnCycleJobsURL    = "/action/subscription/spawn-cycle-jobs/{subscriptionId}"
	BackfillCycleJobsURL = "/action/subscription/backfill-cycle-jobs/{subscriptionId}"

	// 2026-05-01 ad-hoc-subscription-billing plan §5.2 — operator-driven CTA
	// for AD_HOC subscriptions. Pool-Generate-Invoice reuses the existing
	// RecognizeURL; Extend-Pool deferred to v1.5.5 (needs new
	// espyna use case for Subscription.entitled_occurrences_override write).
	RequestUsageURL = "/action/subscription/request-usage/{subscriptionId}"
)

// Routes holds all route paths for subscription views and actions.
type Routes struct {
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

// DefaultRoutes returns a Routes populated from the
// package-level route constants defined in routes.go.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "client",
		ActiveSubNav: "subscriptions",

		ListURL:              ListURL,
		TableURL:             TableURL,
		DetailURL:            DetailURL,
		UnderClientDetailURL: UnderClientDetailURL,
		AddURL:               AddURL,
		EditURL:              EditURL,
		DeleteURL:            DeleteURL,
		BulkDeleteURL:        BulkDeleteURL,
		SetStatusURL:         SetStatusURL,
		BulkSetStatusURL:     BulkSetStatusURL,
		TabActionURL:         TabActionURL,
		SearchPlanURL:        SearchPlanURL,
		SearchClientURL:      SearchClientURL,
		RecognizeURL:         RecognizeURL,
		CustomizePackageURL:  CustomizePackageURL,

		// 2026-04-29 milestone-billing.
		MilestoneMarkReadyURL: MilestoneMarkReadyURL,
		MilestoneWaiveURL:     MilestoneWaiveURL,

		// 20260517-advance-cash-events Plan B Phase 7.
		MilestoneRecognizeURL: MilestoneRecognizeURL,

		// 2026-04-29 auto-spawn-jobs-from-subscription.
		SpawnJobsURL:        SpawnJobsURL,
		SpawnJobsPartialURL: SpawnJobsPartialURL,

		// 2026-04-30 cyclic-subscription-jobs.
		SpawnCycleJobsURL:    SpawnCycleJobsURL,
		BackfillCycleJobsURL: BackfillCycleJobsURL,

		// 2026-05-01 ad-hoc-subscription-billing.
		RequestUsageURL: RequestUsageURL,

		// 2026-05-06 revenue-run — per-subscription drawer.
		RevenueRunURL: RevenueRunURL,

		AttachmentUploadURL:   AttachmentUploadURL,
		AttachmentDeleteURL:   AttachmentDeleteURL,
		AttachmentDownloadURL: AttachmentDownloadURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// subscription routes.
func (r Routes) RouteMap() map[string]string {
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
