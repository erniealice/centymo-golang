package detail

import (
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/view"
)

func attachmentConfig(deps *DetailViewDeps) *attachment.Config {
	return &attachment.Config{
		EntityType:       "supplier_contract_price_schedule",
		BucketName:       "attachments",
		UploadURL:        deps.Routes.AttachmentUploadURL,
		DeleteURL:        deps.Routes.AttachmentDeleteURL,
		Labels:           attachment.DefaultLabels(),
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
		NewID:            deps.NewAttachmentID,
		UploadFile:       deps.UploadFile,
		ListAttachments:  deps.ListAttachments,
		CreateAttachment: deps.CreateAttachment,
		DeleteAttachment: deps.DeleteAttachment,
	}
}

// NewAttachmentUploadAction creates the upload handler for supplier_contract_price_schedule attachments.
func NewAttachmentUploadAction(deps *DetailViewDeps) view.View {
	return attachment.NewUploadAction(attachmentConfig(deps))
}

// NewAttachmentDeleteAction creates the delete handler for supplier_contract_price_schedule attachments.
func NewAttachmentDeleteAction(deps *DetailViewDeps) view.View {
	return attachment.NewDeleteAction(attachmentConfig(deps))
}
