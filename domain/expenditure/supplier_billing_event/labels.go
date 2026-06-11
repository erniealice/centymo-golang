package supplier_billing_event

// supplier_billing_event_labels.go — SupplierBillingEvent label structs
// (expenditure-domain). Moved here in centymo W6 from the root advance_labels.go
// (which was then deleted). SupplierBillingEvent is an EXPENDITURE-domain entity
// per esqyma (proto/v1/domain/expenditure/supplier_billing_event/). The treasury
// half of the historic advance_labels.go (AdvanceKind/AdvanceStatus enums, the
// "Advance Schedule" tab, the cash "Advances Dashboard" labels) landed in
// domain/treasury/advance.go in W5. Pure structural move — no behaviour change.
//
// Lyngua JSON source: packages/lyngua/translations/en/general/supplier_billing_event.json

// === SupplierBillingEvent labels ===

// StatusLabels — labels for the 5-state status enum.
type StatusLabels struct {
	Unspecified string `json:"unspecified"`
	Ready       string `json:"ready"`
	Billed      string `json:"billed"`
	Waived      string `json:"waived"`
	Cancelled   string `json:"cancelled"`
}

// TriggerLabels — labels for the trigger enum.
type TriggerLabels struct {
	Unspecified string `json:"unspecified"`
	ManualEarly string `json:"manualEarly"`
	ManualLate  string `json:"manualLate"`
}

// ColumnLabels — list table column headers.
type ColumnLabels struct {
	ID                   string `json:"id"`
	SupplierSubscription string `json:"supplierSubscription"`
	SupplierContract     string `json:"supplierContract"`
	BillableAmount       string `json:"billableAmount"`
	Currency             string `json:"currency"`
	Status               string `json:"status"`
	Trigger              string `json:"trigger"`
	ExpenseRecognition   string `json:"expenseRecognition"`
	DateCreated          string `json:"dateCreated"`
	Actions              string `json:"actions"`
}

// PageLabels — list / detail page strings.
type PageLabels struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
}

// ActionLabels — row-action button labels.
type ActionLabels struct {
	Recognize string `json:"recognize"`
	MarkReady string `json:"markReady"`
	Waive     string `json:"waive"`
	Cancel    string `json:"cancel"`
}

// DetailLabels — detail page tab + section labels.
type DetailLabels struct {
	Title              string `json:"title"`
	TabInfo            string `json:"tabInfo"`
	TabAudit           string `json:"tabAudit"`
	InfoHeading        string `json:"infoHeading"`
	LinkedAdvanceBadge string `json:"linkedAdvanceBadge"`
}

// EmptyLabels — empty-state labels for the list view.
type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// ErrorLabels — error toasts / validations.
type ErrorLabels struct {
	PermissionDenied  string `json:"permissionDenied"`
	NotFound          string `json:"notFound"`
	AlreadyRecognized string `json:"alreadyRecognized"`
	InvalidTransition string `json:"invalidTransition"`
}

// Labels — root struct for supplier_billing_event.json.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Columns ColumnLabels  `json:"columns"`
	Status  StatusLabels  `json:"status"`
	Trigger TriggerLabels `json:"trigger"`
	Actions ActionLabels  `json:"actions"`
	Detail  DetailLabels  `json:"detail"`
	Empty   EmptyLabels   `json:"empty"`
	Errors  ErrorLabels   `json:"errors"`
}

// === Defaults ===

// DefaultLabels — English defaults for the supplier
// billing-event list / detail screens that anchor MILESTONE recognition on
// the buying side.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:   "Supplier Billing Events",
			Caption: "Manual billing triggers for supplier subscriptions and contracts.",
		},
		Columns: ColumnLabels{
			ID:                   "ID",
			SupplierSubscription: "Supplier subscription",
			SupplierContract:     "Supplier contract",
			BillableAmount:       "Billable amount",
			Currency:             "Currency",
			Status:               "Status",
			Trigger:              "Trigger",
			ExpenseRecognition:   "Expense recognition",
			DateCreated:          "Created",
			Actions:              "Actions",
		},
		Status: StatusLabels{
			Unspecified: "Unspecified",
			Ready:       "Ready",
			Billed:      "Billed",
			Waived:      "Waived",
			Cancelled:   "Cancelled",
		},
		Trigger: TriggerLabels{
			Unspecified: "Unspecified",
			ManualEarly: "Manual (early)",
			ManualLate:  "Manual (late)",
		},
		Actions: ActionLabels{
			Recognize: "Recognize",
			MarkReady: "Mark ready",
			Waive:     "Waive",
			Cancel:    "Cancel",
		},
		Detail: DetailLabels{
			Title:              "Supplier billing event",
			TabInfo:            "Info",
			TabAudit:           "Audit",
			InfoHeading:        "Billing event details",
			LinkedAdvanceBadge: "Linked advance",
		},
		Empty: EmptyLabels{
			Title:   "No supplier billing events",
			Message: "Supplier billing events appear here once created.",
		},
		Errors: ErrorLabels{
			PermissionDenied:  "Permission denied",
			NotFound:          "Supplier billing event not found",
			AlreadyRecognized: "This billing event has already been recognized",
			InvalidTransition: "Cannot transition from %s to %s",
		},
	}
}
