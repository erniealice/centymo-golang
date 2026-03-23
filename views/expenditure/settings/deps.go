package settings

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	templateview "github.com/erniealice/hybra-golang/views/template"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/template"
)

// SettingsViewDeps holds view dependencies for the purchases settings (template management) views.
type SettingsViewDeps struct {
	Routes       centymo.ExpenditureRoutes
	Labels       templateview.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Document template CRUD operations (injected by composition root)
	ListDocumentTemplates  func(ctx context.Context, req *documenttemplatepb.ListDocumentTemplatesRequest) (*documenttemplatepb.ListDocumentTemplatesResponse, error)
	CreateDocumentTemplate func(ctx context.Context, req *documenttemplatepb.CreateDocumentTemplateRequest) (*documenttemplatepb.CreateDocumentTemplateResponse, error)
	UpdateDocumentTemplate func(ctx context.Context, req *documenttemplatepb.UpdateDocumentTemplateRequest) (*documenttemplatepb.UpdateDocumentTemplateResponse, error)
	DeleteDocumentTemplate func(ctx context.Context, req *documenttemplatepb.DeleteDocumentTemplateRequest) (*documenttemplatepb.DeleteDocumentTemplateResponse, error)

	// Storage operations (injected by composition root)
	UploadFile func(ctx context.Context, bucketName, objectKey string, content []byte, contentType string) error
}

func templateConfig(deps *SettingsViewDeps) *templateview.Config {
	return &templateview.Config{
		DocumentPurpose:    "purchase_order",
		AllowedExtensions:  []string{".docx"},
		AllowedContentType: templateview.DocxContentType,
		StoragePrefix:      "templates/purchase_order",
		BucketName:         "templates",
		ListURL:            deps.Routes.SettingsTemplatesURL,
		UploadURL:          deps.Routes.SettingsTemplateUploadURL,
		DeleteURL:          deps.Routes.SettingsTemplateDeleteURL,
		SetDefaultURL:      deps.Routes.SettingsTemplateDefaultURL,
		Labels:             deps.Labels,
		CommonLabels:       deps.CommonLabels,
		TableLabels:        deps.TableLabels,
		ActiveNav:          "purchases",
		PageIcon:           "icon-file-text",
		ListDocumentTemplates:  deps.ListDocumentTemplates,
		CreateDocumentTemplate: deps.CreateDocumentTemplate,
		UpdateDocumentTemplate: deps.UpdateDocumentTemplate,
		DeleteDocumentTemplate: deps.DeleteDocumentTemplate,
		UploadFile:             deps.UploadFile,
	}
}
