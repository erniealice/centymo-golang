package procurement

// ---------------------------------------------------------------------------
// P3 — SupplierProductCostPlan labels (inline editor, no full module)
// ---------------------------------------------------------------------------

// SupplierProductCostPlanLabels holds translatable strings for the inline cost plan line editor.
type SupplierProductCostPlanLabels struct {
	Form    SupplierProductCostPlanFormLabels   `json:"form"`
	Columns SupplierProductCostPlanColumnLabels `json:"columns"`
	Empty   SupplierProductCostPlanEmptyLabels  `json:"empty"`
	Actions SupplierProductCostPlanActionLabels `json:"actions"`
	Errors  SupplierProductCostPlanErrorLabels  `json:"errors"`
}

type SupplierProductCostPlanFormLabels struct {
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

type SupplierProductCostPlanColumnLabels struct {
	SupplierProductPlan string `json:"supplierProductPlan"`
	BillingTreatment    string `json:"billingTreatment"`
	Amount              string `json:"amount"`
	Active              string `json:"active"`
}

type SupplierProductCostPlanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	AddLine string `json:"addLine"`
}

type SupplierProductCostPlanActionLabels struct {
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Add          string `json:"add"`
	NoPermission string `json:"noPermission"`
}

type SupplierProductCostPlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
}

// DefaultSupplierProductCostPlanLabels returns English fallback labels.
func DefaultSupplierProductCostPlanLabels() SupplierProductCostPlanLabels {
	return SupplierProductCostPlanLabels{
		Form: SupplierProductCostPlanFormLabels{
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
		Columns: SupplierProductCostPlanColumnLabels{
			SupplierProductPlan: "Product Plan",
			BillingTreatment:    "Treatment",
			Amount:              "Amount",
			Active:              "Status",
		},
		Empty: SupplierProductCostPlanEmptyLabels{
			Title:   "No cost plan lines yet",
			Message: "Add product-level cost lines to this cost plan.",
			AddLine: "Add Line",
		},
		Actions: SupplierProductCostPlanActionLabels{
			Edit:         "Edit",
			Delete:       "Delete",
			Add:          "Add Line",
			NoPermission: "No permission",
		},
		Errors: SupplierProductCostPlanErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Cost plan line not found.",
			IDRequired:       "Cost plan line ID is required.",
		},
	}
}
