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
	Product          string `json:"product"`
	SKU              string `json:"sku"`
	SKUPlaceholder   string `json:"skuPlaceholder"`
	OnHand           string `json:"onHand"`
	Reserved         string `json:"reserved"`
	ReorderLevel     string `json:"reorderLevel"`
	UnitOfMeasure    string `json:"unitOfMeasure"`
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	Active           string `json:"active"`
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
	RevenueReference string `json:"revenueReference"`

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
	Attachments  string `json:"attachments"`
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
	StockLevels          string `json:"stockLevels"`
	RecentActivity       string `json:"recentActivity"`
	ViewAll              string `json:"viewAll"`
	Week                 string `json:"week"`
	Month                string `json:"month"`
	Year                 string `json:"year"`
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
	PermissionDenied          string `json:"permissionDenied"`
	InvalidFormData           string `json:"invalidFormData"`
	NotFound                  string `json:"notFound"`
	IDRequired                string `json:"idRequired"`
	NoIDsProvided             string `json:"noIDsProvided"`
	InvalidStatus             string `json:"invalidStatus"`
	NoPermission              string `json:"noPermission"`
	SerialNotFound            string `json:"serialNotFound"`
	SerialIDRequired          string `json:"serialIDRequired"`
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
	Page      RevenuePageLabels      `json:"page"`
	Buttons   RevenueButtonLabels    `json:"buttons"`
	Columns   RevenueColumnLabels    `json:"columns"`
	Empty     RevenueEmptyLabels     `json:"empty"`
	Form      RevenueFormLabels      `json:"form"`
	Actions   RevenueActionLabels    `json:"actions"`
	Bulk      RevenueBulkLabels      `json:"bulkActions"`
	Detail    RevenueDetailLabels    `json:"detail"`
	Confirm   RevenueConfirmLabels   `json:"confirm"`
	Errors    RevenueErrorLabels     `json:"errors"`
	Dashboard RevenueDashboardLabels `json:"dashboard"`
	Settings  RevenueSettingsLabels  `json:"settings"`
}

type RevenuePageLabels struct {
	Heading          string `json:"heading"`
	HeadingDraft     string `json:"headingDraft"`
	HeadingComplete  string `json:"headingComplete"`
	HeadingCancelled string `json:"headingCancelled"`
	Caption          string `json:"caption"`
	CaptionDraft     string `json:"captionDraft"`
	CaptionComplete  string `json:"captionComplete"`
	CaptionCancelled string `json:"captionCancelled"`
}

type RevenueButtonLabels struct {
	AddSale string `json:"addSale"`
}

type RevenueColumnLabels struct {
	Reference string `json:"reference"`
	Customer  string `json:"customer"`
	Date      string `json:"date"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
}

type RevenueEmptyLabels struct {
	DraftTitle       string `json:"draftTitle"`
	DraftMessage     string `json:"draftMessage"`
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

	// Payment terms and client search labels
	PaymentTerms              string `json:"paymentTerms"`
	SelectPaymentTerm         string `json:"selectPaymentTerm"`
	DueDate                   string `json:"dueDate"`
	CustomerSearchPlaceholder string `json:"customerSearchPlaceholder"`
	CustomerNoResults         string `json:"customerNoResults"`

	// Subscription search labels
	Subscription          string `json:"subscription"`
	SubscriptionNoResults string `json:"subscriptionNoResults"`

	// Placeholders and translated option labels
	CurrencyPlaceholder            string `json:"currencyPlaceholder"`
	CustomerNamePlaceholder        string `json:"customerNamePlaceholder"`
	StatusDraft                    string `json:"statusDraft"`
	StatusComplete                 string `json:"statusComplete"`
	StatusCancelled                string `json:"statusCancelled"`
	PaymentMethod                  string `json:"paymentMethod"`
	ReferenceNumber                string `json:"referenceNumber"`
	TransactionIdPlaceholder       string `json:"transactionIdPlaceholder"`
	ReceivedBy                     string `json:"receivedBy"`
	Role                           string `json:"role"`
	SelectInventoryItem            string `json:"selectInventoryItem"`
	ItemDescriptionPlaceholder     string `json:"itemDescriptionPlaceholder"`
	DiscountDescriptionPlaceholder string `json:"discountDescriptionPlaceholder"`
}

type RevenueActionLabels struct {
	View              string `json:"view"`
	Edit              string `json:"edit"`
	Delete            string `json:"delete"`
	Complete          string `json:"complete"`
	Reactivate        string `json:"reactivate"`
	DownloadInvoice   string `json:"downloadInvoice"`
	SendEmail         string `json:"sendEmail"`
	Cancel            string `json:"cancel"`
	ReclassifyToDraft string `json:"reclassifyToDraft"`
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
	TabBasicInfo   string `json:"tabBasicInfo"`
	TabLineItems   string `json:"tabLineItems"`
	TabPayment     string `json:"tabPayment"`
	TabAttachments string `json:"tabAttachments"`
	TabAuditTrail  string `json:"tabAuditTrail"`

	// Basic info fields
	Customer     string `json:"customer"`
	Date         string `json:"date"`
	Amount       string `json:"amount"`
	Currency     string `json:"currency"`
	Status       string `json:"status"`
	Notes        string `json:"notes"`
	PaymentTerms string `json:"paymentTerms"`
	DueDate      string `json:"dueDate"`

	// Payment fields
	PaymentMethod string `json:"paymentMethod"`
	AmountPaid    string `json:"amountPaid"`
	CardDetails   string `json:"cardDetails"`
	PaymentDate   string `json:"paymentDate"`
	ReceivedBy    string `json:"receivedBy"`
	PaymentInfo   string `json:"paymentInfo"`

	// Audit trail
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditAction          string `json:"auditAction"`
	AuditUser            string `json:"auditUser"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`

	// Totals
	TotalGrossProfit string `json:"totalGrossProfit"`

	// Payment empty/table
	Reference           string `json:"reference"`
	PaymentEmptyTitle   string `json:"paymentEmptyTitle"`
	PaymentEmptyMessage string `json:"paymentEmptyMessage"`

	// Line item management
	AddItem                    string `json:"addItem"`
	AddDiscount                string `json:"addDiscount"`
	EditItem                   string `json:"editItem"`
	RemoveItem                 string `json:"removeItem"`
	ItemType                   string `json:"itemType"`
	ItemTypeItem               string `json:"itemTypeItem"`
	ItemTypeDiscount           string `json:"itemTypeDiscount"`
	InventoryItem              string `json:"inventoryItem"`
	SelectInventoryItem        string `json:"selectInventoryItem"`
	ItemDescriptionPlaceholder string `json:"itemDescriptionPlaceholder"`
	NotesPlaceholder           string `json:"notesPlaceholder"`
	SerialNumber               string `json:"serialNumber"`
	Product                    string `json:"product"`
	ProductNoResults           string `json:"productNoResults"`
	ProductPlaceholder         string `json:"productPlaceholder"`
	ItemEmptyTitle             string `json:"itemEmptyTitle"`
	ItemEmptyMessage           string `json:"itemEmptyMessage"`

	// Payment tab
	TotalPaid                  string `json:"totalPaid"`
	Remaining                  string `json:"remaining"`
	RecordPayment              string `json:"recordPayment"`
	NoPaymentInfo              string `json:"noPaymentInfo"`
	PaymentDetailsNotAvailable string `json:"paymentDetailsNotAvailable"`
}

type RevenueConfirmLabels struct {
	Complete                 string `json:"complete"`
	CompleteMessage          string `json:"completeMessage"`
	Reactivate               string `json:"reactivate"`
	ReactivateMessage        string `json:"reactivateMessage"`
	BulkComplete             string `json:"bulkComplete"`
	BulkCompleteMessage      string `json:"bulkCompleteMessage"`
	BulkReactivate           string `json:"bulkReactivate"`
	BulkReactivateMessage    string `json:"bulkReactivateMessage"`
	SendEmail                string `json:"sendEmail"`
	SendEmailMessage         string `json:"sendEmailMessage"`
	Cancel                   string `json:"cancel"`
	CancelMessage            string `json:"cancelMessage"`
	ReclassifyToDraft        string `json:"reclassifyToDraft"`
	ReclassifyToDraftMessage string `json:"reclassifyToDraftMessage"`
}

type RevenueErrorLabels struct {
	PermissionDenied        string `json:"permissionDenied"`
	InvalidFormData         string `json:"invalidFormData"`
	NotFound                string `json:"notFound"`
	IDRequired              string `json:"idRequired"`
	NoIDsProvided           string `json:"noIDsProvided"`
	InvalidStatus           string `json:"invalidStatus"`
	InvalidTargetStatus     string `json:"invalidTargetStatus"`
	NoItemsCannotComplete   string `json:"noItemsCannotComplete"`
	HasPaymentsCannotCancel string `json:"hasPaymentsCannotCancel"`
	BulkHasPayments         string `json:"bulkHasPayments"`
	BulkNoItems             string `json:"bulkNoItems"`
	PaymentNotFound         string `json:"paymentNotFound"`
	InvalidDiscount         string `json:"invalidDiscount"`
}

type RevenueDashboardLabels struct {
	Title          string `json:"title"`
	TotalRevenue   string `json:"totalRevenue"`
	Revenue        string `json:"revenue"`
	Completed      string `json:"completed"`
	Active         string `json:"active"`
	RevenueTrend   string `json:"revenueTrend"`
	Week           string `json:"week"`
	Month          string `json:"month"`
	Year           string `json:"year"`
	RecentRevenue  string `json:"recentRevenue"`
	ViewAll        string `json:"viewAll"`
	NewRevenueCreated string `json:"newRevenueCreated"`
	RevenueCompleted  string `json:"revenueCompleted"`
	RevenueUpdated    string `json:"revenueUpdated"`
	RevenueCancelled  string `json:"revenueCancelled"`
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
	PricePlaceholder         string `json:"pricePlaceholder"`
	SelectOption             string `json:"selectOption"`
	InitialValues            string `json:"initialValues"`
	InitialValuesPlaceholder string `json:"initialValuesPlaceholder"`
	Required                 string `json:"required"`
	Option                   string `json:"option"`
	SelectAttribute          string `json:"selectAttribute"`
	AllAttributesAssigned    string `json:"allAttributesAssigned"`
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
	Info        string `json:"info"`
	Variants    string `json:"variants"`
	Attributes  string `json:"attributes"`
	Pricing     string `json:"pricing"`
	Options     string `json:"options"`
	Images      string `json:"images"`
	Stock       string `json:"stock"`
	Attachments string `json:"attachments"`
	AuditTrail  string `json:"auditTrail"`
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

type ProductOptionColumnLabels struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	DataType    string `json:"dataType"`
	ValuesCount string `json:"valuesCount"`
	SortOrder   string `json:"sortOrder"`
	Status      string `json:"status"`
}

type ProductOptionFormLabels struct {
	Name                    string `json:"name"`
	NamePlaceholder         string `json:"namePlaceholder"`
	Code                    string `json:"code"`
	CodePlaceholder         string `json:"codePlaceholder"`
	DataType                string `json:"dataType"`
	SortOrder               string `json:"sortOrder"`
	MinValue                string `json:"minValue"`
	MaxValue                string `json:"maxValue"`
	Active                  string `json:"active"`
	InitialValues           string `json:"initialValues"`
	InitialValuesPlaceholder string `json:"initialValuesPlaceholder"`
	Required                string `json:"required"`
	Form                    ProductOptionFormInnerLabels `json:"form"`
}

// ProductOptionFormInnerLabels holds nested form labels referenced by the template as .Labels.Form.*
type ProductOptionFormInnerLabels struct {
	InitialValues           string `json:"initialValues"`
	InitialValuesPlaceholder string `json:"initialValuesPlaceholder"`
	Required                string `json:"required"`
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
	Product         string `json:"product"`
	SelectProduct   string `json:"selectProduct"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`
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
	AddProductLine string `json:"addProductLine"`
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
			AddProductLine: "Add Product Line",
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
// Expenditure labels
// ---------------------------------------------------------------------------

// ExpenditureLabels holds all translatable strings for the expenditure module
// (purchase + expense views).
type ExpenditureLabels struct {
	Labels              ExpenditureLabelNames              `json:"labels"`
	Page                ExpenditurePageLabels              `json:"page"`
	Buttons             ExpenditureButtonLabels            `json:"buttons"`
	Columns             ExpenditureColumnLabels            `json:"columns"`
	Empty               ExpenditureEmptyLabels             `json:"empty"`
	Form                ExpenditureFormLabels              `json:"form"`
	Status              ExpenditureStatusLabels            `json:"status"`
	Types               ExpenditureTypeLabels              `json:"types"`
	Actions             ExpenditureActionLabels            `json:"actions"`
	Bulk                ExpenditureBulkLabels              `json:"bulkActions"`
	Detail              ExpenditureDetailLabels            `json:"detail"`
	Errors              ExpenditureErrorLabels             `json:"errors"`
	Category            ExpenditureCategoryLabels          `json:"category"`
	PaymentMethod       ExpenditurePaymentMethodLabels     `json:"paymentMethod"`
	DisbursementCategory ExpenditureDisbursementCategoryLabels `json:"disbursementCategory"`
	Schedule            ExpenditureScheduleLabels          `json:"schedule"`
	LineItemForm        ExpenditureLineItemFormLabels      `json:"lineItemForm"`
	DisbursementForm    ExpenditureDisbursementFormLabels  `json:"disbursementForm"`
	PurchaseOrder       PurchaseOrderLabels                `json:"purchaseOrder"`
}

// ExpenditureCategoryLabels holds translatable strings for the expenditure
// category settings list and CRUD drawer.
type ExpenditureCategoryLabels struct {
	Page    ExpenditureCategoryPageLabels    `json:"page"`
	Columns ExpenditureCategoryColumnLabels  `json:"columns"`
	Empty   ExpenditureCategoryEmptyLabels   `json:"empty"`
	Form    ExpenditureCategoryFormLabels    `json:"form"`
	Actions ExpenditureCategoryActionLabels  `json:"actions"`
	Errors  ExpenditureCategoryErrorLabels   `json:"errors"`
	Confirm ExpenditureCategoryConfirmLabels `json:"confirm"`
	Buttons ExpenditureCategoryButtonLabels  `json:"buttons"`
}

type ExpenditureCategoryPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type ExpenditureCategoryButtonLabels struct {
	AddCategory string `json:"addCategory"`
}

type ExpenditureCategoryColumnLabels struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type ExpenditureCategoryEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ExpenditureCategoryFormLabels struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ExpenditureCategoryActionLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type ExpenditureCategoryErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	InvalidFormData  string `json:"invalidFormData"`
}

type ExpenditureCategoryConfirmLabels struct {
	DeleteTitle   string `json:"deleteTitle"`
	DeleteMessage string `json:"deleteMessage"`
}

// ExpenditureErrorLabels holds error messages for the expenditure action handlers.
type ExpenditureErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
}

type ExpenditureLabelNames struct {
	Name           string `json:"name"`
	NamePlural     string `json:"namePlural"`
	Purchase       string `json:"purchase"`
	PurchasePlural string `json:"purchasePlural"`
	PurchaseOrder  string `json:"purchaseOrder"`
	Expense        string `json:"expense"`
	ExpensePlural  string `json:"expensePlural"`
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
	VendorName                 string `json:"vendorName"`
	VendorNamePlaceholder      string `json:"vendorNamePlaceholder"`
	ExpenditureDate            string `json:"expenditureDate"`
	TotalAmount                string `json:"totalAmount"`
	Currency                   string `json:"currency"`
	Status                     string `json:"status"`
	ReferenceNumber            string `json:"referenceNumber"`
	ReferenceNumberPlaceholder string `json:"referenceNumberPlaceholder"`
	PaymentTerms               string `json:"paymentTerms"`
	DueDate                    string `json:"dueDate"`
	ApprovedBy                 string `json:"approvedBy"`
	ExpenditureType            string `json:"expenditureType"`
	ExpenditureCategory        string `json:"expenditureCategory"`
	Notes                      string `json:"notes"`
	NotesPlaceholder           string `json:"notesPlaceholder"`
	SectionInfo                string `json:"sectionInfo"`
	SectionVendor              string `json:"sectionVendor"`
	SectionPayment             string `json:"sectionPayment"`
	SectionNotes               string `json:"sectionNotes"`
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
	PurchasePageTitle    string `json:"purchasePageTitle"`
	ExpensePageTitle     string `json:"expensePageTitle"`
	VendorInfo           string `json:"vendorInfo"`
	VendorName           string `json:"vendorName"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Status               string `json:"status"`
	Type                 string `json:"type"`
	Category             string `json:"category"`
	ReferenceNumber      string `json:"referenceNumber"`
	PaymentTerms         string `json:"paymentTerms"`
	DueDate              string `json:"dueDate"`
	ApprovedBy           string `json:"approvedBy"`
	Notes                string `json:"notes"`
	LineItems            string `json:"lineItems"`
	Description          string `json:"description"`
	Quantity             string `json:"quantity"`
	UnitPrice            string `json:"unitPrice"`
	Total                string `json:"total"`
	SubTotal             string `json:"subTotal"`
	GrandTotal           string `json:"grandTotal"`
	TabBasicInfo         string `json:"tabBasicInfo"`
	TabLineItems         string `json:"tabLineItems"`
	TabPayment           string `json:"tabPayment"`
	TabAuditTrail        string `json:"tabAuditTrail"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditAction          string `json:"auditAction"`
	AuditUser            string `json:"auditUser"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	// Additional fields used in the expense detail template
	Title               string `json:"title"`
	InfoSection         string `json:"infoSection"`
	Name                string `json:"name"`
	PaymentSummary      string `json:"paymentSummary"`
	TotalAmount         string `json:"totalAmount"`
	Paid                string `json:"paid"`
	Outstanding         string `json:"outstanding"`
	PaymentStatus       string `json:"paymentStatus"`
	UpdateStatus        string `json:"updateStatus"`
	SaveStatus          string `json:"saveStatus"`
	Payment             string `json:"payment"`
	Pay                 string `json:"pay"`
	AddItem             string `json:"addItem"`
	EmptyTitle          string `json:"emptyTitle"`
	EmptyMessage        string `json:"emptyMessage"`
	TabDetails          string `json:"tabDetails"`
	TabPayments         string `json:"tabPayments"`
}

// ExpenditurePaymentMethodLabels holds translatable strings for disbursement payment methods.
type ExpenditurePaymentMethodLabels struct {
	Cash         string `json:"cash"`
	BankTransfer string `json:"bankTransfer"`
	Check        string `json:"check"`
	GCash        string `json:"gcash"`
	Other        string `json:"other"`
}

// ExpenditureDisbursementCategoryLabels holds translatable strings for disbursement categories.
type ExpenditureDisbursementCategoryLabels struct {
	SupplierPayment string `json:"supplierPayment"`
	Payroll         string `json:"payroll"`
	Rent            string `json:"rent"`
	Utilities       string `json:"utilities"`
	Other           string `json:"other"`
}

// ExpenditureScheduleLabels holds translatable strings for the payment schedule tab.
type ExpenditureScheduleLabels struct {
	Scheduled   string `json:"scheduled"`
	Paid        string `json:"paid"`
	Remaining   string `json:"remaining"`
	DueDate     string `json:"dueDate"`
	AmountDue   string `json:"amountDue"`
	PaidAmount  string `json:"paidAmount"`
	PaidDate    string `json:"paidDate"`
	Reference   string `json:"reference"`
	EmptyTitle  string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// ExpenditureLineItemFormLabels holds translatable strings for the line item drawer form.
type ExpenditureLineItemFormLabels struct {
	EditTitle           string `json:"editTitle"`
	Description         string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Quantity            string `json:"quantity"`
	UnitPrice           string `json:"unitPrice"`
	Notes               string `json:"notes"`
	Save                string `json:"save"`
	Cancel              string `json:"cancel"`
}

// ExpenditureDisbursementFormLabels holds translatable strings for the pay (disbursement) drawer form.
type ExpenditureDisbursementFormLabels struct {
	Reference          string `json:"reference"`
	ReferencePlaceholder string `json:"referencePlaceholder"`
	Payee              string `json:"payee"`
	Amount             string `json:"amount"`
	Currency           string `json:"currency"`
	CurrencyPlaceholder string `json:"currencyPlaceholder"`
	PaymentMethod      string `json:"paymentMethod"`
	Category           string `json:"category"`
	ApprovedBy         string `json:"approvedBy"`
	ApproverPlaceholder string `json:"approverPlaceholder"`
}

// ---------------------------------------------------------------------------
// Purchase Order labels
// ---------------------------------------------------------------------------

// PurchaseOrderLabels holds all translatable strings for the purchase order module.
type PurchaseOrderLabels struct {
	Labels      PurchaseOrderLabelNames      `json:"labels"`
	Page        PurchaseOrderPageLabels      `json:"page"`
	Buttons     PurchaseOrderButtonLabels    `json:"buttons"`
	Columns     PurchaseOrderColumnLabels    `json:"columns"`
	Empty       PurchaseOrderEmptyLabels     `json:"empty"`
	Form        PurchaseOrderFormLabels      `json:"form"`
	Status      PurchaseOrderStatusLabels    `json:"status"`
	POTypes     PurchaseOrderPOTypeLabels    `json:"poTypes"`
	LineTypes   PurchaseOrderLineTypeLabels  `json:"lineTypes"`
	Actions     PurchaseOrderActionLabels    `json:"actions"`
	Bulk        PurchaseOrderBulkLabels      `json:"bulkActions"`
	Detail      PurchaseOrderDetailLabels    `json:"detail"`
	LineItems   PurchaseOrderLineItemLabels  `json:"lineItems"`
	Receipt     PurchaseOrderReceiptLabels   `json:"receipt"`
}

type PurchaseOrderLabelNames struct {
	Name         string `json:"name"`
	NamePlural   string `json:"namePlural"`
	LineItem     string `json:"lineItem"`
	LineItemPlural string `json:"lineItemPlural"`
}

type PurchaseOrderPageLabels struct {
	Heading                  string `json:"heading"`
	Caption                  string `json:"caption"`
	HeadingDraft             string `json:"headingDraft"`
	HeadingPendingApproval   string `json:"headingPendingApproval"`
	HeadingApproved          string `json:"headingApproved"`
	HeadingPartiallyReceived string `json:"headingPartiallyReceived"`
	HeadingFullyReceived     string `json:"headingFullyReceived"`
	HeadingBilled            string `json:"headingBilled"`
	HeadingClosed            string `json:"headingClosed"`
	HeadingCancelled         string `json:"headingCancelled"`
	Dashboard                string `json:"dashboard"`
}

type PurchaseOrderButtonLabels struct {
	Add         string `json:"add"`
	AddLineItem string `json:"addLineItem"`
}

type PurchaseOrderColumnLabels struct {
	PONumber       string `json:"poNumber"`
	POType         string `json:"poType"`
	Supplier       string `json:"supplier"`
	Location       string `json:"location"`
	OrderDate      string `json:"orderDate"`
	Status         string `json:"status"`
	Currency       string `json:"currency"`
	Subtotal       string `json:"subtotal"`
	TaxAmount      string `json:"taxAmount"`
	TotalAmount    string `json:"totalAmount"`
	PaymentTerms   string `json:"paymentTerms"`
	ShippingTerms  string `json:"shippingTerms"`
	ApprovedBy     string `json:"approvedBy"`
	ReferenceNumber string `json:"referenceNumber"`
	Notes          string `json:"notes"`
}

type PurchaseOrderEmptyLabels struct {
	Title                    string `json:"title"`
	Message                  string `json:"message"`
	DraftTitle               string `json:"draftTitle"`
	DraftMessage             string `json:"draftMessage"`
	PendingApprovalTitle     string `json:"pendingApprovalTitle"`
	PendingApprovalMessage   string `json:"pendingApprovalMessage"`
	ApprovedTitle            string `json:"approvedTitle"`
	ApprovedMessage          string `json:"approvedMessage"`
	PartiallyReceivedTitle   string `json:"partiallyReceivedTitle"`
	PartiallyReceivedMessage string `json:"partiallyReceivedMessage"`
	FullyReceivedTitle       string `json:"fullyReceivedTitle"`
	FullyReceivedMessage     string `json:"fullyReceivedMessage"`
	BilledTitle              string `json:"billedTitle"`
	BilledMessage            string `json:"billedMessage"`
	ClosedTitle              string `json:"closedTitle"`
	ClosedMessage            string `json:"closedMessage"`
	CancelledTitle           string `json:"cancelledTitle"`
	CancelledMessage         string `json:"cancelledMessage"`
}

type PurchaseOrderFormLabels struct {
	PONumber                 string `json:"poNumber"`
	PONumberPlaceholder      string `json:"poNumberPlaceholder"`
	POType                   string `json:"poType"`
	SelectPOType             string `json:"selectPoType"`
	Supplier                 string `json:"supplier"`
	SelectSupplier           string `json:"selectSupplier"`
	Location                 string `json:"location"`
	SelectLocation           string `json:"selectLocation"`
	OrderDate                string `json:"orderDate"`
	Currency                 string `json:"currency"`
	Subtotal                 string `json:"subtotal"`
	TaxAmount                string `json:"taxAmount"`
	TotalAmount              string `json:"totalAmount"`
	PaymentTerms             string `json:"paymentTerms"`
	ShippingTerms            string `json:"shippingTerms"`
	ApprovedBy               string `json:"approvedBy"`
	ReferenceNumber          string `json:"referenceNumber"`
	ReferenceNumberPlaceholder string `json:"referenceNumberPlaceholder"`
	Notes                    string `json:"notes"`
	NotesPlaceholder         string `json:"notesPlaceholder"`
	SectionInfo              string `json:"sectionInfo"`
	SectionSupplier          string `json:"sectionSupplier"`
	SectionFinancials        string `json:"sectionFinancials"`
	SectionNotes             string `json:"sectionNotes"`
}

type PurchaseOrderStatusLabels struct {
	Draft              string `json:"draft"`
	PendingApproval    string `json:"pending_approval"`
	Approved           string `json:"approved"`
	PartiallyReceived  string `json:"partially_received"`
	FullyReceived      string `json:"fully_received"`
	Billed             string `json:"billed"`
	Closed             string `json:"closed"`
	Cancelled          string `json:"cancelled"`
}

type PurchaseOrderPOTypeLabels struct {
	Standard string `json:"standard"`
	Blanket  string `json:"blanket"`
	Contract string `json:"contract"`
}

type PurchaseOrderLineTypeLabels struct {
	Goods   string `json:"goods"`
	Service string `json:"service"`
	Expense string `json:"expense"`
}

type PurchaseOrderActionLabels struct {
	Cancel         string `json:"cancel"`
	Close          string `json:"close"`
	ConfirmReceipt string `json:"confirmReceipt"`
	Create         string `json:"create"`
	Delete         string `json:"delete"`
	Edit           string `json:"edit"`
	Approve        string `json:"approve"`
	Receive        string `json:"receive"`
	Reject         string `json:"reject"`
	View           string `json:"view"`
}

type PurchaseOrderBulkLabels struct {
	Delete  string `json:"delete"`
	Approve string `json:"approve"`
	Close   string `json:"close"`
}

// PurchaseOrderDetailLabels holds translatable strings for the PO detail page.
type PurchaseOrderDetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	Title                string `json:"title"`
	InfoSection          string `json:"supplierInfo"`
	Supplier             string `json:"supplier"`
	Location             string `json:"location"`
	OrderDate            string `json:"orderDate"`
	PONumber             string `json:"poNumber"`
	POType               string `json:"poType"`
	Status               string `json:"status"`
	Currency             string `json:"currency"`
	Subtotal             string `json:"subtotal"`
	TaxAmount            string `json:"taxAmount"`
	TotalAmount          string `json:"totalAmount"`
	PaymentTerms         string `json:"paymentTerms"`
	ShippingTerms        string `json:"shippingTerms"`
	ApprovedBy           string `json:"approvedBy"`
	ReferenceNumber      string `json:"referenceNumber"`
	Notes                string `json:"notes"`
	LineItems            string `json:"lineItems"`
	Description          string `json:"description"`
	LineType             string `json:"lineType"`
	LineNumber           string `json:"lineNumber"`
	QuantityOrdered      string `json:"quantityOrdered"`
	QuantityReceived     string `json:"quantityReceived"`
	QuantityBilled       string `json:"quantityBilled"`
	UnitPrice            string `json:"unitPrice"`
	TotalPrice           string `json:"totalPrice"`
	SubTotal             string `json:"subTotal"`
	GrandTotal           string `json:"grandTotal"`
	TabBasicInfo         string `json:"tabBasicInfo"`
	TabLineItems         string `json:"tabLineItems"`
	TabReceiving         string `json:"tabReceiving"`
	TabAuditTrail        string `json:"tabAuditTrail"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditAction          string `json:"auditAction"`
	AuditUser            string `json:"auditUser"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	Total                string `json:"total"`
	AddLineItem          string `json:"addLineItem"`
	NoLineItems          string `json:"noLineItems"`
	ConfirmReceiptBtn    string `json:"confirmReceiptBtn"`
}

// PurchaseOrderLineItemLabels holds translatable strings for the PO line item drawer form.
type PurchaseOrderLineItemLabels struct {
	AddItem                string `json:"addItem"`
	AddLineItem            string `json:"addLineItem"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	EditItem               string `json:"editItem"`
	EditLineItem           string `json:"editLineItem"`
	InventoryItem          string `json:"inventoryItem"`
	LineNumber             string `json:"lineNumber"`
	LineType               string `json:"lineType"`
	Location               string `json:"location"`
	Locked                 string `json:"locked"`
	NoItems                string `json:"noItems"`
	Notes                  string `json:"notes"`
	Product                string `json:"product"`
	QtyOrdered             string `json:"qtyOrdered"`
	QuantityBilled         string `json:"quantityBilled"`
	QuantityOrdered        string `json:"quantityOrdered"`
	QuantityReceived       string `json:"quantityReceived"`
	RemoveItem             string `json:"removeItem"`
	RemoveLineItem         string `json:"removeLineItem"`
	SelectItem             string `json:"selectItem"`
	TotalPrice             string `json:"totalPrice"`
	TypeExpense            string `json:"typeExpense"`
	TypeGoods              string `json:"typeGoods"`
	TypeService            string `json:"typeService"`
	UnitPrice              string `json:"unitPrice"`
	Type                   string `json:"type"`
	ProductID              string `json:"productId"`
	InventoryItemID        string `json:"inventoryItemId"`
	LocationID             string `json:"locationId"`
	Save                   string `json:"save"`
	Cancel                 string `json:"cancel"`
}

// PurchaseOrderReceiptLabels holds translatable strings for the confirm receipt drawer form.
type PurchaseOrderReceiptLabels struct {
	AutoConfirmed       string `json:"autoConfirmed"`
	NoLines             string `json:"noLines"`
	OverReceiptError    string `json:"overReceiptError"`
	PartialSuccess      string `json:"partialSuccess"`
	QtyToReceive        string `json:"qtyToReceive"`
	ReceiptDate         string `json:"receiptDate"`
	ReceivingLocation   string `json:"receivingLocation"`
	ServiceRendered     string `json:"serviceRendered"`
	Success             string `json:"success"`
	Title               string `json:"title"`
	AllReceived         string `json:"allReceived"`
	Description         string `json:"description"`
	Type                string `json:"type"`
	Ordered             string `json:"ordered"`
	Received            string `json:"received"`
	Remaining           string `json:"remaining"`
	ConfirmButton       string `json:"confirmButton"`
	Cancel              string `json:"cancel"`
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
	PageTitle            string `json:"pageTitle"`
	TitlePrefix          string `json:"titlePrefix"`
	PaymentInfo          string `json:"paymentInfo"`
	Customer             string `json:"customer"`
	Date                 string `json:"date"`
	Amount               string `json:"amount"`
	Currency             string `json:"currency"`
	Status               string `json:"status"`
	Method               string `json:"method"`
	Reference            string `json:"reference"`
	Notes                string `json:"notes"`
	TabBasicInfo         string `json:"tabBasicInfo"`
	TabAttachments       string `json:"tabAttachments"`
	TabAuditTrail        string `json:"tabAuditTrail"`
	AuditAction          string `json:"auditAction"`
	AuditUser            string `json:"auditUser"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditTrailDesc       string `json:"auditTrailDesc"`
}

type CollectionStatusLabels struct {
	Pending   string `json:"pending"`
	Completed string `json:"completed"`
	Failed    string `json:"failed"`
}

type CollectionConfirmLabels struct {
	MarkComplete          string `json:"markComplete"`
	MarkCompleteMessage   string `json:"markCompleteMessage"`
	Reactivate            string `json:"reactivate"`
	ReactivateMessage     string `json:"reactivateMessage"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	BulkComplete          string `json:"bulkComplete"`
	BulkCompleteMessage   string `json:"bulkCompleteMessage"`
	BulkReactivate        string `json:"bulkReactivate"`
	BulkReactivateMessage string `json:"bulkReactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
}

type CollectionErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
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
			Customer:                "Customer",
			Date:                    "Date",
			Amount:                  "Amount",
			Currency:                "Currency",
			Reference:               "Reference",
			ReferencePlaceholder:    "e.g. INV-001",
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
			PageTitle:            "Collection Details",
			TitlePrefix:          "Collection #",
			PaymentInfo:          "Payment Information",
			Customer:             "Customer",
			Date:                 "Date",
			Amount:               "Amount",
			Currency:             "Currency",
			Status:               "Status",
			Method:               "Payment Method",
			Reference:            "Reference",
			Notes:                "Notes",
			TabBasicInfo:         "Basic Info",
			TabAttachments:       "Attachments",
			TabAuditTrail:        "Audit Trail",
			AuditAction:          "Action",
			AuditUser:            "User",
			AuditEmptyTitle:      "No audit records",
			AuditEmptyMessage:    "No audit trail entries yet.",
			AuditTrailComingSoon: "Audit trail coming soon.",
			AuditTrailDesc:       "Audit trail for collection changes is coming soon.",
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
	Payee                   string `json:"payee"`
	PayeePlaceholder        string `json:"payeePlaceholder"`
	Date                    string `json:"date"`
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	Reference               string `json:"reference"`
	ReferencePlaceholder    string `json:"referencePlaceholder"`
	PaymentMethod           string `json:"paymentMethod"`
	Category                string `json:"category"`
	Status                  string `json:"status"`
	Notes                   string `json:"notes"`
	NotesPlaceholder        string `json:"notesPlaceholder"`
	ApprovedBy              string `json:"approvedBy"`
	AmountPlaceholder       string `json:"amountPlaceholder"`
	CurrencyPlaceholder     string `json:"currencyPlaceholder"`
	MethodCash              string `json:"methodCash"`
	MethodBankTransfer      string `json:"methodBankTransfer"`
	MethodCheck             string `json:"methodCheck"`
	MethodGCash             string `json:"methodGCash"`
	MethodOther             string `json:"methodOther"`
	StatusDraft             string `json:"statusDraft"`
	StatusPending           string `json:"statusPending"`
	StatusApproved          string `json:"statusApproved"`
	StatusPaid              string `json:"statusPaid"`
	StatusCancelled         string `json:"statusCancelled"`
	TypeSupplierPayment     string `json:"typeSupplierPayment"`
	TypePayroll             string `json:"typePayroll"`
	TypeRent                string `json:"typeRent"`
	TypeUtilities           string `json:"typeUtilities"`
	TypeOther               string `json:"typeOther"`
	ApproverNamePlaceholder string `json:"approverNamePlaceholder"`
	LinkToBill              string `json:"linkToBill"`
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
	TabAttachments    string `json:"tabAttachments"`
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
	Submit                string `json:"submit"`
	SubmitMessage         string `json:"submitMessage"`
	Approve               string `json:"approve"`
	ApproveMessage        string `json:"approveMessage"`
	MarkPaid              string `json:"markPaid"`
	MarkPaidMessage       string `json:"markPaidMessage"`
	Cancel                string `json:"cancel"`
	CancelMessage         string `json:"cancelMessage"`
	Reactivate            string `json:"reactivate"`
	ReactivateMessage     string `json:"reactivateMessage"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	BulkSubmit            string `json:"bulkSubmit"`
	BulkSubmitMessage     string `json:"bulkSubmitMessage"`
	BulkApprove           string `json:"bulkApprove"`
	BulkApproveMessage    string `json:"bulkApproveMessage"`
	BulkMarkPaid          string `json:"bulkMarkPaid"`
	BulkMarkPaidMessage   string `json:"bulkMarkPaidMessage"`
	BulkCancel            string `json:"bulkCancel"`
	BulkCancelMessage     string `json:"bulkCancelMessage"`
	BulkReactivate        string `json:"bulkReactivate"`
	BulkReactivateMessage string `json:"bulkReactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
}

type DisbursementErrorLabels struct {
	PermissionDenied  string `json:"permissionDenied"`
	InvalidFormData   string `json:"invalidFormData"`
	NotFound          string `json:"notFound"`
	IDRequired        string `json:"idRequired"`
	NoIDsProvided     string `json:"noIDsProvided"`
	InvalidStatus     string `json:"invalidStatus"`
	InvalidTransition string `json:"invalidTransition"`
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
			Payee:                   "Payee",
			PayeePlaceholder:        "Enter payee name",
			Date:                    "Date",
			Amount:                  "Amount",
			Currency:                "Currency",
			Reference:               "Reference",
			ReferencePlaceholder:    "e.g. DISB-001",
			PaymentMethod:           "Payment Method",
			Category:                "Category",
			Status:                  "Status",
			Notes:                   "Notes",
			NotesPlaceholder:        "Additional notes...",
			ApprovedBy:              "Approved By",
			AmountPlaceholder:       "0.00",
			CurrencyPlaceholder:     "PHP",
			MethodCash:              "Cash",
			MethodBankTransfer:      "Bank Transfer",
			MethodCheck:             "Check",
			MethodGCash:             "GCash",
			MethodOther:             "Other",
			StatusDraft:             "Draft",
			StatusPending:           "Pending",
			StatusApproved:          "Approved",
			StatusPaid:              "Paid",
			StatusCancelled:         "Cancelled",
			TypeSupplierPayment:     "Supplier Payment",
			TypePayroll:             "Payroll",
			TypeRent:                "Rent",
			TypeUtilities:           "Utilities",
			TypeOther:               "Other",
			ApproverNamePlaceholder: "Approver name",
			LinkToBill:              "Link to Bill",
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
			TabAttachments:    "Attachments",
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
	Page            PlanPageLabels         `json:"page"`
	Buttons         PlanButtonLabels       `json:"buttons"`
	Columns         PlanColumnLabels       `json:"columns"`
	Empty           PlanEmptyLabels        `json:"empty"`
	Form            PlanFormLabels         `json:"form"`
	Actions         PlanActionLabels       `json:"actions"`
	Bulk            PlanBulkLabels         `json:"bulkActions"`
	Status          PlanStatusLabels       `json:"status"`
	Detail          PlanDetailLabels       `json:"detail"`
	Tabs            PlanTabLabels          `json:"tabs"`
	Confirm         PlanConfirmLabels      `json:"confirm"`
	Errors          PlanErrorLabels        `json:"errors"`
	PricePlanForm   PricePlanFormLabels    `json:"pricePlanForm"`
	ProductPlanForm ProductPlanFormLabels  `json:"productPlanForm"`
}

// ProductPlanFormLabels holds translatable labels for the ProductPlan add/edit form within a plan.
type ProductPlanFormLabels struct {
	Product            string `json:"product"`
	ProductPlaceholder string `json:"productPlaceholder"`
	SelectProduct      string `json:"selectProduct"`
	Active             string `json:"active"`
}

// PricePlanFormLabels holds translatable labels for the PricePlan add/edit form.
type PricePlanFormLabels struct {
	Name                string `json:"name"`
	NamePlaceholder     string `json:"namePlaceholder"`
	Description         string `json:"description"`
	DescPlaceholder     string `json:"descriptionPlaceholder"`
	Amount              string `json:"amount"`
	AmountPlaceholder   string `json:"amountPlaceholder"`
	Currency            string `json:"currency"`
	CurrencyPlaceholder string `json:"currencyPlaceholder"`
	DurationValue       string `json:"durationValue"`
	DurationUnit        string `json:"durationUnit"`
	Location            string `json:"location"`
	LocationPlaceholder string `json:"locationPlaceholder"`
	SelectLocation      string `json:"selectLocation"`
	Active              string `json:"active"`
}

type PlanPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type PlanButtonLabels struct {
	AddPlan       string `json:"addPlan"`
	AddPricePlan  string `json:"addPricePlan"`
	AddProduct    string `json:"addProduct"`
}

type PlanColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Interval    string `json:"interval"`
	Price       string `json:"price"`
	Status      string `json:"status"`
	Product     string `json:"product"`
	PricePlan   string `json:"pricePlan"`
	Duration    string `json:"duration"`
	Location    string `json:"location"`
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
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type PlanBulkLabels struct {
	Delete string `json:"delete"`
}

type PlanStatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type PlanErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
	NoPermission     string `json:"noPermission"`
	CannotDelete     string `json:"cannotDelete"`
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
	Sections            PlanFormSectionLabels `json:"sections"`

	// Fulfillment type option labels
	TypeSchedule string `json:"typeSchedule"`
	TypeLicense  string `json:"typeLicense"`
	TypeContent  string `json:"typeContent"`
	TypePhysical string `json:"typePhysical"`
}

type PlanDetailLabels struct {
	PageTitle             string `json:"pageTitle"`
	Price                 string `json:"price"`
	Currency              string `json:"currency"`
	Status                string `json:"status"`
	Description           string `json:"description"`
	FulfillmentType       string `json:"fulfillmentType"`
	CreatedDate           string `json:"createdDate"`
	ModifiedDate          string `json:"modifiedDate"`
	NoProductsAssigned    string `json:"noProductsAssigned"`
	NoProductsAssignedMsg string `json:"noProductsAssignedMsg"`
	NoProductsDesc        string `json:"noProductsDesc"`
	NoPricePlans          string `json:"noPricePlans"`
	NoPricePlansMsg       string `json:"noPricePlansMsg"`
	NoPricePlansDesc      string `json:"noPricePlansDesc"`
	AuditTrailComingSoon  string `json:"auditTrailComingSoon"`
}

type PlanTabLabels struct {
	Info        string `json:"info"`
	Products    string `json:"products"`
	PriceLists  string `json:"priceLists"`
	Attachments string `json:"attachments"`
	AuditTrail  string `json:"auditTrail"`
}

type PlanConfirmLabels struct {
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
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
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
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
	Customer                  string `json:"customer"`
	CustomerPlaceholder       string `json:"customerPlaceholder"`
	Plan                      string `json:"plan"`
	PlanPlaceholder           string `json:"planPlaceholder"`
	StartDate                 string `json:"startDate"`
	EndDate                   string `json:"endDate"`
	Active                    string `json:"active"`
	Notes                     string `json:"notes"`
	NotesPlaceholder          string `json:"notesPlaceholder"`
	CustomerSearchPlaceholder string `json:"customerSearchPlaceholder"`
	PlanSearchPlaceholder     string `json:"planSearchPlaceholder"`
	CustomerNoResults         string `json:"customerNoResults"`
	PlanNoResults             string `json:"planNoResults"`
	Code                      string `json:"code"`
	CodePlaceholder           string `json:"codePlaceholder"`
}

type SubscriptionDetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	Customer             string `json:"customer"`
	Plan                 string `json:"plan"`
	StartDate            string `json:"startDate"`
	EndDate              string `json:"endDate"`
	Status               string `json:"status"`
	CreatedDate          string `json:"createdDate"`
	ModifiedDate         string `json:"modifiedDate"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditTrailDesc       string `json:"auditTrailDesc"`
}

type SubscriptionTabLabels struct {
	Info        string `json:"info"`
	History     string `json:"history"`
	Attachments string `json:"attachments"`
	AuditTrail  string `json:"auditTrail"`
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
			AddPlan:      "Add Plan",
			AddPricePlan: "Add Price Plan",
			AddProduct:   "Add Product",
		},
		Columns: PlanColumnLabels{
			Name:        "Name",
			Description: "Description",
			Interval:    "Interval",
			Price:       "Price",
			Status:      "Status",
			Product:     "Product",
			PricePlan:   "Price Plan",
			Duration:    "Duration",
			Location:    "Location",
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
			TypeSchedule:        "Schedule",
			TypeLicense:         "License",
			TypeContent:         "Content",
			TypePhysical:        "Physical",
			Sections: PlanFormSectionLabels{
				Basic:    "Basic Information",
				Services: "Assigned Products",
			},
		},
		Actions: PlanActionLabels{
			View:       "View Plan",
			Edit:       "Edit Plan",
			Delete:     "Delete Plan",
			Activate:   "Activate Plan",
			Deactivate: "Deactivate Plan",
		},
		Bulk: PlanBulkLabels{
			Delete: "Delete Selected",
		},
		Status: PlanStatusLabels{
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Detail: PlanDetailLabels{
			PageTitle:             "Plan Details",
			Price:                 "Price",
			Currency:              "Currency",
			Status:                "Status",
			Description:           "Description",
			FulfillmentType:       "Fulfillment Type",
			CreatedDate:           "Created",
			ModifiedDate:          "Last Modified",
			NoProductsAssigned:    "No products assigned",
			NoProductsAssignedMsg: "No products have been linked to this plan yet.",
			NoProductsDesc:        "No products have been linked to this plan yet.",
			NoPricePlans:          "No price plans",
			NoPricePlansMsg:       "No price plans have been configured for this plan yet.",
			NoPricePlansDesc:      "No price plans have been configured for this plan yet.",
			AuditTrailComingSoon:  "Audit trail coming soon.",
		},
		Tabs: PlanTabLabels{
			Info:        "Information",
			Products:    "Products",
			PriceLists:  "Price Lists",
			Attachments: "Attachments",
			AuditTrail:  "Audit Trail",
		},
		Confirm: PlanConfirmLabels{
			Delete:                "Delete Plan",
			DeleteMessage:         "Are you sure you want to delete \"%s\"? This action cannot be undone.",
			Activate:              "Activate Plan",
			ActivateMessage:       "Are you sure you want to activate \"%s\"?",
			Deactivate:            "Deactivate Plan",
			DeactivateMessage:     "Are you sure you want to deactivate \"%s\"?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Are you sure you want to activate the selected plans?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Are you sure you want to deactivate the selected plans?",
			BulkDelete:            "Delete Selected",
			BulkDeleteMessage:     "Are you sure you want to delete the selected plans? This action cannot be undone.",
		},
		Errors: PlanErrorLabels{
			PermissionDenied: "You do not have permission to perform this action",
			InvalidFormData:  "Invalid form data. Please check your inputs and try again.",
			NotFound:         "Plan not found",
			IDRequired:       "Plan ID is required",
			NoIDsProvided:    "No plan IDs provided",
			InvalidStatus:    "Invalid status",
			NoPermission:     "No permission",
			CannotDelete:     "This plan cannot be deleted because it has products or rate cards assigned",
		},
		PricePlanForm: PricePlanFormLabels{
			Name:                "Price Plan Name",
			NamePlaceholder:     "Enter price plan name",
			Description:         "Description",
			DescPlaceholder:     "Enter description...",
			Amount:              "Amount",
			AmountPlaceholder:   "0.00",
			Currency:            "Currency",
			CurrencyPlaceholder: "e.g. PHP",
			DurationValue:       "Duration",
			DurationUnit:        "Unit",
			Location:            "Location",
			LocationPlaceholder: "Select a location...",
			SelectLocation:      "— No location (all locations) —",
			Active:              "Active",
		},
		ProductPlanForm: ProductPlanFormLabels{
			Product:            "Product",
			ProductPlaceholder: "Select a product...",
			SelectProduct:      "— Select a product —",
			Active:             "Active",
		},
	}
}

// ---------------------------------------------------------------------------
// Price Plan labels
// ---------------------------------------------------------------------------

// PricePlanLabels holds all labels for the standalone price plan (rate card) module.
type PricePlanLabels struct {
	Page    PricePlanPageLabels    `json:"page"`
	Buttons PricePlanButtonLabels  `json:"buttons"`
	Columns PricePlanColumnLabels2 `json:"columns"`
	Empty   PricePlanEmptyLabels   `json:"empty"`
	Form    PricePlanFormLabels    `json:"form"`
	Actions PricePlanActionLabels  `json:"actions"`
	Bulk    PricePlanBulkLabels    `json:"bulk"`
	Detail  PricePlanDetailLabels2 `json:"detail"`
	Tabs    PricePlanTabLabels2    `json:"tabs"`
	Confirm PricePlanConfirmLabels `json:"confirm"`
	Errors  PricePlanErrorLabels   `json:"errors"`
}

type PricePlanPageLabels struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	ActiveTitle   string `json:"activeTitle"`
	InactiveTitle string `json:"inactiveTitle"`
}

type PricePlanButtonLabels struct {
	View       string `json:"view"`
	Add        string `json:"add"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	BulkDelete string `json:"bulkDelete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type PricePlanColumnLabels2 struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Duration    string `json:"duration"`
	Location    string `json:"location"`
	Plan        string `json:"plan"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
	Actions     string `json:"actions"`
}

type PricePlanEmptyLabels struct {
	Title       string `json:"title"`
	Message     string `json:"message"`
	Description string `json:"description"`
	ActionLabel string `json:"actionLabel"`
}

type PricePlanActionLabels struct {
	CreateSuccess string `json:"createSuccess"`
	CreateError   string `json:"createError"`
	UpdateSuccess string `json:"updateSuccess"`
	UpdateError   string `json:"updateError"`
	DeleteSuccess string `json:"deleteSuccess"`
	DeleteError   string `json:"deleteError"`
}

type PricePlanBulkLabels struct {
	DeleteTitle   string `json:"deleteTitle"`
	DeleteMessage string `json:"deleteMessage"`
	StatusTitle   string `json:"statusTitle"`
	StatusMessage string `json:"statusMessage"`
}

type PricePlanDetailLabels2 struct {
	Title          string `json:"title"`
	InfoTab        string `json:"infoTab"`
	AttachmentsTab string `json:"attachmentsTab"`
	AuditTab       string `json:"auditTab"`
	ProductsTab    string `json:"productsTab"`
}

type PricePlanTabLabels2 struct {
	Info        string `json:"info"`
	Products    string `json:"products"`
	Attachments string `json:"attachments"`
	Audit       string `json:"audit"`
}

type PricePlanConfirmLabels struct {
	DeleteTitle       string `json:"deleteTitle"`
	DeleteMessage     string `json:"deleteMessage"`
	DeactivateTitle   string `json:"deactivateTitle"`
	DeactivateMessage string `json:"deactivateMessage"`
}

type PricePlanErrorLabels struct {
	NotFound     string `json:"notFound"`
	LoadFailed   string `json:"loadFailed"`
	Unauthorized string `json:"unauthorized"`
	CreateFailed string `json:"createFailed"`
	UpdateFailed string `json:"updateFailed"`
	DeleteFailed string `json:"deleteFailed"`
}

// DefaultPricePlanLabels returns PricePlanLabels with sensible English defaults.
func DefaultPricePlanLabels() PricePlanLabels {
	return PricePlanLabels{
		Page: PricePlanPageLabels{
			Title:         "Rate Cards",
			Subtitle:      "Manage your rate cards",
			ActiveTitle:   "Active Rate Cards",
			InactiveTitle: "Inactive Rate Cards",
		},
		Buttons: PricePlanButtonLabels{
			View:       "View",
			Add:        "Add Rate Card",
			Edit:       "Edit Rate Card",
			Delete:     "Delete Rate Card",
			BulkDelete: "Delete Rate Cards",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: PricePlanColumnLabels2{
			Name:        "Name",
			Description: "Description",
			Amount:      "Amount",
			Currency:    "Currency",
			Duration:    "Duration",
			Location:    "Location",
			Plan:        "Plan",
			Status:      "Status",
			DateCreated: "Date Created",
			Actions:     "Actions",
		},
		Empty: PricePlanEmptyLabels{
			Title:       "No Rate Cards",
			Message:     "No rate cards to display.",
			Description: "Add a rate card to define pricing for your plans.",
			ActionLabel: "Add Rate Card",
		},
		Form: PricePlanFormLabels{
			Name:                "Price Plan Name",
			NamePlaceholder:     "Enter price plan name",
			Description:         "Description",
			DescPlaceholder:     "Enter description...",
			Amount:              "Amount",
			AmountPlaceholder:   "0.00",
			Currency:            "Currency",
			CurrencyPlaceholder: "e.g. PHP",
			DurationValue:       "Duration",
			DurationUnit:        "Unit",
			Location:            "Location",
			LocationPlaceholder: "Select a location...",
			SelectLocation:      "— No location (all locations) —",
			Active:              "Active",
		},
		Actions: PricePlanActionLabels{
			CreateSuccess: "Rate card created successfully.",
			CreateError:   "Failed to create rate card.",
			UpdateSuccess: "Rate card updated successfully.",
			UpdateError:   "Failed to update rate card.",
			DeleteSuccess: "Rate card deleted successfully.",
			DeleteError:   "Failed to delete rate card.",
		},
		Bulk: PricePlanBulkLabels{
			DeleteTitle:   "Delete Rate Cards",
			DeleteMessage: "Are you sure you want to delete the selected rate cards?",
			StatusTitle:   "Update Status",
			StatusMessage: "Are you sure you want to update the status of the selected rate cards?",
		},
		Detail: PricePlanDetailLabels2{
			Title:          "Rate Card Details",
			InfoTab:        "Information",
			AttachmentsTab: "Attachments",
			AuditTab:       "Audit Trail",
			ProductsTab:    "Products",
		},
		Tabs: PricePlanTabLabels2{
			Info:        "Information",
			Products:    "Products",
			Attachments: "Attachments",
			Audit:       "Audit Trail",
		},
		Confirm: PricePlanConfirmLabels{
			DeleteTitle:       "Delete Rate Card",
			DeleteMessage:     "Are you sure you want to delete this rate card? This action cannot be undone.",
			DeactivateTitle:   "Deactivate Rate Card",
			DeactivateMessage: "Are you sure you want to deactivate this rate card?",
		},
		Errors: PricePlanErrorLabels{
			NotFound:     "Rate card not found.",
			LoadFailed:   "Failed to load rate cards.",
			Unauthorized: "You do not have permission to access this resource.",
			CreateFailed: "Failed to create rate card.",
			UpdateFailed: "Failed to update rate card.",
			DeleteFailed: "Failed to delete rate card.",
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
			Customer:                  "Customer",
			CustomerPlaceholder:       "Select customer...",
			Plan:                      "Plan",
			PlanPlaceholder:           "Select plan...",
			StartDate:                 "Start Date",
			EndDate:                   "End Date",
			Active:                    "Active",
			Notes:                     "Notes",
			NotesPlaceholder:          "Enter notes...",
			CustomerSearchPlaceholder: "Search customers...",
			PlanSearchPlaceholder:     "Search plans...",
			CustomerNoResults:         "No customers found",
			PlanNoResults:             "No plans found",
			Code:                      "Code",
			CodePlaceholder:           "e.g. A3K7PXR",
		},
		Actions: SubscriptionActionLabels{
			View:   "View Subscription",
			Edit:   "Edit Subscription",
			Cancel: "Cancel Subscription",
		},
		Detail: SubscriptionDetailLabels{
			PageTitle:            "Subscription Details",
			Customer:             "Customer",
			Plan:                 "Plan",
			StartDate:            "Start Date",
			EndDate:              "End Date",
			Status:               "Status",
			CreatedDate:          "Created",
			ModifiedDate:         "Last Modified",
			AuditTrailComingSoon: "Audit trail coming soon.",
			AuditTrailDesc:       "Audit trail for subscription changes is coming soon.",
		},
		Tabs: SubscriptionTabLabels{
			Info:        "Information",
			History:     "History",
			Attachments: "Attachments",
			AuditTrail:  "Audit Trail",
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
		Export:             common.Table.Export,
		DensityLabel:       common.Table.Density.Title,
		DensityDense:       common.Table.Density.Dense,
		DensityDefault:     common.Table.Density.Default,
		DensityComfortable: common.Table.Density.Comfortable,
		DensityCompact:     common.Table.Density.Compact,
		EntriesPerPage:     common.Table.EntriesLabel,
		Show:               common.Table.Show,
		Entries:            common.Table.Entries,
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
