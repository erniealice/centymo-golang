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

// DefaultCurrencyOptions returns the in-code fallback list, sorted by ISO
// code. Tiers should override the lyngua `currency.options` JSON array to
// add or reorder entries. Label format: "ISO — Full Name (Symbol)".
// Description carries the country / region the currency is tied to so the
// form-group select can render a per-option hint.
//
// The shape (types.SelectOption) is owned by pyeza; the curated list is
// owned here in centymo. Lyngua serves the JSON; nobody owns the words
// "PHP", "USD", etc. except this function.
func DefaultCurrencyOptions() []types.SelectOption {
	return []types.SelectOption{
		{Value: "AUD", Label: "AUD — Australian Dollar (A$)", Description: "Australia"},
		{Value: "CAD", Label: "CAD — Canadian Dollar (C$)", Description: "Canada"},
		{Value: "CHF", Label: "CHF — Swiss Franc (CHF)", Description: "Switzerland, Liechtenstein"},
		{Value: "CNY", Label: "CNY — Chinese Yuan (¥)", Description: "Mainland China"},
		{Value: "EUR", Label: "EUR — Euro (€)", Description: "Eurozone"},
		{Value: "GBP", Label: "GBP — British Pound (£)", Description: "United Kingdom"},
		{Value: "HKD", Label: "HKD — Hong Kong Dollar (HK$)", Description: "Hong Kong SAR"},
		{Value: "IDR", Label: "IDR — Indonesian Rupiah (Rp)", Description: "Indonesia"},
		{Value: "INR", Label: "INR — Indian Rupee (₹)", Description: "India"},
		{Value: "JPY", Label: "JPY — Japanese Yen (¥)", Description: "Japan"},
		{Value: "KRW", Label: "KRW — South Korean Won (₩)", Description: "South Korea"},
		{Value: "MYR", Label: "MYR — Malaysian Ringgit (RM)", Description: "Malaysia"},
		{Value: "NZD", Label: "NZD — New Zealand Dollar (NZ$)", Description: "New Zealand"},
		{Value: "PHP", Label: "PHP — Philippine Peso (₱)", Description: "Philippines"},
		{Value: "SGD", Label: "SGD — Singapore Dollar (S$)", Description: "Singapore"},
		{Value: "THB", Label: "THB — Thai Baht (฿)", Description: "Thailand"},
		{Value: "TWD", Label: "TWD — New Taiwan Dollar (NT$)", Description: "Taiwan"},
		{Value: "USD", Label: "USD — US Dollar ($)", Description: "United States"},
		{Value: "VND", Label: "VND — Vietnamese Dong (₫)", Description: "Vietnam"},
	}
}

// BuildCurrencyOptions returns the select options for a currency field.
// Sources the list from the supplied lyngua-loaded labels; falls back to
// DefaultCurrencyOptions() when the bundle didn't ship the list. Used by
// every monetary drawer form.
func BuildCurrencyOptions(cl pyeza.CurrencyLabels) []types.SelectOption {
	if len(cl.Options) == 0 {
		return DefaultCurrencyOptions()
	}
	return cl.Options
}
