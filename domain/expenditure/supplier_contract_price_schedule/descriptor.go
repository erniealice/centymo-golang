package supplier_contract_price_schedule

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.supplier_contract_price_schedule",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "supplier_contract_price_schedule"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "supplier_contract_price_schedule.json", Key: "supplierContractPriceSchedule"},
		LabelName: "SupplierContractPriceScheduleLabels",
		Templates: TemplatesFS,
	}
}
