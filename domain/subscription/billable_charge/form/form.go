// Package form holds the drawer view models of the billable charge actions.
package form

import (
	"github.com/erniealice/pyeza-golang/types"

	bc "github.com/erniealice/centymo-golang/domain/subscription/billable_charge"
)

// ChargeRow is one selectable open charge in the issue drawer.
type ChargeRow struct {
	ID           string
	Client       string
	Subscription string
	Period       string
	Amount       string
	Checked      bool
}

// IssueData is the template data for the issue-recovery-documents drawer.
type IssueData struct {
	FormAction  string
	WorkspaceID string // injected by the ViewAdapter (action_workspace_guard)

	Labels    bc.Labels
	Series    []types.SelectOption
	Charges   []ChargeRow
	IssueDate string
	DueDate   string
	// IssuanceKey is minted once per drawer render; a double-submit posts the
	// same key and the use case replays the original documents (D8).
	IssuanceKey string
	Summary     string
	NoSeries    bool
	NoCharges   bool

	CommonLabels any
}

// AdjustData is the template data for the adjust-charge drawer.
type AdjustData struct {
	FormAction  string
	WorkspaceID string

	Labels bc.Labels
	ID     string
	Client string
	// CurrentAmount is the display string of the charge's current amount.
	CurrentAmount string
	// NewAmount is the decimal form value (major units).
	NewAmount string
	Currency  string
	Reason    string

	CommonLabels any
}
