package form

import (
	plan_group "github.com/erniealice/centymo-golang/domain/product/plan_group"
)

// Data is the template data for the plan_group drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce       string // CSP nonce for inline <script nonce>
	IsEdit      bool
	ID          string

	Name string
	Code string

	// Hierarchy — optional parent_id picker.
	ParentID      string
	ParentLabel   string
	ParentOptions []map[string]any

	Active bool

	Labels       plan_group.FormLabels
	CommonLabels any
}

// BuildParentOptions converts a list of plan groups into the map shape
// expected by the auto-complete component.
func BuildParentOptions(pairs []Pair, selectedID string) []map[string]any {
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

// Pair is a simple id/label option for the parent group picker.
type Pair struct {
	ID    string
	Label string
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
