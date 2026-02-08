package centymo

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
}
