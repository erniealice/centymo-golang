package form

import (
	"context"
	"log"
	"strings"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
)

// ResolvePlanGroupForClientLabel renders the {{.ClientName}}-templated
// "For {ClientName}" group header. Falls back gracefully when the label
// has no template directive or the client name is empty.
func ResolvePlanGroupForClientLabel(template, clientName string) string {
	if clientName == "" {
		return template
	}
	return strings.ReplaceAll(template, "{{.ClientName}}", clientName)
}

// LoadClientOptions fetches the client list and converts to select options.
func LoadClientOptions(ctx context.Context, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)) []map[string]string {
	if listClients == nil {
		return nil
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		log.Printf("Failed to load clients for dropdown: %v", err)
		return nil
	}
	var options []map[string]string
	for _, c := range resp.GetData() {
		label := c.GetId()
		if u := c.GetUser(); u != nil {
			first := u.GetFirstName()
			last := u.GetLastName()
			if first != "" || last != "" {
				label = first + " " + last
			}
		}
		options = append(options, map[string]string{
			"Value": c.GetId(),
			"Label": label,
		})
	}
	return options
}

// LoadPlanOptions fetches the plan list and converts to select options.
func LoadPlanOptions(ctx context.Context, listPlans func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)) []map[string]string {
	if listPlans == nil {
		return nil
	}
	resp, err := listPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		log.Printf("Failed to load plans for dropdown: %v", err)
		return nil
	}
	var options []map[string]string
	for _, p := range resp.GetData() {
		if !p.GetActive() {
			continue
		}
		options = append(options, map[string]string{
			"Value": p.GetId(),
			"Label": p.GetName(),
		})
	}
	return options
}

// LoadPlanOptionGroups builds the grouped Plan picker for the subscription
// drawer per plan §5.1. Group order: client-scoped first ("For {ClientName}"),
// general ("General packages") second. Empty groups are omitted.
//
// The list is filtered post-fetch in Go because TypedFilter doesn't yet
// expose a NULL/NOT-NULL primitive on string fields. Volume is small.
func LoadPlanOptionGroups(
	ctx context.Context,
	listPlans func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error),
	clientID, clientName string,
	l Labels,
) []OptionGroup {
	if listPlans == nil {
		return nil
	}
	resp, err := listPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		log.Printf("Failed to load plans for grouped picker: %v", err)
		return nil
	}

	var clientPlans, masterPlans []map[string]string
	for _, p := range resp.GetData() {
		if !p.GetActive() {
			continue
		}
		entry := map[string]string{"Value": p.GetId(), "Label": p.GetName()}
		switch cid := p.GetClientId(); {
		case cid == "":
			masterPlans = append(masterPlans, entry)
		case clientID != "" && cid == clientID:
			clientPlans = append(clientPlans, entry)
		}
	}

	var groups []OptionGroup
	if len(clientPlans) > 0 {
		groups = append(groups, OptionGroup{
			GroupLabel: ResolvePlanGroupForClientLabel(l.PlanGroupForClient, clientName),
			Options:    clientPlans,
		})
	}
	if len(masterPlans) > 0 {
		groups = append(groups, OptionGroup{
			GroupLabel: l.PlanGroupGeneral,
			Options:    masterPlans,
		})
	}
	return groups
}

// LoadPricePlanOptionGroups is the same shape as LoadPlanOptionGroups but
// keyed off PricePlan.client_id OR the joined PriceSchedule.client_id. The
// schedule join captures rate-card-level scoping where the underlying Plan
// is generic but the operator pinned the rate card to a specific client.
// Used by the subscription edit drawer's PricePlan picker.
func LoadPricePlanOptionGroups(
	ctx context.Context,
	listPricePlans func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error),
	listPriceSchedules func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error),
	clientID, clientName string,
	l Labels,
) []OptionGroup {
	if listPricePlans == nil {
		return nil
	}
	resp, err := listPricePlans(ctx, &priceplanpb.ListPricePlansRequest{})
	if err != nil {
		log.Printf("Failed to load price plans for grouped picker: %v", err)
		return nil
	}

	scheduleClientByID := map[string]string{}
	if listPriceSchedules != nil {
		schedResp, schedErr := listPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
		if schedErr != nil {
			log.Printf("Failed to load price schedules for grouped picker: %v", schedErr)
		} else {
			for _, s := range schedResp.GetData() {
				scheduleClientByID[s.GetId()] = s.GetClientId()
			}
		}
	}

	var clientPP, masterPP []map[string]string
	for _, pp := range resp.GetData() {
		if !pp.GetActive() {
			continue
		}
		label := pp.GetName()
		if label == "" {
			if pl := pp.GetPlan(); pl != nil {
				label = pl.GetName()
			}
			if label == "" {
				label = pp.GetId()
			}
		}
		entry := map[string]string{"Value": pp.GetId(), "Label": label}
		ppClient := pp.GetClientId()
		schedClient := scheduleClientByID[pp.GetPriceScheduleId()]
		switch {
		case clientID != "" && (ppClient == clientID || schedClient == clientID):
			clientPP = append(clientPP, entry)
		case ppClient == "" && schedClient == "":
			masterPP = append(masterPP, entry)
		}
	}

	var groups []OptionGroup
	if len(clientPP) > 0 {
		groups = append(groups, OptionGroup{
			GroupLabel: ResolvePlanGroupForClientLabel(l.PlanGroupForClient, clientName),
			Options:    clientPP,
		})
	}
	if len(masterPP) > 0 {
		groups = append(groups, OptionGroup{
			GroupLabel: l.PlanGroupGeneral,
			Options:    masterPP,
		})
	}
	return groups
}
