package detail

import (
	"context"
	"testing"

	jobpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job"
)

func TestBuildOperationsRows_OriginCompositionOrderWithoutParents(t *testing.T) {
	math, science := "math", "science"
	jobs := []*jobpb.Job{
		{Id: "job-science", Name: "Science", JobTemplateId: &science},
		{Id: "job-math", Name: "Math", JobTemplateId: &math},
	}
	rows := buildOperationsRows(context.Background(), &DetailViewDeps{}, jobs, map[string]int32{"math": 1, "science": 2})
	if len(rows) != 2 {
		t.Fatalf("rows=%d", len(rows))
	}
	if rows[0].JobTemplateID != "math" || rows[1].JobTemplateID != "science" {
		t.Fatalf("composition order not applied: %+v", rows)
	}
	if !rows[0].IsRoot || !rows[1].IsRoot {
		t.Fatalf("parentless composition entries must render as roots")
	}
}

func TestApplyJobsTabData_CompositionOrderForRootlessBundle(t *testing.T) {
	math, science := "math", "science"
	page := &PageData{}
	applyJobsTabData(context.Background(), &DetailViewDeps{}, page, []*jobpb.Job{
		{Id: "job-science", Name: "Science", JobTemplateId: &science},
		{Id: "job-math", Name: "Math", JobTemplateId: &math},
	}, map[string]int32{"math": 1, "science": 2})
	if page.Jobs == nil || len(page.Jobs.Rows) != 2 {
		t.Fatalf("rows=%v", page.Jobs)
	}
	if page.Jobs.Rows[0].JobTemplateID != "math" || page.Jobs.Rows[1].JobTemplateID != "science" {
		t.Fatalf("composition order not applied: %+v", page.Jobs.Rows)
	}
}

func TestApplyJobsTabData_LegacyRootRemainsFirst(t *testing.T) {
	rootTemplate, childTemplate, rootID := "root", "child", "job-root"
	page := &PageData{}
	applyJobsTabData(context.Background(), &DetailViewDeps{}, page, []*jobpb.Job{
		{Id: "job-child", Name: "Child", JobTemplateId: &childTemplate, ParentJobId: &rootID},
		{Id: rootID, Name: "Root", JobTemplateId: &rootTemplate},
	}, map[string]int32{"child": 1})
	if page.Jobs == nil || len(page.Jobs.Rows) != 2 {
		t.Fatalf("rows=%v", page.Jobs)
	}
	if page.Jobs.Rows[0].JobID != rootID {
		t.Fatalf("legacy root should remain first: %+v", page.Jobs.Rows)
	}
}
