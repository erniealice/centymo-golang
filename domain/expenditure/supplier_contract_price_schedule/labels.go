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
	NamePlural string `json:"name_plural"`
	Line       string `json:"line"`
	LinePlural string `json:"line_plural"`
}

type PageLabels struct {
	Heading           string `json:"heading"`
	Caption           string `json:"caption"`
	HeadingScheduled  string `json:"heading_scheduled"`
	HeadingActive     string `json:"heading_active"`
	HeadingSuperseded string `json:"heading_superseded"`
	HeadingCancelled  string `json:"heading_cancelled"`
	TabTitle          string `json:"tab_title"`
}

type ButtonLabels struct {
	Add       string `json:"add"`
	AddLine   string `json:"add_line"`
	Activate  string `json:"activate"`
	Supersede string `json:"supersede"`
	Cancel    string `json:"cancel"`
}

type FilterLabels struct {
	All                 string `json:"all"`
	Status              string `json:"status"`
	AnyStatus           string `json:"any_status"`
	SupplierContract    string `json:"supplier_contract"`
	AnySupplierContract string `json:"any_supplier_contract"`
	DateRange           string `json:"date_range"`
}

type ColumnLabels struct {
	InternalID       string `json:"internal_id"`
	Name             string `json:"name"`
	SupplierContract string `json:"supplier_contract"`
	SequenceNumber   string `json:"sequence_number"`
	DateStart        string `json:"date_start"`
	DateEnd          string `json:"date_end"`
	Status           string `json:"status"`
	Currency         string `json:"currency"`
	LineCount        string `json:"line_count"`
	Total            string `json:"total"`
}

type EmptyLabels struct {
	Title             string `json:"title"`
	Message           string `json:"message"`
	ScheduledTitle    string `json:"scheduled_title"`
	ScheduledMessage  string `json:"scheduled_message"`
	ActiveTitle       string `json:"active_title"`
	ActiveMessage     string `json:"active_message"`
	SupersededTitle   string `json:"superseded_title"`
	SupersededMessage string `json:"superseded_message"`
	CancelledTitle    string `json:"cancelled_title"`
	CancelledMessage  string `json:"cancelled_message"`
}

type FormLabels struct {
	// Section headers
	SectionIdentity  string `json:"section_identity"`
	SectionValidity  string `json:"section_validity"`
	SectionScoping   string `json:"section_scoping"`
	SectionLifecycle string `json:"section_lifecycle"`
	SectionNotes     string `json:"section_notes"`

	// Identity
	Name                   string `json:"name"`
	NamePlaceholder        string `json:"name_placeholder"`
	NameInfo               string `json:"name_info"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"description_placeholder"`
	InternalID             string `json:"internal_id"`
	InternalIDPlaceholder  string `json:"internal_id_placeholder"`
	InternalIDInfo         string `json:"internal_id_info"`

	// Scoping
	SupplierContract       string `json:"supplier_contract"`
	SelectSupplierContract string `json:"select_supplier_contract"`
	SupplierContractInfo   string `json:"supplier_contract_info"`

	// Validity
	DateStart          string `json:"date_start"`
	DateStartInfo      string `json:"date_start_info"`
	DateEnd            string `json:"date_end"`
	DateEndPlaceholder string `json:"date_end_placeholder"`
	DateEndInfo        string `json:"date_end_info"`

	// Currency / location
	Currency            string `json:"currency"`
	CurrencyPlaceholder string `json:"currency_placeholder"`
	CurrencyInfo        string `json:"currency_info"`
	Location            string `json:"location"`
	SelectLocation      string `json:"select_location"`
	LocationInfo        string `json:"location_info"`

	// Lifecycle
	Status                    string `json:"status"`
	SelectStatus              string `json:"select_status"`
	StatusInfo                string `json:"status_info"`
	SequenceNumber            string `json:"sequence_number"`
	SequenceNumberPlaceholder string `json:"sequence_number_placeholder"`
	SequenceNumberInfo        string `json:"sequence_number_info"`

	// Notes
	Notes            string `json:"notes"`
	NotesPlaceholder string `json:"notes_placeholder"`
	NotesInfo        string `json:"notes_info"`
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
	AddLine             string         `json:"add_line"`
	ColumnContractLine  string         `json:"column_contract_line"`
	ColumnUnitPrice     string         `json:"column_unit_price"`
	ColumnQuantity      string         `json:"column_quantity"`
	ColumnMinimumAmount string         `json:"column_minimum_amount"`
	ColumnCurrency      string         `json:"column_currency"`
	ColumnCycleOverride string         `json:"column_cycle_override"`
	LineForm            LineFormLabels `json:"line_form"`
}

type LineFormLabels struct {
	SectionLink                   string `json:"section_link"`
	SectionPricing                string `json:"section_pricing"`
	SectionCycle                  string `json:"section_cycle"`
	SupplierContractLine          string `json:"supplier_contract_line"`
	SelectSupplierContractLine    string `json:"select_supplier_contract_line"`
	SupplierContractLineInfo      string `json:"supplier_contract_line_info"`
	UnitPrice                     string `json:"unit_price"`
	UnitPricePlaceholder          string `json:"unit_price_placeholder"`
	UnitPriceInfo                 string `json:"unit_price_info"`
	MinimumAmount                 string `json:"minimum_amount"`
	MinimumAmountPlaceholder      string `json:"minimum_amount_placeholder"`
	MinimumAmountInfo             string `json:"minimum_amount_info"`
	Quantity                      string `json:"quantity"`
	QuantityPlaceholder           string `json:"quantity_placeholder"`
	QuantityInfo                  string `json:"quantity_info"`
	Currency                      string `json:"currency"`
	CurrencyPlaceholder           string `json:"currency_placeholder"`
	CycleValueOverride            string `json:"cycle_value_override"`
	CycleValueOverridePlaceholder string `json:"cycle_value_override_placeholder"`
	CycleValueOverrideInfo        string `json:"cycle_value_override_info"`
	CycleUnitOverride             string `json:"cycle_unit_override"`
	CycleUnitOverridePlaceholder  string `json:"cycle_unit_override_placeholder"`
	CycleUnitOverrideInfo         string `json:"cycle_unit_override_info"`
}

type DetailLabels struct {
	PageTitle            string `json:"page_title"`
	Title                string `json:"title"`
	InfoSection          string `json:"info_section"`
	LinesSection         string `json:"lines_section"`
	AuditTrailComingSoon string `json:"audit_trail_coming_soon"`
	AuditEmptyTitle      string `json:"audit_empty_title"`
	AuditEmptyMessage    string `json:"audit_empty_message"`
	TabAttachments       string `json:"tab_attachments"`
}

type ErrorLabels struct {
	PermissionDenied    string `json:"permission_denied"`
	InvalidFormData     string `json:"invalid_form_data"`
	NotFound            string `json:"not_found"`
	IDRequired          string `json:"id_required"`
	NoPermission        string `json:"no_permission"`
	CannotDelete        string `json:"cannot_delete"`
	InUse               string `json:"in_use"`
	CreationFailed      string `json:"creation_failed"`
	UpdateFailed        string `json:"update_failed"`
	DeletionFailed      string `json:"deletion_failed"`
	ListFailed          string `json:"list_failed"`
	AuthorizationFailed string `json:"authorization_failed"`
	ActivationFailed    string `json:"activation_failed"`
	SupersedeFailed     string `json:"supersede_failed"`
	OverlapDetected     string `json:"overlap_detected"`
	LoadFailed          string `json:"load_failed"`
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
