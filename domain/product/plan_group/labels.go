package plan_group

// ---------------------------------------------------------------------------
// Plan Group labels
// ---------------------------------------------------------------------------

// Labels holds all labels for the plan_group module.
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
	Name        string `json:"name"`
	Code        string `json:"code"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
	Actions     string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels. Plan group vocabulary —
// a plan_group is a stable taxonomy node that groups plans across periods.
type FormLabels struct {
	SectionIdentity  string `json:"sectionIdentity"`
	SectionHierarchy string `json:"sectionHierarchy"`
	Name             string `json:"name"`
	NamePlaceholder  string `json:"namePlaceholder"`
	NameInfo         string `json:"nameInfo"`
	Code             string `json:"code"`
	CodePlaceholder  string `json:"codePlaceholder"`
	CodeInfo         string `json:"codeInfo"`
	ParentGroup      string `json:"parentGroup"`
	ParentGroupPH    string `json:"parentGroupPlaceholder"`
	ParentGroupSrch  string `json:"parentGroupSearch"`
	ParentGroupInfo  string `json:"parentGroupInfo"`
	Active           string `json:"active"`
	ActiveInfo       string `json:"activeInfo"`
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
	NoParent     string `json:"noParent"`
	NoCode       string `json:"noCode"`
	NoSubtitle   string `json:"noSubtitle"`
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

// DefaultLabels returns Labels with sensible English defaults using plan group
// vocabulary. Tiers override field names via lyngua.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:         "Plan Groups",
			Subtitle:      "Manage your plan group taxonomy",
			ActiveTitle:   "Active Plan Groups",
			InactiveTitle: "Inactive Plan Groups",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Plan Group",
			Edit:       "Edit Plan Group",
			Delete:     "Delete Plan Group",
			BulkDelete: "Delete Plan Groups",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: ColumnLabels{
			Name:        "Name",
			Code:        "Code",
			Status:      "Status",
			DateCreated: "Date Created",
			Actions:     "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Plan Groups",
			Message: "No plan groups to display.",
		},
		Form: FormLabels{
			SectionIdentity:  "Group details",
			SectionHierarchy: "Hierarchy",
			Name:             "Name",
			NamePlaceholder:  "Enter plan group name",
			NameInfo:         "A short display name for this plan group (e.g. \"Junior High\").",
			Code:             "Code",
			CodePlaceholder:  "Enter a short code (optional)",
			CodeInfo:         "Optional short identifier for this group used in reporting.",
			ParentGroup:      "Parent group",
			ParentGroupPH:    "Select a parent group...",
			ParentGroupSrch:  "Filter...",
			ParentGroupInfo:  "Optional parent in the plan group hierarchy. Leave unset for a top-level group.",
			Active:           "Active",
			ActiveInfo:       "Inactive plan groups are hidden from new plan assignments.",
		},
		Bulk: BulkLabels{
			DeleteTitle:       "Delete Plan Groups",
			DeleteMessage:     "Permanently delete the selected plan groups? This cannot be undone.",
			ActivateTitle:     "Activate Plan Groups",
			ActivateMessage:   "Activate the selected plan groups?",
			DeactivateTitle:   "Deactivate Plan Groups",
			DeactivateMessage: "Deactivate the selected plan groups?",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:       "Delete Plan Group",
			DeleteMessage:     "Permanently delete this plan group? This cannot be undone.",
			ActivateTitle:     "Activate Plan Group",
			ActivateMessage:   "Activate {{name}}?",
			DeactivateTitle:   "Deactivate Plan Group",
			DeactivateMessage: "Deactivate {{name}}?",
		},
		Tabs: TabLabels{
			Info: "Info",
		},
		Detail: DetailLabels{
			Title:        "Plan Group",
			DateCreated:  "Date Created",
			DateModified: "Date Modified",
			NoParent:     "No parent group",
			NoCode:       "—",
			NoSubtitle:   "No description provided",
		},
		Errors: ErrorLabels{
			NotFound:     "Plan group not found",
			LoadFailed:   "Failed to load plan group",
			Unauthorized: "You are not authorized to perform this action",
			CreateFailed: "Failed to create plan group",
			UpdateFailed: "Failed to update plan group",
			DeleteFailed: "Failed to delete plan group",
			InUse:        "This plan group is in use and cannot be deleted.",
		},
	}
}
