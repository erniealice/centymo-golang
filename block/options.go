// Package block — option types and toggles for centymo.Block.
//
// This file owns the public BlockOption surface and the package-private
// blockConfig flags. Adding a new optional centymo module means: (1) add a
// bool field to blockConfig, (2) add a `WithX() BlockOption` here, (3) add
// a `wantX()` accessor here, (4) check `cfg.wantX()` in block.go and wire
// the module. Nothing else in this file is load-bearing — it is a flat list
// by design so a reader can scan every option in one screen.
package block

// ---------------------------------------------------------------------------
// BlockOption — per-module granular selection
// ---------------------------------------------------------------------------

// BlockOption enables specific centymo sub-modules within Block().
type BlockOption func(*blockConfig)

type blockConfig struct {
	enableAll    bool
	useCases     *UseCases
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
	// Phase 4 (20260506-subscription-invoice-run) — revenue-run history pages.
	revenueRun bool
	// P3 (20260506-supplier-subscriptions) — six new procurement modules.
	costSchedule            bool
	supplierPlan            bool
	costPlan                bool
	supplierProductPlan     bool
	supplierProductCostPlan bool
	supplierSubscription    bool
	// clientDetailURL is the absolute path template (e.g.
	// "/app/clients/detail/{id}") used for the subscription detail's
	// page-header breadcrumb when accessed via the under-client nested route.
	// Centymo cannot import entydad (wrong dep direction); the consumer
	// supplies it via WithClientDetailURL.
	clientDetailURL string
	// clientRevenueRunDrawerURL is the path template for the Surface-A per-client
	// revenue-run drawer (e.g. "/action/client/revenue-run/{id}").
	// Phase 7 (Surface B) — the queue page drills into this drawer per row.
	clientRevenueRunDrawerURL string
	// jobDetailURL is the absolute path template (e.g. "/app/jobs/detail/{id}")
	// used by the subscription detail's Operations tab to deep-link to fayna
	// Job detail. Centymo cannot import fayna; the consumer supplies it via
	// WithJobDetailURL. Optional — Operations tab renders job rows without a
	// link when unset.
	jobDetailURL string
}

// WithUseCases supplies the typed use-case aggregate for centymo.Block.
// It must be provided before any module options. Block() will fail at startup
// with a descriptive error if required fields are nil for the enabled modules.
func WithUseCases(uc *UseCases) BlockOption {
	return func(c *blockConfig) { c.useCases = uc }
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

// WithRevenueRun enables the revenue-run history list + detail pages (Surface D).
// Phase 4 of the 20260506-subscription-invoice-run plan.
func WithRevenueRun() BlockOption { return func(c *blockConfig) { c.revenueRun = true } }

// P3 (20260506-supplier-subscriptions) — six new procurement module toggles.
func WithCostSchedule() BlockOption {
	return func(c *blockConfig) { c.costSchedule = true }
}
func WithSupplierPlan() BlockOption {
	return func(c *blockConfig) { c.supplierPlan = true }
}
func WithCostPlan() BlockOption {
	return func(c *blockConfig) { c.costPlan = true }
}
func WithSupplierProductPlan() BlockOption {
	return func(c *blockConfig) { c.supplierProductPlan = true }
}
func WithSupplierProductCostPlan() BlockOption {
	return func(c *blockConfig) { c.supplierProductCostPlan = true }
}
func WithSupplierSubscription() BlockOption {
	return func(c *blockConfig) { c.supplierSubscription = true }
}

// WithClientDetailURL supplies the entydad client-detail path template (e.g.
// "/app/clients/detail/{id}") so the subscription detail page can render a
// "client → subscription" breadcrumb when accessed under a client context.
// Optional — when unset the breadcrumb label still renders (sourced from the
// joined client) but isn't a link.
func WithClientDetailURL(url string) BlockOption {
	return func(c *blockConfig) { c.clientDetailURL = url }
}

// WithClientRevenueRunDrawerURL supplies the entydad client-revenue-run drawer
// path template (e.g. "/action/client/revenue-run/{id}") so the queue page
// (Surface B) can render a per-row [Run] action that opens the Surface-A drawer.
// Optional — the per-row action is omitted when unset.
// Phase 7 (20260506-subscription-invoice-run Surface B).
func WithClientRevenueRunDrawerURL(url string) BlockOption {
	return func(c *blockConfig) { c.clientRevenueRunDrawerURL = url }
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

// Phase 4 (20260506-subscription-invoice-run).
func (c *blockConfig) wantRevenueRun() bool { return c.enableAll || c.revenueRun }

// P3 (20260506-supplier-subscriptions) — six new procurement module want() helpers.
func (c *blockConfig) wantCostSchedule() bool {
	return c.enableAll || c.costSchedule
}
func (c *blockConfig) wantSupplierPlan() bool {
	return c.enableAll || c.supplierPlan
}
func (c *blockConfig) wantCostPlan() bool {
	return c.enableAll || c.costPlan
}
func (c *blockConfig) wantSupplierProductPlan() bool {
	return c.enableAll || c.supplierProductPlan
}
func (c *blockConfig) wantSupplierProductCostPlan() bool {
	return c.enableAll || c.supplierProductCostPlan
}
func (c *blockConfig) wantSupplierSubscription() bool {
	return c.enableAll || c.supplierSubscription
}
