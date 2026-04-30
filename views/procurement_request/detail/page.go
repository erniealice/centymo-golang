package detail

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

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

	// Workflow invocations — block.go injects closures that call the espyna use
	// cases (with user-id sourced from request context). All optional — when
	// nil the action handler degrades to a redirect-only no-op.
	SubmitProcurementRequest  func(ctx context.Context, id string) error
	ApproveProcurementRequest func(ctx context.Context, id string) error
	RejectProcurementRequest  func(ctx context.Context, id string, reason string) error
	SpawnPurchaseOrder        func(ctx context.Context, id string) (newPOID string, err error)
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
			"status_label":           statusDisplayLabel(l, req.GetStatus().String()),
			"requester_user_id":      req.GetRequesterUserId(),
			"supplier_name":          supplierName,
			"currency":               req.GetCurrency(),
			"estimated_total_amount": req.GetEstimatedTotalAmount(),
			"needed_by_date":         req.GetNeededByDate(),
			"justification":          req.GetJustification(),
			"approved_by":            req.GetApprovedBy(),
			"date_created_string":    req.GetDateCreatedString(),
			// SPS Wave 3 — F3 strategy + policy_decision_log surfaces
			"fulfillment_strategy":       req.GetFulfillmentStrategy().String(),
			"fulfillment_strategy_label": strategyDisplayLabel(l, req.GetFulfillmentStrategy().String()),
			"policy_decision_log":        req.GetPolicyDecisionLog(),
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
			"status_label":           statusDisplayLabel(l, req.GetStatus().String()),
			"requester_user_id":      req.GetRequesterUserId(),
			"supplier_name":          supplierName,
			"currency":               req.GetCurrency(),
			"estimated_total_amount": req.GetEstimatedTotalAmount(),
			"needed_by_date":         req.GetNeededByDate(),
			"justification":          req.GetJustification(),
			"approved_by":            req.GetApprovedBy(),
			"date_created_string":    req.GetDateCreatedString(),
			"fulfillment_strategy":       req.GetFulfillmentStrategy().String(),
			"fulfillment_strategy_label": strategyDisplayLabel(l, req.GetFulfillmentStrategy().String()),
			"policy_decision_log":        req.GetPolicyDecisionLog(),
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
// Invokes SubmitProcurementRequest use case (DRAFT → SUBMITTED) before
// HX-Redirecting back to the detail page.
func NewSubmitAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		if deps.SubmitProcurementRequest != nil {
			if err := deps.SubmitProcurementRequest(ctx, id); err != nil {
				log.Printf("SubmitProcurementRequest %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
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
// Invokes ApproveProcurementRequest use case (PENDING_APPROVAL → APPROVED)
// before HX-Redirecting back to the detail page.
func NewApproveAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		if deps.ApproveProcurementRequest != nil {
			if err := deps.ApproveProcurementRequest(ctx, id); err != nil {
				log.Printf("ApproveProcurementRequest %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
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
// Invokes RejectProcurementRequest use case (PENDING_APPROVAL → REJECTED)
// before HX-Redirecting back to the detail page. Optional rejection_reason
// from form body is forwarded to the use case.
func NewRejectAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		reason := ""
		if err := viewCtx.Request.ParseForm(); err == nil {
			reason = viewCtx.Request.FormValue("rejection_reason")
		}
		if deps.RejectProcurementRequest != nil {
			if err := deps.RejectProcurementRequest(ctx, id, reason); err != nil {
				log.Printf("RejectProcurementRequest %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
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
// Invokes SpawnPurchaseOrder use case (creates a new PurchaseOrder linked to
// this APPROVED request) and redirects back to the detail page on the
// "spawned-pos" tab. The new PO ID is currently rendered via the linked tab,
// not pushed in the redirect.
func NewSpawnPOAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		if deps.SpawnPurchaseOrder != nil {
			if _, err := deps.SpawnPurchaseOrder(ctx, id); err != nil {
				log.Printf("SpawnPurchaseOrder %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
		}
		detailURL := buildURL(deps.Routes.DetailURL, id) + "?tab=spawned-pos"
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
	case "PROCUREMENT_REQUEST_STATUS_PENDING_APPROVAL",
		"PROCUREMENT_REQUEST_STATUS_APPROVED_PENDING_SPAWN":
		// CRIT-3: APPROVED_PENDING_SPAWN reuses the warning palette to convey
		// "in progress" semantics until every line spawns its downstream artifact.
		return "warning"
	case "PROCUREMENT_REQUEST_STATUS_APPROVED", "PROCUREMENT_REQUEST_STATUS_FULFILLED":
		return "success"
	case "PROCUREMENT_REQUEST_STATUS_REJECTED", "PROCUREMENT_REQUEST_STATUS_CANCELLED":
		return "danger"
	default:
		return "default"
	}
}

// statusDisplayLabel maps the proto status enum string to the localized label.
func statusDisplayLabel(l centymo.ProcurementRequestLabels, status string) string {
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
	return ""
}

// fulfillmentModeDisplay maps the proto enum string to a (label, badgeMode-token)
// pair for badge rendering. Token feeds CSS via data-mode={token}.
func fulfillmentModeDisplay(l centymo.ProcurementRequestLabels, mode string) (label, token string) {
	switch mode {
	case "PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_OUTRIGHT":
		return l.FulfillmentMode.Outright, "outright"
	case "PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_STOCKABLE":
		return l.FulfillmentMode.Stockable, "stockable"
	case "PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_RECURRING":
		return l.FulfillmentMode.Recurring, "recurring"
	case "PROCUREMENT_REQUEST_LINE_FULFILLMENT_MODE_PETTY":
		return l.FulfillmentMode.Petty, "petty"
	}
	return "—", "unspecified"
}

// spawnStatusDisplay maps the proto spawn-status enum string to a localized
// label + status-badge token for variant styling.
func spawnStatusDisplay(l centymo.ProcurementRequestLabels, status string) (label, badge string) {
	switch status {
	case "PROCUREMENT_REQUEST_LINE_SPAWN_STATUS_PENDING":
		return l.Spawn.StatusPending, "default"
	case "PROCUREMENT_REQUEST_LINE_SPAWN_STATUS_SPAWNING":
		return l.Spawn.StatusSpawning, "info"
	case "PROCUREMENT_REQUEST_LINE_SPAWN_STATUS_SPAWNED":
		return l.Spawn.StatusSpawned, "success"
	case "PROCUREMENT_REQUEST_LINE_SPAWN_STATUS_FAILED":
		return l.Spawn.StatusFailed, "danger"
	}
	return l.Spawn.StatusUnspecified, "default"
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

	// SPS Wave 3 — F1 mode badge, F2 spawn back-FK link, CRIT-3 spawn status +
	// retry CTA all live in the lines table.
	columns := []types.TableColumn{
		{Key: "description", Label: l.Lines.Description},
		{Key: "line_type", Label: l.Lines.LineType, WidthClass: "col-2xl"},
		{Key: "fulfillment_mode", Label: l.Spawn.ModeColumn, WidthClass: "col-2xl"},
		{Key: "quantity", Label: l.Lines.Quantity, Align: "right", WidthClass: "col-xl"},
		{Key: "estimated_unit_price", Label: l.Lines.EstimatedUnitPrice, Align: "right", WidthClass: "col-3xl"},
		{Key: "estimated_total_price", Label: l.Lines.EstimatedTotalPrice, Align: "right", WidthClass: "col-3xl"},
		{Key: "spawn_status", Label: l.Spawn.StatusColumn, WidthClass: "col-2xl"},
		{Key: "spawned_link", Label: l.Spawn.SpawnedColumn, WidthClass: "col-3xl"},
	}

	rows := []types.TableRow{}
	for _, line := range resp.GetData() {
		modeStr := line.GetFulfillmentMode().String()
		modeLabel, modeToken := fulfillmentModeDisplay(l, modeStr)

		spawnStr := line.GetSpawnStatus().String()
		spawnLabel, spawnVariant := spawnStatusDisplay(l, spawnStr)

		spawnedHTML := buildSpawnedCellHTML(line, l, requestID, deps.Routes)

		modeCell := types.TableCell{
			Type:    "html",
			HTML:    template.HTML(fmt.Sprintf(`<span class="badge" data-mode="%s" data-testid="prl-mode-badge">%s</span>`, htmlEsc(modeToken), htmlEsc(modeLabel))),
			Value:   modeLabel,
		}

		spawnCell := types.TableCell{
			Type:    "badge",
			Value:   spawnLabel,
			Variant: spawnVariant,
		}

		rows = append(rows, types.TableRow{
			ID: line.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: line.GetDescription()},
				{Type: "badge", Value: line.GetLineType(), Variant: "default"},
				modeCell,
				{Type: "number", Value: fmt.Sprintf("%.2f", line.GetQuantity())},
				types.MoneyCell(float64(line.GetEstimatedUnitPrice()), "", true),
				types.MoneyCell(float64(line.GetEstimatedTotalPrice()), "", true),
				spawnCell,
				spawnedHTML,
			},
			DataAttrs: map[string]string{
				"fulfillment_mode": modeStr,
				"spawn_status":     spawnStr,
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

// buildSpawnedCellHTML renders the F2 spawned-artifact link OR the CRIT-3
// failure note + retry button in a single HTML cell. The cell value depends
// on (spawn_status, fulfillment_mode):
//   - SPAWNED: link to the populated back-FK artifact (PO line / contract /
//     expenditure).
//   - FAILED: render `spawn_error` + a Retry button that POSTs to the
//     placeholder retry route (handler is a stub for SPS Wave 3).
//   - PENDING / SPAWNING / UNSPECIFIED: small em-dash placeholder.
func buildSpawnedCellHTML(line *procurementrequestlinepb.ProcurementRequestLine, l centymo.ProcurementRequestLabels, requestID string, routes centymo.ProcurementRequestRoutes) types.TableCell {
	spawnStatus := line.GetSpawnStatus().String()

	switch spawnStatus {
	case "PROCUREMENT_REQUEST_LINE_SPAWN_STATUS_SPAWNED":
		var (
			href, label string
		)
		if v := line.GetSpawnedPurchaseOrderLineItemId(); v != "" {
			href = "/app/purchase-orders/detail/" + v
			label = l.Spawn.LinkPO
		} else if v := line.GetSpawnedSupplierContractId(); v != "" {
			href = "/app/supplier-contracts/detail/" + v
			label = l.Spawn.LinkContract
		} else if v := line.GetSpawnedExpenditureId(); v != "" {
			href = "/app/expenditures/detail/" + v
			label = l.Spawn.LinkExpenditure
		}
		if href == "" {
			return types.TableCell{Type: "text", Value: l.Spawn.NotApplicable}
		}
		return types.TableCell{
			Type: "html",
			HTML: template.HTML(fmt.Sprintf(`<a href="%s" class="btn-link" data-testid="prl-spawned-link">%s</a>`, htmlEsc(href), htmlEsc(label))),
			Value: label,
		}

	case "PROCUREMENT_REQUEST_LINE_SPAWN_STATUS_FAILED":
		retryURL := strings.NewReplacer("{id}", requestID, "{lid}", line.GetId()).Replace(routes.LineRetrySpawnURL)
		errMsg := line.GetSpawnError()
		var html string
		if errMsg != "" {
			html = fmt.Sprintf(
				`<div class="prl-spawn-failed"><span class="form-hint" data-testid="prl-spawn-error">%s: %s</span><button type="button" class="btn btn-secondary btn-sm" hx-post="%s" hx-swap="none" hx-confirm="%s" data-testid="prl-retry-spawn">%s</button></div>`,
				htmlEsc(l.Spawn.ErrorPrefix), htmlEsc(errMsg), htmlEsc(retryURL), htmlEsc(l.Spawn.RetryConfirm), htmlEsc(l.Spawn.RetryButton),
			)
		} else {
			html = fmt.Sprintf(
				`<button type="button" class="btn btn-secondary btn-sm" hx-post="%s" hx-swap="none" hx-confirm="%s" data-testid="prl-retry-spawn">%s</button>`,
				htmlEsc(retryURL), htmlEsc(l.Spawn.RetryConfirm), htmlEsc(l.Spawn.RetryButton),
			)
		}
		return types.TableCell{Type: "html", HTML: template.HTML(html), Value: l.Spawn.RetryButton}
	}

	return types.TableCell{Type: "text", Value: l.Spawn.NotApplicable}
}

// htmlEsc escapes a string for safe interpolation into an HTML attribute or text.
func htmlEsc(s string) string {
	return template.HTMLEscapeString(s)
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
