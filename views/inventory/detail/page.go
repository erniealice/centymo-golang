package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	inventorydepreciationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_depreciation"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	inventorytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_transaction"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes                     centymo.InventoryRoutes
	ReadInventoryItem          func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	ListInventorySerials       func(ctx context.Context, req *inventoryserialpb.ListInventorySerialsRequest) (*inventoryserialpb.ListInventorySerialsResponse, error)
	ListInventoryTransactions  func(ctx context.Context, req *inventorytransactionpb.ListInventoryTransactionsRequest) (*inventorytransactionpb.ListInventoryTransactionsResponse, error)
	ListInventoryDepreciations func(ctx context.Context, req *inventorydepreciationpb.ListInventoryDepreciationsRequest) (*inventorydepreciationpb.ListInventoryDepreciationsResponse, error)
	ListProductVariantOptions  func(ctx context.Context, req *productvariantoptionpb.ListProductVariantOptionsRequest) (*productvariantoptionpb.ListProductVariantOptionsResponse, error)
	ListProductOptionValues    func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	ListProductOptions         func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	Labels                     centymo.InventoryLabels
	CommonLabels               pyeza.CommonLabels
	TableLabels                types.TableLabels

	attachment.AttachmentOps
	auditlog.AuditOps
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
	ContentTemplate     string
	Item                map[string]any
	Labels              centymo.InventoryLabels
	ActiveTab           string
	TabItems            []pyeza.TabItem
	IsSerialized        bool
	TrackingMode        string
	TrackingModeLabel   string
	TrackingModeVariant string
	LocationName        string
	AvailableQty        string
	Attributes          []AttributeEntry
	SerialTable         *types.TableConfig
	SerialSummary       *SerialSummary
	TransactionTable    *types.TableConfig
	Depreciation        *DepreciationInfo
	AuditTable          *types.TableConfig
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// NewView creates the inventory detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		resp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
			Data: &inventoryitempb.InventoryItem{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read inventory_item %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load inventory item: %w", err))
		}
		items := resp.GetData()
		if len(items) == 0 {
			return view.Error(fmt.Errorf("inventory item not found"))
		}
		item := items[0]

		name := item.GetName()
		locationID := item.GetLocationId()
		locationName := centymo.LocationDisplayName(locationID)
		headerTitle := name + " \u2014 " + locationName

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}

		l := deps.Labels
		trackingMode := item.GetProduct().GetTrackingMode()
		if trackingMode == "" {
			trackingMode = "bulk"
		}
		isSerialized := trackingMode == "serialized"
		tabItems := buildTabItems(l, id, isSerialized, deps.Routes)

		available := computeAvailable(item.GetQuantityOnHand(), item.GetQuantityReserved())

		// Build a map[string]any for backward compatibility with templates
		itemMap := inventoryItemToMap(item)

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
			ContentTemplate:   "inventory-detail-content",
			Item:              itemMap,
			Labels:            l,
			ActiveTab:         activeTab,
			TabItems:          tabItems,
			IsSerialized:      isSerialized,
			TrackingMode:      trackingMode,
			TrackingModeLabel: trackingModeDisplayLabel(trackingMode, l),
			TrackingModeVariant: trackingModeDisplayVariant(trackingMode),
			LocationName:      locationName,
			AvailableQty:      available,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "inventory-detail"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		// Load tab-specific data
		switch activeTab {
		case "info":
			// No extra data needed — item map has everything

		case "attributes":
			pageData.Attributes = loadAttributes(ctx, deps, item)

		case "serials":
			perms := view.GetUserPermissions(ctx)
			serials := loadSerials(ctx, deps, id)
			pageData.SerialTable = buildSerialTable(serials, l, deps.TableLabels, id, deps.Routes, perms)
			pageData.SerialSummary = computeSerialSummary(serials)

		case "transactions":
			perms := view.GetUserPermissions(ctx)
			pageData.TransactionTable = buildTransactionTable(ctx, deps, id, l, deps.TableLabels, deps.Routes, perms)

		case "depreciation":
			pageData.Depreciation = loadDepreciation(ctx, deps, id, l)

		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)

		case "attachments":
			if deps.ListAttachments != nil {
				cfg := attachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
				if err != nil {
					log.Printf("Failed to list attachments for %s %s: %v", cfg.EntityType, id, err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)

		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.QueryParams["cursor"]
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "inventory_item",
					EntityID:    id,
					Limit:       20,
					CursorToken: cursor,
				})
				if err != nil {
					log.Printf("Failed to load audit history: %v", err)
				}
				if auditResp != nil {
					pageData.AuditEntries = auditResp.Entries
					pageData.AuditHasNext = auditResp.HasNext
					pageData.AuditNextCursor = auditResp.NextCursor
				}
			}
			pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
		}

		return view.OK("inventory-detail", pageData)
	})
}

// NewTabAction creates an HTMX tab action view that returns only the tab content partial.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")

		resp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
			Data: &inventoryitempb.InventoryItem{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read inventory_item %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load inventory item: %w", err))
		}
		items := resp.GetData()
		if len(items) == 0 {
			return view.Error(fmt.Errorf("inventory item not found"))
		}
		item := items[0]

		l := deps.Labels
		trackingMode := item.GetProduct().GetTrackingMode()
		if trackingMode == "" {
			trackingMode = "bulk"
		}

		available := computeAvailable(item.GetQuantityOnHand(), item.GetQuantityReserved())
		itemMap := inventoryItemToMap(item)

		pageData := &PageData{
			Item:                itemMap,
			Labels:              l,
			ActiveTab:           tab,
			IsSerialized:        trackingMode == "serialized",
			TrackingMode:        trackingMode,
			TrackingModeLabel:   trackingModeDisplayLabel(trackingMode, l),
			TrackingModeVariant: trackingModeDisplayVariant(trackingMode),
			LocationName:        centymo.LocationDisplayName(item.GetLocationId()),
			AvailableQty:        available,
		}

		switch tab {
		case "info":
			// item map has everything
		case "attributes":
			pageData.Attributes = loadAttributes(ctx, deps, item)
		case "serials":
			perms := view.GetUserPermissions(ctx)
			serials := loadSerials(ctx, deps, id)
			pageData.SerialTable = buildSerialTable(serials, l, deps.TableLabels, id, deps.Routes, perms)
			pageData.SerialSummary = computeSerialSummary(serials)
		case "transactions":
			perms := view.GetUserPermissions(ctx)
			pageData.TransactionTable = buildTransactionTable(ctx, deps, id, l, deps.TableLabels, deps.Routes, perms)
		case "depreciation":
			pageData.Depreciation = loadDepreciation(ctx, deps, id, l)
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)

		case "attachments":
			if deps.ListAttachments != nil {
				cfg := attachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
				if err != nil {
					log.Printf("Failed to list attachments for %s %s: %v", cfg.EntityType, id, err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)

		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.QueryParams["cursor"]
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "inventory_item",
					EntityID:    id,
					Limit:       20,
					CursorToken: cursor,
				})
				if err != nil {
					log.Printf("Failed to load audit history: %v", err)
				}
				if auditResp != nil {
					pageData.AuditEntries = auditResp.Entries
					pageData.AuditHasNext = auditResp.HasNext
					pageData.AuditNextCursor = auditResp.NextCursor
				}
			}
			pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
		}

		templateName := "inventory-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// inventoryItemToMap converts a proto InventoryItem to a map[string]any for template backward compat.
func inventoryItemToMap(item *inventoryitempb.InventoryItem) map[string]any {
	m := map[string]any{
		"id":                 item.GetId(),
		"name":               item.GetName(),
		"active":             item.GetActive(),
		"sku":                item.GetSku(),
		"quantity_on_hand":   item.GetQuantityOnHand(),
		"quantity_reserved":  item.GetQuantityReserved(),
		"quantity_available": item.GetQuantityAvailable(),
		"reorder_level":      item.GetReorderLevel(),
		"unit_of_measure":    item.GetUnitOfMeasure(),
		"notes":              item.GetNotes(),
		"tracking_mode":      item.GetProduct().GetTrackingMode(),
		"location_id":        item.GetLocationId(),
		"product_id":         item.GetProductId(),
		"product_variant_id": item.GetProductVariantId(),
	}
	return m
}

func buildTabItems(l centymo.InventoryLabels, id string, isSerialized bool, routes centymo.InventoryRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	tabs := []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "attributes", Label: l.Tabs.Attributes, Href: base + "?tab=attributes", HxGet: action + "attributes", Icon: "icon-layers"},
	}
	if isSerialized {
		tabs = append(tabs, pyeza.TabItem{Key: "serials", Label: l.Tabs.Serials, Href: base + "?tab=serials", HxGet: action + "serials", Icon: "icon-hash"})
	}
	tabs = append(tabs,
		pyeza.TabItem{Key: "transactions", Label: l.Tabs.Transactions, Href: base + "?tab=transactions", HxGet: action + "transactions", Icon: "icon-repeat"},
		pyeza.TabItem{Key: "depreciation", Label: l.Tabs.Depreciation, Href: base + "?tab=depreciation", HxGet: action + "depreciation", Icon: "icon-trending-down"},
		pyeza.TabItem{Key: "audit", Label: l.Tabs.Audit, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
		pyeza.TabItem{Key: "attachments", Label: l.Tabs.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		pyeza.TabItem{Key: "audit-history", Label: "History", Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	)
	return tabs
}

func trackingModeDisplayLabel(trackingMode string, l centymo.InventoryLabels) string {
	switch trackingMode {
	case "none":
		return l.TrackingMode.None
	case "bulk":
		return l.TrackingMode.Bulk
	case "serialized":
		return l.TrackingMode.Serialized
	default:
		return trackingMode
	}
}

func trackingModeDisplayVariant(trackingMode string) string {
	switch trackingMode {
	case "none":
		return "neutral"
	case "bulk":
		return "info"
	case "serialized":
		return "success"
	default:
		return "default"
	}
}

// ---------------------------------------------------------------------------
// Attributes tab
// ---------------------------------------------------------------------------

func loadAttributes(ctx context.Context, deps *DetailViewDeps, item *inventoryitempb.InventoryItem) []AttributeEntry {
	variantID := item.GetProductVariantId()
	if variantID == "" {
		return nil
	}

	// 1. Get product_variant_option records for this variant
	pvoResp, err := deps.ListProductVariantOptions(ctx, &productvariantoptionpb.ListProductVariantOptionsRequest{})
	if err != nil {
		log.Printf("Failed to list product_variant_option: %v", err)
		return nil
	}
	allPVO := pvoResp.GetData()
	var variantOptions []*productvariantoptionpb.ProductVariantOption
	for _, pvo := range allPVO {
		if pvo.GetProductVariantId() == variantID {
			variantOptions = append(variantOptions, pvo)
		}
	}
	if len(variantOptions) == 0 {
		return nil
	}

	// 2. Collect option_value IDs
	valueIDs := map[string]bool{}
	for _, pvo := range variantOptions {
		vid := pvo.GetProductOptionValueId()
		if vid != "" {
			valueIDs[vid] = true
		}
	}

	// 3. Load product_option_value records and build lookup
	povResp, err := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
	if err != nil {
		log.Printf("Failed to list product_option_value: %v", err)
		return nil
	}
	allPOV := povResp.GetData()
	valueMap := map[string]*productoptionvaluepb.ProductOptionValue{}
	optionIDs := map[string]bool{}
	for _, pov := range allPOV {
		if valueIDs[pov.GetId()] {
			valueMap[pov.GetId()] = pov
			oid := pov.GetProductOptionId()
			if oid != "" {
				optionIDs[oid] = true
			}
		}
	}

	// 4. Load product_option records for names and sort order
	poResp, err := deps.ListProductOptions(ctx, &productoptionpb.ListProductOptionsRequest{})
	if err != nil {
		log.Printf("Failed to list product_option: %v", err)
		return nil
	}
	allPO := poResp.GetData()
	optionMap := map[string]*productoptionpb.ProductOption{}
	for _, po := range allPO {
		if optionIDs[po.GetId()] {
			optionMap[po.GetId()] = po
		}
	}

	// 5. Build sorted entries: option name + value label
	type sortedEntry struct {
		sortOrder int32
		entry     AttributeEntry
	}
	var sorted []sortedEntry
	for _, pvo := range variantOptions {
		vid := pvo.GetProductOptionValueId()
		pov := valueMap[vid]
		if pov == nil {
			continue
		}
		oid := pov.GetProductOptionId()
		po := optionMap[oid]
		if po == nil {
			continue
		}
		name := po.GetName()
		label := pov.GetLabel()
		order := po.GetSortOrder()
		if name != "" && label != "" {
			sorted = append(sorted, sortedEntry{
				sortOrder: order,
				entry:     AttributeEntry{Name: name, Value: label},
			})
		}
	}

	// Sort by option sort_order
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].sortOrder < sorted[i].sortOrder {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	entries := make([]AttributeEntry, len(sorted))
	for i, s := range sorted {
		entries[i] = s.entry
	}
	return entries
}

// ---------------------------------------------------------------------------
// Serials tab
// ---------------------------------------------------------------------------

func loadSerials(ctx context.Context, deps *DetailViewDeps, inventoryItemID string) []*inventoryserialpb.InventorySerial {
	resp, err := deps.ListInventorySerials(ctx, &inventoryserialpb.ListInventorySerialsRequest{
		InventoryItemId: &inventoryItemID,
	})
	if err != nil {
		log.Printf("Failed to list inventory_serial: %v", err)
		return nil
	}
	return resp.GetData()
}

func buildSerialTable(serials []*inventoryserialpb.InventorySerial, l centymo.InventoryLabels, tableLabels types.TableLabels, inventoryItemID string, routes centymo.InventoryRoutes, perms *types.UserPermissions) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "serial_number", Label: l.Detail.SerialNumber},
		{Key: "imei", Label: l.Detail.IMEI, NoSort: true, WidthClass: "col-6xl"},
		{Key: "status", Label: l.Detail.SerialStatus, WidthClass: "col-2xl"},
		{Key: "warranty_end", Label: l.Detail.WarrantyEnd, WidthClass: "col-3xl"},
		{Key: "purchase_order", Label: l.Detail.PurchaseOrder, NoSort: true, WidthClass: "col-3xl"},
	}

	rows := []types.TableRow{}
	for _, s := range serials {
		id := s.GetId()
		serial := s.GetSerialNumber()
		imei := s.GetImei()
		status := s.GetStatus()
		warrantyEnd := s.GetWarrantyEnd()
		po := s.GetPurchaseOrder()

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: serial},
				{Type: "text", Value: imei},
				{Type: "badge", Value: status, Variant: serialStatusVariant(status)},
				{Type: "text", Value: warrantyEnd},
				{Type: "text", Value: po},
			},
			Actions: []types.TableAction{
				{Type: "edit", Label: l.Serial.Edit, Action: "edit", URL: route.ResolveURL(routes.SerialEditURL, "id", inventoryItemID, "sid", id), DrawerTitle: l.Serial.Edit, Disabled: !perms.Can("inventory_item", "update"), DisabledTooltip: l.Errors.PermissionDenied},
				{Type: "delete", Label: l.Serial.Remove, Action: "delete", URL: route.ResolveURL(routes.SerialRemoveURL, "id", inventoryItemID), ItemName: serial, Disabled: !perms.Can("inventory_item", "delete"), DisabledTooltip: l.Errors.PermissionDenied},
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   "serial-table",
		RefreshURL:           route.ResolveURL(routes.SerialTableURL, "id", inventoryItemID),
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
			Label:           l.Serial.Assign,
			ActionURL:       route.ResolveURL(routes.SerialAssignURL, "id", inventoryItemID),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("inventory_item", "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
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

func computeSerialSummary(serials []*inventoryserialpb.InventorySerial) *SerialSummary {
	summary := &SerialSummary{Total: len(serials)}
	for _, s := range serials {
		switch s.GetStatus() {
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

func buildTransactionTable(ctx context.Context, deps *DetailViewDeps, inventoryItemID string, l centymo.InventoryLabels, tableLabels types.TableLabels, routes centymo.InventoryRoutes, perms *types.UserPermissions) *types.TableConfig {
	resp, err := deps.ListInventoryTransactions(ctx, &inventorytransactionpb.ListInventoryTransactionsRequest{
		InventoryItemId: &inventoryItemID,
	})
	if err != nil {
		log.Printf("Failed to list inventory_transaction: %v", err)
	}
	var txns []*inventorytransactionpb.InventoryTransaction
	if resp != nil {
		txns = resp.GetData()
	}

	columns := []types.TableColumn{
		{Key: "transaction_date", Label: l.Detail.Date, WidthClass: "col-3xl"},
		{Key: "transaction_type", Label: l.Detail.Type, WidthClass: "col-2xl"},
		{Key: "quantity", Label: l.Detail.Quantity, WidthClass: "col-lg"},
		{Key: "reference", Label: l.Detail.Reference, NoSort: true},
		{Key: "serial_number", Label: l.Detail.Serial, NoSort: true, WidthClass: "col-4xl"},
		{Key: "performed_by", Label: l.Detail.PerformedBy, NoSort: true, WidthClass: "col-4xl"},
	}

	rows := []types.TableRow{}
	for _, t := range txns {
		id := t.GetId()
		txDate := t.GetTransactionDateString()
		txType := t.GetTransactionType()
		qty := formatQuantity(t.GetQuantity(), txType)
		ref := t.GetReferenceType()
		serial := t.GetSerialNumber()
		performer := t.GetPerformedBy()

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
		RefreshURL:           route.ResolveURL(routes.TransactionTableURL, "id", inventoryItemID),
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
			Label:           l.Transaction.Record,
			ActionURL:       route.ResolveURL(routes.TransactionAssignURL, "id", inventoryItemID),
			Icon:            "icon-plus",
			Disabled:        !perms.Can("inventory_item", "create"),
			DisabledTooltip: l.Errors.PermissionDenied,
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
func formatQuantity(qty float64, txType string) string {
	val := fmt.Sprintf("%g", qty)
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

func loadDepreciation(ctx context.Context, deps *DetailViewDeps, inventoryItemID string, l centymo.InventoryLabels) *DepreciationInfo {
	resp, err := deps.ListInventoryDepreciations(ctx, &inventorydepreciationpb.ListInventoryDepreciationsRequest{
		InventoryItemId: &inventoryItemID,
	})
	if err != nil {
		log.Printf("Failed to list inventory_depreciation: %v", err)
		return nil
	}

	records := resp.GetData()
	if len(records) == 0 {
		return nil
	}

	r := records[0]
	return &DepreciationInfo{
		ID:          r.GetId(),
		Method:      depreciationMethodLabel(r.GetMethod(), l),
		CostBasis:   fmt.Sprintf("%g", r.GetCostBasis()),
		SalvageVal:  fmt.Sprintf("%g", r.GetSalvageValue()),
		UsefulLife:  fmt.Sprintf("%d %s", r.GetUsefulLifeMonths(), l.Depreciation.MonthsUnit),
		StartDate:   r.GetStartDate(),
		Accumulated: fmt.Sprintf("%g", r.GetAccumulatedDepreciation()),
		BookValue:   fmt.Sprintf("%g", r.GetBookValue()),
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
		{Key: "date", Label: l.Detail.Date, WidthClass: "col-5xl"},
		{Key: "action", Label: l.Detail.AuditAction},
		{Key: "user", Label: l.Detail.AuditUser, WidthClass: "col-6xl"},
		{Key: "description", Label: l.Detail.Description, NoSort: true},
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

func computeAvailable(onHand, reserved float64) string {
	avail := onHand - reserved
	if avail < 0 {
		avail = 0
	}
	if avail == float64(int64(avail)) {
		return fmt.Sprintf("%d", int64(avail))
	}
	return fmt.Sprintf("%.2f", avail)
}
