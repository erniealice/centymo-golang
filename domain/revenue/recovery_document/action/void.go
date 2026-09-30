// Package action holds the recovery document write handlers (void).
//
// The handler re-checks its permission first (layer 2, fail closed); the use
// case remains the authoritative gate (strict verb; refusals has_applications,
// has_credit_notes, already_void, void_reason_required).
package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	rd "github.com/erniealice/centymo-golang/domain/revenue/recovery_document"
)

// Deps holds the action handler dependencies.
type Deps struct {
	Routes       rd.Routes
	Labels       rd.Labels
	CommonLabels pyeza.CommonLabels
	UseCases     *rd.UseCases
}

// VoidData is the template data for the void drawer.
type VoidData struct {
	FormAction  string
	WorkspaceID string // injected by the ViewAdapter (action_workspace_guard)
	ID          string
	Number      string
	Labels      rd.Labels

	CommonLabels any
}

func refuse(deps *Deps, err error) view.ViewResult {
	kind := rd.ErrorKind(err)
	if kind == rd.ErrKindPermissionDenied {
		return view.HTMXError(deps.CommonLabels.Errors.PermissionDenied)
	}
	if kind == rd.ErrKindUnknown {
		log.Printf("recovery_document void: %v", err)
	}
	if m := deps.Labels.ErrorMessage(kind); m != "" {
		return view.HTMXError(m)
	}
	return view.HTMXError(deps.CommonLabels.Errors.General)
}

// NewVoidAction handles GET (drawer) and POST (VoidRecoveryDocument).
func NewVoidAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("recovery_document", "void") {
			return view.HTMXError(deps.CommonLabels.Errors.PermissionDenied)
		}
		uc := deps.UseCases
		if uc == nil || uc.VoidRecoveryDocument == nil || uc.ReadRecoveryDocument == nil {
			return view.ViewResult{
				StatusCode: http.StatusServiceUnavailable,
				Headers:    map[string]string{"HX-Error-Message": deps.CommonLabels.Errors.General},
			}
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.HTMXError(deps.CommonLabels.Errors.InvalidFormData)
		}
		formAction := route.ResolveURL(deps.Routes.VoidURL, "id", id)

		if viewCtx.Request.Method == http.MethodGet {
			resp, err := uc.ReadRecoveryDocument(ctx, &recoverydocumentpb.ReadRecoveryDocumentRequest{
				Data: &recoverydocumentpb.RecoveryDocument{Id: id},
			})
			if err != nil || len(resp.GetData()) == 0 {
				if err != nil {
					return refuse(deps, err)
				}
				return view.HTMXError(deps.Labels.Errors.NotFound)
			}
			return view.OK("recovery-document-void-drawer", &VoidData{
				FormAction: formAction, ID: id, Number: resp.GetData()[0].GetDocumentNumber(), Labels: deps.Labels,
			})
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.CommonLabels.Errors.InvalidFormData)
		}
		reason := strings.TrimSpace(viewCtx.Request.FormValue("reason"))
		if reason == "" {
			return view.HTMXError(deps.Labels.Errors.VoidReasonRequired)
		}
		if _, err := uc.VoidRecoveryDocument(ctx, &recoverydocumentpb.VoidRecoveryDocumentRequest{
			RecoveryDocumentId: id, Reason: reason,
		}); err != nil {
			return refuse(deps, err)
		}
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger":  `{"formSuccess":true}`,
				"HX-Redirect": route.ResolveURL(deps.Routes.DetailURL, "id", id),
			},
		}
	})
}
