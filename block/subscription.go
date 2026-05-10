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

	consumer "github.com/erniealice/espyna-golang/consumer"
	"github.com/erniealice/espyna-golang/reference"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"

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
func wireSubscriptionModule(ctx *pyeza.AppContext, cfg *blockConfig, useCases *consumer.UseCases, w subscriptionWiring) {
	if !cfg.wantSubscription() {
		return
	}

	subListDeps := &subscriptionlist.ListViewDeps{
		Routes:       w.subscriptionRoutes,
		Labels:       w.subscriptionLabels,
		CommonLabels: ctx.Common,
		TableLabels:  w.centymoTableLabels,
	}
	if useCases.Subscription != nil && useCases.Subscription.Subscription != nil {
		subListDeps.GetSubscriptionListPageData = useCases.Subscription.Subscription.GetSubscriptionListPageData.Execute
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
	if useCases.Subscription != nil && useCases.Subscription.Subscription != nil && useCases.Subscription.Subscription.CreateSubscription != nil {
		subActionDeps := &subscriptionaction.Deps{
			Routes:             w.subscriptionRoutes,
			Labels:             w.subscriptionLabels,
			CommonLabels:       ctx.Common,
			CreateSubscription: useCases.Subscription.Subscription.CreateSubscription.Execute,
			ReadSubscription:   useCases.Subscription.Subscription.ReadSubscription.Execute,
			UpdateSubscription: useCases.Subscription.Subscription.UpdateSubscription.Execute,
			DeleteSubscription: useCases.Subscription.Subscription.DeleteSubscription.Execute,
			// SetSubscriptionActive uses raw DB update (proto3 omits false booleans)
			SetSubscriptionActive: func(fctx context.Context, id string, active bool) error {
				_, err := w.db.Update(fctx, "subscription", id, map[string]any{"active": active})
				return err
			},
		}
		if w.refChecker != nil {
			subActionDeps.GetInUseIDs = w.refChecker.GetSubscriptionInUseIDs
		}
		if useCases.Subscription.Subscription.GetSubscriptionItemPageData != nil {
			subActionDeps.GetSubscriptionItemPageData = useCases.Subscription.Subscription.GetSubscriptionItemPageData.Execute
		}
		if useCases.Entity != nil && useCases.Entity.Client != nil {
			subActionDeps.ListClients = useCases.Entity.Client.ListClients.Execute
			if useCases.Entity.Client.SearchClientsByName != nil {
				subActionDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName.Execute
			}
		}
		if useCases.Subscription.Plan != nil {
			subActionDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
			if useCases.Subscription.Plan.ReadPlan != nil {
				subActionDeps.ReadPlan = useCases.Subscription.Plan.ReadPlan.Execute
			}
			if useCases.Subscription.Plan.SearchPlansByName != nil {
				subActionDeps.SearchPlansByName = useCases.Subscription.Plan.SearchPlansByName.Execute
			}
		}
		if useCases.Subscription.PricePlan != nil {
			subActionDeps.ListPricePlans = useCases.Subscription.PricePlan.ListPricePlans.Execute
			if useCases.Subscription.PricePlan.ReadPricePlan != nil {
				subActionDeps.ReadPricePlan = useCases.Subscription.PricePlan.ReadPricePlan.Execute
			}
		}
		if useCases.Subscription.PriceSchedule != nil && useCases.Subscription.PriceSchedule.ListPriceSchedules != nil {
			subActionDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
		}
		// Wire the espyna recognize-revenue use case so the new
		// drawer + the existing manual-revenue-add auto-populate
		// path share one source of truth.
		if useCases.Revenue != nil && useCases.Revenue.Revenue != nil &&
			useCases.Revenue.Revenue.RecognizeRevenueFromSubscription != nil {
			subActionDeps.RecognizeRevenueFromSubscription =
				useCases.Revenue.Revenue.RecognizeRevenueFromSubscription.Execute
		}

		// 2026-05-06 revenue-run plan Phase 6 (Surface C) — wire the
		// per-subscription Invoice Run drawer callbacks.
		// Both use cases must be present; the drawer gates on nil callbacks.
		if useCases.Revenue != nil && useCases.Revenue.Revenue != nil &&
			useCases.Revenue.Revenue.ListRevenueRunCandidates != nil &&
			useCases.Revenue.Revenue.GenerateRevenueRun != nil {
			subActionDeps.ListRevenueRunCandidates = func(fctx context.Context, scope subscriptionaction.RevenueRunScopeAction) ([]subscriptionaction.RevenueRunCandidateAction, string, error) {
				candidates, nextCursor, err := consumer.ListRevenueRunCandidates(useCases, fctx, consumer.RevenueRunScope{
					WorkspaceID:    scope.WorkspaceID,
					ClientID:       scope.ClientID,
					SubscriptionID: scope.SubscriptionID,
					AsOfDate:       scope.AsOfDate,
					Cursor:         scope.Cursor,
					Limit:          scope.Limit,
				})
				if err != nil {
					return nil, "", err
				}
				out := make([]subscriptionaction.RevenueRunCandidateAction, 0, len(candidates))
				for _, c := range candidates {
					amtDisplay := fmt.Sprintf("%.2f", float64(c.Amount)/100)
					out = append(out, subscriptionaction.RevenueRunCandidateAction{
						SubscriptionID:    c.SubscriptionID,
						SubscriptionName:  c.SubscriptionName,
						ClientID:          c.ClientID,
						ClientName:        c.ClientName,
						PlanName:          c.PlanName,
						BillingCycleLabel: c.BillingCycleLabel,
						Currency:          c.Currency,
						PeriodStart:       c.PeriodStart,
						PeriodEnd:         c.PeriodEnd,
						PeriodLabel:       c.PeriodLabel,
						PeriodMarker:      c.PeriodMarker,
						Amount:            c.Amount,
						AmountDisplay:     amtDisplay,
						LineItemCount:     c.LineItemCount,
						Eligible:          c.Eligible,
						BlockerReason:     c.BlockerReason,
					})
				}
				return out, nextCursor, nil
			}
			subActionDeps.GenerateRevenueRun = func(fctx context.Context, scope subscriptionaction.RevenueRunScopeAction, sels subscriptionaction.RevenueRunSelectionsAction) (*subscriptionaction.RevenueRunResultAction, error) {
				consumerSels := consumer.RevenueRunSelections{
					FilterToken: sels.FilterToken,
				}
				for _, s := range sels.ExplicitList {
					consumerSels.ExplicitList = append(consumerSels.ExplicitList, consumer.SelectedRevenueRunCandidate{
						SubscriptionID: s.SubscriptionID,
						PeriodStart:    s.PeriodStart,
						PeriodEnd:      s.PeriodEnd,
						PeriodMarker:   s.PeriodMarker,
					})
				}
				result, err := consumer.GenerateRevenueRun(useCases, fctx, consumer.RevenueRunScope{
					WorkspaceID:    scope.WorkspaceID,
					ClientID:       scope.ClientID,
					SubscriptionID: scope.SubscriptionID,
					AsOfDate:       scope.AsOfDate,
				}, consumerSels)
				if err != nil || result == nil {
					return nil, err
				}
				run := result.Run
				runID := ""
				runStatus := ""
				if run != nil {
					runID = run.GetId()
					runStatus = run.GetStatus().String()
				}
				var created, skipped, errored int32
				for _, a := range result.Attempts {
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
		// mark-ready/waive handlers. nil-safe: the espyna subscription
		// composition exposes the server pointer directly (no use-case
		// wrapper yet).
		if useCases.Subscription.BillingEvent != nil {
			be := useCases.Subscription.BillingEvent
			subActionDeps.ListBillingEventsBySubscription = be.ListBySubscription
			subActionDeps.SetBillingEventStatus = be.SetStatus
		}

		// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 / Phase D —
		// wire the JobTemplate read deps that drive the Spawn Jobs
		// section detection on the subscription create drawer + the
		// retroactive spawn drawer. nil-safe.
		if useCases.Operation != nil {
			if uc := useCases.Operation.JobTemplate; uc != nil && uc.ReadJobTemplate != nil {
				subActionDeps.ReadJobTemplate = uc.ReadJobTemplate.Execute
			}
			if uc := useCases.Operation.JobTemplatePhase; uc != nil && uc.ListByJobTemplate != nil {
				subActionDeps.ListJobTemplatePhases = uc.ListByJobTemplate.Execute
			}
			if uc := useCases.Operation.JobTemplateTask; uc != nil && uc.ListByPhase != nil {
				subActionDeps.ListJobTemplateTasks = uc.ListByPhase.Execute
			}
			if useCases.Operation.JobTemplateRelation != nil {
				subActionDeps.ListJobTemplateRelations = useCases.Operation.JobTemplateRelation.ListByParent
			}
		}
		if useCases.Subscription.MaterializeJobsForSubscription != nil {
			subActionDeps.MaterializeJobsForSubscription = func(fctx context.Context, subID string, spawn bool) (int, string, error) {
				resp, err := consumer.MaterializeJobsForSubscription(useCases, fctx, &consumer.MaterializeJobsForSubscriptionRequest{
					SubscriptionID: subID,
					SpawnJobs:      spawn,
				})
				if err != nil {
					return 0, "", err
				}
				if resp == nil {
					return 0, "", nil
				}
				return resp.JobCount, resp.SkippedReason, nil
			}
		}
		// 2026-04-30 cyclic-subscription-jobs plan §5.3 / Phase D —
		// wire espyna's MaterializeInstanceJobsForSubscription through
		// a centymo-side adapter so the Operations tab "Spawn this
		// cycle now" + "Backfill missing cycles" handlers can call it
		// without importing espyna directly. nil-safe: the cycle-spawn
		// and backfill action handlers gate on the adapter pointer.
		if useCases.Subscription.MaterializeInstanceJobsForSubscription != nil {
			subActionDeps.MaterializeInstanceJobsForSubscription = func(fctx context.Context, req *subscriptionaction.MaterializeInstanceJobsRequest) (*subscriptionaction.MaterializeInstanceJobsResponse, error) {
				if req == nil {
					return nil, nil
				}
				resp, err := consumer.MaterializeInstanceJobsForSubscription(useCases, fctx, &consumer.MaterializeInstanceJobsForSubscriptionRequest{
					SubscriptionID:   req.SubscriptionID,
					CyclePeriodStart: req.CyclePeriodStart,
					Backfill:         req.Backfill,
				})
				if err != nil {
					return nil, err
				}
				if resp == nil {
					return &subscriptionaction.MaterializeInstanceJobsResponse{}, nil
				}
				return &subscriptionaction.MaterializeInstanceJobsResponse{
					SpawnedCycleCount:         resp.SpawnedCycleCount,
					SpawnedJobCount:           resp.SpawnedJobCount,
					OnceAtStartJobCount:       resp.OnceAtStartJobCount,
					EngagementWasNewlyCreated: resp.EngagementWasNewlyCreated,
					SkippedReason:             resp.SkippedReason,
					BackfillCappedAt:          resp.BackfillCappedAt,
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
	if useCases.Subscription != nil && useCases.Subscription.Subscription != nil && useCases.Subscription.Subscription.ReadSubscription != nil {
		subDetailDeps := &subscriptiondetail.DetailViewDeps{
			Routes:           w.subscriptionRoutes,
			ReadSubscription: useCases.Subscription.Subscription.ReadSubscription.Execute,
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
		if useCases.Subscription.Subscription.GetSubscriptionItemPageData != nil {
			subDetailDeps.GetSubscriptionItemPageData = useCases.Subscription.Subscription.GetSubscriptionItemPageData.Execute
		}
		if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ReadClient != nil {
			subDetailDeps.ReadClient = useCases.Entity.Client.ReadClient.Execute
		}
		if useCases.Revenue != nil && useCases.Revenue.Revenue != nil && useCases.Revenue.Revenue.GetRevenueListPageData != nil {
			subDetailDeps.GetRevenueListPageData = useCases.Revenue.Revenue.GetRevenueListPageData.Execute
		}
		// 2026-04-29 milestone-billing — wire BillingEvent listing into
		// the subscription detail Package tab.
		if useCases.Subscription.BillingEvent != nil {
			subDetailDeps.ListBillingEventsBySubscription = useCases.Subscription.BillingEvent.ListBySubscription
		}
		// 2026-04-29 auto-spawn-jobs-from-subscription Phase D — wire
		// the Operations tab data ops + spawn-jobs CTA URL.
		if useCases.Operation != nil {
			if uc := useCases.Operation.Job; uc != nil && uc.GetJobsByOrigin != nil {
				subDetailDeps.GetJobsByOrigin = uc.GetJobsByOrigin.Execute
			}
			if uc := useCases.Operation.JobPhase; uc != nil && uc.ListByJob != nil {
				subDetailDeps.ListJobPhasesByJob = uc.ListByJob.Execute
			}
		}
		subDetailDeps.SpawnJobsURL = w.subscriptionRoutes.SpawnJobsURL
		subDetailDeps.JobDetailURL = cfg.jobDetailURL
		subDetailDeps.ClientDetailURL = cfg.clientDetailURL
		// 2026-05-04 — engagement breadcrumb (rate-card → plan).
		subDetailDeps.PriceScheduleDetailURL = w.priceScheduleRoutes.DetailURL
		subDetailDeps.PricePlanDetailURL = w.priceScheduleRoutes.PlanDetailURL
		if useCases.Subscription.PriceSchedule != nil && useCases.Subscription.PriceSchedule.ReadPriceSchedule != nil {
			subDetailDeps.ReadPriceSchedule = useCases.Subscription.PriceSchedule.ReadPriceSchedule.Execute
		}
		if useCases.Subscription.PricePlan != nil && useCases.Subscription.PricePlan.ReadPricePlan != nil {
			subDetailDeps.ReadPricePlan = useCases.Subscription.PricePlan.ReadPricePlan.Execute
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

// wireCustomizePlanForClient threads the espyna Plan.CustomizePlanForClient use
// case into the centymo subscription action Deps. The espyna side ships an
// independent request/response shape; the centymo side uses an in-package
// type-narrow shape so its handlers don't depend directly on the espyna
// generated proto/use-case structs.
//
// When the use case isn't wired (composition layer didn't initialize it), we
// leave the function pointer nil; the handler falls through to a generic
// `customize_failed` toast.
//
// 2026-04-27 plan-client-scope plan §4. Same adapter pattern as
// RecognizeRevenueFromSubscription above.
func wireCustomizePlanForClient(useCases *consumer.UseCases, subActionDeps *subscriptionaction.Deps) {
	if useCases == nil || useCases.Subscription == nil || useCases.Subscription.Plan == nil {
		return
	}
	customizeUC := useCases.Subscription.Plan.CustomizePlanForClient
	if customizeUC == nil {
		return
	}
	_ = customizeUC
	subActionDeps.CustomizePlanForClient = func(
		ctx context.Context, req *subscriptionaction.CustomizePlanForClientRequest,
	) (*subscriptionaction.CustomizePlanForClientResponse, error) {
		resp, err := consumer.CustomizePlanForClient(useCases, ctx, &consumer.CustomizePlanForClientRequest{
			SourcePlanID:      req.SourcePlanID,
			SourcePricePlanID: req.SourcePricePlanID,
			ClientID:          req.ClientID,
			SubscriptionID:    req.SubscriptionID,
			NewScheduleName:   req.NewScheduleName,
		})
		if err != nil {
			return nil, err
		}
		if resp == nil {
			return &subscriptionaction.CustomizePlanForClientResponse{}, nil
		}
		return &subscriptionaction.CustomizePlanForClientResponse{
			NewPlanID:      resp.NewPlanID,
			NewPricePlanID: resp.NewPricePlanID,
			NewScheduleID:  resp.NewScheduleID,
			Reused:         resp.Reused,
		}, nil
	}
}
