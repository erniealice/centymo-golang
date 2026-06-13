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
		Nav: compose.NavContrib{
			Permission: "supplier_contract_price_schedule:list",
			Items: []compose.NavItem{
				// supplier app — sub-link beneath Supplier Contracts (SPS Wave 4)
				{Key: "contracts-price-schedules", Route: "supplier_contract_price_schedule.list",
					Params: map[string]string{"status": "scheduled"},
					Label: "Price Schedules", Icon: "icon-calendar",
					Permission: "supplier_contract_price_schedule:list"},
			},
		},
	}
}
