package expense_recognition

// routes.go — expense_recognition entity route constants + Routes config struct. Extracted from the
// domain-level routes.go during the per-entity restructure. Entity-local
// naming (ExpenseRecognition prefix stripped). Pure structural move — route strings are
// byte-identical.

const (
	// ---------------------------------------------------------------------------
	// SPS P10 — ExpenseRecognition + ExpenseRecognitionLine route constants
	// ---------------------------------------------------------------------------

	// ExpenseRecognition master routes (no add/edit drawer — created BY use case)
	ListURL                     = "/expense-recognitions/list/{status}"
	DetailURL                   = "/expense-recognitions/detail/{id}"
	DeleteURL                   = "/action/expense-recognition/delete"
	TabActionURL                = "/action/expense-recognition/detail/{id}/tab/{tab}"
	AttachmentUploadURL         = "/action/expense-recognition/detail/{id}/attachments/upload"
	AttachmentDeleteURL         = "/action/expense-recognition/detail/{id}/attachments/delete"
	ReverseURL                  = "/action/expense-recognition/reverse/{id}"
	RecognizeFromExpenditureURL = "/action/expense-recognition/recognize-from-expenditure"
	RecognizeFromContractURL    = "/action/expense-recognition/recognize-from-contract"

	// ExpenseRecognitionLine routes (child of recognition detail — inline CRUD)
	LineAddURL    = "/action/expense-recognition/{id}/lines/add"
	LineEditURL   = "/action/expense-recognition/{id}/lines/edit/{lid}"
	LineDeleteURL = "/action/expense-recognition/{id}/lines/delete"
)

type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL             string `json:"list_url"`
	DetailURL           string `json:"detail_url"`
	DeleteURL           string `json:"delete_url"`
	TabActionURL        string `json:"tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Workflow
	ReverseURL                  string `json:"reverse_url"`
	RecognizeFromExpenditureURL string `json:"recognize_from_expenditure_url"`
	RecognizeFromContractURL    string `json:"recognize_from_contract_url"`

	// Recognition line actions (child entity — inline CRUD)
	LineAddURL    string `json:"line_add_url"`
	LineEditURL   string `json:"line_edit_url"`
	LineDeleteURL string `json:"line_delete_url"`
}

// DefaultRoutes returns an Routes using the
// package-level URL constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:                   "expense-recognitions",
		ActiveSubNav:                "posted",
		ListURL:                     ListURL,
		DetailURL:                   DetailURL,
		DeleteURL:                   DeleteURL,
		TabActionURL:                TabActionURL,
		AttachmentUploadURL:         AttachmentUploadURL,
		AttachmentDeleteURL:         AttachmentDeleteURL,
		ReverseURL:                  ReverseURL,
		RecognizeFromExpenditureURL: RecognizeFromExpenditureURL,
		RecognizeFromContractURL:    RecognizeFromContractURL,
		LineAddURL:                  LineAddURL,
		LineEditURL:                 LineEditURL,
		LineDeleteURL:               LineDeleteURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"expense_recognition.list":                       r.ListURL,
		"expense_recognition.detail":                     r.DetailURL,
		"expense_recognition.delete":                     r.DeleteURL,
		"expense_recognition.attachment.upload":          r.AttachmentUploadURL,
		"expense_recognition.attachment.delete":          r.AttachmentDeleteURL,
		"expense_recognition.reverse":                    r.ReverseURL,
		"expense_recognition.recognize_from_expenditure": r.RecognizeFromExpenditureURL,
		"expense_recognition.recognize_from_contract":    r.RecognizeFromContractURL,
		"expense_recognition.line.add":                   r.LineAddURL,
		"expense_recognition.line.edit":                  r.LineEditURL,
		"expense_recognition.line.delete":                r.LineDeleteURL,
	}
}
