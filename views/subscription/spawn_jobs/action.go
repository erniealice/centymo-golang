// Package spawn_jobs handles the retroactive "Spawn Jobs" feature for subscriptions.
// Drawer template: subscription-spawn-jobs-drawer-form.html (stays flat at view root).
package spawn_jobs

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	spawnform "github.com/erniealice/centymo-golang/views/subscription/spawn_jobs/form"

	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// Deps is the dependency subset needed by the spawn_jobs feature.
type Deps struct {
	Routes centymo.SubscriptionRoutes
	Labels centymo.SubscriptionLabels

	ReadSubscription               func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
	MaterializeJobsForSubscription func(ctx context.Context, subscriptionID string, spawnJobs bool) (jobCount int, skippedReason string, err error)

	// DetectSpawnJobs resolves PricePlan → Plan → JobTemplate for the drawer.
	// Provided by the action package since it requires full action.Deps.
	DetectSpawnJobs func(ctx context.Context, pricePlanID string) DetectionResult
}

// DetectionResult carries the spawn-jobs detection outcome.
type DetectionResult struct {
	Available     bool
	TemplateNames []string
	JobCount      int
	PhaseCount    int
	TaskCount     int
}

// NewAction handles GET /action/subscription/{subscriptionId}/spawn-jobs
// (drawer) and POST (commit). 2026-04-29 auto-spawn-jobs-from-subscription
// plan §5.3.
func NewAction(deps *Deps) view.View {
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
			det, rows, rootName := detectDrawer(ctx, deps, subscriptionID)
			subLabel := resolveSubscriptionLabel(ctx, deps, subscriptionID)
			drawerLabels := spawnform.Labels{
				Title:             deps.Labels.Spawn.Title,
				DetectedTemplates: deps.Labels.Spawn.DetectedTemplates,
				RootTemplate:      deps.Labels.Spawn.RootTemplate,
				Cancel:            deps.Labels.Spawn.Cancel,
				Confirm:           deps.Labels.Spawn.Confirm,
				Skipped:           deps.Labels.Spawn.Skipped,
			}
			return view.OK("subscription-spawn-jobs-drawer-form", &spawnform.Data{
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

// detectDrawer extends the detection result with per-template phase/task counts,
// returning rows suitable for the drawer template.
func detectDrawer(ctx context.Context, deps *Deps, subscriptionID string) (DetectionResult, []spawnform.TemplateRow, string) {
	if deps == nil || deps.ReadSubscription == nil {
		return DetectionResult{}, nil, ""
	}
	subResp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
		Data: &subscriptionpb.Subscription{Id: subscriptionID},
	})
	if err != nil || subResp == nil || len(subResp.GetData()) == 0 {
		return DetectionResult{}, nil, ""
	}
	pricePlanID := subResp.GetData()[0].GetPricePlanId()
	var det DetectionResult
	if deps.DetectSpawnJobs != nil {
		det = deps.DetectSpawnJobs(ctx, pricePlanID)
	}
	rows := make([]spawnform.TemplateRow, 0, len(det.TemplateNames))
	for i, name := range det.TemplateNames {
		rows = append(rows, spawnform.TemplateRow{
			TemplateID:   name,
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
