package inventory

import "github.com/erniealice/pyeza-golang/compose"

// Describe returns the composition-v2 descriptor for the inventory entity.
// Labels are not yet exposed via DefaultLabels() — the LabelJSON binding is
// left empty until a DefaultLabels factory is added.
func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "inventory.inventory",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "inventory"},
		Templates: TemplatesFS,
	}
}
