// Package list — minimal list view for supplier_billing_event rows.
//
// 20260517-advance-cash-events Plan B Phase 7 — SupplierBillingEvent is the
// buying-side mirror of BillingEvent and was authored in Phase 0 per
// Decision 10 of the Plan B design. The list page surfaces all events
// regardless of supplier_subscription so operators can quickly find the
// row to Recognize when it links to an advance Disbursement.
package list

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/erniealice/centymo-golang/domain/expenditure"
	"github.com/erniealice/centymo-golang/domain/treasury"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	supplierbillingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_billing_event"
)

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes       treasury.TreasuryAdvancesRoutes
	Labels       expenditure.SupplierBillingEventLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListSupplierBillingEvents func(ctx context.Context, req *supplierbillingeventpb.ListSupplierBillingEventsRequest) (*supplierbillingeventpb.ListSupplierBillingEventsResponse, error)
}

// EventRow is the per-row view shape used by the template.
type EventRow struct {
	ID                   string
	SupplierSubscription string
	SupplierContract     string
	BillableDisplay      string
	Currency             string
	StatusKey            string
	StatusLabel          string
	TriggerLabel         string
	ExpenseRecognitionID string
	DetailURL            string
	ShowRecognize        bool
	RecognizeURL         string
}

// PageData holds the data for the supplier_billing_event list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Rows            []EventRow
	Labels          expenditure.SupplierBillingEventLabels
	Status          string
}

// NewView creates the supplier_billing_event list view (full page).
//
// URL mount: SupplierBillingEventListURL = /app/supplier-billing-events/list/{status}.
// {status} filters the rows by status — "all" returns everything; specific
// statuses map to the enum value.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if perms != nil && !perms.Can("supplier_billing_event", "list") && !perms.Can("supplier_subscription", "list") {
			return view.Forbidden("supplier_billing_event:list")
		}
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "all"
		}
		rows := loadRows(ctx, deps, status)
		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Page.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    l.Page.Title,
				HeaderSubtitle: l.Page.Caption,
				HeaderIcon:     "icon-milestone",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "supplier-billing-event-list-content",
			Rows:            rows,
			Labels:          l,
			Status:          status,
		}
		return view.OK("supplier-billing-event-list", pageData)
	})
}

func loadRows(ctx context.Context, deps *ListViewDeps, status string) []EventRow {
	if deps.ListSupplierBillingEvents == nil {
		return nil
	}
	resp, err := deps.ListSupplierBillingEvents(ctx, &supplierbillingeventpb.ListSupplierBillingEventsRequest{})
	if err != nil {
		log.Printf("Failed to list supplier_billing_events: %v", err)
		return nil
	}
	if resp == nil {
		return nil
	}
	wantStatus := strings.ToLower(status)
	rows := make([]EventRow, 0, len(resp.GetData()))
	for _, ev := range resp.GetData() {
		statusK := statusKey(ev.GetStatus())
		if wantStatus != "all" && wantStatus != "" && statusK != wantStatus {
			continue
		}
		showRecognize := ev.GetStatus() == supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_BILLED &&
			strings.TrimSpace(ev.GetExpenseRecognitionId()) == ""
		row := EventRow{
			ID:                   ev.GetId(),
			SupplierSubscription: ev.GetSupplierSubscriptionId(),
			SupplierContract:     ev.GetSupplierContractId(),
			BillableDisplay:      fmt.Sprintf("%.2f", float64(ev.GetBillableAmount())/100),
			Currency:             ev.GetBillingCurrency(),
			StatusKey:            statusK,
			StatusLabel:          statusLabel(ev.GetStatus(), deps.Labels.Status),
			TriggerLabel:         triggerLabel(ev.GetTrigger(), deps.Labels.Trigger),
			ExpenseRecognitionID: ev.GetExpenseRecognitionId(),
			DetailURL:            route.ResolveURL(deps.Routes.SupplierBillingEventDetailURL, "id", ev.GetId()),
			ShowRecognize:        showRecognize,
			RecognizeURL:         route.ResolveURL(deps.Routes.SupplierBillingEventRecognizeURL, "id", ev.GetId()),
		}
		rows = append(rows, row)
	}
	return rows
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

func statusLabel(s supplierbillingeventpb.SupplierBillingEventStatus, l expenditure.SupplierBillingEventStatusLabels) string {
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

func triggerLabel(t supplierbillingeventpb.SupplierBillingEventTrigger, l expenditure.SupplierBillingEventTriggerLabels) string {
	switch t {
	case supplierbillingeventpb.SupplierBillingEventTrigger_SUPPLIER_BILLING_EVENT_TRIGGER_MANUAL_EARLY:
		return l.ManualEarly
	case supplierbillingeventpb.SupplierBillingEventTrigger_SUPPLIER_BILLING_EVENT_TRIGGER_MANUAL_LATE:
		return l.ManualLate
	default:
		return l.Unspecified
	}
}
