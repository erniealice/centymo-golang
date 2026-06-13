package collection

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "treasury.collection",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "treasury_collection"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "collection.json", Key: "collection"},
		LabelName: "CollectionLabels",
		Templates: TemplatesFS,
	}
}
