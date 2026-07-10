package accrued_expense

// ---------------------------------------------------------------------------
// AccruedExpense labels  (SPS P10)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the accrued_expense
// + accrued_expense_settlement modules. Loaded from lyngua key root
// "accruedExpense" with settlement subkeys merged in via composition.
type Labels struct {
	Page        PageLabels       `json:"page"`
	Buttons     ButtonLabels     `json:"buttons"`
	Columns     ColumnLabels     `json:"columns"`
	Tabs        TabLabels        `json:"tabs"`
	Detail      DetailLabels     `json:"detail"`
	Settlements SettlementLabels `json:"settlements"`
	Form        FormLabels       `json:"form"`
	Status      StatusLabels     `json:"status"`
	Actions     ActionLabels     `json:"actions"`
	Confirm     ConfirmLabels    `json:"confirm"`
	Balances    BalanceLabels    `json:"balances"`
	Empty       EmptyLabels      `json:"empty"`
	Errors      ErrorLabels      `json:"errors"`
}

type PageLabels struct {
	Heading            string `json:"heading"`
	Caption            string `json:"caption"`
	HeadingOutstanding string `json:"heading_outstanding"`
	HeadingPartial     string `json:"heading_partial"`
	HeadingSettled     string `json:"heading_settled"`
	HeadingReversed    string `json:"heading_reversed"`
	Dashboard          string `json:"dashboard"`
}

type ButtonLabels struct {
	Add                string `json:"add"`
	AccrueFromContract string `json:"accrue_from_contract"`
	Settle             string `json:"settle"`
	Reverse            string `json:"reverse"`
	AddSettlement      string `json:"add_settlement"`
}

type ColumnLabels struct {
	InternalID       string `json:"internal_id"`
	Name             string `json:"name"`
	Supplier         string `json:"supplier"`
	SupplierContract string `json:"supplier_contract"`
	RecognitionDate  string `json:"recognition_date"`
	PeriodStart      string `json:"period_start"`
	PeriodEnd        string `json:"period_end"`
	CycleDate        string `json:"cycle_date"`
	Currency         string `json:"currency"`
	AccruedAmount    string `json:"accrued_amount"`
	SettledAmount    string `json:"settled_amount"`
	RemainingAmount  string `json:"remaining_amount"`
	Status           string `json:"status"`
}

type TabLabels struct {
	Info        string `json:"info"`
	Settlements string `json:"settlements"`
	Source      string `json:"source"`
	Activity    string `json:"activity"`
}

type DetailLabels struct {
	PageTitle            string `json:"page_title"`
	Title                string `json:"title"`
	InfoSection          string `json:"info_section"`
	SettlementsSection   string `json:"settlements_section"`
	SourceSection        string `json:"source_section"`
	AuditTrailComingSoon string `json:"audit_trail_coming_soon"`
	AuditEmptyTitle      string `json:"audit_empty_title"`
	AuditEmptyMessage    string `json:"audit_empty_message"`
	TabAttachments       string `json:"tab_attachments"`

	// Info-tab + source-tab field labels (4.4)
	Notes          string `json:"notes"`
	SourceContract string `json:"source_contract"`
	Supplier       string `json:"supplier"`
	ExpenseAccount string `json:"expense_account"`
	AccrualAccount string `json:"accrual_account"`
}

type SettlementLabels struct {
	Expenditure        string `json:"expenditure"`
	AmountSettled      string `json:"amount_settled"`
	Currency           string `json:"currency"`
	FxRate             string `json:"fx_rate"`
	FxAdjustmentAmount string `json:"fx_adjustment_amount"`
	SettledAt          string `json:"settled_at"`
	Reversal           string `json:"reversal"`
	EmptyTitle         string `json:"empty_title"`
	EmptyMessage       string `json:"empty_message"`
	AddSettlement      string `json:"add_settlement"`

	// Drawer form labels
	FormExpenditure            string `json:"form_expenditure"`
	FormExpenditurePlaceholder string `json:"form_expenditure_placeholder"`
	FormAmountSettled          string `json:"form_amount_settled"`
	FormCurrency               string `json:"form_currency"`
	FormFxRate                 string `json:"form_fx_rate"`
	FormFxRateInfo             string `json:"form_fx_rate_info"`
	FormReversalReason         string `json:"form_reversal_reason"`
}

type FormLabels struct {
	// Section headers
	SectionIdentity   string `json:"section_identity"`
	SectionSource     string `json:"section_source"`
	SectionPeriod     string `json:"section_period"`
	SectionMoney      string `json:"section_money"`
	SectionAccounting string `json:"section_accounting"`
	SectionLifecycle  string `json:"section_lifecycle"`
	SectionNotes      string `json:"section_notes"`

	// §1 Identity
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"name_placeholder"`
	NameInfo               string `json:"name_info"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	InternalID             string `json:"internal_id"`
	InternalIDPlaceholder  string `json:"internal_id_placeholder"`

	// §2 Source
	SupplierContract       string `json:"supplier_contract"`
	SelectSupplierContract string `json:"select_supplier_contract"`
	SupplierContractInfo   string `json:"supplier_contract_info"`
	Supplier               string `json:"supplier"`
	SelectSupplier         string `json:"select_supplier"`
	SupplierInfo           string `json:"supplier_info"`

	// §3 Period
	RecognitionDate      string `json:"recognition_date"`
	RecognitionDateInfo  string `json:"recognition_date_info"`
	PeriodStart          string `json:"period_start"`
	PeriodStartInfo      string `json:"period_start_info"`
	PeriodEnd            string `json:"period_end"`
	PeriodEndInfo        string `json:"period_end_info"`
	CycleDate            string `json:"cycle_date"`
	CycleDatePlaceholder string `json:"cycle_date_placeholder"`
	CycleDateInfo        string `json:"cycle_date_info"`

	// §4 Money
	Currency                 string `json:"currency"`
	CurrencyPlaceholder      string `json:"currency_placeholder"`
	CurrencyInfo             string `json:"currency_info"`
	AccruedAmount            string `json:"accrued_amount"`
	AccruedAmountPlaceholder string `json:"accrued_amount_placeholder"`
	AccruedAmountInfo        string `json:"accrued_amount_info"`
	SettledAmount            string `json:"settled_amount"`
	SettledAmountInfo        string `json:"settled_amount_info"`
	RemainingAmount          string `json:"remaining_amount"`
	RemainingAmountInfo      string `json:"remaining_amount_info"`

	// §5 Lifecycle
	Status            string `json:"status"`
	SelectStatus      string `json:"select_status"`
	StatusInfo        string `json:"status_info"`
	StatusOutstanding string `json:"status_outstanding"`
	StatusPartial     string `json:"status_partial"`
	StatusSettled     string `json:"status_settled"`
	StatusReversed    string `json:"status_reversed"`

	// §6 Accounting
	ExpenseAccount       string `json:"expense_account"`
	SelectExpenseAccount string `json:"select_expense_account"`
	ExpenseAccountInfo   string `json:"expense_account_info"`
	AccrualAccount       string `json:"accrual_account"`
	SelectAccrualAccount string `json:"select_accrual_account"`
	AccrualAccountInfo   string `json:"accrual_account_info"`

	// §7 Notes
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notes_placeholder"`
	NotesInfo        string `json:"notes_info"`

	// Buttons
	Edit      string `json:"edit"`
	EditTitle string `json:"edit_title"`
	Active    string `json:"active"`
}

type StatusLabels struct {
	Outstanding string `json:"outstanding"`
	Partial     string `json:"partial"`
	Settled     string `json:"settled"`
	Reversed    string `json:"reversed"`
}

type ActionLabels struct {
	View               string `json:"view"`
	Edit               string `json:"edit"`
	Delete             string `json:"delete"`
	AccrueFromContract string `json:"accrue_from_contract"`
	Settle             string `json:"settle"`
	Reverse            string `json:"reverse"`
	AddSettlement      string `json:"add_settlement"`
	NoPermission       string `json:"no_permission"`
}

type ConfirmLabels struct {
	Delete         string `json:"delete"`
	DeleteMessage  string `json:"delete_message"`
	Settle         string `json:"settle"`
	SettleMessage  string `json:"settle_message"`
	Reverse        string `json:"reverse"`
	ReverseMessage string `json:"reverse_message"`
}

type BalanceLabels struct {
	Title       string `json:"title"`
	Accrued     string `json:"accrued"`
	Settled     string `json:"settled"`
	Remaining   string `json:"remaining"`
	Utilization string `json:"utilization"`
}

type EmptyLabels struct {
	Title              string `json:"title"`
	Message            string `json:"message"`
	OutstandingTitle   string `json:"outstanding_title"`
	OutstandingMessage string `json:"outstanding_message"`
	PartialTitle       string `json:"partial_title"`
	PartialMessage     string `json:"partial_message"`
	SettledTitle       string `json:"settled_title"`
	SettledMessage     string `json:"settled_message"`
	ReversedTitle      string `json:"reversed_title"`
	ReversedMessage    string `json:"reversed_message"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permission_denied"`
	InvalidFormData  string `json:"invalid_form_data"`
	NotFound         string `json:"not_found"`
	IDRequired       string `json:"id_required"`
	NoPermission     string `json:"no_permission"`
	CreationFailed   string `json:"creation_failed"`
	UpdateFailed     string `json:"update_failed"`
	DeletionFailed   string `json:"deletion_failed"`
	ListFailed       string `json:"list_failed"`
	SettleFailed     string `json:"settle_failed"`
	ReverseFailed    string `json:"reverse_failed"`
	BalanceDrift     string `json:"balance_drift"`
	LoadFailed       string `json:"load_failed"`
}

// DefaultLabels returns English fallback labels.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:            "Accrued Expenses",
			Caption:            "Recognized supplier obligations awaiting the actual bill",
			HeadingOutstanding: "Outstanding Accruals",
			HeadingPartial:     "Partially Settled",
			HeadingSettled:     "Settled Accruals",
			HeadingReversed:    "Reversed Accruals",
			Dashboard:          "Accrued Expense Dashboard",
		},
		Buttons: ButtonLabels{
			Add:                "New Accrual",
			AccrueFromContract: "Accrue from Contract",
			Settle:             "Settle",
			Reverse:            "Reverse",
			AddSettlement:      "Record Settlement",
		},
		Columns: ColumnLabels{
			InternalID:       "ID",
			Name:             "Name",
			Supplier:         "Supplier",
			SupplierContract: "Contract",
			RecognitionDate:  "Recognition Date",
			PeriodStart:      "Period Start",
			PeriodEnd:        "Period End",
			CycleDate:        "Cycle",
			Currency:         "Currency",
			AccruedAmount:    "Accrued",
			SettledAmount:    "Settled",
			RemainingAmount:  "Remaining",
			Status:           "Status",
		},
		Tabs: TabLabels{
			Info:        "Information",
			Settlements: "Settlements",
			Source:      "Source",
			Activity:    "Activity",
		},
		Detail: DetailLabels{
			PageTitle:            "Accrual Details",
			Title:                "Accrual Detail",
			InfoSection:          "Accrual Information",
			SettlementsSection:   "Settlements",
			SourceSection:        "Source",
			AuditTrailComingSoon: "Activity log feature coming soon.",
			AuditEmptyTitle:      "No activity entries",
			AuditEmptyMessage:    "Activity logs for this accrual will appear here.",
			TabAttachments:       "Attachments",
		},
		Settlements: SettlementLabels{
			Expenditure:                "Bill",
			AmountSettled:              "Amount Settled",
			Currency:                   "Currency",
			FxRate:                     "FX Rate",
			FxAdjustmentAmount:         "FX Adjustment",
			SettledAt:                  "Settled At",
			Reversal:                   "Reversal",
			EmptyTitle:                 "No settlements yet",
			EmptyMessage:               "Settlements applied against this accrual will appear here.",
			AddSettlement:              "Record Settlement",
			FormExpenditure:            "Bill",
			FormExpenditurePlaceholder: "Select bill...",
			FormAmountSettled:          "Amount Settled",
			FormCurrency:               "Currency",
			FormFxRate:                 "FX Rate",
			FormFxRateInfo:             "Bill currency to accrual currency conversion rate.",
			FormReversalReason:         "Reversal Reason",
		},
		Form: FormLabels{
			SectionIdentity:          "Accrual Identity",
			SectionSource:            "Source",
			SectionPeriod:            "Period",
			SectionMoney:             "Money",
			SectionAccounting:        "Accounting",
			SectionLifecycle:         "Lifecycle",
			SectionNotes:             "Notes",
			Name:                     "Accrual Name",
			NamePlaceholder:          "e.g. Utilities — May 2026 (estimate)",
			NameInfo:                 "Descriptive label for the accrued obligation.",
			Description:              "Description",
			DescriptionPlaceholder:   "Optional details about this accrual...",
			InternalID:               "Internal ID",
			InternalIDPlaceholder:    "Auto-generated",
			SupplierContract:         "Source Contract",
			SelectSupplierContract:   "Select contract...",
			SupplierContractInfo:     "The contract whose cycle is being accrued.",
			Supplier:                 "Supplier",
			SelectSupplier:           "Select supplier...",
			SupplierInfo:             "Supplier the obligation is owed to.",
			RecognitionDate:          "Recognition Date",
			RecognitionDateInfo:      "Date the obligation is booked into the period.",
			PeriodStart:              "Period Start",
			PeriodStartInfo:          "Start of the period the accrual covers.",
			PeriodEnd:                "Period End",
			PeriodEndInfo:            "End of the period the accrual covers.",
			CycleDate:                "Cycle Date",
			CycleDatePlaceholder:     "YYYY-MM-DD",
			CycleDateInfo:            "Cycle bucket used for idempotency.",
			Currency:                 "Currency",
			CurrencyPlaceholder:      "PHP",
			CurrencyInfo:             "ISO 4217 currency for the accrual.",
			AccruedAmount:            "Accrued Amount",
			AccruedAmountPlaceholder: "0.00",
			AccruedAmountInfo:        "Estimated obligation for the period.",
			SettledAmount:            "Settled Amount",
			SettledAmountInfo:        "Sum of settlements applied so far. Read-only.",
			RemainingAmount:          "Remaining",
			RemainingAmountInfo:      "Accrued minus settled. Read-only.",
			Status:                   "Status",
			SelectStatus:             "Select status...",
			StatusInfo:               "Lifecycle state.",
			StatusOutstanding:        "Outstanding",
			StatusPartial:            "Partially Settled",
			StatusSettled:            "Settled",
			StatusReversed:           "Reversed",
			ExpenseAccount:           "Expense Account",
			SelectExpenseAccount:     "Select expense account...",
			ExpenseAccountInfo:       "GL account to debit on recognition.",
			AccrualAccount:           "Accrual Account",
			SelectAccrualAccount:     "Select accrual account...",
			AccrualAccountInfo:       "GL account to credit while outstanding.",
			Notes:                    "Notes",
			NotesPlaceholder:         "Internal notes about this accrual...",
			NotesInfo:                "Internal remarks only.",
			Edit:                     "Edit",
			EditTitle:                "Edit Accrual",
			Active:                   "Active",
		},
		Status: StatusLabels{
			Outstanding: "Outstanding",
			Partial:     "Partially Settled",
			Settled:     "Settled",
			Reversed:    "Reversed",
		},
		Actions: ActionLabels{
			View:               "View Accrual",
			Edit:               "Edit Accrual",
			Delete:             "Delete Accrual",
			AccrueFromContract: "Accrue from Contract",
			Settle:             "Settle Accrual",
			Reverse:            "Reverse Accrual",
			AddSettlement:      "Record Settlement",
			NoPermission:       "No permission",
		},
		Confirm: ConfirmLabels{
			Delete:         "Delete Accrual",
			DeleteMessage:  "Are you sure you want to delete this accrual?",
			Settle:         "Settle Accrual",
			SettleMessage:  "Apply a settlement against this accrual?",
			Reverse:        "Reverse Accrual",
			ReverseMessage: "Reversing flips this accrual to Reversed and posts a reversing journal. Continue?",
		},
		Balances: BalanceLabels{
			Title:       "Settlement Progress",
			Accrued:     "Accrued",
			Settled:     "Settled",
			Remaining:   "Remaining",
			Utilization: "Settlement Progress",
		},
		Empty: EmptyLabels{
			Title:              "No accrued expenses yet",
			Message:            "Accrue from a contract cycle when the period closes before the supplier bill arrives.",
			OutstandingTitle:   "No outstanding accruals",
			OutstandingMessage: "Outstanding accruals waiting for a supplier bill will appear here.",
			PartialTitle:       "No partially settled accruals",
			PartialMessage:     "Accruals with at least one settlement will appear here.",
			SettledTitle:       "No settled accruals",
			SettledMessage:     "Fully reconciled accruals will appear here.",
			ReversedTitle:      "No reversed accruals",
			ReversedMessage:    "Reversed accruals will appear here.",
		},
		Errors: ErrorLabels{
			PermissionDenied: "You do not have permission to perform this action.",
			InvalidFormData:  "Invalid form data. Please check your inputs and try again.",
			NotFound:         "Accrued expense not found.",
			IDRequired:       "Accrual ID is required.",
			NoPermission:     "No permission.",
			CreationFailed:   "Accrual creation failed",
			UpdateFailed:     "Accrual update failed",
			DeletionFailed:   "Accrual deletion failed",
			ListFailed:       "Failed to retrieve accruals",
			SettleFailed:     "Settlement recording failed",
			ReverseFailed:    "Accrual reversal failed",
			BalanceDrift:     "Accrual balance drift detected; recompute required",
			LoadFailed:       "Failed to load accrual",
		},
	}
}
