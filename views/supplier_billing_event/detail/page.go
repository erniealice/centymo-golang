// Package detail — minimal detail view for supplier_billing_event rows.
//
// 20260517-advance-cash-events Plan B Phase 7. Renders the event metadata
// (status, trigger, billable amount, supplier subscription / contract, the
// expense_recognition back-edge if set) and surfaces the Recognize CTA when
// the event is BILLED and unconsumed.
package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	supplierbillingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_billing_event"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes       centymo.TreasuryAdvancesRoutes
	Labels       centymo.SupplierBillingEventLabels
	CommonLabels pyeza.CommonLabels

	ReadSupplierBillingEvent func(ctx context.Context, req *supplierbillingeventpb.ReadSupplierBillingEventRequest) (*supplierbillingeventpb.ReadSupplierBillingEventResponse, error)
}

// PageData holds the data for the supplier_billing_event detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          centymo.SupplierBillingEventLabels
	Event           map[string]any
}

// NewView creates the supplier_billing_event detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if perms != nil && !perms.Can("supplier_billing_event", "read") && !perms.Can("supplier_subscription", "read") {
			return view.Forbidden("supplier_billing_event:read")
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("id is required"))
		}
		if deps.ReadSupplierBillingEvent == nil {
			return view.Error(fmt.Errorf("supplier_billing_event read unavailable"))
		}
		resp, err := deps.ReadSupplierBillingEvent(ctx, &supplierbillingeventpb.ReadSupplierBillingEventRequest{
			Data: &supplierbillingeventpb.SupplierBillingEvent{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read supplier_billing_event %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load event: %w", err))
		}
		if resp == nil || len(resp.GetData()) == 0 {
			return view.Error(fmt.Errorf("supplier_billing_event not found"))
		}
		ev := resp.GetData()[0]
		l := deps.Labels
		eventMap := eventToMap(ev, l, deps.Routes)
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Detail.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    l.Detail.Title,
				HeaderSubtitle: ev.GetId(),
				HeaderIcon:     "icon-milestone",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "supplier-billing-event-detail-content",
			Labels:          l,
			Event:           eventMap,
		}
		return view.OK("supplier-billing-event-detail", pageData)
	})
}

func eventToMap(ev *supplierbillingeventpb.SupplierBillingEvent, l centymo.SupplierBillingEventLabels, routes centymo.TreasuryAdvancesRoutes) map[string]any {
	status := ev.GetStatus()
	statusK := statusKey(status)
	hasRec := ev.GetExpenseRecognitionId() != ""
	showRecognize := status == supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_BILLED && !hasRec
	return map[string]any{
		"id":                     ev.GetId(),
		"supplier_subscription":  ev.GetSupplierSubscriptionId(),
		"supplier_contract":      ev.GetSupplierContractId(),
		"billable_amount":        ev.GetBillableAmount(),
		"billable_display":       fmt.Sprintf("%.2f", float64(ev.GetBillableAmount())/100),
		"currency":               ev.GetBillingCurrency(),
		"status_key":             statusK,
		"status_label":           statusLabel(status, l.Status),
		"trigger_label":          triggerLabel(ev.GetTrigger(), l.Trigger),
		"expense_recognition_id": ev.GetExpenseRecognitionId(),
		"show_recognize":         showRecognize,
		"recognize_url":          route.ResolveURL(routes.SupplierBillingEventRecognizeURL, "id", ev.GetId()),
	}
}

func statusKey(s supplierbillingeventpb.SupplierBillingEventStatus) string {
	switch s {
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_READY:
		return "ready"
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_BILLED:
		return "billed"
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_WAIVED:
		return "waived"
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_CANCELLED:
		return "cancelled"
	default:
		return "unspecified"
	}
}

func statusLabel(s supplierbillingeventpb.SupplierBillingEventStatus, l centymo.SupplierBillingEventStatusLabels) string {
	switch s {
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_READY:
		return l.Ready
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_BILLED:
		return l.Billed
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_WAIVED:
		return l.Waived
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_CANCELLED:
		return l.Cancelled
	default:
		return l.Unspecified
	}
}

func triggerLabel(t supplierbillingeventpb.SupplierBillingEventTrigger, l centymo.SupplierBillingEventTriggerLabels) string {
	switch t {
	case supplierbillingeventpb.SupplierBillingEventTrigger_SUPPLIER_BILLING_EVENT_TRIGGER_MANUAL_EARLY:
		return l.ManualEarly
	case supplierbillingeventpb.SupplierBillingEventTrigger_SUPPLIER_BILLING_EVENT_TRIGGER_MANUAL_LATE:
		return l.ManualLate
	default:
		return l.Unspecified
	}
}
