package inventory

import (
	sibProductProduct "github.com/erniealice/centymo-golang/domain/product/product"
)

// ---------------------------------------------------------------------------
// Inventory labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the inventory module.
type Labels struct {
	Page         PageLabels                           `json:"page"`
	Buttons      ButtonLabels                         `json:"buttons"`
	Columns      ColumnLabels                         `json:"columns"`
	Empty        EmptyLabels                          `json:"empty"`
	Form         FormLabels                           `json:"form"`
	Actions      ActionLabels                         `json:"actions"`
	Bulk         BulkLabels                           `json:"bulk_actions"`
	Detail       DetailLabels                         `json:"detail"`
	Tabs         TabLabels                            `json:"tabs"`
	TrackingMode sibProductProduct.TrackingModeLabels `json:"tracking_mode"`
	Status       StatusLabels                         `json:"status"`
	Serial       SerialLabels                         `json:"serial"`
	Transaction  TransactionLabels                    `json:"transaction"`
	Depreciation DepreciationLabels                   `json:"depreciation"`
	Dashboard    DashboardLabels                      `json:"dashboard"`
	Movements    MovementsLabels                      `json:"movements"`
	Confirm      ConfirmLabels                        `json:"confirm"`
	Errors       ErrorLabels                          `json:"errors"`
	Breadcrumb   BreadcrumbLabels                     `json:"breadcrumb"`
}

type PageLabels struct {
	Heading  string `json:"heading"`
	Caption  string `json:"caption"`
	Location string `json:"location"`
}

type ButtonLabels struct {
	AddItem string `json:"add_item"`
}

type ColumnLabels struct {
	ProductName string `json:"product_name"`
	SKU         string `json:"sku"`
	OnHand      string `json:"on_hand"`
	Available   string `json:"available"`
	ReorderLvl  string `json:"reorder_level"`
	Status      string `json:"status"`
	Type        string `json:"type"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type FormLabels struct {
	Product          string `json:"product"`
	SKU              string `json:"sku"`
	SKUPlaceholder   string `json:"sku_placeholder"`
	OnHand           string `json:"on_hand"`
	Reserved         string `json:"reserved"`
	ReorderLevel     string `json:"reorder_level"`
	UnitOfMeasure    string `json:"unit_of_measure"`
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notes_placeholder"`
	Active           string `json:"active"`

	// Field-level info text surfaced via an info button beside each label.
	ProductInfo       string `json:"product_info"`
	SKUInfo           string `json:"sku_info"`
	OnHandInfo        string `json:"on_hand_info"`
	ReservedInfo      string `json:"reserved_info"`
	ReorderLevelInfo  string `json:"reorder_level_info"`
	UnitOfMeasureInfo string `json:"unit_of_measure_info"`
	NotesInfo         string `json:"notes_info"`
	ActiveInfo        string `json:"active_info"`
}

type ActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

// DetailLabels holds all translatable strings for the inventory detail page.
type DetailLabels struct {
	TitlePrefix string `json:"title_prefix"`
	MonthsUnit  string `json:"months_unit"`
	TypesUnit   string `json:"types_unit"`

	// Tab labels
	TabBasicInfo    string `json:"tab_basic_info"`
	TabAttributes   string `json:"tab_attributes"`
	TabSerials      string `json:"tab_serials"`
	TabTransactions string `json:"tab_transactions"`
	TabAuditTrail   string `json:"tab_audit_trail"`

	// Info fields
	ItemInfo      string `json:"item_info"`
	ProductName   string `json:"product_name"`
	SKU           string `json:"sku"`
	Location      string `json:"location"`
	OnHand        string `json:"on_hand"`
	Reserved      string `json:"reserved"`
	Available     string `json:"available"`
	ReorderLevel  string `json:"reorder_level"`
	UnitOfMeasure string `json:"unit_of_measure"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`

	// Attribute labels
	AttributeName  string `json:"attribute_name"`
	AttributeValue string `json:"attribute_value"`

	// Serial columns
	SerialNumber     string `json:"serial_number"`
	IMEI             string `json:"imei"`
	SerialStatus     string `json:"serial_status"`
	WarrantyEnd      string `json:"warranty_end"`
	PurchaseOrder    string `json:"purchase_order"`
	RevenueReference string `json:"revenue_reference"`

	// Serial summary
	TotalUnits     string `json:"total_units"`
	AvailableUnits string `json:"available_units"`
	SoldUnits      string `json:"sold_units"`
	ReservedUnits  string `json:"reserved_units"`

	// Transaction columns
	Date        string `json:"date"`
	Type        string `json:"type"`
	Quantity    string `json:"quantity"`
	Reference   string `json:"reference"`
	Serial      string `json:"serial"`
	PerformedBy string `json:"performed_by"`

	// Audit columns
	AuditAction string `json:"audit_action"`
	AuditUser   string `json:"audit_user"`
	Description string `json:"description"`

	// Info field labels (not shared with transaction columns)
	Product       string `json:"product"`
	ViewProduct   string `json:"view_product"`
	Active        string `json:"active"`
	Inactive      string `json:"inactive"`
	SerialNumbers string `json:"serial_numbers"`

	// Empty states
	AttributeEmptyTitle     string `json:"attribute_empty_title"`
	AttributeEmptyMessage   string `json:"attribute_empty_message"`
	SerialEmptyTitle        string `json:"serial_empty_title"`
	SerialEmptyMessage      string `json:"serial_empty_message"`
	TransactionEmptyTitle   string `json:"transaction_empty_title"`
	TransactionEmptyMessage string `json:"transaction_empty_message"`
	AuditEmptyTitle         string `json:"audit_empty_title"`
	AuditEmptyMessage       string `json:"audit_empty_message"`
}

type TabLabels struct {
	Info         string `json:"info"`
	Attributes   string `json:"attributes"`
	Serials      string `json:"serials"`
	Transactions string `json:"transactions"`
	Depreciation string `json:"depreciation"`
	Audit        string `json:"audit"`
	Attachments  string `json:"attachments"`
	AuditHistory string `json:"audit_history"`
}

type StatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type SerialLabels struct {
	Title           string `json:"title"`
	SerialNumber    string `json:"serial_number"`
	IMEI            string `json:"imei"`
	Status          string `json:"status"`
	WarrantyStart   string `json:"warranty_start"`
	WarrantyEnd     string `json:"warranty_end"`
	PurchaseOrder   string `json:"purchase_order"`
	SoldReference   string `json:"sold_reference"`
	Assign          string `json:"assign"`
	Edit            string `json:"edit"`
	Remove          string `json:"remove"`
	Empty           string `json:"empty"`
	StatusAvailable string `json:"status_available"`
	StatusSold      string `json:"status_sold"`
	StatusReserved  string `json:"status_reserved"`
	StatusDefective string `json:"status_defective"`
	StatusReturned  string `json:"status_returned"`

	// Field-level info text surfaced via an info button beside each label.
	SerialNumberInfo  string `json:"serial_number_info"`
	IMEIInfo          string `json:"imei_info"`
	StatusInfo        string `json:"status_info"`
	WarrantyStartInfo string `json:"warranty_start_info"`
	WarrantyEndInfo   string `json:"warranty_end_info"`
	PurchaseOrderInfo string `json:"purchase_order_info"`
	SoldReferenceInfo string `json:"sold_reference_info"`
}

type TransactionLabels struct {
	Title           string `json:"title"`
	Type            string `json:"type"`
	Quantity        string `json:"quantity"`
	Date            string `json:"date"`
	Reference       string `json:"reference"`
	PerformedBy     string `json:"performed_by"`
	Record          string `json:"record"`
	Empty           string `json:"empty"`
	TypeReceived    string `json:"type_received"`
	TypeSold        string `json:"type_sold"`
	TypeAdjusted    string `json:"type_adjusted"`
	TypeTransferred string `json:"type_transferred"`
	TypeReturned    string `json:"type_returned"`
	TypeWriteOff    string `json:"type_write_off"`

	// Field-level info text surfaced via an info button beside each label.
	TypeInfo      string `json:"type_info"`
	QuantityInfo  string `json:"quantity_info"`
	DateInfo      string `json:"date_info"`
	ReferenceInfo string `json:"reference_info"`
}

type DepreciationLabels struct {
	Title                  string `json:"title"`
	Method                 string `json:"method"`
	CostBasis              string `json:"cost_basis"`
	SalvageValue           string `json:"salvage_value"`
	UsefulLife             string `json:"useful_life"`
	StartDate              string `json:"start_date"`
	Accumulated            string `json:"accumulated"`
	BookValue              string `json:"book_value"`
	Configure              string `json:"configure"`
	Edit                   string `json:"edit"`
	NotConfigured          string `json:"not_configured"`
	MethodStraightLine     string `json:"method_straight_line"`
	MethodDecliningBalance string `json:"method_declining_balance"`
	MethodSumOfYears       string `json:"method_sum_of_years"`
	MonthsUnit             string `json:"months_unit"`

	// Field-level info text surfaced via an info button beside each label.
	MethodInfo       string `json:"method_info"`
	CostBasisInfo    string `json:"cost_basis_info"`
	SalvageValueInfo string `json:"salvage_value_info"`
	UsefulLifeInfo   string `json:"useful_life_info"`
	StartDateInfo    string `json:"start_date_info"`
}

type DashboardLabels struct {
	Title                string `json:"title"`
	TotalStockValue      string `json:"total_stock_value"`
	LowStockAlerts       string `json:"low_stock_alerts"`
	StockTurnover        string `json:"stock_turnover"`
	ItemsByLocation      string `json:"items_by_location"`
	DepreciationSummary  string `json:"depreciation_summary"`
	SerialUnitStatus     string `json:"serial_unit_status"`
	RecentMovements      string `json:"recent_movements"`
	CategoryDistribution string `json:"category_distribution"`
	TypesUnit            string `json:"types_unit"`
	StockLevels          string `json:"stock_levels"`
	RecentActivity       string `json:"recent_activity"`
	ViewAll              string `json:"view_all"`
	Week                 string `json:"week"`
	Month                string `json:"month"`
	Year                 string `json:"year"`
	// Quick-action labels — populated for the pyeza dashboard block.
	QuickNewItem   string `json:"quick_new_item"`
	QuickViewAll   string `json:"quick_view_all"`
	QuickMovements string `json:"quick_movements"`
}

type MovementsLabels struct {
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle"`
	DateRange      string `json:"date_range"`
	LocationFilter string `json:"location_filter"`
	TypeFilter     string `json:"type_filter"`
	ProductSearch  string `json:"product_search"`
	ClearAll       string `json:"clear_all"`
	ExportCsv      string `json:"export_csv"`
	AllLocations   string `json:"all_locations"`
	AllTypes       string `json:"all_types"`
	ProductColumn  string `json:"product_column"`
	VariantSKU     string `json:"variant_sku"`
}

type ConfirmLabels struct {
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activate_message"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivate_message"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"delete_message"`
	BulkActivate          string `json:"bulk_activate"`
	BulkActivateMessage   string `json:"bulk_activate_message"`
	BulkDeactivate        string `json:"bulk_deactivate"`
	BulkDeactivateMessage string `json:"bulk_deactivate_message"`
	BulkDelete            string `json:"bulk_delete"`
	BulkDeleteMessage     string `json:"bulk_delete_message"`
}

type ErrorLabels struct {
	PermissionDenied          string `json:"permission_denied"`
	InvalidFormData           string `json:"invalid_form_data"`
	NotFound                  string `json:"not_found"`
	IDRequired                string `json:"id_required"`
	NoIDsProvided             string `json:"no_ids_provided"`
	InvalidStatus             string `json:"invalid_status"`
	NoPermission              string `json:"no_permission"`
	SerialNotFound            string `json:"serial_not_found"`
	SerialIDRequired          string `json:"serial_idrequired"`
	InvalidDepreciationMethod string `json:"invalid_depreciation_method"`
}

type BreadcrumbLabels struct {
	Products string `json:"products"`
	Product  string `json:"product"`
}

// DefaultLabels returns the zero-value label set. Every rendered string for
// this entity must come from the lyngua cascade (general -> business-type
// tier); there are no Go-side default strings to fall back on.
func DefaultLabels() Labels { return Labels{} }
