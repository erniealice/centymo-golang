package settings

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"
	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/template"
	templateview "github.com/erniealice/hybra-golang/views/template"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// SettingsViewDeps holds view dependencies for the sales settings (template management) views.
type SettingsViewDeps struct {
	Routes       centymo.RevenueRoutes
	Labels       centymo.RevenueLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Document template CRUD operations (injected by composition root)
	ListDocumentTemplates  func(ctx context.Context, req *documenttemplatepb.ListDocumentTemplatesRequest) (*documenttemplatepb.ListDocumentTemplatesResponse, error)
	CreateDocumentTemplate func(ctx context.Context, req *documenttemplatepb.CreateDocumentTemplateRequest) (*documenttemplatepb.CreateDocumentTemplateResponse, error)
	UpdateDocumentTemplate func(ctx context.Context, req *documenttemplatepb.UpdateDocumentTemplateRequest) (*documenttemplatepb.UpdateDocumentTemplateResponse, error)
	DeleteDocumentTemplate func(ctx context.Context, req *documenttemplatepb.DeleteDocumentTemplateRequest) (*documenttemplatepb.DeleteDocumentTemplateResponse, error)

	// Storage operations (injected by composition root)
	UploadTemplate func(ctx context.Context, bucketName, objectKey string, content []byte, contentType string) error
}

func templateConfig(deps *SettingsViewDeps) *templateview.Config {
	return &templateview.Config{
		DocumentPurpose:        "invoice",
		AllowedExtensions:      []string{".docx"},
		AllowedContentType:     templateview.DocxContentType,
		StoragePrefix:          "templates/invoice",
		BucketName:             "templates",
		ListURL:                deps.Routes.SettingsTemplatesURL,
		UploadURL:              deps.Routes.SettingsTemplateUploadURL,
		DeleteURL:              deps.Routes.SettingsTemplateDeleteURL,
		SetDefaultURL:          deps.Routes.SettingsTemplateDefaultURL,
		Labels:                 convertLabels(deps.Labels.Settings),
		CommonLabels:           deps.CommonLabels,
		TableLabels:            deps.TableLabels,
		ActiveNav:              "revenue",
		PageIcon:               "icon-file-text",
		ListDocumentTemplates:  deps.ListDocumentTemplates,
		CreateDocumentTemplate: deps.CreateDocumentTemplate,
		UpdateDocumentTemplate: deps.UpdateDocumentTemplate,
		DeleteDocumentTemplate: deps.DeleteDocumentTemplate,
		UploadFile:             deps.UploadTemplate,
	}
}

func convertLabels(l centymo.RevenueSettingsLabels) templateview.Labels {
	return templateview.Labels{
		PageTitle:      l.PageTitle,
		Caption:        l.Caption,
		UploadTemplate: l.UploadTemplate,
		TemplateName:   l.TemplateName,
		TemplateType:   l.TemplateType,
		Purpose:        l.Purpose,
		SetDefault:     l.SetDefault,
		Delete:         l.Delete,
		DefaultBadge:   l.DefaultBadge,
		EmptyTitle:     l.EmptyTitle,
		EmptyMessage:   l.EmptyMessage,
		UploadSuccess:  l.UploadSuccess,
		DeleteConfirm:  l.DeleteConfirm,
	}
}
