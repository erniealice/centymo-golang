package form

import (
	line_workspace_user "github.com/erniealice/centymo-golang/domain/product/line_workspace_user"
)

// Data is the template data for the line_workspace_user drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce       string // CSP nonce for inline <script nonce>
	IsEdit      bool
	ID          string

	WorkspaceUserId string
	LineId          string
	Scope           string
	Role            string
	IsOwner         bool
	Active          bool

	Labels       line_workspace_user.FormLabels
	CommonLabels any
}
