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

	consumer "github.com/erniealice/espyna-golang/consumer"

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
func wireAccruedExpenseModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *consumer.UseCases, w accruedExpenseWiring) {
	// AccruedExpense module
	if cfg.wantAccruedExpense() {
		aeDeps := &accruedexpensemod.ModuleDeps{
			Routes:       w.accruedExpenseRoutes,
			Labels:       w.accruedExpenseLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
		}
		if useCases.Expenditure != nil && useCases.Expenditure.AccruedExpense != nil {
			uc := useCases.Expenditure.AccruedExpense
			if uc.ListAccruedExpenses != nil {
				aeDeps.ListAccruedExpenses = uc.ListAccruedExpenses.Execute
			}
			if uc.ReadAccruedExpense != nil {
				aeDeps.ReadAccruedExpense = uc.ReadAccruedExpense.Execute
			}
			if uc.CreateAccruedExpense != nil {
				aeDeps.CreateAccruedExpense = uc.CreateAccruedExpense.Execute
			}
			if uc.UpdateAccruedExpense != nil {
				aeDeps.UpdateAccruedExpense = uc.UpdateAccruedExpense.Execute
			}
			if uc.DeleteAccruedExpense != nil {
				aeDeps.DeleteAccruedExpense = uc.DeleteAccruedExpense.Execute
			}
			if uc.SettleAccrual != nil {
				settleUC := uc.SettleAccrual
				aeDeps.SettleAccrual = func(fctx context.Context, req *accruedexpensepb.SettleAccrualRequest) error {
					_, err := settleUC.SettleAccrual(fctx, req)
					return err
				}
			}
			if uc.ReverseAccrual != nil {
				reverseUC := uc.ReverseAccrual
				aeDeps.ReverseAccrual = func(fctx context.Context, id, reason string) error {
					req := &accruedexpensepb.ReverseAccrualRequest{AccruedExpenseId: id}
					if reason != "" {
						req.Reason = &reason
					}
					_, err := reverseUC.Execute(fctx, req)
					return err
				}
			}
			if uc.AccrueFromContract != nil {
				afcUC := uc.AccrueFromContract
				aeDeps.AccrueFromContract = func(fctx context.Context, req *accruedexpensepb.AccrueFromContractRequest) (*accruedexpensepb.AccrueFromContractResponse, error) {
					return afcUC.Execute(fctx, req)
				}
			}
		}
		// Inline child — settlements.
		if useCases.Expenditure != nil && useCases.Expenditure.AccruedExpenseSettlement != nil {
			if uc := useCases.Expenditure.AccruedExpenseSettlement.ListAccruedExpenseSettlements; uc != nil {
				aeDeps.ListAccruedExpenseSettlements = uc.Execute
			}
		}
		// Dropdowns for the manual-create drawer + filter pickers.
		if useCases.Entity != nil && useCases.Entity.Supplier != nil &&
			useCases.Entity.Supplier.ListSuppliers != nil {
			aeDeps.ListSuppliers = useCases.Entity.Supplier.ListSuppliers.Execute
		}
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContract != nil &&
			useCases.Expenditure.SupplierContract.ListSupplierContracts != nil {
			aeDeps.ListSupplierContracts = useCases.Expenditure.SupplierContract.ListSupplierContracts.Execute
		}
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
		if useCases.Expenditure != nil && useCases.Expenditure.AccruedExpenseSettlement != nil {
			uc := useCases.Expenditure.AccruedExpenseSettlement
			if uc.CreateAccruedExpenseSettlement != nil {
				aesDeps.CreateAccruedExpenseSettlement = uc.CreateAccruedExpenseSettlement.Execute
			}
			if uc.ReadAccruedExpenseSettlement != nil {
				aesDeps.ReadAccruedExpenseSettlement = uc.ReadAccruedExpenseSettlement.Execute
			}
			if uc.UpdateAccruedExpenseSettlement != nil {
				aesDeps.UpdateAccruedExpenseSettlement = uc.UpdateAccruedExpenseSettlement.Execute
			}
			if uc.DeleteAccruedExpenseSettlement != nil {
				aesDeps.DeleteAccruedExpenseSettlement = uc.DeleteAccruedExpenseSettlement.Execute
			}
		}
		// Settling-Expenditure picker for the settlement drawer form.
		if useCases.Expenditure != nil && useCases.Expenditure.Expenditure != nil &&
			useCases.Expenditure.Expenditure.ListExpenditures != nil {
			aesDeps.ListExpenditures = useCases.Expenditure.Expenditure.ListExpenditures.Execute
		}
		accruedexpensesettlementmod.NewModule(aesDeps).RegisterRoutes(ctx.Routes)
	}
}
