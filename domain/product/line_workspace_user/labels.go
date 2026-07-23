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
	Title   TitleLabels   `json:"title"`
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
	Capacity        string `json:"capacity"`
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
	SectionServicing  string `json:"section_servicing"`

	WorkspaceUserId            string `json:"workspace_user_id"`
	WorkspaceUserIdPlaceholder string `json:"workspace_user_id_placeholder"`
	WorkspaceUserIdInfo        string `json:"workspace_user_id_info"`

	LineId            string `json:"line_id"`
	LineIdPlaceholder string `json:"line_id_placeholder"`
	LineIdInfo        string `json:"line_id_info"`

	Capacity        string `json:"capacity"`
	CapacityInfo    string `json:"capacity_info"`
	CapacityPrimary string `json:"capacity_primary"`
	CapacityAccess  string `json:"capacity_access"`

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

// TitleLabels holds the servicing-title strings derived from the generic
// (capacity, is_owner) axes. Label is the field/row caption; the other three
// map to the generic keys {primary_owner, primary, access}. Vertical titles
// are supplied only by the lyngua businessType overlay, never by this code.
type TitleLabels struct {
	Label        string `json:"label"`
	PrimaryOwner string `json:"primary_owner"`
	Primary      string `json:"primary"`
	Access       string `json:"access"`
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
			Capacity:        "Capacity",
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
			SectionServicing:  "Servicing",

			WorkspaceUserId:            "Workspace User ID",
			WorkspaceUserIdPlaceholder: "Enter workspace user ID",
			WorkspaceUserIdInfo:        "The operator (workspace user) being assigned to the line node.",

			LineId:            "Line ID",
			LineIdPlaceholder: "Enter line ID",
			LineIdInfo:        "The line node this operator is pinned to.",

			Capacity:        "Capacity",
			CapacityInfo:    "Whether this operator services the node (Primary) or only has view access (Access).",
			CapacityPrimary: "Primary (servicer)",
			CapacityAccess:  "Access (view-only)",

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
		Title: TitleLabels{
			Label:        "Title",
			PrimaryOwner: "Owner",
			Primary:      "Servicer",
			Access:       "Member",
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

// CapacityLabel returns the human-readable label for a generic capacity token.
// The domain is closed ({primary, access}); any other value (incl. "") floors
// to the least-privilege access label.
func CapacityLabel(l FormLabels, capacity string) string {
	if capacity == "primary" {
		return l.CapacityPrimary
	}
	return l.CapacityAccess
}

// DeriveTitle maps the generic (capacity, is_owner) axes to a servicing title.
// The key is generic ({primary_owner, primary, access}); the vertical wording
// is supplied by the lyngua overlay, never by this code.
func DeriveTitle(l TitleLabels, capacity string, isOwner bool) string {
	if capacity == "primary" {
		if isOwner {
			return l.PrimaryOwner
		}
		return l.Primary
	}
	return l.Access
}
