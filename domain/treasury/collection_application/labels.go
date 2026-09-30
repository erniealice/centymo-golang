package collection_application

// Lyngua file: general/collection_application.json (root key collection_application), leasing overlay: values only. JSON tags byte-match the key tree
// (TestLabelsMatchLynguaKeyTree); every field has a general-tier key (TestDefaultLabelsEqualGeneralLyngua).

// Labels holds every translatable string of the receive-and-apply flow.
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
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type TabsLabels struct {
	Applications string `json:"applications"`
}

type ColumnsLabels struct {
	Target         string `json:"target"`
	TargetKind     string `json:"target_kind"`
	DocumentNumber string `json:"document_number"`
	DueDate        string `json:"due_date"`
	OpenBalance    string `json:"open_balance"`
	Applied        string `json:"applied"`
	Amount         string `json:"amount"`
	AppliedOn      string `json:"applied_on"`
	AppliedBy      string `json:"applied_by"`
	Status         string `json:"status"`
	Kind           string `json:"kind"`
	Certificate    string `json:"certificate"`
}

type ButtonsLabels struct {
	ReceiveApply string `json:"receive_apply"`
	Apply        string `json:"apply"`
	Reverse      string `json:"reverse"`
	Preview      string `json:"preview"`
}

type FormLabels struct {
	ClientLabel          string `json:"client_label"`
	ClientPlaceholder    string `json:"client_placeholder"`
	AmountLabel          string `json:"amount_label"`
	DateLabel            string `json:"date_label"`
	MethodLabel          string `json:"method_label"`
	MethodPlaceholder    string `json:"method_placeholder"`
	ReferenceLabel       string `json:"reference_label"`
	ReferencePlaceholder string `json:"reference_placeholder"`
	CurrencyLabel        string `json:"currency_label"`
	PreviewTitle         string `json:"preview_title"`
	ReverseReasonLabel   string `json:"reverse_reason_label"`
}

type EnumsLabels struct {
	StatusApplied                    string `json:"status_applied"`
	StatusReversed                   string `json:"status_reversed"`
	TargetKindRevenue                string `json:"target_kind_revenue"`
	TargetKindRecoveryDocument       string `json:"target_kind_recovery_document"`
	ApplicationKindCash              string `json:"application_kind_cash"`
	ApplicationKindNonCashSettlement string `json:"application_kind_non_cash_settlement"`
}

type DetailLabels struct {
	Unapplied    string `json:"unapplied"`
	OrderNote    string `json:"order_note"`
	TotalApplied string `json:"total_applied"`
	OpenItems    string `json:"open_items"`
	ReversalOf   string `json:"reversal_of"`
}

type ConfirmLabels struct {
	ApplyTitle   string `json:"apply_title"`
	ApplyMsg     string `json:"apply_msg"`
	ReverseTitle string `json:"reverse_title"`
	ReverseMsg   string `json:"reverse_msg"`
}

type EmptyLabels struct {
	OpenItemsTitle    string `json:"open_items_title"`
	OpenItemsMessage  string `json:"open_items_message"`
	ApplicationsTitle string `json:"applications_title"`
}

type ErrorsLabels struct {
	CurrencyMismatch    string `json:"currency_mismatch"`
	NotFound            string `json:"not_found"`
	AmountInvalid       string `json:"amount_invalid"`
	ClientRequired      string `json:"client_required"`
	AlreadyReversed     string `json:"already_reversed"`
	MissingPosting      string `json:"missing_posting"`
	TransactionRequired string `json:"transaction_required"`
	DateInvalid         string `json:"date_invalid"`
}

// DefaultLabels returns the English fallback (equals the general Lyngua tier).
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:    "Receive and apply",
			Subtitle: "Record a payment and apply it to what the customer owes",
		},
		Tabs: TabsLabels{
			Applications: "Applications",
		},
		Columns: ColumnsLabels{
			Target:         "Applied to",
			TargetKind:     "Type",
			DocumentNumber: "Number",
			DueDate:        "Due date",
			OpenBalance:    "Open balance",
			Applied:        "Applied",
			Amount:         "Amount",
			AppliedOn:      "Applied on",
			AppliedBy:      "Applied by",
			Status:         "Status",
			Kind:           "Application",
			Certificate:    "Certificate",
		},
		Buttons: ButtonsLabels{
			ReceiveApply: "Receive and apply",
			Apply:        "Apply",
			Reverse:      "Reverse",
			Preview:      "Preview",
		},
		Form: FormLabels{
			ClientLabel:          "Customer",
			ClientPlaceholder:    "Select a customer",
			AmountLabel:          "Amount received",
			DateLabel:            "Payment date",
			MethodLabel:          "Method",
			MethodPlaceholder:    "Select a method",
			ReferenceLabel:       "Reference",
			ReferencePlaceholder: "e.g. cheque or transfer number",
			CurrencyLabel:        "Currency",
			PreviewTitle:         "How this payment will be applied",
			ReverseReasonLabel:   "Reason",
		},
		Enums: EnumsLabels{
			StatusApplied:                    "Applied",
			StatusReversed:                   "Reversed",
			TargetKindRevenue:                "Invoice",
			TargetKindRecoveryDocument:       "Recovery document",
			ApplicationKindCash:              "Cash",
			ApplicationKindNonCashSettlement: "Withholding certificate",
		},
		Detail: DetailLabels{
			Unapplied:    "Unapplied credit",
			OrderNote:    "Applied to the oldest due date first. On the same date, invoices before recovery documents.",
			TotalApplied: "Total applied",
			OpenItems:    "Open items",
			ReversalOf:   "Reverses {0}",
		},
		Confirm: ConfirmLabels{
			ApplyTitle:   "Receive and apply this payment?",
			ApplyMsg:     "The payment is recorded and applied as shown in the preview.",
			ReverseTitle: "Reverse this application?",
			ReverseMsg:   "The amount returns to the customer's open balance. The original stays on record.",
		},
		Empty: EmptyLabels{
			OpenItemsTitle:    "Nothing is open for this customer",
			OpenItemsMessage:  "The full amount will be kept as unapplied credit.",
			ApplicationsTitle: "No applications yet",
		},
		Errors: ErrorsLabels{
			CurrencyMismatch:    "The currencies do not match.",
			NotFound:            "Application not found.",
			AmountInvalid:       "Enter an amount above zero.",
			ClientRequired:      "Select a customer.",
			AlreadyReversed:     "This application is already reversed.",
			MissingPosting:      "The charge policy has no account mapping for this step.",
			TransactionRequired: "This action needs a database transaction, which is not available.",
			DateInvalid:         "Enter the payment date as YYYY-MM-DD.",
		},
	}
}
