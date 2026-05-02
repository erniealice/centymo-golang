package action

import (
	"context"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	productform "github.com/erniealice/centymo-golang/views/product/form"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	linepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/line"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productoptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_option"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
)

// Deps holds dependencies for product action handlers.
type Deps struct {
	Routes           centymo.ProductRoutes
	Labels           centymo.ProductLabels
	CreateProduct    func(ctx context.Context, req *productpb.CreateProductRequest) (*productpb.CreateProductResponse, error)
	ReadProduct      func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	UpdateProduct    func(ctx context.Context, req *productpb.UpdateProductRequest) (*productpb.UpdateProductResponse, error)
	DeleteProduct    func(ctx context.Context, req *productpb.DeleteProductRequest) (*productpb.DeleteProductResponse, error)
	SetProductActive func(ctx context.Context, id string, active bool) error
	ListLines        func(ctx context.Context, req *linepb.ListLinesRequest) (*linepb.ListLinesResponse, error)
	// ListProductOptions / ListProductVariants — used by the edit handler to
	// decide whether the variant_mode toggle should be locked (Bug 4: a product
	// that already has options or variants cannot be flipped back to "none"
	// without deleting its children first).
	ListProductOptions  func(ctx context.Context, req *productoptionpb.ListProductOptionsRequest) (*productoptionpb.ListProductOptionsResponse, error)
	ListProductVariants func(ctx context.Context, req *productvariantpb.ListProductVariantsRequest) (*productvariantpb.ListProductVariantsResponse, error)
	// DefaultProductKind is applied when the form does not supply a product_kind value.
	// Set to "service" for the service UI mount so new products appear in the services list.
	DefaultProductKind string
	// DefaultDeliveryMode is applied when the form does not supply a delivery_mode value.
	DefaultDeliveryMode string
	// DefaultTrackingMode is applied when the form does not supply a tracking_mode value.
	DefaultTrackingMode string

	// Allowed*: per-mount restriction on which enum values the drawer exposes
	// as <select> options. Nil or empty means "all values". The services mount
	// passes AllowedProductKinds=["service"] so the product_kind select is
	// rendered single-option-disabled; the supplies mount passes ["consumable"];
	// inventory passes ["stocked_good","non_stocked_good"]. DeliveryMode and
	// TrackingMode stay open to all values per mount today — the list is here
	// for future narrowing if a mount needs it.
	AllowedProductKinds  []string
	AllowedDeliveryModes []string
	AllowedTrackingModes []string

	// PermissionEntity is the first argument to perms.Can(entity, action).
	// Default "product" when empty. Mounts override: services → "service",
	// supplies → "supplies". See ModuleDeps.PermissionEntity doc for rationale.
	PermissionEntity string
}

// permEntity returns the configured PermissionEntity with a safe default so
// callers never need nil-guard. Kept private to keep the fallback rule in
// exactly one place.
func (d *Deps) permEntity() string {
	if d == nil || d.PermissionEntity == "" {
		return "product"
	}
	return d.PermissionEntity
}

// Canonical enum universes — source of truth for the drawer-form select
// options. Ordered to match the lyngua productKind/deliveryMode/trackingMode
// blocks so the UI reads left-to-right the same way.
var (
	allProductKinds  = []string{"service", "stocked_good", "non_stocked_good", "consumable"}
	allDeliveryModes = []string{"instant", "scheduled", "shipped", "digital", "project", "subscription"}
	allTrackingModes = []string{"none", "bulk", "serialized"}
)

// buildEnumOptions filters the given enum universe to just the allowed
// values (or returns every value when allowed is empty/nil), attaching the
// per-value label from labelMap. Per-value help text lives in infoMap and
// surfaces via the <option>'s Description field — form-group.html propagates
// it as data-description so a future per-option tooltip can read it without
// a schema change. The Selected flag marks the current value.
func buildEnumOptions(universe, allowed []string, current string, labelMap, infoMap map[string]string) []types.SelectOption {
	keep := universe
	if len(allowed) > 0 {
		keep = allowed
	}
	opts := make([]types.SelectOption, 0, len(keep))
	for _, v := range keep {
		label := labelMap[v]
		if label == "" {
			label = v
		}
		opts = append(opts, types.SelectOption{
			Value:       v,
			Label:       label,
			Description: infoMap[v],
			Selected:    v == current,
		})
	}
	return opts
}

// enumLabelMap extracts a value→label map from typed label structs (e.g.
// centymo.ProductKindLabels) using reflection over the struct's JSON tags.
// Done this way so the drawer form inherits the exact same tier-cascaded
// strings that show elsewhere in the UI without re-listing them here.
func productKindLabelMap(l centymo.ProductKindLabels) map[string]string {
	return map[string]string{
		"service":          l.Service,
		"stocked_good":     l.StockedGood,
		"non_stocked_good": l.NonStockedGood,
		"consumable":       l.Consumable,
	}
}

func deliveryModeLabelMap(l centymo.DeliveryModeLabels) map[string]string {
	return map[string]string{
		"instant":      l.Instant,
		"scheduled":    l.Scheduled,
		"shipped":      l.Shipped,
		"digital":      l.Digital,
		"project":      l.Project,
		"subscription": l.Subscription,
	}
}

func trackingModeLabelMap(l centymo.TrackingModeLabels) map[string]string {
	return map[string]string{
		"none":       l.None,
		"bulk":       l.Bulk,
		"serialized": l.Serialized,
	}
}

// firstNonEmpty returns the first non-empty string from its args. Used to
// pick "stored value, else mount default, else first-allowed" precedence on
// the drawer form's classifier selects.
func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// firstAllowed returns the first element of allowed (or of universe when
// allowed is empty/nil). Used to give the drawer form a safe initial value
// for a mount that hasn't been given a DefaultXxx.
func firstAllowed(allowed, universe []string) string {
	if len(allowed) > 0 {
		return allowed[0]
	}
	if len(universe) > 0 {
		return universe[0]
	}
	return ""
}

// countOptionsAndVariants returns how many product_option + product_variant
// rows reference the given product. Used to lock the variant_mode toggle +
// reject illegal transitions in the POST handler.
func countOptionsAndVariants(ctx context.Context, deps *Deps, productID string) (int, int) {
	optionCount := 0
	variantCount := 0
	if deps.ListProductOptions != nil {
		if resp, err := deps.ListProductOptions(ctx, &productoptionpb.ListProductOptionsRequest{}); err == nil {
			for _, o := range resp.GetData() {
				if o != nil && o.GetProductId() == productID {
					optionCount++
				}
			}
		}
	}
	if deps.ListProductVariants != nil {
		if resp, err := deps.ListProductVariants(ctx, &productvariantpb.ListProductVariantsRequest{}); err == nil {
			for _, v := range resp.GetData() {
				if v != nil && v.GetProductId() == productID {
					variantCount++
				}
			}
		}
	}
	return optionCount, variantCount
}

// loadLineOptions fetches all active lines and returns them as SelectOption slice.
// selectedID marks the option that should be pre-selected (for edit mode).
func loadLineOptions(ctx context.Context, deps *Deps, selectedID string) []types.SelectOption {
	if deps.ListLines == nil {
		return nil
	}
	resp, err := deps.ListLines(ctx, &linepb.ListLinesRequest{})
	if err != nil {
		log.Printf("Failed to load lines for product form: %v", err)
		return nil
	}
	options := make([]types.SelectOption, 0, len(resp.GetData()))
	for _, line := range resp.GetData() {
		if line == nil {
			continue
		}
		if !line.GetActive() && line.GetId() != selectedID {
			continue
		}
		options = append(options, types.SelectOption{
			Value:    line.GetId(),
			Label:    line.GetName(),
			Selected: line.GetId() == selectedID,
		})
	}
	return options
}

// NewAddAction creates the product add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(deps.permEntity(), "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			// Mount defaults seed the currently-selected value for each axis.
			// Falls back to the first allowed value so the select never opens
			// with an empty selection (important for disabled single-option
			// renders where the submitted value must match stored state).
			productKind := firstNonEmpty(deps.DefaultProductKind, firstAllowed(deps.AllowedProductKinds, allProductKinds))
			deliveryMode := firstNonEmpty(deps.DefaultDeliveryMode, firstAllowed(deps.AllowedDeliveryModes, allDeliveryModes))
			trackingMode := firstNonEmpty(deps.DefaultTrackingMode, firstAllowed(deps.AllowedTrackingModes, allTrackingModes))
			return view.OK("product-drawer-form", &productform.Data{
				FormAction:           deps.Routes.AddURL,
				Active:               true,
				Currency:             "PHP",
				VariantMode:          "none",
				CanToggleVariantMode: true, // new products have no options/variants yet
				LineOptions:          loadLineOptions(ctx, deps, ""),
				ProductKind:          productKind,
				ProductKindOptions:   buildEnumOptions(allProductKinds, deps.AllowedProductKinds, productKind, productKindLabelMap(deps.Labels.ProductKind), deps.Labels.Form.ProductKindValueInfo),
				DeliveryMode:         deliveryMode,
				DeliveryModeOptions:  buildEnumOptions(allDeliveryModes, deps.AllowedDeliveryModes, deliveryMode, deliveryModeLabelMap(deps.Labels.DeliveryMode), deps.Labels.Form.DeliveryModeValueInfo),
				TrackingMode:         trackingMode,
				TrackingModeOptions:  buildEnumOptions(allTrackingModes, deps.AllowedTrackingModes, trackingMode, trackingModeLabelMap(deps.Labels.TrackingMode), deps.Labels.Form.TrackingModeValueInfo),
				Labels:               deps.Labels.Form,
				CommonLabels:         nil, // injected by ViewAdapter
			})
		}

		// POST — create product
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		desc := r.FormValue("description")

		productKind := r.FormValue("product_kind")
		if productKind == "" {
			productKind = deps.DefaultProductKind
		}
		deliveryMode := r.FormValue("delivery_mode")
		if deliveryMode == "" {
			deliveryMode = deps.DefaultDeliveryMode
		}
		trackingMode := r.FormValue("tracking_mode")
		if trackingMode == "" {
			trackingMode = deps.DefaultTrackingMode
		}

		// Model D — variant_mode drives optional pricing. The hidden fallback
		// posts "none" when the toggle is unchecked; when checked the toggle
		// posts "configurable" with a later form value (Go's ParseForm keeps
		// both under the same key, last-wins).
		variantMode := r.FormValue("variant_mode")
		if variantMode == "" {
			variantMode = "none"
		}
		// Pick the last submitted value (the toggle overrides the hidden input
		// when checked).
		if vals := r.Form["variant_mode"]; len(vals) > 0 {
			variantMode = vals[len(vals)-1]
		}
		if variantMode != "configurable" {
			variantMode = "none"
		}

		productData := &productpb.Product{
			Name:         r.FormValue("name"),
			Description:  &desc,
			Currency:     r.FormValue("currency"),
			Active:       active,
			ProductKind:  productKind,
			DeliveryMode: deliveryMode,
			TrackingMode: trackingMode,
			VariantMode:  variantMode,
		}
		// Price is optional: write only when variant_mode = "none". When
		// configurable, per-variant price_override on ProductVariant is
		// authoritative and product.price should stay null.
		if variantMode == "none" {
			if priceStr := r.FormValue("price"); priceStr != "" {
				priceF, _ := strconv.ParseFloat(priceStr, 64)
				price := int64(math.Round(priceF * 100))
				productData.Price = &price
			}
		}
		if unit := r.FormValue("unit"); unit != "" {
			productData.Unit = &unit
		}
		if lineID := r.FormValue("line_id"); lineID != "" {
			productData.LineId = &lineID
		}

		_, err := deps.CreateProduct(ctx, &productpb.CreateProductRequest{
			Data: productData,
		})
		if err != nil {
			log.Printf("Failed to create product: %v", err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("products-table")
	})
}

// NewEditAction creates the product edit action (GET = form, POST = update).
// When the GET request includes ?clone=1, the handler returns the drawer form
// pre-populated from the source record but wired to AddURL (submission creates
// a new product) with " (Copy)" appended to the name.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can(deps.permEntity(), requiredAction) {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
				Data: &productpb.Product{Id: id},
			})
			if err != nil {
				log.Printf("Failed to read product %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}

			data := resp.GetData()
			if len(data) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			p := data[0]

			currentLineID := p.GetLineId()
			name := p.GetName()
			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id)
			formID := id
			if isClone {
				name = strings.TrimSpace(name) + viewCtx.T("actions.copySuffix")
				formAction = deps.Routes.AddURL
				formID = ""
			}
			// VariantMode defaults to "none" for legacy rows that predate Model D.
			variantMode := p.GetVariantMode()
			if variantMode == "" {
				variantMode = "none"
			}
			// Price is optional on the proto — only render a value when set.
			priceStr := ""
			if p.Price != nil {
				priceStr = strconv.FormatFloat(float64(p.GetPrice())/100.0, 'f', 2, 64)
			}
			// Bug 4: lock the toggle when the product already has options or
			// variants. Clone always starts fresh, so permit toggling on clone.
			canToggle := true
			if !isClone {
				optionCount, variantCount := countOptionsAndVariants(ctx, deps, id)
				canToggle = optionCount == 0 && variantCount == 0
			}
			// Classifier selects — current value prefers the stored product,
			// falls back to mount default, then to first-allowed. buildEnumOptions
			// filters the universe by the per-mount Allowed* list so the services
			// mount renders product_kind with a single disabled "Service" option,
			// the supplies mount with a single disabled "Consumable", and the
			// inventory mount with {stocked_good, non_stocked_good}.
			productKind := firstNonEmpty(p.GetProductKind(), deps.DefaultProductKind, firstAllowed(deps.AllowedProductKinds, allProductKinds))
			deliveryMode := firstNonEmpty(p.GetDeliveryMode(), deps.DefaultDeliveryMode, firstAllowed(deps.AllowedDeliveryModes, allDeliveryModes))
			trackingMode := firstNonEmpty(p.GetTrackingMode(), deps.DefaultTrackingMode, firstAllowed(deps.AllowedTrackingModes, allTrackingModes))
			return view.OK("product-drawer-form", &productform.Data{
				FormAction:           formAction,
				IsEdit:               !isClone,
				ID:                   formID,
				Name:                 name,
				Description:          p.GetDescription(),
				Price:                priceStr,
				Currency:             p.GetCurrency(),
				Active:               p.GetActive(),
				VariantMode:          variantMode,
				Unit:                 p.GetUnit(),
				CanToggleVariantMode: canToggle,
				LineID:               currentLineID,
				LineOptions:          loadLineOptions(ctx, deps, currentLineID),
				ProductKind:          productKind,
				ProductKindOptions:   buildEnumOptions(allProductKinds, deps.AllowedProductKinds, productKind, productKindLabelMap(deps.Labels.ProductKind), deps.Labels.Form.ProductKindValueInfo),
				DeliveryMode:         deliveryMode,
				DeliveryModeOptions:  buildEnumOptions(allDeliveryModes, deps.AllowedDeliveryModes, deliveryMode, deliveryModeLabelMap(deps.Labels.DeliveryMode), deps.Labels.Form.DeliveryModeValueInfo),
				TrackingMode:         trackingMode,
				TrackingModeOptions:  buildEnumOptions(allTrackingModes, deps.AllowedTrackingModes, trackingMode, trackingModeLabelMap(deps.Labels.TrackingMode), deps.Labels.Form.TrackingModeValueInfo),
				Labels:               deps.Labels.Form,
				CommonLabels:         nil, // injected by ViewAdapter
			})
		}

		// POST — update product
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		desc := r.FormValue("description")

		// variant_mode: prefer the last posted value (toggle overrides hidden
		// fallback when checked; otherwise the fallback's "none" wins).
		// Bug 3 hardening: if neither the hidden nor the toggle was submitted
		// (e.g. toggle rendered as disabled and stripped from the form), fall
		// through to the stored value so the UPDATE does not silently clobber
		// the column with the form-group default.
		var variantMode string
		if vals := r.Form["variant_mode"]; len(vals) > 0 {
			variantMode = vals[len(vals)-1]
		}
		if variantMode != "configurable" && variantMode != "none" {
			// No submission — retain the stored value. Read the product fresh
			// instead of trusting client-provided defaults.
			if currentResp, cerr := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
				Data: &productpb.Product{Id: id},
			}); cerr == nil && len(currentResp.GetData()) > 0 {
				variantMode = currentResp.GetData()[0].GetVariantMode()
			}
			if variantMode == "" {
				variantMode = "none"
			}
		}

		// Bug 4: reject illegal transitions. When the product already has
		// options/variants, variant_mode is locked — only accept a submission
		// that matches the stored value. Don't rely on UI disabling alone.
		optionCount, variantCount := countOptionsAndVariants(ctx, deps, id)
		if optionCount > 0 || variantCount > 0 {
			stored := ""
			if currentResp, cerr := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
				Data: &productpb.Product{Id: id},
			}); cerr == nil && len(currentResp.GetData()) > 0 {
				stored = currentResp.GetData()[0].GetVariantMode()
			}
			if stored == "" {
				stored = "configurable" // legacy rows with options default to configurable for the lock check
			}
			if variantMode != stored {
				msg := deps.Labels.Form.VariantModeLockedError
				if msg == "" {
					msg = "Delete all variants and options before changing the variant mode."
				}
				return centymo.HTMXError(msg)
			}
		}

		// Read the three taxonomy axes from the submitted form, falling back
		// to the stored value when a disabled select was stripped from the
		// submission. Never fall through to the mount default — that would
		// silently migrate a product's classification every edit if a legacy
		// row ever lands under a mount whose default differs.
		productKind := r.FormValue("product_kind")
		deliveryMode := r.FormValue("delivery_mode")
		trackingMode := r.FormValue("tracking_mode")
		if productKind == "" || deliveryMode == "" || trackingMode == "" {
			if currentResp, cerr := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
				Data: &productpb.Product{Id: id},
			}); cerr == nil && len(currentResp.GetData()) > 0 {
				stored := currentResp.GetData()[0]
				if productKind == "" {
					productKind = stored.GetProductKind()
				}
				if deliveryMode == "" {
					deliveryMode = stored.GetDeliveryMode()
				}
				if trackingMode == "" {
					trackingMode = stored.GetTrackingMode()
				}
			}
		}

		updatedProduct := &productpb.Product{
			Id:           id,
			Name:         r.FormValue("name"),
			Description:  &desc,
			Currency:     r.FormValue("currency"),
			Active:       active,
			ProductKind:  productKind,
			DeliveryMode: deliveryMode,
			TrackingMode: trackingMode,
			VariantMode:  variantMode,
		}
		if variantMode == "none" {
			if priceStr := r.FormValue("price"); priceStr != "" {
				priceF, _ := strconv.ParseFloat(priceStr, 64)
				price := int64(math.Round(priceF * 100))
				updatedProduct.Price = &price
			}
		}
		if unit := r.FormValue("unit"); unit != "" {
			updatedProduct.Unit = &unit
		}
		if lineID := r.FormValue("line_id"); lineID != "" {
			updatedProduct.LineId = &lineID
		}

		_, err := deps.UpdateProduct(ctx, &productpb.UpdateProductRequest{
			Data: updatedProduct,
		})
		if err != nil {
			log.Printf("Failed to update product %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("products-table")
	})
}

// NewDeleteAction creates the product delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(deps.permEntity(), "delete") {
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

		_, err := deps.DeleteProduct(ctx, &productpb.DeleteProductRequest{
			Data: &productpb.Product{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete product %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("products-table")
	})
}

// NewBulkDeleteAction creates the product bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(deps.permEntity(), "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
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
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(deps.permEntity(), "update") {
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
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		if err := deps.SetProductActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update product status %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("products-table")
	})
}

// NewBulkSetStatusAction creates the product bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(deps.permEntity(), "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
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
