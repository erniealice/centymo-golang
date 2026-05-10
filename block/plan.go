// Package block — plan, price-plan, price-schedule, price-list, and plan-bundle domain wiring.
//
// Holds wirePlanModules (the lifted bodies of the wantPricePlan, wantPriceSchedule,
// wantPriceList, and wantPlan branches of Block()).
//
// Phase 4a of the 20260510-block-go-splitting-strategy.
package block

import (
	"context"

	consumer "github.com/erniealice/espyna-golang/consumer"
	"github.com/erniealice/espyna-golang/reference"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"

	"github.com/erniealice/hybra-golang/views/attachment"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	centymo "github.com/erniealice/centymo-golang"
	planaction "github.com/erniealice/centymo-golang/views/plan/action"
	plandetail "github.com/erniealice/centymo-golang/views/plan/detail"
	planlist "github.com/erniealice/centymo-golang/views/plan/list"
	pricelistmod "github.com/erniealice/centymo-golang/views/pricelist"
	priceplanmod "github.com/erniealice/centymo-golang/views/price_plan"
	priceschedulemod "github.com/erniealice/centymo-golang/views/price_schedule"
	priceschedulepricepldetail "github.com/erniealice/centymo-golang/views/price_schedule/detail/plan"
)

// planWiring holds everything wirePlanModules needs from the surrounding Block()
// scope. More than 6 fields → struct. Kept private; never re-exported.
type planWiring struct {
	db           centymo.DataSource
	refChecker   reference.Checker
	// Attachment ops
	uploadFile      func(context.Context, string, string, []byte, string) error
	downloadFile    func(context.Context, string, string) ([]byte, error)
	readAttachment  func(context.Context, *attachmentpb.ReadAttachmentRequest) (*attachmentpb.ReadAttachmentResponse, error)
	listAttachments func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	newAttachmentID func() string
	// Routes
	pricePlanRoutes              centymo.PricePlanRoutes
	priceScheduleRoutes          centymo.PriceScheduleRoutes
	priceScheduleInventoryRoutes centymo.PriceScheduleRoutes
	priceListRoutes              centymo.PriceListRoutes
	planRoutes                   centymo.PlanRoutes
	planBundleRoutes             centymo.PlanRoutes
	subscriptionRoutes           centymo.SubscriptionRoutes
	// Labels
	pricePlanLabels           centymo.PricePlanLabels
	productPricePlanLabels    centymo.ProductPricePlanLabels
	priceScheduleLabels       centymo.PriceScheduleLabels
	priceListLabels           centymo.PriceListLabels
	planLabels                centymo.PlanLabels
	centymoTableLabels        types.TableLabels
}

// wirePlanModules lifts the bodies of the four plan-related `if cfg.wantXxx()`
// branches (PricePlan, PriceSchedule, PriceList, Plan+PlanBundle) from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once at the position where
// the plan wiring used to be.
func wirePlanModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *consumer.UseCases, w planWiring) {
	// =====================================================================
	// Price Plan module (standalone — separate from plan-nested price plans)
	// =====================================================================

	if cfg.wantPricePlan() {
		if useCases.Subscription != nil && useCases.Subscription.PricePlan != nil {
			uc := useCases.Subscription.PricePlan
			var getPricePlanInUseIDs func(context.Context, []string) (map[string]bool, error)
			if w.refChecker != nil {
				getPricePlanInUseIDs = w.refChecker.GetPricePlanInUseIDs
			}
			// 2026-04-27 plan-client-scope plan §6.7 — closure used to look
			// up the parent PriceSchedule's client name for the info banner
			// rendered on the price-plan drawer.
			var ppListClientNames func(ctx context.Context) map[string]string
			if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
				lc := useCases.Entity.Client.ListClients.Execute
				ppListClientNames = func(fctx context.Context) map[string]string {
					out := map[string]string{}
					resp, err := lc(fctx, &clientpb.ListClientsRequest{})
					if err != nil {
						return out
					}
					for _, c := range resp.GetData() {
						label := c.GetName()
						if label == "" {
							if u := c.GetUser(); u != nil {
								label = u.GetFirstName() + " " + u.GetLastName()
							}
						}
						out[c.GetId()] = label
					}
					return out
				}
			}

			pricePlanDeps := &priceplanmod.ModuleDeps{
				Routes:                    w.pricePlanRoutes,
				Labels:                    w.pricePlanLabels,
				ProductPricePlanLabels:    w.productPricePlanLabels,
				PriceScheduleDetailLabels: w.priceScheduleLabels.Detail,
				CommonLabels:              ctx.Common,
				TableLabels:               w.centymoTableLabels,
				ListPricePlans:         uc.ListPricePlans.Execute,
				ReadPricePlan:          uc.ReadPricePlan.Execute,
				CreatePricePlan:        uc.CreatePricePlan.Execute,
				UpdatePricePlan:        uc.UpdatePricePlan.Execute,
				DeletePricePlan:        uc.DeletePricePlan.Execute,
				GetPricePlanInUseIDs:   getPricePlanInUseIDs,
				ListClientNames:        ppListClientNames,
			}
			// Price schedule listing — parent container (owns location + date range)
			if useCases.Subscription.PriceSchedule != nil {
				pricePlanDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
			}
			// Add plan listing if available
			if useCases.Subscription != nil && useCases.Subscription.Plan != nil {
				pricePlanDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
			}
			// Add product listing for detail page product selector
			if useCases.Product != nil && useCases.Product.Product != nil {
				pricePlanDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
			}
			// Add product plan listing for scoping product selector to plan's products
			if useCases.Product != nil && useCases.Product.ProductPlan != nil {
				pricePlanDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
			}
			// Add ProductPricePlan CRUD for detail page
			if useCases.Subscription.ProductPricePlan != nil {
				ppp := useCases.Subscription.ProductPricePlan
				pricePlanDeps.ListProductPricePlans = ppp.ListProductPricePlans.Execute
				pricePlanDeps.CreateProductPricePlan = ppp.CreateProductPricePlan.Execute
				pricePlanDeps.UpdateProductPricePlan = ppp.UpdateProductPricePlan.Execute
				pricePlanDeps.DeleteProductPricePlan = ppp.DeleteProductPricePlan.Execute
			}
			// 2026-04-29 milestone-billing plan §5 / Phase D — milestone phase
			// select on the PPP drawer needs ReadPlan (to resolve job_template_id)
			// and ListByJobTemplate (to load phase rows).
			if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
				pricePlanDeps.ReadPlan = useCases.Subscription.Plan.ReadPlan.Execute
			}
			if useCases.Operation != nil && useCases.Operation.JobTemplatePhase != nil && useCases.Operation.JobTemplatePhase.ListByJobTemplate != nil {
				pricePlanDeps.ListJobTemplatePhasesByJobTemplate = useCases.Operation.JobTemplatePhase.ListByJobTemplate.Execute
			}
			pricePlanDeps.UploadFile = w.uploadFile
			pricePlanDeps.ListAttachments = w.listAttachments
			pricePlanDeps.CreateAttachment = w.createAttachment
			pricePlanDeps.DeleteAttachment = w.deleteAttachment
			pricePlanDeps.NewAttachmentID = w.newAttachmentID
			priceplanmod.NewModule(pricePlanDeps).RegisterRoutes(ctx.Routes)
		}
	}

	// =====================================================================
	// Price Schedule module
	// =====================================================================

	if cfg.wantPriceSchedule() {
		if useCases.Subscription != nil && useCases.Subscription.PriceSchedule != nil {
			uc := useCases.Subscription.PriceSchedule
			var getPriceScheduleInUseIDs func(context.Context, []string) (map[string]bool, error)
			if w.refChecker != nil {
				getPriceScheduleInUseIDs = w.refChecker.GetPriceScheduleInUseIDs
			}
			// 2026-04-27 plan-client-scope plan §6.1 / §4.4.1 — schedule list
			// Client column lookup + drawer Client picker. Same listClientNames
			// helper used by the plan list.
			var psListClientNames func(ctx context.Context) map[string]string
			if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
				lc := useCases.Entity.Client.ListClients.Execute
				psListClientNames = func(fctx context.Context) map[string]string {
					out := map[string]string{}
					resp, err := lc(fctx, &clientpb.ListClientsRequest{})
					if err != nil {
						return out
					}
					for _, c := range resp.GetData() {
						label := c.GetName()
						if label == "" {
							if u := c.GetUser(); u != nil {
								label = u.GetFirstName() + " " + u.GetLastName()
							}
						}
						out[c.GetId()] = label
					}
					return out
				}
			}

			priceScheduleDeps := &priceschedulemod.ModuleDeps{
				Routes:                   w.priceScheduleRoutes,
				Labels:                   w.priceScheduleLabels,
				PricePlanLabels:          w.pricePlanLabels,
				ProductPricePlanLabels:   w.productPricePlanLabels,
				CommonLabels:             ctx.Common,
				TableLabels:              w.centymoTableLabels,
				ListPriceSchedules:       uc.ListPriceSchedules.Execute,
				ReadPriceSchedule:        uc.ReadPriceSchedule.Execute,
				CreatePriceSchedule:      uc.CreatePriceSchedule.Execute,
				UpdatePriceSchedule:      uc.UpdatePriceSchedule.Execute,
				DeletePriceSchedule:      uc.DeletePriceSchedule.Execute,
				GetPriceScheduleInUseIDs: getPriceScheduleInUseIDs,
				ListClientNames:          psListClientNames,
			}
			// 2026-04-27 plan-client-scope plan §6.7 / §4.4.1 — Client picker
			// for the schedule add/edit drawer.
			if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
				priceScheduleDeps.ListClients = useCases.Entity.Client.ListClients.Execute
			}
			// Add location listing if available
			if useCases.Entity != nil && useCases.Entity.Location != nil {
				priceScheduleDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
			}
			// Plans tab on the detail page lists price_plans filtered by price_schedule_id FK.
			if useCases.Subscription.PricePlan != nil {
				priceScheduleDeps.ListPricePlans = useCases.Subscription.PricePlan.ListPricePlans.Execute
				priceScheduleDeps.CreatePricePlan = useCases.Subscription.PricePlan.CreatePricePlan.Execute
				priceScheduleDeps.ReadPricePlan = useCases.Subscription.PricePlan.ReadPricePlan.Execute
				priceScheduleDeps.UpdatePricePlan = useCases.Subscription.PricePlan.UpdatePricePlan.Execute
				priceScheduleDeps.DeletePricePlan = useCases.Subscription.PricePlan.DeletePricePlan.Execute
			}
			// Reference checker for in-use guard (disables row Delete + locks pricing fields
			// on the edit drawer when a price_plan is referenced by active subscriptions).
			if w.refChecker != nil {
				priceScheduleDeps.GetPricePlanInUseIDs = w.refChecker.GetPricePlanInUseIDs
			}
			if useCases.Subscription.Plan != nil {
				priceScheduleDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
			}
			// Schedule-scoped plan detail (info + product-prices tabs) needs product lookups + ProductPricePlan CRUD
			if useCases.Product != nil && useCases.Product.Product != nil {
				priceScheduleDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
			}
			if useCases.Product != nil && useCases.Product.ProductPlan != nil {
				priceScheduleDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
			}
			if useCases.Subscription.ProductPricePlan != nil {
				ppp := useCases.Subscription.ProductPricePlan
				priceScheduleDeps.ListProductPricePlans = ppp.ListProductPricePlans.Execute
				priceScheduleDeps.CreateProductPricePlan = ppp.CreateProductPricePlan.Execute
				priceScheduleDeps.UpdateProductPricePlan = ppp.UpdateProductPricePlan.Execute
				priceScheduleDeps.DeleteProductPricePlan = ppp.DeleteProductPricePlan.Execute
			}
			// 2026-05-04 — Engagements (subscriptions) tab on the schedule-scoped
			// price_plan detail page. See docs/plan/20260504-price-plan-engagements-tab/.
			if useCases.Subscription.Subscription != nil && useCases.Subscription.Subscription.ListSubscriptionsByPricePlan != nil {
				priceScheduleDeps.ListSubscriptionsByPricePlan = useCases.Subscription.Subscription.ListSubscriptionsByPricePlan.Execute
			}
			priceScheduleDeps.SubscriptionDetailURL = w.subscriptionRoutes.DetailURL
			priceScheduleDeps.SubscriptionEditURL = w.subscriptionRoutes.EditURL
			priceScheduleDeps.SubscriptionDeleteURL = w.subscriptionRoutes.DeleteURL
			priceScheduleDeps.UploadFile = w.uploadFile
			priceScheduleDeps.ListAttachments = w.listAttachments
			priceScheduleDeps.CreateAttachment = w.createAttachment
			priceScheduleDeps.DeleteAttachment = w.deleteAttachment
			priceScheduleDeps.NewAttachmentID = w.newAttachmentID
			priceschedulemod.NewModule(priceScheduleDeps).RegisterRoutes(ctx.Routes)

			// =====================================================================
			// PriceSchedule inventory-mount (second registration on distinct URLs)
			// =====================================================================
			// Reuses the same PriceSchedule views but on /app/inventory/price-schedules/*.
			// Gate: if a lyngua price_schedule_inventory override ever collapses ListURL
			// back onto the services mount, skip to avoid a ServeMux duplicate-route panic.
			if w.priceScheduleInventoryRoutes.ListURL != w.priceScheduleRoutes.ListURL {
				priceScheduleInventoryDeps := &priceschedulemod.ModuleDeps{
					Routes:                   w.priceScheduleInventoryRoutes,
					Labels:                   w.priceScheduleLabels,
					PricePlanLabels:          w.pricePlanLabels,
					ProductPricePlanLabels:   w.productPricePlanLabels,
					CommonLabels:             ctx.Common,
					TableLabels:              w.centymoTableLabels,
					ListPriceSchedules:       uc.ListPriceSchedules.Execute,
					ReadPriceSchedule:        uc.ReadPriceSchedule.Execute,
					CreatePriceSchedule:      uc.CreatePriceSchedule.Execute,
					UpdatePriceSchedule:      uc.UpdatePriceSchedule.Execute,
					DeletePriceSchedule:      uc.DeletePriceSchedule.Execute,
					GetPriceScheduleInUseIDs: getPriceScheduleInUseIDs,
					ListLocations:            priceScheduleDeps.ListLocations,
					ListPricePlans:           priceScheduleDeps.ListPricePlans,
					CreatePricePlan:          priceScheduleDeps.CreatePricePlan,
					ReadPricePlan:            priceScheduleDeps.ReadPricePlan,
					UpdatePricePlan:          priceScheduleDeps.UpdatePricePlan,
					DeletePricePlan:          priceScheduleDeps.DeletePricePlan,
					GetPricePlanInUseIDs:     priceScheduleDeps.GetPricePlanInUseIDs,
					ListPlans:                priceScheduleDeps.ListPlans,
					ListProducts:             priceScheduleDeps.ListProducts,
					ListProductPlans:         priceScheduleDeps.ListProductPlans,
					ListProductPricePlans:    priceScheduleDeps.ListProductPricePlans,
					CreateProductPricePlan:   priceScheduleDeps.CreateProductPricePlan,
					UpdateProductPricePlan:   priceScheduleDeps.UpdateProductPricePlan,
					DeleteProductPricePlan:   priceScheduleDeps.DeleteProductPricePlan,
					UploadFile:               w.uploadFile,
					ListAttachments:          w.listAttachments,
					CreateAttachment:         w.createAttachment,
					DeleteAttachment:         w.deleteAttachment,
					NewAttachmentID:          w.newAttachmentID,
				}
				priceschedulemod.NewModule(priceScheduleInventoryDeps).RegisterRoutes(ctx.Routes)
			}
		}
	}

	// =====================================================================
	// PriceList module
	// =====================================================================

	if cfg.wantPriceList() {
		var getPriceListInUseIDs func(context.Context, []string) (map[string]bool, error)
		if w.refChecker != nil {
			getPriceListInUseIDs = w.refChecker.GetPriceListInUseIDs
		}

		pricelistmod.NewModule(&pricelistmod.ModuleDeps{
			Routes:             w.priceListRoutes,
			Labels:             w.priceListLabels,
			CommonLabels:       ctx.Common,
			TableLabels:        w.centymoTableLabels,
			GetInUseIDs:        getPriceListInUseIDs,
			ListPriceLists:     useCases.Product.PriceList.ListPriceLists.Execute,
			ReadPriceList:      useCases.Product.PriceList.ReadPriceList.Execute,
			CreatePriceList:    useCases.Product.PriceList.CreatePriceList.Execute,
			UpdatePriceList:    useCases.Product.PriceList.UpdatePriceList.Execute,
			DeletePriceList:    useCases.Product.PriceList.DeletePriceList.Execute,
			ListPriceProducts:  useCases.Product.PriceProduct.ListPriceProducts.Execute,
			CreatePriceProduct: useCases.Product.PriceProduct.CreatePriceProduct.Execute,
			DeletePriceProduct: useCases.Product.PriceProduct.DeletePriceProduct.Execute,
			ListProducts:       useCases.Product.Product.ListProducts.Execute,
			// Attachments
			UploadFile:       w.uploadFile,
			ListAttachments:  w.listAttachments,
			CreateAttachment: w.createAttachment,
			DeleteAttachment: w.deleteAttachment,
			NewID:            w.newAttachmentID,
		}).RegisterRoutes(ctx.Routes)
	}

	// =====================================================================
	// Plan (inline — not a module, uses planlist/planaction/plandetail directly)
	// =====================================================================

	if cfg.wantPlan() {
		// 2026-04-27 plan-client-scope plan §6.1 / §6.2 — client name lookup
		// for the optional Client column on the plan list and for the
		// plan-drawer Client picker label resolution. Falls back to the
		// bare client_id when no use case is wired (e.g. tests).
		var listClientNames func(ctx context.Context) map[string]string
		if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
			lc := useCases.Entity.Client.ListClients.Execute
			listClientNames = func(fctx context.Context) map[string]string {
				out := map[string]string{}
				resp, err := lc(fctx, &clientpb.ListClientsRequest{})
				if err != nil {
					return out
				}
				for _, c := range resp.GetData() {
					label := c.GetName()
					if label == "" {
						if u := c.GetUser(); u != nil {
							label = u.GetFirstName() + " " + u.GetLastName()
						}
					}
					out[c.GetId()] = label
				}
				return out
			}
		}

		// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — job
		// template name lookup for the Job Template column on the plan
		// list. Mirrors the listClientNames pattern; falls back to the
		// bare job_template_id when the use case is unwired.
		var listJobTemplateNames func(ctx context.Context) map[string]string
		if useCases.Operation != nil && useCases.Operation.JobTemplate != nil && useCases.Operation.JobTemplate.ListJobTemplates != nil {
			ljt := useCases.Operation.JobTemplate.ListJobTemplates.Execute
			listJobTemplateNames = func(fctx context.Context) map[string]string {
				out := map[string]string{}
				resp, err := ljt(fctx, &jobtemplatepb.ListJobTemplatesRequest{})
				if err != nil {
					return out
				}
				for _, t := range resp.GetData() {
					if t == nil {
						continue
					}
					out[t.GetId()] = t.GetName()
				}
				return out
			}
		}

		planListDeps := &planlist.ListViewDeps{
			Routes:               w.planRoutes,
			Labels:               w.planLabels,
			CommonLabels:         ctx.Common,
			TableLabels:          w.centymoTableLabels,
			ListClientNames:      listClientNames,
			ListJobTemplateNames: listJobTemplateNames,
		}
		if useCases.Subscription != nil && useCases.Subscription.Plan != nil {
			planListDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
		}
		if w.refChecker != nil {
			planListDeps.GetInUseIDs = w.refChecker.GetPlanInUseIDs
		}
		ctx.Routes.GET(w.planRoutes.ListURL, planlist.NewView(planListDeps))
		ctx.Routes.GET(w.planRoutes.TableURL, planlist.NewTableView(planListDeps))

		// Plan CRUD actions
		if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.CreatePlan != nil {
			planActionDeps := &planaction.Deps{
				Routes:     w.planRoutes,
				Labels:     w.planLabels,
				CreatePlan: useCases.Subscription.Plan.CreatePlan.Execute,
				ReadPlan:   useCases.Subscription.Plan.ReadPlan.Execute,
				UpdatePlan: useCases.Subscription.Plan.UpdatePlan.Execute,
				DeletePlan: useCases.Subscription.Plan.DeletePlan.Execute,
				// SetPlanActive uses raw DB update (proto3 omits false booleans)
				SetPlanActive: func(fctx context.Context, id string, active bool) error {
					_, err := w.db.Update(fctx, "plan", id, map[string]any{"active": active})
					return err
				},
			}
			// 2026-04-27 plan-client-scope plan §6.2 — Client picker support
			// + reference-checker lock state.
			if useCases.Entity != nil && useCases.Entity.Client != nil {
				if useCases.Entity.Client.ListClients != nil {
					planActionDeps.ListClients = useCases.Entity.Client.ListClients.Execute
				}
				if useCases.Entity.Client.SearchClientsByName != nil {
					planActionDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName.Execute
				}
			}
			if w.refChecker != nil {
				planActionDeps.GetPlanClientScopeLockedIDs = w.refChecker.GetPlanClientScopeLockedIDs
			}
			// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 —
			// JobTemplate select for Plan.job_template_id assignment.
			if useCases.Operation != nil && useCases.Operation.JobTemplate != nil && useCases.Operation.JobTemplate.ListJobTemplates != nil {
				planActionDeps.ListJobTemplates = useCases.Operation.JobTemplate.ListJobTemplates.Execute
			}
			ctx.Routes.GET(w.planRoutes.AddURL, planaction.NewAddAction(planActionDeps))
			ctx.Routes.POST(w.planRoutes.AddURL, planaction.NewAddAction(planActionDeps))
			ctx.Routes.GET(w.planRoutes.EditURL, planaction.NewEditAction(planActionDeps))
			ctx.Routes.POST(w.planRoutes.EditURL, planaction.NewEditAction(planActionDeps))
			ctx.Routes.POST(w.planRoutes.DeleteURL, planaction.NewDeleteAction(planActionDeps))
			ctx.Routes.POST(w.planRoutes.BulkDeleteURL, planaction.NewBulkDeleteAction(planActionDeps))
			ctx.Routes.POST(w.planRoutes.SetStatusURL, planaction.NewSetStatusAction(planActionDeps))
			ctx.Routes.POST(w.planRoutes.BulkSetStatusURL, planaction.NewBulkSetStatusAction(planActionDeps))
		}

		// Plan detail page + tab action
		if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
			planDetailDeps := &plandetail.DetailViewDeps{
				Routes:                     w.planRoutes,
				PriceSchedulePlanDetailURL: w.priceScheduleRoutes.PlanDetailURL,
				ReadPlan:                   useCases.Subscription.Plan.ReadPlan.Execute,
				Labels:                     w.planLabels,
				CommonLabels:               ctx.Common,
				TableLabels:                w.centymoTableLabels,
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
			if useCases.Product != nil && useCases.Product.ProductPlan != nil {
				planDetailDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
			}
			if useCases.Product != nil && useCases.Product.Product != nil && useCases.Product.Product.ListProducts != nil {
				planDetailDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
			}
			if useCases.Product != nil && useCases.Product.ProductVariant != nil {
				planDetailDeps.ListProductVariants = useCases.Product.ProductVariant.ListProductVariants.Execute
			}
			if useCases.Subscription.PricePlan != nil {
				planDetailDeps.ListPricePlans = useCases.Subscription.PricePlan.ListPricePlans.Execute
			}
			if useCases.Entity != nil && useCases.Entity.Location != nil {
				planDetailDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
			}
			if useCases.Subscription.PriceSchedule != nil {
				planDetailDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
			}
			// 2026-04-28 plan-client-scope — Info tab Client row needs to
			// resolve the plan's client_id label and (optionally) link to
			// the entydad client-detail page.
			if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
				planDetailDeps.ListClients = useCases.Entity.Client.ListClients.Execute
			}
			planDetailDeps.ClientDetailURL = cfg.clientDetailURL
			// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — Info
			// tab JobTemplate row resolution.
			if useCases.Operation != nil && useCases.Operation.JobTemplate != nil && useCases.Operation.JobTemplate.ReadJobTemplate != nil {
				planDetailDeps.ReadJobTemplate = useCases.Operation.JobTemplate.ReadJobTemplate.Execute
			}
			ctx.Routes.GET(w.planRoutes.DetailURL, plandetail.NewView(planDetailDeps))
			ctx.Routes.GET(w.planRoutes.TabActionURL, plandetail.NewTabAction(planDetailDeps))

			// Plan-scoped PricePlan detail (/app/plans/detail/{id}/price/{ppid}).
			// Reuses the schedule-scoped detail body but anchors ActiveNav to
			// Services > Packages and points the breadcrumb back at the
			// package's package-prices tab. The {id} path value is plan_id;
			// the handler resolves schedule_id from the price_plan record.
			if w.planRoutes.PricePlanDetailURL != "" && useCases.Subscription.PricePlan != nil {
				// The plan detail's "Package prices" tab is registered under the
				// `pricePlan` key in plan tab labels; the lyngua professional
				// override surfaces it as the slug "package-prices" in the URL.
				packagePricesSlug := w.planLabels.Tabs.ResolveTabSlug("pricePlan")
				planScopedDeps := &priceschedulepricepldetail.DetailViewDeps{
					Routes:                 w.priceScheduleRoutes,
					ScheduleLabels:         w.priceScheduleLabels,
					PlanLabels:             w.pricePlanLabels,
					ProductPricePlanLabels: w.productPricePlanLabels,
					CommonLabels:           ctx.Common,
					TableLabels:            w.centymoTableLabels,
					ReadPricePlan:          useCases.Subscription.PricePlan.ReadPricePlan.Execute,
					// Mount overrides — keep the page anchored to Packages.
					ActiveNavOverride:      w.planRoutes.ActiveNav,
					ActiveSubNavOverride:   w.planRoutes.ActiveSubNav,
					PlanDetailBackURL:      w.planRoutes.DetailURL,
					PlanDetailBackTab:      packagePricesSlug,
					PlanScopedDetailURL:    w.planRoutes.PricePlanDetailURL,
					PlanScopedTabActionURL: w.planRoutes.PricePlanTabActionURL,
				}
				if useCases.Subscription.PriceSchedule != nil {
					planScopedDeps.ReadPriceSchedule = useCases.Subscription.PriceSchedule.ReadPriceSchedule.Execute
				}
				if useCases.Subscription.Plan != nil {
					planScopedDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
				}
				if useCases.Product != nil && useCases.Product.Product != nil {
					planScopedDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
				}
				if useCases.Product != nil && useCases.Product.ProductPlan != nil {
					planScopedDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
				}
				if useCases.Product != nil && useCases.Product.ProductVariant != nil {
					planScopedDeps.ListProductVariants = useCases.Product.ProductVariant.ListProductVariants.Execute
				}
				if useCases.Subscription.ProductPricePlan != nil {
					planScopedDeps.ListProductPricePlans = useCases.Subscription.ProductPricePlan.ListProductPricePlans.Execute
				}
				ctx.Routes.GET(w.planRoutes.PricePlanDetailURL, priceschedulepricepldetail.NewPlanScopedView(planScopedDeps))
				if w.planRoutes.PricePlanTabActionURL != "" {
					ctx.Routes.GET(w.planRoutes.PricePlanTabActionURL, priceschedulepricepldetail.NewPlanScopedTabAction(planScopedDeps))
				}
			}
			// Plan attachments
			if w.uploadFile != nil {
				ctx.Routes.GET(w.planRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planDetailDeps))
				ctx.Routes.POST(w.planRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planDetailDeps))
				ctx.Routes.POST(w.planRoutes.AttachmentDeleteURL, plandetail.NewAttachmentDeleteAction(planDetailDeps))
			}
			// PricePlan CRUD within plan detail
			if useCases.Subscription.PricePlan != nil && useCases.Subscription.PricePlan.CreatePricePlan != nil {
				ppActionDeps := &planaction.PricePlanDeps{
					Routes:              w.planRoutes,
					Labels:              w.planLabels,
					PricePlanLabels:     w.pricePlanLabels,
					PriceScheduleLabels: w.priceScheduleLabels,
					CommonLabels:        ctx.Common,
					CreatePricePlan: useCases.Subscription.PricePlan.CreatePricePlan.Execute,
					ReadPricePlan:   useCases.Subscription.PricePlan.ReadPricePlan.Execute,
					UpdatePricePlan: useCases.Subscription.PricePlan.UpdatePricePlan.Execute,
					DeletePricePlan: useCases.Subscription.PricePlan.DeletePricePlan.Execute,
				}
				if useCases.Subscription.PriceSchedule != nil {
					ppActionDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
				}
				if useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
					ppActionDeps.ReadPlan = useCases.Subscription.Plan.ReadPlan.Execute
				}
				// Plan §6.7 — ListClients powers the readonly schedule
				// label + lock tooltip when the parent Plan is client-scoped.
				if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
					ppActionDeps.ListClients = useCases.Entity.Client.ListClients.Execute
				}
				if useCases.Entity != nil && useCases.Entity.Location != nil && useCases.Entity.Location.ListLocations != nil {
					ppActionDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
				}
				if w.refChecker != nil {
					ppActionDeps.GetPricePlanInUseIDs = w.refChecker.GetPricePlanInUseIDs
				}
				if useCases.Product != nil && useCases.Product.Product != nil && useCases.Product.Product.ListProducts != nil {
					ppActionDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
				}
				if useCases.Product != nil && useCases.Product.ProductPlan != nil {
					ppActionDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
				}
				if useCases.Subscription.ProductPricePlan != nil {
					ppActionDeps.CreateProductPricePlan = useCases.Subscription.ProductPricePlan.CreateProductPricePlan.Execute
					ppActionDeps.ListProductPricePlans = useCases.Subscription.ProductPricePlan.ListProductPricePlans.Execute
				}
				ctx.Routes.GET(w.planRoutes.PricePlanAddURL, planaction.NewPricePlanAddAction(ppActionDeps))
				ctx.Routes.POST(w.planRoutes.PricePlanAddURL, planaction.NewPricePlanAddAction(ppActionDeps))
				ctx.Routes.GET(w.planRoutes.PricePlanEditURL, planaction.NewPricePlanEditAction(ppActionDeps))
				ctx.Routes.POST(w.planRoutes.PricePlanEditURL, planaction.NewPricePlanEditAction(ppActionDeps))
				ctx.Routes.POST(w.planRoutes.PricePlanDeleteURL, planaction.NewPricePlanDeleteAction(ppActionDeps))
			}
			// ProductPlan CRUD within plan detail
			if useCases.Product != nil && useCases.Product.ProductPlan != nil && useCases.Product.ProductPlan.CreateProductPlan != nil {
				productPlanActionDeps := &planaction.ProductPlanDeps{
					Routes:            w.planRoutes,
					Labels:            w.planLabels,
					CreateProductPlan: useCases.Product.ProductPlan.CreateProductPlan.Execute,
					ReadProductPlan:   useCases.Product.ProductPlan.ReadProductPlan.Execute,
					UpdateProductPlan: useCases.Product.ProductPlan.UpdateProductPlan.Execute,
					DeleteProductPlan: useCases.Product.ProductPlan.DeleteProductPlan.Execute,
				}
				if useCases.Product.Product != nil {
					productPlanActionDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
				}
				if useCases.Product.ProductPlan.ListProductPlans != nil {
					productPlanActionDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
				}
				if useCases.Product.ProductVariant != nil {
					productPlanActionDeps.ListProductVariants = useCases.Product.ProductVariant.ListProductVariants.Execute
				}
				if useCases.Product.ProductVariantOption != nil {
					productPlanActionDeps.ListProductVariantOptions = useCases.Product.ProductVariantOption.ListProductVariantOptions.Execute
				}
				if useCases.Product.ProductOptionValue != nil {
					productPlanActionDeps.ListProductOptionValues = useCases.Product.ProductOptionValue.ListProductOptionValues.Execute
				}
				if useCases.Product.ProductOption != nil {
					productPlanActionDeps.ListProductOptions = useCases.Product.ProductOption.ListProductOptions.Execute
				}
				ctx.Routes.GET(w.planRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(productPlanActionDeps))
				ctx.Routes.POST(w.planRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(productPlanActionDeps))
				ctx.Routes.GET(w.planRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(productPlanActionDeps))
				ctx.Routes.POST(w.planRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(productPlanActionDeps))
				ctx.Routes.POST(w.planRoutes.ProductPlanDeleteURL, planaction.NewProductPlanDeleteAction(productPlanActionDeps))
				ctx.Routes.GET(w.planRoutes.ProductPlanPickerURL, planaction.NewProductPlanPickerAction(productPlanActionDeps))
			}
		}

		// =====================================================================
		// Plan bundle inventory-mount (second registration on distinct URLs)
		// =====================================================================
		// Reuses the same plan views but on /app/inventory/bundles/* URLs.
		// Gate: if a lyngua plan_bundle override ever collapses ListURL back
		// onto the services mount, skip to avoid a ServeMux duplicate-route panic.
		if cfg.wantPlan() && w.planBundleRoutes.ListURL != w.planRoutes.ListURL {
			planBundleListDeps := &planlist.ListViewDeps{
				Routes:               w.planBundleRoutes,
				Labels:               w.planLabels,
				CommonLabels:         ctx.Common,
				TableLabels:          w.centymoTableLabels,
				ListClientNames:      listClientNames,
				ListJobTemplateNames: listJobTemplateNames,
			}
			if useCases.Subscription != nil && useCases.Subscription.Plan != nil {
				planBundleListDeps.ListPlans = useCases.Subscription.Plan.ListPlans.Execute
			}
			if w.refChecker != nil {
				planBundleListDeps.GetInUseIDs = w.refChecker.GetPlanInUseIDs
			}
			ctx.Routes.GET(w.planBundleRoutes.ListURL, planlist.NewView(planBundleListDeps))
			ctx.Routes.GET(w.planBundleRoutes.TableURL, planlist.NewTableView(planBundleListDeps))

			if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.CreatePlan != nil {
				planBundleActionDeps := &planaction.Deps{
					Routes:     w.planBundleRoutes,
					Labels:     w.planLabels,
					CreatePlan: useCases.Subscription.Plan.CreatePlan.Execute,
					ReadPlan:   useCases.Subscription.Plan.ReadPlan.Execute,
					UpdatePlan: useCases.Subscription.Plan.UpdatePlan.Execute,
					DeletePlan: useCases.Subscription.Plan.DeletePlan.Execute,
					// SetPlanActive uses raw DB update (proto3 omits false booleans)
					SetPlanActive: func(fctx context.Context, id string, active bool) error {
						_, err := w.db.Update(fctx, "plan", id, map[string]any{"active": active})
						return err
					},
				}
				// 2026-04-27 plan-client-scope plan §6.2 — same Client picker
				// + lock state on the bundle mount.
				if useCases.Entity != nil && useCases.Entity.Client != nil {
					if useCases.Entity.Client.ListClients != nil {
						planBundleActionDeps.ListClients = useCases.Entity.Client.ListClients.Execute
					}
					if useCases.Entity.Client.SearchClientsByName != nil {
						planBundleActionDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName.Execute
					}
				}
				if w.refChecker != nil {
					planBundleActionDeps.GetPlanClientScopeLockedIDs = w.refChecker.GetPlanClientScopeLockedIDs
				}
				// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 —
				// JobTemplate select on the bundle mount.
				if useCases.Operation != nil && useCases.Operation.JobTemplate != nil && useCases.Operation.JobTemplate.ListJobTemplates != nil {
					planBundleActionDeps.ListJobTemplates = useCases.Operation.JobTemplate.ListJobTemplates.Execute
				}
				ctx.Routes.GET(w.planBundleRoutes.AddURL, planaction.NewAddAction(planBundleActionDeps))
				ctx.Routes.POST(w.planBundleRoutes.AddURL, planaction.NewAddAction(planBundleActionDeps))
				ctx.Routes.GET(w.planBundleRoutes.EditURL, planaction.NewEditAction(planBundleActionDeps))
				ctx.Routes.POST(w.planBundleRoutes.EditURL, planaction.NewEditAction(planBundleActionDeps))
				ctx.Routes.POST(w.planBundleRoutes.DeleteURL, planaction.NewDeleteAction(planBundleActionDeps))
				ctx.Routes.POST(w.planBundleRoutes.BulkDeleteURL, planaction.NewBulkDeleteAction(planBundleActionDeps))
				ctx.Routes.POST(w.planBundleRoutes.SetStatusURL, planaction.NewSetStatusAction(planBundleActionDeps))
				ctx.Routes.POST(w.planBundleRoutes.BulkSetStatusURL, planaction.NewBulkSetStatusAction(planBundleActionDeps))
			}

			if useCases.Subscription != nil && useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
				planBundleDetailDeps := &plandetail.DetailViewDeps{
					Routes:                     w.planBundleRoutes,
					PriceSchedulePlanDetailURL: w.priceScheduleRoutes.PlanDetailURL,
					ReadPlan:                   useCases.Subscription.Plan.ReadPlan.Execute,
					Labels:                     w.planLabels,
					CommonLabels:               ctx.Common,
					TableLabels:                w.centymoTableLabels,
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
				if useCases.Product != nil && useCases.Product.ProductPlan != nil {
					planBundleDetailDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
				}
				if useCases.Product != nil && useCases.Product.Product != nil && useCases.Product.Product.ListProducts != nil {
					planBundleDetailDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
				}
				if useCases.Product != nil && useCases.Product.ProductVariant != nil {
					planBundleDetailDeps.ListProductVariants = useCases.Product.ProductVariant.ListProductVariants.Execute
				}
				if useCases.Subscription.PricePlan != nil {
					planBundleDetailDeps.ListPricePlans = useCases.Subscription.PricePlan.ListPricePlans.Execute
				}
				if useCases.Entity != nil && useCases.Entity.Location != nil {
					planBundleDetailDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
				}
				if useCases.Subscription.PriceSchedule != nil {
					planBundleDetailDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
				}
				// 2026-04-28 plan-client-scope — same Info tab Client row
				// wiring on the bundle mount.
				if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
					planBundleDetailDeps.ListClients = useCases.Entity.Client.ListClients.Execute
				}
				planBundleDetailDeps.ClientDetailURL = cfg.clientDetailURL
				// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 —
				// Info tab JobTemplate row on the bundle mount.
				if useCases.Operation != nil && useCases.Operation.JobTemplate != nil && useCases.Operation.JobTemplate.ReadJobTemplate != nil {
					planBundleDetailDeps.ReadJobTemplate = useCases.Operation.JobTemplate.ReadJobTemplate.Execute
				}
				ctx.Routes.GET(w.planBundleRoutes.DetailURL, plandetail.NewView(planBundleDetailDeps))
				ctx.Routes.GET(w.planBundleRoutes.TabActionURL, plandetail.NewTabAction(planBundleDetailDeps))
				if w.uploadFile != nil {
					ctx.Routes.GET(w.planBundleRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planBundleDetailDeps))
					ctx.Routes.POST(w.planBundleRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planBundleDetailDeps))
					ctx.Routes.POST(w.planBundleRoutes.AttachmentDeleteURL, plandetail.NewAttachmentDeleteAction(planBundleDetailDeps))
				}
				if useCases.Subscription.PricePlan != nil && useCases.Subscription.PricePlan.CreatePricePlan != nil {
					ppBundleDeps := &planaction.PricePlanDeps{
						Routes:              w.planBundleRoutes,
						Labels:              w.planLabels,
						PricePlanLabels:     w.pricePlanLabels,
						PriceScheduleLabels: w.priceScheduleLabels,
						CommonLabels:        ctx.Common,
						CreatePricePlan: useCases.Subscription.PricePlan.CreatePricePlan.Execute,
						ReadPricePlan:   useCases.Subscription.PricePlan.ReadPricePlan.Execute,
						UpdatePricePlan: useCases.Subscription.PricePlan.UpdatePricePlan.Execute,
						DeletePricePlan: useCases.Subscription.PricePlan.DeletePricePlan.Execute,
					}
					if useCases.Subscription.PriceSchedule != nil {
						ppBundleDeps.ListPriceSchedules = useCases.Subscription.PriceSchedule.ListPriceSchedules.Execute
					}
					if useCases.Subscription.Plan != nil && useCases.Subscription.Plan.ReadPlan != nil {
						ppBundleDeps.ReadPlan = useCases.Subscription.Plan.ReadPlan.Execute
					}
					// Plan §6.7 — ListClients powers the readonly schedule
					// label + lock tooltip on the bundle-mount drawer.
					if useCases.Entity != nil && useCases.Entity.Client != nil && useCases.Entity.Client.ListClients != nil {
						ppBundleDeps.ListClients = useCases.Entity.Client.ListClients.Execute
					}
					if useCases.Entity != nil && useCases.Entity.Location != nil && useCases.Entity.Location.ListLocations != nil {
						ppBundleDeps.ListLocations = useCases.Entity.Location.ListLocations.Execute
					}
					if w.refChecker != nil {
						ppBundleDeps.GetPricePlanInUseIDs = w.refChecker.GetPricePlanInUseIDs
					}
					if useCases.Product != nil && useCases.Product.Product != nil && useCases.Product.Product.ListProducts != nil {
						ppBundleDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
					}
					if useCases.Product != nil && useCases.Product.ProductPlan != nil {
						ppBundleDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
					}
					if useCases.Subscription.ProductPricePlan != nil {
						ppBundleDeps.CreateProductPricePlan = useCases.Subscription.ProductPricePlan.CreateProductPricePlan.Execute
						ppBundleDeps.ListProductPricePlans = useCases.Subscription.ProductPricePlan.ListProductPricePlans.Execute
					}
					ctx.Routes.GET(w.planBundleRoutes.PricePlanAddURL, planaction.NewPricePlanAddAction(ppBundleDeps))
					ctx.Routes.POST(w.planBundleRoutes.PricePlanAddURL, planaction.NewPricePlanAddAction(ppBundleDeps))
					ctx.Routes.GET(w.planBundleRoutes.PricePlanEditURL, planaction.NewPricePlanEditAction(ppBundleDeps))
					ctx.Routes.POST(w.planBundleRoutes.PricePlanEditURL, planaction.NewPricePlanEditAction(ppBundleDeps))
					ctx.Routes.POST(w.planBundleRoutes.PricePlanDeleteURL, planaction.NewPricePlanDeleteAction(ppBundleDeps))
				}
				// Bundle-mount sibling of services-mount `productPlanActionDeps` (~line 1111).
				// Keep these two registrations field-for-field identical (only Routes
				// differs). Unlike PricePlanDeps, ProductPlanDeps has a single Labels
				// field — all form-label data is nested under centymo.PlanLabels
				// (`Labels.ProductPlanForm`), so threading `Labels: planLabels` is
				// sufficient. If a future change adds a separate label struct (e.g.
				// ProductPlanLabels), thread it into BOTH registrations.
				if useCases.Product != nil && useCases.Product.ProductPlan != nil && useCases.Product.ProductPlan.CreateProductPlan != nil {
					ppBundleProductPlanDeps := &planaction.ProductPlanDeps{
						Routes:            w.planBundleRoutes,
						Labels:            w.planLabels,
						CreateProductPlan: useCases.Product.ProductPlan.CreateProductPlan.Execute,
						ReadProductPlan:   useCases.Product.ProductPlan.ReadProductPlan.Execute,
						UpdateProductPlan: useCases.Product.ProductPlan.UpdateProductPlan.Execute,
						DeleteProductPlan: useCases.Product.ProductPlan.DeleteProductPlan.Execute,
					}
					if useCases.Product.Product != nil {
						ppBundleProductPlanDeps.ListProducts = useCases.Product.Product.ListProducts.Execute
					}
					if useCases.Product.ProductPlan.ListProductPlans != nil {
						ppBundleProductPlanDeps.ListProductPlans = useCases.Product.ProductPlan.ListProductPlans.Execute
					}
					if useCases.Product.ProductVariant != nil {
						ppBundleProductPlanDeps.ListProductVariants = useCases.Product.ProductVariant.ListProductVariants.Execute
					}
					if useCases.Product.ProductVariantOption != nil {
						ppBundleProductPlanDeps.ListProductVariantOptions = useCases.Product.ProductVariantOption.ListProductVariantOptions.Execute
					}
					if useCases.Product.ProductOptionValue != nil {
						ppBundleProductPlanDeps.ListProductOptionValues = useCases.Product.ProductOptionValue.ListProductOptionValues.Execute
					}
					if useCases.Product.ProductOption != nil {
						ppBundleProductPlanDeps.ListProductOptions = useCases.Product.ProductOption.ListProductOptions.Execute
					}
					ctx.Routes.GET(w.planBundleRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(ppBundleProductPlanDeps))
					ctx.Routes.POST(w.planBundleRoutes.ProductPlanAddURL, planaction.NewProductPlanAddAction(ppBundleProductPlanDeps))
					ctx.Routes.GET(w.planBundleRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(ppBundleProductPlanDeps))
					ctx.Routes.POST(w.planBundleRoutes.ProductPlanEditURL, planaction.NewProductPlanEditAction(ppBundleProductPlanDeps))
					ctx.Routes.POST(w.planBundleRoutes.ProductPlanDeleteURL, planaction.NewProductPlanDeleteAction(ppBundleProductPlanDeps))
					ctx.Routes.GET(w.planBundleRoutes.ProductPlanPickerURL, planaction.NewProductPlanPickerAction(ppBundleProductPlanDeps))
				}
			}
		}
	}
}
