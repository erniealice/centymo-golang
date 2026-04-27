// Package form owns the shared template data shape for the canonical
// price-plan-drawer-form.html template. All three callers — plan-detail
// rate-cards tab, price-schedule-detail package-prices tab, and the
// standalone price-plan list — build a form.Data and render the same
// template, differing only in the Context discriminator.
package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// Context identifies which URL context the drawer was opened from. The
// template uses this to decide which field is a locked display row vs. an
// auto-complete. See price-plan-drawer-form.html for the branching.
type Context string

const (
	ContextPlan       Context = "Plan"
	ContextSchedule   Context = "Schedule"
	ContextStandalone Context = "Standalone"
)

// Data is the template shape for price-plan-drawer-form.html.
//
// Both PlanID and ScheduleID are always set (either pre-seeded from path
// context or picked by the user on submit). The template emits a hidden
// input for whichever one matches Context and an auto-complete for the
// other; in Standalone mode both are auto-completes.
type Data struct {
	FormAction string
	IsEdit     bool
	Context    Context
	ID         string // PricePlan.ID — present on edit, empty on add

	// Context-locked IDs.
	PlanID     string
	ScheduleID string

	// Locked display values shown in the disabled form-group in their
	// respective contexts. The rate-card auto-complete surfaces each
	// schedule's location as a per-option Description — no separate
	// LocationName field is needed here.
	PlanName     string
	ScheduleName string

	// Form values (edit preload / add defaults).
	Name          string
	Description   string
	Amount        string // decimal display, e.g. "1500.00"
	Currency      string
	DurationValue string // DEPRECATED: Phase 1 dual-write; keep for read-back
	DurationUnit  string // DEPRECATED: Phase 1 dual-write; keep for read-back
	Active        bool

	// Wave 2: new billing semantics fields (Phase 1 dual-write alongside DurationValue/Unit).
	BillingKind        string
	BillingKindOptions []types.SelectOption
	AmountBasis        string
	AmountBasisOptions []types.SelectOption
	BillingCycleValue   string // int32 as string for form field
	BillingCycleUnit    string
	TermValue           string // int32 as string for form field; backs default_term_value on the wire
	TermUnit            string // backs default_term_unit on the wire
	DurationUnitOptions []types.SelectOption // reused for both billing_cycle_unit and term unit

	// Auto-complete option lists. Each entry is {Value, Label, Selected?}.
	// PlanOptions is consumed in Schedule + Standalone contexts;
	// ScheduleOptions in Plan + Standalone contexts.
	PlanOptions     []map[string]any
	ScheduleOptions []map[string]any

	SelectedPlanID        string
	SelectedPlanLabel     string
	SelectedScheduleID    string
	SelectedScheduleLabel string

	// Pricing-lock signal from the reference checker. When InUse is true
	// on an edit load, the template disables amount/currency/duration_*.
	InUse       bool
	LockMessage string

	// 2026-04-27 plan-client-scope plan §6.7 — surfaced when the parent
	// PriceSchedule is client-scoped. The template renders an info banner
	// "{ClientName} owns this rate card..." above the first form section
	// and emits a hidden client_id input populated server-side. Empty when
	// the parent schedule is master.
	ParentScheduleClientID   string
	ParentScheduleClientName string

	Labels       Labels
	CommonLabels pyeza.CommonLabels
}

// Labels are the template-facing flat labels consumed by the drawer.
type Labels struct {
	SectionBasic           string
	SectionPricing         string
	NameLabel              string
	NamePlaceholder        string
	DescriptionLabel       string
	DescriptionPlaceholder string
	AmountLabel            string
	AmountPlaceholder      string
	CurrencyLabel          string
	CurrencyPlaceholder    string
	CurrencyPHP            string
	CurrencyUSD            string
	DurationLabel          string
	DurationUnitLabel      string
	ActiveLabel            string
	PlanLabel              string
	PlanPlaceholder        string
	PlanSearch             string
	ScheduleLabel          string
	SchedulePlaceholder    string
	ScheduleSearch         string
	LocationHintPrefix     string

	// Wave 2: new billing semantics labels.
	BillingKindLabel        string
	AmountBasisLabel        string
	BillingCycleLabel       string
	BillingCyclePlaceholder string
	TermLabel               string
	TermPlaceholder         string
	TermOpenEndedHelp       string

	// Field-level info text surfaced via an info button beside each label.
	// Hover/click opens a popover explaining what the field means.
	PlanInfo         string
	ScheduleInfo     string
	NameInfo         string
	DescriptionInfo  string
	BillingKindInfo  string
	AmountBasisInfo  string
	AmountInfo       string
	CurrencyInfo     string
	BillingCycleInfo string
	TermInfo         string
	ActiveInfo       string

	// 2026-04-27 plan-client-scope plan §6.7 — info banner template.
	// Templated via {{.ClientName}}. Blank means "no banner".
	ParentScheduleClientNotice string
}

// LabelsFromPriceSchedule maps the price-schedule-side PlanForm labels into
// the flat template-facing Labels shape. Used by callers in a Schedule context
// (add/edit from a PriceSchedule detail page).
func LabelsFromPriceSchedule(pf centymo.PriceSchedulePlanFormLabels) Labels {
	return Labels{
		SectionBasic:           pf.SectionPackage,
		SectionPricing:         pf.SectionPricing,
		NameLabel:              pf.NameLabel,
		NamePlaceholder:        pf.NamePlaceholder,
		DescriptionLabel:       pf.DescriptionLabel,
		DescriptionPlaceholder: pf.DescriptionPlaceholder,
		AmountLabel:            pf.AmountLabel,
		AmountPlaceholder:      pf.AmountPlaceholder,
		CurrencyLabel:          pf.CurrencyLabel,
		CurrencyPlaceholder:    pf.CurrencyPlaceholder,
		CurrencyPHP:            pf.CurrencyPHP,
		CurrencyUSD:            pf.CurrencyUSD,
		DurationLabel:          pf.DurationLabel,
		DurationUnitLabel:      pf.UnitLabel,
		ActiveLabel:            pf.ActiveLabel,
		PlanLabel:              pf.PackageLabel,
		PlanPlaceholder:        pf.PackagePlaceholder,
		PlanSearch:             pf.PackageSearch,
		ScheduleLabel:          pf.PriceScheduleField,
		SchedulePlaceholder:    pf.SchedulePlaceholder,
		ScheduleSearch:         pf.ScheduleSearch,
		LocationHintPrefix:     pf.LocationHintPrefix,
		// Wave 2: billing labels not yet on PriceSchedulePlanFormLabels —
		// leave empty here; PricePlanFormLabels path provides them when
		// the standalone action builds the form.
	}
}

// LabelsFromPricePlan maps centymo.PricePlanFormLabels (the tier-aware
// struct populated by lyngua) into the flat template-facing Labels shape.
// Fields are sourced entirely from the struct; no hardcoded English fallbacks.
// Lyngua is the single source of truth — if a key is missing, the field renders
// empty and the missing key should be fixed in the JSON, not here.
func LabelsFromPricePlan(pp centymo.PricePlanFormLabels) Labels {
	return Labels{
		SectionBasic:           pp.SectionBasic,
		SectionPricing:         pp.SectionPricing,
		NameLabel:              pp.Name,
		NamePlaceholder:        pp.NamePlaceholder,
		DescriptionLabel:       pp.Description,
		DescriptionPlaceholder: pp.DescPlaceholder,
		AmountLabel:            pp.Amount,
		AmountPlaceholder:      pp.AmountPlaceholder,
		CurrencyLabel:          pp.Currency,
		CurrencyPlaceholder:    pp.CurrencyPlaceholder,
		CurrencyPHP:            pp.CurrencyPHP,
		CurrencyUSD:            pp.CurrencyUSD,
		DurationLabel:          pp.DurationValue,
		DurationUnitLabel:      pp.DurationUnit,
		ActiveLabel:            pp.Active,
		PlanLabel:              pp.PlanLabel,
		PlanPlaceholder:        pp.PlanPlaceholder,
		PlanSearch:             pp.PlanSearch,
		ScheduleLabel:          pp.Schedule,
		SchedulePlaceholder:    pp.SchedulePlaceholder,
		ScheduleSearch:         pp.ScheduleSearch,
		LocationHintPrefix:     pp.LocationHintPrefix,
		// Wave 2 new fields
		BillingKindLabel:        pp.BillingKindLabel,
		AmountBasisLabel:        pp.AmountBasisLabel,
		BillingCycleLabel:       pp.BillingCycleLabel,
		BillingCyclePlaceholder: pp.BillingCyclePlaceholder,
		TermLabel:               pp.TermLabel,
		TermPlaceholder:         pp.TermPlaceholder,
		TermOpenEndedHelp:       pp.TermOpenEndedHelp,
		// Field-level info popovers
		PlanInfo:         pp.PlanInfo,
		ScheduleInfo:     pp.ScheduleInfo,
		NameInfo:         pp.NameInfo,
		DescriptionInfo:  pp.DescriptionInfo,
		BillingKindInfo:  pp.BillingKindInfo,
		AmountBasisInfo:  pp.AmountBasisInfo,
		AmountInfo:       pp.AmountInfo,
		CurrencyInfo:     pp.CurrencyInfo,
		BillingCycleInfo: pp.BillingCycleInfo,
		TermInfo:         pp.TermInfo,
		ActiveInfo:       pp.ActiveInfo,
		// 2026-04-27 plan-client-scope plan §6.7.
		ParentScheduleClientNotice: pp.ParentScheduleClientNotice,
	}
}

// BuildOptions converts a slice of (id, name, description) tuples into the
// auto-complete option map shape the template expects. Description is
// surfaced by the auto-complete as a .form-hint right below the field,
// updating as the user switches selections.
func BuildOptions(entries []Option, selectedID string) []map[string]any {
	opts := make([]map[string]any, 0, len(entries))
	for _, e := range entries {
		opts = append(opts, map[string]any{
			"Value":       e.ID,
			"Label":       e.Name,
			"Description": e.Description,
			"Selected":    e.ID == selectedID,
		})
	}
	return opts
}

// Option is a simple tuple used when building auto-complete lists.
// Description is optional and rendered as a form-hint below the field
// whenever the option is the current selection.
type Option struct {
	ID          string
	Name        string
	Description string
}

// FindLabel returns the Name for a given ID, or empty string if not found.
func FindLabel(entries []Option, id string) string {
	for _, e := range entries {
		if e.ID == id {
			return e.Name
		}
	}
	return ""
}
