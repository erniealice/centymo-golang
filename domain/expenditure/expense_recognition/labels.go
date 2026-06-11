package expense_recognition

// ---------------------------------------------------------------------------
// ExpenseRecognition labels  (SPS P10)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the
// expense_recognition module. Loaded from lyngua key root "expenseRecognition".
type Labels struct {
	Page    PageLabels    `json:"page"`
	Buttons ButtonLabels  `json:"buttons"`
	Columns ColumnLabels  `json:"columns"`
	Tabs    TabLabels     `json:"tabs"`
	Detail  DetailLabels  `json:"detail"`
	Lines   LineLabels    `json:"lines"`
	Source  SourceLabels  `json:"source"`
	Status  StatusLabels  `json:"status"`
	Actions ActionLabels  `json:"actions"`
	Confirm ConfirmLabels `json:"confirm"`
	Empty   EmptyLabels   `json:"empty"`
	Errors  ErrorLabels   `json:"errors"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	Caption         string `json:"caption"`
	HeadingDraft    string `json:"headingDraft"`
	HeadingPosted   string `json:"headingPosted"`
	HeadingReversed string `json:"headingReversed"`
	Dashboard       string `json:"dashboard"`
}

type ButtonLabels struct {
	Add                      string `json:"add"`
	RecognizeFromExpenditure string `json:"recognizeFromExpenditure"`
	RecognizeFromContract    string `json:"recognizeFromContract"`
	Reverse                  string `json:"reverse"`
}

type ColumnLabels struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	RecognitionDate  string `json:"recognitionDate"`
	PeriodStart      string `json:"periodStart"`
	PeriodEnd        string `json:"periodEnd"`
	CycleDate        string `json:"cycleDate"`
	Supplier         string `json:"supplier"`
	SupplierContract string `json:"supplierContract"`
	Expenditure      string `json:"expenditure"`
	Currency         string `json:"currency"`
	TotalAmount      string `json:"totalAmount"`
	Status           string `json:"status"`
	Source           string `json:"source"`
	IdempotencyKey   string `json:"idempotencyKey"`
}

type TabLabels struct {
	Info     string `json:"info"`
	Lines    string `json:"lines"`
	Source   string `json:"source"`
	Activity string `json:"activity"`
}

type DetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	Title                string `json:"title"`
	InfoSection          string `json:"infoSection"`
	SourceSection        string `json:"sourceSection"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	TabAttachments       string `json:"tabAttachments"`

	// Info-tab + source-tab field labels (4.4)
	Notes           string `json:"notes"`
	SourceContract  string `json:"sourceContract"`
	SourceBill      string `json:"sourceBill"`
	DeferredExpense string `json:"deferredExpense"`
	SourceAccrual   string `json:"sourceAccrual"`
	ReversalOf      string `json:"reversalOf"`
}

type LineLabels struct {
	Description    string `json:"description"`
	Quantity       string `json:"quantity"`
	UnitAmount     string `json:"unitAmount"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Product        string `json:"product"`
	ExpenseAccount string `json:"expenseAccount"`
	EmptyTitle     string `json:"emptyTitle"`
	EmptyMessage   string `json:"emptyMessage"`
	AddLine        string `json:"addLine"`

	// Drawer form labels
	FormDescription            string `json:"formDescription"`
	FormDescriptionPlaceholder string `json:"formDescriptionPlaceholder"`
	FormQuantity               string `json:"formQuantity"`
	FormUnitAmount             string `json:"formUnitAmount"`
	FormAmount                 string `json:"formAmount"`
	FormCurrency               string `json:"formCurrency"`
}

type SourceLabels struct {
	Recurrence  string `json:"recurrence"`
	Expenditure string `json:"expenditure"`
	Manual      string `json:"manual"`
	Reversal    string `json:"reversal"`
}

type StatusLabels struct {
	Draft    string `json:"draft"`
	Posted   string `json:"posted"`
	Reversed string `json:"reversed"`
}

type ActionLabels struct {
	View                     string `json:"view"`
	Edit                     string `json:"edit"`
	Delete                   string `json:"delete"`
	Reverse                  string `json:"reverse"`
	RecognizeFromExpenditure string `json:"recognizeFromExpenditure"`
	RecognizeFromContract    string `json:"recognizeFromContract"`
	NoPermission             string `json:"noPermission"`
}

type ConfirmLabels struct {
	Delete         string `json:"delete"`
	DeleteMessage  string `json:"deleteMessage"`
	Reverse        string `json:"reverse"`
	ReverseMessage string `json:"reverseMessage"`
}

type EmptyLabels struct {
	Title           string `json:"title"`
	Message         string `json:"message"`
	DraftTitle      string `json:"draftTitle"`
	DraftMessage    string `json:"draftMessage"`
	PostedTitle     string `json:"postedTitle"`
	PostedMessage   string `json:"postedMessage"`
	ReversedTitle   string `json:"reversedTitle"`
	ReversedMessage string `json:"reversedMessage"`
}

type ErrorLabels struct {
	PermissionDenied     string `json:"permissionDenied"`
	InvalidFormData      string `json:"invalidFormData"`
	NotFound             string `json:"notFound"`
	IDRequired           string `json:"idRequired"`
	NoPermission         string `json:"noPermission"`
	CreationFailed       string `json:"creation_failed"`
	UpdateFailed         string `json:"update_failed"`
	DeletionFailed       string `json:"deletion_failed"`
	ListFailed           string `json:"list_failed"`
	ReverseFailed        string `json:"reverse_failed"`
	IdempotencyCollision string `json:"idempotency_collision"`
	LoadFailed           string `json:"load_failed"`
}

// DefaultLabels returns English fallback labels.
// Tier overrides belong in lyngua JSON.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:         "Expense Recognition",
			Caption:         "Period in which a supplier cost is recognized",
			HeadingDraft:    "Draft Recognitions",
			HeadingPosted:   "Posted Recognitions",
			HeadingReversed: "Reversed Recognitions",
			Dashboard:       "Expense Recognition Dashboard",
		},
		Buttons: ButtonLabels{
			Add:                      "New Recognition",
			RecognizeFromExpenditure: "Recognize from Bill",
			RecognizeFromContract:    "Recognize from Contract",
			Reverse:                  "Reverse Recognition",
		},
		Columns: ColumnLabels{
			InternalID:       "ID",
			Name:             "Name",
			RecognitionDate:  "Recognition Date",
			PeriodStart:      "Period Start",
			PeriodEnd:        "Period End",
			CycleDate:        "Cycle",
			Supplier:         "Supplier",
			SupplierContract: "Contract",
			Expenditure:      "Source Bill",
			Currency:         "Currency",
			TotalAmount:      "Amount",
			Status:           "Status",
			Source:           "Source",
			IdempotencyKey:   "Idempotency Key",
		},
		Tabs: TabLabels{
			Info:     "Information",
			Lines:    "Recognition Lines",
			Source:   "Source",
			Activity: "Activity",
		},
		Detail: DetailLabels{
			PageTitle:            "Recognition Details",
			Title:                "Recognition Detail",
			InfoSection:          "Recognition Information",
			SourceSection:        "Source",
			AuditTrailComingSoon: "Activity log feature coming soon.",
			AuditEmptyTitle:      "No activity entries",
			AuditEmptyMessage:    "Activity logs for this recognition will appear here.",
			TabAttachments:       "Attachments",
		},
		Lines: LineLabels{
			Description:                "Description",
			Quantity:                   "Quantity",
			UnitAmount:                 "Unit Amount",
			Amount:                     "Amount",
			Currency:                   "Currency",
			Product:                    "Product",
			ExpenseAccount:             "Expense Account",
			EmptyTitle:                 "No recognition lines",
			EmptyMessage:               "Lines breaking down this recognition will appear here.",
			AddLine:                    "Add Line",
			FormDescription:            "Description",
			FormDescriptionPlaceholder: "e.g. Cloud hosting — May 2026",
			FormQuantity:               "Quantity",
			FormUnitAmount:             "Unit Amount",
			FormAmount:                 "Amount",
			FormCurrency:               "Currency",
		},
		Source: SourceLabels{
			Recurrence:  "Recurrence Engine",
			Expenditure: "From Bill",
			Manual:      "Manual",
			Reversal:    "Reversal",
		},
		Status: StatusLabels{
			Draft:    "Draft",
			Posted:   "Posted",
			Reversed: "Reversed",
		},
		Actions: ActionLabels{
			View:                     "View Recognition",
			Edit:                     "Edit Recognition",
			Delete:                   "Delete Recognition",
			Reverse:                  "Reverse",
			RecognizeFromExpenditure: "Recognize from Bill",
			RecognizeFromContract:    "Recognize from Contract",
			NoPermission:             "No permission",
		},
		Confirm: ConfirmLabels{
			Delete:         "Delete Recognition",
			DeleteMessage:  "Are you sure you want to delete this recognition? Only Draft recognitions can be deleted.",
			Reverse:        "Reverse Recognition",
			ReverseMessage: "Reversing creates a counter-entry and marks this recognition as Reversed. Continue?",
		},
		Empty: EmptyLabels{
			Title:           "No recognitions yet",
			Message:         "Recognize an expense from a posted bill or directly from a contract cycle.",
			DraftTitle:      "No draft recognitions",
			DraftMessage:    "Draft recognitions will appear here.",
			PostedTitle:     "No posted recognitions",
			PostedMessage:   "Posted recognitions will appear here.",
			ReversedTitle:   "No reversed recognitions",
			ReversedMessage: "Reversed recognitions will appear here.",
		},
		Errors: ErrorLabels{
			PermissionDenied:     "You do not have permission to perform this action.",
			InvalidFormData:      "Invalid form data. Please check your inputs and try again.",
			NotFound:             "Recognition not found.",
			IDRequired:           "Recognition ID is required.",
			NoPermission:         "No permission.",
			CreationFailed:       "Recognition creation failed",
			UpdateFailed:         "Recognition update failed",
			DeletionFailed:       "Recognition deletion failed",
			ListFailed:           "Failed to retrieve recognitions",
			ReverseFailed:        "Recognition reversal failed",
			IdempotencyCollision: "A recognition for this source and period already exists",
			LoadFailed:           "Failed to load recognition",
		},
	}
}
