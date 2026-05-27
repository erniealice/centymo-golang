package centymo

// labels_treasury_method.go — Stage-1 (Wave 4) Go label structs for the
// treasury-domain-rebuild Method management views (collection_method +
// disbursement_method). English defaults inline per the wave brief; full
// lyngua JSON + per-tier seeds are a SEPARATE follow-up wave.
//
// Field-help convention (pages.md §A-1): *Label = field name, *Info = ⓘ
// static guidance. The reactive *OptionHint / per-fragment intro keys are
// stubbed minimally here (Stage 6 lyngua wave fleshes them out).

// ---------------------------------------------------------------------------
// CollectionMethodLabels (selling side)
// ---------------------------------------------------------------------------

// CollectionMethodLabels holds all translatable strings for the collection_method module.
type CollectionMethodLabels struct {
	Page            CollectionMethodPageLabels            `json:"page"`
	Columns         CollectionMethodColumnLabels          `json:"columns"`
	Tabs            CollectionMethodTabLabels             `json:"tabs"`
	Detail          CollectionMethodDetailLabels          `json:"detail"`
	Form            CollectionMethodFormLabels            `json:"form"`
	Fragment        CollectionMethodFragmentLabels        `json:"fragment"`
	Empty           CollectionMethodEmptyLabels           `json:"empty"`
	EligibilityRule CollectionMethodEligibilityRuleLabels `json:"eligibilityRule"`
}

type CollectionMethodPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingDraft    string `json:"headingDraft"`
	HeadingArchived string `json:"headingArchived"`
	Caption         string `json:"caption"`
	AddButton       string `json:"addButton"`
	DetailSubtitle  string `json:"detailSubtitle"`
}

type CollectionMethodColumnLabels struct {
	TemplateCode string `json:"templateCode"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	PostingKind  string `json:"postingKind"`
	AudienceMode string `json:"audienceMode"`
	Lifecycle    string `json:"lifecycle"`
	Source       string `json:"source"`
	Revision     string `json:"revision"`
}

type CollectionMethodTabLabels struct {
	Info          string `json:"info"`
	Eligibility   string `json:"eligibility"`
	Grants        string `json:"grants"`
	SubStatusTags string `json:"subStatusTags"`
	Approvals     string `json:"approvals"`
	Instances     string `json:"instances"`
	Profiles      string `json:"profiles"`
	Transitions   string `json:"transitions"`
	Versions      string `json:"versions"`
	Activity      string `json:"activity"`
	StagePending  string `json:"stagePending"`
}

type CollectionMethodDetailLabels struct {
	OverviewSection string `json:"overviewSection"`
	Name            string `json:"name"`
	TemplateCode    string `json:"templateCode"`
	Category        string `json:"category"`
	PostingKind     string `json:"postingKind"`
	AudienceMode    string `json:"audienceMode"`
	TaxEffectKind   string `json:"taxEffectKind"`
	Lifecycle       string `json:"lifecycle"`
	Source          string `json:"source"`
	Revision        string `json:"revision"`
	VersionStatus   string `json:"versionStatus"`
	BalanceAccount  string `json:"balanceAccount"`
	TargetAccount   string `json:"targetAccount"`
	EligibilityRule string `json:"eligibilityRule"`
	KindSummary     string `json:"kindSummary"`
}

type CollectionMethodFormLabels struct {
	AddTitle  string `json:"addTitle"`
	EditTitle string `json:"editTitle"`

	SectionCommon  string `json:"sectionCommon"`
	SectionGL      string `json:"sectionGl"`
	SectionVersion string `json:"sectionVersion"`
	SectionKind    string `json:"sectionKind"`

	Name             string `json:"name"`
	NamePlaceholder  string `json:"namePlaceholder"`
	NameInfo         string `json:"nameInfo"`
	Category         string `json:"category"`
	CategoryInfo     string `json:"categoryInfo"`
	PostingKind      string `json:"postingKind"`
	PostingKindInfo  string `json:"postingKindInfo"`
	AudienceMode     string `json:"audienceMode"`
	AudienceModeInfo string `json:"audienceModeInfo"`
	TaxEffectKind    string `json:"taxEffectKind"`
	TaxEffectInfo    string `json:"taxEffectInfo"`
	EligibilityRule  string `json:"eligibilityRule"`
	BalanceAccount   string `json:"balanceAccount"`
	TargetAccount    string `json:"targetAccount"`
	Lifecycle        string `json:"lifecycle"`
	LifecycleInfo    string `json:"lifecycleInfo"`
	Source           string `json:"source"`
	TemplateCode     string `json:"templateCode"`
	TemplateCodeInfo string `json:"templateCodeInfo"`
	Revision         string `json:"revision"`
	VersionStatus    string `json:"versionStatus"`
	Supersedes       string `json:"supersedes"`

	// Category option labels (human text for the enum values)
	CategoryStandard string `json:"categoryStandard"`
	CategoryVoucher  string `json:"categoryVoucher"`
	CategoryAdvance  string `json:"categoryAdvance"`
	CategoryCard     string `json:"categoryCard"`
}

// CollectionMethodFragmentLabels covers the 5 §A-2 kind fragments + intros.
type CollectionMethodFragmentLabels struct {
	// Per-fragment plain-language intro blocks (pages.md §A-2 table)
	CardIntro    string `json:"cardIntro"`
	VoucherIntro string `json:"voucherIntro"`
	AdvanceIntro string `json:"advanceIntro"`
	CashIntro    string `json:"cashIntro"`
	BankIntro    string `json:"bankIntro"`

	// Voucher-program fields
	DefaultFaceValue   string `json:"defaultFaceValue"`
	DefaultExpiryDays  string `json:"defaultExpiryDays"`
	AllowedBearerModes string `json:"allowedBearerModes"`

	// Advance-program fields
	AdvanceKind        string `json:"advanceKind"`
	DefaultBalanceAcct string `json:"defaultBalanceAcct"`
	DefaultTargetAcct  string `json:"defaultTargetAcct"`
	DefaultPeriodCount string `json:"defaultPeriodCount"`
	DefaultPeriodUnit  string `json:"defaultPeriodUnit"`
	DefaultProration   string `json:"defaultProration"`

	// Bank-account fields
	BankName      string `json:"bankName"`
	AccountFormat string `json:"accountFormat"`
}

type CollectionMethodEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultCollectionMethodLabels returns English fallback labels.
func DefaultCollectionMethodLabels() CollectionMethodLabels {
	return CollectionMethodLabels{
		Page: CollectionMethodPageLabels{
			Heading:         "Collection Methods",
			HeadingActive:   "Active Collection Methods",
			HeadingDraft:    "Draft Collection Methods",
			HeadingArchived: "Archived Collection Methods",
			Caption:         "Templates that define how you collect from clients.",
			AddButton:       "New Collection Method",
			DetailSubtitle:  "Collection method template",
		},
		Columns: CollectionMethodColumnLabels{
			TemplateCode: "Code",
			Name:         "Name",
			Category:     "Category",
			PostingKind:  "Posting",
			AudienceMode: "Audience",
			Lifecycle:    "Lifecycle",
			Source:       "Source",
			Revision:     "Rev.",
		},
		Tabs: CollectionMethodTabLabels{
			Info:          "Info",
			Eligibility:   "Eligibility Rules",
			Grants:        "Grants",
			SubStatusTags: "Sub-status Tags",
			Approvals:     "Approval Rules",
			Instances:     "Instances",
			Profiles:      "Collection Profiles",
			Transitions:   "Transition Requests",
			Versions:      "Versions",
			Activity:      "Activity",
			StagePending:  "This section is coming in a later release.",
		},
		Detail: CollectionMethodDetailLabels{
			OverviewSection: "Template Overview",
			Name:            "Name",
			TemplateCode:    "Template Code",
			Category:        "Category",
			PostingKind:     "Posting Kind",
			AudienceMode:    "Audience Mode",
			TaxEffectKind:   "Tax Effect",
			Lifecycle:       "Lifecycle",
			Source:          "Source",
			Revision:        "Revision",
			VersionStatus:   "Version Status",
			BalanceAccount:  "Balance Account",
			TargetAccount:   "Target Account",
			EligibilityRule: "Default Eligibility Rule",
			KindSummary:     "Kind Configuration",
		},
		Form: CollectionMethodFormLabels{
			AddTitle:         "New Collection Method",
			EditTitle:        "Edit Collection Method",
			SectionCommon:    "Method Details",
			SectionGL:        "Accounting Defaults",
			SectionVersion:   "Versioning",
			SectionKind:      "Kind-specific Settings",
			Name:             "Name",
			NamePlaceholder:  "e.g. Retail Gift Card Spring 2026",
			NameInfo:         "A short, recognizable name for this method template.",
			Category:         "Category",
			CategoryInfo:     "What kind of method this is. Changing it updates the settings below.",
			PostingKind:      "Posting Kind",
			PostingKindInfo:  "How collections through this method are recorded.",
			AudienceMode:     "Audience Mode",
			AudienceModeInfo: "Who is allowed to use this method.",
			TaxEffectKind:    "Tax Effect",
			TaxEffectInfo:    "Whether amounts include or exclude tax.",
			EligibilityRule:  "Default Eligibility Rule",
			BalanceAccount:   "Balance Account",
			TargetAccount:    "Target Account",
			Lifecycle:        "Lifecycle",
			LifecycleInfo:    "New templates start as Draft.",
			Source:           "Source",
			TemplateCode:     "Template Code",
			TemplateCodeInfo: "Stable code shared across revisions of this template.",
			Revision:         "Revision",
			VersionStatus:    "Version Status",
			Supersedes:       "Supersedes",
			CategoryStandard: "Standard",
			CategoryVoucher:  "Voucher Program",
			CategoryAdvance:  "Advance Template",
			CategoryCard:     "Card",
		},
		Fragment: CollectionMethodFragmentLabels{
			CardIntro:          "This sets up which cards you accept. When a customer saves a card, it's stored against that customer — the card number, last four, and expiry stay with the saved card, not here.",
			VoucherIntro:       "This is a voucher program — the template for a batch of vouchers, like a gift-card campaign. The settings below are just the starting defaults for each voucher you issue.",
			AdvanceIntro:       "This is an advance template — for customer prepayments. The settings below are the starting defaults for each prepayment you record.",
			CashIntro:          "Cash methods aren't issued to individual customers — they're a plain cash flow. There's nothing extra to set up here.",
			BankIntro:          "This sets up the kind of bank account you accept. A customer's actual account details are saved against that customer, not here.",
			DefaultFaceValue:   "Default Face Value",
			DefaultExpiryDays:  "Default Expiry (days)",
			AllowedBearerModes: "Allowed Bearer Modes",
			AdvanceKind:        "Advance Kind",
			DefaultBalanceAcct: "Default Balance Account",
			DefaultTargetAcct:  "Default Target Account",
			DefaultPeriodCount: "Default Period Count",
			DefaultPeriodUnit:  "Default Period Unit",
			DefaultProration:   "Default Proration Policy",
			BankName:           "Bank Name",
			AccountFormat:      "Account Format Rules",
		},
		Empty: CollectionMethodEmptyLabels{
			Title:   "No collection methods yet",
			Message: "Create a method template to start collecting from clients.",
		},
		EligibilityRule: DefaultCollectionMethodEligibilityRuleLabels(),
	}
}

// ---------------------------------------------------------------------------
// CollectionMethodEligibilityRuleLabels (Stage 2 — Eligibility Rules tab)
// ---------------------------------------------------------------------------

// CollectionMethodEligibilityRuleLabels holds translatable strings for the
// eligibility-rule CRUD sub-module (pages.md §B-5 tab #2 + drawer-form fields
// per entities.md §E-3). JSON tags match collection_method_eligibility_rule.json.
type CollectionMethodEligibilityRuleLabels struct {
	Tab     CollectionMethodEligibilityRuleTabLabels    `json:"tab"`
	Columns CollectionMethodEligibilityRuleColumnLabels `json:"columns"`
	Form    CollectionMethodEligibilityRuleFormLabels   `json:"form"`
	Empty   CollectionMethodEligibilityRuleEmptyLabels  `json:"empty"`
	Delete  CollectionMethodEligibilityRuleDeleteLabels `json:"delete"`
}

type CollectionMethodEligibilityRuleTabLabels struct {
	Title string `json:"title"`
}

type CollectionMethodEligibilityRuleColumnLabels struct {
	Name             string `json:"name"`
	BearerMode       string `json:"bearerMode"`
	ValidFrom        string `json:"validFrom"`
	ValidUntil       string `json:"validUntil"`
	StackingPolicy   string `json:"stackingPolicy"`
	JurisdictionCode string `json:"jurisdictionCode"`
	MaxPerInstance   string `json:"maxPerInstance"`
	MaxPerClient     string `json:"maxPerClient"`
}

type CollectionMethodEligibilityRuleFormLabels struct {
	AddTitle  string `json:"addTitle"`
	EditTitle string `json:"editTitle"`

	Name                          string `json:"name"`
	NamePlaceholder               string `json:"namePlaceholder"`
	NameInfo                      string `json:"nameInfo"`
	BearerMode                    string `json:"bearerMode"`
	BearerModeInfo                string `json:"bearerModeInfo"`
	BearerModeHolderBound         string `json:"bearerModeHolderBound"`
	BearerModeTransferable        string `json:"bearerModeTransferable"`
	ValidFromDate                 string `json:"validFromDate"`
	ValidFromDateInfo             string `json:"validFromDateInfo"`
	ValidUntilDate                string `json:"validUntilDate"`
	ValidUntilDateInfo            string `json:"validUntilDateInfo"`
	ExpiryDaysAfterIssuance       string `json:"expiryDaysAfterIssuance"`
	ExpiryDaysInfo                string `json:"expiryDaysInfo"`
	MinAmountCentavos             string `json:"minAmountCentavos"`
	MinAmountInfo                 string `json:"minAmountInfo"`
	MaxAmountCentavos             string `json:"maxAmountCentavos"`
	MaxAmountInfo                 string `json:"maxAmountInfo"`
	StackingPolicy                string `json:"stackingPolicy"`
	StackingPolicyInfo            string `json:"stackingPolicyInfo"`
	StackingExclusive             string `json:"stackingExclusive"`
	StackingStackable             string `json:"stackingStackable"`
	StackingFirstOnly             string `json:"stackingFirstOnly"`
	JurisdictionCode              string `json:"jurisdictionCode"`
	JurisdictionInfo              string `json:"jurisdictionInfo"`
	MaxRedemptionsPerInstance     string `json:"maxRedemptionsPerInstance"`
	MaxRedemptionsPerInstanceInfo string `json:"maxRedemptionsPerInstanceInfo"`
	MaxRedemptionsPerClient       string `json:"maxRedemptionsPerClient"`
	MaxRedemptionsPerClientInfo   string `json:"maxRedemptionsPerClientInfo"`
	TermsURL                      string `json:"termsUrl"`
	TermsURLInfo                  string `json:"termsUrlInfo"`
	TermsSummary                  string `json:"termsSummary"`
	TermsSummaryInfo              string `json:"termsSummaryInfo"`
}

type CollectionMethodEligibilityRuleEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type CollectionMethodEligibilityRuleDeleteLabels struct {
	Confirm string `json:"confirm"`
	Button  string `json:"button"`
}

// DefaultCollectionMethodEligibilityRuleLabels returns English fallback labels.
func DefaultCollectionMethodEligibilityRuleLabels() CollectionMethodEligibilityRuleLabels {
	return CollectionMethodEligibilityRuleLabels{
		Tab: CollectionMethodEligibilityRuleTabLabels{
			Title: "Eligibility Rules",
		},
		Columns: CollectionMethodEligibilityRuleColumnLabels{
			Name:             "Name",
			BearerMode:       "Bearer Mode",
			ValidFrom:        "Valid From",
			ValidUntil:       "Valid Until",
			StackingPolicy:   "Stacking",
			JurisdictionCode: "Jurisdiction",
			MaxPerInstance:   "Max / Instance",
			MaxPerClient:     "Max / Client",
		},
		Form: CollectionMethodEligibilityRuleFormLabels{
			AddTitle:  "New Eligibility Rule",
			EditTitle: "Edit Eligibility Rule",

			Name:                          "Name",
			NamePlaceholder:               "e.g. Gift Card Standard Terms",
			NameInfo:                      "A short internal name for this rule. Operators see it; customers do not.",
			BearerMode:                    "Bearer Mode",
			BearerModeInfo:                "Whether the issued instrument is tied to one customer or can be transferred to another.",
			BearerModeHolderBound:         "Holder bound (non-transferable)",
			BearerModeTransferable:        "Holder transferable",
			ValidFromDate:                 "Valid From",
			ValidFromDateInfo:             "Optional. The earliest date an instrument under this rule can be used.",
			ValidUntilDate:                "Valid Until",
			ValidUntilDateInfo:            "Optional. The last date an instrument under this rule can be used. Leave blank for no hard cutoff.",
			ExpiryDaysAfterIssuance:       "Expires After (days)",
			ExpiryDaysInfo:                "Optional. Number of days after issuance before the instrument expires. Overrides Valid Until when both are set.",
			MinAmountCentavos:             "Minimum Amount",
			MinAmountInfo:                 "Optional. The smallest face value allowed when issuing under this rule.",
			MaxAmountCentavos:             "Maximum Amount",
			MaxAmountInfo:                 "Optional. The largest face value allowed when issuing under this rule.",
			StackingPolicy:                "Stacking Policy",
			StackingPolicyInfo:            "How this instrument interacts with other discounts or instruments at the point of redemption.",
			StackingExclusive:             "Exclusive — cannot be combined with others",
			StackingStackable:             "Stackable — can be combined freely",
			StackingFirstOnly:             "First only — only applied if no other discount is active",
			JurisdictionCode:              "Jurisdiction",
			JurisdictionInfo:              "Optional ISO 3166-2 code (e.g. PH, US-CA) if this rule applies to a specific region.",
			MaxRedemptionsPerInstance:     "Max Uses per Instrument",
			MaxRedemptionsPerInstanceInfo: "Optional. How many times a single issued instrument can be redeemed. Leave blank for unlimited.",
			MaxRedemptionsPerClient:       "Max Uses per Customer",
			MaxRedemptionsPerClientInfo:   "Optional. How many times any one customer can use instruments under this rule. Leave blank for unlimited.",
			TermsURL:                      "Terms URL",
			TermsURLInfo:                  "Optional link to the full terms and conditions for this rule.",
			TermsSummary:                  "Terms Summary",
			TermsSummaryInfo:              "Optional plain-language summary of the key conditions.",
		},
		Empty: CollectionMethodEligibilityRuleEmptyLabels{
			Title:   "No eligibility rules yet",
			Message: "Create an eligibility rule to control how instruments issued under this method can be used.",
		},
		Delete: CollectionMethodEligibilityRuleDeleteLabels{
			Confirm: "Are you sure you want to delete this eligibility rule? This cannot be undone.",
			Button:  "Delete Rule",
		},
	}
}

// ---------------------------------------------------------------------------
// DisbursementMethodLabels (buying side; D-4.9: no audience_mode/grants)
// ---------------------------------------------------------------------------

// DisbursementMethodLabels holds all translatable strings for the disbursement_method module.
type DisbursementMethodLabels struct {
	Page     DisbursementMethodPageLabels     `json:"page"`
	Columns  DisbursementMethodColumnLabels   `json:"columns"`
	Tabs     DisbursementMethodTabLabels      `json:"tabs"`
	Detail   DisbursementMethodDetailLabels   `json:"detail"`
	Form     DisbursementMethodFormLabels     `json:"form"`
	Fragment DisbursementMethodFragmentLabels `json:"fragment"`
	Empty    DisbursementMethodEmptyLabels    `json:"empty"`
}

type DisbursementMethodPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingDraft    string `json:"headingDraft"`
	HeadingArchived string `json:"headingArchived"`
	Caption         string `json:"caption"`
	AddButton       string `json:"addButton"`
	DetailSubtitle  string `json:"detailSubtitle"`
}

type DisbursementMethodColumnLabels struct {
	TemplateCode string `json:"templateCode"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	PostingKind  string `json:"postingKind"`
	Lifecycle    string `json:"lifecycle"`
	Source       string `json:"source"`
	Revision     string `json:"revision"`
}

type DisbursementMethodTabLabels struct {
	Info         string `json:"info"`
	Approvals    string `json:"approvals"`
	Instances    string `json:"instances"`
	Profiles     string `json:"profiles"`
	Transitions  string `json:"transitions"`
	Versions     string `json:"versions"`
	Activity     string `json:"activity"`
	StagePending string `json:"stagePending"`
}

type DisbursementMethodDetailLabels struct {
	OverviewSection string `json:"overviewSection"`
	Name            string `json:"name"`
	TemplateCode    string `json:"templateCode"`
	Category        string `json:"category"`
	PostingKind     string `json:"postingKind"`
	TaxEffectKind   string `json:"taxEffectKind"`
	Lifecycle       string `json:"lifecycle"`
	Source          string `json:"source"`
	Revision        string `json:"revision"`
	VersionStatus   string `json:"versionStatus"`
	BalanceAccount  string `json:"balanceAccount"`
	TargetAccount   string `json:"targetAccount"`
	KindSummary     string `json:"kindSummary"`
}

type DisbursementMethodFormLabels struct {
	AddTitle  string `json:"addTitle"`
	EditTitle string `json:"editTitle"`

	SectionCommon  string `json:"sectionCommon"`
	SectionGL      string `json:"sectionGl"`
	SectionVersion string `json:"sectionVersion"`
	SectionKind    string `json:"sectionKind"`

	Name             string `json:"name"`
	NamePlaceholder  string `json:"namePlaceholder"`
	NameInfo         string `json:"nameInfo"`
	Category         string `json:"category"`
	CategoryInfo     string `json:"categoryInfo"`
	PostingKind      string `json:"postingKind"`
	PostingKindInfo  string `json:"postingKindInfo"`
	TaxEffectKind    string `json:"taxEffectKind"`
	TaxEffectInfo    string `json:"taxEffectInfo"`
	BalanceAccount   string `json:"balanceAccount"`
	TargetAccount    string `json:"targetAccount"`
	Lifecycle        string `json:"lifecycle"`
	LifecycleInfo    string `json:"lifecycleInfo"`
	Source           string `json:"source"`
	TemplateCode     string `json:"templateCode"`
	TemplateCodeInfo string `json:"templateCodeInfo"`
	Revision         string `json:"revision"`
	VersionStatus    string `json:"versionStatus"`
	Supersedes       string `json:"supersedes"`

	// Category option labels — buying side: bank-account / check / advance
	CategoryStandard string `json:"categoryStandard"`
	CategoryAdvance  string `json:"categoryAdvance"`
	CategoryCard     string `json:"categoryCard"`
	CategoryVoucher  string `json:"categoryVoucher"`
}

// DisbursementMethodFragmentLabels covers buying-side fragments
// (bank-account / check / advance-program per §A-2 disbursement side).
type DisbursementMethodFragmentLabels struct {
	BankIntro    string `json:"bankIntro"`
	CheckIntro   string `json:"checkIntro"`
	AdvanceIntro string `json:"advanceIntro"`

	// Bank-account fields
	BankName      string `json:"bankName"`
	AccountFormat string `json:"accountFormat"`

	// Check fields
	CheckSeries   string `json:"checkSeries"`
	SigningPolicy string `json:"signingPolicy"`

	// Advance-program fields
	AdvanceKind        string `json:"advanceKind"`
	DefaultBalanceAcct string `json:"defaultBalanceAcct"`
	DefaultTargetAcct  string `json:"defaultTargetAcct"`
	DefaultPeriodCount string `json:"defaultPeriodCount"`
	DefaultPeriodUnit  string `json:"defaultPeriodUnit"`
	DefaultProration   string `json:"defaultProration"`
}

type DisbursementMethodEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultDisbursementMethodLabels returns English fallback labels.
func DefaultDisbursementMethodLabels() DisbursementMethodLabels {
	return DisbursementMethodLabels{
		Page: DisbursementMethodPageLabels{
			Heading:         "Payment Methods",
			HeadingActive:   "Active Payment Methods",
			HeadingDraft:    "Draft Payment Methods",
			HeadingArchived: "Archived Payment Methods",
			Caption:         "Templates that define how you pay suppliers.",
			AddButton:       "New Payment Method",
			DetailSubtitle:  "Payment method template",
		},
		Columns: DisbursementMethodColumnLabels{
			TemplateCode: "Code",
			Name:         "Name",
			Category:     "Category",
			PostingKind:  "Posting",
			Lifecycle:    "Lifecycle",
			Source:       "Source",
			Revision:     "Rev.",
		},
		Tabs: DisbursementMethodTabLabels{
			Info:         "Info",
			Approvals:    "Approval Rules",
			Instances:    "Instances",
			Profiles:     "Disbursement Profiles",
			Transitions:  "Transition Requests",
			Versions:     "Versions",
			Activity:     "Activity",
			StagePending: "This section is coming in a later release.",
		},
		Detail: DisbursementMethodDetailLabels{
			OverviewSection: "Template Overview",
			Name:            "Name",
			TemplateCode:    "Template Code",
			Category:        "Category",
			PostingKind:     "Posting Kind",
			TaxEffectKind:   "Tax Effect",
			Lifecycle:       "Lifecycle",
			Source:          "Source",
			Revision:        "Revision",
			VersionStatus:   "Version Status",
			BalanceAccount:  "Balance Account",
			TargetAccount:   "Target Account",
			KindSummary:     "Kind Configuration",
		},
		Form: DisbursementMethodFormLabels{
			AddTitle:         "New Payment Method",
			EditTitle:        "Edit Payment Method",
			SectionCommon:    "Method Details",
			SectionGL:        "Accounting Defaults",
			SectionVersion:   "Versioning",
			SectionKind:      "Kind-specific Settings",
			Name:             "Name",
			NamePlaceholder:  "e.g. Corporate Check Series A",
			NameInfo:         "A short, recognizable name for this method template.",
			Category:         "Category",
			CategoryInfo:     "What kind of method this is. Changing it updates the settings below.",
			PostingKind:      "Posting Kind",
			PostingKindInfo:  "How disbursements through this method are recorded.",
			TaxEffectKind:    "Tax Effect",
			TaxEffectInfo:    "Whether amounts include or exclude tax.",
			BalanceAccount:   "Balance Account",
			TargetAccount:    "Target Account",
			Lifecycle:        "Lifecycle",
			LifecycleInfo:    "New templates start as Draft.",
			Source:           "Source",
			TemplateCode:     "Template Code",
			TemplateCodeInfo: "Stable code shared across revisions of this template.",
			Revision:         "Revision",
			VersionStatus:    "Version Status",
			Supersedes:       "Supersedes",
			CategoryStandard: "Standard / Bank Account",
			CategoryAdvance:  "Advance Template",
			CategoryCard:     "Card",
			CategoryVoucher:  "Voucher Program",
		},
		Fragment: DisbursementMethodFragmentLabels{
			BankIntro:          "This sets up the kind of bank account you pay from or to. A supplier's actual account details are saved against that supplier, not here.",
			CheckIntro:         "This sets up a check series and signing policy. Individual checks are issued from this series.",
			AdvanceIntro:       "This is an advance template — for prepayments to suppliers. The settings below are the starting defaults for each prepayment you record.",
			BankName:           "Bank Name",
			AccountFormat:      "Account Format Rules",
			CheckSeries:        "Check Series",
			SigningPolicy:      "Signing Policy",
			AdvanceKind:        "Advance Kind",
			DefaultBalanceAcct: "Default Balance Account",
			DefaultTargetAcct:  "Default Target Account",
			DefaultPeriodCount: "Default Period Count",
			DefaultPeriodUnit:  "Default Period Unit",
			DefaultProration:   "Default Proration Policy",
		},
		Empty: DisbursementMethodEmptyLabels{
			Title:   "No payment methods yet",
			Message: "Create a method template to start paying suppliers.",
		},
	}
}
