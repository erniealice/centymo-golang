package procurement_request

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.procurement_request",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "procurement_request"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "procurement_request.json", Key: "procurementRequest"},
		LabelName: "ProcurementRequestLabels",
		Templates: TemplatesFS,
	}
}
