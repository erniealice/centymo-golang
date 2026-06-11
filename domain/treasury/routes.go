package treasury

// routes.go — treasury-domain route constants + Routes config structs (centymo W5).
//
// The Collection, Disbursement, and TreasuryAdvances route constants + Routes
// types + Default* constructors + RouteMap methods that formerly lived here have
// been extracted into their per-entity packages under
// domain/treasury/<entity>/routes.go as part of the domain-first restructure:
//   - CollectionRoutes        -> collection.Routes
//   - DisbursementRoutes      -> disbursement.Routes
//   - TreasuryAdvancesRoutes  -> advancesdashboard.Routes
//
// The advance-feature route consts moved to the entity that owns them:
//   - selling-side (TreasuryCollection*URL)        -> collection
//   - buying-side  (TreasuryDisbursement*URL)      -> disbursement
//   - workspace    (Advances*URL / Advance*ListURL) -> advancesdashboard
