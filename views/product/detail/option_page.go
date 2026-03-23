package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
)

// Breadcrumb represents a navigation breadcrumb item.
type Breadcrumb struct {
	Label string
	Href  string
}

// OptionPageData holds data for the option values page.
type OptionPageData struct {
	types.PageData
	ContentTemplate string
	Breadcrumbs     []Breadcrumb
	ProductID       string
	OptionID        string
	OptionName      string
	OptionCode      string
	OptionDataType  string
	OptionStatus    string
	StatusVariant   string
	ValuesTable     *types.TableConfig
	Labels          centymo.ProductLabels
}

// NewOptionPageView creates the option values full-page view.
// Route: /app/products/detail/{id}/options/{oid}
func NewOptionPageView(deps *OptionsDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")
		optionID := viewCtx.Request.PathValue("oid")

		// Load product name for breadcrumb
		l := deps.Labels
		productName := l.Breadcrumb.Product
		if productID != "" && deps.ReadProduct != nil {
			prodResp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
				Data: &productpb.Product{Id: productID},
			})
			if err == nil && len(prodResp.GetData()) > 0 {
				if name := prodResp.GetData()[0].GetName(); name != "" {
					productName = name
				}
			} else if err != nil {
				log.Printf("Failed to read product %s: %v", productID, err)
			}
		}

		// Load option details
		optionName := l.Breadcrumb.Option
		optionCode := ""
		optionDataType := ""
		optionActive := true

		optResp, err := deps.ReadProductOption(ctx, &productoptionpb.ReadProductOptionRequest{
			Data: &productoptionpb.ProductOption{Id: optionID},
		})
		if err != nil || len(optResp.GetData()) == 0 {
			log.Printf("Failed to read product_option %s: %v", optionID, err)
			return view.Error(fmt.Errorf("failed to load product option: %w", err))
		}
		option := optResp.GetData()[0]

		if name := option.GetName(); name != "" {
			optionName = name
		}
		optionCode = option.GetCode()
		optionDataType = option.GetDataType()
		optionActive = option.GetActive()

		// Map data_type to display name
		dataTypeDisplay := dataTypeDisplayName(optionDataType, l.Options.DataTypes)

		// Map active/inactive to status and variant
		optionStatus := "active"
		sVariant := "success"
		if !optionActive {
			optionStatus = "inactive"
			sVariant = "warning"
		}

		// Build breadcrumbs
		breadcrumbs := []Breadcrumb{
			{Label: l.Breadcrumb.Products, Href: route.ResolveURL(deps.Routes.ListURL, "status", "active")},
			{Label: productName, Href: route.ResolveURL(deps.Routes.DetailURL, "id", productID) + "?tab=options"},
			{Label: optionName, Href: ""},
		}

		// Build option values table
		detailDeps := &DetailViewDeps{
			Routes:                  deps.Routes,
			DB:                      deps.DB,
			Labels:                  deps.Labels,
			TableLabels:             deps.TableLabels,
			ListProductOptions:      deps.ListProductOptions,
			ListProductOptionValues: deps.ListProductOptionValues,
		}
		valuesTable := buildOptionValuesTable(ctx, detailDeps, productID, optionID)

		// Type-assert CommonLabels
		commonLabels, _ := deps.CommonLabels.(pyeza.CommonLabels)

		pageData := &OptionPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          optionName,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    optionName,
				HeaderSubtitle: optionCode,
				HeaderIcon:     "icon-settings",
				CommonLabels:   commonLabels,
			},
			ContentTemplate: "option-detail-content",
			Breadcrumbs:     breadcrumbs,
			ProductID:       productID,
			OptionID:        optionID,
			OptionName:      optionName,
			OptionCode:      optionCode,
			OptionDataType:  dataTypeDisplay,
			OptionStatus:    optionStatus,
			StatusVariant:   sVariant,
			ValuesTable:     valuesTable,
			Labels:          l,
		}

		return view.OK("option-detail", pageData)
	})
}
