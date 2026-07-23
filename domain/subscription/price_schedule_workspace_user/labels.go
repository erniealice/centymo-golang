package price_schedule_workspace_user

// ---------------------------------------------------------------------------
// PriceScheduleWorkspaceUser labels — operator access records at period nodes
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
	PriceScheduleId string `json:"price_schedule_id"`
	WorkspaceUserId string `json:"workspace_user_id"`
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

// FormLabels holds the drawer-form field labels. A price_schedule_workspace_user
// pins an operator at a period node.
type FormLabels struct {
	SectionCoordinator  string `json:"section_coordinator"`
	SectionAccess       string `json:"section_access"`
	PriceScheduleId     string `json:"price_schedule_id"`
	PriceScheduleIdPH   string `json:"price_schedule_id_placeholder"`
	PriceScheduleIdInfo string `json:"price_schedule_id_info"`
	WorkspaceUserId     string `json:"workspace_user_id"`
	WorkspaceUserIdPH   string `json:"workspace_user_id_placeholder"`
	WorkspaceUserIdInfo string `json:"workspace_user_id_info"`
	Capacity            string `json:"capacity"`
	CapacityInfo        string `json:"capacity_info"`
	CapacityPrimary     string `json:"capacity_primary"`
	CapacityAccess      string `json:"capacity_access"`
	IsOwner             string `json:"is_owner"`
	IsOwnerInfo         string `json:"is_owner_info"`
	Active              string `json:"active"`
	ActiveInfo          string `json:"active_info"`
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
	NoSchedule   string `json:"no_schedule"`
	NoUser       string `json:"no_user"`
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
			Capacity:        "Capacity",
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
			Capacity:            "Capacity",
			CapacityInfo:        "Whether this operator services the node (Primary) or only has view access (Access).",
			CapacityPrimary:     "Primary (servicer)",
			CapacityAccess:      "Access (view-only)",
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
		Title: TitleLabels{
			Label:        "Title",
			PrimaryOwner: "Owner",
			Primary:      "Servicer",
			Access:       "Member",
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
