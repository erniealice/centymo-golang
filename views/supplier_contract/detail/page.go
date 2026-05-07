package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"
	"github.com/erniealice/hybra-golang/views/attachment"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	suppliercontractlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_line"
	scpspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract_price_schedule"
)

// DetailViewDeps holds all dependencies for the supplier contract detail page.
type DetailViewDeps struct {
	Routes       centymo.SupplierContractRoutes
	Labels       centymo.SupplierContractLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadSupplierContract      func(ctx context.Context, req *suppliercontractpb.ReadSupplierContractRequest) (*suppliercontractpb.ReadSupplierContractResponse, error)
	ListSupplierContractLines func(ctx context.Context, req *suppliercontractlinepb.ListSupplierContractLinesRequest) (*suppliercontractlinepb.ListSupplierContractLinesResponse, error)

	// Linked POs — optional; list POs where supplier_contract_id = this contract
	ListPurchaseOrders func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)

	// Linked Expenditures — optional; list expenditures where supplier_contract_id = this contract
	ListExpenditures func(ctx context.Context, req *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error)

	// Price Schedules tab (SPS P7) — optional; list SCPS rows where
	// supplier_contract_id = this contract. The tab badge and add-CTA URL
	// are sourced from PriceScheduleRoutes (master-thread-injected).
	ListSupplierContractPriceSchedules func(ctx context.Context, req *scpspb.ListSupplierContractPriceSchedulesRequest) (*scpspb.ListSupplierContractPriceSchedulesResponse, error)
	PriceScheduleListURL               string // /app/supplier-contract-price-schedules/list/{status}
	PriceScheduleDetailURL             string // /app/supplier-contract-price-schedules/detail/{id}
	PriceScheduleAddURL                string // /action/supplier-contract-price-schedule/add (with ?supplier_contract_id=)

	// Workflow invocations — block.go injects closures that call the espyna use
	// cases (with user-id sourced from request context). All optional — when
	// nil the action handler degrades to a redirect-only no-op.
	ApproveSupplierContract   func(ctx context.Context, id string) error
	TerminateSupplierContract func(ctx context.Context, id string, reason string) error

	attachment.AttachmentOps
}

// PageData holds the template data for the supplier contract detail page.
type PageData struct {
	types.PageData
	ContentTemplate string

	// Contract record
	Contract      map[string]any
	StatusVariant string

	// Tab navigation
	TabItems  []pyeza.TabItem
	ActiveTab string

	// Lines tab
	LineItemTable  *types.TableConfig
	LineItemAddURL string

	// Linked POs tab
	LinkedPOTable *types.TableConfig

	// Linked Expenditures tab
	LinkedExpenditureTable *types.TableConfig

	// Price Schedules tab (SPS P7)
	PriceScheduleTable  *types.TableConfig
	PriceScheduleAddURL string

	// Attachments tab
	AttachmentTable *types.TableConfig

	// Action URLs
	ApproveURL   string
	TerminateURL string
	EditURL      string
}

const (
	tabInfo            = "info"
	tabLines           = "lines"
	tabLinkedPOs       = "linked-pos"
	tabLinkedExp       = "linked-expenditures"
	tabPriceSchedules  = "price-schedules"
	tabActivity        = "activity"
	tabAttachments     = "attachments"
)

// NewView creates the supplier contract detail page view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Redirect(deps.Routes.ListURL)
		}

		resp, err := deps.ReadSupplierContract(ctx, &suppliercontractpb.ReadSupplierContractRequest{
			Data: &suppliercontractpb.SupplierContract{Id: id},
		})
		if err != nil {
			log.Printf("ReadSupplierContract %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load supplier contract: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("supplier contract not found"))
		}
		contract := data[0]

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = tabInfo
		}

		l := deps.Labels

		contractMap := map[string]any{
			"id":               contract.GetId(),
			"name":             contract.GetName(),
			"kind":             contract.GetKind().String(),
			"status":           contract.GetStatus().String(),
			"supplier_name":    supplierName(contract),
			"start_date":       contract.GetDateTimeStart(),
			"end_date":         contract.GetDateTimeEnd(),
			"auto_renew":       contract.GetAutoRenew(),
			"currency":         contract.GetCurrency(),
			"committed_amount": contract.GetCommittedAmount(),
			"released_amount":  contract.GetReleasedAmount(),
			"billed_amount":    contract.GetBilledAmount(),
			"remaining_amount": contract.GetRemainingAmount(),
			"notes":            contract.GetNotes(),
		}

		attachmentsLabel := l.Detail.TabAttachments
		if attachmentsLabel == "" {
			attachmentsLabel = "Attachments"
		}
		tabItems := []pyeza.TabItem{
			{Key: tabInfo, Label: l.Tabs.Info},
			{Key: tabLines, Label: l.Tabs.Lines},
			{Key: tabPriceSchedules, Label: l.Tabs.PriceSchedules, Count: priceScheduleActiveCount(ctx, deps, id)},
			{Key: tabLinkedPOs, Label: l.Tabs.LinkedPOs},
			{Key: tabLinkedExp, Label: l.Tabs.LinkedExpenditures},
			{Key: tabActivity, Label: l.Tabs.Activity},
			{Key: tabAttachments, Label: attachmentsLabel, Icon: "icon-paperclip"},
		}

		pd := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          contract.GetName(),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    contract.GetName(),
				HeaderSubtitle: l.Page.DetailSubtitle,
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "supplier-contract-detail-content",
			Contract:        contractMap,
			StatusVariant:   contractStatusVariant(contract.GetStatus().String()),
			TabItems:        tabItems,
			ActiveTab:       activeTab,
			ApproveURL:      buildActionURL(deps.Routes.ApproveURL, id),
			TerminateURL:    buildActionURL(deps.Routes.TerminateURL, id),
			EditURL:         buildActionURL(deps.Routes.EditURL, id),
		}

		// Lines tab — load if active
		if activeTab == tabLines && deps.ListSupplierContractLines != nil {
			pd.LineItemTable = buildLineItemTable(ctx, deps, id, l)
			pd.LineItemAddURL = buildLineAddURL(deps.Routes.LineAddURL, id)
		}

		// Linked POs tab
		if activeTab == tabLinkedPOs && deps.ListPurchaseOrders != nil {
			pd.LinkedPOTable = buildLinkedPOTable(ctx, deps, id, l)
		}

		// Linked Expenditures tab
		if activeTab == tabLinkedExp && deps.ListExpenditures != nil {
			pd.LinkedExpenditureTable = buildLinkedExpTable(ctx, deps, id, l)
		}

		// Price Schedules tab — SPS P7
		if activeTab == tabPriceSchedules && deps.ListSupplierContractPriceSchedules != nil {
			pd.PriceScheduleTable = buildPriceScheduleTable(ctx, deps, id, l)
			pd.PriceScheduleAddURL = buildPriceScheduleAddURL(deps.PriceScheduleAddURL, id)
		}

		// Attachments tab
		if activeTab == tabAttachments && deps.ListAttachments != nil {
			cfg := attachmentConfig(deps)
			var attachItems []*attachmentpb.Attachment
			if resp, err := deps.ListAttachments(ctx, cfg.EntityType, id); err == nil && resp != nil {
				attachItems = resp.GetData()
			}
			pd.AttachmentTable = attachment.BuildTable(attachItems, cfg, id)
		}

		return view.OK("supplier-contract-detail", pd)
	})
}

// NewTabAction handles HTMX tab switch requests (/action/supplier-contract/detail/{id}/tab/{tab}).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if id == "" || tab == "" {
			return view.Error(fmt.Errorf("missing id or tab"))
		}

		resp, err := deps.ReadSupplierContract(ctx, &suppliercontractpb.ReadSupplierContractRequest{
			Data: &suppliercontractpb.SupplierContract{Id: id},
		})
		if err != nil {
			return view.Error(fmt.Errorf("failed to load supplier contract: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("supplier contract not found"))
		}
		contract := data[0]

		l := deps.Labels
		contractMap := map[string]any{
			"id":               contract.GetId(),
			"name":             contract.GetName(),
			"kind":             contract.GetKind().String(),
			"status":           contract.GetStatus().String(),
			"supplier_name":    supplierName(contract),
			"start_date":       contract.GetDateTimeStart(),
			"end_date":         contract.GetDateTimeEnd(),
			"auto_renew":       contract.GetAutoRenew(),
			"currency":         contract.GetCurrency(),
			"committed_amount": contract.GetCommittedAmount(),
			"released_amount":  contract.GetReleasedAmount(),
			"billed_amount":    contract.GetBilledAmount(),
			"remaining_amount": contract.GetRemainingAmount(),
			"notes":            contract.GetNotes(),
		}

		attachmentsLabel := l.Detail.TabAttachments
		if attachmentsLabel == "" {
			attachmentsLabel = "Attachments"
		}
		tabItems := []pyeza.TabItem{
			{Key: tabInfo, Label: l.Tabs.Info},
			{Key: tabLines, Label: l.Tabs.Lines},
			{Key: tabPriceSchedules, Label: l.Tabs.PriceSchedules, Count: priceScheduleActiveCount(ctx, deps, id)},
			{Key: tabLinkedPOs, Label: l.Tabs.LinkedPOs},
			{Key: tabLinkedExp, Label: l.Tabs.LinkedExpenditures},
			{Key: tabActivity, Label: l.Tabs.Activity},
			{Key: tabAttachments, Label: attachmentsLabel, Icon: "icon-paperclip"},
		}

		pd := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "supplier-contract-detail-content",
			Contract:        contractMap,
			StatusVariant:   contractStatusVariant(contract.GetStatus().String()),
			TabItems:        tabItems,
			ActiveTab:       tab,
			ApproveURL:      buildActionURL(deps.Routes.ApproveURL, id),
			TerminateURL:    buildActionURL(deps.Routes.TerminateURL, id),
			EditURL:         buildActionURL(deps.Routes.EditURL, id),
		}

		templateName := "supplier-contract-tab-content"
		switch tab {
		case tabLines:
			if deps.ListSupplierContractLines != nil {
				pd.LineItemTable = buildLineItemTable(ctx, deps, id, l)
				pd.LineItemAddURL = buildLineAddURL(deps.Routes.LineAddURL, id)
			}
		case tabLinkedPOs:
			if deps.ListPurchaseOrders != nil {
				pd.LinkedPOTable = buildLinkedPOTable(ctx, deps, id, l)
			}
		case tabLinkedExp:
			if deps.ListExpenditures != nil {
				pd.LinkedExpenditureTable = buildLinkedExpTable(ctx, deps, id, l)
			}
		case tabPriceSchedules:
			if deps.ListSupplierContractPriceSchedules != nil {
				pd.PriceScheduleTable = buildPriceScheduleTable(ctx, deps, id, l)
				pd.PriceScheduleAddURL = buildPriceScheduleAddURL(deps.PriceScheduleAddURL, id)
			}
		case tabAttachments:
			if deps.ListAttachments != nil {
				cfg := attachmentConfig(deps)
				var attachItems []*attachmentpb.Attachment
				if resp, err := deps.ListAttachments(ctx, cfg.EntityType, id); err == nil && resp != nil {
					attachItems = resp.GetData()
				}
				pd.AttachmentTable = attachment.BuildTable(attachItems, cfg, id)
			}
			templateName = "attachment-tab"
		}

		return view.OK(templateName, pd)
	})
}

// NewApproveAction handles POST /action/supplier-contract/approve/{id}.
// Invokes ApproveSupplierContract use case (PENDING_APPROVAL → APPROVED) and
// HX-Redirects back to the detail page.
func NewApproveAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		if deps.ApproveSupplierContract != nil {
			if err := deps.ApproveSupplierContract(ctx, id); err != nil {
				log.Printf("ApproveSupplierContract %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
		}
		detailURL := buildDetailURL(deps.Routes.DetailURL, id)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Redirect": detailURL,
			},
		}
	})
}

// NewTerminateAction handles POST /action/supplier-contract/terminate/{id}.
// Invokes TerminateSupplierContract use case and HX-Redirects back to the
// detail page. Optional reason from form body is forwarded to the use case.
func NewTerminateAction(deps *DetailViewDeps) view.View {
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
			reason = viewCtx.Request.FormValue("reason")
			if reason == "" {
				reason = viewCtx.Request.FormValue("rejection_reason")
			}
		}
		if deps.TerminateSupplierContract != nil {
			if err := deps.TerminateSupplierContract(ctx, id, reason); err != nil {
				log.Printf("TerminateSupplierContract %s: %v", id, err)
				return centymo.HTMXError(err.Error())
			}
		}
		detailURL := buildDetailURL(deps.Routes.DetailURL, id)
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Redirect": detailURL,
			},
		}
	})
}

// --- helpers -----------------------------------------------------------------

func supplierName(c *suppliercontractpb.SupplierContract) string {
	if s := c.GetSupplier(); s != nil {
		return s.GetName()
	}
	return ""
}

func optionalStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func buildActionURL(template, id string) string {
	if template == "" {
		return ""
	}
	return route.ResolveURL(template, "id", id)
}

func buildDetailURL(template, id string) string {
	if template == "" {
		return ""
	}
	return route.ResolveURL(template, "id", id)
}

func buildLineAddURL(template, contractID string) string {
	if template == "" {
		return ""
	}
	return route.ResolveURL(template, "id", contractID)
}

func contractStatusVariant(status string) string {
	switch status {
	case "SUPPLIER_CONTRACT_STATUS_ACTIVE":
		return "success"
	case "SUPPLIER_CONTRACT_STATUS_APPROVED":
		return "info"
	case "SUPPLIER_CONTRACT_STATUS_PENDING_APPROVAL",
		"SUPPLIER_CONTRACT_STATUS_REQUESTED",
		"SUPPLIER_CONTRACT_STATUS_EXPIRING",
		"SUPPLIER_CONTRACT_STATUS_SUSPENDED":
		return "warning"
	case "SUPPLIER_CONTRACT_STATUS_EXPIRED",
		"SUPPLIER_CONTRACT_STATUS_TERMINATED",
		"SUPPLIER_CONTRACT_STATUS_REJECTED":
		return "danger"
	default:
		return "default"
	}
}

func buildLineItemTable(ctx context.Context, deps *DetailViewDeps, contractID string, l centymo.SupplierContractLabels) *types.TableConfig {
	cIDPtr := contractID
	resp, err := deps.ListSupplierContractLines(ctx, &suppliercontractlinepb.ListSupplierContractLinesRequest{
		SupplierContractId: &cIDPtr,
	})
	if err != nil {
		log.Printf("ListSupplierContractLines for %s: %v", contractID, err)
		return nil
	}

	columns := []types.TableColumn{
		{Key: "description", Label: l.Lines.Description},
		{Key: "line_type", Label: l.Lines.LineType, WidthClass: "col-2xl"},
		{Key: "quantity", Label: l.Lines.Quantity, Align: "right", WidthClass: "col-xl"},
		{Key: "unit_price", Label: l.Lines.UnitPrice, Align: "right", WidthClass: "col-3xl"},
		{Key: "total", Label: l.Lines.Total, Align: "right", WidthClass: "col-3xl"},
		{Key: "treatment", Label: l.Lines.Treatment, WidthClass: "col-2xl"},
	}

	rows := []types.TableRow{}
	for _, line := range resp.GetData() {
		rows = append(rows, types.TableRow{
			ID: line.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: line.GetDescription()},
				{Type: "badge", Value: line.GetLineType(), Variant: "default"},
				{Type: "number", Value: fmt.Sprintf("%.2f", line.GetQuantity())},
				types.MoneyCell(float64(line.GetUnitPrice()), "", true),
				types.MoneyCell(float64(line.GetTotalAmount()), "", true),
				{Type: "badge", Value: line.GetTreatment().String(), Variant: "default"},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:          "contract-lines-table",
		Columns:     columns,
		Rows:        rows,
		ShowActions: true,
		EmptyState: types.TableEmptyState{
			Title:   l.Lines.EmptyTitle,
			Message: l.Lines.EmptyMessage,
		},
	}
}

func buildLinkedPOTable(ctx context.Context, deps *DetailViewDeps, contractID string, l centymo.SupplierContractLabels) *types.TableConfig {
	resp, err := deps.ListPurchaseOrders(ctx, &purchaseorderpb.ListPurchaseOrdersRequest{})
	if err != nil {
		log.Printf("ListPurchaseOrders for contract %s: %v", contractID, err)
		return nil
	}

	// Filter to POs linked to this contract
	var linked []*purchaseorderpb.PurchaseOrder
	for _, po := range resp.GetData() {
		if po.GetSupplierContractId() == contractID {
			linked = append(linked, po)
		}
	}

	columns := []types.TableColumn{
		{Key: "po_number", Label: l.LinkedPOs.PONumber},
		{Key: "status", Label: l.LinkedPOs.Status, WidthClass: "col-2xl"},
		{Key: "total_amount", Label: l.LinkedPOs.TotalAmount, Align: "right", WidthClass: "col-3xl"},
		{Key: "order_date", Label: l.LinkedPOs.OrderDate, WidthClass: "col-3xl"},
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
		ID:      "linked-pos-table",
		Columns: columns,
		Rows:    rows,
		EmptyState: types.TableEmptyState{
			Title:   l.LinkedPOs.EmptyTitle,
			Message: l.LinkedPOs.EmptyMessage,
		},
	}
}

func buildLinkedExpTable(ctx context.Context, deps *DetailViewDeps, contractID string, l centymo.SupplierContractLabels) *types.TableConfig {
	resp, err := deps.ListExpenditures(ctx, &expenditurepb.ListExpendituresRequest{})
	if err != nil {
		log.Printf("ListExpenditures for contract %s: %v", contractID, err)
		return nil
	}

	// Filter to expenditures linked to this contract
	var linked []*expenditurepb.Expenditure
	for _, e := range resp.GetData() {
		if e.GetSupplierContractId() == contractID {
			linked = append(linked, e)
		}
	}

	columns := []types.TableColumn{
		{Key: "reference", Label: l.LinkedExpenditures.Reference},
		{Key: "status", Label: l.LinkedExpenditures.Status, WidthClass: "col-2xl"},
		{Key: "amount", Label: l.LinkedExpenditures.Amount, Align: "right", WidthClass: "col-3xl"},
		{Key: "date", Label: l.LinkedExpenditures.Date, WidthClass: "col-3xl"},
	}

	rows := []types.TableRow{}
	for _, e := range linked {
		rows = append(rows, types.TableRow{
			ID: e.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: e.GetReferenceNumber()},
				{Type: "badge", Value: e.GetStatus(), Variant: "default"},
				types.MoneyCell(float64(e.GetTotalAmount()), e.GetCurrency(), true),
				types.DateTimeCell(e.GetExpenditureDateString(), types.DateReadable),
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "linked-expenditures-table",
		Columns: columns,
		Rows:    rows,
		EmptyState: types.TableEmptyState{
			Title:   l.LinkedExpenditures.EmptyTitle,
			Message: l.LinkedExpenditures.EmptyMessage,
		},
	}
}

// --- SPS P7 — Price Schedules tab helpers -----------------------------------

// priceScheduleActiveCount counts ACTIVE schedules attached to this contract.
// Returns 0 on error or when the dep isn't wired (graceful degradation).
func priceScheduleActiveCount(ctx context.Context, deps *DetailViewDeps, contractID string) int {
	if deps.ListSupplierContractPriceSchedules == nil {
		return 0
	}
	cIDPtr := contractID
	resp, err := deps.ListSupplierContractPriceSchedules(ctx, &scpspb.ListSupplierContractPriceSchedulesRequest{
		SupplierContractId: &cIDPtr,
	})
	if err != nil {
		return 0
	}
	count := 0
	for _, s := range resp.GetData() {
		if s.GetStatus() == scpspb.SupplierContractPriceScheduleStatus_SUPPLIER_CONTRACT_PRICE_SCHEDULE_STATUS_ACTIVE {
			count++
		}
	}
	return count
}

func buildPriceScheduleTable(ctx context.Context, deps *DetailViewDeps, contractID string, l centymo.SupplierContractLabels) *types.TableConfig {
	cIDPtr := contractID
	resp, err := deps.ListSupplierContractPriceSchedules(ctx, &scpspb.ListSupplierContractPriceSchedulesRequest{
		SupplierContractId: &cIDPtr,
	})
	if err != nil {
		log.Printf("ListSupplierContractPriceSchedules for contract %s: %v", contractID, err)
		return nil
	}

	columns := []types.TableColumn{
		{Key: "name", Label: l.Tabs.PriceSchedules},
		{Key: "sequence", Label: "Seq.", Align: "right", WidthClass: "col-xs"},
		{Key: "period", Label: "Period", WidthClass: "col-3xl"},
		{Key: "status", Label: "Status", WidthClass: "col-2xl"},
		{Key: "currency", Label: "Currency", WidthClass: "col-xs"},
	}

	rows := []types.TableRow{}
	for _, s := range resp.GetData() {
		statusStr := s.GetStatus().String()
		startStr := ""
		if t := s.GetDateTimeStart(); t != nil && t.IsValid() {
			startStr = t.AsTime().UTC().Format("2006-01-02")
		}
		endStr := "—"
		if t := s.GetDateTimeEnd(); t != nil && t.IsValid() {
			endStr = t.AsTime().UTC().Format("2006-01-02")
		}
		period := startStr
		if endStr != "" {
			period = startStr + " → " + endStr
		}
		detailURL := ""
		if deps.PriceScheduleDetailURL != "" {
			detailURL = route.ResolveURL(deps.PriceScheduleDetailURL, "id", s.GetId())
		}
		nameCell := types.TableCell{Type: "text", Value: s.GetName()}
		if detailURL != "" {
			nameCell = types.TableCell{Type: "link", Value: s.GetName(), Href: detailURL}
		}
		rows = append(rows, types.TableRow{
			ID: s.GetId(),
			Cells: []types.TableCell{
				nameCell,
				{Type: "number", Value: fmt.Sprintf("%d", s.GetSequenceNumber())},
				{Type: "text", Value: period},
				{Type: "badge", Value: statusStr, Variant: priceScheduleStatusVariant(statusStr)},
				{Type: "text", Value: s.GetCurrency()},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "contract-price-schedules-table",
		Columns: columns,
		Rows:    rows,
		EmptyState: types.TableEmptyState{
			Title:   l.Tabs.PriceSchedules,
			Message: l.Tabs.PriceSchedulesEmpty,
		},
	}
}

func buildPriceScheduleAddURL(template, contractID string) string {
	if template == "" {
		return ""
	}
	// Pre-fill the parent contract via query string so the drawer pre-selects
	// the correct supplier_contract_id when opened from this tab.
	if contractID == "" {
		return template
	}
	return template + "?supplier_contract_id=" + contractID
}

func priceScheduleStatusVariant(status string) string {
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
