package accrued_expense

import "github.com/erniealice/pyeza-golang/compose"

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
	}
}
