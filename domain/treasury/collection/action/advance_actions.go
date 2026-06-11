// 20260517-advance-cash-events Plan B Phase 4 — Settle / Refund / Cancel
// drawer handlers for UNSCHEDULED advance Collections.
//
// Each handler is a single GET (drawer form) / POST (submit) split. The GET
// returns a drawer-form template; the POST calls the view-typed workflow
// closure the block layer wires from the espyna SettleUnscheduledAdvance /
// RefundUnscheduledAdvance / CancelAdvance use cases. When the closure is
// nil (the wiring hasn't landed in service-admin yet) the POST returns a
// safe "not configured" error and the GET still renders so operators can see
// the drawer shape.
package action

import (
	"context"
	"log"
	"net/http"

	collection "github.com/erniealice/centymo-golang/domain/treasury/collection"
	shared "github.com/erniealice/centymo-golang/domain/treasury/shared"
	"github.com/erniealice/pyeza-golang/view"
)

// AdvanceActionDeps holds the workflow closures + label tables the advance
// drawers need. Kept separate from the base Deps struct because the workflow
// closures aren't proto-shaped and the existing CRUD Deps would otherwise
// grow unrelated fields.
type AdvanceActionDeps struct {
	Routes        collection.Routes
	Labels        collection.Labels
	AdvanceLabels shared.TreasuryAdvanceLabels
	EnumLabels    shared.AdvanceEnumLabels
	CommonLabels  any

	// Workflow closures (nil-safe). Bound by service-admin's adapter.
	SettleUnscheduled func(ctx context.Context, in shared.AdvanceSettleViewInput) (*shared.AdvanceSettleViewOutput, error)
	RefundUnscheduled func(ctx context.Context, in shared.AdvanceRefundViewInput) (*shared.AdvanceRefundViewOutput, error)
	Cancel            func(ctx context.Context, in shared.AdvanceCancelViewInput) (*shared.AdvanceCancelViewOutput, error)
}

// AdvanceDrawerData is the template data shape for the Settle/Refund/Cancel
// drawer-form partials.
type AdvanceDrawerData struct {
	FormAction        string
	WorkspaceID       string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	AdvanceID         string
	Action            string // "settle" | "refund" | "cancel"
	Labels            shared.TreasuryAdvanceLabels
	EnumLabels        shared.AdvanceEnumLabels
	CommonLabels      any
	ShowAmount        bool
	ShowTargetAccount bool
	ShowReason        bool
	ShowDestination   bool
	ConfirmTitle      string
}

// NewSettleAction is the Settle drawer (GET = form, POST = use case).
func NewSettleAction(deps *AdvanceActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("treasury_collection", "settle") && !perms.Can("collection", "update") {
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
		input := shared.AdvanceSettleViewInput{
			AdvanceID:       id,
			Amount:          parseAmount(r.FormValue("amount")),
			TargetAccountID: r.FormValue("target_account_id"),
			Reason:          r.FormValue("reason"),
		}
		if _, err := deps.SettleUnscheduled(ctx, input); err != nil {
			log.Printf("Failed to settle advance %s: %v", id, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("collections-table")
	})
}

// NewRefundAction is the Refund drawer (GET = form, POST = use case).
func NewRefundAction(deps *AdvanceActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("treasury_collection", "refund") && !perms.Can("collection", "update") {
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
		input := shared.AdvanceRefundViewInput{
			AdvanceID:          id,
			Amount:             parseAmount(r.FormValue("amount")),
			RefundMethod:       r.FormValue("refund_method"),
			DestinationAccount: r.FormValue("destination_account"),
			Reason:             r.FormValue("reason"),
		}
		if _, err := deps.RefundUnscheduled(ctx, input); err != nil {
			log.Printf("Failed to refund advance %s: %v", id, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("collections-table")
	})
}

// NewCancelAction is the Cancel drawer (GET = form, POST = use case).
func NewCancelAction(deps *AdvanceActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("treasury_collection", "cancel") && !perms.Can("collection", "update") {
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
		input := shared.AdvanceCancelViewInput{
			AdvanceID: id,
			Reason:    viewCtx.Request.FormValue("reason"),
		}
		if _, err := deps.Cancel(ctx, input); err != nil {
			log.Printf("Failed to cancel advance %s: %v", id, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("collections-table")
	})
}

// formActionFor substitutes {id} in a URL pattern; if not found returns the
// pattern unchanged. Kept inline to avoid pulling in pyeza/route here.
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
