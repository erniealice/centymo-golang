package action

import (
	"context"
	"fmt"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/view"

	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
	suppliersubscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_subscription"
)

// RecognizeExpenseDeps holds deps for the Recognize Expense CTA on the
// supplier_subscription detail page.
type RecognizeExpenseDeps struct {
	Routes centymo.SupplierSubscriptionRoutes
	Labels centymo.SupplierSubscriptionLabels

	// GetSupplierSubscriptionItemPageData is used to resolve the supplier_subscription
	// record (to find the most-recent linked expenditure_id for recognition).
	// Optional — when nil the handler falls back to expenditure_id from the form.
	GetSupplierSubscriptionItemPageData func(ctx context.Context, req *suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataRequest) (*suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataResponse, error)

	// RecognizeFromExpenditure is the espyna use case function.
	// Required — 422 returned when nil.
	RecognizeFromExpenditure func(ctx context.Context, req *expenserecognitionpb.RecognizeFromExpenditureRequest) (*expenserecognitionpb.RecognizeFromExpenditureResponse, error)
}

// NewRecognizeExpenseAction handles POST /action/supplier-subscription/recognize-expense/{id}.
//
// This is the buying-side mirror of the "Recognize Revenue" CTA on the
// subscription detail page. It calls RecognizeFromExpenditure with:
//   - expenditure_id from the form body (required)
//   - recognition_period from the form body (optional)
//   - supplier_subscription_id is threaded by the use case (Wave 2 / Agent E)
//     from the Expenditure's supplier_subscription_id FK.
//
// On success returns HX-Trigger("expense-recognitions-table") so the
// Linked Recognitions tab refreshes inline (mirrors selling-side pattern).
func NewRecognizeExpenseAction(deps *RecognizeExpenseDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense_recognition", "create") {
			return centymo.HTMXError("Missing permission: expense_recognition:create")
		}
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		if deps == nil || deps.RecognizeFromExpenditure == nil {
			return centymo.HTMXError("recognize-expense handler not wired")
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("invalid form data")
		}

		expenditureID := viewCtx.Request.FormValue("expenditure_id")
		if expenditureID == "" {
			return centymo.HTMXError("expenditure_id is required")
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

		if _, err := deps.RecognizeFromExpenditure(ctx, req); err != nil {
			log.Printf("RecognizeExpense (supplier_subscription %s, expenditure %s): %v",
				viewCtx.Request.PathValue("id"), expenditureID, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("expense-recognitions-table")
	})
}
