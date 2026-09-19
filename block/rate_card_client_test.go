package block

import (
	"context"
	schedule "github.com/erniealice/centymo-golang/domain/subscription/price_schedule"
	schedulelist "github.com/erniealice/centymo-golang/domain/subscription/price_schedule/list"
	"github.com/erniealice/espyna-golang/consumer/compose"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	schedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	"net/http/httptest"
	"testing"
)

func TestPriceScheduleClientNameWiring(t *testing.T) {
	uc := &UseCases{}
	clientID := "client-1"
	clientName := "Client Corporation"
	uc.PriceSchedule.ListPriceSchedules = func(context.Context, *schedulepb.ListPriceSchedulesRequest) (*schedulepb.ListPriceSchedulesResponse, error) {
		return &schedulepb.ListPriceSchedulesResponse{Data: []*schedulepb.PriceSchedule{{Id: "schedule-1", ClientId: &clientID, Active: true}}}, nil
	}
	uc.Entity.Client.ListClients = func(context.Context, *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error) {
		return &clientpb.ListClientsResponse{Data: []*clientpb.Client{{Id: clientID, Name: &clientName, Active: true}}}, nil
	}
	u := PriceScheduleUnit(uc, &Infra{})
	r := &inventoryCatalogRegistrar{get: map[string]view.View{}}
	if err := u.Mount(&compose.MountContext{Routes: r}); err != nil {
		t.Fatal(err)
	}
	routes := u.Routes.(*schedule.Routes)
	req := httptest.NewRequest("GET", "/price-schedules/list/active", nil)
	req.SetPathValue("status", "active")
	ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"price_schedule:list"}))
	result := r.get[routes.ListURL].Handle(ctx, &view.ViewContext{Request: req})
	data, ok := result.Data.(*schedulelist.PageData)
	if !ok {
		t.Fatalf("unexpected result %+v", result)
	}
	if len(data.Table.Rows) != 1 || data.Table.Rows[0].Cells[4].Value != "Client Corporation" {
		t.Fatalf("unexpected rows %+v", data.Table.Rows)
	}
}
