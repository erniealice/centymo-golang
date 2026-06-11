package price_plan

// ---------------------------------------------------------------------------
// Price Plan labels
// ---------------------------------------------------------------------------

// Labels holds all labels for the standalone price plan (rate card) module.
type Labels struct {
	Page         PageLabels         `json:"page"`
	Buttons      ButtonLabels       `json:"buttons"`
	Columns      ColumnLabels2      `json:"columns"`
	Empty        EmptyLabels        `json:"empty"`
	Form         FormLabels         `json:"form"`
	Actions      ActionLabels       `json:"actions"`
	Bulk         BulkLabels         `json:"bulk"`
	Detail       DetailLabels2      `json:"detail"`
	Tabs         TabLabels2         `json:"tabs"`
	Confirm      ConfirmLabels      `json:"confirm"`
	Errors       ErrorLabels        `json:"errors"`
	ProductPrice ProductPriceLabels `json:"productPrice"`
	Messages     MessageLabels      `json:"messages"`
}

// ProductPriceLabels holds labels for product-price sub-table actions and empty state.
type ProductPriceLabels struct {
	EditTitle   string `json:"editTitle"`
	DeleteTitle string `json:"deleteTitle"`
	EmptyTitle  string `json:"emptyTitle"`
	EmptyMsg    string `json:"emptyMsg"`
}

// MessageLabels holds translatable message strings used in the price plan
// and price schedule plan views (pricing-lock notices, validation errors).
type MessageLabels struct {
	PricingLockedReason     string `json:"pricingLockedReason"`
	ItemPricingLockedReason string `json:"itemPricingLockedReason"`
	CreateNotAvailable      string `json:"createNotAvailable"`
	UpdateNotAvailable      string `json:"updateNotAvailable"`
	ProductRequired         string `json:"productRequired"`
	InvalidPrice            string `json:"invalidPrice"`
	InUseCannotModify       string `json:"inUseCannotModify"`
	IDRequired              string `json:"idRequired"`
	DeleteNotAvailable      string `json:"deleteNotAvailable"`
	CurrencyMismatch        string `json:"currencyMismatch"`
}

type PageLabels struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	ActiveTitle   string `json:"activeTitle"`
	InactiveTitle string `json:"inactiveTitle"`
}

type ButtonLabels struct {
	View       string `json:"view"`
	Add        string `json:"add"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	BulkDelete string `json:"bulkDelete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ColumnLabels2 struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Duration    string `json:"duration"`
	Location    string `json:"location"`
	Schedule    string `json:"schedule"`
	Plan        string `json:"plan"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
	Actions     string `json:"actions"`
}

type EmptyLabels struct {
	Title       string `json:"title"`
	Message     string `json:"message"`
	Description string `json:"description"`
	ActionLabel string `json:"actionLabel"`
}

type ActionLabels struct {
	CreateSuccess string `json:"createSuccess"`
	CreateError   string `json:"createError"`
	UpdateSuccess string `json:"updateSuccess"`
	UpdateError   string `json:"updateError"`
	DeleteSuccess string `json:"deleteSuccess"`
	DeleteError   string `json:"deleteError"`
}

type BulkLabels struct {
	DeleteTitle   string `json:"deleteTitle"`
	DeleteMessage string `json:"deleteMessage"`
	StatusTitle   string `json:"statusTitle"`
	StatusMessage string `json:"statusMessage"`
}

type DetailLabels2 struct {
	Title          string `json:"title"`
	InfoTab        string `json:"infoTab"`
	AttachmentsTab string `json:"attachmentsTab"`
	AuditTab       string `json:"auditTab"`
	ProductsTab    string `json:"productsTab"`

	// Info-tab field labels (price-schedule-plan-tab-info).
	Heading       string `json:"heading"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Duration      string `json:"duration"`
	ScheduleLabel string `json:"scheduleLabel"`
	Status        string `json:"status"`
	DateCreated   string `json:"dateCreated"`
	DateModified  string `json:"dateModified"`
	Edit          string `json:"edit"`
	EditTitle     string `json:"editTitle"`

	// 2026-04-30 cyclic-subscription-jobs plan §20 — Billing model summary
	// rendered on the info tab. Lyngua key: `pricePlan.detail.summary*`.
	SummaryHeading            string             `json:"summaryHeading"`
	CustomerHeading           string             `json:"customerHeading"`
	OperationsHeading         string             `json:"operationsHeading"`
	RevenueRecognitionHeading string             `json:"revenueRecognitionHeading"`
	Summary                   BillingSummaryCopy `json:"summary"`
	Warning                   BillingSummaryWarn `json:"warning"`

	// 2026-05-04 — Subscriptions/Engagements tab on the price-plan detail.
	// See docs/plan/20260504-price-plan-engagements-tab/.
	Subscriptions SubscriptionsSectionLabels `json:"subscriptions"`
}

// BillingSummaryCopy carries the per-(kind × basis) lyngua copy that
// `buildBillingModelSummary` projects into the info-sections grid.
// Cyclic plan ships rows 1-6 (oneTime / recurring / contract / milestone);
// AD_HOC plan adds adHoc.* in a follow-up. Each entry has 3 lines:
// customer, operations, revenue.
type BillingSummaryCopy struct {
	OneTime   SummaryByBasis `json:"oneTime"`
	Recurring SummaryByBasis `json:"recurring"`
	Contract  SummaryByBasis `json:"contract"`
	Milestone SummaryByBasis `json:"milestone"`
	AdHoc     SummaryByBasis `json:"adHoc"`
}

// SummaryByBasis groups the text lines per basis. Empty
// strings on a basis means "no copy for that combo" — view skips it.
type SummaryByBasis struct {
	PerCycle         SummaryLines `json:"perCycle"`
	TotalPackage     SummaryLines `json:"totalPackage"`
	DerivedFromLines SummaryLines `json:"derivedFromLines"`
	PerOccurrence    SummaryLines `json:"perOccurrence"`
}

// SummaryLines holds the 3 lines for a kind × basis cell.
type SummaryLines struct {
	Customer   string `json:"customer"`
	Operations string `json:"operations"`
	Revenue    string `json:"revenue"`
}

// BillingSummaryWarn carries the warning-row copy keyed by symbol
// per plan §20.3. View only renders entries whose preconditions trip.
type BillingSummaryWarn struct {
	MilestoneNoTemplate           string `json:"milestoneNoTemplate"`
	RecurringNoTemplate           string `json:"recurringNoTemplate"`
	VisitsPerCycleInvalidKind     string `json:"visitsPerCycleInvalidKind"`
	AdHocPoolNoTemplate           string `json:"adHocPoolNoTemplate"`
	AdHocPerCallNoTemplate        string `json:"adHocPerCallNoTemplate"`
	AdHocNoEntitlement            string `json:"adHocNoEntitlement"`
	AdHocBillingCycleNotAllowed   string `json:"adHocBillingCycleNotAllowed"`
	AdHocVisitsPerCycleNotAllowed string `json:"adHocVisitsPerCycleNotAllowed"`
}

type TabLabels2 struct {
	Info          string `json:"info"`
	Products      string `json:"products"`
	Subscriptions string `json:"subscriptions"`
	Attachments   string `json:"attachments"`
	Audit         string `json:"audit"`
}

// SubscriptionsSectionLabels holds the column headers, empty state,
// and confirm-delete copy for the price-plan detail "Subscriptions" tab —
// professional tier overrides this block to use the engagement vocabulary.
type SubscriptionsSectionLabels struct {
	ColumnName           string `json:"columnName"`
	ColumnClient         string `json:"columnClient"`
	ColumnPlan           string `json:"columnPlan"`
	ColumnStartDate      string `json:"columnStartDate"`
	ColumnEndDate        string `json:"columnEndDate"`
	EmptyTitle           string `json:"emptyTitle"`
	EmptyMessage         string `json:"emptyMessage"`
	ConfirmDeleteTitle   string `json:"confirmDeleteTitle"`
	ConfirmDeleteMessage string `json:"confirmDeleteMessage"`
}

type ConfirmLabels struct {
	DeleteTitle       string `json:"deleteTitle"`
	DeleteMessage     string `json:"deleteMessage"`
	DeactivateTitle   string `json:"deactivateTitle"`
	DeactivateMessage string `json:"deactivateMessage"`

	// 2026-04-27 plan-client-scope plan §3.5 — fired by the centymo confirm
	// dialog when an operator changes monetary fields on a client-scoped
	// PricePlan that has N > 1 active subscriptions. Templated via
	// {{.Count}} and {{.ClientName}}.
	EditAmountMultipleSubscriptions string `json:"editAmountMultipleSubscriptions"`
}

type ErrorLabels struct {
	NotFound     string `json:"notFound"`
	LoadFailed   string `json:"loadFailed"`
	Unauthorized string `json:"unauthorized"`
	CreateFailed string `json:"createFailed"`
	UpdateFailed string `json:"updateFailed"`
	DeleteFailed string `json:"deleteFailed"`
	InUse        string `json:"inUse"`

	// 2026-04-27 plan-client-scope plan §7. Surfaced when an UpdatePricePlan
	// body sends a client_id that doesn't match the parent Plan's client_id.
	ClientScopeMismatch string `json:"clientScopeMismatch"`
	// 2026-04-28 — surfaced when the operator picks a price_schedule whose
	// client_id belongs to a different client than the parent Plan. Master
	// schedules (sched.client_id == "") are still accepted; only the
	// cross-client cases get rejected.
	ScheduleClientMismatch string `json:"scheduleClientMismatch"`
	// 2026-04-28 — surfaced when an operator submits a PricePlan with no
	// price_schedule_id under a client-scoped Plan. The use case used to
	// auto-create a schedule with a synthetic now() date; reverted because
	// that hid real operator intent. Operator must pick or create a client
	// rate card first.
	ScheduleRequiredForClientScope string `json:"scheduleRequiredForClientScope"`
	// Server-side-only error key — the centymo confirm dialog catches the
	// N>1-engagements gate before this surfaces.
	MultiSubscriptionConfirmRequired string `json:"multiSubscriptionConfirmRequired"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:         "Rate Cards",
			Subtitle:      "Manage your rate cards",
			ActiveTitle:   "Active Rate Cards",
			InactiveTitle: "Inactive Rate Cards",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Rate Card",
			Edit:       "Edit Rate Card",
			Delete:     "Delete Rate Card",
			BulkDelete: "Delete Rate Cards",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: ColumnLabels2{
			Name:        "Name",
			Description: "Description",
			Amount:      "Amount",
			Currency:    "Currency",
			Duration:    "Duration",
			Location:    "Location",
			Schedule:    "Schedule",
			Plan:        "Plan",
			Status:      "Status",
			DateCreated: "Date Created",
			Actions:     "Actions",
		},
		Empty: EmptyLabels{
			Title:       "No Rate Cards",
			Message:     "No rate cards to display.",
			Description: "Add a rate card to define pricing for your plans.",
			ActionLabel: "Add Rate Card",
		},
		Form: FormLabels{
			Name:                "Price Plan Name",
			NamePlaceholder:     "Enter price plan name",
			Description:         "Description",
			DescPlaceholder:     "Enter description...",
			Amount:              "Amount",
			AmountPlaceholder:   "0.00",
			Currency:            "Currency",
			CurrencyPlaceholder: "e.g. PHP",
			DurationValue:       "Duration",
			DurationUnit:        "Unit",
			Schedule:            "Price Schedule",
			SchedulePlaceholder: "Select a schedule...",
			ScheduleSearch:      "Filter...",
			Location:            "Location",
			LocationPlaceholder: "Select a location...",
			LocationHintPrefix:  "Location: ",
			SelectLocation:      "— No location (all locations) —",
			Active:              "Active",
			PlanLabel:           "Package",
			PlanPlaceholder:     "Select a package...",
			PlanSearch:          "Filter...",
			// Wave 2 new fields
			SectionBasic:                "Basic info",
			SectionPricing:              "Pricing",
			BillingKindLabel:            "Billing model",
			BillingKindOneTime:          "One-time",
			BillingKindRecurring:        "Recurring retainer",
			BillingKindContract:         "Fixed-term engagement",
			BillingKindMilestone:        "Milestone",
			AmountBasisLabel:            "Amount basis",
			AmountBasisPerCycle:         "Per cycle",
			AmountBasisTotalPackage:     "Total package",
			AmountBasisDerivedFromLines: "Sum of items",
			BillingCycleLabel:           "Billing cycle",
			BillingCyclePlaceholder:     "e.g. every 1 month",
			TermLabel:                   "Term",
			TermPlaceholder:             "e.g. 12 months",
			TermOpenEndedHelp:           "Leave empty for open-ended / no expiration",
			// Field-level info popovers — use proto-generic wording; business-type
			// tiers override via lyngua (e.g. "plan" → "package" / "rate card").
			PlanInfo:         "The plan this price plan belongs to. Locked from the parent page.",
			ScheduleInfo:     "The price schedule (date range + location) this price plan belongs to.",
			NameInfo:         "Optional — defaults to the plan name when left blank.",
			DescriptionInfo:  "Optional notes shown alongside the price plan in detail views.",
			BillingKindInfo:  "One-time = charged once. Recurring = billed every cycle. Fixed-term = recurring with an end date.",
			AmountBasisInfo:  "Per cycle = amount charged each billing cycle. Total package = amount charged across the full term. Sum of items = derived from the per-item breakdown.",
			AmountInfo:       "Price in the selected currency. For Sum of items, this is computed automatically.",
			CurrencyInfo:     "Currency applied to this price plan and any auto-seeded product price plans.",
			BillingCycleInfo: "How often the recurring charge is issued (e.g. every 1 month).",
			TermInfo:         "How long the engagement lasts. Leave empty for open-ended / no expiration.",
			ActiveInfo:       "Inactive price plans stay on record but are hidden from new subscriptions.",
			// 2026-04-27 plan-client-scope plan §6.7 — info banner shown above
			// the form when its parent PriceSchedule is client-scoped. Loaded
			// from lyngua price_plan.json#price_plan.form.parentScheduleClientNotice
			// (under "form", not "fields" — the Go struct lives on FormLabels).
			ParentScheduleClientNotice: "This price schedule belongs to {{.ClientName}}. Price plans created here will be available only for engagements with this client.",
			// 2026-04-27 plan-client-scope plan §6.7 — tooltip on the readonly
			// Schedule label when the parent Plan is client-scoped. Proto-generic
			// wording; tier overrides live in lyngua.
			ScheduleLockedTooltip:  "This price plan is bound to {{.ClientName}}'s price schedule.",
			ScheduleAutoCreateHint: "No price schedule exists for {{.ClientName}} yet — one will be created automatically when you save.",
			ScheduleAutoReuseHint:  "This price plan will be added to the existing price schedule for {{.ClientName}}.",
			// 2026-05-03 — Plan picker scope notices on the schedule-scoped add drawer.
			ScheduleClientPickerNotice:  "Only packages assigned to this client are available below.",
			ScheduleGeneralPickerNotice: "Packages assigned specifically to clients are not available for assignment to general scoped rate cards.",
			// 2026-04-30 cyclic-subscription-jobs plan §9.4.
			MilestoneCyclicBlock: "Milestone billing is not supported on cyclic plans (RECURRING / CONTRACT × PER_CYCLE / multi-visit).",
		},
		Actions: ActionLabels{
			CreateSuccess: "Rate card created successfully.",
			CreateError:   "Failed to create rate card.",
			UpdateSuccess: "Rate card updated successfully.",
			UpdateError:   "Failed to update rate card.",
			DeleteSuccess: "Rate card deleted successfully.",
			DeleteError:   "Failed to delete rate card.",
		},
		Bulk: BulkLabels{
			DeleteTitle:   "Delete Rate Cards",
			DeleteMessage: "Are you sure you want to delete the selected rate cards?",
			StatusTitle:   "Update Status",
			StatusMessage: "Are you sure you want to update the status of the selected rate cards?",
		},
		Detail: DetailLabels2{
			Title:          "Rate Card Details",
			InfoTab:        "Information",
			AttachmentsTab: "Attachments",
			AuditTab:       "Audit Trail",
			ProductsTab:    "Products",
			Heading:        "Plan Info",
			Name:           "Name",
			Description:    "Description",
			Amount:         "Amount",
			Currency:       "Currency",
			Duration:       "Duration",
			ScheduleLabel:  "Schedule",
			Status:         "Status",
			DateCreated:    "Date Created",
			DateModified:   "Date Modified",
			Edit:           "Edit Price Plan",
			EditTitle:      "Edit Price Plan",
			// 2026-04-30 cyclic-subscription-jobs plan §20.
			SummaryHeading:            "Billing model summary",
			CustomerHeading:           "Customer experience",
			OperationsHeading:         "Operations impact",
			RevenueRecognitionHeading: "Revenue recognition",
			Summary: BillingSummaryCopy{
				OneTime: SummaryByBasis{
					TotalPackage: SummaryLines{
						Customer:   "Pays {{.Amount}} once at signup. No further charges.",
						Operations: "Engagement spawns 1 lifetime Job with phases (if Plan has a JobTemplate).",
						Revenue:    "One Revenue at Subscription.Create covering the full amount.",
					},
				},
				Recurring: SummaryByBasis{
					PerCycle: SummaryLines{
						Customer:   "Charged {{.Amount}} every {{.CycleLabel}}. Subscription auto-renews until cancelled.",
						Operations: "Each cycle spawns {{.VisitsPerCycle}} cycle Job(s) (if Plan has a JobTemplate). Operations tab shows cycle accordions.",
						Revenue:    "One Revenue per cycle. Recognize Revenue creates the invoice and (via piggyback) spawns the cycle Job if missing.",
					},
					DerivedFromLines: SummaryLines{
						Customer:   "Charged the sum of itemised lines every {{.CycleLabel}}.",
						Operations: "Each cycle spawns 1+ cycle Jobs. Operations tracking flows through Plan's JobTemplate.",
						Revenue:    "Revenue total computed from ProductPricePlan rows; one Revenue per cycle.",
					},
				},
				Contract: SummaryByBasis{
					PerCycle: SummaryLines{
						Customer:   "Charged {{.Amount}} every {{.CycleLabel}} for {{.TermLength}}. Auto-deactivates at term end.",
						Operations: "Same as recurring + the engagement closes when the {{.TermLength}} term completes.",
						Revenue:    "Same as recurring. Operator can extend the term to spawn additional cycles.",
					},
					TotalPackage: SummaryLines{
						Customer:   "Pays {{.Amount}} upfront for {{.TermLength}} of service.",
						Operations: "Engagement spawns 1 lifetime Job (or cycle Jobs if cyclic — see Plan's visits_per_cycle).",
						Revenue:    "One Revenue at signup; cycle Jobs are operational only.",
					},
				},
				Milestone: SummaryByBasis{
					TotalPackage: SummaryLines{
						Customer:   "Pays {{.Amount}} total. Invoice fires per milestone (engagement phase) as work completes.",
						Operations: "Lifetime engagement Job with phases. BillingEvent rows gate per-milestone invoicing.",
						Revenue:    "Revenue per milestone trigger; sum across milestones equals the total package.",
					},
				},
			},
			Warning: BillingSummaryWarn{
				MilestoneNoTemplate:       "Milestone billing requires the Plan to have a JobTemplate. Configure it on the Plan first.",
				RecurringNoTemplate:       "This subscription will not have operational tracking. Add a JobTemplate to the Plan to enable cycle Jobs.",
				VisitsPerCycleInvalidKind: "visits_per_cycle is only valid for cyclic plans. Reset to 1 or change the billing kind.",
			},
			Subscriptions: SubscriptionsSectionLabels{
				ColumnName:           "Subscription",
				ColumnClient:         "Client",
				ColumnPlan:           "Plan",
				ColumnStartDate:      "Start Date",
				ColumnEndDate:        "End Date",
				EmptyTitle:           "No subscriptions yet",
				EmptyMessage:         "Subscriptions referencing this price plan will appear here.",
				ConfirmDeleteTitle:   "Delete Subscription",
				ConfirmDeleteMessage: "Are you sure you want to delete subscription %s? This action cannot be undone.",
			},
		},
		Tabs: TabLabels2{
			Info:          "Information",
			Products:      "Products",
			Subscriptions: "Subscriptions",
			Attachments:   "Attachments",
			Audit:         "Audit Trail",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:       "Delete Rate Card",
			DeleteMessage:     "Are you sure you want to delete this rate card? This action cannot be undone.",
			DeactivateTitle:   "Deactivate Rate Card",
			DeactivateMessage: "Are you sure you want to deactivate this rate card?",
			// 2026-04-27 plan-client-scope plan §3.5 / §7.
			EditAmountMultipleSubscriptions: "This price plan is attached to {{.Count}} active subscriptions for {{.ClientName}}. Changing the amount or cycle will affect all of them on the next bill cycle. Continue?",
		},
		Errors: ErrorLabels{
			NotFound:                         "Rate card not found.",
			LoadFailed:                       "Failed to load rate cards.",
			Unauthorized:                     "You do not have permission to access this resource.",
			CreateFailed:                     "Failed to create rate card.",
			UpdateFailed:                     "Failed to update rate card.",
			DeleteFailed:                     "Failed to delete rate card.",
			InUse:                            "This price plan is in use by active subscriptions and cannot be deleted.",
			ClientScopeMismatch:              "Price plan client must match its parent plan's client.",
			ScheduleClientMismatch:           "Selected schedule belongs to a different client and cannot be attached to this price plan.",
			ScheduleRequiredForClientScope:   "This package is scoped to a client. Pick or create a rate card for that client before adding a price plan.",
			MultiSubscriptionConfirmRequired: "Confirmation required — multiple attached subscriptions and monetary fields changing.",
		},
		ProductPrice: ProductPriceLabels{
			EditTitle:   "Edit Product Price",
			DeleteTitle: "Delete Product Price",
			EmptyTitle:  "No Product Prices",
			EmptyMsg:    "No product prices have been configured for this rate card yet.",
		},
		Messages: MessageLabels{
			PricingLockedReason:     "This plan is in use by active subscriptions. Pricing changes are disabled. You can still rename or reassign the package.",
			ItemPricingLockedReason: "This package is in use by active engagements. Item price and currency are locked to keep billing consistent.",
			CreateNotAvailable:      "Product price plan create is not available.",
			UpdateNotAvailable:      "Product price plan update is not available.",
			ProductRequired:         "Product is required.",
			InvalidPrice:            "Invalid price value.",
			InUseCannotModify:       "This package is in use by active engagements. Item price and currency are locked.",
			IDRequired:              "ID is required.",
			DeleteNotAvailable:      "Product price plan delete is not available.",
			CurrencyMismatch:        "Currency must match the rate card currency.",
		},
	}
}

// ---------------------------------------------------------------------------
// PricePlan add/edit form labels (relocated from plan_labels.go god-file)
// ---------------------------------------------------------------------------

// FormLabels holds translatable labels for the PricePlan add/edit form.
type FormLabels struct {
	Name                string `json:"name"`
	NamePlaceholder     string `json:"namePlaceholder"`
	Description         string `json:"description"`
	DescPlaceholder     string `json:"descriptionPlaceholder"`
	Amount              string `json:"amount"`
	AmountPlaceholder   string `json:"amountPlaceholder"`
	Currency            string `json:"currency"`
	CurrencyPlaceholder string `json:"currencyPlaceholder"`
	DurationValue       string `json:"durationValue"`
	DurationUnit        string `json:"durationUnit"`
	Schedule            string `json:"schedule"`
	SchedulePlaceholder string `json:"schedulePlaceholder"`
	ScheduleSearch      string `json:"scheduleSearch"`
	Location            string `json:"location"`
	LocationPlaceholder string `json:"locationPlaceholder"`
	LocationHintPrefix  string `json:"locationHintPrefix"`
	SelectLocation      string `json:"selectLocation"`
	Active              string `json:"active"`
	PlanLabel           string `json:"planLabel"`
	PlanPlaceholder     string `json:"planPlaceholder"`
	PlanSearch          string `json:"planSearch"`

	// Wave 2 — new billing semantics fields (from lyngua price_plan.json → price_plan.form)
	SectionBasic         string `json:"sectionBasic"`
	SectionPricing       string `json:"sectionPricing"`
	BillingKindLabel     string `json:"billingKindLabel"`
	BillingKindOneTime   string `json:"billingKindOneTime"`
	BillingKindRecurring string `json:"billingKindRecurring"`
	BillingKindContract  string `json:"billingKindContract"`
	BillingKindMilestone string `json:"billingKindMilestone"`
	BillingKindAdHoc     string `json:"billingKindAdHoc"`
	// Per-option hint copy surfaced inline below the billing_kind select as the
	// operator picks. Matches the multi-vertical convention — general/ tier ships
	// neutral phrasing, professional/ overrides with engagement vocabulary.
	BillingKindOneTimeHint      string `json:"billingKindOneTimeHint"`
	BillingKindRecurringHint    string `json:"billingKindRecurringHint"`
	BillingKindContractHint     string `json:"billingKindContractHint"`
	BillingKindMilestoneHint    string `json:"billingKindMilestoneHint"`
	BillingKindAdHocHint        string `json:"billingKindAdHocHint"`
	AmountBasisLabel            string `json:"amountBasisLabel"`
	AmountBasisPerCycle         string `json:"amountBasisPerCycle"`
	AmountBasisTotalPackage     string `json:"amountBasisTotalPackage"`
	AmountBasisDerivedFromLines string `json:"amountBasisDerivedFromLines"`
	AmountBasisPerOccurrence    string `json:"amountBasisPerOccurrence"`
	// Per-option hint copy for amount_basis (mirrors billing_kind pattern).
	AmountBasisPerCycleHint         string `json:"amountBasisPerCycleHint"`
	AmountBasisTotalPackageHint     string `json:"amountBasisTotalPackageHint"`
	AmountBasisDerivedFromLinesHint string `json:"amountBasisDerivedFromLinesHint"`
	AmountBasisPerOccurrenceHint    string `json:"amountBasisPerOccurrenceHint"`
	EntitledOccurrencesLabel        string `json:"entitledOccurrencesLabel"`
	EntitledOccurrencesPlaceholder  string `json:"entitledOccurrencesPlaceholder"`
	EntitledOccurrencesInfo         string `json:"entitledOccurrencesInfo"`
	BillingCycleLabel               string `json:"billingCycleLabel"`
	BillingCyclePlaceholder         string `json:"billingCyclePlaceholder"`
	TermLabel                       string `json:"termLabel"`
	TermPlaceholder                 string `json:"termPlaceholder"`
	TermOpenEndedHelp               string `json:"termOpenEndedHelp"`

	// Field-level info text surfaced via an info button beside each label.
	PlanInfo         string `json:"planInfo"`
	ScheduleInfo     string `json:"scheduleInfo"`
	NameInfo         string `json:"nameInfo"`
	DescriptionInfo  string `json:"descriptionInfo"`
	BillingKindInfo  string `json:"billingKindInfo"`
	AmountBasisInfo  string `json:"amountBasisInfo"`
	AmountInfo       string `json:"amountInfo"`
	CurrencyInfo     string `json:"currencyInfo"`
	BillingCycleInfo string `json:"billingCycleInfo"`
	TermInfo         string `json:"termInfo"`
	ActiveInfo       string `json:"activeInfo"`

	// 2026-04-27 plan-client-scope plan §6.7 — info banner shown above the
	// PricePlan add/edit form when its parent PriceSchedule is client-scoped.
	// Templated via Go's text/template ({{.ClientName}}).
	ParentScheduleClientNotice string `json:"parentScheduleClientNotice"`

	// 2026-04-27 plan-client-scope plan §6.7 — tooltip surfaced beside the
	// readonly Schedule label when the PricePlan's parent Plan is
	// client-scoped (the schedule field is locked to the resolved/derived
	// client schedule). Templated via Go's text/template ({{.ClientName}}).
	ScheduleLockedTooltip string `json:"scheduleLockedTooltip"`
	// 2026-04-28 — info-row hints rendered beneath the readonly Schedule
	// label so the operator knows what happens on save:
	//   ScheduleAutoCreateHint — no client rate card exists yet; one will be
	//     created with this client's name + the lyngua suffix.
	//   ScheduleAutoReuseHint  — an existing client rate card was found; the
	//     new price plan will attach to it.
	// Both templated with {{.ClientName}}.
	ScheduleAutoCreateHint string `json:"scheduleAutoCreateHint"`
	ScheduleAutoReuseHint  string `json:"scheduleAutoReuseHint"`

	// 2026-05-03 — info banner rendered below the readonly Schedule display.
	// ScheduleClientPickerNotice fires when the parent schedule is
	// client-scoped: picker shows only plans assigned to that client.
	// ScheduleGeneralPickerNotice fires when the parent schedule is
	// general-scope: picker shows only general-scope plans (client-specific
	// plans cannot attach to a general schedule).
	ScheduleClientPickerNotice  string `json:"scheduleClientPickerNotice"`
	ScheduleGeneralPickerNotice string `json:"scheduleGeneralPickerNotice"`

	// 2026-04-30 cyclic-subscription-jobs plan §9.4 — client-side block
	// surfaced as a tooltip on the disabled MILESTONE option in the
	// billing_kind dropdown when the parent Plan is cyclic.
	MilestoneCyclicBlock string `json:"milestoneCyclicBlock"`

	// 2026-05-01 ad-hoc-subscription-billing plan §6 — client-side guards
	// surfaced as drawer warnings / tooltips on the disabled options. The
	// server enforces the same rules in validate_ad_hoc.go.
	AdHocPoolNoTemplate           string `json:"adHocPoolNoTemplate"`
	AdHocPerCallNoTemplate        string `json:"adHocPerCallNoTemplate"`
	AdHocNoEntitlement            string `json:"adHocNoEntitlement"`
	AdHocBillingCycleNotAllowed   string `json:"adHocBillingCycleNotAllowed"`
	AdHocVisitsPerCycleNotAllowed string `json:"adHocVisitsPerCycleNotAllowed"`
}
