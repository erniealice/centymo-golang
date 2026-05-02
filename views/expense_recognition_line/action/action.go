// Package action implements inline CRUD action handlers for
// ExpenseRecognitionLine. Lines are rendered inside the parent
// ExpenseRecognition detail page's Lines tab; there is no standalone list.
package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	expenserecognitionlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition_line"

	"github.com/erniealice/centymo-golang/views/expense_recognition_line/form"
)

// Deps holds dependencies for line action handlers.
type Deps struct {
	Routes       centymo.ExpenseRecognitionRoutes
	Labels       centymo.ExpenseRecognitionLabels
	CommonLabels pyeza.CommonLabels

	CreateExpenseRecognitionLine func(ctx context.Context, req *expenserecognitionlinepb.CreateExpenseRecognitionLineRequest) (*expenserecognitionlinepb.CreateExpenseRecognitionLineResponse, error)
	ReadExpenseRecognitionLine   func(ctx context.Context, req *expenserecognitionlinepb.ReadExpenseRecognitionLineRequest) (*expenserecognitionlinepb.ReadExpenseRecognitionLineResponse, error)
	UpdateExpenseRecognitionLine func(ctx context.Context, req *expenserecognitionlinepb.UpdateExpenseRecognitionLineRequest) (*expenserecognitionlinepb.UpdateExpenseRecognitionLineResponse, error)
	DeleteExpenseRecognitionLine func(ctx context.Context, req *expenserecognitionlinepb.DeleteExpenseRecognitionLineRequest) (*expenserecognitionlinepb.DeleteExpenseRecognitionLineResponse, error)
}

// NewAddAction handles GET (form) + POST (create) on the line add URL.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		recognitionID := viewCtx.Request.PathValue("id")
		if recognitionID == "" {
			return centymo.HTMXError("missing recognition id")
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("expense-recognition-line-drawer-form", &form.Data{
				FormAction:           route.ResolveURL(deps.Routes.LineAddURL, "id", recognitionID),
				ExpenseRecognitionID: recognitionID,
				Quantity:             "1",
				Labels:               deps.Labels.Lines,
				CommonLabels:         deps.CommonLabels,
			})
		}

		// POST
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("invalid form data")
		}
		r := viewCtx.Request

		quantityF, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
		if quantityF == 0 {
			quantityF = 1
		}
		unitAmountF, _ := strconv.ParseFloat(r.FormValue("unit_amount"), 64)
		amountF := quantityF * unitAmountF

		_, err := deps.CreateExpenseRecognitionLine(ctx, &expenserecognitionlinepb.CreateExpenseRecognitionLineRequest{
			Data: &expenserecognitionlinepb.ExpenseRecognitionLine{
				ExpenseRecognitionId: recognitionID,
				Description:          r.FormValue("description"),
				Quantity:             quantityF,
				UnitAmount:           int64(math.Round(unitAmountF * 100)),
				Amount:               int64(math.Round(amountF * 100)),
				Currency:             r.FormValue("currency"),
			},
		})
		if err != nil {
			log.Printf("CreateExpenseRecognitionLine: %v", err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("expense-recognition-lines-table")
	})
}

// NewEditAction handles GET (form) + POST (update) on the line edit URL.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		recognitionID := viewCtx.Request.PathValue("id")
		lineID := viewCtx.Request.PathValue("lid")
		if recognitionID == "" || lineID == "" {
			return centymo.HTMXError("missing id or lid")
		}

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadExpenseRecognitionLine(ctx, &expenserecognitionlinepb.ReadExpenseRecognitionLineRequest{
				Data: &expenserecognitionlinepb.ExpenseRecognitionLine{Id: lineID},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				return centymo.HTMXError("recognition line not found")
			}
			line := readResp.GetData()[0]
			return view.OK("expense-recognition-line-drawer-form", &form.Data{
				FormAction:           route.ResolveURL(deps.Routes.LineEditURL, "id", recognitionID, "lid", lineID),
				IsEdit:               true,
				ID:                   lineID,
				ExpenseRecognitionID: recognitionID,
				Description:          line.GetDescription(),
				Quantity:             fmt.Sprintf("%.2f", line.GetQuantity()),
				UnitAmount:           fmt.Sprintf("%.2f", float64(line.GetUnitAmount())/100.0),
				Amount:               fmt.Sprintf("%.2f", float64(line.GetAmount())/100.0),
				Currency:             line.GetCurrency(),
				Labels:               deps.Labels.Lines,
				CommonLabels:         deps.CommonLabels,
			})
		}

		// POST
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("invalid form data")
		}
		r := viewCtx.Request

		quantityF, _ := strconv.ParseFloat(r.FormValue("quantity"), 64)
		if quantityF == 0 {
			quantityF = 1
		}
		unitAmountF, _ := strconv.ParseFloat(r.FormValue("unit_amount"), 64)
		amountF := quantityF * unitAmountF

		_, err := deps.UpdateExpenseRecognitionLine(ctx, &expenserecognitionlinepb.UpdateExpenseRecognitionLineRequest{
			Data: &expenserecognitionlinepb.ExpenseRecognitionLine{
				Id:          lineID,
				Description: r.FormValue("description"),
				Quantity:    quantityF,
				UnitAmount:  int64(math.Round(unitAmountF * 100)),
				Amount:      int64(math.Round(amountF * 100)),
				Currency:    r.FormValue("currency"),
			},
		})
		if err != nil {
			log.Printf("UpdateExpenseRecognitionLine %s: %v", lineID, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("expense-recognition-lines-table")
	})
}

// NewDeleteAction handles POST .../lines/delete.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return centymo.HTMXError("method not allowed")
		}
		lineID := viewCtx.Request.URL.Query().Get("lid")
		if lineID == "" {
			_ = viewCtx.Request.ParseForm()
			lineID = viewCtx.Request.FormValue("lid")
			if lineID == "" {
				lineID = viewCtx.Request.FormValue("id")
			}
		}
		if lineID == "" {
			return centymo.HTMXError("missing line id")
		}

		_, err := deps.DeleteExpenseRecognitionLine(ctx, &expenserecognitionlinepb.DeleteExpenseRecognitionLineRequest{
			Data: &expenserecognitionlinepb.ExpenseRecognitionLine{Id: lineID},
		})
		if err != nil {
			log.Printf("DeleteExpenseRecognitionLine %s: %v", lineID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("expense-recognition-lines-table")
	})
}
