package detail

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	expenditurelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_line_item"
)

// DetailViewDeps holds view dependencies for the expense detail page.
type DetailViewDeps struct {
	Routes       centymo.ExpenditureRoutes
	Labels       centymo.ExpenditureLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadExpenditure          func(ctx context.Context, req *expenditurepb.ReadExpenditureRequest) (*expenditurepb.ReadExpenditureResponse, error)
	ListExpenditureLineItems func(ctx context.Context, req *expenditurelineitempb.ListExpenditureLineItemsRequest) (*expenditurelineitempb.ListExpenditureLineItemsResponse, error)
}

// LineItemDeps holds dependencies for line item action handlers.
type LineItemDeps struct {
	Routes       centymo.ExpenditureRoutes
	Labels       centymo.ExpenditureLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadExpenditure           func(ctx context.Context, req *expenditurepb.ReadExpenditureRequest) (*expenditurepb.ReadExpenditureResponse, error)
	UpdateExpenditure         func(ctx context.Context, req *expenditurepb.UpdateExpenditureRequest) (*expenditurepb.UpdateExpenditureResponse, error)
	CreateExpenditureLineItem func(ctx context.Context, req *expenditurelineitempb.CreateExpenditureLineItemRequest) (*expenditurelineitempb.CreateExpenditureLineItemResponse, error)
	ReadExpenditureLineItem   func(ctx context.Context, req *expenditurelineitempb.ReadExpenditureLineItemRequest) (*expenditurelineitempb.ReadExpenditureLineItemResponse, error)
	UpdateExpenditureLineItem func(ctx context.Context, req *expenditurelineitempb.UpdateExpenditureLineItemRequest) (*expenditurelineitempb.UpdateExpenditureLineItemResponse, error)
	DeleteExpenditureLineItem func(ctx context.Context, req *expenditurelineitempb.DeleteExpenditureLineItemRequest) (*expenditurelineitempb.DeleteExpenditureLineItemResponse, error)
	ListExpenditureLineItems  func(ctx context.Context, req *expenditurelineitempb.ListExpenditureLineItemsRequest) (*expenditurelineitempb.ListExpenditureLineItemsResponse, error)
}

// LineItemFormData is the template data for the line item drawer form.
type LineItemFormData struct {
	FormAction    string
	IsEdit        bool
	ID            string
	ExpenditureID string
	Description   string
	Quantity      string
	UnitPrice     string
	Notes         string
	CommonLabels  any
	Labels        centymo.ExpenditureLabels
}

// PageData holds the data for the expense detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Expense         map[string]any
	Labels          centymo.ExpenditureLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	LineItemTable   *types.TableConfig
	LineItemAddURL  string
	TotalAmount     string
	SetStatusURL    string
}

// expenditureToMap converts an Expenditure proto to a map for template use.
func expenditureToMap(e *expenditurepb.Expenditure) map[string]any {
	return map[string]any{
		"id":                   e.GetId(),
		"name":                 e.GetName(),
		"reference_number":     e.GetReferenceNumber(),
		"expenditure_type":     e.GetExpenditureType(),
		"total_amount":         centymo.FormatWithCommas(float64(e.GetTotalAmount()) / 100.0),
		"currency":             e.GetCurrency(),
		"status":               e.GetStatus(),
		"notes":                e.GetNotes(),
		"active":               e.GetActive(),
		"date_created_string":  e.GetDateCreatedString(),
		"date_modified_string": e.GetDateModifiedString(),
	}
}

// NewView creates the expense detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		resp, err := deps.ReadExpenditure(ctx, &expenditurepb.ReadExpenditureRequest{
			Data: &expenditurepb.Expenditure{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read expenditure %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load expense: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("expense not found"))
		}
		expense := expenditureToMap(data[0])

		refNumber, _ := expense["reference_number"].(string)
		headerTitle := refNumber
		if headerTitle == "" {
			headerTitle, _ = expense["name"].(string)
		}

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}
		tabItems := buildTabItems(deps.Labels, id, deps.Routes)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "expenses",
				HeaderTitle:    headerTitle,
				HeaderSubtitle: deps.Labels.Page.ExpenseHeading,
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "expense-detail-content",
			Expense:         expense,
			Labels:          deps.Labels,
			ActiveTab:       activeTab,
			TabItems:        tabItems,
			SetStatusURL:    deps.Routes.SetStatusURL,
		}

		// Load tab-specific data
		switch activeTab {
		case "info":
			// expense map has everything
		case "items":
			if deps.ListExpenditureLineItems != nil {
				perms := view.GetUserPermissions(ctx)
				lineItems := listLineItemMaps(ctx, deps.ListExpenditureLineItems, id)
				currency, _ := expense["currency"].(string)
				pageData.LineItemTable = buildLineItemTable(lineItems, deps.Labels, deps.TableLabels, currency, id, deps.Routes, perms)
				pageData.LineItemAddURL = route.ResolveURL(deps.Routes.LineItemAddURL, "id", id)
				totalAmount, _ := expense["total_amount"].(string)
				pageData.TotalAmount = currency + " " + totalAmount
			}
		}

		return view.OK("expense-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		resp, err := deps.ReadExpenditure(ctx, &expenditurepb.ReadExpenditureRequest{
			Data: &expenditurepb.Expenditure{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read expenditure %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load expense: %w", err))
		}
		data := resp.GetData()
		if len(data) == 0 {
			return view.Error(fmt.Errorf("expense not found"))
		}
		expense := expenditureToMap(data[0])

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Expense:      expense,
			Labels:       deps.Labels,
			ActiveTab:    tab,
			TabItems:     buildTabItems(deps.Labels, id, deps.Routes),
			SetStatusURL: deps.Routes.SetStatusURL,
		}

		switch tab {
		case "info":
			// expense map has everything
		case "items":
			if deps.ListExpenditureLineItems != nil {
				perms := view.GetUserPermissions(ctx)
				lineItems := listLineItemMaps(ctx, deps.ListExpenditureLineItems, id)
				currency, _ := expense["currency"].(string)
				pageData.LineItemTable = buildLineItemTable(lineItems, deps.Labels, deps.TableLabels, currency, id, deps.Routes, perms)
				pageData.LineItemAddURL = route.ResolveURL(deps.Routes.LineItemAddURL, "id", id)
				totalAmount, _ := expense["total_amount"].(string)
				pageData.TotalAmount = currency + " " + totalAmount
			}
		}

		templateName := "expense-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

// NewLineItemTableView returns a view that renders only the line items table (HTMX refresh).
func NewLineItemTableView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		expenditureID := viewCtx.Request.PathValue("id")

		lineItems := listLineItemMaps(ctx, deps.ListExpenditureLineItems, expenditureID)
		perms := view.GetUserPermissions(ctx)
		table := buildLineItemTable(lineItems, deps.Labels, deps.TableLabels, "", expenditureID, deps.Routes, perms)
		return view.OK("table-card", table)
	})
}

// NewLineItemAddView creates the line item add action (GET = form, POST = create).
func NewLineItemAddView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense", "update") {
			return lineItemHTMXError("Permission denied")
		}

		expenditureID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("expense-line-item-drawer-form", &LineItemFormData{
				FormAction:    route.ResolveURL(deps.Routes.LineItemAddURL, "id", expenditureID),
				ExpenditureID: expenditureID,
				Quantity:      "1",
				Labels:        deps.Labels,
				CommonLabels:  nil,
			})
		}

		// POST — create line item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError("Invalid form data")
		}

		r := viewCtx.Request
		quantity := r.FormValue("quantity")
		unitPrice := r.FormValue("unit_price")

		quantityF, _ := strconv.ParseFloat(quantity, 64)
		unitPriceF, _ := strconv.ParseFloat(unitPrice, 64)
		if quantityF == 0 {
			quantityF = 1
		}
		total := quantityF * unitPriceF

		notesCreate := r.FormValue("notes")
		_, err := deps.CreateExpenditureLineItem(ctx, &expenditurelineitempb.CreateExpenditureLineItemRequest{
			Data: &expenditurelineitempb.ExpenditureLineItem{
				ExpenditureId: expenditureID,
				Description:   r.FormValue("description"),
				Quantity:      quantityF,
				UnitPrice:     int64(math.Round(unitPriceF * 100)),
				TotalPrice:    int64(math.Round(total * 100)),
				Notes:         &notesCreate,
			},
		})
		if err != nil {
			log.Printf("Failed to create expenditure line item: %v", err)
			return lineItemHTMXError(err.Error())
		}

		// Recalculate total
		recalculateExpenseTotal(ctx, deps.ListExpenditureLineItems, deps.UpdateExpenditure, expenditureID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// NewLineItemEditView creates the line item edit action (GET = form, POST = update).
func NewLineItemEditView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense", "update") {
			return lineItemHTMXError("Permission denied")
		}

		expenditureID := viewCtx.Request.PathValue("id")
		itemID := viewCtx.Request.PathValue("itemId")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadExpenditureLineItem(ctx, &expenditurelineitempb.ReadExpenditureLineItemRequest{
				Data: &expenditurelineitempb.ExpenditureLineItem{Id: itemID},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				return lineItemHTMXError("Line item not found")
			}
			item := readResp.GetData()[0]

			return view.OK("expense-line-item-drawer-form", &LineItemFormData{
				FormAction:    route.ResolveURL(deps.Routes.LineItemEditURL, "id", expenditureID, "itemId", itemID),
				IsEdit:        true,
				ID:            itemID,
				ExpenditureID: expenditureID,
				Description:   item.GetDescription(),
				Quantity:      fmt.Sprintf("%.0f", item.GetQuantity()),
				UnitPrice:     fmt.Sprintf("%.2f", float64(item.GetUnitPrice())/100.0),
				Notes:         item.GetNotes(),
				Labels:        deps.Labels,
				CommonLabels:  nil,
			})
		}

		// POST — update line item
		if err := viewCtx.Request.ParseForm(); err != nil {
			return lineItemHTMXError("Invalid form data")
		}

		r := viewCtx.Request
		quantity := r.FormValue("quantity")
		unitPrice := r.FormValue("unit_price")

		quantityF, _ := strconv.ParseFloat(quantity, 64)
		unitPriceF, _ := strconv.ParseFloat(unitPrice, 64)
		if quantityF == 0 {
			quantityF = 1
		}
		total := quantityF * unitPriceF

		notesUpdate := r.FormValue("notes")
		_, err := deps.UpdateExpenditureLineItem(ctx, &expenditurelineitempb.UpdateExpenditureLineItemRequest{
			Data: &expenditurelineitempb.ExpenditureLineItem{
				Id:          itemID,
				Description: r.FormValue("description"),
				Quantity:    quantityF,
				UnitPrice:   int64(math.Round(unitPriceF * 100)),
				TotalPrice:  int64(math.Round(total * 100)),
				Notes:       &notesUpdate,
			},
		})
		if err != nil {
			log.Printf("Failed to update expenditure line item: %v", err)
			return lineItemHTMXError(err.Error())
		}

		recalculateExpenseTotal(ctx, deps.ListExpenditureLineItems, deps.UpdateExpenditure, expenditureID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// NewLineItemRemoveView creates the line item remove action (POST only).
func NewLineItemRemoveView(deps *LineItemDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expense", "update") {
			return lineItemHTMXError("Permission denied")
		}

		expenditureID := viewCtx.Request.PathValue("id")

		itemID := viewCtx.Request.URL.Query().Get("itemId")
		if itemID == "" {
			_ = viewCtx.Request.ParseForm()
			itemID = viewCtx.Request.FormValue("itemId")
		}
		if itemID == "" {
			return lineItemHTMXError("Item ID required")
		}

		_, err := deps.DeleteExpenditureLineItem(ctx, &expenditurelineitempb.DeleteExpenditureLineItemRequest{
			Data: &expenditurelineitempb.ExpenditureLineItem{Id: itemID},
		})
		if err != nil {
			log.Printf("Failed to delete expenditure line item: %v", err)
			return lineItemHTMXError(err.Error())
		}

		recalculateExpenseTotal(ctx, deps.ListExpenditureLineItems, deps.UpdateExpenditure, expenditureID)

		return lineItemHTMXSuccess("line-items-table")
	})
}

// buildTabItems builds the tab navigation for the expense detail page.
func buildTabItems(l centymo.ExpenditureLabels, id string, routes centymo.ExpenditureRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: "Details", Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "items", Label: "Line Items", Href: base + "?tab=items", HxGet: action + "items", Icon: "icon-list"},
	}
}

// buildLineItemTable builds the line items table config.
func buildLineItemTable(items []map[string]any, l centymo.ExpenditureLabels, tableLabels types.TableLabels, currency string, expenditureID string, routes centymo.ExpenditureRoutes, perms *types.UserPermissions) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "description", Label: "Description", Sortable: false},
		{Key: "quantity", Label: "Qty", Sortable: false, Width: "80px"},
		{Key: "unit_price", Label: "Unit Price", Sortable: false, Width: "130px"},
		{Key: "total", Label: "Total", Sortable: false, Width: "130px"},
	}

	rows := []types.TableRow{}
	for _, item := range items {
		id, _ := item["id"].(string)
		description, _ := item["description"].(string)
		quantity, _ := item["quantity"].(string)
		unitPrice, _ := item["unit_price"].(string)
		total, _ := item["total"].(string)

		actions := []types.TableAction{
			{
				Type:            "edit",
				Label:           "Edit",
				Action:          "edit",
				URL:             route.ResolveURL(routes.LineItemEditURL, "id", expenditureID, "itemId", id),
				DrawerTitle:     "Edit Line Item",
				Disabled:        !perms.Can("expense", "update"),
				DisabledTooltip: "No permission",
			},
			{
				Type:            "delete",
				Label:           "Remove",
				Action:          "delete",
				URL:             route.ResolveURL(routes.LineItemRemoveURL, "id", expenditureID) + "?itemId=" + id,
				ItemName:        description,
				Disabled:        !perms.Can("expense", "update"),
				DisabledTooltip: "No permission",
			},
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: description},
				{Type: "text", Value: quantity},
				{Type: "text", Value: currency + " " + unitPrice},
				{Type: "text", Value: currency + " " + total},
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "line-items-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   "No line items",
			Message: "Add line items to this expense.",
		},
	}
}

// listLineItemMaps lists line items for an expenditure and returns as maps.
func listLineItemMaps(ctx context.Context, listFn func(context.Context, *expenditurelineitempb.ListExpenditureLineItemsRequest) (*expenditurelineitempb.ListExpenditureLineItemsResponse, error), expenditureID string) []map[string]any {
	resp, err := listFn(ctx, &expenditurelineitempb.ListExpenditureLineItemsRequest{
		ExpenditureId: &expenditureID,
	})
	if err != nil {
		log.Printf("Failed to list line items for expenditure %s: %v", expenditureID, err)
		return []map[string]any{}
	}
	items := []map[string]any{}
	for _, item := range resp.GetData() {
		if item.GetExpenditureId() == expenditureID {
			items = append(items, map[string]any{
				"id":             item.GetId(),
				"expenditure_id": item.GetExpenditureId(),
				"description":    item.GetDescription(),
				"quantity":       fmt.Sprintf("%.0f", item.GetQuantity()),
				"unit_price":     centymo.FormatWithCommas(float64(item.GetUnitPrice()) / 100.0),
				"total":          centymo.FormatWithCommas(float64(item.GetTotalPrice()) / 100.0),
				"notes":          item.GetNotes(),
			})
		}
	}
	return items
}

// recalculateExpenseTotal recalculates and updates the expenditure total from line items.
func recalculateExpenseTotal(
	ctx context.Context,
	listFn func(context.Context, *expenditurelineitempb.ListExpenditureLineItemsRequest) (*expenditurelineitempb.ListExpenditureLineItemsResponse, error),
	updateFn func(context.Context, *expenditurepb.UpdateExpenditureRequest) (*expenditurepb.UpdateExpenditureResponse, error),
	expenditureID string,
) {
	if listFn == nil || updateFn == nil {
		return
	}
	resp, err := listFn(ctx, &expenditurelineitempb.ListExpenditureLineItemsRequest{
		ExpenditureId: &expenditureID,
	})
	if err != nil {
		log.Printf("Failed to list line items for total recalculation: %v", err)
		return
	}

	var total int64
	for _, item := range resp.GetData() {
		if item.GetExpenditureId() == expenditureID {
			total += item.GetTotalPrice()
		}
	}

	_, err = updateFn(ctx, &expenditurepb.UpdateExpenditureRequest{
		Data: &expenditurepb.Expenditure{
			Id:          expenditureID,
			TotalAmount: total,
		},
	})
	if err != nil {
		log.Printf("Failed to update expenditure total: %v", err)
	}
}

// lineItemHTMXSuccess returns a success HTMX response.
func lineItemHTMXSuccess(tableID string) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": fmt.Sprintf(`{"formSuccess":true,"refreshTable":"%s"}`, tableID),
		},
	}
}

// lineItemHTMXError returns an error HTMX response.
func lineItemHTMXError(message string) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusUnprocessableEntity,
		Headers: map[string]string{
			"HX-Error-Message": message,
		},
	}
}
