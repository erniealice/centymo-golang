package form

import (
	"time"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	scpspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// StatusOption is a select option for the schedule status enum.
// Description is required by the shared form-group select template (reads
// .Description on every option to set data-description) — keep the field
// even when unused so template execution doesn't panic.
type StatusOption struct {
	Value       string
	Label       string
	Selected    bool
	Description string
}

// Data is the template data for the SCPS drawer form.
type Data struct {
	FormAction string
	IsEdit     bool
	ID         string

	// Identity
	Name        string
	Description string
	InternalID  string

	// Scoping
	SupplierContractID    string
	SupplierContractLabel string
	SupplierContracts     []types.SelectOption

	// Validity (date-only — half-open [start, end), end blank = open-ended)
	DateStart string
	DateEnd   string
	OpenEnded bool

	// Money / location
	Currency   string
	LocationID string
	Locations  []types.SelectOption

	// Lifecycle
	Status         string
	StatusOptions  []StatusOption
	SequenceNumber string

	// Notes
	Notes string

	Labels       centymo.SupplierContractPriceScheduleFormLabels
	StatusLabels centymo.SupplierContractPriceScheduleStatusLabels
	CommonLabels pyeza.CommonLabels
}

// ParseDateUTC parses YYYY-MM-DD as midnight UTC.
func ParseDateUTC(date string, _ bool) *timestamppb.Timestamp {
	if date == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil
	}
	return timestamppb.New(t.UTC())
}

// ParseEndDate handles the open-ended checkbox: when checked, end is nil.
// Otherwise parses as YYYY-MM-DD at end-of-day UTC (23:59:59).
func ParseEndDate(date string, openEnded bool) *timestamppb.Timestamp {
	if openEnded || date == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil
	}
	// End is exclusive in the half-open window; render as end-of-day UTC.
	endOfDay := t.UTC().Add(24*time.Hour - time.Second)
	return timestamppb.New(endOfDay)
}

// FormatDateUTC formats a timestamppb.Timestamp as YYYY-MM-DD in UTC.
func FormatDateUTC(ts *timestamppb.Timestamp) string {
	if ts == nil || !ts.IsValid() {
		return ""
	}
	return ts.AsTime().UTC().Format("2006-01-02")
}

// ParseStatus maps a proto enum string name to the SupplierContractPriceScheduleStatus enum.
// Falls back to SCHEDULED on unknown values.
func ParseStatus(s string) scpspb.SupplierContractPriceScheduleStatus {
	if v, ok := scpspb.SupplierContractPriceScheduleStatus_value[s]; ok {
		return scpspb.SupplierContractPriceScheduleStatus(v)
	}
	return scpspb.SupplierContractPriceScheduleStatus_SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SCHEDULED
}
