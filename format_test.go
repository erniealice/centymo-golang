package centymo

import (
	"testing"
)

func TestFormatCentavoAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		centavos int64
		currency string
		want     string
	}{
		{
			name:     "standard PHP amount",
			centavos: 5000000,
			currency: "PHP",
			want:     "PHP 50,000.00",
		},
		{
			name:     "zero amount",
			centavos: 0,
			currency: "PHP",
			want:     "PHP 0.00",
		},
		{
			name:     "small amount 1 centavo",
			centavos: 1,
			currency: "PHP",
			want:     "PHP 0.01",
		},
		{
			name:     "exact peso no cents",
			centavos: 10000,
			currency: "PHP",
			want:     "PHP 100.00",
		},
		{
			name:     "negative amount",
			centavos: -5000000,
			currency: "PHP",
			want:     "PHP -50,000.00",
		},
		{
			name:     "negative small amount",
			centavos: -1,
			currency: "PHP",
			want:     "PHP -0.01",
		},
		{
			name:     "empty currency defaults to PHP",
			centavos: 10000,
			currency: "",
			want:     "PHP 100.00",
		},
		{
			name:     "USD currency",
			centavos: 123456,
			currency: "USD",
			want:     "USD 1,234.56",
		},
		{
			name:     "large amount with commas",
			centavos: 123456789012,
			currency: "PHP",
			want:     "PHP 1,234,567,890.12",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := FormatCentavoAmount(tc.centavos, tc.currency)
			if got != tc.want {
				t.Errorf("FormatCentavoAmount(%v, %q) = %q, want %q",
					tc.centavos, tc.currency, got, tc.want)
			}
		})
	}
}

func TestFormatWithCommas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{"zero", 0, "0.00"},
		{"small number", 1.23, "1.23"},
		{"hundreds", 999.99, "999.99"},
		{"thousands", 1234.56, "1,234.56"},
		{"ten thousands", 50000.00, "50,000.00"},
		{"millions", 1234567.89, "1,234,567.89"},
		{"negative thousands", -1234.56, "-1,234.56"},
		{"negative millions", -1000000.50, "-1,000,000.50"},
		{"very small decimal", 0.01, "0.01"},
		{"integer value", 5000.00, "5,000.00"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := FormatWithCommas(tc.value)
			if got != tc.want {
				t.Errorf("FormatWithCommas(%v) = %q, want %q", tc.value, got, tc.want)
			}
		})
	}
}

func TestFormatCentavoAmount_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		centavos int64
		currency string
		want     string
	}{
		{
			name:     "large value",
			centavos: 99999999999900,
			currency: "PHP",
			want:     "PHP 999,999,999,999.00",
		},
		{
			name:     "large negative value",
			centavos: -99999999999900,
			currency: "PHP",
			want:     "PHP -999,999,999,999.00",
		},
		{
			name:     "amount 1 centavo",
			centavos: 1,
			currency: "PHP",
			want:     "PHP 0.01",
		},
		{
			name:     "amount 99 centavos",
			centavos: 99,
			currency: "PHP",
			want:     "PHP 0.99",
		},
		{
			name:     "amount 100 centavos is exactly 1 unit",
			centavos: 100,
			currency: "PHP",
			want:     "PHP 1.00",
		},
		{
			name:     "very large negative amount",
			centavos: -999999999999,
			currency: "USD",
			want:     "USD -9,999,999,999.99",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := FormatCentavoAmount(tc.centavos, tc.currency)
			if got != tc.want {
				t.Errorf("FormatCentavoAmount(%v, %q) = %q, want %q",
					tc.centavos, tc.currency, got, tc.want)
			}
		})
	}
}

func TestFormatWithCommas_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{"very large positive", 9999999999999.00, "9,999,999,999,999.00"},
		{"very large negative", -9999999999999.00, "-9,999,999,999,999.00"},
		{"billion", 1000000000.00, "1,000,000,000.00"},
		{"very small negative", -0.01, "-0.01"},
		{"exactly one", 1.00, "1.00"},
		{"large number with cents", 1234567890.12, "1,234,567,890.12"},
		{"3 digits no comma", 999.00, "999.00"},
		{"4 digits with comma", 1000.00, "1,000.00"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := FormatWithCommas(tc.value)
			if got != tc.want {
				t.Errorf("FormatWithCommas(%v) = %q, want %q", tc.value, got, tc.want)
			}
		})
	}
}

func TestFormatIntegerWithCommas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input int64
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{12, "12"},
		{123, "123"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
		{1000000000, "1,000,000,000"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()
			got := formatIntegerWithCommas(tc.input)
			if got != tc.want {
				t.Errorf("formatIntegerWithCommas(%d) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestFormatIntegerWithCommas_EdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input int64
		want  string
	}{
		{"max int64", 9223372036854775807, "9,223,372,036,854,775,807"},
		{"large round number", 1000000000000, "1,000,000,000,000"},
		{"exactly 4 digits", 1000, "1,000"},
		{"exactly 7 digits", 1000000, "1,000,000"},
		{"single digit", 9, "9"},
		{"two digits", 99, "99"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := formatIntegerWithCommas(tc.input)
			if got != tc.want {
				t.Errorf("formatIntegerWithCommas(%d) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
