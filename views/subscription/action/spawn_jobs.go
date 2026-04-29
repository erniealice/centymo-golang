package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// SpawnJobsDrawerData is the template data for the retroactive spawn drawer.
// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.3.
type SpawnJobsDrawerData struct {
	FormAction        string
	SubscriptionID    string
	SubscriptionLabel string
	// Detected templates (root + active children). RootName highlights which
	// template is the root for the operator radio group. Empty Templates =
	// nothing to spawn (operator sees the skipped notice).
	Templates  []SpawnTemplateRow
	RootName   string
	HasContent bool

	// Resolved labels for the drawer (avoid label drift across tiers).
	Labels SpawnDrawerLabels
	// Common labels supplied by the view adapter (cancel/save buttons).
	CommonLabels any
}

// SpawnTemplateRow is one detected JobTemplate rendered in the drawer.
type SpawnTemplateRow struct {
	TemplateID   string
	TemplateName string
	IsRoot       bool
	PhaseCount   int
	TaskCount    int
}

// SpawnDrawerLabels carries the typed strings consumed by the drawer template.
// 2026-04-29 auto-spawn-jobs-from-subscription plan §9.3.
type SpawnDrawerLabels struct {
	Title             string
	DetectedTemplates string
	RootTemplate      string
	Cancel            string
	Confirm           string
	Skipped           string
}

// NewSpawnJobsPartialAction handles GET /action/subscription/_partial/spawn-jobs-section.
// It re-renders the Spawn Jobs section inside the create drawer when the
// operator changes the selected Plan/PricePlan. Reads `price_plan_id` from
// the query string (auto-complete trigger sends it via hx-include).
//
// Returns the section HTML wrapped in #spawn-jobs-wrapper's inner content
// (the wrapper itself is preserved on the client).
func NewSpawnJobsPartialAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pricePlanID := viewCtx.Request.URL.Query().Get("price_plan_id")
		labels := formLabels(deps.Labels)

		det := detectSpawnJobs(ctx, deps, pricePlanID)
		summary := resolveSpawnJobsSummary(labels.SpawnJobsSummary, det)

		data := &FormData{
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

// NewSpawnJobsAction handles GET /action/subscription/{subscriptionId}/spawn-jobs
// (drawer) and POST (commit). 2026-04-29 auto-spawn-jobs-from-subscription
// plan §5.3.
func NewSpawnJobsAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if perms != nil && !perms.Can("subscription", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		subscriptionID := viewCtx.Request.PathValue("subscriptionId")
		if subscriptionID == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		formAction := strings.ReplaceAll(deps.Routes.SpawnJobsURL, "{subscriptionId}", subscriptionID)

		if viewCtx.Request.Method == http.MethodGet {
			det, rows, rootName := detectSpawnJobsDrawer(ctx, deps, subscriptionID)
			subLabel := resolveSubscriptionLabel(ctx, deps, subscriptionID)
			drawerLabels := SpawnDrawerLabels{
				Title:             deps.Labels.Spawn.Title,
				DetectedTemplates: deps.Labels.Spawn.DetectedTemplates,
				RootTemplate:      deps.Labels.Spawn.RootTemplate,
				Cancel:            deps.Labels.Spawn.Cancel,
				Confirm:           deps.Labels.Spawn.Confirm,
				Skipped:           deps.Labels.Spawn.Skipped,
			}
			return view.OK("subscription-spawn-jobs-drawer-form", &SpawnJobsDrawerData{
				FormAction:        formAction,
				SubscriptionID:    subscriptionID,
				SubscriptionLabel: subLabel,
				Templates:         rows,
				RootName:          rootName,
				HasContent:        det.Available,
				Labels:            drawerLabels,
				CommonLabels:      nil,
			})
		}

		// POST — invoke the use case.
		if deps.MaterializeJobsForSubscription == nil {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		_, _, err := deps.MaterializeJobsForSubscription(ctx, subscriptionID, true)
		if err != nil {
			log.Printf("Failed to spawn jobs for subscription %s: %v", subscriptionID, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("subscription-operations-tab")
	})
}

// detectSpawnJobsDrawer extends detectSpawnJobs with the per-template phase /
// task counts already aggregated, returning rows suitable for the drawer
// template.
func detectSpawnJobsDrawer(ctx context.Context, deps *Deps, subscriptionID string) (SpawnJobsDetection, []SpawnTemplateRow, string) {
	if deps == nil || deps.ReadSubscription == nil {
		return SpawnJobsDetection{}, nil, ""
	}
	subResp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
		Data: &subscriptionpb.Subscription{Id: subscriptionID},
	})
	if err != nil || subResp == nil || len(subResp.GetData()) == 0 {
		return SpawnJobsDetection{}, nil, ""
	}
	pricePlanID := subResp.GetData()[0].GetPricePlanId()
	det := detectSpawnJobs(ctx, deps, pricePlanID)
	rows := make([]SpawnTemplateRow, 0, len(det.TemplateNames))
	for i, name := range det.TemplateNames {
		rows = append(rows, SpawnTemplateRow{
			TemplateID:   name, // template names are the only stable handle exposed by detectSpawnJobs; v2 enrichment can carry the ID separately
			TemplateName: name,
			IsRoot:       i == 0,
		})
	}
	root := ""
	if len(rows) > 0 {
		root = rows[0].TemplateName
	}
	return det, rows, root
}

// resolveSubscriptionLabel returns a "code · name" label for the drawer header.
func resolveSubscriptionLabel(ctx context.Context, deps *Deps, subscriptionID string) string {
	if deps == nil || deps.ReadSubscription == nil {
		return subscriptionID
	}
	resp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
		Data: &subscriptionpb.Subscription{Id: subscriptionID},
	})
	if err != nil || resp == nil || len(resp.GetData()) == 0 {
		return subscriptionID
	}
	s := resp.GetData()[0]
	if c := s.GetCode(); c != "" {
		if n := s.GetName(); n != "" {
			return c + " · " + n
		}
		return c
	}
	if n := s.GetName(); n != "" {
		return n
	}
	return subscriptionID
}
