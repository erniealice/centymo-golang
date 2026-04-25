package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
)

// OptionFormData is the template data for the option drawer form.
type OptionFormData struct {
	FormAction      string
	IsEdit          bool
	ID              string
	ProductID       string
	Name            string
	Code            string
	Description     string
	DataType        string
	SortOrder       string
	MinValue        string
	MaxValue        string
	Required        bool
	Active          bool
	Labels          centymo.ProductOptionFormLabels
	CommonLabels    any
	DataTypeOptions []types.SelectOption
	// ManageValuesURL points at the option detail page, where values
	// (rows in product_option_value) are added/edited/removed. Only set
	// in Edit mode — the option must exist before its values do.
	ManageValuesURL string
	// IsListType is true when the option's data_type is one whose
	// per-value rows show on the detail page (text_list / number_list /
	// color_list). The "Manage values" link is suppressed for free_text
	// / free_number — those don't have a predefined value set.
	IsListType bool
}

// nextOptionSortOrder returns the smallest unused positive sort_order for the
// given product, suggested as the default when the Add Option drawer opens.
// Returns max(existing) + 1, or 1 if the product has no options yet.
func nextOptionSortOrder(ctx context.Context, deps *OptionsDeps, productID string) int32 {
	if deps == nil || deps.ListProductOptions == nil {
		return 1
	}
	resp, err := deps.ListProductOptions(ctx, &productoptionpb.ListProductOptionsRequest{})
	if err != nil {
		return 1
	}
	var max int32
	for _, o := range resp.GetData() {
		if o == nil || o.GetProductId() != productID {
			continue
		}
		if so := o.GetSortOrder(); so > max {
			max = so
		}
	}
	return max + 1
}

// buildDataTypeOptions builds the list of data type choices for the form select
// dropdown. Returns []types.SelectOption (not a local struct) so the shared
// form-group template can read the Description field it expects on every
// option — form-group.html renders `data-description="{{.Description}}"`
// unconditionally, so missing the field triggers a template execution error.
func buildDataTypeOptions(labels centymo.ProductOptionDataTypeLabels) []types.SelectOption {
	return []types.SelectOption{
		{Value: "text_list", Label: labels.TextList},
		{Value: "number_list", Label: labels.NumberRange},
		{Value: "color_list", Label: labels.ColorList},
		{Value: "free_text", Label: labels.FreeText},
		{Value: "free_number", Label: labels.FreeNumber},
	}
}

// NewOptionAddView creates the option add action (GET = form, POST = create).
func NewOptionAddView(deps *OptionsDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")
		ol := deps.Labels.Options

		if viewCtx.Request.Method == http.MethodGet {
			// Pre-fill sort_order with the next slot above existing options so the
			// user can drop a new option at the end without renumbering. Falls
			// back to "1" if no options exist yet.
			nextSort := nextOptionSortOrder(ctx, deps, productID)
			return view.OK("option-drawer-form", &OptionFormData{
				FormAction:      route.ResolveURL(deps.Routes.OptionAddURL, "id", productID),
				ProductID:       productID,
				SortOrder:       fmt.Sprintf("%d", nextSort),
				Required:        false,
				Active:          true,
				Labels:          ol.Form,
				CommonLabels:    nil, // injected by ViewAdapter
				DataTypeOptions: buildDataTypeOptions(ol.DataTypes),
			})
		}

		// POST — create option
		if err := viewCtx.Request.ParseForm(); err != nil {
			return HtmxError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		required := r.FormValue("required") == "true"

		optionData := &productoptionpb.ProductOption{
			ProductId: productID,
			Name:      r.FormValue("name"),
			Code:      r.FormValue("code"),
			DataType:  r.FormValue("data_type"),
			Required:  required,
			Active:    active,
		}
		if v := r.FormValue("description"); v != "" {
			optionData.Description = &v
		}
		if v := r.FormValue("sort_order"); v != "" {
			if so, err := strconv.ParseInt(v, 10, 32); err == nil {
				optionData.SortOrder = int32(so)
			}
		}
		if v := r.FormValue("min_value"); v != "" {
			if mv, err := strconv.ParseFloat(v, 64); err == nil {
				optionData.MinValue = &mv
			}
		}
		if v := r.FormValue("max_value"); v != "" {
			if mv, err := strconv.ParseFloat(v, 64); err == nil {
				optionData.MaxValue = &mv
			}
		}

		if _, err := deps.CreateProductOption(ctx, &productoptionpb.CreateProductOptionRequest{
			Data: optionData,
		}); err != nil {
			log.Printf("Failed to create product option: %v", err)
			return HtmxError(err.Error())
		}

		return HtmxSuccess("product-options-table")
	})
}

// NewOptionEditView creates the option edit action (GET = form, POST = update).
func NewOptionEditView(deps *OptionsDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")
		optionID := viewCtx.Request.PathValue("oid")
		ol := deps.Labels.Options

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadProductOption(ctx, &productoptionpb.ReadProductOptionRequest{
				Data: &productoptionpb.ProductOption{Id: optionID},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				log.Printf("Failed to read product option %s: %v", optionID, err)
				return HtmxError(deps.Labels.Errors.NotFound)
			}
			record := readResp.GetData()[0]

			name := record.GetName()
			code := record.GetCode()
			description := record.GetDescription()
			dataType := record.GetDataType()
			sortOrder := ""
			if so := record.GetSortOrder(); so != 0 {
				sortOrder = fmt.Sprintf("%d", so)
			}
			minValue := ""
			if mv := record.GetMinValue(); mv != 0 {
				minValue = fmt.Sprintf("%v", mv)
			}
			maxValue := ""
			if mv := record.GetMaxValue(); mv != 0 {
				maxValue = fmt.Sprintf("%v", mv)
			}
			active := record.GetActive()
			required := record.GetRequired()
			isListType := dataType == "text_list" || dataType == "number_list" || dataType == "color_list"

			return view.OK("option-drawer-form", &OptionFormData{
				FormAction:      route.ResolveURL(deps.Routes.OptionEditURL, "id", productID, "oid", optionID),
				IsEdit:          true,
				ID:              optionID,
				ProductID:       productID,
				Name:            name,
				Code:            code,
				Description:     description,
				DataType:        dataType,
				SortOrder:       sortOrder,
				MinValue:        minValue,
				MaxValue:        maxValue,
				Required:        required,
				Active:          active,
				Labels:          ol.Form,
				CommonLabels:    nil, // injected by ViewAdapter
				DataTypeOptions: buildDataTypeOptions(ol.DataTypes),
				ManageValuesURL: route.ResolveURL(deps.Routes.OptionDetailURL, "id", productID, "oid", optionID),
				IsListType:      isListType,
			})
		}

		// POST — update option
		if err := viewCtx.Request.ParseForm(); err != nil {
			return HtmxError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		required := r.FormValue("required") == "true"

		optionData := &productoptionpb.ProductOption{
			Id:       optionID,
			Name:     r.FormValue("name"),
			Code:     r.FormValue("code"),
			DataType: r.FormValue("data_type"),
			Required: required,
			Active:   active,
		}
		if v := r.FormValue("description"); v != "" {
			optionData.Description = &v
		}
		if v := r.FormValue("sort_order"); v != "" {
			if so, err := strconv.ParseInt(v, 10, 32); err == nil {
				optionData.SortOrder = int32(so)
			}
		}
		if v := r.FormValue("min_value"); v != "" {
			if mv, err := strconv.ParseFloat(v, 64); err == nil {
				optionData.MinValue = &mv
			}
		}
		if v := r.FormValue("max_value"); v != "" {
			if mv, err := strconv.ParseFloat(v, 64); err == nil {
				optionData.MaxValue = &mv
			}
		}

		_, err := deps.UpdateProductOption(ctx, &productoptionpb.UpdateProductOptionRequest{
			Data: optionData,
		})
		if err != nil {
			log.Printf("Failed to update product option %s: %v", optionID, err)
			return HtmxError(err.Error())
		}

		return HtmxSuccess("product-options-table")
	})
}

// NewOptionDeleteView creates the option delete action (POST only, with dialog confirmation).
func NewOptionDeleteView(deps *OptionsDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return HtmxError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeleteProductOption(ctx, &productoptionpb.DeleteProductOptionRequest{
			Data: &productoptionpb.ProductOption{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete product option %s: %v", id, err)
			return HtmxError(err.Error())
		}

		return HtmxSuccess("product-options-table")
	})
}
