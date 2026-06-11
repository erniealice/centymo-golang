package procurement

// supplier_subscription_labels.go — extracted verbatim from the root labels.go
// (centymo W7). Pure structural move — no behaviour change.

// ---------------------------------------------------------------------------
// P3 — SupplierSubscription labels (20260506-supplier-subscriptions)
// ---------------------------------------------------------------------------

// SupplierSubscriptionLabels holds all translatable strings for the supplier_subscription module.
type SupplierSubscriptionLabels struct {
	Page    SupplierSubscriptionPageLabels    `json:"page"`
	Columns SupplierSubscriptionColumnLabels  `json:"columns"`
	Tabs    SupplierSubscriptionTabLabels     `json:"tabs"`
	Detail  SupplierSubscriptionDetailLabels  `json:"detail"`
	Form    SupplierSubscriptionFormLabels    `json:"form"`
	Actions SupplierSubscriptionActionLabels  `json:"actions"`
	Confirm SupplierSubscriptionConfirmLabels `json:"confirm"`
	Buttons SupplierSubscriptionButtonLabels  `json:"buttons"`
	Bulk    SupplierSubscriptionBulkLabels    `json:"bulk"`
	Status  SupplierSubscriptionStatusLabels  `json:"status"`
	Empty   SupplierSubscriptionEmptyLabels   `json:"empty"`
	Errors  SupplierSubscriptionErrorLabels   `json:"errors"`
}

type SupplierSubscriptionPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type SupplierSubscriptionColumnLabels struct {
	Name      string `json:"name"`
	Supplier  string `json:"supplier"`
	CostPlan  string `json:"costPlan"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Active    string `json:"active"`
	AutoRenew string `json:"autoRenew"`
	Code      string `json:"code"`
}

type SupplierSubscriptionTabLabels struct {
	Info                 string `json:"info"`
	CostPlan             string `json:"costPlan"`
	LinkedExpenditures   string `json:"linkedExpenditures"`
	LinkedPurchaseOrders string `json:"linkedPurchaseOrders"`
	LinkedRecognitions   string `json:"linkedRecognitions"`
	Activity             string `json:"activity"`
}

type SupplierSubscriptionDetailLabels struct {
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
	Recognitions SupplierSubscriptionRecognitionsLabels `json:"recognitions"`
}

// SupplierSubscriptionRecognitionsLabels labels the linked-recognitions tab
// table headers and empty state on the supplier_subscription detail page.
type SupplierSubscriptionRecognitionsLabels struct {
	Name            string `json:"name"`
	Status          string `json:"status"`
	RecognitionDate string `json:"recognitionDate"`
	Amount          string `json:"amount"`
	EmptyTitle      string `json:"emptyTitle"`
	EmptyMessage    string `json:"emptyMessage"`
}

type SupplierSubscriptionFormLabels struct {
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

type SupplierSubscriptionActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type SupplierSubscriptionConfirmLabels struct {
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

type SupplierSubscriptionButtonLabels struct {
	AddSupplierSubscription string `json:"addSupplierSubscription"`
	RecognizeExpense        string `json:"recognizeExpense"`
}

type SupplierSubscriptionBulkLabels struct {
	Delete string `json:"delete"`
}

type SupplierSubscriptionStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type SupplierSubscriptionEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SupplierSubscriptionErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultSupplierSubscriptionLabels returns English fallback labels for the supplier_subscription module.
func DefaultSupplierSubscriptionLabels() SupplierSubscriptionLabels {
	return SupplierSubscriptionLabels{
		Page: SupplierSubscriptionPageLabels{
			Heading:         "Supplier Subscriptions",
			HeadingActive:   "Active Supplier Subscriptions",
			HeadingInactive: "Inactive Supplier Subscriptions",
			Caption:         "Recurring supplier commitments",
			CaptionActive:   "Active recurring supplier commitments",
			CaptionInactive: "Inactive recurring supplier commitments",
			PageTitle:       "Supplier Subscription",
		},
		Columns: SupplierSubscriptionColumnLabels{
			Name:      "Name",
			Supplier:  "Supplier",
			CostPlan:  "Cost Plan",
			StartDate: "Start Date",
			EndDate:   "End Date",
			Active:    "Status",
			Code:      "Code",
		},
		Tabs: SupplierSubscriptionTabLabels{
			Info:                 "Info",
			CostPlan:             "Cost Plan",
			LinkedExpenditures:   "Expenditures",
			LinkedPurchaseOrders: "Purchase Orders",
			LinkedRecognitions:   "Recognitions",
			Activity:             "Activity",
		},
		Detail: SupplierSubscriptionDetailLabels{
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
		Form: SupplierSubscriptionFormLabels{
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
		Actions: SupplierSubscriptionActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: SupplierSubscriptionConfirmLabels{
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
		Buttons: SupplierSubscriptionButtonLabels{
			AddSupplierSubscription: "Add Supplier Subscription",
			RecognizeExpense:        "Recognize Expense",
		},
		Bulk: SupplierSubscriptionBulkLabels{
			Delete: "Delete",
		},
		Status: SupplierSubscriptionStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: SupplierSubscriptionEmptyLabels{
			Title:   "No supplier subscriptions yet",
			Message: "Add a supplier subscription to start tracking recurring vendor commitments.",
		},
		Errors: SupplierSubscriptionErrorLabels{
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
