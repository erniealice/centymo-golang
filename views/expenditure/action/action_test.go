package action

import (
	"testing"

	centymo "github.com/erniealice/centymo-golang"
)

// TestFormLabelsFromExpenditureAllFieldsPopulated asserts that every field on
// the result Labels struct is populated when the source ExpenditureLabels has
// every relevant field set to a non-empty value. This test prevents the
// silent-empty-field bug class: if a future field is added to the source or
// destination and the mapper forgets to copy it, this test fails.
func TestFormLabelsFromExpenditureAllFieldsPopulated(t *testing.T) {
	src := centymo.ExpenditureLabels{
		Form: centymo.ExpenditureFormLabels{
			VendorName:          "VendorName",
			VendorNamePlaceholder: "VendorNamePlaceholder",
			ExpenditureCategory: "ExpenditureCategory",
			ExpenditureDate:     "ExpenditureDate",
			TotalAmount:         "TotalAmount",
			Currency:            "Currency",
			ReferenceNumber:     "ReferenceNumber",
			Notes:               "Notes",
			NotesPlaceholder:    "NotesPlaceholder",
			Status:              "Status",
			ExpenditureType:     "ExpenditureType",
			// Info fields
			NameInfo:            "NameInfo",
			ExpenditureTypeInfo: "ExpenditureTypeInfo",
			CategoryInfo:        "CategoryInfo",
			DateInfo:            "DateInfo",
			AmountInfo:          "AmountInfo",
			CurrencyInfo:        "CurrencyInfo",
			ReferenceNumberInfo: "ReferenceNumberInfo",
			SupplierInfo:        "SupplierInfo",
			NotesInfo:           "NotesInfo",
		},
		Types: centymo.ExpenditureTypeLabels{
			Expense:  "Expense",
			Purchase: "Purchase",
		},
		Status: centymo.ExpenditureStatusLabels{
			Pending:   "Pending",
			Approved:  "Approved",
			Paid:      "Paid",
			Cancelled: "Cancelled",
		},
	}

	result := formLabels(src)

	tests := []struct {
		name  string
		value string
	}{
		// Fields sourced from l.Form.*
		{"Name", result.Name},
		{"NamePlaceholder", result.NamePlaceholder},
		{"Category", result.Category},
		// Hardcoded sentinel — always non-empty from mapper itself.
		{"Supplier", result.Supplier},
		{"Date", result.Date},
		{"Amount", result.Amount},
		{"Currency", result.Currency},
		{"ReferenceNumber", result.ReferenceNumber},
		{"Notes", result.Notes},
		{"NotesPlaceholder", result.NotesPlaceholder},
		{"Status", result.Status},
		{"ExpenditureType", result.ExpenditureType},
		// Fields sourced from l.Types.*
		{"TypeExpense", result.TypeExpense},
		{"TypePurchase", result.TypePurchase},
		// Fields sourced from l.Status.*
		{"StatusPending", result.StatusPending},
		{"StatusApproved", result.StatusApproved},
		{"StatusPaid", result.StatusPaid},
		{"StatusCancelled", result.StatusCancelled},
		// Hardcoded sentinel fields — always non-empty from mapper itself.
		{"CurrencyPlaceholder", result.CurrencyPlaceholder},
		{"PaymentTerms", result.PaymentTerms},
		{"SelectPaymentTerm", result.SelectPaymentTerm},
		{"DueDate", result.DueDate},
		{"LinkToPurchaseOrder", result.LinkToPurchaseOrder},
		// Info fields sourced from l.Form.*
		{"NameInfo", result.NameInfo},
		{"ExpenditureTypeInfo", result.ExpenditureTypeInfo},
		{"CategoryInfo", result.CategoryInfo},
		{"DateInfo", result.DateInfo},
		{"AmountInfo", result.AmountInfo},
		{"CurrencyInfo", result.CurrencyInfo},
		{"ReferenceNumberInfo", result.ReferenceNumberInfo},
		{"SupplierInfo", result.SupplierInfo},
		{"NotesInfo", result.NotesInfo},
	}

	for _, tc := range tests {
		if tc.value == "" {
			t.Errorf("field %s is empty; mapper may have omitted a field", tc.name)
		}
	}
}
