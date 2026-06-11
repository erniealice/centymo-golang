package price_list

// ---------------------------------------------------------------------------
// Price List labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the price list module.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Buttons ButtonLabels  `json:"buttons"`
	Columns ColumnLabels  `json:"columns"`
	Empty   EmptyLabels   `json:"empty"`
	Form    FormLabels    `json:"form"`
	Actions ActionLabels  `json:"actions"`
	Bulk    BulkLabels    `json:"bulkActions"`
	Detail  DetailLabels  `json:"detail"`
	Confirm ConfirmLabels `json:"confirm"`
	Errors  ErrorLabels   `json:"errors"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type ButtonLabels struct {
	AddPriceList string `json:"addPriceList"`
}

type ColumnLabels struct {
	Name      string `json:"name"`
	DateStart string `json:"dateStart"`
	DateEnd   string `json:"dateEnd"`
	Status    string `json:"status"`
}

type EmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type FormLabels struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DescPlaceholder string `json:"descriptionPlaceholder"`
	DateStart       string `json:"dateStart"`
	DateEnd         string `json:"dateEnd"`
	Active          string `json:"active"`
	Product         string `json:"product"`
	SelectProduct   string `json:"selectProduct"`
	Amount          string `json:"amount"`
	Currency        string `json:"currency"`

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string `json:"nameInfo"`
	DescriptionInfo string `json:"descriptionInfo"`
	DateStartInfo   string `json:"dateStartInfo"`
	DateEndInfo     string `json:"dateEndInfo"`
	ActiveInfo      string `json:"activeInfo"`
	// Price-product sub-drawer info fields.
	AmountInfo   string `json:"amountInfo"`
	CurrencyInfo string `json:"currencyInfo"`
}

type ActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

type BulkLabels struct {
	Delete string `json:"delete"`
}

type DetailLabels struct {
	PageTitle          string `json:"pageTitle"`
	BasicInfo          string `json:"basicInfo"`
	Prices             string `json:"prices"`
	TabAttachments     string `json:"tabAttachments"`
	TabAuditHistory    string `json:"tabAuditHistory"`
	ProductName        string `json:"productName"`
	Amount             string `json:"amount"`
	Currency           string `json:"currency"`
	AddPrice           string `json:"addPrice"`
	RemoveLabel        string `json:"removeLabel"`
	EmptyTitle         string `json:"emptyTitle"`
	EmptyMessage       string `json:"emptyMessage"`
	ActiveBadge        string `json:"activeBadge"`
	InactiveBadge      string `json:"inactiveBadge"`
	NoPricesConfigured string `json:"noPricesConfigured"`
	NoPricesDesc       string `json:"noPricesDesc"`
}

type ConfirmLabels struct {
	Activate          string `json:"activate"`
	ActivateMessage   string `json:"activateMessage"`
	Deactivate        string `json:"deactivate"`
	DeactivateMessage string `json:"deactivateMessage"`
	Delete            string `json:"delete"`
	DeleteMessage     string `json:"deleteMessage"`
	BulkDelete        string `json:"bulkDelete"`
	BulkDeleteMessage string `json:"bulkDeleteMessage"`
}

type ErrorLabels struct {
	PermissionDenied string `json:"permissionDenied"`
	InvalidFormData  string `json:"invalidFormData"`
	NotFound         string `json:"notFound"`
	IDRequired       string `json:"idRequired"`
	NoIDsProvided    string `json:"noIDsProvided"`
	CannotDelete     string `json:"cannotDelete"`
	ProductRequired  string `json:"productRequired"`
	AmountRequired   string `json:"amountRequired"`
}
