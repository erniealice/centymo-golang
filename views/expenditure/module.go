package expenditure

import (
	"context"
	"database/sql"
	"time"

	centymo "github.com/erniealice/centymo-golang"
	templateviewform "github.com/erniealice/hybra-golang/views/template/form"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	expenditureaction "github.com/erniealice/centymo-golang/views/expenditure/action"
	expenditurecategory "github.com/erniealice/centymo-golang/views/expenditure/category"
	expendituredetail "github.com/erniealice/centymo-golang/views/expenditure/detail"
	expenseboard "github.com/erniealice/centymo-golang/views/expenditure/expense_dashboard"
	expenditurelist "github.com/erniealice/centymo-golang/views/expenditure/list"
	expenditurepay "github.com/erniealice/centymo-golang/views/expenditure/pay"
	purchaseboard "github.com/erniealice/centymo-golang/views/expenditure/purchase_dashboard"
	expendituresettings "github.com/erniealice/centymo-golang/views/expenditure/settings"
	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/template"
	supplierpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/supplier"
	accruedexpensepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/accrued_expense"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	expenditurecategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_category"
	expenditurelineitempb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure_line_item"
	expenserecognitionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expense_recognition"
	purchaseorderpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/purchase_order"
	disbursementpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement"
	disbursementschedulepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/disbursement_schedule"
)


// PaymentTermOption is re-exported from action for use by callers wiring ModuleDeps.
// The underlying type lives in views/expenditure/form; action re-exports it as an alias.
type PaymentTermOption = expenditureaction.PaymentTermOption

// ModuleDeps holds all dependencies for the expenditure module.
type ModuleDeps struct {
	Routes           centymo.ExpenditureRoutes
	DB               centymo.DataSource
	SqlDB            *sql.DB
	ListExpenditures func(ctx context.Context, req *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error)
	Labels           centymo.ExpenditureLabels
	TemplateLabels   templateviewform.Labels
	CommonLabels     pyeza.CommonLabels
	TableLabels      types.TableLabels

	// Payment terms dropdown (optional — gracefully degrades when nil)
	ListPaymentTerms func(ctx context.Context) ([]*PaymentTermOption, error)

	// Expense CRUD operations (for action handlers)
	CreateExpenditure func(ctx context.Context, req *expenditurepb.CreateExpenditureRequest) (*expenditurepb.CreateExpenditureResponse, error)
	ReadExpenditure   func(ctx context.Context, req *expenditurepb.ReadExpenditureRequest) (*expenditurepb.ReadExpenditureResponse, error)
	UpdateExpenditure func(ctx context.Context, req *expenditurepb.UpdateExpenditureRequest) (*expenditurepb.UpdateExpenditureResponse, error)
	DeleteExpenditure func(ctx context.Context, req *expenditurepb.DeleteExpenditureRequest) (*expenditurepb.DeleteExpenditureResponse, error)

	// Category listing (optional — gracefully degrades when nil)
	ListExpenditureCategories func(ctx context.Context, req *expenditurecategorypb.ListExpenditureCategoriesRequest) (*expenditurecategorypb.ListExpenditureCategoriesResponse, error)

	// Supplier listing (optional — gracefully degrades when nil)
	ListSuppliers func(ctx context.Context, req *supplierpb.ListSuppliersRequest) (*supplierpb.ListSuppliersResponse, error)

	// Purchase order listing (optional — used to populate PO dropdown on expense form)
	ListPurchaseOrders func(ctx context.Context, req *purchaseorderpb.ListPurchaseOrdersRequest) (*purchaseorderpb.ListPurchaseOrdersResponse, error)

	// Category CRUD (optional — only built when provided)
	CreateExpenditureCategory func(ctx context.Context, req *expenditurecategorypb.CreateExpenditureCategoryRequest) (*expenditurecategorypb.CreateExpenditureCategoryResponse, error)
	ReadExpenditureCategory   func(ctx context.Context, req *expenditurecategorypb.ReadExpenditureCategoryRequest) (*expenditurecategorypb.ReadExpenditureCategoryResponse, error)
	UpdateExpenditureCategory func(ctx context.Context, req *expenditurecategorypb.UpdateExpenditureCategoryRequest) (*expenditurecategorypb.UpdateExpenditureCategoryResponse, error)
	DeleteExpenditureCategory func(ctx context.Context, req *expenditurecategorypb.DeleteExpenditureCategoryRequest) (*expenditurecategorypb.DeleteExpenditureCategoryResponse, error)

	// Expense line item CRUD (optional — only built when provided)
	CreateExpenditureLineItem func(ctx context.Context, req *expenditurelineitempb.CreateExpenditureLineItemRequest) (*expenditurelineitempb.CreateExpenditureLineItemResponse, error)
	ReadExpenditureLineItem   func(ctx context.Context, req *expenditurelineitempb.ReadExpenditureLineItemRequest) (*expenditurelineitempb.ReadExpenditureLineItemResponse, error)
	UpdateExpenditureLineItem func(ctx context.Context, req *expenditurelineitempb.UpdateExpenditureLineItemRequest) (*expenditurelineitempb.UpdateExpenditureLineItemResponse, error)
	DeleteExpenditureLineItem func(ctx context.Context, req *expenditurelineitempb.DeleteExpenditureLineItemRequest) (*expenditurelineitempb.DeleteExpenditureLineItemResponse, error)
	ListExpenditureLineItems  func(ctx context.Context, req *expenditurelineitempb.ListExpenditureLineItemsRequest) (*expenditurelineitempb.ListExpenditureLineItemsResponse, error)

	// Document template CRUD
	ListDocumentTemplates  func(ctx context.Context, req *documenttemplatepb.ListDocumentTemplatesRequest) (*documenttemplatepb.ListDocumentTemplatesResponse, error)
	CreateDocumentTemplate func(ctx context.Context, req *documenttemplatepb.CreateDocumentTemplateRequest) (*documenttemplatepb.CreateDocumentTemplateResponse, error)
	UpdateDocumentTemplate func(ctx context.Context, req *documenttemplatepb.UpdateDocumentTemplateRequest) (*documenttemplatepb.UpdateDocumentTemplateResponse, error)
	DeleteDocumentTemplate func(ctx context.Context, req *documenttemplatepb.DeleteDocumentTemplateRequest) (*documenttemplatepb.DeleteDocumentTemplateResponse, error)
	UploadFile             func(ctx context.Context, bucket, key string, content []byte, contentType string) error

	// Disbursement creation (optional — enables Pay action on expense detail)
	DisbursementRoutes centymo.DisbursementRoutes
	DisbursementLabels centymo.DisbursementLabels
	CreateDisbursement func(ctx context.Context, req *disbursementpb.CreateDisbursementRequest) (*disbursementpb.CreateDisbursementResponse, error)

	// SPS Wave 4 — Recognition + Accrual tabs on the expense detail page.
	// All optional; nil-safe — when missing, the tabs render empty states.
	ReadExpenseRecognition      func(ctx context.Context, req *expenserecognitionpb.ReadExpenseRecognitionRequest) (*expenserecognitionpb.ReadExpenseRecognitionResponse, error)
	ListAccruedExpenses         func(ctx context.Context, req *accruedexpensepb.ListAccruedExpensesRequest) (*accruedexpensepb.ListAccruedExpensesResponse, error)
	ExpenseRecognitionDetailURL string // /app/expense-recognitions/detail/{id}
	AccruedExpenseDetailURL     string // /app/accrued-expenses/detail/{id}
	RecognizeFromExpenditureURL string // /action/expense-recognition/recognize-from-expenditure (POST trigger)

	// Phase 5 — purchase/expense dashboard data callbacks. Nil-safe; the
	// dashboards fall back to zero values when the orchestrator hasn't
	// wired the espyna expenditure dashboard use case yet.
	GetPurchaseDashboardPageData func(ctx context.Context, req *purchaseboard.Request) (*purchaseboard.Response, error)
	GetExpenseDashboardPageData  func(ctx context.Context, req *expenseboard.Request) (*expenseboard.Response, error)
}

// Module holds all constructed expenditure views.
type Module struct {
	routes            centymo.ExpenditureRoutes
	PurchaseList      view.View
	PurchaseDashboard view.View
	ExpenseList       view.View
	ExpenseDashboard  view.View

	// Expense detail page
	ExpenseDetail    view.View
	ExpenseTabAction view.View

	// Expense pay action (creates pre-linked disbursement)
	ExpensePay view.View

	// Expense CRUD actions
	ExpenseAdd       view.View
	ExpenseEdit      view.View
	ExpenseDelete    view.View
	ExpenseSetStatus view.View

	// Expense line item actions
	ExpenseLineItemAdd    view.View
	ExpenseLineItemEdit   view.View
	ExpenseLineItemRemove view.View
	ExpenseLineItemTable  view.View

	// Settings (template management)
	SettingsTemplates  view.View
	SettingsUpload     view.View
	SettingsDelete     view.View
	SettingsSetDefault view.View

	// Category CRUD
	CategoryList   view.View
	CategoryAdd    view.View
	CategoryEdit   view.View
	CategoryDelete view.View
}

// NewModule creates the expenditure module with purchase and expense views.
func NewModule(deps *ModuleDeps) *Module {
	actionDeps := &expenditureaction.Deps{
		Routes:                    deps.Routes,
		Labels:                    deps.Labels,
		ListPaymentTerms:          deps.ListPaymentTerms,
		CreateExpenditure:         deps.CreateExpenditure,
		ReadExpenditure:           deps.ReadExpenditure,
		UpdateExpenditure:         deps.UpdateExpenditure,
		DeleteExpenditure:         deps.DeleteExpenditure,
		ListExpenditureCategories: deps.ListExpenditureCategories,
		ListSuppliers:             deps.ListSuppliers,
		ListPurchaseOrders:        deps.ListPurchaseOrders,
	}

	m := &Module{
		routes: deps.Routes,
		PurchaseList: expenditurelist.NewView(&expenditurelist.ListViewDeps{
			ListExpenditures: deps.ListExpenditures,
			RefreshURL:       deps.Routes.PurchaseListURL,
			ExpenditureType:  "purchase",
			Labels:           deps.Labels,
			CommonLabels:     deps.CommonLabels,
			TableLabels:      deps.TableLabels,
		}),
		ExpenseList: expenditurelist.NewView(&expenditurelist.ListViewDeps{
			ListExpenditures: deps.ListExpenditures,
			RefreshURL:       deps.Routes.ExpenseListURL,
			ExpenditureType:  "expense",
			AddURL:           deps.Routes.AddURL,
			Labels:           deps.Labels,
			CommonLabels:     deps.CommonLabels,
			TableLabels:      deps.TableLabels,
		}),
		// Phase 5 — real dashboards backed by espyna expenditure dashboard
		// use case (kind=purchase / kind=expense via separate callbacks).
		PurchaseDashboard: purchaseboard.NewView(&purchaseboard.Deps{
			Routes:       deps.Routes,
			Labels:       deps.Labels,
			CommonLabels: deps.CommonLabels,
			GetPageData:  deps.GetPurchaseDashboardPageData,
		}),
		ExpenseDashboard: expenseboard.NewView(&expenseboard.Deps{
			Routes:       deps.Routes,
			Labels:       deps.Labels,
			CommonLabels: deps.CommonLabels,
			GetPageData:  deps.GetExpenseDashboardPageData,
		}),
	}

	// Expense CRUD actions (nil-guarded — only built when CRUD deps are provided)
	if deps.CreateExpenditure != nil {
		m.ExpenseAdd = expenditureaction.NewAddAction(actionDeps)
		m.ExpenseEdit = expenditureaction.NewEditAction(actionDeps)
		m.ExpenseDelete = expenditureaction.NewDeleteAction(actionDeps)
		m.ExpenseSetStatus = expenditureaction.NewSetStatusAction(actionDeps)
	}

	// Category views (nil-guarded — only built when category deps are provided)
	if deps.ListExpenditureCategories != nil {
		m.CategoryList = expenditurecategory.NewView(&expenditurecategory.ListViewDeps{
			Routes:                    deps.Routes,
			ListExpenditureCategories: deps.ListExpenditureCategories,
			Labels:                    deps.Labels,
			CommonLabels:              deps.CommonLabels,
			TableLabels:               deps.TableLabels,
		})
	}
	if deps.CreateExpenditureCategory != nil {
		catActionDeps := &expenditurecategory.ActionDeps{
			Routes:                    deps.Routes,
			Labels:                    deps.Labels,
			CreateExpenditureCategory: deps.CreateExpenditureCategory,
			ReadExpenditureCategory:   deps.ReadExpenditureCategory,
			UpdateExpenditureCategory: deps.UpdateExpenditureCategory,
			DeleteExpenditureCategory: deps.DeleteExpenditureCategory,
		}
		m.CategoryAdd = expenditurecategory.NewAddAction(catActionDeps)
		m.CategoryEdit = expenditurecategory.NewEditAction(catActionDeps)
		m.CategoryDelete = expenditurecategory.NewDeleteAction(catActionDeps)
	}

	// Expense detail page (nil-guarded — only built when ReadExpenditure is provided)
	if deps.ReadExpenditure != nil {
		detailDeps := &expendituredetail.DetailViewDeps{
			Routes:                      deps.Routes,
			Labels:                      deps.Labels,
			CommonLabels:                deps.CommonLabels,
			TableLabels:                 deps.TableLabels,
			ReadExpenditure:             deps.ReadExpenditure,
			ReadExpenseRecognition:      deps.ReadExpenseRecognition,
			ListAccruedExpenses:         deps.ListAccruedExpenses,
			ExpenseRecognitionDetailURL: deps.ExpenseRecognitionDetailURL,
			AccruedExpenseDetailURL:     deps.AccruedExpenseDetailURL,
			RecognizeFromExpenditureURL: deps.RecognizeFromExpenditureURL,
		}
		if deps.ListExpenditureLineItems != nil {
			detailDeps.ListExpenditureLineItems = deps.ListExpenditureLineItems
		}
		if deps.SqlDB != nil {
			sqlDB := deps.SqlDB
			detailDeps.GetPaidAmount = func(ctx context.Context, expenditureID string) (int64, error) {
				var total int64
				err := sqlDB.QueryRowContext(ctx,
					`SELECT COALESCE(SUM(amount), 0) FROM treasury_disbursement
					 WHERE expenditure_id = $1 AND active = true AND status IN ('paid', 'completed')`,
					expenditureID,
				).Scan(&total)
				return total, err
			}
		}
		if deps.SqlDB != nil {
			sqlDB2 := deps.SqlDB
			detailDeps.ListDisbursementSchedules = func(ctx context.Context, expenditureID string) ([]*disbursementschedulepb.DisbursementSchedule, error) {
				rows, err := sqlDB2.QueryContext(ctx,
					`SELECT id, sequence, amount, due_date, status,
					        paid_amount, paid_date, disbursement_id
					 FROM disbursement_schedule
					 WHERE expenditure_id = $1 AND active = true
					 ORDER BY sequence ASC`,
					expenditureID,
				)
				if err != nil {
					return nil, err
				}
				defer rows.Close()

				var schedules []*disbursementschedulepb.DisbursementSchedule
				for rows.Next() {
					var (
						id             string
						sequence       int32
						amount         int64
						dueDateMillis  int64
						status         string
						paidAmount     *int64
						paidDate       *int64
						disbursementID *string
					)
					if err := rows.Scan(&id, &sequence, &amount, &dueDateMillis, &status, &paidAmount, &paidDate, &disbursementID); err != nil {
						return nil, err
					}
					s := &disbursementschedulepb.DisbursementSchedule{
						Id:             id,
						ExpenditureId:  expenditureID,
						Sequence:       sequence,
						Amount:         amount,
						DueDate:        formatDisbEpochMillis(dueDateMillis),
						Status:         status,
						PaidAmount:     paidAmount,
						PaidDate:       paidDate,
						DisbursementId: disbursementID,
					}
					schedules = append(schedules, s)
				}
				return schedules, rows.Err()
			}
		}
		m.ExpenseDetail = expendituredetail.NewView(detailDeps)
		m.ExpenseTabAction = expendituredetail.NewTabAction(detailDeps)
	}

	// Expense pay action (nil-guarded — only built when CreateDisbursement and ReadExpenditure are provided)
	if deps.CreateDisbursement != nil && deps.ReadExpenditure != nil {
		m.ExpensePay = expenditurepay.NewPayAction(&expenditurepay.Deps{
			ExpenditureRoutes:  deps.Routes,
			DisbursementRoutes: deps.DisbursementRoutes,
			DisbursementLabels: deps.DisbursementLabels,
			ReadExpenditure:    deps.ReadExpenditure,
			CreateDisbursement: deps.CreateDisbursement,
		})
	}

	// Expense line item actions (nil-guarded)
	if deps.CreateExpenditureLineItem != nil {
		lineItemDeps := &expendituredetail.LineItemDeps{
			Routes:                    deps.Routes,
			Labels:                    deps.Labels,
			CommonLabels:              deps.CommonLabels,
			TableLabels:               deps.TableLabels,
			ReadExpenditure:           deps.ReadExpenditure,
			UpdateExpenditure:         deps.UpdateExpenditure,
			CreateExpenditureLineItem: deps.CreateExpenditureLineItem,
			ReadExpenditureLineItem:   deps.ReadExpenditureLineItem,
			UpdateExpenditureLineItem: deps.UpdateExpenditureLineItem,
			DeleteExpenditureLineItem: deps.DeleteExpenditureLineItem,
			ListExpenditureLineItems:  deps.ListExpenditureLineItems,
		}
		m.ExpenseLineItemAdd = expendituredetail.NewLineItemAddView(lineItemDeps)
		m.ExpenseLineItemEdit = expendituredetail.NewLineItemEditView(lineItemDeps)
		m.ExpenseLineItemRemove = expendituredetail.NewLineItemRemoveView(lineItemDeps)
		m.ExpenseLineItemTable = expendituredetail.NewLineItemTableView(lineItemDeps)
	}

	// Settings views (nil-guarded — only built when document template deps are provided)
	if deps.ListDocumentTemplates != nil {
		settingsDeps := &expendituresettings.SettingsViewDeps{
			Routes:                 deps.Routes,
			Labels:                 deps.TemplateLabels,
			CommonLabels:           deps.CommonLabels,
			TableLabels:            deps.TableLabels,
			ListDocumentTemplates:  deps.ListDocumentTemplates,
			CreateDocumentTemplate: deps.CreateDocumentTemplate,
			UpdateDocumentTemplate: deps.UpdateDocumentTemplate,
			DeleteDocumentTemplate: deps.DeleteDocumentTemplate,
			UploadFile:             deps.UploadFile,
		}
		m.SettingsTemplates = expendituresettings.NewView(settingsDeps)
		m.SettingsUpload = expendituresettings.NewUploadAction(settingsDeps)
		m.SettingsDelete = expendituresettings.NewDeleteAction(settingsDeps)
		m.SettingsSetDefault = expendituresettings.NewSetDefaultAction(settingsDeps)
	}

	return m
}

// RegisterRoutes registers all expenditure routes.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.PurchaseListURL, m.PurchaseList)
	r.GET(m.routes.PurchaseDashboardURL, m.PurchaseDashboard)
	r.GET(m.routes.ExpenseListURL, m.ExpenseList)
	r.GET(m.routes.ExpenseDashboardURL, m.ExpenseDashboard)

	// Expense detail page (nil-guarded)
	if m.ExpenseDetail != nil {
		r.GET(m.routes.DetailURL, m.ExpenseDetail)
		r.GET(m.routes.TabActionURL, m.ExpenseTabAction)
	}

	// Expense pay action routes (nil-guarded)
	if m.ExpensePay != nil {
		r.GET(m.routes.PayURL, m.ExpensePay)
		r.POST(m.routes.PayURL, m.ExpensePay)
	}

	// Expense CRUD action routes (nil-guarded)
	if m.ExpenseAdd != nil {
		r.GET(m.routes.AddURL, m.ExpenseAdd)
		r.POST(m.routes.AddURL, m.ExpenseAdd)
		r.GET(m.routes.EditURL, m.ExpenseEdit)
		r.POST(m.routes.EditURL, m.ExpenseEdit)
		r.POST(m.routes.DeleteURL, m.ExpenseDelete)
		r.POST(m.routes.SetStatusURL, m.ExpenseSetStatus)
	}

	// Expense line item action routes (nil-guarded)
	if m.ExpenseLineItemAdd != nil {
		r.GET(m.routes.LineItemAddURL, m.ExpenseLineItemAdd)
		r.POST(m.routes.LineItemAddURL, m.ExpenseLineItemAdd)
		r.GET(m.routes.LineItemEditURL, m.ExpenseLineItemEdit)
		r.POST(m.routes.LineItemEditURL, m.ExpenseLineItemEdit)
		r.POST(m.routes.LineItemRemoveURL, m.ExpenseLineItemRemove)
		r.GET(m.routes.LineItemTableURL, m.ExpenseLineItemTable)
	}

	// Settings routes (nil-guarded)
	if m.SettingsTemplates != nil {
		r.GET(m.routes.SettingsTemplatesURL, m.SettingsTemplates)
		r.GET(m.routes.SettingsTemplateUploadURL, m.SettingsUpload)
		r.POST(m.routes.SettingsTemplateUploadURL, m.SettingsUpload)
		r.POST(m.routes.SettingsTemplateDeleteURL, m.SettingsDelete)
		r.POST(m.routes.SettingsTemplateDefaultURL, m.SettingsSetDefault)
	}

	// Category routes (nil-guarded)
	if m.CategoryList != nil {
		r.GET(m.routes.ExpenseCategoryListURL, m.CategoryList)
		r.GET(m.routes.ExpenseCategoryTableURL, m.CategoryList)
	}
	if m.CategoryAdd != nil {
		r.GET(m.routes.ExpenseCategoryAddURL, m.CategoryAdd)
		r.POST(m.routes.ExpenseCategoryAddURL, m.CategoryAdd)
		r.GET(m.routes.ExpenseCategoryEditURL, m.CategoryEdit)
		r.POST(m.routes.ExpenseCategoryEditURL, m.CategoryEdit)
		r.POST(m.routes.ExpenseCategoryDeleteURL, m.CategoryDelete)
	}
}

// formatDisbEpochMillis converts epoch milliseconds to a YYYY-MM-DD string.
func formatDisbEpochMillis(ms int64) string {
	if ms == 0 {
		return ""
	}
	return time.UnixMilli(ms).UTC().Format("2006-01-02")
}
