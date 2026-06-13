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
		Nav: compose.NavContrib{
			Permission: "supplier_contract:list",
			Items: []compose.NavItem{
				// supplier app — Supplier Contracts section (master data)
				{Key: "contracts-active", Route: "supplier_contract.list", Params: map[string]string{"status": "active"},
					Label: "Active", Icon: "icon-check-circle", Permission: "supplier_contract:list"},
				{Key: "contracts-expiring", Route: "supplier_contract.list", Params: map[string]string{"status": "expiring"},
					Label: "Expiring", Icon: "icon-alert-triangle", Permission: "supplier_contract:list"},
				{Key: "contracts-pending", Route: "supplier_contract.list", Params: map[string]string{"status": "pending_approval"},
					Label: "Pending Approval", Icon: "icon-clock", Permission: "supplier_contract:list"},
				{Key: "contracts-draft", Route: "supplier_contract.list", Params: map[string]string{"status": "draft"},
					Label: "Draft", Icon: "icon-file-text", Permission: "supplier_contract:list"},
				{Key: "contracts-terminated", Route: "supplier_contract.list", Params: map[string]string{"status": "terminated"},
					Label: "Terminated", Icon: "icon-x-circle", Permission: "supplier_contract:list"},
			},
		},
	}
}
