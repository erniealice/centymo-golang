package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/view"

	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"
)

// AccrueFromContractFunc is the espyna use-case function pointer threaded
// through ModuleDeps.
type AccrueFromContractFunc func(ctx context.Context, req *accruedexpensepb.AccrueFromContractRequest) (*accruedexpensepb.AccrueFromContractResponse, error)

// NewAccrueFromContractAction handles
// POST /action/accrued-expense/accrue-from-contract.
//
// Form fields:
//   - supplier_contract_id (required)
//   - cycle_date (required) — YYYY-MM-DD
//   - accrued_amount (optional) — centavos; resolved from schedule if omitted
//
// Returns 422 on missing required fields or use-case error; 200 on success.
func NewAccrueFromContractAction(fn AccrueFromContractFunc) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		if fn == nil {
			return centymo.HTMXError("accrue-from-contract handler not wired")
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
		req := &accruedexpensepb.AccrueFromContractRequest{
			SupplierContractId: contractID,
			CycleDate:          cycleDate,
		}
		if amtStr := viewCtx.Request.FormValue("accrued_amount"); amtStr != "" {
			if amt, err := strconv.ParseInt(amtStr, 10, 64); err == nil {
				req.AccruedAmount = &amt
			}
		}
		if _, err := fn(ctx, req); err != nil {
			log.Printf("AccrueFromContract %s @ %s: %v", contractID, cycleDate, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("accrued-expenses-table")
	})
}
