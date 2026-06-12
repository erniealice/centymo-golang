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
func LocationDisplayName(slug string) string {
	return slug
}
