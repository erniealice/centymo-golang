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

	consumer "github.com/erniealice/espyna-golang/consumer"

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
	db                          centymo.DataSource
	costScheduleRoutes          centymo.CostScheduleRoutes
	costScheduleLabels          centymo.CostScheduleLabels
	supplierPlanRoutes          centymo.SupplierPlanRoutes
	supplierPlanLabels          centymo.SupplierPlanLabels
	costPlanRoutes              centymo.CostPlanRoutes
	costPlanLabels              centymo.CostPlanLabels
	supplierProductPlanRoutes   centymo.SupplierProductPlanRoutes
	supplierProductPlanLabels   centymo.SupplierProductPlanLabels
	supplierProductCostPlanLabels centymo.SupplierProductCostPlanLabels
	supplierSubscriptionRoutes  centymo.SupplierSubscriptionRoutes
	supplierSubscriptionLabels  centymo.SupplierSubscriptionLabels
	centymoTableLabels          types.TableLabels
}

// wireSupplierSubscriptionModules lifts the bodies of the six P3 procurement
// `if cfg.wantXxx()` branches from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once at the position where
// the P3 supplier-subscription wiring used to be.
func wireSupplierSubscriptionModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *consumer.UseCases, w supplierSubscriptionWiring) {
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
		if useCases.Procurement != nil && useCases.Procurement.CostSchedule != nil {
			uc := useCases.Procurement.CostSchedule
			if uc.CreateCostSchedule != nil {
				csDeps.CreateCostSchedule = uc.CreateCostSchedule.Execute
			}
			if uc.ReadCostSchedule != nil {
				csDeps.ReadCostSchedule = uc.ReadCostSchedule.Execute
			}
			if uc.UpdateCostSchedule != nil {
				csDeps.UpdateCostSchedule = uc.UpdateCostSchedule.Execute
			}
			if uc.DeleteCostSchedule != nil {
				csDeps.DeleteCostSchedule = uc.DeleteCostSchedule.Execute
			}
			if uc.GetCostScheduleListPageData != nil {
				csDeps.GetCostScheduleListPageData = uc.GetCostScheduleListPageData.Execute
			}
			if uc.GetCostScheduleItemPageData != nil {
				csDeps.GetCostScheduleItemPageData = uc.GetCostScheduleItemPageData.Execute
			}
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
		if useCases.Procurement != nil && useCases.Procurement.SupplierPlan != nil {
			uc := useCases.Procurement.SupplierPlan
			if uc.CreateSupplierPlan != nil {
				spDeps.CreateSupplierPlan = uc.CreateSupplierPlan.Execute
			}
			if uc.ReadSupplierPlan != nil {
				spDeps.ReadSupplierPlan = uc.ReadSupplierPlan.Execute
			}
			if uc.UpdateSupplierPlan != nil {
				spDeps.UpdateSupplierPlan = uc.UpdateSupplierPlan.Execute
			}
			if uc.DeleteSupplierPlan != nil {
				spDeps.DeleteSupplierPlan = uc.DeleteSupplierPlan.Execute
			}
			if uc.GetSupplierPlanListPageData != nil {
				spDeps.GetSupplierPlanListPageData = uc.GetSupplierPlanListPageData.Execute
			}
			if uc.GetSupplierPlanItemPageData != nil {
				spDeps.GetSupplierPlanItemPageData = uc.GetSupplierPlanItemPageData.Execute
			}
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
		if useCases.Procurement != nil && useCases.Procurement.CostPlan != nil {
			uc := useCases.Procurement.CostPlan
			if uc.CreateCostPlan != nil {
				cpDeps.CreateCostPlan = uc.CreateCostPlan.Execute
			}
			if uc.ReadCostPlan != nil {
				cpDeps.ReadCostPlan = uc.ReadCostPlan.Execute
			}
			if uc.UpdateCostPlan != nil {
				cpDeps.UpdateCostPlan = uc.UpdateCostPlan.Execute
			}
			if uc.DeleteCostPlan != nil {
				cpDeps.DeleteCostPlan = uc.DeleteCostPlan.Execute
			}
			if uc.GetCostPlanListPageData != nil {
				cpDeps.GetCostPlanListPageData = uc.GetCostPlanListPageData.Execute
			}
			if uc.GetCostPlanItemPageData != nil {
				cpDeps.GetCostPlanItemPageData = uc.GetCostPlanItemPageData.Execute
			}
		}
		if useCases.Procurement != nil && useCases.Procurement.SupplierProductCostPlan != nil {
			uc := useCases.Procurement.SupplierProductCostPlan
			if uc.CreateSupplierProductCostPlan != nil {
				cpDeps.CreateSupplierProductCostPlan = uc.CreateSupplierProductCostPlan.Execute
			}
			if uc.ReadSupplierProductCostPlan != nil {
				cpDeps.ReadSupplierProductCostPlan = uc.ReadSupplierProductCostPlan.Execute
			}
			if uc.UpdateSupplierProductCostPlan != nil {
				cpDeps.UpdateSupplierProductCostPlan = uc.UpdateSupplierProductCostPlan.Execute
			}
			if uc.DeleteSupplierProductCostPlan != nil {
				cpDeps.DeleteSupplierProductCostPlan = uc.DeleteSupplierProductCostPlan.Execute
			}
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
		if useCases.Procurement != nil && useCases.Procurement.SupplierProductPlan != nil {
			uc := useCases.Procurement.SupplierProductPlan
			if uc.CreateSupplierProductPlan != nil {
				sppDeps.CreateSupplierProductPlan = uc.CreateSupplierProductPlan.Execute
			}
			if uc.ReadSupplierProductPlan != nil {
				sppDeps.ReadSupplierProductPlan = uc.ReadSupplierProductPlan.Execute
			}
			if uc.UpdateSupplierProductPlan != nil {
				sppDeps.UpdateSupplierProductPlan = uc.UpdateSupplierProductPlan.Execute
			}
			if uc.DeleteSupplierProductPlan != nil {
				sppDeps.DeleteSupplierProductPlan = uc.DeleteSupplierProductPlan.Execute
			}
			if uc.GetSupplierProductPlanListPageData != nil {
				sppDeps.GetSupplierProductPlanListPageData = uc.GetSupplierProductPlanListPageData.Execute
			}
			if uc.GetSupplierProductPlanItemPageData != nil {
				sppDeps.GetSupplierProductPlanItemPageData = uc.GetSupplierProductPlanItemPageData.Execute
			}
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
		if useCases.Procurement != nil && useCases.Procurement.SupplierProductCostPlan != nil {
			uc := useCases.Procurement.SupplierProductCostPlan
			if uc.CreateSupplierProductCostPlan != nil {
				spcpDeps.CreateSupplierProductCostPlan = uc.CreateSupplierProductCostPlan.Execute
			}
			if uc.ReadSupplierProductCostPlan != nil {
				spcpDeps.ReadSupplierProductCostPlan = uc.ReadSupplierProductCostPlan.Execute
			}
			if uc.UpdateSupplierProductCostPlan != nil {
				spcpDeps.UpdateSupplierProductCostPlan = uc.UpdateSupplierProductCostPlan.Execute
			}
			if uc.DeleteSupplierProductCostPlan != nil {
				spcpDeps.DeleteSupplierProductCostPlan = uc.DeleteSupplierProductCostPlan.Execute
			}
			if uc.GetSupplierProductCostPlanItemPageData != nil {
				spcpDeps.GetSupplierProductCostPlanItemPageData = uc.GetSupplierProductCostPlanItemPageData.Execute
			}
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
		if useCases.Procurement != nil && useCases.Procurement.SupplierSubscription != nil {
			uc := useCases.Procurement.SupplierSubscription
			if uc.CreateSupplierSubscription != nil {
				ssDeps.CreateSupplierSubscription = uc.CreateSupplierSubscription.Execute
			}
			if uc.ReadSupplierSubscription != nil {
				ssDeps.ReadSupplierSubscription = uc.ReadSupplierSubscription.Execute
			}
			if uc.UpdateSupplierSubscription != nil {
				ssDeps.UpdateSupplierSubscription = uc.UpdateSupplierSubscription.Execute
			}
			if uc.DeleteSupplierSubscription != nil {
				ssDeps.DeleteSupplierSubscription = uc.DeleteSupplierSubscription.Execute
			}
			if uc.GetSupplierSubscriptionListPageData != nil {
				ssDeps.GetSupplierSubscriptionListPageData = uc.GetSupplierSubscriptionListPageData.Execute
			}
			if uc.GetSupplierSubscriptionItemPageData != nil {
				ssDeps.GetSupplierSubscriptionItemPageData = uc.GetSupplierSubscriptionItemPageData.Execute
			}
		}
		suppliersubscriptionmod.NewModule(ssDeps).RegisterRoutes(ctx.Routes)
	}
}
