package block

import (
	"context"

	"github.com/erniealice/espyna-golang/ports"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/template"
)

// Infra carries the subset of AppContext that view modules need beyond the
// typed UseCases: attachment ops, document generation, template CRUD,
// email dispatch, image/template uploads, reference checker, and cross-package
// URL patterns. Built once by service-admin and passed into each catalog binder.
//
// Unlike fayna's Infra, centymo's Infra does NOT carry a DB field — the
// centymo DataSource duck was deleted in 20260612-datasource-typed-path W6.
// All formerly DB-backed calls now flow through typed closures on *UseCases.
type Infra struct {
	// Attachment ops
	UploadFile       func(context.Context, string, string, []byte, string) error
	DownloadFile     func(context.Context, string, string) ([]byte, error)
	ListAttachments  func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	ReadAttachment   func(context.Context, *attachmentpb.ReadAttachmentRequest) (*attachmentpb.ReadAttachmentResponse, error)
	DeleteAttachment func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewAttachmentID  func() string

	// Storage / email / doc generation ops
	UploadImage    func(context.Context, string, string, []byte, string) error
	UploadTemplate func(context.Context, string, string, []byte, string) error
	SendEmail      func(context.Context, []string, string, string, string, string, []byte) error
	GenerateDoc    func([]byte, map[string]any) ([]byte, error)

	// Document template CRUD
	ListDocTemplates   func(context.Context, *documenttemplatepb.ListDocumentTemplatesRequest) (*documenttemplatepb.ListDocumentTemplatesResponse, error)
	CreateDocTemplate  func(context.Context, *documenttemplatepb.CreateDocumentTemplateRequest) (*documenttemplatepb.CreateDocumentTemplateResponse, error)
	UpdateDocTemplate  func(context.Context, *documenttemplatepb.UpdateDocumentTemplateRequest) (*documenttemplatepb.UpdateDocumentTemplateResponse, error)
	DeleteDocTemplate  func(context.Context, *documenttemplatepb.DeleteDocumentTemplateRequest) (*documenttemplatepb.DeleteDocumentTemplateResponse, error)

	// Reference checker (optional — nil-safe)
	RefChecker ports.Checker

	// Cross-package URL patterns (wired by service-admin)
	SupplierDetailURL                      string
	SupplierExpenseRecognitionRunDrawerURL string
}
