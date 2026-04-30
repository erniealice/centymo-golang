package detail

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	centymo "github.com/erniealice/centymo-golang"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	enums "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/enums"
	jobpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job"
	jobphasepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_phase"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	billingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billing_event"
	priceplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/price_plan"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes           centymo.SubscriptionRoutes
	ReadSubscription func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)

	// 2026-04-29 milestone-billing — list events for the Package tab. nil-safe.
	ListBillingEventsBySubscription func(ctx context.Context, req *billingeventpb.ListBillingEventsBySubscriptionRequest) (*billingeventpb.ListBillingEventsBySubscriptionResponse, error)

	// 2026-04-29 auto-spawn-jobs-from-subscription Phase D — Operations tab
	// data ops. nil-safe; tab degrades to empty state.
	GetJobsByOrigin   func(ctx context.Context, req *jobpb.GetJobsByOriginRequest) (*jobpb.GetJobsByOriginResponse, error)
	ListJobPhasesByJob func(ctx context.Context, req *jobphasepb.ListJobPhasesByJobRequest) (*jobphasepb.ListJobPhasesByJobResponse, error)
	// JobDetailURL is the absolute URL pattern (e.g. /app/jobs/detail/{id})
	// used to deep-link to fayna's Job detail page from the Operations tab.
	// Empty means no link.
	JobDetailURL string
	// SpawnJobsURL is the centymo retroactive spawn drawer entry. Used by the
	// Operations tab empty-state CTA. Empty disables the CTA.
	SpawnJobsURL string
	// GetSubscriptionItemPageData returns the subscription with its joined
	// Client (+ User) and PricePlan (+ Plan) populated. Preferred over
	// ReadSubscription for the detail view so customer + package fields render
	// without extra round-trips.
	GetSubscriptionItemPageData func(ctx context.Context, req *subscriptionpb.GetSubscriptionItemPageDataRequest) (*subscriptionpb.GetSubscriptionItemPageDataResponse, error)
	// ReadClient is used by the nested-route variant ("under client") to look
	// up the client name for the page-header breadcrumb. Optional — when nil
	// the embedded client from GetSubscriptionItemPageData is used as the
	// fallback label, and the URL falls back to the flat subscription list.
	ReadClient func(ctx context.Context, req *clientpb.ReadClientRequest) (*clientpb.ReadClientResponse, error)
	// GetRevenueListPageData fetches revenue records (invoices) for the
	// invoices tab. Filtered by subscription_id. Optional — tab renders empty
	// state when nil.
	GetRevenueListPageData func(ctx context.Context, req *revenuepb.GetRevenueListPageDataRequest) (*revenuepb.GetRevenueListPageDataResponse, error)
	// ClientDetailURL is the absolute path template for the client detail page
	// (e.g. "/app/clients/detail/{id}"); used to build the breadcrumb link.
	// Empty string disables the breadcrumb link (label still renders).
	ClientDetailURL string
	Labels          centymo.SubscriptionLabels
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels

	attachment.AttachmentOps
	auditlog.AuditOps
}

// loadSubscriptionWithRelations fetches the subscription joined with client
// and price plan. Falls back to plain ReadSubscription when the page-data
// dep is unwired.
func loadSubscriptionWithRelations(ctx context.Context, deps *DetailViewDeps, id string) (*subscriptionpb.Subscription, error) {
	if deps.GetSubscriptionItemPageData != nil {
		resp, err := deps.GetSubscriptionItemPageData(ctx, &subscriptionpb.GetSubscriptionItemPageDataRequest{
			SubscriptionId: id,
		})
		if err != nil {
			return nil, err
		}
		if resp == nil || resp.GetSubscription() == nil {
			return nil, fmt.Errorf("subscription not found")
		}
		return resp.GetSubscription(), nil
	}
	resp, err := deps.ReadSubscription(ctx, &subscriptionpb.ReadSubscriptionRequest{
		Data: &subscriptionpb.Subscription{Id: id},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.GetData()) == 0 {
		return nil, fmt.Errorf("subscription not found")
	}
	return resp.GetData()[0], nil
}

// PageData holds the data for the subscription detail page.
type PageData struct {
	types.PageData
	ContentTemplate     string
	Subscription        map[string]any
	Labels              centymo.SubscriptionLabels
	ActiveTab           string
	TabItems            []pyeza.TabItem
	// Invoices tab
	Invoices        *types.TableConfig
	AttachmentTable     *types.TableConfig
	AttachmentUploadURL string
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string

	// 2026-04-27 plan-client-scope plan §6.5 — Package tab.
	// CTA shown on the package tab when pricePlan.client_id IS NULL.
	PackageCustomizeURL    string // POST endpoint for the customize CTA
	PackageCustomizeLabel  string // pre-resolved label with {{.ClientName}}
	PackageCustomizeShown  bool   // false hides the CTA (already client-scoped)
	PackageCustomizeDisabled bool // true grays the CTA (no permission, etc.)
	PackageClientName      string
	PackagePricePlan       map[string]any // {id, name, currency, amount, plan_id, client_id}

	// 2026-04-29 milestone-billing plan §5 / Phase D — Milestones section
	// inside Package tab. Rendered only when pricePlan.billing_kind = MILESTONE.
	MilestonesShown        bool
	Milestones             []MilestoneRow
	TotalInvoicedDisplay   string // formatted "₱430,000.00"
	MilestoneCurrency      string

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.2 / Phase D —
	// Operations tab data. OperationsHasJobs flips the empty-state CTA;
	// SpawnJobsURL drives the retroactive-spawn drawer.
	OperationsHasJobs   bool
	OperationsRootJobs  []OperationsJobRow
	OperationsEmptyText string // pre-resolved {{.SpawnAction}} substitution
	OperationsSpawnURL  string

	// 2026-04-30 cyclic-subscription-jobs plan §7 / Phase D — branch the
	// Operations tab when the subscription's PricePlan is cyclic
	// (RECURRING or CONTRACT-with-cycle). When IsCyclic is true the tab
	// renders the cycle accordion + spawn / backfill CTAs; the legacy flat
	// rendering above stays for non-cyclic engagements (regression guard
	// per progress.md "key gotchas").
	IsCyclic                  bool
	OperationsCyclic          *SubscriptionCyclesData
	SpawnCycleJobsURL         string // POST endpoint for "Spawn this cycle now"
	BackfillCycleJobsDrawerURL string // GET endpoint for the backfill drawer

	// 2026-04-30 cyclic-subscription-jobs plan §21 / Phase D — flat Jobs
	// tab. Hidden when no Jobs exist (Jobs.HasJobs == false).
	Jobs *SubscriptionJobsTabData
}

// SubscriptionCyclesData carries the cycle-accordion view rows for a cyclic
// subscription's Operations tab. Built by buildSubscriptionCyclesData per
// cyclic-subscription-jobs plan §7.1.
type SubscriptionCyclesData struct {
	// EngagementJob is the parent shell Job (parent_job_id == NULL for cyclic
	// subscriptions). Empty struct when the engagement hasn't been spawned
	// yet (legacy subscriptions created pre-this-plan).
	EngagementJobID   string
	EngagementName    string
	EngagementHeading string // pre-resolved {{.Started}} / {{.Name}}

	// OnceAtStartJobs are children of the engagement with cycle_index=NULL
	// (e.g. onboarding fired by JOB_TEMPLATE_RELATION_TYPE_ONCE_AT_ENGAGEMENT_START).
	OnceAtStartJobs []OperationsJobRow

	// Cycles holds one entry per cycle (sorted descending by CycleIndex so
	// the most-recent cycle is on top). Empty when no cycle Jobs exist yet.
	Cycles []SubscriptionCycleView

	// CycleEmpty is the lyngua-resolved "No cycles yet" string surfaced when
	// Cycles is empty.
	CycleEmpty string

	// MissingCycleCount is the number of cycle windows from sub.date_time_start
	// → today that have no Jobs yet. > 0 surfaces the backfill banner.
	MissingCycleCount int
	BackfillBannerText string // pre-resolved {{.Count}} substitution
}

// SubscriptionCycleView is one cycle accordion row. View-side only — no
// proto changes (cycle metadata read from Job.cycle_* nullable fields).
type SubscriptionCycleView struct {
	CycleIndex   int32
	PeriodStart  string                 // ISO 8601 (YYYY-MM-DD)
	PeriodEnd    string
	PeriodLabel  string                 // human-readable, "May 2026"
	HeadingText  string                 // pre-resolved cycleHeading template
	StatusKey    string                 // pending | inProgress | completed | overdue
	StatusLabel  string                 // lyngua-resolved
	StatusVariant string                // pyeza badge variant (success / warning / etc.)
	Jobs         []OperationsJobRow     // cycle Jobs (1 or N for multi-visit)
	InvoiceID    string                 // matched Revenue.id (empty if not yet recognized)
	InvoiceLabel string                 // pre-resolved cycleInvoiceLinked OR cycleNoInvoice
	IsPlaceholder bool                  // true for the next-un-spawned cycle (operator click → spawn)
	OpenByDefault bool                  // current cycle expanded; past cycles collapsed
}

// SubscriptionJobsTabData backs the new flat Jobs tab (plan §21.5).
type SubscriptionJobsTabData struct {
	Rows         []SubscriptionJobRow
	HasJobs      bool                   // tab hidden when false
	StatusCounts map[string]int         // for filter pill counts (by status key)
	TypeCounts   map[string]int         // for filter pill counts (by type key)
}

// SubscriptionJobRow is one row in the flat Jobs tab table.
type SubscriptionJobRow struct {
	JobID         string
	JobName       string
	JobType       string // engagement | cycle | onboarding | visit
	JobTypeLabel  string // lyngua-resolved
	Status        string
	StatusLabel   string
	StatusVariant string
	PeriodLabel   string // empty for engagement / onboarding
	CycleIndex    int32  // 0 for engagement / onboarding
	DetailURL     string // deep link to Job detail in fayna
}

// OperationsJobRow is one Job row rendered on the subscription detail's
// Operations tab. Children are rendered inline via the Children slice.
// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.2.
type OperationsJobRow struct {
	JobID            string
	JobName          string
	IsRoot           bool
	StatusKey        string // lowercase status, e.g. "planned"
	StatusVariant    string // pyeza badge variant (success / warning / etc.)
	BillingRuleKey   string // lowercase billing rule type for the badge
	PhaseSummaryText string // resolved "{Complete} / {Total} phases complete"
	JobDetailURL     string // empty when JobDetailURL dep is unwired
	Children         []OperationsJobRow
}

// MilestoneRow is a single BillingEvent row rendered inside the Package tab's
// Milestones section. flow.md §10 selectors:
//   - [data-testid="milestone-row"][data-event-id="ev-XXX"]
//   - [data-testid="milestone-status-{pending|ready|billed|...}"]
type MilestoneRow struct {
	EventID         string
	StatusKey       string // "pending" | "ready" | "billed" | "waived" | "deferred" | "cancelled"
	StatusLabel     string
	SequenceLabel   string
	BillableAmount  int64
	BillableDisplay string
	Currency        string
	RevenueID       string
	RevenueURL      string
	MarkReadyURL    string
	WaiveURL        string
	ShowMarkReady   bool
	ShowWaive       bool
	ShowRevenueLink bool
}

// subscriptionToMap converts a Subscription protobuf to a map[string]any for template use.
// Expects the subscription to have its joined Client (+ User) and PricePlan
// (+ Plan) populated — see loadSubscriptionWithRelations.
func subscriptionToMap(ctx context.Context, s *subscriptionpb.Subscription) map[string]any {
	// Customer = company name, falling back to representative full name.
	// Empty when the join is missing — never use the subscription's own name
	// here (it's a plan-derived label, not a customer label).
	customer := ""
	if c := s.GetClient(); c != nil {
		if companyName := c.GetName(); companyName != "" {
			customer = companyName
		} else if u := c.GetUser(); u != nil {
			first := u.GetFirstName()
			last := u.GetLastName()
			if first != "" || last != "" {
				customer = first + " " + last
			}
		}
	}

	// Get plan name from nested price_plan → plan
	planName := ""
	if pp := s.GetPricePlan(); pp != nil {
		if p := pp.GetPlan(); p != nil {
			planName = p.GetName()
		}
		if planName == "" {
			planName = pp.GetName()
		}
	}

	status := "active"
	if !s.GetActive() {
		status = "inactive"
	}

	tz := types.LocationFromContext(ctx)

	return map[string]any{
		"id":                   s.GetId(),
		"name":                 s.GetName(),
		"customer":             customer,
		"plan":                 planName,
		"price_plan_id":        s.GetPricePlanId(),
		"client_id":            s.GetClientId(),
		"date_start_string":    types.FormatTimestampInTZ(s.GetDateTimeStart(), tz, types.DateTimeReadable),
		"date_end_string":      types.FormatTimestampInTZ(s.GetDateTimeEnd(), tz, types.DateTimeReadable),
		"status":               status,
		"active":               s.GetActive(),
		"date_created_string":  s.GetDateCreatedString(),
		"date_modified_string": s.GetDateModifiedString(),
		"quantity":             s.GetQuantity(),
		"assigned_count":       s.GetAssignedCount(),
		"available_count":      s.GetAvailableCount(),
	}
}

// NewView creates the subscription detail view. Handles both the flat
// /app/subscriptions/detail/{id} URL and the nested
// /app/clients/detail/{client_id}/subscriptions/{id} URL — when the latter
// path param is set, the page-header renders a "client name → subscription"
// breadcrumb.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		clientIDFromPath := viewCtx.Request.PathValue("client_id")

		sub, err := loadSubscriptionWithRelations(ctx, deps, id)
		if err != nil {
			log.Printf("Failed to read subscription %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load subscription: %w", err))
		}
		subscription := subscriptionToMap(ctx, sub)
		breadcrumbLabel, breadcrumbURL := resolveClientBreadcrumb(ctx, deps, clientIDFromPath, sub)

		subName, _ := subscription["name"].(string)
		customer, _ := subscription["customer"].(string)
		dateStartStr, _ := subscription["date_start_string"].(string)
		dateEndStr, _ := subscription["date_end_string"].(string)

		// Header: title = subscription name; subtitle = customer · start[ — end].
		headerTitle := subName
		var subtitleParts []string
		if customer != "" {
			subtitleParts = append(subtitleParts, customer)
		}
		switch {
		case dateStartStr != "" && dateEndStr != "":
			subtitleParts = append(subtitleParts, dateStartStr+" — "+dateEndStr)
		case dateStartStr != "":
			subtitleParts = append(subtitleParts, dateStartStr)
		case dateEndStr != "":
			subtitleParts = append(subtitleParts, "until "+dateEndStr)
		}

		l := deps.Labels
		headerSubtitle := strings.Join(subtitleParts, " · ")
		if headerSubtitle == "" {
			headerSubtitle = l.Detail.PageTitle
		}

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}

		// 2026-04-30 cyclic-subscription-jobs plan §7 — branch the Operations
		// tab on the subscription's PricePlan billing kind. IsCyclic mirrors
		// espyna's eligibleForInstanceSpawn predicate (RECURRING OR
		// CONTRACT-with-cycle).
		isCyclic := computeIsCyclic(sub)

		// 2026-04-30 cyclic-subscription-jobs plan §21.3 — Jobs tab visibility
		// gate: shown only when at least one Job exists for this subscription.
		// One-shot upfront load so the gate matches what the tab itself will
		// render (no flicker when the operator navigates from Operations →
		// Jobs).
		allJobs := loadSubscriptionJobs(ctx, deps, id)
		jobsTabVisible := len(allJobs) > 0
		tabItems := buildTabItems(l, id, deps.Routes, jobsTabVisible)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:        viewCtx.CacheVersion,
				Title:               headerTitle,
				CurrentPath:         viewCtx.CurrentPath,
				ActiveNav:           deps.Routes.ActiveNav,
				ActiveSubNav:        deps.Routes.ActiveSubNav,
				HeaderTitle:         headerTitle,
				HeaderSubtitle:      headerSubtitle,
				HeaderIcon:          "icon-refresh-cw",
				HeaderBreadcrumb:    breadcrumbLabel,
				HeaderBreadcrumbURL: breadcrumbURL,
				CommonLabels:        deps.CommonLabels,
			},
			ContentTemplate: "subscription-detail-content",
			Subscription:    subscription,
			Labels:          l,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
			IsCyclic:        isCyclic,
			SpawnCycleJobsURL: strings.ReplaceAll(
				deps.Routes.SpawnCycleJobsURL, "{subscriptionId}", id),
			BackfillCycleJobsDrawerURL: strings.ReplaceAll(
				deps.Routes.BackfillCycleJobsURL, "{subscriptionId}", id),
		}

		// Inject the tab-content URL into the subscription map so the invoices
		// tab can refresh inline (HX-Trigger refresh-invoices listens here).
		subscription["tab_invoices_url"] = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "invoices"

		perms := view.GetUserPermissions(ctx)
		// nil perms = no restrictions (dev / mock mode). Consistent with
		// the comment on UserPermissions.Can() and most other view checks.
		canRecognize := perms == nil || perms.Can("revenue", "create")
		subscriptionActive, _ := subscription["active"].(bool)

		switch activeTab {
		case "package":
			canCustomize := perms != nil && (perms.Can("revenue", "create") || perms.Can("plan", "create"))
			customizeURL := route.ResolveURL(deps.Routes.CustomizePackageURL, "id", id)
			shown, disabled, label, clientName, ppData := buildPackageTabData(sub, customizeURL, l.Actions.CustomizePackage, canCustomize)
			pageData.PackageCustomizeURL = customizeURL
			pageData.PackageCustomizeShown = shown && subscriptionActive
			pageData.PackageCustomizeDisabled = disabled
			pageData.PackageCustomizeLabel = label
			pageData.PackageClientName = clientName
			pageData.PackagePricePlan = ppData
			// 2026-04-29 milestone-billing — milestones list (only on MILESTONE plans).
			applyMilestoneTabData(ctx, deps, sub, pageData, id)
		case "operations":
			applyOperationsTabData(ctx, deps, pageData, id, sub, allJobs, isCyclic)
		case "jobs":
			applyJobsTabData(ctx, deps, pageData, allJobs)
		case "invoices":
			revenues := loadSubscriptionInvoices(ctx, deps, id)
			pageData.Invoices = buildInvoicesTable(
				revenues, l, deps.TableLabels, centymo.RevenueDetailURL,
				deps.Routes.RecognizeURL, id,
				canRecognize, subscriptionActive,
				resolveRecognizeDisabledTooltip(canRecognize, subscriptionActive, l),
			)
		case "attachments":
			if deps.ListAttachments != nil {
				cfg := attachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
				if err != nil {
					log.Printf("Failed to list attachments: %v", err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)
		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.QueryParams["cursor"]
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "subscription",
					EntityID:    id,
					Limit:       20,
					CursorToken: cursor,
				})
				if err != nil {
					log.Printf("Failed to load audit history: %v", err)
				}
				if auditResp != nil {
					pageData.AuditEntries = auditResp.Entries
					pageData.AuditHasNext = auditResp.HasNext
					pageData.AuditNextCursor = auditResp.NextCursor
				}
			}
			pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
		}

		return view.OK("subscription-detail", pageData)
	})
}

// resolveRecognizeDisabledTooltip returns a sensible tooltip explaining why
// the recognize action is disabled. Empty string when the action is enabled.
//
// We re-use existing label keys (PermissionDenied, InvalidStatus) so this
// stays lyngua-driven without inventing a new key just for the disabled
// hover state.
func resolveRecognizeDisabledTooltip(canRecognize, active bool, l centymo.SubscriptionLabels) string {
	switch {
	case !canRecognize:
		return l.Errors.PermissionDenied
	case !active:
		return l.Errors.InvalidStatus
	default:
		return ""
	}
}

// resolveClientBreadcrumb returns the (label, href) pair for the page-header
// breadcrumb when the subscription is being viewed under a client context.
// Resolution order:
//   1. clientIDFromPath set and ReadClient available → live lookup (most reliable
//      because the joined client may be missing or stale).
//   2. clientIDFromPath set, no ReadClient → use the joined client's name if any.
//   3. clientIDFromPath empty → no breadcrumb (returns empty strings).
// The href points at the client's engagements tab so "back" lands where the
// operator clicked through from.
func resolveClientBreadcrumb(ctx context.Context, deps *DetailViewDeps, clientIDFromPath string, sub *subscriptionpb.Subscription) (string, string) {
	if clientIDFromPath == "" {
		return "", ""
	}
	label := ""
	if deps.ReadClient != nil {
		if resp, err := deps.ReadClient(ctx, &clientpb.ReadClientRequest{
			Data: &clientpb.Client{Id: clientIDFromPath},
		}); err == nil && len(resp.GetData()) > 0 {
			c := resp.GetData()[0]
			if name := c.GetName(); name != "" {
				label = name
			} else if u := c.GetUser(); u != nil {
				label = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			}
		}
	}
	if label == "" {
		if c := sub.GetClient(); c != nil {
			if name := c.GetName(); name != "" {
				label = name
			} else if u := c.GetUser(); u != nil {
				label = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			}
		}
	}
	if label == "" {
		label = clientIDFromPath
	}
	href := ""
	if deps.ClientDetailURL != "" {
		href = route.ResolveURL(deps.ClientDetailURL, "id", clientIDFromPath) + "?tab=engagements"
	}
	return label, href
}

func buildTabItems(l centymo.SubscriptionLabels, id string, routes centymo.SubscriptionRoutes, jobsTabVisible bool) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	items := []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		// 2026-04-27 plan-client-scope plan §6.5 — Package tab.
		{Key: "package", Label: l.Detail.Plan, Href: base + "?tab=package", HxGet: action + "package", Icon: "icon-package"},
		// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.2 — Operations tab.
		{Key: "operations", Label: l.Tabs.Operations, Href: base + "?tab=operations", HxGet: action + "operations", Icon: "icon-briefcase"},
	}
	// 2026-04-30 cyclic-subscription-jobs plan §21.3 — flat Jobs tab; hidden
	// when COUNT(jobs) == 0 (SaaS / advisory subscriptions).
	if jobsTabVisible {
		items = append(items, pyeza.TabItem{
			Key: "jobs", Label: l.Tabs.Jobs, Href: base + "?tab=jobs", HxGet: action + "jobs", Icon: "icon-list",
		})
	}
	items = append(items,
		pyeza.TabItem{Key: "invoices", Label: l.Tabs.Invoices, Href: base + "?tab=invoices", HxGet: action + "invoices", Icon: "icon-file-text"},
		pyeza.TabItem{Key: "attachments", Label: l.Tabs.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		pyeza.TabItem{Key: "audit", Label: l.Tabs.AuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
		pyeza.TabItem{Key: "audit-history", Label: l.Tabs.AuditHistory, Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	)
	return items
}

// buildPackageTabData populates the per-tab fields used by the Package tab
// per plan §6.5. The CTA is shown when pricePlan.client_id == "" (master);
// hidden when the PricePlan is already client-scoped. Caller passes the
// pre-resolved customize URL (centymo.SubscriptionCustomizePackageURL with
// the {id} substituted) — keeps this helper free of route knowledge.
func buildPackageTabData(sub *subscriptionpb.Subscription, customizeURL, customizeLabelTemplate string, canCustomize bool) (shown, disabled bool, label, clientName string, pp map[string]any) {
	if sub == nil {
		return false, false, "", "", nil
	}
	if c := sub.GetClient(); c != nil {
		clientName = c.GetName()
		if clientName == "" {
			if u := c.GetUser(); u != nil {
				clientName = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			}
		}
	}
	if pricePlan := sub.GetPricePlan(); pricePlan != nil {
		pp = map[string]any{
			"id":         pricePlan.GetId(),
			"name":       pricePlan.GetName(),
			"plan_id":    pricePlan.GetPlanId(),
			"client_id":  pricePlan.GetClientId(),
			"amount":     pricePlan.GetBillingAmount(),
			"currency":   pricePlan.GetBillingCurrency(),
		}
		// CTA gating per plan §6.5 / decision #6:
		//   - master (client_id == "") → show.
		//   - already client-scoped → hide (offer Edit instead).
		shown = pricePlan.GetClientId() == ""
	}
	if shown {
		label = strings.ReplaceAll(customizeLabelTemplate, "{{.ClientName}}", clientName)
		disabled = !canCustomize
	}
	_ = customizeURL // template reads this from PageData.PackageCustomizeURL
	return shown, disabled, label, clientName, pp
}

// loadMilestoneRows fetches BillingEvent rows for a subscription and converts
// them into the per-row template shape per flow.md §10. Returns the rendered
// rows, the running total invoiced (centavos), and the inferred currency.
//
// 2026-04-29 milestone-billing plan §5 / Phase D.
func loadMilestoneRows(ctx context.Context, deps *DetailViewDeps, subscriptionID string) ([]MilestoneRow, int64, string) {
	if deps.ListBillingEventsBySubscription == nil {
		return nil, 0, ""
	}
	resp, err := deps.ListBillingEventsBySubscription(ctx, &billingeventpb.ListBillingEventsBySubscriptionRequest{
		SubscriptionId: subscriptionID,
	})
	if err != nil || resp == nil {
		return nil, 0, ""
	}
	events := resp.GetBillingEvents()
	mLabels := deps.Labels.Milestone
	currency := ""
	var totalInvoiced int64
	rows := make([]MilestoneRow, 0, len(events))
	for _, ev := range events {
		if currency == "" {
			currency = ev.GetBillingCurrency()
		}
		status := ev.GetStatus()
		statusKey := statusKeyForBillingEventDetail(status)
		statusLabel := statusLabelForBillingEventDetail(status, mLabels)
		seq := strings.TrimSpace(ev.GetSequenceLabel())
		if seq == "" {
			id := ev.GetId()
			if len(id) > 8 {
				seq = "Event " + id[len(id)-6:]
			} else {
				seq = "Event " + id
			}
		}
		row := MilestoneRow{
			EventID:         ev.GetId(),
			StatusKey:       statusKey,
			StatusLabel:     statusLabel,
			SequenceLabel:   seq,
			BillableAmount:  ev.GetBillableAmount(),
			BillableDisplay: formatCentavoDisplay(ev.GetBillableAmount()),
			Currency:        ev.GetBillingCurrency(),
		}
		// Mark-ready when status = UNSPECIFIED (pending) or DEFERRED.
		row.ShowMarkReady = status == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_UNSPECIFIED ||
			status == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_DEFERRED
		// Waive when UNSPECIFIED or READY.
		row.ShowWaive = status == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_UNSPECIFIED ||
			status == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_READY
		if status == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_BILLED {
			row.ShowRevenueLink = true
			row.RevenueID = ev.GetRevenueId()
			if row.RevenueID != "" {
				row.RevenueURL = strings.ReplaceAll(centymo.RevenueDetailURL, "{id}", row.RevenueID)
			}
			totalInvoiced += ev.GetBillableAmount()
		}
		// Resolve URLs for the action buttons.
		row.MarkReadyURL = resolveBillingEventURL(deps.Routes.MilestoneMarkReadyURL, subscriptionID, ev.GetId())
		row.WaiveURL = resolveBillingEventURL(deps.Routes.MilestoneWaiveURL, subscriptionID, ev.GetId())
		rows = append(rows, row)
	}
	return rows, totalInvoiced, currency
}

func resolveBillingEventURL(template, subscriptionID, eventID string) string {
	if template == "" {
		return ""
	}
	r := strings.ReplaceAll(template, "{id}", subscriptionID)
	r = strings.ReplaceAll(r, "{eventId}", eventID)
	return r
}

func statusKeyForBillingEventDetail(s billingeventpb.BillingEventStatus) string {
	switch s {
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_READY:
		return "ready"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_BILLED:
		return "billed"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_WAIVED:
		return "waived"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_DEFERRED:
		return "deferred"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_CANCELLED:
		return "cancelled"
	default:
		return "pending"
	}
}

func statusLabelForBillingEventDetail(s billingeventpb.BillingEventStatus, l centymo.SubscriptionMilestoneLabels) string {
	switch s {
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_READY:
		return l.StatusReady
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_BILLED:
		return l.StatusBilled
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_WAIVED:
		return l.StatusWaived
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_DEFERRED:
		return l.StatusDeferred
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_CANCELLED:
		return l.StatusCancelled
	default:
		return l.StatusPending
	}
}

func formatCentavoDisplay(c int64) string {
	whole := c / 100
	frac := c % 100
	if frac < 0 {
		frac = -frac
	}
	return fmt.Sprintf("%d.%02d", whole, frac)
}

// applyMilestoneTabData populates the package-tab milestone fields when the
// subscription's PricePlan is MILESTONE. No-op for other billing kinds.
func applyMilestoneTabData(ctx context.Context, deps *DetailViewDeps, sub *subscriptionpb.Subscription, pageData *PageData, subscriptionID string) {
	if sub == nil {
		return
	}
	pp := sub.GetPricePlan()
	if pp == nil || pp.GetBillingKind() != priceplanpb.BillingKind_BILLING_KIND_MILESTONE {
		return
	}
	rows, total, currency := loadMilestoneRows(ctx, deps, subscriptionID)
	pageData.MilestonesShown = true
	pageData.Milestones = rows
	pageData.TotalInvoicedDisplay = formatCentavoDisplay(total)
	pageData.MilestoneCurrency = currency
}

// loadSubscriptionInvoices fetches revenue records filtered by subscription_id.
// Returns an empty slice on error or when the dep is nil.
func loadSubscriptionInvoices(ctx context.Context, deps *DetailViewDeps, subscriptionID string) []*revenuepb.Revenue {
	if deps.GetRevenueListPageData == nil {
		return nil
	}
	resp, err := deps.GetRevenueListPageData(ctx, &revenuepb.GetRevenueListPageDataRequest{
		Filters: &commonpb.FilterRequest{
			Filters: []*commonpb.TypedFilter{
				{
					Field: "rv.subscription_id",
					FilterType: &commonpb.TypedFilter_StringFilter{
						StringFilter: &commonpb.StringFilter{
							Value:         subscriptionID,
							Operator:      commonpb.StringOperator_STRING_EQUALS,
							CaseSensitive: true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to load invoices for subscription %s: %v", subscriptionID, err)
		return nil
	}
	return resp.GetRevenueList()
}

// buildInvoicesTable builds a TableConfig for the invoices tab.
// Columns: reference number (code), date, amount, status.
//
// Surfaces a "Recognize Revenue" PrimaryAction on the toolbar AND on the
// empty-state when the operator has revenue:create AND the subscription is
// active. Disabled state degrades to a tooltip explaining the gate.
//
// Note: NO page-header "Recognize Revenue" button — per plan §11.2 / §4.2,
// cause and effect stay adjacent on the invoices tab.
func buildInvoicesTable(
	revenues []*revenuepb.Revenue,
	l centymo.SubscriptionLabels,
	tableLabels types.TableLabels,
	detailURLTemplate string,
	recognizeURL string,
	subscriptionID string,
	canRecognize bool,
	subscriptionActive bool,
	disabledTooltip string,
) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "reference_number", Label: l.Invoices.ColumnCode, Sortable: true},
		{Key: "revenue_date_string", Label: l.Invoices.ColumnDate, Sortable: true, WidthClass: "col-3xl"},
		{Key: "total_amount", Label: l.Invoices.ColumnAmount, Sortable: true, WidthClass: "col-3xl", Align: "right"},
		{Key: "status", Label: l.Invoices.ColumnStatus, Sortable: true, WidthClass: "col-2xl"},
	}

	var rows []types.TableRow
	for _, r := range revenues {
		id := r.GetId()
		refNumber := r.GetReferenceNumber()
		currency := r.GetCurrency()
		status := r.GetStatus()

		statusVariant := "default"
		switch status {
		case "draft":
			statusVariant = "warning"
		case "complete":
			statusVariant = "success"
		case "cancelled":
			statusVariant = "danger"
		}

		var detailHref string
		if detailURLTemplate != "" {
			detailHref = route.ResolveURL(detailURLTemplate, "id", id)
		}

		rows = append(rows, types.TableRow{
			ID:   id,
			Href: detailHref,
			Cells: []types.TableCell{
				{Type: "text", Value: refNumber},
				types.DateTimeCell(r.GetRevenueDate(), types.DateReadable),
				types.MoneyCell(float64(r.GetTotalAmount()), currency, true),
				{Type: "badge", Value: status, Variant: statusVariant},
			},
		})
	}

	types.ApplyColumnStyles(columns, rows)

	tc := &types.TableConfig{
		ID:                   "subscription-invoices-table",
		Columns:              columns,
		Rows:                 rows,
		Labels:               tableLabels,
		ShowSearch:           false,
		ShowActions:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "revenue_date_string",
		DefaultSortDirection: "desc",
		EmptyState: types.TableEmptyState{
			Title:   l.Invoices.Title,
			Message: l.Invoices.Empty,
		},
	}

	// PrimaryAction: tab toolbar CTA + same action surfaced on the empty state
	// (the table-card template renders the primary action in both places).
	if recognizeURL != "" && subscriptionID != "" {
		actionURL := route.ResolveURL(recognizeURL, "id", subscriptionID)
		disabled := !canRecognize || !subscriptionActive
		tc.PrimaryAction = &types.PrimaryAction{
			Label:           l.Invoices.RecognizeAction,
			ActionURL:       actionURL,
			Icon:            "icon-file-plus",
			Disabled:        disabled,
			DisabledTooltip: disabledTooltip,
		}
	}

	types.ApplyTableSettings(tc)
	return tc
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		sub, err := loadSubscriptionWithRelations(ctx, deps, id)
		if err != nil {
			log.Printf("Failed to read subscription %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load subscription: %w", err))
		}
		subscription := subscriptionToMap(ctx, sub)

		l := deps.Labels

		// 2026-04-30 cyclic-subscription-jobs plan — pre-load the Job set so
		// the tab handler shares the same fan-out as the full-page handler.
		isCyclic := computeIsCyclic(sub)
		allJobs := loadSubscriptionJobs(ctx, deps, id)
		jobsTabVisible := len(allJobs) > 0

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Subscription: subscription,
			Labels:       l,
			ActiveTab:    tab,
			TabItems:     buildTabItems(l, id, deps.Routes, jobsTabVisible),
			IsCyclic:     isCyclic,
			SpawnCycleJobsURL: strings.ReplaceAll(
				deps.Routes.SpawnCycleJobsURL, "{subscriptionId}", id),
			BackfillCycleJobsDrawerURL: strings.ReplaceAll(
				deps.Routes.BackfillCycleJobsURL, "{subscriptionId}", id),
		}

		// Same tab_invoices_url + perms gating as the full-page handler so a
		// tab-only refresh sees a consistent PrimaryAction state.
		subscription["tab_invoices_url"] = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "invoices"
		perms := view.GetUserPermissions(ctx)
		// nil perms = no restrictions (dev / mock mode). Consistent with
		// the comment on UserPermissions.Can() and most other view checks.
		canRecognize := perms == nil || perms.Can("revenue", "create")
		subscriptionActive, _ := subscription["active"].(bool)

		switch tab {
		case "package":
			canCustomize := perms != nil && (perms.Can("revenue", "create") || perms.Can("plan", "create"))
			customizeURL := route.ResolveURL(deps.Routes.CustomizePackageURL, "id", id)
			shown, disabled, label, clientName, ppData := buildPackageTabData(sub, customizeURL, l.Actions.CustomizePackage, canCustomize)
			pageData.PackageCustomizeURL = customizeURL
			pageData.PackageCustomizeShown = shown && subscriptionActive
			pageData.PackageCustomizeDisabled = disabled
			pageData.PackageCustomizeLabel = label
			pageData.PackageClientName = clientName
			pageData.PackagePricePlan = ppData
			// 2026-04-29 milestone-billing — milestones list (only on MILESTONE plans).
			applyMilestoneTabData(ctx, deps, sub, pageData, id)
		case "operations":
			applyOperationsTabData(ctx, deps, pageData, id, sub, allJobs, isCyclic)
		case "jobs":
			applyJobsTabData(ctx, deps, pageData, allJobs)
		case "invoices":
			revenues := loadSubscriptionInvoices(ctx, deps, id)
			pageData.Invoices = buildInvoicesTable(
				revenues, l, deps.TableLabels, centymo.RevenueDetailURL,
				deps.Routes.RecognizeURL, id,
				canRecognize, subscriptionActive,
				resolveRecognizeDisabledTooltip(canRecognize, subscriptionActive, l),
			)
		case "attachments":
			if deps.ListAttachments != nil {
				cfg := attachmentConfig(deps)
				resp, err := deps.ListAttachments(ctx, cfg.EntityType, id)
				if err != nil {
					log.Printf("Failed to list attachments: %v", err)
				}
				var items []*attachmentpb.Attachment
				if resp != nil {
					items = resp.GetData()
				}
				pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
			}
			pageData.AttachmentUploadURL = route.ResolveURL(deps.Routes.AttachmentUploadURL, "id", id)
		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.QueryParams["cursor"]
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "subscription",
					EntityID:    id,
					Limit:       20,
					CursorToken: cursor,
				})
				if err != nil {
					log.Printf("Failed to load audit history: %v", err)
				}
				if auditResp != nil {
					pageData.AuditEntries = auditResp.Entries
					pageData.AuditHasNext = auditResp.HasNext
					pageData.AuditNextCursor = auditResp.NextCursor
				}
			}
			pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
		}

		templateName := "subscription-tab-" + tab
		if tab == "invoices" {
			templateName = "subscription-tab-invoices"
		}
		if tab == "package" {
			templateName = "subscription-tab-package"
		}
		if tab == "operations" {
			templateName = "subscription-tab-operations"
		}
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// applyOperationsTabData populates the Operations tab fields on PageData.
// Reads spawned Jobs via GetJobsByOrigin (origin_type = SUBSCRIPTION,
// origin_id = subscription.id), then enriches each Job with its phase rollup
// from ListJobPhasesByJob. Jobs with parent_job_id are nested under their
// parent. nil-safe: when deps are unwired the tab degrades to the empty
// state.
//
// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.2 +
// 2026-04-30 cyclic-subscription-jobs plan §7.1 — branches on isCyclic.
//
// jobs slice is pre-loaded by the page handler (single fan-out so the Jobs
// tab + Operations tab + visibility gate share one read). When non-cyclic
// the legacy parent-child rendering is preserved verbatim (regression
// guard: `09-non-cyclic-unaffected.spec.ts` is the canary).
func applyOperationsTabData(
	ctx context.Context,
	deps *DetailViewDeps,
	pageData *PageData,
	subscriptionID string,
	sub *subscriptionpb.Subscription,
	jobs []*jobpb.Job,
	isCyclic bool,
) {
	pageData.OperationsSpawnURL = strings.ReplaceAll(deps.SpawnJobsURL, "{subscriptionId}", subscriptionID)
	pageData.OperationsEmptyText = strings.ReplaceAll(
		deps.Labels.Operations.EmptyMessage,
		"{{.SpawnAction}}",
		deps.Labels.Operations.SpawnAction,
	)
	if isCyclic {
		// Cyclic branch — render the cycle accordion. Legacy fields stay
		// zero-valued so the template's IsCyclic-gated section takes over.
		pageData.OperationsCyclic = buildSubscriptionCyclesData(ctx, deps, sub, jobs)
		// HasJobs stays true when ANY job exists (engagement shell counts).
		// The template falls back to the empty-state CTA only when the
		// cyclic data block has no engagement and no cycles AND no Jobs.
		pageData.OperationsHasJobs = len(jobs) > 0
		return
	}
	// Non-cyclic — preserve existing rendering path verbatim.
	if len(jobs) == 0 {
		return
	}
	rows := buildOperationsRows(ctx, deps, jobs)
	pageData.OperationsHasJobs = len(rows) > 0
	pageData.OperationsRootJobs = rows
}

// loadSubscriptionJobs reads ALL Jobs for a subscription (origin_type =
// SUBSCRIPTION, origin_id = subscription.id). Cached at the page-handler
// level so Operations tab + Jobs tab + visibility gate share one round-trip.
//
// 2026-04-30 cyclic-subscription-jobs plan §21.5.
func loadSubscriptionJobs(ctx context.Context, deps *DetailViewDeps, subscriptionID string) []*jobpb.Job {
	if deps.GetJobsByOrigin == nil || subscriptionID == "" {
		return nil
	}
	resp, err := deps.GetJobsByOrigin(ctx, &jobpb.GetJobsByOriginRequest{
		OriginType: enums.OriginType_ORIGIN_TYPE_SUBSCRIPTION,
		OriginId:   subscriptionID,
	})
	if err != nil {
		log.Printf("Failed to load jobs for subscription %s: %v", subscriptionID, err)
		return nil
	}
	if resp == nil {
		return nil
	}
	return resp.GetJobs()
}

// computeIsCyclic mirrors espyna's `eligibleForInstanceSpawn` predicate:
// RECURRING OR (CONTRACT AND billing_cycle_value > 0). See
// cyclic-subscription-jobs plan §3.1 / §19.3 (eligibility-gate refactor for
// AD_HOC forward-compat).
func computeIsCyclic(sub *subscriptionpb.Subscription) bool {
	if sub == nil {
		return false
	}
	pp := sub.GetPricePlan()
	if pp == nil {
		return false
	}
	kind := pp.GetBillingKind()
	if kind == priceplanpb.BillingKind_BILLING_KIND_RECURRING {
		return true
	}
	if kind == priceplanpb.BillingKind_BILLING_KIND_CONTRACT && pp.GetBillingCycleValue() > 0 {
		return true
	}
	return false
}

// buildSubscriptionCyclesData partitions the Job list into engagement shell
// + once-at-start children + cycle accordions per cyclic-subscription-jobs
// plan §7.1. Sorted descending by CycleIndex (most-recent on top).
func buildSubscriptionCyclesData(
	ctx context.Context,
	deps *DetailViewDeps,
	sub *subscriptionpb.Subscription,
	jobs []*jobpb.Job,
) *SubscriptionCyclesData {
	data := &SubscriptionCyclesData{
		CycleEmpty: deps.Labels.Operations.CycleEmpty,
	}

	// Pass 1 — find engagement shell (parent_job_id == "").
	var engagementJob *jobpb.Job
	for _, j := range jobs {
		if j.GetParentJobId() == "" {
			engagementJob = j
			break
		}
	}
	if engagementJob != nil {
		data.EngagementJobID = engagementJob.GetId()
		data.EngagementName = engagementJob.GetName()
		started := ""
		if ts := sub.GetDateTimeStart(); ts != nil && ts.IsValid() {
			started = ts.AsTime().Format("2006-01-02")
		}
		r := strings.NewReplacer(
			"{{.Started}}", started,
			"{{.Name}}", engagementJob.GetName(),
		)
		data.EngagementHeading = r.Replace(deps.Labels.Operations.EngagementHeading)
	}

	// Pass 2 — bucket children: cycle Jobs (have cycle_index) vs onboarding
	// (parent_job_id set but cycle_index == 0 / nil).
	cyclesByIndex := map[int32][]*jobpb.Job{}
	cycleIndices := []int32{}
	for _, j := range jobs {
		parentID := j.GetParentJobId()
		if parentID == "" {
			continue
		}
		if j.GetCycleIndex() == 0 {
			// Onboarding / once-at-start — cycle_index is NULL on the proto
			// so the getter returns 0.
			data.OnceAtStartJobs = append(data.OnceAtStartJobs, jobToOperationsRow(ctx, deps, j))
			continue
		}
		idx := j.GetCycleIndex()
		if _, exists := cyclesByIndex[idx]; !exists {
			cycleIndices = append(cycleIndices, idx)
		}
		cyclesByIndex[idx] = append(cyclesByIndex[idx], j)
	}

	// Sort indices descending (most recent on top).
	sort.Slice(cycleIndices, func(i, j int) bool {
		return cycleIndices[i] > cycleIndices[j]
	})

	for i, idx := range cycleIndices {
		group := cyclesByIndex[idx]
		view := buildSubscriptionCycleView(ctx, deps, idx, group)
		// Most-recent cycle expanded by default; older cycles collapsed.
		view.OpenByDefault = (i == 0)
		data.Cycles = append(data.Cycles, view)
	}
	return data
}

// buildSubscriptionCycleView constructs one cycle accordion entry from the
// Jobs in that cycle (1 for visits_per_cycle=1, N for multi-visit). The first
// Job's cycle_period_* fields represent the cycle's overall window.
func buildSubscriptionCycleView(
	ctx context.Context,
	deps *DetailViewDeps,
	cycleIndex int32,
	group []*jobpb.Job,
) SubscriptionCycleView {
	v := SubscriptionCycleView{
		CycleIndex: cycleIndex,
	}
	if len(group) == 0 {
		return v
	}
	first := group[0]
	v.PeriodStart = first.GetCyclePeriodStart()
	v.PeriodEnd = first.GetCyclePeriodEnd()
	v.PeriodLabel = formatCyclePeriodLabel(v.PeriodStart, v.PeriodEnd)

	// Heading template is "Cycle {{.CycleIndex}} — {{.PeriodLabel}}".
	r := strings.NewReplacer(
		"{{.CycleIndex}}", intStr(int(cycleIndex)),
		"{{.PeriodLabel}}", v.PeriodLabel,
	)
	v.HeadingText = r.Replace(deps.Labels.Operations.CycleHeading)

	// Status rollup — view-side aggregation per plan §7.1.
	v.StatusKey, v.StatusLabel, v.StatusVariant = rollupCycleStatus(group, deps.Labels.Operations)

	// Per-Job rows (operations rendering — phase summary etc).
	v.Jobs = make([]OperationsJobRow, 0, len(group))
	for _, j := range group {
		v.Jobs = append(v.Jobs, jobToOperationsRow(ctx, deps, j))
	}

	// InvoiceLabel — cycleNoInvoice when no Revenue is matched yet. The
	// matched-revenue path requires a separate Revenue lookup by date join
	// (plan §2.2 — no proto FK in v1); v1 surfaces "Not yet invoiced" until
	// the Revenue lookup is wired through.
	v.InvoiceLabel = deps.Labels.Operations.CycleNoInvoice
	return v
}

// formatCyclePeriodLabel renders a "May 2026"-style label from ISO date
// strings. Falls back to the raw start string if parsing fails.
func formatCyclePeriodLabel(start, _ string) string {
	if start == "" {
		return ""
	}
	// Use a stable format key recognised by Go's time package. Keep parsing
	// loose — operators may store DateTime values too.
	for _, layout := range []string{"2006-01-02", "2006-01-02T15:04:05Z07:00", time.RFC3339} {
		if t, err := time.Parse(layout, start); err == nil {
			return t.Format("Jan 2006")
		}
	}
	return start
}

// rollupCycleStatus aggregates per-Job statuses into a cycle-level rollup.
// Pending if all Jobs are PLANNED/PENDING/DRAFT; In progress if any ACTIVE;
// Completed when all CLOSED/COMPLETED; Overdue when at least one is past-due
// (proxy: PAUSED). Lyngua-resolved labels.
func rollupCycleStatus(group []*jobpb.Job, l centymo.SubscriptionOperationsLabels) (string, string, string) {
	hasActive := false
	hasOverdue := false
	allDone := true
	for _, j := range group {
		switch j.GetStatus() {
		case enums.JobStatus_JOB_STATUS_ACTIVE:
			hasActive = true
			allDone = false
		case enums.JobStatus_JOB_STATUS_PAUSED:
			hasOverdue = true
			allDone = false
		case enums.JobStatus_JOB_STATUS_COMPLETED, enums.JobStatus_JOB_STATUS_CLOSED:
			// done
		default:
			allDone = false
		}
	}
	switch {
	case hasOverdue:
		return "overdue", l.CycleStatusOverdue, "danger"
	case hasActive:
		return "inProgress", l.CycleStatusInProgress, "success"
	case allDone && len(group) > 0:
		return "completed", l.CycleStatusCompleted, "info"
	default:
		return "pending", l.CycleStatusPending, "warning"
	}
}

// applyJobsTabData populates the new flat Jobs tab. Hidden (HasJobs=false)
// when no Jobs exist; the page handler also gates the tab nav item upstream.
//
// 2026-04-30 cyclic-subscription-jobs plan §21.5.
func applyJobsTabData(
	ctx context.Context,
	deps *DetailViewDeps,
	pageData *PageData,
	jobs []*jobpb.Job,
) {
	if pageData.Jobs == nil {
		pageData.Jobs = &SubscriptionJobsTabData{
			HasJobs:      len(jobs) > 0,
			StatusCounts: map[string]int{},
			TypeCounts:   map[string]int{},
		}
	}
	if len(jobs) == 0 {
		return
	}
	l := deps.Labels.Jobs
	rows := make([]SubscriptionJobRow, 0, len(jobs))
	for _, j := range jobs {
		jobType := jobTypeKey(j)
		typeLabel := jobTypeLabel(jobType, l)
		statusKey, statusVariant := operationsJobStatusInfo(j.GetStatus())
		statusLabel := statusLabelForJobStatus(j.GetStatus(), deps.Labels.Operations)
		periodLabel := ""
		if j.GetCycleIndex() != 0 {
			periodLabel = formatCyclePeriodLabel(j.GetCyclePeriodStart(), j.GetCyclePeriodEnd())
		}
		row := SubscriptionJobRow{
			JobID:         j.GetId(),
			JobName:       j.GetName(),
			JobType:       jobType,
			JobTypeLabel:  typeLabel,
			Status:        statusKey,
			StatusLabel:   statusLabel,
			StatusVariant: statusVariant,
			PeriodLabel:   periodLabel,
			CycleIndex:    j.GetCycleIndex(),
		}
		if deps.JobDetailURL != "" {
			row.DetailURL = strings.ReplaceAll(deps.JobDetailURL, "{id}", j.GetId())
		}
		rows = append(rows, row)
		pageData.Jobs.StatusCounts[statusKey]++
		pageData.Jobs.TypeCounts[jobType]++
	}
	// Sort by cycle_index descending (engagement first, then most-recent
	// cycles); ties broken by name. Engagement shell has cycle_index=0 and
	// parent_job_id="" — pin it to top by giving it a sentinel.
	sort.SliceStable(rows, func(i, jj int) bool {
		if rows[i].JobType == "engagement" {
			return true
		}
		if rows[jj].JobType == "engagement" {
			return false
		}
		if rows[i].CycleIndex != rows[jj].CycleIndex {
			return rows[i].CycleIndex > rows[jj].CycleIndex
		}
		return rows[i].JobName < rows[jj].JobName
	})
	pageData.Jobs.Rows = rows
}

// jobTypeKey classifies a Job for the Jobs tab Type column / filter chips.
// engagement = parent_job_id empty (shell);
// onboarding = parent_job_id set AND cycle_index == 0 (ONCE_AT_ENGAGEMENT_START);
// cycle      = parent_job_id set AND cycle_index > 0.
// AD_HOC plan reinterprets "cycle" as "visit" — that override lives in the
// downstream plan, not here.
func jobTypeKey(j *jobpb.Job) string {
	if j.GetParentJobId() == "" {
		return "engagement"
	}
	if j.GetCycleIndex() == 0 {
		return "onboarding"
	}
	return "cycle"
}

func jobTypeLabel(key string, l centymo.SubscriptionJobsTabLabels) string {
	switch key {
	case "engagement":
		return l.TypeEngagement
	case "onboarding":
		return l.TypeOnboarding
	case "cycle":
		return l.TypeCycle
	case "visit":
		return l.TypeVisit
	}
	return key
}

func statusLabelForJobStatus(s enums.JobStatus, l centymo.SubscriptionOperationsLabels) string {
	switch s {
	case enums.JobStatus_JOB_STATUS_ACTIVE:
		return l.CycleStatusInProgress
	case enums.JobStatus_JOB_STATUS_COMPLETED, enums.JobStatus_JOB_STATUS_CLOSED:
		return l.CycleStatusCompleted
	case enums.JobStatus_JOB_STATUS_PAUSED:
		return l.CycleStatusOverdue
	default:
		return l.CycleStatusPending
	}
}

// buildOperationsRows converts a flat Job slice into a parent-child tree
// suitable for the Operations tab template. Each row carries its phase
// summary string already rendered.
func buildOperationsRows(ctx context.Context, deps *DetailViewDeps, jobs []*jobpb.Job) []OperationsJobRow {
	byID := map[string]*OperationsJobRow{}
	roots := make([]OperationsJobRow, 0)
	// First pass: build node map.
	for _, j := range jobs {
		row := jobToOperationsRow(ctx, deps, j)
		byID[j.GetId()] = &row
	}
	// Second pass: link children, collect roots.
	for _, j := range jobs {
		row := byID[j.GetId()]
		if row == nil {
			continue
		}
		parentID := j.GetParentJobId()
		if parentID == "" {
			roots = append(roots, *row)
			continue
		}
		parent, ok := byID[parentID]
		if !ok {
			// Orphan child — render at the root for visibility.
			roots = append(roots, *row)
			continue
		}
		parent.Children = append(parent.Children, *row)
	}
	// Re-link children from the byID map after the appends (we need the
	// updated parent in the roots slice — Go maps store pointers but value
	// receivers in roots copied. Recompute the roots' Children from byID).
	for i := range roots {
		if updated, ok := byID[roots[i].JobID]; ok {
			roots[i].Children = updated.Children
		}
	}
	return roots
}

func jobToOperationsRow(ctx context.Context, deps *DetailViewDeps, j *jobpb.Job) OperationsJobRow {
	statusKey, statusVariant := operationsJobStatusInfo(j.GetStatus())
	row := OperationsJobRow{
		JobID:          j.GetId(),
		JobName:        j.GetName(),
		IsRoot:         j.GetParentJobId() == "",
		StatusKey:      statusKey,
		StatusVariant:  statusVariant,
		BillingRuleKey: operationsBillingRuleKey(j.GetBillingRuleType()),
	}
	if deps.JobDetailURL != "" {
		row.JobDetailURL = strings.ReplaceAll(deps.JobDetailURL, "{id}", j.GetId())
	}
	row.PhaseSummaryText = renderPhaseSummary(ctx, deps, j.GetId())
	return row
}

func renderPhaseSummary(ctx context.Context, deps *DetailViewDeps, jobID string) string {
	if deps.ListJobPhasesByJob == nil {
		return ""
	}
	resp, err := deps.ListJobPhasesByJob(ctx, &jobphasepb.ListJobPhasesByJobRequest{JobId: jobID})
	if err != nil || resp == nil {
		return ""
	}
	phases := resp.GetJobPhases()
	total := len(phases)
	complete := 0
	for _, p := range phases {
		if p.GetStatus() == jobphasepb.PhaseStatus_PHASE_STATUS_COMPLETED {
			complete++
		}
	}
	tmpl := deps.Labels.Operations.PhaseSummary
	r := strings.NewReplacer(
		"{{.Complete}}", intStr(complete),
		"{{.Total}}", intStr(total),
	)
	return r.Replace(tmpl)
}

func operationsJobStatusInfo(s enums.JobStatus) (string, string) {
	switch s {
	case enums.JobStatus_JOB_STATUS_DRAFT:
		return "draft", "default"
	case enums.JobStatus_JOB_STATUS_PENDING:
		return "pending", "warning"
	case enums.JobStatus_JOB_STATUS_PLANNED:
		return "planned", "default"
	case enums.JobStatus_JOB_STATUS_ACTIVE:
		return "active", "success"
	case enums.JobStatus_JOB_STATUS_PAUSED:
		return "paused", "warning"
	case enums.JobStatus_JOB_STATUS_COMPLETED:
		return "completed", "info"
	case enums.JobStatus_JOB_STATUS_CLOSED:
		return "closed", "default"
	default:
		return "draft", "default"
	}
}

func operationsBillingRuleKey(b enums.BillingRuleType) string {
	switch b {
	case enums.BillingRuleType_BILLING_RULE_TYPE_MILESTONE:
		return "milestone"
	case enums.BillingRuleType_BILLING_RULE_TYPE_T_AND_M:
		return "t_and_m"
	case enums.BillingRuleType_BILLING_RULE_TYPE_NON_BILLABLE:
		return "non_billable"
	case enums.BillingRuleType_BILLING_RULE_TYPE_FIXED_FEE:
		return "fixed_fee"
	case enums.BillingRuleType_BILLING_RULE_TYPE_INCLUDED:
		return "included"
	default:
		return "unspecified"
	}
}

func intStr(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
