package product_price_plan

import "testing"

// ParseOptionalCentavos shares the one strict centavos parser (C22): blank stays unset, and
// Inf/NaN/negative/fractional-cent inputs are refused instead of being stored.
func TestParseOptionalCentavosIsStrict(t *testing.T) {
	if v, ok := ParseOptionalCentavos("  "); !ok || v != nil {
		t.Errorf("blank = %v,%v", v, ok)
	}
	if v, ok := ParseOptionalCentavos("30000.50"); !ok || v == nil || *v != 3000050 {
		t.Errorf("valid = %v,%v", v, ok)
	}
	for _, bad := range []string{"Inf", "-Inf", "NaN", "+Inf", "-1", "-0.01", "1.005", "abc", "1e3"} {
		if v, ok := ParseOptionalCentavos(bad); ok || v != nil {
			t.Errorf("%q accepted: %v", bad, v)
		}
	}
}
