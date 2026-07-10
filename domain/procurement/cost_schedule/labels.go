package cost_schedule

// cost_schedule_labels.go — extracted verbatim from the root labels.go
// (centymo W7). Pure structural move — no behaviour change.

// ---------------------------------------------------------------------------
// P3 — CostSchedule labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the cost_schedule module.
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
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Location  string `json:"location"`
	Active    string `json:"active"`
}

type TabLabels struct {
	Info      string `json:"info"`
	CostPlans string `json:"cost_plans"`
	Activity  string `json:"activity"`
}

type DetailLabels struct {
	InfoSection string `json:"info_section"`
	Name        string `json:"name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Active      string `json:"active"`
	Inactive    string `json:"inactive"`
}

type FormLabels struct {
	SectionIdentification string `json:"section_identification"`
	SectionRelationships  string `json:"section_relationships"`
	SectionConfiguration  string `json:"section_configuration"`
	SectionSchedule       string `json:"section_schedule"`
	SectionNotes          string `json:"section_notes"`

	Name                string `json:"name"`
	NamePlaceholder     string `json:"name_placeholder"`
	Description         string `json:"description"`
	DescPlaceholder     string `json:"desc_placeholder"`
	StartDate           string `json:"start_date"`
	EndDate             string `json:"end_date"`
	Location            string `json:"location"`
	LocationPlaceholder string `json:"location_placeholder"`
	Active              string `json:"active"`
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
	AddCostSchedule string `json:"add_cost_schedule"`
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
			Heading:         "Cost Schedules",
			HeadingActive:   "Active Cost Schedules",
			HeadingInactive: "Inactive Cost Schedules",
			Caption:         "Date-bounded supplier pricing windows",
			CaptionActive:   "Active pricing windows",
			CaptionInactive: "Inactive pricing windows",
			PageTitle:       "Cost Schedule",
		},
		Columns: ColumnLabels{
			Name:      "Name",
			StartDate: "Start Date",
			EndDate:   "End Date",
			Location:  "Location",
			Active:    "Status",
		},
		Tabs: TabLabels{
			Info:      "Info",
			CostPlans: "Cost Plans",
			Activity:  "Activity",
		},
		Detail: DetailLabels{
			InfoSection: "Schedule Details",
			Name:        "Name",
			StartDate:   "Start Date",
			EndDate:     "End Date",
			Location:    "Location",
			Description: "Description",
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
			NamePlaceholder:       "e.g. Q1 2026 Supplier Rates",
			Description:           "Description",
			DescPlaceholder:       "Internal notes about this cost schedule",
			StartDate:             "Start Date",
			EndDate:               "End Date",
			Location:              "Location",
			LocationPlaceholder:   "Select location",
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
		Buttons: ButtonLabels{
			AddCostSchedule: "Add Cost Schedule",
		},
		Bulk: BulkLabels{Delete: "Delete"},
		Status: StatusLabels{
			Active:     "Active",
			Inactive:   "Inactive",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Empty: EmptyLabels{
			Title:   "No cost schedules yet",
			Message: "Add a cost schedule to group supplier cost plans by date range.",
		},
		Errors: ErrorLabels{
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
