package supplier_contract

// ---------------------------------------------------------------------------
// SupplierContract labels  (P3a)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the supplier_contract module.
type Labels struct {
	Page               PageLabels              `json:"page"`
	Columns            ColumnLabels            `json:"columns"`
	Tabs               TabLabels               `json:"tabs"`
	Detail             DetailLabels            `json:"detail"`
	Lines              LineLabels              `json:"lines"`
	LinkedPOs          LinkedPOLabels          `json:"linked_pos"`
	LinkedExpenditures LinkedExpenditureLabels `json:"linked_expenditures"`
	Form               FormLabels              `json:"form"`
	Empty              EmptyLabels             `json:"empty"`
}

type PageLabels struct {
	Heading           string `json:"heading"`
	HeadingDraft      string `json:"heading_draft"`
	HeadingActive     string `json:"heading_active"`
	HeadingExpiring   string `json:"heading_expiring"`
	HeadingExpired    string `json:"heading_expired"`
	HeadingTerminated string `json:"heading_terminated"`
	Caption           string `json:"caption"`
	AddButton         string `json:"add_button"`
	DetailSubtitle    string `json:"detail_subtitle"`
}

type ColumnLabels struct {
	Name      string `json:"name"`
	Supplier  string `json:"supplier"`
	Kind      string `json:"kind"`
	Status    string `json:"status"`
	Validity  string `json:"validity"`
	Committed string `json:"committed"`
	Released  string `json:"released"`
	Billed    string `json:"billed"`
	Remaining string `json:"remaining"`
}

type TabLabels struct {
	Info                string `json:"info"`
	Lines               string `json:"lines"`
	LinkedPOs           string `json:"linked_pos"`
	LinkedExpenditures  string `json:"linked_expenditures"`
	PriceSchedules      string `json:"price_schedules"`
	Activity            string `json:"activity"`
	ActivityEmpty       string `json:"activity_empty"`
	PriceSchedulesEmpty string `json:"price_schedules_empty"`
}

type DetailLabels struct {
	InfoSection     string `json:"info_section"`
	Name            string `json:"name"`
	Kind            string `json:"kind"`
	Status          string `json:"status"`
	Supplier        string `json:"supplier"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	AutoRenew       string `json:"auto_renew"`
	Currency        string `json:"currency"`
	CommittedAmount string `json:"committed_amount"`
	ReleasedAmount  string `json:"released_amount"`
	BilledAmount    string `json:"billed_amount"`
	RemainingAmount string `json:"remaining_amount"`
	Notes           string `json:"notes"`
	TabAttachments  string `json:"tab_attachments"`
}

type LineLabels struct {
	// Column labels
	Description  string `json:"description"`
	LineType     string `json:"line_type"`
	Quantity     string `json:"quantity"`
	UnitPrice    string `json:"unit_price"`
	Total        string `json:"total"`
	Treatment    string `json:"treatment"`
	EmptyTitle   string `json:"empty_title"`
	EmptyMessage string `json:"empty_message"`
	AddLine      string `json:"add_line"`

	// Enum label values for treatment
	TreatmentRecurring         string `json:"treatment_recurring"`
	TreatmentOneTime           string `json:"treatment_one_time"`
	TreatmentUsageBased        string `json:"treatment_usage_based"`
	TreatmentMinimumCommitment string `json:"treatment_minimum_commitment"`

	// Enum label values for line_type
	LineTypeGoods   string `json:"line_type_goods"`
	LineTypeService string `json:"line_type_service"`
	LineTypeExpense string `json:"line_type_expense"`

	// Drawer form labels
	FormDescription               string `json:"form_description"`
	FormDescriptionPlaceholder    string `json:"form_description_placeholder"`
	FormLineType                  string `json:"form_line_type"`
	FormLineTypeInfo              string `json:"form_line_type_info"`
	FormTreatment                 string `json:"form_treatment"`
	FormTreatmentInfo             string `json:"form_treatment_info"`
	FormProduct                   string `json:"form_product"`
	FormProductPlaceholder        string `json:"form_product_placeholder"`
	FormQuantity                  string `json:"form_quantity"`
	FormQuantityInfo              string `json:"form_quantity_info"`
	FormUnitPrice                 string `json:"form_unit_price"`
	FormUnitPriceInfo             string `json:"form_unit_price_info"`
	FormExpenseAccount            string `json:"form_expense_account"`
	FormExpenseAccountPlaceholder string `json:"form_expense_account_placeholder"`
	FormStartDate                 string `json:"form_start_date"`
	FormStartDateHint             string `json:"form_start_date_hint"`
	FormEndDate                   string `json:"form_end_date"`
	FormLineNumber                string `json:"form_line_number"`
}

type LinkedPOLabels struct {
	PONumber     string `json:"po_number"`
	Status       string `json:"status"`
	TotalAmount  string `json:"total_amount"`
	OrderDate    string `json:"order_date"`
	EmptyTitle   string `json:"empty_title"`
	EmptyMessage string `json:"empty_message"`
}

type LinkedExpenditureLabels struct {
	Reference    string `json:"reference"`
	Status       string `json:"status"`
	Amount       string `json:"amount"`
	Date         string `json:"date"`
	EmptyTitle   string `json:"empty_title"`
	EmptyMessage string `json:"empty_message"`
}

// FormLabels holds all form-level labels for the drawer form.
type FormLabels struct {
	// Section headers (5-section parity layout)
	SectionIdentity       string `json:"section_identity"`
	SectionValidity       string `json:"section_validity"`
	SectionMoney          string `json:"section_money"`
	SectionCategorization string `json:"section_categorization"`
	SectionOthers         string `json:"section_others"`

	// §1 Identity
	Name                      string `json:"name"`
	NamePlaceholder           string `json:"name_placeholder"`
	NameInfo                  string `json:"name_info"`
	ContractNumber            string `json:"contract_number"`
	ContractNumberPlaceholder string `json:"contract_number_placeholder"`
	Kind                      string `json:"kind"`
	KindInfo                  string `json:"kind_info"`
	KindSubscription          string `json:"kind_subscription"`
	KindRetainer              string `json:"kind_retainer"`
	KindLease                 string `json:"kind_lease"`
	KindUtility               string `json:"kind_utility"`
	KindFramework             string `json:"kind_framework"`
	KindBlanket               string `json:"kind_blanket"`
	KindOneTime               string `json:"kind_one_time"`
	KindOther                 string `json:"kind_other"`
	Supplier                  string `json:"supplier"`
	SupplierPlaceholder       string `json:"supplier_placeholder"`
	SupplierInfo              string `json:"supplier_info"`

	// §2 Validity & Recurrence
	StartDate             string `json:"start_date"`
	EndDate               string `json:"end_date"`
	EndDateHint           string `json:"end_date_hint"`
	BillingCycleValue     string `json:"billing_cycle_value"`
	BillingCycleUnit      string `json:"billing_cycle_unit"`
	BillingCycleInfo      string `json:"billing_cycle_info"`
	CycleUnitDay          string `json:"cycle_unit_day"`
	CycleUnitWeek         string `json:"cycle_unit_week"`
	CycleUnitMonth        string `json:"cycle_unit_month"`
	CycleUnitYear         string `json:"cycle_unit_year"`
	AutoRenew             string `json:"auto_renew"`
	RenewalNoticeDays     string `json:"renewal_notice_days"`
	RenewalNoticeDaysHint string `json:"renewal_notice_days_hint"`

	// §3 Money & Approval
	Currency               string `json:"currency"`
	CurrencyInfo           string `json:"currency_info"`
	Status                 string `json:"status"`
	StatusInfo             string `json:"status_info"`
	StatusDraft            string `json:"status_draft"`
	StatusRequested        string `json:"status_requested"`
	StatusPendingApproval  string `json:"status_pending_approval"`
	StatusApproved         string `json:"status_approved"`
	StatusActive           string `json:"status_active"`
	StatusExpiring         string `json:"status_expiring"`
	StatusSuspended        string `json:"status_suspended"`
	StatusExpired          string `json:"status_expired"`
	StatusTerminated       string `json:"status_terminated"`
	StatusRejected         string `json:"status_rejected"`
	CommittedAmount        string `json:"committed_amount"`
	CommittedAmountInfo    string `json:"committed_amount_info"`
	CycleAmount            string `json:"cycle_amount"`
	CycleAmountHint        string `json:"cycle_amount_hint"`
	PaymentTerm            string `json:"payment_term"`
	PaymentTermPlaceholder string `json:"payment_term_placeholder"`
	ApprovedBy             string `json:"approved_by"`
	ApprovedDate           string `json:"approved_date"`

	// §4 Categorization
	ExpenditureCategory            string `json:"expenditure_category"`
	ExpenditureCategoryPlaceholder string `json:"expenditure_category_placeholder"`
	ExpenseAccount                 string `json:"expense_account"`
	ExpenseAccountPlaceholder      string `json:"expense_account_placeholder"`
	Location                       string `json:"location"`
	LocationPlaceholder            string `json:"location_placeholder"`

	// §5 Others
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notes_placeholder"`
	Active           string `json:"active"`

	// Action buttons on detail page
	Edit      string `json:"edit"`
	EditTitle string `json:"edit_title"`
	Approve   string `json:"approve"`
	Terminate string `json:"terminate"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultLabels returns English fallback labels.
// Uses proto-generic naming — tier overrides belong in lyngua JSON.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:           "Supplier Contracts",
			HeadingDraft:      "Draft Contracts",
			HeadingActive:     "Active Contracts",
			HeadingExpiring:   "Expiring Contracts",
			HeadingExpired:    "Expired Contracts",
			HeadingTerminated: "Terminated Contracts",
			Caption:           "Standing agreements with suppliers",
			AddButton:         "New Contract",
			DetailSubtitle:    "Contract details",
		},
		Columns: ColumnLabels{
			Name:      "Name",
			Supplier:  "Supplier",
			Kind:      "Kind",
			Status:    "Status",
			Validity:  "Validity",
			Committed: "Committed",
			Released:  "Released",
			Billed:    "Billed",
			Remaining: "Remaining",
		},
		Tabs: TabLabels{
			Info:                "Info",
			Lines:               "Lines",
			LinkedPOs:           "Linked POs",
			LinkedExpenditures:  "Linked Expenditures",
			PriceSchedules:      "Price Schedules",
			Activity:            "Activity",
			ActivityEmpty:       "No activity recorded yet.",
			PriceSchedulesEmpty: "No price schedules yet. Add a schedule to layer multi-year pricing on this contract.",
		},
		Detail: DetailLabels{
			InfoSection:     "Contract Information",
			Name:            "Name",
			Kind:            "Kind",
			Status:          "Status",
			Supplier:        "Supplier",
			StartDate:       "Start Date",
			EndDate:         "End Date",
			AutoRenew:       "Auto Renew",
			Currency:        "Currency",
			CommittedAmount: "Committed Amount",
			ReleasedAmount:  "Released Amount",
			BilledAmount:    "Billed Amount",
			RemainingAmount: "Remaining Amount",
			Notes:           "Notes",
			TabAttachments:  "Attachments",
		},
		Lines: LineLabels{
			Description:                   "Description",
			LineType:                      "Line Type",
			Quantity:                      "Quantity",
			UnitPrice:                     "Unit Price",
			Total:                         "Total",
			Treatment:                     "Treatment",
			EmptyTitle:                    "No lines yet",
			EmptyMessage:                  "Add a line to this contract.",
			AddLine:                       "Add Line",
			TreatmentRecurring:            "Recurring",
			TreatmentOneTime:              "One Time",
			TreatmentUsageBased:           "Usage Based",
			TreatmentMinimumCommitment:    "Minimum Commitment",
			LineTypeGoods:                 "Goods",
			LineTypeService:               "Service",
			LineTypeExpense:               "Expense",
			FormDescription:               "Description",
			FormDescriptionPlaceholder:    "e.g. Cloud hosting — 50 seats",
			FormLineType:                  "Line Type",
			FormLineTypeInfo:              "Goods = physical items; Service = intangible; Expense = direct cost",
			FormTreatment:                 "Treatment",
			FormTreatmentInfo:             "How this line is billed: recurring, one-time, usage-based, or minimum commitment",
			FormProduct:                   "Product",
			FormProductPlaceholder:        "Select a product (optional)",
			FormQuantity:                  "Quantity",
			FormQuantityInfo:              "For recurring lines, this is the per-cycle quantity.",
			FormUnitPrice:                 "Unit Price",
			FormUnitPriceInfo:             "Amount in centavos ÷ 100 for display.",
			FormExpenseAccount:            "Expense Account",
			FormExpenseAccountPlaceholder: "GL account ID",
			FormStartDate:                 "Start Date",
			FormStartDateHint:             "Leave empty to inherit from contract.",
			FormEndDate:                   "End Date",
			FormLineNumber:                "Line Number",
		},
		LinkedPOs: LinkedPOLabels{
			PONumber:     "PO Number",
			Status:       "Status",
			TotalAmount:  "Total Amount",
			OrderDate:    "Order Date",
			EmptyTitle:   "No linked purchase orders",
			EmptyMessage: "POs created against this contract will appear here.",
		},
		LinkedExpenditures: LinkedExpenditureLabels{
			Reference:    "Reference",
			Status:       "Status",
			Amount:       "Amount",
			Date:         "Date",
			EmptyTitle:   "No linked expenditures",
			EmptyMessage: "Expenditures linked to this contract will appear here.",
		},
		Form: FormLabels{
			SectionIdentity:                "Identity Details",
			SectionValidity:                "Validity & Recurrence",
			SectionMoney:                   "Money & Approval",
			SectionCategorization:          "Categorization",
			SectionOthers:                  "Others",
			Name:                           "Contract Name",
			NamePlaceholder:                "e.g. AWS Hosting MSA 2026",
			NameInfo:                       "A short descriptive name for this contract.",
			ContractNumber:                 "Contract Number",
			ContractNumberPlaceholder:      "Supplier's reference number",
			Kind:                           "Kind",
			KindInfo:                       "Subscription = recurring time-based; Blanket = quantity-based commitment; Framework = pricing agreement only.",
			KindSubscription:               "Subscription",
			KindRetainer:                   "Retainer",
			KindLease:                      "Lease",
			KindUtility:                    "Utility",
			KindFramework:                  "Framework",
			KindBlanket:                    "Blanket",
			KindOneTime:                    "One Time",
			KindOther:                      "Other",
			Supplier:                       "Supplier",
			SupplierPlaceholder:            "Select supplier",
			SupplierInfo:                   "The vendor or service provider you are committing to.",
			StartDate:                      "Start Date",
			EndDate:                        "End Date",
			EndDateHint:                    "Leave empty for open-ended.",
			BillingCycleValue:              "Billing Cycle",
			BillingCycleUnit:               "Cycle Unit",
			BillingCycleInfo:               "How often this contract generates an expenditure (for recurring kinds).",
			CycleUnitDay:                   "Day",
			CycleUnitWeek:                  "Week",
			CycleUnitMonth:                 "Month",
			CycleUnitYear:                  "Year",
			AutoRenew:                      "Auto Renew",
			RenewalNoticeDays:              "Renewal Notice (days)",
			RenewalNoticeDaysHint:          "How many days before expiry to send a renewal reminder.",
			Currency:                       "Currency",
			CurrencyInfo:                   "ISO 4217 currency code (e.g. PHP, USD).",
			Status:                         "Status",
			StatusInfo:                     "Lifecycle stage. draft → requested → pending_approval → approved → active → expiring/expired/terminated.",
			StatusDraft:                    "Draft",
			StatusRequested:                "Requested",
			StatusPendingApproval:          "Pending Approval",
			StatusApproved:                 "Approved",
			StatusActive:                   "Active",
			StatusExpiring:                 "Expiring",
			StatusSuspended:                "Suspended",
			StatusExpired:                  "Expired",
			StatusTerminated:               "Terminated",
			StatusRejected:                 "Rejected",
			CommittedAmount:                "Committed Amount",
			CommittedAmountInfo:            "Total value committed at signing (centavos). Immutable after approval.",
			CycleAmount:                    "Cycle Amount",
			CycleAmountHint:                "Expected per-cycle charge for recurring contracts.",
			PaymentTerm:                    "Payment Term",
			PaymentTermPlaceholder:         "Select payment term",
			ApprovedBy:                     "Approved By",
			ApprovedDate:                   "Approved Date",
			ExpenditureCategory:            "Expenditure Category",
			ExpenditureCategoryPlaceholder: "Select category",
			ExpenseAccount:                 "Expense Account",
			ExpenseAccountPlaceholder:      "GL account ID",
			Location:                       "Location",
			LocationPlaceholder:            "Branch or cost center",
			Notes:                          "Notes",
			NotesPlaceholder:               "Additional notes or context",
			Active:                         "Active",
			Edit:                           "Edit",
			EditTitle:                      "Edit Supplier Contract",
			Approve:                        "Approve",
			Terminate:                      "Terminate",
		},
		Empty: EmptyLabels{
			Title:   "No supplier contracts",
			Message: "Create your first supplier contract to start tracking commitments.",
		},
	}
}
