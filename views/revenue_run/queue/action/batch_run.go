// Package action implements the POST handler for the revenue-run batch-submit
// endpoint (Surface B — workspace queue page, Phase 7).
package action

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// BatchRunResult holds the aggregate outcome of a batch run submission.
type BatchRunResult struct {
	RunIDs  []string
	Created int
	Skipped int
	Errored int
}

// GenerateRevenueRunInput is the minimal input shape for one client's run.
type GenerateRevenueRunInput struct {
	ClientID string
	AsOfDate string
}

// GenerateRevenueRunOutput is the minimal output shape from one client's run.
type GenerateRevenueRunOutput struct {
	RunID   string
	Created int
	Skipped int
	Errored int
}

// BatchRunDeps holds dependencies for the batch-run POST handler.
type BatchRunDeps struct {
	Routes centymo.RevenueRunRoutes
	Labels centymo.RevenueRunLabels

	// GenerateRevenueRun executes the revenue run for a single client.
	// Returns a minimal output struct; block.go shim translates from consumer types.
	GenerateRevenueRun func(ctx context.Context, in GenerateRevenueRunInput) (*GenerateRevenueRunOutput, error)
}

// NewBatchRunAction creates the POST handler for RevenueRunSubmitBatchURL.
func NewBatchRunAction(deps *BatchRunDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.Error(err)
		}

		perms := view.GetUserPermissions(ctx)
		if !perms.Can("revenue", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		selectionMode := viewCtx.Request.FormValue("selection_mode")
		asOfDate := viewCtx.Request.FormValue("as_of_date")

		switch selectionMode {
		case "all_matching":
			// FilterToken path — deferred. Stub-reject with label string (Wave 3).
			return centymo.HTMXError(deps.Labels.Errors.RunAllMatchingNotImplemented)

		case "selected":
			return handleRunSelected(ctx, viewCtx, deps, asOfDate)

		default:
			return centymo.HTMXError(deps.Labels.Errors.InvalidSelection)
		}
	})
}

// handleRunSelected processes the "run for selected clients" path.
func handleRunSelected(
	ctx context.Context,
	viewCtx *view.ViewContext,
	deps *BatchRunDeps,
	asOfDate string,
) view.ViewResult {
	l := deps.Labels

	// bulk-action.js posts repeated "id" fields (the pyeza convention used by
	// every other bulk action handler in the codebase).
	clientIDs := viewCtx.Request.Form["id"]
	// Deduplicate and validate.
	clientIDs = deduplicateStrings(clientIDs)
	if len(clientIDs) == 0 {
		return centymo.HTMXError(l.Errors.InvalidSelection)
	}
	if len(clientIDs) > 50 {
		return centymo.HTMXError(l.Errors.CapExceeded)
	}

	if deps.GenerateRevenueRun == nil {
		return centymo.HTMXError(l.Errors.UseCaseUnavailable)
	}

	var result BatchRunResult
	for _, clientID := range clientIDs {
		out, err := deps.GenerateRevenueRun(ctx, GenerateRevenueRunInput{
			ClientID: clientID,
			AsOfDate: asOfDate,
		})
		if err != nil {
			log.Printf("revenue-run batch: GenerateRevenueRun error for client %s: %v", clientID, err)
			result.Errored++
			continue
		}
		if out == nil {
			result.Skipped++
			continue
		}
		if out.RunID != "" {
			result.RunIDs = append(result.RunIDs, out.RunID)
		}
		result.Created += out.Created
		result.Skipped += out.Skipped
		result.Errored += out.Errored
	}

	// Build HX-Trigger payload — emit a structured pyeza:toast event that the
	// centralized lf.Toast module (pyeza/toast.js) renders directly. Single-run
	// batches carry a "View run" link; multi-run batches omit it.
	triggerPayload := buildBatchCompletedTrigger(deps.Labels, result)

	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": triggerPayload,
		},
	}
}

// buildBatchCompletedTrigger encodes the batch result as an HX-Trigger JSON
// payload using the shared "pyeza:toast" event shape. The toast text is
// fully resolved server-side (lyngua substitution) so JS receives only
// rendered strings.
func buildBatchCompletedTrigger(labels centymo.RevenueRunLabels, result BatchRunResult) string {
	message := strings.NewReplacer(
		"{{.Created}}", fmt.Sprintf("%d", result.Created),
		"{{.Skipped}}", fmt.Sprintf("%d", result.Skipped),
		"{{.Errored}}", fmt.Sprintf("%d", result.Errored),
	).Replace(labels.ToastBatchSuccess)

	state := "success"
	if result.Errored > 0 && result.Created == 0 {
		state = "error"
	} else if result.Errored > 0 {
		state = "warning"
	}

	toast := map[string]any{
		"message": message,
		"state":   state,
	}
	// Single-run batch → link straight to the run detail. Multi-run batches
	// omit the link (no good single target; sidebar Revenue Run history serves).
	if len(result.RunIDs) == 1 && labels.ViewRunLink != "" {
		runDetailURL := strings.ReplaceAll(centymo.RevenueRunDetailURL, "{id}", result.RunIDs[0])
		toast["link"] = map[string]any{
			"url":   runDetailURL,
			"label": labels.ViewRunLink,
		}
	}

	envelope, err := json.Marshal(map[string]any{
		"pyeza:toast": toast,
	})
	if err != nil {
		// Fallback — never happens with plain map values.
		return `{"pyeza:toast":{"message":"","state":"info"}}`
	}
	return string(envelope)
}

// deduplicateStrings returns a slice with duplicate strings removed.
func deduplicateStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
