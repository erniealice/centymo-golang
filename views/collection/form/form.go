package form

import (
	centymo "github.com/erniealice/centymo-golang"
)

// Data is the template data for the collection drawer form.
//
// 20260517-advance-cash-events Plan B Phase 4 — also carries the advance_kind
// + advance_proration_policy fields so the drawer can render the conditional
// proration-policy dropdown (visible only when advance_kind == TIME_BASED).
// The values default to empty (no chip selected = NONE).
type Data struct {
	FormAction             string
	WorkspaceID             string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit                 bool
	ID                     string
	Customer               string
	ReferenceNumber        string
	Amount                 string
	Currency               string
	CollectionMethod       string
	Date                   string
	ReceivedBy             string
	ReceivedRole           string
	Notes                  string
	CollectionType         string
	Status                 string
	AdvanceKind            string
	AdvanceProrationPolicy string
	Labels                 centymo.CollectionFormLabels
	// EnumLabels carries the AdvanceKind / AdvanceProrationPolicy option
	// labels (loaded from advance_kind.json) so the drawer-form template can
	// render the dropdown options without hardcoded English.
	// 20260517-advance-cash-events Plan B Phase 4.
	EnumLabels   centymo.AdvanceEnumLabels
	CommonLabels any
}
