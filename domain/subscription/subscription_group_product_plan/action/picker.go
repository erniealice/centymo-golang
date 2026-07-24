package action

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sort"
	"strings"

	sgpp "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"
	jobtemplaterelationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_relation"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
	sgpppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan"
)

// sectionStaffTabPanelID is the S1 section-tab panel id (subscription_group/
// detail/staff.go's staffTabID, kept a literal string here — S1's edits are a
// separate workstream, plan.md package tree §1 — so success responses that
// mount under S1 later refresh correctly without a cross-package import).
const sectionStaffTabPanelID = "subscription-group-staff-tab"

// PickerDeps holds all dependencies for the S2 add-offerings picker
// (centymo.md §4). Every read is a batched LIST_IN call (espyna.md §1b); the
// template resolver is the grounded chain from plan.md §3: plan.job_template_id
// (root) -> job_template_relation (parent=root, active) -> child job_templates
// filtered output_product_id == offering.product_id -> MIN(sequence_order),
// all candidates surfaced.
type PickerDeps struct {
	Routes       sgpp.Routes
	Labels       sgpp.Labels
	CommonLabels pyeza.CommonLabels
	Options      sgpp.Options

	ReadSubscriptionGroup              func(ctx context.Context, req *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error)
	ListPlans                          func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules                 func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	ListProductPlans                   func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProductVariants                func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	ListSubscriptionGroupProductPlans  func(ctx context.Context, req *sgpppb.ListSubscriptionGroupProductPlansRequest) (*sgpppb.ListSubscriptionGroupProductPlansResponse, error)
	ListJobTemplateRelations           func(ctx context.Context, req *jobtemplaterelationpb.ListJobTemplateRelationsRequest) (*jobtemplaterelationpb.ListJobTemplateRelationsResponse, error)
	ListJobTemplates                   func(ctx context.Context, req *jobtemplatepb.ListJobTemplatesRequest) (*jobtemplatepb.ListJobTemplatesResponse, error)
	CreateSubscriptionGroupProductPlan func(ctx context.Context, req *sgpppb.CreateSubscriptionGroupProductPlanRequest) (*sgpppb.CreateSubscriptionGroupProductPlanResponse, error)
}

// templateCandidate is one job_template_relation child surviving the
// output_product_id filter (plan.md §3 "ALL candidates logged").
type templateCandidate struct {
	ID    string
	Label string
	Order int32
}

// pickerRow is one candidate offering row (an active plan offering with no
// active class row on this section yet).
type pickerRow struct {
	ProductPlanID      string
	Label              string
	VariantLabel       string
	CurriculumLabel    string // resolved (MIN sequence_order) template name, "" when unresolved
	TemplateCandidates []templateCandidate
	Ambiguous          bool
	Checked            bool
}

// PickerData is the template data for the S2 picker drawer.
type PickerData struct {
	FormAction       string
	SectionID        string
	SectionName      string
	PlanName         string
	PeriodName       string
	Rows             []pickerRow
	SelectAllChecked bool
	Labels           sgpp.Labels
	CommonLabels     pyeza.CommonLabels
	WorkspaceID      string
	Nonce            string
}

// pickerContext bundles the section-scoped lookups shared by GET and POST.
type pickerContext struct {
	sectionName     string
	planName        string
	periodName      string
	rootTemplateID  string
	offeringsByID   map[string]*productplanpb.ProductPlan
	existingPlanIDs map[string]bool
}

func loadPickerContext(ctx context.Context, deps *PickerDeps, sectionID string) (*pickerContext, error) {
	sgResp, err := deps.ReadSubscriptionGroup(ctx, &subscriptiongrouppb.ReadSubscriptionGroupRequest{
		Data: &subscriptiongrouppb.SubscriptionGroup{Id: sectionID},
	})
	if err != nil || sgResp == nil || len(sgResp.GetData()) == 0 {
		return nil, errors.New(deps.Labels.Errors.NotFound)
	}
	sg := sgResp.GetData()[0]
	planID := sg.GetPlanId()

	out := &pickerContext{
		sectionName:     sg.GetName(),
		offeringsByID:   map[string]*productplanpb.ProductPlan{},
		existingPlanIDs: map[string]bool{},
	}

	if planID != "" && deps.ListPlans != nil {
		if resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{Filters: sgpp.StringEqFilter("id", planID)}); err == nil && resp != nil {
			for _, p := range resp.GetData() {
				if p.GetId() == planID {
					out.planName = p.GetName()
					out.rootTemplateID = p.GetJobTemplateId()
				}
			}
		}
	}
	if sid := sg.GetPriceScheduleId(); sid != "" && deps.ListPriceSchedules != nil {
		if resp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{Filters: sgpp.StringEqFilter("id", sid)}); err == nil && resp != nil {
			for _, s := range resp.GetData() {
				if s.GetId() == sid {
					out.periodName = s.GetName()
				}
			}
		}
	}
	if planID != "" && deps.ListProductPlans != nil {
		if resp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
			Filters: sgpp.AndFilters(sgpp.StringEqFilter("plan_id", planID), sgpp.BoolEqFilter("active", true)),
		}); err == nil && resp != nil {
			for _, pp := range resp.GetData() {
				if pp != nil {
					out.offeringsByID[pp.GetId()] = pp
				}
			}
		}
	}
	if deps.ListSubscriptionGroupProductPlans != nil {
		if resp, err := deps.ListSubscriptionGroupProductPlans(ctx, &sgpppb.ListSubscriptionGroupProductPlansRequest{
			Filters: sgpp.AndFilters(sgpp.StringEqFilter("subscription_group_id", sectionID), sgpp.BoolEqFilter("active", true)),
		}); err == nil && resp != nil {
			for _, c := range resp.GetData() {
				if c != nil && c.GetSubscriptionGroupId() == sectionID {
					out.existingPlanIDs[c.GetProductPlanId()] = true
				}
			}
		}
	}
	return out, nil
}

// resolveTemplateCandidates is the M0 resolver (plan.md §3, verbatim): root ->
// active job_template_relation children -> job_template filtered by
// output_product_id -> MIN(sequence_order) tie-break, all candidates
// returned for the picker's constrained select.
func resolveTemplateCandidates(ctx context.Context, deps *PickerDeps, rootTemplateID, productID string) (candidates []templateCandidate, chosenID string) {
	if rootTemplateID == "" || productID == "" || deps.ListJobTemplateRelations == nil || deps.ListJobTemplates == nil {
		return nil, ""
	}
	relResp, err := deps.ListJobTemplateRelations(ctx, &jobtemplaterelationpb.ListJobTemplateRelationsRequest{
		Filters: sgpp.AndFilters(sgpp.StringEqFilter("parent_template_id", rootTemplateID), sgpp.BoolEqFilter("active", true)),
	})
	if err != nil || relResp == nil {
		return nil, ""
	}
	seqByChild := map[string]int32{}
	childIDs := make([]string, 0, len(relResp.GetData()))
	for _, rel := range relResp.GetData() {
		if rel == nil || !rel.GetActive() || rel.GetParentTemplateId() != rootTemplateID {
			continue
		}
		cid := rel.GetChildTemplateId()
		if cid == "" {
			continue
		}
		childIDs = append(childIDs, cid)
		seqByChild[cid] = rel.GetSequenceOrder()
	}
	if len(childIDs) == 0 {
		return nil, ""
	}
	tmplResp, err := deps.ListJobTemplates(ctx, &jobtemplatepb.ListJobTemplatesRequest{
		Filters: sgpp.AndFilters(sgpp.StringInFilter("id", childIDs), sgpp.BoolEqFilter("active", true), sgpp.StringEqFilter("output_product_id", productID)),
	})
	if err != nil || tmplResp == nil {
		return nil, ""
	}
	for _, t := range tmplResp.GetData() {
		if t == nil || t.GetOutputProductId() != productID {
			continue
		}
		candidates = append(candidates, templateCandidate{ID: t.GetId(), Label: t.GetName(), Order: seqByChild[t.GetId()]})
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Order != candidates[j].Order {
			return candidates[i].Order < candidates[j].Order
		}
		return candidates[i].ID < candidates[j].ID
	})
	if len(candidates) > 0 {
		chosenID = candidates[0].ID
	}
	return candidates, chosenID
}

func offeringLabel(pp *productplanpb.ProductPlan) string {
	if pp == nil {
		return ""
	}
	if n := strings.TrimSpace(pp.GetName()); n != "" {
		return n
	}
	return pp.GetId()
}

func buildPickerData(ctx context.Context, deps *PickerDeps, sectionID string) (*PickerData, error) {
	pc, err := loadPickerContext(ctx, deps, sectionID)
	if err != nil {
		return nil, err
	}

	variantIDs := make([]string, 0)
	for _, pp := range pc.offeringsByID {
		if v := pp.GetProductVariantId(); v != "" {
			variantIDs = append(variantIDs, v)
		}
	}
	variantLabels := map[string]string{}
	if len(variantIDs) > 0 && deps.ListProductVariants != nil {
		if resp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{Filters: sgpp.StringInFilter("id", variantIDs)}); err == nil && resp != nil {
			for _, v := range resp.GetData() {
				label := v.GetSku()
				if label == "" {
					label = v.GetId()
				}
				variantLabels[v.GetId()] = label
			}
		}
	}

	selectAllDefault := deps.Options.Picker.SelectAll()
	rows := make([]pickerRow, 0, len(pc.offeringsByID))
	for id, pp := range pc.offeringsByID {
		if pc.existingPlanIDs[id] {
			continue
		}
		candidates, chosen := resolveTemplateCandidates(ctx, deps, pc.rootTemplateID, pp.GetProductId())
		curriculumLabel := ""
		if chosen != "" {
			for _, c := range candidates {
				if c.ID == chosen {
					curriculumLabel = c.Label
					break
				}
			}
		}
		rows = append(rows, pickerRow{
			ProductPlanID:      id,
			Label:              offeringLabel(pp),
			VariantLabel:       variantLabels[pp.GetProductVariantId()],
			CurriculumLabel:    curriculumLabel,
			TemplateCandidates: candidates,
			Ambiguous:          len(candidates) > 1,
			Checked:            selectAllDefault,
		})
	}
	sort.SliceStable(rows, func(i, j int) bool { return strings.ToLower(rows[i].Label) < strings.ToLower(rows[j].Label) })

	return &PickerData{
		FormAction:       route.ResolveURL(deps.Routes.PickerURL, "id", sectionID),
		SectionID:        sectionID,
		SectionName:      pc.sectionName,
		PlanName:         pc.planName,
		PeriodName:       pc.periodName,
		Rows:             rows,
		SelectAllChecked: selectAllDefault,
		Labels:           deps.Labels,
		CommonLabels:     deps.CommonLabels,
	}, nil
}

// createFromPicker re-derives candidates server-side (never trusts the body
// for which ids are addable — plan.md §5 IDOR guard) and creates one ACTIVE,
// unstaffed class row per checked, still-valid candidate (centymo.md §4).
// Non-candidates are skipped + logged, never erroring the whole batch.
func createFromPicker(ctx context.Context, deps *PickerDeps, sectionID string, r *http.Request) (int, error) {
	pc, err := loadPickerContext(ctx, deps, sectionID)
	if err != nil {
		return 0, err
	}
	checked := r.Form["product_plan_id"]
	created := 0
	for _, ppID := range checked {
		ppID = strings.TrimSpace(ppID)
		if ppID == "" {
			continue
		}
		if pc.existingPlanIDs[ppID] {
			log.Printf("picker: skipping already-added offering %s on section %s", ppID, sectionID)
			continue
		}
		pp, ok := pc.offeringsByID[ppID]
		if !ok {
			log.Printf("picker: skipping non-candidate offering %s on section %s", ppID, sectionID)
			continue
		}
		candidates, autoChosen := resolveTemplateCandidates(ctx, deps, pc.rootTemplateID, pp.GetProductId())
		templateID := autoChosen
		if picked := strings.TrimSpace(r.FormValue("template_id_" + ppID)); picked != "" {
			for _, c := range candidates {
				if c.ID == picked {
					templateID = picked
					break
				}
			}
		}
		_, err := deps.CreateSubscriptionGroupProductPlan(ctx, &sgpppb.CreateSubscriptionGroupProductPlanRequest{
			Data: &sgpppb.SubscriptionGroupProductPlan{
				SubscriptionGroupId: sectionID,
				ProductPlanId:       ppID,
				JobTemplateId:       templateID,
				Status:              sgpppb.SubscriptionGroupProductPlanStatus_SUBSCRIPTION_GROUP_PRODUCT_PLAN_STATUS_ACTIVE,
			},
		})
		if err != nil {
			log.Printf("picker: failed to add offering %s to section %s: %v", ppID, sectionID, err)
			continue
		}
		created++
	}
	if created == 0 {
		return 0, errors.New(deps.Labels.Errors.CreateFailed)
	}
	return created, nil
}

// NewPickerAction handles GET (candidates + resolved template per row) / POST
// (create batch) for the S2 add-offerings picker.
func NewPickerAction(deps *PickerDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(sgppEntity, "create") {
			return view.HTMXError(deps.Labels.Errors.Unauthorized)
		}
		sectionID := viewCtx.Request.PathValue("id")
		if sectionID == "" {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		if viewCtx.Request.Method == http.MethodGet {
			data, err := buildPickerData(ctx, deps, sectionID)
			if err != nil {
				return view.HTMXError(err.Error())
			}
			return view.OK("sgpp-picker-drawer", data)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.CreateFailed)
		}
		if _, err := createFromPicker(ctx, deps, sectionID, viewCtx.Request); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess(sectionStaffTabPanelID)
	})
}
