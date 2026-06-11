// Package block — expense-recognition-run domain wiring (Surfaces B + D).
//
// Holds wireExpenseRecognitionRunModule (the lifted body of the
// `if cfg.wantExpenseRecognitionRun()` branch of Block()) plus the proto↔view
// translators that only the expense-recognition-run wiring calls.
//
// Mirror of packages/centymo-golang/block/revenue_run.go.
// Plan A 20260517-expense-run Phase 4.
package block

import (
	"context"
	expenserunmodmodule "github.com/erniealice/centymo-golang/domain/expenditure/expense_recognition_run/module"
	"log"
	"os"
	"strconv"
	"time"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
	expenserecognitionrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition_run"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"

	expendituredomain "github.com/erniealice/centymo-golang/domain/expenditure"
)

// expenseRecognitionRunWiring holds everything wireExpenseRecognitionRunModule
// needs from the surrounding Block() scope.
type expenseRecognitionRunWiring struct {
	routes                                 expendituredomain.ExpenseRecognitionRunRoutes
	labels                                 expendituredomain.ExpenseRecognitionRunLabels
	expenditureRoutes                      expendituredomain.ExpenditureRoutes
	centymoTableLabels                     types.TableLabels
	supplierDetailURL                      string
	supplierExpenseRecognitionRunDrawerURL string
}

// wireExpenseRecognitionRunModule wires Surfaces B + D of the buying-side
// Expense Recognition Run. block.go calls this exactly once when
// cfg.wantExpenseRecognitionRun() is true and all required deps are set.
func wireExpenseRecognitionRunModule(ctx *pyeza.AppContext, cfg *blockConfig, useCases *UseCases, w expenseRecognitionRunWiring) {
	deps := &expenserunmodmodule.ModuleDeps{
		Routes:                    w.routes,
		Labels:                    w.labels,
		CommonLabels:              ctx.Common,
		TableLabels:               w.centymoTableLabels,
		SupplierDetailURLTemplate: w.supplierDetailURL,
		SupplierDrawerURLTemplate: w.supplierExpenseRecognitionRunDrawerURL,
	}

	// Surface D — ListExpenseRecognitionRuns translator.
	if useCases.ExpenseRecognitionRun.ListExpenseRecognitionRuns != nil {
		deps.ListExpenseRecognitionRuns = func(fctx context.Context, scope expenserunmodmodule.ListExpenseRecognitionRunsScope) ([]expenserunmodmodule.ExpenseRecognitionRunRow, string, error) {
			req := &expenserecognitionrunpb.ListExpenseRecognitionRunsRequest{}
			resp, err := useCases.ExpenseRecognitionRun.ListExpenseRecognitionRuns(fctx, req)
			if err != nil {
				return nil, "", err
			}
			if resp == nil {
				return []expenserunmodmodule.ExpenseRecognitionRunRow{}, "", nil
			}
			rows := make([]expenserunmodmodule.ExpenseRecognitionRunRow, 0, len(resp.GetData()))
			for _, r := range resp.GetData() {
				row := protoExpenseRunToRow(r)
				if scope.Status != "" && row.Status != scope.Status {
					continue
				}
				rows = append(rows, row)
			}
			return rows, "", nil
		}
	}

	// Surface D — ReadExpenseRecognitionRun translator.
	if useCases.ExpenseRecognitionRun.ReadExpenseRecognitionRun != nil {
		deps.ReadExpenseRecognitionRun = func(fctx context.Context, id string) (*expenserunmodmodule.ExpenseRecognitionRunWithAttempts, error) {
			resp, err := useCases.ExpenseRecognitionRun.ReadExpenseRecognitionRun(fctx, &expenserecognitionrunpb.ReadExpenseRecognitionRunRequest{
				Data: &expenserecognitionrunpb.ExpenseRecognitionRun{Id: id},
			})
			if err != nil {
				return nil, err
			}
			if resp == nil || len(resp.GetData()) == 0 {
				return nil, nil
			}
			run := protoExpenseRunToRow(resp.GetData()[0])

			var attempts []expenserunmodmodule.ExpenseRecognitionRunAttemptRow
			if useCases.ExpenseRecognitionRun.ListExpenseRecognitionRunAttempts != nil {
				attResp, err := useCases.ExpenseRecognitionRun.ListExpenseRecognitionRunAttempts(fctx, &expenserecognitionrunpb.ListExpenseRecognitionRunAttemptsRequest{
					RunId: id,
				})
				if err != nil {
					log.Printf("centymo.Block: failed to load attempts for expense run %s: %v", id, err)
				} else if attResp != nil {
					attempts = make([]expenserunmodmodule.ExpenseRecognitionRunAttemptRow, 0, len(attResp.GetData()))
					for _, a := range attResp.GetData() {
						attempts = append(attempts, protoExpenseRunAttemptToRow(a))
					}
				}
			}
			return &expenserunmodmodule.ExpenseRecognitionRunWithAttempts{Run: run, Attempts: attempts}, nil
		}
	}

	// Surface D — ListExpendituresByRunID translator (Bills tab).
	if useCases.Expenditure.ListExpenditures != nil {
		expenditureDetailURLPattern := w.expenditureRoutes.DetailURL
		deps.ListExpendituresByRunID = func(fctx context.Context, runID string) ([]expenserunmodmodule.ExpenditureRow, error) {
			resp, err := useCases.Expenditure.ListExpenditures(fctx, &expenditurepb.ListExpendituresRequest{
				Filters: &commonpb.FilterRequest{
					Filters: []*commonpb.TypedFilter{
						{
							Field: "run_id",
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
			rows := make([]expenserunmodmodule.ExpenditureRow, 0)
			if resp != nil {
				for _, e := range resp.GetData() {
					detailURL := ""
					if expenditureDetailURLPattern != "" {
						detailURL = route.ResolveURL(expenditureDetailURLPattern, "id", e.GetId())
					}
					rows = append(rows, expenserunmodmodule.ExpenditureRow{
						ID:              e.GetId(),
						ReferenceNumber: e.GetReferenceNumber(),
						ExpenditureDate: expenditureDateString(e),
						TotalAmount:     int64(e.GetTotalAmount()),
						Currency:        e.GetCurrency(),
						Status:          e.GetStatus(),
						DetailURL:       detailURL,
					})
				}
			}
			return rows, nil
		}
	}

	// Surface D — ListExpenseRecognitionsByRunID translator (Recognitions tab).
	if useCases.Expenditure.ListExpenseRecognitions != nil {
		deps.ListExpenseRecognitionsByRunID = func(fctx context.Context, runID string) ([]expenserunmodmodule.ExpenseRecognitionRow, error) {
			resp, err := useCases.Expenditure.ListExpenseRecognitions(fctx, &expenserecognitionpb.ListExpenseRecognitionsRequest{
				Filters: &commonpb.FilterRequest{
					Filters: []*commonpb.TypedFilter{
						{
							Field: "run_id",
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
			rows := make([]expenserunmodmodule.ExpenseRecognitionRow, 0)
			if resp != nil {
				for _, r := range resp.GetData() {
					sourceKind := ""
					if r.GetSupplierSubscriptionId() != "" {
						sourceKind = "subscription"
					} else if r.GetAdvanceDisbursementId() != "" {
						sourceKind = "advance_disbursement"
					}
					recDate := ""
					if r.RecognitionDate != nil {
						recDate = r.GetRecognitionDate().AsTime().Format("2006-01-02")
					}
					rows = append(rows, expenserunmodmodule.ExpenseRecognitionRow{
						ID:              r.GetId(),
						ReferenceNumber: r.GetName(),
						RecognitionDate: recDate,
						TotalAmount:     r.GetTotalAmount(),
						Currency:        r.GetCurrency(),
						Status:          r.GetStatus().String(),
						SourceKind:      sourceKind,
					})
				}
			}
			return rows, nil
		}
	}

	// Surface B — ListSuppliers translator.
	if useCases.Entity.Supplier.ListSuppliers != nil {
		deps.ListSuppliers = func(fctx context.Context, cursor string) ([]expenserunmodmodule.QueueSupplierRecord, string, error) {
			resp, err := useCases.Entity.Supplier.ListSuppliers(fctx, &supplierpb.ListSuppliersRequest{})
			if err != nil {
				return nil, "", err
			}
			if resp == nil {
				return []expenserunmodmodule.QueueSupplierRecord{}, "", nil
			}
			records := make([]expenserunmodmodule.QueueSupplierRecord, 0, len(resp.GetData()))
			for _, s := range resp.GetData() {
				records = append(records, expenserunmodmodule.QueueSupplierRecord{
					ID:   s.GetId(),
					Name: s.GetName(),
				})
			}
			return records, "", nil
		}
	}

	// Surface B — ListExpenseRunCandidates translator.
	if useCases.ExpenseRecognitionRun.ListExpenseRunCandidates != nil {
		deps.ListExpenseRunCandidates = func(fctx context.Context, supplierID, asOfDate string) ([]expenserunmodmodule.QueueCandidateInput, error) {
			supplierIDCopy := supplierID
			asOfDateCopy := asOfDate
			resp, err := useCases.ExpenseRecognitionRun.ListExpenseRunCandidates(fctx, &expenserecognitionrunpb.ListExpenseRunCandidatesRequest{
				Scope: &expenserecognitionrunpb.ExpenseRecognitionRunScopeMsg{
					SupplierId: &supplierIDCopy,
					AsOfDate:   &asOfDateCopy,
				},
			})
			if err != nil {
				return nil, err
			}
			if resp == nil {
				return nil, nil
			}
			out := make([]expenserunmodmodule.QueueCandidateInput, 0, len(resp.GetData()))
			for _, c := range resp.GetData() {
				sourceKind := ""
				switch c.GetSourceKind() {
				case expenserecognitionrunpb.ExpenseRecognitionRunSourceKind_EXPENSE_RECOGNITION_RUN_SOURCE_KIND_SUBSCRIPTION_CYCLE:
					sourceKind = "subscription"
				case expenserecognitionrunpb.ExpenseRecognitionRunSourceKind_EXPENSE_RECOGNITION_RUN_SOURCE_KIND_ADVANCE_DISBURSEMENT:
					sourceKind = "advance_disbursement"
				}
				out = append(out, expenserunmodmodule.QueueCandidateInput{
					SourceKind:             sourceKind,
					SupplierSubscriptionID: c.GetSupplierSubscriptionId(),
					AdvanceDisbursementID:  c.GetAdvanceDisbursementId(),
					Currency:               c.GetCurrency(),
					Amount:                 c.GetAmount(),
					Eligible:               c.GetEligible(),
				})
			}
			return out, nil
		}
	}

	// Surface A + C — drawer candidate listing translator.
	if useCases.ExpenseRecognitionRun.ListExpenseRunCandidates != nil {
		deps.ListExpenseRunCandidatesForDrawer = func(fctx context.Context, scope expenserunmodmodule.DrawerListScope) ([]expenserunmodmodule.DrawerCandidateRow, error) {
			scopeMsg := &expenserecognitionrunpb.ExpenseRecognitionRunScopeMsg{
				AsOfDate: &scope.AsOfDate,
			}
			if scope.SupplierID != "" {
				supplierID := scope.SupplierID
				scopeMsg.SupplierId = &supplierID
			}
			if scope.SupplierSubscriptionID != "" {
				subID := scope.SupplierSubscriptionID
				scopeMsg.SupplierSubscriptionId = &subID
			}
			resp, err := useCases.ExpenseRecognitionRun.ListExpenseRunCandidates(fctx, &expenserecognitionrunpb.ListExpenseRunCandidatesRequest{
				Scope: scopeMsg,
			})
			if err != nil {
				return nil, err
			}
			if resp == nil {
				return nil, nil
			}
			out := make([]expenserunmodmodule.DrawerCandidateRow, 0, len(resp.GetData()))
			for _, c := range resp.GetData() {
				out = append(out, expenserunmodmodule.DrawerCandidateRow{
					SourceKind:             expenseRunSourceKindString(c.GetSourceKind()),
					SourceLabel:            c.GetSourceLabel(),
					SupplierSubscriptionID: c.GetSupplierSubscriptionId(),
					AdvanceDisbursementID:  c.GetAdvanceDisbursementId(),
					PeriodStart:            c.GetPeriodStart(),
					PeriodEnd:              c.GetPeriodEnd(),
					PeriodLabel:            c.GetPeriodLabel(),
					PeriodMarker:           c.GetPeriodMarker(),
					Amount:                 c.GetAmount(),
					Currency:               c.GetCurrency(),
					Eligible:               c.GetEligible(),
					BlockerReason:          c.GetBlockerReason(),
					SuppressingAdvanceID:   c.GetSuppressingAdvanceDisbursementId(),
				})
			}
			return out, nil
		}
	}

	// Surface A + C — drawer Generate translator.
	if useCases.ExpenseRecognitionRun.GenerateExpenseRun != nil {
		deps.GenerateExpenseRunForDrawer = func(fctx context.Context, in expenserunmodmodule.DrawerGenerateInput) (*expenserunmodmodule.BatchRunOutput, error) {
			scopeMsg := &expenserecognitionrunpb.ExpenseRecognitionRunScopeMsg{
				AsOfDate: &in.AsOfDate,
			}
			if in.SupplierID != "" {
				supplierID := in.SupplierID
				scopeMsg.SupplierId = &supplierID
			}
			if in.SupplierSubscriptionID != "" {
				subID := in.SupplierSubscriptionID
				scopeMsg.SupplierSubscriptionId = &subID
			}
			resp, err := useCases.ExpenseRecognitionRun.GenerateExpenseRun(fctx, &expenserecognitionrunpb.GenerateExpenseRunRequest{
				Scope: scopeMsg,
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
				case expenserecognitionrunpb.ExpenseRecognitionRunAttemptOutcome_EXPENSE_RECOGNITION_RUN_ATTEMPT_OUTCOME_CREATED:
					created++
				case expenserecognitionrunpb.ExpenseRecognitionRunAttemptOutcome_EXPENSE_RECOGNITION_RUN_ATTEMPT_OUTCOME_SKIPPED:
					skipped++
				case expenserecognitionrunpb.ExpenseRecognitionRunAttemptOutcome_EXPENSE_RECOGNITION_RUN_ATTEMPT_OUTCOME_ERRORED:
					errored++
				}
			}
			return &expenserunmodmodule.BatchRunOutput{
				RunID:   resp.GetRun().GetId(),
				Created: created,
				Skipped: skipped,
				Errored: errored,
			}, nil
		}
	}

	// Surface B — GenerateExpenseRun translator.
	if useCases.ExpenseRecognitionRun.GenerateExpenseRun != nil {
		deps.GenerateExpenseRun = func(fctx context.Context, in expenserunmodmodule.BatchRunInput) (*expenserunmodmodule.BatchRunOutput, error) {
			supplierID := in.SupplierID
			asOfDate := in.AsOfDate
			resp, err := useCases.ExpenseRecognitionRun.GenerateExpenseRun(fctx, &expenserecognitionrunpb.GenerateExpenseRunRequest{
				Scope: &expenserecognitionrunpb.ExpenseRecognitionRunScopeMsg{
					SupplierId: &supplierID,
					AsOfDate:   &asOfDate,
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
				case expenserecognitionrunpb.ExpenseRecognitionRunAttemptOutcome_EXPENSE_RECOGNITION_RUN_ATTEMPT_OUTCOME_CREATED:
					created++
				case expenserecognitionrunpb.ExpenseRecognitionRunAttemptOutcome_EXPENSE_RECOGNITION_RUN_ATTEMPT_OUTCOME_SKIPPED:
					skipped++
				case expenserecognitionrunpb.ExpenseRecognitionRunAttemptOutcome_EXPENSE_RECOGNITION_RUN_ATTEMPT_OUTCOME_ERRORED:
					errored++
				}
			}
			return &expenserunmodmodule.BatchRunOutput{
				RunID:   resp.GetRun().GetId(),
				Created: created,
				Skipped: skipped,
				Errored: errored,
			}, nil
		}
	}

	expenserunmodmodule.NewModule(deps).RegisterRoutes(ctx.Routes)
}

// ---------------------------------------------------------------------------
// Phase 4 — Expense Run proto → view-type shim helpers
// ---------------------------------------------------------------------------

// expenseRunMillisToRFC3339 converts a proto epoch-millisecond int64 to
// RFC3339 UTC. Returns "" when ms is zero.
func expenseRunMillisToRFC3339(ms int64) string {
	if ms == 0 {
		return ""
	}
	return time.UnixMilli(ms).UTC().Format(time.RFC3339)
}

// expenseRunScopeKindString maps the proto Scope enum to a short lowercase string.
func expenseRunScopeKindString(sk expenserecognitionrunpb.ExpenseRecognitionRunScope) string {
	switch sk {
	case expenserecognitionrunpb.ExpenseRecognitionRunScope_EXPENSE_RECOGNITION_RUN_SCOPE_SUPPLIER:
		return "supplier"
	case expenserecognitionrunpb.ExpenseRecognitionRunScope_EXPENSE_RECOGNITION_RUN_SCOPE_SUBSCRIPTION:
		return "subscription"
	case expenserecognitionrunpb.ExpenseRecognitionRunScope_EXPENSE_RECOGNITION_RUN_SCOPE_WORKSPACE:
		return "workspace"
	default:
		return ""
	}
}

// expenseRunStatusString maps the proto Status enum to the lowercase string.
func expenseRunStatusString(s expenserecognitionrunpb.ExpenseRecognitionRunStatus) string {
	switch s {
	case expenserecognitionrunpb.ExpenseRecognitionRunStatus_EXPENSE_RECOGNITION_RUN_STATUS_PENDING:
		return "pending"
	case expenserecognitionrunpb.ExpenseRecognitionRunStatus_EXPENSE_RECOGNITION_RUN_STATUS_COMPLETE:
		return "complete"
	case expenserecognitionrunpb.ExpenseRecognitionRunStatus_EXPENSE_RECOGNITION_RUN_STATUS_FAILED:
		return "failed"
	default:
		return ""
	}
}

// expenseRunAttemptOutcomeString maps the proto Outcome enum to lowercase.
func expenseRunAttemptOutcomeString(o expenserecognitionrunpb.ExpenseRecognitionRunAttemptOutcome) string {
	switch o {
	case expenserecognitionrunpb.ExpenseRecognitionRunAttemptOutcome_EXPENSE_RECOGNITION_RUN_ATTEMPT_OUTCOME_CREATED:
		return "created"
	case expenserecognitionrunpb.ExpenseRecognitionRunAttemptOutcome_EXPENSE_RECOGNITION_RUN_ATTEMPT_OUTCOME_SKIPPED:
		return "skipped"
	case expenserecognitionrunpb.ExpenseRecognitionRunAttemptOutcome_EXPENSE_RECOGNITION_RUN_ATTEMPT_OUTCOME_ERRORED:
		return "errored"
	default:
		return ""
	}
}

// expenseRunSourceKindString maps the proto SourceKind enum to lowercase.
func expenseRunSourceKindString(sk expenserecognitionrunpb.ExpenseRecognitionRunSourceKind) string {
	switch sk {
	case expenserecognitionrunpb.ExpenseRecognitionRunSourceKind_EXPENSE_RECOGNITION_RUN_SOURCE_KIND_SUBSCRIPTION_CYCLE:
		return "subscription"
	case expenserecognitionrunpb.ExpenseRecognitionRunSourceKind_EXPENSE_RECOGNITION_RUN_SOURCE_KIND_ADVANCE_DISBURSEMENT:
		return "advance_disbursement"
	default:
		return ""
	}
}

// protoExpenseRunToRow translates a *expenserecognitionrunpb.ExpenseRecognitionRun
// proto message to the view-typed expenserunmodmodule.ExpenseRecognitionRunRow.
// IsStalePending is computed using EXPENSE_RUN_PENDING_STALE_MINUTES (default 5).
func protoExpenseRunToRow(r *expenserecognitionrunpb.ExpenseRecognitionRun) expenserunmodmodule.ExpenseRecognitionRunRow {
	if r == nil {
		return expenserunmodmodule.ExpenseRecognitionRunRow{}
	}
	initiatedAt := expenseRunMillisToRFC3339(r.GetInitiatedAt())
	completedAt := expenseRunMillisToRFC3339(r.GetCompletedAt())
	status := expenseRunStatusString(r.GetStatus())

	isPending := r.GetStatus() == expenserecognitionrunpb.ExpenseRecognitionRunStatus_EXPENSE_RECOGNITION_RUN_STATUS_PENDING
	isStalePending := false
	if isPending && r.GetInitiatedAt() > 0 {
		age := time.Since(time.UnixMilli(r.GetInitiatedAt()))
		isStalePending = age > expenseRunStaleThreshold()
	}

	return expenserunmodmodule.ExpenseRecognitionRunRow{
		ID:                     r.GetId(),
		ScopeKind:              expenseRunScopeKindString(r.GetScope()),
		SupplierID:             r.GetSupplierId(),
		SupplierSubscriptionID: r.GetSupplierSubscriptionId(),
		AsOfDate:               r.GetAsOfDate(),
		Initiator:              r.GetInitiatedBy(),
		InitiatedAt:            initiatedAt,
		CompletedAt:            completedAt,
		Status:                 status,
		SelectionCount:         r.GetSelectionCount(),
		CreatedCount:           r.GetCreatedCount(),
		SkippedCount:           r.GetSkippedCount(),
		ErroredCount:           r.GetErroredCount(),
		IsStalePending:         isStalePending,
		Notes:                  r.GetNotes(),
	}
}

// protoExpenseRunAttemptToRow translates a proto attempt to the view-typed row.
func protoExpenseRunAttemptToRow(a *expenserecognitionrunpb.ExpenseRecognitionRunAttempt) expenserunmodmodule.ExpenseRecognitionRunAttemptRow {
	if a == nil {
		return expenserunmodmodule.ExpenseRecognitionRunAttemptRow{}
	}
	return expenserunmodmodule.ExpenseRecognitionRunAttemptRow{
		ID:                     a.GetId(),
		RunID:                  a.GetRunId(),
		SourceKind:             expenseRunSourceKindString(a.GetSourceKind()),
		SupplierSubscriptionID: a.GetSupplierSubscriptionId(),
		AdvanceDisbursementID:  a.GetAdvanceDisbursementId(),
		PeriodStart:            a.GetPeriodStart(),
		PeriodEnd:              a.GetPeriodEnd(),
		PeriodMarker:           a.GetPeriodMarker(),
		AttemptedAt:            expenseRunMillisToRFC3339(a.GetAttemptedAt()),
		Outcome:                expenseRunAttemptOutcomeString(a.GetOutcome()),
		ExpenseRecognitionID:   a.GetExpenseRecognitionId(),
		ExpenditureID:          a.GetExpenditureId(),
		ErrorCode:              a.GetErrorCode(),
		ErrorMessage:           a.GetErrorMessage(),
	}
}

func expenseRunStaleThreshold() time.Duration {
	if v := os.Getenv("EXPENSE_RUN_PENDING_STALE_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Minute
		}
	}
	return 5 * time.Minute
}

// expenditureDateString resolves a YYYY-MM-DD date string from a proto
// Expenditure record. ExpenditureDate is an epoch-millisecond int64; returns
// "" when unset (zero).
func expenditureDateString(e *expenditurepb.Expenditure) string {
	if e == nil {
		return ""
	}
	ms := e.GetExpenditureDate()
	if ms == 0 {
		return ""
	}
	return time.UnixMilli(ms).UTC().Format("2006-01-02")
}
