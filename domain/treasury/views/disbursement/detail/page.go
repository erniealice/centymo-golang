package detail

import (
	"context"
	"fmt"
	"log"
	"strings"

	treasury "github.com/erniealice/centymo-golang/domain/treasury"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	supplierbillingeventpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_billing_event"
	disbursementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement"
	// 20260517-advance-cash-events Plan B Phase 7 — junction reads.
	junctionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement_supplier_billing_event"
)

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	Routes           treasury.DisbursementRoutes
	ReadDisbursement func(ctx context.Context, req *disbursementpb.ReadDisbursementRequest) (*disbursementpb.ReadDisbursementResponse, error)
	Labels           treasury.DisbursementLabels
	CommonLabels     pyeza.CommonLabels
	TableLabels      types.TableLabels

	// 20260517-advance-cash-events Plan B Phase 4 — Advance Schedule tab
	// typed-label sources. AdvanceLabels carries the tab + column + section
	// strings; AdvanceEnumLabels carries the AdvanceKind / AdvanceStatus /
	// AdvanceProrationPolicy option labels rendered as badge text.
	AdvanceLabels     treasury.TreasuryAdvanceLabels
	AdvanceEnumLabels treasury.AdvanceEnumLabels

	// 20260517-advance-cash-events Plan B Phase 7 — MILESTONE cross-link.
	// When advance_kind == MILESTONE, the Advance Schedule tab lists the
	// linked SupplierBillingEvent rows + per-event Recognize button.
	// Both deps are nil-safe.
	ListDisbursementSupplierBillingEvents func(ctx context.Context, req *junctionpb.ListDisbursementSupplierBillingEventsRequest) (*junctionpb.ListDisbursementSupplierBillingEventsResponse, error)
	ReadSupplierBillingEvent              func(ctx context.Context, req *supplierbillingeventpb.ReadSupplierBillingEventRequest) (*supplierbillingeventpb.ReadSupplierBillingEventResponse, error)
	// SupplierBillingEventRecognizeURL is the route template for the
	// per-event Recognize POST (`/action/supplier-billing-event/recognize/{id}`).
	// Empty disables the button.
	SupplierBillingEventRecognizeURL string

	attachment.AttachmentOps
	auditlog.AuditOps
}

// MilestoneLinkRow is a single linked-SupplierBillingEvent row shown in the
// Advance Schedule tab when advance_kind == MILESTONE (buying side).
type MilestoneLinkRow struct {
	SupplierBillingEventID string
	SupplierSubscriptionID string
	TrancheAmount          int64
	TrancheDisplay         string
	Currency               string
	StatusKey              string
	StatusLabel            string
	ExpenseRecognitionID   string
	ShowRecognize          bool
	RecognizeURL           string
}

// PageData holds the data for the disbursement detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Disbursement    map[string]any
	Labels          treasury.DisbursementLabels
	// 20260517-advance-cash-events Plan B Phase 4 — typed-label sources for
	// the Advance Schedule tab.
	AdvanceLabels treasury.TreasuryAdvanceLabels
	AdvanceEnum   treasury.AdvanceEnumLabels
	ActiveTab     string
	TabItems      []pyeza.TabItem

	// Convenience fields for template rendering
	Reference     string
	StatusLabel   string
	StatusVariant string
	Amount        types.TableCell
	Currency      string

	AuditTable      *types.TableConfig
	AttachmentTable *types.TableConfig
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
	// 20260517-advance-cash-events Plan B Phase 7.
	MilestoneLinks []MilestoneLinkRow
}

// disbursementToMap converts a Disbursement protobuf to a map[string]any for template use.
//
// 20260517-advance-cash-events Plan B Phase 4 — also projects the advance_*
// fields so the detail template can render the Advance Schedule tab + the
// UNSCHEDULED Settle/Refund/Cancel CTAs. Phase 4 label cleanup adds resolved
// _label fields populated from the resolved AdvanceEnumLabels.
func disbursementToMap(d *disbursementpb.Disbursement, enumLabels treasury.AdvanceEnumLabels) map[string]any {
	advanceKind := d.GetAdvanceKind().String()
	advanceStatus := d.GetAdvanceStatus().String()
	advanceProrationPolicy := d.GetAdvanceProrationPolicy().String()
	isAdvance := advanceKind != "" &&
		advanceKind != "ADVANCE_KIND_UNSPECIFIED" &&
		advanceKind != "ADVANCE_KIND_NONE" &&
		advanceKind != "NONE"
	isUnscheduledActive := advanceKind == "ADVANCE_KIND_UNSCHEDULED" &&
		(advanceStatus == "ADVANCE_STATUS_ACTIVE" || advanceStatus == "ADVANCE_STATUS_PARTIALLY_SETTLED" ||
			advanceStatus == "ACTIVE" || advanceStatus == "PARTIALLY_SETTLED")
	return map[string]any{
		"id":                             d.GetId(),
		"name":                           d.GetName(),
		"reference_number":               d.GetReferenceNumber(),
		"amount":                         types.MoneyCell(float64(d.GetAmount()), d.GetCurrency(), true),
		"currency":                       d.GetCurrency(),
		"status":                         d.GetStatus(),
		"disbursement_method_id":         d.GetDisbursementMethodId(),
		"disbursement_type":              d.GetDisbursementType(),
		"expenditure_id":                 d.GetExpenditureId(),
		"approved_by":                    d.GetApprovedBy(),
		"active":                         d.GetActive(),
		"date_created_string":            d.GetDateCreatedString(),
		"date_modified_string":           d.GetDateModifiedString(),
		"advance_kind":                   advanceKind,
		"advance_status":                 advanceStatus,
		"advance_proration_policy":       advanceProrationPolicy,
		"advance_kind_label":             advanceKindLabel(advanceKind, enumLabels.Kind),
		"advance_status_label":           advanceStatusLabel(advanceStatus, enumLabels.Status),
		"advance_proration_policy_label": advanceProrationPolicyLabel(advanceProrationPolicy, enumLabels.ProrationPolicy),
		"advance_start_date":             d.GetAdvanceStartDate(),
		"advance_end_date":               d.GetAdvanceEndDate(),
		"advance_period_count":           d.GetAdvancePeriodCount(),
		"advance_period_unit":            d.GetAdvancePeriodUnit(),
		"advance_total_amount":           d.GetAdvanceTotalAmount(),
		"advance_remaining":              d.GetAdvanceRemainingAmount(),
		"advance_recognized":             d.GetAdvanceRecognizedAmount(),
		"advance_balance_account":        d.GetAdvanceBalanceAccountId(),
		"advance_target_account":         d.GetAdvanceTargetAccountId(),
		"is_advance":                     isAdvance,
		"is_unscheduled_active":          isUnscheduledActive,
	}
}

// advanceKindLabel resolves the AdvanceKind proto-enum string to its
// operator-facing label, accepting both the fully-qualified
// (ADVANCE_KIND_TIME_BASED) and bare (TIME_BASED) variants.
func advanceKindLabel(s string, l treasury.AdvanceKindLabels) string {
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
func advanceStatusLabel(s string, l treasury.AdvanceStatusLabels) string {
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
func advanceProrationPolicyLabel(s string, l treasury.AdvanceProrationPolicyLabels) string {
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

// NewView creates the disbursement detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement", "read") {
			return view.Forbidden("disbursement:read")
		}
		_ = perms
		id := viewCtx.Request.PathValue("id")

		resp, err := deps.ReadDisbursement(ctx, &disbursementpb.ReadDisbursementRequest{
			Data: &disbursementpb.Disbursement{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read disbursement %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load disbursement: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			log.Printf("Disbursement %s not found", id)
			return view.Error(fmt.Errorf("disbursement not found"))
		}
		record := data[0]
		disbursement := disbursementToMap(record, deps.AdvanceEnumLabels)

		refNumber := record.GetReferenceNumber()
		status := record.GetStatus()
		currency := record.GetCurrency()
		amount := types.MoneyCell(float64(record.GetAmount()), record.GetCurrency(), true)

		l := deps.Labels
		headerTitle := l.Detail.TitlePrefix + refNumber

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}
		isAdvance, _ := disbursement["is_advance"].(bool)
		tabItems := buildTabItemsWithAdvance(l, id, deps.Routes, isAdvance, deps.AdvanceLabels.Tab)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "cash",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.PageTitle,
				HeaderIcon:     "icon-arrow-up-right",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "disbursement-detail-content",
			Disbursement:    disbursement,
			Labels:          l,
			AdvanceLabels:   deps.AdvanceLabels,
			AdvanceEnum:     deps.AdvanceEnumLabels,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
			Reference:       refNumber,
			StatusLabel:     status,
			StatusVariant:   statusVariant(status),
			Amount:          amount,
			Currency:        currency,
		}

		// 20260517-advance-cash-events Plan B Phase 7 — preload MILESTONE
		// link rows for the Advance Schedule tab.
		if advanceKind, _ := disbursement["advance_kind"].(string); advanceKind == "ADVANCE_KIND_MILESTONE" || advanceKind == "MILESTONE" {
			pageData.MilestoneLinks = loadMilestoneLinks(ctx, deps, id)
		}

		// Load tab-specific data
		switch activeTab {
		case "info":
			// Disbursement map has everything
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
					EntityType:  "disbursement",
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

		return view.OK("disbursement-detail", pageData)
	})
}

func buildTabItems(l treasury.DisbursementLabels, id string, routes treasury.DisbursementRoutes) []pyeza.TabItem {
	return buildTabItemsWithAdvance(l, id, routes, false, "")
}

// buildTabItemsWithAdvance conditionally adds the "Advance Schedule" tab.
// 20260517-advance-cash-events Plan B Phase 4. The advanceTabLabel arrives
// from TreasuryAdvanceLabels.Tab (loaded from treasury_disbursement.json
// `advance.tab`); empty string falls back to the default English wording.
func buildTabItemsWithAdvance(l treasury.DisbursementLabels, id string, routes treasury.DisbursementRoutes, isAdvance bool, advanceTabLabel string) []pyeza.TabItem {
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
		pyeza.TabItem{Key: "audit-history", Label: l.Detail.TabAuditHistory, Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	)
	return items
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement", "read") {
			return view.Forbidden("disbursement:read")
		}
		_ = perms
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		resp, err := deps.ReadDisbursement(ctx, &disbursementpb.ReadDisbursementRequest{
			Data: &disbursementpb.Disbursement{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read disbursement %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load disbursement: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			log.Printf("Disbursement %s not found", id)
			return view.Error(fmt.Errorf("disbursement not found"))
		}
		record := data[0]
		disbursement := disbursementToMap(record, deps.AdvanceEnumLabels)

		status := record.GetStatus()
		currency := record.GetCurrency()
		amount := types.MoneyCell(float64(record.GetAmount()), record.GetCurrency(), true)
		refNumber := record.GetReferenceNumber()

		l := deps.Labels
		isAdvance, _ := disbursement["is_advance"].(bool)
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Disbursement:  disbursement,
			Labels:        l,
			AdvanceLabels: deps.AdvanceLabels,
			AdvanceEnum:   deps.AdvanceEnumLabels,
			ActiveTab:     tab,
			TabItems:      buildTabItemsWithAdvance(l, id, deps.Routes, isAdvance, deps.AdvanceLabels.Tab),
			Reference:     refNumber,
			StatusLabel:   status,
			StatusVariant: statusVariant(status),
			Amount:        amount,
			Currency:      currency,
		}
		// 20260517-advance-cash-events Plan B Phase 7 — preload milestone
		// links when the active tab is advance-schedule and kind=MILESTONE.
		if tab == "advance-schedule" {
			if advanceKind, _ := disbursement["advance_kind"].(string); advanceKind == "ADVANCE_KIND_MILESTONE" || advanceKind == "MILESTONE" {
				pageData.MilestoneLinks = loadMilestoneLinks(ctx, deps, id)
			}
		}

		switch tab {
		case "info":
			// disbursement map has everything
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
					EntityType:  "disbursement",
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

		templateName := "disbursement-tab-" + tab
		if tab == "attachments" {
			templateName = "attachment-tab"
		}
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

func buildAuditTable(l treasury.DisbursementLabels, tableLabels types.TableLabels) *types.TableConfig {
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

func statusVariant(status string) string {
	switch status {
	case "draft":
		return "default"
	case "pending":
		return "warning"
	case "approved":
		return "info"
	case "paid":
		return "success"
	case "cancelled":
		return "danger"
	case "overdue":
		return "danger"
	default:
		return "default"
	}
}

// loadMilestoneLinks fetches the
// disbursement_supplier_billing_event junction rows for this advance
// Disbursement, hydrates each row's SupplierBillingEvent state + per-event
// Recognize URL, and returns the per-row view shape rendered in the Advance
// Schedule tab.
//
// Nil-safe: empty slice when either dep is unwired.
//
// 20260517-advance-cash-events Plan B Phase 7.
func loadMilestoneLinks(ctx context.Context, deps *DetailViewDeps, disbursementID string) []MilestoneLinkRow {
	if deps.ListDisbursementSupplierBillingEvents == nil {
		return nil
	}
	resp, err := deps.ListDisbursementSupplierBillingEvents(ctx, &junctionpb.ListDisbursementSupplierBillingEventsRequest{
		Filters: &commonpb.FilterRequest{
			Filters: []*commonpb.TypedFilter{
				{
					Field: "treasury_disbursement_id",
					FilterType: &commonpb.TypedFilter_StringFilter{
						StringFilter: &commonpb.StringFilter{
							Value:    disbursementID,
							Operator: commonpb.StringOperator_STRING_EQUALS,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("loadMilestoneLinks: list junctions for %s: %v", disbursementID, err)
		return nil
	}
	if resp == nil {
		return nil
	}
	out := make([]MilestoneLinkRow, 0, len(resp.GetData()))
	for _, j := range resp.GetData() {
		row := MilestoneLinkRow{
			SupplierBillingEventID: j.GetSupplierBillingEventId(),
			TrancheAmount:          j.GetTrancheAmount(),
			TrancheDisplay:         fmt.Sprintf("%.2f", float64(j.GetTrancheAmount())/100),
			ExpenseRecognitionID:   j.GetExpenseRecognitionId(),
		}
		recognized := strings.TrimSpace(row.ExpenseRecognitionID) != ""
		var sbeStatus supplierbillingeventpb.SupplierBillingEventStatus
		if deps.ReadSupplierBillingEvent != nil {
			beResp, err := deps.ReadSupplierBillingEvent(ctx, &supplierbillingeventpb.ReadSupplierBillingEventRequest{
				Data: &supplierbillingeventpb.SupplierBillingEvent{Id: row.SupplierBillingEventID},
			})
			if err == nil && beResp != nil && len(beResp.GetData()) > 0 {
				be := beResp.GetData()[0]
				row.SupplierSubscriptionID = be.GetSupplierSubscriptionId()
				row.Currency = be.GetBillingCurrency()
				sbeStatus = be.GetStatus()
			}
		}
		row.StatusKey, row.StatusLabel = supplierMilestoneStatusKeyLabel(sbeStatus)
		if !recognized && sbeStatus == supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_BILLED && deps.SupplierBillingEventRecognizeURL != "" {
			row.ShowRecognize = true
			row.RecognizeURL = strings.ReplaceAll(deps.SupplierBillingEventRecognizeURL, "{id}", row.SupplierBillingEventID)
		}
		out = append(out, row)
	}
	return out
}

// supplierMilestoneStatusKeyLabel maps a SupplierBillingEventStatus to the
// (key, label) pair used by the template badge.
func supplierMilestoneStatusKeyLabel(s supplierbillingeventpb.SupplierBillingEventStatus) (string, string) {
	switch s {
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_READY:
		return "ready", "Ready"
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_BILLED:
		return "billed", "Billed"
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_WAIVED:
		return "waived", "Waived"
	case supplierbillingeventpb.SupplierBillingEventStatus_SUPPLIER_BILLING_EVENT_STATUS_CANCELLED:
		return "cancelled", "Cancelled"
	default:
		return "pending", "Pending"
	}
}
