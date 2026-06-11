// 20260517-advance-cash-events Plan B Phase 4 — Settle / Refund / Cancel
// drawer handlers for UNSCHEDULED advance Disbursements. Mirrors the
// collection-side advance_actions.go file exactly (same shapes, same
// behaviour); only the permission codes / template names differ.
package action

import (
	"context"
	"log"
	"net/http"

	disbursement "github.com/erniealice/centymo-golang/domain/treasury/disbursement"
	shared "github.com/erniealice/centymo-golang/domain/treasury/shared"
	"github.com/erniealice/pyeza-golang/view"
)

// AdvanceActionDeps holds the workflow closures + label tables for the
// disbursement advance drawers.
type AdvanceActionDeps struct {
	Routes        disbursement.Routes
	Labels        disbursement.Labels
	AdvanceLabels shared.TreasuryAdvanceLabels
	EnumLabels    shared.AdvanceEnumLabels
	CommonLabels  any

	SettleUnscheduled func(ctx context.Context, in shared.AdvanceSettleViewInput) (*shared.AdvanceSettleViewOutput, error)
	RefundUnscheduled func(ctx context.Context, in shared.AdvanceRefundViewInput) (*shared.AdvanceRefundViewOutput, error)
	Cancel            func(ctx context.Context, in shared.AdvanceCancelViewInput) (*shared.AdvanceCancelViewOutput, error)
}

// AdvanceDrawerData is the per-template data carrier (mirrors collection-side).
type AdvanceDrawerData struct {
	FormAction        string
	WorkspaceID       string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	AdvanceID         string
	Action            string
	Labels            shared.TreasuryAdvanceLabels
	EnumLabels        shared.AdvanceEnumLabels
	CommonLabels      any
	ShowAmount        bool
	ShowTargetAccount bool
	ShowReason        bool
	ShowDestination   bool
	ConfirmTitle      string
}

// NewSettleAction wires the Settle drawer for advance Disbursements.
func NewSettleAction(deps *AdvanceActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("treasury_disbursement", "settle") && !perms.Can("disbursement", "update") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("advance-settle-drawer-form", &AdvanceDrawerData{
				FormAction:        formActionFor(deps.Routes.SettleURL, id),
				AdvanceID:         id,
				Action:            "settle",
				Labels:            deps.AdvanceLabels,
				EnumLabels:        deps.EnumLabels,
				CommonLabels:      deps.CommonLabels,
				ShowAmount:        true,
				ShowTargetAccount: true,
				ShowReason:        true,
				ConfirmTitle:      deps.AdvanceLabels.Actions.SettleConfirm,
			})
		}

		if deps.SettleUnscheduled == nil {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		if _, err := deps.SettleUnscheduled(ctx, shared.AdvanceSettleViewInput{
			AdvanceID:       id,
			Amount:          parseAmount(r.FormValue("amount")),
			TargetAccountID: r.FormValue("target_account_id"),
			Reason:          r.FormValue("reason"),
		}); err != nil {
			log.Printf("Failed to settle advance disbursement %s: %v", id, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("disbursements-table")
	})
}

// NewRefundAction wires the Refund drawer for advance Disbursements.
func NewRefundAction(deps *AdvanceActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("treasury_disbursement", "refund") && !perms.Can("disbursement", "update") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("advance-refund-drawer-form", &AdvanceDrawerData{
				FormAction:      formActionFor(deps.Routes.RefundURL, id),
				AdvanceID:       id,
				Action:          "refund",
				Labels:          deps.AdvanceLabels,
				EnumLabels:      deps.EnumLabels,
				CommonLabels:    deps.CommonLabels,
				ShowAmount:      true,
				ShowDestination: true,
				ShowReason:      true,
				ConfirmTitle:    deps.AdvanceLabels.Actions.RefundConfirm,
			})
		}

		if deps.RefundUnscheduled == nil {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		if _, err := deps.RefundUnscheduled(ctx, shared.AdvanceRefundViewInput{
			AdvanceID:          id,
			Amount:             parseAmount(r.FormValue("amount")),
			RefundMethod:       r.FormValue("refund_method"),
			DestinationAccount: r.FormValue("destination_account"),
			Reason:             r.FormValue("reason"),
		}); err != nil {
			log.Printf("Failed to refund advance disbursement %s: %v", id, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("disbursements-table")
	})
}

// NewCancelAction wires the Cancel drawer for advance Disbursements.
func NewCancelAction(deps *AdvanceActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("treasury_disbursement", "cancel") && !perms.Can("disbursement", "update") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("advance-cancel-drawer-form", &AdvanceDrawerData{
				FormAction:   formActionFor(deps.Routes.CancelURL, id),
				AdvanceID:    id,
				Action:       "cancel",
				Labels:       deps.AdvanceLabels,
				EnumLabels:   deps.EnumLabels,
				CommonLabels: deps.CommonLabels,
				ShowReason:   true,
				ConfirmTitle: deps.AdvanceLabels.Actions.CancelConfirm,
			})
		}

		if deps.Cancel == nil {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		if _, err := deps.Cancel(ctx, shared.AdvanceCancelViewInput{
			AdvanceID: id,
			Reason:    viewCtx.Request.FormValue("reason"),
		}); err != nil {
			log.Printf("Failed to cancel advance disbursement %s: %v", id, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("disbursements-table")
	})
}

// formActionFor substitutes {id} in the URL pattern; unchanged if absent.
func formActionFor(pattern, id string) string {
	if pattern == "" {
		return ""
	}
	idx := -1
	needle := "{id}"
	for i := 0; i+len(needle) <= len(pattern); i++ {
		if pattern[i:i+len(needle)] == needle {
			idx = i
			break
		}
	}
	if idx < 0 {
		return pattern
	}
	return pattern[:idx] + id + pattern[idx+len(needle):]
}
