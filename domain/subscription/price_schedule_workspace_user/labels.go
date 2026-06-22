package price_schedule_workspace_user

// ---------------------------------------------------------------------------
// PriceScheduleWorkspaceUser labels — "year coordinator" access records
// ---------------------------------------------------------------------------

// Labels holds all labels for the price_schedule_workspace_user module.
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
	PriceScheduleId string `json:"priceScheduleId"`
	WorkspaceUserId string `json:"workspaceUserId"`
	Scope           string `json:"scope"`
	Role            string `json:"role"`
	IsOwner         string `json:"isOwner"`
	Status          string `json:"status"`
	DateCreated     string `json:"dateCreated"`
	Actions         string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels. Year-coordinator vocabulary —
// a price_schedule_workspace_user pins an operator at a period node.
type FormLabels struct {
	SectionCoordinator  string `json:"sectionCoordinator"`
	SectionAccess       string `json:"sectionAccess"`
	PriceScheduleId     string `json:"priceScheduleId"`
	PriceScheduleIdPH   string `json:"priceScheduleIdPlaceholder"`
	PriceScheduleIdInfo string `json:"priceScheduleIdInfo"`
	WorkspaceUserId     string `json:"workspaceUserId"`
	WorkspaceUserIdPH   string `json:"workspaceUserIdPlaceholder"`
	WorkspaceUserIdInfo string `json:"workspaceUserIdInfo"`
	Scope               string `json:"scope"`
	ScopePH             string `json:"scopePlaceholder"`
	ScopeInfo           string `json:"scopeInfo"`
	Role                string `json:"role"`
	RolePH              string `json:"rolePlaceholder"`
	RoleInfo            string `json:"roleInfo"`
	IsOwner             string `json:"isOwner"`
	IsOwnerInfo         string `json:"isOwnerInfo"`
	Active              string `json:"active"`
	ActiveInfo          string `json:"activeInfo"`
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
	Title        string `json:"title"`
	DateCreated  string `json:"dateCreated"`
	DateModified string `json:"dateModified"`
	NoSchedule   string `json:"noSchedule"`
	NoUser       string `json:"noUser"`
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
// year-coordinator vocabulary.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:         "Period Coordinators",
			Subtitle:      "Manage operator access at period (price schedule) level",
			ActiveTitle:   "Active Period Coordinators",
			InactiveTitle: "Inactive Period Coordinators",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Coordinator",
			Edit:       "Edit Coordinator",
			Delete:     "Delete Coordinator",
			BulkDelete: "Delete Coordinators",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: ColumnLabels{
			PriceScheduleId: "Period",
			WorkspaceUserId: "Operator",
			Scope:           "Scope",
			Role:            "Role",
			IsOwner:         "Owner",
			Status:          "Status",
			DateCreated:     "Date Created",
			Actions:         "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Period Coordinators",
			Message: "No period coordinator records to display.",
		},
		Form: FormLabels{
			SectionCoordinator:  "Coordinator details",
			SectionAccess:       "Access settings",
			PriceScheduleId:     "Period",
			PriceScheduleIdPH:   "Enter period ID",
			PriceScheduleIdInfo: "The price schedule (period / academic year) this coordinator is pinned to.",
			WorkspaceUserId:     "Operator",
			WorkspaceUserIdPH:   "Enter operator ID",
			WorkspaceUserIdInfo: "The workspace user who acts as coordinator for this period.",
			Scope:               "Scope",
			ScopePH:             "e.g. academic_year",
			ScopeInfo:           "Optional scope tag narrowing the coordinator's visibility within the period.",
			Role:                "Role",
			RolePH:              "e.g. coordinator",
			RoleInfo:            "Role label for this coordinator (free-text).",
			IsOwner:             "Owner",
			IsOwnerInfo:         "Mark this operator as the primary owner of the period.",
			Active:              "Active",
			ActiveInfo:          "Inactive coordinators are hidden from period-level views.",
		},
		Bulk: BulkLabels{
			DeleteTitle:       "Delete Period Coordinators",
			DeleteMessage:     "Permanently delete the selected coordinator records? This cannot be undone.",
			ActivateTitle:     "Activate Period Coordinators",
			ActivateMessage:   "Activate the selected coordinator records?",
			DeactivateTitle:   "Deactivate Period Coordinators",
			DeactivateMessage: "Deactivate the selected coordinator records?",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:       "Delete Coordinator",
			DeleteMessage:     "Permanently delete this coordinator record? This cannot be undone.",
			ActivateTitle:     "Activate Coordinator",
			ActivateMessage:   "Activate {{id}}?",
			DeactivateTitle:   "Deactivate Coordinator",
			DeactivateMessage: "Deactivate {{id}}?",
		},
		Tabs: TabLabels{
			Info: "Info",
		},
		Detail: DetailLabels{
			Title:        "Period Coordinator",
			DateCreated:  "Date Created",
			DateModified: "Date Modified",
			NoSchedule:   "No period",
			NoUser:       "No operator",
		},
		Errors: ErrorLabels{
			NotFound:     "Coordinator record not found",
			LoadFailed:   "Failed to load coordinator record",
			Unauthorized: "You are not authorized to perform this action",
			CreateFailed: "Failed to create coordinator record",
			UpdateFailed: "Failed to update coordinator record",
			DeleteFailed: "Failed to delete coordinator record",
			InUse:        "This coordinator record is in use and cannot be deleted.",
		},
	}
}
