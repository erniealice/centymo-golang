package action

import (
	"context"
	"math"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/supplier_product_cost_plan/form"
	pyeza "github.com/erniealice/pyeza-golang"

	supplierproductcostplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_cost_plan"
)

// Deps holds dependencies for supplier_product_cost_plan inline action handlers.
type Deps struct {
	// CostPlanRoutes for the parent entity (add/edit/delete URLs are scoped under /cost-plan/{id}/...).
	CostPlanRoutes centymo.CostPlanRoutes
	Labels         centymo.SupplierProductCostPlanLabels
	CommonLabels   pyeza.CommonLabels

	CreateSupplierProductCostPlan          func(ctx context.Context, req *supplierproductcostplanpb.CreateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.CreateSupplierProductCostPlanResponse, error)
	ReadSupplierProductCostPlan            func(ctx context.Context, req *supplierproductcostplanpb.ReadSupplierProductCostPlanRequest) (*supplierproductcostplanpb.ReadSupplierProductCostPlanResponse, error)
	UpdateSupplierProductCostPlan          func(ctx context.Context, req *supplierproductcostplanpb.UpdateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.UpdateSupplierProductCostPlanResponse, error)
	DeleteSupplierProductCostPlan          func(ctx context.Context, req *supplierproductcostplanpb.DeleteSupplierProductCostPlanRequest) (*supplierproductcostplanpb.DeleteSupplierProductCostPlanResponse, error)
	GetSupplierProductCostPlanItemPageData func(ctx context.Context, req *supplierproductcostplanpb.GetSupplierProductCostPlanItemPageDataRequest) (*supplierproductcostplanpb.GetSupplierProductCostPlanItemPageDataResponse, error)

	// SearchSupplierProductPlanURL for the autocomplete in the drawer.
	SearchSupplierProductPlanURL string
}

// buildFormLabels converts centymo.SupplierProductCostPlanLabels into form.Labels.
func buildFormLabels(l centymo.SupplierProductCostPlanLabels) form.Labels {
	return form.Labels{
		SectionIdentification:          l.Form.SectionIdentification,
		SectionRelationships:           l.Form.SectionRelationships,
		SectionConfiguration:           l.Form.SectionConfiguration,
		SectionSchedule:                l.Form.SectionSchedule,
		SectionNotes:                   l.Form.SectionNotes,
		SupplierProductPlan:            l.Form.SupplierProductPlan,
		SupplierProductPlanPlaceholder: l.Form.SupplierProductPlanPlaceholder,
		BillingTreatment:               l.Form.BillingTreatment,
		Amount:                         l.Form.Amount,
		AmountPlaceholder:              l.Form.AmountPlaceholder,
		Active:                         l.Form.Active,
		TreatmentRecurring:             l.Form.TreatmentRecurring,
		TreatmentOneTimeInitial:        l.Form.TreatmentOneTimeInitial,
		TreatmentUsageBased:            l.Form.TreatmentUsageBased,
		TreatmentMinimumCommitment:     l.Form.TreatmentMinimumCommitment,
	}
}

// parseAmount converts a display string to centavos.
func parseAmount(s string) int64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f * 100))
}

// formatAmount converts centavos int64 to display string.
func formatAmount(centavos int64) string {
	return strconv.FormatFloat(float64(centavos)/100.0, 'f', 2, 64)
}
