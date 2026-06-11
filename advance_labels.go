package centymo

// advance_labels.go — SupplierBillingEvent label structs (expenditure-domain).
//
// NOTE (centymo W5): the treasury-domain advance content that historically lived
// here — AdvanceKind / AdvanceStatus / AdvanceProrationPolicy enum labels, the
// "Advance Schedule" tab (TreasuryAdvanceLabels), the cash-app "Advances
// Dashboard" labels, and their Default* constructors — moved to
// domain/treasury/advance.go (Decision D, structure-mantra §8). What remains here
// is the SupplierBillingEvent* labels, which are EXPENDITURE-domain per esqyma
// (proto/v1/domain/expenditure/supplier_billing_event/) and relocate in W6.
//
// Lyngua JSON source: packages/lyngua/translations/en/general/supplier_billing_event.json

// === SupplierBillingEvent labels ===

// SupplierBillingEventStatusLabels — labels for the 5-state status enum.
type SupplierBillingEventStatusLabels struct {
	Unspecified string `json:"unspecified"`
	Ready       string `json:"ready"`
	Billed      string `json:"billed"`
	Waived      string `json:"waived"`
	Cancelled   string `json:"cancelled"`
}

// SupplierBillingEventTriggerLabels — labels for the trigger enum.
type SupplierBillingEventTriggerLabels struct {
	Unspecified string `json:"unspecified"`
	ManualEarly string `json:"manualEarly"`
	ManualLate  string `json:"manualLate"`
}

// SupplierBillingEventColumnLabels — list table column headers.
type SupplierBillingEventColumnLabels struct {
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

// SupplierBillingEventPageLabels — list / detail page strings.
type SupplierBillingEventPageLabels struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
}

// SupplierBillingEventActionLabels — row-action button labels.
type SupplierBillingEventActionLabels struct {
	Recognize string `json:"recognize"`
	MarkReady string `json:"markReady"`
	Waive     string `json:"waive"`
	Cancel    string `json:"cancel"`
}

// SupplierBillingEventDetailLabels — detail page tab + section labels.
type SupplierBillingEventDetailLabels struct {
	Title              string `json:"title"`
	TabInfo            string `json:"tabInfo"`
	TabAudit           string `json:"tabAudit"`
	InfoHeading        string `json:"infoHeading"`
	LinkedAdvanceBadge string `json:"linkedAdvanceBadge"`
}

// SupplierBillingEventEmptyLabels — empty-state labels for the list view.
type SupplierBillingEventEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// SupplierBillingEventErrorLabels — error toasts / validations.
type SupplierBillingEventErrorLabels struct {
	PermissionDenied  string `json:"permissionDenied"`
	NotFound          string `json:"notFound"`
	AlreadyRecognized string `json:"alreadyRecognized"`
	InvalidTransition string `json:"invalidTransition"`
}

// SupplierBillingEventLabels — root struct for supplier_billing_event.json.
type SupplierBillingEventLabels struct {
	Page    SupplierBillingEventPageLabels    `json:"page"`
	Columns SupplierBillingEventColumnLabels  `json:"columns"`
	Status  SupplierBillingEventStatusLabels  `json:"status"`
	Trigger SupplierBillingEventTriggerLabels `json:"trigger"`
	Actions SupplierBillingEventActionLabels  `json:"actions"`
	Detail  SupplierBillingEventDetailLabels  `json:"detail"`
	Empty   SupplierBillingEventEmptyLabels   `json:"empty"`
	Errors  SupplierBillingEventErrorLabels   `json:"errors"`
}

// === Defaults ===

// DefaultSupplierBillingEventLabels — English defaults for the supplier
// billing-event list / detail screens that anchor MILESTONE recognition on
// the buying side.
func DefaultSupplierBillingEventLabels() SupplierBillingEventLabels {
	return SupplierBillingEventLabels{
		Page: SupplierBillingEventPageLabels{
			Title:   "Supplier Billing Events",
			Caption: "Manual billing triggers for supplier subscriptions and contracts.",
		},
		Columns: SupplierBillingEventColumnLabels{
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
		Status: SupplierBillingEventStatusLabels{
			Unspecified: "Unspecified",
			Ready:       "Ready",
			Billed:      "Billed",
			Waived:      "Waived",
			Cancelled:   "Cancelled",
		},
		Trigger: SupplierBillingEventTriggerLabels{
			Unspecified: "Unspecified",
			ManualEarly: "Manual (early)",
			ManualLate:  "Manual (late)",
		},
		Actions: SupplierBillingEventActionLabels{
			Recognize: "Recognize",
			MarkReady: "Mark ready",
			Waive:     "Waive",
			Cancel:    "Cancel",
		},
		Detail: SupplierBillingEventDetailLabels{
			Title:              "Supplier billing event",
			TabInfo:            "Info",
			TabAudit:           "Audit",
			InfoHeading:        "Billing event details",
			LinkedAdvanceBadge: "Linked advance",
		},
		Empty: SupplierBillingEventEmptyLabels{
			Title:   "No supplier billing events",
			Message: "Supplier billing events appear here once created.",
		},
		Errors: SupplierBillingEventErrorLabels{
			PermissionDenied:  "Permission denied",
			NotFound:          "Supplier billing event not found",
			AlreadyRecognized: "This billing event has already been recognized",
			InvalidTransition: "Cannot transition from %s to %s",
		},
	}
}
