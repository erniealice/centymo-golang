package allocation_batch

// Labels holds every translatable string of the cost allocation (allocation_batch) screens.
// Lyngua file: general/allocation_batch.json, root key "allocation_batch" (leasing overlays values only).
// Every field has a general-tier Lyngua key; the Go defaults equal the general tier
// (a reflective test compares every string field, C14).
type Labels struct {
	Page    PageLabels    `json:"page"`
	Tabs    TabsLabels    `json:"tabs"`
	Columns ColumnsLabels `json:"columns"`
	Buttons ButtonsLabels `json:"buttons"`
	Form    FormLabels    `json:"form"`
	Enums   EnumsLabels   `json:"enums"`
	Detail  DetailLabels  `json:"detail"`
	Confirm ConfirmLabels `json:"confirm"`
	Empty   EmptyLabels   `json:"empty"`
	Errors  ErrorsLabels  `json:"errors"`
}

type PageLabels struct {
	Title           string `json:"title"`
	Subtitle        string `json:"subtitle"`
	TitleDraft      string `json:"title_draft"`
	TitlePublished  string `json:"title_published"`
	TitleSuperseded string `json:"title_superseded"`
}

type TabsLabels struct {
	Shares  string `json:"shares"`
	History string `json:"history"`
	Charges string `json:"charges"`
}

type ColumnsLabels struct {
	Revision     string `json:"revision"`
	Status       string `json:"status"`
	ShareKind    string `json:"share_kind"`
	Subscription string `json:"subscription"`
	Client       string `json:"client"`
	Numerator    string `json:"numerator"`
	Denominator  string `json:"denominator"`
	Percentage   string `json:"percentage"`
	Amount       string `json:"amount"`
	Rule         string `json:"rule"`
	PublishedOn  string `json:"published_on"`
	PublishedBy  string `json:"published_by"`
	PreparedBy   string `json:"prepared_by"`
	Total        string `json:"total"`
}

type ButtonsLabels struct {
	Allocate    string `json:"allocate"`
	AddShare    string `json:"add_share"`
	RemoveShare string `json:"remove_share"`
	Preview     string `json:"preview"`
	SaveDraft   string `json:"save_draft"`
	Publish     string `json:"publish"`
	Revise      string `json:"revise"`
	ViewCharges string `json:"view_charges"`
}

type FormLabels struct {
	ShareKindLabel          string `json:"share_kind_label"`
	SubscriptionLabel       string `json:"subscription_label"`
	SubscriptionPlaceholder string `json:"subscription_placeholder"`
	NumeratorLabel          string `json:"numerator_label"`
	NumeratorPlaceholder    string `json:"numerator_placeholder"`
	DenominatorLabel        string `json:"denominator_label"`
	RuleLabel               string `json:"rule_label"`
	PreviewTitle            string `json:"preview_title"`
	OwnUseRowLabel          string `json:"own_use_row_label"`
	VacancyRowLabel         string `json:"vacancy_row_label"`
	CommonLossRowLabel      string `json:"common_loss_row_label"`
}

type EnumsLabels struct {
	StatusDraft          string `json:"status_draft"`
	StatusPublished      string `json:"status_published"`
	StatusSuperseded     string `json:"status_superseded"`
	ShareKindRecoverable string `json:"share_kind_recoverable"`
	ShareKindOwnUse      string `json:"share_kind_own_use"`
	ShareKindVacancy     string `json:"share_kind_vacancy"`
	ShareKindCommonLoss  string `json:"share_kind_common_loss"`
	RuleLargestRemainder string `json:"rule_largest_remainder"`
}

type DetailLabels struct {
	TotalAllocated string `json:"total_allocated"`
	SourceAmount   string `json:"source_amount"`
	RemainderNote  string `json:"remainder_note"`
	ChargesCreated string `json:"charges_created"`
	RevisionLabel  string `json:"revision_label"`
}

type ConfirmLabels struct {
	PublishTitle string `json:"publish_title"`
	PublishMsg   string `json:"publish_msg"`
	ReviseTitle  string `json:"revise_title"`
	ReviseMsg    string `json:"revise_msg"`
}

type EmptyLabels struct {
	DraftTitle    string `json:"draft_title"`
	SharesTitle   string `json:"shares_title"`
	SharesMessage string `json:"shares_message"`
}

type ErrorsLabels struct {
	SourceClaimedByRecognition string `json:"source_claimed_by_recognition"`
	SourceClaimedByAllocation  string `json:"source_claimed_by_allocation"`
	TermBoundaryCrossed        string `json:"term_boundary_crossed"`
	AlreadyPublished           string `json:"already_published"`
	NotDraft                   string `json:"not_draft"`
	NotFound                   string `json:"not_found"`
	SharesTotalMismatch        string `json:"shares_total_mismatch"`
	DenominatorInvalid         string `json:"denominator_invalid"`
	NoRecoverableShare         string `json:"no_recoverable_share"`
	Validation                 string `json:"validation"`
	ServicePeriodInvalid       string `json:"service_period_invalid"`
	PolicyComponentInvalid     string `json:"policy_component_invalid"`
	TransactionRequired        string `json:"transaction_required"`
	LockUnavailable            string `json:"lock_unavailable"`
	AgreementTermMissing       string `json:"agreement_term_missing"`
	Overlap                    string `json:"overlap"`
	Generic                    string `json:"generic"`
	Unavailable                string `json:"unavailable"`
	FormInvalid                string `json:"form_invalid"`
}

// DefaultLabels returns the English (general tier) defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:           "Cost allocation",
			Subtitle:        "Split a cost line across the agreements that share it",
			TitleDraft:      "Draft allocation",
			TitlePublished:  "Published allocation",
			TitleSuperseded: "Superseded allocations",
		},
		Tabs: TabsLabels{
			Shares:  "Shares",
			History: "Revisions",
			Charges: "Charges",
		},
		Columns: ColumnsLabels{
			Revision:     "Revision",
			Status:       "Status",
			ShareKind:    "Share",
			Subscription: "Agreement",
			Client:       "Customer",
			Numerator:    "Weight",
			Denominator:  "Total weight",
			Percentage:   "Share of total",
			Amount:       "Amount",
			Rule:         "Rule",
			PublishedOn:  "Published on",
			PublishedBy:  "Published by",
			PreparedBy:   "Prepared by",
			Total:        "Total",
		},
		Buttons: ButtonsLabels{
			Allocate:    "Allocate cost",
			AddShare:    "Add share",
			RemoveShare: "Remove",
			Preview:     "Preview",
			SaveDraft:   "Save draft",
			Publish:     "Publish",
			Revise:      "Revise allocation",
			ViewCharges: "View charges",
		},
		Form: FormLabels{
			ShareKindLabel:          "Share",
			SubscriptionLabel:       "Agreement",
			SubscriptionPlaceholder: "Select an agreement",
			NumeratorLabel:          "Weight",
			NumeratorPlaceholder:    "e.g. 250",
			DenominatorLabel:        "Total weight",
			RuleLabel:               "Split rule",
			PreviewTitle:            "Preview",
			OwnUseRowLabel:          "Own use",
			VacancyRowLabel:         "Vacancy",
			CommonLossRowLabel:      "Common-area loss",
		},
		Enums: EnumsLabels{
			StatusDraft:          "Draft",
			StatusPublished:      "Published",
			StatusSuperseded:     "Superseded",
			ShareKindRecoverable: "Recoverable",
			ShareKindOwnUse:      "Own use",
			ShareKindVacancy:     "Vacancy",
			ShareKindCommonLoss:  "Common-area loss",
			RuleLargestRemainder: "Proportional (largest remainder)",
		},
		Detail: DetailLabels{
			TotalAllocated: "Total allocated",
			SourceAmount:   "Cost line amount",
			RemainderNote:  "Rounding remainders go to the largest fractions so the shares add up exactly.",
			ChargesCreated: "{0} charges created",
			RevisionLabel:  "Revision {0}",
		},
		Confirm: ConfirmLabels{
			PublishTitle: "Publish this allocation?",
			PublishMsg:   "Charges will be created for each recoverable share. A published allocation can only be replaced by a new revision.",
			ReviseTitle:  "Start a new revision?",
			ReviseMsg:    "The published allocation stays on record until the new revision is published.",
		},
		Empty: EmptyLabels{
			DraftTitle:    "No draft allocation",
			SharesTitle:   "No shares yet",
			SharesMessage: "Add the agreements that share this cost.",
		},
		Errors: ErrorsLabels{
			SourceClaimedByRecognition: "This cost is already recognised as an expense.",
			SourceClaimedByAllocation:  "This cost is already allocated.",
			TermBoundaryCrossed:        "The service period crosses a change in the agreement's charge terms.",
			AlreadyPublished:           "Already published.",
			NotDraft:                   "Only drafts can be changed.",
			NotFound:                   "Allocation not found.",
			SharesTotalMismatch:        "The shares do not add up to the cost line amount.",
			DenominatorInvalid:         "Total weight must be above zero.",
			NoRecoverableShare:         "Add at least one recoverable share.",
			Validation:                 "The allocation is not valid.",
			ServicePeriodInvalid:       "The cost component has no valid service period.",
			PolicyComponentInvalid:     "The charge policy version must have exactly one complete component.",
			TransactionRequired:        "This action needs a database transaction, which is not available.",
			LockUnavailable:            "This action needs row locking, which the storage provider does not support.",
			AgreementTermMissing:       "The agreement has no charge term for this customer.",
			Overlap:                    "Charge terms of one line cannot overlap.",
			Generic:                    "Something went wrong. Please try again.",
			Unavailable:                "This action is not available right now.",
			FormInvalid:                "Please check the form and try again.",
		},
	}
}

// ErrorMessage maps a use-case refusal code (ErrorKind) to its label; unknown codes
// fall back to Errors.Generic.
func (l Labels) ErrorMessage(code string) string {
	switch code {
	case ErrSourceClaimedByRecognition:
		return l.Errors.SourceClaimedByRecognition
	case ErrSourceClaimedByAllocation:
		return l.Errors.SourceClaimedByAllocation
	case ErrTermBoundaryCrossed:
		return l.Errors.TermBoundaryCrossed
	case ErrAlreadyPublished:
		return l.Errors.AlreadyPublished
	case ErrNotDraft:
		return l.Errors.NotDraft
	case ErrNotFound:
		return l.Errors.NotFound
	case ErrSharesTotalMismatch:
		return l.Errors.SharesTotalMismatch
	case ErrDenominatorInvalid:
		return l.Errors.DenominatorInvalid
	case ErrNoRecoverableShare:
		return l.Errors.NoRecoverableShare
	case ErrValidation:
		return l.Errors.Validation
	case ErrServicePeriodInvalid:
		return l.Errors.ServicePeriodInvalid
	case ErrPolicyComponentInvalid:
		return l.Errors.PolicyComponentInvalid
	case ErrTransactionRequired:
		return l.Errors.TransactionRequired
	case ErrLockUnavailable:
		return l.Errors.LockUnavailable
	case ErrAgreementTermMissing:
		return l.Errors.AgreementTermMissing
	case ErrOverlap:
		return l.Errors.Overlap
	default:
		return l.Errors.Generic
	}
}
