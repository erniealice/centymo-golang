package action

import (
	"testing"

	centymo "github.com/erniealice/centymo-golang"
)

// TestFormLabelsFromPurchaseOrderAllFieldsPopulated asserts that every field on
// the result Labels struct is populated when the source ExpenditureLabels has
// every relevant field set to a non-empty value. This test prevents the
// silent-empty-field bug class: if a future field is added to the source or
// destination and the mapper forgets to copy it, this test fails.
func TestFormLabelsFromPurchaseOrderAllFieldsPopulated(t *testing.T) {
	src := centymo.ExpenditureLabels{
		Form: centymo.ExpenditureFormLabels{
			ExpenditureDate:  "ExpenditureDate",
			Currency:         "Currency",
			PaymentTerms:     "PaymentTerms",
			Notes:            "Notes",
			NotesPlaceholder: "NotesPlaceholder",
			Status:           "Status",
		},
		PurchaseOrder: centymo.PurchaseOrderLabels{
			Form: centymo.PurchaseOrderFormLabels{
				PONumberInfo:     "PONumberInfo",
				SupplierInfo:     "SupplierInfo",
				POTypeInfo:       "POTypeInfo",
				OrderDateInfo:    "OrderDateInfo",
				CurrencyInfo:     "CurrencyInfo",
				PaymentTermsInfo: "PaymentTermsInfo",
				NotesInfo:        "NotesInfo",
			},
		},
	}

	result := formLabels(src)

	tests := []struct {
		name  string
		value string
	}{
		// Hardcoded sentinel fields — always non-empty from the mapper itself.
		{"PoNumber", result.PoNumber},
		{"SupplierID", result.SupplierID},
		{"PoType", result.PoType},
		// Fields sourced from l.Form.*
		{"OrderDate", result.OrderDate},
		{"Currency", result.Currency},
		{"PaymentTerms", result.PaymentTerms},
		{"Notes", result.Notes},
		{"NotesPlaceholder", result.NotesPlaceholder},
		{"Status", result.Status},
		// Info fields sourced from l.PurchaseOrder.Form.*
		{"PoNumberInfo", result.PoNumberInfo},
		{"SupplierIDInfo", result.SupplierIDInfo},
		{"PoTypeInfo", result.PoTypeInfo},
		{"OrderDateInfo", result.OrderDateInfo},
		{"CurrencyInfo", result.CurrencyInfo},
		{"PaymentTermsInfo", result.PaymentTermsInfo},
		{"NotesInfo", result.NotesInfo},
	}

	for _, tc := range tests {
		if tc.value == "" {
			t.Errorf("field %s is empty; mapper may have omitted a field", tc.name)
		}
	}
}
