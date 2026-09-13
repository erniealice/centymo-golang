package block

import (
	productpkg "github.com/erniealice/centymo-golang/domain/product/product"
	priceschedulepkg "github.com/erniealice/centymo-golang/domain/subscription/price_schedule"
	compose "github.com/erniealice/espyna-golang/consumer/compose"
)

// WithInventoryCatalogMounts registers products, supplies and price schedules
// under their distinct inventory namespaces. Hosts must omit nav-only units
// with these same keys when enabling the real modules.
func WithInventoryCatalogMounts() EngineOption {
	return func(c *engineConfig) { c.inventoryCatalogMounts = true }
}

func ProductInventoryUnit(uc *UseCases, infra *Infra) compose.Unit {
	productInventoryRoutes := productpkg.DefaultInventoryRoutes()
	labels := productpkg.DefaultLabels()
	u := compose.Unit{
		Key:            "product.product_inventory",
		Routes:         &productInventoryRoutes,
		Labels:         &labels,
		LabelJSON:      compose.JSONBinding{File: "product_inventory.json", Key: "product_inventory"},
		RouteJSON:      compose.JSONBinding{File: "route.json", Key: "product_inventory"},
		RouteKeyPrefix: "product_inventory",
		Nav: compose.NavContrib{
			Permission: "product:list",
			Items: []compose.NavItem{
				{Key: "masterlist", Route: "product_inventory.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-box", Permission: "product:list", LabelKey: "active_label", IconKey: "masterlist_icon"},
				{Key: "masterlist-inactive", Route: "product_inventory.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-box", Permission: "product:list", LabelKey: "inactive_label", IconKey: "masterlist_icon"},
			},
		},
	}
	return productCatalogUnit(uc, infra, u, "inventory", []string{"stocked_good", "non_stocked_good"}, "shipped", "bulk")
}

func ProductSuppliesUnit(uc *UseCases, infra *Infra) compose.Unit {
	productSuppliesRoutes := productpkg.DefaultSuppliesRoutes()
	labels := productpkg.DefaultLabels()
	u := compose.Unit{
		Key:            "product.product_supplies",
		Routes:         &productSuppliesRoutes,
		Labels:         &labels,
		LabelJSON:      compose.JSONBinding{File: "product_supplies.json", Key: "product_supplies"},
		RouteJSON:      compose.JSONBinding{File: "route.json", Key: "product_supplies"},
		RouteKeyPrefix: "product_supplies",
		Nav: compose.NavContrib{
			Permission: "supplies:list",
			Items: []compose.NavItem{
				{Key: "supplies", Route: "product_supplies.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-box", Permission: "supplies:list", LabelKey: "active_label", IconKey: "masterlist_icon"},
				{Key: "supplies-inactive", Route: "product_supplies.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-box", Permission: "supplies:list", LabelKey: "inactive_label", IconKey: "masterlist_icon"},
			},
		},
	}
	return productCatalogUnit(uc, infra, u, "supplies", []string{"consumable"}, "shipped", "bulk")
}

func PriceScheduleInventoryUnit(uc *UseCases, infra *Infra) compose.Unit {
	priceScheduleInventoryRoutes := priceschedulepkg.DefaultInventoryRoutes()
	priceScheduleInventoryRoutes.ActiveNav = "inventory"
	priceScheduleInventoryRoutes.ActiveSubNav = "inventory-price-schedules-active"
	labels := priceschedulepkg.DefaultLabels()
	u := compose.Unit{
		Key:            "subscription.price_schedule_inventory",
		Routes:         &priceScheduleInventoryRoutes,
		Labels:         &labels,
		LabelJSON:      compose.JSONBinding{File: "price_schedule.json", Key: "price_schedule"},
		RouteJSON:      compose.JSONBinding{File: "route.json", Key: "price_schedule_inventory"},
		RouteKeyPrefix: "price_schedule_inventory",
		Nav: compose.NavContrib{
			Permission: "price_schedule:list",
			Items: []compose.NavItem{
				{Key: "inventory-price-schedules-active", Route: "price_schedule_inventory.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-layers", Permission: "price_schedule:list", LabelKey: "active_label", IconKey: "plans_active_icon"},
				{Key: "inventory-price-schedules-inactive", Route: "price_schedule_inventory.list", Params: map[string]string{"status": "inactive"}, Label: "Inactive", Icon: "icon-layers", Permission: "price_schedule:list", LabelKey: "inactive_label", IconKey: "plans_inactive_icon"},
			},
		},
	}
	return priceScheduleCatalogUnit(uc, infra, u)
}
