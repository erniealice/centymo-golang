package action

import (
	"context"
	"net/http"
	"sort"
	"strings"

	"github.com/erniealice/centymo-golang/domain/subscription/plan"
	"github.com/erniealice/centymo-golang/domain/subscription/plan/form"
	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"
	planjobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/plan_job_template"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"
)

type PlanJobTemplateDeps struct {
	Routes           plan.Routes
	Labels           plan.Labels
	ListJobTemplates func(context.Context, *jobtemplatepb.ListJobTemplatesRequest) (*jobtemplatepb.ListJobTemplatesResponse, error)
	Create           func(context.Context, *planjobtemplatepb.CreatePlanJobTemplateRequest) (*planjobtemplatepb.CreatePlanJobTemplateResponse, error)
	Update           func(context.Context, *planjobtemplatepb.UpdatePlanJobTemplateRequest) (*planjobtemplatepb.UpdatePlanJobTemplateResponse, error)
	Delete           func(context.Context, *planjobtemplatepb.DeletePlanJobTemplateRequest) (*planjobtemplatepb.DeletePlanJobTemplateResponse, error)
	ListByPlan       func(context.Context, *planjobtemplatepb.ListPlanJobTemplatesByPlanRequest) (*planjobtemplatepb.ListPlanJobTemplatesByPlanResponse, error)
}

func loadPlanCompositionData(ctx context.Context, deps *PlanJobTemplateDeps, planID string) *form.PlanCompositionData {
	if deps == nil || deps.ListByPlan == nil || planID == "" {
		return nil
	}
	labels := deps.Labels.Composition
	data := &form.PlanCompositionData{
		PlanID:    planID,
		AddURL:    route.ResolveURL(deps.Routes.PlanJobTemplateAddURL, "id", planID),
		PickerURL: route.ResolveURL(deps.Routes.PlanJobTemplatePickerURL, "id", planID),
		Labels:    form.PlanCompositionLabels{Heading: labels.Page.Heading, Caption: labels.Page.Caption, AddEntry: labels.Buttons.AddEntry, MoveUp: labels.Buttons.MoveUp, MoveDown: labels.Buttons.MoveDown, Remove: labels.Buttons.Remove, Template: labels.Form.JobTemplateID, Pattern: labels.Form.CompositionEntryPattern, EmptyTitle: labels.Empty.Title, EmptyMessage: labels.Empty.Message},
		PatternOptions: []form.PlanCompositionOption{
			{Value: planjobtemplatepb.PlanJobTemplateCompositionEntryPattern_PLAN_JOB_TEMPLATE_COMPOSITION_ENTRY_PATTERN_BUNDLE_ENTRY.String(), Label: labels.Patterns.BundleEntry},
			{Value: planjobtemplatepb.PlanJobTemplateCompositionEntryPattern_PLAN_JOB_TEMPLATE_COMPOSITION_ENTRY_PATTERN_STANDALONE_ENTRY.String(), Label: labels.Patterns.StandaloneEntry},
		},
	}
	names := map[string]string{}
	if deps.ListJobTemplates != nil {
		if resp, err := deps.ListJobTemplates(ctx, &jobtemplatepb.ListJobTemplatesRequest{}); err == nil {
			for _, row := range resp.GetData() {
				if row != nil && row.GetActive() {
					names[row.GetId()] = row.GetName()
					data.TemplateOptions = append(data.TemplateOptions, form.JobTemplateOption{Value: row.GetId(), Label: row.GetName()})
				}
			}
			sort.SliceStable(data.TemplateOptions, func(i, j int) bool { return data.TemplateOptions[i].Label < data.TemplateOptions[j].Label })
		}
	}
	resp, err := deps.ListByPlan(ctx, &planjobtemplatepb.ListPlanJobTemplatesByPlanRequest{PlanId: planID})
	if err != nil || resp == nil {
		return data
	}
	rows := resp.GetPlanJobTemplates()
	for i, row := range rows {
		name := names[row.GetJobTemplateId()]
		if name == "" {
			name = row.GetJobTemplateId()
		}
		patternLabel := labels.Patterns.BundleEntry
		if row.GetCompositionEntryPattern() == planjobtemplatepb.PlanJobTemplateCompositionEntryPattern_PLAN_JOB_TEMPLATE_COMPOSITION_ENTRY_PATTERN_STANDALONE_ENTRY {
			patternLabel = labels.Patterns.StandaloneEntry
		}
		editURL := route.ResolveURL(deps.Routes.PlanJobTemplateEditURL, "id", planID, "pjtid", row.GetId())
		data.Entries = append(data.Entries, form.PlanCompositionEntry{ID: row.GetId(), JobTemplateID: row.GetJobTemplateId(), JobTemplateName: name, Pattern: row.GetCompositionEntryPattern().String(), PatternLabel: patternLabel, SequenceOrder: row.GetSequenceOrder(), MoveUpURL: editURL, MoveDownURL: editURL, DeleteURL: route.ResolveURL(deps.Routes.PlanJobTemplateDeleteURL, "id", planID), CanMoveUp: i > 0, CanMoveDown: i+1 < len(rows)})
	}
	return data
}

func NewPlanJobTemplateAddAction(deps *PlanJobTemplateDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, vc *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("plan_job_template", "create") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		planID := vc.Request.PathValue("id")
		if planID == "" {
			return view.HTMXError(deps.Labels.Errors.IDRequired)
		}
		if vc.Request.Method == http.MethodGet {
			return view.OK("plan-composition-entry-list", loadPlanCompositionData(ctx, deps, planID))
		}
		if err := vc.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		templateID := strings.TrimSpace(vc.Request.FormValue("composition_job_template_id"))
		patternName := strings.TrimSpace(vc.Request.FormValue("composition_entry_pattern"))
		if templateID == "" {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		pattern := planjobtemplatepb.PlanJobTemplateCompositionEntryPattern(planjobtemplatepb.PlanJobTemplateCompositionEntryPattern_value[patternName])
		if pattern == 0 {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		current := loadPlanCompositionData(ctx, deps, planID)
		order := int32(len(current.Entries) + 1)
		if _, err := deps.Create(ctx, &planjobtemplatepb.CreatePlanJobTemplateRequest{Data: &planjobtemplatepb.PlanJobTemplate{PlanId: planID, JobTemplateId: templateID, SequenceOrder: order, CompositionEntryPattern: pattern}}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.OK("plan-composition-entry-list", loadPlanCompositionData(ctx, deps, planID))
	})
}

func NewPlanJobTemplateEditAction(deps *PlanJobTemplateDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, vc *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("plan_job_template", "update") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		planID, entryID := vc.Request.PathValue("id"), vc.Request.PathValue("pjtid")
		if planID == "" || entryID == "" {
			return view.HTMXError(deps.Labels.Errors.IDRequired)
		}
		if err := vc.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		direction := vc.Request.FormValue("direction")
		resp, err := deps.ListByPlan(ctx, &planjobtemplatepb.ListPlanJobTemplatesByPlanRequest{PlanId: planID})
		if err != nil {
			return view.HTMXError(err.Error())
		}
		rows := resp.GetPlanJobTemplates()
		idx := -1
		for i, row := range rows {
			if row.GetId() == entryID {
				idx = i
				break
			}
		}
		target := idx
		if direction == "up" {
			target--
		} else if direction == "down" {
			target++
		} else {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		if idx < 0 || target < 0 || target >= len(rows) {
			return view.OK("plan-composition-entry-list", loadPlanCompositionData(ctx, deps, planID))
		}
		a, b := rows[idx], rows[target]
		ao, bo := a.GetSequenceOrder(), b.GetSequenceOrder()
		if _, err = deps.Update(ctx, &planjobtemplatepb.UpdatePlanJobTemplateRequest{Data: &planjobtemplatepb.PlanJobTemplate{Id: a.GetId(), SequenceOrder: bo}}); err != nil {
			return view.HTMXError(err.Error())
		}
		if _, err = deps.Update(ctx, &planjobtemplatepb.UpdatePlanJobTemplateRequest{Data: &planjobtemplatepb.PlanJobTemplate{Id: b.GetId(), SequenceOrder: ao}}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.OK("plan-composition-entry-list", loadPlanCompositionData(ctx, deps, planID))
	})
}

func NewPlanJobTemplateDeleteAction(deps *PlanJobTemplateDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, vc *view.ViewContext) view.ViewResult {
		if !view.GetUserPermissions(ctx).Can("plan_job_template", "delete") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		planID := vc.Request.PathValue("id")
		_ = vc.Request.ParseForm()
		entryID := vc.Request.FormValue("entry_id")
		if entryID == "" {
			return view.HTMXError(deps.Labels.Errors.IDRequired)
		}
		if _, err := deps.Delete(ctx, &planjobtemplatepb.DeletePlanJobTemplateRequest{Data: &planjobtemplatepb.PlanJobTemplate{Id: entryID}}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.OK("plan-composition-entry-list", loadPlanCompositionData(ctx, deps, planID))
	})
}

func NewPlanJobTemplatePickerAction(deps *PlanJobTemplateDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, vc *view.ViewContext) view.ViewResult {
		return view.OK("plan-composition-entry-list", loadPlanCompositionData(ctx, deps, vc.Request.PathValue("id")))
	})
}
