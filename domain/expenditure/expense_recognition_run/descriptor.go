package expense_recognition_run

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.expense_recognition_run",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "expense_recognition_run"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "expense_recognition_run.json", Key: "expenseRecognitionRun"},
		LabelName: "ExpenseRecognitionRunLabels",
		Templates: TemplatesFS,
	}
}
