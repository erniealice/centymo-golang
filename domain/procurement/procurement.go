// Package procurement is the procurement-domain consumer facade (centymo restructure).
//
// PURE RE-EXPORT — zero behaviour. The procurement domain's data/route types,
// Default* constructors, and URL consts moved into per-entity packages under
// domain/procurement/<entity>/ with entity-local names (the <Entity> prefix stripped).
// This facade re-adds the original prefixed names so existing consumers
// (block/, service-admin) keep resolving procurement.<Entity>Labels /
// procurement.Default<Entity>Routes() / procurement.<Entity>ListURL unchanged.
//
// An entity package MUST NEVER import this facade (that would be an import
// cycle procurement -> <entity> -> procurement); cross-entity references go DIRECT to the
// sibling package.
package procurement

import (
	costplanpkg "github.com/erniealice/centymo-golang/domain/procurement/cost_plan"
	costschedulepkg "github.com/erniealice/centymo-golang/domain/procurement/cost_schedule"
	procurementdashboardpkg "github.com/erniealice/centymo-golang/domain/procurement/procurementdashboard"
	supplierplanpkg "github.com/erniealice/centymo-golang/domain/procurement/supplier_plan"
	supplierproductcostplanpkg "github.com/erniealice/centymo-golang/domain/procurement/supplier_product_cost_plan"
	supplierproductplanpkg "github.com/erniealice/centymo-golang/domain/procurement/supplier_product_plan"
	suppliersubscriptionpkg "github.com/erniealice/centymo-golang/domain/procurement/supplier_subscription"
)

// Re-exported data/route types (type aliases — identity-preserving).
type (
	CostPlanActionLabels                   = costplanpkg.ActionLabels
	CostPlanBulkLabels                     = costplanpkg.BulkLabels
	CostPlanButtonLabels                   = costplanpkg.ButtonLabels
	CostPlanColumnLabels                   = costplanpkg.ColumnLabels
	CostPlanConfirmLabels                  = costplanpkg.ConfirmLabels
	CostPlanDetailLabels                   = costplanpkg.DetailLabels
	CostPlanEmptyLabels                    = costplanpkg.EmptyLabels
	CostPlanErrorLabels                    = costplanpkg.ErrorLabels
	CostPlanFormLabels                     = costplanpkg.FormLabels
	CostPlanLabels                         = costplanpkg.Labels
	CostPlanPageLabels                     = costplanpkg.PageLabels
	CostPlanRoutes                         = costplanpkg.Routes
	CostPlanStatusLabels                   = costplanpkg.StatusLabels
	CostPlanTabLabels                      = costplanpkg.TabLabels
	CostScheduleActionLabels               = costschedulepkg.ActionLabels
	CostScheduleBulkLabels                 = costschedulepkg.BulkLabels
	CostScheduleButtonLabels               = costschedulepkg.ButtonLabels
	CostScheduleColumnLabels               = costschedulepkg.ColumnLabels
	CostScheduleConfirmLabels              = costschedulepkg.ConfirmLabels
	CostScheduleDetailLabels               = costschedulepkg.DetailLabels
	CostScheduleEmptyLabels                = costschedulepkg.EmptyLabels
	CostScheduleErrorLabels                = costschedulepkg.ErrorLabels
	CostScheduleFormLabels                 = costschedulepkg.FormLabels
	CostScheduleLabels                     = costschedulepkg.Labels
	CostSchedulePageLabels                 = costschedulepkg.PageLabels
	CostScheduleRoutes                     = costschedulepkg.Routes
	CostScheduleStatusLabels               = costschedulepkg.StatusLabels
	CostScheduleTabLabels                  = costschedulepkg.TabLabels
	ProcurementLabels                      = procurementdashboardpkg.Labels
	ProcurementRoutes                      = procurementdashboardpkg.Routes
	SupplierPlanActionLabels               = supplierplanpkg.ActionLabels
	SupplierPlanBulkLabels                 = supplierplanpkg.BulkLabels
	SupplierPlanButtonLabels               = supplierplanpkg.ButtonLabels
	SupplierPlanColumnLabels               = supplierplanpkg.ColumnLabels
	SupplierPlanConfirmLabels              = supplierplanpkg.ConfirmLabels
	SupplierPlanDetailLabels               = supplierplanpkg.DetailLabels
	SupplierPlanEmptyLabels                = supplierplanpkg.EmptyLabels
	SupplierPlanErrorLabels                = supplierplanpkg.ErrorLabels
	SupplierPlanFormLabels                 = supplierplanpkg.FormLabels
	SupplierPlanLabels                     = supplierplanpkg.Labels
	SupplierPlanPageLabels                 = supplierplanpkg.PageLabels
	SupplierPlanRoutes                     = supplierplanpkg.Routes
	SupplierPlanStatusLabels               = supplierplanpkg.StatusLabels
	SupplierPlanTabLabels                  = supplierplanpkg.TabLabels
	SupplierProductCostPlanActionLabels    = supplierproductcostplanpkg.ActionLabels
	SupplierProductCostPlanColumnLabels    = supplierproductcostplanpkg.ColumnLabels
	SupplierProductCostPlanEmptyLabels     = supplierproductcostplanpkg.EmptyLabels
	SupplierProductCostPlanErrorLabels     = supplierproductcostplanpkg.ErrorLabels
	SupplierProductCostPlanFormLabels      = supplierproductcostplanpkg.FormLabels
	SupplierProductCostPlanLabels          = supplierproductcostplanpkg.Labels
	SupplierProductPlanActionLabels        = supplierproductplanpkg.ActionLabels
	SupplierProductPlanBulkLabels          = supplierproductplanpkg.BulkLabels
	SupplierProductPlanButtonLabels        = supplierproductplanpkg.ButtonLabels
	SupplierProductPlanColumnLabels        = supplierproductplanpkg.ColumnLabels
	SupplierProductPlanConfirmLabels       = supplierproductplanpkg.ConfirmLabels
	SupplierProductPlanDetailLabels        = supplierproductplanpkg.DetailLabels
	SupplierProductPlanEmptyLabels         = supplierproductplanpkg.EmptyLabels
	SupplierProductPlanErrorLabels         = supplierproductplanpkg.ErrorLabels
	SupplierProductPlanFormLabels          = supplierproductplanpkg.FormLabels
	SupplierProductPlanLabels              = supplierproductplanpkg.Labels
	SupplierProductPlanPageLabels          = supplierproductplanpkg.PageLabels
	SupplierProductPlanRoutes              = supplierproductplanpkg.Routes
	SupplierProductPlanStatusLabels        = supplierproductplanpkg.StatusLabels
	SupplierProductPlanTabLabels           = supplierproductplanpkg.TabLabels
	SupplierSubscriptionActionLabels       = suppliersubscriptionpkg.ActionLabels
	SupplierSubscriptionBulkLabels         = suppliersubscriptionpkg.BulkLabels
	SupplierSubscriptionButtonLabels       = suppliersubscriptionpkg.ButtonLabels
	SupplierSubscriptionColumnLabels       = suppliersubscriptionpkg.ColumnLabels
	SupplierSubscriptionConfirmLabels      = suppliersubscriptionpkg.ConfirmLabels
	SupplierSubscriptionDetailLabels       = suppliersubscriptionpkg.DetailLabels
	SupplierSubscriptionEmptyLabels        = suppliersubscriptionpkg.EmptyLabels
	SupplierSubscriptionErrorLabels        = suppliersubscriptionpkg.ErrorLabels
	SupplierSubscriptionFormLabels         = suppliersubscriptionpkg.FormLabels
	SupplierSubscriptionLabels             = suppliersubscriptionpkg.Labels
	SupplierSubscriptionPageLabels         = suppliersubscriptionpkg.PageLabels
	SupplierSubscriptionRecognitionsLabels = suppliersubscriptionpkg.RecognitionsLabels
	SupplierSubscriptionRoutes             = suppliersubscriptionpkg.Routes
	SupplierSubscriptionStatusLabels       = suppliersubscriptionpkg.StatusLabels
	SupplierSubscriptionTabLabels          = suppliersubscriptionpkg.TabLabels
)

// Re-exported URL route consts (const-identity preserved).
const (
	CostPlanAddURL                          = costplanpkg.AddURL
	CostPlanBulkDeleteURL                   = costplanpkg.BulkDeleteURL
	CostPlanBulkSetStatusURL                = costplanpkg.BulkSetStatusURL
	CostPlanDeleteURL                       = costplanpkg.DeleteURL
	CostPlanDetailURL                       = costplanpkg.DetailURL
	CostPlanEditURL                         = costplanpkg.EditURL
	CostPlanListURL                         = costplanpkg.ListURL
	CostPlanProductCostAddURL               = costplanpkg.ProductCostAddURL
	CostPlanProductCostDeleteURL            = costplanpkg.ProductCostDeleteURL
	CostPlanProductCostEditURL              = costplanpkg.ProductCostEditURL
	CostPlanSetStatusURL                    = costplanpkg.SetStatusURL
	CostPlanTabActionURL                    = costplanpkg.TabActionURL
	CostPlanTableURL                        = costplanpkg.TableURL
	CostScheduleAddURL                      = costschedulepkg.AddURL
	CostScheduleBulkDeleteURL               = costschedulepkg.BulkDeleteURL
	CostScheduleBulkSetStatusURL            = costschedulepkg.BulkSetStatusURL
	CostScheduleDeleteURL                   = costschedulepkg.DeleteURL
	CostScheduleDetailURL                   = costschedulepkg.DetailURL
	CostScheduleEditURL                     = costschedulepkg.EditURL
	CostScheduleListURL                     = costschedulepkg.ListURL
	CostScheduleSetStatusURL                = costschedulepkg.SetStatusURL
	CostScheduleTabActionURL                = costschedulepkg.TabActionURL
	CostScheduleTableURL                    = costschedulepkg.TableURL
	ProcurementDashboardURL                 = procurementdashboardpkg.DashboardURL
	ProcurementRecurrenceDraftsURL          = procurementdashboardpkg.RecurrenceDraftsURL
	ProcurementRenewalCalendarURL           = procurementdashboardpkg.RenewalCalendarURL
	ProcurementUtilizationURL               = procurementdashboardpkg.UtilizationURL
	ProcurementVarianceURL                  = procurementdashboardpkg.VarianceURL
	SupplierPlanAddURL                      = supplierplanpkg.AddURL
	SupplierPlanBulkDeleteURL               = supplierplanpkg.BulkDeleteURL
	SupplierPlanBulkSetStatusURL            = supplierplanpkg.BulkSetStatusURL
	SupplierPlanDeleteURL                   = supplierplanpkg.DeleteURL
	SupplierPlanDetailURL                   = supplierplanpkg.DetailURL
	SupplierPlanEditURL                     = supplierplanpkg.EditURL
	SupplierPlanListURL                     = supplierplanpkg.ListURL
	SupplierPlanSetStatusURL                = supplierplanpkg.SetStatusURL
	SupplierPlanTabActionURL                = supplierplanpkg.TabActionURL
	SupplierPlanTableURL                    = supplierplanpkg.TableURL
	SupplierProductPlanAddURL               = supplierproductplanpkg.AddURL
	SupplierProductPlanBulkDeleteURL        = supplierproductplanpkg.BulkDeleteURL
	SupplierProductPlanBulkSetStatusURL     = supplierproductplanpkg.BulkSetStatusURL
	SupplierProductPlanDeleteURL            = supplierproductplanpkg.DeleteURL
	SupplierProductPlanDetailURL            = supplierproductplanpkg.DetailURL
	SupplierProductPlanEditURL              = supplierproductplanpkg.EditURL
	SupplierProductPlanListURL              = supplierproductplanpkg.ListURL
	SupplierProductPlanSetStatusURL         = supplierproductplanpkg.SetStatusURL
	SupplierProductPlanTabActionURL         = supplierproductplanpkg.TabActionURL
	SupplierProductPlanTableURL             = supplierproductplanpkg.TableURL
	SupplierSubscriptionAddURL              = suppliersubscriptionpkg.AddURL
	SupplierSubscriptionBulkDeleteURL       = suppliersubscriptionpkg.BulkDeleteURL
	SupplierSubscriptionBulkSetStatusURL    = suppliersubscriptionpkg.BulkSetStatusURL
	SupplierSubscriptionDeleteURL           = suppliersubscriptionpkg.DeleteURL
	SupplierSubscriptionDetailURL           = suppliersubscriptionpkg.DetailURL
	SupplierSubscriptionEditURL             = suppliersubscriptionpkg.EditURL
	SupplierSubscriptionListURL             = suppliersubscriptionpkg.ListURL
	SupplierSubscriptionRecognizeExpenseURL = suppliersubscriptionpkg.RecognizeExpenseURL
	SupplierSubscriptionSearchCostPlanURL   = suppliersubscriptionpkg.SearchCostPlanURL
	SupplierSubscriptionSearchSupplierURL   = suppliersubscriptionpkg.SearchSupplierURL
	SupplierSubscriptionSetStatusURL        = suppliersubscriptionpkg.SetStatusURL
	SupplierSubscriptionTabActionURL        = suppliersubscriptionpkg.TabActionURL
	SupplierSubscriptionTableURL            = suppliersubscriptionpkg.TableURL
)

// Re-exported Default* constructors (function values).
var (
	DefaultCostPlanLabels                = costplanpkg.DefaultLabels
	DefaultCostPlanRoutes                = costplanpkg.DefaultRoutes
	DefaultCostScheduleLabels            = costschedulepkg.DefaultLabels
	DefaultCostScheduleRoutes            = costschedulepkg.DefaultRoutes
	DefaultProcurementRoutes             = procurementdashboardpkg.DefaultRoutes
	DefaultSupplierPlanLabels            = supplierplanpkg.DefaultLabels
	DefaultSupplierPlanRoutes            = supplierplanpkg.DefaultRoutes
	DefaultSupplierProductCostPlanLabels = supplierproductcostplanpkg.DefaultLabels
	DefaultSupplierProductPlanLabels     = supplierproductplanpkg.DefaultLabels
	DefaultSupplierProductPlanRoutes     = supplierproductplanpkg.DefaultRoutes
	DefaultSupplierSubscriptionLabels    = suppliersubscriptionpkg.DefaultLabels
	DefaultSupplierSubscriptionRoutes    = suppliersubscriptionpkg.DefaultRoutes
)
