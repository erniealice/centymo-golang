// Package block — supplier-contract-price-schedule domain wiring.
//
// Holds wireSupplierContractPriceScheduleModules (the lifted bodies of
// the `if cfg.wantSupplierContractPriceSchedule()` and
// `if cfg.wantSupplierContractPriceScheduleLine()` branches of Block()).
//
// SPS Wave 4 — supplier-side pricing graph (20260427-supplier-commitments).
package block

import (
	"context"

	consumer "github.com/erniealice/espyna-golang/consumer"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	scpspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	suppliercontractpriceschedulemod "github.com/erniealice/centymo-golang/views/supplier_contract_price_schedule"
	suppliercontractpricescheduleinemod "github.com/erniealice/centymo-golang/views/supplier_contract_price_schedule_line"
)

// supplierContractPriceScheduleWiring holds everything wireSupplierContractPriceScheduleModules
// needs from the surrounding Block() scope. More than 6 fields → struct.
// Kept private; never re-exported.
type supplierContractPriceScheduleWiring struct {
	supplierContractPriceScheduleRoutes centymo.SupplierContractPriceScheduleRoutes
	supplierContractPriceScheduleLabels centymo.SupplierContractPriceScheduleLabels
	centymoTableLabels                  types.TableLabels
	uploadFile                          func(context.Context, string, string, []byte, string) error
	listAttachments                     func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment                    func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment                    func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID                     func() string
}

// wireSupplierContractPriceScheduleModules lifts the bodies of the two
// `if cfg.wantSupplierContractPriceScheduleXxx()` branches from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once at the position where
// the SPS Wave 4 wiring used to be.
func wireSupplierContractPriceScheduleModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *consumer.UseCases, w supplierContractPriceScheduleWiring) {
	// SupplierContractPriceSchedule module
	if cfg.wantSupplierContractPriceSchedule() {
		scpsDeps := &suppliercontractpriceschedulemod.ModuleDeps{
			Routes:       w.supplierContractPriceScheduleRoutes,
			Labels:       w.supplierContractPriceScheduleLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
		}
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContractPriceSchedule != nil {
			uc := useCases.Expenditure.SupplierContractPriceSchedule
			if uc.ListSupplierContractPriceSchedules != nil {
				scpsDeps.ListSupplierContractPriceSchedules = uc.ListSupplierContractPriceSchedules.Execute
			}
			if uc.ReadSupplierContractPriceSchedule != nil {
				scpsDeps.ReadSupplierContractPriceSchedule = uc.ReadSupplierContractPriceSchedule.Execute
			}
			if uc.CreateSupplierContractPriceSchedule != nil {
				scpsDeps.CreateSupplierContractPriceSchedule = uc.CreateSupplierContractPriceSchedule.Execute
			}
			if uc.UpdateSupplierContractPriceSchedule != nil {
				scpsDeps.UpdateSupplierContractPriceSchedule = uc.UpdateSupplierContractPriceSchedule.Execute
			}
			if uc.DeleteSupplierContractPriceSchedule != nil {
				scpsDeps.DeleteSupplierContractPriceSchedule = uc.DeleteSupplierContractPriceSchedule.Execute
			}
			// Workflow — wrap Execute() with the closure shapes the view expects.
			if uc.ActivateSupplierContractPriceSchedule != nil {
				activateUC := uc.ActivateSupplierContractPriceSchedule
				scpsDeps.ActivateSupplierContractPriceSchedule = func(fctx context.Context, id string) error {
					userID := consumer.ExtractUserIDFromContext(fctx)
					_, err := activateUC.Execute(fctx, &scpspb.ActivateSupplierContractPriceScheduleRequest{
						SupplierContractPriceScheduleId: id,
						ActivatedBy:                     userID,
					})
					return err
				}
			}
			if uc.SupersedeSupplierContractPriceSchedule != nil {
				supersedeUC := uc.SupersedeSupplierContractPriceSchedule
				scpsDeps.SupersedeSupplierContractPriceSchedule = func(fctx context.Context, id, reason string) error {
					req := &scpspb.SupersedeSupplierContractPriceScheduleRequest{SupplierContractPriceScheduleId: id}
					if reason != "" {
						req.Reason = &reason
					}
					_, err := supersedeUC.Execute(fctx, req)
					return err
				}
			}
		}
		// Schedule lines — list query for the schedule detail's Lines tab.
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContractPriceScheduleLine != nil {
			if uc := useCases.Expenditure.SupplierContractPriceScheduleLine.ListSupplierContractPriceScheduleLines; uc != nil {
				scpsDeps.ListSupplierContractPriceScheduleLines = uc.Execute
			}
		}
		// Parent contract picker for the drawer form + line picker on detail.
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContract != nil &&
			useCases.Expenditure.SupplierContract.ListSupplierContracts != nil {
			scpsDeps.ListSupplierContracts = useCases.Expenditure.SupplierContract.ListSupplierContracts.Execute
		}
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContractLine != nil &&
			useCases.Expenditure.SupplierContractLine.ListSupplierContractLines != nil {
			scpsDeps.ListSupplierContractLines = useCases.Expenditure.SupplierContractLine.ListSupplierContractLines.Execute
		}
		scpsDeps.UploadFile = w.uploadFile
		scpsDeps.ListAttachments = w.listAttachments
		scpsDeps.CreateAttachment = w.createAttachment
		scpsDeps.DeleteAttachment = w.deleteAttachment
		scpsDeps.NewAttachmentID = w.newAttachmentID
		suppliercontractpriceschedulemod.NewModule(scpsDeps).RegisterRoutes(ctx.Routes)
	}

	// SupplierContractPriceScheduleLine module — child rows of SupplierContractPriceSchedule.
	if cfg.wantSupplierContractPriceScheduleLine() {
		scpslDeps := &suppliercontractpricescheduleinemod.ModuleDeps{
			Routes:       w.supplierContractPriceScheduleRoutes,
			Labels:       w.supplierContractPriceScheduleLabels,
			CommonLabels: ctx.Common,
		}
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContractPriceScheduleLine != nil {
			uc := useCases.Expenditure.SupplierContractPriceScheduleLine
			if uc.CreateSupplierContractPriceScheduleLine != nil {
				scpslDeps.CreateSupplierContractPriceScheduleLine = uc.CreateSupplierContractPriceScheduleLine.Execute
			}
			if uc.ReadSupplierContractPriceScheduleLine != nil {
				scpslDeps.ReadSupplierContractPriceScheduleLine = uc.ReadSupplierContractPriceScheduleLine.Execute
			}
			if uc.UpdateSupplierContractPriceScheduleLine != nil {
				scpslDeps.UpdateSupplierContractPriceScheduleLine = uc.UpdateSupplierContractPriceScheduleLine.Execute
			}
			if uc.DeleteSupplierContractPriceScheduleLine != nil {
				scpslDeps.DeleteSupplierContractPriceScheduleLine = uc.DeleteSupplierContractPriceScheduleLine.Execute
			}
		}
		// Parent contract-line picker for the drawer form (line drawer needs a contract-line FK).
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContractLine != nil &&
			useCases.Expenditure.SupplierContractLine.ListSupplierContractLines != nil {
			scpslDeps.ListSupplierContractLines = useCases.Expenditure.SupplierContractLine.ListSupplierContractLines.Execute
		}
		suppliercontractpricescheduleinemod.NewModule(scpslDeps).RegisterRoutes(ctx.Routes)
	}
}
