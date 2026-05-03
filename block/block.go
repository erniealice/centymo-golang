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
	"github.com/erniealice/espyna-golang/reference"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/template"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"

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
	priceschedulepricepldetail "github.com/erniealice/centymo-golang/views/price_schedule/detail/plan"
	procurementmod "github.com/erniealice/centymo-golang/views/procurement"
	procurementrequestmod "github.com/erniealice/centymo-golang/views/procurement_request"
	procurementrequestlinemod "github.com/erniealice/centymo-golang/views/procurement_request_line"
	resourcemod "github.com/erniealice/centymo-golang/views/resource"
	revenuemod "github.com/erniealice/centymo-golang/views/revenue"
	subscriptionaction "github.com/erniealice/centymo-golang/views/subscription/action"
	subscriptiondetail "github.com/erniealice/centymo-golang/views/subscription/detail"
	subscriptionlist "github.com/erniealice/centymo-golang/views/subscription/list"
	suppliercontractmod "github.com/erniealice/centymo-golang/views/supplier_contract"
	suppliercontractlinemod "github.com/erniealice/centymo-golang/views/supplier_contract_line"
	// SPS Wave 4 — six new view modules (price-schedule master + line, recognition + line, accrued + settlement).
	accruedexpensemod "github.com/erniealice/centymo-golang/views/accrued_expense"
	accruedexpensesettlementmod "github.com/erniealice/centymo-golang/views/accrued_expense_settlement"
	expenserecognitionmod "github.com/erniealice/centymo-golang/views/expense_recognition"
	expenserecognitionlinemod "github.com/erniealice/centymo-golang/views/expense_recognition_line"
	suppliercontractpriceschedulemod "github.com/erniealice/centymo-golang/views/supplier_contract_price_schedule"
	suppliercontractpricescheduleinemod "github.com/erniealice/centymo-golang/views/supplier_contract_price_schedule_line"

	procurementrequestpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	// SPS Wave 4 — proto packages for the six new view modules.
	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"
	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
	scpspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule"
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
	// 20260427-supplier-commitments P3a/P3b — five new modules wired by Block.
	supplierContract       bool
	supplierContractLine   bool
	procurementRequest     bool
	procurementRequestLine bool
	procurement            bool
	// SPS Wave 4 — six new modules wired by Block.
	supplierContractPriceSchedule     bool
	supplierContractPriceScheduleLine bool
	expenseRecognition                bool
	expenseRecognitionLine            bool
	accruedExpense                    bool
	accruedExpenseSettlement          bool
	// clientDetailURL is the absolute path template (e.g.
	// "/app/clients/detail/{id}") used for the subscription detail's
	// page-header breadcrumb when accessed via the under-client nested route.
	// Centymo cannot import entydad (wrong dep direction); the consumer
	// supplies it via WithClientDetailURL.
	clientDetailURL string
	// jobDetailURL is the absolute path template (e.g. "/app/jobs/detail/{id}")
	// used by the subscription detail's Operations tab to deep-link to fayna
	// Job detail. Centymo cannot import fayna; the consumer supplies it via
	// WithJobDetailURL. Optional — Operations tab renders job rows without a
	// link when unset.
	jobDetailURL string
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

// 20260427-supplier-commitments — five new module toggles.
func WithSupplierContract() BlockOption       { return func(c *blockConfig) { c.supplierContract = true } }
func WithSupplierContractLine() BlockOption   { return func(c *blockConfig) { c.supplierContractLine = true } }
func WithProcurementRequest() BlockOption     { return func(c *blockConfig) { c.procurementRequest = true } }
func WithProcurementRequestLine() BlockOption { return func(c *blockConfig) { c.procurementRequestLine = true } }
func WithProcurement() BlockOption            { return func(c *blockConfig) { c.procurement = true } }

// SPS Wave 4 — six new module toggles (supplier-side pricing graph + accrual layer).
func WithSupplierContractPriceSchedule() BlockOption {
	return func(c *blockConfig) { c.supplierContractPriceSchedule = true }
}
func WithSupplierContractPriceScheduleLine() BlockOption {
	return func(c *blockConfig) { c.supplierContractPriceScheduleLine = true }
}
func WithExpenseRecognition() BlockOption {
	return func(c *blockConfig) { c.expenseRecognition = true }
}
func WithExpenseRecognitionLine() BlockOption {
	return func(c *blockConfig) { c.expenseRecognitionLine = true }
}
func WithAccruedExpense() BlockOption {
	return func(c *blockConfig) { c.accruedExpense = true }
}
func WithAccruedExpenseSettlement() BlockOption {
	return func(c *blockConfig) { c.accruedExpenseSettlement = true }
}

// WithClientDetailURL supplies the entydad client-detail path template (e.g.
// "/app/clients/detail/{id}") so the subscription detail page can render a
// "client → subscription" breadcrumb when accessed under a client context.
// Optional — when unset the breadcrumb label still renders (sourced from the
// joined client) but isn't a link.
func WithClientDetailURL(url string) BlockOption {
	return func(c *blockConfig) { c.clientDetailURL = url }
}

// WithJobDetailURL supplies the fayna job-detail path template (e.g.
// "/app/jobs/detail/{id}") so the subscription detail's Operations tab can
// render a deep link to each spawned Job. Optional — when unset rows render
// without a link.
// 2026-04-29 auto-spawn-jobs-from-subscription Phase D.
func WithJobDetailURL(url string) BlockOption {
	return func(c *blockConfig) { c.jobDetailURL = url }
}

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

// 20260427-supplier-commitments — five new module toggles.
func (c *blockConfig) wantSupplierContract() bool       { return c.enableAll || c.supplierContract }
func (c *blockConfig) wantSupplierContractLine() bool   { return c.enableAll || c.supplierContractLine }
func (c *blockConfig) wantProcurementRequest() bool     { return c.enableAll || c.procurementRequest }
func (c *blockConfig) wantProcurementRequestLine() bool { return c.enableAll || c.procurementRequestLine }
func (c *blockConfig) wantProcurement() bool            { return c.enableAll || c.procurement }

// SPS Wave 4 — six new module want() helpers.
func (c *blockConfig) wantSupplierContractPriceSchedule() bool {
	return c.enableAll || c.supplierContractPriceSchedule
}
func (c *blockConfig) wantSupplierContractPriceScheduleLine() bool {
	return c.enableAll || c.supplierContractPriceScheduleLine
}
func (c *blockConfig) wantExpenseRecognition() bool {
	return c.enableAll || c.expenseRecognition
}
func (c *blockConfig) wantExpenseRecognitionLine() bool {
	return c.enableAll || c.expenseRecognitionLine
}
func (c *blockConfig) wantAccruedExpense() bool {
	return c.enableAll || c.accruedExpense
}
func (c *blockConfig) wantAccruedExpenseSettlement() bool {
	return c.enableAll || c.accruedExpenseSettlement
}

// ---------------------------------------------------------------------------
// Block — the main Lego entry point
// ---------------------------------------------------------------------------

// Block registers centymo domain modules (commerce: inventory, revenue, product,
// product line, pricelist, plan, subscription, collection, disbursement, expenditure, and inline
// report routes). Call with no options to register ALL modules. Call with specific
// WithX() options for a subset.
func Block(opts ...BlockOption) pyeza.AppOption {
	cfg := &blockConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	// "Enable all modules" is a derived flag — true when no module-toggling
	// option was passed. Non-module options (e.g. WithClientDetailURL) must
	// NOT flip this off, otherwise passing only a config option (no module
	// toggle) silently disables every module.
	moduleSelected := cfg.inventory || cfg.revenue || cfg.product || cfg.productLine ||
		cfg.pricePlan || cfg.priceSchedule || cfg.priceList || cfg.plan ||
		cfg.subscription || cfg.collection || cfg.disbursement || cfg.expenditure ||
		cfg.resource ||
		cfg.supplierContract || cfg.supplierContractLine ||
		cfg.procurementRequest || cfg.procurementRequestLine || cfg.procurement ||
		cfg.supplierContractPriceSchedule || cfg.supplierContractPriceScheduleLine ||
		cfg.expenseRecognition || cfg.expenseRecognitionLine ||
		cfg.accruedExpense || cfg.accruedExpenseSettlement
	cfg.enableAll = !moduleSelected

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
		var refChecker reference.Checker
		if ctx.RefChecker != nil {
			refChecker, _ = ctx.RefChecker.(reference.Checker)
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

		// Supplies mount — third Product module registration scoped to
		// product_kind = 'consumable' (used-in-service-delivery items: gauze,
		// lotion, cleaning solution, coffee beans). Sibling to product_inventory
		// (resold goods) under the Inventory accordion. Lyngua
		// `product_supplies` can override individual URLs on top of the
		// /app/inventory/supplies/* namespace.
		productSuppliesRoutes := centymo.DefaultProductSuppliesRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "product_supplies", &productSuppliesRoutes)

		productLineRoutes := centymo.DefaultProductLineRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "product_line", &productLineRoutes)

		pricePlanRoutes := centymo.DefaultPricePlanRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "price_plan", &pricePlanRoutes)

		priceScheduleRoutes := centymo.DefaultPriceScheduleRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "price_schedule", &priceScheduleRoutes)

		// Inventory-mount variant — namespace-shifted onto /app/inventory/price-schedules/*.
		// Anchors ActiveNav to "inventory" so the inventory accordion sidebar stays open
		// when users browse price schedules from the inventory side. A lyngua
		// `price_schedule_inventory` override can layer additional tweaks on top.
		priceScheduleInventoryRoutes := centymo.DefaultPriceScheduleInventoryRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "price_schedule_inventory", &priceScheduleInventoryRoutes)

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

		// Supplies-flavoured labels. Same overlay pattern: baseline from the
		// service product labels, sparse-overlay product_supplies.json for
		// headings / CTA ("Add Supply") that should differ from both the
		// services mount and the inventory mount.
		productSuppliesLabels := productLabels
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "product_supplies.json", "product_supplies", &productSuppliesLabels)

		productLineLabels := centymo.DefaultProductLineLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "product_line.json", "product_line", &productLineLabels)

		pricePlanLabels := centymo.DefaultPricePlanLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "price_plan.json", "price_plan", &pricePlanLabels)

		productPricePlanLabels := centymo.DefaultProductPricePlanLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "product_price_plan.json", "product_price_plan", &productPricePlanLabels)

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

		// 20260427-supplier-commitments — load routes + labels for the five new view modules.
		supplierContractRoutes := centymo.DefaultSupplierContractRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "supplier_contract", &supplierContractRoutes)
		supplierContractLabels := centymo.DefaultSupplierContractLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "supplier_contract.json", "supplierContract", &supplierContractLabels)

		procurementRequestRoutes := centymo.DefaultProcurementRequestRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "procurement_request", &procurementRequestRoutes)
		procurementRequestLabels := centymo.DefaultProcurementRequestLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "procurement_request.json", "procurementRequest", &procurementRequestLabels)

		procurementRoutes := centymo.DefaultProcurementRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "procurement", &procurementRoutes)
		// Procurement Operations composition app — no Default*Labels factory
		// yet (P4 lyngua wiring deferred); zero-value struct is fine until
		// translations land. LoadPathIfExists is a no-op if the file is absent.
		var procurementLabels centymo.ProcurementLabels
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "procurement.json", "procurement", &procurementLabels)

		// SPS Wave 4 — Routes + Labels for the six new view modules.
		supplierContractPriceScheduleRoutes := centymo.DefaultSupplierContractPriceScheduleRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "supplier_contract_price_schedule", &supplierContractPriceScheduleRoutes)
		supplierContractPriceScheduleLabels := centymo.DefaultSupplierContractPriceScheduleLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "supplier_contract_price_schedule.json", "supplierContractPriceSchedule", &supplierContractPriceScheduleLabels)

		expenseRecognitionRoutes := centymo.DefaultExpenseRecognitionRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "expense_recognition", &expenseRecognitionRoutes)
		expenseRecognitionLabels := centymo.DefaultExpenseRecognitionLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "expense_recognition.json", "expenseRecognition", &expenseRecognitionLabels)

		accruedExpenseRoutes := centymo.DefaultAccruedExpenseRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "accrued_expense", &accruedExpenseRoutes)
		accruedExpenseLabels := centymo.DefaultAccruedExpenseLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "accrued_expense.json", "accruedExpense", &accruedExpenseLabels)

		// =====================================================================
		// Inventory module
		// =====================================================================

		if cfg.wantInventory() {
			invDeps := &inventorymod.ModuleDeps{
				Routes:       inventoryRoutes,
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
					if uc.GetInventoryMovementsListPageData != nil {
						invDeps.GetInventoryMovementsListPageData = uc.GetInventoryMovementsListPageData.Execute
					}
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

			// Recognize-revenue use case — shared with the subscription
			// recognize-drawer flow. When the manual revenue-add picks a
			// subscription, autoPopulateLineItems delegates to this same use
			// case (skip_header=true mode) so both paths converge.
			if useCases.Revenue != nil && useCases.Revenue.Revenue != nil &&
				useCases.Revenue.Revenue.RecognizeRevenueFromSubscription != nil {
				revDeps.RecognizeRevenueFromSubscription =
					useCases.Revenue.Revenue.RecognizeRevenueFromSubscription.Execute
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
			if productInventoryRoutes.ListURL != productRoutes.ListURL {
				productInventoryDeps := *productDeps
				productInventoryDeps.Routes = productInventoryRoutes
				productInventoryDeps.Mode = "inventory"
				productInventoryDeps.Labels = productInventoryLabels
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
			if productSuppliesRoutes.ListURL != productRoutes.ListURL &&
				productSuppliesRoutes.ListURL != productInventoryRoutes.ListURL {
				productSuppliesDeps := *productDeps
				productSuppliesDeps.Routes = productSuppliesRoutes
				productSuppliesDeps.Mode = "supplies"
				productSuppliesDeps.Labels = productSuppliesLabels
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
				// 2026-04-27 plan-client-scope plan §6.7 — closure used to look
				// up the parent PriceSchedule's client name for the info banner
				// rendered on the price-plan drawer.
				var ppListClientNames func(ctx context.Context) map[string]string
				if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
					lc := useCases.Entity.Client.ListClients.Execute
					ppListClientNames = func(fctx context.Context) map[string]string {
						out := map[string]string{}
						resp, err := lc(fctx, &clientpb.ListClientsRequest{})
						if err != nil {
							return out
						}
						for _, c := range resp.GetData() {
							label := c.GetName()
							if label == "" {
								if u := c.GetUser(); u != nil {
									label = u.GetFirstName() + " " + u.GetLastName()
								}
							}
							out[c.GetId()] = label
						}
						return out
					}
				}

				pricePlanDeps := &priceplanmod.ModuleDeps{
					Routes:                    pricePlanRoutes,
					Labels:                    pricePlanLabels,
					ProductPricePlanLabels:    productPricePlanLabels,
					PriceScheduleDetailLabels: priceScheduleLabels.Detail,
					CommonLabels:              ctx.Common,
					TableLabels:               centymoTableLabels,
					ListPricePlans:         uc.ListPricePlans.Execute,
					ReadPricePlan:          uc.ReadPricePlan.Execute,
					CreatePricePlan:        uc.CreatePricePlan.Execute,
					UpdatePricePlan:        uc.UpdatePricePlan.Execute,
					DeletePricePlan:        uc.DeletePricePlan.Execute,
					GetPricePlanInUseIDs:   getPricePlanInUseIDs,
					ListClientNames:        ppListClientNames,
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
				// 2026-04-29 milestone-billing plan §5 / Phase D — milestone phase
				// select on the PPP drawer needs ReadPlan (to resolve job_template_id)
				// and ListByJobTemplate (to load phase rows).
				if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
					pricePlanDeps.ReadPlan = useCases.Subscription.Plan.ReadPlan.Execute
				}
				if useCases.Operation != nil && useCases.Operation.JobTemplatePhase != nil && useCases.Operation.JobTemplatePhase.ListByJobTemplate != nil {
					pricePlanDeps.ListJobTemplatePhasesByJobTemplate = useCases.Operation.JobTemplatePhase.ListByJobTemplate.Execute
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
				// 2026-04-27 plan-client-scope plan §6.1 / §4.4.1 — schedule list
				// Client column lookup + drawer Client picker. Same listClientNames
				// helper used by the plan list.
				var psListClientNames func(ctx context.Context) map[string]string
				if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
					lc := useCases.Entity.Client.ListClients.Execute
					psListClientNames = func(fctx context.Context) map[string]string {
						out := map[string]string{}
						resp, err := lc(fctx, &clientpb.ListClientsRequest{})
						if err != nil {
							return out
						}
						for _, c := range resp.GetData() {
							label := c.GetName()
							if label == "" {
								if u := c.GetUser(); u != nil {
									label = u.GetFirstName() + " " + u.GetLastName()
								}
							}
							out[c.GetId()] = label
						}
						return out
					}
				}

				priceScheduleDeps := &priceschedulemod.ModuleDeps{
					Routes:                   priceScheduleRoutes,
					Labels:                   priceScheduleLabels,
					PricePlanLabels:          pricePlanLabels,
					ProductPricePlanLabels:   productPricePlanLabels,
					CommonLabels:             ctx.Common,
					TableLabels:              centymoTableLabels,
					ListPriceSchedules:       uc.ListPriceSchedules.Execute,
					ReadPriceSchedule:        uc.ReadPriceSchedule.Execute,
					CreatePriceSchedule:      uc.CreatePriceSchedule.Execute,
					UpdatePriceSchedule:      uc.UpdatePriceSchedule.Execute,
					DeletePriceSchedule:      uc.DeletePriceSchedule.Execute,
					GetPriceScheduleInUseIDs: getPriceScheduleInUseIDs,
					ListClientNames:          psListClientNames,
				}
				// 2026-04-27 plan-client-scope plan §6.7 / §4.4.1 — Client picker
				// for the schedule add/edit drawer.
				if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
					priceScheduleDeps.ListClients = useCases.Entity.Client.ListClients.Execute
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

				// =====================================================================
				// PriceSchedule inventory-mount (second registration on distinct URLs)
				// =====================================================================
				// Reuses the same PriceSchedule views but on /app/inventory/price-schedules/*.
				// Gate: if a lyngua price_schedule_inventory override ever collapses ListURL
				// back onto the services mount, skip to avoid a ServeMux duplicate-route panic.
				if priceScheduleInventoryRoutes.ListURL != priceScheduleRoutes.ListURL {
					priceScheduleInventoryDeps := &priceschedulemod.ModuleDeps{
						Routes:                   priceScheduleInventoryRoutes,
						Labels:                   priceScheduleLabels,
						PricePlanLabels:          pricePlanLabels,
						ProductPricePlanLabels:   productPricePlanLabels,
						CommonLabels:             ctx.Common,
						TableLabels:              centymoTableLabels,
						ListPriceSchedules:       uc.ListPriceSchedules.Execute,
						ReadPriceSchedule:        uc.ReadPriceSchedule.Execute,
						CreatePriceSchedule:      uc.CreatePriceSchedule.Execute,
						UpdatePriceSchedule:      uc.UpdatePriceSchedule.Execute,
						DeletePriceSchedule:      uc.DeletePriceSchedule.Execute,
						GetPriceScheduleInUseIDs: getPriceScheduleInUseIDs,
						ListLocations:            priceScheduleDeps.ListLocations,
						ListPricePlans:           priceScheduleDeps.ListPricePlans,
						CreatePricePlan:          priceScheduleDeps.CreatePricePlan,
						ReadPricePlan:            priceScheduleDeps.ReadPricePlan,
						UpdatePricePlan:          priceScheduleDeps.UpdatePricePlan,
						DeletePricePlan:          priceScheduleDeps.DeletePricePlan,
						GetPricePlanInUseIDs:     priceScheduleDeps.GetPricePlanInUseIDs,
						ListPlans:                priceScheduleDeps.ListPlans,
						ListProducts:             priceScheduleDeps.ListProducts,
						ListProductPlans:         priceScheduleDeps.ListProductPlans,
						ListProductPricePlans:    priceScheduleDeps.ListProductPricePlans,
						CreateProductPricePlan:   priceScheduleDeps.CreateProductPricePlan,
						UpdateProductPricePlan:   priceScheduleDeps.UpdateProductPricePlan,
						DeleteProductPricePlan:   priceScheduleDeps.DeleteProductPricePlan,
					}
					priceschedulemod.NewModule(priceScheduleInventoryDeps).RegisterRoutes(ctx.Routes)
				}
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
			// 2026-04-27 plan-client-scope plan §6.1 / §6.2 — client name lookup
			// for the optional Client column on the plan list and for the
			// plan-drawer Client picker label resolution. Falls back to the
			// bare client_id when no use case is wired (e.g. tests).
			var listClientNames func(ctx context.Context) map[string]string
			if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
				lc := useCases.Entity.Client.ListClients.Execute
				listClientNames = func(fctx context.Context) map[string]string {
					out := map[string]string{}
					resp, err := lc(fctx, &clientpb.ListClientsRequest{})
					if err != nil {
						return out
					}
					for _, c := range resp.GetData() {
						label := c.GetName()
						if label == "" {
							if u := c.GetUser(); u != nil {
								label = u.GetFirstName() + " " + u.GetLastName()
							}
						}
						out[c.GetId()] = label
					}
					return out
				}
			}

			// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — job
			// template name lookup for the Job Template column on the plan
			// list. Mirrors the listClientNames pattern; falls back to the
			// bare job_template_id when the use case is unwired.
			var listJobTemplateNames func(ctx context.Context) map[string]string
			if useCases.Operation != nil && useCases.Operation.JobTemplate != nil && useCases.Operation.JobTemplate.ListJobTemplates != nil {
				ljt := useCases.Operation.JobTemplate.ListJobTemplates.Execute
				listJobTemplateNames = func(fctx context.Context) map[string]string {
					out := map[string]string{}
					resp, err := ljt(fctx, &jobtemplatepb.ListJobTemplatesRequest{})
					if err != nil {
						return out
					}
					for _, t := range resp.GetData() {
						if t == nil {
							continue
						}
						out[t.GetId()] = t.GetName()
					}
					return out
				}
			}

			planListDeps := &planlist.ListViewDeps{
				Routes:               planRoutes,
				Labels:               planLabels,
				CommonLabels:         ctx.Common,
				TableLabels:          centymoTableLabels,
				ListClientNames:      listClientNames,
				ListJobTemplateNames: listJobTemplateNames,
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
					// SetPlanActive uses raw DB update (proto3 omits false booleans)
					SetPlanActive: func(fctx context.Context, id string, active bool) error {
						_, err := db.Update(fctx, "plan", id, map[string]any{"active": active})
						return err
					},
				}
				// 2026-04-27 plan-client-scope plan §6.2 — Client picker support
				// + reference-checker lock state.
				if useCases.Entity != nil && useCases.Entity.Client != nil {
					if useCases.Entity.Client.ListClients != nil {
						planActionDeps.ListClients = useCases.Entity.Client.ListClients.Execute
					}
					if useCases.Entity.Client.SearchClientsByName != nil {
						planActionDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName.Execute
					}
				}
				if refChecker != nil {
					planActionDeps.GetPlanClientScopeLockedIDs = refChecker.GetPlanClientScopeLockedIDs
				}
				// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 —
				// JobTemplate select for Plan.job_template_id assignment.
				if useCases.Operation != nil && useCases.Operation.JobTemplate != nil && useCases.Operation.JobTemplate.ListJobTemplates != nil {
					planActionDeps.ListJobTemplates = useCases.Operation.JobTemplate.ListJobTemplates.Execute
				}
				ctx.Routes.GET(planRoutes.AddURL, planaction.NewAddAction(planActionDeps))
				ctx.Routes.POST(planRoutes.AddURL, planaction.NewAddAction(planActionDeps))
				ctx.Routes.GET(planRoutes.EditURL, planaction.NewEditAction(planActionDeps))
				ctx.Routes.POST(planRoutes.EditURL, planaction.NewEditAction(planActionDeps))
				ctx.Routes.POST(planRoutes.DeleteURL, planaction.NewDeleteAction(planActionDeps))
				ctx.Routes.POST(planRoutes.BulkDeleteURL, planaction.NewBulkDeleteAction(planActionDeps))
				ctx.Routes.POST(planRoutes.SetStatusURL, planaction.NewSetStatusAction(planActionDeps))
				ctx.Routes.POST(planRoutes.BulkSetStatusURL, planaction.NewBulkSetStatusAction(planActionDeps))
			}

			// Plan detail page + tab action
			if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
				planDetailDeps := &plandetail.DetailViewDeps{
					Routes:                     planRoutes,
					PriceSchedulePlanDetailURL: priceScheduleRoutes.PlanDetailURL,
					ReadPlan:                   useCases.Subscription.Plan.ReadPlan.Execute,
					Labels:                     planLabels,
					CommonLabels:               ctx.Common,
					TableLabels:                centymoTableLabels,
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
				if useCases.Product != nil && useCases.Product.Product != nil && useCases.Product.Product.ListProducts != nil {
					planDetailDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
				}
				if useCases.Product != nil && useCases.Product.ProductVariant != nil {
					planDetailDeps.ListProductVariants = useCases.Product.ProductVariant.ListProductVariants.Execute
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
				// 2026-04-28 plan-client-scope — Info tab Client row needs to
				// resolve the plan's client_id label and (optionally) link to
				// the entydad client-detail page.
				if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
					planDetailDeps.ListClients = useCases.Entity.Client.ListClients.Execute
				}
				planDetailDeps.ClientDetailURL = cfg.clientDetailURL
				// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — Info
				// tab JobTemplate row resolution.
				if useCases.Operation != nil && useCases.Operation.JobTemplate != nil && useCases.Operation.JobTemplate.ReadJobTemplate != nil {
					planDetailDeps.ReadJobTemplate = useCases.Operation.JobTemplate.ReadJobTemplate.Execute
				}
				ctx.Routes.GET(planRoutes.DetailURL, plandetail.NewView(planDetailDeps))
				ctx.Routes.GET(planRoutes.TabActionURL, plandetail.NewTabAction(planDetailDeps))

				// Plan-scoped PricePlan detail (/app/plans/detail/{id}/price/{ppid}).
				// Reuses the schedule-scoped detail body but anchors ActiveNav to
				// Services > Packages and points the breadcrumb back at the
				// package's package-prices tab. The {id} path value is plan_id;
				// the handler resolves schedule_id from the price_plan record.
				if planRoutes.PricePlanDetailURL != "" && useCases.Subscription.PricePlan != nil {
					// The plan detail's "Package prices" tab is registered under the
					// `pricePlan` key in plan tab labels; the lyngua professional
					// override surfaces it as the slug "package-prices" in the URL.
					packagePricesSlug := planLabels.Tabs.ResolveTabSlug("pricePlan")
					planScopedDeps := &priceschedulepricepldetail.DetailViewDeps{
						Routes:                 priceScheduleRoutes,
						ScheduleLabels:         priceScheduleLabels,
						PlanLabels:             pricePlanLabels,
						ProductPricePlanLabels: productPricePlanLabels,
						CommonLabels:           ctx.Common,
						TableLabels:            centymoTableLabels,
						ReadPricePlan:          useCases.Subscription.PricePlan.ReadPricePlan.Execute,
						// Mount overrides — keep the page anchored to Packages.
						ActiveNavOverride:      planRoutes.ActiveNav,
						ActiveSubNavOverride:   planRoutes.ActiveSubNav,
						PlanDetailBackURL:      planRoutes.DetailURL,
						PlanDetailBackTab:      packagePricesSlug,
						PlanScopedDetailURL:    planRoutes.PricePlanDetailURL,
						PlanScopedTabActionURL: planRoutes.PricePlanTabActionURL,
					}
					if useCases.Subscription.PriceSchedule != nil {
						planScopedDeps.ReadPriceSchedule = useCases.Subscription.PriceSchedule.ReadPriceSchedule.Execute
					}
					if useCases.Subscription.Plan != nil {
						planScopedDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
					}
					if useCases.Product != nil && useCases.Product.Product != nil {
						planScopedDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
					}
					if useCases.Product != nil && useCases.Product.ProductPlan != nil {
						planScopedDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
					}
					if useCases.Product != nil && useCases.Product.ProductVariant != nil {
						planScopedDeps.ListProductVariants = useCases.Product.ProductVariant.ListProductVariants.Execute
					}
					if useCases.Subscription.ProductPricePlan != nil {
						planScopedDeps.ListProductPricePlans = useCases.Subscription.ProductPricePlan.ListProductPricePlans.Execute
					}
					ctx.Routes.GET(planRoutes.PricePlanDetailURL, priceschedulepricepldetail.NewPlanScopedView(planScopedDeps))
					if planRoutes.PricePlanTabActionURL != "" {
						ctx.Routes.GET(planRoutes.PricePlanTabActionURL, priceschedulepricepldetail.NewPlanScopedTabAction(planScopedDeps))
					}
				}
				// Plan attachments
				if uploadFile != nil {
					ctx.Routes.GET(planRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planDetailDeps))
					ctx.Routes.POST(planRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planDetailDeps))
					ctx.Routes.POST(planRoutes.AttachmentDeleteURL, plandetail.NewAttachmentDeleteAction(planDetailDeps))
				}
				// PricePlan CRUD within plan detail
				if useCases.Subscription.PricePlan != nil && useCases.Subscription.PricePlan.CreatePricePlan != nil {
					ppActionDeps := &planaction.PricePlanDeps{
						Routes:              planRoutes,
						Labels:              planLabels,
						PricePlanLabels:     pricePlanLabels,
						PriceScheduleLabels: priceScheduleLabels,
						CommonLabels:        ctx.Common,
						CreatePricePlan: useCases.Subscription.PricePlan.CreatePricePlan.Execute,
						ReadPricePlan:   useCases.Subscription.PricePlan.ReadPricePlan.Execute,
						UpdatePricePlan: useCases.Subscription.PricePlan.UpdatePricePlan.Execute,
						DeletePricePlan: useCases.Subscription.PricePlan.DeletePricePlan.Execute,
					}
					if useCases.Subscription.PriceSchedule != nil {
						ppActionDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
					}
					if useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
						ppActionDeps.ReadPlan = useCases.Subscription.Plan.ReadPlan.Execute
					}
					// Plan §6.7 — ListClients powers the readonly schedule
					// label + lock tooltip when the parent Plan is client-scoped.
					if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
						ppActionDeps.ListClients = useCases.Entity.Client.ListClients.Execute
					}
					if useCases.Entity != nil && useCases.Entity.Location != nil && useCases.Entity.Location.ListLocations != nil {
						ppActionDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
					}
					if refChecker != nil {
						ppActionDeps.GetPricePlanInUseIDs = refChecker.GetPricePlanInUseIDs
					}
					if useCases.Product != nil && useCases.Product.Product != nil && useCases.Product.Product.ListProducts != nil {
						ppActionDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
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
					if useCases.Product.ProductVariant != nil {
						productPlanActionDeps.ListProductVariants = useCases.Product.ProductVariant.ListProductVariants.Execute
					}
					if useCases.Product.ProductVariantOption != nil {
						productPlanActionDeps.ListProductVariantOptions = useCases.Product.ProductVariantOption.ListProductVariantOptions.Execute
					}
					if useCases.Product.ProductOptionValue != nil {
						productPlanActionDeps.ListProductOptionValues = useCases.Product.ProductOptionValue.ListProductOptionValues.Execute
					}
					if useCases.Product.ProductOption != nil {
						productPlanActionDeps.ListProductOptions = useCases.Product.ProductOption.ListProductOptions.Execute
					}
					ctx.Routes.GET(planRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(productPlanActionDeps))
					ctx.Routes.POST(planRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(productPlanActionDeps))
					ctx.Routes.GET(planRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(productPlanActionDeps))
					ctx.Routes.POST(planRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(productPlanActionDeps))
					ctx.Routes.POST(planRoutes.ProductPlanDeleteURL, planaction.NewProductPlanDeleteAction(productPlanActionDeps))
					ctx.Routes.GET(planRoutes.ProductPlanPickerURL, planaction.NewProductPlanPickerAction(productPlanActionDeps))
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
					Routes:               planBundleRoutes,
					Labels:               planLabels,
					CommonLabels:         ctx.Common,
					TableLabels:          centymoTableLabels,
					ListClientNames:      listClientNames,
					ListJobTemplateNames: listJobTemplateNames,
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
						// SetPlanActive uses raw DB update (proto3 omits false booleans)
						SetPlanActive: func(fctx context.Context, id string, active bool) error {
							_, err := db.Update(fctx, "plan", id, map[string]any{"active": active})
							return err
						},
					}
					// 2026-04-27 plan-client-scope plan §6.2 — same Client picker
					// + lock state on the bundle mount.
					if useCases.Entity != nil && useCases.Entity.Client != nil {
						if useCases.Entity.Client.ListClients != nil {
							planBundleActionDeps.ListClients = useCases.Entity.Client.ListClients.Execute
						}
						if useCases.Entity.Client.SearchClientsByName != nil {
							planBundleActionDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName.Execute
						}
					}
					if refChecker != nil {
						planBundleActionDeps.GetPlanClientScopeLockedIDs = refChecker.GetPlanClientScopeLockedIDs
					}
					// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 —
					// JobTemplate select on the bundle mount.
					if useCases.Operation != nil && useCases.Operation.JobTemplate != nil && useCases.Operation.JobTemplate.ListJobTemplates != nil {
						planBundleActionDeps.ListJobTemplates = useCases.Operation.JobTemplate.ListJobTemplates.Execute
					}
					ctx.Routes.GET(planBundleRoutes.AddURL, planaction.NewAddAction(planBundleActionDeps))
					ctx.Routes.POST(planBundleRoutes.AddURL, planaction.NewAddAction(planBundleActionDeps))
					ctx.Routes.GET(planBundleRoutes.EditURL, planaction.NewEditAction(planBundleActionDeps))
					ctx.Routes.POST(planBundleRoutes.EditURL, planaction.NewEditAction(planBundleActionDeps))
					ctx.Routes.POST(planBundleRoutes.DeleteURL, planaction.NewDeleteAction(planBundleActionDeps))
					ctx.Routes.POST(planBundleRoutes.BulkDeleteURL, planaction.NewBulkDeleteAction(planBundleActionDeps))
					ctx.Routes.POST(planBundleRoutes.SetStatusURL, planaction.NewSetStatusAction(planBundleActionDeps))
					ctx.Routes.POST(planBundleRoutes.BulkSetStatusURL, planaction.NewBulkSetStatusAction(planBundleActionDeps))
				}

				if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
					planBundleDetailDeps := &plandetail.DetailViewDeps{
						Routes:                     planBundleRoutes,
						PriceSchedulePlanDetailURL: priceScheduleRoutes.PlanDetailURL,
						ReadPlan:                   useCases.Subscription.Plan.ReadPlan.Execute,
						Labels:                     planLabels,
						CommonLabels:               ctx.Common,
						TableLabels:                centymoTableLabels,
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
					if useCases.Product != nil && useCases.Product.Product != nil && useCases.Product.Product.ListProducts != nil {
						planBundleDetailDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
					}
					if useCases.Product != nil && useCases.Product.ProductVariant != nil {
						planBundleDetailDeps.ListProductVariants = useCases.Product.ProductVariant.ListProductVariants.Execute
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
					// 2026-04-28 plan-client-scope — same Info tab Client row
					// wiring on the bundle mount.
					if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
						planBundleDetailDeps.ListClients = useCases.Entity.Client.ListClients.Execute
					}
					planBundleDetailDeps.ClientDetailURL = cfg.clientDetailURL
					// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 —
					// Info tab JobTemplate row on the bundle mount.
					if useCases.Operation != nil && useCases.Operation.JobTemplate != nil && useCases.Operation.JobTemplate.ReadJobTemplate != nil {
						planBundleDetailDeps.ReadJobTemplate = useCases.Operation.JobTemplate.ReadJobTemplate.Execute
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
							Routes:              planBundleRoutes,
							Labels:              planLabels,
							PricePlanLabels:     pricePlanLabels,
							PriceScheduleLabels: priceScheduleLabels,
							CommonLabels:        ctx.Common,
							CreatePricePlan: useCases.Subscription.PricePlan.CreatePricePlan.Execute,
							ReadPricePlan:   useCases.Subscription.PricePlan.ReadPricePlan.Execute,
							UpdatePricePlan: useCases.Subscription.PricePlan.UpdatePricePlan.Execute,
							DeletePricePlan: useCases.Subscription.PricePlan.DeletePricePlan.Execute,
						}
						if useCases.Subscription.PriceSchedule != nil {
							ppBundleDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
						}
						if useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
							ppBundleDeps.ReadPlan = useCases.Subscription.Plan.ReadPlan.Execute
						}
						// Plan §6.7 — ListClients powers the readonly schedule
						// label + lock tooltip on the bundle-mount drawer.
						if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
							ppBundleDeps.ListClients = useCases.Entity.Client.ListClients.Execute
						}
						if useCases.Entity != nil && useCases.Entity.Location != nil && useCases.Entity.Location.ListLocations != nil {
							ppBundleDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
						}
						if refChecker != nil {
							ppBundleDeps.GetPricePlanInUseIDs = refChecker.GetPricePlanInUseIDs
						}
						if useCases.Product != nil && useCases.Product.Product != nil && useCases.Product.Product.ListProducts != nil {
							ppBundleDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
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
					// Bundle-mount sibling of services-mount `productPlanActionDeps` (~line 1111).
					// Keep these two registrations field-for-field identical (only Routes
					// differs). Unlike PricePlanDeps, ProductPlanDeps has a single Labels
					// field — all form-label data is nested under centymo.PlanLabels
					// (`Labels.ProductPlanForm`), so threading `Labels: planLabels` is
					// sufficient. If a future change adds a separate label struct (e.g.
					// ProductPlanLabels), thread it into BOTH registrations.
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
						if useCases.Product.ProductVariant != nil {
							ppBundleProductPlanDeps.ListProductVariants = useCases.Product.ProductVariant.ListProductVariants.Execute
						}
						if useCases.Product.ProductVariantOption != nil {
							ppBundleProductPlanDeps.ListProductVariantOptions = useCases.Product.ProductVariantOption.ListProductVariantOptions.Execute
						}
						if useCases.Product.ProductOptionValue != nil {
							ppBundleProductPlanDeps.ListProductOptionValues = useCases.Product.ProductOptionValue.ListProductOptionValues.Execute
						}
						if useCases.Product.ProductOption != nil {
							ppBundleProductPlanDeps.ListProductOptions = useCases.Product.ProductOption.ListProductOptions.Execute
						}
						ctx.Routes.GET(planBundleRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(ppBundleProductPlanDeps))
						ctx.Routes.POST(planBundleRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(ppBundleProductPlanDeps))
						ctx.Routes.GET(planBundleRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(ppBundleProductPlanDeps))
						ctx.Routes.POST(planBundleRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(ppBundleProductPlanDeps))
						ctx.Routes.POST(planBundleRoutes.ProductPlanDeleteURL, planaction.NewProductPlanDeleteAction(ppBundleProductPlanDeps))
						ctx.Routes.GET(planBundleRoutes.ProductPlanPickerURL, planaction.NewProductPlanPickerAction(ppBundleProductPlanDeps))
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
			if refChecker != nil {
				subListDeps.GetInUseIDs = refChecker.GetSubscriptionInUseIDs
			}
			ctx.Routes.GET(subscriptionRoutes.ListURL, subscriptionlist.NewView(subListDeps))
			// Table-only endpoint — used by sheet.js refreshTable() after
			// activate/deactivate/delete so HTMX swaps the table-card partial,
			// not the whole page.
			if subscriptionRoutes.TableURL != "" {
				ctx.Routes.GET(subscriptionRoutes.TableURL, subscriptionlist.NewTableView(subListDeps))
			}

			// Subscription CRUD actions
			if useCases.Subscription != nil && useCases.Subscription.Subscription != nil && useCases.Subscription.Subscription.CreateSubscription != nil {
				subActionDeps := &subscriptionaction.Deps{
					Routes:             subscriptionRoutes,
					Labels:             subscriptionLabels,
					CreateSubscription: useCases.Subscription.Subscription.CreateSubscription.Execute,
					ReadSubscription:   useCases.Subscription.Subscription.ReadSubscription.Execute,
					UpdateSubscription: useCases.Subscription.Subscription.UpdateSubscription.Execute,
					DeleteSubscription: useCases.Subscription.Subscription.DeleteSubscription.Execute,
					// SetSubscriptionActive uses raw DB update (proto3 omits false booleans)
					SetSubscriptionActive: func(fctx context.Context, id string, active bool) error {
						_, err := db.Update(fctx, "subscription", id, map[string]any{"active": active})
						return err
					},
				}
				if refChecker != nil {
					subActionDeps.GetInUseIDs = refChecker.GetSubscriptionInUseIDs
				}
				if useCases.Subscription.Subscription.GetSubscriptionItemPageData != nil {
					subActionDeps.GetSubscriptionItemPageData = useCases.Subscription.Subscription.GetSubscriptionItemPageData.Execute
				}
				if useCases.Entity != nil && useCases.Entity.Client != nil {
					subActionDeps.ListClients = useCases.Entity.Client.ListClients.Execute
					if useCases.Entity.Client.SearchClientsByName != nil {
						subActionDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName.Execute
					}
				}
				if useCases.Subscription.Plan != nil {
					subActionDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
					if useCases.Subscription.Plan.ReadPlan != nil {
						subActionDeps.ReadPlan = useCases.Subscription.Plan.ReadPlan.Execute
					}
					if useCases.Subscription.Plan.SearchPlansByName != nil {
						subActionDeps.SearchPlansByName = useCases.Subscription.Plan.SearchPlansByName.Execute
					}
				}
				if useCases.Subscription.PricePlan != nil {
					subActionDeps.ListPricePlans = useCases.Subscription.PricePlan.ListPricePlans.Execute
					if useCases.Subscription.PricePlan.ReadPricePlan != nil {
						subActionDeps.ReadPricePlan = useCases.Subscription.PricePlan.ReadPricePlan.Execute
					}
				}
				if useCases.Subscription.PriceSchedule != nil && useCases.Subscription.PriceSchedule.ListPriceSchedules != nil {
					subActionDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
				}
				// Wire the espyna recognize-revenue use case so the new
				// drawer + the existing manual-revenue-add auto-populate
				// path share one source of truth.
				if useCases.Revenue != nil && useCases.Revenue.Revenue != nil &&
					useCases.Revenue.Revenue.RecognizeRevenueFromSubscription != nil {
					subActionDeps.RecognizeRevenueFromSubscription =
						useCases.Revenue.Revenue.RecognizeRevenueFromSubscription.Execute
				}

				// 2026-04-27 plan-client-scope plan §4 / §6.5 — wire the
				// CustomizePlanForClient use case via a thin adapter that
				// converts the centymo-side request shape to whatever the
				// espyna-golang use case expects. The espyna use case is in
				// flight in a parallel agent's branch; until its signature
				// stabilizes the Plan.CustomizePlanForClient pointer is
				// optional here — when nil, the customize CTA returns
				// `customize_failed` and the drawer error toast surfaces.
				subActionDeps.CustomClientPriceScheduleLabelSuffix =
					priceScheduleLabels.Form.CustomClientPriceScheduleLabelSuffix
				wireCustomizePlanForClient(useCases, subActionDeps)

				// 2026-04-29 milestone-billing plan §5 / Phase D — wire the
				// BillingEvent server through to the recognize drawer +
				// mark-ready/waive handlers. nil-safe: the espyna subscription
				// composition exposes the server pointer directly (no use-case
				// wrapper yet).
				if useCases.Subscription.BillingEvent != nil {
					be := useCases.Subscription.BillingEvent
					subActionDeps.ListBillingEventsBySubscription = be.ListBySubscription
					subActionDeps.SetBillingEventStatus = be.SetStatus
				}

				// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / Phase D —
				// wire the JobTemplate read deps that drive the Spawn Jobs
				// section detection on the subscription create drawer + the
				// retroactive spawn drawer. nil-safe.
				if useCases.Operation != nil {
					if uc := useCases.Operation.JobTemplate; uc != nil && uc.ReadJobTemplate != nil {
						subActionDeps.ReadJobTemplate = uc.ReadJobTemplate.Execute
					}
					if uc := useCases.Operation.JobTemplatePhase; uc != nil && uc.ListByJobTemplate != nil {
						subActionDeps.ListJobTemplatePhases = uc.ListByJobTemplate.Execute
					}
					if uc := useCases.Operation.JobTemplateTask; uc != nil && uc.ListByPhase != nil {
						subActionDeps.ListJobTemplateTasks = uc.ListByPhase.Execute
					}
					if useCases.Operation.JobTemplateRelation != nil {
						subActionDeps.ListJobTemplateRelations = useCases.Operation.JobTemplateRelation.ListByParent
					}
				}
				if useCases.Subscription.MaterializeJobsForSubscription != nil {
					subActionDeps.MaterializeJobsForSubscription = func(fctx context.Context, subID string, spawn bool) (int, string, error) {
						resp, err := consumer.MaterializeJobsForSubscription(useCases, fctx, &consumer.MaterializeJobsForSubscriptionRequest{
							SubscriptionID: subID,
							SpawnJobs:      spawn,
						})
						if err != nil {
							return 0, "", err
						}
						if resp == nil {
							return 0, "", nil
						}
						return resp.JobCount, resp.SkippedReason, nil
					}
				}
				// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D —
				// wire espyna's MaterializeInstanceJobsForSubscription through
				// a centymo-side adapter so the Operations tab "Spawn this
				// cycle now" + "Backfill missing cycles" handlers can call it
				// without importing espyna directly. nil-safe: the cycle-spawn
				// and backfill action handlers gate on the adapter pointer.
				if useCases.Subscription.MaterializeInstanceJobsForSubscription != nil {
					subActionDeps.MaterializeInstanceJobsForSubscription = func(fctx context.Context, req *subscriptionaction.MaterializeInstanceJobsRequest) (*subscriptionaction.MaterializeInstanceJobsResponse, error) {
						if req == nil {
							return nil, nil
						}
						resp, err := consumer.MaterializeInstanceJobsForSubscription(useCases, fctx, &consumer.MaterializeInstanceJobsForSubscriptionRequest{
							SubscriptionID:   req.SubscriptionID,
							CyclePeriodStart: req.CyclePeriodStart,
							Backfill:         req.Backfill,
						})
						if err != nil {
							return nil, err
						}
						if resp == nil {
							return &subscriptionaction.MaterializeInstanceJobsResponse{}, nil
						}
						return &subscriptionaction.MaterializeInstanceJobsResponse{
							SpawnedCycleCount:         resp.SpawnedCycleCount,
							SpawnedJobCount:           resp.SpawnedJobCount,
							OnceAtStartJobCount:       resp.OnceAtStartJobCount,
							EngagementWasNewlyCreated: resp.EngagementWasNewlyCreated,
							SkippedReason:             resp.SkippedReason,
							BackfillCappedAt:          resp.BackfillCappedAt,
						}, nil
					}
				}

				ctx.Routes.GET(subscriptionRoutes.AddURL, subscriptionaction.NewAddAction(subActionDeps))
				ctx.Routes.POST(subscriptionRoutes.AddURL, subscriptionaction.NewAddAction(subActionDeps))
				ctx.Routes.GET(subscriptionRoutes.EditURL, subscriptionaction.NewEditAction(subActionDeps))
				ctx.Routes.POST(subscriptionRoutes.EditURL, subscriptionaction.NewEditAction(subActionDeps))
				ctx.Routes.POST(subscriptionRoutes.DeleteURL, subscriptionaction.NewDeleteAction(subActionDeps))
				ctx.Routes.POST(subscriptionRoutes.BulkDeleteURL, subscriptionaction.NewBulkDeleteAction(subActionDeps))
				ctx.Routes.POST(subscriptionRoutes.SetStatusURL, subscriptionaction.NewSetStatusAction(subActionDeps))
				ctx.Routes.POST(subscriptionRoutes.BulkSetStatusURL, subscriptionaction.NewBulkSetStatusAction(subActionDeps))
				// Recognize-revenue drawer (GET = preview, POST = generate). Per
				// plan §11.1, POST returns HTMXSuccess + refresh-invoices so the
				// invoices table refreshes inline.
				if subActionDeps.RecognizeRevenueFromSubscription != nil && subscriptionRoutes.RecognizeURL != "" {
					ctx.Routes.GET(subscriptionRoutes.RecognizeURL, subscriptionaction.NewRecognizeAction(subActionDeps))
					ctx.Routes.POST(subscriptionRoutes.RecognizeURL, subscriptionaction.NewRecognizeAction(subActionDeps))
				}
				// 2026-04-27 plan-client-scope plan §6.5 — Customize package
				// CTA on subscription detail's Package tab.
				if subscriptionRoutes.CustomizePackageURL != "" {
					ctx.Routes.POST(subscriptionRoutes.CustomizePackageURL, subscriptionaction.NewCustomizePackageAction(subActionDeps))
				}
				// 2026-04-29 milestone-billing plan §5 / Phase D — mark-ready +
				// waive handlers for BillingEvent rows on the subscription
				// Package tab. Only registered when the BillingEvent server
				// is wired (espyna subscription provider has the adapter).
				if subActionDeps.SetBillingEventStatus != nil {
					if subscriptionRoutes.MilestoneMarkReadyURL != "" {
						ctx.Routes.POST(subscriptionRoutes.MilestoneMarkReadyURL,
							subscriptionaction.NewMilestoneMarkReadyAction(
								subActionDeps.SetBillingEventStatus,
								subActionDeps.Labels.Errors))
					}
					if subscriptionRoutes.MilestoneWaiveURL != "" {
						ctx.Routes.POST(subscriptionRoutes.MilestoneWaiveURL,
							subscriptionaction.NewMilestoneWaiveAction(
								subActionDeps.SetBillingEventStatus,
								subActionDeps.Labels.Errors))
					}
				}

				// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / Phase D —
				// HTMX partial that re-renders the Spawn Jobs section on Plan
				// select change + the retroactive spawn drawer.
				if subscriptionRoutes.SpawnJobsPartialURL != "" {
					ctx.Routes.GET(subscriptionRoutes.SpawnJobsPartialURL, subscriptionaction.NewSpawnJobsPartialAction(subActionDeps))
				}
				if subscriptionRoutes.SpawnJobsURL != "" {
					ctx.Routes.GET(subscriptionRoutes.SpawnJobsURL, subscriptionaction.NewSpawnJobsAction(subActionDeps))
					ctx.Routes.POST(subscriptionRoutes.SpawnJobsURL, subscriptionaction.NewSpawnJobsAction(subActionDeps))
				}
				// 2026-04-30 cyclic-subscription-jobs plan §5.3 — Operations
				// tab CTAs: "Spawn this cycle now" (POST) + "Backfill missing
				// cycles" (GET drawer / POST commit). Both gate on the adapter
				// being wired (handlers also re-check internally). The detail
				// page tab template hides the buttons when the URL fields are
				// empty, so nil-safety is double-bottomed.
				if subscriptionRoutes.SpawnCycleJobsURL != "" {
					ctx.Routes.POST(subscriptionRoutes.SpawnCycleJobsURL, subscriptionaction.NewSpawnCycleJobsAction(subActionDeps))
				}
				if subscriptionRoutes.BackfillCycleJobsURL != "" {
					ctx.Routes.GET(subscriptionRoutes.BackfillCycleJobsURL, subscriptionaction.NewBackfillCyclesAction(subActionDeps))
					ctx.Routes.POST(subscriptionRoutes.BackfillCycleJobsURL, subscriptionaction.NewBackfillCyclesAction(subActionDeps))
				}
				// 2026-05-01 ad-hoc-subscription-billing — Request Usage CTA on
				// the AD_HOC Operations tab. Pool-Generate-Invoice reuses the
				// existing /action/subscription/recognize-revenue/{id} endpoint
				// (espyna's executeAdHoc dispatches by PricePlan kind).
				if subscriptionRoutes.RequestUsageURL != "" {
					ctx.Routes.POST(subscriptionRoutes.RequestUsageURL, subscriptionaction.NewRequestUsageAction(subActionDeps))
				}
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
				if useCases.Subscription.Subscription.GetSubscriptionItemPageData != nil {
					subDetailDeps.GetSubscriptionItemPageData = useCases.Subscription.Subscription.GetSubscriptionItemPageData.Execute
				}
				if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ReadClient != nil {
					subDetailDeps.ReadClient = useCases.Entity.Client.ReadClient.Execute
				}
				if useCases.Revenue != nil && useCases.Revenue.Revenue != nil && useCases.Revenue.Revenue.GetRevenueListPageData != nil {
					subDetailDeps.GetRevenueListPageData = useCases.Revenue.Revenue.GetRevenueListPageData.Execute
				}
				// 2026-04-29 milestone-billing — wire BillingEvent listing into
				// the subscription detail Package tab.
				if useCases.Subscription.BillingEvent != nil {
					subDetailDeps.ListBillingEventsBySubscription = useCases.Subscription.BillingEvent.ListBySubscription
				}
				// 2026-04-29 auto-spawn-jobs-from-subscription Phase D — wire
				// the Operations tab data ops + spawn-jobs CTA URL.
				if useCases.Operation != nil {
					if uc := useCases.Operation.Job; uc != nil && uc.GetJobsByOrigin != nil {
						subDetailDeps.GetJobsByOrigin = uc.GetJobsByOrigin.Execute
					}
					if uc := useCases.Operation.JobPhase; uc != nil && uc.ListByJob != nil {
						subDetailDeps.ListJobPhasesByJob = uc.ListByJob.Execute
					}
				}
				subDetailDeps.SpawnJobsURL = subscriptionRoutes.SpawnJobsURL
				subDetailDeps.JobDetailURL = cfg.jobDetailURL
				subDetailDeps.ClientDetailURL = cfg.clientDetailURL
				ctx.Routes.GET(subscriptionRoutes.DetailURL, subscriptiondetail.NewView(subDetailDeps))
				ctx.Routes.GET(subscriptionRoutes.TabActionURL, subscriptiondetail.NewTabAction(subDetailDeps))
				// Nested route — same view, breadcrumb activated via path param.
				if subscriptionRoutes.UnderClientDetailURL != "" {
					ctx.Routes.GET(subscriptionRoutes.UnderClientDetailURL, subscriptiondetail.NewView(subDetailDeps))
				}
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
				collDeps := &collectionmod.ModuleDeps{
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
				}
				wireCashDashboard(collDeps, useCases)
				collectionmod.NewModule(collDeps).RegisterRoutes(ctx.Routes)
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
			// SPS Wave 4 — Recognition + Accrual tabs on the expense detail page.
			// Nil-safe — when the use case is missing, the tab renders an empty state.
			if useCases.Expenditure != nil && useCases.Expenditure.ExpenseRecognition != nil {
				if uc := useCases.Expenditure.ExpenseRecognition.ReadExpenseRecognition; uc != nil {
					expDeps.ReadExpenseRecognition = uc.Execute
				}
			}
			if useCases.Expenditure != nil && useCases.Expenditure.AccruedExpense != nil {
				if uc := useCases.Expenditure.AccruedExpense.ListAccruedExpenses; uc != nil {
					expDeps.ListAccruedExpenses = uc.Execute
				}
			}
			expDeps.ExpenseRecognitionDetailURL = expenseRecognitionRoutes.DetailURL
			expDeps.AccruedExpenseDetailURL = accruedExpenseRoutes.DetailURL
			// RecognizeFromExpenditureURL is the espyna trigger surfaced as the
			// empty-state CTA. The espyna RecognizeFromExpenditure use case is
			// exposed on the API surface; the centymo route_config does not
			// expose a /action/* mirror because recognition is created by use
			// case (no UI form). Leaving empty by default; verticals can wire
			// a custom trigger URL via lyngua override.
			wirePurchaseDashboard(expDeps, useCases)
			wireExpenseDashboard(expDeps, useCases)
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

		// =====================================================================
		// 20260427-supplier-commitments — five new modules (P3a + P3b)
		//
		// Mirrors the expendituremod pattern (lines ~1937 above): construct
		// ModuleDeps, plumb each available use case through, register routes.
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
				Routes:       supplierContractRoutes,
				Labels:       supplierContractLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
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
			scDeps.PriceScheduleListURL = supplierContractPriceScheduleRoutes.ListURL
			scDeps.PriceScheduleDetailURL = supplierContractPriceScheduleRoutes.DetailURL
			scDeps.PriceScheduleAddURL = supplierContractPriceScheduleRoutes.AddURL
			suppliercontractmod.NewModule(scDeps).RegisterRoutes(ctx.Routes)
		}

		// SupplierContractLine module — child rows of SupplierContract.
		if cfg.wantSupplierContractLine() {
			sclDeps := &suppliercontractlinemod.ModuleDeps{
				Routes:       supplierContractRoutes,
				Labels:       supplierContractLabels,
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
				Routes:       procurementRequestRoutes,
				Labels:       procurementRequestLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
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
			procurementrequestmod.NewModule(prDeps).RegisterRoutes(ctx.Routes)
		}

		// ProcurementRequestLine module — child rows of ProcurementRequest.
		if cfg.wantProcurementRequestLine() {
			prlDeps := &procurementrequestlinemod.ModuleDeps{
				Routes:       procurementRequestRoutes,
				Labels:       procurementRequestLabels,
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
				Routes:       procurementRoutes,
				Labels:       procurementLabels,
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

		// =====================================================================
		// SPS Wave 4 — supplier-side pricing graph + accrual layer.
		//
		// Six new view modules wired in dependency order:
		//   1. SupplierContractPriceSchedule (master, mirrors price_schedule)
		//   2. SupplierContractPriceScheduleLine (inline child of #1)
		//   3. ExpenseRecognition (master, no add/edit drawer — created BY use case)
		//   4. ExpenseRecognitionLine (inline child of #3)
		//   5. AccruedExpense (master, manual create + AccrueFromContract)
		//   6. AccruedExpenseSettlement (inline child of #5, shares parent routes)
		//
		// Mirrors the supplier_contract / supplier_contract_line wiring pattern
		// above. All use-case threading is nil-safe — missing use cases fall back
		// to the empty/disabled view state instead of panicking.
		// =====================================================================

		// SupplierContractPriceSchedule module
		if cfg.wantSupplierContractPriceSchedule() {
			scpsDeps := &suppliercontractpriceschedulemod.ModuleDeps{
				Routes:       supplierContractPriceScheduleRoutes,
				Labels:       supplierContractPriceScheduleLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
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
			suppliercontractpriceschedulemod.NewModule(scpsDeps).RegisterRoutes(ctx.Routes)
		}

		// SupplierContractPriceScheduleLine module — child rows of SupplierContractPriceSchedule.
		if cfg.wantSupplierContractPriceScheduleLine() {
			scpslDeps := &suppliercontractpricescheduleinemod.ModuleDeps{
				Routes:       supplierContractPriceScheduleRoutes,
				Labels:       supplierContractPriceScheduleLabels,
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

		// ExpenseRecognition module — no Add/Edit drawer (created BY use case).
		if cfg.wantExpenseRecognition() {
			erDeps := &expenserecognitionmod.ModuleDeps{
				Routes:       expenseRecognitionRoutes,
				Labels:       expenseRecognitionLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
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
			expenserecognitionmod.NewModule(erDeps).RegisterRoutes(ctx.Routes)
		}

		// ExpenseRecognitionLine module — inline child of ExpenseRecognition.
		if cfg.wantExpenseRecognitionLine() {
			erlDeps := &expenserecognitionlinemod.ModuleDeps{
				Routes:       expenseRecognitionRoutes,
				Labels:       expenseRecognitionLabels,
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

		// AccruedExpense module
		if cfg.wantAccruedExpense() {
			aeDeps := &accruedexpensemod.ModuleDeps{
				Routes:       accruedExpenseRoutes,
				Labels:       accruedExpenseLabels,
				CommonLabels: ctx.Common,
				TableLabels:  centymoTableLabels,
			}
			if useCases.Expenditure != nil && useCases.Expenditure.AccruedExpense != nil {
				uc := useCases.Expenditure.AccruedExpense
				if uc.ListAccruedExpenses != nil {
					aeDeps.ListAccruedExpenses = uc.ListAccruedExpenses.Execute
				}
				if uc.ReadAccruedExpense != nil {
					aeDeps.ReadAccruedExpense = uc.ReadAccruedExpense.Execute
				}
				if uc.CreateAccruedExpense != nil {
					aeDeps.CreateAccruedExpense = uc.CreateAccruedExpense.Execute
				}
				if uc.UpdateAccruedExpense != nil {
					aeDeps.UpdateAccruedExpense = uc.UpdateAccruedExpense.Execute
				}
				if uc.DeleteAccruedExpense != nil {
					aeDeps.DeleteAccruedExpense = uc.DeleteAccruedExpense.Execute
				}
				if uc.SettleAccrual != nil {
					settleUC := uc.SettleAccrual
					aeDeps.SettleAccrual = func(fctx context.Context, req *accruedexpensepb.SettleAccrualRequest) error {
						_, err := settleUC.SettleAccrual(fctx, req)
						return err
					}
				}
				if uc.ReverseAccrual != nil {
					reverseUC := uc.ReverseAccrual
					aeDeps.ReverseAccrual = func(fctx context.Context, id, reason string) error {
						req := &accruedexpensepb.ReverseAccrualRequest{AccruedExpenseId: id}
						if reason != "" {
							req.Reason = &reason
						}
						_, err := reverseUC.Execute(fctx, req)
						return err
					}
				}
				if uc.AccrueFromContract != nil {
					afcUC := uc.AccrueFromContract
					aeDeps.AccrueFromContract = func(fctx context.Context, req *accruedexpensepb.AccrueFromContractRequest) (*accruedexpensepb.AccrueFromContractResponse, error) {
						return afcUC.Execute(fctx, req)
					}
				}
			}
			// Inline child — settlements.
			if useCases.Expenditure != nil && useCases.Expenditure.AccruedExpenseSettlement != nil {
				if uc := useCases.Expenditure.AccruedExpenseSettlement.ListAccruedExpenseSettlements; uc != nil {
					aeDeps.ListAccruedExpenseSettlements = uc.Execute
				}
			}
			// Dropdowns for the manual-create drawer + filter pickers.
			if useCases.Entity != nil && useCases.Entity.Supplier != nil &&
				useCases.Entity.Supplier.ListSuppliers != nil {
				aeDeps.ListSuppliers = useCases.Entity.Supplier.ListSuppliers.Execute
			}
			if useCases.Expenditure != nil && useCases.Expenditure.SupplierContract != nil &&
				useCases.Expenditure.SupplierContract.ListSupplierContracts != nil {
				aeDeps.ListSupplierContracts = useCases.Expenditure.SupplierContract.ListSupplierContracts.Execute
			}
			accruedexpensemod.NewModule(aeDeps).RegisterRoutes(ctx.Routes)
		}

		// AccruedExpenseSettlement module — inline child of AccruedExpense (shares parent routes).
		if cfg.wantAccruedExpenseSettlement() {
			aesDeps := &accruedexpensesettlementmod.ModuleDeps{
				Routes:       accruedExpenseRoutes,
				Labels:       accruedExpenseLabels,
				CommonLabels: ctx.Common,
			}
			if useCases.Expenditure != nil && useCases.Expenditure.AccruedExpenseSettlement != nil {
				uc := useCases.Expenditure.AccruedExpenseSettlement
				if uc.CreateAccruedExpenseSettlement != nil {
					aesDeps.CreateAccruedExpenseSettlement = uc.CreateAccruedExpenseSettlement.Execute
				}
				if uc.ReadAccruedExpenseSettlement != nil {
					aesDeps.ReadAccruedExpenseSettlement = uc.ReadAccruedExpenseSettlement.Execute
				}
				if uc.UpdateAccruedExpenseSettlement != nil {
					aesDeps.UpdateAccruedExpenseSettlement = uc.UpdateAccruedExpenseSettlement.Execute
				}
				if uc.DeleteAccruedExpenseSettlement != nil {
					aesDeps.DeleteAccruedExpenseSettlement = uc.DeleteAccruedExpenseSettlement.Execute
				}
			}
			// Settling-Expenditure picker for the settlement drawer form.
			if useCases.Expenditure != nil && useCases.Expenditure.Expenditure != nil &&
				useCases.Expenditure.Expenditure.ListExpenditures != nil {
				aesDeps.ListExpenditures = useCases.Expenditure.Expenditure.ListExpenditures.Execute
			}
			accruedexpensesettlementmod.NewModule(aesDeps).RegisterRoutes(ctx.Routes)
		}

		log.Println("  centymo commerce domain initialized")
		return nil
	}
}

// wireCustomizePlanForClient threads the espyna Plan.CustomizePlanForClient use
// case into the centymo subscription action Deps. The espyna side ships an
// independent request/response shape; the centymo side uses an in-package
// type-narrow shape so its handlers don't depend directly on the espyna
// generated proto/use-case structs.
//
// When the use case isn't wired (composition layer didn't initialize it), we
// leave the function pointer nil; the handler falls through to a generic
// `customize_failed` toast.
//
// 2026-04-27 plan-client-scope plan §4. Same adapter pattern as
// RecognizeRevenueFromSubscription above.
func wireCustomizePlanForClient(useCases *consumer.UseCases, subActionDeps *subscriptionaction.Deps) {
	if useCases == nil || useCases.Subscription == nil || useCases.Subscription.Plan == nil {
		return
	}
	customizeUC := useCases.Subscription.Plan.CustomizePlanForClient
	if customizeUC == nil {
		return
	}
	_ = customizeUC
	subActionDeps.CustomizePlanForClient = func(
		ctx context.Context, req *subscriptionaction.CustomizePlanForClientRequest,
	) (*subscriptionaction.CustomizePlanForClientResponse, error) {
		resp, err := consumer.CustomizePlanForClient(useCases, ctx, &consumer.CustomizePlanForClientRequest{
			SourcePlanID:      req.SourcePlanID,
			SourcePricePlanID: req.SourcePricePlanID,
			ClientID:          req.ClientID,
			SubscriptionID:    req.SubscriptionID,
			NewScheduleName:   req.NewScheduleName,
		})
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return &subscriptionaction.CustomizePlanForClientResponse{}, nil
		}
		return &subscriptionaction.CustomizePlanForClientResponse{
			NewPlanID:      resp.NewPlanID,
			NewPricePlanID: resp.NewPricePlanID,
			NewScheduleID:  resp.NewScheduleID,
			Reused:         resp.Reused,
		}, nil
	}
}
