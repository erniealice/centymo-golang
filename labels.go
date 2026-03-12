package centymo

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// ---------------------------------------------------------------------------
// Inventory labels
// ---------------------------------------------------------------------------

// InventoryLabels holds all translatable strings for the inventory module.
type InventoryLabels struct {
	Page         InventoryPageLabels         `json:"page"`
	Buttons      InventoryButtonLabels       `json:"buttons"`
	Columns      InventoryColumnLabels       `json:"columns"`
	Empty        InventoryEmptyLabels        `json:"empty"`
	Form         InventoryFormLabels         `json:"form"`
	Actions      InventoryActionLabels       `json:"actions"`
	Bulk         InventoryBulkLabels         `json:"bulkActions"`
	Detail       InventoryDetailLabels       `json:"detail"`
	Tabs         InventoryTabLabels          `json:"tabs"`
	ItemType     InventoryItemTypeLabels     `json:"itemType"`
	Status       InventoryStatusLabels       `json:"status"`
	Serial       InventorySerialLabels       `json:"serial"`
	Transaction  InventoryTransactionLabels  `json:"transaction"`
	Depreciation InventoryDepreciationLabels `json:"depreciation"`
	Dashboard    InventoryDashboardLabels    `json:"dashboard"`
	Movements    InventoryMovementsLabels    `json:"movements"`
	Confirm      InventoryConfirmLabels      `json:"confirm"`
	Errors       InventoryErrorLabels        `json:"errors"`
	Breadcrumb   InventoryBreadcrumbLabels   `json:"breadcrumb"`
}

type InventoryPageLabels struct {
	Heading  string `json:"heading"`
	Caption  string `json:"caption"`
	Location string `json:"location"`
}

type InventoryButtonLabels struct {
	AddItem string `json:"addItem"`
}

type InventoryColumnLabels struct {
	ProductName string `json:"productName"`
	SKU         string `json:"sku"`
	OnHand      string `json:"onHand"`
	Available   string `json:"available"`
	ReorderLvl  string `json:"reorderLevel"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}

type InventoryEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type InventoryFormLabels struct {
	Product       string `json:"product"`
	SKU           string `json:"sku"`
	SKUPlaceholder string `json:"skuPlaceholder"`
	OnHand        string `json:"onHand"`
	Reserved      string `json:"reserved"`
	ReorderLevel  string `json:"reorderLevel"`
	UnitOfMeasure string `json:"unitOfMeasure"`
	Notes         string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	Active        string `json:"active"`
}

type InventoryActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type InventoryBulkLabels struct {
	Delete string `json:"delete"`
}

// InventoryDetailLabels holds all translatable strings for the inventory detail page.
type InventoryDetailLabels struct {
	TitlePrefix string `json:"titlePrefix"`
	MonthsUnit  string `json:"monthsUnit"`
	TypesUnit   string `json:"typesUnit"`

	// Tab labels
	TabBasicInfo    string `json:"tabBasicInfo"`
	TabAttributes   string `json:"tabAttributes"`
	TabSerials      string `json:"tabSerials"`
	TabTransactions string `json:"tabTransactions"`
	TabAuditTrail   string `json:"tabAuditTrail"`

	// Info fields
	ItemInfo      string `json:"itemInfo"`
	ProductName   string `json:"productName"`
	SKU           string `json:"sku"`
	Location      string `json:"location"`
	OnHand        string `json:"onHand"`
	Reserved      string `json:"reserved"`
	Available     string `json:"available"`
	ReorderLevel  string `json:"reorderLevel"`
	UnitOfMeasure string `json:"unitOfMeasure"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`

	// Attribute labels
	AttributeName  string `json:"attributeName"`
	AttributeValue string `json:"attributeValue"`

	// Serial columns
	SerialNumber  string `json:"serialNumber"`
	IMEI          string `json:"imei"`
	SerialStatus  string `json:"serialStatus"`
	WarrantyEnd   string `json:"warrantyEnd"`
	PurchaseOrder string `json:"purchaseOrder"`
	SaleReference string `json:"saleReference"`

	// Serial summary
	TotalUnits     string `json:"totalUnits"`
	AvailableUnits string `json:"availableUnits"`
	SoldUnits      string `json:"soldUnits"`
	ReservedUnits  string `json:"reservedUnits"`

	// Transaction columns
	Date        string `json:"date"`
	Type        string `json:"type"`
	Quantity    string `json:"quantity"`
	Reference   string `json:"reference"`
	Serial      string `json:"serial"`
	PerformedBy string `json:"performedBy"`

	// Audit columns
	AuditAction string `json:"auditAction"`
	AuditUser   string `json:"auditUser"`
	Description string `json:"description"`

	// Info field labels (not shared with transaction columns)
	Product       string `json:"product"`
	ViewProduct   string `json:"viewProduct"`
	Active        string `json:"active"`
	Inactive      string `json:"inactive"`
	SerialNumbers string `json:"serialNumbers"`

	// Empty states
	AttributeEmptyTitle     string `json:"attributeEmptyTitle"`
	AttributeEmptyMessage   string `json:"attributeEmptyMessage"`
	SerialEmptyTitle        string `json:"serialEmptyTitle"`
	SerialEmptyMessage      string `json:"serialEmptyMessage"`
	TransactionEmptyTitle   string `json:"transactionEmptyTitle"`
	TransactionEmptyMessage string `json:"transactionEmptyMessage"`
	AuditEmptyTitle         string `json:"auditEmptyTitle"`
	AuditEmptyMessage       string `json:"auditEmptyMessage"`
}

type InventoryTabLabels struct {
	Info         string `json:"info"`
	Attributes   string `json:"attributes"`
	Serials      string `json:"serials"`
	Transactions string `json:"transactions"`
	Depreciation string `json:"depreciation"`
	Audit        string `json:"audit"`
}

type InventoryItemTypeLabels struct {
	Serialized    string `json:"serialized"`
	NonSerialized string `json:"nonSerialized"`
	Consumable    string `json:"consumable"`
}

type InventoryStatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type InventorySerialLabels struct {
	Title           string `json:"title"`
	SerialNumber    string `json:"serialNumber"`
	IMEI            string `json:"imei"`
	Status          string `json:"status"`
	WarrantyStart   string `json:"warrantyStart"`
	WarrantyEnd     string `json:"warrantyEnd"`
	PurchaseOrder   string `json:"purchaseOrder"`
	SoldReference   string `json:"soldReference"`
	Assign          string `json:"assign"`
	Edit            string `json:"edit"`
	Remove          string `json:"remove"`
	Empty           string `json:"empty"`
	StatusAvailable string `json:"statusAvailable"`
	StatusSold      string `json:"statusSold"`
	StatusReserved  string `json:"statusReserved"`
	StatusDefective string `json:"statusDefective"`
	StatusReturned  string `json:"statusReturned"`
}

type InventoryTransactionLabels struct {
	Title           string `json:"title"`
	Type            string `json:"type"`
	Quantity        string `json:"quantity"`
	Date            string `json:"date"`
	Reference       string `json:"reference"`
	PerformedBy     string `json:"performedBy"`
	Record          string `json:"record"`
	Empty           string `json:"empty"`
	TypeReceived    string `json:"typeReceived"`
	TypeSold        string `json:"typeSold"`
	TypeAdjusted    string `json:"typeAdjusted"`
	TypeTransferred string `json:"typeTransferred"`
	TypeReturned    string `json:"typeReturned"`
	TypeWriteOff    string `json:"typeWriteOff"`
}

type InventoryDepreciationLabels struct {
	Title                  string `json:"title"`
	Method                 string `json:"method"`
	CostBasis              string `json:"costBasis"`
	SalvageValue           string `json:"salvageValue"`
	UsefulLife             string `json:"usefulLife"`
	StartDate              string `json:"startDate"`
	Accumulated            string `json:"accumulated"`
	BookValue              string `json:"bookValue"`
	Configure              string `json:"configure"`
	Edit                   string `json:"edit"`
	NotConfigured          string `json:"notConfigured"`
	MethodStraightLine     string `json:"methodStraightLine"`
	MethodDecliningBalance string `json:"methodDecliningBalance"`
	MethodSumOfYears       string `json:"methodSumOfYears"`
	MonthsUnit             string `json:"monthsUnit"`
}

type InventoryDashboardLabels struct {
	Title                string `json:"title"`
	TotalStockValue      string `json:"totalStockValue"`
	LowStockAlerts       string `json:"lowStockAlerts"`
	StockTurnover        string `json:"stockTurnover"`
	ItemsByLocation      string `json:"itemsByLocation"`
	DepreciationSummary  string `json:"depreciationSummary"`
	SerialUnitStatus     string `json:"serialUnitStatus"`
	RecentMovements      string `json:"recentMovements"`
	CategoryDistribution string `json:"categoryDistribution"`
	TypesUnit            string `json:"typesUnit"`
}

type InventoryMovementsLabels struct {
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle"`
	DateRange      string `json:"dateRange"`
	LocationFilter string `json:"locationFilter"`
	TypeFilter     string `json:"typeFilter"`
	ProductSearch  string `json:"productSearch"`
	ClearAll       string `json:"clearAll"`
	ExportCsv      string `json:"exportCsv"`
	AllLocations   string `json:"allLocations"`
	AllTypes       string `json:"allTypes"`
	ProductColumn  string `json:"productColumn"`
	VariantSKU     string `json:"variantSku"`
}

type InventoryConfirmLabels struct {
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activateMessage"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivateMessage"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
}

type InventoryErrorLabels struct {
	PermissionDenied  string `json:"permissionDenied"`
	InvalidFormData   string `json:"invalidFormData"`
	NotFound          string `json:"notFound"`
	IDRequired        string `json:"idRequired"`
	NoIDsProvided     string `json:"noIDsProvided"`
	InvalidStatus     string `json:"invalidStatus"`
	NoPermission      string `json:"noPermission"`
	SerialNotFound    string `json:"serialNotFound"`
	SerialIDRequired  string `json:"serialIDRequired"`
	InvalidDepreciationMethod string `json:"invalidDepreciationMethod"`
}

type InventoryBreadcrumbLabels struct {
	Products string `json:"products"`
	Product  string `json:"product"`
}

// ---------------------------------------------------------------------------
// Revenue labels
// ---------------------------------------------------------------------------

// RevenueLabels holds all translatable strings for the revenue module.
type RevenueLabels struct {
	Page    RevenuePageLabels    `json:"page"`
	Buttons RevenueButtonLabels  `json:"buttons"`
	Columns RevenueColumnLabels  `json:"columns"`
	Empty   RevenueEmptyLabels   `json:"empty"`
	Form    RevenueFormLabels    `json:"form"`
	Actions RevenueActionLabels  `json:"actions"`
	Bulk    RevenueBulkLabels    `json:"bulkActions"`
	Detail  RevenueDetailLabels  `json:"detail"`
	Confirm RevenueConfirmLabels `json:"confirm"`
	Errors  RevenueErrorLabels   `json:"errors"`
	Dashboard RevenueDashboardLabels `json:"dashboard"`
	Settings  RevenueSettingsLabels `json:"settings"`
}

type RevenuePageLabels struct {
	Heading          string `json:"heading"`
	HeadingOngoing   string `json:"headingOngoing"`
	HeadingComplete  string `json:"headingComplete"`
	HeadingCancelled string `json:"headingCancelled"`
	Caption          string `json:"caption"`
	CaptionOngoing   string `json:"captionOngoing"`
	CaptionComplete  string `json:"captionComplete"`
	CaptionCancelled string `json:"captionCancelled"`
}

type RevenueButtonLabels struct {
	AddSale string `json:"addSale"`
}

type RevenueColumnLabels struct {
	Reference  string `json:"reference"`
	Customer   string `json:"customer"`
	Date       string `json:"date"`
	Amount     string `json:"amount"`
	Status     string `json:"status"`
}

type RevenueEmptyLabels struct {
	OngoingTitle     string `json:"ongoingTitle"`
	OngoingMessage   string `json:"ongoingMessage"`
	CompleteTitle    string `json:"completeTitle"`
	CompleteMessage  string `json:"completeMessage"`
	CancelledTitle   string `json:"cancelledTitle"`
	CancelledMessage string `json:"cancelledMessage"`
}

type RevenueFormLabels struct {
	Customer             string `json:"customer"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Reference            string `json:"reference"`
	ReferencePlaceholder string `json:"referencePlaceholder"`
	Status               string `json:"status"`
	Notes                string `json:"notes"`
	NotesPlaceholder     string `json:"notesPlaceholder"`
	Active               string `json:"active"`
	Location             string `json:"location"`
}

type RevenueActionLabels struct {
	View            string `json:"view"`
	Edit            string `json:"edit"`
	Delete          string `json:"delete"`
	Complete        string `json:"complete"`
	Reactivate      string `json:"reactivate"`
	DownloadInvoice string `json:"downloadInvoice"`
	SendEmail       string `json:"sendEmail"`
}

type RevenueBulkLabels struct {
	Delete string `json:"delete"`
}

type RevenueDetailLabels struct {
	PageTitle   string `json:"pageTitle"`
	TitlePrefix string `json:"titlePrefix"`
	InvoiceInfo string `json:"invoiceInfo"`
	LineItems   string `json:"lineItems"`
	Description string `json:"description"`
	Quantity    string `json:"quantity"`
	UnitPrice   string `json:"unitPrice"`
	CostPrice   string `json:"costPrice"`
	GrossProfit string `json:"grossProfit"`
	Total       string `json:"total"`
	Discount    string `json:"discount"`
	SubTotal    string `json:"subTotal"`
	GrandTotal  string `json:"grandTotal"`

	// Tab labels
	TabBasicInfo  string `json:"tabBasicInfo"`
	TabLineItems  string `json:"tabLineItems"`
	TabPayment    string `json:"tabPayment"`
	TabAuditTrail string `json:"tabAuditTrail"`

	// Basic info fields
	Customer string `json:"customer"`
	Date     string `json:"date"`
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
	Notes    string `json:"notes"`

	// Payment fields
	PaymentMethod  string `json:"paymentMethod"`
	AmountPaid     string `json:"amountPaid"`
	CardDetails    string `json:"cardDetails"`
	PaymentDate    string `json:"paymentDate"`
	ReceivedBy     string `json:"receivedBy"`
	PaymentInfo    string `json:"paymentInfo"`

	// Audit trail
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditAction          string `json:"auditAction"`
	AuditUser            string `json:"auditUser"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`

	// Totals
	TotalGrossProfit string `json:"totalGrossProfit"`

	// Payment empty/table
	Reference            string `json:"reference"`
	PaymentEmptyTitle    string `json:"paymentEmptyTitle"`
	PaymentEmptyMessage  string `json:"paymentEmptyMessage"`

	// Line item management
	AddItem         string `json:"addItem"`
	AddDiscount     string `json:"addDiscount"`
	EditItem        string `json:"editItem"`
	RemoveItem      string `json:"removeItem"`
	ItemType        string `json:"itemType"`
	ItemTypeItem    string `json:"itemTypeItem"`
	ItemTypeDiscount string `json:"itemTypeDiscount"`
	InventoryItem   string `json:"inventoryItem"`
	SerialNumber    string `json:"serialNumber"`
	ItemEmptyTitle   string `json:"itemEmptyTitle"`
	ItemEmptyMessage string `json:"itemEmptyMessage"`
}

type RevenueConfirmLabels struct {
	Complete         string `json:"complete"`
	CompleteMessage  string `json:"completeMessage"`
	Reactivate       string `json:"reactivate"`
	ReactivateMessage string `json:"reactivateMessage"`
	BulkComplete        string `json:"bulkComplete"`
	BulkCompleteMessage string `json:"bulkCompleteMessage"`
	BulkReactivate       string `json:"bulkReactivate"`
	BulkReactivateMessage string `json:"bulkReactivateMessage"`
	SendEmail        string `json:"sendEmail"`
	SendEmailMessage string `json:"sendEmailMessage"`
}

type RevenueErrorLabels struct {
	PermissionDenied     string `json:"permissionDenied"`
	InvalidFormData      string `json:"invalidFormData"`
	NotFound             string `json:"notFound"`
	IDRequired           string `json:"idRequired"`
	NoIDsProvided        string `json:"noIDsProvided"`
	InvalidStatus        string `json:"invalidStatus"`
	InvalidTargetStatus  string `json:"invalidTargetStatus"`
	NoItemsCannotComplete string `json:"noItemsCannotComplete"`
	HasPaymentsCannotCancel string `json:"hasPaymentsCannotCancel"`
	BulkHasPayments      string `json:"bulkHasPayments"`
	BulkNoItems          string `json:"bulkNoItems"`
	PaymentNotFound      string `json:"paymentNotFound"`
	InvalidDiscount      string `json:"invalidDiscount"`
}

type RevenueDashboardLabels struct {
	Title          string `json:"title"`
	TotalSales     string `json:"totalSales"`
	Revenue        string `json:"revenue"`
	Completed      string `json:"completed"`
	Active         string `json:"active"`
	RevenueTrend   string `json:"revenueTrend"`
	Week           string `json:"week"`
	Month          string `json:"month"`
	Year           string `json:"year"`
	RecentSales    string `json:"recentSales"`
	ViewAll        string `json:"viewAll"`
	NewSaleCreated string `json:"newSaleCreated"`
	SaleCompleted  string `json:"saleCompleted"`
	SaleUpdated    string `json:"saleUpdated"`
	SaleCancelled  string `json:"saleCancelled"`
}

// RevenueSettingsLabels holds translatable strings for the revenue settings page
// (invoice template management).
type RevenueSettingsLabels struct {
	PageTitle      string `json:"pageTitle"`
	Caption        string `json:"caption"`
	UploadTemplate string `json:"uploadTemplate"`
	TemplateName   string `json:"templateName"`
	TemplateType   string `json:"templateType"`
	Purpose        string `json:"purpose"`
	SetDefault     string `json:"setDefault"`
	Delete         string `json:"delete"`
	DefaultBadge   string `json:"defaultBadge"`
	EmptyTitle     string `json:"emptyTitle"`
	EmptyMessage   string `json:"emptyMessage"`
	UploadSuccess  string `json:"uploadSuccess"`
	DeleteConfirm  string `json:"deleteConfirm"`
}

// ---------------------------------------------------------------------------
// Product labels
// ---------------------------------------------------------------------------

// ProductLabels holds all translatable strings for the product module.
type ProductLabels struct {
	Page    ProductPageLabels    `json:"page"`
	Buttons ProductButtonLabels  `json:"buttons"`
	Columns ProductColumnLabels  `json:"columns"`
	Empty   ProductEmptyLabels   `json:"empty"`
	Form    ProductFormLabels    `json:"form"`
	Actions ProductActionLabels  `json:"actions"`
	Bulk    ProductBulkLabels    `json:"bulkActions"`
	Tabs    ProductTabLabels     `json:"tabs"`
	Detail  ProductDetailLabels  `json:"detail"`
	Status  ProductStatusLabels  `json:"status"`
	Variant ProductVariantLabels `json:"variant"`
	Attribute ProductAttributeLabels `json:"attribute"`
	Options   ProductOptionLabels    `json:"options"`
	Confirm   ProductConfirmLabels   `json:"confirm"`
	Errors    ProductErrorLabels     `json:"errors"`
	Breadcrumb ProductBreadcrumbLabels `json:"breadcrumb"`
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
	Info       string `json:"info"`
	Variants   string `json:"variants"`
	Attributes string `json:"attributes"`
	Pricing    string `json:"pricing"`
	Options    string `json:"options"`
	Images     string `json:"images"`
	Stock      string `json:"stock"`
	AuditTrail string `json:"auditTrail"`
}

type ProductDetailLabels struct {
	Price        string `json:"price"`
	Currency     string `json:"currency"`
	Collections  string `json:"collections"`
	VariantCount string `json:"variantCount"`
	Status       string `json:"status"`
	OptionsLabel string `json:"optionsLabel"`
	EmptyVariantsMessage string `json:"emptyVariantsMessage"`
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
}

type ProductAttributeLabels struct {
	Title        string `json:"title"`
	DefaultValue string `json:"defaultValue"`
	Assign       string `json:"assign"`
	Remove       string `json:"remove"`
	Empty        string `json:"empty"`
}

type ProductConfirmLabels struct {
	Activate            string `json:"activate"`
	ActivateMessage     string `json:"activateMessage"`
	Deactivate          string `json:"deactivate"`
	DeactivateMessage   string `json:"deactivateMessage"`
	BulkActivate        string `json:"bulkActivate"`
	BulkActivateMessage string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
	BulkDelete          string `json:"bulkDelete"`
	BulkDeleteMessage   string `json:"bulkDeleteMessage"`
	RemoveVariant       string `json:"removeVariant"`
	RemoveVariantMessage string `json:"removeVariantMessage"`
}

type ProductErrorLabels struct {
	PermissionDenied   string `json:"permissionDenied"`
	InvalidFormData    string `json:"invalidFormData"`
	NotFound           string `json:"notFound"`
	IDRequired         string `json:"idRequired"`
	NoIDsProvided      string `json:"noIDsProvided"`
	InvalidStatus      string `json:"invalidStatus"`
	CannotDelete       string `json:"cannotDelete"`
	NameRequired       string `json:"nameRequired"`
	FieldRequired      string `json:"fieldRequired"`
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
	Tab       ProductOptionTabLabels       `json:"tab"`
	Columns   ProductOptionColumnLabels    `json:"columns"`
	Form      ProductOptionFormLabels      `json:"form"`
	DataTypes ProductOptionDataTypeLabels  `json:"dataTypes"`
	Value     ProductOptionValueLabels     `json:"value"`
	Actions   ProductOptionActionLabels    `json:"actions"`
	Empty     ProductOptionEmptyLabels     `json:"empty"`
	Confirm   ProductOptionConfirmLabels   `json:"confirm"`
}

type ProductOptionTabLabels struct {
	Title string `json:"title"`
}

type ProductOptionColumnLabels struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	DataType    string `json:"dataType"`
	ValuesCount string `json:"valuesCount"`
	SortOrder   string `json:"sortOrder"`
	Status      string `json:"status"`
}

type ProductOptionFormLabels struct {
	Name            string `json:"name"`
	NamePlaceholder string `json:"namePlaceholder"`
	Code            string `json:"code"`
	CodePlaceholder string `json:"codePlaceholder"`
	DataType        string `json:"dataType"`
	SortOrder       string `json:"sortOrder"`
	MinValue        string `json:"minValue"`
	MaxValue        string `json:"maxValue"`
	Active          string `json:"active"`
}

type ProductOptionDataTypeLabels struct {
	TextList   string `json:"textList"`
	NumberList string `json:"numberList"`
	ColorList  string `json:"colorList"`
	EnumList   string `json:"enumList"`
	FreeText   string `json:"freeText"`
	FreeNumber string `json:"freeNumber"`
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
}

type ProductOptionActionLabels struct {
	AddOption    string `json:"addOption"`
	EditOption   string `json:"editOption"`
	DeleteOption string `json:"deleteOption"`
	ViewValues   string `json:"viewValues"`
	AddValue     string `json:"addValue"`
	EditValue    string `json:"editValue"`
	DeleteValue  string `json:"deleteValue"`
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
	BasicInfo    string `json:"basicInfo"`
	Prices       string `json:"prices"`
	ProductName  string `json:"productName"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	AddPrice     string `json:"addPrice"`
	RemoveLabel  string `json:"removeLabel"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
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
// Expenditure labels
// ---------------------------------------------------------------------------

// ExpenditureLabels holds all translatable strings for the expenditure module
// (purchase + expense views).
type ExpenditureLabels struct {
	Labels  ExpenditureLabelNames   `json:"labels"`
	Page    ExpenditurePageLabels   `json:"page"`
	Buttons ExpenditureButtonLabels `json:"buttons"`
	Columns ExpenditureColumnLabels `json:"columns"`
	Empty   ExpenditureEmptyLabels  `json:"empty"`
	Form    ExpenditureFormLabels   `json:"form"`
	Status  ExpenditureStatusLabels `json:"status"`
	Types   ExpenditureTypeLabels   `json:"types"`
	Actions ExpenditureActionLabels `json:"actions"`
	Bulk    ExpenditureBulkLabels   `json:"bulkActions"`
	Detail  ExpenditureDetailLabels `json:"detail"`
}

type ExpenditureLabelNames struct {
	Name          string `json:"name"`
	NamePlural    string `json:"namePlural"`
	Purchase      string `json:"purchase"`
	PurchasePlural string `json:"purchasePlural"`
	PurchaseOrder string `json:"purchaseOrder"`
	Expense       string `json:"expense"`
	ExpensePlural string `json:"expensePlural"`
}

type ExpenditurePageLabels struct {
	PurchaseHeading          string `json:"purchaseHeading"`
	PurchaseCaption          string `json:"purchaseCaption"`
	PurchaseHeadingDraft     string `json:"purchaseHeadingDraft"`
	PurchaseHeadingPending   string `json:"purchaseHeadingPending"`
	PurchaseHeadingApproved  string `json:"purchaseHeadingApproved"`
	PurchaseHeadingPaid      string `json:"purchaseHeadingPaid"`
	PurchaseHeadingCancelled string `json:"purchaseHeadingCancelled"`
	PurchaseHeadingOverdue   string `json:"purchaseHeadingOverdue"`
	ExpenseHeading           string `json:"expenseHeading"`
	ExpenseCaption           string `json:"expenseCaption"`
	ExpenseHeadingDraft      string `json:"expenseHeadingDraft"`
	ExpenseHeadingPending    string `json:"expenseHeadingPending"`
	ExpenseHeadingApproved   string `json:"expenseHeadingApproved"`
	ExpenseHeadingPaid       string `json:"expenseHeadingPaid"`
	ExpenseHeadingCancelled  string `json:"expenseHeadingCancelled"`
	ExpenseHeadingOverdue    string `json:"expenseHeadingOverdue"`
	DashboardPurchase        string `json:"dashboardPurchase"`
	DashboardExpense         string `json:"dashboardExpense"`
}

type ExpenditureButtonLabels struct {
	AddPurchase string `json:"addPurchase"`
	AddExpense  string `json:"addExpense"`
}

type ExpenditureColumnLabels struct {
	Reference string `json:"reference"`
	Vendor    string `json:"vendor"`
	Amount    string `json:"amount"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Category  string `json:"category"`
}

type ExpenditureEmptyLabels struct {
	PurchaseTitle            string `json:"purchaseTitle"`
	PurchaseMessage          string `json:"purchaseMessage"`
	PurchaseDraftTitle       string `json:"purchaseDraftTitle"`
	PurchaseDraftMessage     string `json:"purchaseDraftMessage"`
	PurchasePendingTitle     string `json:"purchasePendingTitle"`
	PurchasePendingMessage   string `json:"purchasePendingMessage"`
	PurchaseApprovedTitle    string `json:"purchaseApprovedTitle"`
	PurchaseApprovedMessage  string `json:"purchaseApprovedMessage"`
	PurchasePaidTitle        string `json:"purchasePaidTitle"`
	PurchasePaidMessage      string `json:"purchasePaidMessage"`
	PurchaseCancelledTitle   string `json:"purchaseCancelledTitle"`
	PurchaseCancelledMessage string `json:"purchaseCancelledMessage"`
	PurchaseOverdueTitle     string `json:"purchaseOverdueTitle"`
	PurchaseOverdueMessage   string `json:"purchaseOverdueMessage"`
	ExpenseTitle             string `json:"expenseTitle"`
	ExpenseMessage           string `json:"expenseMessage"`
	ExpenseDraftTitle        string `json:"expenseDraftTitle"`
	ExpenseDraftMessage      string `json:"expenseDraftMessage"`
	ExpensePendingTitle      string `json:"expensePendingTitle"`
	ExpensePendingMessage    string `json:"expensePendingMessage"`
	ExpenseApprovedTitle     string `json:"expenseApprovedTitle"`
	ExpenseApprovedMessage   string `json:"expenseApprovedMessage"`
	ExpensePaidTitle         string `json:"expensePaidTitle"`
	ExpensePaidMessage       string `json:"expensePaidMessage"`
	ExpenseCancelledTitle    string `json:"expenseCancelledTitle"`
	ExpenseCancelledMessage  string `json:"expenseCancelledMessage"`
	ExpenseOverdueTitle      string `json:"expenseOverdueTitle"`
	ExpenseOverdueMessage    string `json:"expenseOverdueMessage"`
}

type ExpenditureFormLabels struct {
	VendorName               string `json:"vendorName"`
	VendorNamePlaceholder    string `json:"vendorNamePlaceholder"`
	ExpenditureDate          string `json:"expenditureDate"`
	TotalAmount              string `json:"totalAmount"`
	Currency                 string `json:"currency"`
	Status                   string `json:"status"`
	ReferenceNumber          string `json:"referenceNumber"`
	ReferenceNumberPlaceholder string `json:"referenceNumberPlaceholder"`
	PaymentTerms             string `json:"paymentTerms"`
	DueDate                  string `json:"dueDate"`
	ApprovedBy               string `json:"approvedBy"`
	ExpenditureType          string `json:"expenditureType"`
	ExpenditureCategory      string `json:"expenditureCategory"`
	Notes                    string `json:"notes"`
	NotesPlaceholder         string `json:"notesPlaceholder"`
	SectionInfo              string `json:"sectionInfo"`
	SectionVendor            string `json:"sectionVendor"`
	SectionPayment           string `json:"sectionPayment"`
	SectionNotes             string `json:"sectionNotes"`
}

type ExpenditureStatusLabels struct {
	Draft     string `json:"draft"`
	Pending   string `json:"pending"`
	Approved  string `json:"approved"`
	Paid      string `json:"paid"`
	Cancelled string `json:"cancelled"`
	Overdue   string `json:"overdue"`
}

type ExpenditureTypeLabels struct {
	Purchase string `json:"purchase"`
	Expense  string `json:"expense"`
	Refund   string `json:"refund"`
	Payroll  string `json:"payroll"`
}

type ExpenditureActionLabels struct {
	Add            string `json:"add"`
	Edit           string `json:"edit"`
	Delete         string `json:"delete"`
	Approve        string `json:"approve"`
	Reject         string `json:"reject"`
	MarkPaid       string `json:"markPaid"`
	ViewPurchase   string `json:"viewPurchase"`
	EditPurchase   string `json:"editPurchase"`
	DeletePurchase string `json:"deletePurchase"`
	ViewExpense    string `json:"viewExpense"`
	EditExpense    string `json:"editExpense"`
	DeleteExpense  string `json:"deleteExpense"`
}

type ExpenditureBulkLabels struct {
	Delete   string `json:"delete"`
	Approve  string `json:"approve"`
	MarkPaid string `json:"markPaid"`
}

type ExpenditureDetailLabels struct {
	PurchasePageTitle string `json:"purchasePageTitle"`
	ExpensePageTitle  string `json:"expensePageTitle"`
	VendorInfo        string `json:"vendorInfo"`
	VendorName        string `json:"vendorName"`
	Date              string `json:"date"`
	Amount            string `json:"amount"`
	Currency          string `json:"currency"`
	Status            string `json:"status"`
	Type              string `json:"type"`
	Category          string `json:"category"`
	ReferenceNumber   string `json:"referenceNumber"`
	PaymentTerms      string `json:"paymentTerms"`
	DueDate           string `json:"dueDate"`
	ApprovedBy        string `json:"approvedBy"`
	Notes             string `json:"notes"`
	LineItems         string `json:"lineItems"`
	Description       string `json:"description"`
	Quantity          string `json:"quantity"`
	UnitPrice         string `json:"unitPrice"`
	Total             string `json:"total"`
	SubTotal          string `json:"subTotal"`
	GrandTotal        string `json:"grandTotal"`
	TabBasicInfo      string `json:"tabBasicInfo"`
	TabLineItems      string `json:"tabLineItems"`
	TabPayment        string `json:"tabPayment"`
	TabAuditTrail     string `json:"tabAuditTrail"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditAction       string `json:"auditAction"`
	AuditUser         string `json:"auditUser"`
	AuditEmptyTitle   string `json:"auditEmptyTitle"`
	AuditEmptyMessage string `json:"auditEmptyMessage"`
}

// ---------------------------------------------------------------------------
// Collection labels (money IN — payment collections, receivables)
// ---------------------------------------------------------------------------

// CollectionLabels holds all translatable strings for the collection module.
type CollectionLabels struct {
	Page    CollectionPageLabels    `json:"page"`
	Buttons CollectionButtonLabels  `json:"buttons"`
	Columns CollectionColumnLabels  `json:"columns"`
	Empty   CollectionEmptyLabels   `json:"empty"`
	Form    CollectionFormLabels    `json:"form"`
	Actions CollectionActionLabels  `json:"actions"`
	Bulk    CollectionBulkLabels    `json:"bulkActions"`
	Detail  CollectionDetailLabels  `json:"detail"`
	Status  CollectionStatusLabels  `json:"status"`
	Confirm CollectionConfirmLabels `json:"confirm"`
	Errors  CollectionErrorLabels   `json:"errors"`
}

type CollectionPageLabels struct {
	Heading          string `json:"heading"`
	HeadingPending   string `json:"headingPending"`
	HeadingCompleted string `json:"headingCompleted"`
	HeadingFailed    string `json:"headingFailed"`
	Caption          string `json:"caption"`
	CaptionPending   string `json:"captionPending"`
	CaptionCompleted string `json:"captionCompleted"`
	CaptionFailed    string `json:"captionFailed"`
	Dashboard        string `json:"dashboard"`
}

type CollectionButtonLabels struct {
	AddCollection string `json:"addCollection"`
}

type CollectionColumnLabels struct {
	Reference string `json:"reference"`
	Customer  string `json:"customer"`
	Amount    string `json:"amount"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Method    string `json:"method"`
}

type CollectionEmptyLabels struct {
	PendingTitle     string `json:"pendingTitle"`
	PendingMessage   string `json:"pendingMessage"`
	CompletedTitle   string `json:"completedTitle"`
	CompletedMessage string `json:"completedMessage"`
	FailedTitle      string `json:"failedTitle"`
	FailedMessage    string `json:"failedMessage"`
}

type CollectionFormLabels struct {
	Customer             string `json:"customer"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Reference            string `json:"reference"`
	ReferencePlaceholder string `json:"referencePlaceholder"`
	PaymentMethod        string `json:"paymentMethod"`
	Status               string `json:"status"`
	Notes                string `json:"notes"`
	NotesPlaceholder     string `json:"notesPlaceholder"`
}

type CollectionActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	MarkComplete string `json:"markComplete"`
	Reactivate   string `json:"reactivate"`
}

type CollectionBulkLabels struct {
	Delete string `json:"delete"`
}

type CollectionDetailLabels struct {
	PageTitle    string `json:"pageTitle"`
	TitlePrefix  string `json:"titlePrefix"`
	PaymentInfo  string `json:"paymentInfo"`
	Customer     string `json:"customer"`
	Date         string `json:"date"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Status       string `json:"status"`
	Method       string `json:"method"`
	Reference    string `json:"reference"`
	Notes        string `json:"notes"`
	TabBasicInfo string `json:"tabBasicInfo"`
	TabAuditTrail string `json:"tabAuditTrail"`
	AuditAction       string `json:"auditAction"`
	AuditUser         string `json:"auditUser"`
	AuditEmptyTitle   string `json:"auditEmptyTitle"`
	AuditEmptyMessage string `json:"auditEmptyMessage"`
}

type CollectionStatusLabels struct {
	Pending   string `json:"pending"`
	Completed string `json:"completed"`
	Failed    string `json:"failed"`
}

type CollectionConfirmLabels struct {
	MarkComplete        string `json:"markComplete"`
	MarkCompleteMessage string `json:"markCompleteMessage"`
	Reactivate          string `json:"reactivate"`
	ReactivateMessage   string `json:"reactivateMessage"`
	Delete              string `json:"delete"`
	DeleteMessage       string `json:"deleteMessage"`
	BulkComplete        string `json:"bulkComplete"`
	BulkCompleteMessage string `json:"bulkCompleteMessage"`
	BulkReactivate       string `json:"bulkReactivate"`
	BulkReactivateMessage string `json:"bulkReactivateMessage"`
	BulkDelete           string `json:"bulkDelete"`
	BulkDeleteMessage    string `json:"bulkDeleteMessage"`
}

type CollectionErrorLabels struct {
	PermissionDenied  string `json:"permissionDenied"`
	InvalidFormData   string `json:"invalidFormData"`
	NotFound          string `json:"notFound"`
	IDRequired        string `json:"idRequired"`
	NoIDsProvided     string `json:"noIDsProvided"`
	InvalidStatus     string `json:"invalidStatus"`
}

// DefaultCollectionLabels returns CollectionLabels with sensible English defaults.
func DefaultCollectionLabels() CollectionLabels {
	return CollectionLabels{
		Page: CollectionPageLabels{
			Heading:          "Collections",
			HeadingPending:   "Pending Collections",
			HeadingCompleted: "Completed Collections",
			HeadingFailed:    "Failed Collections",
			Caption:          "Manage payment collections",
			CaptionPending:   "Payments awaiting collection",
			CaptionCompleted: "Successfully collected payments",
			CaptionFailed:    "Failed payment attempts",
			Dashboard:        "Collections Dashboard",
		},
		Buttons: CollectionButtonLabels{
			AddCollection: "Add Collection",
		},
		Columns: CollectionColumnLabels{
			Reference: "Reference",
			Customer:  "Customer",
			Amount:    "Amount",
			Date:      "Date",
			Status:    "Status",
			Method:    "Method",
		},
		Empty: CollectionEmptyLabels{
			PendingTitle:     "No pending collections",
			PendingMessage:   "No pending collections to display.",
			CompletedTitle:   "No completed collections",
			CompletedMessage: "No completed collections to display.",
			FailedTitle:      "No failed collections",
			FailedMessage:    "No failed collections to display.",
		},
		Form: CollectionFormLabels{
			Customer:             "Customer",
			Date:                 "Date",
			Amount:               "Amount",
			Currency:             "Currency",
			Reference:            "Reference",
			ReferencePlaceholder: "e.g. INV-001",
			PaymentMethod:        "Payment Method",
			Status:               "Status",
			Notes:                "Notes",
			NotesPlaceholder:     "Additional notes...",
		},
		Actions: CollectionActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			MarkComplete: "Mark Complete",
			Reactivate:   "Reactivate",
		},
		Bulk: CollectionBulkLabels{
			Delete: "Delete Selected",
		},
		Detail: CollectionDetailLabels{
			PageTitle:         "Collection Details",
			TitlePrefix:       "Collection #",
			PaymentInfo:       "Payment Information",
			Customer:          "Customer",
			Date:              "Date",
			Amount:            "Amount",
			Currency:          "Currency",
			Status:            "Status",
			Method:            "Payment Method",
			Reference:         "Reference",
			Notes:             "Notes",
			TabBasicInfo:      "Basic Info",
			TabAuditTrail:     "Audit Trail",
			AuditAction:       "Action",
			AuditUser:         "User",
			AuditEmptyTitle:   "No audit records",
			AuditEmptyMessage: "No audit trail entries yet.",
		},
		Status: CollectionStatusLabels{
			Pending:   "Pending",
			Completed: "Completed",
			Failed:    "Failed",
		},
		Confirm: CollectionConfirmLabels{
			MarkComplete:          "Mark Complete",
			MarkCompleteMessage:   "Are you sure you want to mark %s as complete?",
			Reactivate:            "Reactivate",
			ReactivateMessage:     "Are you sure you want to reactivate %s?",
			Delete:                "Delete",
			DeleteMessage:         "Are you sure you want to delete %s?",
			BulkComplete:          "Mark Complete",
			BulkCompleteMessage:   "Are you sure you want to mark {{count}} collection(s) as complete?",
			BulkReactivate:        "Reactivate",
			BulkReactivateMessage: "Are you sure you want to reactivate {{count}} collection(s)?",
			BulkDelete:            "Delete Collections",
			BulkDeleteMessage:     "Are you sure you want to delete {{count}} collection(s)?",
		},
		Errors: CollectionErrorLabels{
			PermissionDenied: "Permission denied",
			InvalidFormData:  "Invalid form data",
			NotFound:         "Collection not found",
			IDRequired:       "Collection ID is required",
			NoIDsProvided:    "No collection IDs provided",
			InvalidStatus:    "Invalid status",
		},
	}
}

// ---------------------------------------------------------------------------
// Disbursement labels (money OUT — payments, refunds, payouts)
// ---------------------------------------------------------------------------

// DisbursementLabels holds all translatable strings for the disbursement module.
type DisbursementLabels struct {
	Page    DisbursementPageLabels    `json:"page"`
	Buttons DisbursementButtonLabels  `json:"buttons"`
	Columns DisbursementColumnLabels  `json:"columns"`
	Empty   DisbursementEmptyLabels   `json:"empty"`
	Form    DisbursementFormLabels    `json:"form"`
	Actions DisbursementActionLabels  `json:"actions"`
	Bulk    DisbursementBulkLabels    `json:"bulkActions"`
	Detail  DisbursementDetailLabels  `json:"detail"`
	Status  DisbursementStatusLabels  `json:"status"`
	Confirm DisbursementConfirmLabels `json:"confirm"`
	Errors  DisbursementErrorLabels   `json:"errors"`
}

type DisbursementPageLabels struct {
	Heading          string `json:"heading"`
	HeadingDraft     string `json:"headingDraft"`
	HeadingPending   string `json:"headingPending"`
	HeadingApproved  string `json:"headingApproved"`
	HeadingPaid      string `json:"headingPaid"`
	HeadingCancelled string `json:"headingCancelled"`
	Caption          string `json:"caption"`
	CaptionDraft     string `json:"captionDraft"`
	CaptionPending   string `json:"captionPending"`
	CaptionApproved  string `json:"captionApproved"`
	CaptionPaid      string `json:"captionPaid"`
	CaptionCancelled string `json:"captionCancelled"`
	Dashboard        string `json:"dashboard"`
}

type DisbursementButtonLabels struct {
	AddDisbursement string `json:"addDisbursement"`
}

type DisbursementColumnLabels struct {
	Reference string `json:"reference"`
	Payee     string `json:"payee"`
	Amount    string `json:"amount"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Method    string `json:"method"`
	Category  string `json:"category"`
}

type DisbursementEmptyLabels struct {
	DraftTitle       string `json:"draftTitle"`
	DraftMessage     string `json:"draftMessage"`
	PendingTitle     string `json:"pendingTitle"`
	PendingMessage   string `json:"pendingMessage"`
	ApprovedTitle    string `json:"approvedTitle"`
	ApprovedMessage  string `json:"approvedMessage"`
	PaidTitle        string `json:"paidTitle"`
	PaidMessage      string `json:"paidMessage"`
	CancelledTitle   string `json:"cancelledTitle"`
	CancelledMessage string `json:"cancelledMessage"`
}

type DisbursementFormLabels struct {
	Payee                string `json:"payee"`
	PayeePlaceholder     string `json:"payeePlaceholder"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Reference            string `json:"reference"`
	ReferencePlaceholder string `json:"referencePlaceholder"`
	PaymentMethod        string `json:"paymentMethod"`
	Category             string `json:"category"`
	Status               string `json:"status"`
	Notes                string `json:"notes"`
	NotesPlaceholder     string `json:"notesPlaceholder"`
	ApprovedBy           string `json:"approvedBy"`
}

type DisbursementActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Approve    string `json:"approve"`
	MarkPaid   string `json:"markPaid"`
	Cancel     string `json:"cancel"`
	Submit     string `json:"submit"`
	Reactivate string `json:"reactivate"`
}

type DisbursementBulkLabels struct {
	Delete   string `json:"delete"`
	Approve  string `json:"approve"`
	MarkPaid string `json:"markPaid"`
}

type DisbursementDetailLabels struct {
	PageTitle         string `json:"pageTitle"`
	TitlePrefix       string `json:"titlePrefix"`
	PaymentInfo       string `json:"paymentInfo"`
	Payee             string `json:"payee"`
	Date              string `json:"date"`
	Amount            string `json:"amount"`
	Currency          string `json:"currency"`
	Status            string `json:"status"`
	Method            string `json:"method"`
	Category          string `json:"category"`
	Reference         string `json:"reference"`
	ApprovedBy        string `json:"approvedBy"`
	Notes             string `json:"notes"`
	TabBasicInfo      string `json:"tabBasicInfo"`
	TabAuditTrail     string `json:"tabAuditTrail"`
	AuditAction       string `json:"auditAction"`
	AuditUser         string `json:"auditUser"`
	AuditEmptyTitle   string `json:"auditEmptyTitle"`
	AuditEmptyMessage string `json:"auditEmptyMessage"`
}

type DisbursementStatusLabels struct {
	Draft     string `json:"draft"`
	Pending   string `json:"pending"`
	Approved  string `json:"approved"`
	Paid      string `json:"paid"`
	Cancelled string `json:"cancelled"`
}

type DisbursementConfirmLabels struct {
	Submit              string `json:"submit"`
	SubmitMessage       string `json:"submitMessage"`
	Approve             string `json:"approve"`
	ApproveMessage      string `json:"approveMessage"`
	MarkPaid            string `json:"markPaid"`
	MarkPaidMessage     string `json:"markPaidMessage"`
	Cancel              string `json:"cancel"`
	CancelMessage       string `json:"cancelMessage"`
	Reactivate          string `json:"reactivate"`
	ReactivateMessage   string `json:"reactivateMessage"`
	Delete              string `json:"delete"`
	DeleteMessage       string `json:"deleteMessage"`
	BulkSubmit          string `json:"bulkSubmit"`
	BulkSubmitMessage   string `json:"bulkSubmitMessage"`
	BulkApprove         string `json:"bulkApprove"`
	BulkApproveMessage  string `json:"bulkApproveMessage"`
	BulkMarkPaid        string `json:"bulkMarkPaid"`
	BulkMarkPaidMessage string `json:"bulkMarkPaidMessage"`
	BulkCancel          string `json:"bulkCancel"`
	BulkCancelMessage   string `json:"bulkCancelMessage"`
	BulkReactivate        string `json:"bulkReactivate"`
	BulkReactivateMessage string `json:"bulkReactivateMessage"`
	BulkDelete           string `json:"bulkDelete"`
	BulkDeleteMessage    string `json:"bulkDeleteMessage"`
}

type DisbursementErrorLabels struct {
	PermissionDenied    string `json:"permissionDenied"`
	InvalidFormData     string `json:"invalidFormData"`
	NotFound            string `json:"notFound"`
	IDRequired          string `json:"idRequired"`
	NoIDsProvided       string `json:"noIDsProvided"`
	InvalidStatus       string `json:"invalidStatus"`
	InvalidTransition   string `json:"invalidTransition"`
}

// DefaultDisbursementLabels returns DisbursementLabels with sensible English defaults.
func DefaultDisbursementLabels() DisbursementLabels {
	return DisbursementLabels{
		Page: DisbursementPageLabels{
			Heading:          "Disbursements",
			HeadingDraft:     "Draft Disbursements",
			HeadingPending:   "Pending Disbursements",
			HeadingApproved:  "Approved Disbursements",
			HeadingPaid:      "Paid Disbursements",
			HeadingCancelled: "Cancelled Disbursements",
			Caption:          "Manage disbursements and payouts",
			CaptionDraft:     "Draft disbursements awaiting submission",
			CaptionPending:   "Disbursements awaiting approval",
			CaptionApproved:  "Approved disbursements ready for payment",
			CaptionPaid:      "Completed disbursement payments",
			CaptionCancelled: "Cancelled disbursements",
			Dashboard:        "Disbursements Dashboard",
		},
		Buttons: DisbursementButtonLabels{
			AddDisbursement: "Add Disbursement",
		},
		Columns: DisbursementColumnLabels{
			Reference: "Reference",
			Payee:     "Payee",
			Amount:    "Amount",
			Date:      "Date",
			Status:    "Status",
			Method:    "Method",
			Category:  "Category",
		},
		Empty: DisbursementEmptyLabels{
			DraftTitle:       "No draft disbursements",
			DraftMessage:     "No draft disbursements to display.",
			PendingTitle:     "No pending disbursements",
			PendingMessage:   "No pending disbursements to display.",
			ApprovedTitle:    "No approved disbursements",
			ApprovedMessage:  "No approved disbursements to display.",
			PaidTitle:        "No paid disbursements",
			PaidMessage:      "No paid disbursements to display.",
			CancelledTitle:   "No cancelled disbursements",
			CancelledMessage: "No cancelled disbursements to display.",
		},
		Form: DisbursementFormLabels{
			Payee:                "Payee",
			PayeePlaceholder:     "Enter payee name",
			Date:                 "Date",
			Amount:               "Amount",
			Currency:             "Currency",
			Reference:            "Reference",
			ReferencePlaceholder: "e.g. DISB-001",
			PaymentMethod:        "Payment Method",
			Category:             "Category",
			Status:               "Status",
			Notes:                "Notes",
			NotesPlaceholder:     "Additional notes...",
			ApprovedBy:           "Approved By",
		},
		Actions: DisbursementActionLabels{
			View:       "View",
			Edit:       "Edit",
			Delete:     "Delete",
			Approve:    "Approve",
			MarkPaid:   "Mark as Paid",
			Cancel:     "Cancel",
			Submit:     "Submit",
			Reactivate: "Reactivate",
		},
		Bulk: DisbursementBulkLabels{
			Delete:   "Delete Selected",
			Approve:  "Approve Selected",
			MarkPaid: "Mark Selected as Paid",
		},
		Detail: DisbursementDetailLabels{
			PageTitle:         "Disbursement Details",
			TitlePrefix:       "Disbursement #",
			PaymentInfo:       "Payment Information",
			Payee:             "Payee",
			Date:              "Date",
			Amount:            "Amount",
			Currency:          "Currency",
			Status:            "Status",
			Method:            "Payment Method",
			Category:          "Category",
			Reference:         "Reference",
			ApprovedBy:        "Approved By",
			Notes:             "Notes",
			TabBasicInfo:      "Basic Info",
			TabAuditTrail:     "Audit Trail",
			AuditAction:       "Action",
			AuditUser:         "User",
			AuditEmptyTitle:   "No audit records",
			AuditEmptyMessage: "No audit trail entries yet.",
		},
		Status: DisbursementStatusLabels{
			Draft:     "Draft",
			Pending:   "Pending",
			Approved:  "Approved",
			Paid:      "Paid",
			Cancelled: "Cancelled",
		},
		Confirm: DisbursementConfirmLabels{
			Submit:                "Submit",
			SubmitMessage:         "Are you sure you want to submit {{count}} disbursement(s)?",
			Approve:               "Approve",
			ApproveMessage:        "Are you sure you want to approve {{count}} disbursement(s)?",
			MarkPaid:              "Mark as Paid",
			MarkPaidMessage:       "Are you sure you want to mark {{count}} disbursement(s) as paid?",
			Cancel:                "Cancel",
			CancelMessage:         "Are you sure you want to cancel {{count}} disbursement(s)?",
			Reactivate:            "Reactivate",
			ReactivateMessage:     "Are you sure you want to reactivate {{count}} disbursement(s)?",
			Delete:                "Delete",
			DeleteMessage:         "Are you sure you want to delete {{count}} disbursement(s)?",
			BulkSubmit:            "Submit Disbursements",
			BulkSubmitMessage:     "Are you sure you want to submit {{count}} disbursement(s)?",
			BulkApprove:           "Approve Disbursements",
			BulkApproveMessage:    "Are you sure you want to approve {{count}} disbursement(s)?",
			BulkMarkPaid:          "Mark as Paid",
			BulkMarkPaidMessage:   "Are you sure you want to mark {{count}} disbursement(s) as paid?",
			BulkCancel:            "Cancel Disbursements",
			BulkCancelMessage:     "Are you sure you want to cancel {{count}} disbursement(s)?",
			BulkReactivate:        "Reactivate Disbursements",
			BulkReactivateMessage: "Are you sure you want to reactivate {{count}} disbursement(s)?",
			BulkDelete:            "Delete Disbursements",
			BulkDeleteMessage:     "Are you sure you want to delete {{count}} disbursement(s)?",
		},
		Errors: DisbursementErrorLabels{
			PermissionDenied:  "Permission denied",
			InvalidFormData:   "Invalid form data",
			NotFound:          "Disbursement not found",
			IDRequired:        "Disbursement ID is required",
			NoIDsProvided:     "No disbursement IDs provided",
			InvalidStatus:     "Invalid target status",
			InvalidTransition: "Cannot transition from %s to %s",
		},
	}
}

// ---------------------------------------------------------------------------
// Plan labels
// ---------------------------------------------------------------------------

// PlanLabels holds all translatable strings for the plan module.
type PlanLabels struct {
	Page    PlanPageLabels    `json:"page"`
	Buttons PlanButtonLabels  `json:"buttons"`
	Columns PlanColumnLabels  `json:"columns"`
	Empty   PlanEmptyLabels   `json:"empty"`
	Form    PlanFormLabels    `json:"form"`
	Actions PlanActionLabels  `json:"actions"`
	Detail  PlanDetailLabels  `json:"detail"`
	Tabs    PlanTabLabels     `json:"tabs"`
	Confirm PlanConfirmLabels `json:"confirm"`
	Errors  PlanErrorLabels   `json:"errors"`
}

type PlanPageLabels struct {
	Heading        string `json:"heading"`
	HeadingActive  string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption        string `json:"caption"`
	CaptionActive  string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type PlanButtonLabels struct {
	AddPlan string `json:"addPlan"`
}

type PlanColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Interval    string `json:"interval"`
	Price       string `json:"price"`
	Status      string `json:"status"`
}

type PlanEmptyLabels struct {
	Title           string `json:"title"`
	Message         string `json:"message"`
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type PlanActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type PlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
}

// ---------------------------------------------------------------------------
// Plan form, detail, tabs, confirm sub-labels
// ---------------------------------------------------------------------------

type PlanFormSectionLabels struct {
	Basic    string `json:"basic"`
	Services string `json:"services"`
}

type PlanFormLabels struct {
	Name                string                `json:"name"`
	NamePlaceholder     string                `json:"namePlaceholder"`
	Description         string                `json:"description"`
	DescPlaceholder     string                `json:"descriptionPlaceholder"`
	FulfillmentType     string                `json:"fulfillmentType"`
	Active              string                `json:"active"`
	Products            string                `json:"products"`
	ProductsPlaceholder string                `json:"productsPlaceholder"`
	ProductsSearch      string                `json:"productsSearch"`
	Sections            PlanFormSectionLabels  `json:"sections"`
}

type PlanDetailLabels struct {
	PageTitle       string `json:"pageTitle"`
	Price           string `json:"price"`
	Currency        string `json:"currency"`
	Status          string `json:"status"`
	Description     string `json:"description"`
	FulfillmentType string `json:"fulfillmentType"`
	CreatedDate     string `json:"createdDate"`
	ModifiedDate    string `json:"modifiedDate"`
}

type PlanTabLabels struct {
	Info       string `json:"info"`
	Products   string `json:"products"`
	PriceLists string `json:"priceLists"`
	AuditTrail string `json:"auditTrail"`
}

type PlanConfirmLabels struct {
	Delete            string `json:"delete"`
	DeleteMessage     string `json:"deleteMessage"`
	Deactivate        string `json:"deactivate"`
	DeactivateMessage string `json:"deactivateMessage"`
}

// ---------------------------------------------------------------------------
// Subscription labels
// ---------------------------------------------------------------------------

// SubscriptionLabels holds all translatable strings for the subscription module.
type SubscriptionLabels struct {
	Page    SubscriptionPageLabels    `json:"page"`
	Buttons SubscriptionButtonLabels  `json:"buttons"`
	Columns SubscriptionColumnLabels  `json:"columns"`
	Empty   SubscriptionEmptyLabels   `json:"empty"`
	Form    SubscriptionFormLabels    `json:"form"`
	Actions SubscriptionActionLabels  `json:"actions"`
	Detail  SubscriptionDetailLabels  `json:"detail"`
	Tabs    SubscriptionTabLabels     `json:"tabs"`
	Confirm SubscriptionConfirmLabels `json:"confirm"`
	Errors  SubscriptionErrorLabels   `json:"errors"`
}

type SubscriptionPageLabels struct {
	Heading        string `json:"heading"`
	HeadingActive  string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption        string `json:"caption"`
	CaptionActive  string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type SubscriptionButtonLabels struct {
	AddSubscription string `json:"addSubscription"`
}

type SubscriptionColumnLabels struct {
	Customer  string `json:"customer"`
	Plan      string `json:"plan"`
	StartDate string `json:"startDate"`
	Status    string `json:"status"`
}

type SubscriptionEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SubscriptionActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Cancel string `json:"cancel"`
}

type SubscriptionErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
}

// ---------------------------------------------------------------------------
// Subscription form, detail, tabs, confirm sub-labels
// ---------------------------------------------------------------------------

type SubscriptionFormLabels struct {
	Customer            string `json:"customer"`
	CustomerPlaceholder string `json:"customerPlaceholder"`
	Plan                string `json:"plan"`
	PlanPlaceholder     string `json:"planPlaceholder"`
	StartDate           string `json:"startDate"`
	EndDate             string `json:"endDate"`
	Active              string `json:"active"`
	Notes                     string `json:"notes"`
	NotesPlaceholder          string `json:"notesPlaceholder"`
	CustomerSearchPlaceholder string `json:"customerSearchPlaceholder"`
	PlanSearchPlaceholder     string `json:"planSearchPlaceholder"`
	CustomerNoResults         string `json:"customerNoResults"`
	PlanNoResults             string `json:"planNoResults"`
}

type SubscriptionDetailLabels struct {
	PageTitle    string `json:"pageTitle"`
	Customer     string `json:"customer"`
	Plan         string `json:"plan"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	Status       string `json:"status"`
	CreatedDate  string `json:"createdDate"`
	ModifiedDate string `json:"modifiedDate"`
}

type SubscriptionTabLabels struct {
	Info       string `json:"info"`
	History    string `json:"history"`
	AuditTrail string `json:"auditTrail"`
}

type SubscriptionConfirmLabels struct {
	Cancel        string `json:"cancel"`
	CancelMessage string `json:"cancelMessage"`
}

// DefaultPlanLabels returns PlanLabels with sensible English defaults.
func DefaultPlanLabels() PlanLabels {
	return PlanLabels{
		Page: PlanPageLabels{
			Heading:         "Plans",
			HeadingActive:   "Active Plans",
			HeadingInactive: "Inactive Plans",
			Caption:         "Manage your plans",
			CaptionActive:   "Manage your active plans",
			CaptionInactive: "View inactive or archived plans",
		},
		Buttons: PlanButtonLabels{
			AddPlan: "Add Plan",
		},
		Columns: PlanColumnLabels{
			Name:        "Name",
			Description: "Description",
			Interval:    "Interval",
			Price:       "Price",
			Status:      "Status",
		},
		Empty: PlanEmptyLabels{
			Title:           "No plans found",
			Message:         "No plans to display.",
			ActiveTitle:     "No active plans",
			ActiveMessage:   "Create your first plan to get started.",
			InactiveTitle:   "No inactive plans",
			InactiveMessage: "Discontinued plans will appear here.",
		},
		Form: PlanFormLabels{
			Name:                "Plan Name",
			NamePlaceholder:     "Enter plan name",
			Description:         "Description",
			DescPlaceholder:     "Enter plan description...",
			FulfillmentType:     "Fulfillment Type",
			Active:              "Active",
			Products:            "Products",
			ProductsPlaceholder: "Select products...",
			ProductsSearch:      "Search products...",
			Sections: PlanFormSectionLabels{
				Basic:    "Basic Information",
				Services: "Assigned Products",
			},
		},
		Actions: PlanActionLabels{
			View:   "View Plan",
			Edit:   "Edit Plan",
			Delete: "Delete Plan",
		},
		Detail: PlanDetailLabels{
			PageTitle:       "Plan Details",
			Price:           "Price",
			Currency:        "Currency",
			Status:          "Status",
			Description:     "Description",
			FulfillmentType: "Fulfillment Type",
			CreatedDate:     "Created",
			ModifiedDate:    "Last Modified",
		},
		Tabs: PlanTabLabels{
			Info:       "Information",
			Products:   "Products",
			PriceLists: "Price Lists",
			AuditTrail: "Audit Trail",
		},
		Confirm: PlanConfirmLabels{
			Delete:            "Delete Plan",
			DeleteMessage:     "Are you sure you want to delete \"%s\"? This action cannot be undone.",
			Deactivate:        "Deactivate Plan",
			DeactivateMessage: "Are you sure you want to deactivate \"%s\"?",
		},
		Errors: PlanErrorLabels{
			PermissionDenied: "You do not have permission to perform this action",
			InvalidFormData:  "Invalid form data. Please check your inputs and try again.",
			NotFound:         "Plan not found",
			IDRequired:       "Plan ID is required",
			NoPermission:     "No permission",
		},
	}
}

// DefaultSubscriptionLabels returns SubscriptionLabels with sensible English defaults.
func DefaultSubscriptionLabels() SubscriptionLabels {
	return SubscriptionLabels{
		Page: SubscriptionPageLabels{
			Heading:         "Subscriptions",
			HeadingActive:   "Active Subscriptions",
			HeadingInactive: "Inactive Subscriptions",
			Caption:         "Subscription management",
			CaptionActive:   "Manage your active subscriptions",
			CaptionInactive: "View cancelled or expired subscriptions",
		},
		Buttons: SubscriptionButtonLabels{
			AddSubscription: "Add Subscription",
		},
		Columns: SubscriptionColumnLabels{
			Customer:  "Customer",
			Plan:      "Plan",
			StartDate: "Start Date",
			Status:    "Status",
		},
		Empty: SubscriptionEmptyLabels{
			Title:   "No subscriptions found",
			Message: "No subscriptions to display.",
		},
		Form: SubscriptionFormLabels{
			Customer:            "Customer",
			CustomerPlaceholder: "Select customer...",
			Plan:                "Plan",
			PlanPlaceholder:     "Select plan...",
			StartDate:           "Start Date",
			EndDate:             "End Date",
			Active:              "Active",
			Notes:               "Notes",
			NotesPlaceholder:          "Enter notes...",
			CustomerSearchPlaceholder: "Search customers...",
			PlanSearchPlaceholder:     "Search plans...",
			CustomerNoResults:         "No customers found",
			PlanNoResults:             "No plans found",
		},
		Actions: SubscriptionActionLabels{
			View:   "View Subscription",
			Edit:   "Edit Subscription",
			Cancel: "Cancel Subscription",
		},
		Detail: SubscriptionDetailLabels{
			PageTitle:    "Subscription Details",
			Customer:     "Customer",
			Plan:         "Plan",
			StartDate:    "Start Date",
			EndDate:      "End Date",
			Status:       "Status",
			CreatedDate:  "Created",
			ModifiedDate: "Last Modified",
		},
		Tabs: SubscriptionTabLabels{
			Info:       "Information",
			History:    "History",
			AuditTrail: "Audit Trail",
		},
		Confirm: SubscriptionConfirmLabels{
			Cancel:        "Cancel Subscription",
			CancelMessage: "Are you sure you want to cancel this subscription? This action cannot be undone.",
		},
		Errors: SubscriptionErrorLabels{
			PermissionDenied: "You do not have permission to perform this action",
			InvalidFormData:  "Invalid form data. Please check your inputs and try again.",
			NotFound:         "Subscription not found",
			IDRequired:       "Subscription ID is required",
			NoPermission:     "No permission",
		},
	}
}

// ---------------------------------------------------------------------------
// Mapping helpers
// ---------------------------------------------------------------------------

// MapTableLabels maps common labels into the flat types.TableLabels structure.
func MapTableLabels(common pyeza.CommonLabels) types.TableLabels {
	return types.TableLabels{
		Search:             common.Table.Search,
		SearchPlaceholder:  common.Table.SearchPlaceholder,
		Filters:            common.Table.Filters,
		FilterConditions:   common.Table.FilterConditions,
		ClearAll:           common.Table.ClearAll,
		AddCondition:       common.Table.AddCondition,
		Clear:              common.Table.Clear,
		ApplyFilters:       common.Table.ApplyFilters,
		Sort:               common.Table.Sort,
		Columns:            common.Table.Columns,
		Export:              common.Table.Export,
		DensityDefault:     common.Table.Density.Default,
		DensityComfortable: common.Table.Density.Comfortable,
		DensityCompact:     common.Table.Density.Compact,
		Show:               common.Table.Show,
		Entries:             common.Table.Entries,
		Showing:            common.Table.Showing,
		To:                 common.Table.To,
		Of:                 common.Table.Of,
		EntriesLabel:       common.Table.EntriesLabel,
		SelectAll:          common.Table.SelectAll,
		Actions:            common.Table.Actions,
		Prev:               common.Pagination.Prev,
		Next:               common.Pagination.Next,
	}
}

// MapBulkConfig returns a BulkActionsConfig with labels from common bulk labels.
func MapBulkConfig(common pyeza.CommonLabels) types.BulkActionsConfig {
	return types.BulkActionsConfig{
		Enabled:        true,
		SelectAllLabel: common.Bulk.SelectAll,
		SelectedLabel:  common.Bulk.Selected,
		CancelLabel:    common.Bulk.ClearSelection,
	}
}

// LocationMap maps location slugs to display names.
var LocationMap = map[string]string{
	"ayala-central-bloc": "Ayala Central Bloc",
	"sm-city-cebu":       "SM City Cebu",
	"ayala-center-cebu":  "Ayala Center Cebu",
	"robinsons-galleria": "Robinsons Galleria",
}

// LocationDisplayName returns the display name for a location slug.
func LocationDisplayName(slug string) string {
	if name, ok := LocationMap[slug]; ok {
		return name
	}
	return slug
}
