// Package block — expense-recognition domain wiring.
//
// Holds wireExpenseRecognitionModules (the lifted bodies of
// the `if cfg.wantExpenseRecognition()` and
// `if cfg.wantExpenseRecognitionLine()` branches of Block()).
//
// SPS Wave 4 — expense-recognition + line (20260427-supplier-commitments).
package block

import (
	"context"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	expendituredomain "github.com/erniealice/centymo-golang/domain/expenditure"
	expenserecognitionmodmodule "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition/module"
	expenserecognitionlinemod "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_line"
)

// expenseRecognitionWiring holds everything wireExpenseRecognitionModules
// needs from the surrounding Block() scope. More than 6 fields → struct.
// Kept private; never re-exported.
type expenseRecognitionWiring struct {
	expenseRecognitionRoutes expendituredomain.ExpenseRecognitionRoutes
	expenseRecognitionLabels expendituredomain.ExpenseRecognitionLabels
	centymoTableLabels       types.TableLabels
	uploadFile               func(context.Context, string, string, []byte, string) error
	listAttachments          func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment         func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment         func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID          func() string
}

// wireExpenseRecognitionModules lifts the bodies of the two
// `if cfg.wantExpenseRecognitionXxx()` branches from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once at the position where
// the SPS Wave 4 ExpRec wiring used to be.
func wireExpenseRecognitionModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *UseCases, w expenseRecognitionWiring) {
	// ExpenseRecognition module — no Add/Edit drawer (created BY use case).
	if cfg.wantExpenseRecognition() {
		erDeps := &expenserecognitionmodmodule.ModuleDeps{
			Routes:       w.expenseRecognitionRoutes,
			Labels:       w.expenseRecognitionLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
		}
		erDeps.ListExpenseRecognitions = useCases.Expenditure.ListExpenseRecognitions
		erDeps.ReadExpenseRecognition = useCases.Expenditure.ReadExpenseRecognition
		erDeps.DeleteExpenseRecognition = useCases.Expenditure.DeleteExpenseRecognition
		if useCases.Expenditure.ReverseExpenseRecognition != nil {
			reverseUC := useCases.Expenditure.ReverseExpenseRecognition
			erDeps.ReverseExpenseRecognition = func(fctx context.Context, id, reason string) error {
				req := &expenserecognitionpb.ReverseExpenseRecognitionRequest{ExpenseRecognitionId: id}
				if reason != "" {
					req.Reason = &reason
				}
				_, err := reverseUC(fctx, req)
				return err
			}
		}
		if useCases.Expenditure.RecognizeFromExpenditure != nil {
			rfeUC := useCases.Expenditure.RecognizeFromExpenditure
			erDeps.RecognizeFromExpenditure = func(fctx context.Context, req *expenserecognitionpb.RecognizeFromExpenditureRequest) (*expenserecognitionpb.RecognizeFromExpenditureResponse, error) {
				return rfeUC(fctx, req)
			}
		}
		if useCases.Expenditure.RecognizeFromContract != nil {
			rfcUC := useCases.Expenditure.RecognizeFromContract
			erDeps.RecognizeFromContract = func(fctx context.Context, req *expenserecognitionpb.RecognizeFromContractRequest) (*expenserecognitionpb.RecognizeFromContractResponse, error) {
				return rfcUC(fctx, req)
			}
		}
		// Inline child — recognition lines.
		erDeps.ListExpenseRecognitionLines = useCases.Expenditure.ListExpenseRecognitionLines
		erDeps.UploadFile = w.uploadFile
		erDeps.ListAttachments = w.listAttachments
		erDeps.CreateAttachment = w.createAttachment
		erDeps.DeleteAttachment = w.deleteAttachment
		erDeps.NewAttachmentID = w.newAttachmentID
		expenserecognitionmodmodule.NewModule(erDeps).RegisterRoutes(ctx.Routes)
	}

	// ExpenseRecognitionLine module — inline child of ExpenseRecognition.
	if cfg.wantExpenseRecognitionLine() {
		erlDeps := &expenserecognitionlinemod.ModuleDeps{
			Routes:       w.expenseRecognitionRoutes,
			Labels:       w.expenseRecognitionLabels,
			CommonLabels: ctx.Common,
		}
		erlDeps.CreateExpenseRecognitionLine = useCases.Expenditure.CreateExpenseRecognitionLine
		erlDeps.ReadExpenseRecognitionLine = useCases.Expenditure.ReadExpenseRecognitionLine
		erlDeps.UpdateExpenseRecognitionLine = useCases.Expenditure.UpdateExpenseRecognitionLine
		erlDeps.DeleteExpenseRecognitionLine = useCases.Expenditure.DeleteExpenseRecognitionLine
		expenserecognitionlinemod.NewModule(erlDeps).RegisterRoutes(ctx.Routes)
	}
}
