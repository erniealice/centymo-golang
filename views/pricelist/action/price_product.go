package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	priceproductpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_product"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
)

// ProductOption represents a product in the dropdown.
type ProductOption struct {
	ID   string
	Name string
}

// PriceProductFormData is the template data for the price product add drawer.
type PriceProductFormData struct {
	FormAction   string
	PriceListID  string
	Products     []ProductOption
	CommonLabels any
}

// PriceProductDeps holds dependencies for price product action handlers.
type PriceProductDeps struct {
	CreatePriceProduct func(ctx context.Context, req *priceproductpb.CreatePriceProductRequest) (*priceproductpb.CreatePriceProductResponse, error)
	DeletePriceProduct func(ctx context.Context, req *priceproductpb.DeletePriceProductRequest) (*priceproductpb.DeletePriceProductResponse, error)
	ListProducts       func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
}

// NewPriceProductAddAction creates the price product add action (GET = form, POST = create).
func NewPriceProductAddAction(deps *PriceProductDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		priceListID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			// Load product list for dropdown
			var products []ProductOption
			if deps.ListProducts != nil {
				resp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
				if err != nil {
					log.Printf("Failed to list products for price product form: %v", err)
				} else {
					for _, p := range resp.GetData() {
						if p.GetActive() {
							products = append(products, ProductOption{
								ID:   p.GetId(),
								Name: p.GetName(),
							})
						}
					}
				}
			}

			return view.OK("price-product-drawer-form", &PriceProductFormData{
				FormAction:   fmt.Sprintf("/action/price-lists/%s/products/add", priceListID),
				PriceListID:  priceListID,
				Products:     products,
				CommonLabels: nil,
			})
		}

		// POST -- create price product
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		productID := r.FormValue("product_id")
		name := r.FormValue("name")
		currency := r.FormValue("currency")
		amountStr := r.FormValue("amount")

		if productID == "" {
			return centymo.HTMXError("Product is required")
		}

		var amount int64
		if amountStr != "" {
			a, err := strconv.ParseInt(amountStr, 10, 64)
			if err != nil {
				return centymo.HTMXError("Amount must be a valid number")
			}
			amount = a
		}

		_, err := deps.CreatePriceProduct(ctx, &priceproductpb.CreatePriceProductRequest{
			Data: &priceproductpb.PriceProduct{
				ProductId:   productID,
				Name:        name,
				Amount:      amount,
				Currency:    currency,
				PriceListId: &priceListID,
				Active:      true,
			},
		})
		if err != nil {
			log.Printf("Failed to create price product: %v", err)
			return centymo.HTMXError("Failed to add product price")
		}

		return centymo.HTMXSuccess("price-products-table")
	})
}

// NewPriceProductDeleteAction creates the price product delete action (POST only).
func NewPriceProductDeleteAction(deps *PriceProductDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError("Price product ID is required")
		}

		_, err := deps.DeletePriceProduct(ctx, &priceproductpb.DeletePriceProductRequest{
			Data: &priceproductpb.PriceProduct{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete price product %s: %v", id, err)
			return centymo.HTMXError("Failed to remove product price")
		}

		return centymo.HTMXSuccess("price-products-table")
	})
}
