package list

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"
)

// ListViewDeps holds dependencies for the accrued_expense list view.
type ListViewDeps struct {
	Routes              centymo.AccruedExpenseRoutes
	ListAccruedExpenses func(ctx context.Context, req *accruedexpensepb.ListAccruedExpensesRequest) (*accruedexpensepb.ListAccruedExpensesResponse, error)
	Labels              centymo.AccruedExpenseLabels
	CommonLabels        pyeza.CommonLabels
	TableLabels         types.TableLabels
}

// PageData holds the data for the accrued_expense list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the accrued_expense list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "all"
		}

		resp, err := deps.ListAccruedExpenses(ctx, &accruedexpensepb.ListAccruedExpensesRequest{})
		if err != nil {
			log.Printf("ListAccruedExpenses: %v", err)
			return view.Error(fmt.Errorf("failed to load accrued expenses: %w", err))
		}

		items := resp.GetData()
		if status != "all" {
			var filtered []*accruedexpensepb.AccruedExpense
			for _, a := range items {
				if statusKey(a.GetStatus().String()) == status {
					filtered = append(filtered, a)
				}
			}
			items = filtered
		}

		l := deps.Labels
		columns := accrualColumns(l)
		rows := buildTableRows(items)
		types.ApplyColumnStyles(columns, rows)

		var primaryAction *types.PrimaryAction
		if deps.Routes.AddURL != "" {
			primaryAction = &types.PrimaryAction{
				Label:     l.Buttons.Add,
				ActionURL: deps.Routes.AddURL,
			}
		}

		tableConfig := &types.TableConfig{
			ID:                   "accrued-expenses-table",
			RefreshURL:           deps.Routes.ListURL,
			Columns:              columns,
			Rows:                 rows,
			PrimaryAction:        primaryAction,
			ShowSearch:           true,
			ShowActions:          true,
			ShowFilters:          true,
			ShowSort:             true,
			ShowColumns:          true,
			ShowExport:           true,
			ShowDensity:          true,
			ShowEntries:          true,
			DefaultSortColumn:    "recognition_date",
			DefaultSortDirection: "desc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   l.Empty.Title,
				Message: l.Empty.Message,
			},
		}
		types.ApplyTableSettings(tableConfig)

		heading := statusPageTitle(l, status)
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   status,
				HeaderTitle:    heading,
				HeaderSubtitle: l.Page.Caption,
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "accrued-expense-list-content",
			Table:           tableConfig,
		}

		return view.OK("accrued-expense-list", pageData)
	})
}

func accrualColumns(l centymo.AccruedExpenseLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
		{Key: "period", Label: l.Columns.PeriodStart, WidthClass: "col-3xl"},
		{Key: "supplier_contract", Label: l.Columns.SupplierContract, WidthClass: "col-3xl"},
		{Key: "accrued_amount", Label: l.Columns.AccruedAmount, WidthClass: "col-3xl", Align: "right"},
		{Key: "settled_amount", Label: l.Columns.SettledAmount, WidthClass: "col-3xl", Align: "right"},
		{Key: "remaining_amount", Label: l.Columns.RemainingAmount, WidthClass: "col-3xl", Align: "right"},
	}
}

func buildTableRows(items []*accruedexpensepb.AccruedExpense) []types.TableRow {
	rows := []types.TableRow{}
	for _, a := range items {
		id := a.GetId()
		name := a.GetName()
		statusStr := a.GetStatus().String()
		currency := a.GetCurrency()

		period := ""
		if a.PeriodStart != nil {
			period = a.GetPeriodStart().AsTime().Format("2006-01-02")
		}
		if a.PeriodEnd != nil {
			period += " → " + a.GetPeriodEnd().AsTime().Format("2006-01-02")
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "badge", Value: statusStr, Variant: accrualStatusVariant(statusStr)},
				{Type: "text", Value: period},
				{Type: "text", Value: a.GetSupplierContractId()},
				types.MoneyCell(float64(a.GetAccruedAmount()), currency, true),
				types.MoneyCell(float64(a.GetSettledAmount()), currency, true),
				types.MoneyCell(float64(a.GetRemainingAmount()), currency, true),
			},
			DataAttrs: map[string]string{
				"name":   name,
				"status": statusStr,
			},
		})
	}
	return rows
}

func accrualStatusVariant(status string) string {
	switch status {
	case "ACCRUED_EXPENSE_STATUS_OUTSTANDING":
		return "warning"
	case "ACCRUED_EXPENSE_STATUS_PARTIAL":
		return "info"
	case "ACCRUED_EXPENSE_STATUS_SETTLED":
		return "success"
	case "ACCRUED_EXPENSE_STATUS_REVERSED":
		return "danger"
	default:
		return "default"
	}
}

func statusKey(s string) string {
	switch s {
	case "ACCRUED_EXPENSE_STATUS_OUTSTANDING":
		return "outstanding"
	case "ACCRUED_EXPENSE_STATUS_PARTIAL":
		return "partial"
	case "ACCRUED_EXPENSE_STATUS_SETTLED":
		return "settled"
	case "ACCRUED_EXPENSE_STATUS_REVERSED":
		return "reversed"
	default:
		return ""
	}
}

func statusPageTitle(l centymo.AccruedExpenseLabels, status string) string {
	switch status {
	case "outstanding":
		return l.Page.HeadingOutstanding
	case "partial":
		return l.Page.HeadingPartial
	case "settled":
		return l.Page.HeadingSettled
	case "reversed":
		return l.Page.HeadingReversed
	default:
		return l.Page.Heading
	}
}
