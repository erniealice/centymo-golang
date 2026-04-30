package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/view"

	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
)

// RecognizeFromContractFunc is the espyna use-case function pointer threaded
// through ModuleDeps.
type RecognizeFromContractFunc func(ctx context.Context, req *expenserecognitionpb.RecognizeFromContractRequest) (*expenserecognitionpb.RecognizeFromContractResponse, error)

// NewRecognizeFromContractAction handles
// POST /action/expense-recognition/recognize-from-contract.
//
// Form fields:
//   - supplier_contract_id (required) — UUID of the supplier contract
//   - cycle_date (required) — YYYY-MM-DD
//   - amount (optional) — centavos; resolved from schedule if omitted
//   - idempotency_key (optional)
//
// Returns 422 on missing required fields or use-case error; 200 on success.
func NewRecognizeFromContractAction(fn RecognizeFromContractFunc) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		if fn == nil {
			return centymo.HTMXError("recognize-from-contract handler not wired")
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("invalid form data")
		}
		contractID := viewCtx.Request.FormValue("supplier_contract_id")
		cycleDate := viewCtx.Request.FormValue("cycle_date")
		if contractID == "" {
			return centymo.HTMXError("supplier_contract_id is required")
		}
		if cycleDate == "" {
			return centymo.HTMXError("cycle_date is required")
		}
		req := &expenserecognitionpb.RecognizeFromContractRequest{
			SupplierContractId: contractID,
			CycleDate:          cycleDate,
		}
		if amtStr := viewCtx.Request.FormValue("amount"); amtStr != "" {
			if amt, err := strconv.ParseInt(amtStr, 10, 64); err == nil {
				req.Amount = &amt
			}
		}
		if ik := viewCtx.Request.FormValue("idempotency_key"); ik != "" {
			req.IdempotencyKey = &ik
		}
		if _, err := fn(ctx, req); err != nil {
			log.Printf("RecognizeFromContract %s @ %s: %v", contractID, cycleDate, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("expense-recognitions-table")
	})
}
