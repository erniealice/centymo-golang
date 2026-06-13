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
		Nav: compose.NavContrib{
			Permission: "disbursement:list",
			Items: []compose.NavItem{
				// cash app — Disbursements section
				{Key: "disbursements-draft", Route: "disbursement.list", Params: map[string]string{"status": "draft"},
					Label: "Draft", Icon: "icon-file-text", Permission: "disbursement:list"},
				{Key: "disbursements-pending", Route: "disbursement.list", Params: map[string]string{"status": "pending"},
					Label: "Pending", Icon: "icon-clock", Permission: "disbursement:list"},
				{Key: "disbursements-approved", Route: "disbursement.list", Params: map[string]string{"status": "approved"},
					Label: "Approved", Icon: "icon-check-circle", Permission: "disbursement:list"},
				{Key: "disbursements-paid", Route: "disbursement.list", Params: map[string]string{"status": "paid"},
					Label: "Paid", Icon: "icon-dollar-sign", Permission: "disbursement:list"},
			},
		},
	}
}
