package item

import (
	"github.com/erniealice/pyeza-golang/attachment"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/centymo-golang/views/product/detail/variant"
)

func stockAttachmentConfig(deps *variant.Deps) *attachment.Config {
	return &attachment.Config{
		EntityType:       "stock-item",
		BucketName:       "attachments",
		UploadURL:        deps.Routes.VariantStockAttachmentUploadURL,
		DeleteURL:        deps.Routes.VariantStockAttachmentDeleteURL,
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

// NewAttachmentUploadAction creates the upload handler for stock item attachments.
func NewAttachmentUploadAction(deps *variant.Deps) view.View {
	return attachment.NewUploadAction(stockAttachmentConfig(deps))
}

// NewAttachmentDeleteAction creates the delete handler for stock item attachments.
func NewAttachmentDeleteAction(deps *variant.Deps) view.View {
	return attachment.NewDeleteAction(stockAttachmentConfig(deps))
}
