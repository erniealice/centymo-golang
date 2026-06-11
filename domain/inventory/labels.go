package inventory

import (
	productdom "github.com/erniealice/centymo-golang/domain/product"
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
	TrackingMode productdom.TrackingModeLabels  `json:"trackingMode"`
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
	SerialNumber     string `json:"serialNumber"`
	IMEI             string `json:"imei"`
	SerialStatus     string `json:"serialStatus"`
	WarrantyEnd      string `json:"warrantyEnd"`
	PurchaseOrder    string `json:"purchaseOrder"`
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
	AuditHistory string `json:"auditHistory"`
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
	// Quick-action labels — populated for the pyeza dashboard block.
	QuickNewItem   string `json:"quickNewItem"`
	QuickViewAll   string `json:"quickViewAll"`
	QuickMovements string `json:"quickMovements"`
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
