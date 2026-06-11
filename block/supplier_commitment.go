// Package block — supplier-commitment domain wiring.
//
// Holds wireSupplierCommitmentModules (the lifted body of the five
// `if cfg.wantSupplierXxx()` branches for: SupplierContract,
// SupplierContractLine, ProcurementRequest, ProcurementRequestLine,
// and Procurement). These five modules form the supplier-side
// commitments wave (20260427-supplier-commitments P3a/P3b).
//
// All use-case threading is nil-safe — when the espyna composition layer
// didn't initialize a use case, the corresponding view falls back to its
// empty/disabled state instead of panicking.
package block

import (
	"context"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	procurementrequestpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	expendituredomain "github.com/erniealice/centymo-golang/domain/expenditure"
	procurementrequestmod "github.com/erniealice/centymo-golang/domain/expenditure/views/procurement_request"
	procurementrequestlinemod "github.com/erniealice/centymo-golang/domain/expenditure/views/procurement_request_line"
	suppliercontractmod "github.com/erniealice/centymo-golang/domain/expenditure/views/supplier_contract"
	suppliercontractlinemod "github.com/erniealice/centymo-golang/domain/expenditure/views/supplier_contract_line"
	procurementmod "github.com/erniealice/centymo-golang/views/procurement"
)

// supplierCommitmentWiring holds everything wireSupplierCommitmentModules
// needs from the surrounding Block() scope. More than 6 fields → struct.
// Kept private; never re-exported.
type supplierCommitmentWiring struct {
	supplierContractRoutes              expendituredomain.SupplierContractRoutes
	supplierContractLabels              expendituredomain.SupplierContractLabels
	supplierContractPriceScheduleRoutes expendituredomain.SupplierContractPriceScheduleRoutes
	procurementRequestRoutes            expendituredomain.ProcurementRequestRoutes
	procurementRequestLabels            expendituredomain.ProcurementRequestLabels
	procurementRoutes                   centymo.ProcurementRoutes
	procurementLabels                   centymo.ProcurementLabels
	centymoTableLabels                  types.TableLabels
	uploadFile                          func(context.Context, string, string, []byte, string) error
	listAttachments                     func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment                    func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment                    func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID                     func() string
}

// wireSupplierCommitmentModules lifts the bodies of the five
// `if cfg.wantSupplierXxx()` / `if cfg.wantProcurementXxx()` branches
// from Block(). Behaviour-preserving: same construction order, same
// registration order, same callbacks. block.go calls this exactly once.
func wireSupplierCommitmentModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *UseCases, w supplierCommitmentWiring) {
	// =====================================================================
	// 20260427-supplier-commitments — five new modules (P3a + P3b)
	//
	// Mirrors the expendituremod pattern: construct ModuleDeps, plumb each
	// available use case through, register routes.
	// All use-case threading is nil-safe — when the espyna composition
	// layer didn't initialize a use case, the corresponding view falls back
	// to its empty/disabled state instead of panicking.
	//
	// Workflow action closures (Submit/Approve/Reject/SpawnPO and
	// Approve/Terminate) source ApprovedBy from the request context via
	// useCases.ExtractUserID so the centymo views package stays
	// free of espyna ctx imports.
	// =====================================================================

	// SupplierContract module
	if cfg.wantSupplierContract() {
		scDeps := &suppliercontractmod.ModuleDeps{
			Routes:       w.supplierContractRoutes,
			Labels:       w.supplierContractLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
		}
		scDeps.CreateSupplierContract = useCases.SupplierContract.CreateSupplierContract
		scDeps.ReadSupplierContract = useCases.SupplierContract.ReadSupplierContract
		scDeps.UpdateSupplierContract = useCases.SupplierContract.UpdateSupplierContract
		scDeps.DeleteSupplierContract = useCases.SupplierContract.DeleteSupplierContract
		scDeps.ListSupplierContracts = useCases.SupplierContract.ListSupplierContracts
		// Workflow actions — wrap with closures that source
		// the approver/user identity from ctx (set by the session
		// middleware in the composition layer).
		if useCases.SupplierContract.ApproveSupplierContract != nil {
			approveUC := useCases.SupplierContract.ApproveSupplierContract
			scDeps.ApproveSupplierContract = func(fctx context.Context, id string) error {
				userID := useCases.ExtractUserID(fctx)
				_, err := approveUC(fctx, &suppliercontractpb.ApproveSupplierContractRequest{
					SupplierContractId: id,
					ApprovedBy:         userID,
				})
				return err
			}
		}
		if useCases.SupplierContract.TerminateSupplierContract != nil {
			terminateUC := useCases.SupplierContract.TerminateSupplierContract
			scDeps.TerminateSupplierContract = func(fctx context.Context, id, reason string) error {
				req := &suppliercontractpb.TerminateSupplierContractRequest{
					SupplierContractId: id,
				}
				if reason != "" {
					req.Reason = &reason
				}
				_, err := terminateUC(fctx, req)
				return err
			}
		}
		// Lines query for the Lines tab on the contract detail page.
		scDeps.ListSupplierContractLines = useCases.SupplierContract.ListSupplierContractLines
		// Suppliers dropdown.
		scDeps.ListSuppliers = useCases.Entity.Supplier.ListSuppliers
		// Linked POs and Linked Expenditures tabs on the detail page.
		scDeps.ListPurchaseOrders = useCases.Expenditure.ListPurchaseOrders
		scDeps.ListExpenditures = useCases.Expenditure.ListExpenditures
		// SPS Wave 4 — Price Schedules tab on the contract detail page.
		scDeps.ListSupplierContractPriceSchedules = useCases.SupplierContract.ListSupplierContractPriceSchedules
		scDeps.PriceScheduleListURL = w.supplierContractPriceScheduleRoutes.ListURL
		scDeps.PriceScheduleDetailURL = w.supplierContractPriceScheduleRoutes.DetailURL
		scDeps.PriceScheduleAddURL = w.supplierContractPriceScheduleRoutes.AddURL
		scDeps.UploadFile = w.uploadFile
		scDeps.ListAttachments = w.listAttachments
		scDeps.CreateAttachment = w.createAttachment
		scDeps.DeleteAttachment = w.deleteAttachment
		scDeps.NewAttachmentID = w.newAttachmentID
		suppliercontractmod.NewModule(scDeps).RegisterRoutes(ctx.Routes)
	}

	// SupplierContractLine module — child rows of SupplierContract.
	if cfg.wantSupplierContractLine() {
		sclDeps := &suppliercontractlinemod.ModuleDeps{
			Routes:       w.supplierContractRoutes,
			Labels:       w.supplierContractLabels,
			CommonLabels: ctx.Common,
		}
		sclDeps.CreateSupplierContractLine = useCases.SupplierContract.CreateSupplierContractLine
		sclDeps.ReadSupplierContractLine = useCases.SupplierContract.ReadSupplierContractLine
		sclDeps.UpdateSupplierContractLine = useCases.SupplierContract.UpdateSupplierContractLine
		sclDeps.DeleteSupplierContractLine = useCases.SupplierContract.DeleteSupplierContractLine
		// Product picker for the line drawer form.
		sclDeps.ListProducts = useCases.Product.ListProducts
		suppliercontractlinemod.NewModule(sclDeps).RegisterRoutes(ctx.Routes)
	}

	// ProcurementRequest module
	if cfg.wantProcurementRequest() {
		prDeps := &procurementrequestmod.ModuleDeps{
			Routes:       w.procurementRequestRoutes,
			Labels:       w.procurementRequestLabels,
			CommonLabels: ctx.Common,
			TableLabels:  w.centymoTableLabels,
		}
		prDeps.CreateProcurementRequest = useCases.SupplierContract.CreateProcurementRequest
		prDeps.ReadProcurementRequest = useCases.SupplierContract.ReadProcurementRequest
		prDeps.UpdateProcurementRequest = useCases.SupplierContract.UpdateProcurementRequest
		prDeps.DeleteProcurementRequest = useCases.SupplierContract.DeleteProcurementRequest
		prDeps.ListProcurementRequests = useCases.SupplierContract.ListProcurementRequests
		// Workflow action closures — source ApprovedBy from ctx.
		if useCases.SupplierContract.SubmitProcurementRequest != nil {
			submitUC := useCases.SupplierContract.SubmitProcurementRequest
			prDeps.SubmitProcurementRequest = func(fctx context.Context, id string) error {
				_, err := submitUC(fctx, &procurementrequestpb.SubmitProcurementRequestRequest{
					ProcurementRequestId: id,
				})
				return err
			}
		}
		if useCases.SupplierContract.ApproveProcurementRequest != nil {
			approveUC := useCases.SupplierContract.ApproveProcurementRequest
			prDeps.ApproveProcurementRequest = func(fctx context.Context, id string) error {
				userID := useCases.ExtractUserID(fctx)
				_, err := approveUC(fctx, &procurementrequestpb.ApproveProcurementRequestRequest{
					ProcurementRequestId: id,
					ApprovedBy:           userID,
				})
				return err
			}
		}
		if useCases.SupplierContract.RejectProcurementRequest != nil {
			rejectUC := useCases.SupplierContract.RejectProcurementRequest
			prDeps.RejectProcurementRequest = func(fctx context.Context, id, reason string) error {
				req := &procurementrequestpb.RejectProcurementRequestRequest{
					ProcurementRequestId: id,
				}
				if reason != "" {
					req.RejectionReason = &reason
				}
				_, err := rejectUC(fctx, req)
				return err
			}
		}
		if useCases.SupplierContract.SpawnProcurementRequestPO != nil {
			spawnUC := useCases.SupplierContract.SpawnProcurementRequestPO
			prDeps.SpawnPurchaseOrder = func(fctx context.Context, id string) (string, error) {
				resp, err := spawnUC(fctx, &procurementrequestpb.SpawnPurchaseOrderRequest{
					ProcurementRequestId: id,
				})
				if err != nil {
					return "", err
				}
				return resp.GetPurchaseOrderId(), nil
			}
		}
		// Lines query for the Lines tab on the request detail page.
		prDeps.ListProcurementRequestLines = useCases.SupplierContract.ListProcurementRequestLines
		// Suppliers dropdown (nullable for RFQ flow).
		prDeps.ListSuppliers = useCases.Entity.Supplier.ListSuppliers
		// Spawned POs tab.
		prDeps.ListPurchaseOrders = useCases.Expenditure.ListPurchaseOrders
		prDeps.UploadFile = w.uploadFile
		prDeps.ListAttachments = w.listAttachments
		prDeps.CreateAttachment = w.createAttachment
		prDeps.DeleteAttachment = w.deleteAttachment
		prDeps.NewAttachmentID = w.newAttachmentID
		procurementrequestmod.NewModule(prDeps).RegisterRoutes(ctx.Routes)
	}

	// ProcurementRequestLine module — child rows of ProcurementRequest.
	if cfg.wantProcurementRequestLine() {
		prlDeps := &procurementrequestlinemod.ModuleDeps{
			Routes:       w.procurementRequestRoutes,
			Labels:       w.procurementRequestLabels,
			CommonLabels: ctx.Common,
		}
		prlDeps.CreateProcurementRequestLine = useCases.SupplierContract.CreateProcurementRequestLine
		prlDeps.ReadProcurementRequestLine = useCases.SupplierContract.ReadProcurementRequestLine
		prlDeps.UpdateProcurementRequestLine = useCases.SupplierContract.UpdateProcurementRequestLine
		prlDeps.DeleteProcurementRequestLine = useCases.SupplierContract.DeleteProcurementRequestLine
		// Product picker for the line drawer form.
		prlDeps.ListProducts = useCases.Product.ListProducts
		procurementrequestlinemod.NewModule(prlDeps).RegisterRoutes(ctx.Routes)
	}

	// Procurement Operations composition app (read-only — no proto entity).
	// Nil-safe: missing list closures render empty states in each view.
	if cfg.wantProcurement() {
		procDeps := &procurementmod.ModuleDeps{
			Routes:       w.procurementRoutes,
			Labels:       w.procurementLabels,
			CommonLabels: ctx.Common,
		}
		procDeps.ListSupplierContracts = useCases.SupplierContract.ListSupplierContracts
		procDeps.ListProcurementRequests = useCases.SupplierContract.ListProcurementRequests
		procDeps.ListExpenditures = useCases.Expenditure.ListExpenditures
		procurementmod.NewModule(procDeps).RegisterRoutes(ctx.Routes)
	}
}
