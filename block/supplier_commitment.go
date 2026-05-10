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

	consumer "github.com/erniealice/espyna-golang/consumer"

	procurementrequestpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	procurementmod "github.com/erniealice/centymo-golang/views/procurement"
	procurementrequestmod "github.com/erniealice/centymo-golang/views/procurement_request"
	procurementrequestlinemod "github.com/erniealice/centymo-golang/views/procurement_request_line"
	suppliercontractmod "github.com/erniealice/centymo-golang/views/supplier_contract"
	suppliercontractlinemod "github.com/erniealice/centymo-golang/views/supplier_contract_line"
)

// supplierCommitmentWiring holds everything wireSupplierCommitmentModules
// needs from the surrounding Block() scope. More than 6 fields → struct.
// Kept private; never re-exported.
type supplierCommitmentWiring struct {
	supplierContractRoutes            centymo.SupplierContractRoutes
	supplierContractLabels            centymo.SupplierContractLabels
	supplierContractPriceScheduleRoutes centymo.SupplierContractPriceScheduleRoutes
	procurementRequestRoutes          centymo.ProcurementRequestRoutes
	procurementRequestLabels          centymo.ProcurementRequestLabels
	procurementRoutes                 centymo.ProcurementRoutes
	procurementLabels                 centymo.ProcurementLabels
	centymoTableLabels                types.TableLabels
	uploadFile                        func(context.Context, string, string, []byte, string) error
	listAttachments                   func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment                  func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment                  func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID                   func() string
}

// wireSupplierCommitmentModules lifts the bodies of the five
// `if cfg.wantSupplierXxx()` / `if cfg.wantProcurementXxx()` branches
// from Block(). Behaviour-preserving: same construction order, same
// registration order, same callbacks. block.go calls this exactly once.
func wireSupplierCommitmentModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *consumer.UseCases, w supplierCommitmentWiring) {
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
	// consumer.ExtractUserIDFromContext so the centymo views package stays
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
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContract != nil {
			uc := useCases.Expenditure.SupplierContract
			if uc.CreateSupplierContract != nil {
				scDeps.CreateSupplierContract = uc.CreateSupplierContract.Execute
			}
			if uc.ReadSupplierContract != nil {
				scDeps.ReadSupplierContract = uc.ReadSupplierContract.Execute
			}
			if uc.UpdateSupplierContract != nil {
				scDeps.UpdateSupplierContract = uc.UpdateSupplierContract.Execute
			}
			if uc.DeleteSupplierContract != nil {
				scDeps.DeleteSupplierContract = uc.DeleteSupplierContract.Execute
			}
			if uc.ListSupplierContracts != nil {
				scDeps.ListSupplierContracts = uc.ListSupplierContracts.Execute
			}
			// Workflow actions — wrap Execute with closures that source
			// the approver/user identity from ctx (set by the session
			// middleware in the composition layer).
			if uc.ApproveSupplierContract != nil {
				approveUC := uc.ApproveSupplierContract
				scDeps.ApproveSupplierContract = func(fctx context.Context, id string) error {
					userID := consumer.ExtractUserIDFromContext(fctx)
					_, err := approveUC.Execute(fctx, &suppliercontractpb.ApproveSupplierContractRequest{
						SupplierContractId: id,
						ApprovedBy:         userID,
					})
					return err
				}
			}
			if uc.TerminateSupplierContract != nil {
				terminateUC := uc.TerminateSupplierContract
				scDeps.TerminateSupplierContract = func(fctx context.Context, id, reason string) error {
					req := &suppliercontractpb.TerminateSupplierContractRequest{
						SupplierContractId: id,
					}
					if reason != "" {
						req.Reason = &reason
					}
					_, err := terminateUC.Execute(fctx, req)
					return err
				}
			}
		}
		// Lines query for the Lines tab on the contract detail page.
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContractLine != nil {
			if uc := useCases.Expenditure.SupplierContractLine.ListSupplierContractLines; uc != nil {
				scDeps.ListSupplierContractLines = uc.Execute
			}
		}
		// Suppliers dropdown.
		if useCases.Entity != nil && useCases.Entity.Supplier != nil &&
			useCases.Entity.Supplier.ListSuppliers != nil {
			scDeps.ListSuppliers = useCases.Entity.Supplier.ListSuppliers.Execute
		}
		// Linked POs and Linked Expenditures tabs on the detail page.
		if useCases.Expenditure != nil && useCases.Expenditure.PurchaseOrder != nil &&
			useCases.Expenditure.PurchaseOrder.ListPurchaseOrders != nil {
			scDeps.ListPurchaseOrders = useCases.Expenditure.PurchaseOrder.ListPurchaseOrders.Execute
		}
		if useCases.Expenditure != nil && useCases.Expenditure.Expenditure != nil &&
			useCases.Expenditure.Expenditure.ListExpenditures != nil {
			scDeps.ListExpenditures = useCases.Expenditure.Expenditure.ListExpenditures.Execute
		}
		// SPS Wave 4 — Price Schedules tab on the contract detail page.
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContractPriceSchedule != nil {
			if uc := useCases.Expenditure.SupplierContractPriceSchedule.ListSupplierContractPriceSchedules; uc != nil {
				scDeps.ListSupplierContractPriceSchedules = uc.Execute
			}
		}
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
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContractLine != nil {
			uc := useCases.Expenditure.SupplierContractLine
			if uc.CreateSupplierContractLine != nil {
				sclDeps.CreateSupplierContractLine = uc.CreateSupplierContractLine.Execute
			}
			if uc.ReadSupplierContractLine != nil {
				sclDeps.ReadSupplierContractLine = uc.ReadSupplierContractLine.Execute
			}
			if uc.UpdateSupplierContractLine != nil {
				sclDeps.UpdateSupplierContractLine = uc.UpdateSupplierContractLine.Execute
			}
			if uc.DeleteSupplierContractLine != nil {
				sclDeps.DeleteSupplierContractLine = uc.DeleteSupplierContractLine.Execute
			}
		}
		// Product picker for the line drawer form.
		if useCases.Product != nil && useCases.Product.Product != nil &&
			useCases.Product.Product.ListProducts != nil {
			sclDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
		}
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
		if useCases.Expenditure != nil && useCases.Expenditure.ProcurementRequest != nil {
			uc := useCases.Expenditure.ProcurementRequest
			if uc.CreateProcurementRequest != nil {
				prDeps.CreateProcurementRequest = uc.CreateProcurementRequest.Execute
			}
			if uc.ReadProcurementRequest != nil {
				prDeps.ReadProcurementRequest = uc.ReadProcurementRequest.Execute
			}
			if uc.UpdateProcurementRequest != nil {
				prDeps.UpdateProcurementRequest = uc.UpdateProcurementRequest.Execute
			}
			if uc.DeleteProcurementRequest != nil {
				prDeps.DeleteProcurementRequest = uc.DeleteProcurementRequest.Execute
			}
			if uc.ListProcurementRequests != nil {
				prDeps.ListProcurementRequests = uc.ListProcurementRequests.Execute
			}
			// Workflow action closures — sourced ApprovedBy from ctx.
			if uc.SubmitProcurementRequest != nil {
				submitUC := uc.SubmitProcurementRequest
				prDeps.SubmitProcurementRequest = func(fctx context.Context, id string) error {
					_, err := submitUC.Execute(fctx, &procurementrequestpb.SubmitProcurementRequestRequest{
						ProcurementRequestId: id,
					})
					return err
				}
			}
			if uc.ApproveProcurementRequest != nil {
				approveUC := uc.ApproveProcurementRequest
				prDeps.ApproveProcurementRequest = func(fctx context.Context, id string) error {
					userID := consumer.ExtractUserIDFromContext(fctx)
					_, err := approveUC.Execute(fctx, &procurementrequestpb.ApproveProcurementRequestRequest{
						ProcurementRequestId: id,
						ApprovedBy:           userID,
					})
					return err
				}
			}
			if uc.RejectProcurementRequest != nil {
				rejectUC := uc.RejectProcurementRequest
				prDeps.RejectProcurementRequest = func(fctx context.Context, id, reason string) error {
					req := &procurementrequestpb.RejectProcurementRequestRequest{
						ProcurementRequestId: id,
					}
					if reason != "" {
						req.RejectionReason = &reason
					}
					_, err := rejectUC.Execute(fctx, req)
					return err
				}
			}
			if uc.SpawnPurchaseOrder != nil {
				spawnUC := uc.SpawnPurchaseOrder
				prDeps.SpawnPurchaseOrder = func(fctx context.Context, id string) (string, error) {
					resp, err := spawnUC.Execute(fctx, &procurementrequestpb.SpawnPurchaseOrderRequest{
						ProcurementRequestId: id,
					})
					if err != nil {
						return "", err
					}
					return resp.GetPurchaseOrderId(), nil
				}
			}
		}
		// Lines query for the Lines tab on the request detail page.
		if useCases.Expenditure != nil && useCases.Expenditure.ProcurementRequestLine != nil {
			if uc := useCases.Expenditure.ProcurementRequestLine.ListProcurementRequestLines; uc != nil {
				prDeps.ListProcurementRequestLines = uc.Execute
			}
		}
		// Suppliers dropdown (nullable for RFQ flow).
		if useCases.Entity != nil && useCases.Entity.Supplier != nil &&
			useCases.Entity.Supplier.ListSuppliers != nil {
			prDeps.ListSuppliers = useCases.Entity.Supplier.ListSuppliers.Execute
		}
		// Spawned POs tab.
		if useCases.Expenditure != nil && useCases.Expenditure.PurchaseOrder != nil &&
			useCases.Expenditure.PurchaseOrder.ListPurchaseOrders != nil {
			prDeps.ListPurchaseOrders = useCases.Expenditure.PurchaseOrder.ListPurchaseOrders.Execute
		}
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
		if useCases.Expenditure != nil && useCases.Expenditure.ProcurementRequestLine != nil {
			uc := useCases.Expenditure.ProcurementRequestLine
			if uc.CreateProcurementRequestLine != nil {
				prlDeps.CreateProcurementRequestLine = uc.CreateProcurementRequestLine.Execute
			}
			if uc.ReadProcurementRequestLine != nil {
				prlDeps.ReadProcurementRequestLine = uc.ReadProcurementRequestLine.Execute
			}
			if uc.UpdateProcurementRequestLine != nil {
				prlDeps.UpdateProcurementRequestLine = uc.UpdateProcurementRequestLine.Execute
			}
			if uc.DeleteProcurementRequestLine != nil {
				prlDeps.DeleteProcurementRequestLine = uc.DeleteProcurementRequestLine.Execute
			}
		}
		// Product picker for the line drawer form.
		if useCases.Product != nil && useCases.Product.Product != nil &&
			useCases.Product.Product.ListProducts != nil {
			prlDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
		}
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
		if useCases.Expenditure != nil && useCases.Expenditure.SupplierContract != nil &&
			useCases.Expenditure.SupplierContract.ListSupplierContracts != nil {
			procDeps.ListSupplierContracts = useCases.Expenditure.SupplierContract.ListSupplierContracts.Execute
		}
		if useCases.Expenditure != nil && useCases.Expenditure.ProcurementRequest != nil &&
			useCases.Expenditure.ProcurementRequest.ListProcurementRequests != nil {
			procDeps.ListProcurementRequests = useCases.Expenditure.ProcurementRequest.ListProcurementRequests.Execute
		}
		if useCases.Expenditure != nil && useCases.Expenditure.Expenditure != nil &&
			useCases.Expenditure.Expenditure.ListExpenditures != nil {
			procDeps.ListExpenditures = useCases.Expenditure.Expenditure.ListExpenditures.Execute
		}
		procurementmod.NewModule(procDeps).RegisterRoutes(ctx.Routes)
	}
}
