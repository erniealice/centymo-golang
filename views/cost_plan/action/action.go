package action

import (
	"context"
	"math"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/cost_plan/form"
	pyeza "github.com/erniealice/pyeza-golang"

	costplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_plan"
)

// Deps holds dependencies for cost_plan action handlers.
type Deps struct {
	Routes       centymo.CostPlanRoutes
	Labels       centymo.CostPlanLabels
	CommonLabels pyeza.CommonLabels

	CreateCostPlan          func(ctx context.Context, req *costplanpb.CreateCostPlanRequest) (*costplanpb.CreateCostPlanResponse, error)
	ReadCostPlan            func(ctx context.Context, req *costplanpb.ReadCostPlanRequest) (*costplanpb.ReadCostPlanResponse, error)
	UpdateCostPlan          func(ctx context.Context, req *costplanpb.UpdateCostPlanRequest) (*costplanpb.UpdateCostPlanResponse, error)
	DeleteCostPlan          func(ctx context.Context, req *costplanpb.DeleteCostPlanRequest) (*costplanpb.DeleteCostPlanResponse, error)
	GetCostPlanItemPageData func(ctx context.Context, req *costplanpb.GetCostPlanItemPageDataRequest) (*costplanpb.GetCostPlanItemPageDataResponse, error)

	// SetCostPlanActive performs a raw DB update to toggle active.
	SetCostPlanActive func(ctx context.Context, id string, active bool) error

	// Autocomplete URLs for supplier_plan and cost_schedule selects.
	SearchSupplierPlanURL string
	SearchCostScheduleURL string
}

// buildFormLabels converts centymo.CostPlanLabels into form.Labels.
func buildFormLabels(l centymo.CostPlanLabels) form.Labels {
	return form.Labels{
		SectionIdentification:       l.Form.SectionIdentification,
		SectionRelationships:        l.Form.SectionRelationships,
		SectionConfiguration:        l.Form.SectionConfiguration,
		SectionSchedule:             l.Form.SectionSchedule,
		SectionNotes:                l.Form.SectionNotes,
		Name:                        l.Form.Name,
		NamePlaceholder:             l.Form.NamePlaceholder,
		Description:                 l.Form.Description,
		DescPlaceholder:             l.Form.DescPlaceholder,
		SupplierPlan:                l.Form.SupplierPlan,
		SupplierPlanPlaceholder:     l.Form.SupplierPlanPlaceholder,
		CostSchedule:                l.Form.CostSchedule,
		CostSchedulePlaceholder:     l.Form.CostSchedulePlaceholder,
		BillingKind:                 l.Form.BillingKind,
		AmountBasis:                 l.Form.AmountBasis,
		Amount:                      l.Form.Amount,
		AmountPlaceholder:           l.Form.AmountPlaceholder,
		Currency:                    l.Form.Currency,
		CurrencyPlaceholder:         l.Form.CurrencyPlaceholder,
		BillingCycle:                l.Form.BillingCycle,
		BillingCyclePlaceholder:     l.Form.BillingCyclePlaceholder,
		DefaultTerm:                 l.Form.DefaultTerm,
		DefaultTermPlaceholder:      l.Form.DefaultTermPlaceholder,
		Active:                      l.Form.Active,
		BillingKindOneTime:          l.Form.BillingKindOneTime,
		BillingKindRecurring:        l.Form.BillingKindRecurring,
		BillingKindContract:         l.Form.BillingKindContract,
		BillingKindUsageBased:       l.Form.BillingKindUsageBased,
		BillingKindAdHoc:            l.Form.BillingKindAdHoc,
		AmountBasisPerCycle:         l.Form.AmountBasisPerCycle,
		AmountBasisTotalPackage:     l.Form.AmountBasisTotalPackage,
		AmountBasisDerivedFromLines: l.Form.AmountBasisDerivedFromLines,
		AmountBasisPerOccurrence:    l.Form.AmountBasisPerOccurrence,
		DurationUnitDay:             l.Form.DurationUnitDay,
		DurationUnitWeek:            l.Form.DurationUnitWeek,
		DurationUnitMonth:           l.Form.DurationUnitMonth,
		DurationUnitYear:            l.Form.DurationUnitYear,
	}
}

// strPtr returns a pointer to s.
func strPtr(s string) *string { return &s }

// parseAmount converts a display string to centavos (e.g. "10.00" → 1000).
func parseAmount(s string) int64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f * 100))
}

// formatAmount converts centavos int64 to display string (e.g. 1000 → "10.00").
func formatAmount(centavos int64) string {
	return strconv.FormatFloat(float64(centavos)/100.0, 'f', 2, 64)
}
