package subscription_group_product_plan

// Options — app-configurable presentation for the subscription_group_product_plan
// surfaces, set by the consuming app through the block's EngineBlock option
// (mirrors subscription_group.Options: a generic, fail-safe reference grammar
// — the zero value is today's default behavior, so apps that set no options
// are unaffected). No vertical noun; every field is a display/UX knob only,
// never a business rule.
type Options struct {
	// Picker configures the S2 add-offerings picker's default candidate state.
	Picker PickerOptions
}

// PickerOptions — the S2 picker's default candidate selection state.
type PickerOptions struct {
	// SelectAllByDefault pre-checks every candidate row (the sketch's default:
	// "[x] Select all (2)"). nil (unset) resolves to true; set to force
	// explicit per-row opt-in instead.
	SelectAllByDefault *bool
}

// SelectAll resolves the effective default (fail-safe: nil -> true).
func (p PickerOptions) SelectAll() bool {
	if p.SelectAllByDefault == nil {
		return true
	}
	return *p.SelectAllByDefault
}
