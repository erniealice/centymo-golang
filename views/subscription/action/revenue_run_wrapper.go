package action

// revenue_run_wrapper.go provides the thin shim that plugs the per-subscription
// Invoice Run drawer into the subscription action aggregator. Block.go calls
// NewRevenueRunAction(subActionDeps) to register the drawer route.

import (
	"context"

	"github.com/erniealice/pyeza-golang/view"

	revenuerundeps "github.com/erniealice/centymo-golang/views/subscription/revenue_run"
)

// NewRevenueRunAction is the shim for block.go.
// Delegates to revenue_run.NewAction using a sub-set of action.Deps, converting
// the action-package type shapes into the sub-package's local shapes.
func NewRevenueRunAction(deps *Deps) view.View {
	innerDeps := &revenuerundeps.Deps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
	}

	// Wire ListRevenueRunCandidates when present.
	if deps.ListRevenueRunCandidates != nil {
		outerList := deps.ListRevenueRunCandidates
		innerDeps.ListRevenueRunCandidates = func(ctx context.Context, scope revenuerundeps.RevenueRunScope) ([]revenuerundeps.RevenueRunCandidate, string, error) {
			outerCandidates, nextCursor, err := outerList(ctx, RevenueRunScopeAction{
				WorkspaceID:    scope.WorkspaceID,
				ClientID:       scope.ClientID,
				SubscriptionID: scope.SubscriptionID,
				AsOfDate:       scope.AsOfDate,
				Cursor:         scope.Cursor,
				Limit:          scope.Limit,
			})
			if err != nil {
				return nil, "", err
			}
			out := make([]revenuerundeps.RevenueRunCandidate, 0, len(outerCandidates))
			for _, c := range outerCandidates {
				out = append(out, revenuerundeps.RevenueRunCandidate{
					SubscriptionID:    c.SubscriptionID,
					SubscriptionName:  c.SubscriptionName,
					ClientID:          c.ClientID,
					ClientName:        c.ClientName,
					PlanName:          c.PlanName,
					BillingCycleLabel: c.BillingCycleLabel,
					Currency:          c.Currency,
					PeriodStart:       c.PeriodStart,
					PeriodEnd:         c.PeriodEnd,
					PeriodLabel:       c.PeriodLabel,
					PeriodMarker:      c.PeriodMarker,
					Amount:            c.Amount,
					AmountDisplay:     c.AmountDisplay,
					LineItemCount:     c.LineItemCount,
					Eligible:          c.Eligible,
					BlockerReason:     c.BlockerReason,
				})
			}
			return out, nextCursor, nil
		}
	}

	// Wire GenerateRevenueRun when present.
	if deps.GenerateRevenueRun != nil {
		outerGenerate := deps.GenerateRevenueRun
		innerDeps.GenerateRevenueRun = func(ctx context.Context, scope revenuerundeps.RevenueRunScope, sels revenuerundeps.RevenueRunSelections) (*revenuerundeps.RevenueRunResult, error) {
			outerSels := RevenueRunSelectionsAction{
				FilterToken: sels.FilterToken,
			}
			for _, s := range sels.ExplicitList {
				outerSels.ExplicitList = append(outerSels.ExplicitList, SelectedRevenueRunCandidateAction{
					SubscriptionID: s.SubscriptionID,
					PeriodStart:    s.PeriodStart,
					PeriodEnd:      s.PeriodEnd,
					PeriodMarker:   s.PeriodMarker,
				})
			}
			result, err := outerGenerate(ctx, RevenueRunScopeAction{
				WorkspaceID:    scope.WorkspaceID,
				ClientID:       scope.ClientID,
				SubscriptionID: scope.SubscriptionID,
				AsOfDate:       scope.AsOfDate,
				Cursor:         scope.Cursor,
				Limit:          scope.Limit,
			}, outerSels)
			if err != nil {
				return nil, err
			}
			if result == nil {
				return nil, nil
			}
			return &revenuerundeps.RevenueRunResult{
				RunID:   result.RunID,
				Status:  result.Status,
				Created: result.Created,
				Skipped: result.Skipped,
				Errored: result.Errored,
			}, nil
		}
	}

	return revenuerundeps.NewAction(innerDeps)
}
