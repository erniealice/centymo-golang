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
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type ColumnLabels struct {
	Name      string `json:"name"`
	Supplier  string `json:"supplier"`
	CostPlan  string `json:"costPlan"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Active    string `json:"active"`
	AutoRenew string `json:"autoRenew"`
	Code      string `json:"code"`
}

type TabLabels struct {
	Info                 string `json:"info"`
	CostPlan             string `json:"costPlan"`
	LinkedExpenditures   string `json:"linkedExpenditures"`
	LinkedPurchaseOrders string `json:"linkedPurchaseOrders"`
	LinkedRecognitions   string `json:"linkedRecognitions"`
	Activity             string `json:"activity"`
}

type DetailLabels struct {
	InfoSection string `json:"infoSection"`
	Name        string `json:"name"`
	Supplier    string `json:"supplier"`
	CostPlan    string `json:"costPlan"`
	Code        string `json:"code"`
	Status      string `json:"status"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Active      string `json:"active"`
	Inactive    string `json:"inactive"`
	AutoRenew   string `json:"autoRenew"`
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
	RecognitionDate string `json:"recognitionDate"`
	Amount          string `json:"amount"`
	EmptyTitle      string `json:"emptyTitle"`
	EmptyMessage    string `json:"emptyMessage"`
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
	SupplierSearch      string `json:"supplierSearch"`
	SupplierNoResults   string `json:"supplierNoResults"`
	CostPlan            string `json:"costPlan"`
	CostPlanPlaceholder string `json:"costPlanPlaceholder"`
	CostPlanSearch      string `json:"costPlanSearch"`
	CostPlanNoResults   string `json:"costPlanNoResults"`
	AutoRenew           string `json:"autoRenew"`
	Active              string `json:"active"`
	StartDate           string `json:"startDate"`
	StartTime           string `json:"startTime"`
	EndDate             string `json:"endDate"`
	EndTime             string `json:"endTime"`
	TimePlaceholder     string `json:"timePlaceholder"`
	Notes               string `json:"notes"`
	NotesPlaceholder    string `json:"notesPlaceholder"`
	CurrencyError       string `json:"currencyError"`
	EditLockedReason    string `json:"editLockedReason"`
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
	AddSupplierSubscription string `json:"addSupplierSubscription"`
	RecognizeExpense        string `json:"recognizeExpense"`
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
