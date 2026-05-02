// Package form owns the template data shape for the subscription backfill-cycles
// drawer (subscription-backfill-cycles-drawer-form.html). Pure types only —
// no Deps, no context.Context, no repository imports.
//
// The spawn_cycle_jobs feature covers two handlers: NewSpawnCycleJobsAction
// (POST-only, no drawer) and NewBackfillCyclesAction (GET drawer + POST commit).
// This form/ package holds the drawer data for the backfill cycles GET path.
package form

// Labels carries the typed strings consumed by the backfill drawer template.
type Labels struct {
	Title       string
	Description string
	CountLabel  string
	Confirm     string
	Cancel      string
	MaxWarning  string
}

// Data is the template shape for subscription-backfill-cycles-drawer-form.html.
// 2026-04-30 cyclic-subscription-jobs plan §5.3.
type Data struct {
	FormAction        string
	SubscriptionID    string
	SubscriptionLabel string

	// MaxCycles caps the number input (see plan §15 risk mitigation —
	// 24 cycles per request).
	MaxCycles int
	// DefaultCycles is the prefilled value (1 = "spawn the next missing").
	DefaultCycles int

	Labels       Labels
	CommonLabels any
}
