package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	DB           centymo.DataSource
	Labels       centymo.InventoryLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// AttributeEntry holds a name-value pair for display.
type AttributeEntry struct {
	Name  string
	Value string
}

// SerialSummary holds serial count totals.
type SerialSummary struct {
	Total     int
	Available int
	Sold      int
	Reserved  int
}

// DepreciationInfo holds depreciation policy data for display.
type DepreciationInfo struct {
	ID          string
	Method      string
	CostBasis   string
	SalvageVal  string
	UsefulLife  string
	StartDate   string
	Accumulated string
	BookValue   string
}

// PageData holds the data for the inventory detail page.
type PageData struct {
	types.PageData
	ContentTemplate  string
	Item             map[string]any
	Labels           centymo.InventoryLabels
	ActiveTab        string
	TabItems         []pyeza.TabItem
	IsSerialized     bool
	ItemType         string
	ItemTypeLabel    string
	ItemTypeVariant  string
	LocationName     string
	AvailableQty     string
	Attributes       []AttributeEntry
	SerialTable      *types.TableConfig
	SerialSummary    *SerialSummary
	TransactionTable *types.TableConfig
	Depreciation     *DepreciationInfo
	AuditTable       *types.TableConfig
}

// NewView creates the inventory detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		item, err := deps.DB.Read(ctx, "inventory_item", id)
		if err != nil {
			log.Printf("Failed to read inventory_item %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load inventory item: %w", err))
		}

		name, _ := item["name"].(string)
		locationID, _ := item["location_id"].(string)
		locationName := centymo.LocationDisplayName(locationID)
		headerTitle := name + " \u2014 " + locationName

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}

		l := deps.Labels
		itemType, _ := item["item_type"].(string)
		if itemType == "" {
			itemType = "non_serialized"
		}
		isSerialized := itemType == "serialized"
		tabItems := buildTabItems(l, id, isSerialized)

		available := computeAvailable(item["quantity_on_hand"], item["quantity_reserved"])

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "inventory",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.ItemInfo,
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "inventory-detail-content",
			Item:            item,
			Labels:          l,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
			IsSerialized:    isSerialized,
			ItemType:        itemType,
			ItemTypeLabel:   itemTypeDisplayLabel(itemType, l),
			ItemTypeVariant: itemTypeDisplayVariant(itemType),
			LocationName:    locationName,
			AvailableQty:    available,
		}

		// Load tab-specific data
		switch activeTab {
		case "info":
			// No extra data needed — item map has everything

		case "attributes":
			pageData.Attributes = loadAttributes(ctx, deps.DB, item)

		case "serials":
			serials := loadSerials(ctx, deps.DB, id)
			pageData.SerialTable = buildSerialTable(serials, l, deps.TableLabels, id)
			pageData.SerialSummary = computeSerialSummary(serials)

		case "transactions":
			pageData.TransactionTable = buildTransactionTable(ctx, deps.DB, id, l, deps.TableLabels)

		case "depreciation":
			pageData.Depreciation = loadDepreciation(ctx, deps.DB, id, l)

		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		}

		return view.OK("inventory-detail", pageData)
	})
}

// NewTabAction creates an HTMX tab action view that returns only the tab content partial.
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")

		item, err := deps.DB.Read(ctx, "inventory_item", id)
		if err != nil {
			log.Printf("Failed to read inventory_item %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load inventory item: %w", err))
		}

		l := deps.Labels
		itemType, _ := item["item_type"].(string)
		if itemType == "" {
			itemType = "non_serialized"
		}

		available := computeAvailable(item["quantity_on_hand"], item["quantity_reserved"])

		pageData := &PageData{
			Item:            item,
			Labels:          l,
			ActiveTab:       tab,
			IsSerialized:    itemType == "serialized",
			ItemType:        itemType,
			ItemTypeLabel:   itemTypeDisplayLabel(itemType, l),
			ItemTypeVariant: itemTypeDisplayVariant(itemType),
			LocationName:    centymo.LocationDisplayName(mustString(item["location_id"])),
			AvailableQty:    available,
		}

		switch tab {
		case "info":
			// item map has everything
		case "attributes":
			pageData.Attributes = loadAttributes(ctx, deps.DB, item)
		case "serials":
			serials := loadSerials(ctx, deps.DB, id)
			pageData.SerialTable = buildSerialTable(serials, l, deps.TableLabels, id)
			pageData.SerialSummary = computeSerialSummary(serials)
		case "transactions":
			pageData.TransactionTable = buildTransactionTable(ctx, deps.DB, id, l, deps.TableLabels)
		case "depreciation":
			pageData.Depreciation = loadDepreciation(ctx, deps.DB, id, l)
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		}

		templateName := "inventory-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

func buildTabItems(l centymo.InventoryLabels, id string, isSerialized bool) []pyeza.TabItem {
	base := "/app/inventory/detail/" + id
	tabs := []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", Icon: "icon-info"},
		{Key: "attributes", Label: l.Tabs.Attributes, Href: base + "?tab=attributes", Icon: "icon-layers"},
	}
	if isSerialized {
		tabs = append(tabs, pyeza.TabItem{Key: "serials", Label: l.Tabs.Serials, Href: base + "?tab=serials", Icon: "icon-hash"})
	}
	tabs = append(tabs,
		pyeza.TabItem{Key: "transactions", Label: l.Tabs.Transactions, Href: base + "?tab=transactions", Icon: "icon-repeat"},
		pyeza.TabItem{Key: "depreciation", Label: l.Tabs.Depreciation, Href: base + "?tab=depreciation", Icon: "icon-trending-down"},
		pyeza.TabItem{Key: "audit", Label: l.Tabs.Audit, Href: base + "?tab=audit", Icon: "icon-clock"},
	)
	return tabs
}

func itemTypeDisplayLabel(itemType string, l centymo.InventoryLabels) string {
	switch itemType {
	case "serialized":
		return l.ItemType.Serialized
	case "non_serialized":
		return l.ItemType.NonSerialized
	case "consumable":
		return l.ItemType.Consumable
	default:
		return itemType
	}
}

func itemTypeDisplayVariant(itemType string) string {
	switch itemType {
	case "serialized":
		return "info"
	case "non_serialized":
		return "default"
	case "consumable":
		return "success"
	default:
		return "default"
	}
}

// ---------------------------------------------------------------------------
// Attributes tab
// ---------------------------------------------------------------------------

func loadAttributes(ctx context.Context, db centymo.DataSource, item map[string]any) []AttributeEntry {
	productID, _ := item["product_id"].(string)
	itemID, _ := item["id"].(string)

	if productID == "" {
		return nil
	}

	// 1. Get product_attribute records for this product
	allPA, err := db.ListSimple(ctx, "product_attribute")
	if err != nil {
		log.Printf("Failed to list product_attribute: %v", err)
		return nil
	}
	productAttrs := filterByField(allPA, "product_id", productID)

	// 2. Collect attribute IDs and build lookup
	attrIDs := map[string]bool{}
	for _, pa := range productAttrs {
		if aid, ok := pa["attribute_id"].(string); ok {
			attrIDs[aid] = true
		}
	}

	// 3. Load all attributes and build name map
	allAttrs, err := db.ListSimple(ctx, "attribute")
	if err != nil {
		log.Printf("Failed to list attribute: %v", err)
		return nil
	}
	attrNameMap := map[string]string{}
	for _, a := range allAttrs {
		id, _ := a["id"].(string)
		if attrIDs[id] {
			name, _ := a["name"].(string)
			attrNameMap[id] = name
		}
	}

	// 4. Load inventory_attribute records for this inventory item
	allIA, err := db.ListSimple(ctx, "inventory_attribute")
	if err != nil {
		log.Printf("Failed to list inventory_attribute: %v", err)
		return nil
	}
	itemAttrs := filterByField(allIA, "inventory_item_id", itemID)

	// 5. Build result: join attribute name + inventory value
	iaMap := map[string]string{} // attribute_id -> value
	for _, ia := range itemAttrs {
		aid, _ := ia["attribute_id"].(string)
		val, _ := ia["value"].(string)
		iaMap[aid] = val
	}

	entries := []AttributeEntry{}
	for _, pa := range productAttrs {
		aid, _ := pa["attribute_id"].(string)
		name := attrNameMap[aid]
		value := iaMap[aid]
		if value == "" {
			// Fallback to product default
			value, _ = pa["default_value"].(string)
		}
		if name != "" {
			entries = append(entries, AttributeEntry{Name: name, Value: value})
		}
	}

	return entries
}

// ---------------------------------------------------------------------------
// Serials tab
// ---------------------------------------------------------------------------

func loadSerials(ctx context.Context, db centymo.DataSource, inventoryItemID string) []map[string]any {
	all, err := db.ListSimple(ctx, "inventory_serial")
	if err != nil {
		log.Printf("Failed to list inventory_serial: %v", err)
		return nil
	}
	return filterByField(all, "inventory_item_id", inventoryItemID)
}

func buildSerialTable(serials []map[string]any, l centymo.InventoryLabels, tableLabels types.TableLabels, inventoryItemID string) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "serial_number", Label: l.Detail.SerialNumber, Sortable: true},
		{Key: "imei", Label: l.Detail.IMEI, Sortable: false, Width: "180px"},
		{Key: "status", Label: l.Detail.SerialStatus, Sortable: true, Width: "120px"},
		{Key: "warranty_end", Label: l.Detail.WarrantyEnd, Sortable: true, Width: "140px"},
		{Key: "purchase_order", Label: l.Detail.PurchaseOrder, Sortable: false, Width: "140px"},
		{Key: "sold_reference", Label: l.Detail.SaleReference, Sortable: false, Width: "140px"},
	}

	rows := []types.TableRow{}
	for _, s := range serials {
		id, _ := s["id"].(string)
		serial, _ := s["serial_number"].(string)
		imei, _ := s["imei"].(string)
		status, _ := s["status"].(string)
		warrantyEnd, _ := s["warranty_end"].(string)
		po, _ := s["purchase_order"].(string)
		soldRef, _ := s["sold_reference"].(string)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: serial},
				{Type: "text", Value: imei},
				{Type: "badge", Value: status, Variant: serialStatusVariant(status)},
				{Type: "text", Value: warrantyEnd},
				{Type: "text", Value: po},
				{Type: "text", Value: soldRef},
			},
			Actions: []types.TableAction{
				{Type: "edit", Label: l.Serial.Edit, Action: "edit", URL: "/action/inventory/detail/" + inventoryItemID + "/serials/edit/" + id, DrawerTitle: l.Serial.Edit},
				{Type: "delete", Label: l.Serial.Remove, Action: "delete", URL: "/action/inventory/detail/" + inventoryItemID + "/serials/remove", ItemName: serial},
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   "serial-table",
		RefreshURL:           "/action/inventory/detail/" + inventoryItemID + "/serials/table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowEntries:          true,
		DefaultSortColumn:    "serial_number",
		DefaultSortDirection: "asc",
		Labels:               tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.SerialEmptyTitle,
			Message: l.Detail.SerialEmptyMessage,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     l.Serial.Assign,
			ActionURL: "/action/inventory/detail/" + inventoryItemID + "/serials/assign",
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(cfg)

	return cfg
}

func serialStatusVariant(status string) string {
	switch status {
	case "available":
		return "success"
	case "sold":
		return "default"
	case "reserved":
		return "warning"
	case "defective":
		return "danger"
	case "returned":
		return "warning"
	default:
		return "default"
	}
}

func computeSerialSummary(serials []map[string]any) *SerialSummary {
	summary := &SerialSummary{Total: len(serials)}
	for _, s := range serials {
		status, _ := s["status"].(string)
		switch status {
		case "available":
			summary.Available++
		case "sold":
			summary.Sold++
		case "reserved":
			summary.Reserved++
		}
	}
	return summary
}

// ---------------------------------------------------------------------------
// Transactions tab
// ---------------------------------------------------------------------------

func buildTransactionTable(ctx context.Context, db centymo.DataSource, inventoryItemID string, l centymo.InventoryLabels, tableLabels types.TableLabels) *types.TableConfig {
	all, err := db.ListSimple(ctx, "inventory_transaction")
	if err != nil {
		log.Printf("Failed to list inventory_transaction: %v", err)
		all = []map[string]any{}
	}
	txns := filterByField(all, "inventory_item_id", inventoryItemID)

	columns := []types.TableColumn{
		{Key: "transaction_date", Label: l.Detail.Date, Sortable: true, Width: "130px"},
		{Key: "transaction_type", Label: l.Detail.Type, Sortable: true, Width: "120px"},
		{Key: "quantity", Label: l.Detail.Quantity, Sortable: true, Width: "100px"},
		{Key: "reference", Label: l.Detail.Reference, Sortable: false},
		{Key: "serial_number", Label: l.Detail.Serial, Sortable: false, Width: "150px"},
		{Key: "performed_by", Label: l.Detail.PerformedBy, Sortable: false, Width: "150px"},
	}

	rows := []types.TableRow{}
	for _, t := range txns {
		id, _ := t["id"].(string)
		txDate, _ := t["transaction_date"].(string)
		txType, _ := t["transaction_type"].(string)
		qty := formatQuantity(t["quantity"], txType)
		ref, _ := t["reference"].(string)
		serial, _ := t["serial_number"].(string)
		performer, _ := t["performed_by"].(string)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: txDate},
				{Type: "badge", Value: txType, Variant: txTypeVariant(txType)},
				{Type: "text", Value: qty},
				{Type: "text", Value: ref},
				{Type: "text", Value: serial},
				{Type: "text", Value: performer},
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   "transaction-table",
		RefreshURL:           "/action/inventory/detail/" + inventoryItemID + "/transactions/table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowEntries:          true,
		DefaultSortColumn:    "transaction_date",
		DefaultSortDirection: "desc",
		Labels:               tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.TransactionEmptyTitle,
			Message: l.Detail.TransactionEmptyMessage,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     l.Transaction.Record,
			ActionURL: "/action/inventory/detail/" + inventoryItemID + "/transactions/assign",
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(cfg)

	return cfg
}

func txTypeVariant(txType string) string {
	switch txType {
	case "received":
		return "success"
	case "sold":
		return "default"
	case "adjusted":
		return "info"
	case "transferred":
		return "warning"
	case "returned":
		return "danger"
	case "write_off":
		return "danger"
	default:
		return "default"
	}
}

// formatQuantity formats the quantity with +/- prefix based on transaction type.
func formatQuantity(qty any, txType string) string {
	val := fmt.Sprintf("%v", qty)
	switch txType {
	case "received", "returned":
		return "+" + val
	case "sold", "transferred", "write_off":
		return "-" + val
	default:
		return val
	}
}

// ---------------------------------------------------------------------------
// Depreciation tab
// ---------------------------------------------------------------------------

func loadDepreciation(ctx context.Context, db centymo.DataSource, inventoryItemID string, l centymo.InventoryLabels) *DepreciationInfo {
	all, err := db.ListSimple(ctx, "inventory_depreciation")
	if err != nil {
		log.Printf("Failed to list inventory_depreciation: %v", err)
		return nil
	}

	records := filterByField(all, "inventory_item_id", inventoryItemID)
	if len(records) == 0 {
		return nil
	}

	r := records[0]
	id, _ := r["id"].(string)
	method, _ := r["method"].(string)
	costBasis := fmt.Sprintf("%v", r["cost_basis"])
	salvageVal := fmt.Sprintf("%v", r["salvage_value"])
	usefulLife := fmt.Sprintf("%v", r["useful_life_months"])
	startDate, _ := r["start_date"].(string)
	accumulated := fmt.Sprintf("%v", r["accumulated_depreciation"])
	bookValue := fmt.Sprintf("%v", r["book_value"])

	return &DepreciationInfo{
		ID:          id,
		Method:      depreciationMethodLabel(method, l),
		CostBasis:   costBasis,
		SalvageVal:  salvageVal,
		UsefulLife:  usefulLife + " months",
		StartDate:   startDate,
		Accumulated: accumulated,
		BookValue:   bookValue,
	}
}

func depreciationMethodLabel(method string, l centymo.InventoryLabels) string {
	switch method {
	case "straight_line":
		return l.Depreciation.MethodStraightLine
	case "declining_balance":
		return l.Depreciation.MethodDecliningBalance
	case "sum_of_years":
		return l.Depreciation.MethodSumOfYears
	default:
		return method
	}
}

// ---------------------------------------------------------------------------
// Audit tab
// ---------------------------------------------------------------------------

func buildAuditTable(l centymo.InventoryLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: l.Detail.Date, Sortable: true, Width: "160px"},
		{Key: "action", Label: l.Detail.AuditAction, Sortable: true},
		{Key: "user", Label: l.Detail.AuditUser, Sortable: true, Width: "180px"},
		{Key: "description", Label: l.Detail.Description, Sortable: false},
	}

	rows := []types.TableRow{}
	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   "audit-trail-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowEntries:          true,
		DefaultSortColumn:    "date",
		DefaultSortDirection: "desc",
		Labels:               tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.AuditEmptyTitle,
			Message: l.Detail.AuditEmptyMessage,
		},
	}
	types.ApplyTableSettings(cfg)

	return cfg
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func filterByField(records []map[string]any, field, value string) []map[string]any {
	result := []map[string]any{}
	for _, r := range records {
		v, _ := r[field].(string)
		if v == value {
			result = append(result, r)
		}
	}
	return result
}

func mustString(v any) string {
	s, _ := v.(string)
	return s
}

func computeAvailable(onHandVal, reservedVal any) string {
	onHand := toFloat64(onHandVal)
	reserved := toFloat64(reservedVal)
	avail := onHand - reserved
	if avail < 0 {
		avail = 0
	}
	if avail == float64(int64(avail)) {
		return fmt.Sprintf("%d", int64(avail))
	}
	return fmt.Sprintf("%.2f", avail)
}

func toFloat64(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int64:
		return float64(n)
	default:
		return 0
	}
}
