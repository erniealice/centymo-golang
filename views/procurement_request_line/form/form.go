package form

import (
	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	procurementrequestlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request_line"
)

// Data is the template data for the procurement request line drawer form.
type Data struct {
	FormAction           string
	IsEdit               bool
	ID                   string
	ProcurementRequestID string

	// Core fields
	Description           string
	LineType              string
	ProductID             string
	Quantity              string
	EstimatedUnitPrice    string
	EstimatedTotalPrice   string
	ExpenditureCategoryID string
	LocationID            string
	LineNumber            string

	// SPS Wave 3 — F1 fulfillment_mode picker + RECURRING fields
	// FulfillmentMode is the short canonical token: "outright" | "stockable" | "recurring" | "petty"
	FulfillmentMode        string
	FulfillmentModeOptions []types.RadioOption

	// RECURRING-only fields (rendered conditionally when FulfillmentMode == "recurring")
	RecurringCycleValue  string
	RecurringCycleUnit   string
	RecurringTermValue   string
	RecurringTermUnit    string
	RecurringUnitOptions []types.SelectOption

	// Options
	Products []types.SelectOption

	Labels       centymo.ProcurementRequestLabels
	CommonLabels pyeza.CommonLabels
}

// ParseFulfillmentMode maps the form's short token into the proto enum value.
// Falls back to UNSPECIFIED on unknown tokens.
func ParseFulfillmentMode(token string) procurementrequestlinepb.ProcurementRequestLineFulfillmentMode {
	switch token {
	case "outright":
		return procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_OUTRIGHT
	case "stockable":
		return procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_STOCKABLE
	case "recurring":
		return procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_RECURRING
	case "petty":
		return procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_PETTY
	}
	return procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_UNSPECIFIED
}

// FulfillmentModeToToken is the inverse of ParseFulfillmentMode for edit-form
// pre-population — returns the short canonical token used in the radio group.
func FulfillmentModeToToken(mode procurementrequestlinepb.ProcurementRequestLineFulfillmentMode) string {
	switch mode {
	case procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_OUTRIGHT:
		return "outright"
	case procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_STOCKABLE:
		return "stockable"
	case procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_RECURRING:
		return "recurring"
	case procurementrequestlinepb.ProcurementRequestLineFulfillmentMode_PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_PETTY:
		return "petty"
	}
	return ""
}
