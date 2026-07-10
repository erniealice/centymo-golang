package plan_group_plan

// ---------------------------------------------------------------------------
// Plan Group Plan labels
// ---------------------------------------------------------------------------

// Labels holds all labels for the plan_group_plan module.
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
	PlanGroupID   string `json:"plan_group_id"`
	PlanID        string `json:"plan_id"`
	SequenceOrder string `json:"sequence_order"`
	Status        string `json:"status"`
	DateCreated   string `json:"date_created"`
	Actions       string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels.
type FormLabels struct {
	SectionIdentity          string `json:"section_identity"`
	SectionOrdering          string `json:"section_ordering"`
	PlanGroupID              string `json:"plan_group_id"`
	PlanGroupIDPlaceholder   string `json:"plan_group_id_placeholder"`
	PlanGroupIDInfo          string `json:"plan_group_id_info"`
	PlanID                   string `json:"plan_id"`
	PlanIDPlaceholder        string `json:"plan_id_placeholder"`
	PlanIDInfo               string `json:"plan_id_info"`
	SequenceOrder            string `json:"sequence_order"`
	SequenceOrderPlaceholder string `json:"sequence_order_placeholder"`
	SequenceOrderInfo        string `json:"sequence_order_info"`
	Active                   string `json:"active"`
	ActiveInfo               string `json:"active_info"`
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
			Title:         "Plan Group Plans",
			Subtitle:      "Manage plans within plan groups",
			ActiveTitle:   "Active Plan Group Plans",
			InactiveTitle: "Inactive Plan Group Plans",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Plan Group Plan",
			Edit:       "Edit Plan Group Plan",
			Delete:     "Delete Plan Group Plan",
			BulkDelete: "Delete Plan Group Plans",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: ColumnLabels{
			PlanGroupID:   "Plan Group",
			PlanID:        "Plan",
			SequenceOrder: "Order",
			Status:        "Status",
			DateCreated:   "Date Created",
			Actions:       "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Plan Group Plans",
			Message: "No plan group plans to display.",
		},
		Form: FormLabels{
			SectionIdentity:          "Plan group plan details",
			SectionOrdering:          "Ordering",
			PlanGroupID:              "Plan Group",
			PlanGroupIDPlaceholder:   "Enter plan group ID",
			PlanGroupIDInfo:          "The plan group this plan belongs to.",
			PlanID:                   "Plan",
			PlanIDPlaceholder:        "Enter plan ID",
			PlanIDInfo:               "The plan assigned to this plan group.",
			SequenceOrder:            "Sequence order",
			SequenceOrderPlaceholder: "e.g. 1",
			SequenceOrderInfo:        "Optional display order within the plan group. Leave blank for no specific order.",
			Active:                   "Active",
			ActiveInfo:               "Inactive plan group plans are hidden from selection.",
		},
		Bulk: BulkLabels{
			DeleteTitle:       "Delete Plan Group Plans",
			DeleteMessage:     "Permanently delete the selected plan group plans? This cannot be undone.",
			ActivateTitle:     "Activate Plan Group Plans",
			ActivateMessage:   "Activate the selected plan group plans?",
			DeactivateTitle:   "Deactivate Plan Group Plans",
			DeactivateMessage: "Deactivate the selected plan group plans?",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:       "Delete Plan Group Plan",
			DeleteMessage:     "Permanently delete this plan group plan? This cannot be undone.",
			ActivateTitle:     "Activate Plan Group Plan",
			ActivateMessage:   "Activate this plan group plan?",
			DeactivateTitle:   "Deactivate Plan Group Plan",
			DeactivateMessage: "Deactivate this plan group plan?",
		},
		Tabs: TabLabels{
			Info: "Info",
		},
		Detail: DetailLabels{
			Title:        "Plan Group Plan",
			DateCreated:  "Date Created",
			DateModified: "Date Modified",
			NoSubtitle:   "No description provided",
		},
		Errors: ErrorLabels{
			NotFound:     "Plan group plan not found",
			LoadFailed:   "Failed to load plan group plan",
			Unauthorized: "You are not authorized to perform this action",
			CreateFailed: "Failed to create plan group plan",
			UpdateFailed: "Failed to update plan group plan",
			DeleteFailed: "Failed to delete plan group plan",
			InUse:        "This plan group plan is in use and cannot be deleted.",
		},
	}
}
