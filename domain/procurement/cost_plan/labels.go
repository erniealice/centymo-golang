package cost_plan

// cost_plan_labels.go — extracted verbatim from the root labels.go
// (centymo W7). Pure structural move — no behaviour change.

// ---------------------------------------------------------------------------
// P3 — CostPlan labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the cost_plan module.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Columns ColumnLabels  `json:"columns"`
	Tabs    TabLabels     `json:"tabs"`
	Detail  DetailLabels  `json:"detail"`
	Form    FormLabels    `json:"form"`
	Actions ActionLabels  `json:"actions"`
	Confirm ConfirmLabels `json:"confirm"`
	Buttons ButtonLabels  `json:"buttons"`
	Bulk    BulkLabels    `json:"bulk"`
	Status  StatusLabels  `json:"status"`
	Empty   EmptyLabels   `json:"empty"`
	Errors  ErrorLabels   `json:"errors"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"heading_active"`
	HeadingInactive string `json:"heading_inactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"caption_active"`
	CaptionInactive string `json:"caption_inactive"`
	PageTitle       string `json:"page_title"`
}

type ColumnLabels struct {
	Name         string `json:"name"`
	BillingKind  string `json:"billing_kind"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	SupplierPlan string `json:"supplier_plan"`
	CostSchedule string `json:"cost_schedule"`
	Active       string `json:"active"`
}

type TabLabels struct {
	Info                string `json:"info"`
	Lines               string `json:"lines"`
	LinkedSubscriptions string `json:"linked_subscriptions"`
	Activity            string `json:"activity"`
}

type DetailLabels struct {
	InfoSection  string `json:"info_section"`
	Name         string `json:"name"`
	BillingKind  string `json:"billing_kind"`
	AmountBasis  string `json:"amount_basis"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	BillingCycle string `json:"billing_cycle"`
	DefaultTerm  string `json:"default_term"`
	SupplierPlan string `json:"supplier_plan"`
	CostSchedule string `json:"cost_schedule"`
	Active       string `json:"active"`
	Inactive     string `json:"inactive"`
}

type FormLabels struct {
	SectionIdentification string `json:"section_identification"`
	SectionRelationships  string `json:"section_relationships"`
	SectionConfiguration  string `json:"section_configuration"`
	SectionSchedule       string `json:"section_schedule"`
	SectionNotes          string `json:"section_notes"`

	Name                    string `json:"name"`
	NamePlaceholder         string `json:"name_placeholder"`
	Description             string `json:"description"`
	DescPlaceholder         string `json:"desc_placeholder"`
	SupplierPlan            string `json:"supplier_plan"`
	SupplierPlanPlaceholder string `json:"supplier_plan_placeholder"`
	CostSchedule            string `json:"cost_schedule"`
	CostSchedulePlaceholder string `json:"cost_schedule_placeholder"`
	BillingKind             string `json:"billing_kind"`
	AmountBasis             string `json:"amount_basis"`
	Amount                  string `json:"amount"`
	AmountPlaceholder       string `json:"amount_placeholder"`
	Currency                string `json:"currency"`
	CurrencyPlaceholder     string `json:"currency_placeholder"`
	BillingCycle            string `json:"billing_cycle"`
	BillingCyclePlaceholder string `json:"billing_cycle_placeholder"`
	DefaultTerm             string `json:"default_term"`
	DefaultTermPlaceholder  string `json:"default_term_placeholder"`
	Active                  string `json:"active"`

	// BillingKind option labels
	BillingKindOneTime    string `json:"billing_kind_one_time"`
	BillingKindRecurring  string `json:"billing_kind_recurring"`
	BillingKindContract   string `json:"billing_kind_contract"`
	BillingKindUsageBased string `json:"billing_kind_usage_based"`
	BillingKindAdHoc      string `json:"billing_kind_ad_hoc"`

	// AmountBasis option labels
	AmountBasisPerCycle         string `json:"amount_basis_per_cycle"`
	AmountBasisTotalPackage     string `json:"amount_basis_total_package"`
	AmountBasisDerivedFromLines string `json:"amount_basis_derived_from_lines"`
	AmountBasisPerOccurrence    string `json:"amount_basis_per_occurrence"`

	// Duration unit option labels (shared by billing_cycle_unit and default_term_unit)
	DurationUnitDay   string `json:"duration_unit_day"`
	DurationUnitWeek  string `json:"duration_unit_week"`
	DurationUnitMonth string `json:"duration_unit_month"`
	DurationUnitYear  string `json:"duration_unit_year"`
}

type ActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"no_permission"`
}

type ConfirmLabels struct {
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"delete_message"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activate_message"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivate_message"`
	BulkDelete            string `json:"bulk_delete"`
	BulkDeleteMessage     string `json:"bulk_delete_message"`
	BulkActivate          string `json:"bulk_activate"`
	BulkActivateMessage   string `json:"bulk_activate_message"`
	BulkDeactivate        string `json:"bulk_deactivate"`
	BulkDeactivateMessage string `json:"bulk_deactivate_message"`
}

type ButtonLabels struct {
	AddCostPlan string `json:"add_cost_plan"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

type StatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	NoPermission     string `json:"no_permission"`
	InUse            string `json:"in_use"`
	LoadFailed       string `json:"load_failed"`
	NoIDsProvided    string `json:"no_ids_provided"`
}

// DefaultLabels returns English fallback labels.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:         "Cost Plans",
			HeadingActive:   "Active Cost Plans",
			HeadingInactive: "Inactive Cost Plans",
			Caption:         "Supplier pricing plans and billing schedules",
			CaptionActive:   "Active cost plans",
			CaptionInactive: "Inactive cost plans",
			PageTitle:       "Cost Plan",
		},
		Columns: ColumnLabels{
			Name:         "Name",
			BillingKind:  "Billing Kind",
			Amount:       "Amount",
			Currency:     "Currency",
			SupplierPlan: "Supplier Plan",
			CostSchedule: "Cost Schedule",
			Active:       "Status",
		},
		Tabs: TabLabels{
			Info:                "Info",
			Lines:               "Lines",
			LinkedSubscriptions: "Subscriptions",
			Activity:            "Activity",
		},
		Detail: DetailLabels{
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
		Form: FormLabels{
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
		Actions: ActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: ConfirmLabels{
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
		Buttons: ButtonLabels{AddCostPlan: "Add Cost Plan"},
		Bulk:    BulkLabels{Delete: "Delete"},
		Status: StatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: EmptyLabels{
			Title:   "No cost plans yet",
			Message: "Add a cost plan to define billing terms for a supplier engagement.",
		},
		Errors: ErrorLabels{
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
