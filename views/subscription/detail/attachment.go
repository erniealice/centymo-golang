package detail

import (
	"net/http"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/view"
)

func attachmentConfig(deps *DetailViewDeps) *attachment.Config {
	return &attachment.Config{
		EntityType:       "subscription",
		BucketName:       "attachments",
		UploadURL:        deps.Routes.AttachmentUploadURL,
		DeleteURL:        deps.Routes.AttachmentDeleteURL,
		DownloadURL:      deps.Routes.AttachmentDownloadURL,
		Labels:           attachment.DefaultLabels(),
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
		NewID:            deps.NewAttachmentID,
		UploadFile:       deps.UploadFile,
		DownloadFile:     deps.DownloadFile,
		ListAttachments:  deps.ListAttachments,
		CreateAttachment: deps.CreateAttachment,
		ReadAttachment:   deps.ReadAttachment,
		DeleteAttachment: deps.DeleteAttachment,
	}
}

// NewAttachmentUploadAction creates the upload handler for subscription attachments.
func NewAttachmentUploadAction(deps *DetailViewDeps) view.View {
	return attachment.NewUploadAction(attachmentConfig(deps))
}

// NewAttachmentDeleteAction creates the delete handler for subscription attachments.
func NewAttachmentDeleteAction(deps *DetailViewDeps) view.View {
	return attachment.NewDeleteAction(attachmentConfig(deps))
}

// NewAttachmentDownloadHandler creates the GET preview/download handler. It is
// an http.HandlerFunc rather than a view.View because it streams raw bytes
// instead of rendering a template.
func NewAttachmentDownloadHandler(deps *DetailViewDeps) http.HandlerFunc {
	return attachment.NewDownloadHandler(attachmentConfig(deps))
}
