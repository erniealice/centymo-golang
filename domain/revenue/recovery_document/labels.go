package recovery_document

// Lyngua file: general/recovery_document.json (root key recovery_document), leasing overlay: values only. JSON tags byte-match the key tree
// (TestLabelsMatchLynguaKeyTree); every field has a general-tier key (TestDefaultLabelsEqualGeneralLyngua).

// Labels holds every translatable string of the recovery document pages.
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
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	TitleIssued string `json:"title_issued"`
	TitleVoid   string `json:"title_void"`
}

type TabsLabels struct {
	Issued       string `json:"issued"`
	Void         string `json:"void"`
	Info         string `json:"info"`
	Lines        string `json:"lines"`
	Applications string `json:"applications"`
	CreditNotes  string `json:"credit_notes"`
}

type ColumnsLabels struct {
	DocumentNumber string `json:"document_number"`
	DocumentType   string `json:"document_type"`
	Client         string `json:"client"`
	Subscription   string `json:"subscription"`
	IssueDate      string `json:"issue_date"`
	DueDate        string `json:"due_date"`
	Total          string `json:"total"`
	Paid           string `json:"paid"`
	Balance        string `json:"balance"`
	Status         string `json:"status"`
	Description    string `json:"description"`
	ServicePeriod  string `json:"service_period"`
	Amount         string `json:"amount"`
	Corrects       string `json:"corrects"`
	IssuedBy       string `json:"issued_by"`
	VoidReason     string `json:"void_reason"`
}

type ButtonsLabels struct {
	Void           string `json:"void"`
	Print          string `json:"print"`
	ViewCharge     string `json:"view_charge"`
	ApplyPayment   string `json:"apply_payment"`
	ViewCreditNote string `json:"view_credit_note"`
}

type FormLabels struct {
	VoidReasonLabel       string `json:"void_reason_label"`
	VoidReasonPlaceholder string `json:"void_reason_placeholder"`
}

type EnumsLabels struct {
	StatusIssued   string `json:"status_issued"`
	StatusVoid     string `json:"status_void"`
	TypeStatement  string `json:"type_statement"`
	TypeCreditNote string `json:"type_credit_note"`
}

type DetailLabels struct {
	IssuedOn         string `json:"issued_on"`
	DueOn            string `json:"due_on"`
	CorrectsDocument string `json:"corrects_document"`
	VoidNotice       string `json:"void_notice"`
	BalanceDue       string `json:"balance_due"`
	CreditBalance    string `json:"credit_balance"`
}

type ConfirmLabels struct {
	VoidTitle string `json:"void_title"`
	VoidMsg   string `json:"void_msg"`
}

type EmptyLabels struct {
	IssuedTitle       string `json:"issued_title"`
	IssuedMessage     string `json:"issued_message"`
	VoidTitle         string `json:"void_title"`
	ApplicationsTitle string `json:"applications_title"`
	CreditNotesTitle  string `json:"credit_notes_title"`
}

type ErrorsLabels struct {
	HasApplications       string `json:"has_applications"`
	MissingPosting        string `json:"missing_posting"`
	SeriesRetired         string `json:"series_retired"`
	NotFound              string `json:"not_found"`
	AlreadyVoid           string `json:"already_void"`
	VoidReasonRequired    string `json:"void_reason_required"`
	NothingSelected       string `json:"nothing_selected"`
	NotOpen               string `json:"not_open"`
	NotIssued             string `json:"not_issued"`
	CurrencyMismatch      string `json:"currency_mismatch"`
	SeriesKindMismatch    string `json:"series_kind_mismatch"`
	ComponentKindMismatch string `json:"component_kind_mismatch"`
	IssuanceConflict      string `json:"issuance_conflict"`
	HasCreditNotes        string `json:"has_credit_notes"`
	Validation            string `json:"validation"`
	TransactionRequired   string `json:"transaction_required"`
}

// DefaultLabels returns the English fallback (equals the general Lyngua tier).
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:       "Recovery documents",
			Subtitle:    "Documents issued to recover costs from customers",
			TitleIssued: "Issued documents",
			TitleVoid:   "Void documents",
		},
		Tabs: TabsLabels{
			Issued:       "Issued",
			Void:         "Void",
			Info:         "Info",
			Lines:        "Lines",
			Applications: "Payments applied",
			CreditNotes:  "Credit notes",
		},
		Columns: ColumnsLabels{
			DocumentNumber: "Number",
			DocumentType:   "Type",
			Client:         "Customer",
			Subscription:   "Agreement",
			IssueDate:      "Issue date",
			DueDate:        "Due date",
			Total:          "Total",
			Paid:           "Applied",
			Balance:        "Balance",
			Status:         "Status",
			Description:    "Description",
			ServicePeriod:  "Service period",
			Amount:         "Amount",
			Corrects:       "Corrects",
			IssuedBy:       "Issued by",
			VoidReason:     "Void reason",
		},
		Buttons: ButtonsLabels{
			Void:           "Void",
			Print:          "Print",
			ViewCharge:     "View charge",
			ApplyPayment:   "Receive and apply",
			ViewCreditNote: "View credit note",
		},
		Form: FormLabels{
			VoidReasonLabel:       "Reason",
			VoidReasonPlaceholder: "Why is this document void?",
		},
		Enums: EnumsLabels{
			StatusIssued:   "Issued",
			StatusVoid:     "Void",
			TypeStatement:  "Statement",
			TypeCreditNote: "Credit note",
		},
		Detail: DetailLabels{
			IssuedOn:         "Issued on {0}",
			DueOn:            "Due on {0}",
			CorrectsDocument: "Corrects {0}",
			VoidNotice:       "This document is void. It stays on record.",
			BalanceDue:       "Balance due",
			CreditBalance:    "Credit balance",
		},
		Confirm: ConfirmLabels{
			VoidTitle: "Void this document?",
			VoidMsg:   "The document stays on record and its charges are cancelled. This cannot be undone.",
		},
		Empty: EmptyLabels{
			IssuedTitle:       "No issued documents",
			IssuedMessage:     "Issue documents from the billable charges list.",
			VoidTitle:         "No void documents",
			ApplicationsTitle: "No payments applied yet",
			CreditNotesTitle:  "No credit notes",
		},
		Errors: ErrorsLabels{
			HasApplications:       "Payments have been applied. Reverse them first.",
			MissingPosting:        "The charge policy has no account mapping for this step.",
			SeriesRetired:         "This document series is retired.",
			NotFound:              "Document not found.",
			AlreadyVoid:           "This document is already void.",
			VoidReasonRequired:    "Enter a reason.",
			NothingSelected:       "Select at least one charge.",
			NotOpen:               "Only open charges can be issued.",
			NotIssued:             "The document being corrected is not issued.",
			CurrencyMismatch:      "The currencies do not match.",
			SeriesKindMismatch:    "This document series does not issue this kind of document.",
			ComponentKindMismatch: "A charge on this document cannot be issued on this kind of document.",
			IssuanceConflict:      "These charges were already issued with different content.",
			HasCreditNotes:        "Credit notes correct this document. Void them first.",
			Validation:            "The document is not valid.",
			TransactionRequired:   "This action needs a database transaction, which is not available.",
		},
	}
}
