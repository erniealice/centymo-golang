// Package block — composition-v2 catalog.
//
// Each XxxUnit function returns a compose.Unit whose Mount closure mirrors the
// wiring that block.go's wireXxx helpers perform today. Callers build the list
// via AllUnits and hand it to the compose engine.
//
// Design rules (matching fayna/block/catalog.go):
//   - Every XxxUnit func calls Describe() from the entity sub-package to obtain
//     the descriptor (routes + labels pointers, lyngua JSON bindings, TemplatesFS).
//   - The Mount closure type-asserts the post-overlay Routes/Labels pointers back
//     to the concrete entity types (safe: the descriptor owns the pointer).
//   - Wiring mirrors block.go's existing wireXxx helpers exactly — same deps
//     construction order, same use-case field assignments, same nil-guards.
//   - compose.HandleFunc is used for raw http.HandlerFunc routes.
package block

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/compose"

	expendituredomain "github.com/erniealice/centymo-golang/domain/expenditure"
	accruedexpensepkg "github.com/erniealice/centymo-golang/domain/expenditure/accrued_expense"
	expenditurepkg "github.com/erniealice/centymo-golang/domain/expenditure/expenditure"
	expenserecognitionpkg "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition"
	expenserecognitionrunpkg "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run"
	procurementrequestpkg "github.com/erniealice/centymo-golang/domain/expenditure/procurement_request"
	supplierbillingeventpkg "github.com/erniealice/centymo-golang/domain/expenditure/supplier_billing_event"
	suppliercontractpkg "github.com/erniealice/centymo-golang/domain/expenditure/supplier_contract"
	suppliercontractpriceschedulepkg "github.com/erniealice/centymo-golang/domain/expenditure/supplier_contract_price_schedule"
	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
	inventorydomain "github.com/erniealice/centymo-golang/domain/inventory"
	inventorypkg "github.com/erniealice/centymo-golang/domain/inventory/inventory"
	procurementdomain "github.com/erniealice/centymo-golang/domain/procurement"
	costplanpkg "github.com/erniealice/centymo-golang/domain/procurement/cost_plan"
	costschedulepkg "github.com/erniealice/centymo-golang/domain/procurement/cost_schedule"
	procurementdashboardpkg "github.com/erniealice/centymo-golang/domain/procurement/procurementdashboard"
	supplierplanpkg "github.com/erniealice/centymo-golang/domain/procurement/supplier_plan"
	supplierproductplanpkg "github.com/erniealice/centymo-golang/domain/procurement/supplier_product_plan"
	suppliersubscriptionpkg "github.com/erniealice/centymo-golang/domain/procurement/supplier_subscription"
	productdom "github.com/erniealice/centymo-golang/domain/product"
	pricelistpkg "github.com/erniealice/centymo-golang/domain/product/price_list"
	productpkg "github.com/erniealice/centymo-golang/domain/product/product"
	resourcepkg "github.com/erniealice/centymo-golang/domain/product/resource"
	revenuedomain "github.com/erniealice/centymo-golang/domain/revenue"
	revenuepkg "github.com/erniealice/centymo-golang/domain/revenue/revenue"
	revenuerunpkg "github.com/erniealice/centymo-golang/domain/revenue/revenue_run"
	subscriptiondom "github.com/erniealice/centymo-golang/domain/subscription"
	planpkg "github.com/erniealice/centymo-golang/domain/subscription/plan"
	priceplanpkg "github.com/erniealice/centymo-golang/domain/subscription/price_plan"
	priceschedulepkg "github.com/erniealice/centymo-golang/domain/subscription/price_schedule"
	subscriptionpkg "github.com/erniealice/centymo-golang/domain/subscription/subscription"
	treasurydomain "github.com/erniealice/centymo-golang/domain/treasury"
	collectionpkg "github.com/erniealice/centymo-golang/domain/treasury/collection"
	disbursementpkg "github.com/erniealice/centymo-golang/domain/treasury/disbursement"
	advancesdashboardpkg "github.com/erniealice/centymo-golang/domain/treasury/treasuryadvancesdashboard"
)

// allEnabledConfig returns a blockConfig with all modules enabled, used by
// PlanUnit and SubscriptionUnit when calling the private wire helpers.
func allEnabledConfig() *blockConfig {
	return &blockConfig{enableAll: true}
}

// labelsOrZero returns the concrete label value from a compose.Unit pointer, or
// the zero value of T when the descriptor has no Labels (partial descriptor).
// Usage: labelsOrZero[inventorypkg.Labels](&u)
func labelsOrZero[T any](u *compose.Unit) T {
	if u.Labels != nil {
		if v, ok := u.Labels.(*T); ok {
			return *v
		}
	}
	var zero T
	return zero
}

// ---------------------------------------------------------------------------
// Inventory
// ---------------------------------------------------------------------------

func InventoryUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := inventorypkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*inventorypkg.Routes)

		// Inventory descriptor is partial (no LabelJSON); labels are zero-value
		// until the compose engine gains a full lyngua overlay for inventory.
		labels := labelsOrZero[inventorypkg.Labels](&u)

		deps := &inventorydomain.InventoryModuleDeps{
			Routes:           *r,
			Labels:           labels,
			CommonLabels:     mc.Common,
			TableLabels:      mc.Table,
			SetItemActive:    setActiveClosure(uc, "inventory_item"),
			UploadFile:       infra.UploadFile,
			ListAttachments:  infra.ListAttachments,
			CreateAttachment: infra.CreateAttachment,
			DeleteAttachment: infra.DeleteAttachment,
			NewID:            infra.NewAttachmentID,
		}
		deps.ListInventoryItems = uc.Inventory.ListInventoryItems
		deps.CreateInventoryItem = uc.Inventory.CreateInventoryItem
		deps.ReadInventoryItem = uc.Inventory.ReadInventoryItem
		deps.UpdateInventoryItem = uc.Inventory.UpdateInventoryItem
		deps.DeleteInventoryItem = uc.Inventory.DeleteInventoryItem
		deps.ListInventorySerials = uc.Inventory.ListInventorySerials
		deps.CreateInventorySerial = uc.Inventory.CreateInventorySerial
		deps.ReadInventorySerial = uc.Inventory.ReadInventorySerial
		deps.UpdateInventorySerial = uc.Inventory.UpdateInventorySerial
		deps.DeleteInventorySerial = uc.Inventory.DeleteInventorySerial
		deps.ListInventoryTransactions = uc.Inventory.ListInventoryTransactions
		deps.CreateInventoryTransaction = uc.Inventory.CreateInventoryTransaction
		deps.GetInventoryMovementsListPageData = uc.Inventory.GetInventoryMovementsListPageData
		deps.ListInventoryDepreciations = uc.Inventory.ListInventoryDepreciations
		deps.CreateInventoryDepreciation = uc.Inventory.CreateInventoryDepreciation
		deps.ReadInventoryDepreciation = uc.Inventory.ReadInventoryDepreciation
		deps.UpdateInventoryDepreciation = uc.Inventory.UpdateInventoryDepreciation
		deps.ReadProduct = uc.Product.ReadProduct
		deps.ListProductVariantOptions = uc.Product.ListProductVariantOptions
		deps.ListProductOptionValues = uc.Product.ListProductOptionValues
		deps.ListProductOptions = uc.Product.ListProductOptions
		deps.ListLocations = uc.Entity.Location.ListLocations
		deps.LocationName = buildLocationResolver(uc)

		invMod := inventorydomain.NewInventoryModule(deps)
		invMod.RegisterRoutes(mc.Routes)
		compose.HandleFunc(mc.Routes, "GET", r.MovementsExportURL, invMod.MovementsExport)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Revenue
// ---------------------------------------------------------------------------

func RevenueUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := revenuepkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*revenuepkg.Routes)
		// Revenue descriptor is partial (no LabelJSON); labels are zero-value
		// until a DefaultLabels factory is added to the revenue entity package.
		labels := labelsOrZero[revenuepkg.Labels](&u)

		deps := &revenuedomain.RevenueModuleDeps{
			Routes:                 *r,
			Labels:                 labels,
			CommonLabels:           mc.Common,
			TableLabels:            mc.Table,
			GenerateDoc:            infra.GenerateDoc,
			ListDocumentTemplates:  infra.ListDocTemplates,
			CreateDocumentTemplate: infra.CreateDocTemplate,
			UpdateDocumentTemplate: infra.UpdateDocTemplate,
			DeleteDocumentTemplate: infra.DeleteDocTemplate,
			UploadTemplate:         infra.UploadTemplate,
			SendEmail:              infra.SendEmail,
			UploadFile:             infra.UploadFile,
			ListAttachments:        infra.ListAttachments,
			CreateAttachment:       infra.CreateAttachment,
			DeleteAttachment:       infra.DeleteAttachment,
			NewID:                  infra.NewAttachmentID,
		}
		wireRevenueDeps(deps, uc)
		revenueMod := revenuedomain.NewRevenueModule(deps)
		revenueMod.RegisterRoutes(mc.Routes)
		compose.HandleFunc(mc.Routes, "GET", r.InvoiceDownloadURL, revenueMod.InvoiceDownload)
		compose.HandleFunc(mc.Routes, "POST", r.SendEmailURL, revenueMod.SendEmailHandler)
		compose.HandleFunc(mc.Routes, "GET", r.SearchClientURL, revenueMod.SearchClients)
		compose.HandleFunc(mc.Routes, "GET", r.SearchSubscriptionURL, revenueMod.SearchSubscriptions)
		compose.HandleFunc(mc.Routes, "GET", r.SearchLocationURL, revenueMod.SearchLocations)
		compose.HandleFunc(mc.Routes, "GET", r.SearchProductURL, revenueMod.SearchProducts)
		compose.HandleFunc(mc.Routes, "GET", r.PriceLookupURL, revenueMod.PriceLookup)
		compose.HandleFunc(mc.Routes, "POST", r.RecomputeTaxesURL, revenueMod.RecomputeTaxes)
		return nil
	}
	return u
}

// wireRevenueDeps assigns the use-case closures onto RevenueModuleDeps.
// Extracted so RevenueUnit stays readable.
func wireRevenueDeps(deps *revenuedomain.RevenueModuleDeps, uc *UseCases) {
	deps.ListClients = uc.Entity.Client.ListClients
	deps.SearchClientsByName = uc.Entity.Client.SearchClientsByName
	deps.ListSubscriptions = uc.Subscription.ListSubscriptions
	deps.ReadSubscription = uc.Subscription.ReadSubscription
	deps.ReadPricePlan = uc.PricePlan.ReadPricePlan
	deps.ListProductPricePlans = uc.PricePlan.ListProductPricePlans
	deps.ReadProduct = uc.Product.ReadProduct
	deps.ListProducts = uc.Product.ListProducts
	deps.GetListPageData = uc.Revenue.GetListPageData
	deps.CreateRevenue = uc.Revenue.CreateRevenue
	deps.ReadRevenue = uc.Revenue.ReadRevenue
	deps.UpdateRevenue = uc.Revenue.UpdateRevenue
	deps.DeleteRevenue = uc.Revenue.DeleteRevenue
	deps.CreateRevenueLineItem = uc.Revenue.CreateRevenueLineItem
	deps.ReadRevenueLineItem = uc.Revenue.ReadRevenueLineItem
	deps.UpdateRevenueLineItem = uc.Revenue.UpdateRevenueLineItem
	deps.DeleteRevenueLineItem = uc.Revenue.DeleteRevenueLineItem
	deps.ListRevenueLineItems = uc.Revenue.ListRevenueLineItems
	deps.ReadInventoryItem = uc.Inventory.ReadInventoryItem
	deps.UpdateInventoryItem = uc.Inventory.UpdateInventoryItem
	deps.ListInventoryItems = uc.Inventory.ListInventoryItems
	deps.UpdateInventorySerial = uc.Inventory.UpdateInventorySerial
	deps.CreateInventorySerialHistory = uc.Inventory.CreateInventorySerialHistory
	deps.FindApplicablePriceList = uc.Product.FindApplicablePriceList
	deps.ListPriceProducts = uc.Product.ListPriceProducts
	deps.ReadJobActivity = uc.Operation.JobActivity.ReadJobActivity
	deps.RecognizeRevenueFromSubscription = uc.Revenue.RecognizeRevenueFromSubscription
	deps.ListRevenueTaxLines = uc.Revenue.ListRevenueTaxLines
	deps.CreateRevenuePayment = uc.Revenue.RevenuePayment.CreateRevenuePayment
	deps.ReadRevenuePayment = uc.Revenue.RevenuePayment.ReadRevenuePayment
	deps.UpdateRevenuePayment = uc.Revenue.RevenuePayment.UpdateRevenuePayment
	deps.DeleteRevenuePayment = uc.Revenue.RevenuePayment.DeleteRevenuePayment
	deps.ListRevenuePayments = uc.Revenue.RevenuePayment.ListRevenuePayments
	deps.ReadCollectionMethod = uc.CollectionMethod.ReadCollectionMethod
	deps.ListCollectionMethods = uc.CollectionMethod.ListCollectionMethods
	deps.ListLocations = uc.Entity.Location.ListLocations
}

// ---------------------------------------------------------------------------
// Product
// ---------------------------------------------------------------------------

func ProductUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := productpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*productpkg.Routes)

		// Product descriptor is partial (no LabelJSON); labels are zero-value
		// until a DefaultLabels factory is added to the product entity package.
		productLabels := labelsOrZero[productpkg.Labels](&u)

		var getProductInUseIDs func(context.Context, []string) (map[string]bool, error)
		if infra.RefChecker != nil {
			getProductInUseIDs = infra.RefChecker.GetProductInUseIDs
		}

		deps := &productdom.ProductModuleDeps{
			Routes:           *r,
			Labels:           productLabels,
			Mode:             "service",
			CommonLabels:     mc.Common,
			TableLabels:      mc.Table,
			GetInUseIDs:      getProductInUseIDs,
			SetProductActive: setActiveClosure(uc, "product"),
			PermissionEntity: "service",
			UploadImage:      infra.UploadImage,
			UploadFile:       infra.UploadFile,
			DownloadFile:     infra.DownloadFile,
			ReadAttachment:   infra.ReadAttachment,
			ListAttachments:  infra.ListAttachments,
			CreateAttachment: infra.CreateAttachment,
			DeleteAttachment: infra.DeleteAttachment,
			NewID:            infra.NewAttachmentID,
		}
		wireProductDeps(deps, uc)
		wireServiceDashboard(deps, uc)
		deps.LocationName = buildLocationResolver(uc)
		productdom.NewProductModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

func wireProductDeps(deps *productdom.ProductModuleDeps, uc *UseCases) {
	deps.ListProducts = uc.Product.ListProducts
	deps.ReadProduct = uc.Product.ReadProduct
	deps.CreateProduct = uc.Product.CreateProduct
	deps.UpdateProduct = uc.Product.UpdateProduct
	deps.DeleteProduct = uc.Product.DeleteProduct
	deps.ListProductVariants = uc.Product.ListProductVariants
	deps.ReadProductVariant = uc.Product.ReadProductVariant
	deps.CreateProductVariant = uc.Product.CreateProductVariant
	deps.UpdateProductVariant = uc.Product.UpdateProductVariant
	deps.DeleteProductVariant = uc.Product.DeleteProductVariant
	deps.ListProductVariantOptions = uc.Product.ListProductVariantOptions
	deps.CreateProductVariantOption = uc.Product.CreateProductVariantOption
	deps.DeleteProductVariantOption = uc.Product.DeleteProductVariantOption
	deps.ListProductOptions = uc.Product.ListProductOptions
	deps.ReadProductOption = uc.Product.ReadProductOption
	deps.CreateProductOption = uc.Product.CreateProductOption
	deps.UpdateProductOption = uc.Product.UpdateProductOption
	deps.DeleteProductOption = uc.Product.DeleteProductOption
	deps.ListProductOptionValues = uc.Product.ListProductOptionValues
	deps.ReadProductOptionValue = uc.Product.ReadProductOptionValue
	deps.CreateProductOptionValue = uc.Product.CreateProductOptionValue
	deps.UpdateProductOptionValue = uc.Product.UpdateProductOptionValue
	deps.DeleteProductOptionValue = uc.Product.DeleteProductOptionValue
	deps.ListProductAttributes = uc.Product.ListProductAttributes
	deps.CreateProductAttribute = uc.Product.CreateProductAttribute
	deps.DeleteProductAttribute = uc.Product.DeleteProductAttribute
	deps.ListProductVariantImages = uc.Product.ListProductVariantImages
	deps.CreateProductVariantImage = uc.Product.CreateProductVariantImage
	deps.DeleteProductVariantImage = uc.Product.DeleteProductVariantImage
	deps.ListProductPlans = uc.Product.ListProductPlans
	deps.ListAttributes = uc.Common.ListAttributes
	deps.ReadAttribute = uc.Common.ReadAttribute
}

// ---------------------------------------------------------------------------
// Resource
// ---------------------------------------------------------------------------

func ResourceUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := resourcepkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*resourcepkg.Routes)
		l := u.Labels.(*resourcepkg.Labels)

		deps := &productdom.ResourceModuleDeps{
			Routes:         *r,
			Labels:         *l,
			CommonLabels:   mc.Common,
			TableLabels:    mc.Table,
			ListResources:  uc.Product.ListResources,
			ReadResource:   uc.Product.ReadResource,
			CreateResource: uc.Product.CreateResource,
			UpdateResource: uc.Product.UpdateResource,
			DeleteResource: uc.Product.DeleteResource,
		}
		productdom.NewResourceModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// PriceList
// ---------------------------------------------------------------------------

func PriceListUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := pricelistpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*pricelistpkg.Routes)
		labels := labelsOrZero[pricelistpkg.Labels](&u)

		deps := &productdom.PriceListModuleDeps{
			Routes:         *r,
			Labels:         labels,
			CommonLabels:   mc.Common,
			TableLabels:    mc.Table,
			ListPriceLists: uc.Product.ListPriceLists,
			ReadPriceList:           uc.Product.ReadPriceList,
			CreatePriceList:         uc.Product.CreatePriceList,
			UpdatePriceList:         uc.Product.UpdatePriceList,
			DeletePriceList:         uc.Product.DeletePriceList,
			ListPriceProducts:       uc.Product.ListPriceProducts,
			CreatePriceProduct:      uc.Product.CreatePriceProduct,
			DeletePriceProduct:      uc.Product.DeletePriceProduct,
		}
		productdom.NewPriceListModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// PricePlan
// ---------------------------------------------------------------------------

func PricePlanUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := priceplanpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*priceplanpkg.Routes)
		l := u.Labels.(*priceplanpkg.Labels)

		productPricePlanLabels := subscriptiondom.DefaultProductPricePlanLabels()
		priceScheduleLabels := priceschedulepkg.DefaultLabels()

		var getPricePlanInUseIDs func(context.Context, []string) (map[string]bool, error)
		if infra.RefChecker != nil {
			getPricePlanInUseIDs = infra.RefChecker.GetPricePlanInUseIDs
		}

		deps := &subscriptiondom.PricePlanModuleDeps{
			Routes:                    *r,
			Labels:                    *l,
			ProductPricePlanLabels:    productPricePlanLabels,
			PriceScheduleDetailLabels: priceScheduleLabels.Detail,
			CommonLabels:              mc.Common,
			TableLabels:               mc.Table,
			GetPricePlanInUseIDs:      getPricePlanInUseIDs,
			ListPricePlans:            uc.PricePlan.ListPricePlans,
			ReadPricePlan:             uc.PricePlan.ReadPricePlan,
			CreatePricePlan:           uc.PricePlan.CreatePricePlan,
			UpdatePricePlan:           uc.PricePlan.UpdatePricePlan,
			DeletePricePlan:           uc.PricePlan.DeletePricePlan,
			ListProductPricePlans:     uc.PricePlan.ListProductPricePlans,
			CreateProductPricePlan:    uc.PricePlan.CreateProductPricePlan,
			UpdateProductPricePlan:    uc.PricePlan.UpdateProductPricePlan,
			DeleteProductPricePlan:    uc.PricePlan.DeleteProductPricePlan,
			ListPlans:                 uc.Plan.ListPlans,
			ListPriceSchedules:        uc.PriceSchedule.ListPriceSchedules,
			ReadPlan:                  uc.Plan.ReadPlan,
			ListProducts:              uc.Product.ListProducts,
			ListProductPlans:          uc.Product.ListProductPlans,
			ListProductVariants:       uc.Product.ListProductVariants,
			ListProductOptions:        uc.Product.ListProductOptions,
			ListProductOptionValues:   uc.Product.ListProductOptionValues,
			ListProductVariantOptions: uc.Product.ListProductVariantOptions,
			UploadFile:                infra.UploadFile,
			ListAttachments:           infra.ListAttachments,
			CreateAttachment:          infra.CreateAttachment,
			DeleteAttachment:          infra.DeleteAttachment,
			NewAttachmentID:           infra.NewAttachmentID,
		}
		subscriptiondom.NewPricePlanModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// PriceSchedule
// ---------------------------------------------------------------------------

func PriceScheduleUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := priceschedulepkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*priceschedulepkg.Routes)
		l := u.Labels.(*priceschedulepkg.Labels)

		pricePlanLabels := priceplanpkg.DefaultLabels()
		productPricePlanLabels := subscriptiondom.DefaultProductPricePlanLabels()

		var getPriceScheduleInUseIDs func(context.Context, []string) (map[string]bool, error)
		var getPricePlanInUseIDs func(context.Context, []string) (map[string]bool, error)
		if infra.RefChecker != nil {
			getPriceScheduleInUseIDs = infra.RefChecker.GetPriceScheduleInUseIDs
			getPricePlanInUseIDs = infra.RefChecker.GetPricePlanInUseIDs
		}

		deps := &subscriptiondom.PriceScheduleModuleDeps{
			Routes:                   *r,
			Labels:                   *l,
			PricePlanLabels:          pricePlanLabels,
			ProductPricePlanLabels:   productPricePlanLabels,
			CommonLabels:             mc.Common,
			TableLabels:              mc.Table,
			GetPriceScheduleInUseIDs: getPriceScheduleInUseIDs,
			GetPricePlanInUseIDs:     getPricePlanInUseIDs,
			ListPriceSchedules:       uc.PriceSchedule.ListPriceSchedules,
			ReadPriceSchedule:        uc.PriceSchedule.ReadPriceSchedule,
			CreatePriceSchedule:      uc.PriceSchedule.CreatePriceSchedule,
			UpdatePriceSchedule:      uc.PriceSchedule.UpdatePriceSchedule,
			DeletePriceSchedule:      uc.PriceSchedule.DeletePriceSchedule,
			ListLocations:            uc.Entity.Location.ListLocations,
			ListClients:              uc.Entity.Client.ListClients,
			ListPlans:                uc.Plan.ListPlans,
			ListPricePlans:           uc.PricePlan.ListPricePlans,
			ReadPricePlan:            uc.PricePlan.ReadPricePlan,
			CreatePricePlan:          uc.PricePlan.CreatePricePlan,
			UpdatePricePlan:          uc.PricePlan.UpdatePricePlan,
			DeletePricePlan:          uc.PricePlan.DeletePricePlan,
			ListProducts:             uc.Product.ListProducts,
			ListProductPlans:         uc.Product.ListProductPlans,
			ListProductVariants:      uc.Product.ListProductVariants,
			ListProductPricePlans:    uc.PricePlan.ListProductPricePlans,
			CreateProductPricePlan:   uc.PricePlan.CreateProductPricePlan,
			UpdateProductPricePlan:   uc.PricePlan.UpdateProductPricePlan,
			DeleteProductPricePlan:   uc.PricePlan.DeleteProductPricePlan,
			ListSubscriptionsByPricePlan: uc.PriceSchedule.ListSubscriptionsByPricePlan,
			UploadFile:               infra.UploadFile,
			ListAttachments:          infra.ListAttachments,
			CreateAttachment:         infra.CreateAttachment,
			DeleteAttachment:         infra.DeleteAttachment,
			NewAttachmentID:          infra.NewAttachmentID,
		}
		subscriptiondom.NewPriceScheduleModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Plan
// ---------------------------------------------------------------------------

// PlanUnit wires the plan + plan-bundle domain by delegating to the existing
// wirePlanModules helper (same helper Block() uses). This avoids duplicating
// the 800-line sub-package registration logic inline; catalog.go stays in
// the same block package so the private helpers are directly accessible.
func PlanUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := planpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*planpkg.Routes)
		l := u.Labels.(*planpkg.Labels)

		pricePlanLabels := priceplanpkg.DefaultLabels()
		productPricePlanLabels := subscriptiondom.DefaultProductPricePlanLabels()
		priceScheduleLabels := priceschedulepkg.DefaultLabels()
		priceScheduleInventoryRoutes := priceschedulepkg.DefaultInventoryRoutes()

		// Look up subscription routes for breadcrumb/subscription-tab URL wiring.
		subRoutes := subscriptionpkg.DefaultRoutes()
		if sr, ok := compose.RoutesOf[*subscriptionpkg.Routes](mc, "subscription.subscription"); ok {
			subRoutes = *sr
		}
		// Look up price-schedule routes if the sibling unit was registered.
		var psRoutes priceschedulepkg.Routes
		if psr, ok := compose.RoutesOf[*priceschedulepkg.Routes](mc, "subscription.price_schedule"); ok {
			psRoutes = *psr
		} else {
			psRoutes = priceschedulepkg.DefaultRoutes()
		}

		minCtx := &pyeza.AppContext{
			Routes: mc.Routes,
			Common: mc.Common,
		}
		wirePlanModules(minCtx, allEnabledConfig(), uc, planWiring{
			refChecker:                   infra.RefChecker,
			uploadFile:                   infra.UploadFile,
			downloadFile:                 infra.DownloadFile,
			readAttachment:               infra.ReadAttachment,
			listAttachments:              infra.ListAttachments,
			createAttachment:             infra.CreateAttachment,
			deleteAttachment:             infra.DeleteAttachment,
			newAttachmentID:              infra.NewAttachmentID,
			pricePlanRoutes:              subscriptiondom.PricePlanRoutes(priceplanpkg.DefaultRoutes()),
			priceScheduleRoutes:          subscriptiondom.PriceScheduleRoutes(psRoutes),
			priceScheduleInventoryRoutes: subscriptiondom.PriceScheduleRoutes(priceScheduleInventoryRoutes),
			priceListRoutes:              productdom.PriceListRoutes(pricelistpkg.DefaultRoutes()),
			planRoutes:                   subscriptiondom.PlanRoutes(*r),
			planBundleRoutes:             subscriptiondom.PlanRoutes(planpkg.DefaultBundleRoutes()),
			subscriptionRoutes:           subscriptiondom.SubscriptionRoutes(subRoutes),
			pricePlanLabels:              subscriptiondom.PricePlanLabels(pricePlanLabels),
			productPricePlanLabels:       subscriptiondom.ProductPricePlanLabels(productPricePlanLabels),
			priceScheduleLabels:          subscriptiondom.PriceScheduleLabels(priceScheduleLabels),
			priceListLabels:              productdom.PriceListLabels{},
			planLabels:                   subscriptiondom.PlanLabels(*l),
			centymoTableLabels:           mc.Table,
		})
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Subscription
// ---------------------------------------------------------------------------

// SubscriptionUnit wires the subscription domain by delegating to the existing
// wireSubscriptionModule helper (same helper Block() uses). This avoids
// duplicating the 540-line sub-package registration logic; catalog.go is in
// the block package so the private helpers are accessible directly.
func SubscriptionUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := subscriptionpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*subscriptionpkg.Routes)
		l := u.Labels.(*subscriptionpkg.Labels)

		priceScheduleLabels := priceschedulepkg.DefaultLabels()
		if psr, ok := compose.RoutesOf[*priceschedulepkg.Routes](mc, "subscription.price_schedule"); ok {
			// Reuse the sibling's labels if available — they carry any lyngua overlay.
			_ = psr // routes already set in wiring below
		}

		minCtx := &pyeza.AppContext{
			Routes: mc.Routes,
			Common: mc.Common,
		}

		// Look up price-schedule routes for the breadcrumb/subscription-tab URLs.
		psRoutes := priceschedulepkg.DefaultRoutes()
		if psr, ok := compose.RoutesOf[*priceschedulepkg.Routes](mc, "subscription.price_schedule"); ok {
			psRoutes = *psr
		}

		wireSubscriptionModule(minCtx, allEnabledConfig(), uc, subscriptionWiring{
			refChecker:          infra.RefChecker,
			uploadFile:          infra.UploadFile,
			downloadFile:        infra.DownloadFile,
			readAttachment:      infra.ReadAttachment,
			listAttachments:     infra.ListAttachments,
			createAttachment:    infra.CreateAttachment,
			deleteAttachment:    infra.DeleteAttachment,
			newAttachmentID:     infra.NewAttachmentID,
			subscriptionRoutes:  subscriptiondom.SubscriptionRoutes(*r),
			priceScheduleRoutes: subscriptiondom.PriceScheduleRoutes(psRoutes),
			subscriptionLabels:  subscriptiondom.SubscriptionLabels(*l),
			priceScheduleLabels: subscriptiondom.PriceScheduleLabels(priceScheduleLabels),
			centymoTableLabels:  mc.Table,
		})
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Collection
// ---------------------------------------------------------------------------

func CollectionUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := collectionpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*collectionpkg.Routes)
		l := u.Labels.(*collectionpkg.Labels)

		collectionAdvanceLabels := treasurydomain.DefaultTreasuryCollectionAdvanceLabels()
		advanceEnumLabels := treasurydomain.DefaultAdvanceEnumLabels()

		collDeps := &treasurydomain.CollectionModuleDeps{
			Routes:                   *r,
			Labels:                   *l,
			CommonLabels:             mc.Common,
			TableLabels:              mc.Table,
			CreateCollection:         uc.Collection.CreateCollection,
			ReadCollection:           uc.Collection.ReadCollection,
			UpdateCollection:         uc.Collection.UpdateCollection,
			DeleteCollection:         uc.Collection.DeleteCollection,
			ListCollections:          uc.Collection.ListCollections,
			UploadFile:               infra.UploadFile,
			ListAttachments:          infra.ListAttachments,
			CreateAttachment:         infra.CreateAttachment,
			DeleteAttachment:         infra.DeleteAttachment,
			NewID:                    infra.NewAttachmentID,
			AdvanceLabels:            collectionAdvanceLabels,
			AdvanceEnumLabels:        advanceEnumLabels,
			SettleUnscheduledAdvance: bridgeSettleAdvance(uc.Collection.SettleUnscheduledAdvance),
			RefundUnscheduledAdvance: bridgeRefundAdvance(uc.Collection.RefundUnscheduledAdvance),
			CancelAdvance:            bridgeCancelAdvance(uc.Collection.CancelAdvance),
		}
		wireCashDashboard(collDeps, uc)
		collDeps.GetFunctionalCurrency = func(fctx context.Context) string {
			return getFunctionalCurrency(fctx, uc)
		}
		treasurydomain.NewCollectionModule(collDeps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Disbursement
// ---------------------------------------------------------------------------

func DisbursementUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := disbursementpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*disbursementpkg.Routes)
		l := u.Labels.(*disbursementpkg.Labels)

		disbursementAdvanceLabels := treasurydomain.DefaultTreasuryDisbursementAdvanceLabels()
		advanceEnumLabels := treasurydomain.DefaultAdvanceEnumLabels()

		treasurydomain.NewDisbursementModule(&treasurydomain.DisbursementModuleDeps{
			Routes:                   *r,
			Labels:                   *l,
			CommonLabels:             mc.Common,
			TableLabels:              mc.Table,
			CreateDisbursement:       uc.Disbursement.CreateDisbursement,
			ReadDisbursement:         uc.Disbursement.ReadDisbursement,
			UpdateDisbursement:       uc.Disbursement.UpdateDisbursement,
			DeleteDisbursement:       uc.Disbursement.DeleteDisbursement,
			ListDisbursements:        uc.Disbursement.ListDisbursements,
			UploadFile:               infra.UploadFile,
			ListAttachments:          infra.ListAttachments,
			CreateAttachment:         infra.CreateAttachment,
			DeleteAttachment:         infra.DeleteAttachment,
			NewID:                    infra.NewAttachmentID,
			AdvanceLabels:            disbursementAdvanceLabels,
			AdvanceEnumLabels:        advanceEnumLabels,
			SettleUnscheduledAdvance: bridgeSettleAdvance(uc.Disbursement.SettleUnscheduledAdvance),
			RefundUnscheduledAdvance: bridgeRefundAdvance(uc.Disbursement.RefundUnscheduledAdvance),
			CancelAdvance:            bridgeCancelAdvance(uc.Disbursement.CancelAdvance),
		}).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Expenditure
// ---------------------------------------------------------------------------

func ExpenditureUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := expenditurepkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*expenditurepkg.Routes)
		// Expenditure descriptor is partial (no LabelJSON).
		labels := labelsOrZero[expenditurepkg.Labels](&u)

		var disbursementRoutes treasurydomain.DisbursementRoutes
		if dr, ok := compose.RoutesOf[*disbursementpkg.Routes](mc, "treasury.disbursement"); ok {
			disbursementRoutes = treasurydomain.DisbursementRoutes(*dr)
		}

		deps := &expendituredomain.ExpenditureModuleDeps{
			Routes:                 *r,
			Labels:                 labels,
			CommonLabels:           mc.Common,
			TableLabels:            mc.Table,
			ListDocumentTemplates:  infra.ListDocTemplates,
			CreateDocumentTemplate: infra.CreateDocTemplate,
			UpdateDocumentTemplate: infra.UpdateDocTemplate,
			DeleteDocumentTemplate: infra.DeleteDocTemplate,
			UploadFile:             infra.UploadTemplate,
			ListAttachments:        infra.ListAttachments,
			CreateAttachment:       infra.CreateAttachment,
			DeleteAttachment:       infra.DeleteAttachment,
			NewAttachmentID:        infra.NewAttachmentID,
		}
		if disbursementRoutes.DetailURL != "" && uc.Disbursement.CreateDisbursement != nil {
			deps.DisbursementRoutes = disbursementRoutes
			deps.CreateDisbursement = uc.Disbursement.CreateDisbursement
		}
		wireExpenditureDeps(deps, uc)
		wirePurchaseDashboard(deps, uc)
		wireExpenseDashboard(deps, uc)
		deps.GetFunctionalCurrency = func(fctx context.Context) string {
			return getFunctionalCurrency(fctx, uc)
		}
		expendituredomain.NewExpenditureModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

func wireExpenditureDeps(deps *expendituredomain.ExpenditureModuleDeps, uc *UseCases) {
	deps.ListExpenditures = uc.Expenditure.ListExpenditures
	deps.CreateExpenditure = uc.Expenditure.CreateExpenditure
	deps.ReadExpenditure = uc.Expenditure.ReadExpenditure
	deps.UpdateExpenditure = uc.Expenditure.UpdateExpenditure
	deps.DeleteExpenditure = uc.Expenditure.DeleteExpenditure
	deps.ListExpenditureCategories = uc.Expenditure.ListExpenditureCategories
	deps.CreateExpenditureCategory = uc.Expenditure.CreateExpenditureCategory
	deps.ReadExpenditureCategory = uc.Expenditure.ReadExpenditureCategory
	deps.UpdateExpenditureCategory = uc.Expenditure.UpdateExpenditureCategory
	deps.DeleteExpenditureCategory = uc.Expenditure.DeleteExpenditureCategory
	deps.CreateExpenditureLineItem = uc.Expenditure.CreateExpenditureLineItem
	deps.ReadExpenditureLineItem = uc.Expenditure.ReadExpenditureLineItem
	deps.UpdateExpenditureLineItem = uc.Expenditure.UpdateExpenditureLineItem
	deps.DeleteExpenditureLineItem = uc.Expenditure.DeleteExpenditureLineItem
	deps.ListExpenditureLineItems = uc.Expenditure.ListExpenditureLineItems
	deps.ListSuppliers = uc.Entity.Supplier.ListSuppliers
}

// ---------------------------------------------------------------------------
// ExpenseRecognition
// ---------------------------------------------------------------------------

func ExpenseRecognitionUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := expenserecognitionpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*expenserecognitionpkg.Routes)
		l := u.Labels.(*expenserecognitionpkg.Labels)

		erDeps := &expendituredomain.ExpenseRecognitionModuleDeps{
			Routes:       *r,
			Labels:       *l,
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,
			UploadFile:   infra.UploadFile,
			ListAttachments:  infra.ListAttachments,
			CreateAttachment: infra.CreateAttachment,
			DeleteAttachment: infra.DeleteAttachment,
			NewAttachmentID:  infra.NewAttachmentID,
		}
		erDeps.ListExpenseRecognitions = uc.Expenditure.ListExpenseRecognitions
		erDeps.ReadExpenseRecognition = uc.Expenditure.ReadExpenseRecognition
		erDeps.DeleteExpenseRecognition = uc.Expenditure.DeleteExpenseRecognition
		if uc.Expenditure.ReverseExpenseRecognition != nil {
			reverseUC := uc.Expenditure.ReverseExpenseRecognition
			erDeps.ReverseExpenseRecognition = func(fctx context.Context, id, reason string) error {
				req := &expenserecognitionpb.ReverseExpenseRecognitionRequest{ExpenseRecognitionId: id}
				if reason != "" {
					req.Reason = &reason
				}
				_, err := reverseUC(fctx, req)
				return err
			}
		}
		erDeps.ListExpenseRecognitionLines = uc.Expenditure.ListExpenseRecognitionLines
		expendituredomain.NewExpenseRecognitionModule(erDeps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// ExpenseRecognitionRun
// ---------------------------------------------------------------------------

func ExpenseRecognitionRunUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := expenserecognitionrunpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*expenserecognitionrunpkg.Routes)
		l := u.Labels.(*expenserecognitionrunpkg.Labels)

		// Resolve expenditure routes from the sibling unit if available.
		expenditureRoutes := expendituredomain.DefaultExpenditureRoutes()
		if er, ok := compose.RoutesOf[*expenditurepkg.Routes](mc, "expenditure.expenditure"); ok {
			expenditureRoutes = expendituredomain.ExpenditureRoutes(*er)
		}

		minCtx := &pyeza.AppContext{
			Routes: mc.Routes,
			Common: mc.Common,
		}
		wireExpenseRecognitionRunModule(minCtx, allEnabledConfig(), uc, expenseRecognitionRunWiring{
			routes:             expendituredomain.ExpenseRecognitionRunRoutes(*r),
			labels:             expendituredomain.ExpenseRecognitionRunLabels(*l),
			expenditureRoutes:  expenditureRoutes,
			centymoTableLabels: mc.Table,
		})
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// AccruedExpense
// ---------------------------------------------------------------------------

func AccruedExpenseUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := accruedexpensepkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*accruedexpensepkg.Routes)
		l := u.Labels.(*accruedexpensepkg.Labels)

		minCtx := &pyeza.AppContext{
			Routes: mc.Routes,
			Common: mc.Common,
		}
		wireAccruedExpenseModules(minCtx, allEnabledConfig(), uc, accruedExpenseWiring{
			accruedExpenseRoutes: expendituredomain.AccruedExpenseRoutes(*r),
			accruedExpenseLabels: expendituredomain.AccruedExpenseLabels(*l),
			centymoTableLabels:   mc.Table,
			uploadFile:           infra.UploadFile,
			listAttachments:      infra.ListAttachments,
			createAttachment:     infra.CreateAttachment,
			deleteAttachment:     infra.DeleteAttachment,
			newAttachmentID:      infra.NewAttachmentID,
		})
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// RevenueRun
// ---------------------------------------------------------------------------

func RevenueRunUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := revenuerunpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*revenuerunpkg.Routes)
		l := u.Labels.(*revenuerunpkg.Labels)

		// Resolve revenue routes from the sibling unit if available.
		revenueRoutes := revenuedomain.DefaultRevenueRoutes()
		if rr, ok := compose.RoutesOf[*revenuepkg.Routes](mc, "revenue.revenue"); ok {
			revenueRoutes = revenuedomain.RevenueRoutes(*rr)
		}

		minCtx := &pyeza.AppContext{
			Routes: mc.Routes,
			Common: mc.Common,
		}
		wireRevenueRunModule(minCtx, allEnabledConfig(), uc, revenueRunWiring{
			revenueRunRoutes:  revenuedomain.RevenueRunRoutes(*r),
			revenueRunLabels:  revenuedomain.RevenueRunLabels(*l),
			revenueRoutes:     revenueRoutes,
			centymoTableLabels: mc.Table,
			uploadFile:        infra.UploadFile,
			listAttachments:   infra.ListAttachments,
			createAttachment:  infra.CreateAttachment,
			deleteAttachment:  infra.DeleteAttachment,
			newAttachmentID:   infra.NewAttachmentID,
		})
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// SupplierContract
// ---------------------------------------------------------------------------

func SupplierContractUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := suppliercontractpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*suppliercontractpkg.Routes)
		l := u.Labels.(*suppliercontractpkg.Labels)

		// Resolve price-schedule routes for the cross-tab URLs from the sibling unit.
		scpsRoutes := expendituredomain.DefaultSupplierContractPriceScheduleRoutes()
		if pr, ok := compose.RoutesOf[*suppliercontractpriceschedulepkg.Routes](mc, "expenditure.supplier_contract_price_schedule"); ok {
			scpsRoutes = expendituredomain.SupplierContractPriceScheduleRoutes(*pr)
		}

		minCtx := &pyeza.AppContext{
			Routes: mc.Routes,
			Common: mc.Common,
		}
		wireSupplierCommitmentModules(minCtx, &blockConfig{supplierContract: true, supplierContractLine: true}, uc, supplierCommitmentWiring{
			supplierContractRoutes:              expendituredomain.SupplierContractRoutes(*r),
			supplierContractLabels:              expendituredomain.SupplierContractLabels(*l),
			supplierContractPriceScheduleRoutes: scpsRoutes,
			centymoTableLabels:                  mc.Table,
			uploadFile:                          infra.UploadFile,
			listAttachments:                     infra.ListAttachments,
			createAttachment:                    infra.CreateAttachment,
			deleteAttachment:                    infra.DeleteAttachment,
			newAttachmentID:                     infra.NewAttachmentID,
		})
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// ProcurementRequest
// ---------------------------------------------------------------------------

func ProcurementRequestUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := procurementrequestpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*procurementrequestpkg.Routes)
		l := u.Labels.(*procurementrequestpkg.Labels)

		minCtx := &pyeza.AppContext{
			Routes: mc.Routes,
			Common: mc.Common,
		}
		wireSupplierCommitmentModules(minCtx, &blockConfig{procurementRequest: true, procurementRequestLine: true}, uc, supplierCommitmentWiring{
			procurementRequestRoutes: expendituredomain.ProcurementRequestRoutes(*r),
			procurementRequestLabels: expendituredomain.ProcurementRequestLabels(*l),
			centymoTableLabels:       mc.Table,
			uploadFile:               infra.UploadFile,
			listAttachments:          infra.ListAttachments,
			createAttachment:         infra.CreateAttachment,
			deleteAttachment:         infra.DeleteAttachment,
			newAttachmentID:          infra.NewAttachmentID,
		})
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// SupplierContractPriceSchedule
// ---------------------------------------------------------------------------

func SupplierContractPriceScheduleUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := suppliercontractpriceschedulepkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*suppliercontractpriceschedulepkg.Routes)
		l := u.Labels.(*suppliercontractpriceschedulepkg.Labels)

		minCtx := &pyeza.AppContext{
			Routes: mc.Routes,
			Common: mc.Common,
		}
		wireSupplierContractPriceScheduleModules(minCtx, &blockConfig{supplierContractPriceSchedule: true, supplierContractPriceScheduleLine: true}, uc, supplierContractPriceScheduleWiring{
			supplierContractPriceScheduleRoutes: expendituredomain.SupplierContractPriceScheduleRoutes(*r),
			supplierContractPriceScheduleLabels: expendituredomain.SupplierContractPriceScheduleLabels(*l),
			centymoTableLabels:                  mc.Table,
			uploadFile:                          infra.UploadFile,
			listAttachments:                     infra.ListAttachments,
			createAttachment:                    infra.CreateAttachment,
			deleteAttachment:                    infra.DeleteAttachment,
			newAttachmentID:                     infra.NewAttachmentID,
		})
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// CostSchedule
// ---------------------------------------------------------------------------

func CostScheduleUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := costschedulepkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*costschedulepkg.Routes)
		l := u.Labels.(*costschedulepkg.Labels)

		csDeps := &procurementdomain.CostScheduleModuleDeps{
			Routes:                *r,
			Labels:                *l,
			CommonLabels:          mc.Common,
			TableLabels:           mc.Table,
			SetCostScheduleActive: setActiveClosure(uc, "cost_schedule"),
		}
		cs := uc.Procurement.CostSchedule
		csDeps.CreateCostSchedule = cs.CreateCostSchedule
		csDeps.ReadCostSchedule = cs.ReadCostSchedule
		csDeps.UpdateCostSchedule = cs.UpdateCostSchedule
		csDeps.DeleteCostSchedule = cs.DeleteCostSchedule
		csDeps.GetCostScheduleListPageData = cs.GetCostScheduleListPageData
		csDeps.GetCostScheduleItemPageData = cs.GetCostScheduleItemPageData
		procurementdomain.NewCostScheduleModule(csDeps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// SupplierPlan
// ---------------------------------------------------------------------------

func SupplierPlanUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := supplierplanpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*supplierplanpkg.Routes)
		l := u.Labels.(*supplierplanpkg.Labels)

		spDeps := &procurementdomain.SupplierPlanModuleDeps{
			Routes:                *r,
			Labels:                *l,
			CommonLabels:          mc.Common,
			TableLabels:           mc.Table,
			SetSupplierPlanActive: setActiveClosure(uc, "supplier_plan"),
		}
		sp := uc.Procurement.SupplierPlan
		spDeps.CreateSupplierPlan = sp.CreateSupplierPlan
		spDeps.ReadSupplierPlan = sp.ReadSupplierPlan
		spDeps.UpdateSupplierPlan = sp.UpdateSupplierPlan
		spDeps.DeleteSupplierPlan = sp.DeleteSupplierPlan
		spDeps.GetSupplierPlanListPageData = sp.GetSupplierPlanListPageData
		spDeps.GetSupplierPlanItemPageData = sp.GetSupplierPlanItemPageData
		procurementdomain.NewSupplierPlanModule(spDeps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// CostPlan
// ---------------------------------------------------------------------------

func CostPlanUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := costplanpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*costplanpkg.Routes)
		l := u.Labels.(*costplanpkg.Labels)

		cpDeps := &procurementdomain.CostPlanModuleDeps{
			Routes:       *r,
			Labels:       *l,
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,
		}
		cp := uc.Procurement.CostPlan
		cpDeps.CreateCostPlan = cp.CreateCostPlan
		cpDeps.ReadCostPlan = cp.ReadCostPlan
		cpDeps.UpdateCostPlan = cp.UpdateCostPlan
		cpDeps.DeleteCostPlan = cp.DeleteCostPlan
		cpDeps.GetCostPlanListPageData = cp.GetCostPlanListPageData
		cpDeps.GetCostPlanItemPageData = cp.GetCostPlanItemPageData
		procurementdomain.NewCostPlanModule(cpDeps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// SupplierProductPlan
// ---------------------------------------------------------------------------

func SupplierProductPlanUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := supplierproductplanpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*supplierproductplanpkg.Routes)
		l := u.Labels.(*supplierproductplanpkg.Labels)

		sppDeps := &procurementdomain.SupplierProductPlanModuleDeps{
			Routes:       *r,
			Labels:       *l,
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,
		}
		spp := uc.Procurement.SupplierProductPlan
		sppDeps.CreateSupplierProductPlan = spp.CreateSupplierProductPlan
		sppDeps.ReadSupplierProductPlan = spp.ReadSupplierProductPlan
		sppDeps.UpdateSupplierProductPlan = spp.UpdateSupplierProductPlan
		sppDeps.DeleteSupplierProductPlan = spp.DeleteSupplierProductPlan
		sppDeps.GetSupplierProductPlanListPageData = spp.GetSupplierProductPlanListPageData
		sppDeps.GetSupplierProductPlanItemPageData = spp.GetSupplierProductPlanItemPageData
		procurementdomain.NewSupplierProductPlanModule(sppDeps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// SupplierSubscription
// ---------------------------------------------------------------------------

func SupplierSubscriptionUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := suppliersubscriptionpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*suppliersubscriptionpkg.Routes)
		l := u.Labels.(*suppliersubscriptionpkg.Labels)

		ssDeps := &procurementdomain.SupplierSubscriptionModuleDeps{
			Routes:       *r,
			Labels:       *l,
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,
		}
		ss := uc.Procurement.SupplierSubscription
		ssDeps.CreateSupplierSubscription = ss.CreateSupplierSubscription
		ssDeps.ReadSupplierSubscription = ss.ReadSupplierSubscription
		ssDeps.UpdateSupplierSubscription = ss.UpdateSupplierSubscription
		ssDeps.DeleteSupplierSubscription = ss.DeleteSupplierSubscription
		ssDeps.GetSupplierSubscriptionListPageData = ss.GetSupplierSubscriptionListPageData
		ssDeps.GetSupplierSubscriptionItemPageData = ss.GetSupplierSubscriptionItemPageData
		ssDeps.ListExpenseRecognitions = uc.Expenditure.ListExpenseRecognitions
		procurementdomain.NewSupplierSubscriptionModule(ssDeps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// ProcurementDashboard
// ---------------------------------------------------------------------------

func ProcurementDashboardUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := procurementdashboardpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*procurementdashboardpkg.Routes)
		// Procurement dashboard descriptor is partial (no LabelJSON).
		labels := labelsOrZero[procurementdashboardpkg.Labels](&u)

		deps := &procurementdomain.ProcurementdashboardModuleDeps{
			Routes:       *r,
			Labels:       labels,
			CommonLabels: mc.Common,
		}
		// Nil-safe: missing list closures render empty states per view.
		deps.ListSupplierContracts = uc.SupplierContract.ListSupplierContracts
		deps.ListProcurementRequests = uc.SupplierContract.ListProcurementRequests
		deps.ListExpenditures = uc.Expenditure.ListExpenditures
		procurementdomain.NewProcurementdashboardModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// AdvancesDashboard (treasury)
// ---------------------------------------------------------------------------

func AdvancesDashboardUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := advancesdashboardpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*advancesdashboardpkg.Routes)
		l := u.Labels.(*advancesdashboardpkg.Labels)

		var collectionDetailURL, disbursementDetailURL string
		if cr, ok := compose.RoutesOf[*collectionpkg.Routes](mc, "treasury.collection"); ok {
			collectionDetailURL = cr.DetailURL
		}
		if dr, ok := compose.RoutesOf[*disbursementpkg.Routes](mc, "treasury.disbursement"); ok {
			disbursementDetailURL = dr.DetailURL
		}

		deps := &advancesdashboardpkg.ModuleDeps{
			Routes:                        *r,
			Labels:                        *l,
			CommonLabels:                  mc.Common,
			TableLabels:                   mc.Table,
			CollectionDetailURLTemplate:   collectionDetailURL,
			DisbursementDetailURLTemplate: disbursementDetailURL,
			GetFunctionalCurrency: func(fctx context.Context) string {
				return getFunctionalCurrency(fctx, uc)
			},
		}
		if uc.TreasuryAdvances.GetAdvancesDashboard != nil {
			wireAdvancesDashboardGet(deps, uc, collectionDetailURL, disbursementDetailURL)
		}
		advancesdashboardpkg.NewModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// wireAdvancesDashboardGet wires the GetDashboard closure into the deps.
// Extracted from AdvancesDashboardUnit for readability.
func wireAdvancesDashboardGet(
	deps *advancesdashboardpkg.ModuleDeps,
	uc *UseCases,
	collectionDetailURL, disbursementDetailURL string,
) {
	cb := uc.TreasuryAdvances.GetAdvancesDashboard
	deps.GetDashboard = func(fctx context.Context, req advancesdashboardpkg.DashboardRequest) (*advancesdashboardpkg.DashboardResponse, error) {
		data, err := cb(fctx, req.AsOfDate)
		if err != nil || data == nil {
			return &advancesdashboardpkg.DashboardResponse{}, err
		}
		return &advancesdashboardpkg.DashboardResponse{
			Outflows: advancesdashboardpkg.AdvancesSection{
				Rows: convertAdvancesRows(data.Outflows, disbursementDetailURL),
			},
			Inflows: advancesdashboardpkg.AdvancesSection{
				Rows: convertAdvancesRows(data.Inflows, collectionDetailURL),
			},
			Position: advancesdashboardpkg.AdvancesPosition{
				OutflowTotalRemaining:  data.OutflowTotalRemaining,
				InflowTotalRemaining:   data.InflowTotalRemaining,
				OutflowActiveCount:     data.OutflowActiveCount,
				InflowActiveCount:      data.InflowActiveCount,
				OutflowFullyRecognized: data.OutflowFullyRecognized,
				InflowFullyRecognized:  data.InflowFullyRecognized,
				Currency:               data.Currency,
			},
		}, nil
	}
}

// ---------------------------------------------------------------------------
// SupplierBillingEvent
// ---------------------------------------------------------------------------

func SupplierBillingEventUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := supplierbillingeventpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		l := u.Labels.(*supplierbillingeventpkg.Labels)

		// Routes for supplier_billing_event come from the TreasuryAdvancesRoutes
		// (the URLs live in the advances namespace). Resolve them from the
		// advances-dashboard sibling unit if available.
		var advancesRoutes treasurydomain.TreasuryAdvancesRoutes
		if ar, ok := compose.RoutesOf[*advancesdashboardpkg.Routes](mc, "treasury.treasuryadvancesdashboard"); ok {
			advancesRoutes = treasurydomain.TreasuryAdvancesRoutes(*ar)
		}

		deps := expendituredomain.SupplierBillingEventModuleDeps{
			Routes:                    advancesRoutes,
			Labels:                    expendituredomain.DefaultSupplierBillingEventLabels(),
			CommonLabels:              mc.Common,
			ListSupplierBillingEvents: uc.Expenditure.ListSupplierBillingEvents,
			ReadSupplierBillingEvent:  uc.Expenditure.ReadSupplierBillingEvent,
			Recognize:                 uc.TreasuryAdvances.RecognizeMilestoneAdvanceDisbursement,
		}
		_ = l // Labels pointer held for future compose overlay support
		expendituredomain.NewSupplierBillingEventModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// AllUnits — ordered to match Block() registration sequence
// ---------------------------------------------------------------------------

// AllUnits returns the complete curated unit list for all centymo commerce
// domains, in the same registration order as Block().
func AllUnits(uc *UseCases, infra *Infra) []compose.Unit {
	return []compose.Unit{
		InventoryUnit(uc, infra),
		RevenueUnit(uc, infra),
		ProductUnit(uc, infra),
		ResourceUnit(uc, infra),
		PriceListUnit(uc, infra),
		PricePlanUnit(uc, infra),
		PriceScheduleUnit(uc, infra),
		PlanUnit(uc, infra),
		SubscriptionUnit(uc, infra),
		CollectionUnit(uc, infra),
		DisbursementUnit(uc, infra),
		ExpenditureUnit(uc, infra),
		SupplierContractUnit(uc, infra),
		ProcurementRequestUnit(uc, infra),
		SupplierContractPriceScheduleUnit(uc, infra),
		ExpenseRecognitionUnit(uc, infra),
		AccruedExpenseUnit(uc, infra),
		RevenueRunUnit(uc, infra),
		ExpenseRecognitionRunUnit(uc, infra),
		AdvancesDashboardUnit(uc, infra),
		SupplierBillingEventUnit(uc, infra),
		CostScheduleUnit(uc, infra),
		SupplierPlanUnit(uc, infra),
		CostPlanUnit(uc, infra),
		SupplierProductPlanUnit(uc, infra),
		SupplierSubscriptionUnit(uc, infra),
		ProcurementDashboardUnit(uc, infra),
	}
}
