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
	procurementrequestpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/procurement_request"
)

// FormData is the template data for the procurement request drawer form.
type FormData struct {
	FormAction string
	IsEdit     bool
	ID         string

	// Section 1 — Identity
	RequestNumber      string
	RequesterUserID    string
	SupplierID         string
	LocationID         string

	// Section 2 — Financial
	Currency             string
	EstimatedTotalAmount string

	// Section 3 — Timing & Approval
	NeededByDate   string
	ApprovedBy     string
	ApprovedAt     string
	Status         string

	// Section 4 — Others
	Justification string
	Notes         string
	Active        bool

	// Dropdown options
	Suppliers     []types.SelectOption
	StatusOptions []types.SelectOption

	Labels       centymo.ProcurementRequestFormLabels
	CommonLabels pyeza.CommonLabels
}

// Deps holds all dependencies for the procurement request action handlers.
type Deps struct {
	Routes                     centymo.ProcurementRequestRoutes
	Labels                     centymo.ProcurementRequestLabels
	CommonLabels               pyeza.CommonLabels
	CreateProcurementRequest   func(ctx context.Context, req *procurementrequestpb.CreateProcurementRequestRequest) (*procurementrequestpb.CreateProcurementRequestResponse, error)
	ReadProcurementRequest     func(ctx context.Context, req *procurementrequestpb.ReadProcurementRequestRequest) (*procurementrequestpb.ReadProcurementRequestResponse, error)
	UpdateProcurementRequest   func(ctx context.Context, req *procurementrequestpb.UpdateProcurementRequestRequest) (*procurementrequestpb.UpdateProcurementRequestResponse, error)
	DeleteProcurementRequest   func(ctx context.Context, req *procurementrequestpb.DeleteProcurementRequestRequest) (*procurementrequestpb.DeleteProcurementRequestResponse, error)
	SetProcurementRequestStatus func(ctx context.Context, id, status string) error
	ListSuppliers               func(ctx context.Context, req *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error)
}

// NewAddAction handles GET+POST /action/procurement-request/add.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			fd := buildEmptyFormData(ctx, deps, l)
			fd.FormAction = deps.Routes.AddURL
			fd.Status = "draft"
			fd.Active = true
			fd.Currency = "PHP"
			return view.OK("procurement-request-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		estimated := parseCentavos(r.FormValue("estimated_total_amount"))

		req := &procurementrequestpb.CreateProcurementRequestRequest{
			Data: &procurementrequestpb.ProcurementRequest{
				RequestNumber:        r.FormValue("request_number"),
				RequesterUserId:      r.FormValue("requester_user_id"),
				SupplierId:           optionalString(r.FormValue("supplier_id")),
				LocationId:           optionalString(r.FormValue("location_id")),
				Currency:             r.FormValue("currency"),
				EstimatedTotalAmount: estimated,
				NeededByDate:         optionalString(r.FormValue("needed_by_date")),
				Justification:        optionalString(r.FormValue("justification")),
				Notes:                optionalString(r.FormValue("notes")),
				Active:               r.FormValue("active") == "true",
			},
		}

		_, err := deps.CreateProcurementRequest(ctx, req)
		if err != nil {
			log.Printf("CreateProcurementRequest: %v", err)
			return view.Error(fmt.Errorf("failed to create procurement request: %w", err))
		}

		return centymo.HTMXSuccess("procurement-requests-table")
	})
}

// NewEditAction handles GET+POST /action/procurement-request/edit/{id}.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadProcurementRequest(ctx, &procurementrequestpb.ReadProcurementRequestRequest{
				Data: &procurementrequestpb.ProcurementRequest{Id: id},
			})
			if err != nil {
				return view.Error(fmt.Errorf("failed to read procurement request: %w", err))
			}
			data := resp.GetData()
			if len(data) == 0 {
				return view.Error(fmt.Errorf("procurement request not found"))
			}
			pr := data[0]

			fd := buildEmptyFormData(ctx, deps, l)
			fd.FormAction = route.ResolveURL(deps.Routes.EditURL, "id", id)
			fd.IsEdit = true
			fd.ID = id
			fd.RequestNumber = pr.GetRequestNumber()
			fd.RequesterUserID = pr.GetRequesterUserId()
			fd.SupplierID = pr.GetSupplierId()
			fd.LocationID = pr.GetLocationId()
			fd.Currency = pr.GetCurrency()
			fd.EstimatedTotalAmount = formatCentavos(pr.GetEstimatedTotalAmount())
			fd.NeededByDate = pr.GetNeededByDate()
			fd.ApprovedBy = pr.GetApprovedBy()
			fd.Justification = pr.GetJustification()
			fd.Notes = pr.GetNotes()
			fd.Status = pr.GetStatus().String()
			fd.Active = pr.GetActive()
			return view.OK("procurement-request-drawer-form", fd)
		}

		// POST
		r := viewCtx.Request
		if err := r.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}

		estimated := parseCentavos(r.FormValue("estimated_total_amount"))

		req := &procurementrequestpb.UpdateProcurementRequestRequest{
			Data: &procurementrequestpb.ProcurementRequest{
				Id:                   id,
				RequestNumber:        r.FormValue("request_number"),
				RequesterUserId:      r.FormValue("requester_user_id"),
				SupplierId:           optionalString(r.FormValue("supplier_id")),
				LocationId:           optionalString(r.FormValue("location_id")),
				Currency:             r.FormValue("currency"),
				EstimatedTotalAmount: estimated,
				NeededByDate:         optionalString(r.FormValue("needed_by_date")),
				Justification:        optionalString(r.FormValue("justification")),
				Notes:                optionalString(r.FormValue("notes")),
				Active:               r.FormValue("active") == "true",
			},
		}

		_, err := deps.UpdateProcurementRequest(ctx, req)
		if err != nil {
			log.Printf("UpdateProcurementRequest %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to update procurement request: %w", err))
		}

		return centymo.HTMXSuccess("procurement-requests-table")
	})
}

// NewDeleteAction handles POST /action/procurement-request/delete.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		id := viewCtx.Request.FormValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing id"))
		}
		_, err := deps.DeleteProcurementRequest(ctx, &procurementrequestpb.DeleteProcurementRequestRequest{
			Data: &procurementrequestpb.ProcurementRequest{Id: id},
		})
		if err != nil {
			log.Printf("DeleteProcurementRequest %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to delete procurement request: %w", err))
		}
		return centymo.HTMXSuccess("procurement-requests-table")
	})
}

// NewSetStatusAction handles POST /action/procurement-request/set-status.
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
		if deps.SetProcurementRequestStatus != nil {
			if err := deps.SetProcurementRequestStatus(ctx, id, status); err != nil {
				return view.Error(fmt.Errorf("failed to set status: %w", err))
			}
		}
		return centymo.HTMXSuccess("procurement-requests-table")
	})
}

// NewBulkSetStatusAction handles POST /action/procurement-request/bulk-set-status.
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
		if deps.SetProcurementRequestStatus != nil {
			for _, id := range ids {
				if err := deps.SetProcurementRequestStatus(ctx, id, targetStatus); err != nil {
					log.Printf("BulkSetStatus %s → %s: %v", id, targetStatus, err)
				}
			}
		}
		return centymo.HTMXSuccess("procurement-requests-table")
	})
}

// --- helpers -----------------------------------------------------------------

func buildEmptyFormData(ctx context.Context, deps *Deps, l centymo.ProcurementRequestLabels) *FormData {
	fd := &FormData{
		Labels:       l.Form,
		CommonLabels: deps.CommonLabels,
	}

	// Status options matching ProcurementRequestStatus enum
	fd.StatusOptions = []types.SelectOption{
		{Value: "draft", Label: l.Form.StatusDraft},
		{Value: "submitted", Label: l.Form.StatusSubmitted},
		{Value: "pending_approval", Label: l.Form.StatusPendingApproval},
		{Value: "approved", Label: l.Form.StatusApproved},
		{Value: "rejected", Label: l.Form.StatusRejected},
		{Value: "fulfilled", Label: l.Form.StatusFulfilled},
		{Value: "cancelled", Label: l.Form.StatusCancelled},
	}

	// Load supplier options (nullable — supplier may be unknown at request time)
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

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func optionalStrVal(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
