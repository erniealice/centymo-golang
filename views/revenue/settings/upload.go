package settings

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	centymo "github.com/erniealice/centymo-golang"

	"github.com/erniealice/pyeza-golang/view"

	"github.com/google/uuid"

	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/document_template"
)

// maxTemplateUploadSize is the maximum size for a template upload (10 MB).
const maxTemplateUploadSize = 10 << 20

// docxContentType is the MIME type for .docx files.
const docxContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"

// UploadFormData is the template data for the upload drawer form.
type UploadFormData struct {
	FormAction   string
	Labels       centymo.RevenueSettingsLabels
	CommonLabels any
}

// NewUploadAction creates the handler for uploading invoice templates.
// GET = drawer form, POST = upload + create record.
// Route: /action/sales/settings/templates/upload
func NewUploadAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("settings-upload-drawer-form", &UploadFormData{
				FormAction:   deps.Routes.SettingsTemplateUploadURL,
				Labels:       deps.Labels.Settings,
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — upload template
		if deps.UploadTemplate == nil || deps.CreateDocumentTemplate == nil {
			log.Printf("Template upload deps not configured")
			return centymo.HTMXError("template upload not configured")
		}

		err := viewCtx.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Printf("Failed to parse multipart form: %v", err)
			return centymo.HTMXError("failed to parse upload")
		}

		// Get template name (required)
		name := viewCtx.Request.FormValue("name")
		if name == "" {
			return centymo.HTMXError("template name is required")
		}

		// Get the uploaded file
		fh, header, err := viewCtx.Request.FormFile("template_file")
		if err != nil {
			log.Printf("Failed to get uploaded file: %v", err)
			return centymo.HTMXError("no file provided")
		}
		defer fh.Close()

		// Validate file size
		if header.Size > maxTemplateUploadSize {
			return centymo.HTMXError(fmt.Sprintf("file too large: %d bytes (max %d)", header.Size, maxTemplateUploadSize))
		}

		// Validate content type (.docx only)
		ct := header.Header.Get("Content-Type")
		if ct != docxContentType {
			return centymo.HTMXError(fmt.Sprintf("invalid file type: only .docx files are accepted"))
		}

		// Read file content
		content, err := io.ReadAll(fh)
		if err != nil {
			log.Printf("Failed to read uploaded file: %v", err)
			return centymo.HTMXError("failed to read file")
		}

		// Generate a new UUID for the object ID and database ID.
		newID := uuid.New().String()

		// Generate a safe storage object key using the new UUID.
		objectKey := fmt.Sprintf("templates/invoice/%s.docx", newID)

		// Upload to storage
		bucketName := "templates"
		err = deps.UploadTemplate(ctx, bucketName, objectKey, content, ct)
		if err != nil {
			log.Printf("Failed to upload template: %v", err)
			return centymo.HTMXError("failed to upload template")
		}

		// Create DB record
		fileSize := header.Size
		storageContainer := "templates"
		originalFilename := header.Filename
		_, err = deps.CreateDocumentTemplate(ctx, &documenttemplatepb.CreateDocumentTemplateRequest{
			Data: &documenttemplatepb.DocumentTemplate{
				Id:               newID,
				Name:             name,
				TemplateType:     "docx",
				DocumentPurpose:  "invoice",
				StorageContainer: &storageContainer,
				StorageKey:       &objectKey,
				OriginalFilename: &originalFilename,
				FileSizeBytes:    &fileSize,
				Status:           "active",
				Active:           true,
			},
		})
		if err != nil {
			log.Printf("Failed to create document template record: %v", err)
			return centymo.HTMXError("failed to save template")
		}

		// Redirect back to settings page to show updated list
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": deps.Routes.SettingsTemplatesURL,
			},
		}
	})
}
