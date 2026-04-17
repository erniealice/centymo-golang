package action

import (
	"context"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
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
	Product            string
	ProductPlaceholder string
	SelectProduct      string
	Active             string
}

// ProductPlanFormData is the template data for the product plan drawer form.
type ProductPlanFormData struct {
	FormAction           string
	IsEdit               bool
	ID                   string
	PlanID               string
	Name                 string
	Mode                 string
	SelectedProductID    string
	SelectedProductLabel string
	Active               bool
	Products             []*ProductOption
	ProductOptions       []map[string]any
	Labels               ProductPlanFormLabels
	CommonLabels         any
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
	ListProductPlans  func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
}

// productPlanFormLabels converts centymo.ProductPlanFormLabels into the local type.
func productPlanFormLabels(l centymo.ProductPlanFormLabels) ProductPlanFormLabels {
	return ProductPlanFormLabels{
		Product:            l.Product,
		ProductPlaceholder: l.ProductPlaceholder,
		SelectProduct:      l.SelectProduct,
		Active:             l.Active,
	}
}

// loadProductOptions fetches the product list and converts to options.
// Returns nil slice on error (graceful degradation).
//
// Currently UNFILTERED — shows all products regardless of product_kind.
// The bundle/package drawer kind-selector design is in research:
// docs/plan/20260417-product-taxonomy-refactor/{claude,codex}-research.md
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

// buildProductAutoCompleteOptions converts []*ProductOption to the auto-complete compatible format.
// disabledIDs is a set of product IDs that should be shown as disabled (already added to the plan).
func buildProductAutoCompleteOptions(products []*ProductOption, selectedID string, disabledIDs map[string]bool) []map[string]any {
	opts := make([]map[string]any, 0, len(products))
	for _, p := range products {
		opt := map[string]any{
			"Value":    p.Id,
			"Label":    p.Name,
			"Selected": p.Id == selectedID,
		}
		if disabledIDs[p.Id] {
			opt["Disabled"] = true
		}
		opts = append(opts, opt)
	}
	return opts
}

// loadExistingProductIDs returns a set of product IDs already added to the given plan.
func loadExistingProductIDs(ctx context.Context, deps *ProductPlanDeps, planID string) map[string]bool {
	if deps.ListProductPlans == nil {
		return nil
	}
	resp, err := deps.ListProductPlans(ctx, &productplanpb.ListProductPlansRequest{
		Filters: &commonpb.FilterRequest{
			Logic: commonpb.FilterLogic_AND,
			Filters: []*commonpb.TypedFilter{
				{
					Field: "plan_id",
					FilterType: &commonpb.TypedFilter_StringFilter{
						StringFilter: &commonpb.StringFilter{
							Value:    planID,
							Operator: commonpb.StringOperator_STRING_EQUALS,
						},
					},
				},
			},
		},
	})
	if err != nil {
		log.Printf("Failed to load existing product plans for plan %s: %v", planID, err)
		return nil
	}
	ids := make(map[string]bool)
	for _, pp := range resp.GetData() {
		ids[pp.GetProductId()] = true
	}
	return ids
}

// findProductLabel returns the name of the product with the given ID, or empty string.
func findProductLabel(products []*ProductOption, id string) string {
	for _, p := range products {
		if p.Id == id {
			return p.Name
		}
	}
	return ""
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
			products := loadProductOptions(ctx, deps)
			disabledIDs := loadExistingProductIDs(ctx, deps, planID)
			return view.OK("product-plan-drawer-form", &ProductPlanFormData{
				FormAction:     route.ResolveURL(deps.Routes.ProductPlanAddURL, "id", planID),
				PlanID:         planID,
				Name:           "",
				Mode:           "service",
				Active:         true,
				Products:       products,
				ProductOptions: buildProductAutoCompleteOptions(products, "", disabledIDs),
				Labels:         productPlanFormLabels(deps.Labels.ProductPlanForm),
				CommonLabels:   nil, // injected by ViewAdapter
			})
		}

		// POST — create product plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		productID := r.FormValue("product_id")

		name := r.FormValue("name")
		if name == "" {
			name = r.FormValue("product_name")
		}
		if name == "" && productID != "" {
			products := loadProductOptions(ctx, deps)
			name = findProductLabel(products, productID)
		}

		// No price/currency at creation — pricing is set via price plans
		pp := &productplanpb.ProductPlan{
			PlanId:    planID,
			ProductId: productID,
			Name:      name,
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

			products := loadProductOptions(ctx, deps)
			selectedProductID := pp.GetProductId()
			return view.OK("product-plan-drawer-form", &ProductPlanFormData{
				FormAction:           route.ResolveURL(deps.Routes.ProductPlanEditURL, "id", planID, "ppid", ppID),
				IsEdit:               true,
				ID:                   ppID,
				PlanID:               planID,
				Name:                 pp.GetName(),
				Mode:                 "service",
				SelectedProductID:    selectedProductID,
				SelectedProductLabel: findProductLabel(products, selectedProductID),
				Active:               pp.GetActive(),
				Products:             products,
				ProductOptions:       buildProductAutoCompleteOptions(products, selectedProductID, nil),
				Labels:               productPlanFormLabels(deps.Labels.ProductPlanForm),
				CommonLabels:         nil, // injected by ViewAdapter
			})
		}

		// POST — update product plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		productID := r.FormValue("product_id")

		name := r.FormValue("name")
		if name == "" {
			name = r.FormValue("product_name")
		}
		if name == "" && productID != "" {
			products := loadProductOptions(ctx, deps)
			name = findProductLabel(products, productID)
		}

		pp := &productplanpb.ProductPlan{
			Id:        ppID,
			PlanId:    planID,
			ProductId: productID,
			Name:      name,
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
