package advancesdashboard

// labels.go — advances-dashboard label structs (centymo W5).
//
// The cash-app "Advances Dashboard" labels, extracted verbatim from the
// treasury-domain advance.go into the per-view advancesdashboard package per
// the domain-first restructure. The Advances Dashboard is a domain-level view
// (not an esqyma entity) so the package is disambiguated as advancesdashboard.
// Pure structural move — no behaviour change. Lyngua JSON load paths unchanged.

// AdvancesDashboardTableLabels — column headers for the per-side
// (outflow / inflow) table on the Advances Dashboard. The selling-side
// counterparty header reads "Customer" while the buying-side reads
// "Supplier"; the Defaults functions provide both.
type TableLabels struct {
	ID           string `json:"id"`
	Counterparty string `json:"counterparty"`
	Kind         string `json:"kind"`
	Total        string `json:"total"`
	Remaining    string `json:"remaining"`
	Status       string `json:"status"`
}

// AdvancesDashboardSectionLabels — labels for one half (outflow OR inflow)
// of the Advances Dashboard.
type SectionLabels struct {
	CardTitle    string      `json:"cardTitle"`
	Table        TableLabels `json:"table"`
	EmptyTitle   string      `json:"emptyTitle"`
	EmptyMessage string      `json:"emptyMessage"`
}

// AdvancesDashboardLabels — root struct for the Advances Dashboard page.
// JSON shape mirrors `advancesDashboard.*` in advances_dashboard.json.
type Labels struct {
	Title                string        `json:"title"`
	AsOfLabel            string        `json:"asOfLabel"`
	TotalOutflow         string        `json:"totalOutflow"`
	TotalInflow          string        `json:"totalInflow"`
	UtilizationLabel     string        `json:"utilizationLabel"`
	ActiveCount          string        `json:"activeCount"`
	FullyRecognizedCount string        `json:"fullyRecognizedCount"`
	Outflow              SectionLabels `json:"outflow"`
	Inflow               SectionLabels `json:"inflow"`
}

// DefaultLabels returns English defaults for the cash-app Advances Dashboard.
func DefaultLabels() Labels {
	return Labels{
		Title:                "Advances Dashboard",
		AsOfLabel:            "As of",
		TotalOutflow:         "Total prepaid (asset)",
		TotalInflow:          "Total deferred (liability)",
		UtilizationLabel:     "Utilization",
		ActiveCount:          "Active",
		FullyRecognizedCount: "Fully recognized",
		Outflow: SectionLabels{
			CardTitle: "Outflows (Advance Disbursements)",
			Table: TableLabels{
				ID:           "Advance",
				Counterparty: "Supplier",
				Kind:         "Kind",
				Total:        "Total",
				Remaining:    "Remaining",
				Status:       "Status",
			},
			EmptyTitle:   "No outflow advances",
			EmptyMessage: "Advance disbursements appear here as they are recorded.",
		},
		Inflow: SectionLabels{
			CardTitle: "Inflows (Advance Collections)",
			Table: TableLabels{
				ID:           "Advance",
				Counterparty: "Customer",
				Kind:         "Kind",
				Total:        "Total",
				Remaining:    "Remaining",
				Status:       "Status",
			},
			EmptyTitle:   "No inflow advances",
			EmptyMessage: "Advance collections appear here as they are recorded.",
		},
	}
}
