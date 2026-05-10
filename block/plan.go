// Package block — plan, price-plan, price-schedule, price-list, and plan-bundle domain wiring.
//
// Holds wirePlanModules (the lifted bodies of the wantPricePlan, wantPriceSchedule,
// wantPriceList, and wantPlan branches of Block()).
//
// Phase 4a of the 20260510-block-go-splitting-strategy.
package block

import (
	"context"

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
func wirePlanModules(ctx *pyeza.AppContext, cfg *blockConfig, useCases *UseCases, w planWiring) {
	// =====================================================================
	// Price Plan module (standalone — separate from plan-nested price plans)
	// =====================================================================

	if cfg.wantPricePlan() {
		if useCases.PricePlan.ListPricePlans != nil {
			var getPricePlanInUseIDs func(context.Context, []string) (map[string]bool, error)
			if w.refChecker != nil {
				getPricePlanInUseIDs = w.refChecker.GetPricePlanInUseIDs
			}
			// 2026-04-27 plan-client-scope plan §6.7 — closure used to look
			// up the parent PriceSchedule's client name for the info banner
			// rendered on the price-plan drawer.
			var ppListClientNames func(ctx context.Context) map[string]string
			if useCases.Entity.Client.ListClients != nil {
				lc := useCases.Entity.Client.ListClients
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
				ListPricePlans:         useCases.PricePlan.ListPricePlans,
				ReadPricePlan:          useCases.PricePlan.ReadPricePlan,
				CreatePricePlan:        useCases.PricePlan.CreatePricePlan,
				UpdatePricePlan:        useCases.PricePlan.UpdatePricePlan,
				DeletePricePlan:        useCases.PricePlan.DeletePricePlan,
				GetPricePlanInUseIDs:   getPricePlanInUseIDs,
				ListClientNames:        ppListClientNames,
			}
			// Price schedule listing — parent container (owns location + date range)
			if useCases.PriceSchedule.ListPriceSchedules != nil {
				pricePlanDeps.ListPriceSchedules = useCases.PriceSchedule.ListPriceSchedules
			}
			// Add plan listing if available
			if useCases.Plan.ListPlans != nil {
				pricePlanDeps.ListPlans = useCases.Plan.ListPlans
			}
			// Add product listing for detail page product selector
			if useCases.Product.ListProducts != nil {
				pricePlanDeps.ListProducts = useCases.Product.ListProducts
			}
			// Add product plan listing for scoping product selector to plan's products
			if useCases.Product.ListProductPlans != nil {
				pricePlanDeps.ListProductPlans = useCases.Product.ListProductPlans
			}
			// Add ProductPricePlan CRUD for detail page
			if useCases.PricePlan.ListProductPricePlans != nil {
				pricePlanDeps.ListProductPricePlans = useCases.PricePlan.ListProductPricePlans
				pricePlanDeps.CreateProductPricePlan = useCases.PricePlan.CreateProductPricePlan
				pricePlanDeps.UpdateProductPricePlan = useCases.PricePlan.UpdateProductPricePlan
				pricePlanDeps.DeleteProductPricePlan = useCases.PricePlan.DeleteProductPricePlan
			}
			// 2026-04-29 milestone-billing plan §5 / Phase D — milestone phase
			// select on the PPP drawer needs ReadPlan (to resolve job_template_id)
			// and ListByJobTemplate (to load phase rows).
			if useCases.Plan.ReadPlan != nil {
				pricePlanDeps.ReadPlan = useCases.Plan.ReadPlan
			}
			if useCases.Operation.JobTemplatePhase.ListByJobTemplate != nil {
				pricePlanDeps.ListJobTemplatePhasesByJobTemplate = useCases.Operation.JobTemplatePhase.ListByJobTemplate
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
		if useCases.PriceSchedule.ListPriceSchedules != nil {
			var getPriceScheduleInUseIDs func(context.Context, []string) (map[string]bool, error)
			if w.refChecker != nil {
				getPriceScheduleInUseIDs = w.refChecker.GetPriceScheduleInUseIDs
			}
			// 2026-04-27 plan-client-scope plan §6.1 / §4.4.1 — schedule list
			// Client column lookup + drawer Client picker. Same listClientNames
			// helper used by the plan list.
			var psListClientNames func(ctx context.Context) map[string]string
			if useCases.Entity.Client.ListClients != nil {
				lc := useCases.Entity.Client.ListClients
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
				ListPriceSchedules:       useCases.PriceSchedule.ListPriceSchedules,
				ReadPriceSchedule:        useCases.PriceSchedule.ReadPriceSchedule,
				CreatePriceSchedule:      useCases.PriceSchedule.CreatePriceSchedule,
				UpdatePriceSchedule:      useCases.PriceSchedule.UpdatePriceSchedule,
				DeletePriceSchedule:      useCases.PriceSchedule.DeletePriceSchedule,
				GetPriceScheduleInUseIDs: getPriceScheduleInUseIDs,
				ListClientNames:          psListClientNames,
			}
			// 2026-04-27 plan-client-scope plan §6.7 / §4.4.1 — Client picker
			// for the schedule add/edit drawer.
			if useCases.Entity.Client.ListClients != nil {
				priceScheduleDeps.ListClients = useCases.Entity.Client.ListClients
			}
			// Add location listing if available
			if useCases.Entity.Location.ListLocations != nil {
				priceScheduleDeps.ListLocations = useCases.Entity.Location.ListLocations
			}
			// Plans tab on the detail page lists price_plans filtered by price_schedule_id FK.
			if useCases.PricePlan.ListPricePlans != nil {
				priceScheduleDeps.ListPricePlans = useCases.PricePlan.ListPricePlans
				priceScheduleDeps.CreatePricePlan = useCases.PricePlan.CreatePricePlan
				priceScheduleDeps.ReadPricePlan = useCases.PricePlan.ReadPricePlan
				priceScheduleDeps.UpdatePricePlan = useCases.PricePlan.UpdatePricePlan
				priceScheduleDeps.DeletePricePlan = useCases.PricePlan.DeletePricePlan
			}
			// Reference checker for in-use guard (disables row Delete + locks pricing fields
			// on the edit drawer when a price_plan is referenced by active subscriptions).
			if w.refChecker != nil {
				priceScheduleDeps.GetPricePlanInUseIDs = w.refChecker.GetPricePlanInUseIDs
			}
			if useCases.Plan.ListPlans != nil {
				priceScheduleDeps.ListPlans = useCases.Plan.ListPlans
			}
			// Schedule-scoped plan detail (info + product-prices tabs) needs product lookups + ProductPricePlan CRUD
			if useCases.Product.ListProducts != nil {
				priceScheduleDeps.ListProducts = useCases.Product.ListProducts
			}
			if useCases.Product.ListProductPlans != nil {
				priceScheduleDeps.ListProductPlans = useCases.Product.ListProductPlans
			}
			if useCases.PricePlan.ListProductPricePlans != nil {
				priceScheduleDeps.ListProductPricePlans = useCases.PricePlan.ListProductPricePlans
				priceScheduleDeps.CreateProductPricePlan = useCases.PricePlan.CreateProductPricePlan
				priceScheduleDeps.UpdateProductPricePlan = useCases.PricePlan.UpdateProductPricePlan
				priceScheduleDeps.DeleteProductPricePlan = useCases.PricePlan.DeleteProductPricePlan
			}
			// 2026-05-04 — Engagements (subscriptions) tab on the schedule-scoped
			// price_plan detail page. See docs/plan/20260504-price-plan-engagements-tab/.
			if useCases.PriceSchedule.ListSubscriptionsByPricePlan != nil {
				priceScheduleDeps.ListSubscriptionsByPricePlan = useCases.PriceSchedule.ListSubscriptionsByPricePlan
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
					ListPriceSchedules:       useCases.PriceSchedule.ListPriceSchedules,
					ReadPriceSchedule:        useCases.PriceSchedule.ReadPriceSchedule,
					CreatePriceSchedule:      useCases.PriceSchedule.CreatePriceSchedule,
					UpdatePriceSchedule:      useCases.PriceSchedule.UpdatePriceSchedule,
					DeletePriceSchedule:      useCases.PriceSchedule.DeletePriceSchedule,
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
			ListPriceLists:     useCases.Product.ListPriceLists,
			ReadPriceList:      useCases.Product.ReadPriceList,
			CreatePriceList:    useCases.Product.CreatePriceList,
			UpdatePriceList:    useCases.Product.UpdatePriceList,
			DeletePriceList:    useCases.Product.DeletePriceList,
			ListPriceProducts:  useCases.Product.ListPriceProducts,
			CreatePriceProduct: useCases.Product.CreatePriceProduct,
			DeletePriceProduct: useCases.Product.DeletePriceProduct,
			ListProducts:       useCases.Product.ListProducts,
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
		if useCases.Entity.Client.ListClients != nil {
			lc := useCases.Entity.Client.ListClients
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
		if useCases.Operation.JobTemplate.ListJobTemplates != nil {
			ljt := useCases.Operation.JobTemplate.ListJobTemplates
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
		if useCases.Plan.ListPlans != nil {
			planListDeps.ListPlans = useCases.Plan.ListPlans
		}
		if w.refChecker != nil {
			planListDeps.GetInUseIDs = w.refChecker.GetPlanInUseIDs
		}
		ctx.Routes.GET(w.planRoutes.ListURL, planlist.NewView(planListDeps))
		ctx.Routes.GET(w.planRoutes.TableURL, planlist.NewTableView(planListDeps))

		// Plan CRUD actions
		if useCases.Plan.CreatePlan != nil {
			planActionDeps := &planaction.Deps{
				Routes:     w.planRoutes,
				Labels:     w.planLabels,
				CreatePlan: useCases.Plan.CreatePlan,
				ReadPlan:   useCases.Plan.ReadPlan,
				UpdatePlan: useCases.Plan.UpdatePlan,
				DeletePlan: useCases.Plan.DeletePlan,
				// SetPlanActive uses raw DB update (proto3 omits false booleans)
				SetPlanActive: func(fctx context.Context, id string, active bool) error {
					_, err := w.db.Update(fctx, "plan", id, map[string]any{"active": active})
					return err
				},
			}
			// 2026-04-27 plan-client-scope plan §6.2 — Client picker support
			// + reference-checker lock state.
			if useCases.Entity.Client.ListClients != nil {
				planActionDeps.ListClients = useCases.Entity.Client.ListClients
			}
			if useCases.Entity.Client.SearchClientsByName != nil {
				planActionDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName
			}
			if w.refChecker != nil {
				planActionDeps.GetPlanClientScopeLockedIDs = w.refChecker.GetPlanClientScopeLockedIDs
			}
			// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 —
			// JobTemplate select for Plan.job_template_id assignment.
			if useCases.Operation.JobTemplate.ListJobTemplates != nil {
				planActionDeps.ListJobTemplates = useCases.Operation.JobTemplate.ListJobTemplates
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
		if useCases.Plan.ReadPlan != nil {
			planDetailDeps := &plandetail.DetailViewDeps{
				Routes:                     w.planRoutes,
				PriceSchedulePlanDetailURL: w.priceScheduleRoutes.PlanDetailURL,
				ReadPlan:                   useCases.Plan.ReadPlan,
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
			if useCases.Product.ListProductPlans != nil {
				planDetailDeps.ListProductPlans = useCases.Product.ListProductPlans
			}
			if useCases.Product.ListProducts != nil {
				planDetailDeps.ListProducts = useCases.Product.ListProducts
			}
			if useCases.Product.ListProductVariants != nil {
				planDetailDeps.ListProductVariants = useCases.Product.ListProductVariants
			}
			if useCases.PricePlan.ListPricePlans != nil {
				planDetailDeps.ListPricePlans = useCases.PricePlan.ListPricePlans
			}
			if useCases.Entity.Location.ListLocations != nil {
				planDetailDeps.ListLocations = useCases.Entity.Location.ListLocations
			}
			if useCases.PriceSchedule.ListPriceSchedules != nil {
				planDetailDeps.ListPriceSchedules = useCases.PriceSchedule.ListPriceSchedules
			}
			// 2026-04-28 plan-client-scope — Info tab Client row needs to
			// resolve the plan's client_id label and (optionally) link to
			// the entydad client-detail page.
			if useCases.Entity.Client.ListClients != nil {
				planDetailDeps.ListClients = useCases.Entity.Client.ListClients
			}
			planDetailDeps.ClientDetailURL = cfg.clientDetailURL
			// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — Info
			// tab JobTemplate row resolution.
			if useCases.Operation.JobTemplate.ReadJobTemplate != nil {
				planDetailDeps.ReadJobTemplate = useCases.Operation.JobTemplate.ReadJobTemplate
			}
			ctx.Routes.GET(w.planRoutes.DetailURL, plandetail.NewView(planDetailDeps))
			ctx.Routes.GET(w.planRoutes.TabActionURL, plandetail.NewTabAction(planDetailDeps))

			// Plan-scoped PricePlan detail (/app/plans/detail/{id}/price/{ppid}).
			// Reuses the schedule-scoped detail body but anchors ActiveNav to
			// Services > Packages and points the breadcrumb back at the
			// package's package-prices tab. The {id} path value is plan_id;
			// the handler resolves schedule_id from the price_plan record.
			if w.planRoutes.PricePlanDetailURL != "" && useCases.PricePlan.ReadPricePlan != nil {
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
					ReadPricePlan:          useCases.PricePlan.ReadPricePlan,
					// Mount overrides — keep the page anchored to Packages.
					ActiveNavOverride:      w.planRoutes.ActiveNav,
					ActiveSubNavOverride:   w.planRoutes.ActiveSubNav,
					PlanDetailBackURL:      w.planRoutes.DetailURL,
					PlanDetailBackTab:      packagePricesSlug,
					PlanScopedDetailURL:    w.planRoutes.PricePlanDetailURL,
					PlanScopedTabActionURL: w.planRoutes.PricePlanTabActionURL,
				}
				if useCases.PriceSchedule.ReadPriceSchedule != nil {
					planScopedDeps.ReadPriceSchedule = useCases.PriceSchedule.ReadPriceSchedule
				}
				if useCases.Plan.ListPlans != nil {
					planScopedDeps.ListPlans = useCases.Plan.ListPlans
				}
				if useCases.Product.ListProducts != nil {
					planScopedDeps.ListProducts = useCases.Product.ListProducts
				}
				if useCases.Product.ListProductPlans != nil {
					planScopedDeps.ListProductPlans = useCases.Product.ListProductPlans
				}
				if useCases.Product.ListProductVariants != nil {
					planScopedDeps.ListProductVariants = useCases.Product.ListProductVariants
				}
				if useCases.PricePlan.ListProductPricePlans != nil {
					planScopedDeps.ListProductPricePlans = useCases.PricePlan.ListProductPricePlans
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
			if useCases.PricePlan.CreatePricePlan != nil {
				ppActionDeps := &planaction.PricePlanDeps{
					Routes:              w.planRoutes,
					Labels:              w.planLabels,
					PricePlanLabels:     w.pricePlanLabels,
					PriceScheduleLabels: w.priceScheduleLabels,
					CommonLabels:        ctx.Common,
					CreatePricePlan: useCases.PricePlan.CreatePricePlan,
					ReadPricePlan:   useCases.PricePlan.ReadPricePlan,
					UpdatePricePlan: useCases.PricePlan.UpdatePricePlan,
					DeletePricePlan: useCases.PricePlan.DeletePricePlan,
				}
				if useCases.PriceSchedule.ListPriceSchedules != nil {
					ppActionDeps.ListPriceSchedules = useCases.PriceSchedule.ListPriceSchedules
				}
				if useCases.Plan.ReadPlan != nil {
					ppActionDeps.ReadPlan = useCases.Plan.ReadPlan
				}
				// Plan §6.7 — ListClients powers the readonly schedule
				// label + lock tooltip when the parent Plan is client-scoped.
				if useCases.Entity.Client.ListClients != nil {
					ppActionDeps.ListClients = useCases.Entity.Client.ListClients
				}
				if useCases.Entity.Location.ListLocations != nil {
					ppActionDeps.ListLocations = useCases.Entity.Location.ListLocations
				}
				if w.refChecker != nil {
					ppActionDeps.GetPricePlanInUseIDs = w.refChecker.GetPricePlanInUseIDs
				}
				if useCases.Product.ListProducts != nil {
					ppActionDeps.ListProducts = useCases.Product.ListProducts
				}
				if useCases.Product.ListProductPlans != nil {
					ppActionDeps.ListProductPlans = useCases.Product.ListProductPlans
				}
				if useCases.PricePlan.ListProductPricePlans != nil {
					ppActionDeps.CreateProductPricePlan = useCases.PricePlan.CreateProductPricePlan
					ppActionDeps.ListProductPricePlans = useCases.PricePlan.ListProductPricePlans
				}
				ctx.Routes.GET(w.planRoutes.PricePlanAddURL, planaction.NewPricePlanAddAction(ppActionDeps))
				ctx.Routes.POST(w.planRoutes.PricePlanAddURL, planaction.NewPricePlanAddAction(ppActionDeps))
				ctx.Routes.GET(w.planRoutes.PricePlanEditURL, planaction.NewPricePlanEditAction(ppActionDeps))
				ctx.Routes.POST(w.planRoutes.PricePlanEditURL, planaction.NewPricePlanEditAction(ppActionDeps))
				ctx.Routes.POST(w.planRoutes.PricePlanDeleteURL, planaction.NewPricePlanDeleteAction(ppActionDeps))
			}
			// ProductPlan CRUD within plan detail
			if useCases.Product.CreateProductPlan != nil {
				productPlanActionDeps := &planaction.ProductPlanDeps{
					Routes:            w.planRoutes,
					Labels:            w.planLabels,
					CreateProductPlan: useCases.Product.CreateProductPlan,
					ReadProductPlan:   useCases.Product.ReadProductPlan,
					UpdateProductPlan: useCases.Product.UpdateProductPlan,
					DeleteProductPlan: useCases.Product.DeleteProductPlan,
				}
				if useCases.Product.ListProducts != nil {
					productPlanActionDeps.ListProducts = useCases.Product.ListProducts
				}
				if useCases.Product.ListProductPlans != nil {
					productPlanActionDeps.ListProductPlans = useCases.Product.ListProductPlans
				}
				if useCases.Product.ListProductVariants != nil {
					productPlanActionDeps.ListProductVariants = useCases.Product.ListProductVariants
				}
				if useCases.Product.ListProductVariantOptions != nil {
					productPlanActionDeps.ListProductVariantOptions = useCases.Product.ListProductVariantOptions
				}
				if useCases.Product.ListProductOptionValues != nil {
					productPlanActionDeps.ListProductOptionValues = useCases.Product.ListProductOptionValues
				}
				if useCases.Product.ListProductOptions != nil {
					productPlanActionDeps.ListProductOptions = useCases.Product.ListProductOptions
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
			if useCases.Plan.ListPlans != nil {
				planBundleListDeps.ListPlans = useCases.Plan.ListPlans
			}
			if w.refChecker != nil {
				planBundleListDeps.GetInUseIDs = w.refChecker.GetPlanInUseIDs
			}
			ctx.Routes.GET(w.planBundleRoutes.ListURL, planlist.NewView(planBundleListDeps))
			ctx.Routes.GET(w.planBundleRoutes.TableURL, planlist.NewTableView(planBundleListDeps))

			if useCases.Plan.CreatePlan != nil {
				planBundleActionDeps := &planaction.Deps{
					Routes:     w.planBundleRoutes,
					Labels:     w.planLabels,
					CreatePlan: useCases.Plan.CreatePlan,
					ReadPlan:   useCases.Plan.ReadPlan,
					UpdatePlan: useCases.Plan.UpdatePlan,
					DeletePlan: useCases.Plan.DeletePlan,
					// SetPlanActive uses raw DB update (proto3 omits false booleans)
					SetPlanActive: func(fctx context.Context, id string, active bool) error {
						_, err := w.db.Update(fctx, "plan", id, map[string]any{"active": active})
						return err
					},
				}
				// 2026-04-27 plan-client-scope plan §6.2 — same Client picker
				// + lock state on the bundle mount.
				if useCases.Entity.Client.ListClients != nil {
					planBundleActionDeps.ListClients = useCases.Entity.Client.ListClients
				}
				if useCases.Entity.Client.SearchClientsByName != nil {
					planBundleActionDeps.SearchClientsByName = useCases.Entity.Client.SearchClientsByName
				}
				if w.refChecker != nil {
					planBundleActionDeps.GetPlanClientScopeLockedIDs = w.refChecker.GetPlanClientScopeLockedIDs
				}
				// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 —
				// JobTemplate select on the bundle mount.
				if useCases.Operation.JobTemplate.ListJobTemplates != nil {
					planBundleActionDeps.ListJobTemplates = useCases.Operation.JobTemplate.ListJobTemplates
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

			if useCases.Plan.ReadPlan != nil {
				planBundleDetailDeps := &plandetail.DetailViewDeps{
					Routes:                     w.planBundleRoutes,
					PriceSchedulePlanDetailURL: w.priceScheduleRoutes.PlanDetailURL,
					ReadPlan:                   useCases.Plan.ReadPlan,
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
				if useCases.Product.ListProductPlans != nil {
					planBundleDetailDeps.ListProductPlans = useCases.Product.ListProductPlans
				}
				if useCases.Product.ListProducts != nil {
					planBundleDetailDeps.ListProducts = useCases.Product.ListProducts
				}
				if useCases.Product.ListProductVariants != nil {
					planBundleDetailDeps.ListProductVariants = useCases.Product.ListProductVariants
				}
				if useCases.PricePlan.ListPricePlans != nil {
					planBundleDetailDeps.ListPricePlans = useCases.PricePlan.ListPricePlans
				}
				if useCases.Entity.Location.ListLocations != nil {
					planBundleDetailDeps.ListLocations = useCases.Entity.Location.ListLocations
				}
				if useCases.PriceSchedule.ListPriceSchedules != nil {
					planBundleDetailDeps.ListPriceSchedules = useCases.PriceSchedule.ListPriceSchedules
				}
				// 2026-04-28 plan-client-scope — same Info tab Client row
				// wiring on the bundle mount.
				if useCases.Entity.Client.ListClients != nil {
					planBundleDetailDeps.ListClients = useCases.Entity.Client.ListClients
				}
				planBundleDetailDeps.ClientDetailURL = cfg.clientDetailURL
				// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 —
				// Info tab JobTemplate row on the bundle mount.
				if useCases.Operation.JobTemplate.ReadJobTemplate != nil {
					planBundleDetailDeps.ReadJobTemplate = useCases.Operation.JobTemplate.ReadJobTemplate
				}
				ctx.Routes.GET(w.planBundleRoutes.DetailURL, plandetail.NewView(planBundleDetailDeps))
				ctx.Routes.GET(w.planBundleRoutes.TabActionURL, plandetail.NewTabAction(planBundleDetailDeps))
				if w.uploadFile != nil {
					ctx.Routes.GET(w.planBundleRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planBundleDetailDeps))
					ctx.Routes.POST(w.planBundleRoutes.AttachmentUploadURL, plandetail.NewAttachmentUploadAction(planBundleDetailDeps))
					ctx.Routes.POST(w.planBundleRoutes.AttachmentDeleteURL, plandetail.NewAttachmentDeleteAction(planBundleDetailDeps))
				}
				if useCases.PricePlan.CreatePricePlan != nil {
					ppBundleDeps := &planaction.PricePlanDeps{
						Routes:              w.planBundleRoutes,
						Labels:              w.planLabels,
						PricePlanLabels:     w.pricePlanLabels,
						PriceScheduleLabels: w.priceScheduleLabels,
						CommonLabels:        ctx.Common,
						CreatePricePlan: useCases.PricePlan.CreatePricePlan,
						ReadPricePlan:   useCases.PricePlan.ReadPricePlan,
						UpdatePricePlan: useCases.PricePlan.UpdatePricePlan,
						DeletePricePlan: useCases.PricePlan.DeletePricePlan,
					}
					if useCases.PriceSchedule.ListPriceSchedules != nil {
						ppBundleDeps.ListPriceSchedules = useCases.PriceSchedule.ListPriceSchedules
					}
					if useCases.Plan.ReadPlan != nil {
						ppBundleDeps.ReadPlan = useCases.Plan.ReadPlan
					}
					// Plan §6.7 — ListClients powers the readonly schedule
					// label + lock tooltip on the bundle-mount drawer.
					if useCases.Entity.Client.ListClients != nil {
						ppBundleDeps.ListClients = useCases.Entity.Client.ListClients
					}
					if useCases.Entity.Location.ListLocations != nil {
						ppBundleDeps.ListLocations = useCases.Entity.Location.ListLocations
					}
					if w.refChecker != nil {
						ppBundleDeps.GetPricePlanInUseIDs = w.refChecker.GetPricePlanInUseIDs
					}
					if useCases.Product.ListProducts != nil {
						ppBundleDeps.ListProducts = useCases.Product.ListProducts
					}
					if useCases.Product.ListProductPlans != nil {
						ppBundleDeps.ListProductPlans = useCases.Product.ListProductPlans
					}
					if useCases.PricePlan.ListProductPricePlans != nil {
						ppBundleDeps.CreateProductPricePlan = useCases.PricePlan.CreateProductPricePlan
						ppBundleDeps.ListProductPricePlans = useCases.PricePlan.ListProductPricePlans
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
				if useCases.Product.CreateProductPlan != nil {
					ppBundleProductPlanDeps := &planaction.ProductPlanDeps{
						Routes:            w.planBundleRoutes,
						Labels:            w.planLabels,
						CreateProductPlan: useCases.Product.CreateProductPlan,
						ReadProductPlan:   useCases.Product.ReadProductPlan,
						UpdateProductPlan: useCases.Product.UpdateProductPlan,
						DeleteProductPlan: useCases.Product.DeleteProductPlan,
					}
					if useCases.Product.ListProducts != nil {
						ppBundleProductPlanDeps.ListProducts = useCases.Product.ListProducts
					}
					if useCases.Product.ListProductPlans != nil {
						ppBundleProductPlanDeps.ListProductPlans = useCases.Product.ListProductPlans
					}
					if useCases.Product.ListProductVariants != nil {
						ppBundleProductPlanDeps.ListProductVariants = useCases.Product.ListProductVariants
					}
					if useCases.Product.ListProductVariantOptions != nil {
						ppBundleProductPlanDeps.ListProductVariantOptions = useCases.Product.ListProductVariantOptions
					}
					if useCases.Product.ListProductOptionValues != nil {
						ppBundleProductPlanDeps.ListProductOptionValues = useCases.Product.ListProductOptionValues
					}
					if useCases.Product.ListProductOptions != nil {
						ppBundleProductPlanDeps.ListProductOptions = useCases.Product.ListProductOptions
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
