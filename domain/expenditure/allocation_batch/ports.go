package allocation_batch

import (
	"context"
	"errors"
	"fmt"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	allocationbatchpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_batch"
	allocationsharepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/allocation_share"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	agreementlinetermpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/agreement_line_term"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// Error codes carried by a use-case refusal (`ErrorCode()`); each maps to the
// Lyngua key `allocation_batch.errors.<code>` (Labels.ErrorMessage). The codes
// are the espyna allocation_batch use-case codes (progress-f3-s1-allocation.md).
const (
	ErrNone                       = ""
	ErrValidation                 = "validation"
	ErrNotFound                   = "not_found"
	ErrNotDraft                   = "not_draft"
	ErrAlreadyPublished           = "already_published"
	ErrSourceClaimedByRecognition = "source_claimed_by_recognition"
	ErrSourceClaimedByAllocation  = "source_claimed_by_allocation"
	ErrSharesTotalMismatch        = "shares_total_mismatch"
	ErrDenominatorInvalid         = "denominator_invalid"
	ErrNoRecoverableShare         = "no_recoverable_share"
	ErrServicePeriodInvalid       = "service_period_invalid"
	ErrPolicyComponentInvalid     = "policy_component_invalid"
	ErrTransactionRequired        = "transaction_required"
	ErrLockUnavailable            = "lock_unavailable"
	ErrAgreementTermMissing       = "agreement_term_missing"
	ErrTermBoundaryCrossed        = "term_boundary_crossed"
	ErrOverlap                    = "overlap"
	ErrPermissionDenied           = "permission_denied" // coded strict-gate denial (espyna actiongate)
	ErrUnknown                    = "unknown"
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

// UseCases is the set of espyna closures the allocation drawers consume, exactly
// the proto request/response shape of `uc.Expenditure.AllocationBatch.<Field>.Execute`
// (+ the existing reads that build the participants list). The block binds them; a
// nil closure disables the feature that needs it and the view fails closed.
type UseCases struct {
	CreateAllocationBatch          func(context.Context, *allocationbatchpb.CreateAllocationBatchRequest) (*allocationbatchpb.CreateAllocationBatchResponse, error)
	UpdateAllocationBatchShares    func(context.Context, *allocationbatchpb.UpdateAllocationBatchSharesRequest) (*allocationbatchpb.UpdateAllocationBatchSharesResponse, error)
	PublishAllocationBatch         func(context.Context, *allocationbatchpb.PublishAllocationBatchRequest) (*allocationbatchpb.PublishAllocationBatchResponse, error)
	GetAllocationBatchListPageData func(context.Context, *allocationbatchpb.GetAllocationBatchListPageDataRequest) (*allocationbatchpb.GetAllocationBatchListPageDataResponse, error)
	GetAllocationShareListPageData func(context.Context, *allocationsharepb.GetAllocationShareListPageDataRequest) (*allocationsharepb.GetAllocationShareListPageDataResponse, error)

	// Participant reads: the component being allocated, the agreement terms that
	// decide who can share it, and the subscription (+ client) display names.
	ReadCostSourceComponent     func(context.Context, *costsourcecomponentpb.ReadCostSourceComponentRequest) (*costsourcecomponentpb.ReadCostSourceComponentResponse, error)
	ListAgreementLineTerms      func(context.Context, *agreementlinetermpb.ListAgreementLineTermsRequest) (*agreementlinetermpb.ListAgreementLineTermsResponse, error)
	GetSubscriptionItemPageData func(context.Context, *subscriptionpb.GetSubscriptionItemPageDataRequest) (*subscriptionpb.GetSubscriptionItemPageDataResponse, error)
}

// maxPages bounds every page walker in this package (R5 m7): past it the walk fails closed.
const maxPages = 1000

var errPageBound = fmt.Errorf("allocation_batch: page bound (%d pages) exceeded", maxPages)

// pageRequest is the paged request every loop in this package sends.
func pageRequest(page int32) *commonpb.PaginationRequest {
	return &commonpb.PaginationRequest{Limit: 100, Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}}}
}

func stringEquals(field, value string) *commonpb.FilterRequest {
	return &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
		Field: field,
		FilterType: &commonpb.TypedFilter_StringFilter{StringFilter: &commonpb.StringFilter{
			Value: value, Operator: commonpb.StringOperator_STRING_EQUALS, CaseSensitive: true,
		}},
	}}}
}
