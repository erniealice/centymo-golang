package action

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
)

// ProductOption is a minimal struct for rendering product options in the product plan form.
type ProductOption struct {
	Id   string
	Name string
}

// ProductPlanFormLabels holds i18n labels for the product plan drawer form template.
type ProductPlanFormLabels struct {
	Product             string
	ProductPlaceholder  string
	SelectProduct       string
	Price               string
	PricePlaceholder    string
	Currency            string
	CurrencyPlaceholder string
	Active              string
}

// ProductPlanFormData is the template data for the product plan drawer form.
type ProductPlanFormData struct {
	FormAction        string
	IsEdit            bool
	ID                string
	PlanID            string
	Name              string
	SelectedProductID string
	Price             string
	Currency          string
	Active            bool
	Products          []*ProductOption
	Labels            ProductPlanFormLabels
	CommonLabels      any
}

// ProductPlanDeps holds dependencies for product plan action handlers.
type ProductPlanDeps struct {
	Routes            centymo.PlanRoutes
	Labels            centymo.PlanLabels
	CreateProductPlan func(ctx context.Context, req *productplanpb.CreateProductPlanRequest) (*productplanpb.CreateProductPlanResponse, error)
	ReadProductPlan   func(ctx context.Context, req *productplanpb.ReadProductPlanRequest) (*productplanpb.ReadProductPlanResponse, error)
	UpdateProductPlan func(ctx context.Context, req *productplanpb.UpdateProductPlanRequest) (*productplanpb.UpdateProductPlanResponse, error)
	DeleteProductPlan func(ctx context.Context, req *productplanpb.DeleteProductPlanRequest) (*productplanpb.DeleteProductPlanResponse, error)
	ListProducts      func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
}

// productPlanFormLabels converts centymo.ProductPlanFormLabels into the local type.
func productPlanFormLabels(l centymo.ProductPlanFormLabels) ProductPlanFormLabels {
	return ProductPlanFormLabels{
		Product:             l.Product,
		ProductPlaceholder:  l.ProductPlaceholder,
		SelectProduct:       l.SelectProduct,
		Price:               l.Price,
		PricePlaceholder:    l.PricePlaceholder,
		Currency:            l.Currency,
		CurrencyPlaceholder: l.CurrencyPlaceholder,
		Active:              l.Active,
	}
}

// loadProductOptions fetches the product list and converts to options.
// Returns nil slice on error (graceful degradation).
func loadProductOptions(ctx context.Context, deps *ProductPlanDeps) []*ProductOption {
	if deps.ListProducts == nil {
		return nil
	}
	resp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		log.Printf("Failed to load products for product plan form: %v", err)
		return nil
	}
	var options []*ProductOption
	for _, p := range resp.GetData() {
		options = append(options, &ProductOption{
			Id:   p.GetId(),
			Name: p.GetName(),
		})
	}
	return options
}

// NewProductPlanAddAction creates the product plan add action (GET = form, POST = create).
// URL: /action/plans/{id}/products/add
func NewProductPlanAddAction(deps *ProductPlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		planID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("product-plan-drawer-form", &ProductPlanFormData{
				FormAction:   route.ResolveURL(deps.Routes.ProductPlanAddURL, "id", planID),
				PlanID:       planID,
				Name:         "",
				Active:       true,
				Currency:     "PHP",
				Products:     loadProductOptions(ctx, deps),
				Labels:       productPlanFormLabels(deps.Labels.ProductPlanForm),
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create product plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		price := int64(0)
		if v, err := strconv.ParseFloat(r.FormValue("price"), 64); err == nil {
			price = int64(v * 100)
		}

		// Use product name as the product plan name (required field, min 3 chars enforced by usecase)
		name := r.FormValue("name")
		if name == "" {
			name = r.FormValue("product_name")
		}

		pp := &productplanpb.ProductPlan{
			PlanId:    planID,
			ProductId: r.FormValue("product_id"),
			Name:      name,
			Price:     price,
			Currency:  r.FormValue("currency"),
			Active:    active,
		}

		_, err := deps.CreateProductPlan(ctx, &productplanpb.CreateProductPlanRequest{
			Data: pp,
		})
		if err != nil {
			log.Printf("Failed to create product plan for plan %s: %v", planID, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plan-products-table")
	})
}

// NewProductPlanEditAction creates the product plan edit action (GET = form, POST = update).
// URL: /action/plans/{id}/products/edit/{ppid}
func NewProductPlanEditAction(deps *ProductPlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		planID := viewCtx.Request.PathValue("id")
		ppID := viewCtx.Request.PathValue("ppid")

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadProductPlan(ctx, &productplanpb.ReadProductPlanRequest{
				Data: &productplanpb.ProductPlan{Id: ppID},
			})
			if err != nil {
				log.Printf("Failed to read product plan %s: %v", ppID, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			data := resp.GetData()
			if len(data) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			pp := data[0]

			priceStr := strconv.FormatFloat(float64(pp.GetPrice())/100.0, 'f', 2, 64)

			return view.OK("product-plan-drawer-form", &ProductPlanFormData{
				FormAction:        route.ResolveURL(deps.Routes.ProductPlanEditURL, "id", planID, "ppid", ppID),
				IsEdit:            true,
				ID:                ppID,
				PlanID:            planID,
				Name:              pp.GetName(),
				SelectedProductID: pp.GetProductId(),
				Price:             priceStr,
				Currency:          pp.GetCurrency(),
				Active:            pp.GetActive(),
				Products:          loadProductOptions(ctx, deps),
				Labels:            productPlanFormLabels(deps.Labels.ProductPlanForm),
				CommonLabels:      nil, // injected by ViewAdapter
			})
		}

		// POST — update product plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		price := int64(0)
		if v, err := strconv.ParseFloat(r.FormValue("price"), 64); err == nil {
			price = int64(v * 100)
		}

		name := r.FormValue("name")
		if name == "" {
			name = r.FormValue("product_name")
		}

		pp := &productplanpb.ProductPlan{
			Id:        ppID,
			PlanId:    planID,
			ProductId: r.FormValue("product_id"),
			Name:      name,
			Price:     price,
			Currency:  r.FormValue("currency"),
			Active:    active,
		}

		_, err := deps.UpdateProductPlan(ctx, &productplanpb.UpdateProductPlanRequest{
			Data: pp,
		})
		if err != nil {
			log.Printf("Failed to update product plan %s: %v", ppID, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plan-products-table")
	})
}

// NewProductPlanDeleteAction creates the product plan delete action (POST only).
// URL: /action/plans/{id}/products/delete  (product plan id via query param or form)
func NewProductPlanDeleteAction(deps *ProductPlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		ppID := viewCtx.Request.URL.Query().Get("id")
		if ppID == "" {
			_ = viewCtx.Request.ParseForm()
			ppID = viewCtx.Request.FormValue("id")
		}
		if ppID == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeleteProductPlan(ctx, &productplanpb.DeleteProductPlanRequest{
			Data: &productplanpb.ProductPlan{Id: ppID},
		})
		if err != nil {
			log.Printf("Failed to delete product plan %s: %v", ppID, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plan-products-table")
	})
}
