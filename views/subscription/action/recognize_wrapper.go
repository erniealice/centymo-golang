package action

// recognize_wrapper.go provides backward-compatible shim constructors that
// keep block.go's subscriptionaction.New* call sites unchanged while the
// implementation now lives in the recognize/ sub-package.

import (
	"github.com/erniealice/pyeza-golang/view"

	recognizepkg "github.com/erniealice/centymo-golang/views/subscription/recognize"
)

// NewRecognizeAction is the backward-compatible shim for block.go.
// Delegates to recognize.NewAction using a sub-set of action.Deps.
func NewRecognizeAction(deps *Deps) view.View {
	return recognizepkg.NewAction(&recognizepkg.Deps{
		Routes:                           deps.Routes,
		Labels:                           deps.Labels,
		ReadSubscription:                 deps.ReadSubscription,
		ListClients:                      deps.ListClients,
		ReadPricePlan:                    deps.ReadPricePlan,
		RecognizeRevenueFromSubscription: deps.RecognizeRevenueFromSubscription,
		ListBillingEventsBySubscription:  deps.ListBillingEventsBySubscription,
	})
}
