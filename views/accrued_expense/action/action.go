package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/erniealice/centymo-golang/views/accrued_expense/form"
)

// Deps holds dependencies for accrued_expense action handlers.
type Deps struct {
	Routes       centymo.AccruedExpenseRoutes
	Labels       centymo.AccruedExpenseLabels
	CommonLabels pyeza.CommonLabels

	CreateAccruedExpense       func(ctx context.Context, req *accruedexpensepb.CreateAccruedExpenseRequest) (*accruedexpensepb.CreateAccruedExpenseResponse, error)
	ReadAccruedExpense         func(ctx context.Context, req *accruedexpensepb.ReadAccruedExpenseRequest) (*accruedexpensepb.ReadAccruedExpenseResponse, error)
	UpdateAccruedExpense       func(ctx context.Context, req *accruedexpensepb.UpdateAccruedExpenseRequest) (*accruedexpensepb.UpdateAccruedExpenseResponse, error)
	DeleteAccruedExpense       func(ctx context.Context, req *accruedexpensepb.DeleteAccruedExpenseRequest) (*accruedexpensepb.DeleteAccruedExpenseResponse, error)
	SetAccruedExpenseStatus    func(ctx context.Context, id, status string) error

	// Dropdowns
	ListSuppliers         func(ctx context.Context, req *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error)
	ListSupplierContracts func(ctx context.Context, req *suppliercontractpb.ListSupplierContractsRequest) (*suppliercontractpb.ListSupplierContractsResponse, error)
}

// NewAddAction handles GET+POST /action/accrued-expense/add.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyFormData(ctx, deps, l)
			fd.FormAction = deps.Routes.AddURL
			fd.Status = "ACCRUED_EXPENSE_STATUS_OUTSTANDING"
			fd.Active = true
			fd.Currency = "PHP"
			fd.RecognitionDate = time.Now().Format("2006-01-02")
			return view.OK("accrued-expense-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		accrued := parseCentavos(r.FormValue("accrued_amount"))

		req := &accruedexpensepb.CreateAccruedExpenseRequest{
			Data: &accruedexpensepb.AccruedExpense{
				Name:               r.FormValue("name"),
				Description:        optionalString(r.FormValue("description")),
				SupplierContractId: r.FormValue("supplier_contract_id"),
				SupplierId:         optionalString(r.FormValue("supplier_id")),
				RecognitionDate:    parseTimestamp(r.FormValue("recognition_date")),
				PeriodStart:        parseTimestamp(r.FormValue("period_start")),
				PeriodEnd:          parseTimestamp(r.FormValue("period_end")),
				CycleDate:          optionalString(r.FormValue("cycle_date")),
				Currency:           r.FormValue("currency"),
				AccruedAmount:      accrued,
				ExpenseAccountId:   optionalString(r.FormValue("expense_account_id")),
				AccrualAccountId:   optionalString(r.FormValue("accrual_account_id")),
				Notes:              optionalString(r.FormValue("notes")),
				Active:             r.FormValue("active") == "true",
			},
		}

		_, err := deps.CreateAccruedExpense(ctx, req)
		if err != nil {
			log.Printf("CreateAccruedExpense: %v", err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("accrued-expenses-table")
	})
}

// NewEditAction handles GET+POST /action/accrued-expense/edit/{id}.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadAccruedExpense(ctx, &accruedexpensepb.ReadAccruedExpenseRequest{
				Data: &accruedexpensepb.AccruedExpense{Id: id},
			})
			if err != nil {
				return view.Error(fmt.Errorf("failed to read accrual: %w", err))
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.Error(fmt.Errorf("accrual not found"))
			}
			a := data[0]

			fd := buildEmptyFormData(ctx, deps, l)
			fd.FormAction = route.ResolveURL(deps.Routes.EditURL, "id", id)
			fd.IsEdit = true
			fd.ID = id
			fd.Name = a.GetName()
			fd.Description = a.GetDescription()
			fd.SupplierContractID = a.GetSupplierContractId()
			fd.SupplierID = a.GetSupplierId()
			fd.RecognitionDate = formatTimestamp(a.GetRecognitionDate())
			fd.PeriodStart = formatTimestamp(a.GetPeriodStart())
			fd.PeriodEnd = formatTimestamp(a.GetPeriodEnd())
			fd.CycleDate = a.GetCycleDate()
			fd.Currency = a.GetCurrency()
			fd.AccruedAmount = formatCentavos(a.GetAccruedAmount())
			fd.SettledAmount = formatCentavos(a.GetSettledAmount())
			fd.RemainingAmount = formatCentavos(a.GetRemainingAmount())
			fd.Status = a.GetStatus().String()
			fd.ExpenseAccountID = a.GetExpenseAccountId()
			fd.AccrualAccountID = a.GetAccrualAccountId()
			fd.Notes = a.GetNotes()
			fd.Active = a.GetActive()
			return view.OK("accrued-expense-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}
		accrued := parseCentavos(r.FormValue("accrued_amount"))

		req := &accruedexpensepb.UpdateAccruedExpenseRequest{
			Data: &accruedexpensepb.AccruedExpense{
				Id:                 id,
				Name:               r.FormValue("name"),
				Description:        optionalString(r.FormValue("description")),
				SupplierContractId: r.FormValue("supplier_contract_id"),
				SupplierId:         optionalString(r.FormValue("supplier_id")),
				RecognitionDate:    parseTimestamp(r.FormValue("recognition_date")),
				PeriodStart:        parseTimestamp(r.FormValue("period_start")),
				PeriodEnd:          parseTimestamp(r.FormValue("period_end")),
				CycleDate:          optionalString(r.FormValue("cycle_date")),
				Currency:           r.FormValue("currency"),
				AccruedAmount:      accrued,
				ExpenseAccountId:   optionalString(r.FormValue("expense_account_id")),
				AccrualAccountId:   optionalString(r.FormValue("accrual_account_id")),
				Notes:              optionalString(r.FormValue("notes")),
				Active:             r.FormValue("active") == "true",
				// NOTE: settled_amount + remaining_amount are NOT written via this
				// path — single-write boundary is enforced by SettleAccrual / ReverseAccrual.
			},
		}

		_, err := deps.UpdateAccruedExpense(ctx, req)
		if err != nil {
			log.Printf("UpdateAccruedExpense %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("accrued-expenses-table")
	})
}

// NewDeleteAction handles POST /action/accrued-expense/delete.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.FormValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		_, err := deps.DeleteAccruedExpense(ctx, &accruedexpensepb.DeleteAccruedExpenseRequest{
			Data: &accruedexpensepb.AccruedExpense{Id: id},
		})
		if err != nil {
			log.Printf("DeleteAccruedExpense %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}
		return centymo.HTMXSuccess("accrued-expenses-table")
	})
}

// NewSetStatusAction handles POST /action/accrued-expense/set-status.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.URL.Query().Get("id")
		status := viewCtx.Request.URL.Query().Get("status")
		if id == "" || status == "" {
			return view.Error(fmt.Errorf("missing id or status"))
		}
		if deps.SetAccruedExpenseStatus != nil {
			if err := deps.SetAccruedExpenseStatus(ctx, id, status); err != nil {
				return centymo.HTMXError(err.Error())
			}
		}
		return centymo.HTMXSuccess("accrued-expenses-table")
	})
}

// NewBulkSetStatusAction handles POST /action/accrued-expense/bulk-set-status.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}
		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")
		if len(ids) == 0 || targetStatus == "" {
			return view.Error(fmt.Errorf("missing ids or target_status"))
		}
		if deps.SetAccruedExpenseStatus != nil {
			for _, id := range ids {
				if err := deps.SetAccruedExpenseStatus(ctx, id, targetStatus); err != nil {
					log.Printf("BulkSetStatus %s → %s: %v", id, targetStatus, err)
				}
			}
		}
		return centymo.HTMXSuccess("accrued-expenses-table")
	})
}

// --- form helpers ------------------------------------------------------------

func buildEmptyFormData(ctx context.Context, deps *Deps, l centymo.AccruedExpenseLabels) *form.Data {
	fd := &form.Data{
		Labels:       l.Form,
		CommonLabels: deps.CommonLabels,
	}

	// Status options
	fd.StatusOptions = []types.SelectOption{
		{Value: "ACCRUED_EXPENSE_STATUS_OUTSTANDING", Label: l.Form.StatusOutstanding},
		{Value: "ACCRUED_EXPENSE_STATUS_PARTIAL", Label: l.Form.StatusPartial},
		{Value: "ACCRUED_EXPENSE_STATUS_SETTLED", Label: l.Form.StatusSettled},
		{Value: "ACCRUED_EXPENSE_STATUS_REVERSED", Label: l.Form.StatusReversed},
	}

	if deps.ListSuppliers != nil {
		resp, err := deps.ListSuppliers(ctx, &supplierpb.ListSuppliersRequest{})
		if err == nil {
			for _, s := range resp.GetData() {
				fd.Suppliers = append(fd.Suppliers, types.SelectOption{
					Value: s.GetId(),
					Label: s.GetName(),
				})
			}
		}
	}

	if deps.ListSupplierContracts != nil {
		resp, err := deps.ListSupplierContracts(ctx, &suppliercontractpb.ListSupplierContractsRequest{})
		if err == nil {
			for _, c := range resp.GetData() {
				fd.SupplierContracts = append(fd.SupplierContracts, types.SelectOption{
					Value: c.GetId(),
					Label: c.GetName(),
				})
			}
		}
	}

	return fd
}

// --- centavo / timestamp helpers --------------------------------------------

func parseCentavos(s string) int64 {
	if s == "" {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f * 100))
}

func formatCentavos(v int64) string {
	if v == 0 {
		return ""
	}
	return strconv.FormatFloat(float64(v)/100.0, 'f', 2, 64)
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseTimestamp(s string) *timestamppb.Timestamp {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return timestamppb.New(t)
}

func formatTimestamp(t *timestamppb.Timestamp) string {
	if t == nil {
		return ""
	}
	return t.AsTime().Format("2006-01-02")
}
