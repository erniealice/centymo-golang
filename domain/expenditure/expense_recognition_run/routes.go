package expense_recognition_run

// routes.go — expense_recognition_run entity route constants + Routes config struct. Extracted from the
// domain-level routes.go during the per-entity restructure. Entity-local
// naming (ExpenseRecognitionRun prefix stripped). Pure structural move — route strings are
// byte-identical.

const (
	// Expense Recognition Run (buying-side) routes — Plan A 20260517-expense-run.
	QueueURL                 = "/expense-recognition-run/queue"
	QueueTableURL            = "/action/expense-recognition-run/queue/table"
	ListURL                  = "/expense-recognition-run/list/{status}"
	ListTableURL             = "/action/expense-recognition-run/table/{status}"
	DetailURL                = "/expense-recognition-run/detail/{id}"
	DetailTabActionURL       = "/action/expense-recognition-run/detail/{id}/tab/{tab}"
	NewURL                   = "/expense-recognition-run/new"
	GenerateURL              = "/action/expense-recognition-run/generate"
	SubmitBatchURL           = "/action/expense-recognition-run/submit-batch"
	PerSupplierDrawerURL     = "/action/supplier/expense-recognition-run/{id}"
	PerSubscriptionDrawerURL = "/action/supplier-subscription/expense-recognition-run/{id}"
)

type Routes struct {
	// Sidebar navigation context.
	ActiveNav string `json:"active_nav"`

	// Surface B — workspace queue page.
	QueueURL      string `json:"queue_url"`
	QueueTableURL string `json:"queue_table_url"`

	// Surface D — run history list + detail.
	ListURL            string `json:"list_url"`
	ListTableURL       string `json:"list_table_url"`
	DetailURL          string `json:"detail_url"`
	DetailTabActionURL string `json:"detail_tab_action_url"`

	// Action endpoints.
	NewURL         string `json:"new_url"`
	GenerateURL    string `json:"generate_url"`
	SubmitBatchURL string `json:"submit_batch_url"`

	// Surface A — per-supplier drawer (entydad supplier statement tab).
	PerSupplierDrawerURL string `json:"per_supplier_drawer_url"`

	// Surface C — per-SupplierSubscription drawer.
	PerSubscriptionDrawerURL string `json:"per_subscription_drawer_url"`
}

// DefaultRoutes returns Routes
// populated from the package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:                "expense-recognition-run",
		QueueURL:                 QueueURL,
		QueueTableURL:            QueueTableURL,
		ListURL:                  ListURL,
		ListTableURL:             ListTableURL,
		DetailURL:                DetailURL,
		DetailTabActionURL:       DetailTabActionURL,
		NewURL:                   NewURL,
		GenerateURL:              GenerateURL,
		SubmitBatchURL:           SubmitBatchURL,
		PerSupplierDrawerURL:     PerSupplierDrawerURL,
		PerSubscriptionDrawerURL: PerSubscriptionDrawerURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// expense-recognition-run routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"expense_recognition_run.queue":                   r.QueueURL,
		"expense_recognition_run.queue_table":             r.QueueTableURL,
		"expense_recognition_run.list":                    r.ListURL,
		"expense_recognition_run.list_table":              r.ListTableURL,
		"expense_recognition_run.detail":                  r.DetailURL,
		"expense_recognition_run.detail_tab_action":       r.DetailTabActionURL,
		"expense_recognition_run.new":                     r.NewURL,
		"expense_recognition_run.generate":                r.GenerateURL,
		"expense_recognition_run.submit_batch":            r.SubmitBatchURL,
		"expense_recognition_run.per_supplier_drawer":     r.PerSupplierDrawerURL,
		"expense_recognition_run.per_subscription_drawer": r.PerSubscriptionDrawerURL,
	}
}
