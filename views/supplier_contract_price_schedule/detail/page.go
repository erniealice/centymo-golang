package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	suppliercontractlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_line"
	scpspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule"
	scpslpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule_line"
)

const (
	tabInfo     = "info"
	tabLines    = "lines"
	tabActivity = "activity"
)

// DetailViewDeps holds all dependencies for the SCPS detail page.
type DetailViewDeps struct {
	Routes       centymo.SupplierContractPriceScheduleRoutes
	Labels       centymo.SupplierContractPriceScheduleLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadSupplierContractPriceSchedule       func(ctx context.Context, req *scpspb.ReadSupplierContractPriceScheduleRequest) (*scpspb.ReadSupplierContractPriceScheduleResponse, error)
	ListSupplierContractPriceScheduleLines  func(ctx context.Context, req *scpslpb.ListSupplierContractPriceScheduleLinesRequest) (*scpslpb.ListSupplierContractPriceScheduleLinesResponse, error)
	ListSupplierContractLines               func(ctx context.Context, req *suppliercontractlinepb.ListSupplierContractLinesRequest) (*suppliercontractlinepb.ListSupplierContractLinesResponse, error)
}

// PageData holds the template data for the SCPS detail page.
type PageData struct {
	types.PageData
	ContentTemplate string

	Schedule      map[string]any
	StatusVariant string

	TabItems  []pyeza.TabItem
	ActiveTab string

	// Lines tab
	LineItemTable  *types.TableConfig
	LineItemAddURL string

	// Action URLs
	EditURL      string
	ActivateURL  string
	SupersedeURL string
	DeleteURL    string
}

// NewView creates the SCPS detail page view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Redirect(deps.Routes.ListURL)
		}
		schedule, err := readSchedule(ctx, deps, id)
		if err != nil {
			return view.Error(err)
		}

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = tabInfo
		}

		l := deps.Labels
		pd := buildPageData(viewCtx, deps, l, schedule, activeTab, id)

		switch activeTab {
		case tabLines:
			pd.LineItemTable = buildLineItemTable(ctx, deps, id, l)
			pd.LineItemAddURL = route.ResolveURL(deps.Routes.LineAddURL, "id", id)
		}

		return view.OK("supplier-contract-price-schedule-detail", pd)
	})
}

// NewTabAction handles HTMX tab switch requests.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if id == "" || tab == "" {
			return view.Error(fmt.Errorf("missing id or tab"))
		}
		schedule, err := readSchedule(ctx, deps, id)
		if err != nil {
			return view.Error(err)
		}
		l := deps.Labels
		pd := buildPageData(viewCtx, deps, l, schedule, tab, id)
		switch tab {
		case tabLines:
			pd.LineItemTable = buildLineItemTable(ctx, deps, id, l)
			pd.LineItemAddURL = route.ResolveURL(deps.Routes.LineAddURL, "id", id)
		}
		return view.OK("supplier-contract-price-schedule-tab-content", pd)
	})
}

// --- helpers -----------------------------------------------------------------

func readSchedule(ctx context.Context, deps *DetailViewDeps, id string) (*scpspb.SupplierContractPriceSchedule, error) {
	resp, err := deps.ReadSupplierContractPriceSchedule(ctx, &scpspb.ReadSupplierContractPriceScheduleRequest{
		Data: &scpspb.SupplierContractPriceSchedule{Id: id},
	})
	if err != nil {
		log.Printf("ReadSupplierContractPriceSchedule %s: %v", id, err)
		return nil, fmt.Errorf("failed to load price schedule: %w", err)
	}
	data := resp.GetData()
	if len(data) == 0 {
		return nil, fmt.Errorf("price schedule not found")
	}
	return data[0], nil
}

func buildPageData(viewCtx *view.ViewContext, deps *DetailViewDeps, l centymo.SupplierContractPriceScheduleLabels, schedule *scpspb.SupplierContractPriceSchedule, activeTab, id string) *PageData {
	startStr := ""
	if t := schedule.GetDateTimeStart(); t != nil && t.IsValid() {
		startStr = t.AsTime().UTC().Format("2006-01-02")
	}
	endStr := "—"
	if t := schedule.GetDateTimeEnd(); t != nil && t.IsValid() {
		endStr = t.AsTime().UTC().Format("2006-01-02")
	}

	scheduleMap := map[string]any{
		"id":                   schedule.GetId(),
		"name":                 schedule.GetName(),
		"description":          schedule.GetDescription(),
		"internal_id":          schedule.GetInternalId(),
		"status":               schedule.GetStatus().String(),
		"status_label":         scheduleStatusLabel(schedule.GetStatus().String(), l),
		"supplier_contract_id": schedule.GetSupplierContractId(),
		"date_start":           startStr,
		"date_end":             endStr,
		"sequence_number":      schedule.GetSequenceNumber(),
		"currency":             schedule.GetCurrency(),
		"location_id":          schedule.GetLocationId(),
		"notes":                schedule.GetNotes(),
	}

	tabItems := []pyeza.TabItem{
		{Key: tabInfo, Label: l.Tabs.Info},
		{Key: tabLines, Label: l.Tabs.Lines},
		{Key: tabActivity, Label: l.Tabs.Activity},
	}

	return &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          schedule.GetName(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			HeaderTitle:    schedule.GetName(),
			HeaderSubtitle: l.Detail.Title,
			HeaderIcon:     "icon-calendar",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "supplier-contract-price-schedule-detail-content",
		Schedule:        scheduleMap,
		StatusVariant:   scheduleStatusVariant(schedule.GetStatus().String()),
		TabItems:        tabItems,
		ActiveTab:       activeTab,
		EditURL:         route.ResolveURL(deps.Routes.EditURL, "id", id),
		ActivateURL:     route.ResolveURL(deps.Routes.ActivateURL, "id", id),
		SupersedeURL:    route.ResolveURL(deps.Routes.SupersedeURL, "id", id),
		DeleteURL:       deps.Routes.DeleteURL,
	}
}

func scheduleStatusVariant(status string) string {
	switch status {
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SCHEDULED":
		return "info"
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_ACTIVE":
		return "success"
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SUPERSEDED":
		return "default"
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_CANCELLED":
		return "danger"
	default:
		return "default"
	}
}

func scheduleStatusLabel(status string, l centymo.SupplierContractPriceScheduleLabels) string {
	switch status {
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SCHEDULED":
		return l.Status.Scheduled
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_ACTIVE":
		return l.Status.Active
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SUPERSEDED":
		return l.Status.Superseded
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_CANCELLED":
		return l.Status.Cancelled
	default:
		return status
	}
}

func buildLineItemTable(ctx context.Context, deps *DetailViewDeps, scheduleID string, l centymo.SupplierContractPriceScheduleLabels) *types.TableConfig {
	if deps.ListSupplierContractPriceScheduleLines == nil {
		return nil
	}
	idPtr := scheduleID
	resp, err := deps.ListSupplierContractPriceScheduleLines(ctx, &scpslpb.ListSupplierContractPriceScheduleLinesRequest{
		SupplierContractPriceScheduleId: &idPtr,
	})
	if err != nil {
		log.Printf("ListSupplierContractPriceScheduleLines for %s: %v", scheduleID, err)
		return nil
	}

	// Build a contract-line description map so the line table can show
	// human-readable names instead of raw line IDs.
	contractLineNames := map[string]string{}
	if deps.ListSupplierContractLines != nil {
		if scLinesResp, err := deps.ListSupplierContractLines(ctx, &suppliercontractlinepb.ListSupplierContractLinesRequest{}); err == nil {
			for _, line := range scLinesResp.GetData() {
				contractLineNames[line.GetId()] = line.GetDescription()
			}
		}
	}

	columns := []types.TableColumn{
		{Key: "contract_line", Label: l.Lines.ColumnContractLine},
		{Key: "unit_price", Label: l.Lines.ColumnUnitPrice, Align: "right", WidthClass: "col-3xl"},
		{Key: "quantity", Label: l.Lines.ColumnQuantity, Align: "right", WidthClass: "col-xl"},
		{Key: "minimum_amount", Label: l.Lines.ColumnMinimumAmount, Align: "right", WidthClass: "col-3xl"},
		{Key: "currency", Label: l.Lines.ColumnCurrency, WidthClass: "col-xs"},
		{Key: "cycle_override", Label: l.Lines.ColumnCycleOverride, WidthClass: "col-2xl"},
	}

	rows := make([]types.TableRow, 0, len(resp.GetData()))
	for _, line := range resp.GetData() {
		contractLineLabel := contractLineNames[line.GetSupplierContractLineId()]
		if contractLineLabel == "" {
			contractLineLabel = line.GetSupplierContractLineId()
		}
		minStr := "—"
		if line.MinimumAmount != nil {
			minStr = fmt.Sprintf("%.2f", float64(line.GetMinimumAmount())/100.0)
		}
		qtyStr := "—"
		if line.Quantity != nil {
			qtyStr = fmt.Sprintf("%.2f", line.GetQuantity())
		}
		cycleStr := "—"
		if line.CycleValueOverride != nil && line.CycleUnitOverride != nil {
			cycleStr = fmt.Sprintf("%d %s", line.GetCycleValueOverride(), line.GetCycleUnitOverride())
		}
		rows = append(rows, types.TableRow{
			ID: line.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: contractLineLabel},
				types.MoneyCell(float64(line.GetUnitPrice()), line.GetCurrency(), true),
				{Type: "number", Value: qtyStr},
				{Type: "number", Value: minStr},
				{Type: "text", Value: line.GetCurrency()},
				{Type: "text", Value: cycleStr},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:          "schedule-lines-table",
		Columns:     columns,
		Rows:        rows,
		ShowActions: true,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Lines.Empty,
		},
	}
}

// NewActivateAction is wired by the action package — but the redirect helper
// can be placed here too if we ever need a detail-page invocation. Provided
// as a no-op safety net in case block.go pulls from detail.
var _ = http.MethodPost
