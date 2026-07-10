package price_list

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the composition-v2 descriptor for the price_list entity.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "product.price_list",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "price_list"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "pricelist.json", Key: "pricelist"},
		LabelName: "PriceListLabels",
		Templates: TemplatesFS,
	}
}
