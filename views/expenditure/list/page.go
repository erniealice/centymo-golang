package list

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
	DB              centymo.DataSource
	RefreshURL      string
	ExpenditureType string // "purchase" or "expense" — determines which type to filter
	Labels          centymo.ExpenditureLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
}

// PageData holds the data for the expenditure list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the expenditure list view, filtered by type (purchase or expense).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "all"
		}

		records, err := deps.DB.ListSimple(ctx, "expenditure")
		if err != nil {
			log.Printf("Failed to list expenditures: %v", err)
			return view.Error(fmt.Errorf("failed to load expenditures: %w", err))
		}

		// Filter by expenditure_type
		var filtered []map[string]any
		for _, r := range records {
			if t, ok := r["expenditure_type"].(string); ok && t == deps.ExpenditureType {
				filtered = append(filtered, r)
			}
		}

		// Further filter by status if not "all"
		if status != "all" {
			var statusFiltered []map[string]any
			for _, r := range filtered {
				if s, ok := r["status"].(string); ok && s == status {
					statusFiltered = append(statusFiltered, r)
				}
			}
			filtered = statusFiltered
		}

		l := deps.Labels
		columns := expenditureColumns(l, deps.ExpenditureType)
		rows := buildTableRows(filtered, l)
		types.ApplyColumnStyles(columns, rows)

		tableID := "purchases-table"
		activeNav := "purchases"
		heading := statusPageTitle(l, deps.ExpenditureType, status)
		caption := statusPageCaption(l, deps.ExpenditureType, status)
		icon := "icon-shopping-bag"
		emptyTitle := statusEmptyTitle(l, deps.ExpenditureType, status)
		emptyMessage := statusEmptyMessage(l, deps.ExpenditureType, status)

		if deps.ExpenditureType == "expense" {
			tableID = "expenses-table"
			activeNav = "expenses"
			icon = "icon-file-minus"
		}

		tableConfig := &types.TableConfig{
			ID:                   tableID,
			RefreshURL:           deps.RefreshURL,
			Columns:              columns,
			Rows:                 rows,
			ShowSearch:           true,
			ShowActions:          true,
			ShowFilters:          true,
			ShowSort:             true,
			ShowColumns:          true,
			ShowExport:           true,
			ShowDensity:          true,
			ShowEntries:          true,
			DefaultSortColumn:    "date",
			DefaultSortDirection: "desc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   emptyTitle,
				Message: emptyMessage,
			},
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      activeNav,
				ActiveSubNav:   status,
				HeaderTitle:    heading,
				HeaderSubtitle: caption,
				HeaderIcon:     icon,
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "expenditure-list-content",
			Table:           tableConfig,
		}

		return view.OK("expenditure-list", pageData)
	})
}

func expenditureColumns(l centymo.ExpenditureLabels, expenditureType string) []types.TableColumn {
	cols := []types.TableColumn{
		{Key: "reference", Label: l.Columns.Reference, Sortable: true},
		{Key: "vendor", Label: l.Columns.Vendor, Sortable: true},
		{Key: "date", Label: l.Columns.Date, Sortable: true, Width: "140px"},
		{Key: "amount", Label: l.Columns.Amount, Sortable: true, Width: "140px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
	if expenditureType == "expense" {
		// Replace vendor with category for expenses
		cols[1] = types.TableColumn{Key: "category", Label: l.Columns.Category, Sortable: true}
	}
	return cols
}

func buildTableRows(records []map[string]any, l centymo.ExpenditureLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, record := range records {
		id, _ := record["id"].(string)
		refNumber, _ := record["reference_number"].(string)
		vendor, _ := record["vendor_name"].(string)
		date, _ := record["expenditure_date_string"].(string)
		currency, _ := record["currency"].(string)
		recordStatus, _ := record["status"].(string)
		category, _ := record["category"].(string)

		amountDisplay := currency + " " + formatAmount(record["total_amount"])

		// Second column is vendor or category depending on type
		expType, _ := record["expenditure_type"].(string)
		secondCol := vendor
		if expType == "expense" {
			secondCol = category
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: refNumber},
				{Type: "text", Value: secondCol},
				{Type: "text", Value: date},
				{Type: "text", Value: amountDisplay},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"reference": refNumber,
				"vendor":    vendor,
				"date":      date,
				"amount":    amountDisplay,
				"status":    recordStatus,
			},
		})
	}
	return rows
}

func statusPageTitle(l centymo.ExpenditureLabels, expType, status string) string {
	if expType == "purchase" {
		switch status {
		case "draft":
			return l.Page.PurchaseHeadingDraft
		case "pending":
			return l.Page.PurchaseHeadingPending
		case "approved":
			return l.Page.PurchaseHeadingApproved
		case "paid":
			return l.Page.PurchaseHeadingPaid
		case "cancelled":
			return l.Page.PurchaseHeadingCancelled
		case "overdue":
			return l.Page.PurchaseHeadingOverdue
		default:
			return l.Page.PurchaseHeading
		}
	}
	// expense
	switch status {
	case "draft":
		return l.Page.ExpenseHeadingDraft
	case "pending":
		return l.Page.ExpenseHeadingPending
	case "approved":
		return l.Page.ExpenseHeadingApproved
	case "paid":
		return l.Page.ExpenseHeadingPaid
	case "cancelled":
		return l.Page.ExpenseHeadingCancelled
	case "overdue":
		return l.Page.ExpenseHeadingOverdue
	default:
		return l.Page.ExpenseHeading
	}
}

func statusPageCaption(l centymo.ExpenditureLabels, expType, status string) string {
	if expType == "purchase" {
		return l.Page.PurchaseCaption
	}
	return l.Page.ExpenseCaption
}

func statusEmptyTitle(l centymo.ExpenditureLabels, expType, status string) string {
	if expType == "purchase" {
		switch status {
		case "draft":
			return l.Empty.PurchaseDraftTitle
		case "pending":
			return l.Empty.PurchasePendingTitle
		case "approved":
			return l.Empty.PurchaseApprovedTitle
		case "paid":
			return l.Empty.PurchasePaidTitle
		case "cancelled":
			return l.Empty.PurchaseCancelledTitle
		case "overdue":
			return l.Empty.PurchaseOverdueTitle
		default:
			return l.Empty.PurchaseTitle
		}
	}
	// expense
	switch status {
	case "draft":
		return l.Empty.ExpenseDraftTitle
	case "pending":
		return l.Empty.ExpensePendingTitle
	case "approved":
		return l.Empty.ExpenseApprovedTitle
	case "paid":
		return l.Empty.ExpensePaidTitle
	case "cancelled":
		return l.Empty.ExpenseCancelledTitle
	case "overdue":
		return l.Empty.ExpenseOverdueTitle
	default:
		return l.Empty.ExpenseTitle
	}
}

func statusEmptyMessage(l centymo.ExpenditureLabels, expType, status string) string {
	if expType == "purchase" {
		switch status {
		case "draft":
			return l.Empty.PurchaseDraftMessage
		case "pending":
			return l.Empty.PurchasePendingMessage
		case "approved":
			return l.Empty.PurchaseApprovedMessage
		case "paid":
			return l.Empty.PurchasePaidMessage
		case "cancelled":
			return l.Empty.PurchaseCancelledMessage
		case "overdue":
			return l.Empty.PurchaseOverdueMessage
		default:
			return l.Empty.PurchaseMessage
		}
	}
	// expense
	switch status {
	case "draft":
		return l.Empty.ExpenseDraftMessage
	case "pending":
		return l.Empty.ExpensePendingMessage
	case "approved":
		return l.Empty.ExpenseApprovedMessage
	case "paid":
		return l.Empty.ExpensePaidMessage
	case "cancelled":
		return l.Empty.ExpenseCancelledMessage
	case "overdue":
		return l.Empty.ExpenseOverdueMessage
	default:
		return l.Empty.ExpenseMessage
	}
}

func formatAmount(v any) string {
	switch n := v.(type) {
	case float64:
		return fmt.Sprintf("%.2f", n)
	case float32:
		return fmt.Sprintf("%.2f", n)
	case int64:
		return fmt.Sprintf("%d", n)
	case int:
		return fmt.Sprintf("%d", n)
	case string:
		return n
	default:
		return fmt.Sprintf("%v", v)
	}
}

func statusVariant(status string) string {
	switch status {
	case "draft":
		return "default"
	case "pending":
		return "warning"
	case "approved":
		return "info"
	case "paid":
		return "success"
	case "cancelled":
		return "danger"
	case "overdue":
		return "danger"
	default:
		return "default"
	}
}
