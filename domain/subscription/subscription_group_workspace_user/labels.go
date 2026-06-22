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
	WorkspaceUser     string `json:"workspaceUser"`
	SubscriptionGroup string `json:"subscriptionGroup"`
	Scope             string `json:"scope"`
	Role              string `json:"role"`
	IsOwner           string `json:"isOwner"`
	Status            string `json:"status"`
	DateCreated       string `json:"dateCreated"`
	Actions           string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels for the operator-assignment
// entity (subscription_group_workspace_user).
type FormLabels struct {
	SectionAssignment string `json:"sectionAssignment"`
	SectionServicing  string `json:"sectionServicing"`

	WorkspaceUserId         string `json:"workspaceUserId"`
	WorkspaceUserIdPH       string `json:"workspaceUserIdPlaceholder"`
	WorkspaceUserIdInfo     string `json:"workspaceUserIdInfo"`
	SubscriptionGroupId     string `json:"subscriptionGroupId"`
	SubscriptionGroupIdPH   string `json:"subscriptionGroupIdPlaceholder"`
	SubscriptionGroupIdInfo string `json:"subscriptionGroupIdInfo"`
	Scope                   string `json:"scope"`
	ScopePlaceholder        string `json:"scopePlaceholder"`
	ScopeInfo               string `json:"scopeInfo"`
	Role                    string `json:"role"`
	RolePlaceholder         string `json:"rolePlaceholder"`
	RoleInfo                string `json:"roleInfo"`
	IsOwner                 string `json:"isOwner"`
	IsOwnerInfo             string `json:"isOwnerInfo"`
	Active                  string `json:"active"`
	ActiveInfo              string `json:"activeInfo"`
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
	NoGroup      string `json:"noGroup"`
	NoUser       string `json:"noUser"`
	NoScope      string `json:"noScope"`
	NoRole       string `json:"noRole"`
	OwnerYes     string `json:"ownerYes"`
	OwnerNo      string `json:"ownerNo"`
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
			WorkspaceUserIdPH:       "Enter workspace user ID",
			WorkspaceUserIdInfo:     "The operator (workspace user) being assigned to this group.",
			SubscriptionGroupId:     "Subscription Group",
			SubscriptionGroupIdPH:   "Enter group ID",
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
