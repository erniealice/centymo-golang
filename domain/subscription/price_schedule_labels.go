package subscription

import (
	"strings"
)

// ---------------------------------------------------------------------------
// Price Schedule labels
// ---------------------------------------------------------------------------

// PriceScheduleFilterLabels holds translatable labels for the scope filter chip
// on the price schedule list page (§6.1 of the 2026-04-27 plan-client-scope plan).
type PriceScheduleFilterLabels struct {
	ScopeChipLabel string `json:"scopeChipLabel"`
	ScopeMaster    string `json:"scopeMaster"`
	ScopeClient    string `json:"scopeClient"`
	ScopeAll       string `json:"scopeAll"`
}

// PriceScheduleLabels holds all labels for the price schedule module.
type PriceScheduleLabels struct {
	Page     PriceSchedulePageLabels     `json:"page"`
	Buttons  PriceScheduleButtonLabels   `json:"buttons"`
	Columns  PriceScheduleColumnLabels   `json:"columns"`
	Empty    PriceScheduleEmptyLabels    `json:"empty"`
	Form     PriceScheduleFormLabels     `json:"form"`
	PlanForm PriceSchedulePlanFormLabels `json:"planForm"`
	Bulk     PriceScheduleBulkLabels     `json:"bulk"`
	Confirm  PriceScheduleConfirmLabels  `json:"confirm"`
	Tabs     PriceScheduleTabLabels      `json:"tabs"`
	Detail   PriceScheduleDetailLabels   `json:"detail"`
	Errors   PriceScheduleErrorLabels    `json:"errors"`
	Filters  PriceScheduleFilterLabels   `json:"filters"`
}

// PriceSchedulePlanFormLabels holds labels for the "Add Plan" (price_plan) drawer form
// within a price schedule. Professional tier overrides field names (e.g., "Package").
type PriceSchedulePlanFormLabels struct {
	SectionSchedule        string `json:"sectionSchedule"`
	SectionPackage         string `json:"sectionPackage"`
	SectionPricing         string `json:"sectionPricing"`
	PriceScheduleField     string `json:"priceScheduleField"`
	PackageLabel           string `json:"packageLabel"`
	PackagePlaceholder     string `json:"packagePlaceholder"`
	PackageSearch          string `json:"packageSearch"`
	NameLabel              string `json:"nameLabel"`
	NamePlaceholder        string `json:"namePlaceholder"`
	DescriptionLabel       string `json:"descriptionLabel"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	AmountLabel            string `json:"amountLabel"`
	AmountPlaceholder      string `json:"amountPlaceholder"`
	CurrencyLabel          string `json:"currencyLabel"`
	CurrencyPlaceholder    string `json:"currencyPlaceholder"`
	DurationLabel          string `json:"durationLabel"`
	UnitLabel              string `json:"unitLabel"`
	ActiveLabel            string `json:"activeLabel"`
	SchedulePlaceholder    string `json:"schedulePlaceholder"`
	ScheduleSearch         string `json:"scheduleSearch"`
	LocationHintPrefix     string `json:"locationHintPrefix"`
}

type PriceSchedulePageLabels struct {
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	ActiveTitle   string `json:"activeTitle"`
	InactiveTitle string `json:"inactiveTitle"`
}

type PriceScheduleButtonLabels struct {
	View       string `json:"view"`
	Add        string `json:"add"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	BulkDelete string `json:"bulkDelete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type PriceScheduleColumnLabels struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DateStart   string `json:"dateStart"`
	DateEnd     string `json:"dateEnd"`
	Location    string `json:"location"`
	Status      string `json:"status"`
	DateCreated string `json:"dateCreated"`
	Actions     string `json:"actions"`
}

type PriceScheduleEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type PriceScheduleFormLabels struct {
	Name            string `json:"name"`
	NamePlaceholder string `json:"namePlaceholder"`
	Description     string `json:"description"`
	DescPlaceholder string `json:"descPlaceholder"`
	DateStart       string `json:"dateStart"`
	DateEnd         string `json:"dateEnd"`
	// Optional time inputs paired with DateStart/DateEnd (2026-04-28 date+time
	// field plan). TimePlaceholder is shared by both inputs.
	TimeStart           string `json:"timeStart"`
	TimeEnd             string `json:"timeEnd"`
	TimePlaceholder     string `json:"timePlaceholder"`
	Location            string `json:"location"`
	LocationPlaceholder string `json:"locationPlaceholder"`
	SelectLocation      string `json:"selectLocation"`
	Active              string `json:"active"`

	// Wave 2 — section headers (from lyngua price_schedule.json → priceSchedule.form)
	SectionScheduleDetails string `json:"sectionScheduleDetails"`
	SectionDateRange       string `json:"sectionDateRange"`
	SectionLocation        string `json:"sectionLocation"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	DateStartInfo   string `json:"dateStartInfo"`
	DateEndInfo     string `json:"dateEndInfo"`
	TimeStartInfo   string `json:"timeStartInfo"`
	TimeEndInfo     string `json:"timeEndInfo"`
	LocationInfo    string `json:"locationInfo"`
	ActiveInfo      string `json:"activeInfo"`

	// Client-scope fields (2026-04-27 plan-client-scope plan §7).
	// Set on the schedule add/edit drawer Client picker. The suffix is
	// appended to the client's name to produce the default schedule name
	// (e.g. "Cruz Engineering - Rate Cards" on professional tier, or
	// "Cruz Engineering - Price Schedule" on general). See plan §4.4.1.
	ClientLabel                          string `json:"clientLabel"`
	ClientHelp                           string `json:"clientHelp"`
	ClientPlaceholder                    string `json:"clientPlaceholder"`
	ClientSearchPlaceholder              string `json:"clientSearchPlaceholder"`
	ClientNoResults                      string `json:"clientNoResults"`
	ClientInfo                           string `json:"clientInfo"`
	CustomClientPriceScheduleLabelSuffix string `json:"customClientPriceScheduleLabelSuffix"`
	LocationSearchPlaceholder            string `json:"locationSearchPlaceholder"`

	// Scope radio (2026-04-28) — mutually exclusive Location / Client picker.
	ScopeLabel              string `json:"scopeLabel"`
	ScopeInfo               string `json:"scopeInfo"`
	ScopeOptionLocation     string `json:"scopeOptionLocation"`
	ScopeOptionClient       string `json:"scopeOptionClient"`
	ScopeOptionLocationHelp string `json:"scopeOptionLocationHelp"`
	ScopeOptionClientHelp   string `json:"scopeOptionClientHelp"`
}

type PriceScheduleBulkLabels struct {
	DeleteTitle       string `json:"deleteTitle"`
	DeleteMessage     string `json:"deleteMessage"`
	ActivateTitle     string `json:"activateTitle"`
	ActivateMessage   string `json:"activateMessage"`
	DeactivateTitle   string `json:"deactivateTitle"`
	DeactivateMessage string `json:"deactivateMessage"`
}

type PriceScheduleConfirmLabels struct {
	DeleteTitle       string `json:"deleteTitle"`
	DeleteMessage     string `json:"deleteMessage"`
	ActivateTitle     string `json:"activateTitle"`
	ActivateMessage   string `json:"activateMessage"`
	DeactivateTitle   string `json:"deactivateTitle"`
	DeactivateMessage string `json:"deactivateMessage"`
}

type PriceScheduleTabLabels struct {
	Info              string `json:"info"`
	PricePlan         string `json:"pricePlan"`
	PricePlanSlug     string `json:"pricePlanSlug"`
	ProductPrices     string `json:"productPrices"`
	ProductPricesSlug string `json:"productPricesSlug"`
	// 2026-05-04 — Subscriptions/Engagements tab on the schedule-scoped
	// price_plan detail. Professional tier overrides the label to
	// "Engagements"; URL slug stays "subscriptions" across tiers.
	Subscriptions     string `json:"subscriptions"`
	SubscriptionsSlug string `json:"subscriptionsSlug"`
}

// ResolveTabSlug returns the URL slug for a canonical tab key. Today only the
// "pricePlan" tab on the parent detail and "product-prices" on the nested plan
// detail are re-slugged (e.g., professional tier ships "package-prices" /
// "package-item-prices"); other tabs round-trip through as-is.
func (t PriceScheduleTabLabels) ResolveTabSlug(canonical string) string {
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
func (t PriceScheduleTabLabels) CanonicalizeTab(slug string) string {
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

type PriceScheduleDetailLabels struct {
	Title                 string `json:"title"`
	DateCreated           string `json:"dateCreated"`
	DateModified          string `json:"dateModified"`
	NoLocation            string `json:"noLocation"`
	NoDateEnd             string `json:"noDateEnd"`
	NoDescription         string `json:"noDescription"`
	PlansEmptyTitle       string `json:"plansEmptyTitle"`
	PlansEmptyMsg         string `json:"plansEmptyMsg"`
	NoDescriptionSubtitle string `json:"noDescriptionSubtitle"`

	// Product price (per-product breakdown, shown on the schedule-scoped plan detail).
	// Professional tier renames these to "Service Price" via lyngua.
	ProductPriceAdd           string `json:"productPriceAdd"`
	ProductPriceEdit          string `json:"productPriceEdit"`
	ProductPriceDelete        string `json:"productPriceDelete"`
	ProductPriceDeleteConfirm string `json:"productPriceDeleteConfirm"`
	ProductPriceEmptyTitle    string `json:"productPriceEmptyTitle"`
	ProductPriceEmptyMsg      string `json:"productPriceEmptyMsg"`
	ProductPriceSection       string `json:"productPriceSection"` // drawer section title ("Product Price" / "Service Price")
	ProductField              string `json:"productField"`        // drawer product select label ("Product" / "Service")

	// Plans table columns (price-schedule-detail plans tab).
	PlanColumnPlan        string `json:"planColumnPlan"`
	PlanColumnAmount      string `json:"planColumnAmount"`
	PlanColumnBillingKind string `json:"planColumnBillingKind"`
	PlanColumnAmountBasis string `json:"planColumnAmountBasis"`
	PlanColumnCadence     string `json:"planColumnCadence"`
	PlanColumnDuration    string `json:"planColumnDuration"` // deprecated; replaced by PlanColumnCadence
	PlanColumnStatus      string `json:"planColumnStatus"`

	// Cadence cell prefixes per BillingKind (rendered as "{prefix} {cycle}" or
	// just "{prefix}" when no cycle applies).
	CadenceOneTime     string `json:"cadenceOneTime"`   // e.g. "One-time payment"
	CadenceRecurring   string `json:"cadenceRecurring"` // e.g. "Every {cycle}"
	CadenceContract    string `json:"cadenceContract"`  // e.g. "Contract — billed every {cycle}"
	CadenceMilestone   string `json:"cadenceMilestone"` // e.g. "Per milestone"
	CadenceAdHoc       string `json:"cadenceAdHoc"`     // e.g. "Per occurrence"
	CadenceUnspecified string `json:"cadenceUnspecified"`

	// Compact labels for the BillingKind / AmountBasis cells in the plans table.
	BillingKindOneTime          string `json:"billingKindOneTime"`
	BillingKindRecurring        string `json:"billingKindRecurring"`
	BillingKindContract         string `json:"billingKindContract"`
	BillingKindMilestone        string `json:"billingKindMilestone"`
	BillingKindAdHoc            string `json:"billingKindAdHoc"`
	AmountBasisPerCycle         string `json:"amountBasisPerCycle"`
	AmountBasisTotalPackage     string `json:"amountBasisTotalPackage"`
	AmountBasisDerivedFromLines string `json:"amountBasisDerivedFromLines"`
	AmountBasisPerOccurrence    string `json:"amountBasisPerOccurrence"`

	// Plans table row actions + confirms.
	PlanView            string `json:"planView"`
	PlanEdit            string `json:"planEdit"`
	PlanEditDrawerTitle string `json:"planEditDrawerTitle"`
	PlanDelete          string `json:"planDelete"`
	PlanDeleteTitle     string `json:"planDeleteTitle"`
	PlanDeleteMsg       string `json:"planDeleteMsg"`
	PlanInUseTooltip    string `json:"planInUseTooltip"`

	// Plans table primary action + inline error messages.
	PlanAdd      string `json:"planAdd"`
	PlanRequired string `json:"planRequired"`

	// Product prices table columns.
	ProductPriceColumnProduct   string `json:"productPriceColumnProduct"`
	ProductPriceColumnPrice     string `json:"productPriceColumnPrice"`
	ProductPriceColumnCurrency  string `json:"productPriceColumnCurrency"`
	ProductPriceColumnTreatment string `json:"productPriceColumnTreatment"`
	ProductPriceColumnEffective string `json:"productPriceColumnEffective"`

	// Drawer banners explaining how the per-line price relates to the parent
	// PricePlan.amount_basis. Surfaced above the Price input.
	BasisBannerPerCycle     string `json:"basisBannerPerCycle"`
	BasisBannerTotalPackage string `json:"basisBannerTotalPackage"`
	BasisBannerDerived      string `json:"basisBannerDerived"`

	// Drawer section labels used by the schedule-scoped PPP drawer.
	ProductPriceCatalogSection   string `json:"productPriceCatalogSection"`
	ProductPricePricingSection   string `json:"productPricePricingSection"`
	ProductPriceEffectiveSection string `json:"productPriceEffectiveSection"`

	// Attachments tab label for the price_schedule detail page and the nested
	// price_plan (plan) detail page.
	TabAttachments string `json:"tabAttachments"`
}

type PriceScheduleErrorLabels struct {
	NotFound                   string `json:"notFound"`
	LoadFailed                 string `json:"loadFailed"`
	Unauthorized               string `json:"unauthorized"`
	CreateFailed               string `json:"createFailed"`
	UpdateFailed               string `json:"updateFailed"`
	DeleteFailed               string `json:"deleteFailed"`
	InUse                      string `json:"inUse"`
	PricePlanCreateUnavailable string `json:"pricePlanCreateUnavailable"`
}

// DefaultPriceScheduleLabels returns PriceScheduleLabels with sensible English defaults.
func DefaultPriceScheduleLabels() PriceScheduleLabels {
	return PriceScheduleLabels{
		Page: PriceSchedulePageLabels{
			Title:         "Price Schedules",
			Subtitle:      "Manage your price schedules",
			ActiveTitle:   "Active Price Schedules",
			InactiveTitle: "Inactive Price Schedules",
		},
		Buttons: PriceScheduleButtonLabels{
			View:       "View",
			Add:        "Add Price Schedule",
			Edit:       "Edit Price Schedule",
			Delete:     "Delete Price Schedule",
			BulkDelete: "Delete Price Schedules",
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Columns: PriceScheduleColumnLabels{
			Name:        "Name",
			Description: "Description",
			DateStart:   "Start Date",
			DateEnd:     "End Date",
			Location:    "Location",
			Status:      "Status",
			DateCreated: "Date Created",
			Actions:     "Actions",
		},
		Empty: PriceScheduleEmptyLabels{
			Title:   "No Price Schedules",
			Message: "No price schedules to display.",
		},
		Form: PriceScheduleFormLabels{
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
		Bulk: PriceScheduleBulkLabels{
			DeleteTitle:       "Delete Price Schedules",
			DeleteMessage:     "Permanently delete the selected price schedules? This cannot be undone.",
			ActivateTitle:     "Activate Price Schedules",
			ActivateMessage:   "Activate the selected price schedules?",
			DeactivateTitle:   "Deactivate Price Schedules",
			DeactivateMessage: "Deactivate the selected price schedules?",
		},
		Confirm: PriceScheduleConfirmLabels{
			DeleteTitle:       "Delete Price Schedule",
			DeleteMessage:     "Permanently delete this price schedule? This cannot be undone.",
			ActivateTitle:     "Activate Price Schedule",
			ActivateMessage:   "Activate {{name}}?",
			DeactivateTitle:   "Deactivate Price Schedule",
			DeactivateMessage: "Deactivate {{name}}?",
		},
		Tabs: PriceScheduleTabLabels{
			Info:              "Info",
			PricePlan:         "Plans",
			PricePlanSlug:     "",
			ProductPrices:     "Product Prices",
			Subscriptions:     "Subscriptions",
			SubscriptionsSlug: "",
		},
		Detail: PriceScheduleDetailLabels{
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
		PlanForm: PriceSchedulePlanFormLabels{
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
		Errors: PriceScheduleErrorLabels{
			NotFound:                   "Price schedule not found",
			LoadFailed:                 "Failed to load price schedule",
			Unauthorized:               "You are not authorized to perform this action",
			CreateFailed:               "Failed to create price schedule",
			UpdateFailed:               "Failed to update price schedule",
			DeleteFailed:               "Failed to delete price schedule",
			InUse:                      "This price schedule is in use by active subscriptions and cannot be deleted.",
			PricePlanCreateUnavailable: "Adding a price plan is not available. Please contact support.",
		},
		Filters: PriceScheduleFilterLabels{
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
