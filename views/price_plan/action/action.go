package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/price_plan/form"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
)

type PlanOption struct {
	Id   string
	Name string
}

type ScheduleOption struct {
	Id   string
	Name string
}

type FormData struct {
	FormAction            string
	IsEdit                bool
	ID                    string
	Name                  string
	Description           string
	Amount                string
	Currency              string
	DurationValue         string
	DurationUnit          string
	Active                bool
	PlanID                string
	SelectedPlanID        string
	SelectedPlanLabel     string
	Plans                 []*PlanOption
	PlanOptions           []map[string]any
	SelectedScheduleID    string
	SelectedScheduleLabel string
	Schedules             []*ScheduleOption
	ScheduleOptions       []map[string]any
	Labels                centymo.PricePlanFormLabels
	CommonLabels          pyeza.CommonLabels
}

type Deps struct {
	Routes                 centymo.PricePlanRoutes
	Labels                 centymo.PricePlanLabels
	CommonLabels           pyeza.CommonLabels
	CreatePricePlan        func(ctx context.Context, req *priceplanpb.CreatePricePlanRequest) (*priceplanpb.CreatePricePlanResponse, error)
	ReadPricePlan          func(ctx context.Context, req *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error)
	UpdatePricePlan        func(ctx context.Context, req *priceplanpb.UpdatePricePlanRequest) (*priceplanpb.UpdatePricePlanResponse, error)
	DeletePricePlan        func(ctx context.Context, req *priceplanpb.DeletePricePlanRequest) (*priceplanpb.DeletePricePlanResponse, error)
	ListPlans              func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules     func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	GetPricePlanInUseIDs   func(ctx context.Context, ids []string) (map[string]bool, error)
}

func loadPlans(ctx context.Context, deps *Deps) []*PlanOption {
	if deps.ListPlans == nil {
		return nil
	}
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return nil
	}
	opts := make([]*PlanOption, 0, len(resp.GetData()))
	for _, p := range resp.GetData() {
		opts = append(opts, &PlanOption{Id: p.GetId(), Name: p.GetName()})
	}
	return opts
}

func loadSchedules(ctx context.Context, deps *Deps) []*ScheduleOption {
	if deps.ListPriceSchedules == nil {
		return nil
	}
	resp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
	if err != nil {
		return nil
	}
	opts := make([]*ScheduleOption, 0, len(resp.GetData()))
	for _, s := range resp.GetData() {
		if !s.GetActive() {
			continue
		}
		opts = append(opts, &ScheduleOption{Id: s.GetId(), Name: s.GetName()})
	}
	return opts
}

func buildPlanAutoCompleteOptions(plans []*PlanOption, selectedID string) []map[string]any {
	opts := make([]map[string]any, 0, len(plans))
	for _, p := range plans {
		opts = append(opts, map[string]any{
			"Value":    p.Id,
			"Label":    p.Name,
			"Selected": p.Id == selectedID,
		})
	}
	return opts
}

func findPlanLabel(plans []*PlanOption, id string) string {
	for _, p := range plans {
		if p.Id == id {
			return p.Name
		}
	}
	return ""
}

func buildScheduleAutoCompleteOptions(schedules []*ScheduleOption, selectedID string) []map[string]any {
	opts := make([]map[string]any, 0, len(schedules))
	for _, s := range schedules {
		opts = append(opts, map[string]any{
			"Value":    s.Id,
			"Label":    s.Name,
			"Selected": s.Id == selectedID,
		})
	}
	return opts
}

func findScheduleLabel(schedules []*ScheduleOption, id string) string {
	for _, s := range schedules {
		if s.Id == id {
			return s.Name
		}
	}
	return ""
}

func parseAmount(s string) int64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f * 100))
}

func formatAmount(centavos int64) string {
	return strconv.FormatFloat(float64(centavos)/100.0, 'f', 2, 64)
}

func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if viewCtx.Request.Method == http.MethodGet {
			plans := loadPlans(ctx, deps)
			schedules := loadSchedules(ctx, deps)
			return view.OK("price-plan-drawer-form", &form.Data{
				FormAction:      deps.Routes.AddURL,
				Context:         form.ContextStandalone,
				Active:          true,
				Currency:        "PHP",
				DurationUnit:    "months",
				PlanOptions:     buildPlanAutoCompleteOptions(plans, ""),
				ScheduleOptions: buildScheduleAutoCompleteOptions(schedules, ""),
				Labels:          form.LabelsFromPricePlan(deps.Labels.Form),
				CommonLabels:    deps.CommonLabels,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		dv, _ := strconv.ParseInt(r.FormValue("duration_value"), 10, 32)
		scheduleID := r.FormValue("price_schedule_id")
		createName := r.FormValue("name")
		createDescription := r.FormValue("description")
		req := &priceplanpb.CreatePricePlanRequest{
			Data: &priceplanpb.PricePlan{
				PlanId:        r.FormValue("plan_id"),
				Name:          &createName,
				Description:   &createDescription,
				Amount:        parseAmount(r.FormValue("amount")),
				Currency:      r.FormValue("currency"),
				DurationValue: int32(dv),
				DurationUnit:  r.FormValue("duration_unit"),
				Active:        active,
			},
		}
		if scheduleID != "" {
			req.Data.PriceScheduleId = &scheduleID
		}
		if _, err := deps.CreatePricePlan(ctx, req); err != nil {
			log.Printf("Failed to create price plan: %v", err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}

func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}})
			if err != nil || len(resp.GetData()) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := resp.GetData()[0]
			plans := loadPlans(ctx, deps)
			schedules := loadSchedules(ctx, deps)
			selectedPlanID := record.GetPlanId()
			selectedScheduleID := record.GetPriceScheduleId()
			inUse := false
			lockMsg := ""
			if deps.GetPricePlanInUseIDs != nil {
				if m, _ := deps.GetPricePlanInUseIDs(ctx, []string{id}); m[id] {
					inUse = true
					lockMsg = "This price plan is in use by active subscriptions. Pricing changes are disabled."
				}
			}
			return view.OK("price-plan-drawer-form", &form.Data{
				FormAction:            route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                true,
				Context:               form.ContextStandalone,
				ID:                    id,
				PlanID:                selectedPlanID,
				ScheduleID:            selectedScheduleID,
				Name:                  record.GetName(),
				Description:           record.GetDescription(),
				Amount:                formatAmount(record.GetAmount()),
				Currency:              record.GetCurrency(),
				DurationValue:         fmt.Sprintf("%d", record.GetDurationValue()),
				DurationUnit:          record.GetDurationUnit(),
				Active:                record.GetActive(),
				PlanOptions:           buildPlanAutoCompleteOptions(plans, selectedPlanID),
				ScheduleOptions:       buildScheduleAutoCompleteOptions(schedules, selectedScheduleID),
				SelectedPlanID:        selectedPlanID,
				SelectedPlanLabel:     findPlanLabel(plans, selectedPlanID),
				SelectedScheduleID:    selectedScheduleID,
				SelectedScheduleLabel: findScheduleLabel(schedules, selectedScheduleID),
				InUse:                 inUse,
				LockMessage:           lockMsg,
				Labels:                form.LabelsFromPricePlan(deps.Labels.Form),
				CommonLabels:          deps.CommonLabels,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.UpdateFailed)
		}
		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		dv, _ := strconv.ParseInt(r.FormValue("duration_value"), 10, 32)
		scheduleID := r.FormValue("price_schedule_id")
		editName := r.FormValue("name")
		editDescription := r.FormValue("description")
		req := &priceplanpb.UpdatePricePlanRequest{
			Data: &priceplanpb.PricePlan{
				Id:            id,
				PlanId:        r.FormValue("plan_id"),
				Name:          &editName,
				Description:   &editDescription,
				Amount:        parseAmount(r.FormValue("amount")),
				Currency:      r.FormValue("currency"),
				DurationValue: int32(dv),
				DurationUnit:  r.FormValue("duration_unit"),
				Active:        active,
			},
		}
		if scheduleID != "" {
			req.Data.PriceScheduleId = &scheduleID
		}
		if _, err := deps.UpdatePricePlan(ctx, req); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}

func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		if deps.GetPricePlanInUseIDs != nil {
			if inUse, _ := deps.GetPricePlanInUseIDs(ctx, []string{id}); inUse[id] {
				return centymo.HTMXError(deps.Labels.Errors.InUse)
			}
		}
		if _, err := deps.DeletePricePlan(ctx, &priceplanpb.DeletePricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}}); err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}

func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.DeleteFailed)
		}
		for _, id := range viewCtx.Request.Form["id"] {
			if id != "" {
				_, _ = deps.DeletePricePlan(ctx, &priceplanpb.DeletePricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}})
			}
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}

func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			status = viewCtx.Request.FormValue("status")
		}
		readResp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}})
		if err != nil || len(readResp.GetData()) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		record := readResp.GetData()[0]
		statusName := record.GetName()
		statusDescription := record.GetDescription()
		_, err = deps.UpdatePricePlan(ctx, &priceplanpb.UpdatePricePlanRequest{
			Data: &priceplanpb.PricePlan{
				Id: id, PlanId: record.GetPlanId(), Name: &statusName,
				Description: &statusDescription, Amount: record.GetAmount(),
				Currency: record.GetCurrency(), DurationValue: record.GetDurationValue(),
				DurationUnit: record.GetDurationUnit(), Active: status == "active",
				PriceScheduleId: record.PriceScheduleId,
			},
		})
		if err != nil {
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}

func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)
		ids := viewCtx.Request.Form["id"]
		status := viewCtx.Request.FormValue("target_status")
		for _, id := range ids {
			if id == "" {
				continue
			}
			readResp, err := deps.ReadPricePlan(ctx, &priceplanpb.ReadPricePlanRequest{Data: &priceplanpb.PricePlan{Id: id}})
			if err != nil || len(readResp.GetData()) == 0 {
				continue
			}
			record := readResp.GetData()[0]
			bulkName := record.GetName()
			bulkDescription := record.GetDescription()
			_, _ = deps.UpdatePricePlan(ctx, &priceplanpb.UpdatePricePlanRequest{
				Data: &priceplanpb.PricePlan{
					Id: id, PlanId: record.GetPlanId(), Name: &bulkName,
					Description: &bulkDescription, Amount: record.GetAmount(),
					Currency: record.GetCurrency(), DurationValue: record.GetDurationValue(),
					DurationUnit: record.GetDurationUnit(), Active: status == "active",
					PriceScheduleId: record.PriceScheduleId,
				},
			})
		}
		return centymo.HTMXSuccess("price-plans-table")
	})
}
