package form

import (
	product_plan_staff "github.com/erniealice/centymo-golang/domain/product/product_plan_staff"
)

// Data is the template data for the product_plan_staff drawer form.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce       string // CSP nonce for inline <script nonce>
	IsEdit      bool
	ID          string

	// Core fields. StaffId and ProductPlanId are required FKs (non-optional in
	// the proto — always string, not pointer). Role is a free-text discriminator.
	StaffID       string
	ProductPlanID string
	Role          string
	Active        bool

	Labels       product_plan_staff.FormLabels
	CommonLabels any
}
