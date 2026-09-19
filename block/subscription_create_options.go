package block

import subscription "github.com/erniealice/centymo-golang/domain/subscription/subscription"

// WithSubscriptionCreateOptions configures the generic subscription create mode.
func WithSubscriptionCreateOptions(options subscription.CreateOptions) EngineOption {
	return func(c *engineConfig) { c.subscriptionCreateOptions = options }
}

// WithSubscriptionCreatePolicy configures creation for the legacy Block API.
// Callers previously relying on implicit vertical activation must opt in here.
// This is a policy option, not a module-selection option.
func WithSubscriptionCreatePolicy(options subscription.CreateOptions) BlockOption {
	return func(c *blockConfig) { c.subscriptionCreateOptions = options }
}
