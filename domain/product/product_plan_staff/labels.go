package product_plan_staff

// ---------------------------------------------------------------------------
// Product Plan Staff labels
// ---------------------------------------------------------------------------

// Labels holds all labels for the product_plan_staff module.
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
	StaffID     string `json:"staffId"`
	ProductPlan string `json:"productPlan"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
	Actions     string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels for the staff-eligibility
// assignment. A product_plan_staff row assigns a staff member to a
// product_plan with a pool role (eligible / primary / assistant).
type FormLabels struct {
	SectionStaff       string `json:"sectionStaff"`
	SectionAssignment  string `json:"sectionAssignment"`
	StaffID            string `json:"staffId"`
	StaffIDPlaceholder string `json:"staffIdPlaceholder"`
	StaffIDInfo        string `json:"staffIdInfo"`
	ProductPlanID      string `json:"productPlanId"`
	ProductPlanPH      string `json:"productPlanPlaceholder"`
	ProductPlanInfo    string `json:"productPlanInfo"`
	Role               string `json:"role"`
	RolePlaceholder    string `json:"rolePlaceholder"`
	RoleInfo           string `json:"roleInfo"`
	Active             string `json:"active"`
	ActiveInfo         string `json:"activeInfo"`
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
	Title         string `json:"title"`
	DateCreated   string `json:"dateCreated"`
	DateModified  string `json:"dateModified"`
	NoProductPlan string `json:"noProductPlan"`
	NoRole        string `json:"noRole"`
	NoSubtitle    string `json:"noSubtitle"`
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

// DefaultLabels returns Labels with sensible English defaults using
// staff-eligibility vocabulary. Tiers override field names via lyngua.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:         "Plan Staff",
			Subtitle:      "Manage staff eligibility for product plans",
			ActiveTitle:   "Active Plan Staff",
			InactiveTitle: "Inactive Plan Staff",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Staff",
			Edit:       "Edit Staff",
			Delete:     "Delete Staff",
			BulkDelete: "Delete Staff",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: ColumnLabels{
			StaffID:     "Staff",
			ProductPlan: "Product Plan",
			Role:        "Role",
			Status:      "Status",
			DateCreated: "Date Created",
			Actions:     "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Plan Staff",
			Message: "No plan staff assignments to display.",
		},
		Form: FormLabels{
			SectionStaff:       "Staff member",
			SectionAssignment:  "Plan assignment",
			StaffID:            "Staff ID",
			StaffIDPlaceholder: "Enter staff ID",
			StaffIDInfo:        "The workspace_user ID of the eligible staff member.",
			ProductPlanID:      "Product Plan",
			ProductPlanPH:      "Enter product plan ID",
			ProductPlanInfo:    "The product_plan this staff member is eligible for.",
			Role:               "Role",
			RolePlaceholder:    "e.g. eligible, primary, assistant",
			RoleInfo:           "Pool membership role: eligible / primary / assistant (data-driven, not an enum).",
			Active:             "Active",
			ActiveInfo:         "Inactive assignments are excluded from the deliverer picker.",
		},
		Bulk: BulkLabels{
			DeleteTitle:       "Delete Plan Staff",
			DeleteMessage:     "Permanently delete the selected staff assignments? This cannot be undone.",
			ActivateTitle:     "Activate Plan Staff",
			ActivateMessage:   "Activate the selected staff assignments?",
			DeactivateTitle:   "Deactivate Plan Staff",
			DeactivateMessage: "Deactivate the selected staff assignments?",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:       "Delete Staff Assignment",
			DeleteMessage:     "Permanently delete this staff assignment? This cannot be undone.",
			ActivateTitle:     "Activate Staff Assignment",
			ActivateMessage:   "Activate {{name}}?",
			DeactivateTitle:   "Deactivate Staff Assignment",
			DeactivateMessage: "Deactivate {{name}}?",
		},
		Tabs: TabLabels{
			Info: "Info",
		},
		Detail: DetailLabels{
			Title:         "Plan Staff",
			DateCreated:   "Date Created",
			DateModified:  "Date Modified",
			NoProductPlan: "No product plan",
			NoRole:        "—",
			NoSubtitle:    "No role specified",
		},
		Errors: ErrorLabels{
			NotFound:     "Staff assignment not found",
			LoadFailed:   "Failed to load staff assignment",
			Unauthorized: "You are not authorized to perform this action",
			CreateFailed: "Failed to create staff assignment",
			UpdateFailed: "Failed to update staff assignment",
			DeleteFailed: "Failed to delete staff assignment",
			InUse:        "This staff assignment is in use and cannot be deleted.",
		},
	}
}
