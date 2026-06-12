package block

import (
	"context"
	"testing"

	subscriptiondom "github.com/erniealice/centymo-golang/domain/subscription"
	lyngua "github.com/erniealice/lyngua"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"

	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
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

// newPlanUseCases returns a UseCases value with no-op stubs for the required
// Plan CRUD fields, as validated by RequireFor when WithPlan() is enabled.
func newPlanUseCases() *UseCases {
	return &UseCases{
		Plan: PlanUseCases{
			ListPlans: func(_ context.Context, _ *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error) {
				return &planpb.ListPlansResponse{}, nil
			},
			ReadPlan: func(_ context.Context, _ *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error) {
				return &planpb.ReadPlanResponse{}, nil
			},
			CreatePlan: func(_ context.Context, _ *planpb.CreatePlanRequest) (*planpb.CreatePlanResponse, error) {
				return &planpb.CreatePlanResponse{}, nil
			},
			UpdatePlan: func(_ context.Context, _ *planpb.UpdatePlanRequest) (*planpb.UpdatePlanResponse, error) {
				return &planpb.UpdatePlanResponse{}, nil
			},
			DeletePlan: func(_ context.Context, _ *planpb.DeletePlanRequest) (*planpb.DeletePlanResponse, error) {
				return &planpb.DeletePlanResponse{}, nil
			},
		},
	}
}

// newSubscriptionUseCases returns a UseCases value with no-op stubs for the required
// Subscription CRUD fields, as validated by RequireFor when WithSubscription() is enabled.
func newSubscriptionUseCases() *UseCases {
	return &UseCases{
		Subscription: SubscriptionUseCases{
			GetSubscriptionListPageData: func(_ context.Context, _ *subscriptionpb.GetSubscriptionListPageDataRequest) (*subscriptionpb.GetSubscriptionListPageDataResponse, error) {
				return &subscriptionpb.GetSubscriptionListPageDataResponse{}, nil
			},
			CreateSubscription: func(_ context.Context, _ *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error) {
				return &subscriptionpb.CreateSubscriptionResponse{}, nil
			},
			ReadSubscription: func(_ context.Context, _ *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error) {
				return &subscriptionpb.ReadSubscriptionResponse{}, nil
			},
			UpdateSubscription: func(_ context.Context, _ *subscriptionpb.UpdateSubscriptionRequest) (*subscriptionpb.UpdateSubscriptionResponse, error) {
				return &subscriptionpb.UpdateSubscriptionResponse{}, nil
			},
			DeleteSubscription: func(_ context.Context, _ *subscriptionpb.DeleteSubscriptionRequest) (*subscriptionpb.DeleteSubscriptionResponse, error) {
				return &subscriptionpb.DeleteSubscriptionResponse{}, nil
			},
		},
	}
}

func TestBlockLoadsRouteOverridesForSelectedModules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		option      BlockOption
		useCases    *UseCases
		expectPath  string
		defaultPath string
	}{
		{
			name:        "subscription routes use service override",
			option:      WithSubscription(),
			useCases:    newSubscriptionUseCases(),
			expectPath:  "/app/memberships/list/{status}",
			defaultPath: subscriptiondom.SubscriptionListURL,
		},
		{
			name:        "plan routes use service override",
			option:      WithPlan(),
			useCases:    newPlanUseCases(),
			expectPath:  "/app/packages/list/{status}",
			defaultPath: subscriptiondom.PlanListURL,
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
				// ctx.DB intentionally unset — the centymo DataSource duck was
				// deleted (20260612-datasource-typed-path W6); Block() no longer
				// reads ctx.DB.
			}

			if err := Block(tc.option, WithUseCases(tc.useCases))(ctx); err != nil {
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
