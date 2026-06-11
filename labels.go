package centymo

// labels.go — RESIDUAL after centymo W6.
//
// W1-W6 dissolved this god-file per domain. What remains:
//   - LocationMap / LocationDisplayName (entydad-bound; WL deferral, not yet landed)
//   - ProcurementLabels (the Procurement Operations composition app — procurement
//     domain, lands in W7)
//   - the W7 procurement-domain label sections (SupplierSubscription, CostSchedule,
//     SupplierPlan, CostPlan, SupplierProductPlan, SupplierProductCostPlan).
// Expenditure-domain labels moved to domain/expenditure/<entity>_labels.go in W6.

var LocationMap = map[string]string{
	"ayala-central-bloc": "Ayala Central Bloc",
	"sm-city-cebu":       "SM City Cebu",
	"ayala-center-cebu":  "Ayala Center Cebu",
	"robinsons-galleria": "Robinsons Galleria",
}

// LocationDisplayName returns the display name for a location slug.
func LocationDisplayName(slug string) string {
	if name, ok := LocationMap[slug]; ok {
		return name
	}
	return slug
}

// ---------------------------------------------------------------------------
// P3b — Procurement Operations app labels
// (composition surface, no proto entity — mirrors the schedule/cyta pattern)
// ---------------------------------------------------------------------------

// ProcurementLabels holds all translatable strings for the Procurement
// Operations composition app. Populated via lyngua (P4). These keys are
// intentionally generic so they render without overrides when lyngua has not
// yet supplied values.
type ProcurementLabels struct {
	AppLabel              string `json:"app_label"`
	DashboardTitle        string `json:"dashboard_title"`
	PendingApprovalsTitle string `json:"pending_approvals_title"`
	ExpiringTitle         string `json:"expiring_title"`
	VarianceTitle         string `json:"variance_title"`
	RecurrenceTitle       string `json:"recurrence_title"`
	RenewalsTitle         string `json:"renewals_title"`
	UtilizationTitle      string `json:"utilization_title"`
	EmptyRenewals         string `json:"empty_renewals"`
	EmptyVariance         string `json:"empty_variance"`
	EmptyUtilization      string `json:"empty_utilization"`
	EmptyRecurrence       string `json:"empty_recurrence"`
	DaysUntilExpiry       string `json:"days_until_expiry"`
	UtilizationPercent    string `json:"utilization_percent"`
	BudgetPressureLabel   string `json:"budget_pressure_label"`
}

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

// ---------------------------------------------------------------------------
// P3 — CostSchedule labels
// ---------------------------------------------------------------------------

// CostScheduleLabels holds all translatable strings for the cost_schedule module.
type CostScheduleLabels struct {
	Page    CostSchedulePageLabels    `json:"page"`
	Columns CostScheduleColumnLabels  `json:"columns"`
	Tabs    CostScheduleTabLabels     `json:"tabs"`
	Detail  CostScheduleDetailLabels  `json:"detail"`
	Form    CostScheduleFormLabels    `json:"form"`
	Actions CostScheduleActionLabels  `json:"actions"`
	Confirm CostScheduleConfirmLabels `json:"confirm"`
	Buttons CostScheduleButtonLabels  `json:"buttons"`
	Bulk    CostScheduleBulkLabels    `json:"bulk"`
	Status  CostScheduleStatusLabels  `json:"status"`
	Empty   CostScheduleEmptyLabels   `json:"empty"`
	Errors  CostScheduleErrorLabels   `json:"errors"`
}

type CostSchedulePageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type CostScheduleColumnLabels struct {
	Name      string `json:"name"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Location  string `json:"location"`
	Active    string `json:"active"`
}

type CostScheduleTabLabels struct {
	Info      string `json:"info"`
	CostPlans string `json:"costPlans"`
	Activity  string `json:"activity"`
}

type CostScheduleDetailLabels struct {
	InfoSection string `json:"infoSection"`
	Name        string `json:"name"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Active      string `json:"active"`
	Inactive    string `json:"inactive"`
}

type CostScheduleFormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	Name                string `json:"name"`
	NamePlaceholder     string `json:"namePlaceholder"`
	Description         string `json:"description"`
	DescPlaceholder     string `json:"descPlaceholder"`
	StartDate           string `json:"startDate"`
	EndDate             string `json:"endDate"`
	Location            string `json:"location"`
	LocationPlaceholder string `json:"locationPlaceholder"`
	Active              string `json:"active"`
}

type CostScheduleActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type CostScheduleConfirmLabels struct {
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

type CostScheduleButtonLabels struct {
	AddCostSchedule string `json:"addCostSchedule"`
}

type CostScheduleBulkLabels struct {
	Delete string `json:"delete"`
}

type CostScheduleStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type CostScheduleEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type CostScheduleErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultCostScheduleLabels returns English fallback labels.
func DefaultCostScheduleLabels() CostScheduleLabels {
	return CostScheduleLabels{
		Page: CostSchedulePageLabels{
			Heading:         "Cost Schedules",
			HeadingActive:   "Active Cost Schedules",
			HeadingInactive: "Inactive Cost Schedules",
			Caption:         "Date-bounded supplier pricing windows",
			CaptionActive:   "Active pricing windows",
			CaptionInactive: "Inactive pricing windows",
			PageTitle:       "Cost Schedule",
		},
		Columns: CostScheduleColumnLabels{
			Name:      "Name",
			StartDate: "Start Date",
			EndDate:   "End Date",
			Location:  "Location",
			Active:    "Status",
		},
		Tabs: CostScheduleTabLabels{
			Info:      "Info",
			CostPlans: "Cost Plans",
			Activity:  "Activity",
		},
		Detail: CostScheduleDetailLabels{
			InfoSection: "Schedule Details",
			Name:        "Name",
			StartDate:   "Start Date",
			EndDate:     "End Date",
			Location:    "Location",
			Description: "Description",
			Active:      "Active",
			Inactive:    "Inactive",
		},
		Form: CostScheduleFormLabels{
			SectionIdentification: "Identification",
			SectionRelationships:  "Relationships",
			SectionConfiguration:  "Configuration",
			SectionSchedule:       "Schedule",
			SectionNotes:          "Notes",
			Name:                  "Name",
			NamePlaceholder:       "e.g. Q1 2026 Supplier Rates",
			Description:           "Description",
			DescPlaceholder:       "Internal notes about this cost schedule",
			StartDate:             "Start Date",
			EndDate:               "End Date",
			Location:              "Location",
			LocationPlaceholder:   "Select location",
			Active:                "Active",
		},
		Actions: CostScheduleActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: CostScheduleConfirmLabels{
			Delete:                "Delete Cost Schedule",
			DeleteMessage:         "Are you sure you want to delete this cost schedule?",
			Activate:              "Activate Cost Schedule",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Cost Schedule",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Cost Schedules",
			BulkDeleteMessage:     "Delete selected cost schedules?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected cost schedules?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected cost schedules?",
		},
		Buttons: CostScheduleButtonLabels{
			AddCostSchedule: "Add Cost Schedule",
		},
		Bulk: CostScheduleBulkLabels{Delete: "Delete"},
		Status: CostScheduleStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: CostScheduleEmptyLabels{
			Title:   "No cost schedules yet",
			Message: "Add a cost schedule to group supplier cost plans by date range.",
		},
		Errors: CostScheduleErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Cost schedule not found.",
			IDRequired:       "Cost schedule ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This cost schedule has linked cost plans and cannot be deleted.",
			LoadFailed:       "Failed to load cost schedule.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}

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

// ---------------------------------------------------------------------------
// P3 — CostPlan labels
// ---------------------------------------------------------------------------

// CostPlanLabels holds all translatable strings for the cost_plan module.
type CostPlanLabels struct {
	Page    CostPlanPageLabels    `json:"page"`
	Columns CostPlanColumnLabels  `json:"columns"`
	Tabs    CostPlanTabLabels     `json:"tabs"`
	Detail  CostPlanDetailLabels  `json:"detail"`
	Form    CostPlanFormLabels    `json:"form"`
	Actions CostPlanActionLabels  `json:"actions"`
	Confirm CostPlanConfirmLabels `json:"confirm"`
	Buttons CostPlanButtonLabels  `json:"buttons"`
	Bulk    CostPlanBulkLabels    `json:"bulk"`
	Status  CostPlanStatusLabels  `json:"status"`
	Empty   CostPlanEmptyLabels   `json:"empty"`
	Errors  CostPlanErrorLabels   `json:"errors"`
}

type CostPlanPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type CostPlanColumnLabels struct {
	Name         string `json:"name"`
	BillingKind  string `json:"billingKind"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	SupplierPlan string `json:"supplierPlan"`
	CostSchedule string `json:"costSchedule"`
	Active       string `json:"active"`
}

type CostPlanTabLabels struct {
	Info                string `json:"info"`
	Lines               string `json:"lines"`
	LinkedSubscriptions string `json:"linkedSubscriptions"`
	Activity            string `json:"activity"`
}

type CostPlanDetailLabels struct {
	InfoSection  string `json:"infoSection"`
	Name         string `json:"name"`
	BillingKind  string `json:"billingKind"`
	AmountBasis  string `json:"amountBasis"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	BillingCycle string `json:"billingCycle"`
	DefaultTerm  string `json:"defaultTerm"`
	SupplierPlan string `json:"supplierPlan"`
	CostSchedule string `json:"costSchedule"`
	Active       string `json:"active"`
	Inactive     string `json:"inactive"`
}

type CostPlanFormLabels struct {
	SectionIdentification string `json:"sectionIdentification"`
	SectionRelationships  string `json:"sectionRelationships"`
	SectionConfiguration  string `json:"sectionConfiguration"`
	SectionSchedule       string `json:"sectionSchedule"`
	SectionNotes          string `json:"sectionNotes"`

	Name                    string `json:"name"`
	NamePlaceholder         string `json:"namePlaceholder"`
	Description             string `json:"description"`
	DescPlaceholder         string `json:"descPlaceholder"`
	SupplierPlan            string `json:"supplierPlan"`
	SupplierPlanPlaceholder string `json:"supplierPlanPlaceholder"`
	CostSchedule            string `json:"costSchedule"`
	CostSchedulePlaceholder string `json:"costSchedulePlaceholder"`
	BillingKind             string `json:"billingKind"`
	AmountBasis             string `json:"amountBasis"`
	Amount                  string `json:"amount"`
	AmountPlaceholder       string `json:"amountPlaceholder"`
	Currency                string `json:"currency"`
	CurrencyPlaceholder     string `json:"currencyPlaceholder"`
	BillingCycle            string `json:"billingCycle"`
	BillingCyclePlaceholder string `json:"billingCyclePlaceholder"`
	DefaultTerm             string `json:"defaultTerm"`
	DefaultTermPlaceholder  string `json:"defaultTermPlaceholder"`
	Active                  string `json:"active"`

	// BillingKind option labels
	BillingKindOneTime    string `json:"billingKindOneTime"`
	BillingKindRecurring  string `json:"billingKindRecurring"`
	BillingKindContract   string `json:"billingKindContract"`
	BillingKindUsageBased string `json:"billingKindUsageBased"`
	BillingKindAdHoc      string `json:"billingKindAdHoc"`

	// AmountBasis option labels
	AmountBasisPerCycle         string `json:"amountBasisPerCycle"`
	AmountBasisTotalPackage     string `json:"amountBasisTotalPackage"`
	AmountBasisDerivedFromLines string `json:"amountBasisDerivedFromLines"`
	AmountBasisPerOccurrence    string `json:"amountBasisPerOccurrence"`

	// Duration unit option labels (shared by billing_cycle_unit and default_term_unit)
	DurationUnitDay   string `json:"durationUnitDay"`
	DurationUnitWeek  string `json:"durationUnitWeek"`
	DurationUnitMonth string `json:"durationUnitMonth"`
	DurationUnitYear  string `json:"durationUnitYear"`
}

type CostPlanActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type CostPlanConfirmLabels struct {
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

type CostPlanButtonLabels struct {
	AddCostPlan string `json:"addCostPlan"`
}

type CostPlanBulkLabels struct {
	Delete string `json:"delete"`
}

type CostPlanStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type CostPlanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type CostPlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultCostPlanLabels returns English fallback labels.
func DefaultCostPlanLabels() CostPlanLabels {
	return CostPlanLabels{
		Page: CostPlanPageLabels{
			Heading:         "Cost Plans",
			HeadingActive:   "Active Cost Plans",
			HeadingInactive: "Inactive Cost Plans",
			Caption:         "Supplier pricing plans and billing schedules",
			CaptionActive:   "Active cost plans",
			CaptionInactive: "Inactive cost plans",
			PageTitle:       "Cost Plan",
		},
		Columns: CostPlanColumnLabels{
			Name:         "Name",
			BillingKind:  "Billing Kind",
			Amount:       "Amount",
			Currency:     "Currency",
			SupplierPlan: "Supplier Plan",
			CostSchedule: "Cost Schedule",
			Active:       "Status",
		},
		Tabs: CostPlanTabLabels{
			Info:                "Info",
			Lines:               "Lines",
			LinkedSubscriptions: "Subscriptions",
			Activity:            "Activity",
		},
		Detail: CostPlanDetailLabels{
			InfoSection:  "Cost Plan Details",
			Name:         "Name",
			BillingKind:  "Billing Kind",
			AmountBasis:  "Amount Basis",
			Amount:       "Amount",
			Currency:     "Currency",
			BillingCycle: "Billing Cycle",
			DefaultTerm:  "Default Term",
			SupplierPlan: "Supplier Plan",
			CostSchedule: "Cost Schedule",
			Active:       "Active",
			Inactive:     "Inactive",
		},
		Form: CostPlanFormLabels{
			SectionIdentification:       "Identification",
			SectionRelationships:        "Relationships",
			SectionConfiguration:        "Configuration",
			SectionSchedule:             "Schedule",
			SectionNotes:                "Notes",
			Name:                        "Name",
			NamePlaceholder:             "e.g. AWS EC2 Monthly",
			Description:                 "Description",
			DescPlaceholder:             "Internal notes about this cost plan",
			SupplierPlan:                "Supplier Plan",
			SupplierPlanPlaceholder:     "Select supplier plan",
			CostSchedule:                "Cost Schedule",
			CostSchedulePlaceholder:     "Select cost schedule",
			BillingKind:                 "Billing Kind",
			AmountBasis:                 "Amount Basis",
			Amount:                      "Amount",
			AmountPlaceholder:           "0.00",
			Currency:                    "Currency",
			CurrencyPlaceholder:         "e.g. PHP",
			BillingCycle:                "Billing Cycle",
			BillingCyclePlaceholder:     "e.g. 1",
			DefaultTerm:                 "Default Term",
			DefaultTermPlaceholder:      "e.g. 12",
			Active:                      "Active",
			BillingKindOneTime:          "One Time",
			BillingKindRecurring:        "Recurring",
			BillingKindContract:         "Contract",
			BillingKindUsageBased:       "Usage Based",
			BillingKindAdHoc:            "Ad Hoc",
			AmountBasisPerCycle:         "Per Cycle",
			AmountBasisTotalPackage:     "Total Package",
			AmountBasisDerivedFromLines: "Derived From Lines",
			AmountBasisPerOccurrence:    "Per Occurrence",
			DurationUnitDay:             "Day",
			DurationUnitWeek:            "Week",
			DurationUnitMonth:           "Month",
			DurationUnitYear:            "Year",
		},
		Actions: CostPlanActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: CostPlanConfirmLabels{
			Delete:                "Delete Cost Plan",
			DeleteMessage:         "Are you sure you want to delete this cost plan?",
			Activate:              "Activate Cost Plan",
			ActivateMessage:       "Activate %s?",
			Deactivate:            "Deactivate Cost Plan",
			DeactivateMessage:     "Deactivate %s?",
			BulkDelete:            "Delete Cost Plans",
			BulkDeleteMessage:     "Delete selected cost plans?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Activate selected cost plans?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Deactivate selected cost plans?",
		},
		Buttons: CostPlanButtonLabels{AddCostPlan: "Add Cost Plan"},
		Bulk:    CostPlanBulkLabels{Delete: "Delete"},
		Status: CostPlanStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: CostPlanEmptyLabels{
			Title:   "No cost plans yet",
			Message: "Add a cost plan to define billing terms for a supplier engagement.",
		},
		Errors: CostPlanErrorLabels{
			PermissionDenied: "You do not have permission.",
			InvalidFormData:  "Invalid form data.",
			NotFound:         "Cost plan not found.",
			IDRequired:       "Cost plan ID is required.",
			NoPermission:     "No permission.",
			InUse:            "This cost plan has linked subscriptions and cannot be deleted.",
			LoadFailed:       "Failed to load cost plan.",
			NoIDsProvided:    "No IDs provided.",
		},
	}
}

// ---------------------------------------------------------------------------
// P3 — SupplierProductPlan labels
// ---------------------------------------------------------------------------

// SupplierProductPlanLabels holds all translatable strings for the supplier_product_plan module.
type SupplierProductPlanLabels struct {
	Page    SupplierProductPlanPageLabels    `json:"page"`
	Columns SupplierProductPlanColumnLabels  `json:"columns"`
	Tabs    SupplierProductPlanTabLabels     `json:"tabs"`
	Detail  SupplierProductPlanDetailLabels  `json:"detail"`
	Form    SupplierProductPlanFormLabels    `json:"form"`
	Actions SupplierProductPlanActionLabels  `json:"actions"`
	Confirm SupplierProductPlanConfirmLabels `json:"confirm"`
	Buttons SupplierProductPlanButtonLabels  `json:"buttons"`
	Bulk    SupplierProductPlanBulkLabels    `json:"bulk"`
	Status  SupplierProductPlanStatusLabels  `json:"status"`
	Empty   SupplierProductPlanEmptyLabels   `json:"empty"`
	Errors  SupplierProductPlanErrorLabels   `json:"errors"`
}

type SupplierProductPlanPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
	PageTitle       string `json:"pageTitle"`
}

type SupplierProductPlanColumnLabels struct {
	SupplierPlan   string `json:"supplierPlan"`
	Product        string `json:"product"`
	ProductVariant string `json:"productVariant"`
	SupplierSKU    string `json:"supplierSku"`
	SupplierUnit   string `json:"supplierUnit"`
	Active         string `json:"active"`
}

type SupplierProductPlanTabLabels struct {
	Info          string `json:"info"`
	CostPlanLines string `json:"costPlanLines"`
	Activity      string `json:"activity"`
}

type SupplierProductPlanDetailLabels struct {
	InfoSection    string `json:"infoSection"`
	SupplierPlan   string `json:"supplierPlan"`
	Product        string `json:"product"`
	ProductVariant string `json:"productVariant"`
	SupplierSKU    string `json:"supplierSku"`
	SupplierUnit   string `json:"supplierUnit"`
	Active         string `json:"active"`
	Inactive       string `json:"inactive"`
}

type SupplierProductPlanFormLabels struct {
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

type SupplierProductPlanActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Activate     string `json:"activate"`
	Deactivate   string `json:"deactivate"`
	NoPermission string `json:"noPermission"`
}

type SupplierProductPlanConfirmLabels struct {
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

type SupplierProductPlanButtonLabels struct {
	AddSupplierProductPlan string `json:"addSupplierProductPlan"`
}

type SupplierProductPlanBulkLabels struct {
	Delete string `json:"delete"`
}

type SupplierProductPlanStatusLabels struct {
	Active     string `json:"active"`
	Inactive   string `json:"inactive"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type SupplierProductPlanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SupplierProductPlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	InUse            string `json:"inUse"`
	LoadFailed       string `json:"loadFailed"`
	NoIDsProvided    string `json:"noIdsProvided"`
}

// DefaultSupplierProductPlanLabels returns English fallback labels.
func DefaultSupplierProductPlanLabels() SupplierProductPlanLabels {
	return SupplierProductPlanLabels{
		Page: SupplierProductPlanPageLabels{
			Heading:         "Supplier Product Plans",
			HeadingActive:   "Active Supplier Product Plans",
			HeadingInactive: "Inactive Supplier Product Plans",
			Caption:         "Supplier product catalogue line items",
			CaptionActive:   "Active supplier product plans",
			CaptionInactive: "Inactive supplier product plans",
			PageTitle:       "Supplier Product Plan",
		},
		Columns: SupplierProductPlanColumnLabels{
			SupplierPlan:   "Supplier Plan",
			Product:        "Product",
			ProductVariant: "Variant",
			SupplierSKU:    "Supplier SKU",
			SupplierUnit:   "Supplier Unit",
			Active:         "Status",
		},
		Tabs: SupplierProductPlanTabLabels{
			Info:          "Info",
			CostPlanLines: "Cost Plan Lines",
			Activity:      "Activity",
		},
		Detail: SupplierProductPlanDetailLabels{
			InfoSection:    "Product Plan Details",
			SupplierPlan:   "Supplier Plan",
			Product:        "Product",
			ProductVariant: "Variant",
			SupplierSKU:    "Supplier SKU",
			SupplierUnit:   "Supplier Unit",
			Active:         "Active",
			Inactive:       "Inactive",
		},
		Form: SupplierProductPlanFormLabels{
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
		Actions: SupplierProductPlanActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			Activate:     "Activate",
			Deactivate:   "Deactivate",
			NoPermission: "No permission",
		},
		Confirm: SupplierProductPlanConfirmLabels{
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
		Buttons: SupplierProductPlanButtonLabels{AddSupplierProductPlan: "Add Supplier Product Plan"},
		Bulk:    SupplierProductPlanBulkLabels{Delete: "Delete"},
		Status: SupplierProductPlanStatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: SupplierProductPlanEmptyLabels{
			Title:   "No supplier product plans yet",
			Message: "Add a supplier product plan to map vendor catalogue items to your internal products.",
		},
		Errors: SupplierProductPlanErrorLabels{
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
