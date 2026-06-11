package supplier_product_cost_plan

// ---------------------------------------------------------------------------
// P3 — SupplierProductCostPlan labels (inline editor, no full module)
// ---------------------------------------------------------------------------

// Labels holds translatable strings for the inline cost plan line editor.
type Labels struct {
	Form    FormLabels   `json:"form"`
	Columns ColumnLabels `json:"columns"`
	Empty   EmptyLabels  `json:"empty"`
	Actions ActionLabels `json:"actions"`
	Errors  ErrorLabels  `json:"errors"`
}

type FormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	SupplierProductPlan            string `json:"supplierProductPlan"`
	SupplierProductPlanPlaceholder string `json:"supplierProductPlanPlaceholder"`
	BillingTreatment               string `json:"billingTreatment"`
	Amount                         string `json:"amount"`
	AmountPlaceholder              string `json:"amountPlaceholder"`
	MinimumCommitment              string `json:"minimumCommitment"`
	MinimumCommitmentPlaceholder   string `json:"minimumCommitmentPlaceholder"`
	Active                         string `json:"active"`

	// BillingTreatment option labels
	TreatmentRecurring         string `json:"treatmentRecurring"`
	TreatmentOneTimeInitial    string `json:"treatmentOneTimeInitial"`
	TreatmentUsageBased        string `json:"treatmentUsageBased"`
	TreatmentMinimumCommitment string `json:"treatmentMinimumCommitment"`
}

type ColumnLabels struct {
	SupplierProductPlan string `json:"supplierProductPlan"`
	BillingTreatment    string `json:"billingTreatment"`
	Amount              string `json:"amount"`
	Active              string `json:"active"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	AddLine string `json:"addLine"`
}

type ActionLabels struct {
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Add          string `json:"add"`
	NoPermission string `json:"noPermission"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
}

// DefaultLabels returns English fallback labels.
func DefaultLabels() Labels {
	return Labels{
		Form: FormLabels{
			SectionIdentification:          "Identification",
			SectionRelationships:           "Relationships",
			SectionConfiguration:           "Configuration",
			SectionSchedule:                "Schedule",
			SectionNotes:                   "Notes",
			SupplierProductPlan:            "Supplier Product Plan",
			SupplierProductPlanPlaceholder: "Select product plan",
			BillingTreatment:               "Billing Treatment",
			Amount:                         "Amount",
			AmountPlaceholder:              "0.00",
			MinimumCommitment:              "Minimum Commitment",
			MinimumCommitmentPlaceholder:   "0.00",
			Active:                         "Active",
			TreatmentRecurring:             "Recurring",
			TreatmentOneTimeInitial:        "One-Time Initial",
			TreatmentUsageBased:            "Usage Based",
			TreatmentMinimumCommitment:     "Minimum Commitment",
		},
		Columns: ColumnLabels{
			SupplierProductPlan: "Product Plan",
			BillingTreatment:    "Treatment",
			Amount:              "Amount",
			Active:              "Status",
		},
		Empty: EmptyLabels{
			Title:   "No cost plan lines yet",
			Message: "Add product-level cost lines to this cost plan.",
			AddLine: "Add Line",
		},
		Actions: ActionLabels{
			Edit:         "Edit",
			Delete:       "Delete",
			Add:          "Add Line",
			NoPermission: "No permission",
		},
		Errors: ErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Cost plan line not found.",
			IDRequired:       "Cost plan line ID is required.",
		},
	}
}
