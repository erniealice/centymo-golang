package centymo

import (
	"strings"

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
	TrackingMode TrackingModeLabels          `json:"trackingMode"`
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

	// Field-level info text surfaced via an info button beside each label.
	ProductInfo       string `json:"productInfo"`
	SKUInfo           string `json:"skuInfo"`
	OnHandInfo        string `json:"onHandInfo"`
	ReservedInfo      string `json:"reservedInfo"`
	ReorderLevelInfo  string `json:"reorderLevelInfo"`
	UnitOfMeasureInfo string `json:"unitOfMeasureInfo"`
	NotesInfo         string `json:"notesInfo"`
	ActiveInfo        string `json:"activeInfo"`
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

	// Field-level info text surfaced via an info button beside each label.
	SerialNumberInfo  string `json:"serialNumberInfo"`
	IMEIInfo          string `json:"imeiInfo"`
	StatusInfo        string `json:"statusInfo"`
	WarrantyStartInfo string `json:"warrantyStartInfo"`
	WarrantyEndInfo   string `json:"warrantyEndInfo"`
	PurchaseOrderInfo string `json:"purchaseOrderInfo"`
	SoldReferenceInfo string `json:"soldReferenceInfo"`
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

	// Field-level info text surfaced via an info button beside each label.
	TypeInfo      string `json:"typeInfo"`
	QuantityInfo  string `json:"quantityInfo"`
	DateInfo      string `json:"dateInfo"`
	ReferenceInfo string `json:"referenceInfo"`
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

	// Field-level info text surfaced via an info button beside each label.
	MethodInfo       string `json:"methodInfo"`
	CostBasisInfo    string `json:"costBasisInfo"`
	SalvageValueInfo string `json:"salvageValueInfo"`
	UsefulLifeInfo   string `json:"usefulLifeInfo"`
	StartDateInfo    string `json:"startDateInfo"`
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

	// Field-level info text for the payment drawer form.
	PaymentMethodInfo   string `json:"paymentMethodInfo"`
	AmountInfo          string `json:"amountInfo"`
	CurrencyInfo        string `json:"currencyInfo"`
	ReferenceNumberInfo string `json:"referenceNumberInfo"`
	ReceivedByInfo      string `json:"receivedByInfo"`
	RoleInfo            string `json:"roleInfo"`
	NotesInfo           string `json:"notesInfo"`
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

	// Field-level info text for the line-item drawer form.
	ProductInfo     string `json:"productInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	QuantityInfo    string `json:"quantityInfo"`
	UnitPriceInfo   string `json:"unitPriceInfo"`
	CostPriceInfo   string `json:"costPriceInfo"`
	DiscountInfo    string `json:"discountInfo"`
	NotesInfo       string `json:"notesInfo"`

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
	Page         ProductPageLabels       `json:"page"`
	Buttons      ProductButtonLabels     `json:"buttons"`
	Columns      ProductColumnLabels     `json:"columns"`
	Empty        ProductEmptyLabels      `json:"empty"`
	Form         ProductFormLabels       `json:"form"`
	Actions      ProductActionLabels     `json:"actions"`
	Bulk         ProductBulkLabels       `json:"bulkActions"`
	Tabs         ProductTabLabels        `json:"tabs"`
	Detail       ProductDetailLabels     `json:"detail"`
	Status       ProductStatusLabels     `json:"status"`
	Variant      ProductVariantLabels    `json:"variant"`
	Attribute    ProductAttributeLabels  `json:"attribute"`
	Options      ProductOptionLabels     `json:"options"`
	Confirm      ProductConfirmLabels    `json:"confirm"`
	Errors       ProductErrorLabels      `json:"errors"`
	Breadcrumb   ProductBreadcrumbLabels `json:"breadcrumb"`
	// Four-axis product taxonomy enum labels — loaded from lyngua
	// product.json "productKind"/"deliveryMode"/"trackingMode" blocks.
	// Wired here so the drawer-form select uses the exact tier-cascaded
	// display string for each enum value without hardcoding in Go.
	ProductKind  ProductKindLabels  `json:"productKind"`
	DeliveryMode DeliveryModeLabels `json:"deliveryMode"`
	TrackingMode TrackingModeLabels `json:"trackingMode"`
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
	Required                 string `json:"required"`
	Option                   string `json:"option"`
	SelectAttribute          string `json:"selectAttribute"`
	AllAttributesAssigned    string `json:"allAttributesAssigned"`
	OptionNeedsValuesAlert   string `json:"optionNeedsValuesAlert"`

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
	Name                    string `json:"name"`
	NamePlaceholder         string `json:"namePlaceholder"`
	Code                    string `json:"code"`
	CodePlaceholder         string `json:"codePlaceholder"`
	DataType                string `json:"dataType"`
	SortOrder               string `json:"sortOrder"`
	MinValue                string `json:"minValue"`
	MaxValue                string `json:"maxValue"`
	Active                  string `json:"active"`
	Required                string `json:"required"`
	RequiredCaution         string `json:"requiredCaution"`
	Description             string `json:"description"`
	DescriptionPlaceholder  string `json:"descriptionPlaceholder"`
	DescriptionEmpty        string `json:"descriptionEmpty"`

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
	TextList   string `json:"textList"`
	NumberRange string `json:"numberRange"`
	ColorList  string `json:"colorList"`
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

	// Field-level info text surfaced via an info button beside each label.
	CodeInfo        string `json:"codeInfo"`
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
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
	NoPermission     string `json:"noPermission"`
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

	// Field-level info text surfaced via an info button beside each label.
	NameInfo            string `json:"nameInfo"`
	ExpenditureTypeInfo string `json:"expenditureTypeInfo"`
	CategoryInfo        string `json:"categoryInfo"`
	DateInfo            string `json:"dateInfo"`
	AmountInfo          string `json:"amountInfo"`
	CurrencyInfo        string `json:"currencyInfo"`
	ReferenceNumberInfo string `json:"referenceNumberInfo"`
	SupplierInfo        string `json:"supplierInfo"`
	NotesInfo           string `json:"notesInfo"`
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
	// SPS P10 — Recognition + Accrual tabs on expenditure detail
	TabRecognition          string `json:"tabRecognition"`
	TabAccrual              string `json:"tabAccrual"`
	RecognitionEmptyTitle   string `json:"recognitionEmptyTitle"`
	RecognitionEmptyMessage string `json:"recognitionEmptyMessage"`
	RecognitionRecognizeCTA string `json:"recognitionRecognizeCta"`
	AccrualEmptyTitle       string `json:"accrualEmptyTitle"`
	AccrualEmptyMessage     string `json:"accrualEmptyMessage"`
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

// PurchaseOrderErrorLabels holds error messages for the purchase order action handlers.
type PurchaseOrderErrorLabels struct {
	NoPermission string `json:"noPermission"`
}

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
	Errors      PurchaseOrderErrorLabels     `json:"errors"`
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

	// Field-level info text surfaced via an info button beside each label.
	PONumberInfo          string `json:"poNumberInfo"`
	POTypeInfo            string `json:"poTypeInfo"`
	SupplierInfo          string `json:"supplierInfo"`
	OrderDateInfo         string `json:"orderDateInfo"`
	ExpectedDeliveryInfo  string `json:"expectedDeliveryInfo"`
	CurrencyInfo          string `json:"currencyInfo"`
	PaymentTermsInfo      string `json:"paymentTermsInfo"`
	ShippingTermsInfo     string `json:"shippingTermsInfo"`
	ReferenceNumberInfo   string `json:"referenceNumberInfo"`
	NotesInfo             string `json:"notesInfo"`
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

	// Field-level info text surfaced via an info button beside each label.
	ReferenceInfo    string `json:"referenceInfo"`
	CustomerInfo     string `json:"customerInfo"`
	AmountInfo       string `json:"amountInfo"`
	CurrencyInfo     string `json:"currencyInfo"`
	PaymentMethodInfo string `json:"paymentMethodInfo"`
	DateInfo         string `json:"dateInfo"`
	StatusInfo       string `json:"statusInfo"`
	NotesInfo        string `json:"notesInfo"`
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
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			ReferenceInfo:     "Unique reference number for this collection record.",
			CustomerInfo:      "Name of the customer or payer.",
			AmountInfo:        "Total amount collected (in centavos; displayed as amount ÷ 100).",
			CurrencyInfo:      "Currency of the collected amount.",
			PaymentMethodInfo: "How the payment was received.",
			DateInfo:          "Date the payment was collected.",
			StatusInfo:        "Current state of this collection record.",
			NotesInfo:         "Internal remarks — not shown on customer-facing documents.",
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

	// Field-level info text surfaced via an info button beside each label.
	ReferenceInfo    string `json:"referenceInfo"`
	DateInfo         string `json:"dateInfo"`
	PayeeInfo        string `json:"payeeInfo"`
	AmountInfo       string `json:"amountInfo"`
	CurrencyInfo     string `json:"currencyInfo"`
	PaymentMethodInfo string `json:"paymentMethodInfo"`
	StatusInfo       string `json:"statusInfo"`
	CategoryInfo     string `json:"categoryInfo"`
	ApprovedByInfo   string `json:"approvedByInfo"`
	NotesInfo        string `json:"notesInfo"`
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
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			ReferenceInfo:     "Unique reference number for this disbursement.",
			DateInfo:          "Date the disbursement was issued.",
			PayeeInfo:         "Name of the recipient (supplier, payroll, etc.).",
			AmountInfo:        "Total amount disbursed (in centavos; displayed as amount ÷ 100).",
			CurrencyInfo:      "Currency of the disbursed amount.",
			PaymentMethodInfo: "How the payment was made.",
			StatusInfo:        "Current state of this disbursement.",
			CategoryInfo:      "Type of disbursement for categorisation and reporting.",
			ApprovedByInfo:    "Name of the person who authorised this disbursement.",
			NotesInfo:         "Internal remarks — not shown on supplier-facing documents.",
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

// PlanFilterLabels holds translatable labels for the scope filter chip on the
// plan list page (§6.1 of the 2026-04-27 plan-client-scope plan).
type PlanFilterLabels struct {
	ScopeChipLabel string `json:"scopeChipLabel"`
	ScopeMaster    string `json:"scopeMaster"`
	ScopeClient    string `json:"scopeClient"`
	ScopeAll       string `json:"scopeAll"`
}

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
	Filters         PlanFilterLabels       `json:"filters"`
}

// ProductPlanFormLabels holds translatable labels for the ProductPlan add/edit form within a plan.
type ProductPlanFormLabels struct {
	Product            string                  `json:"product"`
	ProductPlaceholder string                  `json:"productPlaceholder"`
	SelectProduct      string                  `json:"selectProduct"`
	Active             string                  `json:"active"`
	ProductKindLabel   string                  `json:"productKindLabel"`
	ProductKind        ProductKindOptionLabels `json:"productKind"`

	// Model D — variant picker on the ProductPlan drawer form
	VariantSelectLabel       string `json:"variantSelectLabel"`
	VariantSelectPlaceholder string `json:"variantSelectPlaceholder"`
	VariantSelectInfo        string `json:"variantSelectInfo"`
}

// ProductKindOptionLabels provides translated labels for each product_kind
// enum value, used to build the kind selector on the add/edit drawer AND
// to map product_kind values to display labels in table cells.
type ProductKindOptionLabels struct {
	Service        string `json:"service"`
	StockedGood    string `json:"stockedGood"`
	NonStockedGood string `json:"nonStockedGood"`
	Consumable     string `json:"consumable"`
}

// Label returns the translated label for a product_kind value
// ("service" | "stocked_good" | "non_stocked_good" | "consumable").
// Unknown values round-trip through as-is so callers always get a string.
func (k ProductKindOptionLabels) Label(kind string) string {
	switch kind {
	case "service":
		return k.Service
	case "stocked_good":
		return k.StockedGood
	case "non_stocked_good":
		return k.NonStockedGood
	case "consumable":
		return k.Consumable
	}
	return kind
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
	CurrencyPHP         string `json:"currencyPHP"`
	CurrencyUSD         string `json:"currencyUSD"`
	DurationValue       string `json:"durationValue"`
	DurationUnit        string `json:"durationUnit"`
	Schedule            string `json:"schedule"`
	SchedulePlaceholder string `json:"schedulePlaceholder"`
	ScheduleSearch      string `json:"scheduleSearch"`
	Location            string `json:"location"`
	LocationPlaceholder string `json:"locationPlaceholder"`
	LocationHintPrefix  string `json:"locationHintPrefix"`
	SelectLocation      string `json:"selectLocation"`
	Active              string `json:"active"`
	PlanLabel           string `json:"planLabel"`
	PlanPlaceholder     string `json:"planPlaceholder"`
	PlanSearch          string `json:"planSearch"`

	// Wave 2 — new billing semantics fields (from lyngua price_plan.json → price_plan.form)
	SectionBasic               string `json:"sectionBasic"`
	SectionPricing             string `json:"sectionPricing"`
	BillingKindLabel           string `json:"billingKindLabel"`
	BillingKindOneTime         string `json:"billingKindOneTime"`
	BillingKindRecurring       string `json:"billingKindRecurring"`
	BillingKindContract        string `json:"billingKindContract"`
	BillingKindMilestone       string `json:"billingKindMilestone"`
	BillingKindAdHoc           string `json:"billingKindAdHoc"`
	// Per-option hint copy surfaced inline below the billing_kind select as the
	// operator picks. Matches the multi-vertical convention — general/ tier ships
	// neutral phrasing, professional/ overrides with engagement vocabulary.
	BillingKindOneTimeHint   string `json:"billingKindOneTimeHint"`
	BillingKindRecurringHint string `json:"billingKindRecurringHint"`
	BillingKindContractHint  string `json:"billingKindContractHint"`
	BillingKindMilestoneHint string `json:"billingKindMilestoneHint"`
	BillingKindAdHocHint     string `json:"billingKindAdHocHint"`
	AmountBasisLabel           string `json:"amountBasisLabel"`
	AmountBasisPerCycle        string `json:"amountBasisPerCycle"`
	AmountBasisTotalPackage    string `json:"amountBasisTotalPackage"`
	AmountBasisDerivedFromLines string `json:"amountBasisDerivedFromLines"`
	AmountBasisPerOccurrence   string `json:"amountBasisPerOccurrence"`
	// Per-option hint copy for amount_basis (mirrors billing_kind pattern).
	AmountBasisPerCycleHint         string `json:"amountBasisPerCycleHint"`
	AmountBasisTotalPackageHint     string `json:"amountBasisTotalPackageHint"`
	AmountBasisDerivedFromLinesHint string `json:"amountBasisDerivedFromLinesHint"`
	AmountBasisPerOccurrenceHint    string `json:"amountBasisPerOccurrenceHint"`
	EntitledOccurrencesLabel       string `json:"entitledOccurrencesLabel"`
	EntitledOccurrencesPlaceholder string `json:"entitledOccurrencesPlaceholder"`
	EntitledOccurrencesInfo        string `json:"entitledOccurrencesInfo"`
	BillingCycleLabel       string `json:"billingCycleLabel"`
	BillingCyclePlaceholder string `json:"billingCyclePlaceholder"`
	TermLabel               string `json:"termLabel"`
	TermPlaceholder         string `json:"termPlaceholder"`
	TermOpenEndedHelp       string `json:"termOpenEndedHelp"`

	// Field-level info text surfaced via an info button beside each label.
	PlanInfo         string `json:"planInfo"`
	ScheduleInfo     string `json:"scheduleInfo"`
	NameInfo         string `json:"nameInfo"`
	DescriptionInfo  string `json:"descriptionInfo"`
	BillingKindInfo  string `json:"billingKindInfo"`
	AmountBasisInfo  string `json:"amountBasisInfo"`
	AmountInfo       string `json:"amountInfo"`
	CurrencyInfo     string `json:"currencyInfo"`
	BillingCycleInfo string `json:"billingCycleInfo"`
	TermInfo         string `json:"termInfo"`
	ActiveInfo       string `json:"activeInfo"`

	// 2026-04-27 plan-client-scope plan §6.7 — info banner shown above the
	// PricePlan add/edit form when its parent PriceSchedule is client-scoped.
	// Templated via Go's text/template ({{.ClientName}}).
	ParentScheduleClientNotice string `json:"parentScheduleClientNotice"`

	// 2026-04-27 plan-client-scope plan §6.7 — tooltip surfaced beside the
	// readonly Schedule label when the PricePlan's parent Plan is
	// client-scoped (the schedule field is locked to the resolved/derived
	// client schedule). Templated via Go's text/template ({{.ClientName}}).
	ScheduleLockedTooltip string `json:"scheduleLockedTooltip"`
	// 2026-04-28 — info-row hints rendered beneath the readonly Schedule
	// label so the operator knows what happens on save:
	//   ScheduleAutoCreateHint — no client rate card exists yet; one will be
	//     created with this client's name + the lyngua suffix.
	//   ScheduleAutoReuseHint  — an existing client rate card was found; the
	//     new price plan will attach to it.
	// Both templated with {{.ClientName}}.
	ScheduleAutoCreateHint string `json:"scheduleAutoCreateHint"`
	ScheduleAutoReuseHint  string `json:"scheduleAutoReuseHint"`

	// 2026-04-30 cyclic-subscription-jobs plan §9.4 — client-side block
	// surfaced as a tooltip on the disabled MILESTONE option in the
	// billing_kind dropdown when the parent Plan is cyclic.
	MilestoneCyclicBlock string `json:"milestoneCyclicBlock"`

	// 2026-05-01 ad-hoc-subscription-billing plan §6 — client-side guards
	// surfaced as drawer warnings / tooltips on the disabled options. The
	// server enforces the same rules in validate_ad_hoc.go.
	AdHocPoolNoTemplate          string `json:"adHocPoolNoTemplate"`
	AdHocPerCallNoTemplate       string `json:"adHocPerCallNoTemplate"`
	AdHocNoEntitlement           string `json:"adHocNoEntitlement"`
	AdHocBillingCycleNotAllowed  string `json:"adHocBillingCycleNotAllowed"`
	AdHocVisitsPerCycleNotAllowed string `json:"adHocVisitsPerCycleNotAllowed"`
}

// ---------------------------------------------------------------------------
// ProductPricePlan labels
// ---------------------------------------------------------------------------

// ProductPricePlanLabels holds all labels for the ProductPricePlan drawer form.
// Wave 2 addition: billing treatment + product/price/currency/date fields.
type ProductPricePlanLabels struct {
	Form ProductPricePlanFormLabels `json:"form"`
}

// ProductPricePlanFormLabels holds translatable labels for the ProductPricePlan
// add/edit drawer form. Keys match lyngua product_price_plan.json → product_price_plan.form.
type ProductPricePlanFormLabels struct {
	BillingTreatmentLabel              string `json:"billingTreatmentLabel"`
	BillingTreatmentRecurring          string `json:"billingTreatmentRecurring"`
	BillingTreatmentRecurringHelp      string `json:"billingTreatmentRecurringHelp"`
	BillingTreatmentOneTimeInitial     string `json:"billingTreatmentOneTimeInitial"`
	BillingTreatmentOneTimeInitialHelp string `json:"billingTreatmentOneTimeInitialHelp"`
	BillingTreatmentUsageBased         string `json:"billingTreatmentUsageBased"`
	BillingTreatmentUsageBasedHelp     string `json:"billingTreatmentUsageBasedHelp"`
	ProductLabel                       string `json:"productLabel"`
	ProductPlaceholder                 string `json:"productPlaceholder"`
	PriceLabel                         string `json:"priceLabel"`
	PricePlaceholder                   string `json:"pricePlaceholder"`
	CurrencyLabel                      string `json:"currencyLabel"`
	CurrencyPlaceholder                string `json:"currencyPlaceholder"`
	CurrencyPHP                        string `json:"currencyPHP"`
	CurrencyUSD                        string `json:"currencyUSD"`
	DateStartLabel                     string `json:"dateStartLabel"`
	DateEndLabel                       string `json:"dateEndLabel"`

	// Field-level info text surfaced via an info button beside each label.
	ProductInfo          string `json:"productInfo"`
	PriceInfo            string `json:"priceInfo"`
	CurrencyInfo         string `json:"currencyInfo"`
	BillingTreatmentInfo string `json:"billingTreatmentInfo"`
	DateStartInfo        string `json:"dateStartInfo"`
	DateEndInfo          string `json:"dateEndInfo"`

	// Model D — catalog line picker (replaces product_id with product_plan_id)
	CatalogLineLabel       string `json:"catalogLineLabel"`
	CatalogLinePlaceholder string `json:"catalogLinePlaceholder"`
	CatalogLineInfo        string `json:"catalogLineInfo"`

	// 2026-04-29 milestone-billing plan §5 / Phase D — milestone (job
	// template phase) select. Surfaced when the parent PricePlan has
	// billing_kind = MILESTONE; an empty selection falls through to the
	// first event for the milestone plan.
	MilestonePhaseLabel       string `json:"milestonePhaseLabel"`
	MilestonePhaseFallthrough string `json:"milestonePhaseFallthrough"`
	MilestonePhaseBillable    string `json:"milestonePhaseBillable"`
}

// DefaultProductPricePlanLabels returns ProductPricePlanLabels with sensible English defaults.
func DefaultProductPricePlanLabels() ProductPricePlanLabels {
	return ProductPricePlanLabels{
		Form: ProductPricePlanFormLabels{
			BillingTreatmentLabel:              "Billing treatment",
			BillingTreatmentRecurring:          "Every cycle",
			BillingTreatmentRecurringHelp:      "Charge this line every billing cycle",
			BillingTreatmentOneTimeInitial:     "First cycle only",
			BillingTreatmentOneTimeInitialHelp: "Charge once on the first invoice (setup fees, welcome gifts)",
			BillingTreatmentUsageBased:         "On use",
			BillingTreatmentUsageBasedHelp:     "Charge when consumed or performed",
			ProductLabel:                       "Product",
			ProductPlaceholder:                 "Select a product",
			PriceLabel:                         "Price",
			PricePlaceholder:                   "0.00",
			CurrencyLabel:                      "Currency",
			CurrencyPlaceholder:                "e.g. PHP",
			CurrencyPHP:                        "PHP (₱)",
			CurrencyUSD:                        "USD ($)",
			DateStartLabel:                     "Effective from",
			DateEndLabel:                       "Effective until",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			ProductInfo:          "The product this price applies to.",
			PriceInfo:            "Price in centavos. Displayed as amount ÷ 100.",
			CurrencyInfo:         "Currency applied to this product price.",
			BillingTreatmentInfo: "Every cycle = charged each billing cycle. First cycle only = setup fee. On use = charged when consumed.",
			DateStartInfo:        "Date from which this product price is effective.",
			DateEndInfo:          "Last date this product price is effective. Leave empty for no end date.",
			// Model D — catalog line picker defaults
			CatalogLineLabel:       "Catalog line",
			CatalogLinePlaceholder: "Select a line from the plan's catalog",
			CatalogLineInfo:        "Prices the chosen catalog line from the parent plan. If the line has a variant, that variant is priced.",
			// 2026-04-29 milestone-billing plan §5 — milestone phase select.
			MilestonePhaseLabel:       "Milestone phase",
			MilestonePhaseFallthrough: "Falls through to first event",
			MilestonePhaseBillable:    "billable",
		},
	}
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
	EditPricePlan string `json:"editPricePlan"`
	AddProduct    string `json:"addProduct"`
}

type PlanColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Interval    string `json:"interval"`
	Price       string `json:"price"`
	Status      string `json:"status"`
	Product     string `json:"product"`
	PricePlan     string `json:"pricePlan"`
	PriceSchedule string `json:"priceSchedule"`
	Duration      string `json:"duration"`
	Location    string `json:"location"`
	ItemType    string `json:"itemType"`
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

	// 2026-04-27 plan-client-scope plan §7 — surfaced when an operator tries
	// to change a Plan's client_id while one of its PricePlans is attached
	// to an active subscription. Hard block; no force-override.
	ClientScopeLocked string `json:"clientScopeLocked"`
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

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	ActiveInfo      string `json:"activeInfo"`

	// Client-scope fields (2026-04-27 plan-client-scope plan §7).
	// Set on the Plan add/edit drawer Client picker.
	ClientLabel              string `json:"clientLabel"`
	ClientHelp               string `json:"clientHelp"`
	ClientPlaceholder        string `json:"clientPlaceholder"`
	ClientSearchPlaceholder  string `json:"clientSearchPlaceholder"`
	ClientNoResults          string `json:"clientNoResults"`
	ClientLockedTooltip      string `json:"clientLockedTooltip"`
	ClientForLabel           string `json:"clientForLabel"` // "For {{.ClientName}}" — read-only badge in client-context entry-point
	ClientInfo               string `json:"clientInfo"`

	// JobTemplate select (2026-04-29 auto-spawn-jobs-from-subscription plan §5
	// — Plan.job_template_id assignment from the drawer). Empty value =
	// advisory-only plan; spawn use case skips silently.
	JobTemplate     string `json:"jobTemplate"`
	JobTemplateNone string `json:"jobTemplateNone"`
	JobTemplateHint string `json:"jobTemplateHint"`

	// 2026-04-30 cyclic-subscription-jobs plan §9.3 — visits_per_cycle field.
	// Number of cycle Job instances spawned per billing cycle (default 1).
	VisitsPerCycleLabel       string `json:"visitsPerCycleLabel"`
	VisitsPerCyclePlaceholder string `json:"visitsPerCyclePlaceholder"`
	VisitsPerCycleHint        string `json:"visitsPerCycleHint"`
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
	Info          string `json:"info"`
	Products      string `json:"products"`
	ProductsSlug  string `json:"productsSlug"`
	PricePlan     string `json:"pricePlan"`
	PricePlanSlug string `json:"pricePlanSlug"`
	Attachments   string `json:"attachments"`
	AuditTrail    string `json:"auditTrail"`
}

// ResolveTabSlug returns the URL slug for a canonical tab key. The "products"
// and "pricePlan" tabs can be re-slugged per tier (e.g. professional ships
// "items" and "package-prices"); other tabs round-trip through as-is.
func (t PlanTabLabels) ResolveTabSlug(canonical string) string {
	switch canonical {
	case "products":
		if s := strings.TrimSpace(t.ProductsSlug); s != "" {
			return s
		}
	case "pricePlan":
		if s := strings.TrimSpace(t.PricePlanSlug); s != "" {
			return s
		}
	}
	return canonical
}

// CanonicalizeTab maps an incoming URL tab slug back to its canonical key so
// internal template lookups and equality checks stay tier-agnostic.
func (t PlanTabLabels) CanonicalizeTab(slug string) string {
	if slug == "" {
		return ""
	}
	if s := strings.TrimSpace(t.ProductsSlug); s != "" && slug == s {
		return "products"
	}
	if s := strings.TrimSpace(t.PricePlanSlug); s != "" && slug == s {
		return "pricePlan"
	}
	return slug
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
	Page       SubscriptionPageLabels       `json:"page"`
	Buttons    SubscriptionButtonLabels     `json:"buttons"`
	Columns    SubscriptionColumnLabels     `json:"columns"`
	Empty      SubscriptionEmptyLabels      `json:"empty"`
	Form       SubscriptionFormLabels       `json:"form"`
	Actions    SubscriptionActionLabels     `json:"actions"`
	Bulk       SubscriptionBulkLabels       `json:"bulkActions"`
	Status     SubscriptionStatusLabels     `json:"status"`
	Detail     SubscriptionDetailLabels     `json:"detail"`
	Tabs       SubscriptionTabLabels        `json:"tabs"`
	Invoices   SubscriptionInvoicesLabels   `json:"invoices"`
	Recognize  SubscriptionRecognizeLabels  `json:"recognize"`
	Milestone  SubscriptionMilestoneLabels  `json:"milestone"`
	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / §9 — Operations
	// tab on the subscription detail page + retroactive spawn drawer copy.
	Operations SubscriptionOperationsLabels `json:"operations"`
	Spawn      SubscriptionSpawnLabels      `json:"spawn"`
	// 2026-04-30 cyclic-subscription-jobs plan §9.2 / §21.3 — Backfill cycle
	// Jobs drawer + flat Jobs tab.
	Backfill SubscriptionBackfillLabels  `json:"backfill"`
	Jobs     SubscriptionJobsTabLabels   `json:"jobs"`
	Confirm  SubscriptionConfirmLabels   `json:"confirm"`
	Errors   SubscriptionErrorLabels     `json:"errors"`
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
	Name      string `json:"name"`
	Client    string `json:"client"`
	Customer  string `json:"customer"` // legacy alias; kept for backward compat with old translations
	Plan      string `json:"plan"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Status    string `json:"status"`
}

type SubscriptionEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type SubscriptionActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Cancel     string `json:"cancel"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`

	// 2026-04-27 plan-client-scope plan §6.5 / §7 — CTA copy on the
	// subscription detail's Package tab. Templated via {{.ClientName}}.
	CustomizePackage string `json:"customizePackage"`
}

type SubscriptionBulkLabels struct {
	Delete     string `json:"delete"`
	Activate   string `json:"bulkActivate"`
	Deactivate string `json:"bulkDeactivate"`
}

type SubscriptionStatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type SubscriptionErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
	NoPermission     string `json:"noPermission"`
	CannotDelete     string `json:"cannotDelete"`
	InUse            string `json:"inUse"`

	// 2026-04-27 plan-client-scope plan §3.3 / §7 — surfaced when the
	// subscription's selected price_plan belongs to a different client.
	PlanClientMismatch string `json:"planClientMismatch"`
	// Surfaced when the customize-package CTA fails (cross-package errors
	// from the espyna use case bubble up here as a generic fallback).
	CustomizeFailed string `json:"customizeFailed"`
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
	StartTime                 string `json:"startTime"`
	EndTime                   string `json:"endTime"`
	TimePlaceholder           string `json:"timePlaceholder"`
	Timezone                  string `json:"timezone"`
	Active                    string `json:"active"`
	Notes                     string `json:"notes"`
	NotesPlaceholder          string `json:"notesPlaceholder"`
	CustomerSearchPlaceholder string `json:"customerSearchPlaceholder"`
	PlanSearchPlaceholder     string `json:"planSearchPlaceholder"`
	CustomerNoResults         string `json:"customerNoResults"`
	PlanNoResults             string `json:"planNoResults"`
	Code                      string `json:"code"`
	CodePlaceholder           string `json:"codePlaceholder"`

	// Field-level info text surfaced via an info button beside each label.
	CustomerInfo  string `json:"customerInfo"`
	PlanInfo      string `json:"planInfo"`
	CodeInfo      string `json:"codeInfo"`
	StartDateInfo string `json:"startDateInfo"`
	EndDateInfo   string `json:"endDateInfo"`
	StartTimeInfo string `json:"startTimeInfo"`
	EndTimeInfo   string `json:"endTimeInfo"`
	NotesInfo     string `json:"notesInfo"`

	// 2026-04-27 plan-client-scope plan §5.1 / §7 — group headers in the
	// grouped Plan / PricePlan auto-complete picker on the subscription
	// drawer. Templated via {{.ClientName}} for the per-client group.
	PlanGroupForClient string `json:"planGroupForClient"`
	PlanGroupGeneral   string `json:"planGroupGeneral"`

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.1 / §9 — Spawn
	// Jobs toggle section on the subscription create drawer.
	SpawnJobsSectionTitle string `json:"spawnJobsSectionTitle"`
	SpawnJobsToggle       string `json:"spawnJobsToggle"`
	SpawnJobsHelpText     string `json:"spawnJobsHelpText"`
	SpawnJobsSummary      string `json:"spawnJobsSummary"`
	SpawnJobsNone         string `json:"spawnJobsNone"`
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
	Info         string `json:"info"`
	Operations   string `json:"operations"`
	// 2026-04-30 cyclic-subscription-jobs plan §21.2 — flat Jobs tab.
	Jobs         string `json:"jobs"`
	Invoices     string `json:"invoices"`
	History      string `json:"history"`
	Attachments  string `json:"attachments"`
	AuditTrail   string `json:"auditTrail"`
	AuditHistory string `json:"auditHistory"`
}

type SubscriptionInvoicesLabels struct {
	Title        string `json:"title"`
	Empty        string `json:"empty"`
	ColumnCode   string `json:"columnCode"`
	ColumnDate   string `json:"columnDate"`
	ColumnAmount string `json:"columnAmount"`
	ColumnStatus string `json:"columnStatus"`

	// Recognize-revenue action surfaced as a primary action on the invoices
	// tab toolbar AND on the empty-state. No page-header button (per plan
	// §11.2 — tab-only).
	RecognizeAction   string `json:"recognizeAction"`
	RecognizeTitle    string `json:"recognizeTitle"`
	RecognizeSubtitle string `json:"recognizeSubtitle"`
}

// SubscriptionRecognizeLabels holds drawer-form labels for the
// "Recognize Revenue" flow. See plan §5 Phase E for the full table; the
// blocking-error keys (currencyMismatchError, idempotencyError) are renamed
// from their advisory counterparts since v1 surfaces them as hard blocks.
type SubscriptionRecognizeLabels struct {
	// Header / context section
	ContextSection string `json:"contextSection"`
	ClientLabel    string `json:"clientLabel"`
	PlanLabel      string `json:"planLabel"`
	QuantityLabel  string `json:"quantityLabel"`

	// Period section
	PeriodSection string `json:"periodSection"`
	PeriodStart   string `json:"periodStart"`
	PeriodEnd     string `json:"periodEnd"`
	RevenueDate   string `json:"revenueDate"`

	// Line items table
	LineItemsSection      string `json:"lineItemsSection"`
	ColumnDescription     string `json:"columnDescription"`
	ColumnUnitPrice       string `json:"columnUnitPrice"`
	ColumnQuantity        string `json:"columnQuantity"`
	ColumnLineTotal       string `json:"columnLineTotal"`
	ColumnTreatment       string `json:"columnTreatment"`
	TotalLabel            string `json:"totalLabel"`
	RemoveLine            string `json:"removeLine"`
	TreatmentRecurring    string `json:"treatmentRecurring"`
	TreatmentFirstCycle   string `json:"treatmentFirstCycle"`
	TreatmentUsageBased   string `json:"treatmentUsageBased"`
	TreatmentOneTime      string `json:"treatmentOneTime"`

	// Notes
	NotesLabel       string `json:"notesLabel"`
	NotesPlaceholder string `json:"notesPlaceholder"`

	// Footer buttons (v1 — single Generate button; "Save as Draft" is dropped
	// per plan Phase D refinement since both paths run the idempotency check.)
	Generate string `json:"generate"`
	Cancel   string `json:"cancel"`

	// Blocking error banners
	CurrencyMismatchError       string `json:"currencyMismatchError"`
	IdempotencyError            string `json:"idempotencyError"`
	IdempotencyExistingLink     string `json:"idempotencyExistingLink"`
	NoLinesError                string `json:"noLinesError"`
	CycleNotConfiguredWarning   string `json:"cycleNotConfiguredWarning"`
	UsageBasedSkippedNotice     string `json:"usageBasedSkippedNotice"`

	// 2026-04-27 plan-client-scope plan §7 — info notice on the recognize
	// drawer when the active subscription's PricePlan is client-scoped.
	// Templated via {{.ClientName}}.
	ClientCustomNotice string `json:"clientCustomNotice"`

	// 2026-04-29 milestone-billing plan §5 / Phase E — milestone-specific
	// drawer fields. Surfaced only when pricePlan.billing_kind = MILESTONE.
	MilestoneSelect            string `json:"milestoneSelect"`
	MilestoneSelectPlaceholder string `json:"milestoneSelectPlaceholder"`
	NoReadyMilestone           string `json:"noReadyMilestone"`
	MilestoneNotApplicable     string `json:"milestoneNotApplicable"`
	BillAmount                 string `json:"billAmount"`
	LeaveRemainderOpen         string `json:"leaveRemainderOpen"`
	CloseShort                 string `json:"closeShort"`
	PartialReason              string `json:"partialReason"`
	PartialReasonRequired      string `json:"partialReasonRequired"`
	OverBillingRejected        string `json:"overBillingRejected"`
}

// SubscriptionMilestoneLabels holds labels for the Subscription Package tab's
// Milestones section + the mark-ready / waive CTAs. Lyngua key:
// `subscription.milestone.*`. See milestone-billing plan §5.
type SubscriptionMilestoneLabels struct {
	Title            string `json:"title"`
	Subtitle         string `json:"subtitle"`
	MarkReady        string `json:"markReady"`
	Waive            string `json:"waive"`
	ViewInvoice      string `json:"viewInvoice"`
	StatusPending    string `json:"statusPending"`
	StatusReady      string `json:"statusReady"`
	StatusBilled     string `json:"statusBilled"`
	StatusWaived     string `json:"statusWaived"`
	StatusDeferred   string `json:"statusDeferred"`
	StatusCancelled  string `json:"statusCancelled"`
	TotalInvoiced    string `json:"totalInvoiced"`
	AmountFull       string `json:"amountFull"`
	AmountPartial    string `json:"amountPartial"`
}

// SubscriptionOperationsLabels holds labels for the Subscription detail's
// Operations tab. Lyngua key: `subscription.detail.operations.*`. See
// auto-spawn-jobs-from-subscription plan §5.2 / §9 and cyclic-subscription-jobs
// plan §9.1 (cycle accordion + backfill keys).
type SubscriptionOperationsLabels struct {
	Title         string `json:"title"`
	EmptyTitle    string `json:"emptyTitle"`
	EmptyMessage  string `json:"emptyMessage"`
	SpawnAction   string `json:"spawnAction"`
	RootJob       string `json:"rootJob"`
	ChildJob      string `json:"childJob"`
	PhaseSummary  string `json:"phaseSummary"`
	ViewJobLink   string `json:"viewJobLink"`

	// 2026-04-30 cyclic-subscription-jobs plan §9.1 — cycle accordion copy.
	EngagementHeading     string `json:"engagementHeading"`
	CycleHeading          string `json:"cycleHeading"`
	CyclePlaceholder      string `json:"cyclePlaceholder"`
	CycleSpawnNow         string `json:"cycleSpawnNow"`
	CycleStatusPending    string `json:"cycleStatusPending"`
	CycleStatusInProgress string `json:"cycleStatusInProgress"`
	CycleStatusCompleted  string `json:"cycleStatusCompleted"`
	CycleStatusOverdue    string `json:"cycleStatusOverdue"`
	CycleInvoiceLinked    string `json:"cycleInvoiceLinked"`
	CycleNoInvoice        string `json:"cycleNoInvoice"`
	CycleEmpty            string `json:"cycleEmpty"`
	BackfillBanner        string `json:"backfillBanner"`
	BackfillCta           string `json:"backfillCta"`

	// 2026-05-01 ad-hoc-subscription-billing plan §5.2 — Operations tab
	// AD_HOC mode keys. Vertical-neutral defaults ("usage", "occurrence")
	// with professional-tier overrides ("service call", "retainer", etc.).
	AdHocPoolHeading        string `json:"adHocPoolHeading"`
	AdHocPerCallHeading     string `json:"adHocPerCallHeading"`
	EntitlementUsed         string `json:"entitlementUsed"`
	EntitlementRemaining    string `json:"entitlementRemaining"`
	EntitlementExhausted    string `json:"entitlementExhausted"`
	RequestUsageCta         string `json:"requestUsageCta"`
	ExtendEntitlementCta    string `json:"extendEntitlementCta"`
	UsageRequestedDate      string `json:"usageRequestedDate"`
	UsageDeliveredDate      string `json:"usageDeliveredDate"`
	UsageOrdinalLabel       string `json:"usageOrdinalLabel"`
	UsageNotDelivered       string `json:"usageNotDelivered"`
	PoolInvoiceLink         string `json:"poolInvoiceLink"`
	PoolInvoicePending      string `json:"poolInvoicePending"`
	PoolGenerateInvoiceCta  string `json:"poolGenerateInvoiceCta"`
	PerCallRecognizeCta     string `json:"perCallRecognizeCta"`
	PerCallInvoiceLink      string `json:"perCallInvoiceLink"`
	PerCallNotReady         string `json:"perCallNotReady"`
}

// SubscriptionBackfillLabels holds labels for the Backfill cycle Jobs drawer.
// Lyngua key: `subscription.detail.backfill.*`. See cyclic-subscription-jobs
// plan §9.2.
type SubscriptionBackfillLabels struct {
	DrawerTitle       string `json:"drawerTitle"`
	DrawerDescription string `json:"drawerDescription"`
	PreviewLine       string `json:"previewLine"`
	CountLabel        string `json:"countLabel"`
	Confirm           string `json:"confirm"`
	Cancel            string `json:"cancel"`
	MaxWarning        string `json:"maxWarning"`
}

// SubscriptionJobsTabLabels holds labels for the new flat Jobs tab on the
// Subscription detail page. Lyngua key: `subscription.detail.jobs.*`. See
// cyclic-subscription-jobs plan §21.
type SubscriptionJobsTabLabels struct {
	Heading          string `json:"heading"`
	Empty            string `json:"empty"`
	FilterStatus     string `json:"filterStatus"`
	FilterType       string `json:"filterType"`
	FilterAll        string `json:"filterAll"`
	SortBy           string `json:"sortBy"`
	SortByCycle      string `json:"sortByCycle"`
	ExportCsv        string `json:"exportCsv"`
	Summary          string `json:"summary"`
	ColumnNumber     string `json:"columnNumber"`
	ColumnName       string `json:"columnName"`
	ColumnType       string `json:"columnType"`
	ColumnPhase      string `json:"columnPhase"`
	ColumnStatus     string `json:"columnStatus"`
	ColumnPeriod     string `json:"columnPeriod"`
	TypeEngagement   string `json:"typeEngagement"`
	TypeOnboarding   string `json:"typeOnboarding"`
	TypeCycle        string `json:"typeCycle"`
	TypeVisit        string `json:"typeVisit"`
	SpawnFailedToast string `json:"spawnFailedToast"`
}

// SubscriptionSpawnLabels holds labels for the retroactive Spawn Jobs drawer.
// Lyngua key: `subscription.spawn.*`. See auto-spawn-jobs-from-subscription
// plan §5.3 / §9.
type SubscriptionSpawnLabels struct {
	Title             string `json:"title"`
	DetectedTemplates string `json:"detectedTemplates"`
	RootTemplate      string `json:"rootTemplate"`
	Cancel            string `json:"cancel"`
	Confirm           string `json:"confirm"`
	SuccessToast      string `json:"successToast"`
	Skipped           string `json:"skipped"`
}

type SubscriptionConfirmLabels struct {
	Cancel                string `json:"cancel"`
	CancelMessage         string `json:"cancelMessage"`
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
			AddPlan:       "Add Plan",
			AddPricePlan:  "Add Price Plan",
			EditPricePlan: "Edit Price Plan",
			AddProduct:    "Add Product",
		},
		Columns: PlanColumnLabels{
			Name:        "Name",
			Description: "Description",
			Interval:    "Interval",
			Price:       "Price",
			Status:      "Status",
			Product:     "Product",
			PricePlan:     "Price Plan",
			PriceSchedule: "Price Schedule",
			Duration:      "Duration",
			Location:    "Location",
			ItemType:    "Item Type",
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
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			NameInfo:        "Display name for this plan. Shown in subscription lists and invoices.",
			DescriptionInfo: "Optional notes about this plan. Visible on detail pages.",
			ActiveInfo:      "Inactive plans are hidden from new subscriptions.",
			// Client-scope fields (2026-04-27 plan-client-scope plan §7).
			ClientLabel:             "Client",
			ClientHelp:              "Leave blank to make this package available for any client. Set a client to make it a custom package for that client only.",
			ClientPlaceholder:       "Leave blank for a general package",
			ClientSearchPlaceholder: "Search clients...",
			ClientNoResults:         "No clients found",
			ClientLockedTooltip:     "Locked — this plan has active subscriptions. Detach them or create a new plan.",
			ClientForLabel:          "For {{.ClientName}}",
			ClientInfo:              "Optional. When set, this plan only appears for engagements with that client.",
			JobTemplate:             "Job Template",
			JobTemplateNone:         "(none — engagement has no operational tracking)",
			JobTemplateHint:         "Select the operational template that defines the work for this engagement. Leave empty for advisory-only plans.",
			// 2026-04-30 cyclic-subscription-jobs plan §9.3.
			VisitsPerCycleLabel:       "Visits per billing cycle",
			VisitsPerCyclePlaceholder: "1",
			VisitsPerCycleHint:        "Number of cycle Job instances per billing cycle. Default 1. Use 2 for biweekly visits billed monthly, 4 for weekly visits billed monthly.",
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
			Info:          "Information",
			Products:      "Products",
			PricePlan:     "Rate Cards",
			PricePlanSlug: "",
			Attachments:   "Attachments",
			AuditTrail:    "Audit Trail",
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
			PermissionDenied:  "You do not have permission to perform this action",
			InvalidFormData:   "Invalid form data. Please check your inputs and try again.",
			NotFound:          "Plan not found",
			IDRequired:        "Plan ID is required",
			NoIDsProvided:     "No plan IDs provided",
			InvalidStatus:     "Invalid status",
			NoPermission:      "No permission",
			CannotDelete:      "This plan cannot be deleted because it has products or rate cards assigned",
			ClientScopeLocked: "Cannot change this plan's client while it has active subscriptions.",
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
			Schedule:            "Price Schedule",
			SchedulePlaceholder: "Select a schedule...",
			Location:            "Location",
			LocationPlaceholder: "Select a location...",
			SelectLocation:      "— No location (all locations) —",
			Active:              "Active",
		},
		ProductPlanForm: ProductPlanFormLabels{
			Product:            "Product",
			ProductPlaceholder: "Select an item...",
			SelectProduct:      "— Select a product —",
			Active:             "Active",
			ProductKindLabel:   "Item Type",
			ProductKind: ProductKindOptionLabels{
				Service:        "Service",
				StockedGood:    "Stocked Good",
				NonStockedGood: "Non-Stocked Good",
				Consumable:     "Consumable",
			},
			// Model D — variant picker defaults
			VariantSelectLabel:       "Variant",
			VariantSelectPlaceholder: "Select a variant",
			VariantSelectInfo:        "Required when the parent product has variants enabled.",
		},
		Filters: PlanFilterLabels{
			ScopeChipLabel: "Show:",
			ScopeMaster:    "Master",
			ScopeClient:    "Client-specific",
			ScopeAll:       "All",
		},
	}
}

// ---------------------------------------------------------------------------
// Price Plan labels
// ---------------------------------------------------------------------------

// PricePlanLabels holds all labels for the standalone price plan (rate card) module.
type PricePlanLabels struct {
	Page         PricePlanPageLabels         `json:"page"`
	Buttons      PricePlanButtonLabels       `json:"buttons"`
	Columns      PricePlanColumnLabels2      `json:"columns"`
	Empty        PricePlanEmptyLabels        `json:"empty"`
	Form         PricePlanFormLabels         `json:"form"`
	Actions      PricePlanActionLabels       `json:"actions"`
	Bulk         PricePlanBulkLabels         `json:"bulk"`
	Detail       PricePlanDetailLabels2      `json:"detail"`
	Tabs         PricePlanTabLabels2         `json:"tabs"`
	Confirm      PricePlanConfirmLabels      `json:"confirm"`
	Errors       PricePlanErrorLabels        `json:"errors"`
	ProductPrice PricePlanProductPriceLabels `json:"productPrice"`
	Messages     PricePlanMessageLabels      `json:"messages"`
}

// PricePlanProductPriceLabels holds labels for product-price sub-table actions and empty state.
type PricePlanProductPriceLabels struct {
	EditTitle   string `json:"editTitle"`
	DeleteTitle string `json:"deleteTitle"`
	EmptyTitle  string `json:"emptyTitle"`
	EmptyMsg    string `json:"emptyMsg"`
}

// PricePlanMessageLabels holds translatable message strings used in the price plan
// and price schedule plan views (pricing-lock notices, validation errors).
type PricePlanMessageLabels struct {
	PricingLockedReason     string `json:"pricingLockedReason"`
	ItemPricingLockedReason string `json:"itemPricingLockedReason"`
	CreateNotAvailable      string `json:"createNotAvailable"`
	UpdateNotAvailable      string `json:"updateNotAvailable"`
	ProductRequired         string `json:"productRequired"`
	InvalidPrice            string `json:"invalidPrice"`
	InUseCannotModify       string `json:"inUseCannotModify"`
	IDRequired              string `json:"idRequired"`
	DeleteNotAvailable      string `json:"deleteNotAvailable"`
	CurrencyMismatch        string `json:"currencyMismatch"`
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
	Schedule    string `json:"schedule"`
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

	// Info-tab field labels (price-schedule-plan-tab-info).
	Heading       string `json:"heading"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Duration      string `json:"duration"`
	ScheduleLabel string `json:"scheduleLabel"`
	Status        string `json:"status"`
	DateCreated   string `json:"dateCreated"`
	DateModified  string `json:"dateModified"`
	Edit          string `json:"edit"`
	EditTitle     string `json:"editTitle"`

	// 2026-04-30 cyclic-subscription-jobs plan §20 — Billing model summary
	// rendered on the info tab. Lyngua key: `pricePlan.detail.summary*`.
	SummaryHeading            string                       `json:"summaryHeading"`
	CustomerHeading           string                       `json:"customerHeading"`
	OperationsHeading         string                       `json:"operationsHeading"`
	RevenueRecognitionHeading string                       `json:"revenueRecognitionHeading"`
	Summary                   PricePlanBillingSummaryCopy  `json:"summary"`
	Warning                   PricePlanBillingSummaryWarn  `json:"warning"`
}

// PricePlanBillingSummaryCopy carries the per-(kind × basis) lyngua copy that
// `buildBillingModelSummary` projects into the info-sections grid.
// Cyclic plan ships rows 1-6 (oneTime / recurring / contract / milestone);
// AD_HOC plan adds adHoc.* in a follow-up. Each entry has 3 lines:
// customer, operations, revenue.
type PricePlanBillingSummaryCopy struct {
	OneTime   PricePlanSummaryByBasis `json:"oneTime"`
	Recurring PricePlanSummaryByBasis `json:"recurring"`
	Contract  PricePlanSummaryByBasis `json:"contract"`
	Milestone PricePlanSummaryByBasis `json:"milestone"`
	AdHoc     PricePlanSummaryByBasis `json:"adHoc"`
}

// PricePlanSummaryByBasis groups the text lines per basis. Empty
// strings on a basis means "no copy for that combo" — view skips it.
type PricePlanSummaryByBasis struct {
	PerCycle         PricePlanSummaryLines `json:"perCycle"`
	TotalPackage     PricePlanSummaryLines `json:"totalPackage"`
	DerivedFromLines PricePlanSummaryLines `json:"derivedFromLines"`
	PerOccurrence    PricePlanSummaryLines `json:"perOccurrence"`
}

// PricePlanSummaryLines holds the 3 lines for a kind × basis cell.
type PricePlanSummaryLines struct {
	Customer   string `json:"customer"`
	Operations string `json:"operations"`
	Revenue    string `json:"revenue"`
}

// PricePlanBillingSummaryWarn carries the warning-row copy keyed by symbol
// per plan §20.3. View only renders entries whose preconditions trip.
type PricePlanBillingSummaryWarn struct {
	MilestoneNoTemplate           string `json:"milestoneNoTemplate"`
	RecurringNoTemplate           string `json:"recurringNoTemplate"`
	VisitsPerCycleInvalidKind     string `json:"visitsPerCycleInvalidKind"`
	AdHocPoolNoTemplate           string `json:"adHocPoolNoTemplate"`
	AdHocPerCallNoTemplate        string `json:"adHocPerCallNoTemplate"`
	AdHocNoEntitlement            string `json:"adHocNoEntitlement"`
	AdHocBillingCycleNotAllowed   string `json:"adHocBillingCycleNotAllowed"`
	AdHocVisitsPerCycleNotAllowed string `json:"adHocVisitsPerCycleNotAllowed"`
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

	// 2026-04-27 plan-client-scope plan §3.5 — fired by the centymo confirm
	// dialog when an operator changes monetary fields on a client-scoped
	// PricePlan that has N > 1 active subscriptions. Templated via
	// {{.Count}} and {{.ClientName}}.
	EditAmountMultipleEngagements string `json:"editAmountMultipleEngagements"`
}

type PricePlanErrorLabels struct {
	NotFound     string `json:"notFound"`
	LoadFailed   string `json:"loadFailed"`
	Unauthorized string `json:"unauthorized"`
	CreateFailed string `json:"createFailed"`
	UpdateFailed string `json:"updateFailed"`
	DeleteFailed string `json:"deleteFailed"`
	InUse        string `json:"inUse"`

	// 2026-04-27 plan-client-scope plan §7. Surfaced when an UpdatePricePlan
	// body sends a client_id that doesn't match the parent Plan's client_id.
	ClientScopeMismatch string `json:"clientScopeMismatch"`
	// 2026-04-28 — surfaced when the operator picks a price_schedule whose
	// client_id belongs to a different client than the parent Plan. Master
	// schedules (sched.client_id == "") are still accepted; only the
	// cross-client cases get rejected.
	ScheduleClientMismatch string `json:"scheduleClientMismatch"`
	// 2026-04-28 — surfaced when an operator submits a PricePlan with no
	// price_schedule_id under a client-scoped Plan. The use case used to
	// auto-create a schedule with a synthetic now() date; reverted because
	// that hid real operator intent. Operator must pick or create a client
	// rate card first.
	ScheduleRequiredForClientScope string `json:"scheduleRequiredForClientScope"`
	// Server-side-only error key — the centymo confirm dialog catches the
	// N>1-engagements gate before this surfaces.
	MultiEngagementConfirmRequired string `json:"multiEngagementConfirmRequired"`
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
			Schedule:    "Schedule",
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
			CurrencyPHP:         "PHP (₱)",
			CurrencyUSD:         "USD ($)",
			DurationValue:       "Duration",
			DurationUnit:        "Unit",
			Schedule:            "Price Schedule",
			SchedulePlaceholder: "Select a schedule...",
			ScheduleSearch:      "Filter...",
			Location:            "Location",
			LocationPlaceholder: "Select a location...",
			LocationHintPrefix:  "Location: ",
			SelectLocation:      "— No location (all locations) —",
			Active:              "Active",
			PlanLabel:           "Package",
			PlanPlaceholder:     "Select a package...",
			PlanSearch:          "Filter...",
			// Wave 2 new fields
			SectionBasic:                "Basic info",
			SectionPricing:              "Pricing",
			BillingKindLabel:            "Billing model",
			BillingKindOneTime:          "One-time",
			BillingKindRecurring:        "Recurring retainer",
			BillingKindContract:         "Fixed-term engagement",
			BillingKindMilestone:        "Milestone",
			AmountBasisLabel:            "Amount basis",
			AmountBasisPerCycle:         "Per cycle",
			AmountBasisTotalPackage:     "Total package",
			AmountBasisDerivedFromLines: "Sum of items",
			BillingCycleLabel:           "Billing cycle",
			BillingCyclePlaceholder:     "e.g. every 1 month",
			TermLabel:                   "Term",
			TermPlaceholder:             "e.g. 12 months",
			TermOpenEndedHelp:           "Leave empty for open-ended / no expiration",
			// Field-level info popovers — use proto-generic wording; business-type
			// tiers override via lyngua (e.g. "plan" → "package" / "rate card").
			PlanInfo:         "The plan this price plan belongs to. Locked from the parent page.",
			ScheduleInfo:     "The price schedule (date range + location) this price plan belongs to.",
			NameInfo:         "Optional — defaults to the plan name when left blank.",
			DescriptionInfo:  "Optional notes shown alongside the price plan in detail views.",
			BillingKindInfo:  "One-time = charged once. Recurring = billed every cycle. Fixed-term = recurring with an end date.",
			AmountBasisInfo:  "Per cycle = amount charged each billing cycle. Total package = amount charged across the full term. Sum of items = derived from the per-item breakdown.",
			AmountInfo:       "Price in the selected currency. For Sum of items, this is computed automatically.",
			CurrencyInfo:     "Currency applied to this price plan and any auto-seeded product price plans.",
			BillingCycleInfo: "How often the recurring charge is issued (e.g. every 1 month).",
			TermInfo:         "How long the engagement lasts. Leave empty for open-ended / no expiration.",
			ActiveInfo:       "Inactive price plans stay on record but are hidden from new subscriptions.",
			// 2026-04-27 plan-client-scope plan §6.7 — info banner shown above
			// the form when its parent PriceSchedule is client-scoped.
			ParentScheduleClientNotice: "This price schedule belongs to {{.ClientName}}. Price plans created here will be available only for engagements with this client.",
			// 2026-04-27 plan-client-scope plan §6.7 — tooltip on the readonly
			// Schedule label when the parent Plan is client-scoped. Proto-generic
			// wording; tier overrides live in lyngua.
			ScheduleLockedTooltip:  "This price plan is bound to {{.ClientName}}'s price schedule.",
			ScheduleAutoCreateHint: "No price schedule exists for {{.ClientName}} yet — one will be created automatically when you save.",
			ScheduleAutoReuseHint:  "This price plan will be added to the existing price schedule for {{.ClientName}}.",
			// 2026-04-30 cyclic-subscription-jobs plan §9.4.
			MilestoneCyclicBlock: "Milestone billing is not supported on cyclic plans (RECURRING / CONTRACT × PER_CYCLE / multi-visit).",
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
			Heading:        "Plan Info",
			Name:           "Name",
			Description:    "Description",
			Amount:         "Amount",
			Currency:       "Currency",
			Duration:       "Duration",
			ScheduleLabel:  "Schedule",
			Status:         "Status",
			DateCreated:    "Date Created",
			DateModified:   "Date Modified",
			Edit:           "Edit Price Plan",
			EditTitle:      "Edit Price Plan",
			// 2026-04-30 cyclic-subscription-jobs plan §20.
			SummaryHeading:            "Billing model summary",
			CustomerHeading:           "Customer experience",
			OperationsHeading:         "Operations impact",
			RevenueRecognitionHeading: "Revenue recognition",
			Summary: PricePlanBillingSummaryCopy{
				OneTime: PricePlanSummaryByBasis{
					TotalPackage: PricePlanSummaryLines{
						Customer:   "Pays {{.Amount}} once at signup. No further charges.",
						Operations: "Engagement spawns 1 lifetime Job with phases (if Plan has a JobTemplate).",
						Revenue:    "One Revenue at Subscription.Create covering the full amount.",
					},
				},
				Recurring: PricePlanSummaryByBasis{
					PerCycle: PricePlanSummaryLines{
						Customer:   "Charged {{.Amount}} every {{.CycleLabel}}. Subscription auto-renews until cancelled.",
						Operations: "Each cycle spawns {{.VisitsPerCycle}} cycle Job(s) (if Plan has a JobTemplate). Operations tab shows cycle accordions.",
						Revenue:    "One Revenue per cycle. Recognize Revenue creates the invoice and (via piggyback) spawns the cycle Job if missing.",
					},
					DerivedFromLines: PricePlanSummaryLines{
						Customer:   "Charged the sum of itemised lines every {{.CycleLabel}}.",
						Operations: "Each cycle spawns 1+ cycle Jobs. Operations tracking flows through Plan's JobTemplate.",
						Revenue:    "Revenue total computed from ProductPricePlan rows; one Revenue per cycle.",
					},
				},
				Contract: PricePlanSummaryByBasis{
					PerCycle: PricePlanSummaryLines{
						Customer:   "Charged {{.Amount}} every {{.CycleLabel}} for {{.TermLength}}. Auto-deactivates at term end.",
						Operations: "Same as recurring + the engagement closes when the {{.TermLength}} term completes.",
						Revenue:    "Same as recurring. Operator can extend the term to spawn additional cycles.",
					},
					TotalPackage: PricePlanSummaryLines{
						Customer:   "Pays {{.Amount}} upfront for {{.TermLength}} of service.",
						Operations: "Engagement spawns 1 lifetime Job (or cycle Jobs if cyclic — see Plan's visits_per_cycle).",
						Revenue:    "One Revenue at signup; cycle Jobs are operational only.",
					},
				},
				Milestone: PricePlanSummaryByBasis{
					TotalPackage: PricePlanSummaryLines{
						Customer:   "Pays {{.Amount}} total. Invoice fires per milestone (engagement phase) as work completes.",
						Operations: "Lifetime engagement Job with phases. BillingEvent rows gate per-milestone invoicing.",
						Revenue:    "Revenue per milestone trigger; sum across milestones equals the total package.",
					},
				},
			},
			Warning: PricePlanBillingSummaryWarn{
				MilestoneNoTemplate:       "Milestone billing requires the Plan to have a JobTemplate. Configure it on the Plan first.",
				RecurringNoTemplate:       "This subscription will not have operational tracking. Add a JobTemplate to the Plan to enable cycle Jobs.",
				VisitsPerCycleInvalidKind: "visits_per_cycle is only valid for cyclic plans. Reset to 1 or change the billing kind.",
			},
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
			// 2026-04-27 plan-client-scope plan §3.5 / §7.
			EditAmountMultipleEngagements: "This price plan is attached to {{.Count}} active subscriptions for {{.ClientName}}. Changing the amount or cycle will affect all of them on the next bill cycle. Continue?",
		},
		Errors: PricePlanErrorLabels{
			NotFound:                       "Rate card not found.",
			LoadFailed:                     "Failed to load rate cards.",
			Unauthorized:                   "You do not have permission to access this resource.",
			CreateFailed:                   "Failed to create rate card.",
			UpdateFailed:                   "Failed to update rate card.",
			DeleteFailed:                   "Failed to delete rate card.",
			InUse:                          "This price plan is in use by active subscriptions and cannot be deleted.",
			ClientScopeMismatch:            "Price plan client must match its parent plan's client.",
			ScheduleClientMismatch:         "Selected schedule belongs to a different client and cannot be attached to this price plan.",
			ScheduleRequiredForClientScope: "This package is scoped to a client. Pick or create a rate card for that client before adding a price plan.",
			MultiEngagementConfirmRequired: "Confirmation required — multiple attached subscriptions and monetary fields changing.",
		},
		ProductPrice: PricePlanProductPriceLabels{
			EditTitle:   "Edit Product Price",
			DeleteTitle: "Delete Product Price",
			EmptyTitle:  "No Product Prices",
			EmptyMsg:    "No product prices have been configured for this rate card yet.",
		},
		Messages: PricePlanMessageLabels{
			PricingLockedReason:     "This plan is in use by active subscriptions. Pricing changes are disabled. You can still rename or reassign the package.",
			ItemPricingLockedReason: "This package is in use by active engagements. Item price and currency are locked to keep billing consistent.",
			CreateNotAvailable:      "Product price plan create is not available.",
			UpdateNotAvailable:      "Product price plan update is not available.",
			ProductRequired:         "Product is required.",
			InvalidPrice:            "Invalid price value.",
			InUseCannotModify:       "This package is in use by active engagements. Item price and currency are locked.",
			IDRequired:              "ID is required.",
			DeleteNotAvailable:      "Product price plan delete is not available.",
			CurrencyMismatch:        "Currency must match the rate card currency.",
		},
	}
}

// ---------------------------------------------------------------------------
// Price Schedule labels
// ---------------------------------------------------------------------------

// PriceScheduleFilterLabels holds translatable labels for the scope filter chip
// on the price schedule list page (§6.1 of the 2026-04-27 plan-client-scope plan).
type PriceScheduleFilterLabels struct {
	ScopeChipLabel string `json:"scopeChipLabel"`
	ScopeMaster    string `json:"scopeMaster"`
	ScopeClient    string `json:"scopeClient"`
	ScopeAll       string `json:"scopeAll"`
}

// PriceScheduleLabels holds all labels for the price schedule module.
type PriceScheduleLabels struct {
	Page     PriceSchedulePageLabels      `json:"page"`
	Buttons  PriceScheduleButtonLabels    `json:"buttons"`
	Columns  PriceScheduleColumnLabels    `json:"columns"`
	Empty    PriceScheduleEmptyLabels     `json:"empty"`
	Form     PriceScheduleFormLabels      `json:"form"`
	PlanForm PriceSchedulePlanFormLabels  `json:"planForm"`
	Bulk     PriceScheduleBulkLabels      `json:"bulk"`
	Confirm  PriceScheduleConfirmLabels   `json:"confirm"`
	Tabs     PriceScheduleTabLabels       `json:"tabs"`
	Detail   PriceScheduleDetailLabels    `json:"detail"`
	Errors   PriceScheduleErrorLabels     `json:"errors"`
	Filters  PriceScheduleFilterLabels    `json:"filters"`
}

// PriceSchedulePlanFormLabels holds labels for the "Add Plan" (price_plan) drawer form
// within a price schedule. Professional tier overrides field names (e.g., "Package").
type PriceSchedulePlanFormLabels struct {
	SectionSchedule        string `json:"sectionSchedule"`
	SectionPackage         string `json:"sectionPackage"`
	SectionPricing         string `json:"sectionPricing"`
	PriceScheduleField     string `json:"priceScheduleField"`
	PackageLabel           string `json:"packageLabel"`
	PackagePlaceholder     string `json:"packagePlaceholder"`
	PackageSearch          string `json:"packageSearch"`
	NameLabel              string `json:"nameLabel"`
	NamePlaceholder        string `json:"namePlaceholder"`
	DescriptionLabel       string `json:"descriptionLabel"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	AmountLabel            string `json:"amountLabel"`
	AmountPlaceholder      string `json:"amountPlaceholder"`
	CurrencyLabel          string `json:"currencyLabel"`
	CurrencyPlaceholder    string `json:"currencyPlaceholder"`
	CurrencyPHP            string `json:"currencyPHP"`
	CurrencyUSD            string `json:"currencyUSD"`
	DurationLabel          string `json:"durationLabel"`
	UnitLabel              string `json:"unitLabel"`
	ActiveLabel            string `json:"activeLabel"`
	SchedulePlaceholder    string `json:"schedulePlaceholder"`
	ScheduleSearch         string `json:"scheduleSearch"`
	LocationHintPrefix     string `json:"locationHintPrefix"`
}

type PriceSchedulePageLabels struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	ActiveTitle   string `json:"activeTitle"`
	InactiveTitle string `json:"inactiveTitle"`
}

type PriceScheduleButtonLabels struct {
	View       string `json:"view"`
	Add        string `json:"add"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	BulkDelete string `json:"bulkDelete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type PriceScheduleColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DateStart   string `json:"dateStart"`
	DateEnd     string `json:"dateEnd"`
	Location    string `json:"location"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
	Actions     string `json:"actions"`
}

type PriceScheduleEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type PriceScheduleFormLabels struct {
	Name                string `json:"name"`
	NamePlaceholder     string `json:"namePlaceholder"`
	Description         string `json:"description"`
	DescPlaceholder     string `json:"descPlaceholder"`
	DateStart           string `json:"dateStart"`
	DateEnd             string `json:"dateEnd"`
	// Optional time inputs paired with DateStart/DateEnd (2026-04-28 date+time
	// field plan). TimePlaceholder is shared by both inputs.
	TimeStart        string `json:"timeStart"`
	TimeEnd          string `json:"timeEnd"`
	TimePlaceholder  string `json:"timePlaceholder"`
	Location            string `json:"location"`
	LocationPlaceholder string `json:"locationPlaceholder"`
	SelectLocation      string `json:"selectLocation"`
	Active              string `json:"active"`

	// Wave 2 — section headers (from lyngua price_schedule.json → priceSchedule.form)
	SectionScheduleDetails string `json:"sectionScheduleDetails"`
	SectionDateRange       string `json:"sectionDateRange"`
	SectionLocation        string `json:"sectionLocation"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	DateStartInfo   string `json:"dateStartInfo"`
	DateEndInfo     string `json:"dateEndInfo"`
	TimeStartInfo   string `json:"timeStartInfo"`
	TimeEndInfo     string `json:"timeEndInfo"`
	LocationInfo    string `json:"locationInfo"`
	ActiveInfo      string `json:"activeInfo"`

	// Client-scope fields (2026-04-27 plan-client-scope plan §7).
	// Set on the schedule add/edit drawer Client picker. The suffix is
	// appended to the client's name to produce the default schedule name
	// (e.g. "Cruz Engineering - Rate Cards" on professional tier, or
	// "Cruz Engineering - Price Schedule" on general). See plan §4.4.1.
	ClientLabel                          string `json:"clientLabel"`
	ClientHelp                           string `json:"clientHelp"`
	ClientPlaceholder                    string `json:"clientPlaceholder"`
	ClientSearchPlaceholder              string `json:"clientSearchPlaceholder"`
	ClientNoResults                      string `json:"clientNoResults"`
	ClientInfo                           string `json:"clientInfo"`
	CustomClientPriceScheduleLabelSuffix string `json:"customClientPriceScheduleLabelSuffix"`
	LocationSearchPlaceholder            string `json:"locationSearchPlaceholder"`

	// Scope radio (2026-04-28) — mutually exclusive Location / Client picker.
	ScopeLabel               string `json:"scopeLabel"`
	ScopeInfo                string `json:"scopeInfo"`
	ScopeOptionLocation      string `json:"scopeOptionLocation"`
	ScopeOptionClient        string `json:"scopeOptionClient"`
	ScopeOptionLocationHelp  string `json:"scopeOptionLocationHelp"`
	ScopeOptionClientHelp    string `json:"scopeOptionClientHelp"`
}

type PriceScheduleBulkLabels struct {
	DeleteTitle       string `json:"deleteTitle"`
	DeleteMessage     string `json:"deleteMessage"`
	ActivateTitle     string `json:"activateTitle"`
	ActivateMessage   string `json:"activateMessage"`
	DeactivateTitle   string `json:"deactivateTitle"`
	DeactivateMessage string `json:"deactivateMessage"`
}

type PriceScheduleConfirmLabels struct {
	DeleteTitle       string `json:"deleteTitle"`
	DeleteMessage     string `json:"deleteMessage"`
	ActivateTitle     string `json:"activateTitle"`
	ActivateMessage   string `json:"activateMessage"`
	DeactivateTitle   string `json:"deactivateTitle"`
	DeactivateMessage string `json:"deactivateMessage"`
}

type PriceScheduleTabLabels struct {
	Info              string `json:"info"`
	PricePlan         string `json:"pricePlan"`
	PricePlanSlug     string `json:"pricePlanSlug"`
	ProductPrices     string `json:"productPrices"`
	ProductPricesSlug string `json:"productPricesSlug"`
}

// ResolveTabSlug returns the URL slug for a canonical tab key. Today only the
// "pricePlan" tab on the parent detail and "product-prices" on the nested plan
// detail are re-slugged (e.g., professional tier ships "package-prices" /
// "package-item-prices"); other tabs round-trip through as-is.
func (t PriceScheduleTabLabels) ResolveTabSlug(canonical string) string {
	switch canonical {
	case "pricePlan":
		if s := strings.TrimSpace(t.PricePlanSlug); s != "" {
			return s
		}
	case "product-prices":
		if s := strings.TrimSpace(t.ProductPricesSlug); s != "" {
			return s
		}
	}
	return canonical
}

// CanonicalizeTab maps an incoming URL tab slug back to its canonical key so
// internal template lookups and equality checks stay tier-agnostic.
func (t PriceScheduleTabLabels) CanonicalizeTab(slug string) string {
	if slug == "" {
		return ""
	}
	if s := strings.TrimSpace(t.PricePlanSlug); s != "" && slug == s {
		return "pricePlan"
	}
	if s := strings.TrimSpace(t.ProductPricesSlug); s != "" && slug == s {
		return "product-prices"
	}
	return slug
}

type PriceScheduleDetailLabels struct {
	Title           string `json:"title"`
	DateCreated     string `json:"dateCreated"`
	DateModified    string `json:"dateModified"`
	NoLocation      string `json:"noLocation"`
	NoDateEnd       string `json:"noDateEnd"`
	NoDescription   string `json:"noDescription"`
	PlansEmptyTitle      string `json:"plansEmptyTitle"`
	PlansEmptyMsg        string `json:"plansEmptyMsg"`
	NoDescriptionSubtitle string `json:"noDescriptionSubtitle"`

	// Product price (per-product breakdown, shown on the schedule-scoped plan detail).
	// Professional tier renames these to "Service Price" via lyngua.
	ProductPriceAdd           string `json:"productPriceAdd"`
	ProductPriceEdit          string `json:"productPriceEdit"`
	ProductPriceDelete        string `json:"productPriceDelete"`
	ProductPriceDeleteConfirm string `json:"productPriceDeleteConfirm"`
	ProductPriceEmptyTitle    string `json:"productPriceEmptyTitle"`
	ProductPriceEmptyMsg      string `json:"productPriceEmptyMsg"`
	ProductPriceSection       string `json:"productPriceSection"` // drawer section title ("Product Price" / "Service Price")
	ProductField              string `json:"productField"`        // drawer product select label ("Product" / "Service")

	// Plans table columns (price-schedule-detail plans tab).
	PlanColumnPlan     string `json:"planColumnPlan"`
	PlanColumnAmount   string `json:"planColumnAmount"`
	PlanColumnDuration string `json:"planColumnDuration"`
	PlanColumnStatus   string `json:"planColumnStatus"`

	// Plans table row actions + confirms.
	PlanView            string `json:"planView"`
	PlanEdit            string `json:"planEdit"`
	PlanEditDrawerTitle string `json:"planEditDrawerTitle"`
	PlanDelete          string `json:"planDelete"`
	PlanDeleteTitle     string `json:"planDeleteTitle"`
	PlanDeleteMsg       string `json:"planDeleteMsg"`
	PlanInUseTooltip    string `json:"planInUseTooltip"`

	// Plans table primary action + inline error messages.
	PlanAdd          string `json:"planAdd"`
	PlanRequired     string `json:"planRequired"`

	// Product prices table columns.
	ProductPriceColumnProduct   string `json:"productPriceColumnProduct"`
	ProductPriceColumnPrice     string `json:"productPriceColumnPrice"`
	ProductPriceColumnCurrency  string `json:"productPriceColumnCurrency"`
	ProductPriceColumnTreatment string `json:"productPriceColumnTreatment"`
	ProductPriceColumnEffective string `json:"productPriceColumnEffective"`

	// Drawer banners explaining how the per-line price relates to the parent
	// PricePlan.amount_basis. Surfaced above the Price input.
	BasisBannerPerCycle     string `json:"basisBannerPerCycle"`
	BasisBannerTotalPackage string `json:"basisBannerTotalPackage"`
	BasisBannerDerived      string `json:"basisBannerDerived"`

	// Drawer section labels used by the schedule-scoped PPP drawer.
	ProductPriceCatalogSection string `json:"productPriceCatalogSection"`
	ProductPricePricingSection string `json:"productPricePricingSection"`
	ProductPriceEffectiveSection string `json:"productPriceEffectiveSection"`
}

type PriceScheduleErrorLabels struct {
	NotFound                    string `json:"notFound"`
	LoadFailed                  string `json:"loadFailed"`
	Unauthorized                string `json:"unauthorized"`
	CreateFailed                string `json:"createFailed"`
	UpdateFailed                string `json:"updateFailed"`
	DeleteFailed                string `json:"deleteFailed"`
	InUse                       string `json:"inUse"`
	PricePlanCreateUnavailable  string `json:"pricePlanCreateUnavailable"`
}

// DefaultPriceScheduleLabels returns PriceScheduleLabels with sensible English defaults.
func DefaultPriceScheduleLabels() PriceScheduleLabels {
	return PriceScheduleLabels{
		Page: PriceSchedulePageLabels{
			Title:         "Price Schedules",
			Subtitle:      "Manage your price schedules",
			ActiveTitle:   "Active Price Schedules",
			InactiveTitle: "Inactive Price Schedules",
		},
		Buttons: PriceScheduleButtonLabels{
			View:       "View",
			Add:        "Add Price Schedule",
			Edit:       "Edit Price Schedule",
			Delete:     "Delete Price Schedule",
			BulkDelete: "Delete Price Schedules",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: PriceScheduleColumnLabels{
			Name:        "Name",
			Description: "Description",
			DateStart:   "Start Date",
			DateEnd:     "End Date",
			Location:    "Location",
			Status:      "Status",
			DateCreated: "Date Created",
			Actions:     "Actions",
		},
		Empty: PriceScheduleEmptyLabels{
			Title:   "No Price Schedules",
			Message: "No price schedules to display.",
		},
		Form: PriceScheduleFormLabels{
			Name:                "Name",
			NamePlaceholder:     "Enter price schedule name",
			Description:         "Description",
			DescPlaceholder:     "Enter description...",
			DateStart:           "Start Date",
			DateEnd:             "End Date",
			TimeStart:           "Start Time (optional)",
			TimeEnd:             "End Time (optional)",
			TimePlaceholder:     "HH:MM",
			Location:            "Location",
			LocationPlaceholder: "Select a location...",
			SelectLocation:      "— No location (all locations) —",
			Active:              "Active",
			// Wave 2 new section headers
			SectionScheduleDetails: "Schedule details",
			SectionDateRange:       "Date range",
			SectionLocation:        "Location",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			NameInfo:        "A short display name for this price schedule.",
			DescriptionInfo: "Optional notes or context for this price schedule.",
			DateStartInfo:   "First date this price schedule becomes effective.",
			DateEndInfo:     "Last date this price schedule is effective. Leave empty for no end date.",
			TimeStartInfo:   "Optional time of day in the operator's display timezone. Leave blank for start of day (00:00).",
			TimeEndInfo:     "Optional time of day in the operator's display timezone. Leave blank for end of day (23:59).",
			LocationInfo:    "Restrict this price schedule to a specific location, or leave empty to apply to all locations.",
			ActiveInfo:      "Inactive price schedules are hidden from new subscriptions.",
			// Client-scope fields (2026-04-27 plan-client-scope plan §7).
			ClientLabel:                          "Client",
			ClientHelp:                           "Leave blank for a general schedule. Set a client to create a bespoke schedule reused across that client's price plans.",
			ClientPlaceholder:                    "Leave blank for a general schedule",
			ClientSearchPlaceholder:              "Search clients...",
			ClientNoResults:                      "No clients found",
			ClientInfo:                           "Optional. When set, this schedule is reserved for that client's bespoke price plans.",
			CustomClientPriceScheduleLabelSuffix: "Price Schedule",
			LocationSearchPlaceholder:            "Filter...",
			// Scope radio (2026-04-28).
			ScopeLabel:              "Scope",
			ScopeInfo:               "Choose whether this schedule is shared across every client at a location, or reserved for one client's bespoke pricing. Switching scope clears the inactive picker on save.",
			ScopeOptionLocation:     "Location-scoped",
			ScopeOptionClient:       "Client-scoped",
			ScopeOptionLocationHelp: "Reusable across all clients at this location.",
			ScopeOptionClientHelp:   "Reserved for one client's bespoke pricing.",
		},
		Bulk: PriceScheduleBulkLabels{
			DeleteTitle:       "Delete Price Schedules",
			DeleteMessage:     "Permanently delete the selected price schedules? This cannot be undone.",
			ActivateTitle:     "Activate Price Schedules",
			ActivateMessage:   "Activate the selected price schedules?",
			DeactivateTitle:   "Deactivate Price Schedules",
			DeactivateMessage: "Deactivate the selected price schedules?",
		},
		Confirm: PriceScheduleConfirmLabels{
			DeleteTitle:       "Delete Price Schedule",
			DeleteMessage:     "Permanently delete this price schedule? This cannot be undone.",
			ActivateTitle:     "Activate Price Schedule",
			ActivateMessage:   "Activate {{name}}?",
			DeactivateTitle:   "Deactivate Price Schedule",
			DeactivateMessage: "Deactivate {{name}}?",
		},
		Tabs: PriceScheduleTabLabels{
			Info:          "Info",
			PricePlan:     "Plans",
			PricePlanSlug: "",
			ProductPrices: "Product Prices",
		},
		Detail: PriceScheduleDetailLabels{
			Title:                     "Price Schedule",
			DateCreated:               "Date Created",
			DateModified:              "Date Modified",
			NoLocation:                "All locations",
			NoDateEnd:                 "No end date",
			NoDescription:             "—",
			PlansEmptyTitle:           "No Plans",
			PlansEmptyMsg:             "No price plans are linked to this schedule yet.",
			NoDescriptionSubtitle:     "No description provided",
			ProductPriceAdd:           "Add Product Price",
			ProductPriceEdit:          "Edit Product Price",
			ProductPriceDelete:        "Delete Product Price",
			ProductPriceDeleteConfirm: "Remove %s from this plan?",
			ProductPriceEmptyTitle:    "No Product Prices",
			ProductPriceEmptyMsg:      "No product prices have been configured for this plan yet.",
			ProductPriceSection:       "Product Price",
			ProductField:              "Product",

			PlanColumnPlan:     "Plan",
			PlanColumnAmount:   "Amount",
			PlanColumnDuration: "Duration",
			PlanColumnStatus:   "Status",

			PlanView:            "View",
			PlanEdit:            "Edit",
			PlanEditDrawerTitle: "Edit Plan",
			PlanDelete:          "Delete",
			PlanDeleteTitle:     "Delete Plan",
			PlanDeleteMsg:       "Permanently delete %s? This cannot be undone.",
			PlanInUseTooltip:    "In use by active subscriptions",

			PlanAdd:      "Add Plan",
			PlanRequired: "Plan is required",

			ProductPriceColumnProduct:    "Product",
			ProductPriceColumnPrice:      "Price",
			ProductPriceColumnCurrency:   "Currency",
			ProductPriceColumnTreatment:  "Billing",
			ProductPriceColumnEffective:  "Effective",
			BasisBannerPerCycle:          "Each line below is charged every billing cycle.",
			BasisBannerTotalPackage:      "These per-line prices are informational. The package is sold at a flat rate; the total here does not have to match.",
			BasisBannerDerived:           "The package price is the sum of these line prices. Editing a line changes the package total.",
			ProductPriceCatalogSection:   "Catalog line",
			ProductPricePricingSection:   "Pricing",
			ProductPriceEffectiveSection: "Effective dates",
		},
		PlanForm: PriceSchedulePlanFormLabels{
			SectionSchedule:        "Schedule",
			SectionPackage:         "Plan",
			SectionPricing:         "Pricing",
			PriceScheduleField:     "Price Schedule",
			PackageLabel:           "Plan",
			PackagePlaceholder:     "Select a plan...",
			PackageSearch:          "Filter...",
			NameLabel:              "Plan Name",
			NamePlaceholder:        "Enter plan name",
			DescriptionLabel:       "Description",
			DescriptionPlaceholder: "Optional notes for this package",
			AmountLabel:            "Amount",
			AmountPlaceholder:      "0.00",
			CurrencyLabel:          "Currency",
			CurrencyPlaceholder:    "e.g. PHP",
			CurrencyPHP:            "PHP (₱)",
			CurrencyUSD:            "USD ($)",
			DurationLabel:          "Duration",
			UnitLabel:              "Unit",
			ActiveLabel:            "Active",
			SchedulePlaceholder:    "Select a rate card...",
			ScheduleSearch:         "Filter...",
			LocationHintPrefix:     "Location: ",
		},
		Errors: PriceScheduleErrorLabels{
			NotFound:                   "Price schedule not found",
			LoadFailed:                 "Failed to load price schedule",
			Unauthorized:               "You are not authorized to perform this action",
			CreateFailed:               "Failed to create price schedule",
			UpdateFailed:               "Failed to update price schedule",
			DeleteFailed:               "Failed to delete price schedule",
			InUse:                      "This price schedule is in use by active subscriptions and cannot be deleted.",
			PricePlanCreateUnavailable: "Adding a price plan is not available. Please contact support.",
		},
		Filters: PriceScheduleFilterLabels{
			ScopeChipLabel: "Show:",
			ScopeMaster:    "Master",
			ScopeClient:    "Client-specific",
			ScopeAll:       "All",
		},
	}
}

// ClientPackagesLabels holds labels for the client detail "Packages" tab —
// the list of client-scoped Plans for a given client, with the
// "Add custom package" CTA. Mounted from entydad's client detail page via
// a centymo helper view (plan §6.6 option 1).
//
// 2026-04-27 plan-client-scope plan §6.3 / §7.
type ClientPackagesLabels struct {
	TabTitle  string `json:"tabTitle"`
	Empty     string `json:"empty"`
	AddAction string `json:"addAction"`

	// Column headers for the table on the tab.
	ColumnName        string `json:"columnName"`
	ColumnSchedule    string `json:"columnSchedule"`
	ColumnEngagements string `json:"columnEngagements"`
}

// DefaultClientPackagesLabels returns ClientPackagesLabels with sensible English
// defaults. Surfaces the labels for the client-detail Packages tab + the
// "Add custom package" CTA. Centymo owns these labels because the cross-block
// helper that renders this tab lives here (see plan §6.6 option 1).
//
// 2026-04-27 plan-client-scope plan §7.
func DefaultClientPackagesLabels() ClientPackagesLabels {
	return ClientPackagesLabels{
		TabTitle:    "Packages",
		Empty:       "No custom packages yet — every engagement uses a general package.",
		AddAction:   "Add custom package",
		ColumnName:     "Name",
		ColumnSchedule: "Rate card",
		ColumnEngagements: "Engagements",
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
			Name:      "Engagement",
			Client:    "Client",
			Customer:  "Customer",
			Plan:      "Plan",
			StartDate: "Start Date",
			EndDate:   "End Date",
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
			StartTime:                 "Start Time (optional)",
			EndTime:                   "End Time (optional)",
			TimePlaceholder:           "HH:MM",
			Timezone:                  "Timezone",
			Active:                    "Active",
			Notes:                     "Notes",
			NotesPlaceholder:          "Enter notes...",
			CustomerSearchPlaceholder: "Search customers...",
			PlanSearchPlaceholder:     "Search plans...",
			CustomerNoResults:         "No customers found",
			PlanNoResults:             "No plans found",
			Code:                      "Code",
			CodePlaceholder:           "e.g. A3K7PXR",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			CustomerInfo:  "The client this subscription is billed to.",
			PlanInfo:      "The price plan this subscription follows. Determines amount, billing cycle, and any per-product prices.",
			CodeInfo:      "Short reference used on invoices and receipts. Leave blank to auto-generate.",
			StartDateInfo: "First day the subscription is active. Billing cycles are counted from this date.",
			EndDateInfo:   "Last day the subscription is active. Leave blank for open-ended.",
			StartTimeInfo: "Optional time of day in the operator's display timezone. Leave blank for start of day (00:00).",
			EndTimeInfo:   "Optional time of day in the operator's display timezone. Leave blank for end of day (23:59).",
			NotesInfo:     "Internal remarks — shown on detail pages but not on customer-facing documents.",
			// 2026-04-27 plan-client-scope plan §5.1 / §7 — grouped picker headers.
			PlanGroupForClient: "For {{.ClientName}}",
			PlanGroupGeneral:   "General packages",
			// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.1 / §9 —
			// Spawn Jobs toggle on subscription create drawer.
			SpawnJobsSectionTitle: "Operations",
			SpawnJobsToggle:       "Spawn Job(s) on Create",
			SpawnJobsHelpText:     "Disable to start without operational tracking (e.g., advisory retainers).",
			SpawnJobsSummary:      "Spawning {{.JobCount}} Job(s) from {{.TemplateNames}} — includes {{.PhaseCount}} phases, {{.TaskCount}} tasks.",
			SpawnJobsNone:         "No JobTemplate is configured for this Plan. The engagement will start without operational tracking.",
		},
		Actions: SubscriptionActionLabels{
			View:       "View Subscription",
			Edit:       "Edit Subscription",
			Cancel:     "Cancel Subscription",
			Delete:     "Delete",
			Activate:   "Activate",
			Deactivate: "Deactivate",
			// 2026-04-27 plan-client-scope plan §6.5 / §7 — Package tab CTA.
			CustomizePackage: "Customize this package for {{.ClientName}}",
		},
		Bulk: SubscriptionBulkLabels{
			Delete:     "Delete Selected",
			Activate:   "Activate Selected",
			Deactivate: "Deactivate Selected",
		},
		Status: SubscriptionStatusLabels{
			Activate:   "Activate",
			Deactivate: "Deactivate",
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
			Info:         "Information",
			Operations:   "Operations",
			// 2026-04-30 cyclic-subscription-jobs plan §21.2 — flat Jobs tab.
			Jobs:         "Jobs",
			Invoices:     "Invoices",
			History:      "History",
			Attachments:  "Attachments",
			AuditTrail:   "Audit Trail",
			AuditHistory: "History",
		},
		Invoices: SubscriptionInvoicesLabels{
			Title:             "Invoices",
			Empty:             "No invoices yet — click Recognize Revenue to generate the first one.",
			ColumnCode:        "Number",
			ColumnDate:        "Date",
			ColumnAmount:      "Amount",
			ColumnStatus:      "Status",
			RecognizeAction:   "Recognize Revenue",
			RecognizeTitle:    "Recognize Revenue",
			RecognizeSubtitle: "Generate an invoice from this subscription's price plan.",
		},
		Recognize: SubscriptionRecognizeLabels{
			ContextSection:        "Subscription",
			ClientLabel:           "Client",
			PlanLabel:             "Plan / Rate Card",
			QuantityLabel:         "Quantity",
			PeriodSection:         "Billing period",
			PeriodStart:           "Period start",
			PeriodEnd:             "Period end",
			RevenueDate:           "Revenue date",
			LineItemsSection:      "Line items",
			ColumnDescription:     "Description",
			ColumnUnitPrice:       "Unit price",
			ColumnQuantity:        "Qty",
			ColumnLineTotal:       "Line total",
			ColumnTreatment:       "Treatment",
			TotalLabel:            "Total",
			RemoveLine:            "Remove",
			TreatmentRecurring:    "Every cycle",
			TreatmentFirstCycle:   "First cycle only",
			TreatmentUsageBased:   "On use",
			TreatmentOneTime:      "One time",
			NotesLabel:            "Notes",
			NotesPlaceholder:      "Notes are auto-prefixed with the period; append any free-text below.",
			Generate:              "Generate",
			Cancel:                "Cancel",
			CurrencyMismatchError: "Client billing currency ({{.ClientCurrency}}) does not match the rate card ({{.PlanCurrency}}). Update one before generating revenue.",
			IdempotencyError:      "An invoice for this period already exists. Cancel the existing one or pick a different period.",
			IdempotencyExistingLink: "View the existing invoice",
			NoLinesError:            "Cannot create an invoice with no line items. Add a price plan with at least one product, or override at least one line.",
			CycleNotConfiguredWarning: "Plan has no billing cycle configured; defaulting to 1 month.",
			UsageBasedSkippedNotice:   "Usage-based lines were skipped — record them via metering.",
			// 2026-04-27 plan-client-scope plan §7 — surfaced when the
			// active subscription's PricePlan is client-scoped.
			ClientCustomNotice: "This engagement uses a custom package for {{.ClientName}}.",
			// 2026-04-29 milestone-billing plan §5 / Phase E.
			MilestoneSelect:            "Milestone",
			MilestoneSelectPlaceholder: "Select a ready milestone",
			NoReadyMilestone:           "No milestone is ready to bill.",
			MilestoneNotApplicable:     "Milestones are only available on milestone-priced plans.",
			BillAmount:                 "Bill amount",
			LeaveRemainderOpen:         "Partial — leave remainder open",
			CloseShort:                 "Partial — close milestone short",
			PartialReason:              "Reason",
			PartialReasonRequired:      "A reason is required when billing partially.",
			OverBillingRejected:        "Cannot bill: total would exceed milestone amount.",
		},
		Milestone: SubscriptionMilestoneLabels{
			Title:           "Billing Schedule",
			Subtitle:        "Milestone events for this engagement",
			MarkReady:       "Mark Ready",
			Waive:           "Waive",
			ViewInvoice:     "View Invoice",
			StatusPending:   "Pending",
			StatusReady:     "Ready",
			StatusBilled:    "Billed",
			StatusWaived:    "Waived",
			StatusDeferred:  "Deferred",
			StatusCancelled: "Cancelled",
			TotalInvoiced:   "Total Invoiced",
			AmountFull:      "Full amount",
			AmountPartial:   "Partial — {{.Billed}} of {{.Full}}",
		},
		// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.2 / §9 +
		// 2026-04-30 cyclic-subscription-jobs plan §9.1 — cycle accordion.
		Operations: SubscriptionOperationsLabels{
			Title:        "Operational Jobs",
			EmptyTitle:   "No operational tracking",
			EmptyMessage: "This engagement has no Jobs. {{.SpawnAction}} to start tracking work.",
			SpawnAction:  "Spawn Jobs",
			RootJob:      "Root Job",
			ChildJob:     "Child Job",
			PhaseSummary: "{{.Complete}} / {{.Total}} phases complete",
			ViewJobLink:  "View in Operations",

			// Cycle accordion + backfill copy.
			EngagementHeading:     "Engagement (since {{.Started}})",
			CycleHeading:          "Cycle {{.CycleIndex}} — {{.PeriodLabel}}",
			CyclePlaceholder:      "Cycle starts {{.PeriodStart}} — Jobs will spawn at cycle start, or click below to spawn now.",
			CycleSpawnNow:         "Spawn this cycle now",
			CycleStatusPending:    "Pending",
			CycleStatusInProgress: "In progress",
			CycleStatusCompleted:  "Completed",
			CycleStatusOverdue:    "Overdue",
			CycleInvoiceLinked:    "Invoice {{.RevenueCode}} · {{.Status}}",
			CycleNoInvoice:        "Not yet invoiced",
			CycleEmpty:            "No cycles yet",
			BackfillBanner:        "{{.Count}} cycle(s) missing operational tracking. Spawn now to backfill.",
			BackfillCta:           "Backfill missing cycles",
		},
		// 2026-04-30 cyclic-subscription-jobs plan §9.2 — backfill drawer.
		Backfill: SubscriptionBackfillLabels{
			DrawerTitle:       "Backfill cycle Jobs",
			DrawerDescription: "Preview the cycles that will be spawned, then confirm to materialize them in one transaction.",
			PreviewLine:       "Cycle {{.Index}} — {{.PeriodLabel}}",
			CountLabel:        "Cycles to spawn",
			Confirm:           "Spawn {{.Count}} cycle(s)",
			Cancel:            "Cancel",
			MaxWarning:        "Backfill is capped at 24 cycles per request. Reduce the range or run multiple backfills.",
		},
		// 2026-04-30 cyclic-subscription-jobs plan §21.3 — flat Jobs tab.
		Jobs: SubscriptionJobsTabLabels{
			Heading:          "Jobs",
			Empty:            "No Jobs yet — this engagement has no operational tracking.",
			FilterStatus:     "Status",
			FilterType:       "Type",
			FilterAll:        "All",
			SortBy:           "Sort",
			SortByCycle:      "Cycle #",
			ExportCsv:        "Export CSV",
			Summary:          "Showing {{.Visible}} of {{.Total}} Jobs",
			ColumnNumber:     "#",
			ColumnName:       "Job Name",
			ColumnType:       "Type",
			ColumnPhase:      "Phase",
			ColumnStatus:     "Status",
			ColumnPeriod:     "Period",
			TypeEngagement:   "Engagement",
			TypeOnboarding:   "Onboarding",
			TypeCycle:        "Cycle",
			TypeVisit:        "Visit",
			SpawnFailedToast: "Cycle Job spawn failed for {{.Period}} — invoice was created but the operational Job will need a manual retry.",
		},
		Spawn: SubscriptionSpawnLabels{
			Title:             "Spawn Operational Jobs",
			DetectedTemplates: "Detected templates",
			RootTemplate:      "Root template",
			Cancel:            "Cancel",
			Confirm:           "Spawn Jobs",
			SuccessToast:      "Spawned {{.JobCount}} Job(s).",
			Skipped:           "Nothing to spawn — no JobTemplate is linked to this Plan.",
		},
		Confirm: SubscriptionConfirmLabels{
			Cancel:                "Cancel Subscription",
			CancelMessage:         "Are you sure you want to cancel this subscription? This action cannot be undone.",
			Delete:                "Delete Subscription",
			DeleteMessage:         "Are you sure you want to delete this subscription? This action cannot be undone.",
			Activate:              "Activate Subscription",
			ActivateMessage:       "Are you sure you want to activate %s?",
			Deactivate:            "Deactivate Subscription",
			DeactivateMessage:     "Are you sure you want to deactivate %s?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Are you sure you want to activate the selected subscriptions?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Are you sure you want to deactivate the selected subscriptions?",
			BulkDelete:            "Delete Selected",
			BulkDeleteMessage:     "Are you sure you want to delete the selected subscriptions? This action cannot be undone.",
		},
		Errors: SubscriptionErrorLabels{
			PermissionDenied:   "You do not have permission to perform this action",
			InvalidFormData:    "Invalid form data. Please check your inputs and try again.",
			NotFound:           "Subscription not found",
			IDRequired:         "Subscription ID is required",
			NoIDsProvided:      "No subscription IDs provided",
			InvalidStatus:      "Invalid status value",
			NoPermission:       "No permission",
			CannotDelete:       "Cannot delete — this engagement has dependent records",
			InUse:              "Cannot delete — this engagement has dependent records (jobs, revenue, invoices, etc.)",
			PlanClientMismatch: "This package belongs to a different client and cannot be attached here.",
			CustomizeFailed:    "Failed to customize this package. Please try again.",
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
		SelectAll:                common.Table.SelectAll,
		BulkSelectAllPage:        common.Table.BulkSelectAllPage,
		BulkSelectAllAcrossPages: common.Table.BulkSelectAllAcrossPages,
		BulkClearSelection:       common.Table.BulkClearSelection,
		ColumnSortLockedHint:     common.Table.ColumnSortLockedHint,
		SortAscText:              common.Table.SortAscText,
		SortDescText:             common.Table.SortDescText,
		SortAscNumber:            common.Table.SortAscNumber,
		SortDescNumber:           common.Table.SortDescNumber,
		SortAscDate:              common.Table.SortAscDate,
		SortDescDate:             common.Table.SortDescDate,
		SortAscEnum:              common.Table.SortAscEnum,
		SortDescEnum:             common.Table.SortDescEnum,
		Actions:                  common.Table.Actions,
		Prev:                     common.Pagination.Prev,
		Next:                     common.Pagination.Next,
	}
}

// ---------------------------------------------------------------------------
// SupplierContract labels  (P3a)
// ---------------------------------------------------------------------------

// SupplierContractLabels holds all translatable strings for the supplier_contract module.
type SupplierContractLabels struct {
	Page               SupplierContractPageLabels              `json:"page"`
	Columns            SupplierContractColumnLabels            `json:"columns"`
	Tabs               SupplierContractTabLabels               `json:"tabs"`
	Detail             SupplierContractDetailLabels            `json:"detail"`
	Lines              SupplierContractLineLabels              `json:"lines"`
	LinkedPOs          SupplierContractLinkedPOLabels          `json:"linkedPos"`
	LinkedExpenditures SupplierContractLinkedExpenditureLabels `json:"linkedExpenditures"`
	Form               SupplierContractFormLabels              `json:"form"`
	Empty              SupplierContractEmptyLabels             `json:"empty"`
}

type SupplierContractPageLabels struct {
	Heading           string `json:"heading"`
	HeadingDraft      string `json:"headingDraft"`
	HeadingActive     string `json:"headingActive"`
	HeadingExpiring   string `json:"headingExpiring"`
	HeadingExpired    string `json:"headingExpired"`
	HeadingTerminated string `json:"headingTerminated"`
	Caption           string `json:"caption"`
	AddButton         string `json:"addButton"`
	DetailSubtitle    string `json:"detailSubtitle"`
}

type SupplierContractColumnLabels struct {
	Name      string `json:"name"`
	Supplier  string `json:"supplier"`
	Kind      string `json:"kind"`
	Status    string `json:"status"`
	Validity  string `json:"validity"`
	Committed string `json:"committed"`
	Released  string `json:"released"`
	Billed    string `json:"billed"`
	Remaining string `json:"remaining"`
}

type SupplierContractTabLabels struct {
	Info                string `json:"info"`
	Lines               string `json:"lines"`
	LinkedPOs           string `json:"linkedPos"`
	LinkedExpenditures  string `json:"linkedExpenditures"`
	PriceSchedules      string `json:"priceSchedules"`
	Activity            string `json:"activity"`
	ActivityEmpty       string `json:"activityEmpty"`
	PriceSchedulesEmpty string `json:"priceSchedulesEmpty"`
}

type SupplierContractDetailLabels struct {
	InfoSection     string `json:"infoSection"`
	Name            string `json:"name"`
	Kind            string `json:"kind"`
	Status          string `json:"status"`
	Supplier        string `json:"supplier"`
	StartDate       string `json:"startDate"`
	EndDate         string `json:"endDate"`
	AutoRenew       string `json:"autoRenew"`
	Currency        string `json:"currency"`
	CommittedAmount string `json:"committedAmount"`
	ReleasedAmount  string `json:"releasedAmount"`
	BilledAmount    string `json:"billedAmount"`
	RemainingAmount string `json:"remainingAmount"`
	Notes           string `json:"notes"`
}

type SupplierContractLineLabels struct {
	// Column labels
	Description  string `json:"description"`
	LineType     string `json:"lineType"`
	Quantity     string `json:"quantity"`
	UnitPrice    string `json:"unitPrice"`
	Total        string `json:"total"`
	Treatment    string `json:"treatment"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
	AddLine      string `json:"addLine"`

	// Enum label values for treatment
	TreatmentRecurring         string `json:"treatmentRecurring"`
	TreatmentOneTime           string `json:"treatmentOneTime"`
	TreatmentUsageBased        string `json:"treatmentUsageBased"`
	TreatmentMinimumCommitment string `json:"treatmentMinimumCommitment"`

	// Enum label values for line_type
	LineTypeGoods   string `json:"lineTypeGoods"`
	LineTypeService string `json:"lineTypeService"`
	LineTypeExpense string `json:"lineTypeExpense"`

	// Drawer form labels
	FormDescription               string `json:"formDescription"`
	FormDescriptionPlaceholder    string `json:"formDescriptionPlaceholder"`
	FormLineType                  string `json:"formLineType"`
	FormLineTypeInfo              string `json:"formLineTypeInfo"`
	FormTreatment                 string `json:"formTreatment"`
	FormTreatmentInfo             string `json:"formTreatmentInfo"`
	FormProduct                   string `json:"formProduct"`
	FormProductPlaceholder        string `json:"formProductPlaceholder"`
	FormQuantity                  string `json:"formQuantity"`
	FormQuantityInfo              string `json:"formQuantityInfo"`
	FormUnitPrice                 string `json:"formUnitPrice"`
	FormUnitPriceInfo             string `json:"formUnitPriceInfo"`
	FormExpenseAccount            string `json:"formExpenseAccount"`
	FormExpenseAccountPlaceholder string `json:"formExpenseAccountPlaceholder"`
	FormStartDate                 string `json:"formStartDate"`
	FormStartDateHint             string `json:"formStartDateHint"`
	FormEndDate                   string `json:"formEndDate"`
	FormLineNumber                string `json:"formLineNumber"`
}

type SupplierContractLinkedPOLabels struct {
	PONumber     string `json:"poNumber"`
	Status       string `json:"status"`
	TotalAmount  string `json:"totalAmount"`
	OrderDate    string `json:"orderDate"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

type SupplierContractLinkedExpenditureLabels struct {
	Reference    string `json:"reference"`
	Status       string `json:"status"`
	Amount       string `json:"amount"`
	Date         string `json:"date"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// SupplierContractFormLabels holds all form-level labels for the drawer form.
type SupplierContractFormLabels struct {
	// Section headers (5-section parity layout)
	SectionIdentity       string `json:"sectionIdentity"`
	SectionValidity       string `json:"sectionValidity"`
	SectionMoney          string `json:"sectionMoney"`
	SectionCategorization string `json:"sectionCategorization"`
	SectionOthers         string `json:"sectionOthers"`

	// §1 Identity
	Name                       string `json:"name"`
	NamePlaceholder            string `json:"namePlaceholder"`
	NameInfo                   string `json:"nameInfo"`
	ContractNumber             string `json:"contractNumber"`
	ContractNumberPlaceholder  string `json:"contractNumberPlaceholder"`
	Kind                       string `json:"kind"`
	KindInfo                   string `json:"kindInfo"`
	KindSubscription           string `json:"kindSubscription"`
	KindRetainer               string `json:"kindRetainer"`
	KindLease                  string `json:"kindLease"`
	KindUtility                string `json:"kindUtility"`
	KindFramework              string `json:"kindFramework"`
	KindBlanket                string `json:"kindBlanket"`
	KindOneTime                string `json:"kindOneTime"`
	KindOther                  string `json:"kindOther"`
	Supplier                   string `json:"supplier"`
	SupplierPlaceholder        string `json:"supplierPlaceholder"`
	SupplierInfo               string `json:"supplierInfo"`

	// §2 Validity & Recurrence
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	EndDateHint           string `json:"endDateHint"`
	BillingCycleValue     string `json:"billingCycleValue"`
	BillingCycleUnit      string `json:"billingCycleUnit"`
	BillingCycleInfo      string `json:"billingCycleInfo"`
	CycleUnitDay          string `json:"cycleUnitDay"`
	CycleUnitWeek         string `json:"cycleUnitWeek"`
	CycleUnitMonth        string `json:"cycleUnitMonth"`
	CycleUnitYear         string `json:"cycleUnitYear"`
	AutoRenew             string `json:"autoRenew"`
	RenewalNoticeDays     string `json:"renewalNoticeDays"`
	RenewalNoticeDaysHint string `json:"renewalNoticeDaysHint"`

	// §3 Money & Approval
	Currency               string `json:"currency"`
	CurrencyInfo           string `json:"currencyInfo"`
	Status                 string `json:"status"`
	StatusInfo             string `json:"statusInfo"`
	StatusDraft            string `json:"statusDraft"`
	StatusRequested        string `json:"statusRequested"`
	StatusPendingApproval  string `json:"statusPendingApproval"`
	StatusApproved         string `json:"statusApproved"`
	StatusActive           string `json:"statusActive"`
	StatusExpiring         string `json:"statusExpiring"`
	StatusSuspended        string `json:"statusSuspended"`
	StatusExpired          string `json:"statusExpired"`
	StatusTerminated       string `json:"statusTerminated"`
	StatusRejected         string `json:"statusRejected"`
	CommittedAmount        string `json:"committedAmount"`
	CommittedAmountInfo    string `json:"committedAmountInfo"`
	CycleAmount            string `json:"cycleAmount"`
	CycleAmountHint        string `json:"cycleAmountHint"`
	PaymentTerm            string `json:"paymentTerm"`
	PaymentTermPlaceholder string `json:"paymentTermPlaceholder"`
	ApprovedBy             string `json:"approvedBy"`
	ApprovedDate           string `json:"approvedDate"`

	// §4 Categorization
	ExpenditureCategory            string `json:"expenditureCategory"`
	ExpenditureCategoryPlaceholder string `json:"expenditureCategoryPlaceholder"`
	ExpenseAccount                 string `json:"expenseAccount"`
	ExpenseAccountPlaceholder      string `json:"expenseAccountPlaceholder"`
	Location                       string `json:"location"`
	LocationPlaceholder            string `json:"locationPlaceholder"`

	// §5 Others
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	Active           string `json:"active"`

	// Action buttons on detail page
	Edit      string `json:"edit"`
	EditTitle string `json:"editTitle"`
	Approve   string `json:"approve"`
	Terminate string `json:"terminate"`
}

type SupplierContractEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultSupplierContractLabels returns English fallback labels.
// Uses proto-generic naming — tier overrides belong in lyngua JSON.
func DefaultSupplierContractLabels() SupplierContractLabels {
	return SupplierContractLabels{
		Page: SupplierContractPageLabels{
			Heading:           "Supplier Contracts",
			HeadingDraft:      "Draft Contracts",
			HeadingActive:     "Active Contracts",
			HeadingExpiring:   "Expiring Contracts",
			HeadingExpired:    "Expired Contracts",
			HeadingTerminated: "Terminated Contracts",
			Caption:           "Standing agreements with suppliers",
			AddButton:         "New Contract",
			DetailSubtitle:    "Contract details",
		},
		Columns: SupplierContractColumnLabels{
			Name:      "Name",
			Supplier:  "Supplier",
			Kind:      "Kind",
			Status:    "Status",
			Validity:  "Validity",
			Committed: "Committed",
			Released:  "Released",
			Billed:    "Billed",
			Remaining: "Remaining",
		},
		Tabs: SupplierContractTabLabels{
			Info:                "Info",
			Lines:               "Lines",
			LinkedPOs:           "Linked POs",
			LinkedExpenditures:  "Linked Expenditures",
			PriceSchedules:      "Price Schedules",
			Activity:            "Activity",
			ActivityEmpty:       "No activity recorded yet.",
			PriceSchedulesEmpty: "No price schedules yet. Add a schedule to layer multi-year pricing on this contract.",
		},
		Detail: SupplierContractDetailLabels{
			InfoSection:     "Contract Information",
			Name:            "Name",
			Kind:            "Kind",
			Status:          "Status",
			Supplier:        "Supplier",
			StartDate:       "Start Date",
			EndDate:         "End Date",
			AutoRenew:       "Auto Renew",
			Currency:        "Currency",
			CommittedAmount: "Committed Amount",
			ReleasedAmount:  "Released Amount",
			BilledAmount:    "Billed Amount",
			RemainingAmount: "Remaining Amount",
			Notes:           "Notes",
		},
		Lines: SupplierContractLineLabels{
			Description:                "Description",
			LineType:                   "Line Type",
			Quantity:                   "Quantity",
			UnitPrice:                  "Unit Price",
			Total:                      "Total",
			Treatment:                  "Treatment",
			EmptyTitle:                 "No lines yet",
			EmptyMessage:               "Add a line to this contract.",
			AddLine:                    "Add Line",
			TreatmentRecurring:         "Recurring",
			TreatmentOneTime:           "One Time",
			TreatmentUsageBased:        "Usage Based",
			TreatmentMinimumCommitment: "Minimum Commitment",
			LineTypeGoods:              "Goods",
			LineTypeService:            "Service",
			LineTypeExpense:            "Expense",
			FormDescription:            "Description",
			FormDescriptionPlaceholder: "e.g. Cloud hosting — 50 seats",
			FormLineType:               "Line Type",
			FormLineTypeInfo:           "Goods = physical items; Service = intangible; Expense = direct cost",
			FormTreatment:              "Treatment",
			FormTreatmentInfo:          "How this line is billed: recurring, one-time, usage-based, or minimum commitment",
			FormProduct:                "Product",
			FormProductPlaceholder:     "Select a product (optional)",
			FormQuantity:               "Quantity",
			FormQuantityInfo:           "For recurring lines, this is the per-cycle quantity.",
			FormUnitPrice:              "Unit Price",
			FormUnitPriceInfo:          "Amount in centavos ÷ 100 for display.",
			FormExpenseAccount:         "Expense Account",
			FormExpenseAccountPlaceholder: "GL account ID",
			FormStartDate:              "Start Date",
			FormStartDateHint:          "Leave empty to inherit from contract.",
			FormEndDate:                "End Date",
			FormLineNumber:             "Line Number",
		},
		LinkedPOs: SupplierContractLinkedPOLabels{
			PONumber:     "PO Number",
			Status:       "Status",
			TotalAmount:  "Total Amount",
			OrderDate:    "Order Date",
			EmptyTitle:   "No linked purchase orders",
			EmptyMessage: "POs created against this contract will appear here.",
		},
		LinkedExpenditures: SupplierContractLinkedExpenditureLabels{
			Reference:    "Reference",
			Status:       "Status",
			Amount:       "Amount",
			Date:         "Date",
			EmptyTitle:   "No linked expenditures",
			EmptyMessage: "Expenditures linked to this contract will appear here.",
		},
		Form: SupplierContractFormLabels{
			SectionIdentity:            "Identity Details",
			SectionValidity:            "Validity & Recurrence",
			SectionMoney:               "Money & Approval",
			SectionCategorization:      "Categorization",
			SectionOthers:              "Others",
			Name:                       "Contract Name",
			NamePlaceholder:            "e.g. AWS Hosting MSA 2026",
			NameInfo:                   "A short descriptive name for this contract.",
			ContractNumber:             "Contract Number",
			ContractNumberPlaceholder:  "Supplier's reference number",
			Kind:                       "Kind",
			KindInfo:                   "Subscription = recurring time-based; Blanket = quantity-based commitment; Framework = pricing agreement only.",
			KindSubscription:           "Subscription",
			KindRetainer:               "Retainer",
			KindLease:                  "Lease",
			KindUtility:                "Utility",
			KindFramework:              "Framework",
			KindBlanket:                "Blanket",
			KindOneTime:                "One Time",
			KindOther:                  "Other",
			Supplier:                   "Supplier",
			SupplierPlaceholder:        "Select supplier",
			SupplierInfo:               "The vendor or service provider you are committing to.",
			StartDate:                  "Start Date",
			EndDate:                    "End Date",
			EndDateHint:                "Leave empty for open-ended.",
			BillingCycleValue:          "Billing Cycle",
			BillingCycleUnit:           "Cycle Unit",
			BillingCycleInfo:           "How often this contract generates an expenditure (for recurring kinds).",
			CycleUnitDay:               "Day",
			CycleUnitWeek:              "Week",
			CycleUnitMonth:             "Month",
			CycleUnitYear:              "Year",
			AutoRenew:                  "Auto Renew",
			RenewalNoticeDays:          "Renewal Notice (days)",
			RenewalNoticeDaysHint:      "How many days before expiry to send a renewal reminder.",
			Currency:                   "Currency",
			CurrencyInfo:               "ISO 4217 currency code (e.g. PHP, USD).",
			Status:                     "Status",
			StatusInfo:                 "Lifecycle stage. draft → requested → pending_approval → approved → active → expiring/expired/terminated.",
			StatusDraft:                "Draft",
			StatusRequested:            "Requested",
			StatusPendingApproval:      "Pending Approval",
			StatusApproved:             "Approved",
			StatusActive:               "Active",
			StatusExpiring:             "Expiring",
			StatusSuspended:            "Suspended",
			StatusExpired:              "Expired",
			StatusTerminated:           "Terminated",
			StatusRejected:             "Rejected",
			CommittedAmount:            "Committed Amount",
			CommittedAmountInfo:        "Total value committed at signing (centavos). Immutable after approval.",
			CycleAmount:                "Cycle Amount",
			CycleAmountHint:            "Expected per-cycle charge for recurring contracts.",
			PaymentTerm:                "Payment Term",
			PaymentTermPlaceholder:     "Select payment term",
			ApprovedBy:                 "Approved By",
			ApprovedDate:               "Approved Date",
			ExpenditureCategory:            "Expenditure Category",
			ExpenditureCategoryPlaceholder: "Select category",
			ExpenseAccount:             "Expense Account",
			ExpenseAccountPlaceholder:  "GL account ID",
			Location:                   "Location",
			LocationPlaceholder:        "Branch or cost center",
			Notes:                      "Notes",
			NotesPlaceholder:           "Additional notes or context",
			Active:                     "Active",
			Edit:                       "Edit",
			EditTitle:                  "Edit Supplier Contract",
			Approve:                    "Approve",
			Terminate:                  "Terminate",
		},
		Empty: SupplierContractEmptyLabels{
			Title:   "No supplier contracts",
			Message: "Create your first supplier contract to start tracking commitments.",
		},
	}
}

// ---------------------------------------------------------------------------
// ProcurementRequest labels  (P3a)
// ---------------------------------------------------------------------------

// ProcurementRequestLabels holds all translatable strings for the procurement_request module.
type ProcurementRequestLabels struct {
	Page       ProcurementRequestPageLabels      `json:"page"`
	Columns    ProcurementRequestColumnLabels    `json:"columns"`
	Tabs       ProcurementRequestTabLabels       `json:"tabs"`
	Detail     ProcurementRequestDetailLabels    `json:"detail"`
	Lines      ProcurementRequestLineLabels      `json:"lines"`
	SpawnedPOs ProcurementRequestSpawnedPOLabels `json:"spawnedPos"`
	Form       ProcurementRequestFormLabels      `json:"form"`
	Empty      ProcurementRequestEmptyLabels     `json:"empty"`

	// SPS Wave 3 — F1/F2/F3 + CRIT-3 spawn lifecycle
	Filters              ProcurementRequestFilterLabels              `json:"filters"`
	FulfillmentStrategy  ProcurementRequestFulfillmentStrategyLabels `json:"fulfillmentStrategy"`
	FulfillmentMode      ProcurementRequestFulfillmentModeLabels     `json:"fulfillmentMode"`
	FulfillmentModeHints ProcurementRequestFulfillmentModeHintLabels `json:"fulfillmentModeHints"`
	Spawn                ProcurementRequestSpawnLabels               `json:"spawn"`
	PolicyDecision       ProcurementRequestPolicyDecisionLabels      `json:"policyDecision"`
}

// ProcurementRequestFilterLabels — F3 filter chips on the list page.
type ProcurementRequestFilterLabels struct {
	All                    string `json:"all"`
	Status                 string `json:"status"`
	FulfillmentStrategy    string `json:"fulfillmentStrategy"`
	FulfillmentMode        string `json:"fulfillmentMode"`
	AnyStatus              string `json:"anyStatus"`
	AnyFulfillmentStrategy string `json:"anyFulfillmentStrategy"`
	AnyFulfillmentMode     string `json:"anyFulfillmentMode"`
}

// ProcurementRequestFulfillmentStrategyLabels — F3 strategy values for header-level rollup.
type ProcurementRequestFulfillmentStrategyLabels struct {
	UniformOutright  string `json:"uniformOutright"`
	UniformStockable string `json:"uniformStockable"`
	UniformRecurring string `json:"uniformRecurring"`
	UniformPetty     string `json:"uniformPetty"`
	Mixed            string `json:"mixed"`
	Hint             string `json:"hint"`
}

// ProcurementRequestFulfillmentModeLabels — F1 line-level mode values.
type ProcurementRequestFulfillmentModeLabels struct {
	Outright  string `json:"outright"`
	Stockable string `json:"stockable"`
	Recurring string `json:"recurring"`
	Petty     string `json:"petty"`
}

// ProcurementRequestFulfillmentModeHintLabels — F1 short hints rendered under each radio choice.
type ProcurementRequestFulfillmentModeHintLabels struct {
	Outright  string `json:"outright"`
	Stockable string `json:"stockable"`
	Recurring string `json:"recurring"`
	Petty     string `json:"petty"`
}

// ProcurementRequestSpawnLabels — CRIT-3 spawn lifecycle UI strings.
type ProcurementRequestSpawnLabels struct {
	StatusColumn      string `json:"statusColumn"`
	StatusPending     string `json:"statusPending"`
	StatusSpawning    string `json:"statusSpawning"`
	StatusSpawned     string `json:"statusSpawned"`
	StatusFailed      string `json:"statusFailed"`
	StatusUnspecified string `json:"statusUnspecified"`
	ModeColumn        string `json:"modeColumn"`
	SpawnedColumn     string `json:"spawnedColumn"`
	LinkPO            string `json:"linkPo"`
	LinkContract      string `json:"linkContract"`
	LinkExpenditure   string `json:"linkExpenditure"`
	NotApplicable     string `json:"notApplicable"`
	ErrorPrefix       string `json:"errorPrefix"`
	RetryButton       string `json:"retryButton"`
	RetryConfirm      string `json:"retryConfirm"`
}

// ProcurementRequestPolicyDecisionLabels — policy_decision_log section on Info tab.
type ProcurementRequestPolicyDecisionLabels struct {
	SectionTitle string `json:"sectionTitle"`
	Toggle       string `json:"toggle"`
	EmptyMessage string `json:"emptyMessage"`
	Info         string `json:"info"`
}

type ProcurementRequestPageLabels struct {
	Heading                string `json:"heading"`
	HeadingDraft           string `json:"headingDraft"`
	HeadingSubmitted       string `json:"headingSubmitted"`
	HeadingPendingApproval string `json:"headingPendingApproval"`
	HeadingApproved        string `json:"headingApproved"`
	HeadingRejected        string `json:"headingRejected"`
	HeadingFulfilled       string `json:"headingFulfilled"`
	HeadingCancelled       string `json:"headingCancelled"`
	Caption                string `json:"caption"`
	AddButton              string `json:"addButton"`
	DetailSubtitle         string `json:"detailSubtitle"`
}

type ProcurementRequestColumnLabels struct {
	RequestNumber  string `json:"requestNumber"`
	Status         string `json:"status"`
	Requester      string `json:"requester"`
	Supplier       string `json:"supplier"`
	EstimatedTotal string `json:"estimatedTotal"`
	NeededBy       string `json:"neededBy"`
	DateCreated    string `json:"dateCreated"`
}

type ProcurementRequestTabLabels struct {
	Info          string `json:"info"`
	Lines         string `json:"lines"`
	SpawnedPOs    string `json:"spawnedPos"`
	Activity      string `json:"activity"`
	ActivityEmpty string `json:"activityEmpty"`
}

type ProcurementRequestDetailLabels struct {
	InfoSection    string `json:"infoSection"`
	RequestNumber  string `json:"requestNumber"`
	Status         string `json:"status"`
	Requester      string `json:"requester"`
	Supplier       string `json:"supplier"`
	Currency       string `json:"currency"`
	EstimatedTotal string `json:"estimatedTotal"`
	NeededBy       string `json:"neededBy"`
	DateCreated    string `json:"dateCreated"`
	ApprovedBy     string `json:"approvedBy"`
	Justification  string `json:"justification"`
}

type ProcurementRequestLineLabels struct {
	// Column labels
	Description         string `json:"description"`
	LineType            string `json:"lineType"`
	Quantity            string `json:"quantity"`
	EstimatedUnitPrice  string `json:"estimatedUnitPrice"`
	EstimatedTotalPrice string `json:"estimatedTotalPrice"`
	EmptyTitle          string `json:"emptyTitle"`
	EmptyMessage        string `json:"emptyMessage"`
	AddLine             string `json:"addLine"`

	// Enum label values for line_type
	LineTypeGoods   string `json:"lineTypeGoods"`
	LineTypeService string `json:"lineTypeService"`
	LineTypeExpense string `json:"lineTypeExpense"`

	// Drawer form labels
	FormDescription                    string `json:"formDescription"`
	FormDescriptionPlaceholder         string `json:"formDescriptionPlaceholder"`
	FormLineType                       string `json:"formLineType"`
	FormLineTypeInfo                   string `json:"formLineTypeInfo"`
	FormProduct                        string `json:"formProduct"`
	FormProductPlaceholder             string `json:"formProductPlaceholder"`
	FormQuantity                       string `json:"formQuantity"`
	FormQuantityInfo                   string `json:"formQuantityInfo"`
	FormEstimatedUnitPrice             string `json:"formEstimatedUnitPrice"`
	FormEstimatedUnitPriceInfo         string `json:"formEstimatedUnitPriceInfo"`
	FormEstimatedTotalPrice            string `json:"formEstimatedTotalPrice"`
	FormEstimatedTotalPriceHint        string `json:"formEstimatedTotalPriceHint"`
	FormExpenditureCategory            string `json:"formExpenditureCategory"`
	FormExpenditureCategoryPlaceholder string `json:"formExpenditureCategoryPlaceholder"`
	FormLocation                       string `json:"formLocation"`
	FormLocationPlaceholder            string `json:"formLocationPlaceholder"`
	FormLineNumber                     string `json:"formLineNumber"`

	// SPS Wave 3 — F1 fulfillment_mode picker + RECURRING fields + PETTY hint
	FormFulfillmentMode     string `json:"formFulfillmentMode"`
	FormFulfillmentModeInfo string `json:"formFulfillmentModeInfo"`
	FormFulfillmentModeHint string `json:"formFulfillmentModeHint"`

	FormRecurringSection    string `json:"formRecurringSection"`
	FormRecurringCycleValue string `json:"formRecurringCycleValue"`
	FormRecurringCycleUnit  string `json:"formRecurringCycleUnit"`
	FormRecurringTermValue  string `json:"formRecurringTermValue"`
	FormRecurringTermUnit   string `json:"formRecurringTermUnit"`
	FormRecurringCycleHint  string `json:"formRecurringCycleHint"`
	FormRecurringTermHint   string `json:"formRecurringTermHint"`
	FormRecurringUnitDay    string `json:"formRecurringUnitDay"`
	FormRecurringUnitWeek   string `json:"formRecurringUnitWeek"`
	FormRecurringUnitMonth  string `json:"formRecurringUnitMonth"`
	FormRecurringUnitYear   string `json:"formRecurringUnitYear"`

	FormPettyHint string `json:"formPettyHint"`

	// CRIT-3 spawn lifecycle column on the lines table
	ModeBadgeColumn string `json:"modeBadgeColumn"`
}

type ProcurementRequestSpawnedPOLabels struct {
	PONumber     string `json:"poNumber"`
	Status       string `json:"status"`
	TotalAmount  string `json:"totalAmount"`
	OrderDate    string `json:"orderDate"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// ProcurementRequestFormLabels holds all form-level labels for the drawer form.
type ProcurementRequestFormLabels struct {
	// Section headers
	SectionIdentity  string `json:"sectionIdentity"`
	SectionFinancial string `json:"sectionFinancial"`
	SectionApproval  string `json:"sectionApproval"`
	SectionOthers    string `json:"sectionOthers"`

	// §1 Identity
	RequestNumber            string `json:"requestNumber"`
	RequestNumberPlaceholder string `json:"requestNumberPlaceholder"`
	RequestNumberInfo        string `json:"requestNumberInfo"`
	RequesterUser            string `json:"requesterUser"`
	RequesterUserPlaceholder string `json:"requesterUserPlaceholder"`
	Supplier                 string `json:"supplier"`
	SupplierPlaceholder      string `json:"supplierPlaceholder"`
	SupplierHint             string `json:"supplierHint"`
	Location                 string `json:"location"`
	LocationPlaceholder      string `json:"locationPlaceholder"`

	// §2 Financial
	Currency           string `json:"currency"`
	CurrencyInfo       string `json:"currencyInfo"`
	EstimatedTotal     string `json:"estimatedTotal"`
	EstimatedTotalInfo string `json:"estimatedTotalInfo"`

	// §3 Timing & Approval
	NeededByDate          string `json:"neededByDate"`
	NeededByDateInfo      string `json:"neededByDateInfo"`
	Status                string `json:"status"`
	StatusInfo            string `json:"statusInfo"`
	StatusDraft                string `json:"statusDraft"`
	StatusSubmitted            string `json:"statusSubmitted"`
	StatusPendingApproval      string `json:"statusPendingApproval"`
	StatusApproved             string `json:"statusApproved"`
	StatusApprovedPendingSpawn string `json:"statusApprovedPendingSpawn"`
	StatusRejected             string `json:"statusRejected"`
	StatusFulfilled            string `json:"statusFulfilled"`
	StatusCancelled            string `json:"statusCancelled"`
	ApprovedBy            string `json:"approvedBy"`

	// §4 Others
	Justification            string `json:"justification"`
	JustificationPlaceholder string `json:"justificationPlaceholder"`
	Notes                    string `json:"notes"`
	NotesPlaceholder         string `json:"notesPlaceholder"`
	Active                   string `json:"active"`

	// Action buttons
	Edit      string `json:"edit"`
	EditTitle string `json:"editTitle"`
	Submit    string `json:"submit"`
	Approve   string `json:"approve"`
	Reject    string `json:"reject"`
	SpawnPO   string `json:"spawnPo"`
}

type ProcurementRequestEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultProcurementRequestLabels returns English fallback labels.
func DefaultProcurementRequestLabels() ProcurementRequestLabels {
	return ProcurementRequestLabels{
		Page: ProcurementRequestPageLabels{
			Heading:                "Procurement Requests",
			HeadingDraft:           "Draft Requests",
			HeadingSubmitted:       "Submitted Requests",
			HeadingPendingApproval: "Pending Approval",
			HeadingApproved:        "Approved Requests",
			HeadingRejected:        "Rejected Requests",
			HeadingFulfilled:       "Fulfilled Requests",
			HeadingCancelled:       "Cancelled Requests",
			Caption:                "Internal purchase intent records",
			AddButton:              "New Request",
			DetailSubtitle:         "Procurement request details",
		},
		Columns: ProcurementRequestColumnLabels{
			RequestNumber:  "Request #",
			Status:         "Status",
			Requester:      "Requester",
			Supplier:       "Supplier",
			EstimatedTotal: "Estimated Total",
			NeededBy:       "Needed By",
			DateCreated:    "Created",
		},
		Tabs: ProcurementRequestTabLabels{
			Info:          "Info",
			Lines:         "Lines",
			SpawnedPOs:    "Spawned POs",
			Activity:      "Activity",
			ActivityEmpty: "No activity recorded yet.",
		},
		Detail: ProcurementRequestDetailLabels{
			InfoSection:    "Request Information",
			RequestNumber:  "Request Number",
			Status:         "Status",
			Requester:      "Requester",
			Supplier:       "Supplier",
			Currency:       "Currency",
			EstimatedTotal: "Estimated Total",
			NeededBy:       "Needed By",
			DateCreated:    "Created",
			ApprovedBy:     "Approved By",
			Justification:  "Justification",
		},
		Lines: ProcurementRequestLineLabels{
			Description:                        "Description",
			LineType:                           "Line Type",
			Quantity:                           "Quantity",
			EstimatedUnitPrice:                 "Est. Unit Price",
			EstimatedTotalPrice:                "Est. Total",
			EmptyTitle:                         "No lines yet",
			EmptyMessage:                       "Add a line to this request.",
			AddLine:                            "Add Line",
			LineTypeGoods:                      "Goods",
			LineTypeService:                    "Service",
			LineTypeExpense:                    "Expense",
			FormDescription:                    "Description",
			FormDescriptionPlaceholder:         "e.g. 50 laptop units",
			FormLineType:                       "Line Type",
			FormLineTypeInfo:                   "Goods = physical items; Service = intangible; Expense = direct cost",
			FormProduct:                        "Product",
			FormProductPlaceholder:             "Select a product (optional)",
			FormQuantity:                       "Quantity",
			FormQuantityInfo:                   "Number of units requested.",
			FormEstimatedUnitPrice:             "Estimated Unit Price",
			FormEstimatedUnitPriceInfo:         "Best estimate in centavos ÷ 100.",
			FormEstimatedTotalPrice:            "Estimated Total Price",
			FormEstimatedTotalPriceHint:        "Auto-calculated. Override if needed.",
			FormExpenditureCategory:            "Expenditure Category",
			FormExpenditureCategoryPlaceholder: "Select category",
			FormLocation:                       "Location",
			FormLocationPlaceholder:            "Branch or cost center",
			FormLineNumber:                     "Line Number",
			FormFulfillmentMode:                "Fulfillment Mode",
			FormFulfillmentModeInfo:            "How this line will be sourced after approval. Drives the downstream artifact created when the request is approved.",
			FormFulfillmentModeHint:            "Pick one — the spawn cascade dispatches per-line based on this choice.",
			FormRecurringSection:               "Recurring Schedule",
			FormRecurringCycleValue:            "Cycle Every",
			FormRecurringCycleUnit:             "Cycle Unit",
			FormRecurringTermValue:             "Term Length",
			FormRecurringTermUnit:               "Term Unit",
			FormRecurringCycleHint:             "Billing/delivery cadence (e.g. every 1 month).",
			FormRecurringTermHint:              "Total contract horizon (e.g. 24 months).",
			FormRecurringUnitDay:               "Day",
			FormRecurringUnitWeek:              "Week",
			FormRecurringUnitMonth:             "Month",
			FormRecurringUnitYear:              "Year",
			FormPettyHint:                      "Petty mode auto-approves under threshold and posts a direct expenditure. No PO, no contract.",
			ModeBadgeColumn:                    "Mode",
		},
		SpawnedPOs: ProcurementRequestSpawnedPOLabels{
			PONumber:     "PO Number",
			Status:       "Status",
			TotalAmount:  "Total Amount",
			OrderDate:    "Order Date",
			EmptyTitle:   "No purchase orders yet",
			EmptyMessage: "POs spawned from this request will appear here after approval.",
		},
		Form: ProcurementRequestFormLabels{
			SectionIdentity:          "Identity",
			SectionFinancial:         "Financial",
			SectionApproval:          "Timing & Approval",
			SectionOthers:            "Others",
			RequestNumber:            "Request Number",
			RequestNumberPlaceholder: "e.g. PR-2026-001",
			RequestNumberInfo:        "A unique identifier for this procurement request.",
			RequesterUser:            "Requester",
			RequesterUserPlaceholder: "User ID of requester",
			Supplier:                 "Supplier",
			SupplierPlaceholder:      "Select supplier (optional for RFQ)",
			SupplierHint:             "Leave empty if supplier is not yet selected (RFQ flow).",
			Location:                 "Location",
			LocationPlaceholder:      "Branch or cost center",
			Currency:                 "Currency",
			CurrencyInfo:             "ISO 4217 currency code (e.g. PHP, USD).",
			EstimatedTotal:           "Estimated Total",
			EstimatedTotalInfo:       "Best estimate of total spend (centavos ÷ 100 for display).",
			NeededByDate:             "Needed By",
			NeededByDateInfo:         "When the goods or services are required.",
			Status:                   "Status",
			StatusInfo:               "Lifecycle stage. draft → submitted → pending_approval → approved/rejected → fulfilled/cancelled.",
			StatusDraft:                "Draft",
			StatusSubmitted:            "Submitted",
			StatusPendingApproval:      "Pending Approval",
			StatusApproved:             "Approved",
			StatusApprovedPendingSpawn: "Approved — Pending Spawn",
			StatusRejected:             "Rejected",
			StatusFulfilled:            "Fulfilled",
			StatusCancelled:            "Cancelled",
			ApprovedBy:               "Approved By",
			Justification:            "Justification",
			JustificationPlaceholder: "Business reason for this request",
			Notes:                    "Notes",
			NotesPlaceholder:         "Additional notes or context",
			Active:                   "Active",
			Edit:                     "Edit",
			EditTitle:                "Edit Procurement Request",
			Submit:                   "Submit for Approval",
			Approve:                  "Approve",
			Reject:                   "Reject",
			SpawnPO:                  "Create PO",
		},
		Empty: ProcurementRequestEmptyLabels{
			Title:   "No procurement requests",
			Message: "Create a procurement request to start the approval workflow.",
		},
		Filters: ProcurementRequestFilterLabels{
			All:                    "All",
			Status:                 "Status",
			FulfillmentStrategy:    "Fulfillment",
			FulfillmentMode:        "Mode",
			AnyStatus:              "Any Status",
			AnyFulfillmentStrategy: "Any Fulfillment",
			AnyFulfillmentMode:     "Any Mode",
		},
		FulfillmentStrategy: ProcurementRequestFulfillmentStrategyLabels{
			UniformOutright:  "Uniform — Outright",
			UniformStockable: "Uniform — Stockable",
			UniformRecurring: "Uniform — Recurring",
			UniformPetty:     "Uniform — Petty",
			Mixed:            "Mixed Modes",
			Hint:             "Auto-derived from per-line fulfillment modes. Mixed = lines split across multiple modes.",
		},
		FulfillmentMode: ProcurementRequestFulfillmentModeLabels{
			Outright:  "Outright",
			Stockable: "Stockable",
			Recurring: "Recurring",
			Petty:     "Petty",
		},
		FulfillmentModeHints: ProcurementRequestFulfillmentModeHintLabels{
			Outright:  "One-shot purchase. Spawns a single purchase order on approval; no recurrence, no inventory side-effect.",
			Stockable: "Replenishment buy. Spawns a purchase order; received goods credit inventory on receipt.",
			Recurring: "Standing agreement. Spawns a supplier contract on approval; the recurrence engine emits cycle bills.",
			Petty:     "Cash-out. Spawns an expenditure directly against petty cash. No PO, no contract.",
		},
		Spawn: ProcurementRequestSpawnLabels{
			StatusColumn:      "Spawn Status",
			StatusPending:     "Pending",
			StatusSpawning:    "Spawning",
			StatusSpawned:     "Spawned",
			StatusFailed:      "Failed",
			StatusUnspecified: "—",
			ModeColumn:        "Mode",
			SpawnedColumn:     "Spawned Artifact",
			LinkPO:            "View PO line",
			LinkContract:      "View contract",
			LinkExpenditure:   "View expenditure",
			NotApplicable:     "—",
			ErrorPrefix:       "Error",
			RetryButton:       "Retry spawn",
			RetryConfirm:      "Retry spawning the downstream artifact for this line?",
		},
		PolicyDecision: ProcurementRequestPolicyDecisionLabels{
			SectionTitle: "Approval Policy Log",
			Toggle:       "Show / Hide",
			EmptyMessage: "No policy decisions logged yet.",
			Info:         "Audit trail of approval policy decisions taken on this request (auto-approve, escalation, override). Read-only.",
		},
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

// ---------------------------------------------------------------------------
// P3b — Procurement Operations app labels
// (composition surface, no proto entity — mirrors the schedule/cyta pattern)
// ---------------------------------------------------------------------------

// ProcurementLabels holds all translatable strings for the Procurement
// Operations composition app. Populated via lyngua (P4). These keys are
// intentionally generic so they render without overrides when lyngua has not
// yet supplied values.
type ProcurementLabels struct {
	AppLabel              string `json:"app_label"`
	DashboardTitle        string `json:"dashboard_title"`
	PendingApprovalsTitle string `json:"pending_approvals_title"`
	ExpiringTitle         string `json:"expiring_title"`
	VarianceTitle         string `json:"variance_title"`
	RecurrenceTitle       string `json:"recurrence_title"`
	RenewalsTitle         string `json:"renewals_title"`
	UtilizationTitle      string `json:"utilization_title"`
	EmptyRenewals         string `json:"empty_renewals"`
	EmptyVariance         string `json:"empty_variance"`
	EmptyUtilization      string `json:"empty_utilization"`
	EmptyRecurrence       string `json:"empty_recurrence"`
	DaysUntilExpiry       string `json:"days_until_expiry"`
	UtilizationPercent    string `json:"utilization_percent"`
	BudgetPressureLabel   string `json:"budget_pressure_label"`
}

// ---------------------------------------------------------------------------
// SPS P7 — SupplierContractPriceSchedule labels
// (mirrors lyngua/translations/en/general/supplier_contract_price_schedule.json
//  root key "supplierContractPriceSchedule")
// ---------------------------------------------------------------------------

// SupplierContractPriceScheduleLabels holds all translatable strings for the
// supplier_contract_price_schedule + child line views.
type SupplierContractPriceScheduleLabels struct {
	Labels  SupplierContractPriceScheduleNounLabels    `json:"labels"`
	Page    SupplierContractPriceSchedulePageLabels    `json:"page"`
	Buttons SupplierContractPriceScheduleButtonLabels  `json:"buttons"`
	Filters SupplierContractPriceScheduleFilterLabels  `json:"filters"`
	Columns SupplierContractPriceScheduleColumnLabels  `json:"columns"`
	Empty   SupplierContractPriceScheduleEmptyLabels   `json:"empty"`
	Form    SupplierContractPriceScheduleFormLabels    `json:"form"`
	Status  SupplierContractPriceScheduleStatusLabels  `json:"status"`
	Tabs    SupplierContractPriceScheduleTabLabels     `json:"tabs"`
	Lines   SupplierContractPriceScheduleLinesLabels   `json:"lines"`
	Detail  SupplierContractPriceScheduleDetailLabels  `json:"detail"`
	Errors  SupplierContractPriceScheduleErrorLabels   `json:"errors"`
}

type SupplierContractPriceScheduleNounLabels struct {
	Name       string `json:"name"`
	NamePlural string `json:"namePlural"`
	Line       string `json:"line"`
	LinePlural string `json:"linePlural"`
}

type SupplierContractPriceSchedulePageLabels struct {
	Heading           string `json:"heading"`
	Caption           string `json:"caption"`
	HeadingScheduled  string `json:"headingScheduled"`
	HeadingActive     string `json:"headingActive"`
	HeadingSuperseded string `json:"headingSuperseded"`
	HeadingCancelled  string `json:"headingCancelled"`
	TabTitle          string `json:"tabTitle"`
}

type SupplierContractPriceScheduleButtonLabels struct {
	Add       string `json:"add"`
	AddLine   string `json:"addLine"`
	Activate  string `json:"activate"`
	Supersede string `json:"supersede"`
	Cancel    string `json:"cancel"`
}

type SupplierContractPriceScheduleFilterLabels struct {
	All                 string `json:"all"`
	Status              string `json:"status"`
	AnyStatus           string `json:"anyStatus"`
	SupplierContract    string `json:"supplierContract"`
	AnySupplierContract string `json:"anySupplierContract"`
	DateRange           string `json:"dateRange"`
}

type SupplierContractPriceScheduleColumnLabels struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	SupplierContract string `json:"supplierContract"`
	SequenceNumber   string `json:"sequenceNumber"`
	DateStart        string `json:"dateStart"`
	DateEnd          string `json:"dateEnd"`
	Status           string `json:"status"`
	Currency         string `json:"currency"`
	LineCount        string `json:"lineCount"`
	Total            string `json:"total"`
}

type SupplierContractPriceScheduleEmptyLabels struct {
	Title              string `json:"title"`
	Message            string `json:"message"`
	ScheduledTitle     string `json:"scheduledTitle"`
	ScheduledMessage   string `json:"scheduledMessage"`
	ActiveTitle        string `json:"activeTitle"`
	ActiveMessage      string `json:"activeMessage"`
	SupersededTitle    string `json:"supersededTitle"`
	SupersededMessage  string `json:"supersededMessage"`
	CancelledTitle     string `json:"cancelledTitle"`
	CancelledMessage   string `json:"cancelledMessage"`
}

type SupplierContractPriceScheduleFormLabels struct {
	// Section headers
	SectionIdentity  string `json:"sectionIdentity"`
	SectionValidity  string `json:"sectionValidity"`
	SectionScoping   string `json:"sectionScoping"`
	SectionLifecycle string `json:"sectionLifecycle"`
	SectionNotes     string `json:"sectionNotes"`

	// Identity
	Name                  string `json:"name"`
	NamePlaceholder       string `json:"namePlaceholder"`
	NameInfo              string `json:"nameInfo"`
	Description           string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	InternalID            string `json:"internalId"`
	InternalIDPlaceholder string `json:"internalIdPlaceholder"`
	InternalIDInfo        string `json:"internalIdInfo"`

	// Scoping
	SupplierContract       string `json:"supplierContract"`
	SelectSupplierContract string `json:"selectSupplierContract"`
	SupplierContractInfo   string `json:"supplierContractInfo"`

	// Validity
	DateStart          string `json:"dateStart"`
	DateStartInfo      string `json:"dateStartInfo"`
	DateEnd            string `json:"dateEnd"`
	DateEndPlaceholder string `json:"dateEndPlaceholder"`
	DateEndInfo        string `json:"dateEndInfo"`

	// Currency / location
	Currency           string `json:"currency"`
	CurrencyPlaceholder string `json:"currencyPlaceholder"`
	CurrencyInfo       string `json:"currencyInfo"`
	Location           string `json:"location"`
	SelectLocation     string `json:"selectLocation"`
	LocationInfo       string `json:"locationInfo"`

	// Lifecycle
	Status                    string `json:"status"`
	SelectStatus              string `json:"selectStatus"`
	StatusInfo                string `json:"statusInfo"`
	SequenceNumber            string `json:"sequenceNumber"`
	SequenceNumberPlaceholder string `json:"sequenceNumberPlaceholder"`
	SequenceNumberInfo        string `json:"sequenceNumberInfo"`

	// Notes
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	NotesInfo        string `json:"notesInfo"`
}

type SupplierContractPriceScheduleStatusLabels struct {
	Scheduled  string `json:"scheduled"`
	Active     string `json:"active"`
	Superseded string `json:"superseded"`
	Cancelled  string `json:"cancelled"`
}

type SupplierContractPriceScheduleTabLabels struct {
	Info     string `json:"info"`
	Lines    string `json:"lines"`
	Activity string `json:"activity"`
}

type SupplierContractPriceScheduleLinesLabels struct {
	Title              string                                          `json:"title"`
	Empty              string                                          `json:"empty"`
	AddLine            string                                          `json:"addLine"`
	ColumnContractLine string                                          `json:"columnContractLine"`
	ColumnUnitPrice    string                                          `json:"columnUnitPrice"`
	ColumnQuantity     string                                          `json:"columnQuantity"`
	ColumnMinimumAmount string                                         `json:"columnMinimumAmount"`
	ColumnCurrency     string                                          `json:"columnCurrency"`
	ColumnCycleOverride string                                         `json:"columnCycleOverride"`
	LineForm           SupplierContractPriceScheduleLineFormLabels     `json:"lineForm"`
}

type SupplierContractPriceScheduleLineFormLabels struct {
	SectionLink                  string `json:"sectionLink"`
	SectionPricing               string `json:"sectionPricing"`
	SectionCycle                 string `json:"sectionCycle"`
	SupplierContractLine         string `json:"supplierContractLine"`
	SelectSupplierContractLine   string `json:"selectSupplierContractLine"`
	SupplierContractLineInfo     string `json:"supplierContractLineInfo"`
	UnitPrice                    string `json:"unitPrice"`
	UnitPricePlaceholder         string `json:"unitPricePlaceholder"`
	UnitPriceInfo                string `json:"unitPriceInfo"`
	MinimumAmount                string `json:"minimumAmount"`
	MinimumAmountPlaceholder     string `json:"minimumAmountPlaceholder"`
	MinimumAmountInfo            string `json:"minimumAmountInfo"`
	Quantity                     string `json:"quantity"`
	QuantityPlaceholder          string `json:"quantityPlaceholder"`
	QuantityInfo                 string `json:"quantityInfo"`
	Currency                     string `json:"currency"`
	CurrencyPlaceholder          string `json:"currencyPlaceholder"`
	CycleValueOverride           string `json:"cycleValueOverride"`
	CycleValueOverridePlaceholder string `json:"cycleValueOverridePlaceholder"`
	CycleValueOverrideInfo       string `json:"cycleValueOverrideInfo"`
	CycleUnitOverride            string `json:"cycleUnitOverride"`
	CycleUnitOverridePlaceholder string `json:"cycleUnitOverridePlaceholder"`
	CycleUnitOverrideInfo        string `json:"cycleUnitOverrideInfo"`
}

type SupplierContractPriceScheduleDetailLabels struct {
	PageTitle             string `json:"pageTitle"`
	Title                 string `json:"title"`
	InfoSection           string `json:"infoSection"`
	LinesSection          string `json:"linesSection"`
	AuditTrailComingSoon  string `json:"auditTrailComingSoon"`
	AuditEmptyTitle       string `json:"auditEmptyTitle"`
	AuditEmptyMessage     string `json:"auditEmptyMessage"`
}

type SupplierContractPriceScheduleErrorLabels struct {
	PermissionDenied    string `json:"permissionDenied"`
	InvalidFormData     string `json:"invalidFormData"`
	NotFound            string `json:"notFound"`
	IDRequired          string `json:"idRequired"`
	NoPermission        string `json:"noPermission"`
	CannotDelete        string `json:"cannotDelete"`
	InUse               string `json:"inUse"`
	CreationFailed      string `json:"creation_failed"`
	UpdateFailed        string `json:"update_failed"`
	DeletionFailed      string `json:"deletion_failed"`
	ListFailed          string `json:"list_failed"`
	AuthorizationFailed string `json:"authorization_failed"`
	ActivationFailed    string `json:"activation_failed"`
	SupersedeFailed     string `json:"supersede_failed"`
	OverlapDetected     string `json:"overlap_detected"`
	LoadFailed          string `json:"loadFailed"`
}

// DefaultSupplierContractPriceScheduleLabels returns English fallback labels.
// Uses proto-generic naming — tier overrides belong in lyngua JSON.
func DefaultSupplierContractPriceScheduleLabels() SupplierContractPriceScheduleLabels {
	return SupplierContractPriceScheduleLabels{
		Labels: SupplierContractPriceScheduleNounLabels{
			Name:       "Price Schedule",
			NamePlural: "Price Schedules",
			Line:       "Schedule Line",
			LinePlural: "Schedule Lines",
		},
		Page: SupplierContractPriceSchedulePageLabels{
			Heading:           "Contract Price Schedules",
			Caption:           "Date-windowed pricing layered on top of a supplier contract for multi-year escalation",
			HeadingScheduled:  "Scheduled Periods",
			HeadingActive:     "Active Periods",
			HeadingSuperseded: "Superseded Periods",
			HeadingCancelled:  "Cancelled Periods",
			TabTitle:          "Price Schedules",
		},
		Buttons: SupplierContractPriceScheduleButtonLabels{
			Add:       "New Schedule",
			AddLine:   "Add Schedule Line",
			Activate:  "Activate",
			Supersede: "Supersede",
			Cancel:    "Cancel Schedule",
		},
		Filters: SupplierContractPriceScheduleFilterLabels{
			All:                 "All",
			Status:              "Status",
			AnyStatus:           "Any Status",
			SupplierContract:    "Contract",
			AnySupplierContract: "Any Contract",
			DateRange:           "Effective Window",
		},
		Columns: SupplierContractPriceScheduleColumnLabels{
			InternalID:       "ID",
			Name:             "Schedule Name",
			SupplierContract: "Contract",
			SequenceNumber:   "Seq.",
			DateStart:        "Start",
			DateEnd:          "End",
			Status:           "Status",
			Currency:         "Currency",
			LineCount:        "Lines",
			Total:            "Total",
		},
		Empty: SupplierContractPriceScheduleEmptyLabels{
			Title:             "No price schedules yet",
			Message:           "Add a schedule period to layer multi-year pricing onto this contract.",
			ScheduledTitle:    "No upcoming schedules",
			ScheduledMessage:  "Future-dated schedules will appear here once added.",
			ActiveTitle:       "No active schedule",
			ActiveMessage:     "The contract is using header pricing — no schedule is in effect right now.",
			SupersededTitle:   "No past schedules",
			SupersededMessage: "Schedules whose window has passed will be archived here.",
			CancelledTitle:    "No cancelled schedules",
			CancelledMessage:  "Schedules cancelled before activation will appear here.",
		},
		Form: SupplierContractPriceScheduleFormLabels{
			SectionIdentity:           "Schedule Identity",
			SectionValidity:           "Validity Window",
			SectionScoping:            "Scoping",
			SectionLifecycle:          "Lifecycle",
			SectionNotes:              "Notes",
			Name:                      "Schedule Name",
			NamePlaceholder:           "e.g. Year 1 (2026)",
			NameInfo:                  "Human-readable label for this pricing window. Often a year or renewal label.",
			Description:               "Description",
			DescriptionPlaceholder:    "Optional details about this pricing window...",
			InternalID:                "Internal ID",
			InternalIDPlaceholder:     "Auto-generated",
			InternalIDInfo:            "Auto-generated unique identifier.",
			SupplierContract:          "Contract",
			SelectSupplierContract:    "Select contract...",
			SupplierContractInfo:      "The supplier contract this pricing window applies to. Schedules are scoped to a single contract.",
			DateStart:                 "Start Date",
			DateStartInfo:             "When this pricing window takes effect. Window is half-open: start is inclusive.",
			DateEnd:                   "End Date",
			DateEndPlaceholder:        "Leave empty for open-ended",
			DateEndInfo:               "When this pricing window ends. Leave blank for the open-ended last bucket. Window is half-open: end is exclusive.",
			Currency:                  "Currency",
			CurrencyPlaceholder:       "PHP",
			CurrencyInfo:              "ISO 4217 currency for prices in this schedule.",
			Location:                  "Location",
			SelectLocation:            "Select location...",
			LocationInfo:              "Optional location override. Defaults to the parent contract's location.",
			Status:                    "Status",
			SelectStatus:              "Select status...",
			StatusInfo:                "Lifecycle state. New schedules default to Scheduled and progress to Active when their window arrives.",
			SequenceNumber:            "Sequence",
			SequenceNumberPlaceholder: "1",
			SequenceNumberInfo:        "Ordering position within the contract (1, 2, 3...).",
			Notes:                     "Notes",
			NotesPlaceholder:          "Internal notes about this pricing window...",
			NotesInfo:                 "Internal remarks only. Not visible to the supplier.",
		},
		Status: SupplierContractPriceScheduleStatusLabels{
			Scheduled:  "Scheduled",
			Active:     "Active",
			Superseded: "Superseded",
			Cancelled:  "Cancelled",
		},
		Tabs: SupplierContractPriceScheduleTabLabels{
			Info:     "Information",
			Lines:    "Schedule Lines",
			Activity: "Activity",
		},
		Lines: SupplierContractPriceScheduleLinesLabels{
			Title:               "Schedule Lines",
			Empty:               "No schedule lines yet. Add lines to override per-line pricing for this window.",
			AddLine:             "Add Line",
			ColumnContractLine:  "Contract Line",
			ColumnUnitPrice:     "Unit Price",
			ColumnQuantity:      "Qty",
			ColumnMinimumAmount: "Minimum",
			ColumnCurrency:      "Currency",
			ColumnCycleOverride: "Cycle Override",
			LineForm: SupplierContractPriceScheduleLineFormLabels{
				SectionLink:                   "Contract Line",
				SectionPricing:                "Pricing",
				SectionCycle:                  "Cycle Override",
				SupplierContractLine:          "Contract Line",
				SelectSupplierContractLine:    "Select contract line...",
				SupplierContractLineInfo:      "The line whose unit price is overridden during this window.",
				UnitPrice:                     "Unit Price",
				UnitPricePlaceholder:          "0.00",
				UnitPriceInfo:                 "Per-unit price during this window, in the schedule's currency.",
				MinimumAmount:                 "Minimum Amount",
				MinimumAmountPlaceholder:      "0.00",
				MinimumAmountInfo:             "For Minimum Commitment lines: the floor charged per cycle.",
				Quantity:                      "Quantity",
				QuantityPlaceholder:           "0",
				QuantityInfo:                  "Optional committed quantity for blanket or minimum-commitment lines.",
				Currency:                      "Currency",
				CurrencyPlaceholder:           "PHP",
				CycleValueOverride:            "Cycle Value Override",
				CycleValueOverridePlaceholder: "e.g. 1",
				CycleValueOverrideInfo:        "Optional cycle-length override. Most lines inherit the contract cycle.",
				CycleUnitOverride:             "Cycle Unit Override",
				CycleUnitOverridePlaceholder:  "month",
				CycleUnitOverrideInfo:         "Optional cycle-unit override (day, week, month, year).",
			},
		},
		Detail: SupplierContractPriceScheduleDetailLabels{
			PageTitle:            "Schedule Details",
			Title:                "Schedule Detail",
			InfoSection:          "Schedule Information",
			LinesSection:         "Per-Line Pricing",
			AuditTrailComingSoon: "Activity log feature coming soon.",
			AuditEmptyTitle:      "No activity entries",
			AuditEmptyMessage:    "Activity logs for this schedule will appear here.",
		},
		Errors: SupplierContractPriceScheduleErrorLabels{
			PermissionDenied:    "You do not have permission to perform this action.",
			InvalidFormData:     "Invalid form data. Please check your inputs and try again.",
			NotFound:            "Price schedule not found.",
			IDRequired:          "Schedule ID is required.",
			NoPermission:        "No permission.",
			CannotDelete:        "Cannot delete — this schedule is currently active or has dependent lines.",
			InUse:               "Cannot delete — this schedule is referenced by existing records.",
			CreationFailed:      "Schedule creation failed",
			UpdateFailed:        "Schedule update failed",
			DeletionFailed:      "Schedule deletion failed",
			ListFailed:          "Failed to retrieve price schedules",
			AuthorizationFailed: "Authorization failed for price schedules",
			ActivationFailed:    "Schedule activation failed",
			SupersedeFailed:     "Schedule supersede failed",
			OverlapDetected:     "Schedule windows overlap; adjust dates and retry",
			LoadFailed:          "Failed to load price schedule",
		},
	}
}

// ---------------------------------------------------------------------------
// ExpenseRecognition labels  (SPS P10)
// ---------------------------------------------------------------------------

// ExpenseRecognitionLabels holds all translatable strings for the
// expense_recognition module. Loaded from lyngua key root "expenseRecognition".
type ExpenseRecognitionLabels struct {
	Page    ExpenseRecognitionPageLabels    `json:"page"`
	Buttons ExpenseRecognitionButtonLabels  `json:"buttons"`
	Columns ExpenseRecognitionColumnLabels  `json:"columns"`
	Tabs    ExpenseRecognitionTabLabels     `json:"tabs"`
	Detail  ExpenseRecognitionDetailLabels  `json:"detail"`
	Lines   ExpenseRecognitionLineLabels    `json:"lines"`
	Source  ExpenseRecognitionSourceLabels  `json:"source"`
	Status  ExpenseRecognitionStatusLabels  `json:"status"`
	Actions ExpenseRecognitionActionLabels  `json:"actions"`
	Confirm ExpenseRecognitionConfirmLabels `json:"confirm"`
	Empty   ExpenseRecognitionEmptyLabels   `json:"empty"`
	Errors  ExpenseRecognitionErrorLabels   `json:"errors"`
}

type ExpenseRecognitionPageLabels struct {
	Heading         string `json:"heading"`
	Caption         string `json:"caption"`
	HeadingDraft    string `json:"headingDraft"`
	HeadingPosted   string `json:"headingPosted"`
	HeadingReversed string `json:"headingReversed"`
	Dashboard       string `json:"dashboard"`
}

type ExpenseRecognitionButtonLabels struct {
	Add                      string `json:"add"`
	RecognizeFromExpenditure string `json:"recognizeFromExpenditure"`
	RecognizeFromContract    string `json:"recognizeFromContract"`
	Reverse                  string `json:"reverse"`
}

type ExpenseRecognitionColumnLabels struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	RecognitionDate  string `json:"recognitionDate"`
	PeriodStart      string `json:"periodStart"`
	PeriodEnd        string `json:"periodEnd"`
	CycleDate        string `json:"cycleDate"`
	Supplier         string `json:"supplier"`
	SupplierContract string `json:"supplierContract"`
	Expenditure      string `json:"expenditure"`
	Currency         string `json:"currency"`
	TotalAmount      string `json:"totalAmount"`
	Status           string `json:"status"`
	Source           string `json:"source"`
	IdempotencyKey   string `json:"idempotencyKey"`
}

type ExpenseRecognitionTabLabels struct {
	Info     string `json:"info"`
	Lines    string `json:"lines"`
	Source   string `json:"source"`
	Activity string `json:"activity"`
}

type ExpenseRecognitionDetailLabels struct {
	PageTitle             string `json:"pageTitle"`
	Title                 string `json:"title"`
	InfoSection           string `json:"infoSection"`
	SourceSection         string `json:"sourceSection"`
	AuditTrailComingSoon  string `json:"auditTrailComingSoon"`
	AuditEmptyTitle       string `json:"auditEmptyTitle"`
	AuditEmptyMessage     string `json:"auditEmptyMessage"`
}

type ExpenseRecognitionLineLabels struct {
	Description    string `json:"description"`
	Quantity       string `json:"quantity"`
	UnitAmount     string `json:"unitAmount"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Product        string `json:"product"`
	ExpenseAccount string `json:"expenseAccount"`
	EmptyTitle     string `json:"emptyTitle"`
	EmptyMessage   string `json:"emptyMessage"`
	AddLine        string `json:"addLine"`

	// Drawer form labels
	FormDescription            string `json:"formDescription"`
	FormDescriptionPlaceholder string `json:"formDescriptionPlaceholder"`
	FormQuantity               string `json:"formQuantity"`
	FormUnitAmount             string `json:"formUnitAmount"`
	FormAmount                 string `json:"formAmount"`
	FormCurrency               string `json:"formCurrency"`
}

type ExpenseRecognitionSourceLabels struct {
	Recurrence  string `json:"recurrence"`
	Expenditure string `json:"expenditure"`
	Manual      string `json:"manual"`
	Reversal    string `json:"reversal"`
}

type ExpenseRecognitionStatusLabels struct {
	Draft    string `json:"draft"`
	Posted   string `json:"posted"`
	Reversed string `json:"reversed"`
}

type ExpenseRecognitionActionLabels struct {
	View                     string `json:"view"`
	Edit                     string `json:"edit"`
	Delete                   string `json:"delete"`
	Reverse                  string `json:"reverse"`
	RecognizeFromExpenditure string `json:"recognizeFromExpenditure"`
	RecognizeFromContract    string `json:"recognizeFromContract"`
	NoPermission             string `json:"noPermission"`
}

type ExpenseRecognitionConfirmLabels struct {
	Delete         string `json:"delete"`
	DeleteMessage  string `json:"deleteMessage"`
	Reverse        string `json:"reverse"`
	ReverseMessage string `json:"reverseMessage"`
}

type ExpenseRecognitionEmptyLabels struct {
	Title           string `json:"title"`
	Message         string `json:"message"`
	DraftTitle      string `json:"draftTitle"`
	DraftMessage    string `json:"draftMessage"`
	PostedTitle     string `json:"postedTitle"`
	PostedMessage   string `json:"postedMessage"`
	ReversedTitle   string `json:"reversedTitle"`
	ReversedMessage string `json:"reversedMessage"`
}

type ExpenseRecognitionErrorLabels struct {
	PermissionDenied      string `json:"permissionDenied"`
	InvalidFormData       string `json:"invalidFormData"`
	NotFound              string `json:"notFound"`
	IDRequired            string `json:"idRequired"`
	NoPermission          string `json:"noPermission"`
	CreationFailed        string `json:"creation_failed"`
	UpdateFailed          string `json:"update_failed"`
	DeletionFailed        string `json:"deletion_failed"`
	ListFailed            string `json:"list_failed"`
	ReverseFailed         string `json:"reverse_failed"`
	IdempotencyCollision  string `json:"idempotency_collision"`
	LoadFailed            string `json:"load_failed"`
}

// DefaultExpenseRecognitionLabels returns English fallback labels.
// Tier overrides belong in lyngua JSON.
func DefaultExpenseRecognitionLabels() ExpenseRecognitionLabels {
	return ExpenseRecognitionLabels{
		Page: ExpenseRecognitionPageLabels{
			Heading:         "Expense Recognition",
			Caption:         "Period in which a supplier cost is recognized",
			HeadingDraft:    "Draft Recognitions",
			HeadingPosted:   "Posted Recognitions",
			HeadingReversed: "Reversed Recognitions",
			Dashboard:       "Expense Recognition Dashboard",
		},
		Buttons: ExpenseRecognitionButtonLabels{
			Add:                      "New Recognition",
			RecognizeFromExpenditure: "Recognize from Bill",
			RecognizeFromContract:    "Recognize from Contract",
			Reverse:                  "Reverse Recognition",
		},
		Columns: ExpenseRecognitionColumnLabels{
			InternalID:       "ID",
			Name:             "Name",
			RecognitionDate:  "Recognition Date",
			PeriodStart:      "Period Start",
			PeriodEnd:        "Period End",
			CycleDate:        "Cycle",
			Supplier:         "Supplier",
			SupplierContract: "Contract",
			Expenditure:      "Source Bill",
			Currency:         "Currency",
			TotalAmount:      "Amount",
			Status:           "Status",
			Source:           "Source",
			IdempotencyKey:   "Idempotency Key",
		},
		Tabs: ExpenseRecognitionTabLabels{
			Info:     "Information",
			Lines:    "Recognition Lines",
			Source:   "Source",
			Activity: "Activity",
		},
		Detail: ExpenseRecognitionDetailLabels{
			PageTitle:            "Recognition Details",
			Title:                "Recognition Detail",
			InfoSection:          "Recognition Information",
			SourceSection:        "Source",
			AuditTrailComingSoon: "Activity log feature coming soon.",
			AuditEmptyTitle:      "No activity entries",
			AuditEmptyMessage:    "Activity logs for this recognition will appear here.",
		},
		Lines: ExpenseRecognitionLineLabels{
			Description:                "Description",
			Quantity:                   "Quantity",
			UnitAmount:                 "Unit Amount",
			Amount:                     "Amount",
			Currency:                   "Currency",
			Product:                    "Product",
			ExpenseAccount:             "Expense Account",
			EmptyTitle:                 "No recognition lines",
			EmptyMessage:               "Lines breaking down this recognition will appear here.",
			AddLine:                    "Add Line",
			FormDescription:            "Description",
			FormDescriptionPlaceholder: "e.g. Cloud hosting — May 2026",
			FormQuantity:               "Quantity",
			FormUnitAmount:             "Unit Amount",
			FormAmount:                 "Amount",
			FormCurrency:               "Currency",
		},
		Source: ExpenseRecognitionSourceLabels{
			Recurrence:  "Recurrence Engine",
			Expenditure: "From Bill",
			Manual:      "Manual",
			Reversal:    "Reversal",
		},
		Status: ExpenseRecognitionStatusLabels{
			Draft:    "Draft",
			Posted:   "Posted",
			Reversed: "Reversed",
		},
		Actions: ExpenseRecognitionActionLabels{
			View:                     "View Recognition",
			Edit:                     "Edit Recognition",
			Delete:                   "Delete Recognition",
			Reverse:                  "Reverse",
			RecognizeFromExpenditure: "Recognize from Bill",
			RecognizeFromContract:    "Recognize from Contract",
			NoPermission:             "No permission",
		},
		Confirm: ExpenseRecognitionConfirmLabels{
			Delete:         "Delete Recognition",
			DeleteMessage:  "Are you sure you want to delete this recognition? Only Draft recognitions can be deleted.",
			Reverse:        "Reverse Recognition",
			ReverseMessage: "Reversing creates a counter-entry and marks this recognition as Reversed. Continue?",
		},
		Empty: ExpenseRecognitionEmptyLabels{
			Title:           "No recognitions yet",
			Message:         "Recognize an expense from a posted bill or directly from a contract cycle.",
			DraftTitle:      "No draft recognitions",
			DraftMessage:    "Draft recognitions will appear here.",
			PostedTitle:     "No posted recognitions",
			PostedMessage:   "Posted recognitions will appear here.",
			ReversedTitle:   "No reversed recognitions",
			ReversedMessage: "Reversed recognitions will appear here.",
		},
		Errors: ExpenseRecognitionErrorLabels{
			PermissionDenied:     "You do not have permission to perform this action.",
			InvalidFormData:      "Invalid form data. Please check your inputs and try again.",
			NotFound:             "Recognition not found.",
			IDRequired:           "Recognition ID is required.",
			NoPermission:         "No permission.",
			CreationFailed:       "Recognition creation failed",
			UpdateFailed:         "Recognition update failed",
			DeletionFailed:       "Recognition deletion failed",
			ListFailed:           "Failed to retrieve recognitions",
			ReverseFailed:        "Recognition reversal failed",
			IdempotencyCollision: "A recognition for this source and period already exists",
			LoadFailed:           "Failed to load recognition",
		},
	}
}

// ---------------------------------------------------------------------------
// AccruedExpense labels  (SPS P10)
// ---------------------------------------------------------------------------

// AccruedExpenseLabels holds all translatable strings for the accrued_expense
// + accrued_expense_settlement modules. Loaded from lyngua key root
// "accruedExpense" with settlement subkeys merged in via composition.
type AccruedExpenseLabels struct {
	Page        AccruedExpensePageLabels        `json:"page"`
	Buttons     AccruedExpenseButtonLabels      `json:"buttons"`
	Columns     AccruedExpenseColumnLabels      `json:"columns"`
	Tabs        AccruedExpenseTabLabels         `json:"tabs"`
	Detail      AccruedExpenseDetailLabels      `json:"detail"`
	Settlements AccruedExpenseSettlementLabels  `json:"settlements"`
	Form        AccruedExpenseFormLabels        `json:"form"`
	Status      AccruedExpenseStatusLabels      `json:"status"`
	Actions     AccruedExpenseActionLabels      `json:"actions"`
	Confirm     AccruedExpenseConfirmLabels     `json:"confirm"`
	Balances    AccruedExpenseBalanceLabels     `json:"balances"`
	Empty       AccruedExpenseEmptyLabels       `json:"empty"`
	Errors      AccruedExpenseErrorLabels       `json:"errors"`
}

type AccruedExpensePageLabels struct {
	Heading            string `json:"heading"`
	Caption            string `json:"caption"`
	HeadingOutstanding string `json:"headingOutstanding"`
	HeadingPartial     string `json:"headingPartial"`
	HeadingSettled     string `json:"headingSettled"`
	HeadingReversed    string `json:"headingReversed"`
	Dashboard          string `json:"dashboard"`
}

type AccruedExpenseButtonLabels struct {
	Add                string `json:"add"`
	AccrueFromContract string `json:"accrueFromContract"`
	Settle             string `json:"settle"`
	Reverse            string `json:"reverse"`
	AddSettlement      string `json:"addSettlement"`
}

type AccruedExpenseColumnLabels struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	Supplier         string `json:"supplier"`
	SupplierContract string `json:"supplierContract"`
	RecognitionDate  string `json:"recognitionDate"`
	PeriodStart      string `json:"periodStart"`
	PeriodEnd        string `json:"periodEnd"`
	CycleDate        string `json:"cycleDate"`
	Currency         string `json:"currency"`
	AccruedAmount    string `json:"accruedAmount"`
	SettledAmount    string `json:"settledAmount"`
	RemainingAmount  string `json:"remainingAmount"`
	Status           string `json:"status"`
}

type AccruedExpenseTabLabels struct {
	Info        string `json:"info"`
	Settlements string `json:"settlements"`
	Source      string `json:"source"`
	Activity    string `json:"activity"`
}

type AccruedExpenseDetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	Title                string `json:"title"`
	InfoSection          string `json:"infoSection"`
	SettlementsSection   string `json:"settlementsSection"`
	SourceSection        string `json:"sourceSection"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
}

type AccruedExpenseSettlementLabels struct {
	Expenditure        string `json:"expenditure"`
	AmountSettled      string `json:"amountSettled"`
	Currency           string `json:"currency"`
	FxRate             string `json:"fxRate"`
	FxAdjustmentAmount string `json:"fxAdjustmentAmount"`
	SettledAt          string `json:"settledAt"`
	Reversal           string `json:"reversal"`
	EmptyTitle         string `json:"emptyTitle"`
	EmptyMessage       string `json:"emptyMessage"`
	AddSettlement      string `json:"addSettlement"`

	// Drawer form labels
	FormExpenditure        string `json:"formExpenditure"`
	FormExpenditurePlaceholder string `json:"formExpenditurePlaceholder"`
	FormAmountSettled      string `json:"formAmountSettled"`
	FormCurrency           string `json:"formCurrency"`
	FormFxRate             string `json:"formFxRate"`
	FormFxRateInfo         string `json:"formFxRateInfo"`
	FormReversalReason     string `json:"formReversalReason"`
}

type AccruedExpenseFormLabels struct {
	// Section headers
	SectionIdentity   string `json:"sectionIdentity"`
	SectionSource     string `json:"sectionSource"`
	SectionPeriod     string `json:"sectionPeriod"`
	SectionMoney      string `json:"sectionMoney"`
	SectionAccounting string `json:"sectionAccounting"`
	SectionLifecycle  string `json:"sectionLifecycle"`
	SectionNotes      string `json:"sectionNotes"`

	// §1 Identity
	Name                  string `json:"name"`
	NamePlaceholder       string `json:"namePlaceholder"`
	NameInfo              string `json:"nameInfo"`
	Description           string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	InternalID            string `json:"internalId"`
	InternalIDPlaceholder string `json:"internalIdPlaceholder"`

	// §2 Source
	SupplierContract        string `json:"supplierContract"`
	SelectSupplierContract  string `json:"selectSupplierContract"`
	SupplierContractInfo    string `json:"supplierContractInfo"`
	Supplier                string `json:"supplier"`
	SelectSupplier          string `json:"selectSupplier"`
	SupplierInfo            string `json:"supplierInfo"`

	// §3 Period
	RecognitionDate     string `json:"recognitionDate"`
	RecognitionDateInfo string `json:"recognitionDateInfo"`
	PeriodStart         string `json:"periodStart"`
	PeriodStartInfo     string `json:"periodStartInfo"`
	PeriodEnd           string `json:"periodEnd"`
	PeriodEndInfo       string `json:"periodEndInfo"`
	CycleDate           string `json:"cycleDate"`
	CycleDatePlaceholder string `json:"cycleDatePlaceholder"`
	CycleDateInfo       string `json:"cycleDateInfo"`

	// §4 Money
	Currency             string `json:"currency"`
	CurrencyPlaceholder  string `json:"currencyPlaceholder"`
	CurrencyInfo         string `json:"currencyInfo"`
	AccruedAmount        string `json:"accruedAmount"`
	AccruedAmountPlaceholder string `json:"accruedAmountPlaceholder"`
	AccruedAmountInfo    string `json:"accruedAmountInfo"`
	SettledAmount        string `json:"settledAmount"`
	SettledAmountInfo    string `json:"settledAmountInfo"`
	RemainingAmount      string `json:"remainingAmount"`
	RemainingAmountInfo  string `json:"remainingAmountInfo"`

	// §5 Lifecycle
	Status              string `json:"status"`
	SelectStatus        string `json:"selectStatus"`
	StatusInfo          string `json:"statusInfo"`
	StatusOutstanding   string `json:"statusOutstanding"`
	StatusPartial       string `json:"statusPartial"`
	StatusSettled       string `json:"statusSettled"`
	StatusReversed      string `json:"statusReversed"`

	// §6 Accounting
	ExpenseAccount       string `json:"expenseAccount"`
	SelectExpenseAccount string `json:"selectExpenseAccount"`
	ExpenseAccountInfo   string `json:"expenseAccountInfo"`
	AccrualAccount       string `json:"accrualAccount"`
	SelectAccrualAccount string `json:"selectAccrualAccount"`
	AccrualAccountInfo   string `json:"accrualAccountInfo"`

	// §7 Notes
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	NotesInfo        string `json:"notesInfo"`

	// Buttons
	Edit      string `json:"edit"`
	EditTitle string `json:"editTitle"`
	Active    string `json:"active"`
}

type AccruedExpenseStatusLabels struct {
	Outstanding string `json:"outstanding"`
	Partial     string `json:"partial"`
	Settled     string `json:"settled"`
	Reversed    string `json:"reversed"`
}

type AccruedExpenseActionLabels struct {
	View               string `json:"view"`
	Edit               string `json:"edit"`
	Delete             string `json:"delete"`
	AccrueFromContract string `json:"accrueFromContract"`
	Settle             string `json:"settle"`
	Reverse            string `json:"reverse"`
	AddSettlement      string `json:"addSettlement"`
	NoPermission       string `json:"noPermission"`
}

type AccruedExpenseConfirmLabels struct {
	Delete         string `json:"delete"`
	DeleteMessage  string `json:"deleteMessage"`
	Settle         string `json:"settle"`
	SettleMessage  string `json:"settleMessage"`
	Reverse        string `json:"reverse"`
	ReverseMessage string `json:"reverseMessage"`
}

type AccruedExpenseBalanceLabels struct {
	Title       string `json:"title"`
	Accrued     string `json:"accrued"`
	Settled     string `json:"settled"`
	Remaining   string `json:"remaining"`
	Utilization string `json:"utilization"`
}

type AccruedExpenseEmptyLabels struct {
	Title              string `json:"title"`
	Message            string `json:"message"`
	OutstandingTitle   string `json:"outstandingTitle"`
	OutstandingMessage string `json:"outstandingMessage"`
	PartialTitle       string `json:"partialTitle"`
	PartialMessage     string `json:"partialMessage"`
	SettledTitle       string `json:"settledTitle"`
	SettledMessage     string `json:"settledMessage"`
	ReversedTitle      string `json:"reversedTitle"`
	ReversedMessage    string `json:"reversedMessage"`
}

type AccruedExpenseErrorLabels struct {
	PermissionDenied      string `json:"permissionDenied"`
	InvalidFormData       string `json:"invalidFormData"`
	NotFound              string `json:"notFound"`
	IDRequired            string `json:"idRequired"`
	NoPermission          string `json:"noPermission"`
	CreationFailed        string `json:"creation_failed"`
	UpdateFailed          string `json:"update_failed"`
	DeletionFailed        string `json:"deletion_failed"`
	ListFailed            string `json:"list_failed"`
	SettleFailed          string `json:"settle_failed"`
	ReverseFailed         string `json:"reverse_failed"`
	BalanceDrift          string `json:"balance_drift"`
	LoadFailed            string `json:"load_failed"`
}

// DefaultAccruedExpenseLabels returns English fallback labels.
func DefaultAccruedExpenseLabels() AccruedExpenseLabels {
	return AccruedExpenseLabels{
		Page: AccruedExpensePageLabels{
			Heading:            "Accrued Expenses",
			Caption:            "Recognized supplier obligations awaiting the actual bill",
			HeadingOutstanding: "Outstanding Accruals",
			HeadingPartial:     "Partially Settled",
			HeadingSettled:     "Settled Accruals",
			HeadingReversed:    "Reversed Accruals",
			Dashboard:          "Accrued Expense Dashboard",
		},
		Buttons: AccruedExpenseButtonLabels{
			Add:                "New Accrual",
			AccrueFromContract: "Accrue from Contract",
			Settle:             "Settle",
			Reverse:            "Reverse",
			AddSettlement:      "Record Settlement",
		},
		Columns: AccruedExpenseColumnLabels{
			InternalID:       "ID",
			Name:             "Name",
			Supplier:         "Supplier",
			SupplierContract: "Contract",
			RecognitionDate:  "Recognition Date",
			PeriodStart:      "Period Start",
			PeriodEnd:        "Period End",
			CycleDate:        "Cycle",
			Currency:         "Currency",
			AccruedAmount:    "Accrued",
			SettledAmount:    "Settled",
			RemainingAmount:  "Remaining",
			Status:           "Status",
		},
		Tabs: AccruedExpenseTabLabels{
			Info:        "Information",
			Settlements: "Settlements",
			Source:      "Source",
			Activity:    "Activity",
		},
		Detail: AccruedExpenseDetailLabels{
			PageTitle:            "Accrual Details",
			Title:                "Accrual Detail",
			InfoSection:          "Accrual Information",
			SettlementsSection:   "Settlements",
			SourceSection:        "Source",
			AuditTrailComingSoon: "Activity log feature coming soon.",
			AuditEmptyTitle:      "No activity entries",
			AuditEmptyMessage:    "Activity logs for this accrual will appear here.",
		},
		Settlements: AccruedExpenseSettlementLabels{
			Expenditure:                "Bill",
			AmountSettled:              "Amount Settled",
			Currency:                   "Currency",
			FxRate:                     "FX Rate",
			FxAdjustmentAmount:         "FX Adjustment",
			SettledAt:                  "Settled At",
			Reversal:                   "Reversal",
			EmptyTitle:                 "No settlements yet",
			EmptyMessage:               "Settlements applied against this accrual will appear here.",
			AddSettlement:              "Record Settlement",
			FormExpenditure:            "Bill",
			FormExpenditurePlaceholder: "Select bill...",
			FormAmountSettled:          "Amount Settled",
			FormCurrency:               "Currency",
			FormFxRate:                 "FX Rate",
			FormFxRateInfo:             "Bill currency to accrual currency conversion rate.",
			FormReversalReason:         "Reversal Reason",
		},
		Form: AccruedExpenseFormLabels{
			SectionIdentity:          "Accrual Identity",
			SectionSource:            "Source",
			SectionPeriod:            "Period",
			SectionMoney:             "Money",
			SectionAccounting:        "Accounting",
			SectionLifecycle:         "Lifecycle",
			SectionNotes:             "Notes",
			Name:                     "Accrual Name",
			NamePlaceholder:          "e.g. Utilities — May 2026 (estimate)",
			NameInfo:                 "Descriptive label for the accrued obligation.",
			Description:              "Description",
			DescriptionPlaceholder:   "Optional details about this accrual...",
			InternalID:               "Internal ID",
			InternalIDPlaceholder:    "Auto-generated",
			SupplierContract:         "Source Contract",
			SelectSupplierContract:   "Select contract...",
			SupplierContractInfo:     "The contract whose cycle is being accrued.",
			Supplier:                 "Supplier",
			SelectSupplier:           "Select supplier...",
			SupplierInfo:             "Supplier the obligation is owed to.",
			RecognitionDate:          "Recognition Date",
			RecognitionDateInfo:      "Date the obligation is booked into the period.",
			PeriodStart:              "Period Start",
			PeriodStartInfo:          "Start of the period the accrual covers.",
			PeriodEnd:                "Period End",
			PeriodEndInfo:            "End of the period the accrual covers.",
			CycleDate:                "Cycle Date",
			CycleDatePlaceholder:     "YYYY-MM-DD",
			CycleDateInfo:            "Cycle bucket used for idempotency.",
			Currency:                 "Currency",
			CurrencyPlaceholder:      "PHP",
			CurrencyInfo:             "ISO 4217 currency for the accrual.",
			AccruedAmount:            "Accrued Amount",
			AccruedAmountPlaceholder: "0.00",
			AccruedAmountInfo:        "Estimated obligation for the period.",
			SettledAmount:            "Settled Amount",
			SettledAmountInfo:        "Sum of settlements applied so far. Read-only.",
			RemainingAmount:          "Remaining",
			RemainingAmountInfo:      "Accrued minus settled. Read-only.",
			Status:                   "Status",
			SelectStatus:             "Select status...",
			StatusInfo:               "Lifecycle state.",
			StatusOutstanding:        "Outstanding",
			StatusPartial:            "Partially Settled",
			StatusSettled:            "Settled",
			StatusReversed:           "Reversed",
			ExpenseAccount:           "Expense Account",
			SelectExpenseAccount:     "Select expense account...",
			ExpenseAccountInfo:       "GL account to debit on recognition.",
			AccrualAccount:           "Accrual Account",
			SelectAccrualAccount:     "Select accrual account...",
			AccrualAccountInfo:       "GL account to credit while outstanding.",
			Notes:                    "Notes",
			NotesPlaceholder:         "Internal notes about this accrual...",
			NotesInfo:                "Internal remarks only.",
			Edit:                     "Edit",
			EditTitle:                "Edit Accrual",
			Active:                   "Active",
		},
		Status: AccruedExpenseStatusLabels{
			Outstanding: "Outstanding",
			Partial:     "Partially Settled",
			Settled:     "Settled",
			Reversed:    "Reversed",
		},
		Actions: AccruedExpenseActionLabels{
			View:               "View Accrual",
			Edit:               "Edit Accrual",
			Delete:             "Delete Accrual",
			AccrueFromContract: "Accrue from Contract",
			Settle:             "Settle Accrual",
			Reverse:            "Reverse Accrual",
			AddSettlement:      "Record Settlement",
			NoPermission:       "No permission",
		},
		Confirm: AccruedExpenseConfirmLabels{
			Delete:         "Delete Accrual",
			DeleteMessage:  "Are you sure you want to delete this accrual?",
			Settle:         "Settle Accrual",
			SettleMessage:  "Apply a settlement against this accrual?",
			Reverse:        "Reverse Accrual",
			ReverseMessage: "Reversing flips this accrual to Reversed and posts a reversing journal. Continue?",
		},
		Balances: AccruedExpenseBalanceLabels{
			Title:       "Settlement Progress",
			Accrued:     "Accrued",
			Settled:     "Settled",
			Remaining:   "Remaining",
			Utilization: "Settlement Progress",
		},
		Empty: AccruedExpenseEmptyLabels{
			Title:              "No accrued expenses yet",
			Message:            "Accrue from a contract cycle when the period closes before the supplier bill arrives.",
			OutstandingTitle:   "No outstanding accruals",
			OutstandingMessage: "Outstanding accruals waiting for a supplier bill will appear here.",
			PartialTitle:       "No partially settled accruals",
			PartialMessage:     "Accruals with at least one settlement will appear here.",
			SettledTitle:       "No settled accruals",
			SettledMessage:     "Fully reconciled accruals will appear here.",
			ReversedTitle:      "No reversed accruals",
			ReversedMessage:    "Reversed accruals will appear here.",
		},
		Errors: AccruedExpenseErrorLabels{
			PermissionDenied: "You do not have permission to perform this action.",
			InvalidFormData:  "Invalid form data. Please check your inputs and try again.",
			NotFound:         "Accrued expense not found.",
			IDRequired:       "Accrual ID is required.",
			NoPermission:     "No permission.",
			CreationFailed:   "Accrual creation failed",
			UpdateFailed:     "Accrual update failed",
			DeletionFailed:   "Accrual deletion failed",
			ListFailed:       "Failed to retrieve accruals",
			SettleFailed:     "Settlement recording failed",
			ReverseFailed:    "Accrual reversal failed",
			BalanceDrift:     "Accrual balance drift detected; recompute required",
			LoadFailed:       "Failed to load accrual",
		},
	}
}
