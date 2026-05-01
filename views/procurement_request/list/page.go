package list

import (
	"context"
	"fmt"
	"log"
	"strings"

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

	// SPS Wave 3 — F3 fulfillment_strategy filter chips. Rendered above the
	// table; current selection echoed back as `ActiveStrategy`.
	StrategyChips   []StrategyChip
	ActiveStrategy  string
	StrategyLabel   string
}

// StrategyChip describes a single F3 fulfillment-strategy filter chip.
type StrategyChip struct {
	Value  string // "any" | UNIFORM_OUTRIGHT | UNIFORM_STOCKABLE | … | MIXED
	Label  string
	Href   string // list URL pre-bound with the right ?strategy= param
	Active bool
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

		// SPS Wave 3 — F3 fulfillment_strategy filter (additional filter on top
		// of the status path-segment). "any" or empty = no filter.
		strategy := viewCtx.Request.URL.Query().Get("strategy")
		if strategy != "" && strategy != "any" {
			var filtered []*procurementrequestpb.ProcurementRequest
			for _, r := range requests {
				if r.GetFulfillmentStrategy().String() == strategy {
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
			StrategyChips:   buildStrategyChips(deps, status, strategy),
			ActiveStrategy:  strategy,
			StrategyLabel:   l.Filters.FulfillmentStrategy,
		}

		return view.OK("procurement-request-list", pageData)
	})
}

func procurementRequestColumns(l centymo.ProcurementRequestLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "request_number", Label: l.Columns.RequestNumber},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
		// SPS Wave 3 — F3 strategy column (header-level rollup).
		{Key: "strategy", Label: l.Filters.FulfillmentStrategy, WidthClass: "col-2xl"},
		{Key: "requester", Label: l.Columns.Requester},
		{Key: "supplier", Label: l.Columns.Supplier},
		{Key: "estimated_total", Label: l.Columns.EstimatedTotal, WidthClass: "col-3xl", Align: "right"},
		{Key: "needed_by", Label: l.Columns.NeededBy, WidthClass: "col-3xl"},
		{Key: "date_created", Label: l.Columns.DateCreated, WidthClass: "col-3xl"},
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

		strategyStr := req.GetFulfillmentStrategy().String()
		strategyLabel := strategyDisplayLabel(l, strategyStr)

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: requestNumber},
				{Type: "badge", Value: statusBadgeLabel(l, statusStr), Variant: requestStatusVariant(statusStr)},
				{Type: "badge", Value: strategyLabel, Variant: "default"},
				{Type: "text", Value: requesterName},
				{Type: "text", Value: supplierName},
				types.MoneyCell(float64(req.GetEstimatedTotalAmount()), currency, true),
				types.DateTimeCell(req.GetNeededByDate(), types.DateReadable),
				types.DateTimeCell(req.GetDateCreatedString(), types.DateReadable),
			},
			DataAttrs: map[string]string{
				"request_number": requestNumber,
				"status":         statusStr,
				"strategy":       strategyStr,
				"requester":      requesterName,
				"supplier":       supplierName,
			},
		})
	}
	return rows
}

// statusBadgeLabel maps the proto status enum string to the localized label.
func statusBadgeLabel(l centymo.ProcurementRequestLabels, status string) string {
	switch status {
	case "PROCUREMENT_REQUEST_STATUS_DRAFT":
		return l.Form.StatusDraft
	case "PROCUREMENT_REQUEST_STATUS_SUBMITTED":
		return l.Form.StatusSubmitted
	case "PROCUREMENT_REQUEST_STATUS_PENDING_APPROVAL":
		return l.Form.StatusPendingApproval
	case "PROCUREMENT_REQUEST_STATUS_APPROVED":
		return l.Form.StatusApproved
	case "PROCUREMENT_REQUEST_STATUS_APPROVED_PENDING_SPAWN":
		return l.Form.StatusApprovedPendingSpawn
	case "PROCUREMENT_REQUEST_STATUS_REJECTED":
		return l.Form.StatusRejected
	case "PROCUREMENT_REQUEST_STATUS_FULFILLED":
		return l.Form.StatusFulfilled
	case "PROCUREMENT_REQUEST_STATUS_CANCELLED":
		return l.Form.StatusCancelled
	}
	return status
}

// strategyDisplayLabel maps the proto strategy enum string to the localized label.
func strategyDisplayLabel(l centymo.ProcurementRequestLabels, strategy string) string {
	switch strategy {
	case "PROCUREMENT_REQUEST_FULFILLMENT_STRATEGY_UNIFORM_OUTRIGHT":
		return l.FulfillmentStrategy.UniformOutright
	case "PROCUREMENT_REQUEST_FULFILLMENT_STRATEGY_UNIFORM_STOCKABLE":
		return l.FulfillmentStrategy.UniformStockable
	case "PROCUREMENT_REQUEST_FULFILLMENT_STRATEGY_UNIFORM_RECURRING":
		return l.FulfillmentStrategy.UniformRecurring
	case "PROCUREMENT_REQUEST_FULFILLMENT_STRATEGY_UNIFORM_PETTY":
		return l.FulfillmentStrategy.UniformPetty
	case "PROCUREMENT_REQUEST_FULFILLMENT_STRATEGY_MIXED":
		return l.FulfillmentStrategy.Mixed
	}
	return "—"
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
	case "PROCUREMENT_REQUEST_STATUS_PENDING_APPROVAL",
		"PROCUREMENT_REQUEST_STATUS_APPROVED_PENDING_SPAWN":
		// CRIT-3: APPROVED_PENDING_SPAWN is in-progress (saga still running or
		// at least one line FAILED awaiting retry). Reuse the warning palette.
		return "warning"
	case "PROCUREMENT_REQUEST_STATUS_APPROVED", "PROCUREMENT_REQUEST_STATUS_FULFILLED":
		return "success"
	case "PROCUREMENT_REQUEST_STATUS_REJECTED", "PROCUREMENT_REQUEST_STATUS_CANCELLED":
		return "danger"
	default:
		return "default"
	}
}

// buildStrategyChips builds the F3 fulfillment-strategy chip set. Order:
// "any" first, then 4 uniform values, then MIXED. Hrefs preserve the current
// status path segment so the chip click only changes the strategy filter.
func buildStrategyChips(deps *ListViewDeps, status, active string) []StrategyChip {
	if active == "" {
		active = "any"
	}
	l := deps.Labels
	base := strings.Replace(deps.Routes.ListURL, "{status}", status, 1)
	withQuery := func(value string) string {
		if value == "any" {
			return base
		}
		return base + "?strategy=" + value
	}
	chips := []StrategyChip{
		{Value: "any", Label: l.Filters.AnyFulfillmentStrategy},
		{Value: "PROCUREMENT_REQUEST_FULFILLMENT_STRATEGY_UNIFORM_OUTRIGHT", Label: l.FulfillmentStrategy.UniformOutright},
		{Value: "PROCUREMENT_REQUEST_FULFILLMENT_STRATEGY_UNIFORM_STOCKABLE", Label: l.FulfillmentStrategy.UniformStockable},
		{Value: "PROCUREMENT_REQUEST_FULFILLMENT_STRATEGY_UNIFORM_RECURRING", Label: l.FulfillmentStrategy.UniformRecurring},
		{Value: "PROCUREMENT_REQUEST_FULFILLMENT_STRATEGY_UNIFORM_PETTY", Label: l.FulfillmentStrategy.UniformPetty},
		{Value: "PROCUREMENT_REQUEST_FULFILLMENT_STRATEGY_MIXED", Label: l.FulfillmentStrategy.Mixed},
	}
	for i := range chips {
		chips[i].Href = withQuery(chips[i].Value)
		chips[i].Active = chips[i].Value == active
	}
	return chips
}

func statusPageTitle(l centymo.ProcurementRequestLabels, status string) string {
	switch status {
	case "PROCUREMENT_REQUEST_STATUS_DRAFT":
		return l.Page.HeadingDraft
	case "PROCUREMENT_REQUEST_STATUS_SUBMITTED":
		return l.Page.HeadingSubmitted
	case "PROCUREMENT_REQUEST_STATUS_PENDING_APPROVAL":
		return l.Page.HeadingPendingApproval
	case "PROCUREMENT_REQUEST_STATUS_APPROVED",
		"PROCUREMENT_REQUEST_STATUS_APPROVED_PENDING_SPAWN":
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
