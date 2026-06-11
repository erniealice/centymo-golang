package procurement

// cost_plan_labels.go — extracted verbatim from the root labels.go
// (centymo W7). Pure structural move — no behaviour change.

// ---------------------------------------------------------------------------
// P3 — CostPlan labels
// ---------------------------------------------------------------------------

// CostPlanLabels holds all translatable strings for the cost_plan module.
type CostPlanLabels struct {
	Page    CostPlanPageLabels    `json:"page"`
	Columns CostPlanColumnLabels  `json:"columns"`
	Tabs    CostPlanTabLabels     `json:"tabs"`
	Detail  CostPlanDetailLabels  `json:"detail"`
	Form    CostPlanFormLabels    `json:"form"`
	Actions CostPlanActionLabels  `json:"actions"`
	Confirm CostPlanConfirmLabels `json:"confirm"`
	Buttons CostPlanButtonLabels  `json:"buttons"`
	Bulk    CostPlanBulkLabels    `json:"bulk"`
	Status  CostPlanStatusLabels  `json:"status"`
	Empty   CostPlanEmptyLabels   `json:"empty"`
	Errors  CostPlanErrorLabels   `json:"errors"`
}

type CostPlanPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type CostPlanColumnLabels struct {
	Name         string `json:"name"`
	BillingKind  string `json:"billingKind"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	SupplierPlan string `json:"supplierPlan"`
	CostSchedule string `json:"costSchedule"`
	Active       string `json:"active"`
}

type CostPlanTabLabels struct {
	Info                string `json:"info"`
	Lines               string `json:"lines"`
	LinkedSubscriptions string `json:"linkedSubscriptions"`
	Activity            string `json:"activity"`
}

type CostPlanDetailLabels struct {
	InfoSection  string `json:"infoSection"`
	Name         string `json:"name"`
	BillingKind  string `json:"billingKind"`
	AmountBasis  string `json:"amountBasis"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	BillingCycle string `json:"billingCycle"`
	DefaultTerm  string `json:"defaultTerm"`
	SupplierPlan string `json:"supplierPlan"`
	CostSchedule string `json:"costSchedule"`
	Active       string `json:"active"`
	Inactive     string `json:"inactive"`
}

type CostPlanFormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	Name                    string `json:"name"`
	NamePlaceholder         string `json:"namePlaceholder"`
	Description             string `json:"description"`
	DescPlaceholder         string `json:"descPlaceholder"`
	SupplierPlan            string `json:"supplierPlan"`
	SupplierPlanPlaceholder string `json:"supplierPlanPlaceholder"`
	CostSchedule            string `json:"costSchedule"`
	CostSchedulePlaceholder string `json:"costSchedulePlaceholder"`
	BillingKind             string `json:"billingKind"`
	AmountBasis             string `json:"amountBasis"`
	Amount                  string `json:"amount"`
	AmountPlaceholder       string `json:"amountPlaceholder"`
	Currency                string `json:"currency"`
	CurrencyPlaceholder     string `json:"currencyPlaceholder"`
	BillingCycle            string `json:"billingCycle"`
	BillingCyclePlaceholder string `json:"billingCyclePlaceholder"`
	DefaultTerm             string `json:"defaultTerm"`
	DefaultTermPlaceholder  string `json:"defaultTermPlaceholder"`
	Active                  string `json:"active"`

	// BillingKind option labels
	BillingKindOneTime    string `json:"billingKindOneTime"`
	BillingKindRecurring  string `json:"billingKindRecurring"`
	BillingKindContract   string `json:"billingKindContract"`
	BillingKindUsageBased string `json:"billingKindUsageBased"`
	BillingKindAdHoc      string `json:"billingKindAdHoc"`

	// AmountBasis option labels
	AmountBasisPerCycle         string `json:"amountBasisPerCycle"`
	AmountBasisTotalPackage     string `json:"amountBasisTotalPackage"`
	AmountBasisDerivedFromLines string `json:"amountBasisDerivedFromLines"`
	AmountBasisPerOccurrence    string `json:"amountBasisPerOccurrence"`

	// Duration unit option labels (shared by billing_cycle_unit and default_term_unit)
	DurationUnitDay   string `json:"durationUnitDay"`
	DurationUnitWeek  string `json:"durationUnitWeek"`
	DurationUnitMonth string `json:"durationUnitMonth"`
	DurationUnitYear  string `json:"durationUnitYear"`
}

type CostPlanActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type CostPlanConfirmLabels struct {
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activateMessage"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
}

type CostPlanButtonLabels struct {
	AddCostPlan string `json:"addCostPlan"`
}

type CostPlanBulkLabels struct {
	Delete string `json:"delete"`
}

type CostPlanStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type CostPlanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type CostPlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultCostPlanLabels returns English fallback labels.
func DefaultCostPlanLabels() CostPlanLabels {
	return CostPlanLabels{
		Page: CostPlanPageLabels{
			Heading:         "Cost Plans",
			HeadingActive:   "Active Cost Plans",
			HeadingInactive: "Inactive Cost Plans",
			Caption:         "Supplier pricing plans and billing schedules",
			CaptionActive:   "Active cost plans",
			CaptionInactive: "Inactive cost plans",
			PageTitle:       "Cost Plan",
		},
		Columns: CostPlanColumnLabels{
			Name:         "Name",
			BillingKind:  "Billing Kind",
			Amount:       "Amount",
			Currency:     "Currency",
			SupplierPlan: "Supplier Plan",
			CostSchedule: "Cost Schedule",
			Active:       "Status",
		},
		Tabs: CostPlanTabLabels{
			Info:                "Info",
			Lines:               "Lines",
			LinkedSubscriptions: "Subscriptions",
			Activity:            "Activity",
		},
		Detail: CostPlanDetailLabels{
			InfoSection:  "Cost Plan Details",
			Name:         "Name",
			BillingKind:  "Billing Kind",
			AmountBasis:  "Amount Basis",
			Amount:       "Amount",
			Currency:     "Currency",
			BillingCycle: "Billing Cycle",
			DefaultTerm:  "Default Term",
			SupplierPlan: "Supplier Plan",
			CostSchedule: "Cost Schedule",
			Active:       "Active",
			Inactive:     "Inactive",
		},
		Form: CostPlanFormLabels{
			SectionIdentification:       "Identification",
			SectionRelationships:        "Relationships",
			SectionConfiguration:        "Configuration",
			SectionSchedule:             "Schedule",
			SectionNotes:                "Notes",
			Name:                        "Name",
			NamePlaceholder:             "e.g. AWS EC2 Monthly",
			Description:                 "Description",
			DescPlaceholder:             "Internal notes about this cost plan",
			SupplierPlan:                "Supplier Plan",
			SupplierPlanPlaceholder:     "Select supplier plan",
			CostSchedule:                "Cost Schedule",
			CostSchedulePlaceholder:     "Select cost schedule",
			BillingKind:                 "Billing Kind",
			AmountBasis:                 "Amount Basis",
			Amount:                      "Amount",
			AmountPlaceholder:           "0.00",
			Currency:                    "Currency",
			CurrencyPlaceholder:         "e.g. PHP",
			BillingCycle:                "Billing Cycle",
			BillingCyclePlaceholder:     "e.g. 1",
			DefaultTerm:                 "Default Term",
			DefaultTermPlaceholder:      "e.g. 12",
			Active:                      "Active",
			BillingKindOneTime:          "One Time",
			BillingKindRecurring:        "Recurring",
			BillingKindContract:         "Contract",
			BillingKindUsageBased:       "Usage Based",
			BillingKindAdHoc:            "Ad Hoc",
			AmountBasisPerCycle:         "Per Cycle",
			AmountBasisTotalPackage:     "Total Package",
			AmountBasisDerivedFromLines: "Derived From Lines",
			AmountBasisPerOccurrence:    "Per Occurrence",
			DurationUnitDay:             "Day",
			DurationUnitWeek:            "Week",
			DurationUnitMonth:           "Month",
			DurationUnitYear:            "Year",
		},
		Actions: CostPlanActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: CostPlanConfirmLabels{
			Delete:                "Delete Cost Plan",
			DeleteMessage:         "Are you sure you want to delete this cost plan?",
			Activate:              "Activate Cost Plan",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Cost Plan",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Cost Plans",
			BulkDeleteMessage:     "Delete selected cost plans?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected cost plans?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected cost plans?",
		},
		Buttons: CostPlanButtonLabels{AddCostPlan: "Add Cost Plan"},
		Bulk:    CostPlanBulkLabels{Delete: "Delete"},
		Status: CostPlanStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: CostPlanEmptyLabels{
			Title:   "No cost plans yet",
			Message: "Add a cost plan to define billing terms for a supplier engagement.",
		},
		Errors: CostPlanErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Cost plan not found.",
			IDRequired:       "Cost plan ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This cost plan has linked subscriptions and cannot be deleted.",
			LoadFailed:       "Failed to load cost plan.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}
