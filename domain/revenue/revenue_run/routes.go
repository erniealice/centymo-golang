package revenuerun

// Default route constants for revenue-run views.
// Consumer apps can use these or define their own via lyngua route.json overrides.
const (
	// Revenue Run (invoice-run) routes
	QueueURL            = "/revenue-run/queue"
	QueueTableURL       = "/action/revenue-run/queue/table"
	ListURL             = "/revenue-run/list/{status}"
	ListTableURL        = "/action/revenue-run/table/{status}"
	DetailURL           = "/revenue-run/detail/{id}"
	DetailTabActionURL  = "/action/revenue-run/detail/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/revenue-run/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/revenue-run/detail/{id}/attachments/delete"
	SubmitBatchURL      = "/action/revenue-run/submit-batch"
)

// Routes holds all route paths for the Revenue Run (invoice-run) module.
// Surface B = workspace queue page; Surface D = run history list + detail pages.
type Routes struct {
	// Sidebar navigation context — set via defaults or routes.json override.
	ActiveNav string `json:"active_nav"`

	QueueURL            string `json:"queue_url"`
	QueueTableURL       string `json:"queue_table_url"`
	ListURL             string `json:"list_url"`
	ListTableURL        string `json:"list_table_url"`
	DetailURL           string `json:"detail_url"`
	DetailTabActionURL  string `json:"detail_tab_action_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
	SubmitBatchURL      string `json:"submit_batch_url"`
}

// DefaultRoutes returns a Routes populated from the
// package-level route constants defined in routes.go.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:           "revenue-run",
		QueueURL:            QueueURL,
		QueueTableURL:       QueueTableURL,
		ListURL:             ListURL,
		ListTableURL:        ListTableURL,
		DetailURL:           DetailURL,
		DetailTabActionURL:  DetailTabActionURL,
		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,
		SubmitBatchURL:      SubmitBatchURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths for all
// revenue-run routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"revenue_run.queue":             r.QueueURL,
		"revenue_run.queue_table":       r.QueueTableURL,
		"revenue_run.list":              r.ListURL,
		"revenue_run.list_table":        r.ListTableURL,
		"revenue_run.detail":            r.DetailURL,
		"revenue_run.detail_tab_action": r.DetailTabActionURL,
		"revenue_run.attachment.upload": r.AttachmentUploadURL,
		"revenue_run.attachment.delete": r.AttachmentDeleteURL,
		"revenue_run.submit_batch":      r.SubmitBatchURL,
	}
}
