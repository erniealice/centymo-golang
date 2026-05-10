// Package block — product and product-line domain wiring.
//
// Holds wireProductModules (the lifted bodies of the wantProduct and
// wantProductLine branches of Block(), including all three product mounts:
// services, inventory, and supplies).
//
// Phase 5 of the 20260510-block-go-splitting-strategy.
package block

import (
	"context"

	consumer "github.com/erniealice/espyna-golang/consumer"
	"github.com/erniealice/espyna-golang/reference"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	productmod "github.com/erniealice/centymo-golang/views/product"
	productlinemod "github.com/erniealice/centymo-golang/views/product/line"
)

// productWiring holds everything wireProductModules needs from the surrounding
// Block() scope. More than 6 fields → struct. Kept private; never re-exported.
type productWiring struct {
	db         centymo.DataSource
	refChecker reference.Checker
	// Image + attachment ops
	uploadImage      func(context.Context, string, string, []byte, string) error
	uploadFile       func(context.Context, string, string, []byte, string) error
	listAttachments  func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID  func() string
	// Routes
	productRoutes              centymo.ProductRoutes
	productInventoryRoutes     centymo.ProductRoutes
	productSuppliesRoutes      centymo.ProductRoutes
	productLineRoutes          centymo.ProductLineRoutes
	productLineInventoryRoutes centymo.ProductLineRoutes
	// Labels
	productLabels          centymo.ProductLabels
	productInventoryLabels centymo.ProductLabels
	productSuppliesLabels  centymo.ProductLabels
	productLineLabels      centymo.ProductLineLabels
	centymoTableLabels     types.TableLabels
}

// wireProductModules lifts the bodies of the two product-related `if cfg.wantXxx()`
// branches (Product with 3 mounts, ProductLine with 2 mounts) from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once at the position where
// the product wiring used to be.
func wireProductModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *consumer.UseCases, w productWiring) {
	// =====================================================================
	// Product module
	// =====================================================================

	if cfg.wantProduct() {
		var getProductInUseIDs func(context.Context, []string) (map[string]bool, error)
		if w.refChecker != nil {
			getProductInUseIDs = w.refChecker.GetProductInUseIDs
		}

		// For professional business types the product list is branded as
		// "services" and filters product_kind = 'service'.
		// Default new products created through this UI to 'service' so they
		// appear in the list immediately without extra steps.
		defaultProductKind := ""
		defaultDeliveryMode := ""
		defaultTrackingMode := ""
		if ctx.BusinessType == "professional" {
			defaultProductKind = "service"
			defaultDeliveryMode = "scheduled"
			defaultTrackingMode = "none"
		}

		productDeps := &productmod.ModuleDeps{
			Routes:              w.productRoutes,
			Mode:                "service",
			DB:                  w.db,
			Labels:              w.productLabels,
			CommonLabels:        ctx.Common,
			TableLabels:         w.centymoTableLabels,
			GetInUseIDs:         getProductInUseIDs,
			DefaultProductKind:  defaultProductKind,
			DefaultDeliveryMode: defaultDeliveryMode,
			DefaultTrackingMode: defaultTrackingMode,
			// Services mount locks product_kind to "service" (single option
			// → drawer renders the select disabled). DeliveryMode and
			// TrackingMode stay fully open so clinic admins can still pick
			// e.g. scheduled vs digital vs project per-service.
			AllowedProductKinds: []string{"service"},
			// Operation-level RBAC: every perms.Can check inside this mount
			// uses "service:*" rather than the shared "product:*". Lets a
			// role grant Services CRUD without implicit grant on Products
			// or Supplies.
			PermissionEntity: "service",
			// SetProductActive uses raw DB update (proto3 omits false booleans)
			SetProductActive: func(fctx context.Context, id string, active bool) error {
				_, err := w.db.Update(fctx, "product", id, map[string]any{"active": active})
				return err
			},
			// Image upload (product variant images)
			UploadImage: w.uploadImage,
			// Attachments
			UploadFile:       w.uploadFile,
			ListAttachments:  w.listAttachments,
			CreateAttachment: w.createAttachment,
			DeleteAttachment: w.deleteAttachment,
			NewID:            w.newAttachmentID,
		}
		if useCases.Product != nil {
			if uc := useCases.Product.Product; uc != nil {
				productDeps.ListProducts = uc.ListProducts.Execute
				productDeps.ReadProduct = uc.ReadProduct.Execute
				productDeps.CreateProduct = uc.CreateProduct.Execute
				productDeps.UpdateProduct = uc.UpdateProduct.Execute
				productDeps.DeleteProduct = uc.DeleteProduct.Execute
			}
			if uc := useCases.Product.ProductVariant; uc != nil {
				productDeps.ListProductVariants = uc.ListProductVariants.Execute
				productDeps.ReadProductVariant = uc.ReadProductVariant.Execute
				productDeps.CreateProductVariant = uc.CreateProductVariant.Execute
				productDeps.UpdateProductVariant = uc.UpdateProductVariant.Execute
				productDeps.DeleteProductVariant = uc.DeleteProductVariant.Execute
			}
			if uc := useCases.Product.ProductVariantOption; uc != nil {
				productDeps.ListProductVariantOptions = uc.ListProductVariantOptions.Execute
				productDeps.CreateProductVariantOption = uc.CreateProductVariantOption.Execute
			}
			if uc := useCases.Product.ProductOption; uc != nil {
				productDeps.ListProductOptions = uc.ListProductOptions.Execute
				productDeps.ReadProductOption = uc.ReadProductOption.Execute
				productDeps.CreateProductOption = uc.CreateProductOption.Execute
				productDeps.UpdateProductOption = uc.UpdateProductOption.Execute
				productDeps.DeleteProductOption = uc.DeleteProductOption.Execute
			}
			if uc := useCases.Product.ProductOptionValue; uc != nil {
				productDeps.ListProductOptionValues = uc.ListProductOptionValues.Execute
				productDeps.ReadProductOptionValue = uc.ReadProductOptionValue.Execute
				productDeps.CreateProductOptionValue = uc.CreateProductOptionValue.Execute
				productDeps.UpdateProductOptionValue = uc.UpdateProductOptionValue.Execute
				productDeps.DeleteProductOptionValue = uc.DeleteProductOptionValue.Execute
			}
			if uc := useCases.Product.ProductAttribute; uc != nil {
				productDeps.ListProductAttributes = uc.ListProductAttributes.Execute
				productDeps.CreateProductAttribute = uc.CreateProductAttribute.Execute
				productDeps.DeleteProductAttribute = uc.DeleteProductAttribute.Execute
			}
			if uc := useCases.Product.Line; uc != nil {
				productDeps.ListLines = uc.ListLines.Execute
			}
			if uc := useCases.Product.ProductLine; uc != nil {
				productDeps.ListProductLines = uc.ListProductLines.Execute
				productDeps.CreateProductLine = uc.CreateProductLine.Execute
				productDeps.UpdateProductLine = uc.UpdateProductLine.Execute
				productDeps.DeleteProductLine = uc.DeleteProductLine.Execute
			}
			if uc := useCases.Product.ProductVariantImage; uc != nil {
				productDeps.ListProductVariantImages = uc.ListProductVariantImages.Execute
				productDeps.CreateProductVariantImage = uc.CreateProductVariantImage.Execute
				productDeps.DeleteProductVariantImage = uc.DeleteProductVariantImage.Execute
			}
		}
		// Common Attribute (for attribute dropdowns in product detail)
		if useCases.Common != nil && useCases.Common.Attribute != nil {
			productDeps.ListAttributes = useCases.Common.Attribute.ListAttributes.Execute
			productDeps.ReadAttribute = useCases.Common.Attribute.ReadAttribute.Execute
		}
		// Inventory (for variant detail page + variant stock detail)
		if useCases.Inventory != nil {
			if uc := useCases.Inventory.InventoryItem; uc != nil {
				productDeps.ListInventoryItems = uc.ListInventoryItems.Execute
				productDeps.ReadInventoryItem = uc.ReadInventoryItem.Execute
			}
			if uc := useCases.Inventory.InventorySerial; uc != nil {
				productDeps.ListInventorySerials = uc.ListInventorySerials.Execute
			}
		}
		// Pricing deps (for variant detail Pricing tab).
		if useCases.Product != nil {
			if uc := useCases.Product.ProductPlan; uc != nil {
				productDeps.ListProductPlans = uc.ListProductPlans.Execute
			}
		}
		if useCases.Subscription != nil {
			if uc := useCases.Subscription.ProductPricePlan; uc != nil {
				productDeps.ListProductPricePlans = uc.ListProductPricePlans.Execute
			}
			if uc := useCases.Subscription.PricePlan; uc != nil {
				productDeps.ListPricePlans = uc.ListPricePlans.Execute
			}
			if uc := useCases.Subscription.PriceSchedule; uc != nil {
				productDeps.ListPriceSchedules = uc.ListPriceSchedules.Execute
			}
			if uc := useCases.Subscription.Plan; uc != nil {
				productDeps.ListPlans = uc.ListPlans.Execute
			}
		}
		wireServiceDashboard(productDeps, useCases)
		productmod.NewModule(productDeps).RegisterRoutes(ctx.Routes)

		// Inventory-flavoured product mount. Reuses the same product module
		// (single view module, Option B from the dual-mount plan) but with
		// Mode="inventory" so the list page filters product_kind
		// IN ('stocked_good','non_stocked_good','consumable'), distinct routes
		// (e.g. /app/inventory/products/list/{status}) and distinct labels
		// sourced from product_inventory.json.
		//
		// Register the inventory-flavoured Product mount on distinct URLs
		// produced by DefaultProductInventoryRoutes. The gate is a
		// defensive check: if a lyngua product_inventory override ever
		// collapses ListURL back onto the service mount, skip the second
		// registration to avoid a ServeMux duplicate-route panic.
		if w.productInventoryRoutes.ListURL != w.productRoutes.ListURL {
			productInventoryDeps := *productDeps
			productInventoryDeps.Routes = w.productInventoryRoutes
			productInventoryDeps.Mode = "inventory"
			productInventoryDeps.Labels = w.productInventoryLabels
			productInventoryDeps.DefaultProductKind = "stocked_good"
			productInventoryDeps.DefaultDeliveryMode = "shipped"
			productInventoryDeps.DefaultTrackingMode = "bulk"
			// Inventory (resold goods) mount exposes two product_kind
			// options so the user picks between stocked vs non-stocked
			// (drop-ship/special order). Consumables belong to the
			// supplies mount and are deliberately excluded here.
			productInventoryDeps.AllowedProductKinds = []string{"stocked_good", "non_stocked_good"}
			// Operation-level RBAC: inventory mount uses "product:*" —
			// historically the default entity, so existing product:*
			// grants keep working on the Products surface without any
			// role-permission migration.
			productInventoryDeps.PermissionEntity = "product"
			productmod.NewModule(&productInventoryDeps).RegisterRoutes(ctx.Routes)
		}

		// Supplies-flavoured product mount. Mode="supplies" narrows the
		// list filter to product_kind = 'consumable', and the routes land
		// under /app/inventory/supplies/* + /action/inventory-supplies/*
		// so it coexists with both the services and inventory mounts on
		// the same ServeMux. Gated only on route distinctness — the same
		// defensive check we use for inventory — so a tier that wipes the
		// supplies route block back onto an existing mount silently drops
		// the registration instead of panicking.
		if w.productSuppliesRoutes.ListURL != w.productRoutes.ListURL &&
			w.productSuppliesRoutes.ListURL != w.productInventoryRoutes.ListURL {
			productSuppliesDeps := *productDeps
			productSuppliesDeps.Routes = w.productSuppliesRoutes
			productSuppliesDeps.Mode = "supplies"
			productSuppliesDeps.Labels = w.productSuppliesLabels
			productSuppliesDeps.DefaultProductKind = "consumable"
			productSuppliesDeps.DefaultDeliveryMode = "shipped"
			productSuppliesDeps.DefaultTrackingMode = "bulk"
			// Supplies mount locks product_kind to "consumable" (single
			// option → drawer renders the select disabled).
			productSuppliesDeps.AllowedProductKinds = []string{"consumable"}
			// Operation-level RBAC: supplies mount uses "supplies:*" so a
			// stock-clerk role can be granted Supplies CRUD without any
			// grant on Products or Services.
			productSuppliesDeps.PermissionEntity = "supplies"
			productmod.NewModule(&productSuppliesDeps).RegisterRoutes(ctx.Routes)
		}
	}

	// =====================================================================
	// Product Line module
	// =====================================================================

	if cfg.wantProductLine() {
		if useCases.Product != nil && useCases.Product.Line != nil {
			uc := useCases.Product.Line
			modDeps := &productlinemod.ModuleDeps{
				Routes:       w.productLineRoutes,
				Labels:       w.productLineLabels,
				CommonLabels: ctx.Common,
				TableLabels:  w.centymoTableLabels,
				ListLines:    uc.ListLines.Execute,
				ReadLine:     uc.ReadLine.Execute,
				CreateLine:   uc.CreateLine.Execute,
				UpdateLine:   uc.UpdateLine.Execute,
				DeleteLine:   uc.DeleteLine.Execute,
			}
			if w.refChecker != nil {
				modDeps.GetInUseIDs = w.refChecker.GetLineInUseIDs
			}
			productlinemod.NewModule(modDeps).RegisterRoutes(ctx.Routes)

			// Inventory-mount ProductLine second registration on distinct URLs.
			// Gate: if a lyngua product_line_inventory override ever collapses
			// ListURL back onto the services mount, skip to avoid a ServeMux
			// duplicate-route panic.
			if w.productLineInventoryRoutes.ListURL != w.productLineRoutes.ListURL {
				productLineInventoryDeps := &productlinemod.ModuleDeps{
					Routes:       w.productLineInventoryRoutes,
					Labels:       w.productLineLabels,
					CommonLabels: ctx.Common,
					TableLabels:  w.centymoTableLabels,
					ListLines:    uc.ListLines.Execute,
					ReadLine:     uc.ReadLine.Execute,
					CreateLine:   uc.CreateLine.Execute,
					UpdateLine:   uc.UpdateLine.Execute,
					DeleteLine:   uc.DeleteLine.Execute,
				}
				if w.refChecker != nil {
					productLineInventoryDeps.GetInUseIDs = w.refChecker.GetLineInUseIDs
				}
				productlinemod.NewModule(productLineInventoryDeps).RegisterRoutes(ctx.Routes)
			}
		}
	}
}
