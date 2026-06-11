package subscription

// ---------------------------------------------------------------------------
// Price Plan labels
// ---------------------------------------------------------------------------

// PricePlanLabels holds all labels for the standalone price plan (rate card) module.
type PricePlanLabels struct {
	Page         PricePlanPageLabels         `json:"page"`
	Buttons      PricePlanButtonLabels       `json:"buttons"`
	Columns      PricePlanColumnLabels2      `json:"columns"`
	Empty        PricePlanEmptyLabels        `json:"empty"`
	Form         PricePlanFormLabels         `json:"form"`
	Actions      PricePlanActionLabels       `json:"actions"`
	Bulk         PricePlanBulkLabels         `json:"bulk"`
	Detail       PricePlanDetailLabels2      `json:"detail"`
	Tabs         PricePlanTabLabels2         `json:"tabs"`
	Confirm      PricePlanConfirmLabels      `json:"confirm"`
	Errors       PricePlanErrorLabels        `json:"errors"`
	ProductPrice PricePlanProductPriceLabels `json:"productPrice"`
	Messages     PricePlanMessageLabels      `json:"messages"`
}

// PricePlanProductPriceLabels holds labels for product-price sub-table actions and empty state.
type PricePlanProductPriceLabels struct {
	EditTitle   string `json:"editTitle"`
	DeleteTitle string `json:"deleteTitle"`
	EmptyTitle  string `json:"emptyTitle"`
	EmptyMsg    string `json:"emptyMsg"`
}

// PricePlanMessageLabels holds translatable message strings used in the price plan
// and price schedule plan views (pricing-lock notices, validation errors).
type PricePlanMessageLabels struct {
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

type PricePlanPageLabels struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	ActiveTitle   string `json:"activeTitle"`
	InactiveTitle string `json:"inactiveTitle"`
}

type PricePlanButtonLabels struct {
	View       string `json:"view"`
	Add        string `json:"add"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	BulkDelete string `json:"bulkDelete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type PricePlanColumnLabels2 struct {
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

type PricePlanEmptyLabels struct {
	Title       string `json:"title"`
	Message     string `json:"message"`
	Description string `json:"description"`
	ActionLabel string `json:"actionLabel"`
}

type PricePlanActionLabels struct {
	CreateSuccess string `json:"createSuccess"`
	CreateError   string `json:"createError"`
	UpdateSuccess string `json:"updateSuccess"`
	UpdateError   string `json:"updateError"`
	DeleteSuccess string `json:"deleteSuccess"`
	DeleteError   string `json:"deleteError"`
}

type PricePlanBulkLabels struct {
	DeleteTitle   string `json:"deleteTitle"`
	DeleteMessage string `json:"deleteMessage"`
	StatusTitle   string `json:"statusTitle"`
	StatusMessage string `json:"statusMessage"`
}

type PricePlanDetailLabels2 struct {
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
	SummaryHeading            string                      `json:"summaryHeading"`
	CustomerHeading           string                      `json:"customerHeading"`
	OperationsHeading         string                      `json:"operationsHeading"`
	RevenueRecognitionHeading string                      `json:"revenueRecognitionHeading"`
	Summary                   PricePlanBillingSummaryCopy `json:"summary"`
	Warning                   PricePlanBillingSummaryWarn `json:"warning"`

	// 2026-05-04 — Subscriptions/Engagements tab on the price-plan detail.
	// See docs/plan/20260504-price-plan-engagements-tab/.
	Subscriptions PricePlanSubscriptionsSectionLabels `json:"subscriptions"`
}

// PricePlanBillingSummaryCopy carries the per-(kind × basis) lyngua copy that
// `buildBillingModelSummary` projects into the info-sections grid.
// Cyclic plan ships rows 1-6 (oneTime / recurring / contract / milestone);
// AD_HOC plan adds adHoc.* in a follow-up. Each entry has 3 lines:
// customer, operations, revenue.
type PricePlanBillingSummaryCopy struct {
	OneTime   PricePlanSummaryByBasis `json:"oneTime"`
	Recurring PricePlanSummaryByBasis `json:"recurring"`
	Contract  PricePlanSummaryByBasis `json:"contract"`
	Milestone PricePlanSummaryByBasis `json:"milestone"`
	AdHoc     PricePlanSummaryByBasis `json:"adHoc"`
}

// PricePlanSummaryByBasis groups the text lines per basis. Empty
// strings on a basis means "no copy for that combo" — view skips it.
type PricePlanSummaryByBasis struct {
	PerCycle         PricePlanSummaryLines `json:"perCycle"`
	TotalPackage     PricePlanSummaryLines `json:"totalPackage"`
	DerivedFromLines PricePlanSummaryLines `json:"derivedFromLines"`
	PerOccurrence    PricePlanSummaryLines `json:"perOccurrence"`
}

// PricePlanSummaryLines holds the 3 lines for a kind × basis cell.
type PricePlanSummaryLines struct {
	Customer   string `json:"customer"`
	Operations string `json:"operations"`
	Revenue    string `json:"revenue"`
}

// PricePlanBillingSummaryWarn carries the warning-row copy keyed by symbol
// per plan §20.3. View only renders entries whose preconditions trip.
type PricePlanBillingSummaryWarn struct {
	MilestoneNoTemplate           string `json:"milestoneNoTemplate"`
	RecurringNoTemplate           string `json:"recurringNoTemplate"`
	VisitsPerCycleInvalidKind     string `json:"visitsPerCycleInvalidKind"`
	AdHocPoolNoTemplate           string `json:"adHocPoolNoTemplate"`
	AdHocPerCallNoTemplate        string `json:"adHocPerCallNoTemplate"`
	AdHocNoEntitlement            string `json:"adHocNoEntitlement"`
	AdHocBillingCycleNotAllowed   string `json:"adHocBillingCycleNotAllowed"`
	AdHocVisitsPerCycleNotAllowed string `json:"adHocVisitsPerCycleNotAllowed"`
}

type PricePlanTabLabels2 struct {
	Info          string `json:"info"`
	Products      string `json:"products"`
	Subscriptions string `json:"subscriptions"`
	Attachments   string `json:"attachments"`
	Audit         string `json:"audit"`
}

// PricePlanSubscriptionsSectionLabels holds the column headers, empty state,
// and confirm-delete copy for the price-plan detail "Subscriptions" tab —
// professional tier overrides this block to use the engagement vocabulary.
type PricePlanSubscriptionsSectionLabels struct {
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

type PricePlanConfirmLabels struct {
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

type PricePlanErrorLabels struct {
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

// DefaultPricePlanLabels returns PricePlanLabels with sensible English defaults.
func DefaultPricePlanLabels() PricePlanLabels {
	return PricePlanLabels{
		Page: PricePlanPageLabels{
			Title:         "Rate Cards",
			Subtitle:      "Manage your rate cards",
			ActiveTitle:   "Active Rate Cards",
			InactiveTitle: "Inactive Rate Cards",
		},
		Buttons: PricePlanButtonLabels{
			View:       "View",
			Add:        "Add Rate Card",
			Edit:       "Edit Rate Card",
			Delete:     "Delete Rate Card",
			BulkDelete: "Delete Rate Cards",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: PricePlanColumnLabels2{
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
		Empty: PricePlanEmptyLabels{
			Title:       "No Rate Cards",
			Message:     "No rate cards to display.",
			Description: "Add a rate card to define pricing for your plans.",
			ActionLabel: "Add Rate Card",
		},
		Form: PricePlanFormLabels{
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
			// (under "form", not "fields" — the Go struct lives on PricePlanFormLabels).
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
		Actions: PricePlanActionLabels{
			CreateSuccess: "Rate card created successfully.",
			CreateError:   "Failed to create rate card.",
			UpdateSuccess: "Rate card updated successfully.",
			UpdateError:   "Failed to update rate card.",
			DeleteSuccess: "Rate card deleted successfully.",
			DeleteError:   "Failed to delete rate card.",
		},
		Bulk: PricePlanBulkLabels{
			DeleteTitle:   "Delete Rate Cards",
			DeleteMessage: "Are you sure you want to delete the selected rate cards?",
			StatusTitle:   "Update Status",
			StatusMessage: "Are you sure you want to update the status of the selected rate cards?",
		},
		Detail: PricePlanDetailLabels2{
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
			Summary: PricePlanBillingSummaryCopy{
				OneTime: PricePlanSummaryByBasis{
					TotalPackage: PricePlanSummaryLines{
						Customer:   "Pays {{.Amount}} once at signup. No further charges.",
						Operations: "Engagement spawns 1 lifetime Job with phases (if Plan has a JobTemplate).",
						Revenue:    "One Revenue at Subscription.Create covering the full amount.",
					},
				},
				Recurring: PricePlanSummaryByBasis{
					PerCycle: PricePlanSummaryLines{
						Customer:   "Charged {{.Amount}} every {{.CycleLabel}}. Subscription auto-renews until cancelled.",
						Operations: "Each cycle spawns {{.VisitsPerCycle}} cycle Job(s) (if Plan has a JobTemplate). Operations tab shows cycle accordions.",
						Revenue:    "One Revenue per cycle. Recognize Revenue creates the invoice and (via piggyback) spawns the cycle Job if missing.",
					},
					DerivedFromLines: PricePlanSummaryLines{
						Customer:   "Charged the sum of itemised lines every {{.CycleLabel}}.",
						Operations: "Each cycle spawns 1+ cycle Jobs. Operations tracking flows through Plan's JobTemplate.",
						Revenue:    "Revenue total computed from ProductPricePlan rows; one Revenue per cycle.",
					},
				},
				Contract: PricePlanSummaryByBasis{
					PerCycle: PricePlanSummaryLines{
						Customer:   "Charged {{.Amount}} every {{.CycleLabel}} for {{.TermLength}}. Auto-deactivates at term end.",
						Operations: "Same as recurring + the engagement closes when the {{.TermLength}} term completes.",
						Revenue:    "Same as recurring. Operator can extend the term to spawn additional cycles.",
					},
					TotalPackage: PricePlanSummaryLines{
						Customer:   "Pays {{.Amount}} upfront for {{.TermLength}} of service.",
						Operations: "Engagement spawns 1 lifetime Job (or cycle Jobs if cyclic — see Plan's visits_per_cycle).",
						Revenue:    "One Revenue at signup; cycle Jobs are operational only.",
					},
				},
				Milestone: PricePlanSummaryByBasis{
					TotalPackage: PricePlanSummaryLines{
						Customer:   "Pays {{.Amount}} total. Invoice fires per milestone (engagement phase) as work completes.",
						Operations: "Lifetime engagement Job with phases. BillingEvent rows gate per-milestone invoicing.",
						Revenue:    "Revenue per milestone trigger; sum across milestones equals the total package.",
					},
				},
			},
			Warning: PricePlanBillingSummaryWarn{
				MilestoneNoTemplate:       "Milestone billing requires the Plan to have a JobTemplate. Configure it on the Plan first.",
				RecurringNoTemplate:       "This subscription will not have operational tracking. Add a JobTemplate to the Plan to enable cycle Jobs.",
				VisitsPerCycleInvalidKind: "visits_per_cycle is only valid for cyclic plans. Reset to 1 or change the billing kind.",
			},
			Subscriptions: PricePlanSubscriptionsSectionLabels{
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
		Tabs: PricePlanTabLabels2{
			Info:          "Information",
			Products:      "Products",
			Subscriptions: "Subscriptions",
			Attachments:   "Attachments",
			Audit:         "Audit Trail",
		},
		Confirm: PricePlanConfirmLabels{
			DeleteTitle:       "Delete Rate Card",
			DeleteMessage:     "Are you sure you want to delete this rate card? This action cannot be undone.",
			DeactivateTitle:   "Deactivate Rate Card",
			DeactivateMessage: "Are you sure you want to deactivate this rate card?",
			// 2026-04-27 plan-client-scope plan §3.5 / §7.
			EditAmountMultipleSubscriptions: "This price plan is attached to {{.Count}} active subscriptions for {{.ClientName}}. Changing the amount or cycle will affect all of them on the next bill cycle. Continue?",
		},
		Errors: PricePlanErrorLabels{
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
		ProductPrice: PricePlanProductPriceLabels{
			EditTitle:   "Edit Product Price",
			DeleteTitle: "Delete Product Price",
			EmptyTitle:  "No Product Prices",
			EmptyMsg:    "No product prices have been configured for this rate card yet.",
		},
		Messages: PricePlanMessageLabels{
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
