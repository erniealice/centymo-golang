package detail

import (
	"context"
	"log"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	productattributepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product_attribute"
)

// AvailableAttribute represents a global attribute that can be assigned to a product.
type AvailableAttribute struct {
	ID   string
	Name string
	Code string
}

// AttributeFormLabels holds labels for the attribute drawer form template.
type AttributeFormLabels struct {
	Title        string
	DefaultValue string
}

// AttributeFormData is the template data for the attribute drawer form.
type AttributeFormData struct {
	FormAction          string
	ProductID           string
	Labels              AttributeFormLabels
	CommonLabels        any
	AvailableAttributes []AvailableAttribute
}

// AttributeDeps holds dependencies for attribute action handlers.
type AttributeDeps struct {
	Routes       centymo.ProductRoutes
	DB           centymo.DataSource
	Labels       centymo.ProductLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Typed proto funcs for global attributes
	ListAttributes func(ctx context.Context, req *commonpb.ListAttributesRequest) (*commonpb.ListAttributesResponse, error)
	ReadAttribute  func(ctx context.Context, req *commonpb.ReadAttributeRequest) (*commonpb.ReadAttributeResponse, error)

	// Typed proto funcs for product_attribute
	ListProductAttributes  func(ctx context.Context, req *productattributepb.ListProductAttributesRequest) (*productattributepb.ListProductAttributesResponse, error)
	CreateProductAttribute func(ctx context.Context, req *productattributepb.CreateProductAttributeRequest) (*productattributepb.CreateProductAttributeResponse, error)
	DeleteProductAttribute func(ctx context.Context, req *productattributepb.DeleteProductAttributeRequest) (*productattributepb.DeleteProductAttributeResponse, error)
}

// NewAttributeAssignView creates the attribute assign action (GET = form, POST = create).
func NewAttributeAssignView(deps *AttributeDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			// Load all global attributes
			var available []AvailableAttribute
			if deps.ListAttributes != nil {
				attrResp, err := deps.ListAttributes(ctx, &commonpb.ListAttributesRequest{})
				if err != nil {
					log.Printf("Failed to list global attributes: %v", err)
				} else {
					// Load already-assigned attributes for this product
					assigned := make(map[string]bool)
					if deps.ListProductAttributes != nil {
						paResp, err := deps.ListProductAttributes(ctx, &productattributepb.ListProductAttributesRequest{})
						if err == nil {
							for _, pa := range paResp.GetData() {
								if pa.GetProductId() == productID {
									if aid := pa.GetAttributeId(); aid != "" {
										assigned[aid] = true
									}
								}
							}
						}
					}

					// Filter to unassigned only
					for _, a := range attrResp.GetData() {
						aid := a.GetId()
						if assigned[aid] {
							continue
						}
						name := a.GetName()
						code := a.GetCode()
						available = append(available, AvailableAttribute{
							ID:   aid,
							Name: name,
							Code: code,
						})
					}
				}
			}

			l := deps.Labels
			return view.OK("attribute-drawer-form", &AttributeFormData{
				FormAction: route.ResolveURL(deps.Routes.AttributeAssignURL, "id", productID),
				ProductID:  productID,
				Labels: AttributeFormLabels{
					Title:        l.Attribute.Title,
					DefaultValue: l.Attribute.DefaultValue,
				},
				CommonLabels:        nil, // injected by ViewAdapter
				AvailableAttributes: available,
			})
		}

		// POST — create product_attribute record
		if err := viewCtx.Request.ParseForm(); err != nil {
			return HtmxError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		attributeID := r.FormValue("attribute_id")
		if attributeID == "" {
			return HtmxError(deps.Labels.Errors.FieldRequired)
		}

		// Look up attribute name and code for denormalized storage
		attrName := ""
		attrCode := ""
		if deps.ReadAttribute != nil {
			readResp, err := deps.ReadAttribute(ctx, &commonpb.ReadAttributeRequest{
				Data: &commonpb.Attribute{Id: attributeID},
			})
			if err == nil && len(readResp.GetData()) > 0 {
				attr := readResp.GetData()[0]
				attrName = attr.GetName()
				attrCode = attr.GetCode()
			}
		}

		defaultValue := r.FormValue("default_value")

		_, err := deps.CreateProductAttribute(ctx, &productattributepb.CreateProductAttributeRequest{
			Data: &productattributepb.ProductAttribute{
				ProductId:   productID,
				AttributeId: attributeID,
				Value:       defaultValue,
			},
		})
		if err != nil {
			log.Printf("Failed to create product attribute: %v", err)
			return HtmxError(err.Error())
		}

		// attrName and attrCode are looked up for logging/debugging but stored via the Attribute relation
		_ = attrName
		_ = attrCode

		return HtmxSuccess("product-attributes-table")
	})
}

// NewAttributeRemoveView creates the attribute remove action (POST only, with dialog confirmation).
func NewAttributeRemoveView(deps *AttributeDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return HtmxError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeleteProductAttribute(ctx, &productattributepb.DeleteProductAttributeRequest{
			Data: &productattributepb.ProductAttribute{Id: id},
		})
		if err != nil {
			log.Printf("Failed to delete product attribute %s: %v", id, err)
			return HtmxError(err.Error())
		}

		return HtmxSuccess("product-attributes-table")
	})
}

// strPtr returns a pointer to the string, or nil if empty.
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
