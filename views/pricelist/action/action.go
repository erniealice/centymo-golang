package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	pricelistpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/price_list"
)

// FormLabels holds i18n labels for the price list drawer form template.
type FormLabels struct {
	Name            string
	Description     string
	DescPlaceholder string
	DateStart       string
	DateEnd         string
	Active          string
}

// FormData is the template data for the price list drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	DateStart    string
	DateEnd      string
	Active       bool
	Labels       FormLabels
	CommonLabels any
}

// Deps holds dependencies for price list action handlers.
type Deps struct {
	CreatePriceList func(ctx context.Context, req *pricelistpb.CreatePriceListRequest) (*pricelistpb.CreatePriceListResponse, error)
	ReadPriceList   func(ctx context.Context, req *pricelistpb.ReadPriceListRequest) (*pricelistpb.ReadPriceListResponse, error)
	UpdatePriceList func(ctx context.Context, req *pricelistpb.UpdatePriceListRequest) (*pricelistpb.UpdatePriceListResponse, error)
	DeletePriceList func(ctx context.Context, req *pricelistpb.DeletePriceListRequest) (*pricelistpb.DeletePriceListResponse, error)
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		Name:            t("pricelist.form.name"),
		Description:     t("pricelist.form.description"),
		DescPlaceholder: t("pricelist.form.descriptionPlaceholder"),
		DateStart:       t("pricelist.form.dateStart"),
		DateEnd:         t("pricelist.form.dateEnd"),
		Active:          t("pricelist.form.active"),
	}
}

// NewAddAction creates the price list add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("pricelist-drawer-form", &FormData{
				FormAction:   "/action/price-lists/add",
				Active:       true,
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- create price list
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		desc := r.FormValue("description")
		dateEnd := r.FormValue("date_end")

		req := &pricelistpb.CreatePriceListRequest{
			Data: &pricelistpb.PriceList{
				Name:            r.FormValue("name"),
				Description:     &desc,
				DateStartString: r.FormValue("date_start"),
				Active:          active,
			},
		}
		if dateEnd != "" {
			req.Data.DateEndString = &dateEnd
		}

		_, err := deps.CreatePriceList(ctx, req)
		if err != nil {
			log.Printf("Failed to create price list: %v", err)
			return centymo.HTMXError("Failed to create price list")
		}

		return centymo.HTMXSuccess("price-lists-table")
	})
}

// NewEditAction creates the price list edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadPriceList(ctx, &pricelistpb.ReadPriceListRequest{
				Data: &pricelistpb.PriceList{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read price list %s: %v", id, err)
				return centymo.HTMXError("Price list not found")
			}

			data := resp.GetData()
			if len(data) == 0 {
				return centymo.HTMXError("Price list not found")
			}
			pl := data[0]

			return view.OK("pricelist-drawer-form", &FormData{
				FormAction:   "/action/price-lists/edit/" + id,
				IsEdit:       true,
				ID:           id,
				Name:         pl.GetName(),
				Description:  pl.GetDescription(),
				DateStart:    pl.GetDateStartString(),
				DateEnd:      pl.GetDateEndString(),
				Active:       pl.GetActive(),
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST -- update price list
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		desc := r.FormValue("description")
		dateEnd := r.FormValue("date_end")

		req := &pricelistpb.UpdatePriceListRequest{
			Data: &pricelistpb.PriceList{
				Id:              id,
				Name:            r.FormValue("name"),
				Description:     &desc,
				DateStartString: r.FormValue("date_start"),
				Active:          active,
			},
		}
		if dateEnd != "" {
			req.Data.DateEndString = &dateEnd
		}

		_, err := deps.UpdatePriceList(ctx, req)
		if err != nil {
			log.Printf("Failed to update price list %s: %v", id, err)
			return centymo.HTMXError("Failed to update price list")
		}

		return centymo.HTMXSuccess("price-lists-table")
	})
}

// NewDeleteAction creates the price list delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError("Price list ID is required")
		}

		_, err := deps.DeletePriceList(ctx, &pricelistpb.DeletePriceListRequest{
			Data: &pricelistpb.PriceList{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete price list %s: %v", id, err)
			return centymo.HTMXError("Failed to delete price list")
		}

		return centymo.HTMXSuccess("price-lists-table")
	})
}

// NewBulkDeleteAction creates the price list bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields from bulk-action.js.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError("No price list IDs provided")
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
