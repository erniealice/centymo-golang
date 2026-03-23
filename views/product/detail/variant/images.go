package variant

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"path/filepath"
	"strings"

	"github.com/erniealice/pyeza-golang/view"

	productvariantimagepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_variant_image"
)

// ImageData holds display data for a single variant image.
type ImageData struct {
	ID        string
	ImageURL  string // object key (e.g., "products/x/variants/y/img.webp")
	AltText   string
	SortOrder int
	IsPrimary bool
}

// StorageImageBaseURL is the URL prefix for serving images via the storage handler.
const StorageImageBaseURL = "/storage/images/"

// maxUploadSize is the maximum size for a single uploaded file (5 MB).
const maxUploadSize = 5 << 20

// loadVariantImages loads all active images for a variant from the database.
func loadVariantImages(ctx context.Context, deps *DetailViewDeps, variantID string) []ImageData {
	if deps.ListProductVariantImages == nil {
		return nil
	}

	resp, err := deps.ListProductVariantImages(ctx, &productvariantimagepb.ListProductVariantImagesRequest{})
	if err != nil {
		log.Printf("Failed to list product_variant_image: %v", err)
		return nil
	}

	var images []ImageData
	for _, img := range resp.GetData() {
		if img.GetProductVariantId() != variantID || !img.GetActive() {
			continue
		}
		images = append(images, ImageData{
			ID:        img.GetId(),
			ImageURL:  StorageImageBaseURL + img.GetImageUrl(),
			AltText:   img.GetAltText(),
			SortOrder: int(img.GetSortOrder()),
			IsPrimary: img.GetIsPrimary(),
		})
	}
	return images
}

// NewImageUploadAction creates the POST handler for uploading variant images.
// Route: POST /action/products/detail/{id}/variant/{vid}/images/upload
func NewImageUploadAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		vid := viewCtx.Request.PathValue("vid")

		if deps.UploadImage == nil || deps.CreateProductVariantImage == nil {
			log.Printf("Image upload deps not configured")
			return view.Error(fmt.Errorf("image upload not configured"))
		}

		err := viewCtx.Request.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Printf("Failed to parse multipart form: %v", err)
			return view.Error(fmt.Errorf("failed to parse upload: %w", err))
		}

		files := viewCtx.Request.MultipartForm.File["files"]
		if len(files) == 0 {
			return view.Error(fmt.Errorf("no files provided"))
		}

		bucketName := "images" // storage container/bucket for images

		// Get current image count for sort_order
		existingImages := loadVariantImages(ctx, deps, vid)
		sortOrder := int32(len(existingImages))

		for _, fh := range files {
			// Validate file size
			if fh.Size > maxUploadSize {
				log.Printf("File %s too large: %d bytes (max %d)", fh.Filename, fh.Size, maxUploadSize)
				continue
			}

			// Validate content type
			ct := fh.Header.Get("Content-Type")
			if !strings.HasPrefix(ct, "image/") {
				log.Printf("File %s has invalid content type: %s", fh.Filename, ct)
				continue
			}

			// Read file content
			f, err := fh.Open()
			if err != nil {
				log.Printf("Failed to open uploaded file %s: %v", fh.Filename, err)
				continue
			}
			content, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				log.Printf("Failed to read uploaded file %s: %v", fh.Filename, err)
				continue
			}

			// Determine extension from content type
			ext := extensionFromContentType(ct)
			if ext == "" {
				ext = filepath.Ext(fh.Filename)
			}

			// Generate object key: products/{productID}/variants/{variantID}/{sortOrder}{ext}
			objectKey := fmt.Sprintf("products/%s/variants/%s/%02d%s", id, vid, sortOrder, ext)

			// Upload to storage
			err = deps.UploadImage(ctx, bucketName, objectKey, content, ct)
			if err != nil {
				log.Printf("Failed to upload image %s: %v", fh.Filename, err)
				continue
			}

			// Create DB record
			isPrimary := sortOrder == 0 && len(existingImages) == 0
			_, err = deps.CreateProductVariantImage(ctx, &productvariantimagepb.CreateProductVariantImageRequest{
				Data: &productvariantimagepb.ProductVariantImage{
					ProductVariantId: vid,
					ImageUrl:         objectKey,
					AltText:          &fh.Filename,
					SortOrder:        sortOrder,
					IsPrimary:        isPrimary,
					Active:           true,
				},
			})
			if err != nil {
				log.Printf("Failed to create product_variant_image record: %v", err)
				continue
			}

			sortOrder++
		}

		// Return updated gallery partial
		images := loadVariantImages(ctx, deps, vid)
		data := &VariantPageData{
			ProductID: id,
			VariantID: vid,
			Images:    images,
		}
		return view.OK("variant-image-gallery", data)
	})
}

// NewImageDeleteAction creates the POST handler for deleting variant images.
// Route: POST /action/products/detail/{id}/variant/{vid}/images/delete
func NewImageDeleteAction(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		vid := viewCtx.Request.PathValue("vid")

		if deps.DeleteProductVariantImage == nil {
			return view.Error(fmt.Errorf("image delete not configured"))
		}

		viewCtx.Request.ParseForm()
		imageIDs := viewCtx.Request.Form["image_ids[]"]
		if len(imageIDs) == 0 {
			imageIDs = viewCtx.Request.Form["image_ids"]
		}

		for _, imgID := range imageIDs {
			if imgID == "" {
				continue
			}
			_, err := deps.DeleteProductVariantImage(ctx, &productvariantimagepb.DeleteProductVariantImageRequest{
				Data: &productvariantimagepb.ProductVariantImage{Id: imgID},
			})
			if err != nil {
				log.Printf("Failed to delete product_variant_image %s: %v", imgID, err)
			}
		}

		// Return updated gallery partial
		images := loadVariantImages(ctx, deps, vid)
		data := &VariantPageData{
			ProductID: id,
			VariantID: vid,
			Images:    images,
		}
		return view.OK("variant-image-gallery", data)
	})
}

// extensionFromContentType returns a file extension for the given MIME content type.
func extensionFromContentType(ct string) string {
	switch ct {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "image/gif":
		return ".gif"
	default:
		exts, _ := mime.ExtensionsByType(ct)
		if len(exts) > 0 {
			return exts[0]
		}
		return ""
	}
}
