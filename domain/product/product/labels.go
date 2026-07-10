package product

// OptionValueSeparator is the canonical separator between concatenated
// product_option_value labels. Used by the variants table on the product
// detail page and by every drawer picker that surfaces a variant's
// option-value tuple inline (e.g., "Red / Large / Cotton").
const OptionValueSeparator = " / "

type TrackingModeLabels struct {
	None       string `json:"none"`
	Bulk       string `json:"bulk"`
	Serialized string `json:"serialized"`
}

// KindLabels holds the translated labels for each product_kind enum
// value. Sourced from lyngua product.json "productKind" block. Wired onto
// Labels so the drawer-form select can render the per-value labels
// using the exact tier-cascaded strings that appear elsewhere in the UI.
type KindLabels struct {
	Service        string `json:"service"`
	StockedGood    string `json:"stocked_good"`
	NonStockedGood string `json:"non_stocked_good"`
	Consumable     string `json:"consumable"`
}

// DeliveryModeLabels mirrors KindLabels for the delivery_mode axis.
type DeliveryModeLabels struct {
	Instant      string `json:"instant"`
	Scheduled    string `json:"scheduled"`
	Shipped      string `json:"shipped"`
	Digital      string `json:"digital"`
	Project      string `json:"project"`
	Subscription string `json:"subscription"`
}

// ---------------------------------------------------------------------------
// Product labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the product module.
type Labels struct {
	Page       PageLabels       `json:"page"`
	Buttons    ButtonLabels     `json:"buttons"`
	Columns    ColumnLabels     `json:"columns"`
	Empty      EmptyLabels      `json:"empty"`
	Form       FormLabels       `json:"form"`
	Actions    ActionLabels     `json:"actions"`
	Bulk       BulkLabels       `json:"bulk_actions"`
	Tabs       TabLabels        `json:"tabs"`
	Detail     DetailLabels     `json:"detail"`
	Status     StatusLabels     `json:"status"`
	Variant    VariantLabels    `json:"variant"`
	Attribute  AttributeLabels  `json:"attribute"`
	Options    OptionLabels     `json:"options"`
	Confirm    ConfirmLabels    `json:"confirm"`
	Errors     ErrorLabels      `json:"errors"`
	Breadcrumb BreadcrumbLabels `json:"breadcrumb"`
	// Four-axis product taxonomy enum labels — loaded from lyngua
	// product.json "productKind"/"deliveryMode"/"trackingMode" blocks.
	// Wired here so the drawer-form select uses the exact tier-cascaded
	// display string for each enum value without hardcoding in Go.
	ProductKind  KindLabels         `json:"product_kind"`
	DeliveryMode DeliveryModeLabels `json:"delivery_mode"`
	TrackingMode TrackingModeLabels `json:"tracking_mode"`

	// Phase 5 — service dashboard (product_kind=service surface).
	ServiceDashboard ServiceDashboardLabels `json:"service_dashboard"`
}

// ServiceDashboardLabels holds translatable strings for the service
// dashboard. The "Service" wording is preferred at the dashboard surface
// because the sidebar key is "service"; the underlying entity is still
// Product filtered to product_kind="service".
type ServiceDashboardLabels struct {
	Title              string `json:"title"`
	Subtitle           string `json:"subtitle"`
	StatTotalActive    string `json:"stat_total_active"`
	StatTopRevenue     string `json:"stat_top_revenue"`
	StatByLineCount    string `json:"stat_by_line_count"`
	StatRecentlyAdded  string `json:"stat_recently_added"`
	WidgetByLine       string `json:"widget_by_line"`
	WidgetTopRevenue   string `json:"widget_top_revenue"`
	WidgetRecent       string `json:"widget_recent"`
	QuickNew           string `json:"quick_new"`
	QuickBundleBuilder string `json:"quick_bundle_builder"`
	QuickTagService    string `json:"quick_tag_service"`
	QuickPriceSchedule string `json:"quick_price_schedule"`
	ViewAll            string `json:"view_all"`
	EmptyRecentTitle   string `json:"empty_recent_title"`
	EmptyRecentDesc    string `json:"empty_recent_desc"`
	EmptyTopRevenue    string `json:"empty_top_revenue"`
	NewService         string `json:"new_service"`
	ColLine            string `json:"col_line"`
	ColRank            string `json:"col_rank"`
	ColService         string `json:"col_service"`
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
	AddProduct string `json:"add_product"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Line        string `json:"line"`
	Price       string `json:"price"`
	Status      string `json:"status"`
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
	Price           string `json:"price"`
	Currency        string `json:"currency"`
	Active          string `json:"active"`
	Line            string `json:"line"`
	LinePlaceholder string `json:"line_placeholder"`

	// Variant / option / attribute form labels
	PricePlaceholder       string `json:"price_placeholder"`
	SelectOption           string `json:"select_option"`
	Required               string `json:"required"`
	Option                 string `json:"option"`
	SelectAttribute        string `json:"select_attribute"`
	AllAttributesAssigned  string `json:"all_attributes_assigned"`
	OptionNeedsValuesAlert string `json:"option_needs_values_alert"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"name_info"`
	DescriptionInfo string `json:"description_info"`
	LineInfo        string `json:"line_info"`
	PriceInfo       string `json:"price_info"`
	CurrencyInfo    string `json:"currency_info"`
	ActiveInfo      string `json:"active_info"`

	// Model D — variant_mode toggle + unit field
	VariantModeLabel        string `json:"variant_mode_label"`
	VariantModeInfo         string `json:"variant_mode_info"`
	VariantModeNone         string `json:"variant_mode_none"`
	VariantModeConfigurable string `json:"variant_mode_configurable"`
	UnitLabel               string `json:"unit_label"`
	UnitInfo                string `json:"unit_info"`
	UnitPlaceholder         string `json:"unit_placeholder"`
	VariantPriceVaries      string `json:"variant_price_varies"`
	// Shown as help text beneath the variant toggle when the product already
	// has option or variant rows, to explain why the toggle is disabled.
	VariantModeLockedHelp string `json:"variant_mode_locked_help"`
	// Error surfaced by the Create/Update handlers when a caller tries to
	// flip variant_mode on a product that still has options/variants.
	VariantModeLockedError string `json:"variant_mode_locked_error"`

	// Four-axis product taxonomy — rendered as selects on the drawer form.
	// Each axis carries its own Label + Info popover text plus per-enum-value
	// Info (XxxValueInfo map) keyed by enum string. When the mount restricts
	// the axis to one allowed value the select is rendered disabled so the
	// user still sees the classification without being able to change it.
	ProductKindLabel      string            `json:"product_kind_label"`
	ProductKindInfo       string            `json:"product_kind_info"`
	ProductKindValueInfo  map[string]string `json:"product_kind_value_info,omitempty"`
	DeliveryModeLabel     string            `json:"delivery_mode_label"`
	DeliveryModeInfo      string            `json:"delivery_mode_info"`
	DeliveryModeValueInfo map[string]string `json:"delivery_mode_value_info,omitempty"`
	TrackingModeLabel     string            `json:"tracking_mode_label"`
	TrackingModeInfo      string            `json:"tracking_mode_info"`
	TrackingModeValueInfo map[string]string `json:"tracking_mode_value_info,omitempty"`

	// Tax section labels (Phase 5)
	SectionTax                  string `json:"section_tax"`
	TaxTreatmentLabel           string `json:"tax_treatment_label"`
	TaxTreatmentPlaceholder     string `json:"tax_treatment_placeholder"`
	TaxTreatmentInfo            string `json:"tax_treatment_info"`
	WithholdingClassLabel       string `json:"withholding_class_label"`
	WithholdingClassPlaceholder string `json:"withholding_class_placeholder"`
	WithholdingClassInfo        string `json:"withholding_class_info"`
}

type ActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

type TabLabels struct {
	Info         string `json:"info"`
	Variants     string `json:"variants"`
	Attributes   string `json:"attributes"`
	Pricing      string `json:"pricing"`
	Options      string `json:"options"`
	Images       string `json:"images"`
	Stock        string `json:"stock"`
	Lines        string `json:"lines"`
	Attachments  string `json:"attachments"`
	AuditTrail   string `json:"audit_trail"`
	AuditHistory string `json:"audit_history"`
	// Inventory item sub-tabs
	Serials        string `json:"serials"`
	PricingHistory string `json:"pricing_history"`
}

type DetailLabels struct {
	Price                string `json:"price"`
	Currency             string `json:"currency"`
	Collections          string `json:"collections"`
	VariantCount         string `json:"variant_count"`
	Status               string `json:"status"`
	OptionsLabel         string `json:"options_label"`
	EmptyVariantsMessage string `json:"empty_variants_message"`
	// Header subtitle fallback when the product has no description.
	// Consumed by buildPageData to override the generic "Welcome back"
	// CommonLabels default on the product detail page header.
	NoDescriptionSubtitle string `json:"no_description_subtitle"`
	// Model D — detail-page rows for unit of measure + variant mode.
	// Falls back to English defaults when lyngua doesn't overlay the key.
	Unit        string `json:"unit"`
	VariantMode string `json:"variant_mode"`
	// Serial table columns
	SerialNumber       string `json:"serial_number"`
	IMEI               string `json:"imei"`
	WarrantyEnd        string `json:"warranty_end"`
	PurchaseOrder      string `json:"purchase_order"`
	NoSerialNumbers    string `json:"no_serial_numbers"`
	NoSerialNumbersMsg string `json:"no_serial_numbers_msg"`

	// Variant detail labels
	VariantInformation  string `json:"variant_information"`
	Options             string `json:"options"`
	VariantPricing      string `json:"variant_pricing"`
	VariantPricingDesc  string `json:"variant_pricing_desc"`
	InventoryStock      string `json:"inventory_stock"`
	InventoryStockDesc  string `json:"inventory_stock_desc"`
	DropImagesHere      string `json:"drop_images_here"`
	ImageFileHint       string `json:"image_file_hint"`
	DeleteSelected      string `json:"delete_selected"`
	PrimaryBadge        string `json:"primary_badge"`
	NoImages            string `json:"no_images"`
	NoImagesDesc        string `json:"no_images_desc"`
	AuditTrail          string `json:"audit_trail"`
	AuditTrailDesc      string `json:"audit_trail_desc"`
	NoSerialNumbersDesc string `json:"no_serial_numbers_desc"`

	// Stock detail labels
	InventoryItem      string `json:"inventory_item"`
	Name               string `json:"name"`
	SKU                string `json:"sku"`
	Type               string `json:"type"`
	Location           string `json:"location"`
	QtyOnHand          string `json:"qty_on_hand"`
	Reserved           string `json:"reserved"`
	Available          string `json:"available"`
	StatTotal          string `json:"stat_total"`
	StatAvailable      string `json:"stat_available"`
	StatSold           string `json:"stat_sold"`
	StatReserved       string `json:"stat_reserved"`
	PricingHistory     string `json:"pricing_history"`
	PricingHistoryDesc string `json:"pricing_history_desc"`

	// Serial detail labels
	SerialInformation string `json:"serial_information"`
}

type StatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type VariantLabels struct {
	Title         string `json:"title"`
	SKU           string `json:"sku"`
	PriceOverride string `json:"price_override"`
	Attributes    string `json:"attributes"`
	Assign        string `json:"assign"`
	Edit          string `json:"edit"`
	Remove        string `json:"remove"`
	Empty         string `json:"empty"`
	// Stock table columns
	Location    string `json:"location"`
	QtyOnHand   string `json:"qty_on_hand"`
	SerialCount string `json:"serial_count"`
	NoStock     string `json:"no_stock"`
	NoStockMsg  string `json:"no_stock_msg"`
	// Pricing tab column headers
	Pricing VariantPricingLabels `json:"pricing"`
}

// VariantPricingLabels holds column header labels for the variant pricing tab table.
type VariantPricingLabels struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Package  string `json:"package"`
	RateCard string `json:"rate_card"`
	Amount   string `json:"amount"`
}

type AttributeLabels struct {
	Title        string `json:"title"`
	DefaultValue string `json:"default_value"`
	Assign       string `json:"assign"`
	Remove       string `json:"remove"`
	Empty        string `json:"empty"`
}

type ConfirmLabels struct {
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activate_message"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivate_message"`
	BulkActivate          string `json:"bulk_activate"`
	BulkActivateMessage   string `json:"bulk_activate_message"`
	BulkDeactivate        string `json:"bulk_deactivate"`
	BulkDeactivateMessage string `json:"bulk_deactivate_message"`
	BulkDelete            string `json:"bulk_delete"`
	BulkDeleteMessage     string `json:"bulk_delete_message"`
	RemoveVariant         string `json:"remove_variant"`
	RemoveVariantMessage  string `json:"remove_variant_message"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	NoIDsProvided    string `json:"no_ids_provided"`
	InvalidStatus    string `json:"invalid_status"`
	CannotDelete     string `json:"cannot_delete"`
	NameRequired     string `json:"name_required"`
	FieldRequired    string `json:"field_required"`
}

type BreadcrumbLabels struct {
	Products string `json:"products"`
	Product  string `json:"product"`
	Option   string `json:"option"`
}

// ---------------------------------------------------------------------------
// Product Option labels
// ---------------------------------------------------------------------------

type OptionLabels struct {
	Tab       OptionTabLabels      `json:"tab"`
	Tabs      OptionTabsLabels     `json:"tabs"`
	Columns   OptionColumnLabels   `json:"columns"`
	Form      OptionFormLabels     `json:"form"`
	DataTypes OptionDataTypeLabels `json:"data_types"`
	Value     OptionValueLabels    `json:"value"`
	Actions   OptionActionLabels   `json:"actions"`
	Empty     OptionEmptyLabels    `json:"empty"`
	Confirm   OptionConfirmLabels  `json:"confirm"`
}

type OptionTabLabels struct {
	Title string `json:"title"`
}

type OptionTabsLabels struct {
	Info   string `json:"info"`
	Values string `json:"values"`
}

type OptionColumnLabels struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	DataType    string `json:"data_type"`
	ValuesCount string `json:"values_count"`
	SortOrder   string `json:"sort_order"`
	Required    string `json:"required"`
	Status      string `json:"status"`
}

type OptionFormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"name_placeholder"`
	Code                   string `json:"code"`
	CodePlaceholder        string `json:"code_placeholder"`
	DataType               string `json:"data_type"`
	SortOrder              string `json:"sort_order"`
	MinValue               string `json:"min_value"`
	MaxValue               string `json:"max_value"`
	Active                 string `json:"active"`
	Required               string `json:"required"`
	RequiredCaution        string `json:"required_caution"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	DescriptionEmpty       string `json:"description_empty"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"name_info"`
	CodeInfo        string `json:"code_info"`
	DataTypeInfo    string `json:"data_type_info"`
	MinValueInfo    string `json:"min_value_info"`
	MaxValueInfo    string `json:"max_value_info"`
	SortOrderInfo   string `json:"sort_order_info"`
	ActiveInfo      string `json:"active_info"`
	DescriptionInfo string `json:"description_info"`
}

type OptionDataTypeLabels struct {
	TextList    string `json:"text_list"`
	NumberRange string `json:"number_range"`
	ColorList   string `json:"color_list"`
	FreeText    string `json:"free_text"`
	FreeNumber  string `json:"free_number"`
}

type OptionValueLabels struct {
	Columns OptionValueColumnLabels `json:"columns"`
	Form    OptionValueFormLabels   `json:"form"`
}

type OptionValueColumnLabels struct {
	Label        string `json:"label"`
	Value        string `json:"value"`
	SortOrder    string `json:"sort_order"`
	ColorPreview string `json:"color_preview"`
	Status       string `json:"status"`
}

type OptionValueFormLabels struct {
	Label               string `json:"label"`
	LabelPlaceholder    string `json:"label_placeholder"`
	Value               string `json:"value"`
	ValuePlaceholder    string `json:"value_placeholder"`
	SortOrder           string `json:"sort_order"`
	ColorHex            string `json:"color_hex"`
	ColorHexPlaceholder string `json:"color_hex_placeholder"`
	Active              string `json:"active"`
	// Context labels surfaced on the value drawer to remind the user
	// which option this value belongs to.
	Option   string `json:"option"`
	Required string `json:"required"`

	// Field-level info text surfaced via an info button beside each label.
	LabelInfo     string `json:"label_info"`
	ValueInfo     string `json:"value_info"`
	SortOrderInfo string `json:"sort_order_info"`
	ColorHexInfo  string `json:"color_hex_info"`
	ActiveInfo    string `json:"active_info"`
}

type OptionActionLabels struct {
	AddOption         string `json:"add_option"`
	EditOption        string `json:"edit_option"`
	EditProductOption string `json:"edit_product_option"`
	DeleteOption      string `json:"delete_option"`
	ViewValues        string `json:"view_values"`
	AddValue          string `json:"add_value"`
	EditValue         string `json:"edit_value"`
	DeleteValue       string `json:"delete_value"`
}

type OptionEmptyLabels struct {
	Title        string `json:"title"`
	Message      string `json:"message"`
	ValueTitle   string `json:"value_title"`
	ValueMessage string `json:"value_message"`
}

type OptionConfirmLabels struct {
	DeleteOption string `json:"delete_option"`
	DeleteValue  string `json:"delete_value"`
}

// ---------------------------------------------------------------------------
// Product Line labels
// ---------------------------------------------------------------------------

// LineLabels holds all translatable strings for the product line module.
type LineLabels struct {
	Page    LinePageLabels    `json:"page"`
	Buttons LineButtonLabels  `json:"buttons"`
	Columns LineColumnLabels  `json:"columns"`
	Empty   LineEmptyLabels   `json:"empty"`
	Form    LineFormLabels    `json:"form"`
	Actions LineActionLabels  `json:"actions"`
	Bulk    LineBulkLabels    `json:"bulk_actions"`
	Tabs    LineTabLabels     `json:"tabs"`
	Detail  LineDetailLabels  `json:"detail"`
	Status  LineStatusLabels  `json:"status"`
	Confirm LineConfirmLabels `json:"confirm"`
	Errors  LineErrorLabels   `json:"errors"`
}

type LinePageLabels struct {
	Heading          string `json:"heading"`
	HeadingActive    string `json:"heading_active"`
	HeadingInactive  string `json:"heading_inactive"`
	HeadingPending   string `json:"heading_pending"`
	HeadingCompleted string `json:"heading_completed"`
	HeadingFailed    string `json:"heading_failed"`
	Caption          string `json:"caption"`
	CaptionActive    string `json:"caption_active"`
	CaptionInactive  string `json:"caption_inactive"`
	CaptionPending   string `json:"caption_pending"`
	CaptionCompleted string `json:"caption_completed"`
	CaptionFailed    string `json:"caption_failed"`
}

type LineButtonLabels struct {
	AddProductLine    string `json:"add_product_line"`
	EditProductLine   string `json:"edit_product_line"`
	DeleteProductLine string `json:"delete_product_line"`
}

type LineColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DateCreated string `json:"date_created"`
	Reference   string `json:"reference"`
	Customer    string `json:"customer"`
	Amount      string `json:"amount"`
	Method      string `json:"method"`
	Date        string `json:"date"`
	Status      string `json:"status"`
}

type LineEmptyLabels struct {
	ActiveTitle      string `json:"active_title"`
	ActiveMessage    string `json:"active_message"`
	InactiveTitle    string `json:"inactive_title"`
	InactiveMessage  string `json:"inactive_message"`
	PendingTitle     string `json:"pending_title"`
	PendingMessage   string `json:"pending_message"`
	CompletedTitle   string `json:"completed_title"`
	CompletedMessage string `json:"completed_message"`
	FailedTitle      string `json:"failed_title"`
	FailedMessage    string `json:"failed_message"`
}

type LineFormLabels struct {
	Name                    string `json:"name"`
	NamePlaceholder         string `json:"name_placeholder"`
	Description             string `json:"description"`
	DescPlaceholder         string `json:"description_placeholder"`
	Active                  string `json:"active"`
	Customer                string `json:"customer"`
	Date                    string `json:"date"`
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	Reference               string `json:"reference"`
	ReferencePlaceholder    string `json:"reference_placeholder"`
	PaymentMethod           string `json:"payment_method"`
	Status                  string `json:"status"`
	Notes                   string `json:"notes"`
	NotesPlaceholder        string `json:"notes_placeholder"`
	CustomerNamePlaceholder string `json:"customer_name_placeholder"`
	AmountPlaceholder       string `json:"amount_placeholder"`
	CurrencyPlaceholder     string `json:"currency_placeholder"`
	MethodCash              string `json:"method_cash"`
	MethodBankTransfer      string `json:"method_bank_transfer"`
	MethodCheck             string `json:"method_check"`
	MethodGCash             string `json:"method_gcash"`
	MethodMaya              string `json:"method_maya"`
	MethodCard              string `json:"method_card"`
	MethodOther             string `json:"method_other"`
	StatusPending           string `json:"status_pending"`
	StatusCompleted         string `json:"status_completed"`
	StatusFailed            string `json:"status_failed"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"name_info"`
	DescriptionInfo string `json:"description_info"`
	ActiveInfo      string `json:"active_info"`
}

type LineActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	MarkComplete string `json:"mark_complete"`
	Reactivate   string `json:"reactivate"`
}

type LineBulkLabels struct {
	Delete string `json:"delete"`
}

type LineStatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type LineTabLabels struct {
	Info string `json:"info"`
}

type LineDetailLabels struct {
	TitlePrefix          string `json:"title_prefix"`
	PageTitle            string `json:"page_title"`
	BasicInfo            string `json:"basic_info"`
	PaymentInfo          string `json:"payment_info"`
	Reference            string `json:"reference"`
	Customer             string `json:"customer"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Method               string `json:"method"`
	Date                 string `json:"date"`
	Status               string `json:"status"`
	Notes                string `json:"notes"`
	CreatedDate          string `json:"created_date"`
	ModifiedDate         string `json:"modified_date"`
	ActiveBadge          string `json:"active_badge"`
	InactiveBadge        string `json:"inactive_badge"`
	TabBasicInfo         string `json:"tab_basic_info"`
	TabAttachments       string `json:"tab_attachments"`
	TabAuditTrail        string `json:"tab_audit_trail"`
	TabAuditHistory      string `json:"tab_audit_history"`
	AuditAction          string `json:"audit_action"`
	AuditUser            string `json:"audit_user"`
	AuditEmptyTitle      string `json:"audit_empty_title"`
	AuditEmptyMessage    string `json:"audit_empty_message"`
	AuditTrailComingSoon string `json:"audit_trail_coming_soon"`
	AuditTrailDesc       string `json:"audit_trail_desc"`
}

type LineConfirmLabels struct {
	MarkComplete          string `json:"mark_complete"`
	MarkCompleteMessage   string `json:"mark_complete_message"`
	Reactivate            string `json:"reactivate"`
	ReactivateMessage     string `json:"reactivate_message"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"delete_message"`
	BulkActivate          string `json:"bulk_activate"`
	BulkActivateMessage   string `json:"bulk_activate_message"`
	BulkDeactivate        string `json:"bulk_deactivate"`
	BulkDeactivateMessage string `json:"bulk_deactivate_message"`
	BulkComplete          string `json:"bulk_complete"`
	BulkCompleteMessage   string `json:"bulk_complete_message"`
	BulkReactivate        string `json:"bulk_reactivate"`
	BulkReactivateMessage string `json:"bulk_reactivate_message"`
	BulkDelete            string `json:"bulk_delete"`
	BulkDeleteMessage     string `json:"bulk_delete_message"`
}

type LineErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	NoIDsProvided    string `json:"no_ids_provided"`
	InvalidStatus    string `json:"invalid_status"`
	CannotDelete     string `json:"cannot_delete"`
}

// DefaultLineLabels returns LineLabels with sensible English defaults.
func DefaultLineLabels() LineLabels {
	return LineLabels{
		Page: LinePageLabels{
			Heading:          "Product Lines",
			HeadingActive:    "Active Product Lines",
			HeadingInactive:  "Inactive Product Lines",
			HeadingPending:   "Pending Product Lines",
			HeadingCompleted: "Completed Product Lines",
			HeadingFailed:    "Failed Product Lines",
			Caption:          "Manage product lines",
			CaptionActive:    "Active product lines",
			CaptionInactive:  "Inactive product lines",
			CaptionPending:   "Product lines awaiting completion",
			CaptionCompleted: "Completed product lines",
			CaptionFailed:    "Failed product lines",
		},
		Buttons: LineButtonLabels{
			AddProductLine:    "Add Product Line",
			EditProductLine:   "Edit Product Line",
			DeleteProductLine: "Delete Product Line",
		},
		Columns: LineColumnLabels{
			Name:        "Name",
			Description: "Description",
			DateCreated: "Date Created",
			Reference:   "Reference",
			Customer:    "Customer",
			Amount:      "Amount",
			Method:      "Method",
			Date:        "Date",
			Status:      "Status",
		},
		Empty: LineEmptyLabels{
			ActiveTitle:      "No active product lines",
			ActiveMessage:    "No active product lines to display.",
			InactiveTitle:    "No inactive product lines",
			InactiveMessage:  "No inactive product lines to display.",
			PendingTitle:     "No pending product lines",
			PendingMessage:   "No pending product lines to display.",
			CompletedTitle:   "No completed product lines",
			CompletedMessage: "No completed product lines to display.",
			FailedTitle:      "No failed product lines",
			FailedMessage:    "No failed product lines to display.",
		},
		Form: LineFormLabels{
			Name:                    "Name",
			NamePlaceholder:         "Product line name",
			Description:             "Description",
			DescPlaceholder:         "Optional description",
			Active:                  "Active",
			Customer:                "Customer",
			Date:                    "Date",
			Amount:                  "Amount",
			Currency:                "Currency",
			Reference:               "Reference",
			ReferencePlaceholder:    "e.g. PL-001",
			PaymentMethod:           "Payment Method",
			Status:                  "Status",
			Notes:                   "Notes",
			NotesPlaceholder:        "Additional notes...",
			CustomerNamePlaceholder: "Customer name",
			AmountPlaceholder:       "0.00",
			CurrencyPlaceholder:     "PHP",
			MethodCash:              "Cash",
			MethodBankTransfer:      "Bank Transfer",
			MethodCheck:             "Check",
			MethodGCash:             "GCash",
			MethodMaya:              "Maya",
			MethodCard:              "Card",
			MethodOther:             "Other",
			StatusPending:           "Pending",
			StatusCompleted:         "Completed",
			StatusFailed:            "Failed",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			NameInfo:        "Display name for this product line.",
			DescriptionInfo: "Optional notes about this product line.",
			ActiveInfo:      "Inactive product lines are hidden from new assignments.",
		},
		Actions: LineActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			MarkComplete: "Mark Complete",
			Reactivate:   "Reactivate",
		},
		Bulk: LineBulkLabels{
			Delete: "Delete",
		},
		Tabs: LineTabLabels{
			Info: "Info",
		},
		Detail: LineDetailLabels{
			TitlePrefix:          "Product Line ",
			PageTitle:            "Product Line",
			BasicInfo:            "Product Line Information",
			PaymentInfo:          "Payment Information",
			Reference:            "Reference",
			Customer:             "Customer",
			Amount:               "Amount",
			Currency:             "Currency",
			Method:               "Method",
			Date:                 "Date",
			Status:               "Status",
			Notes:                "Notes",
			CreatedDate:          "Created Date",
			ModifiedDate:         "Modified Date",
			ActiveBadge:          "Active",
			InactiveBadge:        "Inactive",
			TabBasicInfo:         "Info",
			TabAttachments:       "Attachments",
			TabAuditTrail:        "Audit Trail",
			TabAuditHistory:      "History",
			AuditAction:          "Action",
			AuditUser:            "User",
			AuditEmptyTitle:      "No audit entries",
			AuditEmptyMessage:    "No audit entries to display.",
			AuditTrailComingSoon: "Audit trail coming soon",
			AuditTrailDesc:       "Audit trail is not yet available for this product line.",
		},
		Status: LineStatusLabels{
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Confirm: LineConfirmLabels{
			MarkComplete:          "Mark Complete",
			MarkCompleteMessage:   "Are you sure you want to mark this product line as complete?",
			Reactivate:            "Reactivate",
			ReactivateMessage:     "Are you sure you want to reactivate this product line?",
			Delete:                "Delete Product Line",
			DeleteMessage:         "Are you sure you want to delete this product line?",
			BulkActivate:          "Activate",
			BulkActivateMessage:   "Are you sure you want to activate the selected product lines?",
			BulkDeactivate:        "Deactivate",
			BulkDeactivateMessage: "Are you sure you want to deactivate the selected product lines?",
			BulkComplete:          "Mark Complete",
			BulkCompleteMessage:   "Are you sure you want to mark {{count}} product line(s) as complete?",
			BulkReactivate:        "Reactivate",
			BulkReactivateMessage: "Are you sure you want to reactivate {{count}} product line(s)?",
			BulkDelete:            "Delete Product Lines",
			BulkDeleteMessage:     "Are you sure you want to delete {{count}} product line(s)?",
		},
		Errors: LineErrorLabels{
			PermissionDenied: "Permission denied",
			InvalidFormData:  "Invalid form data",
			NotFound:         "Product line not found",
			IDRequired:       "Product line ID is required",
			NoIDsProvided:    "No product line IDs provided",
			InvalidStatus:    "Invalid status",
			CannotDelete:     "This product line cannot be deleted because it is in use",
		},
	}
}

// DefaultLabels returns the zero-value label set. Every rendered string for
// this entity must come from the lyngua cascade (general -> business-type
// tier); there are no Go-side default strings to fall back on.
func DefaultLabels() Labels { return Labels{} }
