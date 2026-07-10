package price_schedule

import (
	"strings"
)

// ---------------------------------------------------------------------------
// Price Schedule labels
// ---------------------------------------------------------------------------

// FilterLabels holds translatable labels for the scope filter chip
// on the price schedule list page (§6.1 of the 2026-04-27 plan-client-scope plan).
type FilterLabels struct {
	ScopeChipLabel string `json:"scope_chip_label"`
	ScopeMaster    string `json:"scope_master"`
	ScopeClient    string `json:"scope_client"`
	ScopeAll       string `json:"scope_all"`
}

// Labels holds all labels for the price schedule module.
type Labels struct {
	Page     PageLabels     `json:"page"`
	Buttons  ButtonLabels   `json:"buttons"`
	Columns  ColumnLabels   `json:"columns"`
	Empty    EmptyLabels    `json:"empty"`
	Form     FormLabels     `json:"form"`
	PlanForm PlanFormLabels `json:"plan_form"`
	Bulk     BulkLabels     `json:"bulk"`
	Confirm  ConfirmLabels  `json:"confirm"`
	Tabs     TabLabels      `json:"tabs"`
	Detail   DetailLabels   `json:"detail"`
	Errors   ErrorLabels    `json:"errors"`
	Filters  FilterLabels   `json:"filters"`
}

// PlanFormLabels holds labels for the "Add Plan" (price_plan) drawer form
// within a price schedule. Professional tier overrides field names (e.g., "Package").
type PlanFormLabels struct {
	SectionSchedule        string `json:"section_schedule"`
	SectionPackage         string `json:"section_package"`
	SectionPricing         string `json:"section_pricing"`
	PriceScheduleField     string `json:"price_schedule_field"`
	PackageLabel           string `json:"package_label"`
	PackagePlaceholder     string `json:"package_placeholder"`
	PackageSearch          string `json:"package_search"`
	NameLabel              string `json:"name_label"`
	NamePlaceholder        string `json:"name_placeholder"`
	DescriptionLabel       string `json:"description_label"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	AmountLabel            string `json:"amount_label"`
	AmountPlaceholder      string `json:"amount_placeholder"`
	CurrencyLabel          string `json:"currency_label"`
	CurrencyPlaceholder    string `json:"currency_placeholder"`
	DurationLabel          string `json:"duration_label"`
	UnitLabel              string `json:"unit_label"`
	ActiveLabel            string `json:"active_label"`
	SchedulePlaceholder    string `json:"schedule_placeholder"`
	ScheduleSearch         string `json:"schedule_search"`
	LocationHintPrefix     string `json:"location_hint_prefix"`
}

type PageLabels struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	ActiveTitle   string `json:"active_title"`
	InactiveTitle string `json:"inactive_title"`
}

type ButtonLabels struct {
	View       string `json:"view"`
	Add        string `json:"add"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	BulkDelete string `json:"bulk_delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type ColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DateStart   string `json:"date_start"`
	DateEnd     string `json:"date_end"`
	Location    string `json:"location"`
	Status      string `json:"status"`
	DateCreated string `json:"date_created"`
	Actions     string `json:"actions"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type FormLabels struct {
	Name            string `json:"name"`
	NamePlaceholder string `json:"name_placeholder"`
	Description     string `json:"description"`
	DescPlaceholder string `json:"desc_placeholder"`
	DateStart       string `json:"date_start"`
	DateEnd         string `json:"date_end"`
	// Optional time inputs paired with DateStart/DateEnd (2026-04-28 date+time
	// field plan). TimePlaceholder is shared by both inputs.
	TimeStart           string `json:"time_start"`
	TimeEnd             string `json:"time_end"`
	TimePlaceholder     string `json:"time_placeholder"`
	Location            string `json:"location"`
	LocationPlaceholder string `json:"location_placeholder"`
	SelectLocation      string `json:"select_location"`
	Active              string `json:"active"`

	// Wave 2 — section headers (from lyngua price_schedule.json → priceSchedule.form)
	SectionScheduleDetails string `json:"section_schedule_details"`
	SectionDateRange       string `json:"section_date_range"`
	SectionLocation        string `json:"section_location"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"name_info"`
	DescriptionInfo string `json:"description_info"`
	DateStartInfo   string `json:"date_start_info"`
	DateEndInfo     string `json:"date_end_info"`
	TimeStartInfo   string `json:"time_start_info"`
	TimeEndInfo     string `json:"time_end_info"`
	LocationInfo    string `json:"location_info"`
	ActiveInfo      string `json:"active_info"`

	// Client-scope fields (2026-04-27 plan-client-scope plan §7).
	// Set on the schedule add/edit drawer Client picker. The suffix is
	// appended to the client's name to produce the default schedule name
	// (e.g. "Cruz Engineering - Rate Cards" on professional tier, or
	// "Cruz Engineering - Price Schedule" on general). See plan §4.4.1.
	ClientLabel                          string `json:"client_label"`
	ClientHelp                           string `json:"client_help"`
	ClientPlaceholder                    string `json:"client_placeholder"`
	ClientSearchPlaceholder              string `json:"client_search_placeholder"`
	ClientNoResults                      string `json:"client_no_results"`
	ClientInfo                           string `json:"client_info"`
	CustomClientPriceScheduleLabelSuffix string `json:"custom_client_price_schedule_label_suffix"`
	LocationSearchPlaceholder            string `json:"location_search_placeholder"`

	// Scope radio (2026-04-28) — mutually exclusive Location / Client picker.
	ScopeLabel              string `json:"scope_label"`
	ScopeInfo               string `json:"scope_info"`
	ScopeOptionLocation     string `json:"scope_option_location"`
	ScopeOptionClient       string `json:"scope_option_client"`
	ScopeOptionLocationHelp string `json:"scope_option_location_help"`
	ScopeOptionClientHelp   string `json:"scope_option_client_help"`
}

type BulkLabels struct {
	DeleteTitle       string `json:"delete_title"`
	DeleteMessage     string `json:"delete_message"`
	ActivateTitle     string `json:"activate_title"`
	ActivateMessage   string `json:"activate_message"`
	DeactivateTitle   string `json:"deactivate_title"`
	DeactivateMessage string `json:"deactivate_message"`
}

type ConfirmLabels struct {
	DeleteTitle       string `json:"delete_title"`
	DeleteMessage     string `json:"delete_message"`
	ActivateTitle     string `json:"activate_title"`
	ActivateMessage   string `json:"activate_message"`
	DeactivateTitle   string `json:"deactivate_title"`
	DeactivateMessage string `json:"deactivate_message"`
}

type TabLabels struct {
	Info              string `json:"info"`
	PricePlan         string `json:"price_plan"`
	PricePlanSlug     string `json:"price_plan_slug"`
	ProductPrices     string `json:"product_prices"`
	ProductPricesSlug string `json:"product_prices_slug"`
	// 2026-05-04 — Subscriptions/Engagements tab on the schedule-scoped
	// price_plan detail. Professional tier overrides the label to
	// "Engagements"; URL slug stays "subscriptions" across tiers.
	Subscriptions     string `json:"subscriptions"`
	SubscriptionsSlug string `json:"subscriptions_slug"`
}

// ResolveTabSlug returns the URL slug for a canonical tab key. Today only the
// "pricePlan" tab on the parent detail and "product-prices" on the nested plan
// detail are re-slugged (e.g., professional tier ships "package-prices" /
// "package-item-prices"); other tabs round-trip through as-is.
func (t TabLabels) ResolveTabSlug(canonical string) string {
	switch canonical {
	case "pricePlan":
		if s := strings.TrimSpace(t.PricePlanSlug); s != "" {
			return s
		}
	case "product-prices":
		if s := strings.TrimSpace(t.ProductPricesSlug); s != "" {
			return s
		}
	case "subscriptions":
		if s := strings.TrimSpace(t.SubscriptionsSlug); s != "" {
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
	if s := strings.TrimSpace(t.PricePlanSlug); s != "" && slug == s {
		return "pricePlan"
	}
	if s := strings.TrimSpace(t.ProductPricesSlug); s != "" && slug == s {
		return "product-prices"
	}
	if s := strings.TrimSpace(t.SubscriptionsSlug); s != "" && slug == s {
		return "subscriptions"
	}
	return slug
}

type DetailLabels struct {
	Title                 string `json:"title"`
	DateCreated           string `json:"date_created"`
	DateModified          string `json:"date_modified"`
	NoLocation            string `json:"no_location"`
	NoDateEnd             string `json:"no_date_end"`
	NoDescription         string `json:"no_description"`
	PlansEmptyTitle       string `json:"plans_empty_title"`
	PlansEmptyMsg         string `json:"plans_empty_msg"`
	NoDescriptionSubtitle string `json:"no_description_subtitle"`

	// Product price (per-product breakdown, shown on the schedule-scoped plan detail).
	// Professional tier renames these to "Service Price" via lyngua.
	ProductPriceAdd           string `json:"product_price_add"`
	ProductPriceEdit          string `json:"product_price_edit"`
	ProductPriceDelete        string `json:"product_price_delete"`
	ProductPriceDeleteConfirm string `json:"product_price_delete_confirm"`
	ProductPriceEmptyTitle    string `json:"product_price_empty_title"`
	ProductPriceEmptyMsg      string `json:"product_price_empty_msg"`
	ProductPriceSection       string `json:"product_price_section"` // drawer section title ("Product Price" / "Service Price")
	ProductField              string `json:"product_field"`         // drawer product select label ("Product" / "Service")

	// Plans table columns (price-schedule-detail plans tab).
	PlanColumnPlan        string `json:"plan_column_plan"`
	PlanColumnAmount      string `json:"plan_column_amount"`
	PlanColumnBillingKind string `json:"plan_column_billing_kind"`
	PlanColumnAmountBasis string `json:"plan_column_amount_basis"`
	PlanColumnCadence     string `json:"plan_column_cadence"`
	PlanColumnDuration    string `json:"plan_column_duration"` // deprecated; replaced by PlanColumnCadence
	PlanColumnStatus      string `json:"plan_column_status"`

	// Cadence cell prefixes per BillingKind (rendered as "{prefix} {cycle}" or
	// just "{prefix}" when no cycle applies).
	CadenceOneTime     string `json:"cadence_one_time"`  // e.g. "One-time payment"
	CadenceRecurring   string `json:"cadence_recurring"` // e.g. "Every {cycle}"
	CadenceContract    string `json:"cadence_contract"`  // e.g. "Contract — billed every {cycle}"
	CadenceMilestone   string `json:"cadence_milestone"` // e.g. "Per milestone"
	CadenceAdHoc       string `json:"cadence_ad_hoc"`    // e.g. "Per occurrence"
	CadenceUnspecified string `json:"cadence_unspecified"`

	// Compact labels for the BillingKind / AmountBasis cells in the plans table.
	BillingKindOneTime          string `json:"billing_kind_one_time"`
	BillingKindRecurring        string `json:"billing_kind_recurring"`
	BillingKindContract         string `json:"billing_kind_contract"`
	BillingKindMilestone        string `json:"billing_kind_milestone"`
	BillingKindAdHoc            string `json:"billing_kind_ad_hoc"`
	AmountBasisPerCycle         string `json:"amount_basis_per_cycle"`
	AmountBasisTotalPackage     string `json:"amount_basis_total_package"`
	AmountBasisDerivedFromLines string `json:"amount_basis_derived_from_lines"`
	AmountBasisPerOccurrence    string `json:"amount_basis_per_occurrence"`

	// Plans table row actions + confirms.
	PlanView            string `json:"plan_view"`
	PlanEdit            string `json:"plan_edit"`
	PlanEditDrawerTitle string `json:"plan_edit_drawer_title"`
	PlanDelete          string `json:"plan_delete"`
	PlanDeleteTitle     string `json:"plan_delete_title"`
	PlanDeleteMsg       string `json:"plan_delete_msg"`
	PlanInUseTooltip    string `json:"plan_in_use_tooltip"`

	// Plans table primary action + inline error messages.
	PlanAdd      string `json:"plan_add"`
	PlanRequired string `json:"plan_required"`

	// Product prices table columns.
	ProductPriceColumnProduct   string `json:"product_price_column_product"`
	ProductPriceColumnPrice     string `json:"product_price_column_price"`
	ProductPriceColumnCurrency  string `json:"product_price_column_currency"`
	ProductPriceColumnTreatment string `json:"product_price_column_treatment"`
	ProductPriceColumnEffective string `json:"product_price_column_effective"`

	// Drawer banners explaining how the per-line price relates to the parent
	// PricePlan.amount_basis. Surfaced above the Price input.
	BasisBannerPerCycle     string `json:"basis_banner_per_cycle"`
	BasisBannerTotalPackage string `json:"basis_banner_total_package"`
	BasisBannerDerived      string `json:"basis_banner_derived"`

	// Drawer section labels used by the schedule-scoped PPP drawer.
	ProductPriceCatalogSection   string `json:"product_price_catalog_section"`
	ProductPricePricingSection   string `json:"product_price_pricing_section"`
	ProductPriceEffectiveSection string `json:"product_price_effective_section"`

	// Attachments tab label for the price_schedule detail page and the nested
	// price_plan (plan) detail page.
	TabAttachments string `json:"tab_attachments"`
}

type ErrorLabels struct {
	NotFound                   string `json:"not_found"`
	LoadFailed                 string `json:"load_failed"`
	Unauthorized               string `json:"unauthorized"`
	CreateFailed               string `json:"create_failed"`
	UpdateFailed               string `json:"update_failed"`
	DeleteFailed               string `json:"delete_failed"`
	InUse                      string `json:"in_use"`
	PricePlanCreateUnavailable string `json:"price_plan_create_unavailable"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:         "Price Schedules",
			Subtitle:      "Manage your price schedules",
			ActiveTitle:   "Active Price Schedules",
			InactiveTitle: "Inactive Price Schedules",
		},
		Buttons: ButtonLabels{
			View:       "View",
			Add:        "Add Price Schedule",
			Edit:       "Edit Price Schedule",
			Delete:     "Delete Price Schedule",
			BulkDelete: "Delete Price Schedules",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: ColumnLabels{
			Name:        "Name",
			Description: "Description",
			DateStart:   "Start Date",
			DateEnd:     "End Date",
			Location:    "Location",
			Status:      "Status",
			DateCreated: "Date Created",
			Actions:     "Actions",
		},
		Empty: EmptyLabels{
			Title:   "No Price Schedules",
			Message: "No price schedules to display.",
		},
		Form: FormLabels{
			Name:                "Name",
			NamePlaceholder:     "Enter price schedule name",
			Description:         "Description",
			DescPlaceholder:     "Enter description...",
			DateStart:           "Start Date",
			DateEnd:             "End Date",
			TimeStart:           "Start Time (optional)",
			TimeEnd:             "End Time (optional)",
			TimePlaceholder:     "HH:MM",
			Location:            "Location",
			LocationPlaceholder: "Select a location...",
			SelectLocation:      "— No location (all locations) —",
			Active:              "Active",
			// Wave 2 new section headers
			SectionScheduleDetails: "Schedule details",
			SectionDateRange:       "Date range",
			SectionLocation:        "Location",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			NameInfo:        "A short display name for this price schedule.",
			DescriptionInfo: "Optional notes or context for this price schedule.",
			DateStartInfo:   "First date this price schedule becomes effective.",
			DateEndInfo:     "Last date this price schedule is effective. Leave empty for no end date.",
			TimeStartInfo:   "Optional time of day in the operator's display timezone. Leave blank for start of day (00:00).",
			TimeEndInfo:     "Optional time of day in the operator's display timezone. Leave blank for end of day (23:59).",
			LocationInfo:    "Restrict this price schedule to a specific location, or leave empty to apply to all locations.",
			ActiveInfo:      "Inactive price schedules are hidden from new subscriptions.",
			// Client-scope fields (2026-04-27 plan-client-scope plan §7).
			ClientLabel:                          "Client",
			ClientHelp:                           "Leave blank for a general schedule. Set a client to create a bespoke schedule reused across that client's price plans.",
			ClientPlaceholder:                    "Leave blank for a general schedule",
			ClientSearchPlaceholder:              "Search clients...",
			ClientNoResults:                      "No clients found",
			ClientInfo:                           "Optional. When set, this schedule is reserved for that client's bespoke price plans.",
			CustomClientPriceScheduleLabelSuffix: "Price Schedule",
			LocationSearchPlaceholder:            "Filter...",
			// Scope radio (2026-04-28).
			ScopeLabel:              "Scope",
			ScopeInfo:               "Choose whether this schedule is shared across every client at a location, or reserved for one client's bespoke pricing. Switching scope clears the inactive picker on save.",
			ScopeOptionLocation:     "Location-scoped",
			ScopeOptionClient:       "Client-scoped",
			ScopeOptionLocationHelp: "Reusable across all clients at this location.",
			ScopeOptionClientHelp:   "Reserved for one client's bespoke pricing.",
		},
		Bulk: BulkLabels{
			DeleteTitle:       "Delete Price Schedules",
			DeleteMessage:     "Permanently delete the selected price schedules? This cannot be undone.",
			ActivateTitle:     "Activate Price Schedules",
			ActivateMessage:   "Activate the selected price schedules?",
			DeactivateTitle:   "Deactivate Price Schedules",
			DeactivateMessage: "Deactivate the selected price schedules?",
		},
		Confirm: ConfirmLabels{
			DeleteTitle:       "Delete Price Schedule",
			DeleteMessage:     "Permanently delete this price schedule? This cannot be undone.",
			ActivateTitle:     "Activate Price Schedule",
			ActivateMessage:   "Activate {{name}}?",
			DeactivateTitle:   "Deactivate Price Schedule",
			DeactivateMessage: "Deactivate {{name}}?",
		},
		Tabs: TabLabels{
			Info:              "Info",
			PricePlan:         "Plans",
			PricePlanSlug:     "",
			ProductPrices:     "Product Prices",
			Subscriptions:     "Subscriptions",
			SubscriptionsSlug: "",
		},
		Detail: DetailLabels{
			Title:                     "Price Schedule",
			DateCreated:               "Date Created",
			DateModified:              "Date Modified",
			NoLocation:                "All locations",
			NoDateEnd:                 "No end date",
			NoDescription:             "—",
			PlansEmptyTitle:           "No Plans",
			PlansEmptyMsg:             "No price plans are linked to this schedule yet.",
			NoDescriptionSubtitle:     "No description provided",
			ProductPriceAdd:           "Add Product Price",
			ProductPriceEdit:          "Edit Product Price",
			ProductPriceDelete:        "Delete Product Price",
			ProductPriceDeleteConfirm: "Remove %s from this plan?",
			ProductPriceEmptyTitle:    "No Product Prices",
			ProductPriceEmptyMsg:      "No product prices have been configured for this plan yet.",
			ProductPriceSection:       "Product Price",
			ProductField:              "Product",

			PlanColumnPlan:        "Plan",
			PlanColumnAmount:      "Amount",
			PlanColumnBillingKind: "Billing model",
			PlanColumnAmountBasis: "Amount basis",
			PlanColumnCadence:     "Cadence",
			PlanColumnDuration:    "Duration",
			PlanColumnStatus:      "Status",

			CadenceOneTime:     "One-time payment",
			CadenceRecurring:   "Every %s",
			CadenceContract:    "Contract — billed every %s",
			CadenceMilestone:   "Per milestone",
			CadenceAdHoc:       "Per occurrence",
			CadenceUnspecified: "—",

			BillingKindOneTime:          "One-time",
			BillingKindRecurring:        "Recurring",
			BillingKindContract:         "Contract",
			BillingKindMilestone:        "Milestone",
			BillingKindAdHoc:            "Ad hoc",
			AmountBasisPerCycle:         "Per cycle",
			AmountBasisTotalPackage:     "Total package",
			AmountBasisDerivedFromLines: "Derived from lines",
			AmountBasisPerOccurrence:    "Per occurrence",

			PlanView:            "View",
			PlanEdit:            "Edit",
			PlanEditDrawerTitle: "Edit Plan",
			PlanDelete:          "Delete",
			PlanDeleteTitle:     "Delete Plan",
			PlanDeleteMsg:       "Permanently delete %s? This cannot be undone.",
			PlanInUseTooltip:    "In use by active subscriptions",

			PlanAdd:      "Add Plan",
			PlanRequired: "Plan is required",

			ProductPriceColumnProduct:    "Product",
			ProductPriceColumnPrice:      "Price",
			ProductPriceColumnCurrency:   "Currency",
			ProductPriceColumnTreatment:  "Billing",
			ProductPriceColumnEffective:  "Effective",
			BasisBannerPerCycle:          "Each line below is charged every billing cycle.",
			BasisBannerTotalPackage:      "These per-line prices are informational. The package is sold at a flat rate; the total here does not have to match.",
			BasisBannerDerived:           "The package price is the sum of these line prices. Editing a line changes the package total.",
			ProductPriceCatalogSection:   "Catalog line",
			ProductPricePricingSection:   "Pricing",
			ProductPriceEffectiveSection: "Effective dates",
			TabAttachments:               "Attachments",
		},
		PlanForm: PlanFormLabels{
			SectionSchedule:        "Schedule",
			SectionPackage:         "Plan",
			SectionPricing:         "Pricing",
			PriceScheduleField:     "Price Schedule",
			PackageLabel:           "Plan",
			PackagePlaceholder:     "Select a plan...",
			PackageSearch:          "Filter...",
			NameLabel:              "Plan Name",
			NamePlaceholder:        "Enter plan name",
			DescriptionLabel:       "Description",
			DescriptionPlaceholder: "Optional notes for this package",
			AmountLabel:            "Amount",
			AmountPlaceholder:      "0.00",
			CurrencyLabel:          "Currency",
			CurrencyPlaceholder:    "e.g. PHP",
			DurationLabel:          "Duration",
			UnitLabel:              "Unit",
			ActiveLabel:            "Active",
			SchedulePlaceholder:    "Select a rate card...",
			ScheduleSearch:         "Filter...",
			LocationHintPrefix:     "Location: ",
		},
		Errors: ErrorLabels{
			NotFound:                   "Price schedule not found",
			LoadFailed:                 "Failed to load price schedule",
			Unauthorized:               "You are not authorized to perform this action",
			CreateFailed:               "Failed to create price schedule",
			UpdateFailed:               "Failed to update price schedule",
			DeleteFailed:               "Failed to delete price schedule",
			InUse:                      "This price schedule is in use by active subscriptions and cannot be deleted.",
			PricePlanCreateUnavailable: "Adding a price plan is not available. Please contact support.",
		},
		Filters: FilterLabels{
			ScopeChipLabel: "Show:",
			ScopeMaster:    "Master",
			ScopeClient:    "Client-specific",
			ScopeAll:       "All",
		},
	}
}

// ClientPackagesLabels holds labels for the client detail "Packages" tab —
// the list of client-scoped Plans for a given client, with the
// "Add custom package" CTA. Mounted from entydad's client detail page via
// a centymo helper view (plan §6.6 option 1).
//
// 2026-04-27 plan-client-scope plan §6.3 / §7.
