package item

import (
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang/views/product/detail/variant"
)

func stockAttachmentConfig(deps *variant.DetailViewDeps) *attachment.Config {
	return &attachment.Config{
		EntityType:         "stock-item",
		BucketName:         "attachments",
		UploadURL:          deps.Routes.VariantStockAttachmentUploadURL,
		DeleteURL:          deps.Routes.VariantStockAttachmentDeleteURL,
		PrimaryIDPathParam: "iid",
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

// NewAttachmentUploadAction creates the upload handler for stock item attachments.
func NewAttachmentUploadAction(deps *variant.DetailViewDeps) view.View {
	return attachment.NewUploadAction(stockAttachmentConfig(deps))
}

// NewAttachmentDeleteAction creates the delete handler for stock item attachments.
func NewAttachmentDeleteAction(deps *variant.DetailViewDeps) view.View {
	return attachment.NewDeleteAction(stockAttachmentConfig(deps))
}
