package cost_source_component

// Recoverable-cost (cost_source_component) routes. The components are managed
// from the expenditure detail "Recoverable costs" tab, so every route hangs off
// /action/cost-source-component/. {id} is the expenditure id on add/table and the
// component id on edit/delete.
const (
	AddURL    = "/action/cost-source-component/add/{id}"
	EditURL   = "/action/cost-source-component/edit/{id}"
	DeleteURL = "/action/cost-source-component/delete/{id}"
	TableURL  = "/action/cost-source-component/table/{id}"
)

// Routes holds every cost_source_component route. JSON tags are the lyngua
// route.json keys (general/route.json carries "cost_source_component": {} —
// the defaults live here).
type Routes struct {
	AddURL    string `json:"add_url"`
	EditURL   string `json:"edit_url"`
	DeleteURL string `json:"delete_url"`
	TableURL  string `json:"table_url"`
}

// DefaultRoutes returns the route table populated from the constants.
func DefaultRoutes() Routes {
	return Routes{AddURL: AddURL, EditURL: EditURL, DeleteURL: DeleteURL, TableURL: TableURL}
}

// RouteMap returns the dot-notation route keys.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"cost_source_component.add":    r.AddURL,
		"cost_source_component.edit":   r.EditURL,
		"cost_source_component.delete": r.DeleteURL,
		"cost_source_component.table":  r.TableURL,
	}
}
