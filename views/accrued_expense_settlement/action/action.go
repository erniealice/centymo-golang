// Package action implements inline CRUD action handlers for
// AccruedExpenseSettlement. Settlements render inside the parent
// AccruedExpense detail page's Settlements tab; there is no standalone list.
//
// Single-write boundary: settled_amount / remaining_amount on the parent
// AccruedExpense are updated EXCLUSIVELY by the SettleAccrual / ReverseAccrual
// use cases. The plain CRUD handlers below should only be used for
// adjustments where direct manipulation is intentional (e.g., FX revaluation
// historical correction); they bypass status recompute and expect the caller
// to follow up with a balance refresh.
package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/erniealice/centymo-golang/views/accrued_expense_settlement/form"
)

// Deps holds dependencies for settlement action handlers.
type Deps struct {
	Routes       centymo.AccruedExpenseRoutes
	Labels       centymo.AccruedExpenseLabels
	CommonLabels pyeza.CommonLabels

	CreateAccruedExpenseSettlement func(ctx context.Context, req *accruedexpensepb.CreateAccruedExpenseSettlementRequest) (*accruedexpensepb.CreateAccruedExpenseSettlementResponse, error)
	ReadAccruedExpenseSettlement   func(ctx context.Context, req *accruedexpensepb.ReadAccruedExpenseSettlementRequest) (*accruedexpensepb.ReadAccruedExpenseSettlementResponse, error)
	UpdateAccruedExpenseSettlement func(ctx context.Context, req *accruedexpensepb.UpdateAccruedExpenseSettlementRequest) (*accruedexpensepb.UpdateAccruedExpenseSettlementResponse, error)
	DeleteAccruedExpenseSettlement func(ctx context.Context, req *accruedexpensepb.DeleteAccruedExpenseSettlementRequest) (*accruedexpensepb.DeleteAccruedExpenseSettlementResponse, error)

	// Dropdown
	ListExpenditures func(ctx context.Context, req *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error)
}

// NewAddAction handles GET (form) + POST (create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		accrualID := viewCtx.Request.PathValue("id")
		if accrualID == "" {
			return centymo.HTMXError("missing accrual id")
		}

		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyFormData(ctx, deps, deps.Labels.Settlements)
			fd.FormAction = route.ResolveURL(deps.Routes.SettlementAddURL, "id", accrualID)
			fd.AccruedExpenseID = accrualID
			fd.Currency = "PHP"
			return view.OK("accrued-expense-settlement-drawer-form", fd)
		}

		// POST
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("invalid form data")
		}
		r := viewCtx.Request

		amountF, _ := strconv.ParseFloat(r.FormValue("amount_settled"), 64)
		data := &accruedexpensepb.AccruedExpenseSettlement{
			AccruedExpenseId: accrualID,
			ExpenditureId:    r.FormValue("expenditure_id"),
			AmountSettled:    int64(math.Round(amountF * 100)),
			Currency:         r.FormValue("currency"),
			SettledAt:        timestamppb.New(time.Now()),
		}
		if line := r.FormValue("expenditure_line_item_id"); line != "" {
			data.ExpenditureLineItemId = &line
		}
		if fxStr := r.FormValue("fx_rate"); fxStr != "" {
			if fx, err := strconv.ParseFloat(fxStr, 64); err == nil {
				data.FxRate = &fx
			}
		}

		_, err := deps.CreateAccruedExpenseSettlement(ctx, &accruedexpensepb.CreateAccruedExpenseSettlementRequest{Data: data})
		if err != nil {
			log.Printf("CreateAccruedExpenseSettlement: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("accrued-expense-settlements-table")
	})
}

// NewEditAction handles GET (form) + POST (update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		accrualID := viewCtx.Request.PathValue("id")
		settlementID := viewCtx.Request.PathValue("sid")
		if accrualID == "" || settlementID == "" {
			return centymo.HTMXError("missing id or sid")
		}

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadAccruedExpenseSettlement(ctx, &accruedexpensepb.ReadAccruedExpenseSettlementRequest{
				Data: &accruedexpensepb.AccruedExpenseSettlement{Id: settlementID},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				return centymo.HTMXError("settlement not found")
			}
			s := readResp.GetData()[0]

			fd := buildEmptyFormData(ctx, deps, deps.Labels.Settlements)
			fd.FormAction = route.ResolveURL(deps.Routes.SettlementEditURL, "id", accrualID, "sid", settlementID)
			fd.IsEdit = true
			fd.ID = settlementID
			fd.AccruedExpenseID = accrualID
			fd.ExpenditureID = s.GetExpenditureId()
			fd.AmountSettled = fmt.Sprintf("%.2f", float64(s.GetAmountSettled())/100.0)
			fd.Currency = s.GetCurrency()
			if s.FxRate != nil {
				fd.FxRate = fmt.Sprintf("%.4f", s.GetFxRate())
			}
			return view.OK("accrued-expense-settlement-drawer-form", fd)
		}

		// POST
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("invalid form data")
		}
		r := viewCtx.Request
		amountF, _ := strconv.ParseFloat(r.FormValue("amount_settled"), 64)

		data := &accruedexpensepb.AccruedExpenseSettlement{
			Id:               settlementID,
			AccruedExpenseId: accrualID,
			ExpenditureId:    r.FormValue("expenditure_id"),
			AmountSettled:    int64(math.Round(amountF * 100)),
			Currency:         r.FormValue("currency"),
		}
		if fxStr := r.FormValue("fx_rate"); fxStr != "" {
			if fx, err := strconv.ParseFloat(fxStr, 64); err == nil {
				data.FxRate = &fx
			}
		}

		_, err := deps.UpdateAccruedExpenseSettlement(ctx, &accruedexpensepb.UpdateAccruedExpenseSettlementRequest{Data: data})
		if err != nil {
			log.Printf("UpdateAccruedExpenseSettlement %s: %v", settlementID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("accrued-expense-settlements-table")
	})
}

// NewDeleteAction handles POST .../settlements/delete.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return centymo.HTMXError("method not allowed")
		}
		settlementID := viewCtx.Request.URL.Query().Get("sid")
		if settlementID == "" {
			_ = viewCtx.Request.ParseForm()
			settlementID = viewCtx.Request.FormValue("sid")
			if settlementID == "" {
				settlementID = viewCtx.Request.FormValue("id")
			}
		}
		if settlementID == "" {
			return centymo.HTMXError("missing settlement id")
		}
		_, err := deps.DeleteAccruedExpenseSettlement(ctx, &accruedexpensepb.DeleteAccruedExpenseSettlementRequest{
			Data: &accruedexpensepb.AccruedExpenseSettlement{Id: settlementID},
		})
		if err != nil {
			log.Printf("DeleteAccruedExpenseSettlement %s: %v", settlementID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("accrued-expense-settlements-table")
	})
}

// --- form helpers ------------------------------------------------------------

func buildEmptyFormData(ctx context.Context, deps *Deps, l centymo.AccruedExpenseSettlementLabels) *form.Data {
	fd := &form.Data{
		Labels:       l,
		CommonLabels: deps.CommonLabels,
	}
	if deps.ListExpenditures != nil {
		resp, err := deps.ListExpenditures(ctx, &expenditurepb.ListExpendituresRequest{})
		if err == nil {
			for _, e := range resp.GetData() {
				label := e.GetReferenceNumber()
				if label == "" {
					label = e.GetName()
				}
				fd.Expenditures = append(fd.Expenditures, types.SelectOption{
					Value: e.GetId(),
					Label: label,
				})
			}
		}
	}
	return fd
}
