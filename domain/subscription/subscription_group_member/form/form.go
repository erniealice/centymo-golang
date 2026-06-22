package form

import (
	subscription_group_member "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_member"
)

// Data is the template data for the subscription_group_member drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce       string // CSP nonce for inline <script nonce>
	IsEdit      bool
	ID          string

	// FK fields — stored as plain strings (IDs).
	SubscriptionGroupId string
	SubscriptionId      string
	ClientId            string

	Active bool

	Labels       subscription_group_member.FormLabels
	CommonLabels any
}
