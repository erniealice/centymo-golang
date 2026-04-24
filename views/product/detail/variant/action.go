package variant

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	detail "github.com/erniealice/centymo-golang/views/product/detail"

	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
)

// NewTableView returns a view that renders only the variants table (for HTMX refresh).
func NewTableView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")

		detailDeps := &detail.DetailViewDeps{
			Routes:                    deps.Routes,
			DB:                        deps.DB,
			Labels:                    deps.Labels,
			TableLabels:               deps.TableLabels,
			ListProductVariants:       deps.ListProductVariants,
			ListProductVariantOptions: deps.ListProductVariantOptions,
			ListProductOptionValues:   deps.ListProductOptionValues,
		}

		perms := view.GetUserPermissions(ctx)
		tableConfig := detail.BuildVariantsTable(ctx, detailDeps, productID, perms)
		return view.OK("table-card", tableConfig)
	})
}

// NewAssignView creates the variant assign action (GET = form, POST = create).
func NewAssignView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(deps.permEntity(), "create") {
			return detail.HtmxError(deps.Labels.Errors.PermissionDenied)
		}

		productID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			l := deps.Labels
			return view.OK("variant-drawer-form", &VariantFormData{
				FormAction: route.ResolveURL(deps.Routes.VariantAssignURL, "id", productID),
				ProductID:  productID,
				Active:     true,
				Labels: VariantFormLabels{
					SKU:           l.Variant.SKU,
					PriceOverride: l.Variant.PriceOverride,
				},
				CommonLabels:     nil, // injected by ViewAdapter
				OptionSelections: loadOptionSelections(ctx, deps, productID),
			})
		}

		// POST — create variant
		if err := viewCtx.Request.ParseForm(); err != nil {
			return detail.HtmxError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		optionSels := loadOptionSelections(ctx, deps, productID)
		if msg := validateRequiredOptions(optionSels, r.Form); msg != "" {
			return detail.HtmxError(msg)
		}

		active := r.FormValue("active") == "true"

		var priceOverride int64
		if v := r.FormValue("price_override"); v != "" {
			f, _ := strconv.ParseFloat(v, 64)
			priceOverride = int64(math.Round(f * 100))
		}

		resp, err := deps.CreateProductVariant(ctx, &productvariantpb.CreateProductVariantRequest{
			Data: &productvariantpb.ProductVariant{
				ProductId:     productID,
				Sku:           r.FormValue("sku"),
				PriceOverride: priceOverride,
				Active:        active,
			},
		})
		if err != nil {
			log.Printf("Failed to create product variant: %v", err)
			return detail.HtmxError(err.Error())
		}

		if created := resp.GetData(); len(created) > 0 {
			variantID := created[0].GetId()
			if variantID != "" {
				saveVariantOptions(ctx, deps, variantID, r.Form)
			}
		}

		return detail.HtmxSuccess("product-variants-table")
	})
}

// NewEditView creates the variant edit action (GET = form, POST = update).
func NewEditView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(deps.permEntity(), "update") {
			return detail.HtmxError(deps.Labels.Errors.PermissionDenied)
		}

		productID := viewCtx.Request.PathValue("id")
		variantID := viewCtx.Request.PathValue("vid")

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadProductVariant(ctx, &productvariantpb.ReadProductVariantRequest{
				Data: &productvariantpb.ProductVariant{Id: variantID},
			})
			if err != nil || len(readResp.GetData()) == 0 {
				log.Printf("Failed to read product variant %s: %v", variantID, err)
				return detail.HtmxError(deps.Labels.Errors.NotFound)
			}
			record := readResp.GetData()[0]

			sku := record.GetSku()
			priceOverride := ""
			if po := record.GetPriceOverride(); po != 0 {
				priceOverride = fmt.Sprintf("%v", po)
			}
			active := record.GetActive()

			l := deps.Labels
			optionSels := loadOptionSelections(ctx, deps, productID)
			existing := loadVariantOptionSelections(ctx, deps, variantID)
			for i := range optionSels {
				if sel, ok := existing[optionSels[i].OptionID]; ok {
					optionSels[i].Selected = sel
				}
			}
			return view.OK("variant-drawer-form", &VariantFormData{
				FormAction:    route.ResolveURL(deps.Routes.VariantEditURL, "id", productID, "vid", variantID),
				IsEdit:        true,
				ID:            variantID,
				ProductID:     productID,
				SKU:           sku,
				PriceOverride: priceOverride,
				Active:        active,
				Labels: VariantFormLabels{
					SKU:           l.Variant.SKU,
					PriceOverride: l.Variant.PriceOverride,
				},
				CommonLabels:     nil, // injected by ViewAdapter
				OptionSelections: optionSels,
			})
		}

		// POST — update variant
		if err := viewCtx.Request.ParseForm(); err != nil {
			return detail.HtmxError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request

		optionSels := loadOptionSelections(ctx, deps, productID)
		if msg := validateRequiredOptions(optionSels, r.Form); msg != "" {
			return detail.HtmxError(msg)
		}

		active := r.FormValue("active") == "true"

		var editPriceOverride int64
		if v := r.FormValue("price_override"); v != "" {
			f, _ := strconv.ParseFloat(v, 64)
			editPriceOverride = int64(math.Round(f * 100))
		}

		_, err := deps.UpdateProductVariant(ctx, &productvariantpb.UpdateProductVariantRequest{
			Data: &productvariantpb.ProductVariant{
				Id:            variantID,
				Sku:           r.FormValue("sku"),
				PriceOverride: editPriceOverride,
				Active:        active,
			},
		})
		if err != nil {
			log.Printf("Failed to update product variant %s: %v", variantID, err)
			return detail.HtmxError(err.Error())
		}

		deleteVariantOptions(ctx, deps, variantID)
		saveVariantOptions(ctx, deps, variantID, r.Form)

		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": route.ResolveURL(deps.Routes.VariantDetailURL, "id", productID, "vid", variantID),
			},
		}
	})
}

// NewRemoveView creates the variant remove action (POST only, with dialog confirmation).
func NewRemoveView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can(deps.permEntity(), "delete") {
			return detail.HtmxError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return detail.HtmxError(deps.Labels.Errors.IDRequired)
		}

		deleteVariantOptions(ctx, deps, id)

		_, err := deps.DeleteProductVariant(ctx, &productvariantpb.DeleteProductVariantRequest{
			Data: &productvariantpb.ProductVariant{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete product variant %s: %v", id, err)
			return detail.HtmxError(err.Error())
		}

		return detail.HtmxSuccess("product-variants-table")
	})
}
