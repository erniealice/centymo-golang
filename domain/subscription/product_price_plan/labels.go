package product_price_plan

// ---------------------------------------------------------------------------
// ProductPricePlan labels
// ---------------------------------------------------------------------------

// Labels holds all labels for the ProductPricePlan drawer form.
// Wave 2 addition: billing treatment + product/price/currency/date fields.
type Labels struct {
	Form FormLabels `json:"form"`
}

// FormLabels holds translatable labels for the ProductPricePlan
// add/edit drawer form. Keys match lyngua product_price_plan.json → product_price_plan.form.
type FormLabels struct {
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

	// 20260604-performance-evaluation Phase A — advertised rate band (the
	// Offering's rate band). billing_amount_min/max are the advertised band
	// FLOOR/CEILING in centavos (proto fields 22/23, optional). They are
	// display/quote bounds: a seat's contracted_amount must fall within
	// [min, max] when both are set (band CHECK enforced server-side in the
	// subscription_seat create/update use case — NOT here). billing_amount
	// (field 11) remains the operative contracted rate. Both inputs are
	// optional; leaving them blank means "no advertised band".
	BillingAmountMinLabel       string `json:"billingAmountMinLabel"`
	BillingAmountMinPlaceholder string `json:"billingAmountMinPlaceholder"`
	BillingAmountMinInfo        string `json:"billingAmountMinInfo"`
	BillingAmountMaxLabel       string `json:"billingAmountMaxLabel"`
	BillingAmountMaxPlaceholder string `json:"billingAmountMaxPlaceholder"`
	BillingAmountMaxInfo        string `json:"billingAmountMaxInfo"`
	// Section header for the rate-band fields group on the drawer.
	SectionRateBand string `json:"sectionRateBand"`
	// Table column header for the combined band display column.
	ColumnRateBand string `json:"columnRateBand"`

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

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Form: FormLabels{
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
			// 20260604-performance-evaluation Phase A — advertised rate band defaults.
			BillingAmountMinLabel:       "Band floor",
			BillingAmountMinPlaceholder: "0.00",
			BillingAmountMinInfo:        "Advertised band floor in centavos (displayed as amount ÷ 100). The lowest rate a seat on this line may be contracted at. Leave blank for no floor.",
			BillingAmountMaxLabel:       "Band ceiling",
			BillingAmountMaxPlaceholder: "0.00",
			BillingAmountMaxInfo:        "Advertised band ceiling in centavos (displayed as amount ÷ 100). The highest rate a seat on this line may be contracted at. Leave blank for no ceiling.",
			SectionRateBand:             "Advertised rate band",
			ColumnRateBand:              "Rate band",
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
