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

	consumer "github.com/erniealice/espyna-golang/consumer"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	expenserecognitionmod "github.com/erniealice/centymo-golang/views/expense_recognition"
	expenserecognitionlinemod "github.com/erniealice/centymo-golang/views/expense_recognition_line"
)

// expenseRecognitionWiring holds everything wireExpenseRecognitionModules
// needs from the surrounding Block() scope. More than 6 fields → struct.
// Kept private; never re-exported.
type expenseRecognitionWiring struct {
	expenseRecognitionRoutes centymo.ExpenseRecognitionRoutes
	expenseRecognitionLabels centymo.ExpenseRecognitionLabels
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
func wireExpenseRecognitionModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *consumer.UseCases, w expenseRecognitionWiring) {
	// ExpenseRecognition module — no Add/Edit drawer (created BY use case).
	if cfg.wantExpenseRecognition() {
		erDeps := &expenserecognitionmod.ModuleDeps{
			Routes:       w.expenseRecognitionRoutes,
			Labels:       w.expenseRecognitionLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
		}
		if useCases.Expenditure != nil && useCases.Expenditure.ExpenseRecognition != nil {
			uc := useCases.Expenditure.ExpenseRecognition
			if uc.ListExpenseRecognitions != nil {
				erDeps.ListExpenseRecognitions = uc.ListExpenseRecognitions.Execute
			}
			if uc.ReadExpenseRecognition != nil {
				erDeps.ReadExpenseRecognition = uc.ReadExpenseRecognition.Execute
			}
			if uc.DeleteExpenseRecognition != nil {
				erDeps.DeleteExpenseRecognition = uc.DeleteExpenseRecognition.Execute
			}
			if uc.ReverseExpenseRecognition != nil {
				reverseUC := uc.ReverseExpenseRecognition
				erDeps.ReverseExpenseRecognition = func(fctx context.Context, id, reason string) error {
					req := &expenserecognitionpb.ReverseExpenseRecognitionRequest{ExpenseRecognitionId: id}
					if reason != "" {
						req.Reason = &reason
					}
					_, err := reverseUC.Execute(fctx, req)
					return err
				}
			}
			if uc.RecognizeFromExpenditure != nil {
				rfeUC := uc.RecognizeFromExpenditure
				erDeps.RecognizeFromExpenditure = func(fctx context.Context, req *expenserecognitionpb.RecognizeFromExpenditureRequest) (*expenserecognitionpb.RecognizeFromExpenditureResponse, error) {
					return rfeUC.Execute(fctx, req)
				}
			}
			if uc.RecognizeFromContract != nil {
				rfcUC := uc.RecognizeFromContract
				erDeps.RecognizeFromContract = func(fctx context.Context, req *expenserecognitionpb.RecognizeFromContractRequest) (*expenserecognitionpb.RecognizeFromContractResponse, error) {
					return rfcUC.Execute(fctx, req)
				}
			}
		}
		// Inline child — recognition lines.
		if useCases.Expenditure != nil && useCases.Expenditure.ExpenseRecognitionLine != nil {
			if uc := useCases.Expenditure.ExpenseRecognitionLine.ListExpenseRecognitionLines; uc != nil {
				erDeps.ListExpenseRecognitionLines = uc.Execute
			}
		}
		erDeps.UploadFile = w.uploadFile
		erDeps.ListAttachments = w.listAttachments
		erDeps.CreateAttachment = w.createAttachment
		erDeps.DeleteAttachment = w.deleteAttachment
		erDeps.NewAttachmentID = w.newAttachmentID
		expenserecognitionmod.NewModule(erDeps).RegisterRoutes(ctx.Routes)
	}

	// ExpenseRecognitionLine module — inline child of ExpenseRecognition.
	if cfg.wantExpenseRecognitionLine() {
		erlDeps := &expenserecognitionlinemod.ModuleDeps{
			Routes:       w.expenseRecognitionRoutes,
			Labels:       w.expenseRecognitionLabels,
			CommonLabels: ctx.Common,
		}
		if useCases.Expenditure != nil && useCases.Expenditure.ExpenseRecognitionLine != nil {
			uc := useCases.Expenditure.ExpenseRecognitionLine
			if uc.CreateExpenseRecognitionLine != nil {
				erlDeps.CreateExpenseRecognitionLine = uc.CreateExpenseRecognitionLine.Execute
			}
			if uc.ReadExpenseRecognitionLine != nil {
				erlDeps.ReadExpenseRecognitionLine = uc.ReadExpenseRecognitionLine.Execute
			}
			if uc.UpdateExpenseRecognitionLine != nil {
				erlDeps.UpdateExpenseRecognitionLine = uc.UpdateExpenseRecognitionLine.Execute
			}
			if uc.DeleteExpenseRecognitionLine != nil {
				erlDeps.DeleteExpenseRecognitionLine = uc.DeleteExpenseRecognitionLine.Execute
			}
		}
		expenserecognitionlinemod.NewModule(erlDeps).RegisterRoutes(ctx.Routes)
	}
}
