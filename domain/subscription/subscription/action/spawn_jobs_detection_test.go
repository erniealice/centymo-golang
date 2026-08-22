package action

import (
	"context"
	"reflect"
	"testing"

	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"
	planjobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/plan_job_template"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
)

func TestDetectSpawnJobs_CompositionWithoutLegacyRoot(t *testing.T) {
	planID := "plan-1"
	templates := map[string]*jobtemplatepb.JobTemplate{
		"math":    {Id: "math", Name: "Mathematics", Active: true},
		"science": {Id: "science", Name: "Science", Active: true},
	}
	deps := &Deps{
		ReadPricePlan: func(context.Context, *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error) {
			return &priceplanpb.ReadPricePlanResponse{Data: []*priceplanpb.PricePlan{{Id: "price-plan-1", PlanId: planID}}}, nil
		},
		ReadPlan: func(context.Context, *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error) {
			return &planpb.ReadPlanResponse{Data: []*planpb.Plan{{Id: &planID, Active: true}}}, nil
		},
		ListPlanJobTemplatesByPlan: func(_ context.Context, req *planjobtemplatepb.ListPlanJobTemplatesByPlanRequest) (*planjobtemplatepb.ListPlanJobTemplatesByPlanResponse, error) {
			if req.GetPlanId() != planID {
				t.Fatalf("plan id = %q, want %q", req.GetPlanId(), planID)
			}
			return &planjobtemplatepb.ListPlanJobTemplatesByPlanResponse{PlanJobTemplates: []*planjobtemplatepb.PlanJobTemplate{
				{Id: "entry-1", PlanId: planID, JobTemplateId: "math", Active: true, SequenceOrder: 1, CompositionEntryPattern: planjobtemplatepb.PlanJobTemplateCompositionEntryPattern_PLAN_JOB_TEMPLATE_COMPOSITION_ENTRY_PATTERN_BUNDLE_ENTRY},
				{Id: "entry-2", PlanId: planID, JobTemplateId: "science", Active: true, SequenceOrder: 2, CompositionEntryPattern: planjobtemplatepb.PlanJobTemplateCompositionEntryPattern_PLAN_JOB_TEMPLATE_COMPOSITION_ENTRY_PATTERN_BUNDLE_ENTRY},
			}}, nil
		},
		ReadJobTemplate: func(_ context.Context, req *jobtemplatepb.ReadJobTemplateRequest) (*jobtemplatepb.ReadJobTemplateResponse, error) {
			return &jobtemplatepb.ReadJobTemplateResponse{Data: []*jobtemplatepb.JobTemplate{templates[req.GetData().GetId()]}}, nil
		},
	}

	got := detectSpawnJobs(context.Background(), deps, "price-plan-1")
	if !got.Available || got.JobCount != 2 {
		t.Fatalf("detection = %+v, want available with two jobs", got)
	}
	if want := []string{"Mathematics", "Science"}; !reflect.DeepEqual(got.TemplateNames, want) {
		t.Fatalf("template names = %v, want %v", got.TemplateNames, want)
	}
}

func TestDetectSpawnJobs_LegacyRootFallback(t *testing.T) {
	planID, rootID := "plan-legacy", "legacy-root"
	deps := &Deps{
		ReadPricePlan: func(context.Context, *priceplanpb.ReadPricePlanRequest) (*priceplanpb.ReadPricePlanResponse, error) {
			return &priceplanpb.ReadPricePlanResponse{Data: []*priceplanpb.PricePlan{{Id: "price-plan-legacy", PlanId: planID}}}, nil
		},
		ReadPlan: func(context.Context, *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error) {
			return &planpb.ReadPlanResponse{Data: []*planpb.Plan{{Id: &planID, JobTemplateId: &rootID, Active: true}}}, nil
		},
		ListPlanJobTemplatesByPlan: func(context.Context, *planjobtemplatepb.ListPlanJobTemplatesByPlanRequest) (*planjobtemplatepb.ListPlanJobTemplatesByPlanResponse, error) {
			return &planjobtemplatepb.ListPlanJobTemplatesByPlanResponse{}, nil
		},
		ReadJobTemplate: func(context.Context, *jobtemplatepb.ReadJobTemplateRequest) (*jobtemplatepb.ReadJobTemplateResponse, error) {
			return &jobtemplatepb.ReadJobTemplateResponse{Data: []*jobtemplatepb.JobTemplate{{Id: rootID, Name: "Legacy Root", Active: true}}}, nil
		},
	}

	got := detectSpawnJobs(context.Background(), deps, "price-plan-legacy")
	if !got.Available || got.JobCount != 1 || !reflect.DeepEqual(got.TemplateNames, []string{"Legacy Root"}) {
		t.Fatalf("detection = %+v, want legacy fallback", got)
	}
}
