package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
)

// VariantFormLabels holds labels for the variant drawer form template.
type VariantFormLabels struct {
	SKU           string
	PriceOverride string
	Attributes    string
}

// AttributeOption represents an attribute available for selection on a variant.
type AttributeOption struct {
	Code  string
	Name  string
	Value string
}

// VariantFormData is the template data for the variant drawer form.
type VariantFormData struct {
	FormAction       string
	IsEdit           bool
	ID               string
	ProductID        string
	SKU              string
	PriceOverride    string
	Active           bool
	Labels           VariantFormLabels
	CommonLabels     any
	AttributeOptions []AttributeOption
}

// VariantDeps holds dependencies for variant action handlers.
type VariantDeps struct {
	DB           centymo.DataSource
	Labels       centymo.ProductLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// NewVariantsTableView returns a view that renders only the variants table (for HTMX refresh).
func NewVariantsTableView(deps *VariantDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")

		detailDeps := &Deps{
			DB:          deps.DB,
			Labels:      deps.Labels,
			TableLabels: deps.TableLabels,
		}

		tableConfig := buildVariantsTable(ctx, detailDeps, productID)
		return view.OK("table-card", tableConfig)
	})
}

// NewVariantAssignView creates the variant assign action (GET = form, POST = create).
func NewVariantAssignView(deps *VariantDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			l := deps.Labels
			return view.OK("variant-drawer-form", &VariantFormData{
				FormAction: fmt.Sprintf("/action/products/detail/%s/variants/assign", productID),
				ProductID:  productID,
				Active:     true,
				Labels: VariantFormLabels{
					SKU:           l.Variant.SKU,
					PriceOverride: l.Variant.PriceOverride,
					Attributes:    l.Variant.Attributes,
				},
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create variant
		if err := viewCtx.Request.ParseForm(); err != nil {
			return htmxError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		data := map[string]any{
			"product_id":     productID,
			"sku":            r.FormValue("sku"),
			"price_override": r.FormValue("price_override"),
			"active":         active,
		}

		_, err := deps.DB.Create(ctx, "product_variant", data)
		if err != nil {
			log.Printf("Failed to create product variant: %v", err)
			return htmxError("Failed to create variant")
		}

		return htmxSuccess("product-variants-table")
	})
}

// NewVariantEditView creates the variant edit action (GET = form, POST = update).
func NewVariantEditView(deps *VariantDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")
		variantID := viewCtx.Request.PathValue("vid")

		if viewCtx.Request.Method == http.MethodGet {
			record, err := deps.DB.Read(ctx, "product_variant", variantID)
			if err != nil {
				log.Printf("Failed to read product variant %s: %v", variantID, err)
				return htmxError("Variant not found")
			}

			sku, _ := record["sku"].(string)
			priceOverride := ""
			if po, ok := record["price_override"]; ok && po != nil {
				priceOverride = fmt.Sprintf("%v", po)
			}
			active, _ := record["active"].(bool)

			l := deps.Labels
			return view.OK("variant-drawer-form", &VariantFormData{
				FormAction: fmt.Sprintf("/action/products/detail/%s/variants/edit/%s", productID, variantID),
				IsEdit:     true,
				ID:         variantID,
				ProductID:  productID,
				SKU:        sku,
				PriceOverride: priceOverride,
				Active:     active,
				Labels: VariantFormLabels{
					SKU:           l.Variant.SKU,
					PriceOverride: l.Variant.PriceOverride,
					Attributes:    l.Variant.Attributes,
				},
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — update variant
		if err := viewCtx.Request.ParseForm(); err != nil {
			return htmxError("Invalid form data")
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"

		data := map[string]any{
			"sku":            r.FormValue("sku"),
			"price_override": r.FormValue("price_override"),
			"active":         active,
		}

		_, err := deps.DB.Update(ctx, "product_variant", variantID, data)
		if err != nil {
			log.Printf("Failed to update product variant %s: %v", variantID, err)
			return htmxError("Failed to update variant")
		}

		return htmxSuccess("product-variants-table")
	})
}

// NewVariantRemoveView creates the variant remove action (POST only, with dialog confirmation).
func NewVariantRemoveView(deps *VariantDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return htmxError("Variant ID is required")
		}

		err := deps.DB.Delete(ctx, "product_variant", id)
		if err != nil {
			log.Printf("Failed to delete product variant %s: %v", id, err)
			return htmxError("Failed to remove variant")
		}

		return htmxSuccess("product-variants-table")
	})
}

// ---------------------------------------------------------------------------
// HTMX response helpers (local to detail package)
// ---------------------------------------------------------------------------

func htmxSuccess(tableID string) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"HX-Trigger": fmt.Sprintf(`{"formSuccess":true,"refreshTable":"%s"}`, tableID),
		},
	}
}

func htmxError(message string) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusUnprocessableEntity,
		Headers: map[string]string{
			"HX-Error-Message": message,
		},
	}
}
