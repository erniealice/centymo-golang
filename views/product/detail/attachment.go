package detail

import (
	"net/http"

	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/view"
)

func attachmentConfig(deps *DetailViewDeps) *attachment.Config {
	return &attachment.Config{
		EntityType:       "product",
		BucketName:       "attachments",
		RefreshURL:       deps.Routes.TabActionURL,
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

// NewAttachmentUploadAction creates the attachment upload view (GET=drawer form, POST=upload).
func NewAttachmentUploadAction(deps *DetailViewDeps) view.View {
	return attachment.NewUploadAction(attachmentConfig(deps))
}

// NewAttachmentDeleteAction creates the attachment delete view (POST delete).
func NewAttachmentDeleteAction(deps *DetailViewDeps) view.View {
	return attachment.NewDeleteAction(attachmentConfig(deps))
}

// NewAttachmentDownloadHandler creates the GET preview/download handler. Streams
// the stored bytes inline (Content-Disposition: inline) so the browser opens
// supported types in a new tab via window.open(url, '_blank').
func NewAttachmentDownloadHandler(deps *DetailViewDeps) http.HandlerFunc {
	return attachment.NewDownloadHandler(attachmentConfig(deps))
}
