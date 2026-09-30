package billable_charge_test

import (
	"context"
	"fmt"
	"testing"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"

	bctest "github.com/erniealice/centymo-golang/domain/subscription/billable_charge/internal/bctest"
)

// R5 M4: the client name map pages through every client (no silent 100-row cap).
func TestClientNamesPagesThroughEveryClient(t *testing.T) {
	f := bctest.Seed()
	for i := 1; i <= 250; i++ {
		id, name := fmt.Sprintf("c%03d", i), fmt.Sprintf("Client %d", i)
		f.Clients = append(f.Clients, &clientpb.Client{Id: id, Name: &name})
	}
	names := f.UseCases().ClientNames(context.Background())
	if len(names) != 250 || names["c101"] != "Client 101" || names["c250"] != "Client 250" {
		t.Errorf("names = %d entries, c101=%q c250=%q", len(names), names["c101"], names["c250"])
	}
	if f.ListClientPages != 3 {
		t.Errorf("list pages = %d, want 3", f.ListClientPages)
	}
}
