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
func wireSupplierContractPriceScheduleModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *UseCases, w supplierContractPriceScheduleWiring) {
	// SupplierContractPriceSchedule module
	if cfg.wantSupplierContractPriceSchedule() {
		scpsDeps := &suppliercontractpriceschedulemod.ModuleDeps{
			Routes:       w.supplierContractPriceScheduleRoutes,
			Labels:       w.supplierContractPriceScheduleLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
		}
		scpsDeps.ListSupplierContractPriceSchedules = useCases.SupplierContract.ListSupplierContractPriceSchedules
		scpsDeps.ReadSupplierContractPriceSchedule = useCases.SupplierContract.ReadSupplierContractPriceSchedule
		scpsDeps.CreateSupplierContractPriceSchedule = useCases.SupplierContract.CreateSupplierContractPriceSchedule
		scpsDeps.UpdateSupplierContractPriceSchedule = useCases.SupplierContract.UpdateSupplierContractPriceSchedule
		scpsDeps.DeleteSupplierContractPriceSchedule = useCases.SupplierContract.DeleteSupplierContractPriceSchedule
		// Workflow — wrap with the closure shapes the view expects.
		if useCases.SupplierContract.ActivateSupplierContractPriceSchedule != nil {
			activateUC := useCases.SupplierContract.ActivateSupplierContractPriceSchedule
			scpsDeps.ActivateSupplierContractPriceSchedule = func(fctx context.Context, id string) error {
				userID := useCases.ExtractUserID(fctx)
				_, err := activateUC(fctx, &scpspb.ActivateSupplierContractPriceScheduleRequest{
					SupplierContractPriceScheduleId: id,
					ActivatedBy:                     userID,
				})
				return err
			}
		}
		if useCases.SupplierContract.SupersedeSupplierContractPriceSchedule != nil {
			supersedeUC := useCases.SupplierContract.SupersedeSupplierContractPriceSchedule
			scpsDeps.SupersedeSupplierContractPriceSchedule = func(fctx context.Context, id, reason string) error {
				req := &scpspb.SupersedeSupplierContractPriceScheduleRequest{SupplierContractPriceScheduleId: id}
				if reason != "" {
					req.Reason = &reason
				}
				_, err := supersedeUC(fctx, req)
				return err
			}
		}
		// Schedule lines — list query for the schedule detail's Lines tab.
		scpsDeps.ListSupplierContractPriceScheduleLines = useCases.SupplierContract.ListSupplierContractPriceScheduleLines
		// Parent contract picker for the drawer form + line picker on detail.
		scpsDeps.ListSupplierContracts = useCases.SupplierContract.ListSupplierContracts
		scpsDeps.ListSupplierContractLines = useCases.SupplierContract.ListSupplierContractLines
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
		scpslDeps.CreateSupplierContractPriceScheduleLine = useCases.SupplierContract.CreateSupplierContractPriceScheduleLine
		scpslDeps.ReadSupplierContractPriceScheduleLine = useCases.SupplierContract.ReadSupplierContractPriceScheduleLine
		scpslDeps.UpdateSupplierContractPriceScheduleLine = useCases.SupplierContract.UpdateSupplierContractPriceScheduleLine
		scpslDeps.DeleteSupplierContractPriceScheduleLine = useCases.SupplierContract.DeleteSupplierContractPriceScheduleLine
		// Parent contract-line picker for the drawer form (line drawer needs a contract-line FK).
		scpslDeps.ListSupplierContractLines = useCases.SupplierContract.ListSupplierContractLines
		suppliercontractpricescheduleinemod.NewModule(scpslDeps).RegisterRoutes(ctx.Routes)
	}
}
