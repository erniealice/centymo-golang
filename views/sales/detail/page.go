package detail

import (
	"context"
	"fmt"
	"log"
	"strconv"

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
	LineItemTable   *types.TableConfig
	TotalGrossProfit string
	Payment         *PaymentInfo
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
			pageData.LineItemTable = buildLineItemTable(lineItems, l, deps.TableLabels, currency)
			pageData.TotalGrossProfit = calcTotalGrossProfit(lineItems, currency)

		case "payment":
			allPayments, err := deps.DB.ListSimple(ctx, "revenue_payment")
			if err != nil {
				log.Printf("Failed to list payments for revenue %s: %v", id, err)
				allPayments = []map[string]any{}
			}
			pageData.Payment = findPayment(allPayments, id, revenue)

		case "audit":
			// Coming soon — no data to load
		}

		return view.OK("sales-detail", pageData)
	})
}

func buildTabItems(l centymo.SalesLabels, id string) []pyeza.TabItem {
	base := "/app/sales/" + id
	return []pyeza.TabItem{
		{Key: "info", Label: l.Detail.TabBasicInfo, Href: base + "?tab=info", Icon: "icon-info"},
		{Key: "items", Label: l.Detail.TabLineItems, Href: base + "?tab=items", Icon: "icon-list"},
		{Key: "payment", Label: l.Detail.TabPayment, Href: base + "?tab=payment", Icon: "icon-credit-card"},
		{Key: "audit", Label: l.Detail.TabAuditTrail, Href: base + "?tab=audit", Icon: "icon-clock"},
	}
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
		{Key: "gross_profit", Label: l.Detail.GrossProfit, Sortable: false, Width: "130px"},
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

		grossProfit := calcGrossProfit(unitPrice, costPrice, quantity, discount)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: description},
				{Type: "text", Value: quantity},
				{Type: "text", Value: currency + " " + costPrice},
				{Type: "text", Value: currency + " " + unitPrice},
				{Type: "text", Value: discount},
				{Type: "text", Value: currency + " " + total},
				{Type: "text", Value: currency + " " + grossProfit},
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

// calcGrossProfit computes (unitPrice - costPrice) * quantity - discount for a single line item.
func calcGrossProfit(unitPriceStr, costPriceStr, quantityStr, discountStr string) string {
	unitPrice, _ := strconv.ParseFloat(unitPriceStr, 64)
	costPrice, _ := strconv.ParseFloat(costPriceStr, 64)
	quantity, _ := strconv.ParseFloat(quantityStr, 64)
	discount, _ := strconv.ParseFloat(discountStr, 64)

	gp := (unitPrice-costPrice)*quantity - discount
	return fmt.Sprintf("%.2f", gp)
}

// calcTotalGrossProfit sums gross profit across all line items.
func calcTotalGrossProfit(items []map[string]any, currency string) string {
	var total float64
	for _, item := range items {
		unitPrice, _ := item["unit_price"].(string)
		costPrice, _ := item["cost_price"].(string)
		quantity, _ := item["quantity"].(string)
		discount, _ := item["discount"].(string)

		up, _ := strconv.ParseFloat(unitPrice, 64)
		cp, _ := strconv.ParseFloat(costPrice, 64)
		q, _ := strconv.ParseFloat(quantity, 64)
		d, _ := strconv.ParseFloat(discount, 64)

		total += (up-cp)*q - d
	}
	return currency + " " + fmt.Sprintf("%.2f", total)
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
