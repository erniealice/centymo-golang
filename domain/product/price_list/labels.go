package price_list

// ---------------------------------------------------------------------------
// Price List labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the price list module.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Buttons ButtonLabels  `json:"buttons"`
	Columns ColumnLabels  `json:"columns"`
	Empty   EmptyLabels   `json:"empty"`
	Form    FormLabels    `json:"form"`
	Actions ActionLabels  `json:"actions"`
	Bulk    BulkLabels    `json:"bulk_actions"`
	Detail  DetailLabels  `json:"detail"`
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
	AddPriceList string `json:"add_price_list"`
}

type ColumnLabels struct {
	Name      string `json:"name"`
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
	Status    string `json:"status"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"active_title"`
	ActiveMessage   string `json:"active_message"`
	InactiveTitle   string `json:"inactive_title"`
	InactiveMessage string `json:"inactive_message"`
}

type FormLabels struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DescPlaceholder string `json:"description_placeholder"`
	DateStart       string `json:"date_start"`
	DateEnd         string `json:"date_end"`
	Active          string `json:"active"`
	Product         string `json:"product"`
	SelectProduct   string `json:"select_product"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"name_info"`
	DescriptionInfo string `json:"description_info"`
	DateStartInfo   string `json:"date_start_info"`
	DateEndInfo     string `json:"date_end_info"`
	ActiveInfo      string `json:"active_info"`
	// Price-product sub-drawer info fields.
	AmountInfo   string `json:"amount_info"`
	CurrencyInfo string `json:"currency_info"`
}

type ActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

type DetailLabels struct {
	PageTitle          string `json:"page_title"`
	BasicInfo          string `json:"basic_info"`
	Prices             string `json:"prices"`
	TabAttachments     string `json:"tab_attachments"`
	TabAuditHistory    string `json:"tab_audit_history"`
	ProductName        string `json:"product_name"`
	Amount             string `json:"amount"`
	Currency           string `json:"currency"`
	AddPrice           string `json:"add_price"`
	RemoveLabel        string `json:"remove_label"`
	EmptyTitle         string `json:"empty_title"`
	EmptyMessage       string `json:"empty_message"`
	ActiveBadge        string `json:"active_badge"`
	InactiveBadge      string `json:"inactive_badge"`
	NoPricesConfigured string `json:"no_prices_configured"`
	NoPricesDesc       string `json:"no_prices_desc"`
}

type ConfirmLabels struct {
	Activate          string `json:"activate"`
	ActivateMessage   string `json:"activate_message"`
	Deactivate        string `json:"deactivate"`
	DeactivateMessage string `json:"deactivate_message"`
	Delete            string `json:"delete"`
	DeleteMessage     string `json:"delete_message"`
	BulkDelete        string `json:"bulk_delete"`
	BulkDeleteMessage string `json:"bulk_delete_message"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	NoIDsProvided    string `json:"no_ids_provided"`
	CannotDelete     string `json:"cannot_delete"`
	ProductRequired  string `json:"product_required"`
	AmountRequired   string `json:"amount_required"`
}

// DefaultLabels returns the zero-value label set. Every rendered string for
// this entity must come from the lyngua cascade (general -> business-type
// tier); there are no Go-side default strings to fall back on.
func DefaultLabels() Labels { return Labels{} }
