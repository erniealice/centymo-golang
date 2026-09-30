package recovery_document

import (
	"context"
	"log"
	"strings"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	collectionapplicationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_application"
)

// UseCases is the set of espyna closures the recovery document views consume
// (`uc.Revenue.RecoveryDocument.<X>.Execute`, `uc.Treasury.CollectionApplication.<X>.Execute`),
// proto request/response shape. A nil closure disables the feature that needs
// it and the view fails closed.
type UseCases struct {
	ListRecoveryDocuments func(context.Context, *recoverydocumentpb.ListRecoveryDocumentsRequest) (*recoverydocumentpb.ListRecoveryDocumentsResponse, error)
	ReadRecoveryDocument  func(context.Context, *recoverydocumentpb.ReadRecoveryDocumentRequest) (*recoverydocumentpb.ReadRecoveryDocumentResponse, error)
	VoidRecoveryDocument  func(context.Context, *recoverydocumentpb.VoidRecoveryDocumentRequest) (*recoverydocumentpb.VoidRecoveryDocumentResponse, error)

	// ListCollectionApplications feeds the "payments applied" tab.
	ListCollectionApplications func(context.Context, *collectionapplicationpb.ListCollectionApplicationsRequest) (*collectionapplicationpb.ListCollectionApplicationsResponse, error)

	// ListClients resolves display names (optional; a nil closure renders the id).
	ListClients func(context.Context, *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	// ReadClient resolves the one client name of the detail page (optional).
	ReadClient func(context.Context, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
}

// ClientName returns the display name of one client through ReadClient (the
// detail page needs exactly one name, never a full list). "" when the reader is
// unwired or the read fails (logged); callers fall back to the id.
func (u *UseCases) ClientName(ctx context.Context, id string) string {
	if u == nil || u.ReadClient == nil || id == "" {
		return ""
	}
	resp, err := u.ReadClient(ctx, &clientpb.ReadClientRequest{Data: &clientpb.Client{Id: id}})
	if err != nil {
		log.Printf("recovery_document: read client %s: %v", id, err)
		return ""
	}
	if len(resp.GetData()) == 0 {
		return ""
	}
	return displayName(resp.GetData()[0])
}

func displayName(c *clientpb.Client) string {
	name := c.GetName()
	if name == "" {
		if usr := c.GetUser(); usr != nil {
			name = strings.TrimSpace(usr.GetFirstName() + " " + usr.GetLastName())
		}
	}
	return name
}

// ClientNames returns id -> display name for the workspace's clients, paged
// through (C11: no silent 100-row cap). A page error is logged and the names
// read so far are returned; callers fall back to the id.
func (u *UseCases) ClientNames(ctx context.Context) map[string]string {
	out := map[string]string{}
	if u == nil || u.ListClients == nil {
		return out
	}
	for page := int32(1); page <= 1000; page++ {
		resp, err := u.ListClients(ctx, &clientpb.ListClientsRequest{Pagination: pageRequest(page)})
		if err != nil {
			log.Printf("recovery_document: list clients page %d: %v", page, err)
			return out
		}
		for _, c := range resp.GetData() {
			out[c.GetId()] = displayName(c)
		}
		if int32(len(resp.GetData())) < pageLimit {
			return out
		}
	}
	log.Printf("recovery_document: client names page bound exceeded")
	return out
}

const pageLimit = int32(100)

func pageRequest(page int32) *commonpb.PaginationRequest {
	return &commonpb.PaginationRequest{
		Limit:  pageLimit,
		Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}},
	}
}

type pageBoundError struct{}

func (pageBoundError) Error() string { return "recovery_document: page bound exceeded" }

// ApplicationsOf returns every collection application targeting the document
// (paged through, C11; the page bound fails closed).
func (u *UseCases) ApplicationsOf(ctx context.Context, documentID string) ([]*collectionapplicationpb.CollectionApplication, error) {
	filters := &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
		Field: "recovery_document_id",
		FilterType: &commonpb.TypedFilter_StringFilter{
			StringFilter: &commonpb.StringFilter{Value: documentID, Operator: commonpb.StringOperator_STRING_EQUALS},
		},
	}}}
	var out []*collectionapplicationpb.CollectionApplication
	for page := int32(1); page <= 1000; page++ {
		resp, err := u.ListCollectionApplications(ctx, &collectionapplicationpb.ListCollectionApplicationsRequest{
			Filters: filters, Pagination: pageRequest(page),
		})
		if err != nil {
			return nil, err
		}
		items := resp.GetData()
		out = append(out, items...)
		if int32(len(items)) < pageLimit {
			return out, nil
		}
	}
	return nil, pageBoundError{}
}
