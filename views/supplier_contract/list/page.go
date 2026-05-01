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
)

// ListViewDeps holds dependencies for the supplier contract list view.
type ListViewDeps struct {
	Routes                centymo.SupplierContractRoutes
	ListSupplierContracts func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
	Labels                centymo.SupplierContractLabels
	CommonLabels          pyeza.CommonLabels
	TableLabels           types.TableLabels
}

// PageData holds the data for the supplier contract list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the supplier contract list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		resp, err := deps.ListSupplierContracts(ctx, &suppliercontractpb.ListSupplierContractsRequest{})
		if err != nil {
			log.Printf("Failed to list supplier contracts: %v", err)
			return view.Error(fmt.Errorf("failed to load supplier contracts: %w", err))
		}

		contracts := resp.GetData()
		if status != "all" {
			var filtered []*suppliercontractpb.SupplierContract
			for _, c := range contracts {
				if c.GetStatus().String() == status {
					filtered = append(filtered, c)
				}
			}
			contracts = filtered
		}

		l := deps.Labels
		columns := supplierContractColumns(l)
		rows := buildTableRows(contracts, l)
		types.ApplyColumnStyles(columns, rows)

		var primaryAction *types.PrimaryAction
		if deps.Routes.AddURL != "" {
			primaryAction = &types.PrimaryAction{
				Label:     l.Page.AddButton,
				ActionURL: deps.Routes.AddURL,
			}
		}

		tableConfig := &types.TableConfig{
			ID:                   "supplier-contracts-table",
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
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "supplier-contract-list-content",
			Table:           tableConfig,
		}

		return view.OK("supplier-contract-list", pageData)
	})
}

func supplierContractColumns(l centymo.SupplierContractLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name},
		{Key: "supplier", Label: l.Columns.Supplier},
		{Key: "kind", Label: l.Columns.Kind, WidthClass: "col-2xl"},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
		{Key: "validity", Label: l.Columns.Validity, NoSort: true, WidthClass: "col-3xl"},
		{Key: "committed", Label: l.Columns.Committed, WidthClass: "col-3xl", Align: "right"},
		{Key: "released", Label: l.Columns.Released, WidthClass: "col-3xl", Align: "right"},
		{Key: "billed", Label: l.Columns.Billed, WidthClass: "col-3xl", Align: "right"},
		{Key: "remaining", Label: l.Columns.Remaining, WidthClass: "col-3xl", Align: "right"},
	}
}

func buildTableRows(contracts []*suppliercontractpb.SupplierContract, l centymo.SupplierContractLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, c := range contracts {
		id := c.GetId()
		name := c.GetName()
		supplierName := ""
		if s := c.GetSupplier(); s != nil {
			supplierName = s.GetName()
		}
		kindStr := c.GetKind().String()
		statusStr := c.GetStatus().String()
		currency := c.GetCurrency()
		startDate := c.GetDateTimeStart()
		endDate := c.GetDateTimeEnd()
		validity := startDate
		if endDate != "" {
			validity = startDate + " → " + endDate
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: name},
				{Type: "text", Value: supplierName},
				{Type: "badge", Value: kindStr, Variant: "default"},
				{Type: "badge", Value: statusStr, Variant: contractStatusVariant(statusStr)},
				{Type: "text", Value: validity},
				types.MoneyCell(float64(c.GetCommittedAmount()), currency, true),
				types.MoneyCell(float64(c.GetReleasedAmount()), currency, true),
				types.MoneyCell(float64(c.GetBilledAmount()), currency, true),
				types.MoneyCell(float64(c.GetRemainingAmount()), currency, true),
			},
			DataAttrs: map[string]string{
				"name":     name,
				"supplier": supplierName,
				"kind":     kindStr,
				"status":   statusStr,
			},
		})
	}
	return rows
}

func optionalStringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func contractStatusVariant(status string) string {
	switch status {
	case "SUPPLIER_CONTRACT_STATUS_DRAFT":
		return "default"
	case "SUPPLIER_CONTRACT_STATUS_REQUESTED":
		return "warning"
	case "SUPPLIER_CONTRACT_STATUS_PENDING_APPROVAL":
		return "warning"
	case "SUPPLIER_CONTRACT_STATUS_APPROVED":
		return "info"
	case "SUPPLIER_CONTRACT_STATUS_ACTIVE":
		return "success"
	case "SUPPLIER_CONTRACT_STATUS_EXPIRING":
		return "warning"
	case "SUPPLIER_CONTRACT_STATUS_SUSPENDED":
		return "warning"
	case "SUPPLIER_CONTRACT_STATUS_EXPIRED":
		return "danger"
	case "SUPPLIER_CONTRACT_STATUS_TERMINATED":
		return "danger"
	case "SUPPLIER_CONTRACT_STATUS_REJECTED":
		return "danger"
	default:
		return "default"
	}
}

func statusPageTitle(l centymo.SupplierContractLabels, status string) string {
	switch status {
	case "SUPPLIER_CONTRACT_STATUS_DRAFT":
		return l.Page.HeadingDraft
	case "SUPPLIER_CONTRACT_STATUS_ACTIVE":
		return l.Page.HeadingActive
	case "SUPPLIER_CONTRACT_STATUS_EXPIRING":
		return l.Page.HeadingExpiring
	case "SUPPLIER_CONTRACT_STATUS_EXPIRED":
		return l.Page.HeadingExpired
	case "SUPPLIER_CONTRACT_STATUS_TERMINATED":
		return l.Page.HeadingTerminated
	default:
		return l.Page.Heading
	}
}
