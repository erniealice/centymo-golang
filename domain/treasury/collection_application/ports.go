package collection_application

import (
	"context"
	"log"
	"strings"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	collectionapplicationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_application"
	collectionmethodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method"
	"github.com/erniealice/pyeza-golang/types"
)

// UseCases is the set of espyna closures the receive-and-apply views consume
// (`uc.Revenue.ReceiveApply.<X>.Execute`), proto request/response shape. A nil
// closure disables the feature that needs it and the view fails closed.
type UseCases struct {
	ReceiveAndApplyCollection    func(context.Context, *collectionapplicationpb.ReceiveAndApplyCollectionRequest) (*collectionapplicationpb.ReceiveAndApplyCollectionResponse, error)
	PreviewCollectionApplication func(context.Context, *collectionapplicationpb.PreviewCollectionApplicationRequest) (*collectionapplicationpb.PreviewCollectionApplicationResponse, error)
	ReverseCollectionApplication func(context.Context, *collectionapplicationpb.ReverseCollectionApplicationRequest) (*collectionapplicationpb.ReverseCollectionApplicationResponse, error)

	// Option feeders for the drawer (optional; a nil closure renders no options).
	ListClients           func(context.Context, *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ReadClient            func(context.Context, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	ListCollectionMethods func(context.Context, *collectionmethodpb.ListCollectionMethodsRequest) (*collectionmethodpb.ListCollectionMethodsResponse, error)
}

const pageLimit = int32(100)

func pageRequest(page int32) *commonpb.PaginationRequest {
	return &commonpb.PaginationRequest{
		Limit:  pageLimit,
		Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}},
	}
}

func clientName(c *clientpb.Client) string {
	name := c.GetName()
	if name == "" {
		if usr := c.GetUser(); usr != nil {
			name = strings.TrimSpace(usr.GetFirstName() + " " + usr.GetLastName())
		}
	}
	if name == "" {
		name = c.GetId()
	}
	return name
}

// SelectedClient returns the display name of one client (the ?client_id=
// preselect of the async client search; ReadClient, never a full list). "" when
// the id is empty, the reader is unwired or the read fails (logged).
func (u *UseCases) SelectedClient(ctx context.Context, id string) string {
	if u == nil || u.ReadClient == nil || id == "" {
		return ""
	}
	resp, err := u.ReadClient(ctx, &clientpb.ReadClientRequest{Data: &clientpb.Client{Id: id}})
	if err != nil {
		log.Printf("collection_application: read client %s: %v", id, err)
		return ""
	}
	if len(resp.GetData()) == 0 {
		return ""
	}
	return clientName(resp.GetData()[0])
}

// ClientOptions returns every client of the workspace as select options, paged
// through (C11: no silent 100-row cap). It is the fallback picker source only;
// the drawer prefers the async client search (SearchClientURL). A page error is
// logged and the options read so far are returned.
func (u *UseCases) ClientOptions(ctx context.Context, selected string) []types.SelectOption {
	if u == nil || u.ListClients == nil {
		return nil
	}
	var out []types.SelectOption
	for page := int32(1); page <= 1000; page++ {
		resp, err := u.ListClients(ctx, &clientpb.ListClientsRequest{Pagination: pageRequest(page)})
		if err != nil {
			log.Printf("collection_application: list clients page %d: %v", page, err)
			return out
		}
		for _, c := range resp.GetData() {
			out = append(out, types.SelectOption{Value: c.GetId(), Label: clientName(c), Selected: c.GetId() == selected})
		}
		if int32(len(resp.GetData())) < pageLimit {
			return out
		}
	}
	log.Printf("collection_application: client options page bound exceeded")
	return out
}

// MethodOptions returns the workspace's collection methods as select options.
func (u *UseCases) MethodOptions(ctx context.Context) []types.SelectOption {
	if u == nil || u.ListCollectionMethods == nil {
		return nil
	}
	resp, err := u.ListCollectionMethods(ctx, &collectionmethodpb.ListCollectionMethodsRequest{})
	if err != nil {
		return nil
	}
	out := make([]types.SelectOption, 0, len(resp.GetData()))
	for _, m := range resp.GetData() {
		out = append(out, types.SelectOption{Value: m.GetId(), Label: m.GetName()})
	}
	return out
}
