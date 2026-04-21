package action

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	pyeza "github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	expenditurecategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_category"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
)

// PaymentTermOption is a minimal struct for rendering payment term options in the form.
type PaymentTermOption struct {
	Id      string
	Name    string
	NetDays int32
}

// PurchaseOrderOption is a minimal struct for rendering purchase order options in the form.
type PurchaseOrderOption struct {
	Id           string
	PoNumber     string
	SupplierName string
}

// FormLabels holds flat i18n labels for the expense drawer form template.
type FormLabels struct {
	Name                string
	NamePlaceholder     string
	Category            string
	Supplier            string
	Date                string
	Amount              string
	Currency            string
	ReferenceNumber     string
	Notes               string
	NotesPlaceholder    string
	Status              string
	ExpenditureType     string
	TypeExpense         string
	TypePurchase        string
	StatusPending       string
	StatusApproved      string
	StatusPaid          string
	StatusCancelled     string
	CurrencyPlaceholder  string
	PaymentTerms         string
	SelectPaymentTerm    string
	DueDate              string
	LinkToPurchaseOrder  string
}

// FormData is the template data for the expense drawer form.
type FormData struct {
	FormAction            string
	IsEdit                bool
	ID                    string
	Name                  string
	ExpenditureType       string
	ExpenditureCategoryID string
	SupplierID            string
	Date                  string
	TotalAmount           string
	Currency              string
	Status                string
	ReferenceNumber       string
	Notes                 string
	Categories            []pyeza.SelectOption
	Suppliers             []pyeza.SelectOption
	PaymentTerms          []*PaymentTermOption
	SelectedPaymentTermID string
	PurchaseOrders        []*PurchaseOrderOption
	PurchaseOrderID       string
	Labels                FormLabels
	CommonLabels          any
}

// Deps holds dependencies for expense action handlers.
type Deps struct {
	Routes centymo.ExpenditureRoutes
	Labels centymo.ExpenditureLabels

	// Payment terms dropdown (optional — gracefully degrades when nil)
	ListPaymentTerms func(ctx context.Context) ([]*PaymentTermOption, error)

	// Typed expenditure CRUD operations
	CreateExpenditure func(ctx context.Context, req *expenditurepb.CreateExpenditureRequest) (*expenditurepb.CreateExpenditureResponse, error)
	ReadExpenditure   func(ctx context.Context, req *expenditurepb.ReadExpenditureRequest) (*expenditurepb.ReadExpenditureResponse, error)
	UpdateExpenditure func(ctx context.Context, req *expenditurepb.UpdateExpenditureRequest) (*expenditurepb.UpdateExpenditureResponse, error)
	DeleteExpenditure func(ctx context.Context, req *expenditurepb.DeleteExpenditureRequest) (*expenditurepb.DeleteExpenditureResponse, error)

	// Category listing (optional — gracefully degrades to empty list if nil)
	ListExpenditureCategories func(ctx context.Context, req *expenditurecategorypb.ListExpenditureCategoriesRequest) (*expenditurecategorypb.ListExpenditureCategoriesResponse, error)

	// Supplier listing (optional — gracefully degrades to empty list if nil)
	ListSuppliers func(ctx context.Context, req *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error)

	// Purchase order listing (optional — gracefully degrades to empty list if nil)
	ListPurchaseOrders func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)
}

// formLabels maps ExpenditureLabels into the flat FormLabels struct for the template.
func formLabels(l centymo.ExpenditureLabels) FormLabels {
	return FormLabels{
		Name:                l.Form.VendorName,
		NamePlaceholder:     l.Form.VendorNamePlaceholder,
		Category:            l.Form.ExpenditureCategory,
		Supplier:            "Supplier",
		Date:                l.Form.ExpenditureDate,
		Amount:              l.Form.TotalAmount,
		Currency:            l.Form.Currency,
		ReferenceNumber:     l.Form.ReferenceNumber,
		Notes:               l.Form.Notes,
		NotesPlaceholder:    l.Form.NotesPlaceholder,
		Status:              l.Form.Status,
		ExpenditureType:     l.Form.ExpenditureType,
		TypeExpense:         l.Types.Expense,
		TypePurchase:        l.Types.Purchase,
		StatusPending:       l.Status.Pending,
		StatusApproved:      l.Status.Approved,
		StatusPaid:          l.Status.Paid,
		StatusCancelled:     l.Status.Cancelled,
		CurrencyPlaceholder: "e.g. PHP",
		PaymentTerms:        "Payment Terms",
		SelectPaymentTerm:   "Select payment term",
		DueDate:             "Due Date",
		LinkToPurchaseOrder: "Link to Purchase Order",
	}
}

// loadPaymentTerms fetches payment term options. Returns nil on error (graceful degradation).
func loadPaymentTerms(ctx context.Context, deps *Deps) []*PaymentTermOption {
	if deps.ListPaymentTerms == nil {
		return nil
	}
	terms, err := deps.ListPaymentTerms(ctx)
	if err != nil {
		log.Printf("Failed to load payment terms: %v", err)
		return nil
	}
	return terms
}

// loadCategoryOptions loads expenditure categories for the dropdown, pre-selecting selectedID.
func loadCategoryOptions(
	ctx context.Context,
	listFn func(ctx context.Context, req *expenditurecategorypb.ListExpenditureCategoriesRequest) (*expenditurecategorypb.ListExpenditureCategoriesResponse, error),
	selectedID string,
) []pyeza.SelectOption {
	if listFn == nil {
		return nil
	}
	resp, err := listFn(ctx, &expenditurecategorypb.ListExpenditureCategoriesRequest{})
	if err != nil {
		log.Printf("Failed to list expenditure categories: %v", err)
		return nil
	}
	var opts []pyeza.SelectOption
	for _, cat := range resp.GetData() {
		if !cat.GetActive() {
			continue
		}
		opts = append(opts, pyeza.SelectOption{
			Value:    cat.GetId(),
			Label:    cat.GetName(),
			Selected: cat.GetId() == selectedID,
		})
	}
	return opts
}

// loadSupplierOptions loads suppliers for the dropdown, pre-selecting selectedID.
func loadSupplierOptions(
	ctx context.Context,
	listFn func(ctx context.Context, req *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error),
	selectedID string,
) []pyeza.SelectOption {
	if listFn == nil {
		return nil
	}
	resp, err := listFn(ctx, &supplierpb.ListSuppliersRequest{})
	if err != nil {
		log.Printf("Failed to list suppliers: %v", err)
		return nil
	}
	var opts []pyeza.SelectOption
	for _, s := range resp.GetData() {
		if !s.GetActive() {
			continue
		}
		label := s.GetCompanyName()
		if label == "" {
			label = s.GetId()
		}
		opts = append(opts, pyeza.SelectOption{
			Value:    s.GetId(),
			Label:    label,
			Selected: s.GetId() == selectedID,
		})
	}
	return opts
}

// loadPurchaseOrderOptions loads purchase orders for the dropdown.
func loadPurchaseOrderOptions(
	ctx context.Context,
	listFn func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error),
) []*PurchaseOrderOption {
	if listFn == nil {
		return nil
	}
	resp, err := listFn(ctx, &purchaseorderpb.ListPurchaseOrdersRequest{})
	if err != nil {
		log.Printf("Failed to list purchase orders: %v", err)
		return nil
	}
	var opts []*PurchaseOrderOption
	for _, po := range resp.GetData() {
		if !po.GetActive() {
			continue
		}
		supplierName := ""
		if s := po.GetSupplier(); s != nil {
			supplierName = s.GetCompanyName()
		}
		opts = append(opts, &PurchaseOrderOption{
			Id:           po.GetId(),
			PoNumber:     po.GetPoNumber(),
			SupplierName: supplierName,
		})
	}
	return opts
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}

// parseAmount converts a form string amount (decimal) to int64 centavos.
func parseAmount(s string) int64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f * 100))
}

// NewAddAction creates the expense add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expenditure", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			paymentTerms := loadPaymentTerms(ctx, deps)
			return view.OK("expense-drawer-form", &FormData{
				FormAction:      deps.Routes.AddURL,
				ExpenditureType: "expense",
				Currency:        "PHP",
				Status:          "pending",
				Categories:      loadCategoryOptions(ctx, deps.ListExpenditureCategories, ""),
				Suppliers:       loadSupplierOptions(ctx, deps.ListSuppliers, ""),
				PaymentTerms:    paymentTerms,
				PurchaseOrders:  loadPurchaseOrderOptions(ctx, deps.ListPurchaseOrders),
				Labels:          formLabels(deps.Labels),
				CommonLabels:    nil, // injected by ViewAdapter
			})
		}

		// POST — create expense
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		resp, err := deps.CreateExpenditure(ctx, &expenditurepb.CreateExpenditureRequest{
			Data: &expenditurepb.Expenditure{
				Name:                  r.FormValue("name"),
				ExpenditureType:       r.FormValue("expenditure_type"),
				ExpenditureCategoryId: strPtr(r.FormValue("expenditure_category_id")),
				SupplierId:            strPtr(r.FormValue("supplier_id")),
				ExpenditureDateString: strPtr(r.FormValue("expenditure_date_string")),
				TotalAmount:           parseAmount(r.FormValue("total_amount")),
				Currency:              r.FormValue("currency"),
				Status:                r.FormValue("status"),
				ReferenceNumber:       strPtr(r.FormValue("reference_number")),
				Notes:                 strPtr(r.FormValue("notes")),
				PaymentTermId:         strPtr(r.FormValue("payment_term_id")),
				PurchaseOrderId:       strPtr(r.FormValue("purchase_order_id")),
			},
		})
		if err != nil {
			log.Printf("Failed to create expense: %v", err)
			return centymo.HTMXError(err.Error())
		}

		newID := ""
		if respData := resp.GetData(); len(respData) > 0 {
			newID = respData[0].GetId()
		}
		if newID != "" {
			return view.ViewResult{
				StatusCode: http.StatusOK,
				Headers: map[string]string{
					"HX-Trigger":  `{"formSuccess":true}`,
					"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", newID),
				},
			}
		}

		return centymo.HTMXSuccess("expenses-table")
	})
}

// NewEditAction creates the expense edit action (GET = pre-filled form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expenditure", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadExpenditure(ctx, &expenditurepb.ReadExpenditureRequest{
				Data: &expenditurepb.Expenditure{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read expense %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			readData := readResp.GetData()
			if len(readData) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := readData[0]

			paymentTerms := loadPaymentTerms(ctx, deps)
			selectedPaymentTermID := record.GetPaymentTermId()
			return view.OK("expense-drawer-form", &FormData{
				FormAction:            route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                true,
				ID:                    id,
				Name:                  record.GetName(),
				ExpenditureType:       record.GetExpenditureType(),
				ExpenditureCategoryID: record.GetExpenditureCategoryId(),
				SupplierID:            record.GetSupplierId(),
				Date:                  record.GetExpenditureDateString(),
				TotalAmount:           strconv.FormatFloat(float64(record.GetTotalAmount())/100.0, 'f', 2, 64),
				Currency:              record.GetCurrency(),
				Status:                record.GetStatus(),
				ReferenceNumber:       record.GetReferenceNumber(),
				Notes:                 record.GetNotes(),
				Categories:            loadCategoryOptions(ctx, deps.ListExpenditureCategories, record.GetExpenditureCategoryId()),
				Suppliers:             loadSupplierOptions(ctx, deps.ListSuppliers, record.GetSupplierId()),
				PaymentTerms:          paymentTerms,
				SelectedPaymentTermID: selectedPaymentTermID,
				PurchaseOrders:        loadPurchaseOrderOptions(ctx, deps.ListPurchaseOrders),
				PurchaseOrderID:       record.GetPurchaseOrderId(),
				Labels:                formLabels(deps.Labels),
				CommonLabels:          nil, // injected by ViewAdapter
			})
		}

		// POST — update expense
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		_, err := deps.UpdateExpenditure(ctx, &expenditurepb.UpdateExpenditureRequest{
			Data: &expenditurepb.Expenditure{
				Id:                    id,
				Name:                  r.FormValue("name"),
				ExpenditureType:       r.FormValue("expenditure_type"),
				ExpenditureCategoryId: strPtr(r.FormValue("expenditure_category_id")),
				SupplierId:            strPtr(r.FormValue("supplier_id")),
				ExpenditureDateString: strPtr(r.FormValue("expenditure_date_string")),
				TotalAmount:           parseAmount(r.FormValue("total_amount")),
				Currency:              r.FormValue("currency"),
				Status:                r.FormValue("status"),
				ReferenceNumber:       strPtr(r.FormValue("reference_number")),
				Notes:                 strPtr(r.FormValue("notes")),
				PaymentTermId:         strPtr(r.FormValue("payment_term_id")),
				PurchaseOrderId:       strPtr(r.FormValue("purchase_order_id")),
			},
		})
		if err != nil {
			log.Printf("Failed to update expense %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", id),
			},
		}
	})
}

// NewDeleteAction creates the expense delete action (POST only).
// The row ID comes via query param (?id=xxx) or form field.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expenditure", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeleteExpenditure(ctx, &expenditurepb.DeleteExpenditureRequest{
			Data: &expenditurepb.Expenditure{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete expense %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("expenses-table")
	})
}

// NewSetStatusAction creates the expense status update action (POST only).
// Expects query params: ?id={expenseId}&status={pending|approved|paid|cancelled}
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expenditure", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		validStatuses := map[string]bool{
			"pending":   true,
			"approved":  true,
			"paid":      true,
			"cancelled": true,
		}
		if !validStatuses[targetStatus] {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		_, err := deps.UpdateExpenditure(ctx, &expenditurepb.UpdateExpenditureRequest{
			Data: &expenditurepb.Expenditure{Id: id, Status: targetStatus},
		})
		if err != nil {
			log.Printf("Failed to update expense status %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("expenses-table")
	})
}
