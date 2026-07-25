package subscription

import (
	"context"

	subscriptiongroupaction "github.com/erniealice/centymo-golang/domain/subscription/subscription_group/action"
	subscriptiongroupdetail "github.com/erniealice/centymo-golang/domain/subscription/subscription_group/detail"
	subscriptiongrouplist "github.com/erniealice/centymo-golang/domain/subscription/subscription_group/list"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	view "github.com/erniealice/pyeza-golang/view"

	epkg "github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	clientattributepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client_attribute"
	jobtemplatephasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template_phase"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productplanstaffpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan_staff"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
	subscriptiongroupmemberpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_member"
	sgpppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan"
	sgppspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_product_plan_staff"
)

// SubscriptionGroupModuleDeps holds all dependencies for the
// subscription_group module (the education "section / cohort" cohort).
type SubscriptionGroupModuleDeps struct {
	Routes       epkg.Routes
	Labels       epkg.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListSubscriptionGroups  func(ctx context.Context, req *subscriptiongrouppb.ListSubscriptionGroupsRequest) (*subscriptiongrouppb.ListSubscriptionGroupsResponse, error)
	ReadSubscriptionGroup   func(ctx context.Context, req *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error)
	CreateSubscriptionGroup func(ctx context.Context, req *subscriptiongrouppb.CreateSubscriptionGroupRequest) (*subscriptiongrouppb.CreateSubscriptionGroupResponse, error)
	UpdateSubscriptionGroup func(ctx context.Context, req *subscriptiongrouppb.UpdateSubscriptionGroupRequest) (*subscriptiongrouppb.UpdateSubscriptionGroupResponse, error)
	DeleteSubscriptionGroup func(ctx context.Context, req *subscriptiongrouppb.DeleteSubscriptionGroupRequest) (*subscriptiongrouppb.DeleteSubscriptionGroupResponse, error)

	// Program (plan) + period (price_schedule) pickers / display lookups.
	ListPlans          func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)

	// Subscriptions (roster) tab: members + client/subscription name resolution.
	ListSubscriptionGroupMembers func(ctx context.Context, req *subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest) (*subscriptiongroupmemberpb.ListSubscriptionGroupMembersResponse, error)
	ListClients                  func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ListSubscriptions            func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)

	// Roster banding (app-configured gender bands): generic
	// "client_attributes.<code>" option + its workspace-bound attribute closures.
	// Zero value → flat roster (service-admin unaffected).
	Options                  epkg.Options
	ListClientAttributes     func(ctx context.Context, req *clientattributepb.ListClientAttributesRequest) (*clientattributepb.ListClientAttributesResponse, error)
	ResolveAttributeIDByCode func(ctx context.Context, code string) (string, error)

	// Teaching-staff tab (§6.1–§6.3): the section assignment grid + inline upsert.
	ListProductPlans                       func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListSubscriptionGroupProductPlanStaffs func(ctx context.Context, req *sgppspb.ListSubscriptionGroupProductPlanStaffsRequest) (*sgppspb.ListSubscriptionGroupProductPlanStaffsResponse, error)
	ListProductPlanStaffs                  func(ctx context.Context, req *productplanstaffpb.ListProductPlanStaffsRequest) (*productplanstaffpb.ListProductPlanStaffsResponse, error)
	ListStaffNames                         func(ctx context.Context) map[string]string
	AssignGroupServicer                    func(ctx context.Context, subscriptionGroupID, productPlanID, staffID, role string) (string, error)
	ProductPlanStaffListURL                string

	// M4 row-source flip (plan.md §2, centymo.md §3): the tab's primary row
	// source — active subscription_group_product_plan (class) rows.
	ListSubscriptionGroupProductPlans       func(ctx context.Context, req *sgpppb.ListSubscriptionGroupProductPlansRequest) (*sgpppb.ListSubscriptionGroupProductPlansResponse, error)
	ListProductVariants                     func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	ListJobTemplatePhases                   func(ctx context.Context, req *jobtemplatephasepb.ListJobTemplatePhasesRequest) (*jobtemplatephasepb.ListJobTemplatePhasesResponse, error)
	GetSubscriptionGroupProductPlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
	SGPPDetailURL                           func(sectionID, sgppID string) string
	SGPPAssignURL                           func(sgppID string) string
	SGPPSetStatusURL                        func(sgppID, status string) string
	SGPPPickerURL                           func(sectionID string) string
	SGPPDeleteURL                           string
	SGPPBulkDeleteURL                       string

	attachment.AttachmentOps // attachments tab
	auditlog.AuditOps        // audit tab (nil today → renders empty)

	// Optional reference checker; nil disables delete gating (no ref-checker
	// method exists for subscription_group yet).
	GetSubscriptionGroupInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// SubscriptionGroupModule holds all constructed subscription_group views.
type SubscriptionGroupModule struct {
	routes           epkg.Routes
	Dashboard        view.View
	List             view.View
	Table            view.View
	Add              view.View
	Edit             view.View
	Delete           view.View
	BulkDelete       view.View
	SetStatus        view.View
	BulkSetStatus    view.View
	Detail           view.View
	TabAction        view.View
	Assign           view.View
	AttachmentUpload view.View
	AttachmentDelete view.View
}

// NewSubscriptionGroupModule creates the subscription_group module with all
// views wired.
func NewSubscriptionGroupModule(deps *SubscriptionGroupModuleDeps) *SubscriptionGroupModule {
	actionDeps := &subscriptiongroupaction.Deps{
		Routes:                       deps.Routes,
		Labels:                       deps.Labels,
		CreateSubscriptionGroup:      deps.CreateSubscriptionGroup,
		ReadSubscriptionGroup:        deps.ReadSubscriptionGroup,
		UpdateSubscriptionGroup:      deps.UpdateSubscriptionGroup,
		DeleteSubscriptionGroup:      deps.DeleteSubscriptionGroup,
		ListPlans:                    deps.ListPlans,
		ListPriceSchedules:           deps.ListPriceSchedules,
		GetSubscriptionGroupInUseIDs: deps.GetSubscriptionGroupInUseIDs,
	}

	listDeps := &subscriptiongrouplist.ListViewDeps{
		Routes:                       deps.Routes,
		ListSubscriptionGroups:       deps.ListSubscriptionGroups,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
		TableLabels:                  deps.TableLabels,
		GetSubscriptionGroupInUseIDs: deps.GetSubscriptionGroupInUseIDs,
	}
	listView := subscriptiongrouplist.NewView(listDeps)
	tableView := subscriptiongrouplist.NewTableView(listDeps)

	detailDeps := &subscriptiongroupdetail.DetailViewDeps{
		Routes:                       deps.Routes,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
		TableLabels:                  deps.TableLabels,
		ReadSubscriptionGroup:        deps.ReadSubscriptionGroup,
		ListPlans:                    deps.ListPlans,
		ListPriceSchedules:           deps.ListPriceSchedules,
		ListSubscriptionGroupMembers: deps.ListSubscriptionGroupMembers,
		ListClients:                  deps.ListClients,
		ListSubscriptions:            deps.ListSubscriptions,
		Options:                      deps.Options,
		ListClientAttributes:         deps.ListClientAttributes,
		ResolveAttributeIDByCode:     deps.ResolveAttributeIDByCode,

		ListProductPlans:                       deps.ListProductPlans,
		ListSubscriptionGroupProductPlanStaffs: deps.ListSubscriptionGroupProductPlanStaffs,
		ListProductPlanStaffs:                  deps.ListProductPlanStaffs,
		ListStaffNames:                         deps.ListStaffNames,
		AssignGroupServicer:                    deps.AssignGroupServicer,
		ProductPlanStaffListURL:                deps.ProductPlanStaffListURL,

		ListSubscriptionGroupProductPlans:       deps.ListSubscriptionGroupProductPlans,
		ListProductVariants:                     deps.ListProductVariants,
		ListJobTemplatePhases:                   deps.ListJobTemplatePhases,
		GetSubscriptionGroupProductPlanInUseIDs: deps.GetSubscriptionGroupProductPlanInUseIDs,
		SGPPDetailURL:                           deps.SGPPDetailURL,
		SGPPAssignURL:                           deps.SGPPAssignURL,
		SGPPSetStatusURL:                        deps.SGPPSetStatusURL,
		SGPPPickerURL:                           deps.SGPPPickerURL,
		SGPPDeleteURL:                           deps.SGPPDeleteURL,
		SGPPBulkDeleteURL:                       deps.SGPPBulkDeleteURL,

		AttachmentOps: deps.AttachmentOps,
		AuditOps:      deps.AuditOps,
	}

	return &SubscriptionGroupModule{
		routes:           deps.Routes,
		Dashboard:        listView,
		List:             listView,
		Table:            tableView,
		Add:              subscriptiongroupaction.NewAddAction(actionDeps),
		Edit:             subscriptiongroupaction.NewEditAction(actionDeps),
		Delete:           subscriptiongroupaction.NewDeleteAction(actionDeps),
		BulkDelete:       subscriptiongroupaction.NewBulkDeleteAction(actionDeps),
		SetStatus:        subscriptiongroupaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    subscriptiongroupaction.NewBulkSetStatusAction(actionDeps),
		Detail:           subscriptiongroupdetail.NewView(detailDeps),
		TabAction:        subscriptiongroupdetail.NewTabAction(detailDeps),
		Assign:           subscriptiongroupdetail.NewAssignAction(detailDeps),
		AttachmentUpload: subscriptiongroupdetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: subscriptiongroupdetail.NewAttachmentDeleteAction(detailDeps),
	}
}

// RegisterRoutes registers all subscription_group routes.
func (m *SubscriptionGroupModule) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)

	if m.Detail != nil && m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil && m.routes.TabActionURL != "" {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
	if m.Assign != nil && m.routes.AssignURL != "" {
		r.POST(m.routes.AssignURL, m.Assign)
	}
	if m.AttachmentUpload != nil && m.routes.AttachmentUploadURL != "" {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
	}
	if m.AttachmentDelete != nil && m.routes.AttachmentDeleteURL != "" {
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
}
