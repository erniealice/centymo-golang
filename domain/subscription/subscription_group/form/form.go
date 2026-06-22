package form

import (
	subscription_group "github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
)

// Data is the template data for the subscription_group drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce       string // CSP nonce for inline <script nonce>
	IsEdit      bool
	ID          string

	Name string
	Kind string

	// Anchors — program (plan_id) and period (price_schedule_id) pickers.
	PlanID               string
	PlanLabel            string
	PlanOptions          []map[string]any
	PriceScheduleID      string
	PriceScheduleLabel   string
	PriceScheduleOptions []map[string]any
	KindOptions          []map[string]any
	CapacityMode         string // proto enum string, e.g. "CAPACITY_MODE_CAPPED"
	CapacityModeOptions  []map[string]any
	MaxCapacity          string // rendered as a number input; "" when unset
	Active               bool

	Labels       subscription_group.FormLabels
	CommonLabels any
}

// BuildKindOptions returns the cohort-kind <select> options. The proto stores
// kind as a free-text discriminator ("cohort"/"roster"/"panel"/"project_team");
// these are the canonical choices surfaced in the drawer.
func BuildKindOptions(l subscription_group.FormLabels, selected string) []map[string]any {
	defs := []struct{ Value, Label string }{
		{"cohort", l.KindCohort},
		{"roster", l.KindRoster},
		{"panel", l.KindPanel},
		{"project_team", l.KindProjectTeam},
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

// BuildCapacityModeOptions returns the capacity_mode <select> options. Mirrors
// the CapacityMode proto enum; UNSPECIFIED is omitted (treated as UNLIMITED by
// the proto contract) so the operator picks an explicit mode.
func BuildCapacityModeOptions(l subscription_group.FormLabels, selected string) []map[string]any {
	if selected == "" || selected == "CAPACITY_MODE_UNSPECIFIED" {
		selected = "CAPACITY_MODE_UNLIMITED"
	}
	defs := []struct{ Value, Label string }{
		{"CAPACITY_MODE_UNLIMITED", l.CapUnlimited},
		{"CAPACITY_MODE_CLOSED", l.CapClosed},
		{"CAPACITY_MODE_CAPPED", l.CapCapped},
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

// BuildAutoCompleteOptions converts an id→label map (with a select order slice)
// into the map shape expected by the auto-complete component.
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

// Pair is a simple id/label option for the program + period pickers.
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
