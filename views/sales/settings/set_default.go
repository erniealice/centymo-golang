package settings

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/document_template"
)

// NewSetDefaultAction creates the POST handler for setting a template as the default.
// Route: POST /action/sales/settings/templates/set-default/{id}
func NewSetDefaultAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if deps.UpdateDocumentTemplate == nil || deps.ListDocumentTemplates == nil {
			return view.Error(fmt.Errorf("template update not configured"))
		}

		targetID := viewCtx.Request.PathValue("id")
		if targetID == "" {
			return view.Error(fmt.Errorf("template id is required"))
		}

		// List all templates to find and unset any existing default
		resp, err := deps.ListDocumentTemplates(ctx, &documenttemplatepb.ListDocumentTemplatesRequest{})
		if err != nil {
			log.Printf("Failed to list document templates: %v", err)
			return view.Error(fmt.Errorf("failed to list templates: %w", err))
		}

		falseVal := false
		trueVal := true

		for _, t := range resp.GetData() {
			if !t.GetActive() || t.GetDocumentPurpose() != "invoice" {
				continue
			}

			if t.GetIsDefault() && t.GetId() != targetID {
				// Unset existing default
				_, err := deps.UpdateDocumentTemplate(ctx, &documenttemplatepb.UpdateDocumentTemplateRequest{
					Data: &documenttemplatepb.DocumentTemplate{
						Id:        t.GetId(),
						IsDefault: &falseVal,
					},
				})
				if err != nil {
					log.Printf("Failed to unset default on template %s: %v", t.GetId(), err)
				}
			}
		}

		// Set the target template as default
		_, err = deps.UpdateDocumentTemplate(ctx, &documenttemplatepb.UpdateDocumentTemplateRequest{
			Data: &documenttemplatepb.DocumentTemplate{
				Id:        targetID,
				IsDefault: &trueVal,
			},
		})
		if err != nil {
			log.Printf("Failed to set default on template %s: %v", targetID, err)
			return view.Error(fmt.Errorf("failed to set default template: %w", err))
		}

		// Redirect back to settings page to show updated list
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Redirect": deps.Routes.SettingsTemplatesURL,
			},
		}
	})
}
