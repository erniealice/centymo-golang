package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	"google.golang.org/protobuf/types/known/structpb"
)

// OptionValueFormData is the template data for the option value drawer form.
type OptionValueFormData struct {
	FormAction     string
	IsEdit         bool
	ID             string
	ProductID      string
	OptionID       string
	OptionName     string
	OptionRequired bool
	Label          string
	Value          string
	SortOrder      string
	ColorHex       string
	Active         bool
	IsColorType    bool
	Labels         centymo.ProductOptionValueFormLabels
	CommonLabels   any
}

// optionInfo holds parent option metadata needed by the value form.
type optionInfo struct {
	Name     string
	Required bool
	IsColor  bool
}

// readOptionInfo loads the parent option's name, required flag, and data type.
func readOptionInfo(ctx context.Context, deps *OptionsDeps, optionID string) optionInfo {
	if deps.ReadProductOption == nil {
		return optionInfo{}
	}
	readResp, err := deps.ReadProductOption(ctx, &productoptionpb.ReadProductOptionRequest{
		Data: &productoptionpb.ProductOption{Id: optionID},
	})
	if err != nil || len(readResp.GetData()) == 0 {
		return optionInfo{}
	}
	o := readResp.GetData()[0]
	return optionInfo{
		Name:     o.GetName(),
		Required: o.GetRequired(),
		IsColor:  o.GetDataType() == "color_list",
	}
}

// NewOptionValueTableView returns a view that renders only the option values table (for HTMX refresh).
func NewOptionValueTableView(deps *OptionsDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")
		optionID := viewCtx.Request.PathValue("oid")

		detailDeps := &Deps{
			Routes:                  deps.Routes,
			DB:                      deps.DB,
			Labels:                  deps.Labels,
			TableLabels:             deps.TableLabels,
			ListProductOptions:      deps.ListProductOptions,
			ListProductOptionValues: deps.ListProductOptionValues,
		}

		tableConfig := buildOptionValuesTable(ctx, detailDeps, productID, optionID)
		return view.OK("table-card", tableConfig)
	})
}

// NewOptionValueAddView creates the option value add action (GET = form, POST = create).
func NewOptionValueAddView(deps *OptionsDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")
		optionID := viewCtx.Request.PathValue("oid")
		ol := deps.Labels.Options

		if viewCtx.Request.Method == http.MethodGet {
			info := readOptionInfo(ctx, deps, optionID)

			return view.OK("option-value-drawer-form", &OptionValueFormData{
				FormAction:     route.ResolveURL(deps.Routes.OptionValueAddURL, "id", productID, "oid", optionID),
				ProductID:      productID,
				OptionID:       optionID,
				OptionName:     info.Name,
				OptionRequired: info.Required,
				Active:         true,
				IsColorType:    info.IsColor,
				Labels:         ol.Value.Form,
				CommonLabels:   nil, // injected by ViewAdapter
			})
		}

		// POST — create option value
		if err := viewCtx.Request.ParseForm(); err != nil {
			return HtmxError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		valData := &productoptionvaluepb.ProductOptionValue{
			ProductOptionId: optionID,
			Label:           r.FormValue("label"),
			Value:           r.FormValue("value"),
			Active:          active,
		}
		if v := r.FormValue("sort_order"); v != "" {
			if so, err := strconv.ParseInt(v, 10, 32); err == nil {
				valData.SortOrder = int32(so)
			}
		}

		// Store color hex in metadata if provided
		colorHex := r.FormValue("color_hex")
		if colorHex != "" {
			meta, err := structpb.NewStruct(map[string]interface{}{"hex": colorHex})
			if err == nil {
				valData.Metadata = meta
			}
		}

		_, err := deps.CreateProductOptionValue(ctx, &productoptionvaluepb.CreateProductOptionValueRequest{
			Data: valData,
		})
		if err != nil {
			log.Printf("Failed to create product option value: %v", err)
			return HtmxError(err.Error())
		}

		return HtmxSuccess("product-option-values-table")
	})
}

// NewOptionValueEditView creates the option value edit action (GET = form, POST = update).
func NewOptionValueEditView(deps *OptionsDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")
		optionID := viewCtx.Request.PathValue("oid")
		valueID := viewCtx.Request.PathValue("vid")
		ol := deps.Labels.Options

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadProductOptionValue(ctx, &productoptionvaluepb.ReadProductOptionValueRequest{
				Data: &productoptionvaluepb.ProductOptionValue{Id: valueID},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				log.Printf("Failed to read product option value %s: %v", valueID, err)
				return HtmxError(deps.Labels.Errors.NotFound)
			}
			record := readResp.GetData()[0]

			label := record.GetLabel()
			value := record.GetValue()
			sortOrder := ""
			if so := record.GetSortOrder(); so != 0 {
				sortOrder = fmt.Sprintf("%d", so)
			}
			active := record.GetActive()

			// Extract color hex from metadata
			colorHex := ""
			if meta := record.GetMetadata(); meta != nil {
				fields := meta.GetFields()
				if hexVal, ok := fields["hex"]; ok {
					colorHex = hexVal.GetStringValue()
				}
			}

			info := readOptionInfo(ctx, deps, optionID)

			return view.OK("option-value-drawer-form", &OptionValueFormData{
				FormAction:     route.ResolveURL(deps.Routes.OptionValueEditURL, "id", productID, "oid", optionID, "vid", valueID),
				IsEdit:         true,
				ID:             valueID,
				ProductID:      productID,
				OptionID:       optionID,
				OptionName:     info.Name,
				OptionRequired: info.Required,
				Label:          label,
				Value:          value,
				SortOrder:      sortOrder,
				ColorHex:       colorHex,
				Active:         active,
				IsColorType:    info.IsColor,
				Labels:         ol.Value.Form,
				CommonLabels:   nil, // injected by ViewAdapter
			})
		}

		// POST — update option value
		if err := viewCtx.Request.ParseForm(); err != nil {
			return HtmxError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		valData := &productoptionvaluepb.ProductOptionValue{
			Id:     valueID,
			Label:  r.FormValue("label"),
			Value:  r.FormValue("value"),
			Active: active,
		}
		if v := r.FormValue("sort_order"); v != "" {
			if so, err := strconv.ParseInt(v, 10, 32); err == nil {
				valData.SortOrder = int32(so)
			}
		}

		// Update color hex in metadata if provided
		colorHex := r.FormValue("color_hex")
		if colorHex != "" {
			meta, err := structpb.NewStruct(map[string]interface{}{"hex": colorHex})
			if err == nil {
				valData.Metadata = meta
			}
		}

		_, err := deps.UpdateProductOptionValue(ctx, &productoptionvaluepb.UpdateProductOptionValueRequest{
			Data: valData,
		})
		if err != nil {
			log.Printf("Failed to update product option value %s: %v", valueID, err)
			return HtmxError(err.Error())
		}

		return HtmxSuccess("product-option-values-table")
	})
}

// NewOptionValueDeleteView creates the option value delete action (POST only, with dialog confirmation).
func NewOptionValueDeleteView(deps *OptionsDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return HtmxError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeleteProductOptionValue(ctx, &productoptionvaluepb.DeleteProductOptionValueRequest{
			Data: &productoptionvaluepb.ProductOptionValue{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete product option value %s: %v", id, err)
			return HtmxError(err.Error())
		}

		return HtmxSuccess("product-option-values-table")
	})
}
