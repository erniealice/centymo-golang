package form

import (
	plan_group_plan "github.com/erniealice/centymo-golang/domain/product/plan_group_plan"
)

// Data is the template data for the plan_group_plan drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce       string // CSP nonce for inline <script nonce>
	IsEdit      bool
	ID          string

	PlanGroupID   string
	PlanID        string
	SequenceOrder string // rendered as a number input; "" when unset (optional)
	Active        bool

	Labels       plan_group_plan.FormLabels
	CommonLabels any
}
