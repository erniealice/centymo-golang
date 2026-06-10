package centymo

// OptionValueSeparator is the canonical separator between concatenated
// product_option_value labels. Used by the variants table on the product
// detail page and by every drawer picker that surfaces a variant's
// option-value tuple inline (e.g., "Red / Large / Cotton"). Keep this
// definition as the single source of truth — when the design system
// updates the visual style, only this string changes.
const OptionValueSeparator = " / "
