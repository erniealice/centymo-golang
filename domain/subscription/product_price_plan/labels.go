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
	BillingTreatmentLabel              string `json:"billing_treatment_label"`
	BillingTreatmentRecurring          string `json:"billing_treatment_recurring"`
	BillingTreatmentRecurringHelp      string `json:"billing_treatment_recurring_help"`
	BillingTreatmentOneTimeInitial     string `json:"billing_treatment_one_time_initial"`
	BillingTreatmentOneTimeInitialHelp string `json:"billing_treatment_one_time_initial_help"`
	BillingTreatmentUsageBased         string `json:"billing_treatment_usage_based"`
	BillingTreatmentUsageBasedHelp     string `json:"billing_treatment_usage_based_help"`
	ProductLabel                       string `json:"product_label"`
	ProductPlaceholder                 string `json:"product_placeholder"`
	PriceLabel                         string `json:"price_label"`
	PricePlaceholder                   string `json:"price_placeholder"`
	CurrencyLabel                      string `json:"currency_label"`
	CurrencyPlaceholder                string `json:"currency_placeholder"`
	DateStartLabel                     string `json:"date_start_label"`
	DateEndLabel                       string `json:"date_end_label"`

	// 20260604-performance-evaluation Phase A — advertised rate band (the
	// Offering's rate band). billing_amount_min/max are the advertised band
	// FLOOR/CEILING in centavos (proto fields 22/23, optional). They are
	// display/quote bounds: a seat's contracted_amount must fall within
	// [min, max] when both are set (band CHECK enforced server-side in the
	// subscription_seat create/update use case — NOT here). billing_amount
	// (field 11) remains the operative contracted rate. Both inputs are
	// optional; leaving them blank means "no advertised band".
	BillingAmountMinLabel       string `json:"billing_amount_min_label"`
	BillingAmountMinPlaceholder string `json:"billing_amount_min_placeholder"`
	BillingAmountMinInfo        string `json:"billing_amount_min_info"`
	BillingAmountMaxLabel       string `json:"billing_amount_max_label"`
	BillingAmountMaxPlaceholder string `json:"billing_amount_max_placeholder"`
	BillingAmountMaxInfo        string `json:"billing_amount_max_info"`
	// Section header for the rate-band fields group on the drawer.
	SectionRateBand string `json:"section_rate_band"`
	// Table column header for the combined band display column.
	ColumnRateBand string `json:"column_rate_band"`

	// Field-level info text surfaced via an info button beside each label.
	ProductInfo          string `json:"product_info"`
	PriceInfo            string `json:"price_info"`
	CurrencyInfo         string `json:"currency_info"`
	BillingTreatmentInfo string `json:"billing_treatment_info"`
	DateStartInfo        string `json:"date_start_info"`
	DateEndInfo          string `json:"date_end_info"`

	// Model D — catalog line picker (replaces product_id with product_plan_id)
	CatalogLineLabel       string `json:"catalog_line_label"`
	CatalogLinePlaceholder string `json:"catalog_line_placeholder"`
	CatalogLineInfo        string `json:"catalog_line_info"`

	// 2026-04-29 milestone-billing plan §5 / Phase D — milestone (job
	// template phase) select. Surfaced when the parent PricePlan has
	// billing_kind = MILESTONE; an empty selection falls through to the
	// first event for the milestone plan.
	MilestonePhaseLabel       string `json:"milestone_phase_label"`
	MilestonePhaseFallthrough string `json:"milestone_phase_fallthrough"`
	MilestonePhaseBillable    string `json:"milestone_phase_billable"`

	// Tax override labels (Phase 5) — optional per-PPP tax overrides.
	SectionTax                  string `json:"section_tax"`
	TaxTreatmentLabel           string `json:"tax_treatment_label"`
	TaxTreatmentPlaceholder     string `json:"tax_treatment_placeholder"`
	TaxTreatmentInfo            string `json:"tax_treatment_info"`
	WithholdingClassLabel       string `json:"withholding_class_label"`
	WithholdingClassPlaceholder string `json:"withholding_class_placeholder"`
	WithholdingClassInfo        string `json:"withholding_class_info"`

	// Read-only parent-PricePlan context block rendered above the editable
	// fields (ppp-parent-context.html). Shared across the PPP drawer and the
	// price-schedule-scoped product-price drawer.
	ParentContext PricePlanParentContextLabels `json:"parent_context"`
}

// PricePlanParentContextLabels labels the read-only "parent context" rows on
// the ppp-parent-context partial. RateCard uses the proto-generic
// "Price Schedule" by default; the professional/education tiers override it to
// "Rate Card" via lyngua.
type PricePlanParentContextLabels struct {
	MoreDetails  string `json:"more_details"`
	RateCard     string `json:"rate_card"`
	BillingModel string `json:"billing_model"`
	AmountBasis  string `json:"amount_basis"`
	BillingCycle string `json:"billing_cycle"`
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
