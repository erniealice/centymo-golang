package block

import (
	"context"
	"net/http/httptest"
	"reflect"
	"testing"

	productpkg "github.com/erniealice/centymo-golang/domain/product/product"
	productform "github.com/erniealice/centymo-golang/domain/product/product/form"
	compose "github.com/erniealice/espyna-golang/consumer/compose"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

type inventoryCatalogRegistrar struct{ get map[string]view.View }

func (r *inventoryCatalogRegistrar) GET(path string, v view.View, _ ...string) {
	if _, exists := r.get[path]; exists {
		panic("duplicate GET: " + path)
	}
	r.get[path] = v
}
func (r *inventoryCatalogRegistrar) POST(string, view.View, ...string) {}

func TestInventoryCatalogMounts(t *testing.T) {
	r := &inventoryCatalogRegistrar{get: map[string]view.View{}}
	uc, infra := &UseCases{}, &Infra{}
	for _, u := range []compose.Unit{ProductUnit(uc, infra), ProductInventoryUnit(uc, infra), ProductSuppliesUnit(uc, infra), PriceScheduleUnit(uc, infra), PriceScheduleInventoryUnit(uc, infra)} {
		if err := u.Mount(&compose.MountContext{Routes: r}); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{"/inventory/products/list/{status}", "/inventory/supplies/list/{status}", "/inventory/price-schedules/list/{status}", "/products/list/{status}", "/price-schedules/list/{status}"} {
		if r.get[path] == nil {
			t.Errorf("missing mounted page %s", path)
		}
	}
}

func TestInventoryCatalogDefaults(t *testing.T) {
	for _, tc := range []struct {
		name       string
		build      func(*UseCases, *Infra) compose.Unit
		kinds      []string
		permission string
	}{
		{"products", ProductInventoryUnit, []string{"stocked_good", "non_stocked_good"}, "product"},
		{"supplies", ProductSuppliesUnit, []string{"consumable"}, "supplies"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			u := tc.build(&UseCases{}, &Infra{})
			r := &inventoryCatalogRegistrar{get: map[string]view.View{}}
			if err := u.Mount(&compose.MountContext{Routes: r}); err != nil {
				t.Fatal(err)
			}
			path := u.Routes.(*productpkg.Routes).AddURL
			vc := &view.ViewContext{Request: httptest.NewRequest("GET", path, nil)}
			ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{tc.permission + ":create"}))
			result := r.get[path].Handle(ctx, vc)
			data, ok := result.Data.(*productform.Data)
			if !ok {
				t.Fatalf("drawer data %T, result %+v", result.Data, result)
			}
			var kinds []string
			for _, o := range data.ProductKindOptions {
				kinds = append(kinds, o.Value)
			}
			if !reflect.DeepEqual(kinds, tc.kinds) {
				t.Errorf("kinds %v, want %v", kinds, tc.kinds)
			}
			if data.ProductKind != tc.kinds[0] || data.DeliveryMode != "shipped" || data.TrackingMode != "bulk" {
				t.Errorf("wrong defaults: %s/%s/%s", data.ProductKind, data.DeliveryMode, data.TrackingMode)
			}
			denied := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"service:create"}))
			if _, ok := r.get[path].Handle(denied, vc).Data.(*productform.Data); ok {
				t.Fatal("service permission opened inventory drawer")
			}
		})
	}
}

func TestInventoryCatalogOptIn(t *testing.T) {
	baseline := AllUnits(nil, nil)
	enabled := AllUnits(nil, nil, WithInventoryCatalogMounts())
	if len(enabled) != len(baseline)+3 {
		t.Fatalf("unit count %d -> %d", len(baseline), len(enabled))
	}
	for i, u := range baseline {
		if enabled[i].Key != u.Key {
			t.Fatalf("default unit changed at %d", i)
		}
	}
	for _, u := range enabled[len(baseline):] {
		if u.Mount == nil {
			t.Errorf("%s has no mount", u.Key)
		}
	}
}
