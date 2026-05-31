// Package form owns the template data shape for the primary revenue drawer
// (revenue-drawer-form.html). Pure types only — no Deps, no context.Context,
// no repository imports.
package form

// PaymentTermOption is a minimal struct for rendering payment term options in the form.
type PaymentTermOption struct {
	Id      string
	Name    string
	NetDays int32
}

// TaxLineRow is a single tax row displayed in the Taxes section of the revenue drawer.
type TaxLineRow struct {
	ID              string
	Direction       string // "SURCHARGE" or "WITHHOLDING"
	DirectionLabel  string // lyngua-translated display value for direction
	KindLabel       string // lyngua-translated display label for the specific tax kind (M1)
	TaxKindSnapshot string
	RegulatoryCode  string
	RateBasisPoints int64
	TaxableBase     int64
	TaxAmount       int64
	// Display strings (formatted centavos ÷ 100)
	TaxableBaseDisplay string
	TaxAmountDisplay   string
	RateDisplay        string // e.g. "12.00%"
}

// Inner holds nested form labels accessed via .Labels.Form.* in templates.
type Inner struct {
	SectionInfo               string
	CurrencyPlaceholder       string
	StatusDraft               string
	StatusComplete            string
	StatusCancelled           string
	CustomerNamePlaceholder   string
	CustomerSearchPlaceholder string
	CustomerNoResults         string
	LocationPlaceholder       string
	LocationSearchPlaceholder string
	LocationNoResults         string
}

// Labels holds i18n labels for the drawer form template.
type Labels struct {
	Customer                  string
	Date                      string
	Currency                  string
	Reference                 string
	ReferencePlaceholder      string
	Status                    string
	Notes                     string
	NotesPlaceholder          string
	Location                  string
	PaymentTerms              string
	SelectPaymentTerm         string
	DueDate                   string
	Subscription              string
	SubscriptionNoResults     string
	RevenueType               string
	RevenueTypeOneTime        string
	RevenueTypeFromSubscription string
	RevenueTypeFromActivities string
	ActivityIDs               string
	ActivityIDsPlaceholder    string
	Form                      Inner

	// Field-level info text surfaced via an info button beside each label.
	ReferenceInfo    string
	DateInfo         string
	CustomerInfo     string
	LocationInfo     string
	SubscriptionInfo string
	CurrencyInfo     string
	NotesInfo        string

	// Tax section labels (Phase 5)
	SectionTax              string
	TaxDirectionSurcharge   string
	TaxDirectionWithholding string
	TaxKind                 string
	TaxRegCode              string
	TaxRate                 string
	TaxableBase             string
	TaxAmount               string
	NetReceivable           string
	WHTAmount               string
	GrandTotal              string
	SettlementStatus        string
	Recompute               string
	AddWHTCertificate       string
	// FX dual-amount display labels
	Billed             string
	Recorded           string
	Rate               string
	RateSourceOperator string

	// TaxKindLabels maps TaxKindSnapshot enum values (e.g. "VAT_STANDARD") to
	// their translated display names. Populated at handler time from lyngua keys.
	// Used by the template via TaxLineRow.KindLabel (Phase 5 M1).
	TaxKindLabels map[string]string
}

// Data is the template data for the revenue drawer form.
type Data struct {
	FormAction            string
	WorkspaceID            string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Nonce                 string // CSP nonce; populated by ViewAdapter.injectPageData (NonceFromContext) for inline <script nonce>
	IsEdit                bool
	ID                    string
	Name                  string
	ClientID              string
	ClientLabel           string
	SearchClientURL       string
	SubscriptionID        string
	SubscriptionLabel     string
	SearchSubscriptionURL string
	ReferenceNumber       string
	Date                  string
	Currency              string
	Status                string
	Notes                 string
	LocationID            string
	LocationLabel         string
	SearchLocationURL     string
	PaymentTerms          []*PaymentTermOption
	SelectedPaymentTermID string
	DueDateString         string
	RevenueType           string
	ActivityIDs           string

	// Tax fields (Phase 5)
	TaxLines           []TaxLineRow
	GrandTotalAmount   int64
	GrandTotalDisplay  string
	CashAmountExpected int64
	CashAmountDisplay  string
	WhtAmountExpected  int64
	WhtAmountDisplay   string
	SettlementStatus   string
	CanRecompute       bool // admin permission gate
	RecomputeURL       string
	AddWHTCertURL      string

	// FX dual-amount fields (Phase 5)
	// When BillingCurrency is set, show dual-amount display.
	BillingCurrency      string
	BillingAmount        int64
	BillingAmountDisplay string
	ForexRateMicroUnits  int64
	ForexRateDisplay     string // e.g. "56.5000 PHP per USD"
	ForexRateSource      string // e.g. "operator entered, 2026-05-01"

	Labels       Labels
	CommonLabels any
}
