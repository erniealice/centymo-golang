package cost_source_component

import (
	"context"
	"errors"

	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	expenditurelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_line_item"
)

// Error codes carried by a use-case refusal (`ErrorCode()`); each maps to the
// Lyngua key `cost_source_component.errors.<code>` (Labels.ErrorMessage). The codes
// are the espyna cost_source_component use-case codes (progress-f3-s1-allocation.md).
const (
	ErrNone                = ""
	ErrValidation          = "validation"
	ErrNotFound            = "not_found"
	ErrClaimed             = "claimed"
	ErrTransactionRequired = "transaction_required"
	ErrReferenceInvalid    = "reference_invalid"
	ErrLockUnavailable     = "lock_unavailable"
	ErrPermissionDenied    = "permission_denied" // coded strict-gate denial (espyna actiongate)
	ErrUnknown             = "unknown"
)

// ErrorKind classifies a use-case error through its `ErrorCode()` (no import of
// espyna). Unknown / uncoded errors classify as ErrUnknown.
func ErrorKind(err error) string {
	if err == nil {
		return ErrNone
	}
	var coded interface{ ErrorCode() string }
	if errors.As(err, &coded) {
		if c := coded.ErrorCode(); c != "" {
			return c
		}
	}
	return ErrUnknown
}

// UseCases is the set of espyna closures the recoverable-cost views consume,
// exactly the proto request/response shape of
// `uc.Expenditure.CostSourceComponent.<Field>.Execute`. The block binds them; a
// nil closure disables the feature that needs it and the view fails closed.
type UseCases struct {
	CreateCostSourceComponent func(context.Context, *costsourcecomponentpb.CreateCostSourceComponentRequest) (*costsourcecomponentpb.CreateCostSourceComponentResponse, error)
	ReadCostSourceComponent   func(context.Context, *costsourcecomponentpb.ReadCostSourceComponentRequest) (*costsourcecomponentpb.ReadCostSourceComponentResponse, error)
	UpdateCostSourceComponent func(context.Context, *costsourcecomponentpb.UpdateCostSourceComponentRequest) (*costsourcecomponentpb.UpdateCostSourceComponentResponse, error)
	DeleteCostSourceComponent func(context.Context, *costsourcecomponentpb.DeleteCostSourceComponentRequest) (*costsourcecomponentpb.DeleteCostSourceComponentResponse, error)
	// ListCostSourceComponents is paged by the caller (ListForExpenditure): the response has no
	// pagination block and the espyna aggregate has no GetListPageData for this entity, so the
	// view pages until a short page (C11: never a silent cap).
	ListCostSourceComponents func(context.Context, *costsourcecomponentpb.ListCostSourceComponentsRequest) (*costsourcecomponentpb.ListCostSourceComponentsResponse, error)

	// ListExpenditureLineItems feeds the optional "bill line" picker (nil hides it).
	ListExpenditureLineItems func(context.Context, *expenditurelineitempb.ListExpenditureLineItemsRequest) (*expenditurelineitempb.ListExpenditureLineItemsResponse, error)
}
