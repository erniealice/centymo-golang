package plan

import (
	"strings"
)

// ---------------------------------------------------------------------------
// Plan labels
// ---------------------------------------------------------------------------

// FilterLabels holds translatable labels for the scope filter chip on the
// plan list page (§6.1 of the 2026-04-27 plan-client-scope plan).
type FilterLabels struct {
	ScopeChipLabel string `json:"scope_chip_label"`
	ScopeMaster    string `json:"scope_master"`
	ScopeClient    string `json:"scope_client"`
	ScopeAll       string `json:"scope_all"`
}

// Labels holds all translatable strings for the plan module.
type Labels struct {
	Page            PageLabels            `json:"page"`
	Buttons         ButtonLabels          `json:"buttons"`
	Columns         ColumnLabels          `json:"columns"`
	Empty           EmptyLabels           `json:"empty"`
	Form            FormLabels            `json:"form"`
	Actions         ActionLabels          `json:"actions"`
	Bulk            BulkLabels            `json:"bulk_actions"`
	Status          StatusLabels          `json:"status"`
	Detail          DetailLabels          `json:"detail"`
	Tabs            TabLabels             `json:"tabs"`
	Confirm         ConfirmLabels         `json:"confirm"`
	Errors          ErrorLabels           `json:"errors"`
	ProductPlanForm ProductPlanFormLabels `json:"product_plan_form"`
	Filters         FilterLabels          `json:"filters"`
}

// ---------------------------------------------------------------------------
// Plan form, detail, tabs, confirm sub-labels
// ---------------------------------------------------------------------------

type FormSectionLabels struct {
	Basic    string `json:"basic"`
	Services string `json:"services"`
}

type FormLabels struct {
	Name                string            `json:"name"`
	NamePlaceholder     string            `json:"name_placeholder"`
	Description         string            `json:"description"`
	DescPlaceholder     string            `json:"description_placeholder"`
	FulfillmentType     string            `json:"fulfillment_type"`
	Active              string            `json:"active"`
	Products            string            `json:"products"`
	ProductsPlaceholder string            `json:"products_placeholder"`
	ProductsSearch      string            `json:"products_search"`
	Sections            FormSectionLabels `json:"sections"`

	// Fulfillment type option labels
	TypeSchedule string `json:"type_schedule"`
	TypeLicense  string `json:"type_license"`
	TypeContent  string `json:"type_content"`
	TypePhysical string `json:"type_physical"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"name_info"`
	DescriptionInfo string `json:"description_info"`
	ActiveInfo      string `json:"active_info"`

	// Client-scope fields (2026-04-27 plan-client-scope plan §7).
	// Set on the Plan add/edit drawer Client picker.
	ClientLabel             string `json:"client_label"`
	ClientHelp              string `json:"client_help"`
	ClientPlaceholder       string `json:"client_placeholder"`
	ClientSearchPlaceholder string `json:"client_search_placeholder"`
	ClientNoResults         string `json:"client_no_results"`
	ClientLockedTooltip     string `json:"client_locked_tooltip"`
	ClientForLabel          string `json:"client_for_label"` // "For {{.ClientName}}" — read-only badge in client-context entry-point
	ClientInfo              string `json:"client_info"`

	// JobTemplate select (2026-04-29 auto-spawn-jobs-from-subscription plan §5
	// — Plan.job_template_id assignment from the drawer). Empty value =
	// advisory-only plan; spawn use case skips silently.
	JobTemplate     string `json:"job_template"`
	JobTemplateNone string `json:"job_template_none"`
	JobTemplateHint string `json:"job_template_hint"`

	// 2026-04-30 cyclic-subscription-jobs plan §9.3 — visits_per_cycle field.
	// Number of cycle Job instances spawned per billing cycle (default 1).
	VisitsPerCycleLabel       string `json:"visits_per_cycle_label"`
	VisitsPerCyclePlaceholder string `json:"visits_per_cycle_placeholder"`
	VisitsPerCycleHint        string `json:"visits_per_cycle_hint"`

	// Client-scope cascade notice — shown unconditionally below the client picker
	// so operators see the schedule restriction before filling other fields.
	// Tier-specific wording lives in lyngua; default uses proto-generic vocabulary.
	ClientScopeCascadeNotice string `json:"client_scope_cascade_notice"`
}

type DetailLabels struct {
	PageTitle             string `json:"page_title"`
	Price                 string `json:"price"`
	Currency              string `json:"currency"`
	Status                string `json:"status"`
	Description           string `json:"description"`
	FulfillmentType       string `json:"fulfillment_type"`
	CreatedDate           string `json:"created_date"`
	ModifiedDate          string `json:"modified_date"`
	NoProductsAssigned    string `json:"no_products_assigned"`
	NoProductsAssignedMsg string `json:"no_products_assigned_msg"`
	NoProductsDesc        string `json:"no_products_desc"`
	NoPricePlans          string `json:"no_price_plans"`
	NoPricePlansMsg       string `json:"no_price_plans_msg"`
	NoPricePlansDesc      string `json:"no_price_plans_desc"`
	AuditTrailComingSoon  string `json:"audit_trail_coming_soon"`
}

type TabLabels struct {
	Info          string `json:"info"`
	Products      string `json:"products"`
	ProductsSlug  string `json:"products_slug"`
	PricePlan     string `json:"price_plan"`
	PricePlanSlug string `json:"price_plan_slug"`
	Attachments   string `json:"attachments"`
	AuditTrail    string `json:"audit_trail"`
	AuditHistory  string `json:"audit_history"`
}

// ResolveTabSlug returns the URL slug for a canonical tab key. The "products"
// and "pricePlan" tabs can be re-slugged per tier (e.g. professional ships
// "items" and "package-prices"); other tabs round-trip through as-is.
func (t TabLabels) ResolveTabSlug(canonical string) string {
	switch canonical {
	case "products":
		if s := strings.TrimSpace(t.ProductsSlug); s != "" {
			return s
		}
	case "pricePlan":
		if s := strings.TrimSpace(t.PricePlanSlug); s != "" {
			return s
		}
	}
	return canonical
}

// CanonicalizeTab maps an incoming URL tab slug back to its canonical key so
// internal template lookups and equality checks stay tier-agnostic.
func (t TabLabels) CanonicalizeTab(slug string) string {
	if slug == "" {
		return ""
	}
	if s := strings.TrimSpace(t.ProductsSlug); s != "" && slug == s {
		return "products"
	}
	if s := strings.TrimSpace(t.PricePlanSlug); s != "" && slug == s {
		return "pricePlan"
	}
	return slug
}

type ConfirmLabels struct {
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"delete_message"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activate_message"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivate_message"`
	BulkActivate          string `json:"bulk_activate"`
	BulkActivateMessage   string `json:"bulk_activate_message"`
	BulkDeactivate        string `json:"bulk_deactivate"`
	BulkDeactivateMessage string `json:"bulk_deactivate_message"`
	BulkDelete            string `json:"bulk_delete"`
	BulkDeleteMessage     string `json:"bulk_delete_message"`
}

func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:         "Plans",
			HeadingActive:   "Active Plans",
			HeadingInactive: "Inactive Plans",
			Caption:         "Manage your plans",
			CaptionActive:   "Manage your active plans",
			CaptionInactive: "View inactive or archived plans",
		},
		Buttons: ButtonLabels{
			AddPlan:       "Add Plan",
			AddPricePlan:  "Add Price Plan",
			EditPricePlan: "Edit Price Plan",
			AddProduct:    "Add Product",
		},
		Columns: ColumnLabels{
			Name:          "Name",
			Description:   "Description",
			Interval:      "Interval",
			Price:         "Price",
			Status:        "Status",
			Product:       "Product",
			PricePlan:     "Price Plan",
			PriceSchedule: "Price Schedule",
			Duration:      "Duration",
			Location:      "Location",
			ItemType:      "Item Type",
		},
		Empty: EmptyLabels{
			Title:           "No plans found",
			Message:         "No plans to display.",
			ActiveTitle:     "No active plans",
			ActiveMessage:   "Create your first plan to get started.",
			InactiveTitle:   "No inactive plans",
			InactiveMessage: "Discontinued plans will appear here.",
		},
		Form: FormLabels{
			Name:                "Plan Name",
			NamePlaceholder:     "Enter plan name",
			Description:         "Description",
			DescPlaceholder:     "Enter plan description...",
			FulfillmentType:     "Fulfillment Type",
			Active:              "Active",
			Products:            "Products",
			ProductsPlaceholder: "Select products...",
			ProductsSearch:      "Search products...",
			TypeSchedule:        "Schedule",
			TypeLicense:         "License",
			TypeContent:         "Content",
			TypePhysical:        "Physical",
			Sections: FormSectionLabels{
				Basic:    "Basic Information",
				Services: "Assigned Products",
			},
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			NameInfo:        "Display name for this plan. Shown in subscription lists and invoices.",
			DescriptionInfo: "Optional notes about this plan. Visible on detail pages.",
			ActiveInfo:      "Inactive plans are hidden from new subscriptions.",
			// Client-scope fields (2026-04-27 plan-client-scope plan §7).
			ClientLabel:             "Client",
			ClientHelp:              "Leave blank to make this package available for any client. Set a client to make it a custom package for that client only.",
			ClientPlaceholder:       "Leave blank for a general package",
			ClientSearchPlaceholder: "Search clients...",
			ClientNoResults:         "No clients found",
			ClientLockedTooltip:     "Locked — this plan has active subscriptions. Detach them or create a new plan.",
			ClientForLabel:          "For {{.ClientName}}",
			ClientInfo:              "Optional. When set, this plan only appears for engagements with that client.",
			JobTemplate:             "Job Template",
			JobTemplateNone:         "(none — engagement has no operational tracking)",
			JobTemplateHint:         "Select the operational template that defines the work for this engagement. Leave empty for advisory-only plans.",
			// 2026-04-30 cyclic-subscription-jobs plan §9.3.
			VisitsPerCycleLabel:       "Visits per billing cycle",
			VisitsPerCyclePlaceholder: "1",
			VisitsPerCycleHint:        "Number of cycle Job instances per billing cycle. Default 1. Use 2 for biweekly visits billed monthly, 4 for weekly visits billed monthly.",
			// Client-scope cascade notice — proto-generic default; tiers override via lyngua.
			ClientScopeCascadeNotice: "If a client is selected, this plan can only be assigned to a client-scoped price schedule.",
		},
		Actions: ActionLabels{
			View:       "View Plan",
			Edit:       "Edit Plan",
			Delete:     "Delete Plan",
			Activate:   "Activate Plan",
			Deactivate: "Deactivate Plan",
		},
		Bulk: BulkLabels{
			Delete: "Delete Selected",
		},
		Status: StatusLabels{
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Detail: DetailLabels{
			PageTitle:             "Plan Details",
			Price:                 "Price",
			Currency:              "Currency",
			Status:                "Status",
			Description:           "Description",
			FulfillmentType:       "Fulfillment Type",
			CreatedDate:           "Created",
			ModifiedDate:          "Last Modified",
			NoProductsAssigned:    "No products assigned",
			NoProductsAssignedMsg: "No products have been linked to this plan yet.",
			NoProductsDesc:        "No products have been linked to this plan yet.",
			NoPricePlans:          "No price plans",
			NoPricePlansMsg:       "No price plans have been configured for this plan yet.",
			NoPricePlansDesc:      "No price plans have been configured for this plan yet.",
			AuditTrailComingSoon:  "Audit trail coming soon.",
		},
		Tabs: TabLabels{
			Info:          "Information",
			Products:      "Products",
			PricePlan:     "Rate Cards",
			PricePlanSlug: "",
			Attachments:   "Attachments",
			AuditTrail:    "Audit Trail",
			AuditHistory:  "History",
		},
		Confirm: ConfirmLabels{
			Delete:                "Delete Plan",
			DeleteMessage:         "Are you sure you want to delete \"%s\"? This action cannot be undone.",
			Activate:              "Activate Plan",
			ActivateMessage:       "Are you sure you want to activate \"%s\"?",
			Deactivate:            "Deactivate Plan",
			DeactivateMessage:     "Are you sure you want to deactivate \"%s\"?",
			BulkActivate:          "Activate Selected",
			BulkActivateMessage:   "Are you sure you want to activate the selected plans?",
			BulkDeactivate:        "Deactivate Selected",
			BulkDeactivateMessage: "Are you sure you want to deactivate the selected plans?",
			BulkDelete:            "Delete Selected",
			BulkDeleteMessage:     "Are you sure you want to delete the selected plans? This action cannot be undone.",
		},
		Errors: ErrorLabels{
			PermissionDenied:  "You do not have permission to perform this action",
			InvalidFormData:   "Invalid form data. Please check your inputs and try again.",
			NotFound:          "Plan not found",
			IDRequired:        "Plan ID is required",
			NoIDsProvided:     "No plan IDs provided",
			InvalidStatus:     "Invalid status",
			NoPermission:      "No permission",
			CannotDelete:      "This plan cannot be deleted because it has products or rate cards assigned",
			ClientScopeLocked: "Cannot change this plan's client while it has active subscriptions.",
		},
		ProductPlanForm: ProductPlanFormLabels{
			Name:               "Name",
			NamePlaceholder:    "Defaults to the selected product's name",
			Product:            "Product",
			ProductPlaceholder: "Select an item...",
			SelectProduct:      "— Select a product —",
			Active:             "Active",
			ProductKindLabel:   "Item Type",
			ProductKind: ProductKindOptionLabels{
				Service:        "Service",
				StockedGood:    "Stocked Good",
				NonStockedGood: "Non-Stocked Good",
				Consumable:     "Consumable",
			},
			// Model D — variant picker defaults
			VariantSelectLabel:       "Variant",
			VariantSelectPlaceholder: "Select a variant",
			VariantSelectInfo:        "Required when the parent product has variants enabled.",
		},
		Filters: FilterLabels{
			ScopeChipLabel: "Show:",
			ScopeMaster:    "Master",
			ScopeClient:    "Client-specific",
			ScopeAll:       "All",
		},
	}
}

// ---------------------------------------------------------------------------
// Plan list/table sub-labels (relocated from product_price_plan_labels.go god-file)
// ---------------------------------------------------------------------------

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"heading_active"`
	HeadingInactive string `json:"heading_inactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"caption_active"`
	CaptionInactive string `json:"caption_inactive"`
}

type ButtonLabels struct {
	AddPlan       string `json:"add_plan"`
	AddPricePlan  string `json:"add_price_plan"`
	EditPricePlan string `json:"edit_price_plan"`
	AddProduct    string `json:"add_product"`
}

type ColumnLabels struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Interval      string `json:"interval"`
	Price         string `json:"price"`
	Status        string `json:"status"`
	Product       string `json:"product"`
	PricePlan     string `json:"price_plan"`
	PriceSchedule string `json:"price_schedule"`
	Duration      string `json:"duration"`
	Location      string `json:"location"`
	ItemType      string `json:"item_type"`
}

type EmptyLabels struct {
	Title           string `json:"title"`
	Message         string `json:"message"`
	ActiveTitle     string `json:"active_title"`
	ActiveMessage   string `json:"active_message"`
	InactiveTitle   string `json:"inactive_title"`
	InactiveMessage string `json:"inactive_message"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

type StatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	NoIDsProvided    string `json:"no_ids_provided"`
	InvalidStatus    string `json:"invalid_status"`
	NoPermission     string `json:"no_permission"`
	CannotDelete     string `json:"cannot_delete"`

	// 2026-04-27 plan-client-scope plan §7 — surfaced when an operator tries
	// to change a Plan's client_id while one of its PricePlans is attached
	// to an active subscription. Hard block; no force-override.
	ClientScopeLocked string `json:"client_scope_locked"`
}
