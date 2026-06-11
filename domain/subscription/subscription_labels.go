package subscription

// ---------------------------------------------------------------------------
// Subscription labels
// ---------------------------------------------------------------------------

// SubscriptionLabels holds all translatable strings for the subscription module.
type SubscriptionLabels struct {
	Page       SubscriptionPageLabels       `json:"page"`
	Buttons    SubscriptionButtonLabels     `json:"buttons"`
	Columns    SubscriptionColumnLabels     `json:"columns"`
	Empty      SubscriptionEmptyLabels      `json:"empty"`
	Form       SubscriptionFormLabels       `json:"form"`
	Actions    SubscriptionActionLabels     `json:"actions"`
	Bulk       SubscriptionBulkLabels       `json:"bulkActions"`
	Status     SubscriptionStatusLabels     `json:"status"`
	Detail     SubscriptionDetailLabels     `json:"detail"`
	Tabs       SubscriptionTabLabels        `json:"tabs"`
	Invoices   SubscriptionInvoicesLabels   `json:"invoices"`
	Recognize  SubscriptionRecognizeLabels  `json:"recognize"`
	RevenueRun SubscriptionRevenueRunLabels `json:"revenueRun"`
	Milestone  SubscriptionMilestoneLabels  `json:"milestone"`
	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / §9 — Operations
	// tab on the subscription detail page + retroactive spawn drawer copy.
	Operations SubscriptionOperationsLabels `json:"operations"`
	Spawn      SubscriptionSpawnLabels      `json:"spawn"`
	// 2026-04-30 cyclic-subscription-jobs plan §9.2 / §21.3 — Backfill cycle
	// Jobs drawer + flat Jobs tab.
	Backfill SubscriptionBackfillLabels `json:"backfill"`
	Jobs     SubscriptionJobsTabLabels  `json:"jobs"`
	Confirm  SubscriptionConfirmLabels  `json:"confirm"`
	Errors   SubscriptionErrorLabels    `json:"errors"`
}

type SubscriptionPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type SubscriptionButtonLabels struct {
	AddSubscription string `json:"addSubscription"`
}

type SubscriptionColumnLabels struct {
	Name      string `json:"name"`
	Client    string `json:"client"`
	Customer  string `json:"customer"` // legacy alias; kept for backward compat with old translations
	Plan      string `json:"plan"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Status    string `json:"status"`
}

type SubscriptionEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SubscriptionActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Cancel     string `json:"cancel"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`

	// 2026-04-27 plan-client-scope plan §6.5 / §7 — CTA copy on the
	// subscription detail's Package tab. Templated via {{.ClientName}}.
	CustomizePackage string `json:"customizePackage"`
}

type SubscriptionBulkLabels struct {
	Delete     string `json:"delete"`
	Activate   string `json:"bulkActivate"`
	Deactivate string `json:"bulkDeactivate"`
}

type SubscriptionStatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type SubscriptionErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
	NoPermission     string `json:"noPermission"`
	CannotDelete     string `json:"cannotDelete"`
	InUse            string `json:"inUse"`

	// 2026-04-27 plan-client-scope plan §3.3 / §7 — surfaced when the
	// subscription's selected price_plan belongs to a different client.
	PlanClientMismatch string `json:"planClientMismatch"`
	// Surfaced when the customize-package CTA fails (cross-package errors
	// from the espyna use case bubble up here as a generic fallback).
	CustomizeFailed string `json:"customizeFailed"`
}

// ---------------------------------------------------------------------------
// Subscription form, detail, tabs, confirm sub-labels
// ---------------------------------------------------------------------------

type SubscriptionFormLabels struct {
	Customer                  string `json:"customer"`
	CustomerPlaceholder       string `json:"customerPlaceholder"`
	Plan                      string `json:"plan"`
	PlanPlaceholder           string `json:"planPlaceholder"`
	StartDate                 string `json:"startDate"`
	EndDate                   string `json:"endDate"`
	StartTime                 string `json:"startTime"`
	EndTime                   string `json:"endTime"`
	TimePlaceholder           string `json:"timePlaceholder"`
	Timezone                  string `json:"timezone"`
	Active                    string `json:"active"`
	Notes                     string `json:"notes"`
	NotesPlaceholder          string `json:"notesPlaceholder"`
	CustomerSearchPlaceholder string `json:"customerSearchPlaceholder"`
	PlanSearchPlaceholder     string `json:"planSearchPlaceholder"`
	CustomerNoResults         string `json:"customerNoResults"`
	PlanNoResults             string `json:"planNoResults"`
	Code                      string `json:"code"`
	CodePlaceholder           string `json:"codePlaceholder"`

	// Field-level info text surfaced via an info button beside each label.
	CustomerInfo  string `json:"customerInfo"`
	PlanInfo      string `json:"planInfo"`
	CodeInfo      string `json:"codeInfo"`
	StartDateInfo string `json:"startDateInfo"`
	EndDateInfo   string `json:"endDateInfo"`
	StartTimeInfo string `json:"startTimeInfo"`
	EndTimeInfo   string `json:"endTimeInfo"`
	NotesInfo     string `json:"notesInfo"`

	// 2026-05-03 — Row-level help text rendered below the start/end date+time
	// rows. Explains the operational consequence of the date range (which plans
	// are eligible, when invoicing stops). Distinct from StartDateInfo /
	// EndDateInfo (per-field popovers explaining what each field stores).
	StartDateRowHelp string `json:"startDateRowHelp"`
	EndDateRowHelp   string `json:"endDateRowHelp"`

	// 2026-04-27 plan-client-scope plan §5.1 / §7 — group headers in the
	// grouped Plan / PricePlan auto-complete picker on the subscription
	// drawer. Templated via {{.ClientName}} for the per-client group.
	PlanGroupForClient string `json:"planGroupForClient"`
	PlanGroupGeneral   string `json:"planGroupGeneral"`

	// 2026-05-03 — info banner shown below the locked Customer field on the
	// subscription create drawer, explaining that the Plan picker is
	// scoped to plans assigned to this client (general-scope plans
	// excluded, mirroring the search.go filter).
	PlanClientScopeNotice string `json:"planClientScopeNotice"`

	// 2026-05-03 — Edit-drawer lock notice rendered when the subscription is
	// referenced by Revenue / subscription_attribute / Job rows. Editing is
	// disabled to preserve the audit trail.
	EditLockedReason string `json:"editLockedReason"`

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.1 / §9 — Spawn
	// Jobs toggle section on the subscription create drawer.
	SpawnJobsSectionTitle string `json:"spawnJobsSectionTitle"`
	SpawnJobsToggle       string `json:"spawnJobsToggle"`
	SpawnJobsHelpText     string `json:"spawnJobsHelpText"`
	SpawnJobsSummary      string `json:"spawnJobsSummary"`
	SpawnJobsNone         string `json:"spawnJobsNone"`
}

type SubscriptionDetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	Customer             string `json:"customer"`
	Plan                 string `json:"plan"`
	StartDate            string `json:"startDate"`
	EndDate              string `json:"endDate"`
	Status               string `json:"status"`
	CreatedDate          string `json:"createdDate"`
	ModifiedDate         string `json:"modifiedDate"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditTrailDesc       string `json:"auditTrailDesc"`
}

type SubscriptionTabLabels struct {
	Info       string `json:"info"`
	Operations string `json:"operations"`
	// 2026-04-30 cyclic-subscription-jobs plan §21.2 — flat Jobs tab.
	Jobs         string `json:"jobs"`
	Invoices     string `json:"invoices"`
	History      string `json:"history"`
	Attachments  string `json:"attachments"`
	AuditTrail   string `json:"auditTrail"`
	AuditHistory string `json:"auditHistory"`
}

type SubscriptionInvoicesLabels struct {
	Title        string `json:"title"`
	Empty        string `json:"empty"`
	ColumnCode   string `json:"columnCode"`
	ColumnDate   string `json:"columnDate"`
	ColumnAmount string `json:"columnAmount"`
	ColumnStatus string `json:"columnStatus"`

	// Recognize-revenue action surfaced as a primary action on the invoices
	// tab toolbar AND on the empty-state. No page-header button (per plan
	// §11.2 — tab-only).
	RecognizeAction   string `json:"recognizeAction"`
	RecognizeTitle    string `json:"recognizeTitle"`
	RecognizeSubtitle string `json:"recognizeSubtitle"`

	// 2026-05-06 revenue-run plan Phase 1 — CTA labels for the three billing-kind
	// branches on the Invoices tab. resolveInvoicesPrimaryAction (Phase 6) picks
	// the correct one; all three must be pre-populated so no branch returns "".
	RunInvoicesAction   string `json:"runInvoicesAction"`
	PoolRecognizeAction string `json:"poolRecognizeAction"`
	RequestUsageAction  string `json:"requestUsageAction"`

	// 2026-05-11 run-invoices-polish Phase 3 — per-row action labels surfaced
	// on the Invoices tab table (view, send email, print, edit).
	RowActions SubscriptionInvoicesRowActionsLabels `json:"rowActions"`
}

// SubscriptionInvoicesRowActionsLabels holds per-row action button labels for
// the Invoices tab on a subscription detail page.
type SubscriptionInvoicesRowActionsLabels struct {
	View      string `json:"view"`
	SendEmail string `json:"sendEmail"`
	Print     string `json:"print"`
	Edit      string `json:"edit"`
}

// SubscriptionRecognizeLabels holds drawer-form labels for the
// "Recognize Revenue" flow. See plan §5 Phase E for the full table; the
// blocking-error keys (currencyMismatchError, idempotencyError) are renamed
// from their advisory counterparts since v1 surfaces them as hard blocks.
type SubscriptionRecognizeLabels struct {
	// Header / context section
	ContextSection string `json:"contextSection"`
	ClientLabel    string `json:"clientLabel"`
	PlanLabel      string `json:"planLabel"`
	QuantityLabel  string `json:"quantityLabel"`

	// Period section
	PeriodSection string `json:"periodSection"`
	PeriodStart   string `json:"periodStart"`
	PeriodEnd     string `json:"periodEnd"`
	RevenueDate   string `json:"revenueDate"`

	// Line items table
	LineItemsSection    string `json:"lineItemsSection"`
	ColumnDescription   string `json:"columnDescription"`
	ColumnUnitPrice     string `json:"columnUnitPrice"`
	ColumnQuantity      string `json:"columnQuantity"`
	ColumnLineTotal     string `json:"columnLineTotal"`
	ColumnTreatment     string `json:"columnTreatment"`
	TotalLabel          string `json:"totalLabel"`
	RemoveLine          string `json:"removeLine"`
	TreatmentRecurring  string `json:"treatmentRecurring"`
	TreatmentFirstCycle string `json:"treatmentFirstCycle"`
	TreatmentUsageBased string `json:"treatmentUsageBased"`
	TreatmentOneTime    string `json:"treatmentOneTime"`

	// Notes
	NotesLabel       string `json:"notesLabel"`
	NotesPlaceholder string `json:"notesPlaceholder"`

	// Footer buttons (v1 — single Generate button; "Save as Draft" is dropped
	// per plan Phase D refinement since both paths run the idempotency check.)
	Generate string `json:"generate"`
	Cancel   string `json:"cancel"`

	// Blocking error banners
	CurrencyMismatchError     string `json:"currencyMismatchError"`
	IdempotencyError          string `json:"idempotencyError"`
	IdempotencyExistingLink   string `json:"idempotencyExistingLink"`
	NoLinesError              string `json:"noLinesError"`
	CycleNotConfiguredWarning string `json:"cycleNotConfiguredWarning"`
	UsageBasedSkippedNotice   string `json:"usageBasedSkippedNotice"`

	// 2026-04-27 plan-client-scope plan §7 — info notice on the recognize
	// drawer when the active subscription's PricePlan is client-scoped.
	// Templated via {{.ClientName}}.
	ClientCustomNotice string `json:"clientCustomNotice"`

	// 2026-04-29 milestone-billing plan §5 / Phase E — milestone-specific
	// drawer fields. Surfaced only when pricePlan.billing_kind = MILESTONE.
	MilestoneSelect            string `json:"milestoneSelect"`
	MilestoneSelectPlaceholder string `json:"milestoneSelectPlaceholder"`
	NoReadyMilestone           string `json:"noReadyMilestone"`
	MilestoneNotApplicable     string `json:"milestoneNotApplicable"`
	BillAmount                 string `json:"billAmount"`
	LeaveRemainderOpen         string `json:"leaveRemainderOpen"`
	CloseShort                 string `json:"closeShort"`
	PartialReason              string `json:"partialReason"`
	PartialReasonRequired      string `json:"partialReasonRequired"`
	OverBillingRejected        string `json:"overBillingRejected"`

	// Tax preview labels (Phase 5)
	TaxPreviewSection       string `json:"taxPreviewSection"`
	TaxDirectionSurcharge   string `json:"taxDirectionSurcharge"`
	TaxDirectionWithholding string `json:"taxDirectionWithholding"`
	NetReceivable           string `json:"netReceivable"`
	WHTAmount               string `json:"whtAmount"`
	// TaxKindLabels maps tax_kind_snapshot values to localized display names.
	// Populated from lyngua; used by convertPreviewTaxLines in the recognize view.
	TaxKindLabels map[string]string `json:"taxKindLabels"`
}

// SubscriptionRevenueRunLabels holds drawer-form labels for the per-subscription
// "Invoice Run" drawer (Surface C — CYCLE billing_kind only). Lyngua key:
// `subscription.revenueRun`. Drops the engagement column that the client-level
// drawer (Surface A) carries — per-sub context makes it redundant.
type SubscriptionRevenueRunLabels struct {
	Title string `json:"title"`
	// Subtitle is templated with the subscription name — e.g. "Run invoices for {{.Name}}"
	Subtitle string `json:"subtitle"`

	// Read-only context row labels (subscription name + plan name)
	SubscriptionLabel string `json:"subscriptionLabel"`
	PlanLabel         string `json:"planLabel"`
	// ClientHintTemplate is shown as a hint beneath the subscription field.
	// Use {client} as the substitution token. E.g. "Client: {client}".
	ClientHintTemplate string `json:"clientHintTemplate"`

	AsOfDateLabel         string `json:"asOfDateLabel"`
	AsOfDateHint          string `json:"asOfDateHint"`
	BillThroughTodayLabel string `json:"billThroughTodayLabel"`

	// Period table columns
	ColumnPeriod string `json:"columnPeriod"`
	ColumnAmount string `json:"columnAmount"`
	ColumnLines  string `json:"columnLines"`

	// Group headings / empty states
	GroupNoPending        string `json:"groupNoPending"`
	GroupCurrencyMismatch string `json:"groupCurrencyMismatch"`
	EmptyTitle            string `json:"emptyTitle"`
	EmptyMessage          string `json:"emptyMessage"`

	// IntroMessage is shown at the top of the drawer body as an info alert.
	IntroMessage string `json:"introMessage"`

	// Footer buttons
	GenerateButton          string `json:"generateButton"`
	GenerateButtonCountOne  string `json:"generateButtonCountOne"`
	GenerateButtonCountMany string `json:"generateButtonCountMany"`
	CancelButton            string `json:"cancelButton"`

	// Post-submit feedback
	ToastSuccess string `json:"toastSuccess"`
	ToastSkipped string `json:"toastSkipped"`
	ToastErrored string `json:"toastErrored"`
	ViewRunLink  string `json:"viewRunLink"`

	// Inline error messages
	Errors SubscriptionRevenueRunErrorLabels `json:"errors"`
}

// SubscriptionRevenueRunErrorLabels holds inline error strings for the
// per-subscription Invoice Run drawer.
type SubscriptionRevenueRunErrorLabels struct {
	PermissionDenied   string `json:"permissionDenied"`
	IDRequired         string `json:"idRequired"`
	InvalidFormData    string `json:"invalidFormData"`
	UseCaseUnavailable string `json:"useCaseUnavailable"`
	SelectOne          string `json:"selectOne"`
}

// SubscriptionMilestoneLabels holds labels for the Subscription Package tab's
// Milestones section + the mark-ready / waive CTAs. Lyngua key:
// `subscription.milestone.*`. See milestone-billing plan §5.
type SubscriptionMilestoneLabels struct {
	Title           string `json:"title"`
	Subtitle        string `json:"subtitle"`
	MarkReady       string `json:"markReady"`
	Waive           string `json:"waive"`
	ViewInvoice     string `json:"viewInvoice"`
	StatusPending   string `json:"statusPending"`
	StatusReady     string `json:"statusReady"`
	StatusBilled    string `json:"statusBilled"`
	StatusWaived    string `json:"statusWaived"`
	StatusDeferred  string `json:"statusDeferred"`
	StatusCancelled string `json:"statusCancelled"`
	TotalInvoiced   string `json:"totalInvoiced"`
	AmountFull      string `json:"amountFull"`
	AmountPartial   string `json:"amountPartial"`

	// 20260517-advance-cash-events Plan B Phase 7 — Recognize CTA + the
	// "linked to advance" badge that flags milestones tied to an advance
	// Collection (via the collection_billing_event junction).
	Recognize          string `json:"recognize"`
	LinkedAdvanceBadge string `json:"linkedAdvanceBadge"`
}

// SubscriptionOperationsLabels holds labels for the Subscription detail's
// Operations tab. Lyngua key: `subscription.detail.operations.*`. See
// auto-spawn-jobs-from-subscription plan §5.2 / §9 and cyclic-subscription-jobs
// plan §9.1 (cycle accordion + backfill keys).
type SubscriptionOperationsLabels struct {
	Title        string `json:"title"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
	SpawnAction  string `json:"spawnAction"`
	RootJob      string `json:"rootJob"`
	ChildJob     string `json:"childJob"`
	PhaseSummary string `json:"phaseSummary"`
	ViewJobLink  string `json:"viewJobLink"`

	// 2026-04-30 cyclic-subscription-jobs plan §9.1 — cycle accordion copy.
	SubscriptionHeading   string `json:"subscriptionHeading"`
	CycleHeading          string `json:"cycleHeading"`
	CyclePlaceholder      string `json:"cyclePlaceholder"`
	CycleSpawnNow         string `json:"cycleSpawnNow"`
	CycleStatusPending    string `json:"cycleStatusPending"`
	CycleStatusInProgress string `json:"cycleStatusInProgress"`
	CycleStatusCompleted  string `json:"cycleStatusCompleted"`
	CycleStatusOverdue    string `json:"cycleStatusOverdue"`
	CycleInvoiceLinked    string `json:"cycleInvoiceLinked"`
	CycleNoInvoice        string `json:"cycleNoInvoice"`
	CycleEmpty            string `json:"cycleEmpty"`
	BackfillBanner        string `json:"backfillBanner"`
	BackfillCta           string `json:"backfillCta"`

	// 2026-05-01 ad-hoc-subscription-billing plan §5.2 — Operations tab
	// AD_HOC mode keys. Vertical-neutral defaults ("usage", "occurrence")
	// with professional-tier overrides ("service call", "retainer", etc.).
	AdHocPoolHeading       string `json:"adHocPoolHeading"`
	AdHocPerCallHeading    string `json:"adHocPerCallHeading"`
	EntitlementUsed        string `json:"entitlementUsed"`
	EntitlementRemaining   string `json:"entitlementRemaining"`
	EntitlementExhausted   string `json:"entitlementExhausted"`
	RequestUsageCta        string `json:"requestUsageCta"`
	ExtendEntitlementCta   string `json:"extendEntitlementCta"`
	UsageRequestedDate     string `json:"usageRequestedDate"`
	UsageDeliveredDate     string `json:"usageDeliveredDate"`
	UsageOrdinalLabel      string `json:"usageOrdinalLabel"`
	UsageNotDelivered      string `json:"usageNotDelivered"`
	PoolInvoiceLink        string `json:"poolInvoiceLink"`
	PoolInvoicePending     string `json:"poolInvoicePending"`
	PoolGenerateInvoiceCta string `json:"poolGenerateInvoiceCta"`
	PerCallRecognizeCta    string `json:"perCallRecognizeCta"`
	PerCallInvoiceLink     string `json:"perCallInvoiceLink"`
	PerCallNotReady        string `json:"perCallNotReady"`
}

// SubscriptionBackfillLabels holds labels for the Backfill cycle Jobs drawer.
// Lyngua key: `subscription.detail.backfill.*`. See cyclic-subscription-jobs
// plan §9.2.
type SubscriptionBackfillLabels struct {
	DrawerTitle       string `json:"drawerTitle"`
	DrawerDescription string `json:"drawerDescription"`
	PreviewLine       string `json:"previewLine"`
	CountLabel        string `json:"countLabel"`
	Confirm           string `json:"confirm"`
	Cancel            string `json:"cancel"`
	MaxWarning        string `json:"maxWarning"`
}

// SubscriptionJobsTabLabels holds labels for the new flat Jobs tab on the
// Subscription detail page. Lyngua key: `subscription.detail.jobs.*`. See
// cyclic-subscription-jobs plan §21.
type SubscriptionJobsTabLabels struct {
	Heading          string `json:"heading"`
	Empty            string `json:"empty"`
	FilterStatus     string `json:"filterStatus"`
	FilterType       string `json:"filterType"`
	FilterAll        string `json:"filterAll"`
	SortBy           string `json:"sortBy"`
	SortByCycle      string `json:"sortByCycle"`
	ExportCsv        string `json:"exportCsv"`
	Summary          string `json:"summary"`
	ColumnNumber     string `json:"columnNumber"`
	ColumnName       string `json:"columnName"`
	ColumnType       string `json:"columnType"`
	ColumnPhase      string `json:"columnPhase"`
	ColumnStatus     string `json:"columnStatus"`
	ColumnPeriod     string `json:"columnPeriod"`
	TypeSubscription string `json:"typeSubscription"`
	TypeOnboarding   string `json:"typeOnboarding"`
	TypeCycle        string `json:"typeCycle"`
	TypeVisit        string `json:"typeVisit"`
	SpawnFailedToast string `json:"spawnFailedToast"`
}

// SubscriptionSpawnLabels holds labels for the retroactive Spawn Jobs drawer.
// Lyngua key: `subscription.spawn.*`. See auto-spawn-jobs-from-subscription
// plan §5.3 / §9.
type SubscriptionSpawnLabels struct {
	Title             string `json:"title"`
	DetectedTemplates string `json:"detectedTemplates"`
	RootTemplate      string `json:"rootTemplate"`
	Cancel            string `json:"cancel"`
	Confirm           string `json:"confirm"`
	SuccessToast      string `json:"successToast"`
	Skipped           string `json:"skipped"`
}

type SubscriptionConfirmLabels struct {
	Cancel                string `json:"cancel"`
	CancelMessage         string `json:"cancelMessage"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activateMessage"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivateMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
}

// DefaultSubscriptionLabels returns SubscriptionLabels with sensible English defaults.
func DefaultSubscriptionLabels() SubscriptionLabels {
	return SubscriptionLabels{
		Page: SubscriptionPageLabels{
			Heading:         "Subscriptions",
			HeadingActive:   "Active Subscriptions",
			HeadingInactive: "Inactive Subscriptions",
			Caption:         "Subscription management",
			CaptionActive:   "Manage your active subscriptions",
			CaptionInactive: "View cancelled or expired subscriptions",
		},
		Buttons: SubscriptionButtonLabels{
			AddSubscription: "Add Subscription",
		},
		Columns: SubscriptionColumnLabels{
			Name:      "Engagement",
			Client:    "Client",
			Customer:  "Customer",
			Plan:      "Plan",
			StartDate: "Start Date",
			EndDate:   "End Date",
			Status:    "Status",
		},
		Empty: SubscriptionEmptyLabels{
			Title:   "No subscriptions found",
			Message: "No subscriptions to display.",
		},
		Form: SubscriptionFormLabels{
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
		},
		Actions: SubscriptionActionLabels{
			View:       "View Subscription",
			Edit:       "Edit Subscription",
			Cancel:     "Cancel Subscription",
			Delete:     "Delete",
			Activate:   "Activate",
			Deactivate: "Deactivate",
			// 2026-04-27 plan-client-scope plan §6.5 / §7 — Package tab CTA.
			CustomizePackage: "Customize this package for {{.ClientName}}",
		},
		Bulk: SubscriptionBulkLabels{
			Delete:     "Delete Selected",
			Activate:   "Activate Selected",
			Deactivate: "Deactivate Selected",
		},
		Status: SubscriptionStatusLabels{
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Detail: SubscriptionDetailLabels{
			PageTitle:            "Subscription Details",
			Customer:             "Customer",
			Plan:                 "Plan",
			StartDate:            "Start Date",
			EndDate:              "End Date",
			Status:               "Status",
			CreatedDate:          "Created",
			ModifiedDate:         "Last Modified",
			AuditTrailComingSoon: "Audit trail coming soon.",
			AuditTrailDesc:       "Audit trail for subscription changes is coming soon.",
		},
		Tabs: SubscriptionTabLabels{
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
		Invoices: SubscriptionInvoicesLabels{
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
			RowActions: SubscriptionInvoicesRowActionsLabels{
				View:      "View",
				SendEmail: "Send email",
				Print:     "Print",
				Edit:      "Edit",
			},
		},
		RevenueRun: SubscriptionRevenueRunLabels{
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
			Errors: SubscriptionRevenueRunErrorLabels{
				PermissionDenied:   "You do not have permission to run invoices.",
				IDRequired:         "Subscription ID is required.",
				InvalidFormData:    "Invalid form data. Please check your inputs and try again.",
				UseCaseUnavailable: "Invoice run is not available for this subscription type.",
				SelectOne:          "Select at least one period to generate.",
			},
		},
		Recognize: SubscriptionRecognizeLabels{
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
		Milestone: SubscriptionMilestoneLabels{
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
		Operations: SubscriptionOperationsLabels{
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
		Backfill: SubscriptionBackfillLabels{
			DrawerTitle:       "Backfill cycle Jobs",
			DrawerDescription: "Preview the cycles that will be spawned, then confirm to materialize them in one transaction.",
			PreviewLine:       "Cycle {{.Index}} — {{.PeriodLabel}}",
			CountLabel:        "Cycles to spawn",
			Confirm:           "Spawn {{.Count}} cycle(s)",
			Cancel:            "Cancel",
			MaxWarning:        "Backfill is capped at 24 cycles per request. Reduce the range or run multiple backfills.",
		},
		// 2026-04-30 cyclic-subscription-jobs plan §21.3 — flat Jobs tab.
		Jobs: SubscriptionJobsTabLabels{
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
		Spawn: SubscriptionSpawnLabels{
			Title:             "Spawn Operational Jobs",
			DetectedTemplates: "Detected templates",
			RootTemplate:      "Root template",
			Cancel:            "Cancel",
			Confirm:           "Spawn Jobs",
			SuccessToast:      "Spawned {{.JobCount}} Job(s).",
			Skipped:           "Nothing to spawn — no JobTemplate is linked to this Plan.",
		},
		Confirm: SubscriptionConfirmLabels{
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
		Errors: SubscriptionErrorLabels{
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
