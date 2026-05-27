package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// Data is the template data for the disbursement method drawer form (§A-1
// scaffold, buying-side). Mirror of the collection-side form minus audience_mode
// (D-4.9: no buying-side audience model). The kind-specific slot (§A-2 buying
// side: bank-account / check / advance) is rendered by a separate fragment.
type Data struct {
	FormAction  string
	WorkspaceID string
	IsEdit      bool
	ID          string

	FragmentURL string

	// ── Common fields ──
	Name           string
	Category       string
	PostingKind    string
	TaxEffectKind  string
	BalanceAccount string
	TargetAccount  string
	Lifecycle      string
	Source         string
	TemplateCode   string
	Revision       string
	VersionStatus  string
	Supersedes     string

	CategoryOptions      []types.SelectOption
	PostingKindOptions   []types.SelectOption
	TaxEffectOptions     []types.SelectOption
	LifecycleOptions     []types.SelectOption
	SourceOptions        []types.SelectOption
	VersionStatusOptions []types.SelectOption

	Fragment FragmentData

	Labels       centymo.DisbursementMethodFormLabels
	FragLabels   centymo.DisbursementMethodFragmentLabels
	CommonLabels pyeza.CommonLabels
}

// FragmentData drives the buying-side kind-specific slot (§A-2 disbursement).
type FragmentData struct {
	Category   string
	FragLabels centymo.DisbursementMethodFragmentLabels

	// Bank-account
	BankName      string
	AccountFormat string

	// Check
	CheckSeries   string
	SigningPolicy string

	// Advance-program
	AdvanceKind        string
	DefaultBalanceAcct string
	DefaultTargetAcct  string
	DefaultPeriodCount string
	DefaultPeriodUnit  string
}
