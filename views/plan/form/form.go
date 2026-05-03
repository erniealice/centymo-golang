// Package form owns the template data shape for the primary plan drawer
// (plan-drawer-form.html). Pure types only — no Deps, no context.Context,
// no repository imports.
package form

// Labels holds i18n labels for the drawer form template.
type Labels struct {
	Name            string
	NamePlaceholder string
	Description     string
	DescPlaceholder string
	Active          string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string
	DescriptionInfo string
	ActiveInfo      string

	// 2026-04-27 plan-client-scope plan §6.2 / §6.6 — Client picker.
	Client                  string
	ClientHelp              string
	ClientPlaceholder       string
	ClientSearchPlaceholder string
	ClientNoResults         string
	ClientLockedTooltip     string
	ClientForLabel          string // "For {{.ClientName}}" — read-only badge in client-context entry
	ClientInfo              string

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — JobTemplate select.
	JobTemplate     string
	JobTemplateNone string
	JobTemplateHint string

	// 2026-04-30 cyclic-subscription-jobs plan §9.3 — visits_per_cycle field.
	VisitsPerCycleLabel       string
	VisitsPerCyclePlaceholder string
	VisitsPerCycleHint        string

	// Client-scope cascade notice — rendered unconditionally below the client
	// picker so operators see the schedule restriction before filling other fields.
	ClientScopeCascadeNotice string
}

// ClientFieldMode selects how the Client field renders on the drawer per
// plan §6.6:
//   - "picker"   → standard auto-complete (workspace add).
//   - "readonly" → read-only badge "For {ClientName}" (client-context entry,
//     ?context=client&client_id=...).
//   - "locked"   → read-only badge with the lock tooltip (Plan has active
//     subscriptions and client_id is reference-checker locked).
type ClientFieldMode string

const (
	ClientFieldModePicker   ClientFieldMode = "picker"
	ClientFieldModeReadonly ClientFieldMode = "readonly"
	ClientFieldModeLocked   ClientFieldMode = "locked"
)

// JobTemplateOption is a {value, label} pair for the JobTemplate select.
type JobTemplateOption struct {
	Value string
	Label string
}

// Data is the template data for the plan drawer form.
type Data struct {
	FormAction  string
	IsEdit      bool
	ID          string
	Name        string
	Description string
	Active      bool

	// 2026-04-27 plan-client-scope plan §6.2 / §6.6.
	ClientFieldMode ClientFieldMode
	ClientID        string           // existing or pre-filled client_id
	ClientLabel     string           // display name for the chosen client
	ClientOptions   []map[string]any // optgroup-flattened options for the picker
	SearchClientURL string           // auto-complete search endpoint

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — Plan.job_template_id
	// assignment. JobTemplateID is the currently-assigned id (empty on add /
	// when unset); JobTemplateOptions enumerates active JobTemplates for the
	// drawer's <select>.
	JobTemplateID      string
	JobTemplateOptions []JobTemplateOption

	// 2026-04-30 cyclic-subscription-jobs plan §7.3 / §9.3 — Plan.visits_per_cycle.
	// Number of cycle Job instances spawned per billing cycle. Default 1
	// when unset; the drawer template renders 1 in the input either way.
	// Visible only when JobTemplateID is set (template-side JS gate).
	VisitsPerCycle int32

	Labels       Labels
	CommonLabels any
}
