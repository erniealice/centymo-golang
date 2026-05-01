package list

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	scpspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule"
)

// ListViewDeps holds dependencies for the SCPS list view.
type ListViewDeps struct {
	Routes                             centymo.SupplierContractPriceScheduleRoutes
	Labels                             centymo.SupplierContractPriceScheduleLabels
	CommonLabels                       pyeza.CommonLabels
	TableLabels                        types.TableLabels
	ListSupplierContractPriceSchedules func(ctx context.Context, req *scpspb.ListSupplierContractPriceSchedulesRequest) (*scpspb.ListSupplierContractPriceSchedulesResponse, error)
	ListSupplierContracts              func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
}

// PageData holds the data for the SCPS list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the SCPS list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		resp, err := deps.ListSupplierContractPriceSchedules(ctx, &scpspb.ListSupplierContractPriceSchedulesRequest{})
		if err != nil {
			log.Printf("Failed to list price schedules: %v", err)
			return view.Error(fmt.Errorf("failed to load price schedules: %w", err))
		}

		schedules := resp.GetData()
		if status != "all" {
			filtered := make([]*scpspb.SupplierContractPriceSchedule, 0, len(schedules))
			for _, s := range schedules {
				if matchesStatus(s.GetStatus(), status) {
					filtered = append(filtered, s)
				}
			}
			schedules = filtered
		}

		// Load contract names for the Contract column.
		contractNames := loadContractNames(ctx, deps)

		l := deps.Labels
		columns := scheduleColumns(l)
		rows := buildTableRows(schedules, contractNames, l)
		types.ApplyColumnStyles(columns, rows)

		var primaryAction *types.PrimaryAction
		if deps.Routes.AddURL != "" {
			primaryAction = &types.PrimaryAction{
				Label:     l.Buttons.Add,
				ActionURL: deps.Routes.AddURL,
			}
		}

		tableConfig := &types.TableConfig{
			ID:                   "supplier-contract-price-schedules-table",
			RefreshURL:           deps.Routes.ListURL,
			Columns:              columns,
			Rows:                 rows,
			PrimaryAction:        primaryAction,
			ShowSearch:           true,
			ShowActions:          true,
			ShowFilters:          true,
			ShowSort:             true,
			ShowColumns:          true,
			ShowExport:           true,
			ShowDensity:          true,
			ShowEntries:          true,
			DefaultSortColumn:    "date_modified",
			DefaultSortDirection: "desc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   l.Empty.Title,
				Message: l.Empty.Message,
			},
		}
		types.ApplyTableSettings(tableConfig)

		heading := statusPageTitle(l, status)
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   status,
				HeaderTitle:    heading,
				HeaderSubtitle: l.Page.Caption,
				HeaderIcon:     "icon-calendar",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "supplier-contract-price-schedule-list-content",
			Table:           tableConfig,
		}

		return view.OK("supplier-contract-price-schedule-list", pageData)
	})
}

func matchesStatus(actual scpspb.SupplierContractPriceScheduleStatus, want string) bool {
	if want == "" || want == "all" {
		return true
	}
	return actual.String() == want
}

func loadContractNames(ctx context.Context, deps *ListViewDeps) map[string]string {
	out := map[string]string{}
	if deps.ListSupplierContracts == nil {
		return out
	}
	resp, err := deps.ListSupplierContracts(ctx, &suppliercontractpb.ListSupplierContractsRequest{})
	if err != nil {
		return out
	}
	for _, c := range resp.GetData() {
		out[c.GetId()] = c.GetName()
	}
	return out
}

func scheduleColumns(l centymo.SupplierContractPriceScheduleLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "supplier_contract", Label: l.Columns.SupplierContract, WidthClass: "col-3xl"},
		{Key: "sequence_number", Label: l.Columns.SequenceNumber, WidthClass: "col-xs", Align: "right"},
		{Key: "date_start", Label: l.Columns.DateStart, WidthClass: "col-2xl"},
		{Key: "date_end", Label: l.Columns.DateEnd, WidthClass: "col-2xl"},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
		{Key: "currency", Label: l.Columns.Currency, WidthClass: "col-xs"},
	}
}

func buildTableRows(schedules []*scpspb.SupplierContractPriceSchedule, contractNames map[string]string, l centymo.SupplierContractPriceScheduleLabels) []types.TableRow {
	rows := make([]types.TableRow, 0, len(schedules))
	for _, s := range schedules {
		id := s.GetId()
		statusStr := s.GetStatus().String()
		startStr := ""
		if t := s.GetDateTimeStart(); t != nil && t.IsValid() {
			startStr = t.AsTime().UTC().Format("2006-01-02")
		}
		endStr := "—"
		if t := s.GetDateTimeEnd(); t != nil && t.IsValid() {
			endStr = t.AsTime().UTC().Format("2006-01-02")
		}
		contractName := contractNames[s.GetSupplierContractId()]
		if contractName == "" {
			contractName = s.GetSupplierContractId()
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: s.GetName()},
				{Type: "text", Value: contractName},
				{Type: "number", Value: fmt.Sprintf("%d", s.GetSequenceNumber())},
				{Type: "text", Value: startStr},
				{Type: "text", Value: endStr},
				{Type: "badge", Value: scheduleStatusLabel(statusStr, l), Variant: scheduleStatusVariant(statusStr)},
				{Type: "text", Value: s.GetCurrency()},
			},
			DataAttrs: map[string]string{
				"name":              s.GetName(),
				"status":            statusStr,
				"supplier_contract": s.GetSupplierContractId(),
			},
		})
	}
	return rows
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

func statusPageTitle(l centymo.SupplierContractPriceScheduleLabels, status string) string {
	switch status {
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SCHEDULED":
		return l.Page.HeadingScheduled
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_ACTIVE":
		return l.Page.HeadingActive
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_SUPERSEDED":
		return l.Page.HeadingSuperseded
	case "SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_CANCELLED":
		return l.Page.HeadingCancelled
	default:
		return l.Page.Heading
	}
}
