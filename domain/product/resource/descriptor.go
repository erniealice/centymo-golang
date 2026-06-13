package resource

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "product.resource",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "resource"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "resource.json", Key: "resource"},
		LabelName: "ResourceLabels",
		Templates: TemplatesFS,
	}
}
