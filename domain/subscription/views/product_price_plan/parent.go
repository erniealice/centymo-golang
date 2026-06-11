package product_price_plan

// ParentContext captures the parent PricePlan fields the PPP table and drawer
// need to know about: currency (locks the per-line currency), billing_kind
// (decides whether billing_treatment renders), amount_basis (drives the
// banner explaining what the line prices mean), and pre-formatted display
// strings for the read-only context block above the editable fields.
//
// The parent caller (e.g. price_schedule/detail/plan) resolves this value
// via its own loadParentContext helper and passes it into BuildTable and the
// action constructors. The product_price_plan package never reaches into the
// parent package to load it.
type ParentContext struct {
	Currency    string
	BillingKind string
	AmountBasis string

	// Display strings — empty when the corresponding source data is missing.
	BillingKindDisplay    string
	AmountBasisDisplay    string
	BillingCycleDisplay   string
	TermDisplay           string
	ParentCurrencyDisplay string
	RateCardName          string
}
