package disbursement

// labels.go — disbursement-entity label structs (centymo W5).
//
// Disbursement (money OUT) labels, extracted verbatim from the treasury-domain
// labels.go into the per-entity disbursement package per the domain-first
// restructure. Pure structural move — no behaviour change. Lyngua JSON load
// paths are unchanged.

// ---------------------------------------------------------------------------
// Disbursement labels (money OUT — payments, refunds, payouts)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the disbursement module.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Buttons ButtonLabels  `json:"buttons"`
	Columns ColumnLabels  `json:"columns"`
	Empty   EmptyLabels   `json:"empty"`
	Form    FormLabels    `json:"form"`
	Actions ActionLabels  `json:"actions"`
	Bulk    BulkLabels    `json:"bulkActions"`
	Detail  DetailLabels  `json:"detail"`
	Status  StatusLabels  `json:"status"`
	Confirm ConfirmLabels `json:"confirm"`
	Errors  ErrorLabels   `json:"errors"`
}

type PageLabels struct {
	Heading          string `json:"heading"`
	HeadingDraft     string `json:"headingDraft"`
	HeadingPending   string `json:"headingPending"`
	HeadingApproved  string `json:"headingApproved"`
	HeadingPaid      string `json:"headingPaid"`
	HeadingCancelled string `json:"headingCancelled"`
	Caption          string `json:"caption"`
	CaptionDraft     string `json:"captionDraft"`
	CaptionPending   string `json:"captionPending"`
	CaptionApproved  string `json:"captionApproved"`
	CaptionPaid      string `json:"captionPaid"`
	CaptionCancelled string `json:"captionCancelled"`
	Dashboard        string `json:"dashboard"`
}

type ButtonLabels struct {
	AddDisbursement string `json:"addDisbursement"`
}

type ColumnLabels struct {
	Reference string `json:"reference"`
	Payee     string `json:"payee"`
	Amount    string `json:"amount"`
	Date      string `json:"date"`
	Status    string `json:"status"`
	Method    string `json:"method"`
	Category  string `json:"category"`
}

type EmptyLabels struct {
	DraftTitle       string `json:"draftTitle"`
	DraftMessage     string `json:"draftMessage"`
	PendingTitle     string `json:"pendingTitle"`
	PendingMessage   string `json:"pendingMessage"`
	ApprovedTitle    string `json:"approvedTitle"`
	ApprovedMessage  string `json:"approvedMessage"`
	PaidTitle        string `json:"paidTitle"`
	PaidMessage      string `json:"paidMessage"`
	CancelledTitle   string `json:"cancelledTitle"`
	CancelledMessage string `json:"cancelledMessage"`
}

type FormLabels struct {
	Payee                   string `json:"payee"`
	PayeePlaceholder        string `json:"payeePlaceholder"`
	Date                    string `json:"date"`
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	Reference               string `json:"reference"`
	ReferencePlaceholder    string `json:"referencePlaceholder"`
	PaymentMethod           string `json:"paymentMethod"`
	Category                string `json:"category"`
	Status                  string `json:"status"`
	Notes                   string `json:"notes"`
	NotesPlaceholder        string `json:"notesPlaceholder"`
	ApprovedBy              string `json:"approvedBy"`
	AmountPlaceholder       string `json:"amountPlaceholder"`
	CurrencyPlaceholder     string `json:"currencyPlaceholder"`
	MethodCash              string `json:"methodCash"`
	MethodBankTransfer      string `json:"methodBankTransfer"`
	MethodCheck             string `json:"methodCheck"`
	MethodGCash             string `json:"methodGCash"`
	MethodOther             string `json:"methodOther"`
	StatusDraft             string `json:"statusDraft"`
	StatusPending           string `json:"statusPending"`
	StatusApproved          string `json:"statusApproved"`
	StatusPaid              string `json:"statusPaid"`
	StatusCancelled         string `json:"statusCancelled"`
	TypeSupplierPayment     string `json:"typeSupplierPayment"`
	TypePayroll             string `json:"typePayroll"`
	TypeRent                string `json:"typeRent"`
	TypeUtilities           string `json:"typeUtilities"`
	TypeOther               string `json:"typeOther"`
	ApproverNamePlaceholder string `json:"approverNamePlaceholder"`
	LinkToBill              string `json:"linkToBill"`
	NoBillOption            string `json:"noBillOption"`

	// Field-level info text surfaced via an info button beside each label.
	ReferenceInfo     string `json:"referenceInfo"`
	DateInfo          string `json:"dateInfo"`
	PayeeInfo         string `json:"payeeInfo"`
	AmountInfo        string `json:"amountInfo"`
	CurrencyInfo      string `json:"currencyInfo"`
	PaymentMethodInfo string `json:"paymentMethodInfo"`
	StatusInfo        string `json:"statusInfo"`
	CategoryInfo      string `json:"categoryInfo"`
	ApprovedByInfo    string `json:"approvedByInfo"`
	NotesInfo         string `json:"notesInfo"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Approve    string `json:"approve"`
	MarkPaid   string `json:"markPaid"`
	Cancel     string `json:"cancel"`
	Submit     string `json:"submit"`
	Reactivate string `json:"reactivate"`
}

type BulkLabels struct {
	Delete   string `json:"delete"`
	Approve  string `json:"approve"`
	MarkPaid string `json:"markPaid"`
}

type DetailLabels struct {
	PageTitle         string `json:"pageTitle"`
	TitlePrefix       string `json:"titlePrefix"`
	PaymentInfo       string `json:"paymentInfo"`
	Payee             string `json:"payee"`
	Date              string `json:"date"`
	Amount            string `json:"amount"`
	Currency          string `json:"currency"`
	Status            string `json:"status"`
	Method            string `json:"method"`
	Category          string `json:"category"`
	Reference         string `json:"reference"`
	ApprovedBy        string `json:"approvedBy"`
	Notes             string `json:"notes"`
	TabBasicInfo      string `json:"tabBasicInfo"`
	TabAttachments    string `json:"tabAttachments"`
	TabAuditTrail     string `json:"tabAuditTrail"`
	TabAuditHistory   string `json:"tabAuditHistory"`
	AuditAction       string `json:"auditAction"`
	AuditUser         string `json:"auditUser"`
	AuditEmptyTitle   string `json:"auditEmptyTitle"`
	AuditEmptyMessage string `json:"auditEmptyMessage"`
}

type StatusLabels struct {
	Draft     string `json:"draft"`
	Pending   string `json:"pending"`
	Approved  string `json:"approved"`
	Paid      string `json:"paid"`
	Cancelled string `json:"cancelled"`
}

type ConfirmLabels struct {
	Submit                string `json:"submit"`
	SubmitMessage         string `json:"submitMessage"`
	Approve               string `json:"approve"`
	ApproveMessage        string `json:"approveMessage"`
	MarkPaid              string `json:"markPaid"`
	MarkPaidMessage       string `json:"markPaidMessage"`
	Cancel                string `json:"cancel"`
	CancelMessage         string `json:"cancelMessage"`
	Reactivate            string `json:"reactivate"`
	ReactivateMessage     string `json:"reactivateMessage"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"deleteMessage"`
	BulkSubmit            string `json:"bulkSubmit"`
	BulkSubmitMessage     string `json:"bulkSubmitMessage"`
	BulkApprove           string `json:"bulkApprove"`
	BulkApproveMessage    string `json:"bulkApproveMessage"`
	BulkMarkPaid          string `json:"bulkMarkPaid"`
	BulkMarkPaidMessage   string `json:"bulkMarkPaidMessage"`
	BulkCancel            string `json:"bulkCancel"`
	BulkCancelMessage     string `json:"bulkCancelMessage"`
	BulkReactivate        string `json:"bulkReactivate"`
	BulkReactivateMessage string `json:"bulkReactivateMessage"`
	BulkDelete            string `json:"bulkDelete"`
	BulkDeleteMessage     string `json:"bulkDeleteMessage"`
}

type ErrorLabels struct {
	PermissionDenied  string `json:"permissionDenied"`
	InvalidFormData   string `json:"invalidFormData"`
	NotFound          string `json:"notFound"`
	IDRequired        string `json:"idRequired"`
	NoIDsProvided     string `json:"noIDsProvided"`
	InvalidStatus     string `json:"invalidStatus"`
	InvalidTransition string `json:"invalidTransition"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:          "Disbursements",
			HeadingDraft:     "Draft Disbursements",
			HeadingPending:   "Pending Disbursements",
			HeadingApproved:  "Approved Disbursements",
			HeadingPaid:      "Paid Disbursements",
			HeadingCancelled: "Cancelled Disbursements",
			Caption:          "Manage disbursements and payouts",
			CaptionDraft:     "Draft disbursements awaiting submission",
			CaptionPending:   "Disbursements awaiting approval",
			CaptionApproved:  "Approved disbursements ready for payment",
			CaptionPaid:      "Completed disbursement payments",
			CaptionCancelled: "Cancelled disbursements",
			Dashboard:        "Disbursements Dashboard",
		},
		Buttons: ButtonLabels{
			AddDisbursement: "Add Disbursement",
		},
		Columns: ColumnLabels{
			Reference: "Reference",
			Payee:     "Payee",
			Amount:    "Amount",
			Date:      "Date",
			Status:    "Status",
			Method:    "Method",
			Category:  "Category",
		},
		Empty: EmptyLabels{
			DraftTitle:       "No draft disbursements",
			DraftMessage:     "No draft disbursements to display.",
			PendingTitle:     "No pending disbursements",
			PendingMessage:   "No pending disbursements to display.",
			ApprovedTitle:    "No approved disbursements",
			ApprovedMessage:  "No approved disbursements to display.",
			PaidTitle:        "No paid disbursements",
			PaidMessage:      "No paid disbursements to display.",
			CancelledTitle:   "No cancelled disbursements",
			CancelledMessage: "No cancelled disbursements to display.",
		},
		Form: FormLabels{
			Payee:                   "Payee",
			PayeePlaceholder:        "Enter payee name",
			Date:                    "Date",
			Amount:                  "Amount",
			Currency:                "Currency",
			Reference:               "Reference",
			ReferencePlaceholder:    "e.g. DISB-001",
			PaymentMethod:           "Payment Method",
			Category:                "Category",
			Status:                  "Status",
			Notes:                   "Notes",
			NotesPlaceholder:        "Additional notes...",
			ApprovedBy:              "Approved By",
			AmountPlaceholder:       "0.00",
			CurrencyPlaceholder:     "PHP",
			MethodCash:              "Cash",
			MethodBankTransfer:      "Bank Transfer",
			MethodCheck:             "Check",
			MethodGCash:             "GCash",
			MethodOther:             "Other",
			StatusDraft:             "Draft",
			StatusPending:           "Pending",
			StatusApproved:          "Approved",
			StatusPaid:              "Paid",
			StatusCancelled:         "Cancelled",
			TypeSupplierPayment:     "Supplier Payment",
			TypePayroll:             "Payroll",
			TypeRent:                "Rent",
			TypeUtilities:           "Utilities",
			TypeOther:               "Other",
			ApproverNamePlaceholder: "Approver name",
			LinkToBill:              "Link to Bill",
			NoBillOption:            "— No Bill —",
			// Field-level info popovers — use proto-generic wording; tiers override via lyngua.
			ReferenceInfo:     "Unique reference number for this disbursement.",
			DateInfo:          "Date the disbursement was issued.",
			PayeeInfo:         "Name of the recipient (supplier, payroll, etc.).",
			AmountInfo:        "Total amount disbursed (in centavos; displayed as amount ÷ 100).",
			CurrencyInfo:      "Currency of the disbursed amount.",
			PaymentMethodInfo: "How the payment was made.",
			StatusInfo:        "Current state of this disbursement.",
			CategoryInfo:      "Type of disbursement for categorisation and reporting.",
			ApprovedByInfo:    "Name of the person who authorised this disbursement.",
			NotesInfo:         "Internal remarks — not shown on supplier-facing documents.",
		},
		Actions: ActionLabels{
			View:       "View",
			Edit:       "Edit",
			Delete:     "Delete",
			Approve:    "Approve",
			MarkPaid:   "Mark as Paid",
			Cancel:     "Cancel",
			Submit:     "Submit",
			Reactivate: "Reactivate",
		},
		Bulk: BulkLabels{
			Delete:   "Delete Selected",
			Approve:  "Approve Selected",
			MarkPaid: "Mark Selected as Paid",
		},
		Detail: DetailLabels{
			PageTitle:         "Disbursement Details",
			TitlePrefix:       "Disbursement #",
			PaymentInfo:       "Payment Information",
			Payee:             "Payee",
			Date:              "Date",
			Amount:            "Amount",
			Currency:          "Currency",
			Status:            "Status",
			Method:            "Payment Method",
			Category:          "Category",
			Reference:         "Reference",
			ApprovedBy:        "Approved By",
			Notes:             "Notes",
			TabBasicInfo:      "Basic Info",
			TabAttachments:    "Attachments",
			TabAuditTrail:     "Audit Trail",
			TabAuditHistory:   "History",
			AuditAction:       "Action",
			AuditUser:         "User",
			AuditEmptyTitle:   "No audit records",
			AuditEmptyMessage: "No audit trail entries yet.",
		},
		Status: StatusLabels{
			Draft:     "Draft",
			Pending:   "Pending",
			Approved:  "Approved",
			Paid:      "Paid",
			Cancelled: "Cancelled",
		},
		Confirm: ConfirmLabels{
			Submit:                "Submit",
			SubmitMessage:         "Are you sure you want to submit {{count}} disbursement(s)?",
			Approve:               "Approve",
			ApproveMessage:        "Are you sure you want to approve {{count}} disbursement(s)?",
			MarkPaid:              "Mark as Paid",
			MarkPaidMessage:       "Are you sure you want to mark {{count}} disbursement(s) as paid?",
			Cancel:                "Cancel",
			CancelMessage:         "Are you sure you want to cancel {{count}} disbursement(s)?",
			Reactivate:            "Reactivate",
			ReactivateMessage:     "Are you sure you want to reactivate {{count}} disbursement(s)?",
			Delete:                "Delete",
			DeleteMessage:         "Are you sure you want to delete {{count}} disbursement(s)?",
			BulkSubmit:            "Submit Disbursements",
			BulkSubmitMessage:     "Are you sure you want to submit {{count}} disbursement(s)?",
			BulkApprove:           "Approve Disbursements",
			BulkApproveMessage:    "Are you sure you want to approve {{count}} disbursement(s)?",
			BulkMarkPaid:          "Mark as Paid",
			BulkMarkPaidMessage:   "Are you sure you want to mark {{count}} disbursement(s) as paid?",
			BulkCancel:            "Cancel Disbursements",
			BulkCancelMessage:     "Are you sure you want to cancel {{count}} disbursement(s)?",
			BulkReactivate:        "Reactivate Disbursements",
			BulkReactivateMessage: "Are you sure you want to reactivate {{count}} disbursement(s)?",
			BulkDelete:            "Delete Disbursements",
			BulkDeleteMessage:     "Are you sure you want to delete {{count}} disbursement(s)?",
		},
		Errors: ErrorLabels{
			PermissionDenied:  "Permission denied",
			InvalidFormData:   "Invalid form data",
			NotFound:          "Disbursement not found",
			IDRequired:        "Disbursement ID is required",
			NoIDsProvided:     "No disbursement IDs provided",
			InvalidStatus:     "Invalid target status",
			InvalidTransition: "Cannot transition from %s to %s",
		},
	}
}
