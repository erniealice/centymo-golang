package block

import (
	"context"
	"log"

	expenseboard "github.com/erniealice/centymo-golang/domain/expenditure/expenditure/expense_dashboard"
	purchaseboard "github.com/erniealice/centymo-golang/domain/expenditure/expenditure/purchase_dashboard"
	productdashboardview "github.com/erniealice/centymo-golang/domain/product/product/dashboard"
	cashdashboardview "github.com/erniealice/centymo-golang/domain/treasury/collection/dashboard"

	"github.com/erniealice/espyna-golang/consumer"
	consumerapp "github.com/erniealice/espyna-golang/consumer/app"
	"github.com/erniealice/espyna-golang/reference"
	advancekindpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common/advance_kind"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	paymenttermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/payment_term"
	expenserecognitionrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition_run"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
	revenuerunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_run"
	treasurycollectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
	treasurydisbursementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement"
	expendituredashpb "github.com/erniealice/esqyma/pkg/schema/v1/service/dashboard/expenditure"
	productdashpb "github.com/erniealice/esqyma/pkg/schema/v1/service/dashboard/product"
	treasurydashpb "github.com/erniealice/esqyma/pkg/schema/v1/service/dashboard/treasury"
)

// centymoEngineBlock returns a consumerapp.AppOption that registers all centymo
// domain modules via the compose engine (replaces legacy centymoBlock).
func EngineBlock() consumerapp.AppOption {
	return func(ctx *consumerapp.AppContext) error {
		uc, err := consumerapp.RequireUseCases(ctx, "centymoEngineBlock")
		if err != nil {
			return err
		}
		adapted := buildCentymoUseCases(uc, ctx.DB)

		infra := &Infra{}
		infra.UploadFile, _ = ctx.UploadFile.(func(context.Context, string, string, []byte, string) error)
		infra.DownloadFile, _ = ctx.DownloadFile.(func(context.Context, string, string) ([]byte, error))
		infra.ListAttachments, _ = ctx.ListAttachments.(func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error))
		infra.CreateAttachment, _ = ctx.CreateAttachment.(func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error))
		infra.ReadAttachment, _ = ctx.ReadAttachment.(func(context.Context, *attachmentpb.ReadAttachmentRequest) (*attachmentpb.ReadAttachmentResponse, error))
		infra.DeleteAttachment, _ = ctx.DeleteAttachment.(func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error))
		infra.NewAttachmentID, _ = ctx.NewAttachmentID.(func() string)
		infra.UploadImage, _ = ctx.UploadImage.(func(context.Context, string, string, []byte, string) error)
		infra.UploadTemplate, _ = ctx.UploadTemplate.(func(context.Context, string, string, []byte, string) error)
		infra.SendEmail, _ = ctx.SendEmail.(func(context.Context, []string, string, string, string, string, []byte) error)
		infra.GenerateDoc, _ = ctx.GenerateDoc.(func([]byte, map[string]any) ([]byte, error))
		if ctx.RefChecker != nil {
			if rc, ok := ctx.RefChecker.(reference.Checker); ok {
				infra.RefChecker = rc
			}
		}

		units := AllUnits(adapted, infra)
		return consumerapp.AssembleEngineBlock("centymo", units, ctx)
	}
}

// ---------------------------------------------------------------------------
// centymo adapter
// ---------------------------------------------------------------------------

// centymoDBOps is the capability-narrow operations surface the three
// ops-backed centymo closures (SetActive / DeleteProductVariantOption /
// payment_term dropdown) need from the concrete database adapter passed via
// ctx.DB (a *consumer.DatabaseAdapter). It is the typed successor to the
// deleted centymo `DataSource` duck — service-admin asserts ctx.DB to it once
// and binds only the three closures that genuinely need generic-collection ops.
// Signatures MUST match *consumer.DatabaseAdapter exactly.
// 20260612-datasource-typed-path W7.
type centymoDBOps interface {
	Update(ctx context.Context, collection string, id string, data map[string]any) (map[string]any, error)
	HardDelete(ctx context.Context, collection string, id string) error
	ListSimple(ctx context.Context, collection string) ([]map[string]any, error)
}

// buildCentymoUseCases maps espyna's *consumer.UseCases to centymo block's
// typed shape. All sub-group wiring is nil-safe — if espyna hasn't wired a
// sub-domain the corresponding closure stays nil and centymo renders an empty
// state or skips that feature.
//
// db is the concrete database adapter (ctx.DB, a *consumer.DatabaseAdapter).
// It backs the three generic-collection closures (SetActive,
// DeleteProductVariantOption, payment_term dropdown) that the deleted centymo
// DataSource duck used to serve. When the assertion to centymoDBOps fails
// (mock build, nil DB) those three closures are left nil — centymo degrades
// gracefully (toggle no-ops, variant-option delete unavailable, payment-term
// dropdown empty). 20260612-datasource-typed-path W7.
func buildCentymoUseCases(uc *consumer.UseCases, db any) *UseCases {
	result := &UseCases{
		ExtractUserID: consumer.ExtractUserIDFromContext,
	}

	// Assert ctx.DB to the capability-narrow ops surface once. ok==false (mock
	// build / nil DB) → the three ops-backed closures below stay nil.
	ops, opsOK := db.(centymoDBOps)
	if !opsOK {
		log.Printf("buildCentymoUseCases: ctx.DB does not satisfy centymoDBOps — SetActive/DeleteProductVariantOption/payment_term dropdown left unwired (active toggles no-op, variant-option delete unavailable, payment-term dropdown empty)")
	}

	// -- Ops-backed closures (W7) ------------------------------------------------
	// The three generic-collection closures the deleted centymo DataSource duck
	// used to serve. Bound only when ctx.DB satisfied centymoDBOps.
	// 20260612-datasource-typed-path W7.
	if opsOK {
		// SetActive — sets ONLY the `active` boolean on the named collection.
		// Typed proto3 Update can't clear `active` (omits false bools), so this
		// goes through the generic ops.Update with an explicit {"active": active}
		// map. Mirrors the deleted duck's Update(collection, id, {"active": ...}).
		result.SetActive = func(ctx context.Context, collection string, id string, active bool) error {
			_, err := ops.Update(ctx, collection, id, map[string]any{"active": active})
			return err
		}

		// DeleteProductVariantOption — junction-row hard delete keyed only by id.
		// Replaces the duck's HardDelete("product_variant_option", id); returns an
		// empty typed response to satisfy the block's proto closure signature.
		result.Product.DeleteProductVariantOption = func(ctx context.Context, req *productvariantoptionpb.DeleteProductVariantOptionRequest) (*productvariantoptionpb.DeleteProductVariantOptionResponse, error) {
			id := req.GetData().GetId()
			err := ops.HardDelete(ctx, "product_variant_option", id)
			return &productvariantoptionpb.DeleteProductVariantOptionResponse{}, err
		}

		// PaymentTerm.ListPaymentTerms — no espyna ListPaymentTerms use case is
		// exposed on the consumer aggregate (EntityUseCases has no PaymentTerm
		// field), so this falls back to the generic ops.ListSimple("payment_term")
		// + manual proto construction, behaviour-equivalent to the old duck's
		// ListSimple("payment_term") dropdown. The block scopes/filters the rows
		// by entity_scope itself, so an unfiltered list of all terms is returned.
		// Only the fields the block reads (id, name, net_days, entity_scope, active)
		// are populated.
		result.Entity.PaymentTerm.ListPaymentTerms = func(ctx context.Context, _ *paymenttermpb.ListPaymentTermsRequest) (*paymenttermpb.ListPaymentTermsResponse, error) {
			rows, err := ops.ListSimple(ctx, "payment_term")
			if err != nil {
				return nil, err
			}
			out := &paymenttermpb.ListPaymentTermsResponse{Success: true}
			for _, row := range rows {
				if row == nil {
					continue
				}
				pt := &paymenttermpb.PaymentTerm{}
				if v, ok := row["id"].(string); ok {
					pt.Id = v
				}
				if v, ok := row["name"].(string); ok {
					pt.Name = v
				}
				if v, ok := row["entity_scope"].(string); ok {
					pt.EntityScope = v
				}
				if v, ok := paymentTermInt32(row["net_days"]); ok {
					pt.NetDays = v
				}
				if v, ok := row["active"].(bool); ok {
					pt.Active = v
				}
				out.Data = append(out.Data, pt)
			}
			return out, nil
		}
	}

	// -- Common ------------------------------------------------------------------
	if uc.Common != nil {
		if uc.Common.Attribute != nil {
			result.Common.ListAttributes = uc.Common.Attribute.ListAttributes.Execute
			result.Common.ReadAttribute = uc.Common.Attribute.ReadAttribute.Execute
		}
		if uc.Common.Category != nil {
			result.Common.ListCategories = uc.Common.Category.ListCategories.Execute
		}
	}

	// -- Entity ------------------------------------------------------------------
	if uc.Entity != nil {
		if uc.Entity.Client != nil {
			result.Entity.Client.ListClients = uc.Entity.Client.ListClients.Execute
			result.Entity.Client.ReadClient = uc.Entity.Client.ReadClient.Execute
			result.Entity.Client.SearchClientsByName = uc.Entity.Client.SearchClientsByName.Execute
		}
		if uc.Entity.Location != nil {
			result.Entity.Location.ListLocations = uc.Entity.Location.ListLocations.Execute
		}
		if uc.Entity.Supplier != nil {
			result.Entity.Supplier.ListSuppliers = uc.Entity.Supplier.ListSuppliers.Execute
		}
		if uc.Entity.Workspace != nil {
			result.Entity.Workspace.ReadWorkspace = uc.Entity.Workspace.ReadWorkspace.Execute
		}
	}

	// -- Inventory ---------------------------------------------------------------
	if uc.Inventory != nil {
		if uc.Inventory.InventoryItem != nil {
			result.Inventory.ListInventoryItems = uc.Inventory.InventoryItem.ListInventoryItems.Execute
			result.Inventory.CreateInventoryItem = uc.Inventory.InventoryItem.CreateInventoryItem.Execute
			result.Inventory.ReadInventoryItem = uc.Inventory.InventoryItem.ReadInventoryItem.Execute
			result.Inventory.UpdateInventoryItem = uc.Inventory.InventoryItem.UpdateInventoryItem.Execute
			result.Inventory.DeleteInventoryItem = uc.Inventory.InventoryItem.DeleteInventoryItem.Execute
		}
		if uc.Inventory.InventorySerial != nil {
			result.Inventory.ListInventorySerials = uc.Inventory.InventorySerial.ListInventorySerials.Execute
			result.Inventory.CreateInventorySerial = uc.Inventory.InventorySerial.CreateInventorySerial.Execute
			result.Inventory.ReadInventorySerial = uc.Inventory.InventorySerial.ReadInventorySerial.Execute
			result.Inventory.UpdateInventorySerial = uc.Inventory.InventorySerial.UpdateInventorySerial.Execute
			result.Inventory.DeleteInventorySerial = uc.Inventory.InventorySerial.DeleteInventorySerial.Execute
		}
		if uc.Inventory.InventoryTransaction != nil {
			result.Inventory.ListInventoryTransactions = uc.Inventory.InventoryTransaction.ListInventoryTransactions.Execute
			result.Inventory.CreateInventoryTransaction = uc.Inventory.InventoryTransaction.CreateInventoryTransaction.Execute
			result.Inventory.GetInventoryMovementsListPageData = uc.Inventory.InventoryTransaction.GetInventoryMovementsListPageData.Execute
		}
		if uc.Inventory.InventoryDepreciation != nil {
			result.Inventory.ListInventoryDepreciations = uc.Inventory.InventoryDepreciation.ListInventoryDepreciations.Execute
			result.Inventory.CreateInventoryDepreciation = uc.Inventory.InventoryDepreciation.CreateInventoryDepreciation.Execute
			result.Inventory.ReadInventoryDepreciation = uc.Inventory.InventoryDepreciation.ReadInventoryDepreciation.Execute
			result.Inventory.UpdateInventoryDepreciation = uc.Inventory.InventoryDepreciation.UpdateInventoryDepreciation.Execute
		}
		if uc.Inventory.InventorySerialHistory != nil {
			result.Inventory.CreateInventorySerialHistory = uc.Inventory.InventorySerialHistory.CreateInventorySerialHistory.Execute
		}
	}

	// -- Product -----------------------------------------------------------------
	if uc.Product != nil {
		if uc.Product.Product != nil {
			result.Product.ListProducts = uc.Product.Product.ListProducts.Execute
			result.Product.ReadProduct = uc.Product.Product.ReadProduct.Execute
			result.Product.CreateProduct = uc.Product.Product.CreateProduct.Execute
			result.Product.UpdateProduct = uc.Product.Product.UpdateProduct.Execute
			result.Product.DeleteProduct = uc.Product.Product.DeleteProduct.Execute
		}
		if uc.Product.ProductVariant != nil {
			result.Product.ListProductVariants = uc.Product.ProductVariant.ListProductVariants.Execute
			result.Product.ReadProductVariant = uc.Product.ProductVariant.ReadProductVariant.Execute
			result.Product.CreateProductVariant = uc.Product.ProductVariant.CreateProductVariant.Execute
			result.Product.UpdateProductVariant = uc.Product.ProductVariant.UpdateProductVariant.Execute
			result.Product.DeleteProductVariant = uc.Product.ProductVariant.DeleteProductVariant.Execute
		}
		if uc.Product.ProductVariantOption != nil {
			result.Product.ListProductVariantOptions = uc.Product.ProductVariantOption.ListProductVariantOptions.Execute
			result.Product.CreateProductVariantOption = uc.Product.ProductVariantOption.CreateProductVariantOption.Execute
		}
		if uc.Product.ProductOption != nil {
			result.Product.ListProductOptions = uc.Product.ProductOption.ListProductOptions.Execute
			result.Product.ReadProductOption = uc.Product.ProductOption.ReadProductOption.Execute
			result.Product.CreateProductOption = uc.Product.ProductOption.CreateProductOption.Execute
			result.Product.UpdateProductOption = uc.Product.ProductOption.UpdateProductOption.Execute
			result.Product.DeleteProductOption = uc.Product.ProductOption.DeleteProductOption.Execute
		}
		if uc.Product.ProductOptionValue != nil {
			result.Product.ListProductOptionValues = uc.Product.ProductOptionValue.ListProductOptionValues.Execute
			result.Product.ReadProductOptionValue = uc.Product.ProductOptionValue.ReadProductOptionValue.Execute
			result.Product.CreateProductOptionValue = uc.Product.ProductOptionValue.CreateProductOptionValue.Execute
			result.Product.UpdateProductOptionValue = uc.Product.ProductOptionValue.UpdateProductOptionValue.Execute
			result.Product.DeleteProductOptionValue = uc.Product.ProductOptionValue.DeleteProductOptionValue.Execute
		}
		if uc.Product.ProductAttribute != nil {
			result.Product.ListProductAttributes = uc.Product.ProductAttribute.ListProductAttributes.Execute
			result.Product.CreateProductAttribute = uc.Product.ProductAttribute.CreateProductAttribute.Execute
			result.Product.DeleteProductAttribute = uc.Product.ProductAttribute.DeleteProductAttribute.Execute
		}
		if uc.Product.Line != nil {
			result.Product.ListLines = uc.Product.Line.ListLines.Execute
			result.Product.ReadLine = uc.Product.Line.ReadLine.Execute
			result.Product.CreateLine = uc.Product.Line.CreateLine.Execute
			result.Product.UpdateLine = uc.Product.Line.UpdateLine.Execute
			result.Product.DeleteLine = uc.Product.Line.DeleteLine.Execute
		}
		if uc.Product.ProductLine != nil {
			result.Product.ListProductLines = uc.Product.ProductLine.ListProductLines.Execute
			result.Product.ReadProductLine = uc.Product.ProductLine.ReadProductLine.Execute
			result.Product.CreateProductLine = uc.Product.ProductLine.CreateProductLine.Execute
			result.Product.UpdateProductLine = uc.Product.ProductLine.UpdateProductLine.Execute
			result.Product.DeleteProductLine = uc.Product.ProductLine.DeleteProductLine.Execute
		}
		if uc.Product.ProductVariantImage != nil {
			result.Product.ListProductVariantImages = uc.Product.ProductVariantImage.ListProductVariantImages.Execute
			result.Product.CreateProductVariantImage = uc.Product.ProductVariantImage.CreateProductVariantImage.Execute
			result.Product.DeleteProductVariantImage = uc.Product.ProductVariantImage.DeleteProductVariantImage.Execute
		}
		if uc.Product.ProductPlan != nil {
			result.Product.ListProductPlans = uc.Product.ProductPlan.ListProductPlans.Execute
			result.Product.ReadProductPlan = uc.Product.ProductPlan.ReadProductPlan.Execute
			result.Product.CreateProductPlan = uc.Product.ProductPlan.CreateProductPlan.Execute
			result.Product.UpdateProductPlan = uc.Product.ProductPlan.UpdateProductPlan.Execute
			result.Product.DeleteProductPlan = uc.Product.ProductPlan.DeleteProductPlan.Execute
		}
		if uc.Product.PriceList != nil {
			result.Product.FindApplicablePriceList = uc.Product.PriceList.FindApplicablePriceList.Execute
			result.Product.ListPriceLists = uc.Product.PriceList.ListPriceLists.Execute
			result.Product.ReadPriceList = uc.Product.PriceList.ReadPriceList.Execute
			result.Product.CreatePriceList = uc.Product.PriceList.CreatePriceList.Execute
			result.Product.UpdatePriceList = uc.Product.PriceList.UpdatePriceList.Execute
			result.Product.DeletePriceList = uc.Product.PriceList.DeletePriceList.Execute
		}
		if uc.Product.PriceProduct != nil {
			result.Product.ListPriceProducts = uc.Product.PriceProduct.ListPriceProducts.Execute
			result.Product.CreatePriceProduct = uc.Product.PriceProduct.CreatePriceProduct.Execute
			result.Product.DeletePriceProduct = uc.Product.PriceProduct.DeletePriceProduct.Execute
		}
		if uc.Product.Resource != nil {
			result.Product.ListResources = uc.Product.Resource.ListResources.Execute
			result.Product.ReadResource = uc.Product.Resource.ReadResource.Execute
			result.Product.CreateResource = uc.Product.Resource.CreateResource.Execute
			result.Product.UpdateResource = uc.Product.Resource.UpdateResource.Execute
			result.Product.DeleteResource = uc.Product.Resource.DeleteResource.Execute
		}
		// Wave C P1.C.11 — Product (service-kind) dashboard rewired to
		// service-driven path. Q-SDM-DASHBOARD-DOWNSTREAM (LOCKED 2026-05-20):
		// same-commit rewire from `uc.Product.Dashboard.Execute` (entity-
		// layer, RETIRED) to `uc.Service.Dashboard.Product.GetProductDashboard.
		// Execute` (service-layer). The view-layer
		// `productdashboardview.Response` shape is the contract owned by
		// centymo; this closure translates from proto Response to
		// view-layer Response.
		if uc.Service != nil && uc.Service.Dashboard != nil && uc.Service.Dashboard.Product != nil && uc.Service.Dashboard.Product.GetProductDashboard != nil {
			productDash := uc.Service.Dashboard.Product.GetProductDashboard
			result.Product.GetServiceDashboard = func(ctx context.Context, req *productdashboardview.Request) (*productdashboardview.Response, error) {
				wsID := consumer.GetWorkspaceIDFromContext(ctx)
				resp, err := productDash.Execute(ctx, &productdashpb.GetProductDashboardRequest{
					WorkspaceId: wsID,
					Kind:        "service",
				})
				if err != nil {
					return nil, err
				}
				if resp == nil {
					return nil, nil
				}
				out := &productdashboardview.Response{
					Stats: productdashboardview.Stats{
						TotalActive:      resp.GetStats().GetTotalActive(),
						TopRevenueName:   resp.GetStats().GetTopRevenueName(),
						TopRevenueValue:  resp.GetStats().GetTopRevenueValue(),
						LineCount:        resp.GetStats().GetLineCount(),
						RecentlyAddedCnt: resp.GetStats().GetRecentlyAddedCnt(),
					},
					LineLabels: resp.GetLineLabels(),
					LineValues: resp.GetLineValues(),
					Recent:     resp.GetRecent(),
				}
				for _, r := range resp.GetTopRevenue() {
					out.TopRevenue = append(out.TopRevenue, productdashboardview.TopRevenueRow{
						ProductID:   r.GetProductId(),
						ProductName: r.GetProductName(),
						Total:       r.GetTotal(),
					})
				}
				return out, nil
			}
		}
	}

	// -- Subscription ------------------------------------------------------------
	if uc.Subscription != nil {
		if uc.Subscription.Subscription != nil {
			result.Subscription.GetSubscriptionListPageData = uc.Subscription.Subscription.GetSubscriptionListPageData.Execute
			result.Subscription.GetSubscriptionItemPageData = uc.Subscription.Subscription.GetSubscriptionItemPageData.Execute
			result.Subscription.CreateSubscription = uc.Subscription.Subscription.CreateSubscription.Execute
			result.Subscription.ReadSubscription = uc.Subscription.Subscription.ReadSubscription.Execute
			result.Subscription.UpdateSubscription = uc.Subscription.Subscription.UpdateSubscription.Execute
			result.Subscription.DeleteSubscription = uc.Subscription.Subscription.DeleteSubscription.Execute
			result.Subscription.ListSubscriptions = uc.Subscription.Subscription.ListSubscriptions.Execute
		}
		// Phase 3 F7 closure — BillingEvent now lives as a Layer-7 use case
		// sub-aggregate at uc.Subscription.BillingEvent.{ListBySubscription,SetStatus}.
		if uc.Subscription.BillingEvent != nil {
			if uc.Subscription.BillingEvent.ListBySubscription != nil {
				result.Subscription.ListBillingEventsBySubscription = uc.Subscription.BillingEvent.ListBySubscription.Execute
			}
			if uc.Subscription.BillingEvent.SetStatus != nil {
				result.Subscription.SetBillingEventStatus = uc.Subscription.BillingEvent.SetStatus.Execute
			}
		}
		// Phase 3 F6 closure — Materialize use cases now nest under
		// uc.Subscription.Subscription.{MaterializeJobs,MaterializeInstanceJobs}.
		if uc.Subscription.Subscription != nil {
			if uc.Subscription.Subscription.MaterializeJobs != nil {
				result.Subscription.MaterializeJobsForSubscription = uc.Subscription.Subscription.MaterializeJobs.Execute
			}
			if uc.Subscription.Subscription.MaterializeInstanceJobs != nil {
				result.Subscription.MaterializeInstanceJobsForSubscription = uc.Subscription.Subscription.MaterializeInstanceJobs.Execute
			}
		}

		// -- Plan (top-level on centymo UseCases, not nested under Subscription) --
		if uc.Subscription.Plan != nil {
			result.Plan.ListPlans = uc.Subscription.Plan.ListPlans.Execute
			result.Plan.ReadPlan = uc.Subscription.Plan.ReadPlan.Execute
			result.Plan.CreatePlan = uc.Subscription.Plan.CreatePlan.Execute
			result.Plan.UpdatePlan = uc.Subscription.Plan.UpdatePlan.Execute
			result.Plan.DeletePlan = uc.Subscription.Plan.DeletePlan.Execute
			result.Plan.SearchPlansByName = uc.Subscription.Plan.SearchPlansByName.Execute
			result.Plan.CustomizePlanForClient = uc.Subscription.Plan.CustomizePlanForClient.Execute
		}

		// -- PricePlan (top-level on centymo UseCases) --
		if uc.Subscription.PricePlan != nil {
			result.PricePlan.ListPricePlans = uc.Subscription.PricePlan.ListPricePlans.Execute
			result.PricePlan.ReadPricePlan = uc.Subscription.PricePlan.ReadPricePlan.Execute
			result.PricePlan.CreatePricePlan = uc.Subscription.PricePlan.CreatePricePlan.Execute
			result.PricePlan.UpdatePricePlan = uc.Subscription.PricePlan.UpdatePricePlan.Execute
			result.PricePlan.DeletePricePlan = uc.Subscription.PricePlan.DeletePricePlan.Execute
		}
		if uc.Subscription.ProductPricePlan != nil {
			result.PricePlan.ListProductPricePlans = uc.Subscription.ProductPricePlan.ListProductPricePlans.Execute
			result.PricePlan.CreateProductPricePlan = uc.Subscription.ProductPricePlan.CreateProductPricePlan.Execute
			result.PricePlan.UpdateProductPricePlan = uc.Subscription.ProductPricePlan.UpdateProductPricePlan.Execute
			result.PricePlan.DeleteProductPricePlan = uc.Subscription.ProductPricePlan.DeleteProductPricePlan.Execute
		}

		// -- PriceSchedule (top-level on centymo UseCases) --
		if uc.Subscription.PriceSchedule != nil {
			result.PriceSchedule.ListPriceSchedules = uc.Subscription.PriceSchedule.ListPriceSchedules.Execute
			result.PriceSchedule.ReadPriceSchedule = uc.Subscription.PriceSchedule.ReadPriceSchedule.Execute
			result.PriceSchedule.CreatePriceSchedule = uc.Subscription.PriceSchedule.CreatePriceSchedule.Execute
			result.PriceSchedule.UpdatePriceSchedule = uc.Subscription.PriceSchedule.UpdatePriceSchedule.Execute
			result.PriceSchedule.DeletePriceSchedule = uc.Subscription.PriceSchedule.DeletePriceSchedule.Execute
		}
		if uc.Subscription.Subscription != nil {
			result.PriceSchedule.ListSubscriptionsByPricePlan = uc.Subscription.Subscription.ListSubscriptionsByPricePlan.Execute
		}

		// -- SubscriptionGroup (education "section / cohort" cohort) --
		if uc.Subscription.SubscriptionGroup != nil {
			result.SubscriptionGroup.ListSubscriptionGroups = uc.Subscription.SubscriptionGroup.ListSubscriptionGroups.Execute
			result.SubscriptionGroup.ReadSubscriptionGroup = uc.Subscription.SubscriptionGroup.ReadSubscriptionGroup.Execute
			result.SubscriptionGroup.CreateSubscriptionGroup = uc.Subscription.SubscriptionGroup.CreateSubscriptionGroup.Execute
			result.SubscriptionGroup.UpdateSubscriptionGroup = uc.Subscription.SubscriptionGroup.UpdateSubscriptionGroup.Execute
			result.SubscriptionGroup.DeleteSubscriptionGroup = uc.Subscription.SubscriptionGroup.DeleteSubscriptionGroup.Execute
		}
	}

	// -- Revenue -----------------------------------------------------------------
	if uc.Revenue != nil {
		if uc.Revenue.Revenue != nil {
			result.Revenue.GetListPageData = uc.Revenue.Revenue.GetRevenueListPageData.Execute
			result.Revenue.CreateRevenue = uc.Revenue.Revenue.CreateRevenue.Execute
			result.Revenue.ReadRevenue = uc.Revenue.Revenue.ReadRevenue.Execute
			result.Revenue.UpdateRevenue = uc.Revenue.Revenue.UpdateRevenue.Execute
			result.Revenue.DeleteRevenue = uc.Revenue.Revenue.DeleteRevenue.Execute
			result.Revenue.RecognizeRevenueFromSubscription = uc.Revenue.Revenue.RecognizeRevenueFromSubscription.Execute
			result.Revenue.ListRevenueRunCandidates = uc.Revenue.Revenue.ListRevenueRunCandidates.Execute
			result.Revenue.GenerateRevenueRun = uc.Revenue.Revenue.GenerateRevenueRun.Execute
		}
		if uc.Revenue.RevenueLineItem != nil {
			result.Revenue.CreateRevenueLineItem = uc.Revenue.RevenueLineItem.CreateRevenueLineItem.Execute
			result.Revenue.ReadRevenueLineItem = uc.Revenue.RevenueLineItem.ReadRevenueLineItem.Execute
			result.Revenue.UpdateRevenueLineItem = uc.Revenue.RevenueLineItem.UpdateRevenueLineItem.Execute
			result.Revenue.DeleteRevenueLineItem = uc.Revenue.RevenueLineItem.DeleteRevenueLineItem.Execute
			result.Revenue.ListRevenueLineItems = uc.Revenue.RevenueLineItem.ListRevenueLineItems.Execute
		}
		if uc.Revenue.RevenueTaxLine != nil {
			result.Revenue.ListRevenueTaxLines = uc.Revenue.RevenueTaxLine.ListRevenueTaxLines.Execute
		}

		// -- RevenuePayment CRUD (W5 typed path) ---------------------------------
		// Typed revenue_payment CRUD replacing the deleted DataSource duck's
		// ListSimple/Create/Read/Update/Delete on the "revenue_payment"
		// collection. Bound from the espyna consumer aggregate; nil-safe.
		// 20260612-datasource-typed-path W7.
		if uc.Revenue.RevenuePayment != nil {
			result.Revenue.RevenuePayment.CreateRevenuePayment = uc.Revenue.RevenuePayment.CreateRevenuePayment.Execute
			result.Revenue.RevenuePayment.ReadRevenuePayment = uc.Revenue.RevenuePayment.ReadRevenuePayment.Execute
			result.Revenue.RevenuePayment.UpdateRevenuePayment = uc.Revenue.RevenuePayment.UpdateRevenuePayment.Execute
			result.Revenue.RevenuePayment.DeleteRevenuePayment = uc.Revenue.RevenuePayment.DeleteRevenuePayment.Execute
			result.Revenue.RevenuePayment.ListRevenuePayments = uc.Revenue.RevenuePayment.ListRevenuePayments.Execute
		}

		// -- RevenueRun: repo-direct pass-through via GenerateRevenueRun.RevenueRunRepo() --
		// When the revenue_run postgres adapter isn't registered (e.g. mock_db, or
		// when the factory is missing), repo is nil. Install nil-safe stubs that
		// return empty responses — matching the pre-Phase-0 consumer behavior — so
		// centymo's RequireFor passes and the run-history UI degrades to empty state.
		if uc.Revenue.Revenue != nil && uc.Revenue.Revenue.GenerateRevenueRun != nil {
			repo := uc.Revenue.Revenue.GenerateRevenueRun.RevenueRunRepo()
			if repo != nil {
				result.RevenueRun.ListRevenueRuns = repo.ListRevenueRuns
				result.RevenueRun.ReadRevenueRun = repo.ReadRevenueRun
				result.RevenueRun.ListRevenueRunAttempts = repo.ListRevenueRunAttempts
			} else {
				result.RevenueRun.ListRevenueRuns = func(context.Context, *revenuerunpb.ListRevenueRunsRequest) (*revenuerunpb.ListRevenueRunsResponse, error) {
					return &revenuerunpb.ListRevenueRunsResponse{Success: true}, nil
				}
				result.RevenueRun.ReadRevenueRun = func(context.Context, *revenuerunpb.ReadRevenueRunRequest) (*revenuerunpb.ReadRevenueRunResponse, error) {
					return &revenuerunpb.ReadRevenueRunResponse{}, nil
				}
				result.RevenueRun.ListRevenueRunAttempts = func(context.Context, *revenuerunpb.ListRevenueRunAttemptsRequest) (*revenuerunpb.ListRevenueRunAttemptsResponse, error) {
					return &revenuerunpb.ListRevenueRunAttemptsResponse{}, nil
				}
			}
		}
	}

	// -- Collection (Treasury) ---------------------------------------------------
	if uc.Treasury != nil && uc.Treasury.Collection != nil {
		result.Collection.ListCollections = uc.Treasury.Collection.ListCollections.Execute
		result.Collection.ReadCollection = uc.Treasury.Collection.ReadCollection.Execute
		result.Collection.CreateCollection = uc.Treasury.Collection.CreateCollection.Execute
		result.Collection.UpdateCollection = uc.Treasury.Collection.UpdateCollection.Execute
		result.Collection.DeleteCollection = uc.Treasury.Collection.DeleteCollection.Execute
	}

	// -- CollectionMethod (Treasury) — typed reads for the revenue/expenditure
	// collection-method dropdowns. Replaces the deleted DataSource duck's
	// ListSimple/Read on the "collection_method" collection. Nil-safe.
	// 20260612-datasource-typed-path W7.
	if uc.Treasury != nil && uc.Treasury.CollectionMethod != nil {
		result.CollectionMethod.ReadCollectionMethod = uc.Treasury.CollectionMethod.ReadCollectionMethod.Execute
		result.CollectionMethod.ListCollectionMethods = uc.Treasury.CollectionMethod.ListCollectionMethods.Execute
	}

	// Wave B P1.C.5 — Cash dashboard rewired to service-driven path.
	// Q-SDM-DASHBOARD-DOWNSTREAM + Q-SDM-DASHBOARD-COUNT: Cash is the second
	// slice of the unified Treasury candidate; the closure bridges proto
	// `treasurydashpb.GetCashDashboardResponse` → centymo view-layer
	// `*cashdashboardview.Response`. Codex review P1 (codex-review-phase1-
	// round1b-ledger-treasury.md REVISE-MAJOR) noted this closure was missing
	// at Wave B Round 1b landing; same-commit follow-up fix 2026-05-21.
	if uc.Service != nil && uc.Service.Dashboard != nil && uc.Service.Dashboard.Treasury != nil && uc.Service.Dashboard.Treasury.Cash != nil && uc.Service.Dashboard.Treasury.Cash.GetCashDashboard != nil {
		cashDash := uc.Service.Dashboard.Treasury.Cash.GetCashDashboard
		result.Collection.GetCashDashboard = func(ctx context.Context, req *cashdashboardview.Request) (*cashdashboardview.Response, error) {
			wsID := consumer.GetWorkspaceIDFromContext(ctx)
			var nowMillis int64
			if req != nil && !req.Now.IsZero() {
				nowMillis = req.Now.UnixMilli()
			}
			resp, err := cashDash.Execute(ctx, &treasurydashpb.GetCashDashboardRequest{
				WorkspaceId: wsID,
				NowMillis:   &nowMillis,
			})
			if err != nil {
				return nil, err
			}
			if resp == nil {
				return nil, nil
			}
			out := &cashdashboardview.Response{
				Stats: cashdashboardview.Stats{
					Pending:           resp.GetStats().GetPending(),
					Overdue:           resp.GetStats().GetOverdue(),
					CollectedToday:    resp.GetStats().GetCollectedToday(),
					CollectedThisWeek: resp.GetStats().GetCollectedThisWeek(),
				},
				DailyLabels: resp.GetDailyLabels(),
				DailyValues: resp.GetDailyValues(),
				ModeLabels:  resp.GetModeLabels(),
				ModeValues:  resp.GetModeValues(),
				Recent:      resp.GetRecent(),
			}
			return out, nil
		}
	}

	// 20260518-hexagonal-strict-adherence Phase 1.D — settle/refund/cancel
	// workflows now route directly through the entity-nested use cases on the
	// treasury aggregator. The TreasuryAdvancesAdapter facade in
	// consumer/adapter_treasury_advances.go has been deleted.
	if uc.Treasury != nil && uc.Treasury.Collection != nil {
		if uc.Treasury.Collection.SettleUnscheduledAdvance != nil {
			result.Collection.SettleUnscheduledAdvance = func(ctx context.Context, in AdvanceSettleInput) (*AdvanceSettleOutput, error) {
				out, err := uc.Treasury.Collection.SettleUnscheduledAdvance.Execute(ctx, &treasurycollectionpb.SettleUnscheduledAdvanceCollectionRequest{
					TreasuryCollectionId: in.AdvanceID,
					Amount:               in.Amount,
					TargetAccountId:      in.TargetAccountID,
					Reason:               in.Reason,
				})
				if err != nil {
					return nil, err
				}
				return &AdvanceSettleOutput{
					NewRemainingAmount:  out.GetNewRemainingAmount(),
					NewRecognizedAmount: out.GetNewRecognizedAmount(),
					NewStatus:           advanceStatusTail(out.GetNewStatus()),
				}, nil
			}
		}
		if uc.Treasury.Collection.RefundUnscheduledAdvance != nil {
			result.Collection.RefundUnscheduledAdvance = func(ctx context.Context, in AdvanceRefundInput) (*AdvanceRefundOutput, error) {
				out, err := uc.Treasury.Collection.RefundUnscheduledAdvance.Execute(ctx, &treasurycollectionpb.RefundUnscheduledAdvanceCollectionRequest{
					TreasuryCollectionId: in.AdvanceID,
					Amount:               in.Amount,
					RefundMethod:         in.RefundMethod,
					Destination:          in.DestinationAccount,
					Reason:               in.Reason,
				})
				if err != nil {
					return nil, err
				}
				return &AdvanceRefundOutput{
					NewRemainingAmount: out.GetNewRemainingAmount(),
					NewStatus:          advanceStatusTail(out.GetNewStatus()),
				}, nil
			}
		}
		if uc.Treasury.Collection.CancelAdvance != nil {
			result.Collection.CancelAdvance = func(ctx context.Context, in AdvanceCancelInput) (*AdvanceCancelOutput, error) {
				out, err := uc.Treasury.Collection.CancelAdvance.Execute(ctx, &treasurycollectionpb.CancelAdvanceCollectionRequest{
					TreasuryCollectionId: in.AdvanceID,
					Reason:               in.Reason,
				})
				if err != nil {
					return nil, err
				}
				return &AdvanceCancelOutput{NewStatus: advanceStatusTail(out.GetNewStatus())}, nil
			}
		}
	}

	// -- Disbursement (Treasury) -------------------------------------------------
	if uc.Treasury != nil && uc.Treasury.Disbursement != nil {
		result.Disbursement.ListDisbursements = uc.Treasury.Disbursement.ListDisbursements.Execute
		result.Disbursement.ReadDisbursement = uc.Treasury.Disbursement.ReadDisbursement.Execute
		result.Disbursement.CreateDisbursement = uc.Treasury.Disbursement.CreateDisbursement.Execute
		result.Disbursement.UpdateDisbursement = uc.Treasury.Disbursement.UpdateDisbursement.Execute
		result.Disbursement.DeleteDisbursement = uc.Treasury.Disbursement.DeleteDisbursement.Execute
	}

	// 20260518-hexagonal-strict-adherence Phase 1.D — buying-side mirrors.
	if uc.Treasury != nil && uc.Treasury.Disbursement != nil {
		if uc.Treasury.Disbursement.SettleUnscheduledAdvance != nil {
			result.Disbursement.SettleUnscheduledAdvance = func(ctx context.Context, in AdvanceSettleInput) (*AdvanceSettleOutput, error) {
				out, err := uc.Treasury.Disbursement.SettleUnscheduledAdvance.Execute(ctx, &treasurydisbursementpb.SettleUnscheduledAdvanceDisbursementRequest{
					TreasuryDisbursementId: in.AdvanceID,
					Amount:                 in.Amount,
					TargetAccountId:        in.TargetAccountID,
					Reason:                 in.Reason,
				})
				if err != nil {
					return nil, err
				}
				return &AdvanceSettleOutput{
					NewRemainingAmount:  out.GetNewRemainingAmount(),
					NewRecognizedAmount: out.GetNewRecognizedAmount(),
					NewStatus:           advanceStatusTail(out.GetNewStatus()),
				}, nil
			}
		}
		if uc.Treasury.Disbursement.RefundUnscheduledAdvance != nil {
			result.Disbursement.RefundUnscheduledAdvance = func(ctx context.Context, in AdvanceRefundInput) (*AdvanceRefundOutput, error) {
				out, err := uc.Treasury.Disbursement.RefundUnscheduledAdvance.Execute(ctx, &treasurydisbursementpb.RefundUnscheduledAdvanceDisbursementRequest{
					TreasuryDisbursementId: in.AdvanceID,
					Amount:                 in.Amount,
					RefundMethod:           in.RefundMethod,
					Destination:            in.DestinationAccount,
					Reason:                 in.Reason,
				})
				if err != nil {
					return nil, err
				}
				return &AdvanceRefundOutput{
					NewRemainingAmount: out.GetNewRemainingAmount(),
					NewStatus:          advanceStatusTail(out.GetNewStatus()),
				}, nil
			}
		}
		if uc.Treasury.Disbursement.CancelAdvance != nil {
			result.Disbursement.CancelAdvance = func(ctx context.Context, in AdvanceCancelInput) (*AdvanceCancelOutput, error) {
				out, err := uc.Treasury.Disbursement.CancelAdvance.Execute(ctx, &treasurydisbursementpb.CancelAdvanceDisbursementRequest{
					TreasuryDisbursementId: in.AdvanceID,
					Reason:                 in.Reason,
				})
				if err != nil {
					return nil, err
				}
				return &AdvanceCancelOutput{NewStatus: advanceStatusTail(out.GetNewStatus())}, nil
			}
		}
	}

	// 20260518-hexagonal-strict-adherence Phase 1.D — workspace Advances
	// Dashboard. The single cross-entity GetAdvancesDashboard use case is
	// replaced by two entity-side ListAdvancesForDashboard use cases (F5
	// resolution). service-admin stacks the per-side rows + totals into the
	// view-typed AdvancesDashboardData.
	if uc.Treasury != nil &&
		uc.Treasury.Collection != nil && uc.Treasury.Collection.ListAdvancesForDashboard != nil &&
		uc.Treasury.Disbursement != nil && uc.Treasury.Disbursement.ListAdvancesForDashboard != nil {
		result.TreasuryAdvances.GetAdvancesDashboard = func(ctx context.Context, asOfDate string) (*AdvancesDashboardData, error) {
			coll, err := uc.Treasury.Collection.ListAdvancesForDashboard.Execute(ctx, &treasurycollectionpb.ListAdvanceCollectionsForDashboardRequest{
				AsOfDate: asOfDate,
			})
			if err != nil {
				return nil, err
			}
			disb, err := uc.Treasury.Disbursement.ListAdvancesForDashboard.Execute(ctx, &treasurydisbursementpb.ListAdvanceDisbursementsForDashboardRequest{
				AsOfDate: asOfDate,
			})
			if err != nil {
				return nil, err
			}
			out := &AdvancesDashboardData{
				Inflows:                convertCollectionAdvancesRows(coll.GetRows()),
				Outflows:               convertDisbursementAdvancesRows(disb.GetRows()),
				InflowTotalRemaining:   coll.GetTotalRemaining(),
				OutflowTotalRemaining:  disb.GetTotalRemaining(),
				InflowActiveCount:      int(coll.GetActiveCount()),
				OutflowActiveCount:     int(disb.GetActiveCount()),
				InflowFullyRecognized:  int(coll.GetFullyRecognizedCount()),
				OutflowFullyRecognized: int(disb.GetFullyRecognizedCount()),
			}
			return out, nil
		}
	}

	// -- Expenditure -------------------------------------------------------------
	if uc.Expenditure != nil {
		if uc.Expenditure.Expenditure != nil {
			result.Expenditure.ListExpenditures = uc.Expenditure.Expenditure.ListExpenditures.Execute
			result.Expenditure.CreateExpenditure = uc.Expenditure.Expenditure.CreateExpenditure.Execute
			result.Expenditure.ReadExpenditure = uc.Expenditure.Expenditure.ReadExpenditure.Execute
			result.Expenditure.UpdateExpenditure = uc.Expenditure.Expenditure.UpdateExpenditure.Execute
			result.Expenditure.DeleteExpenditure = uc.Expenditure.Expenditure.DeleteExpenditure.Execute
		}
		if uc.Expenditure.ExpenditureCategory != nil {
			result.Expenditure.ListExpenditureCategories = uc.Expenditure.ExpenditureCategory.ListExpenditureCategories.Execute
			result.Expenditure.CreateExpenditureCategory = uc.Expenditure.ExpenditureCategory.CreateExpenditureCategory.Execute
			result.Expenditure.ReadExpenditureCategory = uc.Expenditure.ExpenditureCategory.ReadExpenditureCategory.Execute
			result.Expenditure.UpdateExpenditureCategory = uc.Expenditure.ExpenditureCategory.UpdateExpenditureCategory.Execute
			result.Expenditure.DeleteExpenditureCategory = uc.Expenditure.ExpenditureCategory.DeleteExpenditureCategory.Execute
		}
		if uc.Expenditure.ExpenditureLineItem != nil {
			result.Expenditure.CreateExpenditureLineItem = uc.Expenditure.ExpenditureLineItem.CreateExpenditureLineItem.Execute
			result.Expenditure.ReadExpenditureLineItem = uc.Expenditure.ExpenditureLineItem.ReadExpenditureLineItem.Execute
			result.Expenditure.UpdateExpenditureLineItem = uc.Expenditure.ExpenditureLineItem.UpdateExpenditureLineItem.Execute
			result.Expenditure.DeleteExpenditureLineItem = uc.Expenditure.ExpenditureLineItem.DeleteExpenditureLineItem.Execute
			result.Expenditure.ListExpenditureLineItems = uc.Expenditure.ExpenditureLineItem.ListExpenditureLineItems.Execute
		}
		if uc.Expenditure.PurchaseOrder != nil {
			result.Expenditure.ListPurchaseOrders = uc.Expenditure.PurchaseOrder.ListPurchaseOrders.Execute
		}
		if uc.Expenditure.ExpenseRecognition != nil {
			result.Expenditure.ListExpenseRecognitions = uc.Expenditure.ExpenseRecognition.ListExpenseRecognitions.Execute
			result.Expenditure.ReadExpenseRecognition = uc.Expenditure.ExpenseRecognition.ReadExpenseRecognition.Execute
			result.Expenditure.DeleteExpenseRecognition = uc.Expenditure.ExpenseRecognition.DeleteExpenseRecognition.Execute
			result.Expenditure.ReverseExpenseRecognition = uc.Expenditure.ExpenseRecognition.ReverseExpenseRecognition.Execute
			result.Expenditure.RecognizeFromExpenditure = uc.Expenditure.ExpenseRecognition.RecognizeFromExpenditure.Execute
			result.Expenditure.RecognizeFromContract = uc.Expenditure.ExpenseRecognition.RecognizeFromContract.Execute
		}
		if uc.Expenditure.ExpenseRecognitionLine != nil {
			result.Expenditure.ListExpenseRecognitionLines = uc.Expenditure.ExpenseRecognitionLine.ListExpenseRecognitionLines.Execute
			result.Expenditure.ReadExpenseRecognitionLine = uc.Expenditure.ExpenseRecognitionLine.ReadExpenseRecognitionLine.Execute
			result.Expenditure.CreateExpenseRecognitionLine = uc.Expenditure.ExpenseRecognitionLine.CreateExpenseRecognitionLine.Execute
			result.Expenditure.UpdateExpenseRecognitionLine = uc.Expenditure.ExpenseRecognitionLine.UpdateExpenseRecognitionLine.Execute
			result.Expenditure.DeleteExpenseRecognitionLine = uc.Expenditure.ExpenseRecognitionLine.DeleteExpenseRecognitionLine.Execute
		}
		if uc.Expenditure.AccruedExpense != nil {
			result.Expenditure.ListAccruedExpenses = uc.Expenditure.AccruedExpense.ListAccruedExpenses.Execute
			result.Expenditure.ReadAccruedExpense = uc.Expenditure.AccruedExpense.ReadAccruedExpense.Execute
			result.Expenditure.CreateAccruedExpense = uc.Expenditure.AccruedExpense.CreateAccruedExpense.Execute
			result.Expenditure.UpdateAccruedExpense = uc.Expenditure.AccruedExpense.UpdateAccruedExpense.Execute
			result.Expenditure.DeleteAccruedExpense = uc.Expenditure.AccruedExpense.DeleteAccruedExpense.Execute
			result.Expenditure.AccrueFromContract = uc.Expenditure.AccruedExpense.AccrueFromContract.Execute
			result.Expenditure.ReverseAccrual = uc.Expenditure.AccruedExpense.ReverseAccrual.Execute
			result.Expenditure.SettleAccrual = uc.Expenditure.AccruedExpense.SettleAccrual.Execute
		}
		if uc.Expenditure.AccruedExpenseSettlement != nil {
			result.Expenditure.ListAccruedExpenseSettlements = uc.Expenditure.AccruedExpenseSettlement.ListAccruedExpenseSettlements.Execute
			result.Expenditure.CreateAccruedExpenseSettlement = uc.Expenditure.AccruedExpenseSettlement.CreateAccruedExpenseSettlement.Execute
			result.Expenditure.ReadAccruedExpenseSettlement = uc.Expenditure.AccruedExpenseSettlement.ReadAccruedExpenseSettlement.Execute
			result.Expenditure.UpdateAccruedExpenseSettlement = uc.Expenditure.AccruedExpenseSettlement.UpdateAccruedExpenseSettlement.Execute
			result.Expenditure.DeleteAccruedExpenseSettlement = uc.Expenditure.AccruedExpenseSettlement.DeleteAccruedExpenseSettlement.Execute
		}
		// Wave C P1.C.8 — Expenditure dashboard (purchase + expense surfaces)
		// rewired to service-driven path. Q-SDM-DASHBOARD-DOWNSTREAM (LOCKED
		// 2026-05-20): same-commit rewire from `uc.Expenditure.Dashboard.
		// Execute` (entity-layer, RETIRED) to `uc.Service.Dashboard.
		// Expenditure.GetExpenditureDashboard.Execute` (service-layer). One
		// use case serves both surfaces — the Kind discriminator selects
		// "purchase" vs "expense".
		if uc.Service != nil && uc.Service.Dashboard != nil && uc.Service.Dashboard.Expenditure != nil && uc.Service.Dashboard.Expenditure.GetExpenditureDashboard != nil {
			expDash := uc.Service.Dashboard.Expenditure.GetExpenditureDashboard

			result.Expenditure.GetPurchaseDashboard = func(ctx context.Context, req *purchaseboard.Request) (*purchaseboard.Response, error) {
				wsID := consumer.GetWorkspaceIDFromContext(ctx)
				resp, err := expDash.Execute(ctx, &expendituredashpb.GetExpenditureDashboardRequest{
					WorkspaceId: wsID,
					Kind:        "purchase",
				})
				if err != nil {
					return nil, err
				}
				if resp == nil {
					return nil, nil
				}
				out := &purchaseboard.Response{
					Stats: purchaseboard.Stats{
						OpenCount:        resp.GetStats().GetOpenCount(),
						AwaitingCount:    resp.GetStats().GetAwaitingCount(),
						SpentMTD:         resp.GetStats().GetTotalMtd(),
						TopSupplierName:  resp.GetStats().GetTopSupplierName(),
						TopSupplierTotal: resp.GetStats().GetTopSupplierTotal(),
					},
					MonthLabels: resp.GetMonthLabels(),
					MonthValues: resp.GetMonthValues(),
					Recent:      resp.GetRecent(),
				}
				for _, r := range resp.GetTopSuppliers() {
					out.TopSuppliers = append(out.TopSuppliers, purchaseboard.TopSupplierRow{
						SupplierID:   r.GetSupplierId(),
						SupplierName: r.GetSupplierName(),
						Total:        r.GetTotal(),
					})
				}
				return out, nil
			}

			result.Expenditure.GetExpenseDashboard = func(ctx context.Context, req *expenseboard.Request) (*expenseboard.Response, error) {
				wsID := consumer.GetWorkspaceIDFromContext(ctx)
				resp, err := expDash.Execute(ctx, &expendituredashpb.GetExpenditureDashboardRequest{
					WorkspaceId: wsID,
					Kind:        "expense",
				})
				if err != nil {
					return nil, err
				}
				if resp == nil {
					return nil, nil
				}
				return &expenseboard.Response{
					Stats: expenseboard.Stats{
						PendingApprovalCount: resp.GetStats().GetOpenCount(),
						ApprovedMTD:          resp.GetStats().GetTotalMtd(),
						ReimbursableMTD:      resp.GetStats().GetReimbursableMtd(),
						CategoriesUsed:       resp.GetStats().GetCategoryCount(),
					},
					CategoryLabels: resp.GetCategoryLabels(),
					CategoryValues: resp.GetCategoryValues(),
					Recent:         resp.GetRecent(),
				}, nil
			}
		}

		// -- SupplierContract --------------------------------------------------------
		if uc.Expenditure.SupplierContract != nil {
			result.SupplierContract.ListSupplierContracts = uc.Expenditure.SupplierContract.ListSupplierContracts.Execute
			result.SupplierContract.ReadSupplierContract = uc.Expenditure.SupplierContract.ReadSupplierContract.Execute
			result.SupplierContract.CreateSupplierContract = uc.Expenditure.SupplierContract.CreateSupplierContract.Execute
			result.SupplierContract.UpdateSupplierContract = uc.Expenditure.SupplierContract.UpdateSupplierContract.Execute
			result.SupplierContract.DeleteSupplierContract = uc.Expenditure.SupplierContract.DeleteSupplierContract.Execute
			result.SupplierContract.ApproveSupplierContract = uc.Expenditure.SupplierContract.ApproveSupplierContract.Execute
			result.SupplierContract.TerminateSupplierContract = uc.Expenditure.SupplierContract.TerminateSupplierContract.Execute
		}
		if uc.Expenditure.SupplierContractLine != nil {
			result.SupplierContract.ListSupplierContractLines = uc.Expenditure.SupplierContractLine.ListSupplierContractLines.Execute
			result.SupplierContract.ReadSupplierContractLine = uc.Expenditure.SupplierContractLine.ReadSupplierContractLine.Execute
			result.SupplierContract.CreateSupplierContractLine = uc.Expenditure.SupplierContractLine.CreateSupplierContractLine.Execute
			result.SupplierContract.UpdateSupplierContractLine = uc.Expenditure.SupplierContractLine.UpdateSupplierContractLine.Execute
			result.SupplierContract.DeleteSupplierContractLine = uc.Expenditure.SupplierContractLine.DeleteSupplierContractLine.Execute
		}
		if uc.Expenditure.SupplierContractPriceSchedule != nil {
			result.SupplierContract.ListSupplierContractPriceSchedules = uc.Expenditure.SupplierContractPriceSchedule.ListSupplierContractPriceSchedules.Execute
			result.SupplierContract.CreateSupplierContractPriceSchedule = uc.Expenditure.SupplierContractPriceSchedule.CreateSupplierContractPriceSchedule.Execute
			result.SupplierContract.ReadSupplierContractPriceSchedule = uc.Expenditure.SupplierContractPriceSchedule.ReadSupplierContractPriceSchedule.Execute
			result.SupplierContract.UpdateSupplierContractPriceSchedule = uc.Expenditure.SupplierContractPriceSchedule.UpdateSupplierContractPriceSchedule.Execute
			result.SupplierContract.DeleteSupplierContractPriceSchedule = uc.Expenditure.SupplierContractPriceSchedule.DeleteSupplierContractPriceSchedule.Execute
			result.SupplierContract.ActivateSupplierContractPriceSchedule = uc.Expenditure.SupplierContractPriceSchedule.ActivateSupplierContractPriceSchedule.Execute
			result.SupplierContract.SupersedeSupplierContractPriceSchedule = uc.Expenditure.SupplierContractPriceSchedule.SupersedeSupplierContractPriceSchedule.Execute
		}
		if uc.Expenditure.SupplierContractPriceScheduleLine != nil {
			result.SupplierContract.ListSupplierContractPriceScheduleLines = uc.Expenditure.SupplierContractPriceScheduleLine.ListSupplierContractPriceScheduleLines.Execute
			result.SupplierContract.ReadSupplierContractPriceScheduleLine = uc.Expenditure.SupplierContractPriceScheduleLine.ReadSupplierContractPriceScheduleLine.Execute
			result.SupplierContract.CreateSupplierContractPriceScheduleLine = uc.Expenditure.SupplierContractPriceScheduleLine.CreateSupplierContractPriceScheduleLine.Execute
			result.SupplierContract.UpdateSupplierContractPriceScheduleLine = uc.Expenditure.SupplierContractPriceScheduleLine.UpdateSupplierContractPriceScheduleLine.Execute
			result.SupplierContract.DeleteSupplierContractPriceScheduleLine = uc.Expenditure.SupplierContractPriceScheduleLine.DeleteSupplierContractPriceScheduleLine.Execute
		}
		if uc.Expenditure.ExpenseRecognition != nil {
			result.SupplierContract.ListExpenseRecognitions = uc.Expenditure.ExpenseRecognition.ListExpenseRecognitions.Execute
			result.SupplierContract.ReadExpenseRecognition = uc.Expenditure.ExpenseRecognition.ReadExpenseRecognition.Execute
		}
		if uc.Expenditure.ExpenseRecognitionLine != nil {
			result.SupplierContract.ListExpenseRecognitionLines = uc.Expenditure.ExpenseRecognitionLine.ListExpenseRecognitionLines.Execute
			result.SupplierContract.ReadExpenseRecognitionLine = uc.Expenditure.ExpenseRecognitionLine.ReadExpenseRecognitionLine.Execute
		}
		if uc.Expenditure.ProcurementRequest != nil {
			result.SupplierContract.ListProcurementRequests = uc.Expenditure.ProcurementRequest.ListProcurementRequests.Execute
			result.SupplierContract.ReadProcurementRequest = uc.Expenditure.ProcurementRequest.ReadProcurementRequest.Execute
			result.SupplierContract.CreateProcurementRequest = uc.Expenditure.ProcurementRequest.CreateProcurementRequest.Execute
			result.SupplierContract.UpdateProcurementRequest = uc.Expenditure.ProcurementRequest.UpdateProcurementRequest.Execute
			result.SupplierContract.DeleteProcurementRequest = uc.Expenditure.ProcurementRequest.DeleteProcurementRequest.Execute
			result.SupplierContract.SubmitProcurementRequest = uc.Expenditure.ProcurementRequest.SubmitProcurementRequest.Execute
			result.SupplierContract.ApproveProcurementRequest = uc.Expenditure.ProcurementRequest.ApproveProcurementRequest.Execute
			result.SupplierContract.RejectProcurementRequest = uc.Expenditure.ProcurementRequest.RejectProcurementRequest.Execute
			result.SupplierContract.SpawnProcurementRequestPO = uc.Expenditure.ProcurementRequest.SpawnPurchaseOrder.Execute
		}
		if uc.Expenditure.ProcurementRequestLine != nil {
			result.SupplierContract.ListProcurementRequestLines = uc.Expenditure.ProcurementRequestLine.ListProcurementRequestLines.Execute
			result.SupplierContract.ReadProcurementRequestLine = uc.Expenditure.ProcurementRequestLine.ReadProcurementRequestLine.Execute
			result.SupplierContract.CreateProcurementRequestLine = uc.Expenditure.ProcurementRequestLine.CreateProcurementRequestLine.Execute
			result.SupplierContract.UpdateProcurementRequestLine = uc.Expenditure.ProcurementRequestLine.UpdateProcurementRequestLine.Execute
			result.SupplierContract.DeleteProcurementRequestLine = uc.Expenditure.ProcurementRequestLine.DeleteProcurementRequestLine.Execute
		}

		// -- ExpenseRecognitionRun (Plan A Phase 4) ------------------------------
		// Proto-direct repo pass-through for List / Read / ListAttempts (mirror
		// of RevenueRun above). Application-layer use cases for candidates +
		// generate are routed through espyna's GenerateExpenseRun wrapper which
		// exposes the underlying run repo.
		// Phase 3 F6 closure — ListExpenseRunCandidates + GenerateExpenseRun
		// now nest under uc.Expenditure.ExpenseRecognitionRun.
		if uc.Expenditure.ExpenseRecognitionRun != nil {
			if uc.Expenditure.ExpenseRecognitionRun.GenerateExpenseRun != nil {
				repo := uc.Expenditure.ExpenseRecognitionRun.GenerateExpenseRun.ExpenseRecognitionRunRepo()
				if repo != nil {
					result.ExpenseRecognitionRun.ListExpenseRecognitionRuns = repo.ListExpenseRecognitionRuns
					result.ExpenseRecognitionRun.ReadExpenseRecognitionRun = repo.ReadExpenseRecognitionRun
					result.ExpenseRecognitionRun.ListExpenseRecognitionRunAttempts = repo.ListExpenseRecognitionRunAttempts
				} else {
					result.ExpenseRecognitionRun.ListExpenseRecognitionRuns = func(context.Context, *expenserecognitionrunpb.ListExpenseRecognitionRunsRequest) (*expenserecognitionrunpb.ListExpenseRecognitionRunsResponse, error) {
						return &expenserecognitionrunpb.ListExpenseRecognitionRunsResponse{Success: true}, nil
					}
					result.ExpenseRecognitionRun.ReadExpenseRecognitionRun = func(context.Context, *expenserecognitionrunpb.ReadExpenseRecognitionRunRequest) (*expenserecognitionrunpb.ReadExpenseRecognitionRunResponse, error) {
						return &expenserecognitionrunpb.ReadExpenseRecognitionRunResponse{}, nil
					}
					result.ExpenseRecognitionRun.ListExpenseRecognitionRunAttempts = func(context.Context, *expenserecognitionrunpb.ListExpenseRecognitionRunAttemptsRequest) (*expenserecognitionrunpb.ListExpenseRecognitionRunAttemptsResponse, error) {
						return &expenserecognitionrunpb.ListExpenseRecognitionRunAttemptsResponse{}, nil
					}
				}
				result.ExpenseRecognitionRun.GenerateExpenseRun = uc.Expenditure.ExpenseRecognitionRun.GenerateExpenseRun.Execute
			}
			if uc.Expenditure.ExpenseRecognitionRun.ListExpenseRunCandidates != nil {
				result.ExpenseRecognitionRun.ListExpenseRunCandidates = uc.Expenditure.ExpenseRecognitionRun.ListExpenseRunCandidates.Execute
			}
		}
	}

	// -- Procurement (P3 supplier subscriptions) ---------------------------------
	if uc.Procurement != nil {
		if uc.Procurement.CostSchedule != nil {
			result.Procurement.CostSchedule.ListCostSchedules = uc.Procurement.CostSchedule.ListCostSchedules.Execute
			result.Procurement.CostSchedule.ReadCostSchedule = uc.Procurement.CostSchedule.ReadCostSchedule.Execute
			result.Procurement.CostSchedule.CreateCostSchedule = uc.Procurement.CostSchedule.CreateCostSchedule.Execute
			result.Procurement.CostSchedule.UpdateCostSchedule = uc.Procurement.CostSchedule.UpdateCostSchedule.Execute
			result.Procurement.CostSchedule.DeleteCostSchedule = uc.Procurement.CostSchedule.DeleteCostSchedule.Execute
			result.Procurement.CostSchedule.GetCostScheduleListPageData = uc.Procurement.CostSchedule.GetCostScheduleListPageData.Execute
			result.Procurement.CostSchedule.GetCostScheduleItemPageData = uc.Procurement.CostSchedule.GetCostScheduleItemPageData.Execute
		}
		if uc.Procurement.SupplierPlan != nil {
			result.Procurement.SupplierPlan.ListSupplierPlans = uc.Procurement.SupplierPlan.ListSupplierPlans.Execute
			result.Procurement.SupplierPlan.ReadSupplierPlan = uc.Procurement.SupplierPlan.ReadSupplierPlan.Execute
			result.Procurement.SupplierPlan.CreateSupplierPlan = uc.Procurement.SupplierPlan.CreateSupplierPlan.Execute
			result.Procurement.SupplierPlan.UpdateSupplierPlan = uc.Procurement.SupplierPlan.UpdateSupplierPlan.Execute
			result.Procurement.SupplierPlan.DeleteSupplierPlan = uc.Procurement.SupplierPlan.DeleteSupplierPlan.Execute
			result.Procurement.SupplierPlan.GetSupplierPlanListPageData = uc.Procurement.SupplierPlan.GetSupplierPlanListPageData.Execute
			result.Procurement.SupplierPlan.GetSupplierPlanItemPageData = uc.Procurement.SupplierPlan.GetSupplierPlanItemPageData.Execute
		}
		if uc.Procurement.CostPlan != nil {
			result.Procurement.CostPlan.ListCostPlans = uc.Procurement.CostPlan.ListCostPlans.Execute
			result.Procurement.CostPlan.ReadCostPlan = uc.Procurement.CostPlan.ReadCostPlan.Execute
			result.Procurement.CostPlan.CreateCostPlan = uc.Procurement.CostPlan.CreateCostPlan.Execute
			result.Procurement.CostPlan.UpdateCostPlan = uc.Procurement.CostPlan.UpdateCostPlan.Execute
			result.Procurement.CostPlan.DeleteCostPlan = uc.Procurement.CostPlan.DeleteCostPlan.Execute
			result.Procurement.CostPlan.GetCostPlanListPageData = uc.Procurement.CostPlan.GetCostPlanListPageData.Execute
			result.Procurement.CostPlan.GetCostPlanItemPageData = uc.Procurement.CostPlan.GetCostPlanItemPageData.Execute
		}
		if uc.Procurement.SupplierProductPlan != nil {
			result.Procurement.SupplierProductPlan.ListSupplierProductPlans = uc.Procurement.SupplierProductPlan.ListSupplierProductPlans.Execute
			result.Procurement.SupplierProductPlan.ReadSupplierProductPlan = uc.Procurement.SupplierProductPlan.ReadSupplierProductPlan.Execute
			result.Procurement.SupplierProductPlan.CreateSupplierProductPlan = uc.Procurement.SupplierProductPlan.CreateSupplierProductPlan.Execute
			result.Procurement.SupplierProductPlan.UpdateSupplierProductPlan = uc.Procurement.SupplierProductPlan.UpdateSupplierProductPlan.Execute
			result.Procurement.SupplierProductPlan.DeleteSupplierProductPlan = uc.Procurement.SupplierProductPlan.DeleteSupplierProductPlan.Execute
			result.Procurement.SupplierProductPlan.GetSupplierProductPlanListPageData = uc.Procurement.SupplierProductPlan.GetSupplierProductPlanListPageData.Execute
			result.Procurement.SupplierProductPlan.GetSupplierProductPlanItemPageData = uc.Procurement.SupplierProductPlan.GetSupplierProductPlanItemPageData.Execute
		}
		if uc.Procurement.SupplierProductCostPlan != nil {
			result.Procurement.SupplierProductCostPlan.ListSupplierProductCostPlans = uc.Procurement.SupplierProductCostPlan.ListSupplierProductCostPlans.Execute
			result.Procurement.SupplierProductCostPlan.ReadSupplierProductCostPlan = uc.Procurement.SupplierProductCostPlan.ReadSupplierProductCostPlan.Execute
			result.Procurement.SupplierProductCostPlan.CreateSupplierProductCostPlan = uc.Procurement.SupplierProductCostPlan.CreateSupplierProductCostPlan.Execute
			result.Procurement.SupplierProductCostPlan.UpdateSupplierProductCostPlan = uc.Procurement.SupplierProductCostPlan.UpdateSupplierProductCostPlan.Execute
			result.Procurement.SupplierProductCostPlan.DeleteSupplierProductCostPlan = uc.Procurement.SupplierProductCostPlan.DeleteSupplierProductCostPlan.Execute
			result.Procurement.SupplierProductCostPlan.GetSupplierProductCostPlanItemPageData = uc.Procurement.SupplierProductCostPlan.GetSupplierProductCostPlanItemPageData.Execute
		}
		if uc.Procurement.SupplierSubscription != nil {
			result.Procurement.SupplierSubscription.ListSupplierSubscriptions = uc.Procurement.SupplierSubscription.ListSupplierSubscriptions.Execute
			result.Procurement.SupplierSubscription.ReadSupplierSubscription = uc.Procurement.SupplierSubscription.ReadSupplierSubscription.Execute
			result.Procurement.SupplierSubscription.CreateSupplierSubscription = uc.Procurement.SupplierSubscription.CreateSupplierSubscription.Execute
			result.Procurement.SupplierSubscription.UpdateSupplierSubscription = uc.Procurement.SupplierSubscription.UpdateSupplierSubscription.Execute
			result.Procurement.SupplierSubscription.DeleteSupplierSubscription = uc.Procurement.SupplierSubscription.DeleteSupplierSubscription.Execute
			result.Procurement.SupplierSubscription.GetSupplierSubscriptionListPageData = uc.Procurement.SupplierSubscription.GetSupplierSubscriptionListPageData.Execute
			result.Procurement.SupplierSubscription.GetSupplierSubscriptionItemPageData = uc.Procurement.SupplierSubscription.GetSupplierSubscriptionItemPageData.Execute
		}
	}

	// -- Operation ---------------------------------------------------------------
	if uc.Operation != nil {
		if uc.Operation.JobTemplate != nil {
			result.Operation.JobTemplate.ListJobTemplates = uc.Operation.JobTemplate.ListJobTemplates.Execute
			result.Operation.JobTemplate.ReadJobTemplate = uc.Operation.JobTemplate.ReadJobTemplate.Execute
		}
		if uc.Operation.JobTemplatePhase != nil {
			result.Operation.JobTemplatePhase.ListByJobTemplate = uc.Operation.JobTemplatePhase.ListByJobTemplate.Execute
		}
		if uc.Operation.JobTemplateTask != nil {
			result.Operation.JobTemplateTask.ListByPhase = uc.Operation.JobTemplateTask.ListByPhase.Execute
		}
		// Phase 3 F7 closure — JobTemplateRelation is now a Layer-7 use case
		// sub-aggregate at uc.Operation.JobTemplateRelation.ListByParent.
		if uc.Operation.JobTemplateRelation != nil && uc.Operation.JobTemplateRelation.ListByParent != nil {
			result.Operation.JobTemplateRelation.ListByParent = uc.Operation.JobTemplateRelation.ListByParent.Execute
		}
		if uc.Operation.Job != nil {
			result.Operation.Job.GetJobsByOrigin = uc.Operation.Job.GetJobsByOrigin.Execute
		}
		if uc.Operation.JobPhase != nil {
			result.Operation.JobPhase.ListByJob = uc.Operation.JobPhase.ListByJob.Execute
		}
		if uc.Operation.JobActivity != nil {
			result.Operation.JobActivity.ReadJobActivity = uc.Operation.JobActivity.ReadJobActivity.Execute
		}
	}

	return result
}

// convertCollectionAdvancesRows / convertDisbursementAdvancesRows translate
// the entity-side dashboard rows into the centymo view-typed rows. Two
// converters because the proto types are entity-distinct (selling-side rows
// come from treasury/collection, buying-side from treasury/disbursement) per
// the Phase 1.A proto split.
func convertCollectionAdvancesRows(in []*treasurycollectionpb.AdvanceCollectionDashboardRow) []AdvancesDashboardRow {
	if len(in) == 0 {
		return nil
	}
	out := make([]AdvancesDashboardRow, 0, len(in))
	for _, r := range in {
		if r == nil {
			continue
		}
		out = append(out, AdvancesDashboardRow{
			ID:               r.GetAdvanceId(),
			ReferenceNumber:  r.GetReferenceNumber(),
			CounterpartyName: r.GetCounterpartyName(),
			Kind:             advanceKindTail(r.GetKind()),
			Status:           advanceStatusTail(r.GetStatus()),
			Currency:         r.GetCurrency(),
			TotalAmount:      r.GetTotalAmount(),
			RemainingAmount:  r.GetRemainingAmount(),
			RecognizedAmount: r.GetRecognizedAmount(),
		})
	}
	return out
}

func convertDisbursementAdvancesRows(in []*treasurydisbursementpb.AdvanceDisbursementDashboardRow) []AdvancesDashboardRow {
	if len(in) == 0 {
		return nil
	}
	out := make([]AdvancesDashboardRow, 0, len(in))
	for _, r := range in {
		if r == nil {
			continue
		}
		out = append(out, AdvancesDashboardRow{
			ID:               r.GetAdvanceId(),
			ReferenceNumber:  r.GetReferenceNumber(),
			CounterpartyName: r.GetCounterpartyName(),
			Kind:             advanceKindTail(r.GetKind()),
			Status:           advanceStatusTail(r.GetStatus()),
			Currency:         r.GetCurrency(),
			TotalAmount:      r.GetTotalAmount(),
			RemainingAmount:  r.GetRemainingAmount(),
			RecognizedAmount: r.GetRecognizedAmount(),
		})
	}
	return out
}

// advanceStatusTail strips the "ADVANCE_STATUS_" prefix from the stringified
// proto enum so the view layer's status dictionary resolves it
// (e.g. ADVANCE_STATUS_PARTIALLY_SETTLED → PARTIALLY_SETTLED).
//
// 20260518-hexagonal-strict-adherence Phase 1.D — signature migrated from
// consumer.AdvanceStatus (alias from the deleted adapter_treasury_advances.go)
// to advancekindpb.AdvanceStatus.
func advanceStatusTail(s advancekindpb.AdvanceStatus) string {
	const prefix = "ADVANCE_STATUS_"
	v := s.String()
	if len(v) > len(prefix) && v[:len(prefix)] == prefix {
		return v[len(prefix):]
	}
	return v
}

// paymentTermInt32 coerces a generic ListSimple() map value (whatever numeric
// type the underlying postgres scan produced for an integer column — int64 in
// practice, but float64 / int / int32 / json.Number are tolerated) into an
// int32 for the proto PaymentTerm.NetDays field. Returns ok=false for nil /
// unrecognized types so the caller leaves the proto field at its zero value.
// 20260612-datasource-typed-path W7 (payment_term dropdown fallback).
func paymentTermInt32(v any) (int32, bool) {
	switch n := v.(type) {
	case int32:
		return n, true
	case int64:
		return int32(n), true
	case int:
		return int32(n), true
	case float64:
		return int32(n), true
	case float32:
		return int32(n), true
	default:
		return 0, false
	}
}

// advanceKindTail strips the "ADVANCE_KIND_" prefix (e.g. ADVANCE_KIND_TIME_BASED
// → TIME_BASED) for the same view-dictionary resolution reason.
func advanceKindTail(k advancekindpb.AdvanceKind) string {
	const prefix = "ADVANCE_KIND_"
	v := k.String()
	if len(v) > len(prefix) && v[:len(prefix)] == prefix {
		return v[len(prefix):]
	}
	return v
}
