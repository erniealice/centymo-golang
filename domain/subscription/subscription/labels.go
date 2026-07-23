package subscription

// ---------------------------------------------------------------------------
// Subscription labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the subscription module.
type Labels struct {
	Page       PageLabels       `json:"page"`
	Buttons    ButtonLabels     `json:"buttons"`
	Columns    ColumnLabels     `json:"columns"`
	Empty      EmptyLabels      `json:"empty"`
	Form       FormLabels       `json:"form"`
	Actions    ActionLabels     `json:"actions"`
	Bulk       BulkLabels       `json:"bulk_actions"`
	Status     StatusLabels     `json:"status"`
	Detail     DetailLabels     `json:"detail"`
	Tabs       TabLabels        `json:"tabs"`
	Invoices   InvoicesLabels   `json:"invoices"`
	Recognize  RecognizeLabels  `json:"recognize"`
	RevenueRun RevenueRunLabels `json:"revenue_run"`
	Milestone  MilestoneLabels  `json:"milestone"`
	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / §9 — Operations
	// tab on the subscription detail page + retroactive spawn drawer copy.
	Operations OperationsLabels `json:"operations"`
	Spawn      SpawnLabels      `json:"spawn"`
	// 2026-04-30 cyclic-subscription-jobs plan §9.2 / §21.3 — Backfill cycle
	// Jobs drawer + flat Jobs tab.
	Backfill BackfillLabels `json:"backfill"`
	Jobs     JobsTabLabels  `json:"jobs"`
	Confirm  ConfirmLabels  `json:"confirm"`
	Errors   ErrorLabels    `json:"errors"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"heading_active"`
	HeadingInactive string `json:"heading_inactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"caption_active"`
	CaptionInactive string `json:"caption_inactive"`
}

type ButtonLabels struct {
	AddSubscription string `json:"add_subscription"`
}

type ColumnLabels struct {
	Name      string `json:"name"`
	Client    string `json:"client"`
	Customer  string `json:"customer"` // legacy alias; kept for backward compat with old translations
	Plan      string `json:"plan"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Status    string `json:"status"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Cancel     string `json:"cancel"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`

	// 2026-04-27 plan-client-scope plan §6.5 / §7 — CTA copy on the
	// subscription detail's Package tab. Templated via {{.ClientName}}.
	CustomizePackage string `json:"customize_package"`
}

type BulkLabels struct {
	Delete     string `json:"delete"`
	Activate   string `json:"bulk_activate"`
	Deactivate string `json:"bulk_deactivate"`
}

type StatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	NoIDsProvided    string `json:"no_ids_provided"`
	InvalidStatus    string `json:"invalid_status"`
	NoPermission     string `json:"no_permission"`
	CannotDelete     string `json:"cannot_delete"`
	InUse            string `json:"in_use"`

	// 2026-04-27 plan-client-scope plan §3.3 / §7 — surfaced when the
	// subscription's selected price_plan belongs to a different client.
	PlanClientMismatch string `json:"plan_client_mismatch"`
	// Surfaced when the customize-package CTA fails (cross-package errors
	// from the espyna use case bubble up here as a generic fallback).
	CustomizeFailed string `json:"customize_failed"`
}

// ---------------------------------------------------------------------------
// Subscription form, detail, tabs, confirm sub-labels
// ---------------------------------------------------------------------------

type FormLabels struct {
	Customer                  string `json:"customer"`
	CustomerPlaceholder       string `json:"customer_placeholder"`
	Plan                      string `json:"plan"`
	PlanPlaceholder           string `json:"plan_placeholder"`
	StartDate                 string `json:"start_date"`
	EndDate                   string `json:"end_date"`
	StartTime                 string `json:"start_time"`
	EndTime                   string `json:"end_time"`
	TimePlaceholder           string `json:"time_placeholder"`
	Timezone                  string `json:"timezone"`
	Active                    string `json:"active"`
	Notes                     string `json:"notes"`
	NotesPlaceholder          string `json:"notes_placeholder"`
	CustomerSearchPlaceholder string `json:"customer_search_placeholder"`
	PlanSearchPlaceholder     string `json:"plan_search_placeholder"`
	CustomerNoResults         string `json:"customer_no_results"`
	PlanNoResults             string `json:"plan_no_results"`
	Code                      string `json:"code"`
	CodePlaceholder           string `json:"code_placeholder"`

	// Field-level info text surfaced via an info button beside each label.
	CustomerInfo  string `json:"customer_info"`
	PlanInfo      string `json:"plan_info"`
	CodeInfo      string `json:"code_info"`
	StartDateInfo string `json:"start_date_info"`
	EndDateInfo   string `json:"end_date_info"`
	StartTimeInfo string `json:"start_time_info"`
	EndTimeInfo   string `json:"end_time_info"`
	NotesInfo     string `json:"notes_info"`

	// 2026-05-03 — Row-level help text rendered below the start/end date+time
	// rows. Explains the operational consequence of the date range (which plans
	// are eligible, when invoicing stops). Distinct from StartDateInfo /
	// EndDateInfo (per-field popovers explaining what each field stores).
	StartDateRowHelp string `json:"start_date_row_help"`
	EndDateRowHelp   string `json:"end_date_row_help"`

	// 2026-04-27 plan-client-scope plan §5.1 / §7 — group headers in the
	// grouped Plan / PricePlan auto-complete picker on the subscription
	// drawer. Templated via {{.ClientName}} for the per-client group.
	PlanGroupForClient string `json:"plan_group_for_client"`
	PlanGroupGeneral   string `json:"plan_group_general"`

	// 2026-05-03 — info banner shown below the locked Customer field on the
	// subscription create drawer, explaining that the Plan picker is
	// scoped to plans assigned to this client (general-scope plans
	// excluded, mirroring the search.go filter).
	PlanClientScopeNotice string `json:"plan_client_scope_notice"`

	// 2026-05-03 — Edit-drawer lock notice rendered when the subscription is
	// referenced by Revenue / subscription_attribute / Job rows. Editing is
	// disabled to preserve the audit trail.
	EditLockedReason string `json:"edit_locked_reason"`

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.1 / §9 — Spawn
	// Jobs toggle section on the subscription create drawer.
	SpawnJobsSectionTitle string `json:"spawn_jobs_section_title"`
	SpawnJobsToggle       string `json:"spawn_jobs_toggle"`
	SpawnJobsHelpText     string `json:"spawn_jobs_help_text"`
	SpawnJobsSummary      string `json:"spawn_jobs_summary"`
	SpawnJobsNone         string `json:"spawn_jobs_none"`

	// Require-spawn-success toggle — sibling of the Spawn Jobs toggle, opts
	// the create request into the fail-closed strict path.
	RequireSpawnToggle string `json:"require_spawn_toggle"`
	RequireSpawnHint   string `json:"require_spawn_hint"`
}

type DetailLabels struct {
	PageTitle            string `json:"page_title"`
	Customer             string `json:"customer"`
	Plan                 string `json:"plan"`
	PriceSchedule        string `json:"price_schedule"`
	StartDate            string `json:"start_date"`
	EndDate              string `json:"end_date"`
	Status               string `json:"status"`
	CreatedDate          string `json:"created_date"`
	ModifiedDate         string `json:"modified_date"`
	AuditTrailComingSoon string `json:"audit_trail_coming_soon"`
	AuditTrailDesc       string `json:"audit_trail_desc"`

	// Info-tab section headers — the info grid renders as three sections:
	// who the subscription is for, the subscription terms, and record logs.
	SectionClient string `json:"section_client"`
	SectionTerms  string `json:"section_terms"`
	SectionLogs   string `json:"section_logs"`

	// Person-name split for the customer section. Rendered instead of the
	// single Customer row when the client has a linked user with name parts.
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	// PackageEmpty is the Package tab's empty state (price plan enrichment
	// missing). Previously this fell through to Invoices.Empty, which read
	// as a nonsensical "No invoices yet" message on the package tab.
	PackageEmpty string `json:"package_empty"`
}

type TabLabels struct {
	Info       string `json:"info"`
	Operations string `json:"operations"`
	// 2026-04-30 cyclic-subscription-jobs plan §21.2 — flat Jobs tab.
	Jobs         string `json:"jobs"`
	Invoices     string `json:"invoices"`
	History      string `json:"history"`
	Attachments  string `json:"attachments"`
	AuditTrail   string `json:"audit_trail"`
	AuditHistory string `json:"audit_history"`
}

type InvoicesLabels struct {
	Title        string `json:"title"`
	Empty        string `json:"empty"`
	ColumnCode   string `json:"column_code"`
	ColumnDate   string `json:"column_date"`
	ColumnAmount string `json:"column_amount"`
	ColumnStatus string `json:"column_status"`

	// Recognize-revenue action surfaced as a primary action on the invoices
	// tab toolbar AND on the empty-state. No page-header button (per plan
	// §11.2 — tab-only).
	RecognizeAction   string `json:"recognize_action"`
	RecognizeTitle    string `json:"recognize_title"`
	RecognizeSubtitle string `json:"recognize_subtitle"`

	// 2026-05-06 revenue-run plan Phase 1 — CTA labels for the three billing-kind
	// branches on the Invoices tab. resolveInvoicesPrimaryAction (Phase 6) picks
	// the correct one; all three must be pre-populated so no branch returns "".
	RunInvoicesAction   string `json:"run_invoices_action"`
	PoolRecognizeAction string `json:"pool_recognize_action"`
	RequestUsageAction  string `json:"request_usage_action"`

	// 2026-05-11 run-invoices-polish Phase 3 — per-row action labels surfaced
	// on the Invoices tab table (view, send email, print, edit).
	RowActions InvoicesRowActionsLabels `json:"row_actions"`
}

// InvoicesRowActionsLabels holds per-row action button labels for
// the Invoices tab on a subscription detail page.
type InvoicesRowActionsLabels struct {
	View      string `json:"view"`
	SendEmail string `json:"send_email"`
	Print     string `json:"print"`
	Edit      string `json:"edit"`
}

// RecognizeLabels holds drawer-form labels for the
// "Recognize Revenue" flow. See plan §5 Phase E for the full table; the
// blocking-error keys (currencyMismatchError, idempotencyError) are renamed
// from their advisory counterparts since v1 surfaces them as hard blocks.
type RecognizeLabels struct {
	// Header / context section
	ContextSection string `json:"context_section"`
	ClientLabel    string `json:"client_label"`
	PlanLabel      string `json:"plan_label"`
	QuantityLabel  string `json:"quantity_label"`

	// Period section
	PeriodSection string `json:"period_section"`
	PeriodStart   string `json:"period_start"`
	PeriodEnd     string `json:"period_end"`
	RevenueDate   string `json:"revenue_date"`

	// Line items table
	LineItemsSection    string `json:"line_items_section"`
	ColumnDescription   string `json:"column_description"`
	ColumnUnitPrice     string `json:"column_unit_price"`
	ColumnQuantity      string `json:"column_quantity"`
	ColumnLineTotal     string `json:"column_line_total"`
	ColumnTreatment     string `json:"column_treatment"`
	TotalLabel          string `json:"total_label"`
	RemoveLine          string `json:"remove_line"`
	TreatmentRecurring  string `json:"treatment_recurring"`
	TreatmentFirstCycle string `json:"treatment_first_cycle"`
	TreatmentUsageBased string `json:"treatment_usage_based"`
	TreatmentOneTime    string `json:"treatment_one_time"`

	// Notes
	NotesLabel       string `json:"notes_label"`
	NotesPlaceholder string `json:"notes_placeholder"`

	// Footer buttons (v1 — single Generate button; "Save as Draft" is dropped
	// per plan Phase D refinement since both paths run the idempotency check.)
	Generate string `json:"generate"`
	Cancel   string `json:"cancel"`

	// Blocking error banners
	CurrencyMismatchError     string `json:"currency_mismatch_error"`
	IdempotencyError          string `json:"idempotency_error"`
	IdempotencyExistingLink   string `json:"idempotency_existing_link"`
	NoLinesError              string `json:"no_lines_error"`
	CycleNotConfiguredWarning string `json:"cycle_not_configured_warning"`
	UsageBasedSkippedNotice   string `json:"usage_based_skipped_notice"`

	// 2026-04-27 plan-client-scope plan §7 — info notice on the recognize
	// drawer when the active subscription's PricePlan is client-scoped.
	// Templated via {{.ClientName}}.
	ClientCustomNotice string `json:"client_custom_notice"`

	// 2026-04-29 milestone-billing plan §5 / Phase E — milestone-specific
	// drawer fields. Surfaced only when pricePlan.billing_kind = MILESTONE.
	MilestoneSelect            string `json:"milestone_select"`
	MilestoneSelectPlaceholder string `json:"milestone_select_placeholder"`
	NoReadyMilestone           string `json:"no_ready_milestone"`
	MilestoneNotApplicable     string `json:"milestone_not_applicable"`
	BillAmount                 string `json:"bill_amount"`
	LeaveRemainderOpen         string `json:"leave_remainder_open"`
	CloseShort                 string `json:"close_short"`
	PartialReason              string `json:"partial_reason"`
	PartialReasonRequired      string `json:"partial_reason_required"`
	OverBillingRejected        string `json:"over_billing_rejected"`

	// Tax preview labels (Phase 5)
	TaxPreviewSection       string `json:"tax_preview_section"`
	TaxDirectionSurcharge   string `json:"tax_direction_surcharge"`
	TaxDirectionWithholding string `json:"tax_direction_withholding"`
	NetReceivable           string `json:"net_receivable"`
	WHTAmount               string `json:"wht_amount"`
	// TaxKindLabels maps tax_kind_snapshot values to localized display names.
	// Populated from lyngua; used by convertPreviewTaxLines in the recognize view.
	TaxKindLabels map[string]string `json:"tax_kind_labels"`
}

// RevenueRunLabels holds drawer-form labels for the per-subscription
// "Invoice Run" drawer (Surface C — CYCLE billing_kind only). Lyngua key:
// `subscription.revenueRun`. Drops the engagement column that the client-level
// drawer (Surface A) carries — per-sub context makes it redundant.
type RevenueRunLabels struct {
	Title string `json:"title"`
	// Subtitle is templated with the subscription name — e.g. "Run invoices for {{.Name}}"
	Subtitle string `json:"subtitle"`

	// Read-only context row labels (subscription name + plan name)
	SubscriptionLabel string `json:"subscription_label"`
	PlanLabel         string `json:"plan_label"`
	// ClientHintTemplate is shown as a hint beneath the subscription field.
	// Use {client} as the substitution token. E.g. "Client: {client}".
	ClientHintTemplate string `json:"client_hint_template"`

	AsOfDateLabel         string `json:"as_of_date_label"`
	AsOfDateHint          string `json:"as_of_date_hint"`
	BillThroughTodayLabel string `json:"bill_through_today_label"`

	// Period table columns
	ColumnPeriod string `json:"column_period"`
	ColumnAmount string `json:"column_amount"`
	ColumnLines  string `json:"column_lines"`

	// Group headings / empty states
	GroupNoPending        string `json:"group_no_pending"`
	GroupCurrencyMismatch string `json:"group_currency_mismatch"`
	EmptyTitle            string `json:"empty_title"`
	EmptyMessage          string `json:"empty_message"`

	// IntroMessage is shown at the top of the drawer body as an info alert.
	IntroMessage string `json:"intro_message"`

	// Footer buttons
	GenerateButton          string `json:"generate_button"`
	GenerateButtonCountOne  string `json:"generate_button_count_one"`
	GenerateButtonCountMany string `json:"generate_button_count_many"`
	CancelButton            string `json:"cancel_button"`

	// Post-submit feedback
	ToastSuccess string `json:"toast_success"`
	ToastSkipped string `json:"toast_skipped"`
	ToastErrored string `json:"toast_errored"`
	ViewRunLink  string `json:"view_run_link"`

	// Inline error messages
	Errors RevenueRunErrorLabels `json:"errors"`
}

// RevenueRunErrorLabels holds inline error strings for the
// per-subscription Invoice Run drawer.
type RevenueRunErrorLabels struct {
	PermissionDenied   string `json:"permission_denied"`
	IDRequired         string `json:"id_required"`
	InvalidFormData    string `json:"invalid_form_data"`
	UseCaseUnavailable string `json:"use_case_unavailable"`
	SelectOne          string `json:"select_one"`
}

// MilestoneLabels holds labels for the Subscription Package tab's
// Milestones section + the mark-ready / waive CTAs. Lyngua key:
// `subscription.milestone.*`. See milestone-billing plan §5.
type MilestoneLabels struct {
	Title           string `json:"title"`
	Subtitle        string `json:"subtitle"`
	MarkReady       string `json:"mark_ready"`
	Waive           string `json:"waive"`
	ViewInvoice     string `json:"view_invoice"`
	StatusPending   string `json:"status_pending"`
	StatusReady     string `json:"status_ready"`
	StatusBilled    string `json:"status_billed"`
	StatusWaived    string `json:"status_waived"`
	StatusDeferred  string `json:"status_deferred"`
	StatusCancelled string `json:"status_cancelled"`
	TotalInvoiced   string `json:"total_invoiced"`
	AmountFull      string `json:"amount_full"`
	AmountPartial   string `json:"amount_partial"`

	// 20260517-advance-cash-events Plan B Phase 7 — Recognize CTA + the
	// "linked to advance" badge that flags milestones tied to an advance
	// Collection (via the collection_billing_event junction).
	Recognize          string `json:"recognize"`
	LinkedAdvanceBadge string `json:"linked_advance_badge"`
}

// OperationsLabels holds labels for the Subscription detail's
// Operations tab. Lyngua key: `subscription.detail.operations.*`. See
// auto-spawn-jobs-from-subscription plan §5.2 / §9 and cyclic-subscription-jobs
// plan §9.1 (cycle accordion + backfill keys).
type OperationsLabels struct {
	Title        string `json:"title"`
	EmptyTitle   string `json:"empty_title"`
	EmptyMessage string `json:"empty_message"`
	SpawnAction  string `json:"spawn_action"`
	RootJob      string `json:"root_job"`
	ChildJob     string `json:"child_job"`
	PhaseSummary string `json:"phase_summary"`
	ViewJobLink  string `json:"view_job_link"`

	// 2026-04-30 cyclic-subscription-jobs plan §9.1 — cycle accordion copy.
	SubscriptionHeading   string `json:"subscription_heading"`
	CycleHeading          string `json:"cycle_heading"`
	CyclePlaceholder      string `json:"cycle_placeholder"`
	CycleSpawnNow         string `json:"cycle_spawn_now"`
	CycleStatusPending    string `json:"cycle_status_pending"`
	CycleStatusInProgress string `json:"cycle_status_in_progress"`
	CycleStatusCompleted  string `json:"cycle_status_completed"`
	CycleStatusOverdue    string `json:"cycle_status_overdue"`
	CycleInvoiceLinked    string `json:"cycle_invoice_linked"`
	CycleNoInvoice        string `json:"cycle_no_invoice"`
	CycleEmpty            string `json:"cycle_empty"`
	BackfillBanner        string `json:"backfill_banner"`
	BackfillCta           string `json:"backfill_cta"`

	// 2026-05-01 ad-hoc-subscription-billing plan §5.2 — Operations tab
	// AD_HOC mode keys. Vertical-neutral defaults ("usage", "occurrence")
	// with professional-tier overrides ("service call", "retainer", etc.).
	AdHocPoolHeading       string `json:"ad_hoc_pool_heading"`
	AdHocPerCallHeading    string `json:"ad_hoc_per_call_heading"`
	EntitlementUsed        string `json:"entitlement_used"`
	EntitlementRemaining   string `json:"entitlement_remaining"`
	EntitlementExhausted   string `json:"entitlement_exhausted"`
	RequestUsageCta        string `json:"request_usage_cta"`
	ExtendEntitlementCta   string `json:"extend_entitlement_cta"`
	UsageRequestedDate     string `json:"usage_requested_date"`
	UsageDeliveredDate     string `json:"usage_delivered_date"`
	UsageOrdinalLabel      string `json:"usage_ordinal_label"`
	UsageNotDelivered      string `json:"usage_not_delivered"`
	PoolInvoiceLink        string `json:"pool_invoice_link"`
	PoolInvoicePending     string `json:"pool_invoice_pending"`
	PoolGenerateInvoiceCta string `json:"pool_generate_invoice_cta"`
	PerCallRecognizeCta    string `json:"per_call_recognize_cta"`
	PerCallInvoiceLink     string `json:"per_call_invoice_link"`
	PerCallNotReady        string `json:"per_call_not_ready"`
}

// BackfillLabels holds labels for the Backfill cycle Jobs drawer.
// Lyngua key: `subscription.detail.backfill.*`. See cyclic-subscription-jobs
// plan §9.2.
type BackfillLabels struct {
	DrawerTitle       string `json:"drawer_title"`
	DrawerDescription string `json:"drawer_description"`
	PreviewLine       string `json:"preview_line"`
	CountLabel        string `json:"count_label"`
	Confirm           string `json:"confirm"`
	Cancel            string `json:"cancel"`
	MaxWarning        string `json:"max_warning"`
}

// JobsTabLabels holds labels for the new flat Jobs tab on the
// Subscription detail page. Lyngua key: `subscription.detail.jobs.*`. See
// cyclic-subscription-jobs plan §21.
type JobsTabLabels struct {
	Heading          string `json:"heading"`
	Empty            string `json:"empty"`
	FilterStatus     string `json:"filter_status"`
	FilterType       string `json:"filter_type"`
	FilterAll        string `json:"filter_all"`
	SortBy           string `json:"sort_by"`
	SortByCycle      string `json:"sort_by_cycle"`
	ExportCsv        string `json:"export_csv"`
	Summary          string `json:"summary"`
	ColumnNumber     string `json:"column_number"`
	ColumnName       string `json:"column_name"`
	ColumnType       string `json:"column_type"`
	ColumnPhase      string `json:"column_phase"`
	ColumnStatus     string `json:"column_status"`
	ColumnPeriod     string `json:"column_period"`
	TypeSubscription string `json:"type_subscription"`
	TypeOnboarding   string `json:"type_onboarding"`
	TypeCycle        string `json:"type_cycle"`
	TypeVisit        string `json:"type_visit"`
	SpawnFailedToast string `json:"spawn_failed_toast"`
}

// SpawnLabels holds labels for the retroactive Spawn Jobs drawer.
// Lyngua key: `subscription.spawn.*`. See auto-spawn-jobs-from-subscription
// plan §5.3 / §9.
type SpawnLabels struct {
	Title             string `json:"title"`
	DetectedTemplates string `json:"detected_templates"`
	RootTemplate      string `json:"root_template"`
	Cancel            string `json:"cancel"`
	Confirm           string `json:"confirm"`
	SuccessToast      string `json:"success_toast"`
	Skipped           string `json:"skipped"`

	// Blocked is the retroactive-spawn guard message — shown as a blocking
	// GET-drawer banner and returned as the POST rejection error when the
	// subscription's resolved price_schedule is closed or inactive.
	Blocked string `json:"blocked"`
}

type ConfirmLabels struct {
	Cancel                string `json:"cancel"`
	CancelMessage         string `json:"cancel_message"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"delete_message"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activate_message"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivate_message"`
	BulkActivate          string `json:"bulk_activate"`
	BulkActivateMessage   string `json:"bulk_activate_message"`
	BulkDeactivate        string `json:"bulk_deactivate"`
	BulkDeactivateMessage string `json:"bulk_deactivate_message"`
	BulkDelete            string `json:"bulk_delete"`
	BulkDeleteMessage     string `json:"bulk_delete_message"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:         "Subscriptions",
			HeadingActive:   "Active Subscriptions",
			HeadingInactive: "Inactive Subscriptions",
			Caption:         "Subscription management",
			CaptionActive:   "Manage your active subscriptions",
			CaptionInactive: "View cancelled or expired subscriptions",
		},
		Buttons: ButtonLabels{
			AddSubscription: "Add Subscription",
		},
		Columns: ColumnLabels{
			Name:      "Engagement",
			Client:    "Client",
			Customer:  "Customer",
			Plan:      "Plan",
			StartDate: "Start Date",
			EndDate:   "End Date",
			Status:    "Status",
		},
		Empty: EmptyLabels{
			Title:   "No subscriptions found",
			Message: "No subscriptions to display.",
		},
		Form: FormLabels{
			Customer:                  "Customer",
			CustomerPlaceholder:       "Select customer...",
			Plan:                      "Plan",
			PlanPlaceholder:           "Select plan...",
			StartDate:                 "Start Date",
			EndDate:                   "End Date",
			StartTime:                 "Start Time (optional)",
			EndTime:                   "End Time (optional)",
			TimePlaceholder:           "HH:MM",
			Timezone:                  "Timezone",
			Active:                    "Active",
			Notes:                     "Notes",
			NotesPlaceholder:          "Enter notes...",
			CustomerSearchPlaceholder: "Search customers...",
			PlanSearchPlaceholder:     "Search plans...",
			CustomerNoResults:         "No customers found",
			PlanNoResults:             "No plans found",
			Code:                      "Code",
			CodePlaceholder:           "e.g. A3K7PXR",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			CustomerInfo:     "The client this subscription is billed to.",
			PlanInfo:         "The price plan this subscription follows. Determines amount, billing cycle, and any per-product prices.",
			CodeInfo:         "Short reference used on invoices and receipts. Leave blank to auto-generate.",
			StartDateInfo:    "First day the subscription is active. Billing cycles are counted from this date.",
			EndDateInfo:      "Last day the subscription is active. Leave blank for open-ended.",
			StartDateRowHelp: "Start date and time affect which plans are available below — only plans active in this date range can be selected.",
			EndDateRowHelp:   "End date and time control when recurring invoices stop being issued for this subscription. Leave blank for open-ended billing.",
			StartTimeInfo:    "Optional time of day in the operator's display timezone. Leave blank for start of day (00:00).",
			EndTimeInfo:      "Optional time of day in the operator's display timezone. Leave blank for end of day (23:59).",
			NotesInfo:        "Internal remarks — shown on detail pages but not on customer-facing documents.",
			// 2026-04-27 plan-client-scope plan §5.1 / §7 — grouped picker headers.
			PlanGroupForClient:    "For {{.ClientName}}",
			PlanGroupGeneral:      "General packages",
			PlanClientScopeNotice: "Plans below match this client's billing currency ({{.Currency}}).",
			EditLockedReason:      "This subscription has revenue records and cannot be edited. Reassigning the plan would break the audit trail.",
			// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.1 / §9 —
			// Spawn Jobs toggle on subscription create drawer.
			SpawnJobsSectionTitle: "Operations",
			SpawnJobsToggle:       "Spawn Job(s) on Create",
			SpawnJobsHelpText:     "Disable to start without operational tracking (e.g., advisory retainers).",
			SpawnJobsSummary:      "Spawning {{.JobCount}} Job(s) from {{.TemplateNames}} — includes {{.PhaseCount}} phases, {{.TaskCount}} tasks.",
			SpawnJobsNone:         "No JobTemplate is configured for this Plan. The engagement will start without operational tracking.",
			RequireSpawnToggle:    "Require successful job spawn",
			RequireSpawnHint:      "If spawning jobs fails or produces none, the whole subscription create is cancelled instead of saving without jobs.",
		},
		Actions: ActionLabels{
			View:       "View Subscription",
			Edit:       "Edit Subscription",
			Cancel:     "Cancel Subscription",
			Delete:     "Delete",
			Activate:   "Activate",
			Deactivate: "Deactivate",
			// 2026-04-27 plan-client-scope plan §6.5 / §7 — Package tab CTA.
			CustomizePackage: "Customize this package for {{.ClientName}}",
		},
		Bulk: BulkLabels{
			Delete:     "Delete Selected",
			Activate:   "Activate Selected",
			Deactivate: "Deactivate Selected",
		},
		Status: StatusLabels{
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Detail: DetailLabels{
			PageTitle:            "Subscription Details",
			Customer:             "Customer",
			Plan:                 "Plan",
			PriceSchedule:        "Price Schedule",
			StartDate:            "Start Date",
			EndDate:              "End Date",
			Status:               "Status",
			CreatedDate:          "Created",
			ModifiedDate:         "Last Modified",
			AuditTrailComingSoon: "Audit trail coming soon.",
			AuditTrailDesc:       "Audit trail for subscription changes is coming soon.",
			SectionClient:        "Client",
			SectionTerms:         "Subscription",
			SectionLogs:          "Logs",
			FirstName:            "First Name",
			LastName:             "Last Name",
			PackageEmpty:         "No package details available.",
		},
		Tabs: TabLabels{
			Info:       "Information",
			Operations: "Operations",
			// 2026-04-30 cyclic-subscription-jobs plan §21.2 — flat Jobs tab.
			Jobs:         "Jobs",
			Invoices:     "Invoices",
			History:      "History",
			Attachments:  "Attachments",
			AuditTrail:   "Audit Trail",
			AuditHistory: "History",
		},
		Invoices: InvoicesLabels{
			Title:             "Invoices",
			Empty:             "No invoices yet — click Recognize Revenue to generate the first one.",
			ColumnCode:        "Number",
			ColumnDate:        "Date",
			ColumnAmount:      "Amount",
			ColumnStatus:      "Status",
			RecognizeAction:   "Recognize Revenue",
			RecognizeTitle:    "Recognize Revenue",
			RecognizeSubtitle: "Generate an invoice from this subscription's price plan.",
			// 2026-05-06 revenue-run plan Phase 1 — CTA branch labels.
			RunInvoicesAction:   "Run Invoices",
			PoolRecognizeAction: "Recognize Revenue",
			RequestUsageAction:  "Request Usage",
			// 2026-05-11 run-invoices-polish Phase 3 — row action defaults.
			RowActions: InvoicesRowActionsLabels{
				View:      "View",
				SendEmail: "Send email",
				Print:     "Print",
				Edit:      "Edit",
			},
		},
		RevenueRun: RevenueRunLabels{
			Title:                   "Invoice Run",
			Subtitle:                "Run invoices for {{.Name}}",
			SubscriptionLabel:       "Engagement",
			PlanLabel:               "Plan",
			ClientHintTemplate:      "Client: {client}",
			AsOfDateLabel:           "As of date",
			AsOfDateHint:            "Only periods ending on or before this date will be included.",
			BillThroughTodayLabel:   "Bill through today",
			ColumnPeriod:            "Period",
			ColumnAmount:            "Amount",
			ColumnLines:             "Lines",
			GroupNoPending:          "No pending periods",
			GroupCurrencyMismatch:   "Currency mismatch — cannot run",
			EmptyTitle:              "Nothing to invoice",
			EmptyMessage:            "This subscription has no pending billing periods as of the selected date.",
			IntroMessage:            "Generate invoices for the eligible billing periods of this engagement. Select the periods you want to bill below — each will produce a separate invoice.",
			GenerateButton:          "Generate",
			GenerateButtonCountOne:  "Generate {count} Invoice",
			GenerateButtonCountMany: "Generate {count} Invoices",
			CancelButton:            "Cancel",
			ToastSuccess:            "Invoice run complete — {{.Created}} invoice(s) created.",
			ToastSkipped:            "Invoice run complete — all periods skipped.",
			ToastErrored:            "Invoice run completed with errors — {{.Errored}} period(s) failed.",
			ViewRunLink:             "View run",
			Errors: RevenueRunErrorLabels{
				PermissionDenied:   "You do not have permission to run invoices.",
				IDRequired:         "Subscription ID is required.",
				InvalidFormData:    "Invalid form data. Please check your inputs and try again.",
				UseCaseUnavailable: "Invoice run is not available for this subscription type.",
				SelectOne:          "Select at least one period to generate.",
			},
		},
		Recognize: RecognizeLabels{
			ContextSection:            "Subscription",
			ClientLabel:               "Client",
			PlanLabel:                 "Plan / Rate Card",
			QuantityLabel:             "Quantity",
			PeriodSection:             "Billing period",
			PeriodStart:               "Period start",
			PeriodEnd:                 "Period end",
			RevenueDate:               "Revenue date",
			LineItemsSection:          "Line items",
			ColumnDescription:         "Description",
			ColumnUnitPrice:           "Unit price",
			ColumnQuantity:            "Qty",
			ColumnLineTotal:           "Line total",
			ColumnTreatment:           "Treatment",
			TotalLabel:                "Total",
			RemoveLine:                "Remove",
			TreatmentRecurring:        "Every cycle",
			TreatmentFirstCycle:       "First cycle only",
			TreatmentUsageBased:       "On use",
			TreatmentOneTime:          "One time",
			NotesLabel:                "Notes",
			NotesPlaceholder:          "Notes are auto-prefixed with the period; append any free-text below.",
			Generate:                  "Generate",
			Cancel:                    "Cancel",
			CurrencyMismatchError:     "Client billing currency ({{.ClientCurrency}}) does not match the rate card ({{.PlanCurrency}}). Update one before generating revenue.",
			IdempotencyError:          "An invoice for this period already exists. Cancel the existing one or pick a different period.",
			IdempotencyExistingLink:   "View the existing invoice",
			NoLinesError:              "Cannot create an invoice with no line items. Add a price plan with at least one product, or override at least one line.",
			CycleNotConfiguredWarning: "Plan has no billing cycle configured; defaulting to 1 month.",
			UsageBasedSkippedNotice:   "Usage-based lines were skipped — record them via metering.",
			// 2026-04-27 plan-client-scope plan §7 — surfaced when the
			// active subscription's PricePlan is client-scoped.
			ClientCustomNotice: "This engagement uses a custom package for {{.ClientName}}.",
			// 2026-04-29 milestone-billing plan §5 / Phase E.
			MilestoneSelect:            "Milestone",
			MilestoneSelectPlaceholder: "Select a ready milestone",
			NoReadyMilestone:           "No milestone is ready to bill.",
			MilestoneNotApplicable:     "Milestones are only available on milestone-priced plans.",
			BillAmount:                 "Bill amount",
			LeaveRemainderOpen:         "Partial — leave remainder open",
			CloseShort:                 "Partial — close milestone short",
			PartialReason:              "Reason",
			PartialReasonRequired:      "A reason is required when billing partially.",
			OverBillingRejected:        "Cannot bill: total would exceed milestone amount.",
		},
		Milestone: MilestoneLabels{
			Title:           "Billing Schedule",
			Subtitle:        "Milestone events for this engagement",
			MarkReady:       "Mark Ready",
			Waive:           "Waive",
			ViewInvoice:     "View Invoice",
			StatusPending:   "Pending",
			StatusReady:     "Ready",
			StatusBilled:    "Billed",
			StatusWaived:    "Waived",
			StatusDeferred:  "Deferred",
			StatusCancelled: "Cancelled",
			TotalInvoiced:   "Total Invoiced",
			AmountFull:      "Full amount",
			AmountPartial:   "Partial — {{.Billed}} of {{.Full}}",
			// 20260517-advance-cash-events Plan B Phase 7.
			Recognize:          "Recognize",
			LinkedAdvanceBadge: "Linked advance",
		},
		// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.2 / §9 +
		// 2026-04-30 cyclic-subscription-jobs plan §9.1 — cycle accordion.
		Operations: OperationsLabels{
			Title:        "Operational Jobs",
			EmptyTitle:   "No operational tracking",
			EmptyMessage: "This engagement has no Jobs. {{.SpawnAction}} to start tracking work.",
			SpawnAction:  "Spawn Jobs",
			RootJob:      "Root Job",
			ChildJob:     "Child Job",
			PhaseSummary: "{{.Complete}} / {{.Total}} phases complete",
			ViewJobLink:  "View in Operations",

			// Cycle accordion + backfill copy.
			SubscriptionHeading:   "Engagement (since {{.Started}})",
			CycleHeading:          "Cycle {{.CycleIndex}} — {{.PeriodLabel}}",
			CyclePlaceholder:      "Cycle starts {{.PeriodStart}} — Jobs will spawn at cycle start, or click below to spawn now.",
			CycleSpawnNow:         "Spawn this cycle now",
			CycleStatusPending:    "Pending",
			CycleStatusInProgress: "In progress",
			CycleStatusCompleted:  "Completed",
			CycleStatusOverdue:    "Overdue",
			CycleInvoiceLinked:    "Invoice {{.RevenueCode}} · {{.Status}}",
			CycleNoInvoice:        "Not yet invoiced",
			CycleEmpty:            "No cycles yet",
			BackfillBanner:        "{{.Count}} cycle(s) missing operational tracking. Spawn now to backfill.",
			BackfillCta:           "Backfill missing cycles",
		},
		// 2026-04-30 cyclic-subscription-jobs plan §9.2 — backfill drawer.
		Backfill: BackfillLabels{
			DrawerTitle:       "Backfill cycle Jobs",
			DrawerDescription: "Preview the cycles that will be spawned, then confirm to materialize them in one transaction.",
			PreviewLine:       "Cycle {{.Index}} — {{.PeriodLabel}}",
			CountLabel:        "Cycles to spawn",
			Confirm:           "Spawn {{.Count}} cycle(s)",
			Cancel:            "Cancel",
			MaxWarning:        "Backfill is capped at 24 cycles per request. Reduce the range or run multiple backfills.",
		},
		// 2026-04-30 cyclic-subscription-jobs plan §21.3 — flat Jobs tab.
		Jobs: JobsTabLabels{
			Heading:          "Jobs",
			Empty:            "No Jobs yet — this engagement has no operational tracking.",
			FilterStatus:     "Status",
			FilterType:       "Type",
			FilterAll:        "All",
			SortBy:           "Sort",
			SortByCycle:      "Cycle #",
			ExportCsv:        "Export CSV",
			Summary:          "Showing {{.Visible}} of {{.Total}} Jobs",
			ColumnNumber:     "#",
			ColumnName:       "Job Name",
			ColumnType:       "Type",
			ColumnPhase:      "Phase",
			ColumnStatus:     "Status",
			ColumnPeriod:     "Period",
			TypeSubscription: "Engagement",
			TypeOnboarding:   "Onboarding",
			TypeCycle:        "Cycle",
			TypeVisit:        "Visit",
			SpawnFailedToast: "Cycle Job spawn failed for {{.Period}} — invoice was created but the operational Job will need a manual retry.",
		},
		Spawn: SpawnLabels{
			Title:             "Spawn Operational Jobs",
			DetectedTemplates: "Detected templates",
			RootTemplate:      "Root template",
			Cancel:            "Cancel",
			Confirm:           "Spawn Jobs",
			SuccessToast:      "Spawned {{.JobCount}} Job(s).",
			Skipped:           "Nothing to spawn — no JobTemplate is linked to this Plan.",
			Blocked:           "Job spawn is blocked — the price schedule for this subscription's plan is closed or inactive.",
		},
		Confirm: ConfirmLabels{
			Cancel:                "Cancel Subscription",
			CancelMessage:         "Are you sure you want to cancel this subscription? This action cannot be undone.",
			Delete:                "Delete Subscription",
			DeleteMessage:         "Are you sure you want to delete this subscription? This action cannot be undone.",
			Activate:              "Activate Subscription",
			ActivateMessage:       "Are you sure you want to activate %s?",
			Deactivate:            "Deactivate Subscription",
			DeactivateMessage:     "Are you sure you want to deactivate %s?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Are you sure you want to activate the selected subscriptions?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Are you sure you want to deactivate the selected subscriptions?",
			BulkDelete:            "Delete Selected",
			BulkDeleteMessage:     "Are you sure you want to delete the selected subscriptions? This action cannot be undone.",
		},
		Errors: ErrorLabels{
			PermissionDenied:   "You do not have permission to perform this action",
			InvalidFormData:    "Invalid form data. Please check your inputs and try again.",
			NotFound:           "Subscription not found",
			IDRequired:         "Subscription ID is required",
			NoIDsProvided:      "No subscription IDs provided",
			InvalidStatus:      "Invalid status value",
			NoPermission:       "No permission",
			CannotDelete:       "Cannot delete — this engagement has dependent records",
			InUse:              "Cannot delete — this engagement has dependent records (jobs, revenue, invoices, etc.)",
			PlanClientMismatch: "This package belongs to a different client and cannot be attached here.",
			CustomizeFailed:    "Failed to customize this package. Please try again.",
		},
	}
}
