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
	HeadingOutstanding string `json:"headingOutstanding"`
	HeadingPartial     string `json:"headingPartial"`
	HeadingSettled     string `json:"headingSettled"`
	HeadingReversed    string `json:"headingReversed"`
	Dashboard          string `json:"dashboard"`
}

type ButtonLabels struct {
	Add                string `json:"add"`
	AccrueFromContract string `json:"accrueFromContract"`
	Settle             string `json:"settle"`
	Reverse            string `json:"reverse"`
	AddSettlement      string `json:"addSettlement"`
}

type ColumnLabels struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	Supplier         string `json:"supplier"`
	SupplierContract string `json:"supplierContract"`
	RecognitionDate  string `json:"recognitionDate"`
	PeriodStart      string `json:"periodStart"`
	PeriodEnd        string `json:"periodEnd"`
	CycleDate        string `json:"cycleDate"`
	Currency         string `json:"currency"`
	AccruedAmount    string `json:"accruedAmount"`
	SettledAmount    string `json:"settledAmount"`
	RemainingAmount  string `json:"remainingAmount"`
	Status           string `json:"status"`
}

type TabLabels struct {
	Info        string `json:"info"`
	Settlements string `json:"settlements"`
	Source      string `json:"source"`
	Activity    string `json:"activity"`
}

type DetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	Title                string `json:"title"`
	InfoSection          string `json:"infoSection"`
	SettlementsSection   string `json:"settlementsSection"`
	SourceSection        string `json:"sourceSection"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	TabAttachments       string `json:"tabAttachments"`

	// Info-tab + source-tab field labels (4.4)
	Notes          string `json:"notes"`
	SourceContract string `json:"sourceContract"`
	Supplier       string `json:"supplier"`
	ExpenseAccount string `json:"expenseAccount"`
	AccrualAccount string `json:"accrualAccount"`
}

type SettlementLabels struct {
	Expenditure        string `json:"expenditure"`
	AmountSettled      string `json:"amountSettled"`
	Currency           string `json:"currency"`
	FxRate             string `json:"fxRate"`
	FxAdjustmentAmount string `json:"fxAdjustmentAmount"`
	SettledAt          string `json:"settledAt"`
	Reversal           string `json:"reversal"`
	EmptyTitle         string `json:"emptyTitle"`
	EmptyMessage       string `json:"emptyMessage"`
	AddSettlement      string `json:"addSettlement"`

	// Drawer form labels
	FormExpenditure            string `json:"formExpenditure"`
	FormExpenditurePlaceholder string `json:"formExpenditurePlaceholder"`
	FormAmountSettled          string `json:"formAmountSettled"`
	FormCurrency               string `json:"formCurrency"`
	FormFxRate                 string `json:"formFxRate"`
	FormFxRateInfo             string `json:"formFxRateInfo"`
	FormReversalReason         string `json:"formReversalReason"`
}

type FormLabels struct {
	// Section headers
	SectionIdentity   string `json:"sectionIdentity"`
	SectionSource     string `json:"sectionSource"`
	SectionPeriod     string `json:"sectionPeriod"`
	SectionMoney      string `json:"sectionMoney"`
	SectionAccounting string `json:"sectionAccounting"`
	SectionLifecycle  string `json:"sectionLifecycle"`
	SectionNotes      string `json:"sectionNotes"`

	// §1 Identity
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	NameInfo               string `json:"nameInfo"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	InternalID             string `json:"internalId"`
	InternalIDPlaceholder  string `json:"internalIdPlaceholder"`

	// §2 Source
	SupplierContract       string `json:"supplierContract"`
	SelectSupplierContract string `json:"selectSupplierContract"`
	SupplierContractInfo   string `json:"supplierContractInfo"`
	Supplier               string `json:"supplier"`
	SelectSupplier         string `json:"selectSupplier"`
	SupplierInfo           string `json:"supplierInfo"`

	// §3 Period
	RecognitionDate      string `json:"recognitionDate"`
	RecognitionDateInfo  string `json:"recognitionDateInfo"`
	PeriodStart          string `json:"periodStart"`
	PeriodStartInfo      string `json:"periodStartInfo"`
	PeriodEnd            string `json:"periodEnd"`
	PeriodEndInfo        string `json:"periodEndInfo"`
	CycleDate            string `json:"cycleDate"`
	CycleDatePlaceholder string `json:"cycleDatePlaceholder"`
	CycleDateInfo        string `json:"cycleDateInfo"`

	// §4 Money
	Currency                 string `json:"currency"`
	CurrencyPlaceholder      string `json:"currencyPlaceholder"`
	CurrencyInfo             string `json:"currencyInfo"`
	AccruedAmount            string `json:"accruedAmount"`
	AccruedAmountPlaceholder string `json:"accruedAmountPlaceholder"`
	AccruedAmountInfo        string `json:"accruedAmountInfo"`
	SettledAmount            string `json:"settledAmount"`
	SettledAmountInfo        string `json:"settledAmountInfo"`
	RemainingAmount          string `json:"remainingAmount"`
	RemainingAmountInfo      string `json:"remainingAmountInfo"`

	// §5 Lifecycle
	Status            string `json:"status"`
	SelectStatus      string `json:"selectStatus"`
	StatusInfo        string `json:"statusInfo"`
	StatusOutstanding string `json:"statusOutstanding"`
	StatusPartial     string `json:"statusPartial"`
	StatusSettled     string `json:"statusSettled"`
	StatusReversed    string `json:"statusReversed"`

	// §6 Accounting
	ExpenseAccount       string `json:"expenseAccount"`
	SelectExpenseAccount string `json:"selectExpenseAccount"`
	ExpenseAccountInfo   string `json:"expenseAccountInfo"`
	AccrualAccount       string `json:"accrualAccount"`
	SelectAccrualAccount string `json:"selectAccrualAccount"`
	AccrualAccountInfo   string `json:"accrualAccountInfo"`

	// §7 Notes
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	NotesInfo        string `json:"notesInfo"`

	// Buttons
	Edit      string `json:"edit"`
	EditTitle string `json:"editTitle"`
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
	AccrueFromContract string `json:"accrueFromContract"`
	Settle             string `json:"settle"`
	Reverse            string `json:"reverse"`
	AddSettlement      string `json:"addSettlement"`
	NoPermission       string `json:"noPermission"`
}

type ConfirmLabels struct {
	Delete         string `json:"delete"`
	DeleteMessage  string `json:"deleteMessage"`
	Settle         string `json:"settle"`
	SettleMessage  string `json:"settleMessage"`
	Reverse        string `json:"reverse"`
	ReverseMessage string `json:"reverseMessage"`
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
	OutstandingTitle   string `json:"outstandingTitle"`
	OutstandingMessage string `json:"outstandingMessage"`
	PartialTitle       string `json:"partialTitle"`
	PartialMessage     string `json:"partialMessage"`
	SettledTitle       string `json:"settledTitle"`
	SettledMessage     string `json:"settledMessage"`
	ReversedTitle      string `json:"reversedTitle"`
	ReversedMessage    string `json:"reversedMessage"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoPermission     string `json:"noPermission"`
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
