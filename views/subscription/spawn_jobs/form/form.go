// Package form owns the template data shape for the subscription spawn-jobs
// drawer (subscription-spawn-jobs-drawer-form.html). Pure types only — no Deps,
// no context.Context, no repository imports.
package form

// TemplateRow is one detected JobTemplate rendered in the drawer.
type TemplateRow struct {
	TemplateID   string
	TemplateName string
	IsRoot       bool
	PhaseCount   int
	TaskCount    int
}

// Labels carries the typed strings consumed by the drawer template.
// 2026-04-29 auto-spawn-jobs-from-subscription plan §9.3.
type Labels struct {
	Title             string
	DetectedTemplates string
	RootTemplate      string
	Cancel            string
	Confirm           string
	Skipped           string
}

// Data is the template shape for subscription-spawn-jobs-drawer-form.html.
// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.3.
type Data struct {
	FormAction        string
	SubscriptionID    string
	SubscriptionLabel string
	// Detected templates (root + active children). RootName highlights which
	// template is the root for the operator radio group. Empty Templates =
	// nothing to spawn (operator sees the skipped notice).
	Templates  []TemplateRow
	RootName   string
	HasContent bool

	// Resolved labels for the drawer (avoid label drift across tiers).
	Labels Labels
	// Common labels supplied by the view adapter (cancel/save buttons).
	CommonLabels any
}
