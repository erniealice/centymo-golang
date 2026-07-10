package subscription_group_product_plan_staff

// ---------------------------------------------------------------------------
// Subscription Group Product Plan Staff (class-edge) labels
// ---------------------------------------------------------------------------

// Labels holds all labels for the subscription_group_product_plan_staff module.
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
	SubscriptionGroupID string `json:"subscription_group_id"`
	ProductPlanID       string `json:"product_plan_id"`
	StaffID             string `json:"staff_id"`
	Role                string `json:"role"`
	Status              string `json:"status"`
	DateCreated         string `json:"date_created"`
	Actions             string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels.
// Vocabulary: class-edge — (section × subject × teacher × role).
type FormLabels struct {
	SectionAssignment       string `json:"section_assignment"`
	SubscriptionGroupID     string `json:"subscription_group_id"`
	SubscriptionGroupPH     string `json:"subscription_group_placeholder"`
	SubscriptionGroupSearch string `json:"subscription_group_search"`
	SubscriptionGroupInfo   string `json:"subscription_group_info"`
	ProductPlanID           string `json:"product_plan_id"`
	ProductPlanPH           string `json:"product_plan_placeholder"`
	ProductPlanSearch       string `json:"product_plan_search"`
	ProductPlanInfo         string `json:"product_plan_info"`
	StaffID                 string `json:"staff_id"`
	StaffPH                 string `json:"staff_placeholder"`
	StaffSearch             string `json:"staff_search"`
	StaffInfo               string `json:"staff_info"`
	Role                    string `json:"role"`
	RolePlaceholder         string `json:"role_placeholder"`
	RoleInfo                string `json:"role_info"`
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
	NoPlan       string `json:"no_plan"`
	NoStaff      string `json:"no_staff"`
	NoRole       string `json:"no_role"`
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

// DefaultLabels returns Labels with sensible English defaults using the
// class-edge vocabulary. Tiers override field names via lyngua.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:         "Class Assignments",
			Subtitle:      "Manage staff assignments to section-subject pairs",
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
			SubscriptionGroupID: "Section",
			ProductPlanID:       "Subject",
			StaffID:             "Staff",
			Role:                "Role",
			Status:              "Status",
			DateCreated:         "Date Created",
			Actions:             "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Assignments",
			Message: "No class assignments to display.",
		},
		Form: FormLabels{
			SectionAssignment:       "Assignment details",
			SubscriptionGroupID:     "Section",
			SubscriptionGroupPH:     "Select a section...",
			SubscriptionGroupSearch: "Filter...",
			SubscriptionGroupInfo:   "The section (cohort) this staff member is assigned to.",
			ProductPlanID:           "Subject",
			ProductPlanPH:           "Select a subject...",
			ProductPlanSearch:       "Filter...",
			ProductPlanInfo:         "The subject (product plan) delivered in this section.",
			StaffID:                 "Staff member",
			StaffPH:                 "Select a staff member...",
			StaffSearch:             "Filter...",
			StaffInfo:               "The staff member delivering this subject in this section.",
			Role:                    "Role",
			RolePlaceholder:         "e.g. teacher, co-teacher, tutor",
			RoleInfo:                "The role this staff member holds in this class assignment.",
			Active:                  "Active",
			ActiveInfo:              "Inactive assignments are excluded from grade-sheet scoping.",
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
			Title:        "Class Assignment",
			DateCreated:  "Date Created",
			DateModified: "Date Modified",
			NoGroup:      "No section",
			NoPlan:       "No subject",
			NoStaff:      "No staff",
			NoRole:       "—",
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
