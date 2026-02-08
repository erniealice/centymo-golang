package detail

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
)

// Deps holds view dependencies.
type Deps struct {
	ReadProduct func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
}

// PageData holds the data for the product detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Product         *productpb.Product
	ActiveTab       string
	Tabs            []TabConfig
}

// TabConfig defines a single tab in the detail view.
type TabConfig struct {
	Key    string
	Label  string
	Active bool
	URL    string
}

// NewView creates the product detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		tab := viewCtx.Request.URL.Query().Get("tab")
		if tab == "" {
			tab = "basic"
		}

		resp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
			Data: &productpb.Product{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read product %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load product: %w", err))
		}

		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("product not found"))
		}
		product := data[0]

		name := product.GetName()
		description := product.GetDescription()

		tabs := []TabConfig{
			{Key: "basic", Label: "Basic Information", Active: tab == "basic", URL: fmt.Sprintf("/app/products/%s?tab=basic", id)},
			{Key: "history", Label: "History", Active: tab == "history", URL: fmt.Sprintf("/app/products/%s?tab=history", id)},
			{Key: "depreciation", Label: "Depreciation", Active: tab == "depreciation", URL: fmt.Sprintf("/app/products/%s?tab=depreciation", id)},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          name,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "products",
				HeaderTitle:    name,
				HeaderSubtitle: description,
				HeaderIcon:     "icon-package",
			},
			ContentTemplate: "product-detail-content",
			Product:         product,
			ActiveTab:       tab,
			Tabs:            tabs,
		}

		return view.OK("product-detail", pageData)
	})
}
