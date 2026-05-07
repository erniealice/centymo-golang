package action

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/centymo-golang/views/supplier_product_plan/form"
	pyeza "github.com/erniealice/pyeza-golang"

	supplierproductplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_plan"
)

// Deps holds dependencies for supplier_product_plan action handlers.
type Deps struct {
	Routes       centymo.SupplierProductPlanRoutes
	Labels       centymo.SupplierProductPlanLabels
	CommonLabels pyeza.CommonLabels

	CreateSupplierProductPlan          func(ctx context.Context, req *supplierproductplanpb.CreateSupplierProductPlanRequest) (*supplierproductplanpb.CreateSupplierProductPlanResponse, error)
	ReadSupplierProductPlan            func(ctx context.Context, req *supplierproductplanpb.ReadSupplierProductPlanRequest) (*supplierproductplanpb.ReadSupplierProductPlanResponse, error)
	UpdateSupplierProductPlan          func(ctx context.Context, req *supplierproductplanpb.UpdateSupplierProductPlanRequest) (*supplierproductplanpb.UpdateSupplierProductPlanResponse, error)
	DeleteSupplierProductPlan          func(ctx context.Context, req *supplierproductplanpb.DeleteSupplierProductPlanRequest) (*supplierproductplanpb.DeleteSupplierProductPlanResponse, error)
	GetSupplierProductPlanItemPageData func(ctx context.Context, req *supplierproductplanpb.GetSupplierProductPlanItemPageDataRequest) (*supplierproductplanpb.GetSupplierProductPlanItemPageDataResponse, error)

	// SetSupplierProductPlanActive performs a raw DB update to toggle active.
	SetSupplierProductPlanActive func(ctx context.Context, id string, active bool) error

	// Autocomplete URLs for supplier_plan and product selects.
	SearchSupplierPlanURL string
	SearchProductURL      string
}

// buildFormLabels converts centymo.SupplierProductPlanLabels into form.Labels.
func buildFormLabels(l centymo.SupplierProductPlanLabels) form.Labels {
	return form.Labels{
		SectionIdentification:     l.Form.SectionIdentification,
		SectionRelationships:      l.Form.SectionRelationships,
		SectionConfiguration:      l.Form.SectionConfiguration,
		SectionSchedule:           l.Form.SectionSchedule,
		SectionNotes:              l.Form.SectionNotes,
		SupplierPlan:              l.Form.SupplierPlan,
		SupplierPlanPlaceholder:   l.Form.SupplierPlanPlaceholder,
		Product:                   l.Form.Product,
		ProductPlaceholder:        l.Form.ProductPlaceholder,
		ProductVariant:            l.Form.ProductVariant,
		ProductVariantPlaceholder: l.Form.ProductVariantPlaceholder,
		SupplierSKU:               l.Form.SupplierSKU,
		SupplierSKUPlaceholder:    l.Form.SupplierSKUPlaceholder,
		SupplierUnit:              l.Form.SupplierUnit,
		SupplierUnitPlaceholder:   l.Form.SupplierUnitPlaceholder,
		Active:                    l.Form.Active,
	}
}

// strPtr returns a pointer to s.
func strPtr(s string) *string { return &s }
