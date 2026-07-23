package subscription_group_workspace_user

// Labels holds all labels for the subscription_group_workspace_user module.
// Vocabulary: an assignment of a workspace user (operator) to a subscription
// group (cohort) with a servicing scope and role.
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
	WorkspaceUser     string `json:"workspace_user"`
	SubscriptionGroup string `json:"subscription_group"`
	Scope             string `json:"scope"`
	Role              string `json:"role"`
	IsOwner           string `json:"is_owner"`
	Status            string `json:"status"`
	DateCreated       string `json:"date_created"`
	Actions           string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels for the operator-assignment
// entity (subscription_group_workspace_user).
type FormLabels struct {
	SectionAssignment string `json:"section_assignment"`
	SectionServicing  string `json:"section_servicing"`

	WorkspaceUserId         string `json:"workspace_user_id"`
	WorkspaceUserIdPH       string `json:"workspace_user_id_placeholder"`
	WorkspaceUserSearch     string `json:"workspace_user_search"`
	WorkspaceUserIdInfo     string `json:"workspace_user_id_info"`
	SubscriptionGroupId     string `json:"subscription_group_id"`
	SubscriptionGroupIdPH   string `json:"subscription_group_id_placeholder"`
	SubscriptionGroupSearch string `json:"subscription_group_search"`
	SubscriptionGroupIdInfo string `json:"subscription_group_id_info"`
	Scope                   string `json:"scope"`
	ScopePlaceholder        string `json:"scope_placeholder"`
	ScopeInfo               string `json:"scope_info"`
	Role                    string `json:"role"`
	RolePlaceholder         string `json:"role_placeholder"`
	RoleInfo                string `json:"role_info"`
	IsOwner                 string `json:"is_owner"`
	IsOwnerInfo             string `json:"is_owner_info"`
	Active                  string `json:"active"`
	ActiveInfo              string `json:"active_info"`
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
	NoGroup      string `json:"no_group"`
	NoUser       string `json:"no_user"`
	NoScope      string `json:"no_scope"`
	NoRole       string `json:"no_role"`
	OwnerYes     string `json:"owner_yes"`
	OwnerNo      string `json:"owner_no"`
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
			Title:         "Group Assignments",
			Subtitle:      "Manage operator assignments to subscription groups",
			ActiveTitle:   "Active Assignments",
			InactiveTitle: "Inactive Assignments",
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
			WorkspaceUser:     "Workspace User",
			SubscriptionGroup: "Group",
			Scope:             "Scope",
			Role:              "Role",
			IsOwner:           "Owner",
			Status:            "Status",
			DateCreated:       "Date Created",
			Actions:           "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Assignments",
			Message: "No operator assignments to display.",
		},
		Form: FormLabels{
			SectionAssignment:       "Assignment",
			SectionServicing:        "Servicing",
			WorkspaceUserId:         "Workspace User",
			WorkspaceUserIdPH:       "Select a workspace user...",
			WorkspaceUserSearch:     "Filter...",
			WorkspaceUserIdInfo:     "The operator (workspace user) being assigned to this group.",
			SubscriptionGroupId:     "Subscription Group",
			SubscriptionGroupIdPH:   "Select a group...",
			SubscriptionGroupSearch: "Filter...",
			SubscriptionGroupIdInfo: "The subscription group (cohort) this operator is assigned to.",
			Scope:                   "Scope",
			ScopePlaceholder:        "e.g. coordinator",
			ScopeInfo:               "The servicing scope of this assignment (e.g. coordinator, adviser).",
			Role:                    "Role",
			RolePlaceholder:         "e.g. lead",
			RoleInfo:                "The role within the assigned scope.",
			IsOwner:                 "Owner",
			IsOwnerInfo:             "Marks this operator as the group owner.",
			Active:                  "Active",
			ActiveInfo:              "Inactive assignments are excluded from group-level servicing.",
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
			Title:        "Assignment",
			DateCreated:  "Date Created",
			DateModified: "Date Modified",
			NoGroup:      "No group",
			NoUser:       "No user",
			NoScope:      "—",
			NoRole:       "—",
			OwnerYes:     "Yes",
			OwnerNo:      "No",
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
