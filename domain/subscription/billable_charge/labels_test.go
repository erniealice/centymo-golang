package billable_charge

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func lynguaPath(rel string) string {
	return filepath.Join("..", "..", "..", "..", "lyngua", "translations", "en", rel)
}

// The Go json tags must byte-match the Lyngua key tree: decoding the general and
// leasing billable_charge.json files with unknown fields disallowed proves that
// every Lyngua key is consumed by a Go field.
func TestLabelsMatchLynguaKeyTree(t *testing.T) {
	for _, rel := range []string{"general/billable_charge.json", "leasing/billable_charge.json"} {
		raw, err := os.ReadFile(lynguaPath(rel))
		if err != nil {
			t.Skipf("lyngua file not reachable (%s): %v", rel, err)
		}
		var doc struct {
			BillableCharge Labels `json:"billable_charge"`
		}
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&doc); err != nil {
			t.Errorf("%s does not decode into Labels: %v", rel, err)
		}
	}
}

// C14: the English defaults equal the general Lyngua value of EVERY string field.
func TestDefaultLabelsEqualGeneralLyngua(t *testing.T) {
	raw, err := os.ReadFile(lynguaPath("general/billable_charge.json"))
	if err != nil {
		t.Skipf("lyngua file not reachable: %v", err)
	}
	var doc struct {
		BillableCharge Labels `json:"billable_charge"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
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
	walk("Labels", reflect.ValueOf(DefaultLabels()), reflect.ValueOf(doc.BillableCharge))
}

// Every use-case error code maps to its own message.
func TestEveryErrorCodeHasALabel(t *testing.T) {
	l := DefaultLabels()
	seen := map[string]string{}
	for _, c := range []string{ErrKindValidation, ErrKindNotFound, ErrKindNotIssued, ErrKindAdjustNotDownward,
		ErrKindTransactionRequired, ErrKindLockUnavailable, ErrKindObligationConflict} {
		m := l.ErrorMessage(c)
		if m == "" {
			t.Errorf("code %q has no message", c)
		}
		if prev, dup := seen[m]; dup {
			t.Errorf("codes %q and %q share the message %q", prev, c, m)
		}
		seen[m] = c
	}
	if l.ErrorMessage("nonsense") != "" {
		t.Error("an unmapped code must return the empty string (caller falls back to the general error)")
	}
}

type codedErr string

func (e codedErr) Error() string     { return "x: " + string(e) }
func (e codedErr) ErrorCode() string { return string(e) }

func TestErrorKindUsesErrorCode(t *testing.T) {
	if got := ErrorKind(fmt.Errorf("wrapped: %w", codedErr("not_issued"))); got != ErrKindNotIssued {
		t.Errorf("wrapped coded error = %q", got)
	}
	if got := ErrorKind(errors.New("boom")); got != ErrKindUnknown {
		t.Errorf("uncoded error = %q", got)
	}
	if ErrorKind(nil) != ErrKindNone {
		t.Error("nil error must classify as none")
	}
}

func TestFormatSubstitutesPositionalArgs(t *testing.T) {
	if got := Format("{0} charges for {1} customers", "3", "2"); got != "3 charges for 2 customers" {
		t.Errorf("Format = %q", got)
	}
}
