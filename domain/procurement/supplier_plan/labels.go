package supplier_plan

// supplier_plan_labels.go — extracted verbatim from the root labels.go
// (centymo W7). Pure structural move — no behaviour change.

// ---------------------------------------------------------------------------
// P3 — SupplierPlan labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the supplier_plan module.
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
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type ColumnLabels struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Supplier string `json:"supplier"`
	Active   string `json:"active"`
}

type TabLabels struct {
	Info         string `json:"info"`
	CostPlans    string `json:"costPlans"`
	ProductPlans string `json:"productPlans"`
	Activity     string `json:"activity"`
}

type DetailLabels struct {
	InfoSection string `json:"infoSection"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Supplier    string `json:"supplier"`
	Active      string `json:"active"`
	Inactive    string `json:"inactive"`
}

type FormLabels struct {
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

type ActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type ConfirmLabels struct {
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

type ButtonLabels struct {
	AddSupplierPlan string `json:"addSupplierPlan"`
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
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultLabels returns English fallback labels.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:         "Supplier Plans",
			HeadingActive:   "Active Supplier Plans",
			HeadingInactive: "Inactive Supplier Plans",
			Caption:         "Supplier product and pricing plans",
			CaptionActive:   "Active supplier plans",
			CaptionInactive: "Inactive supplier plans",
			PageTitle:       "Supplier Plan",
		},
		Columns: ColumnLabels{
			Name:     "Name",
			Code:     "Code",
			Supplier: "Supplier",
			Active:   "Status",
		},
		Tabs: TabLabels{
			Info:         "Info",
			CostPlans:    "Cost Plans",
			ProductPlans: "Product Plans",
			Activity:     "Activity",
		},
		Detail: DetailLabels{
			InfoSection: "Plan Details",
			Name:        "Name",
			Code:        "Code",
			Supplier:    "Supplier",
			Active:      "Active",
			Inactive:    "Inactive",
		},
		Form: FormLabels{
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
		Actions: ActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: ConfirmLabels{
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
		Buttons: ButtonLabels{AddSupplierPlan: "Add Supplier Plan"},
		Bulk:    BulkLabels{Delete: "Delete"},
		Status: StatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: EmptyLabels{
			Title:   "No supplier plans yet",
			Message: "Add a supplier plan to group cost plans and product plans for a vendor.",
		},
		Errors: ErrorLabels{
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
