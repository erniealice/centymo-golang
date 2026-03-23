package serial

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
	detail "github.com/erniealice/centymo-golang/views/product/detail"
	"github.com/erniealice/centymo-golang/views/product/detail/variant"
	variantitem "github.com/erniealice/centymo-golang/views/product/detail/variant/item"

	inventoryserialpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_serial"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	productvariantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant"
)

// SerialDetailPageData holds data for the inventory serial detail page.
type SerialDetailPageData struct {
	types.PageData
	ContentTemplate   string
	Breadcrumbs       []detail.Breadcrumb
	ProductID         string
	VariantID         string
	InventoryItemID   string
	SerialID          string
	// Serial fields
	SerialNumber      string
	IMEI              string
	Status            string
	StatusVariant     string
	WarrantyEnd       string
	PurchaseOrder     string
	Labels            centymo.ProductLabels
}

// NewPageView creates the inventory serial detail view (full page).
// Route: /app/products/detail/{id}/variant/{vid}/stock/{iid}/serial/{sid}
func NewPageView(deps *variant.DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")
		variantID := viewCtx.Request.PathValue("vid")
		itemID := viewCtx.Request.PathValue("iid")
		serialID := viewCtx.Request.PathValue("sid")

		// Load product for breadcrumb
		prodResp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
			Data: &productpb.Product{Id: productID},
		})
		if err != nil || len(prodResp.GetData()) == 0 {
			log.Printf("Failed to read product %s: %v", productID, err)
			return view.Error(fmt.Errorf("failed to load product: %w", err))
		}
		product := prodResp.GetData()[0]
		productName := product.GetName()

		// Load variant for breadcrumb
		varResp, err := deps.ReadProductVariant(ctx, &productvariantpb.ReadProductVariantRequest{
			Data: &productvariantpb.ProductVariant{Id: variantID},
		})
		if err != nil || len(varResp.GetData()) == 0 {
			log.Printf("Failed to read product_variant %s: %v", variantID, err)
			return view.Error(fmt.Errorf("failed to load variant: %w", err))
		}
		variantSKU := varResp.GetData()[0].GetSku()

		// Load inventory item for breadcrumb
		item, err := variantitem.ReadInventoryItem(ctx, deps, itemID)
		if err != nil {
			return view.Error(err)
		}
		itemName := item.GetName()
		locationName := centymo.LocationDisplayName(item.GetLocationId())

		// Load serial
		serial, err := readSerial(ctx, deps, itemID, serialID)
		if err != nil {
			return view.Error(err)
		}

		serialNumber := serial.GetSerialNumber()
		headerTitle := serialNumber

		l := deps.Labels

		breadcrumbs := []detail.Breadcrumb{
			{Label: l.Breadcrumb.Products, Href: route.ResolveURL(deps.Routes.ListURL, "status", "active")},
			{Label: productName, Href: route.ResolveURL(deps.Routes.DetailURL, "id", productID) + "?tab=variants"},
			{Label: variantSKU, Href: route.ResolveURL(deps.Routes.VariantDetailURL, "id", productID, "vid", variantID) + "?tab=stock"},
			{Label: itemName + " @ " + locationName, Href: route.ResolveURL(deps.Routes.VariantStockDetailURL, "id", productID, "vid", variantID, "iid", itemID) + "?tab=serials"},
			{Label: serialNumber, Href: ""},
		}

		status := serial.GetStatus()

		pageData := &SerialDetailPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    headerTitle,
				HeaderSubtitle: serial.GetImei(),
				HeaderIcon:     "icon-hash",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "serial-detail-content",
			Breadcrumbs:     breadcrumbs,
			ProductID:       productID,
			VariantID:       variantID,
			InventoryItemID: itemID,
			SerialID:        serialID,
			SerialNumber:    serialNumber,
			IMEI:            serial.GetImei(),
			Status:          status,
			StatusVariant:   variantitem.SerialStatusVariant(status),
			WarrantyEnd:     serial.GetWarrantyEnd(),
			PurchaseOrder:   serial.GetPurchaseOrder(),
			Labels:          l,
		}

		return view.OK("serial-detail", pageData)
	})
}

// readSerial reads a single serial by ID from the inventory item's serials list.
func readSerial(ctx context.Context, deps *variant.DetailViewDeps, itemID, serialID string) (*inventoryserialpb.InventorySerial, error) {
	serials := variantitem.LoadSerials(ctx, deps, itemID)
	for _, s := range serials {
		if s.GetId() == serialID {
			return s, nil
		}
	}
	return nil, fmt.Errorf("serial not found: %s", serialID)
}
