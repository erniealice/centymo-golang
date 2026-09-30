package agreement_line_term

// Labels holds every translatable string of the agreement charge-terms tab.
// Lyngua file: general/agreement_line_term.json, root key "agreement_line_term" (leasing overlays values only).
// Every field has a general-tier Lyngua key; the Go defaults equal the general tier
// (a reflective test compares every string field, C14).
type Labels struct {
	Tabs    TabsLabels    `json:"tabs"`
	Page    PageLabels    `json:"page"`
	Columns ColumnsLabels `json:"columns"`
	Buttons ButtonsLabels `json:"buttons"`
	Form    FormLabels    `json:"form"`
	Enums   EnumsLabels   `json:"enums"`
	Detail  DetailLabels  `json:"detail"`
	Confirm ConfirmLabels `json:"confirm"`
	Empty   EmptyLabels   `json:"empty"`
}

type TabsLabels struct {
	ChargeTerms string `json:"charge_terms"`
}

type PageLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type ColumnsLabels struct {
	Line          string `json:"line"`
	Policy        string `json:"policy"`
	Version       string `json:"version"`
	Markup        string `json:"markup"`
	EffectiveFrom string `json:"effective_from"`
	EffectiveTo   string `json:"effective_to"`
	Origin        string `json:"origin"`
	AcceptedOn    string `json:"accepted_on"`
	AcceptedBy    string `json:"accepted_by"`
}

type ButtonsLabels struct {
}

type FormLabels struct {
}

type EnumsLabels struct {
	OriginCopied     string `json:"origin_copied"`
	OriginNegotiated string `json:"origin_negotiated"`
}

type DetailLabels struct {
	OpenEnded    string `json:"open_ended"`
	VersionLabel string `json:"version_label"`
	PinnedNotice string `json:"pinned_notice"`
	LoadFailed   string `json:"load_failed"`
}

type ConfirmLabels struct {
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultLabels returns the English (general tier) defaults.
func DefaultLabels() Labels {
	return Labels{
		Tabs: TabsLabels{
			ChargeTerms: "Charge terms",
		},
		Page: PageLabels{
			Title:    "Charge terms",
			Subtitle: "The charge policy version each priced line was signed with",
		},
		Columns: ColumnsLabels{
			Line:          "Line",
			Policy:        "Charge policy",
			Version:       "Version",
			Markup:        "Markup",
			EffectiveFrom: "From",
			EffectiveTo:   "To",
			Origin:        "Origin",
			AcceptedOn:    "Accepted on",
			AcceptedBy:    "Accepted by",
		},
		Buttons: ButtonsLabels{},
		Form:    FormLabels{},
		Enums: EnumsLabels{
			OriginCopied:     "Copied from the price plan",
			OriginNegotiated: "Negotiated",
		},
		Detail: DetailLabels{
			OpenEnded:    "Open-ended",
			VersionLabel: "v{0}",
			PinnedNotice: "Each term is pinned to the policy version in force when the agreement started.",
			LoadFailed:   "Charge terms could not be loaded. Try again.",
		},
		Confirm: ConfirmLabels{},
		Empty: EmptyLabels{
			Title:   "No charge terms",
			Message: "This agreement has no lines with a charge policy.",
		},
	}
}
