package price_schedule

import "strings"

// PriceSchedule-domain route constants. Relocated from the subscription
// routes.go god-file (entity-local extraction). Pure structural move.
const (
	// PriceSchedule routes (date-bounded pricing container for plans)
	DashboardURL        = "/price-schedules/dashboard"
	ListURL             = "/price-schedules/list/{status}"
	TableURL            = "/action/price-schedule/table/{status}"
	DetailURL           = "/price-schedules/detail/{id}"
	AddURL              = "/action/price-schedule/add"
	EditURL             = "/action/price-schedule/edit/{id}"
	DeleteURL           = "/action/price-schedule/delete"
	BulkDeleteURL       = "/action/price-schedule/bulk-delete"
	SetStatusURL        = "/action/price-schedule/set-status"
	BulkSetStatusURL    = "/action/price-schedule/bulk-set-status"
	TabActionURL        = "/action/price-schedule/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/price-schedule/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/price-schedule/{id}/attachments/delete"
	PlanAddURL          = "/action/price-schedule/{id}/plan/add"
	// Schedule-scoped price_plan detail. Mirrors /app/price-plans/detail/{id} but nests
	// under the schedule so sidebar context stays on price-schedules (price_plan is no
	// longer a top-level sidebar entry).
	PlanDetailURL             = "/price-schedules/detail/{id}/plan/{ppid}"
	PlanTabActionURL          = "/action/price-schedule/{id}/plan/{ppid}/tab/{tab}"
	PlanEditURL               = "/action/price-schedule/{id}/plan/{ppid}/edit"
	PlanDeleteURL             = "/action/price-schedule/{id}/plan/{ppid}/delete"
	PlanProductPriceAddURL    = "/action/price-schedule/{id}/plan/{ppid}/product-prices/add"
	PlanProductPriceEditURL   = "/action/price-schedule/{id}/plan/{ppid}/product-prices/edit/{pppid}"
	PlanProductPriceDeleteURL = "/action/price-schedule/{id}/plan/{ppid}/product-prices/delete"

	// Attachments on the nested price_schedule/plan detail page.
	PlanAttachmentUploadURL = "/action/price-schedule/{id}/plan/{ppid}/attachments/upload"
	PlanAttachmentDeleteURL = "/action/price-schedule/{id}/plan/{ppid}/attachments/delete"

	// 2026-05-04 — Subscriptions tab on the schedule-scoped
	// price_plan detail page. Same handler as SubscriptionDetailURL; the
	// nested URL alone activates the rate-card → plan → subscription breadcrumb.
	// See docs/plan/20260504-price-plan-engagements-tab/.
	PlanSubscriptionDetailURL = "/price-schedules/detail/{id}/plan/{ppid}/subscription/{eid}"
)

// Routes holds all route paths for price schedule views and actions.
type Routes struct {
	ActiveNav                 string `json:"active_nav"`
	ActiveSubNav              string `json:"active_sub_nav"`
	DashboardURL              string `json:"dashboard_url"`
	ListURL                   string `json:"list_url"`
	TableURL                  string `json:"table_url"`
	DetailURL                 string `json:"detail_url"`
	AddURL                    string `json:"add_url"`
	EditURL                   string `json:"edit_url"`
	DeleteURL                 string `json:"delete_url"`
	BulkDeleteURL             string `json:"bulk_delete_url"`
	SetStatusURL              string `json:"set_status_url"`
	BulkSetStatusURL          string `json:"bulk_set_status_url"`
	TabActionURL              string `json:"tab_action_url"`
	AttachmentUploadURL       string `json:"attachment_upload_url"`
	AttachmentDeleteURL       string `json:"attachment_delete_url"`
	PlanAddURL                string `json:"plan_add_url"`
	PlanDetailURL             string `json:"plan_detail_url"`
	PlanTabActionURL          string `json:"plan_tab_action_url"`
	PlanEditURL               string `json:"plan_edit_url"`
	PlanDeleteURL             string `json:"plan_delete_url"`
	PlanProductPriceAddURL    string `json:"plan_product_price_add_url"`
	PlanProductPriceEditURL   string `json:"plan_product_price_edit_url"`
	PlanProductPriceDeleteURL string `json:"plan_product_price_delete_url"`
	PlanAttachmentUploadURL   string `json:"plan_attachment_upload_url"`
	PlanAttachmentDeleteURL   string `json:"plan_attachment_delete_url"`
	// 2026-05-04 — Subscription detail nested under the
	// schedule-scoped price_plan path. Activates the rate-card → plan →
	// subscription breadcrumb in the subscription detail view. Empty string
	// disables the nested route. See
	// docs/plan/20260504-price-plan-engagements-tab/.
	PlanSubscriptionDetailURL string `json:"plan_subscription_detail_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants defined in routes.go.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:                 "service",
		ActiveSubNav:              "price-schedules",
		DashboardURL:              DashboardURL,
		ListURL:                   ListURL,
		TableURL:                  TableURL,
		DetailURL:                 DetailURL,
		AddURL:                    AddURL,
		EditURL:                   EditURL,
		DeleteURL:                 DeleteURL,
		BulkDeleteURL:             BulkDeleteURL,
		SetStatusURL:              SetStatusURL,
		BulkSetStatusURL:          BulkSetStatusURL,
		TabActionURL:              TabActionURL,
		AttachmentUploadURL:       AttachmentUploadURL,
		AttachmentDeleteURL:       AttachmentDeleteURL,
		PlanAddURL:                PlanAddURL,
		PlanDetailURL:             PlanDetailURL,
		PlanTabActionURL:          PlanTabActionURL,
		PlanEditURL:               PlanEditURL,
		PlanDeleteURL:             PlanDeleteURL,
		PlanProductPriceAddURL:    PlanProductPriceAddURL,
		PlanProductPriceEditURL:   PlanProductPriceEditURL,
		PlanProductPriceDeleteURL: PlanProductPriceDeleteURL,
		PlanAttachmentUploadURL:   PlanAttachmentUploadURL,
		PlanAttachmentDeleteURL:   PlanAttachmentDeleteURL,
		PlanSubscriptionDetailURL: PlanSubscriptionDetailURL,
	}
}

// DefaultInventoryRoutes returns a Routes with every
// URL namespace-shifted from the services namespace onto the inventory accordion
// namespace. Used as the route base for the PriceSchedule inventory-mount
// registration in block.go; a lyngua `price_schedule_inventory` override can
// layer additional tweaks on top.
//
// Shift rules:
//   - "/app/price-schedules/"   → "/app/inventory/price-schedules/"
//   - "/action/price-schedule/" → "/action/inventory-price-schedule/"
func DefaultInventoryRoutes() Routes {
	r := DefaultRoutes()
	r.ActiveNav = "inventory"
	r.ActiveSubNav = "inventory-price-schedules-active"
	// shift matches both pre-P4 (`/app/price-schedules/*`) and post-P4
	// (`/price-schedules/*`) constant shapes — see DefaultProductInventoryRoutes
	// shift comment for the P4 regression context.
	shift := func(s string) string {
		s = strings.Replace(s, "/app/price-schedules/", "/app/inventory/price-schedules/", 1)
		s = strings.Replace(s, "/price-schedules/", "/inventory/price-schedules/", 1)
		s = strings.Replace(s, "/action/price-schedule/", "/action/inventory-price-schedule/", 1)
		return s
	}
	r.DashboardURL = shift(r.DashboardURL)
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
	r.PlanAddURL = shift(r.PlanAddURL)
	r.PlanDetailURL = shift(r.PlanDetailURL)
	r.PlanTabActionURL = shift(r.PlanTabActionURL)
	r.PlanEditURL = shift(r.PlanEditURL)
	r.PlanDeleteURL = shift(r.PlanDeleteURL)
	r.PlanProductPriceAddURL = shift(r.PlanProductPriceAddURL)
	r.PlanProductPriceEditURL = shift(r.PlanProductPriceEditURL)
	r.PlanProductPriceDeleteURL = shift(r.PlanProductPriceDeleteURL)
	r.PlanAttachmentUploadURL = shift(r.PlanAttachmentUploadURL)
	r.PlanAttachmentDeleteURL = shift(r.PlanAttachmentDeleteURL)
	r.PlanSubscriptionDetailURL = shift(r.PlanSubscriptionDetailURL)
	return r
}

// RouteMap returns a map of dot-notation keys to route paths for all
// price schedule routes.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"price_schedule.dashboard":                 r.DashboardURL,
		"price_schedule.list":                      r.ListURL,
		"price_schedule.table":                     r.TableURL,
		"price_schedule.detail":                    r.DetailURL,
		"price_schedule.add":                       r.AddURL,
		"price_schedule.edit":                      r.EditURL,
		"price_schedule.delete":                    r.DeleteURL,
		"price_schedule.bulk_delete":               r.BulkDeleteURL,
		"price_schedule.set_status":                r.SetStatusURL,
		"price_schedule.bulk_set_status":           r.BulkSetStatusURL,
		"price_schedule.tab_action":                r.TabActionURL,
		"price_schedule.attachment.upload":         r.AttachmentUploadURL,
		"price_schedule.attachment.delete":         r.AttachmentDeleteURL,
		"price_schedule.plan.add":                  r.PlanAddURL,
		"price_schedule.plan.detail":               r.PlanDetailURL,
		"price_schedule.plan.tab_action":           r.PlanTabActionURL,
		"price_schedule.plan.edit":                 r.PlanEditURL,
		"price_schedule.plan.delete":               r.PlanDeleteURL,
		"price_schedule.plan.product_price.add":    r.PlanProductPriceAddURL,
		"price_schedule.plan.product_price.edit":   r.PlanProductPriceEditURL,
		"price_schedule.plan.product_price.delete": r.PlanProductPriceDeleteURL,
		"price_schedule.plan.attachment.upload":    r.PlanAttachmentUploadURL,
		"price_schedule.plan.attachment.delete":    r.PlanAttachmentDeleteURL,
		"price_schedule.plan.subscription.detail":  r.PlanSubscriptionDetailURL,
	}
}
