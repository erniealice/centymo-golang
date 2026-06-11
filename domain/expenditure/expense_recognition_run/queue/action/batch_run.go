// Package action implements the POST handler for the
// expense-recognition-run batch-submit endpoint (Surface B Phase 4).
//
// Mirror of packages/centymo-golang/views/revenue_run/queue/action/batch_run.go.
// Plan A 20260517-expense-run.
package action

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run"

	errshared "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run/shared"
	"github.com/erniealice/pyeza-golang/view"
)

// BatchRunResult holds the aggregate outcome of a batch run submission.
type BatchRunResult struct {
	RunIDs  []string
	Created int
	Skipped int
	Errored int
}

// BatchRunDeps holds dependencies for the batch-run POST handler.
type BatchRunDeps struct {
	Routes expense_recognition_run.Routes
	Labels expense_recognition_run.Labels

	// GenerateExpenseRun executes the expense run for a single supplier.
	// Returns a minimal output struct; block.go shim translates from consumer types.
	GenerateExpenseRun func(ctx context.Context, in errshared.BatchRunInput) (*errshared.BatchRunOutput, error)
}

// NewBatchRunAction creates the POST handler for ExpenseRecognitionRunSubmitBatchURL.
func NewBatchRunAction(deps *BatchRunDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.Error(err)
		}

		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense_recognition_run", "create") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		selectionMode := viewCtx.Request.FormValue("selection_mode")
		asOfDate := viewCtx.Request.FormValue("as_of_date")

		switch selectionMode {
		case "all_matching":
			// FilterToken path — deferred. Stub-reject with label string.
			return view.HTMXError(deps.Labels.Errors.RunAllMatchingNotImplemented)

		case "selected":
			return handleRunSelected(ctx, viewCtx, deps, asOfDate)

		default:
			return view.HTMXError(deps.Labels.Errors.InvalidSelection)
		}
	})
}

// handleRunSelected processes the "run for selected suppliers" path.
func handleRunSelected(
	ctx context.Context,
	viewCtx *view.ViewContext,
	deps *BatchRunDeps,
	asOfDate string,
) view.ViewResult {
	l := deps.Labels

	// bulk-action.js posts repeated "id" fields (pyeza convention).
	supplierIDs := viewCtx.Request.Form["id"]
	supplierIDs = deduplicateStrings(supplierIDs)
	if len(supplierIDs) == 0 {
		return view.HTMXError(l.Errors.InvalidSelection)
	}
	if len(supplierIDs) > 50 {
		return view.HTMXError(l.Errors.CapExceeded)
	}

	if deps.GenerateExpenseRun == nil {
		return view.HTMXError(l.Errors.UseCaseUnavailable)
	}

	var result BatchRunResult
	for _, supplierID := range supplierIDs {
		out, err := deps.GenerateExpenseRun(ctx, errshared.BatchRunInput{
			SupplierID: supplierID,
			AsOfDate:   asOfDate,
		})
		if err != nil {
			log.Printf("expense-recognition-run batch: GenerateExpenseRun error for supplier %s: %v", supplierID, err)
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

	triggerPayload := buildBatchCompletedTrigger(deps, result)

	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": triggerPayload,
		},
	}
}

// buildBatchCompletedTrigger encodes the batch result as an HX-Trigger JSON
// payload using the shared "pyeza:toast" event shape.
func buildBatchCompletedTrigger(deps *BatchRunDeps, result BatchRunResult) string {
	labels := deps.Labels
	message := strings.NewReplacer(
		"{{.Created}}", fmt.Sprintf("%d", result.Created),
		"{{.Skipped}}", fmt.Sprintf("%d", result.Skipped),
		"{{.Errored}}", fmt.Sprintf("%d", result.Errored),
		"{{.RunCount}}", fmt.Sprintf("%d", len(result.RunIDs)),
	).Replace(labels.Toast.BatchSuccess)

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
	// Single-run batch → link to that run's detail.
	if len(result.RunIDs) == 1 && labels.Toast.ViewRunLink != "" && deps.Routes.DetailURL != "" {
		runDetailURL := strings.ReplaceAll(deps.Routes.DetailURL, "{id}", result.RunIDs[0])
		toast["link"] = map[string]any{
			"url":   runDetailURL,
			"label": labels.Toast.ViewRunLink,
		}
	}

	envelope, err := json.Marshal(map[string]any{
		"pyeza:toast": toast,
	})
	if err != nil {
		return `{"pyeza:toast":{"message":"","state":"info"}}`
	}
	return string(envelope)
}

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
