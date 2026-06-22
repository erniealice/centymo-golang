package subscription_group

// ---------------------------------------------------------------------------
// Subscription Group (section / cohort) labels
// ---------------------------------------------------------------------------

// Labels holds all labels for the subscription_group module.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Buttons ButtonLabels  `json:"buttons"`
	Columns ColumnLabels  `json:"columns"`
	Empty   EmptyLabels   `json:"empty"`
	Form    FormLabels    `json:"form"`
	Bulk    BulkLabels    `json:"bulk"`
	Confirm ConfirmLabels `json:"confirm"`
	Tabs    TabLabels     `json:"tabs"`
	Detail  DetailLabels  `json:"detail"`
	Errors  ErrorLabels   `json:"errors"`
}

type PageLabels struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	ActiveTitle   string `json:"activeTitle"`
	InactiveTitle string `json:"inactiveTitle"`
}

type ButtonLabels struct {
	View       string `json:"view"`
	Add        string `json:"add"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	BulkDelete string `json:"bulkDelete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Capacity    string `json:"capacity"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
	Actions     string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels. Section/cohort vocabulary —
// a subscription_group is a class roster / patient panel / project team
// anchored to a program (plan) and a period (price_schedule).
type FormLabels struct {
	SectionIdentity   string `json:"sectionIdentity"`
	SectionAnchors    string `json:"sectionAnchors"`
	SectionCapacity   string `json:"sectionCapacity"`
	Name              string `json:"name"`
	NamePlaceholder   string `json:"namePlaceholder"`
	NameInfo          string `json:"nameInfo"`
	Kind              string `json:"kind"`
	KindPlaceholder   string `json:"kindPlaceholder"`
	KindInfo          string `json:"kindInfo"`
	KindCohort        string `json:"kindCohort"`
	KindRoster        string `json:"kindRoster"`
	KindPanel         string `json:"kindPanel"`
	KindProjectTeam   string `json:"kindProjectTeam"`
	Plan              string `json:"plan"`
	PlanPlaceholder   string `json:"planPlaceholder"`
	PlanSearch        string `json:"planSearch"`
	PlanInfo          string `json:"planInfo"`
	PriceSchedule     string `json:"priceSchedule"`
	PriceSchedulePH   string `json:"priceSchedulePlaceholder"`
	PriceScheduleSrch string `json:"priceScheduleSearch"`
	PriceScheduleInfo string `json:"priceScheduleInfo"`
	CapacityMode      string `json:"capacityMode"`
	CapacityModeInfo  string `json:"capacityModeInfo"`
	CapUnlimited      string `json:"capUnlimited"`
	CapClosed         string `json:"capClosed"`
	CapCapped         string `json:"capCapped"`
	MaxCapacity       string `json:"maxCapacity"`
	MaxCapacityPH     string `json:"maxCapacityPlaceholder"`
	MaxCapacityInfo   string `json:"maxCapacityInfo"`
	Active            string `json:"active"`
	ActiveInfo        string `json:"activeInfo"`
}

type BulkLabels struct {
	DeleteTitle       string `json:"deleteTitle"`
	DeleteMessage     string `json:"deleteMessage"`
	ActivateTitle     string `json:"activateTitle"`
	ActivateMessage   string `json:"activateMessage"`
	DeactivateTitle   string `json:"deactivateTitle"`
	DeactivateMessage string `json:"deactivateMessage"`
}

type ConfirmLabels struct {
	DeleteTitle       string `json:"deleteTitle"`
	DeleteMessage     string `json:"deleteMessage"`
	ActivateTitle     string `json:"activateTitle"`
	ActivateMessage   string `json:"activateMessage"`
	DeactivateTitle   string `json:"deactivateTitle"`
	DeactivateMessage string `json:"deactivateMessage"`
}

type TabLabels struct {
	Info string `json:"info"`
}

type DetailLabels struct {
	Title          string `json:"title"`
	DateCreated    string `json:"dateCreated"`
	DateModified   string `json:"dateModified"`
	NoPlan         string `json:"noPlan"`
	NoSchedule     string `json:"noSchedule"`
	NoKind         string `json:"noKind"`
	CapacityValue  string `json:"capacityValue"`  // e.g. "%d seats" (CAPPED)
	CapacityModeNF string `json:"capacityModeNF"` // fallback when mode unspecified
	NoSubtitle     string `json:"noSubtitle"`
}

type ErrorLabels struct {
	NotFound     string `json:"notFound"`
	LoadFailed   string `json:"loadFailed"`
	Unauthorized string `json:"unauthorized"`
	CreateFailed string `json:"createFailed"`
	UpdateFailed string `json:"updateFailed"`
	DeleteFailed string `json:"deleteFailed"`
	InUse        string `json:"inUse"`
}

// DefaultLabels returns Labels with sensible English defaults using the
// section/cohort vocabulary. Tiers override field names via lyngua.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:         "Sections",
			Subtitle:      "Manage your sections and cohorts",
			ActiveTitle:   "Active Sections",
			InactiveTitle: "Inactive Sections",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Section",
			Edit:       "Edit Section",
			Delete:     "Delete Section",
			BulkDelete: "Delete Sections",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: ColumnLabels{
			Name:        "Name",
			Kind:        "Type",
			Capacity:    "Capacity",
			Status:      "Status",
			DateCreated: "Date Created",
			Actions:     "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Sections",
			Message: "No sections to display.",
		},
		Form: FormLabels{
			SectionIdentity:   "Section details",
			SectionAnchors:    "Program & period",
			SectionCapacity:   "Capacity",
			Name:              "Name",
			NamePlaceholder:   "Enter section name",
			NameInfo:          "A short display name for this section or cohort.",
			Kind:              "Type",
			KindPlaceholder:   "Select a type...",
			KindInfo:          "The shape of this cohort: a class roster, a panel, or a project team.",
			KindCohort:        "Cohort",
			KindRoster:        "Roster",
			KindPanel:         "Panel",
			KindProjectTeam:   "Project team",
			Plan:              "Program",
			PlanPlaceholder:   "Select a program...",
			PlanSearch:        "Filter...",
			PlanInfo:          "The program (plan) this section realizes. Section identity is program × period.",
			PriceSchedule:     "Period",
			PriceSchedulePH:   "Select a period...",
			PriceScheduleSrch: "Filter...",
			PriceScheduleInfo: "The billing period this section is anchored to (e.g. the academic-year price schedule).",
			CapacityMode:      "Capacity mode",
			CapacityModeInfo:  "Unlimited admits everyone; Closed admits no one; Capped enforces the seat limit below.",
			CapUnlimited:      "Unlimited",
			CapClosed:         "Closed",
			CapCapped:         "Capped",
			MaxCapacity:       "Maximum seats",
			MaxCapacityPH:     "e.g. 30",
			MaxCapacityInfo:   "Read only when the capacity mode is Capped.",
			Active:            "Active",
			ActiveInfo:        "Inactive sections are hidden from new enrollments.",
		},
		Bulk: BulkLabels{
			DeleteTitle:       "Delete Sections",
			DeleteMessage:     "Permanently delete the selected sections? This cannot be undone.",
			ActivateTitle:     "Activate Sections",
			ActivateMessage:   "Activate the selected sections?",
			DeactivateTitle:   "Deactivate Sections",
			DeactivateMessage: "Deactivate the selected sections?",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:       "Delete Section",
			DeleteMessage:     "Permanently delete this section? This cannot be undone.",
			ActivateTitle:     "Activate Section",
			ActivateMessage:   "Activate {{name}}?",
			DeactivateTitle:   "Deactivate Section",
			DeactivateMessage: "Deactivate {{name}}?",
		},
		Tabs: TabLabels{
			Info: "Info",
		},
		Detail: DetailLabels{
			Title:          "Section",
			DateCreated:    "Date Created",
			DateModified:   "Date Modified",
			NoPlan:         "No program",
			NoSchedule:     "No period",
			NoKind:         "—",
			CapacityValue:  "%d seats",
			CapacityModeNF: "Unlimited",
			NoSubtitle:     "No description provided",
		},
		Errors: ErrorLabels{
			NotFound:     "Section not found",
			LoadFailed:   "Failed to load section",
			Unauthorized: "You are not authorized to perform this action",
			CreateFailed: "Failed to create section",
			UpdateFailed: "Failed to update section",
			DeleteFailed: "Failed to delete section",
			InUse:        "This section is in use and cannot be deleted.",
		},
	}
}
