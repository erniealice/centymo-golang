package action

// customize_wrapper.go provides the backward-compatible shim constructor that
// keeps block.go's subscriptionaction.NewCustomizePackageAction call site
// unchanged while the implementation now lives in the customize/ sub-package.

import (
	"context"

	"github.com/erniealice/pyeza-golang/view"

	customizepkg "github.com/erniealice/centymo-golang/views/subscription/customize"
)

// NewCustomizePackageAction is the backward-compatible shim for block.go.
func NewCustomizePackageAction(deps *Deps) view.View {
	var customizeFn func(ctx context.Context, req *customizepkg.Request) (*customizepkg.Response, error)
	if deps.CustomizePlanForClient != nil {
		customizeFn = func(ctx context.Context, req *customizepkg.Request) (*customizepkg.Response, error) {
			resp, err := deps.CustomizePlanForClient(ctx, &CustomizePlanForClientRequest{
				SourcePlanID:      req.SourcePlanID,
				SourcePricePlanID: req.SourcePricePlanID,
				ClientID:          req.ClientID,
				SubscriptionID:    req.SubscriptionID,
				NewScheduleName:   req.NewScheduleName,
			})
			if err != nil || resp == nil {
				return nil, err
			}
			return &customizepkg.Response{
				NewPlanID:      resp.NewPlanID,
				NewPricePlanID: resp.NewPricePlanID,
				NewScheduleID:  resp.NewScheduleID,
				Reused:         resp.Reused,
			}, nil
		}
	}
	return customizepkg.NewAction(&customizepkg.Deps{
		Labels:                               deps.Labels,
		CustomClientPriceScheduleLabelSuffix: deps.CustomClientPriceScheduleLabelSuffix,
		CustomizePlanForClient:               customizeFn,
		GetSubscriptionItemPageData:          deps.GetSubscriptionItemPageData,
		ReadSubscription:                     deps.ReadSubscription,
		ReadPricePlan:                        deps.ReadPricePlan,
		ListClients:                          deps.ListClients,
	})
}
