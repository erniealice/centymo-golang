package cost_source_component

// Labels holds every translatable string of the recoverable cost line (cost_source_component) screens.
// Lyngua file: general/cost_source_component.json, root key "cost_source_component" (leasing overlays values only).
// Every field has a general-tier Lyngua key; the Go defaults equal the general tier
// (a reflective test compares every string field, C14).
type Labels struct {
	Page    PageLabels    `json:"page"`
	Tabs    TabsLabels    `json:"tabs"`
	Columns ColumnsLabels `json:"columns"`
	Buttons ButtonsLabels `json:"buttons"`
	Form    FormLabels    `json:"form"`
	Enums   EnumsLabels   `json:"enums"`
	Detail  DetailLabels  `json:"detail"`
	Confirm ConfirmLabels `json:"confirm"`
	Empty   EmptyLabels   `json:"empty"`
	Errors  ErrorsLabels  `json:"errors"`
}

type PageLabels struct {
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle"`
	TitleUnclaimed string `json:"title_unclaimed"`
	TitleClaimed   string `json:"title_claimed"`
}

type TabsLabels struct {
	RecoverableCosts string `json:"recoverable_costs"`
	Info             string `json:"info"`
	Allocation       string `json:"allocation"`
}

type ColumnsLabels struct {
	ComponentKind string `json:"component_kind"`
	Description   string `json:"description"`
	Basis         string `json:"basis"`
	Amount        string `json:"amount"`
	TaxFact       string `json:"tax_fact"`
	ServicePeriod string `json:"service_period"`
	Claim         string `json:"claim"`
	Allocated     string `json:"allocated"`
	Unallocated   string `json:"unallocated"`
}

type ButtonsLabels struct {
	Add            string `json:"add"`
	Edit           string `json:"edit"`
	Delete         string `json:"delete"`
	Allocate       string `json:"allocate"`
	ViewAllocation string `json:"view_allocation"`
}

type FormLabels struct {
	ComponentKindLabel         string `json:"component_kind_label"`
	ComponentKindPlaceholder   string `json:"component_kind_placeholder"`
	DescriptionLabel           string `json:"description_label"`
	DescriptionPlaceholder     string `json:"description_placeholder"`
	BasisUnitLabel             string `json:"basis_unit_label"`
	BasisUnitPlaceholder       string `json:"basis_unit_placeholder"`
	BasisQuantityLabel         string `json:"basis_quantity_label"`
	AmountLabel                string `json:"amount_label"`
	TaxTreatmentLabel          string `json:"tax_treatment_label"`
	TaxTreatmentPlaceholder    string `json:"tax_treatment_placeholder"`
	TaxFactLabel               string `json:"tax_fact_label"`
	ServiceFromLabel           string `json:"service_from_label"`
	ServiceToLabel             string `json:"service_to_label"`
	ExpenditureLineLabel       string `json:"expenditure_line_label"`
	ExpenditureLinePlaceholder string `json:"expenditure_line_placeholder"`
}

type EnumsLabels struct {
	ComponentKindEnergy     string `json:"component_kind_energy"`
	ComponentKindWater      string `json:"component_kind_water"`
	ComponentKindServiceFee string `json:"component_kind_service_fee"`
	ComponentKindOther      string `json:"component_kind_other"`
	TaxFactVatable          string `json:"tax_fact_vatable"`
	TaxFactExempt           string `json:"tax_fact_exempt"`
	TaxFactZeroRated        string `json:"tax_fact_zero_rated"`
	TaxFactNotApplicable    string `json:"tax_fact_not_applicable"`
	ClaimKindAllocation     string `json:"claim_kind_allocation"`
	ClaimKindRecognition    string `json:"claim_kind_recognition"`
	ClaimNone               string `json:"claim_none"`
}

type DetailLabels struct {
	TotalComponents string `json:"total_components"`
	TotalBill       string `json:"total_bill"`
	Difference      string `json:"difference"`
	ClaimedNotice   string `json:"claimed_notice"`
	ServicePeriod   string `json:"service_period"`
}

type ConfirmLabels struct {
	DeleteTitle string `json:"delete_title"`
	DeleteMsg   string `json:"delete_msg"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ErrorsLabels struct {
	NotFound            string `json:"not_found"`
	Validation          string `json:"validation"`
	Claimed             string `json:"claimed"`
	TransactionRequired string `json:"transaction_required"`
	ReferenceInvalid    string `json:"reference_invalid"`
	LockUnavailable     string `json:"lock_unavailable"`
	Generic             string `json:"generic"`
	Unavailable         string `json:"unavailable"`
	FormInvalid         string `json:"form_invalid"`
}

// DefaultLabels returns the English (general tier) defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:          "Recoverable costs",
			Subtitle:       "Cost lines from this bill that can be recovered from customers",
			TitleUnclaimed: "Unclaimed costs",
			TitleClaimed:   "Claimed costs",
		},
		Tabs: TabsLabels{
			RecoverableCosts: "Recoverable costs",
			Info:             "Info",
			Allocation:       "Allocation",
		},
		Columns: ColumnsLabels{
			ComponentKind: "Type",
			Description:   "Description",
			Basis:         "Quantity",
			Amount:        "Amount",
			TaxFact:       "Tax",
			ServicePeriod: "Service period",
			Claim:         "Claimed by",
			Allocated:     "Allocated",
			Unallocated:   "Not allocated",
		},
		Buttons: ButtonsLabels{
			Add:            "Add cost line",
			Edit:           "Edit",
			Delete:         "Delete",
			Allocate:       "Allocate",
			ViewAllocation: "View allocation",
		},
		Form: FormLabels{
			ComponentKindLabel:         "Type",
			ComponentKindPlaceholder:   "Select a type",
			DescriptionLabel:           "Description",
			DescriptionPlaceholder:     "e.g. Common-area electricity",
			BasisUnitLabel:             "Unit",
			BasisUnitPlaceholder:       "e.g. kWh",
			BasisQuantityLabel:         "Quantity",
			AmountLabel:                "Amount",
			TaxTreatmentLabel:          "Tax treatment",
			TaxTreatmentPlaceholder:    "Select a tax treatment",
			TaxFactLabel:               "Tax on the bill",
			ServiceFromLabel:           "Service from",
			ServiceToLabel:             "Service to",
			ExpenditureLineLabel:       "Bill line",
			ExpenditureLinePlaceholder: "Whole bill",
		},
		Enums: EnumsLabels{
			ComponentKindEnergy:     "Energy",
			ComponentKindWater:      "Water",
			ComponentKindServiceFee: "Service fee",
			ComponentKindOther:      "Other",
			TaxFactVatable:          "Taxable",
			TaxFactExempt:           "Exempt",
			TaxFactZeroRated:        "Zero-rated",
			TaxFactNotApplicable:    "Not applicable",
			ClaimKindAllocation:     "Allocation",
			ClaimKindRecognition:    "Expense recognition",
			ClaimNone:               "Not claimed",
		},
		Detail: DetailLabels{
			TotalComponents: "Cost lines total",
			TotalBill:       "Bill total",
			Difference:      "Difference",
			ClaimedNotice:   "Claimed cost lines cannot be changed.",
			ServicePeriod:   "{0} to {1}",
		},
		Confirm: ConfirmLabels{
			DeleteTitle: "Delete this cost line?",
			DeleteMsg:   "This removes the cost line from the bill. It cannot be undone.",
		},
		Empty: EmptyLabels{
			Title:   "No recoverable costs",
			Message: "Add the parts of this bill you plan to recover.",
		},
		Errors: ErrorsLabels{
			NotFound:            "Cost line not found.",
			Validation:          "The cost component is not valid.",
			Claimed:             "This cost component is already claimed by an allocation or a recognition and can no longer be changed.",
			TransactionRequired: "This action needs a database transaction, which is not available.",
			ReferenceInvalid:    "A referenced expenditure, line item or tax treatment was not found in this workspace.",
			LockUnavailable:     "This action needs row locking, which the storage provider does not support.",
			Generic:             "Something went wrong. Please try again.",
			Unavailable:         "This action is not available right now.",
			FormInvalid:         "Please check the form and try again.",
		},
	}
}

// ErrorMessage maps a use-case refusal code (ErrorKind) to its label; unknown codes
// fall back to Errors.Generic.
func (l Labels) ErrorMessage(code string) string {
	switch code {
	case ErrNotFound:
		return l.Errors.NotFound
	case ErrValidation:
		return l.Errors.Validation
	case ErrClaimed:
		return l.Errors.Claimed
	case ErrTransactionRequired:
		return l.Errors.TransactionRequired
	case ErrReferenceInvalid:
		return l.Errors.ReferenceInvalid
	case ErrLockUnavailable:
		return l.Errors.LockUnavailable
	default:
		return l.Errors.Generic
	}
}
