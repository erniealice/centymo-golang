package price_list

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the composition-v2 descriptor for the price_list entity.
// Labels are not yet exposed via DefaultLabels() — the LabelJSON binding is
// left empty until a DefaultLabels factory is added.
func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "product.price_list",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "price_list"},
		Templates: TemplatesFS,
	}
}
