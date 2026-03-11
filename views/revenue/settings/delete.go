package settings

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	documenttemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/document_template"
)

// NewDeleteAction creates the POST handler for deleting invoice templates.
// Route: POST /action/sales/settings/templates/delete
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		if deps.DeleteDocumentTemplate == nil {
			return view.Error(fmt.Errorf("template delete not configured"))
		}

		viewCtx.Request.ParseForm()
		templateID := viewCtx.Request.FormValue("template_id")
		if templateID == "" {
			return view.Error(fmt.Errorf("template_id is required"))
		}

		_, err := deps.DeleteDocumentTemplate(ctx, &documenttemplatepb.DeleteDocumentTemplateRequest{
			Data: &documenttemplatepb.DocumentTemplate{Id: templateID},
		})
		if err != nil {
			log.Printf("Failed to delete document template %s: %v", templateID, err)
			return view.Error(fmt.Errorf("failed to delete template: %w", err))
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
