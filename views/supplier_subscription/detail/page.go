package detail

import (
	"context"
	"log"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
	suppliersubscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_subscription"
)

// DetailViewDeps holds view dependencies for the supplier_subscription detail page.
type DetailViewDeps struct {
	Routes       centymo.SupplierSubscriptionRoutes
	Labels       centymo.SupplierSubscriptionLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ReadSupplierSubscription            func(ctx context.Context, req *suppliersubscriptionpb.ReadSupplierSubscriptionRequest) (*suppliersubscriptionpb.ReadSupplierSubscriptionResponse, error)
	GetSupplierSubscriptionItemPageData func(ctx context.Context, req *suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataRequest) (*suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataResponse, error)

	// ListExpenseRecognitions is used to populate the Linked Recognitions tab.
	// Filtered client-side on supplier_subscription_id = $id since the proto
	// request has no supplier_subscription_id filter yet (proto field 60 FK was
	// added in Wave 1; request filter is a P3-followup item in the plan).
	// Nil-safe — tab renders empty state when unset.
	ListExpenseRecognitions func(ctx context.Context, req *expenserecognitionpb.ListExpenseRecognitionsRequest) (*expenserecognitionpb.ListExpenseRecognitionsResponse, error)

	// ExpenseRecognitionDetailURL is the path template (e.g.
	// "/app/expense-recognitions/detail/{id}") used to deep-link from the
	// Recognitions tab rows. Empty disables row click-through.
	ExpenseRecognitionDetailURL string
}

// RecognitionRow is a single row in the Linked Recognitions tab table.
type RecognitionRow struct {
	ID              string
	Name            string
	Status          string
	StatusVariant   string
	RecognitionDate string
	TotalAmount     types.TableCell
	DetailURL       string
}

// PageData holds the data for the supplier_subscription detail page.
type PageData struct {
	types.PageData
	ActiveTab string
	TabItems  []pyeza.TabItem
	Record    *suppliersubscriptionpb.SupplierSubscription
	EditURL   string
	DeleteURL string

	// RecognizeExpenseURL is the POST endpoint for the Recognize Expense CTA.
	// Empty = CTA hidden.
	RecognizeExpenseURL string

	// RecognitionsTabRefreshURL is the HTMX endpoint for the Linked Recognitions
	// tab to self-refresh when an expense-recognitions-table event fires (e.g.
	// after the Recognize Expense CTA succeeds). Built as TabActionURL + "recognitions".
	RecognitionsTabRefreshURL string

	// Recognitions holds the rows for the Linked Recognitions tab.
	Recognitions []RecognitionRow
}

func loadRecord(ctx context.Context, deps *DetailViewDeps, id string) (*suppliersubscriptionpb.SupplierSubscription, error) {
	if deps.GetSupplierSubscriptionItemPageData != nil {
		resp, err := deps.GetSupplierSubscriptionItemPageData(ctx, &suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataRequest{
			SupplierSubscriptionId: id,
		})
		if err != nil || resp == nil || resp.GetSupplierSubscription() == nil {
			return nil, err
		}
		return resp.GetSupplierSubscription(), nil
	}
	resp, err := deps.ReadSupplierSubscription(ctx, &suppliersubscriptionpb.ReadSupplierSubscriptionRequest{
		Data: &suppliersubscriptionpb.SupplierSubscription{Id: id},
	})
	if err != nil || len(resp.GetData()) == 0 {
		return nil, err
	}
	return resp.GetData()[0], nil
}

// NewView creates the supplier_subscription detail page view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		record, err := loadRecord(ctx, deps, id)
		if err != nil || record == nil {
			log.Printf("Failed to load supplier subscription detail %s: %v", id, err)
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}

		l := deps.Labels
		tabItems := buildTabs(l, deps.Routes, id, activeTab)

		pageTitle := record.GetName()
		if pageTitle == "" {
			pageTitle = record.GetCode()
		}

		recognizeURL := ""
		if deps.Routes.RecognizeExpenseURL != "" {
			recognizeURL = route.ResolveURL(deps.Routes.RecognizeExpenseURL, "id", id)
		}
		recognitionsRefreshURL := ""
		if deps.Routes.TabActionURL != "" {
			recognitionsRefreshURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "recognitions")
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          pageTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    pageTitle,
				HeaderSubtitle: l.Detail.InfoSection,
				HeaderIcon:     "icon-refresh-cw",
				CommonLabels:   deps.CommonLabels,
			},
			ActiveTab:                 activeTab,
			TabItems:                  tabItems,
			Record:                    record,
			EditURL:                   route.ResolveURL(deps.Routes.EditURL, "id", id),
			DeleteURL:                 deps.Routes.DeleteURL,
			RecognizeExpenseURL:       recognizeURL,
			RecognitionsTabRefreshURL: recognitionsRefreshURL,
		}

		if activeTab == "recognitions" {
			pageData.Recognitions = loadLinkedRecognitions(ctx, deps, id)
		}

		return view.OK("supplier-subscription-detail", pageData)
	})
}

// NewTabAction handles HTMX tab-swap requests for the detail page.
func NewTabAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		record, err := loadRecord(ctx, deps, id)
		if err != nil || record == nil {
			return centymo.HTMXError(deps.Labels.Errors.NotFound)
		}

		l := deps.Labels

		recognizeURL := ""
		if deps.Routes.RecognizeExpenseURL != "" {
			recognizeURL = route.ResolveURL(deps.Routes.RecognizeExpenseURL, "id", id)
		}
		recognitionsRefreshURL := ""
		if deps.Routes.TabActionURL != "" {
			recognitionsRefreshURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "recognitions")
		}

		pageData := &PageData{
			ActiveTab:                 tab,
			TabItems:                  buildTabs(l, deps.Routes, id, tab),
			Record:                    record,
			EditURL:                   route.ResolveURL(deps.Routes.EditURL, "id", id),
			DeleteURL:                 deps.Routes.DeleteURL,
			RecognizeExpenseURL:       recognizeURL,
			RecognitionsTabRefreshURL: recognitionsRefreshURL,
		}

		if tab == "recognitions" {
			pageData.Recognitions = loadLinkedRecognitions(ctx, deps, id)
		}

		return view.OK("supplier-subscription-detail-tab-"+tab, pageData)
	})
}

func buildTabs(l centymo.SupplierSubscriptionLabels, routes centymo.SupplierSubscriptionRoutes, id string, activeTab string) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info"},
		{Key: "cost_plan", Label: l.Tabs.CostPlan, Href: base + "?tab=cost_plan", HxGet: action + "cost_plan"},
		{Key: "recognitions", Label: l.Tabs.LinkedRecognitions, Href: base + "?tab=recognitions", HxGet: action + "recognitions"},
		{Key: "activity", Label: l.Tabs.Activity, Href: base + "?tab=activity", HxGet: action + "activity"},
	}
}

// loadLinkedRecognitions fetches expense_recognition rows where
// supplier_subscription_id = supplierSubscriptionID.
// The ListExpenseRecognitionsRequest has no supplier_subscription_id filter yet
// (proto request field is a follow-up item); we post-filter in the view layer
// using GetSupplierSubscriptionId() on each returned row.
func loadLinkedRecognitions(ctx context.Context, deps *DetailViewDeps, supplierSubscriptionID string) []RecognitionRow {
	if deps.ListExpenseRecognitions == nil {
		return nil
	}
	resp, err := deps.ListExpenseRecognitions(ctx, &expenserecognitionpb.ListExpenseRecognitionsRequest{})
	if err != nil {
		log.Printf("loadLinkedRecognitions for ss=%s: %v", supplierSubscriptionID, err)
		return nil
	}
	var rows []RecognitionRow
	for _, r := range resp.GetData() {
		if r.GetSupplierSubscriptionId() != supplierSubscriptionID {
			continue
		}
		currency := r.GetCurrency()
		recDate := ""
		if r.RecognitionDate != nil {
			recDate = r.GetRecognitionDate().AsTime().Format("2006-01-02")
		}
		detailURL := ""
		if deps.ExpenseRecognitionDetailURL != "" {
			detailURL = route.ResolveURL(deps.ExpenseRecognitionDetailURL, "id", r.GetId())
		}
		rows = append(rows, RecognitionRow{
			ID:              r.GetId(),
			Name:            r.GetName(),
			Status:          r.GetStatus().String(),
			StatusVariant:   recognitionStatusVariant(r.GetStatus().String()),
			RecognitionDate: recDate,
			TotalAmount:     types.MoneyCell(float64(r.GetTotalAmount()), currency, true),
			DetailURL:       detailURL,
		})
	}
	return rows
}

func recognitionStatusVariant(status string) string {
	switch status {
	case "EXPENSE_RECOGNITION_STATUS_DRAFT":
		return "default"
	case "EXPENSE_RECOGNITION_STATUS_POSTED":
		return "success"
	case "EXPENSE_RECOGNITION_STATUS_REVERSED":
		return "danger"
	default:
		return "default"
	}
}

