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
	SubscriptionGroupID string `json:"subscriptionGroupId"`
	ProductPlanID       string `json:"productPlanId"`
	StaffID             string `json:"staffId"`
	Role                string `json:"role"`
	Status              string `json:"status"`
	DateCreated         string `json:"dateCreated"`
	Actions             string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels.
// Vocabulary: class-edge — (section × subject × teacher × role).
type FormLabels struct {
	SectionAssignment       string `json:"sectionAssignment"`
	SubscriptionGroupID     string `json:"subscriptionGroupId"`
	SubscriptionGroupPH     string `json:"subscriptionGroupPlaceholder"`
	SubscriptionGroupSearch string `json:"subscriptionGroupSearch"`
	SubscriptionGroupInfo   string `json:"subscriptionGroupInfo"`
	ProductPlanID           string `json:"productPlanId"`
	ProductPlanPH           string `json:"productPlanPlaceholder"`
	ProductPlanSearch       string `json:"productPlanSearch"`
	ProductPlanInfo         string `json:"productPlanInfo"`
	StaffID                 string `json:"staffId"`
	StaffPH                 string `json:"staffPlaceholder"`
	StaffSearch             string `json:"staffSearch"`
	StaffInfo               string `json:"staffInfo"`
	Role                    string `json:"role"`
	RolePlaceholder         string `json:"rolePlaceholder"`
	RoleInfo                string `json:"roleInfo"`
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
	NoPlan       string `json:"noPlan"`
	NoStaff      string `json:"noStaff"`
	NoRole       string `json:"noRole"`
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
