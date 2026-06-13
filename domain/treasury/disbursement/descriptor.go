package disbursement

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "treasury.disbursement",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "treasury_disbursement"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "disbursement.json", Key: "disbursement"},
		LabelName: "DisbursementLabels",
		Templates: TemplatesFS,
	}
}
