package product

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang"
	productaction "github.com/erniealice/centymo-golang/views/product/action"
	productdetail "github.com/erniealice/centymo-golang/views/product/detail"
	productvariant "github.com/erniealice/centymo-golang/views/product/detail/variant"
	variantitem "github.com/erniealice/centymo-golang/views/product/detail/variant/item"
	variantitemserial "github.com/erniealice/centymo-golang/views/product/detail/variant/item/serial"
	productlist "github.com/erniealice/centymo-golang/views/product/list"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productattributepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_attribute"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantimagepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_image"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
)

// ModuleDeps holds all dependencies for the product module.
type ModuleDeps struct {
	Routes       centymo.ProductRoutes
	DB           centymo.DataSource
	Labels       centymo.ProductLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Deletable state
	GetInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// DefaultFulfillmentMethod is applied when the product create form does not
	// include a fulfillment_method value. Set to "service" for professional
	// business types so new products created through the services UI appear in
	// the services list (which filters fulfillment_method IN ('service','digital')).
	DefaultFulfillmentMethod string

	// Product CRUD
	ListProducts     func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ReadProduct      func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	CreateProduct    func(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error)
	UpdateProduct    func(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error)
	DeleteProduct    func(ctx context.Context, req *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error)
	SetProductActive func(ctx context.Context, id string, active bool) error

	// Product Variant CRUD
	ListProductVariants  func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	ReadProductVariant   func(ctx context.Context, req *productvariantpb.ReadProductVariantRequest) (*productvariantpb.ReadProductVariantResponse, error)
	CreateProductVariant func(ctx context.Context, req *productvariantpb.CreateProductVariantRequest) (*productvariantpb.CreateProductVariantResponse, error)
	UpdateProductVariant func(ctx context.Context, req *productvariantpb.UpdateProductVariantRequest) (*productvariantpb.UpdateProductVariantResponse, error)
	DeleteProductVariant func(ctx context.Context, req *productvariantpb.DeleteProductVariantRequest) (*productvariantpb.DeleteProductVariantResponse, error)

	// Product Variant Option
	ListProductVariantOptions  func(ctx context.Context, req *productvariantoptionpb.ListProductVariantOptionsRequest) (*productvariantoptionpb.ListProductVariantOptionsResponse, error)
	CreateProductVariantOption func(ctx context.Context, req *productvariantoptionpb.CreateProductVariantOptionRequest) (*productvariantoptionpb.CreateProductVariantOptionResponse, error)

	// Product Option CRUD
	ListProductOptions  func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ReadProductOption   func(ctx context.Context, req *productoptionpb.ReadProductOptionRequest) (*productoptionpb.ReadProductOptionResponse, error)
	CreateProductOption func(ctx context.Context, req *productoptionpb.CreateProductOptionRequest) (*productoptionpb.CreateProductOptionResponse, error)
	UpdateProductOption func(ctx context.Context, req *productoptionpb.UpdateProductOptionRequest) (*productoptionpb.UpdateProductOptionResponse, error)
	DeleteProductOption func(ctx context.Context, req *productoptionpb.DeleteProductOptionRequest) (*productoptionpb.DeleteProductOptionResponse, error)

	// Product Option Value CRUD
	ListProductOptionValues  func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	ReadProductOptionValue   func(ctx context.Context, req *productoptionvaluepb.ReadProductOptionValueRequest) (*productoptionvaluepb.ReadProductOptionValueResponse, error)
	CreateProductOptionValue func(ctx context.Context, req *productoptionvaluepb.CreateProductOptionValueRequest) (*productoptionvaluepb.CreateProductOptionValueResponse, error)
	UpdateProductOptionValue func(ctx context.Context, req *productoptionvaluepb.UpdateProductOptionValueRequest) (*productoptionvaluepb.UpdateProductOptionValueResponse, error)
	DeleteProductOptionValue func(ctx context.Context, req *productoptionvaluepb.DeleteProductOptionValueRequest) (*productoptionvaluepb.DeleteProductOptionValueResponse, error)

	// Product Attribute
	ListProductAttributes  func(ctx context.Context, req *productattributepb.ListProductAttributesRequest) (*productattributepb.ListProductAttributesResponse, error)
	CreateProductAttribute func(ctx context.Context, req *productattributepb.CreateProductAttributeRequest) (*productattributepb.CreateProductAttributeResponse, error)
	DeleteProductAttribute func(ctx context.Context, req *productattributepb.DeleteProductAttributeRequest) (*productattributepb.DeleteProductAttributeResponse, error)

	// Common Attribute (for attribute dropdowns)
	ListAttributes func(ctx context.Context, req *commonpb.ListAttributesRequest) (*commonpb.ListAttributesResponse, error)
	ReadAttribute  func(ctx context.Context, req *commonpb.ReadAttributeRequest) (*commonpb.ReadAttributeResponse, error)

	// Inventory (for variant detail page + variant stock detail)
	ListInventoryItems   func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	ReadInventoryItem    func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	ListInventorySerials func(ctx context.Context, req *inventoryserialpb.ListInventorySerialsRequest) (*inventoryserialpb.ListInventorySerialsResponse, error)

	// Product Variant Image CRUD
	ListProductVariantImages  func(ctx context.Context, req *productvariantimagepb.ListProductVariantImagesRequest) (*productvariantimagepb.ListProductVariantImagesResponse, error)
	CreateProductVariantImage func(ctx context.Context, req *productvariantimagepb.CreateProductVariantImageRequest) (*productvariantimagepb.CreateProductVariantImageResponse, error)
	DeleteProductVariantImage func(ctx context.Context, req *productvariantimagepb.DeleteProductVariantImageRequest) (*productvariantimagepb.DeleteProductVariantImageResponse, error)

	// Storage uploader
	UploadImage func(ctx context.Context, bucketName, objectKey string, content []byte, contentType string) error

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewID            func() string

	// Audit history
	ListAuditHistory func(ctx context.Context, req *auditlog.ListAuditRequest) (*auditlog.ListAuditResponse, error)
}

// Module holds all constructed product views.
type Module struct {
	routes        centymo.ProductRoutes
	List          view.View
	Detail        view.View
	TabAction     view.View
	Add           view.View
	Edit          view.View
	Delete        view.View
	BulkDelete    view.View
	SetStatus     view.View
	BulkSetStatus view.View
	VariantTable  view.View
	VariantAssign view.View
	VariantEdit   view.View
	VariantRemove view.View
	// Option detail page (option values management)
	OptionPage view.View
	// Variant detail page + tab action
	VariantPage           view.View
	VariantTabAction      view.View
	VariantStockDetail    view.View
	VariantStockTabAction view.View
	VariantSerialDetail   view.View
	// Options
	OptionTable       view.View
	OptionAdd         view.View
	OptionEdit        view.View
	OptionDelete      view.View
	OptionValueTable  view.View
	OptionValueAdd    view.View
	OptionValueEdit   view.View
	OptionValueDelete view.View
	// Attributes
	AttributeAssign view.View
	AttributeRemove view.View
	// Variant Images
	VariantImageUpload view.View
	VariantImageDelete view.View
	// Attachments (product, variant, variant stock)
	AttachmentUpload             view.View
	AttachmentDelete             view.View
	VariantAttachmentUpload      view.View
	VariantAttachmentDelete      view.View
	VariantStockAttachmentUpload view.View
	VariantStockAttachmentDelete view.View
}

func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &productaction.Deps{
		Routes:                   deps.Routes,
		Labels:                   deps.Labels,
		CreateProduct:            deps.CreateProduct,
		ReadProduct:              deps.ReadProduct,
		UpdateProduct:            deps.UpdateProduct,
		DeleteProduct:            deps.DeleteProduct,
		SetProductActive:         deps.SetProductActive,
		DefaultFulfillmentMethod: deps.DefaultFulfillmentMethod,
	}
	detailDeps := &productdetail.DetailViewDeps{
		ReadProduct:               deps.ReadProduct,
		Labels:                    deps.Labels,
		CommonLabels:              deps.CommonLabels,
		TableLabels:               deps.TableLabels,
		DB:                        deps.DB,
		ListProductVariants:       deps.ListProductVariants,
		ListProductOptions:        deps.ListProductOptions,
		ListProductOptionValues:   deps.ListProductOptionValues,
		ListProductVariantOptions: deps.ListProductVariantOptions,
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
		AuditOps: auditlog.AuditOps{
			ListAuditHistory: deps.ListAuditHistory,
		},
	}
	variantDeps := &productvariant.DetailViewDeps{
		DB:                         deps.DB,
		Labels:                     deps.Labels,
		CommonLabels:               deps.CommonLabels,
		TableLabels:                deps.TableLabels,
		ReadProductVariant:         deps.ReadProductVariant,
		CreateProductVariant:       deps.CreateProductVariant,
		UpdateProductVariant:       deps.UpdateProductVariant,
		DeleteProductVariant:       deps.DeleteProductVariant,
		ListProductVariantOptions:  deps.ListProductVariantOptions,
		CreateProductVariantOption: deps.CreateProductVariantOption,
		ListProductOptions:         deps.ListProductOptions,
		ListProductOptionValues:    deps.ListProductOptionValues,
		ListProductVariants:        deps.ListProductVariants,
		ReadProduct:                deps.ReadProduct,
		ListInventoryItems:         deps.ListInventoryItems,
		ReadInventoryItem:          deps.ReadInventoryItem,
		ListInventorySerials:       deps.ListInventorySerials,
		ListProductVariantImages:   deps.ListProductVariantImages,
		CreateProductVariantImage:  deps.CreateProductVariantImage,
		DeleteProductVariantImage:  deps.DeleteProductVariantImage,
		UploadImage:                deps.UploadImage,
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
		AuditOps: auditlog.AuditOps{
			ListAuditHistory: deps.ListAuditHistory,
		},
	}
	optionDeps := &productdetail.OptionsDeps{
		DB:                       deps.DB,
		Labels:                   deps.Labels,
		CommonLabels:             deps.CommonLabels,
		TableLabels:              deps.TableLabels,
		ListProductOptions:       deps.ListProductOptions,
		ReadProductOption:        deps.ReadProductOption,
		CreateProductOption:      deps.CreateProductOption,
		UpdateProductOption:      deps.UpdateProductOption,
		DeleteProductOption:      deps.DeleteProductOption,
		ListProductOptionValues:  deps.ListProductOptionValues,
		ReadProductOptionValue:   deps.ReadProductOptionValue,
		CreateProductOptionValue: deps.CreateProductOptionValue,
		UpdateProductOptionValue: deps.UpdateProductOptionValue,
		DeleteProductOptionValue: deps.DeleteProductOptionValue,
		ReadProduct:              deps.ReadProduct,
	}
	attributeDeps := &productdetail.AttributeDeps{
		DB:                     deps.DB,
		Labels:                 deps.Labels,
		CommonLabels:           deps.CommonLabels,
		TableLabels:            deps.TableLabels,
		ListAttributes:         deps.ListAttributes,
		ReadAttribute:          deps.ReadAttribute,
		ListProductAttributes:  deps.ListProductAttributes,
		CreateProductAttribute: deps.CreateProductAttribute,
		DeleteProductAttribute: deps.DeleteProductAttribute,
	}

	return &Module{
		routes: deps.Routes,
		List: productlist.NewView(&productlist.ListViewDeps{
			Routes:       deps.Routes,
			ListProducts: deps.ListProducts,
			GetInUseIDs:  deps.GetInUseIDs,
			Labels:       deps.Labels,
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
		}),
		Detail:                       productdetail.NewView(detailDeps),
		TabAction:                    productdetail.NewTabAction(detailDeps),
		Add:                          productaction.NewAddAction(actionDeps),
		Edit:                         productaction.NewEditAction(actionDeps),
		Delete:                       productaction.NewDeleteAction(actionDeps),
		BulkDelete:                   productaction.NewBulkDeleteAction(actionDeps),
		SetStatus:                    productaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:                productaction.NewBulkSetStatusAction(actionDeps),
		VariantTable:                 productvariant.NewTableView(variantDeps),
		VariantAssign:                productvariant.NewAssignView(variantDeps),
		VariantEdit:                  productvariant.NewEditView(variantDeps),
		VariantRemove:                productvariant.NewRemoveView(variantDeps),
		OptionPage:                   productdetail.NewOptionPageView(optionDeps),
		VariantPage:                  productvariant.NewPageView(variantDeps),
		VariantTabAction:             productvariant.NewTabAction(variantDeps),
		VariantStockDetail:           variantitem.NewPageView(variantDeps),
		VariantStockTabAction:        variantitem.NewTabAction(variantDeps),
		VariantSerialDetail:          variantitemserial.NewPageView(variantDeps),
		OptionTable:                  productdetail.NewOptionsTableView(optionDeps),
		OptionAdd:                    productdetail.NewOptionAddView(optionDeps),
		OptionEdit:                   productdetail.NewOptionEditView(optionDeps),
		OptionDelete:                 productdetail.NewOptionDeleteView(optionDeps),
		OptionValueTable:             productdetail.NewOptionValueTableView(optionDeps),
		OptionValueAdd:               productdetail.NewOptionValueAddView(optionDeps),
		OptionValueEdit:              productdetail.NewOptionValueEditView(optionDeps),
		OptionValueDelete:            productdetail.NewOptionValueDeleteView(optionDeps),
		AttributeAssign:              productdetail.NewAttributeAssignView(attributeDeps),
		AttributeRemove:              productdetail.NewAttributeRemoveView(attributeDeps),
		VariantImageUpload:           productvariant.NewImageUploadAction(variantDeps),
		VariantImageDelete:           productvariant.NewImageDeleteAction(variantDeps),
		AttachmentUpload:             productdetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete:             productdetail.NewAttachmentDeleteAction(detailDeps),
		VariantAttachmentUpload:      productvariant.NewAttachmentUploadAction(variantDeps),
		VariantAttachmentDelete:      productvariant.NewAttachmentDeleteAction(variantDeps),
		VariantStockAttachmentUpload: variantitem.NewAttachmentUploadAction(variantDeps),
		VariantStockAttachmentDelete: variantitem.NewAttachmentDeleteAction(variantDeps),
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	// Option detail page (option values)
	r.GET(m.routes.OptionDetailURL, m.OptionPage)
	// Variant detail page + tab switching
	r.GET(m.routes.VariantDetailURL, m.VariantPage)
	r.GET(m.routes.VariantTabActionURL, m.VariantTabAction)
	// Variant stock detail (serial/IMEI data for inventory item within variant context)
	r.GET(m.routes.VariantStockDetailURL, m.VariantStockDetail)
	r.GET(m.routes.VariantStockTabActionURL, m.VariantStockTabAction)
	r.GET(m.routes.VariantSerialDetailURL, m.VariantSerialDetail)
	// Variants
	r.GET(m.routes.VariantTableURL, m.VariantTable)
	r.GET(m.routes.VariantAssignURL, m.VariantAssign)
	r.POST(m.routes.VariantAssignURL, m.VariantAssign)
	r.GET(m.routes.VariantEditURL, m.VariantEdit)
	r.POST(m.routes.VariantEditURL, m.VariantEdit)
	r.POST(m.routes.VariantRemoveURL, m.VariantRemove)
	// Options
	r.GET(m.routes.OptionTableURL, m.OptionTable)
	r.GET(m.routes.OptionAddURL, m.OptionAdd)
	r.POST(m.routes.OptionAddURL, m.OptionAdd)
	r.GET(m.routes.OptionEditURL, m.OptionEdit)
	r.POST(m.routes.OptionEditURL, m.OptionEdit)
	r.POST(m.routes.OptionDeleteURL, m.OptionDelete)
	// Option Values
	r.GET(m.routes.OptionValueTableURL, m.OptionValueTable)
	r.GET(m.routes.OptionValueAddURL, m.OptionValueAdd)
	r.POST(m.routes.OptionValueAddURL, m.OptionValueAdd)
	r.GET(m.routes.OptionValueEditURL, m.OptionValueEdit)
	r.POST(m.routes.OptionValueEditURL, m.OptionValueEdit)
	r.POST(m.routes.OptionValueDeleteURL, m.OptionValueDelete)
	// Attributes
	r.GET(m.routes.AttributeAssignURL, m.AttributeAssign)
	r.POST(m.routes.AttributeAssignURL, m.AttributeAssign)
	r.POST(m.routes.AttributeRemoveURL, m.AttributeRemove)
	// Variant Images
	r.POST(m.routes.VariantImageUploadURL, m.VariantImageUpload)
	r.POST(m.routes.VariantImageDeleteURL, m.VariantImageDelete)
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
		r.GET(m.routes.VariantAttachmentUploadURL, m.VariantAttachmentUpload)
		r.POST(m.routes.VariantAttachmentUploadURL, m.VariantAttachmentUpload)
		r.POST(m.routes.VariantAttachmentDeleteURL, m.VariantAttachmentDelete)
		r.GET(m.routes.VariantStockAttachmentUploadURL, m.VariantStockAttachmentUpload)
		r.POST(m.routes.VariantStockAttachmentUploadURL, m.VariantStockAttachmentUpload)
		r.POST(m.routes.VariantStockAttachmentDeleteURL, m.VariantStockAttachmentDelete)
	}
}
