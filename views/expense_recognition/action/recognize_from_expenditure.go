package action

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
)

// RecognizeFromExpenditureFunc is the espyna use-case function pointer threaded
// through ModuleDeps.
type RecognizeFromExpenditureFunc func(ctx context.Context, req *expenserecognitionpb.RecognizeFromExpenditureRequest) (*expenserecognitionpb.RecognizeFromExpenditureResponse, error)

// NewRecognizeFromExpenditureAction handles
// POST /action/expense-recognition/recognize-from-expenditure.
//
// Form fields:
//   - expenditure_id (required) — UUID of the source expenditure
//   - recognition_period (optional) — YYYY-MM (or YYYY-MM-DD)
//   - idempotency_key (optional) — overrides default derivation
//
// Returns 422 (HTMXError) on missing expenditure_id or use-case error;
// 200 with HX-Trigger on success.
func NewRecognizeFromExpenditureAction(fn RecognizeFromExpenditureFunc) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense_recognition", "create") {
			return view.HTMXError("Missing permission: expense_recognition:create")
		}
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		if fn == nil {
			return view.HTMXError("recognize-from-expenditure handler not wired")
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError("invalid form data")
		}
		expenditureID := viewCtx.Request.FormValue("expenditure_id")
		if expenditureID == "" {
			return view.HTMXError("expenditure_id is required")
		}
		req := &expenserecognitionpb.RecognizeFromExpenditureRequest{
			ExpenditureId: expenditureID,
		}
		if rp := viewCtx.Request.FormValue("recognition_period"); rp != "" {
			req.RecognitionPeriod = &rp
		}
		if ik := viewCtx.Request.FormValue("idempotency_key"); ik != "" {
			req.IdempotencyKey = &ik
		}
		if _, err := fn(ctx, req); err != nil {
			log.Printf("RecognizeFromExpenditure %s: %v", expenditureID, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("expense-recognitions-table")
	})
}
