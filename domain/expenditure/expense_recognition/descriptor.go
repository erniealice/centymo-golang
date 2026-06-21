package expense_recognition

import "github.com/erniealice/espyna-golang/consumer/compose"

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
		Nav: compose.NavContrib{
			Permission: "expense_recognition:list",
			Items: []compose.NavItem{
				// expense app — SPS Wave 4 accrual-basis recognized cost
				{Key: "expense-recognitions", Route: "expense_recognition.list", Params: map[string]string{"status": "posted"},
					Label: "Expense Recognitions", Icon: "icon-file-text", Permission: "expense_recognition:list"},
			},
		},
	}
}
