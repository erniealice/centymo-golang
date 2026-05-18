package detail

import (
	"context"
	"fmt"
	"log"
	"strings"

	centymo "github.com/erniealice/centymo-golang"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	billingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/billing_event"
	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
	// 20260517-advance-cash-events Plan B Phase 7 — junction reads.
	junctionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/treasury_collection_billing_event"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes         centymo.CollectionRoutes
	ReadCollection func(ctx context.Context, req *collectionpb.ReadCollectionRequest) (*collectionpb.ReadCollectionResponse, error)
	Labels         centymo.CollectionLabels
	CommonLabels   pyeza.CommonLabels
	TableLabels    types.TableLabels

	// 20260517-advance-cash-events Plan B Phase 4 — Advance Schedule tab
	// typed-label sources. AdvanceLabels carries the tab + column + section
	// strings; AdvanceEnumLabels carries the AdvanceKind / AdvanceStatus /
	// AdvanceProrationPolicy option labels rendered as badge text. Both
	// thread up to PageData and are referenced by the
	// collection-tab-advance-schedule template.
	AdvanceLabels     centymo.TreasuryAdvanceLabels
	AdvanceEnumLabels centymo.AdvanceEnumLabels

	// 20260517-advance-cash-events Plan B Phase 7 — MILESTONE cross-link.
	// When advance_kind == MILESTONE, the Advance Schedule tab lists the
	// linked BillingEvent rows + per-event Recognize button. Both deps are
	// nil-safe — when unset, the milestone-link section is omitted.
	ListTreasuryCollectionBillingEvents func(ctx context.Context, req *junctionpb.ListTreasuryCollectionBillingEventsRequest) (*junctionpb.ListTreasuryCollectionBillingEventsResponse, error)
	ReadBillingEvent                    func(ctx context.Context, req *billingeventpb.ReadBillingEventRequest) (*billingeventpb.ReadBillingEventResponse, error)
	// MilestoneRecognizeURL is the route template
	// (`/action/subscription/{id}/billing-event/{eventId}/recognize`)
	// used by the per-row Recognize button. Empty means the button is hidden.
	MilestoneRecognizeURL string

	attachment.AttachmentOps
	auditlog.AuditOps
}

// MilestoneLinkRow is a single linked-BillingEvent row shown in the Advance
// Schedule tab when advance_kind == MILESTONE.
type MilestoneLinkRow struct {
	BillingEventID  string
	SubscriptionID  string
	TrancheAmount   int64
	TrancheDisplay  string
	Currency        string
	StatusKey       string
	StatusLabel     string
	RevenueID       string
	ShowRecognize   bool
	RecognizeURL    string
}

// PageData holds the data for the collection detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Collection      map[string]any
	Labels          centymo.CollectionLabels
	// 20260517-advance-cash-events Plan B Phase 4 — typed-label sources for
	// the Advance Schedule tab. Always populated; the template guards on
	// is_advance / MilestoneLinks for visibility.
	AdvanceLabels   centymo.TreasuryAdvanceLabels
	AdvanceEnum     centymo.AdvanceEnumLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	AuditTable      *types.TableConfig
	AttachmentTable *types.TableConfig
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
	// 20260517-advance-cash-events Plan B Phase 7 — MILESTONE links rendered
	// on the Advance Schedule tab when advance_kind == MILESTONE.
	MilestoneLinks []MilestoneLinkRow
}

// collectionToMap converts a Collection protobuf to a map[string]any for template use.
//
// 20260517-advance-cash-events Plan B Phase 4 — the map now also carries the
// advance_* fields so the detail template can render the Advance Schedule tab
// + decide whether to surface the UNSCHEDULED Settle/Refund/Cancel CTAs. We
// project the kind/status/proration_policy enums to their stable string
// representations so the template branch tests stay readable.
//
// Phase 4 (label cleanup) — also project _label fields populated from the
// resolved AdvanceEnumLabels so the badge / metadata-grid templates render
// human-readable text without hardcoding English.
func collectionToMap(c *collectionpb.Collection, enumLabels centymo.AdvanceEnumLabels) map[string]any {
	advanceKind := c.GetAdvanceKind().String()
	advanceStatus := c.GetAdvanceStatus().String()
	advanceProrationPolicy := c.GetAdvanceProrationPolicy().String()
	isAdvance := advanceKind != "" &&
		advanceKind != "ADVANCE_KIND_UNSPECIFIED" &&
		advanceKind != "ADVANCE_KIND_NONE" &&
		advanceKind != "NONE"
	isUnscheduledActive := advanceKind == "ADVANCE_KIND_UNSCHEDULED" &&
		(advanceStatus == "ADVANCE_STATUS_ACTIVE" || advanceStatus == "ADVANCE_STATUS_PARTIALLY_SETTLED" ||
			advanceStatus == "ACTIVE" || advanceStatus == "PARTIALLY_SETTLED")

	return map[string]any{
		"id":                             c.GetId(),
		"name":                           c.GetName(),
		"reference_number":               c.GetReferenceNumber(),
		"amount":                         types.MoneyCell(float64(c.GetAmount()), c.GetCurrency(), true),
		"currency":                       c.GetCurrency(),
		"status":                         c.GetStatus(),
		"collection_method_id":           c.GetCollectionMethodId(),
		"collection_type":                c.GetCollectionType(),
		"revenue_id":                     c.GetRevenueId(),
		"received_by":                    c.GetReceivedBy(),
		"received_role":                  c.GetReceivedRole(),
		"active":                         c.GetActive(),
		"date_created_string":            c.GetDateCreatedString(),
		"date_modified_string":           c.GetDateModifiedString(),
		"advance_kind":                   advanceKind,
		"advance_status":                 advanceStatus,
		"advance_proration_policy":       advanceProrationPolicy,
		"advance_kind_label":             advanceKindLabel(advanceKind, enumLabels.Kind),
		"advance_status_label":           advanceStatusLabel(advanceStatus, enumLabels.Status),
		"advance_proration_policy_label": advanceProrationPolicyLabel(advanceProrationPolicy, enumLabels.ProrationPolicy),
		"advance_start_date":             c.GetAdvanceStartDate(),
		"advance_end_date":               c.GetAdvanceEndDate(),
		"advance_period_count":           c.GetAdvancePeriodCount(),
		"advance_period_unit":            c.GetAdvancePeriodUnit(),
		"advance_total_amount":           c.GetAdvanceTotalAmount(),
		"advance_remaining":              c.GetAdvanceRemainingAmount(),
		"advance_recognized":             c.GetAdvanceRecognizedAmount(),
		"advance_balance_account":        c.GetAdvanceBalanceAccountId(),
		"advance_target_account":         c.GetAdvanceTargetAccountId(),
		"is_advance":                     isAdvance,
		"is_unscheduled_active":          isUnscheduledActive,
	}
}

// advanceKindLabel resolves the AdvanceKind proto-enum string to its
// operator-facing label, accepting both the fully-qualified
// (ADVANCE_KIND_TIME_BASED) and bare (TIME_BASED) variants.
func advanceKindLabel(s string, l centymo.AdvanceKindLabels) string {
	switch s {
	case "ADVANCE_KIND_NONE", "NONE":
		return l.None
	case "ADVANCE_KIND_TIME_BASED", "TIME_BASED":
		return l.TimeBased
	case "ADVANCE_KIND_BURN_DOWN", "BURN_DOWN":
		return l.BurnDown
	case "ADVANCE_KIND_MILESTONE", "MILESTONE":
		return l.Milestone
	case "ADVANCE_KIND_UNSCHEDULED", "UNSCHEDULED":
		return l.Unscheduled
	}
	return ""
}

// advanceStatusLabel resolves the AdvanceStatus proto-enum string.
func advanceStatusLabel(s string, l centymo.AdvanceStatusLabels) string {
	switch s {
	case "ADVANCE_STATUS_ACTIVE", "ACTIVE":
		return l.Active
	case "ADVANCE_STATUS_FULLY_RECOGNIZED", "FULLY_RECOGNIZED":
		return l.FullyRecognized
	case "ADVANCE_STATUS_FULLY_AMORTIZED", "FULLY_AMORTIZED":
		return l.FullyAmortized
	case "ADVANCE_STATUS_FULLY_DRAWN", "FULLY_DRAWN":
		return l.FullyDrawn
	case "ADVANCE_STATUS_SETTLED", "SETTLED":
		return l.Settled
	case "ADVANCE_STATUS_PARTIALLY_SETTLED", "PARTIALLY_SETTLED":
		return l.PartiallySettled
	case "ADVANCE_STATUS_REFUNDED", "REFUNDED":
		return l.Refunded
	case "ADVANCE_STATUS_CANCELLED", "CANCELLED":
		return l.Cancelled
	case "ADVANCE_STATUS_EXPIRED", "EXPIRED":
		return l.Expired
	}
	return ""
}

// advanceProrationPolicyLabel resolves the AdvanceProrationPolicy proto-enum
// string. UNSPECIFIED collapses to FullTranche per Decision 13.
func advanceProrationPolicyLabel(s string, l centymo.AdvanceProrationPolicyLabels) string {
	switch s {
	case "ADVANCE_PRORATION_POLICY_DAY_PRORATED", "DAY_PRORATED":
		return l.DayProrated
	case "ADVANCE_PRORATION_POLICY_FULL_TRANCHE", "FULL_TRANCHE",
		"ADVANCE_PRORATION_POLICY_UNSPECIFIED", "UNSPECIFIED":
		return l.FullTranche
	case "ADVANCE_PRORATION_POLICY_NEXT_PERIOD_START", "NEXT_PERIOD_START":
		return l.NextPeriodStart
	}
	return ""
}

// NewView creates the collection detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection", "read") {
			return view.Forbidden("collection:read")
		}
		_ = perms
		id := viewCtx.Request.PathValue("id")

		resp, err := deps.ReadCollection(ctx, &collectionpb.ReadCollectionRequest{
			Data: &collectionpb.Collection{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read collection %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load collection: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			log.Printf("Collection %s not found", id)
			return view.Error(fmt.Errorf("collection not found"))
		}
		collection := collectionToMap(data[0], deps.AdvanceEnumLabels)

		refNumber, _ := collection["reference_number"].(string)

		l := deps.Labels
		headerTitle := l.Detail.TitlePrefix + refNumber

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}
		isAdvance, _ := collection["is_advance"].(bool)
		tabItems := buildTabItemsWithAdvance(l, id, deps.Routes, isAdvance, deps.AdvanceLabels.Tab)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "cash",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.PageTitle,
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "collection-detail-content",
			Collection:      collection,
			Labels:          l,
			AdvanceLabels:   deps.AdvanceLabels,
			AdvanceEnum:     deps.AdvanceEnumLabels,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
		}

		// 20260517-advance-cash-events Plan B Phase 7 — when this collection
		// is a MILESTONE advance, preload the linked BillingEvent rows for
		// the Advance Schedule tab. Cheap enough to always preload because
		// the advance-schedule tab is the most common destination for
		// MILESTONE advances.
		if advanceKind, _ := collection["advance_kind"].(string); advanceKind == "ADVANCE_KIND_MILESTONE" || advanceKind == "MILESTONE" {
			pageData.MilestoneLinks = loadMilestoneLinks(ctx, deps, id)
		}

		switch activeTab {
		case "info":
			// collection map has everything
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
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.QueryParams["cursor"]
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "collection",
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

		return view.OK("collection-detail", pageData)
	})
}

func buildTabItems(l centymo.CollectionLabels, id string, routes centymo.CollectionRoutes) []pyeza.TabItem {
	return buildTabItemsWithAdvance(l, id, routes, false, "")
}

// buildTabItemsWithAdvance builds the tab list, conditionally inserting the
// "Advance Schedule" tab when this TreasuryCollection carries an advance kind.
// 20260517-advance-cash-events Plan B Phase 4 — conditional tab is added via
// a helper not an `if` block so the renderer always sees a uniform slice.
// The advanceTabLabel is sourced from TreasuryAdvanceLabels.Tab (loaded from
// treasury_collection.json `advance.tab`); empty string falls back to the
// default English wording.
func buildTabItemsWithAdvance(l centymo.CollectionLabels, id string, routes centymo.CollectionRoutes, isAdvance bool, advanceTabLabel string) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	items := []pyeza.TabItem{
		{Key: "info", Label: l.Detail.TabBasicInfo, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
	}
	if isAdvance {
		label := advanceTabLabel
		if label == "" {
			label = "Advance Schedule"
		}
		items = append(items, pyeza.TabItem{Key: "advance-schedule", Label: label, Href: base + "?tab=advance-schedule", HxGet: action + "advance-schedule", Icon: "icon-calendar"})
	}
	items = append(items,
		pyeza.TabItem{Key: "attachments", Label: l.Detail.TabAttachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip"},
		pyeza.TabItem{Key: "audit", Label: l.Detail.TabAuditTrail, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock"},
		pyeza.TabItem{Key: "audit-history", Label: "History", Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	)
	return items
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection", "read") {
			return view.Forbidden("collection:read")
		}
		_ = perms
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		resp, err := deps.ReadCollection(ctx, &collectionpb.ReadCollectionRequest{
			Data: &collectionpb.Collection{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read collection %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load collection: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			log.Printf("Collection %s not found", id)
			return view.Error(fmt.Errorf("collection not found"))
		}
		collection := collectionToMap(data[0], deps.AdvanceEnumLabels)

		l := deps.Labels
		isAdvance, _ := collection["is_advance"].(bool)
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Collection:    collection,
			Labels:        l,
			AdvanceLabels: deps.AdvanceLabels,
			AdvanceEnum:   deps.AdvanceEnumLabels,
			ActiveTab:     tab,
			TabItems:      buildTabItemsWithAdvance(l, id, deps.Routes, isAdvance, deps.AdvanceLabels.Tab),
		}
		// 20260517-advance-cash-events Plan B Phase 7 — preload milestone
		// links when the active tab is advance-schedule and kind=MILESTONE.
		if tab == "advance-schedule" {
			if advanceKind, _ := collection["advance_kind"].(string); advanceKind == "ADVANCE_KIND_MILESTONE" || advanceKind == "MILESTONE" {
				pageData.MilestoneLinks = loadMilestoneLinks(ctx, deps, id)
			}
		}

		switch tab {
		case "info":
			// collection map has everything
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
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		case "audit-history":
			if deps.ListAuditHistory != nil {
				cursor := viewCtx.QueryParams["cursor"]
				auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
					EntityType:  "collection",
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

		templateName := "collection-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// buildAuditTable creates the audit trail table.
func buildAuditTable(l centymo.CollectionLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: l.Detail.Date, WidthClass: "col-5xl"},
		{Key: "action", Label: l.Detail.AuditAction},
		{Key: "user", Label: l.Detail.AuditUser, WidthClass: "col-6xl"},
	}

	rows := []types.TableRow{}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:                   "audit-trail-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowEntries:          true,
		DefaultSortColumn:    "date",
		DefaultSortDirection: "desc",
		Labels:               tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Detail.AuditEmptyTitle,
			Message: l.Detail.AuditEmptyMessage,
		},
	}
	types.ApplyTableSettings(cfg)

	return cfg
}

// loadMilestoneLinks fetches the treasury_collection_billing_event junction
// rows for this advance Collection, hydrates each row's BillingEvent state +
// per-event Recognize URL, and returns the per-row view shape rendered in
// the Advance Schedule tab.
//
// Nil-safe: empty slice when either dep is unwired.
//
// 20260517-advance-cash-events Plan B Phase 7.
func loadMilestoneLinks(ctx context.Context, deps *DetailViewDeps, collectionID string) []MilestoneLinkRow {
	if deps.ListTreasuryCollectionBillingEvents == nil {
		return nil
	}
	resp, err := deps.ListTreasuryCollectionBillingEvents(ctx, &junctionpb.ListTreasuryCollectionBillingEventsRequest{
		Filters: &commonpb.FilterRequest{
			Filters: []*commonpb.TypedFilter{
				{
					Field: "treasury_collection_id",
					FilterType: &commonpb.TypedFilter_StringFilter{
						StringFilter: &commonpb.StringFilter{
							Value:    collectionID,
							Operator: commonpb.StringOperator_STRING_EQUALS,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("loadMilestoneLinks: list junctions for %s: %v", collectionID, err)
		return nil
	}
	if resp == nil {
		return nil
	}
	out := make([]MilestoneLinkRow, 0, len(resp.GetData()))
	for _, j := range resp.GetData() {
		row := MilestoneLinkRow{
			BillingEventID: j.GetBillingEventId(),
			TrancheAmount:  j.GetTrancheAmount(),
			TrancheDisplay: fmt.Sprintf("%.2f", float64(j.GetTrancheAmount())/100),
			RevenueID:      j.GetRevenueId(),
		}
		recognized := strings.TrimSpace(row.RevenueID) != ""
		// Hydrate BillingEvent state when the dep is wired so the operator
		// can see status (READY / BILLED / WAIVED) and we know whether to
		// show the Recognize button.
		var beStatus billingeventpb.BillingEventStatus
		if deps.ReadBillingEvent != nil {
			beResp, err := deps.ReadBillingEvent(ctx, &billingeventpb.ReadBillingEventRequest{
				Data: &billingeventpb.BillingEvent{Id: row.BillingEventID},
			})
			if err == nil && beResp != nil && len(beResp.GetData()) > 0 {
				be := beResp.GetData()[0]
				row.SubscriptionID = be.GetSubscriptionId()
				row.Currency = be.GetBillingCurrency()
				beStatus = be.GetStatus()
			}
		}
		row.StatusKey, row.StatusLabel = milestoneStatusKeyLabel(beStatus)
		if !recognized && beStatus == billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_BILLED && deps.MilestoneRecognizeURL != "" && row.SubscriptionID != "" {
			row.ShowRecognize = true
			row.RecognizeURL = strings.ReplaceAll(strings.ReplaceAll(deps.MilestoneRecognizeURL, "{id}", row.SubscriptionID), "{eventId}", row.BillingEventID)
		}
		out = append(out, row)
	}
	return out
}

// milestoneStatusKeyLabel maps a BillingEventStatus to the (key, label) pair
// used by the template badge. The label is intentionally bare-English here
// because the Collection labels struct doesn't carry billing-event status
// translations; the operator-facing copy lives downstream in the BillingEvent
// detail page.
func milestoneStatusKeyLabel(s billingeventpb.BillingEventStatus) (string, string) {
	switch s {
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_READY:
		return "ready", "Ready"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_BILLED:
		return "billed", "Billed"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_WAIVED:
		return "waived", "Waived"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_DEFERRED:
		return "deferred", "Deferred"
	case billingeventpb.BillingEventStatus_BILLING_EVENT_STATUS_CANCELLED:
		return "cancelled", "Cancelled"
	default:
		return "pending", "Pending"
	}
}
