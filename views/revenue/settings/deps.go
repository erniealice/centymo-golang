package settings

import (
	"context"

	centymo "github.com/erniealice/centymo-golang"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/document_template"
)

// Deps holds view dependencies for the sales settings (template management) views.
type Deps struct {
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
