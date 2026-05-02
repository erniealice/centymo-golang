// Package form owns the template data shape for the primary subscription
// drawer (subscription-drawer-form.html). Pure types only — no Deps, no
// context.Context, no repository imports.
package form

// Labels holds i18n labels for the subscription drawer form template.
type Labels struct {
	Customer                  string
	CustomerPlaceholder       string
	Plan                      string
	PlanPlaceholder           string
	StartDate                 string
	EndDate                   string
	StartTime                 string
	EndTime                   string
	TimePlaceholder           string
	Timezone                  string
	Notes                     string
	NotesPlaceholder          string
	CustomerSearchPlaceholder string
	PlanSearchPlaceholder     string
	CustomerNoResults         string
	PlanNoResults             string
	Code                      string
	CodePlaceholder           string
	CustomerInfo              string
	PlanInfo                  string
	CodeInfo                  string
	StartDateInfo             string
	EndDateInfo               string
	StartTimeInfo             string
	EndTimeInfo               string
	NotesInfo                 string

	// 2026-04-27 plan-client-scope plan §5.1 / §7 — group headers in the
	// grouped Plan / PricePlan auto-complete picker.
	PlanGroupForClient string // "For {ClientName}" — pre-resolved with ClientName injected.
	PlanGroupGeneral   string

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.1 / §9.
	SpawnJobsSectionTitle string
	SpawnJobsToggle       string
	SpawnJobsHelpText     string
	SpawnJobsSummary      string // {{.JobCount}} / {{.TemplateNames}} / {{.PhaseCount}} / {{.TaskCount}}
	SpawnJobsNone         string
}

// OptionGroup is one optgroup in the grouped Plan/PricePlan auto-complete
// on the subscription drawer (plan §5.1). Field name `GroupLabel` matches the
// pyeza auto-complete component's expected SelectOptionGroup shape — see
// templates/components/auto-complete.html.
type OptionGroup struct {
	GroupLabel string              // group header
	Options    []map[string]string // {Value, Label} entries
}

// Data is the template data for the subscription drawer form.
type Data struct {
	FormAction  string
	IsEdit      bool
	ID          string
	Code        string
	ClientID    string
	PricePlanID string
	// Date/Time form values, split for the two-row date+time grid.
	// Stored in the operator's display TZ (DefaultTZ) for the date/time inputs;
	// JS recombines + converts to UTC RFC 3339 for the hidden field.
	DateStartDate string
	DateStartTime string
	DateEndDate   string
	DateEndTime   string
	// Pre-computed RFC 3339 hidden values; JS overwrites on every change.
	DateStartISO string
	DateEndISO   string
	// DefaultTZ is the IANA name of the operator's display timezone, surfaced as
	// data-default-tz on the form for client-side recombination.
	DefaultTZ string
	Notes     string

	Clients         []map[string]string
	PricePlans      []map[string]string
	SearchClientURL string
	SearchPlanURL   string
	ClientLabel     string
	PlanLabel       string
	ClientLocked    bool
	// ClientBillingCurrency is the selected client's billing currency, passed to
	// the plan search URL so the grouped auto-complete only shows plans in that
	// currency. Empty = no currency filter.
	ClientBillingCurrency string

	// 2026-04-27 plan-client-scope plan §5 — grouped picker options. When
	// non-empty, the template renders the grouped variant instead of the
	// flat search auto-complete.
	PlanOptionGroups []OptionGroup

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.1 — Spawn Jobs
	// section state. SpawnJobsAvailable controls section visibility (true
	// when the selected Plan resolves to one or more JobTemplates).
	// SpawnJobsDefault controls the default checked state of the toggle.
	// SpawnJobsSummary is the resolved summary string, or empty when no
	// templates resolve. SpawnJobsPartialURL is the HTMX endpoint that
	// re-renders the section on Plan select change.
	SpawnJobsAvailable  bool
	SpawnJobsDefault    bool
	SpawnJobsSummary    string
	SpawnJobsPartialURL string

	Labels       Labels
	CommonLabels any
}
