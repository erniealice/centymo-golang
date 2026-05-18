package supplier_subscription

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	suppliersubscriptionaction "github.com/erniealice/centymo-golang/views/supplier_subscription/action"
	suppliersubscriptiondetail "github.com/erniealice/centymo-golang/views/supplier_subscription/detail"
	suppliersubscriptionlist "github.com/erniealice/centymo-golang/views/supplier_subscription/list"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
	costplanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/cost_plan"
	suppliersubscriptionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/procurement/supplier_subscription"
)

// ModuleDeps holds all dependencies for the supplier_subscription module.
type ModuleDeps struct {
	Routes       centymo.SupplierSubscriptionRoutes
	Labels       centymo.SupplierSubscriptionLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ExpenseRecognitionRunLabels supplies the "Run Recognitions" label for the
	// Linked Recognitions tab CTA when CostPlan.billing_kind is
	// RECURRING / CONTRACT-with-cycle. Plan A Surface C 20260517-expense-run.
	ExpenseRecognitionRunLabels centymo.ExpenseRecognitionRunLabels

	// ReadCostPlan resolves the supplier subscription's CostPlan so the detail
	// page can branch the Recognitions CTA on billing_kind. Nil-safe.
	ReadCostPlan func(ctx context.Context, req *costplanpb.ReadCostPlanRequest) (*costplanpb.ReadCostPlanResponse, error)

	CreateSupplierSubscription          func(ctx context.Context, req *suppliersubscriptionpb.CreateSupplierSubscriptionRequest) (*suppliersubscriptionpb.CreateSupplierSubscriptionResponse, error)
	ReadSupplierSubscription            func(ctx context.Context, req *suppliersubscriptionpb.ReadSupplierSubscriptionRequest) (*suppliersubscriptionpb.ReadSupplierSubscriptionResponse, error)
	UpdateSupplierSubscription          func(ctx context.Context, req *suppliersubscriptionpb.UpdateSupplierSubscriptionRequest) (*suppliersubscriptionpb.UpdateSupplierSubscriptionResponse, error)
	DeleteSupplierSubscription          func(ctx context.Context, req *suppliersubscriptionpb.DeleteSupplierSubscriptionRequest) (*suppliersubscriptionpb.DeleteSupplierSubscriptionResponse, error)
	GetSupplierSubscriptionListPageData func(ctx context.Context, req *suppliersubscriptionpb.GetSupplierSubscriptionListPageDataRequest) (*suppliersubscriptionpb.GetSupplierSubscriptionListPageDataResponse, error)
	GetSupplierSubscriptionItemPageData func(ctx context.Context, req *suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataRequest) (*suppliersubscriptionpb.GetSupplierSubscriptionItemPageDataResponse, error)

	// SetSupplierSubscriptionActive performs a raw DB update to toggle active.
	// Required because proto3 omits bool=false on serialization.
	SetSupplierSubscriptionActive func(ctx context.Context, id string, active bool) error

	// ListExpenseRecognitions is used to populate the "Linked Recognitions" tab on
	// the supplier_subscription detail page. Post-filtered by supplier_subscription_id.
	// Nil-safe — tab renders empty state when unset.
	ListExpenseRecognitions func(ctx context.Context, req *expenserecognitionpb.ListExpenseRecognitionsRequest) (*expenserecognitionpb.ListExpenseRecognitionsResponse, error)

	// ExpenseRecognitionDetailURL is the path template for deep-linking to an
	// expense_recognition detail page from the Linked Recognitions tab.
	// Empty string disables row click-through (rows still render, just not linked).
	ExpenseRecognitionDetailURL string

	// RecognizeFromExpenditure is the espyna use case called by the
	// "Recognize Expense" CTA on the detail page.
	// Nil-safe — CTA is hidden when this dep is unset.
	RecognizeFromExpenditure func(ctx context.Context, req *expenserecognitionpb.RecognizeFromExpenditureRequest) (*expenserecognitionpb.RecognizeFromExpenditureResponse, error)
}

// Module holds all constructed supplier_subscription views.
type Module struct {
	routes           centymo.SupplierSubscriptionRoutes
	Dashboard        view.View
	List             view.View
	Table            view.View
	Add              view.View
	Edit             view.View
	Delete           view.View
	BulkDelete       view.View
	SetStatus        view.View
	BulkSetStatus    view.View
	Detail           view.View
	TabAction        view.View
	RecognizeExpense view.View
}

// NewModule creates the supplier_subscription module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &suppliersubscriptionaction.Deps{
		Routes:                              deps.Routes,
		Labels:                              deps.Labels,
		CommonLabels:                        deps.CommonLabels,
		CreateSupplierSubscription:          deps.CreateSupplierSubscription,
		ReadSupplierSubscription:            deps.ReadSupplierSubscription,
		UpdateSupplierSubscription:          deps.UpdateSupplierSubscription,
		DeleteSupplierSubscription:          deps.DeleteSupplierSubscription,
		GetSupplierSubscriptionItemPageData: deps.GetSupplierSubscriptionItemPageData,
		SetSupplierSubscriptionActive:       deps.SetSupplierSubscriptionActive,
	}

	listDeps := &suppliersubscriptionlist.ListViewDeps{
		Routes:                              deps.Routes,
		GetSupplierSubscriptionListPageData: deps.GetSupplierSubscriptionListPageData,
		Labels:                              deps.Labels,
		CommonLabels:                        deps.CommonLabels,
		TableLabels:                         deps.TableLabels,
	}
	listView := suppliersubscriptionlist.NewView(listDeps)
	tableView := suppliersubscriptionlist.NewTableView(listDeps)

	detailDeps := &suppliersubscriptiondetail.DetailViewDeps{
		Routes:                              deps.Routes,
		Labels:                              deps.Labels,
		ExpenseRecognitionRunLabels:         deps.ExpenseRecognitionRunLabels,
		CommonLabels:                        deps.CommonLabels,
		TableLabels:                         deps.TableLabels,
		ReadSupplierSubscription:            deps.ReadSupplierSubscription,
		GetSupplierSubscriptionItemPageData: deps.GetSupplierSubscriptionItemPageData,
		ReadCostPlan:                        deps.ReadCostPlan,
		ListExpenseRecognitions:             deps.ListExpenseRecognitions,
		ExpenseRecognitionDetailURL:         deps.ExpenseRecognitionDetailURL,
	}

	recognizeExpenseDeps := &suppliersubscriptionaction.RecognizeExpenseDeps{
		Routes:                              deps.Routes,
		Labels:                              deps.Labels,
		GetSupplierSubscriptionItemPageData: deps.GetSupplierSubscriptionItemPageData,
		RecognizeFromExpenditure:            deps.RecognizeFromExpenditure,
	}

	return &Module{
		routes:           deps.Routes,
		Dashboard:        listView,
		List:             listView,
		Table:            tableView,
		Add:              suppliersubscriptionaction.NewAddAction(actionDeps),
		Edit:             suppliersubscriptionaction.NewEditAction(actionDeps),
		Delete:           suppliersubscriptionaction.NewDeleteAction(actionDeps),
		BulkDelete:       suppliersubscriptionaction.NewBulkDeleteAction(actionDeps),
		SetStatus:        suppliersubscriptionaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    suppliersubscriptionaction.NewBulkSetStatusAction(actionDeps),
		Detail:           suppliersubscriptiondetail.NewView(detailDeps),
		TabAction:        suppliersubscriptiondetail.NewTabAction(detailDeps),
		RecognizeExpense: suppliersubscriptionaction.NewRecognizeExpenseAction(recognizeExpenseDeps),
	}
}

// RegisterRoutes registers all supplier_subscription routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	if m.Detail != nil && m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, m.Detail)
	}
	if m.TabAction != nil && m.routes.TabActionURL != "" {
		r.GET(m.routes.TabActionURL, m.TabAction)
	}
	if m.RecognizeExpense != nil && m.routes.RecognizeExpenseURL != "" {
		r.POST(m.routes.RecognizeExpenseURL, m.RecognizeExpense)
	}
}
