package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
)

// DataTypeOption represents a single data type choice for the option form select.
type DataTypeOption struct {
	Value string
	Label string
}

// OptionFormData is the template data for the option drawer form.
type OptionFormData struct {
	FormAction      string
	IsEdit          bool
	ID              string
	ProductID       string
	Name            string
	Code            string
	DataType        string
	SortOrder       string
	MinValue        string
	MaxValue        string
	Required        bool
	Active          bool
	Labels          centymo.ProductOptionFormLabels
	CommonLabels    any
	DataTypeOptions []DataTypeOption
}

// buildDataTypeOptions builds the list of data type choices for the form select dropdown.
func buildDataTypeOptions(labels centymo.ProductOptionDataTypeLabels) []DataTypeOption {
	return []DataTypeOption{
		{Value: "text_list", Label: labels.TextList},
		{Value: "number_list", Label: labels.NumberList},
		{Value: "color_list", Label: labels.ColorList},
		{Value: "enum_list", Label: labels.EnumList},
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
			return view.OK("option-drawer-form", &OptionFormData{
				FormAction:      route.ResolveURL(deps.Routes.OptionAddURL, "id", productID),
				ProductID:       productID,
				Required:        false,
				Active:          true,
				Labels:          ol.Form,
				CommonLabels:    nil, // injected by ViewAdapter
				DataTypeOptions: buildDataTypeOptions(ol.DataTypes),
			})
		}

		// POST — create option
		if err := viewCtx.Request.ParseForm(); err != nil {
			return HtmxError("Invalid form data")
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

		createResp, err := deps.CreateProductOption(ctx, &productoptionpb.CreateProductOptionRequest{
			Data: optionData,
		})
		if err != nil {
			log.Printf("Failed to create product option: %v", err)
			return HtmxError("Failed to create option")
		}

		// Create initial option values from comma-separated input
		if created := createResp.GetData(); len(created) > 0 {
			optionID := created[0].GetId()
			if optionID != "" {
				if raw := r.FormValue("initial_values"); raw != "" {
					for i, part := range strings.Split(raw, ",") {
						label := strings.TrimSpace(part)
						if label == "" {
							continue
						}
						value := strings.ToLower(strings.ReplaceAll(label, " ", "_"))
						_, err := deps.CreateProductOptionValue(ctx, &productoptionvaluepb.CreateProductOptionValueRequest{
							Data: &productoptionvaluepb.ProductOptionValue{
								ProductOptionId: optionID,
								Label:           label,
								Value:           value,
								SortOrder:       int32(i + 1),
								Active:          true,
							},
						})
						if err != nil {
							log.Printf("Failed to create option value %q: %v", label, err)
						}
					}
				}
			}
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
				return HtmxError("Option not found")
			}
			record := readResp.GetData()[0]

			name := record.GetName()
			code := record.GetCode()
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

			return view.OK("option-drawer-form", &OptionFormData{
				FormAction:      route.ResolveURL(deps.Routes.OptionEditURL, "id", productID, "oid", optionID),
				IsEdit:          true,
				ID:              optionID,
				ProductID:       productID,
				Name:            name,
				Code:            code,
				DataType:        dataType,
				SortOrder:       sortOrder,
				MinValue:        minValue,
				MaxValue:        maxValue,
				Required:        required,
				Active:          active,
				Labels:          ol.Form,
				CommonLabels:    nil, // injected by ViewAdapter
				DataTypeOptions: buildDataTypeOptions(ol.DataTypes),
			})
		}

		// POST — update option
		if err := viewCtx.Request.ParseForm(); err != nil {
			return HtmxError("Invalid form data")
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
			return HtmxError("Failed to update option")
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
			return HtmxError("Option ID is required")
		}

		_, err := deps.DeleteProductOption(ctx, &productoptionpb.DeleteProductOptionRequest{
			Data: &productoptionpb.ProductOption{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete product option %s: %v", id, err)
			return HtmxError("Failed to delete option")
		}

		return HtmxSuccess("product-options-table")
	})
}
