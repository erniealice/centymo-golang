package accrued_expense

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.accrued_expense",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "accrued_expense"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "accrued_expense.json", Key: "accruedExpense"},
		LabelName: "AccruedExpenseLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "accrued_expense:list",
			Items: []compose.NavItem{
				// expense app — SPS Wave 4 accrued expenses
				{Key: "accrued-expenses", Route: "accrued_expense.list", Params: map[string]string{"status": "outstanding"},
					Label: "Accrued Expenses", Icon: "icon-alert-circle", Permission: "accrued_expense:list"},
			},
		},
	}
}
