package recovery_document

import "strings"

// Recovery document route constants (build-spec 6.5). general/route.json
// carries "recovery_document": {} — defaults live here.
const (
	ListURL      = "/revenue/recovery-documents/list/{status}"
	TableURL     = "/action/recovery-document/table/{status}"
	DetailURL    = "/revenue/recovery-documents/detail/{id}"
	TabActionURL = "/action/recovery-document/detail/{id}/tab/{tab}"
	VoidURL      = "/action/recovery-document/void/{id}"
)

// Routes holds every recovery document route. JSON tags are the lyngua route.json keys.
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL      string `json:"list_url"`
	TableURL     string `json:"table_url"`
	DetailURL    string `json:"detail_url"`
	TabActionURL string `json:"tab_action_url"`
	VoidURL      string `json:"void_url"`
}

// DefaultRoutes returns the route table populated from the constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "revenue",
		ActiveSubNav: "recovery-documents",
		ListURL:      ListURL,
		TableURL:     TableURL,
		DetailURL:    DetailURL,
		TabActionURL: TabActionURL,
		VoidURL:      VoidURL,
	}
}

// RouteMap returns dot-notation route keys.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"recovery_document.list":       r.ListURL,
		"recovery_document.table":      r.TableURL,
		"recovery_document.detail":     r.DetailURL,
		"recovery_document.tab_action": r.TabActionURL,
		"recovery_document.void":       r.VoidURL,
	}
}

// Statuses are the canonical list-status path values.
var Statuses = []string{"issued", "void"}

// NormalizeStatus maps an unknown list status to "issued".
func NormalizeStatus(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, v := range Statuses {
		if s == v {
			return s
		}
	}
	return "issued"
}
