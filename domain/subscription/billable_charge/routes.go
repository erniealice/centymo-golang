package billable_charge

import "strings"

// Billable charge route constants. The proto domain is subscription; the URL
// lives under Revenue (build-spec 6.5). general/route.json carries
// "billable_charge": {} — defaults live here.
const (
	ListURL   = "/revenue/billable-charges/list/{status}"
	TableURL  = "/action/billable-charge/table/{status}"
	IssueURL  = "/action/billable-charge/issue"
	AdjustURL = "/action/billable-charge/adjust/{id}"
)

// Routes holds every billable charge route. JSON tags are the lyngua route.json keys.
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL   string `json:"list_url"`
	TableURL  string `json:"table_url"`
	IssueURL  string `json:"issue_url"`
	AdjustURL string `json:"adjust_url"`
}

// DefaultRoutes returns the route table populated from the constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "revenue",
		ActiveSubNav: "billable-charges",
		ListURL:      ListURL,
		TableURL:     TableURL,
		IssueURL:     IssueURL,
		AdjustURL:    AdjustURL,
	}
}

// RouteMap returns dot-notation route keys.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"billable_charge.list":   r.ListURL,
		"billable_charge.table":  r.TableURL,
		"billable_charge.issue":  r.IssueURL,
		"billable_charge.adjust": r.AdjustURL,
	}
}

// Statuses are the canonical list-status path values (open, issued; cancelled
// is reachable by URL only — voiding a document cancels its charges).
var Statuses = []string{"open", "issued", "cancelled"}

// NormalizeStatus maps an unknown list status to "open".
func NormalizeStatus(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, v := range Statuses {
		if s == v {
			return s
		}
	}
	return "open"
}
