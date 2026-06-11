package supplier_product_plan

// supplier_product_plan_labels.go — extracted verbatim from the root labels.go
// (centymo W7). Pure structural move — no behaviour change.

// ---------------------------------------------------------------------------
// P3 — SupplierProductPlan labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the supplier_product_plan module.
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
	SupplierPlan   string `json:"supplierPlan"`
	Product        string `json:"product"`
	ProductVariant string `json:"productVariant"`
	SupplierSKU    string `json:"supplierSku"`
	SupplierUnit   string `json:"supplierUnit"`
	Active         string `json:"active"`
}

type TabLabels struct {
	Info          string `json:"info"`
	CostPlanLines string `json:"costPlanLines"`
	Activity      string `json:"activity"`
}

type DetailLabels struct {
	InfoSection    string `json:"infoSection"`
	SupplierPlan   string `json:"supplierPlan"`
	Product        string `json:"product"`
	ProductVariant string `json:"productVariant"`
	SupplierSKU    string `json:"supplierSku"`
	SupplierUnit   string `json:"supplierUnit"`
	Active         string `json:"active"`
	Inactive       string `json:"inactive"`
}

type FormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	SupplierPlan              string `json:"supplierPlan"`
	SupplierPlanPlaceholder   string `json:"supplierPlanPlaceholder"`
	Product                   string `json:"product"`
	ProductPlaceholder        string `json:"productPlaceholder"`
	ProductVariant            string `json:"productVariant"`
	ProductVariantPlaceholder string `json:"productVariantPlaceholder"`
	SupplierSKU               string `json:"supplierSku"`
	SupplierSKUPlaceholder    string `json:"supplierSkuPlaceholder"`
	SupplierUnit              string `json:"supplierUnit"`
	SupplierUnitPlaceholder   string `json:"supplierUnitPlaceholder"`
	Active                    string `json:"active"`
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
	AddSupplierProductPlan string `json:"addSupplierProductPlan"`
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
			Heading:         "Supplier Product Plans",
			HeadingActive:   "Active Supplier Product Plans",
			HeadingInactive: "Inactive Supplier Product Plans",
			Caption:         "Supplier product catalogue line items",
			CaptionActive:   "Active supplier product plans",
			CaptionInactive: "Inactive supplier product plans",
			PageTitle:       "Supplier Product Plan",
		},
		Columns: ColumnLabels{
			SupplierPlan:   "Supplier Plan",
			Product:        "Product",
			ProductVariant: "Variant",
			SupplierSKU:    "Supplier SKU",
			SupplierUnit:   "Supplier Unit",
			Active:         "Status",
		},
		Tabs: TabLabels{
			Info:          "Info",
			CostPlanLines: "Cost Plan Lines",
			Activity:      "Activity",
		},
		Detail: DetailLabels{
			InfoSection:    "Product Plan Details",
			SupplierPlan:   "Supplier Plan",
			Product:        "Product",
			ProductVariant: "Variant",
			SupplierSKU:    "Supplier SKU",
			SupplierUnit:   "Supplier Unit",
			Active:         "Active",
			Inactive:       "Inactive",
		},
		Form: FormLabels{
			SectionIdentification:     "Identification",
			SectionRelationships:      "Relationships",
			SectionConfiguration:      "Configuration",
			SectionSchedule:           "Schedule",
			SectionNotes:              "Notes",
			SupplierPlan:              "Supplier Plan",
			SupplierPlanPlaceholder:   "Select supplier plan",
			Product:                   "Product",
			ProductPlaceholder:        "Select product",
			ProductVariant:            "Variant (optional)",
			ProductVariantPlaceholder: "Select variant",
			SupplierSKU:               "Supplier SKU",
			SupplierSKUPlaceholder:    "Supplier's internal SKU code",
			SupplierUnit:              "Supplier Unit",
			SupplierUnitPlaceholder:   "e.g. vCPU·hour",
			Active:                    "Active",
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
			Delete:                "Delete Supplier Product Plan",
			DeleteMessage:         "Are you sure you want to delete this supplier product plan?",
			Activate:              "Activate Supplier Product Plan",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Supplier Product Plan",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Supplier Product Plans",
			BulkDeleteMessage:     "Delete selected supplier product plans?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected supplier product plans?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected supplier product plans?",
		},
		Buttons: ButtonLabels{AddSupplierProductPlan: "Add Supplier Product Plan"},
		Bulk:    BulkLabels{Delete: "Delete"},
		Status: StatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: EmptyLabels{
			Title:   "No supplier product plans yet",
			Message: "Add a supplier product plan to map vendor catalogue items to your internal products.",
		},
		Errors: ErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Supplier product plan not found.",
			IDRequired:       "Supplier product plan ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This supplier product plan has linked cost plan lines and cannot be deleted.",
			LoadFailed:       "Failed to load supplier product plan.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}
