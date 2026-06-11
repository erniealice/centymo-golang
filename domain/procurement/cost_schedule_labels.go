package procurement

// cost_schedule_labels.go — extracted verbatim from the root labels.go
// (centymo W7). Pure structural move — no behaviour change.

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
