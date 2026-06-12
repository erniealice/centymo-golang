// Package shared holds centymo cross-domain leaf contracts that are NOT owned by
// any single commerce domain and are NOT framework concerns (those go to pyeza).
//
// It is a charter'd domain/shared/ leaf: it imports nothing from centymo, so any
// entity package may import it DIRECTLY without risking an import cycle. This is
// the in-module home that lets the per-entity view packages stop importing the
// root `package centymo` (which imports the domain facades via routes_config.go —
// an entity -> root -> facade -> entity cycle).
//
// 20260612-datasource-typed-path W6 — the DataSource duck interface and the
// LocationDisplayName pass-through stub were DELETED here. All former duck users
// migrated to narrow typed closures on the block UseCases; the location-name
// resolution keeps its pass-through fallback inlined into ResolveLocationName
// below. The remaining LocationResolver typed path is entydad-bound (WL deferral).
package shared

import "context"

// LocationResolver maps a location id/slug to a human display name. It is fed at
// composition time by the typed espyna location use-case (ListLocations →
// id→name map). It MUST be pass-through safe: when the id is unknown it returns
// the input unchanged (the inventory list, for instance, passes a human slug
// like "ayala-central-bloc" that has no matching location row).
type LocationResolver func(ctx context.Context, id string) string

// ResolveLocationName resolves a location id/slug to a display name using the
// injected resolver, falling back to a pass-through (return the id unchanged)
// when no resolver is wired. Every view consumer routes through this helper so
// the nil-resolver path stays identical to today's behaviour.
func ResolveLocationName(ctx context.Context, resolver LocationResolver, id string) string {
	if resolver == nil {
		return id
	}
	return resolver(ctx, id)
}
