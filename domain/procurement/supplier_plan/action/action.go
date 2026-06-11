package action

import (
	"context"

	"github.com/erniealice/centymo-golang/domain/procurement/supplier_plan"
	"github.com/erniealice/centymo-golang/domain/procurement/supplier_plan/form"
	pyeza "github.com/erniealice/pyeza-golang"

	supplierplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_plan"
)

// Deps holds dependencies for supplier_plan action handlers.
type Deps struct {
	Routes       supplier_plan.Routes
	Labels       supplier_plan.Labels
	CommonLabels pyeza.CommonLabels

	CreateSupplierPlan          func(ctx context.Context, req *supplierplanpb.CreateSupplierPlanRequest) (*supplierplanpb.CreateSupplierPlanResponse, error)
	ReadSupplierPlan            func(ctx context.Context, req *supplierplanpb.ReadSupplierPlanRequest) (*supplierplanpb.ReadSupplierPlanResponse, error)
	UpdateSupplierPlan          func(ctx context.Context, req *supplierplanpb.UpdateSupplierPlanRequest) (*supplierplanpb.UpdateSupplierPlanResponse, error)
	DeleteSupplierPlan          func(ctx context.Context, req *supplierplanpb.DeleteSupplierPlanRequest) (*supplierplanpb.DeleteSupplierPlanResponse, error)
	GetSupplierPlanItemPageData func(ctx context.Context, req *supplierplanpb.GetSupplierPlanItemPageDataRequest) (*supplierplanpb.GetSupplierPlanItemPageDataResponse, error)

	// SetSupplierPlanActive performs a raw DB update to toggle active.
	SetSupplierPlanActive func(ctx context.Context, id string, active bool) error

	// SearchSupplierURL is used by the form autocomplete for supplier selection.
	SearchSupplierURL string
}

// buildFormLabels converts supplier_plan.Labels into form.Labels.
func buildFormLabels(l supplier_plan.Labels) form.Labels {
	return form.Labels{
		SectionIdentification: l.Form.SectionIdentification,
		SectionRelationships:  l.Form.SectionRelationships,
		SectionConfiguration:  l.Form.SectionConfiguration,
		SectionSchedule:       l.Form.SectionSchedule,
		SectionNotes:          l.Form.SectionNotes,
		Name:                  l.Form.Name,
		NamePlaceholder:       l.Form.NamePlaceholder,
		Supplier:              l.Form.Supplier,
		SupplierPlaceholder:   l.Form.SupplierPlaceholder,
		Active:                l.Form.Active,
	}
}
