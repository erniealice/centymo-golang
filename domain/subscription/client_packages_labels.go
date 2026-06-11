package subscription

// client_packages_labels.go — the client-detail "Packages" tab label set. Its
// name resolves to the entity-domain "client" entity under the placement
// test's mechanical longest-match, but it is a subscription-domain projection
// (the client's Plans/engagements), owned here because the cross-block helper
// view that renders the tab lives in centymo. Excused in placement_test
// legacyAllow by basename pending a W9 naming-resolution pass. centymo W4.

// ClientPackagesLabels holds labels for the client detail "Packages" tab —
// the list of client-scoped Plans for a given client, with the
// "Add custom package" CTA. Mounted from entydad's client detail page via
// a centymo helper view (plan §6.6 option 1).
//
// 2026-04-27 plan-client-scope plan §6.3 / §7.
type ClientPackagesLabels struct {
	TabTitle  string `json:"tabTitle"`
	Empty     string `json:"empty"`
	AddAction string `json:"addAction"`

	// Column headers for the table on the tab.
	ColumnName          string `json:"columnName"`
	ColumnSchedule      string `json:"columnSchedule"`
	ColumnSubscriptions string `json:"columnSubscriptions"`
}

// DefaultClientPackagesLabels returns ClientPackagesLabels with sensible English
// defaults. Surfaces the labels for the client-detail Packages tab + the
// "Add custom package" CTA. Centymo owns these labels because the cross-block
// helper that renders this tab lives here (see plan §6.6 option 1).
//
// 2026-04-27 plan-client-scope plan §7.
func DefaultClientPackagesLabels() ClientPackagesLabels {
	return ClientPackagesLabels{
		TabTitle:            "Packages",
		Empty:               "No custom packages yet — every engagement uses a general package.",
		AddAction:           "Add custom package",
		ColumnName:          "Name",
		ColumnSchedule:      "Rate card",
		ColumnSubscriptions: "Engagements",
	}
}
