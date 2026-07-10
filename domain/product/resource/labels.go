package resource

// ---------------------------------------------------------------------------
// Resource labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the resource module.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Buttons ButtonLabels  `json:"buttons"`
	Columns ColumnLabels  `json:"columns"`
	Empty   EmptyLabels   `json:"empty"`
	Form    FormLabels    `json:"form"`
	Actions ActionLabels  `json:"actions"`
	Bulk    BulkLabels    `json:"bulk_actions"`
	Status  StatusLabels  `json:"status"`
	Confirm ConfirmLabels `json:"confirm"`
	Errors  ErrorLabels   `json:"errors"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"heading_active"`
	HeadingInactive string `json:"heading_inactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"caption_active"`
	CaptionInactive string `json:"caption_inactive"`
}

type ButtonLabels struct {
	Add string `json:"add"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Product     string `json:"product"`
	Status      string `json:"status"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type FormLabels struct {
	Name            string `json:"name"`
	NamePlaceholder string `json:"name_placeholder"`
	Description     string `json:"description"`
	DescPlaceholder string `json:"description_placeholder"`
	ProductId       string `json:"product_id"`
	UserId          string `json:"user_id"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"name_info"`
	DescriptionInfo string `json:"description_info"`
	ProductIdInfo   string `json:"product_id_info"`
	UserIdInfo      string `json:"user_id_info"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

type StatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ConfirmLabels struct {
	Delete              string `json:"delete"`
	DeleteMessage       string `json:"delete_message"`
	Activate            string `json:"activate"`
	ActivateMessage     string `json:"activate_message"`
	Deactivate          string `json:"deactivate"`
	DeactivateMessage   string `json:"deactivate_message"`
	BulkDelete          string `json:"bulk_delete"`
	BulkDeleteMessage   string `json:"bulk_delete_message"`
	BulkActivate        string `json:"bulk_activate"`
	BulkActivateMessage string `json:"bulk_activate_message"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	NoPermission     string `json:"no_permission"`
	CannotDelete     string `json:"cannot_delete"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:         "Resources",
			HeadingActive:   "Active Resources",
			HeadingInactive: "Inactive Resources",
			Caption:         "Manage resources linked to products.",
			CaptionActive:   "Showing active resources.",
			CaptionInactive: "Showing inactive resources.",
		},
		Buttons: ButtonLabels{
			Add: "Add Resource",
		},
		Columns: ColumnLabels{
			Name:        "Name",
			Description: "Description",
			Product:     "Product",
			Status:      "Status",
		},
		Empty: EmptyLabels{
			Title:   "No resources found",
			Message: "Add a resource to get started.",
		},
		Form: FormLabels{
			Name:            "Name",
			NamePlaceholder: "Enter resource name",
			Description:     "Description",
			DescPlaceholder: "Enter description (optional)",
			ProductId:       "Product ID",
			UserId:          "User ID",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			NameInfo:        "Display name for this resource.",
			DescriptionInfo: "Optional notes about this resource.",
			ProductIdInfo:   "The product this resource is linked to (used for activity billing).",
			UserIdInfo:      "Optional — restrict this resource to a specific user.",
		},
		Actions: ActionLabels{
			View:       "View",
			Edit:       "Edit",
			Delete:     "Delete",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Bulk: BulkLabels{
			Delete: "Delete Selected",
		},
		Status: StatusLabels{
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Confirm: ConfirmLabels{
			Delete:              "Delete Resource",
			DeleteMessage:       "Are you sure you want to delete this resource?",
			Activate:            "Activate Resource",
			ActivateMessage:     "Activate resource \"%s\"?",
			Deactivate:          "Deactivate Resource",
			DeactivateMessage:   "Deactivate resource \"%s\"?",
			BulkDelete:          "Delete Selected",
			BulkDeleteMessage:   "Are you sure you want to delete the selected resources?",
			BulkActivate:        "Activate Selected",
			BulkActivateMessage: "Activate the selected resources?",
		},
		Errors: ErrorLabels{
			PermissionDenied: "You do not have permission to perform this action",
			InvalidFormData:  "Invalid form data. Please check your inputs and try again.",
			NotFound:         "Resource not found",
			IDRequired:       "Resource ID is required",
			NoPermission:     "No permission",
			CannotDelete:     "This resource cannot be deleted because it is in use",
		},
	}
}
