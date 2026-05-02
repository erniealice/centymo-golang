// Package form owns the template data shape for the primary revenue drawer
// (revenue-drawer-form.html). Pure types only — no Deps, no context.Context,
// no repository imports.
package form

// PaymentTermOption is a minimal struct for rendering payment term options in the form.
type PaymentTermOption struct {
	Id      string
	Name    string
	NetDays int32
}

// Inner holds nested form labels accessed via .Labels.Form.* in templates.
type Inner struct {
	SectionInfo               string
	CurrencyPlaceholder       string
	StatusDraft               string
	StatusComplete            string
	StatusCancelled           string
	CustomerNamePlaceholder   string
	CustomerSearchPlaceholder string
	CustomerNoResults         string
	LocationPlaceholder       string
	LocationSearchPlaceholder string
	LocationNoResults         string
}

// Labels holds i18n labels for the drawer form template.
type Labels struct {
	Customer                  string
	Date                      string
	Currency                  string
	Reference                 string
	ReferencePlaceholder      string
	Status                    string
	Notes                     string
	NotesPlaceholder          string
	Location                  string
	PaymentTerms              string
	SelectPaymentTerm         string
	DueDate                   string
	Subscription              string
	SubscriptionNoResults     string
	RevenueType               string
	RevenueTypeOneTime        string
	RevenueTypeFromEngagement string
	RevenueTypeFromActivities string
	ActivityIDs               string
	ActivityIDsPlaceholder    string
	Form                      Inner

	// Field-level info text surfaced via an info button beside each label.
	ReferenceInfo    string
	DateInfo         string
	CustomerInfo     string
	LocationInfo     string
	SubscriptionInfo string
	CurrencyInfo     string
	NotesInfo        string
}

// Data is the template data for the revenue drawer form.
type Data struct {
	FormAction            string
	IsEdit                bool
	ID                    string
	Name                  string
	ClientID              string
	ClientLabel           string
	SearchClientURL       string
	SubscriptionID        string
	SubscriptionLabel     string
	SearchSubscriptionURL string
	ReferenceNumber       string
	Date                  string
	Currency              string
	Status                string
	Notes                 string
	LocationID            string
	LocationLabel         string
	SearchLocationURL     string
	PaymentTerms          []*PaymentTermOption
	SelectedPaymentTermID string
	DueDateString         string
	RevenueType           string
	ActivityIDs           string
	Labels                Labels
	CommonLabels          any
}
