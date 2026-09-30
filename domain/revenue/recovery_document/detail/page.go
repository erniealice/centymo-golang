// Package detail renders the recovery document detail page
// (info · lines · payments applied · credit notes).
package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	recoverydocumentlinepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document_line"
	collectionapplicationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_application"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
	ra "github.com/erniealice/centymo-golang/domain/treasury/collection_application"
)

// Deps holds the detail view dependencies.
type Deps struct {
	Routes       rd.Routes
	Labels       rd.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *rd.UseCases

	// ApplicationLabels are the collection_application labels (applications
	// tab columns, reverse confirmation). ReverseURL and ReceiveApplyURL are
	// the sibling unit's route templates; empty hides the action.
	ApplicationLabels ra.Labels
	ReverseURL        string
	ReceiveApplyURL   string

	// ApplicationsTableID is the applications table id (owned by the
	// collection_application unit, passed in so this view imports no sibling view).
	ApplicationsTableID string
}

// PageData is the detail page view model.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          rd.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID             string
	Number         string
	TypeLabel      string
	StatusLabel    string
	StatusVariant  string
	IsVoid         bool
	Client         string
	IssueDate      string
	DueDate        string
	Total          string
	Applied        string
	Balance        string
	BalanceLabel   string
	VoidReason     string
	CorrectsNumber string
	CorrectsURL    string

	VoidURL          string
	CanVoid          bool
	VoidDisabledTip  string
	ApplyPaymentURL  string
	ApplyPaymentText string
	CanApplyPayment  bool
	MissingPermApply string

	LinesTable        *types.TableConfig
	ApplicationsTable *types.TableConfig
	CreditNotesTable  *types.TableConfig
}

// NewView renders the full detail page.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("recovery_document", "read") {
			return view.Forbidden("recovery_document:read")
		}
		id := viewCtx.Request.PathValue("id")
		tab := normalizeTab(viewCtx.Request.URL.Query().Get("tab"))
		pd, status, err := buildPageData(ctx, deps, viewCtx, id, tab, perms)
		if err != nil {
			return failure(status, err)
		}
		return view.OK("recovery-document-detail", pd)
	})
}

// NewTabAction serves the HTMX tab partials.
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("recovery_document", "read") {
			return view.Forbidden("recovery_document:read")
		}
		id := viewCtx.Request.PathValue("id")
		tab := normalizeTab(viewCtx.Request.PathValue("tab"))
		pd, status, err := buildPageData(ctx, deps, viewCtx, id, tab, perms)
		if err != nil {
			return failure(status, err)
		}
		return view.OK("recovery-document-tab-"+tab, pd)
	})
}

func failure(status int, err error) view.ViewResult {
	r := view.Error(err)
	if status != 0 {
		r.StatusCode = status
	}
	return r
}

func normalizeTab(t string) string {
	switch t {
	case "info", "lines", "applications", "credit-notes":
		return t
	default:
		return "info"
	}
}

func typeLabel(l rd.Labels, t recoverydocumentpb.RecoveryDocumentType) string {
	if t == recoverydocumentpb.RecoveryDocumentType_RECOVERY_DOCUMENT_TYPE_CREDIT_NOTE {
		return l.Enums.TypeCreditNote
	}
	return l.Enums.TypeStatement
}

func fmtMillis(ms int64) string {
	if ms == 0 {
		return ""
	}
	return time.UnixMilli(ms).UTC().Format("2006-01-02")
}

// appliedTotal sums the live cash applications (a reversed original and its
// reversal row both carry status REVERSED, so only APPLIED rows count; an inactive
// row never counts, matching recovery_reporting's snapshot filter — A1 m6).
func appliedTotal(apps []*collectionapplicationpb.CollectionApplication) int64 {
	var sum int64
	for _, a := range apps {
		if a.GetActive() && a.GetStatus() == collectionapplicationpb.ApplicationStatus_APPLICATION_STATUS_APPLIED &&
			a.GetApplicationKind() == collectionapplicationpb.ApplicationKind_APPLICATION_KIND_CASH {
			sum += a.GetAmount()
		}
	}
	return sum
}

// creditedTotal sums the live ISSUED credit notes' totals (negative amounts; a
// void note no longer corrects the statement).
func creditedTotal(notes []*recoverydocumentpb.RecoveryDocument) int64 {
	var sum int64
	for _, n := range notes {
		if n.GetActive() && n.GetStatus() == recoverydocumentpb.RecoveryDocumentStatus_RECOVERY_DOCUMENT_STATUS_ISSUED {
			sum += n.GetTotalAmount()
		}
	}
	return sum
}

func buildPageData(ctx context.Context, deps *Deps, viewCtx *view.ViewContext, id, tab string, perms *types.UserPermissions) (*PageData, int, error) {
	l := deps.Labels
	uc := deps.UseCases
	if uc == nil || uc.ReadRecoveryDocument == nil {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("%s", deps.CommonLabels.Errors.General)
	}
	resp, err := uc.ReadRecoveryDocument(ctx, &recoverydocumentpb.ReadRecoveryDocumentRequest{
		Data: &recoverydocumentpb.RecoveryDocument{Id: id},
	})
	if err != nil {
		if rd.ErrorKind(err) == rd.ErrKindNotFound {
			return nil, http.StatusNotFound, fmt.Errorf("%s", l.Errors.NotFound)
		}
		log.Printf("recovery_document detail: read %s: %v", id, err)
		return nil, http.StatusInternalServerError, err
	}
	if len(resp.GetData()) == 0 {
		return nil, http.StatusNotFound, fmt.Errorf("%s", l.Errors.NotFound)
	}
	doc := resp.GetData()[0]

	var apps []*collectionapplicationpb.CollectionApplication
	if uc.ListCollectionApplications != nil {
		if apps, err = uc.ApplicationsOf(ctx, id); err != nil {
			log.Printf("recovery_document detail: applications of %s: %v", id, err)
			return nil, http.StatusInternalServerError, err
		}
	}

	cl := deps.CommonLabels
	client := uc.ClientName(ctx, doc.GetClientId())
	if client == "" {
		client = doc.GetClientId()
	}
	isVoid := doc.GetStatus() == recoverydocumentpb.RecoveryDocumentStatus_RECOVERY_DOCUMENT_STATUS_VOID
	applied := appliedTotal(apps)
	// Same snapshot rule as preview / aging / GetRecoveryBalance (R5 M5): total +
	// credited - applied, where credit notes carry a negative total.
	balance := doc.GetTotalAmount() + creditedTotal(resp.GetCreditNotes()) - applied
	pd := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          doc.GetDocumentNumber(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    doc.GetDocumentNumber(),
			HeaderSubtitle: typeLabel(l, doc.GetDocumentType()),
			HeaderIcon:     "icon-file-text",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "recovery-document-detail-content",
		Labels:          l,
		ActiveTab:       tab,
		ID:              id,
		Number:          doc.GetDocumentNumber(),
		TypeLabel:       typeLabel(l, doc.GetDocumentType()),
		IsVoid:          isVoid,
		Client:          client,
		IssueDate:       doc.GetIssueDate(),
		DueDate:         doc.GetDueDate(),
		Total:           types.FormatMoney(doc.GetTotalAmount(), doc.GetCurrency()),
		Applied:         types.FormatMoney(applied, doc.GetCurrency()),
		Balance:         types.FormatMoney(balance, doc.GetCurrency()),
		BalanceLabel:    l.Detail.BalanceDue,
		VoidReason:      doc.GetVoidReason(),
		VoidURL:         route.ResolveURL(deps.Routes.VoidURL, "id", id),
		CanVoid:         perms.Can("recovery_document", "void") && !isVoid,
	}
	if balance < 0 {
		pd.BalanceLabel = l.Detail.CreditBalance
	}
	if isVoid {
		pd.StatusLabel, pd.StatusVariant = l.Enums.StatusVoid, "muted"
		pd.VoidDisabledTip = l.Errors.AlreadyVoid
	} else {
		pd.StatusLabel, pd.StatusVariant = l.Enums.StatusIssued, "success"
		if !perms.Can("recovery_document", "void") {
			pd.VoidDisabledTip = fmt.Sprintf(cl.Errors.MissingPermission, "recovery_document:void")
		}
	}
	if doc.GetCorrectsDocumentId() != "" {
		pd.CorrectsNumber = doc.GetCorrectsDocumentId()
		// Show the corrected document's number (falls back to its id).
		if r, rerr := uc.ReadRecoveryDocument(ctx, &recoverydocumentpb.ReadRecoveryDocumentRequest{
			Data: &recoverydocumentpb.RecoveryDocument{Id: doc.GetCorrectsDocumentId()},
		}); rerr == nil && len(r.GetData()) > 0 && r.GetData()[0].GetDocumentNumber() != "" {
			pd.CorrectsNumber = r.GetData()[0].GetDocumentNumber()
		}
		pd.CorrectsURL = route.ResolveURL(deps.Routes.DetailURL, "id", doc.GetCorrectsDocumentId())
	}
	if deps.ReceiveApplyURL != "" && !isVoid && doc.GetDocumentType() == recoverydocumentpb.RecoveryDocumentType_RECOVERY_DOCUMENT_TYPE_STATEMENT {
		pd.ApplyPaymentURL = deps.ReceiveApplyURL + "?client_id=" + doc.GetClientId()
		pd.ApplyPaymentText = l.Buttons.ApplyPayment
		pd.CanApplyPayment = perms.Can("collection_application", "create")
		pd.MissingPermApply = fmt.Sprintf(cl.Errors.MissingPermission, "collection_application:create")
	}

	pd.TabItems = tabItems(deps, id, len(resp.GetLines()), len(apps), len(resp.GetCreditNotes()))
	switch tab {
	case "lines":
		pd.LinesTable = linesTable(deps, doc, resp.GetLines())
	case "applications":
		pd.ApplicationsTable = applicationsTable(deps, doc, apps, perms)
	case "credit-notes":
		pd.CreditNotesTable = creditNotesTable(deps, resp.GetCreditNotes())
	}
	return pd, 0, nil
}

func tabItems(deps *Deps, id string, lines, apps, notes int) []pyeza.TabItem {
	l := deps.Labels
	base := route.ResolveURL(deps.Routes.DetailURL, "id", id)
	action := route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "lines", Label: l.Tabs.Lines, Href: base + "?tab=lines", HxGet: action + "lines", Icon: "icon-list", Count: lines},
		{Key: "applications", Label: l.Tabs.Applications, Href: base + "?tab=applications", HxGet: action + "applications", Icon: "icon-credit-card", Count: apps},
		{Key: "credit-notes", Label: l.Tabs.CreditNotes, Href: base + "?tab=credit-notes", HxGet: action + "credit-notes", Icon: "icon-file-minus", Count: notes},
	}
}

func tabRefresh(deps *Deps, id, tab string) string {
	return route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", tab)
}

func linesTable(deps *Deps, doc *recoverydocumentpb.RecoveryDocument, lines []*recoverydocumentlinepb.RecoveryDocumentLine) *types.TableConfig {
	l := deps.Labels
	cols := []types.TableColumn{
		{Key: "description", Label: l.Columns.Description, NoSort: true, NoFilter: true},
		{Key: "service_period", Label: l.Columns.ServicePeriod, NoSort: true, NoFilter: true, WidthClass: "col-4xl"},
		{Key: "amount", Label: l.Columns.Amount, NoSort: true, NoFilter: true, WidthClass: "col-3xl", Align: "right"},
	}
	rows := make([]types.TableRow, 0, len(lines))
	for _, ln := range lines {
		period := ln.GetServiceFrom()
		if ln.GetServiceTo() != "" {
			period += " – " + ln.GetServiceTo()
		}
		rows = append(rows, types.TableRow{
			ID: ln.GetId(),
			Cells: []types.TableCell{
				{Type: "text", Value: ln.GetDescription()},
				{Type: "text", Value: period},
				{Type: "text", Value: types.FormatMoney(ln.GetAmount(), ln.GetCurrency())},
			},
			DataAttrs: map[string]string{"testid": "recovery-document-line-" + ln.GetId()},
		})
	}
	types.ApplyColumnStyles(cols, rows)
	t := &types.TableConfig{
		ID: "recovery-document-lines-table", RefreshURL: tabRefresh(deps, doc.GetId(), "lines"),
		Columns: cols, Rows: rows, ShowSearch: false, ShowActions: false, ShowSort: false, ShowEntries: false,
		Labels:     deps.TableLabels,
		EmptyState: types.TableEmptyState{Title: l.Tabs.Lines},
	}
	types.ApplyTableSettings(t)
	return t
}

func applicationsTable(deps *Deps, doc *recoverydocumentpb.RecoveryDocument, apps []*collectionapplicationpb.CollectionApplication, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	al := deps.ApplicationLabels
	cols := []types.TableColumn{
		{Key: "applied_on", Label: al.Columns.AppliedOn, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "amount", Label: al.Columns.Amount, NoSort: true, NoFilter: true, WidthClass: "col-3xl", Align: "right"},
		{Key: "kind", Label: al.Columns.Kind, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "status", Label: al.Columns.Status, NoSort: true, NoFilter: true, WidthClass: "col-2xl"},
	}
	rows := make([]types.TableRow, 0, len(apps))
	for _, a := range apps {
		id := a.GetId()
		kind := al.Enums.ApplicationKindCash
		if a.GetApplicationKind() == collectionapplicationpb.ApplicationKind_APPLICATION_KIND_NON_CASH_SETTLEMENT {
			kind = al.Enums.ApplicationKindNonCashSettlement
		}
		stLabel, stVariant := al.Enums.StatusApplied, "success"
		if a.GetStatus() == collectionapplicationpb.ApplicationStatus_APPLICATION_STATUS_REVERSED {
			stLabel, stVariant = al.Enums.StatusReversed, "muted"
		}
		var actions []types.TableAction
		reversible := a.GetStatus() == collectionapplicationpb.ApplicationStatus_APPLICATION_STATUS_APPLIED && a.GetReversesApplicationId() == ""
		if reversible && deps.ReverseURL != "" {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: al.Buttons.Reverse, Action: "deactivate", TestID: "collection-application-reverse-" + id,
				URL: route.ResolveURL(deps.ReverseURL, "id", id), ItemName: doc.GetDocumentNumber(),
				ConfirmTitle: al.Confirm.ReverseTitle, ConfirmMessage: al.Confirm.ReverseMsg,
				Disabled:        !perms.Can("collection_application", "reverse"),
				DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_application:reverse"),
			})
		}
		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: fmtMillis(a.GetAppliedAt())},
				{Type: "text", Value: types.FormatMoney(a.GetAmount(), a.GetCurrency())},
				{Type: "text", Value: kind},
				{Type: "badge", Value: stLabel, Variant: stVariant},
			},
			DataAttrs: map[string]string{"testid": "collection-application-row-" + id, "status": a.GetStatus().String()},
			Actions:   actions,
		})
	}
	types.ApplyColumnStyles(cols, rows)
	t := &types.TableConfig{
		ID: deps.ApplicationsTableID, RefreshURL: tabRefresh(deps, doc.GetId(), "applications"),
		Columns: cols, Rows: rows, ShowActions: true, ShowSort: false, ShowEntries: false,
		Labels:     deps.TableLabels,
		EmptyState: types.TableEmptyState{Title: l.Empty.ApplicationsTitle},
	}
	types.ApplyTableSettings(t)
	return t
}

func creditNotesTable(deps *Deps, notes []*recoverydocumentpb.RecoveryDocument) *types.TableConfig {
	l := deps.Labels
	cols := []types.TableColumn{
		{Key: "document_number", Label: l.Columns.DocumentNumber, NoSort: true, NoFilter: true, WidthClass: "col-4xl"},
		{Key: "issue_date", Label: l.Columns.IssueDate, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "total", Label: l.Columns.Total, NoSort: true, NoFilter: true, WidthClass: "col-3xl", Align: "right"},
		{Key: "status", Label: l.Columns.Status, NoSort: true, NoFilter: true, WidthClass: "col-2xl"},
	}
	rows := make([]types.TableRow, 0, len(notes))
	for _, n := range notes {
		id := n.GetId()
		href := route.ResolveURL(deps.Routes.DetailURL, "id", id)
		stLabel, stVariant := l.Enums.StatusIssued, "success"
		if n.GetStatus() == recoverydocumentpb.RecoveryDocumentStatus_RECOVERY_DOCUMENT_STATUS_VOID {
			stLabel, stVariant = l.Enums.StatusVoid, "muted"
		}
		rows = append(rows, types.TableRow{
			ID: id, Href: href,
			Cells: []types.TableCell{
				{Type: "link", Value: n.GetDocumentNumber(), Href: href},
				{Type: "text", Value: n.GetIssueDate()},
				{Type: "text", Value: types.FormatMoney(n.GetTotalAmount(), n.GetCurrency())},
				{Type: "badge", Value: stLabel, Variant: stVariant},
			},
			DataAttrs: map[string]string{"testid": "recovery-document-credit-note-" + id},
			Actions:   []types.TableAction{{Type: "view", Label: l.Buttons.ViewCreditNote, Action: "view", Href: href, TestID: "recovery-document-credit-note-view-" + id}},
		})
	}
	types.ApplyColumnStyles(cols, rows)
	t := &types.TableConfig{
		ID: "recovery-document-credit-notes-table", RefreshURL: "",
		Columns: cols, Rows: rows, ShowActions: true, ShowSort: false, ShowEntries: false,
		Labels:     deps.TableLabels,
		EmptyState: types.TableEmptyState{Title: l.Empty.CreditNotesTitle},
	}
	types.ApplyTableSettings(t)
	return t
}
