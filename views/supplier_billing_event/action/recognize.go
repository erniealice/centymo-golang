// Package action — Recognize handler for a SupplierBillingEvent linked to a
// MILESTONE advance Disbursement.
//
// 20260517-advance-cash-events Plan B Phase 7. Mounted at
// `/action/supplier-billing-event/recognize/{id}`. POST submits to the
// view-typed closure
// `useCases.TreasuryAdvances.RecognizeMilestoneAdvanceDisbursement` which
// service-admin wires from the espyna
// RecognizeMilestoneAdvanceDisbursement use case.
//
// The form payload may carry `advance_id` (preferred — the
// treasury_disbursement that owns the linked junction). If absent the
// handler attempts to derive it from the junction reader; for v1 the form
// path is the supported one.
package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
)

// RecognizeMilestoneAdvanceDisbursementFn is the view-typed closure the
// centymo block adapter wires from espyna's
// RecognizeMilestoneAdvanceDisbursement use case.
type RecognizeMilestoneAdvanceDisbursementFn func(
	ctx context.Context,
	in centymo.AdvanceRecognizeMilestoneInput,
) (*centymo.AdvanceRecognizeMilestoneOutput, error)

// NewRecognizeAction creates the POST handler for
// `/action/supplier-billing-event/recognize/{id}`.
//
// RBAC: `supplier_billing_event:recognize` falls back to
// `expense_recognition:create`.
//
// On success: HX-Trigger refresh-supplier-billing-events +
// refresh-advance-schedule so any open list / detail panes refresh.
func NewRecognizeAction(
	recognize RecognizeMilestoneAdvanceDisbursementFn,
	errLabels centymo.SupplierBillingEventErrorLabels,
) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if perms != nil && !perms.Can("supplier_billing_event", "recognize") && !perms.Can("expense_recognition", "create") {
			return centymo.HTMXError(errLabels.PermissionDenied)
		}
		if viewCtx.Request.Method != http.MethodPost {
			return centymo.HTMXError(errLabels.InvalidTransition)
		}
		if recognize == nil {
			return centymo.HTMXError(errLabels.InvalidTransition)
		}
		eventID := viewCtx.Request.PathValue("id")
		if eventID == "" {
			return centymo.HTMXError(errLabels.NotFound)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(errLabels.InvalidTransition)
		}
		advanceID := strings.TrimSpace(viewCtx.Request.FormValue("advance_id"))
		if advanceID == "" {
			// v1: surface a clear validation message rather than silently
			// looking up junctions — the caller must pass advance_id.
			return centymo.HTMXError(errLabels.NotFound)
		}
		out, err := recognize(ctx, centymo.AdvanceRecognizeMilestoneInput{
			AdvanceID: advanceID,
			EventID:   eventID,
		})
		if err != nil {
			log.Printf("Recognize supplier milestone failed (advance=%s event=%s): %v", advanceID, eventID, err)
			return centymo.HTMXError(err.Error())
		}
		_ = out
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger": `{"refresh-supplier-billing-events":true,"refresh-advance-schedule":true}`,
			},
		}
	})
}
