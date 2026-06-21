package treasuryadvancesdashboard

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "treasury.treasuryadvancesdashboard",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "treasury_advances"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "advances_dashboard.json", Key: "advancesDashboard"},
		LabelName: "AdvancesDashboardLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "treasury_collection:list|treasury_disbursement:list",
			Items: []compose.NavItem{
				// cash app — Advances section (Plan B Phase 3)
				{Key: "advances-dashboard", Route: "treasury_advances.dashboard",
					Label: "Dashboard", Icon: "icon-credit-card",
					Permission: "treasury_collection:list|treasury_disbursement:list"},
				{Key: "advance-collections", Route: "treasury_advances.advance_collection_list",
					Label: "Advance Collections", Icon: "icon-trending-up",
					Permission: "treasury_collection:list|collection:list"},
				{Key: "advance-disbursements", Route: "treasury_advances.advance_disbursement_list",
					Label: "Advance Disbursements", Icon: "icon-trending-down",
					Permission: "treasury_disbursement:list|disbursement:list"},
			},
		},
	}
}
