package form

import (
	pswu "github.com/erniealice/centymo-golang/domain/subscription/price_schedule_workspace_user"
)

// Data is the template data for the price_schedule_workspace_user drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce       string // CSP nonce for inline <script nonce>
	IsEdit      bool
	ID          string

	PriceScheduleId string
	WorkspaceUserId string
	Capacity        string
	CapacityOptions []map[string]any
	IsOwner         bool
	Active          bool

	Labels       pswu.FormLabels
	CommonLabels any
}

// BuildCapacityOptions returns the capacity <select> options for the drawer
// form. The domain is the closed generic set {primary, access}; an empty or
// unknown selection floors to the least-privilege access value.
func BuildCapacityOptions(l pswu.FormLabels, selected string) []map[string]any {
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
