package centymo

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// OptionValueSeparator is the canonical separator between concatenated
// product_option_value labels. Used by the variants table on the product
// detail page and by every drawer picker that surfaces a variant's
// option-value tuple inline (e.g., "Red / Large / Cotton"). Keep this
// definition as the single source of truth — when the design system
// updates the visual style, only this string changes.
const OptionValueSeparator = " / "

// CurrencyLabels carries the translated labels for each supported currency option.
// Each drawer form that renders a currency select constructs this from its lyngua
// labels and passes it to BuildCurrencyOptions.
type CurrencyLabels struct {
	PHP string
	USD string
	// Future currencies (SGD, EUR, etc.) added here propagate to every form that
	// uses BuildCurrencyOptions without any template edits.
}

// BuildCurrencyOptions returns the select options for a currency field.
// Used by every monetary drawer form. Extend the struct + append here to add
// new currencies globally.
func BuildCurrencyOptions(l CurrencyLabels) []types.SelectOption {
	return []types.SelectOption{
		{Value: "PHP", Label: l.PHP},
		{Value: "USD", Label: l.USD},
	}
}

// CurrencyLabelsFromCommon maps the lyngua-loaded CommonLabels.Currency into
// the CurrencyLabels struct required by BuildCurrencyOptions. Call this in each
// drawer form's view-data builder so the select options are always sourced from
// the central translation files.
func CurrencyLabelsFromCommon(cl pyeza.CommonLabels) CurrencyLabels {
	return CurrencyLabels{
		PHP: cl.Currency.PHP,
		USD: cl.Currency.USD,
	}
}
