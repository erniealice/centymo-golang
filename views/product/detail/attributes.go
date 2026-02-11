package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"
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
	DB           centymo.DataSource
	Labels       centymo.ProductLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// NewAttributesTableView returns a view that renders only the attributes table (for HTMX refresh).
func NewAttributesTableView(deps *AttributeDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")

		detailDeps := &Deps{
			DB:          deps.DB,
			Labels:      deps.Labels,
			TableLabels: deps.TableLabels,
		}

		tableConfig := buildAttributesTable(ctx, detailDeps, productID)
		return view.OK("table-card", tableConfig)
	})
}

// NewAttributeAssignView creates the attribute assign action (GET = form, POST = create).
func NewAttributeAssignView(deps *AttributeDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		productID := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			// Load all global attributes
			var available []AvailableAttribute
			if deps.DB != nil {
				allAttrs, err := deps.DB.ListSimple(ctx, "attribute")
				if err != nil {
					log.Printf("Failed to list global attributes: %v", err)
				} else {
					// Load already-assigned attributes for this product
					assigned := make(map[string]bool)
					productAttrs, err := deps.DB.ListSimple(ctx, "product_attribute")
					if err == nil {
						for _, pa := range productAttrs {
							if pid, _ := pa["product_id"].(string); pid == productID {
								if aid, _ := pa["attribute_id"].(string); aid != "" {
									assigned[aid] = true
								}
							}
						}
					}

					// Filter to unassigned only
					for _, a := range allAttrs {
						aid, _ := a["id"].(string)
						if assigned[aid] {
							continue
						}
						name, _ := a["name"].(string)
						code, _ := a["code"].(string)
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
				FormAction: fmt.Sprintf("/action/products/detail/%s/attributes/assign", productID),
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
			return htmxError("Invalid form data")
		}

		r := viewCtx.Request
		attributeID := r.FormValue("attribute_id")
		if attributeID == "" {
			return htmxError("Please select an attribute")
		}

		active := r.FormValue("active") == "true"

		// Look up attribute name and code for denormalized storage
		attrName := ""
		attrCode := ""
		if deps.DB != nil {
			attr, err := deps.DB.Read(ctx, "attribute", attributeID)
			if err == nil {
				attrName, _ = attr["name"].(string)
				attrCode, _ = attr["code"].(string)
			}
		}

		data := map[string]any{
			"product_id":     productID,
			"attribute_id":   attributeID,
			"attribute_name": attrName,
			"attribute_code": attrCode,
			"default_value":  r.FormValue("default_value"),
			"active":         active,
		}

		_, err := deps.DB.Create(ctx, "product_attribute", data)
		if err != nil {
			log.Printf("Failed to create product attribute: %v", err)
			return htmxError("Failed to assign attribute")
		}

		return htmxSuccess("product-attributes-table")
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
			return htmxError("Product attribute ID is required")
		}

		err := deps.DB.Delete(ctx, "product_attribute", id)
		if err != nil {
			log.Printf("Failed to delete product attribute %s: %v", id, err)
			return htmxError("Failed to remove attribute")
		}

		return htmxSuccess("product-attributes-table")
	})
}
