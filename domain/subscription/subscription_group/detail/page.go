package detail

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
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

// DetailViewDeps holds view dependencies for the subscription group detail page.
type DetailViewDeps struct {
	Routes                subscription_group.Routes
	Labels                subscription_group.Labels
	CommonLabels          pyeza.CommonLabels
	TableLabels           types.TableLabels
	ReadSubscriptionGroup func(ctx context.Context, req *subscriptiongrouppb.ReadSubscriptionGroupRequest) (*subscriptiongrouppb.ReadSubscriptionGroupResponse, error)
	ListPlans             func(ctx context.Context, req *planpb.ListPlansRequest) (*planpb.ListPlansResponse, error)
	ListPriceSchedules    func(ctx context.Context, req *priceschedulepb.ListPriceSchedulesRequest) (*priceschedulepb.ListPriceSchedulesResponse, error)

	// Subscriptions tab (the section roster): subscription_group_member rows,
	// with client + subscription display names resolved via a batch map.
	ListSubscriptionGroupMembers func(ctx context.Context, req *subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest) (*subscriptiongroupmemberpb.ListSubscriptionGroupMembersResponse, error)
	ListClients                  func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	ListSubscriptions            func(ctx context.Context, req *subscriptionpb.ListSubscriptionsRequest) (*subscriptionpb.ListSubscriptionsResponse, error)

	// Roster banding — app-configured via the block's EngineBlock option, generic
	// "client_attributes.<code>" grammar. A zero Options (or an unresolvable/
	// unwired code) → today's flat roster, byte-identical (service-admin, which
	// sets nothing, is unaffected). The two closures are workspace-bound at the
	// espyna adapter, mirroring fayna's report-cards grid wiring.
	Options                  subscription_group.Options
	ListClientAttributes     func(ctx context.Context, req *clientattributepb.ListClientAttributesRequest) (*clientattributepb.ListClientAttributesResponse, error)
	ResolveAttributeIDByCode func(ctx context.Context, code string) (string, error)

	// Teaching-staff tab (§6.2): the group grid composes these existing List
	// use cases — offerings (ListProductPlans by plan_id), current assignments
	// (ListSubscriptionGroupProductPlanStaffs by group), eligible pools
	// (ListProductPlanStaffs by product_plan_id) — plus a staff id→name batch.
	// AssignGroupServicer is the §6.3 upsert (workspace read from ctx by the
	// espyna consumer seam; group id from the signed path). All nil-safe: an
	// unwired dep degrades the grid to empty, never an error.
	ListProductPlans                       func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	ListSubscriptionGroupProductPlanStaffs func(ctx context.Context, req *sgppspb.ListSubscriptionGroupProductPlanStaffsRequest) (*sgppspb.ListSubscriptionGroupProductPlanStaffsResponse, error)
	ListProductPlanStaffs                  func(ctx context.Context, req *productplanstaffpb.ListProductPlanStaffsRequest) (*productplanstaffpb.ListProductPlanStaffsResponse, error)
	ListStaffNames                         func(ctx context.Context) map[string]string
	AssignGroupServicer                    func(ctx context.Context, subscriptionGroupID, productPlanID, staffID, role string) (string, error)
	// ProductPlanStaffListURL is the eligibility-management list route (with a
	// {status} placeholder) the empty-pool gate links to; "" hides the link.
	ProductPlanStaffListURL string

	// -- M4 row-source flip (plan.md §2 / centymo.md §3): the tab's primary row
	// source becomes active subscription_group_product_plan (class) rows. Every
	// read is a batched LIST_IN call (espyna.md §1b); a nil closure degrades
	// gracefully (empty rows / disabled actions), never an error. The legacy
	// fields above stay wired for the transitional fallback (zero class rows ⇒
	// the section renders today's derived offering rows, actions disabled).
	ListSubscriptionGroupProductPlans       func(ctx context.Context, req *sgpppb.ListSubscriptionGroupProductPlansRequest) (*sgpppb.ListSubscriptionGroupProductPlansResponse, error)
	ListProductVariants                     func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	ListJobTemplatePhases                   func(ctx context.Context, req *jobtemplatephasepb.ListJobTemplatePhasesRequest) (*jobtemplatephasepb.ListJobTemplatePhasesResponse, error)
	GetSubscriptionGroupProductPlanInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// Cross-domain URL closures for the subscription_group_product_plan (class)
	// module's own routes — resolved by the container via compose.RoutesOf,
	// never hardcoded (plan.md §1.1b). Nil-safe: an unwired closure disables
	// the corresponding row action instead of erroring.
	SGPPDetailURL    func(sectionID, sgppID string) string
	SGPPAssignURL    func(sgppID string) string
	SGPPSetStatusURL func(sgppID, status string) string
	SGPPPickerURL    func(sectionID string) string
	SGPPDeleteURL    string

	attachment.AttachmentOps // attachments tab
	auditlog.AuditOps        // audit tab (ListAuditHistory is nil in centymo today — renders empty)
}

// PageData holds the data for the subscription group detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Group           *subscriptiongrouppb.SubscriptionGroup
	Labels          subscription_group.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID            string
	Name          string
	Kind          string
	PlanName      string
	ScheduleName  string
	Capacity      string
	Status        string
	StatusLabel   string
	StatusVariant string
	CreatedDate   string
	ModifiedDate  string

	// Subscriptions tab (the section roster)
	Subscriptions *types.TableConfig
	// Teaching-staff tab (M4: the class-row list, plan.md §2 row-source flip)
	Staff *SectionSGPPTabData
	// Attachments tab
	AttachmentTable *types.TableConfig
	// Audit tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// NewView creates the subscription group detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("subscription_group", "read") {
			return view.Forbidden("subscription_group:read")
		}
		id := viewCtx.Request.PathValue("id")

		activeTab := deps.Labels.Tabs.CanonicalizeTab(viewCtx.Request.URL.Query().Get("tab"))
		if activeTab == "" {
			activeTab = "info"
		}

		pageData, err := buildPageData(ctx, deps, id, activeTab, viewCtx)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("subscription-group-detail", pageData)
	})
}

// NewTabAction handles GET /action/subscription-group/{id}/tab/{tab}.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := deps.Labels.Tabs.CanonicalizeTab(viewCtx.Request.PathValue("tab"))
		// The Teaching-staff assign drawer rides this already-registered GET route
		// under the reserved staff-assign token (no new route): the per-row action
		// GETs it into #sheetContent. Dispatch before the tab-body path.
		if tab == staffAssignTab {
			return renderAssignDrawer(ctx, deps, viewCtx, id)
		}
		if tab == "" {
			tab = "info"
		}
		pageData, err := buildPageData(ctx, deps, id, tab, viewCtx)
		if err != nil {
			return view.Error(err)
		}
		// attachments + audit reuse the pyeza global blocks; every other tab is a
		// module-defined {{define "subscription-group-tab-<tab>"}}.
		templateName := "subscription-group-tab-" + tab
		switch tab {
		case "attachments":
			templateName = "attachment-tab"
		case "audit":
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

func buildPageData(ctx context.Context, deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext) (*PageData, error) {
	resp, err := deps.ReadSubscriptionGroup(ctx, &subscriptiongrouppb.ReadSubscriptionGroupRequest{
		Data: &subscriptiongrouppb.SubscriptionGroup{Id: id},
	})
	if err != nil {
		log.Printf("Failed to read subscription group %s: %v", id, err)
		return nil, fmt.Errorf("%s", deps.Labels.Errors.LoadFailed)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("%s", deps.Labels.Errors.NotFound)
	}
	sg := data[0]

	l := deps.Labels

	planName := l.Detail.NoPlan
	if pid := sg.GetPlanId(); pid != "" {
		if n := lookupPlanName(ctx, deps, pid); n != "" {
			planName = n
		} else {
			planName = pid
		}
	}
	scheduleName := l.Detail.NoSchedule
	if sid := sg.GetPriceScheduleId(); sid != "" {
		if n := lookupScheduleName(ctx, deps, sid); n != "" {
			scheduleName = n
		} else {
			scheduleName = sid
		}
	}

	kind := sg.GetKind()
	if kind == "" {
		kind = l.Detail.NoKind
	}

	status := "active"
	statusVariant := "success"
	// Badge Value renders verbatim — use the lyngua status labels, not the raw key.
	statusLabel := deps.CommonLabels.Status.Active
	if !sg.GetActive() {
		status = "inactive"
		statusVariant = "warning"
		statusLabel = deps.CommonLabels.Status.Inactive
	}

	base := route.ResolveURL(deps.Routes.DetailURL, "id", id)
	action := route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "")
	// Key stays canonical ("subscriptions"); the URL slug is lyngua-fied per tier
	// (education → "enrollments") via ResolveTabSlug. CanonicalizeTab (in NewView /
	// NewTabAction) maps the slug back so dispatch + template lookups stay canonical.
	subsSlug := l.Tabs.ResolveTabSlug("subscriptions")
	// Enrollments (subscriptions) tab count badge — mirrors the roster table's
	// data source + permission gate exactly so the badge and row counts cannot
	// drift. 0 renders no badge (unpermitted / unwired / empty section).
	enrollmentCount := countSectionEnrollments(ctx, deps, id)
	// Teaching-staff tab count badge — the number of non-excluded class rows
	// (M4 row-source flip), mirroring the tab's own row source + read gate
	// exactly (see countSectionSGPPs) so the badge and row counts cannot
	// drift. 0 renders no badge (unpermitted / unwired / no classes yet).
	offeringCount := countSectionSGPPs(ctx, deps, sg)
	tabItems := []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "subscriptions", Label: l.Tabs.Subscriptions, Href: base + "?tab=" + subsSlug, HxGet: action + subsSlug, Icon: "icon-users", Count: enrollmentCount},
		{Key: "staff", Label: l.Tabs.Staff, Href: base + "?tab=staff", HxGet: action + "staff", Icon: "icon-user-check", Count: offeringCount},
		{Key: "attachments", Label: l.Tabs.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		{Key: "audit", Label: l.Tabs.Audit, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
	}

	tz := types.LocationFromContext(ctx)
	createdDate := ""
	if ms := sg.GetDateCreated(); ms > 0 {
		createdDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}
	modifiedDate := ""
	if ms := sg.GetDateModified(); ms > 0 {
		modifiedDate = types.FormatInTZ(time.UnixMilli(ms), tz, types.DateTimeReadable)
	}

	headerSubtitle := strings.TrimSpace(planName)
	if headerSubtitle == "" || planName == l.Detail.NoPlan {
		headerSubtitle = l.Detail.NoSubtitle
	}

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          sg.GetName(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    sg.GetName(),
			HeaderSubtitle: headerSubtitle,
			HeaderIcon:     "icon-users",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "subscription-group-detail-content",
		Group:           sg,
		Labels:          l,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		Name:            sg.GetName(),
		Kind:            kind,
		PlanName:        planName,
		ScheduleName:    scheduleName,
		Capacity:        formatCapacity(sg, l),
		Status:          status,
		StatusLabel:     statusLabel,
		StatusVariant:   statusVariant,
		CreatedDate:     createdDate,
		ModifiedDate:    modifiedDate,
	}

	// Load the active tab's payload on demand (the info fields above are always set).
	switch activeTab {
	case "subscriptions":
		pageData.Subscriptions = buildSubscriptionsTable(ctx, deps, id, l)
	case "staff":
		pageData.Staff = buildSGPPTabData(ctx, deps, sg, l)
	case "attachments":
		if deps.ListAttachments != nil {
			cfg := attachmentConfig(deps)
			resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
			if err != nil {
				log.Printf("Failed to list attachments for subscription_group %s: %v", id, err)
			}
			var items []*attachmentpb.Attachment
			if resp != nil {
				items = resp.GetData()
			}
			pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
		}
	case "audit":
		if deps.ListAuditHistory != nil {
			cursor := viewCtx.Request.URL.Query().Get("cursor")
			auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
				EntityType:  "subscription_group",
				EntityID:    id,
				Limit:       20,
				CursorToken: cursor,
			})
			if err != nil {
				log.Printf("Failed to load audit history for subscription_group %s: %v", id, err)
			}
			if auditResp != nil {
				pageData.AuditEntries = auditResp.Entries
				pageData.AuditHasNext = auditResp.HasNext
				pageData.AuditNextCursor = auditResp.NextCursor
			}
		}
		pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit"
	}

	return pageData, nil
}

// formatCapacity renders the capacity summary from capacity_mode (+ max_capacity
// when CAPPED). UNSPECIFIED is treated as Unlimited per the proto contract.
func formatCapacity(sg *subscriptiongrouppb.SubscriptionGroup, l subscription_group.Labels) string {
	switch sg.GetCapacityMode() {
	case subscriptiongrouppb.CapacityMode_CAPACITY_MODE_CAPPED:
		return fmt.Sprintf(l.Detail.CapacityValue, sg.GetMaxCapacity())
	case subscriptiongrouppb.CapacityMode_CAPACITY_MODE_CLOSED:
		return l.Form.CapClosed
	case subscriptiongrouppb.CapacityMode_CAPACITY_MODE_UNLIMITED:
		return l.Form.CapUnlimited
	default:
		return l.Detail.CapacityModeNF
	}
}

func lookupPlanName(ctx context.Context, deps *DetailViewDeps, planID string) string {
	if deps.ListPlans == nil || planID == "" {
		return ""
	}
	resp, err := deps.ListPlans(ctx, &planpb.ListPlansRequest{})
	if err != nil {
		return ""
	}
	for _, p := range resp.GetData() {
		if p.GetId() == planID {
			return p.GetName()
		}
	}
	return ""
}

func lookupScheduleName(ctx context.Context, deps *DetailViewDeps, scheduleID string) string {
	if deps.ListPriceSchedules == nil || scheduleID == "" {
		return ""
	}
	resp, err := deps.ListPriceSchedules(ctx, &priceschedulepb.ListPriceSchedulesRequest{})
	if err != nil {
		return ""
	}
	for _, ps := range resp.GetData() {
		if ps.GetId() == scheduleID {
			return ps.GetName()
		}
	}
	return ""
}

// sectionMemberFilter scopes a subscription_group_member list to one section.
// Shared by the roster table body and the enrollments-count badge so the two
// query identically and their counts can never drift.
func sectionMemberFilter(groupID string) *commonpb.FilterRequest {
	return &commonpb.FilterRequest{
		Filters: []*commonpb.TypedFilter{{
			Field: "subscription_group_id",
			FilterType: &commonpb.TypedFilter_StringFilter{
				StringFilter: &commonpb.StringFilter{Value: groupID, Operator: commonpb.StringOperator_STRING_EQUALS},
			},
		}},
	}
}

// sectionMemberFilterInactive is sectionMemberFilter plus an explicit
// active=false term — the second half of the status-agnostic two-call pattern
// (see shared.InactiveFilter).
func sectionMemberFilterInactive(groupID string) *commonpb.FilterRequest {
	f := sectionMemberFilter(groupID)
	f.Logic = commonpb.FilterLogic_AND
	f.Filters = append(f.Filters, &commonpb.TypedFilter{
		Field: "active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: false},
		},
	})
	return f
}

// listSectionMembers fetches the section roster STATUS-AGNOSTICALLY: one bare
// call (the List default filters active=true) + one explicit active=false
// call, merged by id. An inactive section's historical roster is entirely
// inactive member rows — the bare call alone renders it empty. Shared by the
// roster table and the count badge so the two can never drift.
func listSectionMembers(ctx context.Context, deps *DetailViewDeps, groupID string) ([]*subscriptiongroupmemberpb.SubscriptionGroupMember, error) {
	var members []*subscriptiongroupmemberpb.SubscriptionGroupMember
	seen := map[string]bool{}
	for _, req := range []*subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest{
		{Filters: sectionMemberFilter(groupID)},
		{Filters: sectionMemberFilterInactive(groupID)},
	} {
		resp, err := deps.ListSubscriptionGroupMembers(ctx, req)
		if err != nil {
			return nil, err
		}
		for _, m := range resp.GetData() {
			if m == nil || seen[m.GetId()] {
				continue
			}
			seen[m.GetId()] = true
			members = append(members, m)
		}
	}
	return members, nil
}

// countSectionEnrollments returns the number of enrollments (subscription_group_
// member rows) in the section for the Enrollments tab count badge. It mirrors
// buildSubscriptionsTable's data source and permission gate exactly, so the
// badge count matches the rendered row count. Returns 0 (no badge) when the
// roster is not permitted, the dep is unwired, or the section is empty.
func countSectionEnrollments(ctx context.Context, deps *DetailViewDeps, groupID string) int {
	perms := view.GetUserPermissions(ctx)
	if perms == nil || !perms.Can("subscription_group_member", "list") || deps.ListSubscriptionGroupMembers == nil {
		return 0
	}
	members, err := listSectionMembers(ctx, deps, groupID)
	if err != nil {
		log.Printf("Failed to count members for subscription_group %s: %v", groupID, err)
		return 0
	}
	return len(members)
}

// buildSubscriptionsTable renders the section roster: one row per
// subscription_group_member, resolving client + subscription display names via a
// single batch fetch each (never per-row). Fail-closed on
// subscription_group_member:list INSIDE the tab body (empty table, not a full-page
// view.Forbidden — which would be wrong for an HTMX tab swap).
func buildSubscriptionsTable(ctx context.Context, deps *DetailViewDeps, groupID string, l subscription_group.Labels) *types.TableConfig {
	perms := view.GetUserPermissions(ctx)
	columns := []types.TableColumn{
		// Cap the name column at half the table so the subscription column
		// gets enough room to render "Last, First (AY …)" without wrapping.
		{Key: "client", Label: l.Columns.Client, NoSort: true, NoFilter: true, Width: "50%"},
		{Key: "subscription", Label: l.Columns.Subscription, NoSort: true, NoFilter: true},
		{Key: "status", Label: l.Columns.Status, NoSort: true, NoFilter: true, WidthClass: "col-2xl"},
	}
	cfg := &types.TableConfig{
		ID:          "subscription-group-subscriptions-table",
		Columns:     columns,
		Rows:        []types.TableRow{},
		Labels:      deps.TableLabels,
		EmptyState:  types.TableEmptyState{Title: l.Empty.Title, Message: l.Empty.Message},
		ShowSearch:  true,
		ShowColumns: true,
		ShowDensity: true,
		ShowEntries: true,
	}

	if perms == nil || !perms.Can("subscription_group_member", "list") || deps.ListSubscriptionGroupMembers == nil {
		types.ApplyTableSettings(cfg)
		return cfg
	}

	members, err := listSectionMembers(ctx, deps, groupID)
	if err != nil {
		log.Printf("Failed to list members for subscription_group %s: %v", groupID, err)
		types.ApplyTableSettings(cfg)
		return cfg
	}

	subCodes := resolveSubscriptionCodes(ctx, deps, members)

	// Banding is enabled only when Roster.GroupByField parses to a client-
	// attribute code. Unset/foreign → today's flat roster, byte-identical.
	groupCode, banded := deps.Options.Roster.GroupByAttributeCode()

	if !banded {
		clientNames := resolveClientNames(ctx, deps, members)
		rows := make([]types.TableRow, 0, len(members))
		for _, m := range members {
			client := clientNames[m.GetClientId()]
			if client == "" {
				client = m.GetClientId()
			}
			rows = append(rows, rosterRow(deps, m, client, subCodes))
		}
		cfg.Rows = rows
		types.ApplyColumnStyles(columns, rows)
		types.ApplyTableSettings(cfg)
		return cfg
	}

	// Optioned path: "Last, First" display, last-name sort within bands, and
	// value bands via the shared pyeza partition helper. TRAP: a row's ID is the
	// subscription_group_member id, but attribute values key by CLIENT id — the
	// value-by-row map is rebuilt to the member identity below.
	records := resolveClientRecords(ctx, deps, members)
	rows := make([]types.TableRow, 0, len(members))
	for _, m := range members {
		client := records[m.GetClientId()].listName()
		if client == "" {
			client = m.GetClientId()
		}
		rows = append(rows, rosterRow(deps, m, client, subCodes))
	}
	sortRosterRows(rows, members, records, deps.Options.Roster)

	attrByClient := fetchAttributeValues(ctx, deps, members, records)[groupCode]
	if len(attrByClient) == 0 {
		// Sort applied, but no attribute values resolved → flat sorted roster
		// (the fail-safe: banding degrades, never errors).
		cfg.Rows = rows
		types.ApplyColumnStyles(columns, rows)
		types.ApplyTableSettings(cfg)
		return cfg
	}
	valueByRowID := make(map[string]string, len(members))
	for _, m := range members {
		valueByRowID[m.GetId()] = attrByClient[m.GetClientId()]
	}
	cfg.Groups = types.GroupRowsByValue(rows, valueByRowID, types.GroupRowsByValueOptions{
		LeadingOrder: deps.Options.Roster.GroupValueOrder,
		GroupID:      func(v string) string { return "sg-band-" + slug(v) },
	})
	types.ApplyColumnStyles(columns, rows)
	types.ApplyTableSettings(cfg)
	return cfg
}

// rosterRow builds one roster TableRow (client + subscription + status) for a
// member. ID is the member id; client is the pre-resolved display string.
func rosterRow(deps *DetailViewDeps, m *subscriptiongroupmemberpb.SubscriptionGroupMember, client string, subCodes map[string]string) types.TableRow {
	sub := subCodes[m.GetSubscriptionId()]
	if sub == "" {
		sub = m.GetSubscriptionId()
	}
	// Badge Value renders verbatim — use the lyngua status labels, not the raw key.
	st, variant := deps.CommonLabels.Status.Active, "success"
	if !m.GetActive() {
		st, variant = deps.CommonLabels.Status.Inactive, "warning"
	}
	return types.TableRow{
		ID: m.GetId(),
		Cells: []types.TableCell{
			{Type: "text", Value: client},
			{Type: "text", Value: sub},
			{Type: "badge", Value: st, Variant: variant},
		},
	}
}

// idListFilter builds an `id IN (…)` filter — resolves EXACTLY the roster's
// referenced rows, not a paginated first page. LIST_IN → SQL IN in the postgres
// adapter (operations.go buildListFilter); precedent: product/list/page.go.
func idListFilter(ids []string) *commonpb.FilterRequest {
	return &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{
		Field: "id",
		FilterType: &commonpb.TypedFilter_ListFilter{
			ListFilter: &commonpb.ListFilter{Values: ids, Operator: commonpb.ListOperator_LIST_IN},
		},
	}}}
}

// idListFilterInactive is idListFilter plus an explicit active=false term —
// the second half of the status-agnostic two-call pattern (see
// shared.InactiveFilter). Roster display-name maps must resolve rows the List
// default would drop: an inactive section's enrollments reference inactive
// subscription rows.
func idListFilterInactive(ids []string) *commonpb.FilterRequest {
	f := idListFilter(ids)
	f.Logic = commonpb.FilterLogic_AND
	f.Filters = append(f.Filters, &commonpb.TypedFilter{
		Field: "active",
		FilterType: &commonpb.TypedFilter_BooleanFilter{
			BooleanFilter: &commonpb.BooleanFilter{Value: false},
		},
	})
	return f
}

// rosterPageLimit chunks id-list hydration reads. The espyna generic list
// silently defaults to LIMIT 100, and a section roster can hold up to 100 active
// + 100 inactive members, so unchunked `id IN (…)` reads would drop rows past
// the first 100 and mis-band them. Mirrors fayna's section-grid pageLimit.
const rosterPageLimit = 100

// chunkIDs splits ids into batches of at most size (size<=0 → one batch).
func chunkIDs(ids []string, size int) [][]string {
	if size <= 0 || len(ids) <= size {
		return [][]string{ids}
	}
	var out [][]string
	for start := 0; start < len(ids); start += size {
		end := start + size
		if end > len(ids) {
			end = len(ids)
		}
		out = append(out, ids[start:end])
	}
	return out
}

// resolveClientNames maps client_id → display name for exactly the roster's
// clients (LIST_IN by id, both statuses), falling back to first+last when the
// name is blank.
func resolveClientNames(ctx context.Context, deps *DetailViewDeps, members []*subscriptiongroupmemberpb.SubscriptionGroupMember) map[string]string {
	if deps.ListClients == nil {
		return nil
	}
	seen := map[string]bool{}
	var ids []string
	for _, m := range members {
		if id := m.GetClientId(); id != "" && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	m := make(map[string]string, len(ids))
	for _, f := range []*commonpb.FilterRequest{idListFilter(ids), idListFilterInactive(ids)} {
		resp, err := deps.ListClients(ctx, &clientpb.ListClientsRequest{Filters: f})
		if err != nil {
			log.Printf("Failed to list clients for subscription_group roster: %v", err)
			continue
		}
		for _, c := range resp.GetData() {
			name := strings.TrimSpace(c.GetName())
			if name == "" {
				name = strings.TrimSpace(c.GetFirstName() + " " + c.GetLastName())
			}
			m[c.GetId()] = name
		}
	}
	return m
}

// resolveSubscriptionCodes maps subscription_id → code for exactly the roster's
// subscriptions (LIST_IN by id, both statuses).
func resolveSubscriptionCodes(ctx context.Context, deps *DetailViewDeps, members []*subscriptiongroupmemberpb.SubscriptionGroupMember) map[string]string {
	if deps.ListSubscriptions == nil {
		return nil
	}
	seen := map[string]bool{}
	var ids []string
	for _, m := range members {
		if id := m.GetSubscriptionId(); id != "" && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	m := make(map[string]string, len(ids))
	for _, f := range []*commonpb.FilterRequest{idListFilter(ids), idListFilterInactive(ids)} {
		resp, err := deps.ListSubscriptions(ctx, &subscriptionpb.ListSubscriptionsRequest{Filters: f})
		if err != nil {
			log.Printf("Failed to list subscriptions for subscription_group roster: %v", err)
			continue
		}
		for _, s := range resp.GetData() {
			m[s.GetId()] = s.GetCode()
		}
	}
	return m
}

// clientRecord is a roster client's structured identity for the optioned path —
// the display name plus the last/first parts the "Last, First" render and the
// last-name sort need.
type clientRecord struct {
	name      string
	lastName  string
	firstName string
}

// listName renders the class-list name form "{last_name}, {first_name}",
// falling back to the plain display name when either part is missing.
func (r clientRecord) listName() string {
	if r.lastName != "" && r.firstName != "" {
		return r.lastName + ", " + r.firstName
	}
	return r.name
}

// resolveClientRecords maps client_id → structured name for exactly the roster's
// clients (LIST_IN by id, both statuses). The optioned counterpart to
// resolveClientNames; the flat path keeps the plain display string.
func resolveClientRecords(ctx context.Context, deps *DetailViewDeps, members []*subscriptiongroupmemberpb.SubscriptionGroupMember) map[string]clientRecord {
	if deps.ListClients == nil {
		return nil
	}
	seen := map[string]bool{}
	var ids []string
	for _, m := range members {
		if id := m.GetClientId(); id != "" && !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	out := make(map[string]clientRecord, len(ids))
	for _, batch := range chunkIDs(ids, rosterPageLimit) {
		for _, f := range []*commonpb.FilterRequest{idListFilter(batch), idListFilterInactive(batch)} {
			resp, err := deps.ListClients(ctx, &clientpb.ListClientsRequest{Filters: f})
			if err != nil {
				log.Printf("Failed to list clients for subscription_group roster: %v", err)
				continue
			}
			for _, c := range resp.GetData() {
				out[c.GetId()] = clientRecord{
					name:      clientDisplayName(c),
					lastName:  clientLastName(c),
					firstName: clientFirstName(c),
				}
			}
		}
	}
	return out
}

// clientDisplayName prefers the client's own name column, then first+last, then
// the embedded User's first+last, then the id.
func clientDisplayName(c *clientpb.Client) string {
	if name := strings.TrimSpace(c.GetName()); name != "" {
		return name
	}
	if fn := strings.TrimSpace(c.GetFirstName() + " " + c.GetLastName()); fn != "" {
		return fn
	}
	if u := c.GetUser(); u != nil {
		if name := strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName()); name != "" {
			return name
		}
	}
	return c.GetId()
}

// clientLastName prefers the client's own last_name column, then the embedded
// User's last name.
func clientLastName(c *clientpb.Client) string {
	if ln := strings.TrimSpace(c.GetLastName()); ln != "" {
		return ln
	}
	if u := c.GetUser(); u != nil {
		if ln := strings.TrimSpace(u.GetLastName()); ln != "" {
			return ln
		}
	}
	return ""
}

// clientFirstName mirrors clientLastName for the first-name column.
func clientFirstName(c *clientpb.Client) string {
	if fn := strings.TrimSpace(c.GetFirstName()); fn != "" {
		return fn
	}
	if u := c.GetUser(); u != nil {
		if fn := strings.TrimSpace(u.GetFirstName()); fn != "" {
			return fn
		}
	}
	return ""
}

// fetchAttributeValues resolves each Roster-referenced attribute code to its id,
// then loads the roster clients' values — code → (client_id → value). Nil-safe:
// an unwired closure or unresolvable code yields no values (the roster then
// renders sorted-flat, never errors).
//
// SECURITY: the client set is the INTERSECTION of the roster's member client_ids
// with `authorized` — the clients the workspace-scoped resolveClientRecords
// actually returned. client_attribute is an EAV table with no workspace_id
// column, so the list wrapper cannot scope it; deriving the client set straight
// from member.client_id would let a membership carrying a foreign (other-
// workspace) client_id disclose that client's attribute value. Values are also
// re-checked against `authorized` before acceptance (defense in depth).
func fetchAttributeValues(ctx context.Context, deps *DetailViewDeps, members []*subscriptiongroupmemberpb.SubscriptionGroupMember, authorized map[string]clientRecord) map[string]map[string]string {
	codes := deps.Options.AttributeCodes()
	if len(codes) == 0 || deps.ListClientAttributes == nil || deps.ResolveAttributeIDByCode == nil {
		return nil
	}
	seen := map[string]bool{}
	var clientIDs []string
	for _, m := range members {
		id := m.GetClientId()
		if id == "" || seen[id] {
			continue
		}
		if _, ok := authorized[id]; !ok {
			continue // unresolved / foreign client → never queried, never banded
		}
		seen[id] = true
		clientIDs = append(clientIDs, id)
	}
	if len(clientIDs) == 0 {
		return nil
	}
	out := make(map[string]map[string]string, len(codes))
	for _, code := range codes {
		attrID, err := deps.ResolveAttributeIDByCode(ctx, code)
		if err != nil || attrID == "" {
			log.Printf("subscription_group roster: attribute code %q did not resolve (bands ignored for it): %v", code, err)
			continue
		}
		vals := map[string]string{}
		for _, batch := range chunkIDs(clientIDs, rosterPageLimit) {
			resp, err := deps.ListClientAttributes(ctx, &clientattributepb.ListClientAttributesRequest{
				Filters: attributeValueFilter(attrID, batch),
			})
			if err != nil {
				log.Printf("subscription_group roster: list client attributes for %q: %v", code, err)
				continue
			}
			for _, ca := range resp.GetData() {
				cid, v := ca.GetClientId(), strings.TrimSpace(ca.GetValue())
				if cid == "" || v == "" {
					continue
				}
				if _, ok := authorized[cid]; !ok {
					continue // defense in depth: never surface a non-authorized client's value
				}
				vals[cid] = v
			}
		}
		out[code] = vals
	}
	return out
}

// attributeValueFilter scopes a client_attribute list to one attribute across
// the roster's clients: attribute_id EQ AND client_id IN. LIST_IN → SQL IN, the
// same resolve-exactly-the-roster idiom as idListFilter.
func attributeValueFilter(attrID string, clientIDs []string) *commonpb.FilterRequest {
	return &commonpb.FilterRequest{
		Logic: commonpb.FilterLogic_AND,
		Filters: []*commonpb.TypedFilter{
			{
				Field: "attribute_id",
				FilterType: &commonpb.TypedFilter_StringFilter{
					StringFilter: &commonpb.StringFilter{Value: attrID, Operator: commonpb.StringOperator_STRING_EQUALS},
				},
			},
			{
				Field: "client_id",
				FilterType: &commonpb.TypedFilter_ListFilter{
					ListFilter: &commonpb.ListFilter{Values: clientIDs, Operator: commonpb.ListOperator_LIST_IN},
				},
			},
		},
	}
}

// sortRosterRows orders roster rows by the client's last name (the single
// implemented SortField), direction-aware and stable; ties fall back to the
// display name then the row id. Rows key by member id, records by client id, so
// the member→client map bridges the two identities.
func sortRosterRows(rows []types.TableRow, members []*subscriptiongroupmemberpb.SubscriptionGroupMember, records map[string]clientRecord, opts subscription_group.RowOptions) {
	if !opts.SortByLastName() {
		return
	}
	clientByMember := make(map[string]string, len(members))
	for _, m := range members {
		clientByMember[m.GetId()] = m.GetClientId()
	}
	desc := opts.Direction() == "desc"
	sort.SliceStable(rows, func(i, j int) bool {
		ri := records[clientByMember[rows[i].ID]]
		rj := records[clientByMember[rows[j].ID]]
		a, b := strings.ToLower(ri.lastName), strings.ToLower(rj.lastName)
		if a == b {
			// Equal last names → first name (the locked last→first contract),
			// then display name, then row ID as deterministic tie-breakers.
			fa, fb := strings.ToLower(ri.firstName), strings.ToLower(rj.firstName)
			if fa != fb {
				return fa < fb
			}
			an, bn := strings.ToLower(ri.name), strings.ToLower(rj.name)
			if an != bn {
				return an < bn
			}
			return rows[i].ID < rows[j].ID
		}
		// Values present sort before empties regardless of direction.
		if (a == "") != (b == "") {
			return a != ""
		}
		if desc {
			return a > b
		}
		return a < b
	})
}

// slug reduces a band value to a stable id token (lowercase alnum; spaces/
// dashes/underscores → "-"; empty → "none"). Keeps sg-band-* ids reproducible.
func slug(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '_':
			b.WriteByte('-')
		}
	}
	if b.Len() == 0 {
		return "none"
	}
	return b.String()
}
