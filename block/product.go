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
	productmodmodule "github.com/erniealice/centymo-golang/domain/product/product/module"

	"github.com/erniealice/espyna-golang/reference"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	productdom "github.com/erniealice/centymo-golang/domain/product"
	productlinemod "github.com/erniealice/centymo-golang/domain/product/product/line"
)

// productWiring holds everything wireProductModules needs from the surrounding
// Block() scope. More than 6 fields → struct. Kept private; never re-exported.
type productWiring struct {
	db         centymo.DataSource
	refChecker reference.Checker
	// Image + attachment ops
	uploadImage      func(context.Context, string, string, []byte, string) error
	uploadFile       func(context.Context, string, string, []byte, string) error
	downloadFile     func(context.Context, string, string) ([]byte, error)
	listAttachments  func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	readAttachment   func(context.Context, *attachmentpb.ReadAttachmentRequest) (*attachmentpb.ReadAttachmentResponse, error)
	deleteAttachment func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID  func() string
	// Routes
	productRoutes              productdom.ProductRoutes
	productInventoryRoutes     productdom.ProductRoutes
	productSuppliesRoutes      productdom.ProductRoutes
	productLineRoutes          productdom.ProductLineRoutes
	productLineInventoryRoutes productdom.ProductLineRoutes
	// Labels
	productLabels          productdom.ProductLabels
	productInventoryLabels productdom.ProductLabels
	productSuppliesLabels  productdom.ProductLabels
	productLineLabels      productdom.ProductLineLabels
	centymoTableLabels     types.TableLabels
}

// wireProductModules lifts the bodies of the two product-related `if cfg.wantXxx()`
// branches (Product with 3 mounts, ProductLine with 2 mounts) from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once at the position where
// the product wiring used to be.
func wireProductModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *UseCases, w productWiring) {
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

		productDeps := &productmodmodule.ModuleDeps{
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
			DownloadFile:     w.downloadFile,
			ListAttachments:  w.listAttachments,
			CreateAttachment: w.createAttachment,
			ReadAttachment:   w.readAttachment,
			DeleteAttachment: w.deleteAttachment,
			NewID:            w.newAttachmentID,
		}
		if useCases.Product.ListProducts != nil {
			productDeps.ListProducts = useCases.Product.ListProducts
			productDeps.ReadProduct = useCases.Product.ReadProduct
			productDeps.CreateProduct = useCases.Product.CreateProduct
			productDeps.UpdateProduct = useCases.Product.UpdateProduct
			productDeps.DeleteProduct = useCases.Product.DeleteProduct
		}
		if useCases.Product.ListProductVariants != nil {
			productDeps.ListProductVariants = useCases.Product.ListProductVariants
			productDeps.ReadProductVariant = useCases.Product.ReadProductVariant
			productDeps.CreateProductVariant = useCases.Product.CreateProductVariant
			productDeps.UpdateProductVariant = useCases.Product.UpdateProductVariant
			productDeps.DeleteProductVariant = useCases.Product.DeleteProductVariant
		}
		if useCases.Product.ListProductVariantOptions != nil {
			productDeps.ListProductVariantOptions = useCases.Product.ListProductVariantOptions
			productDeps.CreateProductVariantOption = useCases.Product.CreateProductVariantOption
		}
		if useCases.Product.ListProductOptions != nil {
			productDeps.ListProductOptions = useCases.Product.ListProductOptions
			productDeps.ReadProductOption = useCases.Product.ReadProductOption
			productDeps.CreateProductOption = useCases.Product.CreateProductOption
			productDeps.UpdateProductOption = useCases.Product.UpdateProductOption
			productDeps.DeleteProductOption = useCases.Product.DeleteProductOption
		}
		if useCases.Product.ListProductOptionValues != nil {
			productDeps.ListProductOptionValues = useCases.Product.ListProductOptionValues
			productDeps.ReadProductOptionValue = useCases.Product.ReadProductOptionValue
			productDeps.CreateProductOptionValue = useCases.Product.CreateProductOptionValue
			productDeps.UpdateProductOptionValue = useCases.Product.UpdateProductOptionValue
			productDeps.DeleteProductOptionValue = useCases.Product.DeleteProductOptionValue
		}
		if useCases.Product.ListProductAttributes != nil {
			productDeps.ListProductAttributes = useCases.Product.ListProductAttributes
			productDeps.CreateProductAttribute = useCases.Product.CreateProductAttribute
			productDeps.DeleteProductAttribute = useCases.Product.DeleteProductAttribute
		}
		if useCases.Product.ListLines != nil {
			productDeps.ListLines = useCases.Product.ListLines
		}
		if useCases.Product.ListProductLines != nil {
			productDeps.ListProductLines = useCases.Product.ListProductLines
			productDeps.CreateProductLine = useCases.Product.CreateProductLine
			productDeps.UpdateProductLine = useCases.Product.UpdateProductLine
			productDeps.DeleteProductLine = useCases.Product.DeleteProductLine
		}
		if useCases.Product.ListProductVariantImages != nil {
			productDeps.ListProductVariantImages = useCases.Product.ListProductVariantImages
			productDeps.CreateProductVariantImage = useCases.Product.CreateProductVariantImage
			productDeps.DeleteProductVariantImage = useCases.Product.DeleteProductVariantImage
		}
		// Common Attribute (for attribute dropdowns in product detail)
		if useCases.Common.ListAttributes != nil {
			productDeps.ListAttributes = useCases.Common.ListAttributes
			productDeps.ReadAttribute = useCases.Common.ReadAttribute
		}
		// Inventory (for variant detail page + variant stock detail)
		if useCases.Inventory.ListInventoryItems != nil {
			productDeps.ListInventoryItems = useCases.Inventory.ListInventoryItems
			productDeps.ReadInventoryItem = useCases.Inventory.ReadInventoryItem
		}
		if useCases.Inventory.ListInventorySerials != nil {
			productDeps.ListInventorySerials = useCases.Inventory.ListInventorySerials
		}
		// Pricing deps (for variant detail Pricing tab).
		if useCases.Product.ListProductPlans != nil {
			productDeps.ListProductPlans = useCases.Product.ListProductPlans
		}
		if useCases.PricePlan.ListProductPricePlans != nil {
			productDeps.ListProductPricePlans = useCases.PricePlan.ListProductPricePlans
		}
		if useCases.PricePlan.ListPricePlans != nil {
			productDeps.ListPricePlans = useCases.PricePlan.ListPricePlans
		}
		if useCases.PriceSchedule.ListPriceSchedules != nil {
			productDeps.ListPriceSchedules = useCases.PriceSchedule.ListPriceSchedules
		}
		if useCases.Plan.ListPlans != nil {
			productDeps.ListPlans = useCases.Plan.ListPlans
		}
		wireServiceDashboard(productDeps, useCases)
		productModule := productmodmodule.NewModule(productDeps)
		productModule.RegisterRoutes(ctx.Routes)
		// Attachment preview/download streams raw bytes — registered at the
		// block layer because RegisterRoutes only handles view.View, not
		// http.HandlerFunc. Mirrors the subscription pattern in subscription.go.
		if w.downloadFile != nil && w.readAttachment != nil && w.productRoutes.AttachmentDownloadURL != "" {
			handleFunc(ctx.Routes, "GET", w.productRoutes.AttachmentDownloadURL, productModule.AttachmentDownload)
		}

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
			productmodmodule.NewModule(&productInventoryDeps).RegisterRoutes(ctx.Routes)
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
			productmodmodule.NewModule(&productSuppliesDeps).RegisterRoutes(ctx.Routes)
		}
	}

	// =====================================================================
	// Product Line module
	// =====================================================================

	if cfg.wantProductLine() {
		if useCases.Product.ListLines != nil {
			modDeps := &productlinemod.ModuleDeps{
				Routes:       w.productLineRoutes,
				Labels:       w.productLineLabels,
				CommonLabels: ctx.Common,
				TableLabels:  w.centymoTableLabels,
				ListLines:    useCases.Product.ListLines,
				ReadLine:     useCases.Product.ReadLine,
				CreateLine:   useCases.Product.CreateLine,
				UpdateLine:   useCases.Product.UpdateLine,
				DeleteLine:   useCases.Product.DeleteLine,
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
					ListLines:    useCases.Product.ListLines,
					ReadLine:     useCases.Product.ReadLine,
					CreateLine:   useCases.Product.CreateLine,
					UpdateLine:   useCases.Product.UpdateLine,
					DeleteLine:   useCases.Product.DeleteLine,
				}
				if w.refChecker != nil {
					productLineInventoryDeps.GetInUseIDs = w.refChecker.GetLineInUseIDs
				}
				productlinemod.NewModule(productLineInventoryDeps).RegisterRoutes(ctx.Routes)
			}
		}
	}
}
