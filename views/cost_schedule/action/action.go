package action

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/cost_schedule/form"
	pyeza "github.com/erniealice/pyeza-golang"

	costschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_schedule"
)

// Deps holds dependencies for cost_schedule action handlers.
type Deps struct {
	Routes       centymo.CostScheduleRoutes
	Labels       centymo.CostScheduleLabels
	CommonLabels pyeza.CommonLabels

	CreateCostSchedule          func(ctx context.Context, req *costschedulepb.CreateCostScheduleRequest) (*costschedulepb.CreateCostScheduleResponse, error)
	ReadCostSchedule            func(ctx context.Context, req *costschedulepb.ReadCostScheduleRequest) (*costschedulepb.ReadCostScheduleResponse, error)
	UpdateCostSchedule          func(ctx context.Context, req *costschedulepb.UpdateCostScheduleRequest) (*costschedulepb.UpdateCostScheduleResponse, error)
	DeleteCostSchedule          func(ctx context.Context, req *costschedulepb.DeleteCostScheduleRequest) (*costschedulepb.DeleteCostScheduleResponse, error)
	GetCostScheduleItemPageData func(ctx context.Context, req *costschedulepb.GetCostScheduleItemPageDataRequest) (*costschedulepb.GetCostScheduleItemPageDataResponse, error)

	// SetCostScheduleActive performs a raw DB update to toggle active field.
	SetCostScheduleActive func(ctx context.Context, id string, active bool) error
}

// buildFormLabels converts centymo.CostScheduleLabels into form.Labels.
func buildFormLabels(l centymo.CostScheduleLabels) form.Labels {
	return form.Labels{
		SectionIdentification: l.Form.SectionIdentification,
		SectionRelationships:  l.Form.SectionRelationships,
		SectionConfiguration:  l.Form.SectionConfiguration,
		SectionSchedule:       l.Form.SectionSchedule,
		SectionNotes:          l.Form.SectionNotes,
		Name:                  l.Form.Name,
		NamePlaceholder:       l.Form.NamePlaceholder,
		Description:           l.Form.Description,
		DescPlaceholder:       l.Form.DescPlaceholder,
		StartDate:             l.Form.StartDate,
		EndDate:               l.Form.EndDate,
		Location:              l.Form.Location,
		LocationPlaceholder:   l.Form.LocationPlaceholder,
		Active:                l.Form.Active,
	}
}
