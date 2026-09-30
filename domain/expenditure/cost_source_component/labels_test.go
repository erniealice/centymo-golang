package cost_source_component

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func lyngua(t *testing.T, rel string) []byte {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "lyngua", "translations", "en", rel))
	if err != nil {
		t.Fatalf("lyngua file not reachable (%s): %v", rel, err)
	}
	return raw
}

// The Go json tags must byte-match the Lyngua key tree: decoding the general and
// leasing files with unknown fields disallowed proves every key is consumed.
func TestLabelsMatchLynguaKeyTree(t *testing.T) {
	for _, rel := range []string{"general/cost_source_component.json", "leasing/cost_source_component.json"} {
		var doc struct {
			L Labels `json:"cost_source_component"`
		}
		dec := json.NewDecoder(bytes.NewReader(lyngua(t, rel)))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&doc); err != nil {
			t.Errorf("%s does not decode into Labels: %v", rel, err)
		}
	}
}

// C14: the English defaults equal the general Lyngua value of EVERY string field.
func TestDefaultLabelsEqualGeneralLyngua(t *testing.T) {
	var doc struct {
		L Labels `json:"cost_source_component"`
	}
	if err := json.Unmarshal(lyngua(t, "general/cost_source_component.json"), &doc); err != nil {
		t.Fatal(err)
	}
	var walk func(path string, def, got reflect.Value)
	walk = func(path string, def, got reflect.Value) {
		switch def.Kind() {
		case reflect.Struct:
			for i := 0; i < def.NumField(); i++ {
				walk(path+"."+def.Type().Field(i).Name, def.Field(i), got.Field(i))
			}
		case reflect.String:
			if def.String() == "" {
				t.Errorf("%s: empty Go default", path)
			}
			if def.String() != got.String() {
				t.Errorf("%s: default %q != general lyngua %q", path, def.String(), got.String())
			}
		}
	}
	walk("Labels", reflect.ValueOf(DefaultLabels()), reflect.ValueOf(doc.L))
}

type coded struct{ code string }

func (c coded) Error() string     { return "refused: " + c.code }
func (c coded) ErrorCode() string { return c.code }

// Every use-case refusal code maps to its own message (never the generic one).
func TestErrorMessageMapsEveryCode(t *testing.T) {
	l := DefaultLabels()
	for _, code := range []string{ErrValidation, ErrNotFound, ErrClaimed, ErrTransactionRequired, ErrReferenceInvalid, ErrLockUnavailable} {
		got := l.ErrorMessage(ErrorKind(coded{code}))
		if got == "" || got == l.Errors.Generic {
			t.Errorf("code %q maps to %q", code, got)
		}
	}
	if got := l.ErrorMessage(ErrorKind(errors.New("boom"))); got != l.Errors.Generic {
		t.Errorf("uncoded error maps to %q, want generic", got)
	}
}

func TestQuantityAndAmountParsing(t *testing.T) {
	for _, c := range []struct {
		raw     string
		scaled  int64
		scale   int32
		present bool
	}{{"", 0, 0, false}, {"12", 12, 0, true}, {"12.50", 1250, 2, true}, {".5", 5, 1, true}} {
		s, sc, present, valid := ParseQuantity(c.raw)
		if !valid || s != c.scaled || sc != c.scale || present != c.present {
			t.Errorf("ParseQuantity(%q) = %d,%d,%v,%v", c.raw, s, sc, present, valid)
		}
	}
	for _, bad := range []string{"-1", "1.1234567", "x"} {
		if _, _, _, valid := ParseQuantity(bad); valid {
			t.Errorf("ParseQuantity(%q) accepted", bad)
		}
	}
	if got := FormatQuantity(1250, 2); got != "12.50" {
		t.Errorf("FormatQuantity = %q", got)
	}
	if got := FormatQuantity(5, 2); got != "0.05" {
		t.Errorf("FormatQuantity small = %q", got)
	}
	if n, ok := ParseAmountCentavos("1,250.75"); !ok || n != 125075 {
		t.Errorf("ParseAmountCentavos = %d, %v", n, ok)
	}
	for _, bad := range []string{"0", "-5", "1.234", "abc", ""} {
		if _, ok := ParseAmountCentavos(bad); ok {
			t.Errorf("ParseAmountCentavos(%q) accepted", bad)
		}
	}
}

func TestRouteMapKeys(t *testing.T) {
	m := DefaultRoutes().RouteMap()
	for _, k := range []string{"cost_source_component.add", "cost_source_component.edit", "cost_source_component.delete", "cost_source_component.table"} {
		if m[k] == "" {
			t.Errorf("route key %s missing", k)
		}
	}
}

// R5 M6: money goes through pyeza types.ParseCentavos; the malformed forms the old local parser
// accepted ("1.-5" parsed as 95 centavos) are refused.
func TestParseAmountCentavosRefusesSignedFraction(t *testing.T) {
	for _, bad := range []string{"1.-5", "1.+5", "+1.50", "1..5", "1.5.0", " ", "-0.01"} {
		if n, ok := ParseAmountCentavos(bad); ok {
			t.Errorf("ParseAmountCentavos(%q) = %d accepted", bad, n)
		}
	}
	if n, ok := ParseAmountCentavos(" 10 "); !ok || n != 1000 {
		t.Errorf("ParseAmountCentavos(\" 10 \") = %d, %v", n, ok)
	}
}
