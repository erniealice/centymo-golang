package action

import (
	"context"
	"fmt"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"
)

// FormData is the template data for the disbursement drawer form.
type FormData struct {
	FormAction      string
	IsEdit          bool
	ID              string
	ReferenceNumber string
	Payee           string
	Amount          string
	Currency        string
	Method          string
	Date            string
	ApprovedBy      string
	ApprovedRole    string
	Notes           string
	DisbursementType string
	ExpenditureID   string
	Status          string
	Labels          centymo.DisbursementFormLabels
	CommonLabels    any
}

// Deps holds dependencies for disbursement action handlers.
type Deps struct {
	Routes centymo.DisbursementRoutes
	DB     centymo.DataSource
	Labels centymo.DisbursementLabels
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
				Labels:       deps.Labels.Form,
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create disbursement
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		data := map[string]any{
			"reference_number":          r.FormValue("reference_number"),
			"payee":                     r.FormValue("payee"),
			"amount":                    r.FormValue("amount"),
			"currency":                  r.FormValue("currency"),
			"disbursement_method":       r.FormValue("disbursement_method"),
			"disbursement_date_string":  r.FormValue("disbursement_date_string"),
			"approved_by":               r.FormValue("approved_by"),
			"approved_role":             r.FormValue("approved_role"),
			"notes":                     r.FormValue("notes"),
			"disbursement_type":         r.FormValue("disbursement_type"),
			"expenditure_id":            r.FormValue("expenditure_id"),
			"status":                    r.FormValue("status"),
		}

		created, err := deps.DB.Create(ctx, "disbursement", data)
		if err != nil {
			log.Printf("Failed to create disbursement: %v", err)
			return centymo.HTMXError(err.Error())
		}

		newID, _ := created["id"].(string)
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
			record, err := deps.DB.Read(ctx, "disbursement", id)
			if err != nil {
				log.Printf("Failed to read disbursement %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}

			refNumber, _ := record["reference_number"].(string)
			payee, _ := record["payee"].(string)
			amount, _ := record["amount"].(string)
			currency, _ := record["currency"].(string)
			method, _ := record["disbursement_method"].(string)
			date, _ := record["disbursement_date_string"].(string)
			approvedBy, _ := record["approved_by"].(string)
			approvedRole, _ := record["approved_role"].(string)
			notes, _ := record["notes"].(string)
			disbursementType, _ := record["disbursement_type"].(string)
			expenditureID, _ := record["expenditure_id"].(string)
			status, _ := record["status"].(string)

			return view.OK("disbursement-drawer-form", &FormData{
				FormAction:       route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:           true,
				ID:               id,
				ReferenceNumber:  refNumber,
				Payee:            payee,
				Amount:           amount,
				Currency:         currency,
				Method:           method,
				Date:             date,
				ApprovedBy:       approvedBy,
				ApprovedRole:     approvedRole,
				Notes:            notes,
				DisbursementType: disbursementType,
				ExpenditureID:    expenditureID,
				Status:           status,
				Labels:           deps.Labels.Form,
				CommonLabels:     nil, // injected by ViewAdapter
			})
		}

		// POST — update disbursement
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		data := map[string]any{
			"reference_number":          r.FormValue("reference_number"),
			"payee":                     r.FormValue("payee"),
			"amount":                    r.FormValue("amount"),
			"currency":                  r.FormValue("currency"),
			"disbursement_method":       r.FormValue("disbursement_method"),
			"disbursement_date_string":  r.FormValue("disbursement_date_string"),
			"approved_by":               r.FormValue("approved_by"),
			"approved_role":             r.FormValue("approved_role"),
			"notes":                     r.FormValue("notes"),
			"disbursement_type":         r.FormValue("disbursement_type"),
			"expenditure_id":            r.FormValue("expenditure_id"),
			"status":                    r.FormValue("status"),
		}

		_, err := deps.DB.Update(ctx, "disbursement", id, data)
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

		err := deps.DB.Delete(ctx, "disbursement", id)
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
			err := deps.DB.Delete(ctx, "disbursement", id)
			if err != nil {
				log.Printf("Failed to delete disbursement %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("disbursements-table")
	})
}

// NewSetStatusAction creates the disbursement status update action (POST only).
// Validates state transitions:
//   - draft → pending
//   - pending → approved, cancelled
//   - approved → paid, cancelled
//   - cancelled → draft (reactivate)
//   - paid → (none)
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
		record, err := deps.DB.Read(ctx, "disbursement", id)
		if err != nil {
			log.Printf("Failed to read disbursement %s: %v", id, err)
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}

		currentStatus, _ := record["status"].(string)
		if !isValidTransition(currentStatus, targetStatus) {
			return centymo.HTMXError(fmt.Sprintf(deps.Labels.Errors.InvalidTransition, currentStatus, targetStatus))
		}

		if _, err := deps.DB.Update(ctx, "disbursement", id, map[string]any{"status": targetStatus}); err != nil {
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
			record, err := deps.DB.Read(ctx, "disbursement", id)
			if err != nil {
				log.Printf("Failed to read disbursement %s for bulk status: %v", id, err)
				continue
			}

			currentStatus, _ := record["status"].(string)
			if !isValidTransition(currentStatus, targetStatus) {
				log.Printf("Skipping invalid transition %s→%s for disbursement %s", currentStatus, targetStatus, id)
				continue
			}

			if _, err := deps.DB.Update(ctx, "disbursement", id, map[string]any{"status": targetStatus}); err != nil {
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
