package expense_recognition

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.expense_recognition",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "expense_recognition"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "expense_recognition.json", Key: "expenseRecognition"},
		LabelName: "ExpenseRecognitionLabels",
		Templates: TemplatesFS,
	}
}
