package list

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	procurementrequestpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request"
)

// ListViewDeps holds dependencies for the procurement request list view.
type ListViewDeps struct {
	Routes                  centymo.ProcurementRequestRoutes
	ListProcurementRequests func(ctx context.Context, req *procurementrequestpb.ListProcurementRequestsRequest) (*procurementrequestpb.ListProcurementRequestsResponse, error)
	Labels                  centymo.ProcurementRequestLabels
	CommonLabels            pyeza.CommonLabels
	TableLabels             types.TableLabels
}

// PageData holds the data for the procurement request list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the procurement request list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "draft"
		}

		resp, err := deps.ListProcurementRequests(ctx, &procurementrequestpb.ListProcurementRequestsRequest{})
		if err != nil {
			log.Printf("Failed to list procurement requests: %v", err)
			return view.Error(fmt.Errorf("failed to load procurement requests: %w", err))
		}

		requests := resp.GetData()
		if status != "all" {
			var filtered []*procurementrequestpb.ProcurementRequest
			for _, r := range requests {
				if r.GetStatus().String() == status {
					filtered = append(filtered, r)
				}
			}
			requests = filtered
		}

		l := deps.Labels
		columns := procurementRequestColumns(l)
		rows := buildTableRows(requests, l)
		types.ApplyColumnStyles(columns, rows)

		var primaryAction *types.PrimaryAction
		if deps.Routes.AddURL != "" {
			primaryAction = &types.PrimaryAction{
				Label:     l.Page.AddButton,
				ActionURL: deps.Routes.AddURL,
			}
		}

		tableConfig := &types.TableConfig{
			ID:                   "procurement-requests-table",
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
			DefaultSortColumn:    "date_created",
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
				HeaderIcon:     "icon-clipboard",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "procurement-request-list-content",
			Table:           tableConfig,
		}

		return view.OK("procurement-request-list", pageData)
	})
}

func procurementRequestColumns(l centymo.ProcurementRequestLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "request_number", Label: l.Columns.RequestNumber, Sortable: true},
		{Key: "status", Label: l.Columns.Status, Sortable: true, WidthClass: "col-2xl"},
		{Key: "requester", Label: l.Columns.Requester, Sortable: true},
		{Key: "supplier", Label: l.Columns.Supplier, Sortable: true},
		{Key: "estimated_total", Label: l.Columns.EstimatedTotal, Sortable: true, WidthClass: "col-3xl", Align: "right"},
		{Key: "needed_by", Label: l.Columns.NeededBy, Sortable: true, WidthClass: "col-3xl"},
		{Key: "date_created", Label: l.Columns.DateCreated, Sortable: true, WidthClass: "col-3xl"},
	}
}

func buildTableRows(requests []*procurementrequestpb.ProcurementRequest, l centymo.ProcurementRequestLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, req := range requests {
		id := req.GetId()
		requestNumber := req.GetRequestNumber()
		statusStr := req.GetStatus().String()
		currency := req.GetCurrency()

		supplierName := ""
		if s := req.GetSupplier(); s != nil {
			supplierName = s.GetName()
		}

		requesterName := req.GetRequesterUserId() // P4 will resolve to a display name

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: requestNumber},
				{Type: "badge", Value: statusStr, Variant: requestStatusVariant(statusStr)},
				{Type: "text", Value: requesterName},
				{Type: "text", Value: supplierName},
				types.MoneyCell(float64(req.GetEstimatedTotalAmount()), currency, true),
				types.DateTimeCell(req.GetNeededByDate(), types.DateReadable),
				types.DateTimeCell(req.GetDateCreatedString(), types.DateReadable),
			},
			DataAttrs: map[string]string{
				"request_number": requestNumber,
				"status":         statusStr,
				"requester":      requesterName,
				"supplier":       supplierName,
			},
		})
	}
	return rows
}

func optionalStrVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func requestStatusVariant(status string) string {
	switch status {
	case "PROCUREMENT_REQUEST_STATUS_DRAFT":
		return "default"
	case "PROCUREMENT_REQUEST_STATUS_SUBMITTED":
		return "info"
	case "PROCUREMENT_REQUEST_STATUS_PENDING_APPROVAL":
		return "warning"
	case "PROCUREMENT_REQUEST_STATUS_APPROVED", "PROCUREMENT_REQUEST_STATUS_FULFILLED":
		return "success"
	case "PROCUREMENT_REQUEST_STATUS_REJECTED", "PROCUREMENT_REQUEST_STATUS_CANCELLED":
		return "danger"
	default:
		return "default"
	}
}

func statusPageTitle(l centymo.ProcurementRequestLabels, status string) string {
	switch status {
	case "PROCUREMENT_REQUEST_STATUS_DRAFT":
		return l.Page.HeadingDraft
	case "PROCUREMENT_REQUEST_STATUS_SUBMITTED":
		return l.Page.HeadingSubmitted
	case "PROCUREMENT_REQUEST_STATUS_PENDING_APPROVAL":
		return l.Page.HeadingPendingApproval
	case "PROCUREMENT_REQUEST_STATUS_APPROVED":
		return l.Page.HeadingApproved
	case "PROCUREMENT_REQUEST_STATUS_REJECTED":
		return l.Page.HeadingRejected
	case "PROCUREMENT_REQUEST_STATUS_FULFILLED":
		return l.Page.HeadingFulfilled
	case "PROCUREMENT_REQUEST_STATUS_CANCELLED":
		return l.Page.HeadingCancelled
	default:
		return l.Page.Heading
	}
}
