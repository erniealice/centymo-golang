package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"

	"github.com/erniealice/centymo-golang/views/pricelist/form"
)

// Deps holds dependencies for price list action handlers.
type Deps struct {
	Routes          centymo.PriceListRoutes
	Labels          centymo.PriceListLabels
	CreatePriceList func(ctx context.Context, req *pricelistpb.CreatePriceListRequest) (*pricelistpb.CreatePriceListResponse, error)
	ReadPriceList   func(ctx context.Context, req *pricelistpb.ReadPriceListRequest) (*pricelistpb.ReadPriceListResponse, error)
	UpdatePriceList func(ctx context.Context, req *pricelistpb.UpdatePriceListRequest) (*pricelistpb.UpdatePriceListResponse, error)
	DeletePriceList func(ctx context.Context, req *pricelistpb.DeletePriceListRequest) (*pricelistpb.DeletePriceListResponse, error)
}

// NewAddAction creates the price list add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_list", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("pricelist-drawer-form", &form.Data{
				FormAction:   deps.Routes.AddURL,
				Active:       true,
				Labels:       form.BuildLabels(viewCtx.T, deps.Labels.Form),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- create price list
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		desc := r.FormValue("description")
		dateEnd := r.FormValue("date_end")

		req := &pricelistpb.CreatePriceListRequest{
			Data: &pricelistpb.PriceList{
				Name:        r.FormValue("name"),
				Description: &desc,
				DateStart:   r.FormValue("date_start"),
				Active:      active,
			},
		}
		if dateEnd != "" {
			req.Data.DateEnd = &dateEnd
		}

		_, err := deps.CreatePriceList(ctx, req)
		if err != nil {
			log.Printf("Failed to create price list: %v", err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("price-lists-table")
	})
}

// NewEditAction creates the price list edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_list", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPriceList(ctx, &pricelistpb.ReadPriceListRequest{
				Data: &pricelistpb.PriceList{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read price list %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}

			data := resp.GetData()
			if len(data) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			pl := data[0]

			return view.OK("pricelist-drawer-form", &form.Data{
				FormAction:   route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:       true,
				ID:           id,
				Name:         pl.GetName(),
				Description:  pl.GetDescription(),
				DateStart:    pl.GetDateStart(),
				DateEnd:      pl.GetDateEnd(),
				Active:       pl.GetActive(),
				Labels:       form.BuildLabels(viewCtx.T, deps.Labels.Form),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- update price list
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		desc := r.FormValue("description")
		dateEnd := r.FormValue("date_end")

		req := &pricelistpb.UpdatePriceListRequest{
			Data: &pricelistpb.PriceList{
				Id:          id,
				Name:        r.FormValue("name"),
				Description: &desc,
				DateStart:   r.FormValue("date_start"),
				Active:      active,
			},
		}
		if dateEnd != "" {
			req.Data.DateEnd = &dateEnd
		}

		_, err := deps.UpdatePriceList(ctx, req)
		if err != nil {
			log.Printf("Failed to update price list %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("price-lists-table")
	})
}

// NewDeleteAction creates the price list delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_list", "delete") {
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

		_, err := deps.DeletePriceList(ctx, &pricelistpb.DeletePriceListRequest{
			Data: &pricelistpb.PriceList{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete price list %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("price-lists-table")
	})
}

// NewBulkDeleteAction creates the price list bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("price_list", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}

		for _, id := range ids {
			_, err := deps.DeletePriceList(ctx, &pricelistpb.DeletePriceListRequest{
				Data: &pricelistpb.PriceList{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete price list %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("price-lists-table")
	})
}
