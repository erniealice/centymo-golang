package agreement_line_term

import (
	"bytes"
	"encoding/json"
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
	for _, rel := range []string{"general/agreement_line_term.json", "leasing/agreement_line_term.json"} {
		var doc struct {
			L Labels `json:"agreement_line_term"`
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
		L Labels `json:"agreement_line_term"`
	}
	if err := json.Unmarshal(lyngua(t, "general/agreement_line_term.json"), &doc); err != nil {
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
