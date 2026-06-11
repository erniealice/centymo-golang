package centymo

import shared "github.com/erniealice/centymo-golang/domain/shared"

// labels.go — RESIDUAL after centymo W7.
//
// LocationMap / LocationDisplayName are re-exported from domain/shared (centymo
// restructure) so the per-entity view packages import the leaf, not the root —
// breaking the entity -> root -> facade import cycle. External root consumers
// (centymo.LocationDisplayName / centymo.LocationMap) keep resolving unchanged.
// (entydad-bound; WL deferral, not yet landed.)
var LocationMap = shared.LocationMap

// LocationDisplayName returns the display name for a location slug.
func LocationDisplayName(slug string) string { return shared.LocationDisplayName(slug) }
