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

	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	procurementrequestpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request"
	procurementrequestlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request_line"
)

// DetailViewDeps holds dependencies for the procurement request detail page.
type DetailViewDeps struct {
	Routes       centymo.ProcurementRequestRoutes
	Labels       centymo.ProcurementRequestLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadProcurementRequest      func(ctx context.Context, req *procurementrequestpb.ReadProcurementRequestRequest) (*procurementrequestpb.ReadProcurementRequestResponse, error)
	ListProcurementRequestLines func(ctx context.Context, req *procurementrequestlinepb.ListProcurementRequestLinesRequest) (*procurementrequestlinepb.ListProcurementRequestLinesResponse, error)

	// Spawned POs — optional; list POs where procurement_request_id = this request
	ListPurchaseOrders func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)
}

// PageData holds template data for the procurement request detail page.
type PageData struct {
	types.PageData
	ContentTemplate string

	// Request record
	Request       map[string]any
	StatusVariant string

	// Tab navigation
	TabItems  []pyeza.TabItem
	ActiveTab string

	// Lines tab
	LineItemTable  *types.TableConfig
	LineItemAddURL string

	// Spawned POs tab
	SpawnedPOTable *types.TableConfig

	// Action URLs
	SubmitURL  string
	ApproveURL string
	RejectURL  string
	SpawnPOURL string
	EditURL    string
}

const (
	tabInfo      = "info"
	tabLines     = "lines"
	tabSpawnedPO = "spawned-pos"
	tabActivity  = "activity"
)

// NewView creates the procurement request detail page view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Redirect(deps.Routes.ListURL)
		}

		resp, err := deps.ReadProcurementRequest(ctx, &procurementrequestpb.ReadProcurementRequestRequest{
			Data: &procurementrequestpb.ProcurementRequest{Id: id},
		})
		if err != nil {
			log.Printf("ReadProcurementRequest %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load procurement request: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("procurement request not found"))
		}
		req := data[0]

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = tabInfo
		}

		l := deps.Labels

		supplierName := ""
		if s := req.GetSupplier(); s != nil {
			supplierName = s.GetName()
		}

		requestMap := map[string]any{
			"id":                     req.GetId(),
			"request_number":         req.GetRequestNumber(),
			"status":                 req.GetStatus().String(),
			"requester_user_id":      req.GetRequesterUserId(),
			"supplier_name":          supplierName,
			"currency":               req.GetCurrency(),
			"estimated_total_amount": req.GetEstimatedTotalAmount(),
			"needed_by_date":         req.GetNeededByDate(),
			"justification":          req.GetJustification(),
			"approved_by":            req.GetApprovedBy(),
			"date_created_string":    req.GetDateCreatedString(),
		}

		tabItems := []pyeza.TabItem{
			{Key: tabInfo, Label: l.Tabs.Info},
			{Key: tabLines, Label: l.Tabs.Lines},
			{Key: tabSpawnedPO, Label: l.Tabs.SpawnedPOs},
			{Key: tabActivity, Label: l.Tabs.Activity},
		}

		pd := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          req.GetRequestNumber(),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    req.GetRequestNumber(),
				HeaderSubtitle: l.Page.DetailSubtitle,
				HeaderIcon:     "icon-clipboard",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "procurement-request-detail-content",
			Request:         requestMap,
			StatusVariant:   requestStatusVariant(req.GetStatus().String()),
			TabItems:        tabItems,
			ActiveTab:       activeTab,
			SubmitURL:       buildURL(deps.Routes.SubmitURL, id),
			ApproveURL:      buildURL(deps.Routes.ApproveURL, id),
			RejectURL:       buildURL(deps.Routes.RejectURL, id),
			SpawnPOURL:      buildURL(deps.Routes.SpawnPOURL, id),
			EditURL:         buildURL(deps.Routes.EditURL, id),
		}

		if activeTab == tabLines && deps.ListProcurementRequestLines != nil {
			pd.LineItemTable = buildLineTable(ctx, deps, id, l)
			pd.LineItemAddURL = buildURL(deps.Routes.LineAddURL, id)
		}

		if activeTab == tabSpawnedPO && deps.ListPurchaseOrders != nil {
			pd.SpawnedPOTable = buildSpawnedPOTable(ctx, deps, id, l)
		}

		return view.OK("procurement-request-detail", pd)
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

		resp, err := deps.ReadProcurementRequest(ctx, &procurementrequestpb.ReadProcurementRequestRequest{
			Data: &procurementrequestpb.ProcurementRequest{Id: id},
		})
		if err != nil {
			return view.Error(fmt.Errorf("failed to load procurement request: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("procurement request not found"))
		}
		req := data[0]

		l := deps.Labels
		supplierName := ""
		if s := req.GetSupplier(); s != nil {
			supplierName = s.GetName()
		}

		requestMap := map[string]any{
			"id":                     req.GetId(),
			"request_number":         req.GetRequestNumber(),
			"status":                 req.GetStatus().String(),
			"requester_user_id":      req.GetRequesterUserId(),
			"supplier_name":          supplierName,
			"currency":               req.GetCurrency(),
			"estimated_total_amount": req.GetEstimatedTotalAmount(),
			"needed_by_date":         req.GetNeededByDate(),
			"justification":          req.GetJustification(),
			"approved_by":            req.GetApprovedBy(),
			"date_created_string":    req.GetDateCreatedString(),
		}

		tabItems := []pyeza.TabItem{
			{Key: tabInfo, Label: l.Tabs.Info},
			{Key: tabLines, Label: l.Tabs.Lines},
			{Key: tabSpawnedPO, Label: l.Tabs.SpawnedPOs},
			{Key: tabActivity, Label: l.Tabs.Activity},
		}

		pd := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "procurement-request-detail-content",
			Request:         requestMap,
			StatusVariant:   requestStatusVariant(req.GetStatus().String()),
			TabItems:        tabItems,
			ActiveTab:       tab,
			SubmitURL:       buildURL(deps.Routes.SubmitURL, id),
			ApproveURL:      buildURL(deps.Routes.ApproveURL, id),
			RejectURL:       buildURL(deps.Routes.RejectURL, id),
			SpawnPOURL:      buildURL(deps.Routes.SpawnPOURL, id),
			EditURL:         buildURL(deps.Routes.EditURL, id),
		}

		switch tab {
		case tabLines:
			if deps.ListProcurementRequestLines != nil {
				pd.LineItemTable = buildLineTable(ctx, deps, id, l)
				pd.LineItemAddURL = buildURL(deps.Routes.LineAddURL, id)
			}
		case tabSpawnedPO:
			if deps.ListPurchaseOrders != nil {
				pd.SpawnedPOTable = buildSpawnedPOTable(ctx, deps, id, l)
			}
		}

		return view.OK("procurement-request-tab-content", pd)
	})
}

// NewSubmitAction handles POST /action/procurement-request/submit/{id}.
func NewSubmitAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		detailURL := buildURL(deps.Routes.DetailURL, id)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Redirect": detailURL,
			},
		}
	})
}

// NewApproveAction handles POST /action/procurement-request/approve/{id}.
func NewApproveAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		detailURL := buildURL(deps.Routes.DetailURL, id)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Redirect": detailURL,
			},
		}
	})
}

// NewRejectAction handles POST /action/procurement-request/reject/{id}.
func NewRejectAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		detailURL := buildURL(deps.Routes.DetailURL, id)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Redirect": detailURL,
			},
		}
	})
}

// NewSpawnPOAction handles POST /action/procurement-request/spawn-po/{id}.
func NewSpawnPOAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		// P2 will implement the use case; redirect back to detail for now.
		detailURL := buildURL(deps.Routes.DetailURL, id)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Redirect": detailURL,
			},
		}
	})
}

// --- helpers -----------------------------------------------------------------

func buildURL(template, id string) string {
	if template == "" {
		return ""
	}
	return route.ResolveURL(template, "id", id)
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

func buildLineTable(ctx context.Context, deps *DetailViewDeps, requestID string, l centymo.ProcurementRequestLabels) *types.TableConfig {
	prIDPtr := requestID
	resp, err := deps.ListProcurementRequestLines(ctx, &procurementrequestlinepb.ListProcurementRequestLinesRequest{
		ProcurementRequestId: &prIDPtr,
	})
	if err != nil {
		log.Printf("ListProcurementRequestLines for %s: %v", requestID, err)
		return nil
	}

	columns := []types.TableColumn{
		{Key: "description", Label: l.Lines.Description},
		{Key: "line_type", Label: l.Lines.LineType, WidthClass: "col-2xl"},
		{Key: "quantity", Label: l.Lines.Quantity, Align: "right", WidthClass: "col-xl"},
		{Key: "estimated_unit_price", Label: l.Lines.EstimatedUnitPrice, Align: "right", WidthClass: "col-3xl"},
		{Key: "estimated_total_price", Label: l.Lines.EstimatedTotalPrice, Align: "right", WidthClass: "col-3xl"},
	}

	rows := []types.TableRow{}
	for _, line := range resp.GetData() {
		rows = append(rows, types.TableRow{
			ID: line.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: line.GetDescription()},
				{Type: "badge", Value: line.GetLineType(), Variant: "default"},
				{Type: "number", Value: fmt.Sprintf("%.2f", line.GetQuantity())},
				types.MoneyCell(float64(line.GetEstimatedUnitPrice()), "", true),
				types.MoneyCell(float64(line.GetEstimatedTotalPrice()), "", true),
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:          "pr-lines-table",
		Columns:     columns,
		Rows:        rows,
		ShowActions: true,
		EmptyState: types.TableEmptyState{
			Title:   l.Lines.EmptyTitle,
			Message: l.Lines.EmptyMessage,
		},
	}
}

func buildSpawnedPOTable(ctx context.Context, deps *DetailViewDeps, requestID string, l centymo.ProcurementRequestLabels) *types.TableConfig {
	resp, err := deps.ListPurchaseOrders(ctx, &purchaseorderpb.ListPurchaseOrdersRequest{})
	if err != nil {
		log.Printf("ListPurchaseOrders for request %s: %v", requestID, err)
		return nil
	}

	// Filter to POs spawned from this procurement request
	var linked []*purchaseorderpb.PurchaseOrder
	for _, po := range resp.GetData() {
		if po.GetProcurementRequestId() == requestID {
			linked = append(linked, po)
		}
	}

	columns := []types.TableColumn{
		{Key: "po_number", Label: l.SpawnedPOs.PONumber},
		{Key: "status", Label: l.SpawnedPOs.Status, WidthClass: "col-2xl"},
		{Key: "total_amount", Label: l.SpawnedPOs.TotalAmount, Align: "right", WidthClass: "col-3xl"},
		{Key: "order_date", Label: l.SpawnedPOs.OrderDate, WidthClass: "col-3xl"},
	}

	rows := []types.TableRow{}
	for _, po := range linked {
		rows = append(rows, types.TableRow{
			ID: po.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: po.GetPoNumber()},
				{Type: "badge", Value: po.GetStatus(), Variant: "default"},
				types.MoneyCell(float64(po.GetTotalAmount()), po.GetCurrency(), true),
				types.DateTimeCell(po.GetOrderDateString(), types.DateReadable),
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "spawned-pos-table",
		Columns: columns,
		Rows:    rows,
		EmptyState: types.TableEmptyState{
			Title:   l.SpawnedPOs.EmptyTitle,
			Message: l.SpawnedPOs.EmptyMessage,
		},
	}
}
