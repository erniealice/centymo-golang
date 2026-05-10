package block

import (
	"context"
	"testing"

	centymo "github.com/erniealice/centymo-golang"
	consumer "github.com/erniealice/espyna-golang/consumer"
	lyngua "github.com/erniealice/lyngua"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"
)

type testRouteRegistrar struct {
	getPaths map[string]struct{}
}

func newTestRouteRegistrar() *testRouteRegistrar {
	return &testRouteRegistrar{getPaths: make(map[string]struct{})}
}

func (r *testRouteRegistrar) GET(path string, _ view.View, _ ...string) {
	r.getPaths[path] = struct{}{}
}

func (r *testRouteRegistrar) POST(_ string, _ view.View, _ ...string) {}

type testDataSource struct{}

func (d *testDataSource) ListSimple(context.Context, string) ([]map[string]any, error) {
	return nil, nil
}

func (d *testDataSource) Create(context.Context, string, map[string]any) (map[string]any, error) {
	return map[string]any{}, nil
}

func (d *testDataSource) Read(context.Context, string, string) (map[string]any, error) {
	return map[string]any{}, nil
}

func (d *testDataSource) Update(context.Context, string, string, map[string]any) (map[string]any, error) {
	return map[string]any{}, nil
}

func (d *testDataSource) Delete(context.Context, string, string) error {
	return nil
}

func (d *testDataSource) HardDelete(context.Context, string, string) error {
	return nil
}

func TestBlockLoadsRouteOverridesForSelectedModules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		option      BlockOption
		expectPath  string
		defaultPath string
	}{
		{
			name:        "subscription routes use service override",
			option:      WithSubscription(),
			expectPath:  "/app/memberships/list/{status}",
			defaultPath: centymo.SubscriptionListURL,
		},
		{
			name:        "plan routes use service override",
			option:      WithPlan(),
			expectPath:  "/app/packages/list/{status}",
			defaultPath: centymo.PlanListURL,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			routes := newTestRouteRegistrar()
			ctx := &pyeza.AppContext{
				Routes:       routes,
				Common:       pyeza.CommonLabels{},
				BusinessType: "service",
				Translations: lynguaV1.NewTranslationProviderFromFS(lyngua.TranslationsFS),
				UseCases:     &consumer.UseCases{},
				DB:           &testDataSource{},
			}

			if err := Block(tc.option)(ctx); err != nil {
				t.Fatalf("Block() returned error: %v", err)
			}

			if _, ok := routes.getPaths[tc.expectPath]; !ok {
				t.Fatalf("expected route %q to be registered, got %v", tc.expectPath, keys(routes.getPaths))
			}

			if tc.defaultPath != tc.expectPath {
				if _, ok := routes.getPaths[tc.defaultPath]; ok {
					t.Fatalf("default route %q should not be registered when override %q exists", tc.defaultPath, tc.expectPath)
				}
			}
		})
	}
}

func keys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
