package action

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

	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	suppliercontractpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/supplier_contract"
)

// KindOption is a select option for the contract kind enum.
type KindOption struct {
	Value    string
	Label    string
	Selected bool
}

// FormData is the template data for the supplier contract drawer form.
type FormData struct {
	FormAction string
	IsEdit     bool
	ID         string

	// Section 1 — Company / Identity Details
	Name             string
	ReferenceNumber  string
	Kind             string
	KindOptions      []KindOption
	SupplierID       string
	SupplierName     string

	// Section 2 — Validity & Recurrence
	StartDate         string
	EndDate           string
	BillingCycleValue string
	BillingCycleUnit  string
	AutoRenew         bool
	RenewalNoticeDays string

	// Section 3 — Money & Approval
	Currency        string
	CommittedAmount string
	CycleAmount     string
	PaymentTermID   string
	ApprovedBy      string
	ApprovedDate    string
	RequestedBy     string

	// Section 4 — Categorization
	ExpenditureCategoryID string
	ExpenseAccountID      string
	LocationID            string

	// Section 5 — Others
	Notes  string
	Active bool
	Status string

	// Dropdown options
	Suppliers     []types.SelectOption
	PaymentTerms  []types.SelectOption
	StatusOptions []types.SelectOption

	Labels       centymo.SupplierContractFormLabels
	CommonLabels pyeza.CommonLabels
}

// Deps holds all dependencies for the supplier contract action handlers.
type Deps struct {
	Routes                    centymo.SupplierContractRoutes
	Labels                    centymo.SupplierContractLabels
	CommonLabels              pyeza.CommonLabels
	CreateSupplierContract    func(ctx context.Context, req *suppliercontractpb.CreateSupplierContractRequest) (*suppliercontractpb.CreateSupplierContractResponse, error)
	ReadSupplierContract      func(ctx context.Context, req *suppliercontractpb.ReadSupplierContractRequest) (*suppliercontractpb.ReadSupplierContractResponse, error)
	UpdateSupplierContract    func(ctx context.Context, req *suppliercontractpb.UpdateSupplierContractRequest) (*suppliercontractpb.UpdateSupplierContractResponse, error)
	DeleteSupplierContract    func(ctx context.Context, req *suppliercontractpb.DeleteSupplierContractRequest) (*suppliercontractpb.DeleteSupplierContractResponse, error)
	SetSupplierContractStatus func(ctx context.Context, id, status string) error
	ListSuppliers             func(ctx context.Context, req *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error)
}

// NewAddAction handles GET+POST /action/supplier-contract/add.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyFormData(ctx, deps, l)
			fd.FormAction = deps.Routes.AddURL
			fd.Status = "draft"
			fd.Active = true
			fd.Currency = "PHP"
			return view.OK("supplier-contract-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		cycleValue, _ := strconv.ParseInt(r.FormValue("billing_cycle_value"), 10, 32)
		renewalDays, _ := strconv.ParseInt(r.FormValue("renewal_notice_days"), 10, 32)
		autoRenew := r.FormValue("auto_renew") == "true"

		committed := parseCentavos(r.FormValue("committed_amount"))
		cycle := parseCentavos(r.FormValue("cycle_amount"))

		req := &suppliercontractpb.CreateSupplierContractRequest{
			Data: &suppliercontractpb.SupplierContract{
				Name:                  r.FormValue("name"),
				ReferenceNumber:       optionalString(r.FormValue("reference_number")),
				SupplierId:            r.FormValue("supplier_id"),
				DateTimeStart:         r.FormValue("start_date"),
				DateTimeEnd:           optionalString(r.FormValue("end_date")),
				BillingCycleValue:     optionalInt32(int32(cycleValue)),
				BillingCycleUnit:      optionalString(r.FormValue("billing_cycle_unit")),
				AutoRenew:             autoRenew,
				RenewalNoticeDays:     optionalInt32(int32(renewalDays)),
				Currency:              r.FormValue("currency"),
				CommittedAmount:       optionalInt64(committed),
				CycleAmount:           optionalInt64(cycle),
				PaymentTermId:         optionalString(r.FormValue("payment_term_id")),
				ExpenditureCategoryId: optionalString(r.FormValue("expenditure_category_id")),
				ExpenseAccountId:      optionalString(r.FormValue("expense_account_id")),
				LocationId:            optionalString(r.FormValue("location_id")),
				Notes:                 optionalString(r.FormValue("notes")),
				Active:                r.FormValue("active") == "true",
			},
		}

		_, err := deps.CreateSupplierContract(ctx, req)
		if err != nil {
			log.Printf("CreateSupplierContract: %v", err)
			return view.Error(fmt.Errorf("failed to create supplier contract: %w", err))
		}

		return centymo.HTMXSuccess("supplier-contracts-table")
	})
}

// NewEditAction handles GET+POST /action/supplier-contract/edit/{id}.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadSupplierContract(ctx, &suppliercontractpb.ReadSupplierContractRequest{
				Data: &suppliercontractpb.SupplierContract{Id: id},
			})
			if err != nil {
				return view.Error(fmt.Errorf("failed to read supplier contract: %w", err))
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.Error(fmt.Errorf("supplier contract not found"))
			}
			c := data[0]

			fd := buildEmptyFormData(ctx, deps, l)
			fd.FormAction = route.ResolveURL(deps.Routes.EditURL, "id", id)
			fd.IsEdit = true
			fd.ID = id
			fd.Name = c.GetName()
			fd.ReferenceNumber = c.GetReferenceNumber()
			fd.Kind = c.GetKind().String()
			fd.SupplierID = c.GetSupplierId()
			fd.StartDate = c.GetDateTimeStart()
			fd.EndDate = c.GetDateTimeEnd()
			fd.AutoRenew = c.GetAutoRenew()
			fd.RenewalNoticeDays = strconv.FormatInt(int64(c.GetRenewalNoticeDays()), 10)
			fd.BillingCycleValue = strconv.FormatInt(int64(c.GetBillingCycleValue()), 10)
			fd.BillingCycleUnit = c.GetBillingCycleUnit()
			fd.Currency = c.GetCurrency()
			fd.CommittedAmount = formatCentavos(c.GetCommittedAmount())
			fd.CycleAmount = formatCentavos(c.GetCycleAmount())
			fd.PaymentTermID = c.GetPaymentTermId()
			fd.ExpenditureCategoryID = c.GetExpenditureCategoryId()
			fd.ExpenseAccountID = c.GetExpenseAccountId()
			fd.LocationID = c.GetLocationId()
			fd.Notes = c.GetNotes()
			fd.Status = c.GetStatus().String()
			fd.Active = c.GetActive()
			return view.OK("supplier-contract-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		cycleValue, _ := strconv.ParseInt(r.FormValue("billing_cycle_value"), 10, 32)
		renewalDays, _ := strconv.ParseInt(r.FormValue("renewal_notice_days"), 10, 32)
		autoRenew := r.FormValue("auto_renew") == "true"
		committed := parseCentavos(r.FormValue("committed_amount"))
		cycle := parseCentavos(r.FormValue("cycle_amount"))

		req := &suppliercontractpb.UpdateSupplierContractRequest{
			Data: &suppliercontractpb.SupplierContract{
				Id:                    id,
				Name:                  r.FormValue("name"),
				ReferenceNumber:       optionalString(r.FormValue("reference_number")),
				SupplierId:            r.FormValue("supplier_id"),
				DateTimeStart:         r.FormValue("start_date"),
				DateTimeEnd:           optionalString(r.FormValue("end_date")),
				BillingCycleValue:     optionalInt32(int32(cycleValue)),
				BillingCycleUnit:      optionalString(r.FormValue("billing_cycle_unit")),
				AutoRenew:             autoRenew,
				RenewalNoticeDays:     optionalInt32(int32(renewalDays)),
				Currency:              r.FormValue("currency"),
				CommittedAmount:       optionalInt64(committed),
				CycleAmount:           optionalInt64(cycle),
				PaymentTermId:         optionalString(r.FormValue("payment_term_id")),
				ExpenditureCategoryId: optionalString(r.FormValue("expenditure_category_id")),
				ExpenseAccountId:      optionalString(r.FormValue("expense_account_id")),
				LocationId:            optionalString(r.FormValue("location_id")),
				Notes:                 optionalString(r.FormValue("notes")),
				Active:                r.FormValue("active") == "true",
			},
		}

		_, err := deps.UpdateSupplierContract(ctx, req)
		if err != nil {
			log.Printf("UpdateSupplierContract %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to update supplier contract: %w", err))
		}

		return centymo.HTMXSuccess("supplier-contracts-table")
	})
}

// NewDeleteAction handles POST /action/supplier-contract/delete.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.FormValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		_, err := deps.DeleteSupplierContract(ctx, &suppliercontractpb.DeleteSupplierContractRequest{
			Data: &suppliercontractpb.SupplierContract{Id: id},
		})
		if err != nil {
			log.Printf("DeleteSupplierContract %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to delete supplier contract: %w", err))
		}
		return centymo.HTMXSuccess("supplier-contracts-table")
	})
}

// NewSetStatusAction handles POST /action/supplier-contract/set-status.
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
		if deps.SetSupplierContractStatus != nil {
			if err := deps.SetSupplierContractStatus(ctx, id, status); err != nil {
				return view.Error(fmt.Errorf("failed to set status: %w", err))
			}
		}
		return centymo.HTMXSuccess("supplier-contracts-table")
	})
}

// NewBulkSetStatusAction handles POST /action/supplier-contract/bulk-set-status.
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
		if deps.SetSupplierContractStatus != nil {
			for _, id := range ids {
				if err := deps.SetSupplierContractStatus(ctx, id, targetStatus); err != nil {
					log.Printf("BulkSetStatus %s → %s: %v", id, targetStatus, err)
				}
			}
		}
		return centymo.HTMXSuccess("supplier-contracts-table")
	})
}

// --- form helpers ------------------------------------------------------------

func buildEmptyFormData(ctx context.Context, deps *Deps, l centymo.SupplierContractLabels) *FormData {
	fd := &FormData{
		Labels:       l.Form,
		CommonLabels: deps.CommonLabels,
	}

	// Kind options — use proto enum string representation
	fd.KindOptions = []KindOption{
		{Value: "SUPPLIER_CONTRACT_KIND_SUBSCRIPTION", Label: l.Form.KindSubscription},
		{Value: "SUPPLIER_CONTRACT_KIND_RETAINER", Label: l.Form.KindRetainer},
		{Value: "SUPPLIER_CONTRACT_KIND_LEASE", Label: l.Form.KindLease},
		{Value: "SUPPLIER_CONTRACT_KIND_UTILITY", Label: l.Form.KindUtility},
		{Value: "SUPPLIER_CONTRACT_KIND_FRAMEWORK", Label: l.Form.KindFramework},
		{Value: "SUPPLIER_CONTRACT_KIND_BLANKET", Label: l.Form.KindBlanket},
		{Value: "SUPPLIER_CONTRACT_KIND_ONE_TIME", Label: l.Form.KindOneTime},
		{Value: "SUPPLIER_CONTRACT_KIND_OTHER", Label: l.Form.KindOther},
	}

	// Status options
	fd.StatusOptions = []types.SelectOption{
		{Value: "SUPPLIER_CONTRACT_STATUS_DRAFT", Label: l.Form.StatusDraft},
		{Value: "SUPPLIER_CONTRACT_STATUS_REQUESTED", Label: l.Form.StatusRequested},
		{Value: "SUPPLIER_CONTRACT_STATUS_PENDING_APPROVAL", Label: l.Form.StatusPendingApproval},
		{Value: "SUPPLIER_CONTRACT_STATUS_APPROVED", Label: l.Form.StatusApproved},
		{Value: "SUPPLIER_CONTRACT_STATUS_ACTIVE", Label: l.Form.StatusActive},
		{Value: "SUPPLIER_CONTRACT_STATUS_EXPIRING", Label: l.Form.StatusExpiring},
		{Value: "SUPPLIER_CONTRACT_STATUS_SUSPENDED", Label: l.Form.StatusSuspended},
		{Value: "SUPPLIER_CONTRACT_STATUS_EXPIRED", Label: l.Form.StatusExpired},
		{Value: "SUPPLIER_CONTRACT_STATUS_TERMINATED", Label: l.Form.StatusTerminated},
		{Value: "SUPPLIER_CONTRACT_STATUS_REJECTED", Label: l.Form.StatusRejected},
	}

	// Load supplier options
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

	return fd
}

// --- centavo helpers ---------------------------------------------------------

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

func formatOptInt64(p *int64) string {
	if p == nil {
		return ""
	}
	return formatCentavos(*p)
}

func formatOptInt32(p *int32) string {
	if p == nil {
		return ""
	}
	return strconv.FormatInt(int64(*p), 10)
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func optionalInt64(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func optionalInt32(v int32) *int32 {
	if v == 0 {
		return nil
	}
	return &v
}

func optionalStrVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
