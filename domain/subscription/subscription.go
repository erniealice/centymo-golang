// Package subscription is the subscription-domain consumer facade (centymo restructure).
//
// PURE RE-EXPORT — zero behaviour. The subscription domain's data/route types,
// Default* constructors, and URL consts moved into per-entity packages under
// domain/subscription/<entity>/ with entity-local names (the <Entity> prefix stripped).
// This facade re-adds the original prefixed names so existing consumers
// (block/, service-admin) keep resolving subscription.<Entity>Labels /
// subscription.Default<Entity>Routes() / subscription.<Entity>ListURL unchanged.
//
// An entity package MUST NEVER import this facade (that would be an import
// cycle subscription -> <entity> -> subscription); cross-entity references go DIRECT to the
// sibling package.
package subscription

import (
	clientpackagespkg "github.com/erniealice/centymo-golang/domain/subscription/client_packages"
	planpkg "github.com/erniealice/centymo-golang/domain/subscription/plan"
	priceplanpkg "github.com/erniealice/centymo-golang/domain/subscription/price_plan"
	priceschedulepkg "github.com/erniealice/centymo-golang/domain/subscription/price_schedule"
	productpriceplanpkg "github.com/erniealice/centymo-golang/domain/subscription/product_price_plan"
	subscriptionpkg "github.com/erniealice/centymo-golang/domain/subscription/subscription"
)

// Re-exported data/route types (type aliases — identity-preserving).
type (
	ClientPackagesLabels                 = clientpackagespkg.Labels
	PlanActionLabels                     = planpkg.ActionLabels
	PlanBulkLabels                       = planpkg.BulkLabels
	PlanButtonLabels                     = planpkg.ButtonLabels
	PlanColumnLabels                     = planpkg.ColumnLabels
	PlanConfirmLabels                    = planpkg.ConfirmLabels
	PlanDetailLabels                     = planpkg.DetailLabels
	PlanEmptyLabels                      = planpkg.EmptyLabels
	PlanErrorLabels                      = planpkg.ErrorLabels
	PlanFilterLabels                     = planpkg.FilterLabels
	PlanFormLabels                       = planpkg.FormLabels
	PlanFormSectionLabels                = planpkg.FormSectionLabels
	PlanLabels                           = planpkg.Labels
	PlanPageLabels                       = planpkg.PageLabels
	PlanRoutes                           = planpkg.Routes
	PlanStatusLabels                     = planpkg.StatusLabels
	PlanTabLabels                        = planpkg.TabLabels
	PricePlanActionLabels                = priceplanpkg.ActionLabels
	PricePlanBillingSummaryCopy          = priceplanpkg.BillingSummaryCopy
	PricePlanBillingSummaryWarn          = priceplanpkg.BillingSummaryWarn
	PricePlanBulkLabels                  = priceplanpkg.BulkLabels
	PricePlanButtonLabels                = priceplanpkg.ButtonLabels
	PricePlanColumnLabels2               = priceplanpkg.ColumnLabels2
	PricePlanConfirmLabels               = priceplanpkg.ConfirmLabels
	PricePlanDetailLabels2               = priceplanpkg.DetailLabels2
	PricePlanEmptyLabels                 = priceplanpkg.EmptyLabels
	PricePlanErrorLabels                 = priceplanpkg.ErrorLabels
	PricePlanFormLabels                  = priceplanpkg.FormLabels
	PricePlanLabels                      = priceplanpkg.Labels
	PricePlanMessageLabels               = priceplanpkg.MessageLabels
	PricePlanPageLabels                  = priceplanpkg.PageLabels
	PricePlanParentContextLabels         = productpriceplanpkg.PricePlanParentContextLabels
	PricePlanProductPriceLabels          = priceplanpkg.ProductPriceLabels
	PricePlanRoutes                      = priceplanpkg.Routes
	PricePlanSubscriptionsSectionLabels  = priceplanpkg.SubscriptionsSectionLabels
	PricePlanSummaryByBasis              = priceplanpkg.SummaryByBasis
	PricePlanSummaryLines                = priceplanpkg.SummaryLines
	PricePlanTabLabels2                  = priceplanpkg.TabLabels2
	PriceScheduleBulkLabels              = priceschedulepkg.BulkLabels
	PriceScheduleButtonLabels            = priceschedulepkg.ButtonLabels
	PriceScheduleColumnLabels            = priceschedulepkg.ColumnLabels
	PriceScheduleConfirmLabels           = priceschedulepkg.ConfirmLabels
	PriceScheduleDetailLabels            = priceschedulepkg.DetailLabels
	PriceScheduleEmptyLabels             = priceschedulepkg.EmptyLabels
	PriceScheduleErrorLabels             = priceschedulepkg.ErrorLabels
	PriceScheduleFilterLabels            = priceschedulepkg.FilterLabels
	PriceScheduleFormLabels              = priceschedulepkg.FormLabels
	PriceScheduleLabels                  = priceschedulepkg.Labels
	PriceSchedulePageLabels              = priceschedulepkg.PageLabels
	PriceScheduleRoutes                  = priceschedulepkg.Routes
	PriceScheduleTabLabels               = priceschedulepkg.TabLabels
	ProductKindOptionLabels              = planpkg.ProductKindOptionLabels
	ProductPlanFormLabels                = planpkg.ProductPlanFormLabels
	ProductPricePlanFormLabels           = productpriceplanpkg.FormLabels
	ProductPricePlanLabels               = productpriceplanpkg.Labels
	SubscriptionActionLabels             = subscriptionpkg.ActionLabels
	SubscriptionBackfillLabels           = subscriptionpkg.BackfillLabels
	SubscriptionBulkLabels               = subscriptionpkg.BulkLabels
	SubscriptionButtonLabels             = subscriptionpkg.ButtonLabels
	SubscriptionColumnLabels             = subscriptionpkg.ColumnLabels
	SubscriptionConfirmLabels            = subscriptionpkg.ConfirmLabels
	SubscriptionDetailLabels             = subscriptionpkg.DetailLabels
	SubscriptionEmptyLabels              = subscriptionpkg.EmptyLabels
	SubscriptionErrorLabels              = subscriptionpkg.ErrorLabels
	SubscriptionFormLabels               = subscriptionpkg.FormLabels
	SubscriptionInvoicesLabels           = subscriptionpkg.InvoicesLabels
	SubscriptionInvoicesRowActionsLabels = subscriptionpkg.InvoicesRowActionsLabels
	SubscriptionJobsTabLabels            = subscriptionpkg.JobsTabLabels
	SubscriptionLabels                   = subscriptionpkg.Labels
	SubscriptionMilestoneLabels          = subscriptionpkg.MilestoneLabels
	SubscriptionOperationsLabels         = subscriptionpkg.OperationsLabels
	SubscriptionPageLabels               = subscriptionpkg.PageLabels
	SubscriptionRecognizeLabels          = subscriptionpkg.RecognizeLabels
	SubscriptionRevenueRunErrorLabels    = subscriptionpkg.RevenueRunErrorLabels
	SubscriptionRevenueRunLabels         = subscriptionpkg.RevenueRunLabels
	SubscriptionRoutes                   = subscriptionpkg.Routes
	SubscriptionSpawnLabels              = subscriptionpkg.SpawnLabels
	SubscriptionStatusLabels             = subscriptionpkg.StatusLabels
	SubscriptionTabLabels                = subscriptionpkg.TabLabels
)

// Re-exported URL route consts (const-identity preserved).
const (
	MilestoneMarkReadyURL                  = subscriptionpkg.MilestoneMarkReadyURL
	MilestoneRecognizeURL                  = subscriptionpkg.MilestoneRecognizeURL
	MilestoneWaiveURL                      = subscriptionpkg.MilestoneWaiveURL
	PlanAddURL                             = planpkg.AddURL
	PlanAttachmentDeleteURL                = planpkg.AttachmentDeleteURL
	PlanAttachmentUploadURL                = planpkg.AttachmentUploadURL
	PlanBulkDeleteURL                      = planpkg.BulkDeleteURL
	PlanBulkSetStatusURL                   = planpkg.BulkSetStatusURL
	PlanDeleteURL                          = planpkg.DeleteURL
	PlanDetailURL                          = planpkg.DetailURL
	PlanEditURL                            = planpkg.EditURL
	PlanListURL                            = planpkg.ListURL
	PlanPricePlanDeleteURL                 = planpkg.PlanPricePlanDeleteURL
	PlanPricePlanEditURL                   = planpkg.PlanPricePlanEditURL
	PlanProductPlanAddURL                  = planpkg.ProductPlanAddURL
	PlanProductPlanDeleteURL               = planpkg.ProductPlanDeleteURL
	PlanProductPlanEditURL                 = planpkg.ProductPlanEditURL
	PlanProductPlanPickerURL               = planpkg.ProductPlanPickerURL
	PlanSetStatusURL                       = planpkg.SetStatusURL
	PlanTabActionURL                       = planpkg.TabActionURL
	PlanTableURL                           = planpkg.TableURL
	PricePlanAddURL                        = planpkg.PricePlanAddURL
	PricePlanAttachmentDeleteURL           = priceplanpkg.AttachmentDeleteURL
	PricePlanAttachmentUploadURL           = priceplanpkg.AttachmentUploadURL
	PricePlanBulkDeleteURL                 = priceplanpkg.BulkDeleteURL
	PricePlanBulkSetStatusURL              = priceplanpkg.BulkSetStatusURL
	PricePlanDashboardURL                  = priceplanpkg.DashboardURL
	PricePlanDeleteURL                     = planpkg.PricePlanDeleteURL
	PricePlanDetailURL                     = priceplanpkg.DetailURL
	PricePlanEditURL                       = planpkg.PricePlanEditURL
	PricePlanListURL                       = priceplanpkg.ListURL
	PricePlanProductPriceAddURL            = priceplanpkg.ProductPriceAddURL
	PricePlanProductPriceDeleteURL         = priceplanpkg.ProductPriceDeleteURL
	PricePlanProductPriceEditURL           = priceplanpkg.ProductPriceEditURL
	PricePlanSetStatusURL                  = priceplanpkg.SetStatusURL
	PricePlanStandaloneAddURL              = priceplanpkg.StandaloneAddURL
	PricePlanStandaloneDeleteURL           = priceplanpkg.StandaloneDeleteURL
	PricePlanStandaloneEditURL             = priceplanpkg.StandaloneEditURL
	PricePlanTabActionURL                  = priceplanpkg.TabActionURL
	PricePlanTableURL                      = priceplanpkg.TableURL
	PriceScheduleAddURL                    = priceschedulepkg.AddURL
	PriceScheduleAttachmentDeleteURL       = priceschedulepkg.AttachmentDeleteURL
	PriceScheduleAttachmentUploadURL       = priceschedulepkg.AttachmentUploadURL
	PriceScheduleBulkDeleteURL             = priceschedulepkg.BulkDeleteURL
	PriceScheduleBulkSetStatusURL          = priceschedulepkg.BulkSetStatusURL
	PriceScheduleDashboardURL              = priceschedulepkg.DashboardURL
	PriceScheduleDeleteURL                 = priceschedulepkg.DeleteURL
	PriceScheduleDetailURL                 = priceschedulepkg.DetailURL
	PriceScheduleEditURL                   = priceschedulepkg.EditURL
	PriceScheduleListURL                   = priceschedulepkg.ListURL
	PriceSchedulePlanProductPriceAddURL    = priceschedulepkg.PlanProductPriceAddURL
	PriceSchedulePlanProductPriceDeleteURL = priceschedulepkg.PlanProductPriceDeleteURL
	PriceSchedulePlanProductPriceEditURL   = priceschedulepkg.PlanProductPriceEditURL
	PriceSchedulePlanSubscriptionDetailURL = priceschedulepkg.PlanSubscriptionDetailURL
	PriceScheduleSetStatusURL              = priceschedulepkg.SetStatusURL
	PriceScheduleTabActionURL              = priceschedulepkg.TabActionURL
	PriceScheduleTableURL                  = priceschedulepkg.TableURL
	SubscriptionAddURL                     = subscriptionpkg.AddURL
	SubscriptionAttachmentDeleteURL        = subscriptionpkg.AttachmentDeleteURL
	SubscriptionAttachmentDownloadURL      = subscriptionpkg.AttachmentDownloadURL
	SubscriptionAttachmentUploadURL        = subscriptionpkg.AttachmentUploadURL
	SubscriptionBackfillCycleJobsURL       = subscriptionpkg.BackfillCycleJobsURL
	SubscriptionBulkDeleteURL              = subscriptionpkg.BulkDeleteURL
	SubscriptionBulkSetStatusURL           = subscriptionpkg.BulkSetStatusURL
	SubscriptionCustomizePackageURL        = subscriptionpkg.CustomizePackageURL
	SubscriptionDeleteURL                  = subscriptionpkg.DeleteURL
	SubscriptionDetailURL                  = subscriptionpkg.DetailURL
	SubscriptionEditURL                    = subscriptionpkg.EditURL
	SubscriptionListURL                    = subscriptionpkg.ListURL
	SubscriptionRecognizeURL               = subscriptionpkg.RecognizeURL
	SubscriptionRequestUsageURL            = subscriptionpkg.RequestUsageURL
	SubscriptionRevenueRunURL              = subscriptionpkg.RevenueRunURL
	SubscriptionSearchClientURL            = subscriptionpkg.SearchClientURL
	SubscriptionSearchPlanURL              = subscriptionpkg.SearchPlanURL
	SubscriptionSetStatusURL               = subscriptionpkg.SetStatusURL
	SubscriptionSpawnCycleJobsURL          = subscriptionpkg.SpawnCycleJobsURL
	SubscriptionSpawnJobsPartialURL        = subscriptionpkg.SpawnJobsPartialURL
	SubscriptionSpawnJobsURL               = subscriptionpkg.SpawnJobsURL
	SubscriptionTabActionURL               = subscriptionpkg.TabActionURL
	SubscriptionTableURL                   = subscriptionpkg.TableURL
	SubscriptionUnderClientDetailURL       = subscriptionpkg.UnderClientDetailURL
)

// Re-exported Default* constructors (function values).
var (
	DefaultClientPackagesLabels         = clientpackagespkg.DefaultLabels
	DefaultPlanBundleRoutes             = planpkg.DefaultBundleRoutes
	DefaultPlanLabels                   = planpkg.DefaultLabels
	DefaultPlanRoutes                   = planpkg.DefaultRoutes
	DefaultPricePlanLabels              = priceplanpkg.DefaultLabels
	DefaultPricePlanRoutes              = priceplanpkg.DefaultRoutes
	DefaultPriceScheduleInventoryRoutes = priceschedulepkg.DefaultInventoryRoutes
	DefaultPriceScheduleLabels          = priceschedulepkg.DefaultLabels
	DefaultPriceScheduleRoutes          = priceschedulepkg.DefaultRoutes
	DefaultProductPricePlanLabels       = productpriceplanpkg.DefaultLabels
	DefaultSubscriptionLabels           = subscriptionpkg.DefaultLabels
	DefaultSubscriptionRoutes           = subscriptionpkg.DefaultRoutes
)
