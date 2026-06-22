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
	SubscriptionGroupId string `json:"subscriptionGroupId"`
	SubscriptionId      string `json:"subscriptionId"`
	ClientId            string `json:"clientId"`
	Status              string `json:"status"`
	DateCreated         string `json:"dateCreated"`
	Actions             string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FormLabels holds the drawer-form field labels.
type FormLabels struct {
	SectionIdentity       string `json:"sectionIdentity"`
	SubscriptionGroup     string `json:"subscriptionGroup"`
	SubscriptionGroupPH   string `json:"subscriptionGroupPlaceholder"`
	SubscriptionGroupInfo string `json:"subscriptionGroupInfo"`
	Subscription          string `json:"subscription"`
	SubscriptionPH        string `json:"subscriptionPlaceholder"`
	SubscriptionInfo      string `json:"subscriptionInfo"`
	Client                string `json:"client"`
	ClientPH              string `json:"clientPlaceholder"`
	ClientInfo            string `json:"clientInfo"`
	Active                string `json:"active"`
	ActiveInfo            string `json:"activeInfo"`
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
	NoSub        string `json:"noSub"`
	NoClient     string `json:"noClient"`
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
