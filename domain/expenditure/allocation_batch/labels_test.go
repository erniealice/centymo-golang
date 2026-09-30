package allocation_batch

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
	for _, rel := range []string{"general/allocation_batch.json", "leasing/allocation_batch.json"} {
		var doc struct {
			L Labels `json:"allocation_batch"`
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
		L Labels `json:"allocation_batch"`
	}
	if err := json.Unmarshal(lyngua(t, "general/allocation_batch.json"), &doc); err != nil {
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
	for _, code := range []string{ErrValidation, ErrNotFound, ErrNotDraft, ErrAlreadyPublished, ErrSourceClaimedByRecognition, ErrSourceClaimedByAllocation, ErrSharesTotalMismatch, ErrDenominatorInvalid, ErrNoRecoverableShare, ErrServicePeriodInvalid, ErrPolicyComponentInvalid, ErrTransactionRequired, ErrLockUnavailable, ErrAgreementTermMissing, ErrTermBoundaryCrossed, ErrOverlap} {
		got := l.ErrorMessage(ErrorKind(coded{code}))
		if got == "" || got == l.Errors.Generic {
			t.Errorf("code %q maps to %q", code, got)
		}
	}
	if got := l.ErrorMessage(ErrorKind(errors.New("boom"))); got != l.Errors.Generic {
		t.Errorf("uncoded error maps to %q, want generic", got)
	}
}

func TestRouteMapKeys(t *testing.T) {
	m := DefaultRoutes().RouteMap()
	for _, k := range []string{"allocation_batch.allocate", "allocation_batch.preview", "allocation_batch.publish", "allocation_batch.view"} {
		if m[k] == "" {
			t.Errorf("route key %s missing", k)
		}
	}
}
