package action

// NewPriceLookupAction returns an http.HandlerFunc that looks up the applicable
// price for a given product + location + invoice date.
//
// Query params: product_id, location_id, date (YYYY-MM-DD)
//
// Response JSON:
//
//	{"found": false}
//	{"found": true, "price": 1500, "currency": "USD", "price_list_id": "...", "price_product_id": "..."}

import (
	"fmt"
	"log"
	"net/http"

	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"
	priceproductpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_product"
)

// PriceLookupResponse is the JSON shape returned by the price lookup handler.
type PriceLookupResponse struct {
	Found          bool   `json:"found"`
	Price          int64  `json:"price,omitempty"`          // centavos
	Currency       string `json:"currency,omitempty"`
	PriceDisplay   string `json:"price_display,omitempty"`  // e.g. "15.00"
	PriceListID    string `json:"price_list_id,omitempty"`
	PriceProductID string `json:"price_product_id,omitempty"`
}

// NewPriceLookupAction returns an http.HandlerFunc for the price lookup endpoint.
// Given ?product_id=&location_id=&date= it:
//  1. Calls FindApplicablePriceList to find the active price list for the location+date.
//  2. Calls ListPriceProducts and filters by price_list_id + product_id.
//  3. Returns the price + identifiers as JSON.
//
// Returns {"found": false} when any required parameter is missing, no price list
// matches, or no price product is configured for the product.
func NewPriceLookupAction(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		productID := r.URL.Query().Get("product_id")
		locationID := r.URL.Query().Get("location_id")
		date := r.URL.Query().Get("date")

		if productID == "" || locationID == "" || date == "" {
			writeJSON(w, PriceLookupResponse{Found: false})
			return
		}

		if deps.FindApplicablePriceList == nil || deps.ListPriceProducts == nil {
			writeJSON(w, PriceLookupResponse{Found: false})
			return
		}

		// 1. Find the active price list for the given location + date.
		plResp, err := deps.FindApplicablePriceList(ctx, &pricelistpb.FindApplicablePriceListRequest{
			LocationId: locationID,
			Date:       date,
		})
		if err != nil {
			log.Printf("price lookup: FindApplicablePriceList failed: %v", err)
			writeJSON(w, PriceLookupResponse{Found: false})
			return
		}
		if !plResp.GetFound() || plResp.GetPriceList() == nil {
			writeJSON(w, PriceLookupResponse{Found: false})
			return
		}
		priceListID := plResp.GetPriceList().GetId()

		// 2. List price products and find one matching this price_list + product.
		ppResp, err := deps.ListPriceProducts(ctx, &priceproductpb.ListPriceProductsRequest{})
		if err != nil {
			log.Printf("price lookup: ListPriceProducts failed: %v", err)
			writeJSON(w, PriceLookupResponse{Found: false})
			return
		}

		for _, pp := range ppResp.GetData() {
			if pp.GetPriceListId() != priceListID {
				continue
			}
			if pp.GetProductId() != productID {
				continue
			}
			// Found a matching price product.
			writeJSON(w, PriceLookupResponse{
				Found:          true,
				Price:          pp.GetAmount(),
				Currency:       pp.GetCurrency(),
				PriceDisplay:   formatCentavos(pp.GetAmount()),
				PriceListID:    priceListID,
				PriceProductID: pp.GetId(),
			})
			return
		}

		// No price product found for this product in the applicable price list.
		writeJSON(w, PriceLookupResponse{Found: false})
	}
}

// formatCentavos formats an int64 centavo value as a decimal string (e.g. 1500 → "15.00").
func formatCentavos(centavos int64) string {
	if centavos == 0 {
		return "0.00"
	}
	negative := centavos < 0
	if negative {
		centavos = -centavos
	}
	whole := centavos / 100
	frac := centavos % 100
	s := fmt.Sprintf("%d.%02d", whole, frac)
	if negative {
		s = "-" + s
	}
	return s
}
