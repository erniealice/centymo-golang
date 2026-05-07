package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/hybra-golang/views/attachment"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
	expenserecognitionlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition_line"
)

// DetailViewDeps holds dependencies for the expense_recognition detail page.
type DetailViewDeps struct {
	Routes       centymo.ExpenseRecognitionRoutes
	Labels       centymo.ExpenseRecognitionLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadExpenseRecognition       func(ctx context.Context, req *expenserecognitionpb.ReadExpenseRecognitionRequest) (*expenserecognitionpb.ReadExpenseRecognitionResponse, error)
	ListExpenseRecognitionLines  func(ctx context.Context, req *expenserecognitionlinepb.ListExpenseRecognitionLinesRequest) (*expenserecognitionlinepb.ListExpenseRecognitionLinesResponse, error)
	ReverseExpenseRecognition    func(ctx context.Context, id, reason string) error

	attachment.AttachmentOps
}

// PageData holds template data for the expense_recognition detail page.
type PageData struct {
	types.PageData
	ContentTemplate string

	Recognition   map[string]any
	StatusVariant string

	TabItems  []pyeza.TabItem
	ActiveTab string

	LineItemTable  *types.TableConfig
	LineItemAddURL string

	AttachmentTable *types.TableConfig

	ReverseURL string
}

const (
	tabInfo        = "info"
	tabLines       = "lines"
	tabSource      = "source"
	tabActivity    = "activity"
	tabAttachments = "attachments"
)

// NewView creates the expense_recognition detail page view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Redirect(deps.Routes.ListURL)
		}

		resp, err := deps.ReadExpenseRecognition(ctx, &expenserecognitionpb.ReadExpenseRecognitionRequest{
			Data: &expenserecognitionpb.ExpenseRecognition{Id: id},
		})
		if err != nil {
			log.Printf("ReadExpenseRecognition %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load expense recognition: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("expense recognition not found"))
		}
		rec := data[0]

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = tabInfo
		}

		l := deps.Labels
		recMap := recognitionToMap(rec)

		tabItems := []pyeza.TabItem{
			{Key: tabInfo, Label: l.Tabs.Info},
			{Key: tabLines, Label: l.Tabs.Lines},
			{Key: tabSource, Label: l.Tabs.Source},
			{Key: tabActivity, Label: l.Tabs.Activity},
			{Key: tabAttachments, Label: l.Detail.TabAttachments, Icon: "icon-paperclip"},
		}

		pd := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          rec.GetName(),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    rec.GetName(),
				HeaderSubtitle: l.Detail.Title,
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "expense-recognition-detail-content",
			Recognition:     recMap,
			StatusVariant:   recognitionStatusVariant(rec.GetStatus().String()),
			TabItems:        tabItems,
			ActiveTab:       activeTab,
			ReverseURL:      buildActionURL(deps.Routes.ReverseURL, id),
		}

		if activeTab == tabLines && deps.ListExpenseRecognitionLines != nil {
			pd.LineItemTable = buildLineItemTable(ctx, deps, id, l)
			pd.LineItemAddURL = route.ResolveURL(deps.Routes.LineAddURL, "id", id)
		}

		if activeTab == tabAttachments && deps.ListAttachments != nil {
			cfg := attachmentConfig(deps)
			resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
			if err != nil {
				log.Printf("Failed to list attachments for expense_recognition %s: %v", id, err)
			}
			var items []*attachmentpb.Attachment
			if resp != nil {
				items = resp.GetData()
			}
			pd.AttachmentTable = attachment.BuildTable(items, cfg, id)
		}

		return view.OK("expense-recognition-detail", pd)
	})
}

// NewTabAction handles HTMX tab switch.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if id == "" || tab == "" {
			return view.Error(fmt.Errorf("missing id or tab"))
		}

		resp, err := deps.ReadExpenseRecognition(ctx, &expenserecognitionpb.ReadExpenseRecognitionRequest{
			Data: &expenserecognitionpb.ExpenseRecognition{Id: id},
		})
		if err != nil {
			return view.Error(fmt.Errorf("failed to load expense recognition: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("expense recognition not found"))
		}
		rec := data[0]

		l := deps.Labels
		recMap := recognitionToMap(rec)

		tabItems := []pyeza.TabItem{
			{Key: tabInfo, Label: l.Tabs.Info},
			{Key: tabLines, Label: l.Tabs.Lines},
			{Key: tabSource, Label: l.Tabs.Source},
			{Key: tabActivity, Label: l.Tabs.Activity},
			{Key: tabAttachments, Label: l.Detail.TabAttachments, Icon: "icon-paperclip"},
		}

		pd := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "expense-recognition-detail-content",
			Recognition:     recMap,
			StatusVariant:   recognitionStatusVariant(rec.GetStatus().String()),
			TabItems:        tabItems,
			ActiveTab:       tab,
			ReverseURL:      buildActionURL(deps.Routes.ReverseURL, id),
		}

		if tab == tabLines && deps.ListExpenseRecognitionLines != nil {
			pd.LineItemTable = buildLineItemTable(ctx, deps, id, l)
			pd.LineItemAddURL = route.ResolveURL(deps.Routes.LineAddURL, "id", id)
		}

		if tab == tabAttachments && deps.ListAttachments != nil {
			cfg := attachmentConfig(deps)
			resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
			if err != nil {
				log.Printf("Failed to list attachments for expense_recognition %s: %v", id, err)
			}
			var items []*attachmentpb.Attachment
			if resp != nil {
				items = resp.GetData()
			}
			pd.AttachmentTable = attachment.BuildTable(items, cfg, id)
		}

		templateName := "expense-recognition-tab-content"
		if tab == tabAttachments {
			templateName = "attachment-tab"
		}
		return view.OK(templateName, pd)
	})
}

// NewReverseAction handles POST .../reverse/{id}.
func NewReverseAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		reason := ""
		if err := viewCtx.Request.ParseForm(); err == nil {
			reason = viewCtx.Request.FormValue("reason")
		}
		if deps.ReverseExpenseRecognition != nil {
			if err := deps.ReverseExpenseRecognition(ctx, id, reason); err != nil {
				log.Printf("ReverseExpenseRecognition %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
		}
		detailURL := route.ResolveURL(deps.Routes.DetailURL, "id", id)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Redirect": detailURL,
			},
		}
	})
}

// --- helpers -----------------------------------------------------------------

func recognitionToMap(r *expenserecognitionpb.ExpenseRecognition) map[string]any {
	currency := r.GetCurrency()
	periodStart := ""
	if r.PeriodStart != nil {
		periodStart = r.GetPeriodStart().AsTime().Format("2006-01-02")
	}
	periodEnd := ""
	if r.PeriodEnd != nil {
		periodEnd = r.GetPeriodEnd().AsTime().Format("2006-01-02")
	}
	recognitionDate := ""
	if r.RecognitionDate != nil {
		recognitionDate = r.GetRecognitionDate().AsTime().Format("2006-01-02")
	}
	return map[string]any{
		"id":                         r.GetId(),
		"name":                       r.GetName(),
		"description":                r.GetDescription(),
		"status":                     r.GetStatus().String(),
		"currency":                   currency,
		"total_amount":               types.MoneyCell(float64(r.GetTotalAmount()), currency, true),
		"recognition_date":           recognitionDate,
		"period_start":               periodStart,
		"period_end":                 periodEnd,
		"cycle_date":                 r.GetCycleDate(),
		"supplier_contract_id":       r.GetSupplierContractId(),
		"expenditure_id":             r.GetExpenditureId(),
		"deferred_expense_id":        r.GetDeferredExpenseId(),
		"accrued_expense_id":         r.GetAccruedExpenseId(),
		"idempotency_key":            r.GetIdempotencyKey(),
		"reversal_of_recognition_id": r.GetReversalOfRecognitionId(),
		"notes":                      r.GetNotes(),
	}
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

func buildActionURL(template, id string) string {
	if template == "" {
		return ""
	}
	return route.ResolveURL(template, "id", id)
}

func buildLineItemTable(ctx context.Context, deps *DetailViewDeps, recognitionID string, l centymo.ExpenseRecognitionLabels) *types.TableConfig {
	rid := recognitionID
	resp, err := deps.ListExpenseRecognitionLines(ctx, &expenserecognitionlinepb.ListExpenseRecognitionLinesRequest{
		ExpenseRecognitionId: &rid,
	})
	if err != nil {
		log.Printf("ListExpenseRecognitionLines for %s: %v", recognitionID, err)
		return nil
	}

	columns := []types.TableColumn{
		{Key: "description", Label: l.Lines.Description},
		{Key: "quantity", Label: l.Lines.Quantity, Align: "right", WidthClass: "col-xl"},
		{Key: "unit_amount", Label: l.Lines.UnitAmount, Align: "right", WidthClass: "col-3xl"},
		{Key: "amount", Label: l.Lines.Amount, Align: "right", WidthClass: "col-3xl"},
	}

	rows := []types.TableRow{}
	for _, line := range resp.GetData() {
		rows = append(rows, types.TableRow{
			ID: line.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: line.GetDescription()},
				{Type: "number", Value: fmt.Sprintf("%.2f", line.GetQuantity())},
				types.MoneyCell(float64(line.GetUnitAmount()), line.GetCurrency(), true),
				types.MoneyCell(float64(line.GetAmount()), line.GetCurrency(), true),
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:          "expense-recognition-lines-table",
		Columns:     columns,
		Rows:        rows,
		ShowActions: true,
		EmptyState: types.TableEmptyState{
			Title:   l.Lines.EmptyTitle,
			Message: l.Lines.EmptyMessage,
		},
	}
}
