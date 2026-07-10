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
	SectionIdentification string `json:"section_identification"`
	SectionRelationships  string `json:"section_relationships"`
	SectionConfiguration  string `json:"section_configuration"`
	SectionSchedule       string `json:"section_schedule"`
	SectionNotes          string `json:"section_notes"`

	SupplierProductPlan            string `json:"supplier_product_plan"`
	SupplierProductPlanPlaceholder string `json:"supplier_product_plan_placeholder"`
	BillingTreatment               string `json:"billing_treatment"`
	Amount                         string `json:"amount"`
	AmountPlaceholder              string `json:"amount_placeholder"`
	MinimumCommitment              string `json:"minimum_commitment"`
	MinimumCommitmentPlaceholder   string `json:"minimum_commitment_placeholder"`
	Active                         string `json:"active"`

	// BillingTreatment option labels
	TreatmentRecurring         string `json:"treatment_recurring"`
	TreatmentOneTimeInitial    string `json:"treatment_one_time_initial"`
	TreatmentUsageBased        string `json:"treatment_usage_based"`
	TreatmentMinimumCommitment string `json:"treatment_minimum_commitment"`
}

type ColumnLabels struct {
	SupplierProductPlan string `json:"supplier_product_plan"`
	BillingTreatment    string `json:"billing_treatment"`
	Amount              string `json:"amount"`
	Active              string `json:"active"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	AddLine string `json:"add_line"`
}

type ActionLabels struct {
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Add          string `json:"add"`
	NoPermission string `json:"no_permission"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
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
