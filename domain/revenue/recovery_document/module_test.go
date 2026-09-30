package recovery_document_test

import (
	"net/http"
	"testing"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	rdtest "github.com/erniealice/centymo-golang/domain/revenue/recovery_document/internal/rdtest"
)

func TestRoutesRegisterWithoutServeMuxConflict(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration conflict: %v", r)
		}
	}()
	rdtest.NewModule(rdtest.Seed()).RegisterRoutes(rdtest.MuxRegistrar{Mux: http.NewServeMux()})
}

func TestRouteMapKeys(t *testing.T) {
	m := rd.DefaultRoutes().RouteMap()
	want := map[string]string{
		"recovery_document.list":       "/revenue/recovery-documents/list/{status}",
		"recovery_document.table":      "/action/recovery-document/table/{status}",
		"recovery_document.detail":     "/revenue/recovery-documents/detail/{id}",
		"recovery_document.tab_action": "/action/recovery-document/detail/{id}/tab/{tab}",
		"recovery_document.void":       "/action/recovery-document/void/{id}",
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
