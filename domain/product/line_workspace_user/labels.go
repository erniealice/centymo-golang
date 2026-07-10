package line_workspace_user

// ---------------------------------------------------------------------------
// Line Workspace User labels
// ---------------------------------------------------------------------------

// Labels holds all labels for the line_workspace_user module.
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
	ActiveTitle   string `json:"active_title"`
	InactiveTitle string `json:"inactive_title"`
}

type ButtonLabels struct {
	View       string `json:"view"`
	Add        string `json:"add"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	BulkDelete string `json:"bulk_delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ColumnLabels struct {
	WorkspaceUserId string `json:"workspace_user_id"`
	LineId          string `json:"line_id"`
	Scope           string `json:"scope"`
	Role            string `json:"role"`
	IsOwner         string `json:"is_owner"`
	Status          string `json:"status"`
	DateCreated     string `json:"date_created"`
	Actions         string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels for line_workspace_user.
type FormLabels struct {
	SectionAssignment string `json:"section_assignment"`
	SectionRole       string `json:"section_role"`

	WorkspaceUserId            string `json:"workspace_user_id"`
	WorkspaceUserIdPlaceholder string `json:"workspace_user_id_placeholder"`
	WorkspaceUserIdInfo        string `json:"workspace_user_id_info"`

	LineId            string `json:"line_id"`
	LineIdPlaceholder string `json:"line_id_placeholder"`
	LineIdInfo        string `json:"line_id_info"`

	Scope            string `json:"scope"`
	ScopePlaceholder string `json:"scope_placeholder"`
	ScopeInfo        string `json:"scope_info"`

	Role            string `json:"role"`
	RolePlaceholder string `json:"role_placeholder"`
	RoleInfo        string `json:"role_info"`

	IsOwner     string `json:"is_owner"`
	IsOwnerInfo string `json:"is_owner_info"`

	Active     string `json:"active"`
	ActiveInfo string `json:"active_info"`
}

type BulkLabels struct {
	DeleteTitle       string `json:"delete_title"`
	DeleteMessage     string `json:"delete_message"`
	ActivateTitle     string `json:"activate_title"`
	ActivateMessage   string `json:"activate_message"`
	DeactivateTitle   string `json:"deactivate_title"`
	DeactivateMessage string `json:"deactivate_message"`
}

type ConfirmLabels struct {
	DeleteTitle       string `json:"delete_title"`
	DeleteMessage     string `json:"delete_message"`
	ActivateTitle     string `json:"activate_title"`
	ActivateMessage   string `json:"activate_message"`
	DeactivateTitle   string `json:"deactivate_title"`
	DeactivateMessage string `json:"deactivate_message"`
}

type TabLabels struct {
	Info string `json:"info"`
}

type DetailLabels struct {
	Title        string `json:"title"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
	NoSubtitle   string `json:"no_subtitle"`
}

type ErrorLabels struct {
	NotFound     string `json:"not_found"`
	LoadFailed   string `json:"load_failed"`
	Unauthorized string `json:"unauthorized"`
	CreateFailed string `json:"create_failed"`
	UpdateFailed string `json:"update_failed"`
	DeleteFailed string `json:"delete_failed"`
	InUse        string `json:"in_use"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:         "Line Workspace Users",
			Subtitle:      "Manage operator assignments at line nodes",
			ActiveTitle:   "Active Line Workspace Users",
			InactiveTitle: "Inactive Line Workspace Users",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Assignment",
			Edit:       "Edit Assignment",
			Delete:     "Delete Assignment",
			BulkDelete: "Delete Assignments",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: ColumnLabels{
			WorkspaceUserId: "Workspace User",
			LineId:          "Line",
			Scope:           "Scope",
			Role:            "Role",
			IsOwner:         "Owner",
			Status:          "Status",
			DateCreated:     "Date Created",
			Actions:         "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Assignments",
			Message: "No line workspace user assignments to display.",
		},
		Form: FormLabels{
			SectionAssignment: "Assignment details",
			SectionRole:       "Role & scope",

			WorkspaceUserId:            "Workspace User ID",
			WorkspaceUserIdPlaceholder: "Enter workspace user ID",
			WorkspaceUserIdInfo:        "The operator (workspace user) being assigned to the line node.",

			LineId:            "Line ID",
			LineIdPlaceholder: "Enter line ID",
			LineIdInfo:        "The line node this operator is pinned to.",

			Scope:            "Scope",
			ScopePlaceholder: "e.g. coordinator, adviser",
			ScopeInfo:        "The visibility scope discriminator for this assignment.",

			Role:            "Role",
			RolePlaceholder: "e.g. coordinator, adviser",
			RoleInfo:        "The servicing role discriminator for this assignment.",

			IsOwner:     "Owner",
			IsOwnerInfo: "Mark this user as the owner of the line node.",

			Active:     "Active",
			ActiveInfo: "Inactive assignments are hidden from tier-2 group visibility.",
		},
		Bulk: BulkLabels{
			DeleteTitle:       "Delete Assignments",
			DeleteMessage:     "Permanently delete the selected assignments? This cannot be undone.",
			ActivateTitle:     "Activate Assignments",
			ActivateMessage:   "Activate the selected assignments?",
			DeactivateTitle:   "Deactivate Assignments",
			DeactivateMessage: "Deactivate the selected assignments?",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:       "Delete Assignment",
			DeleteMessage:     "Permanently delete this assignment? This cannot be undone.",
			ActivateTitle:     "Activate Assignment",
			ActivateMessage:   "Activate {{name}}?",
			DeactivateTitle:   "Deactivate Assignment",
			DeactivateMessage: "Deactivate {{name}}?",
		},
		Tabs: TabLabels{
			Info: "Info",
		},
		Detail: DetailLabels{
			Title:        "Line Workspace User",
			DateCreated:  "Date Created",
			DateModified: "Date Modified",
			NoSubtitle:   "No description provided",
		},
		Errors: ErrorLabels{
			NotFound:     "Assignment not found",
			LoadFailed:   "Failed to load assignment",
			Unauthorized: "You are not authorized to perform this action",
			CreateFailed: "Failed to create assignment",
			UpdateFailed: "Failed to update assignment",
			DeleteFailed: "Failed to delete assignment",
			InUse:        "This assignment is in use and cannot be deleted.",
		},
	}
}
