package subscription

import (
	"context"

	sgppaction "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan/action"
	sgppdetail "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan/detail"
	sgpplist "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	sgpp "github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan"
	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group_product_plan/form"

	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"
	jobtemplatephasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_phase"
	jobtemplaterelationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_relation"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
	sgpppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan"
	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

// SubscriptionGroupProductPlanModuleDeps holds all dependencies for the
// subscription_group_product_plan module (THE CLASS entity — plan.md §1.1b).
// Fields are deduplicated across the S7 admin surfaces, the S2 picker, the
// S3/S5 assign drawer, S6 set-status and the S4 class page — every read is a
// batched LIST_IN call per espyna.md §1b, never per-row.
type SubscriptionGroupProductPlanModuleDeps struct {
	Routes       sgpp.Routes
	Labels       sgpp.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	Options      sgpp.Options

	// Core CRUD (sgpp).
	ListSubscriptionGroupProductPlans  func(ctx context.Context, req *sgpppb.ListSubscriptionGroupProductPlansRequest) (*sgpppb.ListSubscriptionGroupProductPlansResponse, error)
	ReadSubscriptionGroupProductPlan   func(ctx context.Context, req *sgpppb.ReadSubscriptionGroupProductPlanRequest) (*sgpppb.ReadSubscriptionGroupProductPlanResponse, error)
	CreateSubscriptionGroupProductPlan func(ctx context.Context, req *sgpppb.CreateSubscriptionGroupProductPlanRequest) (*sgpppb.CreateSubscriptionGroupProductPlanResponse, error)
	UpdateSubscriptionGroupProductPlan func(ctx context.Context, req *sgpppb.UpdateSubscriptionGroupProductPlanRequest) (*sgpppb.UpdateSubscriptionGroupProductPlanResponse, error)
	DeleteSubscriptionGroupProductPlan func(ctx context.Context, req *sgpppb.DeleteSubscriptionGroupProductPlanRequest) (*sgpppb.DeleteSubscriptionGroupProductPlanResponse, error)
	// GetSubscriptionGroupProductPlanInUseIDs gates S7 delete + S4 Remove (the
	// server-side guard inside espyna's Delete/Update use cases is the hard
	// backstop; this closure only drives the UI's disabled state).
	GetSubscriptionGroupProductPlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// Assignment edge (sgpps) — S3/S5 the assign drawer.
	ReadSubscriptionGroupProductPlanStaff   func(ctx context.Context, req *sgppspb.ReadSubscriptionGroupProductPlanStaffRequest) (*sgppspb.ReadSubscriptionGroupProductPlanStaffResponse, error)
	ListSubscriptionGroupProductPlanStaffs  func(ctx context.Context, req *sgppspb.ListSubscriptionGroupProductPlanStaffsRequest) (*sgppspb.ListSubscriptionGroupProductPlanStaffsResponse, error)
	CreateSubscriptionGroupProductPlanStaff func(ctx context.Context, req *sgppspb.CreateSubscriptionGroupProductPlanStaffRequest) (*sgppspb.CreateSubscriptionGroupProductPlanStaffResponse, error)
	UpdateSubscriptionGroupProductPlanStaff func(ctx context.Context, req *sgppspb.UpdateSubscriptionGroupProductPlanStaffRequest) (*sgppspb.UpdateSubscriptionGroupProductPlanStaffResponse, error)
	DeleteSubscriptionGroupProductPlanStaff func(ctx context.Context, req *sgppspb.DeleteSubscriptionGroupProductPlanStaffRequest) (*sgppspb.DeleteSubscriptionGroupProductPlanStaffResponse, error)

	// Cross-entity reads shared by list/, action/ and detail/.
	ReadSubscriptionGroup    func(ctx context.Context, req *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error)
	ListPlans                func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules       func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)
	ListProductPlans         func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListProductVariants      func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	ListProductPlanStaffs    func(ctx context.Context, req *productplanstaffpb.ListProductPlanStaffsRequest) (*productplanstaffpb.ListProductPlanStaffsResponse, error)
	ListJobTemplates         func(ctx context.Context, req *jobtemplatepb.ListJobTemplatesRequest) (*jobtemplatepb.ListJobTemplatesResponse, error)
	ListJobTemplatePhases    func(ctx context.Context, req *jobtemplatephasepb.ListJobTemplatePhasesRequest) (*jobtemplatephasepb.ListJobTemplatePhasesResponse, error)
	ListJobTemplateRelations func(ctx context.Context, req *jobtemplaterelationpb.ListJobTemplateRelationsRequest) (*jobtemplaterelationpb.ListJobTemplateRelationsResponse, error)
	ListStaffNames           func(ctx context.Context) map[string]string

	// S7 admin list name-resolution batches + Teachers-count batch.
	ListSubscriptionGroupNames func(ctx context.Context) map[string]string
	ListProductPlanNames       func(ctx context.Context) map[string]string
	ListJobTemplateNames       func(ctx context.Context) map[string]string
	CountAssignmentsBySGPPID   func(ctx context.Context, sgppIDs []string) (map[string]int, error)

	// S7 admin add/edit drawer FK pickers — nil disables the picker.
	ListSubscriptionGroupOptions func(ctx context.Context) []form.Pair
	ListProductPlanOptions       func(ctx context.Context) []form.Pair
	ListJobTemplateOptions       func(ctx context.Context) []form.Pair

	// RouteMap cross-links (plan.md §1.1b — never hardcoded).
	SectionDetailURL func(ctx context.Context, sectionID string) string
	GradeSheetURL    func(ctx context.Context, jobTemplateID string) string
}

// SubscriptionGroupProductPlanModule holds all constructed views.
type SubscriptionGroupProductPlanModule struct {
	routes     sgpp.Routes
	Dashboard  view.View
	List       view.View
	Table      view.View
	Add        view.View
	Edit       view.View
	Delete     view.View
	BulkDelete view.View
	Picker     view.View
	Assign     view.View
	SetStatus  view.View
	Detail     view.View
	TabAction  view.View
}

// NewSubscriptionGroupProductPlanModule creates the module with all views
// wired (centymo.md §1 package tree).
func NewSubscriptionGroupProductPlanModule(deps *SubscriptionGroupProductPlanModuleDeps) *SubscriptionGroupProductPlanModule {
	actionDeps := &sgppaction.Deps{
		Routes:                                  deps.Routes,
		Labels:                                  deps.Labels,
		CreateSubscriptionGroupProductPlan:      deps.CreateSubscriptionGroupProductPlan,
		ReadSubscriptionGroupProductPlan:        deps.ReadSubscriptionGroupProductPlan,
		UpdateSubscriptionGroupProductPlan:      deps.UpdateSubscriptionGroupProductPlan,
		DeleteSubscriptionGroupProductPlan:      deps.DeleteSubscriptionGroupProductPlan,
		GetSubscriptionGroupProductPlanInUseIDs: deps.GetSubscriptionGroupProductPlanInUseIDs,
		ListSubscriptionGroupOptions:            deps.ListSubscriptionGroupOptions,
		ListProductPlanOptions:                  deps.ListProductPlanOptions,
		ListJobTemplateOptions:                  deps.ListJobTemplateOptions,
	}

	listDeps := &sgpplist.ListViewDeps{
		Routes:                                  deps.Routes,
		ListSubscriptionGroupProductPlans:       deps.ListSubscriptionGroupProductPlans,
		Labels:                                  deps.Labels,
		CommonLabels:                            deps.CommonLabels,
		TableLabels:                             deps.TableLabels,
		GetSubscriptionGroupProductPlanInUseIDs: deps.GetSubscriptionGroupProductPlanInUseIDs,
		ListSubscriptionGroupNames:              deps.ListSubscriptionGroupNames,
		ListProductPlanNames:                    deps.ListProductPlanNames,
		ListJobTemplateNames:                    deps.ListJobTemplateNames,
		CountAssignmentsBySGPPID:                deps.CountAssignmentsBySGPPID,
	}
	listView := sgpplist.NewView(listDeps)
	tableView := sgpplist.NewTableView(listDeps)

	pickerDeps := &sgppaction.PickerDeps{
		Routes:                             deps.Routes,
		Labels:                             deps.Labels,
		CommonLabels:                       deps.CommonLabels,
		Options:                            deps.Options,
		ReadSubscriptionGroup:              deps.ReadSubscriptionGroup,
		ListPlans:                          deps.ListPlans,
		ListPriceSchedules:                 deps.ListPriceSchedules,
		ListProductPlans:                   deps.ListProductPlans,
		ListProductVariants:                deps.ListProductVariants,
		ListSubscriptionGroupProductPlans:  deps.ListSubscriptionGroupProductPlans,
		ListJobTemplateRelations:           deps.ListJobTemplateRelations,
		ListJobTemplates:                   deps.ListJobTemplates,
		CreateSubscriptionGroupProductPlan: deps.CreateSubscriptionGroupProductPlan,
	}

	assignDeps := &sgppaction.AssignDeps{
		Routes:                                  deps.Routes,
		Labels:                                  deps.Labels,
		CommonLabels:                            deps.CommonLabels,
		ReadSubscriptionGroupProductPlan:        deps.ReadSubscriptionGroupProductPlan,
		ReadSubscriptionGroup:                   deps.ReadSubscriptionGroup,
		ListProductPlans:                        deps.ListProductPlans,
		ListJobTemplates:                        deps.ListJobTemplates,
		ListJobTemplatePhases:                   deps.ListJobTemplatePhases,
		ListProductPlanStaffs:                   deps.ListProductPlanStaffs,
		ListStaffNames:                          deps.ListStaffNames,
		ListSubscriptionGroupProductPlanStaffs:  deps.ListSubscriptionGroupProductPlanStaffs,
		ReadSubscriptionGroupProductPlanStaff:   deps.ReadSubscriptionGroupProductPlanStaff,
		CreateSubscriptionGroupProductPlanStaff: deps.CreateSubscriptionGroupProductPlanStaff,
		UpdateSubscriptionGroupProductPlanStaff: deps.UpdateSubscriptionGroupProductPlanStaff,
		DeleteSubscriptionGroupProductPlanStaff: deps.DeleteSubscriptionGroupProductPlanStaff,
	}

	setStatusDeps := &sgppaction.SetStatusDeps{
		Routes:                             deps.Routes,
		Labels:                             deps.Labels,
		ReadSubscriptionGroupProductPlan:   deps.ReadSubscriptionGroupProductPlan,
		UpdateSubscriptionGroupProductPlan: deps.UpdateSubscriptionGroupProductPlan,
	}

	detailDeps := &sgppdetail.DetailViewDeps{
		Routes:                                  deps.Routes,
		Labels:                                  deps.Labels,
		CommonLabels:                            deps.CommonLabels,
		TableLabels:                             deps.TableLabels,
		ReadSubscriptionGroupProductPlan:        deps.ReadSubscriptionGroupProductPlan,
		ReadSubscriptionGroup:                   deps.ReadSubscriptionGroup,
		ListPlans:                               deps.ListPlans,
		ListPriceSchedules:                      deps.ListPriceSchedules,
		ListProductPlans:                        deps.ListProductPlans,
		ListProductVariants:                     deps.ListProductVariants,
		ListJobTemplates:                        deps.ListJobTemplates,
		ListJobTemplatePhases:                   deps.ListJobTemplatePhases,
		ListProductPlanStaffs:                   deps.ListProductPlanStaffs,
		ListStaffNames:                          deps.ListStaffNames,
		ListSubscriptionGroupProductPlanStaffs:  deps.ListSubscriptionGroupProductPlanStaffs,
		GetSubscriptionGroupProductPlanInUseIDs: deps.GetSubscriptionGroupProductPlanInUseIDs,
		SectionDetailURL:                        deps.SectionDetailURL,
		GradeSheetURL:                           deps.GradeSheetURL,
	}

	return &SubscriptionGroupProductPlanModule{
		routes:     deps.Routes,
		Dashboard:  listView,
		List:       listView,
		Table:      tableView,
		Add:        sgppaction.NewAddAction(actionDeps),
		Edit:       sgppaction.NewEditAction(actionDeps),
		Delete:     sgppaction.NewDeleteAction(actionDeps),
		BulkDelete: sgppaction.NewBulkDeleteAction(actionDeps),
		Picker:     sgppaction.NewPickerAction(pickerDeps),
		Assign:     sgppaction.NewAssignAction(assignDeps),
		SetStatus:  sgppaction.NewSetStatusAction(setStatusDeps),
		Detail:     sgppdetail.NewView(detailDeps),
		TabAction:  sgppdetail.NewTabAction(detailDeps),
	}
}

// RegisterRoutes registers all subscription_group_product_plan routes.
func (m *SubscriptionGroupProductPlanModule) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)

	if m.Picker != nil && m.routes.PickerURL != "" {
		r.GET(m.routes.PickerURL, m.Picker)
		r.POST(m.routes.PickerURL, m.Picker)
	}
	if m.Assign != nil && m.routes.AssignURL != "" {
		r.GET(m.routes.AssignURL, m.Assign)
		r.POST(m.routes.AssignURL, m.Assign)
	}
	if m.SetStatus != nil && m.routes.SetStatusURL != "" {
		r.POST(m.routes.SetStatusURL, m.SetStatus)
	}
	if m.Detail != nil && m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil && m.routes.TabActionURL != "" {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
}
