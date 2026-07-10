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
	Bulk    BulkLabels    `json:"bulk_actions"`
	Detail  DetailLabels  `json:"detail"`
	Status  StatusLabels  `json:"status"`
	Confirm ConfirmLabels `json:"confirm"`
	Errors  ErrorLabels   `json:"errors"`
}

type PageLabels struct {
	Heading          string `json:"heading"`
	HeadingDraft     string `json:"heading_draft"`
	HeadingPending   string `json:"heading_pending"`
	HeadingApproved  string `json:"heading_approved"`
	HeadingPaid      string `json:"heading_paid"`
	HeadingCancelled string `json:"heading_cancelled"`
	Caption          string `json:"caption"`
	CaptionDraft     string `json:"caption_draft"`
	CaptionPending   string `json:"caption_pending"`
	CaptionApproved  string `json:"caption_approved"`
	CaptionPaid      string `json:"caption_paid"`
	CaptionCancelled string `json:"caption_cancelled"`
	Dashboard        string `json:"dashboard"`
}

type ButtonLabels struct {
	AddDisbursement string `json:"add_disbursement"`
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
	DraftTitle       string `json:"draft_title"`
	DraftMessage     string `json:"draft_message"`
	PendingTitle     string `json:"pending_title"`
	PendingMessage   string `json:"pending_message"`
	ApprovedTitle    string `json:"approved_title"`
	ApprovedMessage  string `json:"approved_message"`
	PaidTitle        string `json:"paid_title"`
	PaidMessage      string `json:"paid_message"`
	CancelledTitle   string `json:"cancelled_title"`
	CancelledMessage string `json:"cancelled_message"`
}

type FormLabels struct {
	Payee                   string `json:"payee"`
	PayeePlaceholder        string `json:"payee_placeholder"`
	Date                    string `json:"date"`
	Amount                  string `json:"amount"`
	Currency                string `json:"currency"`
	Reference               string `json:"reference"`
	ReferencePlaceholder    string `json:"reference_placeholder"`
	PaymentMethod           string `json:"payment_method"`
	Category                string `json:"category"`
	Status                  string `json:"status"`
	Notes                   string `json:"notes"`
	NotesPlaceholder        string `json:"notes_placeholder"`
	ApprovedBy              string `json:"approved_by"`
	AmountPlaceholder       string `json:"amount_placeholder"`
	CurrencyPlaceholder     string `json:"currency_placeholder"`
	MethodCash              string `json:"method_cash"`
	MethodBankTransfer      string `json:"method_bank_transfer"`
	MethodCheck             string `json:"method_check"`
	MethodGCash             string `json:"method_gcash"`
	MethodOther             string `json:"method_other"`
	StatusDraft             string `json:"status_draft"`
	StatusPending           string `json:"status_pending"`
	StatusApproved          string `json:"status_approved"`
	StatusPaid              string `json:"status_paid"`
	StatusCancelled         string `json:"status_cancelled"`
	TypeSupplierPayment     string `json:"type_supplier_payment"`
	TypePayroll             string `json:"type_payroll"`
	TypeRent                string `json:"type_rent"`
	TypeUtilities           string `json:"type_utilities"`
	TypeOther               string `json:"type_other"`
	ApproverNamePlaceholder string `json:"approver_name_placeholder"`
	LinkToBill              string `json:"link_to_bill"`
	NoBillOption            string `json:"no_bill_option"`

	// Field-level info text surfaced via an info button beside each label.
	ReferenceInfo     string `json:"reference_info"`
	DateInfo          string `json:"date_info"`
	PayeeInfo         string `json:"payee_info"`
	AmountInfo        string `json:"amount_info"`
	CurrencyInfo      string `json:"currency_info"`
	PaymentMethodInfo string `json:"payment_method_info"`
	StatusInfo        string `json:"status_info"`
	CategoryInfo      string `json:"category_info"`
	ApprovedByInfo    string `json:"approved_by_info"`
	NotesInfo         string `json:"notes_info"`
}

type ActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Approve    string `json:"approve"`
	MarkPaid   string `json:"mark_paid"`
	Cancel     string `json:"cancel"`
	Submit     string `json:"submit"`
	Reactivate string `json:"reactivate"`
}

type BulkLabels struct {
	Delete   string `json:"delete"`
	Approve  string `json:"approve"`
	MarkPaid string `json:"mark_paid"`
}

type DetailLabels struct {
	PageTitle         string `json:"page_title"`
	TitlePrefix       string `json:"title_prefix"`
	PaymentInfo       string `json:"payment_info"`
	Payee             string `json:"payee"`
	Date              string `json:"date"`
	Amount            string `json:"amount"`
	Currency          string `json:"currency"`
	Status            string `json:"status"`
	Method            string `json:"method"`
	Category          string `json:"category"`
	Reference         string `json:"reference"`
	ApprovedBy        string `json:"approved_by"`
	Notes             string `json:"notes"`
	TabBasicInfo      string `json:"tab_basic_info"`
	TabAttachments    string `json:"tab_attachments"`
	TabAuditTrail     string `json:"tab_audit_trail"`
	TabAuditHistory   string `json:"tab_audit_history"`
	AuditAction       string `json:"audit_action"`
	AuditUser         string `json:"audit_user"`
	AuditEmptyTitle   string `json:"audit_empty_title"`
	AuditEmptyMessage string `json:"audit_empty_message"`
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
	SubmitMessage         string `json:"submit_message"`
	Approve               string `json:"approve"`
	ApproveMessage        string `json:"approve_message"`
	MarkPaid              string `json:"mark_paid"`
	MarkPaidMessage       string `json:"mark_paid_message"`
	Cancel                string `json:"cancel"`
	CancelMessage         string `json:"cancel_message"`
	Reactivate            string `json:"reactivate"`
	ReactivateMessage     string `json:"reactivate_message"`
	Delete                string `json:"delete"`
	DeleteMessage         string `json:"delete_message"`
	BulkSubmit            string `json:"bulk_submit"`
	BulkSubmitMessage     string `json:"bulk_submit_message"`
	BulkApprove           string `json:"bulk_approve"`
	BulkApproveMessage    string `json:"bulk_approve_message"`
	BulkMarkPaid          string `json:"bulk_mark_paid"`
	BulkMarkPaidMessage   string `json:"bulk_mark_paid_message"`
	BulkCancel            string `json:"bulk_cancel"`
	BulkCancelMessage     string `json:"bulk_cancel_message"`
	BulkReactivate        string `json:"bulk_reactivate"`
	BulkReactivateMessage string `json:"bulk_reactivate_message"`
	BulkDelete            string `json:"bulk_delete"`
	BulkDeleteMessage     string `json:"bulk_delete_message"`
}

type ErrorLabels struct {
	PermissionDenied  string `json:"permission_denied"`
	InvalidFormData   string `json:"invalid_form_data"`
	NotFound          string `json:"not_found"`
	IDRequired        string `json:"id_required"`
	NoIDsProvided     string `json:"no_ids_provided"`
	InvalidStatus     string `json:"invalid_status"`
	InvalidTransition string `json:"invalid_transition"`
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
