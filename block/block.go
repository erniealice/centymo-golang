// Package block exposes centymo.Block() — the Lego composition entry point
// for the centymo commerce domain (inventory, revenue, product, pricelist,
// plan, subscription, collection, disbursement, expenditure, and inline report
// routes). Consumer apps import this package and optionally alias it:
//
//	import centymoblock "github.com/erniealice/centymo-golang/block"
//	// ...
//	centymoblock.Block()               // all modules
//	centymoblock.Block(
//	    centymoblock.WithRevenue(),
//	    centymoblock.WithProduct(),
//	)                                   // selective modules
//
// This package lives in a sub-package (not the centymo root) to avoid a Go
// import cycle: centymo/views/* imports centymo (root) for route/label types,
// so Block() cannot live in the root package while also importing centymo/views/*.
package block

import (
	"context"
	"fmt"
	"log"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"

	consumer "github.com/erniealice/espyna-golang/consumer"
	"github.com/erniealice/espyna-golang/contrib/postgres/reference"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/template"

	"github.com/erniealice/hybra-golang/views/attachment"
	templateview "github.com/erniealice/hybra-golang/views/template"

	centymo "github.com/erniealice/centymo-golang"
	collectionmod "github.com/erniealice/centymo-golang/views/collection"
	disbursementmod "github.com/erniealice/centymo-golang/views/disbursement"
	expendituremod "github.com/erniealice/centymo-golang/views/expenditure"
	inventorymod "github.com/erniealice/centymo-golang/views/inventory"
	planaction "github.com/erniealice/centymo-golang/views/plan/action"
	plandetail "github.com/erniealice/centymo-golang/views/plan/detail"
	planlist "github.com/erniealice/centymo-golang/views/plan/list"
	pricelistmod "github.com/erniealice/centymo-golang/views/pricelist"
	productmod "github.com/erniealice/centymo-golang/views/product"
	revenuemod "github.com/erniealice/centymo-golang/views/revenue"
	subscriptionaction "github.com/erniealice/centymo-golang/views/subscription/action"
	subscriptiondetail "github.com/erniealice/centymo-golang/views/subscription/detail"
	subscriptionlist "github.com/erniealice/centymo-golang/views/subscription/list"
	reportmod "github.com/erniealice/fycha-golang/views/reports"
)

// ---------------------------------------------------------------------------
// routeRegistrarFull — optional extension for raw http.HandlerFunc routes
// ---------------------------------------------------------------------------

// routeRegistrarFull extends pyeza.RouteRegistrar with HandleFunc support.
// Consumer apps whose RouteRegistrar implements this interface (e.g. service-admin's
// RouteRegistry) can register raw http.HandlerFunc routes for CSV export, invoice
// download, and email dispatch. Apps that do not implement HandleFunc will skip
// those routes with a log warning.
type routeRegistrarFull interface {
	pyeza.RouteRegistrar
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

// handleFunc is a nil-safe helper that registers an http.HandlerFunc route if the
// RouteRegistrar supports it, otherwise logs a warning and skips.
func handleFunc(r pyeza.RouteRegistrar, method, path string, handler http.HandlerFunc) {
	if handler == nil {
		return
	}
	if full, ok := r.(routeRegistrarFull); ok {
		full.HandleFunc(method, path, handler)
		return
	}
	log.Printf("centymo.Block: RouteRegistrar does not support HandleFunc — skipping %s %s", method, path)
}

// ---------------------------------------------------------------------------
// BlockOption — per-module granular selection
// ---------------------------------------------------------------------------

// BlockOption enables specific centymo sub-modules within Block().
type BlockOption func(*blockConfig)

type blockConfig struct {
	enableAll    bool
	inventory    bool
	revenue      bool
	product      bool
	priceList    bool
	plan         bool
	subscription bool
	collection   bool
	disbursement bool
	expenditure  bool
}

func WithInventory() BlockOption    { return func(c *blockConfig) { c.inventory = true } }
func WithRevenue() BlockOption      { return func(c *blockConfig) { c.revenue = true } }
func WithProduct() BlockOption      { return func(c *blockConfig) { c.product = true } }
func WithPriceList() BlockOption    { return func(c *blockConfig) { c.priceList = true } }
func WithPlan() BlockOption         { return func(c *blockConfig) { c.plan = true } }
func WithSubscription() BlockOption { return func(c *blockConfig) { c.subscription = true } }
func WithCollection() BlockOption   { return func(c *blockConfig) { c.collection = true } }
func WithDisbursement() BlockOption { return func(c *blockConfig) { c.disbursement = true } }
func WithExpenditure() BlockOption  { return func(c *blockConfig) { c.expenditure = true } }

func (c *blockConfig) wantInventory() bool    { return c.enableAll || c.inventory }
func (c *blockConfig) wantRevenue() bool      { return c.enableAll || c.revenue }
func (c *blockConfig) wantProduct() bool      { return c.enableAll || c.product }
func (c *blockConfig) wantPriceList() bool    { return c.enableAll || c.priceList }
func (c *blockConfig) wantPlan() bool         { return c.enableAll || c.plan }
func (c *blockConfig) wantSubscription() bool { return c.enableAll || c.subscription }
func (c *blockConfig) wantCollection() bool   { return c.enableAll || c.collection }
func (c *blockConfig) wantDisbursement() bool { return c.enableAll || c.disbursement }
func (c *blockConfig) wantExpenditure() bool  { return c.enableAll || c.expenditure }

// ---------------------------------------------------------------------------
// Block — the main Lego entry point
// ---------------------------------------------------------------------------

// Block registers centymo domain modules (commerce: inventory, revenue, product,
// pricelist, plan, subscription, collection, disbursement, expenditure, and inline
// report routes). Call with no options to register ALL modules. Call with specific
// WithX() options for a subset.
func Block(opts ...BlockOption) pyeza.AppOption {
	cfg := &blockConfig{enableAll: len(opts) == 0}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(ctx *pyeza.AppContext) error {
		// --- Type-assert translations ---
		translations, ok := ctx.Translations.(*lynguaV1.TranslationProvider)
		if !ok || translations == nil {
			return fmt.Errorf("centymo.Block: ctx.Translations must be *lynguaV1.TranslationProvider")
		}

		// --- Type-assert use cases ---
		useCases, ok := ctx.UseCases.(*consumer.UseCases)
		if !ok || useCases == nil {
			return fmt.Errorf("centymo.Block: ctx.UseCases must be *consumer.UseCases")
		}

		// --- Type-assert DB ---
		db, ok := ctx.DB.(centymo.DataSource)
		if !ok || db == nil {
			return fmt.Errorf("centymo.Block: ctx.DB must implement centymo.DataSource")
		}

		// --- Type-assert reference checker (optional — nil-safe) ---
		var refChecker *reference.Checker
		if ctx.RefChecker != nil {
			refChecker, _ = ctx.RefChecker.(*reference.Checker)
		}

		// --- Type-assert attachment operations ---
		uploadFile, _ := ctx.UploadFile.(func(context.Context, string, string, []byte, string) error)
		listAttachments, _ := ctx.ListAttachments.(func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error))
		createAttachment, _ := ctx.CreateAttachment.(func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error))
		deleteAttachment, _ := ctx.DeleteAttachment.(func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error))
		newAttachmentID, _ := ctx.NewAttachmentID.(func() string)

		// --- Type-assert storage/email/doc operations ---
		uploadImage, _ := ctx.UploadImage.(func(context.Context, string, string, []byte, string) error)
		uploadTemplate, _ := ctx.UploadTemplate.(func(context.Context, string, string, []byte, string) error)
		sendEmail, _ := ctx.SendEmail.(func(context.Context, []string, string, string, string, string, []byte) error)
		generateDoc, _ := ctx.GenerateDoc.(func([]byte, map[string]any) ([]byte, error))
		listDocTemplates, _ := ctx.ListDocTemplates.(func(context.Context, *documenttemplatepb.ListDocumentTemplatesRequest) (*documenttemplatepb.ListDocumentTemplatesResponse, error))
		createDocTemplate, _ := ctx.CreateDocTemplate.(func(context.Context, *documenttemplatepb.CreateDocumentTemplateRequest) (*documenttemplatepb.CreateDocumentTemplateResponse, error))
		updateDocTemplate, _ := ctx.UpdateDocTemplate.(func(context.Context, *documenttemplatepb.UpdateDocumentTemplateRequest) (*documenttemplatepb.UpdateDocumentTemplateResponse, error))
		deleteDocTemplate, _ := ctx.DeleteDocTemplate.(func(context.Context, *documenttemplatepb.DeleteDocumentTemplateRequest) (*documenttemplatepb.DeleteDocumentTemplateResponse, error))

		// --- Shared table labels ---
		centymoTableLabels := centymo.MapTableLabels(ctx.Common)

		// --- Load routes (defaults + optional lyngua overrides) ---
		inventoryRoutes := centymo.DefaultInventoryRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "inventory", &inventoryRoutes)

		revenueRoutes := centymo.DefaultRevenueRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "revenue", &revenueRoutes)

		productRoutes := centymo.DefaultProductRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "product", &productRoutes)

		priceListRoutes := centymo.DefaultPriceListRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "price_list", &priceListRoutes)

		planRoutes := centymo.DefaultPlanRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "plan", &planRoutes)

		subscriptionRoutes := centymo.DefaultSubscriptionRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "subscription", &subscriptionRoutes)

		collectionRoutes := centymo.DefaultCollectionRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "treasury_collection", &collectionRoutes)

		disbursementRoutes := centymo.DefaultDisbursementRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "treasury_disbursement", &disbursementRoutes)

		expenditureRoutes := centymo.DefaultExpenditureRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "expenditure", &expenditureRoutes)

		// --- Load labels ---
		var inventoryLabels centymo.InventoryLabels
		if err := translations.LoadPath("en", ctx.BusinessType, "inventory.json", "inventory", &inventoryLabels); err != nil {
			log.Printf("centymo.Block: warning loading inventory labels: %v", err)
		}

		var revenueLabels centymo.RevenueLabels
		if err := translations.LoadPath("en", ctx.BusinessType, "revenue.json", "revenue", &revenueLabels); err != nil {
			log.Printf("centymo.Block: warning loading revenue labels: %v", err)
		}

		var productLabels centymo.ProductLabels
		if err := translations.LoadPath("en", ctx.BusinessType, "product.json", "product", &productLabels); err != nil {
			log.Printf("centymo.Block: warning loading product labels: %v", err)
		}

		var priceListLabels centymo.PriceListLabels
		if err := translations.LoadPath("en", ctx.BusinessType, "pricelist.json", "pricelist", &priceListLabels); err != nil {
			log.Printf("centymo.Block: warning loading pricelist labels: %v", err)
		}

		var expenditureLabels centymo.ExpenditureLabels
		if err := translations.LoadPath("en", ctx.BusinessType, "expenditure.json", "expenditure", &expenditureLabels); err != nil {
			log.Printf("centymo.Block: warning loading expenditure labels: %v", err)
		}

		collectionLabels := centymo.DefaultCollectionLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "collection.json", "collection", &collectionLabels)

		disbursementLabels := centymo.DefaultDisbursementLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "disbursement.json", "disbursement", &disbursementLabels)

		planLabels := centymo.DefaultPlanLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "plan.json", "plan", &planLabels)

		subscriptionLabels := centymo.DefaultSubscriptionLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "subscription.json", "subscription", &subscriptionLabels)

		// =====================================================================
		// Inventory module
		// =====================================================================

		if cfg.wantInventory() {
			invDeps := &inventorymod.ModuleDeps{
				Routes:       inventoryRoutes,
				SqlDB:        ctx.SqlDB,
				Labels:       inventoryLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
				// SetItemActive uses raw DB update (proto3 omits false booleans)
				SetItemActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "inventory_item", id, map[string]any{"active": active})
					return err
				},
				// Attachments
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
				NewID:            newAttachmentID,
			}
			if useCases.Inventory != nil {
				if uc := useCases.Inventory.InventoryItem; uc != nil {
					invDeps.ListInventoryItems = uc.ListInventoryItems.Execute
					invDeps.CreateInventoryItem = uc.CreateInventoryItem.Execute
					invDeps.ReadInventoryItem = uc.ReadInventoryItem.Execute
					invDeps.UpdateInventoryItem = uc.UpdateInventoryItem.Execute
					invDeps.DeleteInventoryItem = uc.DeleteInventoryItem.Execute
				}
				if uc := useCases.Inventory.InventorySerial; uc != nil {
					invDeps.ListInventorySerials = uc.ListInventorySerials.Execute
					invDeps.CreateInventorySerial = uc.CreateInventorySerial.Execute
					invDeps.ReadInventorySerial = uc.ReadInventorySerial.Execute
					invDeps.UpdateInventorySerial = uc.UpdateInventorySerial.Execute
					invDeps.DeleteInventorySerial = uc.DeleteInventorySerial.Execute
				}
				if uc := useCases.Inventory.InventoryTransaction; uc != nil {
					invDeps.ListInventoryTransactions = uc.ListInventoryTransactions.Execute
					invDeps.CreateInventoryTransaction = uc.CreateInventoryTransaction.Execute
				}
				if uc := useCases.Inventory.InventoryDepreciation; uc != nil {
					invDeps.ListInventoryDepreciations = uc.ListInventoryDepreciations.Execute
					invDeps.CreateInventoryDepreciation = uc.CreateInventoryDepreciation.Execute
					invDeps.ReadInventoryDepreciation = uc.ReadInventoryDepreciation.Execute
					invDeps.UpdateInventoryDepreciation = uc.UpdateInventoryDepreciation.Execute
				}
			}
			// Cross-domain: product options + locations for inventory item views
			if useCases.Product != nil {
				if uc := useCases.Product.Product; uc != nil {
					invDeps.ReadProduct = uc.ReadProduct.Execute
				}
				if uc := useCases.Product.ProductVariantOption; uc != nil {
					invDeps.ListProductVariantOptions = uc.ListProductVariantOptions.Execute
				}
				if uc := useCases.Product.ProductOptionValue; uc != nil {
					invDeps.ListProductOptionValues = uc.ListProductOptionValues.Execute
				}
				if uc := useCases.Product.ProductOption; uc != nil {
					invDeps.ListProductOptions = uc.ListProductOptions.Execute
				}
			}
			if useCases.Entity != nil && useCases.Entity.Location != nil {
				invDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
			}

			invMod := inventorymod.NewModule(invDeps)
			invMod.RegisterRoutes(ctx.Routes)
			// CSV export is a raw http.HandlerFunc (bypasses view/template layer)
			handleFunc(ctx.Routes, "GET", inventoryRoutes.MovementsExportURL, invMod.MovementsExport)
		}

		// =====================================================================
		// Revenue module
		// =====================================================================

		if cfg.wantRevenue() {
			revDeps := &revenuemod.ModuleDeps{
				Routes:       revenueRoutes,
				DB:           db,
				Labels:       revenueLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
				// Document generation + template CRUD
				GenerateDoc:            generateDoc,
				ListDocumentTemplates:  listDocTemplates,
				CreateDocumentTemplate: createDocTemplate,
				UpdateDocumentTemplate: updateDocTemplate,
				DeleteDocumentTemplate: deleteDocTemplate,
				UploadTemplate:         uploadTemplate,
				SendEmail:              sendEmail,
				// Attachments
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
				NewID:            newAttachmentID,
			}
			// Revenue CRUD + list page data
			if useCases.Revenue != nil && useCases.Revenue.Revenue != nil {
				revDeps.GetListPageData = useCases.Revenue.Revenue.GetRevenueListPageData.Execute
				revDeps.CreateRevenue = useCases.Revenue.Revenue.CreateRevenue.Execute
				revDeps.ReadRevenue = useCases.Revenue.Revenue.ReadRevenue.Execute
				revDeps.UpdateRevenue = useCases.Revenue.Revenue.UpdateRevenue.Execute
				revDeps.DeleteRevenue = useCases.Revenue.Revenue.DeleteRevenue.Execute
			}
			// Revenue Line Item CRUD
			if useCases.Revenue != nil && useCases.Revenue.RevenueLineItem != nil {
				uc := useCases.Revenue.RevenueLineItem
				revDeps.CreateRevenueLineItem = uc.CreateRevenueLineItem.Execute
				revDeps.ReadRevenueLineItem = uc.ReadRevenueLineItem.Execute
				revDeps.UpdateRevenueLineItem = uc.UpdateRevenueLineItem.Execute
				revDeps.DeleteRevenueLineItem = uc.DeleteRevenueLineItem.Execute
				revDeps.ListRevenueLineItems = uc.ListRevenueLineItems.Execute
			}
			// Inventory (for stock deduction on status change)
			if useCases.Inventory != nil {
				if uc := useCases.Inventory.InventoryItem; uc != nil {
					revDeps.ReadInventoryItem = uc.ReadInventoryItem.Execute
					revDeps.UpdateInventoryItem = uc.UpdateInventoryItem.Execute
					revDeps.ListInventoryItems = uc.ListInventoryItems.Execute
				}
				if uc := useCases.Inventory.InventorySerial; uc != nil {
					revDeps.UpdateInventorySerial = uc.UpdateInventorySerial.Execute
				}
				if uc := useCases.Inventory.InventorySerialHistory; uc != nil {
					revDeps.CreateInventorySerialHistory = uc.CreateInventorySerialHistory.Execute
				}
			}

			revenueMod := revenuemod.NewModule(revDeps)
			revenueMod.RegisterRoutes(ctx.Routes)
			// Invoice download is http.HandlerFunc (bypasses view/template layer)
			handleFunc(ctx.Routes, "GET", revenueRoutes.InvoiceDownloadURL, revenueMod.InvoiceDownload)
			// Send email is http.HandlerFunc (bypasses view/template layer)
			handleFunc(ctx.Routes, "POST", revenueRoutes.SendEmailURL, revenueMod.SendEmailHandler)
		}

		// =====================================================================
		// Product module
		// =====================================================================

		if cfg.wantProduct() {
			var getProductInUseIDs func(context.Context, []string) (map[string]bool, error)
			if refChecker != nil {
				getProductInUseIDs = refChecker.GetProductInUseIDs
			}

			productDeps := &productmod.ModuleDeps{
				Routes:       productRoutes,
				DB:           db,
				Labels:       productLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
				GetInUseIDs:  getProductInUseIDs,
				// SetProductActive uses raw DB update (proto3 omits false booleans)
				SetProductActive: func(fctx context.Context, id string, active bool) error {
					_, err := db.Update(fctx, "product", id, map[string]any{"active": active})
					return err
				},
				// Image upload (product variant images)
				UploadImage: uploadImage,
				// Attachments
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
				NewID:            newAttachmentID,
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
			productmod.NewModule(productDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// PriceList module
		// =====================================================================

		if cfg.wantPriceList() {
			var getPriceListInUseIDs func(context.Context, []string) (map[string]bool, error)
			if refChecker != nil {
				getPriceListInUseIDs = refChecker.GetPriceListInUseIDs
			}

			pricelistmod.NewModule(&pricelistmod.ModuleDeps{
				Routes:             priceListRoutes,
				Labels:             priceListLabels,
				CommonLabels:       ctx.Common,
				TableLabels:        centymoTableLabels,
				GetInUseIDs:        getPriceListInUseIDs,
				ListPriceLists:     useCases.Product.PriceList.ListPriceLists.Execute,
				ReadPriceList:      useCases.Product.PriceList.ReadPriceList.Execute,
				CreatePriceList:    useCases.Product.PriceList.CreatePriceList.Execute,
				UpdatePriceList:    useCases.Product.PriceList.UpdatePriceList.Execute,
				DeletePriceList:    useCases.Product.PriceList.DeletePriceList.Execute,
				ListPriceProducts:  useCases.Product.PriceProduct.ListPriceProducts.Execute,
				CreatePriceProduct: useCases.Product.PriceProduct.CreatePriceProduct.Execute,
				DeletePriceProduct: useCases.Product.PriceProduct.DeletePriceProduct.Execute,
				ListProducts:       useCases.Product.Product.ListProducts.Execute,
				// Attachments
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
				NewID:            newAttachmentID,
			}).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Plan (inline — not a module, uses planlist/planaction/plandetail directly)
		// =====================================================================

		if cfg.wantPlan() {
			planListDeps := &planlist.ListViewDeps{
				Routes:       planRoutes,
				Labels:       planLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
			}
			if useCases.Subscription != nil && useCases.Subscription.Plan != nil {
				planListDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
			}
			ctx.Routes.GET(planRoutes.ListURL, planlist.NewView(planListDeps))

			// Plan CRUD actions
			if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.CreatePlan != nil {
				planActionDeps := &planaction.Deps{
					Routes:     planRoutes,
					Labels:     planLabels,
					CreatePlan: useCases.Subscription.Plan.CreatePlan.Execute,
					ReadPlan:   useCases.Subscription.Plan.ReadPlan.Execute,
					UpdatePlan: useCases.Subscription.Plan.UpdatePlan.Execute,
					DeletePlan: useCases.Subscription.Plan.DeletePlan.Execute,
				}
				ctx.Routes.GET(planRoutes.AddURL, planaction.NewAddAction(planActionDeps))
				ctx.Routes.POST(planRoutes.AddURL, planaction.NewAddAction(planActionDeps))
				ctx.Routes.GET(planRoutes.EditURL, planaction.NewEditAction(planActionDeps))
				ctx.Routes.POST(planRoutes.EditURL, planaction.NewEditAction(planActionDeps))
				ctx.Routes.POST(planRoutes.DeleteURL, planaction.NewDeleteAction(planActionDeps))
			}

			// Plan detail page + tab action
			if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
				planDetailDeps := &plandetail.DetailViewDeps{
					Routes:       planRoutes,
					ReadPlan:     useCases.Subscription.Plan.ReadPlan.Execute,
					Labels:       planLabels,
					CommonLabels: ctx.Common,
					TableLabels:  centymoTableLabels,
					AttachmentOps: attachment.AttachmentOps{
						UploadFile:       uploadFile,
						ListAttachments:  listAttachments,
						CreateAttachment: createAttachment,
						DeleteAttachment: deleteAttachment,
						NewAttachmentID:  newAttachmentID,
					},
				}
				if useCases.Product != nil && useCases.Product.ProductPlan != nil {
					planDetailDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
				}
				if useCases.Subscription.PricePlan != nil {
					planDetailDeps.ListPricePlans = useCases.Subscription.PricePlan.ListPricePlans.Execute
				}
				ctx.Routes.GET(planRoutes.DetailURL, plandetail.NewView(planDetailDeps))
				ctx.Routes.GET(planRoutes.TabActionURL, plandetail.NewTabAction(planDetailDeps))
				// Plan attachments
				if uploadFile != nil {
					ctx.Routes.GET(planRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planDetailDeps))
					ctx.Routes.POST(planRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planDetailDeps))
					ctx.Routes.POST(planRoutes.AttachmentDeleteURL, plandetail.NewAttachmentDeleteAction(planDetailDeps))
				}
			}
		}

		// =====================================================================
		// Subscription (inline — not a module, uses subscriptionlist/subscriptionaction/subscriptiondetail directly)
		// =====================================================================

		if cfg.wantSubscription() {
			subListDeps := &subscriptionlist.ListViewDeps{
				Routes:       subscriptionRoutes,
				Labels:       subscriptionLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
			}
			if useCases.Subscription != nil && useCases.Subscription.Subscription != nil {
				subListDeps.GetSubscriptionListPageData = useCases.Subscription.Subscription.GetSubscriptionListPageData.Execute
			}
			ctx.Routes.GET(subscriptionRoutes.ListURL, subscriptionlist.NewView(subListDeps))

			// Subscription CRUD actions
			if useCases.Subscription != nil && useCases.Subscription.Subscription != nil && useCases.Subscription.Subscription.CreateSubscription != nil {
				subActionDeps := &subscriptionaction.Deps{
					Routes:             subscriptionRoutes,
					Labels:             subscriptionLabels,
					CreateSubscription: useCases.Subscription.Subscription.CreateSubscription.Execute,
					ReadSubscription:   useCases.Subscription.Subscription.ReadSubscription.Execute,
					UpdateSubscription: useCases.Subscription.Subscription.UpdateSubscription.Execute,
					DeleteSubscription: useCases.Subscription.Subscription.DeleteSubscription.Execute,
				}
				if useCases.Entity != nil && useCases.Entity.Client != nil {
					subActionDeps.ListClients = useCases.Entity.Client.ListClients.Execute
					if useCases.Entity.Client.SearchClientsByName != nil {
						subActionDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName.Execute
					}
				}
				if useCases.Subscription.Plan != nil {
					subActionDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
					if useCases.Subscription.Plan.SearchPlansByName != nil {
						subActionDeps.SearchPlansByName = useCases.Subscription.Plan.SearchPlansByName.Execute
					}
				}
				ctx.Routes.GET(subscriptionRoutes.AddURL, subscriptionaction.NewAddAction(subActionDeps))
				ctx.Routes.POST(subscriptionRoutes.AddURL, subscriptionaction.NewAddAction(subActionDeps))
				ctx.Routes.GET(subscriptionRoutes.EditURL, subscriptionaction.NewEditAction(subActionDeps))
				ctx.Routes.POST(subscriptionRoutes.EditURL, subscriptionaction.NewEditAction(subActionDeps))
				ctx.Routes.POST(subscriptionRoutes.DeleteURL, subscriptionaction.NewDeleteAction(subActionDeps))
				// Auto-complete search (http.HandlerFunc — uses HandleFunc, not GET)
				handleFunc(ctx.Routes, "GET", subscriptionRoutes.SearchClientURL, subscriptionaction.NewSearchClientsAction(subActionDeps))
				handleFunc(ctx.Routes, "GET", subscriptionRoutes.SearchPlanURL, subscriptionaction.NewSearchPlansAction(subActionDeps))
			}

			// Subscription detail page + tab action
			if useCases.Subscription != nil && useCases.Subscription.Subscription != nil && useCases.Subscription.Subscription.ReadSubscription != nil {
				subDetailDeps := &subscriptiondetail.DetailViewDeps{
					Routes:           subscriptionRoutes,
					ReadSubscription: useCases.Subscription.Subscription.ReadSubscription.Execute,
					Labels:           subscriptionLabels,
					CommonLabels:     ctx.Common,
					TableLabels:      centymoTableLabels,
					AttachmentOps: attachment.AttachmentOps{
						UploadFile:       uploadFile,
						ListAttachments:  listAttachments,
						CreateAttachment: createAttachment,
						DeleteAttachment: deleteAttachment,
						NewAttachmentID:  newAttachmentID,
					},
				}
				ctx.Routes.GET(subscriptionRoutes.DetailURL, subscriptiondetail.NewView(subDetailDeps))
				ctx.Routes.GET(subscriptionRoutes.TabActionURL, subscriptiondetail.NewTabAction(subDetailDeps))
				// Subscription attachments
				if uploadFile != nil {
					ctx.Routes.GET(subscriptionRoutes.AttachmentUploadURL, subscriptiondetail.NewAttachmentUploadAction(subDetailDeps))
					ctx.Routes.POST(subscriptionRoutes.AttachmentUploadURL, subscriptiondetail.NewAttachmentUploadAction(subDetailDeps))
					ctx.Routes.POST(subscriptionRoutes.AttachmentDeleteURL, subscriptiondetail.NewAttachmentDeleteAction(subDetailDeps))
				}
			}
		}

		// =====================================================================
		// Collection module (conditional: only when treasury collection use cases are available)
		// =====================================================================

		if cfg.wantCollection() {
			if useCases.Treasury != nil && useCases.Treasury.Collection != nil {
				collectionmod.NewModule(&collectionmod.ModuleDeps{
					Routes:           collectionRoutes,
					Labels:           collectionLabels,
					CommonLabels:     ctx.Common,
					TableLabels:      centymoTableLabels,
					CreateCollection: useCases.Treasury.Collection.CreateCollection.Execute,
					ReadCollection:   useCases.Treasury.Collection.ReadCollection.Execute,
					UpdateCollection: useCases.Treasury.Collection.UpdateCollection.Execute,
					DeleteCollection: useCases.Treasury.Collection.DeleteCollection.Execute,
					ListCollections:  useCases.Treasury.Collection.ListCollections.Execute,
					// Attachments
					UploadFile:       uploadFile,
					ListAttachments:  listAttachments,
					CreateAttachment: createAttachment,
					DeleteAttachment: deleteAttachment,
					NewID:            newAttachmentID,
				}).RegisterRoutes(ctx.Routes)
			}
		}

		// =====================================================================
		// Disbursement module (conditional: only when treasury disbursement use cases are available)
		// =====================================================================

		if cfg.wantDisbursement() {
			if useCases.Treasury != nil && useCases.Treasury.Disbursement != nil {
				disbursementmod.NewModule(&disbursementmod.ModuleDeps{
					Routes:             disbursementRoutes,
					Labels:             disbursementLabels,
					CommonLabels:       ctx.Common,
					TableLabels:        centymoTableLabels,
					CreateDisbursement: useCases.Treasury.Disbursement.CreateDisbursement.Execute,
					ReadDisbursement:   useCases.Treasury.Disbursement.ReadDisbursement.Execute,
					UpdateDisbursement: useCases.Treasury.Disbursement.UpdateDisbursement.Execute,
					DeleteDisbursement: useCases.Treasury.Disbursement.DeleteDisbursement.Execute,
					ListDisbursements:  useCases.Treasury.Disbursement.ListDisbursements.Execute,
					// Attachments
					UploadFile:       uploadFile,
					ListAttachments:  listAttachments,
					CreateAttachment: createAttachment,
					DeleteAttachment: deleteAttachment,
					NewID:            newAttachmentID,
				}).RegisterRoutes(ctx.Routes)
			}
		}

		// =====================================================================
		// Expenditure module (purchase + expense)
		// =====================================================================

		if cfg.wantExpenditure() {
			expDeps := &expendituremod.ModuleDeps{
				Routes:         expenditureRoutes,
				DB:             db,
				Labels:         expenditureLabels,
				TemplateLabels: templateview.DefaultLabels(),
				CommonLabels:   ctx.Common,
				TableLabels:    centymoTableLabels,
				// Document template CRUD
				ListDocumentTemplates:  listDocTemplates,
				CreateDocumentTemplate: createDocTemplate,
				UpdateDocumentTemplate: updateDocTemplate,
				DeleteDocumentTemplate: deleteDocTemplate,
				UploadFile:             uploadTemplate,
			}
			if useCases.Expenditure != nil && useCases.Expenditure.Expenditure != nil {
				expDeps.ListExpenditures = useCases.Expenditure.Expenditure.ListExpenditures.Execute
			}
			expendituremod.NewModule(expDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Inline report routes
		// These centymo-owned report views are registered here because they
		// require centymo SQL labels. Consumer apps that also need the
		// Payables Aging (Supplier route) and Receivables Aging (Client route)
		// reports should register those via the entydad block or directly,
		// since those URL keys live in entydad's route structs.
		// =====================================================================

		if cfg.wantRevenue() {
			// Revenue → Reports → Sales Summary
			ctx.Routes.GET(revenueRoutes.RevenueSummaryURL, reportmod.NewSalesSummaryView(ctx.SqlDB, ctx.Common, ctx.Table))
		}
		if cfg.wantExpenditure() {
			// Purchases → Reports → Purchases Summary
			ctx.Routes.GET(expenditureRoutes.PurchasesSummaryURL, reportmod.NewPurchasesSummaryView(ctx.SqlDB, ctx.Common, ctx.Table))
			// Expenses → Reports → Expenses Summary
			ctx.Routes.GET(expenditureRoutes.ExpensesSummaryURL, reportmod.NewExpensesSummaryView(ctx.SqlDB, ctx.Common, ctx.Table))
		}

		log.Println("  centymo commerce domain initialized")
		return nil
	}
}
