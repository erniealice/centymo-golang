// Package drawer implements the per-supplier (Surface A) and
// per-supplier-subscription (Surface C) expense-recognition-run drawer
// handlers. The same backing handler serves both surfaces — the scope
// discriminator (supplier_id vs supplier_subscription_id) is taken from the
// URL path.
//
// Plan A 20260517-expense-run Phase 4.
package drawer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	errshared "github.com/erniealice/centymo-golang/views/expense_recognition_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Scope discriminates the drawer mode.
type Scope int

const (
	ScopeSupplier Scope = iota
	ScopeSubscription
)

// CandidateRow is a single per-period candidate displayed in the drawer.
type CandidateRow struct {
	SourceKind             string // "subscription" | "advance_disbursement"
	SourceLabel            string
	SupplierSubscriptionID string
	AdvanceDisbursementID  string
	PeriodStart            string
	PeriodEnd              string
	PeriodLabel            string
	PeriodMarker           string
	Amount                 int64
	Currency               string
	Eligible               bool
	BlockerReason          string
	SuppressingAdvanceID   string
}

// DrawerData is the page-data passed to the drawer template.
type DrawerData struct {
	types.PageData
	Scope         string // "supplier" | "subscription"
	ScopeID       string
	AsOfDate      string
	AsOfDateMax   string
	FormAction    string // POST URL to Generate
	WorkspaceID    string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	FragmentURL   string // HTMX inner-swap URL for AsOfDate change
	SubscriptionCandidates []CandidateRow
	AdvanceCandidates      []CandidateRow
	SuppressedCandidates   []CandidateRow
	Labels                 centymo.ExpenseRecognitionRunLabels
	CommonLabels           pyeza.CommonLabels
	// AdvanceSuppressionNoticeURL is the link to an advance Disbursement when
	// the Surface C subscription drawer is fully suppressed by an advance.
	AdvanceSuppressionNoticeURL string
}

// Deps holds the consumer callbacks the drawer needs.
type Deps struct {
	Routes       centymo.ExpenseRecognitionRunRoutes
	Labels       centymo.ExpenseRecognitionRunLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ListExpenseRunCandidates returns candidates for the given scope.
	ListExpenseRunCandidates func(ctx context.Context, scope ListScope) ([]CandidateRow, error)

	// GenerateExpenseRun executes the run for the given scope + selections.
	GenerateExpenseRun func(ctx context.Context, in GenerateInput) (*GenerateOutput, error)
}

// ListScope is the input shape for the candidate listing callback.
type ListScope struct {
	SupplierID             string
	SupplierSubscriptionID string
	AsOfDate               string
}

// GenerateInput is the input shape for the Generate callback.
type GenerateInput struct {
	SupplierID             string
	SupplierSubscriptionID string
	AsOfDate               string
}

// GenerateOutput holds the run-level outcome aggregation.
type GenerateOutput = errshared.BatchRunOutput

// NewSupplierDrawer returns a GET handler that renders the Surface A drawer.
func NewSupplierDrawer(deps *Deps) view.View {
	return newDrawerView(deps, ScopeSupplier)
}

// NewSubscriptionDrawer returns a GET handler that renders the Surface C drawer.
func NewSubscriptionDrawer(deps *Deps) view.View {
	return newDrawerView(deps, ScopeSubscription)
}

// NewGenerateAction returns the POST handler that triggers GenerateExpenseRun.
// Reads scope from the form body (scope=supplier|subscription + id).
func NewGenerateAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense_recognition_run", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.Error(err)
		}
		if deps.GenerateExpenseRun == nil {
			return centymo.HTMXError(deps.Labels.Errors.UseCaseUnavailable)
		}

		in := GenerateInput{
			SupplierID:             viewCtx.Request.FormValue("supplier_id"),
			SupplierSubscriptionID: viewCtx.Request.FormValue("supplier_subscription_id"),
			AsOfDate:               viewCtx.Request.FormValue("as_of_date"),
		}
		out, err := deps.GenerateExpenseRun(ctx, in)
		if err != nil {
			log.Printf("expense-recognition-run drawer: GenerateExpenseRun error: %v", err)
			return centymo.HTMXError(deps.Labels.Errors.GenerationFailed)
		}

		trigger := buildToastTrigger(deps, out)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger": trigger,
			},
		}
	})
}

func newDrawerView(deps *Deps, scope Scope) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense_recognition_run", "create") {
			return view.Forbidden("expense_recognition_run:create")
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidSelection)
		}

		asOfDate := viewCtx.Request.URL.Query().Get("as_of_date")
		if asOfDate == "" {
			asOfDate = time.Now().UTC().Format("2006-01-02")
		}

		scopeInput := ListScope{AsOfDate: asOfDate}
		scopeName := "supplier"
		formAction := deps.Routes.GenerateURL
		fragmentURL := strings.ReplaceAll(deps.Routes.PerSupplierDrawerURL, "{id}", id)
		if scope == ScopeSubscription {
			scopeInput.SupplierSubscriptionID = id
			scopeName = "subscription"
			fragmentURL = strings.ReplaceAll(deps.Routes.PerSubscriptionDrawerURL, "{id}", id)
		} else {
			scopeInput.SupplierID = id
		}

		var allCandidates []CandidateRow
		if deps.ListExpenseRunCandidates != nil {
			var err error
			allCandidates, err = deps.ListExpenseRunCandidates(ctx, scopeInput)
			if err != nil {
				log.Printf("expense-recognition-run drawer: ListExpenseRunCandidates error: %v", err)
				return view.Error(fmt.Errorf("failed to load candidates: %w", err))
			}
		}

		// Split candidates into three buckets.
		var subs, advs, suppressed []CandidateRow
		for _, c := range allCandidates {
			if !c.Eligible && c.BlockerReason == "suppressed_by_advance" {
				suppressed = append(suppressed, c)
				continue
			}
			switch c.SourceKind {
			case "subscription", "subscription_cycle":
				subs = append(subs, c)
			case "advance_disbursement":
				advs = append(advs, c)
			default:
				// Unknown kind — leave on the subscription side as a conservative fallback.
				subs = append(subs, c)
			}
		}

		// Surface C only renders subscription candidates; Surface A renders both.
		if scope == ScopeSubscription {
			advs = nil
		}

		data := &DrawerData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Scope:                  scopeName,
			ScopeID:                id,
			AsOfDate:               asOfDate,
			AsOfDateMax:            time.Now().UTC().Format("2006-01-02"),
			FormAction:             formAction,
			FragmentURL:            fragmentURL,
			SubscriptionCandidates: subs,
			AdvanceCandidates:      advs,
			SuppressedCandidates:   suppressed,
			Labels:                 deps.Labels,
			CommonLabels:           deps.CommonLabels,
		}

		return view.OK("expense-recognition-run-drawer", data)
	})
}

func buildToastTrigger(deps *Deps, out *GenerateOutput) string {
	labels := deps.Labels
	created, skipped, errored := 0, 0, 0
	runID := ""
	if out != nil {
		created = out.Created
		skipped = out.Skipped
		errored = out.Errored
		runID = out.RunID
	}
	total := created + skipped + errored
	message := strings.NewReplacer(
		"{{.Created}}", fmt.Sprintf("%d", created),
		"{{.Skipped}}", fmt.Sprintf("%d", skipped),
		"{{.Errored}}", fmt.Sprintf("%d", errored),
		"{{.Total}}", fmt.Sprintf("%d", total),
	).Replace(labels.Toast.Success)

	state := "success"
	if errored > 0 && created == 0 {
		state = "error"
	} else if errored > 0 {
		state = "warning"
	}

	toast := map[string]any{
		"message": message,
		"state":   state,
	}
	if runID != "" && deps.Routes.DetailURL != "" && labels.Toast.ViewRunLink != "" {
		runDetailURL := strings.ReplaceAll(deps.Routes.DetailURL, "{id}", runID)
		toast["link"] = map[string]any{
			"url":   runDetailURL,
			"label": labels.Toast.ViewRunLink,
		}
	}
	envelope, err := json.Marshal(map[string]any{
		"pyeza:toast":                toast,
		"expense-recognitions-table": map[string]any{}, // refresh tab on subscription detail
	})
	if err != nil {
		return `{"pyeza:toast":{"message":"","state":"info"}}`
	}
	return string(envelope)
}
