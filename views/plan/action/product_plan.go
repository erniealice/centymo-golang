package action

import (
	"context"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productoptionvaluepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option_value"
	productplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_plan"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
	productvariantoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_option"
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

	// Model D — variant picker
	VariantSelectLabel       string
	VariantSelectPlaceholder string
	VariantSelectInfo        string
}

// KindOption is a value/label pair for the product_kind selector on the
// product plan drawer form.
type KindOption struct {
	Value string
	Label string
}

// ProductPlanFormData is the template data for the product plan drawer form.
type ProductPlanFormData struct {
	FormAction           string
	PickerURL            string
	VariantPickerURL     string
	IsEdit               bool
	ID                   string
	PlanID               string
	Name                 string
	Mode                 string
	SelectedKind         string
	KindOptions          []KindOption
	KindLabel            string
	ProductLabel         string
	ProductPlaceholder   string
	SelectedProductID    string
	SelectedProductLabel string
	Active               bool
	Products             []*ProductOption
	ProductOptions       []map[string]any
	// Model D — variant picker state. Rendered only when the selected
	// product has variant_mode = "configurable".
	VariantConfigurable bool
	SelectedVariantID   string
	VariantOptions      []types.SelectOption
	Labels              ProductPlanFormLabels
	CommonLabels        any
}

// ProductPlanDeps holds dependencies for product plan action handlers.
type ProductPlanDeps struct {
	Routes              centymo.PlanRoutes
	Labels              centymo.PlanLabels
	CreateProductPlan   func(ctx context.Context, req *productplanpb.CreateProductPlanRequest) (*productplanpb.CreateProductPlanResponse, error)
	ReadProductPlan     func(ctx context.Context, req *productplanpb.ReadProductPlanRequest) (*productplanpb.ReadProductPlanResponse, error)
	UpdateProductPlan   func(ctx context.Context, req *productplanpb.UpdateProductPlanRequest) (*productplanpb.UpdateProductPlanResponse, error)
	DeleteProductPlan   func(ctx context.Context, req *productplanpb.DeleteProductPlanRequest) (*productplanpb.DeleteProductPlanResponse, error)
	ListProducts        func(ctx context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error)
	ListProductPlans    func(ctx context.Context, req *productplanpb.ListProductPlansRequest) (*productplanpb.ListProductPlansResponse, error)
	// Model D — variant lookup for the drawer's variant sub-picker.
	ListProductVariants func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	// Model D — variant label enrichment. Used to render labels of the form
	// "SKU — Red / Large / Cotton" in the variant sub-picker, with the trailing
	// option-value labels ordered by ProductOption.sort_order ASC.
	ListProductVariantOptions func(ctx context.Context, req *productvariantoptionpb.ListProductVariantOptionsRequest) (*productvariantoptionpb.ListProductVariantOptionsResponse, error)
	ListProductOptionValues   func(ctx context.Context, req *productoptionvaluepb.ListProductOptionValuesRequest) (*productoptionvaluepb.ListProductOptionValuesResponse, error)
	ListProductOptions        func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
}

// productPlanFormLabels converts centymo.ProductPlanFormLabels into the local type.
func productPlanFormLabels(l centymo.ProductPlanFormLabels) ProductPlanFormLabels {
	return ProductPlanFormLabels{
		Product:                  l.Product,
		ProductPlaceholder:       l.ProductPlaceholder,
		SelectProduct:            l.SelectProduct,
		Active:                   l.Active,
		VariantSelectLabel:       l.VariantSelectLabel,
		VariantSelectPlaceholder: l.VariantSelectPlaceholder,
		VariantSelectInfo:        l.VariantSelectInfo,
	}
}

// buildKindOptions returns the ordered list of product_kind selector options,
// with labels sourced from the caller's lyngua-driven ProductKind labels.
func buildKindOptions(l centymo.ProductKindOptionLabels) []KindOption {
	return []KindOption{
		{Value: "service", Label: fallbackLabel(l.Service, "Service")},
		{Value: "stocked_good", Label: fallbackLabel(l.StockedGood, "Stocked Good")},
		{Value: "non_stocked_good", Label: fallbackLabel(l.NonStockedGood, "Non-Stocked Good")},
		{Value: "consumable", Label: fallbackLabel(l.Consumable, "Consumable")},
	}
}

func fallbackLabel(primary, fallback string) string {
	if primary == "" {
		return fallback
	}
	return primary
}

// loadProductOptions fetches the product list and converts to options.
// Returns nil slice on error (graceful degradation).
//
// When kind is non-empty, results are filtered to products whose
// product_kind matches exactly (e.g., "service", "stocked_good",
// "non_stocked_good", "consumable"). Empty kind = no filter.
func loadProductOptions(ctx context.Context, deps *ProductPlanDeps, kind string) []*ProductOption {
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
		if kind != "" && p.GetProductKind() != kind {
			continue
		}
		options = append(options, &ProductOption{
			Id:   p.GetId(),
			Name: p.GetName(),
		})
	}
	sort.SliceStable(options, func(i, j int) bool {
		return strings.ToLower(options[i].Name) < strings.ToLower(options[j].Name)
	})
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

// lookupProductKind resolves a product's product_kind by scanning the unfiltered
// product list. Returns empty string if the product is missing or ListProducts
// is unavailable, so callers can fall back to a sensible default.
func lookupProductKind(ctx context.Context, deps *ProductPlanDeps, productID string) string {
	if productID == "" || deps.ListProducts == nil {
		return ""
	}
	resp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		log.Printf("Failed to look up product kind for %s: %v", productID, err)
		return ""
	}
	for _, p := range resp.GetData() {
		if p.GetId() == productID {
			return p.GetProductKind()
		}
	}
	return ""
}

// lookupProductVariantMode resolves a product's variant_mode ("none" or
// "configurable") for Model D's binary invariant. Missing product returns "".
func lookupProductVariantMode(ctx context.Context, deps *ProductPlanDeps, productID string) string {
	if productID == "" || deps.ListProducts == nil {
		return ""
	}
	resp, err := deps.ListProducts(ctx, &productpb.ListProductsRequest{})
	if err != nil {
		return ""
	}
	for _, p := range resp.GetData() {
		if p.GetId() == productID {
			return p.GetVariantMode()
		}
	}
	return ""
}

// loadVariantOptions builds a select list for the given product's variants.
//
// Label format: "SKU — Red / Large / Cotton" where the trailing concat is the
// variant's option-value labels, ordered by the corresponding
// ProductOption.sort_order ASC, joined with " / ". When a variant has no
// product_variant_option rows (data drift, or product hasn't been configured
// yet), the label degrades to plain SKU (no en-dash, no concat).
//
// SKU falls back to the variant ID when blank. Final []types.SelectOption is
// sorted alphabetically by Label ASC (case-insensitive).
func loadVariantOptions(ctx context.Context, deps *ProductPlanDeps, productID, selectedID string) []types.SelectOption {
	if productID == "" || deps.ListProductVariants == nil {
		return nil
	}
	resp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{})
	if err != nil {
		log.Printf("Failed to list product variants for %s: %v", productID, err)
		return nil
	}

	// Build variantID → []optionValueLabel ordered by ProductOption.sort_order ASC.
	// Mirrors the lookup in views/product/detail/page.go BuildVariantsTable.
	variantOptionLabels := buildVariantOptionLabels(ctx, deps)

	options := []types.SelectOption{}
	for _, v := range resp.GetData() {
		if v == nil || v.GetProductId() != productID {
			continue
		}
		// Skip inactive variants — they shouldn't be assignable to new
		// catalog lines. Edit-mode callers pass selectedID; if the
		// previously-selected variant is now inactive, surface it anyway
		// so the picker round-trips correctly.
		if !v.GetActive() && v.GetId() != selectedID {
			continue
		}
		sku := v.GetSku()
		if sku == "" {
			sku = v.GetId()
		}
		label := sku
		if parts := variantOptionLabels[v.GetId()]; len(parts) > 0 {
			label = sku + " — " + strings.Join(parts, centymo.OptionValueSeparator)
		}
		options = append(options, types.SelectOption{
			Value:    v.GetId(),
			Label:    label,
			Selected: v.GetId() == selectedID,
		})
	}

	sort.SliceStable(options, func(i, j int) bool {
		return strings.ToLower(options[i].Label) < strings.ToLower(options[j].Label)
	})
	return options
}

// buildVariantOptionLabels returns variantID → []optionValueLabel where the
// labels are ordered by their parent ProductOption.sort_order ASC. Labels are
// the human-readable ProductOptionValue.label values (e.g., "Red", "Large").
// Returns an empty map when any of the three lookup deps are missing or when
// the upstream calls fail — callers fall back to plain SKU labels.
func buildVariantOptionLabels(ctx context.Context, deps *ProductPlanDeps) map[string][]string {
	out := map[string][]string{}
	if deps.ListProductVariantOptions == nil || deps.ListProductOptionValues == nil || deps.ListProductOptions == nil {
		return out
	}

	voResp, err := deps.ListProductVariantOptions(ctx, &productvariantoptionpb.ListProductVariantOptionsRequest{})
	if err != nil {
		log.Printf("Failed to list product variant options: %v", err)
		return out
	}
	ovResp, err := deps.ListProductOptionValues(ctx, &productoptionvaluepb.ListProductOptionValuesRequest{})
	if err != nil {
		log.Printf("Failed to list product option values: %v", err)
		return out
	}
	optResp, err := deps.ListProductOptions(ctx, &productoptionpb.ListProductOptionsRequest{})
	if err != nil {
		log.Printf("Failed to list product options: %v", err)
		return out
	}

	// productOptionID → sort_order
	optionSortOrder := map[string]int32{}
	for _, o := range optResp.GetData() {
		if o == nil {
			continue
		}
		optionSortOrder[o.GetId()] = o.GetSortOrder()
	}

	// productOptionValueID → (label, parent sort_order)
	type valueInfo struct {
		label     string
		sortOrder int32
	}
	valueLookup := map[string]valueInfo{}
	for _, ov := range ovResp.GetData() {
		if ov == nil {
			continue
		}
		valueLookup[ov.GetId()] = valueInfo{
			label:     ov.GetLabel(),
			sortOrder: optionSortOrder[ov.GetProductOptionId()],
		}
	}

	// variantID → []{label, sort_order}, then sort and project.
	type entry struct {
		label     string
		sortOrder int32
	}
	bucket := map[string][]entry{}
	for _, vo := range voResp.GetData() {
		if vo == nil {
			continue
		}
		varID := vo.GetProductVariantId()
		info, ok := valueLookup[vo.GetProductOptionValueId()]
		if !ok || varID == "" {
			continue
		}
		bucket[varID] = append(bucket[varID], entry{label: info.label, sortOrder: info.sortOrder})
	}
	for varID, entries := range bucket {
		sort.SliceStable(entries, func(i, j int) bool {
			return entries[i].sortOrder < entries[j].sortOrder
		})
		labels := make([]string, 0, len(entries))
		for _, e := range entries {
			labels = append(labels, e.label)
		}
		out[varID] = labels
	}
	return out
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
			selectedKind := "service"
			products := loadProductOptions(ctx, deps, selectedKind)
			disabledIDs := loadExistingProductIDs(ctx, deps, planID)
			formLabels := deps.Labels.ProductPlanForm
			return view.OK("product-plan-drawer-form", &ProductPlanFormData{
				FormAction:         route.ResolveURL(deps.Routes.ProductPlanAddURL, "id", planID),
				PickerURL:          route.ResolveURL(deps.Routes.ProductPlanPickerURL, "id", planID),
				VariantPickerURL:   route.ResolveURL(deps.Routes.ProductPlanPickerURL, "id", planID),
				PlanID:             planID,
				Name:               "",
				Mode:               "service",
				SelectedKind:       selectedKind,
				KindOptions:        buildKindOptions(formLabels.ProductKind),
				KindLabel:          fallbackLabel(formLabels.ProductKindLabel, "Item Type"),
				ProductLabel:       formLabels.Product,
				ProductPlaceholder: formLabels.ProductPlaceholder,
				Active:             true,
				Products:           products,
				ProductOptions:     buildProductAutoCompleteOptions(products, "", disabledIDs),
				// No product selected yet — variant section is hidden until HTMX
				// re-renders the picker with a chosen product_id.
				VariantConfigurable: false,
				Labels:              productPlanFormLabels(formLabels),
				CommonLabels:        nil, // injected by ViewAdapter
			})
		}

		// POST — create product plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		productID := r.FormValue("product_id")
		variantID := r.FormValue("product_variant_id")

		name := r.FormValue("name")
		if name == "" {
			name = r.FormValue("product_name")
		}
		if name == "" && productID != "" {
			products := loadProductOptions(ctx, deps, "")
			name = findProductLabel(products, productID)
		}

		// Model D binary invariant (defensive — use-case layer also enforces):
		// require variant_id iff the parent product is variant-configurable.
		variantMode := lookupProductVariantMode(ctx, deps, productID)
		if variantMode == "configurable" && variantID == "" {
			return centymo.HTMXError("Please select a variant for this product.")
		}
		if variantMode != "configurable" && variantID != "" {
			// Simple products never carry a variant.
			variantID = ""
		}

		pp := &productplanpb.ProductPlan{
			PlanId:    planID,
			ProductId: productID,
			Name:      name,
			Active:    active,
		}
		if variantID != "" {
			pp.ProductVariantId = &variantID
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

			selectedProductID := pp.GetProductId()
			// Derive kind from the saved product's kind. Fall back to "service"
			// when the product can't be located (legacy data) so the picker
			// still has a sensible default to filter by.
			selectedKind := "service"
			allProducts := loadProductOptions(ctx, deps, "")
			if kind := lookupProductKind(ctx, deps, selectedProductID); kind != "" {
				selectedKind = kind
			}
			products := loadProductOptions(ctx, deps, selectedKind)
			variantMode := lookupProductVariantMode(ctx, deps, selectedProductID)
			selectedVariantID := pp.GetProductVariantId()
			var variantOptions []types.SelectOption
			if variantMode == "configurable" {
				variantOptions = loadVariantOptions(ctx, deps, selectedProductID, selectedVariantID)
			}
			formLabels := deps.Labels.ProductPlanForm
			return view.OK("product-plan-drawer-form", &ProductPlanFormData{
				FormAction:           route.ResolveURL(deps.Routes.ProductPlanEditURL, "id", planID, "ppid", ppID),
				PickerURL:            route.ResolveURL(deps.Routes.ProductPlanPickerURL, "id", planID),
				VariantPickerURL:     route.ResolveURL(deps.Routes.ProductPlanPickerURL, "id", planID),
				IsEdit:               true,
				ID:                   ppID,
				PlanID:               planID,
				Name:                 pp.GetName(),
				Mode:                 "service",
				SelectedKind:         selectedKind,
				KindOptions:          buildKindOptions(formLabels.ProductKind),
				KindLabel:            fallbackLabel(formLabels.ProductKindLabel, "Item Type"),
				ProductLabel:         formLabels.Product,
				ProductPlaceholder:   formLabels.ProductPlaceholder,
				SelectedProductID:    selectedProductID,
				SelectedProductLabel: findProductLabel(allProducts, selectedProductID),
				Active:               pp.GetActive(),
				Products:             products,
				ProductOptions:       buildProductAutoCompleteOptions(products, selectedProductID, nil),
				VariantConfigurable:  variantMode == "configurable",
				SelectedVariantID:    selectedVariantID,
				VariantOptions:       variantOptions,
				Labels:               productPlanFormLabels(formLabels),
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
		variantID := r.FormValue("product_variant_id")

		name := r.FormValue("name")
		if name == "" {
			name = r.FormValue("product_name")
		}
		if name == "" && productID != "" {
			products := loadProductOptions(ctx, deps, "")
			name = findProductLabel(products, productID)
		}

		variantMode := lookupProductVariantMode(ctx, deps, productID)
		if variantMode == "configurable" && variantID == "" {
			return centymo.HTMXError("Please select a variant for this product.")
		}
		if variantMode != "configurable" && variantID != "" {
			variantID = ""
		}

		pp := &productplanpb.ProductPlan{
			Id:        ppID,
			PlanId:    planID,
			ProductId: productID,
			Name:      name,
			Active:    active,
		}
		if variantID != "" {
			pp.ProductVariantId = &variantID
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

// PickerPartialData is the template data for the product-picker-partial template.
// Rendered by NewProductPlanPickerAction when the kind selector changes.
type PickerPartialData struct {
	PlanID               string
	Name                 string
	SelectedKind         string
	SelectedProductID    string
	SelectedProductLabel string
	ProductLabel         string
	ProductPlaceholder   string
	ProductOptions       []map[string]any
	// Model D — variant sub-picker rendered when the selected product is
	// variant-configurable.
	VariantConfigurable      bool
	VariantOptions           []types.SelectOption
	VariantSelectLabel       string
	VariantSelectPlaceholder string
}

// NewProductPlanPickerAction handles GET /action/plan/{id}/products/picker?product_kind=...&product_id=...
// Returns only the product-picker-partial template, filtered by the requested
// kind. When a product_id is supplied and that product has variant_mode =
// "configurable", the partial also surfaces a variant sub-picker populated
// from ListProductVariants(product_id). Swapped into #product-picker-wrapper
// on the drawer form via HTMX.
func NewProductPlanPickerAction(deps *ProductPlanDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("product_plan", "create") && !perms.Can("product_plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		planID := viewCtx.Request.PathValue("id")
		kind := viewCtx.Request.URL.Query().Get("product_kind")
		if kind == "" {
			kind = "service"
		}
		productID := viewCtx.Request.URL.Query().Get("product_id")

		products := loadProductOptions(ctx, deps, kind)
		disabledIDs := loadExistingProductIDs(ctx, deps, planID)
		formLabels := deps.Labels.ProductPlanForm

		variantConfigurable := false
		var variantOptions []types.SelectOption
		if productID != "" {
			if lookupProductVariantMode(ctx, deps, productID) == "configurable" {
				variantConfigurable = true
				variantOptions = loadVariantOptions(ctx, deps, productID, "")
			}
		}

		return view.OK("product-picker-partial", &PickerPartialData{
			PlanID:                   planID,
			Name:                     "",
			SelectedKind:             kind,
			SelectedProductID:        productID,
			SelectedProductLabel:     findProductLabel(products, productID),
			ProductLabel:             formLabels.Product,
			ProductPlaceholder:       formLabels.ProductPlaceholder,
			ProductOptions:           buildProductAutoCompleteOptions(products, productID, disabledIDs),
			VariantConfigurable:      variantConfigurable,
			VariantOptions:           variantOptions,
			VariantSelectLabel:       formLabels.VariantSelectLabel,
			VariantSelectPlaceholder: formLabels.VariantSelectPlaceholder,
		})
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
