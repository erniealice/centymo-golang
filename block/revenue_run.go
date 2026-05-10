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

	consumer "github.com/erniealice/espyna-golang/consumer"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	revenuerunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_run"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	revenuerunmod "github.com/erniealice/centymo-golang/views/revenue_run"
)

// revenueRunWiring holds everything wireRevenueRunModule needs from the
// surrounding Block() scope. More than 6 fields → use a struct per convention.
// Kept private; never re-exported.
type revenueRunWiring struct {
	revenueRunRoutes centymo.RevenueRunRoutes
	revenueRunLabels centymo.RevenueRunLabels
	revenueRoutes    centymo.RevenueRoutes
	centymoTableLabels types.TableLabels
	uploadFile       func(context.Context, string, string, []byte, string) error
	listAttachments  func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID  func() string
}

// wireRevenueRunModule lifts the body of `if cfg.wantRevenueRun()` from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once when cfg.wantRevenueRun().
func wireRevenueRunModule(ctx *pyeza.AppContext, cfg *blockConfig, useCases *consumer.UseCases, w revenueRunWiring) {
	rrDeps := &revenuerunmod.ModuleDeps{
		Routes:       w.revenueRunRoutes,
		Labels:       w.revenueRunLabels,
		CommonLabels: ctx.Common,
		TableLabels:  w.centymoTableLabels,
	}

	// Wire ListRevenueRuns — translate proto response to view-typed rows.
	rrDeps.ListRevenueRuns = func(fctx context.Context, scope revenuerunmod.ListRevenueRunsScope) ([]revenuerunmod.RevenueRunRow, string, error) {
		req := &revenuerunpb.ListRevenueRunsRequest{}
		resp, err := consumer.ListRevenueRuns(useCases, fctx, req)
		if err != nil {
			return nil, "", err
		}
		if resp == nil {
			return []revenuerunmod.RevenueRunRow{}, "", nil
		}
		rows := make([]revenuerunmod.RevenueRunRow, 0, len(resp.GetData()))
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
	rrDeps.ReadRevenueRun = func(fctx context.Context, id string) (*revenuerunmod.RevenueRunWithAttempts, error) {
		runID := id
		resp, err := consumer.ReadRevenueRun(useCases, fctx, &revenuerunpb.ReadRevenueRunRequest{
			Data: &revenuerunpb.RevenueRun{Id: runID},
		})
		if err != nil {
			return nil, err
		}
		if resp == nil || len(resp.GetData()) == 0 {
			return nil, nil
		}
		run := protoRevenueRunToRow(resp.GetData()[0])

		attResp, err := consumer.ListRevenueRunAttempts(useCases, fctx, &revenuerunpb.ListRevenueRunAttemptsRequest{
			RunId: runID,
		})
		if err != nil {
			log.Printf("centymo.Block: failed to load attempts for run %s: %v", id, err)
			attResp = nil
		}
		var attempts []revenuerunmod.RevenueRunAttemptRow
		if attResp != nil {
			attempts = make([]revenuerunmod.RevenueRunAttemptRow, 0, len(attResp.GetData()))
			for _, a := range attResp.GetData() {
				attempts = append(attempts, protoRevenueRunAttemptToRow(a))
			}
		}
		return &revenuerunmod.RevenueRunWithAttempts{Run: run, Attempts: attempts}, nil
	}

	// Wire ListRevenueByRunID — filter revenue list by run_id for the Invoices tab.
	if useCases.Revenue != nil && useCases.Revenue.Revenue != nil &&
		useCases.Revenue.Revenue.GetRevenueListPageData != nil {
		revenueDetailURLPattern := w.revenueRoutes.DetailURL
		rrDeps.ListRevenueByRunID = func(fctx context.Context, runID string) ([]revenuerunmod.RevenueRow, error) {
			resp, err := useCases.Revenue.Revenue.GetRevenueListPageData.Execute(fctx, &revenuepb.GetRevenueListPageDataRequest{
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
			rows := make([]revenuerunmod.RevenueRow, 0, len(resp.GetRevenueList()))
			for _, rv := range resp.GetRevenueList() {
				detailURL := ""
				if revenueDetailURLPattern != "" {
					detailURL = route.ResolveURL(revenueDetailURLPattern, "id", rv.GetId())
				}
				rows = append(rows, revenuerunmod.RevenueRow{
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
	if useCases.Entity != nil && useCases.Entity.Client != nil &&
		useCases.Entity.Client.ListClients != nil {
		lc := useCases.Entity.Client.ListClients.Execute
		rrDeps.ListClients = func(fctx context.Context, cursor string) ([]revenuerunmod.QueueClientRecord, string, error) {
			resp, err := lc(fctx, &clientpb.ListClientsRequest{})
			if err != nil {
				return nil, "", err
			}
			if resp == nil {
				return []revenuerunmod.QueueClientRecord{}, "", nil
			}
			records := make([]revenuerunmod.QueueClientRecord, 0, len(resp.GetData()))
			for _, c := range resp.GetData() {
				records = append(records, revenuerunmod.QueueClientRecord{
					ID:   c.GetId(),
					Name: c.GetName(),
				})
			}
			return records, "", nil
		}
	}

	// ListRevenueRunCandidates — thin shim over consumer.ListRevenueRunCandidates.
	if useCases.Revenue != nil && useCases.Revenue.Revenue != nil &&
		useCases.Revenue.Revenue.ListRevenueRunCandidates != nil {
		rrDeps.ListRevenueRunCandidates = func(fctx context.Context, clientID, asOfDate string) ([]revenuerunmod.QueueCandidateInput, error) {
			candidates, _, err := consumer.ListRevenueRunCandidates(useCases, fctx, consumer.RevenueRunScope{
				ClientID: clientID,
				AsOfDate: asOfDate,
			})
			if err != nil {
				return nil, err
			}
			out := make([]revenuerunmod.QueueCandidateInput, 0, len(candidates))
			for _, c := range candidates {
				out = append(out, revenuerunmod.QueueCandidateInput{
					SubscriptionID: c.SubscriptionID,
					Currency:       c.Currency,
					Amount:         c.Amount,
					Eligible:       c.Eligible,
				})
			}
			return out, nil
		}
	}

	// GenerateRevenueRun — thin shim over consumer.GenerateRevenueRun.
	if useCases.Revenue != nil && useCases.Revenue.Revenue != nil &&
		useCases.Revenue.Revenue.GenerateRevenueRun != nil {
		rrDeps.GenerateRevenueRun = func(fctx context.Context, in revenuerunmod.BatchRunInput) (*revenuerunmod.BatchRunOutput, error) {
			result, err := consumer.GenerateRevenueRun(
				useCases,
				fctx,
				consumer.RevenueRunScope{
					ClientID: in.ClientID,
					AsOfDate: in.AsOfDate,
				},
				consumer.RevenueRunSelections{},
			)
			if err != nil {
				return nil, err
			}
			if result == nil || result.Run == nil {
				return nil, nil
			}
			var created, skipped, errored int
			for _, a := range result.Attempts {
				switch a.GetOutcome() {
				case revenuerunpb.RevenueRunAttemptOutcome_REVENUE_RUN_ATTEMPT_OUTCOME_CREATED:
					created++
				case revenuerunpb.RevenueRunAttemptOutcome_REVENUE_RUN_ATTEMPT_OUTCOME_SKIPPED:
					skipped++
				case revenuerunpb.RevenueRunAttemptOutcome_REVENUE_RUN_ATTEMPT_OUTCOME_ERRORED:
					errored++
				}
			}
			return &revenuerunmod.BatchRunOutput{
				RunID:   result.Run.GetId(),
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
	revenuerunmod.NewModule(rrDeps).RegisterRoutes(ctx.Routes)
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
// to the view-typed revenuerunmod.RevenueRunRow.
// IsStalePending is computed here: status=pending AND initiated_at is older
// than REVENUE_RUN_PENDING_STALE_MINUTES (default 5) minutes ago.
func protoRevenueRunToRow(r *revenuerunpb.RevenueRun) revenuerunmod.RevenueRunRow {
	if r == nil {
		return revenuerunmod.RevenueRunRow{}
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

	return revenuerunmod.RevenueRunRow{
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
// proto message to the view-typed revenuerunmod.RevenueRunAttemptRow.
func protoRevenueRunAttemptToRow(a *revenuerunpb.RevenueRunAttempt) revenuerunmod.RevenueRunAttemptRow {
	if a == nil {
		return revenuerunmod.RevenueRunAttemptRow{}
	}
	return revenuerunmod.RevenueRunAttemptRow{
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
