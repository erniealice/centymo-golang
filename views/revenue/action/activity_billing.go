package action

import (
	"context"
	"fmt"
	"log"
	"strings"

	jobactivitypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_activity"
	revenuelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_line_item"
)

// parseActivityIDs splits a newline-separated string of activity IDs into a
// deduplicated, whitespace-trimmed slice. Empty lines are ignored.
func parseActivityIDs(raw string) []string {
	lines := strings.Split(raw, "\n")
	seen := make(map[string]bool, len(lines))
	ids := make([]string, 0, len(lines))
	for _, line := range lines {
		id := strings.TrimSpace(line)
		if id != "" && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids
}

// autoPopulateFromActivities creates revenue line items from a list of job activity IDs.
// It attempts to read each activity via the ReadJobActivity dep; if unavailable it falls
// back to creating a stub line item referencing the activity ID.
// All errors are logged and silently ignored — this is best-effort enrichment.
func autoPopulateFromActivities(ctx context.Context, deps *Deps, revenueID string, activityIDs []string) {
	if deps.CreateRevenueLineItem == nil {
		return
	}

	for _, activityID := range activityIDs {
		if deps.ReadJobActivity != nil {
			actResp, err := deps.ReadJobActivity(ctx, &jobactivitypb.ReadJobActivityRequest{
				Data: &jobactivitypb.JobActivity{Id: activityID},
			})
			if err != nil {
				log.Printf("autoPopulateFromActivities: failed to read activity %s: %v — creating stub", activityID, err)
				createActivityStubLineItem(ctx, deps, revenueID, activityID)
				continue
			}
			if actResp.GetData() == nil || len(actResp.GetData()) == 0 {
				log.Printf("autoPopulateFromActivities: activity %s not found — creating stub", activityID)
				createActivityStubLineItem(ctx, deps, revenueID, activityID)
				continue
			}

			activity := actResp.GetData()[0]
			createActivityLineItem(ctx, deps, revenueID, activity)
		} else {
			// ReadJobActivity not wired — fall back to stub line item
			createActivityStubLineItem(ctx, deps, revenueID, activityID)
		}
	}
}

// createActivityLineItem creates a revenue line item populated from a resolved job activity.
func createActivityLineItem(ctx context.Context, deps *Deps, revenueID string, activity *jobactivitypb.JobActivity) {
	activityID := activity.GetId()

	// Description: use activity's description, fall back to "Activity: {id}"
	desc := activity.GetDescription()
	if desc == "" {
		desc = fmt.Sprintf("Activity: %s", activityID)
	}

	// Quantity: use activity's quantity; default to 1 if zero or negative
	qty := activity.GetQuantity()
	if qty <= 0 {
		qty = 1
	}

	// Pricing: prefer bill_rate/bill_amount; fall back to unit_cost/total_cost
	var unitPrice, totalPrice int64
	var costPrice *int64
	if activity.BillRate != nil {
		unitPrice = activity.GetBillRate()
	} else {
		unitPrice = activity.GetUnitCost()
	}
	if activity.BillAmount != nil {
		totalPrice = activity.GetBillAmount()
	} else {
		totalPrice = activity.GetTotalCost()
	}
	if uc := activity.GetUnitCost(); uc > 0 {
		costPrice = &uc
	}

	lineItem := &revenuelineitempb.RevenueLineItem{
		RevenueId:     revenueID,
		JobActivityId: &activityID,
		Description:   desc,
		Quantity:      qty,
		UnitPrice:     unitPrice,
		TotalPrice:    totalPrice,
		LineItemType:  "item",
	}
	if costPrice != nil {
		lineItem.CostPrice = costPrice
	}

	_, err := deps.CreateRevenueLineItem(ctx, &revenuelineitempb.CreateRevenueLineItemRequest{
		Data: lineItem,
	})
	if err != nil {
		log.Printf("autoPopulateFromActivities: failed to create line item for activity %s: %v", activityID, err)
	}
}

// createActivityStubLineItem creates a minimal line item referencing an activity by ID only.
// Used when ReadJobActivity is unavailable or the activity cannot be read.
func createActivityStubLineItem(ctx context.Context, deps *Deps, revenueID, activityID string) {
	desc := fmt.Sprintf("Activity: %s", activityID)
	_, err := deps.CreateRevenueLineItem(ctx, &revenuelineitempb.CreateRevenueLineItemRequest{
		Data: &revenuelineitempb.RevenueLineItem{
			RevenueId:     revenueID,
			JobActivityId: &activityID,
			Description:   desc,
			Quantity:      1.0,
			UnitPrice:     0,
			TotalPrice:    0,
			LineItemType:  "item",
		},
	})
	if err != nil {
		log.Printf("autoPopulateFromActivities: failed to create stub line item for activity %s: %v", activityID, err)
	}
}
