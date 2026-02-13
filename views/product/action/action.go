package action

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
)

// FormLabels holds i18n labels for the product drawer form template.
type FormLabels struct {
	Name            string
	Description     string
	DescPlaceholder string
	Price           string
	Currency        string
	Active          string
}

// FormData is the template data for the product drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	Price        string
	Currency     string
	Active       bool
	Labels       FormLabels
	CommonLabels any
}

// Deps holds dependencies for product action handlers.
type Deps struct {
	CreateProduct    func(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error)
	ReadProduct      func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	UpdateProduct    func(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error)
	DeleteProduct    func(ctx context.Context, req *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error)
	SetProductActive func(ctx context.Context, id string, active bool) error
}

func formLabels(t func(string) string) FormLabels {
	return FormLabels{
		Name:            t("product.form.name"),
		Description:     t("product.form.description"),
		DescPlaceholder: t("product.form.descriptionPlaceholder"),
		Price:           t("product.form.price"),
		Currency:        t("product.form.currency"),
		Active:          t("product.form.active"),
	}
}

// NewAddAction creates the product add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("product-drawer-form", &FormData{
				FormAction:   "/action/products/add",
				Active:       true,
				Currency:     "PHP",
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create product
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		desc := r.FormValue("description")

		_, err := deps.CreateProduct(ctx, &productpb.CreateProductRequest{
			Data: &productpb.Product{
				Name:        r.FormValue("name"),
				Description: &desc,
				Price:       price,
				Currency:    r.FormValue("currency"),
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to create product: %v", err)
			return centymo.HTMXError("Failed to create product")
		}

		return centymo.HTMXSuccess("products-table")
	})
}

// NewEditAction creates the product edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
				Data: &productpb.Product{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read product %s: %v", id, err)
				return centymo.HTMXError("Product not found")
			}

			data := resp.GetData()
			if len(data) == 0 {
				return centymo.HTMXError("Product not found")
			}
			p := data[0]

			return view.OK("product-drawer-form", &FormData{
				FormAction:   "/action/products/edit/" + id,
				IsEdit:       true,
				ID:           id,
				Name:         p.GetName(),
				Description:  p.GetDescription(),
				Price:        strconv.FormatFloat(p.GetPrice(), 'f', 2, 64),
				Currency:     p.GetCurrency(),
				Active:       p.GetActive(),
				Labels:       formLabels(viewCtx.T),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — update product
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		desc := r.FormValue("description")

		_, err := deps.UpdateProduct(ctx, &productpb.UpdateProductRequest{
			Data: &productpb.Product{
				Id:          id,
				Name:        r.FormValue("name"),
				Description: &desc,
				Price:       price,
				Currency:    r.FormValue("currency"),
				Active:      active,
			},
		})
		if err != nil {
			log.Printf("Failed to update product %s: %v", id, err)
			return centymo.HTMXError("Failed to update product")
		}

		return centymo.HTMXSuccess("products-table")
	})
}

// NewDeleteAction creates the product delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError("Product ID is required")
		}

		_, err := deps.DeleteProduct(ctx, &productpb.DeleteProductRequest{
			Data: &productpb.Product{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete product %s: %v", id, err)
			return centymo.HTMXError("Failed to delete product")
		}

		return centymo.HTMXSuccess("products-table")
	})
}

// NewBulkDeleteAction creates the product bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError("No product IDs provided")
		}

		for _, id := range ids {
			_, err := deps.DeleteProduct(ctx, &productpb.DeleteProductRequest{
				Data: &productpb.Product{Id: id},
			})
			if err != nil {
				log.Printf("Failed to delete product %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("products-table")
	})
}

// NewSetStatusAction creates the product activate/deactivate action (POST only).
// Expects query params: ?id={productId}&status={active|inactive}
//
// Uses SetProductActive (raw map update) instead of protobuf because
// proto3's protojson omits bool fields with value false, which means
// deactivation (active=false) would silently be skipped.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError("Product ID is required")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError("Invalid status")
		}

		if err := deps.SetProductActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update product status %s: %v", id, err)
			return centymo.HTMXError("Failed to update product status")
		}

		return centymo.HTMXSuccess("products-table")
	})
}

// NewBulkSetStatusAction creates the product bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError("No product IDs provided")
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError("Invalid target status")
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetProductActive(ctx, id, active); err != nil {
				log.Printf("Failed to update product status %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("products-table")
	})
}
