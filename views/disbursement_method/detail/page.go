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

	dmpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement_method"
)

// DetailViewDeps holds dependencies for the disbursement method detail page.
type DetailViewDeps struct {
	Routes       centymo.DisbursementMethodRoutes
	Labels       centymo.DisbursementMethodLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ReadDisbursementMethod is nil until the espyna use cases land.
	ReadDisbursementMethod func(ctx context.Context, req *dmpb.ReadDisbursementMethodRequest) (*dmpb.ReadDisbursementMethodResponse, error)
}

// PageData holds the template data for the disbursement method detail page.
type PageData struct {
	types.PageData
	ContentTemplate string

	Method        map[string]any
	KindRows      []KindRow
	StatusVariant string

	TabItems  []pyeza.TabItem
	ActiveTab string

	EditURL    string
	PublishURL string
	CloseURL   string
	ArchiveURL string
	ReviseURL  string
}

// KindRow is a single label/value pair in the kind-specific Overview pane.
type KindRow struct {
	Label string
	Value string
}

// Buying-side tabs (pages.md §C-5): 7 tabs, lighter than CM (no Eligibility /
// Grants / Sub-status). Stage 1 carries Info / Versions / Activity; the rest
// (Approval Rules / Instances / Disbursement Profiles / Transition Requests)
// are TODO-stubbed for Stages 4/6.
const (
	tabInfo        = "info"
	tabApprovals   = "approvals"
	tabInstances   = "instances"
	tabProfiles    = "profiles"
	tabTransitions = "transitions"
	tabVersions    = "versions"
	tabActivity    = "activity"
)

// NewView creates the disbursement method detail page view (pages.md §C-5 detail).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement_method", "read") {
			return view.Forbidden("disbursement_method:read")
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" || deps.ReadDisbursementMethod == nil {
			return view.Redirect(strings.Replace(deps.Routes.ListURL, "{status}", "active", 1))
		}

		resp, err := deps.ReadDisbursementMethod(ctx, &dmpb.ReadDisbursementMethodRequest{
			Data: &dmpb.DisbursementMethod{Id: id},
		})
		if err != nil {
			log.Printf("ReadDisbursementMethod %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load disbursement method: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("disbursement method not found"))
		}

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = tabInfo
		}
		return view.OK("disbursement-method-detail", buildPageData(deps, data[0], activeTab, viewCtx))
	})
}

// NewTabAction handles HTMX tab swaps (/action/disbursement-method/detail/{id}/tab/{tab}).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement_method", "read") {
			return view.Forbidden("disbursement_method:read")
		}
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if id == "" || tab == "" || deps.ReadDisbursementMethod == nil {
			return view.Error(fmt.Errorf("missing id or tab"))
		}
		resp, err := deps.ReadDisbursementMethod(ctx, &dmpb.ReadDisbursementMethodRequest{
			Data: &dmpb.DisbursementMethod{Id: id},
		})
		if err != nil {
			return view.Error(fmt.Errorf("failed to load disbursement method: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("disbursement method not found"))
		}
		return view.OK("disbursement-method-tab-content", buildPageData(deps, data[0], tab, viewCtx))
	})
}

func buildPageData(deps *DetailViewDeps, m *dmpb.DisbursementMethod, activeTab string, viewCtx *view.ViewContext) *PageData {
	l := deps.Labels
	id := m.GetId()

	methodMap := map[string]any{
		"id":              id,
		"name":            m.GetName(),
		"template_code":   m.GetTemplateCode(),
		"category":        enumShort(m.GetCategory().String()),
		"posting_kind":    enumShort(m.GetPostingKind().String()),
		"tax_effect_kind": enumShort(m.GetTaxEffectKind().String()),
		"lifecycle":       enumShort(m.GetLifecycle().String()),
		"source":          enumShort(m.GetSource().String()),
		"revision":        strconv.FormatInt(int64(m.GetRevision()), 10),
		"version_status":  enumShort(m.GetVersionStatus().String()),
		"balance_account": m.GetBalanceAccountId(),
		"target_account":  m.GetTargetAccountId(),
	}

	tabItems := []pyeza.TabItem{
		{Key: tabInfo, Label: l.Tabs.Info},
		{Key: tabApprovals, Label: l.Tabs.Approvals},
		{Key: tabInstances, Label: l.Tabs.Instances},
		{Key: tabProfiles, Label: l.Tabs.Profiles},
		{Key: tabTransitions, Label: l.Tabs.Transitions},
		{Key: tabVersions, Label: l.Tabs.Versions},
		{Key: tabActivity, Label: l.Tabs.Activity},
	}

	return &PageData{
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
		ContentTemplate: "disbursement-method-detail-content",
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
}

func buildKindRows(m *dmpb.DisbursementMethod, fl centymo.DisbursementMethodFragmentLabels) []KindRow {
	var rows []KindRow
	if ap := m.GetAdvanceProgram(); ap != nil {
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
	} else if ba := m.GetBankAccount(); ba != nil {
		if ba.GetBankName() != "" {
			rows = append(rows, KindRow{Label: fl.BankName, Value: ba.GetBankName()})
		}
	}
	return rows
}

// --- helpers -----------------------------------------------------------------

func enumShort(s string) string {
	if s == "" {
		return ""
	}
	markers := []string{"_CATEGORY_", "_POSTING_KIND_", "_LIFECYCLE_", "_SOURCE_", "_TAX_EFFECT_KIND_", "_VERSION_STATUS_", "ADVANCE_KIND_"}
	for _, mk := range markers {
		if i := strings.Index(s, mk); i >= 0 {
			return s[i+len(mk):]
		}
	}
	return s
}

func lifecycleVariant(lifecycle string) string {
	switch lifecycle {
	case "DISBURSEMENT_METHOD_LIFECYCLE_ACTIVE":
		return "success"
	case "DISBURSEMENT_METHOD_LIFECYCLE_DRAFT":
		return "default"
	case "DISBURSEMENT_METHOD_LIFECYCLE_CLOSED":
		return "warning"
	case "DISBURSEMENT_METHOD_LIFECYCLE_ARCHIVED":
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
