package action

import (
	"context"

	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang/views/subscription/form"
)

// NewSpawnJobsPartialAction handles GET /action/subscription/_partial/spawn-jobs-section.
// It re-renders the Spawn Jobs section inside the create drawer when the
// operator changes the selected Plan/PricePlan. Reads `price_plan_id` from
// the query string (auto-complete trigger sends it via hx-include).
//
// Returns the section HTML wrapped in #spawn-jobs-wrapper's inner content
// (the wrapper itself is preserved on the client).
//
// This is an HTMX partial supporting the Add/Edit drawer — it stays in action/
// per S7: partials supporting parent CRUD handlers live in action/.
func NewSpawnJobsPartialAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pricePlanID := viewCtx.Request.URL.Query().Get("price_plan_id")
		labels := buildFormLabels(deps.Labels)

		det := detectSpawnJobs(ctx, deps, pricePlanID)
		summary := resolveSpawnJobsSummary(labels.SpawnJobsSummary, det)

		data := &form.Data{
			SpawnJobsAvailable:  det.Available,
			SpawnJobsDefault:    det.Available,
			SpawnJobsSummary:    summary,
			SpawnJobsPartialURL: deps.Routes.SpawnJobsPartialURL,
			Labels:              labels,
			CommonLabels:        nil,
		}
		return view.OK("subscription-spawn-jobs-section", data)
	})
}
