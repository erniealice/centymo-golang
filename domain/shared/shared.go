// Package shared holds centymo cross-domain leaf contracts that are NOT owned by
// any single commerce domain and are NOT framework concerns (those go to pyeza).
//
// It is a charter'd domain/shared/ leaf: it imports nothing from centymo, so any
// entity package may import it DIRECTLY without risking an import cycle. This is
// the in-module home that lets the per-entity view packages stop importing the
// root `package centymo` (which imports the domain facades via routes_config.go —
// an entity -> root -> facade -> entity cycle). The symbols here are slated to
// migrate out of centymo entirely (DataSource -> pyeza data port;
// LocationDisplayName -> entydad location entity) under Wave P / WL; until then
// they live in this leaf.
package shared

import "context"

// DataSource provides technology-agnostic data access for views.
// Consumer apps satisfy this interface by wrapping their database adapter.
// espyna's DatabaseAdapter already matches this signature directly.
type DataSource interface {
	ListSimple(ctx context.Context, collection string) ([]map[string]any, error)
	Create(ctx context.Context, collection string, data map[string]any) (map[string]any, error)
	Read(ctx context.Context, collection string, id string) (map[string]any, error)
	Update(ctx context.Context, collection string, id string, data map[string]any) (map[string]any, error)
	Delete(ctx context.Context, collection string, id string) error
	HardDelete(ctx context.Context, collection string, id string) error
}

// LocationDisplayName returns a human label for a location slug/ID.
// TODO(WL): resolve the real name via the entydad location entity (location_area).
// The hardcoded demo LocationMap was removed 2026-06-12 to avoid confusion; until the
// typed location path is wired, this passes the slug/ID through unchanged.
//
// This stub remains the fallback when no LocationResolver is injected (tests, or a
// half-wired composition root). Live composition feeds a DB-backed resolver via the
// typed espyna location use-case (see block wiring); see ResolveLocationName.
func LocationDisplayName(slug string) string {
	return slug
}

// LocationResolver maps a location id/slug to a human display name. It is fed at
// composition time by the typed espyna location use-case (ListLocations →
// id→name map). It MUST be pass-through safe: when the id is unknown it returns
// the input unchanged, preserving the LocationDisplayName stub's behaviour (the
// inventory list, for instance, passes a human slug like "ayala-central-bloc"
// that has no matching location row).
type LocationResolver func(ctx context.Context, id string) string

// ResolveLocationName resolves a location id/slug to a display name using the
// injected resolver, falling back to the LocationDisplayName pass-through stub
// when no resolver is wired. Every view consumer routes through this helper so
// the nil-resolver path stays identical to today's behaviour.
func ResolveLocationName(ctx context.Context, resolver LocationResolver, id string) string {
	if resolver == nil {
		return LocationDisplayName(id)
	}
	return resolver(ctx, id)
}
