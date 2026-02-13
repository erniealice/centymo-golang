package detail

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	DB           centymo.DataSource
	Labels       centymo.SalesLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
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
	ContentTemplate string
	Revenue         map[string]any
	Labels          centymo.SalesLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	LineItemTable       *types.TableConfig
	LineItemAddURL      string
	LineItemDiscountURL string
	TotalAmount         string
	Payment             *PaymentInfo
	PaymentTable     *types.TableConfig
	PaymentAddURL    string
	TotalPaid        string
	RemainingBalance string
	PaymentStatus    string
	AuditTable       *types.TableConfig
}

// NewView creates the sales detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		revenue, err := deps.DB.Read(ctx, "revenue", id)
		if err != nil {
			log.Printf("Failed to read revenue %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load sale: %w", err))
		}

		refNumber, _ := revenue["reference_number"].(string)
		headerTitle := "Sale #" + refNumber

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}

		l := deps.Labels
		tabItems := buildTabItems(l, id)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "sales",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.PageTitle,
				HeaderIcon:     "icon-shopping-bag",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "sales-detail-content",
			Revenue:         revenue,
			Labels:          l,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
		}

		// Load tab-specific data
		switch activeTab {
		case "info":
			// No extra data needed — revenue map has everything
		case "items":
			allLineItems, err := deps.DB.ListSimple(ctx, "revenue_line_item")
			if err != nil {
				log.Printf("Failed to list line items for revenue %s: %v", id, err)
				allLineItems = []map[string]any{}
			}
			lineItems := filterLineItems(allLineItems, id)
			currency, _ := revenue["currency"].(string)
			pageData.LineItemTable = buildLineItemTableWithActions(lineItems, l, deps.TableLabels, currency, id)
			pageData.LineItemAddURL = fmt.Sprintf("/action/sales/detail/%s/items/add", id)
			pageData.LineItemDiscountURL = fmt.Sprintf("/action/sales/detail/%s/items/add-discount", id)
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
			pageData.PaymentTable = buildPaymentTable(payments, l, deps.TableLabels, currency, id)
			pageData.PaymentAddURL = fmt.Sprintf("/action/sales/detail/%s/payment/add", id)

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
		}

		return view.OK("sales-detail", pageData)
	})
}

func buildTabItems(l centymo.SalesLabels, id string) []pyeza.TabItem {
	base := "/app/sales/detail/" + id
	action := "/action/sales/detail/" + id + "/tab/"
	return []pyeza.TabItem{
		{Key: "info", Label: l.Detail.TabBasicInfo, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "items", Label: l.Detail.TabLineItems, Href: base + "?tab=items", HxGet: action + "items", Icon: "icon-list"},
		{Key: "payment", Label: l.Detail.TabPayment, Href: base + "?tab=payment", HxGet: action + "payment", Icon: "icon-credit-card"},
		{Key: "audit", Label: l.Detail.TabAuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
	}
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		revenue, err := deps.DB.Read(ctx, "revenue", id)
		if err != nil {
			log.Printf("Failed to read revenue %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load sale: %w", err))
		}

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Revenue:   revenue,
			Labels:    l,
			ActiveTab: tab,
			TabItems:  buildTabItems(l, id),
		}

		switch tab {
		case "info":
			// revenue map has everything
		case "items":
			allLineItems, err := deps.DB.ListSimple(ctx, "revenue_line_item")
			if err != nil {
				log.Printf("Failed to list line items for revenue %s: %v", id, err)
				allLineItems = []map[string]any{}
			}
			lineItems := filterLineItems(allLineItems, id)
			currency, _ := revenue["currency"].(string)
			pageData.LineItemTable = buildLineItemTableWithActions(lineItems, l, deps.TableLabels, currency, id)
			pageData.LineItemAddURL = fmt.Sprintf("/action/sales/detail/%s/items/add", id)
			pageData.LineItemDiscountURL = fmt.Sprintf("/action/sales/detail/%s/items/add-discount", id)
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
			pageData.PaymentTable = buildPaymentTable(payments, l, deps.TableLabels, currency, id)
			pageData.PaymentAddURL = fmt.Sprintf("/action/sales/detail/%s/payment/add", id)

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
		}

		templateName := "sales-tab-" + tab
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

func buildLineItemTable(items []map[string]any, l centymo.SalesLabels, tableLabels types.TableLabels, currency string) *types.TableConfig {
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
			Title:   "No line items",
			Message: "This sale has no line items.",
		},
	}
}

// buildAuditTable creates the audit trail table.
func buildAuditTable(l centymo.SalesLabels, tableLabels types.TableLabels) *types.TableConfig {
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
func buildPaymentTable(payments []map[string]any, l centymo.SalesLabels, tableLabels types.TableLabels, currency string, revenueID string) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "method", Label: l.Detail.PaymentMethod, Sortable: false},
		{Key: "amount", Label: l.Detail.AmountPaid, Sortable: false, Width: "140px"},
		{Key: "reference", Label: "Reference", Sortable: false, Width: "160px"},
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
				{Type: "edit", Label: "Edit", Action: "edit", URL: fmt.Sprintf("/action/sales/detail/%s/payment/edit/%s", revenueID, id), DrawerTitle: "Edit Payment"},
				{Type: "delete", Label: "Remove", Action: "delete", URL: fmt.Sprintf("/action/sales/detail/%s/payment/remove", revenueID), ItemName: method},
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
			Title:   "No payments recorded",
			Message: "Record a payment to track collections for this sale.",
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

// formatAmount formats a float64 as a 2-decimal string.
func formatAmount(f float64) string {
	f = math.Round(f*100) / 100
	return strconv.FormatFloat(f, 'f', 2, 64)
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
