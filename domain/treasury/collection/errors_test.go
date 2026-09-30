package collection

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// The receipt refusal label is a Lyngua key: the general collection.json errors object decodes
// into ErrorLabels (no unknown key) and its receipt_has_applications value equals the Go default.
func TestReceiptRefusalLabelMatchesLyngua(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "lyngua", "translations", "en", "general", "collection.json"))
	if err != nil {
		t.Skipf("lyngua file not reachable: %v", err)
	}
	var doc struct {
		Collection struct {
			Errors json.RawMessage `json:"errors"`
		} `json:"collection"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	var got ErrorLabels
	dec := json.NewDecoder(bytes.NewReader(doc.Collection.Errors))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&got); err != nil {
		t.Fatalf("collection.errors does not decode into ErrorLabels: %v", err)
	}
	if want := DefaultLabels().Errors.ReceiptHasApplications; got.ReceiptHasApplications != want {
		t.Errorf("receipt_has_applications = %q, default %q", got.ReceiptHasApplications, want)
	}
}

func TestErrorMessageMapsOnlyTheReceiptRefusal(t *testing.T) {
	l := DefaultLabels()
	if l.ErrorMessage(ErrKindReceiptHasApplications) != l.Errors.ReceiptHasApplications || l.Errors.ReceiptHasApplications == "" {
		t.Error("receipt_has_applications not mapped")
	}
	if l.ErrorMessage("other") != "" || l.ErrorMessage(ErrKindNone) != "" {
		t.Error("unmapped code must return the empty string")
	}
}
