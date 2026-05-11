// Package block — supplier-subscription (P3) domain wiring.
//
// Holds wireSupplierSubscriptionModules (the lifted bodies of the six
// `if cfg.wantXxx()` branches for CostSchedule, SupplierPlan, CostPlan,
// SupplierProductPlan, SupplierProductCostPlan, and SupplierSubscription).
//
// P3 (20260506-supplier-subscriptions).
package block

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	costplanmod "github.com/erniealice/centymo-golang/views/cost_plan"
	costplanaction "github.com/erniealice/centymo-golang/views/cost_plan/action"
	costschedulemod "github.com/erniealice/centymo-golang/views/cost_schedule"
	supplierplanmod "github.com/erniealice/centymo-golang/views/supplier_plan"
	supplierproductplanmod "github.com/erniealice/centymo-golang/views/supplier_product_plan"
	suppliersubscriptionmod "github.com/erniealice/centymo-golang/views/supplier_subscription"
)

// supplierSubscriptionWiring holds everything wireSupplierSubscriptionModules
// needs from the surrounding Block() scope. More than 6 fields → struct.
// Kept private; never re-exported.
type supplierSubscriptionWiring struct {
	db                            centymo.DataSource
	costScheduleRoutes            centymo.CostScheduleRoutes
	costScheduleLabels            centymo.CostScheduleLabels
	supplierPlanRoutes            centymo.SupplierPlanRoutes
	supplierPlanLabels            centymo.SupplierPlanLabels
	costPlanRoutes                centymo.CostPlanRoutes
	costPlanLabels                centymo.CostPlanLabels
	supplierProductPlanRoutes     centymo.SupplierProductPlanRoutes
	supplierProductPlanLabels     centymo.SupplierProductPlanLabels
	supplierProductCostPlanLabels centymo.SupplierProductCostPlanLabels
	supplierSubscriptionRoutes    centymo.SupplierSubscriptionRoutes
	supplierSubscriptionLabels    centymo.SupplierSubscriptionLabels
	centymoTableLabels            types.TableLabels
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
		csDeps := &costschedulemod.ModuleDeps{
			Routes:       w.costScheduleRoutes,
			Labels:       w.costScheduleLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
			SetCostScheduleActive: func(fctx context.Context, id string, active bool) error {
				_, err := w.db.Update(fctx, "cost_schedule", id, map[string]any{"active": active})
				return err
			},
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
		costschedulemod.NewModule(csDeps).RegisterRoutes(ctx.Routes)
	}

	// =====================================================================
	// P3 — SupplierPlan module
	// =====================================================================
	if cfg.wantSupplierPlan() {
		spDeps := &supplierplanmod.ModuleDeps{
			Routes:       w.supplierPlanRoutes,
			Labels:       w.supplierPlanLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
			SetSupplierPlanActive: func(fctx context.Context, id string, active bool) error {
				_, err := w.db.Update(fctx, "supplier_plan", id, map[string]any{"active": active})
				return err
			},
			SearchSupplierURL: w.supplierPlanRoutes.SearchSupplierURL,
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
		supplierplanmod.NewModule(spDeps).RegisterRoutes(ctx.Routes)
	}

	// =====================================================================
	// P3 — CostPlan module (with inline SupplierProductCostPlan editor)
	// =====================================================================
	if cfg.wantCostPlan() {
		cpDeps := &costplanmod.ModuleDeps{
			Routes:            w.costPlanRoutes,
			Labels:            w.costPlanLabels,
			ProductCostLabels: w.supplierProductCostPlanLabels,
			CommonLabels:      ctx.Common,
			TableLabels:       w.centymoTableLabels,
			SetCostPlanActive: func(fctx context.Context, id string, active bool) error {
				_, err := w.db.Update(fctx, "cost_plan", id, map[string]any{"active": active})
				return err
			},
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
		cpMod := costplanmod.NewModule(cpDeps)
		cpMod.RegisterRoutes(ctx.Routes)
	}

	// =====================================================================
	// P3 — SupplierProductPlan module
	// =====================================================================
	if cfg.wantSupplierProductPlan() {
		sppDeps := &supplierproductplanmod.ModuleDeps{
			Routes:       w.supplierProductPlanRoutes,
			Labels:       w.supplierProductPlanLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
			SetSupplierProductPlanActive: func(fctx context.Context, id string, active bool) error {
				_, err := w.db.Update(fctx, "supplier_product_plan", id, map[string]any{"active": active})
				return err
			},
			SearchSupplierPlanURL: w.supplierProductPlanRoutes.SearchSupplierPlanURL,
			SearchProductURL:      w.supplierProductPlanRoutes.SearchProductURL,
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
		supplierproductplanmod.NewModule(sppDeps).RegisterRoutes(ctx.Routes)
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
		ssDeps := &suppliersubscriptionmod.ModuleDeps{
			Routes:       w.supplierSubscriptionRoutes,
			Labels:       w.supplierSubscriptionLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
			SetSupplierSubscriptionActive: func(fctx context.Context, id string, active bool) error {
				_, err := w.db.Update(fctx, "supplier_subscription", id, map[string]any{"active": active})
				return err
			},
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
		suppliersubscriptionmod.NewModule(ssDeps).RegisterRoutes(ctx.Routes)
	}
}
