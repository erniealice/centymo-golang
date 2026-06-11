package action

// cost_plan_line.go — action handlers for SupplierProductCostPlan inline
// editor (Add / Edit / Delete). These handlers are mounted under cost-plan
// parent routes so that the Lines tab on a CostPlan detail page can create,
// update, and remove SupplierProductCostPlan rows without a standalone module.
//
// Template rendered: supplier-product-cost-plan-drawer-form
// (from views/supplier_product_cost_plan/templates/ — FS registered in
// apps/service-admin/internal/composition/container.go).

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/centymo-golang/domain/procurement/cost_plan"
	costplanlineform "github.com/erniealice/centymo-golang/domain/procurement/cost_plan/form"
	sib_procurement_supplier_product_cost_plan "github.com/erniealice/centymo-golang/domain/procurement/supplier_product_cost_plan"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	supplierproductcostplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_product_cost_plan"
)

// CostPlanLineDeps holds dependencies for SupplierProductCostPlan inline
// action handlers (scoped under a CostPlan parent route).
type CostPlanLineDeps struct {
	// CostPlanRoutes for the parent entity — add/edit/delete URLs are
	// scoped under /cost-plan/{id}/... via ProductCostAddURL, ProductCostEditURL,
	// ProductCostDeleteURL.
	CostPlanRoutes cost_plan.Routes
	Labels         sib_procurement_supplier_product_cost_plan.Labels
	CommonLabels   pyeza.CommonLabels

	CreateSupplierProductCostPlan          func(ctx context.Context, req *supplierproductcostplanpb.CreateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.CreateSupplierProductCostPlanResponse, error)
	ReadSupplierProductCostPlan            func(ctx context.Context, req *supplierproductcostplanpb.ReadSupplierProductCostPlanRequest) (*supplierproductcostplanpb.ReadSupplierProductCostPlanResponse, error)
	UpdateSupplierProductCostPlan          func(ctx context.Context, req *supplierproductcostplanpb.UpdateSupplierProductCostPlanRequest) (*supplierproductcostplanpb.UpdateSupplierProductCostPlanResponse, error)
	DeleteSupplierProductCostPlan          func(ctx context.Context, req *supplierproductcostplanpb.DeleteSupplierProductCostPlanRequest) (*supplierproductcostplanpb.DeleteSupplierProductCostPlanResponse, error)
	GetSupplierProductCostPlanItemPageData func(ctx context.Context, req *supplierproductcostplanpb.GetSupplierProductCostPlanItemPageDataRequest) (*supplierproductcostplanpb.GetSupplierProductCostPlanItemPageDataResponse, error)

	// SearchSupplierProductPlanURL for the autocomplete in the drawer.
	SearchSupplierProductPlanURL string
}

// buildCostPlanLineFormLabels converts sib_procurement_supplier_product_cost_plan.Labels
// into costplanlineform.CostPlanLineLabels.
func buildCostPlanLineFormLabels(l sib_procurement_supplier_product_cost_plan.Labels) costplanlineform.CostPlanLineLabels {
	return costplanlineform.CostPlanLineLabels{
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

// NewCostPlanLineAddAction creates the SupplierProductCostPlan add action
// (inline within CostPlan detail).
func NewCostPlanLineAddAction(deps *CostPlanLineDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_cost_plan", "create") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		costPlanID := viewCtx.Request.PathValue("id")
		addURL := route.ResolveURL(deps.CostPlanRoutes.ProductCostAddURL, "id", costPlanID)
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("supplier-product-cost-plan-drawer-form", &costplanlineform.CostPlanLineData{
				FormAction:                   addURL,
				CostPlanID:                   costPlanID,
				Active:                       true,
				SearchSupplierProductPlanURL: deps.SearchSupplierProductPlanURL,
				Labels:                       buildCostPlanLineFormLabels(deps.Labels),
				CommonLabels:                 deps.CommonLabels,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		supplierProductPlanID := r.FormValue("supplier_product_plan_id")
		billingTreatment := r.FormValue("billing_treatment")
		amount := parseAmount(r.FormValue("amount"))
		active := r.FormValue("active") != "false"

		spcp := &supplierproductcostplanpb.SupplierProductCostPlan{
			CostPlanId:            costPlanID,
			SupplierProductPlanId: supplierProductPlanID,
			BillingAmount:         amount,
			Active:                active,
		}
		if billingTreatment != "" {
			if bt, ok := supplierproductcostplanpb.SupplierProductCostPlanBillingTreatment_value[billingTreatment]; ok {
				spcp.BillingTreatment = supplierproductcostplanpb.SupplierProductCostPlanBillingTreatment(bt)
			}
		}

		if _, err := deps.CreateSupplierProductCostPlan(ctx, &supplierproductcostplanpb.CreateSupplierProductCostPlanRequest{Data: spcp}); err != nil {
			log.Printf("Failed to create supplier product cost plan: %v", err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("cost-plan-lines-table")
	})
}

// NewCostPlanLineEditAction creates the SupplierProductCostPlan edit action
// (inline within CostPlan detail).
func NewCostPlanLineEditAction(deps *CostPlanLineDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_cost_plan", "update") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		costPlanID := viewCtx.Request.PathValue("id")
		pcid := viewCtx.Request.PathValue("pcid")
		editURL := route.ResolveURL(
			route.ResolveURL(deps.CostPlanRoutes.ProductCostEditURL, "id", costPlanID),
			"pcid", pcid,
		)
		if viewCtx.Request.Method == http.MethodGet {
			var record *supplierproductcostplanpb.SupplierProductCostPlan
			if deps.GetSupplierProductCostPlanItemPageData != nil {
				resp, err := deps.GetSupplierProductCostPlanItemPageData(ctx, &supplierproductcostplanpb.GetSupplierProductCostPlanItemPageDataRequest{
					SupplierProductCostPlanId: pcid,
				})
				if err != nil || resp == nil || resp.GetSupplierProductCostPlan() == nil {
					return view.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetSupplierProductCostPlan()
			} else {
				resp, err := deps.ReadSupplierProductCostPlan(ctx, &supplierproductcostplanpb.ReadSupplierProductCostPlanRequest{
					Data: &supplierproductcostplanpb.SupplierProductCostPlan{Id: pcid},
				})
				if err != nil || len(resp.GetData()) == 0 {
					return view.HTMXError(deps.Labels.Errors.NotFound)
				}
				record = resp.GetData()[0]
			}

			sppLabel := record.GetSupplierProductPlanId()
			if spp := record.GetSupplierProductPlan(); spp != nil && spp.GetName() != "" {
				sppLabel = spp.GetName()
			}

			return view.OK("supplier-product-cost-plan-drawer-form", &costplanlineform.CostPlanLineData{
				FormAction:                   editURL,
				IsEdit:                       true,
				ID:                           pcid,
				CostPlanID:                   costPlanID,
				SupplierProductPlanID:        record.GetSupplierProductPlanId(),
				SupplierProductPlanLabel:     sppLabel,
				BillingTreatment:             record.GetBillingTreatment().String(),
				Amount:                       formatAmount(record.GetBillingAmount()),
				Active:                       record.GetActive(),
				SearchSupplierProductPlanURL: deps.SearchSupplierProductPlanURL,
				Labels:                       buildCostPlanLineFormLabels(deps.Labels),
				CommonLabels:                 deps.CommonLabels,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.InvalidFormData)
		}
		r := viewCtx.Request
		supplierProductPlanID := r.FormValue("supplier_product_plan_id")
		billingTreatment := r.FormValue("billing_treatment")
		amount := parseAmount(r.FormValue("amount"))
		active := r.FormValue("active") != "false"

		spcp := &supplierproductcostplanpb.SupplierProductCostPlan{
			Id:                    pcid,
			CostPlanId:            costPlanID,
			SupplierProductPlanId: supplierProductPlanID,
			BillingAmount:         amount,
			Active:                active,
		}
		if billingTreatment != "" {
			if bt, ok := supplierproductcostplanpb.SupplierProductCostPlanBillingTreatment_value[billingTreatment]; ok {
				spcp.BillingTreatment = supplierproductcostplanpb.SupplierProductCostPlanBillingTreatment(bt)
			}
		}

		if _, err := deps.UpdateSupplierProductCostPlan(ctx, &supplierproductcostplanpb.UpdateSupplierProductCostPlanRequest{Data: spcp}); err != nil {
			log.Printf("Failed to update supplier product cost plan %s: %v", pcid, err)
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("cost-plan-lines-table")
	})
}

// NewCostPlanLineDeleteAction creates the SupplierProductCostPlan delete
// action (inline within CostPlan detail).
func NewCostPlanLineDeleteAction(deps *CostPlanLineDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("supplier_product_cost_plan", "delete") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		if _, err := deps.DeleteSupplierProductCostPlan(ctx, &supplierproductcostplanpb.DeleteSupplierProductCostPlanRequest{
			Data: &supplierproductcostplanpb.SupplierProductCostPlan{Id: id},
		}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess("cost-plan-lines-table")
	})
}
