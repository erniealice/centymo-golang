package supplier_contract_price_schedule

// ---------------------------------------------------------------------------
// SPS P7 — SupplierContractPriceSchedule labels
// (mirrors lyngua/translations/en/general/supplier_contract_price_schedule.json
//  root key "supplierContractPriceSchedule")
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the
// supplier_contract_price_schedule + child line views.
type Labels struct {
	Labels  NounLabels   `json:"labels"`
	Page    PageLabels   `json:"page"`
	Buttons ButtonLabels `json:"buttons"`
	Filters FilterLabels `json:"filters"`
	Columns ColumnLabels `json:"columns"`
	Empty   EmptyLabels  `json:"empty"`
	Form    FormLabels   `json:"form"`
	Status  StatusLabels `json:"status"`
	Tabs    TabLabels    `json:"tabs"`
	Lines   LinesLabels  `json:"lines"`
	Detail  DetailLabels `json:"detail"`
	Errors  ErrorLabels  `json:"errors"`
}

type NounLabels struct {
	Name       string `json:"name"`
	NamePlural string `json:"namePlural"`
	Line       string `json:"line"`
	LinePlural string `json:"linePlural"`
}

type PageLabels struct {
	Heading           string `json:"heading"`
	Caption           string `json:"caption"`
	HeadingScheduled  string `json:"headingScheduled"`
	HeadingActive     string `json:"headingActive"`
	HeadingSuperseded string `json:"headingSuperseded"`
	HeadingCancelled  string `json:"headingCancelled"`
	TabTitle          string `json:"tabTitle"`
}

type ButtonLabels struct {
	Add       string `json:"add"`
	AddLine   string `json:"addLine"`
	Activate  string `json:"activate"`
	Supersede string `json:"supersede"`
	Cancel    string `json:"cancel"`
}

type FilterLabels struct {
	All                 string `json:"all"`
	Status              string `json:"status"`
	AnyStatus           string `json:"anyStatus"`
	SupplierContract    string `json:"supplierContract"`
	AnySupplierContract string `json:"anySupplierContract"`
	DateRange           string `json:"dateRange"`
}

type ColumnLabels struct {
	InternalID       string `json:"internalId"`
	Name             string `json:"name"`
	SupplierContract string `json:"supplierContract"`
	SequenceNumber   string `json:"sequenceNumber"`
	DateStart        string `json:"dateStart"`
	DateEnd          string `json:"dateEnd"`
	Status           string `json:"status"`
	Currency         string `json:"currency"`
	LineCount        string `json:"lineCount"`
	Total            string `json:"total"`
}

type EmptyLabels struct {
	Title             string `json:"title"`
	Message           string `json:"message"`
	ScheduledTitle    string `json:"scheduledTitle"`
	ScheduledMessage  string `json:"scheduledMessage"`
	ActiveTitle       string `json:"activeTitle"`
	ActiveMessage     string `json:"activeMessage"`
	SupersededTitle   string `json:"supersededTitle"`
	SupersededMessage string `json:"supersededMessage"`
	CancelledTitle    string `json:"cancelledTitle"`
	CancelledMessage  string `json:"cancelledMessage"`
}

type FormLabels struct {
	// Section headers
	SectionIdentity  string `json:"sectionIdentity"`
	SectionValidity  string `json:"sectionValidity"`
	SectionScoping   string `json:"sectionScoping"`
	SectionLifecycle string `json:"sectionLifecycle"`
	SectionNotes     string `json:"sectionNotes"`

	// Identity
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"namePlaceholder"`
	NameInfo               string `json:"nameInfo"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	InternalID             string `json:"internalId"`
	InternalIDPlaceholder  string `json:"internalIdPlaceholder"`
	InternalIDInfo         string `json:"internalIdInfo"`

	// Scoping
	SupplierContract       string `json:"supplierContract"`
	SelectSupplierContract string `json:"selectSupplierContract"`
	SupplierContractInfo   string `json:"supplierContractInfo"`

	// Validity
	DateStart          string `json:"dateStart"`
	DateStartInfo      string `json:"dateStartInfo"`
	DateEnd            string `json:"dateEnd"`
	DateEndPlaceholder string `json:"dateEndPlaceholder"`
	DateEndInfo        string `json:"dateEndInfo"`

	// Currency / location
	Currency            string `json:"currency"`
	CurrencyPlaceholder string `json:"currencyPlaceholder"`
	CurrencyInfo        string `json:"currencyInfo"`
	Location            string `json:"location"`
	SelectLocation      string `json:"selectLocation"`
	LocationInfo        string `json:"locationInfo"`

	// Lifecycle
	Status                    string `json:"status"`
	SelectStatus              string `json:"selectStatus"`
	StatusInfo                string `json:"statusInfo"`
	SequenceNumber            string `json:"sequenceNumber"`
	SequenceNumberPlaceholder string `json:"sequenceNumberPlaceholder"`
	SequenceNumberInfo        string `json:"sequenceNumberInfo"`

	// Notes
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notesPlaceholder"`
	NotesInfo        string `json:"notesInfo"`
}

type StatusLabels struct {
	Scheduled  string `json:"scheduled"`
	Active     string `json:"active"`
	Superseded string `json:"superseded"`
	Cancelled  string `json:"cancelled"`
}

type TabLabels struct {
	Info     string `json:"info"`
	Lines    string `json:"lines"`
	Activity string `json:"activity"`
}

type LinesLabels struct {
	Title               string         `json:"title"`
	Empty               string         `json:"empty"`
	AddLine             string         `json:"addLine"`
	ColumnContractLine  string         `json:"columnContractLine"`
	ColumnUnitPrice     string         `json:"columnUnitPrice"`
	ColumnQuantity      string         `json:"columnQuantity"`
	ColumnMinimumAmount string         `json:"columnMinimumAmount"`
	ColumnCurrency      string         `json:"columnCurrency"`
	ColumnCycleOverride string         `json:"columnCycleOverride"`
	LineForm            LineFormLabels `json:"lineForm"`
}

type LineFormLabels struct {
	SectionLink                   string `json:"sectionLink"`
	SectionPricing                string `json:"sectionPricing"`
	SectionCycle                  string `json:"sectionCycle"`
	SupplierContractLine          string `json:"supplierContractLine"`
	SelectSupplierContractLine    string `json:"selectSupplierContractLine"`
	SupplierContractLineInfo      string `json:"supplierContractLineInfo"`
	UnitPrice                     string `json:"unitPrice"`
	UnitPricePlaceholder          string `json:"unitPricePlaceholder"`
	UnitPriceInfo                 string `json:"unitPriceInfo"`
	MinimumAmount                 string `json:"minimumAmount"`
	MinimumAmountPlaceholder      string `json:"minimumAmountPlaceholder"`
	MinimumAmountInfo             string `json:"minimumAmountInfo"`
	Quantity                      string `json:"quantity"`
	QuantityPlaceholder           string `json:"quantityPlaceholder"`
	QuantityInfo                  string `json:"quantityInfo"`
	Currency                      string `json:"currency"`
	CurrencyPlaceholder           string `json:"currencyPlaceholder"`
	CycleValueOverride            string `json:"cycleValueOverride"`
	CycleValueOverridePlaceholder string `json:"cycleValueOverridePlaceholder"`
	CycleValueOverrideInfo        string `json:"cycleValueOverrideInfo"`
	CycleUnitOverride             string `json:"cycleUnitOverride"`
	CycleUnitOverridePlaceholder  string `json:"cycleUnitOverridePlaceholder"`
	CycleUnitOverrideInfo         string `json:"cycleUnitOverrideInfo"`
}

type DetailLabels struct {
	PageTitle            string `json:"pageTitle"`
	Title                string `json:"title"`
	InfoSection          string `json:"infoSection"`
	LinesSection         string `json:"linesSection"`
	AuditTrailComingSoon string `json:"auditTrailComingSoon"`
	AuditEmptyTitle      string `json:"auditEmptyTitle"`
	AuditEmptyMessage    string `json:"auditEmptyMessage"`
	TabAttachments       string `json:"tabAttachments"`
}

type ErrorLabels struct {
	PermissionDenied    string `json:"permissionDenied"`
	InvalidFormData     string `json:"invalidFormData"`
	NotFound            string `json:"notFound"`
	IDRequired          string `json:"idRequired"`
	NoPermission        string `json:"noPermission"`
	CannotDelete        string `json:"cannotDelete"`
	InUse               string `json:"inUse"`
	CreationFailed      string `json:"creation_failed"`
	UpdateFailed        string `json:"update_failed"`
	DeletionFailed      string `json:"deletion_failed"`
	ListFailed          string `json:"list_failed"`
	AuthorizationFailed string `json:"authorization_failed"`
	ActivationFailed    string `json:"activation_failed"`
	SupersedeFailed     string `json:"supersede_failed"`
	OverlapDetected     string `json:"overlap_detected"`
	LoadFailed          string `json:"loadFailed"`
}

// DefaultLabels returns English fallback labels.
// Uses proto-generic naming — tier overrides belong in lyngua JSON.
func DefaultLabels() Labels {
	return Labels{
		Labels: NounLabels{
			Name:       "Price Schedule",
			NamePlural: "Price Schedules",
			Line:       "Schedule Line",
			LinePlural: "Schedule Lines",
		},
		Page: PageLabels{
			Heading:           "Contract Price Schedules",
			Caption:           "Date-windowed pricing layered on top of a supplier contract for multi-year escalation",
			HeadingScheduled:  "Scheduled Periods",
			HeadingActive:     "Active Periods",
			HeadingSuperseded: "Superseded Periods",
			HeadingCancelled:  "Cancelled Periods",
			TabTitle:          "Price Schedules",
		},
		Buttons: ButtonLabels{
			Add:       "New Schedule",
			AddLine:   "Add Schedule Line",
			Activate:  "Activate",
			Supersede: "Supersede",
			Cancel:    "Cancel Schedule",
		},
		Filters: FilterLabels{
			All:                 "All",
			Status:              "Status",
			AnyStatus:           "Any Status",
			SupplierContract:    "Contract",
			AnySupplierContract: "Any Contract",
			DateRange:           "Effective Window",
		},
		Columns: ColumnLabels{
			InternalID:       "ID",
			Name:             "Schedule Name",
			SupplierContract: "Contract",
			SequenceNumber:   "Seq.",
			DateStart:        "Start",
			DateEnd:          "End",
			Status:           "Status",
			Currency:         "Currency",
			LineCount:        "Lines",
			Total:            "Total",
		},
		Empty: EmptyLabels{
			Title:             "No price schedules yet",
			Message:           "Add a schedule period to layer multi-year pricing onto this contract.",
			ScheduledTitle:    "No upcoming schedules",
			ScheduledMessage:  "Future-dated schedules will appear here once added.",
			ActiveTitle:       "No active schedule",
			ActiveMessage:     "The contract is using header pricing — no schedule is in effect right now.",
			SupersededTitle:   "No past schedules",
			SupersededMessage: "Schedules whose window has passed will be archived here.",
			CancelledTitle:    "No cancelled schedules",
			CancelledMessage:  "Schedules cancelled before activation will appear here.",
		},
		Form: FormLabels{
			SectionIdentity:           "Schedule Identity",
			SectionValidity:           "Validity Window",
			SectionScoping:            "Scoping",
			SectionLifecycle:          "Lifecycle",
			SectionNotes:              "Notes",
			Name:                      "Schedule Name",
			NamePlaceholder:           "e.g. Year 1 (2026)",
			NameInfo:                  "Human-readable label for this pricing window. Often a year or renewal label.",
			Description:               "Description",
			DescriptionPlaceholder:    "Optional details about this pricing window...",
			InternalID:                "Internal ID",
			InternalIDPlaceholder:     "Auto-generated",
			InternalIDInfo:            "Auto-generated unique identifier.",
			SupplierContract:          "Contract",
			SelectSupplierContract:    "Select contract...",
			SupplierContractInfo:      "The supplier contract this pricing window applies to. Schedules are scoped to a single contract.",
			DateStart:                 "Start Date",
			DateStartInfo:             "When this pricing window takes effect. Window is half-open: start is inclusive.",
			DateEnd:                   "End Date",
			DateEndPlaceholder:        "Leave empty for open-ended",
			DateEndInfo:               "When this pricing window ends. Leave blank for the open-ended last bucket. Window is half-open: end is exclusive.",
			Currency:                  "Currency",
			CurrencyPlaceholder:       "PHP",
			CurrencyInfo:              "ISO 4217 currency for prices in this schedule.",
			Location:                  "Location",
			SelectLocation:            "Select location...",
			LocationInfo:              "Optional location override. Defaults to the parent contract's location.",
			Status:                    "Status",
			SelectStatus:              "Select status...",
			StatusInfo:                "Lifecycle state. New schedules default to Scheduled and progress to Active when their window arrives.",
			SequenceNumber:            "Sequence",
			SequenceNumberPlaceholder: "1",
			SequenceNumberInfo:        "Ordering position within the contract (1, 2, 3...).",
			Notes:                     "Notes",
			NotesPlaceholder:          "Internal notes about this pricing window...",
			NotesInfo:                 "Internal remarks only. Not visible to the supplier.",
		},
		Status: StatusLabels{
			Scheduled:  "Scheduled",
			Active:     "Active",
			Superseded: "Superseded",
			Cancelled:  "Cancelled",
		},
		Tabs: TabLabels{
			Info:     "Information",
			Lines:    "Schedule Lines",
			Activity: "Activity",
		},
		Lines: LinesLabels{
			Title:               "Schedule Lines",
			Empty:               "No schedule lines yet. Add lines to override per-line pricing for this window.",
			AddLine:             "Add Line",
			ColumnContractLine:  "Contract Line",
			ColumnUnitPrice:     "Unit Price",
			ColumnQuantity:      "Qty",
			ColumnMinimumAmount: "Minimum",
			ColumnCurrency:      "Currency",
			ColumnCycleOverride: "Cycle Override",
			LineForm: LineFormLabels{
				SectionLink:                   "Contract Line",
				SectionPricing:                "Pricing",
				SectionCycle:                  "Cycle Override",
				SupplierContractLine:          "Contract Line",
				SelectSupplierContractLine:    "Select contract line...",
				SupplierContractLineInfo:      "The line whose unit price is overridden during this window.",
				UnitPrice:                     "Unit Price",
				UnitPricePlaceholder:          "0.00",
				UnitPriceInfo:                 "Per-unit price during this window, in the schedule's currency.",
				MinimumAmount:                 "Minimum Amount",
				MinimumAmountPlaceholder:      "0.00",
				MinimumAmountInfo:             "For Minimum Commitment lines: the floor charged per cycle.",
				Quantity:                      "Quantity",
				QuantityPlaceholder:           "0",
				QuantityInfo:                  "Optional committed quantity for blanket or minimum-commitment lines.",
				Currency:                      "Currency",
				CurrencyPlaceholder:           "PHP",
				CycleValueOverride:            "Cycle Value Override",
				CycleValueOverridePlaceholder: "e.g. 1",
				CycleValueOverrideInfo:        "Optional cycle-length override. Most lines inherit the contract cycle.",
				CycleUnitOverride:             "Cycle Unit Override",
				CycleUnitOverridePlaceholder:  "month",
				CycleUnitOverrideInfo:         "Optional cycle-unit override (day, week, month, year).",
			},
		},
		Detail: DetailLabels{
			PageTitle:            "Schedule Details",
			Title:                "Schedule Detail",
			InfoSection:          "Schedule Information",
			LinesSection:         "Per-Line Pricing",
			AuditTrailComingSoon: "Activity log feature coming soon.",
			AuditEmptyTitle:      "No activity entries",
			AuditEmptyMessage:    "Activity logs for this schedule will appear here.",
			TabAttachments:       "Attachments",
		},
		Errors: ErrorLabels{
			PermissionDenied:    "You do not have permission to perform this action.",
			InvalidFormData:     "Invalid form data. Please check your inputs and try again.",
			NotFound:            "Price schedule not found.",
			IDRequired:          "Schedule ID is required.",
			NoPermission:        "No permission.",
			CannotDelete:        "Cannot delete — this schedule is currently active or has dependent lines.",
			InUse:               "Cannot delete — this schedule is referenced by existing records.",
			CreationFailed:      "Schedule creation failed",
			UpdateFailed:        "Schedule update failed",
			DeletionFailed:      "Schedule deletion failed",
			ListFailed:          "Failed to retrieve price schedules",
			AuthorizationFailed: "Authorization failed for price schedules",
			ActivationFailed:    "Schedule activation failed",
			SupersedeFailed:     "Schedule supersede failed",
			OverlapDetected:     "Schedule windows overlap; adjust dates and retry",
			LoadFailed:          "Failed to load price schedule",
		},
	}
}
