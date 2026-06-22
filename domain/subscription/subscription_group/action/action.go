package action

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	subscription_group "github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
)

type Deps struct {
	Routes                       subscription_group.Routes
	Labels                       subscription_group.Labels
	CreateSubscriptionGroup      func(ctx context.Context, req *subscriptiongrouppb.CreateSubscriptionGroupRequest) (*subscriptiongrouppb.CreateSubscriptionGroupResponse, error)
	ReadSubscriptionGroup        func(ctx context.Context, req *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error)
	UpdateSubscriptionGroup      func(ctx context.Context, req *subscriptiongrouppb.UpdateSubscriptionGroupRequest) (*subscriptiongrouppb.UpdateSubscriptionGroupResponse, error)
	DeleteSubscriptionGroup      func(ctx context.Context, req *subscriptiongrouppb.DeleteSubscriptionGroupRequest) (*subscriptiongrouppb.DeleteSubscriptionGroupResponse, error)
	ListPlans                    func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules           func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	GetSubscriptionGroupInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

func loadPlanPairs(ctx context.Context, deps *Deps) []form.Pair {
	if deps.ListPlans == nil {
		return nil
	}
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return nil
	}
	pairs := make([]form.Pair, 0, len(resp.GetData()))
	for _, p := range resp.GetData() {
		if p == nil || !p.GetActive() {
			continue
		}
		label := p.GetName()
		if label == "" {
			label = p.GetId()
		}
		pairs = append(pairs, form.Pair{ID: p.GetId(), Label: label})
	}
	return pairs
}

func loadSchedulePairs(ctx context.Context, deps *Deps) []form.Pair {
	if deps.ListPriceSchedules == nil {
		return nil
	}
	resp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
	if err != nil {
		return nil
	}
	pairs := make([]form.Pair, 0, len(resp.GetData()))
	for _, ps := range resp.GetData() {
		if ps == nil || !ps.GetActive() {
			continue
		}
		label := ps.GetName()
		if label == "" {
			label = ps.GetId()
		}
		pairs = append(pairs, form.Pair{ID: ps.GetId(), Label: label})
	}
	return pairs
}

// applyFormToData writes the POST body onto a SubscriptionGroup. Shared by
// Add (no id) and Edit (id set by caller). Optional FK pointers and
// max_capacity are set only when present so unset fields stay null.
func applyFormToData(r *http.Request) *subscriptiongrouppb.SubscriptionGroup {
	data := &subscriptiongrouppb.SubscriptionGroup{
		Name:   r.FormValue("name"),
		Kind:   r.FormValue("kind"),
		Active: r.FormValue("active") == "true",
	}
	if v := strings.TrimSpace(r.FormValue("plan_id")); v != "" {
		data.PlanId = &v
	}
	if v := strings.TrimSpace(r.FormValue("price_schedule_id")); v != "" {
		data.PriceScheduleId = &v
	}
	if v := r.FormValue("capacity_mode"); v != "" {
		if cm, ok := subscriptiongrouppb.CapacityMode_value[v]; ok {
			data.CapacityMode = subscriptiongrouppb.CapacityMode(cm)
		}
	}
	// max_capacity is read only when CAPPED; still persist whatever the
	// operator typed so switching back to CAPPED keeps the prior seat count.
	if s := strings.TrimSpace(r.FormValue("max_capacity")); s != "" {
		if n, err := strconv.ParseInt(s, 10, 32); err == nil {
			v32 := int32(n)
			data.MaxCapacity = &v32
		}
	}
	return data
}

func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group", "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			planPairs := loadPlanPairs(ctx, deps)
			schedulePairs := loadSchedulePairs(ctx, deps)
			return view.OK("subscription-group-drawer-form", &form.Data{
				FormAction:           deps.Routes.AddURL,
				Active:               true,
				Kind:                 "cohort",
				CapacityMode:         "CAPACITY_MODE_UNLIMITED",
				KindOptions:          form.BuildKindOptions(deps.Labels.Form, "cohort"),
				CapacityModeOptions:  form.BuildCapacityModeOptions(deps.Labels.Form, "CAPACITY_MODE_UNLIMITED"),
				PlanOptions:          form.BuildAutoCompleteOptions(planPairs, ""),
				PriceScheduleOptions: form.BuildAutoCompleteOptions(schedulePairs, ""),
				Labels:               deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		req := &subscriptiongrouppb.CreateSubscriptionGroupRequest{Data: applyFormToData(viewCtx.Request)}
		if _, err := deps.CreateSubscriptionGroup(ctx, req); err != nil {
			log.Printf("Failed to create subscription group: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-groups-table")
	})
}

// NewEditAction creates the subscription-group edit action (GET = form,
// POST = update). When the GET request includes ?clone=1, the handler returns
// the drawer pre-populated from the source record but wired to AddURL.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("subscription_group", requiredAction) {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSubscriptionGroup(ctx, &subscriptiongrouppb.ReadSubscriptionGroupRequest{Data: &subscriptiongrouppb.SubscriptionGroup{Id: id}})
			if err != nil || len(resp.GetData()) == 0 {
				return view.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := resp.GetData()[0]

			name := record.GetName()
			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id)
			formID := id
			if isClone {
				name = strings.TrimSpace(name) + viewCtx.T("actions.copySuffix")
				formAction = deps.Routes.AddURL
				formID = ""
			}

			planPairs := loadPlanPairs(ctx, deps)
			schedulePairs := loadSchedulePairs(ctx, deps)
			selectedPlanID := record.GetPlanId()
			selectedScheduleID := record.GetPriceScheduleId()
			capacityMode := subscriptiongrouppb.CapacityMode_name[int32(record.GetCapacityMode())]
			maxCapacity := ""
			if record.MaxCapacity != nil {
				maxCapacity = strconv.FormatInt(int64(record.GetMaxCapacity()), 10)
			}

			return view.OK("subscription-group-drawer-form", &form.Data{
				FormAction:           formAction,
				IsEdit:               !isClone,
				ID:                   formID,
				Name:                 name,
				Kind:                 record.GetKind(),
				PlanID:               selectedPlanID,
				PlanLabel:            form.FindLabel(planPairs, selectedPlanID),
				PlanOptions:          form.BuildAutoCompleteOptions(planPairs, selectedPlanID),
				PriceScheduleID:      selectedScheduleID,
				PriceScheduleLabel:   form.FindLabel(schedulePairs, selectedScheduleID),
				PriceScheduleOptions: form.BuildAutoCompleteOptions(schedulePairs, selectedScheduleID),
				KindOptions:          form.BuildKindOptions(deps.Labels.Form, record.GetKind()),
				CapacityMode:         capacityMode,
				CapacityModeOptions:  form.BuildCapacityModeOptions(deps.Labels.Form, capacityMode),
				MaxCapacity:          maxCapacity,
				Active:               record.GetActive(),
				Labels:               deps.Labels.Form,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		data := applyFormToData(viewCtx.Request)
		data.Id = id
		if _, err := deps.UpdateSubscriptionGroup(ctx, &subscriptiongrouppb.UpdateSubscriptionGroupRequest{Data: data}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-groups-table")
	})
}

func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group", "delete") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		if deps.GetSubscriptionGroupInUseIDs != nil {
			if inUse, _ := deps.GetSubscriptionGroupInUseIDs(ctx, []string{id}); inUse[id] {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeleteSubscriptionGroup(ctx, &subscriptiongrouppb.DeleteSubscriptionGroupRequest{Data: &subscriptiongrouppb.SubscriptionGroup{Id: id}}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-groups-table")
	})
}

func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group", "delete") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.DeleteFailed)
		}
		ids := viewCtx.Request.Form["id"]
		var attempted []string
		for _, id := range ids {
			if id != "" {
				attempted = append(attempted, id)
			}
		}
		if len(attempted) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		var inUse map[string]bool
		if deps.GetSubscriptionGroupInUseIDs != nil {
			inUse, _ = deps.GetSubscriptionGroupInUseIDs(ctx, attempted)
		}
		var deleted, blocked, failed int
		for _, id := range attempted {
			if inUse[id] {
				blocked++
				continue
			}
			if _, err := deps.DeleteSubscriptionGroup(ctx, &subscriptiongrouppb.DeleteSubscriptionGroupRequest{Data: &subscriptiongrouppb.SubscriptionGroup{Id: id}}); err != nil {
				log.Printf("Failed to delete subscription group %s during bulk: %v", id, err)
				failed++
				continue
			}
			deleted++
		}
		if deleted == 0 {
			if blocked > 0 && failed == 0 {
				return view.HTMXError(deps.Labels.Errors.InUse)
			}
			return view.HTMXError(deps.Labels.Errors.DeleteFailed)
		}
		return view.HTMXSuccess("subscription-groups-table")
	})
}

func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadSubscriptionGroup(ctx, &subscriptiongrouppb.ReadSubscriptionGroupRequest{Data: &subscriptiongrouppb.SubscriptionGroup{Id: id}})
		if err != nil || len(readResp.GetData()) == 0 {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		_, err = deps.UpdateSubscriptionGroup(ctx, &subscriptiongrouppb.UpdateSubscriptionGroupRequest{
			Data: &subscriptiongrouppb.SubscriptionGroup{
				Id:              id,
				Name:            record.GetName(),
				Kind:            record.GetKind(),
				PlanId:          record.PlanId,
				PriceScheduleId: record.PriceScheduleId,
				CapacityMode:    record.GetCapacityMode(),
				MaxCapacity:     record.MaxCapacity,
				Active:          status == "active",
			},
		})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("subscription-groups-table")
	})
}

func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group", "update") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadSubscriptionGroup(ctx, &subscriptiongrouppb.ReadSubscriptionGroupRequest{Data: &subscriptiongrouppb.SubscriptionGroup{Id: id}})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			_, _ = deps.UpdateSubscriptionGroup(ctx, &subscriptiongrouppb.UpdateSubscriptionGroupRequest{
				Data: &subscriptiongrouppb.SubscriptionGroup{
					Id:              id,
					Name:            record.GetName(),
					Kind:            record.GetKind(),
					PlanId:          record.PlanId,
					PriceScheduleId: record.PriceScheduleId,
					CapacityMode:    record.GetCapacityMode(),
					MaxCapacity:     record.MaxCapacity,
					Active:          status == "active",
				},
			})
		}
		return view.HTMXSuccess("subscription-groups-table")
	})
}
