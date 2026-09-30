package collection_application

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the compose descriptor for the receive-and-apply flow. It
// has no sidebar entry: the drawer opens from the collection list's primary
// action and the recovery document detail. The unit is mounted by centymo's
// CollectionApplicationUnit only when the app opts in (block.WithRecoveryCharges).
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "treasury.collection_application",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "collection_application"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "collection_application.json", Key: "collection_application"},
		LabelName: "CollectionApplicationLabels",
		Templates: TemplatesFS,
	}
}
