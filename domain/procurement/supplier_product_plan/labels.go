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
	HeadingActive   string `json:"heading_active"`
	HeadingInactive string `json:"heading_inactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"caption_active"`
	CaptionInactive string `json:"caption_inactive"`
	PageTitle       string `json:"page_title"`
}

type ColumnLabels struct {
	SupplierPlan   string `json:"supplier_plan"`
	Product        string `json:"product"`
	ProductVariant string `json:"product_variant"`
	SupplierSKU    string `json:"supplier_sku"`
	SupplierUnit   string `json:"supplier_unit"`
	Active         string `json:"active"`
}

type TabLabels struct {
	Info          string `json:"info"`
	CostPlanLines string `json:"cost_plan_lines"`
	Activity      string `json:"activity"`
}

type DetailLabels struct {
	InfoSection    string `json:"info_section"`
	SupplierPlan   string `json:"supplier_plan"`
	Product        string `json:"product"`
	ProductVariant string `json:"product_variant"`
	SupplierSKU    string `json:"supplier_sku"`
	SupplierUnit   string `json:"supplier_unit"`
	Active         string `json:"active"`
	Inactive       string `json:"inactive"`
}

type FormLabels struct {
	SectionIdentification string `json:"section_identification"`
	SectionRelationships  string `json:"section_relationships"`
	SectionConfiguration  string `json:"section_configuration"`
	SectionSchedule       string `json:"section_schedule"`
	SectionNotes          string `json:"section_notes"`

	SupplierPlan              string `json:"supplier_plan"`
	SupplierPlanPlaceholder   string `json:"supplier_plan_placeholder"`
	Product                   string `json:"product"`
	ProductPlaceholder        string `json:"product_placeholder"`
	ProductVariant            string `json:"product_variant"`
	ProductVariantPlaceholder string `json:"product_variant_placeholder"`
	SupplierSKU               string `json:"supplier_sku"`
	SupplierSKUPlaceholder    string `json:"supplier_sku_placeholder"`
	SupplierUnit              string `json:"supplier_unit"`
	SupplierUnitPlaceholder   string `json:"supplier_unit_placeholder"`
	Active                    string `json:"active"`
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
	AddSupplierProductPlan string `json:"add_supplier_product_plan"`
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
