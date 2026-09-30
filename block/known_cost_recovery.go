package block

// known_cost_recovery.go — the opt-in "known-cost recovery" surface of the
// usage-and-pass-through plan (Slice B S1): the expenditure "Recoverable costs"
// tab (cost lines + allocation) and the subscription read-only "Charge terms" tab.
//
// The surface is off by default so apps that do not opt in keep byte-identical
// compositions (leasing-admin passes WithKnownCostRecovery(true); service-admin and
// school-admin never do). It follows the WithChargePolicies precedent (fycha
// block/engine_options.go) and the WithProductAssets precedent (product_assets.go):
// an EngineOption + a nil-safe bind of proto closures (`.Execute`, C13) + units
// appended by AllUnits only when the option is on.
//
// R4: every Unit here validates its closures in Mount and returns an error when a
// required one is nil (a Mount error stops boot). Units that exist only on the compose-v2
// path do not extend RequireFor (block anatomy policy).

import (
	"fmt"

	"github.com/erniealice/espyna-golang/consumer/compose"

	expendituredomain "github.com/erniealice/centymo-golang/domain/expenditure"
	abpkg "github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch"
	cscpkg "github.com/erniealice/centymo-golang/domain/expenditure/cost_source_component"
	altpkg "github.com/erniealice/centymo-golang/domain/subscription/agreement_line_term"
)

// WithKnownCostRecovery mounts the known-cost recovery surface (recoverable cost
// lines, cost allocation drawers, the expenditure Recoverable costs tab and the
// subscription Charge terms tab). Off by default.
func WithKnownCostRecovery(enabled bool) EngineOption {
	return func(c *engineConfig) { c.knownCostRecovery = enabled }
}

// ---------------------------------------------------------------------------
// R4 — completeness checks run by the units' Mount
// ---------------------------------------------------------------------------

// requireCostSourceComponent validates the recoverable cost line closures.
func requireCostSourceComponent(u *UseCases) error {
	if u == nil {
		return fmt.Errorf("centymo cost_source_component: UseCases is nil")
	}
	c := &u.CostSourceComponent
	return requireUnit("cost_source_component", map[string]bool{
		"UseCases.CostSourceComponent.CreateCostSourceComponent": c.CreateCostSourceComponent != nil,
		"UseCases.CostSourceComponent.ReadCostSourceComponent":   c.ReadCostSourceComponent != nil,
		"UseCases.CostSourceComponent.UpdateCostSourceComponent": c.UpdateCostSourceComponent != nil,
		"UseCases.CostSourceComponent.DeleteCostSourceComponent": c.DeleteCostSourceComponent != nil,
		"UseCases.CostSourceComponent.ListCostSourceComponents":  c.ListCostSourceComponents != nil,
		"UseCases.Expenditure.ReadExpenditure":                   u.Expenditure.ReadExpenditure != nil,
	})
}

// requireAllocationBatch validates the cost allocation closures.
func requireAllocationBatch(u *UseCases) error {
	if u == nil {
		return fmt.Errorf("centymo allocation_batch: UseCases is nil")
	}
	a := &u.AllocationBatch
	return requireUnit("allocation_batch", map[string]bool{
		"UseCases.AllocationBatch.CreateAllocationBatch":          a.CreateAllocationBatch != nil,
		"UseCases.AllocationBatch.UpdateAllocationBatchShares":    a.UpdateAllocationBatchShares != nil,
		"UseCases.AllocationBatch.PublishAllocationBatch":         a.PublishAllocationBatch != nil,
		"UseCases.AllocationBatch.GetAllocationBatchListPageData": a.GetAllocationBatchListPageData != nil,
		"UseCases.AllocationBatch.GetAllocationShareListPageData": a.GetAllocationShareListPageData != nil,
		"UseCases.AllocationBatch.ReadCostSourceComponent":        a.ReadCostSourceComponent != nil,
		"UseCases.AllocationBatch.ListAgreementLineTerms":         a.ListAgreementLineTerms != nil,
		"UseCases.AllocationBatch.GetSubscriptionItemPageData":    a.GetSubscriptionItemPageData != nil,
	})
}

// requireAgreementLineTerm validates the subscription Charge terms tab closures (the
// name enrichment reads are optional and degrade to ids).
func requireAgreementLineTerm(u *UseCases) error {
	if u == nil {
		return fmt.Errorf("centymo agreement_line_term: UseCases is nil")
	}
	return requireUnit("agreement_line_term", map[string]bool{
		"UseCases.AgreementLineTerm.ListAgreementLineTerms": u.AgreementLineTerm.ListAgreementLineTerms != nil,
	})
}

// ---------------------------------------------------------------------------
// Units (appended by AllUnits only when WithKnownCostRecovery(true))
// ---------------------------------------------------------------------------

// CostSourceComponentUnit wires the recoverable cost line add / edit / delete drawers
// and the table refresh route.
func CostSourceComponentUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := cscpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		if err := requireCostSourceComponent(uc); err != nil {
			return err
		}
		r := u.Routes.(*cscpkg.Routes)
		l := u.Labels.(*cscpkg.Labels)
		allocationRoutes := abpkg.DefaultRoutes()
		if ar, ok := compose.RoutesOf[*abpkg.Routes](mc, "expenditure.allocation_batch"); ok {
			allocationRoutes = *ar
		}
		deps := &expendituredomain.CostSourceComponentModuleDeps{
			Routes: *r, Labels: *l, CommonLabels: mc.Common, TableLabels: mc.Table,
			UseCases: costSourceComponentViewUseCases(uc), AllocationRoutes: allocationRoutes,
			ReadExpenditure: uc.Expenditure.ReadExpenditure,
		}
		expendituredomain.NewCostSourceComponentModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// AllocationBatchUnit wires the cost allocation drawers (allocate, preview, publish, view).
func AllocationBatchUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := abpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		if err := requireAllocationBatch(uc); err != nil {
			return err
		}
		r := u.Routes.(*abpkg.Routes)
		l := u.Labels.(*abpkg.Labels)
		componentLabels := cscpkg.DefaultLabels()
		if cl, ok := compose.LabelsOf[*cscpkg.Labels](mc, "expenditure.cost_source_component"); ok {
			componentLabels = *cl
		}
		expendituredomain.NewAllocationBatchModule(&expendituredomain.AllocationBatchModuleDeps{
			Routes: *r, Labels: *l, ComponentLabels: componentLabels, CommonLabels: mc.Common,
			UseCases: allocationBatchViewUseCases(uc),
		}).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// AgreementLineTermUnit is the data-only unit that carries the Charge terms labels
// (post-overlay) for the subscription unit. Its Mount only validates the closures the
// tab needs.
func AgreementLineTermUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := altpkg.Describe()
	u.Mount = func(*compose.MountContext) error { return requireAgreementLineTerm(uc) }
	return u
}

// knownCostRecoveryUnits returns the units AllUnits appends when the option is on.
func knownCostRecoveryUnits(uc *UseCases, infra *Infra) []compose.Unit {
	return []compose.Unit{
		CostSourceComponentUnit(uc, infra),
		AllocationBatchUnit(uc, infra),
		AgreementLineTermUnit(uc, infra),
	}
}

// costSourceComponentViewUseCases adapts the block-local group into the view package's
// UseCases (recovery_charges.go precedent: the view types never sit on block/usecases.go).
func costSourceComponentViewUseCases(uc *UseCases) *cscpkg.UseCases {
	c := &uc.CostSourceComponent
	return &cscpkg.UseCases{
		CreateCostSourceComponent: c.CreateCostSourceComponent,
		ReadCostSourceComponent:   c.ReadCostSourceComponent,
		UpdateCostSourceComponent: c.UpdateCostSourceComponent,
		DeleteCostSourceComponent: c.DeleteCostSourceComponent,
		ListCostSourceComponents:  c.ListCostSourceComponents,
		ListExpenditureLineItems:  c.ListExpenditureLineItems,
	}
}

// allocationBatchViewUseCases adapts the block-local allocation group into the view UseCases.
func allocationBatchViewUseCases(uc *UseCases) *abpkg.UseCases {
	a := &uc.AllocationBatch
	return &abpkg.UseCases{
		CreateAllocationBatch:          a.CreateAllocationBatch,
		UpdateAllocationBatchShares:    a.UpdateAllocationBatchShares,
		PublishAllocationBatch:         a.PublishAllocationBatch,
		GetAllocationBatchListPageData: a.GetAllocationBatchListPageData,
		GetAllocationShareListPageData: a.GetAllocationShareListPageData,
		ReadCostSourceComponent:        a.ReadCostSourceComponent,
		ListAgreementLineTerms:         a.ListAgreementLineTerms,
		GetSubscriptionItemPageData:    a.GetSubscriptionItemPageData,
	}
}

// agreementLineTermViewUseCases adapts the block-local agreement term group into the view UseCases.
func agreementLineTermViewUseCases(uc *UseCases) *altpkg.UseCases {
	t := &uc.AgreementLineTerm
	return &altpkg.UseCases{
		ListAgreementLineTerms:  t.ListAgreementLineTerms,
		ReadChargePolicyVersion: t.ReadChargePolicyVersion,
		ReadChargePolicy:        t.ReadChargePolicy,
		ListProductPricePlans:   t.ListProductPricePlans,
	}
}

// wireRecoverableCosts turns the Recoverable costs tab on for the expenditure module when
// the cost line unit is mounted (the app opted in). Nothing is set otherwise.
func wireRecoverableCosts(mc *compose.MountContext, deps *expendituredomain.ExpenditureModuleDeps, uc *UseCases) {
	routes, ok := compose.RoutesOf[*cscpkg.Routes](mc, "expenditure.cost_source_component")
	if !ok {
		return
	}
	labels, ok := compose.LabelsOf[*cscpkg.Labels](mc, "expenditure.cost_source_component")
	if !ok {
		return
	}
	allocationRoutes := abpkg.DefaultRoutes()
	if ar, ok := compose.RoutesOf[*abpkg.Routes](mc, "expenditure.allocation_batch"); ok {
		allocationRoutes = *ar
	}
	deps.CostSourceComponents = costSourceComponentViewUseCases(uc)
	deps.CostSourceComponentRoutes = *routes
	deps.CostSourceComponentLabels = *labels
	deps.AllocationBatchRoutes = allocationRoutes
}

// chargeTermsDeps carries the subscription Charge terms tab dependencies.
type chargeTermsDeps struct {
	useCases *altpkg.UseCases
	labels   altpkg.Labels
}

// chargeTermsWiringOf returns the Charge terms tab dependencies when the agreement term
// unit is mounted (the app opted in) and its closures are bound; nil otherwise.
func chargeTermsWiringOf(mc *compose.MountContext, uc *UseCases) *chargeTermsDeps {
	labels, ok := compose.LabelsOf[*altpkg.Labels](mc, "subscription.agreement_line_term")
	if !ok || requireAgreementLineTerm(uc) != nil {
		return nil
	}
	return &chargeTermsDeps{useCases: agreementLineTermViewUseCases(uc), labels: *labels}
}
