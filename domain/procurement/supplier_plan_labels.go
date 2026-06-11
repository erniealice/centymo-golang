package procurement

// supplier_plan_labels.go — extracted verbatim from the root labels.go
// (centymo W7). Pure structural move — no behaviour change.

// ---------------------------------------------------------------------------
// P3 — SupplierPlan labels
// ---------------------------------------------------------------------------

// SupplierPlanLabels holds all translatable strings for the supplier_plan module.
type SupplierPlanLabels struct {
	Page    SupplierPlanPageLabels    `json:"page"`
	Columns SupplierPlanColumnLabels  `json:"columns"`
	Tabs    SupplierPlanTabLabels     `json:"tabs"`
	Detail  SupplierPlanDetailLabels  `json:"detail"`
	Form    SupplierPlanFormLabels    `json:"form"`
	Actions SupplierPlanActionLabels  `json:"actions"`
	Confirm SupplierPlanConfirmLabels `json:"confirm"`
	Buttons SupplierPlanButtonLabels  `json:"buttons"`
	Bulk    SupplierPlanBulkLabels    `json:"bulk"`
	Status  SupplierPlanStatusLabels  `json:"status"`
	Empty   SupplierPlanEmptyLabels   `json:"empty"`
	Errors  SupplierPlanErrorLabels   `json:"errors"`
}

type SupplierPlanPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type SupplierPlanColumnLabels struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Supplier string `json:"supplier"`
	Active   string `json:"active"`
}

type SupplierPlanTabLabels struct {
	Info         string `json:"info"`
	CostPlans    string `json:"costPlans"`
	ProductPlans string `json:"productPlans"`
	Activity     string `json:"activity"`
}

type SupplierPlanDetailLabels struct {
	InfoSection string `json:"infoSection"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Supplier    string `json:"supplier"`
	Active      string `json:"active"`
	Inactive    string `json:"inactive"`
}

type SupplierPlanFormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	Name                string `json:"name"`
	NamePlaceholder     string `json:"namePlaceholder"`
	Code                string `json:"code"`
	CodePlaceholder     string `json:"codePlaceholder"`
	Supplier            string `json:"supplier"`
	SupplierPlaceholder string `json:"supplierPlaceholder"`
	Active              string `json:"active"`
}

type SupplierPlanActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type SupplierPlanConfirmLabels struct {
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

type SupplierPlanButtonLabels struct {
	AddSupplierPlan string `json:"addSupplierPlan"`
}

type SupplierPlanBulkLabels struct {
	Delete string `json:"delete"`
}

type SupplierPlanStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type SupplierPlanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SupplierPlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultSupplierPlanLabels returns English fallback labels.
func DefaultSupplierPlanLabels() SupplierPlanLabels {
	return SupplierPlanLabels{
		Page: SupplierPlanPageLabels{
			Heading:         "Supplier Plans",
			HeadingActive:   "Active Supplier Plans",
			HeadingInactive: "Inactive Supplier Plans",
			Caption:         "Supplier product and pricing plans",
			CaptionActive:   "Active supplier plans",
			CaptionInactive: "Inactive supplier plans",
			PageTitle:       "Supplier Plan",
		},
		Columns: SupplierPlanColumnLabels{
			Name:     "Name",
			Code:     "Code",
			Supplier: "Supplier",
			Active:   "Status",
		},
		Tabs: SupplierPlanTabLabels{
			Info:         "Info",
			CostPlans:    "Cost Plans",
			ProductPlans: "Product Plans",
			Activity:     "Activity",
		},
		Detail: SupplierPlanDetailLabels{
			InfoSection: "Plan Details",
			Name:        "Name",
			Code:        "Code",
			Supplier:    "Supplier",
			Active:      "Active",
			Inactive:    "Inactive",
		},
		Form: SupplierPlanFormLabels{
			SectionIdentification: "Identification",
			SectionRelationships:  "Relationships",
			SectionConfiguration:  "Configuration",
			SectionSchedule:       "Schedule",
			SectionNotes:          "Notes",
			Name:                  "Name",
			NamePlaceholder:       "e.g. AWS Standard Plan",
			Code:                  "Code",
			CodePlaceholder:       "e.g. PLAN-AWS-001",
			Supplier:              "Supplier",
			SupplierPlaceholder:   "Select supplier",
			Active:                "Active",
		},
		Actions: SupplierPlanActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: SupplierPlanConfirmLabels{
			Delete:                "Delete Supplier Plan",
			DeleteMessage:         "Are you sure you want to delete this supplier plan?",
			Activate:              "Activate Supplier Plan",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Supplier Plan",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Supplier Plans",
			BulkDeleteMessage:     "Delete selected supplier plans?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected supplier plans?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected supplier plans?",
		},
		Buttons: SupplierPlanButtonLabels{AddSupplierPlan: "Add Supplier Plan"},
		Bulk:    SupplierPlanBulkLabels{Delete: "Delete"},
		Status: SupplierPlanStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: SupplierPlanEmptyLabels{
			Title:   "No supplier plans yet",
			Message: "Add a supplier plan to group cost plans and product plans for a vendor.",
		},
		Errors: SupplierPlanErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Supplier plan not found.",
			IDRequired:       "Supplier plan ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This supplier plan has linked cost plans or product plans and cannot be deleted.",
			LoadFailed:       "Failed to load supplier plan.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}
