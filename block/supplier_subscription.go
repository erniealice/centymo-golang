// Package block — supplier-subscription (P3) domain wiring.
//
// Holds wireSupplierSubscriptionModules (the lifted bodies of the six
// `if cfg.wantXxx()` branches for CostSchedule, SupplierPlan, CostPlan,
// SupplierProductPlan, SupplierProductCostPlan, and SupplierSubscription).
//
// P3 (20260506-supplier-subscriptions).
package block

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	expendituredomain "github.com/erniealice/centymo-golang/domain/expenditure"
	procurement "github.com/erniealice/centymo-golang/domain/procurement"
	costplanaction "github.com/erniealice/centymo-golang/domain/procurement/cost_plan/action"
)

// supplierSubscriptionWiring holds everything wireSupplierSubscriptionModules
// needs from the surrounding Block() scope. More than 6 fields → struct.
// Kept private; never re-exported.
type supplierSubscriptionWiring struct {
	costScheduleRoutes            procurement.CostScheduleRoutes
	costScheduleLabels            procurement.CostScheduleLabels
	supplierPlanRoutes            procurement.SupplierPlanRoutes
	supplierPlanLabels            procurement.SupplierPlanLabels
	costPlanRoutes                procurement.CostPlanRoutes
	costPlanLabels                procurement.CostPlanLabels
	supplierProductPlanRoutes     procurement.SupplierProductPlanRoutes
	supplierProductPlanLabels     procurement.SupplierProductPlanLabels
	supplierProductCostPlanLabels procurement.SupplierProductCostPlanLabels
	supplierSubscriptionRoutes    procurement.SupplierSubscriptionRoutes
	supplierSubscriptionLabels    procurement.SupplierSubscriptionLabels
	// expenseRecognitionRunLabels supplies the "Run Recognitions" CTA label
	// for the supplier_subscription detail page's Linked Recognitions tab.
	// Plan A 20260517-expense-run Surface C.
	expenseRecognitionRunLabels expendituredomain.ExpenseRecognitionRunLabels
	centymoTableLabels          types.TableLabels
}

// wireSupplierSubscriptionModules lifts the bodies of the six P3 procurement
// `if cfg.wantXxx()` branches from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once at the position where
// the P3 supplier-subscription wiring used to be.
func wireSupplierSubscriptionModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *UseCases, w supplierSubscriptionWiring) {
	// =====================================================================
	// P3 — CostSchedule module
	// =====================================================================
	if cfg.wantCostSchedule() {
		csDeps := &procurement.CostScheduleModuleDeps{
			Routes:                w.costScheduleRoutes,
			Labels:                w.costScheduleLabels,
			CommonLabels:          ctx.Common,
			TableLabels:           w.centymoTableLabels,
			SetCostScheduleActive: setActiveClosure(useCases, "cost_schedule"),
		}
		cs := useCases.Procurement.CostSchedule
		if cs.CreateCostSchedule != nil {
			csDeps.CreateCostSchedule = cs.CreateCostSchedule
		}
		if cs.ReadCostSchedule != nil {
			csDeps.ReadCostSchedule = cs.ReadCostSchedule
		}
		if cs.UpdateCostSchedule != nil {
			csDeps.UpdateCostSchedule = cs.UpdateCostSchedule
		}
		if cs.DeleteCostSchedule != nil {
			csDeps.DeleteCostSchedule = cs.DeleteCostSchedule
		}
		if cs.GetCostScheduleListPageData != nil {
			csDeps.GetCostScheduleListPageData = cs.GetCostScheduleListPageData
		}
		if cs.GetCostScheduleItemPageData != nil {
			csDeps.GetCostScheduleItemPageData = cs.GetCostScheduleItemPageData
		}
		procurement.NewCostScheduleModule(csDeps).RegisterRoutes(ctx.Routes)
	}

	// =====================================================================
	// P3 — SupplierPlan module
	// =====================================================================
	if cfg.wantSupplierPlan() {
		spDeps := &procurement.SupplierPlanModuleDeps{
			Routes:                w.supplierPlanRoutes,
			Labels:                w.supplierPlanLabels,
			CommonLabels:          ctx.Common,
			TableLabels:           w.centymoTableLabels,
			SetSupplierPlanActive: setActiveClosure(useCases, "supplier_plan"),
			SearchSupplierURL:     w.supplierPlanRoutes.SearchSupplierURL,
		}
		sp := useCases.Procurement.SupplierPlan
		if sp.CreateSupplierPlan != nil {
			spDeps.CreateSupplierPlan = sp.CreateSupplierPlan
		}
		if sp.ReadSupplierPlan != nil {
			spDeps.ReadSupplierPlan = sp.ReadSupplierPlan
		}
		if sp.UpdateSupplierPlan != nil {
			spDeps.UpdateSupplierPlan = sp.UpdateSupplierPlan
		}
		if sp.DeleteSupplierPlan != nil {
			spDeps.DeleteSupplierPlan = sp.DeleteSupplierPlan
		}
		if sp.GetSupplierPlanListPageData != nil {
			spDeps.GetSupplierPlanListPageData = sp.GetSupplierPlanListPageData
		}
		if sp.GetSupplierPlanItemPageData != nil {
			spDeps.GetSupplierPlanItemPageData = sp.GetSupplierPlanItemPageData
		}
		procurement.NewSupplierPlanModule(spDeps).RegisterRoutes(ctx.Routes)
	}

	// =====================================================================
	// P3 — CostPlan module (with inline SupplierProductCostPlan editor)
	// =====================================================================
	if cfg.wantCostPlan() {
		cpDeps := &procurement.CostPlanModuleDeps{
			Routes:                       w.costPlanRoutes,
			Labels:                       w.costPlanLabels,
			ProductCostLabels:            w.supplierProductCostPlanLabels,
			CommonLabels:                 ctx.Common,
			TableLabels:                  w.centymoTableLabels,
			SetCostPlanActive:            setActiveClosure(useCases, "cost_plan"),
			SearchSupplierPlanURL:        w.costPlanRoutes.SearchSupplierPlanURL,
			SearchCostScheduleURL:        w.costPlanRoutes.SearchCostScheduleURL,
			SearchSupplierProductPlanURL: w.costPlanRoutes.SearchSupplierProductPlanURL,
		}
		cp := useCases.Procurement.CostPlan
		if cp.CreateCostPlan != nil {
			cpDeps.CreateCostPlan = cp.CreateCostPlan
		}
		if cp.ReadCostPlan != nil {
			cpDeps.ReadCostPlan = cp.ReadCostPlan
		}
		if cp.UpdateCostPlan != nil {
			cpDeps.UpdateCostPlan = cp.UpdateCostPlan
		}
		if cp.DeleteCostPlan != nil {
			cpDeps.DeleteCostPlan = cp.DeleteCostPlan
		}
		if cp.GetCostPlanListPageData != nil {
			cpDeps.GetCostPlanListPageData = cp.GetCostPlanListPageData
		}
		if cp.GetCostPlanItemPageData != nil {
			cpDeps.GetCostPlanItemPageData = cp.GetCostPlanItemPageData
		}
		spcp := useCases.Procurement.SupplierProductCostPlan
		if spcp.CreateSupplierProductCostPlan != nil {
			cpDeps.CreateSupplierProductCostPlan = spcp.CreateSupplierProductCostPlan
		}
		if spcp.ReadSupplierProductCostPlan != nil {
			cpDeps.ReadSupplierProductCostPlan = spcp.ReadSupplierProductCostPlan
		}
		if spcp.UpdateSupplierProductCostPlan != nil {
			cpDeps.UpdateSupplierProductCostPlan = spcp.UpdateSupplierProductCostPlan
		}
		if spcp.DeleteSupplierProductCostPlan != nil {
			cpDeps.DeleteSupplierProductCostPlan = spcp.DeleteSupplierProductCostPlan
		}
		cpMod := procurement.NewCostPlanModule(cpDeps)
		cpMod.RegisterRoutes(ctx.Routes)
	}

	// =====================================================================
	// P3 — SupplierProductPlan module
	// =====================================================================
	if cfg.wantSupplierProductPlan() {
		sppDeps := &procurement.SupplierProductPlanModuleDeps{
			Routes:                       w.supplierProductPlanRoutes,
			Labels:                       w.supplierProductPlanLabels,
			CommonLabels:                 ctx.Common,
			TableLabels:                  w.centymoTableLabels,
			SetSupplierProductPlanActive: setActiveClosure(useCases, "supplier_product_plan"),
			SearchSupplierPlanURL:        w.supplierProductPlanRoutes.SearchSupplierPlanURL,
			SearchProductURL:             w.supplierProductPlanRoutes.SearchProductURL,
		}
		spp := useCases.Procurement.SupplierProductPlan
		if spp.CreateSupplierProductPlan != nil {
			sppDeps.CreateSupplierProductPlan = spp.CreateSupplierProductPlan
		}
		if spp.ReadSupplierProductPlan != nil {
			sppDeps.ReadSupplierProductPlan = spp.ReadSupplierProductPlan
		}
		if spp.UpdateSupplierProductPlan != nil {
			sppDeps.UpdateSupplierProductPlan = spp.UpdateSupplierProductPlan
		}
		if spp.DeleteSupplierProductPlan != nil {
			sppDeps.DeleteSupplierProductPlan = spp.DeleteSupplierProductPlan
		}
		if spp.GetSupplierProductPlanListPageData != nil {
			sppDeps.GetSupplierProductPlanListPageData = spp.GetSupplierProductPlanListPageData
		}
		if spp.GetSupplierProductPlanItemPageData != nil {
			sppDeps.GetSupplierProductPlanItemPageData = spp.GetSupplierProductPlanItemPageData
		}
		procurement.NewSupplierProductPlanModule(sppDeps).RegisterRoutes(ctx.Routes)
	}

	// =====================================================================
	// P3 — SupplierProductCostPlan inline module (standalone routes; also
	// mounted inside CostPlan via cpMod above when both wantCostPlan and
	// wantSupplierProductCostPlan are active).
	// =====================================================================
	if cfg.wantSupplierProductCostPlan() && !cfg.wantCostPlan() {
		// Only register standalone SPCP routes when CostPlan module is NOT
		// already registered (which would mount the same URLs twice).
		// Action handlers now live in cost_plan/action (templates-only refactor
		// per docs/plan/20260509-buying-selling-parity-audit plan.md §5.4 θ).
		spcpDeps := &costplanaction.CostPlanLineDeps{
			CostPlanRoutes:               w.costPlanRoutes,
			Labels:                       w.supplierProductCostPlanLabels,
			CommonLabels:                 ctx.Common,
			SearchSupplierProductPlanURL: w.costPlanRoutes.SearchSupplierProductPlanURL,
		}
		spcp := useCases.Procurement.SupplierProductCostPlan
		if spcp.CreateSupplierProductCostPlan != nil {
			spcpDeps.CreateSupplierProductCostPlan = spcp.CreateSupplierProductCostPlan
		}
		if spcp.ReadSupplierProductCostPlan != nil {
			spcpDeps.ReadSupplierProductCostPlan = spcp.ReadSupplierProductCostPlan
		}
		if spcp.UpdateSupplierProductCostPlan != nil {
			spcpDeps.UpdateSupplierProductCostPlan = spcp.UpdateSupplierProductCostPlan
		}
		if spcp.DeleteSupplierProductCostPlan != nil {
			spcpDeps.DeleteSupplierProductCostPlan = spcp.DeleteSupplierProductCostPlan
		}
		if spcp.GetSupplierProductCostPlanItemPageData != nil {
			spcpDeps.GetSupplierProductCostPlanItemPageData = spcp.GetSupplierProductCostPlanItemPageData
		}
		if w.costPlanRoutes.ProductCostAddURL != "" {
			addView := costplanaction.NewCostPlanLineAddAction(spcpDeps)
			ctx.Routes.GET(w.costPlanRoutes.ProductCostAddURL, addView)
			ctx.Routes.POST(w.costPlanRoutes.ProductCostAddURL, addView)
		}
		if w.costPlanRoutes.ProductCostEditURL != "" {
			editView := costplanaction.NewCostPlanLineEditAction(spcpDeps)
			ctx.Routes.GET(w.costPlanRoutes.ProductCostEditURL, editView)
			ctx.Routes.POST(w.costPlanRoutes.ProductCostEditURL, editView)
		}
		if w.costPlanRoutes.ProductCostDeleteURL != "" {
			ctx.Routes.POST(w.costPlanRoutes.ProductCostDeleteURL, costplanaction.NewCostPlanLineDeleteAction(spcpDeps))
		}
	}

	// =====================================================================
	// P3 — SupplierSubscription module
	// =====================================================================
	if cfg.wantSupplierSubscription() {
		ssDeps := &procurement.SupplierSubscriptionModuleDeps{
			Routes:                        w.supplierSubscriptionRoutes,
			Labels:                        w.supplierSubscriptionLabels,
			ExpenseRecognitionRunLabels:   w.expenseRecognitionRunLabels,
			CommonLabels:                  ctx.Common,
			TableLabels:                   w.centymoTableLabels,
			SetSupplierSubscriptionActive: setActiveClosure(useCases, "supplier_subscription"),
		}
		ss := useCases.Procurement.SupplierSubscription
		if ss.CreateSupplierSubscription != nil {
			ssDeps.CreateSupplierSubscription = ss.CreateSupplierSubscription
		}
		if ss.ReadSupplierSubscription != nil {
			ssDeps.ReadSupplierSubscription = ss.ReadSupplierSubscription
		}
		if ss.UpdateSupplierSubscription != nil {
			ssDeps.UpdateSupplierSubscription = ss.UpdateSupplierSubscription
		}
		if ss.DeleteSupplierSubscription != nil {
			ssDeps.DeleteSupplierSubscription = ss.DeleteSupplierSubscription
		}
		if ss.GetSupplierSubscriptionListPageData != nil {
			ssDeps.GetSupplierSubscriptionListPageData = ss.GetSupplierSubscriptionListPageData
		}
		if ss.GetSupplierSubscriptionItemPageData != nil {
			ssDeps.GetSupplierSubscriptionItemPageData = ss.GetSupplierSubscriptionItemPageData
		}
		// Plan A Surface C — ReadCostPlan for the CTA helper.
		if cp := useCases.Procurement.CostPlan; cp.ReadCostPlan != nil {
			ssDeps.ReadCostPlan = cp.ReadCostPlan
		}
		procurement.NewSupplierSubscriptionModule(ssDeps).RegisterRoutes(ctx.Routes)
	}
}
