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

	WorkspaceUserId     string
	SubscriptionGroupId string
	Scope               string
	Role                string
	IsOwner             bool
	Active              bool

	Labels       sgwu.FormLabels
	CommonLabels any
}
