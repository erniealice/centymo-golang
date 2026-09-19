package subscription

// CreateOptions configures subscription creation independently of vertical vocabulary.
// The zero value preserves the ordinary subscription date/range behavior.
type CreateOptions struct {
	// CommencementPricing selects pricing valid at the start date and enables
	// plan-term suggestions, client currency hydration and matching create checks.
	// It is server-owned configuration, never a request-controlled permission.
	CommencementPricing bool
}
