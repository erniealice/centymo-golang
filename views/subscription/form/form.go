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

	// 2026-05-03 — Row-level help text below the date+time rows.
	StartDateRowHelp string
	EndDateRowHelp   string

	// 2026-04-27 plan-client-scope plan §5.1 / §7 — group headers in the
	// grouped Plan / PricePlan auto-complete picker.
	PlanGroupForClient string // "For {ClientName}" — pre-resolved with ClientName injected.
	PlanGroupGeneral   string

	// 2026-05-03 — info banner below the locked Customer field; explains the
	// Plan picker is filtered to client-scoped plans only.
	PlanClientScopeNotice string

	// 2026-05-03 — Edit-locked notice rendered above the form when the
	// subscription has revenue / job references and cannot be edited.
	EditLockedReason string

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.1 / §9.
	SpawnJobsSectionTitle string
	SpawnJobsToggle       string
	SpawnJobsHelpText     string
	SpawnJobsSummary      string // {{.JobCount}} / {{.TemplateNames}} / {{.PhaseCount}} / {{.TaskCount}}
	SpawnJobsNone         string
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

	// 2026-05-03 — Reference-checker lock signal. When InUse is true (the
	// subscription is referenced by Revenue rows, subscription_attribute, or
	// operation Job rows), the drawer renders all fields read-only and hides
	// the Update button. Reassigning the plan after revenue has been
	// recognised would break the audit trail. LockMessage is the lyngua-loaded
	// notice rendered in place of the form footer.
	InUse       bool
	LockMessage string

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
