package collection_application_test

import (
	"net/http"
	"testing"

	ra "github.com/erniealice/centymo-golang/domain/treasury/collection_application"
	ratest "github.com/erniealice/centymo-golang/domain/treasury/collection_application/internal/ratest"
)

func TestRoutesRegisterWithoutServeMuxConflict(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration conflict: %v", r)
		}
	}()
	ratest.NewModule(&ratest.Fake{}).RegisterRoutes(ratest.MuxRegistrar{Mux: http.NewServeMux()})
}

func TestRouteMapKeys(t *testing.T) {
	m := ra.DefaultRoutes().RouteMap()
	want := map[string]string{
		"collection_application.receive_apply": "/action/collection/receive-apply",
		"collection_application.preview":       "/action/collection/receive-apply/preview",
		"collection_application.reverse":       "/action/collection-application/reverse/{id}",
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
