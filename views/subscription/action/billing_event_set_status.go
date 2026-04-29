package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	billingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billing_event"
)

// SetBillingEventStatusFn is the espyna BillingEvent.SetStatus signature.
// Wired by the centymo block when the BillingEvent server is registered;
// nil-safe on the centymo side — the routes only register when the function
// is non-nil so handlers cannot panic on a missing dependency.
type SetBillingEventStatusFn func(ctx context.Context, req *billingeventpb.SetBillingEventStatusRequest) (*billingeventpb.SetBillingEventStatusResponse, error)

// NewMilestoneMarkReadyAction creates the POST handler for
// `/action/subscription/{id}/billing-event/{eventId}/mark-ready`.
//
// 2026-04-29 milestone-billing plan §5 / Phase D — flips a BillingEvent's
// status to READY with trigger=MANUAL_LATE (operator pressed mark-ready
// rather than the JobPhase.COMPLETED hook firing automatically). RBAC:
// `milestone:set_status`.
//
// On success: HX-Trigger refresh-milestones + refresh-invoices so the
// Package tab + invoices tab refresh inline (matches the recognize-drawer
// success header bundle).
func NewMilestoneMarkReadyAction(setStatus SetBillingEventStatusFn, errLabels centymo.SubscriptionErrorLabels) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if perms != nil && !perms.Can("milestone", "set_status") {
			return centymo.HTMXError(errLabels.PermissionDenied)
		}
		if viewCtx.Request.Method != http.MethodPost {
			return centymo.HTMXError(errLabels.InvalidStatus)
		}
		if setStatus == nil {
			return centymo.HTMXError(errLabels.InvalidFormData)
		}
		eventID := viewCtx.Request.PathValue("eventId")
		if eventID == "" {
			return centymo.HTMXError(errLabels.IDRequired)
		}
		_ = viewCtx.Request.ParseForm()
		reason := strings.TrimSpace(viewCtx.Request.FormValue("reason"))
		req := &billingeventpb.SetBillingEventStatusRequest{
			BillingEventId: eventID,
			Status:         billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_READY,
			Trigger:        billingeventpb.BillingEventTrigger_BILLING_EVENT_TRIGGER_MANUAL_LATE,
		}
		if reason != "" {
			req.Reason = &reason
		}
		if _, err := setStatus(ctx, req); err != nil {
			log.Printf("Mark milestone ready failed for event %s: %v", eventID, err)
			return centymo.HTMXError(err.Error())
		}
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger": `{"refresh-milestones":true,"refresh-invoices":true}`,
			},
		}
	})
}

// NewMilestoneWaiveAction creates the POST handler for
// `/action/subscription/{id}/billing-event/{eventId}/waive`. Sets the event's
// status to WAIVED. RBAC: `milestone:set_status`.
func NewMilestoneWaiveAction(setStatus SetBillingEventStatusFn, errLabels centymo.SubscriptionErrorLabels) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if perms != nil && !perms.Can("milestone", "set_status") {
			return centymo.HTMXError(errLabels.PermissionDenied)
		}
		if viewCtx.Request.Method != http.MethodPost {
			return centymo.HTMXError(errLabels.InvalidStatus)
		}
		if setStatus == nil {
			return centymo.HTMXError(errLabels.InvalidFormData)
		}
		eventID := viewCtx.Request.PathValue("eventId")
		if eventID == "" {
			return centymo.HTMXError(errLabels.IDRequired)
		}
		_ = viewCtx.Request.ParseForm()
		reason := strings.TrimSpace(viewCtx.Request.FormValue("reason"))
		req := &billingeventpb.SetBillingEventStatusRequest{
			BillingEventId: eventID,
			Status:         billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_WAIVED,
			// Trigger preserved by the SetStatus implementation when not set.
			Trigger: billingeventpb.BillingEventTrigger_BILLING_EVENT_TRIGGER_UNSPECIFIED,
		}
		if reason != "" {
			req.Reason = &reason
		}
		if _, err := setStatus(ctx, req); err != nil {
			log.Printf("Waive milestone failed for event %s: %v", eventID, err)
			return centymo.HTMXError(err.Error())
		}
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger": `{"refresh-milestones":true}`,
			},
		}
	})
}
