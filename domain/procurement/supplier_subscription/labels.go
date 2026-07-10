package supplier_subscription

// supplier_subscription_labels.go — extracted verbatim from the root labels.go
// (centymo W7). Pure structural move — no behaviour change.

// ---------------------------------------------------------------------------
// P3 — SupplierSubscription labels (20260506-supplier-subscriptions)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the supplier_subscription module.
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
	Name      string `json:"name"`
	Supplier  string `json:"supplier"`
	CostPlan  string `json:"cost_plan"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Active    string `json:"active"`
	AutoRenew string `json:"auto_renew"`
	Code      string `json:"code"`
}

type TabLabels struct {
	Info                 string `json:"info"`
	CostPlan             string `json:"cost_plan"`
	LinkedExpenditures   string `json:"linked_expenditures"`
	LinkedPurchaseOrders string `json:"linked_purchase_orders"`
	LinkedRecognitions   string `json:"linked_recognitions"`
	Activity             string `json:"activity"`
}

type DetailLabels struct {
	InfoSection string `json:"info_section"`
	Name        string `json:"name"`
	Supplier    string `json:"supplier"`
	CostPlan    string `json:"cost_plan"`
	Code        string `json:"code"`
	Status      string `json:"status"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Active      string `json:"active"`
	Inactive    string `json:"inactive"`
	AutoRenew   string `json:"auto_renew"`
	Location    string `json:"location"`
	Notes       string `json:"notes"`

	// Linked-recognitions tab (4.4)
	Recognitions RecognitionsLabels `json:"recognitions"`
}

// RecognitionsLabels labels the linked-recognitions tab
// table headers and empty state on the supplier_subscription detail page.
type RecognitionsLabels struct {
	Name            string `json:"name"`
	Status          string `json:"status"`
	RecognitionDate string `json:"recognition_date"`
	Amount          string `json:"amount"`
	EmptyTitle      string `json:"empty_title"`
	EmptyMessage    string `json:"empty_message"`
}

type FormLabels struct {
	SectionIdentification string `json:"section_identification"`
	SectionRelationships  string `json:"section_relationships"`
	SectionConfiguration  string `json:"section_configuration"`
	SectionSchedule       string `json:"section_schedule"`
	SectionNotes          string `json:"section_notes"`

	Name                string `json:"name"`
	NamePlaceholder     string `json:"name_placeholder"`
	Code                string `json:"code"`
	CodePlaceholder     string `json:"code_placeholder"`
	Supplier            string `json:"supplier"`
	SupplierPlaceholder string `json:"supplier_placeholder"`
	SupplierSearch      string `json:"supplier_search"`
	SupplierNoResults   string `json:"supplier_no_results"`
	CostPlan            string `json:"cost_plan"`
	CostPlanPlaceholder string `json:"cost_plan_placeholder"`
	CostPlanSearch      string `json:"cost_plan_search"`
	CostPlanNoResults   string `json:"cost_plan_no_results"`
	AutoRenew           string `json:"auto_renew"`
	Active              string `json:"active"`
	StartDate           string `json:"start_date"`
	StartTime           string `json:"start_time"`
	EndDate             string `json:"end_date"`
	EndTime             string `json:"end_time"`
	TimePlaceholder     string `json:"time_placeholder"`
	Notes               string `json:"notes"`
	NotesPlaceholder    string `json:"notes_placeholder"`
	CurrencyError       string `json:"currency_error"`
	EditLockedReason    string `json:"edit_locked_reason"`
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
	AddSupplierSubscription string `json:"add_supplier_subscription"`
	RecognizeExpense        string `json:"recognize_expense"`
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

// DefaultLabels returns English fallback labels for the supplier_subscription module.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:         "Supplier Subscriptions",
			HeadingActive:   "Active Supplier Subscriptions",
			HeadingInactive: "Inactive Supplier Subscriptions",
			Caption:         "Recurring supplier commitments",
			CaptionActive:   "Active recurring supplier commitments",
			CaptionInactive: "Inactive recurring supplier commitments",
			PageTitle:       "Supplier Subscription",
		},
		Columns: ColumnLabels{
			Name:      "Name",
			Supplier:  "Supplier",
			CostPlan:  "Cost Plan",
			StartDate: "Start Date",
			EndDate:   "End Date",
			Active:    "Status",
			Code:      "Code",
		},
		Tabs: TabLabels{
			Info:                 "Info",
			CostPlan:             "Cost Plan",
			LinkedExpenditures:   "Expenditures",
			LinkedPurchaseOrders: "Purchase Orders",
			LinkedRecognitions:   "Recognitions",
			Activity:             "Activity",
		},
		Detail: DetailLabels{
			InfoSection: "Subscription Details",
			Name:        "Name",
			Supplier:    "Supplier",
			CostPlan:    "Cost Plan",
			Code:        "Code",
			StartDate:   "Start Date",
			EndDate:     "End Date",
			Active:      "Active",
			Inactive:    "Inactive",
			AutoRenew:   "Auto-renew",
			Location:    "Location",
			Notes:       "Notes",
		},
		Form: FormLabels{
			SectionIdentification: "Identification",
			SectionRelationships:  "Relationships",
			SectionConfiguration:  "Configuration",
			SectionSchedule:       "Schedule",
			SectionNotes:          "Notes",
			Name:                  "Name",
			NamePlaceholder:       "e.g. Cloud Hosting — AWS",
			Code:                  "Code",
			CodePlaceholder:       "e.g. SUB-2026-001",
			Supplier:              "Supplier",
			SupplierPlaceholder:   "Search supplier…",
			SupplierSearch:        "Search suppliers",
			SupplierNoResults:     "No suppliers found",
			CostPlan:              "Cost Plan",
			CostPlanPlaceholder:   "Search cost plan…",
			CostPlanSearch:        "Search cost plans",
			CostPlanNoResults:     "No cost plans found",
			AutoRenew:             "Auto-renew",
			Active:                "Active",
			StartDate:             "Start Date",
			StartTime:             "Start Time",
			EndDate:               "End Date",
			EndTime:               "End Time",
			TimePlaceholder:       "HH:MM",
			Notes:                 "Notes",
			NotesPlaceholder:      "Internal notes about this subscription",
			CurrencyError:         "The selected cost plan's billing currency does not match the workspace functional currency.",
			EditLockedReason:      "This subscription has linked expenditures and cannot be fully edited.",
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
			Delete:                "Delete Supplier Subscription",
			DeleteMessage:         "Are you sure you want to delete this supplier subscription?",
			Activate:              "Activate Supplier Subscription",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Supplier Subscription",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Supplier Subscriptions",
			BulkDeleteMessage:     "Delete selected supplier subscriptions?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected supplier subscriptions?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected supplier subscriptions?",
		},
		Buttons: ButtonLabels{
			AddSupplierSubscription: "Add Supplier Subscription",
			RecognizeExpense:        "Recognize Expense",
		},
		Bulk: BulkLabels{
			Delete: "Delete",
		},
		Status: StatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: EmptyLabels{
			Title:   "No supplier subscriptions yet",
			Message: "Add a supplier subscription to start tracking recurring vendor commitments.",
		},
		Errors: ErrorLabels{
			PermissionDenied: "You do not have permission to perform this action.",
			InvalidFormData:  "Invalid form data. Please check your inputs and try again.",
			NotFound:         "Supplier subscription not found.",
			IDRequired:       "Supplier subscription ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This subscription is in use and cannot be deleted.",
			LoadFailed:       "Failed to load supplier subscription.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}
