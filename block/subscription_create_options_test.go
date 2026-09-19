package block

import (
	"context"
	subscription "github.com/erniealice/centymo-golang/domain/subscription/subscription"
	subform "github.com/erniealice/centymo-golang/domain/subscription/subscription/form"
	group "github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	consumerapp "github.com/erniealice/espyna-golang/consumer/app"
	"github.com/erniealice/espyna-golang/consumer/compose"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	"github.com/erniealice/lyngua"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	"net/http/httptest"
	"reflect"
	"testing"
)

func createOptionUseCases() *UseCases {
	uc := &UseCases{}
	uc.Subscription.CreateSubscription = func(context.Context, *subscriptionpb.CreateSubscriptionRequest) (*subscriptionpb.CreateSubscriptionResponse, error) {
		return nil, nil
	}
	return uc
}
func assertCreateMode(t *testing.T, r *inventoryCatalogRegistrar, want bool) {
	t.Helper()
	path := subscription.DefaultRoutes().AddURL
	h := r.get[path]
	if h == nil {
		t.Fatal("add handler not mounted")
	}
	ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"subscription:create"}))
	result := h.Handle(ctx, &view.ViewContext{Request: httptest.NewRequest("GET", path, nil), BusinessType: "leasing"})
	d, ok := result.Data.(*subform.Data)
	if !ok || d.SuggestPlanTerm != want {
		t.Fatalf("mode=%v result=%+v", want, result)
	}
}
func TestSubscriptionCreateOptionsReachMountedAction(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		units := AllUnits(createOptionUseCases(), &Infra{}, WithSubscriptionCreateOptions(subscription.CreateOptions{CommencementPricing: enabled}))
		for _, u := range units {
			if u.Key == "subscription.subscription" {
				r := &inventoryCatalogRegistrar{get: map[string]view.View{}}
				if err := u.Mount(&compose.MountContext{Routes: r}); err != nil {
					t.Fatal(err)
				}
				assertCreateMode(t, r, enabled)
			}
		}
	}
}
func TestSubscriptionUnitCreateOptionsLastWins(t *testing.T) {
	for _, tc := range []struct {
		opts []subscription.CreateOptions
		want bool
	}{
		{nil, false}, {[]subscription.CreateOptions{{CommencementPricing: true}}, true},
		{[]subscription.CreateOptions{{CommencementPricing: true}, {}}, false},
	} {
		u := SubscriptionUnit(createOptionUseCases(), &Infra{}, tc.opts...)
		r := &inventoryCatalogRegistrar{get: map[string]view.View{}}
		if err := u.Mount(&compose.MountContext{Routes: r}); err != nil {
			t.Fatal(err)
		}
		assertCreateMode(t, r, tc.want)
	}
}
func TestSubscriptionCreateOptionsCoexist(t *testing.T) {
	opts := []EngineOption{WithProductAssets(), WithInventoryCatalogMounts(), WithSubscriptionGroupOptions(group.Options{Roster: group.RowOptions{GroupByField: "client_attributes.group"}}), WithSubscriptionCreateOptions(subscription.CreateOptions{CommencementPricing: true})}
	cfg := engineConfig{}
	for _, o := range opts {
		o(&cfg)
	}
	if !cfg.productAssets || !cfg.inventoryCatalogMounts || cfg.subscriptionGroupOptions.Roster.GroupByField != "client_attributes.group" {
		t.Fatal("existing options lost")
	}
	units := AllUnits(createOptionUseCases(), &Infra{}, opts...)
	found := false
	for _, u := range units {
		if u.Key == "product.product_inventory" {
			found = true
		}
		if u.Key == "subscription.subscription" {
			r := &inventoryCatalogRegistrar{get: map[string]view.View{}}
			if err := u.Mount(&compose.MountContext{Routes: r}); err != nil {
				t.Fatal(err)
			}
			assertCreateMode(t, r, true)
		}
	}
	if !found {
		t.Fatal("inventory mount lost")
	}
}

// Fill only test dependency functions, preserving the real Block configuration,
// validation and registration path without constructing a database container.
func stubCommercialFunctions(v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if !f.CanSet() {
			continue
		}
		switch f.Kind() {
		case reflect.Struct:
			stubCommercialFunctions(f)
		case reflect.Func:
			if f.IsNil() {
				typ := f.Type()
				f.Set(reflect.MakeFunc(typ, func([]reflect.Value) []reflect.Value {
					out := make([]reflect.Value, typ.NumOut())
					for j := range out {
						out[j] = reflect.Zero(typ.Out(j))
					}
					return out
				}))
			}
		}
	}
}
func TestLegacySubscriptionCreatePolicy(t *testing.T) {
	for _, tc := range []struct {
		name   string
		policy bool
		module bool
	}{{"explicit", true, true}, {"omitted", false, true}, {"policy_only_keeps_all", true, false}} {
		t.Run(tc.name, func(t *testing.T) {
			uc := createOptionUseCases()
			stubCommercialFunctions(reflect.ValueOf(uc).Elem())
			opts := []BlockOption{WithUseCases(uc)}
			if tc.module {
				opts = append(opts, WithSubscription())
			}
			if tc.policy {
				opts = append(opts, WithSubscriptionCreatePolicy(subscription.CreateOptions{CommencementPricing: true}))
			}
			r := &inventoryCatalogRegistrar{get: map[string]view.View{}}
			ctx := &consumerapp.AppContext{Routes: r, BusinessType: "leasing", Translations: lynguaV1.NewTranslationProviderFromFS(lyngua.TranslationsFS)}
			if err := Block(opts...)(ctx); err != nil {
				t.Fatal(err)
			}
			assertCreateMode(t, r, tc.policy)
		})
	}
}
