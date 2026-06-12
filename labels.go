package centymo

import shared "github.com/erniealice/centymo-golang/domain/shared"

// labels.go — RESIDUAL after centymo W7.
//
// LocationDisplayName is re-exported from domain/shared (centymo restructure) so
// the per-entity view packages import the leaf, not the root — breaking the
// entity -> root -> facade import cycle. The external root consumer
// (centymo.LocationDisplayName) keeps resolving unchanged.
// (entydad-bound; WL deferral, not yet landed.)

// LocationDisplayName returns a human label for a location slug/ID (pass-through
// until the typed entydad location path is wired — see domain/shared).
func LocationDisplayName(slug string) string { return shared.LocationDisplayName(slug) }
