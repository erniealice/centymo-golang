package allocation_batch

// Cost allocation (allocation_batch) routes. The allocation is worked from the
// expenditure detail "Recoverable costs" tab; {id} is the cost source component id.
const (
	AllocateURL = "/action/allocation-batch/allocate/{id}"
	PreviewURL  = "/action/allocation-batch/preview/{id}"
	PublishURL  = "/action/allocation-batch/publish/{id}"
	ViewURL     = "/action/allocation-batch/view/{id}"
)

// Routes holds every allocation_batch route. JSON tags are the lyngua route.json
// keys (general/route.json carries "allocation_batch": {} — defaults live here).
type Routes struct {
	AllocateURL string `json:"allocate_url"`
	PreviewURL  string `json:"preview_url"`
	PublishURL  string `json:"publish_url"`
	ViewURL     string `json:"view_url"`
}

// DefaultRoutes returns the route table populated from the constants.
func DefaultRoutes() Routes {
	return Routes{AllocateURL: AllocateURL, PreviewURL: PreviewURL, PublishURL: PublishURL, ViewURL: ViewURL}
}

// RouteMap returns the dot-notation route keys.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"allocation_batch.allocate": r.AllocateURL,
		"allocation_batch.preview":  r.PreviewURL,
		"allocation_batch.publish":  r.PublishURL,
		"allocation_batch.view":     r.ViewURL,
	}
}
