package variant

import (
	"context"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/erniealice/pyeza-golang/route"

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
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantimagepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_image"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	productpriceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/product_price_plan"
)

// VariantFormLabels holds labels for the variant drawer form template.
//
// SKU / PriceOverride are variant-specific (sourced from
// centymo.ProductVariantLabels). PricePlaceholder / SelectOption / Active
// are shared form labels reused from centymo.ProductFormLabels — they show
// up on every drawer form, not just variants, so we don't duplicate the
// lyngua keys.
type VariantFormLabels struct {
	SKU                    string
	PriceOverride          string
	PricePlaceholder       string
	SelectOption           string
	Active                 string
	OptionNeedsValuesAlert string
	ViewValues             string
}

// OptionValueChoice represents a selectable option value in the variant form dropdown.
type OptionValueChoice struct {
	ID    string
	Label string
}

// OptionSelection represents a product option with its available values for the variant form.
//
// DataType drives template rendering and submit-time persistence:
//   - text_list / number_list / color_list → <select> populated from Values;
//     form key "option_value_<optionID>" carries the chosen product_option_value ID.
//   - free_text → <input type="text">;
//     form key "option_free_<optionID>" carries the raw text. On submit the
//     handler upserts a one-off product_option_value row (label=value=raw) and
//     links the variant to it, so the rest of the schema stays uniform.
//   - free_number → <input type="number" min max>;
//     same persistence model as free_text but with numeric bounds from the option.
type OptionSelection struct {
	OptionID      string
	OptionName    string
	DataType      string
	FieldName     string // "option_value_<optionID>" (select types) or "option_free_<optionID>" (free types)
	Values        []OptionValueChoice
	Selected      string  // for edit mode — value ID for select types
	SelectedLabel string  // for edit mode — raw text/number for free types
	MinValue      *float64
	MaxValue      *float64
	Required      bool
	SortOrder     int32
	// NeedsValuesAlert is true when the option is a required *select* type
	// (text_list / number_list / color_list) with zero value rows — the user
	// can't satisfy the required constraint until they add values. Free types
	// don't trigger this since they accept arbitrary input.
	NeedsValuesAlert bool
	// ViewValuesURL is the deep link to the option detail page where values
	// can be added. Surfaced from the alert so the user can fix the gap inline.
	ViewValuesURL string
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
	ListProductOptions       func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ListProductOptionValues  func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	CreateProductOptionValue func(ctx context.Context, req *productoptionvaluepb.CreateProductOptionValueRequest) (*productoptionvaluepb.CreateProductOptionValueResponse, error)

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

	// Pricing (for the Pricing tab on the variant detail page).
	// Join path: product_plan (filtered by product_id + variant_id)
	//   → product_price_plan (filtered by product_plan_id)
	//     → price_plan (for amount/currency/name)
	//       → price_schedule (for date range + rate-card name)
	//         → plan (for fallback package name)
	ListProductPlans     func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProductPricePlans func(ctx context.Context, req *productpriceplanpb.ListProductPricePlansRequest) (*productpriceplanpb.ListProductPricePlansResponse, error)
	ListPricePlans       func(ctx context.Context, req *priceplanpb.ListPricePlansRequest) (*priceplanpb.ListPricePlansResponse, error)
	ListPriceSchedules   func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	ListPlans            func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)

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
		dataType := o.GetDataType()
		fieldName := "option_value_" + oid
		if isFreeType(dataType) {
			fieldName = "option_free_" + oid
		}
		var minV, maxV *float64
		if o.MinValue != nil {
			v := o.GetMinValue()
			minV = &v
		}
		if o.MaxValue != nil {
			v := o.GetMaxValue()
			maxV = &v
		}
		required := o.GetRequired()
		values := valuesByOption[oid]
		needsAlert := required && !isFreeType(dataType) && len(values) == 0
		viewValuesURL := ""
		if needsAlert {
			viewValuesURL = route.ResolveURL(deps.Routes.OptionDetailURL, "id", productID, "oid", oid)
		}
		selections = append(selections, OptionSelection{
			OptionID:         oid,
			OptionName:       o.GetName(),
			DataType:         dataType,
			FieldName:        fieldName,
			Values:           values,
			MinValue:         minV,
			MaxValue:         maxV,
			Required:         required,
			SortOrder:        o.GetSortOrder(),
			NeedsValuesAlert: needsAlert,
			ViewValuesURL:    viewValuesURL,
		})
	}
	// Sort by sort_order ASC, name as tiebreaker — keeps the variant form
	// fields in the same order as the options table on the parent page.
	sort.SliceStable(selections, func(i, j int) bool {
		if selections[i].SortOrder != selections[j].SortOrder {
			return selections[i].SortOrder < selections[j].SortOrder
		}
		return selections[i].OptionName < selections[j].OptionName
	})
	return selections
}

// isFreeType reports whether a data_type stores its variant value as raw text/number
// (no predefined product_option_value list shown in the option detail page).
func isFreeType(dataType string) bool {
	return dataType == "free_text" || dataType == "free_number"
}

// VariantOptionSelection holds the existing option_value link for a variant in
// edit mode. ValueID is what the <select> uses for select-type options;
// ValueLabel is what the free-text/free-number inputs prefill from.
type VariantOptionSelection struct {
	ValueID    string
	ValueLabel string
}

// loadVariantOptionSelections loads existing option value selections for a variant (edit mode).
func loadVariantOptionSelections(ctx context.Context, deps *DetailViewDeps, variantID string) map[string]VariantOptionSelection {
	if deps.ListProductVariantOptions == nil || deps.ListProductOptionValues == nil {
		return nil
	}

	voResp, err := deps.ListProductVariantOptions(ctx, &productvariantoptionpb.ListProductVariantOptionsRequest{})
	if err != nil {
		log.Printf("Failed to list product variant options: %v", err)
		return nil
	}

	// Build valueID → (optionID, label) lookup
	valResp, err := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
	if err != nil {
		log.Printf("Failed to list product option values: %v", err)
		return nil
	}
	type valMeta struct {
		OptionID string
		Label    string
	}
	valueLookup := make(map[string]valMeta)
	for _, v := range valResp.GetData() {
		vid := v.GetId()
		if vid == "" {
			continue
		}
		valueLookup[vid] = valMeta{OptionID: v.GetProductOptionId(), Label: v.GetLabel()}
	}

	result := make(map[string]VariantOptionSelection)
	for _, vo := range voResp.GetData() {
		if vo.GetProductVariantId() != variantID {
			continue
		}
		valueID := vo.GetProductOptionValueId()
		meta, ok := valueLookup[valueID]
		if !ok || meta.OptionID == "" {
			continue
		}
		result[meta.OptionID] = VariantOptionSelection{ValueID: valueID, ValueLabel: meta.Label}
	}
	return result
}

// saveVariantOptions creates product_variant_option rows from form selections.
//
// For select-type options (text_list / number_list / color_list) the form
// carries `option_value_<optionID>` with a chosen product_option_value ID.
//
// For free-type options (free_text / free_number) the form carries
// `option_free_<optionID>` with the raw user input. We upsert a
// product_option_value row keyed on (option_id, value) — reusing one if it
// exists (the table has a unique_together on those columns) — and link the
// variant to that value. This keeps product_variant_option uniform without
// schema changes.
func saveVariantOptions(ctx context.Context, deps *DetailViewDeps, variantID string, selections []OptionSelection, form map[string][]string) {
	if deps.CreateProductVariantOption == nil {
		return
	}
	selByOption := make(map[string]OptionSelection, len(selections))
	for _, s := range selections {
		selByOption[s.OptionID] = s
	}

	for key, vals := range form {
		if len(vals) == 0 || vals[0] == "" {
			continue
		}
		raw := vals[0]
		var optionID, valueID string

		switch {
		case strings.HasPrefix(key, "option_value_"):
			optionID = strings.TrimPrefix(key, "option_value_")
			valueID = raw
		case strings.HasPrefix(key, "option_free_"):
			optionID = strings.TrimPrefix(key, "option_free_")
			sel, ok := selByOption[optionID]
			if !ok {
				continue
			}
			id, err := upsertFreeOptionValue(ctx, deps, optionID, sel.DataType, raw)
			if err != nil || id == "" {
				log.Printf("Failed to upsert free option_value for option %s: %v", optionID, err)
				continue
			}
			valueID = id
		default:
			continue
		}

		if valueID == "" {
			continue
		}
		_, err := deps.CreateProductVariantOption(ctx, &productvariantoptionpb.CreateProductVariantOptionRequest{
			Data: &productvariantoptionpb.ProductVariantOption{
				ProductVariantId:     variantID,
				ProductOptionValueId: valueID,
				Active:               true,
			},
		})
		if err != nil {
			log.Printf("Failed to create product_variant_option: %v", err)
		}
		_ = optionID // optionID is consumed implicitly via valueID; kept for future audit
	}
}

// upsertFreeOptionValue finds-or-creates a product_option_value row for the
// given option keyed on (option_id, value). For free_number we normalize the
// value field via float parsing so "5", "5.0", and " 5 " collapse to one row.
// Returns the value row's ID.
func upsertFreeOptionValue(ctx context.Context, deps *DetailViewDeps, optionID, dataType, raw string) (string, error) {
	value := strings.TrimSpace(raw)
	label := value
	if dataType == "free_number" {
		// Normalize so float-equivalent inputs share a row.
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			value = strconv.FormatFloat(f, 'f', -1, 64)
		}
	}

	if deps.ListProductOptionValues != nil {
		valResp, err := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
		if err == nil {
			for _, v := range valResp.GetData() {
				if v.GetProductOptionId() == optionID && v.GetValue() == value {
					return v.GetId(), nil
				}
			}
		}
	}

	if deps.CreateProductOptionValue == nil {
		return "", nil
	}
	resp, err := deps.CreateProductOptionValue(ctx, &productoptionvaluepb.CreateProductOptionValueRequest{
		Data: &productoptionvaluepb.ProductOptionValue{
			ProductOptionId: optionID,
			Label:           label,
			Value:           value,
			Active:          true,
		},
	})
	if err != nil {
		return "", err
	}
	if data := resp.GetData(); len(data) > 0 {
		return data[0].GetId(), nil
	}
	return "", nil
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

// validateRequiredOptions checks that all required product options have a value
// supplied. The field name (option_value_* vs option_free_*) is already encoded
// per data type on the OptionSelection, so the check is uniform.
func validateRequiredOptions(selections []OptionSelection, form map[string][]string) string {
	for _, sel := range selections {
		if !sel.Required {
			continue
		}
		vals := form[sel.FieldName]
		if len(vals) == 0 || strings.TrimSpace(vals[0]) == "" {
			return sel.OptionName + " is required"
		}
	}
	return ""
}
