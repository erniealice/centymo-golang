package form

import (
	sgpps "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan_staff"
)

// Data is the template data for the subscription_group_product_plan_staff
// drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce       string // CSP nonce for inline <script nonce>
	IsEdit      bool
	ID          string

	// FK pickers
	SubscriptionGroupID    string
	SubscriptionGroupLabel string
	SubscriptionGroupOpts  []map[string]any
	ProductPlanID          string
	ProductPlanLabel       string
	ProductPlanOpts        []map[string]any
	StaffID                string
	StaffLabel             string
	StaffOpts              []map[string]any

	Role   string
	Active bool

	Labels       sgpps.FormLabels
	CommonLabels any
}

// Pair is a simple id/label option for FK pickers.
type Pair struct {
	ID    string
	Label string
}

// BuildAutoCompleteOptions converts pairs into the map shape expected by the
// auto-complete component.
func BuildAutoCompleteOptions(pairs []Pair, selectedID string) []map[string]any {
	opts := make([]map[string]any, 0, len(pairs))
	for _, p := range pairs {
		opts = append(opts, map[string]any{
			"Value":    p.ID,
			"Label":    p.Label,
			"Selected": p.ID == selectedID,
		})
	}
	return opts
}

// FindLabel returns the label of the pair with the given ID, or "" if absent.
func FindLabel(pairs []Pair, id string) string {
	for _, p := range pairs {
		if p.ID == id {
			return p.Label
		}
	}
	return ""
}
