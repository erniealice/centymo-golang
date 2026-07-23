package block

// subscription_group_options.go — the app-configurable roster-banding surface
// for the subscription_group enrollments tab: the EngineBlock option that carries
// the "client_attributes.<code>" grammar, the client-attribute use-case group the
// banding reads, and the use-case wiring that binds them (including the
// ambiguity-safe attribute-code resolver). Kept here rather than in engineblock.go
// / usecases.go so both of those stay under the god-file threshold.

import (
	"context"
	"fmt"

	subscriptiongrouppkg "github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	"github.com/erniealice/espyna-golang/consumer"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
)

// EngineOption configures the engine block from the consuming app (the app's
// view option block). Options are forwarded verbatim to AllUnits.
type EngineOption func(*engineConfig)

// engineConfig collects the per-unit view options an app may set.
type engineConfig struct {
	subscriptionGroupOptions subscriptiongrouppkg.Options
}

// WithSubscriptionGroupOptions sets the subscription_group roster (enrollments
// tab) presentation options: value bands + within-band sort, through the generic
// "client_attributes.<code>" grammar. The zero value renders the current flat
// roster unchanged (backward-compatible for consumers — e.g. service-admin —
// that do not set it).
func WithSubscriptionGroupOptions(o subscriptiongrouppkg.Options) EngineOption {
	return func(c *engineConfig) { c.subscriptionGroupOptions = o }
}

// bindClientAttributeUseCases wires the roster-banding closures from the raw
// aggregate onto the adapted result. The value list lives on the Entity
// aggregate; the code→id resolver is derived from Common's attribute list. Both
// nil-safe — an unwired closure degrades the optioned roster to flat.
func bindClientAttributeUseCases(result *UseCases, uc *consumer.UseCases) {
	if result == nil || uc == nil {
		return
	}
	if uc.Entity != nil && uc.Entity.ClientAttribute != nil && uc.Entity.ClientAttribute.ListClientAttributes != nil {
		result.Entity.ClientAttribute.ListClientAttributes = uc.Entity.ClientAttribute.ListClientAttributes.Execute
	}
	if uc.Common != nil && uc.Common.Attribute != nil && uc.Common.Attribute.ListAttributes != nil {
		result.Entity.ClientAttribute.ResolveAttributeIDByCode = unambiguousAttributeResolver(uc.Common.Attribute.ListAttributes.Execute)
	}
}

// unambiguousAttributeResolver returns a code→id resolver that fails flat unless
// exactly ONE active definition matches the code. attribute.code carries no
// unique constraint, so two active definitions sharing a code (or none) would
// otherwise silently pick a row and mis-band the roster; returning an error lets
// the caller fall back to the flat render.
func unambiguousAttributeResolver(
	listAttributes func(context.Context, *commonpb.ListAttributesRequest) (*commonpb.ListAttributesResponse, error),
) func(context.Context, string) (string, error) {
	return func(ctx context.Context, code string) (string, error) {
		if listAttributes == nil {
			return "", fmt.Errorf("attribute list use case not initialized")
		}
		resp, err := listAttributes(ctx, &commonpb.ListAttributesRequest{
			Filters: &commonpb.FilterRequest{
				Filters: []*commonpb.TypedFilter{{
					Field: "code",
					FilterType: &commonpb.TypedFilter_StringFilter{
						StringFilter: &commonpb.StringFilter{Value: code, Operator: commonpb.StringOperator_STRING_EQUALS},
					},
				}},
			},
		})
		if err != nil {
			return "", fmt.Errorf("list attributes by code %q: %w", code, err)
		}
		var id string
		matches := 0
		for _, a := range resp.GetData() {
			if a == nil || !a.GetActive() || a.GetCode() != code {
				continue
			}
			matches++
			id = a.GetId()
		}
		if matches != 1 {
			return "", fmt.Errorf("attribute code %q resolved to %d active definitions, want exactly 1", code, matches)
		}
		return id, nil
	}
}
