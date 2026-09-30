package billable_charge_test

import (
	"net/http"
	"testing"

	bc "github.com/erniealice/centymo-golang/domain/subscription/billable_charge"
	bctest "github.com/erniealice/centymo-golang/domain/subscription/billable_charge/internal/bctest"
)

// The live ServeMux panics on conflicting patterns, so registering every
// billable charge route on a real mux proves there is no self-conflict.
func TestRoutesRegisterWithoutServeMuxConflict(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration conflict: %v", r)
		}
	}()
	bctest.NewModule(bctest.Seed()).RegisterRoutes(bctest.MuxRegistrar{Mux: http.NewServeMux()})
}

func TestRouteMapKeys(t *testing.T) {
	m := bc.DefaultRoutes().RouteMap()
	want := map[string]string{
		"billable_charge.list":   "/revenue/billable-charges/list/{status}",
		"billable_charge.table":  "/action/billable-charge/table/{status}",
		"billable_charge.issue":  "/action/billable-charge/issue",
		"billable_charge.adjust": "/action/billable-charge/adjust/{id}",
	}
	for k, v := range want {
		if m[k] != v {
			t.Errorf("%s = %q, want %q", k, m[k], v)
		}
	}
	if len(m) != len(want) {
		t.Errorf("route map has %d keys, want %d", len(m), len(want))
	}
}
