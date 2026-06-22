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
	Scope           string
	Role            string
	IsOwner         bool
	Active          bool

	Labels       pswu.FormLabels
	CommonLabels any
}
