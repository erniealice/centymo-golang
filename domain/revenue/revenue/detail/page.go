package detail

import (
	"context"
	"fmt"
	revenuedomain "github.com/erniealice/centymo-golang/domain/revenue/revenue"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	"log"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
	revenuepaymentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_payment"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes       revenuedomain.Routes
	Labels       revenuedomain.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Typed revenue operations
	ReadRevenue func(ctx context.Context, req *revenuepb.ReadRevenueRequest) (*revenuepb.ReadRevenueResponse, error)

	// Typed line item operations
	ListRevenueLineItems func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error)

	// Typed revenue_payment list (payment tab). 20260612-datasource-typed-path
	// W5 — replaces DataSource ListSimple("revenue_payment"). Optional —
	// nil-safe (renders an empty payment table).
	ListRevenuePayments func(ctx context.Context, req *revenuepaymentpb.ListRevenuePaymentsRequest) (*revenuepaymentpb.ListRevenuePaymentsResponse, error)

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
	ContentTemplate      string
	Revenue              map[string]any
	Labels               revenuedomain.Labels
	ActiveTab            string
	TabItems             []pyeza.TabItem
	LineItemTable        *types.TableConfig
	LineItemAddURL       string
	LineItemDiscountURL  string
	TotalAmount          types.TableCell
	Payment              *PaymentInfo
	PaymentTable         *types.TableConfig
	PaymentAddURL        string
	TotalPaid            string
	RemainingBalance     string
	PaymentStatus        string
	PaymentStatusVariant string
	AuditTable           *types.TableConfig
	AttachmentTable      *types.TableConfig
	InvoiceDownloadURL   string
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// NewView creates the sales detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "read") {
			return view.Forbidden("invoice:read")
		}
		_ = perms
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
				ActiveNav:      "revenue",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.PageTitle,
				HeaderIcon:     "icon-shopping-bag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:    "revenue-detail-content",
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
			currency, _ := revenue["currency"].(string)
			lineItems := listLineItemMaps(ctx, deps.ListRevenueLineItems, id, currency)
			pageData.LineItemTable = buildLineItemTableWithActions(lineItems, l, deps.TableLabels, currency, id, deps.Routes, perms)
			pageData.LineItemAddURL = route.ResolveURL(deps.Routes.LineItemAddURL, "id", id)
			pageData.LineItemDiscountURL = route.ResolveURL(deps.Routes.LineItemDiscountURL, "id", id)
			totalAmountCell, _ := revenue["total_amount"].(types.TableCell)
			pageData.TotalAmount = totalAmountCell

		case "payment":
			payments := listRevenuePayments(ctx, deps.ListRevenuePayments, id)
			currency, _ := revenue["currency"].(string)
			perms := view.GetUserPermissions(ctx)
			pageData.PaymentTable = buildPaymentTable(payments, l, deps.TableLabels, currency, id, deps.Routes, perms)
			pageData.PaymentAddURL = route.ResolveURL(deps.Routes.PaymentAddURL, "id", id)

			// Calculate totals
			totalCentavos, _ := revenue["total_amount_centavos"].(int64)
			paidCentavos := sumPaymentsCentavos(payments)
			remainingCentavos := totalCentavos - paidCentavos

			pageData.TotalPaid = types.FormatMoney(paidCentavos, currency)
			pageData.RemainingBalance = types.FormatMoney(remainingCentavos, currency)
			if remainingCentavos <= 0 {
				pageData.PaymentStatus = "paid"
				pageData.PaymentStatusVariant = "success"
			} else if paidCentavos > 0 {
				pageData.PaymentStatus = "partial"
				pageData.PaymentStatusVariant = "warning"
			} else {
				pageData.PaymentStatus = "unpaid"
				pageData.PaymentStatusVariant = "info"
			}
			// Keep legacy field for backward compat
			pageData.Payment = findPayment(payments, id, revenue)

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
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "revenue-detail"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("revenue-detail", pageData)
	})
}

func buildTabItems(l revenuedomain.Labels, id string, routes revenuedomain.Routes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Detail.TabBasicInfo, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "items", Label: l.Detail.TabLineItems, Href: base + "?tab=items", HxGet: action + "items", Icon: "icon-list"},
		{Key: "payment", Label: l.Detail.TabPayment, Href: base + "?tab=payment", HxGet: action + "payment", Icon: "icon-credit-card"},
		{Key: "audit", Label: l.Detail.TabAuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
		{Key: "attachments", Label: l.Detail.TabAttachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		{Key: "audit-history", Label: l.Detail.TabAuditHistory, Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("invoice", "read") {
			return view.Forbidden("invoice:read")
		}
		_ = perms
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
			currency, _ := revenue["currency"].(string)
			lineItems := listLineItemMaps(ctx, deps.ListRevenueLineItems, id, currency)
			pageData.LineItemTable = buildLineItemTableWithActions(lineItems, l, deps.TableLabels, currency, id, deps.Routes, perms)
			pageData.LineItemAddURL = route.ResolveURL(deps.Routes.LineItemAddURL, "id", id)
			pageData.LineItemDiscountURL = route.ResolveURL(deps.Routes.LineItemDiscountURL, "id", id)
			totalAmountCell, _ := revenue["total_amount"].(types.TableCell)
			pageData.TotalAmount = totalAmountCell

		case "payment":
			payments := listRevenuePayments(ctx, deps.ListRevenuePayments, id)
			currency, _ := revenue["currency"].(string)
			perms := view.GetUserPermissions(ctx)
			pageData.PaymentTable = buildPaymentTable(payments, l, deps.TableLabels, currency, id, deps.Routes, perms)
			pageData.PaymentAddURL = route.ResolveURL(deps.Routes.PaymentAddURL, "id", id)

			totalCentavos, _ := revenue["total_amount_centavos"].(int64)
			paidCentavos := sumPaymentsCentavos(payments)
			remainingCentavos := totalCentavos - paidCentavos

			pageData.TotalPaid = types.FormatMoney(paidCentavos, currency)
			pageData.RemainingBalance = types.FormatMoney(remainingCentavos, currency)
			if remainingCentavos <= 0 {
				pageData.PaymentStatus = "paid"
				pageData.PaymentStatusVariant = "success"
			} else if paidCentavos > 0 {
				pageData.PaymentStatus = "partial"
				pageData.PaymentStatusVariant = "warning"
			} else {
				pageData.PaymentStatus = "unpaid"
				pageData.PaymentStatusVariant = "info"
			}
			pageData.Payment = findPayment(payments, id, revenue)

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

		templateName := "revenue-tab-" + tab
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

func buildLineItemTable(items []map[string]any, l revenuedomain.Labels, tableLabels types.TableLabels, currency string) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "description", Label: l.Detail.Description, NoSort: true},
		{Key: "quantity", Label: l.Detail.Quantity, NoSort: true, WidthClass: "col-md"},
		{Key: "cost_price", Label: l.Detail.CostPrice, NoSort: true, WidthClass: "col-3xl"},
		{Key: "unit_price", Label: l.Detail.UnitPrice, NoSort: true, WidthClass: "col-3xl"},
		{Key: "discount", Label: l.Detail.Discount, NoSort: true, WidthClass: "col-lg"},
		{Key: "total", Label: l.Detail.Total, NoSort: true, WidthClass: "col-3xl"},
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
func buildAuditTable(l revenuedomain.Labels, tableLabels types.TableLabels) *types.TableConfig {
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

// listRevenuePayments fetches the revenue's payments via the typed use case
// with a SERVER-SIDE revenue_id filter. 20260612-datasource-typed-path W5 —
// replaces DataSource ListSimple("revenue_payment") + client filter. Nil-safe:
// returns an empty slice when the closure is unwired (mock builds / half-wired
// composition root), so the tab renders an empty payment table rather than
// crashing.
func listRevenuePayments(ctx context.Context, list func(ctx context.Context, req *revenuepaymentpb.ListRevenuePaymentsRequest) (*revenuepaymentpb.ListRevenuePaymentsResponse, error), revenueID string) []*revenuepaymentpb.RevenuePayment {
	if list == nil {
		return []*revenuepaymentpb.RevenuePayment{}
	}
	resp, err := list(ctx, &revenuepaymentpb.ListRevenuePaymentsRequest{
		Filters: &commonpb.FilterRequest{
			Filters: []*commonpb.TypedFilter{{
				Field: "revenue_id",
				FilterType: &commonpb.TypedFilter_StringFilter{
					StringFilter: &commonpb.StringFilter{
						Value:    revenueID,
						Operator: commonpb.StringOperator_STRING_EQUALS,
					},
				},
			}},
		},
	})
	if err != nil {
		log.Printf("Failed to list payments for revenue %s: %v", revenueID, err)
		return []*revenuepaymentpb.RevenuePayment{}
	}
	return filterPayments(resp.GetData(), revenueID)
}

// filterPayments filters payments belonging to a specific revenue. Defensive
// client-side re-filter behind the server-side revenue_id filter (a partial
// adapter may not honour it) — behaviour-preserving with the prior client filter.
func filterPayments(all []*revenuepaymentpb.RevenuePayment, revenueID string) []*revenuepaymentpb.RevenuePayment {
	result := []*revenuepaymentpb.RevenuePayment{}
	for _, p := range all {
		if p.GetRevenueId() == revenueID {
			result = append(result, p)
		}
	}
	return result
}

// buildPaymentTable creates the payment table config for the payment tab.
func buildPaymentTable(payments []*revenuepaymentpb.RevenuePayment, l revenuedomain.Labels, tableLabels types.TableLabels, currency string, revenueID string, routes revenuedomain.Routes, perms *types.UserPermissions) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "method", Label: l.Detail.PaymentMethod, NoSort: true},
		{Key: "amount", Label: l.Detail.AmountPaid, NoSort: true, WidthClass: "col-3xl"},
		{Key: "reference", Label: l.Detail.Reference, NoSort: true, WidthClass: "col-5xl"},
		{Key: "received_by", Label: l.Detail.ReceivedBy, NoSort: true, WidthClass: "col-4xl"},
		{Key: "date", Label: l.Detail.PaymentDate, NoSort: true, WidthClass: "col-3xl"},
	}

	rows := []types.TableRow{}
	for _, p := range payments {
		id := p.GetId()
		method := p.GetPaymentMethod()
		// amount is centavos (int64, Rule #1); FormatMoney renders ÷100 with the
		// "<currency> <amount>" prefix, matching the prior currency+" "+amount.
		amountDisplay := types.FormatMoney(p.GetAmount(), currency)
		refNum := p.GetReferenceNumber()
		receivedBy := p.GetReceivedBy()
		paymentDate := p.GetPaymentDate()

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: method},
				{Type: "text", Value: amountDisplay},
				{Type: "text", Value: refNum},
				{Type: "text", Value: receivedBy},
				{Type: "text", Value: paymentDate},
			},
			Actions: []types.TableAction{
				{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.PaymentEditURL, "id", revenueID, "pid", id), DrawerTitle: l.Actions.Edit, Disabled: !perms.Can("payment", "update"), DisabledTooltip: l.Errors.PermissionDenied},
				{Type: "delete", Label: l.Actions.Delete, Action: "delete", URL: route.ResolveURL(routes.PaymentRemoveURL, "id", revenueID), ItemName: method, Disabled: !perms.Can("payment", "delete"), DisabledTooltip: l.Errors.PermissionDenied},
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:         "payment-table",
		Columns:    columns,
		Rows:       rows,
		Labels:     tableLabels,
		RefreshURL: route.ResolveURL(routes.PaymentTableURL, "id", revenueID),
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.PaymentEmptyTitle,
			Message: l.Detail.PaymentEmptyMessage,
		},
	}
}

// sumPaymentsCentavos totals the payment amount (int64 centavos, Rule #1) across
// all records. Kept in centavos end-to-end — no float/display-string round-trip.
func sumPaymentsCentavos(payments []*revenuepaymentpb.RevenuePayment) int64 {
	var totalCentavos int64
	for _, p := range payments {
		totalCentavos += p.GetAmount()
	}
	return totalCentavos
}

// ---------------------------------------------------------------------------
// Proto-to-map conversion helpers
// ---------------------------------------------------------------------------

// revenueToMap converts a Revenue protobuf to a map[string]any for template use.
func revenueToMap(r *revenuepb.Revenue) map[string]any {
	return map[string]any{
		"id":                    r.GetId(),
		"name":                  r.GetName(),
		"client_id":             r.GetClientId(),
		"revenue_date_string":   r.GetRevenueDate(),
		"total_amount":          types.MoneyCell(float64(r.GetTotalAmount()), r.GetCurrency(), true),
		"total_amount_centavos": r.GetTotalAmount(),
		"currency":              r.GetCurrency(),
		"status":                r.GetStatus(),
		"reference_number":      r.GetReferenceNumber(),
		"notes":                 r.GetNotes(),
		"location_id":           r.GetLocationId(),
		"active":                r.GetActive(),
		"date_created_string":   r.GetDateCreatedString(),
		"date_modified_string":  r.GetDateModifiedString(),
		"payment_term_id":       r.GetPaymentTermId(),
		"due_date_string":       r.GetDueDate(),
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
func lineItemToMap(item *revenuelineitempb.RevenueLineItem, currency string) map[string]any {
	return map[string]any{
		"id":                  item.GetId(),
		"revenue_id":          item.GetRevenueId(),
		"description":         item.GetDescription(),
		"quantity":            fmt.Sprintf("%.0f", item.GetQuantity()),
		"unit_price":          types.MoneyCell(float64(item.GetUnitPrice()), currency, true),
		"cost_price":          types.MoneyCell(float64(item.GetCostPrice()), currency, true),
		"discount":            "0",
		"total":               types.MoneyCell(float64(item.GetTotalPrice()), currency, true),
		"line_item_type":      item.GetLineItemType(),
		"inventory_item_id":   item.GetInventoryItemId(),
		"inventory_serial_id": item.GetInventorySerialId(),
		"notes":               item.GetNotes(),
	}
}

// listLineItemMaps lists line items for a revenue via the typed use case and returns maps.
func listLineItemMaps(ctx context.Context, listFn func(ctx context.Context, req *revenuelineitempb.ListRevenueLineItemsRequest) (*revenuelineitempb.ListRevenueLineItemsResponse, error), revenueID string, currency string) []map[string]any {
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
			items = append(items, lineItemToMap(item, currency))
		}
	}
	return items
}

// findPayment finds the payment record for a given revenue ID.
// 20260612-datasource-typed-path W5 — reads proto getters, not map keys.
// Note: the proto carries no card_last4 field (it never existed on the
// revenue_payment record), so CardLast4 is left empty — behaviour-preserving
// (the old map key was always absent on a revenue_payment row).
func findPayment(payments []*revenuepaymentpb.RevenuePayment, revenueID string, revenue map[string]any) *PaymentInfo {
	currency, _ := revenue["currency"].(string)

	for _, p := range payments {
		if p.GetRevenueId() != revenueID {
			continue
		}

		// amount is centavos (int64, Rule #1); FormatMoney renders ÷100 with the
		// "<currency> <amount>" prefix, matching the prior currency+" "+amount.
		return &PaymentInfo{
			Method:       p.GetPaymentMethod(),
			AmountPaid:   types.FormatMoney(p.GetAmount(), currency),
			Currency:     currency,
			PaymentDate:  p.GetPaymentDate(),
			ReceivedBy:   p.GetReceivedBy(),
			ReceivedRole: p.GetReceivedRole(),
		}
	}

	// Fallback: no dedicated payment record — use revenue-level data
	totalAmountCell, _ := revenue["total_amount"].(types.TableCell)
	return &PaymentInfo{
		Method:     "—",
		AmountPaid: currency + " " + totalAmountCell.Value,
		Currency:   currency,
	}
}
