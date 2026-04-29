package detail

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/erniealice/centymo-golang"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	revenuepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/revenue"
	subscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/subscription"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes           centymo.SubscriptionRoutes
	ReadSubscription func(ctx context.Context, req *subscriptionpb.ReadSubscriptionRequest) (*subscriptionpb.ReadSubscriptionResponse, error)
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
		tabItems := buildTabItems(l, id, deps.Routes)

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

func buildTabItems(l centymo.SubscriptionLabels, id string, routes centymo.SubscriptionRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		// 2026-04-27 plan-client-scope plan §6.5 — Package tab.
		{Key: "package", Label: l.Detail.Plan, Href: base + "?tab=package", HxGet: action + "package", Icon: "icon-package"},
		{Key: "invoices", Label: l.Tabs.Invoices, Href: base + "?tab=invoices", HxGet: action + "invoices", Icon: "icon-file-text"},
		{Key: "attachments", Label: l.Tabs.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		{Key: "audit", Label: l.Tabs.AuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
		{Key: "audit-history", Label: l.Tabs.AuditHistory, Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
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
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Subscription: subscription,
			Labels:       l,
			ActiveTab:    tab,
			TabItems:     buildTabItems(l, id, deps.Routes),
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
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}
