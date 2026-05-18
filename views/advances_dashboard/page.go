// Package advancesdashboard renders the cash-app "Advances Dashboard" page —
// a workspace-level summary of every advance Collection (deferred / liability
// inflow) and every advance Disbursement (prepaid / asset outflow) along with
// utilization bars and per-row links into the existing Treasury detail pages.
//
// This is the v1 surface for the 20260517-advance-cash-events plan; the page
// is fed by view-typed callbacks the block layer constructs from the espyna
// consumer surface so this package never imports espyna internals.
package advancesdashboard

import (
	"context"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// AdvanceRow is the view-layer shape for one row on the Advances Dashboard.
// Numeric amounts are centavos (int64) per the centymo-wide convention.
type AdvanceRow struct {
	ID                string
	ReferenceNumber   string
	CounterpartyName  string // customer (inflow) or supplier (outflow)
	Kind              string // lyngua-mapped (TIME_BASED / MILESTONE / UNSCHEDULED / NONE)
	KindRaw           string // raw enum string for data-attrs / variants
	Status            string // lyngua-mapped
	StatusRaw         string // raw enum string
	Currency          string
	TotalAmount       int64
	RemainingAmount   int64
	RecognizedAmount  int64
	UtilizationPct    int // 0-100; recognized / total
	DetailURL         string
}

// AdvancesPosition is the view-layer summary across both directions for the
// page header / stat cards.
type AdvancesPosition struct {
	OutflowTotalRemaining   int64 // Σ prepaid (asset) remaining
	InflowTotalRemaining    int64 // Σ deferred (liability) remaining
	OutflowActiveCount      int
	InflowActiveCount       int
	OutflowFullyRecognized  int
	InflowFullyRecognized   int
	Currency                string // workspace functional currency for display
}

// DashboardRequest is the input shape for the dashboard data callback.
type DashboardRequest struct {
	AsOfDate string // ISO YYYY-MM-DD; empty = today
}

// DashboardResponse bundles outflow + inflow rows and the position summary.
type DashboardResponse struct {
	Outflows AdvancesSection
	Inflows  AdvancesSection
	Position AdvancesPosition
}

// AdvancesSection holds the per-direction (outflow or inflow) rows.
type AdvancesSection struct {
	Rows []AdvanceRow
}

// ModuleDeps holds the dependencies for the advances-dashboard module.
type ModuleDeps struct {
	Routes       centymo.TreasuryAdvancesRoutes
	Labels       centymo.AdvancesDashboardLabels
	EnumLabels   centymo.AdvanceEnumLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// CollectionDetailURLTemplate / DisbursementDetailURLTemplate are URL
	// templates (e.g. "/app/collections/detail/{id}") used to deep-link
	// each row to its underlying TreasuryCollection / TreasuryDisbursement
	// detail page. Empty values render rows as non-clickable.
	CollectionDetailURLTemplate   string
	DisbursementDetailURLTemplate string

	// GetDashboard is the view-typed callback that returns the dashboard
	// data. Nil-safe: when unset, the view renders empty-state.
	GetDashboard func(ctx context.Context, req DashboardRequest) (*DashboardResponse, error)

	// GetFunctionalCurrency returns the workspace ISO 4217 currency code.
	// Nil-safe — when absent, money strings omit the currency prefix.
	GetFunctionalCurrency func(ctx context.Context) string
}

// Module holds the constructed advances-dashboard view.
type Module struct {
	routes centymo.TreasuryAdvancesRoutes
	Page   view.View
}

// NewModule constructs the advances-dashboard module from the given deps.
func NewModule(deps *ModuleDeps) *Module {
	return &Module{
		routes: deps.Routes,
		Page:   newView(deps),
	}
}

// RegisterRoutes registers the advances-dashboard routes on the given
// route registrar. Only the dashboard URL is registered here — the deep-links
// to the underlying Collection / Disbursement lists are served by those
// modules' existing handlers.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	if m.routes.DashboardURL != "" {
		r.GET(m.routes.DashboardURL, m.Page)
	}
}

// PageData carries the rendered values to the template.
type PageData struct {
	types.PageData
	ContentTemplate    string
	Labels             centymo.AdvancesDashboardLabels
	EnumLabels         centymo.AdvanceEnumLabels
	AsOfDate           string
	HasData            bool
	OutflowTable       *types.TableConfig
	InflowTable        *types.TableConfig
	OutflowRows        []AdvanceRow
	InflowRows         []AdvanceRow
	Position           AdvancesPosition
	OutflowSectionTitle string
	InflowSectionTitle  string
}

// newView builds the GET handler for the dashboard page.
func newView(deps *ModuleDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		// Page is gated by treasury_collection:list OR treasury_disbursement:list;
		// either grants read access to the workspace advances surface. Empty
		// permission codes mean the gate degrades open (legacy behavior).
		canCollection := perms.Can("treasury_collection", "list") || perms.Can("collection", "list")
		canDisbursement := perms.Can("treasury_disbursement", "list") || perms.Can("disbursement", "list")
		if !canCollection && !canDisbursement {
			return view.Forbidden("treasury_collection:list|treasury_disbursement:list")
		}

		asOf := viewCtx.QueryParams["as_of"]

		// Read workspace currency for money formatting (nil-safe).
		currency := ""
		if deps.GetFunctionalCurrency != nil {
			currency = deps.GetFunctionalCurrency(ctx)
		}

		// Fetch dashboard data (nil-safe — empty state when callback unset).
		var resp *DashboardResponse
		if deps.GetDashboard != nil {
			r, err := deps.GetDashboard(ctx, DashboardRequest{AsOfDate: asOf})
			if err != nil {
				log.Printf("advances_dashboard: GetDashboard failed: %v", err)
			} else {
				resp = r
			}
		}
		if resp == nil {
			resp = &DashboardResponse{}
		}
		if resp.Position.Currency == "" {
			resp.Position.Currency = currency
		}

		l := deps.Labels

		// Build per-direction tables. Each table uses pyeza.TableConfig so the
		// shared table-card template handles search / sort / empty state.
		outflowTable := buildSectionTable(
			"outflow",
			l.Outflow,
			deps.TableLabels,
			resp.Outflows.Rows,
		)
		inflowTable := buildSectionTable(
			"inflow",
			l.Inflow,
			deps.TableLabels,
			resp.Inflows.Rows,
		)

		hasData := len(resp.Outflows.Rows) > 0 || len(resp.Inflows.Rows) > 0

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "cash",
				ActiveSubNav:   "advances-dashboard",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.AsOfLabel,
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:     "advances-dashboard-content",
			Labels:              l,
			EnumLabels:          deps.EnumLabels,
			AsOfDate:            asOf,
			HasData:             hasData,
			OutflowTable:        outflowTable,
			InflowTable:         inflowTable,
			OutflowRows:         resp.Outflows.Rows,
			InflowRows:          resp.Inflows.Rows,
			Position:            resp.Position,
			OutflowSectionTitle: l.Outflow.CardTitle,
			InflowSectionTitle:  l.Inflow.CardTitle,
		}

		// Suppress unused-variable warnings when the lists are empty and only
		// the deep-link templates would otherwise reference deps.
		_ = deps.CollectionDetailURLTemplate
		_ = deps.DisbursementDetailURLTemplate

		return view.OK("advances-dashboard", pageData)
	})
}

// buildSectionTable constructs a pyeza.TableConfig for one half of the
// dashboard (outflow or inflow). The columns are the same; only the
// counterparty header text differs (set in the labels payload).
func buildSectionTable(side string, sectionLabels centymo.AdvancesDashboardSectionLabels, tableLabels types.TableLabels, rows []AdvanceRow) *types.TableConfig {
	cols := []types.TableColumn{
		{Key: "id", Label: sectionLabels.Table.ID},
		{Key: "counterparty", Label: sectionLabels.Table.Counterparty},
		{Key: "kind", Label: sectionLabels.Table.Kind, WidthClass: "col-3xl"},
		{Key: "total", Label: sectionLabels.Table.Total, WidthClass: "col-3xl", Align: "right"},
		{Key: "remaining", Label: sectionLabels.Table.Remaining, WidthClass: "col-3xl", Align: "right"},
		{Key: "utilization", Label: "Utilization", WidthClass: "col-3xl"}, // TODO(advance-cash-events): wire to AdvancesDashboardLabels.UtilizationLabel once lyngua-cascade lands.
		{Key: "status", Label: sectionLabels.Table.Status, WidthClass: "col-2xl"},
	}

	body := make([]types.TableRow, 0, len(rows))
	for _, r := range rows {
		statusVariant := advanceStatusVariant(r.StatusRaw)
		body = append(body, types.TableRow{
			ID:   r.ID,
			Href: r.DetailURL,
			Cells: []types.TableCell{
				{Type: "text", Value: r.ReferenceNumber},
				{Type: "text", Value: r.CounterpartyName},
				{Type: "badge", Value: r.Kind, Variant: advanceKindVariant(r.KindRaw)},
				types.MoneyCell(float64(r.TotalAmount), r.Currency, true),
				types.MoneyCell(float64(r.RemainingAmount), r.Currency, true),
				{Type: "text", Value: utilizationCellText(r.UtilizationPct)},
				{Type: "badge", Value: r.Status, Variant: statusVariant},
			},
			DataAttrs: map[string]string{
				"advance-id":         r.ID,
				"advance-kind":       r.KindRaw,
				"advance-status":     r.StatusRaw,
				"utilization-pct":    utilizationCellText(r.UtilizationPct),
				"counterparty-name":  r.CounterpartyName,
				"reference-number":   r.ReferenceNumber,
			},
		})
	}
	types.ApplyColumnStyles(cols, body)

	cfg := &types.TableConfig{
		ID:                   "advances-dashboard-" + side + "-table",
		Columns:              cols,
		Rows:                 body,
		ShowSearch:           true,
		ShowEntries:          true,
		DefaultSortColumn:    "remaining",
		DefaultSortDirection: "desc",
		Labels:               tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   sectionLabels.EmptyTitle,
			Message: sectionLabels.EmptyMessage,
		},
	}
	types.ApplyTableSettings(cfg)
	return cfg
}

// utilizationCellText renders a "{pct}%" string for the utilization column.
// The template upgrades this into a pyeza-progress bar on the dashboard.
func utilizationCellText(pct int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	// Manual itoa to avoid strconv import noise.
	if pct == 0 {
		return "0%"
	}
	digits := []byte{}
	for n := pct; n > 0; n /= 10 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
	}
	return string(digits) + "%"
}

// advanceKindVariant maps an AdvanceKind enum string to a badge variant.
// Kept simple — actual color palette lives in pyeza CSS.
func advanceKindVariant(kind string) string {
	switch kind {
	case "TIME_BASED", "ADVANCE_KIND_TIME_BASED":
		return "info"
	case "MILESTONE", "ADVANCE_KIND_MILESTONE":
		return "primary"
	case "UNSCHEDULED", "ADVANCE_KIND_UNSCHEDULED":
		return "warning"
	case "BURN_DOWN", "ADVANCE_KIND_BURN_DOWN":
		return "default"
	default:
		return "default"
	}
}

// advanceStatusVariant maps an AdvanceStatus enum string to a badge variant.
func advanceStatusVariant(status string) string {
	switch status {
	case "ACTIVE", "ADVANCE_STATUS_ACTIVE":
		return "success"
	case "FULLY_RECOGNIZED", "FULLY_AMORTIZED", "FULLY_DRAWN":
		return "default"
	case "SETTLED":
		return "info"
	case "PARTIALLY_SETTLED":
		return "warning"
	case "REFUNDED":
		return "info"
	case "CANCELLED":
		return "danger"
	case "EXPIRED":
		return "danger"
	default:
		return "default"
	}
}

// Suppress an unused-imports warning when pyeza.route is not referenced —
// keep the import handy for future deep-link enrichment.
var _ = route.ResolveURL
