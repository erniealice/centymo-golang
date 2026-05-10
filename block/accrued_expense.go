// Package block — accrued-expense domain wiring.
//
// Holds wireAccruedExpenseModules (the lifted bodies of
// the `if cfg.wantAccruedExpense()` and
// `if cfg.wantAccruedExpenseSettlement()` branches of Block()).
//
// SPS Wave 4 — accrued-expense master + settlement (20260427-supplier-commitments).
package block

import (
	"context"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	accruedexpensemod "github.com/erniealice/centymo-golang/views/accrued_expense"
	accruedexpensesettlementmod "github.com/erniealice/centymo-golang/views/accrued_expense_settlement"
)

// accruedExpenseWiring holds everything wireAccruedExpenseModules
// needs from the surrounding Block() scope. More than 6 fields → struct.
// Kept private; never re-exported.
type accruedExpenseWiring struct {
	accruedExpenseRoutes centymo.AccruedExpenseRoutes
	accruedExpenseLabels centymo.AccruedExpenseLabels
	centymoTableLabels   types.TableLabels
	uploadFile           func(context.Context, string, string, []byte, string) error
	listAttachments      func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment     func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment     func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID      func() string
}

// wireAccruedExpenseModules lifts the bodies of the two
// `if cfg.wantAccruedExpenseXxx()` branches from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once at the position where
// the SPS Wave 4 AccruedExpense wiring used to be.
func wireAccruedExpenseModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *UseCases, w accruedExpenseWiring) {
	// AccruedExpense module
	if cfg.wantAccruedExpense() {
		aeDeps := &accruedexpensemod.ModuleDeps{
			Routes:       w.accruedExpenseRoutes,
			Labels:       w.accruedExpenseLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
		}
		aeDeps.ListAccruedExpenses = useCases.Expenditure.ListAccruedExpenses
		aeDeps.ReadAccruedExpense = useCases.Expenditure.ReadAccruedExpense
		aeDeps.CreateAccruedExpense = useCases.Expenditure.CreateAccruedExpense
		aeDeps.UpdateAccruedExpense = useCases.Expenditure.UpdateAccruedExpense
		aeDeps.DeleteAccruedExpense = useCases.Expenditure.DeleteAccruedExpense
		if useCases.Expenditure.SettleAccrual != nil {
			settleUC := useCases.Expenditure.SettleAccrual
			aeDeps.SettleAccrual = func(fctx context.Context, req *accruedexpensepb.SettleAccrualRequest) error {
				_, err := settleUC(fctx, req)
				return err
			}
		}
		if useCases.Expenditure.ReverseAccrual != nil {
			reverseUC := useCases.Expenditure.ReverseAccrual
			aeDeps.ReverseAccrual = func(fctx context.Context, id, reason string) error {
				req := &accruedexpensepb.ReverseAccrualRequest{AccruedExpenseId: id}
				if reason != "" {
					req.Reason = &reason
				}
				_, err := reverseUC(fctx, req)
				return err
			}
		}
		if useCases.Expenditure.AccrueFromContract != nil {
			afcUC := useCases.Expenditure.AccrueFromContract
			aeDeps.AccrueFromContract = func(fctx context.Context, req *accruedexpensepb.AccrueFromContractRequest) (*accruedexpensepb.AccrueFromContractResponse, error) {
				return afcUC(fctx, req)
			}
		}
		// Inline child — settlements.
		aeDeps.ListAccruedExpenseSettlements = useCases.Expenditure.ListAccruedExpenseSettlements
		// Dropdowns for the manual-create drawer + filter pickers.
		aeDeps.ListSuppliers = useCases.Entity.Supplier.ListSuppliers
		aeDeps.ListSupplierContracts = useCases.SupplierContract.ListSupplierContracts
		aeDeps.UploadFile = w.uploadFile
		aeDeps.ListAttachments = w.listAttachments
		aeDeps.CreateAttachment = w.createAttachment
		aeDeps.DeleteAttachment = w.deleteAttachment
		aeDeps.NewAttachmentID = w.newAttachmentID
		accruedexpensemod.NewModule(aeDeps).RegisterRoutes(ctx.Routes)
	}

	// AccruedExpenseSettlement module — inline child of AccruedExpense (shares parent routes).
	if cfg.wantAccruedExpenseSettlement() {
		aesDeps := &accruedexpensesettlementmod.ModuleDeps{
			Routes:       w.accruedExpenseRoutes,
			Labels:       w.accruedExpenseLabels,
			CommonLabels: ctx.Common,
		}
		aesDeps.CreateAccruedExpenseSettlement = useCases.Expenditure.CreateAccruedExpenseSettlement
		aesDeps.ReadAccruedExpenseSettlement = useCases.Expenditure.ReadAccruedExpenseSettlement
		aesDeps.UpdateAccruedExpenseSettlement = useCases.Expenditure.UpdateAccruedExpenseSettlement
		aesDeps.DeleteAccruedExpenseSettlement = useCases.Expenditure.DeleteAccruedExpenseSettlement
		// Settling-Expenditure picker for the settlement drawer form.
		aesDeps.ListExpenditures = useCases.Expenditure.ListExpenditures
		accruedexpensesettlementmod.NewModule(aesDeps).RegisterRoutes(ctx.Routes)
	}
}
