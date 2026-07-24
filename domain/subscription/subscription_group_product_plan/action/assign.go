package action

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"

	sgpp "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"
	jobtemplatephasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_phase"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
	sgpppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan"
	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

// sgppsEntity is the permission entity for the assignment edge (the pre-
// existing subscription_group_product_plan_staff grants — plan.md §5 D-5:
// assignment writes stay under this entity's :create/:update/:delete, evolved
// shape, not new grants).
const sgppsEntity = "subscription_group_product_plan_staff"

// roleTypePrimary / roleTypeSecondary are the two generic role tokens the
// assign form writes to sgpps.role (a plain string on the proto — no enum).
// Vertical vocabulary ("Teacher of record" / "Co-teacher") lives only in
// lyngua (Labels.Assign.RolePrimary/RoleSecondary); the wire tokens stay
// vertical-noun-free.
const (
	roleTypePrimary   = "primary"
	roleTypeSecondary = "secondary"
)

// AssignDeps holds all dependencies for the S3/S5 shared assign action
// (centymo.md §5/§7). GET renders the shared drawer: context (section +
// offering + curriculum), the Current assignments list (every row, each with
// its own Clear), a Teacher (eligible pool) + Role + Semester picker for
// adding one new assignment, and (when ?edit=<rowID> is set) that row's
// role/phase pre-filled with the teacher shown read-only (teacher is
// immutable on edit — change teacher = remove + add, centymo.md §7). POST
// creates one fresh (class, pps, phase-or-whole) row, updates role/phase on
// an edited row, or (clear=1) soft-deletes one row — the sgpps Clear idiom.
type AssignDeps struct {
	Routes       sgpp.Routes
	Labels       sgpp.Labels
	CommonLabels pyeza.CommonLabels

	ReadSubscriptionGroupProductPlan        func(ctx context.Context, req *sgpppb.ReadSubscriptionGroupProductPlanRequest) (*sgpppb.ReadSubscriptionGroupProductPlanResponse, error)
	ReadSubscriptionGroup                   func(ctx context.Context, req *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error)
	ListProductPlans                        func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListJobTemplates                        func(ctx context.Context, req *jobtemplatepb.ListJobTemplatesRequest) (*jobtemplatepb.ListJobTemplatesResponse, error)
	ListJobTemplatePhases                   func(ctx context.Context, req *jobtemplatephasepb.ListJobTemplatePhasesRequest) (*jobtemplatephasepb.ListJobTemplatePhasesResponse, error)
	ListProductPlanStaffs                   func(ctx context.Context, req *productplanstaffpb.ListProductPlanStaffsRequest) (*productplanstaffpb.ListProductPlanStaffsResponse, error)
	ListStaffNames                          func(ctx context.Context) map[string]string
	ListSubscriptionGroupProductPlanStaffs  func(ctx context.Context, req *sgppspb.ListSubscriptionGroupProductPlanStaffsRequest) (*sgppspb.ListSubscriptionGroupProductPlanStaffsResponse, error)
	ReadSubscriptionGroupProductPlanStaff   func(ctx context.Context, req *sgppspb.ReadSubscriptionGroupProductPlanStaffRequest) (*sgppspb.ReadSubscriptionGroupProductPlanStaffResponse, error)
	CreateSubscriptionGroupProductPlanStaff func(ctx context.Context, req *sgppspb.CreateSubscriptionGroupProductPlanStaffRequest) (*sgppspb.CreateSubscriptionGroupProductPlanStaffResponse, error)
	UpdateSubscriptionGroupProductPlanStaff func(ctx context.Context, req *sgppspb.UpdateSubscriptionGroupProductPlanStaffRequest) (*sgppspb.UpdateSubscriptionGroupProductPlanStaffResponse, error)
	DeleteSubscriptionGroupProductPlanStaff func(ctx context.Context, req *sgppspb.DeleteSubscriptionGroupProductPlanStaffRequest) (*sgppspb.DeleteSubscriptionGroupProductPlanStaffResponse, error)
}

// AssignRow is one row of the drawer's "Current" list (centymo.md §5).
type AssignRow struct {
	ID         string
	StaffName  string
	RoleLabel  string
	PhaseLabel string
}

// AssignData is the template data for the shared assign drawer.
type AssignData struct {
	FormAction     string
	SGPPID         string
	SectionName    string
	OfferingName   string
	CurriculumName string
	Current        []AssignRow
	StaffOptions   []types.SelectOption
	RoleOptions    []types.SelectOption
	PhaseOptions   []types.SelectOption
	IsEdit         bool
	EditRowID      string
	EditStaffLabel string
	Labels         sgpp.Labels
	CommonLabels   pyeza.CommonLabels
	WorkspaceID    string
	Nonce          string
}

func roleLabel(l sgpp.Labels, role string) string {
	if role == roleTypeSecondary {
		return l.Assign.RoleSecondary
	}
	return l.Assign.RolePrimary
}

func phaseLabel(l sgpp.Labels, phases []sgpp.CoveredPhase, phaseID string) string {
	if phaseID == "" {
		return l.Assign.AllPhasesChip
	}
	for _, p := range phases {
		if p.ID == phaseID {
			return p.Name
		}
	}
	return phaseID
}

func strPtr(s string) *string { return &s }

// buildAssignData assembles the shared drawer (GET). Every read is a single
// batched call scoped to this one class (espyna.md §1b "fetched lazily on the
// drawer GET, not eagerly for every row").
func buildAssignData(ctx context.Context, deps *AssignDeps, sgppID, editRowID string) (*AssignData, error) {
	readResp, err := deps.ReadSubscriptionGroupProductPlan(ctx, &sgpppb.ReadSubscriptionGroupProductPlanRequest{
		Data: &sgpppb.SubscriptionGroupProductPlan{Id: sgppID},
	})
	if err != nil || readResp == nil || len(readResp.GetData()) == 0 {
		return nil, errors.New(deps.Labels.Errors.NotFound)
	}
	record := readResp.GetData()[0]

	sectionName := record.GetSubscriptionGroupId()
	if deps.ReadSubscriptionGroup != nil {
		if resp, err := deps.ReadSubscriptionGroup(ctx, &subscriptiongrouppb.ReadSubscriptionGroupRequest{
			Data: &subscriptiongrouppb.SubscriptionGroup{Id: record.GetSubscriptionGroupId()},
		}); err == nil && resp != nil && len(resp.GetData()) > 0 {
			sectionName = resp.GetData()[0].GetName()
		}
	}

	offeringName := record.GetProductPlanId()
	var offeringVariantID string
	if deps.ListProductPlans != nil {
		if resp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
			Filters: sgpp.StringEqFilter("id", record.GetProductPlanId()),
		}); err == nil && resp != nil {
			for _, pp := range resp.GetData() {
				if pp.GetId() == record.GetProductPlanId() {
					offeringName = offeringLabel(pp)
					offeringVariantID = pp.GetProductVariantId()
				}
			}
		}
	}

	curriculumName := record.GetJobTemplateId()
	if deps.ListJobTemplates != nil && record.GetJobTemplateId() != "" {
		if resp, err := deps.ListJobTemplates(ctx, &jobtemplatepb.ListJobTemplatesRequest{
			Filters: sgpp.StringEqFilter("id", record.GetJobTemplateId()),
		}); err == nil && resp != nil {
			for _, t := range resp.GetData() {
				if t.GetId() == record.GetJobTemplateId() {
					curriculumName = t.GetName()
				}
			}
		}
	}

	var phases []sgpp.CoveredPhase
	if deps.ListJobTemplatePhases != nil && record.GetJobTemplateId() != "" {
		if resp, err := deps.ListJobTemplatePhases(ctx, &jobtemplatephasepb.ListJobTemplatePhasesRequest{
			Filters: sgpp.AndFilters(sgpp.StringEqFilter("job_template_id", record.GetJobTemplateId()), sgpp.BoolEqFilter("active", true)),
		}); err == nil && resp != nil {
			phases = sgpp.ResolveCoveredPhases(resp.GetData(), offeringVariantID)
		}
	}

	var names map[string]string
	if deps.ListStaffNames != nil {
		names = deps.ListStaffNames(ctx)
	}

	staffOptions := []types.SelectOption{}
	// ppsStaffByID resolves a product_plan_staff row id -> its staff id — the
	// ONLY sanctioned path to a display name (espyna.md §1b: "Never read
	// legacy f10 for display — the view must survive M7"). f10 (StaffId) is
	// the dual-write legacy mirror; current-assignment rows below resolve
	// through f13 (ProductPlanStaffId) via this map, never a.GetStaffId().
	ppsStaffByID := map[string]string{}
	if deps.ListProductPlanStaffs != nil {
		if resp, err := deps.ListProductPlanStaffs(ctx, &productplanstaffpb.ListProductPlanStaffsRequest{
			Filters: sgpp.AndFilters(sgpp.StringEqFilter("product_plan_id", record.GetProductPlanId()), sgpp.BoolEqFilter("active", true)),
		}); err == nil && resp != nil {
			for _, pps := range resp.GetData() {
				if pps == nil {
					continue
				}
				ppsStaffByID[pps.GetId()] = pps.GetStaffId()
				name := names[pps.GetStaffId()]
				if name == "" {
					name = pps.GetStaffId()
				}
				staffOptions = append(staffOptions, types.SelectOption{Value: pps.GetId(), Label: name})
			}
		}
	}
	sort.SliceStable(staffOptions, func(i, j int) bool {
		return strings.ToLower(staffOptions[i].Label) < strings.ToLower(staffOptions[j].Label)
	})

	current := []AssignRow{}
	rawByID := map[string]*sgppspb.SubscriptionGroupProductPlanStaff{}
	if deps.ListSubscriptionGroupProductPlanStaffs != nil {
		if resp, err := deps.ListSubscriptionGroupProductPlanStaffs(ctx, &sgppspb.ListSubscriptionGroupProductPlanStaffsRequest{
			Filters: sgpp.AndFilters(sgpp.StringEqFilter("subscription_group_product_plan_id", sgppID), sgpp.BoolEqFilter("active", true)),
		}); err == nil && resp != nil {
			for _, a := range resp.GetData() {
				if a == nil || a.GetSubscriptionGroupProductPlanId() != sgppID {
					continue
				}
				rawByID[a.GetId()] = a
				staffID := ppsStaffByID[a.GetProductPlanStaffId()]
				staffName := staffID
				if n := names[staffID]; n != "" {
					staffName = n
				}
				if staffName == "" {
					staffName = a.GetProductPlanStaffId()
				}
				current = append(current, AssignRow{
					ID:         a.GetId(),
					StaffName:  staffName,
					RoleLabel:  roleLabel(deps.Labels, a.GetRole()),
					PhaseLabel: phaseLabel(deps.Labels, phases, a.GetJobTemplatePhaseId()),
				})
			}
		}
	}
	sort.SliceStable(current, func(i, j int) bool { return current[i].StaffName < current[j].StaffName })

	defaultRole := roleTypePrimary
	defaultPhase := ""
	isEdit := false
	editStaffLabel := ""
	if editRowID != "" {
		if raw, ok := rawByID[editRowID]; ok {
			isEdit = true
			defaultRole = raw.GetRole()
			defaultPhase = raw.GetJobTemplatePhaseId()
			editStaffID := ppsStaffByID[raw.GetProductPlanStaffId()]
			editStaffLabel = names[editStaffID]
			if editStaffLabel == "" {
				editStaffLabel = editStaffID
			}
			if editStaffLabel == "" {
				editStaffLabel = raw.GetProductPlanStaffId()
			}
		}
	}

	roleOptions := []types.SelectOption{
		{Value: roleTypePrimary, Label: deps.Labels.Assign.RolePrimary, Selected: defaultRole != roleTypeSecondary},
		{Value: roleTypeSecondary, Label: deps.Labels.Assign.RoleSecondary, Selected: defaultRole == roleTypeSecondary},
	}
	phaseOptions := []types.SelectOption{{Value: "", Label: deps.Labels.Assign.AllPhasesChip, Selected: defaultPhase == ""}}
	for _, p := range phases {
		phaseOptions = append(phaseOptions, types.SelectOption{Value: p.ID, Label: p.Name, Selected: p.ID == defaultPhase})
	}

	return &AssignData{
		FormAction:     route.ResolveURL(deps.Routes.AssignURL, "sgppId", sgppID),
		SGPPID:         sgppID,
		SectionName:    sectionName,
		OfferingName:   offeringName,
		CurriculumName: curriculumName,
		Current:        current,
		StaffOptions:   staffOptions,
		RoleOptions:    roleOptions,
		PhaseOptions:   phaseOptions,
		IsEdit:         isEdit,
		EditRowID:      editRowID,
		EditStaffLabel: editStaffLabel,
		Labels:         deps.Labels,
		CommonLabels:   deps.CommonLabels,
	}, nil
}

// NewAssignAction handles GET (the shared drawer form) / POST (create one
// assignment, update an edited row's role/phase, or clear=1+id soft-delete
// one row) for both the S3 quick-drawer mount and the S5 Teachers-tab Add/Edit
// drawer (centymo.md §5/§7 — same AssignURL, same validation, addressed by
// {sgppId} so either mount can mount it independently).
func NewAssignAction(deps *AssignDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		sgppID := viewCtx.Request.PathValue("sgppId")
		if sgppID == "" {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}

		if viewCtx.Request.Method == http.MethodGet {
			if !perms.CanAny(sgppsEntity+":create", sgppsEntity+":update") {
				return view.HTMXError(deps.Labels.Assign.Unauthorized)
			}
			editRowID := strings.TrimSpace(viewCtx.Request.URL.Query().Get("edit"))
			data, err := buildAssignData(ctx, deps, sgppID, editRowID)
			if err != nil {
				return view.HTMXError(err.Error())
			}
			return view.OK("sgpp-assign-drawer", data)
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.UpdateFailed)
		}

		// Clear — the sgpps soft-delete idiom (a Current row's inline form).
		if viewCtx.Request.FormValue("clear") == "1" {
			rowID := strings.TrimSpace(viewCtx.Request.FormValue("id"))
			if rowID == "" {
				rowID = strings.TrimSpace(viewCtx.Request.URL.Query().Get("id"))
			}
			if rowID == "" {
				return view.HTMXError(deps.Labels.Errors.NotFound)
			}
			if !perms.Can(sgppsEntity, "delete") {
				return view.HTMXError(deps.Labels.Assign.Unauthorized)
			}
			if _, err := deps.DeleteSubscriptionGroupProductPlanStaff(ctx, &sgppspb.DeleteSubscriptionGroupProductPlanStaffRequest{
				Data: &sgppspb.SubscriptionGroupProductPlanStaff{Id: rowID},
			}); err != nil {
				return view.HTMXError(err.Error())
			}
			return view.HTMXSuccess(sgpp.StaffTabPanelID)
		}

		role := strings.TrimSpace(viewCtx.Request.FormValue("role"))
		if role != roleTypeSecondary {
			role = roleTypePrimary
		}
		phaseID := strings.TrimSpace(viewCtx.Request.FormValue("job_template_phase_id"))
		rowID := strings.TrimSpace(viewCtx.Request.FormValue("row_id"))

		if rowID != "" {
			// Edit — teacher immutable; only role/phase change.
			if !perms.Can(sgppsEntity, "update") {
				return view.HTMXError(deps.Labels.Assign.Unauthorized)
			}
			readResp, err := deps.ReadSubscriptionGroupProductPlanStaff(ctx, &sgppspb.ReadSubscriptionGroupProductPlanStaffRequest{
				Data: &sgppspb.SubscriptionGroupProductPlanStaff{Id: rowID},
			})
			if err != nil || readResp == nil || len(readResp.GetData()) == 0 {
				return view.HTMXError(deps.Labels.Errors.NotFound)
			}
			existing := readResp.GetData()[0]
			update := &sgppspb.SubscriptionGroupProductPlanStaff{
				Id:                             rowID,
				Role:                           role,
				SubscriptionGroupProductPlanId: existing.SubscriptionGroupProductPlanId,
				ProductPlanStaffId:             existing.ProductPlanStaffId,
				JobTemplatePhaseId:             strPtr(phaseID),
			}
			if _, err := deps.UpdateSubscriptionGroupProductPlanStaff(ctx, &sgppspb.UpdateSubscriptionGroupProductPlanStaffRequest{Data: update}); err != nil {
				return view.HTMXError(err.Error())
			}
			return view.HTMXSuccess(sgpp.StaffTabPanelID)
		}

		// Add — a fresh (class, pps, phase-or-whole) row. A distinct pps never
		// collides with an existing occupant's row (co-teach-safe by construction —
		// the unique index is (class, pps, phase), plan.md §5 "Multiple staff per
		// phase allowed").
		if !perms.CanAny(sgppsEntity+":create", sgppsEntity+":update") {
			return view.HTMXError(deps.Labels.Assign.Unauthorized)
		}
		ppsID := strings.TrimSpace(viewCtx.Request.FormValue("product_plan_staff_id"))
		if ppsID == "" {
			return view.HTMXError(deps.Labels.Errors.ValidationFailed)
		}
		create := &sgppspb.SubscriptionGroupProductPlanStaff{
			Role:                           role,
			SubscriptionGroupProductPlanId: strPtr(sgppID),
			ProductPlanStaffId:             strPtr(ppsID),
		}
		if phaseID != "" {
			create.JobTemplatePhaseId = strPtr(phaseID)
		}
		if _, err := deps.CreateSubscriptionGroupProductPlanStaff(ctx, &sgppspb.CreateSubscriptionGroupProductPlanStaffRequest{Data: create}); err != nil {
			return view.HTMXError(err.Error())
		}
		return view.HTMXSuccess(sgpp.StaffTabPanelID)
	})
}
