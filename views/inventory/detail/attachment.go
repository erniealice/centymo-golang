package detail

import (
	"github.com/erniealice/fycha-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/view"
)

func attachmentConfig(deps *Deps) *attachment.Config {
	return &attachment.Config{
		EntityType:       "inventory",
		BucketName:       "attachments",
		UploadURL:        deps.Routes.AttachmentUploadURL,
		DeleteURL:        deps.Routes.AttachmentDeleteURL,
		Labels:           attachment.DefaultLabels(),
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
		NewID:            deps.NewID,
		UploadFile:       deps.UploadFile,
		ListAttachments:  deps.ListAttachments,
		CreateAttachment: deps.CreateAttachment,
		DeleteAttachment: deps.DeleteAttachment,
	}
}

// NewAttachmentUploadAction creates the attachment upload view (GET=drawer form, POST=upload).
func NewAttachmentUploadAction(deps *Deps) view.View {
	return attachment.NewUploadAction(attachmentConfig(deps))
}

// NewAttachmentDeleteAction creates the attachment delete view (POST delete).
func NewAttachmentDeleteAction(deps *Deps) view.View {
	return attachment.NewDeleteAction(attachmentConfig(deps))
}
