package centymo

import (
	pricesched "github.com/erniealice/centymo-golang/domain/subscription/price_schedule"
	subscriptionentity "github.com/erniealice/centymo-golang/domain/subscription/subscription"
	advancesdashboard "github.com/erniealice/centymo-golang/domain/treasury/treasuryadvancesdashboard"
	disbursement "github.com/erniealice/centymo-golang/domain/treasury/disbursement"
	shared "github.com/erniealice/centymo-golang/domain/treasury/shared"
)

// NOTE (centymo restructure): these compatibility aliases point DIRECTLY at the
// owning entity packages (domain/<d>/<entity>) rather than the domain facade
// (domain/<d>). The root package must not import the domain facades — that would
// create an import cycle (root -> facade -> entity -> root, since entity views
// still import the root for DataSource/LocationDisplayName). Pointing at the
// root-free entity packages preserves the exact external symbol surface
// (centymo.SubscriptionRoutes, centymo.DisbursementRoutes, …) with ZERO
// behaviour change. The names below are byte-identical to the prior shim.

// ── centymo W4 subscription-domain compatibility shim ────────────────────────
// The Subscription/PriceSchedule route types + their Default* constructors moved
// to domain/subscription/ (centymo W4). entydad-golang/block/route_loading.go is
// an EXTERNAL consumer (outside this wave's edit scope) that still references
// centymo.SubscriptionRoutes / centymo.PriceScheduleRoutes and their Default*
// constructors. These thin aliases + forwarders keep that consumer compiling
// with ZERO behaviour change (pure type-identity aliases). Remove once entydad
// is re-pointed to domain/subscription directly (W9 / entydad-coordinated).
type SubscriptionRoutes = subscriptionentity.Routes
type PriceScheduleRoutes = pricesched.Routes

func DefaultSubscriptionRoutes() SubscriptionRoutes { return subscriptionentity.DefaultRoutes() }
func DefaultPriceScheduleRoutes() PriceScheduleRoutes {
	return pricesched.DefaultRoutes()
}

// ── centymo W5 treasury-domain compatibility shim ────────────────────────────
// Treasury types (Collection/Disbursement labels+routes, TreasuryAdvancesRoutes,
// the AdvanceRecognizeMilestone view I/O) moved to domain/treasury/ (centymo W5).
// The not-yet-migrated W6 view packages still reference a subset of them via the
// centymo root:
//   - views/expenditure/*            -> DisbursementRoutes / DisbursementLabels /
//     DisbursementFormLabels (the expense "pay"
//     flow creates a pre-linked disbursement)
//   - views/supplier_billing_event/* -> TreasuryAdvancesRoutes (+ its Default*)
//     and AdvanceRecognizeMilestoneInput/Output
//   - domain/subscription/...  -> AdvanceRecognizeMilestoneInput/Output
//     (already-migrated W4 billing-event action)
//
// These thin aliases + forwarders keep those consumers compiling with ZERO
// behaviour change. Removed as each consuming domain migrates (W6 / W9).
type DisbursementRoutes = disbursement.Routes
type DisbursementLabels = disbursement.Labels
type DisbursementFormLabels = disbursement.FormLabels
type TreasuryAdvancesRoutes = advancesdashboard.Routes
type AdvanceRecognizeMilestoneInput = shared.AdvanceRecognizeMilestoneInput
type AdvanceRecognizeMilestoneOutput = shared.AdvanceRecognizeMilestoneOutput

func DefaultTreasuryAdvancesRoutes() TreasuryAdvancesRoutes {
	return advancesdashboard.DefaultRoutes()
}

// Three-level routing system for centymo views:
//
// Level 1: Generic defaults from Go consts (this file).
//   DefaultXxxRoutes() constructors return structs populated from the route
//   constants defined in routes.go. These are sensible defaults that work
//   out of the box for any app.
//
// Level 2: Industry-specific overrides via JSON (loaded by consumer apps).
//   Consumer apps can load a JSON config that partially overrides the
//   default routes. Struct fields carry json tags for unmarshalling.
//
// Level 3: App-specific overrides via Go field assignment (optional).
//   After loading defaults and/or JSON, consumer apps can programmatically
//   set individual fields to further customize routing.
//
// Each route struct also exposes a RouteMap() method that returns a
// map[string]string keyed by dot-notation identifiers (e.g. "product.list"),
// useful for template rendering, URL resolution, and debugging.

// Procurement entity route types and Default* constructors moved to
// domain/procurement/routes.go in W7. Compatibility aliases can be added
// here if external consumers need them (W9 / consumer-coordinated).

// MapTableLabels is a shared helper used across all centymo view modules to
// produce a types.TableLabels from pyeza CommonLabels. Defined here to avoid
// duplication; all block module wirings call this.
func mapTableLabelsFromStrings(search, searchPlaceholder, sortAsc, sortDesc, noResults, loading string) struct{} {
	// Placeholder — actual implementation lives in the block package; this
	// comment documents the cross-module convention.
	return struct{}{}
}
