package action

// search_wrapper.go provides backward-compatible shim constructors that keep
// block.go's subscriptionaction.NewSearchClientsAction and
// NewSearchPlansAction call sites unchanged while the implementation now lives
// in the search/ sub-package.

import (
	"net/http"

	searchpkg "github.com/erniealice/centymo-golang/views/subscription/search"
)

// NewSearchClientsAction is the backward-compatible shim for block.go.
func NewSearchClientsAction(deps *Deps) http.HandlerFunc {
	return searchpkg.NewSearchClientsAction(&searchpkg.Deps{
		SearchClientsByName: deps.SearchClientsByName,
		ListClients:         deps.ListClients,
		ListPricePlans:      deps.ListPricePlans,
		ListPriceSchedules:  deps.ListPriceSchedules,
		ListPlans:           deps.ListPlans,
		ReadPlan:            deps.ReadPlan,
	})
}

// NewSearchPlansAction is the backward-compatible shim for block.go.
func NewSearchPlansAction(deps *Deps) http.HandlerFunc {
	return searchpkg.NewSearchPlansAction(&searchpkg.Deps{
		SearchClientsByName: deps.SearchClientsByName,
		ListClients:         deps.ListClients,
		ListPricePlans:      deps.ListPricePlans,
		ListPriceSchedules:  deps.ListPriceSchedules,
		ListPlans:           deps.ListPlans,
		ReadPlan:            deps.ReadPlan,
	})
}
