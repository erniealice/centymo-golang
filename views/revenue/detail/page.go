package detail

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes       centymo.RevenueRoutes
	DB           centymo.DataSource // KEEP — used for revenue_payment operations
	Labels       centymo.RevenueLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Typed revenue operations
	ReadRevenue func(ctx context.Context, req *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)

	// Typed line item operations
	ListRevenueLineItems func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error)

	attachment.AttachmentOps
	auditlog.AuditOps
}

// PaymentInfo holds payment details for the payment tab.
type PaymentInfo struct {
	Method       string
	AmountPaid   string
	Currency     string
	CardLast4    string
	PaymentDate  string
	ReceivedBy   string
	ReceivedRole string
}

// PageData holds the data for the sales detail page.
type PageData struct {
	types.PageData
	ContentTemplate     string
	Revenue             map[string]any
	Labels              centymo.RevenueLabels
	ActiveTab           string
	TabItems            []pyeza.TabItem
	LineItemTable       *types.TableConfig
	LineItemAddURL      string
	LineItemDiscountURL string
	TotalAmount         string
	Payment             *PaymentInfo
	PaymentTable        *types.TableConfig
	PaymentAddURL       string
	TotalPaid           string
	RemainingBalance    string
	PaymentStatus       string
	AuditTable          *types.TableConfig
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
	InvoiceDownloadURL  string
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// NewView creates the sales detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		resp, err := deps.ReadRevenue(ctx, &revenuepb.ReadRevenueRequest{
			Data: &revenuepb.Revenue{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read revenue %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load sale: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			log.Printf("Revenue %s not found", id)
			return view.Error(fmt.Errorf("sale not found"))
		}
		revenue := revenueToMap(data[0])

		refNumber, _ := revenue["reference_number"].(string)

		l := deps.Labels
		headerTitle := l.Detail.TitlePrefix + refNumber

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}
		tabItems := buildTabItems(l, id, deps.Routes)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "sale",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.PageTitle,
				HeaderIcon:     "icon-shopping-bag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:    "sales-detail-content",
			Revenue:            revenue,
			Labels:             l,
			ActiveTab:          activeTab,
			TabItems:           tabItems,
			InvoiceDownloadURL: route.ResolveURL(deps.Routes.InvoiceDownloadURL, "id", id),
		}

		// Load tab-specific data
		switch activeTab {
		case "info":
			// No extra data needed — revenue map has everything
		case "items":
			perms := view.GetUserPermissions(ctx)
			lineItems := listLineItemMaps(ctx, deps.ListRevenueLineItems, id)
			currency, _ := revenue["currency"].(string)
			pageData.LineItemTable = buildLineItemTableWithActions(lineItems, l, deps.TableLabels, currency, id, deps.Routes, perms)
			pageData.LineItemAddURL = route.ResolveURL(deps.Routes.LineItemAddURL, "id", id)
			pageData.LineItemDiscountURL = route.ResolveURL(deps.Routes.LineItemDiscountURL, "id", id)
			totalAmount, _ := revenue["total_amount"].(string)
			pageData.TotalAmount = currency + " " + totalAmount

		case "payment":
			allPayments, err := deps.DB.ListSimple(ctx, "revenue_payment")
			if err != nil {
				log.Printf("Failed to list payments for revenue %s: %v", id, err)
				allPayments = []map[string]any{}
			}
			payments := filterPayments(allPayments, id)
			currency, _ := revenue["currency"].(string)
			perms := view.GetUserPermissions(ctx)
			pageData.PaymentTable = buildPaymentTable(payments, l, deps.TableLabels, currency, id, deps.Routes, perms)
			pageData.PaymentAddURL = route.ResolveURL(deps.Routes.PaymentAddURL, "id", id)

			// Calculate totals
			totalAmount, _ := revenue["total_amount"].(string)
			totalPaid := sumPayments(payments)
			totalAmountFloat := parseAmount(totalAmount)
			remaining := totalAmountFloat - totalPaid

			pageData.TotalPaid = currency + " " + formatAmount(totalPaid)
			pageData.RemainingBalance = currency + " " + formatAmount(remaining)
			if remaining <= 0 {
				pageData.PaymentStatus = "paid"
			} else if totalPaid > 0 {
				pageData.PaymentStatus = "partial"
			} else {
				pageData.PaymentStatus = "unpaid"
			}
			// Keep legacy field for backward compat
			pageData.Payment = findPayment(allPayments, id, revenue)

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
					EntityType:  "revenue",
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

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "sale-detail"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("sale-detail", pageData)
	})
}

func buildTabItems(l centymo.RevenueLabels, id string, routes centymo.RevenueRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Detail.TabBasicInfo, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "items", Label: l.Detail.TabLineItems, Href: base + "?tab=items", HxGet: action + "items", Icon: "icon-list"},
		{Key: "payment", Label: l.Detail.TabPayment, Href: base + "?tab=payment", HxGet: action + "payment", Icon: "icon-credit-card"},
		{Key: "audit", Label: l.Detail.TabAuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
		{Key: "attachments", Label: l.Detail.TabAttachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		{Key: "audit-history", Label: "History", Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		resp, err := deps.ReadRevenue(ctx, &revenuepb.ReadRevenueRequest{
			Data: &revenuepb.Revenue{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read revenue %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load sale: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			log.Printf("Revenue %s not found", id)
			return view.Error(fmt.Errorf("sale not found"))
		}
		revenue := revenueToMap(data[0])

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Revenue:            revenue,
			Labels:             l,
			ActiveTab:          tab,
			TabItems:           buildTabItems(l, id, deps.Routes),
			InvoiceDownloadURL: route.ResolveURL(deps.Routes.InvoiceDownloadURL, "id", id),
		}

		switch tab {
		case "info":
			// revenue map has everything
		case "items":
			perms := view.GetUserPermissions(ctx)
			lineItems := listLineItemMaps(ctx, deps.ListRevenueLineItems, id)
			currency, _ := revenue["currency"].(string)
			pageData.LineItemTable = buildLineItemTableWithActions(lineItems, l, deps.TableLabels, currency, id, deps.Routes, perms)
			pageData.LineItemAddURL = route.ResolveURL(deps.Routes.LineItemAddURL, "id", id)
			pageData.LineItemDiscountURL = route.ResolveURL(deps.Routes.LineItemDiscountURL, "id", id)
			totalAmount, _ := revenue["total_amount"].(string)
			pageData.TotalAmount = currency + " " + totalAmount

		case "payment":
			allPayments, err := deps.DB.ListSimple(ctx, "revenue_payment")
			if err != nil {
				log.Printf("Failed to list payments for revenue %s: %v", id, err)
				allPayments = []map[string]any{}
			}
			payments := filterPayments(allPayments, id)
			currency, _ := revenue["currency"].(string)
			perms := view.GetUserPermissions(ctx)
			pageData.PaymentTable = buildPaymentTable(payments, l, deps.TableLabels, currency, id, deps.Routes, perms)
			pageData.PaymentAddURL = route.ResolveURL(deps.Routes.PaymentAddURL, "id", id)

			totalAmount, _ := revenue["total_amount"].(string)
			totalPaid := sumPayments(payments)
			totalAmountFloat := parseAmount(totalAmount)
			remaining := totalAmountFloat - totalPaid

			pageData.TotalPaid = currency + " " + formatAmount(totalPaid)
			pageData.RemainingBalance = currency + " " + formatAmount(remaining)
			if remaining <= 0 {
				pageData.PaymentStatus = "paid"
			} else if totalPaid > 0 {
				pageData.PaymentStatus = "partial"
			} else {
				pageData.PaymentStatus = "unpaid"
			}
			pageData.Payment = findPayment(allPayments, id, revenue)

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
					EntityType:  "revenue",
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

		templateName := "sales-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

func filterLineItems(all []map[string]any, revenueID string) []map[string]any {
	items := []map[string]any{}
	for _, item := range all {
		rid, _ := item["revenue_id"].(string)
		if rid == revenueID {
			items = append(items, item)
		}
	}
	return items
}

func buildLineItemTable(items []map[string]any, l centymo.RevenueLabels, tableLabels types.TableLabels, currency string) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "description", Label: l.Detail.Description, Sortable: false},
		{Key: "quantity", Label: l.Detail.Quantity, Sortable: false, Width: "80px"},
		{Key: "cost_price", Label: l.Detail.CostPrice, Sortable: false, Width: "130px"},
		{Key: "unit_price", Label: l.Detail.UnitPrice, Sortable: false, Width: "130px"},
		{Key: "discount", Label: l.Detail.Discount, Sortable: false, Width: "100px"},
		{Key: "total", Label: l.Detail.Total, Sortable: false, Width: "130px"},
	}

	rows := []types.TableRow{}
	for _, item := range items {
		id, _ := item["id"].(string)
		description, _ := item["description"].(string)
		quantity, _ := item["quantity"].(string)
		costPrice, _ := item["cost_price"].(string)
		unitPrice, _ := item["unit_price"].(string)
		discount, _ := item["discount"].(string)
		total, _ := item["total"].(string)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: description},
				{Type: "text", Value: quantity},
				{Type: "text", Value: currency + " " + costPrice},
				{Type: "text", Value: currency + " " + unitPrice},
				{Type: "text", Value: discount},
				{Type: "text", Value: currency + " " + total},
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "line-items-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.ItemEmptyTitle,
			Message: l.Detail.ItemEmptyMessage,
		},
	}
}

// buildAuditTable creates the audit trail table.
func buildAuditTable(l centymo.RevenueLabels, tableLabels types.TableLabels) *types.TableConfig {
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

// filterPayments filters payments belonging to a specific revenue.
func filterPayments(all []map[string]any, revenueID string) []map[string]any {
	result := []map[string]any{}
	for _, p := range all {
		rid, _ := p["revenue_id"].(string)
		if rid == revenueID {
			result = append(result, p)
		}
	}
	return result
}

// buildPaymentTable creates the payment table config for the payment tab.
func buildPaymentTable(payments []map[string]any, l centymo.RevenueLabels, tableLabels types.TableLabels, currency string, revenueID string, routes centymo.RevenueRoutes, perms *types.UserPermissions) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "method", Label: l.Detail.PaymentMethod, Sortable: false},
		{Key: "amount", Label: l.Detail.AmountPaid, Sortable: false, Width: "140px"},
		{Key: "reference", Label: l.Detail.Reference, Sortable: false, Width: "160px"},
		{Key: "received_by", Label: l.Detail.ReceivedBy, Sortable: false, Width: "150px"},
		{Key: "date", Label: l.Detail.PaymentDate, Sortable: false, Width: "140px"},
	}

	rows := []types.TableRow{}
	for _, p := range payments {
		id, _ := p["id"].(string)
		method, _ := p["payment_method"].(string)
		amount, _ := p["amount_paid"].(string)
		refNum, _ := p["reference_number"].(string)
		receivedBy, _ := p["received_by"].(string)
		paymentDate, _ := p["payment_date"].(string)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: method},
				{Type: "text", Value: currency + " " + amount},
				{Type: "text", Value: refNum},
				{Type: "text", Value: receivedBy},
				{Type: "text", Value: paymentDate},
			},
			Actions: []types.TableAction{
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.PaymentEditURL, "id", revenueID, "pid", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("invoice", "update"), DisabledTooltip: l.Errors.PermissionDenied},
				{Type: "delete", Label: l.Actions.Delete, Action: "delete", URL: route.ResolveURL(routes.PaymentRemoveURL, "id", revenueID), ItemName: method, Disabled: !perms.Can("invoice", "update"), DisabledTooltip: l.Errors.PermissionDenied},
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "payment-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.PaymentEmptyTitle,
			Message: l.Detail.PaymentEmptyMessage,
		},
	}
}

// sumPayments totals the amount_paid across all payment records.
func sumPayments(payments []map[string]any) float64 {
	total := 0.0
	for _, p := range payments {
		amount, _ := p["amount_paid"].(string)
		total += parseAmount(amount)
	}
	return total
}

// parseAmount converts a string amount to float64.
func parseAmount(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// formatAmount formats a float64 as a comma-separated 2-decimal string.
func formatAmount(f float64) string {
	return centymo.FormatWithCommas(f)
}

// ---------------------------------------------------------------------------
// Proto-to-map conversion helpers
// ---------------------------------------------------------------------------

// revenueToMap converts a Revenue protobuf to a map[string]any for template use.
func revenueToMap(r *revenuepb.Revenue) map[string]any {
	return map[string]any{
		"id":                   r.GetId(),
		"name":                 r.GetName(),
		"client_id":            r.GetClientId(),
		"revenue_date_string":  r.GetRevenueDate(),
		"total_amount":         centymo.FormatWithCommas(float64(r.GetTotalAmount()) / 100.0),
		"currency":             r.GetCurrency(),
		"status":               r.GetStatus(),
		"reference_number":     r.GetReferenceNumber(),
		"notes":                r.GetNotes(),
		"location_id":          r.GetLocationId(),
		"active":               r.GetActive(),
		"date_created_string":  r.GetDateCreatedString(),
		"date_modified_string": r.GetDateModifiedString(),
		"payment_term_id":      r.GetPaymentTermId(),
		"due_date_string":      r.GetDueDate(),
		"payment_term_name": func() string {
			if pt := r.GetPaymentTerm(); pt != nil {
				return pt.GetName()
			}
			return ""
		}(),
		"subscription_id": r.GetSubscriptionId(),
	}
}

// lineItemToMap converts a RevenueLineItem protobuf to a map[string]any for template use.
func lineItemToMap(item *revenuelineitempb.RevenueLineItem) map[string]any {
	return map[string]any{
		"id":                  item.GetId(),
		"revenue_id":          item.GetRevenueId(),
		"description":         item.GetDescription(),
		"quantity":            fmt.Sprintf("%.0f", item.GetQuantity()),
		"unit_price":          centymo.FormatWithCommas(float64(item.GetUnitPrice()) / 100.0),
		"cost_price":          centymo.FormatWithCommas(float64(item.GetCostPrice()) / 100.0),
		"discount":            "0",
		"total":               centymo.FormatWithCommas(float64(item.GetTotalPrice()) / 100.0),
		"line_item_type":      item.GetLineItemType(),
		"inventory_item_id":   item.GetInventoryItemId(),
		"inventory_serial_id": item.GetInventorySerialId(),
		"notes":               item.GetNotes(),
	}
}

// listLineItemMaps lists line items for a revenue via the typed use case and returns maps.
func listLineItemMaps(ctx context.Context, listFn func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error), revenueID string) []map[string]any {
	resp, err := listFn(ctx, &revenuelineitempb.ListRevenueLineItemsRequest{
		RevenueId: &revenueID,
	})
	if err != nil {
		log.Printf("Failed to list line items for revenue %s: %v", revenueID, err)
		return []map[string]any{}
	}
	items := []map[string]any{}
	for _, item := range resp.GetData() {
		if item.GetRevenueId() == revenueID {
			items = append(items, lineItemToMap(item))
		}
	}
	return items
}

// findPayment finds the payment record for a given revenue ID.
func findPayment(payments []map[string]any, revenueID string, revenue map[string]any) *PaymentInfo {
	currency, _ := revenue["currency"].(string)

	for _, p := range payments {
		rid, _ := p["revenue_id"].(string)
		if rid != revenueID {
			continue
		}

		method, _ := p["payment_method"].(string)
		amount, _ := p["amount_paid"].(string)
		cardLast4, _ := p["card_last4"].(string)
		paymentDate, _ := p["payment_date"].(string)
		receivedBy, _ := p["received_by"].(string)
		receivedRole, _ := p["received_role"].(string)

		return &PaymentInfo{
			Method:       method,
			AmountPaid:   currency + " " + amount,
			Currency:     currency,
			CardLast4:    cardLast4,
			PaymentDate:  paymentDate,
			ReceivedBy:   receivedBy,
			ReceivedRole: receivedRole,
		}
	}

	// Fallback: no dedicated payment record — use revenue-level data
	totalAmount, _ := revenue["total_amount"].(string)
	return &PaymentInfo{
		Method:     "—",
		AmountPaid: currency + " " + totalAmount,
		Currency:   currency,
	}
}
