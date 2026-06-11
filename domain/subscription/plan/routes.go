package plan

import "strings"

// Plan-domain route constants. Relocated from the subscription routes.go
// god-file (entity-local extraction). Pure structural move.
const (
	ListURL             = "/plans/list/{status}"
	DetailURL           = "/plans/detail/{id}"
	AddURL              = "/action/plan/add"
	EditURL             = "/action/plan/edit/{id}"
	DeleteURL           = "/action/plan/delete"
	BulkDeleteURL       = "/action/plan/bulk-delete"
	SetStatusURL        = "/action/plan/set-status"
	BulkSetStatusURL    = "/action/plan/bulk-set-status"
	TableURL            = "/action/plan/table/{status}"
	TabActionURL        = "/action/plan/detail/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/plan/detail/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/plan/detail/{id}/attachments/delete"

	// PricePlan CRUD routes (within plan context)
	PricePlanAddURL    = "/action/plan/{id}/price-plans/add"
	PricePlanEditURL   = "/action/plan/{id}/price-plans/edit/{ppid}"
	PricePlanDeleteURL = "/action/plan/{id}/price-plans/delete"

	// Plan-scoped PricePlan detail — mirrors PriceSchedulePlanDetailURL but
	// keeps users in the Package (Plan) URL namespace so ActiveNav stays
	// anchored to the Services accordion's Packages section.
	// {id}=plan id, {ppid}=price_plan id.
	PricePlanDetailURL             = "/plans/detail/{id}/price/{ppid}"
	PricePlanTabActionURL          = "/action/plan/{id}/price/{ppid}/tab/{tab}"
	PlanPricePlanEditURL           = "/action/plan/{id}/price/{ppid}/edit"
	PlanPricePlanDeleteURL         = "/action/plan/{id}/price/{ppid}/delete"
	PricePlanProductPriceAddURL    = "/action/plan/{id}/price/{ppid}/product-prices/add"
	PricePlanProductPriceEditURL   = "/action/plan/{id}/price/{ppid}/product-prices/edit/{pppid}"
	PricePlanProductPriceDeleteURL = "/action/plan/{id}/price/{ppid}/product-prices/delete"

	// ProductPlan CRUD routes (within plan context)
	ProductPlanAddURL    = "/action/plan/{id}/products/add"
	ProductPlanEditURL   = "/action/plan/{id}/products/edit/{ppid}"
	ProductPlanDeleteURL = "/action/plan/{id}/products/delete"
	ProductPlanPickerURL = "/action/plan/{id}/products/picker"
)

type Routes struct {
	// Sidebar navigation context — set via defaults or routes.json override
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL          string `json:"list_url"`
	TableURL         string `json:"table_url"`
	DetailURL        string `json:"detail_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`
	TabActionURL     string `json:"tab_action_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// PricePlan CRUD routes (within plan context)
	PricePlanAddURL    string `json:"price_plan_add_url"`
	PricePlanEditURL   string `json:"price_plan_edit_url"`
	PricePlanDeleteURL string `json:"price_plan_delete_url"`

	// Plan-scoped PricePlan detail (mirrors PriceSchedulePlanRoutes.DetailURL
	// but anchored under /app/plans/detail/{id}/price/{ppid}). Lets the package
	// detail's package-prices tab keep ActiveNav on Services > Packages instead
	// of jumping to the rate-cards namespace.
	PricePlanDetailURL             string `json:"price_plan_detail_url"`
	PricePlanTabActionURL          string `json:"price_plan_tab_action_url"`
	PricePlanProductPriceAddURL    string `json:"price_plan_product_price_add_url"`
	PricePlanProductPriceEditURL   string `json:"price_plan_product_price_edit_url"`
	PricePlanProductPriceDeleteURL string `json:"price_plan_product_price_delete_url"`

	// ProductPlan CRUD routes (within plan context)
	ProductPlanAddURL    string `json:"product_plan_add_url"`
	ProductPlanEditURL   string `json:"product_plan_edit_url"`
	ProductPlanDeleteURL string `json:"product_plan_delete_url"`
	ProductPlanPickerURL string `json:"product_plan_picker_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants defined in routes.go.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "service",
		ActiveSubNav: "plans",

		ListURL:          ListURL,
		TableURL:         TableURL,
		DetailURL:        DetailURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,
		TabActionURL:     TabActionURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		PricePlanAddURL:    PricePlanAddURL,
		PricePlanEditURL:   PricePlanEditURL,
		PricePlanDeleteURL: PricePlanDeleteURL,

		PricePlanDetailURL:             PricePlanDetailURL,
		PricePlanTabActionURL:          PricePlanTabActionURL,
		PricePlanProductPriceAddURL:    PricePlanProductPriceAddURL,
		PricePlanProductPriceEditURL:   PricePlanProductPriceEditURL,
		PricePlanProductPriceDeleteURL: PricePlanProductPriceDeleteURL,

		ProductPlanAddURL:    ProductPlanAddURL,
		ProductPlanEditURL:   ProductPlanEditURL,
		ProductPlanDeleteURL: ProductPlanDeleteURL,
		ProductPlanPickerURL: ProductPlanPickerURL,
	}
}

// DefaultBundleRoutes returns a Routes with every URL namespace-shifted
// from the services namespace onto the inventory accordion namespace. Used as
// the route base for the Plan inventory-mount registration in block.go; a lyngua
// `plan_bundle` override can layer additional tweaks on top.
//
// Bundle-mount variant of Plan routes — shifts every page + action URL from the
// services namespace (/app/plans/*) onto the inventory accordion namespace
// (/app/inventory/bundles/*). Used as the route base for the Plan
// inventory-mount registration in block.go; a lyngua `plan_bundle` override can
// layer additional tweaks on top.
//
// Shift rules:
//   - "/app/plans/"   → "/app/inventory/bundles/"
//   - "/action/plan/" → "/action/inventory-bundle/"
func DefaultBundleRoutes() Routes {
	r := DefaultRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "bundles-active"
	// shift matches both pre-P4 (`/app/plans/*`) and post-P4 (`/plans/*`)
	// constant shapes — see DefaultProductInventoryRoutes shift comment
	// for the P4 regression context.
	shift := func(s string) string {
		s = strings.Replace(s, "/app/plans/", "/app/inventory/bundles/", 1)
		s = strings.Replace(s, "/plans/", "/inventory/bundles/", 1)
		s = strings.Replace(s, "/action/plan/", "/action/inventory-bundle/", 1)
		return s
	}
	r.ListURL = shift(r.ListURL)
	r.TableURL = shift(r.TableURL)
	r.DetailURL = shift(r.DetailURL)
	r.AddURL = shift(r.AddURL)
	r.EditURL = shift(r.EditURL)
	r.DeleteURL = shift(r.DeleteURL)
	r.BulkDeleteURL = shift(r.BulkDeleteURL)
	r.SetStatusURL = shift(r.SetStatusURL)
	r.BulkSetStatusURL = shift(r.BulkSetStatusURL)
	r.TabActionURL = shift(r.TabActionURL)
	r.AttachmentUploadURL = shift(r.AttachmentUploadURL)
	r.AttachmentDeleteURL = shift(r.AttachmentDeleteURL)
	r.PricePlanAddURL = shift(r.PricePlanAddURL)
	r.PricePlanEditURL = shift(r.PricePlanEditURL)
	r.PricePlanDeleteURL = shift(r.PricePlanDeleteURL)
	r.PricePlanDetailURL = shift(r.PricePlanDetailURL)
	r.PricePlanTabActionURL = shift(r.PricePlanTabActionURL)
	r.PricePlanProductPriceAddURL = shift(r.PricePlanProductPriceAddURL)
	r.PricePlanProductPriceEditURL = shift(r.PricePlanProductPriceEditURL)
	r.PricePlanProductPriceDeleteURL = shift(r.PricePlanProductPriceDeleteURL)
	r.ProductPlanAddURL = shift(r.ProductPlanAddURL)
	r.ProductPlanEditURL = shift(r.ProductPlanEditURL)
	r.ProductPlanDeleteURL = shift(r.ProductPlanDeleteURL)
	r.ProductPlanPickerURL = shift(r.ProductPlanPickerURL)
	return r
}

// RouteMap returns a map of dot-notation keys to route paths for all
// plan routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"plan.list":            r.ListURL,
		"plan.table":           r.TableURL,
		"plan.detail":          r.DetailURL,
		"plan.add":             r.AddURL,
		"plan.edit":            r.EditURL,
		"plan.delete":          r.DeleteURL,
		"plan.bulk_delete":     r.BulkDeleteURL,
		"plan.set_status":      r.SetStatusURL,
		"plan.bulk_set_status": r.BulkSetStatusURL,
		"plan.tab_action":      r.TabActionURL,

		"plan.attachment.upload": r.AttachmentUploadURL,
		"plan.attachment.delete": r.AttachmentDeleteURL,

		"plan.pricelist.add":    r.PricePlanAddURL,
		"plan.pricelist.edit":   r.PricePlanEditURL,
		"plan.pricelist.delete": r.PricePlanDeleteURL,

		"plan.price.detail":               r.PricePlanDetailURL,
		"plan.price.tab_action":           r.PricePlanTabActionURL,
		"plan.price.product_price.add":    r.PricePlanProductPriceAddURL,
		"plan.price.product_price.edit":   r.PricePlanProductPriceEditURL,
		"plan.price.product_price.delete": r.PricePlanProductPriceDeleteURL,

		"plan.product_plan.add":    r.ProductPlanAddURL,
		"plan.product_plan.edit":   r.ProductPlanEditURL,
		"plan.product_plan.delete": r.ProductPlanDeleteURL,
		"plan.product_plan.picker": r.ProductPlanPickerURL,
	}
}
