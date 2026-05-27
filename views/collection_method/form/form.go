package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// Data is the template data for the collection method drawer form (§A-1 scaffold).
// The kind-specific slot (§A-2) is rendered by a separate fragment template;
// the FragmentData below feeds it, both on initial render and on the HTMX
// category-change swap.
type Data struct {
	FormAction  string
	WorkspaceID string // injected by ViewAdapter for action_workspace_guard
	IsEdit      bool
	ID          string

	// FragmentURL is the hx-get target the category select fires on change
	// (§A-1: hx-trigger="change" → swaps #kind-specific-slot).
	FragmentURL string

	// ── Common fields (outer scaffold) ──
	Name            string
	Category        string
	PostingKind     string
	AudienceMode    string
	TaxEffectKind   string
	EligibilityRule string
	BalanceAccount  string
	TargetAccount   string
	Lifecycle       string
	Source          string
	TemplateCode    string
	Revision        string
	VersionStatus   string
	Supersedes      string

	// Option slices
	CategoryOptions      []types.SelectOption
	PostingKindOptions   []types.SelectOption
	AudienceModeOptions  []types.SelectOption
	TaxEffectOptions     []types.SelectOption
	LifecycleOptions     []types.SelectOption
	SourceOptions        []types.SelectOption
	VersionStatusOptions []types.SelectOption

	// Kind-specific fragment (§A-2). Embedded so the drawer can render the
	// initial slot inline AND the fragment endpoint can render it standalone.
	Fragment FragmentData

	Labels       centymo.CollectionMethodFormLabels
	FragLabels   centymo.CollectionMethodFragmentLabels
	CommonLabels pyeza.CommonLabels
}

// FragmentData drives the kind-specific slot (§A-2). Only the fields relevant
// to the active Category are populated.
type FragmentData struct {
	Category   string
	FragLabels centymo.CollectionMethodFragmentLabels

	// Voucher-program
	DefaultFaceValue  string
	DefaultExpiryDays string

	// Advance-program
	AdvanceKind        string
	DefaultBalanceAcct string
	DefaultTargetAcct  string
	DefaultPeriodCount string
	DefaultPeriodUnit  string

	// Bank-account
	BankName      string
	AccountFormat string
}
