package action

import (
	"context"
	"fmt"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/view"

	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
)

// Deps holds dependencies for the expense_recognition action handlers.
//
// Note: there is intentionally NO Add/Edit drawer-form action — recognitions
// are created BY use case (RecognizeFromExpenditure / RecognizeFromContract),
// not by user input. The only operator-driven action is Reverse (handled by
// the detail package).
type Deps struct {
	Routes                  centymo.ExpenseRecognitionRoutes
	Labels                  centymo.ExpenseRecognitionLabels
	DeleteExpenseRecognition func(ctx context.Context, req *expenserecognitionpb.DeleteExpenseRecognitionRequest) (*expenserecognitionpb.DeleteExpenseRecognitionResponse, error)
}

// NewDeleteAction handles POST /action/expense-recognition/delete (Draft only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.FormValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		if deps.DeleteExpenseRecognition == nil {
			return centymo.HTMXError("delete handler not wired")
		}
		_, err := deps.DeleteExpenseRecognition(ctx, &expenserecognitionpb.DeleteExpenseRecognitionRequest{
			Data: &expenserecognitionpb.ExpenseRecognition{Id: id},
		})
		if err != nil {
			log.Printf("DeleteExpenseRecognition %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to delete recognition: %w", err))
		}
		return centymo.HTMXSuccess("expense-recognitions-table")
	})
}
