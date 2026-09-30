package billable_charge

import (
	"context"
	"log"
	"strings"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	documentseriespb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"
	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	billablechargepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billable_charge"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// UseCases is the set of espyna closures the billable charge views consume,
// exactly the proto request/response shape of the espyna use cases
// (`uc.Subscription.BillableCharge.<X>.Execute`, `uc.Revenue.RecoveryDocument.<X>.Execute`,
// `uc.Revenue.DocumentSeries.<X>.Execute`). The block binds them; a nil closure
// disables the feature that needs it and the view fails closed. Permission
// gates live in the views AND in the use cases.
type UseCases struct {
	ListBillableCharges    func(context.Context, *billablechargepb.ListBillableChargesRequest) (*billablechargepb.ListBillableChargesResponse, error)
	AdjustBillableCharge   func(context.Context, *billablechargepb.AdjustBillableChargeRequest) (*billablechargepb.AdjustBillableChargeResponse, error)
	IssueRecoveryDocuments func(context.Context, *recoverydocumentpb.IssueRecoveryDocumentsRequest) (*recoverydocumentpb.IssueRecoveryDocumentsResponse, error)
	ListDocumentSeries     func(context.Context, *documentseriespb.ListDocumentSeriesRequest) (*documentseriespb.ListDocumentSeriesResponse, error)

	// Display-name resolvers (optional; a nil closure renders the id).
	ListClients       func(context.Context, *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ListSubscriptions func(context.Context, *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)
}

const pageLimit = int32(100)

func pageRequest(page int32) *commonpb.PaginationRequest {
	return &commonpb.PaginationRequest{
		Limit:  pageLimit,
		Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}},
	}
}

// pageAll walks every page of a paginated list (no silent 100-row cap, C11).
// The page bound fails closed: reaching it returns an error rather than a
// truncated list.
func pageAll[T any](fetch func(*commonpb.PaginationRequest) ([]T, *commonpb.PaginationResponse, error)) ([]T, error) {
	var out []T
	for page := int32(1); page <= 1000; page++ {
		items, pag, err := fetch(pageRequest(page))
		if err != nil {
			return nil, err
		}
		out = append(out, items...)
		if pag != nil {
			if pag.CurrentPage != nil && pag.TotalPages != nil {
				if pag.GetCurrentPage() >= pag.GetTotalPages() {
					return out, nil
				}
				continue
			}
			if int32(len(out)) >= pag.GetTotalItems() {
				return out, nil
			}
		} else if int32(len(items)) < pageLimit {
			return out, nil
		}
		if len(items) == 0 {
			return out, nil
		}
	}
	return nil, errPageBound
}

type pageBoundError struct{}

func (pageBoundError) Error() string { return "billable_charge: page bound exceeded" }

var errPageBound error = pageBoundError{}

// ListAllByStatus returns every charge in the given status (paged through).
func (u *UseCases) ListAllByStatus(ctx context.Context, status billablechargepb.BillableChargeStatus) ([]*billablechargepb.BillableCharge, error) {
	return pageAll(func(p *commonpb.PaginationRequest) ([]*billablechargepb.BillableCharge, *commonpb.PaginationResponse, error) {
		resp, err := u.ListBillableCharges(ctx, &billablechargepb.ListBillableChargesRequest{Status: &status, Pagination: p})
		if err != nil {
			return nil, nil, err
		}
		return resp.GetData(), resp.GetPagination(), nil
	})
}

// ActiveRecoverySeries returns the ACTIVE document series that issue recovery
// documents (the issue drawer's series picker).
func (u *UseCases) ActiveRecoverySeries(ctx context.Context) ([]*documentseriespb.DocumentSeries, error) {
	all, err := pageAll(func(p *commonpb.PaginationRequest) ([]*documentseriespb.DocumentSeries, *commonpb.PaginationResponse, error) {
		resp, err := u.ListDocumentSeries(ctx, &documentseriespb.ListDocumentSeriesRequest{Pagination: p})
		if err != nil {
			return nil, nil, err
		}
		return resp.GetData(), nil, nil
	})
	if err != nil {
		return nil, err
	}
	out := make([]*documentseriespb.DocumentSeries, 0, len(all))
	for _, s := range all {
		if s.GetStatus() == documentseriespb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_ACTIVE &&
			s.GetDocumentKind() == enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_RECOVERY_DOCUMENT {
			out = append(out, s)
		}
	}
	return out, nil
}

// ClientNames returns id -> display name for the workspace's clients, paged
// through (C11: no silent 100-row cap). A page error is logged and an empty map
// returned; callers fall back to the id.
func (u *UseCases) ClientNames(ctx context.Context) map[string]string {
	out := map[string]string{}
	if u == nil || u.ListClients == nil {
		return out
	}
	items, err := pageAll(func(p *commonpb.PaginationRequest) ([]*clientpb.Client, *commonpb.PaginationResponse, error) {
		resp, err := u.ListClients(ctx, &clientpb.ListClientsRequest{Pagination: p})
		if err != nil {
			return nil, nil, err
		}
		return resp.GetData(), nil, nil
	})
	if err != nil {
		log.Printf("billable_charge: list clients: %v", err)
	}
	for _, c := range items {
		name := c.GetName()
		if name == "" {
			if usr := c.GetUser(); usr != nil {
				name = strings.TrimSpace(usr.GetFirstName() + " " + usr.GetLastName())
			}
		}
		out[c.GetId()] = name
	}
	return out
}

// SubscriptionNames returns id -> name for the workspace's subscriptions, paged
// through; a page error is logged and an empty map returned.
func (u *UseCases) SubscriptionNames(ctx context.Context) map[string]string {
	out := map[string]string{}
	if u == nil || u.ListSubscriptions == nil {
		return out
	}
	items, err := pageAll(func(p *commonpb.PaginationRequest) ([]*subscriptionpb.Subscription, *commonpb.PaginationResponse, error) {
		resp, err := u.ListSubscriptions(ctx, &subscriptionpb.ListSubscriptionsRequest{Pagination: p})
		if err != nil {
			return nil, nil, err
		}
		return resp.GetData(), nil, nil
	})
	if err != nil {
		log.Printf("billable_charge: list subscriptions: %v", err)
	}
	for _, s := range items {
		out[s.GetId()] = s.GetName()
	}
	return out
}

// ReadByID returns the charge with the id (nil when absent), through the list
// use case filtered on id (there is no separate read RPC).
func (u *UseCases) ReadByID(ctx context.Context, id string) (*billablechargepb.BillableCharge, error) {
	resp, err := u.ListBillableCharges(ctx, &billablechargepb.ListBillableChargesRequest{
		Filters: &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
			Field: "id",
			FilterType: &commonpb.TypedFilter_StringFilter{
				StringFilter: &commonpb.StringFilter{Value: id, Operator: commonpb.StringOperator_STRING_EQUALS},
			},
		}}},
		Pagination: pageRequest(1),
	})
	if err != nil {
		return nil, err
	}
	for _, c := range resp.GetData() {
		if c.GetId() == id {
			return c, nil
		}
	}
	return nil, nil
}
