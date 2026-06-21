// Package block — revenue-run domain wiring (Surface D + Surface B).
//
// Holds wireRevenueRunModule (the lifted body of the `if cfg.wantRevenueRun()`
// branch of Block()) plus the proto<->view translators that only the
// revenue-run wiring calls. Co-locating translator + caller keeps the
// row-shape converters on the next page when a reader is following the
// wiring code.
//
// Phase 4 of the 20260506-subscription-invoice-run plan.
package block

import (
	"context"
	"log"
	"time"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuerunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_run"

	consumerapp "github.com/erniealice/espyna-golang/consumer/app"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"

	revenuedomain "github.com/erniealice/centymo-golang/domain/revenue"
)

// revenueRunWiring holds everything wireRevenueRunModule needs from the
// surrounding Block() scope. More than 6 fields → use a struct per convention.
// Kept private; never re-exported.
type revenueRunWiring struct {
	revenueRunRoutes   revenuedomain.RevenueRunRoutes
	revenueRunLabels   revenuedomain.RevenueRunLabels
	revenueRoutes      revenuedomain.RevenueRoutes
	centymoTableLabels types.TableLabels
	uploadFile         func(context.Context, string, string, []byte, string) error
	listAttachments    func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment   func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment   func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID    func() string
}

// wireRevenueRunModule lifts the body of `if cfg.wantRevenueRun()` from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once when cfg.wantRevenueRun().
func wireRevenueRunModule(ctx *consumerapp.AppContext, cfg *blockConfig, useCases *UseCases, w revenueRunWiring) {
	rrDeps := &revenuedomain.RevenueRunModuleDeps{
		Routes:       w.revenueRunRoutes,
		Labels:       w.revenueRunLabels,
		CommonLabels: ctx.Common,
		TableLabels:  w.centymoTableLabels,
	}

	// Wire ListRevenueRuns — translate proto response to view-typed rows.
	rrDeps.ListRevenueRuns = func(fctx context.Context, scope revenuedomain.ListRevenueRunsScope) ([]revenuedomain.RevenueRunRow, string, error) {
		req := &revenuerunpb.ListRevenueRunsRequest{}
		resp, err := useCases.RevenueRun.ListRevenueRuns(fctx, req)
		if err != nil {
			return nil, "", err
		}
		if resp == nil {
			return []revenuedomain.RevenueRunRow{}, "", nil
		}
		rows := make([]revenuedomain.RevenueRunRow, 0, len(resp.GetData()))
		for _, r := range resp.GetData() {
			row := protoRevenueRunToRow(r)
			// Apply status filter (the proto service may not support it directly yet)
			if scope.Status != "" && row.Status != scope.Status {
				continue
			}
			rows = append(rows, row)
		}
		return rows, "", nil
	}

	// Wire ReadRevenueRun — translate proto response to view-typed struct.
	rrDeps.ReadRevenueRun = func(fctx context.Context, id string) (*revenuedomain.RevenueRunWithAttempts, error) {
		runID := id
		resp, err := useCases.RevenueRun.ReadRevenueRun(fctx, &revenuerunpb.ReadRevenueRunRequest{
			Data: &revenuerunpb.RevenueRun{Id: runID},
		})
		if err != nil {
			return nil, err
		}
		if resp == nil || len(resp.GetData()) == 0 {
			return nil, nil
		}
		run := protoRevenueRunToRow(resp.GetData()[0])

		attResp, err := useCases.RevenueRun.ListRevenueRunAttempts(fctx, &revenuerunpb.ListRevenueRunAttemptsRequest{
			RunId: runID,
		})
		if err != nil {
			log.Printf("centymo.Block: failed to load attempts for run %s: %v", id, err)
			attResp = nil
		}
		var attempts []revenuedomain.RevenueRunAttemptRow
		if attResp != nil {
			attempts = make([]revenuedomain.RevenueRunAttemptRow, 0, len(attResp.GetData()))
			for _, a := range attResp.GetData() {
				attempts = append(attempts, protoRevenueRunAttemptToRow(a))
			}
		}
		return &revenuedomain.RevenueRunWithAttempts{Run: run, Attempts: attempts}, nil
	}

	// Wire ListRevenueByRunID — filter revenue list by run_id for the Invoices tab.
	if useCases.Revenue.GetListPageData != nil {
		revenueDetailURLPattern := w.revenueRoutes.DetailURL
		rrDeps.ListRevenueByRunID = func(fctx context.Context, runID string) ([]revenuedomain.RevenueRow, error) {
			resp, err := useCases.Revenue.GetListPageData(fctx, &revenuepb.GetRevenueListPageDataRequest{
				Filters: &commonpb.FilterRequest{
					Filters: []*commonpb.TypedFilter{
						{
							Field: "rv.run_id",
							FilterType: &commonpb.TypedFilter_StringFilter{
								StringFilter: &commonpb.StringFilter{
									Value:    runID,
									Operator: commonpb.StringOperator_STRING_EQUALS,
								},
							},
						},
					},
				},
			})
			if err != nil {
				return nil, err
			}
			rows := make([]revenuedomain.RevenueRow, 0, len(resp.GetRevenueList()))
			for _, rv := range resp.GetRevenueList() {
				detailURL := ""
				if revenueDetailURLPattern != "" {
					detailURL = route.ResolveURL(revenueDetailURLPattern, "id", rv.GetId())
				}
				rows = append(rows, revenuedomain.RevenueRow{
					ID:              rv.GetId(),
					ReferenceNumber: rv.GetReferenceNumber(),
					RevenueDate:     rv.GetRevenueDate(),
					TotalAmount:     int64(rv.GetTotalAmount()),
					Currency:        rv.GetCurrency(),
					Status:          rv.GetStatus(),
					DetailURL:       detailURL,
				})
			}
			return rows, nil
		}
	}

	// --------------------------------------------------------
	// Surface B — workspace queue page (Phase 7).
	// --------------------------------------------------------

	// URL templates from BlockOptions.
	rrDeps.ClientDetailURLTemplate = cfg.clientDetailURL
	rrDeps.ClientDrawerURLTemplate = cfg.clientRevenueRunDrawerURL

	// ListClients — reuse the existing client list use case.
	if useCases.Entity.Client.ListClients != nil {
		lc := useCases.Entity.Client.ListClients
		rrDeps.ListClients = func(fctx context.Context, cursor string) ([]revenuedomain.QueueClientRecord, string, error) {
			resp, err := lc(fctx, &clientpb.ListClientsRequest{})
			if err != nil {
				return nil, "", err
			}
			if resp == nil {
				return []revenuedomain.QueueClientRecord{}, "", nil
			}
			records := make([]revenuedomain.QueueClientRecord, 0, len(resp.GetData()))
			for _, c := range resp.GetData() {
				records = append(records, revenuedomain.QueueClientRecord{
					ID:   c.GetId(),
					Name: c.GetName(),
				})
			}
			return records, "", nil
		}
	}

	// ListRevenueRunCandidates — direct proto call (ex-consumer helper).
	if useCases.Revenue.ListRevenueRunCandidates != nil {
		rrDeps.ListRevenueRunCandidates = func(fctx context.Context, clientID, asOfDate string) ([]revenuedomain.QueueCandidateInput, error) {
			// Plan B Phase 5c — opt-in to advance Collection candidates by default.
			includeAdv := true
			resp, err := useCases.Revenue.ListRevenueRunCandidates(fctx, &revenuerunpb.ListRevenueRunCandidatesRequest{
				Scope: &revenuerunpb.RevenueRunScope{
					ClientId: &clientID,
					AsOfDate: &asOfDate,
				},
				IncludeAdvanceCollections: &includeAdv,
			})
			if err != nil {
				return nil, err
			}
			candidates := resp.GetData()
			out := make([]revenuedomain.QueueCandidateInput, 0, len(candidates))
			for _, c := range candidates {
				out = append(out, revenuedomain.QueueCandidateInput{
					SubscriptionID: c.GetSubscriptionId(),
					Currency:       c.GetCurrency(),
					Amount:         c.GetAmount(),
					Eligible:       c.GetEligible(),
				})
			}
			return out, nil
		}
	}

	// GenerateRevenueRun — direct proto call (ex-consumer helper).
	if useCases.Revenue.GenerateRevenueRun != nil {
		rrDeps.GenerateRevenueRun = func(fctx context.Context, in revenuedomain.BatchRunInput) (*revenuedomain.BatchRunOutput, error) {
			resp, err := useCases.Revenue.GenerateRevenueRun(fctx, &revenuerunpb.GenerateRevenueRunRequest{
				Scope: &revenuerunpb.RevenueRunScope{
					ClientId: &in.ClientID,
					AsOfDate: &in.AsOfDate,
				},
			})
			if err != nil {
				return nil, err
			}
			if resp == nil || resp.GetRun() == nil {
				return nil, nil
			}
			var created, skipped, errored int
			for _, a := range resp.GetAttempts() {
				switch a.GetOutcome() {
				case revenuerunpb.RevenueRunAttemptOutcome_REVENUE_RUN_ATTEMPT_OUTCOME_CREATED:
					created++
				case revenuerunpb.RevenueRunAttemptOutcome_REVENUE_RUN_ATTEMPT_OUTCOME_SKIPPED:
					skipped++
				case revenuerunpb.RevenueRunAttemptOutcome_REVENUE_RUN_ATTEMPT_OUTCOME_ERRORED:
					errored++
				}
			}
			return &revenuedomain.BatchRunOutput{
				RunID:   resp.GetRun().GetId(),
				Created: created,
				Skipped: skipped,
				Errored: errored,
			}, nil
		}
	}

	rrDeps.UploadFile = w.uploadFile
	rrDeps.ListAttachments = w.listAttachments
	rrDeps.CreateAttachment = w.createAttachment
	rrDeps.DeleteAttachment = w.deleteAttachment
	rrDeps.NewAttachmentID = w.newAttachmentID
	revenuedomain.NewRevenueRunModule(rrDeps).RegisterRoutes(ctx.Routes)
}

// ---------------------------------------------------------------------------
// Phase 4 — Revenue Run proto → view-type shim helpers
// ---------------------------------------------------------------------------

// revenueRunMillisToRFC3339 converts a proto epoch-millisecond int64 to
// RFC3339 UTC. Returns "" when ms is zero (field absent / unset).
func revenueRunMillisToRFC3339(ms int64) string {
	if ms == 0 {
		return ""
	}
	return time.UnixMilli(ms).UTC().Format(time.RFC3339)
}

// revenueRunScopeKindString maps the proto ScopeKind enum to a short
// lowercase string used by the view layer ("subscription", "client",
// "workspace", or "" for unspecified).
func revenueRunScopeKindString(sk revenuerunpb.RevenueRunScopeKind) string {
	switch sk {
	case revenuerunpb.RevenueRunScopeKind_REVENUE_RUN_SCOPE_KIND_SUBSCRIPTION:
		return "subscription"
	case revenuerunpb.RevenueRunScopeKind_REVENUE_RUN_SCOPE_KIND_CLIENT:
		return "client"
	case revenuerunpb.RevenueRunScopeKind_REVENUE_RUN_SCOPE_KIND_WORKSPACE:
		return "workspace"
	default:
		return ""
	}
}

// revenueRunStatusString maps the proto Status enum to the lowercase string
// expected by the view layer ("pending", "complete", "failed", or "").
func revenueRunStatusString(s revenuerunpb.RevenueRunStatus) string {
	switch s {
	case revenuerunpb.RevenueRunStatus_REVENUE_RUN_STATUS_PENDING:
		return "pending"
	case revenuerunpb.RevenueRunStatus_REVENUE_RUN_STATUS_COMPLETE:
		return "complete"
	case revenuerunpb.RevenueRunStatus_REVENUE_RUN_STATUS_FAILED:
		return "failed"
	default:
		return ""
	}
}

// revenueRunAttemptOutcomeString maps the proto Outcome enum to the lowercase
// string expected by the view layer ("created", "skipped", "errored", or "").
func revenueRunAttemptOutcomeString(o revenuerunpb.RevenueRunAttemptOutcome) string {
	switch o {
	case revenuerunpb.RevenueRunAttemptOutcome_REVENUE_RUN_ATTEMPT_OUTCOME_CREATED:
		return "created"
	case revenuerunpb.RevenueRunAttemptOutcome_REVENUE_RUN_ATTEMPT_OUTCOME_SKIPPED:
		return "skipped"
	case revenuerunpb.RevenueRunAttemptOutcome_REVENUE_RUN_ATTEMPT_OUTCOME_ERRORED:
		return "errored"
	default:
		return ""
	}
}

// protoRevenueRunToRow translates a *revenuerunpb.RevenueRun proto message
// to the view-typed revenuedomain.RevenueRunRow.
// IsStalePending is computed here: status=pending AND initiated_at is older
// than REVENUE_RUN_PENDING_STALE_MINUTES (default 5) minutes ago.
func protoRevenueRunToRow(r *revenuerunpb.RevenueRun) revenuedomain.RevenueRunRow {
	if r == nil {
		return revenuedomain.RevenueRunRow{}
	}
	initiatedAt := revenueRunMillisToRFC3339(r.GetInitiatedAt())
	completedAt := revenueRunMillisToRFC3339(r.GetCompletedAt())
	status := revenueRunStatusString(r.GetStatus())

	// Compute IsStalePending: pending run whose initiated_at is > 5 minutes ago.
	isPending := r.GetStatus() == revenuerunpb.RevenueRunStatus_REVENUE_RUN_STATUS_PENDING
	isStalePending := false
	if isPending && r.GetInitiatedAt() > 0 {
		age := time.Since(time.UnixMilli(r.GetInitiatedAt()))
		isStalePending = age > 5*time.Minute
	}

	return revenuedomain.RevenueRunRow{
		ID:             r.GetId(),
		ScopeKind:      revenueRunScopeKindString(r.GetScopeKind()),
		ClientID:       r.GetClientId(),
		SubscriptionID: r.GetSubscriptionId(),
		AsOfDate:       r.GetAsOfDate(),
		Initiator:      r.GetInitiatedBy(),
		InitiatedAt:    initiatedAt,
		CompletedAt:    completedAt,
		Status:         status,
		SelectionCount: r.GetSelectionCount(),
		CreatedCount:   r.GetCreatedCount(),
		SkippedCount:   r.GetSkippedCount(),
		ErroredCount:   r.GetErroredCount(),
		IsStalePending: isStalePending,
		Notes:          r.GetNotes(),
	}
}

// protoRevenueRunAttemptToRow translates a *revenuerunpb.RevenueRunAttempt
// proto message to the view-typed revenuedomain.RevenueRunAttemptRow.
func protoRevenueRunAttemptToRow(a *revenuerunpb.RevenueRunAttempt) revenuedomain.RevenueRunAttemptRow {
	if a == nil {
		return revenuedomain.RevenueRunAttemptRow{}
	}
	return revenuedomain.RevenueRunAttemptRow{
		ID:             a.GetId(),
		RunID:          a.GetRunId(),
		SubscriptionID: a.GetSubscriptionId(),
		PeriodStart:    a.GetPeriodStart(),
		PeriodEnd:      a.GetPeriodEnd(),
		PeriodMarker:   a.GetPeriodMarker(),
		AttemptedAt:    revenueRunMillisToRFC3339(a.GetAttemptedAt()),
		Outcome:        revenueRunAttemptOutcomeString(a.GetOutcome()),
		RevenueID:      a.GetRevenueId(),
		ErrorCode:      a.GetErrorCode(),
		ErrorMessage:   a.GetErrorMessage(),
	}
}
