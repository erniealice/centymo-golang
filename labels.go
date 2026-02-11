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
}

type InventoryMovementsLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

// ---------------------------------------------------------------------------
// Sales labels
// ---------------------------------------------------------------------------

// SalesLabels holds all translatable strings for the sales (revenue) module.
type SalesLabels struct {
	Page    SalesPageLabels    `json:"page"`
	Buttons SalesButtonLabels  `json:"buttons"`
	Columns SalesColumnLabels  `json:"columns"`
	Empty   SalesEmptyLabels   `json:"empty"`
	Form    SalesFormLabels    `json:"form"`
	Actions SalesActionLabels  `json:"actions"`
	Bulk    SalesBulkLabels    `json:"bulkActions"`
	Detail  SalesDetailLabels  `json:"detail"`
}

type SalesPageLabels struct {
	Heading          string `json:"heading"`
	HeadingActive    string `json:"headingActive"`
	HeadingCompleted string `json:"headingCompleted"`
	HeadingCancelled string `json:"headingCancelled"`
	Caption          string `json:"caption"`
	CaptionActive    string `json:"captionActive"`
	CaptionCompleted string `json:"captionCompleted"`
	CaptionCancelled string `json:"captionCancelled"`
}

type SalesButtonLabels struct {
	AddSale string `json:"addSale"`
}

type SalesColumnLabels struct {
	Reference  string `json:"reference"`
	Customer   string `json:"customer"`
	Date       string `json:"date"`
	Amount     string `json:"amount"`
	Status     string `json:"status"`
}

type SalesEmptyLabels struct {
	ActiveTitle      string `json:"activeTitle"`
	ActiveMessage    string `json:"activeMessage"`
	CompletedTitle   string `json:"completedTitle"`
	CompletedMessage string `json:"completedMessage"`
	CancelledTitle   string `json:"cancelledTitle"`
	CancelledMessage string `json:"cancelledMessage"`
}

type SalesFormLabels struct {
	Customer          string `json:"customer"`
	Date              string `json:"date"`
	Amount            string `json:"amount"`
	Currency          string `json:"currency"`
	Reference         string `json:"reference"`
	ReferencePlaceholder string `json:"referencePlaceholder"`
	Status            string `json:"status"`
	Notes             string `json:"notes"`
	NotesPlaceholder  string `json:"notesPlaceholder"`
	Active            string `json:"active"`
}

type SalesActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type SalesBulkLabels struct {
	Delete string `json:"delete"`
}

type SalesDetailLabels struct {
	PageTitle   string `json:"pageTitle"`
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
}

type ProductDetailLabels struct {
	Price        string `json:"price"`
	Currency     string `json:"currency"`
	Collections  string `json:"collections"`
	VariantCount string `json:"variantCount"`
	Status       string `json:"status"`
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
}

type ProductAttributeLabels struct {
	Title        string `json:"title"`
	DefaultValue string `json:"defaultValue"`
	Assign       string `json:"assign"`
	Remove       string `json:"remove"`
	Empty        string `json:"empty"`
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
	BasicInfo   string `json:"basicInfo"`
	Prices      string `json:"prices"`
	ProductName string `json:"productName"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
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
