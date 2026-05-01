package list

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
)

// ListViewDeps holds dependencies for the expense_recognition list view.
type ListViewDeps struct {
	Routes                  centymo.ExpenseRecognitionRoutes
	ListExpenseRecognitions func(ctx context.Context, req *expenserecognitionpb.ListExpenseRecognitionsRequest) (*expenserecognitionpb.ListExpenseRecognitionsResponse, error)
	Labels                  centymo.ExpenseRecognitionLabels
	CommonLabels            pyeza.CommonLabels
	TableLabels             types.TableLabels
}

// PageData holds the data for the expense_recognition list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the expense_recognition list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "all"
		}

		resp, err := deps.ListExpenseRecognitions(ctx, &expenserecognitionpb.ListExpenseRecognitionsRequest{})
		if err != nil {
			log.Printf("ListExpenseRecognitions: %v", err)
			return view.Error(fmt.Errorf("failed to load expense recognitions: %w", err))
		}

		recognitions := resp.GetData()
		if status != "all" {
			var filtered []*expenserecognitionpb.ExpenseRecognition
			for _, r := range recognitions {
				if statusKey(r.GetStatus().String()) == status {
					filtered = append(filtered, r)
				}
			}
			recognitions = filtered
		}

		l := deps.Labels
		columns := recognitionColumns(l)
		rows := buildTableRows(recognitions, l)
		types.ApplyColumnStyles(columns, rows)

		tableConfig := &types.TableConfig{
			ID:                   "expense-recognitions-table",
			RefreshURL:           deps.Routes.ListURL,
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
			ContentTemplate: "expense-recognition-list-content",
			Table:           tableConfig,
		}

		return view.OK("expense-recognition-list", pageData)
	})
}

func recognitionColumns(l centymo.ExpenseRecognitionLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
		{Key: "period", Label: l.Columns.PeriodStart, WidthClass: "col-3xl"},
		{Key: "supplier_contract", Label: l.Columns.SupplierContract, WidthClass: "col-3xl"},
		{Key: "total_amount", Label: l.Columns.TotalAmount, WidthClass: "col-3xl", Align: "right"},
		{Key: "idempotency_key", Label: l.Columns.IdempotencyKey, NoSort: true, WidthClass: "col-3xl"},
	}
}

func buildTableRows(items []*expenserecognitionpb.ExpenseRecognition, l centymo.ExpenseRecognitionLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, r := range items {
		id := r.GetId()
		name := r.GetName()
		statusStr := r.GetStatus().String()
		currency := r.GetCurrency()

		// Period — combine period_start / period_end
		period := ""
		if r.PeriodStart != nil {
			period = r.GetPeriodStart().AsTime().Format("2006-01-02")
		}
		if r.PeriodEnd != nil {
			period += " → " + r.GetPeriodEnd().AsTime().Format("2006-01-02")
		}

		contractID := r.GetSupplierContractId()
		idemKey := r.GetIdempotencyKey()

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "badge", Value: statusStr, Variant: recognitionStatusVariant(statusStr)},
				{Type: "text", Value: period},
				{Type: "text", Value: contractID},
				types.MoneyCell(float64(r.GetTotalAmount()), currency, true),
				{Type: "text", Value: idemKey},
			},
			DataAttrs: map[string]string{
				"name":   name,
				"status": statusStr,
			},
		})
	}
	return rows
}

func recognitionStatusVariant(status string) string {
	switch status {
	case "EXPENSE_RECOGNITION_STATUS_DRAFT":
		return "default"
	case "EXPENSE_RECOGNITION_STATUS_POSTED":
		return "success"
	case "EXPENSE_RECOGNITION_STATUS_REVERSED":
		return "danger"
	default:
		return "default"
	}
}

func statusKey(s string) string {
	switch s {
	case "EXPENSE_RECOGNITION_STATUS_DRAFT":
		return "draft"
	case "EXPENSE_RECOGNITION_STATUS_POSTED":
		return "posted"
	case "EXPENSE_RECOGNITION_STATUS_REVERSED":
		return "reversed"
	default:
		return ""
	}
}

func statusPageTitle(l centymo.ExpenseRecognitionLabels, status string) string {
	switch status {
	case "draft":
		return l.Page.HeadingDraft
	case "posted":
		return l.Page.HeadingPosted
	case "reversed":
		return l.Page.HeadingReversed
	default:
		return l.Page.Heading
	}
}
