// Package block — subscription domain wiring.
//
// Holds wireSubscriptionModule (the lifted body of the
// `if cfg.wantSubscription()` branch of Block()) and
// wireCustomizePlanForClient (a helper called exclusively from within it).
//
// Phase 4b of the 20260510-block-go-splitting-strategy.
package block

import (
	"context"
	"fmt"

	"github.com/erniealice/espyna-golang/reference"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	revenuerunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue_run"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"

	"github.com/erniealice/hybra-golang/views/attachment"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	subscriptionaction "github.com/erniealice/centymo-golang/views/subscription/action"
	subscriptiondetail "github.com/erniealice/centymo-golang/views/subscription/detail"
	subscriptionlist "github.com/erniealice/centymo-golang/views/subscription/list"
)

// subscriptionWiring holds everything wireSubscriptionModule needs from the
// surrounding Block() scope. More than 6 fields → struct.
// Kept private; never re-exported.
type subscriptionWiring struct {
	db           centymo.DataSource
	refChecker   reference.Checker
	// Attachment ops
	uploadFile       func(context.Context, string, string, []byte, string) error
	downloadFile     func(context.Context, string, string) ([]byte, error)
	readAttachment   func(context.Context, *attachmentpb.ReadAttachmentRequest) (*attachmentpb.ReadAttachmentResponse, error)
	listAttachments  func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID  func() string
	// Routes + labels
	subscriptionRoutes  centymo.SubscriptionRoutes
	priceScheduleRoutes centymo.PriceScheduleRoutes
	subscriptionLabels  centymo.SubscriptionLabels
	priceScheduleLabels centymo.PriceScheduleLabels
	centymoTableLabels  types.TableLabels
}

// wireSubscriptionModule lifts the body of the `if cfg.wantSubscription()`
// branch from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once at the position where
// the Subscription wiring used to be.
func wireSubscriptionModule(ctx *pyeza.AppContext, cfg *blockConfig, useCases *UseCases, w subscriptionWiring) {
	if !cfg.wantSubscription() {
		return
	}

	subListDeps := &subscriptionlist.ListViewDeps{
		Routes:       w.subscriptionRoutes,
		Labels:       w.subscriptionLabels,
		CommonLabels: ctx.Common,
		TableLabels:  w.centymoTableLabels,
	}
	if useCases.Subscription.GetSubscriptionListPageData != nil {
		subListDeps.GetSubscriptionListPageData = useCases.Subscription.GetSubscriptionListPageData
	}
	if w.refChecker != nil {
		subListDeps.GetInUseIDs = w.refChecker.GetSubscriptionInUseIDs
	}
	ctx.Routes.GET(w.subscriptionRoutes.ListURL, subscriptionlist.NewView(subListDeps))
	// Table-only endpoint — used by sheet.js refreshTable() after
	// activate/deactivate/delete so HTMX swaps the table-card partial,
	// not the whole page.
	if w.subscriptionRoutes.TableURL != "" {
		ctx.Routes.GET(w.subscriptionRoutes.TableURL, subscriptionlist.NewTableView(subListDeps))
	}

	// Subscription CRUD actions
	if useCases.Subscription.CreateSubscription != nil {
		subActionDeps := &subscriptionaction.Deps{
			Routes:             w.subscriptionRoutes,
			Labels:             w.subscriptionLabels,
			CommonLabels:       ctx.Common,
			CreateSubscription: useCases.Subscription.CreateSubscription,
			ReadSubscription:   useCases.Subscription.ReadSubscription,
			UpdateSubscription: useCases.Subscription.UpdateSubscription,
			DeleteSubscription: useCases.Subscription.DeleteSubscription,
			// SetSubscriptionActive uses raw DB update (proto3 omits false booleans)
			SetSubscriptionActive: func(fctx context.Context, id string, active bool) error {
				_, err := w.db.Update(fctx, "subscription", id, map[string]any{"active": active})
				return err
			},
		}
		if w.refChecker != nil {
			subActionDeps.GetInUseIDs = w.refChecker.GetSubscriptionInUseIDs
		}
		if useCases.Subscription.GetSubscriptionItemPageData != nil {
			subActionDeps.GetSubscriptionItemPageData = useCases.Subscription.GetSubscriptionItemPageData
		}
		if useCases.Entity.Client.ListClients != nil {
			subActionDeps.ListClients = useCases.Entity.Client.ListClients
		}
		if useCases.Entity.Client.SearchClientsByName != nil {
			subActionDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName
		}
		if useCases.Plan.ListPlans != nil {
			subActionDeps.ListPlans = useCases.Plan.ListPlans
		}
		if useCases.Plan.ReadPlan != nil {
			subActionDeps.ReadPlan = useCases.Plan.ReadPlan
		}
		if useCases.Plan.SearchPlansByName != nil {
			subActionDeps.SearchPlansByName = useCases.Plan.SearchPlansByName
		}
		if useCases.PricePlan.ListPricePlans != nil {
			subActionDeps.ListPricePlans = useCases.PricePlan.ListPricePlans
		}
		if useCases.PricePlan.ReadPricePlan != nil {
			subActionDeps.ReadPricePlan = useCases.PricePlan.ReadPricePlan
		}
		if useCases.PriceSchedule.ListPriceSchedules != nil {
			subActionDeps.ListPriceSchedules = useCases.PriceSchedule.ListPriceSchedules
		}
		// Wire the espyna recognize-revenue use case so the new
		// drawer + the existing manual-revenue-add auto-populate
		// path share one source of truth.
		if useCases.Revenue.RecognizeRevenueFromSubscription != nil {
			subActionDeps.RecognizeRevenueFromSubscription = useCases.Revenue.RecognizeRevenueFromSubscription
		}

		// 2026-05-06 revenue-run plan Phase 6 (Surface C) — wire the
		// per-subscription Invoice Run drawer callbacks.
		// Both use cases must be present; the drawer gates on nil callbacks.
		if useCases.Revenue.ListRevenueRunCandidates != nil &&
			useCases.Revenue.GenerateRevenueRun != nil {
			listCandidatesUC := useCases.Revenue.ListRevenueRunCandidates
			subActionDeps.ListRevenueRunCandidates = func(fctx context.Context, scope subscriptionaction.RevenueRunScopeAction) ([]subscriptionaction.RevenueRunCandidateAction, string, error) {
				req := &revenuerunpb.ListRevenueRunCandidatesRequest{
					Scope: &revenuerunpb.RevenueRunScope{},
				}
				if scope.Cursor != "" {
					req.Cursor = &scope.Cursor
				}
				if scope.Limit != 0 {
					req.Limit = &scope.Limit
				}
				if scope.WorkspaceID != "" {
					req.Scope.WorkspaceId = &scope.WorkspaceID
				}
				if scope.ClientID != "" {
					req.Scope.ClientId = &scope.ClientID
				}
				if scope.SubscriptionID != "" {
					req.Scope.SubscriptionId = &scope.SubscriptionID
				}
				if scope.AsOfDate != "" {
					req.Scope.AsOfDate = &scope.AsOfDate
				}
				resp, err := listCandidatesUC(fctx, req)
				if err != nil {
					return nil, "", err
				}
				nextCursor := resp.GetNextCursor()
				candidates := resp.GetData()
				out := make([]subscriptionaction.RevenueRunCandidateAction, 0, len(candidates))
				for _, c := range candidates {
					amtDisplay := fmt.Sprintf("%.2f", float64(c.GetAmount())/100)
					out = append(out, subscriptionaction.RevenueRunCandidateAction{
						SubscriptionID:    c.GetSubscriptionId(),
						SubscriptionName:  c.GetSubscriptionName(),
						ClientID:          c.GetClientId(),
						ClientName:        c.GetClientName(),
						PlanName:          c.GetPlanName(),
						BillingCycleLabel: c.GetBillingCycleLabel(),
						Currency:          c.GetCurrency(),
						PeriodStart:       c.GetPeriodStart(),
						PeriodEnd:         c.GetPeriodEnd(),
						PeriodLabel:       c.GetPeriodLabel(),
						PeriodMarker:      c.GetPeriodMarker(),
						Amount:            c.GetAmount(),
						AmountDisplay:     amtDisplay,
						LineItemCount:     int(c.GetLineItemCount()),
						Eligible:          c.GetEligible(),
						BlockerReason:     c.GetBlockerReason(),
					})
				}
				return out, nextCursor, nil
			}
			generateUC := useCases.Revenue.GenerateRevenueRun
			subActionDeps.GenerateRevenueRun = func(fctx context.Context, scope subscriptionaction.RevenueRunScopeAction, sels subscriptionaction.RevenueRunSelectionsAction) (*subscriptionaction.RevenueRunResultAction, error) {
				protoScope := &revenuerunpb.RevenueRunScope{}
				if scope.WorkspaceID != "" {
					protoScope.WorkspaceId = &scope.WorkspaceID
				}
				if scope.ClientID != "" {
					protoScope.ClientId = &scope.ClientID
				}
				if scope.SubscriptionID != "" {
					protoScope.SubscriptionId = &scope.SubscriptionID
				}
				if scope.AsOfDate != "" {
					protoScope.AsOfDate = &scope.AsOfDate
				}
				protoSels := &revenuerunpb.RevenueRunSelections{}
				if sels.FilterToken != "" {
					ft := sels.FilterToken
					protoSels.FilterToken = &ft
				}
				for _, s := range sels.ExplicitList {
					protoSels.ExplicitList = append(protoSels.ExplicitList, &revenuerunpb.SelectedRevenueRunCandidate{
						SubscriptionId: s.SubscriptionID,
						PeriodStart:    s.PeriodStart,
						PeriodEnd:      s.PeriodEnd,
						PeriodMarker:   s.PeriodMarker,
					})
				}
				result, err := generateUC(fctx, &revenuerunpb.GenerateRevenueRunRequest{
					Scope:      protoScope,
					Selections: protoSels,
				})
				if err != nil || result == nil {
					return nil, err
				}
				run := result.GetRun()
				runID := ""
				runStatus := ""
				if run != nil {
					runID = run.GetId()
					runStatus = run.GetStatus().String()
				}
				var created, skipped, errored int32
				for _, a := range result.GetAttempts() {
					switch a.GetOutcome().String() {
					case "REVENUE_RUN_ATTEMPT_OUTCOME_CREATED":
						created++
					case "REVENUE_RUN_ATTEMPT_OUTCOME_SKIPPED":
						skipped++
					default:
						errored++
					}
				}
				return &subscriptionaction.RevenueRunResultAction{
					RunID:   runID,
					Status:  runStatus,
					Created: created,
					Skipped: skipped,
					Errored: errored,
				}, nil
			}
		}

		// 2026-04-27 plan-client-scope plan §4 / §6.5 — wire the
		// CustomizePlanForClient use case via a thin adapter that
		// converts the centymo-side request shape to whatever the
		// espyna-golang use case expects.
		subActionDeps.CustomClientPriceScheduleLabelSuffix =
			w.priceScheduleLabels.Form.CustomClientPriceScheduleLabelSuffix
		wireCustomizePlanForClient(useCases, subActionDeps)

		// 2026-04-29 milestone-billing plan §5 / Phase D — wire the
		// BillingEvent server through to the recognize drawer +
		// mark-ready/waive handlers.
		if useCases.Subscription.ListBillingEventsBySubscription != nil {
			subActionDeps.ListBillingEventsBySubscription = useCases.Subscription.ListBillingEventsBySubscription
			subActionDeps.SetBillingEventStatus = useCases.Subscription.SetBillingEventStatus
		}

		// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / Phase D —
		// wire the JobTemplate read deps that drive the Spawn Jobs
		// section detection on the subscription create drawer + the
		// retroactive spawn drawer. nil-safe.
		if useCases.Operation.JobTemplate.ReadJobTemplate != nil {
			subActionDeps.ReadJobTemplate = useCases.Operation.JobTemplate.ReadJobTemplate
		}
		if useCases.Operation.JobTemplatePhase.ListByJobTemplate != nil {
			subActionDeps.ListJobTemplatePhases = useCases.Operation.JobTemplatePhase.ListByJobTemplate
		}
		if useCases.Operation.JobTemplateTask.ListByPhase != nil {
			subActionDeps.ListJobTemplateTasks = useCases.Operation.JobTemplateTask.ListByPhase
		}
		if useCases.Operation.JobTemplateRelation.ListByParent != nil {
			subActionDeps.ListJobTemplateRelations = useCases.Operation.JobTemplateRelation.ListByParent
		}
		if useCases.Subscription.MaterializeJobsForSubscription != nil {
			materializeUC := useCases.Subscription.MaterializeJobsForSubscription
			subActionDeps.MaterializeJobsForSubscription = func(fctx context.Context, subID string, spawn bool) (int, string, error) {
				resp, err := materializeUC(fctx, &subscriptionpb.MaterializeJobsForSubscriptionRequest{
					SubscriptionId: subID,
					SpawnJobs:      spawn,
				})
				if err != nil {
					return 0, "", err
				}
				if resp == nil {
					return 0, "", nil
				}
				return len(resp.GetSpawnedJobs()), resp.GetSkippedReason(), nil
			}
		}
		// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D —
		// wire espyna's MaterializeInstanceJobsForSubscription through
		// a centymo-side adapter so the Operations tab "Spawn this
		// cycle now" + "Backfill missing cycles" handlers can call it
		// without importing espyna directly. nil-safe: the cycle-spawn
		// and backfill action handlers gate on the adapter pointer.
		if useCases.Subscription.MaterializeInstanceJobsForSubscription != nil {
			materializeInstanceUC := useCases.Subscription.MaterializeInstanceJobsForSubscription
			subActionDeps.MaterializeInstanceJobsForSubscription = func(fctx context.Context, req *subscriptionaction.MaterializeInstanceJobsRequest) (*subscriptionaction.MaterializeInstanceJobsResponse, error) {
				if req == nil {
					return nil, nil
				}
				protoReq := &subscriptionpb.MaterializeInstanceJobsForSubscriptionRequest{
					SubscriptionId: req.SubscriptionID,
					Backfill:       req.Backfill,
				}
				if req.CyclePeriodStart != "" {
					protoReq.CyclePeriodStart = &req.CyclePeriodStart
				}
				resp, err := materializeInstanceUC(fctx, protoReq)
				if err != nil {
					return nil, err
				}
				if resp == nil {
					return &subscriptionaction.MaterializeInstanceJobsResponse{}, nil
				}
				return &subscriptionaction.MaterializeInstanceJobsResponse{
					SpawnedCycleCount:         int(resp.GetSpawnedCycleCount()),
					SpawnedJobCount:           int(resp.GetSpawnedJobCount()),
					OnceAtStartJobCount:       int(resp.GetOnceAtStartJobCount()),
					EngagementWasNewlyCreated: resp.GetEngagementWasNewlyCreated(),
					SkippedReason:             resp.GetSkippedReason(),
					BackfillCappedAt:          resp.GetBackfillCappedAt(),
				}, nil
			}
		}

		ctx.Routes.GET(w.subscriptionRoutes.AddURL, subscriptionaction.NewAddAction(subActionDeps))
		ctx.Routes.POST(w.subscriptionRoutes.AddURL, subscriptionaction.NewAddAction(subActionDeps))
		ctx.Routes.GET(w.subscriptionRoutes.EditURL, subscriptionaction.NewEditAction(subActionDeps))
		ctx.Routes.POST(w.subscriptionRoutes.EditURL, subscriptionaction.NewEditAction(subActionDeps))
		ctx.Routes.POST(w.subscriptionRoutes.DeleteURL, subscriptionaction.NewDeleteAction(subActionDeps))
		ctx.Routes.POST(w.subscriptionRoutes.BulkDeleteURL, subscriptionaction.NewBulkDeleteAction(subActionDeps))
		ctx.Routes.POST(w.subscriptionRoutes.SetStatusURL, subscriptionaction.NewSetStatusAction(subActionDeps))
		ctx.Routes.POST(w.subscriptionRoutes.BulkSetStatusURL, subscriptionaction.NewBulkSetStatusAction(subActionDeps))
		// Recognize-revenue drawer (GET = preview, POST = generate). Per
		// plan §11.1, POST returns HTMXSuccess + refresh-invoices so the
		// invoices table refreshes inline.
		if subActionDeps.RecognizeRevenueFromSubscription != nil && w.subscriptionRoutes.RecognizeURL != "" {
			ctx.Routes.GET(w.subscriptionRoutes.RecognizeURL, subscriptionaction.NewRecognizeAction(subActionDeps))
			ctx.Routes.POST(w.subscriptionRoutes.RecognizeURL, subscriptionaction.NewRecognizeAction(subActionDeps))
		}
		// 2026-05-06 revenue-run Phase 6 (Surface C) — per-subscription
		// Invoice Run drawer (GET = preview candidates, POST = generate run).
		// Gated on both callbacks being wired (set above in the revenue-run
		// wiring block) and the route being configured.
		if subActionDeps.ListRevenueRunCandidates != nil &&
			subActionDeps.GenerateRevenueRun != nil &&
			w.subscriptionRoutes.RevenueRunURL != "" {
			ctx.Routes.GET(w.subscriptionRoutes.RevenueRunURL, subscriptionaction.NewRevenueRunAction(subActionDeps))
			ctx.Routes.POST(w.subscriptionRoutes.RevenueRunURL, subscriptionaction.NewRevenueRunAction(subActionDeps))
		}
		// 2026-04-27 plan-client-scope plan §6.5 — Customize package
		// CTA on subscription detail's Package tab.
		if w.subscriptionRoutes.CustomizePackageURL != "" {
			ctx.Routes.POST(w.subscriptionRoutes.CustomizePackageURL, subscriptionaction.NewCustomizePackageAction(subActionDeps))
		}
		// 2026-04-29 milestone-billing plan §5 / Phase D — mark-ready +
		// waive handlers for BillingEvent rows on the subscription
		// Package tab. Only registered when the BillingEvent server
		// is wired (espyna subscription provider has the adapter).
		if subActionDeps.SetBillingEventStatus != nil {
			if w.subscriptionRoutes.MilestoneMarkReadyURL != "" {
				ctx.Routes.POST(w.subscriptionRoutes.MilestoneMarkReadyURL,
					subscriptionaction.NewMilestoneMarkReadyAction(
						subActionDeps.SetBillingEventStatus,
						subActionDeps.Labels.Errors))
			}
			if w.subscriptionRoutes.MilestoneWaiveURL != "" {
				ctx.Routes.POST(w.subscriptionRoutes.MilestoneWaiveURL,
					subscriptionaction.NewMilestoneWaiveAction(
						subActionDeps.SetBillingEventStatus,
						subActionDeps.Labels.Errors))
			}
		}

		// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / Phase D —
		// HTMX partial that re-renders the Spawn Jobs section on Plan
		// select change + the retroactive spawn drawer.
		if w.subscriptionRoutes.SpawnJobsPartialURL != "" {
			ctx.Routes.GET(w.subscriptionRoutes.SpawnJobsPartialURL, subscriptionaction.NewSpawnJobsPartialAction(subActionDeps))
		}
		if w.subscriptionRoutes.SpawnJobsURL != "" {
			ctx.Routes.GET(w.subscriptionRoutes.SpawnJobsURL, subscriptionaction.NewSpawnJobsAction(subActionDeps))
			ctx.Routes.POST(w.subscriptionRoutes.SpawnJobsURL, subscriptionaction.NewSpawnJobsAction(subActionDeps))
		}
		// 2026-04-30 cyclic-subscription-jobs plan §5.3 — Operations
		// tab CTAs: "Spawn this cycle now" (POST) + "Backfill missing
		// cycles" (GET drawer / POST commit). Both gate on the adapter
		// being wired (handlers also re-check internally). The detail
		// page tab template hides the buttons when the URL fields are
		// empty, so nil-safety is double-bottomed.
		if w.subscriptionRoutes.SpawnCycleJobsURL != "" {
			ctx.Routes.POST(w.subscriptionRoutes.SpawnCycleJobsURL, subscriptionaction.NewSpawnCycleJobsAction(subActionDeps))
		}
		if w.subscriptionRoutes.BackfillCycleJobsURL != "" {
			ctx.Routes.GET(w.subscriptionRoutes.BackfillCycleJobsURL, subscriptionaction.NewBackfillCyclesAction(subActionDeps))
			ctx.Routes.POST(w.subscriptionRoutes.BackfillCycleJobsURL, subscriptionaction.NewBackfillCyclesAction(subActionDeps))
		}
		// 2026-05-01 ad-hoc-subscription-billing — Request Usage CTA on
		// the AD_HOC Operations tab. Pool-Generate-Invoice reuses the
		// existing /action/subscription/recognize-revenue/{id} endpoint
		// (espyna's executeAdHoc dispatches by PricePlan kind).
		if w.subscriptionRoutes.RequestUsageURL != "" {
			ctx.Routes.POST(w.subscriptionRoutes.RequestUsageURL, subscriptionaction.NewRequestUsageAction(subActionDeps))
		}
		// Auto-complete search (http.HandlerFunc — uses HandleFunc, not GET)
		handleFunc(ctx.Routes, "GET", w.subscriptionRoutes.SearchClientURL, subscriptionaction.NewSearchClientsAction(subActionDeps))
		handleFunc(ctx.Routes, "GET", w.subscriptionRoutes.SearchPlanURL, subscriptionaction.NewSearchPlansAction(subActionDeps))
	}

	// Subscription detail page + tab action
	if useCases.Subscription.ReadSubscription != nil {
		subDetailDeps := &subscriptiondetail.DetailViewDeps{
			Routes:           w.subscriptionRoutes,
			ReadSubscription: useCases.Subscription.ReadSubscription,
			Labels:           w.subscriptionLabels,
			CommonLabels:     ctx.Common,
			TableLabels:      w.centymoTableLabels,
			AttachmentOps: attachment.AttachmentOps{
				UploadFile:       w.uploadFile,
				DownloadFile:     w.downloadFile,
				ListAttachments:  w.listAttachments,
				CreateAttachment: w.createAttachment,
				ReadAttachment:   w.readAttachment,
				DeleteAttachment: w.deleteAttachment,
				NewAttachmentID:  w.newAttachmentID,
			},
		}
		if useCases.Subscription.GetSubscriptionItemPageData != nil {
			subDetailDeps.GetSubscriptionItemPageData = useCases.Subscription.GetSubscriptionItemPageData
		}
		if useCases.Entity.Client.ReadClient != nil {
			subDetailDeps.ReadClient = useCases.Entity.Client.ReadClient
		}
		if useCases.Revenue.GetListPageData != nil {
			subDetailDeps.GetRevenueListPageData = useCases.Revenue.GetListPageData
		}
		// 2026-04-29 milestone-billing — wire BillingEvent listing into
		// the subscription detail Package tab.
		if useCases.Subscription.ListBillingEventsBySubscription != nil {
			subDetailDeps.ListBillingEventsBySubscription = useCases.Subscription.ListBillingEventsBySubscription
		}
		// 2026-04-29 auto-spawn-jobs-from-subscription Phase D — wire
		// the Operations tab data ops + spawn-jobs CTA URL.
		if useCases.Operation.Job.GetJobsByOrigin != nil {
			subDetailDeps.GetJobsByOrigin = useCases.Operation.Job.GetJobsByOrigin
		}
		if useCases.Operation.JobPhase.ListByJob != nil {
			subDetailDeps.ListJobPhasesByJob = useCases.Operation.JobPhase.ListByJob
		}
		subDetailDeps.SpawnJobsURL = w.subscriptionRoutes.SpawnJobsURL
		subDetailDeps.JobDetailURL = cfg.jobDetailURL
		subDetailDeps.ClientDetailURL = cfg.clientDetailURL
		// 2026-05-04 — engagement breadcrumb (rate-card → plan).
		subDetailDeps.PriceScheduleDetailURL = w.priceScheduleRoutes.DetailURL
		subDetailDeps.PricePlanDetailURL = w.priceScheduleRoutes.PlanDetailURL
		if useCases.PriceSchedule.ReadPriceSchedule != nil {
			subDetailDeps.ReadPriceSchedule = useCases.PriceSchedule.ReadPriceSchedule
		}
		if useCases.PricePlan.ReadPricePlan != nil {
			subDetailDeps.ReadPricePlan = useCases.PricePlan.ReadPricePlan
		}
		ctx.Routes.GET(w.subscriptionRoutes.DetailURL, subscriptiondetail.NewView(subDetailDeps))
		ctx.Routes.GET(w.subscriptionRoutes.TabActionURL, subscriptiondetail.NewTabAction(subDetailDeps))
		// Nested route — same view, breadcrumb activated via path param.
		if w.subscriptionRoutes.UnderClientDetailURL != "" {
			ctx.Routes.GET(w.subscriptionRoutes.UnderClientDetailURL, subscriptiondetail.NewView(subDetailDeps))
		}
		// 2026-05-04 — Engagement detail nested under the rate-card → plan
		// path. Same view; the URL alone activates the schedule + plan
		// breadcrumb segments inside the subscription detail page.
		if w.priceScheduleRoutes.PlanEngagementDetailURL != "" {
			ctx.Routes.GET(w.priceScheduleRoutes.PlanEngagementDetailURL, subscriptiondetail.NewView(subDetailDeps))
		}
		// Subscription attachments
		if w.uploadFile != nil {
			ctx.Routes.GET(w.subscriptionRoutes.AttachmentUploadURL, subscriptiondetail.NewAttachmentUploadAction(subDetailDeps))
			ctx.Routes.POST(w.subscriptionRoutes.AttachmentUploadURL, subscriptiondetail.NewAttachmentUploadAction(subDetailDeps))
			ctx.Routes.POST(w.subscriptionRoutes.AttachmentDeleteURL, subscriptiondetail.NewAttachmentDeleteAction(subDetailDeps))
			if w.downloadFile != nil && w.readAttachment != nil && w.subscriptionRoutes.AttachmentDownloadURL != "" {
				handleFunc(ctx.Routes, "GET", w.subscriptionRoutes.AttachmentDownloadURL, subscriptiondetail.NewAttachmentDownloadHandler(subDetailDeps))
			}
		}
	}
}

// wireCustomizePlanForClient threads the Plan.CustomizePlanForClient use
// case into the centymo subscription action Deps. The centymo side uses an
// in-package type-narrow shape so its handlers don't depend directly on the
// espyna generated proto/use-case structs.
//
// When the use case isn't wired (composition layer didn't initialize it), we
// leave the function pointer nil; the handler falls through to a generic
// `customize_failed` toast.
//
// 2026-04-27 plan-client-scope plan §4. Same adapter pattern as
// RecognizeRevenueFromSubscription above.
func wireCustomizePlanForClient(useCases *UseCases, subActionDeps *subscriptionaction.Deps) {
	if useCases.Plan.CustomizePlanForClient == nil {
		return
	}
	customizeUC := useCases.Plan.CustomizePlanForClient
	subActionDeps.CustomizePlanForClient = func(
		ctx context.Context, req *subscriptionaction.CustomizePlanForClientRequest,
	) (*subscriptionaction.CustomizePlanForClientResponse, error) {
		protoReq := &planpb.CustomizePlanForClientRequest{
			SourcePlanId:      req.SourcePlanID,
			SourcePricePlanId: req.SourcePricePlanID,
			ClientId:          req.ClientID,
		}
		if req.SubscriptionID != "" {
			protoReq.SubscriptionId = &req.SubscriptionID
		}
		if req.NewScheduleName != "" {
			protoReq.NewScheduleName = &req.NewScheduleName
		}
		resp, err := customizeUC(ctx, protoReq)
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return &subscriptionaction.CustomizePlanForClientResponse{}, nil
		}
		return &subscriptionaction.CustomizePlanForClientResponse{
			NewPlanID:      resp.GetNewPlanId(),
			NewPricePlanID: resp.GetNewPricePlanId(),
			NewScheduleID:  resp.GetNewPriceScheduleId(),
			Reused:         resp.GetReused(),
		}, nil
	}
}
