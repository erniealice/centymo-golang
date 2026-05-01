package detail

import (
	"context"
	"fmt"
	"log"
	"sort"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
)

// OptionsDeps holds dependencies for option action handlers.
type OptionsDeps struct {
	Routes       centymo.ProductRoutes
	DB           centymo.DataSource
	Labels       centymo.ProductLabels
	CommonLabels any
	TableLabels  types.TableLabels

	// Typed proto funcs for product_option
	ListProductOptions  func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ReadProductOption   func(ctx context.Context, req *productoptionpb.ReadProductOptionRequest) (*productoptionpb.ReadProductOptionResponse, error)
	CreateProductOption func(ctx context.Context, req *productoptionpb.CreateProductOptionRequest) (*productoptionpb.CreateProductOptionResponse, error)
	UpdateProductOption func(ctx context.Context, req *productoptionpb.UpdateProductOptionRequest) (*productoptionpb.UpdateProductOptionResponse, error)
	DeleteProductOption func(ctx context.Context, req *productoptionpb.DeleteProductOptionRequest) (*productoptionpb.DeleteProductOptionResponse, error)

	// Typed proto funcs for product_option_value
	ListProductOptionValues  func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	ReadProductOptionValue   func(ctx context.Context, req *productoptionvaluepb.ReadProductOptionValueRequest) (*productoptionvaluepb.ReadProductOptionValueResponse, error)
	CreateProductOptionValue func(ctx context.Context, req *productoptionvaluepb.CreateProductOptionValueRequest) (*productoptionvaluepb.CreateProductOptionValueResponse, error)
	UpdateProductOptionValue func(ctx context.Context, req *productoptionvaluepb.UpdateProductOptionValueRequest) (*productoptionvaluepb.UpdateProductOptionValueResponse, error)
	DeleteProductOptionValue func(ctx context.Context, req *productoptionvaluepb.DeleteProductOptionValueRequest) (*productoptionvaluepb.DeleteProductOptionValueResponse, error)

	// Typed proto func for product (option page breadcrumbs)
	ReadProduct func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
}

// NewOptionsTableView returns a view that renders only the options table (for HTMX refresh).
func NewOptionsTableView(deps *OptionsDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")

		// Forward CommonLabels (Yes/No badges live there now). The OptionsDeps
		// keeps a loosely typed `any` to avoid a cyclic typed dep, so we
		// downcast here when forwarding into the typed DetailViewDeps. A
		// missing/empty value is harmless — buildOptionsTable falls back to
		// hardcoded "Yes"/"No" strings.
		var commonLabels pyeza.CommonLabels
		if cl, ok := deps.CommonLabels.(pyeza.CommonLabels); ok {
			commonLabels = cl
		}
		detailDeps := &DetailViewDeps{
			Routes:                  deps.Routes,
			DB:                      deps.DB,
			Labels:                  deps.Labels,
			CommonLabels:            commonLabels,
			TableLabels:             deps.TableLabels,
			ListProductOptions:      deps.ListProductOptions,
			ListProductOptionValues: deps.ListProductOptionValues,
		}

		tableConfig := buildOptionsTable(ctx, detailDeps, productID)
		return view.OK("table-card", tableConfig)
	})
}

// buildOptionsTable builds the options table config for the options tab.
func buildOptionsTable(ctx context.Context, deps *DetailViewDeps, productID string) *types.TableConfig {
	l := deps.Labels
	ol := l.Options

	columns := []types.TableColumn{
		{Key: "sortOrder", Label: ol.Columns.SortOrder, WidthClass: "col-lg"},
		{Key: "name", Label: ol.Columns.Name},
		{Key: "code", Label: ol.Columns.Code},
		{Key: "dataType", Label: ol.Columns.DataType, WidthClass: "col-3xl"},
		{Key: "valuesCount", Label: ol.Columns.ValuesCount, WidthClass: "col-lg"},
		{Key: "required", Label: ol.Columns.Required, WidthClass: "col-lg"},
		{Key: "status", Label: ol.Columns.Status, WidthClass: "col-2xl"},
	}

	yesLabel := deps.CommonLabels.Badges.Yes
	if yesLabel == "" {
		yesLabel = "Yes"
	}
	noLabel := deps.CommonLabels.Badges.No
	if noLabel == "" {
		noLabel = "No"
	}

	rows := []types.TableRow{}

	if deps.ListProductOptions != nil {
		// Load all options
		optResp, err := deps.ListProductOptions(ctx, &productoptionpb.ListProductOptionsRequest{})
		if err != nil {
			log.Printf("Failed to list product options: %v", err)
		} else {
			// Sort by sort_order ASC so the table matches the variant form ordering.
			// Stable sort + Name as tiebreaker keeps zero/equal sort_order rows alphabetical.
			sortedOpts := make([]*productoptionpb.ProductOption, 0, len(optResp.GetData()))
			for _, o := range optResp.GetData() {
				if o != nil && o.GetProductId() == productID {
					sortedOpts = append(sortedOpts, o)
				}
			}
			sort.SliceStable(sortedOpts, func(i, j int) bool {
				if sortedOpts[i].GetSortOrder() != sortedOpts[j].GetSortOrder() {
					return sortedOpts[i].GetSortOrder() < sortedOpts[j].GetSortOrder()
				}
				return sortedOpts[i].GetName() < sortedOpts[j].GetName()
			})
			// Load all option values for counting
			valueCounts := make(map[string]int)
			if deps.ListProductOptionValues != nil {
				valResp, err := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
				if err == nil {
					for _, v := range valResp.GetData() {
						oid := v.GetProductOptionId()
						if oid != "" {
							valueCounts[oid]++
						}
					}
				}
			}

			for _, o := range sortedOpts {
				oid := o.GetId()
				name := o.GetName()
				code := o.GetCode()
				dataType := o.GetDataType()
				sortOrder := ""
				if so := o.GetSortOrder(); so != 0 {
					sortOrder = fmt.Sprintf("%d", so)
				}
				required := o.GetRequired()
				requiredLabel := noLabel
				requiredVariant := "default"
				if required {
					requiredLabel = yesLabel
					requiredVariant = "info"
				}
				active := o.GetActive()
				status := "active"
				if !active {
					status = "inactive"
				}

				dtDisplay := dataTypeDisplayName(dataType, ol.DataTypes)
				vcDisplay := fmt.Sprintf("%d", valueCounts[oid])
				// free_text and free_number have no predefined value list —
				// the option detail page would just show an empty values
				// table, so disable the View row action for those rows.
				hasValues := dataType != "free_text" && dataType != "free_number"

				actions := []types.TableAction{
					{
						Type: "edit", Label: ol.Actions.EditOption, Action: "edit",
						URL:         route.ResolveURL(deps.Routes.OptionEditURL, "id", productID, "oid", oid),
						DrawerTitle: ol.Actions.EditOption,
					},
					{
						Type: "delete", Label: ol.Actions.DeleteOption, Action: "delete",
						URL:            route.ResolveURL(deps.Routes.OptionDeleteURL, "id", productID),
						ItemName:       name,
						ConfirmTitle:   ol.Actions.DeleteOption,
						ConfirmMessage: fmt.Sprintf("%s %s?", ol.Confirm.DeleteOption, name),
					},
					{
						Type: "view", Label: ol.Actions.ViewValues,
						Href:            route.ResolveURL(deps.Routes.OptionDetailURL, "id", productID, "oid", oid),
						Disabled:        !hasValues,
						DisabledTooltip: "This data type does not store predefined values.",
					},
				}

				rows = append(rows, types.TableRow{
					ID: oid,
					Cells: []types.TableCell{
						{Type: "text", Value: sortOrder},
						{Type: "text", Value: name},
						{Type: "text", Value: code},
						{Type: "badge", Value: dtDisplay, Variant: "info"},
						{Type: "text", Value: vcDisplay},
						{Type: "badge", Value: requiredLabel, Variant: requiredVariant},
						{Type: "badge", Value: status, Variant: StatusVariant(status)},
					},
					DataAttrs: map[string]string{
						"name":     name,
						"code":     code,
						"required": fmt.Sprintf("%t", required),
						"status":   status,
					},
					Actions: actions,
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "product-options-table",
		RefreshURL:           route.ResolveURL(deps.Routes.OptionTableURL, "id", productID),
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "sortOrder",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   ol.Empty.Title,
			Message: ol.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     ol.Actions.AddOption,
			ActionURL: route.ResolveURL(deps.Routes.OptionAddURL, "id", productID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// buildOptionValuesTable builds the option values table for a specific option.
func buildOptionValuesTable(ctx context.Context, deps *DetailViewDeps, productID, optionID string) *types.TableConfig {
	l := deps.Labels
	ol := l.Options

	columns := []types.TableColumn{
		{Key: "label", Label: ol.Value.Columns.Label},
		{Key: "value", Label: ol.Value.Columns.Value},
		{Key: "sortOrder", Label: ol.Value.Columns.SortOrder, WidthClass: "col-lg"},
		{Key: "colorPreview", Label: ol.Value.Columns.ColorPreview, NoSort: true, WidthClass: "col-lg"},
		{Key: "status", Label: ol.Value.Columns.Status, WidthClass: "col-2xl"},
	}

	rows := []types.TableRow{}

	// Check if this is a color_list option
	isColor := false
	if deps.ListProductOptions != nil {
		optResp, err := deps.ListProductOptions(ctx, &productoptionpb.ListProductOptionsRequest{})
		if err == nil {
			for _, o := range optResp.GetData() {
				if o.GetId() == optionID {
					isColor = o.GetDataType() == "color_list"
					break
				}
			}
		}
	}

	if deps.ListProductOptionValues != nil {
		valResp, err := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
		if err != nil {
			log.Printf("Failed to list product option values: %v", err)
		} else {
			for _, v := range valResp.GetData() {
				oid := v.GetProductOptionId()
				if oid != optionID {
					continue
				}

				vid := v.GetId()
				label := v.GetLabel()
				value := v.GetValue()
				sortOrder := ""
				if so := v.GetSortOrder(); so != 0 {
					sortOrder = fmt.Sprintf("%d", so)
				}
				active := v.GetActive()
				status := "active"
				if !active {
					status = "inactive"
				}

				// Color swatch preview
				colorPreview := ""
				if isColor {
					if meta := v.GetMetadata(); meta != nil {
						fields := meta.GetFields()
						if hexVal, ok := fields["hex"]; ok {
							hex := hexVal.GetStringValue()
							if hex != "" {
								colorPreview = fmt.Sprintf(`<span class="color-swatch" style="display:inline-block;width:24px;height:24px;border-radius:4px;background:%s;border:1px solid var(--color-border);"></span>`, hex)
							}
						}
					}
				}

				actions := []types.TableAction{
					{
						Type: "edit", Label: ol.Actions.EditValue, Action: "edit",
						URL:         route.ResolveURL(deps.Routes.OptionValueEditURL, "id", productID, "oid", optionID, "vid", vid),
						DrawerTitle: ol.Actions.EditValue,
					},
					{
						Type: "delete", Label: ol.Actions.DeleteValue, Action: "delete",
						URL:            route.ResolveURL(deps.Routes.OptionValueDeleteURL, "id", productID, "oid", optionID),
						ItemName:       label,
						ConfirmTitle:   ol.Actions.DeleteValue,
						ConfirmMessage: fmt.Sprintf("%s %s?", ol.Confirm.DeleteValue, label),
					},
				}

				rows = append(rows, types.TableRow{
					ID: vid,
					Cells: []types.TableCell{
						{Type: "text", Value: label},
						{Type: "text", Value: value},
						{Type: "text", Value: sortOrder},
						{Type: "html", Value: colorPreview},
						{Type: "badge", Value: status, Variant: StatusVariant(status)},
					},
					DataAttrs: map[string]string{
						"label": label,
						"value": value,
					},
					Actions: actions,
				})
			}
		}
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "product-option-values-table",
		RefreshURL:           route.ResolveURL(deps.Routes.OptionValueTableURL, "id", productID, "oid", optionID),
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "label",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   ol.Empty.ValueTitle,
			Message: ol.Empty.ValueMessage,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:     ol.Actions.AddValue,
			ActionURL: route.ResolveURL(deps.Routes.OptionValueAddURL, "id", productID, "oid", optionID),
			Icon:      "icon-plus",
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// getOptionCountTyped returns the number of options for a product (for tab badge) using typed deps.
func getOptionCountTyped(ctx context.Context, deps *DetailViewDeps, productID string) int {
	if deps.ListProductOptions == nil {
		return 0
	}
	optResp, err := deps.ListProductOptions(ctx, &productoptionpb.ListProductOptionsRequest{})
	if err != nil {
		return 0
	}
	count := 0
	for _, o := range optResp.GetData() {
		if o.GetProductId() == productID {
			count++
		}
	}
	return count
}

// dataTypeDisplayName maps a data_type DB value to its human-readable label.
func dataTypeDisplayName(dataType string, labels centymo.ProductOptionDataTypeLabels) string {
	switch dataType {
	case "text_list":
		return labels.TextList
	case "number_list":
		return labels.NumberRange
	case "color_list":
		return labels.ColorList
	case "free_text":
		return labels.FreeText
	case "free_number":
		return labels.FreeNumber
	default:
		return dataType
	}
}
