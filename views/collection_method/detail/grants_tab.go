package detail

// grants_tab.go — Stage 3 (treasury-domain-rebuild) Grants tab handler for
// the Collection Method detail page (pages.md §B-5 tab #3).
//
// Grants bind a client to a CollectionMethod TEMPLATE. They are CONFIG, never
// an EVENT (Q6 LOCKED): grants do not mutate — the only state change is
// ACTIVE → REVOKED. Therefore there is NO edit action; only bulk_grant (create
// many at once) and revoke are provided.
//
// All closures are nil-safe; when the espyna use cases are absent the tab
// renders an empty state rather than panicking.

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	centymo "github.com/erniealice/centymo-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	grantpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/collection_method_grant"
)

// GrantTabDeps holds the closures needed by the Grants tab.
// All closures are optional (nil-safe): the tab renders an empty/disabled
// state when a closure is absent rather than panicking.
type GrantTabDeps struct {
	Routes       centymo.CollectionMethodRoutes
	Labels       centymo.CollectionMethodLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListCollectionMethodGrants      func(ctx context.Context, req *grantpb.ListCollectionMethodGrantsRequest) (*grantpb.ListCollectionMethodGrantsResponse, error)
	CreateCollectionMethodGrant     func(ctx context.Context, req *grantpb.CreateCollectionMethodGrantRequest) (*grantpb.CreateCollectionMethodGrantResponse, error)
	RevokeCollectionMethodGrant     func(ctx context.Context, req *grantpb.RevokeCollectionMethodGrantRequest) (*grantpb.RevokeCollectionMethodGrantResponse, error)
	BulkGrantCollectionMethodGrants func(ctx context.Context, req *grantpb.BulkGrantCollectionMethodGrantsRequest) (*grantpb.BulkGrantCollectionMethodGrantsResponse, error)
}

// GrantRow is a flat representation of one grant for the template.
type GrantRow struct {
	ID           string
	Subject      string
	ClientID     string
	Status       string
	GrantedBy    string
	RevokedBy    string
	RevokeReason string

	RevokeURL string
}

// GrantBulkFormData holds the data for the bulk-grant drawer-form.
type GrantBulkFormData struct {
	Labels       centymo.CollectionMethodGrantLabels
	CommonLabels pyeza.CommonLabels

	MethodID   string
	FormAction string
	ClientIDs  string // pre-filled on re-render after error
}

// GrantTabData is passed to the grants-tab template.
type GrantTabData struct {
	Labels       centymo.CollectionMethodGrantLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	MethodID     string
	Grants       []GrantRow
	BulkGrantURL string
}

// grantTableID is the HTMX target used by action success responses.
const grantTableID = "grants-table"

// NewGrantTabView returns the tab-content partial for the Grants tab.
// The outer detail page delegates here when activeTab == "grants".
func NewGrantTabView(deps *GrantTabDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method_grant", "list") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method_grant:list"))
		}

		methodID := viewCtx.Request.PathValue("id")
		if methodID == "" {
			return view.Error(fmt.Errorf("missing method id"))
		}

		grants, err := loadGrants(ctx, deps, methodID)
		if err != nil {
			log.Printf("ListCollectionMethodGrants method=%s: %v", methodID, err)
			return view.Error(fmt.Errorf("failed to load grants: %w", err))
		}

		bulkGrantURL := ""
		if perms.Can("collection_method_grant", "bulk_grant") {
			bulkGrantURL = route.ResolveURL(deps.Routes.GrantBulkGrantURL, "method_id", methodID)
		}

		td := &GrantTabData{
			Labels:       deps.Labels.Grant,
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			MethodID:     methodID,
			Grants:       grants,
			BulkGrantURL: bulkGrantURL,
		}
		return view.OK("collection-method-grants-tab", td)
	})
}

// NewGrantBulkGrantAction handles GET + POST for the bulk-grant drawer-form.
func NewGrantBulkGrantAction(deps *GrantTabDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method_grant", "bulk_grant") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method_grant:bulk_grant"))
		}
		methodID := viewCtx.Request.PathValue("method_id")
		if methodID == "" {
			return view.Error(fmt.Errorf("missing method_id"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			fd := &GrantBulkFormData{
				Labels:       deps.Labels.Grant,
				CommonLabels: deps.CommonLabels,
				MethodID:     methodID,
				FormAction:   route.ResolveURL(deps.Routes.GrantBulkGrantURL, "method_id", methodID),
			}
			return view.OK("collection-method-grant-bulk-drawer-form", fd)
		}

		if deps.BulkGrantCollectionMethodGrants == nil {
			return centymo.HTMXError("Bulk grant is not wired yet.")
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}
		raw := viewCtx.Request.FormValue("client_ids")
		lines := splitClientIDs(raw)
		if len(lines) == 0 {
			return centymo.HTMXError("No client IDs provided.")
		}

		var grants []*grantpb.CollectionMethodGrant
		for _, cid := range lines {
			grants = append(grants, &grantpb.CollectionMethodGrant{
				CollectionMethodId: methodID,
				ClientId:           cid,
				Subject:            grantpb.CollectionMethodGrantSubject_COLLECTION_METHOD_GRANT_SUBJECT_CLIENT,
			})
		}
		if _, err := deps.BulkGrantCollectionMethodGrants(ctx, &grantpb.BulkGrantCollectionMethodGrantsRequest{
			Data: grants,
		}); err != nil {
			log.Printf("BulkGrantCollectionMethodGrants method=%s: %v", methodID, err)
			return view.Error(fmt.Errorf("failed to bulk grant: %w", err))
		}
		return centymo.HTMXSuccess(grantTableID)
	})
}

// NewGrantRevokeAction handles POST for the revoke action.
func NewGrantRevokeAction(deps *GrantTabDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("collection_method_grant", "revoke") {
			return centymo.HTMXError(fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "collection_method_grant:revoke"))
		}
		if viewCtx.Request.Method != http.MethodPost {
			return view.Error(fmt.Errorf("method not allowed"))
		}
		methodID := viewCtx.Request.PathValue("method_id")
		if methodID == "" {
			return view.Error(fmt.Errorf("missing method_id"))
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.Error(fmt.Errorf("parse form: %w", err))
		}
		id := viewCtx.Request.FormValue("id")
		if id == "" {
			return view.Error(fmt.Errorf("missing grant id"))
		}
		if deps.RevokeCollectionMethodGrant == nil {
			return centymo.HTMXError("Grant revoke is not wired yet.")
		}
		if _, err := deps.RevokeCollectionMethodGrant(ctx, &grantpb.RevokeCollectionMethodGrantRequest{
			Data: &grantpb.CollectionMethodGrant{Id: id},
		}); err != nil {
			log.Printf("RevokeCollectionMethodGrant %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to revoke grant: %w", err))
		}
		return centymo.HTMXSuccess(grantTableID)
	})
}

// --- helpers -----------------------------------------------------------------

func loadGrants(ctx context.Context, deps *GrantTabDeps, methodID string) ([]GrantRow, error) {
	if deps.ListCollectionMethodGrants == nil {
		return nil, nil
	}
	resp, err := deps.ListCollectionMethodGrants(ctx, &grantpb.ListCollectionMethodGrantsRequest{})
	if err != nil {
		return nil, err
	}
	var rows []GrantRow
	for _, g := range resp.GetData() {
		revokeURL := route.ResolveURL(deps.Routes.GrantRevokeURL, "method_id", methodID)
		row := GrantRow{
			ID:           g.GetId(),
			Subject:      grantSubjectShort(g.GetSubject().String()),
			ClientID:     g.GetClientId(),
			Status:       grantStatusShort(g.GetStatus().String()),
			GrantedBy:    g.GetGrantedByUserId(),
			RevokedBy:    g.GetRevokedByUserId(),
			RevokeReason: g.GetRevokeReason(),
			RevokeURL:    revokeURL,
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// grantSubjectShort strips the long proto prefix from a subject enum string.
func grantSubjectShort(s string) string {
	const pfx = "COLLECTION_METHOD_GRANT_SUBJECT_"
	if len(s) > len(pfx) {
		return s[len(pfx):]
	}
	return s
}

// grantStatusShort strips the long proto prefix from a status enum string.
func grantStatusShort(s string) string {
	const pfx = "COLLECTION_METHOD_GRANT_STATUS_"
	if len(s) > len(pfx) {
		return s[len(pfx):]
	}
	return s
}

// splitClientIDs splits a newline/comma-separated string of client IDs into
// a deduplicated, trimmed slice.
func splitClientIDs(raw string) []string {
	// Normalise commas to newlines, then split on newlines.
	raw = strings.ReplaceAll(raw, ",", "\n")
	parts := strings.Split(raw, "\n")
	seen := make(map[string]struct{})
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}
