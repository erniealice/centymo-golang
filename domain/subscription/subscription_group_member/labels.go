package subscription_group_member

// ---------------------------------------------------------------------------
// Subscription Group Member labels
// ---------------------------------------------------------------------------

// Labels holds all labels for the subscription_group_member module.
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
	SubscriptionGroupId string `json:"subscription_group_id"`
	SubscriptionId      string `json:"subscription_id"`
	ClientId            string `json:"client_id"`
	Status              string `json:"status"`
	DateCreated         string `json:"date_created"`
	Actions             string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels.
type FormLabels struct {
	SectionIdentity       string `json:"section_identity"`
	SubscriptionGroup     string `json:"subscription_group"`
	SubscriptionGroupPH   string `json:"subscription_group_placeholder"`
	SubscriptionGroupInfo string `json:"subscription_group_info"`
	Subscription          string `json:"subscription"`
	SubscriptionPH        string `json:"subscription_placeholder"`
	SubscriptionInfo      string `json:"subscription_info"`
	Client                string `json:"client"`
	ClientPH              string `json:"client_placeholder"`
	ClientInfo            string `json:"client_info"`
	Active                string `json:"active"`
	ActiveInfo            string `json:"active_info"`
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
	NoSub        string `json:"no_sub"`
	NoClient     string `json:"no_client"`
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
			Title:         "Members",
			Subtitle:      "Manage subscription group members",
			ActiveTitle:   "Active Members",
			InactiveTitle: "Inactive Members",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Member",
			Edit:       "Edit Member",
			Delete:     "Delete Member",
			BulkDelete: "Delete Members",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: ColumnLabels{
			SubscriptionGroupId: "Group",
			SubscriptionId:      "Subscription",
			ClientId:            "Client",
			Status:              "Status",
			DateCreated:         "Date Created",
			Actions:             "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Members",
			Message: "No members to display.",
		},
		Form: FormLabels{
			SectionIdentity:       "Member details",
			SubscriptionGroup:     "Subscription group",
			SubscriptionGroupPH:   "Enter group ID",
			SubscriptionGroupInfo: "The subscription group this member belongs to.",
			Subscription:          "Subscription",
			SubscriptionPH:        "Enter subscription ID",
			SubscriptionInfo:      "The subscription linked to this member.",
			Client:                "Client",
			ClientPH:              "Enter client ID",
			ClientInfo:            "The client (account) associated with this member.",
			Active:                "Active",
			ActiveInfo:            "Inactive members are hidden from active group rosters.",
		},
		Bulk: BulkLabels{
			DeleteTitle:       "Delete Members",
			DeleteMessage:     "Permanently delete the selected members? This cannot be undone.",
			ActivateTitle:     "Activate Members",
			ActivateMessage:   "Activate the selected members?",
			DeactivateTitle:   "Deactivate Members",
			DeactivateMessage: "Deactivate the selected members?",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:       "Delete Member",
			DeleteMessage:     "Permanently delete this member? This cannot be undone.",
			ActivateTitle:     "Activate Member",
			ActivateMessage:   "Activate this member?",
			DeactivateTitle:   "Deactivate Member",
			DeactivateMessage: "Deactivate this member?",
		},
		Tabs: TabLabels{
			Info: "Info",
		},
		Detail: DetailLabels{
			Title:        "Member",
			DateCreated:  "Date Created",
			DateModified: "Date Modified",
			NoGroup:      "No group",
			NoSub:        "No subscription",
			NoClient:     "No client",
		},
		Errors: ErrorLabels{
			NotFound:     "Member not found",
			LoadFailed:   "Failed to load member",
			Unauthorized: "You are not authorized to perform this action",
			CreateFailed: "Failed to create member",
			UpdateFailed: "Failed to update member",
			DeleteFailed: "Failed to delete member",
			InUse:        "This member is in use and cannot be deleted.",
		},
	}
}
