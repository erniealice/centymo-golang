package plan

import (
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/view"
)

func attachmentConfig(deps *DetailViewDeps) *attachment.Config {
	return &attachment.Config{
		EntityType:         "price_plan",
		BucketName:         "attachments",
		UploadURL:          deps.Routes.PlanAttachmentUploadURL,
		DeleteURL:          deps.Routes.PlanAttachmentDeleteURL,
		PrimaryIDPathParam: "ppid",
		Labels:             attachment.DefaultLabels(),
		CommonLabels:       deps.CommonLabels,
		TableLabels:        deps.TableLabels,
		NewID:              deps.NewAttachmentID,
		UploadFile:         deps.UploadFile,
		ListAttachments:    deps.ListAttachments,
		CreateAttachment:   deps.CreateAttachment,
		DeleteAttachment:   deps.DeleteAttachment,
	}
}

// NewAttachmentUploadAction creates the upload handler for price_schedule/plan attachments.
func NewAttachmentUploadAction(deps *DetailViewDeps) view.View {
	return attachment.NewUploadAction(attachmentConfig(deps))
}

// NewAttachmentDeleteAction creates the delete handler for price_schedule/plan attachments.
func NewAttachmentDeleteAction(deps *DetailViewDeps) view.View {
	return attachment.NewDeleteAction(attachmentConfig(deps))
}
