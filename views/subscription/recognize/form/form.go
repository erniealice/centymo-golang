// Package form owns the template data shape for the subscription recognize-revenue
// drawer (subscription-recognize-drawer-form.html). Pure types only — no Deps,
// no context.Context, no repository imports.
package form

// Labels mirrors the lyngua-driven label set that the recognize drawer surfaces.
// Built from centymo.SubscriptionRecognizeLabels + centymo.SubscriptionInvoicesLabels
// at handler time so the template only touches a single ".Labels.X" path per string.
type Labels struct {
	Title                 string
	Subtitle              string
	ContextSection        string
	ClientLabel           string
	PlanLabel             string
	QuantityLabel         string
	PeriodSection         string
	PeriodStart           string
	PeriodEnd             string
	RevenueDate           string
	LineItemsSection      string
	ColumnDescription     string
	ColumnUnitPrice       string
	ColumnQuantity        string
	ColumnLineTotal       string
	ColumnTreatment       string
	TotalLabel            string
	RemoveLine            string
	TreatmentRecurring    string
	TreatmentFirstCycle   string
	TreatmentUsageBased   string
	TreatmentOneTime      string
	NotesLabel            string
	NotesPlaceholder      string
	Generate              string
	Cancel                string
	Timezone              string
	StartDateInfo         string
	EndDateInfo           string
	StartTimeInfo         string
	EndTimeInfo           string

	// Blocking error banners.
	CurrencyMismatchError   string
	IdempotencyError        string
	IdempotencyExistingLink string
	NoLinesError            string

	// 2026-04-29 milestone-billing plan §5 / Phase D — milestone-specific
	// drawer labels. Surfaced only when pricePlan.billing_kind = MILESTONE.
	MilestoneSelect            string
	MilestoneSelectPlaceholder string
	NoReadyMilestone           string
	MilestoneNotApplicable     string
	BillAmount                 string
	LeaveRemainderOpen         string
	CloseShort                 string
	PartialReason              string
	PartialReasonRequired      string
	OverBillingRejected        string
}

// MilestoneOption is a single row in the milestone select.
// Selectable = true for READY/DEFERRED; false (disabled) for BILLED so
// already-invoiced events are still visible to operators. UNSPECIFIED
// (pending) events are filtered out by the handler — they don't reach
// the template.
type MilestoneOption struct {
	EventID         string
	SequenceLabel   string // e.g. "M1 Kickoff & Design" — falls back to event ID
	Status          string // proto enum String(), lowercased — "ready" / "billed" / etc.
	StatusLabel     string // localized status label
	BillableAmount  int64
	BillableDisplay string // formatted "₱150,000.00"
	Currency        string
	Selectable      bool // false → render disabled
	Hidden          bool // true → don't render at all (UNSPECIFIED)
	Selected        bool
}

// PreviewLine is the row shape consumed by the drawer template.
// Mirrors revenuepb.PreviewLineItem but exposes only fields the template
// actually renders.
type PreviewLine struct {
	ProductPricePlanID string
	Description        string
	UnitPrice          int64
	Quantity           float64
	TotalPrice         int64
	Currency           string
	Treatment          string
	TreatmentLabel     string
}

// Data is the template data for the recognize-revenue drawer.
type Data struct {
	FormAction       string
	SubscriptionID   string
	SubscriptionName string
	ClientLabel      string
	PlanLabel        string
	Quantity         int32
	Currency         string

	// Period (date + time grid, IANA tz aware — same pattern as the standard
	// subscription drawer).
	PeriodStartDate string
	PeriodStartTime string
	PeriodStartISO  string
	PeriodEndDate   string
	PeriodEndTime   string
	PeriodEndISO    string
	DefaultTZ       string

	// Revenue date (single date input).
	RevenueDate string

	// Notes — auto-prefixed with the period marker.
	Notes string

	// Line items preview.
	PreviewLines []PreviewLine
	TotalAmount  int64

	// Blocking-error state.
	CurrencyMismatch      bool
	ClientCurrency        string
	PlanCurrency          string
	IdempotencyConflict   bool
	ConflictingRevenueID  string
	ConflictingRevenueURL string
	NoLinesToInvoice      bool

	// Non-blocking warnings (e.g. usage-based skipped notice).
	Warnings []string

	// 2026-04-29 milestone-billing plan §5 / Phase D — milestone fields.
	// IsMilestone gates the drawer's milestone-only rows (select +
	// partial-billing controls + period-input suppression). When true the
	// MilestoneOptions slice drives the select; SelectedBillingEventID is the
	// pre-selected (or operator-chosen) BillingEvent.id. PartialDefault is
	// the centavo amount the bill-amount input pre-fills with — equal to the
	// selected event's billable_amount unless the operator typed something
	// else. Reason is required when PartialDefault != ChosenAmount.
	IsMilestone            bool
	MilestoneOptions       []MilestoneOption
	SelectedBillingEventID string
	BillAmountDisplay      string // editable input pre-fill
	LeaveRemainderOpen     bool
	CloseShort             bool
	PartialReasonValue     string
	OverBillingError       bool

	Labels       Labels
	CommonLabels any
}
