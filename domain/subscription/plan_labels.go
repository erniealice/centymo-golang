package subscription

import (
	"strings"
)

// ---------------------------------------------------------------------------
// Plan labels
// ---------------------------------------------------------------------------

// PlanFilterLabels holds translatable labels for the scope filter chip on the
// plan list page (§6.1 of the 2026-04-27 plan-client-scope plan).
type PlanFilterLabels struct {
	ScopeChipLabel string `json:"scopeChipLabel"`
	ScopeMaster    string `json:"scopeMaster"`
	ScopeClient    string `json:"scopeClient"`
	ScopeAll       string `json:"scopeAll"`
}

// PlanLabels holds all translatable strings for the plan module.
type PlanLabels struct {
	Page            PlanPageLabels        `json:"page"`
	Buttons         PlanButtonLabels      `json:"buttons"`
	Columns         PlanColumnLabels      `json:"columns"`
	Empty           PlanEmptyLabels       `json:"empty"`
	Form            PlanFormLabels        `json:"form"`
	Actions         PlanActionLabels      `json:"actions"`
	Bulk            PlanBulkLabels        `json:"bulkActions"`
	Status          PlanStatusLabels      `json:"status"`
	Detail          PlanDetailLabels      `json:"detail"`
	Tabs            PlanTabLabels         `json:"tabs"`
	Confirm         PlanConfirmLabels     `json:"confirm"`
	Errors          PlanErrorLabels       `json:"errors"`
	ProductPlanForm ProductPlanFormLabels `json:"productPlanForm"`
	Filters         PlanFilterLabels      `json:"filters"`
}

// PricePlanFormLabels holds translatable labels for the PricePlan add/edit form.
type PricePlanFormLabels struct {
	Name                string `json:"name"`
	NamePlaceholder     string `json:"namePlaceholder"`
	Description         string `json:"description"`
	DescPlaceholder     string `json:"descriptionPlaceholder"`
	Amount              string `json:"amount"`
	AmountPlaceholder   string `json:"amountPlaceholder"`
	Currency            string `json:"currency"`
	CurrencyPlaceholder string `json:"currencyPlaceholder"`
	DurationValue       string `json:"durationValue"`
	DurationUnit        string `json:"durationUnit"`
	Schedule            string `json:"schedule"`
	SchedulePlaceholder string `json:"schedulePlaceholder"`
	ScheduleSearch      string `json:"scheduleSearch"`
	Location            string `json:"location"`
	LocationPlaceholder string `json:"locationPlaceholder"`
	LocationHintPrefix  string `json:"locationHintPrefix"`
	SelectLocation      string `json:"selectLocation"`
	Active              string `json:"active"`
	PlanLabel           string `json:"planLabel"`
	PlanPlaceholder     string `json:"planPlaceholder"`
	PlanSearch          string `json:"planSearch"`

	// Wave 2 — new billing semantics fields (from lyngua price_plan.json → price_plan.form)
	SectionBasic         string `json:"sectionBasic"`
	SectionPricing       string `json:"sectionPricing"`
	BillingKindLabel     string `json:"billingKindLabel"`
	BillingKindOneTime   string `json:"billingKindOneTime"`
	BillingKindRecurring string `json:"billingKindRecurring"`
	BillingKindContract  string `json:"billingKindContract"`
	BillingKindMilestone string `json:"billingKindMilestone"`
	BillingKindAdHoc     string `json:"billingKindAdHoc"`
	// Per-option hint copy surfaced inline below the billing_kind select as the
	// operator picks. Matches the multi-vertical convention — general/ tier ships
	// neutral phrasing, professional/ overrides with engagement vocabulary.
	BillingKindOneTimeHint      string `json:"billingKindOneTimeHint"`
	BillingKindRecurringHint    string `json:"billingKindRecurringHint"`
	BillingKindContractHint     string `json:"billingKindContractHint"`
	BillingKindMilestoneHint    string `json:"billingKindMilestoneHint"`
	BillingKindAdHocHint        string `json:"billingKindAdHocHint"`
	AmountBasisLabel            string `json:"amountBasisLabel"`
	AmountBasisPerCycle         string `json:"amountBasisPerCycle"`
	AmountBasisTotalPackage     string `json:"amountBasisTotalPackage"`
	AmountBasisDerivedFromLines string `json:"amountBasisDerivedFromLines"`
	AmountBasisPerOccurrence    string `json:"amountBasisPerOccurrence"`
	// Per-option hint copy for amount_basis (mirrors billing_kind pattern).
	AmountBasisPerCycleHint         string `json:"amountBasisPerCycleHint"`
	AmountBasisTotalPackageHint     string `json:"amountBasisTotalPackageHint"`
	AmountBasisDerivedFromLinesHint string `json:"amountBasisDerivedFromLinesHint"`
	AmountBasisPerOccurrenceHint    string `json:"amountBasisPerOccurrenceHint"`
	EntitledOccurrencesLabel        string `json:"entitledOccurrencesLabel"`
	EntitledOccurrencesPlaceholder  string `json:"entitledOccurrencesPlaceholder"`
	EntitledOccurrencesInfo         string `json:"entitledOccurrencesInfo"`
	BillingCycleLabel               string `json:"billingCycleLabel"`
	BillingCyclePlaceholder         string `json:"billingCyclePlaceholder"`
	TermLabel                       string `json:"termLabel"`
	TermPlaceholder                 string `json:"termPlaceholder"`
	TermOpenEndedHelp               string `json:"termOpenEndedHelp"`

	// Field-level info text surfaced via an info button beside each label.
	PlanInfo         string `json:"planInfo"`
	ScheduleInfo     string `json:"scheduleInfo"`
	NameInfo         string `json:"nameInfo"`
	DescriptionInfo  string `json:"descriptionInfo"`
	BillingKindInfo  string `json:"billingKindInfo"`
	AmountBasisInfo  string `json:"amountBasisInfo"`
	AmountInfo       string `json:"amountInfo"`
	CurrencyInfo     string `json:"currencyInfo"`
	BillingCycleInfo string `json:"billingCycleInfo"`
	TermInfo         string `json:"termInfo"`
	ActiveInfo       string `json:"activeInfo"`

	// 2026-04-27 plan-client-scope plan §6.7 — info banner shown above the
	// PricePlan add/edit form when its parent PriceSchedule is client-scoped.
	// Templated via Go's text/template ({{.ClientName}}).
	ParentScheduleClientNotice string `json:"parentScheduleClientNotice"`

	// 2026-04-27 plan-client-scope plan §6.7 — tooltip surfaced beside the
	// readonly Schedule label when the PricePlan's parent Plan is
	// client-scoped (the schedule field is locked to the resolved/derived
	// client schedule). Templated via Go's text/template ({{.ClientName}}).
	ScheduleLockedTooltip string `json:"scheduleLockedTooltip"`
	// 2026-04-28 — info-row hints rendered beneath the readonly Schedule
	// label so the operator knows what happens on save:
	//   ScheduleAutoCreateHint — no client rate card exists yet; one will be
	//     created with this client's name + the lyngua suffix.
	//   ScheduleAutoReuseHint  — an existing client rate card was found; the
	//     new price plan will attach to it.
	// Both templated with {{.ClientName}}.
	ScheduleAutoCreateHint string `json:"scheduleAutoCreateHint"`
	ScheduleAutoReuseHint  string `json:"scheduleAutoReuseHint"`

	// 2026-05-03 — info banner rendered below the readonly Schedule display.
	// ScheduleClientPickerNotice fires when the parent schedule is
	// client-scoped: picker shows only plans assigned to that client.
	// ScheduleGeneralPickerNotice fires when the parent schedule is
	// general-scope: picker shows only general-scope plans (client-specific
	// plans cannot attach to a general schedule).
	ScheduleClientPickerNotice  string `json:"scheduleClientPickerNotice"`
	ScheduleGeneralPickerNotice string `json:"scheduleGeneralPickerNotice"`

	// 2026-04-30 cyclic-subscription-jobs plan §9.4 — client-side block
	// surfaced as a tooltip on the disabled MILESTONE option in the
	// billing_kind dropdown when the parent Plan is cyclic.
	MilestoneCyclicBlock string `json:"milestoneCyclicBlock"`

	// 2026-05-01 ad-hoc-subscription-billing plan §6 — client-side guards
	// surfaced as drawer warnings / tooltips on the disabled options. The
	// server enforces the same rules in validate_ad_hoc.go.
	AdHocPoolNoTemplate           string `json:"adHocPoolNoTemplate"`
	AdHocPerCallNoTemplate        string `json:"adHocPerCallNoTemplate"`
	AdHocNoEntitlement            string `json:"adHocNoEntitlement"`
	AdHocBillingCycleNotAllowed   string `json:"adHocBillingCycleNotAllowed"`
	AdHocVisitsPerCycleNotAllowed string `json:"adHocVisitsPerCycleNotAllowed"`
}

// ---------------------------------------------------------------------------
// Plan form, detail, tabs, confirm sub-labels
// ---------------------------------------------------------------------------

type PlanFormSectionLabels struct {
	Basic    string `json:"basic"`
	Services string `json:"services"`
}

type PlanFormLabels struct {
	Name                string                `json:"name"`
	NamePlaceholder     string                `json:"namePlaceholder"`
	Description         string                `json:"description"`
	DescPlaceholder     string                `json:"descriptionPlaceholder"`
	FulfillmentType     string                `json:"fulfillmentType"`
	Active              string                `json:"active"`
	Products            string                `json:"products"`
	ProductsPlaceholder string                `json:"productsPlaceholder"`
	ProductsSearch      string                `json:"productsSearch"`
	Sections            PlanFormSectionLabels `json:"sections"`

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

type PlanDetailLabels struct {
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

type PlanTabLabels struct {
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
func (t PlanTabLabels) ResolveTabSlug(canonical string) string {
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
func (t PlanTabLabels) CanonicalizeTab(slug string) string {
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

type PlanConfirmLabels struct {
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

func DefaultPlanLabels() PlanLabels {
	return PlanLabels{
		Page: PlanPageLabels{
			Heading:         "Plans",
			HeadingActive:   "Active Plans",
			HeadingInactive: "Inactive Plans",
			Caption:         "Manage your plans",
			CaptionActive:   "Manage your active plans",
			CaptionInactive: "View inactive or archived plans",
		},
		Buttons: PlanButtonLabels{
			AddPlan:       "Add Plan",
			AddPricePlan:  "Add Price Plan",
			EditPricePlan: "Edit Price Plan",
			AddProduct:    "Add Product",
		},
		Columns: PlanColumnLabels{
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
		Empty: PlanEmptyLabels{
			Title:           "No plans found",
			Message:         "No plans to display.",
			ActiveTitle:     "No active plans",
			ActiveMessage:   "Create your first plan to get started.",
			InactiveTitle:   "No inactive plans",
			InactiveMessage: "Discontinued plans will appear here.",
		},
		Form: PlanFormLabels{
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
			Sections: PlanFormSectionLabels{
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
		Actions: PlanActionLabels{
			View:       "View Plan",
			Edit:       "Edit Plan",
			Delete:     "Delete Plan",
			Activate:   "Activate Plan",
			Deactivate: "Deactivate Plan",
		},
		Bulk: PlanBulkLabels{
			Delete: "Delete Selected",
		},
		Status: PlanStatusLabels{
			Activate:   "Activate",
			Deactivate: "Deactivate",
		},
		Detail: PlanDetailLabels{
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
		Tabs: PlanTabLabels{
			Info:          "Information",
			Products:      "Products",
			PricePlan:     "Rate Cards",
			PricePlanSlug: "",
			Attachments:   "Attachments",
			AuditTrail:    "Audit Trail",
			AuditHistory:  "History",
		},
		Confirm: PlanConfirmLabels{
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
		Errors: PlanErrorLabels{
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
		Filters: PlanFilterLabels{
			ScopeChipLabel: "Show:",
			ScopeMaster:    "Master",
			ScopeClient:    "Client-specific",
			ScopeAll:       "All",
		},
	}
}
