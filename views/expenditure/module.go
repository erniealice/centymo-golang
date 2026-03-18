package expenditure

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	templateview "github.com/erniealice/hybra-golang/views/template"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/template"
	expenditurepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/expenditure"
	expenditurelist "github.com/erniealice/centymo-golang/views/expenditure/list"
	expendituresettings "github.com/erniealice/centymo-golang/views/expenditure/settings"
)

// ModuleDeps holds all dependencies for the expenditure module.
type ModuleDeps struct {
	Routes           centymo.ExpenditureRoutes
	DB               centymo.DataSource
	ListExpenditures func(ctx context.Context, req *expenditurepb.ListExpendituresRequest) (*expenditurepb.ListExpendituresResponse, error)
	Labels           centymo.ExpenditureLabels
	TemplateLabels   templateview.Labels
	CommonLabels     pyeza.CommonLabels
	TableLabels      types.TableLabels

	// Document template CRUD
	ListDocumentTemplates  func(ctx context.Context, req *documenttemplatepb.ListDocumentTemplatesRequest) (*documenttemplatepb.ListDocumentTemplatesResponse, error)
	CreateDocumentTemplate func(ctx context.Context, req *documenttemplatepb.CreateDocumentTemplateRequest) (*documenttemplatepb.CreateDocumentTemplateResponse, error)
	UpdateDocumentTemplate func(ctx context.Context, req *documenttemplatepb.UpdateDocumentTemplateRequest) (*documenttemplatepb.UpdateDocumentTemplateResponse, error)
	DeleteDocumentTemplate func(ctx context.Context, req *documenttemplatepb.DeleteDocumentTemplateRequest) (*documenttemplatepb.DeleteDocumentTemplateResponse, error)
	UploadFile             func(ctx context.Context, bucket, key string, content []byte, contentType string) error
}

// Module holds all constructed expenditure views.
type Module struct {
	routes            centymo.ExpenditureRoutes
	PurchaseList      view.View
	PurchaseDashboard view.View
	ExpenseList       view.View
	ExpenseDashboard  view.View

	// Settings (template management)
	SettingsTemplates  view.View
	SettingsUpload     view.View
	SettingsDelete     view.View
	SettingsSetDefault view.View
}

// NewModule creates the expenditure module with purchase and expense views.
func NewModule(deps *ModuleDeps) *Module {
	m := &Module{
		routes: deps.Routes,
		PurchaseList: expenditurelist.NewView(&expenditurelist.Deps{
			ListExpenditures: deps.ListExpenditures,
			RefreshURL:       deps.Routes.PurchaseListURL,
			ExpenditureType:  "purchase",
			Labels:           deps.Labels,
			CommonLabels:     deps.CommonLabels,
			TableLabels:      deps.TableLabels,
		}),
		ExpenseList: expenditurelist.NewView(&expenditurelist.Deps{
			ListExpenditures: deps.ListExpenditures,
			RefreshURL:       deps.Routes.ExpenseListURL,
			ExpenditureType:  "expense",
			Labels:           deps.Labels,
			CommonLabels:     deps.CommonLabels,
			TableLabels:      deps.TableLabels,
		}),
		// Dashboards use same list view for now (will be enhanced later)
		PurchaseDashboard: expenditurelist.NewView(&expenditurelist.Deps{
			ListExpenditures: deps.ListExpenditures,
			RefreshURL:       deps.Routes.PurchaseListURL,
			ExpenditureType:  "purchase",
			Labels:           deps.Labels,
			CommonLabels:     deps.CommonLabels,
			TableLabels:      deps.TableLabels,
		}),
		ExpenseDashboard: expenditurelist.NewView(&expenditurelist.Deps{
			ListExpenditures: deps.ListExpenditures,
			RefreshURL:       deps.Routes.ExpenseListURL,
			ExpenditureType:  "expense",
			Labels:           deps.Labels,
			CommonLabels:     deps.CommonLabels,
			TableLabels:      deps.TableLabels,
		}),
	}

	// Settings views (nil-guarded — only built when document template deps are provided)
	if deps.ListDocumentTemplates != nil {
		settingsDeps := &expendituresettings.Deps{
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

	// Settings routes (nil-guarded)
	if m.SettingsTemplates != nil {
		r.GET(m.routes.SettingsTemplatesURL, m.SettingsTemplates)
		r.GET(m.routes.SettingsTemplateUploadURL, m.SettingsUpload)
		r.POST(m.routes.SettingsTemplateUploadURL, m.SettingsUpload)
		r.POST(m.routes.SettingsTemplateDeleteURL, m.SettingsDelete)
		r.POST(m.routes.SettingsTemplateDefaultURL, m.SettingsSetDefault)
	}
}
