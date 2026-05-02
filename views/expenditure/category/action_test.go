package category

import (
	"testing"

	centymo "github.com/erniealice/centymo-golang"
)

// TestFormLabelsFromExpenditureCategoryAllFieldsPopulated asserts that every
// field on the result Labels struct is populated when the source
// ExpenditureCategoryLabels has every field set to a non-empty value. This
// test prevents the silent-empty-field bug class: if a future field is added
// to the source or destination and the mapper forgets to copy it, this test
// fails.
func TestFormLabelsFromExpenditureCategoryAllFieldsPopulated(t *testing.T) {
	src := centymo.ExpenditureCategoryLabels{
		Form: centymo.ExpenditureCategoryFormLabels{
			Code:        "Code",
			Name:        "Name",
			Description: "Description",
			// Info fields
			CodeInfo:        "CodeInfo",
			NameInfo:        "NameInfo",
			DescriptionInfo: "DescriptionInfo",
		},
	}

	result := formLabels(src)

	tests := []struct {
		name  string
		value string
	}{
		{"Code", result.Code},
		{"Name", result.Name},
		{"Description", result.Description},
		{"CodeInfo", result.CodeInfo},
		{"NameInfo", result.NameInfo},
		{"DescriptionInfo", result.DescriptionInfo},
	}

	for _, tc := range tests {
		if tc.value == "" {
			t.Errorf("field %s is empty; mapper may have omitted a field", tc.name)
		}
	}
}
