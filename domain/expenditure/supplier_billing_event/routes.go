package supplier_billing_event

// routes.go — supplier_billing_event entity route constants. Extracted from the
// domain-level routes.go during the per-entity restructure. Entity-local
// naming (SupplierBillingEvent prefix stripped). Pure structural move — route strings are
// byte-identical.

const (
	// SupplierBillingEvent (buying-side MILESTONE anchor).
	ListURL      = "/supplier-billing-events/list/{status}"
	DetailURL    = "/supplier-billing-events/detail/{id}"
	RecognizeURL = "/action/supplier-billing-event/recognize/{id}"
)
