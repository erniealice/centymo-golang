// 20260517-advance-cash-events Plan B Phase 4 — bridge helpers translating
// block-level AdvanceSettleInput / AdvanceRefundInput / AdvanceCancelInput
// shapes to the view-level AdvanceSettleViewInput / AdvanceRefundViewInput /
// AdvanceCancelViewInput shapes the per-package view module deps speak.
//
// Both sides (treasury_collection + treasury_disbursement) share these
// helpers because the input/output shapes are identical and the only
// difference is which closure on UseCases the service-admin adapter binds.
package block

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
)

// bridgeSettleAdvance translates the block UseCases-level
// SettleUnscheduledAdvance closure (if any) into a view-level closure the
// per-package collection/disbursement modules can call directly without
// importing the block sub-package.
func bridgeSettleAdvance(uc func(ctx context.Context, in AdvanceSettleInput) (*AdvanceSettleOutput, error)) func(ctx context.Context, in centymo.AdvanceSettleViewInput) (*centymo.AdvanceSettleViewOutput, error) {
	if uc == nil {
		return nil
	}
	return func(ctx context.Context, in centymo.AdvanceSettleViewInput) (*centymo.AdvanceSettleViewOutput, error) {
		out, err := uc(ctx, AdvanceSettleInput{
			AdvanceID:       in.AdvanceID,
			Amount:          in.Amount,
			TargetAccountID: in.TargetAccountID,
			Reason:          in.Reason,
		})
		if err != nil || out == nil {
			return nil, err
		}
		return &centymo.AdvanceSettleViewOutput{
			NewRemainingAmount:  out.NewRemainingAmount,
			NewRecognizedAmount: out.NewRecognizedAmount,
			NewStatus:           out.NewStatus,
		}, nil
	}
}

// bridgeRefundAdvance mirrors bridgeSettleAdvance for the Refund closure.
func bridgeRefundAdvance(uc func(ctx context.Context, in AdvanceRefundInput) (*AdvanceRefundOutput, error)) func(ctx context.Context, in centymo.AdvanceRefundViewInput) (*centymo.AdvanceRefundViewOutput, error) {
	if uc == nil {
		return nil
	}
	return func(ctx context.Context, in centymo.AdvanceRefundViewInput) (*centymo.AdvanceRefundViewOutput, error) {
		out, err := uc(ctx, AdvanceRefundInput{
			AdvanceID:          in.AdvanceID,
			Amount:             in.Amount,
			RefundMethod:       in.RefundMethod,
			DestinationAccount: in.DestinationAccount,
			Reason:             in.Reason,
		})
		if err != nil || out == nil {
			return nil, err
		}
		return &centymo.AdvanceRefundViewOutput{
			NewRemainingAmount: out.NewRemainingAmount,
			NewStatus:          out.NewStatus,
		}, nil
	}
}

// bridgeCancelAdvance mirrors bridgeSettleAdvance for the Cancel closure.
func bridgeCancelAdvance(uc func(ctx context.Context, in AdvanceCancelInput) (*AdvanceCancelOutput, error)) func(ctx context.Context, in centymo.AdvanceCancelViewInput) (*centymo.AdvanceCancelViewOutput, error) {
	if uc == nil {
		return nil
	}
	return func(ctx context.Context, in centymo.AdvanceCancelViewInput) (*centymo.AdvanceCancelViewOutput, error) {
		out, err := uc(ctx, AdvanceCancelInput{
			AdvanceID: in.AdvanceID,
			Reason:    in.Reason,
		})
		if err != nil || out == nil {
			return nil, err
		}
		return &centymo.AdvanceCancelViewOutput{NewStatus: out.NewStatus}, nil
	}
}
