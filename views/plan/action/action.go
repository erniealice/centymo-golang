package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	centymo "github.com/erniealice/centymo-golang"

	clientpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/client"
	jobtemplatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/operation/job_template"
	planpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/subscription/plan"
)

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Name            string
	NamePlaceholder string
	Description     string
	DescPlaceholder string
	Active          string

	// Field-level info text surfaced via an info button beside each label.
	NameInfo        string
	DescriptionInfo string
	ActiveInfo      string

	// 2026-04-27 plan-client-scope plan §6.2 / §6.6 — Client picker.
	Client                  string
	ClientHelp              string
	ClientPlaceholder       string
	ClientSearchPlaceholder string
	ClientNoResults         string
	ClientLockedTooltip     string
	ClientForLabel          string // "For {{.ClientName}}" — read-only badge in client-context entry
	ClientInfo              string

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — JobTemplate select.
	JobTemplate     string
	JobTemplateNone string
	JobTemplateHint string
}

// ClientFieldMode selects how the Client field renders on the drawer per
// plan §6.6:
//   - "picker"   → standard auto-complete (workspace add).
//   - "readonly" → read-only badge "For {ClientName}" (client-context entry,
//     ?context=client&client_id=...).
//   - "locked"   → read-only badge with the lock tooltip (Plan has active
//     subscriptions and client_id is reference-checker locked).
type ClientFieldMode string

const (
	ClientFieldModePicker   ClientFieldMode = "picker"
	ClientFieldModeReadonly ClientFieldMode = "readonly"
	ClientFieldModeLocked   ClientFieldMode = "locked"
)

// JobTemplateOption is a {value, label} pair for the JobTemplate select.
type JobTemplateOption struct {
	Value string
	Label string
}

// FormData is the template data for the plan drawer form.
type FormData struct {
	FormAction   string
	IsEdit       bool
	ID           string
	Name         string
	Description  string
	Active       bool

	// 2026-04-27 plan-client-scope plan §6.2 / §6.6.
	ClientFieldMode      ClientFieldMode
	ClientID             string             // existing or pre-filled client_id
	ClientLabel          string             // display name for the chosen client
	ClientOptions        []map[string]any   // optgroup-flattened options for the picker
	SearchClientURL      string             // auto-complete search endpoint

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — Plan.job_template_id
	// assignment. JobTemplateID is the currently-assigned id (empty on add /
	// when unset); JobTemplateOptions enumerates active JobTemplates for the
	// drawer's <select>.
	JobTemplateID      string
	JobTemplateOptions []JobTemplateOption

	Labels       FormLabels
	CommonLabels any
}

// Deps holds dependencies for plan action handlers.
type Deps struct {
	Routes centymo.PlanRoutes
	Labels centymo.PlanLabels

	// Typed plan operations
	CreatePlan    func(ctx context.Context, req *planpb.CreatePlanRequest) (*planpb.CreatePlanResponse, error)
	ReadPlan      func(ctx context.Context, req *planpb.ReadPlanRequest) (*planpb.ReadPlanResponse, error)
	UpdatePlan    func(ctx context.Context, req *planpb.UpdatePlanRequest) (*planpb.UpdatePlanResponse, error)
	DeletePlan    func(ctx context.Context, req *planpb.DeletePlanRequest) (*planpb.DeletePlanResponse, error)
	SetPlanActive func(ctx context.Context, id string, active bool) error

	// 2026-04-27 plan-client-scope plan §6.2.
	// Client picker support — list / search clients for the auto-complete and
	// resolve a single client's display name when editing.
	ListClients         func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)
	SearchClientsByName func(ctx context.Context, req *clientpb.SearchClientsByNameRequest) (*clientpb.SearchClientsByNameResponse, error)
	SearchClientsURL    string

	// GetPlanClientScopeLockedIDs is the espyna reference checker that
	// returns plan IDs whose client_id is locked because at least one of
	// their PricePlans is attached to an active subscription. Optional —
	// when nil, lock state never flips.
	GetPlanClientScopeLockedIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// 2026-04-29 auto-spawn-jobs-from-subscription plan §5 — JobTemplate select
	// on the Plan drawer. Optional: when nil the drawer hides the select and
	// Plan.job_template_id is left untouched.
	ListJobTemplates func(ctx context.Context, req *jobtemplatepb.ListJobTemplatesRequest) (*jobtemplatepb.ListJobTemplatesResponse, error)
}

func formLabels(l centymo.PlanLabels) FormLabels {
	return FormLabels{
		Name:            l.Form.Name,
		NamePlaceholder: l.Form.NamePlaceholder,
		Description:     l.Form.Description,
		DescPlaceholder: l.Form.DescPlaceholder,
		Active:          l.Form.Active,
		// Info fields sourced from centymo.PlanFormLabels (populated from lyngua JSON + defaults).
		NameInfo:        l.Form.NameInfo,
		DescriptionInfo: l.Form.DescriptionInfo,
		ActiveInfo:      l.Form.ActiveInfo,
		// 2026-04-27 plan-client-scope plan §6.2.
		Client:                  l.Form.ClientLabel,
		ClientHelp:              l.Form.ClientHelp,
		ClientPlaceholder:       l.Form.ClientPlaceholder,
		ClientSearchPlaceholder: l.Form.ClientSearchPlaceholder,
		ClientNoResults:         l.Form.ClientNoResults,
		ClientLockedTooltip:     l.Form.ClientLockedTooltip,
		ClientForLabel:          l.Form.ClientForLabel,
		ClientInfo:              l.Form.ClientInfo,
		// 2026-04-29 auto-spawn-jobs-from-subscription plan §5.
		JobTemplate:     l.Form.JobTemplate,
		JobTemplateNone: l.Form.JobTemplateNone,
		JobTemplateHint: l.Form.JobTemplateHint,
	}
}

// loadJobTemplateOptions fetches active JobTemplates for the Plan drawer's
// JobTemplate select. Returns nil when the dep is unwired so the template
// can hide the field gracefully.
func loadJobTemplateOptions(ctx context.Context, listJobTemplates func(ctx context.Context, req *jobtemplatepb.ListJobTemplatesRequest) (*jobtemplatepb.ListJobTemplatesResponse, error)) []JobTemplateOption {
	if listJobTemplates == nil {
		return nil
	}
	resp, err := listJobTemplates(ctx, &jobtemplatepb.ListJobTemplatesRequest{})
	if err != nil {
		log.Printf("Failed to load job templates for plan drawer: %v", err)
		return nil
	}
	opts := make([]JobTemplateOption, 0, len(resp.GetData()))
	for _, t := range resp.GetData() {
		if t == nil {
			continue
		}
		// Skip inactive templates so operators only assign live workflows.
		if !t.GetActive() {
			continue
		}
		label := t.GetName()
		if label == "" {
			label = t.GetId()
		}
		opts = append(opts, JobTemplateOption{
			Value: t.GetId(),
			Label: label,
		})
	}
	return opts
}

// loadClientOptions fetches the workspace's clients and converts them into
// the auto-complete option shape ({Value, Label, Selected}). Returns nil
// when the dep is unwired.
func loadClientOptions(ctx context.Context, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error), selectedID string) []map[string]any {
	if listClients == nil {
		return nil
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		log.Printf("Failed to load clients for plan drawer: %v", err)
		return nil
	}
	opts := make([]map[string]any, 0, len(resp.GetData()))
	for _, c := range resp.GetData() {
		label := c.GetName()
		if label == "" {
			if u := c.GetUser(); u != nil {
				label = strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			}
		}
		if label == "" {
			label = c.GetId()
		}
		opts = append(opts, map[string]any{
			"Value":    c.GetId(),
			"Label":    label,
			"Selected": c.GetId() == selectedID,
		})
	}
	return opts
}

// resolveClientLabel finds the display name for a single client_id. Empty
// when listClients is unwired or the lookup misses.
func resolveClientLabel(ctx context.Context, clientID string, listClients func(ctx context.Context, req *clientpb.ListClientsRequest) (*clientpb.ListClientsResponse, error)) string {
	if clientID == "" || listClients == nil {
		return ""
	}
	resp, err := listClients(ctx, &clientpb.ListClientsRequest{})
	if err != nil {
		return clientID
	}
	for _, c := range resp.GetData() {
		if c.GetId() != clientID {
			continue
		}
		if name := c.GetName(); name != "" {
			return name
		}
		if u := c.GetUser(); u != nil {
			full := strings.TrimSpace(u.GetFirstName() + " " + u.GetLastName())
			if full != "" {
				return full
			}
		}
		return clientID
	}
	return clientID
}

// NewAddAction creates the plan add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan", "create") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			// Plan §6.6 — when opened with ?context=client&client_id=...
			// the Client field renders read-only (with a hidden input).
			ctxParam := viewCtx.Request.URL.Query().Get("context")
			pinnedClientID := viewCtx.Request.URL.Query().Get("client_id")
			fieldMode := ClientFieldModePicker
			clientLabel := ""
			if ctxParam == "client" && pinnedClientID != "" {
				fieldMode = ClientFieldModeReadonly
				clientLabel = resolveClientLabel(ctx, pinnedClientID, deps.ListClients)
			}
			return view.OK("plan-drawer-form", &FormData{
				FormAction:         deps.Routes.AddURL,
				Active:             true,
				ClientFieldMode:    fieldMode,
				ClientID:           pinnedClientID,
				ClientLabel:        clientLabel,
				ClientOptions:      loadClientOptions(ctx, deps.ListClients, pinnedClientID),
				SearchClientURL:    deps.SearchClientsURL,
				JobTemplateOptions: loadJobTemplateOptions(ctx, deps.ListJobTemplates),
				Labels:             formLabels(deps.Labels),
				CommonLabels:       nil, // injected by ViewAdapter
			})
		}

		// POST — create plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		clientID := strings.TrimSpace(r.FormValue("client_id"))
		jobTemplateID := strings.TrimSpace(r.FormValue("job_template_id"))

		planData := &planpb.Plan{
			Name:        r.FormValue("name"),
			Description: strPtr(r.FormValue("description")),
			Active:      active,
		}
		if clientID != "" {
			planData.ClientId = strPtr(clientID)
		}
		if jobTemplateID != "" {
			planData.JobTemplateId = strPtr(jobTemplateID)
		}
		resp, err := deps.CreatePlan(ctx, &planpb.CreatePlanRequest{
			Data: planData,
		})
		if err != nil {
			log.Printf("Failed to create plan: %v", err)
			return centymo.HTMXError(err.Error())
		}

		newID := ""
		if respData := resp.GetData(); len(respData) > 0 {
			newID = respData[0].GetId()
		}
		// Proto3 omits active=false on Create; force-sync via raw update when unchecked.
		if !active && newID != "" && deps.SetPlanActive != nil {
			if err := deps.SetPlanActive(ctx, newID, false); err != nil {
				log.Printf("Failed to set plan inactive after create %s: %v", newID, err)
			}
		}

		return centymo.HTMXSuccess("plans-table")
	})
}

// NewEditAction creates the plan edit action (GET = form, POST = update).
// When the GET request includes ?clone=1, the handler returns the drawer form
// pre-populated from the source record but wired to AddURL (submission creates
// a new plan) with " (Copy)" appended to the name.
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		id := viewCtx.Request.PathValue("id")
		isClone := viewCtx.Request.Method == http.MethodGet && viewCtx.Request.URL.Query().Get("clone") == "1"

		requiredAction := "update"
		if isClone {
			requiredAction = "create"
		}
		if !perms.Can("plan", requiredAction) {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		if viewCtx.Request.Method == http.MethodGet {
			readResp, err := deps.ReadPlan(ctx, &planpb.ReadPlanRequest{
				Data: &planpb.Plan{Id: &id},
			})
			if err != nil {
				log.Printf("Failed to read plan %s: %v", id, err)
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			readData := readResp.GetData()
			if len(readData) == 0 {
				return centymo.HTMXError(deps.Labels.Errors.NotFound)
			}
			record := readData[0]

			name := record.GetName()
			formAction := route.ResolveURL(deps.Routes.EditURL, "id", id)
			formID := id
			if isClone {
				name = strings.TrimSpace(name) + viewCtx.T("actions.copySuffix")
				formAction = deps.Routes.AddURL
				formID = ""
			}

			// 2026-04-27 plan-client-scope plan §3.1 / §6.2 — Client field
			// renders as a picker by default; flips to read-only-with-tooltip
			// when the reference checker reports the plan as locked.
			clientID := record.GetClientId()
			clientLabel := resolveClientLabel(ctx, clientID, deps.ListClients)
			fieldMode := ClientFieldModePicker
			if isClone {
				// Cloning starts fresh — the new plan can be assigned to any
				// client (or none). Skip the lock check.
			} else if deps.GetPlanClientScopeLockedIDs != nil {
				if locked, _ := deps.GetPlanClientScopeLockedIDs(ctx, []string{id}); locked[id] {
					fieldMode = ClientFieldModeLocked
				}
			}
			// Honor ?context=client&client_id=... overrides on edit (e.g. opened
			// from the client detail Packages tab) — a context-pinned drawer
			// always renders read-only regardless of lock state, since the
			// client identity is implied by the entry point.
			if viewCtx.Request.URL.Query().Get("context") == "client" {
				if pinned := viewCtx.Request.URL.Query().Get("client_id"); pinned != "" {
					clientID = pinned
					clientLabel = resolveClientLabel(ctx, pinned, deps.ListClients)
					fieldMode = ClientFieldModeReadonly
				}
			}

			return view.OK("plan-drawer-form", &FormData{
				FormAction:         formAction,
				IsEdit:             !isClone,
				ID:                 formID,
				Name:               name,
				Description:        record.GetDescription(),
				Active:             record.GetActive(),
				ClientFieldMode:    fieldMode,
				ClientID:           clientID,
				ClientLabel:        clientLabel,
				ClientOptions:      loadClientOptions(ctx, deps.ListClients, clientID),
				SearchClientURL:    deps.SearchClientsURL,
				JobTemplateID:      record.GetJobTemplateId(),
				JobTemplateOptions: loadJobTemplateOptions(ctx, deps.ListJobTemplates),
				Labels:             formLabels(deps.Labels),
				CommonLabels:       nil, // injected by ViewAdapter
			})
		}

		// POST — update plan
		if err := viewCtx.Request.ParseForm(); err != nil {
			return centymo.HTMXError(deps.Labels.Errors.InvalidFormData)
		}

		r := viewCtx.Request
		active := r.FormValue("active") == "true"
		clientID := strings.TrimSpace(r.FormValue("client_id"))
		jobTemplateID := strings.TrimSpace(r.FormValue("job_template_id"))

		updateData := &planpb.Plan{
			Id:          &id,
			Name:        r.FormValue("name"),
			Description: strPtr(r.FormValue("description")),
		}
		// Always send the client_id (empty → nil so postgres writes NULL,
		// otherwise the empty string trips plan_client_id_fkey) so the espyna
		// update use case can detect a master ↔ client_id transition and run
		// the §3.1 reference-checker guard.
		if clientID == "" {
			updateData.ClientId = nil
		} else {
			updateData.ClientId = &clientID
		}
		// Always send job_template_id (empty → nil so postgres writes NULL,
		// otherwise the empty string would trip plan_job_template_id_fkey).
		if jobTemplateID == "" {
			updateData.JobTemplateId = nil
		} else {
			updateData.JobTemplateId = &jobTemplateID
		}
		_, err := deps.UpdatePlan(ctx, &planpb.UpdatePlanRequest{
			Data: updateData,
		})
		if err != nil {
			log.Printf("Failed to update plan %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		// Proto3 omits active=false; sync the active state via raw update so the
		// toggle in the drawer actually persists deactivation.
		if deps.SetPlanActive != nil {
			if err := deps.SetPlanActive(ctx, id, active); err != nil {
				log.Printf("Failed to set plan active state %s: %v", id, err)
			}
		}

		// Close drawer and reload current page so the detail view reflects the update.
		return view.ViewResult{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"HX-Trigger": `{"formSuccess":true}`,
				"HX-Refresh": "true",
			},
		}
	})
}

// NewDeleteAction creates the plan delete action (POST only).
// The row ID comes via query param (?id=xxx) appended by table-actions.js.
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}

		_, err := deps.DeletePlan(ctx, &planpb.DeletePlanRequest{
			Data: &planpb.Plan{Id: &id},
		})
		if err != nil {
			log.Printf("Failed to delete plan %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plans-table")
	})
}

// NewBulkDeleteAction creates the plan bulk delete action (POST only).
// Selected IDs come as multiple "id" form fields.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan", "delete") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}

		for _, id := range ids {
			if _, err := deps.DeletePlan(ctx, &planpb.DeletePlanRequest{
				Data: &planpb.Plan{Id: &id},
			}); err != nil {
				log.Printf("Failed to delete plan %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("plans-table")
	})
}

// NewSetStatusAction creates the plan activate/deactivate action (POST only).
// Expects query params: ?id={planId}&status={active|inactive}
//
// Uses SetPlanActive (raw map update) instead of protobuf because
// proto3's protojson omits bool fields with value false, which means
// deactivation (active=false) would silently be skipped.
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return centymo.HTMXError(deps.Labels.Errors.IDRequired)
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		if err := deps.SetPlanActive(ctx, id, targetStatus == "active"); err != nil {
			log.Printf("Failed to update plan status %s: %v", id, err)
			return centymo.HTMXError(err.Error())
		}

		return centymo.HTMXSuccess("plans-table")
	})
}

// NewBulkSetStatusAction creates the plan bulk activate/deactivate action (POST only).
// Selected IDs come as multiple "id" form fields; target status from "target_status" field.
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("plan", "update") {
			return centymo.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return centymo.HTMXError(deps.Labels.Errors.NoIDsProvided)
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return centymo.HTMXError(deps.Labels.Errors.InvalidStatus)
		}

		active := targetStatus == "active"

		for _, id := range ids {
			if err := deps.SetPlanActive(ctx, id, active); err != nil {
				log.Printf("Failed to update plan status %s: %v", id, err)
			}
		}

		return centymo.HTMXSuccess("plans-table")
	})
}

// strPtr returns a pointer to a string.
func strPtr(s string) *string {
	return &s
}
