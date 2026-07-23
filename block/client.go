package block

// client.go — the client-entity use-case type groups (client identity + client-
// attribute banding) that the subscription_group roster and other client-facing
// views read. Split out of usecases.go to keep that file under the god-file
// threshold.

import (
	"context"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientattributepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_attribute"
)

type ClientUseCases struct {
	ListClients         func(context.Context, *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ReadClient          func(context.Context, *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	SearchClientsByName func(context.Context, *clientpb.SearchClientsByNameRequest) (*clientpb.SearchClientsByNameResponse, error)
}

// ClientAttributeUseCases backs the roster's generic client-attribute banding
// (the "client_attributes.<code>" option). Both closures are nil-safe — an
// unwired one degrades the optioned roster to sorted-flat, never errors.
type ClientAttributeUseCases struct {
	ListClientAttributes     func(context.Context, *clientattributepb.ListClientAttributesRequest) (*clientattributepb.ListClientAttributesResponse, error)
	ResolveAttributeIDByCode func(ctx context.Context, code string) (string, error)
}
