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
	SpawnedPOs SpawnedPOLabels `json:"spawnedPos"`
	Form       FormLabels      `json:"form"`
	Empty      EmptyLabels     `json:"empty"`

	// SPS Wave 3 — F1/F2/F3 + CRIT-3 spawn lifecycle
	Filters              FilterLabels              `json:"filters"`
	FulfillmentStrategy  FulfillmentStrategyLabels `json:"fulfillmentStrategy"`
	FulfillmentMode      FulfillmentModeLabels     `json:"fulfillmentMode"`
	FulfillmentModeHints FulfillmentModeHintLabels `json:"fulfillmentModeHints"`
	Spawn                SpawnLabels               `json:"spawn"`
	PolicyDecision       PolicyDecisionLabels      `json:"policyDecision"`
}

// FilterLabels — F3 filter chips on the list page.
type FilterLabels struct {
	All                    string `json:"all"`
	Status                 string `json:"status"`
	FulfillmentStrategy    string `json:"fulfillmentStrategy"`
	FulfillmentMode        string `json:"fulfillmentMode"`
	AnyStatus              string `json:"anyStatus"`
	AnyFulfillmentStrategy string `json:"anyFulfillmentStrategy"`
	AnyFulfillmentMode     string `json:"anyFulfillmentMode"`
}

// FulfillmentStrategyLabels — F3 strategy values for header-level rollup.
type FulfillmentStrategyLabels struct {
	UniformOutright  string `json:"uniformOutright"`
	UniformStockable string `json:"uniformStockable"`
	UniformRecurring string `json:"uniformRecurring"`
	UniformPetty     string `json:"uniformPetty"`
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
	StatusColumn      string `json:"statusColumn"`
	StatusPending     string `json:"statusPending"`
	StatusSpawning    string `json:"statusSpawning"`
	StatusSpawned     string `json:"statusSpawned"`
	StatusFailed      string `json:"statusFailed"`
	StatusUnspecified string `json:"statusUnspecified"`
	ModeColumn        string `json:"modeColumn"`
	SpawnedColumn     string `json:"spawnedColumn"`
	LinkPO            string `json:"linkPo"`
	LinkContract      string `json:"linkContract"`
	LinkExpenditure   string `json:"linkExpenditure"`
	NotApplicable     string `json:"notApplicable"`
	ErrorPrefix       string `json:"errorPrefix"`
	RetryButton       string `json:"retryButton"`
	RetryConfirm      string `json:"retryConfirm"`
}

// PolicyDecisionLabels — policy_decision_log section on Info tab.
type PolicyDecisionLabels struct {
	SectionTitle string `json:"sectionTitle"`
	Toggle       string `json:"toggle"`
	EmptyMessage string `json:"emptyMessage"`
	Info         string `json:"info"`
}

type PageLabels struct {
	Heading                string `json:"heading"`
	HeadingDraft           string `json:"headingDraft"`
	HeadingSubmitted       string `json:"headingSubmitted"`
	HeadingPendingApproval string `json:"headingPendingApproval"`
	HeadingApproved        string `json:"headingApproved"`
	HeadingRejected        string `json:"headingRejected"`
	HeadingFulfilled       string `json:"headingFulfilled"`
	HeadingCancelled       string `json:"headingCancelled"`
	Caption                string `json:"caption"`
	AddButton              string `json:"addButton"`
	DetailSubtitle         string `json:"detailSubtitle"`
}

type ColumnLabels struct {
	RequestNumber  string `json:"requestNumber"`
	Status         string `json:"status"`
	Requester      string `json:"requester"`
	Supplier       string `json:"supplier"`
	EstimatedTotal string `json:"estimatedTotal"`
	NeededBy       string `json:"neededBy"`
	DateCreated    string `json:"dateCreated"`
}

type TabLabels struct {
	Info          string `json:"info"`
	Lines         string `json:"lines"`
	SpawnedPOs    string `json:"spawnedPos"`
	Activity      string `json:"activity"`
	ActivityEmpty string `json:"activityEmpty"`
}

type DetailLabels struct {
	InfoSection    string `json:"infoSection"`
	RequestNumber  string `json:"requestNumber"`
	Status         string `json:"status"`
	Requester      string `json:"requester"`
	Supplier       string `json:"supplier"`
	Currency       string `json:"currency"`
	EstimatedTotal string `json:"estimatedTotal"`
	NeededBy       string `json:"neededBy"`
	DateCreated    string `json:"dateCreated"`
	ApprovedBy     string `json:"approvedBy"`
	Justification  string `json:"justification"`
	TabAttachments string `json:"tabAttachments"`
}

type LineLabels struct {
	// Column labels
	Description         string `json:"description"`
	LineType            string `json:"lineType"`
	Quantity            string `json:"quantity"`
	EstimatedUnitPrice  string `json:"estimatedUnitPrice"`
	EstimatedTotalPrice string `json:"estimatedTotalPrice"`
	EmptyTitle          string `json:"emptyTitle"`
	EmptyMessage        string `json:"emptyMessage"`
	AddLine             string `json:"addLine"`

	// Enum label values for line_type
	LineTypeGoods   string `json:"lineTypeGoods"`
	LineTypeService string `json:"lineTypeService"`
	LineTypeExpense string `json:"lineTypeExpense"`

	// Drawer form labels
	FormDescription                    string `json:"formDescription"`
	FormDescriptionPlaceholder         string `json:"formDescriptionPlaceholder"`
	FormLineType                       string `json:"formLineType"`
	FormLineTypeInfo                   string `json:"formLineTypeInfo"`
	FormProduct                        string `json:"formProduct"`
	FormProductPlaceholder             string `json:"formProductPlaceholder"`
	FormQuantity                       string `json:"formQuantity"`
	FormQuantityInfo                   string `json:"formQuantityInfo"`
	FormEstimatedUnitPrice             string `json:"formEstimatedUnitPrice"`
	FormEstimatedUnitPriceInfo         string `json:"formEstimatedUnitPriceInfo"`
	FormEstimatedTotalPrice            string `json:"formEstimatedTotalPrice"`
	FormEstimatedTotalPriceHint        string `json:"formEstimatedTotalPriceHint"`
	FormExpenditureCategory            string `json:"formExpenditureCategory"`
	FormExpenditureCategoryPlaceholder string `json:"formExpenditureCategoryPlaceholder"`
	FormLocation                       string `json:"formLocation"`
	FormLocationPlaceholder            string `json:"formLocationPlaceholder"`
	FormLineNumber                     string `json:"formLineNumber"`

	// SPS Wave 3 — F1 fulfillment_mode picker + RECURRING fields + PETTY hint
	FormFulfillmentMode     string `json:"formFulfillmentMode"`
	FormFulfillmentModeInfo string `json:"formFulfillmentModeInfo"`
	FormFulfillmentModeHint string `json:"formFulfillmentModeHint"`

	FormRecurringSection    string `json:"formRecurringSection"`
	FormRecurringCycleValue string `json:"formRecurringCycleValue"`
	FormRecurringCycleUnit  string `json:"formRecurringCycleUnit"`
	FormRecurringTermValue  string `json:"formRecurringTermValue"`
	FormRecurringTermUnit   string `json:"formRecurringTermUnit"`
	FormRecurringCycleHint  string `json:"formRecurringCycleHint"`
	FormRecurringTermHint   string `json:"formRecurringTermHint"`
	FormRecurringUnitDay    string `json:"formRecurringUnitDay"`
	FormRecurringUnitWeek   string `json:"formRecurringUnitWeek"`
	FormRecurringUnitMonth  string `json:"formRecurringUnitMonth"`
	FormRecurringUnitYear   string `json:"formRecurringUnitYear"`

	FormPettyHint string `json:"formPettyHint"`

	// CRIT-3 spawn lifecycle column on the lines table
	ModeBadgeColumn string `json:"modeBadgeColumn"`
}

type SpawnedPOLabels struct {
	PONumber     string `json:"poNumber"`
	Status       string `json:"status"`
	TotalAmount  string `json:"totalAmount"`
	OrderDate    string `json:"orderDate"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// FormLabels holds all form-level labels for the drawer form.
type FormLabels struct {
	// Section headers
	SectionIdentity  string `json:"sectionIdentity"`
	SectionFinancial string `json:"sectionFinancial"`
	SectionApproval  string `json:"sectionApproval"`
	SectionOthers    string `json:"sectionOthers"`

	// §1 Identity
	RequestNumber            string `json:"requestNumber"`
	RequestNumberPlaceholder string `json:"requestNumberPlaceholder"`
	RequestNumberInfo        string `json:"requestNumberInfo"`
	RequesterUser            string `json:"requesterUser"`
	RequesterUserPlaceholder string `json:"requesterUserPlaceholder"`
	Supplier                 string `json:"supplier"`
	SupplierPlaceholder      string `json:"supplierPlaceholder"`
	SupplierHint             string `json:"supplierHint"`
	Location                 string `json:"location"`
	LocationPlaceholder      string `json:"locationPlaceholder"`

	// §2 Financial
	Currency           string `json:"currency"`
	CurrencyInfo       string `json:"currencyInfo"`
	EstimatedTotal     string `json:"estimatedTotal"`
	EstimatedTotalInfo string `json:"estimatedTotalInfo"`

	// §3 Timing & Approval
	NeededByDate               string `json:"neededByDate"`
	NeededByDateInfo           string `json:"neededByDateInfo"`
	Status                     string `json:"status"`
	StatusInfo                 string `json:"statusInfo"`
	StatusDraft                string `json:"statusDraft"`
	StatusSubmitted            string `json:"statusSubmitted"`
	StatusPendingApproval      string `json:"statusPendingApproval"`
	StatusApproved             string `json:"statusApproved"`
	StatusApprovedPendingSpawn string `json:"statusApprovedPendingSpawn"`
	StatusRejected             string `json:"statusRejected"`
	StatusFulfilled            string `json:"statusFulfilled"`
	StatusCancelled            string `json:"statusCancelled"`
	ApprovedBy                 string `json:"approvedBy"`

	// §4 Others
	Justification            string `json:"justification"`
	JustificationPlaceholder string `json:"justificationPlaceholder"`
	Notes                    string `json:"notes"`
	NotesPlaceholder         string `json:"notesPlaceholder"`
	Active                   string `json:"active"`

	// Action buttons
	Edit      string `json:"edit"`
	EditTitle string `json:"editTitle"`
	Submit    string `json:"submit"`
	Approve   string `json:"approve"`
	Reject    string `json:"reject"`
	SpawnPO   string `json:"spawnPo"`
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
