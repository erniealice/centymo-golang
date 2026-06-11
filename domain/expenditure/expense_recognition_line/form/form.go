package form

import expense_recognition "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition"

// Data is the template data for the recognition line drawer form.
type Data struct {
	FormAction           string
	WorkspaceID          string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit               bool
	ID                   string
	ExpenseRecognitionID string
	Description          string
	Quantity             string
	UnitAmount           string
	Amount               string
	Currency             string
	CommonLabels         any
	Labels               expense_recognition.LineLabels
}
