package action

import (
	"context"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/supplier_subscription/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"
	pyezatypes "github.com/erniealice/pyeza-golang/types"

	suppliersubscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_subscription"
)

// NewEditAction creates the supplier_subscription edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			// Prefer joined item-page-data for enriched supplier + cost plan fields.
			var record *suppliersubscriptionpb.SupplierSubscription
			if deps.GetSupplierSubscriptionItemPageData != nil {
				resp, err := deps.GetSupplierSubscriptionItemPageData(ctx, &suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataRequest{
					SupplierSubscriptionId: id,
				})
				if err != nil || resp == nil || resp.GetSupplierSubscription() == nil {
					log.Printf("Failed to read supplier subscription %s: %v", id, err)
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetSupplierSubscription()
			} else {
				readResp, err := deps.ReadSupplierSubscription(ctx, &suppliersubscriptionpb.ReadSupplierSubscriptionRequest{
					Data: &suppliersubscriptionpb.SupplierSubscription{Id: id},
				})
				if err != nil || len(readResp.GetData()) == 0 {
					return centymo.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = readResp.GetData()[0]
			}

			tz := pyezatypes.LocationFromContext(ctx)
			startDate, startTime, startISO := splitTimestampForInputs(record.GetDateTimeStart(), tz)
			endDate, endTime, endISO := splitTimestampForInputs(record.GetDateTimeEnd(), tz)

			// Resolve cost plan label from nested object when available.
			costPlanLabel := record.GetCostPlanId()
			if cp := record.GetCostPlan(); cp != nil && cp.GetName() != "" {
				costPlanLabel = cp.GetName()
			}

			// Resolve supplier label from nested object when available.
			supplierLabel := record.GetSupplierId()

			// Extract notes from metadata.
			notes := ""
			if meta := record.GetMetadata(); meta != nil {
				notes = meta["notes"]
			}

			labels := buildFormLabels(deps.Labels)
			return view.OK("supplier-subscription-drawer-form", &form.Data{
				FormAction:        route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:            true,
				ID:                id,
				Name:              record.GetName(),
				Code:              record.GetCode(),
				CostPlanID:        record.GetCostPlanId(),
				CostPlanLabel:     costPlanLabel,
				SupplierID:        record.GetSupplierId(),
				SupplierLabel:     supplierLabel,
				AutoRenew:         record.GetAutoRenew(),
				Active:            record.GetActive(),
				DateStartDate:     startDate,
				DateStartTime:     startTime,
				DateStartISO:      startISO,
				DateEndDate:       endDate,
				DateEndTime:       endTime,
				DateEndISO:        endISO,
				DefaultTZ:         tz.String(),
				Notes:             notes,
				SearchCostPlanURL: deps.Routes.SearchCostPlanURL,
				SearchSupplierURL: deps.Routes.SearchSupplierURL,
				Labels:            labels,
			})
		}

		// POST — update supplier subscription
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		tz := pyezatypes.LocationFromContext(ctx)

		dateTimeStart := parseFormDateTime(
			r.FormValue("date_start_date"),
			r.FormValue("date_start_time"),
			r.FormValue("date_time_start_iso"),
			tz, false,
		)
		dateTimeEnd := parseFormDateTime(
			r.FormValue("date_end_date"),
			r.FormValue("date_end_time"),
			r.FormValue("date_time_end_iso"),
			tz, true,
		)

		code := r.FormValue("code")
		name := r.FormValue("name")
		if name == "" {
			name = code
		}
		autoRenew := r.FormValue("auto_renew") == "true" || r.FormValue("auto_renew") == "on"
		active := r.FormValue("active") != "false"

		req := &suppliersubscriptionpb.UpdateSupplierSubscriptionRequest{
			Data: &suppliersubscriptionpb.SupplierSubscription{
				Id:            id,
				Name:          name,
				Code:          strPtr(code),
				CostPlanId:    r.FormValue("cost_plan_id"),
				SupplierId:    r.FormValue("supplier_id"),
				DateTimeStart: dateTimeStart,
				DateTimeEnd:   dateTimeEnd,
				AutoRenew:     autoRenew,
				Active:        active,
			},
		}

		if notes := r.FormValue("notes"); notes != "" {
			if req.Data.Metadata == nil {
				req.Data.Metadata = map[string]string{}
			}
			req.Data.Metadata["notes"] = notes
		}

		if _, err := deps.UpdateSupplierSubscription(ctx, req); err != nil {
			log.Printf("Failed to update supplier subscription %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("supplier-subscriptions-table")
	})
}
