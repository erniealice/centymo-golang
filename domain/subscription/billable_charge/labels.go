package billable_charge

// Lyngua file: general/billable_charge.json (root key billable_charge), leasing overlay: values only. JSON tags byte-match the key tree
// (TestLabelsMatchLynguaKeyTree); every field has a general-tier key (TestDefaultLabelsEqualGeneralLyngua).

// Labels holds every translatable string of the billable charge pages.
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
	Title          string `json:"title"`
	Subtitle       string `json:"subtitle"`
	TitleOpen      string `json:"title_open"`
	TitleIssued    string `json:"title_issued"`
	TitleCancelled string `json:"title_cancelled"`
}

type TabsLabels struct {
	Open       string `json:"open"`
	Issued     string `json:"issued"`
	Cancelled  string `json:"cancelled"`
	Info       string `json:"info"`
	Components string `json:"components"`
	History    string `json:"history"`
}

type ColumnsLabels struct {
	Select         string `json:"select"`
	Client         string `json:"client"`
	Subscription   string `json:"subscription"`
	ChargeKind     string `json:"charge_kind"`
	ServicePeriod  string `json:"service_period"`
	Amount         string `json:"amount"`
	Status         string `json:"status"`
	Document       string `json:"document"`
	Source         string `json:"source"`
	ObligationKey  string `json:"obligation_key"`
	AccountingDate string `json:"accounting_date"`
	Reason         string `json:"reason"`
	Component      string `json:"component"`
	Presentation   string `json:"presentation"`
	TaxPosition    string `json:"tax_position"`
}

type ButtonsLabels struct {
	Issue      string `json:"issue"`
	Adjust     string `json:"adjust"`
	ViewSource string `json:"view_source"`
	Cancel     string `json:"cancel"`
}

type FormLabels struct {
	SeriesLabel       string `json:"series_label"`
	SeriesPlaceholder string `json:"series_placeholder"`
	IssueDateLabel    string `json:"issue_date_label"`
	DueDateLabel      string `json:"due_date_label"`
	NewAmountLabel    string `json:"new_amount_label"`
	ReasonLabel       string `json:"reason_label"`
	ReasonPlaceholder string `json:"reason_placeholder"`
	IssueSummary      string `json:"issue_summary"`
}

type EnumsLabels struct {
	StatusOpen           string `json:"status_open"`
	StatusIssued         string `json:"status_issued"`
	StatusCancelled      string `json:"status_cancelled"`
	ChargeKindOriginal   string `json:"charge_kind_original"`
	ChargeKindCorrection string `json:"charge_kind_correction"`
}

type DetailLabels struct {
	OriginalAmount string `json:"original_amount"`
	Corrections    string `json:"corrections"`
	CurrentAmount  string `json:"current_amount"`
	Predecessor    string `json:"predecessor"`
	IssueResult    string `json:"issue_result"`
}

type ConfirmLabels struct {
	IssueTitle  string `json:"issue_title"`
	IssueMsg    string `json:"issue_msg"`
	AdjustTitle string `json:"adjust_title"`
	AdjustMsg   string `json:"adjust_msg"`
}

type EmptyLabels struct {
	OpenTitle      string `json:"open_title"`
	OpenMessage    string `json:"open_message"`
	IssuedTitle    string `json:"issued_title"`
	CancelledTitle string `json:"cancelled_title"`
}

type ErrorsLabels struct {
	ObligationConflict  string `json:"obligation_conflict"`
	AdjustNotDownward   string `json:"adjust_not_downward"`
	NotFound            string `json:"not_found"`
	NotIssued           string `json:"not_issued"`
	Validation          string `json:"validation"`
	TransactionRequired string `json:"transaction_required"`
	LockUnavailable     string `json:"lock_unavailable"`
}

// DefaultLabels returns the English fallback (equals the general Lyngua tier).
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:          "Billable charges",
			Subtitle:       "Charges ready to be put on a document",
			TitleOpen:      "Open charges",
			TitleIssued:    "Issued charges",
			TitleCancelled: "Cancelled charges",
		},
		Tabs: TabsLabels{
			Open:       "Open",
			Issued:     "Issued",
			Cancelled:  "Cancelled",
			Info:       "Info",
			Components: "Components",
			History:    "History",
		},
		Columns: ColumnsLabels{
			Select:         "Select",
			Client:         "Customer",
			Subscription:   "Agreement",
			ChargeKind:     "Kind",
			ServicePeriod:  "Service period",
			Amount:         "Amount",
			Status:         "Status",
			Document:       "Document",
			Source:         "Source",
			ObligationKey:  "Reference",
			AccountingDate: "Accounting date",
			Reason:         "Reason",
			Component:      "Component",
			Presentation:   "Presentation",
			TaxPosition:    "Tax position",
		},
		Buttons: ButtonsLabels{
			Issue:      "Issue recovery documents",
			Adjust:     "Adjust",
			ViewSource: "View source cost",
			Cancel:     "Cancel charge",
		},
		Form: FormLabels{
			SeriesLabel:       "Document series",
			SeriesPlaceholder: "Select a series",
			IssueDateLabel:    "Issue date",
			DueDateLabel:      "Due date",
			NewAmountLabel:    "New amount",
			ReasonLabel:       "Reason",
			ReasonPlaceholder: "Why is the charge being lowered?",
			IssueSummary:      "{0} charges for {1} customers",
		},
		Enums: EnumsLabels{
			StatusOpen:           "Open",
			StatusIssued:         "Issued",
			StatusCancelled:      "Cancelled",
			ChargeKindOriginal:   "Original",
			ChargeKindCorrection: "Correction",
		},
		Detail: DetailLabels{
			OriginalAmount: "Original amount",
			Corrections:    "Corrections",
			CurrentAmount:  "Current amount",
			Predecessor:    "Corrects {0}",
			IssueResult:    "{0} documents issued",
		},
		Confirm: ConfirmLabels{
			IssueTitle:  "Issue documents for the selected charges?",
			IssueMsg:    "One document is created for each customer and agreement. Issued charges cannot be edited.",
			AdjustTitle: "Adjust this charge?",
			AdjustMsg:   "The original stays unchanged. A credit document is issued for the reduction.",
		},
		Empty: EmptyLabels{
			OpenTitle:      "No open charges",
			OpenMessage:    "Publish a cost allocation to create charges.",
			IssuedTitle:    "No issued charges",
			CancelledTitle: "No cancelled charges",
		},
		Errors: ErrorsLabels{
			ObligationConflict:  "This charge already exists with different details.",
			AdjustNotDownward:   "The new amount must be lower than the current amount.",
			NotFound:            "Charge not found.",
			NotIssued:           "Only issued charges can be adjusted.",
			Validation:          "The charge is not valid.",
			TransactionRequired: "This action needs a database transaction, which is not available.",
			LockUnavailable:     "This action needs row locking, which the storage provider does not support.",
		},
	}
}
