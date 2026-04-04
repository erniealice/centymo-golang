package list

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
)

// ListViewDeps holds view dependencies for the purchase order list.
type ListViewDeps struct {
	ListPurchaseOrders func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)
	RefreshURL         string
	AddURL             string // action URL for the add drawer
	Labels             centymo.ExpenditureLabels
	CommonLabels       pyeza.CommonLabels
	TableLabels        types.TableLabels
}

// PageData holds the data for the purchase order list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the purchase order list view, optionally filtered by status.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "all"
		}

		resp, err := deps.ListPurchaseOrders(ctx, &purchaseorderpb.ListPurchaseOrdersRequest{})
		if err != nil {
			log.Printf("Failed to list purchase orders: %v", err)
			return view.Error(fmt.Errorf("failed to load purchase orders: %w", err))
		}

		// Filter by status if not "all"
		var filtered []*purchaseorderpb.PurchaseOrder
		for _, po := range resp.GetData() {
			if status == "all" || po.GetStatus() == status {
				filtered = append(filtered, po)
			}
		}

		l := deps.Labels
		columns := purchaseOrderColumns()
		rows := buildTableRows(filtered, l)
		types.ApplyColumnStyles(columns, rows)

		heading := statusPageTitle(l, status)
		caption := l.Page.PurchaseCaption

		var primaryAction *types.PrimaryAction
		if deps.AddURL != "" {
			primaryAction = &types.PrimaryAction{
				Label:     l.Buttons.AddPurchase,
				ActionURL: deps.AddURL,
			}
		}

		tableConfig := &types.TableConfig{
			ID:                   "purchase-orders-table",
			RefreshURL:           deps.RefreshURL,
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
			DefaultSortColumn:    "order_date",
			DefaultSortDirection: "desc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   statusEmptyTitle(l, status),
				Message: statusEmptyMessage(l, status),
			},
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "purchase",
				ActiveSubNav:   status,
				HeaderTitle:    heading,
				HeaderSubtitle: caption,
				HeaderIcon:     "icon-shopping-cart",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "purchase-order-list-content",
			Table:           tableConfig,
		}

		return view.OK("purchase-order-list", pageData)
	})
}

func purchaseOrderColumns() []types.TableColumn {
	return []types.TableColumn{
		{Key: "po_number", Label: "PO Number", Sortable: true},
		{Key: "supplier", Label: "Supplier", Sortable: true},
		{Key: "status", Label: "Status", Sortable: true, Width: "120px"},
		{Key: "total_amount", Label: "Total Amount", Sortable: true, Width: "140px", Align: "right"},
		{Key: "order_date", Label: "Order Date", Sortable: true, Width: "140px"},
		{Key: "expected_delivery", Label: "Expected Delivery", Sortable: true, Width: "160px"},
	}
}

func buildTableRows(orders []*purchaseorderpb.PurchaseOrder, l centymo.ExpenditureLabels) []types.TableRow {
	rows := []types.TableRow{}
	for _, po := range orders {
		id := po.GetId()
		poNumber := po.GetPoNumber()
		currency := po.GetCurrency()
		recordStatus := po.GetStatus()
		totalAmount := centymo.FormatCentavoAmount(po.GetTotalAmount(), currency)
		orderDate := po.GetOrderDateString()
		expectedDelivery := po.GetExpectedDeliveryDateString()

		supplierName := ""
		if supplier := po.GetSupplier(); supplier != nil {
			supplierName = supplier.GetCompanyName()
		}
		if supplierName == "" {
			supplierName = po.GetSupplierId()
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: poNumber},
				{Type: "text", Value: supplierName},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
				{Type: "text", Value: totalAmount},
				{Type: "text", Value: orderDate},
				{Type: "text", Value: expectedDelivery},
			},
			DataAttrs: map[string]string{
				"po_number":         poNumber,
				"supplier":          supplierName,
				"status":            recordStatus,
				"total_amount":      totalAmount,
				"order_date":        orderDate,
				"expected_delivery": expectedDelivery,
			},
		})
	}
	return rows
}

func statusPageTitle(l centymo.ExpenditureLabels, status string) string {
	switch status {
	case "draft":
		return l.Page.PurchaseHeadingDraft
	case "pending":
		return l.Page.PurchaseHeadingPending
	case "approved":
		return l.Page.PurchaseHeadingApproved
	case "paid":
		return l.Page.PurchaseHeadingPaid
	case "cancelled":
		return l.Page.PurchaseHeadingCancelled
	case "overdue":
		return l.Page.PurchaseHeadingOverdue
	default:
		return l.Labels.PurchaseOrder
	}
}

func statusEmptyTitle(l centymo.ExpenditureLabels, status string) string {
	switch status {
	case "draft":
		return l.Empty.PurchaseDraftTitle
	case "pending":
		return l.Empty.PurchasePendingTitle
	case "approved":
		return l.Empty.PurchaseApprovedTitle
	case "paid":
		return l.Empty.PurchasePaidTitle
	case "cancelled":
		return l.Empty.PurchaseCancelledTitle
	case "overdue":
		return l.Empty.PurchaseOverdueTitle
	default:
		return l.Empty.PurchaseTitle
	}
}

func statusEmptyMessage(l centymo.ExpenditureLabels, status string) string {
	switch status {
	case "draft":
		return l.Empty.PurchaseDraftMessage
	case "pending":
		return l.Empty.PurchasePendingMessage
	case "approved":
		return l.Empty.PurchaseApprovedMessage
	case "paid":
		return l.Empty.PurchasePaidMessage
	case "cancelled":
		return l.Empty.PurchaseCancelledMessage
	case "overdue":
		return l.Empty.PurchaseOverdueMessage
	default:
		return l.Empty.PurchaseMessage
	}
}

func statusVariant(status string) string {
	switch status {
	case "draft":
		return "default"
	case "pending":
		return "warning"
	case "approved":
		return "info"
	case "partially_received":
		return "info"
	case "fully_received":
		return "success"
	case "closed":
		return "success"
	case "cancelled":
		return "danger"
	default:
		return "default"
	}
}
