package procurement_request

// ---------------------------------------------------------------------------
// ProcurementRequest labels  (P3a)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the procurement_request module.
type Labels struct {
	Page       PageLabels      `json:"page"`
	Columns    ColumnLabels    `json:"columns"`
	Tabs       TabLabels       `json:"tabs"`
	Detail     DetailLabels    `json:"detail"`
	Lines      LineLabels      `json:"lines"`
	SpawnedPOs SpawnedPOLabels `json:"spawned_pos"`
	Form       FormLabels      `json:"form"`
	Empty      EmptyLabels     `json:"empty"`

	// SPS Wave 3 — F1/F2/F3 + CRIT-3 spawn lifecycle
	Filters              FilterLabels              `json:"filters"`
	FulfillmentStrategy  FulfillmentStrategyLabels `json:"fulfillment_strategy"`
	FulfillmentMode      FulfillmentModeLabels     `json:"fulfillment_mode"`
	FulfillmentModeHints FulfillmentModeHintLabels `json:"fulfillment_mode_hints"`
	Spawn                SpawnLabels               `json:"spawn"`
	PolicyDecision       PolicyDecisionLabels      `json:"policy_decision"`
}

// FilterLabels — F3 filter chips on the list page.
type FilterLabels struct {
	All                    string `json:"all"`
	Status                 string `json:"status"`
	FulfillmentStrategy    string `json:"fulfillment_strategy"`
	FulfillmentMode        string `json:"fulfillment_mode"`
	AnyStatus              string `json:"any_status"`
	AnyFulfillmentStrategy string `json:"any_fulfillment_strategy"`
	AnyFulfillmentMode     string `json:"any_fulfillment_mode"`
}

// FulfillmentStrategyLabels — F3 strategy values for header-level rollup.
type FulfillmentStrategyLabels struct {
	UniformOutright  string `json:"uniform_outright"`
	UniformStockable string `json:"uniform_stockable"`
	UniformRecurring string `json:"uniform_recurring"`
	UniformPetty     string `json:"uniform_petty"`
	Mixed            string `json:"mixed"`
	Hint             string `json:"hint"`
}

// FulfillmentModeLabels — F1 line-level mode values.
type FulfillmentModeLabels struct {
	Outright  string `json:"outright"`
	Stockable string `json:"stockable"`
	Recurring string `json:"recurring"`
	Petty     string `json:"petty"`
}

// FulfillmentModeHintLabels — F1 short hints rendered under each radio choice.
type FulfillmentModeHintLabels struct {
	Outright  string `json:"outright"`
	Stockable string `json:"stockable"`
	Recurring string `json:"recurring"`
	Petty     string `json:"petty"`
}

// SpawnLabels — CRIT-3 spawn lifecycle UI strings.
type SpawnLabels struct {
	StatusColumn      string `json:"status_column"`
	StatusPending     string `json:"status_pending"`
	StatusSpawning    string `json:"status_spawning"`
	StatusSpawned     string `json:"status_spawned"`
	StatusFailed      string `json:"status_failed"`
	StatusUnspecified string `json:"status_unspecified"`
	ModeColumn        string `json:"mode_column"`
	SpawnedColumn     string `json:"spawned_column"`
	LinkPO            string `json:"link_po"`
	LinkContract      string `json:"link_contract"`
	LinkExpenditure   string `json:"link_expenditure"`
	NotApplicable     string `json:"not_applicable"`
	ErrorPrefix       string `json:"error_prefix"`
	RetryButton       string `json:"retry_button"`
	RetryConfirm      string `json:"retry_confirm"`
}

// PolicyDecisionLabels — policy_decision_log section on Info tab.
type PolicyDecisionLabels struct {
	SectionTitle string `json:"section_title"`
	Toggle       string `json:"toggle"`
	EmptyMessage string `json:"empty_message"`
	Info         string `json:"info"`
}

type PageLabels struct {
	Heading                string `json:"heading"`
	HeadingDraft           string `json:"heading_draft"`
	HeadingSubmitted       string `json:"heading_submitted"`
	HeadingPendingApproval string `json:"heading_pending_approval"`
	HeadingApproved        string `json:"heading_approved"`
	HeadingRejected        string `json:"heading_rejected"`
	HeadingFulfilled       string `json:"heading_fulfilled"`
	HeadingCancelled       string `json:"heading_cancelled"`
	Caption                string `json:"caption"`
	AddButton              string `json:"add_button"`
	DetailSubtitle         string `json:"detail_subtitle"`
}

type ColumnLabels struct {
	RequestNumber  string `json:"request_number"`
	Status         string `json:"status"`
	Requester      string `json:"requester"`
	Supplier       string `json:"supplier"`
	EstimatedTotal string `json:"estimated_total"`
	NeededBy       string `json:"needed_by"`
	DateCreated    string `json:"date_created"`
}

type TabLabels struct {
	Info          string `json:"info"`
	Lines         string `json:"lines"`
	SpawnedPOs    string `json:"spawned_pos"`
	Activity      string `json:"activity"`
	ActivityEmpty string `json:"activity_empty"`
}

type DetailLabels struct {
	InfoSection    string `json:"info_section"`
	RequestNumber  string `json:"request_number"`
	Status         string `json:"status"`
	Requester      string `json:"requester"`
	Supplier       string `json:"supplier"`
	Currency       string `json:"currency"`
	EstimatedTotal string `json:"estimated_total"`
	NeededBy       string `json:"needed_by"`
	DateCreated    string `json:"date_created"`
	ApprovedBy     string `json:"approved_by"`
	Justification  string `json:"justification"`
	TabAttachments string `json:"tab_attachments"`
}

type LineLabels struct {
	// Column labels
	Description         string `json:"description"`
	LineType            string `json:"line_type"`
	Quantity            string `json:"quantity"`
	EstimatedUnitPrice  string `json:"estimated_unit_price"`
	EstimatedTotalPrice string `json:"estimated_total_price"`
	EmptyTitle          string `json:"empty_title"`
	EmptyMessage        string `json:"empty_message"`
	AddLine             string `json:"add_line"`

	// Enum label values for line_type
	LineTypeGoods   string `json:"line_type_goods"`
	LineTypeService string `json:"line_type_service"`
	LineTypeExpense string `json:"line_type_expense"`

	// Drawer form labels
	FormDescription                    string `json:"form_description"`
	FormDescriptionPlaceholder         string `json:"form_description_placeholder"`
	FormLineType                       string `json:"form_line_type"`
	FormLineTypeInfo                   string `json:"form_line_type_info"`
	FormProduct                        string `json:"form_product"`
	FormProductPlaceholder             string `json:"form_product_placeholder"`
	FormQuantity                       string `json:"form_quantity"`
	FormQuantityInfo                   string `json:"form_quantity_info"`
	FormEstimatedUnitPrice             string `json:"form_estimated_unit_price"`
	FormEstimatedUnitPriceInfo         string `json:"form_estimated_unit_price_info"`
	FormEstimatedTotalPrice            string `json:"form_estimated_total_price"`
	FormEstimatedTotalPriceHint        string `json:"form_estimated_total_price_hint"`
	FormExpenditureCategory            string `json:"form_expenditure_category"`
	FormExpenditureCategoryPlaceholder string `json:"form_expenditure_category_placeholder"`
	FormLocation                       string `json:"form_location"`
	FormLocationPlaceholder            string `json:"form_location_placeholder"`
	FormLineNumber                     string `json:"form_line_number"`

	// SPS Wave 3 — F1 fulfillment_mode picker + RECURRING fields + PETTY hint
	FormFulfillmentMode     string `json:"form_fulfillment_mode"`
	FormFulfillmentModeInfo string `json:"form_fulfillment_mode_info"`
	FormFulfillmentModeHint string `json:"form_fulfillment_mode_hint"`

	FormRecurringSection    string `json:"form_recurring_section"`
	FormRecurringCycleValue string `json:"form_recurring_cycle_value"`
	FormRecurringCycleUnit  string `json:"form_recurring_cycle_unit"`
	FormRecurringTermValue  string `json:"form_recurring_term_value"`
	FormRecurringTermUnit   string `json:"form_recurring_term_unit"`
	FormRecurringCycleHint  string `json:"form_recurring_cycle_hint"`
	FormRecurringTermHint   string `json:"form_recurring_term_hint"`
	FormRecurringUnitDay    string `json:"form_recurring_unit_day"`
	FormRecurringUnitWeek   string `json:"form_recurring_unit_week"`
	FormRecurringUnitMonth  string `json:"form_recurring_unit_month"`
	FormRecurringUnitYear   string `json:"form_recurring_unit_year"`

	FormPettyHint string `json:"form_petty_hint"`

	// CRIT-3 spawn lifecycle column on the lines table
	ModeBadgeColumn string `json:"mode_badge_column"`
}

type SpawnedPOLabels struct {
	PONumber     string `json:"po_number"`
	Status       string `json:"status"`
	TotalAmount  string `json:"total_amount"`
	OrderDate    string `json:"order_date"`
	EmptyTitle   string `json:"empty_title"`
	EmptyMessage string `json:"empty_message"`
}

// FormLabels holds all form-level labels for the drawer form.
type FormLabels struct {
	// Section headers
	SectionIdentity  string `json:"section_identity"`
	SectionFinancial string `json:"section_financial"`
	SectionApproval  string `json:"section_approval"`
	SectionOthers    string `json:"section_others"`

	// §1 Identity
	RequestNumber            string `json:"request_number"`
	RequestNumberPlaceholder string `json:"request_number_placeholder"`
	RequestNumberInfo        string `json:"request_number_info"`
	RequesterUser            string `json:"requester_user"`
	RequesterUserPlaceholder string `json:"requester_user_placeholder"`
	Supplier                 string `json:"supplier"`
	SupplierPlaceholder      string `json:"supplier_placeholder"`
	SupplierHint             string `json:"supplier_hint"`
	Location                 string `json:"location"`
	LocationPlaceholder      string `json:"location_placeholder"`

	// §2 Financial
	Currency           string `json:"currency"`
	CurrencyInfo       string `json:"currency_info"`
	EstimatedTotal     string `json:"estimated_total"`
	EstimatedTotalInfo string `json:"estimated_total_info"`

	// §3 Timing & Approval
	NeededByDate               string `json:"needed_by_date"`
	NeededByDateInfo           string `json:"needed_by_date_info"`
	Status                     string `json:"status"`
	StatusInfo                 string `json:"status_info"`
	StatusDraft                string `json:"status_draft"`
	StatusSubmitted            string `json:"status_submitted"`
	StatusPendingApproval      string `json:"status_pending_approval"`
	StatusApproved             string `json:"status_approved"`
	StatusApprovedPendingSpawn string `json:"status_approved_pending_spawn"`
	StatusRejected             string `json:"status_rejected"`
	StatusFulfilled            string `json:"status_fulfilled"`
	StatusCancelled            string `json:"status_cancelled"`
	ApprovedBy                 string `json:"approved_by"`

	// §4 Others
	Justification            string `json:"justification"`
	JustificationPlaceholder string `json:"justification_placeholder"`
	Notes                    string `json:"notes"`
	NotesPlaceholder         string `json:"notes_placeholder"`
	Active                   string `json:"active"`

	// Action buttons
	Edit      string `json:"edit"`
	EditTitle string `json:"edit_title"`
	Submit    string `json:"submit"`
	Approve   string `json:"approve"`
	Reject    string `json:"reject"`
	SpawnPO   string `json:"spawn_po"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultLabels returns English fallback labels.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:                "Procurement Requests",
			HeadingDraft:           "Draft Requests",
			HeadingSubmitted:       "Submitted Requests",
			HeadingPendingApproval: "Pending Approval",
			HeadingApproved:        "Approved Requests",
			HeadingRejected:        "Rejected Requests",
			HeadingFulfilled:       "Fulfilled Requests",
			HeadingCancelled:       "Cancelled Requests",
			Caption:                "Internal purchase intent records",
			AddButton:              "New Request",
			DetailSubtitle:         "Procurement request details",
		},
		Columns: ColumnLabels{
			RequestNumber:  "Request #",
			Status:         "Status",
			Requester:      "Requester",
			Supplier:       "Supplier",
			EstimatedTotal: "Estimated Total",
			NeededBy:       "Needed By",
			DateCreated:    "Created",
		},
		Tabs: TabLabels{
			Info:          "Info",
			Lines:         "Lines",
			SpawnedPOs:    "Spawned POs",
			Activity:      "Activity",
			ActivityEmpty: "No activity recorded yet.",
		},
		Detail: DetailLabels{
			InfoSection:    "Request Information",
			RequestNumber:  "Request Number",
			Status:         "Status",
			Requester:      "Requester",
			Supplier:       "Supplier",
			Currency:       "Currency",
			EstimatedTotal: "Estimated Total",
			NeededBy:       "Needed By",
			DateCreated:    "Created",
			ApprovedBy:     "Approved By",
			Justification:  "Justification",
			TabAttachments: "Attachments",
		},
		Lines: LineLabels{
			Description:                        "Description",
			LineType:                           "Line Type",
			Quantity:                           "Quantity",
			EstimatedUnitPrice:                 "Est. Unit Price",
			EstimatedTotalPrice:                "Est. Total",
			EmptyTitle:                         "No lines yet",
			EmptyMessage:                       "Add a line to this request.",
			AddLine:                            "Add Line",
			LineTypeGoods:                      "Goods",
			LineTypeService:                    "Service",
			LineTypeExpense:                    "Expense",
			FormDescription:                    "Description",
			FormDescriptionPlaceholder:         "e.g. 50 laptop units",
			FormLineType:                       "Line Type",
			FormLineTypeInfo:                   "Goods = physical items; Service = intangible; Expense = direct cost",
			FormProduct:                        "Product",
			FormProductPlaceholder:             "Select a product (optional)",
			FormQuantity:                       "Quantity",
			FormQuantityInfo:                   "Number of units requested.",
			FormEstimatedUnitPrice:             "Estimated Unit Price",
			FormEstimatedUnitPriceInfo:         "Best estimate in centavos ÷ 100.",
			FormEstimatedTotalPrice:            "Estimated Total Price",
			FormEstimatedTotalPriceHint:        "Auto-calculated. Override if needed.",
			FormExpenditureCategory:            "Expenditure Category",
			FormExpenditureCategoryPlaceholder: "Select category",
			FormLocation:                       "Location",
			FormLocationPlaceholder:            "Branch or cost center",
			FormLineNumber:                     "Line Number",
			FormFulfillmentMode:                "Fulfillment Mode",
			FormFulfillmentModeInfo:            "How this line will be sourced after approval. Drives the downstream artifact created when the request is approved.",
			FormFulfillmentModeHint:            "Pick one — the spawn cascade dispatches per-line based on this choice.",
			FormRecurringSection:               "Recurring Schedule",
			FormRecurringCycleValue:            "Cycle Every",
			FormRecurringCycleUnit:             "Cycle Unit",
			FormRecurringTermValue:             "Term Length",
			FormRecurringTermUnit:              "Term Unit",
			FormRecurringCycleHint:             "Billing/delivery cadence (e.g. every 1 month).",
			FormRecurringTermHint:              "Total contract horizon (e.g. 24 months).",
			FormRecurringUnitDay:               "Day",
			FormRecurringUnitWeek:              "Week",
			FormRecurringUnitMonth:             "Month",
			FormRecurringUnitYear:              "Year",
			FormPettyHint:                      "Petty mode auto-approves under threshold and posts a direct expenditure. No PO, no contract.",
			ModeBadgeColumn:                    "Mode",
		},
		SpawnedPOs: SpawnedPOLabels{
			PONumber:     "PO Number",
			Status:       "Status",
			TotalAmount:  "Total Amount",
			OrderDate:    "Order Date",
			EmptyTitle:   "No purchase orders yet",
			EmptyMessage: "POs spawned from this request will appear here after approval.",
		},
		Form: FormLabels{
			SectionIdentity:            "Identity",
			SectionFinancial:           "Financial",
			SectionApproval:            "Timing & Approval",
			SectionOthers:              "Others",
			RequestNumber:              "Request Number",
			RequestNumberPlaceholder:   "e.g. PR-2026-001",
			RequestNumberInfo:          "A unique identifier for this procurement request.",
			RequesterUser:              "Requester",
			RequesterUserPlaceholder:   "User ID of requester",
			Supplier:                   "Supplier",
			SupplierPlaceholder:        "Select supplier (optional for RFQ)",
			SupplierHint:               "Leave empty if supplier is not yet selected (RFQ flow).",
			Location:                   "Location",
			LocationPlaceholder:        "Branch or cost center",
			Currency:                   "Currency",
			CurrencyInfo:               "ISO 4217 currency code (e.g. PHP, USD).",
			EstimatedTotal:             "Estimated Total",
			EstimatedTotalInfo:         "Best estimate of total spend (centavos ÷ 100 for display).",
			NeededByDate:               "Needed By",
			NeededByDateInfo:           "When the goods or services are required.",
			Status:                     "Status",
			StatusInfo:                 "Lifecycle stage. draft → submitted → pending_approval → approved/rejected → fulfilled/cancelled.",
			StatusDraft:                "Draft",
			StatusSubmitted:            "Submitted",
			StatusPendingApproval:      "Pending Approval",
			StatusApproved:             "Approved",
			StatusApprovedPendingSpawn: "Approved — Pending Spawn",
			StatusRejected:             "Rejected",
			StatusFulfilled:            "Fulfilled",
			StatusCancelled:            "Cancelled",
			ApprovedBy:                 "Approved By",
			Justification:              "Justification",
			JustificationPlaceholder:   "Business reason for this request",
			Notes:                      "Notes",
			NotesPlaceholder:           "Additional notes or context",
			Active:                     "Active",
			Edit:                       "Edit",
			EditTitle:                  "Edit Procurement Request",
			Submit:                     "Submit for Approval",
			Approve:                    "Approve",
			Reject:                     "Reject",
			SpawnPO:                    "Create PO",
		},
		Empty: EmptyLabels{
			Title:   "No procurement requests",
			Message: "Create a procurement request to start the approval workflow.",
		},
		Filters: FilterLabels{
			All:                    "All",
			Status:                 "Status",
			FulfillmentStrategy:    "Fulfillment",
			FulfillmentMode:        "Mode",
			AnyStatus:              "Any Status",
			AnyFulfillmentStrategy: "Any Fulfillment",
			AnyFulfillmentMode:     "Any Mode",
		},
		FulfillmentStrategy: FulfillmentStrategyLabels{
			UniformOutright:  "Uniform — Outright",
			UniformStockable: "Uniform — Stockable",
			UniformRecurring: "Uniform — Recurring",
			UniformPetty:     "Uniform — Petty",
			Mixed:            "Mixed Modes",
			Hint:             "Auto-derived from per-line fulfillment modes. Mixed = lines split across multiple modes.",
		},
		FulfillmentMode: FulfillmentModeLabels{
			Outright:  "Outright",
			Stockable: "Stockable",
			Recurring: "Recurring",
			Petty:     "Petty",
		},
		FulfillmentModeHints: FulfillmentModeHintLabels{
			Outright:  "One-shot purchase. Spawns a single purchase order on approval; no recurrence, no inventory side-effect.",
			Stockable: "Replenishment buy. Spawns a purchase order; received goods credit inventory on receipt.",
			Recurring: "Standing agreement. Spawns a supplier contract on approval; the recurrence engine emits cycle bills.",
			Petty:     "Cash-out. Spawns an expenditure directly against petty cash. No PO, no contract.",
		},
		Spawn: SpawnLabels{
			StatusColumn:      "Spawn Status",
			StatusPending:     "Pending",
			StatusSpawning:    "Spawning",
			StatusSpawned:     "Spawned",
			StatusFailed:      "Failed",
			StatusUnspecified: "—",
			ModeColumn:        "Mode",
			SpawnedColumn:     "Spawned Artifact",
			LinkPO:            "View PO line",
			LinkContract:      "View contract",
			LinkExpenditure:   "View expenditure",
			NotApplicable:     "—",
			ErrorPrefix:       "Error",
			RetryButton:       "Retry spawn",
			RetryConfirm:      "Retry spawning the downstream artifact for this line?",
		},
		PolicyDecision: PolicyDecisionLabels{
			SectionTitle: "Approval Policy Log",
			Toggle:       "Show / Hide",
			EmptyMessage: "No policy decisions logged yet.",
			Info:         "Audit trail of approval policy decisions taken on this request (auto-approve, escalation, override). Read-only.",
		},
	}
}
