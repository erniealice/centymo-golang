package action

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/erniealice/centymo-golang"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	collectionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection"
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
	// Form holds extended placeholder and option labels used by the template.
	Form centymo.CollectionFormLabels
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
	Routes           centymo.CollectionRoutes
	Labels           centymo.CollectionLabels
	CreateCollection func(ctx context.Context, req *collectionpb.CreateCollectionRequest) (*collectionpb.CreateCollectionResponse, error)
	ReadCollection   func(ctx context.Context, req *collectionpb.ReadCollectionRequest) (*collectionpb.ReadCollectionResponse, error)
	UpdateCollection func(ctx context.Context, req *collectionpb.UpdateCollectionRequest) (*collectionpb.UpdateCollectionResponse, error)
	DeleteCollection func(ctx context.Context, req *collectionpb.DeleteCollectionRequest) (*collectionpb.DeleteCollectionResponse, error)
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
		Form:                 l,
	}
}

// parseAmount converts a form string amount (decimal) to int64 centavos.
func parseAmount(s string) int64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return int64(math.Round(f * 100))
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

		resp, err := deps.CreateCollection(ctx, &collectionpb.CreateCollectionRequest{
			Data: &collectionpb.Collection{
				ReferenceNumber:    r.FormValue("reference_number"),
				Name:               r.FormValue("customer"),
				Amount:             parseAmount(r.FormValue("amount")),
				Currency:           r.FormValue("currency"),
				CollectionMethodId: r.FormValue("collection_method"),
				ReceivedBy:         r.FormValue("received_by"),
				ReceivedRole:       r.FormValue("received_role"),
				CollectionType:     r.FormValue("collection_type"),
				Status:             r.FormValue("status"),
			},
		})
		if err != nil {
			log.Printf("Failed to create collection: %v", err)
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
			readResp, err := deps.ReadCollection(ctx, &collectionpb.ReadCollectionRequest{
				Data: &collectionpb.Collection{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read collection %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			readData := readResp.GetData()
			if len(readData) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := readData[0]

			return view.OK("collection-drawer-form", &FormData{
				FormAction:       route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:           true,
				ID:               id,
				Customer:         record.GetName(),
				ReferenceNumber:  record.GetReferenceNumber(),
				Amount:           fmt.Sprintf("%.2f", float64(record.GetAmount())/100.0),
				Currency:         record.GetCurrency(),
				CollectionMethod: record.GetCollectionMethodId(),
				Date:             record.GetDateCreatedString(),
				ReceivedBy:       record.GetReceivedBy(),
				ReceivedRole:     record.GetReceivedRole(),
				Notes:            "",
				CollectionType:   record.GetCollectionType(),
				Status:           record.GetStatus(),
				Labels:           formLabels(centymo.DefaultCollectionLabels().Form),
				CommonLabels:     nil, // injected by ViewAdapter
			})
		}

		// POST — update collection
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		_, err := deps.UpdateCollection(ctx, &collectionpb.UpdateCollectionRequest{
			Data: &collectionpb.Collection{
				Id:                 id,
				ReferenceNumber:    r.FormValue("reference_number"),
				Name:               r.FormValue("customer"),
				Amount:             parseAmount(r.FormValue("amount")),
				Currency:           r.FormValue("currency"),
				CollectionMethodId: r.FormValue("collection_method"),
				ReceivedBy:         r.FormValue("received_by"),
				ReceivedRole:       r.FormValue("received_role"),
				CollectionType:     r.FormValue("collection_type"),
				Status:             r.FormValue("status"),
			},
		})
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

		_, err := deps.DeleteCollection(ctx, &collectionpb.DeleteCollectionRequest{
			Data: &collectionpb.Collection{Id: id},
		})
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
			_, err := deps.DeleteCollection(ctx, &collectionpb.DeleteCollectionRequest{
				Data: &collectionpb.Collection{Id: id},
			})
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

		if _, err := deps.UpdateCollection(ctx, &collectionpb.UpdateCollectionRequest{
			Data: &collectionpb.Collection{Id: id, Status: targetStatus},
		}); err != nil {
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
			if _, err := deps.UpdateCollection(ctx, &collectionpb.UpdateCollectionRequest{
				Data: &collectionpb.Collection{Id: id, Status: targetStatus},
			}); err != nil {
				log.Printf("Failed to update collection status %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("collections-table")
	})
}
