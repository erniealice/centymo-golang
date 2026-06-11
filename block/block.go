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
// # File layout
//
// The wiring is split across companion files (all in package block):
//
//   - block.go               — Block() entry point, inline modules (inventory, collection, etc.)
//   - options.go             — BlockOption, WithX() funcs, blockConfig
//   - revenue_run.go         — wireRevenueRunModules (revenue run + lines + actions)
//   - supplier_commitment.go — wireSupplierCommitmentModules (PO + receipt + returns)
//   - supplier_contract_price_schedule.go — wireSupplierContractPriceScheduleModules
//   - expense_recognition.go — wireExpenseRecognitionModules (expense recognition + lines)
//   - accrued_expense.go     — wireAccruedExpenseModules (accrued expense + settlement)
//   - supplier_subscription.go — wireSupplierSubscriptionModules (cost/supplier plans + subscriptions)
//   - product.go             — wireProductModules (Product 3-mount + ProductLine 2-mount)
//   - plan.go                — wirePlanModules (PricePlan, PriceSchedule, PriceList, Plan + PlanBundle)
//   - subscription.go        — wireSubscriptionModule (subscription CRUD + detail)
//   - wiring.go              — shared wireServiceDashboard helper
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
	"os"

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"

	"github.com/erniealice/espyna-golang/reference"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/template"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"

	templateview "github.com/erniealice/hybra-golang/views/template"

	centymo "github.com/erniealice/centymo-golang"
	collectionmod "github.com/erniealice/centymo-golang/views/collection"
	disbursementmod "github.com/erniealice/centymo-golang/views/disbursement"
	expendituremod "github.com/erniealice/centymo-golang/views/expenditure"
	inventorydomain "github.com/erniealice/centymo-golang/domain/inventory"
	inventorymod "github.com/erniealice/centymo-golang/domain/inventory/views/inventory"
	revenuedomain "github.com/erniealice/centymo-golang/domain/revenue"
	revenuemod "github.com/erniealice/centymo-golang/domain/revenue/views/revenue"
	resourcemod "github.com/erniealice/centymo-golang/views/resource"
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
		cfg.accruedExpense || cfg.accruedExpenseSettlement ||
		cfg.revenueRun ||
		cfg.costSchedule || cfg.supplierPlan || cfg.costPlan ||
		cfg.supplierProductPlan || cfg.supplierProductCostPlan || cfg.supplierSubscription ||
		cfg.treasuryAdvances
	cfg.enableAll = !moduleSelected

	return func(ctx *pyeza.AppContext) error {
		// --- Type-assert translations ---
		translations, ok := ctx.Translations.(*lynguaV1.TranslationProvider)
		if !ok || translations == nil {
			return fmt.Errorf("centymo.Block: ctx.Translations must be *lynguaV1.TranslationProvider")
		}

		// --- Validate use cases ---
		if cfg.useCases == nil {
			return fmt.Errorf("centymo.Block: WithUseCases(...) was not supplied")
		}
		if err := cfg.useCases.RequireFor(cfg); err != nil {
			return err
		}
		useCases := cfg.useCases

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
		downloadFile, _ := ctx.DownloadFile.(func(context.Context, string, string) ([]byte, error))
		listAttachments, _ := ctx.ListAttachments.(func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error))
		createAttachment, _ := ctx.CreateAttachment.(func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error))
		readAttachment, _ := ctx.ReadAttachment.(func(context.Context, *attachmentpb.ReadAttachmentRequest) (*attachmentpb.ReadAttachmentResponse, error))
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
		centymoTableLabels := pyeza.MapTableLabels(ctx.Common)

		// --- Load routes (defaults + optional lyngua overrides) ---
		inventoryRoutes := inventorydomain.DefaultInventoryRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "inventory", &inventoryRoutes)

		revenueRoutes := revenuedomain.DefaultRevenueRoutes()
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
		var inventoryLabels inventorydomain.InventoryLabels
		// Wave 4.2 — wire the standalone inventory child JSONs into the
		// InventoryLabels sub-fields. These files live ONLY in the general/
		// tier (there is no general/inventory.json), so they supply the
		// baseline serial / transaction / depreciation labels for the general
		// tier and the 6 fallback verticals. They are loaded FIRST so the
		// subsequent inventory.json overlay (professional / retail / service,
		// which embed their own serial/transaction/depreciation subtrees) wins
		// where present, while the general tier keeps these baseline values.
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "inventory_serial.json", "inventory_serial", &inventoryLabels.Serial)
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "inventory_transaction.json", "inventory_transaction", &inventoryLabels.Transaction)
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "inventory_depreciation.json", "inventory_depreciation", &inventoryLabels.Depreciation)
		if err := translations.LoadPath("en", ctx.BusinessType, "inventory.json", "inventory", &inventoryLabels); err != nil {
			log.Printf("centymo.Block: warning loading inventory labels: %v", err)
		}

		var revenueLabels revenuedomain.RevenueLabels
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
		// Wave 4.2 — wire the standalone expenditure_category.json into
		// ExpenditureLabels.Category. expenditure.json carries no `category`
		// subtree, so the expense-category list/form/action views
		// (deps.Labels.Category) render blank without this overlay.
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "expenditure_category.json", "expenditure_category", &expenditureLabels.Category)

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

		// Phase 4 — revenue-run (Surface D).
		revenueRunRoutes := revenuedomain.DefaultRevenueRunRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "revenue_run", &revenueRunRoutes)
		revenueRunLabels := revenuedomain.DefaultRevenueRunLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "revenue.json", "revenueRun", &revenueRunLabels)

		// 20260517-expense-run Plan A Phase 4 — Expense Recognition Run (Surfaces B + D).
		expenseRecognitionRunRoutes := centymo.DefaultExpenseRecognitionRunRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "expense_recognition_run", &expenseRecognitionRunRoutes)
		expenseRecognitionRunLabels := centymo.DefaultExpenseRecognitionRunLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "expense_recognition_run.json", "expenseRecognitionRun", &expenseRecognitionRunLabels)

		// 20260517-advance-cash-events Plan B Phase 3 — Advances Dashboard.
		// Routes flow from the new TreasuryAdvancesRoutes block; labels load
		// from advances_dashboard.json + advance_kind.json (the latter is
		// shared with the per-tab AdvanceKind/AdvanceStatus badge rendering
		// done in the existing collection/disbursement detail pages).
		advancesDashboardRoutes := centymo.DefaultTreasuryAdvancesRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "treasury_advances", &advancesDashboardRoutes)
		advancesDashboardLabels := centymo.DefaultAdvancesDashboardLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "advances_dashboard.json", "advancesDashboard", &advancesDashboardLabels)
		advanceEnumLabels := centymo.DefaultAdvanceEnumLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "advance_kind.json", "advanceKind.labels", &advanceEnumLabels)

		// P3 (20260506-supplier-subscriptions) — Routes + Labels for the six new procurement modules.
		costScheduleRoutes := centymo.DefaultCostScheduleRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "cost_schedule", &costScheduleRoutes)
		costScheduleLabels := centymo.DefaultCostScheduleLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "cost_schedule.json", "costSchedule", &costScheduleLabels)

		supplierPlanRoutes := centymo.DefaultSupplierPlanRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "supplier_plan", &supplierPlanRoutes)
		supplierPlanLabels := centymo.DefaultSupplierPlanLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "supplier_plan.json", "supplierPlan", &supplierPlanLabels)

		costPlanRoutes := centymo.DefaultCostPlanRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "cost_plan", &costPlanRoutes)
		costPlanLabels := centymo.DefaultCostPlanLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "cost_plan.json", "costPlan", &costPlanLabels)

		supplierProductPlanRoutes := centymo.DefaultSupplierProductPlanRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "supplier_product_plan", &supplierProductPlanRoutes)
		supplierProductPlanLabels := centymo.DefaultSupplierProductPlanLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "supplier_product_plan.json", "supplierProductPlan", &supplierProductPlanLabels)

		supplierProductCostPlanLabels := centymo.DefaultSupplierProductCostPlanLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "supplier_product_cost_plan.json", "supplierProductCostPlan", &supplierProductCostPlanLabels)

		supplierSubscriptionRoutes := centymo.DefaultSupplierSubscriptionRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "supplier_subscription", &supplierSubscriptionRoutes)
		supplierSubscriptionLabels := centymo.DefaultSupplierSubscriptionLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "supplier_subscription.json", "supplierSubscription", &supplierSubscriptionLabels)

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
			invDeps.ListInventoryItems = useCases.Inventory.ListInventoryItems
			invDeps.CreateInventoryItem = useCases.Inventory.CreateInventoryItem
			invDeps.ReadInventoryItem = useCases.Inventory.ReadInventoryItem
			invDeps.UpdateInventoryItem = useCases.Inventory.UpdateInventoryItem
			invDeps.DeleteInventoryItem = useCases.Inventory.DeleteInventoryItem
			invDeps.ListInventorySerials = useCases.Inventory.ListInventorySerials
			invDeps.CreateInventorySerial = useCases.Inventory.CreateInventorySerial
			invDeps.ReadInventorySerial = useCases.Inventory.ReadInventorySerial
			invDeps.UpdateInventorySerial = useCases.Inventory.UpdateInventorySerial
			invDeps.DeleteInventorySerial = useCases.Inventory.DeleteInventorySerial
			invDeps.ListInventoryTransactions = useCases.Inventory.ListInventoryTransactions
			invDeps.CreateInventoryTransaction = useCases.Inventory.CreateInventoryTransaction
			invDeps.GetInventoryMovementsListPageData = useCases.Inventory.GetInventoryMovementsListPageData
			invDeps.ListInventoryDepreciations = useCases.Inventory.ListInventoryDepreciations
			invDeps.CreateInventoryDepreciation = useCases.Inventory.CreateInventoryDepreciation
			invDeps.ReadInventoryDepreciation = useCases.Inventory.ReadInventoryDepreciation
			invDeps.UpdateInventoryDepreciation = useCases.Inventory.UpdateInventoryDepreciation
			// Cross-domain: product options + locations for inventory item views
			invDeps.ReadProduct = useCases.Product.ReadProduct
			invDeps.ListProductVariantOptions = useCases.Product.ListProductVariantOptions
			invDeps.ListProductOptionValues = useCases.Product.ListProductOptionValues
			invDeps.ListProductOptions = useCases.Product.ListProductOptions
			invDeps.ListLocations = useCases.Entity.Location.ListLocations

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
			revDeps.ListClients = useCases.Entity.Client.ListClients
			revDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName
			// Subscription search + auto-populate
			revDeps.ListSubscriptions = useCases.Subscription.ListSubscriptions
			revDeps.ReadSubscription = useCases.Subscription.ReadSubscription
			revDeps.ReadPricePlan = useCases.PricePlan.ReadPricePlan
			revDeps.ListProductPricePlans = useCases.PricePlan.ListProductPricePlans
			revDeps.ReadProduct = useCases.Product.ReadProduct
			revDeps.ListProducts = useCases.Product.ListProducts
			// Revenue CRUD + list page data
			revDeps.GetListPageData = useCases.Revenue.GetListPageData
			revDeps.CreateRevenue = useCases.Revenue.CreateRevenue
			revDeps.ReadRevenue = useCases.Revenue.ReadRevenue
			revDeps.UpdateRevenue = useCases.Revenue.UpdateRevenue
			revDeps.DeleteRevenue = useCases.Revenue.DeleteRevenue
			// Revenue Line Item CRUD
			revDeps.CreateRevenueLineItem = useCases.Revenue.CreateRevenueLineItem
			revDeps.ReadRevenueLineItem = useCases.Revenue.ReadRevenueLineItem
			revDeps.UpdateRevenueLineItem = useCases.Revenue.UpdateRevenueLineItem
			revDeps.DeleteRevenueLineItem = useCases.Revenue.DeleteRevenueLineItem
			revDeps.ListRevenueLineItems = useCases.Revenue.ListRevenueLineItems
			// Inventory (for stock deduction on status change)
			revDeps.ReadInventoryItem = useCases.Inventory.ReadInventoryItem
			revDeps.UpdateInventoryItem = useCases.Inventory.UpdateInventoryItem
			revDeps.ListInventoryItems = useCases.Inventory.ListInventoryItems
			revDeps.UpdateInventorySerial = useCases.Inventory.UpdateInventorySerial
			revDeps.CreateInventorySerialHistory = useCases.Inventory.CreateInventorySerialHistory
			// Price lookup for line item (find applicable price list + price product)
			revDeps.FindApplicablePriceList = useCases.Product.FindApplicablePriceList
			revDeps.ListPriceProducts = useCases.Product.ListPriceProducts
			// Job activity lookup for "from_activities" revenue type
			revDeps.ReadJobActivity = useCases.Operation.JobActivity.ReadJobActivity
			// Recognize-revenue use case — shared with the subscription recognize-drawer flow
			revDeps.RecognizeRevenueFromSubscription = useCases.Revenue.RecognizeRevenueFromSubscription
			// Phase 5: wire tax lines read-access
			revDeps.ListRevenueTaxLines = useCases.Revenue.ListRevenueTaxLines

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
			// Tax recompute — 501 stub until Phase 4 wires ComputeTaxesForRevenue
			handleFunc(ctx.Routes, "POST", revenueRoutes.RecomputeTaxesURL, revenueMod.RecomputeTaxes)
		}

		// See product.go for wireProductModules (Product 3-mount + ProductLine 2-mount).
		wireProductModules(ctx, cfg, useCases, productWiring{
			db:                         db,
			refChecker:                 refChecker,
			uploadImage:                uploadImage,
			uploadFile:                 uploadFile,
			downloadFile:               downloadFile,
			listAttachments:            listAttachments,
			createAttachment:           createAttachment,
			readAttachment:             readAttachment,
			deleteAttachment:           deleteAttachment,
			newAttachmentID:            newAttachmentID,
			productRoutes:              productRoutes,
			productInventoryRoutes:     productInventoryRoutes,
			productSuppliesRoutes:      productSuppliesRoutes,
			productLineRoutes:          productLineRoutes,
			productLineInventoryRoutes: productLineInventoryRoutes,
			productLabels:              productLabels,
			productInventoryLabels:     productInventoryLabels,
			productSuppliesLabels:      productSuppliesLabels,
			productLineLabels:          productLineLabels,
			centymoTableLabels:         centymoTableLabels,
		})

		// See plan.go for wirePlanModules (PricePlan, PriceSchedule, PriceList, Plan + PlanBundle).
		wirePlanModules(ctx, cfg, useCases, planWiring{
			db:                           db,
			refChecker:                   refChecker,
			uploadFile:                   uploadFile,
			downloadFile:                 downloadFile,
			readAttachment:               readAttachment,
			listAttachments:              listAttachments,
			createAttachment:             createAttachment,
			deleteAttachment:             deleteAttachment,
			newAttachmentID:              newAttachmentID,
			pricePlanRoutes:              pricePlanRoutes,
			priceScheduleRoutes:          priceScheduleRoutes,
			priceScheduleInventoryRoutes: priceScheduleInventoryRoutes,
			priceListRoutes:              priceListRoutes,
			planRoutes:                   planRoutes,
			planBundleRoutes:             planBundleRoutes,
			subscriptionRoutes:           subscriptionRoutes,
			pricePlanLabels:              pricePlanLabels,
			productPricePlanLabels:       productPricePlanLabels,
			priceScheduleLabels:          priceScheduleLabels,
			priceListLabels:              priceListLabels,
			planLabels:                   planLabels,
			centymoTableLabels:           centymoTableLabels,
		})

		// =====================================================================
		// Subscription (inline — not a module, uses subscriptionlist/subscriptionaction/subscriptiondetail directly)
		// See subscription.go for wireSubscriptionModule.
		// =====================================================================

		wireSubscriptionModule(ctx, cfg, useCases, subscriptionWiring{
			db:                  db,
			refChecker:          refChecker,
			uploadFile:          uploadFile,
			downloadFile:        downloadFile,
			readAttachment:      readAttachment,
			listAttachments:     listAttachments,
			createAttachment:    createAttachment,
			deleteAttachment:    deleteAttachment,
			newAttachmentID:     newAttachmentID,
			subscriptionRoutes:  subscriptionRoutes,
			priceScheduleRoutes: priceScheduleRoutes,
			subscriptionLabels:  subscriptionLabels,
			priceScheduleLabels: priceScheduleLabels,
			centymoTableLabels:  centymoTableLabels,
		})

		// =====================================================================
		// Collection module (conditional: only when treasury collection use cases are available)
		// =====================================================================

		if cfg.wantCollection() {
			// 20260517-advance-cash-events Plan B Phase 4 — Advance Schedule
			// labels + the UNSCHEDULED workflow closures. Loaded from lyngua
			// once per Block() invocation and reused here + on the disbursement
			// side below.
			collectionAdvanceLabels := centymo.DefaultTreasuryCollectionAdvanceLabels()
			_ = translations.LoadPathIfExists("en", ctx.BusinessType, "treasury_collection.json", "advance", &collectionAdvanceLabels)

			collDeps := &collectionmod.ModuleDeps{
				Routes:           collectionRoutes,
				Labels:           collectionLabels,
				CommonLabels:     ctx.Common,
				TableLabels:      centymoTableLabels,
				CreateCollection: useCases.Collection.CreateCollection,
				ReadCollection:   useCases.Collection.ReadCollection,
				UpdateCollection: useCases.Collection.UpdateCollection,
				DeleteCollection: useCases.Collection.DeleteCollection,
				ListCollections:  useCases.Collection.ListCollections,
				// Attachments
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
				NewID:            newAttachmentID,
				// 20260517-advance-cash-events Plan B Phase 4 — UNSCHEDULED workflow.
				AdvanceLabels:            collectionAdvanceLabels,
				AdvanceEnumLabels:        advanceEnumLabels,
				SettleUnscheduledAdvance: bridgeSettleAdvance(useCases.Collection.SettleUnscheduledAdvance),
				RefundUnscheduledAdvance: bridgeRefundAdvance(useCases.Collection.RefundUnscheduledAdvance),
				CancelAdvance:            bridgeCancelAdvance(useCases.Collection.CancelAdvance),
			}
			wireCashDashboard(collDeps, useCases)
			collDeps.GetFunctionalCurrency = func(fctx context.Context) string {
				return getFunctionalCurrency(fctx, useCases)
			}
			collectionmod.NewModule(collDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Disbursement module (conditional: only when treasury disbursement use cases are available)
		// =====================================================================

		if cfg.wantDisbursement() {
			// 20260517-advance-cash-events Plan B Phase 4.
			disbursementAdvanceLabels := centymo.DefaultTreasuryDisbursementAdvanceLabels()
			_ = translations.LoadPathIfExists("en", ctx.BusinessType, "treasury_disbursement.json", "advance", &disbursementAdvanceLabels)

			disbursementmod.NewModule(&disbursementmod.ModuleDeps{
				Routes:             disbursementRoutes,
				Labels:             disbursementLabels,
				CommonLabels:       ctx.Common,
				TableLabels:        centymoTableLabels,
				CreateDisbursement: useCases.Disbursement.CreateDisbursement,
				ReadDisbursement:   useCases.Disbursement.ReadDisbursement,
				UpdateDisbursement: useCases.Disbursement.UpdateDisbursement,
				DeleteDisbursement: useCases.Disbursement.DeleteDisbursement,
				ListDisbursements:  useCases.Disbursement.ListDisbursements,
				// Attachments
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
				NewID:            newAttachmentID,
				// 20260517-advance-cash-events Plan B Phase 4 — UNSCHEDULED workflow.
				AdvanceLabels:            disbursementAdvanceLabels,
				AdvanceEnumLabels:        advanceEnumLabels,
				SettleUnscheduledAdvance: bridgeSettleAdvance(useCases.Disbursement.SettleUnscheduledAdvance),
				RefundUnscheduledAdvance: bridgeRefundAdvance(useCases.Disbursement.RefundUnscheduledAdvance),
				CancelAdvance:            bridgeCancelAdvance(useCases.Disbursement.CancelAdvance),
			}).RegisterRoutes(ctx.Routes)
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
			expDeps.ListExpenditures = useCases.Expenditure.ListExpenditures
			expDeps.CreateExpenditure = useCases.Expenditure.CreateExpenditure
			expDeps.ReadExpenditure = useCases.Expenditure.ReadExpenditure
			expDeps.UpdateExpenditure = useCases.Expenditure.UpdateExpenditure
			expDeps.DeleteExpenditure = useCases.Expenditure.DeleteExpenditure
			expDeps.ListExpenditureCategories = useCases.Expenditure.ListExpenditureCategories
			expDeps.CreateExpenditureCategory = useCases.Expenditure.CreateExpenditureCategory
			expDeps.ReadExpenditureCategory = useCases.Expenditure.ReadExpenditureCategory
			expDeps.UpdateExpenditureCategory = useCases.Expenditure.UpdateExpenditureCategory
			expDeps.DeleteExpenditureCategory = useCases.Expenditure.DeleteExpenditureCategory
			expDeps.CreateExpenditureLineItem = useCases.Expenditure.CreateExpenditureLineItem
			expDeps.ReadExpenditureLineItem = useCases.Expenditure.ReadExpenditureLineItem
			expDeps.UpdateExpenditureLineItem = useCases.Expenditure.UpdateExpenditureLineItem
			expDeps.DeleteExpenditureLineItem = useCases.Expenditure.DeleteExpenditureLineItem
			expDeps.ListExpenditureLineItems = useCases.Expenditure.ListExpenditureLineItems
			expDeps.ListSuppliers = useCases.Entity.Supplier.ListSuppliers
			if useCases.Disbursement.CreateDisbursement != nil {
				expDeps.DisbursementRoutes = disbursementRoutes
				expDeps.DisbursementLabels = disbursementLabels
				expDeps.CreateDisbursement = useCases.Disbursement.CreateDisbursement
			}
			// SPS Wave 4 — Recognition + Accrual tabs on the expense detail page.
			// Nil-safe — when the use case is missing, the tab renders an empty state.
			expDeps.ReadExpenseRecognition = useCases.Expenditure.ReadExpenseRecognition
			expDeps.ListAccruedExpenses = useCases.Expenditure.ListAccruedExpenses
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
			expDeps.GetFunctionalCurrency = func(fctx context.Context) string {
				return getFunctionalCurrency(fctx, useCases)
			}
			expDeps.ListAttachments = listAttachments
			expDeps.CreateAttachment = createAttachment
			expDeps.DeleteAttachment = deleteAttachment
			expDeps.NewAttachmentID = newAttachmentID
			expendituremod.NewModule(expDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Resource module
		// =====================================================================

		if cfg.wantResource() {
			resourcemod.NewModule(&resourcemod.ModuleDeps{
				Routes:         resourceRoutes,
				Labels:         resourceLabels,
				CommonLabels:   ctx.Common,
				TableLabels:    centymoTableLabels,
				ListResources:  useCases.Product.ListResources,
				ReadResource:   useCases.Product.ReadResource,
				CreateResource: useCases.Product.CreateResource,
				UpdateResource: useCases.Product.UpdateResource,
				DeleteResource: useCases.Product.DeleteResource,
			}).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// 20260427-supplier-commitments — five new modules (P3a + P3b)
		// See supplier_commitment.go for wireSupplierCommitmentModules.
		// =====================================================================

		wireSupplierCommitmentModules(ctx, cfg, useCases, supplierCommitmentWiring{
			supplierContractRoutes:              supplierContractRoutes,
			supplierContractLabels:              supplierContractLabels,
			supplierContractPriceScheduleRoutes: supplierContractPriceScheduleRoutes,
			procurementRequestRoutes:            procurementRequestRoutes,
			procurementRequestLabels:            procurementRequestLabels,
			procurementRoutes:                   procurementRoutes,
			procurementLabels:                   procurementLabels,
			centymoTableLabels:                  centymoTableLabels,
			uploadFile:                          uploadFile,
			listAttachments:                     listAttachments,
			createAttachment:                    createAttachment,
			deleteAttachment:                    deleteAttachment,
			newAttachmentID:                     newAttachmentID,
		})

		// =====================================================================
		// SPS Wave 4 — supplier-side pricing graph + accrual layer.
		// See supplier_contract_price_schedule.go for wireSupplierContractPriceScheduleModules.
		// =====================================================================

		wireSupplierContractPriceScheduleModules(ctx, cfg, useCases, supplierContractPriceScheduleWiring{
			supplierContractPriceScheduleRoutes: supplierContractPriceScheduleRoutes,
			supplierContractPriceScheduleLabels: supplierContractPriceScheduleLabels,
			centymoTableLabels:                  centymoTableLabels,
			uploadFile:                          uploadFile,
			listAttachments:                     listAttachments,
			createAttachment:                    createAttachment,
			deleteAttachment:                    deleteAttachment,
			newAttachmentID:                     newAttachmentID,
		})

		// ExpenseRecognition module — no Add/Edit drawer (created BY use case).
		// See expense_recognition.go for wireExpenseRecognitionModules.
		wireExpenseRecognitionModules(ctx, cfg, useCases, expenseRecognitionWiring{
			expenseRecognitionRoutes: expenseRecognitionRoutes,
			expenseRecognitionLabels: expenseRecognitionLabels,
			centymoTableLabels:       centymoTableLabels,
			uploadFile:               uploadFile,
			listAttachments:          listAttachments,
			createAttachment:         createAttachment,
			deleteAttachment:         deleteAttachment,
			newAttachmentID:          newAttachmentID,
		})

		// See accrued_expense.go for wireAccruedExpenseModules.
		wireAccruedExpenseModules(ctx, cfg, useCases, accruedExpenseWiring{
			accruedExpenseRoutes: accruedExpenseRoutes,
			accruedExpenseLabels: accruedExpenseLabels,
			centymoTableLabels:   centymoTableLabels,
			uploadFile:           uploadFile,
			listAttachments:      listAttachments,
			createAttachment:     createAttachment,
			deleteAttachment:     deleteAttachment,
			newAttachmentID:      newAttachmentID,
		})

		// =====================================================================
		// Revenue Run module — Surface D (history list + detail pages)
		// Phase 4 of the 20260506-subscription-invoice-run plan.
		// See revenue_run.go for wireRevenueRunModule + proto-shim helpers.
		// =====================================================================

		if cfg.wantRevenueRun() {
			wireRevenueRunModule(ctx, cfg, useCases, revenueRunWiring{
				revenueRunRoutes:   revenueRunRoutes,
				revenueRunLabels:   revenueRunLabels,
				revenueRoutes:      revenueRoutes,
				centymoTableLabels: centymoTableLabels,
				uploadFile:         uploadFile,
				listAttachments:    listAttachments,
				createAttachment:   createAttachment,
				deleteAttachment:   deleteAttachment,
				newAttachmentID:    newAttachmentID,
			})
		}

		// =====================================================================
		// 20260517-expense-run Plan A Phase 4 — Expense Recognition Run.
		// Surfaces B (queue) + D (list/detail). See expense_recognition_run.go
		// for wireExpenseRecognitionRunModule + proto-shim helpers.
		// =====================================================================

		if cfg.wantExpenseRecognitionRun() {
			wireExpenseRecognitionRunModule(ctx, cfg, useCases, expenseRecognitionRunWiring{
				routes:                                 expenseRecognitionRunRoutes,
				labels:                                 expenseRecognitionRunLabels,
				expenditureRoutes:                      expenditureRoutes,
				centymoTableLabels:                     centymoTableLabels,
				supplierDetailURL:                      cfg.supplierDetailURL,
				supplierExpenseRecognitionRunDrawerURL: cfg.supplierExpenseRecognitionRunDrawerURL,
			})
		}

		// =====================================================================
		// 20260517-advance-cash-events Plan B Phase 3 — Advances Dashboard.
		// See advances_dashboard.go for wireAdvancesDashboardModule.
		// =====================================================================

		if cfg.wantTreasuryAdvances() {
			wireAdvancesDashboardModule(ctx, cfg, useCases, advancesDashboardWiring{
				routes:             advancesDashboardRoutes,
				labels:             advancesDashboardLabels,
				enumLabels:         advanceEnumLabels,
				collectionRoutes:   collectionRoutes,
				disbursementRoutes: disbursementRoutes,
				centymoTableLabels: centymoTableLabels,
				functionalCurrency: func(fctx context.Context) string {
					return getFunctionalCurrency(fctx, useCases)
				},
			})

			// 20260517-advance-cash-events Plan B Phase 7 — buying-side
			// SupplierBillingEvent surfaces (list, detail, Recognize).
			// See supplier_billing_event.go for wireSupplierBillingEventModule.
			wireSupplierBillingEventModule(ctx, useCases, supplierBillingEventWiring{
				routes: advancesDashboardRoutes,
			})
		}

		// See supplier_subscription.go for wireSupplierSubscriptionModules.
		wireSupplierSubscriptionModules(ctx, cfg, useCases, supplierSubscriptionWiring{
			db:                            db,
			costScheduleRoutes:            costScheduleRoutes,
			costScheduleLabels:            costScheduleLabels,
			supplierPlanRoutes:            supplierPlanRoutes,
			supplierPlanLabels:            supplierPlanLabels,
			costPlanRoutes:                costPlanRoutes,
			costPlanLabels:                costPlanLabels,
			supplierProductPlanRoutes:     supplierProductPlanRoutes,
			supplierProductPlanLabels:     supplierProductPlanLabels,
			supplierProductCostPlanLabels: supplierProductCostPlanLabels,
			supplierSubscriptionRoutes:    supplierSubscriptionRoutes,
			supplierSubscriptionLabels:    supplierSubscriptionLabels,
			expenseRecognitionRunLabels:   expenseRecognitionRunLabels,
			centymoTableLabels:            centymoTableLabels,
		})

		log.Println("  centymo commerce domain initialized")
		return nil
	}
}

// ---------------------------------------------------------------------------
// Workspace currency helpers (mirrors fycha-golang/block/block.go convention)
// ---------------------------------------------------------------------------

// getDefaultWorkspaceID reads the DEFAULT_WORKSPACE_ID env var and falls back
// to "default-workspace". Matches the entydad block convention.
func getDefaultWorkspaceID() string {
	if v := os.Getenv("DEFAULT_WORKSPACE_ID"); v != "" {
		return v
	}
	return "default-workspace"
}

// getFunctionalCurrency returns the workspace's functional_currency (ISO 4217)
// for use in money display strings. Returns empty string when the workspace use
// case is not wired or the read fails — types.FormatMoney handles empty
// currency by omitting the prefix, so the worst-case fallback is the bare
// number rather than a hardcoded peso glyph.
func getFunctionalCurrency(ctx context.Context, useCases *UseCases) string {
	if useCases == nil || useCases.Entity.Workspace.ReadWorkspace == nil {
		return ""
	}
	resp, err := useCases.Entity.Workspace.ReadWorkspace(ctx, &workspacepb.ReadWorkspaceRequest{
		Data: &workspacepb.Workspace{Id: getDefaultWorkspaceID()},
	})
	if err != nil || resp == nil {
		return ""
	}
	data := resp.GetData()
	if len(data) == 0 {
		return ""
	}
	return data[0].GetFunctionalCurrency()
}
