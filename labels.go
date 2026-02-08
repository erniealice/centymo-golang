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
	Page    InventoryPageLabels    `json:"page"`
	Buttons InventoryButtonLabels  `json:"buttons"`
	Columns InventoryColumnLabels  `json:"columns"`
	Empty   InventoryEmptyLabels   `json:"empty"`
	Form    InventoryFormLabels    `json:"form"`
	Actions InventoryActionLabels  `json:"actions"`
	Bulk    InventoryBulkLabels    `json:"bulkActions"`
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
	Total       string `json:"total"`
	Discount    string `json:"discount"`
	SubTotal    string `json:"subTotal"`
	GrandTotal  string `json:"grandTotal"`
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
