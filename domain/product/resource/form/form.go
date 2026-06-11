package form

import (
	resource "github.com/erniealice/centymo-golang/domain/product/resource"
)

// Data is the template data for the resource drawer form.
type Data struct {
	FormAction   string
	WorkspaceID  string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	ProductId    string
	UserId       string
	Labels       resource.FormLabels
	CommonLabels any
}
