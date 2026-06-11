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

// ProductKindLabels holds the translated labels for each product_kind enum
// value. Sourced from lyngua product.json "productKind" block. Wired onto
// ProductLabels so the drawer-form select can render the per-value labels
// using the exact tier-cascaded strings that appear elsewhere in the UI.
type ProductKindLabels struct {
	Service        string `json:"service"`
	StockedGood    string `json:"stockedGood"`
	NonStockedGood string `json:"nonStockedGood"`
	Consumable     string `json:"consumable"`
}

// DeliveryModeLabels mirrors ProductKindLabels for the delivery_mode axis.
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

// ProductLabels holds all translatable strings for the product module.
type ProductLabels struct {
	Page       ProductPageLabels       `json:"page"`
	Buttons    ProductButtonLabels     `json:"buttons"`
	Columns    ProductColumnLabels     `json:"columns"`
	Empty      ProductEmptyLabels      `json:"empty"`
	Form       ProductFormLabels       `json:"form"`
	Actions    ProductActionLabels     `json:"actions"`
	Bulk       ProductBulkLabels       `json:"bulkActions"`
	Tabs       ProductTabLabels        `json:"tabs"`
	Detail     ProductDetailLabels     `json:"detail"`
	Status     ProductStatusLabels     `json:"status"`
	Variant    ProductVariantLabels    `json:"variant"`
	Attribute  ProductAttributeLabels  `json:"attribute"`
	Options    ProductOptionLabels     `json:"options"`
	Confirm    ProductConfirmLabels    `json:"confirm"`
	Errors     ProductErrorLabels      `json:"errors"`
	Breadcrumb ProductBreadcrumbLabels `json:"breadcrumb"`
	// Four-axis product taxonomy enum labels — loaded from lyngua
	// product.json "productKind"/"deliveryMode"/"trackingMode" blocks.
	// Wired here so the drawer-form select uses the exact tier-cascaded
	// display string for each enum value without hardcoding in Go.
	ProductKind  ProductKindLabels  `json:"productKind"`
	DeliveryMode DeliveryModeLabels `json:"deliveryMode"`
	TrackingMode TrackingModeLabels `json:"trackingMode"`

	// Phase 5 — service dashboard (product_kind=service surface).
	ServiceDashboard ServiceDashboardLabels `json:"serviceDashboard"`
}

// ServiceDashboardLabels holds translatable strings for the service
// dashboard. The "Service" wording is preferred at the dashboard surface
// because the sidebar key is "service"; the underlying entity is still
// Product filtered to product_kind="service".
type ServiceDashboardLabels struct {
	Title              string `json:"title"`
	Subtitle           string `json:"subtitle"`
	StatTotalActive    string `json:"statTotalActive"`
	StatTopRevenue     string `json:"statTopRevenue"`
	StatByLineCount    string `json:"statByLineCount"`
	StatRecentlyAdded  string `json:"statRecentlyAdded"`
	WidgetByLine       string `json:"widgetByLine"`
	WidgetTopRevenue   string `json:"widgetTopRevenue"`
	WidgetRecent       string `json:"widgetRecent"`
	QuickNew           string `json:"quickNew"`
	QuickBundleBuilder string `json:"quickBundleBuilder"`
	QuickTagService    string `json:"quickTagService"`
	QuickPriceSchedule string `json:"quickPriceSchedule"`
	ViewAll            string `json:"viewAll"`
	EmptyRecentTitle   string `json:"emptyRecentTitle"`
	EmptyRecentDesc    string `json:"emptyRecentDesc"`
	EmptyTopRevenue    string `json:"emptyTopRevenue"`
	NewService         string `json:"newService"`
	ColLine            string `json:"colLine"`
	ColRank            string `json:"colRank"`
	ColService         string `json:"colService"`
}

type ProductPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type ProductButtonLabels struct {
	AddProduct string `json:"addProduct"`
}

type ProductColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Line        string `json:"line"`
	Price       string `json:"price"`
	Status      string `json:"status"`
}

type ProductEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type ProductFormLabels struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DescPlaceholder string `json:"descriptionPlaceholder"`
	Price           string `json:"price"`
	Currency        string `json:"currency"`
	Active          string `json:"active"`
	Line            string `json:"line"`
	LinePlaceholder string `json:"linePlaceholder"`

	// Variant / option / attribute form labels
	PricePlaceholder       string `json:"pricePlaceholder"`
	SelectOption           string `json:"selectOption"`
	Required               string `json:"required"`
	Option                 string `json:"option"`
	SelectAttribute        string `json:"selectAttribute"`
	AllAttributesAssigned  string `json:"allAttributesAssigned"`
	OptionNeedsValuesAlert string `json:"optionNeedsValuesAlert"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	LineInfo        string `json:"lineInfo"`
	PriceInfo       string `json:"priceInfo"`
	CurrencyInfo    string `json:"currencyInfo"`
	ActiveInfo      string `json:"activeInfo"`

	// Model D — variant_mode toggle + unit field
	VariantModeLabel        string `json:"variantModeLabel"`
	VariantModeInfo         string `json:"variantModeInfo"`
	VariantModeNone         string `json:"variantModeNone"`
	VariantModeConfigurable string `json:"variantModeConfigurable"`
	UnitLabel               string `json:"unitLabel"`
	UnitInfo                string `json:"unitInfo"`
	UnitPlaceholder         string `json:"unitPlaceholder"`
	VariantPriceVaries      string `json:"variantPriceVaries"`
	// Shown as help text beneath the variant toggle when the product already
	// has option or variant rows, to explain why the toggle is disabled.
	VariantModeLockedHelp string `json:"variantModeLockedHelp"`
	// Error surfaced by the Create/Update handlers when a caller tries to
	// flip variant_mode on a product that still has options/variants.
	VariantModeLockedError string `json:"variantModeLockedError"`

	// Four-axis product taxonomy — rendered as selects on the drawer form.
	// Each axis carries its own Label + Info popover text plus per-enum-value
	// Info (XxxValueInfo map) keyed by enum string. When the mount restricts
	// the axis to one allowed value the select is rendered disabled so the
	// user still sees the classification without being able to change it.
	ProductKindLabel      string            `json:"productKindLabel"`
	ProductKindInfo       string            `json:"productKindInfo"`
	ProductKindValueInfo  map[string]string `json:"productKindValueInfo,omitempty"`
	DeliveryModeLabel     string            `json:"deliveryModeLabel"`
	DeliveryModeInfo      string            `json:"deliveryModeInfo"`
	DeliveryModeValueInfo map[string]string `json:"deliveryModeValueInfo,omitempty"`
	TrackingModeLabel     string            `json:"trackingModeLabel"`
	TrackingModeInfo      string            `json:"trackingModeInfo"`
	TrackingModeValueInfo map[string]string `json:"trackingModeValueInfo,omitempty"`

	// Tax section labels (Phase 5)
	SectionTax                  string `json:"sectionTax"`
	TaxTreatmentLabel           string `json:"taxTreatmentLabel"`
	TaxTreatmentPlaceholder     string `json:"taxTreatmentPlaceholder"`
	TaxTreatmentInfo            string `json:"taxTreatmentInfo"`
	WithholdingClassLabel       string `json:"withholdingClassLabel"`
	WithholdingClassPlaceholder string `json:"withholdingClassPlaceholder"`
	WithholdingClassInfo        string `json:"withholdingClassInfo"`
}

type ProductActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type ProductBulkLabels struct {
	Delete string `json:"delete"`
}

type ProductTabLabels struct {
	Info         string `json:"info"`
	Variants     string `json:"variants"`
	Attributes   string `json:"attributes"`
	Pricing      string `json:"pricing"`
	Options      string `json:"options"`
	Images       string `json:"images"`
	Stock        string `json:"stock"`
	Lines        string `json:"lines"`
	Attachments  string `json:"attachments"`
	AuditTrail   string `json:"auditTrail"`
	AuditHistory string `json:"auditHistory"`
	// Inventory item sub-tabs
	Serials        string `json:"serials"`
	PricingHistory string `json:"pricingHistory"`
}

type ProductDetailLabels struct {
	Price                string `json:"price"`
	Currency             string `json:"currency"`
	Collections          string `json:"collections"`
	VariantCount         string `json:"variantCount"`
	Status               string `json:"status"`
	OptionsLabel         string `json:"optionsLabel"`
	EmptyVariantsMessage string `json:"emptyVariantsMessage"`
	// Header subtitle fallback when the product has no description.
	// Consumed by buildPageData to override the generic "Welcome back"
	// CommonLabels default on the product detail page header.
	NoDescriptionSubtitle string `json:"noDescriptionSubtitle"`
	// Model D — detail-page rows for unit of measure + variant mode.
	// Falls back to English defaults when lyngua doesn't overlay the key.
	Unit        string `json:"unit"`
	VariantMode string `json:"variantMode"`
	// Serial table columns
	SerialNumber       string `json:"serialNumber"`
	IMEI               string `json:"imei"`
	WarrantyEnd        string `json:"warrantyEnd"`
	PurchaseOrder      string `json:"purchaseOrder"`
	NoSerialNumbers    string `json:"noSerialNumbers"`
	NoSerialNumbersMsg string `json:"noSerialNumbersMsg"`

	// Variant detail labels
	VariantInformation  string `json:"variantInformation"`
	Options             string `json:"options"`
	VariantPricing      string `json:"variantPricing"`
	VariantPricingDesc  string `json:"variantPricingDesc"`
	InventoryStock      string `json:"inventoryStock"`
	InventoryStockDesc  string `json:"inventoryStockDesc"`
	DropImagesHere      string `json:"dropImagesHere"`
	ImageFileHint       string `json:"imageFileHint"`
	DeleteSelected      string `json:"deleteSelected"`
	PrimaryBadge        string `json:"primaryBadge"`
	NoImages            string `json:"noImages"`
	NoImagesDesc        string `json:"noImagesDesc"`
	AuditTrail          string `json:"auditTrail"`
	AuditTrailDesc      string `json:"auditTrailDesc"`
	NoSerialNumbersDesc string `json:"noSerialNumbersDesc"`

	// Stock detail labels
	InventoryItem      string `json:"inventoryItem"`
	Name               string `json:"name"`
	SKU                string `json:"sku"`
	Type               string `json:"type"`
	Location           string `json:"location"`
	QtyOnHand          string `json:"qtyOnHand"`
	Reserved           string `json:"reserved"`
	Available          string `json:"available"`
	StatTotal          string `json:"statTotal"`
	StatAvailable      string `json:"statAvailable"`
	StatSold           string `json:"statSold"`
	StatReserved       string `json:"statReserved"`
	PricingHistory     string `json:"pricingHistory"`
	PricingHistoryDesc string `json:"pricingHistoryDesc"`

	// Serial detail labels
	SerialInformation string `json:"serialInformation"`
}

type ProductStatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ProductVariantLabels struct {
	Title         string `json:"title"`
	SKU           string `json:"sku"`
	PriceOverride string `json:"priceOverride"`
	Attributes    string `json:"attributes"`
	Assign        string `json:"assign"`
	Edit          string `json:"edit"`
	Remove        string `json:"remove"`
	Empty         string `json:"empty"`
	// Stock table columns
	Location    string `json:"location"`
	QtyOnHand   string `json:"qtyOnHand"`
	SerialCount string `json:"serialCount"`
	NoStock     string `json:"noStock"`
	NoStockMsg  string `json:"noStockMsg"`
	// Pricing tab column headers
	Pricing VariantPricingLabels `json:"pricing"`
}

// VariantPricingLabels holds column header labels for the variant pricing tab table.
type VariantPricingLabels struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Package  string `json:"package"`
	RateCard string `json:"rateCard"`
	Amount   string `json:"amount"`
}

type ProductAttributeLabels struct {
	Title        string `json:"title"`
	DefaultValue string `json:"defaultValue"`
	Assign       string `json:"assign"`
	Remove       string `json:"remove"`
	Empty        string `json:"empty"`
}

type ProductConfirmLabels struct {
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activateMessage"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivateMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
	RemoveVariant         string `json:"removeVariant"`
	RemoveVariantMessage  string `json:"removeVariantMessage"`
}

type ProductErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
	CannotDelete     string `json:"cannotDelete"`
	NameRequired     string `json:"nameRequired"`
	FieldRequired    string `json:"fieldRequired"`
}

type ProductBreadcrumbLabels struct {
	Products string `json:"products"`
	Product  string `json:"product"`
	Option   string `json:"option"`
}

// ---------------------------------------------------------------------------
// Product Option labels
// ---------------------------------------------------------------------------

type ProductOptionLabels struct {
	Tab       ProductOptionTabLabels      `json:"tab"`
	Tabs      ProductOptionTabsLabels     `json:"tabs"`
	Columns   ProductOptionColumnLabels   `json:"columns"`
	Form      ProductOptionFormLabels     `json:"form"`
	DataTypes ProductOptionDataTypeLabels `json:"dataTypes"`
	Value     ProductOptionValueLabels    `json:"value"`
	Actions   ProductOptionActionLabels   `json:"actions"`
	Empty     ProductOptionEmptyLabels    `json:"empty"`
	Confirm   ProductOptionConfirmLabels  `json:"confirm"`
}

type ProductOptionTabLabels struct {
	Title string `json:"title"`
}

type ProductOptionTabsLabels struct {
	Info   string `json:"info"`
	Values string `json:"values"`
}

type ProductOptionColumnLabels struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	DataType    string `json:"dataType"`
	ValuesCount string `json:"valuesCount"`
	SortOrder   string `json:"sortOrder"`
	Required    string `json:"required"`
	Status      string `json:"status"`
}

type ProductOptionFormLabels struct {
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	Code                   string `json:"code"`
	CodePlaceholder        string `json:"codePlaceholder"`
	DataType               string `json:"dataType"`
	SortOrder              string `json:"sortOrder"`
	MinValue               string `json:"minValue"`
	MaxValue               string `json:"maxValue"`
	Active                 string `json:"active"`
	Required               string `json:"required"`
	RequiredCaution        string `json:"requiredCaution"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	DescriptionEmpty       string `json:"descriptionEmpty"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	CodeInfo        string `json:"codeInfo"`
	DataTypeInfo    string `json:"dataTypeInfo"`
	MinValueInfo    string `json:"minValueInfo"`
	MaxValueInfo    string `json:"maxValueInfo"`
	SortOrderInfo   string `json:"sortOrderInfo"`
	ActiveInfo      string `json:"activeInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
}

type ProductOptionDataTypeLabels struct {
	TextList    string `json:"textList"`
	NumberRange string `json:"numberRange"`
	ColorList   string `json:"colorList"`
	FreeText    string `json:"freeText"`
	FreeNumber  string `json:"freeNumber"`
}

type ProductOptionValueLabels struct {
	Columns ProductOptionValueColumnLabels `json:"columns"`
	Form    ProductOptionValueFormLabels   `json:"form"`
}

type ProductOptionValueColumnLabels struct {
	Label        string `json:"label"`
	Value        string `json:"value"`
	SortOrder    string `json:"sortOrder"`
	ColorPreview string `json:"colorPreview"`
	Status       string `json:"status"`
}

type ProductOptionValueFormLabels struct {
	Label               string `json:"label"`
	LabelPlaceholder    string `json:"labelPlaceholder"`
	Value               string `json:"value"`
	ValuePlaceholder    string `json:"valuePlaceholder"`
	SortOrder           string `json:"sortOrder"`
	ColorHex            string `json:"colorHex"`
	ColorHexPlaceholder string `json:"colorHexPlaceholder"`
	Active              string `json:"active"`
	// Context labels surfaced on the value drawer to remind the user
	// which option this value belongs to.
	Option   string `json:"option"`
	Required string `json:"required"`

	// Field-level info text surfaced via an info button beside each label.
	LabelInfo     string `json:"labelInfo"`
	ValueInfo     string `json:"valueInfo"`
	SortOrderInfo string `json:"sortOrderInfo"`
	ColorHexInfo  string `json:"colorHexInfo"`
	ActiveInfo    string `json:"activeInfo"`
}

type ProductOptionActionLabels struct {
	AddOption         string `json:"addOption"`
	EditOption        string `json:"editOption"`
	EditProductOption string `json:"editProductOption"`
	DeleteOption      string `json:"deleteOption"`
	ViewValues        string `json:"viewValues"`
	AddValue          string `json:"addValue"`
	EditValue         string `json:"editValue"`
	DeleteValue       string `json:"deleteValue"`
}

type ProductOptionEmptyLabels struct {
	Title        string `json:"title"`
	Message      string `json:"message"`
	ValueTitle   string `json:"valueTitle"`
	ValueMessage string `json:"valueMessage"`
}

type ProductOptionConfirmLabels struct {
	DeleteOption string `json:"deleteOption"`
	DeleteValue  string `json:"deleteValue"`
}

// ---------------------------------------------------------------------------
// Price List labels
// ---------------------------------------------------------------------------

// PriceListLabels holds all translatable strings for the price list module.
type PriceListLabels struct {
	Page    PriceListPageLabels    `json:"page"`
	Buttons PriceListButtonLabels  `json:"buttons"`
	Columns PriceListColumnLabels  `json:"columns"`
	Empty   PriceListEmptyLabels   `json:"empty"`
	Form    PriceListFormLabels    `json:"form"`
	Actions PriceListActionLabels  `json:"actions"`
	Bulk    PriceListBulkLabels    `json:"bulkActions"`
	Detail  PriceListDetailLabels  `json:"detail"`
	Confirm PriceListConfirmLabels `json:"confirm"`
	Errors  PriceListErrorLabels   `json:"errors"`
}

type PriceListPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type PriceListButtonLabels struct {
	AddPriceList string `json:"addPriceList"`
}

type PriceListColumnLabels struct {
	Name      string `json:"name"`
	DateStart string `json:"dateStart"`
	DateEnd   string `json:"dateEnd"`
	Status    string `json:"status"`
}

type PriceListEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type PriceListFormLabels struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DescPlaceholder string `json:"descriptionPlaceholder"`
	DateStart       string `json:"dateStart"`
	DateEnd         string `json:"dateEnd"`
	Active          string `json:"active"`
	Product         string `json:"product"`
	SelectProduct   string `json:"selectProduct"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	DateStartInfo   string `json:"dateStartInfo"`
	DateEndInfo     string `json:"dateEndInfo"`
	ActiveInfo      string `json:"activeInfo"`
	// Price-product sub-drawer info fields.
	AmountInfo   string `json:"amountInfo"`
	CurrencyInfo string `json:"currencyInfo"`
}

type PriceListActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type PriceListBulkLabels struct {
	Delete string `json:"delete"`
}

type PriceListDetailLabels struct {
	PageTitle          string `json:"pageTitle"`
	BasicInfo          string `json:"basicInfo"`
	Prices             string `json:"prices"`
	TabAttachments     string `json:"tabAttachments"`
	TabAuditHistory    string `json:"tabAuditHistory"`
	ProductName        string `json:"productName"`
	Amount             string `json:"amount"`
	Currency           string `json:"currency"`
	AddPrice           string `json:"addPrice"`
	RemoveLabel        string `json:"removeLabel"`
	EmptyTitle         string `json:"emptyTitle"`
	EmptyMessage       string `json:"emptyMessage"`
	ActiveBadge        string `json:"activeBadge"`
	InactiveBadge      string `json:"inactiveBadge"`
	NoPricesConfigured string `json:"noPricesConfigured"`
	NoPricesDesc       string `json:"noPricesDesc"`
}

type PriceListConfirmLabels struct {
	Activate          string `json:"activate"`
	ActivateMessage   string `json:"activateMessage"`
	Deactivate        string `json:"deactivate"`
	DeactivateMessage string `json:"deactivateMessage"`
	Delete            string `json:"delete"`
	DeleteMessage     string `json:"deleteMessage"`
	BulkDelete        string `json:"bulkDelete"`
	BulkDeleteMessage string `json:"bulkDeleteMessage"`
}

type PriceListErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	CannotDelete     string `json:"cannotDelete"`
	ProductRequired  string `json:"productRequired"`
	AmountRequired   string `json:"amountRequired"`
}

// ---------------------------------------------------------------------------
// Product Line labels
// ---------------------------------------------------------------------------

// ProductLineLabels holds all translatable strings for the product line module.
type ProductLineLabels struct {
	Page    ProductLinePageLabels    `json:"page"`
	Buttons ProductLineButtonLabels  `json:"buttons"`
	Columns ProductLineColumnLabels  `json:"columns"`
	Empty   ProductLineEmptyLabels   `json:"empty"`
	Form    ProductLineFormLabels    `json:"form"`
	Actions ProductLineActionLabels  `json:"actions"`
	Bulk    ProductLineBulkLabels    `json:"bulkActions"`
	Tabs    ProductLineTabLabels     `json:"tabs"`
	Detail  ProductLineDetailLabels  `json:"detail"`
	Status  ProductLineStatusLabels  `json:"status"`
	Confirm ProductLineConfirmLabels `json:"confirm"`
	Errors  ProductLineErrorLabels   `json:"errors"`
}

type ProductLinePageLabels struct {
	Heading          string `json:"heading"`
	HeadingActive    string `json:"headingActive"`
	HeadingInactive  string `json:"headingInactive"`
	HeadingPending   string `json:"headingPending"`
	HeadingCompleted string `json:"headingCompleted"`
	HeadingFailed    string `json:"headingFailed"`
	Caption          string `json:"caption"`
	CaptionActive    string `json:"captionActive"`
	CaptionInactive  string `json:"captionInactive"`
	CaptionPending   string `json:"captionPending"`
	CaptionCompleted string `json:"captionCompleted"`
	CaptionFailed    string `json:"captionFailed"`
}

type ProductLineButtonLabels struct {
	AddProductLine    string `json:"addProductLine"`
	EditProductLine   string `json:"editProductLine"`
	DeleteProductLine string `json:"deleteProductLine"`
}

type ProductLineColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DateCreated string `json:"dateCreated"`
	Reference   string `json:"reference"`
	Customer    string `json:"customer"`
	Amount      string `json:"amount"`
	Method      string `json:"method"`
	Date        string `json:"date"`
	Status      string `json:"status"`
}

type ProductLineEmptyLabels struct {
	ActiveTitle      string `json:"activeTitle"`
	ActiveMessage    string `json:"activeMessage"`
	InactiveTitle    string `json:"inactiveTitle"`
	InactiveMessage  string `json:"inactiveMessage"`
	PendingTitle     string `json:"pendingTitle"`
	PendingMessage   string `json:"pendingMessage"`
	CompletedTitle   string `json:"completedTitle"`
	CompletedMessage string `json:"completedMessage"`
	FailedTitle      string `json:"failedTitle"`
	FailedMessage    string `json:"failedMessage"`
}

type ProductLineFormLabels struct {
	Name                    string `json:"name"`
	NamePlaceholder         string `json:"namePlaceholder"`
	Description             string `json:"description"`
	DescPlaceholder         string `json:"descriptionPlaceholder"`
	Active                  string `json:"active"`
	Customer                string `json:"customer"`
	Date                    string `json:"date"`
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	Reference               string `json:"reference"`
	ReferencePlaceholder    string `json:"referencePlaceholder"`
	PaymentMethod           string `json:"paymentMethod"`
	Status                  string `json:"status"`
	Notes                   string `json:"notes"`
	NotesPlaceholder        string `json:"notesPlaceholder"`
	CustomerNamePlaceholder string `json:"customerNamePlaceholder"`
	AmountPlaceholder       string `json:"amountPlaceholder"`
	CurrencyPlaceholder     string `json:"currencyPlaceholder"`
	MethodCash              string `json:"methodCash"`
	MethodBankTransfer      string `json:"methodBankTransfer"`
	MethodCheck             string `json:"methodCheck"`
	MethodGCash             string `json:"methodGCash"`
	MethodMaya              string `json:"methodMaya"`
	MethodCard              string `json:"methodCard"`
	MethodOther             string `json:"methodOther"`
	StatusPending           string `json:"statusPending"`
	StatusCompleted         string `json:"statusCompleted"`
	StatusFailed            string `json:"statusFailed"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	ActiveInfo      string `json:"activeInfo"`
}

type ProductLineActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	MarkComplete string `json:"markComplete"`
	Reactivate   string `json:"reactivate"`
}

type ProductLineBulkLabels struct {
	Delete string `json:"delete"`
}

type ProductLineStatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ProductLineTabLabels struct {
	Info string `json:"info"`
}

type ProductLineDetailLabels struct {
	TitlePrefix          string `json:"titlePrefix"`
	PageTitle            string `json:"pageTitle"`
	BasicInfo            string `json:"basicInfo"`
	PaymentInfo          string `json:"paymentInfo"`
	Reference            string `json:"reference"`
	Customer             string `json:"customer"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Method               string `json:"method"`
	Date                 string `json:"date"`
	Status               string `json:"status"`
	Notes                string `json:"notes"`
	CreatedDate          string `json:"createdDate"`
	ModifiedDate         string `json:"modifiedDate"`
	ActiveBadge          string `json:"activeBadge"`
	InactiveBadge        string `json:"inactiveBadge"`
	TabBasicInfo         string `json:"tabBasicInfo"`
	TabAttachments       string `json:"tabAttachments"`
	TabAuditTrail        string `json:"tabAuditTrail"`
	TabAuditHistory      string `json:"tabAuditHistory"`
	AuditAction          string `json:"auditAction"`
	AuditUser            string `json:"auditUser"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditTrailDesc       string `json:"auditTrailDesc"`
}

type ProductLineConfirmLabels struct {
	MarkComplete          string `json:"markComplete"`
	MarkCompleteMessage   string `json:"markCompleteMessage"`
	Reactivate            string `json:"reactivate"`
	ReactivateMessage     string `json:"reactivateMessage"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
	BulkComplete          string `json:"bulkComplete"`
	BulkCompleteMessage   string `json:"bulkCompleteMessage"`
	BulkReactivate        string `json:"bulkReactivate"`
	BulkReactivateMessage string `json:"bulkReactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
}

type ProductLineErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
	CannotDelete     string `json:"cannotDelete"`
}

// DefaultProductLineLabels returns ProductLineLabels with sensible English defaults.
func DefaultProductLineLabels() ProductLineLabels {
	return ProductLineLabels{
		Page: ProductLinePageLabels{
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
		Buttons: ProductLineButtonLabels{
			AddProductLine:    "Add Product Line",
			EditProductLine:   "Edit Product Line",
			DeleteProductLine: "Delete Product Line",
		},
		Columns: ProductLineColumnLabels{
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
		Empty: ProductLineEmptyLabels{
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
		Form: ProductLineFormLabels{
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
		Actions: ProductLineActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			MarkComplete: "Mark Complete",
			Reactivate:   "Reactivate",
		},
		Bulk: ProductLineBulkLabels{
			Delete: "Delete",
		},
		Tabs: ProductLineTabLabels{
			Info: "Info",
		},
		Detail: ProductLineDetailLabels{
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
		Status: ProductLineStatusLabels{
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Confirm: ProductLineConfirmLabels{
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
		Errors: ProductLineErrorLabels{
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

// ---------------------------------------------------------------------------
// Resource labels
// ---------------------------------------------------------------------------

// ResourceLabels holds all translatable strings for the resource module.
type ResourceLabels struct {
	Page    ResourcePageLabels    `json:"page"`
	Buttons ResourceButtonLabels  `json:"buttons"`
	Columns ResourceColumnLabels  `json:"columns"`
	Empty   ResourceEmptyLabels   `json:"empty"`
	Form    ResourceFormLabels    `json:"form"`
	Actions ResourceActionLabels  `json:"actions"`
	Bulk    ResourceBulkLabels    `json:"bulkActions"`
	Status  ResourceStatusLabels  `json:"status"`
	Confirm ResourceConfirmLabels `json:"confirm"`
	Errors  ResourceErrorLabels   `json:"errors"`
}

type ResourcePageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type ResourceButtonLabels struct {
	Add string `json:"add"`
}

type ResourceColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Product     string `json:"product"`
	Status      string `json:"status"`
}

type ResourceEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ResourceFormLabels struct {
	Name            string `json:"name"`
	NamePlaceholder string `json:"namePlaceholder"`
	Description     string `json:"description"`
	DescPlaceholder string `json:"descriptionPlaceholder"`
	ProductId       string `json:"productId"`
	UserId          string `json:"userId"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	ProductIdInfo   string `json:"productIdInfo"`
	UserIdInfo      string `json:"userIdInfo"`
}

type ResourceActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ResourceBulkLabels struct {
	Delete string `json:"delete"`
}

type ResourceStatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ResourceConfirmLabels struct {
	Delete              string `json:"delete"`
	DeleteMessage       string `json:"deleteMessage"`
	Activate            string `json:"activate"`
	ActivateMessage     string `json:"activateMessage"`
	Deactivate          string `json:"deactivate"`
	DeactivateMessage   string `json:"deactivateMessage"`
	BulkDelete          string `json:"bulkDelete"`
	BulkDeleteMessage   string `json:"bulkDeleteMessage"`
	BulkActivate        string `json:"bulkActivate"`
	BulkActivateMessage string `json:"bulkActivateMessage"`
}

type ResourceErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
	CannotDelete     string `json:"cannotDelete"`
}

// DefaultResourceLabels returns ResourceLabels with sensible English defaults.
func DefaultResourceLabels() ResourceLabels {
	return ResourceLabels{
		Page: ResourcePageLabels{
			Heading:         "Resources",
			HeadingActive:   "Active Resources",
			HeadingInactive: "Inactive Resources",
			Caption:         "Manage resources linked to products.",
			CaptionActive:   "Showing active resources.",
			CaptionInactive: "Showing inactive resources.",
		},
		Buttons: ResourceButtonLabels{
			Add: "Add Resource",
		},
		Columns: ResourceColumnLabels{
			Name:        "Name",
			Description: "Description",
			Product:     "Product",
			Status:      "Status",
		},
		Empty: ResourceEmptyLabels{
			Title:   "No resources found",
			Message: "Add a resource to get started.",
		},
		Form: ResourceFormLabels{
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
		Actions: ResourceActionLabels{
			View:       "View",
			Edit:       "Edit",
			Delete:     "Delete",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Bulk: ResourceBulkLabels{
			Delete: "Delete Selected",
		},
		Status: ResourceStatusLabels{
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Confirm: ResourceConfirmLabels{
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
		Errors: ResourceErrorLabels{
			PermissionDenied: "You do not have permission to perform this action",
			InvalidFormData:  "Invalid form data. Please check your inputs and try again.",
			NotFound:         "Resource not found",
			IDRequired:       "Resource ID is required",
			NoPermission:     "No permission",
			CannotDelete:     "This resource cannot be deleted because it is in use",
		},
	}
}
