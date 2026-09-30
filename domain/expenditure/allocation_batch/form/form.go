// Package form holds the template data of the allocation drawers (allocate,
// preview panel, read-only view).
package form

import (
	ab "github.com/erniealice/centymo-golang/domain/expenditure/allocation_batch"
	pyeza "github.com/erniealice/pyeza-golang"
)

// Row is one weight input row of the allocation drawer.
type Row struct {
	Kind           string // ALLOCATION_SHARE_KIND_* proto name (posted back verbatim)
	KindLabel      string
	SubscriptionID string
	ClientID       string
	Name           string // participant or fixed-row label (posted back for the preview)
	Numerator      string
	Slug           string // testid / id suffix: subscription id or own-use | vacancy | common-loss
}

// PreviewLine is one line of the split preview / published allocation.
type PreviewLine struct {
	KindLabel string
	Name      string
	Weight    string
	Percent   string
	Amount    string
}

// Preview is the server-computed split shown after Preview / Save draft.
type Preview struct {
	Lines    []PreviewLine
	Total    string
	Source   string
	Currency string
	Balanced bool
}

// PreviewPage is the template data of "allocation-batch-preview".
type PreviewPage struct {
	Preview *Preview
	Labels  ab.Labels
	// Error is the refusal message shown instead of the table (the preview request is a plain
	// button request, so its refusal is rendered in the panel, not as an HX-Error-Message).
	Error string
}

// Data is the template data of "allocation-batch-drawer-form".
type Data struct {
	FormAction  string
	PreviewURL  string
	PublishURL  string
	WorkspaceID string // injected by ViewAdapter.injectWorkspaceID for the action workspace guard

	ComponentID    string
	ComponentLabel string
	ServicePeriod  string
	SourceAmount   string
	Currency       string

	Denominator string
	Rows        []Row
	Preview     *Preview
	HasDraft    bool
	Revision    string
	CanPublish  bool
	// PublishTooltip is the disabled Publish tooltip naming the missing permission code.
	PublishTooltip string

	Labels       ab.Labels
	CommonLabels pyeza.CommonLabels
}

// ViewData is the template data of "allocation-batch-view" (read-only).
type ViewData struct {
	ComponentLabel string
	ServicePeriod  string
	Revision       string
	Status         string
	Preview        *Preview
	Labels         ab.Labels
	CommonLabels   pyeza.CommonLabels
}
