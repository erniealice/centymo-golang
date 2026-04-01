package action

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	expenditurecategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_category"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
)

// CategoryOption represents a selectable category in the form dropdown.
type CategoryOption struct {
	Value    string
	Label    string
	Selected bool
}

// SupplierOption represents a selectable supplier in the form dropdown.
type SupplierOption struct {
	Value    string
	Label    string
	Selected bool
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
	CurrencyPlaceholder string
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
	Categories            []CategoryOption
	Suppliers             []SupplierOption
	Labels                FormLabels
	CommonLabels          any
}

// Deps holds dependencies for expense action handlers.
type Deps struct {
	Routes centymo.ExpenditureRoutes
	Labels centymo.ExpenditureLabels

	// Typed expenditure CRUD operations
	CreateExpenditure func(ctx context.Context, req *expenditurepb.CreateExpenditureRequest) (*expenditurepb.CreateExpenditureResponse, error)
	ReadExpenditure   func(ctx context.Context, req *expenditurepb.ReadExpenditureRequest) (*expenditurepb.ReadExpenditureResponse, error)
	UpdateExpenditure func(ctx context.Context, req *expenditurepb.UpdateExpenditureRequest) (*expenditurepb.UpdateExpenditureResponse, error)
	DeleteExpenditure func(ctx context.Context, req *expenditurepb.DeleteExpenditureRequest) (*expenditurepb.DeleteExpenditureResponse, error)

	// Category listing (optional — gracefully degrades to empty list if nil)
	ListExpenditureCategories func(ctx context.Context, req *expenditurecategorypb.ListExpenditureCategoriesRequest) (*expenditurecategorypb.ListExpenditureCategoriesResponse, error)

	// Supplier listing (optional — gracefully degrades to empty list if nil)
	ListSuppliers func(ctx context.Context, req *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error)
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
	}
}

// loadCategoryOptions loads expenditure categories for the dropdown, pre-selecting selectedID.
func loadCategoryOptions(
	ctx context.Context,
	listFn func(ctx context.Context, req *expenditurecategorypb.ListExpenditureCategoriesRequest) (*expenditurecategorypb.ListExpenditureCategoriesResponse, error),
	selectedID string,
) []CategoryOption {
	if listFn == nil {
		return nil
	}
	resp, err := listFn(ctx, &expenditurecategorypb.ListExpenditureCategoriesRequest{})
	if err != nil {
		log.Printf("Failed to list expenditure categories: %v", err)
		return nil
	}
	var opts []CategoryOption
	for _, cat := range resp.GetData() {
		if !cat.GetActive() {
			continue
		}
		opts = append(opts, CategoryOption{
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
) []SupplierOption {
	if listFn == nil {
		return nil
	}
	resp, err := listFn(ctx, &supplierpb.ListSuppliersRequest{})
	if err != nil {
		log.Printf("Failed to list suppliers: %v", err)
		return nil
	}
	var opts []SupplierOption
	for _, s := range resp.GetData() {
		if !s.GetActive() {
			continue
		}
		label := s.GetCompanyName()
		if label == "" {
			label = s.GetId()
		}
		opts = append(opts, SupplierOption{
			Value:    s.GetId(),
			Label:    label,
			Selected: s.GetId() == selectedID,
		})
	}
	return opts
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}

// parseAmount converts a form string amount to float64.
func parseAmount(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

// NewAddAction creates the expense add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expenditure", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("expense-drawer-form", &FormData{
				FormAction:      deps.Routes.AddURL,
				ExpenditureType: "expense",
				Currency:        "PHP",
				Status:          "pending",
				Categories:      loadCategoryOptions(ctx, deps.ListExpenditureCategories, ""),
				Suppliers:       loadSupplierOptions(ctx, deps.ListSuppliers, ""),
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

			return view.OK("expense-drawer-form", &FormData{
				FormAction:            route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:                true,
				ID:                    id,
				Name:                  record.GetName(),
				ExpenditureType:       record.GetExpenditureType(),
				ExpenditureCategoryID: record.GetExpenditureCategoryId(),
				SupplierID:            record.GetSupplierId(),
				Date:                  record.GetExpenditureDateString(),
				TotalAmount:           strconv.FormatFloat(record.GetTotalAmount(), 'f', 2, 64),
				Currency:              record.GetCurrency(),
				Status:                record.GetStatus(),
				ReferenceNumber:       record.GetReferenceNumber(),
				Notes:                 record.GetNotes(),
				Categories:            loadCategoryOptions(ctx, deps.ListExpenditureCategories, record.GetExpenditureCategoryId()),
				Suppliers:             loadSupplierOptions(ctx, deps.ListSuppliers, record.GetSupplierId()),
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
