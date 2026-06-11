package expenditure

// ---------------------------------------------------------------------------
// SupplierContract labels  (P3a)
// ---------------------------------------------------------------------------

// SupplierContractLabels holds all translatable strings for the supplier_contract module.
type SupplierContractLabels struct {
	Page               SupplierContractPageLabels              `json:"page"`
	Columns            SupplierContractColumnLabels            `json:"columns"`
	Tabs               SupplierContractTabLabels               `json:"tabs"`
	Detail             SupplierContractDetailLabels            `json:"detail"`
	Lines              SupplierContractLineLabels              `json:"lines"`
	LinkedPOs          SupplierContractLinkedPOLabels          `json:"linkedPos"`
	LinkedExpenditures SupplierContractLinkedExpenditureLabels `json:"linkedExpenditures"`
	Form               SupplierContractFormLabels              `json:"form"`
	Empty              SupplierContractEmptyLabels             `json:"empty"`
}

type SupplierContractPageLabels struct {
	Heading           string `json:"heading"`
	HeadingDraft      string `json:"headingDraft"`
	HeadingActive     string `json:"headingActive"`
	HeadingExpiring   string `json:"headingExpiring"`
	HeadingExpired    string `json:"headingExpired"`
	HeadingTerminated string `json:"headingTerminated"`
	Caption           string `json:"caption"`
	AddButton         string `json:"addButton"`
	DetailSubtitle    string `json:"detailSubtitle"`
}

type SupplierContractColumnLabels struct {
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

type SupplierContractTabLabels struct {
	Info                string `json:"info"`
	Lines               string `json:"lines"`
	LinkedPOs           string `json:"linkedPos"`
	LinkedExpenditures  string `json:"linkedExpenditures"`
	PriceSchedules      string `json:"priceSchedules"`
	Activity            string `json:"activity"`
	ActivityEmpty       string `json:"activityEmpty"`
	PriceSchedulesEmpty string `json:"priceSchedulesEmpty"`
}

type SupplierContractDetailLabels struct {
	InfoSection     string `json:"infoSection"`
	Name            string `json:"name"`
	Kind            string `json:"kind"`
	Status          string `json:"status"`
	Supplier        string `json:"supplier"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	AutoRenew       string `json:"autoRenew"`
	Currency        string `json:"currency"`
	CommittedAmount string `json:"committedAmount"`
	ReleasedAmount  string `json:"releasedAmount"`
	BilledAmount    string `json:"billedAmount"`
	RemainingAmount string `json:"remainingAmount"`
	Notes           string `json:"notes"`
	TabAttachments  string `json:"tabAttachments"`
}

type SupplierContractLineLabels struct {
	// Column labels
	Description  string `json:"description"`
	LineType     string `json:"lineType"`
	Quantity     string `json:"quantity"`
	UnitPrice    string `json:"unitPrice"`
	Total        string `json:"total"`
	Treatment    string `json:"treatment"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
	AddLine      string `json:"addLine"`

	// Enum label values for treatment
	TreatmentRecurring         string `json:"treatmentRecurring"`
	TreatmentOneTime           string `json:"treatmentOneTime"`
	TreatmentUsageBased        string `json:"treatmentUsageBased"`
	TreatmentMinimumCommitment string `json:"treatmentMinimumCommitment"`

	// Enum label values for line_type
	LineTypeGoods   string `json:"lineTypeGoods"`
	LineTypeService string `json:"lineTypeService"`
	LineTypeExpense string `json:"lineTypeExpense"`

	// Drawer form labels
	FormDescription               string `json:"formDescription"`
	FormDescriptionPlaceholder    string `json:"formDescriptionPlaceholder"`
	FormLineType                  string `json:"formLineType"`
	FormLineTypeInfo              string `json:"formLineTypeInfo"`
	FormTreatment                 string `json:"formTreatment"`
	FormTreatmentInfo             string `json:"formTreatmentInfo"`
	FormProduct                   string `json:"formProduct"`
	FormProductPlaceholder        string `json:"formProductPlaceholder"`
	FormQuantity                  string `json:"formQuantity"`
	FormQuantityInfo              string `json:"formQuantityInfo"`
	FormUnitPrice                 string `json:"formUnitPrice"`
	FormUnitPriceInfo             string `json:"formUnitPriceInfo"`
	FormExpenseAccount            string `json:"formExpenseAccount"`
	FormExpenseAccountPlaceholder string `json:"formExpenseAccountPlaceholder"`
	FormStartDate                 string `json:"formStartDate"`
	FormStartDateHint             string `json:"formStartDateHint"`
	FormEndDate                   string `json:"formEndDate"`
	FormLineNumber                string `json:"formLineNumber"`
}

type SupplierContractLinkedPOLabels struct {
	PONumber     string `json:"poNumber"`
	Status       string `json:"status"`
	TotalAmount  string `json:"totalAmount"`
	OrderDate    string `json:"orderDate"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

type SupplierContractLinkedExpenditureLabels struct {
	Reference    string `json:"reference"`
	Status       string `json:"status"`
	Amount       string `json:"amount"`
	Date         string `json:"date"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// SupplierContractFormLabels holds all form-level labels for the drawer form.
type SupplierContractFormLabels struct {
	// Section headers (5-section parity layout)
	SectionIdentity       string `json:"sectionIdentity"`
	SectionValidity       string `json:"sectionValidity"`
	SectionMoney          string `json:"sectionMoney"`
	SectionCategorization string `json:"sectionCategorization"`
	SectionOthers         string `json:"sectionOthers"`

	// §1 Identity
	Name                      string `json:"name"`
	NamePlaceholder           string `json:"namePlaceholder"`
	NameInfo                  string `json:"nameInfo"`
	ContractNumber            string `json:"contractNumber"`
	ContractNumberPlaceholder string `json:"contractNumberPlaceholder"`
	Kind                      string `json:"kind"`
	KindInfo                  string `json:"kindInfo"`
	KindSubscription          string `json:"kindSubscription"`
	KindRetainer              string `json:"kindRetainer"`
	KindLease                 string `json:"kindLease"`
	KindUtility               string `json:"kindUtility"`
	KindFramework             string `json:"kindFramework"`
	KindBlanket               string `json:"kindBlanket"`
	KindOneTime               string `json:"kindOneTime"`
	KindOther                 string `json:"kindOther"`
	Supplier                  string `json:"supplier"`
	SupplierPlaceholder       string `json:"supplierPlaceholder"`
	SupplierInfo              string `json:"supplierInfo"`

	// §2 Validity & Recurrence
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	EndDateHint           string `json:"endDateHint"`
	BillingCycleValue     string `json:"billingCycleValue"`
	BillingCycleUnit      string `json:"billingCycleUnit"`
	BillingCycleInfo      string `json:"billingCycleInfo"`
	CycleUnitDay          string `json:"cycleUnitDay"`
	CycleUnitWeek         string `json:"cycleUnitWeek"`
	CycleUnitMonth        string `json:"cycleUnitMonth"`
	CycleUnitYear         string `json:"cycleUnitYear"`
	AutoRenew             string `json:"autoRenew"`
	RenewalNoticeDays     string `json:"renewalNoticeDays"`
	RenewalNoticeDaysHint string `json:"renewalNoticeDaysHint"`

	// §3 Money & Approval
	Currency               string `json:"currency"`
	CurrencyInfo           string `json:"currencyInfo"`
	Status                 string `json:"status"`
	StatusInfo             string `json:"statusInfo"`
	StatusDraft            string `json:"statusDraft"`
	StatusRequested        string `json:"statusRequested"`
	StatusPendingApproval  string `json:"statusPendingApproval"`
	StatusApproved         string `json:"statusApproved"`
	StatusActive           string `json:"statusActive"`
	StatusExpiring         string `json:"statusExpiring"`
	StatusSuspended        string `json:"statusSuspended"`
	StatusExpired          string `json:"statusExpired"`
	StatusTerminated       string `json:"statusTerminated"`
	StatusRejected         string `json:"statusRejected"`
	CommittedAmount        string `json:"committedAmount"`
	CommittedAmountInfo    string `json:"committedAmountInfo"`
	CycleAmount            string `json:"cycleAmount"`
	CycleAmountHint        string `json:"cycleAmountHint"`
	PaymentTerm            string `json:"paymentTerm"`
	PaymentTermPlaceholder string `json:"paymentTermPlaceholder"`
	ApprovedBy             string `json:"approvedBy"`
	ApprovedDate           string `json:"approvedDate"`

	// §4 Categorization
	ExpenditureCategory            string `json:"expenditureCategory"`
	ExpenditureCategoryPlaceholder string `json:"expenditureCategoryPlaceholder"`
	ExpenseAccount                 string `json:"expenseAccount"`
	ExpenseAccountPlaceholder      string `json:"expenseAccountPlaceholder"`
	Location                       string `json:"location"`
	LocationPlaceholder            string `json:"locationPlaceholder"`

	// §5 Others
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	Active           string `json:"active"`

	// Action buttons on detail page
	Edit      string `json:"edit"`
	EditTitle string `json:"editTitle"`
	Approve   string `json:"approve"`
	Terminate string `json:"terminate"`
}

type SupplierContractEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultSupplierContractLabels returns English fallback labels.
// Uses proto-generic naming — tier overrides belong in lyngua JSON.
func DefaultSupplierContractLabels() SupplierContractLabels {
	return SupplierContractLabels{
		Page: SupplierContractPageLabels{
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
		Columns: SupplierContractColumnLabels{
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
		Tabs: SupplierContractTabLabels{
			Info:                "Info",
			Lines:               "Lines",
			LinkedPOs:           "Linked POs",
			LinkedExpenditures:  "Linked Expenditures",
			PriceSchedules:      "Price Schedules",
			Activity:            "Activity",
			ActivityEmpty:       "No activity recorded yet.",
			PriceSchedulesEmpty: "No price schedules yet. Add a schedule to layer multi-year pricing on this contract.",
		},
		Detail: SupplierContractDetailLabels{
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
		Lines: SupplierContractLineLabels{
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
		LinkedPOs: SupplierContractLinkedPOLabels{
			PONumber:     "PO Number",
			Status:       "Status",
			TotalAmount:  "Total Amount",
			OrderDate:    "Order Date",
			EmptyTitle:   "No linked purchase orders",
			EmptyMessage: "POs created against this contract will appear here.",
		},
		LinkedExpenditures: SupplierContractLinkedExpenditureLabels{
			Reference:    "Reference",
			Status:       "Status",
			Amount:       "Amount",
			Date:         "Date",
			EmptyTitle:   "No linked expenditures",
			EmptyMessage: "Expenditures linked to this contract will appear here.",
		},
		Form: SupplierContractFormLabels{
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
		Empty: SupplierContractEmptyLabels{
			Title:   "No supplier contracts",
			Message: "Create your first supplier contract to start tracking commitments.",
		},
	}
}
