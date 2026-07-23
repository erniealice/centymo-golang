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
	StaffID     string `json:"staff_id"`
	ProductPlan string `json:"product_plan"`
	Role        string `json:"role"`
	Status      string `json:"status"`
	DateCreated string `json:"date_created"`
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
	SectionStaff       string `json:"section_staff"`
	SectionAssignment  string `json:"section_assignment"`
	StaffID            string `json:"staff_id"`
	StaffIDPlaceholder string `json:"staff_id_placeholder"`
	StaffSearch        string `json:"staff_search"`
	StaffIDInfo        string `json:"staff_id_info"`
	ProductPlanID      string `json:"product_plan_id"`
	ProductPlanPH      string `json:"product_plan_placeholder"`
	ProductPlanSearch  string `json:"product_plan_search"`
	ProductPlanInfo    string `json:"product_plan_info"`
	Role               string `json:"role"`
	RolePlaceholder    string `json:"role_placeholder"`
	RoleInfo           string `json:"role_info"`
	Active             string `json:"active"`
	ActiveInfo         string `json:"active_info"`
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
	Title         string `json:"title"`
	DateCreated   string `json:"date_created"`
	DateModified  string `json:"date_modified"`
	NoProductPlan string `json:"no_product_plan"`
	NoRole        string `json:"no_role"`
	NoSubtitle    string `json:"no_subtitle"`
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
			StaffID:            "Staff member",
			StaffIDPlaceholder: "Select a staff member...",
			StaffSearch:        "Filter...",
			StaffIDInfo:        "The eligible staff member.",
			ProductPlanID:      "Product Plan",
			ProductPlanPH:      "Select a product plan...",
			ProductPlanSearch:  "Filter...",
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
