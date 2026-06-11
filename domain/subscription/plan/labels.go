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
	ScopeChipLabel string `json:"scopeChipLabel"`
	ScopeMaster    string `json:"scopeMaster"`
	ScopeClient    string `json:"scopeClient"`
	ScopeAll       string `json:"scopeAll"`
}

// Labels holds all translatable strings for the plan module.
type Labels struct {
	Page            PageLabels            `json:"page"`
	Buttons         ButtonLabels          `json:"buttons"`
	Columns         ColumnLabels          `json:"columns"`
	Empty           EmptyLabels           `json:"empty"`
	Form            FormLabels            `json:"form"`
	Actions         ActionLabels          `json:"actions"`
	Bulk            BulkLabels            `json:"bulkActions"`
	Status          StatusLabels          `json:"status"`
	Detail          DetailLabels          `json:"detail"`
	Tabs            TabLabels             `json:"tabs"`
	Confirm         ConfirmLabels         `json:"confirm"`
	Errors          ErrorLabels           `json:"errors"`
	ProductPlanForm ProductPlanFormLabels `json:"productPlanForm"`
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
	NamePlaceholder     string            `json:"namePlaceholder"`
	Description         string            `json:"description"`
	DescPlaceholder     string            `json:"descriptionPlaceholder"`
	FulfillmentType     string            `json:"fulfillmentType"`
	Active              string            `json:"active"`
	Products            string            `json:"products"`
	ProductsPlaceholder string            `json:"productsPlaceholder"`
	ProductsSearch      string            `json:"productsSearch"`
	Sections            FormSectionLabels `json:"sections"`

	// Fulfillment type option labels
	TypeSchedule string `json:"typeSchedule"`
	TypeLicense  string `json:"typeLicense"`
	TypeContent  string `json:"typeContent"`
	TypePhysical string `json:"typePhysical"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	ActiveInfo      string `json:"activeInfo"`

	// Client-scope fields (2026-04-27 plan-client-scope plan §7).
	// Set on the Plan add/edit drawer Client picker.
	ClientLabel             string `json:"clientLabel"`
	ClientHelp              string `json:"clientHelp"`
	ClientPlaceholder       string `json:"clientPlaceholder"`
	ClientSearchPlaceholder string `json:"clientSearchPlaceholder"`
	ClientNoResults         string `json:"clientNoResults"`
	ClientLockedTooltip     string `json:"clientLockedTooltip"`
	ClientForLabel          string `json:"clientForLabel"` // "For {{.ClientName}}" — read-only badge in client-context entry-point
	ClientInfo              string `json:"clientInfo"`

	// JobTemplate select (2026-04-29 auto-spawn-jobs-from-subscription plan §5
	// — Plan.job_template_id assignment from the drawer). Empty value =
	// advisory-only plan; spawn use case skips silently.
	JobTemplate     string `json:"jobTemplate"`
	JobTemplateNone string `json:"jobTemplateNone"`
	JobTemplateHint string `json:"jobTemplateHint"`

	// 2026-04-30 cyclic-subscription-jobs plan §9.3 — visits_per_cycle field.
	// Number of cycle Job instances spawned per billing cycle (default 1).
	VisitsPerCycleLabel       string `json:"visitsPerCycleLabel"`
	VisitsPerCyclePlaceholder string `json:"visitsPerCyclePlaceholder"`
	VisitsPerCycleHint        string `json:"visitsPerCycleHint"`

	// Client-scope cascade notice — shown unconditionally below the client picker
	// so operators see the schedule restriction before filling other fields.
	// Tier-specific wording lives in lyngua; default uses proto-generic vocabulary.
	ClientScopeCascadeNotice string `json:"clientScopeCascadeNotice"`
}

type DetailLabels struct {
	PageTitle             string `json:"pageTitle"`
	Price                 string `json:"price"`
	Currency              string `json:"currency"`
	Status                string `json:"status"`
	Description           string `json:"description"`
	FulfillmentType       string `json:"fulfillmentType"`
	CreatedDate           string `json:"createdDate"`
	ModifiedDate          string `json:"modifiedDate"`
	NoProductsAssigned    string `json:"noProductsAssigned"`
	NoProductsAssignedMsg string `json:"noProductsAssignedMsg"`
	NoProductsDesc        string `json:"noProductsDesc"`
	NoPricePlans          string `json:"noPricePlans"`
	NoPricePlansMsg       string `json:"noPricePlansMsg"`
	NoPricePlansDesc      string `json:"noPricePlansDesc"`
	AuditTrailComingSoon  string `json:"auditTrailComingSoon"`
}

type TabLabels struct {
	Info          string `json:"info"`
	Products      string `json:"products"`
	ProductsSlug  string `json:"productsSlug"`
	PricePlan     string `json:"pricePlan"`
	PricePlanSlug string `json:"pricePlanSlug"`
	Attachments   string `json:"attachments"`
	AuditTrail    string `json:"auditTrail"`
	AuditHistory  string `json:"auditHistory"`
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
	DeleteMessage         string `json:"deleteMessage"`
	Activate              string `json:"activate"`
	ActivateMessage       string `json:"activateMessage"`
	Deactivate            string `json:"deactivate"`
	DeactivateMessage     string `json:"deactivateMessage"`
	BulkActivate          string `json:"bulkActivate"`
	BulkActivateMessage   string `json:"bulkActivateMessage"`
	BulkDeactivate        string `json:"bulkDeactivate"`
	BulkDeactivateMessage string `json:"bulkDeactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
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
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type ButtonLabels struct {
	AddPlan       string `json:"addPlan"`
	AddPricePlan  string `json:"addPricePlan"`
	EditPricePlan string `json:"editPricePlan"`
	AddProduct    string `json:"addProduct"`
}

type ColumnLabels struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Interval      string `json:"interval"`
	Price         string `json:"price"`
	Status        string `json:"status"`
	Product       string `json:"product"`
	PricePlan     string `json:"pricePlan"`
	PriceSchedule string `json:"priceSchedule"`
	Duration      string `json:"duration"`
	Location      string `json:"location"`
	ItemType      string `json:"itemType"`
}

type EmptyLabels struct {
	Title           string `json:"title"`
	Message         string `json:"message"`
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
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
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	InvalidStatus    string `json:"invalidStatus"`
	NoPermission     string `json:"noPermission"`
	CannotDelete     string `json:"cannotDelete"`

	// 2026-04-27 plan-client-scope plan §7 — surfaced when an operator tries
	// to change a Plan's client_id while one of its PricePlans is attached
	// to an active subscription. Hard block; no force-override.
	ClientScopeLocked string `json:"clientScopeLocked"`
}
