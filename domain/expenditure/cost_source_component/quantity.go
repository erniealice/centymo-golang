package cost_source_component

import (
	"strconv"
	"strings"

	"github.com/erniealice/pyeza-golang/types"
)

// MaxScale is the largest number of decimals accepted for a basis quantity.
const MaxScale = 6

// FormatQuantity renders a scaled integer quantity ("12500", scale 2 -> "125.00").
func FormatQuantity(scaled int64, scale int32) string {
	if scale <= 0 {
		return strconv.FormatInt(scaled, 10)
	}
	neg := scaled < 0
	if neg {
		scaled = -scaled
	}
	s := strconv.FormatInt(scaled, 10)
	for int32(len(s)) <= scale {
		s = "0" + s
	}
	cut := len(s) - int(scale)
	out := s[:cut] + "." + s[cut:]
	if neg {
		out = "-" + out
	}
	return out
}

// ParseQuantity parses a decimal string ("125.5") into (scaled, scale), keeping the decimals the
// operator typed (at most MaxScale). Blank returns present=false, valid=true; an invalid or
// negative value returns valid=false (the view answers with the form-invalid label).
func ParseQuantity(raw string) (scaled int64, scale int32, present, valid bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, 0, false, true
	}
	whole, frac, _ := strings.Cut(raw, ".")
	if whole == "" {
		whole = "0"
	}
	if len(frac) > MaxScale {
		return 0, 0, false, false
	}
	digits := whole + frac
	n, err := strconv.ParseInt(digits, 10, 64)
	if err != nil || n < 0 || strings.ContainsAny(digits, "+-") {
		return 0, 0, false, false
	}
	return n, int32(len(frac)), true, true
}

// ParseAmountCentavos parses a display amount ("1250.75") into integer centavos through the
// shared pyeza parser (no float, no local copy: R5 M6) and refuses non-positive or malformed
// input (ok=false). Thousands separators are stripped first. Quantities keep ParseQuantity: a
// quantity is a scaled decimal, not money.
func ParseAmountCentavos(raw string) (int64, bool) {
	n, err := types.ParseCentavos(strings.ReplaceAll(raw, ",", ""))
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}

// FormatAmount renders integer centavos as a two-decimal display amount.
func FormatAmount(centavos int64) string {
	return FormatQuantity(centavos, 2)
}
