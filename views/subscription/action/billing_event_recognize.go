// 20260517-advance-cash-events Plan B Phase 7 — Recognize action handler for
// a BillingEvent row that is linked to a MILESTONE advance Collection.
//
// Mounted at `/action/subscription/{id}/billing-event/{eventId}/recognize`.
// POST submits to the view-typed closure
// `useCases.TreasuryAdvances.RecognizeMilestoneAdvanceCollection` which
// service-admin wires from the espyna
// RecognizeMilestoneAdvanceCollection use case.
//
// The form payload carries `advance_id` because a single BillingEvent may
// theoretically link to multiple advances (one per junction row); the
// view-layer caller is responsible for pre-resolving the right advance and
// passing it in the form.
package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
)

// RecognizeMilestoneAdvanceCollectionFn is the view-typed closure shape the
// centymo block adapter wires from espyna's
// RecognizeMilestoneAdvanceCollection use case. Mirrors the shape on
// `useCases.TreasuryAdvances`.
type RecognizeMilestoneAdvanceCollectionFn func(
	ctx context.Context,
	in centymo.AdvanceRecognizeMilestoneInput,
) (*centymo.AdvanceRecognizeMilestoneOutput, error)

// NewMilestoneRecognizeAction creates the POST handler for
// `/action/subscription/{id}/billing-event/{eventId}/recognize`.
//
// RBAC: `milestone:recognize` falls back to `revenue:create` for builds where
// the milestone scope is unregistered.
//
// On success: HX-Trigger refresh-milestones + refresh-advance-schedule so the
// subscription Package tab + the Collection Advance Schedule tab refresh
// inline.
func NewMilestoneRecognizeAction(
	recognize RecognizeMilestoneAdvanceCollectionFn,
	errLabels centymo.SubscriptionErrorLabels,
) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if perms != nil && !perms.Can("milestone", "recognize") && !perms.Can("revenue", "create") {
			return view.HTMXError(errLabels.PermissionDenied)
		}
		if viewCtx.Request.Method != http.MethodPost {
			return view.HTMXError(errLabels.InvalidStatus)
		}
		if recognize == nil {
			return view.HTMXError(errLabels.InvalidFormData)
		}
		eventID := viewCtx.Request.PathValue("eventId")
		if eventID == "" {
			return view.HTMXError(errLabels.IDRequired)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(errLabels.InvalidFormData)
		}
		advanceID := strings.TrimSpace(viewCtx.Request.FormValue("advance_id"))
		if advanceID == "" {
			return view.HTMXError(errLabels.IDRequired)
		}

		out, err := recognize(ctx, centymo.AdvanceRecognizeMilestoneInput{
			AdvanceID: advanceID,
			EventID:   eventID,
		})
		if err != nil {
			log.Printf("Recognize milestone failed (advance=%s event=%s): %v", advanceID, eventID, err)
			return view.HTMXError(err.Error())
		}
		// SKIPPED is a benign outcome — the operator double-clicked or the
		// junction was already consumed. Surface as success so the UI
		// refreshes (the new state is reflected on read-back).
		_ = out
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger": `{"refresh-milestones":true,"refresh-advance-schedule":true,"refresh-revenues":true}`,
			},
		}
	})
}
