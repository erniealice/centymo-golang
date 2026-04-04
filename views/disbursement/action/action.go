package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	centymo "github.com/erniealice/centymo-golang"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	disbursementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement"
)

// ExpenditureOption is a minimal struct for rendering expenditure (bill) options in the form.
type ExpenditureOption struct {
	Id     string
	Name   string
	Amount string
}

// FormData is the template data for the disbursement drawer form.
type FormData struct {
	FormAction       string
	IsEdit           bool
	ID               string
	ReferenceNumber  string
	Payee            string
	Amount           string
	Currency         string
	Method           string
	Date             string
	ApprovedBy       string
	ApprovedRole     string
	Notes            string
	DisbursementType string
	ExpenditureID    string
	Expenditures     []*ExpenditureOption
	Status           string
	Labels           centymo.DisbursementFormLabels
	CommonLabels     any
}

// Deps holds dependencies for disbursement action handlers.
type Deps struct {
	Routes             centymo.DisbursementRoutes
	Labels             centymo.DisbursementLabels
	CreateDisbursement func(ctx context.Context, req *disbursementpb.CreateDisbursementRequest) (*disbursementpb.CreateDisbursementResponse, error)
	ReadDisbursement   func(ctx context.Context, req *disbursementpb.ReadDisbursementRequest) (*disbursementpb.ReadDisbursementResponse, error)
	UpdateDisbursement func(ctx context.Context, req *disbursementpb.UpdateDisbursementRequest) (*disbursementpb.UpdateDisbursementResponse, error)
	DeleteDisbursement func(ctx context.Context, req *disbursementpb.DeleteDisbursementRequest) (*disbursementpb.DeleteDisbursementResponse, error)

	// Expenditure (bill) listing (optional — gracefully degrades to empty list if nil)
	ListExpenditures func(ctx context.Context, req *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error)
}

// loadExpenditureOptions loads unpaid expenditures (bills) for the dropdown.
// Only expenditures with status "pending" or "approved" are included.
func loadExpenditureOptions(
	ctx context.Context,
	listFn func(ctx context.Context, req *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error),
) []*ExpenditureOption {
	if listFn == nil {
		return nil
	}
	resp, err := listFn(ctx, &expenditurepb.ListExpendituresRequest{})
	if err != nil {
		log.Printf("Failed to list expenditures: %v", err)
		return nil
	}
	var opts []*ExpenditureOption
	for _, e := range resp.GetData() {
		status := e.GetStatus()
		if status != "pending" && status != "approved" {
			continue
		}
		amount := fmt.Sprintf("%.2f", float64(e.GetTotalAmount())/100.0)
		opts = append(opts, &ExpenditureOption{
			Id:     e.GetId(),
			Name:   e.GetName(),
			Amount: e.GetCurrency() + " " + amount,
		})
	}
	return opts
}

// parseAmount converts a form string amount (decimal) to int64 centavos.
func parseAmount(s string) int64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f * 100))
}

// NewAddAction creates the disbursement add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("disbursement-drawer-form", &FormData{
				FormAction:   deps.Routes.AddURL,
				Currency:     "PHP",
				Status:       "draft",
				Expenditures: loadExpenditureOptions(ctx, deps.ListExpenditures),
				Labels:       deps.Labels.Form,
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create disbursement
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		resp, err := deps.CreateDisbursement(ctx, &disbursementpb.CreateDisbursementRequest{
			Data: &disbursementpb.Disbursement{
				ReferenceNumber:      r.FormValue("reference_number"),
				Name:                 r.FormValue("payee"),
				Amount:               parseAmount(r.FormValue("amount")),
				Currency:             r.FormValue("currency"),
				DisbursementMethodId: r.FormValue("disbursement_method"),
				ApprovedBy:           r.FormValue("approved_by"),
				DisbursementType:     r.FormValue("disbursement_type"),
				ExpenditureId:        r.FormValue("expenditure_id"),
				Status:               r.FormValue("status"),
			},
		})
		if err != nil {
			log.Printf("Failed to create disbursement: %v", err)
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

		return centymo.HTMXSuccess("disbursements-table")
	})
}

// NewEditAction creates the disbursement edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadDisbursement(ctx, &disbursementpb.ReadDisbursementRequest{
				Data: &disbursementpb.Disbursement{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read disbursement %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			readData := readResp.GetData()
			if len(readData) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := readData[0]

			return view.OK("disbursement-drawer-form", &FormData{
				FormAction:       route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:           true,
				ID:               id,
				ReferenceNumber:  record.GetReferenceNumber(),
				Payee:            record.GetName(),
				Amount:           fmt.Sprintf("%.2f", float64(record.GetAmount())/100.0),
				Currency:         record.GetCurrency(),
				Method:           record.GetDisbursementMethodId(),
				Date:             record.GetDateCreatedString(),
				ApprovedBy:       record.GetApprovedBy(),
				ApprovedRole:     "",
				Notes:            "",
				DisbursementType: record.GetDisbursementType(),
				ExpenditureID:    record.GetExpenditureId(),
				Expenditures:     loadExpenditureOptions(ctx, deps.ListExpenditures),
				Status:           record.GetStatus(),
				Labels:           deps.Labels.Form,
				CommonLabels:     nil, // injected by ViewAdapter
			})
		}

		// POST — update disbursement
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		_, err := deps.UpdateDisbursement(ctx, &disbursementpb.UpdateDisbursementRequest{
			Data: &disbursementpb.Disbursement{
				Id:                   id,
				ReferenceNumber:      r.FormValue("reference_number"),
				Name:                 r.FormValue("payee"),
				Amount:               parseAmount(r.FormValue("amount")),
				Currency:             r.FormValue("currency"),
				DisbursementMethodId: r.FormValue("disbursement_method"),
				ApprovedBy:           r.FormValue("approved_by"),
				DisbursementType:     r.FormValue("disbursement_type"),
				ExpenditureId:        r.FormValue("expenditure_id"),
				Status:               r.FormValue("status"),
			},
		})
		if err != nil {
			log.Printf("Failed to update disbursement %s: %v", id, err)
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

// NewDeleteAction creates the disbursement delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement", "delete") {
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

		_, err := deps.DeleteDisbursement(ctx, &disbursementpb.DeleteDisbursementRequest{
			Data: &disbursementpb.Disbursement{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete disbursement %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("disbursements-table")
	})
}

// NewBulkDeleteAction creates the disbursement bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}

		for _, id := range ids {
			_, err := deps.DeleteDisbursement(ctx, &disbursementpb.DeleteDisbursementRequest{
				Data: &disbursementpb.Disbursement{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete disbursement %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("disbursements-table")
	})
}

// NewSetStatusAction creates the disbursement status update action (POST only).
// Validates state transitions:
//   - draft -> pending
//   - pending -> approved, cancelled
//   - approved -> paid, cancelled
//   - cancelled -> draft (reactivate)
//   - paid -> (none)
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("target_status")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		if !isValidStatus(targetStatus) {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		// Read current record to validate transition
		readResp, err := deps.ReadDisbursement(ctx, &disbursementpb.ReadDisbursementRequest{
			Data: &disbursementpb.Disbursement{Id: id},
		})
		if err != nil {
			log.Printf("Failed to read disbursement %s: %v", id, err)
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}
		readData := readResp.GetData()
		if len(readData) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}

		currentStatus := readData[0].GetStatus()
		if !isValidTransition(currentStatus, targetStatus) {
			return centymo.HTMXError(fmt.Sprintf(deps.Labels.Errors.InvalidTransition, currentStatus, targetStatus))
		}

		if _, err := deps.UpdateDisbursement(ctx, &disbursementpb.UpdateDisbursementRequest{
			Data: &disbursementpb.Disbursement{Id: id, Status: targetStatus},
		}); err != nil {
			log.Printf("Failed to update disbursement status %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("disbursements-table")
	})
}

// NewBulkSetStatusAction creates the disbursement bulk status update action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("disbursement", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}

		if !isValidStatus(targetStatus) {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		for _, id := range ids {
			readResp, err := deps.ReadDisbursement(ctx, &disbursementpb.ReadDisbursementRequest{
				Data: &disbursementpb.Disbursement{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read disbursement %s for bulk status: %v", id, err)
				continue
			}
			readData := readResp.GetData()
			if len(readData) == 0 {
				log.Printf("Disbursement %s not found for bulk status", id)
				continue
			}

			currentStatus := readData[0].GetStatus()
			if !isValidTransition(currentStatus, targetStatus) {
				log.Printf("Skipping invalid transition %s->%s for disbursement %s", currentStatus, targetStatus, id)
				continue
			}

			if _, err := deps.UpdateDisbursement(ctx, &disbursementpb.UpdateDisbursementRequest{
				Data: &disbursementpb.Disbursement{Id: id, Status: targetStatus},
			}); err != nil {
				log.Printf("Failed to update disbursement status %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("disbursements-table")
	})
}

func isValidStatus(s string) bool {
	switch s {
	case "draft", "pending", "approved", "paid", "cancelled", "overdue":
		return true
	}
	return false
}

func isValidTransition(from, to string) bool {
	switch from {
	case "draft":
		return to == "pending"
	case "pending":
		return to == "approved" || to == "cancelled"
	case "approved":
		return to == "paid" || to == "cancelled"
	case "overdue":
		return to == "paid" || to == "cancelled"
	case "cancelled":
		return to == "draft"
	case "paid":
		return false
	}
	return false
}
