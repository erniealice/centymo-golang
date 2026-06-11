// Package shared holds centymo cross-domain leaf contracts that are NOT owned by
// any single commerce domain and are NOT framework concerns (those go to pyeza).
//
// It is a charter'd domain/shared/ leaf: it imports nothing from centymo, so any
// entity package may import it DIRECTLY without risking an import cycle. This is
// the in-module home that lets the per-entity view packages stop importing the
// root `package centymo` (which imports the domain facades via routes_config.go —
// an entity -> root -> facade -> entity cycle). The symbols here are slated to
// migrate out of centymo entirely (DataSource -> pyeza data port; Location* ->
// entydad location entity) under Wave P / WL; until then they live in this leaf.
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

// LocationMap maps a location slug to its display name. (entydad-bound: WL
// deferral — replace the hardcoded map with the location entity's real name.)
var LocationMap = map[string]string{
	"ayala-central-bloc": "Ayala Central Bloc",
	"sm-city-cebu":       "SM City Cebu",
	"ayala-center-cebu":  "Ayala Center Cebu",
	"robinsons-galleria": "Robinsons Galleria",
}

// LocationDisplayName returns the display name for a location slug.
func LocationDisplayName(slug string) string {
	if name, ok := LocationMap[slug]; ok {
		return name
	}
	return slug
}
