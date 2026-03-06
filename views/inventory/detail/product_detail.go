package detail

import (
	"context"
	"fmt"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	inventoryitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/inventory/inventory_item"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
)

// ProductDetailDeps holds dependencies for the product-context inventory detail.
type ProductDetailDeps struct {
	InventoryRoutes centymo.InventoryRoutes
	ProductRoutes   centymo.ProductRoutes
	ReadInventoryItem func(ctx context.Context, req *inventoryitempb.ReadInventoryItemRequest) (*inventoryitempb.ReadInventoryItemResponse, error)
	ReadProduct       func(ctx context.Context, req *productpb.ReadProductRequest) (*productpb.ReadProductResponse, error)
	// Delegate to main Deps for tab data loading
	DetailDeps   *Deps
	Labels       centymo.InventoryLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// ProductDetailPageData extends PageData with product context.
type ProductDetailPageData struct {
	PageData
	ProductID   string
	ProductName string
	Breadcrumbs []Breadcrumb
}

// Breadcrumb represents a single breadcrumb navigation item.
type Breadcrumb struct {
	Label string
	Href  string
}

// NewProductDetailView creates the product-centric inventory detail view.
// Route: /app/product/detail/{pid}/inventory/detail/{iid}
func NewProductDetailView(deps *ProductDetailDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("pid")
		itemID := viewCtx.Request.PathValue("iid")

		// Load inventory item
		itemResp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
			Data: &inventoryitempb.InventoryItem{Id: itemID},
		})
		if err != nil {
			log.Printf("Failed to read inventory_item %s: %v", itemID, err)
			return view.Error(fmt.Errorf("failed to load inventory item: %w", err))
		}
		items := itemResp.GetData()
		if len(items) == 0 {
			return view.Error(fmt.Errorf("inventory item not found"))
		}
		item := items[0]

		// Load product name for breadcrumb
		productName := deps.Labels.Breadcrumb.Product
		if productID != "" {
			prodResp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
				Data: &productpb.Product{Id: productID},
			})
			if err == nil && len(prodResp.GetData()) > 0 {
				name := prodResp.GetData()[0].GetName()
				if name != "" {
					productName = name
				}
			}
		}

		name := item.GetName()
		locationID := item.GetLocationId()
		locationName := centymo.LocationDisplayName(locationID)

		activeTab := viewCtx.QueryParams["tab"]
		if activeTab == "" {
			activeTab = "info"
		}

		l := deps.Labels
		itemType := item.GetProduct().GetItemType()
		if itemType == "" {
			itemType = "non_serialized"
		}
		isSerialized := itemType == "serialized"

		// 3 tabs only: Info, Depreciation, Audit
		base := route.ResolveURL(deps.InventoryRoutes.ProductDetailURL, "pid", productID, "iid", itemID)
		action := route.ResolveURL(deps.InventoryRoutes.ProductTabActionURL, "pid", productID, "iid", itemID, "tab", "")
		tabItems := []pyeza.TabItem{
			{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
			{Key: "depreciation", Label: l.Tabs.Depreciation, Href: base + "?tab=depreciation", HxGet: action + "depreciation", Icon: "icon-trending-down", Count: 0, Disabled: false},
			{Key: "audit", Label: l.Tabs.Audit, Href: base + "?tab=audit", HxGet: action + "audit", Icon: "icon-clock", Count: 0, Disabled: false},
		}

		available := computeAvailable(item.GetQuantityOnHand(), item.GetQuantityReserved())
		itemMap := inventoryItemToMap(item)

		breadcrumbs := []Breadcrumb{
			{Label: l.Breadcrumb.Products, Href: route.ResolveURL(deps.ProductRoutes.ListURL, "status", "active")},
			{Label: productName, Href: route.ResolveURL(deps.ProductRoutes.DetailURL, "id", productID)},
			{Label: name, Href: ""},
		}

		pageData := &ProductDetailPageData{
			PageData: PageData{
				PageData: types.PageData{
					CacheVersion:   viewCtx.CacheVersion,
					Title:          name,
					CurrentPath:    viewCtx.CurrentPath,
					ActiveNav:      "inventory",
					HeaderTitle:    name,
					HeaderSubtitle: locationName,
					HeaderIcon:     "icon-package",
					CommonLabels:   deps.CommonLabels,
				},
				ContentTemplate: "inventory-product-detail-content",
				Item:            itemMap,
				Labels:          l,
				ActiveTab:       activeTab,
				TabItems:        tabItems,
				IsSerialized:    isSerialized,
				ItemType:        itemType,
				ItemTypeLabel:   itemTypeDisplayLabel(itemType, l),
				ItemTypeVariant: itemTypeDisplayVariant(itemType),
				LocationName:    locationName,
				AvailableQty:    available,
			},
			ProductID:   productID,
			ProductName: productName,
			Breadcrumbs: breadcrumbs,
		}

		// Load tab-specific data
		dd := deps.DetailDeps
		switch activeTab {
		case "info":
			// Load option values (attributes) for display
			pageData.Attributes = loadAttributes(ctx, dd, item)
			// For serialized items, embed serial table in info tab
			if isSerialized {
				perms := view.GetUserPermissions(ctx)
				serials := loadSerials(ctx, dd, itemID)
				pageData.SerialTable = buildSerialTable(serials, l, deps.TableLabels, itemID, deps.InventoryRoutes, perms)
				pageData.SerialSummary = computeSerialSummary(serials)
			}
		case "depreciation":
			pageData.Depreciation = loadDepreciation(ctx, dd, itemID, l)
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		}

		return view.OK("inventory-product-detail", pageData)
	})
}

// NewProductDetailTabAction handles HTMX tab switching for the product-context detail.
// Route: /action/product/{pid}/inventory/{iid}/tab/{tab}
func NewProductDetailTabAction(deps *ProductDetailDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("pid")
		itemID := viewCtx.Request.PathValue("iid")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		itemResp, err := deps.ReadInventoryItem(ctx, &inventoryitempb.ReadInventoryItemRequest{
			Data: &inventoryitempb.InventoryItem{Id: itemID},
		})
		if err != nil {
			return view.Error(fmt.Errorf("failed to load inventory item: %w", err))
		}
		items := itemResp.GetData()
		if len(items) == 0 {
			return view.Error(fmt.Errorf("inventory item not found"))
		}
		item := items[0]

		// Load product name
		productName := deps.Labels.Breadcrumb.Product
		if productID != "" {
			prodResp, err := deps.ReadProduct(ctx, &productpb.ReadProductRequest{
				Data: &productpb.Product{Id: productID},
			})
			if err == nil && len(prodResp.GetData()) > 0 {
				name := prodResp.GetData()[0].GetName()
				if name != "" {
					productName = name
				}
			}
		}

		l := deps.Labels
		itemType := item.GetProduct().GetItemType()
		if itemType == "" {
			itemType = "non_serialized"
		}
		isSerialized := itemType == "serialized"
		name := item.GetName()
		itemMap := inventoryItemToMap(item)

		pageData := &ProductDetailPageData{
			PageData: PageData{
				Item:            itemMap,
				Labels:          l,
				ActiveTab:       tab,
				IsSerialized:    isSerialized,
				ItemType:        itemType,
				ItemTypeLabel:   itemTypeDisplayLabel(itemType, l),
				ItemTypeVariant: itemTypeDisplayVariant(itemType),
				LocationName:    centymo.LocationDisplayName(item.GetLocationId()),
				AvailableQty:    computeAvailable(item.GetQuantityOnHand(), item.GetQuantityReserved()),
			},
			ProductID:   productID,
			ProductName: productName,
			Breadcrumbs: []Breadcrumb{
				{Label: l.Breadcrumb.Products, Href: route.ResolveURL(deps.ProductRoutes.ListURL, "status", "active")},
				{Label: productName, Href: route.ResolveURL(deps.ProductRoutes.DetailURL, "id", productID)},
				{Label: name, Href: ""},
			},
		}

		dd := deps.DetailDeps
		switch tab {
		case "info":
			pageData.Attributes = loadAttributes(ctx, dd, item)
			if isSerialized {
				perms := view.GetUserPermissions(ctx)
				serials := loadSerials(ctx, dd, itemID)
				pageData.SerialTable = buildSerialTable(serials, l, deps.TableLabels, itemID, deps.InventoryRoutes, perms)
				pageData.SerialSummary = computeSerialSummary(serials)
			}
		case "depreciation":
			pageData.Depreciation = loadDepreciation(ctx, dd, itemID, l)
		case "audit":
			pageData.AuditTable = buildAuditTable(l, deps.TableLabels)
		}

		templateName := "inventory-product-tab-" + tab
		return view.OK(templateName, pageData)
	})
}
