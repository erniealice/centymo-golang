package detail

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	sgpp "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"
	jobtemplatephasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_phase"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
	sgpppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan"
	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

// sgppEntity / sgppsEntity mirror action's entity constants (kept
// package-local — detail/ and action/ are separate packages).
const (
	sgppEntity  = "subscription_group_product_plan"
	sgppsEntity = "subscription_group_product_plan_staff"
)

// DetailViewDeps holds view dependencies for the S4 class detail page (Info +
// Teachers tabs, centymo.md §6/§7). Every read is a single batched call
// scoped to this one class (espyna.md §1b); RouteMap cross-links
// (SectionDetailURL, GradeSheetURL) are nil-safe closures the container
// resolves — never a hardcoded path (plan.md §1.1b).
type DetailViewDeps struct {
	Routes       sgpp.Routes
	Labels       sgpp.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadSubscriptionGroupProductPlan       func(ctx context.Context, req *sgpppb.ReadSubscriptionGroupProductPlanRequest) (*sgpppb.ReadSubscriptionGroupProductPlanResponse, error)
	ReadSubscriptionGroup                  func(ctx context.Context, req *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error)
	ListPlans                              func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules                     func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	ListProductPlans                       func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProductVariants                    func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	ListJobTemplates                       func(ctx context.Context, req *jobtemplatepb.ListJobTemplatesRequest) (*jobtemplatepb.ListJobTemplatesResponse, error)
	ListJobTemplatePhases                  func(ctx context.Context, req *jobtemplatephasepb.ListJobTemplatePhasesRequest) (*jobtemplatephasepb.ListJobTemplatePhasesResponse, error)
	ListProductPlanStaffs                  func(ctx context.Context, req *productplanstaffpb.ListProductPlanStaffsRequest) (*productplanstaffpb.ListProductPlanStaffsResponse, error)
	ListStaffNames                         func(ctx context.Context) map[string]string
	ListSubscriptionGroupProductPlanStaffs func(ctx context.Context, req *sgppspb.ListSubscriptionGroupProductPlanStaffsRequest) (*sgppspb.ListSubscriptionGroupProductPlanStaffsResponse, error)

	// GetSubscriptionGroupProductPlanInUseIDs gates the Info tab's Remove
	// action (plan.md §2 in-use guard; server-side guard is the hard
	// backstop, this is the UI hint).
	GetSubscriptionGroupProductPlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
	// SectionDetailURL resolves the "→ link back" to the section page.
	SectionDetailURL func(ctx context.Context, sectionID string) string
	// GradeSheetURL resolves fayna's grade-sheet route for this class's
	// curriculum (RouteMap cross-link — plan.md §6.1 "no new route").
	GradeSheetURL func(ctx context.Context, jobTemplateID string) string
}

// InfoAction is one Info-tab guarded button (Exclude/Restore/Remove — S6).
// Action must be a value table-actions.js recognizes ("deactivate" for
// Exclude, "activate" for Restore, "delete" for Remove) so the SAME generic
// confirm+POST behavior every list row uses fires here too (table-actions.js
// delegates from `document`, not a table root).
type InfoAction struct {
	Action          string
	Label           string
	URL             string
	ItemName        string
	ConfirmTitle    string
	ConfirmMessage  string
	Disabled        bool
	DisabledTooltip string
	TestID          string
}

// PageData holds the data for the S4 class detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          sgpp.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID                string
	SectionID         string
	SectionName       string
	SectionURL        string
	OfferingName      string
	OfferingVariantID string // raw product_variant_id — the §2.5 coverage-rule key, never displayed
	VariantLabel      string // resolved display name ("" when the offering is an umbrella)
	CurriculumName    string
	PlanName          string
	PeriodName        string
	Status            string
	StatusLabel       string
	StatusVariant     string
	AssignmentsCount  int
	CreatedDate       string
	ModifiedDate      string
	GradeSheetURL     string

	EditURL     string
	InfoActions []InfoAction

	// Teachers tab (S5)
	Assignments      *types.TableConfig
	AssignmentsEmpty bool
	AddURL           string
}

// NewView creates the S4 class detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(sgppEntity, "read") {
			return view.Forbidden(sgppEntity + ":read")
		}
		sgppID := viewCtx.Request.PathValue("sgppId")
		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}
		pageData, err := buildPageData(ctx, deps, sgppID, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}
		// IDOR guard: the path's {id} (section) must match the class's own
		// section — a forged/mismatched section id in the URL 404s rather than
		// silently rendering a foreign class (plan.md §5).
		if pathSectionID := viewCtx.Request.PathValue("id"); pathSectionID != "" && pathSectionID != pageData.SectionID {
			return view.Error(fmt.Errorf("%s", deps.Labels.Errors.NotFound))
		}
		return view.OK("sgpp-detail", pageData)
	})
}

// NewTabAction handles GET /action/subscription-group-product-plan/{sgppId}/
// tab/{tab}.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		sgppID := viewCtx.Request.PathValue("sgppId")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}
		pageData, err := buildPageData(ctx, deps, sgppID, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("sgpp-tab-"+tab, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, sgppID, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadSubscriptionGroupProductPlan(ctx, &sgpppb.ReadSubscriptionGroupProductPlanRequest{
		Data: &sgpppb.SubscriptionGroupProductPlan{Id: sgppID},
	})
	if err != nil {
		log.Printf("Failed to read subscription_group_product_plan %s: %v", sgppID, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	record := data[0]
	l := deps.Labels

	sectionID := record.GetSubscriptionGroupId()
	sectionName := sectionID
	var planID, priceScheduleID string
	if deps.ReadSubscriptionGroup != nil && sectionID != "" {
		if sgResp, err := deps.ReadSubscriptionGroup(ctx, &subscriptiongrouppb.ReadSubscriptionGroupRequest{
			Data: &subscriptiongrouppb.SubscriptionGroup{Id: sectionID},
		}); err == nil && sgResp != nil && len(sgResp.GetData()) > 0 {
			sg := sgResp.GetData()[0]
			sectionName = sg.GetName()
			planID = sg.GetPlanId()
			priceScheduleID = sg.GetPriceScheduleId()
		}
	}

	offeringName := record.GetProductPlanId()
	offeringVariantID := ""
	if deps.ListProductPlans != nil && record.GetProductPlanId() != "" {
		if ppResp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
			Filters: sgpp.StringEqFilter("id", record.GetProductPlanId()),
		}); err == nil && ppResp != nil {
			for _, pp := range ppResp.GetData() {
				if pp.GetId() == record.GetProductPlanId() {
					if n := pp.GetName(); n != "" {
						offeringName = n
					}
					offeringVariantID = pp.GetProductVariantId()
				}
			}
		}
	}
	variantLabel := ""
	if offeringVariantID != "" && deps.ListProductVariants != nil {
		if vResp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{
			Filters: sgpp.StringEqFilter("id", offeringVariantID),
		}); err == nil && vResp != nil {
			for _, v := range vResp.GetData() {
				if v.GetId() == offeringVariantID {
					variantLabel = v.GetSku()
				}
			}
		}
		if variantLabel == "" {
			variantLabel = offeringVariantID
		}
	}

	curriculumName := record.GetJobTemplateId()
	if curriculumName == "" {
		curriculumName = l.Detail.NoTemplate
	} else if deps.ListJobTemplates != nil {
		if tResp, err := deps.ListJobTemplates(ctx, &jobtemplatepb.ListJobTemplatesRequest{
			Filters: sgpp.StringEqFilter("id", record.GetJobTemplateId()),
		}); err == nil && tResp != nil {
			for _, t := range tResp.GetData() {
				if t.GetId() == record.GetJobTemplateId() && t.GetName() != "" {
					curriculumName = t.GetName()
				}
			}
		}
	}

	planName := l.Detail.NoPlan
	if planID != "" && deps.ListPlans != nil {
		if pResp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{Filters: sgpp.StringEqFilter("id", planID)}); err == nil && pResp != nil {
			for _, p := range pResp.GetData() {
				if p.GetId() == planID && p.GetName() != "" {
					planName = p.GetName()
				}
			}
		}
	}
	periodName := ""
	if priceScheduleID != "" && deps.ListPriceSchedules != nil {
		if psResp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{Filters: sgpp.StringEqFilter("id", priceScheduleID)}); err == nil && psResp != nil {
			for _, ps := range psResp.GetData() {
				if ps.GetId() == priceScheduleID && ps.GetName() != "" {
					periodName = ps.GetName()
				}
			}
		}
	}

	status := "active"
	statusLabel := l.Form.StatusActive
	statusVariant := "success"
	if record.GetStatus() == sgpppb.SubscriptionGroupProductPlanStatus_SUBSCRIPTION_GROUP_PRODUCT_PLAN_STATUS_EXCLUDED {
		status = "excluded"
		statusLabel = l.Form.StatusExcluded
		statusVariant = "warning"
	}

	assignmentsCount := countAssignments(ctx, deps, sgppID)

	inUse := false
	if deps.GetSubscriptionGroupProductPlanInUseIDs != nil {
		if m, _ := deps.GetSubscriptionGroupProductPlanInUseIDs(ctx, []string{sgppID}); m != nil {
			inUse = m[sgppID]
		}
	}

	gradeSheetURL := ""
	if deps.GradeSheetURL != nil && record.GetJobTemplateId() != "" {
		gradeSheetURL = deps.GradeSheetURL(ctx, record.GetJobTemplateId())
	}
	sectionURL := ""
	if deps.SectionDetailURL != nil {
		sectionURL = deps.SectionDetailURL(ctx, sectionID)
	}

	tz := types.LocationFromContext(ctx)
	createdDate := ""
	if ms := record.GetDateCreated(); ms > 0 {
		createdDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}
	modifiedDate := ""
	if ms := record.GetDateModified(); ms > 0 {
		modifiedDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}

	base := route.ResolveURL(deps.Routes.DetailURL, "id", sectionID, "sgppId", sgppID)
	action := route.ResolveURL(deps.Routes.TabActionURL, "sgppId", sgppID, "tab", "")
	tabItems := []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "staff", Label: l.Tabs.Staff, Href: base + "?tab=staff", HxGet: action + "staff", Icon: "icon-user-check", Count: assignmentsCount},
	}

	perms := view.GetUserPermissions(ctx)
	editURL := route.ResolveURL(deps.Routes.EditURL, "id", sgppID)
	infoActions := buildInfoActions(deps, sgppID, offeringName, sectionName, status, inUse, perms)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          l.Detail.Title,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    offeringName,
			HeaderSubtitle: sectionName,
			HeaderIcon:     "icon-book-open",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate:   "sgpp-detail-content",
		Labels:            l,
		ActiveTab:         activeTab,
		TabItems:          tabItems,
		ID:                sgppID,
		SectionID:         sectionID,
		SectionName:       sectionName,
		SectionURL:        sectionURL,
		OfferingName:      offeringName,
		OfferingVariantID: offeringVariantID,
		VariantLabel:      variantLabel,
		CurriculumName:    curriculumName,
		PlanName:          planName,
		PeriodName:        periodName,
		Status:            status,
		StatusLabel:       statusLabel,
		StatusVariant:     statusVariant,
		AssignmentsCount:  assignmentsCount,
		CreatedDate:       createdDate,
		ModifiedDate:      modifiedDate,
		GradeSheetURL:     gradeSheetURL,
		EditURL:           editURL,
		InfoActions:       infoActions,
		AddURL:            route.ResolveURL(deps.Routes.AssignURL, "sgppId", sgppID),
	}

	if activeTab == "staff" {
		pageData.Assignments = buildAssignmentsTable(ctx, deps, sgppID, record, pageData)
		pageData.AssignmentsEmpty = pageData.Assignments == nil || len(pageData.Assignments.Rows) == 0
	}

	return pageData, nil
}

// buildInfoActions builds the Info tab's guarded Exclude/Restore/Remove
// buttons (S6). Edit rides the standard hx-get sheet-open link (rendered
// directly in the template, mirroring subscription_group's own Info tab).
// Confirm messages carry {{offering}}/{{section}} placeholders (lyngua.md
// convention) resolved here via strings.ReplaceAll — never re-templated —
// mirroring subscription_group_member's {{name}} substitution precedent.
func buildInfoActions(deps *DetailViewDeps, sgppID, offeringName, sectionName, status string, inUse bool, perms *types.UserPermissions) []InfoAction {
	l := deps.Labels
	canUpdate := perms.Can(sgppEntity, "update")
	canDelete := perms.Can(sgppEntity, "delete")
	tmpl := func(s string) string {
		s = strings.ReplaceAll(s, "{{offering}}", offeringName)
		s = strings.ReplaceAll(s, "{{section}}", sectionName)
		return s
	}
	var actions []InfoAction

	if status == "excluded" {
		actions = append(actions, InfoAction{
			Action:          "activate",
			Label:           l.Buttons.Restore,
			URL:             route.ResolveURL(deps.Routes.SetStatusURL, "sgppId", sgppID) + "?status=active",
			ItemName:        offeringName,
			ConfirmTitle:    l.Confirm.RestoreTitle,
			ConfirmMessage:  tmpl(l.Confirm.RestoreMessage),
			Disabled:        !canUpdate,
			DisabledTooltip: l.Assign.Unauthorized,
			TestID:          "sgpp-restore-" + sgppID,
		})
	} else {
		actions = append(actions, InfoAction{
			Action:          "deactivate",
			Label:           l.Buttons.Exclude,
			URL:             route.ResolveURL(deps.Routes.SetStatusURL, "sgppId", sgppID) + "?status=excluded",
			ItemName:        offeringName,
			ConfirmTitle:    l.Confirm.ExcludeTitle,
			ConfirmMessage:  tmpl(l.Confirm.ExcludeMessage),
			Disabled:        !canUpdate,
			DisabledTooltip: l.Assign.Unauthorized,
			TestID:          "sgpp-exclude-" + sgppID,
		})
	}

	removeAction := InfoAction{
		Action: "delete", Label: l.Buttons.Delete,
		URL:            deps.Routes.DeleteURL,
		ItemName:       offeringName,
		ConfirmTitle:   l.Confirm.DeleteTitle,
		ConfirmMessage: l.Confirm.DeleteMessage,
		TestID:         "sgpp-remove-" + sgppID,
	}
	switch {
	case inUse:
		removeAction.Disabled = true
		removeAction.DisabledTooltip = l.Confirm.RemoveBlockedInUse
	case !canDelete:
		removeAction.Disabled = true
		removeAction.DisabledTooltip = l.Assign.Unauthorized
	}
	actions = append(actions, removeAction)
	return actions
}

// countAssignments returns the number of active assignment rows on this class
// (the Teachers tab's count badge — mirrors the tab's own row source and
// read gate exactly, subscription_group/detail/staff.go's
// countSectionOfferings precedent, so the badge can never drift from the
// rendered rows).
func countAssignments(ctx context.Context, deps *DetailViewDeps, sgppID string) int {
	perms := view.GetUserPermissions(ctx)
	if perms == nil || !perms.Can(sgppsEntity, "list") {
		return 0
	}
	if deps.ListSubscriptionGroupProductPlanStaffs == nil || sgppID == "" {
		return 0
	}
	resp, err := deps.ListSubscriptionGroupProductPlanStaffs(ctx, &sgppspb.ListSubscriptionGroupProductPlanStaffsRequest{
		Filters: sgpp.AndFilters(sgpp.StringEqFilter("subscription_group_product_plan_id", sgppID), sgpp.BoolEqFilter("active", true)),
	})
	if err != nil || resp == nil {
		return 0
	}
	count := 0
	for _, a := range resp.GetData() {
		if a != nil && a.GetSubscriptionGroupProductPlanId() == sgppID && a.GetActive() {
			count++
		}
	}
	return count
}
