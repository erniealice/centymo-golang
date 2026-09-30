package collection_application_test

import (
	"context"
	"fmt"
	"testing"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"

	ratest "github.com/erniealice/centymo-golang/domain/treasury/collection_application/internal/ratest"
)

// R5 M4: the fallback picker pages through every client, so client #101 exists.
func TestClientOptionsPageThroughEveryClient(t *testing.T) {
	f := &ratest.Fake{}
	for i := 1; i <= 250; i++ {
		id, name := fmt.Sprintf("c%03d", i), fmt.Sprintf("Client %d", i)
		f.Clients = append(f.Clients, &clientpb.Client{Id: id, Name: &name})
	}
	opts := f.UseCases().ClientOptions(context.Background(), "c101")
	if len(opts) != 250 || opts[100].Value != "c101" || !opts[100].Selected || f.ListClientPages != 3 {
		t.Errorf("options=%d pages=%d c101=%+v", len(opts), f.ListClientPages, opts[100])
	}
}
