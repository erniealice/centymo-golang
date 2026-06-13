package supplier_contract

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.supplier_contract",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "supplier_contract"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "supplier_contract.json", Key: "supplierContract"},
		LabelName: "SupplierContractLabels",
		Templates: TemplatesFS,
	}
}
