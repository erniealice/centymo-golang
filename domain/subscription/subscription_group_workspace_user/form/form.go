package form

import (
	sgwu "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_workspace_user"
)

// Data is the template data for the subscription_group_workspace_user drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce       string // CSP nonce for inline <script nonce>
	IsEdit      bool
	ID          string

	// FK pickers. Each carries the selected id, the selected human label, and
	// the auto-complete option list.
	WorkspaceUserId        string
	WorkspaceUserLabel     string
	WorkspaceUserOpts      []map[string]any
	SubscriptionGroupId    string
	SubscriptionGroupLabel string
	SubscriptionGroupOpts  []map[string]any

	Capacity        string
	CapacityOptions []map[string]any
	IsOwner         bool
	Active          bool

	Labels       sgwu.FormLabels
	CommonLabels any
}

// BuildCapacityOptions returns the capacity <select> options for the drawer
// form. The domain is the closed generic set {primary, access}; an empty or
// unknown selection floors to the least-privilege access value.
func BuildCapacityOptions(l sgwu.FormLabels, selected string) []map[string]any {
	if selected != "primary" {
		selected = "access"
	}
	defs := []struct{ Value, Label string }{
		{"primary", l.CapacityPrimary},
		{"access", l.CapacityAccess},
	}
	opts := make([]map[string]any, 0, len(defs))
	for _, d := range defs {
		opts = append(opts, map[string]any{
			"Value":    d.Value,
			"Label":    d.Label,
			"Selected": d.Value == selected,
		})
	}
	return opts
}

// Pair is a simple id/label option for FK pickers. Disabled renders the option
// greyed out and unselectable — used for FK rows that must stay visible (a
// saved reference keeps its label) but must not be picked anew. Pair lists that
// have no such axis leave it false, keeping every option selectable.
type Pair struct {
	ID       string
	Label    string
	Disabled bool
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
			"Disabled": p.Disabled,
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
