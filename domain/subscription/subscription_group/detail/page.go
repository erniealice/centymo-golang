package detail

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	subscription_group "github.com/erniealice/centymo-golang/domain/subscription/subscription_group"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
	priceschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_schedule"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
	subscriptiongrouppb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group"
	subscriptiongroupmemberpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription_group_member"
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
	StatusVariant string
	CreatedDate   string
	ModifiedDate  string

	// Subscriptions tab (the section roster)
	Subscriptions *types.TableConfig
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
	if !sg.GetActive() {
		status = "inactive"
		statusVariant = "warning"
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
	tabItems := []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "subscriptions", Label: l.Tabs.Subscriptions, Href: base + "?tab=" + subsSlug, HxGet: action + subsSlug, Icon: "icon-users", Count: enrollmentCount},
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
		StatusVariant:   statusVariant,
		CreatedDate:     createdDate,
		ModifiedDate:    modifiedDate,
	}

	// Load the active tab's payload on demand (the info fields above are always set).
	switch activeTab {
	case "subscriptions":
		pageData.Subscriptions = buildSubscriptionsTable(ctx, deps, id, l)
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
	resp, err := deps.ListSubscriptionGroupMembers(ctx, &subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest{
		Filters: sectionMemberFilter(groupID),
	})
	if err != nil {
		log.Printf("Failed to count members for subscription_group %s: %v", groupID, err)
		return 0
	}
	return len(resp.GetData())
}

// buildSubscriptionsTable renders the section roster: one row per
// subscription_group_member, resolving client + subscription display names via a
// single batch fetch each (never per-row). Fail-closed on
// subscription_group_member:list INSIDE the tab body (empty table, not a full-page
// view.Forbidden — which would be wrong for an HTMX tab swap).
func buildSubscriptionsTable(ctx context.Context, deps *DetailViewDeps, groupID string, l subscription_group.Labels) *types.TableConfig {
	perms := view.GetUserPermissions(ctx)
	columns := []types.TableColumn{
		{Key: "client", Label: l.Columns.Client, NoSort: true, NoFilter: true},
		{Key: "subscription", Label: l.Columns.Subscription, NoSort: true, NoFilter: true, WidthClass: "col-2xl"},
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

	resp, err := deps.ListSubscriptionGroupMembers(ctx, &subscriptiongroupmemberpb.ListSubscriptionGroupMembersRequest{
		Filters: sectionMemberFilter(groupID),
	})
	if err != nil {
		log.Printf("Failed to list members for subscription_group %s: %v", groupID, err)
		types.ApplyTableSettings(cfg)
		return cfg
	}
	members := resp.GetData()

	clientNames := resolveClientNames(ctx, deps, members)
	subCodes := resolveSubscriptionCodes(ctx, deps, members)

	rows := make([]types.TableRow, 0, len(members))
	for _, m := range members {
		client := clientNames[m.GetClientId()]
		if client == "" {
			client = m.GetClientId()
		}
		sub := subCodes[m.GetSubscriptionId()]
		if sub == "" {
			sub = m.GetSubscriptionId()
		}
		st, variant := "active", "success"
		if !m.GetActive() {
			st, variant = "inactive", "warning"
		}
		rows = append(rows, types.TableRow{
			ID: m.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: client},
				{Type: "text", Value: sub},
				{Type: "badge", Value: st, Variant: variant},
			},
		})
	}
	cfg.Rows = rows
	types.ApplyColumnStyles(columns, rows)
	types.ApplyTableSettings(cfg)
	return cfg
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

// resolveClientNames maps client_id → display name for exactly the roster's
// clients (LIST_IN by id), falling back to first+last when the name is blank.
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
	resp, err := deps.ListClients(ctx, &clientpb.ListClientsRequest{Filters: idListFilter(ids)})
	if err != nil {
		log.Printf("Failed to list clients for subscription_group roster: %v", err)
		return nil
	}
	m := make(map[string]string, len(resp.GetData()))
	for _, c := range resp.GetData() {
		name := strings.TrimSpace(c.GetName())
		if name == "" {
			name = strings.TrimSpace(c.GetFirstName() + " " + c.GetLastName())
		}
		m[c.GetId()] = name
	}
	return m
}

// resolveSubscriptionCodes maps subscription_id → code for exactly the roster's
// subscriptions (LIST_IN by id).
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
	resp, err := deps.ListSubscriptions(ctx, &subscriptionpb.ListSubscriptionsRequest{Filters: idListFilter(ids)})
	if err != nil {
		log.Printf("Failed to list subscriptions for subscription_group roster: %v", err)
		return nil
	}
	m := make(map[string]string, len(resp.GetData()))
	for _, s := range resp.GetData() {
		m[s.GetId()] = s.GetCode()
	}
	return m
}
