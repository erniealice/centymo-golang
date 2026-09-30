package action

import (
	"testing"

	collection "github.com/erniealice/centymo-golang/domain/treasury/collection"
)

type codedErr struct{ code string }

func (e codedErr) Error() string     { return "raw " + e.code }
func (e codedErr) ErrorCode() string { return e.code }

// A coded receipt refusal shows the Lyngua message; an uncoded error keeps its text.
func TestRefusalMessageMapsReceiptHasApplicationsToItsLabel(t *testing.T) {
	deps := &Deps{Labels: collection.DefaultLabels()}
	if got := refusalMessage(deps, codedErr{collection.ErrKindReceiptHasApplications}); got != deps.Labels.Errors.ReceiptHasApplications {
		t.Errorf("coded refusal = %q", got)
	}
	if got := refusalMessage(deps, codedErr{"other"}); got != "raw other" {
		t.Errorf("unmapped refusal = %q", got)
	}
}
