package variant

import (
	"context"
	"log"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantimagepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_image"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
)

// VariantFormLabels holds labels for the variant drawer form template.
type VariantFormLabels struct {
	SKU           string
	PriceOverride string
}

// OptionValueChoice represents a selectable option value in the variant form dropdown.
type OptionValueChoice struct {
	ID    string
	Label string
}

// OptionSelection represents a product option with its available values for the variant form.
type OptionSelection struct {
	OptionID   string
	OptionName string
	FieldName  string // "option_value_{optionID}"
	Values     []OptionValueChoice
	Selected   string // for edit mode
	Required   bool
}

// VariantFormData is the template data for the variant drawer form.
type VariantFormData struct {
	FormAction       string
	IsEdit           bool
	ID               string
	ProductID        string
	SKU              string
	PriceOverride    string
	Active           bool
	Labels           VariantFormLabels
	CommonLabels     any
	OptionSelections []OptionSelection
}

// DetailViewDeps holds dependencies for variant action handlers.
type DetailViewDeps struct {
	Routes       centymo.ProductRoutes
	DB           centymo.DataSource
	Labels       centymo.ProductLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Typed proto funcs for product_variant
	ReadProductVariant   func(ctx context.Context, req *productvariantpb.ReadProductVariantRequest) (*productvariantpb.ReadProductVariantResponse, error)
	CreateProductVariant func(ctx context.Context, req *productvariantpb.CreateProductVariantRequest) (*productvariantpb.CreateProductVariantResponse, error)
	UpdateProductVariant func(ctx context.Context, req *productvariantpb.UpdateProductVariantRequest) (*productvariantpb.UpdateProductVariantResponse, error)
	DeleteProductVariant func(ctx context.Context, req *productvariantpb.DeleteProductVariantRequest) (*productvariantpb.DeleteProductVariantResponse, error)

	// Typed proto funcs for product_variant_option
	ListProductVariantOptions  func(ctx context.Context, req *productvariantoptionpb.ListProductVariantOptionsRequest) (*productvariantoptionpb.ListProductVariantOptionsResponse, error)
	CreateProductVariantOption func(ctx context.Context, req *productvariantoptionpb.CreateProductVariantOptionRequest) (*productvariantoptionpb.CreateProductVariantOptionResponse, error)

	// Typed proto funcs for product_option/value (for dropdowns)
	ListProductOptions      func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ListProductOptionValues func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)

	// Typed proto funcs shared with detail.DetailViewDeps (for table building)
	ListProductVariants func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)

	// Typed proto funcs for variant page (product read, inventory)
	ReadProduct          func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	ListInventoryItems   func(ctx context.Context, req *inventoryitempb.ListInventoryItemsRequest) (*inventoryitempb.ListInventoryItemsResponse, error)
	ReadInventoryItem    func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	ListInventorySerials func(ctx context.Context, req *inventoryserialpb.ListInventorySerialsRequest) (*inventoryserialpb.ListInventorySerialsResponse, error)

	// Product Variant Image CRUD
	ListProductVariantImages  func(ctx context.Context, req *productvariantimagepb.ListProductVariantImagesRequest) (*productvariantimagepb.ListProductVariantImagesResponse, error)
	CreateProductVariantImage func(ctx context.Context, req *productvariantimagepb.CreateProductVariantImageRequest) (*productvariantimagepb.CreateProductVariantImageResponse, error)
	DeleteProductVariantImage func(ctx context.Context, req *productvariantimagepb.DeleteProductVariantImageRequest) (*productvariantimagepb.DeleteProductVariantImageResponse, error)

	// Storage uploader (for file content → object storage)
	UploadImage func(ctx context.Context, bucketName, objectKey string, content []byte, contentType string) error

	// PermissionEntity is the first argument to perms.Can(entity, action) for
	// variant-assign / variant-edit / variant-remove actions. Defaults to
	// "product". See centymo-golang/views/product/module.go ModuleDeps.
	PermissionEntity string

	attachment.AttachmentOps
	auditlog.AuditOps
}

// permEntity returns the configured PermissionEntity with a safe default.
func (d *DetailViewDeps) permEntity() string {
	if d == nil || d.PermissionEntity == "" {
		return "product"
	}
	return d.PermissionEntity
}

// loadOptionSelections loads all active product options with their values for the variant form dropdowns.
func loadOptionSelections(ctx context.Context, deps *DetailViewDeps, productID string) []OptionSelection {
	if deps.ListProductOptions == nil || deps.ListProductOptionValues == nil {
		return nil
	}

	optResp, err := deps.ListProductOptions(ctx, &productoptionpb.ListProductOptionsRequest{})
	if err != nil {
		log.Printf("Failed to list product options: %v", err)
		return nil
	}

	valResp, err := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
	if err != nil {
		log.Printf("Failed to list product option values: %v", err)
		return nil
	}

	// Group values by option ID
	valuesByOption := make(map[string][]OptionValueChoice)
	for _, v := range valResp.GetData() {
		oid := v.GetProductOptionId()
		active := v.GetActive()
		if oid == "" || !active {
			continue
		}
		vid := v.GetId()
		label := v.GetLabel()
		valuesByOption[oid] = append(valuesByOption[oid], OptionValueChoice{ID: vid, Label: label})
	}

	var selections []OptionSelection
	for _, o := range optResp.GetData() {
		pid := o.GetProductId()
		active := o.GetActive()
		if pid != productID || !active {
			continue
		}
		oid := o.GetId()
		name := o.GetName()
		selections = append(selections, OptionSelection{
			OptionID:   oid,
			OptionName: name,
			FieldName:  "option_value_" + oid,
			Values:     valuesByOption[oid],
			Required:   o.GetRequired(),
		})
	}
	return selections
}

// loadVariantOptionSelections loads existing option value selections for a variant (edit mode).
func loadVariantOptionSelections(ctx context.Context, deps *DetailViewDeps, variantID string) map[string]string {
	if deps.ListProductVariantOptions == nil || deps.ListProductOptionValues == nil {
		return nil
	}

	voResp, err := deps.ListProductVariantOptions(ctx, &productvariantoptionpb.ListProductVariantOptionsRequest{})
	if err != nil {
		log.Printf("Failed to list product variant options: %v", err)
		return nil
	}

	// Build valueID → optionID lookup
	valResp, err := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
	if err != nil {
		log.Printf("Failed to list product option values: %v", err)
		return nil
	}
	valueToOption := make(map[string]string)
	for _, v := range valResp.GetData() {
		vid := v.GetId()
		oid := v.GetProductOptionId()
		if vid != "" && oid != "" {
			valueToOption[vid] = oid
		}
	}

	// Map: optionID → selected valueID
	result := make(map[string]string)
	for _, vo := range voResp.GetData() {
		vid := vo.GetProductVariantId()
		if vid != variantID {
			continue
		}
		valueID := vo.GetProductOptionValueId()
		if optionID, ok := valueToOption[valueID]; ok {
			result[optionID] = valueID
		}
	}
	return result
}

// saveVariantOptions creates product_variant_option rows from form selections.
func saveVariantOptions(ctx context.Context, deps *DetailViewDeps, variantID string, form map[string][]string) {
	if deps.CreateProductVariantOption == nil {
		return
	}
	for key, vals := range form {
		if !strings.HasPrefix(key, "option_value_") || len(vals) == 0 || vals[0] == "" {
			continue
		}
		_, err := deps.CreateProductVariantOption(ctx, &productvariantoptionpb.CreateProductVariantOptionRequest{
			Data: &productvariantoptionpb.ProductVariantOption{
				ProductVariantId:     variantID,
				ProductOptionValueId: vals[0],
				Active:               true,
			},
		})
		if err != nil {
			log.Printf("Failed to create product_variant_option: %v", err)
		}
	}
}

// deleteVariantOptions removes all product_variant_option rows for a variant.
func deleteVariantOptions(ctx context.Context, deps *DetailViewDeps, variantID string) {
	if deps.ListProductVariantOptions == nil || deps.DB == nil {
		return
	}
	voResp, err := deps.ListProductVariantOptions(ctx, &productvariantoptionpb.ListProductVariantOptionsRequest{})
	if err != nil {
		log.Printf("Failed to list product_variant_option for delete: %v", err)
		return
	}
	for _, vo := range voResp.GetData() {
		vid := vo.GetProductVariantId()
		if vid != variantID {
			continue
		}
		voID := vo.GetId()
		if voID != "" {
			err := deps.DB.HardDelete(ctx, "product_variant_option", voID)
			if err != nil {
				log.Printf("Failed to hard-delete product_variant_option %s: %v", voID, err)
			}
		}
	}
}

// validateRequiredOptions checks that all required product options have a value selected.
func validateRequiredOptions(selections []OptionSelection, form map[string][]string) string {
	for _, sel := range selections {
		if !sel.Required {
			continue
		}
		vals := form[sel.FieldName]
		if len(vals) == 0 || vals[0] == "" {
			return sel.OptionName + " is required"
		}
	}
	return ""
}
