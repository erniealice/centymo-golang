package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/centymo-golang"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"
)

// FormLabels holds i18n labels for the collection drawer form template.
type FormLabels struct {
	Customer             string
	Date                 string
	Amount               string
	Currency             string
	Reference            string
	ReferencePlaceholder string
	PaymentMethod        string
	Status               string
	Notes                string
	NotesPlaceholder     string
}

// FormData is the template data for the collection drawer form.
type FormData struct {
	FormAction       string
	IsEdit           bool
	ID               string
	Customer         string
	ReferenceNumber  string
	Amount           string
	Currency         string
	CollectionMethod string
	Date             string
	ReceivedBy       string
	ReceivedRole     string
	Notes            string
	CollectionType   string
	Status           string
	Labels           FormLabels
	CommonLabels     any
}

// Deps holds dependencies for collection action handlers.
type Deps struct {
	Routes centymo.CollectionRoutes
	DB     centymo.DataSource
	Labels centymo.CollectionLabels
}

func formLabels(l centymo.CollectionFormLabels) FormLabels {
	return FormLabels{
		Customer:             l.Customer,
		Date:                 l.Date,
		Amount:               l.Amount,
		Currency:             l.Currency,
		Reference:            l.Reference,
		ReferencePlaceholder: l.ReferencePlaceholder,
		PaymentMethod:        l.PaymentMethod,
		Status:               l.Status,
		Notes:                l.Notes,
		NotesPlaceholder:     l.NotesPlaceholder,
	}
}

// NewAddAction creates the collection add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("collection-drawer-form", &FormData{
				FormAction:   deps.Routes.AddURL,
				Currency:     "PHP",
				Status:       "pending",
				Labels:       formLabels(centymo.DefaultCollectionLabels().Form),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create collection
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		data := map[string]any{
			"reference_number":  r.FormValue("reference_number"),
			"customer":          r.FormValue("customer"),
			"amount":            r.FormValue("amount"),
			"currency":          r.FormValue("currency"),
			"collection_method": r.FormValue("collection_method"),
			"date":              r.FormValue("date"),
			"received_by":       r.FormValue("received_by"),
			"received_role":     r.FormValue("received_role"),
			"notes":             r.FormValue("notes"),
			"collection_type":   r.FormValue("collection_type"),
			"status":            r.FormValue("status"),
		}

		created, err := deps.DB.Create(ctx, "collection", data)
		if err != nil {
			log.Printf("Failed to create collection: %v", err)
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

		return centymo.HTMXSuccess("collections-table")
	})
}

// NewEditAction creates the collection edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			record, err := deps.DB.Read(ctx, "collection", id)
			if err != nil {
				log.Printf("Failed to read collection %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}

			customer, _ := record["customer"].(string)
			refNumber, _ := record["reference_number"].(string)
			amount, _ := record["amount"].(string)
			currency, _ := record["currency"].(string)
			method, _ := record["collection_method"].(string)
			date, _ := record["date"].(string)
			receivedBy, _ := record["received_by"].(string)
			receivedRole, _ := record["received_role"].(string)
			notes, _ := record["notes"].(string)
			collectionType, _ := record["collection_type"].(string)
			status, _ := record["status"].(string)

			return view.OK("collection-drawer-form", &FormData{
				FormAction:       route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:           true,
				ID:               id,
				Customer:         customer,
				ReferenceNumber:  refNumber,
				Amount:           amount,
				Currency:         currency,
				CollectionMethod: method,
				Date:             date,
				ReceivedBy:       receivedBy,
				ReceivedRole:     receivedRole,
				Notes:            notes,
				CollectionType:   collectionType,
				Status:           status,
				Labels:           formLabels(centymo.DefaultCollectionLabels().Form),
				CommonLabels:     nil, // injected by ViewAdapter
			})
		}

		// POST — update collection
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		data := map[string]any{
			"reference_number":  r.FormValue("reference_number"),
			"customer":          r.FormValue("customer"),
			"amount":            r.FormValue("amount"),
			"currency":          r.FormValue("currency"),
			"collection_method": r.FormValue("collection_method"),
			"date":              r.FormValue("date"),
			"received_by":       r.FormValue("received_by"),
			"received_role":     r.FormValue("received_role"),
			"notes":             r.FormValue("notes"),
			"collection_type":   r.FormValue("collection_type"),
			"status":            r.FormValue("status"),
		}

		_, err := deps.DB.Update(ctx, "collection", id, data)
		if err != nil {
			log.Printf("Failed to update collection %s: %v", id, err)
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

// NewDeleteAction creates the collection delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection", "delete") {
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

		err := deps.DB.Delete(ctx, "collection", id)
		if err != nil {
			log.Printf("Failed to delete collection %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("collections-table")
	})
}

// NewBulkDeleteAction creates the collection bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}

		for _, id := range ids {
			err := deps.DB.Delete(ctx, "collection", id)
			if err != nil {
				log.Printf("Failed to delete collection %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("collections-table")
	})
}

// NewSetStatusAction creates the collection status update action (POST only).
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection", "update") {
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
		if targetStatus != "pending" && targetStatus != "completed" && targetStatus != "failed" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		if _, err := deps.DB.Update(ctx, "collection", id, map[string]any{"status": targetStatus}); err != nil {
			log.Printf("Failed to update collection status %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("collections-table")
	})
}

// NewBulkSetStatusAction creates the collection bulk status update action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}
		if targetStatus != "pending" && targetStatus != "completed" && targetStatus != "failed" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		for _, id := range ids {
			if _, err := deps.DB.Update(ctx, "collection", id, map[string]any{"status": targetStatus}); err != nil {
				log.Printf("Failed to update collection status %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("collections-table")
	})
}
