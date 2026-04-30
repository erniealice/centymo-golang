package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"
)

// DetailViewDeps holds dependencies for the accrued_expense detail page.
type DetailViewDeps struct {
	Routes       centymo.AccruedExpenseRoutes
	Labels       centymo.AccruedExpenseLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadAccruedExpense           func(ctx context.Context, req *accruedexpensepb.ReadAccruedExpenseRequest) (*accruedexpensepb.ReadAccruedExpenseResponse, error)
	ListAccruedExpenseSettlements func(ctx context.Context, req *accruedexpensepb.ListAccruedExpenseSettlementsRequest) (*accruedexpensepb.ListAccruedExpenseSettlementsResponse, error)

	// Workflow — provided as closures from block.go
	SettleAccrual  func(ctx context.Context, req *accruedexpensepb.SettleAccrualRequest) error
	ReverseAccrual func(ctx context.Context, id, reason string) error
}

// PageData holds template data for the accrued_expense detail page.
type PageData struct {
	types.PageData
	ContentTemplate string

	Accrual       map[string]any
	StatusVariant string

	TabItems  []pyeza.TabItem
	ActiveTab string

	SettlementTable     *types.TableConfig
	SettlementAddURL    string

	SettleURL  string
	ReverseURL string
	EditURL    string
}

const (
	tabInfo        = "info"
	tabSettlements = "settlements"
	tabSource      = "source"
	tabActivity    = "activity"
)

// NewView creates the accrued_expense detail page view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Redirect(deps.Routes.ListURL)
		}

		resp, err := deps.ReadAccruedExpense(ctx, &accruedexpensepb.ReadAccruedExpenseRequest{
			Data: &accruedexpensepb.AccruedExpense{Id: id},
		})
		if err != nil {
			log.Printf("ReadAccruedExpense %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load accrual: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("accrual not found"))
		}
		acc := data[0]

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = tabInfo
		}

		l := deps.Labels
		accMap := accrualToMap(acc)

		tabItems := []pyeza.TabItem{
			{Key: tabInfo, Label: l.Tabs.Info},
			{Key: tabSettlements, Label: l.Tabs.Settlements},
			{Key: tabSource, Label: l.Tabs.Source},
			{Key: tabActivity, Label: l.Tabs.Activity},
		}

		pd := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          acc.GetName(),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    acc.GetName(),
				HeaderSubtitle: l.Detail.Title,
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "accrued-expense-detail-content",
			Accrual:         accMap,
			StatusVariant:   accrualStatusVariant(acc.GetStatus().String()),
			TabItems:        tabItems,
			ActiveTab:       activeTab,
			SettleURL:       buildActionURL(deps.Routes.SettleURL, id),
			ReverseURL:      buildActionURL(deps.Routes.ReverseURL, id),
			EditURL:         buildActionURL(deps.Routes.EditURL, id),
		}

		if activeTab == tabSettlements && deps.ListAccruedExpenseSettlements != nil {
			pd.SettlementTable = buildSettlementTable(ctx, deps, id, l)
			pd.SettlementAddURL = route.ResolveURL(deps.Routes.SettlementAddURL, "id", id)
		}

		return view.OK("accrued-expense-detail", pd)
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

		resp, err := deps.ReadAccruedExpense(ctx, &accruedexpensepb.ReadAccruedExpenseRequest{
			Data: &accruedexpensepb.AccruedExpense{Id: id},
		})
		if err != nil {
			return view.Error(fmt.Errorf("failed to load accrual: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("accrual not found"))
		}
		acc := data[0]

		l := deps.Labels
		accMap := accrualToMap(acc)

		tabItems := []pyeza.TabItem{
			{Key: tabInfo, Label: l.Tabs.Info},
			{Key: tabSettlements, Label: l.Tabs.Settlements},
			{Key: tabSource, Label: l.Tabs.Source},
			{Key: tabActivity, Label: l.Tabs.Activity},
		}

		pd := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "accrued-expense-detail-content",
			Accrual:         accMap,
			StatusVariant:   accrualStatusVariant(acc.GetStatus().String()),
			TabItems:        tabItems,
			ActiveTab:       tab,
			SettleURL:       buildActionURL(deps.Routes.SettleURL, id),
			ReverseURL:      buildActionURL(deps.Routes.ReverseURL, id),
			EditURL:         buildActionURL(deps.Routes.EditURL, id),
		}

		if tab == tabSettlements && deps.ListAccruedExpenseSettlements != nil {
			pd.SettlementTable = buildSettlementTable(ctx, deps, id, l)
			pd.SettlementAddURL = route.ResolveURL(deps.Routes.SettlementAddURL, "id", id)
		}

		return view.OK("accrued-expense-tab-content", pd)
	})
}

// NewSettleAction handles POST .../settle/{id}.
// Operator-driven settlement — accepts expenditure_id + amount_settled + currency
// from the form body and forwards to the SettleAccrual use case.
func NewSettleAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("invalid form data")
		}
		r := viewCtx.Request

		expenditureID := r.FormValue("expenditure_id")
		amountStr := r.FormValue("amount_settled")
		currency := r.FormValue("currency")

		amount, _ := strconv.ParseFloat(amountStr, 64)
		req := &accruedexpensepb.SettleAccrualRequest{
			AccruedExpenseId: id,
			ExpenditureId:    expenditureID,
			AmountSettled:    int64(amount * 100),
			Currency:         currency,
		}
		if line := r.FormValue("expenditure_line_item_id"); line != "" {
			req.ExpenditureLineItemId = &line
		}
		if fxStr := r.FormValue("fx_rate"); fxStr != "" {
			if fx, err := strconv.ParseFloat(fxStr, 64); err == nil {
				req.FxRate = &fx
			}
		}

		if deps.SettleAccrual != nil {
			if err := deps.SettleAccrual(ctx, req); err != nil {
				log.Printf("SettleAccrual %s: %v", id, err)
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
		if deps.ReverseAccrual != nil {
			if err := deps.ReverseAccrual(ctx, id, reason); err != nil {
				log.Printf("ReverseAccrual %s: %v", id, err)
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

func accrualToMap(a *accruedexpensepb.AccruedExpense) map[string]any {
	currency := a.GetCurrency()
	periodStart := ""
	if a.PeriodStart != nil {
		periodStart = a.GetPeriodStart().AsTime().Format("2006-01-02")
	}
	periodEnd := ""
	if a.PeriodEnd != nil {
		periodEnd = a.GetPeriodEnd().AsTime().Format("2006-01-02")
	}
	recognitionDate := ""
	if a.RecognitionDate != nil {
		recognitionDate = a.GetRecognitionDate().AsTime().Format("2006-01-02")
	}
	return map[string]any{
		"id":                   a.GetId(),
		"name":                 a.GetName(),
		"description":          a.GetDescription(),
		"status":               a.GetStatus().String(),
		"currency":             currency,
		"accrued_amount":       types.MoneyCell(float64(a.GetAccruedAmount()), currency, true),
		"settled_amount":       types.MoneyCell(float64(a.GetSettledAmount()), currency, true),
		"remaining_amount":     types.MoneyCell(float64(a.GetRemainingAmount()), currency, true),
		"recognition_date":     recognitionDate,
		"period_start":         periodStart,
		"period_end":           periodEnd,
		"cycle_date":           a.GetCycleDate(),
		"supplier_contract_id": a.GetSupplierContractId(),
		"supplier_id":          a.GetSupplierId(),
		"expense_account_id":   a.GetExpenseAccountId(),
		"accrual_account_id":   a.GetAccrualAccountId(),
		"notes":                a.GetNotes(),
	}
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

func buildActionURL(template, id string) string {
	if template == "" {
		return ""
	}
	return route.ResolveURL(template, "id", id)
}

func buildSettlementTable(ctx context.Context, deps *DetailViewDeps, accrualID string, l centymo.AccruedExpenseLabels) *types.TableConfig {
	aid := accrualID
	resp, err := deps.ListAccruedExpenseSettlements(ctx, &accruedexpensepb.ListAccruedExpenseSettlementsRequest{
		AccruedExpenseId: &aid,
	})
	if err != nil {
		log.Printf("ListAccruedExpenseSettlements for %s: %v", accrualID, err)
		return nil
	}

	columns := []types.TableColumn{
		{Key: "expenditure", Label: l.Settlements.Expenditure},
		{Key: "amount_settled", Label: l.Settlements.AmountSettled, Align: "right", WidthClass: "col-3xl"},
		{Key: "currency", Label: l.Settlements.Currency, WidthClass: "col-xl"},
		{Key: "fx_rate", Label: l.Settlements.FxRate, Align: "right", WidthClass: "col-xl"},
		{Key: "settled_at", Label: l.Settlements.SettledAt, WidthClass: "col-3xl"},
	}

	rows := []types.TableRow{}
	for _, s := range resp.GetData() {
		settledAt := ""
		if s.SettledAt != nil {
			settledAt = s.GetSettledAt().AsTime().Format("2006-01-02")
		}
		fxRate := ""
		if s.FxRate != nil {
			fxRate = fmt.Sprintf("%.4f", s.GetFxRate())
		}
		rows = append(rows, types.TableRow{
			ID: s.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: s.GetExpenditureId()},
				types.MoneyCell(float64(s.GetAmountSettled()), s.GetCurrency(), true),
				{Type: "text", Value: s.GetCurrency()},
				{Type: "text", Value: fxRate},
				{Type: "text", Value: settledAt},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:          "accrued-expense-settlements-table",
		Columns:     columns,
		Rows:        rows,
		ShowActions: true,
		EmptyState: types.TableEmptyState{
			Title:   l.Settlements.EmptyTitle,
			Message: l.Settlements.EmptyMessage,
		},
	}
}
