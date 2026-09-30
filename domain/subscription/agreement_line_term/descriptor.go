package agreement_line_term

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the data-only compose descriptor of the agreement charge terms:
// a Labels pointer + lyngua binding and no Mount, exactly like product_price_plan.
// The subscription unit resolves the POST-OVERLAY labels through
// compose.LabelsOf[*Labels](mc, "subscription.agreement_line_term"); the unit is
// opt-in (block.WithKnownCostRecovery), so the subscription "Charge terms" tab exists
// only where the app opted in.
func Describe() compose.Unit {
	l := DefaultLabels()
	return compose.Unit{
		Key:       "subscription.agreement_line_term",
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "agreement_line_term.json", Key: "agreement_line_term"},
		LabelName: "AgreementLineTermLabels",
	}
}
