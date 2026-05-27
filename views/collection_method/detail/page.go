package detail

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	cmpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method"
)

// DetailViewDeps holds dependencies for the collection method detail page.
type DetailViewDeps struct {
	Routes       centymo.CollectionMethodRoutes
	Labels       centymo.CollectionMethodLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ReadCollectionMethod is nil until the espyna use cases land. When nil,
	// the detail view redirects to the list (graceful degradation).
	ReadCollectionMethod func(ctx context.Context, req *cmpb.ReadCollectionMethodRequest) (*cmpb.ReadCollectionMethodResponse, error)

	// Stage 2 — Eligibility Rules tab closures (nil-safe).
	EligibilityRuleDeps *EligibilityRuleTabDeps
}

// PageData holds the template data for the collection method detail page.
type PageData struct {
	types.PageData
	ContentTemplate string

	Method        map[string]any
	KindRows      []KindRow // kind-specific config (Overview pane)
	StatusVariant string

	TabItems  []pyeza.TabItem
	ActiveTab string

	EditURL    string
	PublishURL string
	CloseURL   string
	ArchiveURL string
	ReviseURL  string

	// EligibilityTab carries the pre-loaded tab data when ActiveTab == "eligibility".
	// Nil for all other tabs (the template dispatches via eq .ActiveTab).
	EligibilityTab *EligibilityRuleTabData
}

// KindRow is a single label/value pair in the kind-specific Overview pane.
type KindRow struct {
	Label string
	Value string
}

// Stage-1 tab keys. Only Info / Versions / Activity carry content this stage
// (pages.md §B-5 tab table marks these Stage 1). The remaining tabs
// (Eligibility / Grants / Sub-status / Approvals / Instances / Profiles /
// Transitions) are TODO-stubbed for Stages 2/3/4/6 — they render a
// "coming later" pane so the IA is visible without dead links.
const (
	tabInfo        = "info"
	tabEligibility = "eligibility"
	tabGrants      = "grants"
	tabSubStatus   = "sub-status"
	tabApprovals   = "approvals"
	tabInstances   = "instances"
	tabProfiles    = "profiles"
	tabTransitions = "transitions"
	tabVersions    = "versions"
	tabActivity    = "activity"
)

// NewView creates the collection method detail page view (pages.md §B-5 detail).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method", "read") {
			return view.Forbidden("collection_method:read")
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" || deps.ReadCollectionMethod == nil {
			return view.Redirect(strings.Replace(deps.Routes.ListURL, "{status}", "active", 1))
		}

		resp, err := deps.ReadCollectionMethod(ctx, &cmpb.ReadCollectionMethodRequest{
			Data: &cmpb.CollectionMethod{Id: id},
		})
		if err != nil {
			log.Printf("ReadCollectionMethod %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load collection method: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("collection method not found"))
		}
		m := data[0]

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = tabInfo
		}

		pd := buildPageDataWithContext(ctx, deps, m, activeTab, viewCtx)
		return view.OK("collection-method-detail", pd)
	})
}

// NewTabAction handles HTMX tab swaps (/action/collection-method/detail/{id}/tab/{tab}).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method", "read") {
			return view.Forbidden("collection_method:read")
		}
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if id == "" || tab == "" {
			return view.Error(fmt.Errorf("missing id or tab"))
		}

		// Stage 2: the Eligibility tab is self-contained; delegate to its own view
		// so it can load rules independently without needing ReadCollectionMethod.
		if tab == tabEligibility && deps.EligibilityRuleDeps != nil {
			return NewEligibilityRuleTabView(deps.EligibilityRuleDeps).Handle(ctx, viewCtx)
		}

		if deps.ReadCollectionMethod == nil {
			return view.Error(fmt.Errorf("missing id or tab"))
		}
		resp, err := deps.ReadCollectionMethod(ctx, &cmpb.ReadCollectionMethodRequest{
			Data: &cmpb.CollectionMethod{Id: id},
		})
		if err != nil {
			return view.Error(fmt.Errorf("failed to load collection method: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("collection method not found"))
		}
		pd := buildPageDataWithContext(ctx, deps, data[0], tab, viewCtx)
		return view.OK("collection-method-tab-content", pd)
	})
}

func buildPageDataWithContext(ctx context.Context, deps *DetailViewDeps, m *cmpb.CollectionMethod, activeTab string, viewCtx *view.ViewContext) *PageData {
	l := deps.Labels
	id := m.GetId()

	methodMap := map[string]any{
		"id":               id,
		"name":             m.GetName(),
		"template_code":    m.GetTemplateCode(),
		"category":         enumShort(m.GetCategory().String()),
		"posting_kind":     enumShort(m.GetPostingKind().String()),
		"audience_mode":    enumShort(m.GetAudienceMode().String()),
		"tax_effect_kind":  enumShort(m.GetTaxEffectKind().String()),
		"lifecycle":        enumShort(m.GetLifecycle().String()),
		"source":           enumShort(m.GetSource().String()),
		"revision":         strconv.FormatInt(int64(m.GetRevision()), 10),
		"version_status":   enumShort(m.GetVersionStatus().String()),
		"balance_account":  m.GetBalanceAccountId(),
		"target_account":   m.GetTargetAccountId(),
		"eligibility_rule": m.GetDefaultEligibilityRuleId(),
	}

	tabItems := []pyeza.TabItem{
		{Key: tabInfo, Label: l.Tabs.Info},
		{Key: tabEligibility, Label: l.Tabs.Eligibility},
		{Key: tabGrants, Label: l.Tabs.Grants},
		{Key: tabSubStatus, Label: l.Tabs.SubStatusTags},
		{Key: tabApprovals, Label: l.Tabs.Approvals},
		{Key: tabInstances, Label: l.Tabs.Instances},
		{Key: tabProfiles, Label: l.Tabs.Profiles},
		{Key: tabTransitions, Label: l.Tabs.Transitions},
		{Key: tabVersions, Label: l.Tabs.Versions},
		{Key: tabActivity, Label: l.Tabs.Activity},
	}

	pd := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          m.GetName(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			HeaderTitle:    m.GetName(),
			HeaderSubtitle: l.Page.DetailSubtitle,
			HeaderIcon:     "icon-credit-card",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "collection-method-detail-content",
		Method:          methodMap,
		KindRows:        buildKindRows(m, l.Fragment),
		StatusVariant:   lifecycleVariant(m.GetLifecycle().String()),
		TabItems:        tabItems,
		ActiveTab:       activeTab,
		EditURL:         resolveID(deps.Routes.EditURL, id),
		PublishURL:      resolveID(deps.Routes.PublishURL, id),
		CloseURL:        resolveID(deps.Routes.CloseURL, id),
		ArchiveURL:      resolveID(deps.Routes.ArchiveURL, id),
		ReviseURL:       resolveID(deps.Routes.ReviseURL, id),
	}

	// Stage 2: populate EligibilityTab when the active tab is "eligibility".
	// Rules are loaded inline here so the full-page (non-HTMX) detail load works.
	if activeTab == tabEligibility && deps.EligibilityRuleDeps != nil {
		addURL := route.ResolveURL(deps.Routes.EligibilityRuleAddURL, "method_id", id)
		rules, err := loadEligibilityRules(ctx, deps.EligibilityRuleDeps, id)
		if err != nil {
			log.Printf("EligibilityTab loadRules method=%s: %v", id, err)
		}
		pd.EligibilityTab = &EligibilityRuleTabData{
			Labels:       l.EligibilityRule,
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			MethodID:     id,
			Rules:        rules,
			AddURL:       addURL,
		}
	}

	return pd
}

// buildKindRows renders the kind-specific Overview from the template_details oneof.
func buildKindRows(m *cmpb.CollectionMethod, fl centymo.CollectionMethodFragmentLabels) []KindRow {
	var rows []KindRow
	if vp := m.GetVoucherProgram(); vp != nil {
		if vp.GetDefaultFaceValueCentavos() > 0 {
			rows = append(rows, KindRow{Label: fl.DefaultFaceValue, Value: formatCentavos(vp.GetDefaultFaceValueCentavos())})
		}
		if vp.GetDefaultExpiryDays() > 0 {
			rows = append(rows, KindRow{Label: fl.DefaultExpiryDays, Value: strconv.FormatInt(int64(vp.GetDefaultExpiryDays()), 10)})
		}
		rows = append(rows, KindRow{Label: fl.AllowedBearerModes, Value: enumShort(vp.GetAllowedBearerModes().String())})
	} else if ap := m.GetAdvanceProgram(); ap != nil {
		rows = append(rows, KindRow{Label: fl.AdvanceKind, Value: enumShort(ap.GetAdvanceKind().String())})
		if ap.GetDefaultBalanceAccountId() != "" {
			rows = append(rows, KindRow{Label: fl.DefaultBalanceAcct, Value: ap.GetDefaultBalanceAccountId()})
		}
		if ap.GetDefaultTargetAccountId() != "" {
			rows = append(rows, KindRow{Label: fl.DefaultTargetAcct, Value: ap.GetDefaultTargetAccountId()})
		}
		if ap.GetDefaultPeriodCount() > 0 {
			rows = append(rows, KindRow{Label: fl.DefaultPeriodCount, Value: strconv.FormatInt(int64(ap.GetDefaultPeriodCount()), 10)})
		}
		if ap.GetDefaultPeriodUnit() != "" {
			rows = append(rows, KindRow{Label: fl.DefaultPeriodUnit, Value: ap.GetDefaultPeriodUnit()})
		}
	}
	// CARD (CollectionMethodCardTypeDetails) carries no template fields (D-4.26).
	return rows
}

// --- helpers -----------------------------------------------------------------

func enumShort(s string) string {
	if s == "" {
		return ""
	}
	markers := []string{"_CATEGORY_", "_POSTING_KIND_", "_AUDIENCE_MODE_", "_LIFECYCLE_", "_SOURCE_", "_TAX_EFFECT_KIND_", "_VERSION_STATUS_", "_BEARER_MODE_", "ADVANCE_KIND_"}
	for _, mk := range markers {
		if i := strings.Index(s, mk); i >= 0 {
			return s[i+len(mk):]
		}
	}
	return s
}

func lifecycleVariant(lifecycle string) string {
	switch lifecycle {
	case "COLLECTION_METHOD_LIFECYCLE_ACTIVE":
		return "success"
	case "COLLECTION_METHOD_LIFECYCLE_DRAFT":
		return "default"
	case "COLLECTION_METHOD_LIFECYCLE_CLOSED":
		return "warning"
	case "COLLECTION_METHOD_LIFECYCLE_ARCHIVED":
		return "danger"
	default:
		return "default"
	}
}

func resolveID(template, id string) string {
	if template == "" {
		return ""
	}
	return route.ResolveURL(template, "id", id)
}

func formatCentavos(v int64) string {
	return strconv.FormatFloat(float64(v)/100.0, 'f', 2, 64)
}
