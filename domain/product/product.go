// Package product is the product-domain consumer facade (centymo restructure).
//
// PURE RE-EXPORT — zero behaviour. The product domain's data/route types,
// Default* constructors, and URL consts moved into per-entity packages under
// domain/product/<entity>/ with entity-local names (the <Entity> prefix stripped).
// This facade re-adds the original prefixed names so existing consumers
// (block/, service-admin) keep resolving product.<Entity>Labels /
// product.Default<Entity>Routes() / product.<Entity>ListURL unchanged.
//
// An entity package MUST NEVER import this facade (that would be an import
// cycle product -> <entity> -> product); cross-entity references go DIRECT to the
// sibling package.
package product

import (
	pricelistpkg "github.com/erniealice/centymo-golang/domain/product/price_list"
	productpkg "github.com/erniealice/centymo-golang/domain/product/product"
	resourcepkg "github.com/erniealice/centymo-golang/domain/product/resource"
)

// Re-exported data/route types (type aliases — identity-preserving).
type (
	DeliveryModeLabels             = productpkg.DeliveryModeLabels
	PriceListActionLabels          = pricelistpkg.ActionLabels
	PriceListBulkLabels            = pricelistpkg.BulkLabels
	PriceListButtonLabels          = pricelistpkg.ButtonLabels
	PriceListColumnLabels          = pricelistpkg.ColumnLabels
	PriceListConfirmLabels         = pricelistpkg.ConfirmLabels
	PriceListDetailLabels          = pricelistpkg.DetailLabels
	PriceListEmptyLabels           = pricelistpkg.EmptyLabels
	PriceListErrorLabels           = pricelistpkg.ErrorLabels
	PriceListFormLabels            = pricelistpkg.FormLabels
	PriceListLabels                = pricelistpkg.Labels
	PriceListPageLabels            = pricelistpkg.PageLabels
	PriceListRoutes                = pricelistpkg.Routes
	ProductActionLabels            = productpkg.ActionLabels
	ProductAttributeLabels         = productpkg.AttributeLabels
	ProductBreadcrumbLabels        = productpkg.BreadcrumbLabels
	ProductBulkLabels              = productpkg.BulkLabels
	ProductButtonLabels            = productpkg.ButtonLabels
	ProductColumnLabels            = productpkg.ColumnLabels
	ProductConfirmLabels           = productpkg.ConfirmLabels
	ProductDetailLabels            = productpkg.DetailLabels
	ProductEmptyLabels             = productpkg.EmptyLabels
	ProductErrorLabels             = productpkg.ErrorLabels
	ProductFormLabels              = productpkg.FormLabels
	ProductKindLabels              = productpkg.KindLabels
	ProductLabels                  = productpkg.Labels
	ProductLineActionLabels        = productpkg.LineActionLabels
	ProductLineBulkLabels          = productpkg.LineBulkLabels
	ProductLineButtonLabels        = productpkg.LineButtonLabels
	ProductLineColumnLabels        = productpkg.LineColumnLabels
	ProductLineConfirmLabels       = productpkg.LineConfirmLabels
	ProductLineDetailLabels        = productpkg.LineDetailLabels
	ProductLineEmptyLabels         = productpkg.LineEmptyLabels
	ProductLineErrorLabels         = productpkg.LineErrorLabels
	ProductLineFormLabels          = productpkg.LineFormLabels
	ProductLineLabels              = productpkg.LineLabels
	ProductLinePageLabels          = productpkg.LinePageLabels
	ProductLineRoutes              = productpkg.LineRoutes
	ProductLineStatusLabels        = productpkg.LineStatusLabels
	ProductLineTabLabels           = productpkg.LineTabLabels
	ProductOptionActionLabels      = productpkg.OptionActionLabels
	ProductOptionColumnLabels      = productpkg.OptionColumnLabels
	ProductOptionConfirmLabels     = productpkg.OptionConfirmLabels
	ProductOptionDataTypeLabels    = productpkg.OptionDataTypeLabels
	ProductOptionEmptyLabels       = productpkg.OptionEmptyLabels
	ProductOptionFormLabels        = productpkg.OptionFormLabels
	ProductOptionLabels            = productpkg.OptionLabels
	ProductOptionTabLabels         = productpkg.OptionTabLabels
	ProductOptionTabsLabels        = productpkg.OptionTabsLabels
	ProductOptionValueColumnLabels = productpkg.OptionValueColumnLabels
	ProductOptionValueFormLabels   = productpkg.OptionValueFormLabels
	ProductOptionValueLabels       = productpkg.OptionValueLabels
	ProductPageLabels              = productpkg.PageLabels
	ProductRoutes                  = productpkg.Routes
	ProductStatusLabels            = productpkg.StatusLabels
	ProductTabLabels               = productpkg.TabLabels
	ProductVariantLabels           = productpkg.VariantLabels
	ResourceActionLabels           = resourcepkg.ActionLabels
	ResourceBulkLabels             = resourcepkg.BulkLabels
	ResourceButtonLabels           = resourcepkg.ButtonLabels
	ResourceColumnLabels           = resourcepkg.ColumnLabels
	ResourceConfirmLabels          = resourcepkg.ConfirmLabels
	ResourceEmptyLabels            = resourcepkg.EmptyLabels
	ResourceErrorLabels            = resourcepkg.ErrorLabels
	ResourceFormLabels             = resourcepkg.FormLabels
	ResourceLabels                 = resourcepkg.Labels
	ResourcePageLabels             = resourcepkg.PageLabels
	ResourceRoutes                 = resourcepkg.Routes
	ResourceStatusLabels           = resourcepkg.StatusLabels
	ServiceDashboardLabels         = productpkg.ServiceDashboardLabels
	TrackingModeLabels             = productpkg.TrackingModeLabels
	VariantPricingLabels           = productpkg.VariantPricingLabels
)

// Re-exported URL route consts (const-identity preserved).
const (
	OptionValueSeparator                   = productpkg.OptionValueSeparator
	PriceListAddURL                        = pricelistpkg.AddURL
	PriceListAttachmentDeleteURL           = pricelistpkg.AttachmentDeleteURL
	PriceListAttachmentUploadURL           = pricelistpkg.AttachmentUploadURL
	PriceListBulkDeleteURL                 = pricelistpkg.BulkDeleteURL
	PriceListDeleteURL                     = pricelistpkg.DeleteURL
	PriceListDetailURL                     = pricelistpkg.DetailURL
	PriceListEditURL                       = pricelistpkg.EditURL
	PriceListListURL                       = pricelistpkg.ListURL
	PriceListTabActionURL                  = pricelistpkg.TabActionURL
	PriceListTableURL                      = pricelistpkg.TableURL
	PriceProductAddURL                     = pricelistpkg.PriceProductAddURL
	PriceProductDeleteURL                  = pricelistpkg.PriceProductDeleteURL
	ProductAddURL                          = productpkg.AddURL
	ProductAttachmentDeleteURL             = productpkg.AttachmentDeleteURL
	ProductAttachmentDownloadURL           = productpkg.AttachmentDownloadURL
	ProductAttachmentUploadURL             = productpkg.AttachmentUploadURL
	ProductAttributeAssignURL              = productpkg.AttributeAssignURL
	ProductAttributeRemoveURL              = productpkg.AttributeRemoveURL
	ProductAttributeTableURL               = productpkg.AttributeTableURL
	ProductBulkDeleteURL                   = productpkg.BulkDeleteURL
	ProductBulkSetStatusURL                = productpkg.BulkSetStatusURL
	ProductDeleteURL                       = productpkg.DeleteURL
	ProductDetailURL                       = productpkg.DetailURL
	ProductEditURL                         = productpkg.EditURL
	ProductLineAddURL                      = productpkg.LineAddURL
	ProductLineAttachmentDeleteURL         = productpkg.LineAttachmentDeleteURL
	ProductLineAttachmentUploadURL         = productpkg.LineAttachmentUploadURL
	ProductLineBulkDeleteURL               = productpkg.LineBulkDeleteURL
	ProductLineBulkSetStatusURL            = productpkg.LineBulkSetStatusURL
	ProductLineDashboardURL                = productpkg.LineDashboardURL
	ProductLineDeleteURL                   = productpkg.LineDeleteURL
	ProductLineDetailURL                   = productpkg.LineDetailURL
	ProductLineEditURL                     = productpkg.LineEditURL
	ProductLineListURL                     = productpkg.LineListURL
	ProductLineSetStatusURL                = productpkg.LineSetStatusURL
	ProductLineTabActionURL                = productpkg.LineTabActionURL
	ProductLineTableURL                    = productpkg.LineTableURL
	ProductListURL                         = productpkg.ListURL
	ProductOptionAddURL                    = productpkg.OptionAddURL
	ProductOptionDeleteURL                 = productpkg.OptionDeleteURL
	ProductOptionDetailURL                 = productpkg.OptionDetailURL
	ProductOptionEditURL                   = productpkg.OptionEditURL
	ProductOptionTableURL                  = productpkg.OptionTableURL
	ProductOptionValueAddURL               = productpkg.OptionValueAddURL
	ProductOptionValueDeleteURL            = productpkg.OptionValueDeleteURL
	ProductOptionValueEditURL              = productpkg.OptionValueEditURL
	ProductOptionValueTableURL             = productpkg.OptionValueTableURL
	ProductSetStatusURL                    = productpkg.SetStatusURL
	ProductTabActionURL                    = productpkg.TabActionURL
	ProductTableURL                        = productpkg.TableURL
	ProductVariantAssignURL                = productpkg.VariantAssignURL
	ProductVariantAttachmentDeleteURL      = productpkg.VariantAttachmentDeleteURL
	ProductVariantAttachmentUploadURL      = productpkg.VariantAttachmentUploadURL
	ProductVariantDetailURL                = productpkg.VariantDetailURL
	ProductVariantEditURL                  = productpkg.VariantEditURL
	ProductVariantImageDeleteURL           = productpkg.VariantImageDeleteURL
	ProductVariantImageUploadURL           = productpkg.VariantImageUploadURL
	ProductVariantRemoveURL                = productpkg.VariantRemoveURL
	ProductVariantSerialDetailURL          = productpkg.VariantSerialDetailURL
	ProductVariantStockAttachmentDeleteURL = productpkg.VariantStockAttachmentDeleteURL
	ProductVariantStockAttachmentUploadURL = productpkg.VariantStockAttachmentUploadURL
	ProductVariantStockDetailURL           = productpkg.VariantStockDetailURL
	ProductVariantStockTabActionURL        = productpkg.VariantStockTabActionURL
	ProductVariantTabActionURL             = productpkg.VariantTabActionURL
	ProductVariantTableURL                 = productpkg.VariantTableURL
	ResourceAddURL                         = resourcepkg.AddURL
	ResourceBulkDeleteURL                  = resourcepkg.BulkDeleteURL
	ResourceBulkSetStatusURL               = resourcepkg.BulkSetStatusURL
	ResourceDeleteURL                      = resourcepkg.DeleteURL
	ResourceDetailURL                      = resourcepkg.DetailURL
	ResourceEditURL                        = resourcepkg.EditURL
	ResourceListURL                        = resourcepkg.ListURL
	ResourceSetStatusURL                   = resourcepkg.SetStatusURL
	ResourceTableURL                       = resourcepkg.TableURL
	ServiceDashboardURL                    = productpkg.ServiceDashboardURL
)

// Re-exported Default* constructors (function values).
var (
	DefaultPriceListRoutes            = pricelistpkg.DefaultRoutes
	DefaultProductInventoryRoutes     = productpkg.DefaultInventoryRoutes
	DefaultProductLineInventoryRoutes = productpkg.DefaultLineInventoryRoutes
	DefaultProductLineLabels          = productpkg.DefaultLineLabels
	DefaultProductLineRoutes          = productpkg.DefaultLineRoutes
	DefaultProductRoutes              = productpkg.DefaultRoutes
	DefaultProductSuppliesRoutes      = productpkg.DefaultSuppliesRoutes
	DefaultResourceLabels             = resourcepkg.DefaultLabels
	DefaultResourceRoutes             = resourcepkg.DefaultRoutes
)
