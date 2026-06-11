package subscription

// ---------------------------------------------------------------------------
// ProductPricePlan labels
// ---------------------------------------------------------------------------

// ProductPricePlanLabels holds all labels for the ProductPricePlan drawer form.
// Wave 2 addition: billing treatment + product/price/currency/date fields.
type ProductPricePlanLabels struct {
	Form ProductPricePlanFormLabels `json:"form"`
}

// ProductPricePlanFormLabels holds translatable labels for the ProductPricePlan
// add/edit drawer form. Keys match lyngua product_price_plan.json → product_price_plan.form.
type ProductPricePlanFormLabels struct {
	BillingTreatmentLabel              string `json:"billingTreatmentLabel"`
	BillingTreatmentRecurring          string `json:"billingTreatmentRecurring"`
	BillingTreatmentRecurringHelp      string `json:"billingTreatmentRecurringHelp"`
	BillingTreatmentOneTimeInitial     string `json:"billingTreatmentOneTimeInitial"`
	BillingTreatmentOneTimeInitialHelp string `json:"billingTreatmentOneTimeInitialHelp"`
	BillingTreatmentUsageBased         string `json:"billingTreatmentUsageBased"`
	BillingTreatmentUsageBasedHelp     string `json:"billingTreatmentUsageBasedHelp"`
	ProductLabel                       string `json:"productLabel"`
	ProductPlaceholder                 string `json:"productPlaceholder"`
	PriceLabel                         string `json:"priceLabel"`
	PricePlaceholder                   string `json:"pricePlaceholder"`
	CurrencyLabel                      string `json:"currencyLabel"`
	CurrencyPlaceholder                string `json:"currencyPlaceholder"`
	DateStartLabel                     string `json:"dateStartLabel"`
	DateEndLabel                       string `json:"dateEndLabel"`

	// Field-level info text surfaced via an info button beside each label.
	ProductInfo          string `json:"productInfo"`
	PriceInfo            string `json:"priceInfo"`
	CurrencyInfo         string `json:"currencyInfo"`
	BillingTreatmentInfo string `json:"billingTreatmentInfo"`
	DateStartInfo        string `json:"dateStartInfo"`
	DateEndInfo          string `json:"dateEndInfo"`

	// Model D — catalog line picker (replaces product_id with product_plan_id)
	CatalogLineLabel       string `json:"catalogLineLabel"`
	CatalogLinePlaceholder string `json:"catalogLinePlaceholder"`
	CatalogLineInfo        string `json:"catalogLineInfo"`

	// 2026-04-29 milestone-billing plan §5 / Phase D — milestone (job
	// template phase) select. Surfaced when the parent PricePlan has
	// billing_kind = MILESTONE; an empty selection falls through to the
	// first event for the milestone plan.
	MilestonePhaseLabel       string `json:"milestonePhaseLabel"`
	MilestonePhaseFallthrough string `json:"milestonePhaseFallthrough"`
	MilestonePhaseBillable    string `json:"milestonePhaseBillable"`

	// Tax override labels (Phase 5) — optional per-PPP tax overrides.
	SectionTax                  string `json:"sectionTax"`
	TaxTreatmentLabel           string `json:"taxTreatmentLabel"`
	TaxTreatmentPlaceholder     string `json:"taxTreatmentPlaceholder"`
	TaxTreatmentInfo            string `json:"taxTreatmentInfo"`
	WithholdingClassLabel       string `json:"withholdingClassLabel"`
	WithholdingClassPlaceholder string `json:"withholdingClassPlaceholder"`
	WithholdingClassInfo        string `json:"withholdingClassInfo"`

	// Read-only parent-PricePlan context block rendered above the editable
	// fields (ppp-parent-context.html). Shared across the PPP drawer and the
	// price-schedule-scoped product-price drawer.
	ParentContext PricePlanParentContextLabels `json:"parentContext"`
}

// PricePlanParentContextLabels labels the read-only "parent context" rows on
// the ppp-parent-context partial. RateCard uses the proto-generic
// "Price Schedule" by default; the professional/education tiers override it to
// "Rate Card" via lyngua.
type PricePlanParentContextLabels struct {
	MoreDetails  string `json:"moreDetails"`
	RateCard     string `json:"rateCard"`
	BillingModel string `json:"billingModel"`
	AmountBasis  string `json:"amountBasis"`
	BillingCycle string `json:"billingCycle"`
	Term         string `json:"term"`
	Currency     string `json:"currency"`
}

// DefaultProductPricePlanLabels returns ProductPricePlanLabels with sensible English defaults.
func DefaultProductPricePlanLabels() ProductPricePlanLabels {
	return ProductPricePlanLabels{
		Form: ProductPricePlanFormLabels{
			BillingTreatmentLabel:              "Billing treatment",
			BillingTreatmentRecurring:          "Every cycle",
			BillingTreatmentRecurringHelp:      "Charge this line every billing cycle",
			BillingTreatmentOneTimeInitial:     "First cycle only",
			BillingTreatmentOneTimeInitialHelp: "Charge once on the first invoice (setup fees, welcome gifts)",
			BillingTreatmentUsageBased:         "On use",
			BillingTreatmentUsageBasedHelp:     "Charge when consumed or performed",
			ProductLabel:                       "Product",
			ProductPlaceholder:                 "Select a product",
			PriceLabel:                         "Price",
			PricePlaceholder:                   "0.00",
			CurrencyLabel:                      "Currency",
			CurrencyPlaceholder:                "e.g. PHP",
			DateStartLabel:                     "Effective from",
			DateEndLabel:                       "Effective until",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			ProductInfo:          "The product this price applies to.",
			PriceInfo:            "Price in centavos. Displayed as amount ÷ 100.",
			CurrencyInfo:         "Currency applied to this product price.",
			BillingTreatmentInfo: "Every cycle = charged each billing cycle. First cycle only = setup fee. On use = charged when consumed.",
			DateStartInfo:        "Date from which this product price is effective.",
			DateEndInfo:          "Last date this product price is effective. Leave empty for no end date.",
			// Model D — catalog line picker defaults
			CatalogLineLabel:       "Catalog line",
			CatalogLinePlaceholder: "Select a line from the plan's catalog",
			CatalogLineInfo:        "Prices the chosen catalog line from the parent plan. If the line has a variant, that variant is priced.",
			// 2026-04-29 milestone-billing plan §5 — milestone phase select.
			MilestonePhaseLabel:       "Milestone phase",
			MilestonePhaseFallthrough: "Falls through to first event",
			MilestonePhaseBillable:    "billable",
			// Parent-context block — proto-generic defaults; tiers override
			// RateCard to "Rate Card" via lyngua professional/education.
			ParentContext: PricePlanParentContextLabels{
				MoreDetails:  "More details",
				RateCard:     "Price Schedule",
				BillingModel: "Billing model",
				AmountBasis:  "Amount basis",
				BillingCycle: "Billing cycle",
				Term:         "Term",
				Currency:     "Currency",
			},
		},
	}
}

type PlanPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type PlanButtonLabels struct {
	AddPlan       string `json:"addPlan"`
	AddPricePlan  string `json:"addPricePlan"`
	EditPricePlan string `json:"editPricePlan"`
	AddProduct    string `json:"addProduct"`
}

type PlanColumnLabels struct {
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

type PlanEmptyLabels struct {
	Title           string `json:"title"`
	Message         string `json:"message"`
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type PlanActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type PlanBulkLabels struct {
	Delete string `json:"delete"`
}

type PlanStatusLabels struct {
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type PlanErrorLabels struct {
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
