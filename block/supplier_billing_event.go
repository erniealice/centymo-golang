// Package block — supplier_billing_event domain wiring
// (20260517-advance-cash-events Plan B Phase 7).
//
// Mirror of advances_dashboard.go for the buying-side MILESTONE anchor.
// Wires the list + detail + Recognize routes from
// `treasurydomain.TreasuryAdvancesRoutes.SupplierBillingEvent*URL`. Recognize is
// fed by `useCases.TreasuryAdvances.RecognizeMilestoneAdvanceDisbursement`.
package block

import (
	pyeza "github.com/erniealice/pyeza-golang"

	expendituredomain "github.com/erniealice/centymo-golang/domain/expenditure"
	treasurydomain "github.com/erniealice/centymo-golang/domain/treasury"
)

// supplierBillingEventWiring holds everything wireSupplierBillingEventModule
// needs from the surrounding Block() scope.
type supplierBillingEventWiring struct {
	routes treasurydomain.TreasuryAdvancesRoutes
}

// wireSupplierBillingEventModule mounts the list, detail, and Recognize
// routes for SupplierBillingEvent rows.
//
// Nil-safe: all three view ops on UseCases.Expenditure /
// UseCases.TreasuryAdvances can be missing — the underlying view module's
// nil-safety handles each independently.
func wireSupplierBillingEventModule(ctx *pyeza.AppContext, useCases *UseCases, w supplierBillingEventWiring) {
	deps := expendituredomain.SupplierBillingEventModuleDeps{
		Routes:       w.routes,
		Labels:       expendituredomain.DefaultSupplierBillingEventLabels(),
		CommonLabels: ctx.Common,
		// SupplierBillingEvent table reads live on UseCases.Expenditure (the
		// repo lives in the expenditure domain). The block adapter wires
		// these when the espyna provider registers the supplier_billing_event
		// adapter.
		ListSupplierBillingEvents: useCases.Expenditure.ListSupplierBillingEvents,
		ReadSupplierBillingEvent:  useCases.Expenditure.ReadSupplierBillingEvent,
		Recognize:                 useCases.TreasuryAdvances.RecognizeMilestoneAdvanceDisbursement,
	}
	module := expendituredomain.NewSupplierBillingEventModule(deps)
	module.RegisterRoutes(ctx.Routes)
}
