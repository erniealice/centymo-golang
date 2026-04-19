// Package block exposes centymo.Block() — the Lego composition entry point
// for the centymo commerce domain (inventory, revenue, product, product line,
// pricelist, plan, subscription, collection, disbursement, expenditure, and inline report
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

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"

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
	productlinemod "github.com/erniealice/centymo-golang/views/product/line"
	priceplanmod "github.com/erniealice/centymo-golang/views/price_plan"
	priceschedulemod "github.com/erniealice/centymo-golang/views/price_schedule"
	resourcemod "github.com/erniealice/centymo-golang/views/resource"
	revenuemod "github.com/erniealice/centymo-golang/views/revenue"
	subscriptionaction "github.com/erniealice/centymo-golang/views/subscription/action"
	subscriptiondetail "github.com/erniealice/centymo-golang/views/subscription/detail"
	subscriptionlist "github.com/erniealice/centymo-golang/views/subscription/list"
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
	productLine  bool
	pricePlan     bool
	priceSchedule bool
	priceList     bool
	plan         bool
	subscription bool
	collection   bool
	disbursement bool
	expenditure  bool
	resource     bool
}

func WithInventory() BlockOption    { return func(c *blockConfig) { c.inventory = true } }
func WithRevenue() BlockOption      { return func(c *blockConfig) { c.revenue = true } }
func WithProduct() BlockOption      { return func(c *blockConfig) { c.product = true } }
func WithProductLine() BlockOption  { return func(c *blockConfig) { c.productLine = true } }
func WithPricePlan() BlockOption     { return func(c *blockConfig) { c.pricePlan = true } }
func WithPriceSchedule() BlockOption { return func(c *blockConfig) { c.priceSchedule = true } }
func WithPriceList() BlockOption     { return func(c *blockConfig) { c.priceList = true } }
func WithPlan() BlockOption         { return func(c *blockConfig) { c.plan = true } }
func WithSubscription() BlockOption { return func(c *blockConfig) { c.subscription = true } }
func WithCollection() BlockOption   { return func(c *blockConfig) { c.collection = true } }
func WithDisbursement() BlockOption { return func(c *blockConfig) { c.disbursement = true } }
func WithExpenditure() BlockOption  { return func(c *blockConfig) { c.expenditure = true } }
func WithResource() BlockOption     { return func(c *blockConfig) { c.resource = true } }

func (c *blockConfig) wantInventory() bool    { return c.enableAll || c.inventory }
func (c *blockConfig) wantRevenue() bool      { return c.enableAll || c.revenue }
func (c *blockConfig) wantProduct() bool      { return c.enableAll || c.product }
func (c *blockConfig) wantProductLine() bool  { return c.enableAll || c.productLine }
func (c *blockConfig) wantPricePlan() bool     { return c.enableAll || c.pricePlan }
func (c *blockConfig) wantPriceSchedule() bool { return c.enableAll || c.priceSchedule }
func (c *blockConfig) wantPriceList() bool     { return c.enableAll || c.priceList }
func (c *blockConfig) wantPlan() bool         { return c.enableAll || c.plan }
func (c *blockConfig) wantSubscription() bool { return c.enableAll || c.subscription }
func (c *blockConfig) wantCollection() bool   { return c.enableAll || c.collection }
func (c *blockConfig) wantDisbursement() bool { return c.enableAll || c.disbursement }
func (c *blockConfig) wantExpenditure() bool  { return c.enableAll || c.expenditure }
func (c *blockConfig) wantResource() bool     { return c.enableAll || c.resource }

// ---------------------------------------------------------------------------
// Block — the main Lego entry point
// ---------------------------------------------------------------------------

// Block registers centymo domain modules (commerce: inventory, revenue, product,
// product line, pricelist, plan, subscription, collection, disbursement, expenditure, and inline
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

		// Inventory-flavoured product list route overrides. Starts from the
		// namespace-shifted inventory defaults (every URL rewritten from
		// /app/products/* → /app/inventory/products/* and /action/product/*
		// → /action/inventory-product/*) so both mounts can coexist on the
		// same ServeMux without duplicate route registrations. Lyngua
		// product_inventory block layers on top as tweaks; the dual-mount is
		// activated purely by the presence of that block.
		productInventoryRoutes := centymo.DefaultProductInventoryRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "product_inventory", &productInventoryRoutes)

		productLineRoutes := centymo.DefaultProductLineRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "product_line", &productLineRoutes)

		pricePlanRoutes := centymo.DefaultPricePlanRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "price_plan", &pricePlanRoutes)

		priceScheduleRoutes := centymo.DefaultPriceScheduleRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "price_schedule", &priceScheduleRoutes)

		priceListRoutes := centymo.DefaultPriceListRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "price_list", &priceListRoutes)

		planRoutes := centymo.DefaultPlanRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "plan", &planRoutes)

		// Bundle-mount Plan routes — namespace-shifted onto /app/inventory/bundles/*.
		// Lyngua plan_bundle block can layer additional tweaks on top.
		planBundleRoutes := centymo.DefaultPlanBundleRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "plan_bundle", &planBundleRoutes)

		// Inventory-mount ProductLine routes — namespace-shifted onto /app/inventory/product-lines/*.
		// Lyngua product_line_inventory block can layer additional tweaks on top.
		productLineInventoryRoutes := centymo.DefaultProductLineInventoryRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "product_line_inventory", &productLineInventoryRoutes)

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

		// Inventory-flavoured product labels. Starts from the already-loaded
		// service product labels (centymo.ProductLabels has no exported
		// DefaultProductLabels() factory — the service product.json is the
		// de-facto baseline) and sparse-overlays product_inventory.json so the
		// inventory mount can use distinct headings/buttons (e.g. "Add Product")
		// without duplicating every key in the service product.json.
		productInventoryLabels := productLabels
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "product_inventory.json", "product_inventory", &productInventoryLabels)

		productLineLabels := centymo.DefaultProductLineLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "product_line.json", "product_line", &productLineLabels)

		pricePlanLabels := centymo.DefaultPricePlanLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "price_plan.json", "price_plan", &pricePlanLabels)

		priceScheduleLabels := centymo.DefaultPriceScheduleLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "price_schedule.json", "priceSchedule", &priceScheduleLabels)

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

		resourceRoutes := centymo.DefaultResourceRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "resource", &resourceRoutes)

		resourceLabels := centymo.DefaultResourceLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "resource.json", "resource", &resourceLabels)

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
				// Payment terms dropdown (client/both scope)
				ListPaymentTerms: func(fctx context.Context) ([]*revenuemod.PaymentTermOption, error) {
					rows, err := db.ListSimple(fctx, "payment_term")
					if err != nil {
						return nil, err
					}
					opts := make([]*revenuemod.PaymentTermOption, 0, len(rows))
					for _, row := range rows {
						id, _ := row["id"].(string)
						name, _ := row["name"].(string)
						entityScope, _ := row["entity_scope"].(string)
						if id == "" {
							continue
						}
						if entityScope != "client" && entityScope != "both" {
							continue
						}
						var netDays int32
						switch v := row["net_days"].(type) {
						case int32:
							netDays = v
						case int64:
							netDays = int32(v)
						case float64:
							netDays = int32(v)
						}
						opts = append(opts, &revenuemod.PaymentTermOption{Id: id, Name: name, NetDays: netDays})
					}
					return opts, nil
				},
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
			// Client search for revenue form autocomplete
			if useCases.Entity != nil && useCases.Entity.Client != nil {
				revDeps.ListClients = useCases.Entity.Client.ListClients.Execute
				if useCases.Entity.Client.SearchClientsByName != nil {
					revDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName.Execute
				}
			}
			// Subscription search for revenue form autocomplete
			if useCases.Subscription != nil && useCases.Subscription.Subscription != nil && useCases.Subscription.Subscription.ListSubscriptions != nil {
				revDeps.ListSubscriptions = useCases.Subscription.Subscription.ListSubscriptions.Execute
			}
			// Subscription auto-populate for revenue add (read subscription + price plan + product price plans)
			if useCases.Subscription != nil {
				if useCases.Subscription.Subscription != nil && useCases.Subscription.Subscription.ReadSubscription != nil {
					revDeps.ReadSubscription = useCases.Subscription.Subscription.ReadSubscription.Execute
				}
				if useCases.Subscription.PricePlan != nil && useCases.Subscription.PricePlan.ReadPricePlan != nil {
					revDeps.ReadPricePlan = useCases.Subscription.PricePlan.ReadPricePlan.Execute
				}
				if useCases.Subscription.ProductPricePlan != nil && useCases.Subscription.ProductPricePlan.ListProductPricePlans != nil {
					revDeps.ListProductPricePlans = useCases.Subscription.ProductPricePlan.ListProductPricePlans.Execute
				}
			}
			if useCases.Product != nil && useCases.Product.Product != nil {
				if useCases.Product.Product.ReadProduct != nil {
					revDeps.ReadProduct = useCases.Product.Product.ReadProduct.Execute
				}
				if useCases.Product.Product.ListProducts != nil {
					revDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
				}
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

			// Price lookup for line item (find applicable price list + price product)
			if useCases.Product != nil && useCases.Product.PriceList != nil && useCases.Product.PriceList.FindApplicablePriceList != nil {
				revDeps.FindApplicablePriceList = useCases.Product.PriceList.FindApplicablePriceList.Execute
			}
			if useCases.Product != nil && useCases.Product.PriceProduct != nil && useCases.Product.PriceProduct.ListPriceProducts != nil {
				revDeps.ListPriceProducts = useCases.Product.PriceProduct.ListPriceProducts.Execute
			}

			// Job activity lookup for "from_activities" revenue type
			if useCases.Operation != nil && useCases.Operation.JobActivity != nil && useCases.Operation.JobActivity.ReadJobActivity != nil {
				revDeps.ReadJobActivity = useCases.Operation.JobActivity.ReadJobActivity.Execute
			}

			revenueMod := revenuemod.NewModule(revDeps)
			revenueMod.RegisterRoutes(ctx.Routes)
			// Invoice download is http.HandlerFunc (bypasses view/template layer)
			handleFunc(ctx.Routes, "GET", revenueRoutes.InvoiceDownloadURL, revenueMod.InvoiceDownload)
			// Send email is http.HandlerFunc (bypasses view/template layer)
			handleFunc(ctx.Routes, "POST", revenueRoutes.SendEmailURL, revenueMod.SendEmailHandler)
			handleFunc(ctx.Routes, "GET", revenueRoutes.SearchClientURL, revenueMod.SearchClients)
			handleFunc(ctx.Routes, "GET", revenueRoutes.SearchSubscriptionURL, revenueMod.SearchSubscriptions)
			handleFunc(ctx.Routes, "GET", revenueRoutes.SearchLocationURL, revenueMod.SearchLocations)
			handleFunc(ctx.Routes, "GET", revenueRoutes.SearchProductURL, revenueMod.SearchProducts)
			handleFunc(ctx.Routes, "GET", revenueRoutes.PriceLookupURL, revenueMod.PriceLookup)
		}

		// =====================================================================
		// Product module
		// =====================================================================

		if cfg.wantProduct() {
			var getProductInUseIDs func(context.Context, []string) (map[string]bool, error)
			if refChecker != nil {
				getProductInUseIDs = refChecker.GetProductInUseIDs
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
				Routes:              productRoutes,
				Mode:                "service",
				DB:                  db,
				Labels:              productLabels,
				CommonLabels:        ctx.Common,
				TableLabels:         centymoTableLabels,
				GetInUseIDs:         getProductInUseIDs,
				DefaultProductKind:  defaultProductKind,
				DefaultDeliveryMode: defaultDeliveryMode,
				DefaultTrackingMode: defaultTrackingMode,
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
			if productInventoryRoutes.ListURL != productRoutes.ListURL {
				productInventoryDeps := *productDeps
				productInventoryDeps.Routes = productInventoryRoutes
				productInventoryDeps.Mode = "inventory"
				productInventoryDeps.Labels = productInventoryLabels
				productInventoryDeps.DefaultProductKind = "stocked_good"
				productInventoryDeps.DefaultDeliveryMode = "shipped"
				productInventoryDeps.DefaultTrackingMode = "bulk"
				productmod.NewModule(&productInventoryDeps).RegisterRoutes(ctx.Routes)
			}
		}

		// =====================================================================
		// Product Line module
		// =====================================================================

		if cfg.wantProductLine() {
			if useCases.Product != nil && useCases.Product.Line != nil {
				uc := useCases.Product.Line
				modDeps := &productlinemod.ModuleDeps{
					Routes:       productLineRoutes,
					Labels:       productLineLabels,
					CommonLabels: ctx.Common,
					TableLabels:  centymoTableLabels,
					ListLines:    uc.ListLines.Execute,
					ReadLine:     uc.ReadLine.Execute,
					CreateLine:   uc.CreateLine.Execute,
					UpdateLine:   uc.UpdateLine.Execute,
					DeleteLine:   uc.DeleteLine.Execute,
				}
				if refChecker != nil {
					modDeps.GetInUseIDs = refChecker.GetLineInUseIDs
				}
				productlinemod.NewModule(modDeps).RegisterRoutes(ctx.Routes)

				// Inventory-mount ProductLine second registration on distinct URLs.
				// Gate: if a lyngua product_line_inventory override ever collapses
				// ListURL back onto the services mount, skip to avoid a ServeMux
				// duplicate-route panic.
				if productLineInventoryRoutes.ListURL != productLineRoutes.ListURL {
					productLineInventoryDeps := &productlinemod.ModuleDeps{
						Routes:       productLineInventoryRoutes,
						Labels:       productLineLabels,
						CommonLabels: ctx.Common,
						TableLabels:  centymoTableLabels,
						ListLines:    uc.ListLines.Execute,
						ReadLine:     uc.ReadLine.Execute,
						CreateLine:   uc.CreateLine.Execute,
						UpdateLine:   uc.UpdateLine.Execute,
						DeleteLine:   uc.DeleteLine.Execute,
					}
					if refChecker != nil {
						productLineInventoryDeps.GetInUseIDs = refChecker.GetLineInUseIDs
					}
					productlinemod.NewModule(productLineInventoryDeps).RegisterRoutes(ctx.Routes)
				}
			}
		}

		// =====================================================================
		// Price Plan module (standalone — separate from plan-nested price plans)
		// =====================================================================

		if cfg.wantPricePlan() {
			if useCases.Subscription != nil && useCases.Subscription.PricePlan != nil {
				uc := useCases.Subscription.PricePlan
				var getPricePlanInUseIDs func(context.Context, []string) (map[string]bool, error)
				if refChecker != nil {
					getPricePlanInUseIDs = refChecker.GetPricePlanInUseIDs
				}
				pricePlanDeps := &priceplanmod.ModuleDeps{
					Routes:               pricePlanRoutes,
					Labels:               pricePlanLabels,
					CommonLabels:         ctx.Common,
					TableLabels:          centymoTableLabels,
					ListPricePlans:       uc.ListPricePlans.Execute,
					ReadPricePlan:        uc.ReadPricePlan.Execute,
					CreatePricePlan:      uc.CreatePricePlan.Execute,
					UpdatePricePlan:      uc.UpdatePricePlan.Execute,
					DeletePricePlan:      uc.DeletePricePlan.Execute,
					GetPricePlanInUseIDs: getPricePlanInUseIDs,
				}
				// Price schedule listing — parent container (owns location + date range)
				if useCases.Subscription.PriceSchedule != nil {
					pricePlanDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
				}
				// Add plan listing if available
				if useCases.Subscription != nil && useCases.Subscription.Plan != nil {
					pricePlanDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
				}
				// Add product listing for detail page product selector
				if useCases.Product != nil && useCases.Product.Product != nil {
					pricePlanDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
				}
				// Add product plan listing for scoping product selector to plan's products
				if useCases.Product != nil && useCases.Product.ProductPlan != nil {
					pricePlanDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
				}
				// Add ProductPricePlan CRUD for detail page
				if useCases.Subscription.ProductPricePlan != nil {
					ppp := useCases.Subscription.ProductPricePlan
					pricePlanDeps.ListProductPricePlans = ppp.ListProductPricePlans.Execute
					pricePlanDeps.CreateProductPricePlan = ppp.CreateProductPricePlan.Execute
					pricePlanDeps.UpdateProductPricePlan = ppp.UpdateProductPricePlan.Execute
					pricePlanDeps.DeleteProductPricePlan = ppp.DeleteProductPricePlan.Execute
				}
				priceplanmod.NewModule(pricePlanDeps).RegisterRoutes(ctx.Routes)
			}
		}

		// =====================================================================
		// Price Schedule module
		// =====================================================================

		if cfg.wantPriceSchedule() {
			if useCases.Subscription != nil && useCases.Subscription.PriceSchedule != nil {
				uc := useCases.Subscription.PriceSchedule
				var getPriceScheduleInUseIDs func(context.Context, []string) (map[string]bool, error)
				if refChecker != nil {
					getPriceScheduleInUseIDs = refChecker.GetPriceScheduleInUseIDs
				}
				priceScheduleDeps := &priceschedulemod.ModuleDeps{
					Routes:                   priceScheduleRoutes,
					Labels:                   priceScheduleLabels,
					CommonLabels:             ctx.Common,
					TableLabels:              centymoTableLabels,
					ListPriceSchedules:       uc.ListPriceSchedules.Execute,
					ReadPriceSchedule:        uc.ReadPriceSchedule.Execute,
					CreatePriceSchedule:      uc.CreatePriceSchedule.Execute,
					UpdatePriceSchedule:      uc.UpdatePriceSchedule.Execute,
					DeletePriceSchedule:      uc.DeletePriceSchedule.Execute,
					GetPriceScheduleInUseIDs: getPriceScheduleInUseIDs,
				}
				// Add location listing if available
				if useCases.Entity != nil && useCases.Entity.Location != nil {
					priceScheduleDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
				}
				// Plans tab on the detail page lists price_plans filtered by price_schedule_id FK.
				if useCases.Subscription.PricePlan != nil {
					priceScheduleDeps.ListPricePlans = useCases.Subscription.PricePlan.ListPricePlans.Execute
					priceScheduleDeps.CreatePricePlan = useCases.Subscription.PricePlan.CreatePricePlan.Execute
					priceScheduleDeps.ReadPricePlan = useCases.Subscription.PricePlan.ReadPricePlan.Execute
					priceScheduleDeps.UpdatePricePlan = useCases.Subscription.PricePlan.UpdatePricePlan.Execute
					priceScheduleDeps.DeletePricePlan = useCases.Subscription.PricePlan.DeletePricePlan.Execute
				}
				// Reference checker for in-use guard (disables row Delete + locks pricing fields
				// on the edit drawer when a price_plan is referenced by active subscriptions).
				if refChecker != nil {
					priceScheduleDeps.GetPricePlanInUseIDs = refChecker.GetPricePlanInUseIDs
				}
				if useCases.Subscription.Plan != nil {
					priceScheduleDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
				}
				// Schedule-scoped plan detail (info + product-prices tabs) needs product lookups + ProductPricePlan CRUD
				if useCases.Product != nil && useCases.Product.Product != nil {
					priceScheduleDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
				}
				if useCases.Product != nil && useCases.Product.ProductPlan != nil {
					priceScheduleDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
				}
				if useCases.Subscription.ProductPricePlan != nil {
					ppp := useCases.Subscription.ProductPricePlan
					priceScheduleDeps.ListProductPricePlans = ppp.ListProductPricePlans.Execute
					priceScheduleDeps.CreateProductPricePlan = ppp.CreateProductPricePlan.Execute
					priceScheduleDeps.UpdateProductPricePlan = ppp.UpdateProductPricePlan.Execute
					priceScheduleDeps.DeleteProductPricePlan = ppp.DeleteProductPricePlan.Execute
				}
				priceschedulemod.NewModule(priceScheduleDeps).RegisterRoutes(ctx.Routes)
			}
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
			if refChecker != nil {
				planListDeps.GetInUseIDs = refChecker.GetPlanInUseIDs
			}
			ctx.Routes.GET(planRoutes.ListURL, planlist.NewView(planListDeps))
			ctx.Routes.GET(planRoutes.TableURL, planlist.NewTableView(planListDeps))

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
				if useCases.Entity != nil && useCases.Entity.Location != nil {
					planDetailDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
				}
				if useCases.Subscription.PriceSchedule != nil {
					planDetailDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
				}
				ctx.Routes.GET(planRoutes.DetailURL, plandetail.NewView(planDetailDeps))
				ctx.Routes.GET(planRoutes.TabActionURL, plandetail.NewTabAction(planDetailDeps))
				// Plan attachments
				if uploadFile != nil {
					ctx.Routes.GET(planRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planDetailDeps))
					ctx.Routes.POST(planRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planDetailDeps))
					ctx.Routes.POST(planRoutes.AttachmentDeleteURL, plandetail.NewAttachmentDeleteAction(planDetailDeps))
				}
				// PricePlan CRUD within plan detail
				if useCases.Subscription.PricePlan != nil && useCases.Subscription.PricePlan.CreatePricePlan != nil {
					ppActionDeps := &planaction.PricePlanDeps{
						Routes:          planRoutes,
						Labels:          planLabels,
						CommonLabels:    ctx.Common,
						CreatePricePlan: useCases.Subscription.PricePlan.CreatePricePlan.Execute,
						ReadPricePlan:   useCases.Subscription.PricePlan.ReadPricePlan.Execute,
						UpdatePricePlan: useCases.Subscription.PricePlan.UpdatePricePlan.Execute,
						DeletePricePlan: useCases.Subscription.PricePlan.DeletePricePlan.Execute,
					}
					if useCases.Subscription.PriceSchedule != nil {
						ppActionDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
					}
					if useCases.Product != nil && useCases.Product.ProductPlan != nil {
						ppActionDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
					}
					if useCases.Subscription.ProductPricePlan != nil {
						ppActionDeps.CreateProductPricePlan = useCases.Subscription.ProductPricePlan.CreateProductPricePlan.Execute
						ppActionDeps.ListProductPricePlans = useCases.Subscription.ProductPricePlan.ListProductPricePlans.Execute
					}
					ctx.Routes.GET(planRoutes.PricePlanAddURL, planaction.NewPricePlanAddAction(ppActionDeps))
					ctx.Routes.POST(planRoutes.PricePlanAddURL, planaction.NewPricePlanAddAction(ppActionDeps))
					ctx.Routes.GET(planRoutes.PricePlanEditURL, planaction.NewPricePlanEditAction(ppActionDeps))
					ctx.Routes.POST(planRoutes.PricePlanEditURL, planaction.NewPricePlanEditAction(ppActionDeps))
					ctx.Routes.POST(planRoutes.PricePlanDeleteURL, planaction.NewPricePlanDeleteAction(ppActionDeps))
				}
				// ProductPlan CRUD within plan detail
				if useCases.Product != nil && useCases.Product.ProductPlan != nil && useCases.Product.ProductPlan.CreateProductPlan != nil {
					productPlanActionDeps := &planaction.ProductPlanDeps{
						Routes:            planRoutes,
						Labels:            planLabels,
						CreateProductPlan: useCases.Product.ProductPlan.CreateProductPlan.Execute,
						ReadProductPlan:   useCases.Product.ProductPlan.ReadProductPlan.Execute,
						UpdateProductPlan: useCases.Product.ProductPlan.UpdateProductPlan.Execute,
						DeleteProductPlan: useCases.Product.ProductPlan.DeleteProductPlan.Execute,
					}
					if useCases.Product.Product != nil {
						productPlanActionDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
					}
					if useCases.Product.ProductPlan.ListProductPlans != nil {
						productPlanActionDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
					}
					ctx.Routes.GET(planRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(productPlanActionDeps))
					ctx.Routes.POST(planRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(productPlanActionDeps))
					ctx.Routes.GET(planRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(productPlanActionDeps))
					ctx.Routes.POST(planRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(productPlanActionDeps))
					ctx.Routes.POST(planRoutes.ProductPlanDeleteURL, planaction.NewProductPlanDeleteAction(productPlanActionDeps))
				}
			}

			// =====================================================================
			// Plan bundle inventory-mount (second registration on distinct URLs)
			// =====================================================================
			// Reuses the same plan views but on /app/inventory/bundles/* URLs.
			// Gate: if a lyngua plan_bundle override ever collapses ListURL back
			// onto the services mount, skip to avoid a ServeMux duplicate-route panic.
			if cfg.wantPlan() && planBundleRoutes.ListURL != planRoutes.ListURL {
				planBundleListDeps := &planlist.ListViewDeps{
					Routes:       planBundleRoutes,
					Labels:       planLabels,
					CommonLabels: ctx.Common,
					TableLabels:  centymoTableLabels,
				}
				if useCases.Subscription != nil && useCases.Subscription.Plan != nil {
					planBundleListDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
				}
				if refChecker != nil {
					planBundleListDeps.GetInUseIDs = refChecker.GetPlanInUseIDs
				}
				ctx.Routes.GET(planBundleRoutes.ListURL, planlist.NewView(planBundleListDeps))
				ctx.Routes.GET(planBundleRoutes.TableURL, planlist.NewTableView(planBundleListDeps))

				if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.CreatePlan != nil {
					planBundleActionDeps := &planaction.Deps{
						Routes:     planBundleRoutes,
						Labels:     planLabels,
						CreatePlan: useCases.Subscription.Plan.CreatePlan.Execute,
						ReadPlan:   useCases.Subscription.Plan.ReadPlan.Execute,
						UpdatePlan: useCases.Subscription.Plan.UpdatePlan.Execute,
						DeletePlan: useCases.Subscription.Plan.DeletePlan.Execute,
					}
					ctx.Routes.GET(planBundleRoutes.AddURL, planaction.NewAddAction(planBundleActionDeps))
					ctx.Routes.POST(planBundleRoutes.AddURL, planaction.NewAddAction(planBundleActionDeps))
					ctx.Routes.GET(planBundleRoutes.EditURL, planaction.NewEditAction(planBundleActionDeps))
					ctx.Routes.POST(planBundleRoutes.EditURL, planaction.NewEditAction(planBundleActionDeps))
					ctx.Routes.POST(planBundleRoutes.DeleteURL, planaction.NewDeleteAction(planBundleActionDeps))
				}

				if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
					planBundleDetailDeps := &plandetail.DetailViewDeps{
						Routes:       planBundleRoutes,
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
						planBundleDetailDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
					}
					if useCases.Subscription.PricePlan != nil {
						planBundleDetailDeps.ListPricePlans = useCases.Subscription.PricePlan.ListPricePlans.Execute
					}
					if useCases.Entity != nil && useCases.Entity.Location != nil {
						planBundleDetailDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
					}
					if useCases.Subscription.PriceSchedule != nil {
						planBundleDetailDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
					}
					ctx.Routes.GET(planBundleRoutes.DetailURL, plandetail.NewView(planBundleDetailDeps))
					ctx.Routes.GET(planBundleRoutes.TabActionURL, plandetail.NewTabAction(planBundleDetailDeps))
					if uploadFile != nil {
						ctx.Routes.GET(planBundleRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planBundleDetailDeps))
						ctx.Routes.POST(planBundleRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planBundleDetailDeps))
						ctx.Routes.POST(planBundleRoutes.AttachmentDeleteURL, plandetail.NewAttachmentDeleteAction(planBundleDetailDeps))
					}
					if useCases.Subscription.PricePlan != nil && useCases.Subscription.PricePlan.CreatePricePlan != nil {
						ppBundleDeps := &planaction.PricePlanDeps{
							Routes:          planBundleRoutes,
							Labels:          planLabels,
							CommonLabels:    ctx.Common,
							CreatePricePlan: useCases.Subscription.PricePlan.CreatePricePlan.Execute,
							ReadPricePlan:   useCases.Subscription.PricePlan.ReadPricePlan.Execute,
							UpdatePricePlan: useCases.Subscription.PricePlan.UpdatePricePlan.Execute,
							DeletePricePlan: useCases.Subscription.PricePlan.DeletePricePlan.Execute,
						}
						if useCases.Subscription.PriceSchedule != nil {
							ppBundleDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
						}
						if useCases.Product != nil && useCases.Product.ProductPlan != nil {
							ppBundleDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
						}
						if useCases.Subscription.ProductPricePlan != nil {
							ppBundleDeps.CreateProductPricePlan = useCases.Subscription.ProductPricePlan.CreateProductPricePlan.Execute
							ppBundleDeps.ListProductPricePlans = useCases.Subscription.ProductPricePlan.ListProductPricePlans.Execute
						}
						ctx.Routes.GET(planBundleRoutes.PricePlanAddURL, planaction.NewPricePlanAddAction(ppBundleDeps))
						ctx.Routes.POST(planBundleRoutes.PricePlanAddURL, planaction.NewPricePlanAddAction(ppBundleDeps))
						ctx.Routes.GET(planBundleRoutes.PricePlanEditURL, planaction.NewPricePlanEditAction(ppBundleDeps))
						ctx.Routes.POST(planBundleRoutes.PricePlanEditURL, planaction.NewPricePlanEditAction(ppBundleDeps))
						ctx.Routes.POST(planBundleRoutes.PricePlanDeleteURL, planaction.NewPricePlanDeleteAction(ppBundleDeps))
					}
					if useCases.Product != nil && useCases.Product.ProductPlan != nil && useCases.Product.ProductPlan.CreateProductPlan != nil {
						ppBundleProductPlanDeps := &planaction.ProductPlanDeps{
							Routes:            planBundleRoutes,
							Labels:            planLabels,
							CreateProductPlan: useCases.Product.ProductPlan.CreateProductPlan.Execute,
							ReadProductPlan:   useCases.Product.ProductPlan.ReadProductPlan.Execute,
							UpdateProductPlan: useCases.Product.ProductPlan.UpdateProductPlan.Execute,
							DeleteProductPlan: useCases.Product.ProductPlan.DeleteProductPlan.Execute,
						}
						if useCases.Product.Product != nil {
							ppBundleProductPlanDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
						}
						if useCases.Product.ProductPlan.ListProductPlans != nil {
							ppBundleProductPlanDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
						}
						ctx.Routes.GET(planBundleRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(ppBundleProductPlanDeps))
						ctx.Routes.POST(planBundleRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(ppBundleProductPlanDeps))
						ctx.Routes.GET(planBundleRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(ppBundleProductPlanDeps))
						ctx.Routes.POST(planBundleRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(ppBundleProductPlanDeps))
						ctx.Routes.POST(planBundleRoutes.ProductPlanDeleteURL, planaction.NewProductPlanDeleteAction(ppBundleProductPlanDeps))
					}
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
				if useCases.Subscription.PricePlan != nil {
					subActionDeps.ListPricePlans = useCases.Subscription.PricePlan.ListPricePlans.Execute
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
				SqlDB:          ctx.SqlDB,
				Labels:         expenditureLabels,
				TemplateLabels: templateview.DefaultLabels(),
				CommonLabels:   ctx.Common,
				TableLabels:    centymoTableLabels,
				// Payment terms dropdown (supplier/both scope)
				ListPaymentTerms: func(fctx context.Context) ([]*expendituremod.PaymentTermOption, error) {
					rows, err := db.ListSimple(fctx, "payment_term")
					if err != nil {
						return nil, err
					}
					opts := make([]*expendituremod.PaymentTermOption, 0, len(rows))
					for _, row := range rows {
						id, _ := row["id"].(string)
						name, _ := row["name"].(string)
						entityScope, _ := row["entity_scope"].(string)
						if id == "" {
							continue
						}
						if entityScope != "supplier" && entityScope != "both" {
							continue
						}
						var netDays int32
						switch v := row["net_days"].(type) {
						case int32:
							netDays = v
						case int64:
							netDays = int32(v)
						case float64:
							netDays = int32(v)
						}
						opts = append(opts, &expendituremod.PaymentTermOption{Id: id, Name: name, NetDays: netDays})
					}
					return opts, nil
				},
				// Document template CRUD
				ListDocumentTemplates:  listDocTemplates,
				CreateDocumentTemplate: createDocTemplate,
				UpdateDocumentTemplate: updateDocTemplate,
				DeleteDocumentTemplate: deleteDocTemplate,
				UploadFile:             uploadTemplate,
			}
			if useCases.Expenditure != nil && useCases.Expenditure.Expenditure != nil {
				uc := useCases.Expenditure.Expenditure
				expDeps.ListExpenditures = uc.ListExpenditures.Execute
				expDeps.CreateExpenditure = uc.CreateExpenditure.Execute
				expDeps.ReadExpenditure = uc.ReadExpenditure.Execute
				expDeps.UpdateExpenditure = uc.UpdateExpenditure.Execute
				expDeps.DeleteExpenditure = uc.DeleteExpenditure.Execute
			}
			if useCases.Expenditure != nil && useCases.Expenditure.ExpenditureCategory != nil {
				uc := useCases.Expenditure.ExpenditureCategory
				expDeps.ListExpenditureCategories = uc.ListExpenditureCategories.Execute
				expDeps.CreateExpenditureCategory = uc.CreateExpenditureCategory.Execute
				expDeps.ReadExpenditureCategory = uc.ReadExpenditureCategory.Execute
				expDeps.UpdateExpenditureCategory = uc.UpdateExpenditureCategory.Execute
				expDeps.DeleteExpenditureCategory = uc.DeleteExpenditureCategory.Execute
			}
			if useCases.Expenditure != nil && useCases.Expenditure.ExpenditureLineItem != nil {
				uc := useCases.Expenditure.ExpenditureLineItem
				expDeps.CreateExpenditureLineItem = uc.CreateExpenditureLineItem.Execute
				expDeps.ReadExpenditureLineItem = uc.ReadExpenditureLineItem.Execute
				expDeps.UpdateExpenditureLineItem = uc.UpdateExpenditureLineItem.Execute
				expDeps.DeleteExpenditureLineItem = uc.DeleteExpenditureLineItem.Execute
				expDeps.ListExpenditureLineItems = uc.ListExpenditureLineItems.Execute
			}
			if useCases.Entity != nil && useCases.Entity.Supplier != nil {
				expDeps.ListSuppliers = useCases.Entity.Supplier.ListSuppliers.Execute
			}
			if useCases.Treasury != nil && useCases.Treasury.Disbursement != nil {
				expDeps.DisbursementRoutes = disbursementRoutes
				expDeps.DisbursementLabels = disbursementLabels
				expDeps.CreateDisbursement = useCases.Treasury.Disbursement.CreateDisbursement.Execute
			}
			expendituremod.NewModule(expDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Resource module
		// =====================================================================

		if cfg.wantResource() {
			resourceDeps := &resourcemod.ModuleDeps{
				Routes:       resourceRoutes,
				Labels:       resourceLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
			}
			if useCases.Product != nil && useCases.Product.Resource != nil {
				uc := useCases.Product.Resource
				if uc.ListResources != nil {
					resourceDeps.ListResources = uc.ListResources.Execute
				}
				if uc.ReadResource != nil {
					resourceDeps.ReadResource = uc.ReadResource.Execute
				}
				if uc.CreateResource != nil {
					resourceDeps.CreateResource = uc.CreateResource.Execute
				}
				if uc.UpdateResource != nil {
					resourceDeps.UpdateResource = uc.UpdateResource.Execute
				}
				if uc.DeleteResource != nil {
					resourceDeps.DeleteResource = uc.DeleteResource.Execute
				}
			}
			resourcemod.NewModule(resourceDeps).RegisterRoutes(ctx.Routes)
		}

		log.Println("  centymo commerce domain initialized")
		return nil
	}
}
