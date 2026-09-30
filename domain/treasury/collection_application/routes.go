package collection_application

// Receive-and-apply route constants (build-spec 6.5). general/route.json carries
// "collection_application": {} — defaults live here. The preview is a GET (a
// read): the drawer's signed action_workspace_guard fields are bound to the one
// POST path the form submits to.
const (
	ReceiveApplyURL = "/action/collection/receive-apply"
	PreviewURL      = "/action/collection/receive-apply/preview"
	ReverseURL      = "/action/collection-application/reverse/{id}"

	// ApplicationsTableID is the applications table (recovery document detail)
	// refreshed after a reversal.
	ApplicationsTableID = "recovery-document-applications-table"
)

// Routes holds every receive-and-apply route. JSON tags are the lyngua route.json keys.
type Routes struct {
	ReceiveApplyURL string `json:"receive_apply_url"`
	PreviewURL      string `json:"preview_url"`
	ReverseURL      string `json:"reverse_url"`
}

// DefaultRoutes returns the route table populated from the constants.
func DefaultRoutes() Routes {
	return Routes{ReceiveApplyURL: ReceiveApplyURL, PreviewURL: PreviewURL, ReverseURL: ReverseURL}
}

// RouteMap returns dot-notation route keys.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"collection_application.receive_apply": r.ReceiveApplyURL,
		"collection_application.preview":       r.PreviewURL,
		"collection_application.reverse":       r.ReverseURL,
	}
}
